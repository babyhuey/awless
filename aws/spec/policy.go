/* Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsspec

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/bootswithdefer/awless/logger"
)

type CreatePolicy struct {
	_           string `action:"create" entity:"policy" awsAPI:"iam" awsCall:"CreatePolicy" awsInput:"iam.CreatePolicyInput" awsOutput:"iam.CreatePolicyOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *iam.Client
	Name        *string   `awsName:"PolicyName" awsType:"awsstr" templateName:"name"`
	Effect      *string   `templateName:"effect"`
	Action      []*string `templateName:"action"`
	Resource    []*string `templateName:"resource"`
	Description *string   `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Document    *string   `awsName:"PolicyDocument" awsType:"awsstr"`
	Conditions  []*string `templateName:"conditions"`
}

func (cmd *CreatePolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("action"), params.Key("effect"), params.Key("name"), params.Key("resource"),
		params.Opt("conditions", "description"),
	))
}

func (cmd *CreatePolicy) BeforeRun(renv env.Running) error {
	stat, err := buildStatementFromParams(cmd.Effect, cmd.Resource, cmd.Action, cmd.Conditions)
	if err != nil {
		return err
	}
	policy := &policyBody{
		Version:   "2012-10-17",
		Statement: []*policyStatement{stat},
	}

	b, err := json.MarshalIndent(policy, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal policy document: %w", err)
	}
	cmd.Document = String(string(b))
	cmd.logger.ExtraVerbosef("policy document json:\n%s\n", string(b))
	return nil
}

func (cmd *CreatePolicy) ExtractResult(i any) string {
	return StringValue(i.(*iam.CreatePolicyOutput).Policy.Arn)
}

type UpdatePolicy struct {
	_              string `action:"update" entity:"policy" awsAPI:"iam" awsCall:"CreatePolicyVersion" awsInput:"iam.CreatePolicyVersionInput" awsOutput:"iam.CreatePolicyVersionOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *iam.Client
	Arn            *string   `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	Effect         *string   `templateName:"effect"`
	Action         []*string `templateName:"action"`
	Resource       []*string `templateName:"resource"`
	Conditions     []*string `templateName:"conditions"`
	Document       *string   `awsName:"PolicyDocument" awsType:"awsstr"`
	DefaultVersion *bool     `awsName:"SetAsDefault" awsType:"awsbool"`
}

func (cmd *UpdatePolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("action"), params.Key("arn"), params.Key("effect"), params.Key("resource"),
		params.Opt("conditions"),
	))
}

func (cmd *UpdatePolicy) BeforeRun(renv env.Running) error {
	document, err := cmd.getPolicyLastVersionDocument(cmd.Arn)
	if err != nil {
		return err
	}
	var defaultPolicyDocument *struct {
		Version    string             `json:",omitempty"`
		ID         string             `json:"Id,omitempty"`
		Statements []*json.RawMessage `json:"Statement,omitempty"`
	}

	if err = json.Unmarshal([]byte(document), &defaultPolicyDocument); err != nil {
		return err
	}
	stat, err := buildStatementFromParams(cmd.Effect, cmd.Resource, cmd.Action, cmd.Conditions)
	if err != nil {
		return err
	}

	var newStatement json.RawMessage
	if newStatement, err = json.Marshal(stat); err != nil {
		return err
	}
	defaultPolicyDocument.Statements = append(defaultPolicyDocument.Statements, &newStatement)

	b, err := json.MarshalIndent(defaultPolicyDocument, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal policy document: %w", err)
	}
	cmd.Document = String(string(b))
	cmd.DefaultVersion = aws.Bool(true)
	cmd.logger.ExtraVerbosef("policy document json:\n%s\n", string(b))
	return nil
}

func (cmd *UpdatePolicy) getPolicyLastVersionDocument(arn *string) (string, error) {
	listVersionsInput := &iam.ListPolicyVersionsInput{
		PolicyArn: arn,
	}
	listVersionsOut, err := cmd.api.ListPolicyVersions(context.Background(), listVersionsInput)
	if err != nil {
		return "", err
	}
	var defaultVersion *iamtypes.PolicyVersion
	for _, version := range listVersionsOut.Versions {
		if version.IsDefaultVersion {
			policyDetailInput := &iam.GetPolicyVersionInput{
				VersionId: version.VersionId,
				PolicyArn: arn,
			}
			var policyDetailOutput *iam.GetPolicyVersionOutput
			if policyDetailOutput, err = cmd.api.GetPolicyVersion(context.Background(), policyDetailInput); err != nil {
				return "", err
			}
			defaultVersion = policyDetailOutput.PolicyVersion
		}
	}
	if defaultVersion == nil {
		return "", fmt.Errorf("update policy: can not find default version for policy with arn '%s'", StringValue(arn))
	}
	document, err := url.QueryUnescape(aws.ToString(defaultVersion.Document))
	if err != nil {
		return "", fmt.Errorf("decoding policy document: %w", err)
	}
	return document, nil
}

type DeletePolicy struct {
	_           string `action:"delete" entity:"policy" awsAPI:"iam"  awsCall:"DeletePolicy" awsInput:"iam.DeletePolicyInput" awsOutput:"iam.DeletePolicyOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *iam.Client
	Arn         *string `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	AllVersions *bool   `templateName:"all-versions"`
}

func (cmd *DeletePolicy) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("arn"),
		params.Opt("all-versions"),
	))
}

func (cmd *DeletePolicy) BeforeRun(renv env.Running) error {
	if BoolValue(cmd.AllVersions) {
		list, err := cmd.api.ListPolicyVersions(renv.RequestContext(), &iam.ListPolicyVersionsInput{PolicyArn: cmd.Arn})
		if err != nil {
			return fmt.Errorf("list all policy versions: %w", err)
		}
		for _, v := range list.Versions {
			if !v.IsDefaultVersion {
				cmd.logger.Verbosef("deleting version '%s' of policy '%s'", aws.ToString(v.VersionId), StringValue(cmd.Arn))
				if _, err := cmd.api.DeletePolicyVersion(renv.RequestContext(), &iam.DeletePolicyVersionInput{PolicyArn: cmd.Arn, VersionId: v.VersionId}); err != nil {
					return fmt.Errorf("delete version %s: %w", aws.ToString(v.VersionId), err)
				}
			}
		}
	}
	return nil
}

type AttachPolicy struct {
	_       string `action:"attach" entity:"policy" awsAPI:"iam"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *iam.Client
	Arn     *string `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	User    *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
	Group   *string `awsName:"GroupName" awsType:"awsstr" templateName:"group"`
	Role    *string `awsName:"RoleName" awsType:"awsstr" templateName:"role"`
	Service *string `templateName:"service"`
	Access  *string `templateName:"access"`
}

func (cmd *AttachPolicy) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(
		params.OnlyOneOf(params.Key("user"), params.Key("role"), params.Key("group")),
		params.OnlyOneOf(params.Key("arn"), params.AllOf(params.Key("access"), params.Key("service"))),
	))
	builder.AddReducer(transformAccessServiceToARN, "access", "service")
	return builder.Done()
}

func transformAccessServiceToARN(values map[string]any) (map[string]any, error) {
	service, hasService := values["service"].(string)
	access, hasAccess := values["access"].(string)

	if hasService && hasAccess {
		pol, err := lookupAWSPolicy(service, access)
		if err != nil {
			return values, err
		}
		return map[string]any{"arn": pol.Arn}, nil
	} else {
		return nil, nil
	}
}

func (cmd *AttachPolicy) ManualRun(renv env.Running) (any, error) {
	start := time.Now()
	switch {
	case cmd.User != nil:
		input := &iam.AttachUserPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.UserName = cmd.User
		output, err := cmd.api.AttachUserPolicy(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("ec2.AttachUserPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Group != nil:
		input := &iam.AttachGroupPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.GroupName = cmd.Group
		output, err := cmd.api.AttachGroupPolicy(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("ec2.AttachGroupPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Role != nil:
		input := &iam.AttachRolePolicyInput{}
		input.PolicyArn = cmd.Arn
		input.RoleName = cmd.Role
		output, err := cmd.api.AttachRolePolicy(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("ec2.AttachRolePolicy call took %s", time.Since(start))
		return output, err
	default:
		return nil, errors.New("missing one of 'user, group, role' param")
	}
}

type DetachPolicy struct {
	_      string `action:"detach" entity:"policy" awsAPI:"iam"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *iam.Client
	Arn    *string `awsName:"PolicyArn" awsType:"awsstr" templateName:"arn"`
	User   *string `awsName:"UserName" awsType:"awsstr" templateName:"user"`
	Group  *string `awsName:"GroupName" awsType:"awsstr" templateName:"group"`
	Role   *string `awsName:"RoleName" awsType:"awsstr" templateName:"role"`
}

func (cmd *DetachPolicy) ParamsSpec() params.Spec {
	builder := params.SpecBuilder(params.AllOf(
		params.OnlyOneOf(params.Key("user"), params.Key("role"), params.Key("group")),
		params.OnlyOneOf(params.Key("arn"), params.AllOf(params.Key("access"), params.Key("service"))),
	))
	builder.AddReducer(transformAccessServiceToARN, "access", "service")
	return builder.Done()
}

func (cmd *DetachPolicy) ManualRun(renv env.Running) (any, error) {
	start := time.Now()
	switch {
	case cmd.User != nil:
		input := &iam.DetachUserPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.UserName = cmd.User
		output, err := cmd.api.DetachUserPolicy(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("ec2.DetachUserPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Group != nil:
		input := &iam.DetachGroupPolicyInput{}
		input.PolicyArn = cmd.Arn
		input.GroupName = cmd.Group
		output, err := cmd.api.DetachGroupPolicy(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("ec2.DetachGroupPolicy call took %s", time.Since(start))
		return output, err
	case cmd.Role != nil:
		input := &iam.DetachRolePolicyInput{}
		input.PolicyArn = cmd.Arn
		input.RoleName = cmd.Role
		output, err := cmd.api.DetachRolePolicy(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("ec2.DetachRolePolicy call took %s", time.Since(start))
		return output, err
	default:
		return nil, errors.New("missing one of 'user, group, role' param")
	}
}

// Marshaled into an IAM policy document, where these key names are required by AWS
// rather than chosen here. Tagged explicitly so a field rename cannot change them.
type policyBody struct {
	Version   string             `json:"Version"`
	Statement []*policyStatement `json:"Statement"`
}

type policyStatement struct {
	Effect     string           `json:",omitempty"`
	Actions    []string         `json:"Action,omitempty"`
	Resources  []string         `json:"Resource,omitempty"`
	Principal  *principal       `json:",omitempty"`
	Conditions policyConditions `json:"Condition,omitempty"`
}

type principal struct {
	AWS     any `json:",omitempty"`
	Service any `json:",omitempty"`
}

type policyCondition struct {
	Type  string
	Key   string
	Value string
}

func buildStatementFromParams(effect *string, resource, action, condition []*string) (*policyStatement, error) {
	stat := &policyStatement{Effect: capitalize(StringValue(effect))}
	if resource != nil {
		res := castStringSlice(resource)
		if len(res) == 1 && res[0] == "all" {
			res[0] = "*"
		}
		stat.Resources = res
	}

	if action != nil {
		stat.Actions = castStringSlice(action)
	}
	if condition != nil {
		condStr := castStringSlice(condition)
		for _, str := range condStr {
			cond, err := parseCondition(str)
			if err != nil {
				return stat, err
			}
			stat.Conditions = append(stat.Conditions, cond)
		}
	}
	return stat, nil
}

type policyConditions []*policyCondition

func (c *policyConditions) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("\"\""), nil
	}
	var buff bytes.Buffer
	buff.WriteRune('{')
	for i, cond := range *c {
		fmt.Fprintf(&buff, "\"%s\":{\"%s\":\"%s\"}", cond.Type, cond.Key, cond.Value)
		if i < len(*c)-1 {
			buff.WriteRune(',')
		}
	}
	buff.WriteRune('}')
	return buff.Bytes(), nil
}

var conditionRegex = regexp.MustCompile(`^([a-zA-Z0-9:_\-\[\]\*]+)(==|!=|=~|!~|<=|>=|<|>)(.*)$`)

func parseCondition(condition string) (*policyCondition, error) {
	matches := conditionRegex.FindStringSubmatch(condition)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid condition '%s'", condition)
	}
	key, operator, value := matches[1], matches[2], matches[3]
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") && len(value) >= 2 {
		value = value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) >= 2 {
		value = value[1 : len(value)-1]
	}

	if strings.ToLower(value) == "null" {
		switch operator {
		case "==":
			return &policyCondition{Type: "Null", Key: key, Value: "true"}, nil
		case "!=":
			return &policyCondition{Type: "Null", Key: key, Value: "false"}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for null value '%s', expected either '==' or '!='", operator, value)
		}
	} else if strings.HasPrefix(value, "arn:") {
		switch operator {
		case "==":
			return &policyCondition{Type: "ArnEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "ArnNotEquals", Key: key, Value: value}, nil
		case "=~":
			return &policyCondition{Type: "ArnLike", Key: key, Value: value}, nil
		case "!~":
			return &policyCondition{Type: "ArnNotLike", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for arn value '%s', expected either '==', '!=', '=~' or '!~'", operator, value)
		}
	} else if _, _, cidrErr := net.ParseCIDR(value); cidrErr == nil || net.ParseIP(value) != nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "IpAddress", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "NotIpAddress", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for IP value '%s', expected either '==' or '!='", operator, value)
		}
	} else if _, err := time.Parse("2006-01-02T15:04:05Z", value); err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "DateEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "DateNotEquals", Key: key, Value: value}, nil
		case "<":
			return &policyCondition{Type: "DateLessThan", Key: key, Value: value}, nil
		case "<=":
			return &policyCondition{Type: "DateLessThanEquals", Key: key, Value: value}, nil
		case ">":
			return &policyCondition{Type: "DateGreaterThan", Key: key, Value: value}, nil
		case ">=":
			return &policyCondition{Type: "DateGreaterThanEquals", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for date value '%s', expected either '==', '!=', '>', '>=', '<' or '<='", operator, value)
		}
	} else if _, err := strconv.Atoi(value); err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "NumericEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "NumericNotEquals", Key: key, Value: value}, nil
		case "<":
			return &policyCondition{Type: "NumericLessThan", Key: key, Value: value}, nil
		case "<=":
			return &policyCondition{Type: "NumericLessThanEquals", Key: key, Value: value}, nil
		case ">":
			return &policyCondition{Type: "NumericGreaterThan", Key: key, Value: value}, nil
		case ">=":
			return &policyCondition{Type: "NumericGreaterThanEquals", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for int value '%s', expected either '==', '!=', '>', '>=', '<' or '<='", operator, value)
		}
	} else if b, err := strconv.ParseBool(value); err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "Bool", Key: key, Value: fmt.Sprint(b)}, nil
		case "!=":
			return &policyCondition{Type: "Bool", Key: key, Value: fmt.Sprint(!b)}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for bool value '%s', expected either '==' or '!='", operator, value)
		}
	} else if _, err := base64.StdEncoding.DecodeString(value); value != "" && err == nil {
		switch operator {
		case "==":
			return &policyCondition{Type: "BinaryEquals", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for binary value '%s', expected '=='", operator, value)
		}
	} else {
		switch operator {
		case "==":
			return &policyCondition{Type: "StringEquals", Key: key, Value: value}, nil
		case "!=":
			return &policyCondition{Type: "StringNotEquals", Key: key, Value: value}, nil
		case "=~":
			return &policyCondition{Type: "StringLike", Key: key, Value: value}, nil
		case "!~":
			return &policyCondition{Type: "StringNotLike", Key: key, Value: value}, nil
		default:
			return nil, fmt.Errorf("invalid operator '%s' for string value '%s', expected either '==', '!=', '=~' or '!~'", operator, value)
		}
	}
}

func lookupAWSPolicy(service, access string) (*policy, error) {
	if access != "readonly" && access != "full" {
		return nil, errors.New("looking up AWS policies: access value can only be 'readonly' or 'full'")
	}

	var suggestions []string
	for _, p := range awsPolicies {
		name := strings.ToLower(p.Name)
		match := fmt.Sprintf("%s%s", strings.ToLower(service), strings.ToLower(access))
		if strings.Contains(name, match) {
			return p, nil
		}
		if strings.Contains(name, strings.ToLower(service)) {
			suggestions = append(suggestions, fmt.Sprintf("\t\tarn=%s", p.Arn))
		}
	}

	errBuff := bytes.NewBufferString(fmt.Sprintf("No AWS policy matching service '%s' and access '%s'", service, access))
	if len(suggestions) > 0 {
		errBuff.WriteString(". Try using the full ARN of those potential matches:\n")
		errBuff.WriteString(strings.Join(suggestions, "\n"))
	}

	return nil, errors.New(errBuff.String())
}

type policy struct {
	Name string `json:"PolicyName"`
	ID   string `json:"PolicyId"`
	Arn  string `json:"Arn"`
}

var awsPolicies = []*policy{
	{
		Name: "AWSDirectConnectReadOnlyAccess",
		ID:   "ANPAI23HZ27SI6FQMGNQ2",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectConnectReadOnlyAccess",
	},
	{
		Name: "AmazonGlacierReadOnlyAccess",
		ID:   "ANPAI2D5NJKMU274MET4E",
		Arn:  "arn:aws:iam::aws:policy/AmazonGlacierReadOnlyAccess",
	},
	{
		Name: "AWSMarketplaceFullAccess",
		ID:   "ANPAI2DV5ULJSO2FYVPYG",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceFullAccess",
	},
	{
		Name: "AutoScalingConsoleReadOnlyAccess",
		ID:   "ANPAI3A7GDXOYQV3VUQMK",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingConsoleReadOnlyAccess",
	},
	{
		Name: "AmazonDMSRedshiftS3Role",
		ID:   "ANPAI3CCUQ4U5WNC5F6B6",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonDMSRedshiftS3Role",
	},
	{
		Name: "AWSQuickSightListIAM",
		ID:   "ANPAI3CH5UUWZN4EKGILO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuickSightListIAM",
	},
	{
		Name: "AWSHealthFullAccess",
		ID:   "ANPAI3CUMPCPEUPCSXC4Y",
		Arn:  "arn:aws:iam::aws:policy/AWSHealthFullAccess",
	},
	{
		Name: "AmazonRDSFullAccess",
		ID:   "ANPAI3R4QMOG6Q5A4VWVG",
		Arn:  "arn:aws:iam::aws:policy/AmazonRDSFullAccess",
	},
	{
		Name: "SupportUser",
		ID:   "ANPAI3V4GSSN5SJY3P2RO",
		Arn:  "arn:aws:iam::aws:policy/job-function/SupportUser",
	},
	{
		Name: "AmazonEC2FullAccess",
		ID:   "ANPAI3VAJF5ZCRZ7MCQE6",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
	},
	{
		Name: "AWSElasticBeanstalkReadOnlyAccess",
		ID:   "ANPAI47KNGXDAXFD4SDHG",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkReadOnlyAccess",
	},
	{
		Name: "AWSCertificateManagerReadOnly",
		ID:   "ANPAI4GSWX6S4MESJ3EWC",
		Arn:  "arn:aws:iam::aws:policy/AWSCertificateManagerReadOnly",
	},
	{
		Name: "AWSQuicksightAthenaAccess",
		ID:   "ANPAI4JB77JXFQXDWNRPM",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuicksightAthenaAccess",
	},
	{
		Name: "AWSCodeCommitPowerUser",
		ID:   "ANPAI4UIINUVGB5SEC57G",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeCommitPowerUser",
	},
	{
		Name: "AWSCodeCommitFullAccess",
		ID:   "ANPAI4VCZ3XPIZLQ5NZV2",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeCommitFullAccess",
	},
	{
		Name: "IAMSelfManageServiceSpecificCredentials",
		ID:   "ANPAI4VT74EMXK2PMQJM2",
		Arn:  "arn:aws:iam::aws:policy/IAMSelfManageServiceSpecificCredentials",
	},
	{
		Name: "AmazonSQSFullAccess",
		ID:   "ANPAI65L554VRJ33ECQS6",
		Arn:  "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
	},
	{
		Name: "AWSLambdaFullAccess",
		ID:   "ANPAI6E2CYYMI4XI7AA5K",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaFullAccess",
	},
	{
		Name: "AWSIoTLogging",
		ID:   "ANPAI6R6Z2FHHGS454W7W",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSIoTLogging",
	},
	{
		Name: "AmazonEC2RoleforSSM",
		ID:   "ANPAI6TL3SMY22S4KMMX6",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM",
	},
	{
		Name: "AWSCloudHSMRole",
		ID:   "ANPAI7QIUU4GC66SF26WE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSCloudHSMRole",
	},
	{
		Name: "IAMFullAccess",
		ID:   "ANPAI7XKCFMBPM3QQRRVQ",
		Arn:  "arn:aws:iam::aws:policy/IAMFullAccess",
	},
	{
		Name: "AmazonInspectorFullAccess",
		ID:   "ANPAI7Y6NTA27NWNA5U5E",
		Arn:  "arn:aws:iam::aws:policy/AmazonInspectorFullAccess",
	},
	{
		Name: "AmazonElastiCacheFullAccess",
		ID:   "ANPAIA2V44CPHAUAAECKG",
		Arn:  "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
	},
	{
		Name: "AWSAgentlessDiscoveryService",
		ID:   "ANPAIA3DIL7BYQ35ISM4K",
		Arn:  "arn:aws:iam::aws:policy/AWSAgentlessDiscoveryService",
	},
	{
		Name: "AWSXrayWriteOnlyAccess",
		ID:   "ANPAIAACM4LMYSRGBCTM6",
		Arn:  "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
	},
	{
		Name: "AutoScalingReadOnlyAccess",
		ID:   "ANPAIAFWUVLC2LPLSFTFG",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingReadOnlyAccess",
	},
	{
		Name: "AutoScalingFullAccess",
		ID:   "ANPAIAWRCSJDDXDXGPCFU",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingFullAccess",
	},
	{
		Name: "AmazonEC2RoleforAWSCodeDeploy",
		ID:   "ANPAIAZKXZ27TAJ4PVWGK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy",
	},
	{
		Name: "AWSMobileHub_ReadOnly",
		ID:   "ANPAIBXVYVL3PWQFBZFGW",
		Arn:  "arn:aws:iam::aws:policy/AWSMobileHub_ReadOnly",
	},
	{
		Name: "CloudWatchEventsBuiltInTargetExecutionAccess",
		ID:   "ANPAIC5AQ5DATYSNF4AUM",
		Arn:  "arn:aws:iam::aws:policy/service-role/CloudWatchEventsBuiltInTargetExecutionAccess",
	},
	{
		Name: "AmazonCloudDirectoryReadOnlyAccess",
		ID:   "ANPAICMSZQGR3O62KMD6M",
		Arn:  "arn:aws:iam::aws:policy/AmazonCloudDirectoryReadOnlyAccess",
	},
	{
		Name: "AWSOpsWorksFullAccess",
		ID:   "ANPAICN26VXMXASXKOQCG",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksFullAccess",
	},
	{
		Name: "AWSOpsWorksCMInstanceProfileRole",
		ID:   "ANPAICSU3OSHCURP2WIZW",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksCMInstanceProfileRole",
	},
	{
		Name: "AWSCodePipelineApproverAccess",
		ID:   "ANPAICXNWK42SQ6LMDXM2",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineApproverAccess",
	},
	{
		Name: "AWSApplicationDiscoveryAgentAccess",
		ID:   "ANPAICZIOVAGC6JPF3WHC",
		Arn:  "arn:aws:iam::aws:policy/AWSApplicationDiscoveryAgentAccess",
	},
	{
		Name: "ViewOnlyAccess",
		ID:   "ANPAID22R6XPJATWOFDK6",
		Arn:  "arn:aws:iam::aws:policy/job-function/ViewOnlyAccess",
	},
	{
		Name: "AmazonElasticMapReduceRole",
		ID:   "ANPAIDI2BQT2LKXZG36TW",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceRole",
	},
	{
		Name: "AmazonRoute53DomainsReadOnlyAccess",
		ID:   "ANPAIDRINP6PPTRXYVQCI",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53DomainsReadOnlyAccess",
	},
	{
		Name: "AWSOpsWorksRole",
		ID:   "ANPAIDUTMOKHJFAPJV45W",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSOpsWorksRole",
	},
	{
		Name: "ApplicationAutoScalingForAmazonAppStreamAccess",
		ID:   "ANPAIEL3HJCCWFVHA6KPG",
		Arn:  "arn:aws:iam::aws:policy/service-role/ApplicationAutoScalingForAmazonAppStreamAccess",
	},
	{
		Name: "AmazonEC2ContainerRegistryFullAccess",
		ID:   "ANPAIESRL7KD7IIVF6V4W",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
	},
	{
		Name: "SimpleWorkflowFullAccess",
		ID:   "ANPAIFE3AV6VE7EANYBVM",
		Arn:  "arn:aws:iam::aws:policy/SimpleWorkflowFullAccess",
	},
	{
		Name: "AmazonS3FullAccess",
		ID:   "ANPAIFIR6V6BVTRAHWINE",
		Arn:  "arn:aws:iam::aws:policy/AmazonS3FullAccess",
	},
	{
		Name: "AWSStorageGatewayReadOnlyAccess",
		ID:   "ANPAIFKCTUVOPD5NICXJK",
		Arn:  "arn:aws:iam::aws:policy/AWSStorageGatewayReadOnlyAccess",
	},
	{
		Name: "Billing",
		ID:   "ANPAIFTHXT6FFMIRT7ZEA",
		Arn:  "arn:aws:iam::aws:policy/job-function/Billing",
	},
	{
		Name: "QuickSightAccessForS3StorageManagementAnalyticsReadOnly",
		ID:   "ANPAIFWG3L3WDMR4I7ZJW",
		Arn:  "arn:aws:iam::aws:policy/service-role/QuickSightAccessForS3StorageManagementAnalyticsReadOnly",
	},
	{
		Name: "AmazonEC2ContainerRegistryReadOnly",
		ID:   "ANPAIFYZPA37OOHVIH7KQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
	},
	{
		Name: "AmazonElasticMapReduceforEC2Role",
		ID:   "ANPAIGALS5RCDLZLB3PGS",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforEC2Role",
	},
	{
		Name: "DatabaseAdministrator",
		ID:   "ANPAIGBMAW4VUQKOQNVT6",
		Arn:  "arn:aws:iam::aws:policy/job-function/DatabaseAdministrator",
	},
	{
		Name: "AmazonRedshiftReadOnlyAccess",
		ID:   "ANPAIGD46KSON64QBSEZM",
		Arn:  "arn:aws:iam::aws:policy/AmazonRedshiftReadOnlyAccess",
	},
	{
		Name: "AmazonEC2ReadOnlyAccess",
		ID:   "ANPAIGDT4SV4GSETWTBZK",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess",
	},
	{
		Name: "AWSXrayReadOnlyAccess",
		ID:   "ANPAIH4OFXWPS6ZX6OPGQ",
		Arn:  "arn:aws:iam::aws:policy/AWSXrayReadOnlyAccess",
	},
	{
		Name: "AWSElasticBeanstalkEnhancedHealth",
		ID:   "ANPAIH5EFJNMOGUUTKLFE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkEnhancedHealth",
	},
	{
		Name: "AmazonElasticMapReduceReadOnlyAccess",
		ID:   "ANPAIHP6NH2S6GYFCOINC",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticMapReduceReadOnlyAccess",
	},
	{
		Name: "AWSDirectoryServiceReadOnlyAccess",
		ID:   "ANPAIHWYO6WSDNCG64M2W",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectoryServiceReadOnlyAccess",
	},
	{
		Name: "AmazonVPCReadOnlyAccess",
		ID:   "ANPAIICZJNOJN36GTG6CM",
		Arn:  "arn:aws:iam::aws:policy/AmazonVPCReadOnlyAccess",
	},
	{
		Name: "CloudWatchEventsReadOnlyAccess",
		ID:   "ANPAIILJPXXA6F7GYLYBS",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchEventsReadOnlyAccess",
	},
	{
		Name: "AmazonAPIGatewayInvokeFullAccess",
		ID:   "ANPAIIWAX2NOOQJ4AIEQ6",
		Arn:  "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
	},
	{
		Name: "AmazonKinesisAnalyticsReadOnly",
		ID:   "ANPAIJIEXZAFUK43U7ARK",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisAnalyticsReadOnly",
	},
	{
		Name: "AmazonMobileAnalyticsFullAccess",
		ID:   "ANPAIJIKLU2IJ7WJ6DZFG",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsFullAccess",
	},
	{
		Name: "AWSMobileHub_FullAccess",
		ID:   "ANPAIJLU43R6AGRBK76DM",
		Arn:  "arn:aws:iam::aws:policy/AWSMobileHub_FullAccess",
	},
	{
		Name: "AmazonAPIGatewayPushToCloudWatchLogs",
		ID:   "ANPAIK4GFO7HLKYN64ASK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs",
	},
	{
		Name: "AWSDataPipelineRole",
		ID:   "ANPAIKCP6XS3ESGF4GLO2",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSDataPipelineRole",
	},
	{
		Name: "CloudWatchFullAccess",
		ID:   "ANPAIKEABORKUXN6DEAZU",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchFullAccess",
	},
	{
		Name: "ServiceCatalogAdminFullAccess",
		ID:   "ANPAIKTX42IAS75B7B7BY",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogAdminFullAccess",
	},
	{
		Name: "AmazonRDSDirectoryServiceAccess",
		ID:   "ANPAIL4KBY57XWMYUHKUU",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonRDSDirectoryServiceAccess",
	},
	{
		Name: "AWSCodePipelineReadOnlyAccess",
		ID:   "ANPAILFKZXIBOTNC5TO2Q",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineReadOnlyAccess",
	},
	{
		Name: "ReadOnlyAccess",
		ID:   "ANPAILL3HVNFSB6DCOWYQ",
		Arn:  "arn:aws:iam::aws:policy/ReadOnlyAccess",
	},
	{
		Name: "AmazonMachineLearningBatchPredictionsAccess",
		ID:   "ANPAILOI4HTQSFTF3GQSC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningBatchPredictionsAccess",
	},
	{
		Name: "AmazonRekognitionReadOnlyAccess",
		ID:   "ANPAILWSUHXUY4ES43SA4",
		Arn:  "arn:aws:iam::aws:policy/AmazonRekognitionReadOnlyAccess",
	},
	{
		Name: "AWSCodeDeployReadOnlyAccess",
		ID:   "ANPAILZHHKCKB4NE7XOIQ",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeDeployReadOnlyAccess",
	},
	{
		Name: "CloudSearchFullAccess",
		ID:   "ANPAIM6OOWKQ7L7VBOZOC",
		Arn:  "arn:aws:iam::aws:policy/CloudSearchFullAccess",
	},
	{
		Name: "AWSCloudHSMFullAccess",
		ID:   "ANPAIMBQYQZM7F63DA2UU",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudHSMFullAccess",
	},
	{
		Name: "AmazonEC2SpotFleetAutoscaleRole",
		ID:   "ANPAIMFFRMIOBGDP2TAVE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetAutoscaleRole",
	},
	{
		Name: "AWSCodeBuildDeveloperAccess",
		ID:   "ANPAIMKTMR34XSBQW45HS",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess",
	},
	{
		Name: "AmazonEC2SpotFleetRole",
		ID:   "ANPAIMRTKHWK7ESSNETSW",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetRole",
	},
	{
		Name: "AWSDataPipeline_PowerUser",
		ID:   "ANPAIMXGLVY6DVR24VTYS",
		Arn:  "arn:aws:iam::aws:policy/AWSDataPipeline_PowerUser",
	},
	{
		Name: "AmazonElasticTranscoderJobsSubmitter",
		ID:   "ANPAIN5WGARIKZ3E2UQOU",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticTranscoderJobsSubmitter",
	},
	{
		Name: "AWSCodeStarServiceRole",
		ID:   "ANPAIN6D4M2KD3NBOC4M4",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSCodeStarServiceRole",
	},
	{
		Name: "AWSDirectoryServiceFullAccess",
		ID:   "ANPAINAW5ANUWTH3R4ANI",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectoryServiceFullAccess",
	},
	{
		Name: "AmazonDynamoDBFullAccess",
		ID:   "ANPAINUGF2JSOSUY76KYA",
		Arn:  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
	},
	{
		Name: "AmazonSESReadOnlyAccess",
		ID:   "ANPAINV2XPFRMWJJNSCGI",
		Arn:  "arn:aws:iam::aws:policy/AmazonSESReadOnlyAccess",
	},
	{
		Name: "AWSWAFReadOnlyAccess",
		ID:   "ANPAINZVDMX2SBF7EU2OC",
		Arn:  "arn:aws:iam::aws:policy/AWSWAFReadOnlyAccess",
	},
	{
		Name: "AutoScalingNotificationAccessRole",
		ID:   "ANPAIO2VMUPGDC5PZVXVA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AutoScalingNotificationAccessRole",
	},
	{
		Name: "AmazonMechanicalTurkReadOnly",
		ID:   "ANPAIO5IY3G3WXSX5PPRM",
		Arn:  "arn:aws:iam::aws:policy/AmazonMechanicalTurkReadOnly",
	},
	{
		Name: "AmazonKinesisReadOnlyAccess",
		ID:   "ANPAIOCMTDT5RLKZ2CAJO",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisReadOnlyAccess",
	},
	{
		Name: "AWSCodeDeployFullAccess",
		ID:   "ANPAIONKN3TJZUKXCHXWC",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeDeployFullAccess",
	},
	{
		Name: "CloudWatchActionsEC2Access",
		ID:   "ANPAIOWD4E3FVSORSZTGU",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchActionsEC2Access",
	},
	{
		Name: "AWSLambdaDynamoDBExecutionRole",
		ID:   "ANPAIP7WNAGMIPYNW4WQG",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaDynamoDBExecutionRole",
	},
	{
		Name: "AmazonRoute53DomainsFullAccess",
		ID:   "ANPAIPAFBMIYUILMOKL6G",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53DomainsFullAccess",
	},
	{
		Name: "AmazonElastiCacheReadOnlyAccess",
		ID:   "ANPAIPDACSNQHSENWAKM2",
		Arn:  "arn:aws:iam::aws:policy/AmazonElastiCacheReadOnlyAccess",
	},
	{
		Name: "AmazonAthenaFullAccess",
		ID:   "ANPAIPJMLMD4C7RYZ6XCK",
		Arn:  "arn:aws:iam::aws:policy/AmazonAthenaFullAccess",
	},
	{
		Name: "AmazonElasticFileSystemReadOnlyAccess",
		ID:   "ANPAIPN5S4NE5JJOKVC4Y",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticFileSystemReadOnlyAccess",
	},
	{
		Name: "CloudFrontFullAccess",
		ID:   "ANPAIPRV52SH6HDCCFY6U",
		Arn:  "arn:aws:iam::aws:policy/CloudFrontFullAccess",
	},
	{
		Name: "AmazonMachineLearningRoleforRedshiftDataSource",
		ID:   "ANPAIQ5UDYYMNN42BM4AK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonMachineLearningRoleforRedshiftDataSource",
	},
	{
		Name: "AmazonMobileAnalyticsNon-financialReportAccess",
		ID:   "ANPAIQLKQ4RXPUBBVVRDE",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsNon-financialReportAccess",
	},
	{
		Name: "AWSCloudTrailFullAccess",
		ID:   "ANPAIQNUJTQYDRJPC3BNK",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudTrailFullAccess",
	},
	{
		Name: "AmazonCognitoDeveloperAuthenticatedIdentities",
		ID:   "ANPAIQOKZ5BGKLCMTXH4W",
		Arn:  "arn:aws:iam::aws:policy/AmazonCognitoDeveloperAuthenticatedIdentities",
	},
	{
		Name: "AWSConfigRole",
		ID:   "ANPAIQRXRDRGJUA33ELIO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSConfigRole",
	},
	{
		Name: "AmazonAppStreamServiceAccess",
		ID:   "ANPAISBRZ7LMMCBYEF3SE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonAppStreamServiceAccess",
	},
	{
		Name: "AmazonRedshiftFullAccess",
		ID:   "ANPAISEKCHH4YDB46B5ZO",
		Arn:  "arn:aws:iam::aws:policy/AmazonRedshiftFullAccess",
	},
	{
		Name: "AmazonZocaloReadOnlyAccess",
		ID:   "ANPAISRCSSJNS3QPKZJPM",
		Arn:  "arn:aws:iam::aws:policy/AmazonZocaloReadOnlyAccess",
	},
	{
		Name: "AWSCloudHSMReadOnlyAccess",
		ID:   "ANPAISVCBSY7YDBOT67KE",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudHSMReadOnlyAccess",
	},
	{
		Name: "SystemAdministrator",
		ID:   "ANPAITJPEZXCYCBXANDSW",
		Arn:  "arn:aws:iam::aws:policy/job-function/SystemAdministrator",
	},
	{
		Name: "AmazonEC2ContainerServiceEventsRole",
		ID:   "ANPAITKFNIUAG27VSYNZ4",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceEventsRole",
	},
	{
		Name: "AmazonRoute53ReadOnlyAccess",
		ID:   "ANPAITOYK2ZAOQFXV2JNC",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53ReadOnlyAccess",
	},
	{
		Name: "AmazonEC2ReportsAccess",
		ID:   "ANPAIU6NBZVF2PCRW36ZW",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ReportsAccess",
	},
	{
		Name: "AmazonEC2ContainerServiceAutoscaleRole",
		ID:   "ANPAIUAP3EGGGXXCPDQKK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceAutoscaleRole",
	},
	{
		Name: "AWSBatchServiceRole",
		ID:   "ANPAIUETIXPCKASQJURFE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole",
	},
	{
		Name: "AWSElasticBeanstalkWebTier",
		ID:   "ANPAIUF4325SJYOREKW3A",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkWebTier",
	},
	{
		Name: "AmazonSQSReadOnlyAccess",
		ID:   "ANPAIUGSSQY362XGCM6KW",
		Arn:  "arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess",
	},
	{
		Name: "AWSMobileHub_ServiceUseOnly",
		ID:   "ANPAIUHPQXBDZUWOP3PSK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSMobileHub_ServiceUseOnly",
	},
	{
		Name: "AmazonKinesisFullAccess",
		ID:   "ANPAIVF32HAMOXCUYRAYE",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisFullAccess",
	},
	{
		Name: "AmazonMachineLearningReadOnlyAccess",
		ID:   "ANPAIW5VYBCGEX56JCINC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningReadOnlyAccess",
	},
	{
		Name: "AmazonRekognitionFullAccess",
		ID:   "ANPAIWDAOK6AIFDVX6TT6",
		Arn:  "arn:aws:iam::aws:policy/AmazonRekognitionFullAccess",
	},
	{
		Name: "RDSCloudHsmAuthorizationRole",
		ID:   "ANPAIWKFXRLQG2ROKKXLE",
		Arn:  "arn:aws:iam::aws:policy/service-role/RDSCloudHsmAuthorizationRole",
	},
	{
		Name: "AmazonMachineLearningFullAccess",
		ID:   "ANPAIWKW6AGSGYOQ5ERHC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningFullAccess",
	},
	{
		Name: "AdministratorAccess",
		ID:   "ANPAIWMBCKSKIEE64ZLYK",
		Arn:  "arn:aws:iam::aws:policy/AdministratorAccess",
	},
	{
		Name: "AmazonMachineLearningRealTimePredictionOnlyAccess",
		ID:   "ANPAIWMCNQPRWMWT36GVQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningRealTimePredictionOnlyAccess",
	},
	{
		Name: "AWSConfigUserAccess",
		ID:   "ANPAIWTTSFJ7KKJE3MWGA",
		Arn:  "arn:aws:iam::aws:policy/AWSConfigUserAccess",
	},
	{
		Name: "AWSIoTConfigAccess",
		ID:   "ANPAIWWGD4LM4EMXNRL7I",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTConfigAccess",
	},
	{
		Name: "SecurityAudit",
		ID:   "ANPAIX2T3QCXHR2OGGCTO",
		Arn:  "arn:aws:iam::aws:policy/SecurityAudit",
	},
	{
		Name: "AWSCodeStarFullAccess",
		ID:   "ANPAIXI233TFUGLZOJBEC",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeStarFullAccess",
	},
	{
		Name: "AWSDataPipeline_FullAccess",
		ID:   "ANPAIXOFIG7RSBMRPHXJ4",
		Arn:  "arn:aws:iam::aws:policy/AWSDataPipeline_FullAccess",
	},
	{
		Name: "AmazonDynamoDBReadOnlyAccess",
		ID:   "ANPAIY2XFNA232XJ6J7X2",
		Arn:  "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
	},
	{
		Name: "AutoScalingConsoleFullAccess",
		ID:   "ANPAIYEN6FJGYYWJFFCZW",
		Arn:  "arn:aws:iam::aws:policy/AutoScalingConsoleFullAccess",
	},
	{
		Name: "AmazonSNSReadOnlyAccess",
		ID:   "ANPAIZGQCQTFOFPMHSB6W",
		Arn:  "arn:aws:iam::aws:policy/AmazonSNSReadOnlyAccess",
	},
	{
		Name: "AmazonElasticMapReduceFullAccess",
		ID:   "ANPAIZP5JFP3AMSGINBB2",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticMapReduceFullAccess",
	},
	{
		Name: "AmazonS3ReadOnlyAccess",
		ID:   "ANPAIZTJ4DXE7G6AGAE6M",
		Arn:  "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
	},
	{
		Name: "AWSElasticBeanstalkFullAccess",
		ID:   "ANPAIZYX2YLLBW2LJVUFW",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkFullAccess",
	},
	{
		Name: "AmazonWorkSpacesAdmin",
		ID:   "ANPAJ26AU6ATUQCT5KVJU",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkSpacesAdmin",
	},
	{
		Name: "AWSCodeDeployRole",
		ID:   "ANPAJ2NKMKD73QS5NBFLA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole",
	},
	{
		Name: "AmazonSESFullAccess",
		ID:   "ANPAJ2P4NXCHAT7NDPNR4",
		Arn:  "arn:aws:iam::aws:policy/AmazonSESFullAccess",
	},
	{
		Name: "CloudWatchLogsReadOnlyAccess",
		ID:   "ANPAJ2YIYDYSNNEHK3VKW",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchLogsReadOnlyAccess",
	},
	{
		Name: "AmazonKinesisFirehoseReadOnlyAccess",
		ID:   "ANPAJ36NT645INW4K24W6",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisFirehoseReadOnlyAccess",
	},
	{
		Name: "AWSOpsWorksRegisterCLI",
		ID:   "ANPAJ3AB5ZBFPCQGTVDU4",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksRegisterCLI",
	},
	{
		Name: "AmazonDynamoDBFullAccesswithDataPipeline",
		ID:   "ANPAJ3ORT7KDISSXGHJXA",
		Arn:  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccesswithDataPipeline",
	},
	{
		Name: "AmazonEC2RoleforDataPipelineRole",
		ID:   "ANPAJ3Z5I2WAJE5DN2J36",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforDataPipelineRole",
	},
	{
		Name: "CloudWatchLogsFullAccess",
		ID:   "ANPAJ3ZGNWK2R5HW5BQFO",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess",
	},
	{
		Name: "AWSElasticBeanstalkMulticontainerDocker",
		ID:   "ANPAJ45SBYG72SD6SHJEY",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkMulticontainerDocker",
	},
	{
		Name: "AmazonElasticTranscoderFullAccess",
		ID:   "ANPAJ4D5OJU75P5ZJZVNY",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticTranscoderFullAccess",
	},
	{
		Name: "IAMUserChangePassword",
		ID:   "ANPAJ4L4MM2A7QIEB56MS",
		Arn:  "arn:aws:iam::aws:policy/IAMUserChangePassword",
	},
	{
		Name: "AmazonAPIGatewayAdministrator",
		ID:   "ANPAJ4PT6VY5NLKTNUYSI",
		Arn:  "arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator",
	},
	{
		Name: "ServiceCatalogEndUserAccess",
		ID:   "ANPAJ56OMCO72RI4J5FSA",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogEndUserAccess",
	},
	{
		Name: "AmazonPollyReadOnlyAccess",
		ID:   "ANPAJ5FENL3CVPL2FPDLA",
		Arn:  "arn:aws:iam::aws:policy/AmazonPollyReadOnlyAccess",
	},
	{
		Name: "AmazonMobileAnalyticsWriteOnlyAccess",
		ID:   "ANPAJ5TAWBBQC2FAL3G6G",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsWriteOnlyAccess",
	},
	{
		Name: "AmazonEC2SpotFleetTaggingRole",
		ID:   "ANPAJ5U6UMLCEYLX5OLC4",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2SpotFleetTaggingRole",
	},
	{
		Name: "DataScientist",
		ID:   "ANPAJ5YHI2BQW7EQFYDXS",
		Arn:  "arn:aws:iam::aws:policy/job-function/DataScientist",
	},
	{
		Name: "AWSMarketplaceMeteringFullAccess",
		ID:   "ANPAJ65YJPG7CC7LDXNA6",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceMeteringFullAccess",
	},
	{
		Name: "AWSOpsWorksCMServiceRole",
		ID:   "ANPAJ6I6MPGJE62URSHCO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSOpsWorksCMServiceRole",
	},
	{
		Name: "AWSConnector",
		ID:   "ANPAJ6YATONJHICG3DJ3U",
		Arn:  "arn:aws:iam::aws:policy/AWSConnector",
	},
	{
		Name: "AWSBatchFullAccess",
		ID:   "ANPAJ7K2KIWB3HZVK3CUO",
		Arn:  "arn:aws:iam::aws:policy/AWSBatchFullAccess",
	},
	{
		Name: "ServiceCatalogAdminReadOnlyAccess",
		ID:   "ANPAJ7XOUSS75M4LIPKO4",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogAdminReadOnlyAccess",
	},
	{
		Name: "AmazonSSMFullAccess",
		ID:   "ANPAJA7V6HI4ISQFMDYAG",
		Arn:  "arn:aws:iam::aws:policy/AmazonSSMFullAccess",
	},
	{
		Name: "AWSCodeCommitReadOnly",
		ID:   "ANPAJACNSXR7Z2VLJW3D6",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeCommitReadOnly",
	},
	{
		Name: "AmazonEC2ContainerServiceFullAccess",
		ID:   "ANPAJALOYVTPDZEMIACSM",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerServiceFullAccess",
	},
	{
		Name: "AmazonCognitoReadOnly",
		ID:   "ANPAJBFTRZD2GQGJHSVQK",
		Arn:  "arn:aws:iam::aws:policy/AmazonCognitoReadOnly",
	},
	{
		Name: "AmazonDMSCloudWatchLogsRole",
		ID:   "ANPAJBG7UXZZXUJD3TDJE",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonDMSCloudWatchLogsRole",
	},
	{
		Name: "AWSApplicationDiscoveryServiceFullAccess",
		ID:   "ANPAJBNJEA6ZXM2SBOPDU",
		Arn:  "arn:aws:iam::aws:policy/AWSApplicationDiscoveryServiceFullAccess",
	},
	{
		Name: "AmazonVPCFullAccess",
		ID:   "ANPAJBWPGNOVKZD3JI2P2",
		Arn:  "arn:aws:iam::aws:policy/AmazonVPCFullAccess",
	},
	{
		Name: "AWSImportExportFullAccess",
		ID:   "ANPAJCQCT4JGTLC6722MQ",
		Arn:  "arn:aws:iam::aws:policy/AWSImportExportFullAccess",
	},
	{
		Name: "AmazonMechanicalTurkFullAccess",
		ID:   "ANPAJDGCL5BET73H5QIQC",
		Arn:  "arn:aws:iam::aws:policy/AmazonMechanicalTurkFullAccess",
	},
	{
		Name: "AmazonEC2ContainerRegistryPowerUser",
		ID:   "ANPAJDNE5PIHROIBGGDDW",
		Arn:  "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser",
	},
	{
		Name: "AmazonMachineLearningCreateOnlyAccess",
		ID:   "ANPAJDRUNIC2RYAMAT3CK",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningCreateOnlyAccess",
	},
	{
		Name: "AWSCloudTrailReadOnlyAccess",
		ID:   "ANPAJDU7KJADWBSEQ3E7S",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudTrailReadOnlyAccess",
	},
	{
		Name: "AWSLambdaExecute",
		ID:   "ANPAJE5FX7FQZSU5XAKGO",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaExecute",
	},
	{
		Name: "AWSIoTRuleActions",
		ID:   "ANPAJEZ6FS7BUZVUHMOKY",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSIoTRuleActions",
	},
	{
		Name: "AWSQuickSightDescribeRedshift",
		ID:   "ANPAJFEM6MLSLTW4ZNBW2",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuickSightDescribeRedshift",
	},
	{
		Name: "VMImportExportRoleForAWSConnector",
		ID:   "ANPAJFLQOOJ6F5XNX4LAW",
		Arn:  "arn:aws:iam::aws:policy/service-role/VMImportExportRoleForAWSConnector",
	},
	{
		Name: "AWSCodePipelineCustomActionAccess",
		ID:   "ANPAJFW5Z32BTVF76VCYC",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineCustomActionAccess",
	},
	{
		Name: "AWSOpsWorksInstanceRegistration",
		ID:   "ANPAJG3LCPVNI4WDZCIMU",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksInstanceRegistration",
	},
	{
		Name: "AmazonCloudDirectoryFullAccess",
		ID:   "ANPAJG3XQK77ATFLCF2CK",
		Arn:  "arn:aws:iam::aws:policy/AmazonCloudDirectoryFullAccess",
	},
	{
		Name: "AWSStorageGatewayFullAccess",
		ID:   "ANPAJG5SSPAVOGK3SIDGU",
		Arn:  "arn:aws:iam::aws:policy/AWSStorageGatewayFullAccess",
	},
	{
		Name: "AmazonLexReadOnly",
		ID:   "ANPAJGBI5LSMAJNDGBNAM",
		Arn:  "arn:aws:iam::aws:policy/AmazonLexReadOnly",
	},
	{
		Name: "AmazonElasticTranscoderReadOnlyAccess",
		ID:   "ANPAJGPP7GPMJRRJMEP3Q",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticTranscoderReadOnlyAccess",
	},
	{
		Name: "AWSIoTConfigReadOnlyAccess",
		ID:   "ANPAJHENEMXGX4XMFOIOI",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTConfigReadOnlyAccess",
	},
	{
		Name: "AmazonWorkMailReadOnlyAccess",
		ID:   "ANPAJHF7J65E2QFKCWAJM",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkMailReadOnlyAccess",
	},
	{
		Name: "AmazonDMSVPCManagementRole",
		ID:   "ANPAJHKIGMBQI4AEFFSYO",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonDMSVPCManagementRole",
	},
	{
		Name: "AWSLambdaKinesisExecutionRole",
		ID:   "ANPAJHOLKJPXV4GBRMJUQ",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaKinesisExecutionRole",
	},
	{
		Name: "ResourceGroupsandTagEditorReadOnlyAccess",
		ID:   "ANPAJHXQTPI5I5JKAIU74",
		Arn:  "arn:aws:iam::aws:policy/ResourceGroupsandTagEditorReadOnlyAccess",
	},
	{
		Name: "AmazonSSMAutomationRole",
		ID:   "ANPAJIBQCTBCXD2XRNB6W",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole",
	},
	{
		Name: "ServiceCatalogEndUserFullAccess",
		ID:   "ANPAJIW7AFFOONVKW75KU",
		Arn:  "arn:aws:iam::aws:policy/ServiceCatalogEndUserFullAccess",
	},
	{
		Name: "AWSStepFunctionsConsoleFullAccess",
		ID:   "ANPAJIYC52YWRX6OSMJWK",
		Arn:  "arn:aws:iam::aws:policy/AWSStepFunctionsConsoleFullAccess",
	},
	{
		Name: "AWSCodeBuildReadOnlyAccess",
		ID:   "ANPAJIZZWN6557F5HVP2K",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeBuildReadOnlyAccess",
	},
	{
		Name: "AmazonMachineLearningManageRealTimeEndpointOnlyAccess",
		ID:   "ANPAJJL3PC3VCSVZP6OCI",
		Arn:  "arn:aws:iam::aws:policy/AmazonMachineLearningManageRealTimeEndpointOnlyAccess",
	},
	{
		Name: "CloudWatchEventsInvocationAccess",
		ID:   "ANPAJJXD6JKJLK2WDLZNO",
		Arn:  "arn:aws:iam::aws:policy/service-role/CloudWatchEventsInvocationAccess",
	},
	{
		Name: "CloudFrontReadOnlyAccess",
		ID:   "ANPAJJZMNYOTZCNQP36LG",
		Arn:  "arn:aws:iam::aws:policy/CloudFrontReadOnlyAccess",
	},
	{
		Name: "AmazonSNSRole",
		ID:   "ANPAJK5GQB7CIK7KHY2GA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonSNSRole",
	},
	{
		Name: "AmazonMobileAnalyticsFinancialReportAccess",
		ID:   "ANPAJKJHO2R27TXKCWBU4",
		Arn:  "arn:aws:iam::aws:policy/AmazonMobileAnalyticsFinancialReportAccess",
	},
	{
		Name: "AWSElasticBeanstalkService",
		ID:   "ANPAJKQ5SN74ZQ4WASXBM",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkService",
	},
	{
		Name: "IAMReadOnlyAccess",
		ID:   "ANPAJKSO7NDY4T57MWDSQ",
		Arn:  "arn:aws:iam::aws:policy/IAMReadOnlyAccess",
	},
	{
		Name: "AmazonRDSReadOnlyAccess",
		ID:   "ANPAJKTTTYV2IIHKLZ346",
		Arn:  "arn:aws:iam::aws:policy/AmazonRDSReadOnlyAccess",
	},
	{
		Name: "AmazonCognitoPowerUser",
		ID:   "ANPAJKW5H2HNCPGCYGR6Y",
		Arn:  "arn:aws:iam::aws:policy/AmazonCognitoPowerUser",
	},
	{
		Name: "AmazonElasticFileSystemFullAccess",
		ID:   "ANPAJKXTMNVQGIDNCKPBC",
		Arn:  "arn:aws:iam::aws:policy/AmazonElasticFileSystemFullAccess",
	},
	{
		Name: "ServerMigrationConnector",
		ID:   "ANPAJKZRWXIPK5HSG3QDQ",
		Arn:  "arn:aws:iam::aws:policy/ServerMigrationConnector",
	},
	{
		Name: "AmazonZocaloFullAccess",
		ID:   "ANPAJLCDXYRINDMUXEVL6",
		Arn:  "arn:aws:iam::aws:policy/AmazonZocaloFullAccess",
	},
	{
		Name: "AWSLambdaReadOnlyAccess",
		ID:   "ANPAJLDG7J3CGUHFN4YN6",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaReadOnlyAccess",
	},
	{
		Name: "AWSAccountUsageReportAccess",
		ID:   "ANPAJLIB4VSBVO47ZSBB6",
		Arn:  "arn:aws:iam::aws:policy/AWSAccountUsageReportAccess",
	},
	{
		Name: "AWSMarketplaceGetEntitlements",
		ID:   "ANPAJLPIMQE4WMHDC2K7C",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceGetEntitlements",
	},
	{
		Name: "AmazonEC2ContainerServiceforEC2Role",
		ID:   "ANPAJLYJCVHC7TQHCSQDS",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role",
	},
	{
		Name: "AmazonAppStreamFullAccess",
		ID:   "ANPAJLZZXU2YQVGL4QDNC",
		Arn:  "arn:aws:iam::aws:policy/AmazonAppStreamFullAccess",
	},
	{
		Name: "AWSIoTDataAccess",
		ID:   "ANPAJM2KI2UJDR24XPS2K",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTDataAccess",
	},
	{
		Name: "AmazonESFullAccess",
		ID:   "ANPAJM6ZTCU24QL5PZCGC",
		Arn:  "arn:aws:iam::aws:policy/AmazonESFullAccess",
	},
	{
		Name: "ServerMigrationServiceRole",
		ID:   "ANPAJMBH3M6BO63XFW2D4",
		Arn:  "arn:aws:iam::aws:policy/service-role/ServerMigrationServiceRole",
	},
	{
		Name: "AWSWAFFullAccess",
		ID:   "ANPAJMIKIAFXZEGOLRH7C",
		Arn:  "arn:aws:iam::aws:policy/AWSWAFFullAccess",
	},
	{
		Name: "AmazonKinesisFirehoseFullAccess",
		ID:   "ANPAJMZQMTZ7FRBFHHAHI",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisFirehoseFullAccess",
	},
	{
		Name: "CloudWatchReadOnlyAccess",
		ID:   "ANPAJN23PDQP7SZQAE3QE",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess",
	},
	{
		Name: "AWSLambdaBasicExecutionRole",
		ID:   "ANPAJNCQGXC42545SKXIK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
	},
	{
		Name: "ResourceGroupsandTagEditorFullAccess",
		ID:   "ANPAJNOS54ZFXN4T2Y34A",
		Arn:  "arn:aws:iam::aws:policy/ResourceGroupsandTagEditorFullAccess",
	},
	{
		Name: "AWSKeyManagementServicePowerUser",
		ID:   "ANPAJNPP7PPPPMJRV2SA4",
		Arn:  "arn:aws:iam::aws:policy/AWSKeyManagementServicePowerUser",
	},
	{
		Name: "AWSImportExportReadOnlyAccess",
		ID:   "ANPAJNTV4OG52ESYZHCNK",
		Arn:  "arn:aws:iam::aws:policy/AWSImportExportReadOnlyAccess",
	},
	{
		Name: "AmazonElasticTranscoderRole",
		ID:   "ANPAJNW3WMKVXFJ2KPIQ2",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticTranscoderRole",
	},
	{
		Name: "AmazonEC2ContainerServiceRole",
		ID:   "ANPAJO53W2XHNACG7V77Q",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole",
	},
	{
		Name: "AWSDeviceFarmFullAccess",
		ID:   "ANPAJO7KEDP4VYJPNT5UW",
		Arn:  "arn:aws:iam::aws:policy/AWSDeviceFarmFullAccess",
	},
	{
		Name: "AmazonSSMReadOnlyAccess",
		ID:   "ANPAJODSKQGGJTHRYZ5FC",
		Arn:  "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess",
	},
	{
		Name: "AWSStepFunctionsReadOnlyAccess",
		ID:   "ANPAJONHB2TJQDJPFW5TM",
		Arn:  "arn:aws:iam::aws:policy/AWSStepFunctionsReadOnlyAccess",
	},
	{
		Name: "AWSMarketplaceRead-only",
		ID:   "ANPAJOOM6LETKURTJ3XZ2",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceRead-only",
	},
	{
		Name: "AWSCodePipelineFullAccess",
		ID:   "ANPAJP5LH77KSAT2KHQGG",
		Arn:  "arn:aws:iam::aws:policy/AWSCodePipelineFullAccess",
	},
	{
		Name: "AWSGreengrassResourceAccessRolePolicy",
		ID:   "ANPAJPKEIMB6YMXDEVRTM",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSGreengrassResourceAccessRolePolicy",
	},
	{
		Name: "NetworkAdministrator",
		ID:   "ANPAJPNMADZFJCVPJVZA2",
		Arn:  "arn:aws:iam::aws:policy/job-function/NetworkAdministrator",
	},
	{
		Name: "AmazonWorkSpacesApplicationManagerAdminAccess",
		ID:   "ANPAJPRL4KYETIH7XGTSS",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkSpacesApplicationManagerAdminAccess",
	},
	{
		Name: "AmazonDRSVPCManagement",
		ID:   "ANPAJPXIBTTZMBEFEX6UA",
		Arn:  "arn:aws:iam::aws:policy/AmazonDRSVPCManagement",
	},
	{
		Name: "AWSXrayFullAccess",
		ID:   "ANPAJQBYG45NSJMVQDB2K",
		Arn:  "arn:aws:iam::aws:policy/AWSXrayFullAccess",
	},
	{
		Name: "AWSElasticBeanstalkWorkerTier",
		ID:   "ANPAJQDLBRSJVKVF4JMSK",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkWorkerTier",
	},
	{
		Name: "AWSDirectConnectFullAccess",
		ID:   "ANPAJQF2QKZSK74KTIHOW",
		Arn:  "arn:aws:iam::aws:policy/AWSDirectConnectFullAccess",
	},
	{
		Name: "AWSCodeBuildAdminAccess",
		ID:   "ANPAJQJGIOIE3CD2TQXDS",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeBuildAdminAccess",
	},
	{
		Name: "AmazonKinesisAnalyticsFullAccess",
		ID:   "ANPAJQOSKHTXP43R7P5AC",
		Arn:  "arn:aws:iam::aws:policy/AmazonKinesisAnalyticsFullAccess",
	},
	{
		Name: "AWSAccountActivityAccess",
		ID:   "ANPAJQRYCWMFX5J3E333K",
		Arn:  "arn:aws:iam::aws:policy/AWSAccountActivityAccess",
	},
	{
		Name: "AmazonGlacierFullAccess",
		ID:   "ANPAJQSTZJWB2AXXAKHVQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonGlacierFullAccess",
	},
	{
		Name: "AmazonWorkMailFullAccess",
		ID:   "ANPAJQVKNMT7SVATQ4AUY",
		Arn:  "arn:aws:iam::aws:policy/AmazonWorkMailFullAccess",
	},
	{
		Name: "AWSMarketplaceManageSubscriptions",
		ID:   "ANPAJRDW2WIFN7QLUAKBQ",
		Arn:  "arn:aws:iam::aws:policy/AWSMarketplaceManageSubscriptions",
	},
	{
		Name: "AWSElasticBeanstalkCustomPlatformforEC2Role",
		ID:   "ANPAJRVFXSS6LEIQGBKDY",
		Arn:  "arn:aws:iam::aws:policy/AWSElasticBeanstalkCustomPlatformforEC2Role",
	},
	{
		Name: "AWSSupportAccess",
		ID:   "ANPAJSNKQX2OW67GF4S7E",
		Arn:  "arn:aws:iam::aws:policy/AWSSupportAccess",
	},
	{
		Name: "AmazonElasticMapReduceforAutoScalingRole",
		ID:   "ANPAJSVXG6QHPE6VHDZ4Q",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonElasticMapReduceforAutoScalingRole",
	},
	{
		Name: "AWSLambdaInvocation-DynamoDB",
		ID:   "ANPAJTHQ3EKCQALQDYG5G",
		Arn:  "arn:aws:iam::aws:policy/AWSLambdaInvocation-DynamoDB",
	},
	{
		Name: "IAMUserSSHKeys",
		ID:   "ANPAJTSHUA4UXGXU7ANUA",
		Arn:  "arn:aws:iam::aws:policy/IAMUserSSHKeys",
	},
	{
		Name: "AWSIoTFullAccess",
		ID:   "ANPAJU2FPGG6PQWN72V2G",
		Arn:  "arn:aws:iam::aws:policy/AWSIoTFullAccess",
	},
	{
		Name: "AWSQuickSightDescribeRDS",
		ID:   "ANPAJU5J6OAMCJD3OO76O",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSQuickSightDescribeRDS",
	},
	{
		Name: "AWSConfigRulesExecutionRole",
		ID:   "ANPAJUB3KIKTA4PU4OYAA",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSConfigRulesExecutionRole",
	},
	{
		Name: "AmazonESReadOnlyAccess",
		ID:   "ANPAJUDMRLOQ7FPAR46FQ",
		Arn:  "arn:aws:iam::aws:policy/AmazonESReadOnlyAccess",
	},
	{
		Name: "AWSCodeDeployDeployerAccess",
		ID:   "ANPAJUWEPOMGLMVXJAPUI",
		Arn:  "arn:aws:iam::aws:policy/AWSCodeDeployDeployerAccess",
	},
	{
		Name: "AmazonPollyFullAccess",
		ID:   "ANPAJUZOYQU6XQYPR7EWS",
		Arn:  "arn:aws:iam::aws:policy/AmazonPollyFullAccess",
	},
	{
		Name: "AmazonSSMMaintenanceWindowRole",
		ID:   "ANPAJV3JNYSTZ47VOXYME",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonSSMMaintenanceWindowRole",
	},
	{
		Name: "AmazonRDSEnhancedMonitoringRole",
		ID:   "ANPAJV7BS425S4PTSSVGK",
		Arn:  "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole",
	},
	{
		Name: "AmazonLexFullAccess",
		ID:   "ANPAJVLXDHKVC23HRTKSI",
		Arn:  "arn:aws:iam::aws:policy/AmazonLexFullAccess",
	},
	{
		Name: "AWSLambdaVPCAccessExecutionRole",
		ID:   "ANPAJVTME3YLVNL72YR2K",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
	},
	{
		Name: "AmazonLexRunBotsOnly",
		ID:   "ANPAJVZGB5CM3N6YWJHBE",
		Arn:  "arn:aws:iam::aws:policy/AmazonLexRunBotsOnly",
	},
	{
		Name: "AmazonSNSFullAccess",
		ID:   "ANPAJWEKLCXXUNT2SOLSG",
		Arn:  "arn:aws:iam::aws:policy/AmazonSNSFullAccess",
	},
	{
		Name: "CloudSearchReadOnlyAccess",
		ID:   "ANPAJWPLX7N7BCC3RZLHW",
		Arn:  "arn:aws:iam::aws:policy/CloudSearchReadOnlyAccess",
	},
	{
		Name: "AWSGreengrassFullAccess",
		ID:   "ANPAJWPV6OBK4QONH4J3O",
		Arn:  "arn:aws:iam::aws:policy/AWSGreengrassFullAccess",
	},
	{
		Name: "AWSCloudFormationReadOnlyAccess",
		ID:   "ANPAJWVBEE4I2POWLODLW",
		Arn:  "arn:aws:iam::aws:policy/AWSCloudFormationReadOnlyAccess",
	},
	{
		Name: "AmazonRoute53FullAccess",
		ID:   "ANPAJWVDLG5RPST6PHQ3A",
		Arn:  "arn:aws:iam::aws:policy/AmazonRoute53FullAccess",
	},
	{
		Name: "AWSLambdaRole",
		ID:   "ANPAJX4DPCRGTC4NFDUXI",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaRole",
	},
	{
		Name: "AWSLambdaENIManagementAccess",
		ID:   "ANPAJXAW2Q3KPTURUT2QC",
		Arn:  "arn:aws:iam::aws:policy/service-role/AWSLambdaENIManagementAccess",
	},
	{
		Name: "AWSOpsWorksCloudWatchLogs",
		ID:   "ANPAJXFIK7WABAY5CPXM4",
		Arn:  "arn:aws:iam::aws:policy/AWSOpsWorksCloudWatchLogs",
	},
	{
		Name: "AmazonAppStreamReadOnlyAccess",
		ID:   "ANPAJXIFDGB4VBX23DX7K",
		Arn:  "arn:aws:iam::aws:policy/AmazonAppStreamReadOnlyAccess",
	},
	{
		Name: "AWSStepFunctionsFullAccess",
		ID:   "ANPAJXKA6VP3UFBVHDPPA",
		Arn:  "arn:aws:iam::aws:policy/AWSStepFunctionsFullAccess",
	},
	{
		Name: "AmazonInspectorReadOnlyAccess",
		ID:   "ANPAJXQNTHTEJ2JFRN2SE",
		Arn:  "arn:aws:iam::aws:policy/AmazonInspectorReadOnlyAccess",
	},
	{
		Name: "AWSCertificateManagerFullAccess",
		ID:   "ANPAJYCHABBP6VQIVBCBQ",
		Arn:  "arn:aws:iam::aws:policy/AWSCertificateManagerFullAccess",
	},
	{
		Name: "PowerUserAccess",
		ID:   "ANPAJYRXTHIB4FOVS3ZXS",
		Arn:  "arn:aws:iam::aws:policy/PowerUserAccess",
	},
	{
		Name: "CloudWatchEventsFullAccess",
		ID:   "ANPAJZLOYLNHESMYOJAFU",
		Arn:  "arn:aws:iam::aws:policy/CloudWatchEventsFullAccess",
	},
}

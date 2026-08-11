/*
Copyright 2017 WALLIX

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
	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// PutConfigRule wraps everything in a nested ConfigRule, so the awsName paths are dotted
// rather than flat. The setter creates the intermediate structs as it walks them.
//
// `owner` defaults to AWS, which is the case worth optimizing for: an AWS-managed rule
// needs only its identifier, whereas a CUSTOM_LAMBDA rule needs a function ARN as the
// identifier and usually source details as well.
type CreateConfigrule struct {
	_               string `action:"create" entity:"configrule" awsAPI:"configservice" awsCall:"PutConfigRule" awsInput:"configservice.PutConfigRuleInput" awsOutput:"configservice.PutConfigRuleOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             *configservice.Client
	Name            *string `awsName:"ConfigRule.ConfigRuleName" awsType:"awsstr" templateName:"name"`
	Source          *string `awsName:"ConfigRule.Source.SourceIdentifier" awsType:"awsstr" templateName:"source"`
	Owner           *string `awsName:"ConfigRule.Source.Owner" awsType:"awsstr" templateName:"owner"`
	Description     *string `awsName:"ConfigRule.Description" awsType:"awsstr" templateName:"description"`
	InputParameters *string `awsName:"ConfigRule.InputParameters" awsType:"awsstr" templateName:"input-parameters"`
	Frequency       *string `awsName:"ConfigRule.MaximumExecutionFrequency" awsType:"awsstr" templateName:"frequency"`
	// Restricting a rule to one resource type is the difference between evaluating a
	// handful of resources and evaluating the whole account.
	ResourceTypes []*string `awsName:"ConfigRule.Scope.ComplianceResourceTypes" awsType:"awsstringslice" templateName:"resource-types"`
}

func (cmd *CreateConfigrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("source"),
		params.Opt(
			params.Suggested("owner", "resource-types"),
			"description", "input-parameters", "frequency",
		),
	))
}

// PutConfigRule has an empty response, and the name is what the graph and delete key on.
func (cmd *CreateConfigrule) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

// PutConfigRule both creates and replaces, so update maps to the same call. Anything
// omitted is cleared rather than kept, which is why source stays required here.
type UpdateConfigrule struct {
	_               string `action:"update" entity:"configrule" awsAPI:"configservice" awsCall:"PutConfigRule" awsInput:"configservice.PutConfigRuleInput" awsOutput:"configservice.PutConfigRuleOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             *configservice.Client
	Name            *string   `awsName:"ConfigRule.ConfigRuleName" awsType:"awsstr" templateName:"name"`
	Source          *string   `awsName:"ConfigRule.Source.SourceIdentifier" awsType:"awsstr" templateName:"source"`
	Owner           *string   `awsName:"ConfigRule.Source.Owner" awsType:"awsstr" templateName:"owner"`
	Description     *string   `awsName:"ConfigRule.Description" awsType:"awsstr" templateName:"description"`
	InputParameters *string   `awsName:"ConfigRule.InputParameters" awsType:"awsstr" templateName:"input-parameters"`
	Frequency       *string   `awsName:"ConfigRule.MaximumExecutionFrequency" awsType:"awsstr" templateName:"frequency"`
	ResourceTypes   []*string `awsName:"ConfigRule.Scope.ComplianceResourceTypes" awsType:"awsstringslice" templateName:"resource-types"`
}

func (cmd *UpdateConfigrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("source"),
		params.Opt("owner", "description", "input-parameters", "frequency", "resource-types"),
	))
}

func (cmd *UpdateConfigrule) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteConfigrule struct {
	_      string `action:"delete" entity:"configrule" awsAPI:"configservice" awsCall:"DeleteConfigRule" awsInput:"configservice.DeleteConfigRuleInput" awsOutput:"configservice.DeleteConfigRuleOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *configservice.Client
	Name   *string `awsName:"ConfigRuleName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteConfigrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

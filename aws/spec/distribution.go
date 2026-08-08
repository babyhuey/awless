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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/smithy-go"

	"github.com/bootswithdefer/awless/logger"
)

var CallerReferenceFunc = func() string {
	return fmt.Sprint(time.Now().UTC().Unix())
}

type CreateDistribution struct {
	_              string `action:"create" entity:"distribution" awsAPI:"cloudfront"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *cloudfront.Client
	OriginDomain   *string   `templateName:"origin-domain"`
	Certificate    *string   `templateName:"certificate"`
	Comment        *string   `templateName:"comment"`
	DefaultFile    *string   `templateName:"default-file"`
	DomainAliases  []*string `templateName:"domain-aliases"`
	Enable         *bool     `templateName:"enable"`
	ForwardCookies *string   `templateName:"forward-cookies"`
	ForwardQueries *bool     `templateName:"forward-queries"`
	HTTPSBehavior  *string   `templateName:"https-behavior"`
	OriginPath     *string   `templateName:"origin-path"`
	PriceClass     *string   `templateName:"price-class"`
	MinTTL         *int64    `templateName:"min-ttl"`
}

func (cmd *CreateDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("origin-domain"),
		params.Opt("certificate", "comment", "default-file", "domain-aliases", "enable", "forward-cookies", "forward-queries", "https-behavior", "min-ttl", "origin-path", "price-class"),
	))
}

func (cmd *CreateDistribution) ManualRun(renv env.Running) (any, error) {
	originID := "orig_1"
	input := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfronttypes.DistributionConfig{
			CallerReference: aws.String(CallerReferenceFunc()),
			Comment:         cmd.OriginDomain,
			DefaultCacheBehavior: &cloudfronttypes.DefaultCacheBehavior{
				MinTTL: aws.Int64(0),
				ForwardedValues: &cloudfronttypes.ForwardedValues{
					Cookies:     &cloudfronttypes.CookiePreference{Forward: cloudfronttypes.ItemSelectionAll},
					QueryString: aws.Bool(true),
				},
				TrustedSigners: &cloudfronttypes.TrustedSigners{
					Enabled:  aws.Bool(false),
					Quantity: aws.Int32(0),
				},
				TargetOriginId:       aws.String(originID),
				ViewerProtocolPolicy: cloudfronttypes.ViewerProtocolPolicyAllowAll,
			},
			Enabled: aws.Bool(true),
			Origins: &cloudfronttypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cloudfronttypes.Origin{
					{Id: aws.String(originID)},
				},
			},
		},
	}

	if domain := StringValue(cmd.OriginDomain); strings.HasSuffix(domain, ".s3.amazonaws.com") || (strings.HasSuffix(domain, ".amazonaws.com") && strings.Contains(domain, ".s3-website-")) {
		input.DistributionConfig.Origins.Items[0].S3OriginConfig = &cloudfronttypes.S3OriginConfig{OriginAccessIdentity: aws.String("")}
	}

	call := &awsCall{
		fnName: "cloudfront.CreateDistribution",
		fn:     cmd.api.CreateDistribution,
		logger: cmd.logger,
		setters: []setter{
			{val: cmd.OriginDomain, fieldPath: "DistributionConfig.Origins.Items[0].DomainName", fieldType: awsstr},
		},
	}

	if cmd.Certificate != nil {
		call.setters = append(call.setters, setter{val: cmd.Certificate, fieldPath: "DistributionConfig.ViewerCertificate.ACMCertificateArn", fieldType: awsstr})
		call.setters = append(call.setters, setter{val: "sni-only", fieldPath: "DistributionConfig.ViewerCertificate.SSLSupportMethod", fieldType: awsstr})
	}

	if cmd.Comment != nil {
		call.setters = append(call.setters, setter{val: cmd.Comment, fieldPath: "DistributionConfig.Comment", fieldType: awsstr})
	}
	if cmd.DefaultFile != nil {
		call.setters = append(call.setters, setter{val: cmd.DefaultFile, fieldPath: "DistributionConfig.DefaultRootObject", fieldType: awsstr})
	}
	if cmd.DomainAliases != nil {
		call.setters = append(call.setters, setter{val: cmd.DomainAliases, fieldPath: "DistributionConfig.Aliases.Items", fieldType: awsstringslice})
		call.setters = append(call.setters, setter{val: len(cmd.DomainAliases), fieldPath: "DistributionConfig.Aliases.Quantity", fieldType: awsint64})
	}
	if cmd.Enable != nil {
		call.setters = append(call.setters, setter{val: cmd.Enable, fieldPath: "DistributionConfig.Enabled", fieldType: awsbool})
	}
	if cmd.ForwardCookies != nil {
		call.setters = append(call.setters, setter{val: cmd.ForwardCookies, fieldPath: "DistributionConfig.DefaultCacheBehavior.ForwardedValues.Cookies.Forward", fieldType: awsstr})
	}
	if cmd.ForwardQueries != nil {
		call.setters = append(call.setters, setter{val: cmd.ForwardQueries, fieldPath: "DistributionConfig.DefaultCacheBehavior.ForwardedValues.QueryString", fieldType: awsbool})
	}
	if cmd.HTTPSBehavior != nil {
		call.setters = append(call.setters, setter{val: cmd.HTTPSBehavior, fieldPath: "DistributionConfig.DefaultCacheBehavior.ViewerProtocolPolicy", fieldType: awsstr})
	}
	if cmd.MinTTL != nil {
		call.setters = append(call.setters, setter{val: cmd.MinTTL, fieldPath: "DistributionConfig.DefaultCacheBehavior.MinTTL", fieldType: awsint64})
	}
	if cmd.OriginPath != nil {
		call.setters = append(call.setters, setter{val: cmd.OriginPath, fieldPath: "DistributionConfig.Origins.Items[0].OriginPath", fieldType: awsstr})
	}
	if cmd.PriceClass != nil {
		call.setters = append(call.setters, setter{val: cmd.PriceClass, fieldPath: "DistributionConfig.PriceClass", fieldType: awsstr})
	}

	return call.execute(renv.RequestContext(), input)
}

func (cmd *CreateDistribution) ExtractResult(i any) string {
	return StringValue(i.(*cloudfront.CreateDistributionOutput).Distribution.Id)
}

type CheckDistribution struct {
	_       string `action:"check" entity:"distribution" awsAPI:"cloudfront"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *cloudfront.Client
	ID      *string `templateName:"id"`
	State   *string `templateName:"state"`
	Timeout *int64  `templateName:"timeout"`
}

func (cmd *CheckDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("id"), params.Key("state"), params.Key("timeout")),
		params.Validators{
			"state": params.IsInEnumIgnoreCase("deployed", "inprogress", notFoundState),
		})
}

func (cmd *CheckDistribution) ManualRun(renv env.Running) (any, error) {
	input := &cloudfront.GetDistributionInput{
		Id: cmd.ID,
	}

	c := &checker{
		description: fmt.Sprintf("distribution %s", StringValue(cmd.ID)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.GetDistribution(renv.RequestContext(), input)
			if err != nil {
				var awserr smithy.APIError
				if errors.As(err, &awserr) {
					if awserr.ErrorCode() == "NoSuchDistribution" {
						return notFoundState, nil
					}
					return "", awserr
				} else {
					return "", err
				}
			} else {
				return aws.ToString(output.Distribution.Status), nil
			}
		},
		expect: StringValue(cmd.State),
		logger: cmd.logger,
	}
	return nil, c.check()
}

type UpdateDistribution struct {
	_              string `action:"update" entity:"distribution" awsAPI:"cloudfront"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *cloudfront.Client
	ID             *string   `awsName:"Id" awsType:"awsstr" templateName:"id"`
	OriginDomain   *string   `templateName:"origin-domain"`
	Certificate    *string   `templateName:"certificate"`
	Comment        *string   `templateName:"comment"`
	DefaultFile    *string   `templateName:"default-file"`
	DomainAliases  []*string `templateName:"domain-aliases"`
	Enable         *bool     `templateName:"enable"`
	ForwardCookies *string   `templateName:"forward-cookies"`
	ForwardQueries *bool     `templateName:"forward-queries"`
	HTTPSBehavior  *string   `templateName:"https-behavior"`
	OriginPath     *string   `templateName:"origin-path"`
	PriceClass     *string   `templateName:"price-class"`
	MinTTL         *int64    `templateName:"min-ttl"`
}

func (cmd *UpdateDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"),
		params.Opt("certificate", "comment", "default-file", "domain-aliases", "enable", "forward-cookies", "forward-queries", "https-behavior", "min-ttl", "origin-domain", "origin-path", "price-class"),
	))
}

func (cmd *UpdateDistribution) ManualRun(renv env.Running) (any, error) {
	distribOutput, err := cmd.api.GetDistribution(renv.RequestContext(), &cloudfront.GetDistributionInput{
		Id: cmd.ID,
	})
	if err != nil {
		return nil, err
	}
	distriToUpdate := distribOutput.Distribution
	configToUpdate := distriToUpdate.DistributionConfig
	etag := distribOutput.ETag
	beforeUpdate := fmt.Sprintf("%v", distribOutput.Distribution.DistributionConfig)

	input := &cloudfront.UpdateDistributionInput{
		IfMatch:            etag,
		DistributionConfig: distriToUpdate.DistributionConfig,
	}

	if err = setFieldWithType(renv.RequestContext(), cmd.ID, input, "Id", awsstr); err != nil {
		return nil, err
	}
	if cmd.Enable != nil && BoolValue(cmd.Enable) != BoolValue(distriToUpdate.DistributionConfig.Enabled) {
		if err = setFieldWithType(renv.RequestContext(), cmd.Enable, input, "DistributionConfig.Enabled", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.OriginDomain != nil || cmd.OriginPath != nil {
		if configToUpdate.Origins == nil || len(configToUpdate.Origins.Items) == 0 {
			configToUpdate.Origins = &cloudfronttypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cloudfronttypes.Origin{
					{Id: aws.String("orig_1")},
				},
			}
		}
		if cmd.OriginDomain != nil {
			if err = setFieldWithType(renv.RequestContext(), cmd.OriginDomain, input, "DistributionConfig.Origins.Items[0].DomainName", awsstr); err != nil {
				return nil, err
			}
			if domain := aws.ToString(input.DistributionConfig.Origins.Items[0].DomainName); strings.HasSuffix(domain, ".s3.amazonaws.com") || (strings.HasSuffix(domain, ".amazonaws.com") && strings.Contains(domain, ".s3-website-")) {
				input.DistributionConfig.Origins.Items[0].S3OriginConfig = &cloudfronttypes.S3OriginConfig{OriginAccessIdentity: aws.String("")}
			}
		}

		if cmd.OriginPath != nil {
			if err = setFieldWithType(renv.RequestContext(), cmd.OriginPath, input, "DistributionConfig.Origins.Items[0].OriginPath", awsstr); err != nil {
				return nil, err
			}
		}
	}

	if cmd.Certificate != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.Certificate, input, "DistributionConfig.ViewerCertificate.ACMCertificateArn", awsstr); err != nil {
			return nil, err
		}
		if err = setFieldWithType(renv.RequestContext(), "sni-only", input, "DistributionConfig.ViewerCertificate.SSLSupportMethod", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.Comment != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.Comment, input, "DistributionConfig.Comment", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.DefaultFile != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.DefaultFile, input, "DistributionConfig.DefaultRootObject", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.DomainAliases != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.DomainAliases, input, "DistributionConfig.Aliases.Items", awsstringslice); err != nil {
			return nil, err
		}
	}
	if cmd.Enable != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.Enable, input, "DistributionConfig.Enabled", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.ForwardCookies != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.ForwardCookies, input, "DistributionConfig.DefaultCacheBehavior.ForwardedValues.Cookies.Forward", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.ForwardQueries != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.ForwardQueries, input, "DistributionConfig.DefaultCacheBehavior.ForwardedValues.QueryString", awsbool); err != nil {
			return nil, err
		}
	}
	if cmd.HTTPSBehavior != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.HTTPSBehavior, input, "DistributionConfig.DefaultCacheBehavior.ViewerProtocolPolicy", awsstr); err != nil {
			return nil, err
		}
	}
	if cmd.MinTTL != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.MinTTL, input, "DistributionConfig.DefaultCacheBehavior.MinTTL", awsint64); err != nil {
			return nil, err
		}
	}
	if cmd.PriceClass != nil {
		if err = setFieldWithType(renv.RequestContext(), cmd.PriceClass, input, "DistributionConfig.PriceClass", awsstr); err != nil {
			return nil, err
		}
	}

	if aliases := input.DistributionConfig.Aliases; aliases != nil {
		aliases.Quantity = aws.Int32(int32(len(aliases.Items)))
	}

	if beforeUpdate == fmt.Sprintf("%v", input.DistributionConfig) {
		cmd.logger.Infof("no property has been changed to distribution '%s'", StringValue(cmd.ID))
		return distribOutput, nil
	}

	start := time.Now()
	var output *cloudfront.UpdateDistributionOutput
	output, err = cmd.api.UpdateDistribution(renv.RequestContext(), input)
	cmd.logger.ExtraVerbosef("cloudfront.UpdateDistribution call took %s", time.Since(start))
	return output, err
}

func (cmd *UpdateDistribution) ExtractResult(i any) string {
	switch ii := i.(type) {
	case *cloudfront.GetDistributionOutput:
		return StringValue(ii.ETag)
	case *cloudfront.UpdateDistributionOutput:
		return StringValue(ii.ETag)
	default:
		return ""
	}
}

type DeleteDistribution struct {
	_      string `action:"delete" entity:"distribution" awsAPI:"cloudfront"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cloudfront.Client
	ID     *string `awsName:"Id" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteDistribution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

func (cmd *DeleteDistribution) ManualRun(renv env.Running) (any, error) {
	cmd.logger.Info("disabling distribution")
	updateDistribution := CommandFactory.Build("updatedistribution")().(*UpdateDistribution)
	entries := map[string]any{
		"id":     cmd.ID,
		"enable": false,
	}
	if err := params.Validate(updateDistribution.ParamsSpec().Validators(), entries); err != nil {
		return nil, err
	}

	var etag string
	if out, err := updateDistribution.Run(renv, entries); err != nil {
		return nil, err
	} else if str, ok := out.(string); ok {
		etag = str
	}

	cmd.logger.Info("check distribution disabling has been propagated")
	checkDistribution := CommandFactory.Build("checkdistribution")().(*CheckDistribution)
	entries = map[string]any{
		"id":      cmd.ID,
		"state":   "Deployed",
		"timeout": 1800,
	}
	if err := params.Validate(checkDistribution.ParamsSpec().Validators(), entries); err != nil {
		return nil, err
	}

	if _, err := checkDistribution.Run(renv, entries); err != nil {
		return nil, err
	}

	input := &cloudfront.DeleteDistributionInput{IfMatch: aws.String(fmt.Sprint(etag))}

	if err := setFieldWithType(renv.RequestContext(), cmd.ID, input, "Id", awsstr); err != nil {
		return nil, err
	}

	start := time.Now()
	output, err := cmd.api.DeleteDistribution(renv.RequestContext(), input)
	cmd.logger.ExtraVerbosef("cloudfront.DeleteDistribution call took %s", time.Since(start))
	return output, err
}

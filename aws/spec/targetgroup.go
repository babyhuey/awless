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
	"time"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/bootswithdefer/awless/logger"
)

type CreateTargetgroup struct {
	_                   string `action:"create" entity:"targetgroup" awsAPI:"elbv2" awsCall:"CreateTargetGroup" awsInput:"elbv2.CreateTargetGroupInput" awsOutput:"elbv2.CreateTargetGroupOutput"`
	logger              *logger.Logger
	graph               cloud.GraphAPI
	api                 *elbv2.Client
	Name                *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Port                *int64  `awsName:"Port" awsType:"awsint64" templateName:"port"`
	Protocol            *string `awsName:"Protocol" awsType:"awsstr" templateName:"protocol"`
	Vpc                 *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
	Healthcheckinterval *int64  `awsName:"HealthCheckIntervalSeconds" awsType:"awsint64" templateName:"healthcheckinterval"`
	Healthcheckpath     *string `awsName:"HealthCheckPath" awsType:"awsstr" templateName:"healthcheckpath"`
	Healthcheckport     *string `awsName:"HealthCheckPort" awsType:"awsstr" templateName:"healthcheckport"`
	Healthcheckprotocol *string `awsName:"HealthCheckProtocol" awsType:"awsstr" templateName:"healthcheckprotocol"`
	Healthchecktimeout  *int64  `awsName:"HealthCheckTimeoutSeconds" awsType:"awsint64" templateName:"healthchecktimeout"`
	Healthythreshold    *int64  `awsName:"HealthyThresholdCount" awsType:"awsint64" templateName:"healthythreshold"`
	Unhealthythreshold  *int64  `awsName:"UnhealthyThresholdCount" awsType:"awsint64" templateName:"unhealthythreshold"`
	Matcher             *string `awsName:"Matcher.HttpCode" awsType:"awsstr" templateName:"matcher"`
}

func (cmd *CreateTargetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("port"), params.Key("protocol"), params.Key("vpc"),
		params.Opt("healthcheckinterval", "healthcheckpath", "healthcheckport", "healthcheckprotocol", "healthchecktimeout", "healthythreshold", "matcher", "unhealthythreshold"),
	))
}

func (cmd *CreateTargetgroup) ExtractResult(i any) string {
	out := i.(*elbv2.CreateTargetGroupOutput)
	if len(out.TargetGroups) == 0 {
		return ""
	}
	return awssdk.ToString(out.TargetGroups[0].TargetGroupArn)
}

type UpdateTargetgroup struct {
	_                   string `action:"update" entity:"targetgroup" awsAPI:"elbv2"`
	logger              *logger.Logger
	graph               cloud.GraphAPI
	api                 *elbv2.Client
	ID                  *string `awsName:"TargetGroupArn" awsType:"awsstr" templateName:"id"`
	Deregistrationdelay *string `awsType:"awsstr" templateName:"deregistrationdelay"`
	Stickiness          *string `awsType:"awsstr" templateName:"stickiness"`
	Stickinessduration  *string `awsType:"awsstr" templateName:"stickinessduration"`
	Healthcheckinterval *int64  `awsName:"HealthCheckIntervalSeconds" awsType:"awsint64" templateName:"healthcheckinterval"`
	Healthcheckpath     *string `awsName:"HealthCheckPath" awsType:"awsstr" templateName:"healthcheckpath"`
	Healthcheckport     *string `awsName:"HealthCheckPort" awsType:"awsstr" templateName:"healthcheckport"`
	Healthcheckprotocol *string `awsName:"HealthCheckProtocol" awsType:"awsstr" templateName:"healthcheckprotocol"`
	Healthchecktimeout  *int64  `awsName:"HealthCheckTimeoutSeconds" awsType:"awsint64" templateName:"healthchecktimeout"`
	Healthythreshold    *int64  `awsName:"HealthyThresholdCount" awsType:"awsint64" templateName:"healthythreshold"`
	Unhealthythreshold  *int64  `awsName:"UnhealthyThresholdCount" awsType:"awsint64" templateName:"unhealthythreshold"`
	Matcher             *string `awsName:"Matcher.HttpCode" awsType:"awsstr" templateName:"matcher"`
}

func (cmd *UpdateTargetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"),
		params.Opt("deregistrationdelay", "healthcheckinterval", "healthcheckpath", "healthcheckport", "healthcheckprotocol", "healthchecktimeout", "healthythreshold", "matcher", "stickiness", "stickinessduration", "unhealthythreshold"),
	))
}

func (cmd *UpdateTargetgroup) ManualRun(renv env.Running) (any, error) {
	tgArn := StringValue(cmd.ID)

	attrsInput := &elbv2.ModifyTargetGroupAttributesInput{}
	var areTargetAttrsModified bool

	if v := cmd.Stickiness; v != nil {
		attrsInput.Attributes = append(attrsInput.Attributes, elbv2types.TargetGroupAttribute{
			Key:   String("stickiness.enabled"),
			Value: v,
		})
		areTargetAttrsModified = true
	}
	if v := cmd.Stickinessduration; v != nil {
		attrsInput.Attributes = append(attrsInput.Attributes, elbv2types.TargetGroupAttribute{
			Key:   String("stickiness.lb_cookie.duration_seconds"),
			Value: v,
		})
		areTargetAttrsModified = true
	}
	if v := cmd.Deregistrationdelay; v != nil {
		attrsInput.Attributes = append(attrsInput.Attributes, elbv2types.TargetGroupAttribute{
			Key:   String("deregistration_delay.timeout_seconds"),
			Value: v,
		})
		areTargetAttrsModified = true
	}

	var err error

	if areTargetAttrsModified {
		if err = setFieldWithType(renv.RequestContext(), tgArn, attrsInput, "TargetGroupArn", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		start := time.Now()
		if _, err = cmd.api.ModifyTargetGroupAttributes(renv.RequestContext(), attrsInput); err != nil {
			return nil, err
		}
		cmd.logger.ExtraVerbosef("elbv2.ModifyTargetGroupAttributes call took %s", time.Since(start))
	}

	input := &elbv2.ModifyTargetGroupInput{}
	var isTargetGroupModified bool

	if v := cmd.Healthcheckinterval; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "HealthCheckIntervalSeconds", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := cmd.Healthcheckpath; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "HealthCheckPath", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := cmd.Healthcheckport; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "HealthCheckPort", awsstr, renv.Context()); err != nil {
			return nil, err
		}
	}
	if v := cmd.Healthcheckprotocol; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "HealthCheckProtocol", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := cmd.Healthchecktimeout; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "HealthCheckTimeoutSeconds", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := cmd.Healthythreshold; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "HealthyThresholdCount", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := cmd.Unhealthythreshold; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "UnhealthyThresholdCount", awsint64, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}
	if v := cmd.Matcher; v != nil {
		if err = setFieldWithType(renv.RequestContext(), v, input, "Matcher.HttpCode", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		isTargetGroupModified = true
	}

	if isTargetGroupModified {
		if err = setFieldWithType(renv.RequestContext(), tgArn, input, "TargetGroupArn", awsstr, renv.Context()); err != nil {
			return nil, err
		}
		start := time.Now()
		output, err := cmd.api.ModifyTargetGroup(renv.RequestContext(), input)
		cmd.logger.ExtraVerbosef("elbv2.ModifyTargetGroup call took %s", time.Since(start))
		return output, err
	}
	return nil, nil
}

type DeleteTargetgroup struct {
	_      string `action:"delete" entity:"targetgroup" awsAPI:"elbv2" awsCall:"DeleteTargetGroup" awsInput:"elbv2.DeleteTargetGroupInput" awsOutput:"elbv2.DeleteTargetGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *elbv2.Client
	ID     *string `awsName:"TargetGroupArn" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteTargetgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

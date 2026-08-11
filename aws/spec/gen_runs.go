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

// DO NOT EDIT
// This file was automatically generated with go generate
package awsspec

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	acm "github.com/aws/aws-sdk-go-v2/service/acm"
	apigatewayv2 "github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	applicationautoscaling "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	autoscaling "github.com/aws/aws-sdk-go-v2/service/autoscaling"
	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudtrail "github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudwatch "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	configservice "github.com/aws/aws-sdk-go-v2/service/configservice"
	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	efs "github.com/aws/aws-sdk-go-v2/service/efs"
	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	elasticache "github.com/aws/aws-sdk-go-v2/service/elasticache"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	eventbridge "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	iam "github.com/aws/aws-sdk-go-v2/service/iam"
	lambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	rds "github.com/aws/aws-sdk-go-v2/service/rds"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	secretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	sfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	wafv2 "github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/smithy-go"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
)

func NewAttachAlarm(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachAlarm {
	cmd := new(AttachAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatch.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachAlarm) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachAlarm) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachAlarm) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("alarm"), nil
}

func (cmd *AttachAlarm) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachClassicLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachClassicLoadbalancer {
	cmd := new(AttachClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachClassicLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachClassicLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.RegisterInstancesWithLoadBalancerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.RegisterInstancesWithLoadBalancerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RegisterInstancesWithLoadBalancer(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elb.RegisterInstancesWithLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachClassicLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("classicloadbalancer"), nil
}

func (cmd *AttachClassicLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachContainertask(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachContainertask {
	cmd := new(AttachContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachContainertask) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachContainertask) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachContainertask) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containertask"), nil
}

func (cmd *AttachContainertask) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachElasticip(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachElasticip {
	cmd := new(AttachElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachElasticip) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachElasticip) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AssociateAddressInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AssociateAddress(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.AssociateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachElasticip) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.AssociateAddressInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.AssociateAddress(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AssociateAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach elasticip ok")
			return fakeDryRunID("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *AttachElasticip) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachEventtarget(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachEventtarget {
	cmd := new(AttachEventtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachEventtarget) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachEventtarget) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.PutTargetsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.PutTargetsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutTargets(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.PutTargets call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach eventtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach eventtarget '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach eventtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachEventtarget) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventtarget"), nil
}

func (cmd *AttachEventtarget) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachInstance {
	cmd := new(AttachInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.RegisterTargetsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.RegisterTargetsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RegisterTargets(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.RegisterTargets call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("instance"), nil
}

func (cmd *AttachInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachInstanceprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachInstanceprofile {
	cmd := new(AttachInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachInstanceprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachInstanceprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInstanceprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachInternetgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachInternetgateway {
	cmd := new(AttachInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachInternetgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachInternetgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachInternetGatewayInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AttachInternetGateway(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.AttachInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachInternetgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.AttachInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.AttachInternetGateway(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AttachInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach internetgateway ok")
			return fakeDryRunID("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *AttachInternetgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachListener(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachListener {
	cmd := new(AttachListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachListener) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachListener) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.AddListenerCertificatesInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.AddListenerCertificatesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AddListenerCertificates(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.AddListenerCertificates call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach listener '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachListener) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("listener"), nil
}

func (cmd *AttachListener) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachMfadevice(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachMfadevice {
	cmd := new(AttachMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachMfadevice) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachMfadevice) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.EnableMFADeviceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.EnableMFADeviceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.EnableMFADevice(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.EnableMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachMfadevice) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("mfadevice"), nil
}

func (cmd *AttachMfadevice) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachNetworkinterface(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachNetworkinterface {
	cmd := new(AttachNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachNetworkinterface) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachNetworkinterface) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachNetworkInterfaceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachNetworkInterfaceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AttachNetworkInterface(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.AttachNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachNetworkinterface) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.AttachNetworkInterfaceInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachNetworkInterfaceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.AttachNetworkInterface(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AttachNetworkInterface call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach networkinterface ok")
			return fakeDryRunID("networkinterface"), nil
		}
	}

	return nil, err
}

func (cmd *AttachNetworkinterface) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachPolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachPolicy {
	cmd := new(AttachPolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachPolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachPolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachPolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("policy"), nil
}

func (cmd *AttachPolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachRole(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachRole {
	cmd := new(AttachRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachRole) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachRole) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.AddRoleToInstanceProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.AddRoleToInstanceProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AddRoleToInstanceProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.AddRoleToInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach role '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachRole) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("role"), nil
}

func (cmd *AttachRole) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachRoutetable(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachRoutetable {
	cmd := new(AttachRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachRoutetable) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachRoutetable) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AssociateRouteTableInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AssociateRouteTable(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.AssociateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachRoutetable) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.AssociateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AssociateRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.AssociateRouteTable(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AssociateRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach routetable ok")
			return fakeDryRunID("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *AttachRoutetable) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachSecuritygroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachSecuritygroup {
	cmd := new(AttachSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachSecuritygroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachSecuritygroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachSecuritygroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("securitygroup"), nil
}

func (cmd *AttachSecuritygroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachUser(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachUser {
	cmd := new(AttachUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachUser) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachUser) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.AddUserToGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.AddUserToGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AddUserToGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.AddUserToGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach user '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachUser) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("user"), nil
}

func (cmd *AttachUser) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAttachVolume(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AttachVolume {
	cmd := new(AttachVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AttachVolume) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AttachVolume) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AttachVolumeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AttachVolume(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.AttachVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("attach volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("attach volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("attach volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AttachVolume) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.AttachVolumeInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AttachVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.AttachVolume(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AttachVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: attach volume ok")
			return fakeDryRunID("volume"), nil
		}
	}

	return nil, err
}

func (cmd *AttachVolume) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewAuthenticateRegistry(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *AuthenticateRegistry {
	cmd := new(AuthenticateRegistry)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecr.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *AuthenticateRegistry) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *AuthenticateRegistry) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("authenticate registry: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("authenticate registry '%s' done", extracted)
	} else {
		renv.Log().Verbose("authenticate registry done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *AuthenticateRegistry) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("registry"), nil
}

func (cmd *AuthenticateRegistry) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckCertificate(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckCertificate {
	cmd := new(CheckCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = acm.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckCertificate) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckCertificate) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check certificate '%s' done", extracted)
	} else {
		renv.Log().Verbose("check certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckCertificate) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("certificate"), nil
}

func (cmd *CheckCertificate) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckDatabase(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckDatabase {
	cmd := new(CheckDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckDatabase) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckDatabase) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check database '%s' done", extracted)
	} else {
		renv.Log().Verbose("check database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckDatabase) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("database"), nil
}

func (cmd *CheckDatabase) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckDistribution(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckDistribution {
	cmd := new(CheckDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudfront.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckDistribution) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckDistribution) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("check distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckDistribution) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("distribution"), nil
}

func (cmd *CheckDistribution) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckInstance {
	cmd := new(CheckInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("check instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("instance"), nil
}

func (cmd *CheckInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckLoadbalancer {
	cmd := new(CheckLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check loadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("check loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loadbalancer"), nil
}

func (cmd *CheckLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckNatgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckNatgateway {
	cmd := new(CheckNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckNatgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckNatgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check natgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("check natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckNatgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("natgateway"), nil
}

func (cmd *CheckNatgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckNetworkinterface(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckNetworkinterface {
	cmd := new(CheckNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckNetworkinterface) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckNetworkinterface) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("check networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckNetworkinterface) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("networkinterface"), nil
}

func (cmd *CheckNetworkinterface) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckScalinggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckScalinggroup {
	cmd := new(CheckScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckScalinggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckScalinggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("check scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckScalinggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("scalinggroup"), nil
}

func (cmd *CheckScalinggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckSecuritygroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckSecuritygroup {
	cmd := new(CheckSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckSecuritygroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckSecuritygroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("check securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckSecuritygroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("securitygroup"), nil
}

func (cmd *CheckSecuritygroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCheckVolume(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CheckVolume {
	cmd := new(CheckVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CheckVolume) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CheckVolume) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("check volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("check volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("check volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CheckVolume) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("volume"), nil
}

func (cmd *CheckVolume) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCopyImage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CopyImage {
	cmd := new(CopyImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CopyImage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CopyImage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CopyImageInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopyImageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CopyImage(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CopyImage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("copy image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("copy image '%s' done", extracted)
	} else {
		renv.Log().Verbose("copy image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CopyImage) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CopyImageInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopyImageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CopyImage(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CopyImage call took %s", time.Since(start))
			renv.Log().Verbose("dry run: copy image ok")
			return fakeDryRunID("image"), nil
		}
	}

	return nil, err
}

func (cmd *CopyImage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCopySnapshot(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CopySnapshot {
	cmd := new(CopySnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CopySnapshot) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CopySnapshot) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CopySnapshotInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopySnapshotInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CopySnapshot(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CopySnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("copy snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("copy snapshot '%s' done", extracted)
	} else {
		renv.Log().Verbose("copy snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CopySnapshot) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CopySnapshotInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CopySnapshotInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CopySnapshot(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CopySnapshot call took %s", time.Since(start))
			renv.Log().Verbose("dry run: copy snapshot ok")
			return fakeDryRunID("snapshot"), nil
		}
	}

	return nil, err
}

func (cmd *CopySnapshot) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateAccesskey(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateAccesskey {
	cmd := new(CreateAccesskey)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAccesskey) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAccesskey) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateAccessKeyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateAccessKeyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateAccessKey(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreateAccessKey call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create accesskey: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create accesskey '%s' done", extracted)
	} else {
		renv.Log().Verbose("create accesskey done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAccesskey) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("accesskey"), nil
}

func (cmd *CreateAccesskey) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateAlarm(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateAlarm {
	cmd := new(CreateAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatch.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAlarm) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAlarm) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.PutMetricAlarmInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.PutMetricAlarmInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutMetricAlarm(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatch.PutMetricAlarm call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("create alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAlarm) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("alarm"), nil
}

func (cmd *CreateAlarm) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateApigateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateApigateway {
	cmd := new(CreateApigateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = apigatewayv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateApigateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateApigateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &apigatewayv2.CreateApiInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in apigatewayv2.CreateApiInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateApi(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("apigatewayv2.CreateApi call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create apigateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create apigateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("create apigateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateApigateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("apigateway"), nil
}

func (cmd *CreateApigateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateApigatewayroute(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateApigatewayroute {
	cmd := new(CreateApigatewayroute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = apigatewayv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateApigatewayroute) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateApigatewayroute) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &apigatewayv2.CreateRouteInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in apigatewayv2.CreateRouteInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateRoute(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("apigatewayv2.CreateRoute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create apigatewayroute: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create apigatewayroute '%s' done", extracted)
	} else {
		renv.Log().Verbose("create apigatewayroute done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateApigatewayroute) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("apigatewayroute"), nil
}

func (cmd *CreateApigatewayroute) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateApigatewaystage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateApigatewaystage {
	cmd := new(CreateApigatewaystage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = apigatewayv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateApigatewaystage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateApigatewaystage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &apigatewayv2.CreateStageInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in apigatewayv2.CreateStageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateStage(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("apigatewayv2.CreateStage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create apigatewaystage: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create apigatewaystage '%s' done", extracted)
	} else {
		renv.Log().Verbose("create apigatewaystage done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateApigatewaystage) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("apigatewaystage"), nil
}

func (cmd *CreateApigatewaystage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateAppscalingpolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateAppscalingpolicy {
	cmd := new(CreateAppscalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = applicationautoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAppscalingpolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAppscalingpolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.PutScalingPolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.PutScalingPolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutScalingPolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("applicationautoscaling.PutScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create appscalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create appscalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("create appscalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAppscalingpolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("appscalingpolicy"), nil
}

func (cmd *CreateAppscalingpolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateAppscalingtarget(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateAppscalingtarget {
	cmd := new(CreateAppscalingtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = applicationautoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateAppscalingtarget) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateAppscalingtarget) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.RegisterScalableTargetInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.RegisterScalableTargetInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RegisterScalableTarget(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("applicationautoscaling.RegisterScalableTarget call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create appscalingtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create appscalingtarget '%s' done", extracted)
	} else {
		renv.Log().Verbose("create appscalingtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateAppscalingtarget) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("appscalingtarget"), nil
}

func (cmd *CreateAppscalingtarget) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateBucket(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateBucket {
	cmd := new(CreateBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = s3.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateBucket) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateBucket) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.CreateBucketInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.CreateBucketInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateBucket(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("s3.CreateBucket call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create bucket '%s' done", extracted)
	} else {
		renv.Log().Verbose("create bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateBucket) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("bucket"), nil
}

func (cmd *CreateBucket) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateCachecluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateCachecluster {
	cmd := new(CreateCachecluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateCachecluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateCachecluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.CreateCacheClusterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.CreateCacheClusterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateCacheCluster(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.CreateCacheCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create cachecluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create cachecluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("create cachecluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateCachecluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("cachecluster"), nil
}

func (cmd *CreateCachecluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateCachesubnetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateCachesubnetgroup {
	cmd := new(CreateCachesubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateCachesubnetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateCachesubnetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.CreateCacheSubnetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.CreateCacheSubnetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateCacheSubnetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.CreateCacheSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create cachesubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create cachesubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create cachesubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateCachesubnetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("cachesubnetgroup"), nil
}

func (cmd *CreateCachesubnetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateCertificate(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateCertificate {
	cmd := new(CreateCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = acm.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateCertificate) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateCertificate) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create certificate '%s' done", extracted)
	} else {
		renv.Log().Verbose("create certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateCertificate) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("certificate"), nil
}

func (cmd *CreateCertificate) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateClassicLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateClassicLoadbalancer {
	cmd := new(CreateClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateClassicLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateClassicLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.CreateLoadBalancerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.CreateLoadBalancerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoadBalancer(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elb.CreateLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("create classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateClassicLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("classicloadbalancer"), nil
}

func (cmd *CreateClassicLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateConfigrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateConfigrule {
	cmd := new(CreateConfigrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = configservice.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateConfigrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateConfigrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &configservice.PutConfigRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in configservice.PutConfigRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutConfigRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("configservice.PutConfigRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create configrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create configrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("create configrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateConfigrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("configrule"), nil
}

func (cmd *CreateConfigrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateContainercluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateContainercluster {
	cmd := new(CreateContainercluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateContainercluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateContainercluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.CreateClusterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.CreateClusterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateCluster(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ecs.CreateCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create containercluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create containercluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("create containercluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateContainercluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containercluster"), nil
}

func (cmd *CreateContainercluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateDatabase(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateDatabase {
	cmd := new(CreateDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDatabase) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDatabase) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create database '%s' done", extracted)
	} else {
		renv.Log().Verbose("create database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDatabase) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("database"), nil
}

func (cmd *CreateDatabase) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateDbsubnetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateDbsubnetgroup {
	cmd := new(CreateDbsubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDbsubnetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDbsubnetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.CreateDBSubnetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.CreateDBSubnetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateDBSubnetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("rds.CreateDBSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create dbsubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create dbsubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create dbsubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDbsubnetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("dbsubnetgroup"), nil
}

func (cmd *CreateDbsubnetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateDistribution(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateDistribution {
	cmd := new(CreateDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudfront.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDistribution) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDistribution) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("create distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDistribution) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("distribution"), nil
}

func (cmd *CreateDistribution) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateDynamodbtable(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateDynamodbtable {
	cmd := new(CreateDynamodbtable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = dynamodb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateDynamodbtable) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateDynamodbtable) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create dynamodbtable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create dynamodbtable '%s' done", extracted)
	} else {
		renv.Log().Verbose("create dynamodbtable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateDynamodbtable) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("dynamodbtable"), nil
}

func (cmd *CreateDynamodbtable) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateEkscluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateEkscluster {
	cmd := new(CreateEkscluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eks.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateEkscluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateEkscluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create ekscluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create ekscluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("create ekscluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateEkscluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ekscluster"), nil
}

func (cmd *CreateEkscluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateEksnodegroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateEksnodegroup {
	cmd := new(CreateEksnodegroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eks.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateEksnodegroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateEksnodegroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create eksnodegroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create eksnodegroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create eksnodegroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateEksnodegroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eksnodegroup"), nil
}

func (cmd *CreateEksnodegroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateElasticip(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateElasticip {
	cmd := new(CreateElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateElasticip) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateElasticip) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.AllocateAddressInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AllocateAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.AllocateAddress(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.AllocateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("create elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateElasticip) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.AllocateAddressInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.AllocateAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.AllocateAddress(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.AllocateAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create elasticip ok")
			return fakeDryRunID("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *CreateElasticip) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateEventbus(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateEventbus {
	cmd := new(CreateEventbus)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateEventbus) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateEventbus) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.CreateEventBusInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.CreateEventBusInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateEventBus(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.CreateEventBus call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create eventbus: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create eventbus '%s' done", extracted)
	} else {
		renv.Log().Verbose("create eventbus done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateEventbus) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventbus"), nil
}

func (cmd *CreateEventbus) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateEventrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateEventrule {
	cmd := new(CreateEventrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateEventrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateEventrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.PutRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.PutRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.PutRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create eventrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create eventrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("create eventrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateEventrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventrule"), nil
}

func (cmd *CreateEventrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateFilesystem(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateFilesystem {
	cmd := new(CreateFilesystem)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = efs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateFilesystem) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateFilesystem) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &efs.CreateFileSystemInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in efs.CreateFileSystemInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateFileSystem(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("efs.CreateFileSystem call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create filesystem: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create filesystem '%s' done", extracted)
	} else {
		renv.Log().Verbose("create filesystem done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateFilesystem) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("filesystem"), nil
}

func (cmd *CreateFilesystem) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateFunction(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateFunction {
	cmd := new(CreateFunction)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = lambda.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateFunction) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateFunction) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &lambda.CreateFunctionInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in lambda.CreateFunctionInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateFunction(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("lambda.CreateFunction call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create function: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create function '%s' done", extracted)
	} else {
		renv.Log().Verbose("create function done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateFunction) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("function"), nil
}

func (cmd *CreateFunction) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateGroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateGroup {
	cmd := new(CreateGroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateGroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateGroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreateGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create group: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create group '%s' done", extracted)
	} else {
		renv.Log().Verbose("create group done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateGroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("group"), nil
}

func (cmd *CreateGroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateImage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateImage {
	cmd := new(CreateImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateImage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateImage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateImageInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateImageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateImage(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateImage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create image '%s' done", extracted)
	} else {
		renv.Log().Verbose("create image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateImage) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateImageInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateImageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateImage(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateImage call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create image ok")
			return fakeDryRunID("image"), nil
		}
	}

	return nil, err
}

func (cmd *CreateImage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateInstance {
	cmd := new(CreateInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.RunInstancesInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RunInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RunInstances(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.RunInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("create instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.RunInstancesInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RunInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.RunInstances(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.RunInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create instance ok")
			return fakeDryRunID("instance"), nil
		}
	}

	return nil, err
}

func (cmd *CreateInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateInstanceprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateInstanceprofile {
	cmd := new(CreateInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateInstanceprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateInstanceprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateInstanceProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateInstanceProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateInstanceProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreateInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("create instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInstanceprofile) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("instanceprofile"), nil
}

func (cmd *CreateInstanceprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateInternetgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateInternetgateway {
	cmd := new(CreateInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateInternetgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateInternetgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateInternetGatewayInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateInternetGateway(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("create internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateInternetgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateInternetGateway(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create internetgateway ok")
			return fakeDryRunID("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *CreateInternetgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateIpset(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateIpset {
	cmd := new(CreateIpset)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = wafv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateIpset) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateIpset) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create ipset: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create ipset '%s' done", extracted)
	} else {
		renv.Log().Verbose("create ipset done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateIpset) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ipset"), nil
}

func (cmd *CreateIpset) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateKeypair(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateKeypair {
	cmd := new(CreateKeypair)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateKeypair) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateKeypair) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ImportKeyPairInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportKeyPairInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ImportKeyPair(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.ImportKeyPair call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create keypair: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create keypair '%s' done", extracted)
	} else {
		renv.Log().Verbose("create keypair done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateKeypair) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("keypair"), nil
}

func (cmd *CreateKeypair) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateLaunchconfiguration(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateLaunchconfiguration {
	cmd := new(CreateLaunchconfiguration)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLaunchconfiguration) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLaunchconfiguration) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.CreateLaunchConfigurationInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.CreateLaunchConfigurationInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateLaunchConfiguration(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.CreateLaunchConfiguration call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create launchconfiguration: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create launchconfiguration '%s' done", extracted)
	} else {
		renv.Log().Verbose("create launchconfiguration done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLaunchconfiguration) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("launchconfiguration"), nil
}

func (cmd *CreateLaunchconfiguration) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateListener(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateListener {
	cmd := new(CreateListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateListener) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateListener) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateListenerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateListenerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateListener(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.CreateListener call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create listener '%s' done", extracted)
	} else {
		renv.Log().Verbose("create listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateListener) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("listener"), nil
}

func (cmd *CreateListener) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateLoadbalancer {
	cmd := new(CreateLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateLoadBalancerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateLoadBalancerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoadBalancer(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.CreateLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create loadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("create loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loadbalancer"), nil
}

func (cmd *CreateLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateLoggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateLoggroup {
	cmd := new(CreateLoggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatchlogs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLoggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLoggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatchlogs.CreateLogGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatchlogs.CreateLogGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateLogGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatchlogs.CreateLogGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create loggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create loggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create loggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loggroup"), nil
}

func (cmd *CreateLoggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateLoginprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateLoginprofile {
	cmd := new(CreateLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateLoginprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateLoginprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateLoginProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateLoginProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateLoginProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreateLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create loginprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("create loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateLoginprofile) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loginprofile"), nil
}

func (cmd *CreateLoginprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateMfadevice(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateMfadevice {
	cmd := new(CreateMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateMfadevice) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateMfadevice) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("create mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateMfadevice) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("mfadevice"), nil
}

func (cmd *CreateMfadevice) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateNatgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateNatgateway {
	cmd := new(CreateNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateNatgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateNatgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateNatGatewayInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNatGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateNatGateway(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateNatGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create natgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("create natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateNatgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("natgateway"), nil
}

func (cmd *CreateNatgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateNetworkinterface(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateNetworkinterface {
	cmd := new(CreateNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateNetworkinterface) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateNetworkinterface) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateNetworkInterfaceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNetworkInterfaceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateNetworkInterface(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("create networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateNetworkinterface) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateNetworkInterfaceInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateNetworkInterfaceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateNetworkInterface(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateNetworkInterface call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create networkinterface ok")
			return fakeDryRunID("networkinterface"), nil
		}
	}

	return nil, err
}

func (cmd *CreateNetworkinterface) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreatePolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreatePolicy {
	cmd := new(CreatePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreatePolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreatePolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreatePolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreatePolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreatePolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreatePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("create policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreatePolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("policy"), nil
}

func (cmd *CreatePolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateQueue(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateQueue {
	cmd := new(CreateQueue)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sqs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateQueue) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateQueue) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sqs.CreateQueueInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sqs.CreateQueueInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateQueue(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sqs.CreateQueue call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create queue: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create queue '%s' done", extracted)
	} else {
		renv.Log().Verbose("create queue done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateQueue) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("queue"), nil
}

func (cmd *CreateQueue) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateRecord(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateRecord {
	cmd := new(CreateRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = route53.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRecord) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRecord) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create record '%s' done", extracted)
	} else {
		renv.Log().Verbose("create record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRecord) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("record"), nil
}

func (cmd *CreateRecord) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateReplicationgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateReplicationgroup {
	cmd := new(CreateReplicationgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateReplicationgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateReplicationgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.CreateReplicationGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.CreateReplicationGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateReplicationGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.CreateReplicationGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create replicationgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create replicationgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create replicationgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateReplicationgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("replicationgroup"), nil
}

func (cmd *CreateReplicationgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateRepository(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateRepository {
	cmd := new(CreateRepository)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecr.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRepository) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRepository) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecr.CreateRepositoryInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ecr.CreateRepositoryInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateRepository(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ecr.CreateRepository call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create repository: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create repository '%s' done", extracted)
	} else {
		renv.Log().Verbose("create repository done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRepository) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("repository"), nil
}

func (cmd *CreateRepository) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateRole(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateRole {
	cmd := new(CreateRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRole) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRole) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create role '%s' done", extracted)
	} else {
		renv.Log().Verbose("create role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRole) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("role"), nil
}

func (cmd *CreateRole) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateRoute(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateRoute {
	cmd := new(CreateRoute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRoute) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRoute) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateRouteInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateRoute(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateRoute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create route: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create route '%s' done", extracted)
	} else {
		renv.Log().Verbose("create route done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRoute) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateRouteInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateRoute(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateRoute call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create route ok")
			return fakeDryRunID("route"), nil
		}
	}

	return nil, err
}

func (cmd *CreateRoute) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateRoutetable(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateRoutetable {
	cmd := new(CreateRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateRoutetable) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateRoutetable) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateRouteTableInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateRouteTable(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("create routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateRoutetable) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateRouteTable(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create routetable ok")
			return fakeDryRunID("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *CreateRoutetable) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateS3object(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateS3object {
	cmd := new(CreateS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = s3.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateS3object) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateS3object) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create s3object '%s' done", extracted)
	} else {
		renv.Log().Verbose("create s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateS3object) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("s3object"), nil
}

func (cmd *CreateS3object) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateScalinggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateScalinggroup {
	cmd := new(CreateScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateScalinggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateScalinggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.CreateAutoScalingGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.CreateAutoScalingGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateAutoScalingGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.CreateAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateScalinggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("scalinggroup"), nil
}

func (cmd *CreateScalinggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateScalingpolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateScalingpolicy {
	cmd := new(CreateScalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateScalingpolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateScalingpolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.PutScalingPolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.PutScalingPolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutScalingPolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.PutScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create scalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create scalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("create scalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateScalingpolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("scalingpolicy"), nil
}

func (cmd *CreateScalingpolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateSecret(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateSecret {
	cmd := new(CreateSecret)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = secretsmanager.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSecret) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSecret) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &secretsmanager.CreateSecretInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in secretsmanager.CreateSecretInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateSecret(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("secretsmanager.CreateSecret call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create secret: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create secret '%s' done", extracted)
	} else {
		renv.Log().Verbose("create secret done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSecret) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("secret"), nil
}

func (cmd *CreateSecret) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateSecuritygroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateSecuritygroup {
	cmd := new(CreateSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSecuritygroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSecuritygroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSecurityGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSecurityGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateSecurityGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateSecurityGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSecuritygroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateSecurityGroupInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSecurityGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateSecurityGroup(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateSecurityGroup call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create securitygroup ok")
			return fakeDryRunID("securitygroup"), nil
		}
	}

	return nil, err
}

func (cmd *CreateSecuritygroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateSnapshot(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateSnapshot {
	cmd := new(CreateSnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSnapshot) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSnapshot) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSnapshotInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSnapshotInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateSnapshot(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateSnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create snapshot '%s' done", extracted)
	} else {
		renv.Log().Verbose("create snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSnapshot) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateSnapshotInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSnapshotInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateSnapshot(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateSnapshot call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create snapshot ok")
			return fakeDryRunID("snapshot"), nil
		}
	}

	return nil, err
}

func (cmd *CreateSnapshot) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateSsmparameter(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateSsmparameter {
	cmd := new(CreateSsmparameter)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ssm.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSsmparameter) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSsmparameter) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ssm.PutParameterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ssm.PutParameterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutParameter(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ssm.PutParameter call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create ssmparameter: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create ssmparameter '%s' done", extracted)
	} else {
		renv.Log().Verbose("create ssmparameter done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSsmparameter) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ssmparameter"), nil
}

func (cmd *CreateSsmparameter) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateStack(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateStack {
	cmd := new(CreateStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudformation.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateStack) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateStack) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.CreateStackInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.CreateStackInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateStack(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudformation.CreateStack call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create stack '%s' done", extracted)
	} else {
		renv.Log().Verbose("create stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateStack) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("stack"), nil
}

func (cmd *CreateStack) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateStatemachine(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateStatemachine {
	cmd := new(CreateStatemachine)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sfn.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateStatemachine) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateStatemachine) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sfn.CreateStateMachineInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sfn.CreateStateMachineInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateStateMachine(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sfn.CreateStateMachine call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create statemachine: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create statemachine '%s' done", extracted)
	} else {
		renv.Log().Verbose("create statemachine done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateStatemachine) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("statemachine"), nil
}

func (cmd *CreateStatemachine) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateSubnet(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateSubnet {
	cmd := new(CreateSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSubnet) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSubnet) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateSubnetInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSubnetInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateSubnet(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateSubnet call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create subnet '%s' done", extracted)
	} else {
		renv.Log().Verbose("create subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSubnet) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateSubnetInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateSubnetInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateSubnet(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateSubnet call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create subnet ok")
			return fakeDryRunID("subnet"), nil
		}
	}

	return nil, err
}

func (cmd *CreateSubnet) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateSubscription(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateSubscription {
	cmd := new(CreateSubscription)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sns.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateSubscription) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateSubscription) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.SubscribeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.SubscribeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.Subscribe(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sns.Subscribe call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create subscription: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create subscription '%s' done", extracted)
	} else {
		renv.Log().Verbose("create subscription done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateSubscription) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("subscription"), nil
}

func (cmd *CreateSubscription) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateTag(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateTag {
	cmd := new(CreateTag)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTag) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTag) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create tag: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create tag '%s' done", extracted)
	} else {
		renv.Log().Verbose("create tag done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTag) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateTargetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateTargetgroup {
	cmd := new(CreateTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTargetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTargetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.CreateTargetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.CreateTargetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateTargetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.CreateTargetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create targetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("create targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTargetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("targetgroup"), nil
}

func (cmd *CreateTargetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateTopic(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateTopic {
	cmd := new(CreateTopic)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sns.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTopic) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTopic) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.CreateTopicInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.CreateTopicInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateTopic(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sns.CreateTopic call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create topic: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create topic '%s' done", extracted)
	} else {
		renv.Log().Verbose("create topic done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTopic) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("topic"), nil
}

func (cmd *CreateTopic) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateTrail(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateTrail {
	cmd := new(CreateTrail)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudtrail.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateTrail) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateTrail) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudtrail.CreateTrailInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudtrail.CreateTrailInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateTrail(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudtrail.CreateTrail call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create trail: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create trail '%s' done", extracted)
	} else {
		renv.Log().Verbose("create trail done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateTrail) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("trail"), nil
}

func (cmd *CreateTrail) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateUser(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateUser {
	cmd := new(CreateUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateUser) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateUser) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreateUserInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreateUserInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateUser(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreateUser call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create user '%s' done", extracted)
	} else {
		renv.Log().Verbose("create user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateUser) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("user"), nil
}

func (cmd *CreateUser) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateVolume(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateVolume {
	cmd := new(CreateVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateVolume) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateVolume) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateVolumeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateVolume(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("create volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateVolume) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateVolumeInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateVolume(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create volume ok")
			return fakeDryRunID("volume"), nil
		}
	}

	return nil, err
}

func (cmd *CreateVolume) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateVpc(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateVpc {
	cmd := new(CreateVpc)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateVpc) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateVpc) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.CreateVpcInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVpcInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateVpc(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.CreateVpc call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create vpc: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create vpc '%s' done", extracted)
	} else {
		renv.Log().Verbose("create vpc done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateVpc) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateVpcInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateVpcInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.CreateVpc(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.CreateVpc call took %s", time.Since(start))
			renv.Log().Verbose("dry run: create vpc ok")
			return fakeDryRunID("vpc"), nil
		}
	}

	return nil, err
}

func (cmd *CreateVpc) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewCreateZone(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *CreateZone {
	cmd := new(CreateZone)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = route53.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *CreateZone) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *CreateZone) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &route53.CreateHostedZoneInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in route53.CreateHostedZoneInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreateHostedZone(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("route53.CreateHostedZone call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("create zone: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("create zone '%s' done", extracted)
	} else {
		renv.Log().Verbose("create zone done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *CreateZone) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("zone"), nil
}

func (cmd *CreateZone) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteAccesskey(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAccesskey {
	cmd := new(DeleteAccesskey)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAccesskey) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAccesskey) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteAccessKeyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteAccessKeyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAccessKey(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeleteAccessKey call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete accesskey: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete accesskey '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete accesskey done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAccesskey) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("accesskey"), nil
}

func (cmd *DeleteAccesskey) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteAlarm(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAlarm {
	cmd := new(DeleteAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatch.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAlarm) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAlarm) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.DeleteAlarmsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.DeleteAlarmsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAlarms(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatch.DeleteAlarms call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAlarm) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("alarm"), nil
}

func (cmd *DeleteAlarm) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteApigateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteApigateway {
	cmd := new(DeleteApigateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = apigatewayv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteApigateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteApigateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &apigatewayv2.DeleteApiInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in apigatewayv2.DeleteApiInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteApi(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("apigatewayv2.DeleteApi call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete apigateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete apigateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete apigateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteApigateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("apigateway"), nil
}

func (cmd *DeleteApigateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteApigatewayroute(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteApigatewayroute {
	cmd := new(DeleteApigatewayroute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = apigatewayv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteApigatewayroute) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteApigatewayroute) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &apigatewayv2.DeleteRouteInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in apigatewayv2.DeleteRouteInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRoute(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("apigatewayv2.DeleteRoute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete apigatewayroute: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete apigatewayroute '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete apigatewayroute done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteApigatewayroute) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("apigatewayroute"), nil
}

func (cmd *DeleteApigatewayroute) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteApigatewaystage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteApigatewaystage {
	cmd := new(DeleteApigatewaystage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = apigatewayv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteApigatewaystage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteApigatewaystage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &apigatewayv2.DeleteStageInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in apigatewayv2.DeleteStageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteStage(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("apigatewayv2.DeleteStage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete apigatewaystage: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete apigatewaystage '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete apigatewaystage done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteApigatewaystage) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("apigatewaystage"), nil
}

func (cmd *DeleteApigatewaystage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteAppscalingpolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAppscalingpolicy {
	cmd := new(DeleteAppscalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = applicationautoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAppscalingpolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAppscalingpolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.DeleteScalingPolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.DeleteScalingPolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteScalingPolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("applicationautoscaling.DeleteScalingPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete appscalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete appscalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete appscalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAppscalingpolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("appscalingpolicy"), nil
}

func (cmd *DeleteAppscalingpolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteAppscalingtarget(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteAppscalingtarget {
	cmd := new(DeleteAppscalingtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = applicationautoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteAppscalingtarget) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteAppscalingtarget) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &applicationautoscaling.DeregisterScalableTargetInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in applicationautoscaling.DeregisterScalableTargetInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterScalableTarget(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("applicationautoscaling.DeregisterScalableTarget call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete appscalingtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete appscalingtarget '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete appscalingtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteAppscalingtarget) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("appscalingtarget"), nil
}

func (cmd *DeleteAppscalingtarget) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteBucket(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteBucket {
	cmd := new(DeleteBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = s3.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteBucket) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteBucket) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.DeleteBucketInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.DeleteBucketInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteBucket(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("s3.DeleteBucket call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete bucket '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteBucket) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("bucket"), nil
}

func (cmd *DeleteBucket) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteCachecluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteCachecluster {
	cmd := new(DeleteCachecluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteCachecluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteCachecluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.DeleteCacheClusterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.DeleteCacheClusterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCacheCluster(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.DeleteCacheCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete cachecluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete cachecluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete cachecluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteCachecluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("cachecluster"), nil
}

func (cmd *DeleteCachecluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteCachesubnetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteCachesubnetgroup {
	cmd := new(DeleteCachesubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteCachesubnetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteCachesubnetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.DeleteCacheSubnetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.DeleteCacheSubnetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCacheSubnetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.DeleteCacheSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete cachesubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete cachesubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete cachesubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteCachesubnetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("cachesubnetgroup"), nil
}

func (cmd *DeleteCachesubnetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteCertificate(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteCertificate {
	cmd := new(DeleteCertificate)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = acm.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteCertificate) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteCertificate) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &acm.DeleteCertificateInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in acm.DeleteCertificateInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCertificate(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("acm.DeleteCertificate call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete certificate: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete certificate '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete certificate done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteCertificate) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("certificate"), nil
}

func (cmd *DeleteCertificate) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteClassicLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteClassicLoadbalancer {
	cmd := new(DeleteClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteClassicLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteClassicLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.DeleteLoadBalancerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.DeleteLoadBalancerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoadBalancer(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elb.DeleteLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteClassicLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("classicloadbalancer"), nil
}

func (cmd *DeleteClassicLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteConfigrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteConfigrule {
	cmd := new(DeleteConfigrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = configservice.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteConfigrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteConfigrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &configservice.DeleteConfigRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in configservice.DeleteConfigRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteConfigRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("configservice.DeleteConfigRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete configrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete configrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete configrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteConfigrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("configrule"), nil
}

func (cmd *DeleteConfigrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteContainercluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteContainercluster {
	cmd := new(DeleteContainercluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteContainercluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteContainercluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.DeleteClusterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.DeleteClusterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCluster(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ecs.DeleteCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete containercluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete containercluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete containercluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteContainercluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containercluster"), nil
}

func (cmd *DeleteContainercluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteContainertask(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteContainertask {
	cmd := new(DeleteContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteContainertask) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteContainertask) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteContainertask) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteDatabase(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDatabase {
	cmd := new(DeleteDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDatabase) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDatabase) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.DeleteDBInstanceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.DeleteDBInstanceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteDBInstance(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("rds.DeleteDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete database '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDatabase) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("database"), nil
}

func (cmd *DeleteDatabase) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteDbsubnetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDbsubnetgroup {
	cmd := new(DeleteDbsubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDbsubnetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDbsubnetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.DeleteDBSubnetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.DeleteDBSubnetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteDBSubnetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("rds.DeleteDBSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete dbsubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete dbsubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete dbsubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDbsubnetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("dbsubnetgroup"), nil
}

func (cmd *DeleteDbsubnetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteDistribution(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDistribution {
	cmd := new(DeleteDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudfront.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDistribution) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDistribution) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDistribution) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("distribution"), nil
}

func (cmd *DeleteDistribution) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteDynamodbtable(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteDynamodbtable {
	cmd := new(DeleteDynamodbtable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = dynamodb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteDynamodbtable) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteDynamodbtable) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &dynamodb.DeleteTableInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in dynamodb.DeleteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTable(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("dynamodb.DeleteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete dynamodbtable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete dynamodbtable '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete dynamodbtable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteDynamodbtable) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("dynamodbtable"), nil
}

func (cmd *DeleteDynamodbtable) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteEkscluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteEkscluster {
	cmd := new(DeleteEkscluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eks.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteEkscluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteEkscluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eks.DeleteClusterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eks.DeleteClusterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteCluster(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eks.DeleteCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete ekscluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete ekscluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete ekscluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteEkscluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ekscluster"), nil
}

func (cmd *DeleteEkscluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteEksnodegroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteEksnodegroup {
	cmd := new(DeleteEksnodegroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eks.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteEksnodegroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteEksnodegroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eks.DeleteNodegroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eks.DeleteNodegroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNodegroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eks.DeleteNodegroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete eksnodegroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete eksnodegroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete eksnodegroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteEksnodegroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eksnodegroup"), nil
}

func (cmd *DeleteEksnodegroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteElasticip(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteElasticip {
	cmd := new(DeleteElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteElasticip) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteElasticip) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ReleaseAddressInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ReleaseAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ReleaseAddress(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.ReleaseAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteElasticip) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.ReleaseAddressInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ReleaseAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.ReleaseAddress(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.ReleaseAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete elasticip ok")
			return fakeDryRunID("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteElasticip) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteEventbus(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteEventbus {
	cmd := new(DeleteEventbus)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteEventbus) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteEventbus) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.DeleteEventBusInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.DeleteEventBusInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteEventBus(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.DeleteEventBus call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete eventbus: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete eventbus '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete eventbus done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteEventbus) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventbus"), nil
}

func (cmd *DeleteEventbus) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteEventrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteEventrule {
	cmd := new(DeleteEventrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteEventrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteEventrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.DeleteRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.DeleteRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.DeleteRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete eventrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete eventrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete eventrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteEventrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventrule"), nil
}

func (cmd *DeleteEventrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteFilesystem(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteFilesystem {
	cmd := new(DeleteFilesystem)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = efs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteFilesystem) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteFilesystem) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &efs.DeleteFileSystemInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in efs.DeleteFileSystemInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteFileSystem(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("efs.DeleteFileSystem call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete filesystem: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete filesystem '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete filesystem done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteFilesystem) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("filesystem"), nil
}

func (cmd *DeleteFilesystem) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteFunction(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteFunction {
	cmd := new(DeleteFunction)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = lambda.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteFunction) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteFunction) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &lambda.DeleteFunctionInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in lambda.DeleteFunctionInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteFunction(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("lambda.DeleteFunction call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete function: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete function '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete function done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteFunction) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("function"), nil
}

func (cmd *DeleteFunction) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteGroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteGroup {
	cmd := new(DeleteGroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteGroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteGroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeleteGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete group: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete group '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete group done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteGroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("group"), nil
}

func (cmd *DeleteGroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteImage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteImage {
	cmd := new(DeleteImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteImage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteImage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete image '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteImage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteInstance {
	cmd := new(DeleteInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.TerminateInstancesInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.TerminateInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.TerminateInstances(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.TerminateInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.TerminateInstancesInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.TerminateInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.TerminateInstances(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.TerminateInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete instance ok")
			return fakeDryRunID("instance"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteInstanceprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteInstanceprofile {
	cmd := new(DeleteInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteInstanceprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteInstanceprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteInstanceProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteInstanceProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteInstanceProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeleteInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInstanceprofile) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("instanceprofile"), nil
}

func (cmd *DeleteInstanceprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteInternetgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteInternetgateway {
	cmd := new(DeleteInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteInternetgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteInternetgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteInternetGatewayInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteInternetGateway(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteInternetgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteInternetGateway(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete internetgateway ok")
			return fakeDryRunID("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteInternetgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteIpset(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteIpset {
	cmd := new(DeleteIpset)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = wafv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteIpset) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteIpset) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete ipset: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete ipset '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete ipset done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteIpset) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ipset"), nil
}

func (cmd *DeleteIpset) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteKeypair(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteKeypair {
	cmd := new(DeleteKeypair)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteKeypair) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteKeypair) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteKeyPairInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteKeyPairInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteKeyPair(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteKeyPair call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete keypair: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete keypair '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete keypair done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteKeypair) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteKeyPairInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteKeyPairInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteKeyPair(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteKeyPair call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete keypair ok")
			return fakeDryRunID("keypair"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteKeypair) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteLaunchconfiguration(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLaunchconfiguration {
	cmd := new(DeleteLaunchconfiguration)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLaunchconfiguration) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLaunchconfiguration) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeleteLaunchConfigurationInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeleteLaunchConfigurationInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLaunchConfiguration(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.DeleteLaunchConfiguration call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete launchconfiguration: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete launchconfiguration '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete launchconfiguration done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLaunchconfiguration) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("launchconfiguration"), nil
}

func (cmd *DeleteLaunchconfiguration) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteListener(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteListener {
	cmd := new(DeleteListener)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteListener) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteListener) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteListenerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteListenerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteListener(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.DeleteListener call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete listener: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete listener '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete listener done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteListener) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("listener"), nil
}

func (cmd *DeleteListener) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLoadbalancer {
	cmd := new(DeleteLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteLoadBalancerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteLoadBalancerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoadBalancer(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.DeleteLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete loadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete loadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete loadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loadbalancer"), nil
}

func (cmd *DeleteLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteLoggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLoggroup {
	cmd := new(DeleteLoggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatchlogs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLoggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLoggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatchlogs.DeleteLogGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatchlogs.DeleteLogGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLogGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatchlogs.DeleteLogGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete loggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete loggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete loggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loggroup"), nil
}

func (cmd *DeleteLoggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteLoginprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteLoginprofile {
	cmd := new(DeleteLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteLoginprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteLoginprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteLoginProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteLoginProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteLoginProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeleteLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete loginprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteLoginprofile) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loginprofile"), nil
}

func (cmd *DeleteLoginprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteMfadevice(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteMfadevice {
	cmd := new(DeleteMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteMfadevice) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteMfadevice) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteVirtualMFADeviceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteVirtualMFADeviceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVirtualMFADevice(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeleteVirtualMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteMfadevice) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("mfadevice"), nil
}

func (cmd *DeleteMfadevice) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteNatgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteNatgateway {
	cmd := new(DeleteNatgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteNatgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteNatgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteNatGatewayInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNatGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNatGateway(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteNatGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete natgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete natgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete natgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteNatgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("natgateway"), nil
}

func (cmd *DeleteNatgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteNetworkinterface(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteNetworkinterface {
	cmd := new(DeleteNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteNetworkinterface) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteNetworkinterface) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteNetworkInterfaceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNetworkInterfaceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteNetworkInterface(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteNetworkInterface call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteNetworkinterface) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteNetworkInterfaceInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteNetworkInterfaceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteNetworkInterface(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteNetworkInterface call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete networkinterface ok")
			return fakeDryRunID("networkinterface"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteNetworkinterface) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeletePolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeletePolicy {
	cmd := new(DeletePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeletePolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeletePolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeletePolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeletePolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeletePolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeletePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeletePolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("policy"), nil
}

func (cmd *DeletePolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteQueue(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteQueue {
	cmd := new(DeleteQueue)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sqs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteQueue) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteQueue) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sqs.DeleteQueueInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sqs.DeleteQueueInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteQueue(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sqs.DeleteQueue call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete queue: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete queue '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete queue done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteQueue) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("queue"), nil
}

func (cmd *DeleteQueue) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteRecord(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRecord {
	cmd := new(DeleteRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = route53.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRecord) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRecord) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete record '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRecord) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("record"), nil
}

func (cmd *DeleteRecord) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteReplicationgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteReplicationgroup {
	cmd := new(DeleteReplicationgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteReplicationgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteReplicationgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.DeleteReplicationGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.DeleteReplicationGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteReplicationGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.DeleteReplicationGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete replicationgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete replicationgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete replicationgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteReplicationgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("replicationgroup"), nil
}

func (cmd *DeleteReplicationgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteRepository(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRepository {
	cmd := new(DeleteRepository)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecr.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRepository) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRepository) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecr.DeleteRepositoryInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ecr.DeleteRepositoryInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRepository(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ecr.DeleteRepository call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete repository: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete repository '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete repository done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRepository) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("repository"), nil
}

func (cmd *DeleteRepository) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteRole(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRole {
	cmd := new(DeleteRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRole) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRole) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete role '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRole) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("role"), nil
}

func (cmd *DeleteRole) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteRoute(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRoute {
	cmd := new(DeleteRoute)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRoute) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRoute) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteRouteInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRoute(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteRoute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete route: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete route '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete route done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRoute) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteRouteInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteRoute(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteRoute call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete route ok")
			return fakeDryRunID("route"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteRoute) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteRoutetable(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteRoutetable {
	cmd := new(DeleteRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteRoutetable) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteRoutetable) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteRouteTableInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteRouteTable(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteRoutetable) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteRouteTableInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteRouteTable(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete routetable ok")
			return fakeDryRunID("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteRoutetable) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteS3object(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteS3object {
	cmd := new(DeleteS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = s3.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteS3object) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteS3object) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.DeleteObjectInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.DeleteObjectInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteObject(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("s3.DeleteObject call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete s3object '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteS3object) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("s3object"), nil
}

func (cmd *DeleteS3object) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteScalinggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteScalinggroup {
	cmd := new(DeleteScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteScalinggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteScalinggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeleteAutoScalingGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeleteAutoScalingGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteAutoScalingGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.DeleteAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteScalinggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("scalinggroup"), nil
}

func (cmd *DeleteScalinggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteScalingpolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteScalingpolicy {
	cmd := new(DeleteScalingpolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteScalingpolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteScalingpolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.DeletePolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.DeletePolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeletePolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.DeletePolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete scalingpolicy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete scalingpolicy '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete scalingpolicy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteScalingpolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("scalingpolicy"), nil
}

func (cmd *DeleteScalingpolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteSecret(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSecret {
	cmd := new(DeleteSecret)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = secretsmanager.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSecret) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSecret) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &secretsmanager.DeleteSecretInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in secretsmanager.DeleteSecretInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSecret(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("secretsmanager.DeleteSecret call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete secret: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete secret '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete secret done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSecret) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("secret"), nil
}

func (cmd *DeleteSecret) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteSecuritygroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSecuritygroup {
	cmd := new(DeleteSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSecuritygroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSecuritygroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSecurityGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSecurityGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSecurityGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteSecurityGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSecuritygroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteSecurityGroupInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSecurityGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSecurityGroup(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteSecurityGroup call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete securitygroup ok")
			return fakeDryRunID("securitygroup"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteSecuritygroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteSnapshot(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSnapshot {
	cmd := new(DeleteSnapshot)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSnapshot) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSnapshot) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSnapshotInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSnapshotInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSnapshot(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteSnapshot call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete snapshot: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete snapshot '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete snapshot done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSnapshot) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteSnapshotInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSnapshotInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSnapshot(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteSnapshot call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete snapshot ok")
			return fakeDryRunID("snapshot"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteSnapshot) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteSsmparameter(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSsmparameter {
	cmd := new(DeleteSsmparameter)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ssm.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSsmparameter) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSsmparameter) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ssm.DeleteParameterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ssm.DeleteParameterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteParameter(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ssm.DeleteParameter call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete ssmparameter: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete ssmparameter '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete ssmparameter done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSsmparameter) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ssmparameter"), nil
}

func (cmd *DeleteSsmparameter) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteStack(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteStack {
	cmd := new(DeleteStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudformation.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteStack) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteStack) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.DeleteStackInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.DeleteStackInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteStack(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudformation.DeleteStack call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete stack '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteStack) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("stack"), nil
}

func (cmd *DeleteStack) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteStatemachine(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteStatemachine {
	cmd := new(DeleteStatemachine)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sfn.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteStatemachine) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteStatemachine) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sfn.DeleteStateMachineInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sfn.DeleteStateMachineInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteStateMachine(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sfn.DeleteStateMachine call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete statemachine: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete statemachine '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete statemachine done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteStatemachine) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("statemachine"), nil
}

func (cmd *DeleteStatemachine) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteSubnet(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSubnet {
	cmd := new(DeleteSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSubnet) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSubnet) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteSubnetInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSubnetInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteSubnet(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteSubnet call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete subnet '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSubnet) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteSubnetInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteSubnetInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteSubnet(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteSubnet call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete subnet ok")
			return fakeDryRunID("subnet"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteSubnet) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteSubscription(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteSubscription {
	cmd := new(DeleteSubscription)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sns.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteSubscription) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteSubscription) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.UnsubscribeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.UnsubscribeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.Unsubscribe(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sns.Unsubscribe call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete subscription: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete subscription '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete subscription done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteSubscription) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("subscription"), nil
}

func (cmd *DeleteSubscription) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteTag(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTag {
	cmd := new(DeleteTag)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTag) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTag) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete tag: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete tag '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete tag done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTag) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteTargetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTargetgroup {
	cmd := new(DeleteTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTargetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTargetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeleteTargetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeleteTargetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTargetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.DeleteTargetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete targetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTargetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("targetgroup"), nil
}

func (cmd *DeleteTargetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteTopic(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTopic {
	cmd := new(DeleteTopic)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sns.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTopic) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTopic) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sns.DeleteTopicInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sns.DeleteTopicInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTopic(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sns.DeleteTopic call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete topic: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete topic '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete topic done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTopic) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("topic"), nil
}

func (cmd *DeleteTopic) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteTrail(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteTrail {
	cmd := new(DeleteTrail)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudtrail.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteTrail) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteTrail) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudtrail.DeleteTrailInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudtrail.DeleteTrailInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteTrail(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudtrail.DeleteTrail call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete trail: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete trail '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete trail done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteTrail) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("trail"), nil
}

func (cmd *DeleteTrail) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteUser(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteUser {
	cmd := new(DeleteUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteUser) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteUser) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeleteUserInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeleteUserInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteUser(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeleteUser call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete user '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteUser) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("user"), nil
}

func (cmd *DeleteUser) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteVolume(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteVolume {
	cmd := new(DeleteVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteVolume) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteVolume) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteVolumeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVolume(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteVolume) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteVolumeInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteVolume(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete volume ok")
			return fakeDryRunID("volume"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteVolume) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteVpc(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteVpc {
	cmd := new(DeleteVpc)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteVpc) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteVpc) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DeleteVpcInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVpcInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteVpc(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DeleteVpc call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete vpc: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete vpc '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete vpc done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteVpc) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteVpcInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteVpcInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DeleteVpc(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DeleteVpc call took %s", time.Since(start))
			renv.Log().Verbose("dry run: delete vpc ok")
			return fakeDryRunID("vpc"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteVpc) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDeleteZone(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DeleteZone {
	cmd := new(DeleteZone)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = route53.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DeleteZone) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DeleteZone) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &route53.DeleteHostedZoneInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in route53.DeleteHostedZoneInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeleteHostedZone(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("route53.DeleteHostedZone call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("delete zone: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("delete zone '%s' done", extracted)
	} else {
		renv.Log().Verbose("delete zone done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DeleteZone) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("zone"), nil
}

func (cmd *DeleteZone) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachAlarm(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachAlarm {
	cmd := new(DetachAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatch.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachAlarm) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachAlarm) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachAlarm) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("alarm"), nil
}

func (cmd *DetachAlarm) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachClassicLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachClassicLoadbalancer {
	cmd := new(DetachClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachClassicLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachClassicLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.DeregisterInstancesFromLoadBalancerInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.DeregisterInstancesFromLoadBalancerInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterInstancesFromLoadBalancer(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elb.DeregisterInstancesFromLoadBalancer call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachClassicLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("classicloadbalancer"), nil
}

func (cmd *DetachClassicLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachContainertask(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachContainertask {
	cmd := new(DetachContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachContainertask) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachContainertask) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachContainertask) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containertask"), nil
}

func (cmd *DetachContainertask) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachElasticip(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachElasticip {
	cmd := new(DetachElasticip)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachElasticip) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachElasticip) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DisassociateAddressInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DisassociateAddress(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DisassociateAddress call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach elasticip: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach elasticip '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach elasticip done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachElasticip) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DisassociateAddressInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateAddressInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DisassociateAddress(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DisassociateAddress call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach elasticip ok")
			return fakeDryRunID("elasticip"), nil
		}
	}

	return nil, err
}

func (cmd *DetachElasticip) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachEventtarget(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachEventtarget {
	cmd := new(DetachEventtarget)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachEventtarget) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachEventtarget) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.RemoveTargetsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.RemoveTargetsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RemoveTargets(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.RemoveTargets call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach eventtarget: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach eventtarget '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach eventtarget done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachEventtarget) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventtarget"), nil
}

func (cmd *DetachEventtarget) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachInstance {
	cmd := new(DetachInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elbv2.DeregisterTargetsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elbv2.DeregisterTargetsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeregisterTargets(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elbv2.DeregisterTargets call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("instance"), nil
}

func (cmd *DetachInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachInstanceprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachInstanceprofile {
	cmd := new(DetachInstanceprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachInstanceprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachInstanceprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach instanceprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach instanceprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach instanceprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInstanceprofile) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("instanceprofile"), nil
}

func (cmd *DetachInstanceprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachInternetgateway(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachInternetgateway {
	cmd := new(DetachInternetgateway)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachInternetgateway) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachInternetgateway) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DetachInternetGatewayInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DetachInternetGateway(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DetachInternetGateway call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach internetgateway: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach internetgateway '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach internetgateway done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachInternetgateway) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DetachInternetGatewayInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachInternetGatewayInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DetachInternetGateway(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DetachInternetGateway call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach internetgateway ok")
			return fakeDryRunID("internetgateway"), nil
		}
	}

	return nil, err
}

func (cmd *DetachInternetgateway) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachMfadevice(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachMfadevice {
	cmd := new(DetachMfadevice)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachMfadevice) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachMfadevice) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.DeactivateMFADeviceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.DeactivateMFADeviceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DeactivateMFADevice(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.DeactivateMFADevice call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach mfadevice: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach mfadevice '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach mfadevice done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachMfadevice) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("mfadevice"), nil
}

func (cmd *DetachMfadevice) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachNetworkinterface(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachNetworkinterface {
	cmd := new(DetachNetworkinterface)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachNetworkinterface) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachNetworkinterface) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach networkinterface: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach networkinterface '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach networkinterface done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachNetworkinterface) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachPolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachPolicy {
	cmd := new(DetachPolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachPolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachPolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachPolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("policy"), nil
}

func (cmd *DetachPolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachRole(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachRole {
	cmd := new(DetachRole)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachRole) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachRole) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.RemoveRoleFromInstanceProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.RemoveRoleFromInstanceProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RemoveRoleFromInstanceProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.RemoveRoleFromInstanceProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach role: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach role '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach role done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachRole) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("role"), nil
}

func (cmd *DetachRole) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachRoutetable(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachRoutetable {
	cmd := new(DetachRoutetable)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachRoutetable) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachRoutetable) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DisassociateRouteTableInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DisassociateRouteTable(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DisassociateRouteTable call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach routetable: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach routetable '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach routetable done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachRoutetable) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DisassociateRouteTableInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DisassociateRouteTableInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DisassociateRouteTable(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DisassociateRouteTable call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach routetable ok")
			return fakeDryRunID("routetable"), nil
		}
	}

	return nil, err
}

func (cmd *DetachRoutetable) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachSecuritygroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachSecuritygroup {
	cmd := new(DetachSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachSecuritygroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachSecuritygroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachSecuritygroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("securitygroup"), nil
}

func (cmd *DetachSecuritygroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachUser(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachUser {
	cmd := new(DetachUser)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachUser) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachUser) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.RemoveUserFromGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.RemoveUserFromGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RemoveUserFromGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.RemoveUserFromGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach user: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach user '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach user done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachUser) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("user"), nil
}

func (cmd *DetachUser) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewDetachVolume(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *DetachVolume {
	cmd := new(DetachVolume)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *DetachVolume) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *DetachVolume) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.DetachVolumeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DetachVolume(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.DetachVolume call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("detach volume: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("detach volume '%s' done", extracted)
	} else {
		renv.Log().Verbose("detach volume done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *DetachVolume) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DetachVolumeInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DetachVolumeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.DetachVolume(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.DetachVolume call took %s", time.Since(start))
			renv.Log().Verbose("dry run: detach volume ok")
			return fakeDryRunID("volume"), nil
		}
	}

	return nil, err
}

func (cmd *DetachVolume) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewImportImage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *ImportImage {
	cmd := new(ImportImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *ImportImage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *ImportImage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ImportImageInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportImageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ImportImage(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.ImportImage call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("import image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("import image '%s' done", extracted)
	} else {
		renv.Log().Verbose("import image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *ImportImage) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.ImportImageInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ImportImageInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.ImportImage(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.ImportImage call took %s", time.Since(start))
			renv.Log().Verbose("dry run: import image ok")
			return fakeDryRunID("image"), nil
		}
	}

	return nil, err
}

func (cmd *ImportImage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewRestartDatabase(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *RestartDatabase {
	cmd := new(RestartDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *RestartDatabase) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *RestartDatabase) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.RebootDBInstanceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.RebootDBInstanceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RebootDBInstance(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("rds.RebootDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("restart database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("restart database '%s' done", extracted)
	} else {
		renv.Log().Verbose("restart database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *RestartDatabase) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("database"), nil
}

func (cmd *RestartDatabase) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewRestartInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *RestartInstance {
	cmd := new(RestartInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *RestartInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *RestartInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.RebootInstancesInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RebootInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.RebootInstances(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.RebootInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("restart instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("restart instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("restart instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *RestartInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.RebootInstancesInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.RebootInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.RebootInstances(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.RebootInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: restart instance ok")
			return fakeDryRunID("instance"), nil
		}
	}

	return nil, err
}

func (cmd *RestartInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartAlarm(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartAlarm {
	cmd := new(StartAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatch.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartAlarm) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartAlarm) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.EnableAlarmActionsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.EnableAlarmActionsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.EnableAlarmActions(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatch.EnableAlarmActions call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("start alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartAlarm) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("alarm"), nil
}

func (cmd *StartAlarm) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartContainertask(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartContainertask {
	cmd := new(StartContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartContainertask) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartContainertask) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("start containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartContainertask) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containertask"), nil
}

func (cmd *StartContainertask) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartDatabase(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartDatabase {
	cmd := new(StartDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartDatabase) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartDatabase) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.StartDBInstanceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.StartDBInstanceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StartDBInstance(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("rds.StartDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start database '%s' done", extracted)
	} else {
		renv.Log().Verbose("start database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartDatabase) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("database"), nil
}

func (cmd *StartDatabase) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartEventrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartEventrule {
	cmd := new(StartEventrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartEventrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartEventrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.EnableRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.EnableRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.EnableRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.EnableRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start eventrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start eventrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("start eventrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartEventrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventrule"), nil
}

func (cmd *StartEventrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartExecution(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartExecution {
	cmd := new(StartExecution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sfn.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartExecution) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartExecution) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sfn.StartExecutionInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sfn.StartExecutionInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StartExecution(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sfn.StartExecution call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start execution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start execution '%s' done", extracted)
	} else {
		renv.Log().Verbose("start execution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartExecution) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("execution"), nil
}

func (cmd *StartExecution) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartInstance {
	cmd := new(StartInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.StartInstancesInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StartInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StartInstances(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.StartInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("start instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.StartInstancesInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StartInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.StartInstances(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.StartInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: start instance ok")
			return fakeDryRunID("instance"), nil
		}
	}

	return nil, err
}

func (cmd *StartInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStartTrail(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StartTrail {
	cmd := new(StartTrail)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudtrail.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StartTrail) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StartTrail) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudtrail.StartLoggingInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudtrail.StartLoggingInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StartLogging(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudtrail.StartLogging call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("start trail: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("start trail '%s' done", extracted)
	} else {
		renv.Log().Verbose("start trail done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StartTrail) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("trail"), nil
}

func (cmd *StartTrail) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopAlarm(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopAlarm {
	cmd := new(StopAlarm)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatch.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopAlarm) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopAlarm) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatch.DisableAlarmActionsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatch.DisableAlarmActionsInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DisableAlarmActions(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatch.DisableAlarmActions call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop alarm: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop alarm '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop alarm done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopAlarm) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("alarm"), nil
}

func (cmd *StopAlarm) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopContainertask(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopContainertask {
	cmd := new(StopContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopContainertask) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopContainertask) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopContainertask) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containertask"), nil
}

func (cmd *StopContainertask) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopDatabase(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopDatabase {
	cmd := new(StopDatabase)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = rds.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopDatabase) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopDatabase) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &rds.StopDBInstanceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in rds.StopDBInstanceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StopDBInstance(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("rds.StopDBInstance call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop database: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop database '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop database done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopDatabase) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("database"), nil
}

func (cmd *StopDatabase) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopEventrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopEventrule {
	cmd := new(StopEventrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopEventrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopEventrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.DisableRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.DisableRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.DisableRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.DisableRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop eventrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop eventrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop eventrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopEventrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventrule"), nil
}

func (cmd *StopEventrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopExecution(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopExecution {
	cmd := new(StopExecution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sfn.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopExecution) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopExecution) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sfn.StopExecutionInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sfn.StopExecutionInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StopExecution(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sfn.StopExecution call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop execution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop execution '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop execution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopExecution) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("execution"), nil
}

func (cmd *StopExecution) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopInstance {
	cmd := new(StopInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.StopInstancesInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StopInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StopInstances(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.StopInstances call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.StopInstancesInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.StopInstancesInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.StopInstances(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.StopInstances call took %s", time.Since(start))
			renv.Log().Verbose("dry run: stop instance ok")
			return fakeDryRunID("instance"), nil
		}
	}

	return nil, err
}

func (cmd *StopInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewStopTrail(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *StopTrail {
	cmd := new(StopTrail)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudtrail.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *StopTrail) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *StopTrail) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudtrail.StopLoggingInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudtrail.StopLoggingInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.StopLogging(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudtrail.StopLogging call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("stop trail: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("stop trail '%s' done", extracted)
	} else {
		renv.Log().Verbose("stop trail done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *StopTrail) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("trail"), nil
}

func (cmd *StopTrail) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateBucket(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateBucket {
	cmd := new(UpdateBucket)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = s3.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateBucket) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateBucket) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update bucket: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update bucket '%s' done", extracted)
	} else {
		renv.Log().Verbose("update bucket done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateBucket) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("bucket"), nil
}

func (cmd *UpdateBucket) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateCachecluster(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateCachecluster {
	cmd := new(UpdateCachecluster)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateCachecluster) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateCachecluster) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.ModifyCacheClusterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.ModifyCacheClusterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ModifyCacheCluster(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.ModifyCacheCluster call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update cachecluster: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update cachecluster '%s' done", extracted)
	} else {
		renv.Log().Verbose("update cachecluster done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateCachecluster) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("cachecluster"), nil
}

func (cmd *UpdateCachecluster) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateCachesubnetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateCachesubnetgroup {
	cmd := new(UpdateCachesubnetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elasticache.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateCachesubnetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateCachesubnetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elasticache.ModifyCacheSubnetGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elasticache.ModifyCacheSubnetGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ModifyCacheSubnetGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elasticache.ModifyCacheSubnetGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update cachesubnetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update cachesubnetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update cachesubnetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateCachesubnetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("cachesubnetgroup"), nil
}

func (cmd *UpdateCachesubnetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateClassicLoadbalancer(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateClassicLoadbalancer {
	cmd := new(UpdateClassicLoadbalancer)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elb.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateClassicLoadbalancer) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateClassicLoadbalancer) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &elb.ConfigureHealthCheckInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in elb.ConfigureHealthCheckInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ConfigureHealthCheck(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("elb.ConfigureHealthCheck call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update classicloadbalancer: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update classicloadbalancer '%s' done", extracted)
	} else {
		renv.Log().Verbose("update classicloadbalancer done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateClassicLoadbalancer) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("classicloadbalancer"), nil
}

func (cmd *UpdateClassicLoadbalancer) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateConfigrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateConfigrule {
	cmd := new(UpdateConfigrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = configservice.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateConfigrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateConfigrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &configservice.PutConfigRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in configservice.PutConfigRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutConfigRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("configservice.PutConfigRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update configrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update configrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("update configrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateConfigrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("configrule"), nil
}

func (cmd *UpdateConfigrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateContainertask(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateContainertask {
	cmd := new(UpdateContainertask)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ecs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateContainertask) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateContainertask) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ecs.UpdateServiceInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ecs.UpdateServiceInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.UpdateService(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ecs.UpdateService call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update containertask: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update containertask '%s' done", extracted)
	} else {
		renv.Log().Verbose("update containertask done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateContainertask) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("containertask"), nil
}

func (cmd *UpdateContainertask) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateDistribution(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateDistribution {
	cmd := new(UpdateDistribution)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudfront.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateDistribution) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateDistribution) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update distribution: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update distribution '%s' done", extracted)
	} else {
		renv.Log().Verbose("update distribution done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateDistribution) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("distribution"), nil
}

func (cmd *UpdateDistribution) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateEventrule(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateEventrule {
	cmd := new(UpdateEventrule)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = eventbridge.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateEventrule) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateEventrule) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &eventbridge.PutRuleInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in eventbridge.PutRuleInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutRule(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("eventbridge.PutRule call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update eventrule: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update eventrule '%s' done", extracted)
	} else {
		renv.Log().Verbose("update eventrule done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateEventrule) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("eventrule"), nil
}

func (cmd *UpdateEventrule) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateImage(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateImage {
	cmd := new(UpdateImage)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateImage) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateImage) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update image: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update image '%s' done", extracted)
	} else {
		renv.Log().Verbose("update image done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateImage) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateInstance(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateInstance {
	cmd := new(UpdateInstance)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateInstance) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateInstance) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ModifyInstanceAttributeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyInstanceAttributeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ModifyInstanceAttribute(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.ModifyInstanceAttribute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update instance: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update instance '%s' done", extracted)
	} else {
		renv.Log().Verbose("update instance done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateInstance) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.ModifyInstanceAttributeInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifyInstanceAttributeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}

	start := time.Now()
	_, err := cmd.api.ModifyInstanceAttribute(renv.RequestContext(), input)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch code := ae.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
			renv.Log().ExtraVerbosef("dry run: ec2.ModifyInstanceAttribute call took %s", time.Since(start))
			renv.Log().Verbose("dry run: update instance ok")
			return fakeDryRunID("instance"), nil
		}
	}

	return nil, err
}

func (cmd *UpdateInstance) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateIpset(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateIpset {
	cmd := new(UpdateIpset)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = wafv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateIpset) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateIpset) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update ipset: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update ipset '%s' done", extracted)
	} else {
		renv.Log().Verbose("update ipset done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateIpset) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ipset"), nil
}

func (cmd *UpdateIpset) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateLoggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateLoggroup {
	cmd := new(UpdateLoggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudwatchlogs.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateLoggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateLoggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudwatchlogs.PutRetentionPolicyInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudwatchlogs.PutRetentionPolicyInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutRetentionPolicy(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudwatchlogs.PutRetentionPolicy call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update loggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update loggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update loggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateLoggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loggroup"), nil
}

func (cmd *UpdateLoggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateLoginprofile(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateLoginprofile {
	cmd := new(UpdateLoginprofile)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateLoginprofile) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateLoginprofile) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.UpdateLoginProfileInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.UpdateLoginProfileInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.UpdateLoginProfile(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.UpdateLoginProfile call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update loginprofile: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update loginprofile '%s' done", extracted)
	} else {
		renv.Log().Verbose("update loginprofile done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateLoginprofile) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("loginprofile"), nil
}

func (cmd *UpdateLoginprofile) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdatePolicy(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdatePolicy {
	cmd := new(UpdatePolicy)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = iam.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdatePolicy) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdatePolicy) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &iam.CreatePolicyVersionInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in iam.CreatePolicyVersionInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.CreatePolicyVersion(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("iam.CreatePolicyVersion call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update policy: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update policy '%s' done", extracted)
	} else {
		renv.Log().Verbose("update policy done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdatePolicy) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("policy"), nil
}

func (cmd *UpdatePolicy) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateRecord(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateRecord {
	cmd := new(UpdateRecord)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = route53.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateRecord) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateRecord) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update record: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update record '%s' done", extracted)
	} else {
		renv.Log().Verbose("update record done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateRecord) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("record"), nil
}

func (cmd *UpdateRecord) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateS3object(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateS3object {
	cmd := new(UpdateS3object)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = s3.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateS3object) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateS3object) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &s3.PutObjectAclInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in s3.PutObjectAclInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutObjectAcl(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("s3.PutObjectAcl call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update s3object: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update s3object '%s' done", extracted)
	} else {
		renv.Log().Verbose("update s3object done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateS3object) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("s3object"), nil
}

func (cmd *UpdateS3object) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateScalinggroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateScalinggroup {
	cmd := new(UpdateScalinggroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = autoscaling.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateScalinggroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateScalinggroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in autoscaling.UpdateAutoScalingGroupInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.UpdateAutoScalingGroup(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("autoscaling.UpdateAutoScalingGroup call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update scalinggroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update scalinggroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update scalinggroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateScalinggroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("scalinggroup"), nil
}

func (cmd *UpdateScalinggroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateSecret(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateSecret {
	cmd := new(UpdateSecret)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = secretsmanager.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateSecret) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateSecret) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &secretsmanager.UpdateSecretInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in secretsmanager.UpdateSecretInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.UpdateSecret(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("secretsmanager.UpdateSecret call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update secret: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update secret '%s' done", extracted)
	} else {
		renv.Log().Verbose("update secret done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSecret) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("secret"), nil
}

func (cmd *UpdateSecret) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateSecuritygroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateSecuritygroup {
	cmd := new(UpdateSecuritygroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateSecuritygroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateSecuritygroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update securitygroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update securitygroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update securitygroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSecuritygroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateSsmparameter(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateSsmparameter {
	cmd := new(UpdateSsmparameter)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ssm.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateSsmparameter) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateSsmparameter) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ssm.PutParameterInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ssm.PutParameterInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.PutParameter(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ssm.PutParameter call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update ssmparameter: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update ssmparameter '%s' done", extracted)
	} else {
		renv.Log().Verbose("update ssmparameter done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSsmparameter) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("ssmparameter"), nil
}

func (cmd *UpdateSsmparameter) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateStack(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateStack {
	cmd := new(UpdateStack)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = cloudformation.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateStack) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateStack) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &cloudformation.UpdateStackInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in cloudformation.UpdateStackInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.UpdateStack(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("cloudformation.UpdateStack call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update stack: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update stack '%s' done", extracted)
	} else {
		renv.Log().Verbose("update stack done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateStack) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("stack"), nil
}

func (cmd *UpdateStack) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateStatemachine(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateStatemachine {
	cmd := new(UpdateStatemachine)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = sfn.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateStatemachine) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateStatemachine) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &sfn.UpdateStateMachineInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in sfn.UpdateStateMachineInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.UpdateStateMachine(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("sfn.UpdateStateMachine call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update statemachine: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update statemachine '%s' done", extracted)
	} else {
		renv.Log().Verbose("update statemachine done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateStatemachine) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("statemachine"), nil
}

func (cmd *UpdateStatemachine) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateSubnet(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateSubnet {
	cmd := new(UpdateSubnet)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = ec2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateSubnet) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateSubnet) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	input := &ec2.ModifySubnetAttributeInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.ModifySubnetAttributeInput: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.ModifySubnetAttribute(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("ec2.ModifySubnetAttribute call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update subnet: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update subnet '%s' done", extracted)
	} else {
		renv.Log().Verbose("update subnet done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateSubnet) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("subnet"), nil
}

func (cmd *UpdateSubnet) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

func NewUpdateTargetgroup(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *UpdateTargetgroup {
	cmd := new(UpdateTargetgroup)
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = elbv2.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *UpdateTargetgroup) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *UpdateTargetgroup) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}

	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}

	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("update targetgroup: AWS command returned nil output")
		}
	}

	if extracted != nil {
		renv.Log().Verbosef("update targetgroup '%s' done", extracted)
	} else {
		renv.Log().Verbose("update targetgroup done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

func (cmd *UpdateTargetgroup) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("targetgroup"), nil
}

func (cmd *UpdateTargetgroup) inject(params map[string]any) error {
	return structSetter(cmd, params)
}

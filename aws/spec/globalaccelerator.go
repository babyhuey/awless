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
	"fmt"
	"math/rand"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

// Global Accelerator is a global service: its control plane lives in us-west-2 regardless of
// where the endpoints are.
//
// The listener entity is acceleratorlistener, because ELBv2 already owns "listener" and the
// two are unrelated.

type CreateAccelerator struct {
	_       string `action:"create" entity:"accelerator" awsAPI:"globalaccelerator" awsCall:"CreateAccelerator" awsInput:"globalaccelerator.CreateAcceleratorInput" awsOutput:"globalaccelerator.CreateAcceleratorOutput"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *globalaccelerator.Client
	Name    *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Enabled *bool   `awsName:"Enabled" awsType:"awsbool" templateName:"enabled"`
	IPType  *string `awsName:"IpAddressType" awsType:"awsstr" templateName:"ip-type"`
	// Required by the API. Generated rather than asked for: it exists to make a retried
	// call idempotent, which is machinery rather than a decision the user should make.
	Token *string `awsName:"IdempotencyToken" awsType:"awsstr" templateName:"token"`
}

func (cmd *CreateAccelerator) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("enabled"), "ip-type", "token"),
	))
}

func (cmd *CreateAccelerator) BeforeRun(renv env.Running) error {
	if StringValue(cmd.Token) == "" {
		cmd.Token = String(fmt.Sprintf("awless-%d", rand.Int63())) //nolint:gosec // idempotency key, not a secret
	}
	return nil
}

func (cmd *CreateAccelerator) ExtractResult(i any) string {
	out, ok := i.(*globalaccelerator.CreateAcceleratorOutput)
	if !ok || out.Accelerator == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Accelerator.AcceleratorArn)
}

// An accelerator must be disabled before it can be deleted, which AWS enforces rather than
// doing for you — hence `update accelerator enabled=false` first.
type DeleteAccelerator struct {
	_      string `action:"delete" entity:"accelerator" awsAPI:"globalaccelerator" awsCall:"DeleteAccelerator" awsInput:"globalaccelerator.DeleteAcceleratorInput" awsOutput:"globalaccelerator.DeleteAcceleratorOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *globalaccelerator.Client
	Arn    *string `awsName:"AcceleratorArn" awsType:"awsstr" templateName:"arn"`
}

func (cmd *DeleteAccelerator) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("arn")))
}

type UpdateAccelerator struct {
	_       string `action:"update" entity:"accelerator" awsAPI:"globalaccelerator" awsCall:"UpdateAccelerator" awsInput:"globalaccelerator.UpdateAcceleratorInput" awsOutput:"globalaccelerator.UpdateAcceleratorOutput"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *globalaccelerator.Client
	Arn     *string `awsName:"AcceleratorArn" awsType:"awsstr" templateName:"arn"`
	Name    *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Enabled *bool   `awsName:"Enabled" awsType:"awsbool" templateName:"enabled"`
	IPType  *string `awsName:"IpAddressType" awsType:"awsstr" templateName:"ip-type"`
}

func (cmd *UpdateAccelerator) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("arn"),
		params.AtLeastOneOf(params.Key("name"), params.Key("enabled"), params.Key("ip-type")),
	))
}

func (cmd *UpdateAccelerator) ExtractResult(i any) string {
	return StringValue(cmd.Arn)
}

type CreateAcceleratorlistener struct {
	_           string `action:"create" entity:"acceleratorlistener" awsAPI:"globalaccelerator" awsCall:"CreateListener" awsInput:"globalaccelerator.CreateListenerInput" awsOutput:"globalaccelerator.CreateListenerOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *globalaccelerator.Client
	Accelerator *string `awsName:"AcceleratorArn" awsType:"awsstr" templateName:"accelerator"`
	Protocol    *string `awsName:"Protocol" awsType:"awsstr" templateName:"protocol"`
	// Port ranges are a list of from/to pairs, so they come from a file rather than being
	// squeezed into a flag.
	PortsFile *string `awsName:"PortRanges" awsType:"awsfiletostruct" templateName:"ports-file"`
	Affinity  *string `awsName:"ClientAffinity" awsType:"awsstr" templateName:"client-affinity"`
	Token     *string `awsName:"IdempotencyToken" awsType:"awsstr" templateName:"token"`
}

func (cmd *CreateAcceleratorlistener) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("accelerator"), params.Key("protocol"), params.Key("ports-file"),
		params.Opt("client-affinity", "token"),
	))
}

func (cmd *CreateAcceleratorlistener) BeforeRun(renv env.Running) error {
	if StringValue(cmd.Token) == "" {
		cmd.Token = String(fmt.Sprintf("awless-%d", rand.Int63())) //nolint:gosec // idempotency key, not a secret
	}
	return nil
}

func (cmd *CreateAcceleratorlistener) ExtractResult(i any) string {
	out, ok := i.(*globalaccelerator.CreateListenerOutput)
	if !ok || out.Listener == nil {
		return ""
	}
	return awssdk.ToString(out.Listener.ListenerArn)
}

type DeleteAcceleratorlistener struct {
	_      string `action:"delete" entity:"acceleratorlistener" awsAPI:"globalaccelerator" awsCall:"DeleteListener" awsInput:"globalaccelerator.DeleteListenerInput" awsOutput:"globalaccelerator.DeleteListenerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *globalaccelerator.Client
	Arn    *string `awsName:"ListenerArn" awsType:"awsstr" templateName:"arn"`
}

func (cmd *DeleteAcceleratorlistener) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("arn")))
}

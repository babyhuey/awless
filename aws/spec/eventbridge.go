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
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateEventbus struct {
	_           string `action:"create" entity:"eventbus" awsAPI:"eventbridge" awsCall:"CreateEventBus" awsInput:"eventbridge.CreateEventBusInput" awsOutput:"eventbridge.CreateEventBusOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *eventbridge.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	// An event bus backed by a partner event source is created from that source's
	// name rather than as an empty custom bus.
	Source *string `awsName:"EventSourceName" awsType:"awsstr" templateName:"source"`
	KmsKey *string `awsName:"KmsKeyIdentifier" awsType:"awsstr" templateName:"kms-key"`
}

func (cmd *CreateEventbus) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("description"), "source", "kms-key"),
	))
}

// CreateEventBus returns the bus ARN rather than its name, but the name is what the
// graph and the delete command key on.
func (cmd *CreateEventbus) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteEventbus struct {
	_      string `action:"delete" entity:"eventbus" awsAPI:"eventbridge" awsCall:"DeleteEventBus" awsInput:"eventbridge.DeleteEventBusInput" awsOutput:"eventbridge.DeleteEventBusOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *eventbridge.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteEventbus) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// PutRule both creates and updates, so create and update map to the same call. A rule
// needs either an event pattern or a schedule expression; AWS rejects a rule with
// neither, and the two are independent rather than exclusive.
type CreateEventrule struct {
	_           string `action:"create" entity:"eventrule" awsAPI:"eventbridge" awsCall:"PutRule" awsInput:"eventbridge.PutRuleInput" awsOutput:"eventbridge.PutRuleOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *eventbridge.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	EventBus    *string `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
	Pattern     *string `awsName:"EventPattern" awsType:"awsstr" templateName:"pattern"`
	Schedule    *string `awsName:"ScheduleExpression" awsType:"awsstr" templateName:"schedule"`
	Role        *string `awsName:"RoleArn" awsType:"awsstr" templateName:"role"`
	State       *string `awsName:"State" awsType:"awsstr" templateName:"state"`
}

func (cmd *CreateEventrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.AtLeastOneOf(params.Key("pattern"), params.Key("schedule")),
		params.Opt(params.Suggested("description"), "eventbus", "role", "state"),
	))
}

func (cmd *CreateEventrule) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type UpdateEventrule struct {
	_           string `action:"update" entity:"eventrule" awsAPI:"eventbridge" awsCall:"PutRule" awsInput:"eventbridge.PutRuleInput" awsOutput:"eventbridge.PutRuleOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *eventbridge.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	EventBus    *string `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
	Pattern     *string `awsName:"EventPattern" awsType:"awsstr" templateName:"pattern"`
	Schedule    *string `awsName:"ScheduleExpression" awsType:"awsstr" templateName:"schedule"`
	Role        *string `awsName:"RoleArn" awsType:"awsstr" templateName:"role"`
	State       *string `awsName:"State" awsType:"awsstr" templateName:"state"`
}

// PutRule replaces the whole rule, so an update that omits the pattern or schedule
// would silently clear it. One of them stays required for that reason.
func (cmd *UpdateEventrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.AtLeastOneOf(params.Key("pattern"), params.Key("schedule")),
		params.Opt("description", "eventbus", "role", "state"),
	))
}

type DeleteEventrule struct {
	_        string `action:"delete" entity:"eventrule" awsAPI:"eventbridge" awsCall:"DeleteRule" awsInput:"eventbridge.DeleteRuleInput" awsOutput:"eventbridge.DeleteRuleOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *eventbridge.Client
	Name     *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	EventBus *string `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
	// A rule with targets still attached cannot be deleted without this.
	Force *bool `awsName:"Force" awsType:"awsbool" templateName:"force"`
}

func (cmd *DeleteEventrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("eventbus", "force"),
	))
}

type StartEventrule struct {
	_        string `action:"start" entity:"eventrule" awsAPI:"eventbridge" awsCall:"EnableRule" awsInput:"eventbridge.EnableRuleInput" awsOutput:"eventbridge.EnableRuleOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *eventbridge.Client
	Name     *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	EventBus *string `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
}

func (cmd *StartEventrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Opt("eventbus")))
}

type StopEventrule struct {
	_        string `action:"stop" entity:"eventrule" awsAPI:"eventbridge" awsCall:"DisableRule" awsInput:"eventbridge.DisableRuleInput" awsOutput:"eventbridge.DisableRuleOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *eventbridge.Client
	Name     *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	EventBus *string `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
}

func (cmd *StopEventrule) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Opt("eventbus")))
}

// Targets are attached to a rule rather than existing independently, so they are an
// attach/detach pair rather than a resource with its own lifecycle. PutTargets takes a
// slice of structs; the indexed paths below all address element zero, which the setter
// merges into a single target.
type AttachEventtarget struct {
	_        string `action:"attach" entity:"eventtarget" awsAPI:"eventbridge" awsCall:"PutTargets" awsInput:"eventbridge.PutTargetsInput" awsOutput:"eventbridge.PutTargetsOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *eventbridge.Client
	Rule     *string `awsName:"Rule" awsType:"awsstr" templateName:"rule"`
	ID       *string `awsName:"Targets[0]Id" awsType:"awsslicestruct" templateName:"id"`
	Arn      *string `awsName:"Targets[0]Arn" awsType:"awsslicestruct" templateName:"arn"`
	Role     *string `awsName:"Targets[0]RoleArn" awsType:"awsslicestruct" templateName:"role"`
	Input    *string `awsName:"Targets[0]Input" awsType:"awsslicestruct" templateName:"input"`
	EventBus *string `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
}

func (cmd *AttachEventtarget) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("rule"), params.Key("id"), params.Key("arn"),
		params.Opt("role", "input", "eventbus"),
	))
}

func (cmd *AttachEventtarget) ExtractResult(i any) string {
	return StringValue(cmd.ID)
}

type DetachEventtarget struct {
	_        string `action:"detach" entity:"eventtarget" awsAPI:"eventbridge" awsCall:"RemoveTargets" awsInput:"eventbridge.RemoveTargetsInput" awsOutput:"eventbridge.RemoveTargetsOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *eventbridge.Client
	Rule     *string   `awsName:"Rule" awsType:"awsstr" templateName:"rule"`
	IDs      []*string `awsName:"Ids" awsType:"awsstringslice" templateName:"id"`
	EventBus *string   `awsName:"EventBusName" awsType:"awsstr" templateName:"eventbus"`
	Force    *bool     `awsName:"Force" awsType:"awsbool" templateName:"force"`
}

func (cmd *DetachEventtarget) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("rule"), params.Key("id"),
		params.Opt("eventbus", "force"),
	))
}

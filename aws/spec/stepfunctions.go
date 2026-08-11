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
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateStatemachine struct {
	_      string `action:"create" entity:"statemachine" awsAPI:"sfn" awsCall:"CreateStateMachine" awsInput:"sfn.CreateStateMachineInput" awsOutput:"sfn.CreateStateMachineOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sfn.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Role   *string `awsName:"RoleArn" awsType:"awsstr" templateName:"role"`
	// The Amazon States Language document. Taken from a file because it is JSON of
	// non-trivial size, which the template grammar cannot carry inline without
	// single-quoting the whole thing.
	DefinitionFile *string `awsName:"Definition" awsType:"awsfiletostring" templateName:"definition-file"`
	Type           *string `awsName:"Type" awsType:"awsstr" templateName:"type"`
	Publish        *bool   `awsName:"Publish" awsType:"awsbool" templateName:"publish"`
}

func (cmd *CreateStatemachine) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("role"), params.Key("definition-file"),
		params.Opt(params.Suggested("type"), "publish"),
	))
}

func (cmd *CreateStatemachine) ExtractResult(i any) string {
	out, ok := i.(*sfn.CreateStateMachineOutput)
	if !ok || out.StateMachineArn == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.StateMachineArn)
}

type DeleteStatemachine struct {
	_      string `action:"delete" entity:"statemachine" awsAPI:"sfn" awsCall:"DeleteStateMachine" awsInput:"sfn.DeleteStateMachineInput" awsOutput:"sfn.DeleteStateMachineOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sfn.Client
	Arn    *string `awsName:"StateMachineArn" awsType:"awsstr" templateName:"arn"`
}

func (cmd *DeleteStatemachine) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("arn")))
}

type UpdateStatemachine struct {
	_              string `action:"update" entity:"statemachine" awsAPI:"sfn" awsCall:"UpdateStateMachine" awsInput:"sfn.UpdateStateMachineInput" awsOutput:"sfn.UpdateStateMachineOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *sfn.Client
	Arn            *string `awsName:"StateMachineArn" awsType:"awsstr" templateName:"arn"`
	Role           *string `awsName:"RoleArn" awsType:"awsstr" templateName:"role"`
	DefinitionFile *string `awsName:"Definition" awsType:"awsfiletostring" templateName:"definition-file"`
	Publish        *bool   `awsName:"Publish" awsType:"awsbool" templateName:"publish"`
}

func (cmd *UpdateStatemachine) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("arn"),
		params.AtLeastOneOf(params.Key("definition-file"), params.Key("role")),
		params.Opt("publish"),
	))
}

// An execution is started and stopped rather than created and deleted, and is not
// modeled as a resource: listing executions needs a state machine ARN per call, so it
// would be an N+1 over every machine on each sync, for run history rather than
// infrastructure.
type StartExecution struct {
	_            string `action:"start" entity:"execution" awsAPI:"sfn" awsCall:"StartExecution" awsInput:"sfn.StartExecutionInput" awsOutput:"sfn.StartExecutionOutput"`
	logger       *logger.Logger
	graph        cloud.GraphAPI
	api          *sfn.Client
	StateMachine *string `awsName:"StateMachineArn" awsType:"awsstr" templateName:"statemachine"`
	Name         *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Input        *string `awsName:"Input" awsType:"awsstr" templateName:"input"`
}

func (cmd *StartExecution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("statemachine"),
		params.Opt(params.Suggested("name"), "input"),
	))
}

func (cmd *StartExecution) ExtractResult(i any) string {
	out, ok := i.(*sfn.StartExecutionOutput)
	if !ok || out.ExecutionArn == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.ExecutionArn)
}

type StopExecution struct {
	_      string `action:"stop" entity:"execution" awsAPI:"sfn" awsCall:"StopExecution" awsInput:"sfn.StopExecutionInput" awsOutput:"sfn.StopExecutionOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sfn.Client
	Arn    *string `awsName:"ExecutionArn" awsType:"awsstr" templateName:"arn"`
	Cause  *string `awsName:"Cause" awsType:"awsstr" templateName:"cause"`
	Error  *string `awsName:"Error" awsType:"awsstr" templateName:"error"`
}

func (cmd *StopExecution) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("arn"),
		params.Opt("cause", "error"),
	))
}

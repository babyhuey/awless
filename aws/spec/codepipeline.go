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
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// Pipelines are listed, deleted, and run. Creation is not offered: CreatePipeline takes a
// whole PipelineDeclaration — role, artifact store and every stage with its actions —
// which is a document rather than a set of flags, and the reflective setter has no way to
// read one. Pipelines are normally defined in CloudFormation or Terraform anyway; what is
// missing from those tools is operating them, which is what these commands do.
//
// `start pipeline` rather than `start execution` because the execution entity belongs to
// Step Functions, and an execution here has no independent existence: it is started by
// naming the pipeline.

type DeletePipeline struct {
	_      string `action:"delete" entity:"pipeline" awsAPI:"codepipeline" awsCall:"DeletePipeline" awsInput:"codepipeline.DeletePipelineInput" awsOutput:"codepipeline.DeletePipelineOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *codepipeline.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeletePipeline) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type StartPipeline struct {
	_      string `action:"start" entity:"pipeline" awsAPI:"codepipeline" awsCall:"StartPipelineExecution" awsInput:"codepipeline.StartPipelineExecutionInput" awsOutput:"codepipeline.StartPipelineExecutionOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *codepipeline.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *StartPipeline) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// The execution id is what `stop pipeline` needs, so it is the useful result.
func (cmd *StartPipeline) ExtractResult(i any) string {
	out, ok := i.(*codepipeline.StartPipelineExecutionOutput)
	if !ok || out.PipelineExecutionId == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.PipelineExecutionId)
}

type StopPipeline struct {
	_         string `action:"stop" entity:"pipeline" awsAPI:"codepipeline" awsCall:"StopPipelineExecution" awsInput:"codepipeline.StopPipelineExecutionInput" awsOutput:"codepipeline.StopPipelineExecutionOutput"`
	logger    *logger.Logger
	graph     cloud.GraphAPI
	api       *codepipeline.Client
	Name      *string `awsName:"PipelineName" awsType:"awsstr" templateName:"name"`
	Execution *string `awsName:"PipelineExecutionId" awsType:"awsstr" templateName:"execution"`
	// Abandon stops tracking the execution immediately instead of waiting for
	// in-flight actions to finish, which can leave a deploy half-applied.
	Abandon *bool   `awsName:"Abandon" awsType:"awsbool" templateName:"abandon"`
	Reason  *string `awsName:"Reason" awsType:"awsstr" templateName:"reason"`
}

func (cmd *StopPipeline) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("execution"),
		params.Opt("abandon", "reason"),
	))
}

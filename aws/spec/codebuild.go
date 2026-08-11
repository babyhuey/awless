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
	"github.com/aws/aws-sdk-go-v2/service/codebuild"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// CreateProject nests source, environment and artifacts in their own structs, which
// dotted awsName paths reach; the setter builds the intermediates. All three are required
// by AWS, so their type fields are required here rather than optional.
type CreateBuildproject struct {
	_              string `action:"create" entity:"buildproject" awsAPI:"codebuild" awsCall:"CreateProject" awsInput:"codebuild.CreateProjectInput" awsOutput:"codebuild.CreateProjectOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *codebuild.Client
	Name           *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Role           *string `awsName:"ServiceRole" awsType:"awsstr" templateName:"role"`
	Description    *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	SourceType     *string `awsName:"Source.Type" awsType:"awsstr" templateName:"source-type"`
	SourceLocation *string `awsName:"Source.Location" awsType:"awsstr" templateName:"source-location"`
	Buildspec      *string `awsName:"Source.Buildspec" awsType:"awsstr" templateName:"buildspec"`
	// The environment defaults reflect what most projects use: a Linux container on
	// the standard managed image. Compute type has no safe default because it decides
	// the bill, so it is suggested rather than assumed.
	EnvType      *string `awsName:"Environment.Type" awsType:"awsstr" templateName:"env-type"`
	Image        *string `awsName:"Environment.Image" awsType:"awsstr" templateName:"image"`
	ComputeType  *string `awsName:"Environment.ComputeType" awsType:"awsstr" templateName:"compute-type"`
	ArtifactType *string `awsName:"Artifacts.Type" awsType:"awsstr" templateName:"artifact-type"`
	Timeout      *int64  `awsName:"TimeoutInMinutes" awsType:"awsint64" templateName:"timeout"`
}

func (cmd *CreateBuildproject) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("role"),
		params.Key("source-type"), params.Key("env-type"), params.Key("image"),
		params.Key("compute-type"), params.Key("artifact-type"),
		params.Opt(
			params.Suggested("source-location", "buildspec"),
			"description", "timeout",
		),
	))
}

func (cmd *CreateBuildproject) ExtractResult(i any) string {
	out, ok := i.(*codebuild.CreateProjectOutput)
	if !ok || out.Project == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Project.Name)
}

// UpdateProject leaves omitted fields alone, unlike the put-style updates elsewhere here,
// so nothing beyond the name is required.
type UpdateBuildproject struct {
	_              string `action:"update" entity:"buildproject" awsAPI:"codebuild" awsCall:"UpdateProject" awsInput:"codebuild.UpdateProjectInput" awsOutput:"codebuild.UpdateProjectOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *codebuild.Client
	Name           *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Role           *string `awsName:"ServiceRole" awsType:"awsstr" templateName:"role"`
	Description    *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	SourceType     *string `awsName:"Source.Type" awsType:"awsstr" templateName:"source-type"`
	SourceLocation *string `awsName:"Source.Location" awsType:"awsstr" templateName:"source-location"`
	Buildspec      *string `awsName:"Source.Buildspec" awsType:"awsstr" templateName:"buildspec"`
	EnvType        *string `awsName:"Environment.Type" awsType:"awsstr" templateName:"env-type"`
	Image          *string `awsName:"Environment.Image" awsType:"awsstr" templateName:"image"`
	ComputeType    *string `awsName:"Environment.ComputeType" awsType:"awsstr" templateName:"compute-type"`
	Timeout        *int64  `awsName:"TimeoutInMinutes" awsType:"awsint64" templateName:"timeout"`
}

func (cmd *UpdateBuildproject) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(
			"role", "description", "source-type", "source-location", "buildspec",
			"env-type", "image", "compute-type", "timeout",
		),
	))
}

func (cmd *UpdateBuildproject) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteBuildproject struct {
	_      string `action:"delete" entity:"buildproject" awsAPI:"codebuild" awsCall:"DeleteProject" awsInput:"codebuild.DeleteProjectInput" awsOutput:"codebuild.DeleteProjectOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *codebuild.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteBuildproject) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// Builds are started and stopped through the project, the same shape as CodePipeline: a
// build has no existence until one is started, and the id needed to stop it comes back
// from the start.
type StartBuildproject struct {
	_      string `action:"start" entity:"buildproject" awsAPI:"codebuild" awsCall:"StartBuild" awsInput:"codebuild.StartBuildInput" awsOutput:"codebuild.StartBuildOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *codebuild.Client
	Name   *string `awsName:"ProjectName" awsType:"awsstr" templateName:"name"`
	// Overrides for a one-off build, so a branch or a different buildspec can be built
	// without editing the project.
	SourceVersion  *string `awsName:"SourceVersion" awsType:"awsstr" templateName:"source-version"`
	BuildspecOverr *string `awsName:"BuildspecOverride" awsType:"awsstr" templateName:"buildspec"`
}

func (cmd *StartBuildproject) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("source-version", "buildspec"),
	))
}

func (cmd *StartBuildproject) ExtractResult(i any) string {
	out, ok := i.(*codebuild.StartBuildOutput)
	if !ok || out.Build == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Build.Id)
}

type StopBuildproject struct {
	_      string `action:"stop" entity:"buildproject" awsAPI:"codebuild" awsCall:"StopBuild" awsInput:"codebuild.StopBuildInput" awsOutput:"codebuild.StopBuildOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *codebuild.Client
	// The build id, not the project name: StopBuild addresses a single running build.
	Build *string `awsName:"Id" awsType:"awsstr" templateName:"build"`
}

func (cmd *StopBuildproject) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("build")))
}

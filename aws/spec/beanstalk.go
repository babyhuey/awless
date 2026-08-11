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
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateApplication struct {
	_           string `action:"create" entity:"application" awsAPI:"elasticbeanstalk" awsCall:"CreateApplication" awsInput:"elasticbeanstalk.CreateApplicationInput" awsOutput:"elasticbeanstalk.CreateApplicationOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *elasticbeanstalk.Client
	Name        *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"name"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateApplication) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("description")),
	))
}

func (cmd *CreateApplication) ExtractResult(i any) string {
	out, ok := i.(*elasticbeanstalk.CreateApplicationOutput)
	if !ok || out.Application == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Application.ApplicationName)
}

type DeleteApplication struct {
	_      string `action:"delete" entity:"application" awsAPI:"elasticbeanstalk" awsCall:"DeleteApplication" awsInput:"elasticbeanstalk.DeleteApplicationInput" awsOutput:"elasticbeanstalk.DeleteApplicationOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *elasticbeanstalk.Client
	Name   *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"name"`
	// An application with live environments cannot be deleted without this, and using
	// it terminates them — so it is opt-in rather than a default.
	Force *bool `awsName:"TerminateEnvByForce" awsType:"awsbool" templateName:"force"`
}

func (cmd *DeleteApplication) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("force"),
	))
}

type UpdateApplication struct {
	_           string `action:"update" entity:"application" awsAPI:"elasticbeanstalk" awsCall:"UpdateApplication" awsInput:"elasticbeanstalk.UpdateApplicationInput" awsOutput:"elasticbeanstalk.UpdateApplicationOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *elasticbeanstalk.Client
	Name        *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"name"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

// Description is the only mutable field on an application, so it is required rather than
// optional: an update without it clears the description instead of doing nothing.
func (cmd *UpdateApplication) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("description")))
}

func (cmd *UpdateApplication) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type CreateEnvironment struct {
	_      string `action:"create" entity:"environment" awsAPI:"elasticbeanstalk" awsCall:"CreateEnvironment" awsInput:"elasticbeanstalk.CreateEnvironmentInput" awsOutput:"elasticbeanstalk.CreateEnvironmentOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *elasticbeanstalk.Client
	Name   *string `awsName:"EnvironmentName" awsType:"awsstr" templateName:"name"`
	App    *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"application"`
	// A platform is given either as a solution stack name or as a platform ARN; the two
	// are alternative spellings of the same thing and AWS rejects both together.
	SolutionStack *string `awsName:"SolutionStackName" awsType:"awsstr" templateName:"solution-stack"`
	Platform      *string `awsName:"PlatformArn" awsType:"awsstr" templateName:"platform"`
	Template      *string `awsName:"TemplateName" awsType:"awsstr" templateName:"config-template"`
	Version       *string `awsName:"VersionLabel" awsType:"awsstr" templateName:"version"`
	Description   *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	CNAMEPrefix   *string `awsName:"CNAMEPrefix" awsType:"awsstr" templateName:"cname-prefix"`
}

func (cmd *CreateEnvironment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("application"),
		params.OnlyOneOf(
			params.Key("solution-stack"),
			params.Key("platform"),
			params.Key("config-template"),
		),
		params.Opt(params.Suggested("version"), "description", "cname-prefix"),
	))
}

func (cmd *CreateEnvironment) ExtractResult(i any) string {
	out, ok := i.(*elasticbeanstalk.CreateEnvironmentOutput)
	if !ok || out.EnvironmentName == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.EnvironmentName)
}

// Terminating is Beanstalk's word for deleting an environment, and it is what the delete
// verb maps to here.
type DeleteEnvironment struct {
	_      string `action:"delete" entity:"environment" awsAPI:"elasticbeanstalk" awsCall:"TerminateEnvironment" awsInput:"elasticbeanstalk.TerminateEnvironmentInput" awsOutput:"elasticbeanstalk.TerminateEnvironmentOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *elasticbeanstalk.Client
	Name   *string `awsName:"EnvironmentName" awsType:"awsstr" templateName:"name"`
	// Keeping the resources behind means the load balancer and instances survive the
	// environment, which is rarely what is wanted but is occasionally deliberate.
	KeepResources *bool `awsName:"TerminateResources" awsType:"awsbool" templateName:"terminate-resources"`
	Force         *bool `awsName:"ForceTerminate" awsType:"awsbool" templateName:"force"`
}

func (cmd *DeleteEnvironment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("terminate-resources", "force"),
	))
}

type UpdateEnvironment struct {
	_           string `action:"update" entity:"environment" awsAPI:"elasticbeanstalk" awsCall:"UpdateEnvironment" awsInput:"elasticbeanstalk.UpdateEnvironmentInput" awsOutput:"elasticbeanstalk.UpdateEnvironmentOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *elasticbeanstalk.Client
	Name        *string `awsName:"EnvironmentName" awsType:"awsstr" templateName:"name"`
	Version     *string `awsName:"VersionLabel" awsType:"awsstr" templateName:"version"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	// Changing the platform in place is how a managed platform update is applied.
	SolutionStack *string `awsName:"SolutionStackName" awsType:"awsstr" templateName:"solution-stack"`
	Platform      *string `awsName:"PlatformArn" awsType:"awsstr" templateName:"platform"`
}

// Deploying a new version is the common case, so version is suggested; an update naming
// nothing to change would be a no-op that still triggers an environment update.
func (cmd *UpdateEnvironment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.AtLeastOneOf(
			params.Key("version"),
			params.Key("description"),
			params.Key("solution-stack"),
			params.Key("platform"),
		),
	))
}

func (cmd *UpdateEnvironment) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

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
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// The entity is deployapplication rather than application, which Elastic Beanstalk owns.
// Both are "applications" in their own service's language, and the template grammar has one
// flat namespace.
type CreateDeployapplication struct {
	_        string `action:"create" entity:"deployapplication" awsAPI:"codedeploy" awsCall:"CreateApplication" awsInput:"codedeploy.CreateApplicationInput" awsOutput:"codedeploy.CreateApplicationOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *codedeploy.Client
	Name     *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"name"`
	Platform *string `awsName:"ComputePlatform" awsType:"awsstr" templateName:"platform"`
}

func (cmd *CreateDeployapplication) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("platform")),
	))
}

// CreateApplication returns an id, but every other CodeDeploy call takes the name.
func (cmd *CreateDeployapplication) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteDeployapplication struct {
	_      string `action:"delete" entity:"deployapplication" awsAPI:"codedeploy" awsCall:"DeleteApplication" awsInput:"codedeploy.DeleteApplicationInput" awsOutput:"codedeploy.DeleteApplicationOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *codedeploy.Client
	Name   *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteDeployapplication) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type CreateDeploymentgroup struct {
	_           string `action:"create" entity:"deploymentgroup" awsAPI:"codedeploy" awsCall:"CreateDeploymentGroup" awsInput:"codedeploy.CreateDeploymentGroupInput" awsOutput:"codedeploy.CreateDeploymentGroupOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *codedeploy.Client
	Name        *string   `awsName:"DeploymentGroupName" awsType:"awsstr" templateName:"name"`
	Application *string   `awsName:"ApplicationName" awsType:"awsstr" templateName:"application"`
	Role        *string   `awsName:"ServiceRoleArn" awsType:"awsstr" templateName:"role"`
	Config      *string   `awsName:"DeploymentConfigName" awsType:"awsstr" templateName:"config"`
	ScalingGrps []*string `awsName:"AutoScalingGroups" awsType:"awsstringslice" templateName:"scalinggroups"`
	// The targeting rules — EC2 tag filters, on-premises filters, ECS services — are
	// nested documents, so they come from a file rather than being flattened.
	Ec2FilterFile *string `awsName:"Ec2TagFilters" awsType:"awsfiletostruct" templateName:"ec2-filters-file"`
	StyleFile     *string `awsName:"DeploymentStyle" awsType:"awsfiletostruct" templateName:"style-file"`
	RollbackFile  *string `awsName:"AutoRollbackConfiguration" awsType:"awsfiletostruct" templateName:"rollback-file"`
}

func (cmd *CreateDeploymentgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("application"), params.Key("role"),
		params.Opt(
			params.Suggested("config"),
			"scalinggroups", "ec2-filters-file", "style-file", "rollback-file",
		),
	))
}

func (cmd *CreateDeploymentgroup) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteDeploymentgroup struct {
	_           string `action:"delete" entity:"deploymentgroup" awsAPI:"codedeploy" awsCall:"DeleteDeploymentGroup" awsInput:"codedeploy.DeleteDeploymentGroupInput" awsOutput:"codedeploy.DeleteDeploymentGroupOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *codedeploy.Client
	Name        *string `awsName:"DeploymentGroupName" awsType:"awsstr" templateName:"name"`
	Application *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"application"`
}

// A deployment group name is only unique within its application, so both are required.
func (cmd *DeleteDeploymentgroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("application")))
}

// A deployment is an event rather than a resource: it is created, runs, and is done. It is
// not listed, because listing needs an application and group per call and what it returns is
// history.
type CreateDeployment struct {
	_           string `action:"create" entity:"deployment" awsAPI:"codedeploy" awsCall:"CreateDeployment" awsInput:"codedeploy.CreateDeploymentInput" awsOutput:"codedeploy.CreateDeploymentOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *codedeploy.Client
	Application *string `awsName:"ApplicationName" awsType:"awsstr" templateName:"application"`
	Group       *string `awsName:"DeploymentGroupName" awsType:"awsstr" templateName:"group"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Config      *string `awsName:"DeploymentConfigName" awsType:"awsstr" templateName:"config"`
	// The revision says what to deploy — an S3 object, a GitHub commit, an ECS task
	// definition — and its shape differs per source, so it is a document.
	RevisionFile *string `awsName:"Revision" awsType:"awsfiletostruct" templateName:"revision-file"`
}

func (cmd *CreateDeployment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("application"), params.Key("group"),
		params.Opt(params.Suggested("revision-file"), "description", "config"),
	))
}

func (cmd *CreateDeployment) ExtractResult(i any) string {
	out, ok := i.(*codedeploy.CreateDeploymentOutput)
	if !ok || out.DeploymentId == nil {
		return ""
	}
	return awssdk.ToString(out.DeploymentId)
}

// Stopping a deployment is how a bad rollout is halted, so it is worth having even though
// there is no delete.
type StopDeployment struct {
	_            string `action:"stop" entity:"deployment" awsAPI:"codedeploy" awsCall:"StopDeployment" awsInput:"codedeploy.StopDeploymentInput" awsOutput:"codedeploy.StopDeploymentOutput"`
	logger       *logger.Logger
	graph        cloud.GraphAPI
	api          *codedeploy.Client
	ID           *string `awsName:"DeploymentId" awsType:"awsstr" templateName:"id"`
	AutoRollback *bool   `awsName:"AutoRollbackEnabled" awsType:"awsbool" templateName:"rollback"`
}

func (cmd *StopDeployment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.Opt(params.Suggested("rollback")),
	))
}

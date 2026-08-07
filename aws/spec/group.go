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
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateGroup struct {
	_      string `action:"create" entity:"group" awsAPI:"iam" awsCall:"CreateGroup" awsInput:"iam.CreateGroupInput" awsOutput:"iam.CreateGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *iam.Client
	Name   *string `awsName:"GroupName" awsType:"awsstr" templateName:"name"`
}

func (cmd *CreateGroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

func (cmd *CreateGroup) ExtractResult(i any) string {
	return awssdk.ToString(i.(*iam.CreateGroupOutput).Group.GroupId)
}

type DeleteGroup struct {
	_      string `action:"delete" entity:"group" awsAPI:"iam" awsCall:"DeleteGroup" awsInput:"iam.DeleteGroupInput" awsOutput:"iam.DeleteGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *iam.Client
	Name   *string `awsName:"GroupName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteGroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

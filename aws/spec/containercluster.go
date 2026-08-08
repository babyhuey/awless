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
	"github.com/aws/aws-sdk-go-v2/service/ecs"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateContainercluster struct {
	_      string `action:"create" entity:"containercluster" awsAPI:"ecs" awsCall:"CreateCluster" awsInput:"ecs.CreateClusterInput" awsOutput:"ecs.CreateClusterOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ecs.Client
	Name   *string `awsName:"ClusterName" awsType:"awsstr" templateName:"name"`
}

func (cmd *CreateContainercluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

func (cmd *CreateContainercluster) ExtractResult(i any) string {
	return awssdk.ToString(i.(*ecs.CreateClusterOutput).Cluster.ClusterArn)
}

type DeleteContainercluster struct {
	_      string `action:"delete" entity:"containercluster" awsAPI:"ecs" awsCall:"DeleteCluster" awsInput:"ecs.DeleteClusterInput" awsOutput:"ecs.DeleteClusterOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ecs.Client
	ID     *string `awsName:"Cluster" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteContainercluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

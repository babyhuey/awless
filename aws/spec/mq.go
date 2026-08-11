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
	"github.com/aws/aws-sdk-go-v2/service/mq"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateBroker struct {
	_      string `action:"create" entity:"broker" awsAPI:"mq" awsCall:"CreateBroker" awsInput:"mq.CreateBrokerInput" awsOutput:"mq.CreateBrokerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *mq.Client
	Name   *string `awsName:"BrokerName" awsType:"awsstr" templateName:"name"`
	Engine *string `awsName:"EngineType" awsType:"awsstr" templateName:"engine"`
	Type   *string `awsName:"HostInstanceType" awsType:"awsstr" templateName:"type"`
	// SINGLE_INSTANCE or one of the active/standby and cluster modes. It decides both the
	// cost and whether the broker survives losing an availability zone, and the number of
	// subnets required depends on it, so there is no safe default.
	Mode           *string   `awsName:"DeploymentMode" awsType:"awsstr" templateName:"mode"`
	EngineVersion  *string   `awsName:"EngineVersion" awsType:"awsstr" templateName:"engine-version"`
	Subnets        []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
	Securitygroups []*string `awsName:"SecurityGroups" awsType:"awsstringslice" templateName:"securitygroups"`
	// Required by the API and security-relevant: a publicly accessible broker is
	// reachable from the internet.
	Public *bool `awsName:"PubliclyAccessible" awsType:"awsbool" templateName:"public"`
	// The initial user's credentials. A document because a broker takes a list of users
	// with per-user console access and group membership, and because putting a password on
	// the command line would also put it in the template log.
	UsersFile *string `awsName:"Users" awsType:"awsfiletostruct" templateName:"users-file"`
}

func (cmd *CreateBroker) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("engine"), params.Key("type"),
		params.Key("mode"), params.Key("public"), params.Key("users-file"),
		params.Opt(
			params.Suggested("engine-version", "subnets"),
			"securitygroups",
		),
	))
}

func (cmd *CreateBroker) ExtractResult(i any) string {
	out, ok := i.(*mq.CreateBrokerOutput)
	if !ok || out.BrokerId == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.BrokerId)
}

type DeleteBroker struct {
	_      string `action:"delete" entity:"broker" awsAPI:"mq" awsCall:"DeleteBroker" awsInput:"mq.DeleteBrokerInput" awsOutput:"mq.DeleteBrokerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *mq.Client
	// The broker id, not the name: broker names are not unique across time and the API
	// takes the id.
	ID *string `awsName:"BrokerId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteBroker) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

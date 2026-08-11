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
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// A peering connection is requested by one VPC and accepted by the other, so creating one
// leaves it pending-acceptance rather than active. `start vpcpeering` is the accept step,
// which is a separate call and often a separate account.
type CreateVpcpeering struct {
	_       string `action:"create" entity:"vpcpeering" awsAPI:"ec2" awsCall:"CreateVpcPeeringConnection" awsInput:"ec2.CreateVpcPeeringConnectionInput" awsOutput:"ec2.CreateVpcPeeringConnectionOutput" awsDryRun:""`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *ec2.Client
	Vpc     *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
	PeerVpc *string `awsName:"PeerVpcId" awsType:"awsstr" templateName:"peer-vpc"`
	// Needed when the other VPC is in a different account or region; the request goes to
	// that account for acceptance.
	PeerOwner  *string `awsName:"PeerOwnerId" awsType:"awsstr" templateName:"peer-owner"`
	PeerRegion *string `awsName:"PeerRegion" awsType:"awsstr" templateName:"peer-region"`
}

func (cmd *CreateVpcpeering) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("vpc"), params.Key("peer-vpc"),
		params.Opt("peer-owner", "peer-region"),
	))
}

func (cmd *CreateVpcpeering) ExtractResult(i any) string {
	out, ok := i.(*ec2.CreateVpcPeeringConnectionOutput)
	if !ok || out.VpcPeeringConnection == nil {
		return ""
	}
	return awssdk.ToString(out.VpcPeeringConnection.VpcPeeringConnectionId)
}

type DeleteVpcpeering struct {
	_      string `action:"delete" entity:"vpcpeering" awsAPI:"ec2" awsCall:"DeleteVpcPeeringConnection" awsInput:"ec2.DeleteVpcPeeringConnectionInput" awsOutput:"ec2.DeleteVpcPeeringConnectionOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"VpcPeeringConnectionId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteVpcpeering) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

// Accepting is spelled `start` rather than `update` because it is a one-way transition with
// no parameters, not a modification.
type StartVpcpeering struct {
	_      string `action:"start" entity:"vpcpeering" awsAPI:"ec2" awsCall:"AcceptVpcPeeringConnection" awsInput:"ec2.AcceptVpcPeeringConnectionInput" awsOutput:"ec2.AcceptVpcPeeringConnectionOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"VpcPeeringConnectionId" awsType:"awsstr" templateName:"id"`
}

func (cmd *StartVpcpeering) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

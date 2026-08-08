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
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateSubnet struct {
	_                string `action:"create" entity:"subnet" awsAPI:"ec2" awsCall:"CreateSubnet" awsInput:"ec2.CreateSubnetInput" awsOutput:"ec2.CreateSubnetOutput" awsDryRun:""`
	logger           *logger.Logger
	graph            cloud.GraphAPI
	api              *ec2.Client
	CIDR             *string `awsName:"CidrBlock" awsType:"awsstr" templateName:"cidr"`
	VPC              *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
	AvailabilityZone *string `awsName:"AvailabilityZone" awsType:"awsstr" templateName:"availabilityzone"`
	Public           *bool   `awsType:"awsboolattribute" templateName:"public"`
	Name             *string `templateName:"name"`
}

func (cmd *CreateSubnet) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("cidr"), params.Key("vpc"), params.Opt(params.Suggested("name"), "availabilityzone", "public")),
		params.Validators{"cidr": params.IsCIDR})
}

func (cmd *CreateSubnet) ExtractResult(i any) string {
	return awssdk.ToString(i.(*ec2.CreateSubnetOutput).Subnet.SubnetId)
}

func (cmd *CreateSubnet) AfterRun(renv env.Running, output any) error {
	subnetID := awssdk.String(cmd.ExtractResult(output))
	if err := createNameTag(subnetID, cmd.Name, renv); err != nil {
		return err
	}

	if BoolValue(cmd.Public) {
		updateSubnet := CommandFactory.Build("updatesubnet")().(*UpdateSubnet)
		updateSubnet.ID = subnetID
		updateSubnet.Public = Bool(true)
		if _, err := updateSubnet.Run(renv, nil); err != nil {
			return err
		}
	}

	return nil
}

type UpdateSubnet struct {
	_      string `action:"update" entity:"subnet" awsAPI:"ec2" awsCall:"ModifySubnetAttribute" awsInput:"ec2.ModifySubnetAttributeInput" awsOutput:"ec2.ModifySubnetAttributeOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"SubnetId" awsType:"awsstr" templateName:"id"`
	Public *bool   `awsName:"MapPublicIpOnLaunch" awsType:"awsboolattribute" templateName:"public"`
}

func (cmd *UpdateSubnet) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Opt("public")))
}

type DeleteSubnet struct {
	_      string `action:"delete" entity:"subnet" awsAPI:"ec2" awsCall:"DeleteSubnet" awsInput:"ec2.DeleteSubnetInput" awsOutput:"ec2.DeleteSubnetOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"SubnetId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteSubnet) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

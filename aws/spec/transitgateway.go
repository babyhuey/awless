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

// Transit gateways and VPC endpoints are on the EC2 API, so they share the client and the
// dry-run support the rest of the EC2 commands have.

type CreateTransitgateway struct {
	_           string `action:"create" entity:"transitgateway" awsAPI:"ec2" awsCall:"CreateTransitGateway" awsInput:"ec2.CreateTransitGatewayInput" awsOutput:"ec2.CreateTransitGatewayOutput" awsDryRun:""`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *ec2.Client
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	// The ASN is fixed at creation. Leaving it unset takes the AWS default of 64512.
	ASN *int64 `awsName:"Options.AmazonSideAsn" awsType:"awsint64" templateName:"asn"`
	// Automatic acceptance and automatic association are what make a transit gateway
	// convenient in a single account and dangerous across many, so they are explicit.
	AutoAccept      *string `awsName:"Options.AutoAcceptSharedAttachments" awsType:"awsstr" templateName:"auto-accept"`
	DefaultAssoc    *string `awsName:"Options.DefaultRouteTableAssociation" awsType:"awsstr" templateName:"default-association"`
	DefaultPropagat *string `awsName:"Options.DefaultRouteTablePropagation" awsType:"awsstr" templateName:"default-propagation"`
	DNSSupport      *string `awsName:"Options.DnsSupport" awsType:"awsstr" templateName:"dns-support"`
}

func (cmd *CreateTransitgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Opt(
			params.Suggested("description"),
			"asn", "auto-accept", "default-association", "default-propagation", "dns-support",
		),
	))
}

func (cmd *CreateTransitgateway) ExtractResult(i any) string {
	out, ok := i.(*ec2.CreateTransitGatewayOutput)
	if !ok || out.TransitGateway == nil {
		return ""
	}
	return awssdk.ToString(out.TransitGateway.TransitGatewayId)
}

type DeleteTransitgateway struct {
	_      string `action:"delete" entity:"transitgateway" awsAPI:"ec2" awsCall:"DeleteTransitGateway" awsInput:"ec2.DeleteTransitGatewayInput" awsOutput:"ec2.DeleteTransitGatewayOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"TransitGatewayId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteTransitgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

// An attachment joins a VPC to a transit gateway. It is a create rather than an attach
// because it is a resource with its own id and lifecycle, not a link between two existing
// things.
type CreateTransitgatewayattachment struct {
	_       string `action:"create" entity:"transitgatewayattachment" awsAPI:"ec2" awsCall:"CreateTransitGatewayVpcAttachment" awsInput:"ec2.CreateTransitGatewayVpcAttachmentInput" awsOutput:"ec2.CreateTransitGatewayVpcAttachmentOutput" awsDryRun:""`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *ec2.Client
	Gateway *string   `awsName:"TransitGatewayId" awsType:"awsstr" templateName:"transitgateway"`
	Vpc     *string   `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
	Subnets []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
}

// One subnet per availability zone is what the gateway uses to reach the VPC, so subnets
// are required rather than defaulted.
func (cmd *CreateTransitgatewayattachment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("transitgateway"), params.Key("vpc"), params.Key("subnets"),
	))
}

func (cmd *CreateTransitgatewayattachment) ExtractResult(i any) string {
	out, ok := i.(*ec2.CreateTransitGatewayVpcAttachmentOutput)
	if !ok || out.TransitGatewayVpcAttachment == nil {
		return ""
	}
	return awssdk.ToString(out.TransitGatewayVpcAttachment.TransitGatewayAttachmentId)
}

type DeleteTransitgatewayattachment struct {
	_      string `action:"delete" entity:"transitgatewayattachment" awsAPI:"ec2" awsCall:"DeleteTransitGatewayVpcAttachment" awsInput:"ec2.DeleteTransitGatewayVpcAttachmentInput" awsOutput:"ec2.DeleteTransitGatewayVpcAttachmentOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"TransitGatewayAttachmentId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteTransitgatewayattachment) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CreateTransitgatewayroutetable struct {
	_       string `action:"create" entity:"transitgatewayroutetable" awsAPI:"ec2" awsCall:"CreateTransitGatewayRouteTable" awsInput:"ec2.CreateTransitGatewayRouteTableInput" awsOutput:"ec2.CreateTransitGatewayRouteTableOutput" awsDryRun:""`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *ec2.Client
	Gateway *string `awsName:"TransitGatewayId" awsType:"awsstr" templateName:"transitgateway"`
}

func (cmd *CreateTransitgatewayroutetable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("transitgateway")))
}

func (cmd *CreateTransitgatewayroutetable) ExtractResult(i any) string {
	out, ok := i.(*ec2.CreateTransitGatewayRouteTableOutput)
	if !ok || out.TransitGatewayRouteTable == nil {
		return ""
	}
	return awssdk.ToString(out.TransitGatewayRouteTable.TransitGatewayRouteTableId)
}

type DeleteTransitgatewayroutetable struct {
	_      string `action:"delete" entity:"transitgatewayroutetable" awsAPI:"ec2" awsCall:"DeleteTransitGatewayRouteTable" awsInput:"ec2.DeleteTransitGatewayRouteTableInput" awsOutput:"ec2.DeleteTransitGatewayRouteTableOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	ID     *string `awsName:"TransitGatewayRouteTableId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteTransitgatewayroutetable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CreateVpcendpoint struct {
	_       string `action:"create" entity:"vpcendpoint" awsAPI:"ec2" awsCall:"CreateVpcEndpoint" awsInput:"ec2.CreateVpcEndpointInput" awsOutput:"ec2.CreateVpcEndpointOutput" awsDryRun:""`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *ec2.Client
	Vpc     *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
	Service *string `awsName:"ServiceName" awsType:"awsstr" templateName:"service"`
	// Gateway endpoints attach to route tables and cost nothing; Interface endpoints
	// place an ENI in each subnet and are billed hourly. The distinction decides which
	// of the two sets of params below applies, so the type is required rather than
	// defaulted to the cheaper one by accident.
	Type           *string   `awsName:"VpcEndpointType" awsType:"awsstr" templateName:"type"`
	Subnets        []*string `awsName:"SubnetIds" awsType:"awsstringslice" templateName:"subnets"`
	Securitygroups []*string `awsName:"SecurityGroupIds" awsType:"awsstringslice" templateName:"securitygroups"`
	RouteTables    []*string `awsName:"RouteTableIds" awsType:"awsstringslice" templateName:"routetables"`
	PrivateDNS     *bool     `awsName:"PrivateDnsEnabled" awsType:"awsbool" templateName:"private-dns"`
	PolicyFile     *string   `awsName:"PolicyDocument" awsType:"awsfiletostring" templateName:"policy-file"`
}

func (cmd *CreateVpcendpoint) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("vpc"), params.Key("service"), params.Key("type"),
		params.Opt(
			params.Suggested("subnets", "securitygroups"),
			"routetables", "private-dns", "policy-file",
		),
	))
}

func (cmd *CreateVpcendpoint) ExtractResult(i any) string {
	out, ok := i.(*ec2.CreateVpcEndpointOutput)
	if !ok || out.VpcEndpoint == nil {
		return ""
	}
	return awssdk.ToString(out.VpcEndpoint.VpcEndpointId)
}

type DeleteVpcendpoint struct {
	_      string `action:"delete" entity:"vpcendpoint" awsAPI:"ec2" awsCall:"DeleteVpcEndpoints" awsInput:"ec2.DeleteVpcEndpointsInput" awsOutput:"ec2.DeleteVpcEndpointsOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ec2.Client
	// The API deletes in bulk, hence the plural field, but one id at a time is what a
	// revert produces and what a user means.
	ID []*string `awsName:"VpcEndpointIds" awsType:"awsstringslice" templateName:"id"`
}

func (cmd *DeleteVpcendpoint) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

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
	"github.com/aws/aws-sdk-go-v2/service/directconnect"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// Direct Connect Gateway — the only resource in this service worth creating via CLI.
// Connections and LAGs are ordered through the AWS console or partner portal.
// VIFs require nested structs that map poorly to the template DSL.

type CreateDirectconnectgateway struct {
	_      string `action:"create" entity:"directconnectgateway" awsAPI:"directconnect" awsCall:"CreateDirectConnectGateway" awsInput:"directconnect.CreateDirectConnectGatewayInput" awsOutput:"directconnect.CreateDirectConnectGatewayOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *directconnect.Client
	Name   *string `awsName:"DirectConnectGatewayName" awsType:"awsstr" templateName:"name"`
	Asn    *int64  `awsName:"AmazonSideAsn" awsType:"awsint64" templateName:"amazon-side-asn"`
}

func (cmd *CreateDirectconnectgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Opt("amazon-side-asn")))
}

func (cmd *CreateDirectconnectgateway) ExtractResult(i any) string {
	out, ok := i.(*directconnect.CreateDirectConnectGatewayOutput)
	if !ok || out.DirectConnectGateway == nil {
		return ""
	}
	return awssdk.ToString(out.DirectConnectGateway.DirectConnectGatewayId)
}

type DeleteDirectconnectgateway struct {
	_      string `action:"delete" entity:"directconnectgateway" awsAPI:"directconnect" awsCall:"DeleteDirectConnectGateway" awsInput:"directconnect.DeleteDirectConnectGatewayInput" awsOutput:"directconnect.DeleteDirectConnectGatewayOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *directconnect.Client
	ID     *string `awsName:"DirectConnectGatewayId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteDirectconnectgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

// Direct Connect Gateway Association — associates a DCGW with a VGW or TGW.

type CreateDirectconnectgatewayassociation struct {
	_              string `action:"create" entity:"directconnectgatewayassociation" awsAPI:"directconnect" awsCall:"CreateDirectConnectGatewayAssociation" awsInput:"directconnect.CreateDirectConnectGatewayAssociationInput" awsOutput:"directconnect.CreateDirectConnectGatewayAssociationOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *directconnect.Client
	Gateway        *string `awsName:"DirectConnectGatewayId" awsType:"awsstr" templateName:"gateway"`
	AssocGateway   *string `awsName:"GatewayId" awsType:"awsstr" templateName:"associated-gateway"`
	VirtualGateway *string `awsName:"VirtualGatewayId" awsType:"awsstr" templateName:"virtual-gateway"`
}

func (cmd *CreateDirectconnectgatewayassociation) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("gateway"),
		params.OnlyOneOf(params.Key("associated-gateway"), params.Key("virtual-gateway")),
	))
}

func (cmd *CreateDirectconnectgatewayassociation) ExtractResult(i any) string {
	out, ok := i.(*directconnect.CreateDirectConnectGatewayAssociationOutput)
	if !ok || out.DirectConnectGatewayAssociation == nil {
		return ""
	}
	return awssdk.ToString(out.DirectConnectGatewayAssociation.AssociationId)
}

type DeleteDirectconnectgatewayassociation struct {
	_      string `action:"delete" entity:"directconnectgatewayassociation" awsAPI:"directconnect" awsCall:"DeleteDirectConnectGatewayAssociation" awsInput:"directconnect.DeleteDirectConnectGatewayAssociationInput" awsOutput:"directconnect.DeleteDirectConnectGatewayAssociationOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *directconnect.Client
	ID     *string `awsName:"AssociationId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteDirectconnectgatewayassociation) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

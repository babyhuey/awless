package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestCreateTransitgateway(t *testing.T) {
	mock := NewMock().On("CreateTransitGateway", &ec2.CreateTransitGatewayOutput{
		TransitGateway: &ec2types.TransitGateway{TransitGatewayId: awssdk.String("tgw-1234")},
	})

	Template("create transitgateway description=hub asn=64512 auto-accept=disable default-association=disable").
		Mock(mock).
		ExpectCalls("CreateTransitGateway").
		ExpectCommandResult("tgw-1234").
		ExpectRevert("delete transitgateway id=tgw-1234").
		Run(t)

	in := mock.InputFor("CreateTransitGateway").(*ec2.CreateTransitGatewayInput)
	if got := awssdk.ToString(in.Description); got != "hub" {
		t.Errorf("Description: got %q, want hub", got)
	}
	// Options is a nested struct, so this checks the dotted paths built it.
	if in.Options == nil {
		t.Fatal("Options was never built")
	}
	if got := awssdk.ToInt64(in.Options.AmazonSideAsn); got != 64512 {
		t.Errorf("Options.AmazonSideAsn: got %d, want 64512", got)
	}
	if got := string(in.Options.AutoAcceptSharedAttachments); got != "disable" {
		t.Errorf("Options.AutoAcceptSharedAttachments: got %q, want disable", got)
	}
	if got := string(in.Options.DefaultRouteTableAssociation); got != "disable" {
		t.Errorf("Options.DefaultRouteTableAssociation: got %q, want disable", got)
	}
}

// Every param is optional, so a bare create must still work — a transit gateway with all
// AWS defaults is a legitimate thing to ask for.
func TestCreateTransitgatewayWithNoParams(t *testing.T) {
	mock := NewMock().On("CreateTransitGateway", &ec2.CreateTransitGatewayOutput{
		TransitGateway: &ec2types.TransitGateway{TransitGatewayId: awssdk.String("tgw-1234")},
	})

	Template("create transitgateway").
		Mock(mock).
		ExpectCalls("CreateTransitGateway").
		ExpectCommandResult("tgw-1234").
		Run(t)
}

func TestDeleteTransitgateway(t *testing.T) {
	mock := NewMock().On("DeleteTransitGateway", &ec2.DeleteTransitGatewayOutput{})

	Template("delete transitgateway id=tgw-1234").
		Mock(mock).
		ExpectCalls("DeleteTransitGateway").
		Run(t)

	in := mock.InputFor("DeleteTransitGateway").(*ec2.DeleteTransitGatewayInput)
	if got := awssdk.ToString(in.TransitGatewayId); got != "tgw-1234" {
		t.Errorf("TransitGatewayId: got %q, want tgw-1234", got)
	}
}

func TestCreateTransitgatewayattachment(t *testing.T) {
	mock := NewMock().On("CreateTransitGatewayVpcAttachment", &ec2.CreateTransitGatewayVpcAttachmentOutput{
		TransitGatewayVpcAttachment: &ec2types.TransitGatewayVpcAttachment{
			TransitGatewayAttachmentId: awssdk.String("tgw-attach-1234"),
		},
	})

	Template("create transitgatewayattachment transitgateway=tgw-1234 vpc=vpc-5678 subnets=subnet-1,subnet-2").
		Mock(mock).
		ExpectCalls("CreateTransitGatewayVpcAttachment").
		ExpectCommandResult("tgw-attach-1234").
		ExpectRevert("delete transitgatewayattachment id=tgw-attach-1234").
		Run(t)

	in := mock.InputFor("CreateTransitGatewayVpcAttachment").(*ec2.CreateTransitGatewayVpcAttachmentInput)
	if got := awssdk.ToString(in.TransitGatewayId); got != "tgw-1234" {
		t.Errorf("TransitGatewayId: got %q", got)
	}
	if got := awssdk.ToString(in.VpcId); got != "vpc-5678" {
		t.Errorf("VpcId: got %q, want vpc-5678", got)
	}
	if len(in.SubnetIds) != 2 || in.SubnetIds[1] != "subnet-2" {
		t.Errorf("SubnetIds: got %v", in.SubnetIds)
	}
}

func TestCreateTransitgatewayroutetable(t *testing.T) {
	mock := NewMock().On("CreateTransitGatewayRouteTable", &ec2.CreateTransitGatewayRouteTableOutput{
		TransitGatewayRouteTable: &ec2types.TransitGatewayRouteTable{
			TransitGatewayRouteTableId: awssdk.String("tgw-rtb-1234"),
		},
	})

	Template("create transitgatewayroutetable transitgateway=tgw-1234").
		Mock(mock).
		ExpectCalls("CreateTransitGatewayRouteTable").
		ExpectCommandResult("tgw-rtb-1234").
		ExpectRevert("delete transitgatewayroutetable id=tgw-rtb-1234").
		Run(t)
}

// A Gateway endpoint attaches to route tables; an Interface endpoint places network
// interfaces in subnets. Both shapes go through the same command.
func TestCreateVpcendpointGateway(t *testing.T) {
	mock := NewMock().On("CreateVpcEndpoint", &ec2.CreateVpcEndpointOutput{
		VpcEndpoint: &ec2types.VpcEndpoint{VpcEndpointId: awssdk.String("vpce-1234")},
	})

	Template("create vpcendpoint vpc=vpc-1234 service=com.amazonaws.us-west-2.s3 type=Gateway routetables=rtb-1234,rtb-5678").
		Mock(mock).
		ExpectCalls("CreateVpcEndpoint").
		ExpectCommandResult("vpce-1234").
		ExpectRevert("delete vpcendpoint id=vpce-1234").
		Run(t)

	in := mock.InputFor("CreateVpcEndpoint").(*ec2.CreateVpcEndpointInput)
	if got := string(in.VpcEndpointType); got != "Gateway" {
		t.Errorf("VpcEndpointType: got %q, want Gateway", got)
	}
	if got := awssdk.ToString(in.ServiceName); got != "com.amazonaws.us-west-2.s3" {
		t.Errorf("ServiceName: got %q", got)
	}
	if len(in.RouteTableIds) != 2 {
		t.Errorf("RouteTableIds: got %v", in.RouteTableIds)
	}
	// A gateway endpoint has no interfaces, so nothing subnet-related should be sent.
	if len(in.SubnetIds) != 0 {
		t.Errorf("SubnetIds should be empty for a Gateway endpoint, got %v", in.SubnetIds)
	}
}

func TestCreateVpcendpointInterface(t *testing.T) {
	mock := NewMock().On("CreateVpcEndpoint", &ec2.CreateVpcEndpointOutput{
		VpcEndpoint: &ec2types.VpcEndpoint{VpcEndpointId: awssdk.String("vpce-5678")},
	})

	Template("create vpcendpoint vpc=vpc-1234 service=com.amazonaws.us-west-2.secretsmanager " +
		"type=Interface subnets=subnet-1,subnet-2 securitygroups=sg-1234 private-dns=true").
		Mock(mock).
		ExpectCalls("CreateVpcEndpoint").
		ExpectCommandResult("vpce-5678").
		Run(t)

	in := mock.InputFor("CreateVpcEndpoint").(*ec2.CreateVpcEndpointInput)
	if got := string(in.VpcEndpointType); got != "Interface" {
		t.Errorf("VpcEndpointType: got %q, want Interface", got)
	}
	if len(in.SubnetIds) != 2 {
		t.Errorf("SubnetIds: got %v", in.SubnetIds)
	}
	if len(in.SecurityGroupIds) != 1 || in.SecurityGroupIds[0] != "sg-1234" {
		t.Errorf("SecurityGroupIds: got %v", in.SecurityGroupIds)
	}
	if !awssdk.ToBool(in.PrivateDnsEnabled) {
		t.Error("expected PrivateDnsEnabled to be true")
	}
}

// The API deletes in bulk, so the single id has to arrive as a one-element slice.
func TestDeleteVpcendpoint(t *testing.T) {
	mock := NewMock().On("DeleteVpcEndpoints", &ec2.DeleteVpcEndpointsOutput{})

	Template("delete vpcendpoint id=vpce-1234").
		Mock(mock).
		ExpectCalls("DeleteVpcEndpoints").
		Run(t)

	in := mock.InputFor("DeleteVpcEndpoints").(*ec2.DeleteVpcEndpointsInput)
	if len(in.VpcEndpointIds) != 1 || in.VpcEndpointIds[0] != "vpce-1234" {
		t.Errorf("VpcEndpointIds: got %v, want [vpce-1234]", in.VpcEndpointIds)
	}
}

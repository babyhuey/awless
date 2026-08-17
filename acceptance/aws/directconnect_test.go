package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	directconnecttypes "github.com/aws/aws-sdk-go-v2/service/directconnect/types"
)

func TestCreateDirectconnectgateway(t *testing.T) {
	mock := NewMock().On("CreateDirectConnectGateway", &directconnect.CreateDirectConnectGatewayOutput{
		DirectConnectGateway: &directconnecttypes.DirectConnectGateway{
			DirectConnectGatewayId: awssdk.String("abcd1234-ab12-cd34-ef56-abcdef123456"),
		},
	})

	Template("create directconnectgateway name=my-dcgw").
		Mock(mock).
		ExpectCalls("CreateDirectConnectGateway").
		ExpectCommandResult("abcd1234-ab12-cd34-ef56-abcdef123456").
		ExpectRevert("delete directconnectgateway id=abcd1234-ab12-cd34-ef56-abcdef123456").
		Run(t)

	in := mock.InputFor("CreateDirectConnectGateway").(*directconnect.CreateDirectConnectGatewayInput)
	if got := awssdk.ToString(in.DirectConnectGatewayName); got != "my-dcgw" {
		t.Errorf("DirectConnectGatewayName: got %q, want my-dcgw", got)
	}
	if in.AmazonSideAsn != nil {
		t.Errorf("AmazonSideAsn: should be nil when not specified, got %v", *in.AmazonSideAsn)
	}
}

func TestCreateDirectconnectgatewayWithAsn(t *testing.T) {
	mock := NewMock().On("CreateDirectConnectGateway", &directconnect.CreateDirectConnectGatewayOutput{
		DirectConnectGateway: &directconnecttypes.DirectConnectGateway{
			DirectConnectGatewayId: awssdk.String("abcd1234"),
		},
	})

	Template("create directconnectgateway name=my-dcgw amazon-side-asn=65000").
		Mock(mock).
		ExpectCalls("CreateDirectConnectGateway").
		Run(t)

	in := mock.InputFor("CreateDirectConnectGateway").(*directconnect.CreateDirectConnectGatewayInput)
	if got := awssdk.ToString(in.DirectConnectGatewayName); got != "my-dcgw" {
		t.Errorf("DirectConnectGatewayName: got %q", got)
	}
	if in.AmazonSideAsn == nil || *in.AmazonSideAsn != 65000 {
		t.Errorf("AmazonSideAsn: want 65000, got %v", in.AmazonSideAsn)
	}
}

func TestDeleteDirectconnectgateway(t *testing.T) {
	mock := NewMock().On("DeleteDirectConnectGateway", &directconnect.DeleteDirectConnectGatewayOutput{})

	Template("delete directconnectgateway id=abcd1234-ab12-cd34-ef56-abcdef123456").
		Mock(mock).
		ExpectCalls("DeleteDirectConnectGateway").
		Run(t)

	in := mock.InputFor("DeleteDirectConnectGateway").(*directconnect.DeleteDirectConnectGatewayInput)
	if got := awssdk.ToString(in.DirectConnectGatewayId); got != "abcd1234-ab12-cd34-ef56-abcdef123456" {
		t.Errorf("DirectConnectGatewayId: got %q", got)
	}
}

func TestCreateDirectconnectgatewayassociation(t *testing.T) {
	mock := NewMock().On("CreateDirectConnectGatewayAssociation", &directconnect.CreateDirectConnectGatewayAssociationOutput{
		DirectConnectGatewayAssociation: &directconnecttypes.DirectConnectGatewayAssociation{
			AssociationId: awssdk.String("assoc-1234"),
		},
	})

	Template("create directconnectgatewayassociation gateway=dcgw-1234 associated-gateway=tgw-5678").
		Mock(mock).
		ExpectCalls("CreateDirectConnectGatewayAssociation").
		ExpectCommandResult("assoc-1234").
		ExpectRevert("delete directconnectgatewayassociation id=assoc-1234").
		Run(t)

	in := mock.InputFor("CreateDirectConnectGatewayAssociation").(*directconnect.CreateDirectConnectGatewayAssociationInput)
	if got := awssdk.ToString(in.DirectConnectGatewayId); got != "dcgw-1234" {
		t.Errorf("DirectConnectGatewayId: got %q", got)
	}
	if got := awssdk.ToString(in.GatewayId); got != "tgw-5678" {
		t.Errorf("GatewayId: got %q", got)
	}
}

func TestDeleteDirectconnectgatewayassociation(t *testing.T) {
	mock := NewMock().On("DeleteDirectConnectGatewayAssociation", &directconnect.DeleteDirectConnectGatewayAssociationOutput{})

	Template("delete directconnectgatewayassociation id=assoc-1234").
		Mock(mock).
		ExpectCalls("DeleteDirectConnectGatewayAssociation").
		Run(t)

	in := mock.InputFor("DeleteDirectConnectGatewayAssociation").(*directconnect.DeleteDirectConnectGatewayAssociationInput)
	if got := awssdk.ToString(in.AssociationId); got != "assoc-1234" {
		t.Errorf("AssociationId: got %q", got)
	}
}

// Empty-response cases — an AWS reply with a nil body must not panic.

func TestCreateDirectconnectgatewayEmptyResponse(t *testing.T) {
	mock := NewMock().On("CreateDirectConnectGateway", &directconnect.CreateDirectConnectGatewayOutput{})

	Template("create directconnectgateway name=my-dcgw").
		Mock(mock).
		ExpectCalls("CreateDirectConnectGateway").
		ExpectCommandResult("").
		Run(t)
}

func TestCreateDirectconnectgatewayassociationEmptyResponse(t *testing.T) {
	mock := NewMock().On("CreateDirectConnectGatewayAssociation", &directconnect.CreateDirectConnectGatewayAssociationOutput{})

	Template("create directconnectgatewayassociation gateway=dcgw-1234 associated-gateway=tgw-5678").
		Mock(mock).
		ExpectCalls("CreateDirectConnectGatewayAssociation").
		ExpectCommandResult("").
		Run(t)
}

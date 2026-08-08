package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2 networking. Several of these take an association id rather than a resource id,
// which is easy to conflate, so the tests assert which one reaches the API.

func TestCreateInternetgateway(t *testing.T) {
	mock := NewMock().On("CreateInternetGateway", &ec2.CreateInternetGatewayOutput{
		InternetGateway: &ec2types.InternetGateway{
			InternetGatewayId: awssdk.String("igw-1234"),
		},
	})

	Template("create internetgateway").
		Mock(mock).
		ExpectCalls("CreateInternetGateway").
		ExpectCommandResult("igw-1234").
		ExpectRevert("delete internetgateway id=igw-1234").
		Run(t)
}

func TestAttachAndDetachInternetgateway(t *testing.T) {
	attach := NewMock().On("AttachInternetGateway", &ec2.AttachInternetGatewayOutput{})
	Template("attach internetgateway id=igw-1234 vpc=vpc-1234").
		Mock(attach).
		ExpectCalls("AttachInternetGateway").
		ExpectRevert("detach internetgateway id=igw-1234 vpc=vpc-1234").
		Run(t)

	in := attach.InputFor("AttachInternetGateway").(*ec2.AttachInternetGatewayInput)
	if got := awssdk.ToString(in.VpcId); got != "vpc-1234" {
		t.Errorf("VpcId: got %q, want vpc-1234", got)
	}

	detach := NewMock().On("DetachInternetGateway", &ec2.DetachInternetGatewayOutput{})
	Template("detach internetgateway id=igw-1234 vpc=vpc-1234").
		Mock(detach).
		ExpectCalls("DetachInternetGateway").
		Run(t)
}

func TestDeleteInternetgateway(t *testing.T) {
	mock := NewMock().On("DeleteInternetGateway", &ec2.DeleteInternetGatewayOutput{})

	Template("delete internetgateway id=igw-1234").
		Mock(mock).
		ExpectCalls("DeleteInternetGateway").
		Run(t)
}

func TestCreateRoutetable(t *testing.T) {
	mock := NewMock().On("CreateRouteTable", &ec2.CreateRouteTableOutput{
		RouteTable: &ec2types.RouteTable{RouteTableId: awssdk.String("rtb-1234")},
	})

	Template("create routetable vpc=vpc-1234").
		Mock(mock).
		ExpectCalls("CreateRouteTable").
		ExpectCommandResult("rtb-1234").
		ExpectRevert("delete routetable id=rtb-1234").
		Run(t)
}

// attach routetable returns an association id, and detach takes that association
// rather than the table id — so the revert has to carry the association through.
func TestAttachRoutetableRevertsWithAssociation(t *testing.T) {
	mock := NewMock().On("AssociateRouteTable", &ec2.AssociateRouteTableOutput{
		AssociationId: awssdk.String("rtbassoc-1234"),
	})

	Template("attach routetable id=rtb-1234 subnet=subnet-1234").
		Mock(mock).
		ExpectCalls("AssociateRouteTable").
		ExpectCommandResult("rtbassoc-1234").
		Run(t)

	in := mock.InputFor("AssociateRouteTable").(*ec2.AssociateRouteTableInput)
	if got := awssdk.ToString(in.SubnetId); got != "subnet-1234" {
		t.Errorf("SubnetId: got %q, want subnet-1234", got)
	}
}

func TestDetachRoutetableTakesAssociation(t *testing.T) {
	mock := NewMock().On("DisassociateRouteTable", &ec2.DisassociateRouteTableOutput{})

	Template("detach routetable association=rtbassoc-1234").
		Mock(mock).
		ExpectCalls("DisassociateRouteTable").
		Run(t)

	in := mock.InputFor("DisassociateRouteTable").(*ec2.DisassociateRouteTableInput)
	if got := awssdk.ToString(in.AssociationId); got != "rtbassoc-1234" {
		t.Errorf("AssociationId: got %q, want rtbassoc-1234", got)
	}
}

func TestDeleteRoutetable(t *testing.T) {
	mock := NewMock().On("DeleteRouteTable", &ec2.DeleteRouteTableOutput{})

	Template("delete routetable id=rtb-1234").
		Mock(mock).
		ExpectCalls("DeleteRouteTable").
		Run(t)
}

func TestCreateAndDeleteRoute(t *testing.T) {
	create := NewMock().On("CreateRoute", &ec2.CreateRouteOutput{})
	Template("create route cidr=0.0.0.0/0 gateway=igw-1234 table=rtb-1234").
		Mock(create).
		ExpectCalls("CreateRoute").
		ExpectRevert("delete route cidr=0.0.0.0/0 table=rtb-1234").
		Run(t)

	in := create.InputFor("CreateRoute").(*ec2.CreateRouteInput)
	if got := awssdk.ToString(in.DestinationCidrBlock); got != "0.0.0.0/0" {
		t.Errorf("DestinationCidrBlock: got %q, want 0.0.0.0/0", got)
	}
	if got := awssdk.ToString(in.GatewayId); got != "igw-1234" {
		t.Errorf("GatewayId: got %q, want igw-1234", got)
	}

	del := NewMock().On("DeleteRoute", &ec2.DeleteRouteOutput{})
	Template("delete route cidr=0.0.0.0/0 table=rtb-1234").
		Mock(del).
		ExpectCalls("DeleteRoute").
		Run(t)
}

func TestCreateNatgateway(t *testing.T) {
	mock := NewMock().On("CreateNatGateway", &ec2.CreateNatGatewayOutput{
		NatGateway: &ec2types.NatGateway{NatGatewayId: awssdk.String("nat-1234")},
	})

	Template("create natgateway elasticip-id=eipalloc-1234 subnet=subnet-1234").
		Mock(mock).
		ExpectCalls("CreateNatGateway").
		ExpectCommandResult("nat-1234").
		ExpectRevert("delete natgateway id=nat-1234").
		Run(t)

	in := mock.InputFor("CreateNatGateway").(*ec2.CreateNatGatewayInput)
	if got := awssdk.ToString(in.AllocationId); got != "eipalloc-1234" {
		t.Errorf("AllocationId: got %q, want eipalloc-1234", got)
	}
}

func TestDeleteNatgateway(t *testing.T) {
	mock := NewMock().On("DeleteNatGateway", &ec2.DeleteNatGatewayOutput{})

	Template("delete natgateway id=nat-1234").
		Mock(mock).
		ExpectCalls("DeleteNatGateway").
		Run(t)
}

func TestCreateElasticip(t *testing.T) {
	mock := NewMock().On("AllocateAddress", &ec2.AllocateAddressOutput{
		AllocationId: awssdk.String("eipalloc-1234"),
		PublicIp:     awssdk.String("52.10.20.30"),
	})

	Template("create elasticip domain=vpc").
		Mock(mock).
		ExpectCalls("AllocateAddress").
		ExpectCommandResult("eipalloc-1234").
		Run(t)
}

// delete elasticip accepts an allocation id or the address itself.
func TestDeleteElasticipById(t *testing.T) {
	mock := NewMock().On("ReleaseAddress", &ec2.ReleaseAddressOutput{})

	Template("delete elasticip id=eipalloc-1234").
		Mock(mock).
		ExpectCalls("ReleaseAddress").
		Run(t)

	in := mock.InputFor("ReleaseAddress").(*ec2.ReleaseAddressInput)
	if got := awssdk.ToString(in.AllocationId); got != "eipalloc-1234" {
		t.Errorf("AllocationId: got %q, want eipalloc-1234", got)
	}
}

func TestAttachElasticip(t *testing.T) {
	mock := NewMock().On("AssociateAddress", &ec2.AssociateAddressOutput{
		AssociationId: awssdk.String("eipassoc-1234"),
	})

	Template("attach elasticip id=eipalloc-1234 instance=i-1234").
		Mock(mock).
		ExpectCalls("AssociateAddress").
		ExpectCommandResult("eipassoc-1234").
		Run(t)

	in := mock.InputFor("AssociateAddress").(*ec2.AssociateAddressInput)
	if got := awssdk.ToString(in.InstanceId); got != "i-1234" {
		t.Errorf("InstanceId: got %q, want i-1234", got)
	}
}

func TestDetachElasticipTakesAssociation(t *testing.T) {
	mock := NewMock().On("DisassociateAddress", &ec2.DisassociateAddressOutput{})

	Template("detach elasticip association=eipassoc-1234").
		Mock(mock).
		ExpectCalls("DisassociateAddress").
		Run(t)

	in := mock.InputFor("DisassociateAddress").(*ec2.DisassociateAddressInput)
	if got := awssdk.ToString(in.AssociationId); got != "eipassoc-1234" {
		t.Errorf("AssociationId: got %q, want eipassoc-1234", got)
	}
}

func TestCreateNetworkinterface(t *testing.T) {
	mock := NewMock().On("CreateNetworkInterface", &ec2.CreateNetworkInterfaceOutput{
		NetworkInterface: &ec2types.NetworkInterface{
			NetworkInterfaceId: awssdk.String("eni-1234"),
		},
	})

	Template("create networkinterface subnet=subnet-1234 description=web privateip=10.0.0.42").
		Mock(mock).
		ExpectCalls("CreateNetworkInterface").
		ExpectCommandResult("eni-1234").
		ExpectRevert("delete networkinterface id=eni-1234").
		Run(t)

	in := mock.InputFor("CreateNetworkInterface").(*ec2.CreateNetworkInterfaceInput)
	if got := awssdk.ToString(in.PrivateIpAddress); got != "10.0.0.42" {
		t.Errorf("PrivateIpAddress: got %q, want 10.0.0.42", got)
	}
}

func TestDeleteNetworkinterface(t *testing.T) {
	mock := NewMock().On("DeleteNetworkInterface", &ec2.DeleteNetworkInterfaceOutput{})

	Template("delete networkinterface id=eni-1234").
		Mock(mock).
		ExpectCalls("DeleteNetworkInterface").
		Run(t)
}

func TestAttachNetworkinterface(t *testing.T) {
	mock := NewMock().On("AttachNetworkInterface", &ec2.AttachNetworkInterfaceOutput{
		AttachmentId: awssdk.String("eni-attach-1234"),
	})

	Template("attach networkinterface id=eni-1234 instance=i-1234 device-index=1").
		Mock(mock).
		ExpectCalls("AttachNetworkInterface").
		ExpectCommandResult("eni-attach-1234").
		Run(t)

	in := mock.InputFor("AttachNetworkInterface").(*ec2.AttachNetworkInterfaceInput)
	if got := awssdk.ToInt32(in.DeviceIndex); got != 1 {
		t.Errorf("DeviceIndex: got %d, want 1", got)
	}
}

func TestUpdateSubnetPublicIPMapping(t *testing.T) {
	mock := NewMock().On("ModifySubnetAttribute", &ec2.ModifySubnetAttributeOutput{})

	Template("update subnet id=subnet-1234 public=true").
		Mock(mock).
		ExpectCalls("ModifySubnetAttribute").
		Run(t)

	in := mock.InputFor("ModifySubnetAttribute").(*ec2.ModifySubnetAttributeInput)
	if in.MapPublicIpOnLaunch == nil || !awssdk.ToBool(in.MapPublicIpOnLaunch.Value) {
		t.Errorf("MapPublicIpOnLaunch: got %#v, want true", in.MapPublicIpOnLaunch)
	}
}

func TestDeleteSubnet(t *testing.T) {
	mock := NewMock().On("DeleteSubnet", &ec2.DeleteSubnetOutput{})

	Template("delete subnet id=subnet-1234").
		Mock(mock).
		ExpectCalls("DeleteSubnet").
		Run(t)
}

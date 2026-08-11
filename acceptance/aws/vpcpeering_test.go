package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestCreateVpcpeering(t *testing.T) {
	mock := NewMock().On("CreateVpcPeeringConnection", &ec2.CreateVpcPeeringConnectionOutput{
		VpcPeeringConnection: &ec2types.VpcPeeringConnection{
			VpcPeeringConnectionId: awssdk.String("pcx-1234abcd"),
		},
	})

	Template("create vpcpeering vpc=vpc-1234 peer-vpc=vpc-5678").
		Mock(mock).
		ExpectCalls("CreateVpcPeeringConnection").
		ExpectCommandResult("pcx-1234abcd").
		ExpectRevert("delete vpcpeering id=pcx-1234abcd").
		Run(t)

	in := mock.InputFor("CreateVpcPeeringConnection").(*ec2.CreateVpcPeeringConnectionInput)
	if got := awssdk.ToString(in.VpcId); got != "vpc-1234" {
		t.Errorf("VpcId: got %q, want vpc-1234", got)
	}
	if got := awssdk.ToString(in.PeerVpcId); got != "vpc-5678" {
		t.Errorf("PeerVpcId: got %q, want vpc-5678", got)
	}
	// Same-account peering needs neither, so they must not be invented.
	if in.PeerOwnerId != nil || in.PeerRegion != nil {
		t.Error("peer owner and region should be unset for a same-account peering")
	}
}

// Cross-account and cross-region peering both need the extra fields, and the request goes to
// the other account for acceptance.
func TestCreateVpcpeeringCrossAccount(t *testing.T) {
	mock := NewMock().On("CreateVpcPeeringConnection", &ec2.CreateVpcPeeringConnectionOutput{
		VpcPeeringConnection: &ec2types.VpcPeeringConnection{
			VpcPeeringConnectionId: awssdk.String("pcx-1234abcd"),
		},
	})

	Template("create vpcpeering vpc=vpc-1234 peer-vpc=vpc-5678 peer-owner=210987654321 peer-region=eu-west-1").
		Mock(mock).
		ExpectCalls("CreateVpcPeeringConnection").
		Run(t)

	in := mock.InputFor("CreateVpcPeeringConnection").(*ec2.CreateVpcPeeringConnectionInput)
	if got := awssdk.ToString(in.PeerOwnerId); got != "210987654321" {
		t.Errorf("PeerOwnerId: got %q", got)
	}
	if got := awssdk.ToString(in.PeerRegion); got != "eu-west-1" {
		t.Errorf("PeerRegion: got %q, want eu-west-1", got)
	}
}

// Accepting is the other half of creating: a connection stays pending until the accepter side
// runs this.
func TestStartVpcpeeringAccepts(t *testing.T) {
	mock := NewMock().On("AcceptVpcPeeringConnection", &ec2.AcceptVpcPeeringConnectionOutput{})

	Template("start vpcpeering id=pcx-1234abcd").
		Mock(mock).
		ExpectCalls("AcceptVpcPeeringConnection").
		Run(t)

	in := mock.InputFor("AcceptVpcPeeringConnection").(*ec2.AcceptVpcPeeringConnectionInput)
	if got := awssdk.ToString(in.VpcPeeringConnectionId); got != "pcx-1234abcd" {
		t.Errorf("VpcPeeringConnectionId: got %q", got)
	}
}

func TestDeleteVpcpeering(t *testing.T) {
	mock := NewMock().On("DeleteVpcPeeringConnection", &ec2.DeleteVpcPeeringConnectionOutput{})

	Template("delete vpcpeering id=pcx-1234abcd").
		Mock(mock).
		ExpectCalls("DeleteVpcPeeringConnection").
		Run(t)

	in := mock.InputFor("DeleteVpcPeeringConnection").(*ec2.DeleteVpcPeeringConnectionInput)
	if got := awssdk.ToString(in.VpcPeeringConnectionId); got != "pcx-1234abcd" {
		t.Errorf("VpcPeeringConnectionId: got %q", got)
	}
}

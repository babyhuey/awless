package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// These are the first acceptance tests this package has ever had. The harness
// previously could not run at all: SDK v2 clients are concrete structs, so the
// generated factory stubbed out all 30 commands with a TODO.
//
// They exercise the whole path — template parse, compile, run, AWS call
// interception, result extraction and revert — with no network access.

func TestCreateVpcIsInterceptedAndResultExtracted(t *testing.T) {
	mock := NewMock().On("CreateVpc", &ec2.CreateVpcOutput{
		Vpc: &ec2types.Vpc{VpcId: awssdk.String("vpc-1234")},
	})

	Template("create vpc cidr=10.0.0.0/16").
		Mock(mock).
		ExpectCalls("CreateVpc").
		ExpectCommandResult("vpc-1234").
		Run(t)

	// The command's input must have reached the mock unmodified.
	in, ok := mock.InputFor("CreateVpc").(*ec2.CreateVpcInput)
	if !ok {
		t.Fatalf("expected a *ec2.CreateVpcInput, got %T", mock.InputFor("CreateVpc"))
	}
	if got := awssdk.ToString(in.CidrBlock); got != "10.0.0.0/16" {
		t.Errorf("CidrBlock: got %q, want 10.0.0.0/16", got)
	}
}

func TestExpectInputIsVerified(t *testing.T) {
	mock := NewMock().On("CreateSubnet", &ec2.CreateSubnetOutput{
		Subnet: &ec2types.Subnet{SubnetId: awssdk.String("subnet-1")},
	})

	Template("create subnet cidr=10.0.0.0/24 vpc=vpc-1234").
		Mock(mock).
		ExpectInput("CreateSubnet", &ec2.CreateSubnetInput{
			CidrBlock: awssdk.String("10.0.0.0/24"),
			VpcId:     awssdk.String("vpc-1234"),
		}).
		ExpectCalls("CreateSubnet").
		Run(t)
}

func TestRevertIsDerived(t *testing.T) {
	mock := NewMock().On("CreateVpc", &ec2.CreateVpcOutput{
		Vpc: &ec2types.Vpc{VpcId: awssdk.String("vpc-9999")},
	})

	Template("create vpc cidr=10.0.0.0/16").
		Mock(mock).
		ExpectCalls("CreateVpc").
		ExpectRevert("delete vpc id=vpc-9999").
		Run(t)
}

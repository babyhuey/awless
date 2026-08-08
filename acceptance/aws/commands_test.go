package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Covers a spread of services through the full template pipeline: parameter
// validation, struct injection into the AWS input, the call itself, result
// extraction and revert derivation. These exercise aws/spec broadly without
// touching the network.

func TestCreateTagsOnInstance(t *testing.T) {
	mock := NewMock().On("CreateTags", &ec2.CreateTagsOutput{})

	Template("create tag resource=i-1234 key=Env value=production").
		Mock(mock).
		ExpectCalls("CreateTags").
		Run(t)

	in := mock.InputFor("CreateTags").(*ec2.CreateTagsInput)
	if len(in.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(in.Tags))
	}
	if got := awssdk.ToString(in.Tags[0].Key); got != "Env" {
		t.Errorf("tag key: got %q, want Env", got)
	}
	if got := awssdk.ToString(in.Tags[0].Value); got != "production" {
		t.Errorf("tag value: got %q, want production", got)
	}
}

func TestDeleteVpcRevert(t *testing.T) {
	mock := NewMock().On("DeleteVpc", &ec2.DeleteVpcOutput{})

	Template("delete vpc id=vpc-abcd").
		Mock(mock).
		ExpectCalls("DeleteVpc").
		Run(t)

	in := mock.InputFor("DeleteVpc").(*ec2.DeleteVpcInput)
	if got := awssdk.ToString(in.VpcId); got != "vpc-abcd" {
		t.Errorf("VpcId: got %q, want vpc-abcd", got)
	}
}

func TestCreateSecurityGroupExtractsId(t *testing.T) {
	mock := NewMock().On("CreateSecurityGroup", &ec2.CreateSecurityGroupOutput{
		GroupId: awssdk.String("sg-5678"),
	})

	Template("create securitygroup vpc=vpc-1 description=web name=web-sg").
		Mock(mock).
		ExpectCalls("CreateSecurityGroup").
		ExpectCommandResult("sg-5678").
		Run(t)
}

func TestCreateSubnetWithAvailabilityZone(t *testing.T) {
	mock := NewMock().On("CreateSubnet", &ec2.CreateSubnetOutput{
		Subnet: &ec2types.Subnet{SubnetId: awssdk.String("subnet-az")},
	})

	Template("create subnet cidr=10.0.1.0/24 vpc=vpc-1 availabilityzone=us-west-2a").
		Mock(mock).
		ExpectCalls("CreateSubnet").
		ExpectCommandResult("subnet-az").
		Run(t)

	in := mock.InputFor("CreateSubnet").(*ec2.CreateSubnetInput)
	if got := awssdk.ToString(in.AvailabilityZone); got != "us-west-2a" {
		t.Errorf("AvailabilityZone: got %q, want us-west-2a", got)
	}
}

func TestCreateBucket(t *testing.T) {
	mock := NewMock().On("CreateBucket", &s3.CreateBucketOutput{})

	Template("create bucket name=my-test-bucket").
		Mock(mock).
		ExpectCalls("CreateBucket").
		Run(t)

	in := mock.InputFor("CreateBucket").(*s3.CreateBucketInput)
	if got := awssdk.ToString(in.Bucket); got != "my-test-bucket" {
		t.Errorf("Bucket: got %q, want my-test-bucket", got)
	}
}

func TestCreateTopicExtractsArn(t *testing.T) {
	arn := "arn:aws:sns:us-west-2:123456789012:alerts"
	mock := NewMock().On("CreateTopic", &sns.CreateTopicOutput{
		TopicArn: awssdk.String(arn),
	})

	Template("create topic name=alerts").
		Mock(mock).
		ExpectCalls("CreateTopic").
		ExpectCommandResult(arn).
		Run(t)
}

func TestCreateQueueExtractsUrl(t *testing.T) {
	url := "https://sqs.us-west-2.amazonaws.com/123456789012/jobs"
	mock := NewMock().On("CreateQueue", &sqs.CreateQueueOutput{
		QueueUrl: awssdk.String(url),
	})

	Template("create queue name=jobs").
		Mock(mock).
		ExpectCalls("CreateQueue").
		ExpectCommandResult(url).
		Run(t)
}

func TestCreateIamUserAndRevert(t *testing.T) {
	mock := NewMock().On("CreateUser", &iam.CreateUserOutput{
		// ExtractResult reads UserId, not UserName.
		User: &iamtypes.User{UserId: awssdk.String("AIDAEXAMPLE"), UserName: awssdk.String("jsmith")},
	})

	Template("create user name=jsmith").
		Mock(mock).
		ExpectCalls("CreateUser").
		ExpectRevert("delete user name=jsmith").
		Run(t)
}

func TestCreateGroup(t *testing.T) {
	mock := NewMock().On("CreateGroup", &iam.CreateGroupOutput{
		Group: &iamtypes.Group{GroupId: awssdk.String("AGPAEXAMPLE"), GroupName: awssdk.String("devs")},
	})

	Template("create group name=devs").
		Mock(mock).
		ExpectCalls("CreateGroup").
		Run(t)
}

// Multiple statements in one template must each be dispatched.
func TestMultiStatementTemplate(t *testing.T) {
	mock := NewMock().
		On("CreateVpc", &ec2.CreateVpcOutput{Vpc: &ec2types.Vpc{VpcId: awssdk.String("vpc-multi")}}).
		On("CreateSubnet", &ec2.CreateSubnetOutput{Subnet: &ec2types.Subnet{SubnetId: awssdk.String("subnet-multi")}})

	Template("myvpc = create vpc cidr=10.0.0.0/16\ncreate subnet cidr=10.0.0.0/24 vpc=$myvpc").
		Mock(mock).
		ExpectCalls("CreateVpc", "CreateSubnet").
		Run(t)

	// The subnet must have received the vpc id produced by the first statement.
	in := mock.InputFor("CreateSubnet").(*ec2.CreateSubnetInput)
	if got := awssdk.ToString(in.VpcId); got != "vpc-multi" {
		t.Errorf("VpcId: got %q, want the id from the first statement (vpc-multi)", got)
	}
}

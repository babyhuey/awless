package awsfetch

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func TestAssignAPIs(t *testing.T) {
	iamClient := &iam.Client{}
	ec2Client := &ec2.Client{}

	conf := NewConfig(iamClient, ec2Client, nil)
	if conf.APIs.IAM == nil {
		t.Fatal("unexpected nil")
	}
	if conf.APIs.Ec2 == nil {
		t.Fatal("unexpected nil")
	}
	if conf.APIs.RDS != nil {
		t.Fatal("expected nil")
	}
	if conf.APIs.ECR != nil {
		t.Fatal("expected nil")
	}
}

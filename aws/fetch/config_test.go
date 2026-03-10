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
	if conf.APIs.Iam == nil {
		t.Fatal("unexpected nil")
	}
	if conf.APIs.Ec2 == nil {
		t.Fatal("unexpected nil")
	}
	if conf.APIs.Rds != nil {
		t.Fatal("expected nil")
	}
	if conf.APIs.Ecr != nil {
		t.Fatal("expected nil")
	}
}

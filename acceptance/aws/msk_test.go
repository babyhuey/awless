package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
)

// The broker node group holds the subnets, instance type, security groups and storage, and
// the storage path is three levels deep. This checks every intermediate struct was built.
func TestCreateKafkacluster(t *testing.T) {
	mock := NewMock().On("CreateCluster", &kafka.CreateClusterOutput{
		ClusterName: awssdk.String("events"),
		ClusterArn:  awssdk.String("arn:aws:kafka:us-west-2:1:cluster/events/abcd"),
	})

	Template("create kafkacluster name=events version=3.6.0 brokers=3 " +
		"subnets=subnet-1,subnet-2,subnet-3 type=kafka.m5.large storage=100 securitygroups=sg-1234").
		Mock(mock).
		ExpectCalls("CreateCluster").
		ExpectCommandResult("events").
		// Delete takes the ARN, which is why the revert reads it from the result rather
		// than reusing the name.
		ExpectRevert("delete kafkacluster arn=events").
		Run(t)

	in := mock.InputFor("CreateCluster").(*kafka.CreateClusterInput)
	if got := awssdk.ToString(in.ClusterName); got != "events" {
		t.Errorf("ClusterName: got %q, want events", got)
	}
	if got := awssdk.ToString(in.KafkaVersion); got != "3.6.0" {
		t.Errorf("KafkaVersion: got %q, want 3.6.0", got)
	}
	if got := awssdk.ToInt32(in.NumberOfBrokerNodes); got != 3 {
		t.Errorf("NumberOfBrokerNodes: got %d, want 3", got)
	}
	if in.BrokerNodeGroupInfo == nil {
		t.Fatal("BrokerNodeGroupInfo was never built")
	}
	if len(in.BrokerNodeGroupInfo.ClientSubnets) != 3 {
		t.Errorf("ClientSubnets: got %v", in.BrokerNodeGroupInfo.ClientSubnets)
	}
	if got := awssdk.ToString(in.BrokerNodeGroupInfo.InstanceType); got != "kafka.m5.large" {
		t.Errorf("InstanceType: got %q", got)
	}
	if len(in.BrokerNodeGroupInfo.SecurityGroups) != 1 {
		t.Errorf("SecurityGroups: got %v", in.BrokerNodeGroupInfo.SecurityGroups)
	}
	// Three levels: BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.VolumeSize.
	if in.BrokerNodeGroupInfo.StorageInfo == nil || in.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo == nil {
		t.Fatal("the nested storage info was not built")
	}
	if got := awssdk.ToInt32(in.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.VolumeSize); got != 100 {
		t.Errorf("VolumeSize: got %d, want 100", got)
	}
}

// A cluster needs its topology stated: AWS requires the broker count to be a multiple of the
// subnet count, so neither can be guessed.
func TestCreateKafkaclusterRequiresTopology(t *testing.T) {
	err := Template("create kafkacluster name=events version=3.6.0").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected brokers, subnets and type to be required")
	}
}

func TestDeleteKafkacluster(t *testing.T) {
	arn := "arn:aws:kafka:us-west-2:123456789012:cluster/events/abcd"
	mock := NewMock().On("DeleteCluster", &kafka.DeleteClusterOutput{})

	Template("delete kafkacluster arn=" + arn).
		Mock(mock).
		ExpectCalls("DeleteCluster").
		Run(t)

	in := mock.InputFor("DeleteCluster").(*kafka.DeleteClusterInput)
	if got := awssdk.ToString(in.ClusterArn); got != arn {
		t.Errorf("ClusterArn: got %q, want %q", got, arn)
	}
}

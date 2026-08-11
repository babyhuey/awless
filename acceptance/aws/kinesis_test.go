package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesistypes "github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

func TestCreateStreamOnDemand(t *testing.T) {
	mock := NewMock().On("CreateStream", &kinesis.CreateStreamOutput{})

	Template("create stream name=clickstream mode=ON_DEMAND").
		Mock(mock).
		ExpectCalls("CreateStream").
		ExpectCommandResult("clickstream").
		ExpectRevert("delete stream name=clickstream").
		Run(t)

	in := mock.InputFor("CreateStream").(*kinesis.CreateStreamInput)
	if got := awssdk.ToString(in.StreamName); got != "clickstream" {
		t.Errorf("StreamName: got %q, want clickstream", got)
	}
	// The mode is nested one level down, so this checks the intermediate struct was
	// created rather than the value dropped.
	if in.StreamModeDetails == nil {
		t.Fatal("StreamModeDetails was never built")
	}
	if got := string(in.StreamModeDetails.StreamMode); got != "ON_DEMAND" {
		t.Errorf("StreamMode: got %q, want ON_DEMAND", got)
	}
	// An ON_DEMAND stream must not carry a shard count; AWS rejects the pair.
	if in.ShardCount != nil {
		t.Errorf("ShardCount should be unset for ON_DEMAND, got %d", awssdk.ToInt32(in.ShardCount))
	}
}

func TestCreateStreamProvisioned(t *testing.T) {
	mock := NewMock().On("CreateStream", &kinesis.CreateStreamOutput{})

	Template("create stream name=clickstream mode=PROVISIONED shards=4").
		Mock(mock).
		ExpectCalls("CreateStream").
		Run(t)

	in := mock.InputFor("CreateStream").(*kinesis.CreateStreamInput)
	if got := awssdk.ToInt32(in.ShardCount); got != 4 {
		t.Errorf("ShardCount: got %d, want 4", got)
	}
	if got := string(in.StreamModeDetails.StreamMode); got != "PROVISIONED" {
		t.Errorf("StreamMode: got %q, want PROVISIONED", got)
	}
}

func TestDeleteStream(t *testing.T) {
	mock := NewMock().On("DeleteStream", &kinesis.DeleteStreamOutput{})

	Template("delete stream name=clickstream force=true").
		Mock(mock).
		ExpectCalls("DeleteStream").
		Run(t)

	in := mock.InputFor("DeleteStream").(*kinesis.DeleteStreamInput)
	if got := awssdk.ToString(in.StreamName); got != "clickstream" {
		t.Errorf("StreamName: got %q, want clickstream", got)
	}
	if !awssdk.ToBool(in.EnforceConsumerDeletion) {
		t.Error("expected EnforceConsumerDeletion to be true")
	}
}

// ScalingType is required by the API and has exactly one accepted value, so it is
// defaulted rather than demanded. Without this the call fails with a validation error the
// user can do nothing useful about.
func TestUpdateStreamDefaultsTheScalingType(t *testing.T) {
	mock := NewMock().On("UpdateShardCount", &kinesis.UpdateShardCountOutput{})

	Template("update stream name=clickstream shards=8").
		Mock(mock).
		ExpectCalls("UpdateShardCount").
		Run(t)

	in := mock.InputFor("UpdateShardCount").(*kinesis.UpdateShardCountInput)
	if got := awssdk.ToInt32(in.TargetShardCount); got != 8 {
		t.Errorf("TargetShardCount: got %d, want 8", got)
	}
	if in.ScalingType != kinesistypes.ScalingTypeUniformScaling {
		t.Errorf("ScalingType: got %q, want UNIFORM_SCALING", in.ScalingType)
	}
}

// An explicit value must still win, so the default cannot mask a future second option.
func TestUpdateStreamKeepsAnExplicitScalingType(t *testing.T) {
	mock := NewMock().On("UpdateShardCount", &kinesis.UpdateShardCountOutput{})

	Template("update stream name=clickstream shards=8 scaling-type=UNIFORM_SCALING").
		Mock(mock).
		ExpectCalls("UpdateShardCount").
		Run(t)

	in := mock.InputFor("UpdateShardCount").(*kinesis.UpdateShardCountInput)
	if in.ScalingType != kinesistypes.ScalingTypeUniformScaling {
		t.Errorf("ScalingType: got %q", in.ScalingType)
	}
}

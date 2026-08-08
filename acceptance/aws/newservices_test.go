package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestCreateDynamodbtableDefaultsToOnDemand(t *testing.T) {
	mock := NewMock().On("CreateTable", &dynamodb.CreateTableOutput{
		TableDescription: &dynamodbtypes.TableDescription{TableName: awssdk.String("users")},
	})

	Template("create dynamodbtable name=users partition-key=id").
		Mock(mock).
		ExpectCalls("CreateTable").
		ExpectCommandResult("users").
		Run(t)

	in := mock.InputFor("CreateTable").(*dynamodb.CreateTableInput)
	if got := in.BillingMode; got != dynamodbtypes.BillingModePayPerRequest {
		t.Errorf("BillingMode: got %q, want PAY_PER_REQUEST", got)
	}
	// PAY_PER_REQUEST rejects a throughput block.
	if in.ProvisionedThroughput != nil {
		t.Error("expected no provisioned throughput for an on-demand table")
	}
	if len(in.KeySchema) != 1 || in.KeySchema[0].KeyType != dynamodbtypes.KeyTypeHash {
		t.Errorf("expected a single HASH key, got %#v", in.KeySchema)
	}
	// The attribute type defaults to string when not given.
	if got := in.AttributeDefinitions[0].AttributeType; got != dynamodbtypes.ScalarAttributeTypeS {
		t.Errorf("AttributeType: got %q, want S", got)
	}
}

func TestCreateDynamodbtableWithSortKey(t *testing.T) {
	mock := NewMock().On("CreateTable", &dynamodb.CreateTableOutput{
		TableDescription: &dynamodbtypes.TableDescription{TableName: awssdk.String("events")},
	})

	Template("create dynamodbtable name=events partition-key=id sort-key=ts sort-type=N").
		Mock(mock).
		ExpectCalls("CreateTable").
		Run(t)

	in := mock.InputFor("CreateTable").(*dynamodb.CreateTableInput)
	if len(in.KeySchema) != 2 {
		t.Fatalf("expected 2 key schema elements, got %d", len(in.KeySchema))
	}
	if in.KeySchema[1].KeyType != dynamodbtypes.KeyTypeRange {
		t.Errorf("second key: got %q, want RANGE", in.KeySchema[1].KeyType)
	}
	// The key schema and attribute definitions must stay in agreement.
	if len(in.AttributeDefinitions) != 2 {
		t.Fatalf("expected 2 attribute definitions, got %d", len(in.AttributeDefinitions))
	}
	if got := in.AttributeDefinitions[1].AttributeType; got != dynamodbtypes.ScalarAttributeTypeN {
		t.Errorf("sort attribute type: got %q, want N", got)
	}
}

func TestCreateDynamodbtableProvisioned(t *testing.T) {
	mock := NewMock().On("CreateTable", &dynamodb.CreateTableOutput{
		TableDescription: &dynamodbtypes.TableDescription{TableName: awssdk.String("sessions")},
	})

	Template("create dynamodbtable name=sessions partition-key=id billing-mode=PROVISIONED read-capacity=10 write-capacity=5").
		Mock(mock).
		ExpectCalls("CreateTable").
		Run(t)

	in := mock.InputFor("CreateTable").(*dynamodb.CreateTableInput)
	if in.ProvisionedThroughput == nil {
		t.Fatal("expected a provisioned throughput block")
	}
	if got := awssdk.ToInt64(in.ProvisionedThroughput.ReadCapacityUnits); got != 10 {
		t.Errorf("read capacity: got %d, want 10", got)
	}
	if got := awssdk.ToInt64(in.ProvisionedThroughput.WriteCapacityUnits); got != 5 {
		t.Errorf("write capacity: got %d, want 5", got)
	}
}

// Capacity without PROVISIONED is a request the API would reject, so it must fail
// before the call.
func TestCreateDynamodbtableRejectsCapacityOnDemand(t *testing.T) {
	mock := NewMock().On("CreateTable", &dynamodb.CreateTableOutput{
		TableDescription: &dynamodbtypes.TableDescription{TableName: awssdk.String("bad")},
	})

	if err := Template("create dynamodbtable name=bad partition-key=id billing-mode=PAY_PER_REQUEST read-capacity=10").Mock(mock).RunExpectingError(t); err == nil {
		t.Error("expected read-capacity with PAY_PER_REQUEST to be refused")
	}
	if calls := mock.Calls()["CreateTable"]; calls != 0 {
		t.Errorf("expected no API call, got %d", calls)
	}
}

func TestCreateEkscluster(t *testing.T) {
	mock := NewMock().On("CreateCluster", &eks.CreateClusterOutput{
		Cluster: &ekstypes.Cluster{Name: awssdk.String("prod")},
	})

	Template("create ekscluster name=prod role=arn:aws:iam::123456789012:role/eksRole subnets=subnet-a,subnet-b").
		Mock(mock).
		ExpectCalls("CreateCluster").
		ExpectCommandResult("prod").
		Run(t)

	in := mock.InputFor("CreateCluster").(*eks.CreateClusterInput)
	if got := len(in.ResourcesVpcConfig.SubnetIds); got != 2 {
		t.Errorf("expected 2 subnet ids, got %d", got)
	}
	if got := awssdk.ToString(in.RoleArn); got != "arn:aws:iam::123456789012:role/eksRole" {
		t.Errorf("RoleArn: got %q", got)
	}
}

// EKS requires subnets in two availability zones; one subnet must fail locally
// rather than producing an opaque API error.
func TestCreateEksclusterRejectsSingleSubnet(t *testing.T) {
	mock := NewMock().On("CreateCluster", &eks.CreateClusterOutput{
		Cluster: &ekstypes.Cluster{Name: awssdk.String("prod")},
	})

	if err := Template("create ekscluster name=prod role=arn:aws:iam::123456789012:role/eksRole subnets=subnet-a").Mock(mock).RunExpectingError(t); err == nil {
		t.Error("expected a single subnet to be refused")
	}
	if calls := mock.Calls()["CreateCluster"]; calls != 0 {
		t.Errorf("expected no API call, got %d", calls)
	}
}

func TestCreateEksnodegroupDerivesScalingConfig(t *testing.T) {
	mock := NewMock().On("CreateNodegroup", &eks.CreateNodegroupOutput{
		Nodegroup: &ekstypes.Nodegroup{NodegroupName: awssdk.String("workers")},
	})

	Template("create eksnodegroup name=workers cluster=prod role=arn:aws:iam::123456789012:role/nodeRole subnets=subnet-a instance-type=t3.medium").
		Mock(mock).
		ExpectCalls("CreateNodegroup").
		ExpectCommandResult("workers").
		Run(t)

	in := mock.InputFor("CreateNodegroup").(*eks.CreateNodegroupInput)
	// All three sizes are required together, so the unspecified ones follow desired.
	if got := awssdk.ToInt32(in.ScalingConfig.DesiredSize); got != 2 {
		t.Errorf("DesiredSize: got %d, want the default 2", got)
	}
	if got := awssdk.ToInt32(in.ScalingConfig.MinSize); got != 2 {
		t.Errorf("MinSize: got %d, want 2", got)
	}
	if got := awssdk.ToInt32(in.ScalingConfig.MaxSize); got != 2 {
		t.Errorf("MaxSize: got %d, want 2", got)
	}
	if len(in.InstanceTypes) != 1 || in.InstanceTypes[0] != "t3.medium" {
		t.Errorf("InstanceTypes: got %#v, want [t3.medium]", in.InstanceTypes)
	}
}

func TestCreateEksnodegroupRejectsInconsistentSizes(t *testing.T) {
	mock := NewMock().On("CreateNodegroup", &eks.CreateNodegroupOutput{
		Nodegroup: &ekstypes.Nodegroup{NodegroupName: awssdk.String("workers")},
	})

	if err := Template("create eksnodegroup name=workers cluster=prod role=arn:aws:iam::1:role/r subnets=subnet-a min-size=5 max-size=3 desired-size=4").Mock(mock).RunExpectingError(t); err == nil {
		t.Error("expected min-size above max-size to be refused")
	}
}

func TestDeleteEksnodegroupNeedsCluster(t *testing.T) {
	mock := NewMock().On("DeleteNodegroup", &eks.DeleteNodegroupOutput{})

	Template("delete eksnodegroup name=workers cluster=prod").
		Mock(mock).
		ExpectCalls("DeleteNodegroup").
		Run(t)

	in := mock.InputFor("DeleteNodegroup").(*eks.DeleteNodegroupInput)
	if got := awssdk.ToString(in.ClusterName); got != "prod" {
		t.Errorf("ClusterName: got %q, want prod", got)
	}
}

func TestCreateFilesystem(t *testing.T) {
	mock := NewMock().On("CreateFileSystem", &efs.CreateFileSystemOutput{
		FileSystemId: awssdk.String("fs-abc123"),
	})

	Template("create filesystem token=my-fs encrypted=true performance-mode=maxIO").
		Mock(mock).
		ExpectCalls("CreateFileSystem").
		ExpectCommandResult("fs-abc123").
		Run(t)

	in := mock.InputFor("CreateFileSystem").(*efs.CreateFileSystemInput)
	if got := awssdk.ToString(in.CreationToken); got != "my-fs" {
		t.Errorf("CreationToken: got %q, want my-fs", got)
	}
	if !awssdk.ToBool(in.Encrypted) {
		t.Error("expected Encrypted to be true")
	}
	if got := string(in.PerformanceMode); got != "maxIO" {
		t.Errorf("PerformanceMode: got %q, want maxIO", got)
	}
}

func TestDeleteFilesystem(t *testing.T) {
	mock := NewMock().On("DeleteFileSystem", &efs.DeleteFileSystemOutput{})

	Template("delete filesystem id=fs-abc123").
		Mock(mock).
		ExpectCalls("DeleteFileSystem").
		Run(t)

	in := mock.InputFor("DeleteFileSystem").(*efs.DeleteFileSystemInput)
	if got := awssdk.ToString(in.FileSystemId); got != "fs-abc123" {
		t.Errorf("FileSystemId: got %q, want fs-abc123", got)
	}
}

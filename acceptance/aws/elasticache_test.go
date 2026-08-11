package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	elasticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

// These assert the input mapping rather than only that a call happened. Params reach the
// AWS request struct by reflection over `awsName` tags, which the compiler does not
// check, and a name that does not resolve is dropped silently.

func TestCreateCachecluster(t *testing.T) {
	mock := NewMock().On("CreateCacheCluster", &elasticache.CreateCacheClusterOutput{
		CacheCluster: &elasticachetypes.CacheCluster{
			CacheClusterId: awssdk.String("sessions"),
		},
	})

	Template("create cachecluster id=sessions engine=redis type=cache.t3.micro nodes=2 " +
		"subnet-group=cache-private port=6379 zone=us-west-2a securitygroups=sg-1234,sg-5678 " +
		"engine-version=7.1 parameter-group=default.redis7 snapshot-retention=5").
		Mock(mock).
		ExpectCalls("CreateCacheCluster").
		ExpectCommandResult("sessions").
		ExpectRevert("delete cachecluster id=sessions").
		Run(t)

	in := mock.InputFor("CreateCacheCluster").(*elasticache.CreateCacheClusterInput)
	if got := awssdk.ToString(in.CacheClusterId); got != "sessions" {
		t.Errorf("CacheClusterId: got %q, want sessions", got)
	}
	if got := awssdk.ToString(in.Engine); got != "redis" {
		t.Errorf("Engine: got %q, want redis", got)
	}
	if got := awssdk.ToString(in.CacheNodeType); got != "cache.t3.micro" {
		t.Errorf("CacheNodeType: got %q, want cache.t3.micro", got)
	}
	// Declared as *int64 on the spec and *int32 on the request; the setter converts.
	if got := awssdk.ToInt32(in.NumCacheNodes); got != 2 {
		t.Errorf("NumCacheNodes: got %d, want 2", got)
	}
	if got := awssdk.ToInt32(in.Port); got != 6379 {
		t.Errorf("Port: got %d, want 6379", got)
	}
	if got := awssdk.ToString(in.CacheSubnetGroupName); got != "cache-private" {
		t.Errorf("CacheSubnetGroupName: got %q, want cache-private", got)
	}
	if got := awssdk.ToString(in.PreferredAvailabilityZone); got != "us-west-2a" {
		t.Errorf("PreferredAvailabilityZone: got %q, want us-west-2a", got)
	}
	if got := awssdk.ToString(in.EngineVersion); got != "7.1" {
		t.Errorf("EngineVersion: got %q, want 7.1", got)
	}
	if got := awssdk.ToString(in.CacheParameterGroupName); got != "default.redis7" {
		t.Errorf("CacheParameterGroupName: got %q, want default.redis7", got)
	}
	if got := awssdk.ToInt32(in.SnapshotRetentionLimit); got != 5 {
		t.Errorf("SnapshotRetentionLimit: got %d, want 5", got)
	}
	// A list of strings: v1 modeled these as []*string and v2 uses []string, which is
	// what made 35 params panic during the SDK migration.
	if len(in.SecurityGroupIds) != 2 || in.SecurityGroupIds[0] != "sg-1234" || in.SecurityGroupIds[1] != "sg-5678" {
		t.Errorf("SecurityGroupIds: got %v, want [sg-1234 sg-5678]", in.SecurityGroupIds)
	}
}

// The replication-group form is mutually exclusive with engine and node type, so it
// exercises the other branch of the params spec.
func TestCreateCacheclusterAsReplica(t *testing.T) {
	mock := NewMock().On("CreateCacheCluster", &elasticache.CreateCacheClusterOutput{
		CacheCluster: &elasticachetypes.CacheCluster{
			CacheClusterId: awssdk.String("sessions-replica-1"),
		},
	})

	Template("create cachecluster id=sessions-replica-1 replication-group=sessions-group").
		Mock(mock).
		ExpectCalls("CreateCacheCluster").
		ExpectCommandResult("sessions-replica-1").
		Run(t)

	in := mock.InputFor("CreateCacheCluster").(*elasticache.CreateCacheClusterInput)
	if got := awssdk.ToString(in.ReplicationGroupId); got != "sessions-group" {
		t.Errorf("ReplicationGroupId: got %q, want sessions-group", got)
	}
	if in.Engine != nil {
		t.Errorf("Engine should be unset on a replica, got %q", awssdk.ToString(in.Engine))
	}
}

// Supplying both forms must be rejected before any AWS call is made.
func TestCreateCacheclusterRejectsBothForms(t *testing.T) {
	err := Template("create cachecluster id=sessions engine=redis type=cache.t3.micro replication-group=grp").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected the standalone and replication-group forms to be mutually exclusive")
	}
}

// An empty CacheCluster in the response must not nil-dereference: ExtractResult falls
// back to the supplied id. Six commands panicked on exactly this shape before.
func TestCreateCacheclusterEmptyResponse(t *testing.T) {
	mock := NewMock().On("CreateCacheCluster", &elasticache.CreateCacheClusterOutput{})

	Template("create cachecluster id=sessions engine=redis type=cache.t3.micro").
		Mock(mock).
		ExpectCalls("CreateCacheCluster").
		ExpectCommandResult("sessions").
		Run(t)
}

func TestDeleteCachecluster(t *testing.T) {
	mock := NewMock().On("DeleteCacheCluster", &elasticache.DeleteCacheClusterOutput{})

	Template("delete cachecluster id=sessions snapshot=sessions-final").
		Mock(mock).
		ExpectCalls("DeleteCacheCluster").
		Run(t)

	in := mock.InputFor("DeleteCacheCluster").(*elasticache.DeleteCacheClusterInput)
	if got := awssdk.ToString(in.CacheClusterId); got != "sessions" {
		t.Errorf("CacheClusterId: got %q, want sessions", got)
	}
	if got := awssdk.ToString(in.FinalSnapshotIdentifier); got != "sessions-final" {
		t.Errorf("FinalSnapshotIdentifier: got %q, want sessions-final", got)
	}
}

func TestUpdateCachecluster(t *testing.T) {
	mock := NewMock().On("ModifyCacheCluster", &elasticache.ModifyCacheClusterOutput{})

	Template("update cachecluster id=sessions nodes=3 type=cache.t3.small engine-version=7.1 " +
		"maintenance-window=sun:23:00-mon:01:30 apply-immediately=true securitygroups=sg-1234").
		Mock(mock).
		ExpectCalls("ModifyCacheCluster").
		Run(t)

	in := mock.InputFor("ModifyCacheCluster").(*elasticache.ModifyCacheClusterInput)
	if got := awssdk.ToInt32(in.NumCacheNodes); got != 3 {
		t.Errorf("NumCacheNodes: got %d, want 3", got)
	}
	if got := awssdk.ToString(in.CacheNodeType); got != "cache.t3.small" {
		t.Errorf("CacheNodeType: got %q, want cache.t3.small", got)
	}
	if got := awssdk.ToString(in.PreferredMaintenanceWindow); got != "sun:23:00-mon:01:30" {
		t.Errorf("PreferredMaintenanceWindow: got %q", got)
	}
	if !awssdk.ToBool(in.ApplyImmediately) {
		t.Error("expected ApplyImmediately to be true")
	}
	if len(in.SecurityGroupIds) != 1 || in.SecurityGroupIds[0] != "sg-1234" {
		t.Errorf("SecurityGroupIds: got %v, want [sg-1234]", in.SecurityGroupIds)
	}
}

func TestCreateCachesubnetgroup(t *testing.T) {
	mock := NewMock().On("CreateCacheSubnetGroup", &elasticache.CreateCacheSubnetGroupOutput{
		CacheSubnetGroup: &elasticachetypes.CacheSubnetGroup{
			CacheSubnetGroupName: awssdk.String("cache-private"),
		},
	})

	Template("create cachesubnetgroup name=cache-private subnets=subnet-1234,subnet-5678 description=Private").
		Mock(mock).
		ExpectCalls("CreateCacheSubnetGroup").
		ExpectCommandResult("cache-private").
		ExpectRevert("delete cachesubnetgroup name=cache-private").
		Run(t)

	in := mock.InputFor("CreateCacheSubnetGroup").(*elasticache.CreateCacheSubnetGroupInput)
	if got := awssdk.ToString(in.CacheSubnetGroupName); got != "cache-private" {
		t.Errorf("CacheSubnetGroupName: got %q, want cache-private", got)
	}
	if got := awssdk.ToString(in.CacheSubnetGroupDescription); got != "Private" {
		t.Errorf("CacheSubnetGroupDescription: got %q, want Private", got)
	}
	if len(in.SubnetIds) != 2 || in.SubnetIds[0] != "subnet-1234" {
		t.Errorf("SubnetIds: got %v, want [subnet-1234 subnet-5678]", in.SubnetIds)
	}
}

func TestDeleteCachesubnetgroup(t *testing.T) {
	mock := NewMock().On("DeleteCacheSubnetGroup", &elasticache.DeleteCacheSubnetGroupOutput{})

	Template("delete cachesubnetgroup name=cache-private").
		Mock(mock).
		ExpectCalls("DeleteCacheSubnetGroup").
		Run(t)

	in := mock.InputFor("DeleteCacheSubnetGroup").(*elasticache.DeleteCacheSubnetGroupInput)
	if got := awssdk.ToString(in.CacheSubnetGroupName); got != "cache-private" {
		t.Errorf("CacheSubnetGroupName: got %q, want cache-private", got)
	}
}

func TestUpdateCachesubnetgroup(t *testing.T) {
	mock := NewMock().On("ModifyCacheSubnetGroup", &elasticache.ModifyCacheSubnetGroupOutput{})

	Template("update cachesubnetgroup name=cache-private subnets=subnet-1234,subnet-9abc description=Updated").
		Mock(mock).
		ExpectCalls("ModifyCacheSubnetGroup").
		Run(t)

	in := mock.InputFor("ModifyCacheSubnetGroup").(*elasticache.ModifyCacheSubnetGroupInput)
	if got := awssdk.ToString(in.CacheSubnetGroupDescription); got != "Updated" {
		t.Errorf("CacheSubnetGroupDescription: got %q, want Updated", got)
	}
	if len(in.SubnetIds) != 2 || in.SubnetIds[1] != "subnet-9abc" {
		t.Errorf("SubnetIds: got %v, want [subnet-1234 subnet-9abc]", in.SubnetIds)
	}
}

func TestCreateReplicationgroup(t *testing.T) {
	mock := NewMock().On("CreateReplicationGroup", &elasticache.CreateReplicationGroupOutput{
		ReplicationGroup: &elasticachetypes.ReplicationGroup{
			ReplicationGroupId: awssdk.String("sessions-group"),
		},
	})

	Template("create replicationgroup id=sessions-group description=Sessions engine=redis " +
		"type=cache.t3.micro clusters=3 automatic-failover=true multi-az=true " +
		"at-rest-encryption=true transit-encryption=true port=6379 subnet-group=cache-private").
		Mock(mock).
		ExpectCalls("CreateReplicationGroup").
		ExpectCommandResult("sessions-group").
		ExpectRevert("delete replicationgroup id=sessions-group").
		Run(t)

	in := mock.InputFor("CreateReplicationGroup").(*elasticache.CreateReplicationGroupInput)
	if got := awssdk.ToString(in.ReplicationGroupId); got != "sessions-group" {
		t.Errorf("ReplicationGroupId: got %q, want sessions-group", got)
	}
	// The description has its own AWS field name, unlike most resources here.
	if got := awssdk.ToString(in.ReplicationGroupDescription); got != "Sessions" {
		t.Errorf("ReplicationGroupDescription: got %q, want Sessions", got)
	}
	if got := awssdk.ToInt32(in.NumCacheClusters); got != 3 {
		t.Errorf("NumCacheClusters: got %d, want 3", got)
	}
	if !awssdk.ToBool(in.AutomaticFailoverEnabled) {
		t.Error("expected AutomaticFailoverEnabled to be true")
	}
	if !awssdk.ToBool(in.MultiAZEnabled) {
		t.Error("expected MultiAZEnabled to be true")
	}
	if !awssdk.ToBool(in.AtRestEncryptionEnabled) {
		t.Error("expected AtRestEncryptionEnabled to be true")
	}
	if !awssdk.ToBool(in.TransitEncryptionEnabled) {
		t.Error("expected TransitEncryptionEnabled to be true")
	}
}

func TestDeleteReplicationgroup(t *testing.T) {
	mock := NewMock().On("DeleteReplicationGroup", &elasticache.DeleteReplicationGroupOutput{})

	Template("delete replicationgroup id=sessions-group retain-primary=true snapshot=final").
		Mock(mock).
		ExpectCalls("DeleteReplicationGroup").
		Run(t)

	in := mock.InputFor("DeleteReplicationGroup").(*elasticache.DeleteReplicationGroupInput)
	if got := awssdk.ToString(in.ReplicationGroupId); got != "sessions-group" {
		t.Errorf("ReplicationGroupId: got %q, want sessions-group", got)
	}
	if !awssdk.ToBool(in.RetainPrimaryCluster) {
		t.Error("expected RetainPrimaryCluster to be true")
	}
	if got := awssdk.ToString(in.FinalSnapshotIdentifier); got != "final" {
		t.Errorf("FinalSnapshotIdentifier: got %q, want final", got)
	}
}

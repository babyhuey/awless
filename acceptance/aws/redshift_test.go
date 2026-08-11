package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	redshifttypes "github.com/aws/aws-sdk-go-v2/service/redshift/types"
)

func TestCreateRedshiftcluster(t *testing.T) {
	mock := NewMock().On("CreateCluster", &redshift.CreateClusterOutput{
		Cluster: &redshifttypes.Cluster{ClusterIdentifier: awssdk.String("analytics")},
	})

	Template("create redshiftcluster id=analytics username=admin type=ra3.xlplus " +
		"manage-password=true nodes=4 encrypted=true subnet-group=warehouse-private " +
		"securitygroups=sg-1234,sg-5678 dbname=warehouse zone=us-west-2a").
		Mock(mock).
		ExpectCalls("CreateCluster").
		ExpectCommandResult("analytics").
		ExpectRevert("delete redshiftcluster id=analytics").
		Run(t)

	in := mock.InputFor("CreateCluster").(*redshift.CreateClusterInput)
	if got := awssdk.ToString(in.ClusterIdentifier); got != "analytics" {
		t.Errorf("ClusterIdentifier: got %q, want analytics", got)
	}
	if got := awssdk.ToString(in.MasterUsername); got != "admin" {
		t.Errorf("MasterUsername: got %q, want admin", got)
	}
	if got := awssdk.ToString(in.NodeType); got != "ra3.xlplus" {
		t.Errorf("NodeType: got %q, want ra3.xlplus", got)
	}
	if !awssdk.ToBool(in.ManageMasterPassword) {
		t.Error("expected ManageMasterPassword to be true")
	}
	// No password must reach the request when Redshift is managing it, or the two
	// conflict and AWS rejects the call.
	if in.MasterUserPassword != nil {
		t.Error("MasterUserPassword should be unset when manage-password is used")
	}
	if got := awssdk.ToInt32(in.NumberOfNodes); got != 4 {
		t.Errorf("NumberOfNodes: got %d, want 4", got)
	}
	if got := awssdk.ToString(in.DBName); got != "warehouse" {
		t.Errorf("DBName: got %q, want warehouse", got)
	}
	if len(in.VpcSecurityGroupIds) != 2 || in.VpcSecurityGroupIds[0] != "sg-1234" {
		t.Errorf("VpcSecurityGroupIds: got %v", in.VpcSecurityGroupIds)
	}
	if !awssdk.ToBool(in.Encrypted) {
		t.Error("expected Encrypted to be true")
	}
}

// A password and a managed password are mutually exclusive, and with neither the call
// cannot succeed. Both cases must be caught before reaching AWS.
func TestCreateRedshiftclusterCredentialIsExclusive(t *testing.T) {
	t.Run("both is rejected", func(t *testing.T) {
		err := Template("create redshiftcluster id=a username=admin type=ra3.xlplus password=Sup3rSecret1 manage-password=true").
			Mock(NewMock()).
			RunExpectingError(t)
		if err == nil {
			t.Fatal("expected password and manage-password to be mutually exclusive")
		}
	})

	t.Run("an explicit password is passed through", func(t *testing.T) {
		mock := NewMock().On("CreateCluster", &redshift.CreateClusterOutput{})
		Template("create redshiftcluster id=a username=admin type=ra3.xlplus password=Sup3rSecret1").
			Mock(mock).
			ExpectCalls("CreateCluster").
			Run(t)
		in := mock.InputFor("CreateCluster").(*redshift.CreateClusterInput)
		if got := awssdk.ToString(in.MasterUserPassword); got != "Sup3rSecret1" {
			t.Errorf("MasterUserPassword: got %q", got)
		}
	})
}

// Redshift will not delete without being told what to do about a final snapshot, and
// skipping one silently would destroy a warehouse's only backup.
func TestDeleteRedshiftclusterRequiresASnapshotDecision(t *testing.T) {
	t.Run("with a final snapshot", func(t *testing.T) {
		mock := NewMock().On("DeleteCluster", &redshift.DeleteClusterOutput{})
		Template("delete redshiftcluster id=analytics snapshot=analytics-final").
			Mock(mock).
			ExpectCalls("DeleteCluster").
			Run(t)
		in := mock.InputFor("DeleteCluster").(*redshift.DeleteClusterInput)
		if got := awssdk.ToString(in.FinalClusterSnapshotIdentifier); got != "analytics-final" {
			t.Errorf("FinalClusterSnapshotIdentifier: got %q", got)
		}
	})

	t.Run("explicitly skipping it", func(t *testing.T) {
		mock := NewMock().On("DeleteCluster", &redshift.DeleteClusterOutput{})
		Template("delete redshiftcluster id=analytics skip-snapshot=true").
			Mock(mock).
			ExpectCalls("DeleteCluster").
			Run(t)
		in := mock.InputFor("DeleteCluster").(*redshift.DeleteClusterInput)
		if !awssdk.ToBool(in.SkipFinalClusterSnapshot) {
			t.Error("expected SkipFinalClusterSnapshot to be true")
		}
	})

	t.Run("neither is rejected", func(t *testing.T) {
		err := Template("delete redshiftcluster id=analytics").
			Mock(NewMock()).
			RunExpectingError(t)
		if err == nil {
			t.Fatal("expected a delete with no snapshot decision to be rejected")
		}
	})
}

func TestUpdateRedshiftcluster(t *testing.T) {
	mock := NewMock().On("ModifyCluster", &redshift.ModifyClusterOutput{})

	Template("update redshiftcluster id=analytics nodes=8 type=ra3.4xlarge public=false").
		Mock(mock).
		ExpectCalls("ModifyCluster").
		Run(t)

	in := mock.InputFor("ModifyCluster").(*redshift.ModifyClusterInput)
	if got := awssdk.ToInt32(in.NumberOfNodes); got != 8 {
		t.Errorf("NumberOfNodes: got %d, want 8", got)
	}
	if got := awssdk.ToString(in.NodeType); got != "ra3.4xlarge" {
		t.Errorf("NodeType: got %q, want ra3.4xlarge", got)
	}
	if awssdk.ToBool(in.PubliclyAccessible) {
		t.Error("expected PubliclyAccessible to be false")
	}
}

func TestCreateRedshiftsubnetgroup(t *testing.T) {
	mock := NewMock().On("CreateClusterSubnetGroup", &redshift.CreateClusterSubnetGroupOutput{
		ClusterSubnetGroup: &redshifttypes.ClusterSubnetGroup{
			ClusterSubnetGroupName: awssdk.String("warehouse-private"),
		},
	})

	Template("create redshiftsubnetgroup name=warehouse-private description=Private subnets=subnet-1234,subnet-5678").
		Mock(mock).
		ExpectCalls("CreateClusterSubnetGroup").
		ExpectCommandResult("warehouse-private").
		ExpectRevert("delete redshiftsubnetgroup name=warehouse-private").
		Run(t)

	in := mock.InputFor("CreateClusterSubnetGroup").(*redshift.CreateClusterSubnetGroupInput)
	if got := awssdk.ToString(in.Description); got != "Private" {
		t.Errorf("Description: got %q, want Private", got)
	}
	if len(in.SubnetIds) != 2 || in.SubnetIds[1] != "subnet-5678" {
		t.Errorf("SubnetIds: got %v", in.SubnetIds)
	}
}

func TestDeleteRedshiftsubnetgroup(t *testing.T) {
	mock := NewMock().On("DeleteClusterSubnetGroup", &redshift.DeleteClusterSubnetGroupOutput{})

	Template("delete redshiftsubnetgroup name=warehouse-private").
		Mock(mock).
		ExpectCalls("DeleteClusterSubnetGroup").
		Run(t)

	in := mock.InputFor("DeleteClusterSubnetGroup").(*redshift.DeleteClusterSubnetGroupInput)
	if got := awssdk.ToString(in.ClusterSubnetGroupName); got != "warehouse-private" {
		t.Errorf("ClusterSubnetGroupName: got %q", got)
	}
}

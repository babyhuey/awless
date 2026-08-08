package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// EC2 compute and storage, plus RDS lifecycle.

func TestDeleteInstance(t *testing.T) {
	mock := NewMock().On("TerminateInstances", &ec2.TerminateInstancesOutput{})

	Template("delete instance ids=i-1234,i-5678").
		Mock(mock).
		ExpectCalls("TerminateInstances").
		Run(t)

	in := mock.InputFor("TerminateInstances").(*ec2.TerminateInstancesInput)
	if len(in.InstanceIds) != 2 {
		t.Errorf("expected both ids, got %#v", in.InstanceIds)
	}
}

func TestUpdateInstanceType(t *testing.T) {
	mock := NewMock().On("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeOutput{})

	Template("update instance id=i-1234 type=t3.large").
		Mock(mock).
		ExpectCalls("ModifyInstanceAttribute").
		Run(t)

	in := mock.InputFor("ModifyInstanceAttribute").(*ec2.ModifyInstanceAttributeInput)
	if in.InstanceType == nil || awssdk.ToString(in.InstanceType.Value) != "t3.large" {
		t.Errorf("InstanceType: got %#v, want t3.large", in.InstanceType)
	}
}

// lock sets DisableApiTermination, which is a different attribute from type; sending
// both in one call is rejected by AWS.
func TestUpdateInstanceLock(t *testing.T) {
	mock := NewMock().On("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeOutput{})

	Template("update instance id=i-1234 lock=true").
		Mock(mock).
		ExpectCalls("ModifyInstanceAttribute").
		Run(t)

	in := mock.InputFor("ModifyInstanceAttribute").(*ec2.ModifyInstanceAttributeInput)
	if in.DisableApiTermination == nil || !awssdk.ToBool(in.DisableApiTermination.Value) {
		t.Errorf("DisableApiTermination: got %#v, want true", in.DisableApiTermination)
	}
}

func TestCreateImage(t *testing.T) {
	mock := NewMock().On("CreateImage", &ec2.CreateImageOutput{
		ImageId: awssdk.String("ami-1234"),
	})

	Template("create image instance=i-1234 name=my-ami description=backup").
		Mock(mock).
		ExpectCalls("CreateImage").
		ExpectCommandResult("ami-1234").
		ExpectRevert("delete image id=ami-1234").
		Run(t)

	in := mock.InputFor("CreateImage").(*ec2.CreateImageInput)
	if got := awssdk.ToString(in.InstanceId); got != "i-1234" {
		t.Errorf("InstanceId: got %q, want i-1234", got)
	}
	if got := awssdk.ToString(in.Name); got != "my-ami" {
		t.Errorf("Name: got %q, want my-ami", got)
	}
}

func TestCopyImageAcrossRegions(t *testing.T) {
	mock := NewMock().On("CopyImage", &ec2.CopyImageOutput{
		ImageId: awssdk.String("ami-copy"),
	})

	Template("copy image name=my-copy source-id=ami-1234 source-region=us-east-1").
		Mock(mock).
		ExpectCalls("CopyImage").
		ExpectCommandResult("ami-copy").
		Run(t)

	in := mock.InputFor("CopyImage").(*ec2.CopyImageInput)
	// The source region is what makes this a cross-region copy; losing it would
	// silently copy within the current region.
	if got := awssdk.ToString(in.SourceRegion); got != "us-east-1" {
		t.Errorf("SourceRegion: got %q, want us-east-1", got)
	}
	if got := awssdk.ToString(in.SourceImageId); got != "ami-1234" {
		t.Errorf("SourceImageId: got %q, want ami-1234", got)
	}
}

func TestCreateSnapshot(t *testing.T) {
	mock := NewMock().On("CreateSnapshot", &ec2.CreateSnapshotOutput{
		SnapshotId: awssdk.String("snap-1234"),
	})

	Template("create snapshot volume=vol-1234 description=nightly").
		Mock(mock).
		ExpectCalls("CreateSnapshot").
		ExpectCommandResult("snap-1234").
		ExpectRevert("delete snapshot id=snap-1234").
		Run(t)

	in := mock.InputFor("CreateSnapshot").(*ec2.CreateSnapshotInput)
	if got := awssdk.ToString(in.VolumeId); got != "vol-1234" {
		t.Errorf("VolumeId: got %q, want vol-1234", got)
	}
}

func TestCopySnapshot(t *testing.T) {
	mock := NewMock().On("CopySnapshot", &ec2.CopySnapshotOutput{
		SnapshotId: awssdk.String("snap-copy"),
	})

	Template("copy snapshot source-id=snap-1234 source-region=us-east-1 encrypted=true").
		Mock(mock).
		ExpectCalls("CopySnapshot").
		ExpectCommandResult("snap-copy").
		Run(t)

	in := mock.InputFor("CopySnapshot").(*ec2.CopySnapshotInput)
	if got := awssdk.ToString(in.SourceRegion); got != "us-east-1" {
		t.Errorf("SourceRegion: got %q, want us-east-1", got)
	}
	if !awssdk.ToBool(in.Encrypted) {
		t.Error("expected Encrypted to be set")
	}
}

func TestDeleteSnapshotAndVolume(t *testing.T) {
	snap := NewMock().On("DeleteSnapshot", &ec2.DeleteSnapshotOutput{})
	Template("delete snapshot id=snap-1234").Mock(snap).ExpectCalls("DeleteSnapshot").Run(t)

	vol := NewMock().On("DeleteVolume", &ec2.DeleteVolumeOutput{})
	Template("delete volume id=vol-1234").Mock(vol).ExpectCalls("DeleteVolume").Run(t)

	in := vol.InputFor("DeleteVolume").(*ec2.DeleteVolumeInput)
	if got := awssdk.ToString(in.VolumeId); got != "vol-1234" {
		t.Errorf("VolumeId: got %q, want vol-1234", got)
	}
}

func TestDetachVolume(t *testing.T) {
	mock := NewMock().On("DetachVolume", &ec2.DetachVolumeOutput{})

	Template("detach volume id=vol-1234 instance=i-1234 device=/dev/sdh force=true").
		Mock(mock).
		ExpectCalls("DetachVolume").
		Run(t)

	in := mock.InputFor("DetachVolume").(*ec2.DetachVolumeInput)
	if got := awssdk.ToString(in.Device); got != "/dev/sdh" {
		t.Errorf("Device: got %q, want /dev/sdh", got)
	}
	if !awssdk.ToBool(in.Force) {
		t.Error("expected Force to be set")
	}
}

func TestCreateSecuritygroup(t *testing.T) {
	mock := NewMock().On("CreateSecurityGroup", &ec2.CreateSecurityGroupOutput{
		GroupId: awssdk.String("sg-1234"),
	})

	Template("create securitygroup name=web description=web-traffic vpc=vpc-1234").
		Mock(mock).
		ExpectCalls("CreateSecurityGroup").
		ExpectCommandResult("sg-1234").
		// The revert waits for the group to become unused before deleting it, since
		// AWS refuses to delete a group that is still attached to anything.
		ExpectRevert("check securitygroup id=sg-1234 state=unused timeout=300\ndelete securitygroup id=sg-1234").
		Run(t)

	in := mock.InputFor("CreateSecurityGroup").(*ec2.CreateSecurityGroupInput)
	if got := awssdk.ToString(in.Description); got != "web-traffic" {
		t.Errorf("Description: got %q, want web-traffic", got)
	}
}

// update securitygroup authorizes or revokes a rule; inbound and outbound go to
// different API calls, and getting them crossed would open the wrong direction.
func TestUpdateSecuritygroupInbound(t *testing.T) {
	mock := NewMock().On("AuthorizeSecurityGroupIngress", &ec2.AuthorizeSecurityGroupIngressOutput{})

	Template("update securitygroup id=sg-1234 inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=443").
		Mock(mock).
		ExpectCalls("AuthorizeSecurityGroupIngress").
		Run(t)

	in := mock.InputFor("AuthorizeSecurityGroupIngress").(*ec2.AuthorizeSecurityGroupIngressInput)
	if got := awssdk.ToString(in.GroupId); got != "sg-1234" {
		t.Errorf("GroupId: got %q, want sg-1234", got)
	}
	if len(in.IpPermissions) != 1 {
		t.Fatalf("expected one permission, got %#v", in.IpPermissions)
	}
	p := in.IpPermissions[0]
	if got := awssdk.ToString(p.IpProtocol); got != "tcp" {
		t.Errorf("IpProtocol: got %q, want tcp", got)
	}
	if got := awssdk.ToInt32(p.FromPort); got != 443 {
		t.Errorf("FromPort: got %d, want 443", got)
	}
	if got := awssdk.ToInt32(p.ToPort); got != 443 {
		t.Errorf("ToPort: got %d, want 443", got)
	}
}

func TestUpdateSecuritygroupOutboundRevoke(t *testing.T) {
	mock := NewMock().On("RevokeSecurityGroupEgress", &ec2.RevokeSecurityGroupEgressOutput{})

	Template("update securitygroup id=sg-1234 outbound=revoke protocol=tcp cidr=10.0.0.0/8 portrange=80").
		Mock(mock).
		ExpectCalls("RevokeSecurityGroupEgress").
		Run(t)
}

// A port range spelled from-to must widen the permission, not set both ends equal.
func TestUpdateSecuritygroupPortRange(t *testing.T) {
	mock := NewMock().On("AuthorizeSecurityGroupIngress", &ec2.AuthorizeSecurityGroupIngressOutput{})

	Template("update securitygroup id=sg-1234 inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=8000-8100").
		Mock(mock).
		ExpectCalls("AuthorizeSecurityGroupIngress").
		Run(t)

	in := mock.InputFor("AuthorizeSecurityGroupIngress").(*ec2.AuthorizeSecurityGroupIngressInput)
	p := in.IpPermissions[0]
	if got := awssdk.ToInt32(p.FromPort); got != 8000 {
		t.Errorf("FromPort: got %d, want 8000", got)
	}
	if got := awssdk.ToInt32(p.ToPort); got != 8100 {
		t.Errorf("ToPort: got %d, want 8100", got)
	}
}

func TestDeleteSecuritygroup(t *testing.T) {
	mock := NewMock().On("DeleteSecurityGroup", &ec2.DeleteSecurityGroupOutput{})

	Template("delete securitygroup id=sg-1234").
		Mock(mock).
		ExpectCalls("DeleteSecurityGroup").
		Run(t)
}

func TestDetachSecuritygroupKeepsOthers(t *testing.T) {
	mock := NewMock().
		On("DescribeInstanceAttribute", &ec2.DescribeInstanceAttributeOutput{
			Groups: []ec2types.GroupIdentifier{
				{GroupId: awssdk.String("sg-keep")},
				{GroupId: awssdk.String("sg-remove")},
			},
		}).
		On("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeOutput{})

	Template("detach securitygroup id=sg-remove instance=i-1234").
		Mock(mock).
		ExpectCalls("DescribeInstanceAttribute", "ModifyInstanceAttribute").
		Run(t)

	in := mock.InputFor("ModifyInstanceAttribute").(*ec2.ModifyInstanceAttributeInput)
	// Detaching one group must leave the others attached.
	if len(in.Groups) != 1 || in.Groups[0] != "sg-keep" {
		t.Errorf("expected only sg-keep to remain, got %#v", in.Groups)
	}
}

func TestDeleteKeypair(t *testing.T) {
	mock := NewMock().On("DeleteKeyPair", &ec2.DeleteKeyPairOutput{})

	Template("delete keypair name=my-key").
		Mock(mock).
		ExpectCalls("DeleteKeyPair").
		Run(t)
}

func TestDeleteTag(t *testing.T) {
	mock := NewMock().On("DeleteTags", &ec2.DeleteTagsOutput{})

	Template("delete tag resource=i-1234 key=Env value=staging").
		Mock(mock).
		ExpectCalls("DeleteTags").
		Run(t)

	in := mock.InputFor("DeleteTags").(*ec2.DeleteTagsInput)
	if len(in.Tags) != 1 {
		t.Fatalf("expected one tag, got %#v", in.Tags)
	}
	if got := awssdk.ToString(in.Tags[0].Key); got != "Env" {
		t.Errorf("tag key: got %q, want Env", got)
	}
}

func TestDeleteVpc(t *testing.T) {
	mock := NewMock().On("DeleteVpc", &ec2.DeleteVpcOutput{})

	Template("delete vpc id=vpc-1234").
		Mock(mock).
		ExpectCalls("DeleteVpc").
		Run(t)
}

func TestStartStopRestartDatabase(t *testing.T) {
	start := NewMock().On("StartDBInstance", &rds.StartDBInstanceOutput{
		DBInstance: &rdstypes.DBInstance{DBInstanceIdentifier: awssdk.String("mydb")},
	})
	Template("start database id=mydb").Mock(start).ExpectCalls("StartDBInstance").Run(t)

	stop := NewMock().On("StopDBInstance", &rds.StopDBInstanceOutput{
		DBInstance: &rdstypes.DBInstance{DBInstanceIdentifier: awssdk.String("mydb")},
	})
	Template("stop database id=mydb").Mock(stop).ExpectCalls("StopDBInstance").Run(t)

	restart := NewMock().On("RebootDBInstance", &rds.RebootDBInstanceOutput{
		DBInstance: &rdstypes.DBInstance{DBInstanceIdentifier: awssdk.String("mydb")},
	})
	Template("restart database id=mydb with-failover=true").
		Mock(restart).
		ExpectCalls("RebootDBInstance").
		Run(t)

	in := restart.InputFor("RebootDBInstance").(*rds.RebootDBInstanceInput)
	if !awssdk.ToBool(in.ForceFailover) {
		t.Error("expected ForceFailover to be set")
	}
}

// Deleting a database without a final snapshot is destructive and must be explicit.
func TestDeleteDatabaseSkipSnapshot(t *testing.T) {
	mock := NewMock().On("DeleteDBInstance", &rds.DeleteDBInstanceOutput{
		DBInstance: &rdstypes.DBInstance{DBInstanceIdentifier: awssdk.String("mydb")},
	})

	Template("delete database id=mydb skip-snapshot=true").
		Mock(mock).
		ExpectCalls("DeleteDBInstance").
		Run(t)

	in := mock.InputFor("DeleteDBInstance").(*rds.DeleteDBInstanceInput)
	if !awssdk.ToBool(in.SkipFinalSnapshot) {
		t.Error("expected SkipFinalSnapshot to be set")
	}
}

func TestCreateDbsubnetgroup(t *testing.T) {
	mock := NewMock().On("CreateDBSubnetGroup", &rds.CreateDBSubnetGroupOutput{
		DBSubnetGroup: &rdstypes.DBSubnetGroup{DBSubnetGroupName: awssdk.String("my-subnet-group")},
	})

	Template("create dbsubnetgroup name=my-subnet-group description=for-rds subnets=subnet-a,subnet-b").
		Mock(mock).
		ExpectCalls("CreateDBSubnetGroup").
		Run(t)

	in := mock.InputFor("CreateDBSubnetGroup").(*rds.CreateDBSubnetGroupInput)
	if len(in.SubnetIds) != 2 {
		t.Errorf("expected both subnets, got %#v", in.SubnetIds)
	}
}

func TestDeleteDbsubnetgroup(t *testing.T) {
	mock := NewMock().On("DeleteDBSubnetGroup", &rds.DeleteDBSubnetGroupOutput{})

	Template("delete dbsubnetgroup name=my-subnet-group").
		Mock(mock).
		ExpectCalls("DeleteDBSubnetGroup").
		Run(t)
}

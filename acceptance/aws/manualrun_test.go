package awsat

import (
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Commands with a ManualRun assemble their own AWS request, so the mapping from
// template params to API input is hand-written and cannot be checked by the
// compiler. These cover that mapping.

func TestCreateCertificateUsesFirstDomainAsName(t *testing.T) {
	arn := "arn:aws:acm:us-west-2:123456789012:certificate/abcd"
	mock := NewMock().On("RequestCertificate", &acm.RequestCertificateOutput{
		CertificateArn: awssdk.String(arn),
	})

	Template("create certificate domains=example.com,www.example.com").
		Mock(mock).
		ExpectCalls("RequestCertificate").
		ExpectCommandResult(arn).
		Run(t)

	in := mock.InputFor("RequestCertificate").(*acm.RequestCertificateInput)
	// The first domain becomes DomainName; the rest become SANs.
	if got := awssdk.ToString(in.DomainName); got != "example.com" {
		t.Errorf("DomainName: got %q, want example.com", got)
	}
	if len(in.SubjectAlternativeNames) != 1 || in.SubjectAlternativeNames[0] != "www.example.com" {
		t.Errorf("SubjectAlternativeNames: got %#v, want [www.example.com]", in.SubjectAlternativeNames)
	}
}

// A single domain must not produce an empty SAN list entry.
func TestCreateCertificateSingleDomain(t *testing.T) {
	mock := NewMock().On("RequestCertificate", &acm.RequestCertificateOutput{
		CertificateArn: awssdk.String("arn:aws:acm:us-west-2:1:certificate/x"),
	})

	Template("create certificate domains=example.com").
		Mock(mock).
		ExpectCalls("RequestCertificate").
		Run(t)

	in := mock.InputFor("RequestCertificate").(*acm.RequestCertificateInput)
	if len(in.SubjectAlternativeNames) != 0 {
		t.Errorf("expected no SANs for a single domain, got %#v", in.SubjectAlternativeNames)
	}
}

func TestCreateDatabase(t *testing.T) {
	mock := NewMock().On("CreateDBInstance", &rds.CreateDBInstanceOutput{
		DBInstance: &rdstypes.DBInstance{DBInstanceIdentifier: awssdk.String("mydb")},
	})

	Template("create database id=mydb engine=postgres password=s3cretPass size=20 type=db.t3.micro username=admin").
		Mock(mock).
		ExpectCalls("CreateDBInstance").
		Run(t)

	in := mock.InputFor("CreateDBInstance").(*rds.CreateDBInstanceInput)
	if got := awssdk.ToString(in.Engine); got != "postgres" {
		t.Errorf("Engine: got %q, want postgres", got)
	}
	if got := awssdk.ToInt32(in.AllocatedStorage); got != 20 {
		t.Errorf("AllocatedStorage: got %d, want 20", got)
	}
}

func TestCreateRecordBuildsChangeBatch(t *testing.T) {
	mock := NewMock().On("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53types.ChangeInfo{Id: awssdk.String("C123")},
	})

	Template(`create record zone=/hostedzone/Z1 name=www.example.com type=A ttl=300 values=1.2.3.4,5.6.7.8`).
		Mock(mock).
		ExpectCalls("ChangeResourceRecordSets").
		Run(t)

	in := mock.InputFor("ChangeResourceRecordSets").(*route53.ChangeResourceRecordSetsInput)
	if len(in.ChangeBatch.Changes) != 1 {
		t.Fatalf("expected one change, got %d", len(in.ChangeBatch.Changes))
	}
	rrs := in.ChangeBatch.Changes[0].ResourceRecordSet
	if got := awssdk.ToString(rrs.Name); got != "www.example.com" {
		t.Errorf("Name: got %q", got)
	}
	if got := awssdk.ToInt64(rrs.TTL); got != 300 {
		t.Errorf("TTL: got %d, want 300", got)
	}
	// Both values must become separate resource records.
	if len(rrs.ResourceRecords) != 2 {
		t.Errorf("expected 2 resource records, got %d", len(rrs.ResourceRecords))
	}
}

func TestCreateRoleBuildsTrustPolicy(t *testing.T) {
	mock := NewMock().On("CreateRole", &iam.CreateRoleOutput{
		Role: &iamtypes.Role{RoleName: awssdk.String("my-role"), Arn: awssdk.String("arn:aws:iam::123456789012:role/my-role")},
	})

	Template("create role name=my-role principal-service=ec2.amazonaws.com").
		Mock(mock).
		ExpectCalls("CreateRole", "CreateInstanceProfile").
		ExpectCommandResult("arn:aws:iam::123456789012:role/my-role").
		Run(t)

	in := mock.InputFor("CreateRole").(*iam.CreateRoleInput)
	doc := awssdk.ToString(in.AssumeRolePolicyDocument)
	// The trust policy is assembled as JSON; IAM requires these exact keys.
	for _, want := range []string{`"Version"`, `"Statement"`, `"sts:AssumeRole"`, "ec2.amazonaws.com"} {
		if !strings.Contains(doc, want) {
			t.Errorf("trust policy missing %s\n  got: %s", want, doc)
		}
	}
}

func TestAttachSecuritygroup(t *testing.T) {
	mock := NewMock().
		On("DescribeInstanceAttribute", &ec2.DescribeInstanceAttributeOutput{
			Groups: []ec2types.GroupIdentifier{{GroupId: awssdk.String("sg-existing")}},
		}).
		On("ModifyInstanceAttribute", &ec2.ModifyInstanceAttributeOutput{})

	Template("attach securitygroup id=sg-new instance=i-1234").
		Mock(mock).
		ExpectCalls("DescribeInstanceAttribute", "ModifyInstanceAttribute").
		Run(t)

	in := mock.InputFor("ModifyInstanceAttribute").(*ec2.ModifyInstanceAttributeInput)
	// Attaching must preserve the groups already on the instance.
	if len(in.Groups) != 2 {
		t.Fatalf("expected both the existing and new group, got %#v", in.Groups)
	}
	var found bool
	for _, g := range in.Groups {
		if g == "sg-existing" {
			found = true
		}
	}
	if !found {
		t.Errorf("attaching dropped the existing group: %#v", in.Groups)
	}
}

func TestDeleteImageWithSnapshots(t *testing.T) {
	mock := NewMock().
		On("DescribeImages", &ec2.DescribeImagesOutput{
			Images: []ec2types.Image{{
				ImageId: awssdk.String("ami-1234"),
				BlockDeviceMappings: []ec2types.BlockDeviceMapping{
					{Ebs: &ec2types.EbsBlockDevice{SnapshotId: awssdk.String("snap-1")}},
				},
			}},
		}).
		On("DeregisterImage", &ec2.DeregisterImageOutput{}).
		On("DeleteSnapshot", &ec2.DeleteSnapshotOutput{})

	Template("delete image id=ami-1234 delete-snapshots=true").
		Mock(mock).
		ExpectCalls("DescribeImages", "DeregisterImage", "DeleteSnapshot").
		Run(t)
}

// Without delete-snapshots the snapshots must be left alone.
func TestDeleteImageKeepsSnapshotsByDefault(t *testing.T) {
	mock := NewMock().
		On("DescribeImages", &ec2.DescribeImagesOutput{
			Images: []ec2types.Image{{ImageId: awssdk.String("ami-1234")}},
		}).
		On("DeregisterImage", &ec2.DeregisterImageOutput{})

	Template("delete image id=ami-1234").
		Mock(mock).
		ExpectCalls("DeregisterImage").
		Run(t)

	if n := mock.Calls()["DeleteSnapshot"]; n != 0 {
		t.Errorf("expected no snapshot deletion, got %d calls", n)
	}
}

func TestCreateVolume(t *testing.T) {
	mock := NewMock().On("CreateVolume", &ec2.CreateVolumeOutput{
		VolumeId: awssdk.String("vol-1234"),
	})

	Template("create volume availabilityzone=us-west-2a size=100").
		Mock(mock).
		ExpectCalls("CreateVolume").
		ExpectCommandResult("vol-1234").
		ExpectRevert("delete volume id=vol-1234").
		Run(t)

	in := mock.InputFor("CreateVolume").(*ec2.CreateVolumeInput)
	if got := awssdk.ToInt32(in.Size); got != 100 {
		t.Errorf("Size: got %d, want 100", got)
	}
}

func TestAttachVolume(t *testing.T) {
	mock := NewMock().On("AttachVolume", &ec2.AttachVolumeOutput{
		Device: awssdk.String("/dev/sdh"),
	})

	Template("attach volume id=vol-1234 instance=i-1234 device=/dev/sdh").
		Mock(mock).
		ExpectCalls("AttachVolume").
		Run(t)

	in := mock.InputFor("AttachVolume").(*ec2.AttachVolumeInput)
	if got := awssdk.ToString(in.Device); got != "/dev/sdh" {
		t.Errorf("Device: got %q, want /dev/sdh", got)
	}
}

func TestCreateKeypair(t *testing.T) {
	// The command writes the private key under this directory, so it must exist.
	t.Setenv("__AWLESS_KEYS_DIR", t.TempDir())

	mock := NewMock().On("ImportKeyPair", &ec2.ImportKeyPairOutput{
		KeyName: awssdk.String("my-key"),
	})

	Template("create keypair name=my-key").
		Mock(mock).
		ExpectCalls("ImportKeyPair").
		Run(t)
}

func TestStartAndStopInstance(t *testing.T) {
	start := NewMock().On("StartInstances", &ec2.StartInstancesOutput{})
	Template("start instance ids=i-1234").Mock(start).ExpectCalls("StartInstances").Run(t)

	in := start.InputFor("StartInstances").(*ec2.StartInstancesInput)
	if len(in.InstanceIds) != 1 || in.InstanceIds[0] != "i-1234" {
		t.Errorf("InstanceIds: got %#v", in.InstanceIds)
	}

	stop := NewMock().On("StopInstances", &ec2.StopInstancesOutput{})
	Template("stop instance ids=i-1234,i-5678").Mock(stop).ExpectCalls("StopInstances").Run(t)

	stopIn := stop.InputFor("StopInstances").(*ec2.StopInstancesInput)
	if len(stopIn.InstanceIds) != 2 {
		t.Errorf("expected both ids, got %#v", stopIn.InstanceIds)
	}
}

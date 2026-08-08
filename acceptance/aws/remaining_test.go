package awsat

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// The remaining commands: target group membership, CloudFront, Lambda, EC2 instance
// creation and the last one-offs.

func TestAttachAndDetachInstanceToTargetgroup(t *testing.T) {
	tg := "arn:aws:elasticloadbalancing:us-west-2:1:targetgroup/my-tg/abcd"

	attach := NewMock().On("RegisterTargets", &elasticloadbalancingv2.RegisterTargetsOutput{})
	Template("attach instance id=i-1234 targetgroup=" + tg + " port=8080").
		Mock(attach).
		ExpectCalls("RegisterTargets").
		Run(t)

	in := attach.InputFor("RegisterTargets").(*elasticloadbalancingv2.RegisterTargetsInput)
	if len(in.Targets) != 1 {
		t.Fatalf("expected one target, got %#v", in.Targets)
	}
	if got := awssdk.ToString(in.Targets[0].Id); got != "i-1234" {
		t.Errorf("target id: got %q, want i-1234", got)
	}
	// A port override lets one instance serve several target groups.
	if got := awssdk.ToInt32(in.Targets[0].Port); got != 8080 {
		t.Errorf("target port: got %d, want 8080", got)
	}

	detach := NewMock().On("DeregisterTargets", &elasticloadbalancingv2.DeregisterTargetsOutput{})
	Template("detach instance id=i-1234 targetgroup=" + tg).
		Mock(detach).
		ExpectCalls("DeregisterTargets").
		Run(t)
}

func TestAttachListenerCertificate(t *testing.T) {
	mock := NewMock().On("AddListenerCertificates", &elasticloadbalancingv2.AddListenerCertificatesOutput{})

	Template("attach listener id=arn:aws:elasticloadbalancing:us-west-2:1:listener/app/my-lb/abcd/ef01 certificate=arn:aws:acm:us-west-2:1:certificate/abcd").
		Mock(mock).
		ExpectCalls("AddListenerCertificates").
		Run(t)

	in := mock.InputFor("AddListenerCertificates").(*elasticloadbalancingv2.AddListenerCertificatesInput)
	if len(in.Certificates) != 1 {
		t.Errorf("Certificates: got %#v, want one entry", in.Certificates)
	}
}

func TestUpdateTargetgroupHealthCheck(t *testing.T) {
	mock := NewMock().On("ModifyTargetGroup", &elasticloadbalancingv2.ModifyTargetGroupOutput{})

	Template("update targetgroup id=arn:aws:elasticloadbalancing:us-west-2:1:targetgroup/my-tg/abcd healthcheckinterval=15 healthchecktimeout=4").
		Mock(mock).
		ExpectCalls("ModifyTargetGroup").
		Run(t)

	in := mock.InputFor("ModifyTargetGroup").(*elasticloadbalancingv2.ModifyTargetGroupInput)
	if got := awssdk.ToInt32(in.HealthCheckIntervalSeconds); got != 15 {
		t.Errorf("HealthCheckIntervalSeconds: got %d, want 15", got)
	}
	if got := awssdk.ToInt32(in.HealthCheckTimeoutSeconds); got != 4 {
		t.Errorf("HealthCheckTimeoutSeconds: got %d, want 4", got)
	}
}

func TestCreateDistribution(t *testing.T) {
	mock := NewMock().On("CreateDistribution", &cloudfront.CreateDistributionOutput{
		Distribution: &cloudfronttypes.Distribution{
			Id:         awssdk.String("E1PA6795SAMPLE"),
			DomainName: awssdk.String("d123.cloudfront.net"),
		},
	})

	Template("create distribution origin-domain=my-bucket.s3.amazonaws.com default-file=index.html https-behavior=redirect-to-https min-ttl=3600").
		Mock(mock).
		ExpectCalls("CreateDistribution").
		ExpectCommandResult("E1PA6795SAMPLE").
		Run(t)

	in := mock.InputFor("CreateDistribution").(*cloudfront.CreateDistributionInput)
	cfg := in.DistributionConfig
	if cfg == nil {
		t.Fatal("expected a distribution config")
	}
	if got := awssdk.ToString(cfg.DefaultRootObject); got != "index.html" {
		t.Errorf("DefaultRootObject: got %q, want index.html", got)
	}
	// The origin has to be nested inside Origins.Items, not set at the top level.
	if cfg.Origins == nil || len(cfg.Origins.Items) != 1 {
		t.Fatalf("expected one origin, got %#v", cfg.Origins)
	}
	if got := awssdk.ToString(cfg.Origins.Items[0].DomainName); got != "my-bucket.s3.amazonaws.com" {
		t.Errorf("origin domain: got %q", got)
	}
	if cfg.DefaultCacheBehavior == nil {
		t.Fatal("expected a default cache behavior")
	}
	if got := string(cfg.DefaultCacheBehavior.ViewerProtocolPolicy); got != "redirect-to-https" {
		t.Errorf("ViewerProtocolPolicy: got %q, want redirect-to-https", got)
	}
}

// update distribution has to read the current config and its ETag first, since
// CloudFront requires the whole config on every update.
func TestUpdateDistribution(t *testing.T) {
	mock := NewMock().
		On("GetDistribution", &cloudfront.GetDistributionOutput{
			ETag: awssdk.String("E2ETAG"),
			Distribution: &cloudfronttypes.Distribution{
				Id:     awssdk.String("E1PA6795SAMPLE"),
				Status: awssdk.String("Deployed"),
				DistributionConfig: &cloudfronttypes.DistributionConfig{
					CallerReference: awssdk.String("ref"),
					Comment:         awssdk.String("old"),
					Enabled:         awssdk.Bool(true),
				},
			},
		}).
		On("UpdateDistribution", &cloudfront.UpdateDistributionOutput{
			Distribution: &cloudfronttypes.Distribution{Id: awssdk.String("E1PA6795SAMPLE")},
		})

	Template("update distribution id=E1PA6795SAMPLE comment=new").
		Mock(mock).
		ExpectCalls("GetDistribution", "UpdateDistribution").
		Run(t)

	in := mock.InputFor("UpdateDistribution").(*cloudfront.UpdateDistributionInput)
	// Without the ETag from the read, CloudFront rejects the update.
	if got := awssdk.ToString(in.IfMatch); got != "E2ETAG" {
		t.Errorf("IfMatch: got %q, want the ETag from GetDistribution", got)
	}
	if got := awssdk.ToString(in.DistributionConfig.Comment); got != "new" {
		t.Errorf("Comment: got %q, want new", got)
	}
}

func TestDeleteDistribution(t *testing.T) {
	mock := NewMock().
		On("GetDistribution", &cloudfront.GetDistributionOutput{
			ETag: awssdk.String("E2ETAG"),
			Distribution: &cloudfronttypes.Distribution{
				Id:     awssdk.String("E1PA6795SAMPLE"),
				Status: awssdk.String("Deployed"),
				DistributionConfig: &cloudfronttypes.DistributionConfig{
					CallerReference: awssdk.String("ref"),
					Enabled:         awssdk.Bool(false),
				},
			},
		}).
		// A distribution must be disabled and deployed before it can be deleted, so
		// the command updates it first.
		On("UpdateDistribution", &cloudfront.UpdateDistributionOutput{
			Distribution: &cloudfronttypes.Distribution{Id: awssdk.String("E1PA6795SAMPLE")},
		}).
		On("DeleteDistribution", &cloudfront.DeleteDistributionOutput{})

	Template("delete distribution id=E1PA6795SAMPLE").
		Mock(mock).
		ExpectCalls("GetDistribution", "GetDistribution", "UpdateDistribution", "DeleteDistribution").
		Run(t)
}

func TestCreateFunctionFromZip(t *testing.T) {
	dir := t.TempDir()
	zip := filepath.Join(dir, "fn.zip")
	if err := os.WriteFile(zip, []byte("not-really-a-zip"), 0600); err != nil {
		t.Fatal(err)
	}

	mock := NewMock().On("CreateFunction", &lambda.CreateFunctionOutput{
		FunctionArn:  awssdk.String("arn:aws:lambda:us-west-2:1:function:my-fn"),
		FunctionName: awssdk.String("my-fn"),
	})

	Template("create function name=my-fn handler=index.handler role=arn:aws:iam::1:role/lambda runtime=python3.12 zipfile=" + zip + " memory=256 timeout=30").
		Mock(mock).
		ExpectCalls("CreateFunction").
		Run(t)

	in := mock.InputFor("CreateFunction").(*lambda.CreateFunctionInput)
	if got := string(in.Runtime); got != "python3.12" {
		t.Errorf("Runtime: got %q, want python3.12", got)
	}
	if got := awssdk.ToInt32(in.MemorySize); got != 256 {
		t.Errorf("MemorySize: got %d, want 256", got)
	}
	if got := awssdk.ToInt32(in.Timeout); got != 30 {
		t.Errorf("Timeout: got %d, want 30", got)
	}
	// The zip has to be read off disk and sent as bytes.
	if in.Code == nil || len(in.Code.ZipFile) == 0 {
		t.Errorf("expected the zip contents in Code.ZipFile, got %#v", in.Code)
	}
}

func TestCreateInstance(t *testing.T) {
	mock := NewMock().
		On("RunInstances", &ec2.RunInstancesOutput{
			Instances: []ec2types.Instance{{InstanceId: awssdk.String("i-new")}},
		}).
		On("CreateTags", &ec2.CreateTagsOutput{})

	Template("create instance image=ami-1234 type=t3.micro count=1 name=web subnet=subnet-1234 keypair=my-key").
		Mock(mock).
		ExpectCalls("RunInstances", "CreateTags").
		ExpectCommandResult("i-new").
		Run(t)

	in := mock.InputFor("RunInstances").(*ec2.RunInstancesInput)
	// count maps to both min and max, or AWS may launch fewer than asked.
	if got := awssdk.ToInt32(in.MinCount); got != 1 {
		t.Errorf("MinCount: got %d, want 1", got)
	}
	if got := awssdk.ToInt32(in.MaxCount); got != 1 {
		t.Errorf("MaxCount: got %d, want 1", got)
	}
	if got := string(in.InstanceType); got != "t3.micro" {
		t.Errorf("InstanceType: got %q, want t3.micro", got)
	}
	if got := awssdk.ToString(in.SubnetId); got != "subnet-1234" {
		t.Errorf("SubnetId: got %q, want subnet-1234", got)
	}
}

func TestCreateLaunchconfiguration(t *testing.T) {
	mock := NewMock().On("CreateLaunchConfiguration", &autoscaling.CreateLaunchConfigurationOutput{})

	Template("create launchconfiguration name=my-lc image=ami-1234 type=t3.micro keypair=my-key").
		Mock(mock).
		ExpectCalls("CreateLaunchConfiguration").
		Run(t)
}

func TestCreateMfadevice(t *testing.T) {
	// The command writes the seed QR code into this directory.
	t.Setenv("__AWLESS_KEYS_DIR", t.TempDir())

	mock := NewMock().On("CreateVirtualMFADevice", &iam.CreateVirtualMFADeviceOutput{
		VirtualMFADevice: &iamtypes.VirtualMFADevice{
			SerialNumber:     awssdk.String("arn:aws:iam::1:mfa/jsmith"),
			Base32StringSeed: []byte("SEED"),
		},
	})

	Template("create mfadevice name=jsmith").
		Mock(mock).
		ExpectCalls("CreateVirtualMFADevice").
		Run(t)
}

func TestCreateRepository(t *testing.T) {
	mock := NewMock().On("CreateRepository", &ecr.CreateRepositoryOutput{
		Repository: &ecrtypes.Repository{
			RepositoryName: awssdk.String("my-repo"),
			RepositoryArn:  awssdk.String("arn:aws:ecr:us-west-2:1:repository/my-repo"),
			RepositoryUri:  awssdk.String("1.dkr.ecr.us-west-2.amazonaws.com/my-repo"),
		},
	})

	Template("create repository name=my-repo").
		Mock(mock).
		ExpectCalls("CreateRepository").
		ExpectRevert("delete repository name=my-repo").
		Run(t)
}

func TestCreateS3object(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "report.txt")
	if err := os.WriteFile(file, []byte("contents"), 0600); err != nil {
		t.Fatal(err)
	}

	mock := NewMock().On("PutObject", &s3.PutObjectOutput{})

	Template("create s3object bucket=my-bucket file=" + file).
		Mock(mock).
		ExpectCalls("PutObject").
		Run(t)

	in := mock.InputFor("PutObject").(*s3.PutObjectInput)
	// Without an explicit name the key defaults to the file's base name.
	if got := awssdk.ToString(in.Key); got != "report.txt" {
		t.Errorf("Key: got %q, want report.txt", got)
	}
	if got := awssdk.ToString(in.Bucket); got != "my-bucket" {
		t.Errorf("Bucket: got %q, want my-bucket", got)
	}
}

func TestCreateZone(t *testing.T) {
	mock := NewMock().On("CreateHostedZone", &route53.CreateHostedZoneOutput{
		HostedZone: &route53types.HostedZone{Id: awssdk.String("/hostedzone/Z1NEW")},
	})

	Template("create zone name=example.com callerreference=ref-2024 comment=managed").
		Mock(mock).
		ExpectCalls("CreateHostedZone").
		ExpectCommandResult("/hostedzone/Z1NEW").
		Run(t)

	in := mock.InputFor("CreateHostedZone").(*route53.CreateHostedZoneInput)
	// CallerReference is Route53's idempotency key and is required.
	if got := awssdk.ToString(in.CallerReference); got != "ref-2024" {
		t.Errorf("CallerReference: got %q, want ref-2024", got)
	}
}

func TestDeleteUser(t *testing.T) {
	mock := NewMock().On("DeleteUser", &iam.DeleteUserOutput{})

	Template("delete user name=jsmith").
		Mock(mock).
		ExpectCalls("DeleteUser").
		Run(t)
}

func TestDetachAlarmRemovesOneAction(t *testing.T) {
	mock := NewMock().
		On("DescribeAlarms", &cloudwatch.DescribeAlarmsOutput{
			MetricAlarms: []cloudwatchtypes.MetricAlarm{{
				AlarmName:    awssdk.String("high-cpu"),
				AlarmArn:     awssdk.String("arn:aws:cloudwatch:us-west-2:1:alarm:high-cpu"),
				AlarmActions: []string{"arn:aws:sns:us-west-2:1:remove", "arn:aws:sns:us-west-2:1:keep"},
			}},
		}).
		On("PutMetricAlarm", &cloudwatch.PutMetricAlarmOutput{})

	Template("detach alarm name=high-cpu action-arn=arn:aws:sns:us-west-2:1:remove").
		Mock(mock).
		ExpectCalls("DescribeAlarms", "PutMetricAlarm").
		Run(t)
}

func TestDetachInstanceprofile(t *testing.T) {
	mock := NewMock().
		On("DescribeIamInstanceProfileAssociations", &ec2.DescribeIamInstanceProfileAssociationsOutput{
			IamInstanceProfileAssociations: []ec2types.IamInstanceProfileAssociation{{
				AssociationId: awssdk.String("iip-assoc-1234"),
				InstanceId:    awssdk.String("i-1234"),
				IamInstanceProfile: &ec2types.IamInstanceProfile{
					Arn: awssdk.String("arn:aws:iam::1:instance-profile/my-profile"),
				},
			}},
		}).
		On("DisassociateIamInstanceProfile", &ec2.DisassociateIamInstanceProfileOutput{
			IamInstanceProfileAssociation: &ec2types.IamInstanceProfileAssociation{
				AssociationId:      awssdk.String("iip-assoc-1234"),
				IamInstanceProfile: &ec2types.IamInstanceProfile{Id: awssdk.String("AIPA1234")},
			},
		})

	Template("detach instanceprofile instance=i-1234 name=my-profile").
		Mock(mock).
		ExpectCalls("DescribeIamInstanceProfileAssociations", "DisassociateIamInstanceProfile").
		Run(t)

	in := mock.InputFor("DisassociateIamInstanceProfile").(*ec2.DisassociateIamInstanceProfileInput)
	// The association has to be looked up; the instance id alone will not do.
	if got := awssdk.ToString(in.AssociationId); got != "iip-assoc-1234" {
		t.Errorf("AssociationId: got %q, want iip-assoc-1234", got)
	}
}

func TestDetachNetworkinterfaceByAttachment(t *testing.T) {
	mock := NewMock().On("DetachNetworkInterface", &ec2.DetachNetworkInterfaceOutput{})

	Template("detach networkinterface attachment=eni-attach-1234 force=true").
		Mock(mock).
		ExpectCalls("DetachNetworkInterface").
		Run(t)

	in := mock.InputFor("DetachNetworkInterface").(*ec2.DetachNetworkInterfaceInput)
	if got := awssdk.ToString(in.AttachmentId); got != "eni-attach-1234" {
		t.Errorf("AttachmentId: got %q, want eni-attach-1234", got)
	}
}

func TestRestartInstance(t *testing.T) {
	mock := NewMock().On("RebootInstances", &ec2.RebootInstancesOutput{})

	Template("restart instance ids=i-1234").
		Mock(mock).
		ExpectCalls("RebootInstances").
		Run(t)
}

func TestUpdateImageSharesWithAccounts(t *testing.T) {
	mock := NewMock().On("ModifyImageAttribute", &ec2.ModifyImageAttributeOutput{})

	Template("update image id=ami-1234 accounts=123456789012 operation=add").
		Mock(mock).
		ExpectCalls("ModifyImageAttribute").
		Run(t)

	in := mock.InputFor("ModifyImageAttribute").(*ec2.ModifyImageAttributeInput)
	if got := awssdk.ToString(in.ImageId); got != "ami-1234" {
		t.Errorf("ImageId: got %q, want ami-1234", got)
	}
}

func TestImportImageFromSnapshot(t *testing.T) {
	mock := NewMock().On("ImportImage", &ec2.ImportImageOutput{
		ImportTaskId: awssdk.String("import-ami-1234"),
	})

	Template("import image snapshot=snap-1234 architecture=x86_64 platform=Linux").
		Mock(mock).
		ExpectCalls("ImportImage").
		Run(t)

	in := mock.InputFor("ImportImage").(*ec2.ImportImageInput)
	if len(in.DiskContainers) != 1 {
		t.Fatalf("expected one disk container, got %#v", in.DiskContainers)
	}
	if got := awssdk.ToString(in.DiskContainers[0].SnapshotId); got != "snap-1234" {
		t.Errorf("SnapshotId: got %q, want snap-1234", got)
	}
}

func TestDeleteRecordById(t *testing.T) {
	mock := NewMock().On("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53types.ChangeInfo{Id: awssdk.String("C123")},
	})

	Template("delete record zone=/hostedzone/Z1 name=www.example.com type=A ttl=300 value=1.2.3.4").
		Mock(mock).
		ExpectCalls("ChangeResourceRecordSets").
		Run(t)

	in := mock.InputFor("ChangeResourceRecordSets").(*route53.ChangeResourceRecordSetsInput)
	if got := string(in.ChangeBatch.Changes[0].Action); got != "DELETE" {
		t.Errorf("Action: got %q, want DELETE", got)
	}
}

func TestDeleteContainertask(t *testing.T) {
	mock := NewMock().
		On("ListTaskDefinitions", &ecs.ListTaskDefinitionsOutput{
			TaskDefinitionArns: []string{"arn:aws:ecs:us-west-2:1:task-definition/my-task:1"},
		}).
		On("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionOutput{
			TaskDefinition: &ecstypes.TaskDefinition{
				Family:            awssdk.String("my-task"),
				TaskDefinitionArn: awssdk.String("arn:aws:ecs:us-west-2:1:task-definition/my-task:1"),
			},
		}).
		On("DeregisterTaskDefinition", &ecs.DeregisterTaskDefinitionOutput{
			TaskDefinition: &ecstypes.TaskDefinition{
				TaskDefinitionArn: awssdk.String("arn:aws:ecs:us-west-2:1:task-definition/my-task:1"),
			},
		})

	Template("delete containertask name=my-task").
		Mock(mock).
		ExpectCalls("DescribeTaskDefinition", "DeregisterTaskDefinition").
		Run(t)
}

package awsconv

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
)

func TestInitResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        interface{}
		expectedType string
		expectedID   string
	}{
		// EC2
		{
			name:         "ec2 Instance",
			input:        ec2types.Instance{InstanceId: awssdk.String("i-12345")},
			expectedType: cloud.Instance,
			expectedID:   "i-12345",
		},
		{
			name:         "ec2 Vpc",
			input:        ec2types.Vpc{VpcId: awssdk.String("vpc-abc")},
			expectedType: cloud.Vpc,
			expectedID:   "vpc-abc",
		},
		{
			name:         "ec2 Subnet",
			input:        ec2types.Subnet{SubnetId: awssdk.String("subnet-001")},
			expectedType: cloud.Subnet,
			expectedID:   "subnet-001",
		},
		{
			name:         "ec2 SecurityGroup",
			input:        ec2types.SecurityGroup{GroupId: awssdk.String("sg-999")},
			expectedType: cloud.SecurityGroup,
			expectedID:   "sg-999",
		},
		{
			name:         "ec2 KeyPairInfo",
			input:        ec2types.KeyPairInfo{KeyName: awssdk.String("my-key")},
			expectedType: cloud.Keypair,
			expectedID:   "my-key",
		},
		{
			name:         "ec2 Volume",
			input:        ec2types.Volume{VolumeId: awssdk.String("vol-123")},
			expectedType: cloud.Volume,
			expectedID:   "vol-123",
		},
		{
			name:         "ec2 Image",
			input:        ec2types.Image{ImageId: awssdk.String("ami-abc")},
			expectedType: cloud.Image,
			expectedID:   "ami-abc",
		},
		{
			name:         "ec2 ImportImageTask",
			input:        ec2types.ImportImageTask{ImportTaskId: awssdk.String("import-task-1")},
			expectedType: cloud.ImportImageTask,
			expectedID:   "import-task-1",
		},
		{
			name:         "ec2 InternetGateway",
			input:        ec2types.InternetGateway{InternetGatewayId: awssdk.String("igw-001")},
			expectedType: cloud.InternetGateway,
			expectedID:   "igw-001",
		},
		{
			name:         "ec2 NatGateway",
			input:        ec2types.NatGateway{NatGatewayId: awssdk.String("nat-001")},
			expectedType: cloud.NatGateway,
			expectedID:   "nat-001",
		},
		{
			name:         "ec2 RouteTable",
			input:        ec2types.RouteTable{RouteTableId: awssdk.String("rtb-001")},
			expectedType: cloud.RouteTable,
			expectedID:   "rtb-001",
		},
		{
			name:         "ec2 AvailabilityZone",
			input:        ec2types.AvailabilityZone{ZoneName: awssdk.String("us-east-1a")},
			expectedType: cloud.AvailabilityZone,
			expectedID:   "us-east-1a",
		},
		{
			name:         "ec2 Address with AllocationId",
			input:        ec2types.Address{AllocationId: awssdk.String("eipalloc-001"), PublicIp: awssdk.String("1.2.3.4")},
			expectedType: cloud.ElasticIP,
			expectedID:   "eipalloc-001",
		},
		{
			name:         "ec2 Address without AllocationId falls back to PublicIp",
			input:        ec2types.Address{PublicIp: awssdk.String("5.6.7.8")},
			expectedType: cloud.ElasticIP,
			expectedID:   "5.6.7.8",
		},
		{
			name:         "ec2 Snapshot",
			input:        ec2types.Snapshot{SnapshotId: awssdk.String("snap-001")},
			expectedType: cloud.Snapshot,
			expectedID:   "snap-001",
		},
		{
			name:         "ec2 NetworkInterface",
			input:        ec2types.NetworkInterface{NetworkInterfaceId: awssdk.String("eni-001")},
			expectedType: cloud.NetworkInterface,
			expectedID:   "eni-001",
		},
		// Loadbalancer
		{
			name:         "classic LoadBalancerDescription",
			input:        elbtypes.LoadBalancerDescription{LoadBalancerName: awssdk.String("my-classic-lb")},
			expectedType: cloud.ClassicLoadBalancer,
			expectedID:   "my-classic-lb",
		},
		{
			name:         "elbv2 LoadBalancer",
			input:        elbv2types.LoadBalancer{LoadBalancerArn: awssdk.String("arn:aws:elasticloadbalancing:us-east-1:123:loadbalancer/app/my-lb/abc")},
			expectedType: cloud.LoadBalancer,
			expectedID:   "arn:aws:elasticloadbalancing:us-east-1:123:loadbalancer/app/my-lb/abc",
		},
		{
			name:         "elbv2 TargetGroup",
			input:        elbv2types.TargetGroup{TargetGroupArn: awssdk.String("arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/tg/abc")},
			expectedType: cloud.TargetGroup,
			expectedID:   "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/tg/abc",
		},
		{
			name:         "elbv2 Listener",
			input:        elbv2types.Listener{ListenerArn: awssdk.String("arn:aws:elasticloadbalancing:us-east-1:123:listener/app/my-lb/abc/def")},
			expectedType: cloud.Listener,
			expectedID:   "arn:aws:elasticloadbalancing:us-east-1:123:listener/app/my-lb/abc/def",
		},
		// Database
		{
			name:         "rds DBInstance",
			input:        rdstypes.DBInstance{DBInstanceIdentifier: awssdk.String("mydb")},
			expectedType: cloud.Database,
			expectedID:   "mydb",
		},
		{
			name:         "rds DBSubnetGroup",
			input:        rdstypes.DBSubnetGroup{DBSubnetGroupArn: awssdk.String("arn:aws:rds:us-east-1:123:subgrp:mysubgrp")},
			expectedType: cloud.DbSubnetGroup,
			expectedID:   "arn:aws:rds:us-east-1:123:subgrp:mysubgrp",
		},
		// Autoscaling
		{
			name:         "autoscaling LaunchConfiguration",
			input:        autoscalingtypes.LaunchConfiguration{LaunchConfigurationARN: awssdk.String("arn:aws:autoscaling:us-east-1:123:launchConfiguration:lc-001")},
			expectedType: cloud.LaunchConfiguration,
			expectedID:   "arn:aws:autoscaling:us-east-1:123:launchConfiguration:lc-001",
		},
		{
			name:         "autoscaling AutoScalingGroup",
			input:        autoscalingtypes.AutoScalingGroup{AutoScalingGroupARN: awssdk.String("arn:aws:autoscaling:us-east-1:123:autoScalingGroup:asg-001")},
			expectedType: cloud.ScalingGroup,
			expectedID:   "arn:aws:autoscaling:us-east-1:123:autoScalingGroup:asg-001",
		},
		{
			name:         "autoscaling ScalingPolicy",
			input:        autoscalingtypes.ScalingPolicy{PolicyARN: awssdk.String("arn:aws:autoscaling:us-east-1:123:scalingPolicy:sp-001")},
			expectedType: cloud.ScalingPolicy,
			expectedID:   "arn:aws:autoscaling:us-east-1:123:scalingPolicy:sp-001",
		},
		// Container
		{
			name:         "ecr Repository",
			input:        ecrtypes.Repository{RepositoryArn: awssdk.String("arn:aws:ecr:us-east-1:123:repository/myrepo")},
			expectedType: cloud.Repository,
			expectedID:   "arn:aws:ecr:us-east-1:123:repository/myrepo",
		},
		{
			name:         "ecs Cluster",
			input:        ecstypes.Cluster{ClusterArn: awssdk.String("arn:aws:ecs:us-east-1:123:cluster/mycluster")},
			expectedType: cloud.ContainerCluster,
			expectedID:   "arn:aws:ecs:us-east-1:123:cluster/mycluster",
		},
		{
			name:         "ecs TaskDefinition",
			input:        ecstypes.TaskDefinition{TaskDefinitionArn: awssdk.String("arn:aws:ecs:us-east-1:123:task-definition/mytask:1")},
			expectedType: cloud.ContainerTask,
			expectedID:   "arn:aws:ecs:us-east-1:123:task-definition/mytask:1",
		},
		{
			name:         "ecs Container",
			input:        ecstypes.Container{ContainerArn: awssdk.String("arn:aws:ecs:us-east-1:123:container/abc")},
			expectedType: cloud.Container,
			expectedID:   "arn:aws:ecs:us-east-1:123:container/abc",
		},
		{
			name:         "ecs ContainerInstance",
			input:        ecstypes.ContainerInstance{ContainerInstanceArn: awssdk.String("arn:aws:ecs:us-east-1:123:container-instance/ci-001")},
			expectedType: cloud.ContainerInstance,
			expectedID:   "arn:aws:ecs:us-east-1:123:container-instance/ci-001",
		},
		// ACM
		{
			name:         "acm CertificateSummary",
			input:        acmtypes.CertificateSummary{CertificateArn: awssdk.String("arn:aws:acm:us-east-1:123:certificate/cert-001")},
			expectedType: cloud.Certificate,
			expectedID:   "arn:aws:acm:us-east-1:123:certificate/cert-001",
		},
		// IAM
		{
			name:         "iam User",
			input:        iamtypes.User{UserId: awssdk.String("AIDA12345")},
			expectedType: cloud.User,
			expectedID:   "AIDA12345",
		},
		{
			name:         "iam UserDetail",
			input:        iamtypes.UserDetail{UserId: awssdk.String("AIDA67890")},
			expectedType: cloud.User,
			expectedID:   "AIDA67890",
		},
		{
			name:         "iam RoleDetail",
			input:        iamtypes.RoleDetail{RoleId: awssdk.String("AROA12345")},
			expectedType: cloud.Role,
			expectedID:   "AROA12345",
		},
		{
			name:         "iam GroupDetail",
			input:        iamtypes.GroupDetail{GroupId: awssdk.String("AGPA12345")},
			expectedType: cloud.Group,
			expectedID:   "AGPA12345",
		},
		{
			name:         "iam Policy",
			input:        iamtypes.Policy{PolicyId: awssdk.String("ANPA12345")},
			expectedType: cloud.Policy,
			expectedID:   "ANPA12345",
		},
		{
			name:         "iam ManagedPolicyDetail",
			input:        iamtypes.ManagedPolicyDetail{PolicyId: awssdk.String("ANPA67890")},
			expectedType: cloud.Policy,
			expectedID:   "ANPA67890",
		},
		{
			name:         "iam AccessKeyMetadata",
			input:        iamtypes.AccessKeyMetadata{AccessKeyId: awssdk.String("AKIA12345")},
			expectedType: cloud.AccessKey,
			expectedID:   "AKIA12345",
		},
		{
			name:         "iam InstanceProfile",
			input:        iamtypes.InstanceProfile{InstanceProfileId: awssdk.String("AIPA12345")},
			expectedType: cloud.InstanceProfile,
			expectedID:   "AIPA12345",
		},
		{
			name:         "iam VirtualMFADevice",
			input:        iamtypes.VirtualMFADevice{SerialNumber: awssdk.String("arn:aws:iam::123:mfa/my-mfa")},
			expectedType: cloud.MFADevice,
			expectedID:   "arn:aws:iam::123:mfa/my-mfa",
		},
		// S3
		{
			name:         "s3 Bucket",
			input:        s3types.Bucket{Name: awssdk.String("my-bucket")},
			expectedType: cloud.Bucket,
			expectedID:   "my-bucket",
		},
		{
			name:         "s3 Object",
			input:        s3types.Object{Key: awssdk.String("path/to/object.txt")},
			expectedType: cloud.S3Object,
			expectedID:   "path/to/object.txt",
		},
		// SNS
		{
			name:         "sns Subscription",
			input:        snstypes.Subscription{Endpoint: awssdk.String("arn:aws:sqs:us-east-1:123:my-queue")},
			expectedType: cloud.Subscription,
			expectedID:   "arn:aws:sqs:us-east-1:123:my-queue",
		},
		{
			name:         "sns Topic",
			input:        snstypes.Topic{TopicArn: awssdk.String("arn:aws:sns:us-east-1:123:my-topic")},
			expectedType: cloud.Topic,
			expectedID:   "arn:aws:sns:us-east-1:123:my-topic",
		},
		// DNS
		{
			name:         "route53 HostedZone",
			input:        route53types.HostedZone{Id: awssdk.String("/hostedzone/Z12345")},
			expectedType: cloud.Zone,
			expectedID:   "/hostedzone/Z12345",
		},
		{
			name:         "route53 ResourceRecordSet uses HashFields for ID",
			input:        route53types.ResourceRecordSet{Name: awssdk.String("example.com."), Type: route53types.RRTypeA},
			expectedType: cloud.Record,
			expectedID:   HashFields("example.com.", string(route53types.RRTypeA)),
		},
		// Lambda
		{
			name:         "lambda FunctionConfiguration",
			input:        lambdatypes.FunctionConfiguration{FunctionArn: awssdk.String("arn:aws:lambda:us-east-1:123:function:my-func")},
			expectedType: cloud.Function,
			expectedID:   "arn:aws:lambda:us-east-1:123:function:my-func",
		},
		// Monitoring
		{
			name:         "cloudwatch Metric uses HashFields for ID",
			input:        cloudwatchtypes.Metric{Namespace: awssdk.String("AWS/EC2"), MetricName: awssdk.String("CPUUtilization")},
			expectedType: cloud.Metric,
			expectedID:   HashFields("AWS/EC2", "CPUUtilization"),
		},
		{
			name:         "cloudwatch MetricAlarm",
			input:        cloudwatchtypes.MetricAlarm{AlarmArn: awssdk.String("arn:aws:cloudwatch:us-east-1:123:alarm:my-alarm")},
			expectedType: cloud.Alarm,
			expectedID:   "arn:aws:cloudwatch:us-east-1:123:alarm:my-alarm",
		},
		// CDN
		{
			name:         "cloudfront DistributionSummary",
			input:        cloudfronttypes.DistributionSummary{Id: awssdk.String("E12345")},
			expectedType: cloud.Distribution,
			expectedID:   "E12345",
		},
		// Cloudformation
		{
			name:         "cloudformation Stack",
			input:        cloudformationtypes.Stack{StackId: awssdk.String("arn:aws:cloudformation:us-east-1:123:stack/my-stack/guid")},
			expectedType: cloud.Stack,
			expectedID:   "arn:aws:cloudformation:us-east-1:123:stack/my-stack/guid",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := InitResource(tc.input)
			if err != nil {
				t.Fatalf("InitResource returned error: %v", err)
			}
			if res == nil {
				t.Fatal("InitResource returned nil resource")
			}
			if got, want := res.Type(), tc.expectedType; got != want {
				t.Errorf("Type: got %q, want %q", got, want)
			}
			if got, want := res.Id(), tc.expectedID; got != want {
				t.Errorf("ID: got %q, want %q", got, want)
			}
		})
	}
}

func TestInitResourceUnknownType(t *testing.T) {
	t.Parallel()
	_, err := InitResource("not-an-aws-type")
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
}

func TestInitResourceNilFields(t *testing.T) {
	t.Parallel()
	// Instance with nil InstanceId should produce a resource with empty ID
	res, err := InitResource(ec2types.Instance{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Id(); got != "" {
		t.Errorf("expected empty ID for nil field, got %q", got)
	}
	if got, want := res.Type(), cloud.Instance; got != want {
		t.Errorf("Type: got %q, want %q", got, want)
	}
}

func TestHashFields(t *testing.T) {
	t.Parallel()

	t.Run("deterministic", func(t *testing.T) {
		t.Parallel()
		h1 := HashFields("hello", "world")
		h2 := HashFields("hello", "world")
		if h1 != h2 {
			t.Errorf("same inputs produced different hashes: %q vs %q", h1, h2)
		}
	})

	t.Run("prefix", func(t *testing.T) {
		t.Parallel()
		h := HashFields("test")
		if len(h) < 5 || h[:5] != "awls-" {
			t.Errorf("hash should start with 'awls-', got %q", h)
		}
	})

	t.Run("different inputs produce different hashes", func(t *testing.T) {
		t.Parallel()
		h1 := HashFields("a", "b")
		h2 := HashFields("c", "d")
		if h1 == h2 {
			t.Errorf("different inputs produced same hash: %q", h1)
		}
	})

	t.Run("order matters", func(t *testing.T) {
		t.Parallel()
		h1 := HashFields("first", "second")
		h2 := HashFields("second", "first")
		if h1 == h2 {
			t.Errorf("different order produced same hash: %q", h1)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		h := HashFields()
		if len(h) < 5 || h[:5] != "awls-" {
			t.Errorf("hash of empty input should still have prefix, got %q", h)
		}
	})

	t.Run("single vs multiple fields", func(t *testing.T) {
		t.Parallel()
		// HashFields("ab") and HashFields("a","b") should be the same since
		// it just concatenates with fmt.Sprint
		h1 := HashFields("ab")
		h2 := HashFields("a", "b")
		if h1 != h2 {
			t.Errorf("concatenated single field vs two fields differ: %q vs %q", h1, h2)
		}
	})

	t.Run("non-string fields", func(t *testing.T) {
		t.Parallel()
		h1 := HashFields(42, true)
		h2 := HashFields(42, true)
		if h1 != h2 {
			t.Errorf("same non-string inputs produced different hashes: %q vs %q", h1, h2)
		}
		h3 := HashFields(42, false)
		if h1 == h3 {
			t.Errorf("different non-string inputs produced same hash: %q", h1)
		}
	})
}

func TestNewResource(t *testing.T) {
	t.Parallel()

	t.Run("SecurityGroup with basic fields", func(t *testing.T) {
		t.Parallel()
		sg := ec2types.SecurityGroup{
			GroupId:     awssdk.String("sg-123"),
			GroupName:   awssdk.String("my-sg"),
			Description: awssdk.String("My security group"),
			OwnerId:     awssdk.String("123456789012"),
			VpcId:       awssdk.String("vpc-abc"),
		}

		res, err := NewResource(sg)
		if err != nil {
			t.Fatalf("NewResource returned error: %v", err)
		}
		if got, want := res.Type(), cloud.SecurityGroup; got != want {
			t.Errorf("Type: got %q, want %q", got, want)
		}
		if got, want := res.Id(), "sg-123"; got != want {
			t.Errorf("ID: got %q, want %q", got, want)
		}

		props := res.Properties()

		if got, want := props[properties.ID], "sg-123"; got != want {
			t.Errorf("properties.ID: got %v, want %v", got, want)
		}
		if got, want := props[properties.Name], "my-sg"; got != want {
			t.Errorf("properties.Name: got %v, want %v", got, want)
		}
		if got, want := props[properties.Description], "My security group"; got != want {
			t.Errorf("properties.Description: got %v, want %v", got, want)
		}
		if got, want := props[properties.Owner], "123456789012"; got != want {
			t.Errorf("properties.Owner: got %v, want %v", got, want)
		}
		if got, want := props[properties.Vpc], "vpc-abc"; got != want {
			t.Errorf("properties.Vpc: got %v, want %v", got, want)
		}
	})

	t.Run("Vpc with fields and tags", func(t *testing.T) {
		t.Parallel()
		vpc := ec2types.Vpc{
			VpcId:     awssdk.String("vpc-999"),
			CidrBlock: awssdk.String("10.0.0.0/16"),
			IsDefault: awssdk.Bool(true),
			State:     ec2types.VpcStateAvailable,
			Tags: []ec2types.Tag{
				{Key: awssdk.String("Name"), Value: awssdk.String("my-vpc")},
				{Key: awssdk.String("env"), Value: awssdk.String("prod")},
			},
		}

		res, err := NewResource(vpc)
		if err != nil {
			t.Fatalf("NewResource returned error: %v", err)
		}
		if got, want := res.Type(), cloud.Vpc; got != want {
			t.Errorf("Type: got %q, want %q", got, want)
		}
		if got, want := res.Id(), "vpc-999"; got != want {
			t.Errorf("ID: got %q, want %q", got, want)
		}

		props := res.Properties()
		if got, want := props[properties.Name], "my-vpc"; got != want {
			t.Errorf("properties.Name: got %v, want %v", got, want)
		}
		if got, want := props[properties.CIDR], "10.0.0.0/16"; got != want {
			t.Errorf("properties.CIDR: got %v, want %v", got, want)
		}
		if got, want := props[properties.Default], true; got != want {
			t.Errorf("properties.Default: got %v, want %v", got, want)
		}
		if got, want := props[properties.State], "available"; got != want {
			t.Errorf("properties.State: got %v, want %v", got, want)
		}
	})

	t.Run("Subnet with fields", func(t *testing.T) {
		t.Parallel()
		subnet := ec2types.Subnet{
			SubnetId:         awssdk.String("subnet-001"),
			VpcId:            awssdk.String("vpc-abc"),
			CidrBlock:        awssdk.String("10.0.1.0/24"),
			AvailabilityZone: awssdk.String("us-east-1a"),
			State:            ec2types.SubnetStateAvailable,
		}

		res, err := NewResource(subnet)
		if err != nil {
			t.Fatalf("NewResource returned error: %v", err)
		}
		if got, want := res.Id(), "subnet-001"; got != want {
			t.Errorf("ID: got %q, want %q", got, want)
		}

		props := res.Properties()
		if got, want := props[properties.Vpc], "vpc-abc"; got != want {
			t.Errorf("properties.Vpc: got %v, want %v", got, want)
		}
		if got, want := props[properties.CIDR], "10.0.1.0/24"; got != want {
			t.Errorf("properties.CIDR: got %v, want %v", got, want)
		}
		if got, want := props[properties.AvailabilityZone], "us-east-1a"; got != want {
			t.Errorf("properties.AvailabilityZone: got %v, want %v", got, want)
		}
	})

	t.Run("S3 Bucket", func(t *testing.T) {
		t.Parallel()
		bucket := s3types.Bucket{
			Name: awssdk.String("my-bucket"),
		}
		res, err := NewResource(bucket)
		if err != nil {
			t.Fatalf("NewResource returned error: %v", err)
		}
		if got, want := res.Type(), cloud.Bucket; got != want {
			t.Errorf("Type: got %q, want %q", got, want)
		}
		if got, want := res.Id(), "my-bucket"; got != want {
			t.Errorf("ID: got %q, want %q", got, want)
		}
		if got, want := res.Properties()[properties.ID], "my-bucket"; got != want {
			t.Errorf("properties.ID: got %v, want %v", got, want)
		}
	})

	t.Run("unknown type returns error", func(t *testing.T) {
		t.Parallel()
		_, err := NewResource(42)
		if err == nil {
			t.Fatal("expected error for unknown type, got nil")
		}
	})

	t.Run("IAM User with name", func(t *testing.T) {
		t.Parallel()
		user := iamtypes.User{
			UserId:   awssdk.String("AIDA12345"),
			UserName: awssdk.String("admin"),
			Arn:      awssdk.String("arn:aws:iam::123:user/admin"),
			Path:     awssdk.String("/"),
		}
		res, err := NewResource(user)
		if err != nil {
			t.Fatalf("NewResource returned error: %v", err)
		}
		if got, want := res.Type(), cloud.User; got != want {
			t.Errorf("Type: got %q, want %q", got, want)
		}
		if got, want := res.Id(), "AIDA12345"; got != want {
			t.Errorf("ID: got %q, want %q", got, want)
		}
		props := res.Properties()
		if got, want := props[properties.Name], "admin"; got != want {
			t.Errorf("properties.Name: got %v, want %v", got, want)
		}
		if got, want := props[properties.Arn], "arn:aws:iam::123:user/admin"; got != want {
			t.Errorf("properties.Arn: got %v, want %v", got, want)
		}
		if got, want := props[properties.Path], "/"; got != want {
			t.Errorf("properties.Path: got %v, want %v", got, want)
		}
	})
}

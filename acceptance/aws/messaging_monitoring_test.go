package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Messaging, storage, monitoring and the remaining services.

func TestDeleteTopicAndSubscription(t *testing.T) {
	topic := NewMock().On("DeleteTopic", &sns.DeleteTopicOutput{})
	Template("delete topic id=arn:aws:sns:us-west-2:1:alerts").
		Mock(topic).
		ExpectCalls("DeleteTopic").
		Run(t)

	sub := NewMock().On("Unsubscribe", &sns.UnsubscribeOutput{})
	Template("delete subscription id=arn:aws:sns:us-west-2:1:alerts:abcd").
		Mock(sub).
		ExpectCalls("Unsubscribe").
		Run(t)

	in := sub.InputFor("Unsubscribe").(*sns.UnsubscribeInput)
	if got := awssdk.ToString(in.SubscriptionArn); got != "arn:aws:sns:us-west-2:1:alerts:abcd" {
		t.Errorf("SubscriptionArn: got %q", got)
	}
}

func TestCreateSubscription(t *testing.T) {
	mock := NewMock().On("Subscribe", &sns.SubscribeOutput{
		SubscriptionArn: awssdk.String("arn:aws:sns:us-west-2:1:alerts:abcd"),
	})

	Template("create subscription topic=arn:aws:sns:us-west-2:1:alerts endpoint=ops@example.com protocol=email").
		Mock(mock).
		ExpectCalls("Subscribe").
		ExpectCommandResult("arn:aws:sns:us-west-2:1:alerts:abcd").
		Run(t)

	in := mock.InputFor("Subscribe").(*sns.SubscribeInput)
	if got := awssdk.ToString(in.Protocol); got != "email" {
		t.Errorf("Protocol: got %q, want email", got)
	}
	if got := awssdk.ToString(in.Endpoint); got != "ops@example.com" {
		t.Errorf("Endpoint: got %q", got)
	}
}

// delete queue takes the queue URL, not its name or ARN.
func TestDeleteQueue(t *testing.T) {
	url := "https://sqs.us-west-2.amazonaws.com/123456789012/my-queue"
	mock := NewMock().On("DeleteQueue", &sqs.DeleteQueueOutput{})

	Template("delete queue url=" + url).
		Mock(mock).
		ExpectCalls("DeleteQueue").
		Run(t)

	in := mock.InputFor("DeleteQueue").(*sqs.DeleteQueueInput)
	if got := awssdk.ToString(in.QueueUrl); got != url {
		t.Errorf("QueueUrl: got %q, want %q", got, url)
	}
}

func TestDeleteBucket(t *testing.T) {
	mock := NewMock().On("DeleteBucket", &s3.DeleteBucketOutput{})

	Template("delete bucket name=my-bucket").
		Mock(mock).
		ExpectCalls("DeleteBucket").
		Run(t)
}

func TestUpdateBucketAcl(t *testing.T) {
	mock := NewMock().On("PutBucketAcl", &s3.PutBucketAclOutput{})

	Template("update bucket name=my-bucket acl=public-read").
		Mock(mock).
		ExpectCalls("PutBucketAcl").
		Run(t)

	in := mock.InputFor("PutBucketAcl").(*s3.PutBucketAclInput)
	if got := string(in.ACL); got != "public-read" {
		t.Errorf("ACL: got %q, want public-read", got)
	}
}

func TestDeleteS3object(t *testing.T) {
	mock := NewMock().On("DeleteObject", &s3.DeleteObjectOutput{})

	Template("delete s3object bucket=my-bucket name=path/to/file.txt").
		Mock(mock).
		ExpectCalls("DeleteObject").
		Run(t)

	in := mock.InputFor("DeleteObject").(*s3.DeleteObjectInput)
	if got := awssdk.ToString(in.Key); got != "path/to/file.txt" {
		t.Errorf("Key: got %q, want path/to/file.txt", got)
	}
}

func TestUpdateS3objectAcl(t *testing.T) {
	mock := NewMock().On("PutObjectAcl", &s3.PutObjectAclOutput{})

	Template("update s3object bucket=my-bucket name=file.txt acl=private").
		Mock(mock).
		ExpectCalls("PutObjectAcl").
		Run(t)
}

// create alarm carries eight required params, several of them numeric and adjacent in
// meaning, so this is the highest-risk mapping in the set.
func TestCreateAlarm(t *testing.T) {
	mock := NewMock().On("PutMetricAlarm", &cloudwatch.PutMetricAlarmOutput{})

	Template("create alarm name=high-cpu metric=CPUUtilization namespace=AWS/EC2 operator=GreaterThanThreshold threshold=80 period=300 evaluation-periods=2 statistic-function=Average dimensions=InstanceId:i-1234").
		Mock(mock).
		ExpectCalls("PutMetricAlarm").
		ExpectCommandResult("high-cpu").
		ExpectRevert("delete alarm name=high-cpu").
		Run(t)

	in := mock.InputFor("PutMetricAlarm").(*cloudwatch.PutMetricAlarmInput)
	if got := awssdk.ToFloat64(in.Threshold); got != 80 {
		t.Errorf("Threshold: got %v, want 80", got)
	}
	// period and evaluation-periods are both integers and easily crossed.
	if got := awssdk.ToInt32(in.Period); got != 300 {
		t.Errorf("Period: got %d, want 300", got)
	}
	if got := awssdk.ToInt32(in.EvaluationPeriods); got != 2 {
		t.Errorf("EvaluationPeriods: got %d, want 2", got)
	}
	if got := string(in.Statistic); got != "Average" {
		t.Errorf("Statistic: got %q, want Average", got)
	}
	if len(in.Dimensions) != 1 {
		t.Fatalf("expected one dimension, got %#v", in.Dimensions)
	}
	if got := awssdk.ToString(in.Dimensions[0].Name); got != "InstanceId" {
		t.Errorf("dimension name: got %q, want InstanceId", got)
	}
}

func TestStartStopDeleteAlarm(t *testing.T) {
	start := NewMock().On("EnableAlarmActions", &cloudwatch.EnableAlarmActionsOutput{})
	Template("start alarm names=high-cpu").Mock(start).ExpectCalls("EnableAlarmActions").Run(t)

	stop := NewMock().On("DisableAlarmActions", &cloudwatch.DisableAlarmActionsOutput{})
	Template("stop alarm names=high-cpu").Mock(stop).ExpectCalls("DisableAlarmActions").Run(t)

	del := NewMock().On("DeleteAlarms", &cloudwatch.DeleteAlarmsOutput{})
	Template("delete alarm name=high-cpu").Mock(del).ExpectCalls("DeleteAlarms").Run(t)

	in := del.InputFor("DeleteAlarms").(*cloudwatch.DeleteAlarmsInput)
	if len(in.AlarmNames) != 1 || in.AlarmNames[0] != "high-cpu" {
		t.Errorf("AlarmNames: got %#v", in.AlarmNames)
	}
}

// attach alarm adds an action to the existing set rather than replacing it, so the
// existing alarm has to be read first.
func TestAttachAlarmPreservesExistingActions(t *testing.T) {
	mock := NewMock().
		On("DescribeAlarms", &cloudwatch.DescribeAlarmsOutput{
			MetricAlarms: []cloudwatchtypes.MetricAlarm{{
				AlarmName:    awssdk.String("high-cpu"),
				AlarmArn:     awssdk.String("arn:aws:cloudwatch:us-west-2:1:alarm:high-cpu"),
				AlarmActions: []string{"arn:aws:sns:us-west-2:1:existing"},
			}},
		}).
		On("PutMetricAlarm", &cloudwatch.PutMetricAlarmOutput{})

	Template("attach alarm name=high-cpu action-arn=arn:aws:sns:us-west-2:1:new").
		Mock(mock).
		ExpectCalls("DescribeAlarms", "PutMetricAlarm").
		Run(t)

	in := mock.InputFor("PutMetricAlarm").(*cloudwatch.PutMetricAlarmInput)
	if len(in.AlarmActions) != 2 {
		t.Errorf("expected the existing action to be kept alongside the new one, got %#v", in.AlarmActions)
	}
}

func TestCreateStack(t *testing.T) {
	mock := NewMock().On("CreateStack", &cloudformation.CreateStackOutput{
		StackId: awssdk.String("arn:aws:cloudformation:us-west-2:1:stack/my-stack/abcd"),
	})

	Template("create stack name=my-stack template-file=/dev/null parameters=InstanceType:t3.micro on-failure=ROLLBACK").
		Mock(mock).
		ExpectCalls("CreateStack").
		Run(t)

	in := mock.InputFor("CreateStack").(*cloudformation.CreateStackInput)
	if got := awssdk.ToString(in.StackName); got != "my-stack" {
		t.Errorf("StackName: got %q, want my-stack", got)
	}
	if len(in.Parameters) != 1 {
		t.Fatalf("expected one parameter, got %#v", in.Parameters)
	}
	if got := awssdk.ToString(in.Parameters[0].ParameterKey); got != "InstanceType" {
		t.Errorf("parameter key: got %q, want InstanceType", got)
	}
}

func TestUpdateAndDeleteStack(t *testing.T) {
	update := NewMock().On("UpdateStack", &cloudformation.UpdateStackOutput{
		StackId: awssdk.String("arn:aws:cloudformation:us-west-2:1:stack/my-stack/abcd"),
	})
	Template("update stack name=my-stack template-file=/dev/null").
		Mock(update).
		ExpectCalls("UpdateStack").
		Run(t)

	del := NewMock().On("DeleteStack", &cloudformation.DeleteStackOutput{})
	Template("delete stack name=my-stack retain-resources=my-bucket").
		Mock(del).
		ExpectCalls("DeleteStack").
		Run(t)

	in := del.InputFor("DeleteStack").(*cloudformation.DeleteStackInput)
	if len(in.RetainResources) != 1 {
		t.Errorf("RetainResources: got %#v, want one entry", in.RetainResources)
	}
}

func TestDeleteZone(t *testing.T) {
	mock := NewMock().On("DeleteHostedZone", &route53.DeleteHostedZoneOutput{
		ChangeInfo: &route53types.ChangeInfo{Id: awssdk.String("C123")},
	})

	Template("delete zone id=/hostedzone/Z3P5QSUBK4POTI").
		Mock(mock).
		ExpectCalls("DeleteHostedZone").
		Run(t)
}

func TestUpdateRecord(t *testing.T) {
	mock := NewMock().On("ChangeResourceRecordSets", &route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53types.ChangeInfo{Id: awssdk.String("C123")},
	})

	Template("update record zone=/hostedzone/Z1 name=www.example.com type=A ttl=60 value=9.9.9.9").
		Mock(mock).
		ExpectCalls("ChangeResourceRecordSets").
		Run(t)

	in := mock.InputFor("ChangeResourceRecordSets").(*route53.ChangeResourceRecordSetsInput)
	rrs := in.ChangeBatch.Changes[0].ResourceRecordSet
	if got := awssdk.ToInt64(rrs.TTL); got != 60 {
		t.Errorf("TTL: got %d, want 60", got)
	}
	// A single `value` must still produce one resource record.
	if len(rrs.ResourceRecords) != 1 {
		t.Errorf("expected one resource record, got %#v", rrs.ResourceRecords)
	}
}

func TestDeleteFunction(t *testing.T) {
	mock := NewMock().On("DeleteFunction", &lambda.DeleteFunctionOutput{})

	Template("delete function id=my-function").
		Mock(mock).
		ExpectCalls("DeleteFunction").
		Run(t)

	in := mock.InputFor("DeleteFunction").(*lambda.DeleteFunctionInput)
	if got := awssdk.ToString(in.FunctionName); got != "my-function" {
		t.Errorf("FunctionName: got %q, want my-function", got)
	}
}

func TestDeleteRepository(t *testing.T) {
	mock := NewMock().On("DeleteRepository", &ecr.DeleteRepositoryOutput{})

	Template("delete repository name=my-repo force=true").
		Mock(mock).
		ExpectCalls("DeleteRepository").
		Run(t)

	in := mock.InputFor("DeleteRepository").(*ecr.DeleteRepositoryInput)
	if !in.Force {
		t.Error("expected Force to be set")
	}
}

func TestDeleteCertificate(t *testing.T) {
	mock := NewMock().On("DeleteCertificate", &acm.DeleteCertificateOutput{})

	Template("delete certificate arn=arn:aws:acm:us-west-2:1:certificate/abcd").
		Mock(mock).
		ExpectCalls("DeleteCertificate").
		Run(t)
}

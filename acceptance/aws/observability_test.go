package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func TestCreateTrail(t *testing.T) {
	mock := NewMock().On("CreateTrail", &cloudtrail.CreateTrailOutput{
		Name: awssdk.String("org-audit"),
	})

	Template("create trail name=org-audit bucket=audit-logs multiregion=true log-validation=true").
		Mock(mock).
		ExpectCalls("CreateTrail").
		ExpectCommandResult("org-audit").
		Run(t)

	in := mock.InputFor("CreateTrail").(*cloudtrail.CreateTrailInput)
	if got := awssdk.ToString(in.S3BucketName); got != "audit-logs" {
		t.Errorf("S3BucketName: got %q, want audit-logs", got)
	}
	if !awssdk.ToBool(in.IsMultiRegionTrail) {
		t.Error("expected IsMultiRegionTrail to be true")
	}
	if !awssdk.ToBool(in.EnableLogFileValidation) {
		t.Error("expected EnableLogFileValidation to be true")
	}
}

// A trail does not log until StartLogging runs, so the lifecycle is worth covering.
func TestStartAndStopTrail(t *testing.T) {
	start := NewMock().On("StartLogging", &cloudtrail.StartLoggingOutput{})
	Template("start trail name=org-audit").Mock(start).ExpectCalls("StartLogging").Run(t)

	stop := NewMock().On("StopLogging", &cloudtrail.StopLoggingOutput{})
	Template("stop trail name=org-audit").Mock(stop).ExpectCalls("StopLogging").Run(t)

	in := stop.InputFor("StopLogging").(*cloudtrail.StopLoggingInput)
	if got := awssdk.ToString(in.Name); got != "org-audit" {
		t.Errorf("Name: got %q, want org-audit", got)
	}
}

func TestCreateLoggroup(t *testing.T) {
	mock := NewMock().On("CreateLogGroup", &cloudwatchlogs.CreateLogGroupOutput{})

	Template("create loggroup name=/aws/lambda/handler").
		Mock(mock).
		ExpectCalls("CreateLogGroup").
		ExpectCommandResult("/aws/lambda/handler").
		Run(t)

	in := mock.InputFor("CreateLogGroup").(*cloudwatchlogs.CreateLogGroupInput)
	if got := awssdk.ToString(in.LogGroupName); got != "/aws/lambda/handler" {
		t.Errorf("LogGroupName: got %q", got)
	}
}

func TestUpdateLoggroupRetention(t *testing.T) {
	mock := NewMock().On("PutRetentionPolicy", &cloudwatchlogs.PutRetentionPolicyOutput{})

	Template("update loggroup name=/aws/lambda/handler retention=30").
		Mock(mock).
		ExpectCalls("PutRetentionPolicy").
		Run(t)

	in := mock.InputFor("PutRetentionPolicy").(*cloudwatchlogs.PutRetentionPolicyInput)
	if got := awssdk.ToInt32(in.RetentionInDays); got != 30 {
		t.Errorf("RetentionInDays: got %d, want 30", got)
	}
}

// CloudWatch Logs accepts only a fixed set of retention values, so a value outside
// it must be refused before the call rather than producing an opaque API error.
func TestUpdateLoggroupRejectsInvalidRetention(t *testing.T) {
	mock := NewMock().On("PutRetentionPolicy", &cloudwatchlogs.PutRetentionPolicyOutput{})

	if err := Template("update loggroup name=/aws/lambda/handler retention=45").
		Mock(mock).RunExpectingError(t); err == nil {
		t.Error("expected retention=45 to be refused")
	}
	if calls := mock.Calls()["PutRetentionPolicy"]; calls != 0 {
		t.Errorf("expected no API call, got %d", calls)
	}
}

func TestCreateApigateway(t *testing.T) {
	mock := NewMock().On("CreateApi", &apigatewayv2.CreateApiOutput{
		ApiId: awssdk.String("abc123"),
	})

	Template("create apigateway name=my-api protocol=HTTP").
		Mock(mock).
		ExpectCalls("CreateApi").
		ExpectCommandResult("abc123").
		Run(t)

	in := mock.InputFor("CreateApi").(*apigatewayv2.CreateApiInput)
	if got := string(in.ProtocolType); got != "HTTP" {
		t.Errorf("ProtocolType: got %q, want HTTP", got)
	}
}

func TestCreateApigatewayRejectsBadProtocol(t *testing.T) {
	mock := NewMock().On("CreateApi", &apigatewayv2.CreateApiOutput{ApiId: awssdk.String("x")})

	if err := Template("create apigateway name=my-api protocol=GRPC").
		Mock(mock).RunExpectingError(t); err == nil {
		t.Error("expected an invalid protocol to be refused")
	}
}

func TestCreateApigatewayroute(t *testing.T) {
	mock := NewMock().On("CreateRoute", &apigatewayv2.CreateRouteOutput{
		RouteId: awssdk.String("route1"),
	})

	Template(`create apigatewayroute api=abc123 route-key="GET /items" target=integrations/xyz`).
		Mock(mock).
		ExpectCalls("CreateRoute").
		ExpectCommandResult("route1").
		Run(t)

	in := mock.InputFor("CreateRoute").(*apigatewayv2.CreateRouteInput)
	if got := awssdk.ToString(in.RouteKey); got != "GET /items" {
		t.Errorf("RouteKey: got %q, want \"GET /items\"", got)
	}
	if got := awssdk.ToString(in.ApiId); got != "abc123" {
		t.Errorf("ApiId: got %q, want abc123", got)
	}
}

func TestCreateApigatewaystage(t *testing.T) {
	mock := NewMock().On("CreateStage", &apigatewayv2.CreateStageOutput{
		StageName: awssdk.String("prod"),
	})

	Template("create apigatewaystage api=abc123 name=prod autodeploy=true").
		Mock(mock).
		ExpectCalls("CreateStage").
		ExpectCommandResult("prod").
		Run(t)

	in := mock.InputFor("CreateStage").(*apigatewayv2.CreateStageInput)
	if !awssdk.ToBool(in.AutoDeploy) {
		t.Error("expected AutoDeploy to be true")
	}
}

func TestDeleteApigatewaystageNeedsApi(t *testing.T) {
	mock := NewMock().On("DeleteStage", &apigatewayv2.DeleteStageOutput{})

	Template("delete apigatewaystage api=abc123 name=prod").
		Mock(mock).
		ExpectCalls("DeleteStage").
		Run(t)

	in := mock.InputFor("DeleteStage").(*apigatewayv2.DeleteStageInput)
	if got := awssdk.ToString(in.ApiId); got != "abc123" {
		t.Errorf("ApiId: got %q, want abc123", got)
	}
}

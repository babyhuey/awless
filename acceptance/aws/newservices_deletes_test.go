package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// The last commands: deletes for the services added in this fork, plus ECR
// authentication.

func TestDeleteApigateway(t *testing.T) {
	mock := NewMock().On("DeleteApi", &apigatewayv2.DeleteApiOutput{})

	Template("delete apigateway id=abc123").
		Mock(mock).
		ExpectCalls("DeleteApi").
		Run(t)

	in := mock.InputFor("DeleteApi").(*apigatewayv2.DeleteApiInput)
	if got := awssdk.ToString(in.ApiId); got != "abc123" {
		t.Errorf("ApiId: got %q, want abc123", got)
	}
}

// A route is scoped to its API, so both ids have to reach the call.
func TestDeleteApigatewayroute(t *testing.T) {
	mock := NewMock().On("DeleteRoute", &apigatewayv2.DeleteRouteOutput{})

	Template("delete apigatewayroute api=abc123 id=route1").
		Mock(mock).
		ExpectCalls("DeleteRoute").
		Run(t)

	in := mock.InputFor("DeleteRoute").(*apigatewayv2.DeleteRouteInput)
	if got := awssdk.ToString(in.ApiId); got != "abc123" {
		t.Errorf("ApiId: got %q, want abc123", got)
	}
	if got := awssdk.ToString(in.RouteId); got != "route1" {
		t.Errorf("RouteId: got %q, want route1", got)
	}
}

func TestDeleteDynamodbtable(t *testing.T) {
	mock := NewMock().On("DeleteTable", &dynamodb.DeleteTableOutput{
		TableDescription: &dynamodbtypes.TableDescription{TableName: awssdk.String("users")},
	})

	Template("delete dynamodbtable name=users").
		Mock(mock).
		ExpectCalls("DeleteTable").
		Run(t)

	in := mock.InputFor("DeleteTable").(*dynamodb.DeleteTableInput)
	if got := awssdk.ToString(in.TableName); got != "users" {
		t.Errorf("TableName: got %q, want users", got)
	}
}

func TestDeleteEkscluster(t *testing.T) {
	mock := NewMock().On("DeleteCluster", &eks.DeleteClusterOutput{
		Cluster: &ekstypes.Cluster{Name: awssdk.String("prod")},
	})

	Template("delete ekscluster name=prod").
		Mock(mock).
		ExpectCalls("DeleteCluster").
		Run(t)

	in := mock.InputFor("DeleteCluster").(*eks.DeleteClusterInput)
	if got := awssdk.ToString(in.Name); got != "prod" {
		t.Errorf("Name: got %q, want prod", got)
	}
}

func TestDeleteLoggroup(t *testing.T) {
	mock := NewMock().On("DeleteLogGroup", &cloudwatchlogs.DeleteLogGroupOutput{})

	Template("delete loggroup name=/aws/lambda/handler").
		Mock(mock).
		ExpectCalls("DeleteLogGroup").
		Run(t)

	in := mock.InputFor("DeleteLogGroup").(*cloudwatchlogs.DeleteLogGroupInput)
	if got := awssdk.ToString(in.LogGroupName); got != "/aws/lambda/handler" {
		t.Errorf("LogGroupName: got %q", got)
	}
}

func TestDeleteTrail(t *testing.T) {
	mock := NewMock().On("DeleteTrail", &cloudtrail.DeleteTrailOutput{})

	Template("delete trail name=org-audit").
		Mock(mock).
		ExpectCalls("DeleteTrail").
		Run(t)

	in := mock.InputFor("DeleteTrail").(*cloudtrail.DeleteTrailInput)
	if got := awssdk.ToString(in.Name); got != "org-audit" {
		t.Errorf("Name: got %q, want org-audit", got)
	}
}

// update ssmparameter is PutParameter with Overwrite; without it the API refuses to
// replace an existing value.
func TestUpdateSsmparameter(t *testing.T) {
	mock := NewMock().On("PutParameter", &ssm.PutParameterOutput{Version: 2})

	Template("update ssmparameter name=/app/db/host value=db2.internal overwrite=true").
		Mock(mock).
		ExpectCalls("PutParameter").
		ExpectCommandResult("/app/db/host").
		Run(t)

	in := mock.InputFor("PutParameter").(*ssm.PutParameterInput)
	if !awssdk.ToBool(in.Overwrite) {
		t.Error("expected Overwrite, without which the update is refused")
	}
	if got := awssdk.ToString(in.Value); got != "db2.internal" {
		t.Errorf("Value: got %q, want db2.internal", got)
	}
}

// authenticate registry fetches a token for `docker login`. no-docker-login stops it
// shelling out, which a test must not do.
func TestAuthenticateRegistry(t *testing.T) {
	mock := NewMock().On("GetAuthorizationToken", &ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []ecrtypes.AuthorizationData{{
			// base64 of "AWS:password"
			AuthorizationToken: awssdk.String("QVdTOnBhc3N3b3Jk"),
			ProxyEndpoint:      awssdk.String("https://1.dkr.ecr.us-west-2.amazonaws.com"),
		}},
	})

	Template("authenticate registry no-docker-login=true").
		Mock(mock).
		ExpectCalls("GetAuthorizationToken").
		Run(t)
}

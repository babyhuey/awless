package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
)

func TestCreateDeployapplication(t *testing.T) {
	mock := NewMock().On("CreateApplication", &codedeploy.CreateApplicationOutput{
		ApplicationId: awssdk.String("app-1234"),
	})

	Template("create deployapplication name=web-api platform=Server").
		Mock(mock).
		ExpectCalls("CreateApplication").
		// The id comes back but every other CodeDeploy call takes the name, so that is
		// what the result has to be.
		ExpectCommandResult("web-api").
		ExpectRevert("delete deployapplication name=web-api").
		Run(t)

	in := mock.InputFor("CreateApplication").(*codedeploy.CreateApplicationInput)
	if got := awssdk.ToString(in.ApplicationName); got != "web-api" {
		t.Errorf("ApplicationName: got %q, want web-api", got)
	}
	if got := string(in.ComputePlatform); got != "Server" {
		t.Errorf("ComputePlatform: got %q, want Server", got)
	}
}

func TestDeleteDeployapplication(t *testing.T) {
	mock := NewMock().On("DeleteApplication", &codedeploy.DeleteApplicationOutput{})

	Template("delete deployapplication name=web-api").
		Mock(mock).
		ExpectCalls("DeleteApplication").
		Run(t)

	in := mock.InputFor("DeleteApplication").(*codedeploy.DeleteApplicationInput)
	if got := awssdk.ToString(in.ApplicationName); got != "web-api" {
		t.Errorf("ApplicationName: got %q, want web-api", got)
	}
}

// A deployment group name is only unique within its application, so the revert has to carry
// both — which the revert contract test would otherwise catch as an unexpected param.
func TestCreateDeploymentgroup(t *testing.T) {
	mock := NewMock().On("CreateDeploymentGroup", &codedeploy.CreateDeploymentGroupOutput{
		DeploymentGroupId: awssdk.String("dgp-1234"),
	})

	Template("create deploymentgroup name=prod application=web-api " +
		"role=arn:aws:iam::1:role/CodeDeploy config=CodeDeployDefault.OneAtATime scalinggroups=web-asg").
		Mock(mock).
		ExpectCalls("CreateDeploymentGroup").
		ExpectCommandResult("prod").
		ExpectRevert("delete deploymentgroup application=web-api name=prod").
		Run(t)

	in := mock.InputFor("CreateDeploymentGroup").(*codedeploy.CreateDeploymentGroupInput)
	if got := awssdk.ToString(in.DeploymentGroupName); got != "prod" {
		t.Errorf("DeploymentGroupName: got %q, want prod", got)
	}
	if got := awssdk.ToString(in.ApplicationName); got != "web-api" {
		t.Errorf("ApplicationName: got %q, want web-api", got)
	}
	if got := awssdk.ToString(in.ServiceRoleArn); got != "arn:aws:iam::1:role/CodeDeploy" {
		t.Errorf("ServiceRoleArn: got %q", got)
	}
	if got := awssdk.ToString(in.DeploymentConfigName); got != "CodeDeployDefault.OneAtATime" {
		t.Errorf("DeploymentConfigName: got %q", got)
	}
	if len(in.AutoScalingGroups) != 1 || in.AutoScalingGroups[0] != "web-asg" {
		t.Errorf("AutoScalingGroups: got %v", in.AutoScalingGroups)
	}
}

func TestDeleteDeploymentgroup(t *testing.T) {
	mock := NewMock().On("DeleteDeploymentGroup", &codedeploy.DeleteDeploymentGroupOutput{})

	Template("delete deploymentgroup name=prod application=web-api").
		Mock(mock).
		ExpectCalls("DeleteDeploymentGroup").
		Run(t)

	in := mock.InputFor("DeleteDeploymentGroup").(*codedeploy.DeleteDeploymentGroupInput)
	if got := awssdk.ToString(in.DeploymentGroupName); got != "prod" {
		t.Errorf("DeploymentGroupName: got %q", got)
	}
	if got := awssdk.ToString(in.ApplicationName); got != "web-api" {
		t.Errorf("ApplicationName: got %q", got)
	}
}

// The revision says what to deploy and its shape differs per source, so it is a document.
func TestCreateDeploymentWithARevisionDocument(t *testing.T) {
	revision := docFile(t, "revision.json", `{
	  "revisionType": "S3",
	  "s3Location": {"bucket": "releases", "key": "web-api-42.zip", "bundleType": "zip"}
	}`)

	mock := NewMock().On("CreateDeployment", &codedeploy.CreateDeploymentOutput{
		DeploymentId: awssdk.String("d-ABCDEF123"),
	})

	Template("create deployment application=web-api group=prod revision-file=" + revision +
		" description=Release42").
		Mock(mock).
		ExpectCalls("CreateDeployment").
		ExpectCommandResult("d-ABCDEF123").
		Run(t)

	in := mock.InputFor("CreateDeployment").(*codedeploy.CreateDeploymentInput)
	if in.Revision == nil {
		t.Fatal("the revision was not decoded")
	}
	if got := string(in.Revision.RevisionType); got != "S3" {
		t.Errorf("Revision.RevisionType: got %q, want S3", got)
	}
	if in.Revision.S3Location == nil {
		t.Fatal("Revision.S3Location was not decoded")
	}
	if got := awssdk.ToString(in.Revision.S3Location.Key); got != "web-api-42.zip" {
		t.Errorf("Revision.S3Location.Key: got %q", got)
	}
	if got := awssdk.ToString(in.Description); got != "Release42" {
		t.Errorf("Description: got %q, want Release42", got)
	}
}

// Halting a bad rollout is the point of this command, and rolling back rather than leaving
// it half applied is usually what is wanted.
func TestStopDeployment(t *testing.T) {
	mock := NewMock().On("StopDeployment", &codedeploy.StopDeploymentOutput{})

	Template("stop deployment id=d-ABCDEF123 rollback=true").
		Mock(mock).
		ExpectCalls("StopDeployment").
		Run(t)

	in := mock.InputFor("StopDeployment").(*codedeploy.StopDeploymentInput)
	if got := awssdk.ToString(in.DeploymentId); got != "d-ABCDEF123" {
		t.Errorf("DeploymentId: got %q", got)
	}
	if !awssdk.ToBool(in.AutoRollbackEnabled) {
		t.Error("expected AutoRollbackEnabled to be true")
	}
}

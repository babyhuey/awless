package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	beanstalktypes "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
)

func TestCreateApplication(t *testing.T) {
	mock := NewMock().On("CreateApplication", &elasticbeanstalk.CreateApplicationOutput{
		Application: &beanstalktypes.ApplicationDescription{
			ApplicationName: awssdk.String("web-api"),
		},
	})

	Template("create application name=web-api description=API").
		Mock(mock).
		ExpectCalls("CreateApplication").
		ExpectCommandResult("web-api").
		ExpectRevert("delete application name=web-api").
		Run(t)

	in := mock.InputFor("CreateApplication").(*elasticbeanstalk.CreateApplicationInput)
	if got := awssdk.ToString(in.ApplicationName); got != "web-api" {
		t.Errorf("ApplicationName: got %q, want web-api", got)
	}
	if got := awssdk.ToString(in.Description); got != "API" {
		t.Errorf("Description: got %q, want API", got)
	}
}

// Deleting an application with live environments terminates them, so the force flag must
// be opt-in and must reach the request when given.
func TestDeleteApplication(t *testing.T) {
	mock := NewMock().On("DeleteApplication", &elasticbeanstalk.DeleteApplicationOutput{})

	Template("delete application name=web-api force=true").
		Mock(mock).
		ExpectCalls("DeleteApplication").
		Run(t)

	in := mock.InputFor("DeleteApplication").(*elasticbeanstalk.DeleteApplicationInput)
	if !awssdk.ToBool(in.TerminateEnvByForce) {
		t.Error("expected TerminateEnvByForce to be true")
	}
}

func TestUpdateApplication(t *testing.T) {
	mock := NewMock().On("UpdateApplication", &elasticbeanstalk.UpdateApplicationOutput{})

	Template("update application name=web-api description=Updated").
		Mock(mock).
		ExpectCalls("UpdateApplication").
		Run(t)

	in := mock.InputFor("UpdateApplication").(*elasticbeanstalk.UpdateApplicationInput)
	if got := awssdk.ToString(in.Description); got != "Updated" {
		t.Errorf("Description: got %q, want Updated", got)
	}
}

func TestCreateEnvironmentWithSolutionStack(t *testing.T) {
	mock := NewMock().On("CreateEnvironment", &elasticbeanstalk.CreateEnvironmentOutput{
		EnvironmentName: awssdk.String("web-api-prod"),
	})

	Template(`create environment name=web-api-prod application=web-api ` +
		`solution-stack="64bit Amazon Linux 2023 v4.0.0 running Go 1" version=build-42 cname-prefix=web-api`).
		Mock(mock).
		ExpectCalls("CreateEnvironment").
		ExpectCommandResult("web-api-prod").
		ExpectRevert("delete environment name=web-api-prod").
		Run(t)

	in := mock.InputFor("CreateEnvironment").(*elasticbeanstalk.CreateEnvironmentInput)
	if got := awssdk.ToString(in.EnvironmentName); got != "web-api-prod" {
		t.Errorf("EnvironmentName: got %q", got)
	}
	if got := awssdk.ToString(in.ApplicationName); got != "web-api" {
		t.Errorf("ApplicationName: got %q, want web-api", got)
	}
	if got := awssdk.ToString(in.SolutionStackName); got != "64bit Amazon Linux 2023 v4.0.0 running Go 1" {
		t.Errorf("SolutionStackName: got %q", got)
	}
	if got := awssdk.ToString(in.VersionLabel); got != "build-42" {
		t.Errorf("VersionLabel: got %q, want build-42", got)
	}
	// Only one platform source may be set, so the others must stay nil.
	if in.PlatformArn != nil || in.TemplateName != nil {
		t.Error("only the solution stack should be set")
	}
}

// The three ways of specifying a platform are alternatives; AWS rejects more than one.
func TestCreateEnvironmentPlatformIsExclusive(t *testing.T) {
	t.Run("two platform sources are rejected", func(t *testing.T) {
		err := Template(`create environment name=e application=a solution-stack="stack" platform=arn:aws:elasticbeanstalk:::platform/custom`).
			Mock(NewMock()).
			RunExpectingError(t)
		if err == nil {
			t.Fatal("expected the platform sources to be mutually exclusive")
		}
	})

	t.Run("none is rejected", func(t *testing.T) {
		err := Template("create environment name=e application=a").
			Mock(NewMock()).
			RunExpectingError(t)
		if err == nil {
			t.Fatal("expected a create with no platform source to be rejected")
		}
	})

	t.Run("a config template is accepted", func(t *testing.T) {
		mock := NewMock().On("CreateEnvironment", &elasticbeanstalk.CreateEnvironmentOutput{})
		Template("create environment name=e application=a config-template=saved").
			Mock(mock).
			ExpectCalls("CreateEnvironment").
			Run(t)
		in := mock.InputFor("CreateEnvironment").(*elasticbeanstalk.CreateEnvironmentInput)
		if got := awssdk.ToString(in.TemplateName); got != "saved" {
			t.Errorf("TemplateName: got %q, want saved", got)
		}
	})
}

// Beanstalk calls it terminating; awless calls it deleting, which is what the verb maps to.
func TestDeleteEnvironmentTerminates(t *testing.T) {
	mock := NewMock().On("TerminateEnvironment", &elasticbeanstalk.TerminateEnvironmentOutput{})

	Template("delete environment name=web-api-prod force=true").
		Mock(mock).
		ExpectCalls("TerminateEnvironment").
		Run(t)

	in := mock.InputFor("TerminateEnvironment").(*elasticbeanstalk.TerminateEnvironmentInput)
	if got := awssdk.ToString(in.EnvironmentName); got != "web-api-prod" {
		t.Errorf("EnvironmentName: got %q", got)
	}
	if !awssdk.ToBool(in.ForceTerminate) {
		t.Error("expected ForceTerminate to be true")
	}
}

func TestUpdateEnvironmentDeploysAVersion(t *testing.T) {
	mock := NewMock().On("UpdateEnvironment", &elasticbeanstalk.UpdateEnvironmentOutput{})

	Template("update environment name=web-api-prod version=build-43").
		Mock(mock).
		ExpectCalls("UpdateEnvironment").
		Run(t)

	in := mock.InputFor("UpdateEnvironment").(*elasticbeanstalk.UpdateEnvironmentInput)
	if got := awssdk.ToString(in.VersionLabel); got != "build-43" {
		t.Errorf("VersionLabel: got %q, want build-43", got)
	}
}

// An update naming nothing to change would still trigger an environment update, which
// takes minutes and restarts instances for no reason.
func TestUpdateEnvironmentNeedsSomethingToChange(t *testing.T) {
	err := Template("update environment name=web-api-prod").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected an update with nothing to change to be rejected")
	}
}

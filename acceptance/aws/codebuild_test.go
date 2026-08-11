package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	codebuildtypes "github.com/aws/aws-sdk-go-v2/service/codebuild/types"
)

// Source, environment and artifacts are each their own struct on the request, so this
// checks the dotted paths built all three rather than dropping the values.
func TestCreateBuildproject(t *testing.T) {
	mock := NewMock().On("CreateProject", &codebuild.CreateProjectOutput{
		Project: &codebuildtypes.Project{Name: awssdk.String("api-build")},
	})

	Template("create buildproject name=api-build role=arn:aws:iam::1:role/CodeBuild " +
		"source-type=GITHUB source-location=https://github.com/acme/api " +
		"env-type=LINUX_CONTAINER image=aws/codebuild/standard:7.0 " +
		"compute-type=BUILD_GENERAL1_SMALL artifact-type=NO_ARTIFACTS timeout=30").
		Mock(mock).
		ExpectCalls("CreateProject").
		ExpectCommandResult("api-build").
		ExpectRevert("delete buildproject name=api-build").
		Run(t)

	in := mock.InputFor("CreateProject").(*codebuild.CreateProjectInput)
	if got := awssdk.ToString(in.Name); got != "api-build" {
		t.Errorf("Name: got %q, want api-build", got)
	}
	if in.Source == nil {
		t.Fatal("Source was never built")
	}
	if got := string(in.Source.Type); got != "GITHUB" {
		t.Errorf("Source.Type: got %q, want GITHUB", got)
	}
	if got := awssdk.ToString(in.Source.Location); got != "https://github.com/acme/api" {
		t.Errorf("Source.Location: got %q", got)
	}
	if in.Environment == nil {
		t.Fatal("Environment was never built")
	}
	if got := string(in.Environment.Type); got != "LINUX_CONTAINER" {
		t.Errorf("Environment.Type: got %q", got)
	}
	if got := awssdk.ToString(in.Environment.Image); got != "aws/codebuild/standard:7.0" {
		t.Errorf("Environment.Image: got %q", got)
	}
	if got := string(in.Environment.ComputeType); got != "BUILD_GENERAL1_SMALL" {
		t.Errorf("Environment.ComputeType: got %q", got)
	}
	if in.Artifacts == nil {
		t.Fatal("Artifacts was never built")
	}
	if got := string(in.Artifacts.Type); got != "NO_ARTIFACTS" {
		t.Errorf("Artifacts.Type: got %q", got)
	}
	if got := awssdk.ToInt32(in.TimeoutInMinutes); got != 30 {
		t.Errorf("TimeoutInMinutes: got %d, want 30", got)
	}
}

// AWS requires source, environment and artifacts, so a create missing one of them must be
// caught here rather than by the API.
func TestCreateBuildprojectRequiresTheThreeBlocks(t *testing.T) {
	err := Template("create buildproject name=api-build role=arn:aws:iam::1:role/CodeBuild source-type=GITHUB").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected a create without environment and artifacts to be rejected")
	}
}

// UpdateProject merges rather than replaces, so a partial update is legitimate and nothing
// beyond the name should be demanded.
func TestUpdateBuildprojectAcceptsAPartialChange(t *testing.T) {
	mock := NewMock().On("UpdateProject", &codebuild.UpdateProjectOutput{})

	Template("update buildproject name=api-build compute-type=BUILD_GENERAL1_MEDIUM timeout=45").
		Mock(mock).
		ExpectCalls("UpdateProject").
		Run(t)

	in := mock.InputFor("UpdateProject").(*codebuild.UpdateProjectInput)
	if got := string(in.Environment.ComputeType); got != "BUILD_GENERAL1_MEDIUM" {
		t.Errorf("Environment.ComputeType: got %q", got)
	}
	if got := awssdk.ToInt32(in.TimeoutInMinutes); got != 45 {
		t.Errorf("TimeoutInMinutes: got %d, want 45", got)
	}
	// Nothing was said about the source, so it must not be sent at all.
	if in.Source != nil {
		t.Error("Source should be unset when the update does not mention it")
	}
}

func TestDeleteBuildproject(t *testing.T) {
	mock := NewMock().On("DeleteProject", &codebuild.DeleteProjectOutput{})

	Template("delete buildproject name=api-build").
		Mock(mock).
		ExpectCalls("DeleteProject").
		Run(t)

	in := mock.InputFor("DeleteProject").(*codebuild.DeleteProjectInput)
	if got := awssdk.ToString(in.Name); got != "api-build" {
		t.Errorf("Name: got %q, want api-build", got)
	}
}

// Start returns the build id, which is what stop takes.
func TestStartBuildprojectReturnsTheBuildID(t *testing.T) {
	mock := NewMock().On("StartBuild", &codebuild.StartBuildOutput{
		Build: &codebuildtypes.Build{Id: awssdk.String("api-build:abc123")},
	})

	Template("start buildproject name=api-build source-version=release-2.0").
		Mock(mock).
		ExpectCalls("StartBuild").
		ExpectCommandResult("api-build:abc123").
		Run(t)

	in := mock.InputFor("StartBuild").(*codebuild.StartBuildInput)
	if got := awssdk.ToString(in.ProjectName); got != "api-build" {
		t.Errorf("ProjectName: got %q, want api-build", got)
	}
	if got := awssdk.ToString(in.SourceVersion); got != "release-2.0" {
		t.Errorf("SourceVersion: got %q, want release-2.0", got)
	}
}

func TestStartBuildprojectEmptyResponse(t *testing.T) {
	mock := NewMock().On("StartBuild", &codebuild.StartBuildOutput{})

	Template("start buildproject name=api-build").
		Mock(mock).
		ExpectCalls("StartBuild").
		ExpectCommandResult("api-build").
		Run(t)
}

// Stop addresses a build, not the project, which is why the param is named for it.
func TestStopBuildproject(t *testing.T) {
	mock := NewMock().On("StopBuild", &codebuild.StopBuildOutput{})

	Template("stop buildproject build=api-build:abc123").
		Mock(mock).
		ExpectCalls("StopBuild").
		Run(t)

	in := mock.InputFor("StopBuild").(*codebuild.StopBuildInput)
	if got := awssdk.ToString(in.Id); got != "api-build:abc123" {
		t.Errorf("Id: got %q, want api-build:abc123", got)
	}
}

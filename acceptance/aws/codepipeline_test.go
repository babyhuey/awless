package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
)

func TestDeletePipeline(t *testing.T) {
	mock := NewMock().On("DeletePipeline", &codepipeline.DeletePipelineOutput{})

	Template("delete pipeline name=build-and-deploy").
		Mock(mock).
		ExpectCalls("DeletePipeline").
		Run(t)

	in := mock.InputFor("DeletePipeline").(*codepipeline.DeletePipelineInput)
	if got := awssdk.ToString(in.Name); got != "build-and-deploy" {
		t.Errorf("Name: got %q, want build-and-deploy", got)
	}
}

// The execution id is the only thing `stop pipeline` can be driven with, so start has to
// surface it rather than echoing the pipeline name.
func TestStartPipelineReturnsTheExecutionID(t *testing.T) {
	mock := NewMock().On("StartPipelineExecution", &codepipeline.StartPipelineExecutionOutput{
		PipelineExecutionId: awssdk.String("11111111-2222-3333-4444-555555555555"),
	})

	Template("start pipeline name=build-and-deploy").
		Mock(mock).
		ExpectCalls("StartPipelineExecution").
		ExpectCommandResult("11111111-2222-3333-4444-555555555555").
		Run(t)

	in := mock.InputFor("StartPipelineExecution").(*codepipeline.StartPipelineExecutionInput)
	if got := awssdk.ToString(in.Name); got != "build-and-deploy" {
		t.Errorf("Name: got %q, want build-and-deploy", got)
	}
}

// An empty response falls back to the pipeline name rather than dereferencing a nil id.
func TestStartPipelineEmptyResponse(t *testing.T) {
	mock := NewMock().On("StartPipelineExecution", &codepipeline.StartPipelineExecutionOutput{})

	Template("start pipeline name=build-and-deploy").
		Mock(mock).
		ExpectCalls("StartPipelineExecution").
		ExpectCommandResult("build-and-deploy").
		Run(t)
}

func TestStopPipeline(t *testing.T) {
	mock := NewMock().On("StopPipelineExecution", &codepipeline.StopPipelineExecutionOutput{})

	// The execution id is quoted: an all-numeric hyphenated value is not accepted bare
	// by the template grammar.
	Template(`stop pipeline name=build-and-deploy execution='11111111-2222-3333-4444-555555555555' reason=Superseded abandon=true`).
		Mock(mock).
		ExpectCalls("StopPipelineExecution").
		Run(t)

	in := mock.InputFor("StopPipelineExecution").(*codepipeline.StopPipelineExecutionInput)
	if got := awssdk.ToString(in.PipelineName); got != "build-and-deploy" {
		t.Errorf("PipelineName: got %q, want build-and-deploy", got)
	}
	if got := awssdk.ToString(in.PipelineExecutionId); got != "11111111-2222-3333-4444-555555555555" {
		t.Errorf("PipelineExecutionId: got %q", got)
	}
	if !in.Abandon {
		t.Error("expected Abandon to be true")
	}
	if got := awssdk.ToString(in.Reason); got != "Superseded" {
		t.Errorf("Reason: got %q, want Superseded", got)
	}
}

// A bare UUID does not parse, which is easy to hit and confusing when it happens. Pinned
// so the documented workaround stays true.
func TestStopPipelineRejectsAnUnquotedUUID(t *testing.T) {
	err := Template("stop pipeline name=b execution=11111111-2222-3333-4444-555555555555").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected an unquoted UUID to fail parsing; if this now works, drop the quoting note from the docs")
	}
}

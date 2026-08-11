package awsat

import (
	"os"
	"path/filepath"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

// definitionFile writes an Amazon States Language document and returns its path. The
// definition is passed as a file rather than inline, so the test has to provide one for
// the awsfiletostring setter to read.
func definitionFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "flow.json")
	body := `{"Comment":"test","StartAt":"Done","States":{"Done":{"Type":"Succeed"}}}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("writing the definition: %s", err)
	}
	return path
}

func TestCreateStatemachine(t *testing.T) {
	arn := "arn:aws:states:us-west-2:123456789012:stateMachine:order-flow"
	mock := NewMock().On("CreateStateMachine", &sfn.CreateStateMachineOutput{
		StateMachineArn: awssdk.String(arn),
	})

	def := definitionFile(t)
	Template("create statemachine name=order-flow role=arn:aws:iam::1:role/sfn " +
		"definition-file=" + def + " type=EXPRESS publish=true").
		Mock(mock).
		ExpectCalls("CreateStateMachine").
		ExpectCommandResult(arn).
		ExpectRevert("delete statemachine arn=" + arn).
		Run(t)

	in := mock.InputFor("CreateStateMachine").(*sfn.CreateStateMachineInput)
	if got := awssdk.ToString(in.Name); got != "order-flow" {
		t.Errorf("Name: got %q, want order-flow", got)
	}
	if got := awssdk.ToString(in.RoleArn); got != "arn:aws:iam::1:role/sfn" {
		t.Errorf("RoleArn: got %q", got)
	}
	// The file's contents must reach the request, not its path.
	if got := awssdk.ToString(in.Definition); got == def || got == "" {
		t.Errorf("Definition should hold the file contents, got %q", got)
	}
	if got := string(in.Type); got != "EXPRESS" {
		t.Errorf("Type: got %q, want EXPRESS", got)
	}
	// Publish is a plain bool on this request rather than a *bool; the setter
	// dereferences the spec's *bool into it.
	if !in.Publish {
		t.Error("expected Publish to be true")
	}
}

// An empty response must fall back to the name rather than dereferencing a nil ARN.
func TestCreateStatemachineEmptyResponse(t *testing.T) {
	mock := NewMock().On("CreateStateMachine", &sfn.CreateStateMachineOutput{})

	Template("create statemachine name=order-flow role=arn:aws:iam::1:role/sfn " +
		"definition-file=" + definitionFile(t)).
		Mock(mock).
		ExpectCalls("CreateStateMachine").
		ExpectCommandResult("order-flow").
		Run(t)
}

func TestDeleteStatemachine(t *testing.T) {
	arn := "arn:aws:states:us-west-2:123456789012:stateMachine:order-flow"
	mock := NewMock().On("DeleteStateMachine", &sfn.DeleteStateMachineOutput{})

	Template("delete statemachine arn=" + arn).
		Mock(mock).
		ExpectCalls("DeleteStateMachine").
		Run(t)

	in := mock.InputFor("DeleteStateMachine").(*sfn.DeleteStateMachineInput)
	if got := awssdk.ToString(in.StateMachineArn); got != arn {
		t.Errorf("StateMachineArn: got %q, want %q", got, arn)
	}
}

func TestUpdateStatemachine(t *testing.T) {
	arn := "arn:aws:states:us-west-2:123456789012:stateMachine:order-flow"
	mock := NewMock().On("UpdateStateMachine", &sfn.UpdateStateMachineOutput{})

	Template("update statemachine arn=" + arn + " definition-file=" + definitionFile(t) + " publish=true").
		Mock(mock).
		ExpectCalls("UpdateStateMachine").
		Run(t)

	in := mock.InputFor("UpdateStateMachine").(*sfn.UpdateStateMachineInput)
	if got := awssdk.ToString(in.StateMachineArn); got != arn {
		t.Errorf("StateMachineArn: got %q", got)
	}
	if awssdk.ToString(in.Definition) == "" {
		t.Error("Definition was not read from the file")
	}
}

// UpdateStateMachine leaves anything omitted untouched, but an update that changes
// nothing is a no-op that still bumps the revision, so one of the two must be given.
func TestUpdateStatemachineNeedsSomethingToChange(t *testing.T) {
	err := Template("update statemachine arn=arn:aws:states:us-west-2:1:stateMachine:f publish=true").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected an update with neither definition nor role to be rejected")
	}
}

func TestStartExecution(t *testing.T) {
	execArn := "arn:aws:states:us-west-2:123456789012:execution:order-flow:order-4711"
	mock := NewMock().On("StartExecution", &sfn.StartExecutionOutput{
		ExecutionArn: awssdk.String(execArn),
	})

	Template(`start execution statemachine=arn:aws:states:us-west-2:1:stateMachine:order-flow ` +
		`name=order-4711 input='{"orderId":4711}'`).
		Mock(mock).
		ExpectCalls("StartExecution").
		ExpectCommandResult(execArn).
		Run(t)

	in := mock.InputFor("StartExecution").(*sfn.StartExecutionInput)
	if got := awssdk.ToString(in.StateMachineArn); got != "arn:aws:states:us-west-2:1:stateMachine:order-flow" {
		t.Errorf("StateMachineArn: got %q", got)
	}
	if got := awssdk.ToString(in.Name); got != "order-4711" {
		t.Errorf("Name: got %q, want order-4711", got)
	}
	if got := awssdk.ToString(in.Input); got != `{"orderId":4711}` {
		t.Errorf("Input: got %q", got)
	}
}

func TestStopExecution(t *testing.T) {
	execArn := "arn:aws:states:us-west-2:123456789012:execution:order-flow:order-4711"
	mock := NewMock().On("StopExecution", &sfn.StopExecutionOutput{})

	Template("stop execution arn=" + execArn + " cause=Superseded error=Replaced").
		Mock(mock).
		ExpectCalls("StopExecution").
		Run(t)

	in := mock.InputFor("StopExecution").(*sfn.StopExecutionInput)
	if got := awssdk.ToString(in.ExecutionArn); got != execArn {
		t.Errorf("ExecutionArn: got %q", got)
	}
	if got := awssdk.ToString(in.Cause); got != "Superseded" {
		t.Errorf("Cause: got %q, want Superseded", got)
	}
	if got := awssdk.ToString(in.Error); got != "Replaced" {
		t.Errorf("Error: got %q, want Replaced", got)
	}
}

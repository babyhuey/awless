package awsspec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
)

func writeDoc(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "doc.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("writing the document: %s", err)
	}
	return path
}

// The AWS SDK v2 types have no json tags, so this depends on Go matching field names
// case-insensitively. That is what lets the camelCase the AWS CLI accepts land on the
// exported Go fields, and it is worth pinning because nothing else in the tree relies on
// it.
func TestAwsfiletostructDecodesAnAWSStyleDocument(t *testing.T) {
	path := writeDoc(t, `{
	  "name": "build-and-deploy",
	  "roleArn": "arn:aws:iam::123456789012:role/CodePipeline",
	  "artifactStore": {"type": "S3", "location": "my-artifacts"},
	  "stages": [
	    {"name": "Source", "actions": [
	      {"name": "Source",
	       "actionTypeId": {"category":"Source","owner":"AWS","provider":"S3","version":"1"},
	       "configuration": {"S3Bucket":"src","S3ObjectKey":"app.zip"}}
	    ]}
	  ]
	}`)

	input := &codepipeline.CreatePipelineInput{}
	if err := setFieldWithType(context.Background(), path, input, "Pipeline", awsfiletostruct); err != nil {
		t.Fatalf("setting the document: %s", err)
	}

	if input.Pipeline == nil {
		t.Fatal("the pipeline declaration was not built")
	}
	if got := awssdk.ToString(input.Pipeline.Name); got != "build-and-deploy" {
		t.Errorf("Name: got %q, want build-and-deploy", got)
	}
	if got := awssdk.ToString(input.Pipeline.RoleArn); got != "arn:aws:iam::123456789012:role/CodePipeline" {
		t.Errorf("RoleArn: got %q", got)
	}
	// A nested struct behind a pointer.
	if input.Pipeline.ArtifactStore == nil {
		t.Fatal("ArtifactStore was not decoded")
	}
	if got := string(input.Pipeline.ArtifactStore.Type); got != "S3" {
		t.Errorf("ArtifactStore.Type: got %q, want S3 — enum types must decode", got)
	}
	// A slice of structs, nested two deep, with a map and another enum inside.
	if len(input.Pipeline.Stages) != 1 {
		t.Fatalf("Stages: got %d, want 1", len(input.Pipeline.Stages))
	}
	actions := input.Pipeline.Stages[0].Actions
	if len(actions) != 1 {
		t.Fatalf("Actions: got %d, want 1", len(actions))
	}
	if got := string(actions[0].ActionTypeId.Category); got != "Source" {
		t.Errorf("ActionTypeId.Category: got %q, want Source", got)
	}
	if got := actions[0].Configuration["S3Bucket"]; got != "src" {
		t.Errorf("Configuration[S3Bucket]: got %q, want src", got)
	}
}

// A misspelled key must fail rather than be dropped. Silently ignoring it would send a
// request missing the thing the user configured — the same failure mode as a wrong
// awsName, and just as invisible.
func TestAwsfiletostructRejectsAnUnknownField(t *testing.T) {
	path := writeDoc(t, `{"name": "p", "roleARNN": "typo"}`)

	err := setFieldWithType(context.Background(), path, &codepipeline.CreatePipelineInput{}, "Pipeline", awsfiletostruct)
	if err == nil {
		t.Fatal("expected an unknown field to be rejected")
	}
	if !strings.Contains(err.Error(), "roleARNN") {
		t.Errorf("the error should name the offending key, got: %s", err)
	}
}

// A wrapper object is a common mistake, since the AWS CLI takes {"pipeline": {...}}. It
// must produce a clear error rather than an empty declaration.
func TestAwsfiletostructRejectsAWrapperObject(t *testing.T) {
	path := writeDoc(t, `{"pipeline": {"name": "p"}}`)

	if err := setFieldWithType(context.Background(), path, &codepipeline.CreatePipelineInput{}, "Pipeline", awsfiletostruct); err == nil {
		t.Fatal("expected a wrapped document to be rejected rather than silently empty")
	}
}

func TestAwsfiletostructErrors(t *testing.T) {
	t.Run("a missing file", func(t *testing.T) {
		err := setFieldWithType(context.Background(), "/nonexistent/doc.json",
			&codepipeline.CreatePipelineInput{}, "Pipeline", awsfiletostruct)
		if err == nil {
			t.Fatal("expected a missing file to error")
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		path := writeDoc(t, `{"name": `)
		err := setFieldWithType(context.Background(), path,
			&codepipeline.CreatePipelineInput{}, "Pipeline", awsfiletostruct)
		if err == nil {
			t.Fatal("expected malformed JSON to error")
		}
	})

	t.Run("a field that does not exist", func(t *testing.T) {
		path := writeDoc(t, `{}`)
		err := setFieldWithType(context.Background(), path,
			&codepipeline.CreatePipelineInput{}, "Nonexistent", awsfiletostruct)
		if err == nil {
			t.Fatal("expected an unresolvable field path to error")
		}
		if !strings.Contains(err.Error(), "Nonexistent") {
			t.Errorf("the error should name the path, got: %s", err)
		}
	})
}

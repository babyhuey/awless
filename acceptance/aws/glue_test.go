package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
)

// CreateDatabase wraps everything in a DatabaseInput, so this checks the dotted paths built
// it rather than dropping the values.
func TestCreateGluedatabase(t *testing.T) {
	mock := NewMock().On("CreateDatabase", &glue.CreateDatabaseOutput{})

	Template("create gluedatabase name=analytics description=Catalog location=s3://lake/analytics").
		Mock(mock).
		ExpectCalls("CreateDatabase").
		ExpectCommandResult("analytics").
		ExpectRevert("delete gluedatabase name=analytics").
		Run(t)

	in := mock.InputFor("CreateDatabase").(*glue.CreateDatabaseInput)
	if in.DatabaseInput == nil {
		t.Fatal("DatabaseInput was never built")
	}
	if got := awssdk.ToString(in.DatabaseInput.Name); got != "analytics" {
		t.Errorf("DatabaseInput.Name: got %q, want analytics", got)
	}
	if got := awssdk.ToString(in.DatabaseInput.Description); got != "Catalog" {
		t.Errorf("DatabaseInput.Description: got %q", got)
	}
	if got := awssdk.ToString(in.DatabaseInput.LocationUri); got != "s3://lake/analytics" {
		t.Errorf("DatabaseInput.LocationUri: got %q", got)
	}
}

func TestDeleteGluedatabase(t *testing.T) {
	mock := NewMock().On("DeleteDatabase", &glue.DeleteDatabaseOutput{})

	Template("delete gluedatabase name=analytics").
		Mock(mock).
		ExpectCalls("DeleteDatabase").
		Run(t)

	in := mock.InputFor("DeleteDatabase").(*glue.DeleteDatabaseInput)
	if got := awssdk.ToString(in.Name); got != "analytics" {
		t.Errorf("Name: got %q, want analytics", got)
	}
}

// Crawl targets differ in shape per source, so they arrive as a document.
func TestCreateCrawlerWithTargets(t *testing.T) {
	targets := docFile(t, "targets.json",
		`{"s3Targets": [{"path": "s3://lake/events/"}, {"path": "s3://lake/users/"}]}`)

	mock := NewMock().On("CreateCrawler", &glue.CreateCrawlerOutput{})

	Template("create crawler name=events-crawler role=AWSGlueServiceRole database=analytics " +
		"targets-file=" + targets + ` schedule="cron(0 6 * * ? *)"`).
		Mock(mock).
		ExpectCalls("CreateCrawler").
		ExpectCommandResult("events-crawler").
		ExpectRevert("delete crawler name=events-crawler").
		Run(t)

	in := mock.InputFor("CreateCrawler").(*glue.CreateCrawlerInput)
	if got := awssdk.ToString(in.Role); got != "AWSGlueServiceRole" {
		t.Errorf("Role: got %q", got)
	}
	if got := awssdk.ToString(in.DatabaseName); got != "analytics" {
		t.Errorf("DatabaseName: got %q, want analytics", got)
	}
	if in.Targets == nil {
		t.Fatal("Targets was not decoded")
	}
	if len(in.Targets.S3Targets) != 2 {
		t.Fatalf("Targets.S3Targets: got %d, want 2", len(in.Targets.S3Targets))
	}
	if got := awssdk.ToString(in.Targets.S3Targets[1].Path); got != "s3://lake/users/" {
		t.Errorf("Targets.S3Targets[1].Path: got %q", got)
	}
	if got := awssdk.ToString(in.Schedule); got != "cron(0 6 * * ? *)" {
		t.Errorf("Schedule: got %q", got)
	}
}

func TestStartAndStopCrawler(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		mock := NewMock().On("StartCrawler", &glue.StartCrawlerOutput{})
		Template("start crawler name=events-crawler").
			Mock(mock).ExpectCalls("StartCrawler").Run(t)
		in := mock.InputFor("StartCrawler").(*glue.StartCrawlerInput)
		if got := awssdk.ToString(in.Name); got != "events-crawler" {
			t.Errorf("Name: got %q", got)
		}
	})

	t.Run("stop", func(t *testing.T) {
		mock := NewMock().On("StopCrawler", &glue.StopCrawlerOutput{})
		Template("stop crawler name=events-crawler").
			Mock(mock).ExpectCalls("StopCrawler").Run(t)
		in := mock.InputFor("StopCrawler").(*glue.StopCrawlerInput)
		if got := awssdk.ToString(in.Name); got != "events-crawler" {
			t.Errorf("Name: got %q", got)
		}
	})
}

// The job command is nested but flat enough to spell out, so this checks all three of its
// fields landed in the same struct rather than in separate ones.
func TestCreateJob(t *testing.T) {
	mock := NewMock().On("CreateJob", &glue.CreateJobOutput{Name: awssdk.String("etl-events")})

	Template("create job name=etl-events role=AWSGlueServiceRole command=glueetl " +
		"script=s3://scripts/etl.py python-version=3 glue-version=4.0 " +
		"worker-type=G.1X workers=2 max-retries=1 timeout=60").
		Mock(mock).
		ExpectCalls("CreateJob").
		ExpectCommandResult("etl-events").
		ExpectRevert("delete job name=etl-events").
		Run(t)

	in := mock.InputFor("CreateJob").(*glue.CreateJobInput)
	if in.Command == nil {
		t.Fatal("Command was never built")
	}
	if got := awssdk.ToString(in.Command.Name); got != "glueetl" {
		t.Errorf("Command.Name: got %q, want glueetl", got)
	}
	if got := awssdk.ToString(in.Command.ScriptLocation); got != "s3://scripts/etl.py" {
		t.Errorf("Command.ScriptLocation: got %q", got)
	}
	if got := awssdk.ToString(in.Command.PythonVersion); got != "3" {
		t.Errorf("Command.PythonVersion: got %q, want 3", got)
	}
	if got := string(in.WorkerType); got != "G.1X" {
		t.Errorf("WorkerType: got %q, want G.1X", got)
	}
	if got := awssdk.ToInt32(in.NumberOfWorkers); got != 2 {
		t.Errorf("NumberOfWorkers: got %d, want 2", got)
	}
	if got := awssdk.ToInt32(in.Timeout); got != 60 {
		t.Errorf("Timeout: got %d, want 60", got)
	}
}

// Starting a job returns a run id, which is what identifies the execution afterwards.
func TestStartJobReturnsTheRunID(t *testing.T) {
	mock := NewMock().On("StartJobRun", &glue.StartJobRunOutput{
		JobRunId: awssdk.String("jr_abc123"),
	})

	Template("start job name=etl-events workers=10 worker-type=G.2X").
		Mock(mock).
		ExpectCalls("StartJobRun").
		ExpectCommandResult("jr_abc123").
		Run(t)

	in := mock.InputFor("StartJobRun").(*glue.StartJobRunInput)
	if got := awssdk.ToString(in.JobName); got != "etl-events" {
		t.Errorf("JobName: got %q", got)
	}
	if got := awssdk.ToInt32(in.NumberOfWorkers); got != 10 {
		t.Errorf("NumberOfWorkers: got %d, want 10", got)
	}
}

func TestDeleteJob(t *testing.T) {
	mock := NewMock().On("DeleteJob", &glue.DeleteJobOutput{})

	Template("delete job name=etl-events").
		Mock(mock).
		ExpectCalls("DeleteJob").
		Run(t)

	in := mock.InputFor("DeleteJob").(*glue.DeleteJobInput)
	if got := awssdk.ToString(in.JobName); got != "etl-events" {
		t.Errorf("JobName: got %q, want etl-events", got)
	}
}

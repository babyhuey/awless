package awstailers

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func TestEventPrint(t *testing.T) {
	stamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	e := &event{
		id:      "act-123",
		element: "my-asg",
		stamp:   stamp,
		message: "Successful: launched instance i-abc",
	}

	var buf bytes.Buffer
	err := e.print(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "my-asg") {
		t.Errorf("expected output to contain element name, got: %s", output)
	}
	if !strings.Contains(output, "Successful: launched instance i-abc") {
		t.Errorf("expected output to contain message, got: %s", output)
	}
	if !strings.Contains(output, "2025") {
		t.Errorf("expected output to contain timestamp year, got: %s", output)
	}
	// Verify format: "timestamp: element\n\tmessage\n"
	if !strings.Contains(output, ": my-asg\n\t") {
		t.Errorf("expected format 'timestamp: element\\n\\tmessage', got: %s", output)
	}
}

func TestNewEventFromScalingActivity(t *testing.T) {
	startTime := time.Date(2025, 3, 10, 14, 0, 0, 0, time.UTC)
	activity := autoscalingtypes.Activity{
		ActivityId:           aws.String("act-456"),
		StartTime:            aws.Time(startTime),
		StatusCode:           autoscalingtypes.ScalingActivityStatusCodeSuccessful,
		Description:          aws.String("Launching a new EC2 instance: i-xyz"),
		AutoScalingGroupName: aws.String("prod-asg"),
	}

	e := newEventFromScalingActivity(activity)

	if e.id != "act-456" {
		t.Errorf("expected id 'act-456', got %q", e.id)
	}
	if !e.stamp.Equal(startTime) {
		t.Errorf("expected stamp %v, got %v", startTime, e.stamp)
	}
	if e.element != "prod-asg" {
		t.Errorf("expected element 'prod-asg', got %q", e.element)
	}
	if !strings.Contains(e.message, "Successful") {
		t.Errorf("expected message to contain status code, got %q", e.message)
	}
	if !strings.Contains(e.message, "Launching a new EC2 instance: i-xyz") {
		t.Errorf("expected message to contain description, got %q", e.message)
	}
}

func TestNewScalingActivitiesTailer(t *testing.T) {
	tailer := NewScalingActivitiesTailer(25, true, 10*time.Second)

	if tailer.nbEvents != 25 {
		t.Errorf("expected nbEvents 25, got %d", tailer.nbEvents)
	}
	if !tailer.follow {
		t.Error("expected follow to be true")
	}
	if tailer.pollingFrequency != 10*time.Second {
		t.Errorf("expected pollingFrequency 10s, got %v", tailer.pollingFrequency)
	}
}

func TestNewCloudformationEventsTailer(t *testing.T) {
	f := filters{StackEventLogicalID, StackEventStatus}
	tailer := NewCloudformationEventsTailer("my-stack", 50, true, 7*time.Second, f, 30*time.Minute, true)

	if tailer.stackName != "my-stack" {
		t.Errorf("expected stackName 'my-stack', got %q", tailer.stackName)
	}
	if tailer.nbEvents != 50 {
		t.Errorf("expected nbEvents 50, got %d", tailer.nbEvents)
	}
	if !tailer.follow {
		t.Error("expected follow to be true")
	}
	if tailer.pollingFrequency != 7*time.Second {
		t.Errorf("expected pollingFrequency 7s, got %v", tailer.pollingFrequency)
	}
	if len(tailer.filters) != 2 {
		t.Errorf("expected 2 filters, got %d", len(tailer.filters))
	}
	if tailer.timeout != 30*time.Minute {
		t.Errorf("expected timeout 30m, got %v", tailer.timeout)
	}
	if !tailer.cancelAfterTimeout {
		t.Error("expected cancelAfterTimeout to be true")
	}
}

func TestStackEventFilter(t *testing.T) {
	ts := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	e := stackEvent{
		cloudformationtypes.StackEvent{
			LogicalResourceId:    aws.String("MyBucket"),
			Timestamp:            aws.Time(ts),
			ResourceStatus:       cloudformationtypes.ResourceStatusCreateComplete,
			ResourceStatusReason: aws.String("Resource creation complete"),
			ResourceType:         aws.String("AWS::S3::Bucket"),
		},
	}

	t.Run("single filter logical id", func(t *testing.T) {
		out := e.filter([]string{StackEventLogicalID})
		if !strings.Contains(string(out), "MyBucket") {
			t.Errorf("expected output to contain 'MyBucket', got %q", string(out))
		}
	})

	t.Run("single filter timestamp", func(t *testing.T) {
		out := e.filter([]string{StackEventTimestamp})
		if !strings.Contains(string(out), "2025-06-01") {
			t.Errorf("expected output to contain timestamp, got %q", string(out))
		}
	})

	t.Run("single filter type", func(t *testing.T) {
		out := e.filter([]string{StackEventType})
		if !strings.Contains(string(out), "AWS::S3::Bucket") {
			t.Errorf("expected output to contain resource type, got %q", string(out))
		}
	})

	t.Run("single filter reason", func(t *testing.T) {
		out := e.filter([]string{StackEventStatusReason})
		if !strings.Contains(string(out), "Resource creation complete") {
			t.Errorf("expected output to contain reason, got %q", string(out))
		}
	})

	t.Run("multiple filters", func(t *testing.T) {
		out := e.filter([]string{StackEventLogicalID, StackEventType, StackEventStatusReason})
		s := string(out)
		if !strings.Contains(s, "MyBucket") {
			t.Errorf("expected 'MyBucket' in output, got %q", s)
		}
		if !strings.Contains(s, "AWS::S3::Bucket") {
			t.Errorf("expected 'AWS::S3::Bucket' in output, got %q", s)
		}
		if !strings.Contains(s, "Resource creation complete") {
			t.Errorf("expected reason in output, got %q", s)
		}
		// Fields should be tab-separated
		if !strings.Contains(s, "\t") {
			t.Errorf("expected tab-separated fields, got %q", s)
		}
	})

	t.Run("ends with newline", func(t *testing.T) {
		out := e.filter([]string{StackEventLogicalID})
		if !strings.HasSuffix(string(out), "\n") {
			t.Errorf("expected output to end with newline, got %q", string(out))
		}
	})
}

func TestStackEventIsDeploymentStart(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   *string
		resourceStatus cloudformationtypes.ResourceStatus
		want           bool
	}{
		{
			name:           "create in progress on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusCreateInProgress,
			want:           true,
		},
		{
			name:           "update in progress on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusUpdateInProgress,
			want:           true,
		},
		{
			name:           "delete in progress on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusDeleteInProgress,
			want:           true,
		},
		{
			name:           "create complete on stack - not start",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusCreateComplete,
			want:           false,
		},
		{
			name:           "create in progress on non-stack resource",
			resourceType:   aws.String("AWS::S3::Bucket"),
			resourceStatus: cloudformationtypes.ResourceStatusCreateInProgress,
			want:           false,
		},
		{
			name:           "nil resource type",
			resourceType:   nil,
			resourceStatus: cloudformationtypes.ResourceStatusCreateInProgress,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &stackEvent{
				cloudformationtypes.StackEvent{
					ResourceType:   tt.resourceType,
					ResourceStatus: tt.resourceStatus,
				},
			}
			got := e.isDeploymentStart()
			if got != tt.want {
				t.Errorf("isDeploymentStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStackEventIsDeploymentFinished(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   *string
		resourceStatus cloudformationtypes.ResourceStatus
		want           bool
	}{
		{
			name:           "create complete on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusCreateComplete,
			want:           true,
		},
		{
			name:           "update complete on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusUpdateComplete,
			want:           true,
		},
		{
			name:           "delete complete on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusDeleteComplete,
			want:           true,
		},
		{
			name:           "create failed on stack",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusCreateFailed,
			want:           true,
		},
		{
			name:           "in progress on stack - not finished",
			resourceType:   aws.String(cfnStackResourceType),
			resourceStatus: cloudformationtypes.ResourceStatusCreateInProgress,
			want:           false,
		},
		{
			name:           "complete on non-stack resource",
			resourceType:   aws.String("AWS::EC2::Instance"),
			resourceStatus: cloudformationtypes.ResourceStatusCreateComplete,
			want:           false,
		},
		{
			name:           "nil resource type",
			resourceType:   nil,
			resourceStatus: cloudformationtypes.ResourceStatusCreateComplete,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &stackEvent{
				cloudformationtypes.StackEvent{
					ResourceType:   tt.resourceType,
					ResourceStatus: tt.resourceStatus,
				},
			}
			got := e.isDeploymentFinished()
			if got != tt.want {
				t.Errorf("isDeploymentFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStackEventIsFailed(t *testing.T) {
	tests := []struct {
		name           string
		resourceStatus cloudformationtypes.ResourceStatus
		want           bool
	}{
		{
			name:           "create failed",
			resourceStatus: cloudformationtypes.ResourceStatusCreateFailed,
			want:           true,
		},
		{
			name:           "update failed",
			resourceStatus: cloudformationtypes.ResourceStatusUpdateFailed,
			want:           true,
		},
		{
			name:           "delete failed",
			resourceStatus: cloudformationtypes.ResourceStatusDeleteFailed,
			want:           true,
		},
		{
			name:           "update rollback in progress",
			resourceStatus: cloudformationtypes.ResourceStatusUpdateRollbackInProgress,
			want:           true,
		},
		{
			name:           "create complete - not failed",
			resourceStatus: cloudformationtypes.ResourceStatusCreateComplete,
			want:           false,
		},
		{
			name:           "create in progress - not failed",
			resourceStatus: cloudformationtypes.ResourceStatusCreateInProgress,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &stackEvent{
				cloudformationtypes.StackEvent{
					ResourceStatus: tt.resourceStatus,
				},
			}
			got := e.isFailed()
			if got != tt.want {
				t.Errorf("isFailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiltersHeader(t *testing.T) {
	tests := []struct {
		name   string
		f      filters
		expect []string
	}{
		{
			name:   "all filters",
			f:      filters{StackEventLogicalID, StackEventTimestamp, StackEventStatus, StackEventStatusReason, StackEventType},
			expect: []string{"Logical ID", "Timestamp", "Status", "Status Reason", "Type"},
		},
		{
			name:   "subset of filters",
			f:      filters{StackEventLogicalID, StackEventStatus},
			expect: []string{"Logical ID", "Status"},
		},
		{
			name:   "single filter",
			f:      filters{StackEventType},
			expect: []string{"Type"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := string(tt.f.header())
			for _, exp := range tt.expect {
				if !strings.Contains(header, exp) {
					t.Errorf("expected header to contain %q, got %q", exp, header)
				}
			}
			if !strings.HasSuffix(header, "\n") {
				t.Errorf("expected header to end with newline, got %q", header)
			}
		})
	}
}

func TestScalingActivitiesTailerName(t *testing.T) {
	tailer := NewScalingActivitiesTailer(10, false, 5*time.Second)
	if tailer.Name() != "scaling-activities" {
		t.Errorf("expected Name() to return 'scaling-activities', got %q", tailer.Name())
	}
}

func TestCloudformationEventsTailerName(t *testing.T) {
	tailer := NewCloudformationEventsTailer("stack", 10, false, 5*time.Second, nil, time.Minute, false)
	if tailer.Name() != "stack-events" {
		t.Errorf("expected Name() to return 'stack-events', got %q", tailer.Name())
	}
}

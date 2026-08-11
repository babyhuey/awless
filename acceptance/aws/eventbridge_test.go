package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

func TestCreateEventbus(t *testing.T) {
	mock := NewMock().On("CreateEventBus", &eventbridge.CreateEventBusOutput{
		EventBusArn: awssdk.String("arn:aws:events:us-west-2:123456789012:event-bus/orders"),
	})

	Template("create eventbus name=orders description=Orders kms-key=alias/events").
		Mock(mock).
		ExpectCalls("CreateEventBus").
		// The ARN is what AWS returns, but the name is what delete and the graph key on.
		ExpectCommandResult("orders").
		ExpectRevert("delete eventbus name=orders").
		Run(t)

	in := mock.InputFor("CreateEventBus").(*eventbridge.CreateEventBusInput)
	if got := awssdk.ToString(in.Name); got != "orders" {
		t.Errorf("Name: got %q, want orders", got)
	}
	if got := awssdk.ToString(in.Description); got != "Orders" {
		t.Errorf("Description: got %q, want Orders", got)
	}
	if got := awssdk.ToString(in.KmsKeyIdentifier); got != "alias/events" {
		t.Errorf("KmsKeyIdentifier: got %q, want alias/events", got)
	}
}

func TestDeleteEventbus(t *testing.T) {
	mock := NewMock().On("DeleteEventBus", &eventbridge.DeleteEventBusOutput{})

	Template("delete eventbus name=orders").
		Mock(mock).
		ExpectCalls("DeleteEventBus").
		Run(t)

	in := mock.InputFor("DeleteEventBus").(*eventbridge.DeleteEventBusInput)
	if got := awssdk.ToString(in.Name); got != "orders" {
		t.Errorf("Name: got %q, want orders", got)
	}
}

func TestCreateEventruleWithSchedule(t *testing.T) {
	mock := NewMock().On("PutRule", &eventbridge.PutRuleOutput{
		RuleArn: awssdk.String("arn:aws:events:us-west-2:123456789012:rule/nightly"),
	})

	Template(`create eventrule name=nightly schedule="rate(1 hour)" description=Nightly state=ENABLED`).
		Mock(mock).
		ExpectCalls("PutRule").
		ExpectCommandResult("nightly").
		ExpectRevert("delete eventrule name=nightly").
		Run(t)

	in := mock.InputFor("PutRule").(*eventbridge.PutRuleInput)
	if got := awssdk.ToString(in.Name); got != "nightly" {
		t.Errorf("Name: got %q, want nightly", got)
	}
	if got := awssdk.ToString(in.ScheduleExpression); got != "rate(1 hour)" {
		t.Errorf("ScheduleExpression: got %q, want rate(1 hour)", got)
	}
	// State is an enum on the request rather than a *string.
	if got := string(in.State); got != "ENABLED" {
		t.Errorf("State: got %q, want ENABLED", got)
	}
}

func TestCreateEventruleWithPattern(t *testing.T) {
	mock := NewMock().On("PutRule", &eventbridge.PutRuleOutput{})

	// Single quotes, because the template grammar has no escape for a double quote
	// inside a double-quoted value — which is the only way to write JSON inline.
	Template(`create eventrule name=on-stop pattern='{"source":["aws.ec2"]}' eventbus=orders role=arn:aws:iam::1:role/events`).
		Mock(mock).
		ExpectCalls("PutRule").
		ExpectCommandResult("on-stop").
		Run(t)

	in := mock.InputFor("PutRule").(*eventbridge.PutRuleInput)
	if got := awssdk.ToString(in.EventPattern); got != `{"source":["aws.ec2"]}` {
		t.Errorf("EventPattern: got %q", got)
	}
	if got := awssdk.ToString(in.EventBusName); got != "orders" {
		t.Errorf("EventBusName: got %q, want orders", got)
	}
	if got := awssdk.ToString(in.RoleArn); got != "arn:aws:iam::1:role/events" {
		t.Errorf("RoleArn: got %q", got)
	}
}

// AWS rejects a rule with neither a pattern nor a schedule, so awless should too rather
// than making a call that cannot succeed.
func TestCreateEventruleNeedsPatternOrSchedule(t *testing.T) {
	err := Template("create eventrule name=empty").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected a rule with neither pattern nor schedule to be rejected")
	}
}

func TestUpdateEventrule(t *testing.T) {
	mock := NewMock().On("PutRule", &eventbridge.PutRuleOutput{})

	Template(`update eventrule name=nightly schedule="rate(30 minutes)" state=DISABLED`).
		Mock(mock).
		ExpectCalls("PutRule").
		Run(t)

	in := mock.InputFor("PutRule").(*eventbridge.PutRuleInput)
	if got := awssdk.ToString(in.ScheduleExpression); got != "rate(30 minutes)" {
		t.Errorf("ScheduleExpression: got %q, want rate(30 minutes)", got)
	}
	if got := string(in.State); got != "DISABLED" {
		t.Errorf("State: got %q, want DISABLED", got)
	}
}

func TestDeleteEventrule(t *testing.T) {
	mock := NewMock().On("DeleteRule", &eventbridge.DeleteRuleOutput{})

	Template("delete eventrule name=nightly eventbus=orders force=true").
		Mock(mock).
		ExpectCalls("DeleteRule").
		Run(t)

	in := mock.InputFor("DeleteRule").(*eventbridge.DeleteRuleInput)
	if got := awssdk.ToString(in.Name); got != "nightly" {
		t.Errorf("Name: got %q, want nightly", got)
	}
	if got := awssdk.ToString(in.EventBusName); got != "orders" {
		t.Errorf("EventBusName: got %q, want orders", got)
	}
	if !in.Force {
		t.Error("expected Force to be true")
	}
}

func TestStartAndStopEventrule(t *testing.T) {
	t.Run("start enables the rule", func(t *testing.T) {
		mock := NewMock().On("EnableRule", &eventbridge.EnableRuleOutput{})
		Template("start eventrule name=nightly").
			Mock(mock).
			ExpectCalls("EnableRule").
			Run(t)
		in := mock.InputFor("EnableRule").(*eventbridge.EnableRuleInput)
		if got := awssdk.ToString(in.Name); got != "nightly" {
			t.Errorf("Name: got %q, want nightly", got)
		}
	})

	t.Run("stop disables the rule", func(t *testing.T) {
		mock := NewMock().On("DisableRule", &eventbridge.DisableRuleOutput{})
		Template("stop eventrule name=nightly eventbus=orders").
			Mock(mock).
			ExpectCalls("DisableRule").
			Run(t)
		in := mock.InputFor("DisableRule").(*eventbridge.DisableRuleInput)
		if got := awssdk.ToString(in.EventBusName); got != "orders" {
			t.Errorf("EventBusName: got %q, want orders", got)
		}
	})
}

// PutTargets takes a slice of structs. Several params address element zero through
// indexed paths, and the setter has to merge them into one target rather than appending
// a separate element per param — a bug of exactly that shape made `create distribution`
// silently drop its origin domain.
func TestAttachEventtargetMergesIndexedParams(t *testing.T) {
	mock := NewMock().On("PutTargets", &eventbridge.PutTargetsOutput{})

	Template(`attach eventtarget rule=nightly id=report-lambda ` +
		`arn=arn:aws:lambda:us-west-2:1:function:report role=arn:aws:iam::1:role/events input="{}"`).
		Mock(mock).
		ExpectCalls("PutTargets").
		ExpectCommandResult("report-lambda").
		Run(t)

	in := mock.InputFor("PutTargets").(*eventbridge.PutTargetsInput)
	if got := awssdk.ToString(in.Rule); got != "nightly" {
		t.Errorf("Rule: got %q, want nightly", got)
	}
	if len(in.Targets) != 1 {
		t.Fatalf("expected the params to merge into 1 target, got %d", len(in.Targets))
	}
	tgt := in.Targets[0]
	if got := awssdk.ToString(tgt.Id); got != "report-lambda" {
		t.Errorf("Targets[0].Id: got %q, want report-lambda", got)
	}
	if got := awssdk.ToString(tgt.Arn); got != "arn:aws:lambda:us-west-2:1:function:report" {
		t.Errorf("Targets[0].Arn: got %q", got)
	}
	if got := awssdk.ToString(tgt.RoleArn); got != "arn:aws:iam::1:role/events" {
		t.Errorf("Targets[0].RoleArn: got %q", got)
	}
	if got := awssdk.ToString(tgt.Input); got != "{}" {
		t.Errorf("Targets[0].Input: got %q, want {}", got)
	}
}

func TestDetachEventtarget(t *testing.T) {
	mock := NewMock().On("RemoveTargets", &eventbridge.RemoveTargetsOutput{})

	Template("detach eventtarget rule=nightly id=report-lambda,other-target eventbus=orders").
		Mock(mock).
		ExpectCalls("RemoveTargets").
		Run(t)

	in := mock.InputFor("RemoveTargets").(*eventbridge.RemoveTargetsInput)
	if got := awssdk.ToString(in.Rule); got != "nightly" {
		t.Errorf("Rule: got %q, want nightly", got)
	}
	if len(in.Ids) != 2 || in.Ids[0] != "report-lambda" || in.Ids[1] != "other-target" {
		t.Errorf("Ids: got %v, want [report-lambda other-target]", in.Ids)
	}
}

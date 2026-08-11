package awsat

import (
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	gatypes "github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
)

// The idempotency token is required by the API but is machinery rather than a decision, so it
// is generated. Without this the call fails validation and the user has nothing useful to
// supply.
func TestCreateAcceleratorGeneratesAnIdempotencyToken(t *testing.T) {
	mock := NewMock().On("CreateAccelerator", &globalaccelerator.CreateAcceleratorOutput{
		Accelerator: &gatypes.Accelerator{
			AcceleratorArn: awssdk.String("arn:aws:globalaccelerator::1:accelerator/abcd"),
		},
	})

	Template("create accelerator name=global-api enabled=true ip-type=IPV4").
		Mock(mock).
		ExpectCalls("CreateAccelerator").
		ExpectCommandResult("arn:aws:globalaccelerator::1:accelerator/abcd").
		ExpectRevert("delete accelerator arn=arn:aws:globalaccelerator::1:accelerator/abcd").
		Run(t)

	in := mock.InputFor("CreateAccelerator").(*globalaccelerator.CreateAcceleratorInput)
	if got := awssdk.ToString(in.Name); got != "global-api" {
		t.Errorf("Name: got %q, want global-api", got)
	}
	if !awssdk.ToBool(in.Enabled) {
		t.Error("expected Enabled to be true")
	}
	token := awssdk.ToString(in.IdempotencyToken)
	if token == "" {
		t.Error("IdempotencyToken was not generated; the API requires it")
	}
	if !strings.HasPrefix(token, "awless-") {
		t.Errorf("IdempotencyToken should be recognizably ours, got %q", token)
	}
}

// An explicit token must win, so a caller retrying deliberately can reuse one.
func TestCreateAcceleratorKeepsAnExplicitToken(t *testing.T) {
	mock := NewMock().On("CreateAccelerator", &globalaccelerator.CreateAcceleratorOutput{})

	Template("create accelerator name=global-api token=my-own-token").
		Mock(mock).
		ExpectCalls("CreateAccelerator").
		ExpectCommandResult("global-api").
		Run(t)

	in := mock.InputFor("CreateAccelerator").(*globalaccelerator.CreateAcceleratorInput)
	if got := awssdk.ToString(in.IdempotencyToken); got != "my-own-token" {
		t.Errorf("IdempotencyToken: got %q, want my-own-token", got)
	}
}

// Disabling is the prerequisite for deleting, so update has to be able to do just that.
func TestUpdateAcceleratorCanDisable(t *testing.T) {
	arn := "arn:aws:globalaccelerator::1:accelerator/abcd"
	mock := NewMock().On("UpdateAccelerator", &globalaccelerator.UpdateAcceleratorOutput{})

	Template("update accelerator arn=" + arn + " enabled=false").
		Mock(mock).
		ExpectCalls("UpdateAccelerator").
		Run(t)

	in := mock.InputFor("UpdateAccelerator").(*globalaccelerator.UpdateAcceleratorInput)
	if awssdk.ToBool(in.Enabled) {
		t.Error("expected Enabled to be false")
	}
	if got := awssdk.ToString(in.AcceleratorArn); got != arn {
		t.Errorf("AcceleratorArn: got %q", got)
	}
}

func TestUpdateAcceleratorNeedsSomethingToChange(t *testing.T) {
	err := Template("update accelerator arn=arn:aws:globalaccelerator::1:accelerator/abcd").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected an update with nothing to change to be rejected")
	}
}

func TestDeleteAccelerator(t *testing.T) {
	arn := "arn:aws:globalaccelerator::1:accelerator/abcd"
	mock := NewMock().On("DeleteAccelerator", &globalaccelerator.DeleteAcceleratorOutput{})

	Template("delete accelerator arn=" + arn).
		Mock(mock).
		ExpectCalls("DeleteAccelerator").
		Run(t)

	in := mock.InputFor("DeleteAccelerator").(*globalaccelerator.DeleteAcceleratorInput)
	if got := awssdk.ToString(in.AcceleratorArn); got != arn {
		t.Errorf("AcceleratorArn: got %q", got)
	}
}

// Port ranges are a list of from/to pairs, which is why they come from a file.
func TestCreateAcceleratorlistener(t *testing.T) {
	ports := docFile(t, "ports.json", `[{"fromPort": 80, "toPort": 80}, {"fromPort": 443, "toPort": 443}]`)

	mock := NewMock().On("CreateListener", &globalaccelerator.CreateListenerOutput{
		Listener: &gatypes.Listener{
			ListenerArn: awssdk.String("arn:aws:globalaccelerator::1:accelerator/abcd/listener/1234"),
		},
	})

	Template("create acceleratorlistener accelerator=arn:aws:globalaccelerator::1:accelerator/abcd " +
		"protocol=TCP client-affinity=SOURCE_IP ports-file=" + ports).
		Mock(mock).
		ExpectCalls("CreateListener").
		ExpectCommandResult("arn:aws:globalaccelerator::1:accelerator/abcd/listener/1234").
		Run(t)

	in := mock.InputFor("CreateListener").(*globalaccelerator.CreateListenerInput)
	if got := string(in.Protocol); got != "TCP" {
		t.Errorf("Protocol: got %q, want TCP", got)
	}
	if got := string(in.ClientAffinity); got != "SOURCE_IP" {
		t.Errorf("ClientAffinity: got %q", got)
	}
	if len(in.PortRanges) != 2 {
		t.Fatalf("PortRanges: got %d, want 2", len(in.PortRanges))
	}
	if got := awssdk.ToInt32(in.PortRanges[1].FromPort); got != 443 {
		t.Errorf("PortRanges[1].FromPort: got %d, want 443", got)
	}
}

func TestDeleteAcceleratorlistener(t *testing.T) {
	arn := "arn:aws:globalaccelerator::1:accelerator/abcd/listener/1234"
	mock := NewMock().On("DeleteListener", &globalaccelerator.DeleteListenerOutput{})

	Template("delete acceleratorlistener arn=" + arn).
		Mock(mock).
		ExpectCalls("DeleteListener").
		Run(t)

	in := mock.InputFor("DeleteListener").(*globalaccelerator.DeleteListenerInput)
	if got := awssdk.ToString(in.ListenerArn); got != arn {
		t.Errorf("ListenerArn: got %q", got)
	}
}

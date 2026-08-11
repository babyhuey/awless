package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mq"
)

func TestCreateBroker(t *testing.T) {
	users := docFile(t, "users.json",
		`[{"username": "admin", "password": "sup3rSecretPassw0rd", "consoleAccess": true}]`)

	mock := NewMock().On("CreateBroker", &mq.CreateBrokerOutput{
		BrokerId: awssdk.String("b-1234abcd"),
	})

	Template("create broker name=orders engine=RABBITMQ type=mq.m5.large mode=SINGLE_INSTANCE " +
		"public=false engine-version=3.13 subnets=subnet-1 securitygroups=sg-1234 users-file=" + users).
		Mock(mock).
		ExpectCalls("CreateBroker").
		ExpectCommandResult("b-1234abcd").
		ExpectRevert("delete broker id=b-1234abcd").
		Run(t)

	in := mock.InputFor("CreateBroker").(*mq.CreateBrokerInput)
	if got := awssdk.ToString(in.BrokerName); got != "orders" {
		t.Errorf("BrokerName: got %q, want orders", got)
	}
	if got := string(in.EngineType); got != "RABBITMQ" {
		t.Errorf("EngineType: got %q, want RABBITMQ", got)
	}
	if got := string(in.DeploymentMode); got != "SINGLE_INSTANCE" {
		t.Errorf("DeploymentMode: got %q", got)
	}
	if awssdk.ToBool(in.PubliclyAccessible) {
		t.Error("expected PubliclyAccessible to be false")
	}
	// The users list is a document, and it must arrive as a slice with the credentials
	// intact — an empty Users would leave the broker unusable.
	if len(in.Users) != 1 {
		t.Fatalf("Users: got %d, want 1", len(in.Users))
	}
	if got := awssdk.ToString(in.Users[0].Username); got != "admin" {
		t.Errorf("Users[0].Username: got %q, want admin", got)
	}
	if !awssdk.ToBool(in.Users[0].ConsoleAccess) {
		t.Error("expected Users[0].ConsoleAccess to be true")
	}
}

// Engine, type, mode, public and the users are all required: the mode decides how many
// subnets are needed and whether the broker is highly available, and there is no safe
// default for internet exposure or for the initial credentials.
func TestCreateBrokerRequiresTheEssentials(t *testing.T) {
	err := Template("create broker name=orders engine=RABBITMQ").
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected type, mode, public and users-file to be required")
	}
}

func TestDeleteBroker(t *testing.T) {
	mock := NewMock().On("DeleteBroker", &mq.DeleteBrokerOutput{})

	Template("delete broker id=b-1234abcd").
		Mock(mock).
		ExpectCalls("DeleteBroker").
		Run(t)

	in := mock.InputFor("DeleteBroker").(*mq.DeleteBrokerInput)
	if got := awssdk.ToString(in.BrokerId); got != "b-1234abcd" {
		t.Errorf("BrokerId: got %q, want b-1234abcd", got)
	}
}

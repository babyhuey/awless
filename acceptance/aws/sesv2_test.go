package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func TestCreateEmailidentity(t *testing.T) {
	mock := NewMock().On("CreateEmailIdentity", &sesv2.CreateEmailIdentityOutput{})

	Template("create emailidentity name=mail.example.com configuration-set=transactional").
		Mock(mock).
		ExpectCalls("CreateEmailIdentity").
		// The API returns a verification status, not the name, and the name is what
		// delete and the graph key on.
		ExpectCommandResult("mail.example.com").
		ExpectRevert("delete emailidentity name=mail.example.com").
		Run(t)

	in := mock.InputFor("CreateEmailIdentity").(*sesv2.CreateEmailIdentityInput)
	if got := awssdk.ToString(in.EmailIdentity); got != "mail.example.com" {
		t.Errorf("EmailIdentity: got %q, want mail.example.com", got)
	}
	if got := awssdk.ToString(in.ConfigurationSetName); got != "transactional" {
		t.Errorf("ConfigurationSetName: got %q, want transactional", got)
	}
}

func TestDeleteEmailidentity(t *testing.T) {
	mock := NewMock().On("DeleteEmailIdentity", &sesv2.DeleteEmailIdentityOutput{})

	Template("delete emailidentity name=mail.example.com").
		Mock(mock).
		ExpectCalls("DeleteEmailIdentity").
		Run(t)

	in := mock.InputFor("DeleteEmailIdentity").(*sesv2.DeleteEmailIdentityInput)
	if got := awssdk.ToString(in.EmailIdentity); got != "mail.example.com" {
		t.Errorf("EmailIdentity: got %q", got)
	}
}

func TestCreateConfigurationset(t *testing.T) {
	mock := NewMock().On("CreateConfigurationSet", &sesv2.CreateConfigurationSetOutput{})

	Template("create configurationset name=transactional").
		Mock(mock).
		ExpectCalls("CreateConfigurationSet").
		ExpectCommandResult("transactional").
		ExpectRevert("delete configurationset name=transactional").
		Run(t)

	in := mock.InputFor("CreateConfigurationSet").(*sesv2.CreateConfigurationSetInput)
	if got := awssdk.ToString(in.ConfigurationSetName); got != "transactional" {
		t.Errorf("ConfigurationSetName: got %q, want transactional", got)
	}
	// Nothing was said about the option blocks, so none should be sent — an empty
	// SendingOptions would switch sending off rather than leave it alone.
	if in.SendingOptions != nil || in.DeliveryOptions != nil {
		t.Error("option structs should be unset when no file was given")
	}
}

// The option blocks are documents, and this checks one reaches the request intact.
func TestCreateConfigurationsetWithOptions(t *testing.T) {
	delivery := docFile(t, "delivery.json", `{"tlsPolicy": "REQUIRE", "sendingPoolName": "dedicated"}`)

	mock := NewMock().On("CreateConfigurationSet", &sesv2.CreateConfigurationSetOutput{})

	Template("create configurationset name=transactional delivery-file=" + delivery).
		Mock(mock).
		ExpectCalls("CreateConfigurationSet").
		Run(t)

	in := mock.InputFor("CreateConfigurationSet").(*sesv2.CreateConfigurationSetInput)
	if in.DeliveryOptions == nil {
		t.Fatal("DeliveryOptions was not decoded")
	}
	if got := string(in.DeliveryOptions.TlsPolicy); got != "REQUIRE" {
		t.Errorf("DeliveryOptions.TlsPolicy: got %q, want REQUIRE", got)
	}
	if got := awssdk.ToString(in.DeliveryOptions.SendingPoolName); got != "dedicated" {
		t.Errorf("DeliveryOptions.SendingPoolName: got %q", got)
	}
}

func TestDeleteConfigurationset(t *testing.T) {
	mock := NewMock().On("DeleteConfigurationSet", &sesv2.DeleteConfigurationSetOutput{})

	Template("delete configurationset name=transactional").
		Mock(mock).
		ExpectCalls("DeleteConfigurationSet").
		Run(t)

	in := mock.InputFor("DeleteConfigurationSet").(*sesv2.DeleteConfigurationSetInput)
	if got := awssdk.ToString(in.ConfigurationSetName); got != "transactional" {
		t.Errorf("ConfigurationSetName: got %q", got)
	}
}

package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func TestCreateSecret(t *testing.T) {
	arn := "arn:aws:secretsmanager:us-west-2:123456789012:secret:db-password-AbCdEf"
	mock := NewMock().On("CreateSecret", &secretsmanager.CreateSecretOutput{
		ARN: awssdk.String(arn),
	})

	Template(`create secret name=db-password secret=s3cret description="db creds"`).
		Mock(mock).
		ExpectCalls("CreateSecret").
		ExpectCommandResult(arn).
		Run(t)

	in := mock.InputFor("CreateSecret").(*secretsmanager.CreateSecretInput)
	if got := awssdk.ToString(in.Name); got != "db-password" {
		t.Errorf("Name: got %q, want db-password", got)
	}
	if got := awssdk.ToString(in.SecretString); got != "s3cret" {
		t.Errorf("SecretString: got %q, want s3cret", got)
	}
}

func TestUpdateSecret(t *testing.T) {
	arn := "arn:aws:secretsmanager:us-west-2:123456789012:secret:db-password-AbCdEf"
	mock := NewMock().On("UpdateSecret", &secretsmanager.UpdateSecretOutput{
		ARN: awssdk.String(arn),
	})

	Template("update secret id=db-password secret=rotated").
		Mock(mock).
		ExpectCalls("UpdateSecret").
		Run(t)

	in := mock.InputFor("UpdateSecret").(*secretsmanager.UpdateSecretInput)
	if got := awssdk.ToString(in.SecretString); got != "rotated" {
		t.Errorf("SecretString: got %q, want rotated", got)
	}
}

// Deletion schedules removal; the recovery window must reach the API as given.
func TestDeleteSecretWithRecoveryWindow(t *testing.T) {
	mock := NewMock().On("DeleteSecret", &secretsmanager.DeleteSecretOutput{})

	Template("delete secret id=db-password recovery-window=7").
		Mock(mock).
		ExpectCalls("DeleteSecret").
		Run(t)

	in := mock.InputFor("DeleteSecret").(*secretsmanager.DeleteSecretInput)
	if got := awssdk.ToInt64(in.RecoveryWindowInDays); got != 7 {
		t.Errorf("RecoveryWindowInDays: got %d, want 7", got)
	}
}

func TestCreateSsmparameter(t *testing.T) {
	mock := NewMock().On("PutParameter", &ssm.PutParameterOutput{Version: 1})

	Template("create ssmparameter name=/app/db/host value=db.internal type=String").
		Mock(mock).
		ExpectCalls("PutParameter").
		ExpectCommandResult("/app/db/host").
		Run(t)

	in := mock.InputFor("PutParameter").(*ssm.PutParameterInput)
	if got := awssdk.ToString(in.Value); got != "db.internal" {
		t.Errorf("Value: got %q, want db.internal", got)
	}
	if got := string(in.Type); got != "String" {
		t.Errorf("Type: got %q, want String", got)
	}
}

func TestCreateSsmparameterSecureString(t *testing.T) {
	mock := NewMock().On("PutParameter", &ssm.PutParameterOutput{Version: 1})

	Template("create ssmparameter name=/app/db/pass value=hunter2 type=SecureString").
		Mock(mock).
		ExpectCalls("PutParameter").
		Run(t)

	in := mock.InputFor("PutParameter").(*ssm.PutParameterInput)
	if got := string(in.Type); got != "SecureString" {
		t.Errorf("Type: got %q, want SecureString", got)
	}
}

func TestDeleteSsmparameter(t *testing.T) {
	mock := NewMock().On("DeleteParameter", &ssm.DeleteParameterOutput{})

	Template("delete ssmparameter name=/app/db/host").
		Mock(mock).
		ExpectCalls("DeleteParameter").
		Run(t)

	in := mock.InputFor("DeleteParameter").(*ssm.DeleteParameterInput)
	if got := awssdk.ToString(in.Name); got != "/app/db/host" {
		t.Errorf("Name: got %q, want /app/db/host", got)
	}
}

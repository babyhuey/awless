package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
)

func TestCreateBackupvault(t *testing.T) {
	mock := NewMock().On("CreateBackupVault", &backup.CreateBackupVaultOutput{
		BackupVaultName: awssdk.String("daily"),
	})

	Template("create backupvault name=daily kms-key=arn:aws:kms:us-west-2:1:key/abcd").
		Mock(mock).
		ExpectCalls("CreateBackupVault").
		ExpectCommandResult("daily").
		ExpectRevert("delete backupvault name=daily").
		Run(t)

	in := mock.InputFor("CreateBackupVault").(*backup.CreateBackupVaultInput)
	if got := awssdk.ToString(in.BackupVaultName); got != "daily" {
		t.Errorf("BackupVaultName: got %q, want daily", got)
	}
	if got := awssdk.ToString(in.EncryptionKeyArn); got != "arn:aws:kms:us-west-2:1:key/abcd" {
		t.Errorf("EncryptionKeyArn: got %q", got)
	}
}

// Without a key the vault uses the AWS-managed one, so nothing should be invented.
func TestCreateBackupvaultWithoutAKey(t *testing.T) {
	mock := NewMock().On("CreateBackupVault", &backup.CreateBackupVaultOutput{})

	Template("create backupvault name=daily").
		Mock(mock).
		ExpectCalls("CreateBackupVault").
		ExpectCommandResult("daily").
		Run(t)

	in := mock.InputFor("CreateBackupVault").(*backup.CreateBackupVaultInput)
	if in.EncryptionKeyArn != nil {
		t.Error("EncryptionKeyArn should be unset when no key was given")
	}
}

func TestDeleteBackupvault(t *testing.T) {
	mock := NewMock().On("DeleteBackupVault", &backup.DeleteBackupVaultOutput{})

	Template("delete backupvault name=daily").
		Mock(mock).
		ExpectCalls("DeleteBackupVault").
		Run(t)

	in := mock.InputFor("DeleteBackupVault").(*backup.DeleteBackupVaultInput)
	if got := awssdk.ToString(in.BackupVaultName); got != "daily" {
		t.Errorf("BackupVaultName: got %q", got)
	}
}

// A plan is a list of rules each with its own schedule, lifecycle and target vault, which is
// why it comes from a file. This asserts the whole structure reached the request.
func TestCreateBackupplanFromADocument(t *testing.T) {
	plan := docFile(t, "plan.json", `{
	  "backupPlanName": "daily-and-weekly",
	  "rules": [
	    {"ruleName": "daily", "targetBackupVaultName": "daily",
	     "scheduleExpression": "cron(0 5 * * ? *)",
	     "lifecycle": {"deleteAfterDays": 35}},
	    {"ruleName": "weekly", "targetBackupVaultName": "compliance",
	     "scheduleExpression": "cron(0 5 ? * SUN *)",
	     "lifecycle": {"deleteAfterDays": 365}}
	  ]
	}`)

	mock := NewMock().On("CreateBackupPlan", &backup.CreateBackupPlanOutput{
		BackupPlanId: awssdk.String("abcd1234-5678"),
	})

	Template("create backupplan definition-file=" + plan).
		Mock(mock).
		ExpectCalls("CreateBackupPlan").
		ExpectCommandResult("abcd1234-5678").
		Run(t)

	in := mock.InputFor("CreateBackupPlan").(*backup.CreateBackupPlanInput)
	if in.BackupPlan == nil {
		t.Fatal("the plan was not decoded")
	}
	if got := awssdk.ToString(in.BackupPlan.BackupPlanName); got != "daily-and-weekly" {
		t.Errorf("BackupPlanName: got %q", got)
	}
	if len(in.BackupPlan.Rules) != 2 {
		t.Fatalf("Rules: got %d, want 2", len(in.BackupPlan.Rules))
	}
	if got := awssdk.ToString(in.BackupPlan.Rules[1].RuleName); got != "weekly" {
		t.Errorf("Rules[1].RuleName: got %q, want weekly", got)
	}
	if got := awssdk.ToString(in.BackupPlan.Rules[1].TargetBackupVaultName); got != "compliance" {
		t.Errorf("Rules[1].TargetBackupVaultName: got %q", got)
	}
	// A lifecycle nested inside a rule inside the plan.
	if in.BackupPlan.Rules[1].Lifecycle == nil {
		t.Fatal("the rule lifecycle was not decoded")
	}
	if got := awssdk.ToInt64(in.BackupPlan.Rules[1].Lifecycle.DeleteAfterDays); got != 365 {
		t.Errorf("Rules[1].Lifecycle.DeleteAfterDays: got %d, want 365", got)
	}
}

func TestDeleteBackupplan(t *testing.T) {
	mock := NewMock().On("DeleteBackupPlan", &backup.DeleteBackupPlanOutput{})

	Template("delete backupplan id=abcd1234-5678").
		Mock(mock).
		ExpectCalls("DeleteBackupPlan").
		Run(t)

	in := mock.InputFor("DeleteBackupPlan").(*backup.DeleteBackupPlanInput)
	if got := awssdk.ToString(in.BackupPlanId); got != "abcd1234-5678" {
		t.Errorf("BackupPlanId: got %q", got)
	}
}

/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsspec

import (
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// A vault is where recovery points are stored; a plan says what to back up, when, and into
// which vault. A plan is a document — a list of rules each with its own schedule, lifecycle
// and target vault — so it comes from a file.

type CreateBackupvault struct {
	_      string `action:"create" entity:"backupvault" awsAPI:"backup" awsCall:"CreateBackupVault" awsInput:"backup.CreateBackupVaultInput" awsOutput:"backup.CreateBackupVaultOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *backup.Client
	Name   *string `awsName:"BackupVaultName" awsType:"awsstr" templateName:"name"`
	// Without a key the vault uses the AWS-managed one. A customer-managed key is what
	// makes a cross-account or cross-region copy possible.
	KmsKey *string `awsName:"EncryptionKeyArn" awsType:"awsstr" templateName:"kms-key"`
}

func (cmd *CreateBackupvault) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("kms-key"),
	))
}

func (cmd *CreateBackupvault) ExtractResult(i any) string {
	out, ok := i.(*backup.CreateBackupVaultOutput)
	if !ok || out.BackupVaultName == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.BackupVaultName)
}

// A vault holding recovery points cannot be deleted, which AWS enforces — so this does not
// destroy backups by accident.
type DeleteBackupvault struct {
	_      string `action:"delete" entity:"backupvault" awsAPI:"backup" awsCall:"DeleteBackupVault" awsInput:"backup.DeleteBackupVaultInput" awsOutput:"backup.DeleteBackupVaultOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *backup.Client
	Name   *string `awsName:"BackupVaultName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteBackupvault) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type CreateBackupplan struct {
	_      string `action:"create" entity:"backupplan" awsAPI:"backup" awsCall:"CreateBackupPlan" awsInput:"backup.CreateBackupPlanInput" awsOutput:"backup.CreateBackupPlanOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *backup.Client
	// The plan carries its own name, so there is no separate name param to contradict it.
	DefinitionFile *string `awsName:"BackupPlan" awsType:"awsfiletostruct" templateName:"definition-file"`
}

func (cmd *CreateBackupplan) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("definition-file")))
}

func (cmd *CreateBackupplan) ExtractResult(i any) string {
	out, ok := i.(*backup.CreateBackupPlanOutput)
	if !ok || out.BackupPlanId == nil {
		return ""
	}
	return awssdk.ToString(out.BackupPlanId)
}

type DeleteBackupplan struct {
	_      string `action:"delete" entity:"backupplan" awsAPI:"backup" awsCall:"DeleteBackupPlan" awsInput:"backup.DeleteBackupPlanInput" awsOutput:"backup.DeleteBackupPlanOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *backup.Client
	ID     *string `awsName:"BackupPlanId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteBackupplan) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

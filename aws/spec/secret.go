/* Copyright 2017 WALLIX

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
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateSecret struct {
	_           string `action:"create" entity:"secret" awsAPI:"secretsmanager" awsCall:"CreateSecret" awsInput:"secretsmanager.CreateSecretInput" awsOutput:"secretsmanager.CreateSecretOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *secretsmanager.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Secret      *string `awsName:"SecretString" awsType:"awsstr" templateName:"secret"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	KmsKey      *string `awsName:"KmsKeyId" awsType:"awsstr" templateName:"kms-key"`
}

func (cmd *CreateSecret) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("secret"),
		params.Opt("description", "kms-key"),
	))
}

func (cmd *CreateSecret) ExtractResult(i any) string {
	return awssdk.ToString(i.(*secretsmanager.CreateSecretOutput).ARN)
}

type UpdateSecret struct {
	_           string `action:"update" entity:"secret" awsAPI:"secretsmanager" awsCall:"UpdateSecret" awsInput:"secretsmanager.UpdateSecretInput" awsOutput:"secretsmanager.UpdateSecretOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *secretsmanager.Client
	Id          *string `awsName:"SecretId" awsType:"awsstr" templateName:"id"`
	Secret      *string `awsName:"SecretString" awsType:"awsstr" templateName:"secret"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	KmsKey      *string `awsName:"KmsKeyId" awsType:"awsstr" templateName:"kms-key"`
}

func (cmd *UpdateSecret) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		params.Opt("secret", "description", "kms-key"),
	))
}

func (cmd *UpdateSecret) ExtractResult(i any) string {
	return awssdk.ToString(i.(*secretsmanager.UpdateSecretOutput).ARN)
}

// DeleteSecret schedules deletion rather than deleting immediately, because that
// is the only behavior the API offers without force-delete. recovery-window
// defaults to the API's 30 days; force=true removes the secret with no recovery
// window and cannot be undone.
type DeleteSecret struct {
	_              string `action:"delete" entity:"secret" awsAPI:"secretsmanager" awsCall:"DeleteSecret" awsInput:"secretsmanager.DeleteSecretInput" awsOutput:"secretsmanager.DeleteSecretOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            *secretsmanager.Client
	Id             *string `awsName:"SecretId" awsType:"awsstr" templateName:"id"`
	RecoveryWindow *int64  `awsName:"RecoveryWindowInDays" awsType:"awsint64" templateName:"recovery-window"`
	Force          *bool   `awsName:"ForceDeleteWithoutRecovery" awsType:"awsbool" templateName:"force"`
}

func (cmd *DeleteSecret) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("id"),
		// The API rejects a request carrying both, so they are mutually exclusive
		// here rather than merely optional.
		params.Opt("recovery-window", "force"),
	))
}

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
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateFilesystem struct {
	_               string `action:"create" entity:"filesystem" awsAPI:"efs" awsCall:"CreateFileSystem" awsInput:"efs.CreateFileSystemInput" awsOutput:"efs.CreateFileSystemOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             *efs.Client
	Token           *string `awsName:"CreationToken" awsType:"awsstr" templateName:"token"`
	PerformanceMode *string `awsName:"PerformanceMode" awsType:"awsstr" templateName:"performance-mode"`
	ThroughputMode  *string `awsName:"ThroughputMode" awsType:"awsstr" templateName:"throughput-mode"`
	Encrypted       *bool   `awsName:"Encrypted" awsType:"awsbool" templateName:"encrypted"`
	KmsKey          *string `awsName:"KmsKeyId" awsType:"awsstr" templateName:"kms-key"`
}

func (cmd *CreateFilesystem) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(
			// CreationToken is the API's idempotency key and is required.
			params.Key("token"),
			params.Opt(
				params.Suggested("encrypted"),
				"performance-mode", "throughput-mode", "kms-key",
			),
		),
		params.Validators{
			"performance-mode": isEfsPerformanceMode,
			"throughput-mode":  isEfsThroughputMode,
		})
}

func (cmd *CreateFilesystem) ExtractResult(i any) string {
	return awssdk.ToString(i.(*efs.CreateFileSystemOutput).FileSystemId)
}

type DeleteFilesystem struct {
	_      string `action:"delete" entity:"filesystem" awsAPI:"efs" awsCall:"DeleteFileSystem" awsInput:"efs.DeleteFileSystemInput" awsOutput:"efs.DeleteFileSystemOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *efs.Client
	Id     *string `awsName:"FileSystemId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteFilesystem) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

func isEfsPerformanceMode(i any, _ map[string]any) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected a string, got %T", i)
	}
	switch s {
	case "generalPurpose", "maxIO":
		return nil
	}
	return fmt.Errorf("invalid performance mode '%s', expected generalPurpose or maxIO", s)
}

func isEfsThroughputMode(i any, _ map[string]any) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected a string, got %T", i)
	}
	switch s {
	case "bursting", "provisioned", "elastic":
		return nil
	}
	return fmt.Errorf("invalid throughput mode '%s', expected bursting, provisioned or elastic", s)
}

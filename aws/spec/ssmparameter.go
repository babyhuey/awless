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

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// CreateSsmparameter writes an SSM parameter. type defaults to String; use
// SecureString for values that should be KMS-encrypted at rest.
type CreateSsmparameter struct {
	_           string `action:"create" entity:"ssmparameter" awsAPI:"ssm" awsCall:"PutParameter" awsInput:"ssm.PutParameterInput" awsOutput:"ssm.PutParameterOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *ssm.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Value       *string `awsName:"Value" awsType:"awsstr" templateName:"value"`
	Type        *string `awsName:"Type" awsType:"awsstr" templateName:"type"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	KmsKey      *string `awsName:"KeyId" awsType:"awsstr" templateName:"kms-key"`
	Tier        *string `awsName:"Tier" awsType:"awsstr" templateName:"tier"`
}

func (cmd *CreateSsmparameter) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(
			params.Key("name"), params.Key("value"),
			params.Opt(params.Suggested("type"), "description", "kms-key", "tier"),
		),
		params.Validators{"type": isSsmParameterType})
}

func (cmd *CreateSsmparameter) ExtractResult(i any) string {
	// PutParameter returns a version, not an identifier; the name is what every
	// other command and the local graph key on.
	return StringValue(cmd.Name)
}

// UpdateSsmparameter is PutParameter with Overwrite set, which is how the API
// models an update — there is no separate operation.
type UpdateSsmparameter struct {
	_           string `action:"update" entity:"ssmparameter" awsAPI:"ssm" awsCall:"PutParameter" awsInput:"ssm.PutParameterInput" awsOutput:"ssm.PutParameterOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *ssm.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Value       *string `awsName:"Value" awsType:"awsstr" templateName:"value"`
	Type        *string `awsName:"Type" awsType:"awsstr" templateName:"type"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	KmsKey      *string `awsName:"KeyId" awsType:"awsstr" templateName:"kms-key"`
	Overwrite   *bool   `awsName:"Overwrite" awsType:"awsbool" templateName:"overwrite"`
}

func (cmd *UpdateSsmparameter) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(
			params.Key("name"), params.Key("value"),
			params.Opt("type", "description", "kms-key", "overwrite"),
		),
		params.Validators{"type": isSsmParameterType})
}

func (cmd *UpdateSsmparameter) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteSsmparameter struct {
	_      string `action:"delete" entity:"ssmparameter" awsAPI:"ssm" awsCall:"DeleteParameter" awsInput:"ssm.DeleteParameterInput" awsOutput:"ssm.DeleteParameterOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *ssm.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteSsmparameter) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// isSsmParameterType rejects an invalid type before the call, since the API error
// for it is opaque.
func isSsmParameterType(i any, _ map[string]any) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected a string, got %T", i)
	}
	switch s {
	case "String", "StringList", "SecureString":
		return nil
	}
	return fmt.Errorf("invalid parameter type '%s', expected String, StringList or SecureString", s)
}

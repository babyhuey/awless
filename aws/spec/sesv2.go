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
	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// An email identity is a domain or a single address that SES will send from. Creating one
// starts verification rather than completing it: a domain needs DNS records published, an
// address needs a link clicked, so the identity exists in a pending state until then.
type CreateEmailidentity struct {
	_      string `action:"create" entity:"emailidentity" awsAPI:"sesv2" awsCall:"CreateEmailIdentity" awsInput:"sesv2.CreateEmailIdentityInput" awsOutput:"sesv2.CreateEmailIdentityOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sesv2.Client
	Name   *string `awsName:"EmailIdentity" awsType:"awsstr" templateName:"name"`
	// Attaching a configuration set at creation is the only way to have it apply to the
	// identity's first sends.
	ConfigurationSet *string `awsName:"ConfigurationSetName" awsType:"awsstr" templateName:"configuration-set"`
}

func (cmd *CreateEmailidentity) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("configuration-set"),
	))
}

// CreateEmailIdentity returns the verification status, not the name.
func (cmd *CreateEmailidentity) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteEmailidentity struct {
	_      string `action:"delete" entity:"emailidentity" awsAPI:"sesv2" awsCall:"DeleteEmailIdentity" awsInput:"sesv2.DeleteEmailIdentityInput" awsOutput:"sesv2.DeleteEmailIdentityOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sesv2.Client
	Name   *string `awsName:"EmailIdentity" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteEmailidentity) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type CreateConfigurationset struct {
	_      string `action:"create" entity:"configurationset" awsAPI:"sesv2" awsCall:"CreateConfigurationSet" awsInput:"sesv2.CreateConfigurationSetInput" awsOutput:"sesv2.CreateConfigurationSetOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sesv2.Client
	Name   *string `awsName:"ConfigurationSetName" awsType:"awsstr" templateName:"name"`
	// Sending, reputation, suppression and delivery options are each their own nested
	// structure, and which combination matters depends entirely on the setup, so they
	// come from files rather than being flattened into a wall of flags.
	SendingFile     *string `awsName:"SendingOptions" awsType:"awsfiletostruct" templateName:"sending-file"`
	DeliveryFile    *string `awsName:"DeliveryOptions" awsType:"awsfiletostruct" templateName:"delivery-file"`
	ReputationFile  *string `awsName:"ReputationOptions" awsType:"awsfiletostruct" templateName:"reputation-file"`
	SuppressionFile *string `awsName:"SuppressionOptions" awsType:"awsfiletostruct" templateName:"suppression-file"`
}

func (cmd *CreateConfigurationset) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("sending-file", "delivery-file", "reputation-file", "suppression-file"),
	))
}

func (cmd *CreateConfigurationset) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteConfigurationset struct {
	_      string `action:"delete" entity:"configurationset" awsAPI:"sesv2" awsCall:"DeleteConfigurationSet" awsInput:"sesv2.DeleteConfigurationSetInput" awsOutput:"sesv2.DeleteConfigurationSetOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *sesv2.Client
	Name   *string `awsName:"ConfigurationSetName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteConfigurationset) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

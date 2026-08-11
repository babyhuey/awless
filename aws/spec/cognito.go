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
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// Cognito is two APIs. A user pool is a directory of users on cognito-idp; an identity pool
// hands out AWS credentials to those users on cognito-identity. They are frequently confused
// and are not interchangeable.

type CreateUserpool struct {
	_      string `action:"create" entity:"userpool" awsAPI:"cognitoidentityprovider" awsCall:"CreateUserPool" awsInput:"cognitoidentityprovider.CreateUserPoolInput" awsOutput:"cognitoidentityprovider.CreateUserPoolOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cognitoidentityprovider.Client
	Name   *string `awsName:"PoolName" awsType:"awsstr" templateName:"name"`
	// Deletion protection defaults to active on a new pool in the console but not in the
	// API, and a user pool holds credentials that cannot be recreated, so it is worth
	// surfacing.
	DeletionProtection *string   `awsName:"DeletionProtection" awsType:"awsstr" templateName:"deletion-protection"`
	AutoVerified       []*string `awsName:"AutoVerifiedAttributes" awsType:"awsstringslice" templateName:"auto-verified"`
	AliasAttributes    []*string `awsName:"AliasAttributes" awsType:"awsstringslice" templateName:"alias-attributes"`
	// The password policy, MFA settings and recovery options are nested structures whose
	// shape depends on the setup, so they come from files.
	PoliciesFile *string `awsName:"Policies" awsType:"awsfiletostruct" templateName:"policies-file"`
	RecoveryFile *string `awsName:"AccountRecoverySetting" awsType:"awsfiletostruct" templateName:"recovery-file"`
}

func (cmd *CreateUserpool) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(
			params.Suggested("deletion-protection"),
			"auto-verified", "alias-attributes", "policies-file", "recovery-file",
		),
	))
}

func (cmd *CreateUserpool) ExtractResult(i any) string {
	out, ok := i.(*cognitoidentityprovider.CreateUserPoolOutput)
	if !ok || out.UserPool == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.UserPool.Id)
}

type DeleteUserpool struct {
	_      string `action:"delete" entity:"userpool" awsAPI:"cognitoidentityprovider" awsCall:"DeleteUserPool" awsInput:"cognitoidentityprovider.DeleteUserPoolInput" awsOutput:"cognitoidentityprovider.DeleteUserPoolOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cognitoidentityprovider.Client
	// The id, not the name: user pool names are not unique and the API takes the id.
	ID *string `awsName:"UserPoolId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteUserpool) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CreateIdentitypool struct {
	_      string `action:"create" entity:"identitypool" awsAPI:"cognitoidentity" awsCall:"CreateIdentityPool" awsInput:"cognitoidentity.CreateIdentityPoolInput" awsOutput:"cognitoidentity.CreateIdentityPoolOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cognitoidentity.Client
	Name   *string `awsName:"IdentityPoolName" awsType:"awsstr" templateName:"name"`
	// Required by the API and security-relevant: unauthenticated identities hand AWS
	// credentials to anyone who asks, so there is no safe default to pick.
	AllowUnauthenticated *bool   `awsName:"AllowUnauthenticatedIdentities" awsType:"awsbool" templateName:"allow-unauthenticated"`
	DeveloperProvider    *string `awsName:"DeveloperProviderName" awsType:"awsstr" templateName:"developer-provider"`
	AllowClassicFlow     *bool   `awsName:"AllowClassicFlow" awsType:"awsbool" templateName:"allow-classic-flow"`
}

func (cmd *CreateIdentitypool) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("allow-unauthenticated"),
		params.Opt("developer-provider", "allow-classic-flow"),
	))
}

func (cmd *CreateIdentitypool) ExtractResult(i any) string {
	out, ok := i.(*cognitoidentity.CreateIdentityPoolOutput)
	if !ok || out.IdentityPoolId == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.IdentityPoolId)
}

type DeleteIdentitypool struct {
	_      string `action:"delete" entity:"identitypool" awsAPI:"cognitoidentity" awsCall:"DeleteIdentityPool" awsInput:"cognitoidentity.DeleteIdentityPoolInput" awsOutput:"cognitoidentity.DeleteIdentityPoolOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cognitoidentity.Client
	ID     *string `awsName:"IdentityPoolId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteIdentitypool) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

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
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	apigwtypes "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateApigateway struct {
	_           string `action:"create" entity:"apigateway" awsAPI:"apigatewayv2" awsCall:"CreateApi" awsInput:"apigatewayv2.CreateApiInput" awsOutput:"apigatewayv2.CreateApiOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *apigatewayv2.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Protocol    *string `awsName:"ProtocolType" awsType:"awsstr" templateName:"protocol"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Target      *string `awsName:"Target" awsType:"awsstr" templateName:"target"`
	Version     *string `awsName:"Version" awsType:"awsstr" templateName:"version"`
	RouteKey    *string `awsName:"RouteKey" awsType:"awsstr" templateName:"route-key"`
}

func (cmd *CreateApigateway) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(
			params.Key("name"), params.Key("protocol"),
			params.Opt("description", "target", "version", "route-key"),
		),
		params.Validators{"protocol": isApigatewayProtocol})
}

func (cmd *CreateApigateway) ExtractResult(i any) string {
	return awssdk.ToString(i.(*apigatewayv2.CreateApiOutput).ApiId)
}

type DeleteApigateway struct {
	_      string `action:"delete" entity:"apigateway" awsAPI:"apigatewayv2" awsCall:"DeleteApi" awsInput:"apigatewayv2.DeleteApiInput" awsOutput:"apigatewayv2.DeleteApiOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *apigatewayv2.Client
	Id     *string `awsName:"ApiId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteApigateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CreateApigatewayroute struct {
	_          string `action:"create" entity:"apigatewayroute" awsAPI:"apigatewayv2" awsCall:"CreateRoute" awsInput:"apigatewayv2.CreateRouteInput" awsOutput:"apigatewayv2.CreateRouteOutput"`
	logger     *logger.Logger
	graph      cloud.GraphAPI
	api        *apigatewayv2.Client
	Api        *string `awsName:"ApiId" awsType:"awsstr" templateName:"api"`
	RouteKey   *string `awsName:"RouteKey" awsType:"awsstr" templateName:"route-key"`
	Target     *string `awsName:"Target" awsType:"awsstr" templateName:"target"`
	Authorizer *string `awsName:"AuthorizerId" awsType:"awsstr" templateName:"authorizer"`
	AuthType   *string `awsName:"AuthorizationType" awsType:"awsstr" templateName:"auth-type"`
}

func (cmd *CreateApigatewayroute) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("api"), params.Key("route-key"),
		params.Opt("target", "authorizer", "auth-type"),
	))
}

func (cmd *CreateApigatewayroute) ExtractResult(i any) string {
	return awssdk.ToString(i.(*apigatewayv2.CreateRouteOutput).RouteId)
}

type DeleteApigatewayroute struct {
	_      string `action:"delete" entity:"apigatewayroute" awsAPI:"apigatewayv2" awsCall:"DeleteRoute" awsInput:"apigatewayv2.DeleteRouteInput" awsOutput:"apigatewayv2.DeleteRouteOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *apigatewayv2.Client
	Api    *string `awsName:"ApiId" awsType:"awsstr" templateName:"api"`
	Id     *string `awsName:"RouteId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteApigatewayroute) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("api"), params.Key("id")))
}

type CreateApigatewaystage struct {
	_          string `action:"create" entity:"apigatewaystage" awsAPI:"apigatewayv2" awsCall:"CreateStage" awsInput:"apigatewayv2.CreateStageInput" awsOutput:"apigatewayv2.CreateStageOutput"`
	logger     *logger.Logger
	graph      cloud.GraphAPI
	api        *apigatewayv2.Client
	Api        *string `awsName:"ApiId" awsType:"awsstr" templateName:"api"`
	Name       *string `awsName:"StageName" awsType:"awsstr" templateName:"name"`
	Autodeploy *bool   `awsName:"AutoDeploy" awsType:"awsbool" templateName:"autodeploy"`
	Deployment *string `awsName:"DeploymentId" awsType:"awsstr" templateName:"deployment"`
	Descr      *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateApigatewaystage) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("api"), params.Key("name"),
		params.Opt(params.Suggested("autodeploy"), "deployment", "description"),
	))
}

func (cmd *CreateApigatewaystage) ExtractResult(i any) string {
	return awssdk.ToString(i.(*apigatewayv2.CreateStageOutput).StageName)
}

type DeleteApigatewaystage struct {
	_      string `action:"delete" entity:"apigatewaystage" awsAPI:"apigatewayv2" awsCall:"DeleteStage" awsInput:"apigatewayv2.DeleteStageInput" awsOutput:"apigatewayv2.DeleteStageOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *apigatewayv2.Client
	Api    *string `awsName:"ApiId" awsType:"awsstr" templateName:"api"`
	Name   *string `awsName:"StageName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteApigatewaystage) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("api"), params.Key("name")))
}

func isApigatewayProtocol(i any, _ map[string]any) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected a string, got %T", i)
	}
	for _, p := range apigwtypes.ProtocolTypeHttp.Values() {
		if string(p) == s {
			return nil
		}
	}
	return fmt.Errorf("invalid protocol '%s', expected HTTP or WEBSOCKET", s)
}

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
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

// Cloud Map is service discovery: a namespace holds services, and a service holds instances.
//
// The entity is discoveryservice rather than service. "service" is already overloaded here —
// ECS has containerservice, and awless itself calls its AWS integrations services — so a bare
// "service" would be ambiguous in a template.
//
// A namespace is created either as HTTP-only or backed by private DNS in a VPC. Those are
// different API calls, so the vpc param chooses between them and creation goes through
// ManualRun.
type CreateNamespace struct {
	_           string `action:"create" entity:"namespace" awsAPI:"servicediscovery"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *servicediscovery.Client
	Name        *string `templateName:"name"`
	Description *string `templateName:"description"`
	// Given a VPC, the namespace is backed by a private hosted zone and instances are
	// resolvable by DNS from inside it. Without one, discovery is HTTP-only through the
	// Cloud Map API.
	Vpc *string `templateName:"vpc"`
}

func (cmd *CreateNamespace) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("vpc"), "description"),
	))
}

func (cmd *CreateNamespace) ManualRun(renv env.Running) (any, error) {
	if vpc := StringValue(cmd.Vpc); vpc != "" {
		return cmd.api.CreatePrivateDnsNamespace(renv.RequestContext(), &servicediscovery.CreatePrivateDnsNamespaceInput{
			Name:        cmd.Name,
			Vpc:         cmd.Vpc,
			Description: cmd.Description,
		})
	}

	return cmd.api.CreateHttpNamespace(renv.RequestContext(), &servicediscovery.CreateHttpNamespaceInput{
		Name:        cmd.Name,
		Description: cmd.Description,
	})
}

// Both calls return an operation id rather than the namespace: creation is asynchronous, and
// the namespace does not exist yet when the call returns.
func (cmd *CreateNamespace) ExtractResult(i any) string {
	switch out := i.(type) {
	case *servicediscovery.CreatePrivateDnsNamespaceOutput:
		return awssdk.ToString(out.OperationId)
	case *servicediscovery.CreateHttpNamespaceOutput:
		return awssdk.ToString(out.OperationId)
	}
	return StringValue(cmd.Name)
}

type DeleteNamespace struct {
	_      string `action:"delete" entity:"namespace" awsAPI:"servicediscovery" awsCall:"DeleteNamespace" awsInput:"servicediscovery.DeleteNamespaceInput" awsOutput:"servicediscovery.DeleteNamespaceOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *servicediscovery.Client
	ID     *string `awsName:"Id" awsType:"awsstr" templateName:"id"`
}

// A namespace with services still in it cannot be deleted, which AWS enforces.
func (cmd *DeleteNamespace) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CreateDiscoveryservice struct {
	_           string `action:"create" entity:"discoveryservice" awsAPI:"servicediscovery" awsCall:"CreateService" awsInput:"servicediscovery.CreateServiceInput" awsOutput:"servicediscovery.CreateServiceOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *servicediscovery.Client
	Name        *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Namespace   *string `awsName:"NamespaceId" awsType:"awsstr" templateName:"namespace"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	// The DNS records the service registers and their TTLs, and the health check, are
	// nested structures, so they come from files.
	DNSFile         *string `awsName:"DnsConfig" awsType:"awsfiletostruct" templateName:"dns-file"`
	HealthCheckFile *string `awsName:"HealthCheckConfig" awsType:"awsfiletostruct" templateName:"healthcheck-file"`
}

func (cmd *CreateDiscoveryservice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("namespace"),
		params.Opt(params.Suggested("dns-file"), "description", "healthcheck-file"),
	))
}

func (cmd *CreateDiscoveryservice) ExtractResult(i any) string {
	out, ok := i.(*servicediscovery.CreateServiceOutput)
	if !ok || out.Service == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Service.Id)
}

type DeleteDiscoveryservice struct {
	_      string `action:"delete" entity:"discoveryservice" awsAPI:"servicediscovery" awsCall:"DeleteService" awsInput:"servicediscovery.DeleteServiceInput" awsOutput:"servicediscovery.DeleteServiceOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *servicediscovery.Client
	ID     *string `awsName:"Id" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteDiscoveryservice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

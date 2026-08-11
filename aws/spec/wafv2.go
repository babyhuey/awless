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
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

// WAF v2 is scoped: REGIONAL resources belong to one region, CLOUDFRONT ones to a global
// namespace only reachable from us-east-1. Every call needs the scope, so every command
// takes it, defaulting to REGIONAL.
//
// Deletes and updates also need a LockToken, WAF's optimistic-concurrency handle. It is
// only obtainable by listing or getting the resource first, so those commands look it up
// in ManualRun rather than making the user run a separate command to find it.

const defaultWafScope = string(wafv2types.ScopeRegional)

func wafScope(given *string) wafv2types.Scope {
	if s := StringValue(given); s != "" {
		return wafv2types.Scope(s)
	}
	return wafv2types.Scope(defaultWafScope)
}

// findIPSet returns the id and lock token for a named IP set within a scope. WAF has no
// lookup by name, so the list is walked.
func findIPSet(renv env.Running, api *wafv2.Client, name string, scope wafv2types.Scope) (id, lockToken string, err error) {
	var next *string
	for {
		out, lerr := api.ListIPSets(renv.RequestContext(), &wafv2.ListIPSetsInput{Scope: scope, NextMarker: next})
		if lerr != nil {
			return "", "", lerr
		}
		for _, set := range out.IPSets {
			if awssdk.ToString(set.Name) == name {
				return awssdk.ToString(set.Id), awssdk.ToString(set.LockToken), nil
			}
		}
		if out.NextMarker == nil || awssdk.ToString(out.NextMarker) == "" {
			return "", "", fmt.Errorf("no ip set named %q in scope %s", name, scope)
		}
		next = out.NextMarker
	}
}

type CreateIpset struct {
	_           string `action:"create" entity:"ipset" awsAPI:"wafv2"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *wafv2.Client
	Name        *string   `templateName:"name"`
	Addresses   []*string `templateName:"addresses"`
	Scope       *string   `templateName:"scope"`
	Description *string   `templateName:"description"`
	IPVersion   *string   `templateName:"ip-version"`
}

func (cmd *CreateIpset) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("addresses"),
		params.Opt(params.Suggested("scope", "ip-version"), "description"),
	))
}

// Hand-built rather than driven by awsName tags: IPAddressVersion is an enum that has to
// be derived from the addresses when the user does not give it, since a set may not mix
// IPv4 and IPv6.
func (cmd *CreateIpset) ManualRun(renv env.Running) (any, error) {
	addresses := castStringSlice(cmd.Addresses)
	version := wafv2types.IPAddressVersion(StringValue(cmd.IPVersion))
	if version == "" {
		version = wafv2types.IPAddressVersionIpv4
		for _, a := range addresses {
			// Any colon means IPv6; CIDR notation cannot otherwise contain one.
			if strings.Contains(a, ":") {
				version = wafv2types.IPAddressVersionIpv6
				break
			}
		}
	}

	input := &wafv2.CreateIPSetInput{
		Name:             cmd.Name,
		Scope:            wafScope(cmd.Scope),
		Addresses:        addresses,
		IPAddressVersion: version,
		Description:      cmd.Description,
	}

	out, err := cmd.api.CreateIPSet(renv.RequestContext(), input)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (cmd *CreateIpset) ExtractResult(i any) string {
	out, ok := i.(*wafv2.CreateIPSetOutput)
	if !ok || out.Summary == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Summary.Name)
}

type DeleteIpset struct {
	_      string `action:"delete" entity:"ipset" awsAPI:"wafv2"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *wafv2.Client
	Name   *string `templateName:"name"`
	Scope  *string `templateName:"scope"`
}

func (cmd *DeleteIpset) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("scope")),
	))
}

func (cmd *DeleteIpset) ManualRun(renv env.Running) (any, error) {
	scope := wafScope(cmd.Scope)
	id, lockToken, err := findIPSet(renv, cmd.api, StringValue(cmd.Name), scope)
	if err != nil {
		return nil, err
	}

	return cmd.api.DeleteIPSet(renv.RequestContext(), &wafv2.DeleteIPSetInput{
		Name:      cmd.Name,
		Scope:     scope,
		Id:        awssdk.String(id),
		LockToken: awssdk.String(lockToken),
	})
}

type UpdateIpset struct {
	_           string `action:"update" entity:"ipset" awsAPI:"wafv2"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *wafv2.Client
	Name        *string   `templateName:"name"`
	Addresses   []*string `templateName:"addresses"`
	Scope       *string   `templateName:"scope"`
	Description *string   `templateName:"description"`
}

// UpdateIPSet replaces the address list outright rather than merging, so addresses stays
// required — an update without it would empty the set.
func (cmd *UpdateIpset) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("addresses"),
		params.Opt(params.Suggested("scope"), "description"),
	))
}

func (cmd *UpdateIpset) ManualRun(renv env.Running) (any, error) {
	scope := wafScope(cmd.Scope)
	id, lockToken, err := findIPSet(renv, cmd.api, StringValue(cmd.Name), scope)
	if err != nil {
		return nil, err
	}

	return cmd.api.UpdateIPSet(renv.RequestContext(), &wafv2.UpdateIPSetInput{
		Name:        cmd.Name,
		Scope:       scope,
		Id:          awssdk.String(id),
		LockToken:   awssdk.String(lockToken),
		Addresses:   castStringSlice(cmd.Addresses),
		Description: cmd.Description,
	})
}

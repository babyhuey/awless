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

// DO NOT EDIT
// This file was automatically generated with go generate
package ast

// entities is the set the template parser accepts. Derived from the `entity:`
// struct tags in aws/spec, plus a small set of entities kept parseable so that
// template history already written to ~/.awless still loads. See
// gen/aws/generators/entities.go.
var entities = map[Entity]struct{}{
	"accesskey":           {},
	"alarm":               {},
	"apigateway":          {},
	"apigatewayroute":     {},
	"apigatewaystage":     {},
	"appscalingpolicy":    {},
	"appscalingtarget":    {},
	"bucket":              {},
	"cachecluster":        {},
	"cachesubnetgroup":    {},
	"certificate":         {},
	"classicloadbalancer": {},
	"container":           {},
	"containercluster":    {},
	"containerservice":    {},
	"containertask":       {},
	"database":            {},
	"dbsubnetgroup":       {},
	"distribution":        {},
	"dynamodbtable":       {},
	"ekscluster":          {},
	"eksnodegroup":        {},
	"elasticip":           {},
	"eventbus":            {},
	"eventrule":           {},
	"eventtarget":         {},
	"filesystem":          {},
	"function":            {},
	"group":               {},
	"image":               {},
	"instance":            {},
	"instanceprofile":     {},
	"internetgateway":     {},
	"keypair":             {},
	"launchconfiguration": {},
	"listener":            {},
	"loadbalancer":        {},
	"loggroup":            {},
	"loginprofile":        {},
	"mfadevice":           {},
	"natgateway":          {},
	"networkinterface":    {},
	"none":                {},
	"policy":              {},
	"queue":               {},
	"record":              {},
	"registry":            {},
	"replicationgroup":    {},
	"repository":          {},
	"role":                {},
	"route":               {},
	"routetable":          {},
	"s3object":            {},
	"scalinggroup":        {},
	"scalingpolicy":       {},
	"secret":              {},
	"securitygroup":       {},
	"snapshot":            {},
	"ssmparameter":        {},
	"stack":               {},
	"subnet":              {},
	"subscription":        {},
	"tag":                 {},
	"targetgroup":         {},
	"topic":               {},
	"trail":               {},
	"user":                {},
	"volume":              {},
	"vpc":                 {},
	"zone":                {},
}

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
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

// CreateEkscluster creates a control plane. The VPC config is a nested struct, so
// the request is assembled here rather than through field tags.
//
// Cluster creation is asynchronous and takes roughly ten minutes; this returns as
// soon as the API accepts the request. Use `awless show ekscluster <name>` to
// watch it reach ACTIVE.
type CreateEkscluster struct {
	_                 string `action:"create" entity:"ekscluster" awsAPI:"eks"`
	logger            *logger.Logger
	graph             cloud.GraphAPI
	api               *eks.Client
	Name              *string   `templateName:"name"`
	Role              *string   `templateName:"role"`
	Subnets           []*string `templateName:"subnets"`
	Securitygroups    []*string `templateName:"securitygroups"`
	Version           *string   `templateName:"version"`
	PublicAccess      *bool     `templateName:"public-access"`
	PrivateAccess     *bool     `templateName:"private-access"`
	PublicAccessCidrs []*string `templateName:"public-access-cidrs"`
}

func (cmd *CreateEkscluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("role"), params.Key("subnets"),
		params.Opt(
			params.Suggested("version"),
			"securitygroups", "public-access", "private-access", "public-access-cidrs",
		),
	))
}

func (cmd *CreateEkscluster) ManualRun(renv env.Running) (any, error) {
	subnets := awssdk.ToStringSlice(cmd.Subnets)
	if len(subnets) == 0 {
		return nil, fmt.Errorf("'subnets' must contain at least one element")
	}
	// EKS requires subnets in at least two availability zones. Caught here because
	// the API error names neither the constraint nor the subnets.
	if len(subnets) < 2 {
		return nil, fmt.Errorf("EKS requires at least two subnets in different availability zones, got %d", len(subnets))
	}

	vpcConfig := &ekstypes.VpcConfigRequest{SubnetIds: subnets}
	if sgs := awssdk.ToStringSlice(cmd.Securitygroups); len(sgs) > 0 {
		vpcConfig.SecurityGroupIds = sgs
	}
	if cmd.PublicAccess != nil {
		vpcConfig.EndpointPublicAccess = cmd.PublicAccess
	}
	if cmd.PrivateAccess != nil {
		vpcConfig.EndpointPrivateAccess = cmd.PrivateAccess
	}
	if cidrs := awssdk.ToStringSlice(cmd.PublicAccessCidrs); len(cidrs) > 0 {
		vpcConfig.PublicAccessCidrs = cidrs
	}

	in := &eks.CreateClusterInput{
		Name:               cmd.Name,
		RoleArn:            cmd.Role,
		ResourcesVpcConfig: vpcConfig,
		Version:            cmd.Version,
	}

	start := time.Now()
	output, err := cmd.api.CreateCluster(renv.RequestContext(), in)
	cmd.logger.ExtraVerbosef("eks.CreateCluster call took %s", time.Since(start))
	if err != nil {
		return nil, fmt.Errorf("create ekscluster: %w", err)
	}
	cmd.logger.Infof("cluster %s is being created; this usually takes about 10 minutes",
		awssdk.ToString(cmd.Name))
	return output, nil
}

func (cmd *CreateEkscluster) ExtractResult(i any) string {
	return awssdk.ToString(i.(*eks.CreateClusterOutput).Cluster.Name)
}

type DeleteEkscluster struct {
	_      string `action:"delete" entity:"ekscluster" awsAPI:"eks" awsCall:"DeleteCluster" awsInput:"eks.DeleteClusterInput" awsOutput:"eks.DeleteClusterOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *eks.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteEkscluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// CreateEksnodegroup adds a managed node group to an existing cluster.
type CreateEksnodegroup struct {
	_            string `action:"create" entity:"eksnodegroup" awsAPI:"eks"`
	logger       *logger.Logger
	graph        cloud.GraphAPI
	api          *eks.Client
	Name         *string   `templateName:"name"`
	Cluster      *string   `templateName:"cluster"`
	Role         *string   `templateName:"role"`
	Subnets      []*string `templateName:"subnets"`
	InstanceType *string   `templateName:"instance-type"`
	MinSize      *int64    `templateName:"min-size"`
	MaxSize      *int64    `templateName:"max-size"`
	DesiredSize  *int64    `templateName:"desired-size"`
	DiskSize     *int64    `templateName:"disk-size"`
	AmiType      *string   `templateName:"ami-type"`
}

func (cmd *CreateEksnodegroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("cluster"), params.Key("role"), params.Key("subnets"),
		params.Opt(
			params.Suggested("instance-type"),
			"min-size", "max-size", "desired-size", "disk-size", "ami-type",
		),
	))
}

func (cmd *CreateEksnodegroup) ManualRun(renv env.Running) (any, error) {
	subnets := awssdk.ToStringSlice(cmd.Subnets)
	if len(subnets) == 0 {
		return nil, fmt.Errorf("'subnets' must contain at least one element")
	}

	// The API requires all three sizes together, so derive the ones not given
	// rather than rejecting a partial scaling config.
	desired := int64(2)
	if cmd.DesiredSize != nil {
		desired = *cmd.DesiredSize
	}
	minSize, maxSize := desired, desired
	if cmd.MinSize != nil {
		minSize = *cmd.MinSize
	}
	if cmd.MaxSize != nil {
		maxSize = *cmd.MaxSize
	}
	if minSize > desired || desired > maxSize {
		return nil, fmt.Errorf("desired-size %d must be between min-size %d and max-size %d", desired, minSize, maxSize)
	}

	in := &eks.CreateNodegroupInput{
		NodegroupName: cmd.Name,
		ClusterName:   cmd.Cluster,
		NodeRole:      cmd.Role,
		Subnets:       subnets,
		ScalingConfig: &ekstypes.NodegroupScalingConfig{
			MinSize:     awssdk.Int32(int32(minSize)),
			MaxSize:     awssdk.Int32(int32(maxSize)),
			DesiredSize: awssdk.Int32(int32(desired)),
		},
	}
	if t := StringValue(cmd.InstanceType); t != "" {
		in.InstanceTypes = []string{t}
	}
	if cmd.DiskSize != nil {
		in.DiskSize = awssdk.Int32(int32(*cmd.DiskSize))
	}
	if a := StringValue(cmd.AmiType); a != "" {
		in.AmiType = ekstypes.AMITypes(a)
	}

	start := time.Now()
	output, err := cmd.api.CreateNodegroup(renv.RequestContext(), in)
	cmd.logger.ExtraVerbosef("eks.CreateNodegroup call took %s", time.Since(start))
	if err != nil {
		return nil, fmt.Errorf("create eksnodegroup: %w", err)
	}
	return output, nil
}

func (cmd *CreateEksnodegroup) ExtractResult(i any) string {
	return awssdk.ToString(i.(*eks.CreateNodegroupOutput).Nodegroup.NodegroupName)
}

type DeleteEksnodegroup struct {
	_       string `action:"delete" entity:"eksnodegroup" awsAPI:"eks" awsCall:"DeleteNodegroup" awsInput:"eks.DeleteNodegroupInput" awsOutput:"eks.DeleteNodegroupOutput"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *eks.Client
	Name    *string `awsName:"NodegroupName" awsType:"awsstr" templateName:"name"`
	Cluster *string `awsName:"ClusterName" awsType:"awsstr" templateName:"cluster"`
}

func (cmd *DeleteEksnodegroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("cluster")))
}

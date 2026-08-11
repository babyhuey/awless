// Auto generated implementation for the AWS cloud service

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

package awsfetch

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"context"

	acm "github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	autoscaling "github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudtrail "github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	cloudwatch "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	efs "github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	elasticache "github.com/aws/aws-sdk-go-v2/service/elasticache"
	elasticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	iam "github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	lambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	rds "github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	secretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	awsconv "github.com/bootswithdefer/awless/aws/conv"
	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
)

func BuildInfraFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualInfraFetchFuncs(conf, funcs)

	funcs["instance"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Instance

		if !conf.getBoolDefaultTrue("aws.infra.instance.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[instance]")
			return resources, objects, nil
		}
		paginator := ec2.NewDescribeInstancesPaginator(conf.APIs.Ec2, &ec2.DescribeInstancesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, all := range out.Reservations {
				for _, output := range all.Instances {
					objects = append(objects, output)
					var res *graph.Resource
					res, err = awsconv.NewResource(output)
					if err != nil {
						return resources, objects, err
					}
					resources = append(resources, res)
				}
			}
		}

		return resources, objects, nil
	}

	funcs["subnet"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Subnet

		if !conf.getBoolDefaultTrue("aws.infra.subnet.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[subnet]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Subnets {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["vpc"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Vpc

		if !conf.getBoolDefaultTrue("aws.infra.vpc.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[vpc]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Vpcs {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["keypair"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.KeyPairInfo

		if !conf.getBoolDefaultTrue("aws.infra.keypair.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[keypair]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.KeyPairs {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["securitygroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.SecurityGroup

		if !conf.getBoolDefaultTrue("aws.infra.securitygroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[securitygroup]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.SecurityGroups {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["volume"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Volume

		if !conf.getBoolDefaultTrue("aws.infra.volume.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[volume]")
			return resources, objects, nil
		}
		paginator := ec2.NewDescribeVolumesPaginator(conf.APIs.Ec2, &ec2.DescribeVolumesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Volumes {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["internetgateway"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.InternetGateway

		if !conf.getBoolDefaultTrue("aws.infra.internetgateway.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[internetgateway]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeInternetGateways(ctx, &ec2.DescribeInternetGatewaysInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.InternetGateways {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["natgateway"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.NatGateway

		if !conf.getBoolDefaultTrue("aws.infra.natgateway.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[natgateway]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.NatGateways {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["routetable"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.RouteTable

		if !conf.getBoolDefaultTrue("aws.infra.routetable.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[routetable]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.RouteTables {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["availabilityzone"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.AvailabilityZone

		if !conf.getBoolDefaultTrue("aws.infra.availabilityzone.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[availabilityzone]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.AvailabilityZones {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["image"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Image

		if !conf.getBoolDefaultTrue("aws.infra.image.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[image]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeImages(ctx, &ec2.DescribeImagesInput{Owners: []string{"self"}})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Images {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["importimagetask"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.ImportImageTask

		if !conf.getBoolDefaultTrue("aws.infra.importimagetask.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[importimagetask]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeImportImageTasks(ctx, &ec2.DescribeImportImageTasksInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.ImportImageTasks {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["elasticip"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Address

		if !conf.getBoolDefaultTrue("aws.infra.elasticip.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[elasticip]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Addresses {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["snapshot"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.Snapshot

		if !conf.getBoolDefaultTrue("aws.infra.snapshot.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[snapshot]")
			return resources, objects, nil
		}
		paginator := ec2.NewDescribeSnapshotsPaginator(conf.APIs.Ec2, &ec2.DescribeSnapshotsInput{OwnerIds: []string{"self"}})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Snapshots {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["networkinterface"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ec2types.NetworkInterface

		if !conf.getBoolDefaultTrue("aws.infra.networkinterface.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[networkinterface]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.NetworkInterfaces {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["classicloadbalancer"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []elbtypes.LoadBalancerDescription

		if !conf.getBoolDefaultTrue("aws.infra.classicloadbalancer.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[classicloadbalancer]")
			return resources, objects, nil
		}
		paginator := elb.NewDescribeLoadBalancersPaginator(conf.APIs.Elb, &elb.DescribeLoadBalancersInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.LoadBalancerDescriptions {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["loadbalancer"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []elbv2types.LoadBalancer

		if !conf.getBoolDefaultTrue("aws.infra.loadbalancer.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[loadbalancer]")
			return resources, objects, nil
		}
		paginator := elbv2.NewDescribeLoadBalancersPaginator(conf.APIs.Elbv2, &elbv2.DescribeLoadBalancersInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.LoadBalancers {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["targetgroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []elbv2types.TargetGroup

		if !conf.getBoolDefaultTrue("aws.infra.targetgroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[targetgroup]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Elbv2.DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.TargetGroups {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["database"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []rdstypes.DBInstance

		if !conf.getBoolDefaultTrue("aws.infra.database.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[database]")
			return resources, objects, nil
		}
		paginator := rds.NewDescribeDBInstancesPaginator(conf.APIs.RDS, &rds.DescribeDBInstancesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.DBInstances {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["dbsubnetgroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []rdstypes.DBSubnetGroup

		if !conf.getBoolDefaultTrue("aws.infra.dbsubnetgroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[dbsubnetgroup]")
			return resources, objects, nil
		}
		paginator := rds.NewDescribeDBSubnetGroupsPaginator(conf.APIs.RDS, &rds.DescribeDBSubnetGroupsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.DBSubnetGroups {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["launchconfiguration"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []autoscalingtypes.LaunchConfiguration

		if !conf.getBoolDefaultTrue("aws.infra.launchconfiguration.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[launchconfiguration]")
			return resources, objects, nil
		}
		paginator := autoscaling.NewDescribeLaunchConfigurationsPaginator(conf.APIs.Autoscaling, &autoscaling.DescribeLaunchConfigurationsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.LaunchConfigurations {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["scalinggroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []autoscalingtypes.AutoScalingGroup

		if !conf.getBoolDefaultTrue("aws.infra.scalinggroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[scalinggroup]")
			return resources, objects, nil
		}
		paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(conf.APIs.Autoscaling, &autoscaling.DescribeAutoScalingGroupsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.AutoScalingGroups {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["scalingpolicy"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []autoscalingtypes.ScalingPolicy

		if !conf.getBoolDefaultTrue("aws.infra.scalingpolicy.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[scalingpolicy]")
			return resources, objects, nil
		}
		paginator := autoscaling.NewDescribePoliciesPaginator(conf.APIs.Autoscaling, &autoscaling.DescribePoliciesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.ScalingPolicies {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["repository"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ecrtypes.Repository

		if !conf.getBoolDefaultTrue("aws.infra.repository.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[repository]")
			return resources, objects, nil
		}
		paginator := ecr.NewDescribeRepositoriesPaginator(conf.APIs.ECR, &ecr.DescribeRepositoriesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Repositories {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["certificate"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []acmtypes.CertificateSummary

		if !conf.getBoolDefaultTrue("aws.infra.certificate.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[certificate]")
			return resources, objects, nil
		}
		paginator := acm.NewListCertificatesPaginator(conf.APIs.ACM, &acm.ListCertificatesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.CertificateSummaryList {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildAccessFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualAccessFetchFuncs(conf, funcs)

	funcs["instanceprofile"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.InstanceProfile

		if !conf.getBoolDefaultTrue("aws.access.instanceprofile.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[instanceprofile]")
			return resources, objects, nil
		}
		paginator := iam.NewListInstanceProfilesPaginator(conf.APIs.IAM, &iam.ListInstanceProfilesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.InstanceProfiles {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["mfadevice"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.VirtualMFADevice

		if !conf.getBoolDefaultTrue("aws.access.mfadevice.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[mfadevice]")
			return resources, objects, nil
		}
		paginator := iam.NewListVirtualMFADevicesPaginator(conf.APIs.IAM, &iam.ListVirtualMFADevicesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.VirtualMFADevices {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildStorageFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualStorageFetchFuncs(conf, funcs)
	return funcs
}
func BuildMessagingFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualMessagingFetchFuncs(conf, funcs)

	funcs["subscription"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []snstypes.Subscription

		if !conf.getBoolDefaultTrue("aws.messaging.subscription.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[subscription]")
			return resources, objects, nil
		}
		paginator := sns.NewListSubscriptionsPaginator(conf.APIs.SNS, &sns.ListSubscriptionsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Subscriptions {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["topic"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []snstypes.Topic

		if !conf.getBoolDefaultTrue("aws.messaging.topic.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[topic]")
			return resources, objects, nil
		}
		paginator := sns.NewListTopicsPaginator(conf.APIs.SNS, &sns.ListTopicsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Topics {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildDNSFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualDNSFetchFuncs(conf, funcs)

	funcs["zone"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []route53types.HostedZone

		if !conf.getBoolDefaultTrue("aws.dns.zone.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource dns[zone]")
			return resources, objects, nil
		}
		paginator := route53.NewListHostedZonesPaginator(conf.APIs.Route53, &route53.ListHostedZonesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.HostedZones {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildLambdaFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualLambdaFetchFuncs(conf, funcs)

	funcs["function"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []lambdatypes.FunctionConfiguration

		if !conf.getBoolDefaultTrue("aws.lambda.function.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource lambda[function]")
			return resources, objects, nil
		}
		paginator := lambda.NewListFunctionsPaginator(conf.APIs.Lambda, &lambda.ListFunctionsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Functions {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildMonitoringFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualMonitoringFetchFuncs(conf, funcs)

	funcs["metric"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []cloudwatchtypes.Metric

		if !conf.getBoolDefaultTrue("aws.monitoring.metric.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource monitoring[metric]")
			return resources, objects, nil
		}
		paginator := cloudwatch.NewListMetricsPaginator(conf.APIs.Cloudwatch, &cloudwatch.ListMetricsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Metrics {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["alarm"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []cloudwatchtypes.MetricAlarm

		if !conf.getBoolDefaultTrue("aws.monitoring.alarm.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource monitoring[alarm]")
			return resources, objects, nil
		}
		paginator := cloudwatch.NewDescribeAlarmsPaginator(conf.APIs.Cloudwatch, &cloudwatch.DescribeAlarmsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.MetricAlarms {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildCDNFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualCDNFetchFuncs(conf, funcs)

	funcs["distribution"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []cloudfronttypes.DistributionSummary

		if !conf.getBoolDefaultTrue("aws.cdn.distribution.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource cdn[distribution]")
			return resources, objects, nil
		}
		paginator := cloudfront.NewListDistributionsPaginator(conf.APIs.Cloudfront, &cloudfront.ListDistributionsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.DistributionList.Items {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildCloudformationFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualCloudformationFetchFuncs(conf, funcs)

	funcs["stack"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []cloudformationtypes.Stack

		if !conf.getBoolDefaultTrue("aws.cloudformation.stack.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource cloudformation[stack]")
			return resources, objects, nil
		}
		paginator := cloudformation.NewDescribeStacksPaginator(conf.APIs.Cloudformation, &cloudformation.DescribeStacksInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Stacks {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildEKSFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualEKSFetchFuncs(conf, funcs)
	return funcs
}
func BuildDynamodbFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualDynamodbFetchFuncs(conf, funcs)
	return funcs
}
func BuildSecretsmanagerFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualSecretsmanagerFetchFuncs(conf, funcs)

	funcs["secret"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []secretsmanagertypes.SecretListEntry

		if !conf.getBoolDefaultTrue("aws.secretsmanager.secret.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource secretsmanager[secret]")
			return resources, objects, nil
		}
		paginator := secretsmanager.NewListSecretsPaginator(conf.APIs.Secretsmanager, &secretsmanager.ListSecretsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.SecretList {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildApigatewayFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualApigatewayFetchFuncs(conf, funcs)
	return funcs
}
func BuildSSMFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualSSMFetchFuncs(conf, funcs)

	funcs["ssmparameter"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ssmtypes.ParameterMetadata

		if !conf.getBoolDefaultTrue("aws.ssm.ssmparameter.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource ssm[ssmparameter]")
			return resources, objects, nil
		}
		paginator := ssm.NewDescribeParametersPaginator(conf.APIs.SSM, &ssm.DescribeParametersInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.Parameters {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildEFSFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualEFSFetchFuncs(conf, funcs)

	funcs["filesystem"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []efstypes.FileSystemDescription

		if !conf.getBoolDefaultTrue("aws.efs.filesystem.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource efs[filesystem]")
			return resources, objects, nil
		}
		paginator := efs.NewDescribeFileSystemsPaginator(conf.APIs.EFS, &efs.DescribeFileSystemsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.FileSystems {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildCloudtrailFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualCloudtrailFetchFuncs(conf, funcs)

	funcs["trail"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []cloudtrailtypes.Trail

		if !conf.getBoolDefaultTrue("aws.cloudtrail.trail.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource cloudtrail[trail]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Cloudtrail.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.TrailList {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildCloudwatchlogsFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualCloudwatchlogsFetchFuncs(conf, funcs)

	funcs["loggroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []cloudwatchlogstypes.LogGroup

		if !conf.getBoolDefaultTrue("aws.cloudwatchlogs.loggroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource cloudwatchlogs[loggroup]")
			return resources, objects, nil
		}
		paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(conf.APIs.Cloudwatchlogs, &cloudwatchlogs.DescribeLogGroupsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.LogGroups {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildElasticacheFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualElasticacheFetchFuncs(conf, funcs)

	funcs["cachecluster"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []elasticachetypes.CacheCluster

		if !conf.getBoolDefaultTrue("aws.elasticache.cachecluster.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource elasticache[cachecluster]")
			return resources, objects, nil
		}
		paginator := elasticache.NewDescribeCacheClustersPaginator(conf.APIs.Elasticache, &elasticache.DescribeCacheClustersInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.CacheClusters {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["replicationgroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []elasticachetypes.ReplicationGroup

		if !conf.getBoolDefaultTrue("aws.elasticache.replicationgroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource elasticache[replicationgroup]")
			return resources, objects, nil
		}
		paginator := elasticache.NewDescribeReplicationGroupsPaginator(conf.APIs.Elasticache, &elasticache.DescribeReplicationGroupsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.ReplicationGroups {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["cachesubnetgroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []elasticachetypes.CacheSubnetGroup

		if !conf.getBoolDefaultTrue("aws.elasticache.cachesubnetgroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource elasticache[cachesubnetgroup]")
			return resources, objects, nil
		}
		paginator := elasticache.NewDescribeCacheSubnetGroupsPaginator(conf.APIs.Elasticache, &elasticache.DescribeCacheSubnetGroupsInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.CacheSubnetGroups {
				objects = append(objects, output)
				var res *graph.Resource
				res, err = awsconv.NewResource(output)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
	return funcs
}
func BuildEventbridgeFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualEventbridgeFetchFuncs(conf, funcs)
	return funcs
}

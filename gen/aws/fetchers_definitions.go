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

package aws

import (
	"strings"

	"github.com/bootswithdefer/awless/cloud"
)

func ApiToInterface(api string) string {
	switch api {
	case "autoscaling":
		return "AutoScalingAPI"
	case "cloudwatch":
		return "CloudWatchAPI"
	case "cloudfront":
		return "CloudFrontAPI"
	case "applicationautoscaling":
		return "ApplicationAutoScalingAPI"
	case "cloudformation":
		return "CloudFormationAPI"
	case "apigatewayv2":
		return "ApiGatewayV2API"
	case "secretsmanager":
		return "SecretsManagerAPI"
	case "cloudtrail":
		return "CloudTrailAPI"
	case "cloudwatchlogs":
		return "CloudWatchLogsAPI"
	case "route53", "lambda":
		return strings.Title(api) + "API"
	default:
		return strings.ToUpper(api) + "API"
	}
}

// SdkModulePath maps the short API name to the AWS SDK v2 Go module path.
// Most services use the same name, but some differ (e.g., elb → elasticloadbalancing).
func SdkModulePath(api string) string {
	switch api {
	case "elb":
		return "elasticloadbalancing"
	case "elbv2":
		return "elasticloadbalancingv2"
	default:
		return api
	}
}

type fetchersDef struct {
	Name     string
	Global   bool
	Api      []string
	Fetchers []fetcher
}

// AutoFetcherAPIs returns the set of API names used by non-manual fetchers.
func (d fetchersDef) AutoFetcherAPIs() []string {
	seen := make(map[string]bool)
	var apis []string
	for _, f := range d.Fetchers {
		if !f.ManualFetcher && !seen[f.Api] {
			seen[f.Api] = true
			apis = append(apis, f.Api)
		}
	}
	return apis
}

type fetcher struct {
	ResourceType                                string
	AWSType                                     string
	ApiMethod, Input                            string
	Output, OutputsContainers, OutputsExtractor string
	ManualFetcher                               bool
	Multipage                                   bool
	NextPageMarker                              string
	Api                                         string
}

var FetchersDefs = []fetchersDef{
	{
		Name: "infra",
		Api:  []string{"ec2", "elbv2", "elb", "rds", "autoscaling", "ecr", "ecs", "applicationautoscaling", "acm"},
		Fetchers: []fetcher{
			{Api: "ec2", ResourceType: cloud.Instance, AWSType: "ec2types.Instance", ApiMethod: "DescribeInstances", Input: "ec2.DescribeInstancesInput{}", Output: "ec2.DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: cloud.Subnet, AWSType: "ec2types.Subnet", ApiMethod: "DescribeSubnets", Input: "ec2.DescribeSubnetsInput{}", Output: "ec2.DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{Api: "ec2", ResourceType: cloud.Vpc, AWSType: "ec2types.Vpc", ApiMethod: "DescribeVpcs", Input: "ec2.DescribeVpcsInput{}", Output: "ec2.DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{Api: "ec2", ResourceType: cloud.Keypair, AWSType: "ec2types.KeyPairInfo", ApiMethod: "DescribeKeyPairs", Input: "ec2.DescribeKeyPairsInput{}", Output: "ec2.DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{Api: "ec2", ResourceType: cloud.SecurityGroup, AWSType: "ec2types.SecurityGroup", ApiMethod: "DescribeSecurityGroups", Input: "ec2.DescribeSecurityGroupsInput{}", Output: "ec2.DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{Api: "ec2", ResourceType: cloud.Volume, AWSType: "ec2types.Volume", ApiMethod: "DescribeVolumes", Input: "ec2.DescribeVolumesInput{}", Output: "ec2.DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: cloud.InternetGateway, AWSType: "ec2types.InternetGateway", ApiMethod: "DescribeInternetGateways", Input: "ec2.DescribeInternetGatewaysInput{}", Output: "ec2.DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{Api: "ec2", ResourceType: cloud.NatGateway, AWSType: "ec2types.NatGateway", ApiMethod: "DescribeNatGateways", Input: "ec2.DescribeNatGatewaysInput{}", Output: "ec2.DescribeNatGatewaysOutput", OutputsExtractor: "NatGateways"},
			{Api: "ec2", ResourceType: cloud.RouteTable, AWSType: "ec2types.RouteTable", ApiMethod: "DescribeRouteTables", Input: "ec2.DescribeRouteTablesInput{}", Output: "ec2.DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{Api: "ec2", ResourceType: cloud.AvailabilityZone, AWSType: "ec2types.AvailabilityZone", ApiMethod: "DescribeAvailabilityZones", Input: "ec2.DescribeAvailabilityZonesInput{}", Output: "ec2.DescribeAvailabilityZonesOutput", OutputsExtractor: "AvailabilityZones"},
			{Api: "ec2", ResourceType: cloud.Image, AWSType: "ec2types.Image", ApiMethod: "DescribeImages", Input: "ec2.DescribeImagesInput{Owners: []string{\"self\"}}", Output: "ec2.DescribeImagesOutput", OutputsExtractor: "Images"},
			{Api: "ec2", ResourceType: cloud.ImportImageTask, AWSType: "ec2types.ImportImageTask", ApiMethod: "DescribeImportImageTasks", Input: "ec2.DescribeImportImageTasksInput{}", Output: "ec2.DescribeImportImageTasksOutput", OutputsExtractor: "ImportImageTasks"},
			{Api: "ec2", ResourceType: cloud.ElasticIP, AWSType: "ec2types.Address", ApiMethod: "DescribeAddresses", Input: "ec2.DescribeAddressesInput{}", Output: "ec2.DescribeAddressesOutput", OutputsExtractor: "Addresses"},
			{Api: "ec2", ResourceType: cloud.Snapshot, AWSType: "ec2types.Snapshot", ApiMethod: "DescribeSnapshots", Input: "ec2.DescribeSnapshotsInput{OwnerIds: []string{\"self\"}}", Output: "ec2.DescribeSnapshotsOutput", OutputsExtractor: "Snapshots", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ec2", ResourceType: cloud.NetworkInterface, AWSType: "ec2types.NetworkInterface", ApiMethod: "DescribeNetworkInterfaces", Input: "ec2.DescribeNetworkInterfacesInput{}", Output: "ec2.DescribeNetworkInterfacesOutput", OutputsExtractor: "NetworkInterfaces"},
			{Api: "elb", ResourceType: cloud.ClassicLoadBalancer, AWSType: "elbtypes.LoadBalancerDescription", ApiMethod: "DescribeLoadBalancers", Input: "elb.DescribeLoadBalancersInput{}", Output: "elb.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancerDescriptions", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "elbv2", ResourceType: cloud.LoadBalancer, AWSType: "elbv2types.LoadBalancer", ApiMethod: "DescribeLoadBalancers", Input: "elbv2.DescribeLoadBalancersInput{}", Output: "elbv2.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancers", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "elbv2", ResourceType: cloud.TargetGroup, AWSType: "elbv2types.TargetGroup", ApiMethod: "DescribeTargetGroups", Input: "elbv2.DescribeTargetGroupsInput{}", Output: "elbv2.DescribeTargetGroupsOutput", OutputsExtractor: "TargetGroups"},
			{Api: "elbv2", ResourceType: cloud.Listener, AWSType: "elbv2types.Listener", ManualFetcher: true},
			{Api: "rds", ResourceType: cloud.Database, AWSType: "rdstypes.DBInstance", ApiMethod: "DescribeDBInstances", Input: "rds.DescribeDBInstancesInput{}", Output: "rds.DescribeDBInstancesOutput", OutputsExtractor: "DBInstances", Multipage: true, NextPageMarker: "Marker"},
			{Api: "rds", ResourceType: cloud.DbSubnetGroup, AWSType: "rdstypes.DBSubnetGroup", ApiMethod: "DescribeDBSubnetGroups", Input: "rds.DescribeDBSubnetGroupsInput{}", Output: "rds.DescribeDBSubnetGroupsOutput", OutputsExtractor: "DBSubnetGroups", Multipage: true, NextPageMarker: "Marker"},
			{Api: "autoscaling", ResourceType: cloud.LaunchConfiguration, AWSType: "autoscalingtypes.LaunchConfiguration", ApiMethod: "DescribeLaunchConfigurations", Input: "autoscaling.DescribeLaunchConfigurationsInput{}", Output: "autoscaling.DescribeLaunchConfigurationsOutput", OutputsExtractor: "LaunchConfigurations", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "autoscaling", ResourceType: cloud.ScalingGroup, AWSType: "autoscalingtypes.AutoScalingGroup", ApiMethod: "DescribeAutoScalingGroups", Input: "autoscaling.DescribeAutoScalingGroupsInput{}", Output: "autoscaling.DescribeAutoScalingGroupsOutput", OutputsExtractor: "AutoScalingGroups", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "autoscaling", ResourceType: cloud.ScalingPolicy, AWSType: "autoscalingtypes.ScalingPolicy", ApiMethod: "DescribePolicies", Input: "autoscaling.DescribePoliciesInput{}", Output: "autoscaling.DescribePoliciesOutput", OutputsExtractor: "ScalingPolicies", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ecr", ResourceType: cloud.Repository, AWSType: "ecrtypes.Repository", ApiMethod: "DescribeRepositories", Input: "ecr.DescribeRepositoriesInput{}", Output: "ecr.DescribeRepositoriesOutput", OutputsExtractor: "Repositories", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "ecs", ResourceType: cloud.ContainerCluster, AWSType: "ecstypes.Cluster", ManualFetcher: true},
			{Api: "ecs", ResourceType: cloud.ContainerTask, AWSType: "ecstypes.TaskDefinition", ManualFetcher: true},
			{Api: "ecs", ResourceType: cloud.Container, AWSType: "ecstypes.Container", ManualFetcher: true},
			{Api: "ecs", ResourceType: cloud.ContainerInstance, AWSType: "ecstypes.ContainerInstance", ManualFetcher: true},
			{Api: "acm", ResourceType: cloud.Certificate, AWSType: "acmtypes.CertificateSummary", ApiMethod: "ListCertificates", Input: "acm.ListCertificatesInput{}", Output: "acm.ListCertificatesOutput", OutputsExtractor: "CertificateSummaryList", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name:   "access",
		Global: true,
		Api:    []string{"iam", "sts"},
		Fetchers: []fetcher{
			{Api: "iam", ResourceType: cloud.User, AWSType: "iamtypes.UserDetail", ManualFetcher: true},
			{Api: "iam", ResourceType: cloud.Group, AWSType: "iamtypes.GroupDetail", ManualFetcher: true},
			{Api: "iam", ResourceType: cloud.Role, AWSType: "iamtypes.RoleDetail", ManualFetcher: true},
			{Api: "iam", ResourceType: cloud.Policy, AWSType: "iamtypes.Policy", ManualFetcher: true},
			{Api: "iam", ResourceType: cloud.AccessKey, AWSType: "iamtypes.AccessKeyMetadata", ManualFetcher: true},
			{Api: "iam", ResourceType: cloud.InstanceProfile, AWSType: "iamtypes.InstanceProfile", ApiMethod: "ListInstanceProfiles", Input: "iam.ListInstanceProfilesInput{}", Output: "iam.ListInstanceProfilesOutput", OutputsExtractor: "InstanceProfiles", Multipage: true, NextPageMarker: "Marker"},
			{Api: "iam", ResourceType: cloud.MFADevice, AWSType: "iamtypes.VirtualMFADevice", ApiMethod: "ListVirtualMFADevices", Input: "iam.ListVirtualMFADevicesInput{}", Output: "iam.ListVirtualMFADevicesOutput", OutputsExtractor: "VirtualMFADevices", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Name: "storage",
		Api:  []string{"s3"},
		Fetchers: []fetcher{
			{Api: "s3", ResourceType: cloud.Bucket, AWSType: "s3types.Bucket", ManualFetcher: true},
			{Api: "s3", ResourceType: cloud.S3Object, AWSType: "s3types.Object", ManualFetcher: true},
		},
	},
	{
		Name: "messaging",
		Api:  []string{"sns", "sqs"},
		Fetchers: []fetcher{
			{Api: "sns", ResourceType: cloud.Subscription, AWSType: "snstypes.Subscription", ApiMethod: "ListSubscriptions", Input: "sns.ListSubscriptionsInput{}", Output: "sns.ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "sns", ResourceType: cloud.Topic, AWSType: "snstypes.Topic", ApiMethod: "ListTopics", Input: "sns.ListTopicsInput{}", Output: "sns.ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "sqs", ResourceType: cloud.Queue, AWSType: "string", ManualFetcher: true},
		},
	},
	{
		Name:   "dns",
		Global: true,
		Api:    []string{"route53"},
		Fetchers: []fetcher{
			{Api: "route53", ResourceType: cloud.Zone, AWSType: "route53types.HostedZone", ApiMethod: "ListHostedZones", Input: "route53.ListHostedZonesInput{}", Output: "route53.ListHostedZonesOutput", OutputsExtractor: "HostedZones", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "route53", ResourceType: cloud.Record, AWSType: "route53types.ResourceRecordSet", ManualFetcher: true},
		},
	},

	{
		Name: "lambda",
		Api:  []string{"lambda"},
		Fetchers: []fetcher{
			{Api: "lambda", ResourceType: cloud.Function, AWSType: "lambdatypes.FunctionConfiguration", ApiMethod: "ListFunctions", Input: "lambda.ListFunctionsInput{}", Output: "lambda.ListFunctionsOutput", OutputsExtractor: "Functions", Multipage: true, NextPageMarker: "NextMarker"},
		},
	},
	{
		Name: "monitoring",
		Api:  []string{"cloudwatch"},
		Fetchers: []fetcher{
			{Api: "cloudwatch", ResourceType: cloud.Metric, AWSType: "cloudwatchtypes.Metric", ApiMethod: "ListMetrics", Input: "cloudwatch.ListMetricsInput{}", Output: "cloudwatch.ListMetricsOutput", OutputsExtractor: "Metrics", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "cloudwatch", ResourceType: cloud.Alarm, AWSType: "cloudwatchtypes.MetricAlarm", ApiMethod: "DescribeAlarms", Input: "cloudwatch.DescribeAlarmsInput{}", Output: "cloudwatch.DescribeAlarmsOutput", OutputsExtractor: "MetricAlarms", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name:   "cdn",
		Global: true,
		Api:    []string{"cloudfront"},
		Fetchers: []fetcher{
			{Api: "cloudfront", ResourceType: cloud.Distribution, AWSType: "cloudfronttypes.DistributionSummary", ApiMethod: "ListDistributions", Input: "cloudfront.ListDistributionsInput{}", Output: "cloudfront.ListDistributionsOutput", OutputsExtractor: "DistributionList.Items", Multipage: true, NextPageMarker: "DistributionList.NextMarker"},
		},
	},
	{
		Name: "cloudformation", //deployment ?
		Api:  []string{"cloudformation"},
		Fetchers: []fetcher{
			{Api: "cloudformation", ResourceType: cloud.Stack, AWSType: "cloudformationtypes.Stack", ApiMethod: "DescribeStacks", Input: "cloudformation.DescribeStacksInput{}", Output: "cloudformation.DescribeStacksOutput", OutputsExtractor: "Stacks", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "eks",
		Api:  []string{"eks"},
		Fetchers: []fetcher{
			{Api: "eks", ResourceType: cloud.EKSCluster, AWSType: "ekstypes.Cluster", ManualFetcher: true},
			{Api: "eks", ResourceType: cloud.EKSNodeGroup, AWSType: "ekstypes.Nodegroup", ManualFetcher: true},
		},
	},
	{
		Name: "dynamodb",
		Api:  []string{"dynamodb"},
		Fetchers: []fetcher{
			{Api: "dynamodb", ResourceType: cloud.DynamoDBTable, AWSType: "dynamodbtypes.TableDescription", ManualFetcher: true},
		},
	},
	{
		Name: "secretsmanager",
		Api:  []string{"secretsmanager", "kms"},
		Fetchers: []fetcher{
			{Api: "secretsmanager", ResourceType: cloud.Secret, AWSType: "secretsmanagertypes.SecretListEntry", ApiMethod: "ListSecrets", Input: "secretsmanager.ListSecretsInput{}", Output: "secretsmanager.ListSecretsOutput", OutputsExtractor: "SecretList", Multipage: true, NextPageMarker: "NextToken"},
			{Api: "kms", ResourceType: cloud.Key, AWSType: "kmstypes.KeyMetadata", ManualFetcher: true},
		},
	},
	{
		Name: "apigateway",
		Api:  []string{"apigatewayv2"},
		Fetchers: []fetcher{
			{Api: "apigatewayv2", ResourceType: cloud.ApiGateway, AWSType: "apigatewayv2types.Api", ManualFetcher: true},
			{Api: "apigatewayv2", ResourceType: cloud.ApiGatewayRoute, AWSType: "apigatewayv2types.Route", ManualFetcher: true},
			{Api: "apigatewayv2", ResourceType: cloud.ApiGatewayStage, AWSType: "apigatewayv2types.Stage", ManualFetcher: true},
		},
	},
	{
		Name: "ssm",
		Api:  []string{"ssm"},
		Fetchers: []fetcher{
			{Api: "ssm", ResourceType: cloud.SSMParameter, AWSType: "ssmtypes.ParameterMetadata", ApiMethod: "DescribeParameters", Input: "ssm.DescribeParametersInput{}", Output: "ssm.DescribeParametersOutput", OutputsExtractor: "Parameters", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "efs",
		Api:  []string{"efs"},
		Fetchers: []fetcher{
			{Api: "efs", ResourceType: cloud.FileSystem, AWSType: "efstypes.FileSystemDescription", ApiMethod: "DescribeFileSystems", Input: "efs.DescribeFileSystemsInput{}", Output: "efs.DescribeFileSystemsOutput", OutputsExtractor: "FileSystems", Multipage: true, NextPageMarker: "NextMarker"},
			{Api: "efs", ResourceType: cloud.MountTarget, AWSType: "efstypes.MountTargetDescription", ManualFetcher: true},
		},
	},
	{
		Name: "cloudtrail",
		Api:  []string{"cloudtrail"},
		Fetchers: []fetcher{
			{Api: "cloudtrail", ResourceType: cloud.Trail, AWSType: "cloudtrailtypes.Trail", ApiMethod: "DescribeTrails", Input: "cloudtrail.DescribeTrailsInput{}", Output: "cloudtrail.DescribeTrailsOutput", OutputsExtractor: "TrailList"},
		},
	},
	{
		Name: "cloudwatchlogs",
		Api:  []string{"cloudwatchlogs"},
		Fetchers: []fetcher{
			{Api: "cloudwatchlogs", ResourceType: cloud.LogGroup, AWSType: "cloudwatchlogstypes.LogGroup", ApiMethod: "DescribeLogGroups", Input: "cloudwatchlogs.DescribeLogGroupsInput{}", Output: "cloudwatchlogs.DescribeLogGroupsOutput", OutputsExtractor: "LogGroups", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
}

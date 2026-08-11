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

func APIToInterface(api string) string {
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
		return capitalize(api) + "API"
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
	API      []string
	Fetchers []fetcher
}

// AutoFetcherAPIs returns the set of API names used by non-manual fetchers.
func (d fetchersDef) AutoFetcherAPIs() []string {
	seen := make(map[string]bool)
	var apis []string
	for _, f := range d.Fetchers {
		if !f.ManualFetcher && !seen[f.API] {
			seen[f.API] = true
			apis = append(apis, f.API)
		}
	}
	return apis
}

type fetcher struct {
	ResourceType                                string
	AWSType                                     string
	APIMethod, Input                            string
	Output, OutputsContainers, OutputsExtractor string
	ManualFetcher                               bool
	Multipage                                   bool
	NextPageMarker                              string
	API                                         string
}

var FetchersDefs = []fetchersDef{
	{
		Name: "infra",
		API:  []string{"ec2", "elbv2", "elb", "rds", "autoscaling", "ecr", "ecs", "applicationautoscaling", "acm"},
		Fetchers: []fetcher{
			// Transit gateways and VPC endpoints ride the EC2 client that the rest of
			// infra already uses, so they need no service of their own.
			{API: "ec2", ResourceType: cloud.TransitGateway, AWSType: "ec2types.TransitGateway", APIMethod: "DescribeTransitGateways", Input: "ec2.DescribeTransitGatewaysInput{}", Output: "ec2.DescribeTransitGatewaysOutput", OutputsExtractor: "TransitGateways", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.TransitGatewayAttachment, AWSType: "ec2types.TransitGatewayVpcAttachment", APIMethod: "DescribeTransitGatewayVpcAttachments", Input: "ec2.DescribeTransitGatewayVpcAttachmentsInput{}", Output: "ec2.DescribeTransitGatewayVpcAttachmentsOutput", OutputsExtractor: "TransitGatewayVpcAttachments", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.TransitGatewayRouteTable, AWSType: "ec2types.TransitGatewayRouteTable", APIMethod: "DescribeTransitGatewayRouteTables", Input: "ec2.DescribeTransitGatewayRouteTablesInput{}", Output: "ec2.DescribeTransitGatewayRouteTablesOutput", OutputsExtractor: "TransitGatewayRouteTables", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.VpcEndpoint, AWSType: "ec2types.VpcEndpoint", APIMethod: "DescribeVpcEndpoints", Input: "ec2.DescribeVpcEndpointsInput{}", Output: "ec2.DescribeVpcEndpointsOutput", OutputsExtractor: "VpcEndpoints", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.Instance, AWSType: "ec2types.Instance", APIMethod: "DescribeInstances", Input: "ec2.DescribeInstancesInput{}", Output: "ec2.DescribeInstancesOutput", OutputsExtractor: "Instances", OutputsContainers: "Reservations", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.Subnet, AWSType: "ec2types.Subnet", APIMethod: "DescribeSubnets", Input: "ec2.DescribeSubnetsInput{}", Output: "ec2.DescribeSubnetsOutput", OutputsExtractor: "Subnets"},
			{API: "ec2", ResourceType: cloud.Vpc, AWSType: "ec2types.Vpc", APIMethod: "DescribeVpcs", Input: "ec2.DescribeVpcsInput{}", Output: "ec2.DescribeVpcsOutput", OutputsExtractor: "Vpcs"},
			{API: "ec2", ResourceType: cloud.Keypair, AWSType: "ec2types.KeyPairInfo", APIMethod: "DescribeKeyPairs", Input: "ec2.DescribeKeyPairsInput{}", Output: "ec2.DescribeKeyPairsOutput", OutputsExtractor: "KeyPairs"},
			{API: "ec2", ResourceType: cloud.SecurityGroup, AWSType: "ec2types.SecurityGroup", APIMethod: "DescribeSecurityGroups", Input: "ec2.DescribeSecurityGroupsInput{}", Output: "ec2.DescribeSecurityGroupsOutput", OutputsExtractor: "SecurityGroups"},
			{API: "ec2", ResourceType: cloud.Volume, AWSType: "ec2types.Volume", APIMethod: "DescribeVolumes", Input: "ec2.DescribeVolumesInput{}", Output: "ec2.DescribeVolumesOutput", OutputsExtractor: "Volumes", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.InternetGateway, AWSType: "ec2types.InternetGateway", APIMethod: "DescribeInternetGateways", Input: "ec2.DescribeInternetGatewaysInput{}", Output: "ec2.DescribeInternetGatewaysOutput", OutputsExtractor: "InternetGateways"},
			{API: "ec2", ResourceType: cloud.NatGateway, AWSType: "ec2types.NatGateway", APIMethod: "DescribeNatGateways", Input: "ec2.DescribeNatGatewaysInput{}", Output: "ec2.DescribeNatGatewaysOutput", OutputsExtractor: "NatGateways"},
			{API: "ec2", ResourceType: cloud.RouteTable, AWSType: "ec2types.RouteTable", APIMethod: "DescribeRouteTables", Input: "ec2.DescribeRouteTablesInput{}", Output: "ec2.DescribeRouteTablesOutput", OutputsExtractor: "RouteTables"},
			{API: "ec2", ResourceType: cloud.AvailabilityZone, AWSType: "ec2types.AvailabilityZone", APIMethod: "DescribeAvailabilityZones", Input: "ec2.DescribeAvailabilityZonesInput{}", Output: "ec2.DescribeAvailabilityZonesOutput", OutputsExtractor: "AvailabilityZones"},
			{API: "ec2", ResourceType: cloud.Image, AWSType: "ec2types.Image", APIMethod: "DescribeImages", Input: "ec2.DescribeImagesInput{Owners: []string{\"self\"}}", Output: "ec2.DescribeImagesOutput", OutputsExtractor: "Images"},
			{API: "ec2", ResourceType: cloud.ImportImageTask, AWSType: "ec2types.ImportImageTask", APIMethod: "DescribeImportImageTasks", Input: "ec2.DescribeImportImageTasksInput{}", Output: "ec2.DescribeImportImageTasksOutput", OutputsExtractor: "ImportImageTasks"},
			{API: "ec2", ResourceType: cloud.ElasticIP, AWSType: "ec2types.Address", APIMethod: "DescribeAddresses", Input: "ec2.DescribeAddressesInput{}", Output: "ec2.DescribeAddressesOutput", OutputsExtractor: "Addresses"},
			{API: "ec2", ResourceType: cloud.Snapshot, AWSType: "ec2types.Snapshot", APIMethod: "DescribeSnapshots", Input: "ec2.DescribeSnapshotsInput{OwnerIds: []string{\"self\"}}", Output: "ec2.DescribeSnapshotsOutput", OutputsExtractor: "Snapshots", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ec2", ResourceType: cloud.NetworkInterface, AWSType: "ec2types.NetworkInterface", APIMethod: "DescribeNetworkInterfaces", Input: "ec2.DescribeNetworkInterfacesInput{}", Output: "ec2.DescribeNetworkInterfacesOutput", OutputsExtractor: "NetworkInterfaces"},
			{API: "elb", ResourceType: cloud.ClassicLoadBalancer, AWSType: "elbtypes.LoadBalancerDescription", APIMethod: "DescribeLoadBalancers", Input: "elb.DescribeLoadBalancersInput{}", Output: "elb.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancerDescriptions", Multipage: true, NextPageMarker: "NextMarker"},
			{API: "elbv2", ResourceType: cloud.LoadBalancer, AWSType: "elbv2types.LoadBalancer", APIMethod: "DescribeLoadBalancers", Input: "elbv2.DescribeLoadBalancersInput{}", Output: "elbv2.DescribeLoadBalancersOutput", OutputsExtractor: "LoadBalancers", Multipage: true, NextPageMarker: "NextMarker"},
			{API: "elbv2", ResourceType: cloud.TargetGroup, AWSType: "elbv2types.TargetGroup", APIMethod: "DescribeTargetGroups", Input: "elbv2.DescribeTargetGroupsInput{}", Output: "elbv2.DescribeTargetGroupsOutput", OutputsExtractor: "TargetGroups"},
			{API: "elbv2", ResourceType: cloud.Listener, AWSType: "elbv2types.Listener", ManualFetcher: true},
			{API: "rds", ResourceType: cloud.Database, AWSType: "rdstypes.DBInstance", APIMethod: "DescribeDBInstances", Input: "rds.DescribeDBInstancesInput{}", Output: "rds.DescribeDBInstancesOutput", OutputsExtractor: "DBInstances", Multipage: true, NextPageMarker: "Marker"},
			{API: "rds", ResourceType: cloud.DBSubnetGroup, AWSType: "rdstypes.DBSubnetGroup", APIMethod: "DescribeDBSubnetGroups", Input: "rds.DescribeDBSubnetGroupsInput{}", Output: "rds.DescribeDBSubnetGroupsOutput", OutputsExtractor: "DBSubnetGroups", Multipage: true, NextPageMarker: "Marker"},
			{API: "autoscaling", ResourceType: cloud.LaunchConfiguration, AWSType: "autoscalingtypes.LaunchConfiguration", APIMethod: "DescribeLaunchConfigurations", Input: "autoscaling.DescribeLaunchConfigurationsInput{}", Output: "autoscaling.DescribeLaunchConfigurationsOutput", OutputsExtractor: "LaunchConfigurations", Multipage: true, NextPageMarker: "NextToken"},
			{API: "autoscaling", ResourceType: cloud.ScalingGroup, AWSType: "autoscalingtypes.AutoScalingGroup", APIMethod: "DescribeAutoScalingGroups", Input: "autoscaling.DescribeAutoScalingGroupsInput{}", Output: "autoscaling.DescribeAutoScalingGroupsOutput", OutputsExtractor: "AutoScalingGroups", Multipage: true, NextPageMarker: "NextToken"},
			{API: "autoscaling", ResourceType: cloud.ScalingPolicy, AWSType: "autoscalingtypes.ScalingPolicy", APIMethod: "DescribePolicies", Input: "autoscaling.DescribePoliciesInput{}", Output: "autoscaling.DescribePoliciesOutput", OutputsExtractor: "ScalingPolicies", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ecr", ResourceType: cloud.Repository, AWSType: "ecrtypes.Repository", APIMethod: "DescribeRepositories", Input: "ecr.DescribeRepositoriesInput{}", Output: "ecr.DescribeRepositoriesOutput", OutputsExtractor: "Repositories", Multipage: true, NextPageMarker: "NextToken"},
			{API: "ecs", ResourceType: cloud.ContainerCluster, AWSType: "ecstypes.Cluster", ManualFetcher: true},
			{API: "ecs", ResourceType: cloud.ContainerTask, AWSType: "ecstypes.TaskDefinition", ManualFetcher: true},
			{API: "ecs", ResourceType: cloud.Container, AWSType: "ecstypes.Container", ManualFetcher: true},
			{API: "ecs", ResourceType: cloud.ContainerInstance, AWSType: "ecstypes.ContainerInstance", ManualFetcher: true},
			{API: "acm", ResourceType: cloud.Certificate, AWSType: "acmtypes.CertificateSummary", APIMethod: "ListCertificates", Input: "acm.ListCertificatesInput{}", Output: "acm.ListCertificatesOutput", OutputsExtractor: "CertificateSummaryList", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name:   "access",
		Global: true,
		API:    []string{"iam", "sts"},
		Fetchers: []fetcher{
			{API: "iam", ResourceType: cloud.User, AWSType: "iamtypes.UserDetail", ManualFetcher: true},
			{API: "iam", ResourceType: cloud.Group, AWSType: "iamtypes.GroupDetail", ManualFetcher: true},
			{API: "iam", ResourceType: cloud.Role, AWSType: "iamtypes.RoleDetail", ManualFetcher: true},
			{API: "iam", ResourceType: cloud.Policy, AWSType: "iamtypes.Policy", ManualFetcher: true},
			{API: "iam", ResourceType: cloud.AccessKey, AWSType: "iamtypes.AccessKeyMetadata", ManualFetcher: true},
			{API: "iam", ResourceType: cloud.InstanceProfile, AWSType: "iamtypes.InstanceProfile", APIMethod: "ListInstanceProfiles", Input: "iam.ListInstanceProfilesInput{}", Output: "iam.ListInstanceProfilesOutput", OutputsExtractor: "InstanceProfiles", Multipage: true, NextPageMarker: "Marker"},
			{API: "iam", ResourceType: cloud.MFADevice, AWSType: "iamtypes.VirtualMFADevice", APIMethod: "ListVirtualMFADevices", Input: "iam.ListVirtualMFADevicesInput{}", Output: "iam.ListVirtualMFADevicesOutput", OutputsExtractor: "VirtualMFADevices", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Name: "storage",
		API:  []string{"s3"},
		Fetchers: []fetcher{
			{API: "s3", ResourceType: cloud.Bucket, AWSType: "s3types.Bucket", ManualFetcher: true},
			{API: "s3", ResourceType: cloud.S3Object, AWSType: "s3types.Object", ManualFetcher: true},
		},
	},
	{
		Name: "messaging",
		API:  []string{"sns", "sqs"},
		Fetchers: []fetcher{
			{API: "sns", ResourceType: cloud.Subscription, AWSType: "snstypes.Subscription", APIMethod: "ListSubscriptions", Input: "sns.ListSubscriptionsInput{}", Output: "sns.ListSubscriptionsOutput", OutputsExtractor: "Subscriptions", Multipage: true, NextPageMarker: "NextToken"},
			{API: "sns", ResourceType: cloud.Topic, AWSType: "snstypes.Topic", APIMethod: "ListTopics", Input: "sns.ListTopicsInput{}", Output: "sns.ListTopicsOutput", OutputsExtractor: "Topics", Multipage: true, NextPageMarker: "NextToken"},
			{API: "sqs", ResourceType: cloud.Queue, AWSType: "string", ManualFetcher: true},
		},
	},
	{
		Name:   "dns",
		Global: true,
		API:    []string{"route53"},
		Fetchers: []fetcher{
			{API: "route53", ResourceType: cloud.Zone, AWSType: "route53types.HostedZone", APIMethod: "ListHostedZones", Input: "route53.ListHostedZonesInput{}", Output: "route53.ListHostedZonesOutput", OutputsExtractor: "HostedZones", Multipage: true, NextPageMarker: "NextMarker"},
			{API: "route53", ResourceType: cloud.Record, AWSType: "route53types.ResourceRecordSet", ManualFetcher: true},
		},
	},

	{
		Name: "lambda",
		API:  []string{"lambda"},
		Fetchers: []fetcher{
			{API: "lambda", ResourceType: cloud.Function, AWSType: "lambdatypes.FunctionConfiguration", APIMethod: "ListFunctions", Input: "lambda.ListFunctionsInput{}", Output: "lambda.ListFunctionsOutput", OutputsExtractor: "Functions", Multipage: true, NextPageMarker: "NextMarker"},
		},
	},
	{
		Name: "monitoring",
		API:  []string{"cloudwatch"},
		Fetchers: []fetcher{
			{API: "cloudwatch", ResourceType: cloud.Metric, AWSType: "cloudwatchtypes.Metric", APIMethod: "ListMetrics", Input: "cloudwatch.ListMetricsInput{}", Output: "cloudwatch.ListMetricsOutput", OutputsExtractor: "Metrics", Multipage: true, NextPageMarker: "NextToken"},
			{API: "cloudwatch", ResourceType: cloud.Alarm, AWSType: "cloudwatchtypes.MetricAlarm", APIMethod: "DescribeAlarms", Input: "cloudwatch.DescribeAlarmsInput{}", Output: "cloudwatch.DescribeAlarmsOutput", OutputsExtractor: "MetricAlarms", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name:   "cdn",
		Global: true,
		API:    []string{"cloudfront"},
		Fetchers: []fetcher{
			{API: "cloudfront", ResourceType: cloud.Distribution, AWSType: "cloudfronttypes.DistributionSummary", APIMethod: "ListDistributions", Input: "cloudfront.ListDistributionsInput{}", Output: "cloudfront.ListDistributionsOutput", OutputsExtractor: "DistributionList.Items", Multipage: true, NextPageMarker: "DistributionList.NextMarker"},
		},
	},
	{
		Name: "cloudformation", //deployment ?
		API:  []string{"cloudformation"},
		Fetchers: []fetcher{
			{API: "cloudformation", ResourceType: cloud.Stack, AWSType: "cloudformationtypes.Stack", APIMethod: "DescribeStacks", Input: "cloudformation.DescribeStacksInput{}", Output: "cloudformation.DescribeStacksOutput", OutputsExtractor: "Stacks", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "eks",
		API:  []string{"eks"},
		Fetchers: []fetcher{
			{API: "eks", ResourceType: cloud.EKSCluster, AWSType: "ekstypes.Cluster", ManualFetcher: true},
			{API: "eks", ResourceType: cloud.EKSNodeGroup, AWSType: "ekstypes.Nodegroup", ManualFetcher: true},
		},
	},
	{
		Name: "dynamodb",
		API:  []string{"dynamodb"},
		Fetchers: []fetcher{
			{API: "dynamodb", ResourceType: cloud.DynamoDBTable, AWSType: "dynamodbtypes.TableDescription", ManualFetcher: true},
		},
	},
	{
		Name: "secretsmanager",
		API:  []string{"secretsmanager", "kms"},
		Fetchers: []fetcher{
			{API: "secretsmanager", ResourceType: cloud.Secret, AWSType: "secretsmanagertypes.SecretListEntry", APIMethod: "ListSecrets", Input: "secretsmanager.ListSecretsInput{}", Output: "secretsmanager.ListSecretsOutput", OutputsExtractor: "SecretList", Multipage: true, NextPageMarker: "NextToken"},
			{API: "kms", ResourceType: cloud.Key, AWSType: "kmstypes.KeyMetadata", ManualFetcher: true},
		},
	},
	{
		Name: "apigateway",
		API:  []string{"apigatewayv2"},
		Fetchers: []fetcher{
			{API: "apigatewayv2", ResourceType: cloud.APIGateway, AWSType: "apigatewayv2types.Api", ManualFetcher: true},
			{API: "apigatewayv2", ResourceType: cloud.APIGatewayRoute, AWSType: "apigatewayv2types.Route", ManualFetcher: true},
			{API: "apigatewayv2", ResourceType: cloud.APIGatewayStage, AWSType: "apigatewayv2types.Stage", ManualFetcher: true},
		},
	},
	{
		Name: "ssm",
		API:  []string{"ssm"},
		Fetchers: []fetcher{
			{API: "ssm", ResourceType: cloud.SSMParameter, AWSType: "ssmtypes.ParameterMetadata", APIMethod: "DescribeParameters", Input: "ssm.DescribeParametersInput{}", Output: "ssm.DescribeParametersOutput", OutputsExtractor: "Parameters", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "efs",
		API:  []string{"efs"},
		Fetchers: []fetcher{
			{API: "efs", ResourceType: cloud.FileSystem, AWSType: "efstypes.FileSystemDescription", APIMethod: "DescribeFileSystems", Input: "efs.DescribeFileSystemsInput{}", Output: "efs.DescribeFileSystemsOutput", OutputsExtractor: "FileSystems", Multipage: true, NextPageMarker: "NextMarker"},
			{API: "efs", ResourceType: cloud.MountTarget, AWSType: "efstypes.MountTargetDescription", ManualFetcher: true},
		},
	},
	{
		Name: "cloudtrail",
		API:  []string{"cloudtrail"},
		Fetchers: []fetcher{
			{API: "cloudtrail", ResourceType: cloud.Trail, AWSType: "cloudtrailtypes.Trail", APIMethod: "DescribeTrails", Input: "cloudtrail.DescribeTrailsInput{}", Output: "cloudtrail.DescribeTrailsOutput", OutputsExtractor: "TrailList"},
		},
	},
	{
		Name: "cloudwatchlogs",
		API:  []string{"cloudwatchlogs"},
		Fetchers: []fetcher{
			{API: "cloudwatchlogs", ResourceType: cloud.LogGroup, AWSType: "cloudwatchlogstypes.LogGroup", APIMethod: "DescribeLogGroups", Input: "cloudwatchlogs.DescribeLogGroupsInput{}", Output: "cloudwatchlogs.DescribeLogGroupsOutput", OutputsExtractor: "LogGroups", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "elasticache",
		API:  []string{"elasticache"},
		Fetchers: []fetcher{
			{API: "elasticache", ResourceType: cloud.CacheCluster, AWSType: "elasticachetypes.CacheCluster", APIMethod: "DescribeCacheClusters", Input: "elasticache.DescribeCacheClustersInput{}", Output: "elasticache.DescribeCacheClustersOutput", OutputsExtractor: "CacheClusters", Multipage: true, NextPageMarker: "Marker"},
			{API: "elasticache", ResourceType: cloud.ReplicationGroup, AWSType: "elasticachetypes.ReplicationGroup", APIMethod: "DescribeReplicationGroups", Input: "elasticache.DescribeReplicationGroupsInput{}", Output: "elasticache.DescribeReplicationGroupsOutput", OutputsExtractor: "ReplicationGroups", Multipage: true, NextPageMarker: "Marker"},
			{API: "elasticache", ResourceType: cloud.CacheSubnetGroup, AWSType: "elasticachetypes.CacheSubnetGroup", APIMethod: "DescribeCacheSubnetGroups", Input: "elasticache.DescribeCacheSubnetGroupsInput{}", Output: "elasticache.DescribeCacheSubnetGroupsOutput", OutputsExtractor: "CacheSubnetGroups", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Name: "eventbridge",
		API:  []string{"eventbridge"},
		Fetchers: []fetcher{
			// EventBridge publishes no paginators, and the generator's only multipage
			// mode drives an SDK paginator. Fetched manually rather than as a single
			// page, which would silently stop at the API's default page size.
			{API: "eventbridge", ResourceType: cloud.EventBus, AWSType: "eventbridgetypes.EventBus", ManualFetcher: true},
			{API: "eventbridge", ResourceType: cloud.EventRule, AWSType: "eventbridgetypes.Rule", ManualFetcher: true},
		},
	},
	{
		Name: "stepfunctions",
		API:  []string{"sfn"},
		Fetchers: []fetcher{
			// The list item carries only name, ARN, type and creation date; the
			// definition and role need a DescribeStateMachine per machine and are not
			// worth an N+1 on every sync.
			{API: "sfn", ResourceType: cloud.StateMachine, AWSType: "sfntypes.StateMachineListItem", APIMethod: "ListStateMachines", Input: "sfn.ListStateMachinesInput{}", Output: "sfn.ListStateMachinesOutput", OutputsExtractor: "StateMachines", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "waf",
		API:  []string{"wafv2"},
		Fetchers: []fetcher{
			// WAF v2 publishes no paginators, and every call is scoped: REGIONAL
			// resources are visible from their own region, CLOUDFRONT ones only from
			// us-east-1. Both scopes are walked manually so neither is silently
			// missing from a sync.
			{API: "wafv2", ResourceType: cloud.WebACL, AWSType: "wafv2types.WebACLSummary", ManualFetcher: true},
			{API: "wafv2", ResourceType: cloud.IPSet, AWSType: "wafv2types.IPSetSummary", ManualFetcher: true},
			{API: "wafv2", ResourceType: cloud.RuleGroup, AWSType: "wafv2types.RuleGroupSummary", ManualFetcher: true},
		},
	},
	{
		Name: "configservice",
		API:  []string{"configservice"},
		Fetchers: []fetcher{
			// Manual so compliance can be merged in. Whether a rule is passing is the
			// reason to look at Config at all, and it comes from a second API; both
			// calls paginate over all rules rather than one call per rule.
			{API: "configservice", ResourceType: cloud.ConfigRule, AWSType: "configservicetypes.ConfigRule", ManualFetcher: true},
		},
	},
	{
		Name: "kinesis",
		API:  []string{"kinesis"},
		Fetchers: []fetcher{
			{API: "kinesis", ResourceType: cloud.Stream, AWSType: "kinesistypes.StreamSummary", APIMethod: "ListStreams", Input: "kinesis.ListStreamsInput{}", Output: "kinesis.ListStreamsOutput", OutputsExtractor: "StreamSummaries", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "redshift",
		API:  []string{"redshift"},
		Fetchers: []fetcher{
			{API: "redshift", ResourceType: cloud.RedshiftCluster, AWSType: "redshifttypes.Cluster", APIMethod: "DescribeClusters", Input: "redshift.DescribeClustersInput{}", Output: "redshift.DescribeClustersOutput", OutputsExtractor: "Clusters", Multipage: true, NextPageMarker: "Marker"},
			{API: "redshift", ResourceType: cloud.RedshiftSubnetGroup, AWSType: "redshifttypes.ClusterSubnetGroup", APIMethod: "DescribeClusterSubnetGroups", Input: "redshift.DescribeClusterSubnetGroupsInput{}", Output: "redshift.DescribeClusterSubnetGroupsOutput", OutputsExtractor: "ClusterSubnetGroups", Multipage: true, NextPageMarker: "Marker"},
		},
	},
	{
		Name: "codepipeline",
		API:  []string{"codepipeline"},
		Fetchers: []fetcher{
			{API: "codepipeline", ResourceType: cloud.Pipeline, AWSType: "codepipelinetypes.PipelineSummary", APIMethod: "ListPipelines", Input: "codepipeline.ListPipelinesInput{}", Output: "codepipeline.ListPipelinesOutput", OutputsExtractor: "Pipelines", Multipage: true, NextPageMarker: "NextToken"},
		},
	},
	{
		Name: "codebuild",
		API:  []string{"codebuild"},
		Fetchers: []fetcher{
			// ListProjects returns names only, so the details need a second call.
			{API: "codebuild", ResourceType: cloud.BuildProject, AWSType: "codebuildtypes.Project", ManualFetcher: true},
		},
	},
	{
		Name: "beanstalk",
		API:  []string{"elasticbeanstalk"},
		Fetchers: []fetcher{
			// DescribeApplications returns everything in one response; it has neither a
			// paginator nor a token.
			{API: "elasticbeanstalk", ResourceType: cloud.Application, AWSType: "elasticbeanstalktypes.ApplicationDescription", APIMethod: "DescribeApplications", Input: "elasticbeanstalk.DescribeApplicationsInput{}", Output: "elasticbeanstalk.DescribeApplicationsOutput", OutputsExtractor: "Applications"},
			// DescribeEnvironments takes a NextToken but has no paginator, so it is
			// walked manually rather than truncated to the first page.
			{API: "elasticbeanstalk", ResourceType: cloud.Environment, AWSType: "elasticbeanstalktypes.EnvironmentDescription", ManualFetcher: true},
		},
	},
}

// capitalize upper-cases the first character of s.
//
// Replaces strings.Title, deprecated in Go 1.18 because it applies Unicode word
// boundaries and title-cases every word. Every input here is a single ASCII
// token — an AWS API name, resource type, template action, or policy effect —
// so this is both correct and narrower than the deprecated behavior.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

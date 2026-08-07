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

package awsservices

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"context"

	acm "github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	autoscaling "github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	_ "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudwatch "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/bootswithdefer/awless/cloud"
)

type mockEc2 struct {
	instances         []ec2types.Instance
	subnets           []ec2types.Subnet
	vpcs              []ec2types.Vpc
	keypairinfos      []ec2types.KeyPairInfo
	securitygroups    []ec2types.SecurityGroup
	volumes           []ec2types.Volume
	internetgateways  []ec2types.InternetGateway
	natgateways       []ec2types.NatGateway
	routetables       []ec2types.RouteTable
	availabilityzones []ec2types.AvailabilityZone
	images            []ec2types.Image
	importimagetasks  []ec2types.ImportImageTask
	addresss          []ec2types.Address
	snapshots         []ec2types.Snapshot
	networkinterfaces []ec2types.NetworkInterface
}

func (m *mockEc2) Name() string {
	return ""
}

func (m *mockEc2) Region() string {
	return ""
}

func (m *mockEc2) Profile() string {
	return ""
}

func (m *mockEc2) Provider() string {
	return ""
}

func (m *mockEc2) ProviderAPI() string {
	return ""
}

func (m *mockEc2) ResourceTypes() []string {
	return []string{}
}

func (m *mockEc2) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEc2) IsSyncDisabled() bool {
	return false
}

func (m *mockEc2) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEc2) DescribeSubnets(ctx context.Context, input *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return &ec2.DescribeSubnetsOutput{Subnets: m.subnets}, nil
}

func (m *mockEc2) DescribeVpcs(ctx context.Context, input *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	return &ec2.DescribeVpcsOutput{Vpcs: m.vpcs}, nil
}

func (m *mockEc2) DescribeKeyPairs(ctx context.Context, input *ec2.DescribeKeyPairsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{KeyPairs: m.keypairinfos}, nil
}

func (m *mockEc2) DescribeSecurityGroups(ctx context.Context, input *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: m.securitygroups}, nil
}

func (m *mockEc2) DescribeVolumes(ctx context.Context, input *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	return &ec2.DescribeVolumesOutput{Volumes: m.volumes}, nil
}

func (m *mockEc2) DescribeInternetGateways(ctx context.Context, input *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error) {
	return &ec2.DescribeInternetGatewaysOutput{InternetGateways: m.internetgateways}, nil
}

func (m *mockEc2) DescribeNatGateways(ctx context.Context, input *ec2.DescribeNatGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNatGatewaysOutput, error) {
	return &ec2.DescribeNatGatewaysOutput{NatGateways: m.natgateways}, nil
}

func (m *mockEc2) DescribeRouteTables(ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error) {
	return &ec2.DescribeRouteTablesOutput{RouteTables: m.routetables}, nil
}

func (m *mockEc2) DescribeAvailabilityZones(ctx context.Context, input *ec2.DescribeAvailabilityZonesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return &ec2.DescribeAvailabilityZonesOutput{AvailabilityZones: m.availabilityzones}, nil
}

func (m *mockEc2) DescribeImages(ctx context.Context, input *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error) {
	return &ec2.DescribeImagesOutput{Images: m.images}, nil
}

func (m *mockEc2) DescribeImportImageTasks(ctx context.Context, input *ec2.DescribeImportImageTasksInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImportImageTasksOutput, error) {
	return &ec2.DescribeImportImageTasksOutput{ImportImageTasks: m.importimagetasks}, nil
}

func (m *mockEc2) DescribeAddresses(ctx context.Context, input *ec2.DescribeAddressesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAddressesOutput, error) {
	return &ec2.DescribeAddressesOutput{Addresses: m.addresss}, nil
}

func (m *mockEc2) DescribeSnapshots(ctx context.Context, input *ec2.DescribeSnapshotsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSnapshotsOutput, error) {
	return &ec2.DescribeSnapshotsOutput{Snapshots: m.snapshots}, nil
}

func (m *mockEc2) DescribeNetworkInterfaces(ctx context.Context, input *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	return &ec2.DescribeNetworkInterfacesOutput{NetworkInterfaces: m.networkinterfaces}, nil
}

type mockElbv2 struct {
	loadbalancers            []elbv2types.LoadBalancer
	targetgroups             []elbv2types.TargetGroup
	listeners                []elbv2types.Listener
	targethealthdescriptions map[string][]elbv2types.TargetHealthDescription
}

func (m *mockElbv2) Name() string {
	return ""
}

func (m *mockElbv2) Region() string {
	return ""
}

func (m *mockElbv2) Profile() string {
	return ""
}

func (m *mockElbv2) Provider() string {
	return ""
}

func (m *mockElbv2) ProviderAPI() string {
	return ""
}

func (m *mockElbv2) ResourceTypes() []string {
	return []string{}
}

func (m *mockElbv2) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElbv2) IsSyncDisabled() bool {
	return false
}

func (m *mockElbv2) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElbv2) DescribeLoadBalancers(ctx context.Context, input *elbv2.DescribeLoadBalancersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeLoadBalancersOutput, error) {
	return &elbv2.DescribeLoadBalancersOutput{LoadBalancers: m.loadbalancers}, nil
}

func (m *mockElbv2) DescribeTargetGroups(ctx context.Context, input *elbv2.DescribeTargetGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error) {
	return &elbv2.DescribeTargetGroupsOutput{TargetGroups: m.targetgroups}, nil
}

type mockElb struct {
	loadbalancerdescriptions []elbtypes.LoadBalancerDescription
}

func (m *mockElb) Name() string {
	return ""
}

func (m *mockElb) Region() string {
	return ""
}

func (m *mockElb) Profile() string {
	return ""
}

func (m *mockElb) Provider() string {
	return ""
}

func (m *mockElb) ProviderAPI() string {
	return ""
}

func (m *mockElb) ResourceTypes() []string {
	return []string{}
}

func (m *mockElb) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElb) IsSyncDisabled() bool {
	return false
}

func (m *mockElb) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockElb) DescribeLoadBalancers(ctx context.Context, input *elb.DescribeLoadBalancersInput, optFns ...func(*elb.Options)) (*elb.DescribeLoadBalancersOutput, error) {
	return &elb.DescribeLoadBalancersOutput{LoadBalancerDescriptions: m.loadbalancerdescriptions}, nil
}

type mockRds struct {
	dbinstances    []rdstypes.DBInstance
	dbsubnetgroups []rdstypes.DBSubnetGroup
}

func (m *mockRds) Name() string {
	return ""
}

func (m *mockRds) Region() string {
	return ""
}

func (m *mockRds) Profile() string {
	return ""
}

func (m *mockRds) Provider() string {
	return ""
}

func (m *mockRds) ProviderAPI() string {
	return ""
}

func (m *mockRds) ResourceTypes() []string {
	return []string{}
}

func (m *mockRds) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRds) IsSyncDisabled() bool {
	return false
}

func (m *mockRds) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRds) DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	return &rds.DescribeDBInstancesOutput{DBInstances: m.dbinstances}, nil
}

func (m *mockRds) DescribeDBSubnetGroups(ctx context.Context, input *rds.DescribeDBSubnetGroupsInput, optFns ...func(*rds.Options)) (*rds.DescribeDBSubnetGroupsOutput, error) {
	return &rds.DescribeDBSubnetGroupsOutput{DBSubnetGroups: m.dbsubnetgroups}, nil
}

type mockAutoscaling struct {
	launchconfigurations []autoscalingtypes.LaunchConfiguration
	autoscalinggroups    []autoscalingtypes.AutoScalingGroup
	scalingpolicys       []autoscalingtypes.ScalingPolicy
}

func (m *mockAutoscaling) Name() string {
	return ""
}

func (m *mockAutoscaling) Region() string {
	return ""
}

func (m *mockAutoscaling) Profile() string {
	return ""
}

func (m *mockAutoscaling) Provider() string {
	return ""
}

func (m *mockAutoscaling) ProviderAPI() string {
	return ""
}

func (m *mockAutoscaling) ResourceTypes() []string {
	return []string{}
}

func (m *mockAutoscaling) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAutoscaling) IsSyncDisabled() bool {
	return false
}

func (m *mockAutoscaling) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAutoscaling) DescribeLaunchConfigurations(ctx context.Context, input *autoscaling.DescribeLaunchConfigurationsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeLaunchConfigurationsOutput, error) {
	return &autoscaling.DescribeLaunchConfigurationsOutput{LaunchConfigurations: m.launchconfigurations}, nil
}

func (m *mockAutoscaling) DescribeAutoScalingGroups(ctx context.Context, input *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	return &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: m.autoscalinggroups}, nil
}

func (m *mockAutoscaling) DescribePolicies(ctx context.Context, input *autoscaling.DescribePoliciesInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribePoliciesOutput, error) {
	return &autoscaling.DescribePoliciesOutput{ScalingPolicies: m.scalingpolicys}, nil
}

type mockAcm struct {
	certificatesummarys []acmtypes.CertificateSummary
}

func (m *mockAcm) Name() string {
	return ""
}

func (m *mockAcm) Region() string {
	return ""
}

func (m *mockAcm) Profile() string {
	return ""
}

func (m *mockAcm) Provider() string {
	return ""
}

func (m *mockAcm) ProviderAPI() string {
	return ""
}

func (m *mockAcm) ResourceTypes() []string {
	return []string{}
}

func (m *mockAcm) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAcm) IsSyncDisabled() bool {
	return false
}

func (m *mockAcm) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockAcm) ListCertificates(ctx context.Context, input *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	return &acm.ListCertificatesOutput{CertificateSummaryList: m.certificatesummarys}, nil
}

type mockIam struct {
	userdetails          []iamtypes.UserDetail
	groupdetails         []iamtypes.GroupDetail
	roledetails          []iamtypes.RoleDetail
	policys              []iamtypes.Policy
	accesskeymetadatas   []iamtypes.AccessKeyMetadata
	instanceprofiles     []iamtypes.InstanceProfile
	managedpolicydetails []iamtypes.ManagedPolicyDetail
	users                []iamtypes.User
	virtualmfadevices    []iamtypes.VirtualMFADevice
}

func (m *mockIam) Name() string {
	return ""
}

func (m *mockIam) Region() string {
	return ""
}

func (m *mockIam) Profile() string {
	return ""
}

func (m *mockIam) Provider() string {
	return ""
}

func (m *mockIam) ProviderAPI() string {
	return ""
}

func (m *mockIam) ResourceTypes() []string {
	return []string{}
}

func (m *mockIam) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockIam) IsSyncDisabled() bool {
	return false
}

func (m *mockIam) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockIam) ListAccessKeys(ctx context.Context, input *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
	return &iam.ListAccessKeysOutput{AccessKeyMetadata: m.accesskeymetadatas}, nil
}

func (m *mockIam) ListInstanceProfiles(ctx context.Context, input *iam.ListInstanceProfilesInput, optFns ...func(*iam.Options)) (*iam.ListInstanceProfilesOutput, error) {
	return &iam.ListInstanceProfilesOutput{InstanceProfiles: m.instanceprofiles}, nil
}

func (m *mockIam) ListVirtualMFADevices(ctx context.Context, input *iam.ListVirtualMFADevicesInput, optFns ...func(*iam.Options)) (*iam.ListVirtualMFADevicesOutput, error) {
	return &iam.ListVirtualMFADevicesOutput{VirtualMFADevices: m.virtualmfadevices}, nil
}

type mockS3 struct {
	buckets map[string][]s3types.Bucket
	objects map[string][]s3types.Object
	grants  map[string][]s3types.Grant
}

func (m *mockS3) Name() string {
	return ""
}

func (m *mockS3) Region() string {
	return ""
}

func (m *mockS3) Profile() string {
	return ""
}

func (m *mockS3) Provider() string {
	return ""
}

func (m *mockS3) ProviderAPI() string {
	return ""
}

func (m *mockS3) ResourceTypes() []string {
	return []string{}
}

func (m *mockS3) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockS3) IsSyncDisabled() bool {
	return false
}

func (m *mockS3) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

type mockSns struct {
	subscriptions []snstypes.Subscription
	topics        []snstypes.Topic
}

func (m *mockSns) Name() string {
	return ""
}

func (m *mockSns) Region() string {
	return ""
}

func (m *mockSns) Profile() string {
	return ""
}

func (m *mockSns) Provider() string {
	return ""
}

func (m *mockSns) ProviderAPI() string {
	return ""
}

func (m *mockSns) ResourceTypes() []string {
	return []string{}
}

func (m *mockSns) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSns) IsSyncDisabled() bool {
	return false
}

func (m *mockSns) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSns) ListSubscriptions(ctx context.Context, input *sns.ListSubscriptionsInput, optFns ...func(*sns.Options)) (*sns.ListSubscriptionsOutput, error) {
	return &sns.ListSubscriptionsOutput{Subscriptions: m.subscriptions}, nil
}

func (m *mockSns) ListTopics(ctx context.Context, input *sns.ListTopicsInput, optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error) {
	return &sns.ListTopicsOutput{Topics: m.topics}, nil
}

type mockSqs struct {
	strings    []string
	attributes map[string]map[string]string
}

func (m *mockSqs) Name() string {
	return ""
}

func (m *mockSqs) Region() string {
	return ""
}

func (m *mockSqs) Profile() string {
	return ""
}

func (m *mockSqs) Provider() string {
	return ""
}

func (m *mockSqs) ProviderAPI() string {
	return ""
}

func (m *mockSqs) ResourceTypes() []string {
	return []string{}
}

func (m *mockSqs) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSqs) IsSyncDisabled() bool {
	return false
}

func (m *mockSqs) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockSqs) ListQueues(ctx context.Context, input *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error) {
	return &sqs.ListQueuesOutput{QueueUrls: m.strings}, nil
}

type mockRoute53 struct {
	hostedzones        []route53types.HostedZone
	resourcerecordsets map[string][]route53types.ResourceRecordSet
}

func (m *mockRoute53) Name() string {
	return ""
}

func (m *mockRoute53) Region() string {
	return ""
}

func (m *mockRoute53) Profile() string {
	return ""
}

func (m *mockRoute53) Provider() string {
	return ""
}

func (m *mockRoute53) ProviderAPI() string {
	return ""
}

func (m *mockRoute53) ResourceTypes() []string {
	return []string{}
}

func (m *mockRoute53) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRoute53) IsSyncDisabled() bool {
	return false
}

func (m *mockRoute53) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockRoute53) ListHostedZones(ctx context.Context, input *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
	return &route53.ListHostedZonesOutput{HostedZones: m.hostedzones}, nil
}

type mockLambda struct {
	functionconfigurations []lambdatypes.FunctionConfiguration
}

func (m *mockLambda) Name() string {
	return ""
}

func (m *mockLambda) Region() string {
	return ""
}

func (m *mockLambda) Profile() string {
	return ""
}

func (m *mockLambda) Provider() string {
	return ""
}

func (m *mockLambda) ProviderAPI() string {
	return ""
}

func (m *mockLambda) ResourceTypes() []string {
	return []string{}
}

func (m *mockLambda) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockLambda) IsSyncDisabled() bool {
	return false
}

func (m *mockLambda) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockLambda) ListFunctions(ctx context.Context, input *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	return &lambda.ListFunctionsOutput{Functions: m.functionconfigurations}, nil
}

type mockCloudwatch struct {
	metrics      []cloudwatchtypes.Metric
	metricalarms []cloudwatchtypes.MetricAlarm
}

func (m *mockCloudwatch) Name() string {
	return ""
}

func (m *mockCloudwatch) Region() string {
	return ""
}

func (m *mockCloudwatch) Profile() string {
	return ""
}

func (m *mockCloudwatch) Provider() string {
	return ""
}

func (m *mockCloudwatch) ProviderAPI() string {
	return ""
}

func (m *mockCloudwatch) ResourceTypes() []string {
	return []string{}
}

func (m *mockCloudwatch) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudwatch) IsSyncDisabled() bool {
	return false
}

func (m *mockCloudwatch) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudwatch) ListMetrics(ctx context.Context, input *cloudwatch.ListMetricsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.ListMetricsOutput, error) {
	return &cloudwatch.ListMetricsOutput{Metrics: m.metrics}, nil
}

func (m *mockCloudwatch) DescribeAlarms(ctx context.Context, input *cloudwatch.DescribeAlarmsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.DescribeAlarmsOutput, error) {
	return &cloudwatch.DescribeAlarmsOutput{MetricAlarms: m.metricalarms}, nil
}

type mockCloudfront struct {
	distributionsummarys []cloudfronttypes.DistributionSummary
}

func (m *mockCloudfront) Name() string {
	return ""
}

func (m *mockCloudfront) Region() string {
	return ""
}

func (m *mockCloudfront) Profile() string {
	return ""
}

func (m *mockCloudfront) Provider() string {
	return ""
}

func (m *mockCloudfront) ProviderAPI() string {
	return ""
}

func (m *mockCloudfront) ResourceTypes() []string {
	return []string{}
}

func (m *mockCloudfront) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudfront) IsSyncDisabled() bool {
	return false
}

func (m *mockCloudfront) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

type mockCloudformation struct {
	stacks []cloudformationtypes.Stack
}

func (m *mockCloudformation) Name() string {
	return ""
}

func (m *mockCloudformation) Region() string {
	return ""
}

func (m *mockCloudformation) Profile() string {
	return ""
}

func (m *mockCloudformation) Provider() string {
	return ""
}

func (m *mockCloudformation) ProviderAPI() string {
	return ""
}

func (m *mockCloudformation) ResourceTypes() []string {
	return []string{}
}

func (m *mockCloudformation) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudformation) IsSyncDisabled() bool {
	return false
}

func (m *mockCloudformation) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockCloudformation) DescribeStacks(ctx context.Context, input *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return &cloudformation.DescribeStacksOutput{Stacks: m.stacks}, nil
}

type mockEcr struct {
	repositorys []ecrtypes.Repository
}

func (m *mockEcr) Name() string {
	return ""
}

func (m *mockEcr) Region() string {
	return ""
}

func (m *mockEcr) Profile() string {
	return ""
}

func (m *mockEcr) Provider() string {
	return ""
}

func (m *mockEcr) ProviderAPI() string {
	return ""
}

func (m *mockEcr) ResourceTypes() []string {
	return []string{}
}

func (m *mockEcr) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcr) IsSyncDisabled() bool {
	return false
}

func (m *mockEcr) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcr) DescribeRepositories(ctx context.Context, input *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error) {
	return &ecr.DescribeRepositoriesOutput{Repositories: m.repositorys}, nil
}

type mockEcs struct {
	clusters                []ecstypes.Cluster
	clusterNames            []string
	taskdefinitions         []ecstypes.TaskDefinition
	taskdefinitionNames     []string
	tasks                   map[string][]ecstypes.Task
	tasksNames              map[string][]string
	containerinstancesNames map[string][]string
	containerinstances      map[string][]ecstypes.ContainerInstance
}

func (m *mockEcs) Name() string {
	return ""
}

func (m *mockEcs) Region() string {
	return ""
}

func (m *mockEcs) Profile() string {
	return ""
}

func (m *mockEcs) Provider() string {
	return ""
}

func (m *mockEcs) ProviderAPI() string {
	return ""
}

func (m *mockEcs) ResourceTypes() []string {
	return []string{}
}

func (m *mockEcs) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcs) IsSyncDisabled() bool {
	return false
}

func (m *mockEcs) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m *mockEcs) ListClusters(ctx context.Context, input *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error) {
	return &ecs.ListClustersOutput{ClusterArns: m.clusterNames}, nil
}

func (m *mockEcs) ListTaskDefinitions(ctx context.Context, input *ecs.ListTaskDefinitionsInput, optFns ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error) {
	return &ecs.ListTaskDefinitionsOutput{TaskDefinitionArns: m.taskdefinitionNames}, nil
}

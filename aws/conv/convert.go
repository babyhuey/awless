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

package awsconv

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
	"net"
	"net/url"
	"reflect"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/graph"
)

func InitResource(source interface{}) (*graph.Resource, error) {
	var res *graph.Resource
	switch ss := source.(type) {
	// EC2
	case ec2types.Instance:
		res = graph.InitResource(cloud.Instance, awssdk.ToString(ss.InstanceId))
	case ec2types.Vpc:
		res = graph.InitResource(cloud.Vpc, awssdk.ToString(ss.VpcId))
	case ec2types.Subnet:
		res = graph.InitResource(cloud.Subnet, awssdk.ToString(ss.SubnetId))
	case ec2types.SecurityGroup:
		res = graph.InitResource(cloud.SecurityGroup, awssdk.ToString(ss.GroupId))
	case ec2types.KeyPairInfo:
		res = graph.InitResource(cloud.Keypair, awssdk.ToString(ss.KeyName))
	case ec2types.Volume:
		res = graph.InitResource(cloud.Volume, awssdk.ToString(ss.VolumeId))
	case ec2types.Image:
		res = graph.InitResource(cloud.Image, awssdk.ToString(ss.ImageId))
	case ec2types.ImportImageTask:
		res = graph.InitResource(cloud.ImportImageTask, awssdk.ToString(ss.ImportTaskId))
	case ec2types.InternetGateway:
		res = graph.InitResource(cloud.InternetGateway, awssdk.ToString(ss.InternetGatewayId))
	case ec2types.NatGateway:
		res = graph.InitResource(cloud.NatGateway, awssdk.ToString(ss.NatGatewayId))
	case ec2types.RouteTable:
		res = graph.InitResource(cloud.RouteTable, awssdk.ToString(ss.RouteTableId))
	case ec2types.AvailabilityZone:
		res = graph.InitResource(cloud.AvailabilityZone, awssdk.ToString(ss.ZoneName))
	case ec2types.Address:
		if awssdk.ToString(ss.AllocationId) != "" {
			res = graph.InitResource(cloud.ElasticIP, awssdk.ToString(ss.AllocationId))
		} else {
			res = graph.InitResource(cloud.ElasticIP, awssdk.ToString(ss.PublicIp))
		}
	case ec2types.Snapshot:
		res = graph.InitResource(cloud.Snapshot, awssdk.ToString(ss.SnapshotId))
	case ec2types.NetworkInterface:
		res = graph.InitResource(cloud.NetworkInterface, awssdk.ToString(ss.NetworkInterfaceId))
	// Loadbalancer
	case elbtypes.LoadBalancerDescription:
		res = graph.InitResource(cloud.ClassicLoadBalancer, awssdk.ToString(ss.LoadBalancerName))
	case elbv2types.LoadBalancer:
		res = graph.InitResource(cloud.LoadBalancer, awssdk.ToString(ss.LoadBalancerArn))
	case elbv2types.TargetGroup:
		res = graph.InitResource(cloud.TargetGroup, awssdk.ToString(ss.TargetGroupArn))
	case elbv2types.Listener:
		res = graph.InitResource(cloud.Listener, awssdk.ToString(ss.ListenerArn))
		// Database
	case rdstypes.DBInstance:
		res = graph.InitResource(cloud.Database, awssdk.ToString(ss.DBInstanceIdentifier))
	case rdstypes.DBSubnetGroup:
		res = graph.InitResource(cloud.DbSubnetGroup, awssdk.ToString(ss.DBSubnetGroupArn))
		// Autoscaling
	case autoscalingtypes.LaunchConfiguration:
		res = graph.InitResource(cloud.LaunchConfiguration, awssdk.ToString(ss.LaunchConfigurationARN))
	case autoscalingtypes.AutoScalingGroup:
		res = graph.InitResource(cloud.ScalingGroup, awssdk.ToString(ss.AutoScalingGroupARN))
	case autoscalingtypes.ScalingPolicy:
		res = graph.InitResource(cloud.ScalingPolicy, awssdk.ToString(ss.PolicyARN))
	// Container
	case ecrtypes.Repository:
		res = graph.InitResource(cloud.Repository, awssdk.ToString(ss.RepositoryArn))
	case ecstypes.Cluster:
		res = graph.InitResource(cloud.ContainerCluster, awssdk.ToString(ss.ClusterArn))
	case ecstypes.TaskDefinition:
		res = graph.InitResource(cloud.ContainerTask, awssdk.ToString(ss.TaskDefinitionArn))
	case ecstypes.Container:
		res = graph.InitResource(cloud.Container, awssdk.ToString(ss.ContainerArn))
	case ecstypes.ContainerInstance:
		res = graph.InitResource(cloud.ContainerInstance, awssdk.ToString(ss.ContainerInstanceArn))
		// ACM
	case acmtypes.CertificateSummary:
		res = graph.InitResource(cloud.Certificate, awssdk.ToString(ss.CertificateArn))
	// IAM
	case iamtypes.User:
		res = graph.InitResource(cloud.User, awssdk.ToString(ss.UserId))
	case iamtypes.UserDetail:
		res = graph.InitResource(cloud.User, awssdk.ToString(ss.UserId))
	case iamtypes.RoleDetail:
		res = graph.InitResource(cloud.Role, awssdk.ToString(ss.RoleId))
	case iamtypes.GroupDetail:
		res = graph.InitResource(cloud.Group, awssdk.ToString(ss.GroupId))
	case iamtypes.Policy:
		res = graph.InitResource(cloud.Policy, awssdk.ToString(ss.PolicyId))
	case iamtypes.ManagedPolicyDetail:
		res = graph.InitResource(cloud.Policy, awssdk.ToString(ss.PolicyId))
	case iamtypes.AccessKeyMetadata:
		res = graph.InitResource(cloud.AccessKey, awssdk.ToString(ss.AccessKeyId))
	case iamtypes.InstanceProfile:
		res = graph.InitResource(cloud.InstanceProfile, awssdk.ToString(ss.InstanceProfileId))
	case iamtypes.VirtualMFADevice:
		res = graph.InitResource(cloud.MFADevice, awssdk.ToString(ss.SerialNumber))
	// S3
	case s3types.Bucket:
		res = graph.InitResource(cloud.Bucket, awssdk.ToString(ss.Name))
	case s3types.Object:
		res = graph.InitResource(cloud.S3Object, awssdk.ToString(ss.Key))
	//SNS
	case snstypes.Subscription:
		res = graph.InitResource(cloud.Subscription, awssdk.ToString(ss.Endpoint))
	case snstypes.Topic:
		res = graph.InitResource(cloud.Topic, awssdk.ToString(ss.TopicArn))
		// DNS
	case route53types.HostedZone:
		res = graph.InitResource(cloud.Zone, awssdk.ToString(ss.Id))
	case route53types.ResourceRecordSet:
		id := HashFields(awssdk.ToString(ss.Name), string(ss.Type))
		res = graph.InitResource(cloud.Record, id)
		// Lambda
	case lambdatypes.FunctionConfiguration:
		res = graph.InitResource(cloud.Function, awssdk.ToString(ss.FunctionArn))
		// Monitoring
	case cloudwatchtypes.Metric:
		id := HashFields(awssdk.ToString(ss.Namespace), awssdk.ToString(ss.MetricName))
		res = graph.InitResource(cloud.Metric, id)
	case cloudwatchtypes.MetricAlarm:
		res = graph.InitResource(cloud.Alarm, awssdk.ToString(ss.AlarmArn))
		// cdn
	case cloudfronttypes.DistributionSummary:
		res = graph.InitResource(cloud.Distribution, awssdk.ToString(ss.Id))
		// cloudformation
	case cloudformationtypes.Stack:
		res = graph.InitResource(cloud.Stack, awssdk.ToString(ss.StackId))
	// EKS
	case ekstypes.Cluster:
		res = graph.InitResource(cloud.EKSCluster, awssdk.ToString(ss.Name))
	case ekstypes.Nodegroup:
		res = graph.InitResource(cloud.EKSNodeGroup, awssdk.ToString(ss.NodegroupArn))
	// DynamoDB
	case dynamodbtypes.TableDescription:
		res = graph.InitResource(cloud.DynamoDBTable, awssdk.ToString(ss.TableName))
	// Secrets Manager
	case secretsmanagertypes.SecretListEntry:
		res = graph.InitResource(cloud.Secret, awssdk.ToString(ss.ARN))
	// KMS
	case kmstypes.KeyMetadata:
		res = graph.InitResource(cloud.Key, awssdk.ToString(ss.KeyId))
	// API Gateway
	case apigatewayv2types.Api:
		res = graph.InitResource(cloud.ApiGateway, awssdk.ToString(ss.ApiId))
	case apigatewayv2types.Route:
		res = graph.InitResource(cloud.ApiGatewayRoute, awssdk.ToString(ss.RouteId))
	case apigatewayv2types.Stage:
		res = graph.InitResource(cloud.ApiGatewayStage, awssdk.ToString(ss.StageName))
	// SSM
	case ssmtypes.ParameterMetadata:
		res = graph.InitResource(cloud.SSMParameter, awssdk.ToString(ss.Name))
	// EFS
	case efstypes.FileSystemDescription:
		res = graph.InitResource(cloud.FileSystem, awssdk.ToString(ss.FileSystemId))
	case efstypes.MountTargetDescription:
		res = graph.InitResource(cloud.MountTarget, awssdk.ToString(ss.MountTargetId))
	// CloudTrail
	case cloudtrailtypes.Trail:
		res = graph.InitResource(cloud.Trail, awssdk.ToString(ss.TrailARN))
	// CloudWatch Logs
	case cloudwatchlogstypes.LogGroup:
		res = graph.InitResource(cloud.LogGroup, awssdk.ToString(ss.LogGroupName))
	default:
		return nil, fmt.Errorf("Unknown type of resource %T", source)
	}
	return res, nil
}

func NewResource(source interface{}) (*graph.Resource, error) {
	res, err := InitResource(source)
	if err != nil {
		return res, err
	}

	res.Properties()[properties.ID] = res.Id()

	value := reflect.ValueOf(source)
	if !value.IsValid() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not valid.", value)
	}
	var nodeV reflect.Value
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("can not fetch cloud resource. %v is a nil pointer.", value)
		}
		nodeV = value.Elem()
	} else if value.Kind() == reflect.Struct {
		nodeV = value
	} else {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid struct or pointer.", value)
	}

	// Bounded and leak-free. The previous version wrote to unbuffered resultc
	// and errc channels while the consumer returned on the first error, leaving
	// remaining goroutines blocked on send forever. It also continued to send a
	// result after sending an error, rather than stopping.
	var mu sync.Mutex
	g := new(errgroup.Group)
	g.SetLimit(maxParallelPropertyTransforms)

	for prop, trans := range awsResourcesDef[res.Type()] {
		p, t := prop, trans
		g.Go(func() error {
			setProp := func(val interface{}) {
				mu.Lock()
				res.Properties()[p] = val
				mu.Unlock()
			}

			if t.transform != nil {
				sourceField := nodeV.FieldByName(t.name)
				if sourceField.IsValid() && !isNilValue(sourceField) {
					val, err := t.transform(sourceField.Interface())
					if err == ErrTagNotFound {
						return nil
					}
					if err != nil {
						return fmt.Errorf("type [%s]: prop '%v': %w", res.Type(), p, err)
					}
					setProp(val)
				}
			}
			if t.fetch != nil {
				val, err := t.fetch(source)
				if err != nil {
					return fmt.Errorf("type [%s]: prop '%v': %w", res.Type(), p, err)
				}
				setProp(val)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return res, err
	}

	return res, nil
}

// maxParallelPropertyTransforms bounds the per-property fan-out when
// converting one AWS API object into a graph resource.
const maxParallelPropertyTransforms = 10

var ErrTagNotFound = errors.New("aws tag key not found")

type propertyTransform struct {
	name      string
	transform transformFn
	fetch     fetchFn
}

type transformFn func(i interface{}) (interface{}, error)
type fetchFn func(i interface{}) (interface{}, error)

var extractValueFn = func(i interface{}) (interface{}, error) {
	iv := reflect.ValueOf(i)
	if iv.Kind() == reflect.Ptr {
		if iv.IsNil() {
			return nil, nil
		}
		return iv.Elem().Interface(), nil
	}
	switch ii := i.(type) {
	case []string:
		return ii, nil
	case []*string:
		return awssdk.ToStringSlice(ii), nil
	case string:
		return ii, nil
	case int32:
		return ii, nil
	case int64:
		return ii, nil
	case bool:
		return ii, nil
	}
	// Handle typed enums (e.g. ec2types.InstanceStateName) which are underlying string/int types
	if iv.Kind() == reflect.String {
		return iv.String(), nil
	}
	return nil, fmt.Errorf("extract value: not a pointer or known type but a %T", i)
}

var extractValueAsStringFn = func(i interface{}) (interface{}, error) {
	val, err := extractValueFn(i)
	return fmt.Sprint(val), err
}

// Extract time forcing timezone to UTC (friendlier when running test in different timezones i.e. travis)
var extractTimeFn = func(i interface{}) (interface{}, error) {
	t, ok := i.(*time.Time)
	if ok {
		return t.UTC(), nil
	}
	s, ok := i.(string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05.000+0000", s)
		if err != nil {
			return nil, err
		}
		return t.UTC(), nil
	}
	sp, ok := i.(*string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05.000+0000", awssdk.ToString(sp))
		if err != nil {
			return nil, err
		}
		return t.UTC(), nil
	}
	return nil, fmt.Errorf("extract time: expected time pointer, got: %T", i)
}

// Extract time from *int64 representing milliseconds since epoch (e.g., CloudWatch Logs CreationTime)
var extractMillisecondEpochTimeFn = func(i interface{}) (interface{}, error) {
	p, ok := i.(*int64)
	if ok {
		if p == nil {
			return nil, nil
		}
		return time.UnixMilli(*p).UTC(), nil
	}
	v, ok := i.(int64)
	if ok {
		return time.UnixMilli(v).UTC(), nil
	}
	return nil, fmt.Errorf("extract millis epoch time: expected *int64 or int64, got: %T", i)
}

// Extract time that have a Z directly after the time without a space which means UTC
// (https://en.wikipedia.org/wiki/ISO_8601#UTC)
var extractTimeWithZSuffixFn = func(i interface{}) (interface{}, error) {
	t, ok := i.(*time.Time)
	if ok {
		return t.UTC(), nil
	}
	s, ok := i.(string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", s)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	sp, ok := i.(*string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", awssdk.ToString(sp))
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, fmt.Errorf("extract time: expected time pointer, got: %T", i)
}

var extractIpPermissionSliceFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]ec2types.IpPermission); !ok {
		return nil, fmt.Errorf("extract ip permission: not a permission slice but a %T", i)
	}
	var rules []*graph.FirewallRule
	for _, ipPerm := range i.([]ec2types.IpPermission) {
		rule := &graph.FirewallRule{}

		protocol := awssdk.ToString(ipPerm.IpProtocol)
		switch protocol {
		case "-1":
			rule.Protocol = "any"
			rule.PortRange = graph.PortRange{Any: true}
		case "tcp", "udp", "icmp", "58":
			rule.Protocol = protocol
			fromPort := int64(awssdk.ToInt32(ipPerm.FromPort))
			toPort := int64(awssdk.ToInt32(ipPerm.ToPort))
			if fromPort == -1 || toPort == -1 {
				rule.PortRange = graph.PortRange{Any: true}
			} else {
				rule.PortRange = graph.PortRange{FromPort: fromPort, ToPort: toPort}
			}

		default:
			rule.Protocol = protocol
			rule.PortRange = graph.PortRange{Any: true}
		}
		for _, r := range ipPerm.IpRanges {
			_, net, err := net.ParseCIDR(awssdk.ToString(r.CidrIp))
			if err != nil {
				return rules, err
			}
			rule.IPRanges = append(rule.IPRanges, net)
		}
		for _, r := range ipPerm.Ipv6Ranges {
			_, net, err := net.ParseCIDR(awssdk.ToString(r.CidrIpv6))
			if err != nil {
				return rules, err
			}
			rule.IPRanges = append(rule.IPRanges, net)
		}
		for _, group := range ipPerm.UserIdGroupPairs {
			rule.Sources = append(rule.Sources, awssdk.ToString(group.GroupId))
		}

		rules = append(rules, rule)
	}
	return rules, nil

}

var extractNameValueFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]cloudwatchtypes.Dimension); !ok {
		return nil, fmt.Errorf("extract ip namevalue: not a dimension slice but a %T", i)
	}
	var nameValues []*graph.KeyValue
	for _, dimension := range i.([]cloudwatchtypes.Dimension) {
		keyval := &graph.KeyValue{KeyName: awssdk.ToString(dimension.Name), Value: awssdk.ToString(dimension.Value)}

		nameValues = append(nameValues, keyval)
	}
	return nameValues, nil
}

var extractECSAttributesFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]ecstypes.Attribute); !ok {
		return nil, fmt.Errorf("extract ECS attributes: not an attribute slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, attribute := range i.([]ecstypes.Attribute) {
		keyval := &graph.KeyValue{KeyName: awssdk.ToString(attribute.Name), Value: awssdk.ToString(attribute.Value)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractRouteTableAssociationsFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]ec2types.RouteTableAssociation); !ok {
		return nil, fmt.Errorf("extract route table associations: not an association slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, assoc := range i.([]ec2types.RouteTableAssociation) {
		keyval := &graph.KeyValue{KeyName: awssdk.ToString(assoc.RouteTableAssociationId), Value: awssdk.ToString(assoc.SubnetId)}
		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractFieldFn = func(field string) transformFn {
	return func(i interface{}) (interface{}, error) {
		value := reflect.ValueOf(i)
		var struc reflect.Value
		if value.Kind() == reflect.Ptr {
			struc = value.Elem()
			if struc.Kind() != reflect.Struct {
				return nil, fmt.Errorf("extract field '%s': not a struct pointer but a %T", field, i)
			}
		} else if value.Kind() == reflect.Struct {
			struc = value
		} else {
			return nil, fmt.Errorf("extract field '%s': not a pointer or struct but a %T", field, i)
		}

		structField := struc.FieldByName(field)

		if !structField.IsValid() {
			return nil, fmt.Errorf("extract field: field not found: %s", field)
		}

		return extractValueFn(structField.Interface())
	}
}

var extractTagsFn = func(i interface{}) (interface{}, error) {
	var out []string
	switch tags := i.(type) {
	case []ec2types.Tag:
		for _, t := range tags {
			out = append(out, fmt.Sprintf("%s=%s", awssdk.ToString(t.Key), awssdk.ToString(t.Value)))
		}
	case []autoscalingtypes.TagDescription:
		for _, t := range tags {
			out = append(out, fmt.Sprintf("%s=%s", awssdk.ToString(t.Key), awssdk.ToString(t.Value)))
		}
	default:
		return nil, fmt.Errorf("extract tags: not a tag slice, but a %T", i)
	}

	return out, nil
}

var extractMapTagsFn = func(i interface{}) (interface{}, error) {
	tags, ok := i.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("extract map tags: not a map[string]string, but a %T", i)
	}
	var out []string
	for k, v := range tags {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out, nil
}

var extractSecretsManagerTagsFn = func(i interface{}) (interface{}, error) {
	tags, ok := i.([]secretsmanagertypes.Tag)
	if !ok {
		return nil, fmt.Errorf("extract secrets manager tags: not a tag slice, but a %T", i)
	}
	var out []string
	for _, t := range tags {
		out = append(out, fmt.Sprintf("%s=%s", awssdk.ToString(t.Key), awssdk.ToString(t.Value)))
	}
	return out, nil
}

var extractEFSTagsFn = func(i interface{}) (interface{}, error) {
	tags, ok := i.([]efstypes.Tag)
	if !ok {
		return nil, fmt.Errorf("extract EFS tags: not a tag slice, but a %T", i)
	}
	var out []string
	for _, t := range tags {
		out = append(out, fmt.Sprintf("%s=%s", awssdk.ToString(t.Key), awssdk.ToString(t.Value)))
	}
	return out, nil
}

var extractEFSTagFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		tags, ok := i.([]efstypes.Tag)
		if !ok {
			return nil, fmt.Errorf("extract EFS tag: not a tag slice, but a %T", i)
		}
		for _, t := range tags {
			if key == awssdk.ToString(t.Key) {
				return awssdk.ToString(t.Value), nil
			}
		}
		return nil, ErrTagNotFound
	}
}

var extractDynamoDBKeySchemaFn = func(i interface{}) (interface{}, error) {
	keys, ok := i.([]dynamodbtypes.KeySchemaElement)
	if !ok {
		return nil, fmt.Errorf("extract key schema: not a KeySchemaElement slice, but a %T", i)
	}
	var out []string
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s:%s", awssdk.ToString(k.AttributeName), k.KeyType))
	}
	return out, nil
}

var extractEKSScalingConfigFn = func(i interface{}) (interface{}, error) {
	sc, ok := i.(*ekstypes.NodegroupScalingConfig)
	if !ok {
		return nil, fmt.Errorf("extract scaling config: not a NodegroupScalingConfig pointer, but a %T", i)
	}
	if sc == nil {
		return nil, nil
	}
	return fmt.Sprintf("min:%d,max:%d,desired:%d", awssdk.ToInt32(sc.MinSize), awssdk.ToInt32(sc.MaxSize), awssdk.ToInt32(sc.DesiredSize)), nil
}

var extractTagFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		tags, ok := i.([]ec2types.Tag)
		if !ok {
			return nil, fmt.Errorf("extract tag: not a tag slice, but a %T", i)
		}
		for _, t := range tags {
			if key == awssdk.ToString(t.Key) {
				return awssdk.ToString(t.Value), nil
			}
		}

		return nil, ErrTagNotFound
	}
}

var extractStringPointerSliceValues = func(i interface{}) (interface{}, error) {
	switch ss := i.(type) {
	case []string:
		return ss, nil
	case []*string:
		return awssdk.ToStringSlice(ss), nil
	}
	return nil, fmt.Errorf("extract string pointer: not a string slice but a %T", i)
}

var extractStringSliceValues = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		var res []string
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("extract slice: not a slice but a %T", i)
		}
		for i := 0; i < value.Len(); i++ {
			e, err := extractFieldFn(key)(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			str, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("extract string slice: not a string but a %T", e)
			}
			res = append(res, str)
		}

		return res, nil
	}
}

var extractClassicLoadbListenerDescriptionsFn = func(i interface{}) (interface{}, error) {
	listeners, ok := i.([]elbtypes.ListenerDescription)
	if !ok {
		return nil, fmt.Errorf("extract classic loadb listener descriptions: unexpected type %T", i)
	}
	var out []string
	for _, d := range listeners {
		if list := d.Listener; list != nil {
			out = append(out, fmt.Sprintf("%s:%d:%s:%d", awssdk.ToString(list.Protocol), list.LoadBalancerPort, awssdk.ToString(list.InstanceProtocol), awssdk.ToInt32(list.InstancePort)))
		}
	}
	return out, nil
}

var extractRoutesSliceFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]ec2types.Route); !ok {
		return nil, fmt.Errorf("extract route: not a route slice but a %T", i)
	}
	var routes []*graph.Route
	for _, r := range i.([]ec2types.Route) {
		route := &graph.Route{}
		var err error
		if notEmptyStr(r.DestinationCidrBlock) {
			if _, route.Destination, err = net.ParseCIDR(awssdk.ToString(r.DestinationCidrBlock)); err != nil {
				return nil, err
			}
		}
		if notEmptyStr(r.DestinationIpv6CidrBlock) {
			if _, route.DestinationIPv6, err = net.ParseCIDR(awssdk.ToString(r.DestinationIpv6CidrBlock)); err != nil {
				return nil, err
			}
		}
		if notEmptyStr(r.DestinationPrefixListId) {
			route.DestinationPrefixListId = awssdk.ToString(r.DestinationPrefixListId)
		}
		if notEmptyStr(r.EgressOnlyInternetGatewayId) {
			routeTarget := &graph.RouteTarget{Type: graph.EgressOnlyInternetGatewayTarget, Ref: awssdk.ToString(r.EgressOnlyInternetGatewayId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmptyStr(r.GatewayId) {
			routeTarget := &graph.RouteTarget{Type: graph.GatewayTarget, Ref: awssdk.ToString(r.GatewayId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmptyStr(r.InstanceId) {
			routeTarget := &graph.RouteTarget{Type: graph.InstanceTarget, Ref: awssdk.ToString(r.InstanceId), Owner: awssdk.ToString(r.InstanceOwnerId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmptyStr(r.NatGatewayId) {
			routeTarget := &graph.RouteTarget{Type: graph.NatTarget, Ref: awssdk.ToString(r.NatGatewayId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmptyStr(r.NetworkInterfaceId) {
			routeTarget := &graph.RouteTarget{Type: graph.NetworkInterfaceTarget, Ref: awssdk.ToString(r.NetworkInterfaceId)}
			route.Targets = append(route.Targets, routeTarget)
		}
		if notEmptyStr(r.VpcPeeringConnectionId) {
			routeTarget := &graph.RouteTarget{Type: graph.VpcPeeringConnectionTarget, Ref: awssdk.ToString(r.VpcPeeringConnectionId)}
			route.Targets = append(route.Targets, routeTarget)
		}

		routes = append(routes, route)
	}
	return routes, nil
}

var extractHasATrueBoolInStructSliceFn = func(key string) transformFn {
	return func(i interface{}) (interface{}, error) {
		var res bool
		value := reflect.ValueOf(i)
		if value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("extract true bool: not a slice but a %T", i)
		}
		for i := 0; i < value.Len(); i++ {
			e, err := extractFieldFn(key)(value.Index(i).Interface())
			if err != nil {
				return res, err
			}
			if e == nil {
				continue //Empty field
			}
			b, ok := e.(bool)
			if !ok {
				return nil, fmt.Errorf("extract true bool: the field %s is not a boolean, but has type: %T", key, e)
			}
			if b {
				res = true
			}
		}

		return res, nil
	}
}

var extractDistributionOriginFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.(*cloudfronttypes.Origins); !ok {
		return nil, fmt.Errorf("extract origins: not a origins pointer but a %T", i)
	}
	var origins []*graph.DistributionOrigin
	for _, o := range i.(*cloudfronttypes.Origins).Items {
		origin := &graph.DistributionOrigin{
			ID:         awssdk.ToString(o.Id),
			PublicDNS:  awssdk.ToString(o.DomainName),
			PathPrefix: awssdk.ToString(o.OriginPath),
		}
		if o.S3OriginConfig != nil && awssdk.ToString(o.S3OriginConfig.OriginAccessIdentity) != "" {
			origin.OriginType = "s3"
			origin.Config = awssdk.ToString(o.S3OriginConfig.OriginAccessIdentity)
		}

		origins = append(origins, origin)
	}
	return origins, nil
}

var extractStackOutputsFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]cloudformationtypes.Output); !ok {
		return nil, fmt.Errorf("extract ouutputs not an output slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, out := range i.([]cloudformationtypes.Output) {
		keyval := &graph.KeyValue{KeyName: awssdk.ToString(out.OutputKey), Value: awssdk.ToString(out.OutputValue)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractStackParametersFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]cloudformationtypes.Parameter); !ok {
		return nil, fmt.Errorf("extract parameters not a parameter slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, out := range i.([]cloudformationtypes.Parameter) {
		keyval := &graph.KeyValue{KeyName: awssdk.ToString(out.ParameterKey), Value: awssdk.ToString(out.ParameterValue)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

var extractContainersImagesFn = func(i interface{}) (interface{}, error) {
	if _, ok := i.([]ecstypes.ContainerDefinition); !ok {
		return nil, fmt.Errorf("extract containers images, not a container definition slice but a %T", i)
	}
	var keyVals []*graph.KeyValue
	for _, out := range i.([]ecstypes.ContainerDefinition) {
		keyval := &graph.KeyValue{KeyName: awssdk.ToString(out.Name), Value: awssdk.ToString(out.Image)}

		keyVals = append(keyVals, keyval)
	}
	return keyVals, nil
}

func extractDocumentDefaultVersion(i interface{}) (interface{}, error) {
	if _, ok := i.([]iamtypes.PolicyVersion); !ok {
		return nil, fmt.Errorf("extract default version of document, not a policy version slice but a %T", i)
	}
	for _, version := range i.([]iamtypes.PolicyVersion) {
		if version.IsDefaultVersion {
			docStr := awssdk.ToString(version.Document)
			if str, err := url.QueryUnescape(docStr); err == nil {
				var buff bytes.Buffer
				err = json.Compact(&buff, []byte(str))
				return buff.String(), err
			}
			return docStr, nil
		}
	}
	return "", nil
}

func extractURLEncodedJson(i interface{}) (interface{}, error) {
	var docStr string
	switch v := i.(type) {
	case *string:
		docStr = awssdk.ToString(v)
	case string:
		docStr = v
	default:
		return nil, fmt.Errorf("extract URL-encoded JSON, not a string but a %T", i)
	}
	if str, err := url.QueryUnescape(docStr); err == nil {
		var buff bytes.Buffer
		err = json.Compact(&buff, []byte(str))
		return buff.String(), err
	}
	return docStr, nil
}

func notEmptyStr(str *string) bool {
	return str != nil && *str != ""
}

func isNilValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

func HashFields(fields ...interface{}) string {
	var buf bytes.Buffer
	for _, field := range fields {
		buf.WriteString(fmt.Sprint(field))
	}
	h := adler32.New()
	buf.WriteTo(h)
	return "awls-" + hex.EncodeToString(h.Sum(nil))
}

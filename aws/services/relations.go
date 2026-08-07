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

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	tstore "github.com/bootswithdefer/triplestore"

	awsconv "github.com/bootswithdefer/awless/aws/conv"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/graph"
)

const (
	PARENT_OF = iota // default
	APPLIES_ON
	DEPENDING_ON
)

type funcBuilder struct {
	parent                              string
	fieldName, listName, stringListName string
	relation                            int
}

type addParentFn func(*graph.Graph, tstore.RDFGraph, string, any) error

var addParentsFns = map[string][]addParentFn{
	// Infra
	cloud.Subnet: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
	},
	cloud.Instance: {
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId"}.build(),
		funcBuilder{parent: cloud.SecurityGroup, fieldName: "GroupId", listName: "SecurityGroups", relation: APPLIES_ON}.build(),
		funcBuilder{parent: cloud.Keypair, fieldName: "KeyName", relation: APPLIES_ON}.build(),
	},
	cloud.SecurityGroup: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
	},
	cloud.InternetGateway: {
		addRegionParent,
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId", listName: "Attachments", relation: DEPENDING_ON}.build(),
	},
	cloud.NatGateway: {
		addRegionParent,
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId", relation: DEPENDING_ON}.build(),
	},
	cloud.RouteTable: {
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId", listName: "Associations", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
	},
	cloud.Volume: {
		funcBuilder{parent: cloud.AvailabilityZone, fieldName: "AvailabilityZone"}.build(),
		funcBuilder{parent: cloud.Instance, fieldName: "InstanceId", listName: "Attachments", relation: DEPENDING_ON}.build(),
	},
	cloud.ElasticIP: {
		addRegionParent,
		funcBuilder{parent: cloud.Instance, fieldName: "InstanceId", relation: DEPENDING_ON}.build(),
	},
	cloud.Snapshot: {
		addRegionParent,
		funcBuilder{parent: cloud.Volume, fieldName: "VolumeId", relation: DEPENDING_ON}.build(),
	},
	cloud.NetworkInterface: {
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId", relation: PARENT_OF}.build(),
		funcBuilder{parent: cloud.SecurityGroup, fieldName: "GroupId", listName: "Groups", relation: APPLIES_ON}.build(),
		funcBuilder{parent: cloud.Instance, fieldName: "Attachment.InstanceId", relation: DEPENDING_ON}.build(),
	},
	// Loadbalancer
	cloud.LoadBalancer: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: cloud.Subnet, fieldName: "SubnetId", listName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.AvailabilityZone, fieldName: "ZoneName", listName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.SecurityGroup, stringListName: "SecurityGroups", relation: APPLIES_ON}.build(),
	},
	cloud.ClassicLoadBalancer: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VPCId"}.build(),
		funcBuilder{parent: cloud.Subnet, stringListName: "Subnets", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.AvailabilityZone, stringListName: "AvailabilityZones", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.SecurityGroup, stringListName: "SecurityGroups", relation: APPLIES_ON}.build(),
	},
	cloud.Listener: {
		funcBuilder{parent: cloud.LoadBalancer, fieldName: "LoadBalancerArn"}.build(),
	},
	cloud.TargetGroup: {
		funcBuilder{parent: cloud.Vpc, fieldName: "VpcId"}.build(),
		funcBuilder{parent: cloud.LoadBalancer, stringListName: "LoadBalancerArns", relation: APPLIES_ON}.build(),
		fetchTargetsAndAddRelations,
	},
	// Database
	cloud.Database: {
		funcBuilder{parent: cloud.AvailabilityZone, fieldName: "AvailabilityZone"}.build(),
		funcBuilder{parent: cloud.SecurityGroup, listName: "VpcSecurityGroups", fieldName: "VpcSecurityGroupId", relation: APPLIES_ON}.build(),
	},
	// Autoscaling
	cloud.LaunchConfiguration: {
		addRegionParent,
		funcBuilder{parent: cloud.Keypair, fieldName: "KeyName", relation: APPLIES_ON}.build(),
	},
	cloud.ScalingGroup: {
		addRegionParent,
		funcBuilder{parent: cloud.AvailabilityZone, stringListName: "AvailabilityZones", relation: APPLIES_ON}.build(),
		funcBuilder{parent: cloud.Instance, fieldName: "InstanceId", listName: "Instances", relation: DEPENDING_ON}.build(),
		funcBuilder{parent: cloud.TargetGroup, stringListName: "TargetGroupARNs", relation: DEPENDING_ON}.build(),
		addScalingGroupSubnets,
	},
	// Container
	cloud.ContainerInstance: {
		funcBuilder{parent: cloud.Instance, fieldName: "Ec2InstanceId", relation: APPLIES_ON}.build(),
	},
	cloud.Subscription: {
		funcBuilder{parent: cloud.Topic, fieldName: "TopicArn"}.build(),
	},
	cloud.Vpc:              {addRegionParent},
	cloud.AvailabilityZone: {addRegionParent},
	cloud.Keypair:          {addRegionParent},
	cloud.Image:            {addRegionParent},
	cloud.Repository:       {addRegionParent},
	cloud.ContainerCluster: {addRegionParent},
	cloud.ContainerTask:    {addRegionParent},
	cloud.Certificate:      {addRegionParent},
	cloud.User:             {userAddGroupsRelations, addManagedPoliciesRelations},
	cloud.Role:             {addManagedPoliciesRelations},
	cloud.Group:            {addManagedPoliciesRelations},
	cloud.Bucket:           {addRegionParent},
	cloud.Function:         {addRegionParent},
	cloud.Topic:            {addRegionParent},
	cloud.Alarm:            {addRegionParent, addAlarmMetric},
	cloud.Metric:           {addRegionParent},
	cloud.Stack:            {addRegionParent},
	cloud.MFADevice: {
		funcBuilder{parent: cloud.User, fieldName: "User.UserId", relation: DEPENDING_ON}.build(),
	},
}

func (fb funcBuilder) build() addParentFn {
	switch {
	case fb.listName != "":
		return fb.addRelationListWithField()
	case fb.stringListName != "":
		return fb.addRelationListWithStringField()
	default:
		return fb.addRelationWithField()
	}
}

func (fb funcBuilder) addRelationWithField() addParentFn {
	return func(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
		val, err := valueAtPath(i, fb.fieldName)
		if err != nil {
			return err
		}
		if val == nil {
			return nil
		}

		var strVal string
		switch s := val.(type) {
		case *string:
			strVal = aws.ToString(s)
		case string:
			strVal = s
		default:
			return fmt.Errorf("add parent to %s: %T not a string", fb.fieldName, val)
		}

		if strVal == "" {
			return nil
		}

		res, err := awsconv.InitResource(i)
		if err != nil {
			return err
		}

		parent := graph.InitResource(fb.parent, strVal)
		return addRelation(g, parent, res, fb.relation)
	}
}

func (fb funcBuilder) addRelationListWithStringField() addParentFn {
	return func(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
		structField, err := verifyValidStructField(i, fb.stringListName)
		if err != nil {
			return err
		}

		res, err := awsconv.InitResource(i)
		if err != nil {
			return err
		}

		if !structField.IsValid() || structField.Kind() != reflect.Slice {
			return fmt.Errorf("add parent to %s: field not a slice: %T", res.Id(), structField.Kind())
		}

		for i := 0; i < structField.Len(); i++ {
			elem := structField.Index(i).Interface()
			var strVal string
			switch s := elem.(type) {
			case *string:
				strVal = aws.ToString(s)
			case string:
				strVal = s
			default:
				return fmt.Errorf("add parent to %s: not a string: %T", res.Id(), elem)
			}

			if strVal == "" {
				continue
			}
			parent := graph.InitResource(fb.parent, strVal)

			if err = addRelation(g, parent, res, fb.relation); err != nil {
				return err
			}
		}
		return nil
	}
}

func (fb funcBuilder) addRelationListWithField() addParentFn {
	return func(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
		structField, err := verifyValidStructField(i, fb.listName)
		if err != nil {
			return err
		}

		res, err := awsconv.InitResource(i)
		if err != nil {
			return err
		}

		if !structField.IsValid() || structField.Kind() != reflect.Slice {
			return fmt.Errorf("add parent to %s: field not a slice: %T", res.Id(), structField.Kind())
		}

		for i := 0; i < structField.Len(); i++ {
			listValue := structField.Index(i)
			var listStruc reflect.Value
			switch listValue.Kind() {
			case reflect.Ptr:
				listStruc = listValue.Elem()
			case reflect.Struct:
				listStruc = listValue
			default:
				return fmt.Errorf("add parent to %s: not a struct or pointer: %s", res.Id(), listValue.Kind())
			}
			if listStruc.Kind() != reflect.Struct {
				return fmt.Errorf("add parent to %s: not a struct: %s", res.Id(), listStruc.Kind())
			}
			listStructField := listStruc.FieldByName(fb.fieldName)
			if !listStructField.IsValid() {
				return fmt.Errorf("add parent to %s: unknown field %s in %d", res.Id(), listStructField, i)
			}

			var strVal string
			switch s := listStructField.Interface().(type) {
			case *string:
				strVal = aws.ToString(s)
			case string:
				strVal = s
			default:
				return fmt.Errorf("add parent to %s: %T is not a string", listStructField, listStructField.Interface())
			}

			if strVal == "" {
				continue
			}
			parent := graph.InitResource(fb.parent, strVal)

			if err = addRelation(g, parent, res, fb.relation); err != nil {
				return err
			}
		}
		return nil
	}
}

func verifyValidStructField(i any, name string) (reflect.Value, error) {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("%T not a pointer", i)
	}
	struc := value.Elem()
	if struc.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%T not a stuct pointer", i)
	}

	structField := struc.FieldByName(name)
	if !structField.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid field %s: ", name)
	}

	return structField, nil
}

func addRelation(g *graph.Graph, first, other *graph.Resource, relation int) error {
	switch relation {
	case PARENT_OF:
		return g.AddParentRelation(first, other)
	case APPLIES_ON:
		return g.AddAppliesOnRelation(first, other)
	case DEPENDING_ON:
		return g.AddAppliesOnRelation(other, first)
	default:
		return errors.New("unknown relation type")
	}
}

func addRegionParent(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
	res, err := awsconv.InitResource(i)
	if err != nil {
		return err
	}
	return g.AddParentRelation(graph.InitResource(cloud.Region, region), res)
}

func addManagedPoliciesRelations(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
	res, err := awsconv.InitResource(i)
	if err != nil {
		return err
	}
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("add parent to %s: unknown type %T", res.Id(), i)
	}
	struc := value.Elem()
	if struc.Kind() != reflect.Struct {
		return fmt.Errorf("add parent to %s: unknown type %T", res.Id(), i)
	}

	structField := struc.FieldByName("AttachedManagedPolicies")
	if !structField.IsValid() {
		return fmt.Errorf("add parent to %s: unknown field %s in %d", res.Id(), structField, i)
	}
	policies, ok := structField.Interface().([]iamtypes.AttachedPolicy)
	if !ok {
		return fmt.Errorf("add parent to %s: not a valid attached policy list: %T", res.Id(), structField.Interface())
	}

	for _, policy := range policies {
		policies, err := graph.ResolveResourcesWithProp(snap, cloud.Policy, "Name", aws.ToString(policy.PolicyName))
		if err != nil {
			return err
		}
		if len(policies) != 1 {
			fmt.Fprintf(os.Stderr, "add parent to '%s/%s': unknown policy named '%s'. Ignoring it.\n", res.Type(), res.Id(), aws.ToString(policy.PolicyName))
			return nil
		}
		if err := g.AddAppliesOnRelation(policies[0], res); err != nil {
			return err
		}
	}
	return nil
}

func userAddGroupsRelations(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
	user, ok := i.(*iamtypes.UserDetail)
	if !ok {
		return fmt.Errorf("aws fetch: not a user, but a %T", i)
	}
	n, err := awsconv.InitResource(user)
	if err != nil {
		return err
	}

	for _, group := range user.GroupList {
		groupName := group
		resources, err := graph.ResolveResourcesWithProp(snap, cloud.Group, "Name", groupName)
		if err != nil {
			return err
		}
		switch len(resources) {
		case 0:
			fmt.Fprintf(os.Stderr, "no group with name %s found for user %s\n", groupName, n.Id())
		case 1:
			if err := g.AddAppliesOnRelation(resources[0], n); err != nil {
				return err
			}
		default:
			fmt.Fprintf(os.Stderr, "multiple groups with name %s found for user %s:%v\n", groupName, n.Id(), resources)
		}
	}
	return nil
}

func fetchTargetsAndAddRelations(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
	group, ok := i.(*elbv2types.TargetGroup)
	if !ok {
		return fmt.Errorf("add targets relation: not a target group, but a %T", i)
	}
	parent, err := awsconv.InitResource(group)
	if err != nil {
		return err
	}

	targets, err := InfraService.(*Infra).Elbv2Client.DescribeTargetHealth(context.Background(), &elbv2.DescribeTargetHealthInput{TargetGroupArn: group.TargetGroupArn})
	if err != nil {
		return err
	}

	for _, t := range targets.TargetHealthDescriptions {
		n := graph.InitResource(cloud.Instance, aws.ToString(t.Target.Id))
		err = g.AddAppliesOnRelation(parent, n)
		if err != nil {
			return err
		}
	}
	return nil
}

func addScalingGroupSubnets(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
	group, ok := i.(*autoscalingtypes.AutoScalingGroup)
	if !ok {
		return fmt.Errorf("add autoscaling group relation: not a autoscaling group, but a %T", i)
	}
	parent, err := awsconv.InitResource(group)
	if err != nil {
		return err
	}
	if subnets := aws.ToString(group.VPCZoneIdentifier); subnets != "" {
		splits := strings.Split(subnets, ",")
		for _, split := range splits {
			n := graph.InitResource(cloud.Subnet, split)
			err = g.AddAppliesOnRelation(parent, n)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// valueAtPath navigates a dot-separated field path through struct fields,
// returning the value found. Replaces awsutil.ValuesAtPath from SDK v1.
func valueAtPath(i any, path string) (any, error) {
	parts := strings.Split(path, ".")
	v := reflect.ValueOf(i)
	for _, part := range parts {
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return nil, nil
			}
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected struct at path segment '%s', got %s", part, v.Kind())
		}
		v = v.FieldByName(part)
		if !v.IsValid() {
			return nil, fmt.Errorf("field '%s' not found", part)
		}
	}
	return v.Interface(), nil
}

func addAlarmMetric(g *graph.Graph, snap tstore.RDFGraph, region string, i any) error {
	alarm, ok := i.(*cloudwatchtypes.MetricAlarm)
	if !ok {
		return fmt.Errorf("add alarm metric relation: not a alarm, but a %T", i)
	}
	parent, err := awsconv.InitResource(alarm)
	if err != nil {
		return err
	}
	if namespace, metric := aws.ToString(alarm.Namespace), aws.ToString(alarm.MetricName); namespace != "" && metric != "" {
		id := awsconv.HashFields(namespace, metric)
		n := graph.InitResource(cloud.Metric, id)
		err = g.AddAppliesOnRelation(parent, n)
		if err != nil {
			return err
		}
	}
	return nil
}

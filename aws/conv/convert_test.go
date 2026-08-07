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
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	"github.com/bootswithdefer/awless/graph"
)

func TestTransformFunctions(t *testing.T) {
	t.Parallel()
	t.Run("extractTag", func(t *testing.T) {
		t.Parallel()
		tag := []ec2types.Tag{
			{Key: awssdk.String("Name"), Value: awssdk.String("instance-name")},
			{Key: awssdk.String("Created with"), Value: awssdk.String("awless")},
		}

		val, _ := extractTagFn("Name")(tag)
		if got, want := fmt.Sprint(val), "instance-name"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		val, _ = extractTagFn("Created with")(tag)
		if got, want := fmt.Sprint(val), "awless"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("extractValue", func(t *testing.T) {
		t.Parallel()
		val, _ := extractValueFn(awssdk.String("any"))
		if got, want := val.(string), "any"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		val, _ = extractValueFn(awssdk.Int(2))
		if got, want := val.(int), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		val, _ = extractValueFn(awssdk.Int64(4))
		if got, want := val.(int64), int64(4); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}

		val, _ = extractValueFn(awssdk.Bool(true))
		if got, want := val.(bool), true; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}
	})

	t.Run("extractField", func(t *testing.T) {
		t.Parallel()
		data := &ec2types.InstanceState{Code: awssdk.Int32(12), Name: ec2types.InstanceStateNameRunning}

		val, _ := extractFieldFn("Code")(data)
		if got, want := val.(int32), int32(12); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		val, _ = extractFieldFn("Name")(data)
		if got, want := fmt.Sprint(val), "running"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("extractClassicLoadbListenerDescriptions", func(t *testing.T) {
		t.Parallel()
		descriptions := []elbtypes.ListenerDescription{
			{Listener: &elbtypes.Listener{LoadBalancerPort: 443, Protocol: awssdk.String("HTTPS"), InstancePort: awssdk.Int32(8080), InstanceProtocol: awssdk.String("HTTP")}},
			{Listener: &elbtypes.Listener{LoadBalancerPort: 5000, Protocol: awssdk.String("SSL"), InstancePort: awssdk.Int32(3000), InstanceProtocol: awssdk.String("TCP")}},
		}

		expected := []string{"HTTPS:443:HTTP:8080", "SSL:5000:TCP:3000"}

		i, err := extractClassicLoadbListenerDescriptionsFn(descriptions)
		if err != nil {
			t.Fatal(err)
		}
		res := i.([]string)
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i], expected[i]; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}

	})

	t.Run("extractIpPermissions", func(t *testing.T) {
		t.Parallel()
		ipPermissions := []ec2types.IpPermission{
			{FromPort: awssdk.Int32(70),
				ToPort:     awssdk.Int32(85),
				IpProtocol: awssdk.String("udp"),
				IpRanges:   []ec2types.IpRange{},
			},
			{FromPort: awssdk.Int32(-1),
				ToPort:     awssdk.Int32(-1),
				IpProtocol: awssdk.String("-1"),
				IpRanges:   []ec2types.IpRange{{CidrIp: awssdk.String("10.192.24.0/24")}},
			},
			{FromPort: awssdk.Int32(12),
				ToPort:     awssdk.Int32(12),
				IpProtocol: awssdk.String("27"),
				IpRanges: []ec2types.IpRange{
					{CidrIp: awssdk.String("1.2.3.4/32")},
					{CidrIp: awssdk.String("2.3.0.0/16")},
				},
			},
			{FromPort: awssdk.Int32(22),
				ToPort:     awssdk.Int32(22),
				IpProtocol: awssdk.String("tcp"),
				IpRanges: []ec2types.IpRange{
					{CidrIp: awssdk.String("1.2.3.4/32")},
				},
				Ipv6Ranges: []ec2types.Ipv6Range{
					{CidrIpv6: awssdk.String("fd34:fe56:7891:2f3a::/64")},
					{CidrIpv6: awssdk.String("2001:db8::/110")},
				},
			},
			{FromPort: awssdk.Int32(0),
				ToPort:     awssdk.Int32(65535),
				IpProtocol: awssdk.String("tcp"),
				UserIdGroupPairs: []ec2types.UserIdGroupPair{
					{GroupId: awssdk.String("group_1")},
					{GroupId: awssdk.String("group_2")},
				},
			},
		}

		expected := []*graph.FirewallRule{
			{
				PortRange: graph.PortRange{FromPort: int64(70), ToPort: int64(85), Any: false},
				Protocol:  "udp",
				IPRanges:  []*net.IPNet{},
			},
			{
				PortRange: graph.PortRange{Any: true},
				Protocol:  "any",
				IPRanges:  []*net.IPNet{{IP: net.IPv4(10, 192, 24, 0), Mask: net.CIDRMask(24, 32)}},
			},
			{
				PortRange: graph.PortRange{Any: true},
				Protocol:  "27",
				IPRanges: []*net.IPNet{
					{IP: net.IPv4(1, 2, 3, 4), Mask: net.CIDRMask(32, 32)},
					{IP: net.IPv4(2, 3, 0, 0), Mask: net.CIDRMask(16, 32)},
				},
			},
			{
				PortRange: graph.PortRange{FromPort: int64(22), ToPort: int64(22), Any: false},
				Protocol:  "tcp",
				IPRanges: []*net.IPNet{
					{IP: net.IPv4(1, 2, 3, 4), Mask: net.CIDRMask(32, 32)},
					{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
					{IP: net.ParseIP("2001:db8::"), Mask: net.CIDRMask(110, 128)},
				},
			},
			{
				PortRange: graph.PortRange{FromPort: int64(0), ToPort: int64(65535), Any: false},
				Protocol:  "tcp",
				Sources: []string{
					"group_1", "group_2",
				},
			},
		}

		i, err := extractIpPermissionSliceFn(ipPermissions)
		if err != nil {
			t.Fatal(err)
		}
		res := i.([]*graph.FirewallRule)
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i].String(), expected[i].String(); got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}
	})

	t.Run("extractSliceValues", func(t *testing.T) {
		t.Parallel()
		slice := []ec2types.GroupIdentifier{
			{GroupId: awssdk.String("MyGroup1"), GroupName: awssdk.String("MyGroupName1")},
			{GroupId: awssdk.String("MyGroup2"), GroupName: awssdk.String("MyGroupName2")},
		}

		val, err := extractStringSliceValues("GroupId")(slice)
		if err != nil {
			t.Fatal(err)
		}
		expectedI := []string{"MyGroup1", "MyGroup2"}
		if got, want := val, expectedI; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		val, err = extractStringSliceValues("GroupName")(slice)
		if err != nil {
			t.Fatal(err)
		}
		expectedI = []string{"MyGroupName1", "MyGroupName2"}
		if got, want := val, expectedI; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("extractRoutesSlice", func(t *testing.T) {
		t.Parallel()
		routes := []ec2types.Route{
			{
				DestinationCidrBlock:        awssdk.String("10.0.0.0/24"),
				EgressOnlyInternetGatewayId: awssdk.String("test_id_1"),
			},
			{
				DestinationIpv6CidrBlock: awssdk.String("fd34:fe56:7891:2f3a::/64"),
				GatewayId:                awssdk.String("test_id_2"),
			},
			{
				DestinationCidrBlock: awssdk.String("0.0.0.0/0"),
				InstanceId:           awssdk.String("test_id_3"),
			},
			{
				DestinationCidrBlock: awssdk.String("0.0.0.0/0"),
				NatGatewayId:         awssdk.String("test_id_4"),
			},
			{
				DestinationCidrBlock: awssdk.String("0.0.0.0/0"),
				NetworkInterfaceId:   awssdk.String("test_id_5"),
			},
			{
				DestinationCidrBlock:   awssdk.String("0.0.0.0/0"),
				VpcPeeringConnectionId: awssdk.String("test_id_6"),
			},
			{
				DestinationCidrBlock:     awssdk.String("10.0.0.0/24"),
				DestinationIpv6CidrBlock: awssdk.String("fd34:fe56:7891:2f3a::/64"),
				VpcPeeringConnectionId:   awssdk.String("test_id_7"),
			},
			{
				DestinationPrefixListId: awssdk.String("pl-0123456"),
				GatewayId:               awssdk.String("test_id_8"),
			},
			{
				DestinationCidrBlock:    awssdk.String("0.0.0.0/0"),
				DestinationPrefixListId: awssdk.String("pl-0123456"),
				InstanceId:              awssdk.String("test_id_9"),
				InstanceOwnerId:         awssdk.String("owner"),
				NetworkInterfaceId:      awssdk.String("eni-123456"),
			},
		}

		expected := []*graph.Route{
			{
				Destination: &net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(24, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.EgressOnlyInternetGatewayTarget, Ref: "test_id_1"},
				},
			},
			{
				DestinationIPv6: &net.IPNet{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
				Targets: []*graph.RouteTarget{
					{Type: graph.GatewayTarget, Ref: "test_id_2"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.InstanceTarget, Ref: "test_id_3"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.NatTarget, Ref: "test_id_4"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.NetworkInterfaceTarget, Ref: "test_id_5"},
				},
			},
			{
				Destination: &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				Targets: []*graph.RouteTarget{
					{Type: graph.VpcPeeringConnectionTarget, Ref: "test_id_6"},
				},
			},
			{
				Destination:     &net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(24, 32)},
				DestinationIPv6: &net.IPNet{IP: net.ParseIP("fd34:fe56:7891:2f3a::"), Mask: net.CIDRMask(64, 128)},
				Targets: []*graph.RouteTarget{
					{Type: graph.VpcPeeringConnectionTarget, Ref: "test_id_7"},
				},
			},
			{
				DestinationPrefixListId: "pl-0123456",
				Targets: []*graph.RouteTarget{
					{Type: graph.GatewayTarget, Ref: "test_id_8"},
				},
			},
			{
				Destination:             &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)},
				DestinationPrefixListId: "pl-0123456",
				Targets: []*graph.RouteTarget{
					{Type: graph.InstanceTarget, Ref: "test_id_9", Owner: "owner"},
					{Type: graph.NetworkInterfaceTarget, Ref: "eni-123456"},
				},
			},
		}

		i, err := extractRoutesSliceFn(routes)
		if err != nil {
			t.Fatal(err)
		}
		res := i.([]*graph.Route)
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i].String(), expected[i].String(); got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}
	})

	t.Run("extractHasATrueBoolInStructSlice", func(t *testing.T) {
		t.Parallel()
		slice := []ec2types.RouteTableAssociation{
			{Main: awssdk.Bool(false), RouteTableAssociationId: awssdk.String("test")},
			{RouteTableId: awssdk.String("test2"), Main: awssdk.Bool(true)},
		}

		val, err := extractHasATrueBoolInStructSliceFn("Main")(slice)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(bool), true; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}

		slice = []ec2types.RouteTableAssociation{
			{Main: awssdk.Bool(false), RouteTableAssociationId: awssdk.String("test")},
			{RouteTableId: awssdk.String("test2"), Main: awssdk.Bool(false)},
		}

		val, err = extractHasATrueBoolInStructSliceFn("Main")(slice)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(bool), false; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}
	})

	t.Run("extractHasATrueBoolInStructSlice not a slice", func(t *testing.T) {
		t.Parallel()
		_, err := extractHasATrueBoolInStructSliceFn("Main")("not-a-slice")
		if err == nil {
			t.Fatal("expected error for non-slice input")
		}
	})

	t.Run("extractTime with time pointer", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2023, 6, 15, 10, 30, 0, 0, time.FixedZone("EST", -5*3600))
		val, err := extractTimeFn(&now)
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Location() != time.UTC {
			t.Errorf("expected UTC, got %v", result.Location())
		}
		if got, want := result, now.UTC(); !got.Equal(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("extractTime with string", func(t *testing.T) {
		t.Parallel()
		val, err := extractTimeFn("2023-06-15T10:30:00.000+0000")
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Year() != 2023 || result.Month() != 6 || result.Day() != 15 {
			t.Errorf("unexpected time: %v", result)
		}
	})

	t.Run("extractTime with *string", func(t *testing.T) {
		t.Parallel()
		s := "2023-06-15T10:30:00.000+0000"
		val, err := extractTimeFn(&s)
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Year() != 2023 {
			t.Errorf("unexpected year: %d", result.Year())
		}
	})

	t.Run("extractTime with bad string", func(t *testing.T) {
		t.Parallel()
		_, err := extractTimeFn("not-a-date")
		if err == nil {
			t.Fatal("expected error for bad date string")
		}
	})

	t.Run("extractTime with bad *string", func(t *testing.T) {
		t.Parallel()
		s := "not-a-date"
		_, err := extractTimeFn(&s)
		if err == nil {
			t.Fatal("expected error for bad *string date")
		}
	})

	t.Run("extractTime with unsupported type", func(t *testing.T) {
		t.Parallel()
		_, err := extractTimeFn(42)
		if err == nil {
			t.Fatal("expected error for unsupported type")
		}
	})

	t.Run("extractMillisecondEpochTime with *int64", func(t *testing.T) {
		t.Parallel()
		ms := int64(1686825000000) // 2023-06-15T13:30:00Z
		val, err := extractMillisecondEpochTimeFn(&ms)
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Year() != 2023 {
			t.Errorf("unexpected year: %d", result.Year())
		}
		if result.Location() != time.UTC {
			t.Errorf("expected UTC, got %v", result.Location())
		}
	})

	t.Run("extractMillisecondEpochTime with nil *int64", func(t *testing.T) {
		t.Parallel()
		var p *int64
		val, err := extractMillisecondEpochTimeFn(p)
		if err != nil {
			t.Fatal(err)
		}
		if val != nil {
			t.Errorf("expected nil for nil pointer, got %v", val)
		}
	})

	t.Run("extractMillisecondEpochTime with int64", func(t *testing.T) {
		t.Parallel()
		val, err := extractMillisecondEpochTimeFn(int64(1686825000000))
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Year() != 2023 {
			t.Errorf("unexpected year: %d", result.Year())
		}
	})

	t.Run("extractMillisecondEpochTime with unsupported type", func(t *testing.T) {
		t.Parallel()
		_, err := extractMillisecondEpochTimeFn("not-an-int")
		if err == nil {
			t.Fatal("expected error for unsupported type")
		}
	})

	t.Run("extractTimeWithZSuffix with time pointer", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC)
		val, err := extractTimeWithZSuffixFn(&now)
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if !result.Equal(now) {
			t.Errorf("got %v, want %v", result, now)
		}
	})

	t.Run("extractTimeWithZSuffix with string", func(t *testing.T) {
		t.Parallel()
		val, err := extractTimeWithZSuffixFn("2023-06-15T10:30:00.000Z")
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Year() != 2023 || result.Month() != 6 || result.Day() != 15 {
			t.Errorf("unexpected time: %v", result)
		}
	})

	t.Run("extractTimeWithZSuffix with *string", func(t *testing.T) {
		t.Parallel()
		s := "2023-06-15T10:30:00.000Z"
		val, err := extractTimeWithZSuffixFn(&s)
		if err != nil {
			t.Fatal(err)
		}
		result := val.(time.Time)
		if result.Year() != 2023 {
			t.Errorf("unexpected year: %d", result.Year())
		}
	})

	t.Run("extractTimeWithZSuffix with bad string", func(t *testing.T) {
		t.Parallel()
		_, err := extractTimeWithZSuffixFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractTimeWithZSuffix with bad *string", func(t *testing.T) {
		t.Parallel()
		s := "bad"
		_, err := extractTimeWithZSuffixFn(&s)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractTimeWithZSuffix with unsupported type", func(t *testing.T) {
		t.Parallel()
		_, err := extractTimeWithZSuffixFn(42)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractNameValue", func(t *testing.T) {
		t.Parallel()
		dims := []cloudwatchtypes.Dimension{
			{Name: awssdk.String("InstanceId"), Value: awssdk.String("i-12345")},
			{Name: awssdk.String("AutoScalingGroupName"), Value: awssdk.String("my-asg")},
		}
		val, err := extractNameValueFn(dims)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.KeyValue)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d entries, want %d", got, want)
		}
		if got, want := result[0].KeyName, "InstanceId"; got != want {
			t.Errorf("got key %s, want %s", got, want)
		}
		if got, want := result[0].Value, "i-12345"; got != want {
			t.Errorf("got value %s, want %s", got, want)
		}
	})

	t.Run("extractNameValue wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractNameValueFn("not-dimensions")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractECSAttributes", func(t *testing.T) {
		t.Parallel()
		attrs := []ecstypes.Attribute{
			{Name: awssdk.String("ecs.instance-type"), Value: awssdk.String("t2.micro")},
			{Name: awssdk.String("ecs.os-type"), Value: awssdk.String("linux")},
		}
		val, err := extractECSAttributesFn(attrs)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.KeyValue)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0].KeyName, "ecs.instance-type"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractECSAttributes wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractECSAttributesFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractRouteTableAssociations", func(t *testing.T) {
		t.Parallel()
		assocs := []ec2types.RouteTableAssociation{
			{RouteTableAssociationId: awssdk.String("rtbassoc-001"), SubnetId: awssdk.String("subnet-abc")},
			{RouteTableAssociationId: awssdk.String("rtbassoc-002"), SubnetId: awssdk.String("subnet-def")},
		}
		val, err := extractRouteTableAssociationsFn(assocs)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.KeyValue)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0].KeyName, "rtbassoc-001"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[0].Value, "subnet-abc"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractRouteTableAssociations wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractRouteTableAssociationsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractTags ec2", func(t *testing.T) {
		t.Parallel()
		tags := []ec2types.Tag{
			{Key: awssdk.String("Name"), Value: awssdk.String("my-instance")},
			{Key: awssdk.String("env"), Value: awssdk.String("prod")},
		}
		val, err := extractTagsFn(tags)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0], "Name=my-instance"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractTags autoscaling", func(t *testing.T) {
		t.Parallel()
		tags := []autoscalingtypes.TagDescription{
			{Key: awssdk.String("Name"), Value: awssdk.String("my-asg")},
		}
		val, err := extractTagsFn(tags)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if got, want := len(result), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0], "Name=my-asg"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractTags wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractTagsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractMapTags", func(t *testing.T) {
		t.Parallel()
		tags := map[string]string{"Name": "my-cluster", "env": "dev"}
		val, err := extractMapTagsFn(tags)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	})

	t.Run("extractMapTags wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractMapTagsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractSecretsManagerTags", func(t *testing.T) {
		t.Parallel()
		tags := []secretsmanagertypes.Tag{
			{Key: awssdk.String("env"), Value: awssdk.String("staging")},
		}
		val, err := extractSecretsManagerTagsFn(tags)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if got, want := len(result), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0], "env=staging"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractSecretsManagerTags wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractSecretsManagerTagsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractEFSTags", func(t *testing.T) {
		t.Parallel()
		tags := []efstypes.Tag{
			{Key: awssdk.String("Name"), Value: awssdk.String("my-fs")},
			{Key: awssdk.String("env"), Value: awssdk.String("prod")},
		}
		val, err := extractEFSTagsFn(tags)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0], "Name=my-fs"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractEFSTags wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractEFSTagsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractEFSTag found", func(t *testing.T) {
		t.Parallel()
		tags := []efstypes.Tag{
			{Key: awssdk.String("Name"), Value: awssdk.String("my-fs")},
			{Key: awssdk.String("env"), Value: awssdk.String("prod")},
		}
		val, err := extractEFSTagFn("Name")(tags)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), "my-fs"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractEFSTag not found", func(t *testing.T) {
		t.Parallel()
		tags := []efstypes.Tag{
			{Key: awssdk.String("env"), Value: awssdk.String("prod")},
		}
		_, err := extractEFSTagFn("Name")(tags)
		if !errors.Is(err, ErrTagNotFound) {
			t.Errorf("expected ErrTagNotFound, got %v", err)
		}
	})

	t.Run("extractEFSTag wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractEFSTagFn("Name")("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractDynamoDBKeySchema", func(t *testing.T) {
		t.Parallel()
		keys := []dynamodbtypes.KeySchemaElement{
			{AttributeName: awssdk.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
			{AttributeName: awssdk.String("sk"), KeyType: dynamodbtypes.KeyTypeRange},
		}
		val, err := extractDynamoDBKeySchemaFn(keys)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0], "pk:HASH"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[1], "sk:RANGE"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractDynamoDBKeySchema wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractDynamoDBKeySchemaFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractEKSScalingConfig", func(t *testing.T) {
		t.Parallel()
		sc := &ekstypes.NodegroupScalingConfig{
			MinSize:     awssdk.Int32(1),
			MaxSize:     awssdk.Int32(10),
			DesiredSize: awssdk.Int32(3),
		}
		val, err := extractEKSScalingConfigFn(sc)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), "min:1,max:10,desired:3"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractEKSScalingConfig nil", func(t *testing.T) {
		t.Parallel()
		var sc *ekstypes.NodegroupScalingConfig
		val, err := extractEKSScalingConfigFn(sc)
		if err != nil {
			t.Fatal(err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("extractEKSScalingConfig wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractEKSScalingConfigFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractTag not found returns ErrTagNotFound", func(t *testing.T) {
		t.Parallel()
		tags := []ec2types.Tag{
			{Key: awssdk.String("env"), Value: awssdk.String("prod")},
		}
		_, err := extractTagFn("Name")(tags)
		if !errors.Is(err, ErrTagNotFound) {
			t.Errorf("expected ErrTagNotFound, got %v", err)
		}
	})

	t.Run("extractTag wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractTagFn("Name")("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractStringPointerSliceValues with []string", func(t *testing.T) {
		t.Parallel()
		val, err := extractStringPointerSliceValues([]string{"a", "b"})
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if !reflect.DeepEqual(result, []string{"a", "b"}) {
			t.Errorf("got %v, want [a b]", result)
		}
	})

	t.Run("extractStringPointerSliceValues with []*string", func(t *testing.T) {
		t.Parallel()
		val, err := extractStringPointerSliceValues([]*string{awssdk.String("x"), awssdk.String("y")})
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if !reflect.DeepEqual(result, []string{"x", "y"}) {
			t.Errorf("got %v, want [x y]", result)
		}
	})

	t.Run("extractStringPointerSliceValues wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractStringPointerSliceValues(42)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractDistributionOrigin", func(t *testing.T) {
		t.Parallel()
		origins := &cloudfronttypes.Origins{
			Items: []cloudfronttypes.Origin{
				{
					Id:         awssdk.String("origin-1"),
					DomainName: awssdk.String("mybucket.s3.amazonaws.com"),
					OriginPath: awssdk.String("/assets"),
					S3OriginConfig: &cloudfronttypes.S3OriginConfig{
						OriginAccessIdentity: awssdk.String("origin-access-identity/cloudfront/E12345"),
					},
				},
				{
					Id:         awssdk.String("origin-2"),
					DomainName: awssdk.String("example.com"),
					OriginPath: awssdk.String(""),
				},
			},
		}
		val, err := extractDistributionOriginFn(origins)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.DistributionOrigin)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0].ID, "origin-1"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[0].OriginType, "s3"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[0].Config, "origin-access-identity/cloudfront/E12345"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[1].ID, "origin-2"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractDistributionOrigin wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractDistributionOriginFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractStackOutputs", func(t *testing.T) {
		t.Parallel()
		outputs := []cloudformationtypes.Output{
			{OutputKey: awssdk.String("VpcId"), OutputValue: awssdk.String("vpc-123")},
			{OutputKey: awssdk.String("SubnetId"), OutputValue: awssdk.String("subnet-456")},
		}
		val, err := extractStackOutputsFn(outputs)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.KeyValue)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0].KeyName, "VpcId"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[0].Value, "vpc-123"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractStackOutputs wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractStackOutputsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractStackParameters", func(t *testing.T) {
		t.Parallel()
		params := []cloudformationtypes.Parameter{
			{ParameterKey: awssdk.String("KeyName"), ParameterValue: awssdk.String("my-key")},
		}
		val, err := extractStackParametersFn(params)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.KeyValue)
		if got, want := len(result), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0].KeyName, "KeyName"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractStackParameters wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractStackParametersFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractContainersImages", func(t *testing.T) {
		t.Parallel()
		containers := []ecstypes.ContainerDefinition{
			{Name: awssdk.String("web"), Image: awssdk.String("nginx:latest")},
			{Name: awssdk.String("app"), Image: awssdk.String("myapp:v1")},
		}
		val, err := extractContainersImagesFn(containers)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]*graph.KeyValue)
		if got, want := len(result), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		if got, want := result[0].KeyName, "web"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		if got, want := result[0].Value, "nginx:latest"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractContainersImages wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractContainersImagesFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractDocumentDefaultVersion", func(t *testing.T) {
		t.Parallel()
		doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","Resource":"*"}]}`
		versions := []iamtypes.PolicyVersion{
			{Document: awssdk.String(doc), IsDefaultVersion: false, VersionId: awssdk.String("v1")},
			{Document: awssdk.String(doc), IsDefaultVersion: true, VersionId: awssdk.String("v2")},
		}
		val, err := extractDocumentDefaultVersion(versions)
		if err != nil {
			t.Fatal(err)
		}
		if val.(string) == "" {
			t.Error("expected non-empty document")
		}
	})

	t.Run("extractDocumentDefaultVersion no default", func(t *testing.T) {
		t.Parallel()
		versions := []iamtypes.PolicyVersion{
			{Document: awssdk.String("{}"), IsDefaultVersion: false},
		}
		val, err := extractDocumentDefaultVersion(versions)
		if err != nil {
			t.Fatal(err)
		}
		if val.(string) != "" {
			t.Errorf("expected empty string for no default, got %v", val)
		}
	})

	t.Run("extractDocumentDefaultVersion wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractDocumentDefaultVersion("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractURLEncodedJson with string", func(t *testing.T) {
		t.Parallel()
		encoded := "%7B%22key%22%3A%22value%22%7D"
		val, err := extractURLEncodedJson(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), `{"key":"value"}`; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractURLEncodedJson with *string", func(t *testing.T) {
		t.Parallel()
		encoded := "%7B%22key%22%3A%22value%22%7D"
		val, err := extractURLEncodedJson(&encoded)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), `{"key":"value"}`; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractURLEncodedJson wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractURLEncodedJson(42)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractURLEncodedJson plain json", func(t *testing.T) {
		t.Parallel()
		plain := `{"key":"value"}`
		val, err := extractURLEncodedJson(plain)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), `{"key":"value"}`; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractValueAsString", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueAsStringFn(awssdk.Int32(42))
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), "42"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractValue with nil pointer", func(t *testing.T) {
		t.Parallel()
		var p *string
		val, err := extractValueFn(p)
		if err != nil {
			t.Fatal(err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("extractValue with string", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn("direct-string")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), "direct-string"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractValue with []string", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn([]string{"a", "b"})
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if !reflect.DeepEqual(result, []string{"a", "b"}) {
			t.Errorf("got %v, want [a b]", result)
		}
	})

	t.Run("extractValue with []*string", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn([]*string{awssdk.String("x")})
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if !reflect.DeepEqual(result, []string{"x"}) {
			t.Errorf("got %v, want [x]", result)
		}
	})

	t.Run("extractValue with int32", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn(int32(99))
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(int32), int32(99); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("extractValue with int64", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn(int64(100))
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(int64), int64(100); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("extractValue with bool", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn(true)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(bool), true; got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("extractValue with typed string enum", func(t *testing.T) {
		t.Parallel()
		val, err := extractValueFn(ec2types.InstanceStateNameRunning)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(string), "running"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractValue with unsupported type", func(t *testing.T) {
		t.Parallel()
		_, err := extractValueFn(struct{}{})
		if err == nil {
			t.Fatal("expected error for unsupported type")
		}
	})

	t.Run("extractField with struct (not pointer)", func(t *testing.T) {
		t.Parallel()
		data := ec2types.InstanceState{Code: awssdk.Int32(16), Name: ec2types.InstanceStateNameRunning}
		val, err := extractFieldFn("Code")(data)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := val.(int32), int32(16); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("extractField non-struct pointer", func(t *testing.T) {
		t.Parallel()
		s := "hello"
		_, err := extractFieldFn("Foo")(&s)
		if err == nil {
			t.Fatal("expected error for non-struct pointer")
		}
	})

	t.Run("extractField non-struct non-pointer", func(t *testing.T) {
		t.Parallel()
		_, err := extractFieldFn("Foo")(42)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractField missing field", func(t *testing.T) {
		t.Parallel()
		data := ec2types.InstanceState{Code: awssdk.Int32(16)}
		_, err := extractFieldFn("NonExistentField")(data)
		if err == nil {
			t.Fatal("expected error for missing field")
		}
	})

	t.Run("extractStringSliceValues not a slice", func(t *testing.T) {
		t.Parallel()
		_, err := extractStringSliceValues("GroupId")("not-a-slice")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractIpPermissions wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractIpPermissionSliceFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractIpPermissions with icmp and -1 ports", func(t *testing.T) {
		t.Parallel()
		perms := []ec2types.IpPermission{
			{
				FromPort:   awssdk.Int32(-1),
				ToPort:     awssdk.Int32(-1),
				IpProtocol: awssdk.String("icmp"),
			},
		}
		val, err := extractIpPermissionSliceFn(perms)
		if err != nil {
			t.Fatal(err)
		}
		rules := val.([]*graph.FirewallRule)
		if got, want := rules[0].PortRange.Any, true; got != want {
			t.Errorf("expected Any=true for icmp with -1 ports")
		}
		if got, want := rules[0].Protocol, "icmp"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("extractRoutesSlice wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractRoutesSliceFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractClassicLoadbListenerDescriptions wrong type", func(t *testing.T) {
		t.Parallel()
		_, err := extractClassicLoadbListenerDescriptionsFn("bad")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extractClassicLoadbListenerDescriptions nil listener", func(t *testing.T) {
		t.Parallel()
		descriptions := []elbtypes.ListenerDescription{
			{Listener: nil},
		}
		val, err := extractClassicLoadbListenerDescriptionsFn(descriptions)
		if err != nil {
			t.Fatal(err)
		}
		result := val.([]string)
		if len(result) != 0 {
			t.Errorf("expected empty result for nil listener, got %v", result)
		}
	})
}

func TestNotEmptyStr(t *testing.T) {
	t.Parallel()

	if notEmptyStr(nil) {
		t.Error("nil should return false")
	}
	empty := ""
	if notEmptyStr(&empty) {
		t.Error("empty string should return false")
	}
	nonEmpty := "hello"
	if !notEmptyStr(&nonEmpty) {
		t.Error("non-empty string should return true")
	}
}

func TestIsNilValue(t *testing.T) {
	t.Parallel()

	// Pointer
	var p *string
	if !isNilValue(reflect.ValueOf(&p).Elem()) {
		t.Error("nil pointer should return true")
	}
	s := "hello"
	if isNilValue(reflect.ValueOf(&s)) {
		t.Error("non-nil pointer should return false")
	}

	// Slice
	var sl []string
	if !isNilValue(reflect.ValueOf(&sl).Elem()) {
		t.Error("nil slice should return true")
	}
	sl = []string{"a"}
	if isNilValue(reflect.ValueOf(sl)) {
		t.Error("non-nil slice should return false")
	}

	// Map
	var m map[string]string
	if !isNilValue(reflect.ValueOf(&m).Elem()) {
		t.Error("nil map should return true")
	}

	// Non-nillable type (int)
	if isNilValue(reflect.ValueOf(42)) {
		t.Error("int should return false")
	}

	// Struct
	if isNilValue(reflect.ValueOf(struct{}{})) {
		t.Error("struct should return false")
	}
}

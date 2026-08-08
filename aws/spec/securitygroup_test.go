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
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestBuildIpPermissionsFromParams(t *testing.T) {
	tcases := []struct {
		params   map[string]any
		expected []ec2types.IpPermission
	}{
		{
			params: map[string]any{
				"protocol":  "tcp",
				"cidr":      "192.168.1.10/24",
				"portrange": 80,
			},
			expected: []ec2types.IpPermission{
				{
					IpProtocol: awssdk.String("tcp"),
					IpRanges:   []ec2types.IpRange{{CidrIp: awssdk.String("192.168.1.10/24")}},
					FromPort:   awssdk.Int32(int32(80)),
					ToPort:     awssdk.Int32(int32(80)),
				},
			},
		},
		{
			params: map[string]any{
				"protocol": "any",
				"cidr":     "192.168.1.18/32",
			},
			expected: []ec2types.IpPermission{
				{
					IpProtocol: awssdk.String("-1"),
					IpRanges:   []ec2types.IpRange{{CidrIp: awssdk.String("192.168.1.18/32")}},
					FromPort:   awssdk.Int32(int32(-1)),
					ToPort:     awssdk.Int32(int32(-1)),
				},
			},
		},
		{
			params: map[string]any{
				"protocol":  "udp",
				"cidr":      "0.0.0.0/0",
				"portrange": "22-23",
			},
			expected: []ec2types.IpPermission{
				{
					IpProtocol: awssdk.String("udp"),
					IpRanges:   []ec2types.IpRange{{CidrIp: awssdk.String("0.0.0.0/0")}},
					FromPort:   awssdk.Int32(int32(22)),
					ToPort:     awssdk.Int32(int32(23)),
				},
			},
		},
		{
			params: map[string]any{
				"protocol":  "icmp",
				"cidr":      "10.0.0.0/16",
				"portrange": "any",
			},
			expected: []ec2types.IpPermission{
				{
					IpProtocol: awssdk.String("icmp"),
					IpRanges:   []ec2types.IpRange{{CidrIp: awssdk.String("10.0.0.0/16")}},
					FromPort:   awssdk.Int32(int32(-1)),
					ToPort:     awssdk.Int32(int32(-1)),
				},
			},
		},
		{
			params: map[string]any{
				"protocol":      "icmp",
				"securitygroup": "sg-12345",
				"portrange":     "any",
			},
			expected: []ec2types.IpPermission{
				{
					IpProtocol:       awssdk.String("icmp"),
					UserIdGroupPairs: []ec2types.UserIdGroupPair{{GroupId: awssdk.String("sg-12345")}},
					FromPort:         awssdk.Int32(int32(-1)),
					ToPort:           awssdk.Int32(int32(-1)),
				},
			},
		},

		{
			params: map[string]any{
				"protocol":      "tcp",
				"securitygroup": "sg-23456",
				"portrange":     80,
			},
			expected: []ec2types.IpPermission{
				{
					IpProtocol:       awssdk.String("tcp"),
					UserIdGroupPairs: []ec2types.UserIdGroupPair{{GroupId: awssdk.String("sg-23456")}},
					FromPort:         awssdk.Int32(int32(80)),
					ToPort:           awssdk.Int32(int32(80)),
				},
			},
		},
	}

	for i, tcase := range tcases {
		cmd := &UpdateSecuritygroup{}
		if err := cmd.inject(tcase.params); err != nil {
			t.Fatal(err)
		}
		ipPermissions, err := cmd.buildIPPermissions()
		if err != nil {
			t.Fatal(i+1, ":", err)
		}
		if got, want := ipPermissions, tcase.expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %+v, want %+v", i+1, got, want)
		}
	}
}

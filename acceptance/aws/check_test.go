package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// The check commands poll a describe call until the resource reaches a state or the
// timeout expires. Each test answers with the wanted state on the first call, so the
// assertion is that the command reads the state out of the right place in the response
// — a shape that differs for every service and is easy to get wrong.
//
// timeout is deliberately small: a mismatch would otherwise hang the suite rather than
// failing it.

func TestCheckInstanceReachesState(t *testing.T) {
	mock := NewMock().On("DescribeInstances", &ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{{
			Instances: []ec2types.Instance{{
				InstanceId: awssdk.String("i-1234"),
				State:      &ec2types.InstanceState{Name: ec2types.InstanceStateNameRunning},
			}},
		}},
	})

	Template("check instance id=i-1234 state=running timeout=1").
		Mock(mock).
		ExpectCalls("DescribeInstances").
		Run(t)
}

// An instance that AWS reports as not found must satisfy state=not-found rather than
// erroring, which is how a revert waits for a terminate to complete.
func TestCheckInstanceNotFound(t *testing.T) {
	mock := NewMock().OnAPIError("DescribeInstances", "InstanceNotFound", "not found")

	Template("check instance id=i-1234 state=not-found timeout=1").
		Mock(mock).
		ExpectCalls("DescribeInstances").
		Run(t)
}

func TestCheckVolumeReachesState(t *testing.T) {
	mock := NewMock().On("DescribeVolumes", &ec2.DescribeVolumesOutput{
		Volumes: []ec2types.Volume{{
			VolumeId: awssdk.String("vol-1234"),
			State:    ec2types.VolumeStateAvailable,
		}},
	})

	Template("check volume id=vol-1234 state=available timeout=1").
		Mock(mock).
		ExpectCalls("DescribeVolumes").
		Run(t)
}

func TestCheckSecuritygroupUnused(t *testing.T) {
	mock := NewMock().On("DescribeNetworkInterfaces", &ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []ec2types.NetworkInterface{},
	})

	// An empty interface list means nothing is using the group, which is what the
	// revert of create securitygroup waits for.
	Template("check securitygroup id=sg-1234 state=unused timeout=1").
		Mock(mock).
		ExpectCalls("DescribeNetworkInterfaces").
		Run(t)
}

func TestCheckNatgatewayReachesState(t *testing.T) {
	mock := NewMock().On("DescribeNatGateways", &ec2.DescribeNatGatewaysOutput{
		NatGateways: []ec2types.NatGateway{{
			NatGatewayId: awssdk.String("nat-1234"),
			State:        ec2types.NatGatewayStateAvailable,
		}},
	})

	Template("check natgateway id=nat-1234 state=available timeout=1").
		Mock(mock).
		ExpectCalls("DescribeNatGateways").
		Run(t)
}

func TestCheckNetworkinterfaceReachesState(t *testing.T) {
	mock := NewMock().On("DescribeNetworkInterfaces", &ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []ec2types.NetworkInterface{{
			NetworkInterfaceId: awssdk.String("eni-1234"),
			Status:             ec2types.NetworkInterfaceStatusAvailable,
		}},
	})

	Template("check networkinterface id=eni-1234 state=available timeout=1").
		Mock(mock).
		ExpectCalls("DescribeNetworkInterfaces").
		Run(t)
}

func TestCheckDatabaseReachesState(t *testing.T) {
	mock := NewMock().On("DescribeDBInstances", &rds.DescribeDBInstancesOutput{
		DBInstances: []rdstypes.DBInstance{{
			DBInstanceIdentifier: awssdk.String("mydb"),
			DBInstanceStatus:     awssdk.String("available"),
		}},
	})

	Template("check database id=mydb state=available timeout=1").
		Mock(mock).
		ExpectCalls("DescribeDBInstances").
		Run(t)
}

func TestCheckLoadbalancerReachesState(t *testing.T) {
	mock := NewMock().On("DescribeLoadBalancers", &elasticloadbalancingv2.DescribeLoadBalancersOutput{
		LoadBalancers: []elbv2types.LoadBalancer{{
			LoadBalancerArn: awssdk.String("arn:aws:elasticloadbalancing:us-west-2:1:loadbalancer/app/my-lb/abcd"),
			State:           &elbv2types.LoadBalancerState{Code: elbv2types.LoadBalancerStateEnumActive},
		}},
	})

	Template("check loadbalancer id=arn:aws:elasticloadbalancing:us-west-2:1:loadbalancer/app/my-lb/abcd state=active timeout=1").
		Mock(mock).
		ExpectCalls("DescribeLoadBalancers").
		Run(t)
}

func TestCheckScalinggroupReachesCount(t *testing.T) {
	mock := NewMock().On("DescribeAutoScalingGroups", &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []autoscalingtypes.AutoScalingGroup{{
			AutoScalingGroupName: awssdk.String("my-asg"),
			Instances: []autoscalingtypes.Instance{
				{LifecycleState: autoscalingtypes.LifecycleStateInService},
				{LifecycleState: autoscalingtypes.LifecycleStateInService},
			},
		}},
	})

	// check scalinggroup counts instances in service rather than reading one state
	// field, so the assertion is on the count.
	Template("check scalinggroup name=my-asg count=2 timeout=1").
		Mock(mock).
		ExpectCalls("DescribeAutoScalingGroups").
		Run(t)
}

func TestCheckCertificateReachesState(t *testing.T) {
	mock := NewMock().On("DescribeCertificate", &acm.DescribeCertificateOutput{
		Certificate: &acmtypes.CertificateDetail{
			CertificateArn: awssdk.String("arn:aws:acm:us-west-2:1:certificate/abcd"),
			Status:         acmtypes.CertificateStatusIssued,
		},
	})

	Template("check certificate arn=arn:aws:acm:us-west-2:1:certificate/abcd state=issued timeout=1").
		Mock(mock).
		ExpectCalls("DescribeCertificate").
		Run(t)
}

func TestCheckDistributionReachesState(t *testing.T) {
	mock := NewMock().On("GetDistribution", &cloudfront.GetDistributionOutput{
		Distribution: &cloudfronttypes.Distribution{
			Id:     awssdk.String("E1PA6795SAMPLE"),
			Status: awssdk.String("Deployed"),
		},
	})

	Template("check distribution id=E1PA6795SAMPLE state=Deployed timeout=1").
		Mock(mock).
		ExpectCalls("GetDistribution").
		Run(t)
}

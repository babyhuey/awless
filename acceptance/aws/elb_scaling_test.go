package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

// Load balancing and auto scaling. Most of these take several numeric params that are
// interchangeable at the type level — min, max and desired capacity; health check
// interval, timeout and thresholds — so the tests pin which value lands where.

func TestCreateLoadbalancer(t *testing.T) {
	arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-lb/abcd"
	mock := NewMock().On("CreateLoadBalancer", &elasticloadbalancingv2.CreateLoadBalancerOutput{
		LoadBalancers: []elbv2types.LoadBalancer{{LoadBalancerArn: awssdk.String(arn)}},
	})

	Template("create loadbalancer name=my-lb subnets=subnet-a,subnet-b scheme=internet-facing").
		Mock(mock).
		ExpectCalls("CreateLoadBalancer").
		ExpectCommandResult(arn).
		Run(t)

	in := mock.InputFor("CreateLoadBalancer").(*elasticloadbalancingv2.CreateLoadBalancerInput)
	if len(in.Subnets) != 2 {
		t.Errorf("expected both subnets, got %#v", in.Subnets)
	}
	if got := string(in.Scheme); got != "internet-facing" {
		t.Errorf("Scheme: got %q, want internet-facing", got)
	}
}

func TestDeleteLoadbalancer(t *testing.T) {
	mock := NewMock().On("DeleteLoadBalancer", &elasticloadbalancingv2.DeleteLoadBalancerOutput{})

	Template("delete loadbalancer id=arn:aws:elasticloadbalancing:us-west-2:1:loadbalancer/app/my-lb/abcd").
		Mock(mock).
		ExpectCalls("DeleteLoadBalancer").
		Run(t)
}

func TestCreateListener(t *testing.T) {
	arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:listener/app/my-lb/abcd/ef01"
	mock := NewMock().On("CreateListener", &elasticloadbalancingv2.CreateListenerOutput{
		Listeners: []elbv2types.Listener{{ListenerArn: awssdk.String(arn)}},
	})

	Template("create listener actiontype=forward loadbalancer=arn:aws:elasticloadbalancing:us-west-2:1:loadbalancer/app/my-lb/abcd port=443 protocol=HTTPS targetgroup=arn:aws:elasticloadbalancing:us-west-2:1:targetgroup/my-tg/abcd certificate=arn:aws:acm:us-west-2:1:certificate/abcd").
		Mock(mock).
		ExpectCalls("CreateListener").
		ExpectCommandResult(arn).
		Run(t)

	in := mock.InputFor("CreateListener").(*elasticloadbalancingv2.CreateListenerInput)
	if got := awssdk.ToInt32(in.Port); got != 443 {
		t.Errorf("Port: got %d, want 443", got)
	}
	if got := string(in.Protocol); got != "HTTPS" {
		t.Errorf("Protocol: got %q, want HTTPS", got)
	}
	// An HTTPS listener without its certificate would be rejected by AWS.
	if len(in.Certificates) != 1 {
		t.Errorf("expected the certificate to reach the input, got %#v", in.Certificates)
	}
}

func TestDeleteListener(t *testing.T) {
	mock := NewMock().On("DeleteListener", &elasticloadbalancingv2.DeleteListenerOutput{})

	Template("delete listener id=arn:aws:elasticloadbalancing:us-west-2:1:listener/app/my-lb/abcd/ef01").
		Mock(mock).
		ExpectCalls("DeleteListener").
		Run(t)
}

func TestCreateTargetgroupHealthCheckFields(t *testing.T) {
	arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd"
	mock := NewMock().On("CreateTargetGroup", &elasticloadbalancingv2.CreateTargetGroupOutput{
		TargetGroups: []elbv2types.TargetGroup{{TargetGroupArn: awssdk.String(arn)}},
	})

	Template("create targetgroup name=my-tg port=80 protocol=HTTP vpc=vpc-1234 healthcheckinterval=30 healthchecktimeout=5 healthythreshold=3 unhealthythreshold=2 healthcheckpath=/health").
		Mock(mock).
		ExpectCalls("CreateTargetGroup").
		ExpectCommandResult(arn).
		Run(t)

	in := mock.InputFor("CreateTargetGroup").(*elasticloadbalancingv2.CreateTargetGroupInput)
	// Four numeric health-check fields, all int32, trivially swappable.
	if got := awssdk.ToInt32(in.HealthCheckIntervalSeconds); got != 30 {
		t.Errorf("HealthCheckIntervalSeconds: got %d, want 30", got)
	}
	if got := awssdk.ToInt32(in.HealthCheckTimeoutSeconds); got != 5 {
		t.Errorf("HealthCheckTimeoutSeconds: got %d, want 5", got)
	}
	if got := awssdk.ToInt32(in.HealthyThresholdCount); got != 3 {
		t.Errorf("HealthyThresholdCount: got %d, want 3", got)
	}
	if got := awssdk.ToInt32(in.UnhealthyThresholdCount); got != 2 {
		t.Errorf("UnhealthyThresholdCount: got %d, want 2", got)
	}
	if got := awssdk.ToString(in.HealthCheckPath); got != "/health" {
		t.Errorf("HealthCheckPath: got %q, want /health", got)
	}
}

func TestDeleteTargetgroup(t *testing.T) {
	mock := NewMock().On("DeleteTargetGroup", &elasticloadbalancingv2.DeleteTargetGroupOutput{})

	Template("delete targetgroup id=arn:aws:elasticloadbalancing:us-west-2:1:targetgroup/my-tg/abcd").
		Mock(mock).
		ExpectCalls("DeleteTargetGroup").
		Run(t)
}

func TestCreateScalinggroupCapacities(t *testing.T) {
	mock := NewMock().On("CreateAutoScalingGroup", &autoscaling.CreateAutoScalingGroupOutput{})

	Template("create scalinggroup name=my-asg launchconfiguration=my-lc min-size=1 max-size=4 desired-capacity=2 subnets=subnet-a,subnet-b cooldown=300").
		Mock(mock).
		ExpectCalls("CreateAutoScalingGroup").
		Run(t)

	in := mock.InputFor("CreateAutoScalingGroup").(*autoscaling.CreateAutoScalingGroupInput)
	// min, max and desired are three int32s; getting them crossed would silently
	// launch the wrong number of instances.
	if got := awssdk.ToInt32(in.MinSize); got != 1 {
		t.Errorf("MinSize: got %d, want 1", got)
	}
	if got := awssdk.ToInt32(in.MaxSize); got != 4 {
		t.Errorf("MaxSize: got %d, want 4", got)
	}
	if got := awssdk.ToInt32(in.DesiredCapacity); got != 2 {
		t.Errorf("DesiredCapacity: got %d, want 2", got)
	}
	if got := awssdk.ToInt32(in.DefaultCooldown); got != 300 {
		t.Errorf("DefaultCooldown: got %d, want 300", got)
	}
}

func TestUpdateScalinggroup(t *testing.T) {
	mock := NewMock().On("UpdateAutoScalingGroup", &autoscaling.UpdateAutoScalingGroupOutput{})

	Template("update scalinggroup name=my-asg desired-capacity=3").
		Mock(mock).
		ExpectCalls("UpdateAutoScalingGroup").
		Run(t)

	in := mock.InputFor("UpdateAutoScalingGroup").(*autoscaling.UpdateAutoScalingGroupInput)
	if got := awssdk.ToInt32(in.DesiredCapacity); got != 3 {
		t.Errorf("DesiredCapacity: got %d, want 3", got)
	}
	// An update must not send a zero min or max that the user did not ask for.
	if in.MinSize != nil || in.MaxSize != nil {
		t.Errorf("expected only the requested field to be sent, got min=%v max=%v", in.MinSize, in.MaxSize)
	}
}

func TestDeleteScalinggroup(t *testing.T) {
	mock := NewMock().On("DeleteAutoScalingGroup", &autoscaling.DeleteAutoScalingGroupOutput{})

	Template("delete scalinggroup name=my-asg force=true").
		Mock(mock).
		ExpectCalls("DeleteAutoScalingGroup").
		Run(t)

	in := mock.InputFor("DeleteAutoScalingGroup").(*autoscaling.DeleteAutoScalingGroupInput)
	if !awssdk.ToBool(in.ForceDelete) {
		t.Error("expected ForceDelete to be set")
	}
}

func TestCreateScalingpolicy(t *testing.T) {
	arn := "arn:aws:autoscaling:us-west-2:123456789012:scalingPolicy:abcd:autoScalingGroupName/my-asg:policyName/scale-out"
	mock := NewMock().On("PutScalingPolicy", &autoscaling.PutScalingPolicyOutput{
		PolicyARN: awssdk.String(arn),
	})

	Template("create scalingpolicy name=scale-out scalinggroup=my-asg adjustment-type=ChangeInCapacity adjustment-scaling=2 cooldown=300").
		Mock(mock).
		ExpectCalls("PutScalingPolicy").
		ExpectCommandResult(arn).
		Run(t)

	in := mock.InputFor("PutScalingPolicy").(*autoscaling.PutScalingPolicyInput)
	if got := awssdk.ToInt32(in.ScalingAdjustment); got != 2 {
		t.Errorf("ScalingAdjustment: got %d, want 2", got)
	}
	if got := awssdk.ToString(in.AdjustmentType); got != "ChangeInCapacity" {
		t.Errorf("AdjustmentType: got %q", got)
	}
}

func TestDeleteScalingpolicy(t *testing.T) {
	mock := NewMock().On("DeletePolicy", &autoscaling.DeletePolicyOutput{})

	Template("delete scalingpolicy id=scale-out").
		Mock(mock).
		ExpectCalls("DeletePolicy").
		Run(t)
}

func TestDeleteLaunchconfiguration(t *testing.T) {
	mock := NewMock().On("DeleteLaunchConfiguration", &autoscaling.DeleteLaunchConfigurationOutput{})

	Template("delete launchconfiguration name=my-lc").
		Mock(mock).
		ExpectCalls("DeleteLaunchConfiguration").
		Run(t)
}

func TestCreateAppscalingtarget(t *testing.T) {
	mock := NewMock().On("RegisterScalableTarget", &applicationautoscaling.RegisterScalableTargetOutput{})

	Template("create appscalingtarget dimension=ecs:service:DesiredCount resource=service/my-cluster/my-svc service-namespace=ecs min-capacity=1 max-capacity=10 role=arn:aws:iam::123456789012:role/ecsAutoscaleRole").
		Mock(mock).
		ExpectCalls("RegisterScalableTarget").
		Run(t)

	in := mock.InputFor("RegisterScalableTarget").(*applicationautoscaling.RegisterScalableTargetInput)
	if got := awssdk.ToInt32(in.MinCapacity); got != 1 {
		t.Errorf("MinCapacity: got %d, want 1", got)
	}
	if got := awssdk.ToInt32(in.MaxCapacity); got != 10 {
		t.Errorf("MaxCapacity: got %d, want 10", got)
	}
	if got := string(in.ScalableDimension); got != "ecs:service:DesiredCount" {
		t.Errorf("ScalableDimension: got %q", got)
	}
}

func TestDeleteAppscalingtarget(t *testing.T) {
	mock := NewMock().On("DeregisterScalableTarget", &applicationautoscaling.DeregisterScalableTargetOutput{})

	Template("delete appscalingtarget dimension=ecs:service:DesiredCount resource=service/my-cluster/my-svc service-namespace=ecs").
		Mock(mock).
		ExpectCalls("DeregisterScalableTarget").
		Run(t)
}

func TestDeleteAppscalingpolicy(t *testing.T) {
	mock := NewMock().On("DeleteScalingPolicy", &applicationautoscaling.DeleteScalingPolicyOutput{})

	Template("delete appscalingpolicy name=scale-out dimension=ecs:service:DesiredCount resource=service/my-cluster/my-svc service-namespace=ecs").
		Mock(mock).
		ExpectCalls("DeleteScalingPolicy").
		Run(t)
}

// The classic load balancer is v1 of the API, kept for older accounts.
func TestCreateClassicloadbalancer(t *testing.T) {
	mock := NewMock().
		On("CreateLoadBalancer", &elasticloadbalancing.CreateLoadBalancerOutput{
			DNSName: awssdk.String("my-clb-123456.us-west-2.elb.amazonaws.com"),
		}).
		// AfterRun configures the health check, since the create API cannot.
		On("ConfigureHealthCheck", &elasticloadbalancing.ConfigureHealthCheckOutput{})

	Template("create classicloadbalancer name=my-clb listeners=HTTP:80:HTTP:80 subnets=subnet-a").
		Mock(mock).
		ExpectCalls("CreateLoadBalancer", "ConfigureHealthCheck").
		Run(t)

	in := mock.InputFor("CreateLoadBalancer").(*elasticloadbalancing.CreateLoadBalancerInput)
	if got := awssdk.ToString(in.LoadBalancerName); got != "my-clb" {
		t.Errorf("LoadBalancerName: got %q, want my-clb", got)
	}
	if len(in.Listeners) != 1 {
		t.Fatalf("expected one listener, got %#v", in.Listeners)
	}
	// listeners is a colon-packed string that has to be parsed into a struct.
	l := in.Listeners[0]
	if got := l.LoadBalancerPort; got != 80 {
		t.Errorf("LoadBalancerPort: got %d, want 80", got)
	}
}

func TestDeleteClassicloadbalancer(t *testing.T) {
	mock := NewMock().On("DeleteLoadBalancer", &elasticloadbalancing.DeleteLoadBalancerOutput{})

	Template("delete classicloadbalancer name=my-clb").
		Mock(mock).
		ExpectCalls("DeleteLoadBalancer").
		Run(t)
}

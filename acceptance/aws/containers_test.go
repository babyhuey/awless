package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// ECS, IAM roles on instance profiles, the classic load balancer and application
// auto scaling.

func TestCreateContainercluster(t *testing.T) {
	mock := NewMock().On("CreateCluster", &ecs.CreateClusterOutput{
		Cluster: &ecstypes.Cluster{ClusterArn: awssdk.String("arn:aws:ecs:us-west-2:1:cluster/my-cluster")},
	})

	Template("create containercluster name=my-cluster").
		Mock(mock).
		ExpectCalls("CreateCluster").
		Run(t)

	in := mock.InputFor("CreateCluster").(*ecs.CreateClusterInput)
	if got := awssdk.ToString(in.ClusterName); got != "my-cluster" {
		t.Errorf("ClusterName: got %q, want my-cluster", got)
	}
}

func TestDeleteContainercluster(t *testing.T) {
	mock := NewMock().On("DeleteCluster", &ecs.DeleteClusterOutput{})

	Template("delete containercluster id=my-cluster").
		Mock(mock).
		ExpectCalls("DeleteCluster").
		Run(t)
}

// attach containertask registers a container definition, where the memory limit and
// the port mapping both have to be parsed out of strings.
func TestAttachContainertask(t *testing.T) {
	mock := NewMock().
		On("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionOutput{
			TaskDefinition: &ecstypes.TaskDefinition{
				Family:               awssdk.String("my-task"),
				ContainerDefinitions: []ecstypes.ContainerDefinition{},
			},
		}).
		On("RegisterTaskDefinition", &ecs.RegisterTaskDefinitionOutput{
			TaskDefinition: &ecstypes.TaskDefinition{
				TaskDefinitionArn: awssdk.String("arn:aws:ecs:us-west-2:1:task-definition/my-task:2"),
			},
		})

	Template("attach containertask name=my-task container-name=web image=nginx:latest memory-hard-limit=512 ports=80:8080 env=ENVIRONMENT:production").
		Mock(mock).
		ExpectCalls("DescribeTaskDefinition", "RegisterTaskDefinition").
		Run(t)

	in := mock.InputFor("RegisterTaskDefinition").(*ecs.RegisterTaskDefinitionInput)
	if len(in.ContainerDefinitions) != 1 {
		t.Fatalf("expected one container definition, got %#v", in.ContainerDefinitions)
	}
	def := in.ContainerDefinitions[0]
	if got := awssdk.ToString(def.Name); got != "web" {
		t.Errorf("container name: got %q, want web", got)
	}
	if got := awssdk.ToString(def.Image); got != "nginx:latest" {
		t.Errorf("image: got %q, want nginx:latest", got)
	}
	if got := awssdk.ToInt32(def.Memory); got != 512 {
		t.Errorf("memory: got %d, want 512", got)
	}
	// ports=80:8080 means host 80 maps to container 8080; reversing them would expose
	// the wrong port.
	if len(def.PortMappings) != 1 {
		t.Fatalf("expected one port mapping, got %#v", def.PortMappings)
	}
	if got := awssdk.ToInt32(def.PortMappings[0].HostPort); got != 80 {
		t.Errorf("HostPort: got %d, want 80", got)
	}
	if got := awssdk.ToInt32(def.PortMappings[0].ContainerPort); got != 8080 {
		t.Errorf("ContainerPort: got %d, want 8080", got)
	}
	if len(def.Environment) != 1 || awssdk.ToString(def.Environment[0].Name) != "ENVIRONMENT" {
		t.Errorf("environment: got %#v", def.Environment)
	}
}

func TestDetachContainertask(t *testing.T) {
	mock := NewMock().
		On("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionOutput{
			TaskDefinition: &ecstypes.TaskDefinition{
				Family: awssdk.String("my-task"),
				ContainerDefinitions: []ecstypes.ContainerDefinition{
					{Name: awssdk.String("web")},
					{Name: awssdk.String("sidecar")},
				},
			},
		}).
		On("RegisterTaskDefinition", &ecs.RegisterTaskDefinitionOutput{
			TaskDefinition: &ecstypes.TaskDefinition{
				TaskDefinitionArn: awssdk.String("arn:aws:ecs:us-west-2:1:task-definition/my-task:3"),
			},
		})

	Template("detach containertask name=my-task container-name=web").
		Mock(mock).
		ExpectCalls("DescribeTaskDefinition", "RegisterTaskDefinition").
		Run(t)

	in := mock.InputFor("RegisterTaskDefinition").(*ecs.RegisterTaskDefinitionInput)
	// Removing one container must leave the others in the new revision.
	if len(in.ContainerDefinitions) != 1 {
		t.Fatalf("expected the sidecar to remain, got %#v", in.ContainerDefinitions)
	}
	if got := awssdk.ToString(in.ContainerDefinitions[0].Name); got != "sidecar" {
		t.Errorf("remaining container: got %q, want sidecar", got)
	}
}

// start containertask type=task runs a one-off task; type=service creates a service.
func TestStartContainertaskAsTask(t *testing.T) {
	mock := NewMock().On("RunTask", &ecs.RunTaskOutput{
		Tasks: []ecstypes.Task{{TaskArn: awssdk.String("arn:aws:ecs:us-west-2:1:task/my-cluster/abcd")}},
	})

	Template("start containertask cluster=my-cluster name=my-task type=task desired-count=2").
		Mock(mock).
		ExpectCalls("RunTask").
		Run(t)

	in := mock.InputFor("RunTask").(*ecs.RunTaskInput)
	if got := awssdk.ToString(in.Cluster); got != "my-cluster" {
		t.Errorf("Cluster: got %q, want my-cluster", got)
	}
	if got := awssdk.ToInt32(in.Count); got != 2 {
		t.Errorf("Count: got %d, want 2", got)
	}
}

func TestStartContainertaskAsService(t *testing.T) {
	mock := NewMock().On("CreateService", &ecs.CreateServiceOutput{
		Service: &ecstypes.Service{ServiceName: awssdk.String("web-svc")},
	})

	Template("start containertask cluster=my-cluster name=my-task type=service desired-count=3 deployment-name=web-svc").
		Mock(mock).
		ExpectCalls("CreateService").
		Run(t)

	in := mock.InputFor("CreateService").(*ecs.CreateServiceInput)
	if got := awssdk.ToInt32(in.DesiredCount); got != 3 {
		t.Errorf("DesiredCount: got %d, want 3", got)
	}
	if got := awssdk.ToString(in.ServiceName); got != "web-svc" {
		t.Errorf("ServiceName: got %q, want web-svc", got)
	}
}

func TestStopContainertaskService(t *testing.T) {
	mock := NewMock().On("DeleteService", &ecs.DeleteServiceOutput{
		Service: &ecstypes.Service{ServiceName: awssdk.String("web-svc")},
	})

	Template("stop containertask cluster=my-cluster type=service deployment-name=web-svc").
		Mock(mock).
		ExpectCalls("DeleteService").
		Run(t)
}

func TestUpdateContainertaskDesiredCount(t *testing.T) {
	mock := NewMock().On("UpdateService", &ecs.UpdateServiceOutput{
		Service: &ecstypes.Service{ServiceName: awssdk.String("web-svc")},
	})

	Template("update containertask cluster=my-cluster deployment-name=web-svc desired-count=5").
		Mock(mock).
		ExpectCalls("UpdateService").
		Run(t)

	in := mock.InputFor("UpdateService").(*ecs.UpdateServiceInput)
	if got := awssdk.ToInt32(in.DesiredCount); got != 5 {
		t.Errorf("DesiredCount: got %d, want 5", got)
	}
}

// delete role also tears down the instance profile create role made, so the pair is
// symmetric: neither leaves an orphan behind.
func TestDeleteRoleAlsoRemovesItsInstanceProfile(t *testing.T) {
	mock := NewMock().
		On("RemoveRoleFromInstanceProfile", &iam.RemoveRoleFromInstanceProfileOutput{}).
		On("DeleteInstanceProfile", &iam.DeleteInstanceProfileOutput{}).
		On("DeleteRole", &iam.DeleteRoleOutput{})

	Template("delete role name=my-role").
		Mock(mock).
		ExpectCalls("RemoveRoleFromInstanceProfile", "DeleteInstanceProfile", "DeleteRole").
		Run(t)

	in := mock.InputFor("RemoveRoleFromInstanceProfile").(*iam.RemoveRoleFromInstanceProfileInput)
	if got := awssdk.ToString(in.RoleName); got != "my-role" {
		t.Errorf("RoleName: got %q, want my-role", got)
	}
}

// A role is attached to an instance profile, not to an instance directly.
func TestAttachAndDetachRole(t *testing.T) {
	attach := NewMock().On("AddRoleToInstanceProfile", &iam.AddRoleToInstanceProfileOutput{})
	Template("attach role name=my-role instanceprofile=my-profile").
		Mock(attach).
		ExpectCalls("AddRoleToInstanceProfile").
		ExpectRevert("detach role instanceprofile=my-profile name=my-role").
		Run(t)

	in := attach.InputFor("AddRoleToInstanceProfile").(*iam.AddRoleToInstanceProfileInput)
	if got := awssdk.ToString(in.InstanceProfileName); got != "my-profile" {
		t.Errorf("InstanceProfileName: got %q, want my-profile", got)
	}
	if got := awssdk.ToString(in.RoleName); got != "my-role" {
		t.Errorf("RoleName: got %q, want my-role", got)
	}

	detach := NewMock().On("RemoveRoleFromInstanceProfile", &iam.RemoveRoleFromInstanceProfileOutput{})
	Template("detach role name=my-role instanceprofile=my-profile").
		Mock(detach).
		ExpectCalls("RemoveRoleFromInstanceProfile").
		Run(t)
}

func TestAttachAndDetachClassicloadbalancer(t *testing.T) {
	attach := NewMock().On("RegisterInstancesWithLoadBalancer",
		&elasticloadbalancing.RegisterInstancesWithLoadBalancerOutput{})
	Template("attach classicloadbalancer name=my-clb instance=i-1234").
		Mock(attach).
		ExpectCalls("RegisterInstancesWithLoadBalancer").
		ExpectRevert("detach classicloadbalancer instance=i-1234 name=my-clb").
		Run(t)

	in := attach.InputFor("RegisterInstancesWithLoadBalancer").(*elasticloadbalancing.RegisterInstancesWithLoadBalancerInput)
	if len(in.Instances) != 1 {
		t.Fatalf("expected one instance, got %#v", in.Instances)
	}
	if got := awssdk.ToString(in.Instances[0].InstanceId); got != "i-1234" {
		t.Errorf("InstanceId: got %q, want i-1234", got)
	}

	detach := NewMock().On("DeregisterInstancesFromLoadBalancer",
		&elasticloadbalancing.DeregisterInstancesFromLoadBalancerOutput{})
	Template("detach classicloadbalancer name=my-clb instance=i-1234").
		Mock(detach).
		ExpectCalls("DeregisterInstancesFromLoadBalancer").
		Run(t)
}

func TestUpdateClassicloadbalancerHealthCheck(t *testing.T) {
	mock := NewMock().On("ConfigureHealthCheck", &elasticloadbalancing.ConfigureHealthCheckOutput{})

	Template("update classicloadbalancer name=my-clb healthy-threshold=3 unhealthy-threshold=2 health-interval=30 health-timeout=5 health-target=HTTP:80/health").
		Mock(mock).
		ExpectCalls("ConfigureHealthCheck").
		Run(t)

	in := mock.InputFor("ConfigureHealthCheck").(*elasticloadbalancing.ConfigureHealthCheckInput)
	hc := in.HealthCheck
	if got := awssdk.ToInt32(hc.HealthyThreshold); got != 3 {
		t.Errorf("HealthyThreshold: got %d, want 3", got)
	}
	if got := awssdk.ToInt32(hc.UnhealthyThreshold); got != 2 {
		t.Errorf("UnhealthyThreshold: got %d, want 2", got)
	}
	// Interval and timeout are both seconds and adjacent in the struct.
	if got := awssdk.ToInt32(hc.Interval); got != 30 {
		t.Errorf("Interval: got %d, want 30", got)
	}
	if got := awssdk.ToInt32(hc.Timeout); got != 5 {
		t.Errorf("Timeout: got %d, want 5", got)
	}
	if got := awssdk.ToString(hc.Target); got != "HTTP:80/health" {
		t.Errorf("Target: got %q", got)
	}
}

func TestCreateAppscalingpolicy(t *testing.T) {
	arn := "arn:aws:autoscaling:us-west-2:1:scalingPolicy:abcd:resource/ecs/service/my-cluster/my-svc:policyName/scale-out"
	mock := NewMock().On("PutScalingPolicy", &applicationautoscaling.PutScalingPolicyOutput{
		PolicyARN: awssdk.String(arn),
	})

	Template("create appscalingpolicy name=scale-out dimension=ecs:service:DesiredCount resource=service/my-cluster/my-svc service-namespace=ecs type=StepScaling stepscaling-adjustment-type=ChangeInCapacity stepscaling-adjustments=0::1").
		Mock(mock).
		ExpectCalls("PutScalingPolicy").
		Run(t)

	in := mock.InputFor("PutScalingPolicy").(*applicationautoscaling.PutScalingPolicyInput)
	if got := string(in.PolicyType); got != "StepScaling" {
		t.Errorf("PolicyType: got %q, want StepScaling", got)
	}
	if in.StepScalingPolicyConfiguration == nil {
		t.Fatal("expected a step scaling configuration")
	}
	// The adjustments arrive as colon-packed from:to:adjustment triples.
	if len(in.StepScalingPolicyConfiguration.StepAdjustments) != 1 {
		t.Errorf("StepAdjustments: got %#v, want one entry", in.StepScalingPolicyConfiguration.StepAdjustments)
	}
}

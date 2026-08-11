package console

import (
	"testing"

	"github.com/bootswithdefer/awless/cloud"
)

// Every listable resource needs both display definitions: the short property list that
// `awless ls` shows by default, and the full column block that gives those columns their
// order, friendly names and formatting.
//
// Nothing enforced this before, so a resource could be fetched and synced while `awless ls`
// had nothing to show for it. Three resources shipped that way — the column entries were
// simply forgotten, and no test noticed.
func TestEveryResourceHasDisplayColumns(t *testing.T) {
	// Resource types with no fetcher. These exist as command entities — you can create
	// and delete them — but nothing syncs them, so `awless ls` never reaches them and
	// display columns would be dead weight. Verified against the binary rather than
	// assumed: none of them appear in `awless list -h`.
	notListed := map[string]bool{
		cloud.Region:           true,
		cloud.ImportImageTask:  true,
		cloud.ACL:              true,
		cloud.LoginProfile:     true,
		cloud.Registry:         true,
		cloud.ContainerService: true,
		cloud.AppScalingTarget: true,
		cloud.AppScalingPolicy: true,
	}

	for _, resource := range allResourceTypes() {
		if notListed[resource] {
			continue
		}
		t.Run(resource, func(t *testing.T) {
			if props := ColumnsInListing[resource]; len(props) == 0 {
				t.Errorf("no entry in ColumnsInListing for %q; "+
					"awless ls would have no columns to show", resource)
			}
			if cols := DefaultsColumnDefinitions[resource]; len(cols) == 0 {
				t.Errorf("no entry in DefaultsColumnDefinitions for %q; "+
					"the columns would have no order or formatting", resource)
			}
		})
	}
}

// allResourceTypes lists every resource constant. Spelled out rather than reflected so a
// new resource has to be considered here, which is the point.
func allResourceTypes() []string {
	return []string{
		cloud.Vpc, cloud.Subnet, cloud.Image, cloud.SecurityGroup, cloud.AvailabilityZone,
		cloud.Keypair, cloud.Volume, cloud.Instance, cloud.InstanceProfile,
		cloud.InternetGateway, cloud.NatGateway, cloud.RouteTable, cloud.ElasticIP,
		cloud.Snapshot, cloud.NetworkInterface, cloud.Certificate,
		cloud.ClassicLoadBalancer, cloud.LoadBalancer, cloud.TargetGroup, cloud.Listener,
		cloud.Database, cloud.DBSubnetGroup,
		cloud.User, cloud.Role, cloud.Group, cloud.Policy, cloud.AccessKey,
		cloud.LoginProfile, cloud.MFADevice,
		cloud.Bucket, cloud.S3Object,
		cloud.Subscription, cloud.Topic, cloud.Queue, cloud.Zone, cloud.Record,
		cloud.Function, cloud.LaunchConfiguration, cloud.ScalingGroup, cloud.ScalingPolicy,
		cloud.Metric, cloud.Alarm, cloud.Distribution, cloud.Stack,
		cloud.Repository, cloud.Registry, cloud.ContainerCluster, cloud.ContainerService,
		cloud.ContainerTask, cloud.Container, cloud.ContainerInstance,
		cloud.AppScalingTarget, cloud.AppScalingPolicy,
		cloud.EKSCluster, cloud.EKSNodeGroup, cloud.DynamoDBTable, cloud.Secret, cloud.Key,
		cloud.APIGateway, cloud.APIGatewayRoute, cloud.APIGatewayStage, cloud.SSMParameter,
		cloud.FileSystem, cloud.MountTarget, cloud.Trail, cloud.LogGroup,
		cloud.CacheCluster, cloud.ReplicationGroup, cloud.CacheSubnetGroup,
		cloud.EventBus, cloud.EventRule, cloud.StateMachine,
		cloud.WebACL, cloud.IPSet, cloud.RuleGroup, cloud.ConfigRule, cloud.Stream,
		cloud.RedshiftCluster, cloud.RedshiftSubnetGroup, cloud.Pipeline,
		cloud.BuildProject, cloud.Application, cloud.Environment,
	}
}

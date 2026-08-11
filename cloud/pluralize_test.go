package cloud

import "testing"

// Pluralize and Singularize must be exact inverses for every resource type, because the
// plural is what a user types (`awless ls eventbuses`) and the singular is what the
// registry is keyed on. A mismatch makes the resource unlistable: `eventbus` pluralized
// to `eventbuss`, which no lookup could resolve back.
//
// The list is every resource constant in this file. It is spelled out rather than
// derived so that a new resource with an awkward ending has to be considered here.
func TestPluralizeRoundTripsForEveryResource(t *testing.T) {
	for _, singular := range []string{
		Region, Vpc, Subnet, Image, ImportImageTask, SecurityGroup, AvailabilityZone,
		Keypair, Volume, Instance, InstanceProfile, InternetGateway, NatGateway,
		RouteTable, ElasticIP, Snapshot, NetworkInterface, Certificate,
		ClassicLoadBalancer, LoadBalancer, TargetGroup, Listener,
		Database, DBSubnetGroup,
		User, Role, Group, Policy, AccessKey, LoginProfile, MFADevice,
		Bucket, S3Object, ACL,
		Subscription, Topic, Queue, Zone, Record, Function,
		LaunchConfiguration, ScalingGroup, ScalingPolicy,
		Metric, Alarm, Distribution, Stack,
		Repository, Registry, ContainerCluster, ContainerService, ContainerTask,
		Container, ContainerInstance,
		AppScalingTarget, AppScalingPolicy,
		EKSCluster, EKSNodeGroup, DynamoDBTable, Secret, Key,
		APIGateway, APIGatewayRoute, APIGatewayStage, SSMParameter,
		FileSystem, MountTarget, Trail, LogGroup,
		CacheCluster, ReplicationGroup, CacheSubnetGroup,
		EventBus, EventRule,
	} {
		plural := PluralizeResource(singular)
		if plural == singular {
			t.Errorf("%q pluralizes to itself", singular)
		}
		if got := SingularizeResource(plural); got != singular {
			t.Errorf("round trip failed: %q -> %q -> %q", singular, plural, got)
		}
	}
}

// The specific forms that motivated the rules, so a future simplification cannot quietly
// break them.
func TestPluralizeSpecificForms(t *testing.T) {
	for _, tc := range []struct{ singular, plural string }{
		{"eventbus", "eventbuses"}, // trailing s takes es
		{"database", "databases"},  // ends in ses when plural, but is not an s-word
		{"policy", "policies"},     // y to ies
		{"registry", "registries"}, // y to ies
		{"instance", "instances"},  // the ordinary case
		{"s3object", "s3objects"},  // digits are not special
		{"storageacl", "storageacls"},
	} {
		if got := PluralizeResource(tc.singular); got != tc.plural {
			t.Errorf("PluralizeResource(%q) = %q, want %q", tc.singular, got, tc.plural)
		}
		if got := SingularizeResource(tc.plural); got != tc.singular {
			t.Errorf("SingularizeResource(%q) = %q, want %q", tc.plural, got, tc.singular)
		}
	}
}

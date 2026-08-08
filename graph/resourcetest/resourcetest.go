package resourcetest

import (
	"fmt"
	"strings"

	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/graph"
)

type rBuilder struct {
	id, typ string
	props   map[string]any
}

func newResource(typ, id string) *rBuilder {
	r := &rBuilder{id: id, typ: typ, props: make(map[string]any)}
	return r.Prop(properties.ID, id)
}

func Region(id string) *rBuilder {
	return newResource("region", id)
}

func Instance(id string) *rBuilder {
	return newResource("instance", id)
}

func Subnet(id string) *rBuilder {
	return newResource("subnet", id)
}

func VPC(id string) *rBuilder {
	return newResource("vpc", id)
}

func SecurityGroup(id string) *rBuilder {
	return newResource("securitygroup", id)
}

func KeyPair(id string) *rBuilder {
	return newResource("keypair", id)
}

func InternetGw(id string) *rBuilder {
	return newResource("internetgateway", id)
}

func NatGw(id string) *rBuilder {
	return newResource("natgateway", id)
}

func RouteTable(id string) *rBuilder {
	return newResource("routetable", id)
}

func LoadBalancer(id string) *rBuilder {
	return newResource("loadbalancer", id)
}

func ClassicLoadBalancer(id string) *rBuilder {
	return newResource("classicloadbalancer", id)
}

func AvailabilityZone(id string) *rBuilder {
	return newResource("availabilityzone", id)
}

func TargetGroup(id string) *rBuilder {
	return newResource("targetgroup", id)
}

func Policy(id string) *rBuilder {
	return newResource("policy", id)
}

func Group(id string) *rBuilder {
	return newResource("group", id)
}

func Role(id string) *rBuilder {
	return newResource("role", id)
}

func User(id string) *rBuilder {
	return newResource("user", id)
}

func MfaDevice(id string) *rBuilder {
	return newResource("mfadevice", id)
}

func Listener(id string) *rBuilder {
	return newResource("listener", id)
}

func Bucket(id string) *rBuilder {
	return newResource("bucket", id)
}

func Zone(id string) *rBuilder {
	return newResource("zone", id)
}

func Record(id string) *rBuilder {
	return newResource("record", id)
}

func ScalingGroup(id string) *rBuilder {
	return newResource("scalinggroup", id)
}

func LaunchConfig(id string) *rBuilder {
	return newResource("launchconfiguration", id)
}

func Subscription(id string) *rBuilder {
	return newResource("subscription", id)
}

func Topic(id string) *rBuilder {
	return newResource("topic", id)
}

func Queue(id string) *rBuilder {
	return newResource("queue", id)
}

func Function(id string) *rBuilder {
	return newResource("function", id)
}

func Alarm(id string) *rBuilder {
	return newResource("alarm", id)
}

func Metric(id string) *rBuilder {
	return newResource("metric", id)
}

func Image(id string) *rBuilder {
	return newResource("image", id)
}

func Distribution(id string) *rBuilder {
	return newResource("distribution", id)
}

func Stack(id string) *rBuilder {
	return newResource("stack", id)
}

func Repository(id string) *rBuilder {
	return newResource("repository", id)
}

func ContainerCluster(id string) *rBuilder {
	return newResource("containercluster", id)
}

func ContainerTask(id string) *rBuilder {
	return newResource("containertask", id)
}

func Container(id string) *rBuilder {
	return newResource("container", id)
}

func ContainerInstance(id string) *rBuilder {
	return newResource("containerinstance", id)
}

func NetworkInterface(id string) *rBuilder {
	return newResource("networkinterface", id)
}

func Certificate(id string) *rBuilder {
	return newResource("certificate", id)
}

func AccessKey(id string) *rBuilder {
	return newResource("accesskey", id)
}

func (b *rBuilder) Prop(key string, value any) *rBuilder {
	b.props[key] = value
	return b
}

func (b *rBuilder) Build() *graph.Resource {
	res := graph.InitResource(b.typ, b.id)
	for k, v := range b.props {
		res.Properties()[k] = v
	}

	return res
}

func AddParents(g *graph.Graph, relations ...string) {
	for _, rel := range relations {
		splits := strings.Split(rel, "->")
		if len(splits) != 2 {
			panic(fmt.Sprintf("invalid relation '%s'", rel))
		}
		r1 := graph.InitResource("", strings.TrimSpace(splits[0]))
		r2 := graph.InitResource("", strings.TrimSpace(splits[1]))
		err := g.AddParentRelation(r1, r2)
		if err != nil {
			panic(err)
		}
	}
}

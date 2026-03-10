package cloud

import (
	"context"
	"testing"
)

func TestPluralizeResourceExtra(t *testing.T) {
	tcases := []struct {
		in, out string
	}{
		{in: "image", out: "images"},
		{in: "securitygroup", out: "securitygroups"},
		{in: "keypair", out: "keypairs"},
		{in: "volume", out: "volumes"},
		{in: "snapshot", out: "snapshots"},
		{in: "certificate", out: "certificates"},
		{in: "loadbalancer", out: "loadbalancers"},
		{in: "targetgroup", out: "targetgroups"},
		{in: "listener", out: "listeners"},
		{in: "database", out: "databases"},
		{in: "bucket", out: "buckets"},
		{in: "s3object", out: "s3objects"},
		{in: "subscription", out: "subscriptions"},
		{in: "topic", out: "topics"},
		{in: "queue", out: "queues"},
		{in: "zone", out: "zones"},
		{in: "record", out: "records"},
		{in: "function", out: "functions"},
		{in: "metric", out: "metrics"},
		{in: "alarm", out: "alarms"},
		{in: "distribution", out: "distributions"},
		{in: "stack", out: "stacks"},
		{in: "container", out: "containers"},
		// -cy/-ry suffix
		{in: "policy", out: "policies"},
		{in: "repository", out: "repositories"},
		{in: "registry", out: "registries"},
	}
	for _, tc := range tcases {
		if got, want := PluralizeResource(tc.in), tc.out; got != want {
			t.Fatalf("PluralizeResource(%q): got %s, want %s", tc.in, got, want)
		}
	}
}

func TestSingularizeResourceExtra(t *testing.T) {
	tcases := []struct {
		in, out string
	}{
		{in: "images", out: "image"},
		{in: "securitygroups", out: "securitygroup"},
		{in: "keypairs", out: "keypair"},
		{in: "volumes", out: "volume"},
		{in: "snapshots", out: "snapshot"},
		{in: "certificates", out: "certificate"},
		{in: "loadbalancers", out: "loadbalancer"},
		{in: "databases", out: "database"},
		{in: "buckets", out: "bucket"},
		{in: "subscriptions", out: "subscription"},
		{in: "queues", out: "queue"},
		{in: "functions", out: "function"},
		{in: "stacks", out: "stack"},
		{in: "containers", out: "container"},
		// -ies suffix
		{in: "policies", out: "policy"},
		{in: "repositories", out: "repository"},
		{in: "registries", out: "registry"},
		// already singular (no trailing s)
		{in: "instance", out: "instance"},
	}
	for _, tc := range tcases {
		if got, want := SingularizeResource(tc.in), tc.out; got != want {
			t.Fatalf("SingularizeResource(%q): got %s, want %s", tc.in, got, want)
		}
	}
}

func TestPluralSingularRoundtrip(t *testing.T) {
	resources := []string{
		Instance, Vpc, Subnet, Image, SecurityGroup, Keypair, Volume,
		Snapshot, User, Role, Group, Policy, Bucket, Topic, Queue,
		Function, Alarm, Stack, Container, Repository, Registry,
		Database, Distribution, Certificate, LoadBalancer,
	}
	for _, r := range resources {
		plural := PluralizeResource(r)
		back := SingularizeResource(plural)
		if back != r {
			t.Fatalf("roundtrip failed for %q: pluralized to %q, singularized back to %q", r, plural, back)
		}
	}
}

func TestResourceConstants(t *testing.T) {
	// Verify key resource type constants have expected values
	expected := map[string]string{
		"Instance":        "instance",
		"Vpc":             "vpc",
		"Subnet":          "subnet",
		"SecurityGroup":   "securitygroup",
		"Policy":          "policy",
		"User":            "user",
		"Role":            "role",
		"Group":           "group",
		"Bucket":          "bucket",
		"Database":        "database",
		"LoadBalancer":    "loadbalancer",
		"Function":        "function",
		"Stack":           "stack",
		"Queue":           "queue",
		"Region":          "region",
		"Image":           "image",
		"Volume":          "volume",
		"Keypair":         "keypair",
		"InternetGateway": "internetgateway",
		"NatGateway":      "natgateway",
		"RouteTable":      "routetable",
		"ElasticIP":       "elasticip",
		"Snapshot":        "snapshot",
		"Certificate":     "certificate",
		"Topic":           "topic",
		"Subscription":    "subscription",
		"Zone":            "zone",
		"Record":          "record",
		"Distribution":    "distribution",
		"Alarm":           "alarm",
		"Metric":          "metric",
	}

	actuals := map[string]string{
		"Instance":        Instance,
		"Vpc":             Vpc,
		"Subnet":          Subnet,
		"SecurityGroup":   SecurityGroup,
		"Policy":          Policy,
		"User":            User,
		"Role":            Role,
		"Group":           Group,
		"Bucket":          Bucket,
		"Database":        Database,
		"LoadBalancer":    LoadBalancer,
		"Function":        Function,
		"Stack":           Stack,
		"Queue":           Queue,
		"Region":          Region,
		"Image":           Image,
		"Volume":          Volume,
		"Keypair":         Keypair,
		"InternetGateway": InternetGateway,
		"NatGateway":      NatGateway,
		"RouteTable":      RouteTable,
		"ElasticIP":       ElasticIP,
		"Snapshot":        Snapshot,
		"Certificate":     Certificate,
		"Topic":           Topic,
		"Subscription":    Subscription,
		"Zone":            Zone,
		"Record":          Record,
		"Distribution":    Distribution,
		"Alarm":           Alarm,
		"Metric":          Metric,
	}

	for name, expVal := range expected {
		if actual, ok := actuals[name]; !ok {
			t.Fatalf("missing constant %s", name)
		} else if actual != expVal {
			t.Fatalf("constant %s: got %s, want %s", name, actual, expVal)
		}
	}
}

func TestNewQuery(t *testing.T) {
	t.Run("single resource type", func(t *testing.T) {
		q := NewQuery(Instance)
		if len(q.ResourceType) != 1 || q.ResourceType[0] != Instance {
			t.Fatalf("expected [instance], got %v", q.ResourceType)
		}
		if q.Matcher != nil {
			t.Fatal("expected nil matcher")
		}
	})

	t.Run("multiple resource types", func(t *testing.T) {
		q := NewQuery(Instance, Vpc, Subnet)
		if len(q.ResourceType) != 3 {
			t.Fatalf("expected 3 resource types, got %d", len(q.ResourceType))
		}
	})

	t.Run("empty query", func(t *testing.T) {
		q := NewQuery()
		if len(q.ResourceType) != 0 {
			t.Fatalf("expected empty resource types, got %v", q.ResourceType)
		}
	})
}

type stubMatcher struct {
	matchResult bool
}

func (m *stubMatcher) Match(r Resource) bool {
	return m.matchResult
}

func TestQueryMatch(t *testing.T) {
	q := NewQuery(Instance)
	m := &stubMatcher{matchResult: true}
	q2 := q.Match(m)
	if q2.Matcher == nil {
		t.Fatal("expected matcher to be set")
	}
	// Original query should not be modified (value receiver)
	if q.Matcher != nil {
		t.Fatal("original query should not be modified")
	}
}

func TestServicesNames(t *testing.T) {
	srvs := Services{&stubService{name: "infra"}, &stubService{name: "access"}, &stubService{name: "storage"}}
	names := srvs.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	expected := []string{"infra", "access", "storage"}
	for i, n := range expected {
		if names[i] != n {
			t.Fatalf("names[%d]: got %s, want %s", i, names[i], n)
		}
	}
}

func TestResourcesMap(t *testing.T) {
	res := Resources{&stubResource{id: "i-123", typ: "instance"}, &stubResource{id: "i-456", typ: "instance"}}
	ids := res.Map(func(r Resource) string { return r.Id() })
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
	if ids[0] != "i-123" || ids[1] != "i-456" {
		t.Fatalf("got %v, want [i-123, i-456]", ids)
	}
}

// --- stubs ---

type stubService struct {
	name string
}

func (s *stubService) Region() string                            { return "us-east-1" }
func (s *stubService) Profile() string                           { return "default" }
func (s *stubService) Name() string                              { return s.name }
func (s *stubService) ResourceTypes() []string                   { return nil }
func (s *stubService) IsSyncDisabled() bool                      { return false }
func (s *stubService) Fetch(_ context.Context) (GraphAPI, error) { return nil, nil }
func (s *stubService) FetchByType(_ context.Context, _ string) (GraphAPI, error) {
	return nil, nil
}

type stubResource struct {
	id, typ string
}

func (r *stubResource) Type() string                        { return r.typ }
func (r *stubResource) Id() string                          { return r.id }
func (r *stubResource) String() string                      { return r.id }
func (r *stubResource) Format(string) string                { return r.id }
func (r *stubResource) Properties() map[string]interface{}  { return nil }
func (r *stubResource) Property(string) (interface{}, bool) { return nil, false }
func (r *stubResource) Meta(string) (interface{}, bool)     { return nil, false }
func (r *stubResource) Same(Resource) bool                  { return false }

package graph

import (
	"reflect"
	"testing"
)

func TestInitResource(t *testing.T) {
	tcases := []struct {
		kind, id string
	}{
		{"instance", "i-12345"},
		{"subnet", "sub-abc"},
		{"vpc", "vpc-xyz"},
		{"", "empty-type"},
		{"securitygroup", ""},
	}
	for _, tc := range tcases {
		r := InitResource(tc.kind, tc.id)
		if got, want := r.Type(), tc.kind; got != want {
			t.Fatalf("Type(): got %q, want %q", got, want)
		}
		if got, want := r.Id(), tc.id; got != want {
			t.Fatalf("Id(): got %q, want %q", got, want)
		}
		if r.Properties() == nil {
			t.Fatal("expected non-nil properties map")
		}
		if got, want := r.Properties()["ID"], tc.id; got != want {
			t.Fatalf("Properties()[ID]: got %q, want %q", got, want)
		}
	}
}

func TestNotFoundResource(t *testing.T) {
	r := NotFoundResource("missing-abc")
	if got, want := r.Id(), "missing-abc"; got != want {
		t.Fatalf("Id(): got %q, want %q", got, want)
	}
	if got, want := r.Type(), notFoundResourceType; got != want {
		t.Fatalf("Type(): got %q, want %q", got, want)
	}
	if r.Properties() == nil {
		t.Fatal("expected non-nil properties map")
	}
	// NotFoundResource should have the "<not-found>" format string
	if got, want := r.Format("%t"), "<not-found>"; got != want {
		t.Fatalf("Format: got %q, want %q", got, want)
	}
}

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if g == nil {
		t.Fatal("expected non-nil graph")
	}
	// A new graph should have no resources of any type
	resources, err := g.GetAllResources("instance")
	if err != nil {
		t.Fatal(err)
	}
	if len(resources) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(resources))
	}
	snap := g.AsRDFGraphSnaphot()
	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
}

func TestParsePortRange(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tcases := []struct {
			input  string
			expect PortRange
		}{
			{"80:80", PortRange{FromPort: 80, ToPort: 80}},
			{"80:443", PortRange{FromPort: 80, ToPort: 443}},
			{"0:65535", PortRange{FromPort: 0, ToPort: 65535}},
			{":", PortRange{Any: true}},
			{"1:1024", PortRange{FromPort: 1, ToPort: 1024}},
		}
		for _, tc := range tcases {
			got, err := ParsePortRange(tc.input)
			if err != nil {
				t.Fatalf("ParsePortRange(%q): unexpected error: %s", tc.input, err)
			}
			if got != tc.expect {
				t.Fatalf("ParsePortRange(%q): got %+v, want %+v", tc.input, got, tc.expect)
			}
		}
	})
	t.Run("invalid", func(t *testing.T) {
		invalids := []string{
			"abc",
			"abc:def",
			"80",
			"80:abc",
			"abc:80",
			"",
		}
		for _, input := range invalids {
			_, err := ParsePortRange(input)
			if err == nil {
				t.Fatalf("ParsePortRange(%q): expected error, got nil", input)
			}
		}
	})
}

func TestPortRangeString(t *testing.T) {
	tcases := []struct {
		pr     PortRange
		expect string
	}{
		{PortRange{Any: true}, ":"},
		{PortRange{FromPort: 80, ToPort: 80}, "80:80"},
		{PortRange{FromPort: 80, ToPort: 443}, "80:443"},
		{PortRange{FromPort: -1, ToPort: 22}, "22:22"},
		{PortRange{FromPort: 22, ToPort: -1}, "22:22"},
	}
	for _, tc := range tcases {
		if got := tc.pr.String(); got != tc.expect {
			t.Fatalf("PortRange(%+v).String(): got %q, want %q", tc.pr, got, tc.expect)
		}
	}
}

func TestParseRouteTarget(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tcases := []struct {
			input      string
			expectType routeTargetType
			expectRef  string
			expectOwn  string
		}{
			{"0|igw-123|owner1", EgressOnlyInternetGatewayTarget, "igw-123", "owner1"},
			{"1|gw-456|me", GatewayTarget, "gw-456", "me"},
			{"2|inst-789|", InstanceTarget, "inst-789", ""},
			{"3|nat-abc|acc", NatTarget, "nat-abc", "acc"},
			{"4|eni-def|owner", NetworkInterfaceTarget, "eni-def", "owner"},
			{"5|pcx-ghi|peer", VpcPeeringConnectionTarget, "pcx-ghi", "peer"},
		}
		for _, tc := range tcases {
			got, err := ParseRouteTarget(tc.input)
			if err != nil {
				t.Fatalf("ParseRouteTarget(%q): unexpected error: %s", tc.input, err)
			}
			if got.Type != tc.expectType {
				t.Fatalf("ParseRouteTarget(%q): Type got %d, want %d", tc.input, got.Type, tc.expectType)
			}
			if got.Ref != tc.expectRef {
				t.Fatalf("ParseRouteTarget(%q): Ref got %q, want %q", tc.input, got.Ref, tc.expectRef)
			}
			if got.Owner != tc.expectOwn {
				t.Fatalf("ParseRouteTarget(%q): Owner got %q, want %q", tc.input, got.Owner, tc.expectOwn)
			}
		}
	})
	t.Run("invalid", func(t *testing.T) {
		invalids := []string{
			"",
			"foo",
			"a|b",
			"notanumber|ref|owner",
			"1|b|c|d",
		}
		for _, input := range invalids {
			_, err := ParseRouteTarget(input)
			if err == nil {
				t.Fatalf("ParseRouteTarget(%q): expected error, got nil", input)
			}
		}
	})
}

func TestRouteTargetString(t *testing.T) {
	rt := &RouteTarget{Type: GatewayTarget, Ref: "gw-123", Owner: "me"}
	if got, want := rt.String(), "1|gw-123|me"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestParsePortRangeRoundtrip(t *testing.T) {
	originals := []PortRange{
		{Any: true},
		{FromPort: 22, ToPort: 22},
		{FromPort: 80, ToPort: 443},
		{FromPort: 0, ToPort: 65535},
	}
	for _, orig := range originals {
		s := orig.String()
		parsed, err := ParsePortRange(s)
		if err != nil {
			t.Fatalf("roundtrip error for %+v (string %q): %s", orig, s, err)
		}
		if parsed != orig {
			t.Fatalf("roundtrip mismatch for %+v: string=%q parsed=%+v", orig, s, parsed)
		}
	}
}

func TestParseRouteTargetRoundtrip(t *testing.T) {
	originals := []*RouteTarget{
		{Type: GatewayTarget, Ref: "gw-1", Owner: "me"},
		{Type: InstanceTarget, Ref: "i-abc", Owner: ""},
		{Type: NatTarget, Ref: "nat-xyz", Owner: "account"},
	}
	for _, orig := range originals {
		s := orig.String()
		parsed, err := ParseRouteTarget(s)
		if err != nil {
			t.Fatalf("roundtrip error for %+v: %s", orig, err)
		}
		if *parsed != *orig {
			t.Fatalf("roundtrip mismatch: got %+v, want %+v", parsed, orig)
		}
	}
}

func TestGraphMerge(t *testing.T) {
	g1 := NewGraph()
	g1.AddResource(instResource("i1").build())

	g2 := NewGraph()
	g2.AddResource(subResource("s1").build())

	if err := g1.Merge(g2); err != nil {
		t.Fatal(err)
	}

	instances, _ := g1.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d instances, want %d", got, want)
	}
	subnets, _ := g1.GetAllResources("subnet")
	if got, want := len(subnets), 1; got != want {
		t.Fatalf("got %d subnets, want %d", got, want)
	}
}

func TestOrFilter(t *testing.T) {
	g := NewGraph()
	i1 := instResource("i1").prop("Name", "redis").build()
	i2 := instResource("i2").prop("Name", "postgres").build()
	i3 := instResource("i3").prop("Name", "memcache").build()
	g.AddResource(i1, i2, i3)

	matchRedis := func(r *Resource) bool { return r.Properties()["Name"] == "redis" }
	matchPostgres := func(r *Resource) bool { return r.Properties()["Name"] == "postgres" }

	filtered, err := g.OrFilter("instance", matchRedis, matchPostgres)
	if err != nil {
		t.Fatal(err)
	}
	results, _ := filtered.GetAllResources("instance")
	if got, want := len(results), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	// OrFilter with no filters should return all
	filtered, err = g.OrFilter("instance")
	if err != nil {
		t.Fatal(err)
	}
	results, _ = filtered.GetAllResources("instance")
	if got, want := len(results), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestResolveResourcesWithProp(t *testing.T) {
	g := NewGraph()
	i1 := instResource("i1").prop("Name", "redis").build()
	i2 := instResource("i2").prop("Name", "postgres").build()
	s1 := subResource("s1").prop("Name", "redis").build()
	g.AddResource(i1, i2, s1)

	snap := g.AsRDFGraphSnaphot()

	// Resolve instances with Name=redis
	results, err := ResolveResourcesWithProp(snap, "instance", "Name", "redis")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(results), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := results[0].Id(), "i1"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	// Resolve subnets with Name=redis
	results, err = ResolveResourcesWithProp(snap, "subnet", "Name", "redis")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(results), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := results[0].Id(), "s1"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	// Resolve instances with Name=nonexistent
	results, err = ResolveResourcesWithProp(snap, "instance", "Name", "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSubtractEmptyAndIdentical(t *testing.T) {
	// Subtract from empty
	result := Subtract(map[string]interface{}{}, map[string]interface{}{"a": 1})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}

	// Subtract empty from non-empty
	one := map[string]interface{}{"a": 1, "b": "two"}
	result = Subtract(one, map[string]interface{}{})
	if got, want := result, one; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	// Subtract identical maps
	m := map[string]interface{}{"x": 10, "y": "hello"}
	result = Subtract(m, m)
	if len(result) != 0 {
		t.Fatalf("expected empty result for identical maps, got %v", result)
	}

	// Both empty
	result = Subtract(map[string]interface{}{}, map[string]interface{}{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestSubtractWithSliceValues(t *testing.T) {
	one := map[string]interface{}{
		"tags": []string{"a", "b"},
		"id":   "123",
	}
	other := map[string]interface{}{
		"tags": []string{"a", "b"},
		"id":   "456",
	}
	result := Subtract(one, other)
	// tags are equal (DeepEqual), so only id should remain
	expected := map[string]interface{}{"id": "123"}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("got %v, want %v", result, expected)
	}
}

func TestAddResourceRelation(t *testing.T) {
	r := InitResource("instance", "i-1")
	child := InitResource("volume", "vol-1")
	r.AddRelation("children", child)

	// Verify the resource tracks the relation
	if got, want := len(r.relations["children"]), 1; got != want {
		t.Fatalf("got %d relations, want %d", got, want)
	}
	if got, want := r.relations["children"][0].Id(), "vol-1"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResourcePropertyAndMeta(t *testing.T) {
	r := InitResource("instance", "i-1")
	r.SetProperty("Name", "my-instance")

	v, ok := r.Property("Name")
	if !ok {
		t.Fatal("expected to find property Name")
	}
	if got, want := v, "my-instance"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	_, ok = r.Property("NonExistent")
	if ok {
		t.Fatal("expected property not found")
	}

	_, ok = r.Meta("anything")
	if ok {
		t.Fatal("expected meta not found on fresh resource")
	}
}

func TestResourceSame(t *testing.T) {
	r1 := InitResource("instance", "i-1")
	r1.SetProperty("Name", "foo")

	r2 := InitResource("instance", "i-1")
	r2.SetProperty("Name", "bar")

	// Same compares only id and type, not properties
	if !r1.Same(r2) {
		t.Fatal("expected Same() to return true for same id+type")
	}

	r3 := InitResource("subnet", "i-1")
	if r1.Same(r3) {
		t.Fatal("expected Same() to return false for different type")
	}
}

func TestMarshalMustMarshalRoundtrip(t *testing.T) {
	g := NewGraph()
	g.AddResource(instResource("i1").prop("Name", "test").build())

	data := g.MustMarshal()
	if data == "" {
		t.Fatal("expected non-empty marshal output")
	}

	g2 := NewGraph()
	if err := g2.Unmarshal([]byte(data)); err != nil {
		t.Fatal(err)
	}

	res, err := g2.GetResource("instance", "i1")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.Properties()["Name"], "test"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

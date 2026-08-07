package graph

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	tstore "github.com/bootswithdefer/triplestore"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/rdf"
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
	if err := g1.AddResource(instResource("i1").build()); err != nil {
		t.Fatal(err)
	}

	g2 := NewGraph()
	if err := g2.AddResource(subResource("s1").build()); err != nil {
		t.Fatal(err)
	}

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
	if err := g.AddResource(i1, i2, i3); err != nil {
		t.Fatal(err)
	}

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
	if err := g.AddResource(i1, i2, s1); err != nil {
		t.Fatal(err)
	}

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
	result := Subtract(map[string]any{}, map[string]any{"a": 1})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}

	// Subtract empty from non-empty
	one := map[string]any{"a": 1, "b": "two"}
	result = Subtract(one, map[string]any{})
	if got, want := result, one; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	// Subtract identical maps
	m := map[string]any{"x": 10, "y": "hello"}
	result = Subtract(m, m)
	if len(result) != 0 {
		t.Fatalf("expected empty result for identical maps, got %v", result)
	}

	// Both empty
	result = Subtract(map[string]any{}, map[string]any{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestSubtractWithSliceValues(t *testing.T) {
	one := map[string]any{
		"tags": []string{"a", "b"},
		"id":   "123",
	}
	other := map[string]any{
		"tags": []string{"a", "b"},
		"id":   "456",
	}
	result := Subtract(one, other)
	// tags are equal (DeepEqual), so only id should remain
	expected := map[string]any{"id": "123"}
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
	if err := g.AddResource(instResource("i1").prop("Name", "test").build()); err != nil {
		t.Fatal(err)
	}

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

func TestMarshalTo(t *testing.T) {
	g := NewGraph()
	if err := g.AddResource(instResource("i1").prop("Name", "test-marshal").build()); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := g.MarshalTo(&buf); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output from MarshalTo")
	}

	// Verify we can unmarshal what was written
	g2 := NewGraph()
	if err := g2.Unmarshal(buf.Bytes()); err != nil {
		t.Fatal(err)
	}
	res, err := g2.GetResource("instance", "i1")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := res.Properties()["Name"], "test-marshal"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestVisitRelations(t *testing.T) {
	g := NewGraph()
	v1 := InitResource("vpc", "vpc_1")
	s1 := InitResource("subnet", "sub_1")
	s2 := InitResource("subnet", "sub_2")
	i1 := InitResource("instance", "inst_1")
	sg1 := InitResource("securitygroup", "secgroup_1")
	if err := g.AddResource(v1, s1, s2, i1, sg1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddParentRelation(v1, s1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddParentRelation(v1, s2); err != nil {
		t.Fatal(err)
	}
	if err := g.AddParentRelation(s1, i1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddAppliesOnRelation(sg1, i1); err != nil {
		t.Fatal(err)
	}

	t.Run("ChildrenOfRel", func(t *testing.T) {
		var collected []string
		err := g.VisitRelations(v1, rdf.ChildrenOfRel, false, func(r cloud.Resource, depth int) error {
			collected = append(collected, r.Id())
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(collected) == 0 {
			t.Fatal("expected children to be visited")
		}
	})

	t.Run("ChildrenOfRel with IncludeFrom", func(t *testing.T) {
		var collected []string
		err := g.VisitRelations(v1, rdf.ChildrenOfRel, true, func(r cloud.Resource, depth int) error {
			collected = append(collected, r.Id())
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		foundRoot := false
		for _, id := range collected {
			if id == "vpc_1" {
				foundRoot = true
			}
		}
		if !foundRoot {
			t.Fatal("expected root to be included when IncludeFrom=true")
		}
	})

	t.Run("DependingOnRel", func(t *testing.T) {
		var collected []string
		err := g.VisitRelations(i1, rdf.DependingOnRel, false, func(r cloud.Resource, depth int) error {
			collected = append(collected, r.Id())
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(collected) != 1 || collected[0] != "secgroup_1" {
			t.Fatalf("expected [secgroup_1], got %v", collected)
		}
	})

	t.Run("ApplyOn", func(t *testing.T) {
		var collected []string
		err := g.VisitRelations(sg1, rdf.ApplyOn, false, func(r cloud.Resource, depth int) error {
			collected = append(collected, r.Id())
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(collected) != 1 || collected[0] != "inst_1" {
			t.Fatalf("expected [inst_1], got %v", collected)
		}
	})

	t.Run("default relation (ParentOf)", func(t *testing.T) {
		var collected []string
		err := g.VisitRelations(i1, rdf.ParentOf, false, func(r cloud.Resource, depth int) error {
			collected = append(collected, r.Id())
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(collected) == 0 {
			t.Fatal("expected parents to be visited")
		}
	})
}

func TestListResourcesDependingOn(t *testing.T) {
	g := NewGraph()
	inst := InitResource("instance", "inst_1")
	sg1 := InitResource("securitygroup", "sg_1")
	sg2 := InitResource("securitygroup", "sg_2")
	if err := g.AddResource(inst, sg1, sg2); err != nil {
		t.Fatal(err)
	}
	if err := g.AddAppliesOnRelation(sg1, inst); err != nil {
		t.Fatal(err)
	}
	if err := g.AddAppliesOnRelation(sg2, inst); err != nil {
		t.Fatal(err)
	}

	resources, err := g.ListResourcesDependingOn(inst)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	ids := map[string]bool{}
	for _, r := range resources {
		ids[r.Id()] = true
	}
	if !ids["sg_1"] || !ids["sg_2"] {
		t.Fatalf("expected sg_1 and sg_2 in results, got %v", ids)
	}
}

func TestListResourcesAppliedOn(t *testing.T) {
	g := NewGraph()
	sg := InitResource("securitygroup", "sg_1")
	inst1 := InitResource("instance", "inst_1")
	inst2 := InitResource("instance", "inst_2")
	if err := g.AddResource(sg, inst1, inst2); err != nil {
		t.Fatal(err)
	}
	if err := g.AddAppliesOnRelation(sg, inst1); err != nil {
		t.Fatal(err)
	}
	if err := g.AddAppliesOnRelation(sg, inst2); err != nil {
		t.Fatal(err)
	}

	resources, err := g.ListResourcesAppliedOn(sg)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	ids := map[string]bool{}
	for _, r := range resources {
		ids[r.Id()] = true
	}
	if !ids["inst_1"] || !ids["inst_2"] {
		t.Fatalf("expected inst_1 and inst_2 in results, got %v", ids)
	}
}

func TestMergeWithNonGraphType(t *testing.T) {
	g := NewGraph()
	err := g.Merge(nil)
	if err == nil {
		t.Fatal("expected error when merging nil")
	}
}

func TestFindWithNoResourceType(t *testing.T) {
	g := NewGraph()
	if err := g.AddResource(instResource("i1").build()); err != nil {
		t.Fatal(err)
	}

	_, err := g.Find(cloud.Query{})
	if err == nil {
		t.Fatal("expected error for query with no resource type")
	}
	if !strings.Contains(err.Error(), "at least one resource type") {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestFindWithMultipleTypesAndMatcher(t *testing.T) {
	g := NewGraph()
	if err := g.AddResource(instResource("i1").build()); err != nil {
		t.Fatal(err)
	}

	// Multiple resource types + matcher should error
	q := cloud.Query{ResourceType: []string{"instance", "subnet"}}
	// First verify no matcher works fine
	res, err := g.Find(q)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1, got %d", len(res))
	}
}

func TestFilterGraphWithInvalidQuery(t *testing.T) {
	g := NewGraph()
	// Query with no resource types
	_, err := g.FilterGraph(cloud.Query{})
	if err == nil {
		t.Fatal("expected error for empty resource type")
	}

	// Query with multiple resource types
	_, err = g.FilterGraph(cloud.Query{ResourceType: []string{"instance", "subnet"}})
	if err == nil {
		t.Fatal("expected error for multiple resource types")
	}
}

func TestAddResourceWithChildrenRelation(t *testing.T) {
	g := NewGraph()
	parent := InitResource("vpc", "vpc_1")
	child := InitResource("subnet", "sub_1")
	parent.AddRelation(rdf.ChildrenOfRel, child)
	err := g.AddResource(parent, child)
	if err != nil {
		t.Fatal(err)
	}
	// The child should be added and the parent relation created
	snap := g.store.Snapshot()
	triples := snap.WithSubjPred(child.Id(), rdf.ParentOf)
	if len(triples) != 1 {
		t.Fatalf("expected 1 parent relation, got %d", len(triples))
	}
}

func TestAddResourceWithDependingOnRelation(t *testing.T) {
	g := NewGraph()
	sg := InitResource("securitygroup", "sg_1")
	inst := InitResource("instance", "inst_1")
	sg.AddRelation(rdf.DependingOnRel, inst)
	err := g.AddResource(sg, inst)
	if err != nil {
		t.Fatal(err)
	}
	snap := g.store.Snapshot()
	triples := snap.WithSubjPred(inst.Id(), rdf.ApplyOn)
	if len(triples) != 1 {
		t.Fatalf("expected 1 apply-on relation, got %d", len(triples))
	}
}

func TestFindResourceMultipleWithSameId(t *testing.T) {
	g := NewGraph()
	// Add two instances with the same ID property value
	i1 := instResource("inst_a").prop("Name", "same").build()
	i2 := instResource("inst_b").prop("Name", "same").build()
	if err := g.AddResource(i1, i2); err != nil {
		t.Fatal(err)
	}

	// FindResourcesByProperty with Name should find both
	res, err := g.FindResourcesByProperty("Name", "same")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestTrimNS(t *testing.T) {
	tcases := []struct {
		input, expect string
	}{
		{"cloud-owl:Instance", "Instance"},
		{"Instance", "Instance"},
		{"a:b:c", "c"},
		{"", ""},
	}
	for _, tc := range tcases {
		if got := trimNS(tc.input); got != tc.expect {
			t.Fatalf("trimNS(%q): got %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestLowerFirstLetter(t *testing.T) {
	tcases := []struct {
		input, expect string
	}{
		{"Instance", "instance"},
		{"Subnet", "subnet"},
		{"a", "a"},
		{"ABC", "aBC"},
	}
	for _, tc := range tcases {
		if got := lowerFirstLetter(tc.input); got != tc.expect {
			t.Fatalf("lowerFirstLetter(%q): got %q, want %q", tc.input, got, tc.expect)
		}
	}
}

func TestKeyValueString(t *testing.T) {
	kv := &KeyValue{KeyName: "mykey", Value: "myval"}
	expected := "[Key:mykey,Value:myval]"
	if got := kv.String(); got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestDistributionOriginString(t *testing.T) {
	tcases := []struct {
		origin *DistributionOrigin
		expect string
	}{
		{
			&DistributionOrigin{ID: "origin1"},
			"[ID:origin1]",
		},
		{
			&DistributionOrigin{ID: "origin2", PublicDNS: "example.com", PathPrefix: "/path", OriginType: "s3", Config: "cfg"},
			"[ID:origin2,PublicDNS:example.com,PathPrefix:/path,Type:s3,Config:cfg]",
		},
		{
			&DistributionOrigin{ID: "origin3", PublicDNS: "dns.com"},
			"[ID:origin3,PublicDNS:dns.com]",
		},
	}
	for i, tc := range tcases {
		if got := tc.origin.String(); got != tc.expect {
			t.Fatalf("%d: got %q, want %q", i, got, tc.expect)
		}
	}
}

func TestFirewallRuleString(t *testing.T) {
	rule := &FirewallRule{
		PortRange: PortRange{FromPort: 80, ToPort: 80},
		Protocol:  "tcp",
	}
	s := rule.String()
	if !strings.Contains(s, "Protocol:tcp") {
		t.Fatalf("expected Protocol:tcp in %q", s)
	}
	if !strings.Contains(s, "PortRange:") {
		t.Fatalf("expected PortRange: in %q", s)
	}
}

func TestRouteString(t *testing.T) {
	route := &Route{
		DestinationPrefixListId: "pl-123",
	}
	s := route.String()
	if !strings.Contains(s, "DestinationPrefixListId:pl-123") {
		t.Fatalf("expected DestinationPrefixListId in %q", s)
	}
}

func TestGrantString(t *testing.T) {
	grant := &Grant{
		Permission: "FULL_CONTROL",
		Grantee: Grantee{
			GranteeID:          "user123",
			GranteeDisplayName: "User",
			GranteeType:        "CanonicalUser",
		},
	}
	s := grant.String()
	if !strings.Contains(s, "Permission:FULL_CONTROL") {
		t.Fatalf("expected Permission in %q", s)
	}
	if !strings.Contains(s, "GranteeID:user123") {
		t.Fatalf("expected GranteeID in %q", s)
	}
}

func TestResourceFormatInvalidVerb(t *testing.T) {
	r := &Resource{id: "test", kind: "instance"}
	out := r.Format("%z")
	if !strings.Contains(out, "invalid verb") {
		t.Fatalf("expected invalid verb error, got %q", out)
	}
}

func TestResourceSameNilCases(t *testing.T) {
	// Both nil
	var r1 *Resource
	if !r1.Same(nil) {
		t.Fatal("expected nil.Same(nil) == true")
	}

	// One nil, one not
	r2 := InitResource("instance", "i-1")
	if r2.Same(nil) {
		t.Fatal("expected non-nil.Same(nil) == false")
	}
}

func TestMarshalUnmarshalKeyValues(t *testing.T) {
	r := testResource("dist1", "distribution").prop("ID", "dist1").prop(
		"Dimensions", []*KeyValue{
			{KeyName: "key1", Value: "value1"},
			{KeyName: "key2", Value: "value2"},
		}).build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)
	rawRes := InitResource(r.Type(), r.Id())
	err = rawRes.unmarshalFullRdf(g.store.Snapshot())
	if err != nil {
		t.Fatal(err)
	}
	kvs, ok := rawRes.Properties()["Dimensions"].([]*KeyValue)
	if !ok {
		t.Fatalf("expected []*KeyValue, got %T", rawRes.Properties()["Dimensions"])
	}
	if got, want := len(kvs), 2; got != want {
		t.Fatalf("got %d key-values, want %d", got, want)
	}
}

func TestMarshalUnmarshalDistributionOrigins(t *testing.T) {
	r := testResource("dist1", "distribution").prop("ID", "dist1").prop(
		"Origins", []*DistributionOrigin{
			{ID: "o1", PublicDNS: "example.com", PathPrefix: "/path", OriginType: "s3", Config: "cfg1"},
			{ID: "o2", PublicDNS: "other.com", PathPrefix: "/", OriginType: "custom", Config: "cfg2"},
		}).build()
	g := NewGraph()
	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	g.store.Add(triples...)
	rawRes := InitResource(r.Type(), r.Id())
	err = rawRes.unmarshalFullRdf(g.store.Snapshot())
	if err != nil {
		t.Fatal(err)
	}
	origins, ok := rawRes.Properties()["Origins"].([]*DistributionOrigin)
	if !ok {
		t.Fatalf("expected []*DistributionOrigin, got %T", rawRes.Properties()["Origins"])
	}
	if got, want := len(origins), 2; got != want {
		t.Fatalf("got %d origins, want %d", got, want)
	}
}

func TestResourcesMap(t *testing.T) {
	res := Resources{
		InitResource("instance", "i1"),
		InitResource("subnet", "s1"),
	}
	ids := res.Map(func(r *Resource) string { return r.Id() })
	if got, want := len(ids), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if ids[0] != "i1" || ids[1] != "s1" {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestFirewallRulesSort(t *testing.T) {
	rules := FirewallRules{
		{PortRange: PortRange{FromPort: 443, ToPort: 443}, Protocol: "tcp"},
		{PortRange: PortRange{FromPort: 22, ToPort: 22}, Protocol: "tcp"},
		{PortRange: PortRange{FromPort: 80, ToPort: 80}, Protocol: "tcp"},
	}
	rules.Sort()
	// After sorting, the order should be by string representation
	for i := 0; i < len(rules)-1; i++ {
		if rules[i].String() > rules[i+1].String() {
			t.Fatalf("rules not sorted: %s > %s", rules[i].String(), rules[i+1].String())
		}
	}
}

func TestRoutesSort(t *testing.T) {
	routes := Routes{
		{DestinationPrefixListId: "pl-zzz"},
		{DestinationPrefixListId: "pl-aaa"},
	}
	routes.Sort()
	if routes[0].DestinationPrefixListId > routes[1].DestinationPrefixListId {
		t.Fatal("routes not sorted properly")
	}
}

func TestGrantsSort(t *testing.T) {
	grants := Grants{
		{Permission: "WRITE"},
		{Permission: "FULL_CONTROL"},
		{Permission: "READ"},
	}
	grants.Sort()
	for i := 0; i < len(grants)-1; i++ {
		if grants[i].String() > grants[i+1].String() {
			t.Fatalf("grants not sorted: %s > %s", grants[i].String(), grants[i+1].String())
		}
	}
}

func TestNamespacedResourceType(t *testing.T) {
	got := namespacedResourceType("instance")
	if !strings.Contains(got, "Instance") {
		t.Fatalf("expected namespaced type to contain 'Instance', got %q", got)
	}
	if !strings.HasPrefix(got, fmt.Sprintf("%s:", rdf.CloudOwlNS)) {
		t.Fatalf("expected prefix %s:, got %q", rdf.CloudOwlNS, got)
	}
}

func TestUnmarshalFromReaders(t *testing.T) {
	g := NewGraph()
	if err := g.AddResource(instResource("i1").prop("Name", "test").build()); err != nil {
		t.Fatal(err)
	}

	// Marshal to a buffer
	var buf bytes.Buffer
	if err := g.MarshalTo(&buf); err != nil {
		t.Fatal(err)
	}

	// Unmarshal from reader
	g2 := NewGraph()
	reader := bytes.NewReader(buf.Bytes())
	if err := g2.UnmarshalFromReaders(reader); err != nil {
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

func TestResourceMarshalWithDiffMeta(t *testing.T) {
	r := instResource("i1").build()
	r.meta["diff"] = "extra"

	triples, err := r.marshalFullRDF()
	if err != nil {
		t.Fatal(err)
	}
	// Should have a meta triple
	foundMeta := false
	for _, tri := range triples {
		if tri.Predicate() == MetaPredicate {
			foundMeta = true
		}
	}
	if !foundMeta {
		t.Fatal("expected a meta triple with diff info")
	}
}

func TestUnmarshalMeta(t *testing.T) {
	g := NewGraph()
	r := instResource("i1").build()
	r.meta["diff"] = "extra"
	if err := g.AddResource(r); err != nil {
		t.Fatal(err)
	}
	// Manually add meta triple
	g.store.Add(
		tstore.SubjPred("i1", MetaPredicate).StringLiteral("extra"),
	)

	res := InitResource("instance", "i1")
	snap := g.store.Snapshot()
	if err := res.unmarshalFullRdf(snap); err != nil {
		t.Fatal(err)
	}
	if err := res.unmarshalMeta(snap); err != nil {
		t.Fatal(err)
	}
	v, ok := res.Meta("diff")
	if !ok {
		t.Fatal("expected diff meta")
	}
	if got, want := v, "extra"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

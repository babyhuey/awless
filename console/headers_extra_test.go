package console

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/fatih/color"

	"github.com/wallix/awless/graph"
)

func TestColoredValueColumnDefinitionFormat(t *testing.T) {
	color.NoColor = true
	col := ColoredValueColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "State"},
		ColoredValues:          map[string]color.Attribute{"running": color.FgGreen, "stopped": color.FgRed},
	}

	if got := col.format("running"); got != "running" {
		t.Fatalf("got %q, want %q", got, "running")
	}
	if got := col.format("stopped"); got != "stopped" {
		t.Fatalf("got %q, want %q", got, "stopped")
	}
	if got := col.format("pending"); got != "pending" {
		t.Fatalf("got %q, want %q", got, "pending")
	}
	if got := col.format(nil); got != "" {
		t.Fatalf("got %q, want %q", got, "")
	}
}

func TestTimeColumnDefinitionFormat(t *testing.T) {
	stamp := time.Date(2017, 6, 15, 10, 30, 0, 0, time.UTC)

	col := TimeColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Created"},
		Format:                 Humanize,
	}
	if got := col.format(nil); got != "" {
		t.Fatalf("Humanize nil: got %q, want empty", got)
	}
	if got := col.format("not-a-time"); got != "invalid time" {
		t.Fatalf("Humanize non-time: got %q, want 'invalid time'", got)
	}

	colBasic := TimeColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Created"},
		Format:                 Basic,
	}
	got := colBasic.format(stamp)
	expected := "Thu, Jun 15, 2017 10:30"
	if got != expected {
		t.Fatalf("Basic format: got %q, want %q", got, expected)
	}

	colShort := TimeColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Created"},
		Format:                 Short,
	}
	got = colShort.format(stamp)
	expected = "6/15/17 10:30"
	if got != expected {
		t.Fatalf("Short format: got %q, want %q", got, expected)
	}
}

func TestStorageColumnDefinitionFormat(t *testing.T) {
	col := StorageColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Size"},
		Unit:                   gb,
	}

	if got := col.format(nil); got != "" {
		t.Fatalf("nil: got %q, want empty", got)
	}
	if got := col.format(int64(2)); got != "2G" {
		t.Fatalf("int64(2) gb: got %q, want '2G'", got)
	}
	if got := col.format(int(5)); got != "5G" {
		t.Fatalf("int(5) gb: got %q, want '5G'", got)
	}
	if got := col.format(uint64(1)); got != "1G" {
		t.Fatalf("uint64(1) gb: got %q, want '1G'", got)
	}
	if got := col.format("string"); got != "invalid size" {
		t.Fatalf("string: got %q, want 'invalid size'", got)
	}
	if got := col.format(true); got != "invalid size" {
		t.Fatalf("bool: got %q, want 'invalid size'", got)
	}

	colBytes := StorageColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Size"},
		Unit:                   b,
	}
	if got := colBytes.format(int64(500)); got != "500B" {
		t.Fatalf("500 bytes: got %q, want '500B'", got)
	}
}

func TestFirewallRulesColumnDefinitionFormat(t *testing.T) {
	col := FirewallRulesColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "InboundRules"},
	}

	if got := col.format(nil); got != "" {
		t.Fatalf("nil: got %q, want empty", got)
	}
	if got := col.format("not-rules"); got != "invalid rules" {
		t.Fatalf("non-rules: got %q, want 'invalid rules'", got)
	}

	_, cidr, _ := net.ParseCIDR("10.0.0.0/8")

	// Protocol "any"
	rules := []*graph.FirewallRule{
		{Protocol: "any", PortRange: graph.PortRange{Any: true}, IPRanges: []*net.IPNet{cidr}},
	}
	got := col.format(rules)
	if got != "[10.0.0.0/8](any) " {
		t.Fatalf("any protocol: got %q", got)
	}

	// Same ports
	rules = []*graph.FirewallRule{
		{Protocol: "tcp", PortRange: graph.PortRange{FromPort: 22, ToPort: 22}, IPRanges: []*net.IPNet{cidr}},
	}
	got = col.format(rules)
	if got != "[10.0.0.0/8](tcp:22) " {
		t.Fatalf("same port: got %q", got)
	}

	// Port range
	rules = []*graph.FirewallRule{
		{Protocol: "tcp", PortRange: graph.PortRange{FromPort: 80, ToPort: 443}, IPRanges: []*net.IPNet{cidr}},
	}
	got = col.format(rules)
	if got != "[10.0.0.0/8](tcp:80-443) " {
		t.Fatalf("port range: got %q", got)
	}

	// Any port range
	rules = []*graph.FirewallRule{
		{Protocol: "tcp", PortRange: graph.PortRange{Any: true}, IPRanges: []*net.IPNet{cidr}},
	}
	got = col.format(rules)
	if got != "[10.0.0.0/8](tcp:any) " {
		t.Fatalf("any port: got %q", got)
	}

	// With sources
	rules = []*graph.FirewallRule{
		{Protocol: "tcp", PortRange: graph.PortRange{FromPort: 22, ToPort: 22}, IPRanges: []*net.IPNet{cidr}, Sources: []string{"sg-123"}},
	}
	got = col.format(rules)
	if got != "[10.0.0.0/8;sg-123](tcp:22) " {
		t.Fatalf("with sources: got %q", got)
	}
}

func TestRoutesColumnDefinitionFormat(t *testing.T) {
	col := RoutesColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Routes"},
	}

	if got := col.format(nil); got != "" {
		t.Fatalf("nil: got %q, want empty", got)
	}
	if got := col.format("not-routes"); got != "invalid routes" {
		t.Fatalf("non-routes: got %q, want 'invalid routes'", got)
	}

	_, cidr, _ := net.ParseCIDR("10.0.0.0/16")
	_, cidr6, _ := net.ParseCIDR("fd00::/64")

	// Basic route with gateway target
	routes := []*graph.Route{
		{
			Destination: cidr,
			Targets:     []*graph.RouteTarget{{Type: graph.GatewayTarget, Ref: "igw-123"}},
		},
	}
	got := col.format(routes)
	if got != "10.0.0.0/16->gw:igw-123 " {
		t.Fatalf("basic route: got %q", got)
	}

	// Route with IPv6
	routes = []*graph.Route{
		{
			Destination:     cidr,
			DestinationIPv6: cidr6,
			Targets:         []*graph.RouteTarget{{Type: graph.NatTarget, Ref: "nat-456"}},
		},
	}
	got = col.format(routes)
	if got != "10.0.0.0/16+fd00::/64->nat:nat-456 " {
		t.Fatalf("ipv6 route: got %q", got)
	}

	// Route with multiple targets
	routes = []*graph.Route{
		{
			Destination: cidr,
			Targets: []*graph.RouteTarget{
				{Type: graph.InstanceTarget, Ref: "i-123"},
				{Type: graph.NetworkInterfaceTarget, Ref: "eni-456"},
			},
		},
	}
	got = col.format(routes)
	if got != "10.0.0.0/16->[inst:i-123 ni:eni-456 ] " {
		t.Fatalf("multiple targets: got %q", got)
	}

	// IPv6 only route
	routes = []*graph.Route{
		{
			DestinationIPv6: cidr6,
			Targets:         []*graph.RouteTarget{{Type: graph.VpcPeeringConnectionTarget, Ref: "pcx-1"}},
		},
	}
	got = col.format(routes)
	if got != "fd00::/64->vpc:pcx-1 " {
		t.Fatalf("ipv6 only: got %q", got)
	}

	// EgressOnlyInternetGateway target
	routes = []*graph.Route{
		{
			Destination: cidr,
			Targets:     []*graph.RouteTarget{{Type: graph.EgressOnlyInternetGatewayTarget, Ref: "eigw-1"}},
		},
	}
	got = col.format(routes)
	if got != "10.0.0.0/16->inbound-internget-gw:eigw-1 " {
		t.Fatalf("egress only igw: got %q", got)
	}
}

func TestGrantsColumnDefinitionFormat(t *testing.T) {
	col := GrantsColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Grants"},
	}

	if got := col.format(nil); got != "" {
		t.Fatalf("nil: got %q, want empty", got)
	}
	if got := col.format("not-grants"); got != "invalid grants" {
		t.Fatalf("non-grants: got %q, want 'invalid grants'", got)
	}

	// CanonicalUser with display name
	grants := []*graph.Grant{
		{
			Permission: "FULL_CONTROL",
			Grantee: graph.Grantee{
				GranteeType:        "CanonicalUser",
				GranteeID:          "abc123",
				GranteeDisplayName: "johndoe",
			},
		},
	}
	got := col.format(grants)
	if got != "FULL_CONTROL[user:johndoe] " {
		t.Fatalf("canonical user with display name: got %q", got)
	}

	// CanonicalUser without display name
	grants = []*graph.Grant{
		{
			Permission: "READ",
			Grantee: graph.Grantee{
				GranteeType: "CanonicalUser",
				GranteeID:   "abc123",
			},
		},
	}
	got = col.format(grants)
	if got != "READ[user:abc123] " {
		t.Fatalf("canonical user without display name: got %q", got)
	}

	// Group
	grants = []*graph.Grant{
		{
			Permission: "READ",
			Grantee: graph.Grantee{
				GranteeType: "Group",
				GranteeID:   "AllUsers",
			},
		},
	}
	got = col.format(grants)
	if got != "READ[group:AllUsers] " {
		t.Fatalf("group: got %q", got)
	}

	// Unknown type
	grants = []*graph.Grant{
		{
			Permission: "WRITE",
			Grantee: graph.Grantee{
				GranteeType: "AmazonCustomerByEmail",
				GranteeID:   "user@example.com",
			},
		},
	}
	got = col.format(grants)
	if got != "WRITE[AmazonCustomerByEmail:user@example.com] " {
		t.Fatalf("unknown type: got %q", got)
	}
}

func TestKeyValuesColumnDefinitionFormat(t *testing.T) {
	color.NoColor = true
	col := KeyValuesColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "Dimensions"},
	}

	if got := col.format(nil); got != "" {
		t.Fatalf("nil: got %q, want empty", got)
	}
	if got := col.format("not-keyvalues"); got != "invalid keyvalue, got string" {
		t.Fatalf("non-keyvalue: got %q", got)
	}

	kvs := []*graph.KeyValue{
		{KeyName: "env", Value: "prod"},
		{KeyName: "region", Value: "us-east-1"},
	}
	got := col.format(kvs)
	if got != "env:prod region:us-east-1" {
		t.Fatalf("key values: got %q", got)
	}

	// Single key value
	kvs = []*graph.KeyValue{
		{KeyName: "key", Value: "val"},
	}
	got = col.format(kvs)
	if got != "key:val" {
		t.Fatalf("single kv: got %q", got)
	}
}

func TestNameOrID(t *testing.T) {
	// Resource with Name
	r := graph.InitResource("instance", "inst_1")
	r.Properties()["Name"] = "my-instance"
	if got := nameOrID(r); got != "my-instance" {
		t.Fatalf("with Name: got %q, want 'my-instance'", got)
	}

	// Resource without Name but with ID property
	r2 := graph.InitResource("instance", "inst_2")
	r2.Properties()["Id"] = "inst_2"
	if got := nameOrID(r2); got != "inst_2" {
		t.Fatalf("with ID: got %q, want 'inst_2'", got)
	}

	// Resource without Name or ID property: falls back to Id() method
	r3 := graph.InitResource("instance", "inst_3")
	if got := nameOrID(r3); got != "inst_3" {
		t.Fatalf("fallback to Id(): got %q, want 'inst_3'", got)
	}

	// Resource with empty Name
	r4 := graph.InitResource("instance", "inst_4")
	r4.Properties()["Name"] = ""
	r4.Properties()["Id"] = "inst_4"
	if got := nameOrID(r4); got != "inst_4" {
		t.Fatalf("empty Name: got %q, want 'inst_4'", got)
	}
}

func TestResolveSortIndexes(t *testing.T) {
	headers := []ColumnDefinition{
		StringColumnDefinition{Prop: "ID"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
	}

	// Default (no args)
	ids, err := resolveSortIndexes(headers)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != 0 {
		t.Fatalf("default: got %v, want [0]", ids)
	}

	// Valid column
	ids, err = resolveSortIndexes(headers, "Name")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != 1 {
		t.Fatalf("Name: got %v, want [1]", ids)
	}

	// Multiple columns
	ids, err = resolveSortIndexes(headers, "State", "Name")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != 2 || ids[1] != 1 {
		t.Fatalf("State,Name: got %v, want [2,1]", ids)
	}

	// Invalid column
	_, err = resolveSortIndexes(headers, "NonExistent")
	if err == nil {
		t.Fatal("expected error for invalid column")
	}
}

func TestDefaultSorterSymbol(t *testing.T) {
	asc := &defaultSorter{sortBy: []int{0}, descending: false}
	if got := asc.symbol(); got != " ▲" {
		t.Fatalf("ascending: got %q, want ' ▲'", got)
	}

	desc := &defaultSorter{sortBy: []int{0}, descending: true}
	if got := desc.symbol(); got != " ▼" {
		t.Fatalf("descending: got %q, want ' ▼'", got)
	}
}

func TestDefaultSorterColumns(t *testing.T) {
	s := &defaultSorter{sortBy: []int{1, 3}}
	if got := s.columns(); len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Fatalf("got %v, want [1,3]", got)
	}
}

func TestValueLowerOrEqualTimeAndSlices(t *testing.T) {
	t1 := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

	// time.After means t1.After(t2) == false, so t1 is "lower"
	// In valueLowerOrEqual, time uses After: aa.After(bb)
	// t1.After(t2) = false => valueLowerOrEqual(t1, t2) = false
	// t2.After(t1) = true => valueLowerOrEqual(t2, t1) = true
	if got := valueLowerOrEqual(t2, t1); got != true {
		t.Fatalf("time t2 after t1: got %t, want true", got)
	}
	if got := valueLowerOrEqual(t1, t2); got != false {
		t.Fatalf("time t1 before t2: got %t, want false", got)
	}

	// Slice comparison
	s1 := []string{"a", "b"}
	s2 := []string{"c", "d"}
	if got := valueLowerOrEqual(s1, s2); got != true {
		t.Fatalf("string slice: got %t, want true", got)
	}

	// Int slice
	i1 := []int{1, 2}
	i2 := []int{3, 4}
	if got := valueLowerOrEqual(i1, i2); got != true {
		t.Fatalf("int slice: got %t, want true", got)
	}

	// nil comparisons
	if got := valueLowerOrEqual(nil, nil); got != true {
		t.Fatalf("both nil: got %t, want true", got)
	}
	if got := valueLowerOrEqual(nil, "a"); got != true {
		t.Fatalf("first nil: got %t, want true", got)
	}
	if got := valueLowerOrEqual("a", nil); got != false {
		t.Fatalf("second nil: got %t, want false", got)
	}
}

func TestValueLowerOrEqualPanicsOnUnknownType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown type")
		}
	}()
	valueLowerOrEqual(complex(1, 2), complex(3, 4))
}

func TestValueLowerOrEqualPanicsOnMismatchedTypes(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for mismatched types")
		}
	}()
	valueLowerOrEqual(1, "string")
}

func TestTSVDisplay(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(
		graph.InitResource("instance", "inst_1"),
	)
	inst := graph.InitResource("instance", "inst_1")
	inst.Properties()["ID"] = "inst_1"
	inst.Properties()["Name"] = "test"

	g2 := graph.NewGraph()
	g2.AddResource(inst)

	columns := []ColumnDefinition{
		StringColumnDefinition{Prop: "ID"},
		StringColumnDefinition{Prop: "Name"},
	}

	displayer, _ := BuildOptions(
		WithRdfType("instance"),
		WithColumnDefinitions(columns),
		WithFormat("tsv"),
	).SetSource(g2).Build()

	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	output := w.String()
	if output == "" {
		t.Fatal("expected non-empty TSV output")
	}
	// TSV should have tab-separated header
	if !bytes.Contains([]byte(output), []byte("ID\tName")) {
		t.Fatalf("expected tab-separated headers, got: %q", output)
	}
}

func TestHumanizeStorageInvalidUnit(t *testing.T) {
	got := HumanizeStorage(100, storageUnit(99))
	if got != "invalid storage unit" {
		t.Fatalf("got %q, want 'invalid storage unit'", got)
	}
}

func TestBuildQueryWithTagFilters(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithTagFilters([]string{"env=prod"}),
	)
	_, err := b.buildQuery()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestBuildQueryWithTagKeyFilters(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithTagKeyFilters([]string{"env"}),
	)
	_, err := b.buildQuery()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestBuildQueryWithTagValueFilters(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithTagValueFilters([]string{"production"}),
	)
	_, err := b.buildQuery()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestBuildQueryWithInvalidFilter(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithFilters([]string{"nonexistent_column=value"}),
	)
	_, err := b.buildQuery()
	if err == nil {
		t.Fatal("expected error for invalid filter key")
	}
}

func TestBuildQueryWithValidFilter(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithColumns([]string{"ID", "Name", "State"}),
		WithFilters([]string{"Name=redis"}),
	)
	_, err := b.buildQuery()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestBuildQueryNoFilters(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
	)
	_, err := b.buildQuery()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestBuildOptionsWithIDsOnly(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithIDsOnly(true),
	)
	if b.format != "porcelain" {
		t.Fatalf("format: got %q, want 'porcelain'", b.format)
	}
	if len(b.columnDefinitions) != 2 {
		t.Fatalf("column defs: got %d, want 2", len(b.columnDefinitions))
	}
}

func TestBuildOptionsWithIDsOnlyFalse(t *testing.T) {
	b := BuildOptions(
		WithRdfType("instance"),
		WithIDsOnly(false),
	)
	if b.format != "table" {
		t.Fatalf("format: got %q, want 'table'", b.format)
	}
}

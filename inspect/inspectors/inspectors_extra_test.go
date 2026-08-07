package inspectors

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/graph/resourcetest"
)

func TestGetRegion(t *testing.T) {
	g := graph.NewGraph()
	if err := g.AddResource(resourcetest.Region("us-east-1").Build()); err != nil {
		t.Fatal(err)
	}

	region, err := getRegion(g)
	if err != nil {
		t.Fatal(err)
	}
	if region != "us-east-1" {
		t.Fatalf("got %q, want 'us-east-1'", region)
	}
}

func TestGetRegionNoRegion(t *testing.T) {
	g := graph.NewGraph()

	_, err := getRegion(g)
	if err == nil {
		t.Fatal("expected error when no region in graph")
	}
	if !strings.Contains(err.Error(), "cannot resolve region") {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestPortScannerPrintPortRange(t *testing.T) {
	g := graph.NewGraph()

	_, cidr, _ := net.ParseCIDR("10.0.0.0/8")

	sg := resourcetest.SecurityGroup("sg-range").Prop("InboundRules", []*graph.FirewallRule{
		{
			PortRange: graph.PortRange{FromPort: 8000, ToPort: 9000},
			Protocol:  "tcp",
			IPRanges:  []*net.IPNet{cidr},
		},
	}).Build()

	if err := g.AddResource(sg); err != nil {
		t.Fatal(err)
	}

	scanner := &PortScanner{}
	if err := scanner.Inspect(g); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	scanner.Print(&buf)

	output := buf.String()
	if !strings.Contains(output, "ports 8000-9000") {
		t.Errorf("expected 'ports 8000-9000' in output, got: %s", output)
	}
	if !strings.Contains(output, "tcp") {
		t.Errorf("expected 'tcp' in output, got: %s", output)
	}
}

func TestPortScannerPrintWithTargets(t *testing.T) {
	g := graph.NewGraph()

	_, cidr, _ := net.ParseCIDR("10.0.0.0/8")

	sg := resourcetest.SecurityGroup("sg-targets").Prop("InboundRules", []*graph.FirewallRule{
		{
			PortRange: graph.PortRange{FromPort: 22, ToPort: 22},
			Protocol:  "tcp",
			IPRanges:  []*net.IPNet{cidr},
		},
	}).Build()

	if err := g.AddResource(sg); err != nil {
		t.Fatal(err)
	}

	scanner := &PortScanner{}
	if err := scanner.Inspect(g); err != nil {
		t.Fatal(err)
	}

	// When there are no targets, should say "nothing"
	var buf bytes.Buffer
	scanner.Print(&buf)

	output := buf.String()
	if !strings.Contains(output, "nothing") {
		t.Errorf("expected 'nothing' for no targets, got: %s", output)
	}
}

func TestPortScannerAllPermissiveNonAllIPs(t *testing.T) {
	g := graph.NewGraph()

	_, cidr, _ := net.ParseCIDR("10.0.0.0/8")

	sg := resourcetest.SecurityGroup("sg-partial").Prop("InboundRules", []*graph.FirewallRule{
		{
			PortRange: graph.PortRange{Any: true},
			Protocol:  "any",
			IPRanges:  []*net.IPNet{cidr},
		},
	}).Build()

	if err := g.AddResource(sg); err != nil {
		t.Fatal(err)
	}

	scanner := &PortScanner{}
	if err := scanner.Inspect(g); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	scanner.Print(&buf)

	output := buf.String()
	if !strings.Contains(output, "all ports via any protocol for IPs:") {
		t.Errorf("expected specific IPs message, got: %s", output)
	}
}

func TestPortScannerNoInboundRules(t *testing.T) {
	g := graph.NewGraph()

	sg := resourcetest.SecurityGroup("sg-norules").Build()

	if err := g.AddResource(sg); err != nil {
		t.Fatal(err)
	}

	scanner := &PortScanner{}
	if err := scanner.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if len(scanner.inbounds) != 0 {
		t.Errorf("expected 0 inbounds for SG without rules, got %d", len(scanner.inbounds))
	}
}

func TestBucketSizerEmptyGraph(t *testing.T) {
	g := graph.NewGraph()

	sizer := &BucketSizer{}
	if err := sizer.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if sizer.total != 0 {
		t.Errorf("expected total 0, got %d", sizer.total)
	}
	if len(sizer.buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(sizer.buckets))
	}
}

func TestOpenBucketsWithBothOpenTypes(t *testing.T) {
	g := graph.NewGraph()

	bothBucket := resourcetest.Bucket("both-bucket").Prop("Grants", []*graph.Grant{
		{
			Permission: "READ",
			Grantee: graph.Grantee{
				GranteeID:   "http://acs.amazonaws.com/groups/global/AllUsers",
				GranteeType: "Group",
			},
		},
		{
			Permission: "WRITE",
			Grantee: graph.Grantee{
				GranteeID:   "http://acs.amazonaws.com/groups/global/AuthenticatedUsers",
				GranteeType: "Group",
			},
		},
	}).Build()

	if err := g.AddResource(bothBucket); err != nil {
		t.Fatal(err)
	}

	inspector := &OpenBuckets{}
	if err := inspector.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if len(inspector.openToAny) != 1 {
		t.Errorf("expected 1 in openToAny, got %d", len(inspector.openToAny))
	}
	if len(inspector.openToAnyAuth) != 1 {
		t.Errorf("expected 1 in openToAnyAuth, got %d", len(inspector.openToAnyAuth))
	}

	var buf bytes.Buffer
	inspector.Print(&buf)
	output := buf.String()
	if !strings.Contains(output, "Buckets open to anybody") {
		t.Errorf("expected 'Buckets open to anybody', got: %s", output)
	}
	if !strings.Contains(output, "Buckets open to anyone with an AWS account") {
		t.Errorf("expected auth message, got: %s", output)
	}
}

func TestOpenBucketsNoGrants(t *testing.T) {
	g := graph.NewGraph()

	// Bucket without Grants property
	buck := graph.InitResource(cloud.Bucket, "no-grants-bucket")
	buck.Properties()["ID"] = "no-grants-bucket"

	if err := g.AddResource(buck); err != nil {
		t.Fatal(err)
	}

	inspector := &OpenBuckets{}
	if err := inspector.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if len(inspector.openToAny) != 0 {
		t.Errorf("expected 0 openToAny, got %d", len(inspector.openToAny))
	}
}

func TestPricerPrintFormatsTotal(t *testing.T) {
	p := &Pricer{
		total: 1.5,
		count: map[string]int{"m5.large": 3},
	}

	var buf bytes.Buffer
	p.Print(&buf)

	output := buf.String()
	// Total per day = 1.5 * 24 = 36.00
	if !strings.Contains(output, "$36.00") {
		t.Errorf("expected '$36.00' in output, got: %s", output)
	}
}

func TestBucketSizerPrintFormat(t *testing.T) {
	sizer := &BucketSizer{
		total: 2500000000,
		buckets: map[string]*bucket{
			"bucket-a": {objects: 5, size: 1500000000},
			"bucket-b": {objects: 3, size: 1000000000},
		},
	}

	var buf bytes.Buffer
	sizer.Print(&buf)

	output := buf.String()
	if !strings.Contains(output, "bucket-a") {
		t.Errorf("expected 'bucket-a' in output, got: %s", output)
	}
	if !strings.Contains(output, "bucket-b") {
		t.Errorf("expected 'bucket-b' in output, got: %s", output)
	}
	if !strings.Contains(output, "2.500000 Gb") {
		t.Errorf("expected total '2.500000 Gb' in output, got: %s", output)
	}
}

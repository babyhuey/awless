package inspectors

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestPortScannerName(t *testing.T) {
	p := &PortScanner{}
	if name := p.Name(); name != "port_scanner" {
		t.Errorf("expected 'port_scanner', got '%s'", name)
	}
}

func TestBucketSizerName(t *testing.T) {
	b := &BucketSizer{}
	if name := b.Name(); name != "bucket_sizer" {
		t.Errorf("expected 'bucket_sizer', got '%s'", name)
	}
}

func TestOpenBucketsName(t *testing.T) {
	o := &OpenBuckets{}
	if name := o.Name(); name != "open_buckets" {
		t.Errorf("expected 'open_buckets', got '%s'", name)
	}
}

func TestPricerName(t *testing.T) {
	p := &Pricer{}
	if name := p.Name(); name != "pricer" {
		t.Errorf("expected 'pricer', got '%s'", name)
	}
}

func TestPortScannerInspect(t *testing.T) {
	g := graph.NewGraph()

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	sg := resourcetest.SecurityGroup("sg-123").Prop("InboundRules", []*graph.FirewallRule{
		{
			PortRange: graph.PortRange{FromPort: 22, ToPort: 22},
			Protocol:  "tcp",
			IPRanges:  []*net.IPNet{cidr},
		},
		{
			PortRange: graph.PortRange{FromPort: 443, ToPort: 443},
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

	if len(scanner.inbounds) != 1 {
		t.Fatalf("expected 1 security group in inbounds, got %d", len(scanner.inbounds))
	}

	rules, ok := scanner.inbounds["sg-123"]
	if !ok {
		t.Fatal("expected security group 'sg-123' in inbounds")
	}

	if len(rules) != 2 {
		t.Errorf("expected 2 inbound rules, got %d", len(rules))
	}
}

func TestPortScannerPrint(t *testing.T) {
	g := graph.NewGraph()

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	sg := resourcetest.SecurityGroup("sg-456").Prop("InboundRules", []*graph.FirewallRule{
		{
			PortRange: graph.PortRange{FromPort: 80, ToPort: 80},
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
	if output == "" {
		t.Error("expected non-empty print output")
	}
	if !strings.Contains(output, "sg-456") {
		t.Error("expected output to contain security group ID 'sg-456'")
	}
	if !strings.Contains(output, "port 80") {
		t.Errorf("expected output to contain 'port 80', got: %s", output)
	}
}

func TestPortScannerAllPermissive(t *testing.T) {
	g := graph.NewGraph()

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	sg := resourcetest.SecurityGroup("sg-open").Prop("InboundRules", []*graph.FirewallRule{
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
	if !strings.Contains(output, "all ports via any protocol for all IPs") {
		t.Errorf("expected all-permissive message, got: %s", output)
	}
}

func TestBucketSizerInspect(t *testing.T) {
	g := graph.NewGraph()

	obj1 := graph.InitResource(cloud.S3Object, "obj-1")
	obj1.Properties()["Size"] = 1000
	obj1.Properties()["Bucket"] = "my-bucket"

	obj2 := graph.InitResource(cloud.S3Object, "obj-2")
	obj2.Properties()["Size"] = 2000
	obj2.Properties()["Bucket"] = "my-bucket"

	obj3 := graph.InitResource(cloud.S3Object, "obj-3")
	obj3.Properties()["Size"] = 500
	obj3.Properties()["Bucket"] = "other-bucket"

	if err := g.AddResource(obj1, obj2, obj3); err != nil {
		t.Fatal(err)
	}

	sizer := &BucketSizer{}
	if err := sizer.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if sizer.total != 3500 {
		t.Errorf("expected total size 3500, got %d", sizer.total)
	}

	if len(sizer.buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(sizer.buckets))
	}

	myBucket := sizer.buckets["my-bucket"]
	if myBucket == nil {
		t.Fatal("expected 'my-bucket' to be present")
	}
	if myBucket.objects != 2 {
		t.Errorf("expected 2 objects in 'my-bucket', got %d", myBucket.objects)
	}
	if myBucket.size != 3000 {
		t.Errorf("expected size 3000 for 'my-bucket', got %d", myBucket.size)
	}

	otherBucket := sizer.buckets["other-bucket"]
	if otherBucket == nil {
		t.Fatal("expected 'other-bucket' to be present")
	}
	if otherBucket.objects != 1 {
		t.Errorf("expected 1 object in 'other-bucket', got %d", otherBucket.objects)
	}
	if otherBucket.size != 500 {
		t.Errorf("expected size 500 for 'other-bucket', got %d", otherBucket.size)
	}
}

func TestBucketSizerPrint(t *testing.T) {
	g := graph.NewGraph()

	obj := graph.InitResource(cloud.S3Object, "obj-1")
	obj.Properties()["Size"] = 1000000000
	obj.Properties()["Bucket"] = "test-bucket"

	if err := g.AddResource(obj); err != nil {
		t.Fatal(err)
	}

	sizer := &BucketSizer{}
	if err := sizer.Inspect(g); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	sizer.Print(&buf)

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty print output")
	}
	if !strings.Contains(output, "test-bucket") {
		t.Errorf("expected output to contain 'test-bucket', got: %s", output)
	}
	if !strings.Contains(output, "Gb") {
		t.Errorf("expected output to contain 'Gb', got: %s", output)
	}
}

func TestOpenBucketsInspect(t *testing.T) {
	g := graph.NewGraph()

	openBucket := resourcetest.Bucket("open-bucket-1").Prop("Grants", []*graph.Grant{
		{
			Permission: "FULL_CONTROL",
			Grantee: graph.Grantee{
				GranteeID:   "http://acs.amazonaws.com/groups/global/AllUsers",
				GranteeType: "Group",
			},
		},
	}).Build()

	authBucket := resourcetest.Bucket("auth-bucket-1").Prop("Grants", []*graph.Grant{
		{
			Permission: "READ",
			Grantee: graph.Grantee{
				GranteeID:   "http://acs.amazonaws.com/groups/global/AuthenticatedUsers",
				GranteeType: "Group",
			},
		},
	}).Build()

	privateBucket := resourcetest.Bucket("private-bucket-1").Prop("Grants", []*graph.Grant{
		{
			Permission: "FULL_CONTROL",
			Grantee: graph.Grantee{
				GranteeID:   "owner-canonical-id",
				GranteeType: "CanonicalUser",
			},
		},
	}).Build()

	if err := g.AddResource(openBucket, authBucket, privateBucket); err != nil {
		t.Fatal(err)
	}

	inspector := &OpenBuckets{}
	if err := inspector.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if len(inspector.openToAny) != 1 {
		t.Errorf("expected 1 bucket open to any, got %d", len(inspector.openToAny))
	} else if inspector.openToAny[0] != "open-bucket-1" {
		t.Errorf("expected 'open-bucket-1' in openToAny, got '%s'", inspector.openToAny[0])
	}

	if len(inspector.openToAnyAuth) != 1 {
		t.Errorf("expected 1 bucket open to authenticated users, got %d", len(inspector.openToAnyAuth))
	} else if inspector.openToAnyAuth[0] != "auth-bucket-1" {
		t.Errorf("expected 'auth-bucket-1' in openToAnyAuth, got '%s'", inspector.openToAnyAuth[0])
	}
}

func TestOpenBucketsNoneFound(t *testing.T) {
	g := graph.NewGraph()

	privateBucket := resourcetest.Bucket("private-1").Prop("Grants", []*graph.Grant{
		{
			Permission: "FULL_CONTROL",
			Grantee: graph.Grantee{
				GranteeID:   "owner-id",
				GranteeType: "CanonicalUser",
			},
		},
	}).Build()

	if err := g.AddResource(privateBucket); err != nil {
		t.Fatal(err)
	}

	inspector := &OpenBuckets{}
	if err := inspector.Inspect(g); err != nil {
		t.Fatal(err)
	}

	if len(inspector.openToAny) != 0 {
		t.Errorf("expected 0 open buckets, got %d", len(inspector.openToAny))
	}
	if len(inspector.openToAnyAuth) != 0 {
		t.Errorf("expected 0 auth-open buckets, got %d", len(inspector.openToAnyAuth))
	}

	var buf bytes.Buffer
	inspector.Print(&buf)
	if !strings.Contains(buf.String(), "none found") {
		t.Errorf("expected 'none found' message, got: %s", buf.String())
	}
}

func TestOpenBucketsPrint(t *testing.T) {
	g := graph.NewGraph()

	openBucket := resourcetest.Bucket("public-bucket").Prop("Grants", []*graph.Grant{
		{
			Permission: "READ",
			Grantee: graph.Grantee{
				GranteeID:   "http://acs.amazonaws.com/groups/global/AllUsers",
				GranteeType: "Group",
			},
		},
	}).Build()

	if err := g.AddResource(openBucket); err != nil {
		t.Fatal(err)
	}

	inspector := &OpenBuckets{}
	if err := inspector.Inspect(g); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	inspector.Print(&buf)

	output := buf.String()
	if !strings.Contains(output, "Buckets open to anybody") {
		t.Errorf("expected 'Buckets open to anybody' in output, got: %s", output)
	}
	if !strings.Contains(output, "public-bucket") {
		t.Errorf("expected 'public-bucket' in output, got: %s", output)
	}
}

func TestPricerPrint(t *testing.T) {
	p := &Pricer{
		total: 0.05,
		count: map[string]int{"t2.micro": 2, "m4.large": 1},
	}

	var buf bytes.Buffer
	p.Print(&buf)

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty print output")
	}
	if !strings.Contains(output, "t2.micro") {
		t.Errorf("expected output to contain 't2.micro', got: %s", output)
	}
	if !strings.Contains(output, "m4.large") {
		t.Errorf("expected output to contain 'm4.large', got: %s", output)
	}
}

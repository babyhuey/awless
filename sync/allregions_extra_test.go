package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/sync"
)

// writeRegionGraph lays down a synced graph for one region, the way `awless sync` would.
func writeRegionGraph(t *testing.T, home, profile, region, service string, resources ...*graph.Resource) {
	t.Helper()

	dir := filepath.Join(home, "aws", "rdf", profile, region)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}

	g := graph.NewGraph()
	if err := g.AddResource(resources...); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(filepath.Join(dir, service+".nt"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := g.MarshalTo(f); err != nil {
		t.Fatal(err)
	}
}

func instance(id, name string) *graph.Resource {
	r := graph.InitResource(cloud.Instance, id)
	r.SetProperty(properties.Name, name)
	return r
}

// The point of the feature: resources from every synced region in one listing, each
// carrying the region it came from.
func TestLoadLocalGraphForTypeInAllRegions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("__AWLESS_HOME", home)

	writeRegionGraph(t, home, "default", "us-west-2", "infra", instance("i-west", "web"))
	writeRegionGraph(t, home, "default", "eu-west-1", "infra", instance("i-east", "api"))

	g, err := sync.LoadLocalGraphForTypeInAllRegions("infra", cloud.Instance, "default")
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery(cloud.Instance))
	if err != nil {
		t.Fatal(err)
	}
	if len(resources) != 2 {
		t.Fatalf("got %d instances, want 2 (one per region)", len(resources))
	}

	regionOf := map[string]string{}
	for _, r := range resources {
		region, _ := r.Property(properties.Region)
		regionOf[r.ID()] = region.(string)
	}
	if regionOf["i-west"] != "us-west-2" {
		t.Errorf("i-west region is %q, want us-west-2", regionOf["i-west"])
	}
	if regionOf["i-east"] != "eu-west-1" {
		t.Errorf("i-east region is %q, want eu-west-1", regionOf["i-east"])
	}
}

// The failure that makes the whole feature useless: every row tagged with one region, so
// the column is there but wrong. Two resources sharing a name in different regions is the
// case that exposes it, and it is also the realistic one — naming conventions repeat
// across regions.
func TestLoadLocalGraphForTypeInAllRegionsKeepsRegionsDistinct(t *testing.T) {
	home := t.TempDir()
	t.Setenv("__AWLESS_HOME", home)

	writeRegionGraph(t, home, "default", "us-west-2", "infra", instance("i-1", "web-server"))
	writeRegionGraph(t, home, "default", "eu-west-1", "infra", instance("i-2", "web-server"))
	writeRegionGraph(t, home, "default", "ap-south-1", "infra", instance("i-3", "web-server"))

	g, err := sync.LoadLocalGraphForTypeInAllRegions("infra", cloud.Instance, "default")
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery(cloud.Instance))
	if err != nil {
		t.Fatal(err)
	}
	if len(resources) != 3 {
		t.Fatalf("got %d instances, want 3", len(resources))
	}

	seen := map[string]int{}
	for _, r := range resources {
		region, ok := r.Property(properties.Region)
		if !ok {
			t.Fatalf("instance %s has no region property", r.ID())
		}
		seen[region.(string)]++
	}
	if len(seen) != 3 {
		t.Errorf("got regions %v; each instance should carry its own region", seen)
	}
	for _, region := range []string{"us-west-2", "eu-west-1", "ap-south-1"} {
		if seen[region] != 1 {
			t.Errorf("region %s appears %d times, want 1", region, seen[region])
		}
	}
}

// A global service is stored once under "global", so asking for it per region would either
// find nothing or report a meaningless region.
func TestLoadLocalGraphForTypeInAllRegionsHandlesGlobalServices(t *testing.T) {
	home := t.TempDir()
	t.Setenv("__AWLESS_HOME", home)

	user := graph.InitResource(cloud.User, "jsmith")
	writeRegionGraph(t, home, "default", "global", "access", user)
	// A regional directory exists too, to prove the global service is not looked for in it.
	writeRegionGraph(t, home, "default", "us-west-2", "infra", instance("i-1", "web"))

	g, err := sync.LoadLocalGraphForTypeInAllRegions("access", cloud.User, "default")
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.Find(cloud.NewQuery(cloud.User))
	if err != nil {
		t.Fatal(err)
	}
	if len(resources) != 1 {
		t.Fatalf("got %d users, want 1", len(resources))
	}
	region, _ := resources[0].Property(properties.Region)
	if region != "global" {
		t.Errorf("user region is %q, want global", region)
	}
}

// Nothing synced is an empty listing, not an error.
func TestLoadLocalGraphForTypeInAllRegionsWithNothingSynced(t *testing.T) {
	home := t.TempDir()
	t.Setenv("__AWLESS_HOME", home)

	g, err := sync.LoadLocalGraphForTypeInAllRegions("infra", cloud.Instance, "default")
	if err != nil {
		t.Fatalf("an unsynced profile should not error: %s", err)
	}
	resources, err := g.Find(cloud.NewQuery(cloud.Instance))
	if err != nil {
		t.Fatal(err)
	}
	if len(resources) != 0 {
		t.Errorf("got %d resources, want none", len(resources))
	}
}

func TestLocalRegionsExcludesGlobalAndFiles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("__AWLESS_HOME", home)

	writeRegionGraph(t, home, "default", "us-west-2", "infra", instance("i-1", "web"))
	writeRegionGraph(t, home, "default", "eu-west-1", "infra", instance("i-2", "web"))
	writeRegionGraph(t, home, "default", "global", "access", graph.InitResource(cloud.User, "jsmith"))
	// The sync directory is a git repository, whose .git must not be read as a region.
	if err := os.MkdirAll(filepath.Join(home, "aws", "rdf", "default", ".git"), 0o700); err != nil {
		t.Fatal(err)
	}

	got := sync.LocalRegions("default")

	want := []string{"eu-west-1", "us-west-2"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v (sorted, no global, no dotfiles)", got, want)
		}
	}
}

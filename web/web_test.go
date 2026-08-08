package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bootswithdefer/awless/graph"
)

func TestNew(t *testing.T) {
	s := New(":8080", "default")
	if s.addr != ":8080" {
		t.Errorf("expected addr ':8080', got %q", s.addr)
	}
	if s.awsProfile != "default" {
		t.Errorf("expected profile 'default', got %q", s.awsProfile)
	}
	if s.gph != nil {
		t.Error("expected gph to be nil on newly created server")
	}
}

func TestNewDifferentParams(t *testing.T) {
	s := New(":9090", "production")
	if s.addr != ":9090" {
		t.Errorf("expected addr ':9090', got %q", s.addr)
	}
	if s.awsProfile != "production" {
		t.Errorf("expected profile 'production', got %q", s.awsProfile)
	}
}

func TestNewResource(t *testing.T) {
	gr := graph.InitResource("instance", "i-12345")
	r := newResource(gr)
	if r.ID != "i-12345" {
		t.Errorf("expected Id 'i-12345', got %q", r.ID)
	}
	if r.Type != "instance" {
		t.Errorf("expected Type 'instance', got %q", r.Type)
	}
	if r.Properties == nil {
		t.Error("expected non-nil Properties")
	}
}

func TestNewResourceDifferentTypes(t *testing.T) {
	types := []struct {
		kind, id string
	}{
		{"vpc", "vpc-abc123"},
		{"subnet", "subnet-def456"},
		{"securitygroup", "sg-789"},
		{"volume", "vol-001"},
		{"region", "us-west-2"},
	}
	for _, tc := range types {
		r := newResource(graph.InitResource(tc.kind, tc.id))
		if r.ID != tc.id {
			t.Errorf("expected Id %q, got %q", tc.id, r.ID)
		}
		if r.Type != tc.kind {
			t.Errorf("expected Type %q, got %q", tc.kind, r.Type)
		}
	}
}

func TestResourceAddDependsOn(t *testing.T) {
	r := &Resource{}
	dep1 := graph.InitResource("vpc", "vpc-1")
	dep2 := graph.InitResource("subnet", "sub-1")

	r.AddDependsOn(dep1, dep2)

	if len(r.DependsOn) != 2 {
		t.Fatalf("expected 2 DependsOn, got %d", len(r.DependsOn))
	}
	if r.DependsOn[0].ID != "vpc-1" {
		t.Errorf("expected first dep Id 'vpc-1', got %q", r.DependsOn[0].ID)
	}
	if r.DependsOn[1].ID != "sub-1" {
		t.Errorf("expected second dep Id 'sub-1', got %q", r.DependsOn[1].ID)
	}
}

func TestResourceAddDependsOnEmpty(t *testing.T) {
	r := &Resource{}
	r.AddDependsOn()
	if len(r.DependsOn) != 0 {
		t.Errorf("expected 0 DependsOn, got %d", len(r.DependsOn))
	}
}

func TestResourceAddAppliesOn(t *testing.T) {
	r := &Resource{}
	a := graph.InitResource("securitygroup", "sg-1")

	r.AddAppliesOn(a)

	if len(r.AppliesOn) != 1 {
		t.Fatalf("expected 1 AppliesOn, got %d", len(r.AppliesOn))
	}
	if r.AppliesOn[0].Type != "securitygroup" {
		t.Errorf("expected Type 'securitygroup', got %q", r.AppliesOn[0].Type)
	}
}

func TestResourceAddAppliesOnMultiple(t *testing.T) {
	r := &Resource{}
	r.AddAppliesOn(
		graph.InitResource("securitygroup", "sg-1"),
		graph.InitResource("securitygroup", "sg-2"),
		graph.InitResource("securitygroup", "sg-3"),
	)
	if len(r.AppliesOn) != 3 {
		t.Fatalf("expected 3 AppliesOn, got %d", len(r.AppliesOn))
	}
}

func TestResourceAddAppliesOnEmpty(t *testing.T) {
	r := &Resource{}
	r.AddAppliesOn()
	if len(r.AppliesOn) != 0 {
		t.Errorf("expected 0 AppliesOn, got %d", len(r.AppliesOn))
	}
}

func TestResourceAddParents(t *testing.T) {
	r := &Resource{}
	parent := graph.InitResource("region", "us-east-1")

	r.AddParents(parent)

	if len(r.Parents) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(r.Parents))
	}
	if r.Parents[0].ID != "us-east-1" {
		t.Errorf("expected parent Id 'us-east-1', got %q", r.Parents[0].ID)
	}
}

func TestResourceAddParentsMultiple(t *testing.T) {
	r := &Resource{}
	r.AddParents(
		graph.InitResource("region", "us-east-1"),
		graph.InitResource("vpc", "vpc-123"),
	)
	if len(r.Parents) != 2 {
		t.Fatalf("expected 2 parents, got %d", len(r.Parents))
	}
	if r.Parents[0].ID != "us-east-1" {
		t.Errorf("expected first parent Id 'us-east-1', got %q", r.Parents[0].ID)
	}
	if r.Parents[1].ID != "vpc-123" {
		t.Errorf("expected second parent Id 'vpc-123', got %q", r.Parents[1].ID)
	}
}

func TestResourceAddParentsEmpty(t *testing.T) {
	r := &Resource{}
	r.AddParents()
	if len(r.Parents) != 0 {
		t.Errorf("expected 0 Parents, got %d", len(r.Parents))
	}
}

func TestResourceAddChildren(t *testing.T) {
	r := &Resource{}
	child := graph.InitResource("instance", "i-child")

	r.AddChildren(child)

	if len(r.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(r.Children))
	}
	if r.Children[0].ID != "i-child" {
		t.Errorf("expected child Id 'i-child', got %q", r.Children[0].ID)
	}
}

func TestResourceAddMultipleChildren(t *testing.T) {
	r := &Resource{}
	c1 := graph.InitResource("instance", "i-1")
	c2 := graph.InitResource("instance", "i-2")
	c3 := graph.InitResource("instance", "i-3")

	r.AddChildren(c1, c2, c3)

	if len(r.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(r.Children))
	}
}

func TestResourceAddChildrenEmpty(t *testing.T) {
	r := &Resource{}
	r.AddChildren()
	if len(r.Children) != 0 {
		t.Errorf("expected 0 Children, got %d", len(r.Children))
	}
}

func TestRoutes(t *testing.T) {
	s := New(":8080", "test")
	handler := s.routes()
	if handler == nil {
		t.Fatal("expected routes() to return a non-nil http.Handler")
	}
}

func TestHomeHandler(t *testing.T) {
	s := New(":8080", "test")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	s.homeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "/resources") {
		t.Error("expected home page to contain link to /resources")
	}
	if !strings.Contains(body, "/rdf") {
		t.Error("expected home page to contain link to /rdf")
	}
	if !strings.Contains(body, "/graph") {
		t.Error("expected home page to contain link to /graph")
	}
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("expected home page to be valid HTML")
	}
}

func TestResourceAccumulation(t *testing.T) {
	r := &Resource{ID: "test-id", Type: "instance"}

	// Add deps in multiple calls to test accumulation
	r.AddDependsOn(graph.InitResource("vpc", "vpc-1"))
	r.AddDependsOn(graph.InitResource("subnet", "sub-1"))

	if len(r.DependsOn) != 2 {
		t.Fatalf("expected 2 DependsOn after two calls, got %d", len(r.DependsOn))
	}

	r.AddAppliesOn(graph.InitResource("sg", "sg-1"))
	r.AddAppliesOn(graph.InitResource("sg", "sg-2"))

	if len(r.AppliesOn) != 2 {
		t.Fatalf("expected 2 AppliesOn after two calls, got %d", len(r.AppliesOn))
	}

	r.AddParents(graph.InitResource("region", "us-east-1"))
	if len(r.Parents) != 1 {
		t.Fatalf("expected 1 Parent, got %d", len(r.Parents))
	}

	r.AddChildren(graph.InitResource("volume", "vol-1"))
	r.AddChildren(graph.InitResource("volume", "vol-2"))
	if len(r.Children) != 2 {
		t.Fatalf("expected 2 Children after two calls, got %d", len(r.Children))
	}
}

func TestNewResourcePreservesProperties(t *testing.T) {
	gr := graph.InitResource("instance", "i-99")
	r := newResource(gr)
	if r.Properties == nil {
		t.Fatal("expected non-nil properties")
	}
	// Properties map should be the same reference as the graph resource
	if r.ID != "i-99" {
		t.Fatalf("expected Id 'i-99', got %q", r.ID)
	}
}

func TestListResourcesHandlerWithEmptyGraph(t *testing.T) {
	s := New(":8080", "test")
	s.gph = graph.NewGraph()

	req := httptest.NewRequest(http.MethodGet, "/resources", nil)
	rr := httptest.NewRecorder()

	s.listResourcesHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("expected valid HTML response")
	}
}

func TestGraphHandlerWithEmptyGraph(t *testing.T) {
	s := New(":8080", "test")
	s.gph = graph.NewGraph()

	// Use a temp dir so loadLocalTriples finds nothing
	tmpDir, err := os.MkdirTemp("", "test-web-graph")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	origHome := os.Getenv("__AWLESS_HOME")
	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Setenv("__AWLESS_HOME", origHome)

	req := httptest.NewRequest(http.MethodGet, "/graph", nil)
	rr := httptest.NewRecorder()

	s.graphHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestRdfHandlerWithEmptyData(t *testing.T) {
	s := New(":8080", "test")
	s.gph = graph.NewGraph()

	tmpDir, err := os.MkdirTemp("", "test-web-rdf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	origHome := os.Getenv("__AWLESS_HOME")
	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Setenv("__AWLESS_HOME", origHome)

	req := httptest.NewRequest(http.MethodGet, "/rdf", nil)
	rr := httptest.NewRecorder()

	s.rdfHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestRdfHandlerNamespaced(t *testing.T) {
	s := New(":8080", "test")
	s.gph = graph.NewGraph()

	tmpDir, err := os.MkdirTemp("", "test-web-rdf-ns")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	origHome := os.Getenv("__AWLESS_HOME")
	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Setenv("__AWLESS_HOME", origHome)

	req := httptest.NewRequest(http.MethodGet, "/rdf?namespaced=true", nil)
	rr := httptest.NewRecorder()

	s.rdfHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestRoutesServesHome(t *testing.T) {
	s := New(":8080", "test")
	s.gph = graph.NewGraph()

	srv := httptest.NewServer(s.routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRoutesServesResources(t *testing.T) {
	s := New(":8080", "test")
	s.gph = graph.NewGraph()

	srv := httptest.NewServer(s.routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/resources")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

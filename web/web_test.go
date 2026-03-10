package web

import (
	"net/http"
	"testing"

	"github.com/wallix/awless/graph"
)

func TestNew(t *testing.T) {
	s := New(":8080", "default")
	if s.port != ":8080" {
		t.Errorf("expected port ':8080', got %q", s.port)
	}
	if s.awsProfile != "default" {
		t.Errorf("expected profile 'default', got %q", s.awsProfile)
	}
	if s.gph != nil {
		t.Error("expected gph to be nil on newly created server")
	}
}

func TestNewResource(t *testing.T) {
	gr := graph.InitResource("instance", "i-12345")
	r := newResource(gr)
	if r.Id != "i-12345" {
		t.Errorf("expected Id 'i-12345', got %q", r.Id)
	}
	if r.Type != "instance" {
		t.Errorf("expected Type 'instance', got %q", r.Type)
	}
	if r.Properties == nil {
		t.Error("expected non-nil Properties")
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
	if r.DependsOn[0].Id != "vpc-1" {
		t.Errorf("expected first dep Id 'vpc-1', got %q", r.DependsOn[0].Id)
	}
	if r.DependsOn[1].Id != "sub-1" {
		t.Errorf("expected second dep Id 'sub-1', got %q", r.DependsOn[1].Id)
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

func TestResourceAddParents(t *testing.T) {
	r := &Resource{}
	parent := graph.InitResource("region", "us-east-1")

	r.AddParents(parent)

	if len(r.Parents) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(r.Parents))
	}
	if r.Parents[0].Id != "us-east-1" {
		t.Errorf("expected parent Id 'us-east-1', got %q", r.Parents[0].Id)
	}
}

func TestResourceAddChildren(t *testing.T) {
	r := &Resource{}
	child := graph.InitResource("instance", "i-child")

	r.AddChildren(child)

	if len(r.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(r.Children))
	}
	if r.Children[0].Id != "i-child" {
		t.Errorf("expected child Id 'i-child', got %q", r.Children[0].Id)
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

func TestRoutes(t *testing.T) {
	s := New(":8080", "test")
	handler := s.routes()
	if handler == nil {
		t.Fatal("expected routes() to return a non-nil http.Handler")
	}
	if _, ok := handler.(http.Handler); !ok {
		t.Fatal("expected routes() to return an http.Handler")
	}
}

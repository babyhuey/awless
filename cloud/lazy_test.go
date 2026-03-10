/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cloud

import (
	"bytes"
	"io"
	"testing"
)

func TestLazyLoadingGraph(t *testing.T) {
	var nbCalls int
	loadingFunc := func() GraphAPI {
		nbCalls++
		return &StubGraph{}
	}

	lazy := &LazyGraph{LoadingFunc: loadingFunc}
	lazy.FindOne(NewQuery(""))
	if got, want := nbCalls, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	lazy.FindOne(NewQuery(""))
	if got, want := nbCalls, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

type StubGraph struct {
}

func (g *StubGraph) Find(Query) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) FindWithProperties(props map[string]interface{}) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) FilterGraph(Query) (GraphAPI, error) {
	return nil, nil
}

func (g *StubGraph) FindOne(Query) (Resource, error) {
	return nil, nil
}

func (g *StubGraph) MarshalTo(io.Writer) error {
	return nil
}

func (g *StubGraph) ResourceRelations(Resource, string, bool) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) VisitRelations(Resource, string, bool, func(Resource, int) error) error {
	return nil
}

func (g *StubGraph) ResourceSiblings(Resource) ([]Resource, error) {
	return nil, nil
}

func (g *StubGraph) Merge(GraphAPI) error {
	return nil
}

func newLazyWithStub() (*LazyGraph, *StubGraph) {
	stub := &StubGraph{}
	lazy := &LazyGraph{LoadingFunc: func() GraphAPI {
		return stub
	}}
	return lazy, stub
}

func TestLazyGraphFind(t *testing.T) {
	lazy, _ := newLazyWithStub()
	res, err := lazy.Find(NewQuery(Instance))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil result, got %v", res)
	}
}

func TestLazyGraphFindWithProperties(t *testing.T) {
	lazy, _ := newLazyWithStub()
	res, err := lazy.FindWithProperties(map[string]interface{}{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil result, got %v", res)
	}
}

func TestLazyGraphFilterGraph(t *testing.T) {
	lazy, _ := newLazyWithStub()
	g, err := lazy.FilterGraph(NewQuery(Instance))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g != nil {
		t.Fatalf("expected nil graph, got %v", g)
	}
}

func TestLazyGraphMarshalTo(t *testing.T) {
	lazy, _ := newLazyWithStub()
	var buf bytes.Buffer
	err := lazy.MarshalTo(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLazyGraphResourceRelations(t *testing.T) {
	lazy, _ := newLazyWithStub()
	res, err := lazy.ResourceRelations(nil, "parentOf", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil result, got %v", res)
	}
}

func TestLazyGraphVisitRelations(t *testing.T) {
	lazy, _ := newLazyWithStub()
	err := lazy.VisitRelations(nil, "parentOf", false, func(r Resource, depth int) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLazyGraphResourceSiblings(t *testing.T) {
	lazy, _ := newLazyWithStub()
	res, err := lazy.ResourceSiblings(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil result, got %v", res)
	}
}

func TestLazyGraphMerge(t *testing.T) {
	lazy, _ := newLazyWithStub()
	err := lazy.Merge(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLazyGraphLoadOnlyOnce(t *testing.T) {
	var callCount int
	lazy := &LazyGraph{LoadingFunc: func() GraphAPI {
		callCount++
		return &StubGraph{}
	}}

	// Call multiple different methods
	lazy.Find(NewQuery(""))
	lazy.FindWithProperties(nil)
	lazy.FilterGraph(NewQuery(""))
	lazy.MarshalTo(io.Discard)
	lazy.ResourceRelations(nil, "", false)
	lazy.VisitRelations(nil, "", false, nil)
	lazy.ResourceSiblings(nil)
	lazy.Merge(nil)

	if callCount != 1 {
		t.Fatalf("loading function called %d times, expected 1", callCount)
	}
}

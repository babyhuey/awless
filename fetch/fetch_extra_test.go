package fetch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
)

func TestErrorAdd(t *testing.T) {
	fe := fetch.WrapError()

	if fe.Any() {
		t.Fatal("new Error should not have any errors")
	}

	fe.Add(nil)
	if fe.Any() {
		t.Fatal("adding nil should not create an error")
	}

	fe.Add(errors.New("first"))
	if !fe.Any() {
		t.Fatal("expected Any() to be true after adding an error")
	}

	fe.Add(errors.New("second"))
	if got := fe.Error(); got != "first\nsecond" {
		t.Fatalf("got %q, want %q", got, "first\nsecond")
	}
}

func TestErrorString(t *testing.T) {
	fe := fetch.WrapError()
	if got := fe.Error(); got != "" {
		t.Fatalf("empty Error should return empty string, got %q", got)
	}

	fe.Add(errors.New("alpha"))
	if got := fe.Error(); got != "alpha" {
		t.Fatalf("got %q, want %q", got, "alpha")
	}
}

func TestWrapErrorNil(t *testing.T) {
	fe := fetch.WrapError(nil)
	if fe.Any() {
		t.Fatal("wrapping nil should produce empty Error")
	}
}

func TestWrapErrorSingle(t *testing.T) {
	fe := fetch.WrapError(errors.New("single"))
	if !fe.Any() {
		t.Fatal("expected Any() true")
	}
	if got := fe.Error(); got != "single" {
		t.Fatalf("got %q, want %q", got, "single")
	}
}

func TestWrapErrorFlattensNestedError(t *testing.T) {
	inner := fetch.WrapError(errors.New("a"), errors.New("b"))
	outer := fetch.WrapError(inner, errors.New("c"))

	if got := outer.Error(); got != "a\nb\nc" {
		t.Fatalf("got %q, want %q", got, "a\nb\nc")
	}
}

func TestNewFetcherFetchReturnsResources(t *testing.T) {
	resources := []*graph.Resource{
		graph.InitResource("thing", "t1"),
		graph.InitResource("thing", "t2"),
	}

	f := fetch.NewFetcher(fetch.Funcs{
		"thing": func(ctx context.Context, c fetch.Cache) ([]*graph.Resource, interface{}, error) {
			return resources, nil, nil
		},
	})

	gph, err := f.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	res, _ := gph.GetResource("thing", "t1")
	if res == nil {
		t.Fatal("expected to find resource t1")
	}
	res, _ = gph.GetResource("thing", "t2")
	if res == nil {
		t.Fatal("expected to find resource t2")
	}
}

func TestNewFetcherGetAndReset(t *testing.T) {
	f := fetch.NewFetcher(fetch.Funcs{})

	// Store and Get
	f.Store("mykey", "myval")
	val, err := f.Get("mykey")
	if err != nil {
		t.Fatal(err)
	}
	if val != "myval" {
		t.Fatalf("got %v, want %v", val, "myval")
	}

	// Get with func
	val, err = f.Get("computed", func() (interface{}, error) {
		return 42, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if val != 42 {
		t.Fatalf("got %v, want 42", val)
	}

	// Reset clears cache
	f.Reset()
	val, err = f.Get("mykey")
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatalf("expected nil after reset, got %v", val)
	}
}

func TestFetcherFetchCollectsErrors(t *testing.T) {
	f := fetch.NewFetcher(fetch.Funcs{
		"failing": func(ctx context.Context, c fetch.Cache) ([]*graph.Resource, interface{}, error) {
			return nil, nil, errors.New("boom")
		},
	})

	gph, err := f.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error from Fetch")
	}
	if gph == nil {
		t.Fatal("expected non-nil graph even on error")
	}
}

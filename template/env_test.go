package template

import (
	"context"
	"sync"
	"testing"

	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
)

// RequestContext is documented as never returning nil so callers can pass it
// straight to the AWS SDK without a guard.
func TestRunEnvRequestContextIsNeverNil(t *testing.T) {
	renv := NewRunEnv(NewEnv().Build())

	if renv.RequestContext() == nil {
		t.Fatal("expected a non-nil context before one is set")
	}

	type ctxKey string
	want := context.WithValue(context.Background(), ctxKey("k"), "v")
	renv.SetRequestContext(want)

	if got := renv.RequestContext(); got != want {
		t.Errorf("expected the context that was set, got %#v", got)
	}
}

// Context returns a copy: a command mutating the map it receives must not corrupt
// the environment for the statements that follow.
func TestRunEnvContextIsACopy(t *testing.T) {
	cenv := NewEnv().Build()
	cenv.Push(env.ResolvedVars, map[string]any{"myvar": "original"})
	renv := NewRunEnv(cenv, map[string]any{"direct": "value"})

	first := renv.Context()
	first["direct"] = "mutated"
	first["added"] = "new"

	second := renv.Context()
	if got := second["direct"]; got != "value" {
		t.Errorf("mutating the returned map changed the env: got %v, want value", got)
	}
	if _, ok := second["added"]; ok {
		t.Error("a key added to the returned map leaked into the env")
	}
}

// AWLESS, Variables and References are three names for the same resolved vars,
// kept for compatibility with templates written against older releases.
func TestRunEnvExposesResolvedVarsUnderAllThreeNames(t *testing.T) {
	cenv := NewEnv().Build()
	cenv.Push(env.ResolvedVars, map[string]any{"instance": "i-1234"})

	ctx := NewRunEnv(cenv).Context()

	for _, key := range []string{"AWLESS", "Variables", "References"} {
		vars, ok := ctx[key].(map[string]any)
		if !ok {
			t.Errorf("%s: expected a map, got %#v", key, ctx[key])
			continue
		}
		if got := vars["instance"]; got != "i-1234" {
			t.Errorf("%s: got %v, want i-1234", key, got)
		}
	}
}

func TestRunEnvDryRun(t *testing.T) {
	renv := NewRunEnv(NewEnv().Build())
	if renv.IsDryRun() {
		t.Error("expected dry run to default to false")
	}
	renv.SetDryRun(true)
	if !renv.IsDryRun() {
		t.Error("expected dry run to be true after being set")
	}
}

// Get returns a copy for the same reason Context does.
func TestDataMapGetIsACopy(t *testing.T) {
	cenv := NewEnv().Build()
	cenv.Push(env.ResolvedVars, map[string]any{"a": 1})

	got := cenv.Get(env.ResolvedVars)
	got["a"] = 99
	got["b"] = 2

	again := cenv.Get(env.ResolvedVars)
	if again["a"] != 1 {
		t.Errorf("mutating the returned map changed the store: got %v, want 1", again["a"])
	}
	if _, ok := again["b"]; ok {
		t.Error("a key added to the returned map leaked into the store")
	}
}

func TestDataMapPushMergesAndSeparatesByType(t *testing.T) {
	cenv := NewEnv().Build()
	cenv.Push(env.ResolvedVars, map[string]any{"a": 1})
	cenv.Push(env.ResolvedVars, map[string]any{"b": 2})
	cenv.Push(env.Fillers, map[string]any{"c": 3})

	resolved := cenv.Get(env.ResolvedVars)
	if len(resolved) != 2 || resolved["a"] != 1 || resolved["b"] != 2 {
		t.Errorf("expected both pushes merged, got %#v", resolved)
	}
	if _, leaked := resolved["c"]; leaked {
		t.Error("a value pushed under Fillers appeared under ResolvedVars")
	}
	if fillers := cenv.Get(env.Fillers); fillers["c"] != 3 {
		t.Errorf("expected c under Fillers, got %#v", fillers)
	}
}

// Get on a type nothing was pushed to must return an empty map, not nil, since
// callers index it directly.
func TestDataMapGetUnsetTypeReturnsEmptyMap(t *testing.T) {
	got := NewEnv().Build().Get(env.ProcessedFillers)
	if got == nil {
		t.Fatal("expected an empty map, got nil")
	}
	if len(got) != 0 {
		t.Errorf("expected an empty map, got %#v", got)
	}
}

// dataMap carries a mutex, so concurrent use is intended. Run under -race.
func TestDataMapConcurrentPushAndGet(t *testing.T) {
	cenv := NewEnv().Build()

	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			cenv.Push(env.ResolvedVars, map[string]any{"k": n})
		}(i)
		go func() {
			defer wg.Done()
			_ = cenv.Get(env.ResolvedVars)
		}()
	}
	wg.Wait()

	if _, ok := cenv.Get(env.ResolvedVars)["k"]; !ok {
		t.Error("expected the key to be present after concurrent pushes")
	}
}

func TestEnvBuilderWiresEachFunc(t *testing.T) {
	aliasCalled, holesCalled, lookupCalled := false, false, false

	cenv := NewEnv().
		WithAliasFunc(func(string, string) string { aliasCalled = true; return "alias" }).
		WithMissingHolesFunc(func(string, []string, bool) string { holesCalled = true; return "hole" }).
		WithLookupCommandFunc(func(...string) any { lookupCalled = true; return "cmd" }).
		WithParamsMode(env.RequiredParamsOnly).
		WithLog(logger.DiscardLogger).
		Build()

	if got := cenv.AliasFunc()("path", "a"); got != "alias" || !aliasCalled {
		t.Errorf("alias func not wired: got %q", got)
	}
	if got := cenv.MissingHolesFunc()("k", nil, false); got != "hole" || !holesCalled {
		t.Errorf("missing holes func not wired: got %q", got)
	}
	if got := cenv.LookupCommandFunc()("create", "vpc"); got != "cmd" || !lookupCalled {
		t.Errorf("lookup command func not wired: got %v", got)
	}
	if got := cenv.ParamsMode(); got != env.RequiredParamsOnly {
		t.Errorf("params mode: got %d, want %d", got, env.RequiredParamsOnly)
	}
	if cenv.Log() == nil {
		t.Error("expected a logger")
	}
}

// A fresh env must be usable without any With* call, since Compile builds one that
// way in several places.
func TestNewEnvDefaultsAreUsable(t *testing.T) {
	cenv := NewEnv().Build()

	if cenv.LookupCommandFunc() == nil {
		t.Error("expected a default lookup func rather than nil")
	}
	if got := cenv.LookupCommandFunc()("anything"); got != nil {
		t.Errorf("expected the default lookup to resolve nothing, got %v", got)
	}
	if cenv.Log() == nil {
		t.Error("expected a default logger")
	}
}

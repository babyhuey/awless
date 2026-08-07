package params

import (
	"errors"
	"reflect"
	"testing"
)

func TestReduce(t *testing.T) {
	data := map[string]any{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	plusOne := func(in map[string]any) (out map[string]any, err error) {
		out = make(map[string]any)
		for k, i := range in {
			out[k] = i.(int) + 1
		}
		return
	}
	red := newReducer(plusOne, "one", "three")
	out, err := red.Reduce(data)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := red.Keys(), []string{"one", "three"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := len(out), 2; got != want {
		t.Fatalf("got length %d, want length %d", got, want)
	}
	if v, ok := out["one"].(int); !ok || v != 2 {
		t.Fatalf("invalid content: %v", out)
	}
	if v, ok := out["three"].(int); !ok || v != 4 {
		t.Fatalf("invalid content: %v", out)
	}
}

func TestReduceNoMatchingKeys(t *testing.T) {
	data := map[string]any{
		"a": 1,
		"b": 2,
	}
	identity := func(in map[string]any) (map[string]any, error) {
		return in, nil
	}
	red := newReducer(identity, "x", "y")
	out, err := red.Reduce(data)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(out); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestReduceError(t *testing.T) {
	data := map[string]any{"key": "val"}
	errFn := func(in map[string]any) (map[string]any, error) {
		return nil, errors.New("reduce error")
	}
	red := newReducer(errFn, "key")
	_, err := red.Reduce(data)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "reduce error" {
		t.Fatalf("got %v, want 'reduce error'", err)
	}
}

func TestReduceEmptyKeys(t *testing.T) {
	identity := func(in map[string]any) (map[string]any, error) {
		return in, nil
	}
	red := newReducer(identity)
	if got := len(red.Keys()); got != 0 {
		t.Fatalf("expected 0 keys, got %d", got)
	}
	out, err := red.Reduce(map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(out); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestReducePartialKeyMatch(t *testing.T) {
	data := map[string]any{
		"one": 1,
		"two": 2,
	}
	identity := func(in map[string]any) (map[string]any, error) {
		return in, nil
	}
	red := newReducer(identity, "one", "missing")
	out, err := red.Reduce(data)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(out); got != 1 {
		t.Fatalf("expected 1 entry, got %d", got)
	}
	if out["one"] != 1 {
		t.Fatalf("expected one=1, got %v", out["one"])
	}
}

func TestSpecRule(t *testing.T) {
	r := AllOf(Key("a"), Key("b"))
	s := newSpec(r)
	if s.Rule() == nil {
		t.Fatal("expected non-nil rule")
	}
}

func TestSpecNilRule(t *testing.T) {
	s := newSpec(nil)
	rule := s.Rule()
	if rule == nil {
		t.Fatal("expected None() rule, not nil")
	}
	// None rule should return no error
	if err := rule.Run([]string{"anything"}); err != nil {
		t.Fatalf("expected no error from None rule, got %v", err)
	}
}

func TestSpecValidators(t *testing.T) {
	v := Validators{"key": MaxLengthOf(10)}
	s := newSpec(None(), v)
	got := s.Validators()
	if got == nil {
		t.Fatal("expected non-nil validators")
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(got))
	}
}

func TestSpecNilValidators(t *testing.T) {
	s := newSpec(None())
	got := s.Validators()
	if got == nil {
		t.Fatal("expected non-nil validators map")
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 validators, got %d", len(got))
	}
}

func TestSpecReducers(t *testing.T) {
	s := newSpec(None())
	if s.Reducers() != nil {
		t.Fatal("expected nil reducers for basic spec")
	}
}

func TestSpecBuilder(t *testing.T) {
	identity := func(in map[string]any) (map[string]any, error) {
		return in, nil
	}
	s := SpecBuilder(AllOf(Key("a")), Validators{"a": MinLengthOf(1)}).
		AddReducer(identity, "a", "b").
		AddReducer(identity, "c").
		Done()

	if s.Rule() == nil {
		t.Fatal("expected non-nil rule")
	}
	if len(s.Reducers()) != 2 {
		t.Fatalf("expected 2 reducers, got %d", len(s.Reducers()))
	}
	if len(s.Validators()) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(s.Validators()))
	}
}

func TestSpecBuilderNoReducers(t *testing.T) {
	s := SpecBuilder(None()).Done()
	if len(s.Reducers()) != 0 {
		t.Fatalf("expected 0 reducers, got %d", len(s.Reducers()))
	}
}

func TestNewSpecPublic(t *testing.T) {
	s := NewSpec(AllOf(Key("x")))
	if s.Rule() == nil {
		t.Fatal("expected non-nil rule")
	}
	if s.Reducers() != nil {
		t.Fatal("expected nil reducers")
	}
}

func TestNoneRule(t *testing.T) {
	n := None()
	if err := n.Run([]string{"anything"}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got := n.Required(); len(got) != 0 {
		t.Fatalf("expected 0 required, got %v", got)
	}
	if got := n.Missing([]string{"a"}); len(got) != 0 {
		t.Fatalf("expected 0 missing, got %v", got)
	}
	if got := n.String(); got != "none" {
		t.Fatalf("expected 'none', got %q", got)
	}
	// Visit should be a no-op
	called := false
	n.Visit(func(r Rule) { called = true })
	if called {
		t.Fatal("expected Visit to be a no-op for None")
	}
}

func TestOptString(t *testing.T) {
	o := Opt("a", "b", "c")
	got := o.String()
	if got != "[a b c]" {
		t.Fatalf("expected '[a b c]', got %q", got)
	}
}

func TestOptWithSuggested(t *testing.T) {
	o := Opt("a", Suggested("b", "c"))
	got := o.String()
	if got != "[a b c]" {
		t.Fatalf("expected '[a b c]', got %q", got)
	}
}

func TestOptRun(t *testing.T) {
	o := Opt("a")
	err := o.Run([]string{"a"})
	if err == nil {
		t.Fatal("Opt.Run should always return optErr")
	}
}

func TestOptMissing(t *testing.T) {
	o := Opt("a")
	miss := o.Missing([]string{})
	if len(miss) != 0 {
		t.Fatalf("Opt.Missing should return empty, got %v", miss)
	}
}

func TestOptRequired(t *testing.T) {
	o := Opt("a", "b")
	req := o.Required()
	if len(req) != 0 {
		t.Fatalf("Opt.Required should return empty, got %v", req)
	}
}

func TestKeyString(t *testing.T) {
	k := Key("mykey")
	if got := k.String(); got != "mykey" {
		t.Fatalf("expected 'mykey', got %q", got)
	}
}

func TestKeyRun(t *testing.T) {
	k := Key("mykey")
	if err := k.Run([]string{"mykey", "other"}); err != nil {
		t.Fatalf("expected no error when key is present, got %v", err)
	}
	if err := k.Run([]string{"other"}); err == nil {
		t.Fatal("expected error when key is missing")
	}
}

func TestKeyRequired(t *testing.T) {
	k := Key("mykey")
	req := k.Required()
	if len(req) != 1 || req[0] != "mykey" {
		t.Fatalf("expected ['mykey'], got %v", req)
	}
}

func TestKeyMissing(t *testing.T) {
	k := Key("mykey")
	miss := k.Missing([]string{"other"})
	if len(miss) != 1 || miss[0] != "mykey" {
		t.Fatalf("expected ['mykey'], got %v", miss)
	}
	miss = k.Missing([]string{"mykey"})
	if len(miss) != 0 {
		t.Fatalf("expected empty, got %v", miss)
	}
}

func TestKeyVisit(t *testing.T) {
	k := Key("x")
	var visited []string
	k.Visit(func(r Rule) {
		visited = append(visited, r.String())
	})
	if len(visited) != 1 || visited[0] != "x" {
		t.Fatalf("expected ['x'], got %v", visited)
	}
}

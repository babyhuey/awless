package console

import (
	"sort"
	"testing"
	"time"
)

// Resource properties are interface{} from heterogeneous AWS responses, so a
// column can hold mixed or unexpected types. Display code must degrade to a
// stable order rather than crash the CLI.
func TestValueLowerOrEqualNeverPanics(t *testing.T) {
	values := []interface{}{
		nil,
		"str", "",
		1, 0, -3,
		1.5, 0.0,
		true, false,
		time.Now(), time.Time{},
		[]string{"a", "b"}, []string{},
		[]int{1, 2},
		map[string]string{"k": "v"},        // no case for this type
		struct{ A int }{1},                 // nor this
		[]interface{}{1, "mixed"},          // nor this
		uint64(7), int64(9), float32(1.25), // numeric types with no case
	}

	for i, a := range values {
		for j, b := range values {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panicked comparing values[%d] (%T) and values[%d] (%T): %v", i, a, j, b, r)
					}
				}()
				_ = valueLowerOrEqual(a, b)
			}()
		}
	}
}

// The mixed-type case that used to panic: a column holding a string for one
// resource and a number for another.
func TestValueLowerOrEqualMixedTypes(t *testing.T) {
	if _, ok := interface{}(valueLowerOrEqual("10", 10)).(bool); !ok {
		t.Fatal("expected a bool result")
	}

	// Must be antisymmetric so it remains a usable sort predicate: for a != b,
	// at most one direction may report "lower or equal".
	a, b := interface{}("abc"), interface{}(42)
	if valueLowerOrEqual(a, b) && valueLowerOrEqual(b, a) {
		t.Error("comparator reports both directions lower-or-equal for distinct values")
	}
}

// sort.Slice with a comparator that panics aborts the whole program; confirm a
// mixed-type column sorts cleanly instead.
func TestSortMixedColumnDoesNotPanic(t *testing.T) {
	col := []interface{}{"beta", 3, nil, true, 1.5, "alpha", []string{"x"}}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("sorting a mixed-type column panicked: %v", r)
		}
	}()

	sort.Slice(col, func(i, j int) bool { return valueLowerOrEqual(col[i], col[j]) })

	if len(col) != 7 {
		t.Errorf("expected 7 elements after sort, got %d", len(col))
	}
}

// Same-type ordering must be unaffected by the fallback.
func TestValueLowerOrEqualSameTypeOrdering(t *testing.T) {
	tcases := []struct {
		desc string
		a, b interface{}
		want bool
	}{
		{desc: "ints ascending", a: 1, b: 2, want: true},
		{desc: "ints descending", a: 2, b: 1, want: false},
		{desc: "equal ints", a: 2, b: 2, want: true},
		{desc: "strings ascending", a: "a", b: "b", want: true},
		{desc: "strings descending", a: "b", b: "a", want: false},
		{desc: "floats ascending", a: 1.0, b: 2.0, want: true},
		{desc: "bool false first", a: false, b: true, want: true},
		{desc: "nil sorts first", a: nil, b: "x", want: true},
		{desc: "nil second", a: "x", b: nil, want: false},
	}

	for _, tc := range tcases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := valueLowerOrEqual(tc.a, tc.b); got != tc.want {
				t.Errorf("valueLowerOrEqual(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

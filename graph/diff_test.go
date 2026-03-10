package graph

import (
	"testing"

	tstore "github.com/wallix/triplestore"

	"github.com/wallix/awless/cloud/rdf"
)

func TestNewDiff(t *testing.T) {
	from := NewGraph()
	to := NewGraph()
	d := NewDiff(from, to)
	if d == nil {
		t.Fatal("expected non-nil diff")
	}
	if got := d.FromGraph(); got != from {
		t.Fatalf("FromGraph: got %v, want %v", got, from)
	}
	if got := d.ToGraph(); got != to {
		t.Fatalf("ToGraph: got %v, want %v", got, to)
	}
	if d.HasDiff() {
		t.Fatal("expected HasDiff() to be false for new diff")
	}
}

func TestDiffMergedGraph(t *testing.T) {
	from := NewGraph()
	from.store.Add(
		tstore.SubjPred("res1", rdf.RdfType).Resource("cloud-owl:Instance"),
		tstore.SubjPred("res1", MetaPredicate).StringLiteral("extra"),
	)

	to := NewGraph()
	to.store.Add(
		tstore.SubjPred("res2", rdf.RdfType).Resource("cloud-owl:Instance"),
	)

	d := NewDiff(from, to)
	merged := d.MergedGraph()
	if merged == nil {
		t.Fatal("expected non-nil merged graph")
	}

	snap := merged.store.Snapshot()
	// The merged graph should contain triples from 'to'
	if !snap.Contains(tstore.SubjPred("res2", rdf.RdfType).Resource("cloud-owl:Instance")) {
		t.Fatal("merged graph should contain 'to' triples")
	}
	// The 'from' meta should be marked as missing
	metaTriples := snap.WithSubjPred("res1", MetaPredicate)
	found := false
	for _, tri := range metaTriples {
		txt, _ := tstore.ParseString(tri.Object())
		if txt == missingLit {
			found = true
		}
	}
	if !found {
		t.Fatal("expected meta triple with 'missing' literal in merged graph")
	}
}

func TestHierarchicDifferRun(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		from := NewGraph()
		to := NewGraph()
		diff, err := DefaultDiffer.Run("root", from, to)
		if err != nil {
			t.Fatal(err)
		}
		if diff.HasDiff() {
			t.Fatal("expected no diff for two empty graphs")
		}
	})

	t.Run("identical graphs", func(t *testing.T) {
		from := NewGraph()
		from.store.Add(
			tstore.SubjPred("root", rdf.ParentOf).Resource("child1"),
			tstore.SubjPred("child1", rdf.RdfType).Resource("cloud-owl:Instance"),
		)
		to := NewGraph()
		to.store.Add(
			tstore.SubjPred("root", rdf.ParentOf).Resource("child1"),
			tstore.SubjPred("child1", rdf.RdfType).Resource("cloud-owl:Instance"),
		)
		diff, err := DefaultDiffer.Run("root", from, to)
		if err != nil {
			t.Fatal(err)
		}
		if diff.HasDiff() {
			t.Fatal("expected no diff for identical graphs")
		}
	})

	t.Run("extra in to", func(t *testing.T) {
		from := NewGraph()
		from.store.Add(
			tstore.SubjPred("root", rdf.ParentOf).Resource("child1"),
		)
		to := NewGraph()
		to.store.Add(
			tstore.SubjPred("root", rdf.ParentOf).Resource("child1"),
			tstore.SubjPred("root", rdf.ParentOf).Resource("child2"),
		)
		diff, err := DefaultDiffer.Run("root", from, to)
		if err != nil {
			t.Fatal(err)
		}
		if !diff.HasDiff() {
			t.Fatal("expected diff when 'to' has extra children")
		}
	})

	t.Run("missing in to", func(t *testing.T) {
		from := NewGraph()
		from.store.Add(
			tstore.SubjPred("root", rdf.ParentOf).Resource("child1"),
			tstore.SubjPred("root", rdf.ParentOf).Resource("child2"),
		)
		to := NewGraph()
		to.store.Add(
			tstore.SubjPred("root", rdf.ParentOf).Resource("child1"),
		)
		diff, err := DefaultDiffer.Run("root", from, to)
		if err != nil {
			t.Fatal(err)
		}
		if !diff.HasDiff() {
			t.Fatal("expected diff when 'to' is missing children")
		}
	})
}

func TestIntersectTriples(t *testing.T) {
	a := []tstore.Triple{
		tstore.SubjPred("s1", "p1").StringLiteral("o1"),
		tstore.SubjPred("s2", "p2").StringLiteral("o2"),
	}
	b := []tstore.Triple{
		tstore.SubjPred("s2", "p2").StringLiteral("o2"),
		tstore.SubjPred("s3", "p3").StringLiteral("o3"),
	}

	result := intersectTriples(a, b)
	if got, want := len(result), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	// Empty intersection
	c := []tstore.Triple{
		tstore.SubjPred("s4", "p4").StringLiteral("o4"),
	}
	result = intersectTriples(a, c)
	if got, want := len(result), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	// Both empty
	result = intersectTriples(nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestSubtractTriples(t *testing.T) {
	a := []tstore.Triple{
		tstore.SubjPred("s1", "p1").StringLiteral("o1"),
		tstore.SubjPred("s2", "p2").StringLiteral("o2"),
	}
	b := []tstore.Triple{
		tstore.SubjPred("s2", "p2").StringLiteral("o2"),
	}

	result := subtractTriples(a, b)
	if got, want := len(result), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	// Subtract from empty
	result = subtractTriples(nil, b)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}

	// Subtract empty
	result = subtractTriples(a, nil)
	if got, want := len(result), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestMax(t *testing.T) {
	tcases := []struct {
		a, b, expected uint32
	}{
		{0, 0, 0},
		{1, 0, 1},
		{0, 1, 1},
		{5, 10, 10},
		{10, 5, 10},
		{42, 42, 42},
	}
	for _, tc := range tcases {
		if got := max(tc.a, tc.b); got != tc.expected {
			t.Fatalf("max(%d, %d): got %d, want %d", tc.a, tc.b, got, tc.expected)
		}
	}
}

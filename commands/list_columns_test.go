package commands

import (
	"strings"
	"testing"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/console"
)

// withListFlags sets the listing flags for one test and restores them, since they are
// package-level and shared with every other listing test.
func withListFlags(t *testing.T, allRegions bool, columns []string) {
	t.Helper()
	prevAll, prevCols := allLocalRegionsFlag, listingColumnsFlag
	allLocalRegionsFlag, listingColumnsFlag = allRegions, columns
	t.Cleanup(func() { allLocalRegionsFlag, listingColumnsFlag = prevAll, prevCols })
}

func contains(cols []string, want string) bool {
	for _, c := range cols {
		if strings.EqualFold(c, want) {
			return true
		}
	}
	return false
}

// Without the flag nothing changes, so ordinary listings are unaffected.
func TestListingColumnsUnchangedWithoutTheFlag(t *testing.T) {
	withListFlags(t, false, []string{"id", "name"})

	got := listingColumns(cloud.Instance)

	if len(got) != 2 || got[0] != "id" || got[1] != "name" {
		t.Errorf("got %v, want the flag value untouched", got)
	}
}

func TestListingColumnsUnchangedWithoutTheFlagAndNoColumns(t *testing.T) {
	withListFlags(t, false, nil)

	if got := listingColumns(cloud.Instance); len(got) != 0 {
		t.Errorf("got %v, want empty so the displayer applies its own defaults", got)
	}
}

// The region is what makes a merged listing interpretable, so it is added to the defaults.
func TestListingColumnsAddsRegionToTheDefaults(t *testing.T) {
	withListFlags(t, true, nil)

	got := listingColumns(cloud.Instance)

	if !contains(got, properties.Region) {
		t.Errorf("got %v, want a region column", got)
	}
	if len(got) != len(console.ColumnsInListing[cloud.Instance])+1 {
		t.Errorf("got %d columns, want the %d defaults plus region", len(got), len(console.ColumnsInListing[cloud.Instance]))
	}
	if got[len(got)-1] != properties.Region {
		t.Errorf("region is at position %d; it should be appended last", len(got)-1)
	}
}

// An explicit --columns still controls the order of what the user asked for.
func TestListingColumnsAppendsRegionToExplicitColumns(t *testing.T) {
	withListFlags(t, true, []string{"id", "name"})

	got := listingColumns(cloud.Instance)

	want := []string{"id", "name", properties.Region}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

// Asking for the region explicitly must not produce it twice.
func TestListingColumnsDoesNotDuplicateRegion(t *testing.T) {
	for _, cols := range [][]string{
		{"id", "region"},
		{"id", "Region"}, // the column flag is case-insensitive elsewhere
	} {
		withListFlags(t, true, cols)

		got := listingColumns(cloud.Instance)

		var n int
		for _, c := range got {
			if strings.EqualFold(c, properties.Region) {
				n++
			}
		}
		if n != 1 {
			t.Errorf("columns %v produced %d region columns, want 1 (got %v)", cols, n, got)
		}
	}
}

// The defaults are a package-level map, so the returned slice must not share its backing
// array — a caller appending to it would leak a column into every later listing of that
// type in the same process.
//
// This uses availabilityzone deliberately: its default columns already include the region,
// so listingColumns returns early without appending, and that early return is the only path
// where the slice could still be the map's own. For a type that does append, the append
// reallocates (the defaults are literals at exact capacity) and hides the difference — a
// test using one of those would pass whether or not the copy existed.
func TestListingColumnsCopiesTheDefaults(t *testing.T) {
	withListFlags(t, true, nil)

	defaults := console.ColumnsInListing[cloud.AvailabilityZone]
	if len(defaults) == 0 || !contains(defaults, properties.Region) {
		t.Skip("availabilityzone no longer defaults to a region column; pick another type that does")
	}

	got := listingColumns(cloud.AvailabilityZone)

	if len(got) == 0 {
		t.Fatal("no columns returned")
	}
	if &got[0] == &defaults[0] {
		t.Error("the returned columns share the default list's backing array; a caller appending would mutate it")
	}
}

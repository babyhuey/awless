package template

import (
	"strings"
	"testing"
)

// IntRangeValue used to be `[0-9]+'-'[0-9]+` with nothing stopping it at a token boundary,
// so it matched the leading `12345678-1234` of a UUID and left the rest unparseable. PEG's
// ordered choice does not reconsider an alternative that already succeeded, so the value
// never fell through to UnquotedParam, which would have accepted the whole thing.
//
// The values below are all real: AWS Backup plan IDs are UUIDs, and a date or timestamp is
// what several params take.
func TestHyphenatedValuesParse(t *testing.T) {
	for _, v := range []string{
		"12345678-1234-1234-1234-123456789012", // a UUID, e.g. a backup plan ID
		"2024-02-09",                           // an ISO date
		"2024-02-09T10:30:00Z",                 // an ISO timestamp
		"1-2-3",                                // three numeric groups
		"22-80abc",                             // a range-looking prefix with a suffix
		"1-2.3",
		"00-00-00",
	} {
		if _, err := Parse("delete backupplan id=" + v); err != nil {
			t.Errorf("%q should parse unquoted: %s", v, err)
		}
	}
}

// The guard must not have cost the case it was protecting: a bare port range is still a
// value, and still a string rather than being split or rejected.
func TestIntRangesStillParse(t *testing.T) {
	for _, v := range []string{"22-80", "1-10", "0-65535", "1-1"} {
		tpl, err := Parse("update securitygroup id=sg-1234 portrange=" + v)
		if err != nil {
			t.Fatalf("%q should still parse: %s", v, err)
		}
		for _, cmd := range tpl.CommandNodesIterator() {
			got, ok := cmd.ToDriverParams()["portrange"]
			if !ok {
				t.Fatalf("%q: no portrange param", v)
			}
			if got != v {
				t.Errorf("%q became %v (%T), want the string unchanged", v, got, got)
			}
		}
	}
}

// Quoting was the documented workaround, and must keep working for anyone whose templates
// or stored log lines already use it.
func TestQuotedHyphenatedValuesStillParse(t *testing.T) {
	for _, v := range []string{
		`'12345678-1234-1234-1234-123456789012'`,
		`"12345678-1234-1234-1234-123456789012"`,
		`'22-80'`,
		`"2024-02-09"`,
	} {
		tpl, err := Parse("delete backupplan id=" + v)
		if err != nil {
			t.Fatalf("%s should parse: %s", v, err)
		}
		want := strings.Trim(v, `'"`)
		for _, cmd := range tpl.CommandNodesIterator() {
			if got := cmd.ToDriverParams()["id"]; got != want {
				t.Errorf("%s became %v, want %q", v, got, want)
			}
		}
	}
}

// The local template log persists command lines and re-parses them on read, so the grammar
// is a compatibility surface rather than only an input parser. A value that parses but does
// not survive the round trip would corrupt `awless log` and any revert built from it.
func TestHyphenatedValuesSurviveTheLogRoundTrip(t *testing.T) {
	for _, line := range []string{
		"delete backupplan id=12345678-1234-1234-1234-123456789012",
		"update securitygroup id=sg-1234 portrange=22-80",
		"delete backupplan id=2024-02-09",
	} {
		tpl, err := Parse(line)
		if err != nil {
			t.Fatalf("%q: %s", line, err)
		}

		// String() is what gets persisted.
		persisted := tpl.String()

		reparsed, err := Parse(persisted)
		if err != nil {
			t.Fatalf("%q persisted as %q, which no longer parses: %s", line, persisted, err)
		}
		if again := reparsed.String(); again != persisted {
			t.Errorf("%q is not stable across the round trip: %q then %q", line, persisted, again)
		}
	}
}

package console

import (
	"testing"
)

func TestToShortArn(t *testing.T) {
	tcases := []struct {
		input, expected string
	}{
		{"arn:aws:iam::123456789012:user/johndoe", "user/johndoe"},
		{"arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/my-tg/abc123", "targetgroup/my-tg/abc123"},
		{"no-colons-here", "no-colons-here"},
		{"", ""},
		{"single:colon", "colon"},
		{"a:b:c:d:e:f", "f"},
	}

	for _, tc := range tcases {
		if got := ToShortArn(tc.input); got != tc.expected {
			t.Fatalf("ToShortArn(%q): got %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestARNLastValueColumnDefinitionFormat(t *testing.T) {
	col := ARNLastValueColumnDefinition{
		Separator:              "/",
		StringColumnDefinition: StringColumnDefinition{Prop: "test"},
	}

	tcases := []struct {
		input    interface{}
		expected string
	}{
		{"arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster", "my-cluster"},
		{"no-slash", "no-slash"},
		{nil, ""},
		{"a/b/c", "c"},
	}

	for _, tc := range tcases {
		if got := col.format(tc.input); got != tc.expected {
			t.Fatalf("ARNLastValueColumnDefinition.format(%v): got %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestStringColumnDefinitionTitle(t *testing.T) {
	// With Friendly name
	col := StringColumnDefinition{Prop: "AvailabilityZone", Friendly: "Zone"}
	if got := col.title(); got != "Zone" {
		t.Fatalf("got %q, want %q", got, "Zone")
	}
	if got := col.title(" asc"); got != "Zone asc" {
		t.Fatalf("got %q, want %q", got, "Zone asc")
	}

	// Without Friendly name
	col2 := StringColumnDefinition{Prop: "ID"}
	if got := col2.title(); got != "ID" {
		t.Fatalf("got %q, want %q", got, "ID")
	}
}

func TestStringColumnDefinitionFormat(t *testing.T) {
	col := StringColumnDefinition{Prop: "test"}

	if got := col.format(nil); got != "" {
		t.Fatalf("format(nil): got %q, want %q", got, "")
	}
	if got := col.format("hello"); got != "hello" {
		t.Fatalf("format(string): got %q, want %q", got, "hello")
	}
	if got := col.format(42); got != "42" {
		t.Fatalf("format(int): got %q, want %q", got, "42")
	}
}

func TestSliceColumnDefinitionFormat(t *testing.T) {
	col := SliceColumnDefinition{
		StringColumnDefinition: StringColumnDefinition{Prop: "test"},
	}

	if got := col.format(nil); got != "" {
		t.Fatalf("format(nil): got %q, want %q", got, "")
	}
	if got := col.format([]string{"a", "b", "c"}); got != "a b c" {
		t.Fatalf("format(slice): got %q, want %q", got, "a b c")
	}
	if got := col.format("not-a-slice"); got != "invalid slice: string" {
		t.Fatalf("format(non-slice): got %q, want %q", got, "invalid slice: string")
	}

	// With ForEach transformer
	colWithForEach := SliceColumnDefinition{
		ForEach:                ToShortArn,
		StringColumnDefinition: StringColumnDefinition{Prop: "test"},
	}
	input := []string{"arn:aws:iam::123:user/alice", "arn:aws:iam::123:user/bob"}
	if got := colWithForEach.format(input); got != "user/alice user/bob" {
		t.Fatalf("format with ForEach: got %q, want %q", got, "user/alice user/bob")
	}
}

func TestColumnDefinitionsResolveKey(t *testing.T) {
	defs := ColumnDefinitions{
		StringColumnDefinition{Prop: "AvailabilityZone", Friendly: "Zone"},
		StringColumnDefinition{Prop: "ID"},
	}

	tcases := []struct {
		input, expected string
	}{
		{"zone", "AvailabilityZone"},
		{"Zone", "AvailabilityZone"},
		{"availabilityzone", "AvailabilityZone"},
		{"id", "ID"},
		{"ID", "ID"},
		{"nonexistent", ""},
	}

	for _, tc := range tcases {
		if got := defs.resolveKey(tc.input); got != tc.expected {
			t.Fatalf("resolveKey(%q): got %q, want %q", tc.input, got, tc.expected)
		}
	}
}

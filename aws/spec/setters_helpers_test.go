package awsspec

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// --- setValueAtPath tests ---

type simpleStruct struct {
	Name  string
	Value int
	Flag  bool
}

type innerStruct struct {
	Label string
	Count int64
}

type outerStruct struct {
	Inner  innerStruct
	InnerP *innerStruct
	Name   string
}

type ptrFieldStruct struct {
	StrPtr   *string
	IntPtr   *int32
	Int64Ptr *int64
	FloatPtr *float64
	BoolPtr  *bool
}

type enumType string

type enumStruct struct {
	Status enumType
}

func TestSetValueAtPathSimpleField(t *testing.T) {
	s := &simpleStruct{}
	setValueAtPath(s, "Name", "hello")
	if s.Name != "hello" {
		t.Errorf("expected Name='hello', got %q", s.Name)
	}

	setValueAtPath(s, "Value", 42)
	if s.Value != 42 {
		t.Errorf("expected Value=42, got %d", s.Value)
	}

	setValueAtPath(s, "Flag", true)
	if s.Flag != true {
		t.Errorf("expected Flag=true, got %v", s.Flag)
	}
}

func TestSetValueAtPathNestedPath(t *testing.T) {
	s := &outerStruct{}
	setValueAtPath(s, "Inner.Label", "nested")
	if s.Inner.Label != "nested" {
		t.Errorf("expected Inner.Label='nested', got %q", s.Inner.Label)
	}

	setValueAtPath(s, "Inner.Count", int64(99))
	if s.Inner.Count != 99 {
		t.Errorf("expected Inner.Count=99, got %d", s.Inner.Count)
	}
}

func TestSetValueAtPathNestedPointer(t *testing.T) {
	s := &outerStruct{}
	// InnerP is nil initially; setValueAtPath should allocate it
	setValueAtPath(s, "InnerP.Label", "via-ptr")
	if s.InnerP == nil {
		t.Fatal("expected InnerP to be allocated")
	}
	if s.InnerP.Label != "via-ptr" {
		t.Errorf("expected InnerP.Label='via-ptr', got %q", s.InnerP.Label)
	}
}

func TestSetValueAtPathPointerFields(t *testing.T) {
	s := &ptrFieldStruct{}

	setValueAtPath(s, "StrPtr", "hello")
	if s.StrPtr == nil || *s.StrPtr != "hello" {
		t.Errorf("expected StrPtr='hello', got %v", s.StrPtr)
	}

	setValueAtPath(s, "Int64Ptr", int64(123))
	if s.Int64Ptr == nil || *s.Int64Ptr != 123 {
		t.Errorf("expected Int64Ptr=123, got %v", s.Int64Ptr)
	}

	setValueAtPath(s, "BoolPtr", true)
	if s.BoolPtr == nil || *s.BoolPtr != true {
		t.Errorf("expected BoolPtr=true, got %v", s.BoolPtr)
	}
}

func TestSetValueAtPathPointerToPointer(t *testing.T) {
	s := &ptrFieldStruct{}
	str := "world"
	setValueAtPath(s, "StrPtr", &str)
	if s.StrPtr == nil || *s.StrPtr != "world" {
		t.Errorf("expected StrPtr='world', got %v", s.StrPtr)
	}
}

func TestSetValueAtPathTypeConversion(t *testing.T) {
	s := &enumStruct{}
	setValueAtPath(s, "Status", "active")
	if s.Status != enumType("active") {
		t.Errorf("expected Status=enumType('active'), got %q", s.Status)
	}
}

func TestSetValueAtPathInt64ToInt32Pointer(t *testing.T) {
	s := &ptrFieldStruct{}
	setValueAtPath(s, "IntPtr", int64(42))
	if s.IntPtr == nil || *s.IntPtr != 42 {
		t.Errorf("expected IntPtr=42, got %v", s.IntPtr)
	}
}

func TestSetValueAtPathEmptyPath(t *testing.T) {
	s := &simpleStruct{Name: "original"}
	setValueAtPath(s, "", "changed")
	if s.Name != "original" {
		t.Errorf("expected Name to remain 'original' with empty path, got %q", s.Name)
	}
}

func TestSetValueAtPathDotOnlyPath(t *testing.T) {
	s := &simpleStruct{Name: "original"}
	setValueAtPath(s, ".", "changed")
	if s.Name != "original" {
		t.Errorf("expected Name to remain 'original' with dot-only path, got %q", s.Name)
	}
}

func TestSetValueAtPathInvalidField(t *testing.T) {
	s := &simpleStruct{}
	// Should not panic on invalid field name
	setValueAtPath(s, "NonExistent", "value")
}

func TestSetValueAtPathNilStruct(t *testing.T) {
	// Should not panic
	setValueAtPath(nil, "Name", "value")
}

func TestSetValueAtPathNilPointerStruct(t *testing.T) {
	var s *simpleStruct
	// Should not panic on nil pointer
	setValueAtPath(s, "Name", "value")
}

func TestSetValueAtPathLeadingDot(t *testing.T) {
	s := &simpleStruct{}
	setValueAtPath(s, ".Name", "stripped")
	if s.Name != "stripped" {
		t.Errorf("expected Name='stripped' with leading dot path, got %q", s.Name)
	}
}

// --- castString tests ---

func TestCastString(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want string
	}{
		{"plain string", "hello", "hello"},
		{"string pointer", aws.String("world"), "world"},
		{"string slice", []string{"a", "b", "c"}, "a,b,c"},
		{"integer", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool", true, "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := castString(tt.in)
			if got != tt.want {
				t.Errorf("castString(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// --- castInt tests ---

func TestCastInt(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    int
		wantErr bool
	}{
		{"string", "42", 42, false},
		{"string pointer", aws.String("99"), 99, false},
		{"int", 7, 7, false},
		{"int pointer", aws.Int(55), 55, false},
		{"int64", int64(100), 100, false},
		{"int64 pointer", aws.Int64(200), 200, false},
		{"invalid string", "abc", 0, true},
		{"invalid string pointer", aws.String("xyz"), 0, true},
		{"unsupported type", 3.14, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := castInt(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("castInt(%v) error=%v, wantErr=%v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("castInt(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

// --- castInt64 tests ---

func TestCastInt64(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    int64
		wantErr bool
	}{
		{"string", "123", 123, false},
		{"int", 50, 50, false},
		{"int pointer", aws.Int(75), 75, false},
		{"int64", int64(999), 999, false},
		{"int64 pointer", aws.Int64(1000), 1000, false},
		{"invalid string", "notanumber", 0, true},
		{"unsupported type", true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := castInt64(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("castInt64(%v) error=%v, wantErr=%v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("castInt64(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

// --- castFloat tests ---

func TestCastFloat(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    float64
		wantErr bool
	}{
		{"string", "3.14", 3.14, false},
		{"float32", float32(1.5), 1.5, false},
		{"float64", 2.71, 2.71, false},
		{"float64 pointer", aws.Float64(9.99), 9.99, false},
		{"int", 5, 5.0, false},
		{"int64", int64(10), 10.0, false},
		{"invalid string", "abc", 0, true},
		{"unsupported type", true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := castFloat(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("castFloat(%v) error=%v, wantErr=%v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("castFloat(%v) = %f, want %f", tt.in, got, tt.want)
			}
		})
	}
}

// --- castBool tests ---

func TestCastBool(t *testing.T) {
	tests := []struct {
		name    string
		in      interface{}
		want    bool
		wantErr bool
	}{
		{"true string", "true", true, false},
		{"false string", "false", false, false},
		{"1 string", "1", true, false},
		{"0 string", "0", false, false},
		{"bool true", true, true, false},
		{"bool false", false, false, false},
		{"bool pointer true", aws.Bool(true), true, false},
		{"bool pointer false", aws.Bool(false), false, false},
		{"invalid string", "notbool", false, true},
		{"unsupported type", 42, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := castBool(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("castBool(%v) error=%v, wantErr=%v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("castBool(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// --- castStringSlice tests ---

func TestCastStringSlice(t *testing.T) {
	t.Run("single string", func(t *testing.T) {
		got := castStringSlice("hello")
		if len(got) != 1 || got[0] != "hello" {
			t.Errorf("expected [hello], got %v", got)
		}
	})

	t.Run("string pointer", func(t *testing.T) {
		got := castStringSlice(aws.String("world"))
		if len(got) != 1 || got[0] != "world" {
			t.Errorf("expected [world], got %v", got)
		}
	})

	t.Run("string slice", func(t *testing.T) {
		got := castStringSlice([]string{"a", "b"})
		if len(got) != 2 || got[0] != "a" || got[1] != "b" {
			t.Errorf("expected [a b], got %v", got)
		}
	})

	t.Run("string pointer slice", func(t *testing.T) {
		got := castStringSlice([]*string{aws.String("x"), aws.String("y")})
		if len(got) != 2 || got[0] != "x" || got[1] != "y" {
			t.Errorf("expected [x y], got %v", got)
		}
	})

	t.Run("interface slice", func(t *testing.T) {
		got := castStringSlice([]interface{}{"str", aws.String("ptr"), 123})
		if len(got) != 3 || got[0] != "str" || got[1] != "ptr" || got[2] != "123" {
			t.Errorf("expected [str ptr 123], got %v", got)
		}
	})

	t.Run("other type", func(t *testing.T) {
		got := castStringSlice(42)
		if len(got) != 1 || got[0] != "42" {
			t.Errorf("expected [42], got %v", got)
		}
	})
}

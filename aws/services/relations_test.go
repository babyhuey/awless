package awsservices

import (
	"reflect"
	"testing"
)

type innerStruct struct {
	InnerField string
}

type testStruct struct {
	Name    string
	Value   *string
	Nested  innerStruct
	PtrNest *innerStruct
}

func strPtr(s string) *string { return &s }

func TestValueAtPath(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		path    string
		want    any
		wantErr bool
	}{
		{
			name:  "simple string field",
			input: &testStruct{Name: "hello"},
			path:  "Name",
			want:  "hello",
		},
		{
			name:  "pointer string field",
			input: &testStruct{Value: strPtr("world")},
			path:  "Value",
			want:  strPtr("world"),
		},
		{
			name:  "nil pointer string field",
			input: &testStruct{Value: nil},
			path:  "Value",
			want:  (*string)(nil),
		},
		{
			name:  "nested struct field",
			input: &testStruct{Nested: innerStruct{InnerField: "deep"}},
			path:  "Nested.InnerField",
			want:  "deep",
		},
		{
			name:  "pointer nested struct field",
			input: &testStruct{PtrNest: &innerStruct{InnerField: "ptrdeep"}},
			path:  "PtrNest.InnerField",
			want:  "ptrdeep",
		},
		{
			name:  "nil pointer nested struct returns nil",
			input: &testStruct{PtrNest: nil},
			path:  "PtrNest.InnerField",
			want:  nil,
		},
		{
			name:    "missing field returns error",
			input:   &testStruct{},
			path:    "NonExistent",
			wantErr: true,
		},
		{
			name:    "missing nested field returns error",
			input:   &testStruct{Nested: innerStruct{}},
			path:    "Nested.NonExistent",
			wantErr: true,
		},
		{
			name:    "non-struct input returns error for nested path",
			input:   "not a struct",
			path:    "Field",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueAtPath(tt.input, tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("valueAtPath() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("valueAtPath() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valueAtPath() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestVerifyValidStructField(t *testing.T) {
	type sample struct {
		FieldA string
		FieldB int
		Items  []string
	}

	tests := []struct {
		name      string
		input     any
		fieldName string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "valid string field",
			input:     &sample{FieldA: "test"},
			fieldName: "FieldA",
			wantValid: true,
		},
		{
			name:      "valid int field",
			input:     &sample{FieldB: 42},
			fieldName: "FieldB",
			wantValid: true,
		},
		{
			name:      "valid slice field",
			input:     &sample{Items: []string{"a", "b"}},
			fieldName: "Items",
			wantValid: true,
		},
		{
			name:      "invalid field name",
			input:     &sample{},
			fieldName: "NonExistent",
			wantErr:   true,
		},
		{
			name:      "non-pointer input",
			input:     sample{},
			fieldName: "FieldA",
			wantErr:   true,
		},
		{
			name:      "pointer to non-struct",
			input:     strPtr("hello"),
			fieldName: "FieldA",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := verifyValidStructField(tt.input, tt.fieldName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("verifyValidStructField() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("verifyValidStructField() unexpected error: %v", err)
			}
			if tt.wantValid && !val.IsValid() {
				t.Errorf("verifyValidStructField() returned invalid value, expected valid")
			}
		})
	}
}

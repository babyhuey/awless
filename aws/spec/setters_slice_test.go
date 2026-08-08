package awsspec

import (
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Covers the SDK v1 -> v2 slice shape change. v1 modeled string lists as
// []*string and v2 models them as []string, and the setters still build
// []*string, so the conversion in setValueAtPath is what keeps 35 params working.

func TestConvertSliceDerefsPointers(t *testing.T) {
	src := reflect.ValueOf([]*string{awssdk.String("a"), awssdk.String("b")})

	got, ok := convertSlice(src, reflect.TypeOf([]string{}))
	if !ok {
		t.Fatal("expected []*string to convert to []string")
	}
	if want := []string{"a", "b"}; !reflect.DeepEqual(got.Interface(), want) {
		t.Errorf("got %#v, want %#v", got.Interface(), want)
	}
}

func TestConvertSliceTakesAddresses(t *testing.T) {
	src := reflect.ValueOf([]string{"x"})

	got, ok := convertSlice(src, reflect.TypeOf([]*string{}))
	if !ok {
		t.Fatal("expected []string to convert to []*string")
	}
	out := got.Interface().([]*string)
	if len(out) != 1 || awssdk.ToString(out[0]) != "x" {
		t.Errorf("got %#v, want one element \"x\"", out)
	}
}

// A nil element cannot be represented in a non-pointer destination, so it must
// become the zero value rather than panicking.
func TestConvertSliceNilElementBecomesZero(t *testing.T) {
	src := reflect.ValueOf([]*string{awssdk.String("a"), nil})

	got, ok := convertSlice(src, reflect.TypeOf([]string{}))
	if !ok {
		t.Fatal("expected conversion to succeed")
	}
	if want := []string{"a", ""}; !reflect.DeepEqual(got.Interface(), want) {
		t.Errorf("got %#v, want %#v", got.Interface(), want)
	}
}

func TestConvertSliceRejectsUnrelatedTypes(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  any
		dest any
	}{
		{"not a slice", "scalar", []string{}},
		{"dest not a slice", []string{"a"}, ""},
		{"incompatible elements", []*string{}, []int{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := convertSlice(reflect.ValueOf(tc.src), reflect.TypeOf(tc.dest)); ok {
				t.Error("expected conversion to be refused")
			}
		})
	}
}

// setFieldWithType is the dispatch every generated command relies on, so cover
// the common awsType conversions against real SDK input structs.
func TestSetFieldWithTypeConversions(t *testing.T) {
	t.Run("awsstringslice into an SDK v2 []string", func(t *testing.T) {
		in := &ec2.CreateTagsInput{}
		if err := setFieldWithType([]any{"i-1", "i-2"}, in, "Resources", awsstringslice); err != nil {
			t.Fatal(err)
		}
		if want := []string{"i-1", "i-2"}; !reflect.DeepEqual(in.Resources, want) {
			t.Errorf("got %#v, want %#v", in.Resources, want)
		}
	})

	t.Run("awsstr", func(t *testing.T) {
		in := &ec2.CreateVpcInput{}
		if err := setFieldWithType("10.0.0.0/16", in, "CidrBlock", awsstr); err != nil {
			t.Fatal(err)
		}
		if got := awssdk.ToString(in.CidrBlock); got != "10.0.0.0/16" {
			t.Errorf("got %q, want 10.0.0.0/16", got)
		}
	})

	t.Run("awsbool", func(t *testing.T) {
		in := &ec2.RunInstancesInput{}
		if err := setFieldWithType("true", in, "EbsOptimized", awsbool); err != nil {
			t.Fatal(err)
		}
		if !awssdk.ToBool(in.EbsOptimized) {
			t.Error("expected EbsOptimized to be true")
		}
	})

	t.Run("awscsvstr joins", func(t *testing.T) {
		in := &ec2.CreateVpcInput{}
		if err := setFieldWithType([]any{"a", "b", "c"}, in, "CidrBlock", awscsvstr); err != nil {
			t.Fatal(err)
		}
		if got := awssdk.ToString(in.CidrBlock); got != "a,b,c" {
			t.Errorf("got %q, want a,b,c", got)
		}
	})

	t.Run("awsdimensionslice parses key:value", func(t *testing.T) {
		in := &cloudwatch.PutMetricAlarmInput{}
		if err := setFieldWithType([]any{"InstanceId:i-123"}, in, "Dimensions", awsdimensionslice); err != nil {
			t.Fatal(err)
		}
		if len(in.Dimensions) != 1 {
			t.Fatalf("expected 1 dimension, got %d", len(in.Dimensions))
		}
		if got := awssdk.ToString(in.Dimensions[0].Name); got != "InstanceId" {
			t.Errorf("name: got %q, want InstanceId", got)
		}
		if got := awssdk.ToString(in.Dimensions[0].Value); got != "i-123" {
			t.Errorf("value: got %q, want i-123", got)
		}
	})

	t.Run("awsdimensionslice rejects a malformed entry", func(t *testing.T) {
		in := &cloudwatch.PutMetricAlarmInput{}
		if err := setFieldWithType([]any{"no-colon"}, in, "Dimensions", awsdimensionslice); err == nil {
			t.Error("expected an error for a dimension without a colon")
		}
	})

	// A nil value must be a no-op rather than clearing the field or panicking.
	t.Run("nil value is a no-op", func(t *testing.T) {
		in := &ec2.CreateVpcInput{CidrBlock: awssdk.String("keep")}
		if err := setFieldWithType(nil, in, "CidrBlock", awsstr); err != nil {
			t.Fatal(err)
		}
		if got := awssdk.ToString(in.CidrBlock); got != "keep" {
			t.Errorf("got %q, want the field left alone", got)
		}
	})

	// An unknown field is silently ignored rather than erroring. That is safe
	// because awsName tags are generated, not hand-written, so a missing field
	// means the SDK struct changed shape rather than a user typo. Asserted so the
	// behavior is deliberate and a change to it is visible.
	t.Run("unknown field is ignored, not an error", func(t *testing.T) {
		in := &ec2.CreateVpcInput{}
		if err := setFieldWithType("v", in, "NoSuchField", awsstr); err != nil {
			t.Errorf("expected a silent no-op, got %s", err)
		}
		if in.CidrBlock != nil {
			t.Error("expected no field to be set")
		}
	})
}

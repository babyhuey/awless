/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsspec

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	applicationautoscalingtypes "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func TestGoTemplatingInUserdata(t *testing.T) {
	text := []byte("file content {{ .name }}")
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	finfo, _ := f.Stat()
	err = os.WriteFile(f.Name(), text, finfo.Mode().Perm())
	if err != nil {
		t.Fatal(f)
	}

	awsparams := &ec2.RunInstancesInput{}

	err = setFieldWithType(context.Background(), f.Name(), awsparams, "UserData", awsuserdatatobase64, map[string]string{"name": "johndoe"})
	if err != nil {
		t.Fatal(err)
	}
	expText := []byte("file content johndoe")
	if got, want := aws.ToString(awsparams.UserData), base64.StdEncoding.EncodeToString(expText); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSetFieldWithTypeAWSFile(t *testing.T) {
	text := []byte("file content")
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	finfo, _ := f.Stat()
	err = os.WriteFile(f.Name(), text, finfo.Mode().Perm())
	if err != nil {
		t.Fatal(f)
	}

	awsparams := &ec2.RunInstancesInput{}

	err = setFieldWithType(context.Background(), f.Name(), awsparams, "UserData", awsuserdatatobase64)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := aws.ToString(awsparams.UserData), base64.StdEncoding.EncodeToString(text); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	functionInput := &lambda.CreateFunctionInput{}

	err = setFieldWithType(context.Background(), f.Name(), functionInput, "Code.ZipFile", awsfiletobyteslice)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(functionInput.Code.ZipFile), string(text); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	stackInput := &cloudformation.CreateStackInput{}

	err = setFieldWithType(context.Background(), f.Name(), stackInput, "TemplateBody", awsfiletostring)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := aws.ToString(stackInput.TemplateBody), string(text); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSetFieldsOnAwsStruct(t *testing.T) {
	awsparams := &ec2.RunInstancesInput{}

	err := setFieldWithType(context.Background(), "ami", awsparams, "ImageId", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	err = setFieldWithType(context.Background(), "t2.micro", awsparams, "InstanceType", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	err = setFieldWithType(context.Background(), "5", awsparams, "MaxCount", awsint64)
	if err != nil {
		t.Fatal(err)
	}
	err = setFieldWithType(context.Background(), 3, awsparams, "MinCount", awsint64)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := aws.ToString(awsparams.ImageId), "ami"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := string(awsparams.InstanceType), "t2.micro"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := aws.ToInt32(awsparams.MaxCount), int32(5); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := aws.ToInt32(awsparams.MinCount), int32(3); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestSetFieldWithMultiType(t *testing.T) {
	subject := struct {
		Field               string
		IntField            int
		FloatField          *float64
		BoolPointerField    *bool
		BoolField           bool
		StringArrayField    []*string
		Int64ArrayField     []*int64
		BooleanValueField   *ec2types.AttributeBooleanValue
		StringValueField    *ec2types.AttributeValue
		DimensionSliceField []cloudwatchtypes.Dimension
		KeyValueSliceField  []ecstypes.KeyValuePair
		StructAttribute     struct {
			Str  *string
			Bool *bool
		}
		SliceStructPointerAttribute []*struct {
			Str1, Str2 *string
			Integer    *int64
		}
		MapAttribute          map[string]*string
		EmptyMapAttribute     map[string]*string
		ParameterList         []cloudformationtypes.Parameter
		PortMappings          []ecstypes.PortMapping
		SubnetMappings        []elbv2types.SubnetMapping
		LoadBalancerListeners []elbtypes.Listener
		StepAdjustments       []applicationautoscalingtypes.StepAdjustment
		CSVString             *string
		SixDigitsString       *string
		ByteSlice             []byte
	}{Field: "initial", MapAttribute: map[string]*string{"test": aws.String("1234")}}

	err := setFieldWithType(context.Background(), "expected", &subject, "Field", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.Field, "expected"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), 5, &subject, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), 42.21, &subject, "FloatField", awsfloat)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.FloatField, 42.21; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), "5", &subject, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), nil, &subject, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), "first", &subject, "StringArrayField", awsstringslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.StringArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.ToString(subject.StringArrayField[0]), "first"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), []string{"one", "two", "three"}, &subject, "StringArrayField", awsstringslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.StringArrayField), 3; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.ToString(subject.StringArrayField[0]), "one"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := aws.ToString(subject.StringArrayField[1]), "two"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := aws.ToString(subject.StringArrayField[2]), "three"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), []any{"four", "five"}, &subject, "StringArrayField", awsstringslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.StringArrayField), 2; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.ToString(subject.StringArrayField[0]), "four"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := aws.ToString(subject.StringArrayField[1]), "five"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), int64(321), &subject, "Int64ArrayField", awsint64slice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.Int64ArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.ToInt64(subject.Int64ArrayField[0]), int64(321); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), 567, &subject, "Int64ArrayField", awsint64slice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.Int64ArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.ToInt64(subject.Int64ArrayField[0]), int64(567); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), "any", nil, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}

	err = setFieldWithType(context.Background(), "true", &subject, "BooleanValueField", awsboolattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aws.ToBool(subject.BooleanValueField.Value), true; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}
	err = setFieldWithType(context.Background(), nil, &subject, "BooleanValueField", awsboolattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aws.ToBool(subject.BooleanValueField.Value), true; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}
	err = setFieldWithType(context.Background(), false, &subject, "BooleanValueField", awsboolattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aws.ToBool(subject.BooleanValueField.Value), false; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}

	err = setFieldWithType(context.Background(), "true", &subject, "BooleanValueField", awsbool)
	if err == nil {
		t.Fatalf("expected error got nil")
	}
	if got, want := err.Error(), "value of type bool"; !strings.Contains(got, want) || !strings.Contains(got, "types.AttributeBooleanValue") {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType(context.Background(), "abcd", &subject, "StringValueField", awsstringattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aws.ToString(subject.StringValueField.Value), "abcd"; got != want {
		t.Fatalf("len: got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), nil, &subject, "StringValueField", awsstringattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aws.ToString(subject.StringValueField.Value), "abcd"; got != want {
		t.Fatalf("len: got %s, want %s", got, want)
	}

	err = setFieldWithType(context.Background(), true, &subject, "BoolField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.BoolField, true; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), false, &subject, "BoolField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.BoolField, false; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), "true", &subject, "BoolPointerField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.BoolPointerField, true; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), false, &subject, "BoolPointerField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.BoolPointerField, false; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), "fieldValue", &subject, "StructAttribute.Str", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.StructAttribute.Str, "fieldValue"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), []string{"one", "two", "three"}, &subject, "StructAttribute.Str", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.StructAttribute.Str, "one,two,three"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), "true", &subject, "StructAttribute.Bool", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.StructAttribute.Bool, true; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(context.Background(), "abc", &subject, "MapAttribute[Field1]", awsstringpointermap)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.MapAttribute), 1+1; got != want { //First "test" key + Field1
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.MapAttribute["Field1"], "abc"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), "def", &subject, "MapAttribute[Field2]", awsstringpointermap)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.MapAttribute), 1+2; got != want { //First "test" key + Field1 and Field2
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.MapAttribute["Field1"], "abc"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := *subject.MapAttribute["Field2"], "def"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), "abcd", &subject, "EmptyMapAttribute[Field1]", awsstringpointermap)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.EmptyMapAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.EmptyMapAttribute["Field1"], "abcd"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(context.Background(), "tata", &subject, "SliceStructPointerAttribute[0]Str1", awsslicestruct)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.SliceStructPointerAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.SliceStructPointerAttribute[0].Str1, "tata"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), "toto", &subject, "SliceStructPointerAttribute[0]Str2", awsslicestruct)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.SliceStructPointerAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.SliceStructPointerAttribute[0].Str2, "toto"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), 10, &subject, "SliceStructPointerAttribute[0]Integer", awsslicestructint64)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.SliceStructPointerAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.SliceStructPointerAttribute[0].Integer, int64(10); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	err = setFieldWithType(context.Background(), "key:value", &subject, "DimensionSliceField", awsdimensionslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.DimensionSliceField), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.DimensionSliceField[0].Name, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.DimensionSliceField[0].Value, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), []string{"key:value", "key1:value1:with:"}, &subject, "DimensionSliceField", awsdimensionslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.DimensionSliceField), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.DimensionSliceField[0].Name, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.DimensionSliceField[0].Value, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.DimensionSliceField[1].Name, "key1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.DimensionSliceField[1].Value, "value1:with:"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType(context.Background(), []string{"key:value", "key1:value1:with:"}, &subject, "ParameterList", awsparameterslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.ParameterList), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.ParameterList[0].ParameterKey, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.ParameterList[0].ParameterValue, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.ParameterList[1].ParameterKey, "key1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.ParameterList[1].ParameterValue, "value1:with:"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), []string{"key:value", "key1:value1:with:"}, &subject, "KeyValueSliceField", awsecskeyvalue)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.KeyValueSliceField), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.KeyValueSliceField[0].Name, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.KeyValueSliceField[0].Value, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.KeyValueSliceField[1].Name, "key1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.KeyValueSliceField[1].Value, "value1:with:"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType(context.Background(), []string{"80:8080", "8082", "1234:8083/udp"}, &subject, "PortMappings", awsportmappings)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.PortMappings), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.PortMappings[0].HostPort, int32(80); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.PortMappings[0].ContainerPort, int32(8080); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.PortMappings[1].ContainerPort, int32(8082); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.PortMappings[2].HostPort, int32(1234); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.PortMappings[2].ContainerPort, int32(8083); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := string(subject.PortMappings[2].Protocol), "udp"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType(context.Background(), []string{"subnet-123:eipalloc-123", "subnet-456:eipalloc-456"}, &subject, "SubnetMappings", awssubnetmappings)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.SubnetMappings), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.SubnetMappings[0].SubnetId, "subnet-123"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.SubnetMappings[0].AllocationId, "eipalloc-123"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.SubnetMappings[1].SubnetId, "subnet-456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.SubnetMappings[1].AllocationId, "eipalloc-456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType(context.Background(), []string{"HTTP:80:UDP:8080", "HTTPS:443:TCP:12345"}, &subject, "LoadBalancerListeners", awsclassicloadblisteners)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.LoadBalancerListeners), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.LoadBalancerListeners[0].Protocol, "HTTP"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := subject.LoadBalancerListeners[0].LoadBalancerPort, int32(80); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.LoadBalancerListeners[0].InstanceProtocol, "UDP"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.LoadBalancerListeners[0].InstancePort, int32(8080); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.LoadBalancerListeners[1].Protocol, "HTTPS"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := subject.LoadBalancerListeners[1].LoadBalancerPort, int32(443); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.LoadBalancerListeners[1].InstanceProtocol, "TCP"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *subject.LoadBalancerListeners[1].InstancePort, int32(12345); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	err = setFieldWithType(context.Background(), []string{"0:0.25:-1", "0.75:1:+1"}, &subject, "StepAdjustments", awsstepadjustments)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(subject.StepAdjustments), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.StepAdjustments[0].MetricIntervalLowerBound, float64(0); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *subject.StepAdjustments[0].MetricIntervalUpperBound, float64(0.25); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *subject.StepAdjustments[0].ScalingAdjustment, int32(-1); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *subject.StepAdjustments[1].MetricIntervalLowerBound, float64(0.75); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *subject.StepAdjustments[1].MetricIntervalUpperBound, float64(1); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *subject.StepAdjustments[1].ScalingAdjustment, int32(+1); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	err = setFieldWithType(context.Background(), []any{"abcdef", "ghijk"}, &subject, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.CSVString, "abcdef,ghijk"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), []string{"abcdef", "ghijk"}, &subject, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.CSVString, "abcdef,ghijk"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), "abcdef", &subject, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.CSVString, "abcdef"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), "abcdef,ghijk", &subject, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.CSVString, "abcdef,ghijk"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), "123456", &subject, "SixDigitsString", aws6digitsstring)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.SixDigitsString, "123456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), "2345", &subject, "SixDigitsString", aws6digitsstring)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *subject.SixDigitsString, "002345"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(context.Background(), []byte("hello"), &subject, "ByteSlice", awsbyteslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := subject.ByteSlice, []byte("hello"); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

type TestStruct struct {
	FieldStringRequired *string   `awsName:"CloudStringRequired" awsType:"awsstr" templateName:"fstringrequired"`
	FieldString         *string   `awsName:"CloudString" awsType:"awsstr" templateName:"fstring"`
	FieldInt64          *int64    `awsName:"CloudInt64" awsType:"awsint64" templateName:"fint"`
	FieldBool           *bool     `awsName:"Embedded.CloudBool" awsType:"awsbool" templateName:"fbool"`
	FieldStringSlice    []*string `awsName:"CloudStringSlice" awsType:"awsstringslice" templateName:"fstrslice"`
	NilField            *string   `awsName:"CloudNilString" awsType:"awsstr" templateName:"fnilstring"`
	MultiCloudField     *int64    `awsName:"CloudField1,CloudField2" awsType:"awsint64" templateName:"fmultistring"`
}

func (ts *TestStruct) ValidateFieldStringRequired() (err error) {
	if len(*ts.FieldStringRequired) == 0 {
		err = errors.New("fstringrequired should not be empty")
	}
	return
}

func (ts *TestStruct) ValidateFieldString() (err error) {
	if len(*ts.FieldString) != 10 {
		err = errors.New("fstring should be 10 chars")
	}
	return
}

func (ts *TestStruct) ValidateFieldInt64() (err error) {
	if *ts.FieldInt64 > 10 {
		err = errors.New("fint should not exceed 10")
	}
	return
}

func TestStructDynamicSetter(t *testing.T) {
	params := map[string]any{
		"fstringrequired": "jdoe",
		"fint":            "345",
		"fbool":           "true",
		"fstrslice":       []any{"one", "two", 3},
	}

	in := &TestStruct{}
	err := structSetter(in, params)
	if err != nil {
		t.Fatal(err)
	}

	exp := &TestStruct{
		FieldStringRequired: aws.String("jdoe"),
		FieldInt64:          aws.Int64(345),
		FieldBool:           aws.Bool(true),
		FieldStringSlice:    aws.StringSlice([]string{"one", "two", "3"}),
	}

	if got, want := in, exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %#v\n\nwant %#v\n", got, want)
	}
}

func TestStructInjector(t *testing.T) {
	in := &TestStruct{
		FieldStringRequired: aws.String("jdoe"),
		FieldInt64:          aws.Int64(345),
		FieldBool:           aws.Bool(true),
		FieldStringSlice:    aws.StringSlice([]string{"one", "two", "3"}),
		MultiCloudField:     aws.Int64(12345),
	}

	type embStruct struct {
		CloudBool *bool
	}
	type outStruct struct {
		CloudStringRequired      *string
		CloudInt64               *int64
		Embedded                 *embStruct
		CloudStringSlice         []*string
		CloudField1, CloudField2 *int64
	}

	out := new(outStruct)

	err := structInjector(in, out, nil)
	if err != nil {
		t.Fatal(err)
	}

	exp := &outStruct{
		CloudStringRequired: aws.String("jdoe"),
		CloudInt64:          aws.Int64(345),
		Embedded:            &embStruct{CloudBool: aws.Bool(true)},
		CloudStringSlice:    aws.StringSlice([]string{"one", "two", "3"}),
		CloudField1:         aws.Int64(12345),
		CloudField2:         aws.Int64(12345),
	}

	if got, want := out, exp; !reflect.DeepEqual(got, want) {
		// pretty.Print(got)
		// pretty.Print(want)
		t.Fatalf("\ngot %#v\n\nwant %#v\n", got, want)
	}
}

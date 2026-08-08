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
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	gotemplate "text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	applicationautoscalingtypes "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/bootswithdefer/awless/logger"
)

const (
	awsstr                   = "awsstr"
	awsint                   = "awsint"
	awsint64                 = "awsint64"
	awsfloat                 = "awsfloat"
	awsbool                  = "awsbool"
	awsboolattribute         = "awsboolattribute"
	awsstringattribute       = "awsstringattribute"
	awsint64slice            = "awsint64slice"
	awsstringslice           = "awsstringslice"
	awsstringpointermap      = "awsstringpointermap"
	awsslicestruct           = "awsslicestruct"
	awsslicestructint64      = "awsslicestructint64"
	awsuserdatatobase64      = "awsuserdatatobase64"
	awsfiletobyteslice       = "awsfiletobyteslice"
	awsfiletostring          = "awsfiletostring"
	awsdimensionslice        = "awsdimensionslice"
	awsparameterslice        = "awsparameterslice"
	awsecskeyvalue           = "awsecskeyvalue"
	awsportmappings          = "awsportmappings"
	awssubnetmappings        = "awssubnetmappings"
	awsclassicloadblisteners = "awsclassicloadblisteners"
	awsstepadjustments       = "awsstepadjustments"
	awscsvstr                = "awscsvstr"
	aws6digitsstring         = "aws6digitsstring"
	awsbyteslice             = "awsbyteslice"
	awstagslice              = "awstagslice"
	awsalarmrollbacktriggers = "awsalarmrollbacktriggers"
)

var (
	mapAttributeRegex = regexp.MustCompile(`(.+)\[(.+)\].*`)
	sliceStructRegex  = regexp.MustCompile(`(.+)\[0\](.*)`)
)

type setter struct {
	val       any
	fieldPath string
	fieldType string
}

func (s setter) set(i any) error {
	return setFieldWithType(s.val, i, s.fieldPath, s.fieldType)
}

func setFieldWithType(v, i any, fieldPath string, destType string, interfs ...any) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("set field %s for %T object: %s", fieldPath, i, e)
		}
	}()
	if v == nil || i == nil {
		return nil
	}
	switch destType {
	case awsstr:
		v = castString(v)
	case awsint64:
		v, err = castInt64(v)
		if err != nil {
			return
		}
	case awsint:
		v, err = castInt(v)
		if err != nil {
			return
		}
	case aws6digitsstring:
		v, err = castInt(v)
		if err != nil {
			return
		}
		v = fmt.Sprintf("%06d", v)
	case awsfloat:
		v, err = castFloat(v)
		if err != nil {
			return
		}
	case awsbool:
		v, err = castBool(v)
		if err != nil {
			return
		}
	case awsstringslice:
		v = castStringPointerSlice(v)
	case awsbyteslice:
	case awscsvstr:
		v = strings.Join(castStringSlice(v), ",")
	case awsdimensionslice:
		if dimensions, isDim := v.([]cloudwatchtypes.Dimension); isDim {
			v = dimensions
		} else {
			dimensions = []cloudwatchtypes.Dimension{}
			sl := castStringSlice(v)
			for _, s := range sl {
				splits := strings.SplitN(s, ":", 2)
				if len(splits) != 2 {
					return fmt.Errorf("invalid dimension '%s', expected 'key:value'", s)
				}
				dimensions = append(dimensions, cloudwatchtypes.Dimension{Name: aws.String(splits[0]), Value: aws.String(splits[1])})
				v = dimensions
			}
		}
	case awsecskeyvalue:
		sl := castStringSlice(v)
		var keyvalues []ecstypes.KeyValuePair
		for _, s := range sl {
			splits := strings.SplitN(s, ":", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid keyvalue '%s', expected 'key:value'", s)
			}
			keyvalues = append(keyvalues, ecstypes.KeyValuePair{Name: aws.String(splits[0]), Value: aws.String(splits[1])})
		}
		v = keyvalues
	case awsparameterslice:
		sl := castStringSlice(v)
		var parameters []cloudformationtypes.Parameter
		for _, s := range sl {
			splits := strings.SplitN(s, ":", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid parameter '%s', expected 'key:value'", s)
			}
			parameters = append(parameters, cloudformationtypes.Parameter{ParameterKey: aws.String(splits[0]), ParameterValue: aws.String(splits[1])})
		}
		v = parameters
	case awssubnetmappings:
		sl := castStringSlice(v)
		var subnetMappings []elbv2types.SubnetMapping
		for i, s := range sl {
			splits := strings.Split(s, ":")
			if len(splits) != 2 {
				return fmt.Errorf("invalid element %d in subnet mapping %v, expect format [subnet-123:eipalloc-321, subnet-234:eipalloc-678, ...]", i+1, splits)
			}
			subnetMappings = append(subnetMappings, elbv2types.SubnetMapping{SubnetId: aws.String(splits[0]), AllocationId: aws.String(splits[1])})
		}
		v = subnetMappings
	case awsclassicloadblisteners:
		var listeners []elbtypes.Listener
		for _, s := range castStringSlice(v) {
			splits := strings.Split(s, ":")
			if len(splits) != 4 {
				return fmt.Errorf("missing value in listeners param '%s', expect format like HTTP:80:HTTP:80", splits)
			}
			loadbPort, err := strconv.ParseInt(splits[1], 10, 64)
			if err != nil {
				return fmt.Errorf("expecting numerical port value for loadbalancer port in '%s', (expect format like HTTP:80:HTTP:80)", splits)
			}
			instancePort, err := strconv.ParseInt(splits[3], 10, 64)
			if err != nil {
				return fmt.Errorf("expecting numerical port value for instance port in '%s', (expect format like HTTP:80:HTTP:80)", splits)
			}
			listeners = append(listeners, elbtypes.Listener{
				Protocol:         aws.String(splits[0]),
				LoadBalancerPort: int32(loadbPort),
				InstanceProtocol: aws.String(splits[2]),
				InstancePort:     aws.Int32(int32(instancePort)),
			})
		}
		v = listeners
	case awsportmappings:
		sl := castStringSlice(v)
		var portMappings []ecstypes.PortMapping
		for _, s := range sl {
			portMapping := &ecstypes.PortMapping{}
			if strings.Contains(s, "-") {
				return fmt.Errorf("invalid port mapping '%s', AWS do not support portrange (from-to)", s)
			}
			var protocol string
			if strings.Contains(s, "/") {
				splits := strings.Split(s, "/")
				protocol = splits[1]
				if protocol != "tcp" && protocol != "udp" {
					return fmt.Errorf("invalid port mapping '%s', invalid protocol, expect tcp or udp, got %s", s, protocol)
				}
				s = strings.TrimRight(s, "/"+protocol)
				portMapping.Protocol = ecstypes.TransportProtocol(protocol)
			}
			splits := strings.Split(s, ":")
			switch len(splits) {
			case 1:
				containerPort, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
				}
				portMapping.ContainerPort = aws.Int32(int32(containerPort))
			case 2:
				hostPort, err := strconv.ParseInt(splits[0], 10, 32)
				if err != nil {
					return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
				}
				containerPort, err := strconv.ParseInt(splits[1], 10, 32)
				if err != nil {
					return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
				}
				portMapping.HostPort = aws.Int32(int32(hostPort))
				portMapping.ContainerPort = aws.Int32(int32(containerPort))
			default:
				return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
			}

			portMappings = append(portMappings, *portMapping)
		}
		v = portMappings
	case awsstepadjustments:
		sl := castStringSlice(v)
		var stepAdjustments []applicationautoscalingtypes.StepAdjustment
		for _, s := range sl {
			splits := strings.Split(s, ":")
			if len(splits) != 3 {
				return fmt.Errorf("invalid step adjustment '%s', expect from:to:scaling-adjustment", s)
			}
			stepAdjustment := &applicationautoscalingtypes.StepAdjustment{}
			if splits[0] != "" {
				lower, err := strconv.ParseFloat(splits[0], 64)
				if err != nil {
					return fmt.Errorf("invalid from '%s' in step adjustment '%s', expect from:to:scaling-adjustment", splits[0], s)
				}
				stepAdjustment.MetricIntervalLowerBound = aws.Float64(lower)
			}
			if splits[1] != "" {
				upper, err := strconv.ParseFloat(splits[1], 64)
				if err != nil {
					return fmt.Errorf("invalid to '%s' in step adjustment '%s', expect from:to:scaling-adjustment", splits[1], s)
				}
				stepAdjustment.MetricIntervalUpperBound = aws.Float64(upper)
			}
			adjustment, err := strconv.ParseInt(splits[2], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid adjustment-adjustment '%s' in step adjustmentstep adjustment '%s', expect from:to:scaling-adjustment", splits[2], s)
			}
			stepAdjustment.ScalingAdjustment = aws.Int32(int32(adjustment))
			stepAdjustments = append(stepAdjustments, *stepAdjustment)
		}
		v = stepAdjustments
	case awsuserdatatobase64:
		var tplData any
		if len(interfs) > 0 {
			tplData = interfs[0]
		}
		v, err = userDataContentAsBase64(v, tplData)
		if err != nil {
			return err
		}
	case awsfiletobyteslice:
		v, err = os.ReadFile(castString(v))
		if err != nil {
			return err
		}
	case awsfiletostring:
		var b []byte
		b, err = os.ReadFile(castString(v))
		if err != nil {
			return err
		}
		v = string(b)
	case awsint64slice:
		var awsint int64
		awsint, err = castInt64(v)
		if err != nil {
			return
		}
		v = []*int64{&awsint}
	case awsboolattribute:
		var b bool
		b, err = castBool(v)
		if err != nil {
			return
		}
		v = &ec2types.AttributeBooleanValue{Value: &b}
	case awsstringattribute:
		str := castString(v)
		v = &ec2types.AttributeValue{Value: &str}
	case awsstringpointermap:
		matches := mapAttributeRegex.FindStringSubmatch(fieldPath)
		if len(matches) < 2 {
			err = fmt.Errorf("set field awsstringmap: path %s does not start with mymap[key]", fieldPath)
			return
		}
		strcr := reflect.Indirect(reflect.ValueOf(i))
		if strcr.Kind() != reflect.Struct {
			err = fmt.Errorf("set field awsstringmap: %T is not a struct, but a %s", i, strcr.Kind())
			return
		}
		field := strcr.FieldByName(matches[1])
		if field.Kind() != reflect.Map {
			err = fmt.Errorf("set field awsstringmap: field %s is not a map, but a %s", matches[0], field.Kind())
			return
		}
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}
		str := castString(v)
		field.SetMapIndex(reflect.ValueOf(matches[2]), reflect.ValueOf(&str))
		return nil
	case awsslicestruct, awsslicestructint64:
		if destType == awsslicestructint64 {
			v, err = castInt64(v)
			if err != nil {
				return
			}
		}
		matches := sliceStructRegex.FindStringSubmatch(fieldPath)
		if len(matches) < 2 {
			err = fmt.Errorf("set field awsslicestruct: path %s does not start with slice[0]", fieldPath)
			return
		}
		strcr := reflect.Indirect(reflect.ValueOf(i))
		if strcr.Kind() != reflect.Struct {
			err = fmt.Errorf("set field awsslicestruct: %T is not a struct, but a %s", i, strcr.Kind())
			return
		}
		sliceField := strcr.FieldByName(matches[1])
		if sliceField.Kind() != reflect.Slice {
			err = fmt.Errorf("set field awsslicestruct: field %s is not a slice, but a %s", matches[0], sliceField.Kind())
			return
		}
		var elemToSet reflect.Value
		if sliceField.Len() > 0 {
			elemToSet = sliceField.Index(0)
		} else {
			elemToSet = reflect.New(sliceField.Type().Elem().Elem())
			sliceField.Set(reflect.Append(sliceField, elemToSet))
		}
		if sliceField.Type().Elem().Kind() != reflect.Pointer {
			err = fmt.Errorf("set field awsslicestruct: field %s is not a slice of struct pointer, but a %s", matches[0], sliceField.Kind())
			return
		}
		setValueAtPath(elemToSet.Interface(), matches[2], v)

		return nil
	case awstagslice:
		var (
			elbTags    []elbtypes.Tag
			cfTags     []cloudformationtypes.Tag
			appendFunc func(s1, s2 string)
			assignFunc func()
		)
		switch i.(type) {
		case *elb.CreateLoadBalancerInput:
			appendFunc = func(s1, s2 string) {
				elbTags = append(elbTags, elbtypes.Tag{Key: aws.String(s1), Value: aws.String(s2)})
			}
			assignFunc = func() { v = elbTags }
		case *cloudformation.CreateStackInput, *cloudformation.UpdateStackInput:
			appendFunc = func(s1, s2 string) {
				cfTags = append(cfTags, cloudformationtypes.Tag{Key: aws.String(s1), Value: aws.String(s2)})
			}
			assignFunc = func() { v = cfTags }
		}
		for _, s := range castStringSlice(v) {
			splits := strings.SplitN(s, ":", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid tag '%s', expected 'key:value'", s)
			}
			appendFunc(splits[0], splits[1])
		}
		assignFunc()
	case awsalarmrollbacktriggers:
		var triggers []cloudformationtypes.RollbackTrigger
		if list := castStringSlice(v); len(list) > 0 {
			for _, t := range list {
				triggers = append(triggers, cloudformationtypes.RollbackTrigger{
					Arn:  aws.String(t),
					Type: aws.String("AWS::CloudWatch::Alarm"),
				})
			}
		}
		v = triggers
	default:
		// An awsType that is not handled above silently skipped conversion, so the
		// raw value reached reflect and failed there with a message about types the
		// caller never wrote — e.g. a typo'd `awsType:"awsint32"` surfaced as
		// "*int64 cannot be converted to *int32". Name the real problem instead.
		return fmt.Errorf("unknown awsType %q for field %s; add a case to setFieldWithType or use an existing type", destType, fieldPath)
	}

	setValueAtPath(i, fieldPath, v)
	return nil
}

func castString(v any) string {
	switch vv := v.(type) {
	case []string:
		return strings.Join(vv, ",")
	case *string:
		return *vv
	default:
		return fmt.Sprint(v)
	}
}

func castFloat(v any) (float64, error) {
	switch vv := v.(type) {
	case string:
		f, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			return f, fmt.Errorf("invalid float value '%s'", vv)
		}
		return f, nil
	case float32:
		return float64(vv), nil
	case float64:
		return vv, nil
	case *float64:
		return aws.ToFloat64(vv), nil
	case int:
		return float64(vv), nil
	case int64:
		return float64(vv), nil
	default:
		return 0, fmt.Errorf("cannot cast %T to float64", v)
	}
}

func castInt(v any) (int, error) {
	switch vv := v.(type) {
	case *string:
		i, err := strconv.Atoi(aws.ToString(vv))
		if err != nil {
			return i, fmt.Errorf("invalid integer value '%s'", aws.ToString(vv))
		}
		return i, nil
	case string:
		i, err := strconv.Atoi(vv)
		if err != nil {
			return i, fmt.Errorf("invalid integer value '%s'", vv)
		}
		return i, nil
	case *int:
		return aws.ToInt(vv), nil
	case int:
		return vv, nil
	case int64:
		return int(vv), nil
	case *int64:
		return int(aws.ToInt64(vv)), nil
	default:
		return 0, fmt.Errorf("cannot cast %T to int", v)
	}
}

func castBool(v any) (bool, error) {
	switch vv := v.(type) {
	case string:
		b, err := strconv.ParseBool(vv)
		if err != nil {
			return b, fmt.Errorf("invalid integer value '%s'", vv)
		}
		return b, nil
	case bool:
		return vv, nil
	case *bool:
		return aws.ToBool(vv), nil
	default:
		return false, fmt.Errorf("cannot cast %T to bool", v)
	}
}

func castInt64(v any) (int64, error) {
	switch vv := v.(type) {
	case string:
		i, err := strconv.Atoi(vv)
		if err != nil {
			return int64(i), fmt.Errorf("invalid integer value '%s'", vv)
		}
		return int64(i), nil
	case int:
		return int64(vv), nil
	case *int:
		return int64(aws.ToInt(vv)), nil
	case int64:
		return vv, nil
	case *int64:
		return aws.ToInt64(vv), nil
	default:
		return int64(0), fmt.Errorf("cannot cast %T to int64", v)
	}
}

func castStringSlice(v any) []string {
	switch vv := v.(type) {
	case string:
		return []string{vv}
	case *string:
		return []string{aws.ToString(vv)}
	case []*string:
		return aws.ToStringSlice(vv)
	case []string:
		return vv
	case []any:
		var slice []string
		for _, i := range vv {
			switch ii := i.(type) {
			case string:
				slice = append(slice, ii)
			case *string:
				slice = append(slice, *ii)
			default:
				slice = append(slice, fmt.Sprint(ii))
			}
		}
		return slice
	default:
		return []string{fmt.Sprint(v)}
	}
}

func castStringPointerSlice(v any) []*string {
	switch vv := v.(type) {
	case string:
		return []*string{&vv}
	case *string:
		return []*string{vv}
	case []*string:
		return vv
	case []string:
		return aws.StringSlice(vv)
	case []any:
		var slice []*string
		for _, i := range vv {
			switch ii := i.(type) {
			case string:
				slice = append(slice, &ii)
			case *string:
				slice = append(slice, ii)
			default:
				str := fmt.Sprint(ii)
				slice = append(slice, &str)
			}
		}
		return slice
	default:
		str := fmt.Sprint(v)
		return []*string{&str}
	}
}

func userDataContentAsBase64(v any, tplData any) (string, error) {
	userdata := castString(v)

	var readErr error
	var content []byte

	if strings.HasPrefix(strings.TrimSpace(userdata), "#") { // userdata are bash content or yml cloud script content (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts)
		r := strings.NewReplacer("\\a", "\a", "\\b", "\b", "\\f", "\f", "\\n", "\n", "\\t", "\t", "\\r", "\r", "\\v", "\v")
		content = []byte(r.Replace(userdata))
	} else if strings.HasPrefix(userdata, "http") {
		client := &http.Client{Timeout: 5 * time.Second}

		logger.ExtraVerbosef("fetching remote userdata at '%s'", userdata)
		resp, err := client.Get(userdata)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
			return "", fmt.Errorf("'%s' when fetching userdata at '%s'", resp.Status, userdata)
		}

		content, readErr = io.ReadAll(resp.Body)
	} else {
		content, readErr = os.ReadFile(userdata)
	}

	if readErr != nil {
		return "", fmt.Errorf("got userdata from '%s' but cannot read content: %w", userdata, readErr)
	}

	if tpl, err := gotemplate.New("userdata").Parse(string(content)); err != nil {
		logger.Warningf("cannot parse userdata as Go template: %s", err)
	} else {
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, tplData); err == nil {
			content = buf.Bytes()
		}
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func structSetter(s any, params map[string]any) error {
	if params == nil {
		return nil
	}
	val := reflect.ValueOf(s).Elem()
	stru := val.Type()

	for i := 0; i < stru.NumField(); i++ {
		field := stru.Field(i)
		tplName := field.Tag.Get("templateName")
		var fieldType string
		if v, ok := params[tplName]; ok {
			kind := field.Type.Kind()
			if kind == reflect.Pointer {
				switch field.Type.Elem().Kind() {
				case reflect.String:
					fieldType = awsstr
				case reflect.Int64:
					fieldType = awsint64
				case reflect.Bool:
					fieldType = awsbool
				case reflect.Float64:
					fieldType = awsfloat
				default:
					return fmt.Errorf("unknown type %s for parameter %s in struct setter", tplName, field.Type.String())
				}
			} else if kind == reflect.Slice && field.Type.Elem().Kind() == reflect.Pointer {
				switch field.Type.Elem().Elem().Kind() {
				case reflect.String:
					fieldType = awsstringslice
				case reflect.Int64:
					fieldType = awsint64slice
				default:
					return fmt.Errorf("unknown type in slice %s for parameter %s", field.Type.String(), tplName)
				}
			}
			if err := setFieldWithType(v, s, field.Name, fieldType); err != nil {
				return fmt.Errorf("%s: %w", tplName, err)
			}
		}
	}
	return nil
}

func structInjector(src, dest any, ctx map[string]any) error {
	val := reflect.ValueOf(src).Elem()
	stru := val.Type()

	for i := 0; i < stru.NumField(); i++ {
		field := stru.Field(i)
		if dstNames, ok := field.Tag.Lookup("awsName"); ok {
			splits := strings.Split(dstNames, ",")
			for _, destName := range splits {
				destName = strings.TrimSpace(destName)
				if dstType, tok := field.Tag.Lookup("awsType"); tok {
					fieldValue := val.Field(i)
					if fieldValue.IsValid() && fieldValue.Interface() != nil && !fieldValue.IsNil() {
						if err := setFieldWithType(fieldValue.Interface(), dest, destName, dstType, ctx); err != nil {
							fieldName := field.Name
							if tplName, ok := field.Tag.Lookup("templateName"); ok {
								fieldName = tplName
							}
							return fmt.Errorf("%s: %w", fieldName, err)
						}
					}
				}
			}
		}
	}
	return nil
}

func contains(arr []string, e string) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

// setValueAtPath sets a value at a dot-separated field path in a struct.
// Replaces awsutil.SetValueAtPath from SDK v1.
// convertSlice adapts a slice value to a destination slice type whose elements
// differ only by indirection — []*string to []string and the reverse.
//
// SDK v1 modeled string lists as []*string and v2 models them as []string. The
// param setters still build []*string, so without this every field declared
// awsstringslice panicked in reflect.Set: 35 fields across create tag, create
// alarm, create database, create loadbalancer and others.
//
// Handled here rather than by changing the setters so that a field genuinely
// modeled as []*T keeps working.
func convertSlice(rv reflect.Value, dest reflect.Type) (reflect.Value, bool) {
	if rv.Kind() != reflect.Slice || dest.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}

	srcElem, destElem := rv.Type().Elem(), dest.Elem()
	deref := srcElem.Kind() == reflect.Pointer && srcElem.Elem().ConvertibleTo(destElem)
	ref := destElem.Kind() == reflect.Pointer && srcElem.ConvertibleTo(destElem.Elem())
	if !deref && !ref {
		return reflect.Value{}, false
	}

	out := reflect.MakeSlice(dest, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		e := rv.Index(i)
		switch {
		case deref:
			if e.IsNil() {
				// A nil element has no value to carry over; the zero value is the
				// only sane choice for a non-pointer destination.
				out = reflect.Append(out, reflect.Zero(destElem))
				continue
			}
			out = reflect.Append(out, e.Elem().Convert(destElem))
		case ref:
			ptr := reflect.New(destElem.Elem())
			ptr.Elem().Set(e.Convert(destElem.Elem()))
			out = reflect.Append(out, ptr)
		}
	}
	return out, true
}

func setValueAtPath(i any, path string, v any) {
	path = strings.TrimPrefix(path, ".")
	if path == "" {
		return
	}
	parts := strings.Split(path, ".")
	val := reflect.ValueOf(i)
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	for idx, part := range parts {
		if val.Kind() != reflect.Struct {
			return
		}
		field := val.FieldByName(part)
		if !field.IsValid() || !field.CanSet() {
			return
		}
		if idx == len(parts)-1 {
			rv := reflect.ValueOf(v)
			if field.Kind() == reflect.Pointer {
				if rv.Kind() == reflect.Pointer {
					field.Set(rv.Convert(field.Type()))
				} else {
					ptr := reflect.New(field.Type().Elem())
					ptr.Elem().Set(rv.Convert(field.Type().Elem()))
					field.Set(ptr)
				}
			} else if rv.Kind() == reflect.Pointer && !rv.IsNil() {
				field.Set(rv.Elem().Convert(field.Type()))
			} else if rv.Type().AssignableTo(field.Type()) {
				field.Set(rv)
			} else if rv.Type().ConvertibleTo(field.Type()) {
				field.Set(rv.Convert(field.Type()))
			} else if converted, ok := convertSlice(rv, field.Type()); ok {
				field.Set(converted)
			} else {
				field.Set(rv)
			}
		} else {
			if field.Kind() == reflect.Pointer {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				val = field.Elem()
			} else {
				val = field
			}
		}
	}
}

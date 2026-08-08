package awsspec

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/aws/smithy-go"

	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"

	"github.com/fatih/color"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
)

const (
	dryRunOperation = "DryRunOperation"
	notFound        = "NotFound"
)

type BeforeRunner interface {
	BeforeRun(env.Running) error
}

type AfterRunner interface {
	AfterRun(env.Running, any) error
}

type ResultExtractor interface {
	ExtractResult(any) string
}

type InputPostProcessor interface {
	PostProcessInput(any)
}

type command interface {
	ParamsSpec() params.Spec
	inject(map[string]any) error
	Run(env.Running, map[string]any) (any, error)
}

func implementsBeforeRun(i any) (BeforeRunner, bool) {
	v, ok := i.(BeforeRunner)
	return v, ok
}

func implementsAfterRun(i any) (AfterRunner, bool) {
	v, ok := i.(AfterRunner)
	return v, ok
}

func implementsResultExtractor(i any) (ResultExtractor, bool) {
	v, ok := i.(ResultExtractor)
	return v, ok
}

func implementsInputPostProcessor(i any) (InputPostProcessor, bool) {
	v, ok := i.(InputPostProcessor)
	return v, ok
}

// fakeDryRunID builds a plausible-looking AWS id for dry-run output.
//
// math/rand is deliberate: this value is only ever displayed, never used as a
// credential, token or key, so it does not need crypto/rand. Go auto-seeds the
// global source as of 1.20, so no explicit seeding is required either.
func fakeDryRunID(entity string) string {
	suffix := rand.Intn(1e6)
	switch entity {
	case cloud.Instance:
		return fmt.Sprintf("i-%d", suffix)
	case cloud.Subnet:
		return fmt.Sprintf("subnet-%d", suffix)
	case cloud.Vpc:
		return fmt.Sprintf("vpc-%d", suffix)
	case cloud.Volume:
		return fmt.Sprintf("vol-%d", suffix)
	case cloud.SecurityGroup:
		return fmt.Sprintf("sg-%d", suffix)
	case cloud.InternetGateway:
		return fmt.Sprintf("igw-%d", suffix)
	case cloud.NatGateway:
		return fmt.Sprintf("nat-%d", suffix)
	case cloud.RouteTable:
		return fmt.Sprintf("rtb-%d", suffix)
	default:
		return fmt.Sprintf("dryrunid-%d", suffix)
	}
}

type awsCall struct {
	fnName  string
	fn      any
	logger  *logger.Logger
	setters []setter
}

// execute calls the AWS SDK function in dc.fn with input, after applying the
// setters.
//
// The call is reflective because fn's concrete type differs per operation. SDK v2
// signatures are (context.Context, *XxxInput, ...func(*Options)), so the context
// must be passed as the first argument — omitting it made reflect.Call panic with
// "Call with too few input arguments", which the deferred recover turned into an
// opaque error rather than a crash.
func (dc *awsCall) execute(ctx context.Context, input any) (output any, err error) {
	defer func() {
		if e := recover(); e != nil {
			output = nil
			err = fmt.Errorf("%s", e)
		}
	}()

	for _, s := range dc.setters {
		if err = s.set(ctx, input); err != nil {
			return nil, err
		}
	}

	fnVal := reflect.ValueOf(dc.fn)
	values := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(input)}

	start := time.Now()
	results := fnVal.Call(values)

	if err, ok := results[1].Interface().(error); ok && err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	dc.logger.ExtraVerbosef("%s call took %s", dc.fnName, time.Since(start))

	output = results[0].Interface()

	return
}

type checker struct {
	description string
	timeout     time.Duration
	frequency   time.Duration
	fetchFunc   func() (string, error)
	expect      string
	logger      *logger.Logger
	checkName   string
}

func (c *checker) check() error {
	now := time.Now().UTC()
	timer := time.NewTimer(c.timeout)
	if c.checkName == "" {
		c.checkName = "status"
	}
	defer timer.Stop()
	defer c.logger.Println()
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout of %s expired", c.timeout)
		default:
		}
		got, err := c.fetchFunc()
		if err != nil {
			return fmt.Errorf("check %s: %w", c.description, err)
		}
		if strings.EqualFold(got, c.expect) {
			c.logger.InteractiveInfof("check %s %s '%s' done", c.description, c.checkName, c.expect)
			return nil
		}
		elapsed := time.Since(now)
		c.logger.InteractiveInfof("%s %s '%s', expect '%s', timeout in %s (retry in %s)", c.description, c.checkName, got, c.expect, color.New(color.FgGreen).Sprint(c.timeout-elapsed.Round(time.Second)), c.frequency)
		time.Sleep(c.frequency)
	}
}

type enumValidator struct {
	expected []string
}

func NewEnumValidator(expected ...string) *enumValidator {
	return &enumValidator{expected: expected}
}

func (v *enumValidator) Validate(in *string) error {
	val := strings.ToLower(StringValue(in))
	for _, e := range v.expected {
		if val == strings.ToLower(e) {
			return nil
		}
	}
	var expString string
	switch len(v.expected) {
	case 0:
		return errors.New("empty enumeration")
	case 1:
		expString = fmt.Sprintf("'%s'", v.expected[0])
	case 2:
		expString = fmt.Sprintf("'%s' or '%s'", v.expected[0], v.expected[1])
	default:
		expString = fmt.Sprintf("'%s' or '%s'", strings.Join(v.expected[0:len(v.expected)-1], "', '"), v.expected[len(v.expected)-1])
	}
	return fmt.Errorf("invalid value '%s' expect %s", StringValue(in), expString)
}

func String(v string) *string {
	return &v
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func Int64(v int64) *int64 {
	return &v
}

func Int64AsIntValue(v *int64) int {
	if v != nil {
		return int(*v)
	}
	return 0
}

func Bool(v bool) *bool {
	return &v
}

func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

func decorateAWSError(err error) error {
	var aerr smithy.APIError
	if errors.As(err, &aerr) {
		return fmt.Errorf("%s: %s", aerr.ErrorCode(), aerr.ErrorMessage())
	}
	return err
}

// capitalize upper-cases the first character of s.
//
// Replaces strings.Title, deprecated in Go 1.18 because it applies Unicode word
// boundaries and title-cases every word. Every input here is a single ASCII
// token — an AWS API name, resource type, template action, or policy effect —
// so this is both correct and narrower than the deprecated behavior.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

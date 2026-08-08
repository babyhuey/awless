package awsat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	awsspec "github.com/bootswithdefer/awless/aws/spec"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template"
)

type ATBuilder struct {
	template     string
	cmdResult    *string
	expectCalls  map[string]int
	expectInput  map[string]any
	ignoredInput map[string]struct{}
	fillers      map[string]string
	expectRevert string
	mock         *Mock
	graph        *graph.Graph
}

func Template(template string) *ATBuilder {
	return &ATBuilder{template: template,
		expectCalls:  make(map[string]int),
		expectInput:  make(map[string]any),
		ignoredInput: make(map[string]struct{}),
	}
}

func (b *ATBuilder) ExpectCommandResult(key string) *ATBuilder {
	b.cmdResult = &key
	return b
}

func (b *ATBuilder) ExpectCalls(expects ...string) *ATBuilder {
	for _, expect := range expects {
		b.expectCalls[expect]++
	}
	return b
}

func (b *ATBuilder) ExpectInput(call string, input any) *ATBuilder {
	b.expectInput[call] = input
	return b
}

func (b *ATBuilder) IgnoreInput(calls ...string) *ATBuilder {
	for _, call := range calls {
		b.ignoredInput[call] = struct{}{}
	}
	return b
}

func (b *ATBuilder) Graph(g *graph.Graph) *ATBuilder {
	b.graph = g
	return b
}

func (b *ATBuilder) Mock(i *Mock) *ATBuilder {
	b.mock = i
	return b
}

func (b *ATBuilder) Fillers(fillers map[string]string) *ATBuilder {
	b.fillers = fillers
	return b
}

func (b *ATBuilder) ExpectRevert(revert string) *ATBuilder {
	b.expectRevert = revert
	return b
}

// Run executes the template and fails the test on any error.
func (b *ATBuilder) Run(t *testing.T, l ...*logger.Logger) {
	t.Helper()
	if err := b.run(t, l...); err != nil {
		t.Fatal(err)
	}
}

// RunExpectingError executes the template and returns the error rather than
// failing, so pre-flight validation can be asserted — a command that refuses a
// bad request before calling AWS is behavior worth testing.
func (b *ATBuilder) RunExpectingError(t *testing.T, l ...*logger.Logger) error {
	t.Helper()
	return b.run(t, l...)
}

func (b *ATBuilder) run(t *testing.T, l ...*logger.Logger) error {
	t.Helper()
	if b.mock == nil {
		b.mock = NewMock()
	}
	b.mock.SetInputs(b.expectInput)
	b.mock.SetIgnored(b.ignoredInput)
	b.mock.SetTesting(t)

	tpl, err := template.Parse(b.template)
	if err != nil {
		return fmt.Errorf("template parsing: %w", err)
	}
	if b.graph == nil {
		b.graph = graph.NewGraph()
	}
	awsspec.CommandFactory = NewAcceptanceFactory(b.mock, b.graph, l...)

	cenv := template.NewEnv().WithLookupCommandFunc(func(tokens ...string) any {
		return awsspec.CommandFactory.Build(strings.Join(tokens, ""))()
	}).WithMissingHolesFunc(func(key string, paramPaths []string, isOptional bool) string {
		return b.fillers[key]
	}).Build()
	compiled, cenv, err := template.Compile(tpl, cenv, template.NewRunnerCompileMode)
	if err != nil {
		return fmt.Errorf("compiling: %w", err)
	}

	ran, err := compiled.Run(template.NewRunEnv(cenv))
	if err != nil {
		return fmt.Errorf("running: %w", err)
	}
	if ran.HasErrors() {
		for _, cmd := range ran.CommandNodesIterator() {
			if cmd.Err() != nil {
				return cmd.Err()
			}
		}
	}
	if len(b.expectCalls) > 0 {
		if got, want := b.mock.Calls(), b.expectCalls; !reflect.DeepEqual(got, want) {
			return fmt.Errorf("calls: got %#v, want %#v", got, want)
		}
	}
	if b.cmdResult != nil {
		if got, want := fmt.Sprint(ran.CommandNodesIterator()[0].Result()), StringValue(b.cmdResult); got != want {
			return fmt.Errorf("command result: got %s, want %s", got, want)
		}
	}
	if b.expectRevert != "" {
		revert, err := ran.Revert()
		if err != nil {
			return fmt.Errorf("revert: %w", err)
		}
		if got, want := revert.String(), b.expectRevert; got != want {
			return fmt.Errorf("revert: got\n%s\nwant\n%s", got, want)
		}
	}

	return nil
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func String(v string) *string {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Float64(v float64) *float64 {
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

package awsat

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
)

// Mock intercepts AWS calls without reaching the network.
//
// AWS SDK v2 exposes concrete *service.Client structs rather than interfaces, so
// the pre-v2 approach of substituting an iface implementation is not available —
// which is why the previous generated mocks were gutted and every factory
// function stubbed out.
//
// Instead this registers a smithy middleware on the Finalize step that
// short-circuits before the request is signed or sent, records the operation, and
// returns a canned result. Commands build their clients with
// service.NewFromConfig(cfg), so injecting APIOptions into that config is enough
// to reach every one of them.
type Mock struct {
	basicMock

	mu        sync.Mutex
	outputs   map[string]any
	errs      map[string]error
	inputs    map[string]any
	dryRunAll bool
}

// NewMock builds an empty mock. Register expectations with On and OnError.
func NewMock() *Mock {
	return &Mock{
		outputs: make(map[string]any),
		errs:    make(map[string]error),
		inputs:  make(map[string]any),
	}
}

// On makes the named operation return output. The operation name is the SDK's,
// e.g. "PutMetricAlarm", and output must be the matching *service.XxxOutput.
func (m *Mock) On(operation string, output any) *Mock {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.outputs[operation] = output
	return m
}

// OnError makes the named operation fail.
func (m *Mock) OnError(operation string, err error) *Mock {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errs[operation] = err
	return m
}

// OnAPIError makes the named operation fail with an AWS API error carrying code.
//
// Needed for the dry-run path: a command run with DryRun set expects the API to
// reject the call with DryRunOperation, and treats that as success.
func (m *Mock) OnAPIError(operation, code, message string) *Mock {
	return m.OnError(operation, &smithy.GenericAPIError{Code: code, Message: message})
}

// OnDryRun makes every operation fail with DryRunOperation, which is what AWS does
// for a request with DryRun set and what the generated dryRun path expects.
func (m *Mock) OnDryRun() *Mock {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dryRunAll = true
	return m
}

// InputFor returns the input the command passed to the named operation, so a
// test can assert on it after the run.
func (m *Mock) InputFor(operation string) any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.inputs[operation]
}

// Config returns an aws.Config that routes every client built from it through
// this mock. A non-empty Region is required because the generated constructors
// only build a client when one is set.
func (m *Mock) Config() aws.Config {
	return aws.Config{
		Region: "us-west-2",
		APIOptions: []func(*middleware.Stack) error{
			m.middleware,
		},
	}
}

func (m *Mock) middleware(stack *middleware.Stack) error {
	// Initialize is the outermost step, so short-circuiting here skips
	// serialization, signing and transport entirely, and its Result becomes the
	// operation's output. It is also the only step that still carries the typed
	// input parameters.
	return stack.Initialize.Add(
		middleware.InitializeMiddlewareFunc("awlessAcceptanceMock", m.handle),
		middleware.Before,
	)
}

func (m *Mock) handle(ctx context.Context, in middleware.InitializeInput, _ middleware.InitializeHandler) (
	middleware.InitializeOutput, middleware.Metadata, error,
) {
	op := middleware.GetOperationName(ctx)

	m.mu.Lock()
	m.inputs[op] = in.Parameters
	err, hasErr := m.errs[op]
	out, hasOut := m.outputs[op]
	dryRunAll := m.dryRunAll
	m.mu.Unlock()

	m.addCall(op)
	m.verifyInput(op, in.Parameters)

	if hasErr {
		return middleware.InitializeOutput{}, middleware.Metadata{}, err
	}
	if dryRunAll {
		return middleware.InitializeOutput{}, middleware.Metadata{}, &smithy.GenericAPIError{
			Code:    "DryRunOperation",
			Message: "Request would have succeeded, but DryRun flag is set",
		}
	}
	if !hasOut {
		// Returning the zero output would silently produce nil-dereference panics
		// deep inside a command's result extraction, so fail loudly instead.
		return middleware.InitializeOutput{}, middleware.Metadata{}, fmt.Errorf(
			"acceptance mock: no output registered for operation %q; add mock.On(%q, &service.%sOutput{...})", op, op, op)
	}
	return middleware.InitializeOutput{Result: out}, middleware.Metadata{}, nil
}

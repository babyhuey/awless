package awsfetch

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	codedeploytypes "github.com/aws/aws-sdk-go-v2/service/codedeploy/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/aws/smithy-go/middleware"

	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/logger"
)

// pagedMock returns a queued sequence of outputs per operation, which is what makes it
// possible to assert that a fetcher follows a continuation token rather than stopping at
// the first page.
//
// Several services publish no paginators, so those fetchers carry a hand-written
// NextToken/NextMarker loop. A loop that never advances silently returns only the first
// page — the user sees a short list and no error, which is the worst shape of bug for a
// tool whose whole job is showing you what you have.
type pagedMock struct {
	mu    sync.Mutex
	pages map[string][]any
	calls map[string]int
}

func newPagedMock() *pagedMock {
	return &pagedMock{pages: map[string][]any{}, calls: map[string]int{}}
}

// on queues outputs for an operation, returned one per call in order.
func (m *pagedMock) on(op string, outputs ...any) *pagedMock {
	m.pages[op] = append(m.pages[op], outputs...)
	return m
}

func (m *pagedMock) callCount(op string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls[op]
}

func (m *pagedMock) config() aws.Config {
	return aws.Config{
		Region:     "us-west-2",
		APIOptions: []func(*middleware.Stack) error{m.middleware},
	}
}

func (m *pagedMock) middleware(stack *middleware.Stack) error {
	return stack.Initialize.Add(
		middleware.InitializeMiddlewareFunc("awlessPagedMock", m.handle),
		middleware.Before,
	)
}

func (m *pagedMock) handle(ctx context.Context, _ middleware.InitializeInput, _ middleware.InitializeHandler) (
	middleware.InitializeOutput, middleware.Metadata, error,
) {
	op := middleware.GetOperationName(ctx)

	m.mu.Lock()
	defer m.mu.Unlock()

	n := m.calls[op]
	m.calls[op]++

	queued, ok := m.pages[op]
	if !ok {
		return middleware.InitializeOutput{}, middleware.Metadata{}, fmt.Errorf("paged mock: no output registered for %q", op)
	}
	if n >= len(queued) {
		// A loop asking for more pages than were queued has not honored the end of
		// pagination, which would spin against a real API.
		return middleware.InitializeOutput{}, middleware.Metadata{}, fmt.Errorf(
			"paged mock: %q called %d times but only %d pages queued; the fetcher is not stopping", op, n+1, len(queued))
	}
	return middleware.InitializeOutput{Result: queued[n]}, middleware.Metadata{}, nil
}

func str(s string) *string { return &s }

// runFetcher builds the fetch funcs the way a service does and runs one of them.
func runFetcher(t *testing.T, conf *Config, addManual func(*Config, map[string]fetch.Func), key string) []*graph.Resource {
	t.Helper()
	funcs := map[string]fetch.Func{}
	addManual(conf, funcs)

	fn, ok := funcs[key]
	if !ok {
		t.Fatalf("no fetch func registered for %q", key)
	}
	res, _, err := fn(context.Background(), fetch.NewFetcher(nil))
	if err != nil {
		t.Fatalf("fetching %q: %s", key, err)
	}
	return res
}

// EventBridge publishes no paginators, so this loop is hand-written and threads NextToken.
func TestEventBusFetcherFollowsNextToken(t *testing.T) {
	mock := newPagedMock().on("ListEventBuses",
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("default"), Arn: str("arn:1")}},
			NextToken:  str("page2"),
		},
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("orders"), Arn: str("arn:2")}},
			// No NextToken: the loop must stop here.
		},
	)

	conf := NewConfig(eventbridge.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualEventbridgeFetchFuncs, "eventbus")

	if got := mock.callCount("ListEventBuses"); got != 2 {
		t.Errorf("ListEventBuses called %d times, want 2 — the second page was not requested", got)
	}
	if len(res) != 2 {
		t.Fatalf("got %d event buses, want 2 across both pages", len(res))
	}
	ids := []string{res[0].ID(), res[1].ID()}
	if ids[0] != "default" || ids[1] != "orders" {
		t.Errorf("got %v, want [default orders]", ids)
	}
}

// An empty-string token is the other way pagination ends, and treating it as a real token
// would loop forever.
func TestEventBusFetcherStopsOnAnEmptyToken(t *testing.T) {
	mock := newPagedMock().on("ListEventBuses",
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("only"), Arn: str("arn:1")}},
			NextToken:  str(""),
		},
	)

	conf := NewConfig(eventbridge.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualEventbridgeFetchFuncs, "eventbus")

	if got := mock.callCount("ListEventBuses"); got != 1 {
		t.Errorf("ListEventBuses called %d times, want 1 — an empty token is not a next page", got)
	}
	if len(res) != 1 {
		t.Errorf("got %d event buses, want 1", len(res))
	}
}

func TestEventRuleFetcherFollowsNextToken(t *testing.T) {
	mock := newPagedMock().on("ListRules",
		&eventbridge.ListRulesOutput{
			Rules:     []eventbridgetypes.Rule{{Name: str("nightly"), Arn: str("arn:1")}},
			NextToken: str("page2"),
		},
		&eventbridge.ListRulesOutput{
			Rules: []eventbridgetypes.Rule{{Name: str("hourly"), Arn: str("arn:2")}},
		},
	)

	conf := NewConfig(eventbridge.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualEventbridgeFetchFuncs, "eventrule")

	if got := mock.callCount("ListRules"); got != 2 {
		t.Errorf("ListRules called %d times, want 2", got)
	}
	if len(res) != 2 {
		t.Errorf("got %d rules, want 2 across both pages", len(res))
	}
}

// WAF v2 uses NextMarker rather than NextToken. Getting the field name wrong is a silent
// single-page fetch, since the request field is simply never set.
//
// The call count assumes one WAF scope, which holds because CLOUDFRONT is only listable
// from us-east-1 and the mock config is us-west-2. A test in us-east-1 would see the
// fetcher paginate twice, once per scope.
func TestWebACLFetcherFollowsNextMarker(t *testing.T) {
	mock := newPagedMock().on("ListWebACLs",
		&wafv2.ListWebACLsOutput{
			WebACLs:    []wafv2types.WebACLSummary{{Id: str("acl-1"), Name: str("edge"), ARN: str("arn:1")}},
			NextMarker: str("marker2"),
		},
		&wafv2.ListWebACLsOutput{
			WebACLs: []wafv2types.WebACLSummary{{Id: str("acl-2"), Name: str("regional"), ARN: str("arn:2")}},
		},
	)

	conf := NewConfig(wafv2.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualWafFetchFuncs, "webacl")

	if got := mock.callCount("ListWebACLs"); got != 2 {
		t.Errorf("ListWebACLs called %d times, want 2 — NextMarker was not threaded", got)
	}
	if len(res) != 2 {
		t.Fatalf("got %d web ACLs, want 2 across both pages", len(res))
	}
}

func TestIPSetFetcherFollowsNextMarker(t *testing.T) {
	mock := newPagedMock().on("ListIPSets",
		&wafv2.ListIPSetsOutput{
			IPSets:     []wafv2types.IPSetSummary{{Id: str("ips-1"), Name: str("blocklist"), ARN: str("arn:1")}},
			NextMarker: str("marker2"),
		},
		&wafv2.ListIPSetsOutput{
			IPSets: []wafv2types.IPSetSummary{{Id: str("ips-2"), Name: str("allowlist"), ARN: str("arn:2")}},
		},
	)

	conf := NewConfig(wafv2.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualWafFetchFuncs, "ipset")

	if got := mock.callCount("ListIPSets"); got != 2 {
		t.Errorf("ListIPSets called %d times, want 2", got)
	}
	if len(res) != 2 {
		t.Errorf("got %d IP sets, want 2 across both pages", len(res))
	}
}

// Three pages, to catch a loop that follows exactly one continuation and then stops.
func TestEventBusFetcherFollowsMoreThanOneToken(t *testing.T) {
	mock := newPagedMock().on("ListEventBuses",
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("a"), Arn: str("arn:a")}},
			NextToken:  str("p2"),
		},
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("b"), Arn: str("arn:b")}},
			NextToken:  str("p3"),
		},
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("c"), Arn: str("arn:c")}},
		},
	)

	conf := NewConfig(eventbridge.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualEventbridgeFetchFuncs, "eventbus")

	if got := mock.callCount("ListEventBuses"); got != 3 {
		t.Errorf("ListEventBuses called %d times, want 3", got)
	}
	if len(res) != 3 {
		t.Errorf("got %d event buses, want 3 across all pages", len(res))
	}
}

// An error partway through pagination must surface rather than being reported as a short
// but successful list.
func TestEventBusFetcherReturnsAPaginationError(t *testing.T) {
	mock := newPagedMock().on("ListEventBuses",
		&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{{Name: str("first"), Arn: str("arn:1")}},
			NextToken:  str("page2"),
		},
		// Nothing queued for the second call, so the mock errors — standing in for an
		// API failure on a later page.
	)

	conf := NewConfig(eventbridge.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	funcs := map[string]fetch.Func{}
	addManualEventbridgeFetchFuncs(conf, funcs)

	_, _, err := funcs["eventbus"](context.Background(), fetch.NewFetcher(nil))
	if err == nil {
		t.Fatal("expected the mid-pagination failure to be returned, not a truncated list")
	}
}

// The other hand-written pattern is a per-parent fan-out, for resources AWS will not list
// globally. Its characteristic bug is collecting only one parent's children — the last one
// to finish, or the first — so this uses two applications with distinct groups and asserts
// both survive.
func TestDeploymentGroupFetcherCollectsEveryApplicationsGroups(t *testing.T) {
	mock := newPagedMock().
		on("ListApplications", &codedeploy.ListApplicationsOutput{
			Applications: []string{"web", "api"},
		}).
		on("ListDeploymentGroups",
			&codedeploy.ListDeploymentGroupsOutput{DeploymentGroups: []string{"web-prod"}},
			&codedeploy.ListDeploymentGroupsOutput{DeploymentGroups: []string{"api-prod", "api-staging"}},
		).
		on("BatchGetDeploymentGroups",
			&codedeploy.BatchGetDeploymentGroupsOutput{DeploymentGroupsInfo: []codedeploytypes.DeploymentGroupInfo{
				{ApplicationName: str("web"), DeploymentGroupName: str("web-prod"), DeploymentGroupId: str("dg-1")},
			}},
			&codedeploy.BatchGetDeploymentGroupsOutput{DeploymentGroupsInfo: []codedeploytypes.DeploymentGroupInfo{
				{ApplicationName: str("api"), DeploymentGroupName: str("api-prod"), DeploymentGroupId: str("dg-2")},
				{ApplicationName: str("api"), DeploymentGroupName: str("api-staging"), DeploymentGroupId: str("dg-3")},
			}},
		)

	conf := NewConfig(codedeploy.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualCodedeployFetchFuncs, "deploymentgroup")

	if got := mock.callCount("ListDeploymentGroups"); got != 2 {
		t.Errorf("ListDeploymentGroups called %d times, want one per application", got)
	}
	if len(res) != 3 {
		t.Fatalf("got %d deployment groups, want 3 across both applications", len(res))
	}
	seen := map[string]bool{}
	for _, r := range res {
		seen[r.ID()] = true
	}
	// The resource is identified by the group name, which is unique per application.
	for _, want := range []string{"web-prod", "api-prod", "api-staging"} {
		if !seen[want] {
			t.Errorf("deployment group %s is missing; only %v were collected", want, seen)
		}
	}
}

// An application with no deployment groups must be skipped rather than sent to
// BatchGetDeploymentGroups, which rejects an empty name list.
func TestDeploymentGroupFetcherSkipsAnApplicationWithNoGroups(t *testing.T) {
	mock := newPagedMock().
		on("ListApplications", &codedeploy.ListApplicationsOutput{Applications: []string{"empty"}}).
		on("ListDeploymentGroups", &codedeploy.ListDeploymentGroupsOutput{DeploymentGroups: nil})
	// No BatchGetDeploymentGroups queued: reaching it would error out of the mock.

	conf := NewConfig(codedeploy.NewFromConfig(mock.config()))
	conf.Log = logger.DiscardLogger

	res := runFetcher(t, conf, addManualCodedeployFetchFuncs, "deploymentgroup")

	if len(res) != 0 {
		t.Errorf("got %d groups, want none", len(res))
	}
	if got := mock.callCount("BatchGetDeploymentGroups"); got != 0 {
		t.Errorf("BatchGetDeploymentGroups called %d times for an application with no groups", got)
	}
}

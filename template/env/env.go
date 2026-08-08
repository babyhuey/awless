package env

import (
	"context"

	"github.com/bootswithdefer/awless/logger"
)

const (
	Fillers = iota
	ProcessedFillers
	ResolvedVars
)

const (
	RequiredAndSuggestedParams = iota
	RequiredParamsOnly
	AllParams
)

type log interface {
	Log() *logger.Logger
}

type Running interface {
	log
	Context() map[string]any
	// RequestContext is the Go context for outbound AWS calls made while running
	// a template. Named to avoid colliding with Context(), which is the template
	// variable map. Never nil: implementations fall back to context.Background().
	//
	// Threading it here rather than through every command signature means the
	// 260-odd call sites that used context.Background() directly — 163 of them
	// generated — get cancellation without an interface break, since commands
	// already receive an env.Running.
	RequestContext() context.Context
	SetRequestContext(context.Context)
	IsDryRun() bool
	SetDryRun(b bool)
}

type Compiling interface {
	log
	LookupCommandFunc() func(...string) any
	AliasFunc() func(paramPath, alias string) string
	MissingHolesFunc() func(string, []string, bool) string
	ParamsMode() int
	Push(int, ...map[string]any)
	Get(int) map[string]any
}

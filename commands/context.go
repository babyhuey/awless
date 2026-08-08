package commands

import (
	"context"
	"sync"
)

// rootContext is the process-wide context, canceled when the user interrupts.
//
// It is package state rather than a parameter because cobra command bodies here
// are closures assigned to Run fields, and the AWS call sites reach it through
// env.Running.RequestContext() rather than through their own signatures. main
// sets it once before Execute.
var (
	rootCtxMu sync.RWMutex
	rootCtx   context.Context
)

// SetRootContext installs the process context. Called once from main.
func SetRootContext(ctx context.Context) {
	rootCtxMu.Lock()
	defer rootCtxMu.Unlock()
	//nolint:fatcontext // Storing the process context once, not deriving a new one
	// per iteration; fatcontext exists to catch context nesting inside loops.
	rootCtx = ctx
}

// RootContext returns the process context, or context.Background() if it was
// never set — which happens in tests that exercise commands directly.
func RootContext() context.Context {
	rootCtxMu.RLock()
	defer rootCtxMu.RUnlock()
	if rootCtx == nil {
		return context.Background()
	}
	return rootCtx
}

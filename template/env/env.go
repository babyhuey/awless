package env

import (
	"github.com/bootswithdefer/awless/logger"
)

const (
	FILLERS = iota
	PROCESSED_FILLERS
	RESOLVED_VARS
)

const (
	REQUIRED_AND_SUGGESTED_PARAMS = iota
	REQUIRED_PARAMS_ONLY
	ALL_PARAMS
)

type log interface {
	Log() *logger.Logger
}

type Running interface {
	log
	Context() map[string]any
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

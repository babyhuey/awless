package template

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
)

// Context returns the runner's context, never nil.
func (ru *Runner) Context() context.Context {
	if ru.Ctx == nil {
		return context.Background()
	}
	return ru.Ctx
}

type Runner struct {
	Template                               *Template
	Locale, Profile, Message, TemplatePath string
	Log                                    *logger.Logger
	Fillers                                []map[string]any
	AliasFunc                              func(paramPath, alias string) string
	MissingHolesFunc                       func(string, []string, bool) string
	CmdLookuper                            func(tokens ...string) any
	Validators                             []Validator
	ParamsSuggested                        int

	// Ctx is the caller's context, propagated to every AWS call the template
	// makes. Leave nil in tests; Context() falls back to context.Background().
	Ctx context.Context

	BeforeRun func(*TemplateExecution) (bool, error)
	AfterRun  func(*TemplateExecution) error
}

func (ru *Runner) Run() error {
	tplExec := &TemplateExecution{
		Template: ru.Template,
		Path:     ru.TemplatePath,
		Locale:   ru.Locale,
		Profile:  ru.Profile,
		Source:   ru.Template.String(),
	}
	tplExec.SetMessage(ru.Message)

	cenv := NewEnv().WithAliasFunc(ru.AliasFunc).WithMissingHolesFunc(ru.MissingHolesFunc).
		WithLookupCommandFunc(ru.CmdLookuper).WithLog(ru.Log).WithParamsMode(ru.ParamsSuggested).Build()
	cenv.Push(env.Fillers, ru.Fillers...)

	var err error
	tplExec.Template, cenv, err = Compile(tplExec.Template, cenv, NewRunnerCompileMode)
	if err != nil {
		return err
	}

	tplExec.Fillers = cenv.Get(env.ProcessedFillers)

	errs := tplExec.Validate(ru.Validators...)
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Warning(err)
		}
		fmt.Fprintln(os.Stderr)
	}

	if tplExec.IsOneLiner() {
		logger.Verbose("Dry running template ...")
	} else {
		logger.Info("Dry running template ...")
	}

	renv := NewRunEnv(cenv)
	// Cancellation reaches every AWS call made by the commands through here; see
	// env.Running.RequestContext.
	renv.SetRequestContext(ru.Context())
	if _, err = tplExec.DryRun(renv); err != nil {
		var t *Errors
		if errors.As(err, &t) {
			errs, _ := t.Errors()
			for _, e := range errs {
				logger.Errorf("%s", e.Error())
			}
		} else {
			logger.Error(err)
		}
		return errors.New("dry run failed")
	}

	ok, err := ru.BeforeRun(tplExec)
	if err != nil {
		return err
	}

	if ok {
		tplExec.Template, err = tplExec.Run(renv)
		if err != nil {
			logger.Errorf("Running template error: %s", err)
		}
		if err := ru.AfterRun(tplExec); err != nil {
			return err
		}
	}

	if ko := tplExec.Stats().KOCount; ko > 0 {
		// Each failure was already reported as it happened, so this only needs to
		// make the exit status reflect them. Returned rather than exited so the
		// caller's deferred work — persisting the template log — still runs.
		return fmt.Errorf("%d command(s) failed", ko)
	}

	return nil
}

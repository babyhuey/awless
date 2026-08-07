package commands

import (
	"fmt"
	"os"
	"strings"

	awsservices "github.com/bootswithdefer/awless/aws/services"
	awsspec "github.com/bootswithdefer/awless/aws/spec"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/config"
	"github.com/bootswithdefer/awless/database"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/sync"
	"github.com/bootswithdefer/awless/template"
	"github.com/bootswithdefer/awless/template/env"
)

func NewRunnerRequiredParamsOnly(tpl *template.Template, msg, tplPath string, fillers ...map[string]any) *template.Runner {
	r := NewRunner(tpl, msg, tplPath, fillers...)
	r.ParamsSuggested = env.REQUIRED_PARAMS_ONLY
	return r
}

func NewRunner(tpl *template.Template, msg, tplPath string, fillers ...map[string]any) *template.Runner {
	runner := &template.Runner{Ctx: RootContext()}

	runner.Template = tpl
	runner.Locale = config.GetAWSRegion()
	runner.Profile = config.GetAWSProfile()
	runner.Log = logger.DefaultLogger
	runner.Message = msg
	runner.TemplatePath = tplPath
	runner.Fillers = fillers
	runner.AliasFunc = resolveAliasFunc
	runner.MissingHolesFunc = missingHolesStdinFunc()
	if allSuggestedParamsFlag {
		runner.ParamsSuggested = env.ALL_PARAMS
	}
	if noSuggestedParamsFlag {
		runner.ParamsSuggested = env.REQUIRED_PARAMS_ONLY
	}

	runner.Validators = []template.Validator{
		&template.UniqueNameValidator{LookupGraph: func(key string) (cloud.GraphAPI, bool) {
			g := sync.LoadLocalGraphForService(awsservices.ServicePerResourceType[key], config.GetAWSProfile(), config.GetAWSRegion())
			return g, true
		}},
		&template.ParamIsSetValidator{Action: "create", Entity: "instance", Param: "keypair", WarningMessage: "This instance has no access keypair. You might not be able to connect to it. Use `awless create instance keypair=my-keypair ...`"},
	}

	runner.CmdLookuper = func(tokens ...string) any {
		newCommandFunc := awsspec.CommandFactory.Build(strings.Join(tokens, ""))
		if newCommandFunc == nil {
			return nil
		}
		return newCommandFunc()
	}

	runner.BeforeRun = func(tplExec *template.TemplateExecution) (bool, error) {
		var yesorno string
		if forceGlobalFlag {
			yesorno = "y"
		} else {
			fmt.Printf("%s\n\n", renderGreenFn(tplExec.Template))
			fmt.Printf("Confirm (region: %s)? [y/N] ", config.GetAWSRegion())
			if _, err := fmt.Scanln(&yesorno); err != nil && err.Error() != "unexpected newline" {
				return false, err
			}
		}

		if strings.TrimSpace(strings.ToLower(yesorno)) == "y" {
			me, err := awsservices.AccessService.(*awsservices.Access).GetIdentity()
			if err != nil {
				logger.Warningf("cannot resolve template author identity: %s", err)
			} else {
				tplExec.Author = me.ResourcePath
				logger.ExtraVerbosef("resolved template author: %s", tplExec.Author)
			}
			return true, nil
		}
		os.Exit(1)
		return false, nil
	}

	runner.AfterRun = func(tplExec *template.TemplateExecution) error {
		if tplExec.Message == "" {
			if tplExec.IsOneLiner() {
				// Redacted: Message is persisted to the local template log, and
				// a one-liner message embeds the full command line, which may
				// carry a secret (e.g. `create loginprofile password=...`).
				tplExec.SetMessage(fmt.Sprintf("Run %s", tplExec.StringRedacted()))
			} else if path := tplExec.Path; path != "" {
				stats := tplExec.Stats()
				if stats.KOCount > 0 {
					tplExec.SetMessage(fmt.Sprintf("Run %d/%d commands from %s", stats.OKCount, stats.CmdCount, path))
				} else {
					tplExec.SetMessage(fmt.Sprintf("Run %d commands from %s", stats.OKCount, path))
				}
			}
		}

		if err := database.Execute(func(db *database.DB) error {
			return db.AddTemplate(tplExec)
		}); err != nil {
			logger.Errorf("Cannot save executed template in awless logs: %s", err)
		}

		if template.IsRevertible(tplExec.Template) {
			fmt.Println()
			logger.Infof("Revert this template with `awless revert %s`", tplExec.ID)
		}

		runSyncFor(tplExec)

		return nil
	}

	return runner
}

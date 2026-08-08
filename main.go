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

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"

	"github.com/bootswithdefer/awless/commands"
)

func main() {
	// Cancel on SIGINT/SIGTERM so in-flight AWS calls stop instead of being
	// abandoned when the user hits Ctrl-C. The context reaches the AWS SDK via
	// env.Running.RequestContext(); see template/env.
	//
	// A second signal aborts immediately: signal.NotifyContext restores the
	// default disposition after the first one, so an unresponsive command can
	// still be killed with a repeated Ctrl-C.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	commands.SetRootContext(ctx)

	// The single exit point. Commands return errors rather than calling os.Exit, so
	// deferred cleanup always runs, and reporting happens in exactly one place —
	// RootCmd sets SilenceErrors so cobra does not also print.
	if err := commands.RootCmd.ExecuteContext(ctx); err != nil {
		if errors.Is(err, commands.ErrExitZero) {
			// The command already told the user what they needed; that is success.
			return
		}
		if !errors.Is(err, commands.ErrReported) {
			// ErrReported means the reason was already printed, usually with a
			// suggested command; printing again would only repeat it.
			fmt.Fprintln(os.Stderr, color.RedString("[error] "), err)
		}
		os.Exit(1)
	}
}

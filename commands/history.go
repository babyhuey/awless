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

package commands

import (
	"fmt"
	"os"

	awsservices "github.com/bootswithdefer/awless/aws/services"

	"github.com/spf13/cobra"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/config"
	"github.com/bootswithdefer/awless/console"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/sync"
)

var (
	showProperties bool
)

func init() {
	RootCmd.AddCommand(historyCmd)

	historyCmd.Flags().BoolVar(&showProperties, "properties", false, "Full diff with resources properties")
}

var historyCmd = &cobra.Command{
	Use:                "history",
	Hidden:             true,
	Short:              "(in progress) Show a infra resource history & changes using your locally sync snapshots",
	PersistentPreRunE:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRunE: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		region := config.GetAWSRegion()

		root := graph.InitResource(cloud.Region, region)

		var diffs []*sync.Diff

		all, err := sync.DefaultSyncer.List()
		if err != nil {
			return err
		}

		for i := 1; i < len(all); i++ {
			from, err := sync.DefaultSyncer.LoadRev(all[i-1].ID)
			if err != nil {
				return err
			}

			to, err := sync.DefaultSyncer.LoadRev(all[i].ID)
			if err != nil {
				return err
			}

			d, err := sync.BuildDiff(from, to, root.ID())
			if err != nil {
				return err
			}

			diffs = append(diffs, d)
		}

		for _, diff := range diffs {
			if err := displayRevisionDiff(diff, awsservices.InfraService.Name(), root, verboseGlobalFlag); err != nil {
				return err
			}
		}

		return nil
	},
}

func displayRevisionDiff(diff *sync.Diff, cloudService string, root *graph.Resource, verbose bool) error {
	fromRevision := "repository creation"
	if diff.From.ID != "" {
		fromRevision = diff.From.ID[:7] + " on " + diff.From.Date.Format("Monday January 2, 15:04")
	}

	var graphdiff *graph.Diff
	if cloudService == awsservices.InfraService.Name() {
		graphdiff = diff.InfraDiff
	}

	if showProperties {
		if graphdiff.HasDiff() {
			fmt.Println("▶", cloudService, "properties, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			displayer, err := console.BuildOptions(
				console.WithFormat("table"),
				console.WithRootNode(root),
			).SetSource(graphdiff).Build()
			if err != nil {
				return err
			}
			if err := displayer.Print(os.Stdout); err != nil {
				return err
			}
			fmt.Println()
		} else if verbose {
			fmt.Println("▶", cloudService, "properties, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			fmt.Println("No changes.")
		}
	} else {
		if graphdiff.HasDiff() {
			fmt.Println("▶", cloudService, "resources, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			displayer, err := console.BuildOptions(
				console.WithFormat("tree"),
				console.WithRootNode(root),
			).SetSource(graphdiff).Build()
			if err != nil {
				return err
			}
			if err := displayer.Print(os.Stdout); err != nil {
				return err
			}
			fmt.Println()
		} else if verbose {
			fmt.Println("▶", cloudService, "resources, from", fromRevision,
				"to", diff.To.ID[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			fmt.Println("No resource changes.")
		}
	}
	return nil
}

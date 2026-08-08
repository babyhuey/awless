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
	"time"

	"github.com/spf13/cobra"

	awstailers "github.com/bootswithdefer/awless/aws/tailers"
)

var tailFollowFrequencyFlag time.Duration
var tailEnableFollowFlag bool
var tailNumberEventsFlag int
var stackEventsFilters []string
var stackEventsTailTimeout time.Duration
var cancelStackUpdateAfterTimeout bool

func init() {
	RootCmd.AddCommand(tailCmd)

	tailCmd.PersistentFlags().DurationVar(&tailFollowFrequencyFlag, "frequency", 10*time.Second, "Fetch refresh frequency")
	tailCmd.PersistentFlags().BoolVar(&tailEnableFollowFlag, "follow", false, "Periodically refresh and append new data to output")
	tailCmd.PersistentFlags().IntVarP(&tailNumberEventsFlag, "number", "n", 10, "Number of events to display")

	tailCmd.AddCommand(scalingActivitiesCmd)

	stackEventsCmd.PersistentFlags().StringArrayVar(&stackEventsFilters, "filters",
		[]string{awstailers.StackEventTimestamp, awstailers.StackEventLogicalID, awstailers.StackEventType, awstailers.StackEventStatus},
		fmt.Sprintf("Filter the output columns. Valid filters: %s, %s, %s, %s, %s",
			awstailers.StackEventLogicalID,
			awstailers.StackEventStatus,
			awstailers.StackEventStatusReason,
			awstailers.StackEventTimestamp,
			awstailers.StackEventType))

	stackEventsCmd.PersistentFlags().BoolVar(&cancelStackUpdateAfterTimeout, "cancel-on-timeout", false, "Cancel stack update when timeout is reached, use with 'timeout' flag")
	stackEventsCmd.PersistentFlags().DurationVar(&stackEventsTailTimeout, "timeout", 1*time.Hour, "Time to wait for stack update to complete, use with 'follow' flag")

	tailCmd.AddCommand(stackEventsCmd)
}

var tailCmd = &cobra.Command{
	Use:                "tail",
	Hidden:             true,
	PersistentPreRunE:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, firstInstallDoneHook),
	PersistentPostRunE: applyHooks(verifyNewVersionHook, networkMonitorHook),
	Short:              "Tail cloud events",
}

var scalingActivitiesCmd = &cobra.Command{
	Use:   "scaling-activities",
	Short: "Watch scaling-activities",

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := awstailers.NewScalingActivitiesTailer(tailNumberEventsFlag, tailEnableFollowFlag, tailFollowFrequencyFlag).Tail(os.Stdout); err != nil {
			return err
		}
		return nil
	},
}

var stackEventsCmd = &cobra.Command{
	Use:   "stack-events",
	Short: "Watch stack-events",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("expecting stack-name string")
		}

		if err := awstailers.NewCloudformationEventsTailer(args[0], tailNumberEventsFlag, tailEnableFollowFlag, tailFollowFrequencyFlag, stackEventsFilters, stackEventsTailTimeout, cancelStackUpdateAfterTimeout).Tail(os.Stdout); err != nil {
			return err
		}
		return nil
	},
}

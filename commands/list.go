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
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	awsservices "github.com/bootswithdefer/awless/aws/services"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/config"
	"github.com/bootswithdefer/awless/console"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/sync"
)

type contextKey string

var (
	listingFormat              string
	listingFiltersFlag         []string
	listingTagFiltersFlag      []string
	listingTagKeyFiltersFlag   []string
	listingTagValueFiltersFlag []string
	listingColumnsFlag         []string
	listOnlyIDs                bool
	noHeadersFlag              bool
	sortBy                     []string
	reverseFlag                bool
	allLocalRegionsFlag        bool
)

func init() {
	RootCmd.AddCommand(listCmd)

	cobra.EnableCommandSorting = false

	for _, srvName := range awsservices.ServiceNames {
		listCmd.AddCommand(listAllResourceInServiceCmd(srvName))
	}

	for _, name := range awsservices.ServiceNames {
		resources := append([]string{}, awsservices.ResourceTypesPerServiceName()[name]...)
		sort.Strings(resources)
		for _, resType := range resources {
			listCmd.AddCommand(listSpecificResourceCmd(resType))
		}
	}

	listCmd.PersistentFlags().StringVar(&listingFormat, "format", "table", "Output format: table, csv, tsv, json (default to table)")
	listCmd.PersistentFlags().StringSliceVar(&listingFiltersFlag, "filter", []string{}, "Filter resources given key/values fields (case insensitive). Ex: --filter type=t2.micro")
	listCmd.PersistentFlags().StringArrayVar(&listingTagFiltersFlag, "tag", []string{}, "Filter EC2 resources given tags (case sensitive!). Ex: --tag Env=Production")
	listCmd.PersistentFlags().StringArrayVar(&listingTagKeyFiltersFlag, "tag-key", []string{}, "Filter EC2 resources given a tag key only (case sensitive!). Ex: --tag-key Env")
	listCmd.PersistentFlags().StringArrayVar(&listingTagValueFiltersFlag, "tag-value", []string{}, "Filter EC2 resources given a tag value only (case sensitive!). Ex: --tag-value Staging")
	listCmd.PersistentFlags().StringSliceVar(&listingColumnsFlag, "columns", []string{}, "Select the properties to display in the columns. Ex: --columns id,name,cidr")
	listCmd.PersistentFlags().BoolVar(&listOnlyIDs, "ids", false, "List only ids")
	listCmd.PersistentFlags().BoolVar(&noHeadersFlag, "no-headers", false, "Do not display headers")
	listCmd.PersistentFlags().BoolVar(&allLocalRegionsFlag, "all-local-regions", false, "List resources from every locally synced region, adding a region column. Implies --local")
	listCmd.PersistentFlags().BoolVar(&reverseFlag, "reverse", false, "Use in conjunction with --sort to reverse sort")
	listCmd.PersistentFlags().StringSliceVar(&sortBy, "sort", []string{"Id"}, "Sort tables by column(s) name(s)")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: `  awless list instances --sort uptime
  awless list users --format csv
  awless list volumes --filter state=in-use --filter type=gp2
  awless list volumes --tag-value Purchased
  awless list vpcs --tag-key Dept --tag-key Internal
  awless list instances --tag Env=Production,Dept=Marketing
  awless list instances --filter state=running,type=t2.micro
  awless list s3objects --filter bucket=pdf-bucket
  awless list loggroups --sort created --reverse
  awless list loggroups --format json
  awless list trails
  awless list eksclusters
  awless list dynamodbtables
  awless list secrets
  awless list ssmparameters
  awless list filesystems
  awless list apigateways
  awless list instances --all-local-regions`,
	PersistentPreRunE:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, firstInstallDoneHook),
	PersistentPostRunE: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),
	Short:              "List resources: sorting, filtering via tag/properties, output formatting, etc...",
}

var listSpecificResourceCmd = func(resType string) *cobra.Command {
	return &cobra.Command{
		Use:   cloud.PluralizeResource(resType),
		Short: fmt.Sprintf("[%s] List %s %s", awsservices.ServicePerResourceType[resType], strings.ToUpper(awsservices.APIPerResourceType[resType]), cloud.PluralizeResource(resType)),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				var plural string
				if len(args) > 1 {
					plural = "s"
				}
				logger.Errorf("invalid parameter%s '%s'", plural, strings.Join(args, " "))
				if strings.Contains(args[0], "=") {
					if !promptConfirmDefaultYes("Did you mean `awless list %s --filter %s`? ", cloud.PluralizeResource(resType), strings.Join(args, " ")) {
						return ErrReported
					}
					listingFiltersFlag = append(listingFiltersFlag, args...)
				} else {
					return ErrReported
				}
			}
			var g cloud.GraphAPI

			switch {
			case allLocalRegionsFlag:
				srvName, ok := awsservices.ServicePerResourceType[resType]
				if !ok {
					return fmt.Errorf("cannot find service for resource type %s", resType)
				}
				var err error
				g, err = sync.LoadLocalGraphForTypeInAllRegions(srvName, resType, config.GetAWSProfile())
				if err != nil {
					return err
				}
			case localGlobalFlag:
				if srvName, ok := awsservices.ServicePerResourceType[resType]; ok {
					g = sync.LoadLocalGraphForService(srvName, config.GetAWSProfile(), config.GetAWSRegion())
				} else {
					return fmt.Errorf("cannot find service for resource type %s", resType)
				}
			default:
				srv, err := cloud.GetServiceForType(resType)
				if err != nil {
					return err
				}
				fetchContext := context.WithValue(RootContext(), contextKey("force"), true)
				g, err = srv.FetchByType(context.WithValue(fetchContext, contextKey("filters"), listingFiltersFlag), resType)
				if err != nil {
					return err
				}
			}

			return printResources(g, resType)
		},
	}
}

var listAllResourceInServiceCmd = func(srvName string) *cobra.Command {
	return &cobra.Command{
		Use:    srvName,
		Short:  fmt.Sprintf("List all %s resources", srvName),
		Hidden: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			g := sync.LoadLocalGraphForService(srvName, config.GetAWSProfile(), config.GetAWSRegion())
			displayer, err := console.BuildOptions(
				console.WithFormat(listingFormat),
				console.WithMaxWidth(console.GetTerminalWidth()),
				console.WithIDsOnly(listOnlyIDs),
			).SetSource(g).Build()
			if err != nil {
				return err
			}
			if err := displayer.Print(os.Stdout); err != nil {
				return err
			}
			return nil
		},
	}
}

func printResources(g cloud.GraphAPI, resType string) error {
	displayer, err := console.BuildOptions(
		console.WithRdfType(resType),
		console.WithColumns(listingColumns(resType)),
		console.WithFilters(listingFiltersFlag),
		console.WithTagFilters(listingTagFiltersFlag),
		console.WithTagKeyFilters(listingTagKeyFiltersFlag),
		console.WithTagValueFilters(listingTagValueFiltersFlag),
		console.WithMaxWidth(console.GetTerminalWidth()),
		console.WithFormat(listingFormat),
		console.WithIDsOnly(listOnlyIDs),
		console.WithSortBy(sortBy...),
		console.WithReverseSort(reverseFlag),
		console.WithNoHeaders(noHeadersFlag),
	).SetSource(g).Build()
	if err != nil {
		return err
	}

	return displayer.Print(os.Stdout)
}

// listingColumns is the requested column set, with a region column added when listing
// across regions.
//
// The region is not decoration there: several resources are identified by a name rather
// than an ARN, so the same name in two regions produces two rows that are otherwise
// identical. It is appended rather than inserted so an explicit --columns still controls
// the order of everything the user asked for, and skipped when already present.
func listingColumns(resType string) []string {
	if !allLocalRegionsFlag {
		return listingColumnsFlag
	}

	columns := listingColumnsFlag
	if len(columns) == 0 {
		// Copy: this is the package-level default set, and appending to it would leak
		// the region column into later listings of the same type.
		columns = append([]string{}, console.ColumnsInListing[resType]...)
	}

	for _, c := range columns {
		if strings.EqualFold(c, properties.Region) {
			return columns
		}
	}
	return append(columns, properties.Region)
}

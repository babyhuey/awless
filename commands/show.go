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
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	awsconfig "github.com/bootswithdefer/awless/aws/config"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/cloud/rdf"
	"github.com/bootswithdefer/awless/config"
	"github.com/bootswithdefer/awless/console"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/sync"
)

var (
	listAllSiblingsFlag          bool
	noAliasFlag                  bool
	showPropertiesValuesOnlyFlag []string
)

func init() {
	RootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&listAllSiblingsFlag, "siblings", false, "List all the resource's siblings")
	showCmd.Flags().BoolVar(&noAliasFlag, "no-alias", false, "Disable the resolution of ID to alias")
	showCmd.Flags().StringSliceVar(&showPropertiesValuesOnlyFlag, "values-for", []string{}, "Output values only for given properties keys")
}

var showCmd = &cobra.Command{
	Use:   "show REFERENCE",
	Short: "Show resources lineage and dependencies given a REFERENCE: name, id, arn, etc...",
	Example: `  awless show i-8d43b21b            # show an instance via its ref
  awless show AIDAJ3Z24GOKHTZO4OIX6 # show a user via its ref
  awless show jsmith                # show a user via its ref
  awless show @jsmith               # forcing search by name
  awless show /aws/lambda/my-func   # show a log group by name
  awless show my-cluster            # show an EKS cluster
  awless show my-table              # show a DynamoDB table`,
	PersistentPreRunE:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRunE: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("REFERENCE required. See examples")
		}

		ref := args[0]
		notFound := fmt.Errorf("resource '%s' not found", deprefix(ref))

		if _, err := awsconfig.ParseRegion(ref); err == nil && ref != config.GetAWSRegion() {
			logger.Errorf("Cannot show region '%s' as you are in region '%s'", ref, config.GetAWSRegion())
			logger.Infof("Use `awless show %s -r %s`", ref, ref)
			return ErrReported
		}

		var resource cloud.Resource
		var gph cloud.GraphAPI
		var err error

		resource, gph, err = findResourceInLocalGraphs(ref)
		if err != nil {
			return err
		}

		if resource == nil && localGlobalFlag {
			if err := decorateWithSuggestion(notFound, ref); err != nil {
				return err
			}
		} else if resource == nil {
			runFullSync()

			if resource, gph, err = findResourceInLocalGraphs(ref); err != nil {
				return err
			} else if resource == nil {
				if err := decorateWithSuggestion(notFound, ref); err != nil {
					return err
				}
			}
		}

		if !localGlobalFlag && config.GetAutosync() {
			var services []cloud.Service
			if resource.Type() == cloud.Region {
				services = append(services, cloud.AllServices()...)
			} else {
				srv, err := cloud.GetServiceForType(resource.Type())
				if err != nil {
					return err
				}
				services = append(services, srv)
			}

			logger.Verbosef("syncing services for %s type", resource.Type())
			if _, err := sync.DefaultSyncer.Sync(RootContext(), services...); err != nil {
				logger.Verbose(err)
			}
			resource, gph, err = findResourceInLocalGraphs(ref)
			if err != nil {
				return err
			}
		}

		if resource != nil {
			if len(showPropertiesValuesOnlyFlag) > 0 {
				if err := showResourceValuesOnlyFor(resource, showPropertiesValuesOnlyFlag); err != nil {
					return err
				}
			} else if err := showResource(resource, gph); err != nil {
				return err
			}
		}

		return nil
	},
}

func showResourceValuesOnlyFor(resource cloud.Resource, propKeys []string) error {
	var normalized []string
	for _, p := range propKeys {
		normalized = append(normalized, strings.ToLower(strings.ReplaceAll(p, " ", "")))
	}

	valuesForKeys := map[string]string{}
	isIncluded := func(s string) (bool, string) {
		for _, n := range normalized {
			if n == strings.ToLower(s) {
				return true, n
			}
		}
		return false, ""
	}
	for k, v := range resource.Properties() {
		if ok, p := isIncluded(k); ok {
			valuesForKeys[p] = fmt.Sprint(v)
		}
	}

	var values []string
	for _, n := range normalized {
		if v, ok := valuesForKeys[n]; ok {
			values = append(values, v)
		}
	}

	if len(values) > 0 {
		fmt.Println(strings.Join(values, ","))
	} else {
		return fmt.Errorf("no values for %q", propKeys)
	}
	return nil
}

func showResource(resource cloud.Resource, gph cloud.GraphAPI) error {
	displayer, err := console.BuildOptions(
		console.WithColumnDefinitions(console.DefaultsColumnDefinitions[resource.Type()]),
		console.WithFormat(listingFormat),
		console.WithMaxWidth(console.GetTerminalWidth()),
	).SetSource(resource).Build()
	if err != nil {
		return err
	}

	if err := displayer.Print(os.Stdout); err != nil {
		return err
	}

	parents, err := gph.ResourceRelations(resource, rdf.ParentOf, true)
	if err != nil {
		return err
	}

	var parentsW bytes.Buffer
	var count int
	for i := len(parents) - 1; i >= 0; i-- {
		if count == 0 {
			fmt.Fprintf(&parentsW, "%s\n", printResourceRef(parents[i]))
		} else {
			fmt.Fprintf(&parentsW, "%s↳ %s\n", strings.Repeat("\t", count), printResourceRef(parents[i]))
		}
		count++
	}

	var childrenW bytes.Buffer
	var hasChildren bool
	printWithTabs := func(r cloud.Resource, distance int) error {
		var tabs bytes.Buffer
		tabs.WriteString(strings.Repeat("\t", count))
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		display := printResourceRef(r)
		if r.Same(resource) {
			display = printResourceRef(resource, renderGreenFn)
		} else {
			hasChildren = true
		}
		fmt.Fprintf(&childrenW, "%s↳ %s\n", tabs.String(), display)
		return nil
	}
	err = gph.VisitRelations(resource, rdf.ChildrenOfRel, true, printWithTabs)
	if err != nil {
		return err
	}

	if len(parents) > 0 || hasChildren {
		fmt.Println(renderCyanBoldFn("\nLineage:"))
		fmt.Printf("%s", parentsW.String())
		fmt.Printf("%s", childrenW.String())
	}

	appliedOn, err := gph.ResourceRelations(resource, rdf.ApplyOn, false)
	if err != nil {
		return err
	}
	printResourceList(renderCyanBoldFn("Applied on"), appliedOn)

	dependingOn, err := gph.ResourceRelations(resource, rdf.DependingOnRel, false)
	if err != nil {
		return err
	}
	printResourceList(renderCyanBoldFn("Depending on"), dependingOn)

	siblings, err := gph.ResourceSiblings(resource)
	if err != nil {
		return err
	}
	printResourceList(renderCyanBoldFn("Siblings"), siblings, "display all with flag --siblings")
	return nil
}

func runFullSync() {
	if !config.GetAutosync() {
		logger.Info("autosync disabled")
		return
	}

	logger.Infof("cannot find resource in existing data synced locally")
	logger.Infof("running sync for region '%s' and profile '%s'", config.GetAWSRegion(), config.GetAWSProfile())

	var services []cloud.Service
	for _, srv := range cloud.ServiceRegistry {
		services = append(services, srv)
	}

	if _, err := sync.DefaultSyncer.Sync(RootContext(), services...); err != nil {
		logger.Verbose(err)
	}
}

func findResourceInLocalGraphs(ref string) (cloud.Resource, cloud.GraphAPI, error) {
	g, resources, _, err := resolveResourceFromRefInCurrentRegion(ref)
	if err != nil {
		return nil, nil, err
	}
	switch len(resources) {
	case 0:
		return nil, nil, nil
	case 1:
		return resources[0], g, nil
	default:
		logger.Infof("%d resources found with name '%s' in region '%s' for profile '%s'. Show a specific resource with:", len(resources), deprefix(ref), config.GetAWSRegion(), config.GetAWSProfile())
		for _, res := range resources {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "\t`awless show %s` to show the %s", res.ID(), res.Type())
			if state, ok := res.Properties()[properties.State].(string); ok {
				fmt.Fprintf(&buf, " (state: '%s')", state)
			}
			logger.Infof("%s", buf.String())
		}

		// The user has been shown each candidate and the command to run; ending here
		// is success, not failure.
		return nil, nil, ErrExitZero
	}
}

func resolveResourceFromRefInCurrentRegion(ref string) (cloud.GraphAPI, []cloud.Resource, string, error) {
	g, err := sync.LoadLocalGraphs(config.GetAWSProfile(), config.GetAWSRegion())
	if err != nil {
		return nil, nil, "", err
	}
	return resolveResourceFromRef(g, ref)
}

func resolveResourceFromRefInAllLocalRegion(ref string) (cloud.GraphAPI, []cloud.Resource, string, error) {
	g, err := sync.LoadAllLocalGraphs(config.GetAWSProfile())
	if err != nil {
		return nil, nil, "", err
	}
	return resolveResourceFromRef(g, ref)
}

func resolveResourceFromRef(g cloud.GraphAPI, ref string) (cloud.GraphAPI, []cloud.Resource, string, error) {
	name := deprefix(ref)

	if strings.HasPrefix(ref, "@") {
		logger.Verbosef("prefixed with @: forcing research by name '%s'", name)
		rs, err := g.FindWithProperties(map[string]any{properties.Name: name})
		if err != nil {
			return nil, nil, "", err
		}
		return g, rs, properties.Name, nil
	}
	rs, err := g.FindWithProperties(map[string]any{properties.ID: name})
	if err != nil {
		return nil, nil, "", err
	}

	if len(rs) > 0 {
		return g, rs, properties.ID, nil
	}

	rs, err = g.FindWithProperties(map[string]any{properties.Arn: name})
	if err != nil {
		return nil, nil, "", err
	}

	if len(rs) > 0 {
		return g, rs, properties.Arn, nil
	}

	rs, err = g.FindWithProperties(map[string]any{properties.Name: name})
	if err != nil {
		return nil, nil, "", err
	}

	return g, rs, properties.Name, nil
}

func deprefix(s string) string {
	return strings.TrimPrefix(s, "@")
}

func decorateWithSuggestion(err error, ref string) error {
	buf := bytes.NewBufferString(fmt.Sprintf("%s in region '%s' for profile '%s'", err.Error(), config.GetAWSRegion(), config.GetAWSProfile()))
	g, resources, _, resolveErr := resolveResourceFromRefInAllLocalRegion(ref)
	if resolveErr != nil {
		// The suggestion is a nicety layered on the original error; failing to build
		// it must not replace what the caller was actually reporting.
		return err
	}
	for _, res := range resources {
		parents, err := g.ResourceRelations(res, rdf.ParentOf, true)
		if err != nil {
			return err
		}
		for _, parent := range parents {
			if parent.Type() == cloud.Region {
				fmt.Fprintf(buf, "\n\tfound previously synced under region '%s' as %s. Show it with `awless show %s -r %s --local`", parent.ID(), res, res.ID(), parent.ID())
			}
		}

	}
	return errors.New(buf.String())
}

func printResourceList(title string, list []cloud.Resource, shortenListMsg ...string) {
	sort.Sort(byTypeAndString{list})
	all := cloud.Resources(list).Map(func(r cloud.Resource) string { return printResourceRef(r) })
	count := len(all)
	max := 3
	if count > 0 {
		if !listAllSiblingsFlag && len(shortenListMsg) > 0 && count > max {
			fmt.Printf("\n%s: %s, ... (%s)\n", title, strings.Join(all[0:max], ", "), shortenListMsg[0])
		} else {
			fmt.Printf("\n%s: %s\n", title, strings.Join(all, ", "))
		}
	}
}

func printResourceRef(r cloud.Resource, idRenderFunc ...func(a ...any) string) string {
	render := fmt.Sprint
	if len(idRenderFunc) > 0 {
		render = idRenderFunc[0]
	}
	if noAliasFlag {
		return r.Format(render("%i") + "[" + color.New(color.FgBlue, color.Bold).SprintFunc()("%t") + "]")
	}
	return r.Format(render("%n") + "[" + color.New(color.FgBlue, color.Bold).SprintFunc()("%t") + "]")
}

type byTypeAndString struct {
	res []cloud.Resource
}

func (b byTypeAndString) Len() int { return len(b.res) }
func (b byTypeAndString) Swap(i, j int) {
	b.res[i], b.res[j] = b.res[j], b.res[i]
}
func (b byTypeAndString) Less(i, j int) bool {
	if b.res[i].Type() != b.res[j].Type() {
		return b.res[i].Type() < b.res[j].Type()
	}
	return b.res[i].String() <= b.res[j].String()
}

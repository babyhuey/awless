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

package sync

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	gosync "sync"
	"time"

	"runtime"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/sync/repo"
)

const fileExt = ".nt"

var DefaultSyncer Syncer

type Syncer interface {
	repo.Repo
	// Sync fetches the given services. ctx cancels the outbound AWS calls.
	Sync(ctx context.Context, services ...cloud.Service) (map[string]cloud.GraphAPI, error)
}

type noopsyncer struct {
	repo.NullRepo
}

func NoOpSyncer() Syncer { return new(noopsyncer) }

func (s *noopsyncer) Sync(context.Context, ...cloud.Service) (map[string]cloud.GraphAPI, error) {
	return map[string]cloud.GraphAPI{}, nil
}

type syncer struct {
	repo.Repo
	logger *logger.Logger
}

// NewSyncer builds a Syncer backed by the local git repo in ~/.awless.
//
// It returns an error rather than panicking: repo.New fails on a real runtime
// condition — an unwritable or unreadable ~/.awless — not a programmer error.
func NewSyncer(l ...*logger.Logger) (Syncer, error) {
	repo, err := repo.New()
	if err != nil {
		return nil, fmt.Errorf("creating syncer: %w", err)
	}

	s := &syncer{Repo: repo}

	if len(l) > 0 {
		s.logger = l[0]
	} else {
		s.logger = logger.DiscardLogger
	}

	return s, nil

}

func (s *syncer) Sync(ctx context.Context, services ...cloud.Service) (map[string]cloud.GraphAPI, error) {
	var workers gosync.WaitGroup

	type result struct {
		service cloud.Service
		gph     cloud.GraphAPI
		start   time.Time
		err     error
	}

	resultc := make(chan *result, len(services))

	for _, service := range services {
		if service.IsSyncDisabled() {
			s.logger.Verbosef("sync: *disabled* for service %s", service.Name())
			continue
		}
		workers.Add(1)
		go func(srv cloud.Service) {
			defer workers.Done()
			start := time.Now()
			g, err := srv.Fetch(ctx)
			resultc <- &result{service: srv, gph: g, start: start, err: err}
		}(service)
	}

	go func() {
		workers.Wait()
		close(resultc)
	}()

	var allErrors []error
	graphs := make(map[string]cloud.GraphAPI)
	servicesByName := make(map[string]cloud.Service)
	for res := range resultc {
		if res.err != nil {
			allErrors = append(allErrors, fmt.Errorf("syncing %s: %w", res.service.Name(), res.err))
		} else {
			s.logger.ExtraVerbosef("sync: fetched %s service took %s", res.service.Name(), time.Since(res.start))
		}
		if serv := res.service; serv != nil {
			servicesByName[serv.Name()] = serv
			if res.gph != nil {
				graphs[serv.Name()] = res.gph
			}
		}
	}

	var filepaths []string

	for name, g := range graphs {
		serviceRegion := servicesByName[name].Region()
		serviceProfile := servicesByName[name].Profile()
		serviceDir := filepath.Join(s.BaseDir(), serviceProfile, serviceRegion)
		if err := os.MkdirAll(serviceDir, 0700); err != nil {
			allErrors = append(allErrors, fmt.Errorf("creating %s: %w", serviceDir, err))
			continue
		}

		fullpath := filepath.Join(serviceDir, fmt.Sprintf("%s%s", name, fileExt))
		f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("opening %s: %w", fullpath, err))
			continue
		}
		closeFile := func() {
			if err := f.Close(); err != nil {
				allErrors = append(allErrors, fmt.Errorf("closing file %s: %w", fullpath, err))
			}
		}
		if err := g.MarshalTo(f); err != nil {
			allErrors = append(allErrors, fmt.Errorf("marshal to %s: %w", fullpath, err))
			closeFile()
			continue
		}
		relPath, err := filepath.Rel(s.BaseDir(), fullpath)
		if err != nil {
			allErrors = append(allErrors, err)
			closeFile()
			continue
		}

		filepaths = append(filepaths, relPath)
		closeFile()
	}

	if runtime.GOOS != "windows" { // https://github.com/bootswithdefer/awless/issues/119
		if err := s.Commit(filepaths...); err != nil {
			allErrors = append(allErrors, fmt.Errorf("committing %s: %w", strings.Join(filepaths, ", "), err))
		}
	}

	return graphs, concatErrors(allErrors)
}

func concatErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	lines := []string{"syncing errors:"}
	for _, err := range errs {
		lines = append(lines, fmt.Sprintf("\t\t%s", err))
	}

	return errors.New(strings.Join(lines, "\n"))
}

func LoadLocalGraphForService(serviceName, profile, region string) cloud.GraphAPI {
	regionDir := region
	if isGlobalService(serviceName) {
		regionDir = "global"
	}
	path := filepath.Join(repo.BaseDir(), profile, regionDir, fmt.Sprintf("%s%s", serviceName, fileExt))
	g, err := graph.NewGraphFromFiles(path)
	if err != nil {
		return graph.NewGraph()
	}
	return g
}

func LoadLocalGraphs(profile, region string) (cloud.GraphAPI, error) {
	var files []string
	globalFiles, _ := filepath.Glob(filepath.Join(repo.BaseDir(), profile, "global", fmt.Sprintf("*%s", fileExt)))
	regionFiles, _ := filepath.Glob(filepath.Join(repo.BaseDir(), profile, region, fmt.Sprintf("*%s", fileExt)))

	files = append(files, globalFiles...)
	files = append(files, regionFiles...)

	return graph.NewGraphFromFiles(files...)
}

func LoadAllLocalGraphs(profile string) (cloud.GraphAPI, error) {
	path := filepath.Join(repo.BaseDir(), profile, "*", fmt.Sprintf("*%s", fileExt))
	files, _ := filepath.Glob(path)

	return graph.NewGraphFromFiles(files...)
}

// LocalRegions lists the regions that have been synced locally for a profile, in sorted
// order. The "global" directory is not a region: it holds the services whose resources are
// not regional (IAM, Route53, CloudFront), and is reported separately by callers that need
// it.
func LocalRegions(profile string) []string {
	entries, err := os.ReadDir(filepath.Join(repo.BaseDir(), profile))
	if err != nil {
		return nil
	}

	var regions []string
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "global" || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		regions = append(regions, e.Name())
	}
	sort.Strings(regions)
	return regions
}

// LoadLocalGraphForTypeInAllRegions merges one resource type from every locally synced
// region into a single graph, recording on each resource the region it came from.
//
// The region has to be applied here rather than read back from the data, because nothing in
// a synced graph records it: the region is expressed only by which directory the file sits
// in. Without this, merging regions would produce a list with no way to tell an instance in
// us-east-1 from one in eu-west-1 — and, for resources whose ID is a name rather than an
// ARN, two same-named resources in different regions would collide into one row.
//
// Resources of a global service are returned with a region of "global", since asking for
// IAM users per region is meaningless.
func LoadLocalGraphForTypeInAllRegions(serviceName, resType, profile string) (cloud.GraphAPI, error) {
	combined := graph.NewGraph()

	regions := LocalRegions(profile)
	if isGlobalService(serviceName) {
		// A global service is stored once, under "global", so there is a single graph to
		// load and its resources are not regional.
		regions = []string{"global"}
	}

	for _, region := range regions {
		g := LoadLocalGraphForService(serviceName, profile, region)

		resources, err := g.Find(cloud.NewQuery(resType))
		if err != nil {
			return combined, fmt.Errorf("loading %s from region %s: %w", resType, region, err)
		}

		for _, r := range resources {
			// Find hands back the interface, but the concrete type is what can carry a
			// new property. Anything else is skipped rather than silently dropped from
			// the count, because it would indicate the graph package changed shape.
			res, ok := r.(*graph.Resource)
			if !ok {
				return combined, fmt.Errorf("loading %s from region %s: unexpected resource implementation %T", resType, region, r)
			}
			res.SetProperty(properties.Region, region)
			if err := combined.AddResource(res); err != nil {
				return combined, fmt.Errorf("merging %s from region %s: %w", resType, region, err)
			}
		}
	}

	return combined, nil
}

// isGlobalService reports whether a service's resources are stored under "global" rather
// than per region, because they are not regional in AWS.
func isGlobalService(serviceName string) bool {
	return serviceName == "access" || serviceName == "dns" || serviceName == "cdn"
}

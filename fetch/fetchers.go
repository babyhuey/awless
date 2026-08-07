package fetch

import (
	"context"
	"fmt"
	"sync"

	"github.com/bootswithdefer/awless/graph"
)

type Fetcher interface {
	Cache
	Fetch(context.Context) (*graph.Graph, error)
	FetchByType(context.Context, string) (*graph.Graph, error)
}

type Cache interface {
	Store(key string, val any)
	Get(key string, funcs ...func() (any, error)) (any, error)
	Reset()
}

type FetchResult struct {
	ResourceType string
	Err          error
	Resources    []*graph.Resource
	Objects      any
}

type Func func(context.Context, Cache) ([]*graph.Resource, any, error)

type Funcs map[string]Func

type fetcher struct {
	*cache
	fetchFuncs    map[string]Func
	resourceTypes []string
}

func NewFetcher(funcs Funcs) *fetcher {
	ftr := &fetcher{
		fetchFuncs: make(Funcs),
		cache:      newCache(),
	}
	for resType, f := range funcs {
		ftr.resourceTypes = append(ftr.resourceTypes, resType)
		ftr.fetchFuncs[resType] = f
	}
	return ftr
}

func (f *fetcher) Fetch(ctx context.Context) (*graph.Graph, error) {
	results := make(chan FetchResult, len(f.resourceTypes))
	var wg sync.WaitGroup

	for _, resType := range f.resourceTypes {
		wg.Add(1)
		go func(t string, co context.Context) {
			f.fetchResource(co, t, results)
			wg.Done()
		}(resType, ctx)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	gph := graph.NewGraph()

	ferr := new(Error)
	for res := range results {
		if err := res.Err; err != nil {
			ferr.Add(err)
		}
		gph.AddResource(res.Resources...)
	}

	if ferr.Any() {
		return gph, ferr
	}

	return gph, nil
}

// ContextKey is a typed key for context values to avoid SA1029.
type ContextKey string

// FetchModeKey is the context key used to indicate fetch-by-type mode.
const FetchModeKey ContextKey = "fetchmode"

func IsFetchingByType(c context.Context) (string, bool) {
	v, ok := c.Value(FetchModeKey).(string)
	return v, len(v) != 0 && ok
}

func (f *fetcher) FetchByType(ctx context.Context, resourceType string) (*graph.Graph, error) {
	results := make(chan FetchResult)
	defer close(results)

	go f.fetchResource(
		context.WithValue(ctx, FetchModeKey, resourceType),
		resourceType,
		results)

	gph := graph.NewGraph()
	res := <-results
	if err := res.Err; err != nil {
		return gph, err
	}
	for _, r := range res.Resources {
		gph.AddResource(r)
	}
	return gph, nil
}

func (f *fetcher) fetchResource(ctx context.Context, resourceType string, results chan<- FetchResult) {
	var err error
	var objects any
	resources := make([]*graph.Resource, 0)

	fn, ok := f.fetchFuncs[resourceType]
	if ok {
		resources, objects, err = fn(ctx, f.cache)
	} else {
		err = fmt.Errorf("no fetch func defined for resource type '%s'", resourceType)
	}

	f.cache.Store(fmt.Sprintf("%s_objects", resourceType), objects)

	results <- FetchResult{
		ResourceType: resourceType,
		Err:          err,
		Resources:    resources,
		Objects:      objects,
	}
}

type cache struct {
	mu     sync.RWMutex
	cached map[string]*keyCache
}

func newCache() *cache {
	return &cache{
		cached: make(map[string]*keyCache),
	}
}

type keyCache struct {
	once   sync.Once
	err    error
	result any
}

func (c *cache) Get(key string, funcs ...func() (any, error)) (any, error) {
	c.mu.Lock()
	cache, ok := c.cached[key]
	if !ok {
		cache = &keyCache{}
		c.cached[key] = cache
	}
	c.mu.Unlock()

	if len(funcs) > 0 {
		cache.once.Do(func() {
			cache.result, cache.err = funcs[0]()
		})
	}

	return cache.result, cache.err
}

func (c *cache) Store(key string, val any) {
	c.mu.Lock()
	c.cached[key] = &keyCache{result: val}
	c.mu.Unlock()
}

func (c *cache) Reset() {
	c.mu.Lock()
	c.cached = make(map[string]*keyCache)
	c.mu.Unlock()
}

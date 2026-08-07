package params

type Reducer interface {
	Keys() []string
	Reduce(map[string]any) (map[string]any, error)
}

type reduceFunc func(map[string]any) (map[string]any, error)

func newReducer(fn reduceFunc, keys ...string) Reducer {
	return &reducer{reduce: fn, keys: keys}
}

type reducer struct {
	keys   []string
	reduce reduceFunc
}

func (r *reducer) Keys() []string {
	return r.keys
}

func (r *reducer) Reduce(all map[string]any) (map[string]any, error) {
	in := make(map[string]any)
	for _, k := range r.keys {
		if v, ok := all[k]; ok {
			in[k] = v
		}
	}
	return r.reduce(in)
}

package harness

// Selector lazily resolves and caches a filtered set of objects from a Context.
// Results are cached per Context identity — if the Context changes (e.g., across
// separate Check calls), the selector re-evaluates.
type Selector[T any] struct {
	build    func(*Context) []T
	boundCtx *Context
	cache    []T
	resolved bool
}

func newSelector[T any](build func(*Context) []T) *Selector[T] {
	return &Selector[T]{build: build}
}

// Resolve returns the filtered objects, evaluating lazily on first call per Context.
func (s *Selector[T]) Resolve(ctx *Context) []T {
	if !s.resolved || s.boundCtx != ctx {
		s.cache = s.build(ctx)
		s.boundCtx = ctx
		s.resolved = true
	}
	return s.cache
}

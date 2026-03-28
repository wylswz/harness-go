package harness

import "github.com/wylswz/harness-go/internal/analysis"

// Context holds the loaded package graph and analysis caches.
// It is created by Check and threaded through selectors and rules.
type Context struct {
	graph *analysis.Context
}

func newContext(dir string, patterns ...string) (*Context, error) {
	ac, err := analysis.NewContext(dir, patterns...)
	if err != nil {
		return nil, err
	}
	return &Context{graph: ac}, nil
}

package query

import "github.com/wylswz/harness-go/pkg/analysis"

// Context holds the loaded package graph and analysis caches.
// It is created by Check and threaded through selectors and rules.
type Context struct {
	graph *analysis.Context
}

func ContextFrom(ac *analysis.Context) *Context {
	return &Context{graph: ac}
}

package analysis

import "golang.org/x/tools/go/packages"

// Context holds the complete analysis state for a set of loaded packages.
// It owns the Store (extracted objects) and retains the raw package graph
// for heavier analysis passes (SSA, call graph) that may be added later.
type Context struct {
	Store    *Store
	Packages []*packages.Package
}

// NewContext loads the given patterns and returns a fully populated Context.
// If dir is non-empty, it sets the working directory for the package loader.
func NewContext(dir string, patterns ...string) (*Context, error) {
	pkgs, store, err := Load(dir, patterns)
	if err != nil {
		return nil, err
	}
	return &Context{
		Store:    store,
		Packages: pkgs,
	}, nil
}

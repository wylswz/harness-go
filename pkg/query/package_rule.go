package query

import (
	"github.com/wylswz/harness-go/pkg/analysis"
	"github.com/wylswz/harness-go/pkg/object"
	"github.com/wylswz/harness-go/pkg/rule"
)

type packageImportRule struct {
	source    *PackageQuery
	forbidden *Selector[*object.PackageObj]
}

func (r *packageImportRule) Description() string {
	return "packages must not import forbidden packages"
}

func (r *packageImportRule) Check(ctx *analysis.Context) *rule.Result {
	// TODO: resolve source packages, resolve forbidden set, check imports
	hctx := ContextFrom(ctx)
	source := r.source.Select().Resolve(hctx)
	forbidden := r.forbidden.Resolve(hctx)
	for _, p := range source {
		for _, f := range forbidden {
			for _, importPath := range p.ImportPaths() {
				if importPath == f.PkgPath() {
					return &rule.Result{Violations: []rule.Violation{{RuleName: r.Description(), Message: "package imports forbidden package", Pos: p.Pos()}}}
				}
			}
		}
	}
	return &rule.Result{}
}

type packageAllowlistRule struct {
	source  *PackageQuery
	allowed *Selector[*object.PackageObj]
}

func (r *packageAllowlistRule) Description() string {
	return "packages must only import allowed packages"
}

func (r *packageAllowlistRule) Check(ctx *analysis.Context) *rule.Result {
	// TODO: resolve source packages, resolve allowed set, check imports
	return &rule.Result{}
}

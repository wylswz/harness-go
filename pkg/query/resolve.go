package query

import (
	"github.com/wylswz/harness-go/pkg/analysis"
	"github.com/wylswz/harness-go/pkg/object"
)

// Resolve helpers bridge analysis.Store to public object types.

func resolvePackages(ctx *Context, conds []Condition[*object.PackageObj]) []*object.PackageObj {
	var out []*object.PackageObj
	for _, info := range ctx.graph.Store.AllPackages() {
		obj := object.NewPackageFromInfo(info)
		if matchesAll(obj, conds) {
			out = append(out, obj)
		}
	}
	return out
}

func resolveFuncs(ctx *Context, pkgSel *Selector[*object.PackageObj], conds []Condition[*object.FuncObj]) []*object.FuncObj {
	var infos []analysis.FuncInfo

	if pkgSel != nil {
		pkgs := pkgSel.Resolve(ctx)
		for _, p := range pkgs {
			infos = append(infos, ctx.graph.Store.FuncsByPkg(p.PkgPath())...)
		}
	} else {
		infos = ctx.graph.Store.AllFuncs()
	}

	var out []*object.FuncObj
	for _, info := range infos {
		obj := object.NewFuncFromInfo(info)
		if matchesAll(obj, conds) {
			out = append(out, obj)
		}
	}
	return out
}

func resolveStructs(ctx *Context, pkgSel *Selector[*object.PackageObj], conds []Condition[*object.StructObj]) []*object.StructObj {
	var infos []analysis.StructInfo

	if pkgSel != nil {
		pkgs := pkgSel.Resolve(ctx)
		for _, p := range pkgs {
			infos = append(infos, ctx.graph.Store.StructsByPkg(p.PkgPath())...)
		}
	} else {
		infos = ctx.graph.Store.AllStructs()
	}

	var out []*object.StructObj
	for _, info := range infos {
		obj := object.NewStructFromInfo(info)
		if matchesAll(obj, conds) {
			out = append(out, obj)
		}
	}
	return out
}

func matchesAll[T any](obj T, conds []Condition[T]) bool {
	for _, c := range conds {
		if !c(obj) {
			return false
		}
	}
	return true
}

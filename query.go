package harness

import (
	"strings"

	"github.com/wylswz/harness-go/internal/analysis"
)

// ---------------------------------------------------------------------------
// PackageQuery
// ---------------------------------------------------------------------------

// PackageQuery builds a filtered view of packages.
// All filter methods return a new query (builders are immutable).
type PackageQuery struct {
	conditions []Condition[*PackageObj]
}

func Packages() *PackageQuery {
	return &PackageQuery{}
}

func (q *PackageQuery) clone() *PackageQuery {
	out := &PackageQuery{conditions: make([]Condition[*PackageObj], len(q.conditions))}
	copy(out.conditions, q.conditions)
	return out
}

func (q *PackageQuery) WithPrefix(prefix string) *PackageQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(p *PackageObj) bool {
		return strings.HasPrefix(p.PkgPath(), prefix)
	})
	return nq
}

func (q *PackageQuery) WithName(name string) *PackageQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(p *PackageObj) bool {
		return p.Name() == name
	})
	return nq
}

func (q *PackageQuery) Matching(cond Condition[*PackageObj]) *PackageQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, cond)
	return nq
}

func (q *PackageQuery) Select() *Selector[*PackageObj] {
	conds := q.conditions
	return newSelector(func(ctx *Context) []*PackageObj {
		return resolvePackages(ctx, conds)
	})
}

func (q *PackageQuery) MustNotImport(sel *Selector[*PackageObj]) Rule {
	return &packageImportRule{source: q, forbidden: sel}
}

func (q *PackageQuery) MustOnlyImport(sel *Selector[*PackageObj]) Rule {
	return &packageAllowlistRule{source: q, allowed: sel}
}

// ---------------------------------------------------------------------------
// FuncQuery
// ---------------------------------------------------------------------------

// FuncQuery builds a filtered view of functions and methods.
type FuncQuery struct {
	conditions []Condition[*FuncObj]
	residesIn  *Selector[*PackageObj]
}

func Functions() *FuncQuery {
	return &FuncQuery{}
}

func (q *FuncQuery) clone() *FuncQuery {
	out := &FuncQuery{
		conditions: make([]Condition[*FuncObj], len(q.conditions)),
		residesIn:  q.residesIn,
	}
	copy(out.conditions, q.conditions)
	return out
}

func (q *FuncQuery) ResidesIn(sel *Selector[*PackageObj]) *FuncQuery {
	nq := q.clone()
	nq.residesIn = sel
	return nq
}

func (q *FuncQuery) WithName(name string) *FuncQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(f *FuncObj) bool {
		return f.Name() == name
	})
	return nq
}

func (q *FuncQuery) Exported() *FuncQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(f *FuncObj) bool {
		return f.IsExported()
	})
	return nq
}

func (q *FuncQuery) Matching(cond Condition[*FuncObj]) *FuncQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, cond)
	return nq
}

func (q *FuncQuery) Select() *Selector[*FuncObj] {
	conds := q.conditions
	pkgs := q.residesIn
	return newSelector(func(ctx *Context) []*FuncObj {
		return resolveFuncs(ctx, pkgs, conds)
	})
}

func (q *FuncQuery) MustNotCall(sel *Selector[*FuncObj]) Rule {
	return &funcCallRule{source: q, forbidden: sel}
}

// ---------------------------------------------------------------------------
// StructQuery
// ---------------------------------------------------------------------------

// StructQuery builds a filtered view of struct types.
type StructQuery struct {
	conditions []Condition[*StructObj]
	residesIn  *Selector[*PackageObj]
}

func Structs() *StructQuery {
	return &StructQuery{}
}

func (q *StructQuery) clone() *StructQuery {
	out := &StructQuery{
		conditions: make([]Condition[*StructObj], len(q.conditions)),
		residesIn:  q.residesIn,
	}
	copy(out.conditions, q.conditions)
	return out
}

func (q *StructQuery) ResidesIn(sel *Selector[*PackageObj]) *StructQuery {
	nq := q.clone()
	nq.residesIn = sel
	return nq
}

func (q *StructQuery) WithName(name string) *StructQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(s *StructObj) bool {
		return s.Name() == name
	})
	return nq
}

func (q *StructQuery) Exported() *StructQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(s *StructObj) bool {
		return s.IsExported()
	})
	return nq
}

func (q *StructQuery) Matching(cond Condition[*StructObj]) *StructQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, cond)
	return nq
}

func (q *StructQuery) Select() *Selector[*StructObj] {
	conds := q.conditions
	pkgs := q.residesIn
	return newSelector(func(ctx *Context) []*StructObj {
		return resolveStructs(ctx, pkgs, conds)
	})
}

func (q *StructQuery) MustNotImport(sel *Selector[*PackageObj]) Rule {
	return &structImportRule{source: q, forbidden: sel}
}

func (q *StructQuery) MustNotEmbed(sel *Selector[*StructObj]) Rule {
	return &structEmbedRule{source: q, forbidden: sel}
}

// ---------------------------------------------------------------------------
// Resolve helpers — bridge internal analysis.Store to public object types
// ---------------------------------------------------------------------------

func packageFromInfo(info analysis.PackageInfo) *PackageObj {
	return &PackageObj{
		name:    info.Name,
		pkgPath: info.PkgPath,
		imports: info.ImportPaths,
	}
}

func funcFromInfo(info analysis.FuncInfo) *FuncObj {
	return &FuncObj{
		name:     info.Name,
		pkgPath:  info.PkgPath,
		pos:      info.Pos,
		exported: info.Exported,
		receiver: info.Receiver,
	}
}

func structFromInfo(info analysis.StructInfo) *StructObj {
	return &StructObj{
		name:     info.Name,
		pkgPath:  info.PkgPath,
		pos:      info.Pos,
		exported: info.Exported,
	}
}

func resolvePackages(ctx *Context, conds []Condition[*PackageObj]) []*PackageObj {
	var out []*PackageObj
	for _, info := range ctx.graph.Store.AllPackages() {
		obj := packageFromInfo(info)
		if matchesAll(obj, conds) {
			out = append(out, obj)
		}
	}
	return out
}

func resolveFuncs(ctx *Context, pkgSel *Selector[*PackageObj], conds []Condition[*FuncObj]) []*FuncObj {
	var infos []analysis.FuncInfo

	if pkgSel != nil {
		pkgs := pkgSel.Resolve(ctx)
		for _, p := range pkgs {
			infos = append(infos, ctx.graph.Store.FuncsByPkg(p.PkgPath())...)
		}
	} else {
		infos = ctx.graph.Store.AllFuncs()
	}

	var out []*FuncObj
	for _, info := range infos {
		obj := funcFromInfo(info)
		if matchesAll(obj, conds) {
			out = append(out, obj)
		}
	}
	return out
}

func resolveStructs(ctx *Context, pkgSel *Selector[*PackageObj], conds []Condition[*StructObj]) []*StructObj {
	var infos []analysis.StructInfo

	if pkgSel != nil {
		pkgs := pkgSel.Resolve(ctx)
		for _, p := range pkgs {
			infos = append(infos, ctx.graph.Store.StructsByPkg(p.PkgPath())...)
		}
	} else {
		infos = ctx.graph.Store.AllStructs()
	}

	var out []*StructObj
	for _, info := range infos {
		obj := structFromInfo(info)
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

// ---------------------------------------------------------------------------
// Rule implementations (stubs — assertion logic in later phases)
// ---------------------------------------------------------------------------

type packageImportRule struct {
	source    *PackageQuery
	forbidden *Selector[*PackageObj]
}

func (r *packageImportRule) Description() string { return "packages must not import forbidden packages" }
func (r *packageImportRule) Check(ctx *Context) *Result {
	// TODO: resolve source packages, resolve forbidden set, check imports
	return &Result{}
}

type packageAllowlistRule struct {
	source  *PackageQuery
	allowed *Selector[*PackageObj]
}

func (r *packageAllowlistRule) Description() string {
	return "packages must only import allowed packages"
}
func (r *packageAllowlistRule) Check(ctx *Context) *Result {
	// TODO: resolve source packages, resolve allowed set, check imports
	return &Result{}
}

type funcCallRule struct {
	source    *FuncQuery
	forbidden *Selector[*FuncObj]
}

func (r *funcCallRule) Description() string { return "functions must not call forbidden functions" }
func (r *funcCallRule) Check(ctx *Context) *Result {
	// TODO: resolve source functions, build call graph, check edges
	return &Result{}
}

type structImportRule struct {
	source    *StructQuery
	forbidden *Selector[*PackageObj]
}

func (r *structImportRule) Description() string {
	return "struct fields must not reference types from forbidden packages"
}
func (r *structImportRule) Check(ctx *Context) *Result {
	// TODO: resolve source structs, resolve forbidden packages, check field types
	return &Result{}
}

type structEmbedRule struct {
	source    *StructQuery
	forbidden *Selector[*StructObj]
}

func (r *structEmbedRule) Description() string { return "structs must not embed forbidden structs" }
func (r *structEmbedRule) Check(ctx *Context) *Result {
	// TODO: resolve source structs, resolve forbidden set, check embedded types
	return &Result{}
}

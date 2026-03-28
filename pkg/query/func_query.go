package query

import (
	"github.com/wylswz/harness-go/pkg/analysis"
	"github.com/wylswz/harness-go/pkg/object"
	"github.com/wylswz/harness-go/pkg/rule"
)

// FuncQuery builds a filtered view of functions and methods.
type FuncQuery struct {
	conditions []Condition[*object.FuncObj]
	residesIn  *Selector[*object.PackageObj]
}

// Functions returns a new empty function query.
func Functions() *FuncQuery {
	return &FuncQuery{}
}

func (q *FuncQuery) clone() *FuncQuery {
	out := &FuncQuery{
		conditions: make([]Condition[*object.FuncObj], len(q.conditions)),
		residesIn:  q.residesIn,
	}
	copy(out.conditions, q.conditions)
	return out
}

func (q *FuncQuery) ResidesIn(sel *Selector[*object.PackageObj]) *FuncQuery {
	nq := q.clone()
	nq.residesIn = sel
	return nq
}

func (q *FuncQuery) WithName(name string) *FuncQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(f *object.FuncObj) bool {
		return f.Name() == name
	})
	return nq
}

func (q *FuncQuery) Exported() *FuncQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(f *object.FuncObj) bool {
		return f.IsExported()
	})
	return nq
}

func (q *FuncQuery) Matching(cond Condition[*object.FuncObj]) *FuncQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, cond)
	return nq
}

func (q *FuncQuery) Select() *Selector[*object.FuncObj] {
	conds := q.conditions
	pkgs := q.residesIn
	return newSelector(func(ctx *Context) []*object.FuncObj {
		return resolveFuncs(ctx, pkgs, conds)
	})
}

func (q *FuncQuery) MustNotCall(sel *Selector[*object.FuncObj]) rule.Rule {
	return &funcCallRule{source: q, forbidden: sel}
}

type funcCallRule struct {
	source    *FuncQuery
	forbidden *Selector[*object.FuncObj]
}

func (r *funcCallRule) Description() string { return "functions must not call forbidden functions" }

func (r *funcCallRule) Check(ctx *analysis.Context) *rule.Result {
	// TODO: resolve source functions, build call graph, check edges
	return &rule.Result{}
}

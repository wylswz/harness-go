package query

import (
	"github.com/wylswz/harness-go/pkg/analysis"
	"github.com/wylswz/harness-go/pkg/object"
	"github.com/wylswz/harness-go/pkg/rule"
)

// StructQuery builds a filtered view of struct types.
type StructQuery struct {
	conditions []Condition[*object.StructObj]
	residesIn  *Selector[*object.PackageObj]
}

// Structs returns a new empty struct query.
func Structs() *StructQuery {
	return &StructQuery{}
}

func (q *StructQuery) clone() *StructQuery {
	out := &StructQuery{
		conditions: make([]Condition[*object.StructObj], len(q.conditions)),
		residesIn:  q.residesIn,
	}
	copy(out.conditions, q.conditions)
	return out
}

func (q *StructQuery) ResidesIn(sel *Selector[*object.PackageObj]) *StructQuery {
	nq := q.clone()
	nq.residesIn = sel
	return nq
}

func (q *StructQuery) WithName(name string) *StructQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(s *object.StructObj) bool {
		return s.Name() == name
	})
	return nq
}

func (q *StructQuery) Exported() *StructQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(s *object.StructObj) bool {
		return s.IsExported()
	})
	return nq
}

func (q *StructQuery) Matching(cond Condition[*object.StructObj]) *StructQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, cond)
	return nq
}

func (q *StructQuery) Select() *Selector[*object.StructObj] {
	conds := q.conditions
	pkgs := q.residesIn
	return newSelector(func(ctx *Context) []*object.StructObj {
		return resolveStructs(ctx, pkgs, conds)
	})
}

func (q *StructQuery) MustNotImport(sel *Selector[*object.PackageObj]) rule.Rule {
	return &structImportRule{source: q, forbidden: sel}
}

func (q *StructQuery) MustNotEmbed(sel *Selector[*object.StructObj]) rule.Rule {
	return &structEmbedRule{source: q, forbidden: sel}
}

type structImportRule struct {
	source    *StructQuery
	forbidden *Selector[*object.PackageObj]
}

func (r *structImportRule) Description() string {
	return "struct fields must not reference types from forbidden packages"
}

func (r *structImportRule) Check(ctx *analysis.Context) *rule.Result {
	// TODO: resolve source structs, resolve forbidden packages, check field types
	return &rule.Result{}
}

type structEmbedRule struct {
	source    *StructQuery
	forbidden *Selector[*object.StructObj]
}

func (r *structEmbedRule) Description() string { return "structs must not embed forbidden structs" }

func (r *structEmbedRule) Check(ctx *analysis.Context) *rule.Result {
	// TODO: resolve source structs, resolve forbidden set, check embedded types
	return &rule.Result{}
}

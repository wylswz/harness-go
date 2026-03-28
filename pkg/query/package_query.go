package query

import (
	"github.com/wylswz/harness-go/pkg/object"
	"github.com/wylswz/harness-go/pkg/rule"
)

// PackageQuery builds a filtered view of packages.
// All filter methods return a new query (builders are immutable).
type PackageQuery struct {
	conditions []Condition[*object.PackageObj]
}

// Packages returns a new empty package query.
func Packages() *PackageQuery {
	return &PackageQuery{}
}

func (q *PackageQuery) clone() *PackageQuery {
	out := &PackageQuery{conditions: make([]Condition[*object.PackageObj], len(q.conditions))}
	copy(out.conditions, q.conditions)
	return out
}

func (q *PackageQuery) WithPrefix(prefix string) *PackageQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(p *object.PackageObj) bool {
		return p.HasPrefix(prefix)
	})
	return nq
}

func (q *PackageQuery) WithName(name string) *PackageQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, func(p *object.PackageObj) bool {
		return p.Name() == name
	})
	return nq
}

func (q *PackageQuery) Matching(cond Condition[*object.PackageObj]) *PackageQuery {
	nq := q.clone()
	nq.conditions = append(nq.conditions, cond)
	return nq
}

func (q *PackageQuery) Select() *Selector[*object.PackageObj] {
	conds := q.conditions
	return newSelector(func(ctx *Context) []*object.PackageObj {
		return resolvePackages(ctx, conds)
	})
}

func (q *PackageQuery) MustNotImport(sel *Selector[*object.PackageObj]) rule.Rule {
	return &packageImportRule{source: q, forbidden: sel}
}

func (q *PackageQuery) MustOnlyImport(sel *Selector[*object.PackageObj]) rule.Rule {
	return &packageAllowlistRule{source: q, allowed: sel}
}

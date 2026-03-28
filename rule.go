package harness

import (
	"fmt"
	"go/token"
	"strings"
)

// Rule is an architectural constraint that can be checked against a Context.
type Rule interface {
	Check(ctx *Context) *Result
	Description() string
}

// Violation represents a single architectural constraint breach.
type Violation struct {
	RuleName string
	Message  string
	Pos      token.Position
	Symbols  []string
}

func (v Violation) String() string {
	loc := v.Pos.String()
	if !v.Pos.IsValid() {
		loc = "<unknown>"
	}
	return fmt.Sprintf("%s: %s [%s]", loc, v.Message, v.RuleName)
}

// Result holds the outcome of checking one Rule.
type Result struct {
	Violations []Violation
}

func (r *Result) Failed() bool {
	return len(r.Violations) > 0
}

func (r *Result) String() string {
	if !r.Failed() {
		return "OK"
	}
	var b strings.Builder
	for _, v := range r.Violations {
		b.WriteString(v.String())
		b.WriteByte('\n')
	}
	return b.String()
}

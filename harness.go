package harness

import (
	"testing"

	"github.com/wylswz/harness-go/internal/analysis"
	"github.com/wylswz/harness-go/pkg/rule"
)

// Check evaluates architectural rules and reports violations as test failures.
// Packages are loaded from the patterns in HARNESS_PACKAGES (comma-separated)
// or "./..." by default.
func Check(t *testing.T, ctx *analysis.Context, rules ...rule.Rule) []*rule.Violation {
	t.Helper()
	violations := make([]*rule.Violation, 0)

	for _, rule := range rules {
		result := rule.Check(ctx)
		if result.Failed() {
			for _, v := range result.Violations {
				violations = append(violations, &v)
			}
		}
	}
	return violations
}

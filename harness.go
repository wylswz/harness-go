package harness

import (
	"os"
	"strings"
	"testing"
)

// Check evaluates architectural rules and reports violations as test failures.
// Packages are loaded from the patterns in HARNESS_PACKAGES (comma-separated)
// or "./..." by default.
func Check(t *testing.T, rules ...Rule) {
	t.Helper()

	patterns := []string{"./..."}
	if env := os.Getenv("HARNESS_PACKAGES"); env != "" {
		patterns = strings.Split(env, ",")
	}

	ctx, err := newContext("", patterns...)
	if err != nil {
		t.Fatalf("harness: failed to load packages: %v", err)
	}

	for _, rule := range rules {
		result := rule.Check(ctx)
		if result.Failed() {
			for _, v := range result.Violations {
				t.Errorf("%s", v.String())
			}
		}
	}
}

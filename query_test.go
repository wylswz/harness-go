package harness

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wylswz/harness-go/internal/analysis"
	"github.com/wylswz/harness-go/internal/testdata"
)

func TestPackageRules(t *testing.T) {
	ctx, err := analysis.NewContext(testdata.Root(), testdata.Patterns()...)
	if err != nil {
		t.Fatalf("harness: failed to load packages: %v", err)
	}
	t.Run("package restrictions", func(t *testing.T) {

		violations := Check(t, ctx, Packages().WithPrefix(
			testdata.PackagePrefix("user"),
		).MustNotImport(
			Packages().WithPrefix(
				testdata.PackagePrefix("file", "repo"),
			).Select(),
		))
		assert.Equal(t, 1, len(violations))
	},
	)
}

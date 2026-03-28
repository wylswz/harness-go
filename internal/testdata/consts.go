package testdata

import (
	"os"
)

// Root returns the module directory for the minimal testdata project
// (the folder containing testdata/go.mod). Tests set TEST_DATA_ROOT automatically.
func Root() string {
	root := os.Getenv("TEST_DATA_ROOT")
	if root == "" {
		panic("TEST_DATA_ROOT is not set")
	}
	return root
}

func Patterns() []string {
	return []string{PackagePrefix("...")}
}

func PackagePrefix(children ...string) string {
	prefix := "github.com/wylswz/harness-go/testdata/pkg/"
	for _, child := range children {
		prefix += child + "/"
	}
	return prefix
}

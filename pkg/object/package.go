package object

import (
	"go/token"
	"strings"

	"github.com/wylswz/harness-go/pkg/analysis"
)

// PackageObj represents a Go package.
type PackageObj struct {
	name    string
	pkgPath string
	imports []string
}

// NewPackage returns a package object with the given name, import path, and import paths.
func NewPackage(name, pkgPath string, imports []string) *PackageObj {
	return &PackageObj{name: name, pkgPath: pkgPath, imports: imports}
}

// NewPackageFromInfo builds a PackageObj from analysis metadata.
func NewPackageFromInfo(info analysis.PackageInfo) *PackageObj {
	return NewPackage(info.Name, info.PkgPath, info.ImportPaths)
}

func (p *PackageObj) Kind() ObjectKind      { return KindPackage }
func (p *PackageObj) Name() string          { return p.name }
func (p *PackageObj) PkgPath() string       { return p.pkgPath }
func (p *PackageObj) Pos() token.Position   { return token.Position{} }
func (p *PackageObj) ImportPaths() []string { return p.imports }
func (p *PackageObj) HasPrefix(prefix string) bool {

	return strings.HasPrefix(p.pkgPath+"/", prefix)
}

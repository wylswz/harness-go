package analysis

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"

	"github.com/wylswz/harness-go/internal/logging"
	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedSyntax

// Load loads Go packages matching the given patterns and extracts all
// packages, functions, and structs into a Store.
// If dir is non-empty, it is used as the working directory for the
// package loader (required when loading a foreign module).
func Load(dir string, patterns []string) ([]*packages.Package, *Store, error) {
	cfg := &packages.Config{
		Mode: loadMode,
		Dir:  dir,
	}

	logging.Debug("loading packages", "dir", dir, "patterns", patterns)

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, nil, fmt.Errorf("packages.Load: %w", err)
	}

	logging.Debug("initializing store")

	store := NewStore()

	seen := make(map[string]bool)
	packages.Visit(pkgs, func(pkg *packages.Package) bool {
		if len(pkg.Errors) > 0 {
			logging.Error("package errors", "package", pkg.PkgPath, "errors", pkg.Errors)
			return false
		}
		if seen[pkg.PkgPath] {
			return false
		}
		return true

	}, func(pkg *packages.Package) {
		logging.Debug("extracting package", "package", pkg.PkgPath)
		seen[pkg.PkgPath] = true
		extractPackage(store, pkg)
	})

	return pkgs, store, nil
}

func extractPackage(store *Store, pkg *packages.Package) {
	importPaths := make([]string, 0, len(pkg.Imports))
	for path := range pkg.Imports {
		importPaths = append(importPaths, path)
	}
	sort.Strings(importPaths)

	store.AddPackage(PackageInfo{
		Name:        pkg.Name,
		PkgPath:     pkg.PkgPath,
		ImportPaths: importPaths,
	})

	for _, file := range pkg.Syntax {
		extractFile(store, pkg, file)
	}
}

func extractFile(store *Store, pkg *packages.Package, file *ast.File) {
	logging.Debug("extracting file", "file", file.Name.Name, "pkg", pkg.PkgPath)
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			extractFunc(store, pkg.PkgPath, pkg.Fset, d)
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if _, isStruct := ts.Type.(*ast.StructType); isStruct {
					extractStruct(store, pkg, ts)
				}
			}
		}
	}
}

func extractFunc(store *Store, pkgPath string, fset *token.FileSet, fd *ast.FuncDecl) {
	var receiver string
	if fd.Recv != nil && len(fd.Recv.List) > 0 {
		receiver = exprString(fd.Recv.List[0].Type)
	}
	store.AddFunc(FuncInfo{
		Name:     fd.Name.Name,
		PkgPath:  pkgPath,
		Pos:      fset.Position(fd.Name.Pos()),
		Exported: fd.Name.IsExported(),
		Receiver: receiver,
	})
}

func extractStruct(store *Store, pkg *packages.Package, ts *ast.TypeSpec) {
	store.AddStruct(StructInfo{
		Name:     ts.Name.Name,
		PkgPath:  pkg.PkgPath,
		Pos:      pkg.Fset.Position(ts.Name.Pos()),
		Exported: ts.Name.IsExported(),
	})
}

// exprString returns a short string representation of a type expression
// (used for receiver types like *T or T).
func exprString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + exprString(e.X)
	case *ast.IndexExpr:
		return exprString(e.X)
	case *ast.IndexListExpr:
		return exprString(e.X)
	default:
		return ""
	}
}

package analysis

import "go/token"

// PackageInfo holds extracted metadata for a single Go package.
type PackageInfo struct {
	Name        string
	PkgPath     string
	ImportPaths []string
}

// FuncInfo holds extracted metadata for a function or method declaration.
type FuncInfo struct {
	Name     string
	PkgPath  string
	Pos      token.Position
	Exported bool
	Receiver string // empty for top-level functions
}

// StructInfo holds extracted metadata for a struct type declaration.
type StructInfo struct {
	Name     string
	PkgPath  string
	Pos      token.Position
	Exported bool
}

// Store holds indexed code elements extracted from a loaded package graph.
// Elements are stored by kind, with secondary indexes on package path for
// efficient scoped lookups (e.g. "all functions in package X").
type Store struct {
	pkgs    []PackageInfo
	funcs   []FuncInfo
	structs []StructInfo

	funcsByPkg   map[string][]int // pkgPath → indices into funcs
	structsByPkg map[string][]int // pkgPath → indices into structs
}

func NewStore() *Store {
	return &Store{
		funcsByPkg:   make(map[string][]int),
		structsByPkg: make(map[string][]int),
	}
}

func (s *Store) AddPackage(info PackageInfo) {
	s.pkgs = append(s.pkgs, info)
}

func (s *Store) AddFunc(info FuncInfo) {
	idx := len(s.funcs)
	s.funcs = append(s.funcs, info)
	s.funcsByPkg[info.PkgPath] = append(s.funcsByPkg[info.PkgPath], idx)
}

func (s *Store) AddStruct(info StructInfo) {
	idx := len(s.structs)
	s.structs = append(s.structs, info)
	s.structsByPkg[info.PkgPath] = append(s.structsByPkg[info.PkgPath], idx)
}

func (s *Store) AllPackages() []PackageInfo {
	return s.pkgs
}

func (s *Store) AllFuncs() []FuncInfo {
	return s.funcs
}

// FuncsByPkg returns functions declared in the given package.
func (s *Store) FuncsByPkg(pkgPath string) []FuncInfo {
	indices := s.funcsByPkg[pkgPath]
	out := make([]FuncInfo, len(indices))
	for i, idx := range indices {
		out[i] = s.funcs[idx]
	}
	return out
}

func (s *Store) AllStructs() []StructInfo {
	return s.structs
}

// StructsByPkg returns structs declared in the given package.
func (s *Store) StructsByPkg(pkgPath string) []StructInfo {
	indices := s.structsByPkg[pkgPath]
	out := make([]StructInfo, len(indices))
	for i, idx := range indices {
		out[i] = s.structs[idx]
	}
	return out
}

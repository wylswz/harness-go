package harness

import "go/token"

type ObjectKind int

const (
	KindPackage ObjectKind = iota
	KindFunc
	KindStruct
)

func (k ObjectKind) String() string {
	switch k {
	case KindPackage:
		return "package"
	case KindFunc:
		return "function"
	case KindStruct:
		return "struct"
	default:
		return "unknown"
	}
}

// Object is a Go code element that can be queried and constrained.
type Object interface {
	Kind() ObjectKind
	Name() string
	PkgPath() string
	Pos() token.Position
}

// PackageObj represents a Go package.
type PackageObj struct {
	name    string
	pkgPath string
	imports []string
}

func (p *PackageObj) Kind() ObjectKind     { return KindPackage }
func (p *PackageObj) Name() string         { return p.name }
func (p *PackageObj) PkgPath() string      { return p.pkgPath }
func (p *PackageObj) Pos() token.Position  { return token.Position{} }
func (p *PackageObj) ImportPaths() []string { return p.imports }

// FuncObj represents a Go function or method.
type FuncObj struct {
	name     string
	pkgPath  string
	pos      token.Position
	exported bool
	receiver string
}

func (f *FuncObj) Kind() ObjectKind    { return KindFunc }
func (f *FuncObj) Name() string        { return f.name }
func (f *FuncObj) PkgPath() string     { return f.pkgPath }
func (f *FuncObj) Pos() token.Position { return f.pos }
func (f *FuncObj) IsExported() bool    { return f.exported }
func (f *FuncObj) IsMethod() bool      { return f.receiver != "" }
func (f *FuncObj) Receiver() string    { return f.receiver }

// StructObj represents a Go struct type.
type StructObj struct {
	name     string
	pkgPath  string
	pos      token.Position
	exported bool
}

func (s *StructObj) Kind() ObjectKind    { return KindStruct }
func (s *StructObj) Name() string        { return s.name }
func (s *StructObj) PkgPath() string     { return s.pkgPath }
func (s *StructObj) Pos() token.Position { return s.pos }
func (s *StructObj) IsExported() bool    { return s.exported }

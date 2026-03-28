package object

import (
	"go/token"

	"github.com/wylswz/harness-go/pkg/analysis"
)

// StructObj represents a Go struct type.
type StructObj struct {
	name     string
	pkgPath  string
	pos      token.Position
	exported bool
}

// NewStruct returns a struct type object.
func NewStruct(name, pkgPath string, pos token.Position, exported bool) *StructObj {
	return &StructObj{
		name:     name,
		pkgPath:  pkgPath,
		pos:      pos,
		exported: exported,
	}
}

// NewStructFromInfo builds a StructObj from analysis metadata.
func NewStructFromInfo(info analysis.StructInfo) *StructObj {
	return NewStruct(info.Name, info.PkgPath, info.Pos, info.Exported)
}

func (s *StructObj) Kind() ObjectKind    { return KindStruct }
func (s *StructObj) Name() string        { return s.name }
func (s *StructObj) PkgPath() string     { return s.pkgPath }
func (s *StructObj) Pos() token.Position { return s.pos }
func (s *StructObj) IsExported() bool    { return s.exported }

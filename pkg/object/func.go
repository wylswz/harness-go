package object

import (
	"go/token"

	"github.com/wylswz/harness-go/pkg/analysis"
)

// FuncObj represents a Go function or method.
type FuncObj struct {
	name     string
	pkgPath  string
	pos      token.Position
	exported bool
	receiver string
}

// NewFunc returns a function or method object.
func NewFunc(name, pkgPath string, pos token.Position, exported bool, receiver string) *FuncObj {
	return &FuncObj{
		name:     name,
		pkgPath:  pkgPath,
		pos:      pos,
		exported: exported,
		receiver: receiver,
	}
}

// NewFuncFromInfo builds a FuncObj from analysis metadata.
func NewFuncFromInfo(info analysis.FuncInfo) *FuncObj {
	return NewFunc(info.Name, info.PkgPath, info.Pos, info.Exported, info.Receiver)
}

func (f *FuncObj) Kind() ObjectKind    { return KindFunc }
func (f *FuncObj) Name() string        { return f.name }
func (f *FuncObj) PkgPath() string     { return f.pkgPath }
func (f *FuncObj) Pos() token.Position { return f.pos }
func (f *FuncObj) IsExported() bool    { return f.exported }
func (f *FuncObj) IsMethod() bool      { return f.receiver != "" }
func (f *FuncObj) Receiver() string    { return f.receiver }

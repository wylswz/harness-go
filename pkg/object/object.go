package object

import "go/token"

// Object is a Go code element that can be queried and constrained.
type Object interface {
	Kind() ObjectKind
	Name() string
	PkgPath() string
	Pos() token.Position
}

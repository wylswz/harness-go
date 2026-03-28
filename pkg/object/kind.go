package object

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

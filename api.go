package harness

import "github.com/wylswz/harness-go/pkg/query"

func Packages() *query.PackageQuery {
	return query.Packages()
}

func Functions() *query.FuncQuery {
	return query.Functions()
}

func Structs() *query.StructQuery {
	return query.Structs()
}

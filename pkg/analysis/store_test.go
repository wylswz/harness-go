package analysis

import (
	"go/token"
	"testing"
)

func TestStore_AddAndRetrieve(t *testing.T) {
	s := NewStore()

	s.AddPackage(PackageInfo{Name: "foo", PkgPath: "example.com/foo", ImportPaths: []string{"fmt"}})
	s.AddPackage(PackageInfo{Name: "bar", PkgPath: "example.com/bar", ImportPaths: nil})

	s.AddFunc(FuncInfo{Name: "DoFoo", PkgPath: "example.com/foo", Exported: true})
	s.AddFunc(FuncInfo{Name: "helper", PkgPath: "example.com/foo", Exported: false})
	s.AddFunc(FuncInfo{Name: "DoBar", PkgPath: "example.com/bar", Exported: true})

	s.AddStruct(StructInfo{Name: "Foo", PkgPath: "example.com/foo", Exported: true})
	s.AddStruct(StructInfo{Name: "bar", PkgPath: "example.com/bar", Exported: false})

	if got := len(s.AllPackages()); got != 2 {
		t.Errorf("AllPackages: got %d, want 2", got)
	}
	if got := len(s.AllFuncs()); got != 3 {
		t.Errorf("AllFuncs: got %d, want 3", got)
	}
	if got := len(s.AllStructs()); got != 2 {
		t.Errorf("AllStructs: got %d, want 2", got)
	}
}

func TestStore_FuncsByPkg(t *testing.T) {
	s := NewStore()
	s.AddFunc(FuncInfo{Name: "A", PkgPath: "p1"})
	s.AddFunc(FuncInfo{Name: "B", PkgPath: "p2"})
	s.AddFunc(FuncInfo{Name: "C", PkgPath: "p1"})

	got := s.FuncsByPkg("p1")
	if len(got) != 2 {
		t.Fatalf("FuncsByPkg(p1): got %d, want 2", len(got))
	}

	names := map[string]bool{got[0].Name: true, got[1].Name: true}
	if !names["A"] || !names["C"] {
		t.Errorf("expected A and C, got %v", names)
	}

	got = s.FuncsByPkg("p2")
	if len(got) != 1 || got[0].Name != "B" {
		t.Errorf("FuncsByPkg(p2): unexpected result %v", got)
	}

	got = s.FuncsByPkg("nonexistent")
	if len(got) != 0 {
		t.Errorf("FuncsByPkg(nonexistent): expected empty, got %d", len(got))
	}
}

func TestStore_StructsByPkg(t *testing.T) {
	s := NewStore()
	s.AddStruct(StructInfo{Name: "X", PkgPath: "p1", Pos: token.Position{Line: 10}})
	s.AddStruct(StructInfo{Name: "Y", PkgPath: "p2"})
	s.AddStruct(StructInfo{Name: "Z", PkgPath: "p1"})

	got := s.StructsByPkg("p1")
	if len(got) != 2 {
		t.Fatalf("StructsByPkg(p1): got %d, want 2", len(got))
	}

	if got[0].Pos.Line != 10 {
		t.Errorf("expected first struct to have Line=10, got %d", got[0].Pos.Line)
	}
}

# Harness-Go Implementation TODO

All public types live in the top-level `package harness`. Internal helpers go in `internal/` subpackages as needed.

## Phase 0 — Foundation

- [x] Fix module path in `go.mod` to `github.com/wylswz/harness-go`
- [x] Create top-level `package harness` with skeleton types
- [x] Define `Object` interface, `PackageObj`, `FuncObj`, `StructObj`
- [x] Define generic `Condition[T]` with `And`, `Or`, `Not` combinators
- [x] Define generic `Selector[T]` with lazy per-context caching
- [x] Define `Rule` interface, `Violation`, `Result`
- [x] Define immutable query builders: `PackageQuery`, `FuncQuery`, `StructQuery`
- [x] Define `Check(t, rules...)` entry point
- [ ] Add `golang.org/x/tools/go/packages` as a dependency
- [ ] Write a smoke test that loads the current module and asserts at least one package is found

## Phase 1 — Context & Package Loading

- [ ] Add `[]*packages.Package` and lookup maps to `Context`
- [ ] Implement `Context.Load()` using `go/packages.Load` with `NeedName | NeedFiles | NeedImports | NeedDeps | NeedTypes | NeedSyntax`
- [ ] Extract `[]PackageObj` from loaded packages
- [ ] Extract `[]FuncObj` by walking `ast.FuncDecl` nodes per package
- [ ] Extract `[]StructObj` by walking `ast.TypeSpec` (where type is `*ast.StructType`) per package
- [ ] Implement `PackageQuery.Select()` resolve function: filter packages by conditions
- [ ] Implement `FuncQuery.Select()` resolve function: filter by `residesIn` + conditions
- [ ] Implement `StructQuery.Select()` resolve function: filter by `residesIn` + conditions
- [ ] Write unit tests for each selector with a `testdata/` fixture module

## Phase 2 — Package Import Rules

- [ ] Implement `packageImportRule.Check`: resolve source packages and forbidden set, compare import paths
- [ ] Implement `packageAllowlistRule.Check`: resolve source packages and allowed set, flag imports outside the set
- [ ] Create `testdata/importrule/` fixture module with known import violations
- [ ] Write tests: passing rules, failing rules, verify `Violation` positions and messages

## Phase 3 — Function Call Rules

- [ ] Add `golang.org/x/tools/go/ssa` and `golang.org/x/tools/go/callgraph` dependencies
- [ ] Lazily build SSA program and call graph (VTA algorithm) in `Context`
- [ ] Implement `funcCallRule.Check`: resolve source functions and forbidden set, walk call graph edges
- [ ] Create `testdata/callrule/` fixture module with known call violations
- [ ] Write tests: direct calls, method calls, interface dispatch (document precision limits)

## Phase 4 — Struct Dependency Rules

- [ ] Implement `structImportRule.Check`: resolve source structs, walk field types, check package origin
- [ ] Implement `structEmbedRule.Check`: resolve source structs, check embedded types against forbidden set
- [ ] Create `testdata/structrule/` fixture module and tests

## Phase 5 — Ergonomics & Polish

- [ ] Write `example_test.go` demonstrating real-world usage patterns
- [ ] Ensure `go test ./...` in a host project is the only command needed
- [ ] Improve rule description strings to include query details (prefix, name, etc.)
- [ ] Add `WithPathMatching(regex)` to `PackageQuery` for regex-based path filtering
- [ ] Add `Methods()` entry point or `FuncQuery.MethodsOnly()` filter

## Phase 6 — Performance & Caching

- [ ] Benchmark selector evaluation on a medium-size module (~50 packages)
- [ ] Profile call-graph construction; document recommended load flags per rule type
- [ ] Consider incremental analysis: only re-analyze changed packages

## Phase 7 — Optional CLI

- [ ] `cmd/harness` binary that runs rules from Go test files
- [ ] Exit code 1 on violations, structured JSON output with `--json`
- [ ] GitHub Actions integration example in README

## Cross-Cutting Concerns

- [ ] Every `Violation` must include file, line, symbol name, and rule text
- [ ] Godoc on every exported type and function
- [ ] Exported vs unexported symbols: `Exported()` / `Matching()` covers this — document patterns
- [ ] Build tags / OS constraints: document how `go/packages` load patterns interact with tags

# harness-go

A Go structural testing framework that enforces architectural constraints as regular `go test` assertions.

Write rules to ensure packages don't import what they shouldn't, functions don't call what they shouldn't, and structs don't depend on what they shouldn't — all checked at test time using static analysis.

## Installation

```bash
go get github.com/wylswz/harness-go
```

## Quick Start

Create an architecture test file anywhere in your project:

```go
package myproject_test

import (
    "testing"

    h "github.com/wylswz/harness-go"
)

func TestArchitecture(t *testing.T) {
    h.Check(t,
        // API layer must not import internal packages
        h.Packages().WithPrefix("github.com/myorg/myproject/api").MustNotImport(
            h.Packages().WithPrefix("github.com/myorg/myproject/internal").Select(),
        ),
    )
}
```

Run it:

```bash
go test ./...
```

## API Examples

### Package import constraints

```go
// Prevent circular or layered dependency violations
h.Check(t,
    h.Packages().WithPrefix("github.com/myorg/myproject/domain").MustNotImport(
        h.Packages().WithPrefix("github.com/myorg/myproject/infra").Select(),
    ),
)
```

### Reusable selectors

Selectors are lazy and cached — build them once and share across rules:

```go
internalPkgs := h.Packages().WithPrefix("github.com/myorg/myproject/internal").Select()
apiPkgs := h.Packages().WithPrefix("github.com/myorg/myproject/api").Select()

h.Check(t,
    h.Packages().WithPrefix("github.com/myorg/myproject/api").MustNotImport(internalPkgs),
    h.Structs().ResidesIn(apiPkgs).MustNotImport(internalPkgs),
)
```

### Function call constraints

```go
internalPkgs := h.Packages().WithPrefix("github.com/myorg/myproject/internal").Select()

h.Check(t,
    // Functions in internal packages must not call DeprecatedHelper
    h.Functions().ResidesIn(internalPkgs).MustNotCall(
        h.Functions().WithName("DeprecatedHelper").Select(),
    ),
)
```

### Struct embedding constraints

```go
h.Check(t,
    // Domain structs must not embed ORM base models
    h.Structs().ResidesIn(
        h.Packages().WithPrefix("github.com/myorg/myproject/domain").Select(),
    ).MustNotEmbed(
        h.Structs().WithName("BaseModel").Select(),
    ),
)
```

### Custom conditions

Use `Matching` to apply arbitrary predicates:

```go
h.Check(t,
    // Exported functions in the API layer must not call any unexported helpers in infra
    h.Functions().ResidesIn(apiPkgs).Exported().MustNotCall(
        h.Functions().ResidesIn(infraPkgs).Matching(func(f *h.FuncObj) bool {
            return !f.IsExported()
        }).Select(),
    ),
)
```

### Composable conditions

```go
isUtil := h.Or(
    func(p *h.PackageObj) bool { return p.Name() == "util" },
    func(p *h.PackageObj) bool { return p.Name() == "helpers" },
)

h.Check(t,
    h.Packages().Matching(isUtil).MustNotImport(
        h.Packages().WithPrefix("github.com/myorg/myproject/domain").Select(),
    ),
)
```

### Allowlist variant

```go
h.Check(t,
    // The HTTP handler layer may only import these packages
    h.Packages().WithPrefix("github.com/myorg/myproject/handler").MustOnlyImport(
        h.Packages().Matching(h.Or(
            func(p *h.PackageObj) bool { return p.PkgPath() == "net/http" },
            func(p *h.PackageObj) bool { return strings.HasPrefix(p.PkgPath(), "github.com/myorg/myproject/domain") },
            func(p *h.PackageObj) bool { return strings.HasPrefix(p.PkgPath(), "github.com/myorg/myproject/api") },
        )).Select(),
    ),
)
```

## Design

### Objects

Objects represent analyzable Go code elements: **packages**, **functions/methods**, and **structs**. Each carries its name, package path, and source position.

### Queries

Entry points `Packages()`, `Functions()`, `Structs()` create query builders. Chainable filters narrow the selection:

| Method | Available on | Description |
|---|---|---|
| `WithPrefix(s)` | `PackageQuery` | Package path starts with `s` |
| `WithName(s)` | all | Exact name match |
| `Exported()` | `FuncQuery`, `StructQuery` | Only exported symbols |
| `ResidesIn(sel)` | `FuncQuery`, `StructQuery` | Symbol lives in a matching package |
| `Matching(fn)` | all | Arbitrary predicate |

Query builders are **immutable** — each filter method returns a new builder, so a base query can be safely reused.

### Selectors

Calling `.Select()` on a query produces a lazy `Selector[T]` that resolves on first use and caches per `Context`. Selectors can be shared across multiple rules without redundant work.

### Rules & Assertions

Terminal methods on queries produce `Rule` values:

| Assertion | On | Meaning |
|---|---|---|
| `MustNotImport(sel)` | `PackageQuery`, `StructQuery` | Must not import / reference types from matched packages |
| `MustOnlyImport(sel)` | `PackageQuery` | May only import matched packages (allowlist) |
| `MustNotCall(sel)` | `FuncQuery` | Must not call matched functions |
| `MustNotEmbed(sel)` | `StructQuery` | Must not embed matched structs |

### Laziness

Everything is deferred until `Check()` runs. `Check` creates a `Context`, loads the package graph once, then evaluates all rules against it. Selectors resolve once per context and cache their results.

### Integration

`Check(t *testing.T, rules ...Rule)` is the single entry point. It loads packages matching `HARNESS_PACKAGES` (comma-separated) or `./...` by default, evaluates every rule, and calls `t.Errorf` for each violation.

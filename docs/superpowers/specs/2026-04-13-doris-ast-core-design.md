# Doris AST Core Design

DAG nodes: F1 (ast-core) + C1 (corpus-legacy)

## Goal

Bootstrap `doris/ast/` with the minimal foundation types that the lexer (F2), parser (F4), and all downstream nodes depend on. Also lift the legacy test corpus into `doris/parser/testdata/legacy/`.

## Scope

### In scope

- `Node` interface, `Loc` position struct, `NodeTag` type
- `Visitor` interface and `Walk` function
- `File` root node (batch of statements)
- `ObjectName` node (qualified multi-part identifier)
- Copy 51 legacy SQL test files from `bytebase/parser/doris/examples/`

### Out of scope

- Statement types (T1.1+)
- Expression types (T1.3)
- DDL-specific nodes (T2.x)
- `outfuncs.go` / string representation
- `enums.go` / SQL option enums
- Code generator for walk dispatch

## Reference

Template: `snowflake/ast/` (9 files, 33KB, same minimal-start pattern).

## Package: `doris/ast/`

### `node.go` (~60 LOC)

Core interface and position type:

```go
package ast

// Node is the interface implemented by all AST nodes.
type Node interface {
    Tag() NodeTag
}

// Loc tracks the byte-offset span of a node in the source text.
// -1 means the position is unknown.
type Loc struct {
    Start int
    End   int
}

func (l Loc) LocStart() int { return l.Start }
func (l Loc) LocEnd() int   { return l.End }
```

### `nodetags.go` (~20 LOC)

```go
package ast

type NodeTag uint16

const (
    TagFile       NodeTag = iota + 1
    TagObjectName
)
```

New tags are appended as later DAG nodes add statement/expression types.

### `parsenodes.go` (~40 LOC)

Initial node structs:

```go
package ast

// File is the root node representing a parsed SQL file (batch of statements).
type File struct {
    Loc
    Stmts []Node
}

func (n *File) Tag() NodeTag { return TagFile }

// ObjectName is a qualified multi-part identifier (e.g., catalog.db.table).
type ObjectName struct {
    Loc
    Parts []string
}

func (n *ObjectName) Tag() NodeTag { return TagObjectName }
```

### `walk.go` (~40 LOC)

Visitor pattern:

```go
package ast

// Visitor visits AST nodes. Visit is called for each node encountered by Walk.
// If Visit returns a non-nil Visitor, Walk recurses into the node's children
// with the returned Visitor, then calls Visit(nil) for post-order processing.
type Visitor interface {
    Visit(node Node) Visitor
}

// Walk traverses an AST in depth-first order.
func Walk(v Visitor, node Node) {
    if node == nil {
        return
    }
    w := v.Visit(node)
    if w == nil {
        return
    }
    walkChildren(w, node)
    w.Visit(nil)
}

func walkNodes(v Visitor, nodes []Node) {
    for _, n := range nodes {
        Walk(v, n)
    }
}
```

### `walk_generated.go` (~20 LOC)

Manual dispatch (generated once scale demands a tool):

```go
package ast

func walkChildren(v Visitor, node Node) {
    switch n := node.(type) {
    case *File:
        walkNodes(v, n.Stmts)
    case *ObjectName:
        // leaf node, no children
    }
}
```

### `loc.go` (~30 LOC)

Position extraction helpers (following snowflake/ast/loc.go pattern):

```go
package ast

// NodeLoc returns the Loc for any AST node.
func NodeLoc(node Node) Loc {
    switch n := node.(type) {
    case *File:
        return n.Loc
    case *ObjectName:
        return n.Loc
    default:
        return Loc{Start: -1, End: -1}
    }
}
```

## Test corpus: `doris/parser/testdata/legacy/`

Copy all `.sql` files from `/Users/h3n4l/OpenSource/parser/doris/examples/` (51 files) into `doris/parser/testdata/legacy/`. These are the regression baseline -- the omni parser must eventually parse all of them without error.

No test harness is added in this node (that comes with the parser in F4). The files are placed here so they're available when F3/F4 need them.

## Testing

Since this is a pure type-definition package with minimal logic, testing is lightweight:

- Verify `File` and `ObjectName` implement `Node` interface (compile check)
- Verify `Walk` traverses `File.Stmts` correctly (unit test)
- Verify `NodeLoc` returns correct positions (unit test)
- Run `go build ./doris/...` to confirm compilation
- Run `go vet ./doris/...` for correctness

## File summary

| File | Purpose | ~LOC |
|------|---------|------|
| `doris/ast/node.go` | Node interface, Loc struct | 60 |
| `doris/ast/nodetags.go` | NodeTag type + constants | 20 |
| `doris/ast/parsenodes.go` | File, ObjectName structs | 40 |
| `doris/ast/walk.go` | Visitor interface, Walk function | 40 |
| `doris/ast/walk_generated.go` | walkChildren dispatch | 20 |
| `doris/ast/loc.go` | NodeLoc helper | 30 |
| `doris/parser/testdata/legacy/*.sql` | 51 legacy SQL files | N/A |

Total new Go code: ~210 LOC across 6 files.

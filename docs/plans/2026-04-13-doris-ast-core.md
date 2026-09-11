# Doris AST Core + Corpus Legacy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bootstrap `doris/ast/` with foundation types (Node, Loc, Walk, File, ObjectName) and copy the legacy test corpus into `doris/parser/testdata/legacy/`.

**Architecture:** Mirror `snowflake/ast/` exactly -- same Node interface, Loc struct, Visitor/Walk pattern, walkChildren dispatch. Doris differences: backtick-quoted identifiers (not double-quote), case-preserving normalization (lowercase, not uppercase), `Parts []string` on ObjectName instead of typed `Ident` fields (Doris uses `multipartIdentifier` which is a flat list of parts, not a fixed 3-part structure).

**Tech Stack:** Go 1.25, module `github.com/bytebase/omni`

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core` (branch `feat/doris/ast-core`)

---

## File Structure

| File | Responsibility |
|------|----------------|
| `doris/ast/node.go` | `Node` interface, `Loc` struct, `NoLoc()`, `Loc.IsValid()`, `Loc.Contains()`, `Loc.Merge()` |
| `doris/ast/nodetags.go` | `NodeTag` type, `T_Invalid`, `T_File`, `T_ObjectName` constants, `String()` |
| `doris/ast/parsenodes.go` | `File` and `ObjectName` node structs with `Tag()` methods |
| `doris/ast/walk.go` | `Visitor` interface, `Walk()`, `Inspect()`, `walkNodes()` helper |
| `doris/ast/walk_generated.go` | `walkChildren()` dispatch covering File and ObjectName |
| `doris/ast/loc.go` | `NodeLoc()` dispatcher, `SpanNodes()` helper |
| `doris/ast/ast_test.go` | Unit tests for Walk, NodeLoc, ObjectName, Loc |
| `doris/parser/testdata/legacy/*.sql` | 50 SQL files + `regression/` subdir (2 files) from legacy parser |

---

### Task 1: Create `doris/ast/node.go`

**Files:**
- Create: `doris/ast/node.go`

- [ ] **Step 1: Create the file**

```go
// Package ast defines Doris parse-tree node types.
//
// The package mirrors omni's snowflake/ast and mysql/ast conventions:
// every concrete node implements Node by exposing a Tag() NodeTag method,
// every node carries a Loc field for source-position tracking, and
// walk_generated.go dispatches walker traversal via a type switch.
package ast

// Node is the interface implemented by every Doris parse-tree node.
//
// Tag returns a unique NodeTag identifying the concrete type. Use Tag for
// fast switch dispatch in hot paths; use a Go type assertion when you need
// to access the concrete fields.
type Node interface {
	Tag() NodeTag
}

// Loc represents a source byte range. -1 means "unknown" for either field.
//
// Loc is a value type embedded as a plain field on every concrete node;
// it is NOT part of the Node interface (mirroring snowflake/ast). Helpers
// in loc.go extract Loc from any Node.
type Loc struct {
	Start int // inclusive start byte offset (-1 if unknown)
	End   int // exclusive end byte offset   (-1 if unknown)
}

// NoLoc returns a Loc with both Start and End set to -1 (unknown).
func NoLoc() Loc {
	return Loc{Start: -1, End: -1}
}

// IsValid reports whether both Start and End are non-negative.
func (l Loc) IsValid() bool {
	return l.Start >= 0 && l.End >= 0
}

// Contains reports whether l fully contains other (inclusive on both ends).
// Returns false if either Loc is invalid.
func (l Loc) Contains(other Loc) bool {
	if !l.IsValid() || !other.IsValid() {
		return false
	}
	return l.Start <= other.Start && other.End <= l.End
}

// Merge returns the smallest Loc that contains both l and other.
// If either side is invalid, returns the other side. If both are invalid,
// returns NoLoc().
func (l Loc) Merge(other Loc) Loc {
	if !l.IsValid() {
		return other
	}
	if !other.IsValid() {
		return l
	}
	out := l
	if other.Start < out.Start {
		out.Start = other.Start
	}
	if other.End > out.End {
		out.End = other.End
	}
	return out
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go build ./doris/ast/...`

Expected: No output (success). Will fail because package has only one file and no references yet -- that's fine, we'll verify full compilation after all files are in place.

- [ ] **Step 3: Commit**

```bash
git add doris/ast/node.go
git commit -m "feat(doris/ast): add Node interface and Loc position type"
```

---

### Task 2: Create `doris/ast/nodetags.go`

**Files:**
- Create: `doris/ast/nodetags.go`

- [ ] **Step 1: Create the file**

```go
package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under doris/ast must declare a unique
// NodeTag constant in this file and return it from Tag(). This enables fast
// switch dispatch and code-generated walker support.
//
// The numeric values are NOT stable -- do not persist them. Tags are assigned
// by source order; reorder freely as the package evolves.
type NodeTag int

const (
	// T_Invalid is the zero-value tag, returned only by uninitialized nodes
	// or test stubs that have no need for a real tag.
	T_Invalid NodeTag = iota

	// T_File is the tag for *File, the root statement-list container
	// returned by the parser entry point.
	T_File

	// T_ObjectName is the tag for *ObjectName, a qualified multi-part name
	// (e.g., catalog.db.table).
	T_ObjectName
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_ObjectName:
		return "ObjectName"
	default:
		return "Unknown"
	}
}
```

- [ ] **Step 2: Commit**

```bash
git add doris/ast/nodetags.go
git commit -m "feat(doris/ast): add NodeTag type with File and ObjectName tags"
```

---

### Task 3: Create `doris/ast/parsenodes.go`

**Files:**
- Create: `doris/ast/parsenodes.go`

- [ ] **Step 1: Create the file**

Doris uses backtick-quoted identifiers and `multipartIdentifier` is a flat list of dot-separated parts. ObjectName stores `Parts []string` rather than typed Database/Schema/Name fields -- this matches the Doris grammar's `multipartIdentifier` rule and is simpler than Snowflake's 3-field approach.

```go
package ast

import "strings"

// This file holds the concrete Doris parse-tree node types. F1 ships only
// the File root container and ObjectName; later migration nodes (T1.1+)
// populate the rest.

// File is the root node of a parsed Doris source file. It holds the
// top-level statement list and the byte range covering the entire file.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (n *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// Identifier types
// ---------------------------------------------------------------------------

// ObjectName represents a qualified multi-part identifier as used by Doris's
// multipartIdentifier grammar rule. Parts are stored in order: for a 3-part
// name like catalog.db.table, Parts = ["catalog", "db", "table"].
//
// Parts stores the raw text from source: for backtick-quoted identifiers,
// the content between the backticks; for unquoted identifiers, the source
// bytes with case preserved.
//
// ObjectName is a Node and is visited by the AST walker, but has no Node
// children to descend into.
type ObjectName struct {
	Parts []string
	Loc   Loc
}

// Tag implements Node.
func (n *ObjectName) Tag() NodeTag { return T_ObjectName }

// Compile-time assertion that *ObjectName satisfies Node.
var _ Node = (*ObjectName)(nil)

// String returns the dotted form of the name (e.g., "catalog.db.table").
func (n *ObjectName) String() string {
	return strings.Join(n.Parts, ".")
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go build ./doris/ast/...`

Expected: Will fail because `walkChildren` is not yet defined. That's expected -- next tasks complete the package.

- [ ] **Step 3: Commit**

```bash
git add doris/ast/parsenodes.go
git commit -m "feat(doris/ast): add File and ObjectName node types"
```

---

### Task 4: Create `doris/ast/walk.go` and `doris/ast/walk_generated.go`

**Files:**
- Create: `doris/ast/walk.go`
- Create: `doris/ast/walk_generated.go`

- [ ] **Step 1: Create walk.go**

```go
package ast

// Visitor is implemented by clients of Walk.
//
// Visit is called for each node during a depth-first walk. If Visit returns
// a non-nil w, Walk recurses into the node's children with w, then calls
// w.Visit(nil) to signal end-of-children (post-order). If Visit returns nil,
// children are not visited.
//
// This is the same shape as snowflake/ast.Visitor and pg/ast.Visitor.
type Visitor interface {
	Visit(node Node) Visitor
}

// Walk traverses an AST in depth-first order. It calls v.Visit(node);
// if that returns a non-nil visitor w, it walks each child node with w,
// then calls w.Visit(nil) to mark post-order.
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

// Inspect traverses an AST in depth-first order, calling f for each node.
// If f returns true, Inspect recurses into the node's children.
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if node != nil && f(node) {
		return f
	}
	return nil
}

// walkNodes visits each entry in nodes by calling Walk on it. Used by
// walk_generated.go to traverse slice-typed child fields like []Node.
func walkNodes(v Visitor, nodes []Node) {
	for _, n := range nodes {
		Walk(v, n)
	}
}
```

- [ ] **Step 2: Create walk_generated.go**

```go
package ast

// walkChildren walks the child nodes of node, calling Walk(v, child)
// for each child. Maintained manually until the node count warrants a
// code generator.
func walkChildren(v Visitor, node Node) {
	switch n := node.(type) {
	case *File:
		walkNodes(v, n.Stmts)
	case *ObjectName:
		// leaf node, no children
	}
}
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go build ./doris/ast/...`

Expected: Success (no output). The package now has all required symbols.

- [ ] **Step 4: Commit**

```bash
git add doris/ast/walk.go doris/ast/walk_generated.go
git commit -m "feat(doris/ast): add Visitor/Walk pattern and walkChildren dispatch"
```

---

### Task 5: Create `doris/ast/loc.go`

**Files:**
- Create: `doris/ast/loc.go`

- [ ] **Step 1: Create the file**

```go
package ast

// NodeLoc returns the source location of n, or NoLoc() if n is nil or its
// concrete type carries no Loc field. Every concrete node type added under
// doris/ast must add a case here.
//
// The pattern matches snowflake/ast.NodeLoc.
func NodeLoc(n Node) Loc {
	if n == nil {
		return NoLoc()
	}
	switch v := n.(type) {
	case *File:
		return v.Loc
	case *ObjectName:
		return v.Loc
	default:
		return NoLoc()
	}
}

// SpanNodes returns the smallest Loc that covers every node in nodes.
// Nil entries and nodes whose Loc is invalid are skipped. Returns NoLoc()
// when no node has a valid Loc (including the empty-args case).
func SpanNodes(nodes ...Node) Loc {
	out := NoLoc()
	for _, n := range nodes {
		if n == nil {
			continue
		}
		out = out.Merge(NodeLoc(n))
	}
	return out
}
```

- [ ] **Step 2: Verify compilation and vet**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go build ./doris/ast/... && go vet ./doris/ast/...`

Expected: Success (no output).

- [ ] **Step 3: Commit**

```bash
git add doris/ast/loc.go
git commit -m "feat(doris/ast): add NodeLoc and SpanNodes position helpers"
```

---

### Task 6: Write tests in `doris/ast/ast_test.go`

**Files:**
- Create: `doris/ast/ast_test.go`

- [ ] **Step 1: Create the test file**

```go
package ast

import "testing"

func TestNodeInterface(t *testing.T) {
	// Verify that File and ObjectName satisfy the Node interface.
	var _ Node = (*File)(nil)
	var _ Node = (*ObjectName)(nil)
}

func TestNodeTags(t *testing.T) {
	tests := []struct {
		node Node
		want NodeTag
	}{
		{&File{}, T_File},
		{&ObjectName{Parts: []string{"t"}}, T_ObjectName},
	}
	for _, tt := range tests {
		if got := tt.node.Tag(); got != tt.want {
			t.Errorf("Tag() = %v, want %v", got, tt.want)
		}
	}
}

func TestNodeTagString(t *testing.T) {
	tests := []struct {
		tag  NodeTag
		want string
	}{
		{T_Invalid, "Invalid"},
		{T_File, "File"},
		{T_ObjectName, "ObjectName"},
		{NodeTag(9999), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.tag.String(); got != tt.want {
			t.Errorf("NodeTag(%d).String() = %q, want %q", tt.tag, got, tt.want)
		}
	}
}

func TestLocIsValid(t *testing.T) {
	tests := []struct {
		loc  Loc
		want bool
	}{
		{Loc{0, 10}, true},
		{Loc{-1, 10}, false},
		{Loc{0, -1}, false},
		{NoLoc(), false},
	}
	for _, tt := range tests {
		if got := tt.loc.IsValid(); got != tt.want {
			t.Errorf("Loc{%d,%d}.IsValid() = %v, want %v", tt.loc.Start, tt.loc.End, got, tt.want)
		}
	}
}

func TestLocContains(t *testing.T) {
	outer := Loc{0, 100}
	inner := Loc{10, 50}
	disjoint := Loc{200, 300}

	if !outer.Contains(inner) {
		t.Error("outer should contain inner")
	}
	if inner.Contains(outer) {
		t.Error("inner should not contain outer")
	}
	if outer.Contains(disjoint) {
		t.Error("outer should not contain disjoint")
	}
	if NoLoc().Contains(inner) {
		t.Error("NoLoc should not contain anything")
	}
}

func TestLocMerge(t *testing.T) {
	a := Loc{10, 20}
	b := Loc{5, 30}
	merged := a.Merge(b)
	if merged.Start != 5 || merged.End != 30 {
		t.Errorf("Merge = %v, want {5, 30}", merged)
	}

	// Merge with NoLoc returns the valid side.
	if got := a.Merge(NoLoc()); got != a {
		t.Errorf("Merge(NoLoc) = %v, want %v", got, a)
	}
	if got := NoLoc().Merge(b); got != b {
		t.Errorf("NoLoc.Merge = %v, want %v", got, b)
	}
}

func TestObjectNameString(t *testing.T) {
	tests := []struct {
		parts []string
		want  string
	}{
		{[]string{"t"}, "t"},
		{[]string{"db", "t"}, "db.t"},
		{[]string{"catalog", "db", "t"}, "catalog.db.t"},
	}
	for _, tt := range tests {
		n := &ObjectName{Parts: tt.parts}
		if got := n.String(); got != tt.want {
			t.Errorf("ObjectName(%v).String() = %q, want %q", tt.parts, got, tt.want)
		}
	}
}

func TestWalkFile(t *testing.T) {
	// Build a small tree: File with two ObjectName children.
	child1 := &ObjectName{Parts: []string{"a"}, Loc: Loc{0, 1}}
	child2 := &ObjectName{Parts: []string{"b"}, Loc: Loc{2, 3}}
	file := &File{Stmts: []Node{child1, child2}, Loc: Loc{0, 3}}

	var visited []NodeTag
	Inspect(file, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	want := []NodeTag{T_File, T_ObjectName, T_ObjectName}
	if len(visited) != len(want) {
		t.Fatalf("visited %d nodes, want %d", len(visited), len(want))
	}
	for i, tag := range visited {
		if tag != want[i] {
			t.Errorf("visited[%d] = %v, want %v", i, tag, want[i])
		}
	}
}

func TestWalkNilNode(t *testing.T) {
	// Walk(nil) should not panic.
	Walk(inspector(func(Node) bool { return true }), nil)
}

func TestNodeLoc(t *testing.T) {
	file := &File{Loc: Loc{0, 100}}
	obj := &ObjectName{Parts: []string{"t"}, Loc: Loc{5, 10}}

	if got := NodeLoc(file); got != file.Loc {
		t.Errorf("NodeLoc(File) = %v, want %v", got, file.Loc)
	}
	if got := NodeLoc(obj); got != obj.Loc {
		t.Errorf("NodeLoc(ObjectName) = %v, want %v", got, obj.Loc)
	}
	if got := NodeLoc(nil); got != NoLoc() {
		t.Errorf("NodeLoc(nil) = %v, want NoLoc", got)
	}
}

func TestSpanNodes(t *testing.T) {
	a := &ObjectName{Parts: []string{"a"}, Loc: Loc{10, 20}}
	b := &ObjectName{Parts: []string{"b"}, Loc: Loc{5, 30}}

	span := SpanNodes(a, b)
	if span.Start != 5 || span.End != 30 {
		t.Errorf("SpanNodes = %v, want {5, 30}", span)
	}

	// With nil entries.
	span2 := SpanNodes(nil, a, nil)
	if span2 != a.Loc {
		t.Errorf("SpanNodes(nil, a, nil) = %v, want %v", span2, a.Loc)
	}

	// Empty args.
	if got := SpanNodes(); got != NoLoc() {
		t.Errorf("SpanNodes() = %v, want NoLoc", got)
	}
}
```

- [ ] **Step 2: Run tests**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go test ./doris/ast/... -v`

Expected: All tests PASS.

- [ ] **Step 3: Commit**

```bash
git add doris/ast/ast_test.go
git commit -m "test(doris/ast): add unit tests for Node, Loc, Walk, ObjectName"
```

---

### Task 7: Copy legacy test corpus (C1)

**Files:**
- Create: `doris/parser/testdata/legacy/*.sql` (50 files + `regression/` subdir with 2 files)

- [ ] **Step 1: Create directory structure and copy files**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core
mkdir -p doris/parser/testdata/legacy/regression
cp /Users/h3n4l/OpenSource/parser/doris/examples/*.sql doris/parser/testdata/legacy/
cp /Users/h3n4l/OpenSource/parser/doris/examples/regression/*.sql doris/parser/testdata/legacy/regression/
```

- [ ] **Step 2: Verify file count**

Run: `find doris/parser/testdata/legacy -name '*.sql' | wc -l`

Expected: `52` (50 top-level + 2 in regression/)

- [ ] **Step 3: Commit**

```bash
git add doris/parser/testdata/legacy/
git commit -m "test(doris): lift legacy ANTLR4 test corpus (52 SQL files)"
```

---

### Task 8: Final verification

- [ ] **Step 1: Run full build and vet**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go build ./doris/... && go vet ./doris/...`

Expected: Success (no output).

- [ ] **Step 2: Run all doris tests**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go test ./doris/... -v`

Expected: All tests PASS.

- [ ] **Step 3: Verify no regressions in existing packages**

Run: `cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-doris-ast-core && go build ./...`

Expected: Success (no output). The new `doris/` package should not affect any existing engine.

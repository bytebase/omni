# CosmosDB Query Analysis Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `cosmosdb/analysis/` package that extracts field-level projection and predicate information from parsed CosmosDB SELECT queries.

**Architecture:** Walk the raw parsed AST (`*ast.SelectStmt`) to build an alias resolution map from FROM/JOIN, then extract field paths from SELECT projections and WHERE predicates. Returns a self-contained `QueryAnalysis` struct.

**Tech Stack:** Pure Go, zero dependencies beyond `cosmosdb/ast`.

**Spec:** `docs/superpowers/specs/2026-03-30-cosmosdb-analysis-design.md`

---

## File Map

| File | Responsibility |
|------|---------------|
| `cosmosdb/analysis/fieldpath.go` | `Selector` and `FieldPath` types |
| `cosmosdb/analysis/analysis.go` | `Projection`, `QueryAnalysis` types, `Analyze()` entry point, alias map building |
| `cosmosdb/analysis/extract.go` | `extractFieldPaths()` recursive AST walker |
| `cosmosdb/analysis/analysis_test.go` | Table-driven tests: SQL -> expected QueryAnalysis |

---

### Task 1: FieldPath Types

**Files:**
- Create: `cosmosdb/analysis/fieldpath.go`

- [ ] **Step 1: Create `cosmosdb/analysis/fieldpath.go`**

```go
// Package analysis extracts field-level information from parsed CosmosDB queries.
package analysis

import "strings"

// Selector represents one step in a document property path.
type Selector struct {
	Name       string // property name
	ArrayIndex int    // -1 for item access (.name), >= 0 for array index ([n])
}

// ItemSelector creates a Selector for property access.
func ItemSelector(name string) Selector {
	return Selector{Name: name, ArrayIndex: -1}
}

// ArraySelector creates a Selector for array index access.
func ArraySelector(name string, index int) Selector {
	return Selector{Name: name, ArrayIndex: index}
}

// IsArray returns true if this selector represents an array index access.
func (s Selector) IsArray() bool {
	return s.ArrayIndex >= 0
}

// FieldPath represents a path through a JSON document.
// e.g., container.addresses[1].country is:
//
//	FieldPath{ItemSelector("container"), ArraySelector("addresses", 1), ItemSelector("country")}
type FieldPath []Selector

// String returns a human-readable representation like "container.addresses[1].country".
func (fp FieldPath) String() string {
	var sb strings.Builder
	for i, s := range fp {
		if i > 0 && !s.IsArray() {
			sb.WriteByte('.')
		}
		if s.IsArray() {
			sb.WriteString("[" + itoa(s.ArrayIndex) + "]")
		} else {
			sb.WriteString(s.Name)
		}
	}
	return sb.String()
}

func itoa(n int) string {
	// Avoid importing strconv for a single use.
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/analysis/`
Expected: Clean build.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/analysis/fieldpath.go
git commit -m "feat(cosmosdb): add FieldPath and Selector types for query analysis"
```

---

### Task 2: QueryAnalysis Types and Analyze Entry Point

**Files:**
- Create: `cosmosdb/analysis/analysis.go`

- [ ] **Step 1: Create `cosmosdb/analysis/analysis.go`**

```go
package analysis

import (
	"strconv"

	"github.com/bytebase/omni/cosmosdb/ast"
)

// Projection represents one item in the SELECT list with its output name
// and the source field paths it references.
type Projection struct {
	Name        string      // output alias; empty if none and not inferrable
	SourcePaths []FieldPath // source field paths this projection references
}

// QueryAnalysis is the result of analyzing a parsed CosmosDB SELECT statement.
type QueryAnalysis struct {
	Projections []Projection // SELECT list items (empty when SelectStar is true)
	SelectStar  bool         // true if SELECT *
	Predicates  []FieldPath  // field paths referenced in WHERE
}

// Analyze extracts field-level information from a parsed CosmosDB SELECT statement.
func Analyze(stmt *ast.SelectStmt) *QueryAnalysis {
	aliases := buildAliasMap(stmt)
	qa := &QueryAnalysis{}

	// SELECT * detection.
	if stmt.Star {
		qa.SelectStar = true
	} else {
		qa.Projections = extractProjections(stmt.Targets, aliases)
	}

	// WHERE predicates.
	if stmt.Where != nil {
		qa.Predicates = extractFieldPaths(stmt.Where, aliases)
	}

	return qa
}

// aliasMap maps alias names to their resolved source field paths.
// e.g., "p" -> [ItemSelector("products")], "t" -> [ItemSelector("products"), ItemSelector("tags")]
type aliasMap map[string]FieldPath

// buildAliasMap walks FROM and JOINs to build the alias resolution map.
func buildAliasMap(stmt *ast.SelectStmt) aliasMap {
	m := make(aliasMap)
	if stmt.From == nil {
		return m
	}

	// Walk the primary FROM source.
	registerFromSource(stmt.From, m)

	// Walk JOINs.
	for _, join := range stmt.Joins {
		registerFromSource(join.Source, m)
	}

	return m
}

// registerFromSource registers aliases from a single FROM/JOIN source.
func registerFromSource(te ast.TableExpr, m aliasMap) {
	switch src := te.(type) {
	case *ast.AliasedTableExpr:
		containerName := containerRefName(src.Source)
		if containerName != "" && src.Alias != "" {
			m[src.Alias] = FieldPath{ItemSelector(containerName)}
		}
	case *ast.ContainerRef:
		if src.Name != "" {
			m[src.Name] = FieldPath{ItemSelector(src.Name)}
		}
	case *ast.ArrayIterationExpr:
		// e.g., "t IN p.tags" — resolve p.tags through the alias map,
		// then register "t" as pointing to that resolved path.
		sourcePath := resolveTableExprPath(src.Source, m)
		if src.Alias != "" && len(sourcePath) > 0 {
			m[src.Alias] = sourcePath
		}
	}
}

// containerRefName extracts the container name from a TableExpr.
func containerRefName(te ast.TableExpr) string {
	switch src := te.(type) {
	case *ast.ContainerRef:
		return src.Name
	case *ast.DotAccessExpr:
		// e.g., products.sizes — the root is the container name.
		return exprRootName(src)
	default:
		return ""
	}
}

// resolveTableExprPath resolves a TableExpr to a FieldPath using the alias map.
// e.g., DotAccessExpr{ColumnRef("p"), "tags"} with alias p->products
// resolves to [ItemSelector("products"), ItemSelector("tags")].
func resolveTableExprPath(te ast.TableExpr, m aliasMap) FieldPath {
	switch src := te.(type) {
	case *ast.ContainerRef:
		if resolved, ok := m[src.Name]; ok {
			return copyPath(resolved)
		}
		return FieldPath{ItemSelector(src.Name)}
	case *ast.DotAccessExpr:
		base := resolveExprPath(src.Expr, m)
		return append(base, ItemSelector(src.Property))
	case *ast.BracketAccessExpr:
		base := resolveExprPath(src.Expr, m)
		return append(base, bracketSelector(src.Index))
	default:
		return nil
	}
}

// resolveExprPath resolves an ExprNode to a FieldPath using the alias map.
func resolveExprPath(expr ast.ExprNode, m aliasMap) FieldPath {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if resolved, ok := m[e.Name]; ok {
			return copyPath(resolved)
		}
		return FieldPath{ItemSelector(e.Name)}
	case *ast.DotAccessExpr:
		base := resolveExprPath(e.Expr, m)
		return append(base, ItemSelector(e.Property))
	case *ast.BracketAccessExpr:
		base := resolveExprPath(e.Expr, m)
		return append(base, bracketSelector(e.Index))
	default:
		return nil
	}
}

// exprRootName extracts the root identifier name from a chain of dot/bracket accesses.
func exprRootName(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		return e.Name
	case *ast.DotAccessExpr:
		return exprRootName(e.Expr)
	case *ast.BracketAccessExpr:
		return exprRootName(e.Expr)
	default:
		return ""
	}
}

// bracketSelector creates a Selector from a bracket index expression.
func bracketSelector(index ast.ExprNode) Selector {
	switch idx := index.(type) {
	case *ast.NumberLit:
		if n, err := strconv.Atoi(idx.Val); err == nil {
			return ArraySelector(idx.Val, n)
		}
		return ItemSelector(idx.Val)
	case *ast.StringLit:
		return ItemSelector(idx.Val)
	default:
		return ItemSelector("?")
	}
}

// extractProjections builds the projection list from SELECT targets.
func extractProjections(targets []*ast.TargetEntry, aliases aliasMap) []Projection {
	var projections []Projection
	for _, t := range targets {
		paths := extractFieldPaths(t.Expr, aliases)
		name := ""
		if t.Alias != nil {
			name = *t.Alias
		} else if len(paths) == 1 && len(paths[0]) > 0 {
			// Infer name from the last selector in the path.
			last := paths[0][len(paths[0])-1]
			if !last.IsArray() {
				name = last.Name
			}
		}
		projections = append(projections, Projection{
			Name:        name,
			SourcePaths: paths,
		})
	}
	return projections
}

// copyPath returns a shallow copy of a FieldPath to avoid alias mutation.
func copyPath(fp FieldPath) FieldPath {
	cp := make(FieldPath, len(fp))
	copy(cp, fp)
	return cp
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/analysis/`
Expected: Will fail because `extractFieldPaths` is not yet defined. Verify syntax only:
Run: `gofmt -e cosmosdb/analysis/analysis.go > /dev/null`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/analysis/analysis.go
git commit -m "feat(cosmosdb): add QueryAnalysis types and Analyze entry point"
```

---

### Task 3: Field Path Extraction Walker

**Files:**
- Create: `cosmosdb/analysis/extract.go`

- [ ] **Step 1: Create `cosmosdb/analysis/extract.go`**

```go
package analysis

import (
	"github.com/bytebase/omni/cosmosdb/ast"
)

// extractFieldPaths recursively collects all field paths referenced in an expression.
func extractFieldPaths(expr ast.ExprNode, aliases aliasMap) []FieldPath {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *ast.ColumnRef:
		if resolved, ok := aliases[e.Name]; ok {
			return []FieldPath{copyPath(resolved)}
		}
		return []FieldPath{{ItemSelector(e.Name)}}

	case *ast.DotAccessExpr:
		bases := extractFieldPaths(e.Expr, aliases)
		for i := range bases {
			bases[i] = append(bases[i], ItemSelector(e.Property))
		}
		return bases

	case *ast.BracketAccessExpr:
		bases := extractFieldPaths(e.Expr, aliases)
		sel := bracketSelector(e.Index)
		for i := range bases {
			bases[i] = append(bases[i], sel)
		}
		return bases

	case *ast.BinaryExpr:
		return collectFromExprs(aliases, e.Left, e.Right)

	case *ast.UnaryExpr:
		return extractFieldPaths(e.Operand, aliases)

	case *ast.TernaryExpr:
		return collectFromExprs(aliases, e.Cond, e.Then, e.Else)

	case *ast.InExpr:
		paths := extractFieldPaths(e.Expr, aliases)
		for _, item := range e.List {
			paths = append(paths, extractFieldPaths(item, aliases)...)
		}
		return paths

	case *ast.BetweenExpr:
		return collectFromExprs(aliases, e.Expr, e.Low, e.High)

	case *ast.LikeExpr:
		paths := collectFromExprs(aliases, e.Expr, e.Pattern)
		if e.Escape != nil {
			paths = append(paths, extractFieldPaths(e.Escape, aliases)...)
		}
		return paths

	case *ast.FuncCall:
		return collectFromExprSlice(e.Args, aliases)

	case *ast.UDFCall:
		return collectFromExprSlice(e.Args, aliases)

	case *ast.CreateArrayExpr:
		return collectFromExprSlice(e.Elements, aliases)

	case *ast.CreateObjectExpr:
		var paths []FieldPath
		for _, f := range e.Fields {
			paths = append(paths, extractFieldPaths(f.Value, aliases)...)
		}
		return paths

	// Literals and parameters have no field paths.
	case *ast.StringLit, *ast.NumberLit, *ast.BoolLit,
		*ast.NullLit, *ast.UndefinedLit, *ast.InfinityLit,
		*ast.NanLit, *ast.ParamRef:
		return nil

	// Subqueries are not traversed (matches current Bytebase behavior).
	case *ast.SubLink, *ast.ExistsExpr, *ast.ArrayExpr:
		return nil

	default:
		return nil
	}
}

// collectFromExprs collects field paths from multiple expressions.
func collectFromExprs(aliases aliasMap, exprs ...ast.ExprNode) []FieldPath {
	var all []FieldPath
	for _, e := range exprs {
		all = append(all, extractFieldPaths(e, aliases)...)
	}
	return all
}

// collectFromExprSlice collects field paths from a slice of expressions.
func collectFromExprSlice(exprs []ast.ExprNode, aliases aliasMap) []FieldPath {
	var all []FieldPath
	for _, e := range exprs {
		all = append(all, extractFieldPaths(e, aliases)...)
	}
	return all
}
```

- [ ] **Step 2: Verify the full package compiles**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./cosmosdb/analysis/`
Expected: Clean build.

- [ ] **Step 3: Commit**

```bash
git add cosmosdb/analysis/extract.go
git commit -m "feat(cosmosdb): add field path extraction walker"
```

---

### Task 4: Tests

**Files:**
- Create: `cosmosdb/analysis/analysis_test.go`

- [ ] **Step 1: Create `cosmosdb/analysis/analysis_test.go`**

```go
package analysis_test

import (
	"testing"

	"github.com/bytebase/omni/cosmosdb"
	"github.com/bytebase/omni/cosmosdb/analysis"
	"github.com/bytebase/omni/cosmosdb/ast"
)

func parseSelect(t *testing.T, sql string) *ast.SelectStmt {
	t.Helper()
	stmts, err := cosmosdb.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, len(stmts))
	}
	sel, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): AST is %T, want *ast.SelectStmt", sql, stmts[0].AST)
	}
	return sel
}

func pathStr(fp analysis.FieldPath) string {
	return fp.String()
}

func TestAnalyzeSelectStar(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, "SELECT * FROM c"))
	if !qa.SelectStar {
		t.Error("expected SelectStar = true")
	}
	if len(qa.Projections) != 0 {
		t.Errorf("expected no projections, got %d", len(qa.Projections))
	}
}

func TestAnalyzeSimpleProjection(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, "SELECT c.name FROM c"))
	if qa.SelectStar {
		t.Error("expected SelectStar = false")
	}
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	p := qa.Projections[0]
	if p.Name != "name" {
		t.Errorf("projection name = %q, want %q", p.Name, "name")
	}
	if len(p.SourcePaths) != 1 {
		t.Fatalf("expected 1 source path, got %d", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.name" {
		t.Errorf("source path = %q, want %q", got, "c.name")
	}
}

func TestAnalyzeAliasedProjection(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, "SELECT c.name AS n FROM c"))
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	if qa.Projections[0].Name != "n" {
		t.Errorf("projection name = %q, want %q", qa.Projections[0].Name, "n")
	}
}

func TestAnalyzeAliasResolution(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, "SELECT t.name FROM container t"))
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("expected 1 source path, got %d", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "container.name" {
		t.Errorf("source path = %q, want %q", got, "container.name")
	}
}

func TestAnalyzeDotChain(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, "SELECT c.address.city FROM c"))
	p := qa.Projections[0]
	if got := pathStr(p.SourcePaths[0]); got != "c.address.city" {
		t.Errorf("source path = %q, want %q", got, "c.address.city")
	}
	if p.Name != "city" {
		t.Errorf("projection name = %q, want %q", p.Name, "city")
	}
}

func TestAnalyzeBracketAccess(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT c.addresses[1].country FROM c`))
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("expected 1 source path, got %d", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.addresses[1].country" {
		t.Errorf("source path = %q, want %q", got, "c.addresses[1].country")
	}
}

func TestAnalyzeBracketStringAccess(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT c["field"] FROM c`))
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("expected 1 source path, got %d", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.field" {
		t.Errorf("source path = %q, want %q", got, "c.field")
	}
}

func TestAnalyzeExpressionMultiplePaths(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT c.name ?? c.nickname AS displayName FROM c`))
	p := qa.Projections[0]
	if p.Name != "displayName" {
		t.Errorf("projection name = %q, want %q", p.Name, "displayName")
	}
	if len(p.SourcePaths) != 2 {
		t.Fatalf("expected 2 source paths, got %d", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.name" {
		t.Errorf("source path 0 = %q, want %q", got, "c.name")
	}
	if got := pathStr(p.SourcePaths[1]); got != "c.nickname" {
		t.Errorf("source path 1 = %q, want %q", got, "c.nickname")
	}
}

func TestAnalyzeFunctionCall(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT CONCAT(c.first, c.last) AS fullName FROM c`))
	p := qa.Projections[0]
	if p.Name != "fullName" {
		t.Errorf("projection name = %q, want %q", p.Name, "fullName")
	}
	if len(p.SourcePaths) != 2 {
		t.Fatalf("expected 2 source paths, got %d", len(p.SourcePaths))
	}
}

func TestAnalyzeUDFCall(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT udf.myFunc(c.name) AS result FROM c`))
	p := qa.Projections[0]
	if len(p.SourcePaths) != 1 {
		t.Fatalf("expected 1 source path, got %d", len(p.SourcePaths))
	}
	if got := pathStr(p.SourcePaths[0]); got != "c.name" {
		t.Errorf("source path = %q, want %q", got, "c.name")
	}
}

func TestAnalyzeWherePredicate(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT * FROM c WHERE c.country = "US"`))
	if len(qa.Predicates) != 1 {
		t.Fatalf("expected 1 predicate, got %d", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.country" {
		t.Errorf("predicate path = %q, want %q", got, "c.country")
	}
}

func TestAnalyzeWhereMultiplePredicates(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT * FROM c WHERE c.age > 18 AND c.active = true`))
	if len(qa.Predicates) != 2 {
		t.Fatalf("expected 2 predicates, got %d", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.age" {
		t.Errorf("predicate 0 = %q, want %q", got, "c.age")
	}
	if got := pathStr(qa.Predicates[1]); got != "c.active" {
		t.Errorf("predicate 1 = %q, want %q", got, "c.active")
	}
}

func TestAnalyzeWhereBetween(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT * FROM c WHERE c.population BETWEEN 100000 AND 5000000`))
	if len(qa.Predicates) != 1 {
		t.Fatalf("expected 1 predicate, got %d", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.population" {
		t.Errorf("predicate = %q, want %q", got, "c.population")
	}
}

func TestAnalyzeWhereIn(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT * FROM c WHERE c.country IN ("US", "UK")`))
	if len(qa.Predicates) != 1 {
		t.Fatalf("expected 1 predicate, got %d", len(qa.Predicates))
	}
	if got := pathStr(qa.Predicates[0]); got != "c.country" {
		t.Errorf("predicate = %q, want %q", got, "c.country")
	}
}

func TestAnalyzeJoin(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT p.name, t.tag FROM products p JOIN t IN p.tags`))
	if len(qa.Projections) != 2 {
		t.Fatalf("expected 2 projections, got %d", len(qa.Projections))
	}
	// p.name -> products.name
	if got := pathStr(qa.Projections[0].SourcePaths[0]); got != "products.name" {
		t.Errorf("projection 0 path = %q, want %q", got, "products.name")
	}
	// t.tag -> products.tags.tag
	if got := pathStr(qa.Projections[1].SourcePaths[0]); got != "products.tags.tag" {
		t.Errorf("projection 1 path = %q, want %q", got, "products.tags.tag")
	}
}

func TestAnalyzeJoinWherePredicate(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT p.name FROM products p JOIN t IN p.tags WHERE t.active = true`))
	if len(qa.Predicates) != 1 {
		t.Fatalf("expected 1 predicate, got %d", len(qa.Predicates))
	}
	// t.active -> products.tags.active
	if got := pathStr(qa.Predicates[0]); got != "products.tags.active" {
		t.Errorf("predicate = %q, want %q", got, "products.tags.active")
	}
}

func TestAnalyzeNestedJoin(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT s.size FROM products p JOIN t IN p.tags JOIN s IN t.sizes`))
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	// s.size -> products.tags.sizes.size
	if got := pathStr(qa.Projections[0].SourcePaths[0]); got != "products.tags.sizes.size" {
		t.Errorf("path = %q, want %q", got, "products.tags.sizes.size")
	}
}

func TestAnalyzeSelectValue(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT VALUE c.name FROM c`))
	if qa.SelectStar {
		t.Error("expected SelectStar = false")
	}
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	if got := pathStr(qa.Projections[0].SourcePaths[0]); got != "c.name" {
		t.Errorf("path = %q, want %q", got, "c.name")
	}
}

func TestAnalyzeNoFrom(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT 1`))
	if qa.SelectStar {
		t.Error("expected SelectStar = false")
	}
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	if len(qa.Projections[0].SourcePaths) != 0 {
		t.Errorf("expected 0 source paths, got %d", len(qa.Projections[0].SourcePaths))
	}
}

func TestAnalyzeValueCount(t *testing.T) {
	qa := analysis.Analyze(parseSelect(t, `SELECT VALUE COUNT(1) FROM c`))
	if len(qa.Projections) != 1 {
		t.Fatalf("expected 1 projection, got %d", len(qa.Projections))
	}
	// COUNT(1) has no field paths — 1 is a literal.
	if len(qa.Projections[0].SourcePaths) != 0 {
		t.Errorf("expected 0 source paths, got %d", len(qa.Projections[0].SourcePaths))
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./cosmosdb/analysis/ -v`
Expected: All tests pass.

- [ ] **Step 3: Fix any failures, then run again**

Debug and fix until all tests pass.

- [ ] **Step 4: Commit**

```bash
git add cosmosdb/analysis/analysis_test.go
git commit -m "test(cosmosdb): add query analysis tests"
```

---

### Task 5: Integration Verification

- [ ] **Step 1: Run all CosmosDB tests**

Run: `cd /Users/h3n4l/OpenSource/omni && go test ./cosmosdb/... -v`
Expected: All tests pass (parser + analysis).

- [ ] **Step 2: Run go vet**

Run: `cd /Users/h3n4l/OpenSource/omni && go vet ./cosmosdb/...`
Expected: No issues.

- [ ] **Step 3: Verify no impact on existing engines**

Run: `cd /Users/h3n4l/OpenSource/omni && go build ./...`
Expected: Clean build.

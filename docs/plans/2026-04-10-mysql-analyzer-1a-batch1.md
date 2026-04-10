# MySQL Analyzer Phase 1a Batch 1 — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the foundation of `AnalyzeSelectStmt` — making 5 basic single-table SELECT scenarios produce correct Query IR, verified against golden expected output.

**Architecture:** `AnalyzeSelectStmt` in `mysql/catalog` takes an `*ast.SelectStmt` (from the parser) and the `*Catalog` (DDL state) and returns a `*Query` (analyzed IR). It works in three steps: (1) build scope from FROM clause, (2) resolve target list expressions against scope, (3) assemble the Query. Phase 1a Batch 1 only handles single-table FROM (no JOINs) and basic expressions (column refs, literals, star expansion).

**Tech Stack:** Go, `mysql/ast` (parser AST), `mysql/catalog` (catalog state + IR types), `testcontainers-go` (MySQL 8.0 oracle)

---

## Reference files

- Parser AST types: `mysql/ast/parsenodes.go`
  - `SelectStmt` (line 8): `TargetList []ExprNode`, `From []TableExpr`, `Where ExprNode`, etc.
  - `ResTarget` (line 962): `Name string` (alias), `Val ExprNode` (expression) — only wraps when AS alias present
  - `StarExpr` (line 1777): bare `SELECT *`
  - `ColumnRef` (line 761): `Table string`, `Schema string`, `Column string`, `Star bool` — `Star=true` for `t.*`
  - `IntLit` (line 814): `Value int64`
  - `StringLit` (line 832): `Value string`
  - `TableRef` (line 773): `Schema string` (db), `Name string`, `Alias string`
- Catalog types: `mysql/catalog/catalog.go`, `table.go`, `database.go`
  - `Catalog.GetDatabase(name) *Database`, `Database.GetTable(name) *Table` (case-insensitive)
  - `Table.Columns []*Column`, `Table.colByName map[string]int`
  - `Column.Name`, `Column.Position` (0-based in storage, but we'll use 1-based AttNum in IR)
- Query IR types: `mysql/catalog/query.go`
  - `Query`, `TargetEntryQ`, `RangeTableEntryQ`, `JoinTreeQ`, `RangeTableRefQ`, `VarExprQ`, `ConstExprQ`
- Existing test helpers: `mysql/catalog/wt_helpers_test.go` (`wtSetup`, `wtExec`)
- Parser entry: `mysql/parser/parser.go` — `Parse(sql string) (*nodes.List, error)`

## Key parser behavior (verified empirically)

```
SELECT *         → TargetList[0] = *ast.StarExpr{}
SELECT e.*       → TargetList[0] = *ast.ColumnRef{Table:"e", Star:true}
SELECT 1         → TargetList[0] = *ast.IntLit{Value:1}        (NO ResTarget wrapper)
SELECT name      → TargetList[0] = *ast.ColumnRef{Column:"name"} (NO ResTarget wrapper)
SELECT name AS n → TargetList[0] = *ast.ResTarget{Name:"n", Val:*ColumnRef{Column:"name"}}
```

**ResTarget only appears when there's an explicit AS alias.** All other expressions appear bare in TargetList.

---

## Task 1: Analyzer entry point + scope skeleton

**Files:**
- Create: `mysql/catalog/analyze.go`
- Create: `mysql/catalog/scope.go`
- Test: `mysql/catalog/analyze_test.go`

### Step 1: Write the failing test for scenario 1.1 (SELECT 1)

```go
// mysql/catalog/analyze_test.go
package catalog

import (
	"testing"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// parseSelect parses a single SELECT statement and returns the SelectStmt.
func parseSelect(t *testing.T, sql string) *ast.SelectStmt {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(list.Items))
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt, got %T", list.Items[0])
	}
	return sel
}

func TestAnalyze_1_1_BareLiteral(t *testing.T) {
	c := wtSetup(t)
	sel := parseSelect(t, "SELECT 1")

	q, err := c.AnalyzeSelectStmt(sel)
	if err != nil {
		t.Fatalf("AnalyzeSelectStmt error: %v", err)
	}

	// Query shape
	if q.CommandType != CmdSelect {
		t.Errorf("CommandType = %d, want CmdSelect", q.CommandType)
	}
	if len(q.RangeTable) != 0 {
		t.Errorf("RangeTable len = %d, want 0", len(q.RangeTable))
	}
	if q.JoinTree == nil {
		t.Fatal("JoinTree is nil")
	}
	if len(q.JoinTree.FromList) != 0 {
		t.Errorf("FromList len = %d, want 0", len(q.JoinTree.FromList))
	}
	if q.JoinTree.Quals != nil {
		t.Errorf("Quals is non-nil, want nil")
	}

	// TargetList
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList len = %d, want 1", len(q.TargetList))
	}
	te := q.TargetList[0]
	if te.ResNo != 1 {
		t.Errorf("ResNo = %d, want 1", te.ResNo)
	}
	if te.ResName != "1" {
		t.Errorf("ResName = %q, want %q", te.ResName, "1")
	}
	if te.ResJunk {
		t.Error("ResJunk = true, want false")
	}
	constExpr, ok := te.Expr.(*ConstExprQ)
	if !ok {
		t.Fatalf("Expr type = %T, want *ConstExprQ", te.Expr)
	}
	if constExpr.Value != "1" {
		t.Errorf("ConstExprQ.Value = %q, want %q", constExpr.Value, "1")
	}
}
```

### Step 2: Run test to verify it fails

Run: `go test ./mysql/catalog/ -run TestAnalyze_1_1 -v -count=1`
Expected: FAIL — `c.AnalyzeSelectStmt` doesn't exist.

### Step 3: Write minimal scope.go and analyze.go

```go
// mysql/catalog/scope.go
package catalog

import "strings"

// analyzerScope maps table alias/name → RangeTable index + column metadata
// for the current FROM clause. Used by the analyzer to resolve column references.
type analyzerScope struct {
	entries []scopeEntry
	byName  map[string]int // lowered effective name → index in entries
}

type scopeEntry struct {
	name    string // effective name (alias if present, else table name)
	rteIdx  int    // index into Query.RangeTable
	columns []*Column
}

func newScope() *analyzerScope {
	return &analyzerScope{
		byName: make(map[string]int),
	}
}

func (s *analyzerScope) add(name string, rteIdx int, columns []*Column) {
	key := strings.ToLower(name)
	idx := len(s.entries)
	s.entries = append(s.entries, scopeEntry{name: name, rteIdx: rteIdx, columns: columns})
	s.byName[key] = idx
}

// resolveColumn resolves an unqualified column name against the scope.
// Returns (rteIdx, attNum 1-based, error).
func (s *analyzerScope) resolveColumn(colName string) (int, int, error) {
	lower := strings.ToLower(colName)
	found := -1
	foundRTE := -1
	foundAtt := -1
	for i, entry := range s.entries {
		for j, col := range entry.columns {
			if strings.ToLower(col.Name) == lower {
				if found >= 0 {
					return 0, 0, &Error{Code: 1052, Message: "Column '" + colName + "' in field list is ambiguous"}
				}
				found = i
				foundRTE = entry.rteIdx
				foundAtt = j + 1 // 1-based
			}
		}
	}
	if found < 0 {
		return 0, 0, &Error{Code: 1054, Message: "Unknown column '" + colName + "' in 'field list'"}
	}
	return foundRTE, foundAtt, nil
}

// resolveQualifiedColumn resolves a table-qualified column (table.column).
// Returns (rteIdx, attNum 1-based, error).
func (s *analyzerScope) resolveQualifiedColumn(tableName, colName string) (int, int, error) {
	key := strings.ToLower(tableName)
	idx, ok := s.byName[key]
	if !ok {
		return 0, 0, &Error{Code: 1054, Message: "Unknown column '" + tableName + "." + colName + "' in 'field list'"}
	}
	entry := s.entries[idx]
	lower := strings.ToLower(colName)
	for j, col := range entry.columns {
		if strings.ToLower(col.Name) == lower {
			return entry.rteIdx, j + 1, nil
		}
	}
	return 0, 0, &Error{Code: 1054, Message: "Unknown column '" + tableName + "." + colName + "' in 'field list'"}
}

// getColumns returns the columns for a scope entry by effective name.
func (s *analyzerScope) getColumns(tableName string) []*Column {
	key := strings.ToLower(tableName)
	idx, ok := s.byName[key]
	if !ok {
		return nil
	}
	return s.entries[idx].columns
}

// allEntries returns all scope entries in insertion order.
func (s *analyzerScope) allEntries() []scopeEntry {
	return s.entries
}
```

```go
// mysql/catalog/analyze.go
package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// AnalyzeSelectStmt analyzes a parsed SELECT statement and returns a
// semantically resolved Query IR.
//
// The analysis proceeds in three steps:
//  1. Build scope from FROM clause (populate RangeTable + scope)
//  2. Resolve target list expressions against scope (populate TargetList)
//  3. Resolve WHERE/GROUP BY/HAVING/ORDER BY/LIMIT
func (c *Catalog) AnalyzeSelectStmt(stmt *nodes.SelectStmt) (*Query, error) {
	q := &Query{
		CommandType: CmdSelect,
	}

	scope := newScope()

	// Step 1: FROM clause → RangeTable + scope
	if err := c.analyzeFromClause(stmt.From, q, scope); err != nil {
		return nil, err
	}

	// Step 2: Target list
	if err := c.analyzeTargetList(stmt.TargetList, q, scope); err != nil {
		return nil, err
	}

	// Step 3: WHERE
	if stmt.Where != nil {
		expr, err := c.analyzeExpr(stmt.Where, scope)
		if err != nil {
			return nil, err
		}
		if q.JoinTree == nil {
			q.JoinTree = &JoinTreeQ{}
		}
		q.JoinTree.Quals = expr
	}

	// Ensure JoinTree is never nil
	if q.JoinTree == nil {
		q.JoinTree = &JoinTreeQ{}
	}

	return q, nil
}

// analyzeFromClause processes the FROM clause into the Query's RangeTable.
func (c *Catalog) analyzeFromClause(from []nodes.TableExpr, q *Query, scope *analyzerScope) error {
	for _, tableExpr := range from {
		switch t := tableExpr.(type) {
		case *nodes.TableRef:
			rte, columns, err := c.analyzeTableRef(t)
			if err != nil {
				return err
			}
			rteIdx := len(q.RangeTable)
			q.RangeTable = append(q.RangeTable, rte)

			// Add to scope with effective name
			effectiveName := rte.ERef
			scope.add(effectiveName, rteIdx, columns)

			// Add to JoinTree
			if q.JoinTree == nil {
				q.JoinTree = &JoinTreeQ{}
			}
			q.JoinTree.FromList = append(q.JoinTree.FromList, &RangeTableRefQ{RTIndex: rteIdx})

		default:
			return fmt.Errorf("unsupported FROM clause element: %T", tableExpr)
		}
	}
	return nil
}

// analyzeTableRef resolves a table reference to an RTE and its columns.
func (c *Catalog) analyzeTableRef(ref *nodes.TableRef) (*RangeTableEntryQ, []*Column, error) {
	dbName := ref.Schema
	if dbName == "" {
		dbName = c.currentDB
	}
	db := c.GetDatabase(dbName)
	if db == nil {
		return nil, nil, &Error{Code: 1049, Message: "Unknown database '" + dbName + "'"}
	}

	// Check tables first, then views
	table := db.GetTable(ref.Name)
	if table != nil {
		colNames := make([]string, len(table.Columns))
		for i, col := range table.Columns {
			colNames[i] = col.Name
		}

		eref := ref.Name
		if ref.Alias != "" {
			eref = ref.Alias
		}

		rte := &RangeTableEntryQ{
			Kind:      RTERelation,
			DBName:    db.Name,
			TableName: table.Name,
			Alias:     ref.Alias,
			ERef:      eref,
			ColNames:  colNames,
		}
		return rte, table.Columns, nil
	}

	// Check views
	view := db.Views[toLower(ref.Name)]
	if view != nil {
		// For now, create an opaque RTE for views (Phase 2 will add AnalyzedQuery)
		eref := ref.Name
		if ref.Alias != "" {
			eref = ref.Alias
		}
		colNames := make([]string, len(view.Columns))
		copy(colNames, view.Columns)
		rte := &RangeTableEntryQ{
			Kind:          RTERelation,
			DBName:        db.Name,
			TableName:     view.Name,
			Alias:         ref.Alias,
			ERef:          eref,
			ColNames:      colNames,
			IsView:        true,
			ViewAlgorithm: viewAlgorithmFromString(view.Algorithm),
		}
		// Views don't have Column objects — create stub slice for scope
		// Phase 3 will fill in types from view analysis
		var stubCols []*Column
		for i, name := range view.Columns {
			stubCols = append(stubCols, &Column{Position: i, Name: name})
		}
		return rte, stubCols, nil
	}

	return nil, nil, &Error{Code: 1146, Message: "Table '" + dbName + "." + ref.Name + "' doesn't exist"}
}

func viewAlgorithmFromString(s string) ViewAlgorithm {
	switch strings.ToUpper(s) {
	case "MERGE":
		return ViewAlgMerge
	case "TEMPTABLE":
		return ViewAlgTemptable
	case "UNDEFINED", "":
		return ViewAlgUndefined
	default:
		return ViewAlgUndefined
	}
}
```

### Step 4: Write analyze_targetlist.go

```go
// mysql/catalog/analyze_targetlist.go
package catalog

import (
	"fmt"
	"strconv"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// analyzeTargetList processes the SELECT list into Query.TargetList.
func (c *Catalog) analyzeTargetList(targetList []nodes.ExprNode, q *Query, scope *analyzerScope) error {
	resNo := 1
	for _, item := range targetList {
		entries, err := c.analyzeTargetEntry(item, q, scope, &resNo)
		if err != nil {
			return err
		}
		q.TargetList = append(q.TargetList, entries...)
	}
	return nil
}

// analyzeTargetEntry processes one item in the SELECT list. May produce
// multiple TargetEntryQ values when expanding * or t.*.
func (c *Catalog) analyzeTargetEntry(item nodes.ExprNode, q *Query, scope *analyzerScope, resNo *int) ([]*TargetEntryQ, error) {
	// Handle ResTarget (aliased expression): SELECT expr AS alias
	if rt, ok := item.(*nodes.ResTarget); ok {
		expr, err := c.analyzeExpr(rt.Val, scope)
		if err != nil {
			return nil, err
		}
		te := &TargetEntryQ{
			Expr:    expr,
			ResNo:   *resNo,
			ResName: rt.Name,
			ResJunk: false,
		}
		// Fill provenance if Expr is a simple VarExprQ
		fillProvenance(te, q)
		*resNo++
		return []*TargetEntryQ{te}, nil
	}

	// Handle bare star: SELECT *
	if _, ok := item.(*nodes.StarExpr); ok {
		return c.expandStar("", q, scope, resNo)
	}

	// Handle qualified star: SELECT t.*
	if cr, ok := item.(*nodes.ColumnRef); ok && cr.Star {
		return c.expandStar(cr.Table, q, scope, resNo)
	}

	// Handle bare expression (no alias): SELECT expr
	expr, err := c.analyzeExpr(item, scope)
	if err != nil {
		return nil, err
	}
	te := &TargetEntryQ{
		Expr:    expr,
		ResNo:   *resNo,
		ResName: deriveColumnName(item, expr),
		ResJunk: false,
	}
	fillProvenance(te, q)
	*resNo++
	return []*TargetEntryQ{te}, nil
}

// expandStar expands * or table.* into individual TargetEntryQ values.
// If tableName is empty, expands all tables in scope.
func (c *Catalog) expandStar(tableName string, q *Query, scope *analyzerScope, resNo *int) ([]*TargetEntryQ, error) {
	var entries []*TargetEntryQ

	if tableName != "" {
		// Qualified: t.*
		cols := scope.getColumns(tableName)
		if cols == nil {
			return nil, &Error{Code: 1054, Message: "Unknown table '" + tableName + "'"}
		}
		rteIdx, _, err := scope.resolveQualifiedColumn(tableName, cols[0].Name)
		if err != nil {
			return nil, err
		}
		rte := q.RangeTable[rteIdx]
		for i, col := range cols {
			te := &TargetEntryQ{
				Expr:         &VarExprQ{RangeIdx: rteIdx, AttNum: i + 1},
				ResNo:        *resNo,
				ResName:      col.Name,
				ResJunk:      false,
				ResOrigDB:    rte.DBName,
				ResOrigTable: rte.TableName,
				ResOrigCol:   col.Name,
			}
			entries = append(entries, te)
			*resNo++
		}
	} else {
		// Unqualified: SELECT *
		for _, entry := range scope.allEntries() {
			rte := q.RangeTable[entry.rteIdx]
			for i, col := range entry.columns {
				te := &TargetEntryQ{
					Expr:         &VarExprQ{RangeIdx: entry.rteIdx, AttNum: i + 1},
					ResNo:        *resNo,
					ResName:      col.Name,
					ResJunk:      false,
					ResOrigDB:    rte.DBName,
					ResOrigTable: rte.TableName,
					ResOrigCol:   col.Name,
				}
				entries = append(entries, te)
				*resNo++
			}
		}
	}

	if len(entries) == 0 {
		return nil, &Error{Code: 1064, Message: "No tables used"}
	}
	return entries, nil
}

// deriveColumnName generates a column name for an unaliased SELECT expression.
// MySQL uses the original expression text, but we approximate with simple rules.
func deriveColumnName(astNode nodes.ExprNode, analyzed AnalyzedExpr) string {
	switch n := astNode.(type) {
	case *nodes.ColumnRef:
		return n.Column
	case *nodes.IntLit:
		return strconv.FormatInt(n.Value, 10)
	case *nodes.StringLit:
		return n.Value
	case *nodes.FloatLit:
		return n.Value
	case *nodes.NullLit:
		return "NULL"
	case *nodes.BoolLit:
		if n.Value {
			return "TRUE"
		}
		return "FALSE"
	case *nodes.FuncCallExpr:
		return fmt.Sprintf("%s(...)", n.Name)
	default:
		// Fallback: use a generic name
		return fmt.Sprintf("%T", astNode)
	}
}

// fillProvenance sets ResOrigDB/Table/Col when Expr is a direct VarExprQ.
func fillProvenance(te *TargetEntryQ, q *Query) {
	v, ok := te.Expr.(*VarExprQ)
	if !ok || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	if v.AttNum < 1 || v.AttNum > len(rte.ColNames) {
		return
	}
	te.ResOrigDB = rte.DBName
	te.ResOrigTable = rte.TableName
	te.ResOrigCol = rte.ColNames[v.AttNum-1]
}
```

### Step 5: Write analyze_expr.go (minimal — just column refs and literals)

```go
// mysql/catalog/analyze_expr.go
package catalog

import (
	"fmt"
	"strconv"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// analyzeExpr analyzes a single expression node from the AST and returns
// the corresponding AnalyzedExpr in the Query IR.
func (c *Catalog) analyzeExpr(expr nodes.ExprNode, scope *analyzerScope) (AnalyzedExpr, error) {
	if expr == nil {
		return nil, nil
	}
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return c.analyzeColumnRef(e, scope)
	case *nodes.IntLit:
		return &ConstExprQ{Value: strconv.FormatInt(e.Value, 10)}, nil
	case *nodes.StringLit:
		return &ConstExprQ{Value: e.Value}, nil
	case *nodes.FloatLit:
		return &ConstExprQ{Value: e.Value}, nil
	case *nodes.NullLit:
		return &ConstExprQ{IsNull: true, Value: "NULL"}, nil
	case *nodes.BoolLit:
		if e.Value {
			return &ConstExprQ{Value: "TRUE"}, nil
		}
		return &ConstExprQ{Value: "FALSE"}, nil
	case *nodes.FuncCallExpr:
		return c.analyzeFuncCall(e, scope)
	case *nodes.ParenExpr:
		return c.analyzeExpr(e.Expr, scope)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// analyzeColumnRef resolves a column reference to a VarExprQ.
func (c *Catalog) analyzeColumnRef(ref *nodes.ColumnRef, scope *analyzerScope) (AnalyzedExpr, error) {
	if ref.Table != "" {
		// Qualified: table.column or schema.table.column
		tableName := ref.Table
		if ref.Schema != "" {
			// Three-part: schema.table.column — resolve schema part separately
			// For now, just use table part for scope resolution
			tableName = ref.Table
		}
		rteIdx, attNum, err := scope.resolveQualifiedColumn(tableName, ref.Column)
		if err != nil {
			return nil, err
		}
		return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum}, nil
	}

	// Unqualified column
	rteIdx, attNum, err := scope.resolveColumn(ref.Column)
	if err != nil {
		return nil, err
	}
	return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum}, nil
}

// analyzeFuncCall analyzes a function call expression.
func (c *Catalog) analyzeFuncCall(fc *nodes.FuncCallExpr, scope *analyzerScope) (AnalyzedExpr, error) {
	var args []AnalyzedExpr
	for _, arg := range fc.Args {
		a, err := c.analyzeExpr(arg, scope)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}

	return &FuncCallExprQ{
		Name:        fc.Name,
		Args:        args,
		IsAggregate: isAggregateFunc(fc.Name),
		Distinct:    fc.Distinct,
	}, nil
}

// isAggregateFunc returns true if the function name is a known aggregate.
func isAggregateFunc(name string) bool {
	switch name {
	case "count", "COUNT",
		"sum", "SUM",
		"avg", "AVG",
		"min", "MIN",
		"max", "MAX",
		"group_concat", "GROUP_CONCAT",
		"bit_and", "BIT_AND",
		"bit_or", "BIT_OR",
		"bit_xor", "BIT_XOR",
		"std", "STD",
		"stddev", "STDDEV",
		"stddev_pop", "STDDEV_POP",
		"stddev_samp", "STDDEV_SAMP",
		"var_pop", "VAR_POP",
		"var_samp", "VAR_SAMP",
		"variance", "VARIANCE",
		"json_arrayagg", "JSON_ARRAYAGG",
		"json_objectagg", "JSON_OBJECTAGG":
		return true
	}
	return false
}
```

### Step 6: Run test to verify scenario 1.1 passes

Run: `go test ./mysql/catalog/ -run TestAnalyze_1_1 -v -count=1`
Expected: PASS

### Step 7: Commit

```bash
git add mysql/catalog/analyze.go mysql/catalog/analyze_targetlist.go mysql/catalog/analyze_expr.go mysql/catalog/scope.go mysql/catalog/analyze_test.go
git commit -m "feat(mysql/catalog): add AnalyzeSelectStmt foundation (Phase 1a Batch 1)

Implements the analyzer entry point and basic expression/scope resolution.
Passes scenario 1.1 (SELECT 1)."
```

---

## Task 2: Scenarios 1.2 and 1.3 — column references and aliases

**Files:**
- Modify: `mysql/catalog/analyze_test.go`

### Step 1: Write the failing test for 1.2 and 1.3

```go
// Add to analyze_test.go

func TestAnalyze_1_2_ColumnFromTable(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(200),
		department_id INT NOT NULL,
		salary DECIMAL(10,2) NOT NULL,
		hire_date DATE NOT NULL,
		is_active TINYINT(1) NOT NULL DEFAULT 1,
		notes TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	sel := parseSelect(t, "SELECT name FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// RangeTable
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable len = %d, want 1", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Kind != RTERelation {
		t.Errorf("Kind = %d, want RTERelation", rte.Kind)
	}
	if rte.DBName != "testdb" {
		t.Errorf("DBName = %q, want %q", rte.DBName, "testdb")
	}
	if rte.TableName != "employees" {
		t.Errorf("TableName = %q, want %q", rte.TableName, "employees")
	}
	if rte.ERef != "employees" {
		t.Errorf("ERef = %q, want %q", rte.ERef, "employees")
	}
	if len(rte.ColNames) != 9 {
		t.Fatalf("ColNames len = %d, want 9", len(rte.ColNames))
	}

	// JoinTree
	if len(q.JoinTree.FromList) != 1 {
		t.Fatalf("FromList len = %d, want 1", len(q.JoinTree.FromList))
	}
	ref, ok := q.JoinTree.FromList[0].(*RangeTableRefQ)
	if !ok {
		t.Fatalf("FromList[0] type = %T, want *RangeTableRefQ", q.JoinTree.FromList[0])
	}
	if ref.RTIndex != 0 {
		t.Errorf("RTIndex = %d, want 0", ref.RTIndex)
	}

	// TargetList
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList len = %d, want 1", len(q.TargetList))
	}
	te := q.TargetList[0]
	if te.ResName != "name" {
		t.Errorf("ResName = %q, want %q", te.ResName, "name")
	}
	v, ok := te.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("Expr type = %T, want *VarExprQ", te.Expr)
	}
	if v.RangeIdx != 0 || v.AttNum != 2 {
		t.Errorf("VarExprQ = (RangeIdx:%d, AttNum:%d), want (0, 2)", v.RangeIdx, v.AttNum)
	}

	// Provenance
	if te.ResOrigDB != "testdb" || te.ResOrigTable != "employees" || te.ResOrigCol != "name" {
		t.Errorf("Provenance = (%q, %q, %q), want (testdb, employees, name)",
			te.ResOrigDB, te.ResOrigTable, te.ResOrigCol)
	}
}

func TestAnalyze_1_3_MultipleColumnsWithAlias(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(200),
		department_id INT NOT NULL,
		salary DECIMAL(10,2) NOT NULL,
		hire_date DATE NOT NULL,
		is_active TINYINT(1) NOT NULL DEFAULT 1,
		notes TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	sel := parseSelect(t, "SELECT id, name AS employee_name, salary FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 3 {
		t.Fatalf("TargetList len = %d, want 3", len(q.TargetList))
	}

	// id
	te0 := q.TargetList[0]
	if te0.ResName != "id" {
		t.Errorf("TargetList[0].ResName = %q, want %q", te0.ResName, "id")
	}
	if v, ok := te0.Expr.(*VarExprQ); !ok || v.AttNum != 1 {
		t.Errorf("TargetList[0].Expr = %+v, want VarExprQ{AttNum:1}", te0.Expr)
	}

	// name AS employee_name
	te1 := q.TargetList[1]
	if te1.ResName != "employee_name" {
		t.Errorf("TargetList[1].ResName = %q, want %q", te1.ResName, "employee_name")
	}
	if v, ok := te1.Expr.(*VarExprQ); !ok || v.AttNum != 2 {
		t.Errorf("TargetList[1].Expr = %+v, want VarExprQ{AttNum:2}", te1.Expr)
	}

	// salary
	te2 := q.TargetList[2]
	if te2.ResName != "salary" {
		t.Errorf("TargetList[2].ResName = %q, want %q", te2.ResName, "salary")
	}
	if v, ok := te2.Expr.(*VarExprQ); !ok || v.AttNum != 5 {
		t.Errorf("TargetList[2].Expr = %+v, want VarExprQ{AttNum:5}", te2.Expr)
	}

	// ResNo sequence
	for i, te := range q.TargetList {
		if te.ResNo != i+1 {
			t.Errorf("TargetList[%d].ResNo = %d, want %d", i, te.ResNo, i+1)
		}
	}
}
```

### Step 2: Run tests — should pass immediately (code from Task 1 handles these)

Run: `go test ./mysql/catalog/ -run "TestAnalyze_1_[23]" -v -count=1`
Expected: PASS

### Step 3: Commit

```bash
git add mysql/catalog/analyze_test.go
git commit -m "test(mysql/catalog): add scenarios 1.2–1.3 (column refs, aliases)"
```

---

## Task 3: Scenarios 1.4 and 1.5 — star expansion

**Files:**
- Modify: `mysql/catalog/analyze_test.go`

### Step 1: Write the failing test for 1.4 and 1.5

```go
// Add to analyze_test.go

func TestAnalyze_1_4_StarExpansion(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(200),
		department_id INT NOT NULL,
		salary DECIMAL(10,2) NOT NULL,
		hire_date DATE NOT NULL,
		is_active TINYINT(1) NOT NULL DEFAULT 1,
		notes TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	sel := parseSelect(t, "SELECT * FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// Should expand to 9 columns
	if len(q.TargetList) != 9 {
		t.Fatalf("TargetList len = %d, want 9", len(q.TargetList))
	}

	expectedCols := []string{"id", "name", "email", "department_id", "salary", "hire_date", "is_active", "notes", "created_at"}
	for i, te := range q.TargetList {
		if te.ResName != expectedCols[i] {
			t.Errorf("TargetList[%d].ResName = %q, want %q", i, te.ResName, expectedCols[i])
		}
		v, ok := te.Expr.(*VarExprQ)
		if !ok {
			t.Errorf("TargetList[%d].Expr type = %T, want *VarExprQ", i, te.Expr)
			continue
		}
		if v.RangeIdx != 0 || v.AttNum != i+1 {
			t.Errorf("TargetList[%d] VarExprQ = (RangeIdx:%d, AttNum:%d), want (0, %d)",
				i, v.RangeIdx, v.AttNum, i+1)
		}
		if te.ResOrigCol != expectedCols[i] {
			t.Errorf("TargetList[%d].ResOrigCol = %q, want %q", i, te.ResOrigCol, expectedCols[i])
		}
	}
}

func TestAnalyze_1_5_QualifiedStar(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(200),
		department_id INT NOT NULL,
		salary DECIMAL(10,2) NOT NULL,
		hire_date DATE NOT NULL,
		is_active TINYINT(1) NOT NULL DEFAULT 1,
		notes TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	sel := parseSelect(t, "SELECT e.* FROM employees AS e")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// RTE should have alias
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable len = %d, want 1", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Alias != "e" {
		t.Errorf("Alias = %q, want %q", rte.Alias, "e")
	}
	if rte.ERef != "e" {
		t.Errorf("ERef = %q, want %q", rte.ERef, "e")
	}

	// Same 9 columns as 1.4
	if len(q.TargetList) != 9 {
		t.Fatalf("TargetList len = %d, want 9", len(q.TargetList))
	}
	for i, te := range q.TargetList {
		v, ok := te.Expr.(*VarExprQ)
		if !ok {
			t.Errorf("TargetList[%d].Expr type = %T, want *VarExprQ", i, te.Expr)
			continue
		}
		if v.RangeIdx != 0 || v.AttNum != i+1 {
			t.Errorf("TargetList[%d] VarExprQ = (RangeIdx:%d, AttNum:%d), want (0, %d)",
				i, v.RangeIdx, v.AttNum, i+1)
		}
	}
}
```

### Step 2: Run tests

Run: `go test ./mysql/catalog/ -run "TestAnalyze_1_[45]" -v -count=1`
Expected: PASS (star expansion logic already in Task 1)

### Step 3: Commit

```bash
git add mysql/catalog/analyze_test.go
git commit -m "test(mysql/catalog): add scenarios 1.4–1.5 (star expansion)"
```

---

## Task 4: Run all 5 scenarios together + fix any issues

### Step 1: Run all Batch 1 tests

Run: `go test ./mysql/catalog/ -run "TestAnalyze_1_" -v -count=1`
Expected: All 5 PASS

### Step 2: Run full package test suite (ensure no regressions)

Run: `go test ./mysql/catalog/ -count=1 -timeout 120s`
Expected: All existing tests still pass

### Step 3: Run go vet

Run: `go vet ./mysql/catalog/...`
Expected: Clean

### Step 4: Commit if any fixes were needed

---

## Task 5: Verify no dead code, clean up

### Step 1: Check for unused imports / functions

Run: `go vet ./mysql/catalog/...`

### Step 2: Ensure all new files have the standard package header

### Step 3: Final commit

```bash
git add -A mysql/catalog/
git commit -m "feat(mysql/catalog): Phase 1a Batch 1 complete — analyzer foundation

Implements AnalyzeSelectStmt with single-table scope resolution, target
list analysis (columns, aliases, * expansion), and basic expression
analysis (column refs, literals, function calls).

Passes scenarios 1.1–1.5 from SCENARIOS-mysql-analyzer-1a.md."
```

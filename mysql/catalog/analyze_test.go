package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// parseSelect parses a single SELECT statement and returns the SelectStmt node.
func parseSelect(t *testing.T, sql string) *nodes.SelectStmt {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", list.Len())
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *ast.SelectStmt, got %T", list.Items[0])
	}
	return sel
}

// TestAnalyze_1_1_BareLiteral tests SELECT 1: no tables, just a literal.
func TestAnalyze_1_1_BareLiteral(t *testing.T) {
	c := wtSetup(t)
	sel := parseSelect(t, "SELECT 1")

	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// CommandType
	if q.CommandType != CmdSelect {
		t.Errorf("CommandType: want CmdSelect, got %v", q.CommandType)
	}

	// RangeTable should be empty (no FROM clause).
	if len(q.RangeTable) != 0 {
		t.Errorf("RangeTable: want 0 entries, got %d", len(q.RangeTable))
	}

	// JoinTree should be non-nil with empty FromList and nil Quals.
	if q.JoinTree == nil {
		t.Fatal("JoinTree: want non-nil, got nil")
	}
	if len(q.JoinTree.FromList) != 0 {
		t.Errorf("JoinTree.FromList: want 0 entries, got %d", len(q.JoinTree.FromList))
	}
	if q.JoinTree.Quals != nil {
		t.Errorf("JoinTree.Quals: want nil, got %v", q.JoinTree.Quals)
	}

	// TargetList should have exactly 1 entry.
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}

	te := q.TargetList[0]

	// ResNo = 1
	if te.ResNo != 1 {
		t.Errorf("ResNo: want 1, got %d", te.ResNo)
	}

	// ResName = "1"
	if te.ResName != "1" {
		t.Errorf("ResName: want %q, got %q", "1", te.ResName)
	}

	// ResJunk = false
	if te.ResJunk {
		t.Errorf("ResJunk: want false, got true")
	}

	// Expr should be ConstExprQ with Value "1".
	constExpr, ok := te.Expr.(*ConstExprQ)
	if !ok {
		t.Fatalf("Expr: want *ConstExprQ, got %T", te.Expr)
	}
	if constExpr.Value != "1" {
		t.Errorf("ConstExprQ.Value: want %q, got %q", "1", constExpr.Value)
	}
	if constExpr.IsNull {
		t.Errorf("ConstExprQ.IsNull: want false, got true")
	}
}

// TestAnalyze_1_2_ColumnFromTable tests SELECT name FROM employees.
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

	// RangeTable: 1 entry
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1 entry, got %d", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Kind != RTERelation {
		t.Errorf("RTE.Kind: want RTERelation, got %v", rte.Kind)
	}
	if rte.DBName != "testdb" {
		t.Errorf("RTE.DBName: want %q, got %q", "testdb", rte.DBName)
	}
	if rte.TableName != "employees" {
		t.Errorf("RTE.TableName: want %q, got %q", "employees", rte.TableName)
	}
	if rte.ERef != "employees" {
		t.Errorf("RTE.ERef: want %q, got %q", "employees", rte.ERef)
	}

	// RTE ColNames: 9 entries
	if len(rte.ColNames) != 9 {
		t.Errorf("RTE.ColNames: want 9 entries, got %d", len(rte.ColNames))
	}

	// JoinTree.FromList: 1 RangeTableRefQ
	if q.JoinTree == nil {
		t.Fatal("JoinTree: want non-nil, got nil")
	}
	if len(q.JoinTree.FromList) != 1 {
		t.Fatalf("JoinTree.FromList: want 1 entry, got %d", len(q.JoinTree.FromList))
	}
	rtRef, ok := q.JoinTree.FromList[0].(*RangeTableRefQ)
	if !ok {
		t.Fatalf("FromList[0]: want *RangeTableRefQ, got %T", q.JoinTree.FromList[0])
	}
	if rtRef.RTIndex != 0 {
		t.Errorf("RangeTableRefQ.RTIndex: want 0, got %d", rtRef.RTIndex)
	}

	// TargetList: 1 entry
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}
	te := q.TargetList[0]

	varExpr, ok := te.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("Expr: want *VarExprQ, got %T", te.Expr)
	}
	if varExpr.RangeIdx != 0 {
		t.Errorf("VarExprQ.RangeIdx: want 0, got %d", varExpr.RangeIdx)
	}
	if varExpr.AttNum != 2 {
		t.Errorf("VarExprQ.AttNum: want 2, got %d", varExpr.AttNum)
	}
	if te.ResName != "name" {
		t.Errorf("ResName: want %q, got %q", "name", te.ResName)
	}

	// Provenance
	if te.ResOrigDB != "testdb" {
		t.Errorf("ResOrigDB: want %q, got %q", "testdb", te.ResOrigDB)
	}
	if te.ResOrigTable != "employees" {
		t.Errorf("ResOrigTable: want %q, got %q", "employees", te.ResOrigTable)
	}
	if te.ResOrigCol != "name" {
		t.Errorf("ResOrigCol: want %q, got %q", "name", te.ResOrigCol)
	}
}

// TestAnalyze_1_3_MultipleColumnsWithAlias tests SELECT id, name AS employee_name, salary FROM employees.
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

	// TargetList: 3 entries
	if len(q.TargetList) != 3 {
		t.Fatalf("TargetList: want 3 entries, got %d", len(q.TargetList))
	}

	// Entry 0: id
	te0 := q.TargetList[0]
	v0, ok := te0.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *VarExprQ, got %T", te0.Expr)
	}
	if v0.AttNum != 1 {
		t.Errorf("TargetList[0] AttNum: want 1, got %d", v0.AttNum)
	}
	if te0.ResName != "id" {
		t.Errorf("TargetList[0] ResName: want %q, got %q", "id", te0.ResName)
	}
	if te0.ResNo != 1 {
		t.Errorf("TargetList[0] ResNo: want 1, got %d", te0.ResNo)
	}

	// Entry 1: name AS employee_name
	te1 := q.TargetList[1]
	v1, ok := te1.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *VarExprQ, got %T", te1.Expr)
	}
	if v1.AttNum != 2 {
		t.Errorf("TargetList[1] AttNum: want 2, got %d", v1.AttNum)
	}
	if te1.ResName != "employee_name" {
		t.Errorf("TargetList[1] ResName: want %q, got %q", "employee_name", te1.ResName)
	}
	if te1.ResNo != 2 {
		t.Errorf("TargetList[1] ResNo: want 2, got %d", te1.ResNo)
	}

	// Entry 2: salary
	te2 := q.TargetList[2]
	v2, ok := te2.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("TargetList[2].Expr: want *VarExprQ, got %T", te2.Expr)
	}
	if v2.AttNum != 5 {
		t.Errorf("TargetList[2] AttNum: want 5, got %d", v2.AttNum)
	}
	if te2.ResName != "salary" {
		t.Errorf("TargetList[2] ResName: want %q, got %q", "salary", te2.ResName)
	}
	if te2.ResNo != 3 {
		t.Errorf("TargetList[2] ResNo: want 3, got %d", te2.ResNo)
	}

	// All should have provenance
	for i, te := range q.TargetList {
		if te.ResOrigDB == "" {
			t.Errorf("TargetList[%d] ResOrigDB: want non-empty, got empty", i)
		}
		if te.ResOrigTable == "" {
			t.Errorf("TargetList[%d] ResOrigTable: want non-empty, got empty", i)
		}
		if te.ResOrigCol == "" {
			t.Errorf("TargetList[%d] ResOrigCol: want non-empty, got empty", i)
		}
	}
}

// TestAnalyze_1_4_StarExpansion tests SELECT * FROM employees.
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

	// TargetList: 9 entries (one per column)
	if len(q.TargetList) != 9 {
		t.Fatalf("TargetList: want 9 entries, got %d", len(q.TargetList))
	}

	wantNames := []string{"id", "name", "email", "department_id", "salary", "hire_date", "is_active", "notes", "created_at"}

	for i, te := range q.TargetList {
		v, ok := te.Expr.(*VarExprQ)
		if !ok {
			t.Fatalf("TargetList[%d].Expr: want *VarExprQ, got %T", i, te.Expr)
		}
		wantAttNum := i + 1
		if v.AttNum != wantAttNum {
			t.Errorf("TargetList[%d] AttNum: want %d, got %d", i, wantAttNum, v.AttNum)
		}
		if te.ResName != wantNames[i] {
			t.Errorf("TargetList[%d] ResName: want %q, got %q", i, wantNames[i], te.ResName)
		}
		if te.ResOrigCol == "" {
			t.Errorf("TargetList[%d] ResOrigCol: want non-empty, got empty", i)
		}
	}
}

// TestAnalyze_1_5_QualifiedStar tests SELECT e.* FROM employees AS e.
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

	// RangeTable[0]: Alias="e", ERef="e"
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1 entry, got %d", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Alias != "e" {
		t.Errorf("RTE.Alias: want %q, got %q", "e", rte.Alias)
	}
	if rte.ERef != "e" {
		t.Errorf("RTE.ERef: want %q, got %q", "e", rte.ERef)
	}

	// TargetList: 9 entries same as 1.4
	if len(q.TargetList) != 9 {
		t.Fatalf("TargetList: want 9 entries, got %d", len(q.TargetList))
	}

	for i, te := range q.TargetList {
		v, ok := te.Expr.(*VarExprQ)
		if !ok {
			t.Fatalf("TargetList[%d].Expr: want *VarExprQ, got %T", i, te.Expr)
		}
		if v.RangeIdx != 0 {
			t.Errorf("TargetList[%d] RangeIdx: want 0, got %d", i, v.RangeIdx)
		}
		wantAttNum := i + 1
		if v.AttNum != wantAttNum {
			t.Errorf("TargetList[%d] AttNum: want %d, got %d", i, wantAttNum, v.AttNum)
		}
	}
}

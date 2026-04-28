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

// employeesTableDDL is the shared DDL for the employees table used across Batch 2 tests.
const employeesTableDDL = `CREATE TABLE employees (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(200),
	department_id INT NOT NULL,
	salary DECIMAL(10,2) NOT NULL,
	hire_date DATE NOT NULL,
	is_active TINYINT(1) NOT NULL DEFAULT 1,
	notes TEXT,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`

// TestAnalyze_2_1_WhereSimpleEquality tests WHERE id = 1.
func TestAnalyze_2_1_WhereSimpleEquality(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name FROM employees WHERE id = 1")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.JoinTree == nil {
		t.Fatal("JoinTree: want non-nil, got nil")
	}
	if q.JoinTree.Quals == nil {
		t.Fatal("JoinTree.Quals: want non-nil, got nil")
	}

	opExpr, ok := q.JoinTree.Quals.(*OpExprQ)
	if !ok {
		t.Fatalf("Quals: want *OpExprQ, got %T", q.JoinTree.Quals)
	}
	if opExpr.Op != "=" {
		t.Errorf("OpExprQ.Op: want %q, got %q", "=", opExpr.Op)
	}

	leftVar, ok := opExpr.Left.(*VarExprQ)
	if !ok {
		t.Fatalf("OpExprQ.Left: want *VarExprQ, got %T", opExpr.Left)
	}
	if leftVar.AttNum != 1 {
		t.Errorf("Left.AttNum: want 1, got %d", leftVar.AttNum)
	}

	rightConst, ok := opExpr.Right.(*ConstExprQ)
	if !ok {
		t.Fatalf("OpExprQ.Right: want *ConstExprQ, got %T", opExpr.Right)
	}
	if rightConst.Value != "1" {
		t.Errorf("Right.Value: want %q, got %q", "1", rightConst.Value)
	}
}

// TestAnalyze_2_2_WhereAndOr tests WHERE is_active = 1 AND (department_id = 1 OR department_id = 2).
func TestAnalyze_2_2_WhereAndOr(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name FROM employees WHERE is_active = 1 AND (department_id = 1 OR department_id = 2)")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.JoinTree == nil || q.JoinTree.Quals == nil {
		t.Fatal("JoinTree.Quals: want non-nil")
	}

	// Top level: BoolExprQ with BoolAnd
	andExpr, ok := q.JoinTree.Quals.(*BoolExprQ)
	if !ok {
		t.Fatalf("Quals: want *BoolExprQ, got %T", q.JoinTree.Quals)
	}
	if andExpr.Op != BoolAnd {
		t.Errorf("BoolExprQ.Op: want BoolAnd, got %v", andExpr.Op)
	}
	if len(andExpr.Args) != 2 {
		t.Fatalf("BoolExprQ.Args: want 2 entries, got %d", len(andExpr.Args))
	}

	// Left: is_active = 1
	leftOp, ok := andExpr.Args[0].(*OpExprQ)
	if !ok {
		t.Fatalf("Args[0]: want *OpExprQ, got %T", andExpr.Args[0])
	}
	if leftOp.Op != "=" {
		t.Errorf("Args[0].Op: want %q, got %q", "=", leftOp.Op)
	}

	// Right: BoolExprQ with BoolOr
	orExpr, ok := andExpr.Args[1].(*BoolExprQ)
	if !ok {
		t.Fatalf("Args[1]: want *BoolExprQ, got %T", andExpr.Args[1])
	}
	if orExpr.Op != BoolOr {
		t.Errorf("Args[1].Op: want BoolOr, got %v", orExpr.Op)
	}
	if len(orExpr.Args) != 2 {
		t.Fatalf("BoolOr.Args: want 2 entries, got %d", len(orExpr.Args))
	}

	// Each OR arg should be OpExprQ with "="
	for i, arg := range orExpr.Args {
		op, ok := arg.(*OpExprQ)
		if !ok {
			t.Fatalf("BoolOr.Args[%d]: want *OpExprQ, got %T", i, arg)
		}
		if op.Op != "=" {
			t.Errorf("BoolOr.Args[%d].Op: want %q, got %q", i, "=", op.Op)
		}
	}
}

// TestAnalyze_2_3_WhereIn tests WHERE department_id IN (1, 2, 3).
func TestAnalyze_2_3_WhereIn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name FROM employees WHERE department_id IN (1, 2, 3)")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.JoinTree == nil || q.JoinTree.Quals == nil {
		t.Fatal("JoinTree.Quals: want non-nil")
	}

	inExpr, ok := q.JoinTree.Quals.(*InListExprQ)
	if !ok {
		t.Fatalf("Quals: want *InListExprQ, got %T", q.JoinTree.Quals)
	}

	argVar, ok := inExpr.Arg.(*VarExprQ)
	if !ok {
		t.Fatalf("InListExprQ.Arg: want *VarExprQ, got %T", inExpr.Arg)
	}
	if argVar.AttNum != 4 {
		t.Errorf("Arg.AttNum: want 4, got %d", argVar.AttNum)
	}

	if len(inExpr.List) != 3 {
		t.Fatalf("InListExprQ.List: want 3 entries, got %d", len(inExpr.List))
	}
	for i, item := range inExpr.List {
		if _, ok := item.(*ConstExprQ); !ok {
			t.Errorf("List[%d]: want *ConstExprQ, got %T", i, item)
		}
	}

	if inExpr.Negated {
		t.Errorf("InListExprQ.Negated: want false, got true")
	}
}

// TestAnalyze_2_4_WhereBetween tests WHERE salary BETWEEN 50000 AND 100000.
func TestAnalyze_2_4_WhereBetween(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name FROM employees WHERE salary BETWEEN 50000 AND 100000")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.JoinTree == nil || q.JoinTree.Quals == nil {
		t.Fatal("JoinTree.Quals: want non-nil")
	}

	betExpr, ok := q.JoinTree.Quals.(*BetweenExprQ)
	if !ok {
		t.Fatalf("Quals: want *BetweenExprQ, got %T", q.JoinTree.Quals)
	}

	argVar, ok := betExpr.Arg.(*VarExprQ)
	if !ok {
		t.Fatalf("BetweenExprQ.Arg: want *VarExprQ, got %T", betExpr.Arg)
	}
	if argVar.AttNum != 5 {
		t.Errorf("Arg.AttNum: want 5, got %d", argVar.AttNum)
	}

	lowerConst, ok := betExpr.Lower.(*ConstExprQ)
	if !ok {
		t.Fatalf("BetweenExprQ.Lower: want *ConstExprQ, got %T", betExpr.Lower)
	}
	if lowerConst.Value != "50000" {
		t.Errorf("Lower.Value: want %q, got %q", "50000", lowerConst.Value)
	}

	upperConst, ok := betExpr.Upper.(*ConstExprQ)
	if !ok {
		t.Fatalf("BetweenExprQ.Upper: want *ConstExprQ, got %T", betExpr.Upper)
	}
	if upperConst.Value != "100000" {
		t.Errorf("Upper.Value: want %q, got %q", "100000", upperConst.Value)
	}

	if betExpr.Negated {
		t.Errorf("BetweenExprQ.Negated: want false, got true")
	}
}

// TestAnalyze_2_5_WhereIsNull tests WHERE email IS NOT NULL AND notes IS NULL.
func TestAnalyze_2_5_WhereIsNull(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name FROM employees WHERE email IS NOT NULL AND notes IS NULL")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.JoinTree == nil || q.JoinTree.Quals == nil {
		t.Fatal("JoinTree.Quals: want non-nil")
	}

	andExpr, ok := q.JoinTree.Quals.(*BoolExprQ)
	if !ok {
		t.Fatalf("Quals: want *BoolExprQ, got %T", q.JoinTree.Quals)
	}
	if andExpr.Op != BoolAnd {
		t.Errorf("BoolExprQ.Op: want BoolAnd, got %v", andExpr.Op)
	}
	if len(andExpr.Args) != 2 {
		t.Fatalf("BoolExprQ.Args: want 2 entries, got %d", len(andExpr.Args))
	}

	// Args[0]: email IS NOT NULL → NullTestExprQ{IsNull: false}
	nt0, ok := andExpr.Args[0].(*NullTestExprQ)
	if !ok {
		t.Fatalf("Args[0]: want *NullTestExprQ, got %T", andExpr.Args[0])
	}
	if nt0.IsNull {
		t.Errorf("Args[0].IsNull: want false, got true")
	}

	// Args[1]: notes IS NULL → NullTestExprQ{IsNull: true}
	nt1, ok := andExpr.Args[1].(*NullTestExprQ)
	if !ok {
		t.Fatalf("Args[1]: want *NullTestExprQ, got %T", andExpr.Args[1])
	}
	if !nt1.IsNull {
		t.Errorf("Args[1].IsNull: want true, got false")
	}
}

// TestAnalyze_3_1_GroupByColumn tests GROUP BY with a column reference.
func TestAnalyze_3_1_GroupByColumn(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT department_id, COUNT(*) FROM employees GROUP BY department_id")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// TargetList: 2 entries
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	// TargetList[0]: department_id → VarExprQ{AttNum:4}
	te0 := q.TargetList[0]
	v0, ok := te0.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *VarExprQ, got %T", te0.Expr)
	}
	if v0.AttNum != 4 {
		t.Errorf("TargetList[0] AttNum: want 4, got %d", v0.AttNum)
	}
	if te0.ResName != "department_id" {
		t.Errorf("TargetList[0] ResName: want %q, got %q", "department_id", te0.ResName)
	}

	// TargetList[1]: COUNT(*) → FuncCallExprQ with IsAggregate=true
	te1 := q.TargetList[1]
	fc1, ok := te1.Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *FuncCallExprQ, got %T", te1.Expr)
	}
	if fc1.Name != "count" {
		t.Errorf("FuncCallExprQ.Name: want %q, got %q", "count", fc1.Name)
	}
	if !fc1.IsAggregate {
		t.Errorf("FuncCallExprQ.IsAggregate: want true, got false")
	}

	// GroupClause: 1 entry referencing TargetIdx=1
	if len(q.GroupClause) != 1 {
		t.Fatalf("GroupClause: want 1 entry, got %d", len(q.GroupClause))
	}
	if q.GroupClause[0].TargetIdx != 1 {
		t.Errorf("GroupClause[0].TargetIdx: want 1, got %d", q.GroupClause[0].TargetIdx)
	}

	// HasAggs = true
	if !q.HasAggs {
		t.Errorf("HasAggs: want true, got false")
	}
}

// TestAnalyze_3_2_GroupByMultipleAggregates tests GROUP BY with multiple aggregates.
func TestAnalyze_3_2_GroupByMultipleAggregates(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT department_id, COUNT(*) AS cnt, SUM(salary) AS total_salary, AVG(salary) AS avg_salary FROM employees GROUP BY department_id")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// TargetList: 4 entries
	if len(q.TargetList) != 4 {
		t.Fatalf("TargetList: want 4 entries, got %d", len(q.TargetList))
	}

	// TargetList[0]: department_id
	if _, ok := q.TargetList[0].Expr.(*VarExprQ); !ok {
		t.Fatalf("TargetList[0].Expr: want *VarExprQ, got %T", q.TargetList[0].Expr)
	}

	// TargetList[1..3]: all FuncCallExprQ with IsAggregate=true
	for i := 1; i <= 3; i++ {
		fc, ok := q.TargetList[i].Expr.(*FuncCallExprQ)
		if !ok {
			t.Fatalf("TargetList[%d].Expr: want *FuncCallExprQ, got %T", i, q.TargetList[i].Expr)
		}
		if !fc.IsAggregate {
			t.Errorf("TargetList[%d] IsAggregate: want true, got false", i)
		}
	}

	// GroupClause references TargetIdx=1
	if len(q.GroupClause) != 1 {
		t.Fatalf("GroupClause: want 1 entry, got %d", len(q.GroupClause))
	}
	if q.GroupClause[0].TargetIdx != 1 {
		t.Errorf("GroupClause[0].TargetIdx: want 1, got %d", q.GroupClause[0].TargetIdx)
	}
}

// TestAnalyze_3_3_Having tests HAVING clause with aggregate condition.
func TestAnalyze_3_3_Having(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT department_id, COUNT(*) AS cnt FROM employees GROUP BY department_id HAVING COUNT(*) > 5")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// HavingQual should be OpExprQ{Op:">", Left:FuncCallExprQ, Right:ConstExprQ}
	if q.HavingQual == nil {
		t.Fatal("HavingQual: want non-nil, got nil")
	}
	opExpr, ok := q.HavingQual.(*OpExprQ)
	if !ok {
		t.Fatalf("HavingQual: want *OpExprQ, got %T", q.HavingQual)
	}
	if opExpr.Op != ">" {
		t.Errorf("OpExprQ.Op: want %q, got %q", ">", opExpr.Op)
	}

	leftFC, ok := opExpr.Left.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("OpExprQ.Left: want *FuncCallExprQ, got %T", opExpr.Left)
	}
	if !leftFC.IsAggregate {
		t.Errorf("Left.IsAggregate: want true, got false")
	}

	rightConst, ok := opExpr.Right.(*ConstExprQ)
	if !ok {
		t.Fatalf("OpExprQ.Right: want *ConstExprQ, got %T", opExpr.Right)
	}
	if rightConst.Value != "5" {
		t.Errorf("Right.Value: want %q, got %q", "5", rightConst.Value)
	}

	// HasAggs should be true
	if !q.HasAggs {
		t.Errorf("HasAggs: want true, got false")
	}
}

// TestAnalyze_3_4_CountDistinct tests COUNT(DISTINCT column).
func TestAnalyze_3_4_CountDistinct(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT COUNT(DISTINCT department_id) FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}

	fc, ok := q.TargetList[0].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *FuncCallExprQ, got %T", q.TargetList[0].Expr)
	}
	if fc.Name != "count" {
		t.Errorf("FuncCallExprQ.Name: want %q, got %q", "count", fc.Name)
	}
	if !fc.IsAggregate {
		t.Errorf("FuncCallExprQ.IsAggregate: want true, got false")
	}
	if !fc.Distinct {
		t.Errorf("FuncCallExprQ.Distinct: want true, got false")
	}

	// Args should contain one VarExprQ with AttNum=4 (department_id)
	if len(fc.Args) != 1 {
		t.Fatalf("FuncCallExprQ.Args: want 1 entry, got %d", len(fc.Args))
	}
	argVar, ok := fc.Args[0].(*VarExprQ)
	if !ok {
		t.Fatalf("Args[0]: want *VarExprQ, got %T", fc.Args[0])
	}
	if argVar.AttNum != 4 {
		t.Errorf("Args[0].AttNum: want 4, got %d", argVar.AttNum)
	}

	// HasAggs should be true
	if !q.HasAggs {
		t.Errorf("HasAggs: want true, got false")
	}
}

// TestAnalyze_3_5_GroupByOrdinal tests GROUP BY with ordinal reference.
func TestAnalyze_3_5_GroupByOrdinal(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT department_id, COUNT(*) FROM employees GROUP BY 1")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// GroupClause: 1 entry referencing TargetIdx=1 (ordinal 1 → first target entry)
	if len(q.GroupClause) != 1 {
		t.Fatalf("GroupClause: want 1 entry, got %d", len(q.GroupClause))
	}
	if q.GroupClause[0].TargetIdx != 1 {
		t.Errorf("GroupClause[0].TargetIdx: want 1, got %d", q.GroupClause[0].TargetIdx)
	}

	// HasAggs should be true
	if !q.HasAggs {
		t.Errorf("HasAggs: want true, got false")
	}
}

// TestAnalyze_5_1_ArithmeticAlias tests SELECT name, salary * 12 AS annual_salary FROM employees.
func TestAnalyze_5_1_ArithmeticAlias(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name, salary * 12 AS annual_salary FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	te := q.TargetList[1]
	if te.ResName != "annual_salary" {
		t.Errorf("ResName: want %q, got %q", "annual_salary", te.ResName)
	}
	// Computed column: no provenance
	if te.ResOrigDB != "" {
		t.Errorf("ResOrigDB: want empty, got %q", te.ResOrigDB)
	}
	if te.ResOrigTable != "" {
		t.Errorf("ResOrigTable: want empty, got %q", te.ResOrigTable)
	}
	if te.ResOrigCol != "" {
		t.Errorf("ResOrigCol: want empty, got %q", te.ResOrigCol)
	}

	opExpr, ok := te.Expr.(*OpExprQ)
	if !ok {
		t.Fatalf("Expr: want *OpExprQ, got %T", te.Expr)
	}
	if opExpr.Op != "*" {
		t.Errorf("OpExprQ.Op: want %q, got %q", "*", opExpr.Op)
	}
	left, ok := opExpr.Left.(*VarExprQ)
	if !ok {
		t.Fatalf("Left: want *VarExprQ, got %T", opExpr.Left)
	}
	if left.AttNum != 5 {
		t.Errorf("Left.AttNum: want 5, got %d", left.AttNum)
	}
	right, ok := opExpr.Right.(*ConstExprQ)
	if !ok {
		t.Fatalf("Right: want *ConstExprQ, got %T", opExpr.Right)
	}
	if right.Value != "12" {
		t.Errorf("Right.Value: want %q, got %q", "12", right.Value)
	}
}

// TestAnalyze_5_2_ConcatFunc tests SELECT CONCAT(name, ' <', email, '>') AS display FROM employees.
func TestAnalyze_5_2_ConcatFunc(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT CONCAT(name, ' <', email, '>') AS display FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}

	te := q.TargetList[0]
	if te.ResName != "display" {
		t.Errorf("ResName: want %q, got %q", "display", te.ResName)
	}

	fc, ok := te.Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("Expr: want *FuncCallExprQ, got %T", te.Expr)
	}
	if fc.Name != "concat" {
		t.Errorf("FuncCallExprQ.Name: want %q, got %q", "concat", fc.Name)
	}
	if fc.IsAggregate {
		t.Errorf("FuncCallExprQ.IsAggregate: want false, got true")
	}
	if len(fc.Args) != 4 {
		t.Fatalf("FuncCallExprQ.Args: want 4 entries, got %d", len(fc.Args))
	}

	// Args[0]: VarExprQ (name)
	if _, ok := fc.Args[0].(*VarExprQ); !ok {
		t.Errorf("Args[0]: want *VarExprQ, got %T", fc.Args[0])
	}
	// Args[1]: ConstExprQ (' <')
	if c1, ok := fc.Args[1].(*ConstExprQ); !ok {
		t.Errorf("Args[1]: want *ConstExprQ, got %T", fc.Args[1])
	} else if c1.Value != " <" {
		t.Errorf("Args[1].Value: want %q, got %q", " <", c1.Value)
	}
	// Args[2]: VarExprQ (email)
	if _, ok := fc.Args[2].(*VarExprQ); !ok {
		t.Errorf("Args[2]: want *VarExprQ, got %T", fc.Args[2])
	}
	// Args[3]: ConstExprQ ('>')
	if c3, ok := fc.Args[3].(*ConstExprQ); !ok {
		t.Errorf("Args[3]: want *ConstExprQ, got %T", fc.Args[3])
	} else if c3.Value != ">" {
		t.Errorf("Args[3].Value: want %q, got %q", ">", c3.Value)
	}
}

// TestAnalyze_5_3_SearchedCase tests searched CASE WHEN with ELSE.
func TestAnalyze_5_3_SearchedCase(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name, CASE WHEN salary > 100000 THEN 'high' WHEN salary > 50000 THEN 'mid' ELSE 'low' END AS salary_band FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	te := q.TargetList[1]
	if te.ResName != "salary_band" {
		t.Errorf("ResName: want %q, got %q", "salary_band", te.ResName)
	}

	caseExpr, ok := te.Expr.(*CaseExprQ)
	if !ok {
		t.Fatalf("Expr: want *CaseExprQ, got %T", te.Expr)
	}
	if caseExpr.TestExpr != nil {
		t.Errorf("TestExpr: want nil (searched CASE), got %T", caseExpr.TestExpr)
	}
	if len(caseExpr.Args) != 2 {
		t.Fatalf("Args: want 2 CaseWhenQs, got %d", len(caseExpr.Args))
	}

	// WHEN salary > 100000 THEN 'high'
	w0 := caseExpr.Args[0]
	if _, ok := w0.Cond.(*OpExprQ); !ok {
		t.Errorf("Args[0].Cond: want *OpExprQ, got %T", w0.Cond)
	}
	if then0, ok := w0.Then.(*ConstExprQ); !ok {
		t.Errorf("Args[0].Then: want *ConstExprQ, got %T", w0.Then)
	} else if then0.Value != "high" {
		t.Errorf("Args[0].Then.Value: want %q, got %q", "high", then0.Value)
	}

	// WHEN salary > 50000 THEN 'mid'
	w1 := caseExpr.Args[1]
	if _, ok := w1.Cond.(*OpExprQ); !ok {
		t.Errorf("Args[1].Cond: want *OpExprQ, got %T", w1.Cond)
	}
	if then1, ok := w1.Then.(*ConstExprQ); !ok {
		t.Errorf("Args[1].Then: want *ConstExprQ, got %T", w1.Then)
	} else if then1.Value != "mid" {
		t.Errorf("Args[1].Then.Value: want %q, got %q", "mid", then1.Value)
	}

	// ELSE 'low'
	if caseExpr.Default == nil {
		t.Fatal("Default: want non-nil, got nil")
	}
	defConst, ok := caseExpr.Default.(*ConstExprQ)
	if !ok {
		t.Fatalf("Default: want *ConstExprQ, got %T", caseExpr.Default)
	}
	if defConst.Value != "low" {
		t.Errorf("Default.Value: want %q, got %q", "low", defConst.Value)
	}
}

// TestAnalyze_5_4_SimpleCase tests simple CASE with operand and no ELSE.
func TestAnalyze_5_4_SimpleCase(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT CASE department_id WHEN 1 THEN 'eng' WHEN 2 THEN 'sales' END FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}

	te := q.TargetList[0]
	caseExpr, ok := te.Expr.(*CaseExprQ)
	if !ok {
		t.Fatalf("Expr: want *CaseExprQ, got %T", te.Expr)
	}

	// TestExpr: VarExprQ for department_id (attnum=4)
	testVar, ok := caseExpr.TestExpr.(*VarExprQ)
	if !ok {
		t.Fatalf("TestExpr: want *VarExprQ, got %T", caseExpr.TestExpr)
	}
	if testVar.AttNum != 4 {
		t.Errorf("TestExpr.AttNum: want 4, got %d", testVar.AttNum)
	}

	if len(caseExpr.Args) != 2 {
		t.Fatalf("Args: want 2 CaseWhenQs, got %d", len(caseExpr.Args))
	}

	// WHEN 1 THEN 'eng'
	if c0, ok := caseExpr.Args[0].Cond.(*ConstExprQ); !ok {
		t.Errorf("Args[0].Cond: want *ConstExprQ, got %T", caseExpr.Args[0].Cond)
	} else if c0.Value != "1" {
		t.Errorf("Args[0].Cond.Value: want %q, got %q", "1", c0.Value)
	}
	if t0, ok := caseExpr.Args[0].Then.(*ConstExprQ); !ok {
		t.Errorf("Args[0].Then: want *ConstExprQ, got %T", caseExpr.Args[0].Then)
	} else if t0.Value != "eng" {
		t.Errorf("Args[0].Then.Value: want %q, got %q", "eng", t0.Value)
	}

	// No ELSE
	if caseExpr.Default != nil {
		t.Errorf("Default: want nil, got %T", caseExpr.Default)
	}
}

// TestAnalyze_5_5_CoalesceIfnull tests COALESCE and IFNULL producing CoalesceExprQ.
func TestAnalyze_5_5_CoalesceIfnull(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT COALESCE(email, 'no-email') AS email, IFNULL(notes, '') AS notes FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	// TargetList[0]: COALESCE(email, 'no-email')
	te0 := q.TargetList[0]
	if te0.ResName != "email" {
		t.Errorf("TargetList[0].ResName: want %q, got %q", "email", te0.ResName)
	}
	coal0, ok := te0.Expr.(*CoalesceExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *CoalesceExprQ, got %T", te0.Expr)
	}
	if len(coal0.Args) != 2 {
		t.Fatalf("coal0.Args: want 2 entries, got %d", len(coal0.Args))
	}
	if v, ok := coal0.Args[0].(*VarExprQ); !ok {
		t.Errorf("coal0.Args[0]: want *VarExprQ, got %T", coal0.Args[0])
	} else if v.AttNum != 3 {
		t.Errorf("coal0.Args[0].AttNum: want 3, got %d", v.AttNum)
	}
	if c0, ok := coal0.Args[1].(*ConstExprQ); !ok {
		t.Errorf("coal0.Args[1]: want *ConstExprQ, got %T", coal0.Args[1])
	} else if c0.Value != "no-email" {
		t.Errorf("coal0.Args[1].Value: want %q, got %q", "no-email", c0.Value)
	}

	// TargetList[1]: IFNULL(notes, '')
	te1 := q.TargetList[1]
	if te1.ResName != "notes" {
		t.Errorf("TargetList[1].ResName: want %q, got %q", "notes", te1.ResName)
	}
	coal1, ok := te1.Expr.(*CoalesceExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *CoalesceExprQ, got %T", te1.Expr)
	}
	if len(coal1.Args) != 2 {
		t.Fatalf("coal1.Args: want 2 entries, got %d", len(coal1.Args))
	}
	if v, ok := coal1.Args[0].(*VarExprQ); !ok {
		t.Errorf("coal1.Args[0]: want *VarExprQ, got %T", coal1.Args[0])
	} else if v.AttNum != 8 {
		t.Errorf("coal1.Args[0].AttNum: want 8, got %d", v.AttNum)
	}
	if c1, ok := coal1.Args[1].(*ConstExprQ); !ok {
		t.Errorf("coal1.Args[1]: want *ConstExprQ, got %T", coal1.Args[1])
	} else if c1.Value != "" {
		t.Errorf("coal1.Args[1].Value: want %q, got %q", "", c1.Value)
	}
}

// TestAnalyze_5_6_CastSigned tests CAST(salary AS SIGNED) producing CastExprQ.
func TestAnalyze_5_6_CastSigned(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT CAST(salary AS SIGNED) AS salary_int FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}

	te := q.TargetList[0]
	if te.ResName != "salary_int" {
		t.Errorf("ResName: want %q, got %q", "salary_int", te.ResName)
	}

	castExpr, ok := te.Expr.(*CastExprQ)
	if !ok {
		t.Fatalf("Expr: want *CastExprQ, got %T", te.Expr)
	}

	// Arg: VarExprQ for salary (attnum=5)
	argVar, ok := castExpr.Arg.(*VarExprQ)
	if !ok {
		t.Fatalf("Arg: want *VarExprQ, got %T", castExpr.Arg)
	}
	if argVar.AttNum != 5 {
		t.Errorf("Arg.AttNum: want 5, got %d", argVar.AttNum)
	}

	// TargetType: SIGNED -> BaseTypeBigInt
	if castExpr.TargetType == nil {
		t.Fatal("TargetType: want non-nil, got nil")
	}
	if castExpr.TargetType.BaseType != BaseTypeBigInt {
		t.Errorf("TargetType.BaseType: want BaseTypeBigInt, got %v", castExpr.TargetType.BaseType)
	}
}

// TestAnalyze_4_1_OrderByDesc tests ORDER BY with DESC on a column in the SELECT list.
func TestAnalyze_4_1_OrderByDesc(t *testing.T) {
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

	sel := parseSelect(t, "SELECT name, salary FROM employees ORDER BY salary DESC")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// TargetList: 2 entries (no junk needed since salary is in SELECT list).
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	// SortClause: 1 entry.
	if len(q.SortClause) != 1 {
		t.Fatalf("SortClause: want 1 entry, got %d", len(q.SortClause))
	}
	sc := q.SortClause[0]
	if sc.TargetIdx != 2 {
		t.Errorf("SortClause[0].TargetIdx: want 2, got %d", sc.TargetIdx)
	}
	if !sc.Descending {
		t.Errorf("SortClause[0].Descending: want true, got false")
	}
	// DESC → NullsFirst=false (MySQL default).
	if sc.NullsFirst {
		t.Errorf("SortClause[0].NullsFirst: want false, got true")
	}
}

// TestAnalyze_4_2_OrderByJunk tests ORDER BY a column NOT in the SELECT list.
func TestAnalyze_4_2_OrderByJunk(t *testing.T) {
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

	sel := parseSelect(t, "SELECT name FROM employees ORDER BY salary")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// TargetList: 2 entries — name (non-junk) + salary (junk).
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	te0 := q.TargetList[0]
	if te0.ResName != "name" {
		t.Errorf("TargetList[0].ResName: want %q, got %q", "name", te0.ResName)
	}
	if te0.ResJunk {
		t.Errorf("TargetList[0].ResJunk: want false, got true")
	}

	te1 := q.TargetList[1]
	if te1.ResName != "salary" {
		t.Errorf("TargetList[1].ResName: want %q, got %q", "salary", te1.ResName)
	}
	if !te1.ResJunk {
		t.Errorf("TargetList[1].ResJunk: want true, got false")
	}

	// SortClause: 1 entry pointing to junk target.
	if len(q.SortClause) != 1 {
		t.Fatalf("SortClause: want 1 entry, got %d", len(q.SortClause))
	}
	sc := q.SortClause[0]
	if sc.TargetIdx != 2 {
		t.Errorf("SortClause[0].TargetIdx: want 2, got %d", sc.TargetIdx)
	}
	// ASC default → NullsFirst=true.
	if !sc.NullsFirst {
		t.Errorf("SortClause[0].NullsFirst: want true, got false")
	}
}

// TestAnalyze_4_3_OrderByLimitOffset tests ORDER BY + LIMIT + OFFSET.
func TestAnalyze_4_3_OrderByLimitOffset(t *testing.T) {
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

	sel := parseSelect(t, "SELECT name FROM employees ORDER BY id LIMIT 10 OFFSET 20")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// LimitCount = ConstExprQ{Value:"10"}
	if q.LimitCount == nil {
		t.Fatal("LimitCount: want non-nil, got nil")
	}
	lc, ok := q.LimitCount.(*ConstExprQ)
	if !ok {
		t.Fatalf("LimitCount: want *ConstExprQ, got %T", q.LimitCount)
	}
	if lc.Value != "10" {
		t.Errorf("LimitCount.Value: want %q, got %q", "10", lc.Value)
	}

	// LimitOffset = ConstExprQ{Value:"20"}
	if q.LimitOffset == nil {
		t.Fatal("LimitOffset: want non-nil, got nil")
	}
	lo, ok := q.LimitOffset.(*ConstExprQ)
	if !ok {
		t.Fatalf("LimitOffset: want *ConstExprQ, got %T", q.LimitOffset)
	}
	if lo.Value != "20" {
		t.Errorf("LimitOffset.Value: want %q, got %q", "20", lo.Value)
	}
}

// TestAnalyze_4_4_LimitComma tests LIMIT offset,count syntax.
func TestAnalyze_4_4_LimitComma(t *testing.T) {
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

	sel := parseSelect(t, "SELECT name FROM employees LIMIT 20, 10")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// Parser normalizes LIMIT 20,10 to Count=10, Offset=20.
	if q.LimitCount == nil {
		t.Fatal("LimitCount: want non-nil, got nil")
	}
	lc, ok := q.LimitCount.(*ConstExprQ)
	if !ok {
		t.Fatalf("LimitCount: want *ConstExprQ, got %T", q.LimitCount)
	}
	if lc.Value != "10" {
		t.Errorf("LimitCount.Value: want %q, got %q", "10", lc.Value)
	}

	if q.LimitOffset == nil {
		t.Fatal("LimitOffset: want non-nil, got nil")
	}
	lo, ok := q.LimitOffset.(*ConstExprQ)
	if !ok {
		t.Fatalf("LimitOffset: want *ConstExprQ, got %T", q.LimitOffset)
	}
	if lo.Value != "20" {
		t.Errorf("LimitOffset.Value: want %q, got %q", "20", lo.Value)
	}
}

// TestAnalyze_4_5_Distinct tests SELECT DISTINCT.
func TestAnalyze_4_5_Distinct(t *testing.T) {
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

	sel := parseSelect(t, "SELECT DISTINCT department_id FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// Distinct should be true.
	if !q.Distinct {
		t.Errorf("Distinct: want true, got false")
	}

	// TargetList: 1 non-junk entry.
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}
	te := q.TargetList[0]
	if te.ResJunk {
		t.Errorf("TargetList[0].ResJunk: want false, got true")
	}
	if te.ResName != "department_id" {
		t.Errorf("TargetList[0].ResName: want %q, got %q", "department_id", te.ResName)
	}
}

// TestAnalyze_6_1_ColumnAliasAmbiguity tests that WHERE resolves against base columns, not SELECT aliases.
func TestAnalyze_6_1_ColumnAliasAmbiguity(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT name AS id FROM employees WHERE id = 1")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// TargetList[0]: name column aliased as "id"
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}
	te := q.TargetList[0]
	if te.ResName != "id" {
		t.Errorf("ResName: want %q, got %q", "id", te.ResName)
	}
	varExpr, ok := te.Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("Expr: want *VarExprQ, got %T", te.Expr)
	}
	if varExpr.AttNum != 2 {
		t.Errorf("TargetList[0] AttNum: want 2 (name column), got %d", varExpr.AttNum)
	}

	// WHERE id = 1: id must resolve to the base column id (AttNum 1), not the alias.
	if q.JoinTree == nil || q.JoinTree.Quals == nil {
		t.Fatal("JoinTree.Quals: want non-nil")
	}
	opExpr, ok := q.JoinTree.Quals.(*OpExprQ)
	if !ok {
		t.Fatalf("Quals: want *OpExprQ, got %T", q.JoinTree.Quals)
	}
	leftVar, ok := opExpr.Left.(*VarExprQ)
	if !ok {
		t.Fatalf("Left: want *VarExprQ, got %T", opExpr.Left)
	}
	if leftVar.AttNum != 1 {
		t.Errorf("WHERE id AttNum: want 1 (base column id), got %d", leftVar.AttNum)
	}
}

// TestAnalyze_6_2_SameColumnTwice tests SELECT referencing the same column twice.
func TestAnalyze_6_2_SameColumnTwice(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT salary, salary + 1000 AS raised FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2 entries, got %d", len(q.TargetList))
	}

	// TargetList[0]: plain salary column
	v0, ok := q.TargetList[0].Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *VarExprQ, got %T", q.TargetList[0].Expr)
	}
	if v0.RangeIdx != 0 {
		t.Errorf("TargetList[0] RangeIdx: want 0, got %d", v0.RangeIdx)
	}
	if v0.AttNum != 5 {
		t.Errorf("TargetList[0] AttNum: want 5, got %d", v0.AttNum)
	}

	// TargetList[1]: salary + 1000
	opExpr, ok := q.TargetList[1].Expr.(*OpExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *OpExprQ, got %T", q.TargetList[1].Expr)
	}
	if opExpr.Op != "+" {
		t.Errorf("OpExprQ.Op: want %q, got %q", "+", opExpr.Op)
	}

	leftVar, ok := opExpr.Left.(*VarExprQ)
	if !ok {
		t.Fatalf("Left: want *VarExprQ, got %T", opExpr.Left)
	}
	if leftVar.RangeIdx != 0 {
		t.Errorf("Left.RangeIdx: want 0, got %d", leftVar.RangeIdx)
	}
	if leftVar.AttNum != 5 {
		t.Errorf("Left.AttNum: want 5, got %d", leftVar.AttNum)
	}

	rightConst, ok := opExpr.Right.(*ConstExprQ)
	if !ok {
		t.Fatalf("Right: want *ConstExprQ, got %T", opExpr.Right)
	}
	if rightConst.Value != "1000" {
		t.Errorf("Right.Value: want %q, got %q", "1000", rightConst.Value)
	}

	// The two VarExprQ references to salary should be distinct pointers but same RangeIdx/AttNum.
	if v0 == leftVar {
		t.Errorf("VarExprQ pointers: want distinct objects, got same pointer")
	}
	if v0.RangeIdx != leftVar.RangeIdx || v0.AttNum != leftVar.AttNum {
		t.Errorf("VarExprQ values: want same RangeIdx/AttNum, got (%d,%d) vs (%d,%d)",
			v0.RangeIdx, v0.AttNum, leftVar.RangeIdx, leftVar.AttNum)
	}
}

// TestAnalyze_6_3_ThreePartQualified tests three-part column qualification: schema.table.column.
func TestAnalyze_6_3_ThreePartQualified(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)

	sel := parseSelect(t, "SELECT testdb.employees.name FROM testdb.employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// RangeTable: 1 entry with DBName="testdb", TableName="employees"
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1 entry, got %d", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.DBName != "testdb" {
		t.Errorf("RTE.DBName: want %q, got %q", "testdb", rte.DBName)
	}
	if rte.TableName != "employees" {
		t.Errorf("RTE.TableName: want %q, got %q", "employees", rte.TableName)
	}

	// TargetList: 1 entry — name resolves to AttNum 2
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1 entry, got %d", len(q.TargetList))
	}
	varExpr, ok := q.TargetList[0].Expr.(*VarExprQ)
	if !ok {
		t.Fatalf("Expr: want *VarExprQ, got %T", q.TargetList[0].Expr)
	}
	if varExpr.RangeIdx != 0 {
		t.Errorf("RangeIdx: want 0, got %d", varExpr.RangeIdx)
	}
	if varExpr.AttNum != 2 {
		t.Errorf("AttNum: want 2, got %d", varExpr.AttNum)
	}
}

// TestAnalyze_6_4_NoFromClause tests SELECT with no FROM clause — pure expressions.
func TestAnalyze_6_4_NoFromClause(t *testing.T) {
	c := wtSetup(t)

	sel := parseSelect(t, "SELECT 1 + 2, 'hello', NOW()")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// RangeTable empty
	if len(q.RangeTable) != 0 {
		t.Errorf("RangeTable: want 0 entries, got %d", len(q.RangeTable))
	}

	// JoinTree.FromList empty
	if q.JoinTree == nil {
		t.Fatal("JoinTree: want non-nil, got nil")
	}
	if len(q.JoinTree.FromList) != 0 {
		t.Errorf("JoinTree.FromList: want 0 entries, got %d", len(q.JoinTree.FromList))
	}

	// TargetList: 3 entries
	if len(q.TargetList) != 3 {
		t.Fatalf("TargetList: want 3 entries, got %d", len(q.TargetList))
	}

	// [0]: 1 + 2 → OpExprQ
	opExpr, ok := q.TargetList[0].Expr.(*OpExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *OpExprQ, got %T", q.TargetList[0].Expr)
	}
	if opExpr.Op != "+" {
		t.Errorf("OpExprQ.Op: want %q, got %q", "+", opExpr.Op)
	}
	leftConst, ok := opExpr.Left.(*ConstExprQ)
	if !ok {
		t.Fatalf("Left: want *ConstExprQ, got %T", opExpr.Left)
	}
	if leftConst.Value != "1" {
		t.Errorf("Left.Value: want %q, got %q", "1", leftConst.Value)
	}
	rightConst, ok := opExpr.Right.(*ConstExprQ)
	if !ok {
		t.Fatalf("Right: want *ConstExprQ, got %T", opExpr.Right)
	}
	if rightConst.Value != "2" {
		t.Errorf("Right.Value: want %q, got %q", "2", rightConst.Value)
	}

	// [1]: 'hello' → ConstExprQ
	strConst, ok := q.TargetList[1].Expr.(*ConstExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *ConstExprQ, got %T", q.TargetList[1].Expr)
	}
	if strConst.Value != "hello" {
		t.Errorf("ConstExprQ.Value: want %q, got %q", "hello", strConst.Value)
	}

	// [2]: NOW() → FuncCallExprQ
	fc, ok := q.TargetList[2].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("TargetList[2].Expr: want *FuncCallExprQ, got %T", q.TargetList[2].Expr)
	}
	if fc.Name != "now" {
		t.Errorf("FuncCallExprQ.Name: want %q, got %q", "now", fc.Name)
	}
	if len(fc.Args) != 0 {
		t.Errorf("FuncCallExprQ.Args: want 0 entries, got %d", len(fc.Args))
	}
	if fc.IsAggregate {
		t.Errorf("FuncCallExprQ.IsAggregate: want false, got true")
	}
}

// ---------------------------------------------------------------------------
// Phase 1b — Batches 7-10: JOINs, USING/NATURAL, FROM subqueries
// ---------------------------------------------------------------------------

const departmentsTableDDL = `CREATE TABLE departments (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	budget DECIMAL(12,2)
)`

const projectsTableDDL = `CREATE TABLE projects (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	department_id INT NOT NULL,
	lead_id INT,
	start_date DATE
)`

// setupJoinTables creates employees, departments, and projects tables.
func setupJoinTables(t *testing.T) *Catalog {
	t.Helper()
	c := wtSetup(t)
	wtExec(t, c, employeesTableDDL)
	wtExec(t, c, departmentsTableDDL)
	wtExec(t, c, projectsTableDDL)
	return c
}

// --- Batch 7: Basic JOINs ---

// TestAnalyze_7_1_InnerJoin tests INNER JOIN with ON condition.
func TestAnalyze_7_1_InnerJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name AS dept_name FROM employees e INNER JOIN departments d ON e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// 3 RTEs: employees, departments, RTEJoin
	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3 entries, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTERelation || q.RangeTable[0].ERef != "e" {
		t.Errorf("RTE[0]: want RTERelation 'e', got kind=%d eref=%q", q.RangeTable[0].Kind, q.RangeTable[0].ERef)
	}
	if q.RangeTable[1].Kind != RTERelation || q.RangeTable[1].ERef != "d" {
		t.Errorf("RTE[1]: want RTERelation 'd', got kind=%d eref=%q", q.RangeTable[1].Kind, q.RangeTable[1].ERef)
	}
	if q.RangeTable[2].Kind != RTEJoin {
		t.Errorf("RTE[2]: want RTEJoin, got kind=%d", q.RangeTable[2].Kind)
	}
	if q.RangeTable[2].JoinType != JoinInner {
		t.Errorf("RTE[2].JoinType: want JoinInner, got %d", q.RangeTable[2].JoinType)
	}

	// JoinTree.FromList should have 1 JoinExprNodeQ.
	if len(q.JoinTree.FromList) != 1 {
		t.Fatalf("FromList: want 1 entry, got %d", len(q.JoinTree.FromList))
	}
	je, ok := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if !ok {
		t.Fatalf("FromList[0]: want *JoinExprNodeQ, got %T", q.JoinTree.FromList[0])
	}
	if je.JoinType != JoinInner {
		t.Errorf("JoinExprNodeQ.JoinType: want JoinInner, got %d", je.JoinType)
	}
	if je.Natural {
		t.Errorf("JoinExprNodeQ.Natural: want false, got true")
	}
	if je.Quals == nil {
		t.Error("JoinExprNodeQ.Quals: want non-nil ON condition, got nil")
	}
	// ON condition should be OpExprQ with "="
	onOp, ok := je.Quals.(*OpExprQ)
	if !ok {
		t.Fatalf("Quals: want *OpExprQ, got %T", je.Quals)
	}
	if onOp.Op != "=" {
		t.Errorf("ON Op: want %q, got %q", "=", onOp.Op)
	}

	// TargetList: 2 entries.
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2, got %d", len(q.TargetList))
	}
	if q.TargetList[0].ResName != "name" {
		t.Errorf("TargetList[0].ResName: want %q, got %q", "name", q.TargetList[0].ResName)
	}
	if q.TargetList[1].ResName != "dept_name" {
		t.Errorf("TargetList[1].ResName: want %q, got %q", "dept_name", q.TargetList[1].ResName)
	}
}

// TestAnalyze_7_2_LeftJoin tests LEFT JOIN.
func TestAnalyze_7_2_LeftJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name AS dept_name FROM employees e LEFT JOIN departments d ON e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}
	if q.RangeTable[2].JoinType != JoinLeft {
		t.Errorf("RTE[2].JoinType: want JoinLeft, got %d", q.RangeTable[2].JoinType)
	}
	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if je.JoinType != JoinLeft {
		t.Errorf("JoinExprNodeQ.JoinType: want JoinLeft, got %d", je.JoinType)
	}
}

// TestAnalyze_7_3_RightJoin tests RIGHT JOIN.
func TestAnalyze_7_3_RightJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name AS dept_name FROM employees e RIGHT JOIN departments d ON e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}
	if q.RangeTable[2].JoinType != JoinRight {
		t.Errorf("RTE[2].JoinType: want JoinRight, got %d", q.RangeTable[2].JoinType)
	}
	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if je.JoinType != JoinRight {
		t.Errorf("JoinExprNodeQ.JoinType: want JoinRight, got %d", je.JoinType)
	}
}

// TestAnalyze_7_4_CrossJoin tests CROSS JOIN with no condition.
func TestAnalyze_7_4_CrossJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name FROM employees e CROSS JOIN departments d`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}
	if q.RangeTable[2].JoinType != JoinCross {
		t.Errorf("RTE[2].JoinType: want JoinCross, got %d", q.RangeTable[2].JoinType)
	}
	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if je.JoinType != JoinCross {
		t.Errorf("JoinExprNodeQ.JoinType: want JoinCross, got %d", je.JoinType)
	}
	if je.Quals != nil {
		t.Errorf("JoinExprNodeQ.Quals: want nil for CROSS JOIN, got %v", je.Quals)
	}
}

// TestAnalyze_7_5_CommaJoin tests comma-separated tables (implicit cross join).
func TestAnalyze_7_5_CommaJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name FROM employees e, departments d WHERE e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// Comma join: 2 RTEs only (no RTEJoin).
	if len(q.RangeTable) != 2 {
		t.Fatalf("RangeTable: want 2, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTERelation {
		t.Errorf("RTE[0]: want RTERelation, got %d", q.RangeTable[0].Kind)
	}
	if q.RangeTable[1].Kind != RTERelation {
		t.Errorf("RTE[1]: want RTERelation, got %d", q.RangeTable[1].Kind)
	}

	// JoinTree.FromList has 2 RangeTableRefQ entries.
	if len(q.JoinTree.FromList) != 2 {
		t.Fatalf("FromList: want 2, got %d", len(q.JoinTree.FromList))
	}
	for i, item := range q.JoinTree.FromList {
		if _, ok := item.(*RangeTableRefQ); !ok {
			t.Errorf("FromList[%d]: want *RangeTableRefQ, got %T", i, item)
		}
	}

	// WHERE in JoinTree.Quals.
	if q.JoinTree.Quals == nil {
		t.Error("JoinTree.Quals: want non-nil WHERE, got nil")
	}
}

func TestAnalyze_7_6_ParenthesizedTableReferenceList(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name FROM (employees e, departments d) WHERE e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3 entries, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTERelation || q.RangeTable[0].ERef != "e" {
		t.Errorf("RTE[0]: want RTERelation 'e', got kind=%d eref=%q", q.RangeTable[0].Kind, q.RangeTable[0].ERef)
	}
	if q.RangeTable[1].Kind != RTERelation || q.RangeTable[1].ERef != "d" {
		t.Errorf("RTE[1]: want RTERelation 'd', got kind=%d eref=%q", q.RangeTable[1].Kind, q.RangeTable[1].ERef)
	}
	if q.RangeTable[2].Kind != RTEJoin || q.RangeTable[2].JoinType != JoinInner {
		t.Errorf("RTE[2]: want inner RTEJoin, got kind=%d join_type=%d", q.RangeTable[2].Kind, q.RangeTable[2].JoinType)
	}
	if len(q.JoinTree.FromList) != 1 {
		t.Fatalf("FromList: want 1 entry, got %d", len(q.JoinTree.FromList))
	}
	je, ok := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if !ok {
		t.Fatalf("FromList[0]: want *JoinExprNodeQ, got %T", q.JoinTree.FromList[0])
	}
	if je.JoinType != JoinInner {
		t.Errorf("JoinExprNodeQ.JoinType: want JoinInner, got %d", je.JoinType)
	}
	if q.JoinTree.Quals == nil {
		t.Error("JoinTree.Quals: want non-nil WHERE, got nil")
	}
}

// --- Batch 8: USING/NATURAL ---

// TestAnalyze_8_1_JoinUsing tests JOIN ... USING (col).
func TestAnalyze_8_1_JoinUsing(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, p.name AS project_name FROM employees e JOIN projects p USING (department_id)`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}

	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if len(je.UsingClause) != 1 || je.UsingClause[0] != "department_id" {
		t.Errorf("UsingClause: want [department_id], got %v", je.UsingClause)
	}
	if je.Natural {
		t.Errorf("Natural: want false, got true")
	}

	// RTEJoin should also have JoinUsing.
	rteJoin := q.RangeTable[2]
	if len(rteJoin.JoinUsing) != 1 || rteJoin.JoinUsing[0] != "department_id" {
		t.Errorf("RTE JoinUsing: want [department_id], got %v", rteJoin.JoinUsing)
	}

	// TargetList: 2 entries.
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2, got %d", len(q.TargetList))
	}
}

// TestAnalyze_8_2_NaturalJoin tests NATURAL JOIN with star expansion.
func TestAnalyze_8_2_NaturalJoin(t *testing.T) {
	c := setupJoinTables(t)
	// employees has: id, name, email, department_id, salary, hire_date, is_active, notes, created_at (9 cols)
	// departments has: id, name, budget (3 cols)
	// Shared columns: id, name → coalesced
	// Result: id, name (from left), email, department_id, salary, hire_date, is_active, notes, created_at, budget
	// = 10 columns
	sel := parseSelect(t, `SELECT * FROM employees e NATURAL JOIN departments d`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if !je.Natural {
		t.Errorf("Natural: want true, got false")
	}
	if je.JoinType != JoinInner {
		t.Errorf("JoinType: want JoinInner, got %d", je.JoinType)
	}

	// Check USING columns are the shared ones: id, name.
	if len(je.UsingClause) != 2 {
		t.Fatalf("UsingClause: want 2, got %d", len(je.UsingClause))
	}

	// Star expansion: should produce 10 columns (9 from employees + 1 remaining from departments).
	// The right-side id and name are coalesced away.
	if len(q.TargetList) != 10 {
		t.Errorf("TargetList: want 10 columns (NATURAL JOIN coalesced), got %d", len(q.TargetList))
	}
}

// TestAnalyze_8_3_NaturalLeftJoin tests NATURAL LEFT JOIN.
func TestAnalyze_8_3_NaturalLeftJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT * FROM employees e NATURAL LEFT JOIN departments d`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if !je.Natural {
		t.Errorf("Natural: want true, got false")
	}
	if je.JoinType != JoinLeft {
		t.Errorf("JoinType: want JoinLeft, got %d", je.JoinType)
	}

	// Same column count as 8.2: 10 columns.
	if len(q.TargetList) != 10 {
		t.Errorf("TargetList: want 10 columns, got %d", len(q.TargetList))
	}
}

// TestAnalyze_8_4_JoinUsingMultiple tests USING with multiple columns.
func TestAnalyze_8_4_JoinUsingMultiple(t *testing.T) {
	c := setupJoinTables(t)
	// employees and departments share: id, name. USING (id, name).
	sel := parseSelect(t, `SELECT * FROM employees e JOIN departments d USING (id, name)`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if len(je.UsingClause) != 2 {
		t.Fatalf("UsingClause: want 2, got %d", len(je.UsingClause))
	}
	if je.UsingClause[0] != "id" || je.UsingClause[1] != "name" {
		t.Errorf("UsingClause: want [id, name], got %v", je.UsingClause)
	}

	// Star expansion: 10 columns (same as NATURAL — same shared columns).
	if len(q.TargetList) != 10 {
		t.Errorf("TargetList: want 10 columns, got %d", len(q.TargetList))
	}
}

// --- Batch 9: FROM subqueries ---

// TestAnalyze_9_1_SimpleSubquery tests FROM (SELECT ...) AS sub.
func TestAnalyze_9_1_SimpleSubquery(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT sub.total FROM (SELECT COUNT(*) AS total FROM employees) AS sub`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// 1 RTE: RTESubquery.
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1, got %d", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Kind != RTESubquery {
		t.Errorf("RTE Kind: want RTESubquery, got %d", rte.Kind)
	}
	if rte.ERef != "sub" {
		t.Errorf("RTE ERef: want %q, got %q", "sub", rte.ERef)
	}
	if len(rte.ColNames) != 1 || rte.ColNames[0] != "total" {
		t.Errorf("ColNames: want [total], got %v", rte.ColNames)
	}

	// Inner query should have HasAggs = true.
	if rte.Subquery == nil {
		t.Fatal("Subquery: want non-nil, got nil")
	}
	if !rte.Subquery.HasAggs {
		t.Errorf("Inner Query HasAggs: want true, got false")
	}

	// Outer TargetList: 1 column.
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1, got %d", len(q.TargetList))
	}
	if q.TargetList[0].ResName != "total" {
		t.Errorf("ResName: want %q, got %q", "total", q.TargetList[0].ResName)
	}
}

// TestAnalyze_9_2_SubqueryWithGroupBy tests FROM subquery with GROUP BY.
func TestAnalyze_9_2_SubqueryWithGroupBy(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT x.dept, x.cnt FROM (SELECT department_id AS dept, COUNT(*) AS cnt FROM employees GROUP BY department_id) AS x`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1, got %d", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Kind != RTESubquery {
		t.Errorf("RTE Kind: want RTESubquery, got %d", rte.Kind)
	}
	if len(rte.ColNames) != 2 {
		t.Fatalf("ColNames: want 2, got %d", len(rte.ColNames))
	}
	if rte.ColNames[0] != "dept" || rte.ColNames[1] != "cnt" {
		t.Errorf("ColNames: want [dept, cnt], got %v", rte.ColNames)
	}

	// Outer TargetList: 2 columns.
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2, got %d", len(q.TargetList))
	}
}

// TestAnalyze_9_3_SubqueryJoinedWithTable tests JOIN between a table and a FROM subquery.
func TestAnalyze_9_3_SubqueryJoinedWithTable(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, sub.avg_sal FROM employees e JOIN (SELECT department_id, AVG(salary) AS avg_sal FROM employees GROUP BY department_id) AS sub ON e.department_id = sub.department_id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// RTEs: employees (0), RTESubquery (1), RTEJoin (2).
	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTERelation {
		t.Errorf("RTE[0] Kind: want RTERelation, got %d", q.RangeTable[0].Kind)
	}
	if q.RangeTable[1].Kind != RTESubquery {
		t.Errorf("RTE[1] Kind: want RTESubquery, got %d", q.RangeTable[1].Kind)
	}
	if q.RangeTable[2].Kind != RTEJoin {
		t.Errorf("RTE[2] Kind: want RTEJoin, got %d", q.RangeTable[2].Kind)
	}

	// Inner subquery has 2 columns.
	subRTE := q.RangeTable[1]
	if len(subRTE.ColNames) != 2 {
		t.Fatalf("Sub ColNames: want 2, got %d", len(subRTE.ColNames))
	}

	// JoinExprNodeQ with ON condition.
	je := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if je.Quals == nil {
		t.Error("ON condition: want non-nil, got nil")
	}

	// Outer TargetList: 2 columns.
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2, got %d", len(q.TargetList))
	}
}

// --- Batch 10: Multi-table edges ---

// TestAnalyze_10_1_ThreeWayJoin tests a three-way JOIN.
func TestAnalyze_10_1_ThreeWayJoin(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT e.name, d.name, p.name FROM employees e INNER JOIN departments d ON e.department_id = d.id INNER JOIN projects p ON d.id = p.department_id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// 5 RTEs: employees(0), departments(1), RTEJoin for first join(2), projects(3), RTEJoin for second join(4).
	if len(q.RangeTable) != 5 {
		t.Fatalf("RangeTable: want 5, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTERelation {
		t.Errorf("RTE[0]: want RTERelation, got %d", q.RangeTable[0].Kind)
	}
	if q.RangeTable[1].Kind != RTERelation {
		t.Errorf("RTE[1]: want RTERelation, got %d", q.RangeTable[1].Kind)
	}
	if q.RangeTable[2].Kind != RTEJoin {
		t.Errorf("RTE[2]: want RTEJoin, got %d", q.RangeTable[2].Kind)
	}
	if q.RangeTable[3].Kind != RTERelation {
		t.Errorf("RTE[3]: want RTERelation, got %d", q.RangeTable[3].Kind)
	}
	if q.RangeTable[4].Kind != RTEJoin {
		t.Errorf("RTE[4]: want RTEJoin, got %d", q.RangeTable[4].Kind)
	}

	// JoinTree.FromList should have 1 outer JoinExprNodeQ.
	if len(q.JoinTree.FromList) != 1 {
		t.Fatalf("FromList: want 1, got %d", len(q.JoinTree.FromList))
	}
	outerJoin, ok := q.JoinTree.FromList[0].(*JoinExprNodeQ)
	if !ok {
		t.Fatalf("FromList[0]: want *JoinExprNodeQ, got %T", q.JoinTree.FromList[0])
	}
	// The left side of the outer join should be another JoinExprNodeQ (the first join).
	innerJoin, ok := outerJoin.Left.(*JoinExprNodeQ)
	if !ok {
		t.Fatalf("OuterJoin.Left: want *JoinExprNodeQ, got %T", outerJoin.Left)
	}
	_ = innerJoin

	// TargetList: 3 columns.
	if len(q.TargetList) != 3 {
		t.Fatalf("TargetList: want 3, got %d", len(q.TargetList))
	}
}

// TestAnalyze_10_2_StarJoin tests SELECT * with a two-table JOIN.
func TestAnalyze_10_2_StarJoin(t *testing.T) {
	c := setupJoinTables(t)
	// employees: 9 cols, departments: 3 cols → 12 total.
	sel := parseSelect(t, `SELECT * FROM employees e JOIN departments d ON e.department_id = d.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// Star expansion for JOIN without USING: all columns from both tables.
	// employees(9) + departments(3) = 12.
	if len(q.TargetList) != 12 {
		t.Errorf("TargetList: want 12 columns, got %d", len(q.TargetList))
	}
}

// TestAnalyze_10_3_AmbiguousColumn tests that unqualified 'name' is ambiguous across two tables.
func TestAnalyze_10_3_AmbiguousColumn(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees e JOIN departments d ON e.department_id = d.id`)
	_, err := c.AnalyzeSelectStmt(sel)
	assertError(t, err, 1052) // ambiguous column
}

// ---------------------------------------------------------------------------
// Phase 1c — Batches 11-13: WHERE subqueries, CTEs, set operations
// ---------------------------------------------------------------------------

// --- Batch 11: WHERE subqueries ---

// TestAnalyze_11_1_ScalarSubqueryInWhere tests WHERE salary > (SELECT AVG(salary) FROM employees).
func TestAnalyze_11_1_ScalarSubqueryInWhere(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees WHERE salary > (SELECT AVG(salary) FROM employees)`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// WHERE should be an OpExprQ with ">" operator.
	op, ok := q.JoinTree.Quals.(*OpExprQ)
	if !ok {
		t.Fatalf("Quals: want *OpExprQ, got %T", q.JoinTree.Quals)
	}
	if op.Op != ">" {
		t.Errorf("Op: want >, got %s", op.Op)
	}

	// Left side: VarExprQ for salary.
	if _, ok := op.Left.(*VarExprQ); !ok {
		t.Errorf("Left: want *VarExprQ, got %T", op.Left)
	}

	// Right side: SubLinkExprQ (scalar subquery).
	subLink, ok := op.Right.(*SubLinkExprQ)
	if !ok {
		t.Fatalf("Right: want *SubLinkExprQ, got %T", op.Right)
	}
	if subLink.Kind != SubLinkScalar {
		t.Errorf("Kind: want SubLinkScalar, got %d", subLink.Kind)
	}
	if subLink.Subquery == nil {
		t.Fatal("Subquery: want non-nil, got nil")
	}
	if !subLink.Subquery.HasAggs {
		t.Errorf("Subquery.HasAggs: want true, got false")
	}
}

// TestAnalyze_11_2_InSubquery tests WHERE department_id IN (SELECT id FROM departments WHERE budget > 100000).
func TestAnalyze_11_2_InSubquery(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees WHERE department_id IN (SELECT id FROM departments WHERE budget > 100000)`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// WHERE should be a SubLinkExprQ with Kind=SubLinkIn.
	subLink, ok := q.JoinTree.Quals.(*SubLinkExprQ)
	if !ok {
		t.Fatalf("Quals: want *SubLinkExprQ, got %T", q.JoinTree.Quals)
	}
	if subLink.Kind != SubLinkIn {
		t.Errorf("Kind: want SubLinkIn, got %d", subLink.Kind)
	}
	if subLink.Op != "=" {
		t.Errorf("Op: want =, got %s", subLink.Op)
	}

	// TestExpr: VarExprQ for department_id (AttNum=4).
	testVar, ok := subLink.TestExpr.(*VarExprQ)
	if !ok {
		t.Fatalf("TestExpr: want *VarExprQ, got %T", subLink.TestExpr)
	}
	if testVar.AttNum != 4 {
		t.Errorf("TestExpr.AttNum: want 4, got %d", testVar.AttNum)
	}

	// Subquery should have a WHERE qual.
	if subLink.Subquery == nil {
		t.Fatal("Subquery: want non-nil, got nil")
	}
	if subLink.Subquery.JoinTree.Quals == nil {
		t.Error("Subquery WHERE: want non-nil, got nil")
	}
}

// TestAnalyze_11_3_ExistsCorrelated tests EXISTS with a correlated subquery.
func TestAnalyze_11_3_ExistsCorrelated(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees e WHERE EXISTS (SELECT 1 FROM projects p WHERE p.lead_id = e.id)`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// WHERE should be SubLinkExprQ with Kind=SubLinkExists.
	subLink, ok := q.JoinTree.Quals.(*SubLinkExprQ)
	if !ok {
		t.Fatalf("Quals: want *SubLinkExprQ, got %T", q.JoinTree.Quals)
	}
	if subLink.Kind != SubLinkExists {
		t.Errorf("Kind: want SubLinkExists, got %d", subLink.Kind)
	}
	if subLink.TestExpr != nil {
		t.Errorf("TestExpr: want nil for EXISTS, got %v", subLink.TestExpr)
	}

	// Inner query WHERE should reference outer scope.
	innerQ := subLink.Subquery
	if innerQ == nil {
		t.Fatal("Subquery: want non-nil, got nil")
	}
	innerQuals := innerQ.JoinTree.Quals
	if innerQuals == nil {
		t.Fatal("inner Quals: want non-nil, got nil")
	}

	// Should be OpExprQ: p.lead_id = e.id
	innerOp, ok := innerQuals.(*OpExprQ)
	if !ok {
		t.Fatalf("inner Quals: want *OpExprQ, got %T", innerQuals)
	}

	// One side should have LevelsUp=1 (correlated reference to outer e.id).
	leftVar, leftOk := innerOp.Left.(*VarExprQ)
	rightVar, rightOk := innerOp.Right.(*VarExprQ)
	if !leftOk || !rightOk {
		t.Fatalf("inner Op sides: want both *VarExprQ, got %T and %T", innerOp.Left, innerOp.Right)
	}

	// e.id is from the outer scope (LevelsUp=1); p.lead_id is from inner scope (LevelsUp=0).
	if leftVar.LevelsUp == 0 && rightVar.LevelsUp == 1 {
		// right is correlated
		if rightVar.AttNum != 1 { // e.id is column 1
			t.Errorf("correlated ref AttNum: want 1 (id), got %d", rightVar.AttNum)
		}
	} else if leftVar.LevelsUp == 1 && rightVar.LevelsUp == 0 {
		// left is correlated
		if leftVar.AttNum != 1 {
			t.Errorf("correlated ref AttNum: want 1 (id), got %d", leftVar.AttNum)
		}
	} else {
		t.Errorf("expected one LevelsUp=1 (correlated), got left=%d right=%d", leftVar.LevelsUp, rightVar.LevelsUp)
	}
}

// TestAnalyze_11_4_NotInSubquery tests NOT IN (SELECT ...).
func TestAnalyze_11_4_NotInSubquery(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees WHERE department_id NOT IN (SELECT department_id FROM projects)`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// NOT IN subquery is represented as BoolExprQ{BoolNot, [SubLinkExprQ{SubLinkIn}]}.
	boolExpr, ok := q.JoinTree.Quals.(*BoolExprQ)
	if !ok {
		t.Fatalf("Quals: want *BoolExprQ, got %T", q.JoinTree.Quals)
	}
	if boolExpr.Op != BoolNot {
		t.Errorf("BoolOp: want BoolNot, got %d", boolExpr.Op)
	}
	if len(boolExpr.Args) != 1 {
		t.Fatalf("BoolExpr.Args: want 1, got %d", len(boolExpr.Args))
	}

	subLink, ok := boolExpr.Args[0].(*SubLinkExprQ)
	if !ok {
		t.Fatalf("BoolExpr.Args[0]: want *SubLinkExprQ, got %T", boolExpr.Args[0])
	}
	if subLink.Kind != SubLinkIn {
		t.Errorf("Kind: want SubLinkIn, got %d", subLink.Kind)
	}
}

// TestAnalyze_11_5_ScalarSubqueryInSelect tests a scalar subquery in the SELECT list.
func TestAnalyze_11_5_ScalarSubqueryInSelect(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name, (SELECT COUNT(*) FROM projects p WHERE p.lead_id = e.id) AS project_count FROM employees e`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2, got %d", len(q.TargetList))
	}

	// Second column should be a SubLinkExprQ.
	subLink, ok := q.TargetList[1].Expr.(*SubLinkExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *SubLinkExprQ, got %T", q.TargetList[1].Expr)
	}
	if subLink.Kind != SubLinkScalar {
		t.Errorf("Kind: want SubLinkScalar, got %d", subLink.Kind)
	}
	if q.TargetList[1].ResName != "project_count" {
		t.Errorf("ResName: want project_count, got %s", q.TargetList[1].ResName)
	}

	// Inner query should have HasAggs=true (COUNT).
	if !subLink.Subquery.HasAggs {
		t.Errorf("Subquery.HasAggs: want true, got false")
	}
}

// --- Batch 12: CTEs ---

// TestAnalyze_12_1_SimpleCTE tests a simple CTE with WITH clause.
func TestAnalyze_12_1_SimpleCTE(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `WITH dept_stats AS (SELECT department_id, COUNT(*) AS cnt FROM employees GROUP BY department_id) SELECT * FROM dept_stats`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// CTEList should have 1 entry.
	if len(q.CTEList) != 1 {
		t.Fatalf("CTEList: want 1, got %d", len(q.CTEList))
	}
	cte := q.CTEList[0]
	if cte.Name != "dept_stats" {
		t.Errorf("CTE Name: want dept_stats, got %s", cte.Name)
	}
	if cte.Recursive {
		t.Errorf("CTE Recursive: want false, got true")
	}
	if cte.Query == nil {
		t.Fatal("CTE Query: want non-nil, got nil")
	}

	// RangeTable should have an RTECTE entry.
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1, got %d", len(q.RangeTable))
	}
	rte := q.RangeTable[0]
	if rte.Kind != RTECTE {
		t.Errorf("RTE Kind: want RTECTE, got %d", rte.Kind)
	}
	if rte.CTEName != "dept_stats" {
		t.Errorf("RTE CTEName: want dept_stats, got %s", rte.CTEName)
	}
	if rte.CTEIndex != 0 {
		t.Errorf("RTE CTEIndex: want 0, got %d", rte.CTEIndex)
	}

	// Star expansion from CTE should produce columns from the CTE body.
	if len(q.TargetList) != 2 {
		t.Errorf("TargetList: want 2, got %d", len(q.TargetList))
	}
}

// TestAnalyze_12_2_MultipleCTEs tests multiple CTEs.
func TestAnalyze_12_2_MultipleCTEs(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `
		WITH
			active_emps AS (SELECT id, name FROM employees WHERE is_active = 1),
			big_depts AS (SELECT id, name FROM departments WHERE budget > 100000)
		SELECT a.name, b.name AS dept_name FROM active_emps a JOIN big_depts b ON a.id = b.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.CTEList) != 2 {
		t.Fatalf("CTEList: want 2, got %d", len(q.CTEList))
	}
	if q.CTEList[0].Name != "active_emps" {
		t.Errorf("CTE[0].Name: want active_emps, got %s", q.CTEList[0].Name)
	}
	if q.CTEList[1].Name != "big_depts" {
		t.Errorf("CTE[1].Name: want big_depts, got %s", q.CTEList[1].Name)
	}

	// RangeTable: 2 RTECTE + 1 RTEJoin = 3.
	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTECTE {
		t.Errorf("RTE[0]: want RTECTE, got %d", q.RangeTable[0].Kind)
	}
	if q.RangeTable[1].Kind != RTECTE {
		t.Errorf("RTE[1]: want RTECTE, got %d", q.RangeTable[1].Kind)
	}

	if len(q.TargetList) != 2 {
		t.Errorf("TargetList: want 2, got %d", len(q.TargetList))
	}
}

// TestAnalyze_12_3_CTEWithExplicitColumns tests CTE with explicit column aliases.
func TestAnalyze_12_3_CTEWithExplicitColumns(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `WITH emp_summary(dept, cnt) AS (SELECT department_id, COUNT(*) FROM employees GROUP BY department_id) SELECT dept, cnt FROM emp_summary`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.CTEList) != 1 {
		t.Fatalf("CTEList: want 1, got %d", len(q.CTEList))
	}
	cte := q.CTEList[0]
	if len(cte.ColumnNames) != 2 || cte.ColumnNames[0] != "dept" || cte.ColumnNames[1] != "cnt" {
		t.Errorf("CTE ColumnNames: want [dept, cnt], got %v", cte.ColumnNames)
	}

	// RTECTE should use the explicit column names.
	rte := q.RangeTable[0]
	if rte.Kind != RTECTE {
		t.Fatalf("RTE Kind: want RTECTE, got %d", rte.Kind)
	}
	if len(rte.ColNames) != 2 || rte.ColNames[0] != "dept" || rte.ColNames[1] != "cnt" {
		t.Errorf("RTE ColNames: want [dept, cnt], got %v", rte.ColNames)
	}

	// TargetList should resolve dept and cnt.
	if len(q.TargetList) != 2 {
		t.Fatalf("TargetList: want 2, got %d", len(q.TargetList))
	}
	if q.TargetList[0].ResName != "dept" {
		t.Errorf("TargetList[0].ResName: want dept, got %s", q.TargetList[0].ResName)
	}
	if q.TargetList[1].ResName != "cnt" {
		t.Errorf("TargetList[1].ResName: want cnt, got %s", q.TargetList[1].ResName)
	}
}

// TestAnalyze_12_4_CTEReferencedTwice tests a CTE referenced twice in the FROM clause.
func TestAnalyze_12_4_CTEReferencedTwice(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `WITH emp_ids AS (SELECT id, name FROM employees) SELECT a.name, b.name FROM emp_ids a JOIN emp_ids b ON a.id = b.id`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.CTEList) != 1 {
		t.Fatalf("CTEList: want 1, got %d", len(q.CTEList))
	}

	// Two RTECTE entries (one per reference), plus one RTEJoin.
	if len(q.RangeTable) != 3 {
		t.Fatalf("RangeTable: want 3, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTECTE {
		t.Errorf("RTE[0]: want RTECTE, got %d", q.RangeTable[0].Kind)
	}
	if q.RangeTable[1].Kind != RTECTE {
		t.Errorf("RTE[1]: want RTECTE, got %d", q.RangeTable[1].Kind)
	}
	// Both should reference CTEIndex=0.
	if q.RangeTable[0].CTEIndex != 0 {
		t.Errorf("RTE[0].CTEIndex: want 0, got %d", q.RangeTable[0].CTEIndex)
	}
	if q.RangeTable[1].CTEIndex != 0 {
		t.Errorf("RTE[1].CTEIndex: want 0, got %d", q.RangeTable[1].CTEIndex)
	}
}

// TestAnalyze_12_5_RecursiveCTE tests WITH RECURSIVE.
func TestAnalyze_12_5_RecursiveCTE(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE categories (id INT PRIMARY KEY, name VARCHAR(100), parent_id INT)`)

	sel := parseSelect(t, `
		WITH RECURSIVE cat_tree(id, name, parent_id) AS (
			SELECT id, name, parent_id FROM categories WHERE parent_id IS NULL
			UNION ALL
			SELECT c.id, c.name, c.parent_id FROM categories c INNER JOIN cat_tree ct ON c.parent_id = ct.id
		)
		SELECT * FROM cat_tree`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if !q.IsRecursive {
		t.Errorf("IsRecursive: want true, got false")
	}
	if len(q.CTEList) != 1 {
		t.Fatalf("CTEList: want 1, got %d", len(q.CTEList))
	}
	cte := q.CTEList[0]
	if !cte.Recursive {
		t.Errorf("CTE Recursive: want true, got false")
	}

	// CTE body should be a set-op query.
	if cte.Query.SetOp != SetOpUnion {
		t.Errorf("CTE SetOp: want SetOpUnion, got %d", cte.Query.SetOp)
	}
	if !cte.Query.AllSetOp {
		t.Errorf("CTE AllSetOp: want true, got false")
	}

	// RTECTE in main query.
	if len(q.RangeTable) != 1 {
		t.Fatalf("RangeTable: want 1, got %d", len(q.RangeTable))
	}
	if q.RangeTable[0].Kind != RTECTE {
		t.Errorf("RTE Kind: want RTECTE, got %d", q.RangeTable[0].Kind)
	}

	// Star expansion: 3 columns from the CTE.
	if len(q.TargetList) != 3 {
		t.Errorf("TargetList: want 3, got %d", len(q.TargetList))
	}
}

// --- Batch 13: Set operations ---

// TestAnalyze_13_1_Union tests UNION (without ALL).
func TestAnalyze_13_1_Union(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees UNION SELECT name FROM departments`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.SetOp != SetOpUnion {
		t.Errorf("SetOp: want SetOpUnion, got %d", q.SetOp)
	}
	if q.AllSetOp {
		t.Errorf("AllSetOp: want false, got true")
	}
	if q.LArg == nil {
		t.Fatal("LArg: want non-nil, got nil")
	}
	if q.RArg == nil {
		t.Fatal("RArg: want non-nil, got nil")
	}

	// Result columns from left arm.
	if len(q.TargetList) != 1 {
		t.Fatalf("TargetList: want 1, got %d", len(q.TargetList))
	}
	if q.TargetList[0].ResName != "name" {
		t.Errorf("TargetList[0].ResName: want name, got %s", q.TargetList[0].ResName)
	}
}

// TestAnalyze_13_2_UnionAll tests UNION ALL.
func TestAnalyze_13_2_UnionAll(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees UNION ALL SELECT name FROM departments`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.SetOp != SetOpUnion {
		t.Errorf("SetOp: want SetOpUnion, got %d", q.SetOp)
	}
	if !q.AllSetOp {
		t.Errorf("AllSetOp: want true, got false")
	}
}

// TestAnalyze_13_3_Intersect tests INTERSECT.
func TestAnalyze_13_3_Intersect(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees INTERSECT SELECT name FROM departments`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.SetOp != SetOpIntersect {
		t.Errorf("SetOp: want SetOpIntersect, got %d", q.SetOp)
	}
	if q.AllSetOp {
		t.Errorf("AllSetOp: want false, got true")
	}
	if q.LArg == nil || q.RArg == nil {
		t.Fatal("LArg/RArg: want non-nil")
	}
}

// TestAnalyze_13_4_Except tests EXCEPT.
func TestAnalyze_13_4_Except(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees EXCEPT SELECT name FROM departments`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.SetOp != SetOpExcept {
		t.Errorf("SetOp: want SetOpExcept, got %d", q.SetOp)
	}
	if q.AllSetOp {
		t.Errorf("AllSetOp: want false, got true")
	}
}

// TestAnalyze_13_5_UnionAllOrderByLimit tests UNION ALL with ORDER BY and LIMIT.
func TestAnalyze_13_5_UnionAllOrderByLimit(t *testing.T) {
	c := setupJoinTables(t)
	sel := parseSelect(t, `SELECT name FROM employees UNION ALL SELECT name FROM departments ORDER BY name LIMIT 10`)
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if q.SetOp != SetOpUnion {
		t.Errorf("SetOp: want SetOpUnion, got %d", q.SetOp)
	}
	if !q.AllSetOp {
		t.Errorf("AllSetOp: want true, got false")
	}

	// ORDER BY should be populated.
	if len(q.SortClause) != 1 {
		t.Fatalf("SortClause: want 1, got %d", len(q.SortClause))
	}

	// LIMIT should be populated.
	if q.LimitCount == nil {
		t.Error("LimitCount: want non-nil, got nil")
	}
	constExpr, ok := q.LimitCount.(*ConstExprQ)
	if !ok {
		t.Fatalf("LimitCount: want *ConstExprQ, got %T", q.LimitCount)
	}
	if constExpr.Value != "10" {
		t.Errorf("LimitCount value: want 10, got %s", constExpr.Value)
	}
}

// ---------- Phase 3: Standalone expression analysis, function types, DDL hookups ----------

// TestAnalyze_15_1_CheckConstraintAnalyzed tests that CHECK constraint expressions
// are analyzed and stored as CheckAnalyzed on the Constraint.
func TestAnalyze_15_1_CheckConstraintAnalyzed(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (a INT, b INT, CONSTRAINT chk CHECK (a > 0 AND b > 0))`)

	db := c.GetDatabase("testdb")
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}

	// Find the CHECK constraint named "chk".
	var con *Constraint
	for _, cc := range tbl.Constraints {
		if cc.Type == ConCheck && cc.Name == "chk" {
			con = cc
			break
		}
	}
	if con == nil {
		t.Fatal("CHECK constraint 'chk' not found")
	}
	if con.CheckAnalyzed == nil {
		t.Fatal("CheckAnalyzed: want non-nil, got nil")
	}

	// Top-level should be BoolExprQ with BoolAnd.
	boolExpr, ok := con.CheckAnalyzed.(*BoolExprQ)
	if !ok {
		t.Fatalf("CheckAnalyzed: want *BoolExprQ, got %T", con.CheckAnalyzed)
	}
	if boolExpr.Op != BoolAnd {
		t.Errorf("BoolExprQ.Op: want BoolAnd, got %v", boolExpr.Op)
	}
	if len(boolExpr.Args) != 2 {
		t.Fatalf("BoolExprQ.Args: want 2, got %d", len(boolExpr.Args))
	}

	// Each arg should be OpExprQ with Op ">".
	for i, arg := range boolExpr.Args {
		opExpr, ok := arg.(*OpExprQ)
		if !ok {
			t.Fatalf("Args[%d]: want *OpExprQ, got %T", i, arg)
		}
		if opExpr.Op != ">" {
			t.Errorf("Args[%d].Op: want >, got %s", i, opExpr.Op)
		}
	}
}

// TestAnalyze_15_2_DefaultAnalyzed tests that DEFAULT expressions are analyzed.
func TestAnalyze_15_2_DefaultAnalyzed(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (a INT DEFAULT 42, b VARCHAR(100) DEFAULT 'hello')`)

	db := c.GetDatabase("testdb")
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}

	// Column "a" default: ConstExprQ{Value:"42"}
	colA := tbl.GetColumn("a")
	if colA == nil {
		t.Fatal("column a not found")
	}
	if colA.DefaultAnalyzed == nil {
		t.Fatal("a.DefaultAnalyzed: want non-nil, got nil")
	}
	constA, ok := colA.DefaultAnalyzed.(*ConstExprQ)
	if !ok {
		t.Fatalf("a.DefaultAnalyzed: want *ConstExprQ, got %T", colA.DefaultAnalyzed)
	}
	if constA.Value != "42" {
		t.Errorf("a.DefaultAnalyzed.Value: want 42, got %s", constA.Value)
	}

	// Column "b" default: ConstExprQ{Value:"hello"}
	colB := tbl.GetColumn("b")
	if colB == nil {
		t.Fatal("column b not found")
	}
	if colB.DefaultAnalyzed == nil {
		t.Fatal("b.DefaultAnalyzed: want non-nil, got nil")
	}
	constB, ok := colB.DefaultAnalyzed.(*ConstExprQ)
	if !ok {
		t.Fatalf("b.DefaultAnalyzed: want *ConstExprQ, got %T", colB.DefaultAnalyzed)
	}
	if constB.Value != "hello" {
		t.Errorf("b.DefaultAnalyzed.Value: want hello, got %s", constB.Value)
	}
}

// TestAnalyze_15_3_GeneratedAnalyzed tests that GENERATED ALWAYS AS expressions are analyzed.
func TestAnalyze_15_3_GeneratedAnalyzed(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) STORED)`)

	db := c.GetDatabase("testdb")
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("table t not found")
	}

	colC := tbl.GetColumn("c")
	if colC == nil {
		t.Fatal("column c not found")
	}
	if colC.GeneratedAnalyzed == nil {
		t.Fatal("c.GeneratedAnalyzed: want non-nil, got nil")
	}

	opExpr, ok := colC.GeneratedAnalyzed.(*OpExprQ)
	if !ok {
		t.Fatalf("c.GeneratedAnalyzed: want *OpExprQ, got %T", colC.GeneratedAnalyzed)
	}
	if opExpr.Op != "+" {
		t.Errorf("OpExprQ.Op: want +, got %s", opExpr.Op)
	}

	// Left should be VarExprQ for column a (AttNum=1).
	leftVar, ok := opExpr.Left.(*VarExprQ)
	if !ok {
		t.Fatalf("Left: want *VarExprQ, got %T", opExpr.Left)
	}
	if leftVar.AttNum != 1 {
		t.Errorf("Left.AttNum: want 1, got %d", leftVar.AttNum)
	}

	// Right should be VarExprQ for column b (AttNum=2).
	rightVar, ok := opExpr.Right.(*VarExprQ)
	if !ok {
		t.Fatalf("Right: want *VarExprQ, got %T", opExpr.Right)
	}
	if rightVar.AttNum != 2 {
		t.Errorf("Right.AttNum: want 2, got %d", rightVar.AttNum)
	}
}

// TestAnalyze_15_4_FunctionReturnTypes tests that function return types are populated.
func TestAnalyze_15_4_FunctionReturnTypes(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE employees (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	)`)

	sel := parseSelect(t, "SELECT COUNT(*), CONCAT(name, '!'), NOW() FROM employees")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	if len(q.TargetList) != 3 {
		t.Fatalf("TargetList: want 3, got %d", len(q.TargetList))
	}

	// COUNT(*) -> BaseTypeBigInt
	fc0, ok := q.TargetList[0].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("TargetList[0].Expr: want *FuncCallExprQ, got %T", q.TargetList[0].Expr)
	}
	if fc0.ResultType == nil {
		t.Fatal("COUNT(*) ResultType: want non-nil, got nil")
	}
	if fc0.ResultType.BaseType != BaseTypeBigInt {
		t.Errorf("COUNT(*) ResultType.BaseType: want BaseTypeBigInt, got %d", fc0.ResultType.BaseType)
	}

	// CONCAT(name, '!') -> BaseTypeVarchar
	fc1, ok := q.TargetList[1].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("TargetList[1].Expr: want *FuncCallExprQ, got %T", q.TargetList[1].Expr)
	}
	if fc1.ResultType == nil {
		t.Fatal("CONCAT() ResultType: want non-nil, got nil")
	}
	if fc1.ResultType.BaseType != BaseTypeVarchar {
		t.Errorf("CONCAT() ResultType.BaseType: want BaseTypeVarchar, got %d", fc1.ResultType.BaseType)
	}

	// NOW() -> BaseTypeDateTime
	fc2, ok := q.TargetList[2].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("TargetList[2].Expr: want *FuncCallExprQ, got %T", q.TargetList[2].Expr)
	}
	if fc2.ResultType == nil {
		t.Fatal("NOW() ResultType: want non-nil, got nil")
	}
	if fc2.ResultType.BaseType != BaseTypeDateTime {
		t.Errorf("NOW() ResultType.BaseType: want BaseTypeDateTime, got %d", fc2.ResultType.BaseType)
	}
}

// TestAnalyze_15_5_ViewFunctionTypeFlow tests that function return types flow through views.
func TestAnalyze_15_5_ViewFunctionTypeFlow(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE t (id INT, name VARCHAR(100))`)
	wtExec(t, c, `CREATE VIEW v AS SELECT id, name, COUNT(*) OVER () AS cnt FROM t`)

	sel := parseSelect(t, "SELECT * FROM v")
	q, err := c.AnalyzeSelectStmt(sel)
	assertNoError(t, err)

	// Should have 3 columns: id, name, cnt
	if len(q.TargetList) != 3 {
		t.Fatalf("TargetList: want 3, got %d", len(q.TargetList))
	}
	if q.TargetList[0].ResName != "id" {
		t.Errorf("TargetList[0].ResName: want id, got %s", q.TargetList[0].ResName)
	}
	if q.TargetList[1].ResName != "name" {
		t.Errorf("TargetList[1].ResName: want name, got %s", q.TargetList[1].ResName)
	}
	if q.TargetList[2].ResName != "cnt" {
		t.Errorf("TargetList[2].ResName: want cnt, got %s", q.TargetList[2].ResName)
	}
}

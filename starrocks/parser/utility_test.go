package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

func TestDescribeTable(t *testing.T) {
	file, errs := Parse("DESC test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.DescribeStmt)
	if !ok {
		t.Fatalf("expected *ast.DescribeStmt, got %T", file.Stmts[0])
	}
	if n.Target == nil || n.Target.Parts[0] != "test_table" {
		t.Errorf("Target = %v, want test_table", n.Target)
	}
	if n.Full {
		t.Error("Full should be false")
	}
}

func TestDescribeKeyword(t *testing.T) {
	file, errs := Parse("DESCRIBE test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if n.Target == nil || n.Target.Parts[0] != "test_table" {
		t.Errorf("Target = %v, want test_table", n.Target)
	}
}

func TestDescribeFullTable(t *testing.T) {
	file, errs := Parse("DESCRIBE FULL test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if !n.Full {
		t.Error("Full should be true")
	}
}

func TestDescribeAllVerbose(t *testing.T) {
	file, errs := Parse("DESCRIBE test_table ALL VERBOSE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if !n.AllVerbose {
		t.Error("AllVerbose should be true")
	}
}

func TestDescribeLocValid(t *testing.T) {
	file, errs := Parse("DESC t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if !n.Loc.IsValid() {
		t.Error("DescribeStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

func TestExplainSelect(t *testing.T) {
	file, errs := Parse("EXPLAIN SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.ExplainStmt)
	if !ok {
		t.Fatalf("expected *ast.ExplainStmt, got %T", file.Stmts[0])
	}
	if n.Type != "" {
		t.Errorf("Type = %q, want empty", n.Type)
	}
	if n.Query == nil {
		t.Error("Query should not be nil")
	}
}

func TestExplainVerbose(t *testing.T) {
	file, errs := Parse("EXPLAIN VERBOSE SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "VERBOSE" {
		t.Errorf("Type = %q, want VERBOSE", n.Type)
	}
}

func TestExplainGraph(t *testing.T) {
	file, errs := Parse("EXPLAIN GRAPH SELECT 1 FROM t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "GRAPH" {
		t.Errorf("Type = %q, want GRAPH", n.Type)
	}
}

func TestExplainPlan(t *testing.T) {
	file, errs := Parse("EXPLAIN PLAN SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "PLAN" {
		t.Errorf("Type = %q, want PLAN", n.Type)
	}
}

func TestExplainShape(t *testing.T) {
	file, errs := Parse("EXPLAIN SHAPE SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "SHAPE" {
		t.Errorf("Type = %q, want SHAPE", n.Type)
	}
}

func TestExplainMemo(t *testing.T) {
	file, errs := Parse("EXPLAIN MEMO SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "MEMO" {
		t.Errorf("Type = %q, want MEMO", n.Type)
	}
}

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

func TestUseDatabase(t *testing.T) {
	file, errs := Parse("USE demo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.UseStmt)
	if !ok {
		t.Fatalf("expected *ast.UseStmt, got %T", file.Stmts[0])
	}
	if n.Database != "demo" {
		t.Errorf("Database = %q, want demo", n.Database)
	}
	if n.Catalog != "" {
		t.Errorf("Catalog = %q, want empty", n.Catalog)
	}
}

func TestUseCatalogDB(t *testing.T) {
	file, errs := Parse("USE hms_catalog.demo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UseStmt)
	if n.Catalog != "hms_catalog" {
		t.Errorf("Catalog = %q, want hms_catalog", n.Catalog)
	}
	if n.Database != "demo" {
		t.Errorf("Database = %q, want demo", n.Database)
	}
}

func TestUseDBAtCluster(t *testing.T) {
	file, errs := Parse("USE mydb@mycluster")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UseStmt)
	if n.Database != "mydb" {
		t.Errorf("Database = %q, want mydb", n.Database)
	}
	if n.Cluster != "mycluster" {
		t.Errorf("Cluster = %q, want mycluster", n.Cluster)
	}
}

func TestUseLocValid(t *testing.T) {
	file, errs := Parse("USE demo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UseStmt)
	if !n.Loc.IsValid() {
		t.Error("UseStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// SET (generic variable assignment)
// ---------------------------------------------------------------------------

func TestSetVariable(t *testing.T) {
	file, errs := Parse(`SET time_zone = "Asia/Shanghai"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.SetStmt)
	if !ok {
		t.Fatalf("expected *ast.SetStmt, got %T", file.Stmts[0])
	}
	if n.Type != "VARIABLE" {
		t.Errorf("Type = %q, want VARIABLE", n.Type)
	}
	if len(n.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(n.Items))
	}
	if n.Items[0].Name != "time_zone" {
		t.Errorf("Name = %q, want time_zone", n.Items[0].Name)
	}
}

func TestSetGlobalVariable(t *testing.T) {
	file, errs := Parse("SET GLOBAL exec_mem_limit = 137438953472")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if len(n.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(n.Items))
	}
	item := n.Items[0]
	if item.Scope != "GLOBAL" {
		t.Errorf("Scope = %q, want GLOBAL", item.Scope)
	}
	if item.Name != "exec_mem_limit" {
		t.Errorf("Name = %q, want exec_mem_limit", item.Name)
	}
}

func TestSetDoubleAtVariable(t *testing.T) {
	file, errs := Parse("SET @@exec_mem_limit = 137438953472")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if len(n.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(n.Items))
	}
	item := n.Items[0]
	if item.Name != "exec_mem_limit" {
		t.Errorf("Name = %q, want exec_mem_limit", item.Name)
	}
}

func TestSetNames(t *testing.T) {
	file, errs := Parse("SET NAMES 'utf8'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if n.Type != "NAMES" {
		t.Errorf("Type = %q, want NAMES", n.Type)
	}
	if len(n.Items) != 1 || n.Items[0].Raw != "utf8" {
		t.Errorf("Items[0].Raw = %q, want utf8", n.Items[0].Raw)
	}
}

func TestSetCharset(t *testing.T) {
	file, errs := Parse("SET CHARSET 'utf8mb4'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if n.Type != "CHARSET" {
		t.Errorf("Type = %q, want CHARSET", n.Type)
	}
}

func TestSetTransaction(t *testing.T) {
	file, errs := Parse("SET TRANSACTION READ ONLY")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if n.Type != "TRANSACTION" {
		t.Errorf("Type = %q, want TRANSACTION", n.Type)
	}
}

func TestSetLocValid(t *testing.T) {
	file, errs := Parse("SET x = 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if !n.Loc.IsValid() {
		t.Error("SetStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

func TestUnsetVariable(t *testing.T) {
	file, errs := Parse("UNSET VARIABLE myvar")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.UnsetStmt)
	if !ok {
		t.Fatalf("expected *ast.UnsetStmt, got %T", file.Stmts[0])
	}
	if n.Type != "VARIABLE" {
		t.Errorf("Type = %q, want VARIABLE", n.Type)
	}
	if len(n.Names) != 1 || n.Names[0] != "myvar" {
		t.Errorf("Names = %v, want [myvar]", n.Names)
	}
}

func TestUnsetVariableAll(t *testing.T) {
	file, errs := Parse("UNSET VARIABLE ALL")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if !n.All {
		t.Error("All should be true")
	}
}

func TestUnsetGlobalVariable(t *testing.T) {
	file, errs := Parse("UNSET GLOBAL VARIABLE myvar")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if n.Scope != "GLOBAL" {
		t.Errorf("Scope = %q, want GLOBAL", n.Scope)
	}
	if len(n.Names) != 1 || n.Names[0] != "myvar" {
		t.Errorf("Names = %v, want [myvar]", n.Names)
	}
}

func TestUnsetSessionVariable(t *testing.T) {
	file, errs := Parse("UNSET SESSION VARIABLE x")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if n.Scope != "SESSION" {
		t.Errorf("Scope = %q, want SESSION", n.Scope)
	}
}

func TestUnsetLocValid(t *testing.T) {
	file, errs := Parse("UNSET VARIABLE x")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if !n.Loc.IsValid() {
		t.Error("UnsetStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// HELP
// ---------------------------------------------------------------------------

func TestHelp(t *testing.T) {
	file, errs := Parse("HELP 'SELECT'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.HelpStmt)
	if !ok {
		t.Fatalf("expected *ast.HelpStmt, got %T", file.Stmts[0])
	}
	if n.Mask != "SELECT" {
		t.Errorf("Mask = %q, want SELECT", n.Mask)
	}
}

func TestHelpIdentifier(t *testing.T) {
	file, errs := Parse("HELP contents")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.HelpStmt)
	if n.Mask != "contents" {
		t.Errorf("Mask = %q, want contents", n.Mask)
	}
}

func TestHelpLocValid(t *testing.T) {
	file, errs := Parse("HELP 'x'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.HelpStmt)
	if !n.Loc.IsValid() {
		t.Error("HelpStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// NodeTag checks
// ---------------------------------------------------------------------------

func TestUtilityNodeTags(t *testing.T) {
	tests := []struct {
		input   string
		wantTag ast.NodeTag
	}{
		{"SHOW TABLES", ast.T_ShowStmt},
		{"DESC t", ast.T_DescribeStmt},
		{"DESCRIBE t", ast.T_DescribeStmt},
		{"EXPLAIN SELECT 1", ast.T_ExplainStmt},
		{"USE db", ast.T_UseStmt},
		{"SET x = 1", ast.T_SetStmt},
		{"UNSET VARIABLE x", ast.T_UnsetStmt},
		{"HELP 'x'", ast.T_HelpStmt},
	}
	for _, tt := range tests {
		file, errs := Parse(tt.input)
		if len(errs) != 0 {
			t.Errorf("Parse(%q) errors: %v", tt.input, errs)
			continue
		}
		if len(file.Stmts) == 0 {
			t.Errorf("Parse(%q): no stmts", tt.input)
			continue
		}
		if got := file.Stmts[0].Tag(); got != tt.wantTag {
			t.Errorf("Parse(%q): Tag() = %v, want %v", tt.input, got, tt.wantTag)
		}
	}
}

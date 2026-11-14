package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ANALYZE
// ---------------------------------------------------------------------------

func TestAnalyze_Table(t *testing.T) {
	file, errs := parseForTest("ANALYZE TABLE lineitem")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.AnalyzeStmt, got %T", file.Stmts[0])
	}
	if stmt.TargetType != "TABLE" {
		t.Errorf("TargetType = %q, want TABLE", stmt.TargetType)
	}
	if stmt.Target == nil || stmt.Target.String() != "lineitem" {
		t.Errorf("Target = %v, want lineitem", stmt.Target)
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("expected no columns, got %v", stmt.Columns)
	}
}

func TestAnalyze_TableWithColumns(t *testing.T) {
	file, errs := parseForTest("ANALYZE TABLE t (col1, col2)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if len(stmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Columns))
	}
	if stmt.Columns[0] != "col1" || stmt.Columns[1] != "col2" {
		t.Errorf("columns = %v, want [col1 col2]", stmt.Columns)
	}
}

// ---------------------------------------------------------------------------
// SHOW ANALYZE
// ---------------------------------------------------------------------------

func TestShowAnalyze_Basic(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW ANALYZE test1 WHERE STATE=\"FINISHED\"")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowAnalyzeStmt, got %T", file.Stmts[0])
	}
	if stmt.For == nil || stmt.For.String() != "test1" {
		t.Errorf("For = %v, want test1", stmt.For)
	}
	if stmt.Where == nil {
		t.Error("expected WHERE clause")
	}
}

func TestShowAnalyze_ByJobID(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW ANALYZE 1738725887903")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowAnalyzeStmt, got %T", file.Stmts[0])
	}
	if stmt.JobID != 1738725887903 {
		t.Errorf("JobID = %d, want 1738725887903", stmt.JobID)
	}
}

func TestShowAnalyze_All(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW ALL ANALYZE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowAnalyzeStmt, got %T", file.Stmts[0])
	}
	if !stmt.All {
		t.Error("expected All = true")
	}
}

func TestShowAnalyze_Queued(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW QUEUED ANALYZE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowAnalyzeStmt, got %T", file.Stmts[0])
	}
	if !stmt.Queued {
		t.Error("expected Queued = true")
	}
}

func TestShowAnalyze_TaskStatus(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW ANALYZE TASK STATUS 12345")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowAnalyzeStmt, got %T", file.Stmts[0])
	}
	if !stmt.IsTask {
		t.Error("expected IsTask = true")
	}
	if stmt.JobID != 12345 {
		t.Errorf("JobID = %d, want 12345", stmt.JobID)
	}
}

func TestShowAnalyze_ForTable(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW ANALYZE FOR mydb.t1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if stmt.For == nil || stmt.For.String() != "mydb.t1" {
		t.Errorf("For = %v, want mydb.t1", stmt.For)
	}
}

// ---------------------------------------------------------------------------
// SHOW STATS
// ---------------------------------------------------------------------------

func TestShowTableStats(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW TABLE STATS test1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ShowStatsStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowStatsStmt, got %T", file.Stmts[0])
	}
	if stmt.Type != "TABLE" {
		t.Errorf("Type = %q, want TABLE", stmt.Type)
	}
	if stmt.Target == nil || stmt.Target.String() != "test1" {
		t.Errorf("Target = %v, want test1", stmt.Target)
	}
}

func TestShowColumnStats(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW COLUMN STATS mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowStatsStmt)
	if stmt.Type != "COLUMN" {
		t.Errorf("Type = %q, want COLUMN", stmt.Type)
	}
}

func TestShowStats_NoType(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW STATS mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowStatsStmt)
	if stmt.Type != "" {
		t.Errorf("Type = %q, want empty", stmt.Type)
	}
	if stmt.Target == nil || stmt.Target.String() != "mytable" {
		t.Errorf("Target = %v, want mytable", stmt.Target)
	}
}

// ---------------------------------------------------------------------------
// SHOW CONSTRAINTS
// ---------------------------------------------------------------------------

func TestShowConstraints_FromTable(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW CONSTRAINTS FROM mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ShowConstraintsStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowConstraintsStmt, got %T", file.Stmts[0])
	}
	if stmt.Table == nil || stmt.Table.String() != "mytable" {
		t.Errorf("Table = %v, want mytable", stmt.Table)
	}
}

func TestShowConstraints_QualifiedTable(t *testing.T) {
	file, errs := parseDorisFormForTest("SHOW CONSTRAINTS FROM mydb.mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowConstraintsStmt)
	if stmt.Table.String() != "mydb.mytable" {
		t.Errorf("Table = %v, want mydb.mytable", stmt.Table)
	}
}

// ---------------------------------------------------------------------------
// DROP STATS
// ---------------------------------------------------------------------------

func TestDropStats_Table(t *testing.T) {
	file, errs := parseForTest("DROP STATS table1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropStatsStmt)
	if !ok {
		t.Fatalf("expected *ast.DropStatsStmt, got %T", file.Stmts[0])
	}
	if stmt.Variant != "" {
		t.Errorf("Variant = %q, want empty", stmt.Variant)
	}
	if stmt.Target == nil || stmt.Target.String() != "table1" {
		t.Errorf("Target = %v, want table1", stmt.Target)
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("expected no columns, got %v", stmt.Columns)
	}
}

// ---------------------------------------------------------------------------
// KILL ANALYZE
// ---------------------------------------------------------------------------

func TestKillAnalyze(t *testing.T) {
	file, errs := parseForTest("KILL ANALYZE 12345")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.KillAnalyzeStmt)
	if !ok {
		t.Fatalf("expected *ast.KillAnalyzeStmt, got %T", file.Stmts[0])
	}
	if stmt.JobID != 12345 {
		t.Errorf("JobID = %d, want 12345", stmt.JobID)
	}
}

// ---------------------------------------------------------------------------
// ADD / DROP CONSTRAINT
// ---------------------------------------------------------------------------

func TestAlterTableAddConstraintPrimaryKey(t *testing.T) {
	file, errs := parseForTest("ALTER TABLE t ADD CONSTRAINT pk_id PRIMARY KEY (id)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AddConstraintStmt)
	if !ok {
		t.Fatalf("expected *ast.AddConstraintStmt, got %T", file.Stmts[0])
	}
	if stmt.Table == nil || stmt.Table.String() != "t" {
		t.Errorf("Table = %v, want t", stmt.Table)
	}
	if stmt.Name != "pk_id" {
		t.Errorf("Name = %q, want pk_id", stmt.Name)
	}
	if stmt.Type != "PRIMARY KEY" {
		t.Errorf("Type = %q, want PRIMARY KEY", stmt.Type)
	}
	if len(stmt.Columns) != 1 || stmt.Columns[0] != "id" {
		t.Errorf("Columns = %v, want [id]", stmt.Columns)
	}
}

func TestAlterTableAddConstraintUnique(t *testing.T) {
	file, errs := parseForTest("ALTER TABLE t ADD CONSTRAINT uq_email UNIQUE (email)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AddConstraintStmt)
	if !ok {
		t.Fatalf("expected *ast.AddConstraintStmt, got %T", file.Stmts[0])
	}
	if stmt.Type != "UNIQUE" {
		t.Errorf("Type = %q, want UNIQUE", stmt.Type)
	}
	if len(stmt.Columns) != 1 || stmt.Columns[0] != "email" {
		t.Errorf("Columns = %v, want [email]", stmt.Columns)
	}
}

func TestAlterTableAddConstraintForeignKey(t *testing.T) {
	file, errs := parseForTest("ALTER TABLE orders ADD CONSTRAINT fk_cust FOREIGN KEY (customer_id) REFERENCES customers (id)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AddConstraintStmt)
	if !ok {
		t.Fatalf("expected *ast.AddConstraintStmt, got %T", file.Stmts[0])
	}
	if stmt.Type != "FOREIGN KEY" {
		t.Errorf("Type = %q, want FOREIGN KEY", stmt.Type)
	}
	if len(stmt.Columns) != 1 || stmt.Columns[0] != "customer_id" {
		t.Errorf("Columns = %v, want [customer_id]", stmt.Columns)
	}
	if stmt.RefTable == nil || stmt.RefTable.String() != "customers" {
		t.Errorf("RefTable = %v, want customers", stmt.RefTable)
	}
	if len(stmt.RefColumns) != 1 || stmt.RefColumns[0] != "id" {
		t.Errorf("RefColumns = %v, want [id]", stmt.RefColumns)
	}
}

func TestAlterTableDropConstraint(t *testing.T) {
	file, errs := parseForTest("ALTER TABLE t DROP CONSTRAINT pk_id")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropConstraintStmt)
	if !ok {
		t.Fatalf("expected *ast.DropConstraintStmt, got %T", file.Stmts[0])
	}
	if stmt.Table == nil || stmt.Table.String() != "t" {
		t.Errorf("Table = %v, want t", stmt.Table)
	}
	if stmt.Name != "pk_id" {
		t.Errorf("Name = %q, want pk_id", stmt.Name)
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus forms (from statistics.sql)
// ---------------------------------------------------------------------------

func TestLegacyCorpus_DropStats(t *testing.T) {
	// DROP STATS table1;
	file, errs := parseForTest("DROP STATS table1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStatsStmt)
	if stmt.Target.String() != "table1" {
		t.Errorf("Target = %v, want table1", stmt.Target)
	}
}

func TestLegacyCorpus_ShowTableStats(t *testing.T) {
	// SHOW TABLE STATS test1;
	file, errs := parseDorisFormForTest("SHOW TABLE STATS test1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowStatsStmt)
	if stmt.Type != "TABLE" || stmt.Target.String() != "test1" {
		t.Errorf("Type=%q Target=%v", stmt.Type, stmt.Target)
	}
}

func TestLegacyCorpus_ShowAnalyzeWithWhere(t *testing.T) {
	// SHOW ANALYZE test1 WHERE STATE="FINISHED";
	file, errs := parseDorisFormForTest(`SHOW ANALYZE test1 WHERE STATE="FINISHED"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if stmt.For == nil || stmt.For.String() != "test1" {
		t.Errorf("For = %v, want test1", stmt.For)
	}
	if stmt.Where == nil {
		t.Error("expected WHERE clause")
	}
}

func TestLegacyCorpus_ShowAnalyzeByJobID(t *testing.T) {
	// SHOW ANALYZE 1738725887903
	file, errs := parseDorisFormForTest("SHOW ANALYZE 1738725887903")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowAnalyzeStmt)
	if stmt.JobID != 1738725887903 {
		t.Errorf("JobID = %d, want 1738725887903", stmt.JobID)
	}
}

// ---------------------------------------------------------------------------
// Location tracking sanity
// ---------------------------------------------------------------------------

func TestAnalyze_LocIsValid(t *testing.T) {
	file, errs := parseForTest("ANALYZE TABLE t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if !stmt.Loc.IsValid() {
		t.Error("expected valid Loc")
	}
}

func TestDropStats_LocIsValid(t *testing.T) {
	file, errs := parseForTest("DROP STATS mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStatsStmt)
	if !stmt.Loc.IsValid() {
		t.Error("expected valid Loc")
	}
}

func TestNodeTagsForStatsNodes(t *testing.T) {
	tests := []struct {
		input string
		tag   ast.NodeTag
	}{
		{"ANALYZE TABLE t", ast.T_AnalyzeStmt},
		{"SHOW ANALYZE", ast.T_ShowAnalyzeStmt},
		{"DROP STATS t", ast.T_DropStatsStmt},
		{"KILL ANALYZE 1", ast.T_KillAnalyzeStmt},
		{"ALTER TABLE t ADD CONSTRAINT c PRIMARY KEY (id)", ast.T_AddConstraintStmt},
		{"ALTER TABLE t DROP CONSTRAINT c", ast.T_DropConstraintStmt},
	}

	for _, tt := range tests {
		file, errs := parseForTest(tt.input)
		if len(errs) != 0 {
			t.Errorf("%q: unexpected errors: %v", tt.input, errs)
			continue
		}
		if len(file.Stmts) != 1 {
			t.Errorf("%q: expected 1 stmt, got %d", tt.input, len(file.Stmts))
			continue
		}
		if got := file.Stmts[0].Tag(); got != tt.tag {
			t.Errorf("%q: Tag() = %v, want %v", tt.input, got, tt.tag)
		}
	}
}

// parseDorisFormForTest reads one of the Doris SHOW forms the readers above
// still build, which StarRocks.g4 does not have (SHOW TABLE STATS, SHOW
// CONSTRAINTS, SHOW ANALYZE <table>): Parse rejects them, best effort reads.
func parseDorisFormForTest(sql string) (*ast.File, []ParseError) {
	r := ParseBestEffort(sql)
	return r.File, r.Errors
}

// TestAnalyzeGrammar holds ANALYZE and the stats statements to StarRocks.g4
// (:1297-1354); the rejected forms are Doris's, which this file once pinned.
func TestAnalyzeGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"ANALYZE FULL TABLE db.t (c1, s.c2) WITH ASYNC MODE PROPERTIES('k' = 'v')", true},
		{"ANALYZE SAMPLE TABLE t WITH SYNC MODE", true},
		{"ANALYZE TABLE t UPDATE HISTOGRAM ON c1, c2 WITH 32 BUCKETS PROPERTIES('k' = 'v') WITH SYNC MODE", true},
		{"ANALYZE TABLE t DROP HISTOGRAM ON c1", true},
		{"ANALYZE PROFILE FROM 'q', 1, 2", true},
		{"CREATE ANALYZE FULL ALL PROPERTIES('k' = 'v')", true},
		{"CREATE ANALYZE SAMPLE DATABASE d", true},
		{"CREATE ANALYZE TABLE t (c)", true},
		{"CREATE ANALYZE TABLE t UPDATE HISTOGRAM ON c", true},
		{"CREATE ANALYZE TABLE t DROP HISTOGRAM ON c", false},
		{"CREATE ANALYZE TABLE t WITH SYNC MODE", false},
		{"DROP ANALYZE 1", true},
		{"ANALYZE TABLE t WITH SYNC", false},
		{"ANALYZE TABLE t()", false},
		{"ANALYZE TABLE lineitem WITH SAMPLE PERCENT 10", false},
		{"ANALYZE DATABASE d", false},
		{"ANALYZE t", false},
		{"ANALYZE", false},
		{"DROP STATS t (c1, c2)", false},
		{"DROP EXPIRED STATS", false},
		{"DROP CACHED STATS t", false},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}

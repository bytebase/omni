package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ANALYZE
// ---------------------------------------------------------------------------

func TestAnalyze_Table(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE lineitem")
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

func TestAnalyze_TableWithSamplePercent(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE lineitem WITH SAMPLE PERCENT 10")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "SAMPLE PERCENT" {
		t.Errorf("property key = %q, want SAMPLE PERCENT", stmt.Properties[0].Key)
	}
	if stmt.Properties[0].Value != "10" {
		t.Errorf("property value = %q, want 10", stmt.Properties[0].Value)
	}
}

func TestAnalyze_TableWithSampleRows(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE lineitem WITH SAMPLE ROWS 100000")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "SAMPLE ROWS" {
		t.Errorf("property key = %q, want SAMPLE ROWS", stmt.Properties[0].Key)
	}
	if stmt.Properties[0].Value != "100000" {
		t.Errorf("property value = %q, want 100000", stmt.Properties[0].Value)
	}
}

func TestAnalyze_TableWithSync(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE t WITH SYNC")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "SYNC" {
		t.Errorf("expected SYNC property, got %v", stmt.Properties)
	}
}

func TestAnalyze_TableWithIncremental(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE t WITH INCREMENTAL")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "INCREMENTAL" {
		t.Errorf("expected INCREMENTAL property, got %v", stmt.Properties)
	}
}

func TestAnalyze_TableWithColumns(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE t (col1, col2)")
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

func TestAnalyze_Database(t *testing.T) {
	file, errs := Parse("ANALYZE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if stmt.TargetType != "DATABASE" {
		t.Errorf("TargetType = %q, want DATABASE", stmt.TargetType)
	}
	if stmt.Target.String() != "mydb" {
		t.Errorf("Target = %v, want mydb", stmt.Target)
	}
}

func TestAnalyze_QualifiedTable(t *testing.T) {
	file, errs := Parse("ANALYZE TABLE mydb.lineitem WITH SAMPLE PERCENT 10")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if stmt.Target.String() != "mydb.lineitem" {
		t.Errorf("Target = %v, want mydb.lineitem", stmt.Target)
	}
}

// ---------------------------------------------------------------------------
// SHOW ANALYZE
// ---------------------------------------------------------------------------

func TestShowAnalyze_Basic(t *testing.T) {
	file, errs := Parse("SHOW ANALYZE test1 WHERE STATE=\"FINISHED\"")
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
	file, errs := Parse("SHOW ANALYZE 1738725887903")
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
	file, errs := Parse("SHOW ALL ANALYZE")
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
	file, errs := Parse("SHOW QUEUED ANALYZE")
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
	file, errs := Parse("SHOW ANALYZE TASK STATUS 12345")
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
	file, errs := Parse("SHOW ANALYZE FOR mydb.t1")
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
	file, errs := Parse("SHOW TABLE STATS test1")
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
	file, errs := Parse("SHOW COLUMN STATS mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.ShowStatsStmt)
	if stmt.Type != "COLUMN" {
		t.Errorf("Type = %q, want COLUMN", stmt.Type)
	}
}

func TestShowStats_NoType(t *testing.T) {
	file, errs := Parse("SHOW STATS mytable")
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
	file, errs := Parse("SHOW CONSTRAINTS FROM mytable")
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
	file, errs := Parse("SHOW CONSTRAINTS FROM mydb.mytable")
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
	file, errs := Parse("DROP STATS table1")
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

func TestDropStats_WithColumns(t *testing.T) {
	file, errs := Parse("DROP STATS table1 (col1, col2)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStatsStmt)
	if len(stmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Columns))
	}
	if stmt.Columns[0] != "col1" || stmt.Columns[1] != "col2" {
		t.Errorf("columns = %v, want [col1 col2]", stmt.Columns)
	}
}

func TestDropStats_Expired(t *testing.T) {
	file, errs := Parse("DROP EXPIRED STATS mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropStatsStmt)
	if !ok {
		t.Fatalf("expected *ast.DropStatsStmt, got %T", file.Stmts[0])
	}
	if stmt.Variant != "EXPIRED" {
		t.Errorf("Variant = %q, want EXPIRED", stmt.Variant)
	}
}

func TestDropStats_Cached(t *testing.T) {
	file, errs := Parse("DROP CACHED STATS mytable")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DropStatsStmt)
	if !ok {
		t.Fatalf("expected *ast.DropStatsStmt, got %T", file.Stmts[0])
	}
	if stmt.Variant != "CACHED" {
		t.Errorf("Variant = %q, want CACHED", stmt.Variant)
	}
	if stmt.Target.String() != "mytable" {
		t.Errorf("Target = %v, want mytable", stmt.Target)
	}
}

// ---------------------------------------------------------------------------
// KILL ANALYZE
// ---------------------------------------------------------------------------

func TestKillAnalyze(t *testing.T) {
	file, errs := Parse("KILL ANALYZE 12345")
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
	file, errs := Parse("ALTER TABLE t ADD CONSTRAINT pk_id PRIMARY KEY (id)")
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
	file, errs := Parse("ALTER TABLE t ADD CONSTRAINT uq_email UNIQUE (email)")
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
	file, errs := Parse("ALTER TABLE orders ADD CONSTRAINT fk_cust FOREIGN KEY (customer_id) REFERENCES customers (id)")
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
	file, errs := Parse("ALTER TABLE t DROP CONSTRAINT pk_id")
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

func TestLegacyCorpus_AnalyzeWithSamplePercent(t *testing.T) {
	// ANALYZE TABLE lineitem WITH SAMPLE PERCENT 10;
	file, errs := Parse("ANALYZE TABLE lineitem WITH SAMPLE PERCENT 10")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if stmt.Target.String() != "lineitem" {
		t.Errorf("Target = %v, want lineitem", stmt.Target)
	}
	if stmt.Properties[0].Key != "SAMPLE PERCENT" || stmt.Properties[0].Value != "10" {
		t.Errorf("property = %v", stmt.Properties[0])
	}
}

func TestLegacyCorpus_AnalyzeWithSampleRows(t *testing.T) {
	// ANALYZE TABLE lineitem WITH SAMPLE ROWS 100000;
	file, errs := Parse("ANALYZE TABLE lineitem WITH SAMPLE ROWS 100000")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if stmt.Properties[0].Key != "SAMPLE ROWS" || stmt.Properties[0].Value != "100000" {
		t.Errorf("property = %v", stmt.Properties[0])
	}
}

func TestLegacyCorpus_DropStats(t *testing.T) {
	// DROP STATS table1;
	file, errs := Parse("DROP STATS table1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStatsStmt)
	if stmt.Target.String() != "table1" {
		t.Errorf("Target = %v, want table1", stmt.Target)
	}
}

func TestLegacyCorpus_DropStatsWithColumns(t *testing.T) {
	// DROP STATS table1 (col1, col2);
	file, errs := Parse("DROP STATS table1 (col1, col2)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropStatsStmt)
	if len(stmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Columns))
	}
}

func TestLegacyCorpus_ShowTableStats(t *testing.T) {
	// SHOW TABLE STATS test1;
	file, errs := Parse("SHOW TABLE STATS test1")
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
	file, errs := Parse(`SHOW ANALYZE test1 WHERE STATE="FINISHED"`)
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
	file, errs := Parse("SHOW ANALYZE 1738725887903")
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
	file, errs := Parse("ANALYZE TABLE t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AnalyzeStmt)
	if !stmt.Loc.IsValid() {
		t.Error("expected valid Loc")
	}
}

func TestDropStats_LocIsValid(t *testing.T) {
	file, errs := Parse("DROP STATS mytable")
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
		{"SHOW TABLE STATS t", ast.T_ShowStatsStmt},
		{"SHOW CONSTRAINTS FROM t", ast.T_ShowConstraintsStmt},
		{"DROP STATS t", ast.T_DropStatsStmt},
		{"KILL ANALYZE 1", ast.T_KillAnalyzeStmt},
		{"ALTER TABLE t ADD CONSTRAINT c PRIMARY KEY (id)", ast.T_AddConstraintStmt},
		{"ALTER TABLE t DROP CONSTRAINT c", ast.T_DropConstraintStmt},
	}

	for _, tt := range tests {
		file, errs := Parse(tt.input)
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

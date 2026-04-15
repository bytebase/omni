package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// mustParseInsert parses input and returns the first statement as *ast.InsertStmt.
func mustParseInsert(t *testing.T, input string) *ast.InsertStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.InsertStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.InsertStmt", input, file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Basic INSERT INTO
// ---------------------------------------------------------------------------

func TestInsertBasicValues(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t VALUES (1, 'a')")

	if stmt.Overwrite {
		t.Error("Overwrite = true, want false")
	}
	if stmt.Target == nil || stmt.Target.Parts[0] != "t" {
		t.Errorf("Target = %v, want t", stmt.Target)
	}
	if len(stmt.Values) != 1 {
		t.Fatalf("Values rows = %d, want 1", len(stmt.Values))
	}
	if len(stmt.Values[0]) != 2 {
		t.Fatalf("Values[0] len = %d, want 2", len(stmt.Values[0]))
	}
	lit0, ok := stmt.Values[0][0].(*ast.Literal)
	if !ok || lit0.Kind != ast.LitInt || lit0.Value != "1" {
		t.Errorf("Values[0][0] = %v, want Literal{Int,1}", stmt.Values[0][0])
	}
	lit1, ok := stmt.Values[0][1].(*ast.Literal)
	if !ok || lit1.Kind != ast.LitString || lit1.Value != "a" {
		t.Errorf("Values[0][1] = %v, want Literal{String,a}", stmt.Values[0][1])
	}
	if stmt.Query != nil {
		t.Error("Query should be nil for VALUES form")
	}
}

// ---------------------------------------------------------------------------
// Multi-row VALUES
// ---------------------------------------------------------------------------

func TestInsertMultiRowValues(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t VALUES (1, 'a'), (2, 'b')")

	if len(stmt.Values) != 2 {
		t.Fatalf("Values rows = %d, want 2", len(stmt.Values))
	}
	if len(stmt.Values[0]) != 2 {
		t.Fatalf("row 0 len = %d, want 2", len(stmt.Values[0]))
	}
	if len(stmt.Values[1]) != 2 {
		t.Fatalf("row 1 len = %d, want 2", len(stmt.Values[1]))
	}
}

// ---------------------------------------------------------------------------
// VALUES with expressions
// ---------------------------------------------------------------------------

func TestInsertValuesWithExpr(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO test VALUES (1, 2), (3, 2 + 2)")

	if len(stmt.Values) != 2 {
		t.Fatalf("Values rows = %d, want 2", len(stmt.Values))
	}
	// Second row, second expression should be a binary expression
	_, ok := stmt.Values[1][1].(*ast.BinaryExpr)
	if !ok {
		t.Errorf("Values[1][1] = %T, want *ast.BinaryExpr", stmt.Values[1][1])
	}
}

// ---------------------------------------------------------------------------
// Column list
// ---------------------------------------------------------------------------

func TestInsertWithColumns(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t (c1, c2) VALUES (1, 'a')")

	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns = %d, want 2", len(stmt.Columns))
	}
	if stmt.Columns[0] != "c1" || stmt.Columns[1] != "c2" {
		t.Errorf("Columns = %v, want [c1 c2]", stmt.Columns)
	}
	if len(stmt.Values) != 1 {
		t.Fatalf("Values rows = %d, want 1", len(stmt.Values))
	}
}

// ---------------------------------------------------------------------------
// DEFAULT keyword in VALUES
// ---------------------------------------------------------------------------

func TestInsertValuesDefault(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO test (c1, c2) VALUES (1, DEFAULT)")

	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns = %d, want 2", len(stmt.Columns))
	}
	if len(stmt.Values) != 1 || len(stmt.Values[0]) != 2 {
		t.Fatalf("Values = %v, want 1 row with 2 items", stmt.Values)
	}
	lit, ok := stmt.Values[0][1].(*ast.Literal)
	if !ok || lit.Kind != ast.LitKeyword || lit.Value != "DEFAULT" {
		t.Errorf("Values[0][1] = %v, want Literal{Keyword,DEFAULT}", stmt.Values[0][1])
	}
}

// ---------------------------------------------------------------------------
// SELECT source
// ---------------------------------------------------------------------------

func TestInsertSelectSource(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t SELECT * FROM other")

	if stmt.Overwrite {
		t.Error("Overwrite = true, want false")
	}
	if stmt.Target.Parts[0] != "t" {
		t.Errorf("Target = %v, want t", stmt.Target)
	}
	if stmt.Values != nil {
		t.Error("Values should be nil for SELECT form")
	}
	sel, ok := stmt.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
	if len(sel.From) != 1 {
		t.Fatalf("FROM = %d, want 1", len(sel.From))
	}
}

func TestInsertSelectWithColumns(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO test (c1, c2) SELECT * FROM test2")

	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns = %d, want 2", len(stmt.Columns))
	}
	if _, ok := stmt.Query.(*ast.SelectStmt); !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
}

// ---------------------------------------------------------------------------
// INSERT OVERWRITE TABLE
// ---------------------------------------------------------------------------

func TestInsertOverwriteValues(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT OVERWRITE TABLE test VALUES (1, 2)")

	if !stmt.Overwrite {
		t.Error("Overwrite = false, want true")
	}
	if stmt.Target.Parts[0] != "test" {
		t.Errorf("Target = %v, want test", stmt.Target)
	}
	if len(stmt.Values) != 1 {
		t.Fatalf("Values rows = %d, want 1", len(stmt.Values))
	}
}

func TestInsertOverwriteSelect(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT OVERWRITE TABLE t SELECT * FROM src")

	if !stmt.Overwrite {
		t.Error("Overwrite = false, want true")
	}
	if _, ok := stmt.Query.(*ast.SelectStmt); !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
}

// ---------------------------------------------------------------------------
// PARTITION clause
// ---------------------------------------------------------------------------

func TestInsertWithPartition(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t PARTITION(p1, p2) SELECT * FROM other")

	if len(stmt.Partition) != 2 {
		t.Fatalf("Partition = %v, want [p1 p2]", stmt.Partition)
	}
	if stmt.Partition[0] != "p1" || stmt.Partition[1] != "p2" {
		t.Errorf("Partition = %v, want [p1 p2]", stmt.Partition)
	}
	if stmt.PartitionStar {
		t.Error("PartitionStar = true, want false")
	}
}

func TestInsertOverwritePartition(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT OVERWRITE TABLE test PARTITION(p1, p2) VALUES (1, 2)")

	if !stmt.Overwrite {
		t.Error("Overwrite = false, want true")
	}
	if len(stmt.Partition) != 2 {
		t.Fatalf("Partition = %v, want [p1 p2]", stmt.Partition)
	}
}

func TestInsertPartitionStar(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT OVERWRITE TABLE test PARTITION(*) VALUES (3)")

	if !stmt.PartitionStar {
		t.Error("PartitionStar = false, want true")
	}
	if len(stmt.Partition) != 0 {
		t.Errorf("Partition = %v, want nil", stmt.Partition)
	}
}

// ---------------------------------------------------------------------------
// WITH LABEL clause
// ---------------------------------------------------------------------------

func TestInsertWithLabel(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t WITH LABEL my_label VALUES (1)")

	if stmt.Label != "my_label" {
		t.Errorf("Label = %q, want my_label", stmt.Label)
	}
}

func TestInsertWithLabelBacktick(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO test PARTITION(p1, p2) WITH LABEL `label1` SELECT * FROM test2")

	if stmt.Label != "label1" {
		t.Errorf("Label = %q, want label1", stmt.Label)
	}
	if len(stmt.Partition) != 2 {
		t.Fatalf("Partition = %v, want [p1 p2]", stmt.Partition)
	}
}

func TestInsertWithLabelAndColumns(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO test WITH LABEL `label1` (c1, c2) SELECT * FROM test2")

	if stmt.Label != "label1" {
		t.Errorf("Label = %q, want label1", stmt.Label)
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns = %v, want [c1 c2]", stmt.Columns)
	}
}

// ---------------------------------------------------------------------------
// CTE (WITH ... SELECT) source
// ---------------------------------------------------------------------------

func TestInsertWithCTE(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t WITH c AS (SELECT 1) SELECT * FROM c")

	if stmt.Query == nil {
		t.Fatal("Query is nil, want a *ast.SelectStmt with WITH clause")
	}
	sel, ok := stmt.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}
	if sel.With == nil {
		t.Error("Query.With is nil, want WithClause")
	}
	if len(sel.With.CTEs) != 1 {
		t.Fatalf("CTEs = %d, want 1", len(sel.With.CTEs))
	}
	if sel.With.CTEs[0].Name != "c" {
		t.Errorf("CTE name = %q, want c", sel.With.CTEs[0].Name)
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus forms
// ---------------------------------------------------------------------------

// TestInsertLegacyCorpusDML verifies that all statements from dml_insert.sql
// parse without errors.
func TestInsertLegacyCorpusDML(t *testing.T) {
	cases := []string{
		"INSERT INTO test VALUES (1, 2)",
		"INSERT INTO test (c1, c2) VALUES (1, 2)",
		"INSERT INTO test (c1, c2) VALUES (1, DEFAULT)",
		"INSERT INTO test (c1) VALUES (1)",
		"INSERT INTO test VALUES (1, 2), (3, 2 + 2)",
		"INSERT INTO test (c1, c2) VALUES (1, 2), (3, 2 * 2)",
		"INSERT INTO test (c1) VALUES (1), (3)",
		"INSERT INTO test (c1, c2) VALUES (1, DEFAULT), (3, DEFAULT)",
		"INSERT INTO test SELECT * FROM test2",
		"INSERT INTO test (c1, c2) SELECT * FROM test2",
		"INSERT INTO test PARTITION(p1, p2) WITH LABEL `label1` SELECT * FROM test2",
		"INSERT INTO test WITH LABEL `label1` (c1, c2) SELECT * FROM test2",
		"INSERT INTO tbl1 SELECT * FROM empty_tbl",
		"INSERT INTO tbl1 SELECT * FROM tbl2",
		"INSERT INTO tbl1 WITH LABEL my_label1 SELECT * FROM tbl2",
		`INSERT INTO tbl1 SELECT * FROM tbl2 WHERE k1 = "a"`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) > 0 {
				t.Fatalf("Parse(%q) errors: %v", sql, errs)
			}
		})
	}
}

// TestInsertLegacyCorpusOverwrite verifies all statements from dml_insert_overwrite.sql.
func TestInsertLegacyCorpusOverwrite(t *testing.T) {
	cases := []string{
		"INSERT OVERWRITE TABLE test VALUES (1, 2)",
		"INSERT OVERWRITE TABLE test (c1, c2) VALUES (1, 2)",
		"INSERT OVERWRITE TABLE test (c1, c2) VALUES (1, DEFAULT)",
		"INSERT OVERWRITE TABLE test (c1) VALUES (1)",
		"INSERT OVERWRITE TABLE test VALUES (1, 2), (3, 2 + 2)",
		"INSERT OVERWRITE TABLE test (c1, c2) VALUES (1, 2), (3, 2 * 2)",
		"INSERT OVERWRITE TABLE test (c1, c2) VALUES (1, DEFAULT), (3, DEFAULT)",
		"INSERT OVERWRITE TABLE test (c1) VALUES (1), (3)",
		"INSERT OVERWRITE TABLE test SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test (c1, c2) SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test WITH LABEL `label1` SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test WITH LABEL `label2` (c1, c2) SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) VALUES (1, 2)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1, c2) VALUES (1, 2)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1, c2) VALUES (1, DEFAULT)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1) VALUES (1)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) VALUES (1, 2), (4, 2 + 2)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1, c2) VALUES (1, 2), (4, 2 * 2)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1, c2) VALUES (1, DEFAULT), (4, DEFAULT)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1) VALUES (1), (4)",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) (c1, c2) SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) WITH LABEL `label3` SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test PARTITION(p1, p2) WITH LABEL `label4` (c1, c2) SELECT * FROM test2",
		"INSERT OVERWRITE TABLE test PARTITION(*) VALUES (3), (1234)",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) > 0 {
				t.Fatalf("Parse(%q) errors: %v", sql, errs)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Loc coverage
// ---------------------------------------------------------------------------

func TestInsertLocCoverage(t *testing.T) {
	input := "INSERT INTO t VALUES (1)"
	stmt := mustParseInsert(t, input)

	if !stmt.Loc.IsValid() {
		t.Error("Loc is invalid")
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("Loc.End (%d) <= Loc.Start (%d)", stmt.Loc.End, stmt.Loc.Start)
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestInsertNodeTag(t *testing.T) {
	stmt := mustParseInsert(t, "INSERT INTO t VALUES (1)")
	if stmt.Tag() != ast.T_InsertStmt {
		t.Errorf("Tag() = %v, want T_InsertStmt", stmt.Tag())
	}
}

package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// helper: parse a single statement and expect no errors, returning the node.
func parseIndexStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	return file.Stmts[0]
}

func TestCreateIndexBasic(t *testing.T) {
	sql := "CREATE INDEX idx ON t (col1)"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.Name.String() != "idx" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "idx")
	}
	if stmt.Table.String() != "t" {
		t.Errorf("Table = %q, want %q", stmt.Table.String(), "t")
	}
	if len(stmt.Columns) != 1 || stmt.Columns[0] != "col1" {
		t.Errorf("Columns = %v, want [col1]", stmt.Columns)
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if stmt.IndexType != "" {
		t.Errorf("IndexType = %q, want empty", stmt.IndexType)
	}
}

func TestCreateIndexIfNotExists(t *testing.T) {
	sql := "CREATE INDEX IF NOT EXISTS idx ON t (col1, col2)"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns = %v, want 2 columns", stmt.Columns)
	}
	if stmt.Columns[0] != "col1" || stmt.Columns[1] != "col2" {
		t.Errorf("Columns = %v, want [col1 col2]", stmt.Columns)
	}
}

func TestCreateIndexUsingBitmap(t *testing.T) {
	sql := "CREATE INDEX idx ON t (col1) USING BITMAP"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.IndexType != "bitmap" {
		t.Errorf("IndexType = %q, want %q", stmt.IndexType, "bitmap")
	}
}

func TestCreateIndexUsingInverted(t *testing.T) {
	sql := "CREATE INDEX idx ON t (col1) USING INVERTED"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.IndexType != "inverted" {
		t.Errorf("IndexType = %q, want %q", stmt.IndexType, "inverted")
	}
}

func TestCreateIndexNgramBfWithProperties(t *testing.T) {
	sql := `CREATE INDEX idx ON t (col1) USING NGRAM_BF PROPERTIES("gram_size"="3")`
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.IndexType != "ngram_bf" {
		t.Errorf("IndexType = %q, want %q", stmt.IndexType, "ngram_bf")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties = %v, want 1 property", stmt.Properties)
	}
	if stmt.Properties[0].Key != "gram_size" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "gram_size")
	}
	if stmt.Properties[0].Value != "3" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "3")
	}
}

func TestCreateIndexWithComment(t *testing.T) {
	sql := "CREATE INDEX idx ON t (col1) COMMENT 'my index'"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.Comment != "my index" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "my index")
	}
}

func TestCreateIndexQualifiedTable(t *testing.T) {
	sql := "CREATE INDEX idx ON db.t (col1) USING BITMAP"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.Table.String() != "db.t" {
		t.Errorf("Table = %q, want %q", stmt.Table.String(), "db.t")
	}
}

func TestCreateIndexMultipleProperties(t *testing.T) {
	sql := `CREATE INDEX idx ON t (col1) USING NGRAM_BF PROPERTIES("gram_size"="3", "bf_size"="64")`
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("Properties = %v, want 2 properties", stmt.Properties)
	}
	if stmt.Properties[1].Key != "bf_size" {
		t.Errorf("Properties[1].Key = %q, want %q", stmt.Properties[1].Key, "bf_size")
	}
}

func TestDropIndexBasic(t *testing.T) {
	sql := "DROP INDEX idx ON t"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.DropIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.DropIndexStmt, got %T", node)
	}
	if stmt.Name.String() != "idx" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "idx")
	}
	if stmt.Table.String() != "t" {
		t.Errorf("Table = %q, want %q", stmt.Table.String(), "t")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropIndexIfExists(t *testing.T) {
	sql := "DROP INDEX IF EXISTS idx ON t"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.DropIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.DropIndexStmt, got %T", node)
	}
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "idx" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "idx")
	}
	if stmt.Table.String() != "t" {
		t.Errorf("Table = %q, want %q", stmt.Table.String(), "t")
	}
}

func TestBuildIndexBasic(t *testing.T) {
	sql := "BUILD INDEX idx ON t"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.BuildIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.BuildIndexStmt, got %T", node)
	}
	if stmt.Name.String() != "idx" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "idx")
	}
	if stmt.Table.String() != "t" {
		t.Errorf("Table = %q, want %q", stmt.Table.String(), "t")
	}
	if len(stmt.Partitions) != 0 {
		t.Errorf("Partitions = %v, want empty", stmt.Partitions)
	}
}

func TestBuildIndexWithPartitions(t *testing.T) {
	sql := "BUILD INDEX idx ON t PARTITIONS(p1, p2)"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.BuildIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.BuildIndexStmt, got %T", node)
	}
	if len(stmt.Partitions) != 2 {
		t.Fatalf("Partitions = %v, want 2 partitions", stmt.Partitions)
	}
	if stmt.Partitions[0] != "p1" || stmt.Partitions[1] != "p2" {
		t.Errorf("Partitions = %v, want [p1 p2]", stmt.Partitions)
	}
}

func TestBuildIndexQualifiedTable(t *testing.T) {
	sql := "BUILD INDEX idx ON db.t PARTITIONS(p1)"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.BuildIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.BuildIndexStmt, got %T", node)
	}
	if stmt.Table.String() != "db.t" {
		t.Errorf("Table = %q, want %q", stmt.Table.String(), "db.t")
	}
}

func TestCreateIndexLocIsSet(t *testing.T) {
	sql := "CREATE INDEX idx ON t (col1)"
	node := parseIndexStmt(t, sql)

	stmt := node.(*ast.CreateIndexStmt)
	if !stmt.Loc.IsValid() {
		t.Errorf("Loc = %v, want valid location", stmt.Loc)
	}
	// Loc.Start should be 0 (start of CREATE keyword).
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
}

func TestDropIndexLocIsSet(t *testing.T) {
	sql := "DROP INDEX idx ON t"
	node := parseIndexStmt(t, sql)

	stmt := node.(*ast.DropIndexStmt)
	if !stmt.Loc.IsValid() {
		t.Errorf("Loc = %v, want valid location", stmt.Loc)
	}
}

func TestBuildIndexLocIsSet(t *testing.T) {
	sql := "BUILD INDEX idx ON t"
	node := parseIndexStmt(t, sql)

	stmt := node.(*ast.BuildIndexStmt)
	if !stmt.Loc.IsValid() {
		t.Errorf("Loc = %v, want valid location", stmt.Loc)
	}
}

func TestCreateIndexANN(t *testing.T) {
	sql := "CREATE INDEX idx ON t (embedding_col) USING ANN"
	node := parseIndexStmt(t, sql)

	stmt, ok := node.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateIndexStmt, got %T", node)
	}
	if stmt.IndexType != "ann" {
		t.Errorf("IndexType = %q, want %q", stmt.IndexType, "ann")
	}
}

// TestCreateIndexNodeTag verifies Tag() returns the correct value.
func TestCreateIndexNodeTag(t *testing.T) {
	stmt := &ast.CreateIndexStmt{}
	if stmt.Tag() != ast.T_CreateIndexStmt {
		t.Errorf("Tag() = %v, want T_CreateIndexStmt", stmt.Tag())
	}
}

func TestDropIndexNodeTag(t *testing.T) {
	stmt := &ast.DropIndexStmt{}
	if stmt.Tag() != ast.T_DropIndexStmt {
		t.Errorf("Tag() = %v, want T_DropIndexStmt", stmt.Tag())
	}
}

func TestBuildIndexNodeTag(t *testing.T) {
	stmt := &ast.BuildIndexStmt{}
	if stmt.Tag() != ast.T_BuildIndexStmt {
		t.Errorf("Tag() = %v, want T_BuildIndexStmt", stmt.Tag())
	}
}

// TestCreateIndexStillUnsupportedForNonIndex verifies that CREATE TABLE
// still returns the "not yet supported" error path.
func TestCreateIndexStillUnsupportedForNonIndex(t *testing.T) {
	_, errs := Parse("CREATE TABLE t (id INT)")
	if len(errs) == 0 {
		t.Fatal("expected error for CREATE TABLE")
	}
	want := "CREATE statement parsing is not yet supported"
	if errs[0].Msg != want {
		t.Errorf("got %q, want %q", errs[0].Msg, want)
	}
}

// TestDropIndexStillUnsupportedForNonIndex verifies that DROP TABLE
// still returns the "not yet supported" error path.
func TestDropIndexStillUnsupportedForNonIndex(t *testing.T) {
	_, errs := Parse("DROP TABLE t")
	if len(errs) == 0 {
		t.Fatal("expected error for DROP TABLE")
	}
	want := "DROP statement parsing is not yet supported"
	if errs[0].Msg != want {
		t.Errorf("got %q, want %q", errs[0].Msg, want)
	}
}

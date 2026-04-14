package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseCreateViewStmt(t *testing.T, sql string) *ast.CreateViewStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateViewStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateViewStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseAlterViewStmt(t *testing.T, sql string) *ast.AlterViewStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterViewStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterViewStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseDropViewStmt(t *testing.T, sql string) *ast.DropViewStmt {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropViewStmt)
	if !ok {
		t.Fatalf("expected *ast.DropViewStmt, got %T", file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// CREATE VIEW
// ---------------------------------------------------------------------------

func TestCreateView_Basic(t *testing.T) {
	stmt := parseCreateViewStmt(t, "CREATE VIEW v AS SELECT * FROM t")
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
	if stmt.OrReplace {
		t.Error("OrReplace should be false")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("Columns: got %d, want 0", len(stmt.Columns))
	}
	if stmt.Comment != "" {
		t.Errorf("Comment = %q, want empty", stmt.Comment)
	}
	if stmt.Query == nil {
		t.Fatal("Query is nil")
	}
	if stmt.Query.Tag() != ast.T_SelectStmt {
		t.Errorf("Query.Tag() = %v, want T_SelectStmt", stmt.Query.Tag())
	}
}

func TestCreateView_OrReplace(t *testing.T) {
	stmt := parseCreateViewStmt(t, "CREATE OR REPLACE VIEW v AS SELECT * FROM t")
	if !stmt.OrReplace {
		t.Error("OrReplace should be true")
	}
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
}

func TestCreateView_IfNotExists(t *testing.T) {
	stmt := parseCreateViewStmt(t, "CREATE VIEW IF NOT EXISTS v AS SELECT * FROM t")
	if stmt.IfNotExists != true {
		t.Error("IfNotExists should be true")
	}
	if stmt.OrReplace {
		t.Error("OrReplace should be false")
	}
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
}

func TestCreateView_WithColumnList(t *testing.T) {
	stmt := parseCreateViewStmt(t, "CREATE VIEW v (id, name) AS SELECT id, name FROM t")
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns: got %d, want 2", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "id" {
		t.Errorf("Columns[0].Name = %q, want %q", stmt.Columns[0].Name, "id")
	}
	if stmt.Columns[0].Comment != "" {
		t.Errorf("Columns[0].Comment = %q, want empty", stmt.Columns[0].Comment)
	}
	if stmt.Columns[1].Name != "name" {
		t.Errorf("Columns[1].Name = %q, want %q", stmt.Columns[1].Name, "name")
	}
}

func TestCreateView_WithColumnComments(t *testing.T) {
	stmt := parseCreateViewStmt(t, `CREATE VIEW v (id COMMENT 'identifier', name) AS SELECT id, name FROM t`)
	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns: got %d, want 2", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "id" {
		t.Errorf("Columns[0].Name = %q, want %q", stmt.Columns[0].Name, "id")
	}
	if stmt.Columns[0].Comment != "identifier" {
		t.Errorf("Columns[0].Comment = %q, want %q", stmt.Columns[0].Comment, "identifier")
	}
	if stmt.Columns[1].Comment != "" {
		t.Errorf("Columns[1].Comment = %q, want empty", stmt.Columns[1].Comment)
	}
}

func TestCreateView_WithViewComment(t *testing.T) {
	stmt := parseCreateViewStmt(t, `CREATE VIEW v COMMENT 'my view' AS SELECT * FROM t`)
	if stmt.Comment != "my view" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "my view")
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("Columns: got %d, want 0", len(stmt.Columns))
	}
}

func TestCreateView_QualifiedName(t *testing.T) {
	stmt := parseCreateViewStmt(t, "CREATE VIEW example_db.example_view (k1, k2, k3, v1) AS SELECT c1 AS k1, k2, k3, SUM(v1) FROM example_table WHERE k1 = 20160112 GROUP BY k1, k2, k3")
	if stmt.Name.String() != "example_db.example_view" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "example_db.example_view")
	}
	if len(stmt.Columns) != 4 {
		t.Fatalf("Columns: got %d, want 4", len(stmt.Columns))
	}
	cols := []string{"k1", "k2", "k3", "v1"}
	for i, want := range cols {
		if stmt.Columns[i].Name != want {
			t.Errorf("Columns[%d].Name = %q, want %q", i, stmt.Columns[i].Name, want)
		}
	}
}

func TestCreateView_AllColumnComments(t *testing.T) {
	sql := `CREATE VIEW example_db.example_view(k1 COMMENT "first key", k2 COMMENT "second key", k3 COMMENT "third key", v1 COMMENT "first value") COMMENT "my first view" AS SELECT c1 AS k1, k2, k3, SUM(v1) FROM example_table WHERE k1 = 20160112 GROUP BY k1, k2, k3`
	stmt := parseCreateViewStmt(t, sql)
	if stmt.Name.String() != "example_db.example_view" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "example_db.example_view")
	}
	if len(stmt.Columns) != 4 {
		t.Fatalf("Columns: got %d, want 4", len(stmt.Columns))
	}
	if stmt.Columns[0].Comment != "first key" {
		t.Errorf("Columns[0].Comment = %q, want %q", stmt.Columns[0].Comment, "first key")
	}
	if stmt.Columns[1].Comment != "second key" {
		t.Errorf("Columns[1].Comment = %q, want %q", stmt.Columns[1].Comment, "second key")
	}
	if stmt.Comment != "my first view" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "my first view")
	}
}

// ---------------------------------------------------------------------------
// ALTER VIEW
// ---------------------------------------------------------------------------

func TestAlterView_Basic(t *testing.T) {
	stmt := parseAlterViewStmt(t, "ALTER VIEW v AS SELECT id FROM t")
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
	if len(stmt.Columns) != 0 {
		t.Errorf("Columns: got %d, want 0", len(stmt.Columns))
	}
	if stmt.Query == nil {
		t.Fatal("Query is nil")
	}
	if stmt.Query.Tag() != ast.T_SelectStmt {
		t.Errorf("Query.Tag() = %v, want T_SelectStmt", stmt.Query.Tag())
	}
}

func TestAlterView_WithColumns(t *testing.T) {
	stmt := parseAlterViewStmt(t, "ALTER VIEW v (a, b) AS SELECT a, b FROM t")
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns: got %d, want 2", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "a" {
		t.Errorf("Columns[0].Name = %q, want %q", stmt.Columns[0].Name, "a")
	}
	if stmt.Columns[1].Name != "b" {
		t.Errorf("Columns[1].Name = %q, want %q", stmt.Columns[1].Name, "b")
	}
}

// ---------------------------------------------------------------------------
// DROP VIEW
// ---------------------------------------------------------------------------

func TestDropView_Basic(t *testing.T) {
	stmt := parseDropViewStmt(t, "DROP VIEW v")
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropView_IfExists(t *testing.T) {
	stmt := parseDropViewStmt(t, "DROP VIEW IF EXISTS v")
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "v")
	}
}

func TestDropView_QualifiedName(t *testing.T) {
	stmt := parseDropViewStmt(t, "DROP VIEW db.v")
	if stmt.Name.String() != "db.v" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "db.v")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

// ---------------------------------------------------------------------------
// Node tags
// ---------------------------------------------------------------------------

func TestViewNodeTags(t *testing.T) {
	cv := &ast.CreateViewStmt{}
	if cv.Tag() != ast.T_CreateViewStmt {
		t.Errorf("CreateViewStmt.Tag() = %v, want T_CreateViewStmt", cv.Tag())
	}
	av := &ast.AlterViewStmt{}
	if av.Tag() != ast.T_AlterViewStmt {
		t.Errorf("AlterViewStmt.Tag() = %v, want T_AlterViewStmt", av.Tag())
	}
	dv := &ast.DropViewStmt{}
	if dv.Tag() != ast.T_DropViewStmt {
		t.Errorf("DropViewStmt.Tag() = %v, want T_DropViewStmt", dv.Tag())
	}
	vc := &ast.ViewColumn{}
	if vc.Tag() != ast.T_ViewColumn {
		t.Errorf("ViewColumn.Tag() = %v, want T_ViewColumn", vc.Tag())
	}
}

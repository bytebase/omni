package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseInlineExternalTable(t *testing.T) {
	result := ParseAndCheck(t, `
SELECT ext.id, ext.name
FROM EXTERNAL (
  (id NUMBER, name VARCHAR2(100))
  TYPE ORACLE_LOADER
  DEFAULT DIRECTORY data_dir
  ACCESS PARAMETERS (FIELDS TERMINATED BY ',')
  LOCATION ('file.csv')
  REJECT LIMIT UNLIMITED
) ext`)

	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.SelectStmt)
	if stmt.FromClause == nil || stmt.FromClause.Len() != 1 {
		t.Fatalf("FromClause len = %d, want 1", stmt.FromClause.Len())
	}

	ext, ok := stmt.FromClause.Items[0].(*ast.InlineExternalTable)
	if !ok {
		t.Fatalf("FROM item = %T, want *ast.InlineExternalTable", stmt.FromClause.Items[0])
	}
	if ext.Columns == nil || ext.Columns.Len() != 2 {
		t.Fatalf("Columns len = %d, want 2", ext.Columns.Len())
	}
	col0 := ext.Columns.Items[0].(*ast.ColumnDef)
	if col0.Name != "ID" || col0.TypeName == nil {
		t.Fatalf("first column = %#v, want ID NUMBER", col0)
	}
	col1 := ext.Columns.Items[1].(*ast.ColumnDef)
	if col1.Name != "NAME" || col1.TypeName == nil {
		t.Fatalf("second column = %#v, want NAME VARCHAR2(100)", col1)
	}
	if ext.Type != "ORACLE_LOADER" {
		t.Fatalf("Type = %q, want ORACLE_LOADER", ext.Type)
	}
	if ext.Directory != "DATA_DIR" {
		t.Fatalf("Directory = %q, want DATA_DIR", ext.Directory)
	}
	if ext.AccessParams == "" {
		t.Fatal("AccessParams is empty, want captured ACCESS PARAMETERS text")
	}
	if ext.Location != "file.csv" {
		t.Fatalf("Location = %q, want file.csv", ext.Location)
	}
	if _, ok := ext.RejectLimit.(*ast.ColumnRef); !ok {
		t.Fatalf("RejectLimit = %T, want *ast.ColumnRef for UNLIMITED", ext.RejectLimit)
	}
	if ext.Alias == nil || ext.Alias.Name != "EXT" {
		t.Fatalf("Alias = %#v, want EXT", ext.Alias)
	}
}

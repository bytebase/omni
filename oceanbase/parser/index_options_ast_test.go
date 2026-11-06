package parser

import (
	"testing"

	ast "github.com/bytebase/omni/oceanbase/ast"
)

// TestIndexOptionsAST: each option lands in the list with its value, and the other options
// around it are not lost.
func TestIndexOptionsAST(t *testing.T) {
	sql := "CREATE INDEX i ON t (a) PARALLEL = 4 STORING (b, `c`) NOPARALLEL WITH ROWID MAX_USED_PART_ID 7 STORING (d) COMMENT 'x'"
	stmt := oneNode(t, sql).(*ast.CreateIndexStmt)
	names := make([]string, 0, len(stmt.Options))
	for _, opt := range stmt.Options {
		names = append(names, opt.Name)
	}
	want := []string{"PARALLEL", "STORING", "NOPARALLEL", "WITH ROWID", "MAX_USED_PART_ID", "STORING", "COMMENT"}
	if len(names) != len(want) {
		t.Fatalf("options %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("options %v, want %v", names, want)
		}
	}
	if v, ok := stmt.Options[0].Value.(*ast.IntLit); !ok || v.Value != 4 {
		t.Errorf("PARALLEL: %+v", stmt.Options[0])
	}
	if cols := stmt.Options[1].Columns; len(cols) != 2 || cols[0] != "b" || cols[1] != "c" {
		t.Errorf("STORING: %+v", stmt.Options[1])
	}
	if v, ok := stmt.Options[4].Value.(*ast.IntLit); !ok || v.Value != 7 {
		t.Errorf("MAX_USED_PART_ID: %+v", stmt.Options[4])
	}
	if got := sql[stmt.Options[1].Loc.Start:stmt.Options[1].Loc.End]; got != "STORING (b, `c`)" {
		t.Errorf("STORING Loc spans %q", got)
	}
	if got := sql[stmt.Options[3].Loc.Start:stmt.Options[3].Loc.End]; got != "WITH ROWID" {
		t.Errorf("WITH ROWID Loc spans %q", got)
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Loc.End=%d, want %d", stmt.Loc.End, len(sql))
	}
}

func TestStorageCachePolicyAST(t *testing.T) {
	global := oneNode(t, "CREATE INDEX i ON t (a) STORAGE_CACHE_POLICY (GLOBAL = 'hot')").(*ast.CreateIndexStmt).Options[0]
	if global.Name != "STORAGE_CACHE_POLICY" || len(global.Properties) != 1 || global.Properties[0].Name != "global" {
		t.Fatalf("GLOBAL: %+v", global)
	}
	if v, ok := global.Properties[0].Value.(*ast.StringLit); !ok || v.Value != "hot" {
		t.Errorf("GLOBAL value: %+v", global.Properties[0])
	}
	list := oneNode(t, "CREATE INDEX i ON t (a) STORAGE_CACHE_POLICY = (BOUNDARY_COLUMN = a, HOT_RETENTION = 1 DAY)").(*ast.CreateIndexStmt).Options[0]
	if len(list.Properties) != 2 {
		t.Fatalf("list: %+v", list.Properties)
	}
	if v, ok := list.Properties[1].Value.(*ast.StringLit); list.Properties[1].Name != "hot_retention" || !ok || v.Value != "1 DAY" {
		t.Errorf("HOT_RETENTION: %+v", list.Properties[1])
	}
}

// TestQualifiedIndexNameAST: the parts of the name land apart and IndexName stays the name.
func TestQualifiedIndexNameAST(t *testing.T) {
	for _, tc := range []struct{ sql, catalog, schema, name, link string }{
		{"CREATE INDEX i ON t (a)", "", "", "i", ""},
		{"CREATE INDEX d.i ON t (a)", "", "d", "i", ""},
		{"CREATE INDEX `d`.`i` ON t (a)", "", "d", "i", ""},
		{"CREATE INDEX c.d.i ON t (a)", "c", "d", "i", ""},
		{"CREATE INDEX d.select ON t (a)", "", "d", "select", ""},
		{"CREATE INDEX select.i ON t (a)", "", "select", "i", ""},
		{"CREATE INDEX i@l ON t (a)", "", "", "i", "l"},
		{"CREATE INDEX d.i@l ON t (a)", "", "d", "i", "l"},
	} {
		stmt := oneNode(t, tc.sql).(*ast.CreateIndexStmt)
		if stmt.IndexCatalog != tc.catalog || stmt.IndexSchema != tc.schema || stmt.IndexName != tc.name || stmt.IndexDBLink != tc.link {
			t.Errorf("%q: catalog %q schema %q name %q link %q", tc.sql, stmt.IndexCatalog, stmt.IndexSchema, stmt.IndexName, stmt.IndexDBLink)
		}
		if stmt.Table == nil || stmt.Table.Name != "t" {
			t.Errorf("%q: the table was lost: %+v", tc.sql, stmt.Table)
		}
	}
}

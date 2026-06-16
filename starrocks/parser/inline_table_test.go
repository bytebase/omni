package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// mustParseInlineTable parses input, expects exactly one FROM item which must
// be an *ast.InlineTable, and returns it.
func mustParseInlineTable(t *testing.T, input string) *ast.InlineTable {
	t.Helper()
	stmt := mustParseSelect(t, input)
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	it, ok := stmt.From[0].(*ast.InlineTable)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.InlineTable", stmt.From[0])
	}
	return it
}

// TestInlineTableValues covers the headline form: a VALUES table constructor
// with an alias and a column-alias list.
func TestInlineTableValues(t *testing.T) {
	it := mustParseInlineTable(t, "SELECT * FROM (VALUES (1, 'a'), (2, 'b')) AS v(id, name)")
	if len(it.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(it.Rows))
	}
	if len(it.Rows[0]) != 2 {
		t.Errorf("row[0] cols = %d, want 2", len(it.Rows[0]))
	}
	if it.Alias != "v" {
		t.Errorf("alias = %q, want v", it.Alias)
	}
	if got := it.ColumnAliases; len(got) != 2 || got[0] != "id" || got[1] != "name" {
		t.Errorf("columnAliases = %v, want [id name]", got)
	}
}

// TestInlineTableNoAlias confirms the alias is optional in FROM position
// (unlike a bare top-level VALUES, which StarRocks does not accept).
func TestInlineTableNoAlias(t *testing.T) {
	it := mustParseInlineTable(t, "SELECT * FROM (VALUES (1), (2))")
	if len(it.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(it.Rows))
	}
	if it.Alias != "" {
		t.Errorf("alias = %q, want empty", it.Alias)
	}
	if it.ColumnAliases != nil {
		t.Errorf("columnAliases = %v, want nil", it.ColumnAliases)
	}
}

// TestInlineTableAliasNoColumns covers an alias without a column-alias list.
func TestInlineTableAliasNoColumns(t *testing.T) {
	it := mustParseInlineTable(t, "SELECT * FROM (VALUES (1, 'a'), (2, 'b')) v")
	if it.Alias != "v" {
		t.Errorf("alias = %q, want v", it.Alias)
	}
	if it.ColumnAliases != nil {
		t.Errorf("columnAliases = %v, want nil", it.ColumnAliases)
	}
}

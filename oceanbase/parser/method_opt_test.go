package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase's grammar reaches method_opt as a statement of its own (stmt:
// method_opt), so DBMS_STATS' statistics method parses as one:
//
//	FOR ALL [INDEXED | HIDDEN] COLUMNS [SIZE {AUTO | REPEAT | SKEWONLY | n}]
//	FOR COLUMNS [column [{,} column] ... [SIZE ...]]
//
// Every case below is held to the pinned server by the method_opt,
// method_list, method, for_all, opt_indexed_hiddden, opt_size_clause,
// size_clause, for_columns, for_columns_list and for_columns_item rows of
// testdata/compat.tsv, whose carriers are these statements.

// TestMethodOptForms pins which method_opt statements parse.
func TestMethodOptForms(t *testing.T) {
	for _, sql := range []string{
		"FOR COLUMNS",
		"FOR COLUMNS a",
		"FOR COLUMNS a a",
		"FOR COLUMNS a, a",
		"FOR COLUMNS SIZE AUTO",
		"FOR COLUMNS SIZE REPEAT",
		"FOR COLUMNS SIZE SKEWONLY",
		"FOR COLUMNS SIZE 1",
		"FOR COLUMNS SIZE 1.5",
		"FOR COLUMNS a SIZE AUTO",
		"FOR COLUMNS a SIZE AUTO zzz",
		"FOR COLUMNS a SKEWONLY AUTO",
		"FOR COLUMNS AUTO",
		"FOR COLUMNS (a, b) SIZE 10",
		"FOR ALL COLUMNS",
		"FOR ALL COLUMNS SIZE AUTO",
		"FOR ALL INDEXED COLUMNS SIZE AUTO",
		"FOR ALL HIDDEN COLUMNS",
		"FOR ALL COLUMNS SIZE AUTO FOR COLUMNS a SIZE 1",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"FOR",
		"FOR ALL",
		"FOR ALL COLUMNS AUTO", // a size_clause needs its SIZE
		"FOR ALL COLUMNS SIZE",
		"FOR ALL INDEXED",
		"FOR COLUMNS 1",          // no size_clause and no column: 1 is not a name
		"FOR COLUMNS SIZE",       // SIZE with nothing after it
		"FOR COLUMNS a SIZE",     // the same, after a column
		"FOR COLUMNS , a",        // the list never opens with a comma
		"FOR COLUMNS a,",         // nor ends with one
		"FOR COLUMNS REPEAT",     // REPEAT is reserved, so it is no column name
		"FOR COLUMNS SKEWONLY 1", // a number is no column
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}

// TestMethodOptTree pins what a method_opt statement parses into: its methods,
// the size_clause words, and the Location the statement covers.
func TestMethodOptTree(t *testing.T) {
	list, err := Parse("FOR ALL INDEXED COLUMNS SIZE AUTO FOR COLUMNS a SIZE 1, (b, c) SIZE SKEWONLY")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("got %d statements, want 1", len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.MethodOptStmt)
	if !ok {
		t.Fatalf("got %T, want *nodes.MethodOptStmt", list.Items[0])
	}
	if len(stmt.Methods) != 2 {
		t.Fatalf("got %d methods, want 2", len(stmt.Methods))
	}

	all := stmt.Methods[0]
	if !all.ForAll || all.Indexed != "INDEXED" || all.Size != "AUTO" || len(all.Items) != 0 {
		t.Errorf("FOR ALL INDEXED COLUMNS SIZE AUTO = %#v, want ForAll, INDEXED, AUTO and no items", all)
	}

	cols := stmt.Methods[1]
	if cols.ForAll || len(cols.Items) != 2 {
		t.Fatalf("FOR COLUMNS ... = %#v, want two items", cols)
	}
	first, second := cols.Items[0], cols.Items[1]
	if len(first.Columns) != 1 || first.Columns[0] != "a" || first.Size != "1" {
		t.Errorf("first item = %#v, want column a of size 1", first)
	}
	if len(second.Columns) != 2 || second.Columns[0] != "b" || second.Columns[1] != "c" || second.Size != "SKEWONLY" {
		t.Errorf("second item = %#v, want extension (b, c) of size SKEWONLY", second)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len("FOR ALL INDEXED COLUMNS SIZE AUTO FOR COLUMNS a SIZE 1, (b, c) SIZE SKEWONLY") {
		t.Errorf("Loc = %+v, want the whole statement", stmt.Loc)
	}
}

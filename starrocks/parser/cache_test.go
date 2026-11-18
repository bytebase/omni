package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// TestDataCacheSelect covers dataCacheSelectStatement (StarRocks.g4:1867).
func TestDataCacheSelect(t *testing.T) {
	file, errs := parseForTest("CACHE SELECT a, b AS c FROM db.t WHERE a = 1 PROPERTIES('x' = 'y')")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.DataCacheSelectStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.DataCacheSelectStmt", file.Stmts[0])
	}
	if len(stmt.Items) != 2 || stmt.Items[1].Alias != "c" {
		t.Errorf("Items = %+v, want two items with the second aliased c", stmt.Items)
	}
	if stmt.From == nil || stmt.From.String() != "db.t" {
		t.Errorf("From = %v, want db.t", stmt.From)
	}
	if stmt.Where == nil {
		t.Error("Where is nil")
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties = %+v, want one pair", stmt.Properties)
	}
}

func TestDataCacheSelectRejects(t *testing.T) {
	for _, sql := range []string{
		"CACHE SELECT FROM a",
		"CACHE SELECT a",
		"CACHE SELECT a FROM",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

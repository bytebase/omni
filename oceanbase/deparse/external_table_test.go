package deparse

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

// TestDeparseExternalTableExprs: a FROM holds what OceanBase reads, so a table
// read from a location, RB_ITERATE, UNNEST and a derived table's flashback all
// survive a deparse, and what comes out parses and deparses to itself.
func TestDeparseExternalTableExprs(t *testing.T) {
	for _, c := range []struct{ input, expected string }{
		{"SELECT a FROM SOURCE(TYPE = 's')", "select `a` AS `a` from source(type = 's')"},
		{"SELECT a FROM SOURCE(TYPE = 's', COMPRESSION_CODE = a) AS x", "select `a` AS `a` from source(type = 's', compression_code = `a`) as `x`"},
		{"SELECT a FROM 's'(FORMAT(TYPE = 'csv') PATTERN 'p')", "select `a` AS `a` from 's'(format(type = 'csv') pattern 'p')"},
		{"SELECT a FROM FILES(LOCATION 's', FORMAT(TYPE = 'csv'), PATTERN 'p') AS f", "select `a` AS `a` from files(location 's', format(type = 'csv'), pattern 'p') as `f`"},
		{"SELECT a FROM 's'(FORMAT(TYPE = 'csv') PATTERN '')", "select `a` AS `a` from 's'(format(type = 'csv') pattern '')"},
		{"SELECT a FROM 's'(FORMAT(TYPE = 'csv', NULL_IF = ('a', 'b'), COLUMN_BLOOM_FILTER = (1, 2)))", "select `a` AS `a` from 's'(format(type = 'csv', null_if = ('a', 'b'), column_bloom_filter = (1, 2)))"},
		{"SELECT a FROM (SELECT a) AS OF SNAPSHOT 1 t", "select `a` AS `a` from (select `a` AS `a`) as of snapshot 1 `t`"},
		{"SELECT a FROM RB_ITERATE(a) AS b(c)", "select `a` AS `a` from rb_iterate(`a`) `b`(`c`)"},
		{"SELECT a FROM UNNEST(a, b) AS u(c, d)", "select `a` AS `a` from unnest(`a`, `b`) `u`(`c`, `d`)"},
	} {
		list, err := parser.Parse(c.input)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.input, err)
			continue
		}
		got := DeparseSelect(list.Items[0].(*ast.SelectStmt))
		if got != c.expected {
			t.Errorf("DeparseSelect(%q)\n got %s\nwant %s", c.input, got, c.expected)
		}
		again, err := parser.Parse(got)
		if err != nil {
			t.Errorf("Parse(%q), the deparse of %q: %v", got, c.input, err)
			continue
		}
		if twice := DeparseSelect(again.Items[0].(*ast.SelectStmt)); twice != got {
			t.Errorf("deparse of %q is not stable:\n once  %s\n twice %s", c.input, got, twice)
		}
	}
}

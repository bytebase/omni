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
		// A side with no name keeps its USING rather than writing an empty
		// qualifier.
		{"SELECT t.a FROM t JOIN UNNEST(t.a) USING (a)", "select `t`.`a` AS `a` from (`t` join unnest(`t`.`a`) using(`a`))"},
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

// TestResolverStarOverExternalSource: a star over a source whose columns are
// known only when the query runs stays a star, instead of expanding to nothing.
func TestResolverStarOverExternalSource(t *testing.T) {
	const sql = "SELECT * FROM SOURCE(TYPE = 's') AS e"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	RewriteSelectStmt(sel)
	resolver := &Resolver{Lookup: func(string) *ResolverTable { return nil }}
	resolved := resolver.Resolve(sel)
	if len(resolved.TargetList) != 1 {
		t.Fatalf("got %d targets, want the star kept as one", len(resolved.TargetList))
	}
	if _, ok := resolved.TargetList[0].(*ast.StarExpr); !ok {
		t.Errorf("target = %T, want the star kept", resolved.TargetList[0])
	}
}

// TestResolverStarThroughDerivedSource: a derived table over a source with
// unknown columns cannot expand d.* either, so the star stays.
func TestResolverStarThroughDerivedSource(t *testing.T) {
	const sql = "SELECT d.* FROM (SELECT * FROM SOURCE(TYPE = 's')) AS d"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sel := list.Items[0].(*ast.SelectStmt)
	RewriteSelectStmt(sel)
	resolver := &Resolver{Lookup: func(string) *ResolverTable { return nil }}
	resolved := resolver.Resolve(sel)
	if len(resolved.TargetList) != 1 {
		t.Fatalf("got %d targets, want d.* kept as one", len(resolved.TargetList))
	}
	col, ok := resolved.TargetList[0].(*ast.ColumnRef)
	if !ok || !col.Star {
		t.Errorf("target = %T, want the qualified star kept", resolved.TargetList[0])
	}
}

package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAnalyzeColumnScopeOptionsParse(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		option string
	}{
		{"ANALYZE venue PREDICATE COLUMNS", "predicate_columns"},
		{"ANALYZE listing ALL COLUMNS", "all_columns"},
		{"ANALYZE VERBOSE public.listing PREDICATE COLUMNS", "predicate_columns"},
	} {
		stmt := singleStmt(t, tc.sql)
		analyze, ok := stmt.(*nodes.VacuumStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected VacuumStmt, got %T", tc.sql, stmt)
		}
		if analyze.IsVacuumCmd {
			t.Fatalf("Parse(%q): expected ANALYZE VacuumStmt, got vacuum", tc.sql)
		}
		if findDefElem(analyze.Options, tc.option) == nil {
			t.Fatalf("Parse(%q): expected option %q in %#v", tc.sql, tc.option, analyze.Options)
		}
	}
}

func TestRedshiftAnalyzeCompressionParse(t *testing.T) {
	for _, tc := range []struct {
		sql          string
		relname      string
		wantComprows int64
	}{
		{"ANALYZE COMPRESSION", "", 0},
		{"ANALYZE COMPRESSION listing", "listing", 0},
		{"ANALYZE COMPRESSION public.sales COMPROWS 200000", "sales", 200000},
		{"ANALYZE COMPRESSION venue(venuename, venuecity) COMPROWS 50000", "venue", 50000},
	} {
		stmt := singleStmt(t, tc.sql)
		analyze, ok := stmt.(*nodes.VacuumStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected VacuumStmt, got %T", tc.sql, stmt)
		}
		if findDefElem(analyze.Options, "compression") == nil {
			t.Fatalf("Parse(%q): expected compression option in %#v", tc.sql, analyze.Options)
		}
		if tc.wantComprows > 0 {
			assertDefElemInt(t, analyze.Options, "comprows", tc.wantComprows)
		}
		if tc.relname == "" {
			if analyze.Rels != nil {
				t.Fatalf("Parse(%q): rels = %#v, want nil", tc.sql, analyze.Rels)
			}
			continue
		}
		if analyze.Rels == nil || len(analyze.Rels.Items) != 1 {
			t.Fatalf("Parse(%q): rels = %#v, want one relation", tc.sql, analyze.Rels)
		}
		rel, ok := analyze.Rels.Items[0].(*nodes.VacuumRelation)
		if !ok {
			t.Fatalf("Parse(%q): relation = %T, want VacuumRelation", tc.sql, analyze.Rels.Items[0])
		}
		if rel.Relation == nil || rel.Relation.Relname != tc.relname {
			t.Fatalf("Parse(%q): relname = %#v, want %q", tc.sql, rel.Relation, tc.relname)
		}
	}
}

// TestRedshiftAnalyzeCompressionColumnListParse pins r_ANALYZE_COMPRESSION's
// column list, which the docs leave outside the optional table name:
// ANALYZE COMPRESSION (a) names columns with no table to hold them.
func TestRedshiftAnalyzeCompressionColumnListParse(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		rel    string
		column int
	}{
		{"ANALYZE COMPRESSION(a)", "", 1},
		{"ANALYZE COMPRESSION(a, b)", "", 2},
		{"ANALYZE COMPRESSION venue(venuename, venuecity)", "venue", 2},
		{"ANALYZE COMPRESSION(a) COMPROWS 100", "", 1},
	} {
		stmt := singleStmt(t, tc.sql)
		analyze, ok := stmt.(*nodes.VacuumStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected VacuumStmt, got %T", tc.sql, stmt)
		}
		if findDefElem(analyze.Options, "compression") == nil {
			t.Fatalf("Parse(%q): expected compression option in %#v", tc.sql, analyze.Options)
		}
		if analyze.Rels == nil || len(analyze.Rels.Items) != 1 {
			t.Fatalf("Parse(%q): rels = %#v, want one relation", tc.sql, analyze.Rels)
		}
		rel, ok := analyze.Rels.Items[0].(*nodes.VacuumRelation)
		if !ok {
			t.Fatalf("Parse(%q): relation = %T, want VacuumRelation", tc.sql, analyze.Rels.Items[0])
		}
		if tc.rel == "" {
			if rel.Relation != nil {
				t.Fatalf("Parse(%q): relation name = %#v, want none", tc.sql, rel.Relation)
			}
		} else if rel.Relation == nil || rel.Relation.Relname != tc.rel {
			t.Fatalf("Parse(%q): relation name = %#v, want %q", tc.sql, rel.Relation, tc.rel)
		}
		if rel.VaCols == nil || len(rel.VaCols.Items) != tc.column {
			t.Fatalf("Parse(%q): columns = %#v, want %d", tc.sql, rel.VaCols, tc.column)
		}
	}
}

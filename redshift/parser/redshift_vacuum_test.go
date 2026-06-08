package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftVacuumModesParse(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		mode      string
		relname   string
		threshold int64
		boost     bool
	}{
		{"VACUUM FULL sales TO 75 PERCENT", "full", "sales", 75, false},
		{"VACUUM SORT ONLY sales TO 75 PERCENT BOOST", "sort_only", "sales", 75, true},
		{"VACUUM DELETE ONLY inventory_table", "delete_only", "inventory_table", 0, false},
		{"VACUUM REINDEX listing TO 80 PERCENT", "reindex", "listing", 80, false},
		{"VACUUM RECLUSTER listing BOOST", "recluster", "listing", 0, true},
		{"VACUUM TO 100 PERCENT", "", "", 100, false},
		{"VACUUM BOOST", "", "", 0, true},
		{"VACUUM SORT ONLY", "sort_only", "", 0, false},
	} {
		stmt := singleStmt(t, tc.sql)
		vacuum, ok := stmt.(*nodes.VacuumStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected VacuumStmt, got %T", tc.sql, stmt)
		}
		if !vacuum.IsVacuumCmd {
			t.Fatalf("Parse(%q): expected VACUUM statement", tc.sql)
		}
		if tc.mode != "" {
			assertDefElemString(t, vacuum.Options, "mode", tc.mode)
		}
		if tc.threshold > 0 {
			assertDefElemInt(t, vacuum.Options, "threshold_percent", tc.threshold)
		}
		if tc.boost {
			if findDefElem(vacuum.Options, "boost") == nil {
				t.Fatalf("Parse(%q): expected BOOST option in %#v", tc.sql, vacuum.Options)
			}
		}
		if tc.relname == "" {
			if vacuum.Rels != nil {
				t.Fatalf("Parse(%q): rels = %#v, want nil", tc.sql, vacuum.Rels)
			}
			continue
		}
		if vacuum.Rels == nil || len(vacuum.Rels.Items) != 1 {
			t.Fatalf("Parse(%q): rels = %#v, want one relation", tc.sql, vacuum.Rels)
		}
		rel, ok := vacuum.Rels.Items[0].(*nodes.VacuumRelation)
		if !ok {
			t.Fatalf("Parse(%q): relation = %T, want VacuumRelation", tc.sql, vacuum.Rels.Items[0])
		}
		if rel.Relation == nil || rel.Relation.Relname != tc.relname {
			t.Fatalf("Parse(%q): relname = %#v, want %q", tc.sql, rel.Relation, tc.relname)
		}
	}
}

func TestRedshiftVacuumQualifiedAndQuotedRelationParse(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		relname string
	}{
		{`VACUUM "MixedCase_Table"`, "MixedCase_Table"},
		{`VACUUM FULL "schema"."table"`, "table"},
	} {
		stmt := singleStmt(t, tc.sql)
		vacuum, ok := stmt.(*nodes.VacuumStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected VacuumStmt, got %T", tc.sql, stmt)
		}
		if vacuum.Rels == nil || len(vacuum.Rels.Items) != 1 {
			t.Fatalf("Parse(%q): rels = %#v, want one relation", tc.sql, vacuum.Rels)
		}
		rel, ok := vacuum.Rels.Items[0].(*nodes.VacuumRelation)
		if !ok {
			t.Fatalf("Parse(%q): relation = %T, want VacuumRelation", tc.sql, vacuum.Rels.Items[0])
		}
		if rel.Relation == nil || rel.Relation.Relname != tc.relname {
			t.Fatalf("Parse(%q): relname = %#v, want %q", tc.sql, rel.Relation, tc.relname)
		}
	}
}

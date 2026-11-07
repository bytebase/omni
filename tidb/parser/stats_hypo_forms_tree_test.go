package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What the extended-statistics, hypothetical-index and TiFlash-replica
// extensions leave in the tree. The corpus records only whether each parses;
// these read back the fields, so a HYPO flag dropped, a stats type read as its
// neighbour or a column list lost would fail here even while the corpus stays
// green.

func parseOneAlterCmd(t *testing.T, sql string) *ast.AlterTableCmd {
	t.Helper()
	stmt, ok := parseOneStmt(t, sql).(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave a statement of another kind", sql)
	}
	if len(stmt.Commands) != 1 {
		t.Fatalf("Parse(%q) gave %d commands, want 1", sql, len(stmt.Commands))
	}
	return stmt.Commands[0]
}

// TestSetHypoTiFlashReplicaTree reads back the HYPO modifier on SET TIFLASH
// REPLICA. HYPO pins a hypothetical replica for the optimizer; the plain form
// must leave it off.
func TestSetHypoTiFlashReplicaTree(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		count  int
		hypo   bool
		labels []string
	}{
		{sql: "alter table t set hypo tiflash replica 1", count: 1, hypo: true},
		{sql: "alter table t set tiflash replica 1", count: 1, hypo: false},
		{sql: "alter table t set hypo tiflash replica 2 location labels 'z'", count: 2, hypo: true, labels: []string{"z"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			cmd := parseOneAlterCmd(t, tc.sql)
			if cmd.Type != ast.ATSetTiFlashReplica {
				t.Fatalf("Type = %v, want ATSetTiFlashReplica", cmd.Type)
			}
			if cmd.Hypo != tc.hypo {
				t.Errorf("Hypo = %v, want %v", cmd.Hypo, tc.hypo)
			}
			if cmd.TiFlashReplica != tc.count {
				t.Errorf("TiFlashReplica = %d, want %d", cmd.TiFlashReplica, tc.count)
			}
			wantList(t, "Labels", cmd.TiFlashLocationLabels, tc.labels)
		})
	}
}

// TestAddStatsExtendedTree reads back the name, the type and the column list
// of ADD STATS_EXTENDED. The three types are three values and the columns may
// be qualified, as in CREATE STATISTICS.
func TestAddStatsExtendedTree(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		ifNotExists bool
		name        string
		typ         ast.StatisticsType
		columns     []string
	}{
		{sql: "alter table t add stats_extended s cardinality(a)", name: "s", typ: ast.StatisticsCardinality, columns: []string{"a"}},
		{sql: "alter table t add stats_extended s dependency(a,b)", name: "s", typ: ast.StatisticsDependency, columns: []string{"a", "b"}},
		{sql: "alter table t add stats_extended s correlation(a, b)", name: "s", typ: ast.StatisticsCorrelation, columns: []string{"a", "b"}},
		{sql: "alter table t add stats_extended if not exists s cardinality(a)", ifNotExists: true, name: "s", typ: ast.StatisticsCardinality, columns: []string{"a"}},
		{sql: "alter table t add stats_extended s cardinality(t.a)", name: "s", typ: ast.StatisticsCardinality, columns: []string{"t.a"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			cmd := parseOneAlterCmd(t, tc.sql)
			if cmd.Type != ast.ATAddStatistics {
				t.Fatalf("Type = %v, want ATAddStatistics", cmd.Type)
			}
			if cmd.IfNotExists != tc.ifNotExists {
				t.Errorf("IfNotExists = %v, want %v", cmd.IfNotExists, tc.ifNotExists)
			}
			if cmd.StatsName != tc.name {
				t.Errorf("StatsName = %q, want %q", cmd.StatsName, tc.name)
			}
			if cmd.StatsType != tc.typ {
				t.Errorf("StatsType = %s, want %s", cmd.StatsType, tc.typ)
			}
			var cols []string
			for _, c := range cmd.StatsColumns {
				var b strings.Builder
				if c.Schema != "" {
					b.WriteString(c.Schema + ".")
				}
				if c.Table != "" {
					b.WriteString(c.Table + ".")
				}
				b.WriteString(c.Column)
				cols = append(cols, b.String())
			}
			wantList(t, "StatsColumns", cols, tc.columns)
		})
	}
}

// TestDropStatsExtendedTree reads back the name and the IF EXISTS flag of DROP
// STATS_EXTENDED.
func TestDropStatsExtendedTree(t *testing.T) {
	for _, tc := range []struct {
		sql      string
		name     string
		ifExists bool
	}{
		{sql: "alter table t drop stats_extended s", name: "s"},
		{sql: "alter table t drop stats_extended if exists s", name: "s", ifExists: true},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			cmd := parseOneAlterCmd(t, tc.sql)
			if cmd.Type != ast.ATDropStatistics {
				t.Fatalf("Type = %v, want ATDropStatistics", cmd.Type)
			}
			if cmd.StatsName != tc.name {
				t.Errorf("StatsName = %q, want %q", cmd.StatsName, tc.name)
			}
			if cmd.IfExists != tc.ifExists {
				t.Errorf("IfExists = %v, want %v", cmd.IfExists, tc.ifExists)
			}
		})
	}
}

// TestDropHypoIndexTree reads back the HYPO flag on DROP INDEX, the name and
// the table. The plain form must leave IsHypo off.
func TestDropHypoIndexTree(t *testing.T) {
	for _, tc := range []struct {
		sql      string
		isHypo   bool
		name     string
		table    string
		ifExists bool
	}{
		{sql: "drop hypo index i on t", isHypo: true, name: "i", table: "t"},
		{sql: "drop hypo index if exists i on db.t", isHypo: true, name: "i", table: "db.t", ifExists: true},
		{sql: "drop index i on t", isHypo: false, name: "i", table: "t"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, tc.sql).(*ast.DropIndexStmt)
			if !ok {
				t.Fatalf("Parse(%q) gave a statement of another kind", tc.sql)
			}
			if stmt.IsHypo != tc.isHypo {
				t.Errorf("IsHypo = %v, want %v", stmt.IsHypo, tc.isHypo)
			}
			if stmt.Name != tc.name {
				t.Errorf("Name = %q, want %q", stmt.Name, tc.name)
			}
			if stmt.IfExists != tc.ifExists {
				t.Errorf("IfExists = %v, want %v", stmt.IfExists, tc.ifExists)
			}
			wantList(t, "Table", tableNames([]*ast.TableRef{stmt.Table}), []string{tc.table})
		})
	}
}

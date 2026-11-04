package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What PLAN REPLAYER leaves in the tree.
//
// The corpus records only whether a statement parses, and this statement is
// where that matters most in the family: four shapes share one node, and
// nothing in the corpus says which shape a row landed as. A CAPTURE whose two
// digests came out swapped, a LOAD read as a DUMP of a file, or a SLOW QUERY
// read as a dump of nothing would each parse every corpus row. So would an
// Analyze that failed to set -- and that one is not cosmetic: Analyze is what
// says whether the server runs the statement the dump carries, which the
// catalog applies and review's walk follows.

func parseOnePlanReplayer(t *testing.T, sql string) *ast.PlanReplayerStmt {
	t.Helper()
	stmt, ok := parseOneStmt(t, sql).(*ast.PlanReplayerStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave a statement of another kind", sql)
	}
	return stmt
}

// shape is every flag and string of a PlanReplayerStmt, so a reading that
// confused two arms shows up as a different shape rather than as a nil check
// that still passes.
type planReplayerShape struct {
	analyze    bool
	slowQuery  bool
	load       bool
	capture    bool
	remove     bool
	file       string
	sqlDigest  string
	planDigest string
	hasStmt    bool
	hasStats   bool
}

func shapeOf(s *ast.PlanReplayerStmt) planReplayerShape {
	return planReplayerShape{
		analyze:    s.Analyze,
		slowQuery:  s.SlowQuery,
		load:       s.Load,
		capture:    s.Capture,
		remove:     s.Remove,
		file:       s.File,
		sqlDigest:  s.SQLDigest,
		planDigest: s.PlanDigest,
		hasStmt:    s.Stmt != nil,
		hasStats:   s.HistoricalStats != nil,
	}
}

// TestPlanReplayerShapes pins each arm's whole shape by value. Every field is
// compared, so a flag that leaked from one arm into another fails here even
// when the field the test was written for is right.
func TestPlanReplayerShapes(t *testing.T) {
	tests := []struct {
		sql  string
		want planReplayerShape
	}{
		{"plan replayer dump explain select 1", planReplayerShape{hasStmt: true}},
		{"plan replayer dump explain analyze select 1", planReplayerShape{analyze: true, hasStmt: true}},
		{"plan replayer dump explain alter table t add column c int", planReplayerShape{hasStmt: true}},
		{"plan replayer dump explain analyze alter table t add column c int", planReplayerShape{analyze: true, hasStmt: true}},
		// SLOW QUERY carries no statement, and every tail of it is
		// optional, so the flag is the only thing that says which arm it is.
		{"plan replayer dump explain slow query", planReplayerShape{slowQuery: true}},
		{"plan replayer dump explain analyze slow query", planReplayerShape{analyze: true, slowQuery: true}},
		{"plan replayer dump explain slow query limit 10", planReplayerShape{slowQuery: true}},
		// The file arm carries a name and no statement.
		{"plan replayer dump explain 'sql.txt'", planReplayerShape{file: "sql.txt"}},
		{"plan replayer dump explain analyze 'sql.txt'", planReplayerShape{analyze: true, file: "sql.txt"}},
		// LOAD takes the same field and is a different arm, which Load is
		// what says.
		{"plan replayer load '/tmp/a.zip'", planReplayerShape{load: true, file: "/tmp/a.zip"}},
		// CAPTURE and its REMOVE form are two flags, never both.
		{"plan replayer capture 'sql1' 'plan1'", planReplayerShape{capture: true, sqlDigest: "sql1", planDigest: "plan1"}},
		{"plan replayer capture remove 'sql1' 'plan1'", planReplayerShape{remove: true, sqlDigest: "sql1", planDigest: "plan1"}},
		// WITH STATS stands in front of EXPLAIN and reaches every arm.
		{"plan replayer dump with stats as of timestamp '2024-01-01' explain select 1", planReplayerShape{hasStmt: true, hasStats: true}},
		{"plan replayer dump with stats as of timestamp now() explain analyze select 1", planReplayerShape{analyze: true, hasStmt: true, hasStats: true}},
		{"plan replayer dump with stats as of timestamp '2024-01-01' explain slow query limit 10", planReplayerShape{slowQuery: true, hasStats: true}},
		{"plan replayer dump with stats as of timestamp '2024-01-01' explain 'sql.txt'", planReplayerShape{file: "sql.txt", hasStats: true}},
	}
	for _, c := range tests {
		t.Run(c.sql, func(t *testing.T) {
			if got := shapeOf(parseOnePlanReplayer(t, c.sql)); got != c.want {
				t.Errorf("shape = %+v, want %+v", got, c.want)
			}
		})
	}
}

// TestPlanReplayerCaptureDigestOrder pins which digest is which. Both are
// string literals in the same arm, so a reading that swapped them parses
// every corpus row -- the corpus writes '123' '123' for both.
func TestPlanReplayerCaptureDigestOrder(t *testing.T) {
	for _, sql := range []string{
		"plan replayer capture 'the-sql-digest' 'the-plan-digest'",
		"plan replayer capture remove 'the-sql-digest' 'the-plan-digest'",
	} {
		t.Run(sql, func(t *testing.T) {
			stmt := parseOnePlanReplayer(t, sql)
			if stmt.SQLDigest != "the-sql-digest" {
				t.Errorf("SQLDigest = %q, want the-sql-digest", stmt.SQLDigest)
			}
			if stmt.PlanDigest != "the-plan-digest" {
				t.Errorf("PlanDigest = %q, want the-plan-digest", stmt.PlanDigest)
			}
		})
	}
}

// TestPlanReplayerCarriesInnerStatement pins the statement a dump holds, by
// the node it is and by what that node says. ALTER TABLE is the case that
// matters: it is the one DDL the dump can carry, so a command list that came
// out empty would have the catalog apply nothing while the server applies the
// change.
func TestPlanReplayerCarriesInnerStatement(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		stmt := parseOnePlanReplayer(t, "plan replayer dump explain select a from t")
		sel, ok := stmt.Stmt.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.SelectStmt", stmt.Stmt)
		}
		if got := exprTableNames(sel.From); join(got) != join([]string{"t"}) {
			t.Errorf("the dumped SELECT reads %s, want [t]", join(got))
		}
	})

	t.Run("alter table", func(t *testing.T) {
		for _, sql := range []string{
			"plan replayer dump explain alter table t add column c int",
			"plan replayer dump explain analyze alter table t add column c int",
		} {
			stmt := parseOnePlanReplayer(t, sql)
			alt, ok := stmt.Stmt.(*ast.AlterTableStmt)
			if !ok {
				t.Fatalf("%q carries %T, want *ast.AlterTableStmt", sql, stmt.Stmt)
			}
			if alt.Table == nil || alt.Table.Name != "t" {
				t.Errorf("%q alters %#v, want table t", sql, alt.Table)
			}
			if len(alt.Commands) != 1 {
				t.Fatalf("%q carries %d commands, want 1", sql, len(alt.Commands))
			}
			if alt.Commands[0].Type != ast.ATAddColumn {
				t.Errorf("%q carries command %v, want ADD COLUMN", sql, alt.Commands[0].Type)
			}
		}
	})

	t.Run("explain analyze of table", func(t *testing.T) {
		// `... EXPLAIN ANALYZE TABLE t` is EXPLAIN ANALYZE of the TABLE
		// statement, not an ANALYZE TABLE: ANALYZE TABLE is not explainable.
		// The two readings parse the same text, and only the shape tells
		// them apart.
		stmt := parseOnePlanReplayer(t, "plan replayer dump explain analyze table t")
		if !stmt.Analyze {
			t.Error("Analyze is false; ANALYZE was read as part of the statement")
		}
		if _, ok := stmt.Stmt.(*ast.TableStmt); !ok {
			t.Errorf("Stmt is %T, want *ast.TableStmt", stmt.Stmt)
		}
	})

	t.Run("dml", func(t *testing.T) {
		for _, c := range []struct {
			sql  string
			want string
		}{
			{"plan replayer dump explain insert into t values (1)", "*ast.InsertStmt"},
			{"plan replayer dump explain replace into t values (1)", "*ast.InsertStmt"},
			{"plan replayer dump explain update t set a = 1", "*ast.UpdateStmt"},
			{"plan replayer dump explain delete from t", "*ast.DeleteStmt"},
			{"plan replayer dump explain table t", "*ast.TableStmt"},
			{"plan replayer dump explain values row(1)", "*ast.ValuesStmt"},
			{"plan replayer dump explain with c as (select 1) select * from c", "*ast.SelectStmt"},
		} {
			stmt := parseOnePlanReplayer(t, c.sql)
			if got := nodeTypeName(stmt.Stmt); got != c.want {
				t.Errorf("%q carries %s, want %s", c.sql, got, c.want)
			}
		}
	})
}

// TestPlanReplayerSlowQueryTail pins the three clauses that narrow a
// slow-query dump, by value. Each is optional, so a reading that dropped one
// or put a LIMIT's two numbers the wrong way round parses every corpus row.
func TestPlanReplayerSlowQueryTail(t *testing.T) {
	t.Run("where order limit", func(t *testing.T) {
		stmt := parseOnePlanReplayer(t, "plan replayer dump explain slow query where a > 10 and t < 1 order by t limit 10")
		if stmt.Where == nil {
			t.Error("WHERE was dropped")
		}
		if len(stmt.OrderBy) != 1 {
			t.Fatalf("ORDER BY has %d items, want 1", len(stmt.OrderBy))
		}
		if col, ok := stmt.OrderBy[0].Expr.(*ast.ColumnRef); !ok || col.Column != "t" {
			t.Errorf("ORDER BY reads %#v, want column t", stmt.OrderBy[0].Expr)
		}
		if stmt.Limit == nil {
			t.Fatal("LIMIT was dropped")
		}
		if got := litInt(t, stmt.Limit.Count); got != 10 {
			t.Errorf("LIMIT count = %d, want 10", got)
		}
		if stmt.Limit.Offset != nil {
			t.Errorf("LIMIT offset = %#v, want none", stmt.Limit.Offset)
		}
	})

	t.Run("limit offset, count", func(t *testing.T) {
		// LIMIT 1, 10 is offset 1 count 10, as it is in a SELECT.
		stmt := parseOnePlanReplayer(t, "plan replayer dump explain slow query limit 1, 10")
		if stmt.Limit == nil {
			t.Fatal("LIMIT was dropped")
		}
		if got := litInt(t, stmt.Limit.Count); got != 10 {
			t.Errorf("LIMIT count = %d, want 10", got)
		}
		if got := litInt(t, stmt.Limit.Offset); got != 1 {
			t.Errorf("LIMIT offset = %d, want 1", got)
		}
	})

	t.Run("limit count offset n", func(t *testing.T) {
		stmt := parseOnePlanReplayer(t, "plan replayer dump explain slow query where a > 10 limit 10 offset 5")
		if got := litInt(t, stmt.Limit.Count); got != 10 {
			t.Errorf("LIMIT count = %d, want 10", got)
		}
		if got := litInt(t, stmt.Limit.Offset); got != 5 {
			t.Errorf("LIMIT offset = %d, want 5", got)
		}
	})

	t.Run("no tail", func(t *testing.T) {
		stmt := parseOnePlanReplayer(t, "plan replayer dump explain slow query")
		if stmt.Where != nil || stmt.OrderBy != nil || stmt.Limit != nil {
			t.Errorf("a bare SLOW QUERY carries where=%#v order=%#v limit=%#v, want none of the three",
				stmt.Where, stmt.OrderBy, stmt.Limit)
		}
	})
}

// TestPlanReplayerHistoricalStats pins the expression WITH STATS AS OF
// TIMESTAMP holds. It is the only place the clause's value lands, so a
// reading that consumed the expression without keeping it would parse every
// form in the forms test.
func TestPlanReplayerHistoricalStats(t *testing.T) {
	stmt := parseOnePlanReplayer(t, "plan replayer dump with stats as of timestamp '2024-01-01 00:00:00' explain select 1")
	lit, ok := stmt.HistoricalStats.(*ast.StringLit)
	if !ok {
		t.Fatalf("HistoricalStats is %T, want *ast.StringLit", stmt.HistoricalStats)
	}
	if lit.Value != "2024-01-01 00:00:00" {
		t.Errorf("HistoricalStats = %q, want 2024-01-01 00:00:00", lit.Value)
	}
}

// TestPlanReplayerLoc pins the statement's own byte range: the whole
// statement, from PLAN to the last token of what it carries.
func TestPlanReplayerLoc(t *testing.T) {
	const sql = "plan replayer dump explain analyze select 1"
	stmt := parseOnePlanReplayer(t, sql)
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(sql))
	}
	inner := stmt.Stmt.(*ast.SelectStmt)
	if got := sql[inner.Loc.Start:inner.Loc.End]; got != "select 1" {
		t.Errorf("the dumped statement spans %q, want %q", got, "select 1")
	}
}

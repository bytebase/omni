package parser

import (
	"context"
	"testing"
)

// TestJoinGrammarOceanBase pins the join-chain rules this parser follows against
// the OceanBase engine: an outer (LEFT/RIGHT) join requires an ON/USING
// condition, a NATURAL join forbids one, and a conditionless plain/inner/outer
// join lets its right operand absorb the joins that follow, so a trailing ON
// binds to the outer join (`a JOIN b JOIN c ON x ON y`).
//
// Unlike MySQL and TiDB, OceanBase also lets STRAIGHT_JOIN absorb its right
// operand, so `a STRAIGHT_JOIN b JOIN c ON x ON y` parses; this is the one join
// rule where omni-oceanbase diverges from its MySQL ancestry, and the engine
// decides it (read 2026-11-08).
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestJoinGrammarOceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS objoin",
		"USE objoin",
		"DROP TABLE IF EXISTS t1, t2, t3",
		"CREATE TABLE t1 (id INT, a INT, b INT)",
		"CREATE TABLE t2 (id INT, a INT, b INT)",
		"CREATE TABLE t3 (id INT, a INT, b INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS objoin")
	})

	check := func(sql string, want bool) {
		t.Helper()
		engineTakes, err := engineAccepts(ctx, db, sql)
		if err != nil {
			t.Fatalf("the oracle gave no verdict on %q: %v", sql, err)
		}
		if engineTakes != want {
			t.Errorf("OceanBase accepts=%v, want %v for %q", engineTakes, want, sql)
		}
		_, perr := Parse(sql)
		if (perr == nil) != engineTakes {
			t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v", perr == nil, engineTakes, sql, perr)
		}
	}

	// An outer join with no condition is rejected.
	check("SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id LEFT JOIN t3", false)
	check("SELECT * FROM t1 RIGHT JOIN t2 USING (id) LEFT JOIN t3", false)
	// A conditionless plain/inner/outer join absorbs its right operand, so the
	// stacked trailing ON clauses bind to the nested joins.
	check("SELECT * FROM t3 JOIN t1 JOIN t2 ON t1.a=t2.a ON t3.b=t2.b", true)
	check("SELECT * FROM t1 LEFT JOIN t2 JOIN t3 ON t2.a=t3.a ON t1.b=t2.b", true)
	// A NATURAL join takes no condition and does not absorb its right operand.
	check("SELECT * FROM t1 NATURAL JOIN t2 ON t1.a=t2.a", false)
	check("SELECT * FROM t1 NATURAL JOIN t2 JOIN t3 ON t2.a=t3.a ON t1.b=t2.b", false)
	// OceanBase, unlike MySQL/TiDB, lets STRAIGHT_JOIN absorb its right operand.
	check("SELECT * FROM t1 STRAIGHT_JOIN t2 JOIN t3 ON t2.a=t3.a ON t1.b=t2.b", true)
	check("SELECT * FROM t1 STRAIGHT_JOIN t2 STRAIGHT_JOIN t3 ON t2.a=t3.a ON t1.b=t2.b", true)
}

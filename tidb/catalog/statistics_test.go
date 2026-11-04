package catalog

import "testing"

// TiDB's standalone statistics statements reach Exec because a change script
// may hold them next to its DDL: DROP STATS and LOAD STATS beside a
// migration that rebuilds a table, CREATE and DROP STATISTICS beside the
// index they are collected for.
//
// None changes a schema the catalog holds, so what they must do is run
// without error and leave everything as it was, and the statement after them
// must still see the table it had. Each is dispatched explicitly rather than
// left to processUtility's default arm, where an unmodelled statement and a
// deliberate no-op read the same.
//
// ANALYZE TABLE is here too, which is older than the dispatch: its wide TiDB
// forms -- a partition list, an index list, a column choice, a WITH list --
// are new, and a catalog that read any of them as a schema change would show
// up here.
//
// These tests run unconditionally: they need no container and so are outside
// the TIDB_CATALOG_PARITY quarantine in container_test.go.

func TestExecStatisticsStatementsChangeNoSchema(t *testing.T) {
	for _, sql := range []string{
		`ANALYZE TABLE t`,
		`ANALYZE TABLE t INDEX idx1`,
		`ANALYZE TABLE t INDEX`,
		`ANALYZE TABLE t ALL COLUMNS WITH 4 BUCKETS, 0.5 SAMPLERATE`,
		`ANALYZE TABLE t PREDICATE COLUMNS`,
		`ANALYZE TABLE t COLUMNS a, b`,
		`ANALYZE TABLE t PARTITION p0 INDEX idx1 WITH 4 CMSKETCH WIDTH`,
		`ANALYZE INCREMENTAL TABLE t INDEX idx1`,
		`ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS`,
		`ANALYZE TABLE t DROP HISTOGRAM ON a`,
		`DROP STATS t`,
		`DROP STATS t GLOBAL`,
		`DROP STATS t PARTITION p0`,
		`LOAD STATS '/tmp/stats.json'`,
		`CREATE STATISTICS s1 (CARDINALITY) ON t(a, b)`,
		`CREATE STATISTICS IF NOT EXISTS s2 (CORRELATION) ON t(a, b)`,
		`DROP STATISTICS s1`,
	} {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)

			before := snapshotTable(t, c)
			mustExec(t, c, sql)
			if got := snapshotTable(t, c); got != before {
				t.Errorf("the catalog changed:\nbefore %s\nafter  %s", before, got)
			}

			// The statement is not skipped the way a DML statement is: it
			// is applied, as a no-op.
			results, err := c.Exec(sql, nil)
			if err != nil {
				t.Fatalf("Exec(%q): %v", sql, err)
			}
			if len(results) != 1 {
				t.Fatalf("results = %d, want 1", len(results))
			}
			if results[0].Skipped {
				t.Error("the statement was skipped, want applied as a no-op")
			}
			if results[0].Error != nil {
				t.Errorf("Error = %v, want none", results[0].Error)
			}

			// The DDL that follows still sees the table it had, which is
			// what a catalog that lost the schema here would break.
			mustExec(t, c, `ALTER TABLE t ADD COLUMN c int`)
			if col := c.GetDatabase("testdb").GetTable("t").GetColumn("c"); col == nil {
				t.Error("the column the next statement adds is missing")
			}
		})
	}
}

// TestExecStatisticsStatementsKeepTheirNeighbours pins that a statistics
// statement in the middle of a script changes nothing about the statements
// around it: the table one above it made is there below it, and the drop
// below it still lands.
func TestExecStatisticsStatementsKeepTheirNeighbours(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)
	results, err := c.Exec(`CREATE TABLE u (id int);
ANALYZE TABLE u INDEX WITH 4 TOPN;
DROP STATS u;
CREATE STATISTICS s1 (CARDINALITY) ON u(id);
LOAD STATS '/tmp/s.json';
DROP STATISTICS s1;
ALTER TABLE u ADD COLUMN v int;`, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != 7 {
		t.Fatalf("results = %d, want 7", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("statement %d: %v", i+1, r.Error)
		}
	}
	tbl := c.GetDatabase("testdb").GetTable("u")
	if tbl == nil {
		t.Fatal("table u is missing")
	}
	if tbl.GetColumn("v") == nil {
		t.Error("column v is missing: the ALTER TABLE below the statistics statements did not land")
	}
}

package catalog

import "testing"

// LOAD DATA and IMPORT INTO reach Exec because a change script may hold
// either next to its DDL: a bulk load beside the migration that made the
// table it fills.
//
// Neither changes a schema the catalog holds, so what they must do is run
// without error and leave everything as it was, and the statement after them
// must still see the table it had. Each is dispatched explicitly rather than
// left to processUtility's default arm, where an unmodelled statement and a
// deliberate no-op read the same -- which is the whole difference a corpus
// recording only whether a statement parsed cannot see.
//
// These tests run unconditionally: they need no container and so are outside
// the TIDB_CATALOG_PARITY quarantine in container_test.go.

func TestExecBulkLoadStatementsChangeNoSchema(t *testing.T) {
	for _, sql := range []string{
		`LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t`,
		`LOAD DATA LOCAL INFILE '/tmp/t.csv' REPLACE INTO TABLE t`,
		`LOAD DATA INFILE '/tmp/t.csv' FORMAT 'delimited data' INTO TABLE t`,
		`LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t FIELDS TERMINATED BY X'6B6B' DEFINED NULL BY 'nil'`,
		`LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t LINES STARTING BY X'3E' TERMINATED BY X'0A' IGNORE 1 LINES`,
		`LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t (a, @v) SET b = @v`,
		`LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t WITH threads=10, detached`,
		`IMPORT INTO t FROM '/tmp/t.csv'`,
		`IMPORT INTO t (a, @1) SET b = @1 FROM '/tmp/t.csv' FORMAT 'csv' WITH thread=2`,
		`IMPORT INTO t FROM SELECT a, b FROM t`,
		`IMPORT INTO t FROM TABLE t`,
	} {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)

			before := snapshotTable(t, c)
			mustExec(t, c, sql)
			if got := snapshotTable(t, c); got != before {
				t.Errorf("the catalog changed:\nbefore %s\nafter  %s", before, got)
			}

			// The statement is applied, as a no-op, rather than skipped.
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

// TestExecBulkLoadStatementsKeepTheirNeighbours pins that a bulk load in the
// middle of a script changes nothing about the statements around it: the
// table one above it made is there below it, and the ALTER below it lands.
func TestExecBulkLoadStatementsKeepTheirNeighbours(t *testing.T) {
	c := setupWithDB(t)
	results, err := c.Exec(`CREATE TABLE u (id int);
LOAD DATA INFILE '/tmp/u.csv' FORMAT 'csv' INTO TABLE u WITH detached;
IMPORT INTO u (id) FROM '/tmp/u.csv' WITH thread=4;
ALTER TABLE u ADD COLUMN v int;`, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("results = %d, want 4", len(results))
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
		t.Error("column v is missing: the ALTER TABLE below the bulk loads did not land")
	}
}

package catalog

import "testing"

// SPLIT ... TABLE and DISTRIBUTE TABLE reach Exec because a change script may
// hold them next to its DDL. Neither changes a schema the catalog holds: the
// first pre-splits regions, the second schedules replicas. The catalog is
// their consumer either way, so what they must do is run without error and
// leave everything as it was, and the statement after them must still see the
// table it had.
//
// Each is dispatched explicitly rather than left to processUtility's default
// arm, where an unmodelled statement and a deliberate no-op read the same.
// These tests need no container.

func TestExecSplitAndDistributeChangeNoSchema(t *testing.T) {
	for _, sql := range []string{
		`SPLIT TABLE t BY ('a'),('b')`,
		`SPLIT TABLE t INDEX idx1 BY ('a')`,
		`SPLIT REGION FOR PARTITION TABLE t PARTITION (p0) INDEX idx1 BETWEEN ('a') AND ('z') REGIONS 10`,
		`SPLIT TABLE t BETWEEN () AND () REGIONS 10`,
		`DISTRIBUTE TABLE t RULE = 'leader-scatter' ENGINE = 'tikv'`,
		`DISTRIBUTE TABLE t PARTITION (p0) RULE = 'peer-scatter' ENGINE = 'tiflash' TIMEOUT = '30m'`,
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

// snapshotTable renders t's columns and indexes, so a change anywhere in them
// shows up as a different string.
func snapshotTable(t *testing.T, c *Catalog) string {
	t.Helper()
	db := c.GetDatabase("testdb")
	if db == nil {
		t.Fatal("database testdb is missing")
	}
	tbl := db.GetTable("t")
	if tbl == nil {
		t.Fatal("table t is missing")
	}
	s := "t("
	for i, col := range tbl.Columns {
		if i > 0 {
			s += ", "
		}
		s += col.Name + " " + col.DataType
	}
	s += ")"
	for _, idx := range tbl.Indexes {
		s += " index " + idx.Name
	}
	return s
}

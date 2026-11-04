package catalog

import "testing"

// TiDB's backup, restore, recover and flashback statements reach Exec because
// a change script may hold them next to its DDL. None of them tells the
// catalog what a schema becomes: BACKUP reads and changes nothing, while
// RESTORE, RECOVER TABLE and the bring-back FLASHBACK forms put back an
// object whose definition lives in a backup or in a dropped table the catalog
// keeps no record of, and FLASHBACK ... TO TIMESTAMP rolls data back to a
// point it never saw. So each must run without error and leave everything as
// it was, and the statement after it must still see the table it had.
//
// FLASHBACK TABLE t TO t1 is the one that would be easy to get wrong, because
// it looks like a rename. It is not: t has to be gone already for the
// statement to succeed and t1 is created, which its sibling makes plain --
// v8.5.5 answers 1007, "Can't create database 'd'; database exists", for
// FLASHBACK DATABASE d when d is still there. Reading it as a rename would
// drop t from the catalog and add a t1 whose columns the statement never gave.
//
// Each is dispatched explicitly rather than left to processUtility's default
// arm, where an unmodelled statement and a deliberate no-op read the same.
// These tests run unconditionally: they need no container and so are outside
// the TIDB_CATALOG_PARITY quarantine in container_test.go.

func TestExecBRIEAndFlashbackChangeNoSchema(t *testing.T) {
	for _, sql := range []string{
		`BACKUP DATABASE * TO 'noop://'`,
		`BACKUP DATABASE testdb TO 'noop://' snapshot 5 minute ago`,
		`BACKUP TABLE t TO 'noop://' concurrency 4 checksum off`,
		`BACKUP LOGS TO 'noop://' start_ts = '20220304'`,
		`RESTORE DATABASE * FROM 'noop://'`,
		`RESTORE TABLE t FROM 'noop://' online 1`,
		`RESTORE POINT FROM 'noop://log_backup' full_backup_storage = 'noop://full'`,
		`RECOVER TABLE BY JOB 11`,
		`RECOVER TABLE t`,
		`RECOVER TABLE t 100`,
		`FLASHBACK TABLE t`,
		`FLASHBACK TABLE t TO t1`,
		`FLASHBACK DATABASE testdb`,
		`FLASHBACK DATABASE testdb TO other`,
		`FLASHBACK SCHEMA testdb TO other`,
		`FLASHBACK CLUSTER TO TIMESTAMP '2021-05-26 16:45:26'`,
		`FLASHBACK CLUSTER TO TSO 445494955052105721`,
		`FLASHBACK TABLE t TO TIMESTAMP '2021-05-26 16:45:26'`,
		`FLASHBACK DATABASE testdb TO TSO 445494955052105724`,
	} {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)

			before := snapshotTable(t, c)
			mustExec(t, c, sql)
			if got := snapshotTable(t, c); got != before {
				t.Errorf("the catalog changed:\nbefore %s\nafter  %s", before, got)
			}

			// The statement is not skipped the way a DML statement is:
			// it is applied, as a no-op.
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

// TestExecFlashbackTableIsNotARename pins the one failure the test above
// would miss if FLASHBACK TABLE ... TO ... were read as a rename: the table
// it names stays where it was, under its own name, and the new name is not
// added to the catalog from a definition the statement never gave.
func TestExecFlashbackTableIsNotARename(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)
	mustExec(t, c, `FLASHBACK TABLE t TO t1`)

	db := c.GetDatabase("testdb")
	if db.GetTable("t") == nil {
		t.Error("table t is gone; FLASHBACK TABLE t TO t1 renamed it")
	}
	if db.GetTable("t1") != nil {
		t.Error("table t1 is in the catalog, which holds no definition for it")
	}

	// A database is the same: the name stays, and the new one is not made.
	mustExec(t, c, `FLASHBACK DATABASE testdb TO other`)
	if c.GetDatabase("testdb") == nil {
		t.Error("database testdb is gone; FLASHBACK DATABASE testdb TO other renamed it")
	}
	if c.GetDatabase("other") != nil {
		t.Error("database other is in the catalog, which holds no definition for it")
	}
}

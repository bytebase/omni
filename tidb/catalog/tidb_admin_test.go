package catalog

import (
	"errors"
	"testing"
)

// TiDB's RECOMMEND INDEX, SET CONFIG, CANCEL DISTRIBUTION JOB, KILL and FLUSH, and
// the log task and BR job statements, reach Exec because a change script may hold them
// next to its DDL. None changes a schema the catalog holds, so what they must do is run
// without error and leave everything as it was, and the statement after them must still
// see the table it had. RECOMMEND INDEX APPLY is the one that sounds as if it adds an
// index: v8.5.5 answers 1105 "unsupported action apply" for it when it runs.
//
// ALTER RANGE is the exception, and has a test of its own below.
//
// These tests need no container.

func TestExecTiDBAdminStatementsChangeNoSchema(t *testing.T) {
	for _, sql := range []string{
		`RECOMMEND INDEX RUN`,
		`RECOMMEND INDEX RUN FOR 'select * from t where a = 1'`,
		`RECOMMEND INDEX RUN FOR 'select * from nosuch' WITH max_num_index = 1`,
		`RECOMMEND INDEX SHOW`,
		`RECOMMEND INDEX APPLY 1`,
		`RECOMMEND INDEX IGNORE 1`,
		`RECOMMEND INDEX SET max_num_index = 1`,
		`SET CONFIG tikv log.level = 'info'`,
		`SET CONFIG '127.0.0.1:3306' auto-compaction-mode = TRUE`,
		`CANCEL DISTRIBUTION JOB 1`,
		`KILL 23123`,
		`KILL QUERY 23123`,
		`KILL TIDB 23123`,
		`KILL TIDB CONNECTION 23123`,
		`KILL TIDB QUERY 23123`,
		`KILL TIDB CONNECTION_ID()`,
		`FLUSH TIDB PLUGINS a, b`,
		`FLUSH PRIVILEGES`,
		`FLUSH TABLES t WITH READ LOCK`,
		`STOP BACKUP LOGS`,
		`PAUSE BACKUP LOGS gc_ttl = '20220304'`,
		`RESUME BACKUP LOGS`,
		`PURGE BACKUP LOGS FROM 'noop://' until_ts = '2012122304'`,
		`SHOW BACKUP LOGS STATUS`,
		`SHOW BACKUP LOGS METADATA FROM 'noop://'`,
		`SHOW BACKUP METADATA FROM 'noop://'`,
		`SHOW BR JOB 1234`,
		`SHOW BR JOB QUERY 1234`,
		`CANCEL BR JOB 1234`,
		// Clearing a range refers to nothing, so the schema is untouched too.
		`ALTER RANGE global PLACEMENT POLICY DEFAULT`,
	} {
		t.Run(sql, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE TABLE t (a varchar(64), b int, INDEX idx1 (a))`)

			before := snapshotTable(t, c)
			mustExec(t, c, sql)
			if got := snapshotTable(t, c); got != before {
				t.Errorf("the catalog changed:\nbefore %s\nafter  %s", before, got)
			}

			// The statement is not skipped the way a DML statement is: it is
			// applied, as a no-op, and reported on its own line.
			results, err := c.Exec("\n\n"+sql, nil)
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
			if results[0].Line != 3 {
				t.Errorf("Line = %d, want 3", results[0].Line)
			}

			// The DDL that follows still sees the table it had.
			mustExec(t, c, `ALTER TABLE t ADD COLUMN c int`)
		})
	}
}

// execOne runs one statement and returns the error Exec reported for it.
func execOne(t *testing.T, c *Catalog, sql string) error {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	if len(results) != 1 {
		t.Fatalf("Exec(%q): %d results, want 1", sql, len(results))
	}
	return results[0].Error
}

// TestExecAlterRange pins ALTER RANGE to what pingcap/tidb:v8.5.5 does with a range
// and a policy, all of it measured on 2026-11-05:
//
//   - a range other than global and meta answers 1105 "range name is not supported",
//     and answers it before the policy is looked at
//   - a policy that does not exist answers what ALTER DATABASE answers for one
//   - every spelling of default clears the range: DEFAULT, SET DEFAULT, `default`,
//     'default', 'DEFAULT'
//   - DROP PLACEMENT POLICY answers "still in use" while a range refers to the policy,
//     and succeeds once the range is cleared
//
// The catalog validates a policy reference on a database and a table already, so a
// range that it did not record would let a script drop a policy the server would not.
func TestExecAlterRange(t *testing.T) {
	var catErr *Error
	assertCode := func(t *testing.T, err error, code int) {
		t.Helper()
		if !errors.As(err, &catErr) || catErr.Code != code {
			t.Errorf("error = %v, want code %d", err, code)
		}
	}

	t.Run("range names", func(t *testing.T) {
		c := setupWithDB(t)
		for _, sql := range []string{
			"ALTER RANGE global PLACEMENT POLICY DEFAULT",
			"ALTER RANGE meta PLACEMENT POLICY DEFAULT",
			"ALTER RANGE GLOBAL PLACEMENT POLICY DEFAULT",
			"ALTER RANGE META PLACEMENT POLICY DEFAULT",
			"ALTER RANGE `global` PLACEMENT POLICY DEFAULT",
		} {
			if err := execOne(t, c, sql); err != nil {
				t.Errorf("%s: %v", sql, err)
			}
		}
		assertCode(t, execOne(t, c, "ALTER RANGE foo PLACEMENT POLICY DEFAULT"), ErrUnknownError)
		// The name is checked first: no policy has to exist to get this answer.
		assertCode(t, execOne(t, c, "ALTER RANGE foo PLACEMENT POLICY nosuch"), ErrUnknownError)
	})

	t.Run("a policy that does not exist", func(t *testing.T) {
		c := setupWithDB(t)
		assertCode(t, execOne(t, c, "ALTER RANGE global PLACEMENT POLICY nosuch"), ErrPlacementPolicyMissing)
		assertCode(t, execOne(t, c, "ALTER RANGE meta PLACEMENT POLICY = 'nosuch'"), ErrPlacementPolicyMissing)
	})

	for _, clear := range []string{"DEFAULT", "SET DEFAULT", "= DEFAULT", "`default`", "'default'", "= 'DEFAULT'"} {
		t.Run("a range holds its policy until "+clear, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, `CREATE PLACEMENT POLICY p PRIMARY_REGION="r" REGIONS="r"`)
			mustExec(t, c, "ALTER RANGE global PLACEMENT POLICY p")
			assertCode(t, execOne(t, c, "DROP PLACEMENT POLICY p"), ErrPlacementPolicyInUse)

			// The other range is its own reference, and case is not part of a name.
			mustExec(t, c, "ALTER RANGE META PLACEMENT POLICY = 'P'")
			mustExec(t, c, "ALTER RANGE global PLACEMENT POLICY "+clear)
			assertCode(t, execOne(t, c, "DROP PLACEMENT POLICY p"), ErrPlacementPolicyInUse)
			mustExec(t, c, "ALTER RANGE meta PLACEMENT POLICY "+clear)

			mustExec(t, c, "DROP PLACEMENT POLICY p")
			if c.GetPlacementPolicy("p") != nil {
				t.Error("the policy is still there")
			}
		})
	}

	t.Run("a second ALTER RANGE replaces the reference", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, `CREATE PLACEMENT POLICY p PRIMARY_REGION="r" REGIONS="r"`)
		mustExec(t, c, `CREATE PLACEMENT POLICY q PRIMARY_REGION="r" REGIONS="r"`)
		mustExec(t, c, "ALTER RANGE global PLACEMENT POLICY p")
		mustExec(t, c, "ALTER RANGE global PLACEMENT POLICY q")
		mustExec(t, c, "DROP PLACEMENT POLICY p")
		assertCode(t, execOne(t, c, "DROP PLACEMENT POLICY q"), ErrPlacementPolicyInUse)
	})

	t.Run("a statement that failed leaves the reference", func(t *testing.T) {
		c := setupWithDB(t)
		mustExec(t, c, `CREATE PLACEMENT POLICY p PRIMARY_REGION="r" REGIONS="r"`)
		mustExec(t, c, "ALTER RANGE global PLACEMENT POLICY p")
		assertCode(t, execOne(t, c, "ALTER RANGE global PLACEMENT POLICY nosuch"), ErrPlacementPolicyMissing)
		assertCode(t, execOne(t, c, "DROP PLACEMENT POLICY p"), ErrPlacementPolicyInUse)
	})
}

package catalog

import (
	"errors"
	"testing"
)

func execErrorCode(t *testing.T, c *Catalog, sql string) int {
	t.Helper()
	results, err := c.Exec(sql, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatalf("Exec(%q) parse error: %v", sql, err)
	}
	if len(results) == 0 || results[0].Error == nil {
		t.Fatalf("Exec(%q): expected error, got nil", sql)
	}
	var ce *Error
	if !errors.As(results[0].Error, &ce) {
		t.Fatalf("Exec(%q): expected *Error, got %T (%v)", sql, results[0].Error, results[0].Error)
	}
	return ce.Code
}

func TestViewDDLObjectKindAndRename(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE d; CREATE DATABASE d2; USE d;")
		mustExec(t, c, "CREATE TABLE t (b INT PRIMARY KEY, c INT); CREATE VIEW v AS SELECT b, c FROM t;")
		return c
	}

	t.Run("wrong object 1347 on base-table operations", func(t *testing.T) {
		for _, stmt := range []string{
			"ALTER TABLE v ADD COLUMN x INT",
			"ALTER TABLE v REBUILD PARTITION ALL",
			"CREATE INDEX idx ON v (c)",
			"DROP INDEX idx ON v",
			"CREATE TABLE t2 LIKE v",
		} {
			c := setup(t)
			if got := execErrorCode(t, c, stmt); got != ErrWrongObject {
				t.Errorf("%s: got error %d, want %d", stmt, got, ErrWrongObject)
			}
		}
	})

	t.Run("TRUNCATE PARTITION on view returns 1505", func(t *testing.T) {
		c := setup(t)
		if got := execErrorCode(t, c, "ALTER TABLE v TRUNCATE PARTITION ALL"); got != ErrPartitionMgmtOnNonPartitioned {
			t.Errorf("got error %d, want %d", got, ErrPartitionMgmtOnNonPartitioned)
		}
	})

	t.Run("RENAME TABLE supports views", func(t *testing.T) {
		c := setup(t)
		if got := execErrorCode(t, c, "RENAME TABLE v TO v"); got != ErrDupTable {
			t.Errorf("RENAME TABLE v TO v: got error %d, want %d", got, ErrDupTable)
		}
		if got := execErrorCode(t, c, "RENAME TABLE v TO d2.v"); got != ErrForbidSchemaChange {
			t.Errorf("RENAME TABLE v TO d2.v: got error %d, want %d", got, ErrForbidSchemaChange)
		}
		if got := execErrorCode(t, c, "ALTER TABLE t RENAME TO v"); got != ErrDupTable {
			t.Errorf("ALTER TABLE t RENAME TO v: got error %d, want %d", got, ErrDupTable)
		}
		mustExec(t, c, "CREATE TABLE IF NOT EXISTS v (a INT)")
		mustExec(t, c, "RENAME TABLE v TO v2")
		db := c.GetDatabase("d")
		if db.Views["v"] != nil || db.Views["v2"] == nil || db.Views["v2"].Name != "v2" {
			t.Fatalf("expected view v renamed to v2, got views=%v", db.Views)
		}
	})
}

func TestPartitionedTableValidations(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d;")
		mustExec(t, c, `CREATE TABLE p (
			b INT,
			k INT,
			s VARCHAR(20),
			PRIMARY KEY (b, k)
		) PARTITION BY RANGE (k) (
			PARTITION p0 VALUES LESS THAN (100),
			PARTITION p1 VALUES LESS THAN (200)
		)`)
		return c
	}

	cases := []struct {
		stmt string
		want int
	}{
		{"CREATE FULLTEXT INDEX ft ON p (k)", ErrTableCantHandleFt},
		{"CREATE FULLTEXT INDEX ft ON p (s)", ErrTableCantHandleFt},
		{"ALTER TABLE p ADD FULLTEXT (k)", ErrTableCantHandleFt},
		{"ALTER TABLE p ADD FULLTEXT (s)", ErrTableCantHandleFt},
		{"ALTER TABLE p ADD PARTITION", 1492},
		{"ALTER TABLE p ADD PARTITION PARTITIONS 1", 1492},
		{"ALTER TABLE p ADD PARTITION (PARTITION p2)", ErrPartitionRequiresValues},
		{"ALTER TABLE p ADD PARTITION (PARTITION p2 TABLESPACE missing_ts)", ErrPartitionRequiresValues},
		{"ALTER TABLE p ADD PARTITION (PARTITION p1 VALUES LESS THAN MAXVALUE)", ErrSameNamePartition},
		{"ALTER TABLE p REORGANIZE PARTITION", ErrReorgNoParam},
		{"ALTER TABLE p REORGANIZE PARTITION p1 INTO (PARTITION p1)", ErrPartitionRequiresValues},
		{"ALTER TABLE p COALESCE PARTITION 1", ErrCoalesceOnlyOnHashPartition},
		{"ALTER TABLE p DROP PARTITION p1, p1", ErrDropLastPartition},
		{"ALTER TABLE p DROP PARTITION p0, p1", ErrDropLastPartition},
		{"ALTER TABLE p DROP COLUMN k", ErrDependentByPartitionFunc},
		{"ALTER TABLE p RENAME COLUMN k TO k", ErrDependentByPartitionFunc},
		{"ALTER TABLE p RENAME COLUMN k TO k2", ErrDependentByPartitionFunc},
		{"ALTER TABLE p CHANGE COLUMN k k2 INT", ErrDependentByPartitionFunc},
	}
	for _, tc := range cases {
		c := setup(t)
		if got := execErrorCode(t, c, tc.stmt); got != tc.want {
			t.Errorf("%s: got error %d, want %d", tc.stmt, got, tc.want)
		}
	}

	// Same-name CHANGE / MODIFY on partitioning column succeed, and DROP COLUMN
	// with REMOVE PARTITIONING succeeds.
	c := setup(t)
	mustExec(t, c, "ALTER TABLE p CHANGE COLUMN k k INT")
	mustExec(t, c, "ALTER TABLE p MODIFY COLUMN k INT")
	mustExec(t, c, "ALTER TABLE p DROP COLUMN k REMOVE PARTITIONING")
	if c.GetDatabase("d").GetTable("p").GetColumn("k") != nil {
		t.Fatalf("expected column k dropped after REMOVE PARTITIONING")
	}

	// CREATE TABLE with FULLTEXT + PARTITION BY is 1214.
	c2 := New()
	mustExec(t, c2, "CREATE DATABASE d; USE d;")
	if got := execErrorCode(t, c2, "CREATE TABLE t (a INT, s VARCHAR(20), FULLTEXT (s)) PARTITION BY HASH (a)"); got != ErrTableCantHandleFt {
		t.Errorf("CREATE TABLE FULLTEXT+PARTITION: got error %d, want %d", got, ErrTableCantHandleFt)
	}
}

func TestForeignKeyPartitionAndGeneratedBackingIndex(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d;")
		mustExec(t, c, "CREATE TABLE parent (b INT PRIMARY KEY);")
		mustExec(t, c, `CREATE TABLE child (
			b INT PRIMARY KEY,
			pid INT,
			d INT,
			KEY idx_d (d),
			CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent (b),
			CONSTRAINT chk_d CHECK (d > 0)
		)`)
		return c
	}

	t.Run("partitioning on FK table returns 1506 with MySQL precedence", func(t *testing.T) {
		for _, tc := range []struct {
			stmt string
			want int
		}{
			{"ALTER TABLE child PARTITION BY KEY()", ErrForeignKeyOnPartitioned},
			{"ALTER TABLE child PARTITION BY KEY(zzz)", ErrForeignKeyOnPartitioned},
			{"ALTER TABLE child PARTITION BY HASH(pid)", ErrForeignKeyOnPartitioned},
			{"ALTER TABLE child PARTITION BY KEY(pid, pid)", ErrDupPartitionField},
			{"ALTER TABLE child PARTITION BY KEY() SUBPARTITION BY HASH(pid)", ErrPartitionSubpartMix},
			{"ALTER TABLE child PARTITION BY RANGE(pid)", 1492},
			{"ALTER TABLE parent PARTITION BY KEY()", ErrForeignKeyOnPartitioned},
		} {
			c := setup(t)
			if got := execErrorCode(t, c, tc.stmt); got != tc.want {
				t.Errorf("%s: got error %d, want %d", tc.stmt, got, tc.want)
			}
		}
		// Dropping the FK in the same ALTER allows partitioning.
		c := setup(t)
		mustExec(t, c, "ALTER TABLE child DROP FOREIGN KEY fk_pid PARTITION BY KEY()")
	})

	t.Run("ALTER CONSTRAINT on non-CHECK returns 3941", func(t *testing.T) {
		c := setup(t)
		if got := execErrorCode(t, c, "ALTER TABLE child ALTER CONSTRAINT fk_pid ENFORCED"); got != ErrAlterConstraintNotSupported {
			t.Errorf("got error %d, want %d", got, ErrAlterConstraintNotSupported)
		}
		if got := execErrorCode(t, c, "ALTER TABLE child ALTER CHECK fk_pid ENFORCED"); got != ErrCheckConstraintNotFound {
			t.Errorf("got error %d, want %d", got, ErrCheckConstraintNotFound)
		}
		if got := execErrorCode(t, c, "ALTER TABLE child ALTER CONSTRAINT missing ENFORCED"); got != ErrConstraintDoesNotExist {
			t.Errorf("got error %d, want %d", got, ErrConstraintDoesNotExist)
		}
	})

	t.Run("dropping index needed by FK returns 1553", func(t *testing.T) {
		for _, stmt := range []string{
			"DROP INDEX fk_pid ON child",
			"ALTER TABLE child DROP INDEX fk_pid",
			"ALTER TABLE parent DROP PRIMARY KEY",
		} {
			c := setup(t)
			if got := execErrorCode(t, c, stmt); got != ErrDropIndexNeededInFK {
				t.Errorf("%s: got error %d, want %d", stmt, got, ErrDropIndexNeededInFK)
			}
		}
		// Dropping FK and index in the same ALTER succeeds.
		c := setup(t)
		mustExec(t, c, "ALTER TABLE child DROP INDEX fk_pid, DROP FOREIGN KEY fk_pid")
	})

	t.Run("generated FK backing index is replaced by explicit index", func(t *testing.T) {
		c := setup(t)
		mustExec(t, c, "CREATE INDEX fk_pid ON child (pid)")
		tbl := c.GetDatabase("d").GetTable("child")
		count := 0
		for _, idx := range tbl.Indexes {
			if idx.Name == "fk_pid" {
				count++
				if idx.Generated {
					t.Errorf("explicitly created index fk_pid should not stay Generated")
				}
			}
		}
		if count != 1 {
			t.Fatalf("expected 1 index named fk_pid, got %d", count)
		}
		// A second CREATE INDEX of the same name now fails with 1061.
		if got := execErrorCode(t, c, "CREATE INDEX fk_pid ON child (pid)"); got != ErrDupKeyName {
			t.Errorf("second CREATE INDEX fk_pid: got %d, want %d", got, ErrDupKeyName)
		}

		// Unnamed ADD KEY replaces generated backing index without allocating _2.
		c2 := setup(t)
		mustExec(t, c2, "ALTER TABLE child ADD KEY (pid) COMMENT 'explicit'")
		tbl2 := c2.GetDatabase("d").GetTable("child")
		for _, idx := range tbl2.Indexes {
			if idx.Name == "fk_pid" {
				t.Errorf("generated backing index fk_pid should have been dropped when explicit KEY (pid) was added")
			}
		}

		// FULLTEXT / SPATIAL over generated backing index name report 1283 / 1687, not 1061.
		c3 := setup(t)
		if got := execErrorCode(t, c3, "CREATE FULLTEXT INDEX fk_pid ON child (pid)"); got != ErrCannotPartOfFulltext {
			t.Errorf("CREATE FULLTEXT INDEX fk_pid: got %d, want %d", got, ErrCannotPartOfFulltext)
		}
		if got := execErrorCode(t, c3, "CREATE SPATIAL INDEX fk_pid ON child (pid)"); got != ErrSpatialColumnNotGeometrical {
			t.Errorf("CREATE SPATIAL INDEX fk_pid: got %d, want %d", got, ErrSpatialColumnNotGeometrical)
		}

		// Adding a non-covering index preserves the generated FK backing index.
		c4 := setup(t)
		mustExec(t, c4, "ALTER TABLE child ADD KEY ((1))")
		tbl4 := c4.GetDatabase("d").GetTable("child")
		if !indexNameExists(tbl4, "fk_pid") {
			t.Errorf("expected generated backing index fk_pid to be preserved after non-covering ADD KEY")
		}
	})
}

package catalog

import (
	"errors"
	"strings"
	"testing"
)

func execExpectCode(t *testing.T, setup, stmt string, wantCode int) {
	t.Helper()
	c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
	if err != nil {
		t.Fatalf("LoadSQL(%q): %v", setup, err)
	}
	results, err := c.Exec(stmt, nil)
	if err != nil {
		t.Fatalf("parse(%q): %v", stmt, err)
	}
	var gotErr error
	for _, r := range results {
		if r.Error != nil {
			gotErr = r.Error
			break
		}
	}
	if wantCode == 0 {
		if gotErr != nil {
			t.Fatalf("%s: unexpected error %v", stmt, gotErr)
		}
		return
	}
	if gotErr == nil {
		t.Fatalf("%s: expected error %d, got nil", stmt, wantCode)
	}
	var catErr *Error
	if !errors.As(gotErr, &catErr) || catErr.Code != wantCode {
		t.Fatalf("%s: got error %v, want code %d", stmt, gotErr, wantCode)
	}
}

func TestViewTableDDLRefusals(t *testing.T) {
	const setup = `
		CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
		CREATE VIEW a AS SELECT b, a, d FROM t;
		CREATE DATABASE d2;
	`
	cases := []struct {
		stmt string
		code int
	}{
		// ALTER TABLE on a view answers 1347, even with IF EXISTS.
		{"ALTER TABLE a;", 1347},
		{"ALTER TABLE IF EXISTS a;", 1347},
		{"ALTER TABLE a ADD COLUMN x INT;", 1347},
		{"ALTER TABLE a REMOVE PARTITIONING;", 1347},
		// TRUNCATE PARTITION checks partitioning before table-vs-view (1505).
		{"ALTER TABLE a TRUNCATE PARTITION ALL;", 1505},
		{"ALTER TABLE IF EXISTS a TRUNCATE PARTITION ALL;", 1505},
		// Partition maintenance commands succeed silently on views too.
		{"ALTER TABLE a OPTIMIZE PARTITION ALL;", 0},
		{"ALTER TABLE a ANALYZE PARTITION ALL;", 0},
		{"ALTER TABLE a CHECK PARTITION ALL;", 0},
		{"ALTER TABLE a REPAIR PARTITION ALL;", 0},
		// CREATE/DROP INDEX on a view answer 1347.
		{"CREATE INDEX idx ON a (a);", 1347},
		{"DROP INDEX idx ON a;", 1347},
		{"DROP INDEX IF EXISTS idx ON a;", 1347},
		// DROP TABLE on a view answers 1965 (unless IF EXISTS or TEMPORARY).
		{"DROP TABLE a;", 1965},
		{"DROP TABLE a CASCADE;", 1965},
		{"DROP TABLE a RESTRICT;", 1965},
		{"DROP TABLE IF EXISTS a;", 0},
		{"DROP TEMPORARY TABLE a;", 1051},
		// RENAME TABLE works on views within the same DB, refuses existing target (1050)
		// and cross-DB view rename (1450).
		{"RENAME TABLE a TO a;", 1050},
		{"RENAME TABLE a TO d2.v2;", 1450},
		{"RENAME TABLE a TO v2;", 0},
	}
	for _, tc := range cases {
		t.Run(tc.stmt, func(t *testing.T) {
			execExpectCode(t, setup, tc.stmt, tc.code)
		})
	}

	// Verify DROP TABLE IF EXISTS preserves the view and RENAME TABLE moves it.
	c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Exec("DROP TABLE IF EXISTS a;", nil); err != nil {
		t.Fatal(err)
	}
	if c.GetDatabase("d").Views["a"] == nil {
		t.Fatal("DROP TABLE IF EXISTS should preserve view `a`")
	}
	if _, err := c.Exec("RENAME TABLE a TO v2;", nil); err != nil {
		t.Fatal(err)
	}
	db := c.GetDatabase("d")
	if db.Views["a"] != nil || db.Views["v2"] == nil || db.Views["v2"].Name != "v2" {
		t.Fatalf("expected view renamed from `a` to `v2`, got views=%v", db.Views)
	}
}

func TestPartitionedTableDDL(t *testing.T) {
	const setup = `
		CREATE TABLE a (b INT, a INT, d INT, PRIMARY KEY (b, a), KEY idx_d (d))
			PARTITION BY RANGE (a) (
				PARTITION p0 VALUES LESS THAN (100),
				PARTITION a VALUES LESS THAN (200)
			);
	`
	const setup3 = `
		CREATE TABLE a (b INT, a INT, PRIMARY KEY (b, a))
			PARTITION BY RANGE (a) (
				PARTITION p0 VALUES LESS THAN (100),
				PARTITION p1 VALUES LESS THAN (150),
				PARTITION a VALUES LESS THAN (200)
			);
	`
	const setupHash = `CREATE TABLE a (b INT) PARTITION BY HASH (b) PARTITIONS 4;`

	cases := []struct {
		setup string
		stmt  string
		code  int
	}{
		// FK operations on a partitioned table answer 1506 before column/constraint lookups.
		{setup, "ALTER TABLE a ADD FOREIGN KEY (a) REFERENCES a;", 1506},
		{setup, "ALTER TABLE a ADD a INT REFERENCES a;", 1506},
		{setup, "ALTER TABLE a DROP FOREIGN KEY a;", 1506},
		// DROP CONSTRAINT (not DROP FOREIGN KEY) answers 1091.
		{setup, "ALTER TABLE a DROP CONSTRAINT a;", 1091},
		// SEQUENCE 1 on partitioned table answers 1478 (partition engine).
		{setup, "ALTER TABLE a SEQUENCE 1;", 1478},
		// Partition management validations on RANGE tables.
		{setup, "ALTER TABLE a ADD PARTITION;", 1492},
		{setup, "ALTER TABLE a ADD PARTITION PARTITIONS 1;", 1492},
		{setup, "ALTER TABLE a DROP PARTITION a, a;", 1508},
		{setup, "ALTER TABLE a DROP PARTITION p0, a;", 1508},
		{setup, "ALTER TABLE a DROP PARTITION nope;", 1507},
		{setup3, "ALTER TABLE a DROP PARTITION a, a;", 1507},
		{setup, "ALTER TABLE a COALESCE PARTITION 1;", 1509},
		{setup, "ALTER TABLE a REORGANIZE PARTITION;", 1511},
		{setup, "ALTER TABLE a REORGANIZE PARTITION a INTO (PARTITION a);", 1479},
		// REORGANIZE PARTITION without parameters succeeds on HASH-partitioned tables.
		{setupHash, "ALTER TABLE a REORGANIZE PARTITION;", 0},
	}
	for _, tc := range cases {
		t.Run(tc.stmt, func(t *testing.T) {
			execExpectCode(t, tc.setup, tc.stmt, tc.code)
		})
	}

	// MODIFY/CHANGE on a composite PRIMARY KEY column preserves NOT NULL.
	for _, stmt := range []string{"ALTER TABLE a MODIFY a INT;", "ALTER TABLE a CHANGE a a INT;", "ALTER TABLE a MODIFY a INT NULL;"} {
		c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
		if err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec(stmt, nil)
		if err != nil || results[0].Error != nil {
			t.Fatalf("%s failed: parse=%v exec=%v", stmt, err, results)
		}
		col := c.GetDatabase("d").GetTable("a").GetColumn("a")
		if col.Nullable {
			t.Errorf("%s: PK column `a` should remain NOT NULL", stmt)
		}
	}

	// Dropping the PK in the same ALTER allows the column to become nullable.
	c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
	if err != nil {
		t.Fatal(err)
	}
	results, err := c.Exec("ALTER TABLE a DROP PRIMARY KEY, MODIFY a INT;", nil)
	if err != nil || results[0].Error != nil {
		t.Fatalf("DROP PRIMARY KEY + MODIFY failed: parse=%v exec=%v", err, results)
	}
	if !c.GetDatabase("d").GetTable("a").GetColumn("a").Nullable {
		t.Errorf("column `a` should be nullable once PK is dropped")
	}
}

func TestForeignKeyImplicitIndexAndConstraints(t *testing.T) {
	const setup = `
		CREATE TABLE parent (b INT PRIMARY KEY);
		CREATE TABLE a (
			b INT PRIMARY KEY,
			a INT,
			d INT,
			KEY idx_d (d),
			CONSTRAINT a FOREIGN KEY (a) REFERENCES parent (b),
			CONSTRAINT chk_d CHECK (d > 0)
		);
	`
	const setupExplicit = `
		CREATE TABLE parent (b INT PRIMARY KEY);
		CREATE TABLE a (
			b INT PRIMARY KEY,
			a INT,
			KEY a (a),
			CONSTRAINT a FOREIGN KEY (a) REFERENCES parent (b)
		);
	`

	errorCases := []struct {
		setup string
		stmt  string
		code  int
	}{
		// Implicit FK backing index is replaced by covering BTREE/default index, even with the same name.
		{setup, "CREATE INDEX a ON a (a);", 0},
		{setup, "CREATE UNIQUE INDEX a ON a (a);", 0},
		{setup, "ALTER TABLE a ADD KEY a (a);", 0},
		{setup, "ALTER TABLE a ADD KEY a TYPE BTREE (a);", 0},
		{setup, "ALTER TABLE a ADD UNIQUE a TYPE BTREE (a);", 0},
		// Non-covering index with the same name still fails 1061.
		{setup, "ALTER TABLE a ADD KEY a (d);", 1061},
		// Explicit (non-implicit) index is NOT replaced -> 1061.
		{setupExplicit, "CREATE INDEX a ON a (a);", 1061},
		// Dropping the sole covering FK index fails with 1553 when foreign_key_checks=1.
		{setup, "ALTER TABLE a DROP INDEX a;", 1553},
		{setup, "DROP INDEX a ON a;", 1553},
		// Dropping the FK or adding a replacement index in the same ALTER succeeds.
		{setup, "ALTER TABLE a DROP FOREIGN KEY a, DROP INDEX a;", 0},
		{setup, "ALTER TABLE a DROP INDEX a, ADD INDEX k_a (a);", 0},
		// Dropping with foreign_key_checks=0 succeeds.
		{setup, "SET foreign_key_checks = 0; ALTER TABLE a DROP INDEX a;", 0},
		// Dropping a column needed by an FK fails with 1553 when foreign_key_checks=1
		// and 1828 when foreign_key_checks=0.
		{setup, "ALTER TABLE a DROP COLUMN a;", 1553},
		{setup, "SET foreign_key_checks = 0; ALTER TABLE a DROP COLUMN a;", 1828},
		{setup, "ALTER TABLE a DROP FOREIGN KEY a, DROP COLUMN a;", 0},
		// CHECK constraint name collides with FK constraint name on the same table (1826),
		// even with IF NOT EXISTS (which only matches existing CHECK constraints).
		{setup, "ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS a CHECK (a);", 1826},
		{setup, "ALTER TABLE a ADD CONSTRAINT a CHECK (a);", 1826},
		{setup, "ALTER TABLE a ADD CONSTRAINT IF NOT EXISTS chk_d CHECK (d > 0);", 0},
		// Partitioning a table that has or is referenced by an FK answers 1217.
		{setup, "ALTER TABLE a PARTITION BY SYSTEM_TIME;", 1217},
		{setup, "ALTER TABLE parent PARTITION BY SYSTEM_TIME;", 1217},
	}
	for _, tc := range errorCases {
		t.Run(tc.stmt, func(t *testing.T) {
			execExpectCode(t, tc.setup, tc.stmt, tc.code)
		})
	}

	// Verify schema effects of implicit FK index replacement vs preservation.
	t.Run("ADD KEY zzz(a) drops implicit FK index a", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.Exec("ALTER TABLE a ADD KEY zzz (a);", nil); err != nil {
			t.Fatal(err)
		}
		ddl := c.ShowCreateTable("d", "a")
		if strings.Contains(ddl, "KEY `a` (`a`)") {
			t.Errorf("implicit FK index `a` should be dropped, got:\n%s", ddl)
		}
		if !strings.Contains(ddl, "KEY `zzz` (`a`)") {
			t.Errorf("expected KEY `zzz` (`a`), got:\n%s", ddl)
		}
	})

	t.Run("ADD KEY USING HASH(a) preserves implicit FK index a", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.Exec("ALTER TABLE a ADD KEY USING HASH (a);", nil); err != nil {
			t.Fatal(err)
		}
		ddl := c.ShowCreateTable("d", "a")
		if !strings.Contains(ddl, "KEY `a` (`a`)") || !strings.Contains(ddl, "KEY `a_2` (`a`) USING HASH") {
			t.Errorf("expected both KEY `a` and KEY `a_2` USING HASH, got:\n%s", ddl)
		}
	})

	t.Run("ADD KEY(a) IGNORED drops implicit FK index a", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d;\n" + setup)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.Exec("ALTER TABLE a ADD KEY (a) IGNORED;", nil); err != nil {
			t.Fatal(err)
		}
		ddl := c.ShowCreateTable("d", "a")
		if strings.Contains(ddl, "KEY `a_2`") || !strings.Contains(ddl, "KEY `a` (`a`) IGNORED") {
			t.Errorf("expected KEY `a` (`a`) IGNORED replacing implicit key, got:\n%s", ddl)
		}
	})
}

package catalog

import (
	"strings"
	"testing"
)

// A hex or bit literal default on a column that is not BIT is the value it
// spells: text for a string or binary string column, an integer for a numeric
// one (getDefaultValue in TiDB's pkg/ddl/add_column.go; v8.5.5's SHOW CREATE
// TABLE for this statement, read on 2026-11-06).
func TestShowCreateTableBinaryLiteralDefaults(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE h (a VARCHAR(10) DEFAULT 0x41, b INT DEFAULT 0x41, c VARCHAR(10) DEFAULT b'01000010', d CHAR(3) DEFAULT X'4344', e INT DEFAULT b'101', f VARBINARY(10) DEFAULT 0x41, g BINARY(2) DEFAULT X'4344')")
	got := c.ShowCreateTable("testdb", "h")
	for _, want := range []string{
		"`a` varchar(10) DEFAULT 'A'",
		"`b` int DEFAULT '65'",
		"`c` varchar(10) DEFAULT 'B'",
		"`d` char(3) DEFAULT 'CD'",
		"`e` int DEFAULT '5'",
		"`f` varbinary(10) DEFAULT 'A'",
		"`g` binary(2) DEFAULT 'CD'",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %s in:\n%s", want, got)
		}
	}
}

// TiDB names a foreign key's backing index after the key and refuses the
// statement with 1061 when an index on other columns already has that name, in
// CREATE TABLE and in ALTER TABLE alike (v8.5.5, read on 2026-11-06).
func TestForeignKeyBackingIndexNameTaken(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE p (id INT PRIMARY KEY)")
	results, err := c.Exec("CREATE TABLE c (a INT, b INT, KEY fk_1(b), FOREIGN KEY(a) REFERENCES p(id))", nil)
	if err != nil {
		t.Fatal(err)
	}
	assertError(t, results[0].Error, ErrDupKeyName)
	if c.GetDatabase("testdb").GetTable("c") != nil {
		t.Errorf("the refused CREATE TABLE left a table")
	}

	wtExec(t, c, "CREATE TABLE c2 (a INT, b INT, KEY fk_1(b))")
	results, err = c.Exec("ALTER TABLE c2 ADD FOREIGN KEY (a) REFERENCES p(id)", nil)
	if err != nil {
		t.Fatal(err)
	}
	assertError(t, results[0].Error, ErrDupKeyName)
	for _, con := range c.GetDatabase("testdb").GetTable("c2").Constraints {
		if con.Type == ConForeignKey {
			t.Errorf("the refused ALTER TABLE left foreign key %s", con.Name)
		}
	}
}

// A refused ALTER DATABASE leaves the database as it was.
func TestAlterDatabaseRefusedLeavesDatabase(t *testing.T) {
	c := wtSetup(t)
	for _, sql := range []string{
		"ALTER DATABASE testdb CHARACTER SET nope",
		"ALTER DATABASE testdb COLLATE latin1_swedish_ci",
		"ALTER DATABASE testdb CHARACTER SET latin1 COLLATE utf8mb4_bin",
	} {
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatal(err)
		}
		if results[0].Error == nil {
			t.Errorf("%s was accepted", sql)
		}
		if db := c.GetDatabase("testdb"); db.Charset != "utf8mb4" || db.Collation != "utf8mb4_bin" {
			t.Errorf("after %s: charset %q collation %q", sql, db.Charset, db.Collation)
		}
	}
}

// A partition's values come out in one pair of parentheses, so a LIST COLUMNS
// partition with a single tuple keeps the tuple's own, and a string value is
// written as it was, escaped quote and all; what is written must read back.
func TestShowCreateTablePartitionValues(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE l (a INT, b INT) PARTITION BY LIST COLUMNS (a, b) (PARTITION p0 VALUES IN ((1,2)), PARTITION p1 VALUES IN ((3,4),(5,6)))`)
	wtExec(t, c, `CREATE TABLE s (a VARCHAR(10)) PARTITION BY LIST COLUMNS (a) (PARTITION p0 VALUES IN ('a\'x, y', 'z'))`)
	got := c.ShowCreateTable("testdb", "l")
	for _, want := range []string{"PARTITION `p0` VALUES IN ((1,2))", "PARTITION `p1` VALUES IN ((3,4),(5,6))"} {
		if !strings.Contains(got, want) {
			t.Errorf("want %s in:\n%s", want, got)
		}
	}
	for _, table := range []string{"l", "s"} {
		ddl := c.ShowCreateTable("testdb", table)
		again := wtSetup(t)
		results, err := again.Exec(ddl, nil)
		if err != nil || results[0].Error != nil {
			t.Fatalf("%s does not read back: %v %v", ddl, err, results)
		}
		if back := again.ShowCreateTable("testdb", table); back != ddl {
			t.Errorf("round trip changed %s:\n%s\nto\n%s", table, ddl, back)
		}
	}
}

// A partitioned table with a TTL writes the TTL on the options line, before
// the partition clause, as v8.5.5 does; the server refuses table options
// after the partitions, so the other order would not read back.
func TestShowCreateTableTTLBeforePartitions(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR PARTITION BY HASH(id) PARTITIONS 2")
	want := "CREATE TABLE `t` (\n" +
		"  `id` int DEFAULT NULL,\n" +
		"  `created_at` datetime DEFAULT NULL\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![ttl] TTL=`created_at` + INTERVAL 1 YEAR */ /*T![ttl] TTL_ENABLE='ON' */ /*T![ttl] TTL_JOB_INTERVAL='1h' */\n" +
		"PARTITION BY HASH (`id`) PARTITIONS 2"
	assertShowCreate(t, c, "t", want)
	again := wtSetup(t)
	results, err := again.Exec(want, nil)
	if err != nil || results[0].Error != nil {
		t.Fatalf("does not read back: %v %v", err, results)
	}
}

// A binary default holding bytes that are not UTF-8 is written byte for byte.
func TestShowCreateTableBinaryDefaultBytes(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a VARBINARY(4) DEFAULT 0xFF27)")
	if got := c.ShowCreateTable("testdb", "t"); !strings.Contains(got, "`a` varbinary(4) DEFAULT '\xff'''") {
		t.Errorf("want the bytes ff 27 escaped in place, got:\n%q", got)
	}
}

// A refused ADD FOREIGN KEY uses no number: TiDB keeps MaxForeignKeyID in the
// table's meta, which only a finished statement changes.
func TestRefusedAddForeignKeyKeepsCounter(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE c (a INT)")
	results, err := c.Exec("ALTER TABLE c ADD FOREIGN KEY (a) REFERENCES nosuch(id)", nil)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Error == nil {
		t.Fatal("a foreign key to a missing table was accepted")
	}
	wtExec(t, c, "CREATE TABLE p (id INT PRIMARY KEY)")
	wtExec(t, c, "ALTER TABLE c ADD FOREIGN KEY (a) REFERENCES p(id)")
	for _, con := range c.GetDatabase("testdb").GetTable("c").Constraints {
		if con.Type == ConForeignKey && con.Name != "fk_1" {
			t.Errorf("foreign key named %s, want fk_1", con.Name)
		}
	}
}

// With tidb_enable_check_constraint off, DROP CHECK and DROP CONSTRAINT may
// name nothing, but DROP FOREIGN KEY of a missing key is still 1091 (v8.5.5,
// read on 2026-11-06).
func TestDropMissingForeignKeyWithChecksOff(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE f (a INT)")
	wtExec(t, c, "ALTER TABLE f DROP CHECK nosuch")
	wtExec(t, c, "ALTER TABLE f DROP CONSTRAINT nosuch")
	results, err := c.Exec("ALTER TABLE f DROP FOREIGN KEY nosuch", nil)
	if err != nil {
		t.Fatal(err)
	}
	assertError(t, results[0].Error, ErrCantDropKey)
}

// A TTL first added by ALTER TABLE is enabled, as one given at CREATE is.
func TestAlterTableAddsEnabledTTL(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE tt (created_at DATETIME)")
	wtExec(t, c, "ALTER TABLE tt TTL = created_at + INTERVAL 1 DAY")
	assertShowCreate(t, c, "tt", "CREATE TABLE `tt` (\n"+
		"  `created_at` datetime DEFAULT NULL\n"+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![ttl] TTL=`created_at` + INTERVAL 1 DAY */ /*T![ttl] TTL_ENABLE='ON' */ /*T![ttl] TTL_JOB_INTERVAL='1h' */")
}

// CREATE TABLE LIKE keeps the source primary key's clustering.
func TestCreateTableLikeKeepsClustering(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE src (id INT PRIMARY KEY NONCLUSTERED)")
	wtExec(t, c, "CREATE TABLE dst LIKE src")
	assertShowCreate(t, c, "dst", "CREATE TABLE `dst` (\n"+
		"  `id` int NOT NULL,\n"+
		"  PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */\n"+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")
}

// CONVERT TO CHARACTER SET DEFAULT parses, and v8.5.5 refuses it with 1105
// "Unknown charset " (read on 2026-11-06), leaving the table as it was.
func TestConvertToDefaultCharsetRefused(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE cv CONVERT TO CHARSET DEFAULT",
		"ALTER TABLE cv CONVERT TO CHARACTER SET DEFAULT COLLATE utf8mb4_general_ci",
	} {
		c := wtSetup(t)
		wtExec(t, c, "CREATE TABLE cv (a VARCHAR(5)) CHARSET latin1")
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatal(err)
		}
		assertError(t, results[0].Error, ErrUnknownError)
		if tbl := c.GetDatabase("testdb").GetTable("cv"); tbl.Charset != "latin1" {
			t.Errorf("%s changed the charset to %q", sql, tbl.Charset)
		}
	}
}

// TiDB refuses a literal default on a BLOB, TEXT or JSON column with 1101 in
// every statement that sets one, and takes DEFAULT NULL and an expression
// default (v8.5.5, read on 2026-11-06).
func TestBlobTextJSONLiteralDefaultRefused(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (b BLOB DEFAULT 0x4142)",
		"CREATE TABLE t (b BLOB DEFAULT '')",
		"CREATE TABLE t (t TINYTEXT NOT NULL DEFAULT 'a')",
		"CREATE TABLE t (j JSON DEFAULT '[]')",
		"CREATE TABLE t (a INT); ALTER TABLE t ADD COLUMN b TEXT DEFAULT 'x'",
		"CREATE TABLE t (b TEXT); ALTER TABLE t MODIFY b TEXT DEFAULT 'x'",
		"CREATE TABLE t (b TEXT); ALTER TABLE t ALTER COLUMN b SET DEFAULT 'x'",
	} {
		c := wtSetup(t)
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		assertError(t, results[len(results)-1].Error, ErrBlobCantHaveDefault)
	}
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE ok (j JSON DEFAULT (JSON_ARRAY()), b LONGBLOB DEFAULT NULL)")
	if got := c.ShowCreateTable("testdb", "ok"); !strings.Contains(got, "`b` longblob DEFAULT NULL") {
		t.Errorf("got:\n%s", got)
	}
}

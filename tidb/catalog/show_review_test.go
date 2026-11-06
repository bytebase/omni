package catalog

import (
	"strings"
	"testing"
)

// A hex or bit literal default on a column that is not BIT is the value it
// spells: text for a string column, an integer for a numeric one
// (getDefaultValue in TiDB's pkg/ddl/add_column.go).
func TestShowCreateTableBinaryLiteralDefaults(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE h (a VARCHAR(10) DEFAULT 0x41, b INT DEFAULT 0x41, c VARCHAR(10) DEFAULT b'01000010', d CHAR(3) DEFAULT X'4344', e INT DEFAULT b'101')")
	got := c.ShowCreateTable("testdb", "h")
	for _, want := range []string{
		"`a` varchar(10) DEFAULT 'A'",
		"`b` int DEFAULT '65'",
		"`c` varchar(10) DEFAULT 'B'",
		"`d` char(3) DEFAULT 'CD'",
		"`e` int DEFAULT '5'",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %s in:\n%s", want, got)
		}
	}
}

// TiDB names a foreign key's backing index after the key and refuses the
// statement with 1061 when an index on other columns already has that name, in
// CREATE TABLE and in ALTER TABLE alike.
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

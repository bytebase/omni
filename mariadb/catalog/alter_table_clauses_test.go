package catalog

import (
	"strings"
	"testing"
)

// These tests assert what the ALTER TABLE clauses the parser newly reads do to
// the catalog, not that they parse. A clause that parses and applies nothing
// looks like a closed gap to a parser test. Neither test boots a container.

// execAlterClauses runs a script and fails on the first statement that errors.
func execAlterClauses(t *testing.T, c *Catalog, sql string) {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("Exec(%q): line %d: %v", sql, r.Line, r.Error)
		}
	}
}

// TestAlterTableAddColumnTableElementList: the parenthesized ADD reads a table
// element list, so the constraints written beside the columns apply too. The
// reader took only the columns before, which parsed and added no index.
func TestAlterTableAddColumnTableElementList(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT);")
	execAlterClauses(t, c, "ALTER TABLE t ADD COLUMN (b INT, c INT, PRIMARY KEY (a), UNIQUE KEY u (b));")

	ddl := c.ShowCreateTable("d", "t")
	for _, want := range []string{"`b` int", "`c` int", "PRIMARY KEY (`a`)", "UNIQUE KEY `u` (`b`)"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}
}

// TestAlterTablePartDefStorageEngine: STORAGE ENGINE is the other spelling of
// a partition's ENGINE option and names the same engine.
func TestAlterTablePartDefStorageEngine(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10));")
	execAlterClauses(t, c, "ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) STORAGE ENGINE = InnoDB);")

	tbl := c.databases["d"].Tables["t"]
	last := tbl.Partitioning.Partitions[len(tbl.Partitioning.Partitions)-1]
	if last.Name != "p1" {
		t.Fatalf("the last partition is %q, want p1", last.Name)
	}
	if last.Engine != "InnoDB" {
		t.Errorf("p1 engine = %q, want \"InnoDB\"", last.Engine)
	}
}

// TestAlterTableRenameToQualifiedName: the destination of a RENAME may carry a
// database, and then the table moves there.
func TestAlterTableRenameToQualifiedName(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; CREATE DATABASE d2; USE d; CREATE TABLE t (a INT);")
	execAlterClauses(t, c, "ALTER TABLE t RENAME TO d2.t2;")
	if c.databases["d"].Tables["t"] != nil {
		t.Errorf("d.t is still there")
	}
	if c.databases["d"].Tables["d2.t2"] != nil {
		t.Errorf("a table named \"d2.t2\" was made in d")
	}
	moved := c.databases["d2"].Tables["t2"]
	if moved == nil {
		t.Fatalf("d2.t2 was not created; d2 holds %v", c.databases["d2"].Tables)
	}
	if moved.Database == nil || moved.Database.Name != "d2" {
		t.Errorf("the moved table's database = %v, want d2", moved.Database)
	}
}

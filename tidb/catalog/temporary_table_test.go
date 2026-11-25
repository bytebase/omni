package catalog

import (
	"strings"
	"testing"
)

// TestLocalTemporaryTables holds the catalog to what TiDB v8.5.5 does with
// local and global temporary tables, each case read on the server: a local
// temporary table has a namespace of its own, shadows a table of its name,
// and refuses the DDL TiDB does not run on it.
func TestLocalTemporaryTables(t *testing.T) {
	setup := `CREATE TABLE t (b INT PRIMARY KEY, d INT, KEY idx_d (d));
CREATE TABLE pt (a INT) PARTITION BY HASH (a) PARTITIONS 2;
CREATE VIEW v AS SELECT 1 AS a;
CREATE SEQUENCE sq;
CREATE GLOBAL TEMPORARY TABLE g (a INT) ON COMMIT DELETE ROWS;
CREATE TEMPORARY TABLE tmp (a INT, KEY ia (a))`
	cases := []struct {
		name  string
		steps []string
		code  int // of the last step; 0 is accepted
	}{
		{"shadows a table", []string{"CREATE TEMPORARY TABLE t LIKE t"}, 0},
		{"shares a view's name", []string{"CREATE TEMPORARY TABLE v (a INT)"}, 0},
		{"shares a sequence's name", []string{"CREATE TEMPORARY TABLE sq (a INT)"}, 0},
		{"shares a global temporary table's name", []string{"CREATE TEMPORARY TABLE g (a INT)"}, 0},
		{"a table beside it", []string{"CREATE TABLE tmp (a INT)"}, 0},
		{"another of its name", []string{"CREATE TEMPORARY TABLE tmp (a INT)"}, ErrDupTable},
		{"another under IF NOT EXISTS", []string{"CREATE TEMPORARY TABLE IF NOT EXISTS tmp (a INT)"}, 0},
		{"partitions", []string{"CREATE TEMPORARY TABLE p (a INT) PARTITION BY HASH (a) PARTITIONS 2"}, ErrPartitionNoTemporary},
		{"global with partitions", []string{"CREATE GLOBAL TEMPORARY TABLE p (a INT) PARTITION BY HASH (a) PARTITIONS 2 ON COMMIT DELETE ROWS"}, ErrPartitionNoTemporary},
		{"LIKE a partitioned table", []string{"CREATE TEMPORARY TABLE p LIKE pt"}, ErrPartitionNoTemporary},
		{"global LIKE a partitioned table", []string{"CREATE GLOBAL TEMPORARY TABLE p LIKE pt ON COMMIT DELETE ROWS"}, ErrPartitionNoTemporary},
		{"LIKE a view", []string{"CREATE TEMPORARY TABLE x LIKE v"}, ErrWrongObject},
		{"table LIKE it", []string{"CREATE TABLE x LIKE tmp"}, ErrOptOnTemporaryTable},
		{"temporary table LIKE it", []string{"CREATE TEMPORARY TABLE x LIKE tmp"}, ErrOptOnTemporaryTable},
		{"LIKE a table it shadows", []string{"CREATE TEMPORARY TABLE t LIKE t", "CREATE TABLE x LIKE t"}, ErrOptOnTemporaryTable},
		{"ALTER TABLE", []string{"ALTER TABLE tmp ADD COLUMN z INT"}, ErrUnsupportedDDL},
		{"ALTER TABLE of a shadowed table", []string{"CREATE TEMPORARY TABLE t LIKE t", "ALTER TABLE t ADD COLUMN z INT"}, ErrUnsupportedDDL},
		{"CREATE INDEX", []string{"CREATE INDEX iz ON tmp (a)"}, ErrUnsupportedDDL},
		{"DROP INDEX", []string{"DROP INDEX ia ON tmp"}, ErrUnsupportedDDL},
		{"RENAME TABLE", []string{"RENAME TABLE tmp TO tmp2"}, ErrUnsupportedDDL},
		{"TRUNCATE", []string{"TRUNCATE TABLE tmp"}, 0},
		{"DROP TEMPORARY TABLE", []string{"DROP TEMPORARY TABLE tmp"}, 0},
		{"DROP TEMPORARY TABLE of a table", []string{"DROP TEMPORARY TABLE t"}, ErrUnknownTable},
		{"DROP TEMPORARY TABLE of a view", []string{"DROP TEMPORARY TABLE v"}, ErrUnknownTable},
		{"DROP TEMPORARY TABLE of a global one", []string{"DROP TEMPORARY TABLE g"}, ErrUnknownTable},
		{"DROP TEMPORARY TABLE IF EXISTS", []string{"DROP TEMPORARY TABLE IF EXISTS nope"}, 0},
		{"DROP TEMPORARY TABLE with one missing", []string{"DROP TEMPORARY TABLE tmp, nope"}, ErrUnknownTable},
		{"DROP GLOBAL TEMPORARY TABLE", []string{"DROP GLOBAL TEMPORARY TABLE g"}, 0},
		{"DROP GLOBAL TEMPORARY TABLE of a table", []string{"DROP GLOBAL TEMPORARY TABLE t"}, ErrDropTableOnTemporaryTable},
		{"DROP GLOBAL TEMPORARY TABLE of a view", []string{"DROP GLOBAL TEMPORARY TABLE v"}, ErrDropTableOnTemporaryTable},
		{"DROP GLOBAL TEMPORARY TABLE of a sequence", []string{"DROP GLOBAL TEMPORARY TABLE sq"}, ErrDropTableOnTemporaryTable},
		{"DROP GLOBAL TEMPORARY TABLE of a local one", []string{"DROP GLOBAL TEMPORARY TABLE tmp"}, ErrDropTableOnTemporaryTable},
		{"DROP GLOBAL TEMPORARY TABLE of a shadowed one", []string{"CREATE TEMPORARY TABLE g (a INT)", "DROP GLOBAL TEMPORARY TABLE g"}, ErrDropTableOnTemporaryTable},
		{"DROP GLOBAL TEMPORARY TABLE of nothing", []string{"DROP GLOBAL TEMPORARY TABLE nope"}, ErrUnknownTable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			mustExec(t, c, "CREATE DATABASE test")
			c.SetCurrentDatabase("test")
			mustExec(t, c, setup)
			last := len(tc.steps) - 1
			for _, step := range tc.steps[:last] {
				mustExec(t, c, step)
			}
			if got := execCode(t, c, tc.steps[last]); got != tc.code {
				t.Errorf("%s: code %d, want %d", tc.steps[last], got, tc.code)
			}
		})
	}
}

// TestLocalTemporaryTableShadows checks what SHOW CREATE TABLE and DROP TABLE
// read when a local temporary table shares a table's name: the temporary
// table, until it is dropped (read on v8.5.5). The schema itself never holds
// it.
func TestLocalTemporaryTableShadows(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t (b INT PRIMARY KEY, d INT)")
	mustExec(t, c, "CREATE TEMPORARY TABLE t LIKE t")
	showsTemporary := func(want bool) {
		t.Helper()
		got := strings.HasPrefix(c.ShowCreateTable("test", "t"), "CREATE TEMPORARY TABLE ")
		if got != want {
			t.Fatalf("SHOW CREATE TABLE t: %s", c.ShowCreateTable("test", "t"))
		}
	}
	showsTemporary(true)
	db := c.GetDatabase("test")
	if tbl := db.Tables["t"]; tbl == nil || tbl.Temporary {
		t.Fatalf("the schema's t is %+v", tbl)
	}
	if meta := c.DatabaseMetadata("test"); len(meta.Schemas[0].Tables) != 1 {
		t.Fatalf("metadata lists %d tables, want the schema's one", len(meta.Schemas[0].Tables))
	}

	mustExec(t, c, "DROP TABLE t")
	showsTemporary(false)
	mustExec(t, c, "DROP TABLE t")
	if db.GetTable("t") != nil {
		t.Fatal("the second DROP TABLE left t")
	}
}

// TestLocalTemporaryTableLeavesSchemaAlone checks that a local temporary table
// shadowing a schema table changes nothing the schema answers: TRUNCATE of it
// is not held back by the shadowed table's foreign keys, and a view over the
// shadowed table keeps that table's column types in the metadata.
func TestLocalTemporaryTableLeavesSchemaAlone(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, `CREATE TABLE p (a INT PRIMARY KEY);
CREATE TABLE ch (a INT, FOREIGN KEY (a) REFERENCES p (a));
CREATE VIEW v AS SELECT a FROM p;
CREATE TEMPORARY TABLE p (a VARCHAR(20))`)
	if got := execCode(t, c, "TRUNCATE TABLE p"); got != 0 {
		t.Errorf("TRUNCATE TABLE of the temporary p: code %d, want 0", got)
	}
	meta := c.DatabaseMetadata("test")
	view := meta.Schemas[0].Views[0]
	if len(view.Columns) != 1 || view.Columns[0].Type != "int" {
		t.Errorf("v's columns are %v, want a of p's int", view.Columns)
	}
}

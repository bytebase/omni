package catalog

import (
	"strings"
	"testing"
)

// These tests assert what the ALTER TABLE clauses do to the catalog, not that
// they parse. A clause that parses and applies nothing looks like a closed gap
// in tidb/parser's corpus, which records only a verdict. None of them boots a
// container, so none is inside the TIDB_CATALOG_PARITY quarantine that gates
// the differential tests in container_test.go.

// execAll runs a script and fails on the first statement that errors.
func execAll(t *testing.T, c *Catalog, sql string) {
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

// execExpectError runs a script and returns the first statement error.
func execExpectError(t *testing.T, c *Catalog, sql string) error {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}

// newCatalogWithTable returns a catalog holding d.t with the given definition.
func newCatalogWithTable(t *testing.T, def string) *Catalog {
	t.Helper()
	c := New()
	execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t "+def+";")
	return c
}

// TestAlterTableAddColumnTableElementList: the parenthesized ADD reads a table
// element list, so the constraints written beside the columns apply too. An
// earlier form of this reader took only the columns, which parsed and added no
// index.
func TestAlterTableAddColumnTableElementList(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT)")
	execAll(t, c, "ALTER TABLE t ADD COLUMN (b INT, c INT, PRIMARY KEY (a), UNIQUE KEY u (b));")

	ddl := c.ShowCreateTable("d", "t")
	for _, want := range []string{"`b` int", "`c` int", "PRIMARY KEY (`a`)", "UNIQUE KEY `u` (`b`)"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}
}

// TestAlterTableAddColumnIfNotExists: IF NOT EXISTS makes a duplicate column a
// no-op rather than an error, and leaves the column that is there untouched.
func TestAlterTableAddColumnIfNotExists(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, b VARCHAR(10))")
	execAll(t, c, "ALTER TABLE t ADD COLUMN IF NOT EXISTS b INT;")
	ddl := c.ShowCreateTable("d", "t")
	if !strings.Contains(ddl, "`b` varchar(10)") {
		t.Errorf("IF NOT EXISTS replaced the existing column:\n%s", ddl)
	}

	// Without IF NOT EXISTS the same statement is an error.
	if err := execExpectError(t, c, "ALTER TABLE t ADD COLUMN b INT;"); err == nil {
		t.Errorf("ADD COLUMN b on a table that has b was accepted")
	}

	// A column that is not there is still added.
	execAll(t, c, "ALTER TABLE t ADD COLUMN IF NOT EXISTS zz INT;")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "`zz` int") {
		t.Errorf("ADD COLUMN IF NOT EXISTS zz added nothing:\n%s", ddl)
	}
}

// TestAlterTableModifyChangeColumnIfExists: IF EXISTS makes a missing column a
// no-op rather than an error.
func TestAlterTableModifyChangeColumnIfExists(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT)")
	execAll(t, c, "ALTER TABLE t MODIFY COLUMN IF EXISTS zz VARCHAR(5); ALTER TABLE t CHANGE COLUMN IF EXISTS zz yy VARCHAR(5);")
	ddl := c.ShowCreateTable("d", "t")
	for _, unwanted := range []string{"`zz`", "`yy`"} {
		if strings.Contains(ddl, unwanted) {
			t.Errorf("the IF EXISTS form invented %s:\n%s", unwanted, ddl)
		}
	}

	// A column that is there is still modified.
	execAll(t, c, "ALTER TABLE t MODIFY COLUMN IF EXISTS a BIGINT;")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "`a` bigint") {
		t.Errorf("MODIFY COLUMN IF EXISTS a did not change a:\n%s", ddl)
	}

	// Without IF EXISTS the missing column is an error.
	if err := execExpectError(t, c, "ALTER TABLE t MODIFY COLUMN zz INT;"); err == nil {
		t.Errorf("MODIFY COLUMN zz on a table without zz was accepted")
	}
}

// TestAlterTableDropIndexIfExists: IF EXISTS makes a missing index a no-op,
// and the index that is there still goes.
func TestAlterTableDropIndexIfExists(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, KEY i (a))")
	execAll(t, c, "ALTER TABLE t DROP INDEX IF EXISTS nosuch;")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "KEY `i`") {
		t.Errorf("DROP INDEX IF EXISTS nosuch took the wrong index:\n%s", ddl)
	}
	execAll(t, c, "ALTER TABLE t DROP INDEX IF EXISTS i;")
	if ddl := c.ShowCreateTable("d", "t"); strings.Contains(ddl, "KEY `i`") {
		t.Errorf("DROP INDEX IF EXISTS i left the index:\n%s", ddl)
	}
	if err := execExpectError(t, c, "ALTER TABLE t DROP INDEX nosuch;"); err == nil {
		t.Errorf("DROP INDEX nosuch was accepted")
	}
}

// TestAlterTableAddPartitionIfNotExists: IF NOT EXISTS skips a partition the
// table already has instead of adding a second one of the same name.
func TestAlterTableAddPartitionIfNotExists(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))")
	execAll(t, c, "ALTER TABLE t ADD PARTITION IF NOT EXISTS (PARTITION p0 VALUES LESS THAN (20), PARTITION p1 VALUES LESS THAN (30));")

	tbl := c.databases["d"].Tables["t"]
	if tbl.Partitioning == nil {
		t.Fatal("the table lost its partitioning")
	}
	var names []string
	for _, pd := range tbl.Partitioning.Partitions {
		names = append(names, pd.Name)
	}
	if len(names) != 2 || names[0] != "p0" || names[1] != "p1" {
		t.Errorf("partitions = %v, want [p0 p1]", names)
	}
	// p0 keeps the bound it was created with.
	if got := tbl.Partitioning.Partitions[0].ValueExpr; !strings.Contains(got, "10") {
		t.Errorf("p0 bound = %q, want the original 10", got)
	}
}

// TestAlterTableDropPartitionIfExists: IF EXISTS drops what is there and
// ignores what is not.
func TestAlterTableDropPartitionIfExists(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))")
	execAll(t, c, "ALTER TABLE t DROP PARTITION IF EXISTS p0, nosuch;")
	tbl := c.databases["d"].Tables["t"]
	if n := len(tbl.Partitioning.Partitions); n != 1 || tbl.Partitioning.Partitions[0].Name != "p1" {
		t.Errorf("partitions = %v, want [p1]", tbl.Partitioning.Partitions)
	}
}

// TestAlterTableRenameToQualifiedName: the destination of a RENAME may carry a
// database, and then the table moves there. Reading the name as written would
// have made a table called "d2.t2" in the database the statement started in.
func TestAlterTableRenameToQualifiedName(t *testing.T) {
	for _, sep := range []string{"TO", "AS", "=", ""} {
		t.Run("sep="+sep, func(t *testing.T) {
			c := New()
			execAll(t, c, "CREATE DATABASE d; CREATE DATABASE d2; USE d; CREATE TABLE t (a INT);")
			execAll(t, c, "ALTER TABLE t RENAME "+sep+" d2.t2;")

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
			if moved.Name != "t2" {
				t.Errorf("the moved table is named %q, want \"t2\"", moved.Name)
			}
			if moved.Database == nil || moved.Database.Name != "d2" {
				t.Errorf("the moved table's database = %v, want d2", moved.Database)
			}
		})
	}

	// An unqualified destination still renames inside the same database.
	c := New()
	execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT); ALTER TABLE t RENAME = t1;")
	if c.databases["d"].Tables["t1"] == nil {
		t.Errorf("RENAME = t1 did not rename inside d; d holds %v", c.databases["d"].Tables)
	}
}

// TestAlterTablePartDefStorageEngine: STORAGE ENGINE is the other spelling of
// a partition's ENGINE option and names the same engine.
func TestAlterTablePartDefStorageEngine(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))")
	execAll(t, c, "ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) STORAGE ENGINE = InnoDB);")
	tbl := c.databases["d"].Tables["t"]
	last := tbl.Partitioning.Partitions[len(tbl.Partitioning.Partitions)-1]
	if last.Name != "p1" {
		t.Fatalf("the last partition is %q, want p1", last.Name)
	}
	if last.Engine != "InnoDB" {
		t.Errorf("p1 engine = %q, want \"InnoDB\"", last.Engine)
	}
}

// TestAlterTableForceAutoIncrement: FORCE AUTO_INCREMENT is a table option and
// applies the value, as the plain option does.
func TestAlterTableForceAutoIncrement(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT AUTO_INCREMENT PRIMARY KEY)")
	execAll(t, c, "ALTER TABLE t FORCE AUTO_INCREMENT = 50;")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "AUTO_INCREMENT=50") {
		t.Errorf("FORCE AUTO_INCREMENT = 50 applied nothing:\n%s", ddl)
	}
}

// TestAlterTableGlobalLocalIndexOption: an index option this parser newly
// reads must not displace the options beside it, and CLUSTERED still reaches
// the constraint it describes.
func TestAlterTableGlobalLocalIndexOption(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, b INT)")
	execAll(t, c, "ALTER TABLE t ADD INDEX i (b) GLOBAL COMMENT 'kept';")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "COMMENT 'kept'") {
		t.Errorf("GLOBAL swallowed the COMMENT beside it:\n%s", ddl)
	}

	c2 := newCatalogWithTable(t, "(a INT)")
	execAll(t, c2, "ALTER TABLE t ADD PRIMARY KEY (a) NONCLUSTERED GLOBAL;")
	tbl := c2.databases["d"].Tables["t"]
	var pk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey {
			pk = con
		}
	}
	if pk == nil {
		t.Fatalf("no primary key among %v", tbl.Constraints)
	}
	if pk.Clustered == nil || *pk.Clustered {
		t.Errorf("Clustered = %v, want false for NONCLUSTERED", pk.Clustered)
	}
}

// TestAlterTableNoEffectClauses: the TiDB actions that change no schema the
// catalog models must still leave the table intact and raise no error. An
// action that errored here would make every statement carrying it fail.
func TestAlterTableNoEffectClauses(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE t CACHE",
		"ALTER TABLE t NOCACHE",
		"ALTER TABLE t CACHE, NOCACHE",
		"ALTER TABLE t READ ONLY",
		"ALTER TABLE t READ WRITE",
		"ALTER TABLE t COMPACT",
		"ALTER TABLE t COMPACT TIFLASH REPLICA",
		"ALTER TABLE t ATTRIBUTES = 'str1,str2'",
		"ALTER TABLE t ATTRIBUTES = DEFAULT",
		"ALTER TABLE t STATS_OPTIONS = 'str'",
		"ALTER TABLE t ANALYZE PARTITION p0 WITH 4 BUCKETS",
		"ALTER TABLE t PARTITION p0 ATTRIBUTES = 'str'",
		"ALTER TABLE t PARTITION p0 PLACEMENT POLICY = x",
		"ALTER TABLE t LAST PARTITION LESS THAN (1000)",
		"ALTER TABLE t FIRST PARTITION LESS THAN (1)",
		"ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1000)",
		"ALTER TABLE t MERGE FIRST PARTITION LESS THAN (1000)",
	} {
		t.Run(sql, func(t *testing.T) {
			c := newCatalogWithTable(t, "(a INT, b INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))")
			before := c.ShowCreateTable("d", "t")
			execAll(t, c, sql+";")
			if after := c.ShowCreateTable("d", "t"); after != before {
				t.Errorf("%q changed the table:\nbefore:\n%s\nafter:\n%s", sql, before, after)
			}
		})
	}
}

// TestAlterTableClausesKeepTheRestOfTheStatement: a clause in a multi-clause
// statement must not swallow the one beside it. The no-effect clauses are the
// risk: a reader that consumed one token too many would lose the action after
// it without failing.
func TestAlterTableClausesKeepTheRestOfTheStatement(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT)")
	execAll(t, c, "ALTER TABLE t CACHE, ADD COLUMN b INT, READ ONLY, ADD COLUMN c INT, NOCACHE;")
	ddl := c.ShowCreateTable("d", "t")
	for _, want := range []string{"`b` int", "`c` int"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}

	c = newCatalogWithTable(t, "(a INT)")
	execAll(t, c, "ALTER TABLE t ATTRIBUTES = 'x', ADD COLUMN b INT, STATS_OPTIONS = 'y', COMMENT 'c';")
	ddl = c.ShowCreateTable("d", "t")
	for _, want := range []string{"`b` int", "COMMENT='c'"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}
}

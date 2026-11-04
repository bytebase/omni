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

// TestAlterTableAddIndexIfNotExists: IF NOT EXISTS on an index makes a
// duplicate name a no-op rather than 1061, and keeps the index that is there.
// TiDB v8.5.5 raises note 1061 and leaves the table alone.
func TestAlterTableAddIndexIfNotExists(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, b INT, KEY i (a))")
	execAll(t, c, "ALTER TABLE t ADD INDEX IF NOT EXISTS i (b);")
	ddl := c.ShowCreateTable("d", "t")
	if !strings.Contains(ddl, "KEY `i` (`a`)") {
		t.Errorf("IF NOT EXISTS replaced the existing index:\n%s", ddl)
	}
	if strings.Count(ddl, "KEY `i`") != 1 {
		t.Errorf("IF NOT EXISTS added a second index of the name:\n%s", ddl)
	}

	// Without it the same statement is an error.
	if err := execExpectError(t, c, "ALTER TABLE t ADD INDEX i (b);"); err == nil {
		t.Errorf("ADD INDEX i on a table that has i was accepted")
	}

	// A name that is not there is still added.
	execAll(t, c, "ALTER TABLE t ADD KEY IF NOT EXISTS j (b);")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "KEY `j` (`b`)") {
		t.Errorf("ADD KEY IF NOT EXISTS j added nothing:\n%s", ddl)
	}
}

// TestAlterTableAddVectorIndex: ADD VECTOR INDEX adds an index. It parsed and
// applied nothing before, which reads as a successful no-op for a statement
// that changes the schema.
func TestAlterTableAddVectorIndex(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, b INT)")
	execAll(t, c, "ALTER TABLE t ADD VECTOR INDEX idx ((lower(a))) USING HNSW;")

	tbl := c.databases["d"].Tables["t"]
	var vec *Index
	for _, idx := range tbl.Indexes {
		if idx.Vector {
			vec = idx
		}
	}
	if vec == nil {
		t.Fatalf("no vector index among %v", tbl.Indexes)
	}
	if vec.Name != "idx" {
		t.Errorf("Name = %q, want \"idx\"", vec.Name)
	}
	if vec.IndexType != "HNSW" {
		t.Errorf("IndexType = %q, want \"HNSW\"", vec.IndexType)
	}
	if len(vec.Columns) != 1 {
		t.Errorf("Columns = %v, want one key part", vec.Columns)
	}
	// SHOW CREATE TABLE writes the word TiDB writes: ConstraintVector in
	// v8.5.5 pkg/parser/ast/ddl.go restores `VECTOR INDEX`.
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "VECTOR INDEX `idx`") {
		t.Errorf("SHOW CREATE TABLE lacks the vector index:\n%s", ddl)
	}

	// A duplicate name is 1061, and IF NOT EXISTS makes it a no-op.
	if err := execExpectError(t, c, "ALTER TABLE t ADD VECTOR INDEX idx ((lower(b)));"); err == nil {
		t.Errorf("a second vector index named idx was accepted")
	}
	execAll(t, c, "ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS idx ((lower(b)));")

	// CREATE TABLE declares one too.
	c2 := New()
	execAll(t, c2, "CREATE DATABASE d2; USE d2; CREATE TABLE t (a INT, VECTOR INDEX v ((lower(a))));")
	if ddl := c2.ShowCreateTable("d2", "t"); !strings.Contains(ddl, "VECTOR INDEX `v`") {
		t.Errorf("CREATE TABLE lost the vector index:\n%s", ddl)
	}
}

// TestAlterTableRenameRollback: a cross-database rename beside an action that
// fails must leave the table where it was, and in one database only. The
// rollback repaired the starting database's map alone, so the destination kept
// an entry pointing at the restored table -- and when the move kept the name,
// the starting database was left without it.
func TestAlterTableRenameRollback(t *testing.T) {
	for _, c := range []struct {
		name string
		sql  string
		want string // the name the table keeps in d
	}{
		{"a new name in another database", "ALTER TABLE t RENAME TO e.u, DROP COLUMN missing;", "t"},
		{"the same name in another database", "ALTER TABLE t RENAME TO e.t, DROP COLUMN missing;", "t"},
		{"a new name in the same database", "ALTER TABLE t RENAME TO u, DROP COLUMN missing;", "t"},
	} {
		t.Run(c.name, func(t *testing.T) {
			cat := New()
			execAll(t, cat, "CREATE DATABASE d; CREATE DATABASE e; USE d; CREATE TABLE t (a INT);")
			if err := execExpectError(t, cat, c.sql); err == nil {
				t.Fatalf("%q was accepted", c.sql)
			}
			d, e := cat.databases["d"], cat.databases["e"]
			if len(e.Tables) != 0 {
				t.Errorf("e holds %v after the rollback, want nothing", e.Tables)
			}
			if len(d.Tables) != 1 || d.Tables[c.want] == nil {
				t.Fatalf("d holds %v, want just %q", d.Tables, c.want)
			}
			tbl := d.Tables[c.want]
			if tbl.Name != "t" {
				t.Errorf("the table is named %q, want \"t\"", tbl.Name)
			}
			if tbl.Database == nil || tbl.Database.Name != "d" {
				t.Errorf("the table's database = %v, want d", tbl.Database)
			}
			// The catalog must still work on it afterwards.
			execAll(t, cat, "ALTER TABLE d.t ADD COLUMN b INT;")
			if ddl := cat.ShowCreateTable("d", "t"); !strings.Contains(ddl, "`b` int") {
				t.Errorf("the restored table does not take a later ALTER:\n%s", ddl)
			}
		})
	}
}

// TestAlterTableRenameToCurrentDatabase: a destination with no qualifier lands
// in the current database, not in the renamed table's. `USE e;
// ALTER TABLE d.t RENAME TO u` leaves u in e on TiDB v8.5.5 and MySQL 8.0.46.
func TestAlterTableRenameToCurrentDatabase(t *testing.T) {
	c := New()
	execAll(t, c, "CREATE DATABASE d; CREATE DATABASE e; CREATE TABLE d.t (a INT); USE e;")
	execAll(t, c, "ALTER TABLE d.t RENAME TO u;")
	if c.databases["d"].Tables["t"] != nil {
		t.Errorf("d.t is still there")
	}
	if c.databases["d"].Tables["u"] != nil {
		t.Errorf("u landed in d, not in the current database")
	}
	if c.databases["e"].Tables["u"] == nil {
		t.Fatalf("e.u was not created; e holds %v", c.databases["e"].Tables)
	}
}

// TestAlterTableAddTableElementListIsAtomic: the parenthesized ADD is one
// command carrying several table elements, so a failure part-way through must
// not leave the earlier ones installed. TiDB v8.5.5 leaves the table
// untouched for `ALTER TABLE t ADD (b INT, KEY i (b))` against a table that
// already has i (verified on 2026-11-04); this used to take the
// single-command fast path and keep b.
func TestAlterTableAddTableElementListIsAtomic(t *testing.T) {
	for _, c := range []struct {
		name string
		sql  string
	}{
		{"a constraint fails after a column", "ALTER TABLE t ADD (b INT, KEY i (b));"},
		{"a column fails after a column", "ALTER TABLE t ADD (b INT, a INT);"},
		{"a constraint fails after a constraint", "ALTER TABLE t ADD (KEY j (a), KEY i (a));"},
	} {
		t.Run(c.name, func(t *testing.T) {
			cat := New()
			execAll(t, cat, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, KEY i (a));")
			before := cat.ShowCreateTable("d", "t")
			if err := execExpectError(t, cat, c.sql); err == nil {
				t.Fatalf("%q was accepted", c.sql)
			}
			if after := cat.ShowCreateTable("d", "t"); after != before {
				t.Errorf("%q left the table changed:\nbefore:\n%s\nafter:\n%s", c.sql, before, after)
			}
		})
	}
}

// TestAlterTableGlobalUniqueIndex: GLOBAL exempts a unique key from covering
// the partitioning columns, which is the whole point of the option. v8.5.5
// runs `ADD UNIQUE INDEX i (b) GLOBAL` on a table partitioned by a and
// answers 8264 without it, and writes the index back with its own version
// comment. Reading GLOBAL as an index option without carrying it into the
// catalog left the statement failing with the error the option avoids.
func TestAlterTableGlobalUniqueIndex(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, b INT, c INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))")
	execAll(t, c, "ALTER TABLE t ADD UNIQUE INDEX i (b) GLOBAL COMMENT 'kept';")

	ddl := c.ShowCreateTable("d", "t")
	// The comment comes last, after COMMENT, as v8.5.5 writes it.
	if want := "UNIQUE KEY `i` (`b`) COMMENT 'kept' /*T![global_index] GLOBAL */"; !strings.Contains(ddl, want) {
		t.Errorf("SHOW CREATE TABLE lacks %q:\n%s", want, ddl)
	}

	// A non-unique index takes GLOBAL too.
	execAll(t, c, "ALTER TABLE t ADD INDEX j (c) GLOBAL;")
	if want := "KEY `j` (`c`) /*T![global_index] GLOBAL */"; !strings.Contains(c.ShowCreateTable("d", "t"), want) {
		t.Errorf("SHOW CREATE TABLE lacks %q:\n%s", want, c.ShowCreateTable("d", "t"))
	}

	// Without the option the same key is still 8264, so the exemption is the
	// option's and not a hole in the check.
	err := execExpectError(t, c, "ALTER TABLE t ADD UNIQUE INDEX k (b);")
	if err == nil {
		t.Fatalf("ADD UNIQUE INDEX k (b) without GLOBAL was accepted")
	}
	if !strings.Contains(err.Error(), "8264") {
		t.Errorf("error = %v, want 8264", err)
	}

	// LOCAL is the default, writes nothing, and keeps the requirement.
	if err := execExpectError(t, c, "ALTER TABLE t ADD UNIQUE INDEX l (b) LOCAL;"); err == nil {
		t.Errorf("ADD UNIQUE INDEX l (b) LOCAL was accepted")
	}
	// A later LOCAL after GLOBAL wins, as the grammar's action does.
	if err := execExpectError(t, c, "ALTER TABLE t ADD UNIQUE INDEX m (b) GLOBAL LOCAL;"); err == nil {
		t.Errorf("GLOBAL LOCAL was accepted; the last option must win")
	}
}

// TestAlterTableForeignKeyIfNotExistsNamespace: a foreign key's IF NOT EXISTS
// is not about the index namespace. v8.5.5 runs
// `ADD CONSTRAINT fk FOREIGN KEY IF NOT EXISTS (a) REFERENCES p (a)` against
// a table that already has `KEY fk (a)` and keeps both, and answers 1826 for a
// foreign key of the name even with IF NOT EXISTS written. Reading the index
// namespace here dropped the foreign key and reported success.
func TestAlterTableForeignKeyIfNotExistsNamespace(t *testing.T) {
	c := New()
	execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE p (a INT PRIMARY KEY); CREATE TABLE t (a INT, KEY fk (a));")
	execAll(t, c, "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY IF NOT EXISTS (a) REFERENCES p (a);")

	tbl := c.databases["d"].Tables["t"]
	var fk *Constraint
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			fk = con
		}
	}
	if fk == nil {
		t.Fatalf("the foreign key was dropped; constraints = %v", tbl.Constraints)
	}
	if fk.Name != "fk" {
		t.Errorf("Name = %q, want \"fk\"", fk.Name)
	}
	ddl := c.ShowCreateTable("d", "t")
	for _, want := range []string{"KEY `fk` (`a`)", "CONSTRAINT `fk` FOREIGN KEY (`a`) REFERENCES `p` (`a`)"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}

	// A foreign key of the name is not covered by the flag either -- v8.5.5
	// answers 1826 for a second one even with IF NOT EXISTS written -- but
	// this catalog models no foreign-key name namespace at all and so raises
	// no 1826 for a duplicate with or without the flag. That gap is older
	// than this family and is not asserted here; what matters for the flag is
	// that it no longer reads the index namespace for a foreign key.

	// An index of the name is still covered, which is what the flag is for.
	execAll(t, c, "ALTER TABLE t ADD INDEX IF NOT EXISTS fk (a);")
	if got := strings.Count(c.ShowCreateTable("d", "t"), "KEY `fk` (`a`)"); got != 1 {
		t.Errorf("a second index named fk was added (%d):\n%s", got, c.ShowCreateTable("d", "t"))
	}
}

// TestAlterTableUpdateIndexes: UPDATE INDEXES says what each named index
// becomes once the table is repartitioned, so it has to apply before the
// coverage check. v8.5.5 runs the repartition with the clause and answers 8264
// without it, and reads the index back as GLOBAL.
func TestAlterTableUpdateIndexes(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT, b INT, UNIQUE KEY i (b))")
	execAll(t, c, "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES (i GLOBAL);")
	if want := "UNIQUE KEY `i` (`b`) /*T![global_index] GLOBAL */"; !strings.Contains(c.ShowCreateTable("d", "t"), want) {
		t.Errorf("SHOW CREATE TABLE lacks %q:\n%s", want, c.ShowCreateTable("d", "t"))
	}

	// Without the clause the same repartition is 8264.
	c2 := newCatalogWithTable(t, "(a INT, b INT, UNIQUE KEY i (b))")
	err := execExpectError(t, c2, "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2;")
	if err == nil {
		t.Fatalf("the repartition without UPDATE INDEXES was accepted")
	}
	if !strings.Contains(err.Error(), "8264") {
		t.Errorf("error = %v, want 8264", err)
	}

	// LOCAL keeps the requirement, so the clause is read and not just
	// consumed.
	c3 := newCatalogWithTable(t, "(a INT, b INT, UNIQUE KEY i (b))")
	if err := execExpectError(t, c3, "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES (i LOCAL);"); err == nil {
		t.Errorf("UPDATE INDEXES (i LOCAL) was accepted")
	}

	// A name the clause gives twice is 1280 as well -- v8.5.5 refuses
	// `(i GLOBAL, i LOCAL)` and even `(i GLOBAL, i GLOBAL)` -- so the list
	// names each index at most once and the undo cannot depend on its order.
	c5 := newCatalogWithTable(t, "(a INT, b INT, UNIQUE KEY i (b))")
	err = execExpectError(t, c5, "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES (i GLOBAL, i LOCAL);")
	if err == nil {
		t.Fatalf("UPDATE INDEXES (i GLOBAL, i LOCAL) was accepted")
	}
	if !strings.Contains(err.Error(), "1280") {
		t.Errorf("error = %v, want 1280", err)
	}
	if strings.Contains(c5.ShowCreateTable("d", "t"), "GLOBAL") {
		t.Errorf("the refused statement left a GLOBAL behind:\n%s", c5.ShowCreateTable("d", "t"))
	}

	// Two distinct names both apply: v8.5.5 runs
	// `UPDATE INDEXES (i GLOBAL, j GLOBAL)` and reads both back as GLOBAL.
	c6 := newCatalogWithTable(t, "(a INT, b INT, c INT, UNIQUE KEY i (b), UNIQUE KEY j (c))")
	execAll(t, c6, "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES (i GLOBAL, j GLOBAL);")
	if got := strings.Count(c6.ShowCreateTable("d", "t"), "/*T![global_index] GLOBAL */"); got != 2 {
		t.Errorf("%d indexes read back GLOBAL, want 2:\n%s", got, c6.ShowCreateTable("d", "t"))
	}

	// A name the table does not have is 1280, and the refused statement
	// leaves no setting behind.
	c4 := newCatalogWithTable(t, "(a INT, b INT, UNIQUE KEY i (b))")
	err = execExpectError(t, c4, "ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES (nosuch GLOBAL);")
	if err == nil {
		t.Fatalf("UPDATE INDEXES (nosuch GLOBAL) was accepted")
	}
	if !strings.Contains(err.Error(), "1280") {
		t.Errorf("error = %v, want 1280", err)
	}
	if strings.Contains(c4.ShowCreateTable("d", "t"), "GLOBAL") {
		t.Errorf("the refused statement left a GLOBAL behind:\n%s", c4.ShowCreateTable("d", "t"))
	}
}

// TestAlterTableRenameLeavesItsCurrentDatabase: a second rename in the same
// statement starts from where the first one left the table, not from the
// database the statement began in. TiDB v8.5.5 answers 8200 to two renames in
// one ALTER, but MySQL 8.0.46 and MariaDB 11.8.8 both run it and end with one
// entry, and reading the starting database here left the intermediate
// destination installed.
func TestAlterTableRenameLeavesItsCurrentDatabase(t *testing.T) {
	c := New()
	execAll(t, c, "CREATE DATABASE d; CREATE DATABASE e; USE d; CREATE TABLE t (a INT);")
	execAll(t, c, "ALTER TABLE t RENAME TO e.u, RENAME TO e.v;")

	d, e := c.databases["d"], c.databases["e"]
	if len(d.Tables) != 0 {
		t.Errorf("d holds %v, want nothing", d.Tables)
	}
	if len(e.Tables) != 1 || e.Tables["v"] == nil {
		t.Fatalf("e holds %v, want just v", e.Tables)
	}
	if e.Tables["u"] != nil {
		t.Errorf("the intermediate name u is still in e")
	}
	if tbl := e.Tables["v"]; tbl.Name != "v" || tbl.Database == nil || tbl.Database.Name != "e" {
		t.Errorf("e.v = %q in %v, want v in e", tbl.Name, tbl.Database)
	}
}

// TestAlterTableAddColumnWithConstraintIsAtomic: a column's constraints are
// applied after the column itself is inserted, and they can fail -- a second
// primary key, and in the forks a duplicate CHECK name. The servers leave the
// table untouched for both, and a single-column ADD took the no-rollback fast
// path, so the column survived the error.
func TestAlterTableAddColumnWithConstraintIsAtomic(t *testing.T) {
	c := newCatalogWithTable(t, "(a INT PRIMARY KEY)")
	before := c.ShowCreateTable("d", "t")
	if err := execExpectError(t, c, "ALTER TABLE t ADD b INT PRIMARY KEY;"); err == nil {
		t.Fatalf("a second primary key was accepted")
	}
	if after := c.ShowCreateTable("d", "t"); after != before {
		t.Errorf("the failed ADD left the table changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}

	// A column with no constraint keeps the fast path and still applies.
	execAll(t, c, "ALTER TABLE t ADD c INT;")
	if !strings.Contains(c.ShowCreateTable("d", "t"), "`c` int") {
		t.Errorf("ADD c INT applied nothing:\n%s", c.ShowCreateTable("d", "t"))
	}
}

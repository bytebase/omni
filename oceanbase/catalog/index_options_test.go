package catalog

import (
	"errors"
	"strings"
	"testing"
)

// These hold what the index options leave in the catalog. The expected text is what SHOW
// CREATE TABLE printed on 4.3.5.6 (read on 2026-11-05), past the BLOCK_SIZE and scope it
// prints on every key; TestOracle_OceanBaseIndexOptions holds it to the server.

func TestStoringShowCreate(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c))", "KEY `k` (`a`) STORING (`b`, `c`)"},
		{"CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b) STORING (c))", "KEY `k` (`a`) STORING (`b`, `c`)"},
		{"CREATE TABLE t (a INT, b INT, KEY k (a) STORING (B))", "KEY `k` (`a`) STORING (`b`)"},
		{"CREATE TABLE t (a INT, b INT, UNIQUE KEY k (a) STORING (b))", "UNIQUE KEY `k` (`a`) STORING (`b`)"},
		{"CREATE TABLE t (a INT, b INT); CREATE INDEX k ON t (a) STORING (b) PARALLEL 2", "KEY `k` (`a`) STORING (`b`)"},
		{"CREATE TABLE t (a INT, b INT); ALTER TABLE t ADD INDEX k (a) STORING (b) NOPARALLEL", "KEY `k` (`a`) STORING (`b`)"},
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; "+tc.sql+";")
		if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, tc.want) {
			t.Errorf("%s: SHOW CREATE TABLE lacks %s:\n%s", tc.sql, tc.want, ddl)
		}
	}
}

// TestOptionsThatLeaveNoTrace: these change nothing SHOW CREATE TABLE prints.
func TestOptionsThatLeaveNoTrace(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, KEY k (a) PARALLEL 3, KEY k2 (a) NOPARALLEL, KEY k3 (a) WITH ROWID);")
	ddl := c.ShowCreateTable("d", "t")
	for _, not := range []string{"PARALLEL", "ROWID", "STORING"} {
		if strings.Contains(ddl, not) {
			t.Errorf("SHOW CREATE TABLE has %s:\n%s", not, ddl)
		}
	}
}

// TestStoringFollowsItsColumns: dropping a stored column keeps the index, renaming one
// renames it, and LIKE copies it.
func TestStoringFollowsItsColumns(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c));")
	mustExec(t, c, "CREATE TABLE u LIKE t;")
	if ddl := c.ShowCreateTable("d", "u"); !strings.Contains(ddl, "STORING (`b`, `c`)") {
		t.Errorf("LIKE lost STORING:\n%s", ddl)
	}
	mustExec(t, c, "ALTER TABLE t CHANGE b e INT;")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "STORING (`e`, `c`)") {
		t.Errorf("CHANGE did not rename the stored column:\n%s", ddl)
	}
	mustExec(t, c, "ALTER TABLE t DROP COLUMN e;")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "KEY `k` (`a`) STORING (`c`)") {
		t.Errorf("DROP COLUMN:\n%s", ddl)
	}
	mustExec(t, c, "ALTER TABLE t DROP COLUMN c;")
	ddl := c.ShowCreateTable("d", "t")
	if !strings.Contains(ddl, "KEY `k` (`a`)") || strings.Contains(ddl, "STORING") {
		t.Errorf("the last stored column:\n%s", ddl)
	}
	// LIKE copied the list and does not share it.
	if ddl := c.ShowCreateTable("d", "u"); !strings.Contains(ddl, "STORING (`b`, `c`)") {
		t.Errorf("the copy followed the original:\n%s", ddl)
	}
}

// TestStoringSurvivesAFailedAlter holds the rollback snapshot to its own copy of the
// stored columns: a statement whose later command fails leaves the index as it was,
// whether an earlier command renamed a stored column or dropped one.
func TestStoringSurvivesAFailedAlter(t *testing.T) {
	for _, stmt := range []string{
		"ALTER TABLE t CHANGE b x INT, DROP COLUMN nosuch",
		"ALTER TABLE t DROP COLUMN b, DROP COLUMN nosuch",
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c));")
		if err := execError(t, c, stmt); err == nil {
			t.Fatalf("%s: no error", stmt)
		}
		if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "KEY `k` (`a`) STORING (`b`, `c`)") {
			t.Errorf("%s changed the index:\n%s", stmt, ddl)
		}
	}
}

func execError(t *testing.T, c *Catalog, sql string) error {
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

func TestIndexOptionRefusals(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"CREATE INDEX i ON t (a) STORING (nosuch);", 1054},
		{"ALTER TABLE t ADD INDEX i (a) STORING (nosuch);", 1054},
		{"CREATE TABLE u (a INT, KEY k (a) STORING (nosuch));", 1054},
		{"CREATE INDEX nosuch.i ON t (a);", 1049},
		{"CREATE INDEX other.i ON t (a);", 1235},
		{"CREATE INDEX d.i ON other.u (a);", 1235},
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE other; CREATE TABLE other.u (a INT); CREATE DATABASE d; USE d; CREATE TABLE t (a INT);")
		err := execError(t, c, tc.sql)
		var cerr *Error
		if !errors.As(err, &cerr) || cerr.Code != tc.code {
			t.Errorf("%s: error %v, want %d", tc.sql, err, tc.code)
		}
	}
	c := New()
	mustExec(t, c, "CREATE DATABASE other; CREATE TABLE other.u (a INT); CREATE DATABASE d; USE d; CREATE TABLE t (a INT);")
	mustExec(t, c, "CREATE INDEX d.i ON t (a); CREATE INDEX other.j ON other.u (a);")
	if ddl := c.ShowCreateTable("d", "t"); !strings.Contains(ddl, "KEY `i` (`a`)") {
		t.Errorf("the qualified name:\n%s", ddl)
	}
	if ddl := c.ShowCreateTable("other", "u"); !strings.Contains(ddl, "KEY `j` (`a`)") {
		t.Errorf("the qualified name on another database:\n%s", ddl)
	}
}

// TestStoringIsPartOfTheIndex: which columns an index stores is part of it to the diff, and
// the order they were written in is not.
func TestStoringIsPartOfTheIndex(t *testing.T) {
	load := func(storing string) *Catalog {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT, c INT, KEY k (a)"+storing+");")
		return c
	}
	if d := Diff(load(" STORING (b, c)"), load(" STORING (C, b)")); !d.IsEmpty() {
		t.Errorf("the same columns read as a change: %s", describeDiff(d))
	}
	if d := Diff(load(" STORING (b)"), load(" STORING (c)")); d.IsEmpty() {
		t.Errorf("other columns read as no change")
	}
	if d := Diff(load(""), load(" STORING (b)")); d.IsEmpty() {
		t.Errorf("a new STORING read as no change")
	}
	if d := Diff(load(""), load(" PARALLEL 2")); !d.IsEmpty() {
		t.Errorf("PARALLEL read as a change: %s", describeDiff(d))
	}
}

// TestOracle_OceanBaseIndexOptions holds the text above to the server, and the catalog's
// reading of the DDL to its reading of the server's echo of it.
func TestOracle_OceanBaseIndexOptions(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := oracleCtx(t)
	for _, probe := range []struct{ id, create, want string }{
		{"storing", "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c))", "KEY `k` (`a`) STORING (`b`, `c`)"},
		{"twice", "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b) STORING (c))", "KEY `k` (`a`) STORING (`b`, `c`)"},
		{"case", "CREATE TABLE t (a INT, b INT, KEY k (a) STORING (B))", "KEY `k` (`a`) STORING (`b`)"},
		{"parallel", "CREATE TABLE t (a INT, KEY k (a) PARALLEL 3)", "KEY `k` (`a`)"},
		{"noparallel", "CREATE TABLE t (a INT, KEY k (a) NOPARALLEL)", "KEY `k` (`a`)"},
		{"rowid", "CREATE TABLE t (a INT, KEY k (a) WITH ROWID)", "KEY `k` (`a`)"},
	} {
		t.Run(probe.id, func(t *testing.T) {
			readback, err := storeAndRead(ctx, o, "obopt_"+probe.id, probe.create, "t")
			if err != nil {
				t.Fatalf("OceanBase would not take the form: %v", err)
			}
			if !strings.Contains(readback, probe.want) {
				t.Fatalf("the server printed:\n%s\nwant %s", readback, probe.want)
			}
			for _, not := range []string{"PARALLEL", "ROWID"} {
				if strings.Contains(readback, not) {
					t.Errorf("the server printed %s:\n%s", not, readback)
				}
			}
			sc := serverCharsetFor(o.version)
			n := NormalizerFor(o.version)
			userCat, _ := loadOneTable(t, sc, probe.create, "t")
			storedCat, _ := loadOneTable(t, sc, readback, "t")
			d := DiffWithNormalizer(userCat, storedCat, n)
			for _, te := range d.Tables {
				if len(te.Indexes) != 0 {
					t.Errorf("the index read as a change against its readback: %s", describeDiff(d))
				}
			}
		})
	}
}

// TestStoringMigration: a plan carries the stored columns, however the index comes to be
// written, and its text reads back as the target.
func TestStoringMigration(t *testing.T) {
	load := func(ddl string) *Catalog {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; "+ddl+";")
		return c
	}
	plan := func(from, to *Catalog) string {
		return GenerateMigration(from, to, Diff(from, to)).SQL()
	}
	const (
		bare   = "CREATE TABLE t (a INT, b INT, c INT, KEY k (a))"
		stored = "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c) COMMENT 'x')"
	)
	got := plan(load(bare), load(stored))
	if !strings.Contains(got, "STORING (`b`,`c`)") {
		t.Errorf("a new STORING planned as:\n%s", got)
	}
	replay := load(bare)
	mustExec(t, replay, "USE d; "+got+";")
	if d := Diff(replay, load(stored)); !d.IsEmpty() {
		t.Errorf("the plan's text left a difference: %s", describeDiff(d))
	}

	got = plan(load(bare), load("CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (C))"))
	if !strings.Contains(got, "STORING (`c`)") {
		t.Errorf("a column spelled otherwise planned as:\n%s", got)
	}

	got = plan(load(stored), load(bare))
	if strings.Contains(got, "STORING") {
		t.Errorf("a dropped STORING planned as:\n%s", got)
	}

	empty := New()
	mustExec(t, empty, "CREATE DATABASE d;")
	got = plan(empty, load(stored))
	if !strings.Contains(got, "STORING (`b`,`c`)") {
		t.Errorf("a new table's index planned as:\n%s", got)
	}
}

// TestVectorIndexStoringMigration: a vector index has a renderer of its own, which carries the
// stored columns too.
func TestVectorIndexStoringMigration(t *testing.T) {
	load := func(sql string) *Catalog {
		c, _ := loadOneTable(t, "utf8mb4", sql, "t")
		return c
	}
	bare := load("CREATE TABLE t (a INT, b INT, v VECTOR(3) NOT NULL)")
	stored := load("CREATE TABLE t (a INT, b INT, v VECTOR(3) NOT NULL, VECTOR KEY k (v) WITH (distance=l2, type=hnsw) STORING (b))")
	got := GenerateMigration(bare, stored, Diff(bare, stored)).SQL()
	if !strings.Contains(got, "STORING (`b`)") {
		t.Errorf("a new vector index planned as:\n%s", got)
	}
}

// TestStoringKeepsAnIndexOutOfTheBackingSet: an index that stores columns is no index a foreign
// key made for itself, however it is named, so the diff and the plan both see it.
func TestStoringKeepsAnIndexOutOfTheBackingSet(t *testing.T) {
	load := func(storing string) *Catalog {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE p (id INT PRIMARY KEY); "+
			"CREATE TABLE t (pid INT, extra INT, KEY pid (pid)"+storing+", FOREIGN KEY (pid) REFERENCES p (id));")
		return c
	}
	if d := Diff(load(""), load(" STORING (extra)")); d.IsEmpty() {
		t.Error("a STORING on an index named like a foreign key's read as no change")
	}
	empty := New()
	mustExec(t, empty, "CREATE DATABASE d;")
	target := load(" STORING (extra)")
	if got := GenerateMigration(empty, target, Diff(empty, target)).SQL(); !strings.Contains(got, "STORING (`extra`)") {
		t.Errorf("a new table planned as:\n%s", got)
	}
}

// TestStoringNamesDoNotRun: a quoted name can hold the separator, and is still one column.
func TestStoringNamesDoNotRun(t *testing.T) {
	load := func(storing string) *Catalog {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT, c INT, bc INT, `b,c` INT, KEY k (a)"+storing+"));")
		return c
	}
	for _, other := range []string{" STORING (`b,c`", " STORING (bc"} {
		if d := Diff(load(other), load(" STORING (b, c")); d.IsEmpty() {
			t.Errorf("%s) and STORING (b, c) read as no change", other)
		}
	}
}

// TestOracle_OceanBaseStoringMigration applies a plan on the server: the server takes the text,
// and what it prints reads as the target. USING BTREE and INVISIBLE are left out, since neither
// reads back with or without STORING: the server does not print USING BTREE, and a plan's
// /*!80000 INVISIBLE */ is a comment to a server that reports 5.7.25.
func TestOracle_OceanBaseStoringMigration(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := oracleCtx(t)
	const plainTable = "CREATE TABLE t (a INT, b INT, c INT, KEY k (a))"
	for _, target := range []struct{ id, bare, create string }{
		{"plain", plainTable, "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b))"},
		{"two", plainTable, "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b, c))"},
		{"comment", plainTable, "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b) COMMENT 'x')"},
		{"global", plainTable, "CREATE TABLE t (a INT, b INT, c INT, KEY k (a) STORING (b) GLOBAL)"},
		{"unique", plainTable, "CREATE TABLE t (a INT, b INT, c INT, UNIQUE KEY k (a) STORING (b))"},
		// No vector index: OceanBase 4.3.5.6 answers 1235 ("setting vector index storing column is
		// not supported") for one with STORING, so TestVectorIndexStoringMigration holds that renderer.
	} {
		t.Run(target.id, func(t *testing.T) {
			bare := target.bare
			sc := serverCharsetFor(o.version)
			from, _ := loadOneTable(t, sc, bare, "t")
			to, _ := loadOneTable(t, sc, target.create, "t")
			plan := GenerateMigration(from, to, Diff(from, to)).SQL()
			if !strings.Contains(plan, "STORING") {
				t.Fatalf("the plan lost STORING:\n%s", plan)
			}
			schema := "obstore_" + target.id
			if _, err := storeAndRead(ctx, o, schema, bare, "t"); err != nil {
				t.Fatalf("the server would not take the table: %v", err)
			}
			for _, stmt := range strings.Split(plan, ";\n") {
				stmt = strings.ReplaceAll(strings.TrimSuffix(strings.TrimSpace(stmt), ";"), "`diffdb`.", "`"+schema+"`.")
				if stmt == "" {
					continue
				}
				if _, err := o.db.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("the server would not take the plan:\n%s\n%v", stmt, err)
				}
			}
			var name, readback string
			if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE `"+schema+"`.`t`").Scan(&name, &readback); err != nil {
				t.Fatal(err)
			}
			stored, _ := loadOneTable(t, sc, readback, "t")
			// The server's echo adds table options of its own, so the index is what is held.
			for _, te := range DiffWithNormalizer(to, stored, NormalizerFor(o.version)).Tables {
				if len(te.Indexes) != 0 {
					t.Errorf("the applied plan left an index difference:\n%s", readback)
				}
			}
		})
	}
}

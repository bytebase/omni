package catalog

import (
	"context"
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
	ctx := context.Background()
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

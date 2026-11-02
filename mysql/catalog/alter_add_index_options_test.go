package catalog

import (
	"errors"
	"testing"
)

// alterAddIndexBase and alterAddIndexAlters add one index of every kind by ALTER
// TABLE, each with the options MySQL 8.0 keeps for it. alterAddConstraint used to
// build the index and drop con.IndexOptions for all but the primary key, while the
// CREATE TABLE path has always applied them, so these came back as plain indexes.
const alterAddIndexBase = "CREATE TABLE t (a INT, b INT, c TEXT, v VARCHAR(10), g GEOMETRY NOT NULL)"

var alterAddIndexAlters = []string{
	"ALTER TABLE t ADD INDEX i (a) COMMENT 'hello' INVISIBLE",
	"ALTER TABLE t ADD UNIQUE KEY u (b) COMMENT 'u note' KEY_BLOCK_SIZE=8",
	"ALTER TABLE t ADD UNIQUE INDEX u2 (a, b) INVISIBLE",
	"ALTER TABLE t ADD CONSTRAINT cu UNIQUE (v) COMMENT 'cu note'",
	"ALTER TABLE t ADD KEY k (v) KEY_BLOCK_SIZE=4",
	"ALTER TABLE t ADD FULLTEXT INDEX f (c) COMMENT 'ft note' INVISIBLE",
	"ALTER TABLE t ADD SPATIAL INDEX s (g) COMMENT 'sp note'",
}

// TestAlterAddIndexKeepsOptions pins the catalog state, KEY_BLOCK_SIZE included:
// SHOW CREATE TABLE does not print an index's block size, so the container test
// below cannot see it.
func TestAlterAddIndexKeepsOptions(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, alterAddIndexBase)
	for _, stmt := range alterAddIndexAlters {
		mustExec(t, c, stmt)
	}
	want := map[string]struct {
		comment      string
		visible      bool
		keyBlockSize int
	}{
		"i":  {"hello", false, 0},
		"u":  {"u note", true, 8},
		"u2": {"", false, 0},
		"cu": {"cu note", true, 0},
		"k":  {"", true, 4},
		"f":  {"ft note", false, 0},
		"s":  {"sp note", true, 0},
	}
	tbl := c.GetDatabase("testdb").GetTable("t")
	if len(tbl.Indexes) != len(want) {
		t.Fatalf("got %d indexes, want %d", len(tbl.Indexes), len(want))
	}
	for _, idx := range tbl.Indexes {
		w, ok := want[idx.Name]
		if !ok {
			t.Errorf("unexpected index %q", idx.Name)
			continue
		}
		if idx.Comment != w.comment || idx.Visible != w.visible || idx.KeyBlockSize != w.keyBlockSize {
			t.Errorf("index %s: Comment=%q Visible=%v KeyBlockSize=%d, want %q %v %d",
				idx.Name, idx.Comment, idx.Visible, idx.KeyBlockSize, w.comment, w.visible, w.keyBlockSize)
		}
	}
}

// TestContainer_AlterAddIndexOptions runs the same statements on MySQL 8.0 and on
// the catalog and compares SHOW CREATE TABLE, which prints each index's comment
// and visibility (container test).
func TestContainer_AlterAddIndexOptions(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	for _, stmt := range append([]string{alterAddIndexBase}, alterAddIndexAlters...) {
		if err := ctr.execSQL(stmt); err != nil {
			t.Fatalf("container exec %q: %v", stmt, err)
		}
		mustExec(t, c, stmt)
	}
	ctrDDL, err := ctr.showCreateTable("t")
	if err != nil {
		t.Fatalf("SHOW CREATE TABLE: %v", err)
	}
	if omniDDL := c.ShowCreateTable("test", "t"); normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
		t.Errorf("mismatch:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL, omniDDL)
	}

	// An invisible primary key is refused whole: 3522, and the column stays nullable.
	mustExec(t, c, "CREATE TABLE tn (id INT)")
	if err := ctr.execSQL("CREATE TABLE tn (id INT)"); err != nil {
		t.Fatalf("container exec: %v", err)
	}
	if err := ctr.execSQL(invisiblePrimaryKey); err == nil {
		t.Errorf("MySQL accepted %q", invisiblePrimaryKey)
	}
	if results, _ := c.Exec(invisiblePrimaryKey, nil); len(results) != 1 || results[0].Error == nil {
		t.Errorf("the catalog accepted %q", invisiblePrimaryKey)
	}
	ctrDDL, err = ctr.showCreateTable("tn")
	if err != nil {
		t.Fatalf("SHOW CREATE TABLE tn: %v", err)
	}
	if omniDDL := c.ShowCreateTable("test", "tn"); normalizeWhitespace(ctrDDL) != normalizeWhitespace(omniDDL) {
		t.Errorf("mismatch after the refused key:\n--- container ---\n%s\n--- omni ---\n%s", ctrDDL, omniDDL)
	}
}

const invisiblePrimaryKey = "ALTER TABLE tn ADD PRIMARY KEY (id) INVISIBLE"

// TestAlterAddInvisiblePrimaryKey pins that the refused key leaves the table as it
// was. The 3522 check used to run after the columns were marked NOT NULL.
func TestAlterAddInvisiblePrimaryKey(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE tn (id INT)")
	results, err := c.Exec(invisiblePrimaryKey, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	var catErr *Error
	if len(results) != 1 || !errors.As(results[0].Error, &catErr) || catErr.Code != 3522 {
		t.Fatalf("got %+v, want error 3522", results)
	}
	tbl := c.GetDatabase("testdb").GetTable("tn")
	if len(tbl.Indexes) != 0 || !tbl.GetColumn("id").Nullable {
		t.Errorf("the refused key changed the table:\n%s", c.ShowCreateTable("testdb", "tn"))
	}
}

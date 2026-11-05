package catalog

import (
	"context"
	"testing"
)

// alterAddIndexBase and alterAddIndexAlters add one index of every kind by ALTER
// TABLE, each with the options MariaDB keeps for it. alterAddConstraint used to
// build the index and drop con.IndexOptions for all but the primary key, while the
// CREATE TABLE path has always applied them, so these came back as plain indexes.
// MariaDB hides an index with IGNORED and answers 1064 for INVISIBLE, so visibility
// is not among the options here.
const alterAddIndexBase = "CREATE TABLE t (a INT, b INT, c TEXT, v VARCHAR(10), g GEOMETRY NOT NULL)"

var alterAddIndexAlters = []string{
	"ALTER TABLE t ADD INDEX i (a) COMMENT 'hello'",
	"ALTER TABLE t ADD UNIQUE KEY u (b) COMMENT 'u note' KEY_BLOCK_SIZE=8",
	"ALTER TABLE t ADD CONSTRAINT cu UNIQUE (v) COMMENT 'cu note'",
	"ALTER TABLE t ADD KEY k (v) KEY_BLOCK_SIZE=4",
	"ALTER TABLE t ADD FULLTEXT INDEX f (c) COMMENT 'ft note'",
	"ALTER TABLE t ADD SPATIAL INDEX s (g) COMMENT 'sp note'",
}

type indexOptions struct {
	comment      string
	keyBlockSize int
}

func tableIndexOptions(t *testing.T, c *Catalog, db string) map[string]indexOptions {
	t.Helper()
	got := map[string]indexOptions{}
	for _, idx := range c.GetDatabase(db).GetTable("t").Indexes {
		got[idx.Name] = indexOptions{idx.Comment, idx.KeyBlockSize}
	}
	return got
}

// TestAlterAddIndexKeepsOptions pins the catalog state the ALTERs leave.
func TestAlterAddIndexKeepsOptions(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, alterAddIndexBase)
	for _, stmt := range alterAddIndexAlters {
		mustExec(t, c, stmt)
	}
	want := map[string]indexOptions{
		"i":  {"hello", 0},
		"u":  {"u note", 8},
		"cu": {"cu note", 0},
		"k":  {"", 4},
		"f":  {"ft note", 0},
		"s":  {"sp note", 0},
	}
	got := tableIndexOptions(t, c, "testdb")
	if len(got) != len(want) {
		t.Fatalf("got indexes %v, want %v", got, want)
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("index %s: got %+v, want %+v", name, got[name], w)
		}
	}
}

// TestContainerAlterAddIndexOptions runs the statements on MariaDB 11.8.8, loads
// the table MariaDB reads back with SHOW CREATE TABLE into a fresh catalog, and
// requires every index's comment and block size to match what the ALTERs left
// in this one (container test). MariaDB prints an index's KEY_BLOCK_SIZE, which
// MySQL does not.
func TestContainerAlterAddIndexOptions(t *testing.T) {
	db := mariadbDatabase(t)
	ctx := context.Background()

	c := setupWithDB(t)
	for _, stmt := range append([]string{alterAddIndexBase}, alterAddIndexAlters...) {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("MariaDB exec %q: %v", stmt, err)
		}
		mustExec(t, c, stmt)
	}
	var name, readback string
	if err := db.QueryRowContext(ctx, "SHOW CREATE TABLE t").Scan(&name, &readback); err != nil {
		t.Fatalf("SHOW CREATE TABLE: %v", err)
	}
	engine := setupWithDB(t)
	mustExec(t, engine, readback)

	want := tableIndexOptions(t, engine, "testdb")
	got := tableIndexOptions(t, c, "testdb")
	if len(got) != len(want) {
		t.Fatalf("catalog indexes %v, MariaDB reads back %v:\n%s", got, want, readback)
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("index %s: catalog %+v, MariaDB reads back %+v", name, got[name], w)
		}
	}
}

package catalog

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// alterAddIndexBase and alterAddIndexAlters add a primary key and secondary
// indexes by ALTER TABLE, each with the options TiDB keeps for it.
// alterAddConstraint used to build every index and never read con.IndexOptions,
// while the CREATE TABLE path has always applied them, so these came back as
// plain indexes. TiDB v8.5.5 has no SPATIAL index and accepts a FULLTEXT one
// without keeping it, so neither is here.
const alterAddIndexBase = "CREATE TABLE t (id INT NOT NULL, a INT, b INT, v VARCHAR(10))"

var alterAddIndexAlters = []string{
	"ALTER TABLE t ADD PRIMARY KEY (id) COMMENT 'pk note'",
	"ALTER TABLE t ADD INDEX i (a) COMMENT 'hello' INVISIBLE",
	"ALTER TABLE t ADD UNIQUE KEY u (b) COMMENT 'u note' KEY_BLOCK_SIZE=8",
	"ALTER TABLE t ADD UNIQUE INDEX u2 (a, b) INVISIBLE",
	"ALTER TABLE t ADD CONSTRAINT cu UNIQUE (v) COMMENT 'cu note'",
	"ALTER TABLE t ADD KEY k (v) KEY_BLOCK_SIZE=4",
}

type indexOptions struct {
	comment      string
	visible      bool
	keyBlockSize int
}

func tableIndexOptions(c *Catalog, db string) map[string]indexOptions {
	got := map[string]indexOptions{}
	for _, idx := range c.GetDatabase(db).GetTable("t").Indexes {
		got[idx.Name] = indexOptions{idx.Comment, idx.Visible, idx.KeyBlockSize}
	}
	return got
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
	want := map[string]indexOptions{
		"PRIMARY": {"pk note", true, 0},
		"i":       {"hello", false, 0},
		"u":       {"u note", true, 8},
		"u2":      {"", false, 0},
		"cu":      {"cu note", true, 0},
		"k":       {"", true, 4},
	}
	got := tableIndexOptions(c, "testdb")
	if len(got) != len(want) {
		t.Fatalf("got indexes %v, want %v", got, want)
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("index %s: got %+v, want %+v", name, got[name], w)
		}
	}
}

// TestAlterAddInvisiblePrimaryKey pins TiDB's refusal of an invisible primary key.
// Reading the options on this path is what lets the catalog see INVISIBLE at all,
// and adding the key anyway would leave a table TiDB never has. TiDB refuses the
// statement whole, so the column it names stays nullable.
func TestAlterAddInvisiblePrimaryKey(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE tn (id INT)")
	results, err := c.Exec("ALTER TABLE tn ADD PRIMARY KEY (id) INVISIBLE", nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	var catErr *Error
	if len(results) != 1 || !errors.As(results[0].Error, &catErr) || catErr.Code != 3522 {
		t.Fatalf("got %+v, want error 3522", results)
	}
	tbl := c.GetDatabase("testdb").GetTable("tn")
	if n := len(tbl.Indexes); n != 0 {
		t.Errorf("the refused key was added: %d indexes", n)
	}
	if !tbl.GetColumn("id").Nullable {
		t.Error("the refused key made id NOT NULL")
	}
}

// TestContainerAlterAddIndexOptions runs the statements on TiDB v8.5.5, loads the
// table TiDB reads back with SHOW CREATE TABLE into a fresh catalog, and requires
// every index's comment and visibility to match what the ALTERs left in this one.
// It also asks TiDB about the invisible primary key (container test).
func TestContainerAlterAddIndexOptions(t *testing.T) {
	tc := startTiDBForCatalog(t)
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS omni_alter_index",
		"CREATE DATABASE omni_alter_index",
		"USE omni_alter_index",
	} {
		mustExecTiDB(t, tc, stmt)
	}
	defer tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_alter_index")

	c := setupWithDB(t)
	for _, stmt := range append([]string{alterAddIndexBase}, alterAddIndexAlters...) {
		mustExecTiDB(t, tc, stmt)
		mustExec(t, c, stmt)
	}
	var name, readback string
	if err := tc.db.QueryRowContext(tc.ctx, "SHOW CREATE TABLE t").Scan(&name, &readback); err != nil {
		t.Fatalf("SHOW CREATE TABLE: %v", err)
	}
	engine := setupWithDB(t)
	mustExec(t, engine, readback)

	want := tableIndexOptions(engine, "testdb")
	got := tableIndexOptions(c, "testdb")
	if len(got) != len(want) {
		t.Fatalf("catalog indexes %v, TiDB reads back %v:\n%s", got, want, readback)
	}
	for name, w := range want {
		// TiDB does not print an index's KEY_BLOCK_SIZE, so the readback cannot
		// answer for it; TestAlterAddIndexKeepsOptions does.
		g := got[name]
		g.keyBlockSize, w.keyBlockSize = 0, 0
		if g != w {
			t.Errorf("index %s: catalog %+v, TiDB reads back %+v", name, g, w)
		}
	}

	mustExecTiDB(t, tc, "CREATE TABLE tn (id INT)")
	_, err := tc.db.ExecContext(tc.ctx, "ALTER TABLE tn ADD PRIMARY KEY (id) INVISIBLE")
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) || myErr.Number != 3522 {
		t.Errorf("TiDB answered %v for an invisible primary key, want 3522", err)
	}
	if err := tc.db.QueryRowContext(tc.ctx, "SHOW CREATE TABLE tn").Scan(&name, &readback); err != nil {
		t.Fatalf("SHOW CREATE TABLE tn: %v", err)
	}
	if !strings.Contains(readback, "`id` int DEFAULT NULL") {
		t.Errorf("TiDB changed id when it refused the key:\n%s", readback)
	}
}

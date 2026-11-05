package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// These hold what the VECTOR type and a vector index leave in the catalog, not that they
// parse. v8.5.5 runs the column forms, so their expected text is what its SHOW CREATE
// TABLE printed (read on 2026-11-05). It refuses a vector index with 8200 on a server
// with no TiFlash, which the test images have none of, so the index forms are held to the
// grammar's restore text in parser_test.go instead.

func newVectorCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b VECTOR(3));")
	return c
}

func TestVectorColumnShowCreate(t *testing.T) {
	c := New()
	execAll(t, c, `CREATE DATABASE d; USE d;
CREATE TABLE v1 (id INT PRIMARY KEY, a VECTOR(3), b VECTOR<FLOAT>(3), c VECTOR<FLOAT>, e VECTOR NOT NULL);
ALTER TABLE v1 ADD COLUMN f VECTOR(4), MODIFY e VECTOR(2), ADD COLUMN g VECTOR<FLOAT>(2) DEFAULT NULL COMMENT 'cc';`)
	ddl := c.ShowCreateTable("d", "v1")
	for _, want := range []string{
		"`a` vector(3) DEFAULT NULL",
		"`b` vector(3) DEFAULT NULL",
		"`c` vector DEFAULT NULL",
		"`e` vector(2) DEFAULT NULL",
		"`f` vector(4) DEFAULT NULL",
		"`g` vector(2) DEFAULT NULL COMMENT 'cc'",
	} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}
}

// TestCreateVectorIndexIsAVectorIndex: the statement adds an index the catalog knows for
// one, with the type the grammar's action writes when the statement names none, and it
// can be dropped again.
func TestCreateVectorIndexIsAVectorIndex(t *testing.T) {
	c := newVectorCatalog(t)
	execAll(t, c, "CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(b))) COMMENT 'c';")

	tbl := c.databases["d"].Tables["t"]
	if len(tbl.Indexes) != 1 {
		t.Fatalf("Indexes=%+v", tbl.Indexes)
	}
	idx := tbl.Indexes[0]
	if !idx.Vector || idx.Name != "i" || idx.IndexType != "HNSW" || idx.Comment != "c" || idx.Unique {
		t.Errorf("index: %+v", idx)
	}
	ddl := c.ShowCreateTable("d", "t")
	if !strings.Contains(ddl, "VECTOR INDEX `i` (") || !strings.Contains(ddl, "USING HNSW") {
		t.Errorf("SHOW CREATE TABLE lost the vector index:\n%s", ddl)
	}

	if err := execExpectError(t, c, "CREATE VECTOR INDEX i ON t ((VEC_L2_DISTANCE(b)));"); err == nil {
		t.Errorf("a second vector index named i was accepted")
	}
	execAll(t, c, "CREATE VECTOR INDEX IF NOT EXISTS i ON t ((VEC_L2_DISTANCE(b)));")
	if len(tbl.Indexes) != 1 {
		t.Errorf("IF NOT EXISTS added an index: %+v", tbl.Indexes)
	}
	execAll(t, c, "DROP INDEX i ON t;")
	if len(tbl.Indexes) != 0 {
		t.Errorf("DROP INDEX left %+v", tbl.Indexes)
	}
}

// TestVectorIndexUnnamed: an unnamed key over an expression takes the name v8.5.5 gives
// any: `KEY ((a+1)), KEY ((a+2))` reads back expression_index and expression_index_2.
func TestVectorIndexUnnamed(t *testing.T) {
	c := New()
	execAll(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))));")
	execAll(t, c, "ALTER TABLE t ADD VECTOR INDEX ((VEC_L2_DISTANCE(b)));")
	tbl := c.databases["d"].Tables["t"]
	if len(tbl.Indexes) != 2 || tbl.Indexes[0].Name != "expression_index" || tbl.Indexes[1].Name != "expression_index_2" {
		t.Errorf("Indexes=%+v", tbl.Indexes)
	}
	for _, idx := range tbl.Indexes {
		if !idx.Vector || idx.IndexType != "HNSW" {
			t.Errorf("index: %+v", idx)
		}
	}
}

// TestVectorIndexLike: CREATE TABLE LIKE copies a vector index as one.
func TestVectorIndexLike(t *testing.T) {
	c := newVectorCatalog(t)
	execAll(t, c, "CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(b))); CREATE TABLE u LIKE t;")
	if ddl := c.ShowCreateTable("d", "u"); !strings.Contains(ddl, "VECTOR INDEX `i` (") {
		t.Errorf("LIKE lost the vector index:\n%s", ddl)
	}
}

// TestVectorIndexSyntaxRefused: what the grammar's actions answer 1149 for reaches the
// catalog as parsed text, and leaves nothing behind.
func TestVectorIndexSyntaxRefused(t *testing.T) {
	for _, sql := range []string{
		"CREATE VECTOR INDEX i ON t (b);",
		"CREATE VECTOR INDEX i ON t ((b), a);",
		"CREATE VECTOR INDEX i ON t ((b)) USING BTREE;",
		"CREATE INDEX i ON t (a) USING HNSW;",
		"ALTER TABLE t ADD VECTOR INDEX i (b);",
		"ALTER TABLE t ADD VECTOR INDEX i ((b)) USING HASH;",
		"ALTER TABLE t ADD VECTOR INDEX i (a, (b));",
		"CREATE TABLE u (b VECTOR(3), VECTOR INDEX i (b));",
	} {
		c := newVectorCatalog(t)
		err := execExpectError(t, c, sql)
		var cerr *Error
		if !errors.As(err, &cerr) || cerr.Code != ErrSyntaxError {
			t.Errorf("%s: error %v, want 1149", sql, err)
		}
		if n := len(c.databases["d"].Tables["t"].Indexes); n != 0 {
			t.Errorf("%s: left %d indexes", sql, n)
		}
	}
	// The server answers 1149 first, and 1146 only for the shape it reads.
	c := newVectorCatalog(t)
	var cerr *Error
	if err := execExpectError(t, c, "CREATE VECTOR INDEX i ON nosuch (b);"); !errors.As(err, &cerr) || cerr.Code != ErrSyntaxError {
		t.Errorf("a refused shape on a missing table: %v, want 1149", err)
	}
	if err := execExpectError(t, c, "CREATE VECTOR INDEX i ON nosuch ((b));"); !errors.As(err, &cerr) || cerr.Code != ErrNoSuchTable {
		t.Errorf("a good shape on a missing table: %v, want 1146", err)
	}
}

// TestVectorIndexFromMetadata: a synced vector index stays one. It syncs as an index of
// type HNSW, and read as an ordinary index SHOW CREATE TABLE would call it KEY.
func TestVectorIndexFromMetadata(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}, {Name: "b", Type: "vector(3)", Nullable: true}},
			Indexes: []*metadata.IndexMetadata{
				{Name: "vi", Type: "HNSW", Visible: true, Expressions: []string{"VEC_COSINE_DISTANCE(`b`)"}},
				{Name: "ki", Type: "BTREE", Visible: true, Expressions: []string{"a"}},
			},
		}},
	}))
	requireReport(t, report)
	tbl := requireTable(t, c, "t")
	var vi, ki *Index
	for _, idx := range tbl.Indexes {
		switch idx.Name {
		case "vi":
			vi = idx
		case "ki":
			ki = idx
		}
	}
	if vi == nil || !vi.Vector || vi.IndexType != "HNSW" {
		t.Errorf("vi: %+v", vi)
	}
	if ki == nil || ki.Vector {
		t.Errorf("ki: %+v", ki)
	}
}

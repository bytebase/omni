package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	mysqlparser "github.com/bytebase/omni/tidb/parser"
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

// TestVectorIndexShapeRefused: what the grammar's actions refuse reaches the catalog as a
// parse error, and leaves nothing behind.
func TestVectorIndexShapeRefused(t *testing.T) {
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
		var perr *mysqlparser.ParseError
		if _, err := c.Exec(sql, nil); !errors.As(err, &perr) {
			t.Errorf("%s: error %v, want a parse error", sql, err)
		}
		if n := len(c.databases["d"].Tables["t"].Indexes); n != 0 {
			t.Errorf("%s: left %d indexes", sql, n)
		}
	}
	// The server refuses that shape first, and answers 1146 only for the shape it reads.
	c := newVectorCatalog(t)
	var perr *mysqlparser.ParseError
	if _, err := c.Exec("CREATE VECTOR INDEX i ON nosuch (b);", nil); !errors.As(err, &perr) {
		t.Errorf("a refused shape on a missing table: %v, want a parse error", err)
	}
	var cerr *Error
	if err := execExpectError(t, c, "CREATE VECTOR INDEX i ON nosuch ((b));"); !errors.As(err, &cerr) || cerr.Code != ErrNoSuchTable {
		t.Errorf("a good shape on a missing table: %v, want 1146", err)
	}
}

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

// TestVectorCastViewText: a view's stored text drops a CAST's dimension, because TiDB's does.
// v8.5.5's SHOW CREATE VIEW writes CAST(... AS VECTOR) for VECTOR(3), for VECTOR<FLOAT>(3) and
// for a bare VECTOR (read on 2026-11-05).
func TestVectorCastViewText(t *testing.T) {
	for _, sql := range []string{
		"CREATE VIEW vv AS SELECT CAST('[1,2,3]' AS VECTOR(3)) AS v;",
		"CREATE VIEW vv AS SELECT CAST('[1,2,3]' AS VECTOR<FLOAT>(3)) AS v;",
		"CREATE VIEW vv AS SELECT CONVERT('[1,2,3]', VECTOR<FLOAT>(3)) AS v;",
		"CREATE VIEW vv AS SELECT CAST('[1,2,3]' AS VECTOR) AS v;",
	} {
		c := newVectorCatalog(t)
		execAll(t, c, sql)
		def := strings.ToLower(c.GetDatabase("d").Views["vv"].Definition)
		if !strings.Contains(def, " vector)") && !strings.Contains(def, ", vector)") {
			t.Errorf("%s: stored %s, want the type with no dimension", sql, def)
		}
		if strings.Contains(def, "vector(") || strings.Contains(def, "vector<") {
			t.Errorf("%s: stored %s, which keeps what v8.5.5 drops", sql, def)
		}
	}
}

// TestVectorIndexTypeReadsBack: a definition the parser takes has the one type a vector index has,
// whichever position named which type, so what the catalog writes reads back.
func TestVectorIndexTypeReadsBack(t *testing.T) {
	const part = "((VEC_COSINE_DISTANCE(b)))"
	for name, stmt := range map[string]string{
		"constraint, the leading type wins": "CREATE TABLE u (b VECTOR(3), VECTOR INDEX i USING HNSW " + part + " USING BTREE)",
		"constraint, no leading type":       "CREATE TABLE u (b VECTOR(3), VECTOR INDEX i " + part + " USING HNSW)",
		"constraint, no type":               "CREATE TABLE u (b VECTOR(3), VECTOR INDEX i " + part + ")",
		"ALTER, the leading type wins":      "CREATE TABLE u (b VECTOR(3)); ALTER TABLE u ADD VECTOR INDEX i USING HNSW " + part + " USING BTREE",
		"CREATE INDEX, the option wins":     "CREATE TABLE u (b VECTOR(3)); CREATE VECTOR INDEX i USING BTREE ON u " + part + " USING HNSW",
	} {
		t.Run(name, func(t *testing.T) {
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+stmt, nil); err != nil {
				t.Fatal(err)
			}
			text := c.ShowCreateTable("d", "u")
			if !strings.Contains(text, "USING HNSW") || strings.Contains(text, "BTREE") {
				t.Fatalf("the index is no HNSW one:\n%s", text)
			}
			again := New()
			if _, err := again.Exec("CREATE DATABASE d; USE d; "+text, nil); err != nil {
				t.Fatalf("%s does not read back: %v", text, err)
			}
			if got := again.ShowCreateTable("d", "u"); got != text {
				t.Errorf("read back as:\n%s\nwant:\n%s", got, text)
			}
		})
	}
}

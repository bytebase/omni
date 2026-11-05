package catalog

import (
	"context"
	"strings"
	"testing"
)

// These hold what the VECTOR type and a vector index leave in the catalog, not that they
// parse. Their expected text is what SHOW CREATE TABLE printed on 4.3.5.6 (read on
// 2026-11-05); TestOracle_OceanBaseVectorIndex holds it to the server.

func TestVectorShowCreate(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))",
			"VECTOR KEY `i` (`v`) WITH (DISTANCE=COSINE, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64)"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (type=hnsw, lib=vsag, m=8, distance=l2))",
			"VECTOR KEY `i` (`v`) WITH (TYPE=HNSW, LIB=VSAG, M=8, DISTANCE=L2, EF_CONSTRUCTION=200, EF_SEARCH=64)"},
		{"CREATE TABLE t (v VECTOR(3) NOT NULL); CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'",
			"VECTOR KEY `i` (`v`) WITH (DISTANCE=L2, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64)"},
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; "+tc.sql+";")
		ddl := c.ShowCreateTable("d", "t")
		if !strings.Contains(ddl, tc.want) {
			t.Errorf("%s: SHOW CREATE TABLE lacks %s:\n%s", tc.sql, tc.want, ddl)
		}
		if !strings.Contains(ddl, "`v` VECTOR(3) NOT NULL") {
			t.Errorf("%s: SHOW CREATE TABLE lacks the column in upper case:\n%s", tc.sql, ddl)
		}
		idx := c.databases["d"].Tables["t"].Indexes[0]
		if !idx.Vector || idx.Name != "i" {
			t.Errorf("%s: index %+v", tc.sql, idx)
		}
	}
}

// TestVectorIndexLike: CREATE TABLE LIKE copies a vector index as one.
func TestVectorIndexLike(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw)); CREATE TABLE u LIKE t;")
	ddl := c.ShowCreateTable("d", "u")
	if !strings.Contains(ddl, "VECTOR KEY `i` (`v`) WITH (DISTANCE=L2, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64)") {
		t.Errorf("LIKE lost the vector index:\n%s", ddl)
	}
}

func TestVectorIndexDrop(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (v VECTOR(3) NOT NULL); CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw);")
	if n := len(c.databases["d"].Tables["t"].Indexes); n != 1 {
		t.Fatalf("%d indexes", n)
	}
	mustExec(t, c, "DROP INDEX i ON t;")
	if n := len(c.databases["d"].Tables["t"].Indexes); n != 0 {
		t.Errorf("DROP INDEX left %d indexes", n)
	}
}

// TestVectorIndexWithIsPartOfTheIndex: two indexes that differ only in their parameters
// are two different indexes to the diff.
func TestVectorIndexWithIsPartOfTheIndex(t *testing.T) {
	load := func(with string) *Catalog {
		c := New()
		mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH ("+with+"));")
		return c
	}
	if d := Diff(load("distance=l2, type=hnsw"), load("type=hnsw, distance=l2, lib=vsag, m=16")); !d.IsEmpty() {
		t.Errorf("the same parameters read as a change: %s", describeDiff(d))
	}
	if d := Diff(load("distance=l2, type=hnsw"), load("distance=cosine, type=hnsw")); d.IsEmpty() {
		t.Errorf("a changed distance read as no change")
	}
	if d := Diff(load("distance=l2, type=hnsw"), load("distance=l2, type=hnsw, m=8")); d.IsEmpty() {
		t.Errorf("a changed m read as no change")
	}
}

// TestOracle_OceanBaseVectorIndex holds the text above to the server, and the catalog's
// reading of the DDL to its reading of the server's echo of it, so that a vector index
// does not read as a change against its own readback.
func TestOracle_OceanBaseVectorIndex(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, probe := range []struct{ id, create, want string }{
		{"defaults", "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))",
			"VECTOR KEY `i` (`v`) WITH (DISTANCE=COSINE, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64)"},
		{"reordered", "CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (type=hnsw, lib=vsag, m=8, distance=l2))",
			"VECTOR KEY `i` (`v`) WITH (TYPE=HNSW, LIB=VSAG, M=8, DISTANCE=L2, EF_CONSTRUCTION=200, EF_SEARCH=64)"},
	} {
		t.Run(probe.id, func(t *testing.T) {
			readback, err := storeAndRead(ctx, o, "obvec_"+probe.id, probe.create, "t")
			if err != nil {
				t.Fatalf("OceanBase would not take the form: %v", err)
			}
			if !strings.Contains(readback, probe.want) || !strings.Contains(readback, "`v` VECTOR(3) NOT NULL") {
				t.Fatalf("the server printed:\n%s\nwant %s", readback, probe.want)
			}
			sc := serverCharsetFor(o.version)
			n := NormalizerFor(o.version)
			userCat, _ := loadOneTable(t, sc, probe.create, "t")
			storedCat, _ := loadOneTable(t, sc, readback, "t")
			d := DiffWithNormalizer(userCat, storedCat, n)
			for _, te := range d.Tables {
				if len(te.Indexes) != 0 {
					t.Errorf("the vector index read as a change against its readback: %s", describeDiff(d))
				}
			}
		})
	}
}

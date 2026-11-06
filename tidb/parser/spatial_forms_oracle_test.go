package parser

import "testing"

// What TiDB v8.5.5 reads of MySQL's spatial syntax: almost nothing. Its grammar has
// no spatial column type, no SRID column attribute, no SPATIAL REFERENCE SYSTEM
// statement, and no SPATIAL index inside a table definition or an ALTER TABLE ADD.
// The one place SPATIAL appears is IndexKeyTypeOpt, so `CREATE SPATIAL INDEX` is
// read (and refused later, by the DDL). The type names stay unreserved keywords, so
// they still read as names. Every verdict was read from pingcap/tidb:v8.5.5 on
// 2026-11-05.
//
// Ref: IndexKeyTypeOpt and CreateIndexStmt in TiDB v8.5.5 pkg/parser/parser.y, which
// has no GEOMETRY, POINT, LINESTRING, POLYGON, MULTI*, GEOMETRYCOLLECTION or SRID
// token at all.

type spatialFormCase struct {
	sql    string
	accept bool
}

var spatialFormCases = []spatialFormCase{
	// column types.
	{"CREATE TABLE t (a GEOMETRY)", false},
	{"CREATE TABLE t (a POINT)", false},
	{"CREATE TABLE t (a LINESTRING)", false},
	{"CREATE TABLE t (a POLYGON)", false},
	{"CREATE TABLE t (a MULTIPOINT)", false},
	{"CREATE TABLE t (a MULTILINESTRING)", false},
	{"CREATE TABLE t (a MULTIPOLYGON)", false},
	{"CREATE TABLE t (a GEOMETRYCOLLECTION)", false},
	{"CREATE TABLE t (a GEOMCOLLECTION)", false},
	{"CREATE TABLE t (a POINT NOT NULL SRID 4326)", false},
	{"CREATE TABLE t (a INT SRID 4326)", false},
	{"ALTER TABLE t ADD COLUMN a GEOMETRY", false},
	{"ALTER TABLE t MODIFY a POINT", false},
	{"SELECT CAST(a AS GEOMETRY) FROM t", false},

	// the type names as names.
	{"CREATE TABLE t (geometry INT, point INT, linestring INT, polygon INT)", true},
	{"CREATE TABLE geometry (multipoint INT, multilinestring INT, multipolygon INT)", true},
	{"CREATE TABLE t (geometrycollection INT, geomcollection INT, srid INT)", true},
	{"SELECT point, polygon FROM geometry", true},

	// SPATIAL indexes.
	{"CREATE SPATIAL INDEX i ON t (a)", true},
	{"CREATE SPATIAL INDEX IF NOT EXISTS i ON t (a) USING BTREE", true},
	{"CREATE SPATIAL KEY i ON t (a)", false},
	{"CREATE SPATIAL i ON t (a)", false},
	{"CREATE TEMPORARY SPATIAL INDEX i ON t (a)", false},
	// CreateIndexStmt spells the keyword INDEX after every IndexKeyTypeOpt, and
	// has no TEMPORARY in front of any of them.
	{"CREATE FULLTEXT INDEX i ON t (a)", true},
	{"CREATE FULLTEXT KEY i ON t (a)", false},
	{"CREATE TEMPORARY FULLTEXT INDEX i ON t (a)", false},
	{"CREATE TEMPORARY UNIQUE INDEX i ON t (a)", false},
	{"CREATE TEMPORARY INDEX i ON t (a)", false},
	{"CREATE TABLE t (a INT, SPATIAL INDEX i (a))", false},
	{"CREATE TABLE t (a INT, SPATIAL KEY i (a))", false},
	{"CREATE TABLE t (a INT, SPATIAL (a))", false},
	{"ALTER TABLE t ADD SPATIAL INDEX i (a)", false},
	{"ALTER TABLE t ADD SPATIAL KEY i (a)", false},
	{"ALTER TABLE t ADD SPATIAL (a)", false},

	// spatial reference systems.
	{"CREATE SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'", false},
	{"CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'x' DEFINITION 'y'", false},
	{"DROP SPATIAL REFERENCE SYSTEM 4120", false},
	{"DROP SPATIAL REFERENCE SYSTEM IF EXISTS 4120", false},

	// spatial functions are ordinary calls.
	{"SELECT ST_GeomFromText('POINT(1 1)')", true},
	{"SELECT ST_Distance(a, b) FROM t", true},
}

func TestSpatialFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range spatialFormCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = spatialFormCases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, c.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, c.accept, c.sql, perr)
			}
		})
	}
}

// TestSpatialForms holds the parser to the recorded verdicts without the server, so a
// regression fails however the suite is run.
func TestSpatialForms(t *testing.T) {
	for _, c := range spatialFormCases {
		t.Run(c.sql, func(t *testing.T) {
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}

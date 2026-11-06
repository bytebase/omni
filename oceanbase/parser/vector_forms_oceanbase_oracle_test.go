package parser

import (
	"context"
	"testing"
)

// The VECTOR column type, the vector index and the WITH (...) list, as OceanBase 4.3.5.6
// reads them. MySQL 8.0 has none of this and TiDB and MariaDB read other forms of the
// same words, so it is inferred from neither: every verdict was read from
// oceanbase/oceanbase-ce:4.3.5-lts on 2026-11-05.
//
// Left out because they are no part of this family: `SELECT vector(3)`, refused like
// `SELECT comment(3)`; and `ORDER BY a <=> b APPROXIMATE`, the vector search forms.
// index_options_oceanbase_oracle_test.go holds what a vector index takes from every index.

type vectorFormCase struct {
	sql    string
	accept bool
}

var vectorFormCases = []vectorFormCase{
	// column type.
	{"CREATE TABLE t (a VECTOR)", false},
	{"CREATE TABLE t (a VECTOR(3))", true},
	{"CREATE TABLE t (a VECTOR(0))", true},
	{"CREATE TABLE t (a VECTOR(16383))", true},
	{"CREATE TABLE t (a VECTOR(16384))", true},
	{"CREATE TABLE t (a VECTOR(99999999999999999999))", false},
	{"CREATE TABLE t (a VECTOR<FLOAT>)", false},
	{"CREATE TABLE t (a VECTOR<FLOAT>(3))", false},
	{"CREATE TABLE t (a VECTOR(5)<FLOAT>)", false},
	{"CREATE TABLE t (a VECTOR())", false},
	{"CREATE TABLE t (a VECTOR(3,2))", false},
	{"CREATE TABLE t (a VECTOR('3'))", false},
	{"CREATE TABLE t (a VECTOR(b))", false},
	{"CREATE TABLE t (a VECTOR(1.5))", false},
	{"CREATE TABLE t (a VECTOR(0x3))", false},
	{"CREATE TABLE t (a VECTOR(-1))", false},
	{"CREATE TABLE t (a VECTOR(3)", false},
	{"CREATE TABLE t (a VECTOR(3, FLOAT))", false},
	{"CREATE TABLE t (a VECTOR(3) NOT NULL)", true},
	{"CREATE TABLE t (a VECTOR(3) NULL)", true},
	{"CREATE TABLE t (a VECTOR(3) COMMENT 'x')", true},
	{"CREATE TABLE t (a VECTOR(3) DEFAULT NULL)", true},
	{"CREATE TABLE t (a VECTOR(3) DEFAULT '[1,2,3]')", true},
	{"CREATE TABLE t (a VECTOR(3) UNSIGNED)", false},
	{"CREATE TABLE t (a VECTOR(3) ZEROFILL)", false},
	{"CREATE TABLE t (a VECTOR(3) BINARY)", false},
	{"CREATE TABLE t (a VECTOR(3) CHARSET utf8mb4)", false},
	{"CREATE TABLE t (a VECTOR(3) CHARACTER SET utf8mb4)", false},
	{"CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin)", true},
	{"CREATE TABLE t (a VECTOR(3) PRIMARY KEY)", true},
	{"CREATE TABLE t (a VECTOR(3) ARRAY)", false},
	{"CREATE TABLE t (a VECTOR(3) GENERATED ALWAYS AS ('[1]'))", true},
	{"CREATE TABLE t (a VECTOR(3) AS ('[1]') STORED)", true},
	{"CREATE TABLE t (a `VECTOR`)", false},
	{"CREATE TABLE t (a `VECTOR`(3))", false},
	{"CREATE TABLE t (a vector(3))", true},
	{"CREATE TEMPORARY TABLE tt (a VECTOR(3))", true},
	// table element.
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY (v))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR (v))", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR i (v))", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX IF NOT EXISTS i (v))", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i USING HNSW (v))", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING HNSW)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING BTREE)", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) TYPE HNSW)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX ((VEC_COSINE_DISTANCE(v))))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v, w))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v(3)))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v DESC))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M 6)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=abc)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M='6')", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=euclidean)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=cosine)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=foo)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE='cosine')", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=default)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=default)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c')", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) VISIBLE)", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) IGNORED)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) NOT IGNORED)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) KEY_BLOCK_SIZE=1)", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH PARSER p)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) FOO=1)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag, m=16, ef_construction=200, ef_search=64))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=inner_product, type=hnsw_sq, lib=vsag))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=ivf_flat))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH ())", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=foo, type=bar))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (foo=bar))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2 type=hnsw))", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH distance=l2)", false},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c')", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v) WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE t (v VECTOR(3) NOT NULL, CONSTRAINT c VECTOR INDEX (v))", false},
	// CREATE VECTOR INDEX.
	{"CREATE VECTOR INDEX i ON t (v)", true},
	{"CREATE VECTOR KEY i ON t (v)", false},
	{"CREATE VECTOR i ON t (v)", false},
	{"CREATE VECTOR INDEX ON t (v)", false},
	{"CREATE VECTOR INDEX", false},
	{"CREATE VECTOR INDEX i ON t", false},
	{"CREATE VECTOR INDEX i ON t ()", false},
	{"CREATE VECTOR", false},
	{"CREATE OR REPLACE VECTOR INDEX i ON t (v)", false},
	{"CREATE VECTOR INDEX IF NOT EXISTS i ON t (v)", true},
	{"CREATE OR REPLACE VECTOR INDEX IF NOT EXISTS i ON t (v)", false},
	{"CREATE UNIQUE VECTOR INDEX i ON t (v)", false},
	{"CREATE VECTOR UNIQUE INDEX i ON t (v)", false},
	{"CREATE TEMPORARY VECTOR INDEX i ON t (v)", false},
	{"CREATE VECTOR INDEX i USING HNSW ON t (v)", false},
	{"CREATE VECTOR INDEX i TYPE HNSW ON t (v)", false},
	{"CREATE VECTOR INDEX i USING BTREE ON t (v)", true},
	{"CREATE VECTOR INDEX i ON t (v) USING HNSW", false},
	{"CREATE VECTOR INDEX i ON t (v) USING BTREE", true},
	{"CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(v)))", true},
	{"CREATE VECTOR INDEX i ON t (v, w)", true},
	{"CREATE VECTOR INDEX i ON t (v(3))", true},
	{"CREATE VECTOR INDEX i ON t (v DESC)", true},
	{"CREATE VECTOR INDEX i ON t (v) M=6", false},
	{"CREATE VECTOR INDEX i ON t (v) M=6 DISTANCE=cosine", false},
	{"CREATE VECTOR INDEX i ON t (v) DISTANCE=euclidean M=6", false},
	{"CREATE VECTOR INDEX i ON t (v) M=abc", false},
	{"CREATE VECTOR INDEX i ON t (v) DISTANCE=foo", false},
	{"CREATE VECTOR INDEX i ON t (v) COMMENT 'c'", true},
	{"CREATE VECTOR INDEX i ON t (v) ALGORITHM=DEFAULT", false},
	{"CREATE VECTOR INDEX i ON t (v) ALGORITHM=INPLACE LOCK=NONE", false},
	{"CREATE VECTOR INDEX i ON t (v) LOCK=NONE", false},
	{"CREATE VECTOR INDEX i ON t (v) WAIT 3", false},
	{"CREATE VECTOR INDEX i ON t (v) NOWAIT", false},
	{"CREATE VECTOR INDEX i WAIT 3 ON t (v)", false},
	{"CREATE VECTOR INDEX i ON t (v) M=6 ALGORITHM=DEFAULT", false},
	{"CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)", true},
	{"CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'", true},
	{"CREATE VECTOR INDEX i ON t (v) WITH (distance=cosine, type=hnsw)", true},
	{"CREATE VECTOR INDEX i ON t (v) WITH (foo=bar)", true},
	{"CREATE VECTOR INDEX i ON t (v) WITH ()", false},
	{"CREATE VECTOR INDEX i ON t (v) WITH (distance=l2,)", false},
	{"CREATE VECTOR INDEX i ON t (v) WITH (distance=l2) ALGORITHM=INPLACE", false},
	{"CREATE VECTOR INDEX i ON t (v) SELECT 1", false},
	{"CREATE VECTOR INDEX i ON t (v) x", false},
	{"CREATE VECTOR INDEX i ON t (v),", false},
	{"CREATE VECTOR INDEX vector ON vector (vector)", true},
	{"CREATE VECTOR INDEX `index` ON t (v)", true},
	{"CREATE VECTOR INDEX select ON t (v)", false},
	// as a word.
	{"CREATE TABLE vector (vector INT)", true},
	{"CREATE TABLE vector (vector vector(3))", true},
	{"CREATE TABLE t (vector)", false},
	{"CREATE TABLE t (vector index)", false},
	{"CREATE TABLE t (vector index INT)", false},
	{"CREATE TABLE t (vector INT, INDEX (vector))", true},
	{"SELECT 1 vector", true},
	{"CREATE INDEX vector ON t (vector)", true},
	{"DROP INDEX vector ON t", true},
	// other places.
	{"SELECT CAST('[1]' AS VECTOR)", false},
	{"SELECT CAST('[1]' AS VECTOR(3))", false},
	{"SELECT CAST('[1]' AS VECTOR(3) ARRAY)", false},
	{"SELECT CONVERT('[1]', VECTOR)", false},
	{"SELECT CONVERT('[1]', VECTOR(3))", false},
	{"ALTER TABLE t ADD COLUMN a VECTOR(3)", true},
	{"ALTER TABLE t ADD a VECTOR(3)", true},
	{"ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)", true},
	{"ALTER TABLE t ADD vector vector(3)", true},
	{"ALTER TABLE t ADD vector INT", true},
	{"ALTER TABLE t ADD COLUMN vector INT", true},
	{"ALTER TABLE t ADD VECTOR", false},
	{"ALTER TABLE t ADD VECTOR (a)", false},
	{"ALTER TABLE t ADD VECTOR KEY (a)", false},
	{"ALTER TABLE t ADD VECTOR KEY i (a)", false},
	{"ALTER TABLE t ADD VECTOR INDEX (a)", false},
	{"ALTER TABLE t ADD VECTOR INDEX i (a)", false},
	{"ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS i (a)", false},
	{"ALTER TABLE t ADD VECTOR i (a)", false},
	{"ALTER TABLE t ADD VECTOR INDEX i (a) M=6 DISTANCE=cosine", false},
	{"ALTER TABLE t ADD VECTOR INDEX i (a) WITH (distance=l2, type=hnsw, lib=vsag)", false},
	{"ALTER TABLE t ADD VECTOR KEY i (a) WITH (distance=l2, type=hnsw, lib=vsag)", false},
	{"ALTER TABLE t ADD VECTOR INDEX i (a) USING HNSW", false},
	{"ALTER TABLE t ADD VECTOR INDEX ((a))", false},
	{"ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX (a)", false},
	{"ALTER TABLE t MODIFY a VECTOR(3)", true},
	{"ALTER TABLE t MODIFY COLUMN a VECTOR(3) NOT NULL", true},
	{"ALTER TABLE t MODIFY a VECTOR", false},
	{"ALTER TABLE t CHANGE a b VECTOR(3)", true},
	{"ALTER TABLE t DROP INDEX i", true},
	{"CREATE VECTOR INDEX i1 ON t (v) COMMENT 'c' WITH (distance=l2, type=hnsw, lib=vsag)", true},
	{"CREATE VECTOR INDEX i2 ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) WITH (distance=l2, type=hnsw, lib=vsag)", true},
	{"CREATE VECTOR INDEX i3 ON t (v) USING BTREE WITH (distance=l2, type=hnsw, lib=vsag)", true},
	{"CREATE VECTOR INDEX i4 ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) USING BTREE", true},
	{"CREATE VECTOR INDEX i5 USING BTREE ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)", true},
	{"CREATE VECTOR INDEX i6 USING HASH ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)", true},
	{"CREATE VECTOR INDEX i7 ON t (v) WITH (distance=l2, distance=l2)", true},
	{"CREATE VECTOR INDEX i8 ON t (v) WITH (distance='l2', type='hnsw', lib='vsag')", false},
	{"CREATE VECTOR INDEX i9 ON t (v) WITH (m=16)", true},
	{"CREATE VECTOR INDEX i10 ON t (v) WITH (m='16')", false},
	{"CREATE VECTOR INDEX i11 ON t (v) WITH (ef_search=1.5)", true},
	{"CREATE VECTOR INDEX i12 ON t (v) WITH (`distance`=l2)", true},
	{"CREATE VECTOR INDEX i13 ON t (v) WITH (distance=`l2`)", true},
	{"CREATE VECTOR INDEX i14 ON t (v) WITH (distance l2)", false},
	{"CREATE VECTOR INDEX i15 ON t (v) WITH (distance=select)", false},
	{"CREATE VECTOR INDEX i16 ON t (v) WITH (select=1)", false},
	{"CREATE VECTOR INDEX i17 ON t (v) WITH (distance=)", false},
	{"CREATE VECTOR INDEX i18 ON t (v) WITH (=l2)", false},
	{"CREATE VECTOR INDEX i19 ON t (v) WITH (a=-1)", false},
	{"CREATE VECTOR INDEX i20 ON t (v) WITH (a=+1)", false},
	{"CREATE VECTOR INDEX i21 ON t (v) WITH (a=NULL)", false},
	{"CREATE VECTOR INDEX i22 ON t (v) WITH (a=DEFAULT)", false},
	{"CREATE VECTOR INDEX i23 ON t (v) WITH (a=TRUE)", true},
	{"CREATE VECTOR INDEX i24 ON t (v) WITH (a=(1))", false},
	{"CREATE VECTOR INDEX i25 ON t (v) WITH (a=1+1)", false},
	{"CREATE VECTOR INDEX i26 ON t (v) WITH (a=b.c)", false},
	{"CREATE VECTOR INDEX i30 ON t (v) WITH (type=hnsw) GLOBAL", true},
	{"CREATE VECTOR INDEX i31 ON t (v) WITH (type=hnsw) LOCAL", true},
	{"CREATE VECTOR INDEX i32 ON t (v) GLOBAL WITH (type=hnsw)", true},
	{"CREATE VECTOR INDEX i34 ON t (v) WITH (type=hnsw) BLOCK_SIZE 16384", true},
	{"CREATE VECTOR INDEX i35 ON t (v) WITH (type=hnsw) KEY_BLOCK_SIZE=1", true},
	{"CREATE VECTOR INDEX i36 ON t (v) WITH (type=hnsw) VISIBLE", true},
	{"CREATE VECTOR INDEX i37 ON t (v) INVISIBLE WITH (type=hnsw)", true},
	{"CREATE VECTOR INDEX i38 ON t (v) WITH PARSER p", false},
	{"CREATE VECTOR INDEX i39 ON t (v) WITH", false},
	{"CREATE VECTOR INDEX i40 ON t (v) WITH (", false},
	{"CREATE VECTOR INDEX i41 ON t (v) WITH (a=1", false},
	{"CREATE VECTOR INDEX i42 ON t (v) WITH (a=1),", false},
	{"CREATE VECTOR INDEX i43 ON t (v) WITH (a=1) x", false},
	{"CREATE VECTOR INDEX i44 ON d.t (v) WITH (a=1)", true},
	{"CREATE VECTOR INDEX IF EXISTS i46 ON t (v)", false},
	{"CREATE VECTOR INDEX i47 ON t (v ASC)", true},
	{"CREATE VECTOR INDEX i48 ON t (v) WITH (TYPE=HNSW, DISTANCE=L2, LIB=VSAG)", true},
	{"ALTER TABLE t ADD (VECTOR INDEX i (v))", false},
	{"SELECT `vector`(3)", true},
	{"CREATE TABLE s6 (a VECTOR(2147483647))", true},
	{"CREATE TABLE s7 (a VECTOR(2147483648))", true},
	{"CREATE TABLE s8 (a VECTOR(4294967295))", true},
	{"CREATE TABLE s9 (a VECTOR(4294967296))", true},
	{"CREATE TABLE s10 (a VECTOR(9223372036854775807))", true},
	{"CREATE TABLE s11 (a VECTOR(9223372036854775808))", true},
	{"CREATE TABLE s12 (a VECTOR(18446744073709551615))", true},
	{"CREATE TABLE s13 (a VECTOR(18446744073709551616))", false},
	{"CREATE VECTOR INDEX w1 ON t (v) WITH (a=FALSE)", true},
	{"CREATE VECTOR INDEX w2 ON t (v) WITH (a=0x10)", false},
	{"CREATE VECTOR INDEX w3 ON t (v) WITH (a=1e5)", true},
	{"CREATE VECTOR INDEX w4 ON t (v) WITH (a=.5)", true},
	{"SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a VECTOR(3) PATH '$')) AS j", true},
	{"SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a INT PATH '$')) AS j", true},
	{"CREATE FUNCTION f1(a VECTOR(3)) RETURNS INT RETURN 1", false},
	// WITH (...) is an option of every index, refused later on any but a vector one.
	{"CREATE INDEX i50 ON t (a) WITH (foo=bar)", true},
	{"CREATE UNIQUE INDEX i51 ON t (a) WITH (foo=bar)", true},
	{"ALTER TABLE t ADD INDEX i56 (a) WITH (foo=bar)", true},
	{"CREATE TABLE tk1 (a INT, KEY k (a) WITH (foo=bar))", true},
	{"CREATE TABLE tk2 (a INT, UNIQUE KEY k (a) WITH (foo=bar))", true},
	// a vector index as a table element, in full.
	{"CREATE TABLE tv1 (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c' WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE tv2 (v VECTOR(3) NOT NULL, VECTOR INDEX `key` (v) WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE tv3 (v VECTOR(3) NOT NULL, VECTOR INDEX select (v))", false},
	{"CREATE TABLE tv4 (vector VECTOR(3) NOT NULL, VECTOR INDEX vector (vector) WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE tv5 (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) LOCAL WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE tv6 (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) GLOBAL WITH (distance=l2, type=hnsw, lib=vsag))", true},
	{"CREATE TABLE tv7 (v VECTOR(3) NOT NULL, VECTOR INDEX i USING BTREE (v) WITH (distance=l2, type=hnsw, lib=vsag))", false},
	{"CREATE TABLE tv8 (v VECTOR(3) NOT NULL, VECTOR KEY USING BTREE (v) WITH (distance=l2, type=hnsw, lib=vsag))", false},
	{"CREATE TABLE tv9 (v VECTOR(3) NOT NULL, VECTOR INDEX USING HASH (v) WITH (distance=l2, type=hnsw, lib=vsag))", false},
	{"CREATE TABLE tv10 (vector INT KEY, b INT)", true},
	{"CREATE TABLE tv11 (vector key)", false},
	{"CREATE TABLE tv12 (vector key INT)", false},
	{"ALTER TABLE t ADD (b2 VECTOR(3), VECTOR INDEX i (b2) WITH (distance=l2, type=hnsw, lib=vsag))", false},
	// VECTOR is a column type: JSON_TABLE takes it, a routine does not.
	{"CREATE FUNCTION f2() RETURNS VECTOR(3) RETURN 1", false},
	{"CREATE PROCEDURE p2(a VECTOR(3)) SELECT 1", false},
	{"CREATE PROCEDURE p3() BEGIN DECLARE b VECTOR(3); END", false},
}

// TestVectorForms holds the parser to the recorded verdicts, without the server CI does
// not have.
func TestVectorForms(t *testing.T) {
	for _, c := range vectorFormCases {
		t.Run(c.sql, func(t *testing.T) {
			omniAccepts, err := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestVectorFormsOceanBaseOracle holds the record to the server.
func TestVectorFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS obvec",
		"CREATE DATABASE obvec",
		"USE obvec",
		"CREATE TABLE t (a INT, b INT, v VECTOR(3) NOT NULL)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() { _, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obvec") })
	for _, c := range vectorFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := engineAccepts(ctx, db, c.sql)
			if err != nil {
				t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("OceanBase accepts %q = %v, the record says %v", c.sql, got, c.accept)
			}
		})
	}
}

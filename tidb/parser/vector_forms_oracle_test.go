package parser

import "testing"

// The VECTOR column type and the VECTOR INDEX statement, as TiDB v8.5.5 reads them.
//
// The corpus pins the forms pingcap's own parser test writes. This holds the rest, and
// above all the shapes the server refuses, since a corpus of accepting rows cannot catch
// an over-broad reader. Every verdict was read from pingcap/tidb:v8.5.5 on 2026-11-05.
//
// A refusal here means what it means in the corpus: 1064, or a code a grammar action
// raises while parsing. CreateIndexStmt's and ConstraintVectorIndex's answer 1149 for a
// vector index that is not one functional key part of type HNSW (prepareVerdict).
//
// Left out because they are no part of this family: `SELECT vector(3)`, refused like
// every other unreserved keyword as a function name (`SELECT comment(3)` is 1064 and
// omni reads both), JSON_TABLE, which v8.5.5 has for no column type, and CAST(a AS INT
// ARRAY).
//
// Ref: CreateIndexStmt, IndexKeyTypeOpt, ConstraintVectorIndex, OptVectorElementType and
// the two `"VECTOR" OptVectorElementType OptFieldLen` arms in TiDB v8.5.5
// pkg/parser/parser.y.

type vectorFormCase struct {
	sql    string
	accept bool
}

var vectorFormCases = []vectorFormCase{
	// index shapes.
	{"CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a)))", true},
	{"CREATE VECTOR INDEX idx ON t (a)", false},
	{"CREATE VECTOR INDEX idx ON t (a) USING HNSW", false},
	{"CREATE VECTOR INDEX idx ON t (a, b) USING HNSW", false},
	{"CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE BTREE", false},
	{"CREATE VECTOR INDEX idx USING BTREE ON t ((VEC_COSINE_DISTANCE(a)))", false},
	{"CREATE VECTOR INDEX idx USING HASH ON t ((VEC_COSINE_DISTANCE(a)))", false},
	{"CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) USING HYPO", false},
	{"CREATE VECTOR INDEX idx ON t USING HNSW ((VEC_COSINE_DISTANCE(a)))", false},
	{"CREATE VECTOR idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW", false},
	{"CREATE VECTOR KEY idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW", false},
	{"CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE HNSW", true},
	{"CREATE VECTOR INDEX IF NOT idx ON t ((a))", false},
	{"CREATE VECTOR INDEX IF EXISTS idx ON t ((a))", false},
	{"CREATE VECTOR INDEX ident TYPE HNSW ON d_n.t_n ((VEC_COSINE_DISTANCE(a)))", true},
	{"CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a)))", true},
	{"CREATE VECTOR INDEX idx USING `HNSW` ON t ((VEC_COSINE_DISTANCE(a)))", false},
	{"CREATE VECTOR INDEX idx USING hnsw ON t ((VEC_COSINE_DISTANCE(a))) TYPE hnsw", true},
	{"CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a))) USING BTREE", false},
	{"CREATE VECTOR INDEX idx USING BTREE ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW", true},
	{"CREATE VECTOR INDEX ON t ((a))", false},
	{"CREATE VECTOR INDEX", false},
	{"CREATE VECTOR INDEX idx", false},
	{"CREATE VECTOR INDEX idx ON", false},
	{"CREATE VECTOR INDEX idx ON t", false},
	{"CREATE VECTOR INDEX idx ON t (", false},
	{"CREATE VECTOR INDEX idx ON t ()", false},
	{"CREATE VECTOR", false},
	{"CREATE VECTOR UNIQUE INDEX idx ON t ((a))", false},
	{"CREATE UNIQUE VECTOR INDEX idx ON t ((a))", false},
	{"CREATE FULLTEXT VECTOR INDEX idx ON t ((a))", false},
	{"CREATE TEMPORARY VECTOR INDEX idx ON t ((a))", false},
	{"CREATE GLOBAL TEMPORARY VECTOR INDEX idx ON t ((a))", false},
	{"CREATE OR REPLACE VECTOR INDEX idx ON t ((a))", false},
	{"CREATE VECTOR INDEX vector ON vector ((vec(vector)))", true},
	{"CREATE VECTOR INDEX `index` ON t ((a))", true},
	{"CREATE VECTOR INDEX select ON t ((a))", false},
	{"CREATE VECTOR INDEX idx ON t ((a) DESC)", true},
	{"CREATE VECTOR INDEX idx ON t ((a) ASC)", true},
	{"CREATE VECTOR INDEX idx ON t (a(3))", false},
	{"CREATE VECTOR INDEX idx ON t ((a) DESC) COMMENT 'c' ALGORITHM = INPLACE LOCK = NONE", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) COMMENT 'c'", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) USING HNSW COMMENT 'c' VISIBLE", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) INVISIBLE", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) GLOBAL", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) LOCAL", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) CLUSTERED", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) KEY_BLOCK_SIZE = 1", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) WITH PARSER p", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) ALGORITHM = DEFAULT", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) LOCK = SHARED ALGORITHM = COPY", true},
	{"CREATE VECTOR INDEX idx ON t ((a)) ALGORITHM = INPLACE ALGORITHM = COPY", false},
	{"CREATE VECTOR INDEX idx ON t ((a)) WITH (distance=l2)", false},
	{"CREATE VECTOR INDEX idx ON t ((a)) M=6", false},
	{"CREATE VECTOR INDEX idx ON t ((a)) DISTANCE=cosine", false},
	{"CREATE VECTOR INDEX idx ON t ((a)) ADD_COLUMNAR_REPLICA_ON_DEMAND", false},
	{"CREATE VECTOR INDEX idx ON t ((a)),", false},
	{"CREATE VECTOR INDEX idx ON t ((a)) SELECT 1", false},
	{"CREATE VECTOR INDEX idx ON t ((a)) x", false},
	{"CREATE VECTOR INDEX idx ON t ((a), (b)) USING HNSW", false},
	// HNSW on another kind of index.
	{"CREATE INDEX i USING HNSW ON t (a)", false},
	{"CREATE INDEX i ON t (a) USING HNSW", false},
	{"CREATE UNIQUE INDEX ident USING HNSW ON d_n.t_n ( ident , ident ASC )", false},
	{"CREATE FULLTEXT INDEX i ON t (a) USING HNSW", false},
	{"CREATE SPATIAL INDEX i ON t (a) USING HNSW", false},
	{"ALTER TABLE t ADD INDEX (a) USING HNSW", true},
	{"CREATE TABLE t (a INT, INDEX (a) USING HNSW)", true},
	// type shapes.
	{"CREATE TABLE t (a VECTOR)", true},
	{"CREATE TABLE t (a VECTOR(3))", true},
	{"CREATE TABLE t (a VECTOR(0))", true},
	{"CREATE TABLE t (a VECTOR(16383))", true},
	{"CREATE TABLE t (a VECTOR(16384))", true},
	{"CREATE TABLE t (a VECTOR(99999999999999999999))", false},
	{"CREATE TABLE t (a VECTOR<FLOAT>)", true},
	{"CREATE TABLE t (a VECTOR<float>)", true},
	{"CREATE TABLE t (a VECTOR<FLOAT>(3))", true},
	{"CREATE TABLE t (a VECTOR < FLOAT > ( 3 ))", true},
	{"CREATE TABLE t (a VECTOR<DOUBLE>)", false},
	{"CREATE TABLE t (a VECTOR<DOUBLE>(3))", false},
	{"CREATE TABLE t (a VECTOR<INT>)", false},
	{"CREATE TABLE t (a VECTOR<ABC>)", false},
	{"CREATE TABLE t (a VECTOR<>)", false},
	{"CREATE TABLE t (a VECTOR<FLOAT)", false},
	{"CREATE TABLE t (a VECTOR FLOAT>)", false},
	{"CREATE TABLE t (a VECTOR<FLOAT>>)", false},
	{"CREATE TABLE t (a VECTOR<<FLOAT>)", false},
	{"CREATE TABLE t (a VECTOR<`FLOAT`>)", false},
	{"CREATE TABLE t (a VECTOR<FLOAT(3)>)", false},
	{"CREATE TABLE t (a VECTOR(5)<FLOAT>)", false},
	{"CREATE TABLE t (a VECTOR())", false},
	{"CREATE TABLE t (a VECTOR(3,2))", false},
	{"CREATE TABLE t (a VECTOR('3'))", false},
	{"CREATE TABLE t (a VECTOR(b))", false},
	{"CREATE TABLE t (a VECTOR(1.5))", false},
	{"CREATE TABLE t (a VECTOR(0x3))", false},
	{"CREATE TABLE t (a VECTOR(-1))", false},
	{"CREATE TABLE t (a VECTOR(3)", false},
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
	{"CREATE TABLE t (a VECTOR AS ('[1]') STORED)", true},
	{"CREATE TABLE t (a `VECTOR`)", false},
	{"CREATE TABLE t (a VEC)", false},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))", true},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX i TYPE HNSW ((VEC_COSINE_DISTANCE(b))))", true},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX i USING HNSW ((VEC_COSINE_DISTANCE(b))))", true},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX TYPE HNSW ((VEC_COSINE_DISTANCE(b))))", false},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX USING HNSW ((VEC_COSINE_DISTANCE(b))))", true},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX IF NOT EXISTS i ((VEC_COSINE_DISTANCE(b))))", true},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX (b))", false},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))) USING BTREE)", false},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR KEY ((VEC_COSINE_DISTANCE(b))))", false},
	{"CREATE TABLE t (a int, b VECTOR(3), VECTOR (b))", false},
	{"CREATE TABLE t (a int, b VECTOR(3), CONSTRAINT c VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))", false},
	{"CREATE TABLE t (a int, b VECTOR(3), CONSTRAINT VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))", true},
	// VECTOR as a word.
	{"CREATE TABLE vector (vector vector)", true},
	{"CREATE TABLE vector (vector vector(3), vector2 vector<float>(3), VECTOR INDEX vector ((vec(vector))))", true},
	{"CREATE TABLE t (vector INT)", true},
	{"CREATE TABLE t (vector INT, INDEX (vector))", true},
	{"CREATE TABLE t (vector)", false},
	{"CREATE TABLE t (vector index)", false},
	{"CREATE TABLE t (vector index INT)", false},
	{"CREATE TABLE t (`vector` INT)", true},
	{"SELECT 1 vector", true},
	{"CREATE INDEX vector ON t (vector)", true},
	{"CREATE TABLE t (a int) COMMENT 'x' vector", false},
	{"DROP INDEX vector ON t", true},
	// other places a type goes.
	{"SELECT CAST(a AS VECTOR)", true},
	{"SELECT CAST(a AS VECTOR(3))", true},
	{"SELECT CAST(a AS VECTOR<FLOAT>)", true},
	{"SELECT CAST(a AS VECTOR<FLOAT>(3))", true},
	{"SELECT CAST(a AS VECTOR<DOUBLE>)", false},
	{"SELECT CAST(a AS VECTOR<DOUBLE>(3))", false},
	{"SELECT CAST(a AS VECTOR(3)<FLOAT>)", false},
	{"SELECT CAST(a AS VECTOR())", false},
	{"SELECT CAST(a AS VECTOR(3) ARRAY)", true},
	{"SELECT CAST(a AS VECTOR(3) UNSIGNED)", false},
	{"SELECT CONVERT(a, VECTOR)", true},
	{"SELECT CONVERT(a, VECTOR(3))", true},
	{"SELECT CONVERT(a, VECTOR<FLOAT>(3))", true},
	{"SELECT CONVERT(a USING VECTOR)", true},
	{"SELECT a::VECTOR", false},
	{"ALTER TABLE t ADD COLUMN a VECTOR(3)", true},
	{"ALTER TABLE t ADD a VECTOR", true},
	{"ALTER TABLE t ADD a VECTOR<FLOAT>", true},
	{"ALTER TABLE t ADD a VECTOR<DOUBLE>", false},
	{"ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)", true},
	{"ALTER TABLE t ADD vector vector(3)", false},
	{"ALTER TABLE t ADD vector INT", false},
	{"ALTER TABLE t ADD COLUMN vector INT", true},
	{"ALTER TABLE t ADD vector int, ADD VECTOR INDEX ((a))", false},
	{"ALTER TABLE t ADD VECTOR", false},
	{"ALTER TABLE t ADD VECTOR (a)", false},
	{"ALTER TABLE t ADD VECTOR KEY ((a))", false},
	{"ALTER TABLE t ADD VECTOR INDEX ((a))", true},
	{"ALTER TABLE t ADD VECTOR INDEX (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX ((a))", false},
	{"ALTER TABLE t MODIFY a VECTOR(3)", true},
	{"ALTER TABLE t MODIFY COLUMN a VECTOR<FLOAT>(3) NOT NULL", true},
	{"ALTER TABLE t MODIFY a VECTOR<INT>", false},
	{"ALTER TABLE t CHANGE a b VECTOR(3)", true},
	{"ALTER TABLE t CHANGE COLUMN a b VECTOR", true},
	{"CREATE FUNCTION f(a VECTOR) RETURNS VECTOR RETURN a", false},
	{"CREATE TABLE t2 LIKE t", true},
	{"CREATE TEMPORARY TABLE tt (a VECTOR(3))", true},
	{"SELECT VEC_COSINE_DISTANCE(a, '[1,2,3]') FROM t ORDER BY VEC_COSINE_DISTANCE(a, '[1,2,3]') LIMIT 3", true},
	{"SELECT VEC_FROM_TEXT('[1]'), VEC_DIMS(a), VEC_L2_NORM(a), VEC_AS_TEXT(a) FROM t", true},
	{"ALTER TABLE t SET TIFLASH REPLICA 1", true},
	// the column named vector, which a bare ADD cannot read and everything else can.
	{"ALTER TABLE t ADD `vector` INT", true},
	{"ALTER TABLE t ADD comment INT", true},
	{"ALTER TABLE t ADD COLUMN IF NOT EXISTS vector INT", true},
	{"ALTER TABLE t MODIFY vector INT", true},
	{"ALTER TABLE t CHANGE vector vector INT", true},
	{"ALTER TABLE t ADD (vector INT)", true},
	{"ALTER TABLE t ADD (vector VECTOR(3))", true},
	{"ALTER TABLE t ADD (VECTOR INDEX ((a)))", true},
	{"ALTER TABLE t ADD COLUMN (VECTOR INDEX ((a)))", true},
	{"ALTER TABLE t ADD COLUMN IF NOT EXISTS (VECTOR INDEX ((a)))", true},
	// the dimension is any unsigned 64-bit number here, and the server checks it later.
	{"CREATE TABLE t (a VECTOR(4294967296))", true},
	{"CREATE TABLE t (a VECTOR(9223372036854775808))", true},
	{"CREATE TABLE t (a VECTOR(18446744073709551615))", true},
	{"CREATE TABLE t (a VECTOR(18446744073709551616))", false},
	{"CREATE TABLE t (a VECTOR(1,))", false},
	{"CREATE TABLE t (a VECTOR<FLOAT>())", false},
	{"CREATE TABLE t (a VECTOR<FLOAT> NOT NULL)", true},
	{"CREATE TABLE t (a VECTOR/* c */<FLOAT>/* d */(3))", true},
	{"SELECT CAST(a AS VECTOR<FLOAT> ARRAY)", true},
	{"CREATE TEMPORARY TABLE t (a int, VECTOR INDEX ((a)))", true},
	{"CREATE TABLE t (a VECTOR(3) CHECK (a IS NOT NULL))", true},
	{"CREATE VIEW vv AS SELECT CAST(a AS VECTOR)", true},
	// CONSTRAINT takes a symbol named vector and an ordinary index.
	{"CREATE TABLE t (a VECTOR(3), CONSTRAINT vector UNIQUE (a))", true},
}

func TestVectorFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	for _, c := range vectorFormCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbRelease, accepts, c.accept, c.sql)
			}
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbRelease, c.accept, c.sql, perr)
			}
		})
	}
}

// TestVectorForms holds the parser to the recorded verdicts without the server, so a
// regression fails however the suite is run.
func TestVectorForms(t *testing.T) {
	for _, c := range vectorFormCases {
		t.Run(c.sql, func(t *testing.T) {
			omniAccepts, perr := parsesAsOneStatement(c.sql)
			if omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, perr)
			}
		})
	}
}

package parser

import (
	"context"
	"testing"
)

// The index_option list and the index name of OceanBase, past what MySQL has: PARALLEL,
// NOPARALLEL, STORING, WITH ROWID, the table ids and STORAGE_CACHE_POLICY, a qualified
// name on CREATE INDEX, and CONSTRAINT with the four kinds of key it goes with and no
// others. Every verdict was read from oceanbase/oceanbase-ce:4.3.5-lts on 2026-11-05.
//
// Ref: index_option, parallel_option, normal_relation_factor and table_element in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
//
// Left out because they are no index option: CREATE INDEX ... PARTITION BY and WITH COLUMN
// GROUP, which follow the option list, PARALLEL as a table option, and `ADD (KEY k (a))`, a
// key in parentheses with any option or none. All of them parse on 4.3.5.6 and omni does
// not read them yet.

type indexOptionFormCase struct {
	sql    string
	accept bool
}

var indexOptionFormCases = []indexOptionFormCase{
	// PARALLEL.
	{"CREATE INDEX p1 ON t (a) PARALLEL 4", true},
	{"CREATE INDEX p2 ON t (a) PARALLEL = 4", true},
	{"CREATE INDEX p3 ON t (a) PARALLEL=1", true},
	{"CREATE INDEX p4 ON t (a) PARALLEL a", false},
	{"CREATE INDEX p5 ON t (a) PARALLEL 0", false},
	{"CREATE INDEX p6 ON t (a) PARALLEL", false},
	{"CREATE INDEX p7 ON t (a) PARALLEL 1.5", false},
	{"CREATE INDEX p8 ON t (a) PARALLEL '4'", false},
	{"CREATE INDEX p9 ON t (a) PARALLEL -1", false},
	{"CREATE INDEX p10 ON t (a) PARALLEL =", false},
	{"CREATE INDEX p11 ON t (a) PARALLEL 2 PARALLEL 3", true},
	{"CREATE INDEX p12 ON t (a) PARALLEL 2 COMMENT 'c' NOPARALLEL", true},
	{"CREATE INDEX p13 ON t (a) NOPARALLEL", true},
	{"CREATE INDEX p14 ON t (a) NOPARALLEL 4", false},
	{"CREATE INDEX p15 ON t (a) NOPARALLEL = 1", false},
	{"CREATE INDEX `parallel` ON t (a)", true},
	{"CREATE INDEX parallel ON t (a)", true},
	{"CREATE INDEX noparallel ON t (a)", true},
	{"CREATE INDEX storing ON t (a)", true},
	{"CREATE INDEX rowid ON t (a)", true},
	{"CREATE TABLE parallel (parallel INT, noparallel INT, storing INT, rowid INT, data_table_id INT, index_table_id INT, virtual_column_id INT, max_used_part_id INT, storage_cache_policy INT)", true},
	// STORING.
	{"CREATE INDEX s1 ON t (a) STORING (b)", true},
	{"CREATE INDEX s2 ON t (a) STORING (b, c)", true},
	{"CREATE INDEX s3 ON t (a) STORING ()", false},
	{"CREATE INDEX s4 ON t (a) STORING b", false},
	{"CREATE INDEX s5 ON t (a) STORING (`b`)", true},
	{"CREATE INDEX s6 ON t (a) STORING (t.b)", false},
	{"CREATE INDEX s7 ON t (a) STORING (nosuch)", true},
	{"CREATE INDEX s8 ON t (a) STORING (b,)", false},
	{"CREATE INDEX s9 ON t (a) STORING (b DESC)", false},
	{"CREATE INDEX s10 ON t (a) STORING (b(3))", false},
	{"CREATE INDEX s11 ON t (a) STORING (b) STORING (c)", true},
	{"CREATE INDEX s12 ON t (a) STORING", false},
	{"CREATE INDEX s13 ON t (a) STORING (select)", false},
	{"CREATE INDEX s14 ON t (a) `STORING` (b)", false},
	{"CREATE INDEX s15 ON t (a) `PARALLEL` 2", false},
	{"CREATE INDEX s16 ON t (a) `NOPARALLEL`", false},
	// ROWID and the ids.
	{"CREATE INDEX w1 ON t (a) WITH ROWID", true},
	{"CREATE INDEX w2 ON t (a) WITH ROW", false},
	{"CREATE INDEX w3 ON t (a) WITH", false},
	{"CREATE INDEX i1 ON t (a) DATA_TABLE_ID = 1", true},
	{"CREATE INDEX i2 ON t (a) DATA_TABLE_ID 1", true},
	{"CREATE INDEX i3 ON t (a) DATA_TABLE_ID a", false},
	{"CREATE INDEX i4 ON t (a) INDEX_TABLE_ID = 1", true},
	{"CREATE INDEX i5 ON t (a) VIRTUAL_COLUMN_ID = 1", false},
	{"CREATE INDEX i6 ON t (a) MAX_USED_PART_ID = 1", true},
	{"CREATE INDEX i7 ON t (a) INDEX_TABLE_ID", false},
	{"CREATE INDEX c1 ON t (a) STORAGE_CACHE_POLICY = (GLOBAL = 'hot')", true},
	{"CREATE INDEX c2 ON t (a) STORAGE_CACHE_POLICY (GLOBAL = 'hot')", true},
	{"CREATE INDEX c3 ON t (a) STORAGE_CACHE_POLICY = (GLOBAL = hot)", false},
	{"CREATE INDEX c4 ON t (a) STORAGE_CACHE_POLICY = ()", false},
	{"CREATE INDEX c5 ON t (a) STORAGE_CACHE_POLICY = (BOUNDARY_COLUMN = a, HOT_RETENTION = 1 DAY)", true},
	{"CREATE INDEX c6 ON t (a) STORAGE_CACHE_POLICY = (BOUNDARY_COLUMN = a, BOUNDARY_COLUMN_UNIT = 's', HOT_RETENTION = 1 DAY)", true},
	{"CREATE INDEX c7 ON t (a) STORAGE_CACHE_POLICY = 'hot'", false},
	// every index position.
	{"CREATE TABLE n1 (a INT, b INT, KEY k (a) PARALLEL 2)", true},
	{"CREATE TABLE n2 (a INT, b INT, KEY k (a) NOPARALLEL)", true},
	{"CREATE TABLE n3 (a INT, b INT, KEY k (a) STORING (b))", true},
	{"CREATE TABLE n4 (a INT, b INT, KEY k (a) WITH ROWID)", true},
	{"CREATE TABLE n5 (a INT, b INT, UNIQUE KEY k (a) PARALLEL 2)", true},
	{"CREATE TABLE n6 (a INT, b INT, UNIQUE KEY k (a) STORING (b))", true},
	{"CREATE TABLE n7 (a INT, c TEXT, FULLTEXT KEY k (c) PARALLEL 2)", true},
	{"CREATE TABLE n8 (a INT, c TEXT, FULLTEXT KEY k (c) STORING (a))", true},
	{"CREATE TABLE n9 (a INT, g GEOMETRY NOT NULL, SPATIAL KEY k (g) PARALLEL 2)", true},
	{"CREATE TABLE n10 (a INT, g GEOMETRY NOT NULL, SPATIAL KEY k (g) STORING (a))", true},
	{"CREATE TABLE n11 (a INT, b INT, PRIMARY KEY (a) PARALLEL 2)", false},
	{"CREATE TABLE n12 (a INT, b INT, PRIMARY KEY (a) STORING (b))", false},
	{"CREATE TABLE n13 (a INT, b INT, CONSTRAINT c PRIMARY KEY (a) NOPARALLEL)", false},
	{"CREATE TABLE n14 (a INT, b INT, CONSTRAINT c UNIQUE KEY k (a) PARALLEL 2)", true},
	{"CREATE TABLE n15 (a INT, b INT, CONSTRAINT c UNIQUE KEY k (a) STORING (b))", true},
	{"CREATE TABLE n16 (a INT PARALLEL 2)", false},
	{"ALTER TABLE t ADD INDEX a1 (a) PARALLEL 4", true},
	{"ALTER TABLE t ADD INDEX a2 (a) NOPARALLEL", true},
	{"ALTER TABLE t ADD INDEX a3 (a) STORING (b)", true},
	{"ALTER TABLE t ADD UNIQUE KEY a4 (a) STORING (b)", true},
	{"ALTER TABLE t ADD UNIQUE KEY a5 (a) PARALLEL 4", true},
	{"ALTER TABLE t ADD FULLTEXT KEY a6 (c) PARALLEL 4", true},
	{"ALTER TABLE t ADD CONSTRAINT c7 UNIQUE KEY a7 (a) PARALLEL 4", true},
	{"ALTER TABLE t ADD PRIMARY KEY (a) PARALLEL 4", true},
	{"ALTER TABLE t ADD CONSTRAINT c8 PRIMARY KEY (a) STORING (b)", true},
	{"ALTER TABLE t ADD SPATIAL KEY a9 (g) STORING (a)", true},
	{"ALTER TABLE t ADD INDEX a11 (a) PARALLEL 2, ADD INDEX a12 (b) STORING (a)", true},
	{"CREATE UNIQUE INDEX u1 ON t (a) STORING (b)", true},
	{"CREATE UNIQUE INDEX u2 ON t (a) PARALLEL 2", true},
	{"CREATE FULLTEXT INDEX f1 ON t (c) PARALLEL 2", true},
	{"CREATE FULLTEXT INDEX f2 ON t (c) STORING (a)", true},
	{"CREATE SPATIAL INDEX g1 ON t (g) PARALLEL 2", true},
	{"CREATE SPATIAL INDEX g2 ON t (g) STORING (a)", true},
	{"CREATE VECTOR INDEX v1 ON t (a) WITH (distance=l2, type=hnsw) PARALLEL 2", true},
	{"CREATE VECTOR INDEX v2 ON t (a) STORING (b) WITH (distance=l2, type=hnsw)", true},
	// qualified names.
	{"CREATE INDEX probe.q1 ON t (a)", true},
	{"CREATE INDEX obo.q2 ON t (a)", true},
	{"CREATE INDEX nosuch.q3 ON t (a)", true},
	{"CREATE INDEX x.y.q4 ON t (a)", true},
	{"CREATE INDEX w.x.y.q5 ON t (a)", false},
	{"CREATE INDEX .q6 ON t (a)", false},
	{"CREATE INDEX probe.select ON t (a)", true},
	{"CREATE INDEX select.q7 ON t (a)", true},
	{"CREATE INDEX `probe`.`q8` ON t (a)", true},
	{"CREATE INDEX probe. q9 ON t (a)", true},
	{"CREATE INDEX probe.q10 ON obo.ot (a)", true},
	{"CREATE INDEX obo.q11 ON obo.ot (a)", true},
	{"CREATE INDEX IF NOT EXISTS probe.q12 ON t (a)", true},
	{"CREATE UNIQUE INDEX probe.q13 ON t (a)", true},
	{"CREATE FULLTEXT INDEX probe.q14 ON t (c)", true},
	{"CREATE VECTOR INDEX probe.q15 ON t (a) WITH (distance=l2, type=hnsw)", true},
	{"CREATE INDEX probe.q16 USING BTREE ON t (a)", true},
	{"CREATE INDEX q17@l ON t (a)", true},
	{"CREATE INDEX probe.q18@l ON t (a)", true},
	{"CREATE TABLE q19 (a INT, KEY probe.k (a))", false},
	{"ALTER TABLE t ADD INDEX probe.q20 (a)", false},
	{"DROP INDEX probe.q1 ON t", false},
	// CONSTRAINT.
	{"CREATE TABLE m1 (a INT, CONSTRAINT c INDEX (a))", false},
	{"CREATE TABLE m2 (a INT, CONSTRAINT c KEY (a))", false},
	{"CREATE TABLE m3 (a INT, CONSTRAINT INDEX (a))", false},
	{"CREATE TABLE m4 (a INT, CONSTRAINT KEY (a))", false},
	{"CREATE TABLE m5 (a INT, CONSTRAINT c KEY k (a))", false},
	{"CREATE TABLE m6 (c TEXT, CONSTRAINT x FULLTEXT (c))", false},
	{"CREATE TABLE m7 (c TEXT, CONSTRAINT FULLTEXT (c))", false},
	{"CREATE TABLE m8 (g GEOMETRY NOT NULL, CONSTRAINT x SPATIAL (g))", false},
	{"CREATE TABLE m9 (g GEOMETRY NOT NULL, CONSTRAINT SPATIAL (g))", false},
	{"CREATE TABLE m10 (a INT, CONSTRAINT VECTOR INDEX (a))", false},
	{"CREATE TABLE m11 (a INT, CONSTRAINT x VECTOR INDEX (a))", false},
	{"CREATE TABLE m12 (a INT, CONSTRAINT x UNIQUE (a))", true},
	{"CREATE TABLE m13 (a INT, CONSTRAINT UNIQUE (a))", true},
	{"CREATE TABLE m14 (a INT, CONSTRAINT x PRIMARY KEY (a))", true},
	{"CREATE TABLE m15 (a INT, CONSTRAINT x FOREIGN KEY (a) REFERENCES obo.ot (a))", true},
	{"CREATE TABLE m16 (a INT, CONSTRAINT x CHECK (a > 0))", true},
	{"CREATE TABLE m17 (a INT, CONSTRAINT)", false},
	{"CREATE TABLE m18 (a INT, CONSTRAINT x)", false},
	{"ALTER TABLE t ADD CONSTRAINT x INDEX (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT x KEY k (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT INDEX (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT KEY (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT x FULLTEXT (c)", false},
	{"ALTER TABLE t ADD CONSTRAINT x SPATIAL (g)", false},
	{"ALTER TABLE t ADD CONSTRAINT FULLTEXT (c)", false},
	{"ALTER TABLE t ADD (CONSTRAINT x INDEX (a))", false},
	{"ALTER TABLE t ADD COLUMN (d INT, CONSTRAINT x KEY (d))", false},
}

// TestIndexOptionForms holds the parser to the recorded verdicts, without the server CI
// does not have.
func TestIndexOptionForms(t *testing.T) {
	for _, c := range indexOptionFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if omniAccepts := err == nil; omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestIndexOptionFormsOceanBaseOracle holds the record to the server.
func TestIndexOptionFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS obo",
		"DROP DATABASE IF EXISTS probe",
		"CREATE DATABASE obo",
		"CREATE DATABASE probe",
		"USE probe",
		"CREATE TABLE t (a INT, b INT, c TEXT, g GEOMETRY NOT NULL)",
		"CREATE TABLE obo.ot (a INT, b INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obo")
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS probe")
	})
	for _, c := range indexOptionFormCases {
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

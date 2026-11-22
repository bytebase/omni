package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestWave5CompatReadings pins the Wave 5 grammar and byte-offset families
// moved to Oracle 23.26.3's exact byte offsets, together with valid Oracle SQL
// controls beside each refusal.
func TestWave5CompatReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// CREATE PROPERTY GRAPH duplicate element tables/aliases, KEY, EXCEPT, and PROPERTIES.
		{`CREATE PROPERTY GRAPH a VERTEX TABLES(a) EDGE TABLES(a SOURCE a(a) DESTINATION a(a))`, 53},
		{`CREATE PROPERTY GRAPH a VERTEX TABLES(a KEY(a, a))`, 47},
		{`CREATE PROPERTY GRAPH a VERTEX TABLES(a PROPERTIES ALL COLUMNS EXCEPT(a, a))`, 73},
		{`CREATE PROPERTY GRAPH a VERTEX TABLES(a PROPERTIES(a, a))`, 54},

		// CREATE VECTOR INDEX INCLUDE duplicate columns & CREATE DIMENSION duplicate LEVEL.
		{`CREATE VECTOR INDEX a ON a(a) INCLUDE(a, a)`, 38},
		{`CREATE DIMENSION a LEVEL a IS(a.a) LEVEL a IS(a.a) JOIN KEY(a) REFERENCES a`, 41},

		// CREATE TABLE / ALTER TABLE duplicate lists, virtual column datatypes, identity, bounds & clauses.
		{`CREATE TABLE a(a, a)`, 18},
		{`CREATE TABLE a(UNIQUE(a, a))`, 25},
		{`CREATE TABLE t ("A" CLOB) LOB ("A", a) STORE AS (DISABLE STORAGE IN ROW)`, 36},
		{`ALTER TABLE t ADD ("A" CLOB) LOB ("A", a) STORE AS (DISABLE STORAGE IN ROW)`, 39},
		{`CREATE TABLE a(a zzz AS(1))`, 15},
		{`CREATE TABLE a(a ALWAYS AS IDENTITY)`, 35},
		{`CREATE TABLE a(a NUMBER) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN(1, 2))`, 81},
		{`CREATE TABLE a(CONSTRAINT a zzz UNIQUE(a))`, 32},
		{`CREATE TABLE a(SUPPLEMENTAL GROUP a(a))`, 35},
		{`ALTER TABLE ADD (a NUMBER)`, 12},
		{`ALTER TABLE a ENABLE UNIQUE(a, a)`, 31},
		{`ALTER TABLE a DROP PARTITION FORCE(1)`, 37},
		{`ALTER TABLE a ADD(a zzz AS(1))`, 18},
		{`ALTER TABLE a ADD DEFERRABLE INITIALLY DEFERRED`, 47},
		{`ALTER TABLE a IMMEDIATE DELETE_ALL`, 24},
		{`ALTER TABLE a UPGRADE INCLUDING DATA(a)`, 37},

		// CREATE / ALTER VIEW, MATERIALIZED VIEW, ZONEMAP, CLUSTER, INDEX, OPERATOR, INDEXTYPE.
		{`CREATE VIEW a(a, a) AS SELECT 1 FROM dual`, 17},
		{`CREATE VIEW a ANNOTATIONS(a, a) AS SELECT 1 FROM dual`, 32},
		{`CREATE MATERIALIZED VIEW a(a, a) AS SELECT 1 FROM dual`, 27},
		{`ALTER MATERIALIZED VIEW a ALLOCATE EXTENT(SIZE)`, 46},
		{`ALTER MATERIALIZED VIEW a MERGE SUBPARTITIONS FORCE(1), a`, 54},
		{`CREATE MATERIALIZED ZONEMAP a ON a(a, a)`, 38},
		{`CREATE CLUSTER a(a NUMBER, a NUMBER)`, 27},
		{`ALTER INDEX a UNUSABLE zzz`, 14},
		{`CREATE OPERATOR a BINDING(OBJECT) RETURN NUMBER USING a`, 26},
		{`ALTER INDEXTYPE a ADD a(NUMBER, OBJECT)`, 32},

		// Profile / resource cost duplicates & PFILE / SPFILE syntax.
		{`CREATE PROFILE a LIMIT SESSIONS_PER_USER 1 SESSIONS_PER_USER 1`, 61},
		{`ALTER RESOURCE COST CPU_PER_SESSION 1 CPU_PER_SESSION 1`, 54},
		{`CREATE PFILE = FROM SPFILE`, 15},
		{`CREATE SPFILE FROM PFILE AS`, 25},

		// Query & DML readings (EXTERNAL duplicate columns, INSERT SET duplicates, FETCH, UPDATE FOR).
		{`SELECT * FROM EXTERNAL((a NUMBER, a NUMBER) TYPE a DEFAULT DIRECTORY a LOCATION('s'))`, 34},
		{`INSERT INTO a SET a = 1, a = 1`, 25},
		{`SELECT * FETCH FIRST 1 zzz ROW ONLY`, 27},
		{`UPDATE a PARTITIONING FOR(1) SET a = 1`, 25},

		// Utility statements (LOCK TABLE, FLASHBACK TABLE, TRUNCATE).
		{`LOCK TABLE a IN MODE`, 20},
		{`FLASHBACK TABLE a SCN 1`, 22},
		{`TRUNCATE TABLE a PRESERVE VIEW LOG`, 31},

		// Other DDL (DISKGROUP, DATABASE, PDB, USER, KEY MANAGEMENT, AUDIT POLICY, INMEMORY JOIN GROUP, JAVA).
		{`ALTER DISKGROUP a ADD VOLUME 1`, 30},
		{`ALTER DATABASE COMMIT SWITCHOVER TO STANDBY`, 33},
		{`ALTER PLUGGABLE DATABASE 's' ONLINE`, 29},
		{`CREATE USER a IDENTIFIED EXTERNALLY ENABLE`, 42},
		{`ADMINISTER KEY MANAGEMENT SET KEY IDENTIFIED BY REVERSE MIGRATE USING 's'`, 56},
		{`CREATE AUDIT POLICY a ACTIONS COMPONENT = DV ALL zzz`, 45},
		{`CREATE INMEMORY JOIN GROUP a(a(a) zzz)`, 37},
		{`CREATE JAVA SOURCE IF EXISTS NAMED omni_jie_R1 AS class X {}`, 19},
		{`CREATE OR REPLACE JAVA SOURCE IF NOT EXISTS NAMED omni_jor_R1 AS class X {}`, 44},

		// Lexer error offsets & reserved DATE in expressions.
		{`SELECT 1 "" FROM dual`, 9},
		{`SELECT q'[a' FROM dual`, 7},
		{`SELECT 1 /* a FROM dual`, 9},
		{`SELECT DATE 1 FROM dual`, 7},

		// Parenthesized trailing ORDER BY blames ORDER (ORA-00907).
		{`(SELECT 1 a FROM dual ORDER BY 1)`, 22},
		{`SELECT 1 FROM dual WHERE 1 IN (SELECT 1 FROM dual ORDER BY 1)`, 50},
		{`CREATE TABLE omni_bad_ctas_R1 AS (SELECT 1 a FROM dual ORDER BY 1)`, 55},
	}

	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		if err == nil {
			t.Errorf("Parse(%q) accepted; want rejection at byte %d", tc.sql, tc.want)
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Position != tc.want {
			t.Errorf("Parse(%q) failed with %v; want ParseError at byte %d", tc.sql, err, tc.want)
		}
	}

	// Valid Oracle SQL controls beside the Wave 5 refusals.
	accepted := []string{
		`CREATE PROPERTY GRAPH g VERTEX TABLES (v1 KEY (id) PROPERTIES (p1, p2)) EDGE TABLES (e1 SOURCE KEY (src_id) REFERENCES v1 (id) DESTINATION KEY (dst_id) REFERENCES v1 (id) PROPERTIES ALL COLUMNS EXCEPT (x, y))`,
		`CREATE VECTOR INDEX vi ON t (v) INCLUDE (c1, c2) ORGANIZATION NEIGHBOR PARTITIONS`,
		`CREATE DIMENSION d LEVEL l1 IS (t.c1) LEVEL l2 IS (t.c2) HIERARCHY h (l1 CHILD OF l2)`,
		`CREATE TABLE t (a NUMBER, b VARCHAR2(10), c NUMBER GENERATED ALWAYS AS (a + 1), d NUMBER GENERATED ALWAYS AS IDENTITY, UNIQUE (a, b)) PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (10))`,
		`CREATE TABLE t (a NUMBER, a NUMBER)`,
		`CREATE TABLE t ("a" CLOB, "A" CLOB) LOB ("a", "A") STORE AS (DISABLE STORAGE IN ROW)`,
		`ALTER TABLE t ADD ("a" CLOB, "A" CLOB) LOB ("a", "A") STORE AS (DISABLE STORAGE IN ROW)`,
		`ALTER TABLE t ADD (x NUMBER, x NUMBER)`,
		`ALTER TABLE t DROP (n, n)`,
		`ALTER TABLE t ADD (c NUMBER GENERATED ALWAYS AS (a + 1))`,
		`ALTER TABLE t ENABLE UNIQUE (a, b)`,
		`ALTER TABLE t DROP PARTITION FOR (1)`,
		`CREATE VIEW v (a, b) ANNOTATIONS (ann1, ann2) AS SELECT 1, 2 FROM dual`,
		`CREATE MATERIALIZED VIEW mv (a, b) AS SELECT 1, 2 FROM dual`,
		`ALTER MATERIALIZED VIEW mv ALLOCATE EXTENT (SIZE 10M)`,
		`CREATE MATERIALIZED ZONEMAP zm ON t (a, b)`,
		`CREATE CLUSTER c (a NUMBER, b NUMBER)`,
		`CREATE INDEX idx ON t (a, b)`,
		`CREATE INDEX idx ON t (a, a)`,
		`ALTER INDEX idx UNUSABLE`,
		`CREATE OPERATOR op BINDING (NUMBER, VARCHAR2) RETURN NUMBER USING fn`,
		`ALTER INDEXTYPE it ADD op(NUMBER, VARCHAR2)`,
		`CREATE PROFILE p LIMIT SESSIONS_PER_USER 10 CPU_PER_SESSION 1000`,
		`ALTER RESOURCE COST CPU_PER_SESSION 10 LOGICAL_READS_PER_SESSION 20`,
		`CREATE PFILE = '/tmp/init.ora' FROM SPFILE`,
		`CREATE PFILE '/tmp/init.ora' FROM SPFILE '/tmp/spfile.ora'`,
		`CREATE SPFILE FROM PFILE AS COPY`,
		`SELECT * FROM EXTERNAL ((a NUMBER, b NUMBER) TYPE ORACLE_LOADER DEFAULT DIRECTORY ext_dir LOCATION ('f.csv'))`,
		`INSERT INTO t SET a = 1, b = 2`,
		`UPDATE t SET a = 1, a = 2`,
		`SELECT * FROM t FETCH FIRST 10 ROWS ONLY`,
		`SELECT * FROM t FETCH FIRST 10 PERCENT ROWS ONLY`,
		`UPDATE t PARTITION FOR (1) SET a = 2`,
		`LOCK TABLE t IN EXCLUSIVE MODE`,
		`FLASHBACK TABLE t TO SCN 100`,
		`TRUNCATE TABLE t PRESERVE MATERIALIZED VIEW LOG`,
		`ALTER DISKGROUP dg ADD VOLUME vol1 SIZE 100M`,
		`ALTER DATABASE COMMIT TO SWITCHOVER TO PHYSICAL STANDBY`,
		`ALTER PLUGGABLE DATABASE pdb1 OPEN`,
		`CREATE USER u IDENTIFIED EXTERNALLY`,
		`CREATE AUDIT POLICY pol ACTIONS COMPONENT = DV ALL`,
		`CREATE INMEMORY JOIN GROUP jg (t1(c1), t2(c2))`,
		`CREATE JAVA SOURCE IF NOT EXISTS NAMED js AS class X {}`,
		`CREATE OR REPLACE JAVA SOURCE NAMED js AS class X {}`,
		`SELECT DATE '2026-01-01', TIMESTAMP '2026-01-01 00:00:00', timestamp FROM dual`,
		`SELECT 1 "a", q'[hello]' /* comment */ FROM dual`,
		`SELECT * FROM (SELECT 1 a FROM dual ORDER BY 1)`,
		`CREATE TABLE ct AS SELECT 1 a FROM dual ORDER BY 1`,
	}

	for _, sql := range accepted {
		list, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", sql, err)
			continue
		}
		if list == nil || len(list.Items) != 1 {
			t.Errorf("Parse(%q) returned %d statements; want 1", sql, len(list.Items))
			continue
		}
		raw, ok := list.Items[0].(*ast.RawStmt)
		if !ok || raw.Stmt == nil {
			t.Errorf("Parse(%q) item is %T; want non-nil *ast.RawStmt", sql, list.Items[0])
		}
	}
}

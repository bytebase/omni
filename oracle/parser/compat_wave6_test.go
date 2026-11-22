package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestWave6CompatReadings pins the Wave 6 grammar and byte-offset families
// moved to Oracle 23.26.3's exact byte offsets, together with valid Oracle SQL
// controls beside each refusal.
func TestWave6CompatReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// CREATE CLUSTER PARTITION BY RANGE requires a hash cluster (HASHKEYS, ORA-14179 at '(').
		{`CREATE CLUSTER a(a NUMBER) PARTITION BY RANGE(a)(PARTITION VALUES LESS THAN(1))`, 45},
		{`CREATE CLUSTER a(a NUMBER) PARTITION BY RANGE()(PARTITION VALUES LESS THAN(1))`, 45},

		// CREATE TABLE ORGANIZATION INDEX requires a PRIMARY KEY constraint (ORA-25175 at INDEX).
		{`CREATE TABLE a ORGANIZATION INDEX`, 28},
		{`CREATE TABLE a(b NUMBER) ORGANIZATION INDEX`, 38},
		{`CREATE TABLE a ORGANIZATION INDEX OVERFLOW`, 28},

		// TDE encryption and integrity algorithm validation (ORA-28339).
		{`CREATE TABLE a(a ENCRYPT USING 's')`, 31},
		{`ALTER TABLE a MODIFY(a ENCRYPT USING 's')`, 37},
		{`CREATE MATERIALIZED VIEW a(a ENCRYPT USING 's') AS SELECT 1 FROM dual`, 43},
		{`CREATE MATERIALIZED VIEW a(a ENCRYPT 's') AS SELECT 1 FROM dual`, 37},
		{`CREATE TABLESPACE a ENCRYPTION USING 's' ENCRYPT`, 37},
		{`CREATE TABLESPACE SET a USING TEMPLATE(ENCRYPTION USING 's' ENCRYPT)`, 56},
		{`ALTER TABLESPACE a ENCRYPTION ONLINE USING 's'`, 43},
		{`ALTER TABLESPACE SET a ENCRYPTION ONLINE USING 's'`, 47},

		// DEALLOCATE UNUSED KEEP / ALLOCATE EXTENT SIZE rejects P and E units at integer (ORA-03277).
		{`ALTER TABLE a DEALLOCATE UNUSED KEEP 1 P`, 37},
		{`ALTER INDEX a DEALLOCATE UNUSED KEEP 1 P`, 37},
		{`ALTER CLUSTER a DEALLOCATE UNUSED KEEP 1 P`, 39},
		{`ALTER MATERIALIZED VIEW a DEALLOCATE UNUSED KEEP 1 P`, 49},
		{`ALTER MATERIALIZED VIEW LOG ON a ALLOCATE EXTENT(SIZE 1 P)`, 54},

		// CREATE CLUSTER SIZE rejects M/G/T/P/E units at integer (ORA-02229).
		{`CREATE CLUSTER a(a NUMBER) SIZE 1 M`, 32},
		{`CREATE CLUSTER a(a NUMBER) SIZE 1 G`, 32},
		{`CREATE CLUSTER a(a NUMBER) SIZE 1 T`, 32},
		{`CREATE CLUSTER a(a NUMBER) SIZE 1 P`, 32},

		// ALTER TABLESPACE SET OMF-only ADD/DROP DATAFILE rejects file specs/numbers (ORA-02580).
		{`ALTER TABLESPACE SET a ADD DATAFILE 's'`, 36},
		{`ALTER TABLESPACE SET a DROP DATAFILE 's'`, 37},
		{`ALTER TABLESPACE SET a DROP DATAFILE 1`, 37},

		// CREATE/ALTER INDEX and CLUSTER require INITRANS >= 2 (ORA-02207).
		{`CREATE INDEX a ON CLUSTER a INITRANS 1`, 38},
		{`ALTER INDEX a INITRANS 1`, 24},
		{`CREATE CLUSTER a(a NUMBER) INITRANS 1`, 37},
		{`ALTER CLUSTER a INITRANS 1`, 26},

		// XMLTYPE STORE AS OBJECT RELATIONAL requires XMLSCHEMA/ELEMENT (ORA-19002).
		{`CREATE TABLE a XMLTYPE a STORE AS OBJECT RELATIONAL`, 51},
		{`ALTER TABLE a UPGRADE INCLUDING DATA XMLTYPE a STORE AS OBJECT RELATIONAL`, 73},
		{`CREATE MATERIALIZED VIEW a XMLTYPE a STORE AS OBJECT RELATIONAL AS SELECT 1 FROM dual`, 64},

		// CREATE MATERIALIZED ZONEMAP SCALE must be in 3..20 (ORA-31951).
		{`CREATE MATERIALIZED ZONEMAP a SCALE 1 AS SELECT 1 FROM dual`, 36},
		{`CREATE MATERIALIZED ZONEMAP a SCALE 21 AS SELECT 1 FROM dual`, 36},

		// CREATE MATERIALIZED VIEW refuses PARTITION BY SYSTEM (ORA-32362).
		{`CREATE MATERIALIZED VIEW a PARTITION BY SYSTEM AS SELECT 1 FROM dual`, 47},

		// INSERT ALL INTO refuses duplicate column names (ORA-00957).
		{`INSERT ALL INTO a(a, a) SELECT 1 FROM dual`, 18},

		// suite:lexer and suite:reference exact byte-offset fixes.
		{`SELECT 1..2 FROM dual`, 9},
		{`SELECT 1 FROM t, s WHERE systimestamp(+) = s.x`, 38},
		{`CREATE TABLE omni_bad_nch_R1 (x NCHAR(2 CHAR))`, 40},
		{`CREATE TABLE omni_bad_dec_R1 (x DEC(0))`, 37},
		{`CREATE TABLE omni_bad_rlnk_R1 (n NUMBER REFERENCES omni_missing@omni_lnk (a))`, 63},
		{`(WITH w AS (SELECT 1 a FROM dual) SELECT a FROM w) UNION SELECT 1 FROM dual`, 49},
		{`SELECT 1 FROM dual MINUS ALL WITH w AS (SELECT 2 x FROM dual) SELECT x FROM w`, 62},
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

	// Valid Oracle SQL controls beside the Wave 6 refusals (including ASSOCIATE STATISTICS 3-part column).
	accepted := []string{
		`CREATE CLUSTER c (a NUMBER) HASHKEYS 10 PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (10))`,
		`CREATE TABLE t (id NUMBER PRIMARY KEY) ORGANIZATION INDEX`,
		`CREATE TABLE t (id NUMBER, CONSTRAINT pk PRIMARY KEY (id)) ORGANIZATION INDEX OVERFLOW`,
		`CREATE TABLE t (a VARCHAR2(20) ENCRYPT USING 'AES256')`,
		`ALTER TABLE t MODIFY (a ENCRYPT USING 'aes128')`,
		`CREATE MATERIALIZED VIEW mv (a ENCRYPT USING 'AES192' 'SHA-1') AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED VIEW mv (a ENCRYPT 'NOMAC') AS SELECT 1 FROM dual`,
		`CREATE TABLESPACE ts DATAFILE 'f.dbf' SIZE 10M ENCRYPTION USING 'AES256' ENCRYPT`,
		`ALTER TABLESPACE ts ENCRYPTION ONLINE USING 'AES256' ENCRYPT`,
		`ALTER TABLE t DEALLOCATE UNUSED KEEP 1 T`,
		`ALTER INDEX idx DEALLOCATE UNUSED KEEP 1 G`,
		`ALTER CLUSTER c DEALLOCATE UNUSED KEEP 1 M`,
		`ALTER MATERIALIZED VIEW LOG ON t ALLOCATE EXTENT (SIZE 1 K)`,
		`CREATE CLUSTER c (a NUMBER) SIZE 1 K`,
		`ALTER TABLESPACE SET ts_set ADD DATAFILE`,
		`ALTER TABLESPACE SET ts_set DROP DATAFILE`,
		`ALTER TABLESPACE ts ADD DATAFILE 'f.dbf' SIZE 10M`,
		`ALTER TABLESPACE ts DROP DATAFILE 'f.dbf'`,
		`CREATE INDEX idx ON CLUSTER c INITRANS 2`,
		`ALTER INDEX idx INITRANS 4`,
		`CREATE CLUSTER c (a NUMBER) INITRANS 2`,
		`ALTER CLUSTER c INITRANS 4`,
		`CREATE TABLE t (x XMLTYPE) XMLTYPE COLUMN x STORE AS OBJECT RELATIONAL XMLSCHEMA s ELEMENT e`,
		`CREATE MATERIALIZED VIEW mv XMLTYPE x STORE AS OBJECT RELATIONAL XMLSCHEMA s ELEMENT e AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED ZONEMAP zm SCALE 4 AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED VIEW mv PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (10)) AS SELECT 1 a FROM dual`,
		`INSERT ALL INTO t (a, b) SELECT 1, 2 FROM dual`,
		`ASSOCIATE STATISTICS WITH COLUMNS s.t.c USING stats_type`,
		`BEGIN FOR i IN 1..3 LOOP NULL; END LOOP; END;`,
		`CREATE TABLE t (x NCHAR(2), y DEC(10, 2), z NUMBER REFERENCES other_t (id))`,
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

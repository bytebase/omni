package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestWave7CompatReadings pins the Wave 7 grammar and byte-offset families
// moved to Oracle 23.26.3's exact byte offsets, together with valid Oracle SQL
// controls beside each refusal.
func TestWave7CompatReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// CREATE TABLE ENABLE/DISABLE PRIMARY KEY with trailing clauses requires
		// a PRIMARY KEY defined on the table (ORA-02432 at the token after PRIMARY KEY).
		{`CREATE TABLE a ENABLE PRIMARY KEY USING INDEX`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY ENABLE PRIMARY KEY`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY EXCEPTIONS INTO a`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY CASCADE`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY KEEP INDEX`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY DROP INDEX`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY USING INDEX GLOBAL PARTITION BY HASH(a) PARTITIONS 1`, 34},
		{`CREATE TABLE a ENABLE PRIMARY KEY USING INDEX LOCAL`, 34},

		// CREATE TABLE ENABLE/DISABLE CONSTRAINT name requires that constraint to be
		// defined on the table (ORA-02250).
		{`CREATE TABLE a ENABLE CONSTRAINT a`, 34},
		{`CREATE TABLE a ENABLE CONSTRAINT a zzz`, 35},
		{`CREATE TABLE a(id NUMBER) ENABLE CONSTRAINT missing_ck`, 54},

		// CREATE TABLE PARENT cannot reference the table itself (ORA-02523).
		{`CREATE TABLE a PARENT a`, 22},
		{`CREATE TABLE a PARENT a zzz`, 22},
		{`CREATE TABLE hr.a(id NUMBER) PARENT hr.a`, 36},

		// CREATE SHARDED TABLE without columns reports the table name token.
		{`CREATE SHARDED TABLE a`, 21},

		// Column COLLATE requires a collation name.
		{`ALTER TABLE a ADD(a NUMBER COLLATE)`, 34},

		// CREATE MATERIALIZED VIEW XMLTYPE STORE AS BINARY XML without XMLSCHEMA/ELEMENT (ORA-19001).
		{`CREATE MATERIALIZED VIEW a XMLTYPE a STORE AS BINARY XML AS SELECT 1 FROM dual`, 57},

		// ALTER TABLE MODIFY SUBPARTITION LOB rejects RETENTION (ORA-14160).
		{`ALTER TABLE a MODIFY SUBPARTITION a LOB(a)(RETENTION)`, 43},

		// ALTER MATERIALIZED VIEW ADD PARTITION UPDATE INDEXES rejects nameless PARTITION NOCOMPRESS (ORA-14006).
		{`ALTER MATERIALIZED VIEW a ADD PARTITION UPDATE INDEXES(a(PARTITION NOCOMPRESS))`, 40},

		// CREATE/ALTER TABLESPACE [SET] DEFAULT ILM ADD POLICY rejects ROW compression scope (ORA-38336).
		{`CREATE TABLESPACE a DEFAULT ILM ADD POLICY COMPRESS ROW AFTER 1 DAY OF CREATION`, 52},
		{`ALTER TABLESPACE a DEFAULT ILM ADD POLICY COMPRESS ROW AFTER 1 DAY OF CREATION`, 51},
		{`ALTER TABLESPACE SET a DEFAULT ILM ADD POLICY COMPRESS ROW AFTER 1 DAY OF CREATION`, 55},

		// GRANT exact byte offsets.
		{`GRANT ALL zzz ON SCHEMA a TO a`, 6},
		{`GRANT a, a ON SCHEMA a TO a`, 9},
		{`GRANT a, a zzz ON a TO a`, 6},

		// CALL INTO bare identifier reports EOF.
		{`CALL a() INTO a`, 15},

		// PL/SQL truncated block and CAST precision exact byte offsets (suite:reference & suite:manifest).
		{`BEGIN NULL;`, 10},
		{`DECLARE x NUMBER;`, 16},
		{`DECLARE v VARCHAR2(CAST(1 AS NUMBER(foo))); BEGIN NULL; END;`, 35},
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

	// Valid Oracle SQL controls beside the Wave 7 refusals.
	accepted := []string{
		`CREATE TABLE t (id NUMBER PRIMARY KEY) ENABLE PRIMARY KEY USING INDEX`,
		`CREATE TABLE t (id NUMBER, CONSTRAINT pk_t PRIMARY KEY (id)) ENABLE PRIMARY KEY USING INDEX pk_idx`,
		`CREATE TABLE t (id NUMBER, CONSTRAINT ck_t CHECK (id > 0)) DISABLE CONSTRAINT ck_t`,
		`CREATE TABLE t (id NUMBER CONSTRAINT nn_t NOT NULL) ENABLE CONSTRAINT nn_t`,
		`CREATE TABLE a ENABLE PRIMARY KEY`,
		`ALTER TABLE t ENABLE PRIMARY KEY USING INDEX TABLESPACE ts1`,
		`ALTER TABLE t ENABLE CONSTRAINT ck_t`,
		`CREATE TABLE child_tab (id NUMBER) PARENT parent_tab`,
		`CREATE TABLE child_tab (id NUMBER) PARENT hr.child_tab`,
		`CREATE SHARDED TABLE customers (id NUMBER, name VARCHAR2(200))`,
		`ALTER TABLE t ADD (name VARCHAR2(50) COLLATE BINARY_CI)`,
		`CREATE TABLE t (x XMLTYPE) XMLTYPE x STORE AS SECUREFILE BINARY XML`,
		`CREATE MATERIALIZED VIEW mv XMLTYPE x STORE AS CLOB AS SELECT 1 FROM dual`,
		`ALTER TABLE t MODIFY PARTITION p1 LOB (c) (RETENTION)`,
		`ALTER TABLE t MODIFY SUBPARTITION sp1 LOB (c) (CACHE)`,
		`ALTER MATERIALIZED VIEW mv ADD PARTITION UPDATE INDEXES (idx (PARTITION COMPRESS))`,
		`ALTER TABLE t MOVE UPDATE INDEXES (idx (PARTITION p1 NOCOMPRESS))`,
		`CREATE TABLESPACE ts DEFAULT ILM ADD POLICY COMPRESS SEGMENT AFTER 1 DAY OF CREATION`,
		`ALTER TABLESPACE ts DEFAULT ILM ADD POLICY COMPRESS GROUP AFTER 1 DAY OF CREATION`,
		`CREATE TABLE lifecycle_t (id NUMBER) ILM ADD POLICY ROW STORE COMPRESS ADVANCED ROW AFTER 30 DAYS OF NO MODIFICATION`,
		`GRANT ALL PRIVILEGES ON SCHEMA s TO u`,
		`GRANT SELECT, INSERT ON t TO u`,
		`CALL proc() INTO :res`,
		`BEGIN NULL; END;`,
		`DECLARE x NUMBER; BEGIN NULL; END;`,
		`SELECT CAST(1 AS NUMBER(10)) FROM dual`,
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

package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestBytePositionGrammarReadings pins the grammar readings that move 241
// compat.tsv rows to the byte offsets Oracle 23.26.3 reports, together with
// the accepted controls beside each refusal.
func TestBytePositionGrammarReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// CONTAINERS and SHARDS after a table name are read as its alias, so
		// Oracle stops at the '(' that follows.
		{`SELECT * FROM a CONTAINERS(a)`, 26},
		{`SELECT * FROM a SHARDS(a)`, 22},
		{`DELETE a FROM a CONTAINERS(a)`, 26},
		{`DELETE a FROM a SHARDS(a)`, 22},

		// RETURNING inside a SELECT FROM clause names a table alias, so
		// INSERT ... SELECT ... RETURNING stops at what follows RETURNING.
		{`INSERT INTO a SELECT 1 FROM dual RETURNING 1 INTO a`, 43},
		{`INSERT INTO a SELECT 1 FROM dual RETURNING INTO a`, 43},

		// DML subquery_restriction_clause enters on WITH after a FROM query.
		{`DELETE(SELECT 1 FROM dual WITH ONLY)`, 31},
		{`DELETE a FROM(SELECT 1 FROM dual WITH ONLY)`, 38},
		{`UPDATE(SELECT 1 FROM dual WITH READS ONLY) SET a = 1`, 31},
		{`UPDATE(SELECT 1 FROM dual WITH CHECKPOINT OPTION) SET a = 1`, 31},

		// CREATE TABLE object_type_col_properties and COLUMN STORE COMPRESS.
		{`CREATE TABLE a COLUMN a ELEMENT IS OF(a)`, 24},
		{`CREATE TABLE a COLUMN SUBSTITUTABLE AT ALL LEVELS`, 36},
		{`CREATE TABLE a COLUMN a AT ALL LEVELS`, 24},
		{`CREATE TABLE a COLUMN a SUBSTITUTABLE AT ALL LEVELS zzz`, 52},
		{`CREATE TABLE a COLUMN a NOT zzz SUBSTITUTABLE AT ALL LEVELS`, 28},
		{`CREATE TABLE a COLUMN COMPRESS FOR QUERY`, 22},
		{`CREATE TABLE a COLUMN STORE COMPRESS FOR`, 40},

		// MODIFY [SUB]PARTITION VARRAY reads varray_col_properties (STORE AS).
		{`ALTER TABLE a MODIFY PARTITION a VARRAY a(RETENTION)`, 41},
		{`ALTER TABLE a MODIFY SUBPARTITION a VARRAY a(RETENTION)`, 44},
		{`ALTER MATERIALIZED VIEW a MODIFY PARTITION a VARRAY a()`, 53},
		{`ALTER MATERIALIZED VIEW a MODIFY SUBPARTITION a VARRAY a()`, 56},

		// ALLOW and DISALLOW expect CLUSTERING on ALTER TABLE, its OVERFLOW,
		// and MODIFY [SUB]PARTITION, but stay refused on UPGRADE and mview OVERFLOW.
		{`ALTER TABLE a ALLOW ANYSCHEMA`, 20},
		{`ALTER TABLE a DISALLOW NONSCHEMA`, 23},
		{`ALTER TABLE a ALLOW`, 19},
		{`ALTER TABLE a OVERFLOW ALLOW EXTENT`, 29},
		{`ALTER TABLE a MODIFY PARTITION a ALLOW EXTENT`, 39},
		{`ALTER TABLE a MODIFY PARTITION a OVERFLOW ALLOW EXTENT`, 48},
		{`ALTER TABLE a MODIFY SUBPARTITION a ALLOW EXTENT`, 42},
		{`ALTER MATERIALIZED VIEW a MODIFY PARTITION a ALLOW EXTENT`, 51},
		{`ALTER MATERIALIZED VIEW a MODIFY SUBPARTITION a ALLOW EXTENT`, 54},
		{`ALTER MATERIALIZED VIEW a OVERFLOW ALLOW EXTENT`, 35},
		{`ALTER TABLE a UPGRADE INCLUDING DATA XMLTYPE a ELEMENT a ALLOW NONSCHEMA`, 57},

		// CREATE [SHARED] [PUBLIC] DATABASE LINK and AUTHENTICATED BY.
		{`CREATE SHARED DATABASE LINK a USING 's'`, 30},
		{`CREATE SHARED zzz DATABASE LINK a USING 's'`, 14},
		{`CREATE DATABASE LINK a AUTHENTICATED BY a IDENTIFIED BY a USING 's'`, 37},
		{`CREATE DATABASE LINK a AUTHENTICATED`, 36},

		// CREATE SYNONYM SHARING accepts METADATA or NONE and refuses other modes after '='.
		{`CREATE SYNONYM a SHARING = DATA FOR a`, 27},
		{`CREATE SYNONYM a SHARING METADATA FOR a`, 25},
		{`CREATE SYNONYM a SHARING =`, 26},

		// INMEMORY NO in table/mview attributes and ILM policies.
		{`CREATE TABLE a INMEMORY NO`, 26},
		{`CREATE TABLE a INMEMORY NO(a)`, 26},
		{`CREATE MATERIALIZED VIEW a INMEMORY NO AS SELECT 1 FROM dual`, 39},
		{`CREATE MATERIALIZED VIEW a INMEMORY(a) INMEMORY NO AS SELECT 1 FROM dual`, 51},
		{`ALTER TABLE a ILM ADD POLICY SET INMEMORY NO SEGMENT AFTER 1 DAY OF CREATION`, 45},
		{`ALTER TABLE a ILM ADD POLICY MODIFY INMEMORY NO SEGMENT AFTER 1 DAY OF CREATION`, 48},
		{`CREATE TABLE a ILM ADD POLICY MODIFY INMEMORY NO SEGMENT ON a`, 49},

		// CREATE/ALTER TABLE & VIEW constraint, DEFAULT COLLATION, NOFORCE and OF ... NOT.
		{`CREATE TABLE a DEFAULT a`, 23},
		{`CREATE TABLE a OF a NOT zzz SUBSTITUTABLE AT ALL LEVELS`, 24},
		{`CREATE TABLE a(PRIMARY KEY(a) zzz)`, 30},
		{`CREATE TABLE a(FOREIGN KEY(a) REFERENCES a ON CASCADE)`, 46},
		{`ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a ON CASCADE`, 49},
		{`CREATE TABLE a(a REFERENCES a ON CASCADE)`, 33},
		{`CREATE TABLE a(a WITH)`, 21},
		{`CREATE TABLE a(REF(a) WITH)`, 26},
		{`CREATE TABLE a OF a(a WITH)`, 26},
		{`CREATE TABLE a OF a(REF(a) WITH)`, 31},
		{`CREATE TABLE a(a GENERATED ALWAYS AS IDENTITY(WITH 1))`, 46},
		{`CREATE VIEW a(a, FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE) AS SELECT 1 FROM dual`, 45},
		{`ALTER VIEW a ADD FOREIGN KEY(a) REFERENCES a ON DELETE CASCADE`, 45},
		{`CREATE NOFORCE zzz VIEW a AS SELECT 1 FROM dual`, 15},
	}
	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.want {
			t.Errorf("Parse(%q): refused at %d, want %d", tc.sql, pe.Position, tc.want)
		}
	}

	// Controls: each valid form that Oracle 23.26.3 parses beside those refusals.
	controls := []string{
		`SELECT * FROM dual CONTAINERS`,
		`SELECT * FROM dual SHARDS`,
		`SELECT 1 FROM dual RETURNING`,
		`DELETE FROM t RETURNING a INTO :x`,
		`WITH FUNCTION f RETURN NUMBER IS x NUMBER; BEGIN DELETE FROM t RETURNING a INTO x; RETURN x; END; SELECT f() FROM dual`,
		`CREATE TABLE ct (a obj_t) COLUMN a SUBSTITUTABLE AT ALL LEVELS`,
		`CREATE TABLE ct (a obj_t) COLUMN a NOT SUBSTITUTABLE AT ALL LEVELS`,
		`CREATE TABLE ct (a obj_t) COLUMN a IS OF (sub_t)`,
		`CREATE TABLE ct (a obj_t) COLUMN a IS OF TYPE (ONLY sub_t)`,
		`ALTER TABLE t MODIFY PARTITION p VARRAY v STORE AS LOB l`,
		`ALTER TABLE t MODIFY SUBPARTITION sp VARRAY v STORE AS LOB l`,
		`ALTER MATERIALIZED VIEW mv MODIFY PARTITION p VARRAY v STORE AS LOB l`,
		`ALTER TABLE t ALLOW CLUSTERING LOGGING`,
		`ALTER TABLE t DISALLOW CLUSTERING LOGGING`,
		`ALTER TABLE t OVERFLOW ALLOW CLUSTERING`,
		`ALTER TABLE t MODIFY PARTITION p ALLOW CLUSTERING`,
		`ALTER TABLE t MODIFY SUBPARTITION sp DISALLOW CLUSTERING`,
		`CREATE SHARED DATABASE LINK d AUTHENTICATED BY u IDENTIFIED BY p USING 's'`,
		`CREATE SHARED PUBLIC DATABASE LINK d CONNECT TO u1 IDENTIFIED BY p1 AUTHENTICATED BY u2 IDENTIFIED BY p2 USING 's'`,
		`CREATE SYNONYM syn SHARING = NONE FOR t`,
		`CREATE VIEW v (a CONSTRAINTS c UNIQUE RELY DISABLE NOVALIDATE) AS SELECT 1 FROM dual`,
		`CREATE VIEW v (a, CONSTRAINTS c UNIQUE(a) RELY DISABLE NOVALIDATE) AS SELECT 1 FROM dual`,
	}
	for _, sql := range controls {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}

	// Verify CreateDatabaseLinkStmt records SHARED and AUTHENTICATED BY fields.
	list, err := Parse(`CREATE SHARED PUBLIC DATABASE LINK db1 CONNECT TO u1 IDENTIFIED BY p1 AUTHENTICATED BY u2 IDENTIFIED BY p2 USING 'svc'`)
	if err != nil {
		t.Fatalf("Parse shared dblink: %v", err)
	}
	raw := list.Items[0].(*ast.RawStmt)
	dbl, ok := raw.Stmt.(*ast.CreateDatabaseLinkStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateDatabaseLinkStmt, got %T", raw.Stmt)
	}
	if !dbl.Shared || !dbl.Public || dbl.AuthenticatedUser != "U2" || dbl.AuthenticatedPass != "P2" {
		t.Fatalf("unexpected CreateDatabaseLinkStmt fields: %+v", dbl)
	}
	rendered := ast.NodeToString(dbl)
	if !strings.Contains(rendered, ":shared true") || !strings.Contains(rendered, `:authenticatedUser "U2"`) {
		t.Fatalf("unexpected NodeToString output: %s", rendered)
	}
}

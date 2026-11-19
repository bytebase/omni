package parser

import (
	"errors"
	"testing"
)

// TestOverAcceptsRefused holds the rows of compat.tsv that omni accepted and
// the pinned Oracle refused, where Oracle's answer is the grammar's and not a
// check of the parsed statement: each statement is refused at the byte Oracle
// reports, and the form beside it, which Oracle takes, still parses. The
// measured run is oracle/parser/testdata/README.md's "Rows only omni accepts,
// measured".
func TestOverAcceptsRefused(t *testing.T) {
	// refused: the statement, and the byte Oracle's DBMS_SQL.PARSE reports.
	refused := []struct {
		sql string
		at  int
	}{
		// ANALYZE names TABLE, INDEX or CLUSTER, and Oracle answers ORA-01490
		// at the word the action starts with.
		{"ANALYZE a VALIDATE STRUCTURE", 10},
		{"ANALYZE TABLES a VALIDATE STRUCTURE", 15},
		{"ANALYZE INDEXES a VALIDATE STRUCTURE", 16},
		{"ANALYZE CLUSTERING a VALIDATE STRUCTURE", 19},
		// CREATE SCHEMA requires AUTHORIZATION (ORA-02420).
		{"CREATE SCHEMA a CREATE TABLE a (a NUMBER)", 14},
		{"SCHEMA AUTHORIZATION a CREATE TABLE a (a NUMBER)", 0},
		// DROP USER requires a user (ORA-01935).
		{"DROP USER", 9},
		// DROP AUDIT POLICY requires a policy (ORA-46353).
		{"DROP AUDIT POLICY", 17},
		// ALTER TRIGGER RENAME TO requires a new name (ORA-04070).
		{"ALTER TRIGGER a RENAME TO", 25},
		// PURGE takes TABLE, INDEX, TABLESPACE, RECYCLEBIN or DBA_RECYCLEBIN,
		// and nothing else: ORA-38302 at the word.
		{"PURGE a", 6},
		{"PURGE TABLES a", 6},
		{"PURGE INDEXES a", 6},
		{"PURGE TAG a", 6},
		{"PURGE TAG SET a", 6},
		{"PURGE SET a", 6},
		{"PURGE REDO", 6},
		{"PURGE DBTIMEZONE", 6},
		// SET PARTITIONING takes AUTOMATIC or MANUAL (ORA-14853).
		{"ALTER TABLE a SET PARTITIONING", 30},
		{"ALTER MATERIALIZED VIEW a SET PARTITIONING", 42},
		// IF NOT EXISTS is the whole clause (ORA-11543 at the word where
		// EXISTS belongs).
		{"CREATE FUNCTION IF NOT a RETURN NUMBER IS BEGIN RETURN 1; END;", 23},
		// A materialized view's inmemory_column_clause writes its columns
		// first and then INMEMORY or NO INMEMORY (ORA-00922 where the query
		// starts).
		{"CREATE MATERIALIZED VIEW a INMEMORY(a) AS SELECT 1 FROM dual", 39},
		{"CREATE MATERIALIZED VIEW a INMEMORY(a) INMEMORY(a) INMEMORY AS SELECT 1 FROM dual", 47},
		// Oracle's SQL lexer reads ~= as ~ then = (ORA-00911 at the ~).
		{"SELECT 1 FROM dual WHERE 1 ~= 2", 27},
		{"SELECT CASE WHEN 'a' ~= 'b' THEN 1 ELSE 0 END FROM dual", 21},
		// A quoted identifier holds no double quotation mark, doubled or not
		// (ORA-25716 at the identifier).
		{`SELECT 1 "a""b" FROM dual`, 9},
		{`CREATE TABLE "a""b" (x NUMBER)`, 13},
	}
	for _, tc := range refused {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.at {
			t.Errorf("Parse(%q): refused at %d (%s), want %d", tc.sql, pe.Position, pe.Message, tc.at)
		}
	}

	// accepted: the forms beside them, which Oracle takes.
	accepted := []string{
		"ANALYZE TABLE a VALIDATE STRUCTURE",
		"ANALYZE INDEX a VALIDATE STRUCTURE",
		"ANALYZE CLUSTER a VALIDATE STRUCTURE",
		"ANALYZE TABLE a COMPUTE STATISTICS",
		"ANALYZE TABLE a PARTITION (p) VALIDATE STRUCTURE",
		"CREATE SCHEMA AUTHORIZATION a CREATE TABLE a (a NUMBER)",
		"CREATE SCHEMA AUTHORIZATION a",
		"DROP USER a CASCADE",
		"DROP USER IF EXISTS a",
		"DROP AUDIT POLICY a",
		"ALTER TRIGGER a RENAME TO b",
		"PURGE TABLE a",
		"PURGE INDEX a",
		"PURGE TABLESPACE a",
		"PURGE RECYCLEBIN",
		"PURGE DBA_RECYCLEBIN",
		"ALTER TABLE a SET PARTITIONING AUTOMATIC",
		"ALTER TABLE a SET PARTITIONING MANUAL",
		"ALTER MATERIALIZED VIEW a SET PARTITIONING AUTOMATIC",
		"CREATE FUNCTION IF NOT EXISTS a RETURN NUMBER IS BEGIN RETURN 1; END;",
		"CREATE FUNCTION a RETURN NUMBER IS BEGIN RETURN 1; END;",
		// A materialized view's column clause, as its page writes it.
		"CREATE MATERIALIZED VIEW a INMEMORY(a) INMEMORY AS SELECT 1 FROM dual",
		"CREATE MATERIALIZED VIEW a INMEMORY(a) NO INMEMORY AS SELECT 1 FROM dual",
		"CREATE MATERIALIZED VIEW a INMEMORY NO MEMCOMPRESS AS SELECT 1 FROM dual",
		"CREATE MATERIALIZED VIEW a NO INMEMORY AS SELECT 1 FROM dual",
		// A table's, which writes the keyword first, and a column list alone.
		"CREATE TABLE a(a NUMBER) INMEMORY(a)",
		"CREATE TABLE a(a NUMBER) INMEMORY(a) NO INMEMORY",
		// PL/SQL's ~=, which its lexer reads whole.
		"BEGIN IF 1 ~= 2 THEN NULL; END IF; END;",
		"DECLARE x NUMBER := 1; BEGIN WHILE x ~= 2 LOOP NULL; END LOOP; END;",
		`SELECT 1 "a b" FROM dual`,
	}
	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want it taken", sql, err)
		}
	}
}

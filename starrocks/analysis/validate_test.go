package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/starrocks/parser"
)

// A statement cut out of a script reports its position in the script, not in
// the segment: a caller highlighting the error must land on the right text.
func TestValidateQueryReportsScriptPositions(t *testing.T) {
	const script = "SELECT 1;\nSELECT FROM"
	_, _, err := ValidateQuery(script)
	if err == nil {
		t.Fatal("ValidateQuery accepted a syntax error")
	}
	errs := parser.AllErrors(err)
	if len(errs) == 0 {
		t.Fatalf("ValidateQuery error = %v, want a ParseError", err)
	}
	if want := strings.Index(script, "FROM"); errs[0].Position != want {
		t.Errorf("error position = %d, want %d", errs[0].Position, want)
	}
}

// A statement that changes state outside the session is not read-only, and the
// session forms stay allowed.
func TestValidateQueryServerState(t *testing.T) {
	tests := []struct {
		statement   string
		readOnly    bool
		returnsData bool
	}{
		{"SET query_timeout = 3600", true, false},
		{"SET @@session.query_timeout = 3600", true, false},
		{"SET GLOBAL query_timeout = 3600", false, false},
		{"SET @@global.query_timeout = 3600", false, false},
		{"SELECT 1", true, true},
		{"DESC FILES('path' = 's3://x')", true, true},
		{"DESC PIPE db.p", true, true},
		{"ADMIN SHOW FRONTEND CONFIG", false, true},
	}
	for _, tc := range tests {
		readOnly, returnsData, err := ValidateQuery(tc.statement)
		if err != nil {
			t.Errorf("%q: %v", tc.statement, err)
			continue
		}
		if readOnly != tc.readOnly || returnsData != tc.returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (%v, %v)", tc.statement, readOnly, returnsData, tc.readOnly, tc.returnsData)
		}
	}
}

// CACHE SELECT returns the cache warm-up metrics, so the driver reads rows.
func TestValidateQueryDataCacheSelect(t *testing.T) {
	readOnly, returnsData, err := ValidateQuery("CACHE SELECT a FROM t")
	if err != nil {
		t.Fatalf("ValidateQuery: %v", err)
	}
	if readOnly || !returnsData {
		t.Errorf("CACHE SELECT = (readOnly %v, returnsData %v), want (false, true)", readOnly, returnsData)
	}
}

// SUBMIT TASK returns the task name and its status, so the driver must read
// rows even though the statement is not read-only.
func TestValidateQuerySubmitTask(t *testing.T) {
	readOnly, returnsData, err := ValidateQuery("SUBMIT TASK AS INSERT INTO t SELECT 1")
	if err != nil {
		t.Fatalf("ValidateQuery: %v", err)
	}
	if readOnly || !returnsData {
		t.Errorf("SUBMIT TASK = (readOnly %v, returnsData %v), want (false, true)", readOnly, returnsData)
	}
}

// EXPLAIN and DESC reject queries that carry INTO OUTFILE: the editor refuses
// by way of EXPLAIN/DESC what it refuses when asked directly. EXPLAIN ANALYZE
// still returns its profile result set when run outside read-only mode.
func TestValidateQueryExplainDescOutfile(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		returnsData bool
	}{
		{`EXPLAIN SELECT a FROM t INTO OUTFILE "s3://b/x"`, false},
		{`EXPLAIN VERBOSE SELECT a FROM t INTO OUTFILE "s3://b/x"`, false},
		{`EXPLAIN (SELECT a FROM t) INTO OUTFILE "s3://b/x"`, false},
		{`DESC SELECT 1 INTO OUTFILE 'x'`, false},
		{`DESCRIBE SELECT a FROM t INTO OUTFILE "s3://b/x"`, false},
		{`DESC VERBOSE SELECT a FROM t INTO OUTFILE "s3://b/x"`, false},
		{`EXPLAIN ANALYZE SELECT a FROM t INTO OUTFILE "s3://b/x"`, true},
		{`EXPLAIN ANALYZE SELECT 1 UNION SELECT a FROM t INTO OUTFILE "s3://b/x"`, true},
		{`EXPLAIN ANALYZE (SELECT 1) INTO OUTFILE "s3://b/x"`, true},
	} {
		readOnly, returnsData, err := ValidateQuery(tc.sql)
		if err != nil {
			t.Errorf("%q: %v", tc.sql, err)
			continue
		}
		if readOnly || returnsData != tc.returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (false, %v)", tc.sql, readOnly, returnsData, tc.returnsData)
		}
	}
}

// EXPLAIN ANALYZE over DML is not read-only, but its profile is still a
// result set the driver must read rather than Exec.
func TestValidateQueryAnalyzeInsert(t *testing.T) {
	for _, sql := range []string{
		"EXPLAIN ANALYZE INSERT INTO t SELECT 1",
		"DESC ANALYZE INSERT INTO t SELECT 1",
	} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Errorf("%q: %v", sql, err)
			continue
		}
		if readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (false, true)", sql, readOnly, returnsData)
		}
	}
}

// TRACE plans the query it wraps and never runs it, so it is read-only like
// EXPLAIN, including over a query that writes out a file.
func TestValidateQueryOptimizerTrace(t *testing.T) {
	for _, sql := range []string{
		"TRACE ALL SELECT 1",
		"TRACE LOGS SELECT * FROM t",
		"TRACE ALL SELECT 1 INTO OUTFILE 's3://bucket/x'",
	} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Errorf("%q: %v", sql, err)
			continue
		}
		if !readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (true, true)", sql, readOnly, returnsData)
		}
	}
}

// DESC STORAGE VOLUME answers with a storage volume's settings, so the
// editor's read-only mode admits it the way it admits DESC over a table.
func TestValidateQueryDescStorageVolume(t *testing.T) {
	for _, sql := range []string{"DESC STORAGE VOLUME v", "DESCRIBE STORAGE VOLUME v"} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Errorf("%q: %v", sql, err)
			continue
		}
		if !readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (true, true)", sql, readOnly, returnsData)
		}
	}
}

// An explained statement is a plan, and only StarRocks's EXPLAIN ANALYZE runs
// the statement it plans.
func TestDescribeExplain(t *testing.T) {
	for _, statement := range []string{"SELECT 1", "SHOW TABLES"} {
		if _, _, ok := DescribeExplain(statement); ok {
			t.Errorf("%q describes a plan", statement)
		}
	}
	format, executed, ok := DescribeExplain("EXPLAIN SELECT 1")
	if !ok || format != "TEXT" || executed {
		t.Errorf("EXPLAIN = (%q, %v, %v), want (TEXT, false, true)", format, executed, ok)
	}
	format, executed, ok = DescribeExplain("EXPLAIN VERBOSE SELECT 1")
	if !ok || format != "VERBOSE" || executed {
		t.Errorf("EXPLAIN VERBOSE = (%q, %v, %v), want (VERBOSE, false, true)", format, executed, ok)
	}
	format, executed, ok = DescribeExplain("TRACE ALL SELECT 1")
	if !ok || format != "TRACE" || executed {
		t.Errorf("TRACE = (%q, %v, %v), want (TRACE, false, true)", format, executed, ok)
	}
}

// A buffer that cannot lex is not read-only anything: the chunk it sits in
// holds no statement, so parser.Split drops it and a per-segment parse never
// sees it. The check runs on the whole text before splitting.
func TestValidateQueryFailsClosedOnLexError(t *testing.T) {
	for _, script := range []string{"/*", "SELECT 1; /*"} {
		if readOnly, returnsData, err := ValidateQuery(script); err == nil {
			t.Errorf("ValidateQuery(%q) = %v, %v, nil; want the lex error", script, readOnly, returnsData)
		}
	}
}

// A segment carries its terminator into the parse: TRANSLATE's translateSQL
// is `.+`, so `TRANSLATE presto;` reads but a bare `TRANSLATE presto` at the
// end of the input does not.
func TestValidateQueryKeepsTheTerminator(t *testing.T) {
	if _, _, err := ValidateQuery("TRANSLATE presto;"); err != nil {
		t.Errorf("ValidateQuery(TRANSLATE presto;): %v", err)
	}
	if _, _, err := ValidateQuery("TRANSLATE presto"); err == nil {
		t.Error("ValidateQuery(TRANSLATE presto) accepted, want a syntax error")
	}
}

// EXECUTE may run a prepared SELECT, so the driver reads rows from it rather
// than Exec'ing them away.
func TestValidateQueryExecuteReturnsData(t *testing.T) {
	readOnly, returnsData, err := ValidateQuery("EXECUTE p")
	if err != nil {
		t.Fatalf("ValidateQuery(EXECUTE p): %v", err)
	}
	if readOnly || !returnsData {
		t.Errorf("EXECUTE = (readOnly %v, returnsData %v), want (false, true)", readOnly, returnsData)
	}
}

// TRANSLATE answers with the translated SQL instead of running it, so the
// driver reads rows from it.
func TestValidateQueryTranslateReturnsData(t *testing.T) {
	for _, sql := range []string{"TRANSLATE trino SELECT 1", "TRANSLATE trino SELECT 1;"} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Fatalf("ValidateQuery(%q): %v", sql, err)
		}
		if readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (false, true)", sql, readOnly, returnsData)
		}
	}
}

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		statement   string
		valid       bool
		description string
	}{
		{statement: "SELECT * FROM users", valid: true, description: "Basic SELECT query"},
		{statement: "SELECT * FROM users WHERE id = 1", valid: true, description: "SELECT with WHERE clause"},
		{statement: "SHOW DATA", valid: true, description: "SHOW DATA statement"},
		{statement: "SHOW DATA FROM db1", valid: true, description: "SHOW DATA with FROM clause"},
		{statement: "SHOW DATABASES", valid: true, description: "SHOW DATABASES statement"},
		{statement: "SHOW TABLES", valid: true, description: "SHOW TABLES statement"},
		{statement: "SHOW TABLETS FROM table1", valid: true, description: "SHOW TABLETS statement"},
		{statement: "SHOW VARIABLES", valid: true, description: "SHOW VARIABLES statement"},
		{statement: "SHOW CREATE TABLE users", valid: true, description: "SHOW CREATE TABLE statement"},
		{statement: "SHOW CREATE DATABASE db1", valid: true, description: "SHOW CREATE DATABASE statement"},
		{statement: "INSERT INTO users (id, name) VALUES (1, 'test')", valid: false, description: "INSERT statement should be invalid"},
		{statement: "UPDATE users SET name = 'test' WHERE id = 1", valid: false, description: "UPDATE statement should be invalid"},
		{statement: "DELETE FROM users WHERE id = 1", valid: false, description: "DELETE statement should be invalid"},
		{statement: "CREATE TABLE test (id INT)", valid: false, description: "CREATE TABLE should be invalid"},
		{statement: "DROP TABLE users", valid: false, description: "DROP TABLE should be invalid"},
		{statement: "EXPLAIN INSERT INTO users (id, name) VALUES (1, 'test')", valid: true, description: "EXPLAIN INSERT should be valid (read-only)"},
		{statement: "EXPLAIN UPDATE users SET name = 'test' WHERE id = 1", valid: true, description: "EXPLAIN UPDATE should be valid (read-only)"},
		{statement: "EXPLAIN DELETE FROM users WHERE id = 1", valid: true, description: "EXPLAIN DELETE should be valid (read-only)"},
		{statement: "WITH c AS (SELECT 1) SELECT * FROM c", valid: true, description: "CTE-prefixed SELECT is read-only"},
		{statement: "((SELECT 1))", valid: true, description: "Nested parenthesized SELECT is read-only"},
		{statement: "(SELECT 1) LIMIT 5", valid: true, description: "Parenthesized SELECT with trailing LIMIT is read-only"},
		{statement: "(WITH c AS (SELECT 1) SELECT * FROM c)", valid: true, description: "Parenthesized WITH query is read-only"},
		{statement: `SELECT 1 UNION (SELECT a FROM t INTO OUTFILE "s3://b/x")`, valid: false, description: "INTO OUTFILE inside a parenthesized set-op arm must not slip the export gate"},
		{statement: `((SELECT a FROM t INTO OUTFILE "s3://b/x"))`, valid: false, description: "INTO OUTFILE inside nested parens must not slip the export gate"},
		{statement: "EXPLAIN (SELECT 1)", valid: true, description: "EXPLAIN over a parenthesized query is read-only"},
		{statement: "EXPLAIN ((SELECT 1))", valid: true, description: "EXPLAIN over nested parens is read-only"},
		{statement: "EXPLAIN ANALYZE SELECT a FROM t", valid: true, description: "EXPLAIN ANALYZE over a SELECT is read-only"},
		{statement: "EXPLAIN ANALYZE (SELECT 1)", valid: true, description: "EXPLAIN ANALYZE over a parenthesized query is read-only"},
		{statement: "EXPLAIN ANALYZE ((SELECT 1))", valid: true, description: "EXPLAIN ANALYZE over nested parens is read-only"},
		{statement: "EXPLAIN ANALYZE SELECT 1 UNION SELECT 2", valid: true, description: "EXPLAIN ANALYZE over a set operation is read-only"},
		{statement: "EXPLAIN VERBOSE INSERT INTO users (id, name) VALUES (1, 'test')", valid: true, description: "a plan-only modifier stays read-only over DML"},
		{statement: "EXPLAIN ANALYZE INSERT INTO users (id, name) VALUES (1, 'test')", valid: false, description: "EXPLAIN ANALYZE over INSERT must be rejected"},
		{statement: "EXPLAIN ANALYZE INSERT INTO users SELECT id, name FROM other", valid: false, description: "EXPLAIN ANALYZE over INSERT ... SELECT must be rejected"},
		{statement: "SELECT EXTRACT(YEAR FROM MONTHS_ADD(NOW(), -1))", valid: true, description: "EXTRACT(unit FROM expr) is read-only"},
		{statement: "SELECT EXTRACT(\n  YEAR\n  FROM\n    MONTHS_ADD(NOW(), -1)\n)", valid: true, description: "Multi-line EXTRACT is read-only"},
		{statement: "WITH c AS (SELECT EXTRACT(MONTH FROM created_at) m FROM t) SELECT * FROM c", valid: true, description: "EXTRACT inside a CTE body is read-only"},
		{statement: "SELECT TRIM(s), LEFT(s, 2), RIGHT(s, 2), IF(a, 1, 2) FROM t", valid: true, description: "Reserved keywords used as function names are read-only"},
		{statement: "SELECT DATABASE()", valid: true, description: "DATABASE() is read-only"},
		{statement: "SELECT CURRENT_TIMESTAMP(3)", valid: true, description: "Datetime function with precision is read-only"},
		{statement: "SELECT @@version_comment, @@session.query_timeout, @uservar", valid: true, description: "Session and user variables are read-only"},
		{statement: "SELECT [1, 2, 3], map{'a': 1}, x'1f', b'101'", valid: true, description: "Array, map, hex and bit literals are read-only"},
		{statement: "SELECT array_map(x -> x + 1, arr) FROM t", valid: true, description: "Lambda in a higher-order array function is read-only"},
		{statement: "SELECT /*+ SET_VAR(query_timeout = 100) */ a FROM t", valid: true, description: "Statement with an optimizer hint is read-only"},
		{statement: "SELECT a, SUM(b) FROM t GROUP BY CUBE(a)", valid: true, description: "GROUP BY CUBE is read-only"},
		{statement: "SELECT CONVERT(s, SIGNED), BINARY s FROM t", valid: true, description: "CONVERT and BINARY forms are read-only"},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			valid, _, err := ValidateQuery(tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.valid, valid, "statement: %s", tc.statement)
		})
	}

	rejectCases := []struct {
		statement   string
		description string
	}{
		{statement: "WITH c AS (SELECT 1) UPDATE t SET x = 1", description: "CTE-prefixed UPDATE must be rejected"},
		{statement: "WITH c AS (SELECT 1) DELETE FROM t", description: "CTE-prefixed DELETE must be rejected"},
		{statement: "WITH c AS (SELECT 1) INSERT INTO t SELECT * FROM c", description: "CTE-prefixed INSERT must be rejected"},
		{statement: "SELECT a > (select max(a) from t1) FROM", description: "Truncated SELECT must be rejected as syntax error"},
		{statement: "SHOW", description: "Bare SHOW must be rejected"},
		{statement: "DESCRIBE", description: "Bare DESCRIBE must be rejected"},
		{statement: "EXPLAIN", description: "Bare EXPLAIN must be rejected"},
		{statement: "EXPLAIN DROP TABLE t", description: "EXPLAIN over DDL must be rejected"},
		{statement: "SET DEFAULT STORAGE VAULT v", description: "SET DEFAULT STORAGE VAULT must be rejected"},
		{statement: "EXPLAIN ANALYZE UPDATE users SET name = 'test' WHERE id = 1", description: "EXPLAIN ANALYZE over UPDATE must be rejected"},
		{statement: "EXPLAIN ANALYZE DELETE FROM users WHERE id = 1", description: "EXPLAIN ANALYZE over DELETE must be rejected"},
		{statement: "EXPLAIN ANALYZE TRUNCATE TABLE users", description: "EXPLAIN ANALYZE over TRUNCATE must be rejected"},
		{statement: `EXPLAIN ANALYZE (SELECT a FROM t INTO OUTFILE "s3://b/x")`, description: "EXPLAIN ANALYZE over a parenthesized export must be rejected"},
		{statement: `EXPLAIN ANALYZE ((SELECT a FROM t INTO OUTFILE "s3://b/x"))`, description: "EXPLAIN ANALYZE over a nested parenthesized export must be rejected"},
		{statement: "SELECT a FROM t xx yy zz", description: "Trailing junk after a valid prefix must be rejected"},
	}
	for _, tc := range rejectCases {
		t.Run(tc.description, func(t *testing.T) {
			valid, _, err := ValidateQuery(tc.statement)
			require.False(t, valid && err == nil, "expected rejection, statement: %s", tc.statement)
		})
	}
}

func TestValidateQuerySetStatements(t *testing.T) {
	for _, tc := range []struct {
		name             string
		statement        string
		wantValid        bool
		wantExplicitRead bool
	}{
		{"a session variable", "SET time_zone = 'UTC'", true, false},
		{"the @@ spelling", "SET @@time_zone = 'UTC'", true, false},
		{"an explicit session scope", "SET SESSION time_zone = 'UTC'", true, false},
		{"a global scope", "SET GLOBAL time_zone = 'UTC'", false, false},
		{"SET NAMES", "SET NAMES utf8mb4", true, false},
		{"a read stays explicitly read-only", "SELECT 1", true, true},
		{"a SET anywhere in a script needs executing", "SELECT 1; SET time_zone = 'UTC'", true, false},
		{"a SET does not launder a DDL beside it", "SET time_zone = 'UTC'; DROP TABLE t", false, false},
		{"SET PASSWORD administers the server", "SET PASSWORD FOR 'u' = 'x'", false, false},
		{"SET DEFAULT STORAGE VOLUME administers the cluster", "SET v AS DEFAULT STORAGE VOLUME", false, false},
		{"SET CATALOG switches the session", "SET CATALOG hive", true, false},
		{"SET WAREHOUSE switches the session", "SET WAREHOUSE w", true, false},
		{"SET ROLE activates roles for the session", "SET ROLE r1, r2", true, false},
		{"SET ROLE ALL", "SET ROLE ALL EXCEPT r1", true, false},
		{"SET DEFAULT ROLE changes a user's default roles", "SET DEFAULT ROLE r TO u", false, false},
		{"SET DEFAULT ROLE ALL for several users", "SET DEFAULT ROLE ALL TO 'u1'@'%', u2", false, false},
		{"SET LDAP_ADMIN_PASSWORD stores the LDAP bind password", "SET LDAP_ADMIN_PASSWORD = 'x'", false, false},
		{"SET LDAP_ADMIN_PASSWORD in lowercase with PASSWORD()", "set ldap_admin_password = password('x')", false, false},
		{"SET LDAP_ADMIN_PASSWORD beside a session variable", "SET time_zone = 'UTC', LDAP_ADMIN_PASSWORD = 'x'", false, false},
		{"a user variable named ldap_admin_password", "SET @ldap_admin_password = 'x'", true, false},
		{"SET TRANSACTION", "SET TRANSACTION READ ONLY", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid, explicitReadOnly, err := ValidateQuery(tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.wantValid, valid)
			require.Equal(t, tc.wantExplicitRead, explicitReadOnly)
		})
	}
}

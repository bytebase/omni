package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/doris/parser"
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
		{"UNSET VARIABLE ALL", true, false},
		{"UNSET GLOBAL VARIABLE ALL", false, false},
		{"SELECT 1", true, true},
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
		{statement: "SELECT EXTRACT(YEAR FROM MONTHS_ADD(NOW(), -1))", valid: true, description: "EXTRACT(unit FROM expr) is read-only"},
		{statement: "SELECT EXTRACT(\n  YEAR\n  FROM\n    MONTHS_ADD(NOW(), -1)\n)", valid: true, description: "Multi-line EXTRACT is read-only"},
		{statement: "WITH c AS (SELECT EXTRACT(MONTH FROM created_at) m FROM t) SELECT * FROM c", valid: true, description: "EXTRACT inside a CTE body is read-only"},
		{statement: "SELECT TRIM(s), LEFT(s, 2), RIGHT(s, 2), IF(a, 1, 2) FROM t", valid: true, description: "Reserved keywords used as function names are read-only"},
		{statement: "SELECT DATABASE()", valid: true, description: "DATABASE() is read-only"},
		{statement: "SELECT CURRENT_TIMESTAMP(3)", valid: true, description: "Datetime function with precision is read-only"},
		{statement: "SELECT @@version_comment, @@session.query_timeout, @uservar", valid: true, description: "Session and user variables are read-only"},
		{statement: "SELECT [1, 2, 3], {'a': 1}, x'1f', b'101'", valid: true, description: "Array, map, hex and bit literals are read-only"},
		{statement: "SELECT array_map(x -> x + 1, arr) FROM t", valid: true, description: "Lambda in a higher-order array function is read-only"},
		{statement: "SELECT /*+ SET_VAR(query_timeout = 100) */ a FROM t", valid: true, description: "Statement with an optimizer hint is read-only"},
		{statement: "SELECT a, SUM(b) FROM t GROUP BY CUBE(a)", valid: true, description: "GROUP BY CUBE is read-only"},
		{statement: "SELECT CONVERT(s USING utf8), CONVERT(s, SIGNED), BINARY s FROM t", valid: true, description: "CONVERT and BINARY forms are read-only"},
		{statement: "(SELECT 1) LIMIT 1", valid: true, description: "Parenthesized query with trailing LIMIT is read-only"},
		{statement: "(SELECT 1 LIMIT 1) LIMIT 2", valid: true, description: "GroupedQuery wrapper (repeated clause groups) is read-only"},
		{statement: "SELECT 1 LIMIT 5 ORDER BY 1", valid: true, description: "Engine-lenient clause order wraps in GroupedQuery, still read-only"},
		{statement: "(WITH c AS (SELECT 1) SELECT * FROM c)", valid: true, description: "Parenthesized WITH group (CTE scope boundary) is read-only"},
		{statement: "SELECT 1 UNION (SELECT 2) LIMIT 5", valid: true, description: "Set operation with parenthesized operand and trailing LIMIT is read-only"},
		{statement: "EXPLAIN (SELECT 1)", valid: true, description: "EXPLAIN over a parenthesized query is read-only"},
		{statement: "EXPLAIN (SELECT 1 LIMIT 1) LIMIT 2", valid: true, description: "EXPLAIN over a grouped query is read-only"},
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
		{"UNSET VARIABLE", "UNSET VARIABLE time_zone", true, false},
		{"a read stays explicitly read-only", "SELECT 1", true, true},
		{"a SET anywhere in a script needs executing", "SELECT 1; SET time_zone = 'UTC'", true, false},
		{"a SET does not launder a DDL beside it", "SET time_zone = 'UTC'; DROP TABLE t", false, false},
		{"SET PASSWORD administers the server", "SET PASSWORD FOR 'u' = 'x'", false, false},
		{"SET DEFAULT STORAGE VAULT administers the server", "SET v AS DEFAULT STORAGE VAULT", false, false},
		{"SET PROPERTY sets a user's properties", "SET PROPERTY FOR 'u' 'max_user_connections' = '1000'", false, false},
		{"SET PROPERTY for the current user", "SET PROPERTY 'max_user_connections' = '1000'", false, false},
		{"a variable named property", "SET PROPERTY = 1", true, false},
		{"SET LDAP_ADMIN_PASSWORD stores the LDAP bind password", "SET LDAP_ADMIN_PASSWORD = 'x'", false, false},
		{"SET LDAP_ADMIN_PASSWORD with PASSWORD()", "SET LDAP_ADMIN_PASSWORD = PASSWORD('x')", false, false},
		{"SET LDAP_ADMIN_PASSWORD beside a session variable", "SET time_zone = 'UTC', LDAP_ADMIN_PASSWORD = 'x'", false, false},
		{"a user variable named ldap_admin_password", "SET @ldap_admin_password = 'x'", true, false},
		{"SET TRANSACTION", "SET TRANSACTION READ ONLY", true, false},
		{"a user variable", "SET @a = 1", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid, explicitReadOnly, err := ValidateQuery(tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.wantValid, valid)
			require.Equal(t, tc.wantExplicitRead, explicitReadOnly)
		})
	}
}

func TestValidateQueryFileExport(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		wantValid bool
	}{
		{"plain", "SELECT 1 INTO OUTFILE '/tmp/x'", false},
		{"s3 with format and properties", `SELECT * FROM t INTO OUTFILE "s3://bucket/export/" FORMAT AS PARQUET PROPERTIES("s3.endpoint" = "example.com")`, false},
		{"after ORDER BY and LIMIT", "SELECT * FROM t ORDER BY a LIMIT 10 INTO OUTFILE '/tmp/x'", false},
		{"on a set operation", "SELECT a FROM t1 UNION SELECT a FROM t2 INTO OUTFILE '/tmp/x'", false},
		{"after parentheses", "(SELECT * FROM t) INTO OUTFILE '/tmp/x'", false},
		{"on a grouped query", "(SELECT 1 ORDER BY 1) ORDER BY 1 INTO OUTFILE '/tmp/x'", false},
		{"with a CTE", "WITH c AS (SELECT 1 AS a) SELECT * FROM c INTO OUTFILE '/tmp/x'", false},
		{"under EXPLAIN", "EXPLAIN SELECT 1 INTO OUTFILE '/tmp/x'", false},
		{"beside a read in one script", "SELECT 1; SELECT 2 INTO OUTFILE '/tmp/x'", false},
		{"in a string literal", `SELECT a FROM t WHERE a = 'into outfile'`, true},
		{"in a quoted identifier", "SELECT `INTO OUTFILE` FROM t", true},
		{"in a block comment", "SELECT * FROM t /* INTO OUTFILE */", true},
		{"in a line comment", "SELECT * FROM t -- INTO OUTFILE 's3://bucket/'", true},
		{"a plain read", "SELECT * FROM t", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid, _, err := ValidateQuery(tc.statement)
			require.NoError(t, err)
			require.Equal(t, tc.wantValid, valid)
		})
	}

	for _, statement := range []string{
		"SELECT 1 INTO DUMPFILE '/tmp/x'",
		"SELECT 1 INTO @v",
		"SELECT * FROM (SELECT a FROM t INTO OUTFILE '/tmp/x') x",
		"SELECT 1 UNION (SELECT a FROM t INTO OUTFILE '/tmp/x')",
		"SELECT 'unterminated",
	} {
		t.Run("syntax error: "+statement, func(t *testing.T) {
			_, _, err := ValidateQuery(statement)
			require.Error(t, err)
		})
	}
}

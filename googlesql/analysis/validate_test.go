package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/googlesql/parser"
)

func TestValidateQueryShared(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		readOnly    bool
		returnsData bool
		wantErr     bool
	}{
		// --- Valid read-only queries ---
		{name: "simple select", sql: "SELECT 1", readOnly: true, returnsData: true},
		{name: "select star no space", sql: "select* from t", readOnly: true, returnsData: true},
		{name: "select with where and order", sql: "SELECT a, b FROM ds.t WHERE a > 1 ORDER BY b LIMIT 10", readOnly: true, returnsData: true},
		{name: "with cte select", sql: "WITH c AS (SELECT 1 AS n) SELECT n FROM c", readOnly: true, returnsData: true},
		{name: "multiple ctes", sql: "WITH a AS (SELECT 1 AS x), b AS (SELECT x FROM a) SELECT * FROM b;", readOnly: true, returnsData: true},
		{name: "recursive cte", sql: "WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT * FROM r", readOnly: true, returnsData: true},
		{name: "union all", sql: "SELECT a FROM t1 UNION ALL SELECT a FROM t2", readOnly: true, returnsData: true},
		{name: "intersect distinct", sql: "SELECT a FROM t1 INTERSECT DISTINCT SELECT a FROM t2", readOnly: true, returnsData: true},
		{name: "except distinct", sql: "SELECT a FROM t1 EXCEPT DISTINCT SELECT a FROM t2", readOnly: true, returnsData: true},
		{name: "parenthesized set op", sql: "(SELECT 1) UNION ALL (SELECT 2)", readOnly: true, returnsData: true},
		{name: "table query", sql: "TABLE t", readOnly: true, returnsData: true},
		{name: "show tables", sql: "SHOW TABLES", readOnly: true, returnsData: true},
		{name: "describe table", sql: "DESCRIBE t", readOnly: true, returnsData: true},
		{name: "desc table", sql: "DESC t", readOnly: true, returnsData: true},
		{name: "explain select", sql: "EXPLAIN SELECT * FROM t;", readOnly: true, returnsData: true},
		{name: "explain with cte", sql: "EXPLAIN WITH c AS (SELECT 1 AS x) SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "explain set op", sql: "EXPLAIN (SELECT 1) UNION ALL (SELECT 2)", readOnly: true, returnsData: true},
		{name: "graph query", sql: "GRAPH g MATCH (n) RETURN n", readOnly: true, returnsData: true},
		{name: "pipe read query", sql: "SELECT x FROM t |> WHERE x > 1 |> AGGREGATE COUNT(*) GROUP BY x", readOnly: true, returnsData: true},
		{name: "pipe tvf call", sql: "SELECT x FROM t |> CALL f()", readOnly: true, returnsData: true},
		{name: "show with read pipe", sql: "SHOW TABLES |> WHERE true", readOnly: true, returnsData: true},
		{name: "explain show", sql: "EXPLAIN SHOW TABLES", readOnly: true, returnsData: true},
		{name: "explain describe", sql: "EXPLAIN DESCRIBE t", readOnly: true, returnsData: true},
		{name: "multi-statement reads", sql: "SELECT 1; SELECT 2;", readOnly: true, returnsData: true},

		// --- Strings, backtick identifiers, and comments ---
		{name: "semicolon and comment in string literal", sql: "SELECT * FROM t WHERE a = 'foo; DROP TABLE t; -- bar'", readOnly: true, returnsData: true},
		{name: "dml keywords in string, identifier, and comment", sql: "WITH t AS (SELECT * FROM t1 WHERE a = 'insert'), tx AS (SELECT * FROM `delete`) /* UPDATE */ SELECT `update` FROM t;", readOnly: true, returnsData: true},
		{name: "triple-quoted string with semicolon", sql: `SELECT """a; DROP TABLE t""" FROM t`, readOnly: true, returnsData: true},
		{name: "raw string with semicolon", sql: `SELECT r'a; DELETE FROM t' FROM t`, readOnly: true, returnsData: true},
		{name: "hash and dash line comments", sql: "# DROP TABLE t\n-- DELETE FROM t\nSELECT 1", readOnly: true, returnsData: true},
		{name: "statement hint before select", sql: "@{optimizer_version=1} SELECT * FROM t", readOnly: true, returnsData: true},

		// --- Empty / whitespace / comment-only ---
		{name: "empty input", sql: "", readOnly: true, returnsData: true},
		{name: "whitespace only", sql: "   \n\t  ", readOnly: true, returnsData: true},
		{name: "line comment only", sql: "-- just a comment\n", readOnly: true, returnsData: true},
		{name: "block comment only", sql: "/* block comment */", readOnly: true, returnsData: true},
		{name: "semicolons only", sql: ";;;", readOnly: true, returnsData: true},

		// --- DML writes ---
		{name: "insert", sql: "INSERT INTO t (a) VALUES (1)"},
		{name: "insert select", sql: "INSERT INTO t SELECT a FROM u"},
		{name: "update", sql: "UPDATE t SET a = 1 WHERE b = 2"},
		{name: "delete", sql: "DELETE FROM t WHERE a = 1"},
		{name: "merge", sql: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE"},

		// --- DDL & DCL ---
		{name: "create view", sql: "CREATE VIEW v AS SELECT 1"},
		{name: "alter table", sql: "ALTER TABLE t ADD COLUMN b INT64"},
		{name: "drop table", sql: "DROP TABLE t"},
		{name: "drop view", sql: "DROP VIEW v"},
		{name: "rename table", sql: "RENAME TABLE a TO b"},
		{name: "undrop schema", sql: "UNDROP SCHEMA s"},
		{name: "grant", sql: "GRANT SELECT ON TABLE t TO ROLE r"},
		{name: "revoke", sql: "REVOKE SELECT ON TABLE t FROM ROLE r"},

		// --- Data movement, calls, dynamic SQL, session, transactions ---
		{name: "export data", sql: "EXPORT DATA OPTIONS(uri='gs://bucket/file-*.csv', format='CSV') AS SELECT * FROM t"},
		{name: "export model", sql: "EXPORT MODEL ds.m OPTIONS(uri='gs://bucket/model')"},
		{name: "load data", sql: "LOAD DATA INTO ds.t FROM FILES(format='CSV', uris=['gs://bucket/file.csv'])"},
		{name: "clone data", sql: "CLONE DATA INTO ds.t FROM ds.s"},
		{name: "call procedure", sql: "CALL ds.proc(1)"},
		{name: "call with pipe", sql: "CALL ds.proc(1) |> WHERE true"},
		{name: "execute immediate", sql: "EXECUTE IMMEDIATE 'SELECT 1'"},
		{name: "execute immediate with pipe", sql: "EXECUTE IMMEDIATE 'SELECT 1' |> WHERE true"},
		{name: "set variable", sql: "SET x = 1"},
		{name: "set system variable", sql: "SET @@dataset_id = 'ds'"},
		{name: "set transaction", sql: "SET TRANSACTION READ ONLY"},
		{name: "begin transaction", sql: "BEGIN TRANSACTION"},
		{name: "commit", sql: "COMMIT"},
		{name: "rollback", sql: "ROLLBACK"},
		{name: "analyze", sql: "ANALYZE"},
		{name: "assert", sql: "ASSERT 1 = 1"},

		// --- Writing pipe operators ---
		{name: "pipe insert", sql: "SELECT 1 |> INSERT t (x)"},
		{name: "pipe create table", sql: "SELECT 1 |> CREATE TABLE t (x INT64)"},
		{name: "pipe export data", sql: "SELECT 1 |> EXPORT DATA OPTIONS(uri='gs://bucket/x')"},
		{name: "pipe if write branch", sql: "SELECT 1 |> IF true THEN (|> INSERT t) ELSE (|> DISTINCT)"},
		{name: "pipe fork write branch", sql: "SELECT 1 |> FORK (|> INSERT t), (|> DISTINCT)"},
		{name: "pipe tee export branch", sql: "SELECT 1 |> TEE (|> EXPORT DATA OPTIONS(uri='gs://bucket/x')), (|> DISTINCT)"},
		{name: "pipe log write", sql: "SELECT 1 |> LOG (|> INSERT t)"},
		{name: "cte with pipe insert", sql: "WITH c AS (SELECT 1 |> INSERT t) SELECT * FROM c"},
		{name: "show with pipe insert", sql: "SHOW TABLES |> INSERT t"},

		// --- EXPLAIN over non-read-only statements ---
		{name: "explain analyze", sql: "EXPLAIN ANALYZE"},
		{name: "explain insert", sql: "EXPLAIN INSERT INTO t (a) VALUES (1)"},
		{name: "explain update", sql: "EXPLAIN UPDATE t SET a = 1 WHERE b = 2"},
		{name: "explain delete", sql: "EXPLAIN DELETE FROM t WHERE a = 1"},
		{name: "explain drop table", sql: "EXPLAIN DROP TABLE t"},
		{name: "explain call", sql: "EXPLAIN CALL ds.proc()"},
		{name: "explain export data", sql: "EXPLAIN EXPORT DATA OPTIONS(uri='gs://bucket/x') AS SELECT 1"},
		{name: "explain pipe insert", sql: "EXPLAIN SELECT 1 |> INSERT t"},
		{name: "nested explain", sql: "EXPLAIN EXPLAIN SELECT 1"},

		// --- Multi-statement injection attempts ---
		{name: "select then drop", sql: "SELECT 1; DROP TABLE t"},
		{name: "select then delete", sql: "SELECT 1; DELETE FROM employee WHERE true"},
		{name: "select then insert", sql: "SELECT 1; INSERT INTO t (a) VALUES (1)"},
		{name: "write then select", sql: "DELETE FROM t WHERE a = 1; SELECT 1"},
		{name: "select then call", sql: "SELECT 1; CALL ds.proc()"},
		{name: "select then set", sql: "SELECT 1; SET @@dataset_id = 'other'"},

		// --- Syntax and lex errors ---
		{name: "explain analyze select", sql: "EXPLAIN ANALYZE SELECT * FROM t", wantErr: true},
		{name: "cte with update", sql: "WITH t AS (SELECT * FROM t1) UPDATE t SET a = 1", wantErr: true},
		{name: "cte with insert", sql: "WITH t AS (SELECT * FROM t1) INSERT INTO t VALUES (1)", wantErr: true},
		{name: "malformed select", sql: "SELECT FROM WHERE", wantErr: true},
		{name: "trailing junk", sql: "SELECT 1 FROM t a b c", wantErr: true},
		{name: "unterminated string", sql: "SELECT 'unterminated", wantErr: true},
		{name: "unterminated block comment", sql: "/*", wantErr: true},
		{name: "unterminated block comment after select", sql: "SELECT 1; /*", wantErr: true},
		{name: "unterminated backtick", sql: "SELECT `col FROM t", wantErr: true},
		{name: "select then syntax error", sql: "SELECT 1; SELECT FROM", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tc.sql)
			if tc.wantErr {
				require.Error(t, err, "sql: %s", tc.sql)
				require.False(t, readOnly)
				require.False(t, returnsData)
				var perr *parser.ParseError
				require.ErrorAs(t, err, &perr)
				return
			}
			require.NoError(t, err, "sql: %s", tc.sql)
			require.Equal(t, tc.readOnly, readOnly, "readOnly for %q", tc.sql)
			require.Equal(t, tc.returnsData, returnsData, "returnsData for %q", tc.sql)
		})
	}
}

func TestValidateQueryDialectSpecific(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		readOnly    bool
		returnsData bool
		wantErr     bool
	}{
		// BigQuery-specific syntax

		// Spanner-specific syntax
		{name: "spanner qualify as table alias", sql: "SELECT qualify.a FROM t qualify", readOnly: true, returnsData: true},
		{name: "spanner create table is write", sql: "CREATE TABLE t (a INT64) PRIMARY KEY (a)"},
		{name: "spanner for update rejected as non-read-only", sql: "SELECT * FROM t FOR UPDATE"},
		{name: "spanner subquery for update rejected", sql: "SELECT * FROM (SELECT * FROM t FOR UPDATE)"},
		{name: "spanner scalar subquery for update rejected", sql: "SELECT (SELECT a FROM t FOR UPDATE)"},
		{name: "spanner explain for update rejected", sql: "EXPLAIN SELECT * FROM t FOR UPDATE"},
		{name: "spanner create change stream is write", sql: "CREATE CHANGE STREAM cs FOR ALL"},
		{name: "spanner rejects describe with pipe", sql: "DESCRIBE t |> AGGREGATE COUNT(*)", wantErr: true},
		{name: "spanner rejects begin end block", sql: "BEGIN SELECT 1; END;", wantErr: true},
		{name: "spanner rejects bigquery from-first", sql: "FROM t |> SELECT x", wantErr: true},
		{name: "spanner rejects bigquery qualify clause", sql: "SELECT x FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY x) = 1", wantErr: true},
		{name: "spanner rejects declare", sql: "DECLARE x INT64 DEFAULT 1", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tc.sql)
			if tc.wantErr {
				require.Error(t, err, "sql: %s", tc.sql)
				require.False(t, readOnly)
				require.False(t, returnsData)
				return
			}
			require.NoError(t, err, "sql: %s", tc.sql)
			require.Equal(t, tc.readOnly, readOnly, "readOnly for %q", tc.sql)
			require.Equal(t, tc.returnsData, returnsData, "returnsData for %q", tc.sql)
		})
	}
}

func TestClassifySQLDialectSpanner(t *testing.T) {
	// Spanner-specific syntax that requires parser.DialectSpanner:
	// unreserved QUALIFY table alias and LANGUAGE REMOTE functions.
	require.Equal(t, Select, ClassifySQL("SELECT qualify.a FROM t qualify"))
	require.Equal(t, DDL, ClassifySQL("CREATE FUNCTION f(x INT64) RETURNS INT64 LANGUAGE REMOTE OPTIONS (endpoint='https://example.com')"))
}

package analysis

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateQuery(t *testing.T) {
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
		{name: "select paren no space", sql: "SELECT(1)", readOnly: true, returnsData: true},
		{name: "select with catalog and schema", sql: "SELECT a, b FROM main.app.t WHERE a > 1 ORDER BY b LIMIT 10;", readOnly: true, returnsData: true},
		{name: "with cte select", sql: "WITH c AS (SELECT 1 AS n) SELECT n FROM c", readOnly: true, returnsData: true},
		{name: "with cte column list", sql: "WITH c (a, b) AS (SELECT 1, 2) SELECT a, b FROM c", readOnly: true, returnsData: true},
		{name: "with backtick cte name", sql: "WITH `my cte` AS (SELECT 1) SELECT * FROM `my cte`", readOnly: true, returnsData: true},
		{name: "multiple ctes", sql: "WITH a AS (SELECT 1 AS x), b AS (SELECT x FROM a) SELECT * FROM b;", readOnly: true, returnsData: true},
		{name: "recursive cte", sql: "WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT * FROM r", readOnly: true, returnsData: true},
		{name: "cte with materialized hint", sql: "WITH c AS MATERIALIZED (SELECT 1), d AS NOT MATERIALIZED (SELECT 2) SELECT * FROM c UNION ALL SELECT * FROM d", readOnly: true, returnsData: true},
		{name: "cte with values", sql: "WITH c AS (VALUES (1, 'a')) SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "cte followed by values", sql: "WITH c AS (SELECT 1) VALUES (2)", readOnly: true, returnsData: true},
		{name: "cte followed by parenthesized union", sql: "WITH c AS (SELECT 1) (SELECT * FROM c) UNION ALL (SELECT 2)", readOnly: true, returnsData: true},
		{name: "values", sql: "VALUES (1, 'a'), (2, 'b')", readOnly: true, returnsData: true},
		{name: "table shorthand", sql: "TABLE main.app.users", readOnly: true, returnsData: true},
		{name: "parenthesized select", sql: "(SELECT 1)", readOnly: true, returnsData: true},
		{name: "nested parenthesized select", sql: "((SELECT 1))", readOnly: true, returnsData: true},
		{name: "parenthesized union all", sql: "(SELECT 1) UNION ALL (SELECT 2)", readOnly: true, returnsData: true},
		{name: "parenthesized union by name", sql: "(SELECT 1 AS a) UNION ALL BY NAME (SELECT 2 AS a)", readOnly: true, returnsData: true},
		{name: "parenthesized select with limit", sql: "(SELECT 1) ORDER BY 1 LIMIT 1", readOnly: true, returnsData: true},
		{name: "parenthesized select with pipe", sql: "(SELECT 1 AS x) |> WHERE x > 0", readOnly: true, returnsData: true},
		{name: "from-first pipeline select", sql: "FROM t |> SELECT *", readOnly: true, returnsData: true},
		{name: "from-first pipeline aggregate", sql: "FROM main.app.t |> WHERE x > 1 |> AGGREGATE COUNT(*) GROUP BY y", readOnly: true, returnsData: true},
		{name: "from-first pipeline extend drop set columns", sql: "FROM t |> EXTEND a + 1 AS b |> DROP a |> SET b = b * 2", readOnly: true, returnsData: true},
		{name: "from-first pipeline select scalar replace", sql: "FROM t |> SELECT replace(name, 'a', 'b')", readOnly: true, returnsData: true},
		{name: "from-first pipeline extend scalar replace", sql: "FROM t |> EXTEND replace(name, 'a', 'b') AS normalized", readOnly: true, returnsData: true},
		{name: "from-first pipeline parenthesized scalar replace", sql: "FROM t |> SELECT (replace(name, 'a', 'b'))", readOnly: true, returnsData: true},
		{name: "from-first pipeline join on scalar replace", sql: "FROM t JOIN u ON replace(t.a, 'x', 'y') = u.a |> SELECT *", readOnly: true, returnsData: true},
		{name: "from-first pipeline over subquery", sql: "FROM (SELECT 1 AS a) |> SELECT a", readOnly: true, returnsData: true},
		{name: "cte with from-first pipeline body", sql: "WITH c AS (FROM t |> WHERE x > 1) SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "cte feeding from-first pipeline", sql: "WITH c AS (FROM t |> WHERE x > 1) FROM c |> SELECT *", readOnly: true, returnsData: true},
		{name: "parenthesized from-first pipeline union", sql: "(FROM t |> SELECT *) UNION ALL (FROM u |> SELECT *)", readOnly: true, returnsData: true},
		{name: "explain from-first pipeline", sql: "EXPLAIN FORMATTED FROM t |> WHERE x > 1 |> SELECT *", readOnly: true, returnsData: true},
		{name: "describe query from-first pipeline", sql: "DESCRIBE QUERY FROM t |> SELECT *", readOnly: true, returnsData: true},
		{name: "show tables", sql: "SHOW TABLES", readOnly: true, returnsData: true},
		{name: "show create table", sql: "SHOW CREATE TABLE main.app.t", readOnly: true, returnsData: true},
		{name: "show schemas", sql: "SHOW SCHEMAS IN main", readOnly: true, returnsData: true},
		{name: "describe table", sql: "DESCRIBE t", readOnly: true, returnsData: true},
		{name: "desc table", sql: "DESC t", readOnly: true, returnsData: true},
		{name: "describe extended", sql: "DESCRIBE EXTENDED main.app.t", readOnly: true, returnsData: true},
		{name: "describe detail", sql: "DESCRIBE DETAIL main.app.t", readOnly: true, returnsData: true},
		{name: "describe history", sql: "DESCRIBE HISTORY main.app.t", readOnly: true, returnsData: true},
		{name: "describe query select", sql: "DESCRIBE QUERY SELECT 1", readOnly: true, returnsData: true},
		{name: "describe query with cte", sql: "DESCRIBE QUERY WITH c AS (SELECT 1) SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "explain select", sql: "EXPLAIN SELECT * FROM t;", readOnly: true, returnsData: true},
		{name: "explain formatted select", sql: "EXPLAIN FORMATTED SELECT * FROM t", readOnly: true, returnsData: true},
		{name: "explain extended cte", sql: "EXPLAIN EXTENDED WITH c AS (SELECT 1) SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "explain cost values", sql: "EXPLAIN COST VALUES (1, 2)", readOnly: true, returnsData: true},
		{name: "explain codegen table", sql: "EXPLAIN CODEGEN TABLE t", readOnly: true, returnsData: true},
		{name: "explain show", sql: "EXPLAIN SHOW TABLES", readOnly: true, returnsData: true},
		{name: "explain describe", sql: "EXPLAIN DESCRIBE t", readOnly: true, returnsData: true},
		{name: "multi-statement reads single line", sql: "SELECT 1; SELECT 2;", readOnly: true, returnsData: true},
		{name: "multi-statement reads multi line", sql: "SELECT 1;\nSHOW TABLES;\nDESCRIBE t;", readOnly: true, returnsData: true},

		// --- Strings, identifiers, raw strings, and nested comments ---
		{name: "semicolon and dml in single quoted string", sql: "SELECT * FROM t WHERE a = 'foo; DROP TABLE t; -- bar'", readOnly: true, returnsData: true},
		{name: "semicolon and dml in double quoted string", sql: `SELECT * FROM t WHERE a = "foo; DROP TABLE t; /* bar"`, readOnly: true, returnsData: true},
		{name: "escaped single quotes", sql: `SELECT 'it\'s; DROP TABLE t', 'it''s; DELETE FROM t'`, readOnly: true, returnsData: true},
		{name: "escaped double quotes", sql: `SELECT "a\"; DROP TABLE t", "a""; DELETE FROM t"`, readOnly: true, returnsData: true},
		{name: "raw string with semicolon inside", sql: `SELECT r'foo; DROP TABLE t', R"bar; DELETE FROM t"`, readOnly: true, returnsData: true},
		{name: "hex literal", sql: `SELECT X'1A2B', x'ff'`, readOnly: true, returnsData: true},
		{name: "dml keywords in string identifier and comment inside cte", sql: "WITH t AS (SELECT * FROM t1 WHERE a = 'insert'), tx AS (SELECT * FROM `delete`) /* UPDATE */ SELECT `update` FROM t;", readOnly: true, returnsData: true},
		{name: "backtick escaped identifier", sql: "SELECT `a``; DROP TABLE t` FROM t", readOnly: true, returnsData: true},
		{name: "leading line comment", sql: "-- DROP TABLE t\nSELECT 1", readOnly: true, returnsData: true},
		{name: "line comment terminated by CR before select", sql: "-- DROP TABLE t\rSELECT 1", readOnly: true, returnsData: true},
		{name: "line comment with backslash-LF continuation before select", sql: "-- first\\\nstill comment\nSELECT 1", readOnly: true, returnsData: true},
		{name: "leading inline block comment", sql: "/* comment */ SELECT 1;", readOnly: true, returnsData: true},
		{name: "nested block comment before select", sql: "/* outer /* inner ; DROP TABLE t */ still comment */ SELECT 1", readOnly: true, returnsData: true},
		{name: "line comment containing block open", sql: "-- /* unclosed\nSELECT 1", readOnly: true, returnsData: true},
		{name: "block comment containing line comment", sql: "/* -- */ SELECT 1", readOnly: true, returnsData: true},

		// --- Empty / whitespace / comment-only / semicolons-only ---
		{name: "empty input", sql: "", readOnly: true, returnsData: true},
		{name: "whitespace only", sql: "   \n\t\r  ", readOnly: true, returnsData: true},
		{name: "line comment only", sql: "-- just a comment\n", readOnly: true, returnsData: true},
		{name: "block comment only", sql: "/* block comment */", readOnly: true, returnsData: true},
		{name: "nested block comment only", sql: "/* outer /* inner */ */", readOnly: true, returnsData: true},
		{name: "semicolons only", sql: ";;;", readOnly: true, returnsData: true},
		{name: "comments and semicolons only", sql: "; -- c\n /* a /* b */ */ ;", readOnly: true, returnsData: true},

		// --- DML writes ---
		{name: "insert values", sql: "INSERT INTO t VALUES (1)"},
		{name: "insert overwrite", sql: "INSERT OVERWRITE t SELECT * FROM s"},
		{name: "update", sql: "UPDATE t SET a = 1 WHERE b = 2"},
		{name: "delete", sql: "DELETE FROM t WHERE a = 1"},
		{name: "merge", sql: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE"},
		{name: "copy into", sql: "COPY INTO t FROM 's3://bucket/path' FILEFORMAT = PARQUET"},

		// --- DDL & Delta/catalog maintenance ---
		{name: "create table", sql: "CREATE TABLE t (a INT)"},
		{name: "create view", sql: "CREATE VIEW v AS SELECT 1"},
		{name: "alter table", sql: "ALTER TABLE t ADD COLUMNS (b STRING)"},
		{name: "drop table", sql: "DROP TABLE t"},
		{name: "truncate table", sql: "TRUNCATE TABLE t"},
		{name: "optimize", sql: "OPTIMIZE t ZORDER BY (a)"},
		{name: "vacuum", sql: "VACUUM t RETAIN 168 HOURS"},
		{name: "restore", sql: "RESTORE TABLE t TO VERSION AS OF 1"},
		{name: "refresh", sql: "REFRESH TABLE t"},
		{name: "msck", sql: "MSCK REPAIR TABLE t"},
		{name: "replace table", sql: "REPLACE TABLE t AS SELECT 1"},
		{name: "analyze table", sql: "ANALYZE TABLE t COMPUTE STATISTICS"},

		// --- DCL, session, procedural ---
		{name: "grant", sql: "GRANT SELECT ON TABLE t TO `user@example.com`"},
		{name: "revoke", sql: "REVOKE SELECT ON TABLE t FROM `user@example.com`"},
		{name: "deny", sql: "DENY SELECT ON TABLE t TO `user@example.com`"},
		{name: "set config", sql: "SET spark.sql.ansi.enabled = true"},
		{name: "use catalog", sql: "USE CATALOG main"},
		{name: "use schema", sql: "USE SCHEMA app"},
		{name: "call procedure", sql: "CALL app.my_proc(1)"},
		{name: "execute immediate", sql: "EXECUTE IMMEDIATE 'DROP TABLE t'"},
		{name: "from-first insert", sql: "FROM s INSERT INTO t SELECT *"},
		{name: "from-first insert overwrite", sql: "FROM s INSERT OVERWRITE t SELECT *"},
		{name: "bare from without pipeline", sql: "FROM t"},
		{name: "from pipeline missing relation", sql: "FROM |> SELECT *"},
		{name: "from pipeline trailing pipe", sql: "FROM t |>"},
		{name: "from pipeline with insert stage", sql: "FROM t |> INSERT INTO u"},
		{name: "from pipeline with export stage", sql: "FROM t |> EXPORT DATA"},
		{name: "from pipeline with trailing insert", sql: "FROM t |> WHERE x > 1 INSERT INTO u SELECT *"},
		{name: "from-first insert followed by pipe", sql: "FROM s INSERT INTO t SELECT * |> SELECT *"},
		{name: "from pipeline union with from-first insert", sql: "FROM t |> UNION ALL (FROM u INSERT INTO w SELECT *)"},
		{name: "from pipeline with mutating cte subquery", sql: "FROM (WITH c AS (DELETE FROM t) SELECT * FROM c) |> SELECT *"},
		{name: "parenthesized select with pipe insert", sql: "(SELECT 1) |> INSERT INTO t"},
		{name: "explain analyze from-first pipeline", sql: "EXPLAIN ANALYZE FROM t |> SELECT *"},

		// --- CTEs with writes or malformed CTEs ---
		{name: "cte feeding insert", sql: "WITH c AS (SELECT 1) INSERT INTO t SELECT * FROM c"},
		{name: "cte feeding update", sql: "WITH c AS (SELECT 1) UPDATE t SET a = 1"},
		{name: "cte feeding delete", sql: "WITH c AS (SELECT 1) DELETE FROM t"},
		{name: "cte feeding merge", sql: "WITH c AS (SELECT 1) MERGE INTO t USING c ON t.id = c.id WHEN MATCHED THEN DELETE"},
		{name: "cte body insert", sql: "WITH c AS (INSERT INTO t VALUES (1) RETURNING *) SELECT * FROM c"},
		{name: "cte body update", sql: "WITH c AS (UPDATE t SET a = 1) SELECT * FROM c"},
		{name: "cte body delete", sql: "WITH c AS (DELETE FROM t) SELECT * FROM c"},
		{name: "cte body merge", sql: "WITH c AS (MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE) SELECT * FROM c"},
		{name: "second cte body delete", sql: "WITH a AS (SELECT 1), b AS (DELETE FROM t) SELECT * FROM a"},
		{name: "nested cte body delete", sql: "WITH a AS (WITH b AS (DELETE FROM t) SELECT * FROM b) SELECT * FROM a"},
		{name: "cte missing main query", sql: "WITH c AS (SELECT 1)"},
		{name: "cte with invalid NOT without MATERIALIZED", sql: "WITH c AS NOT (SELECT 1) SELECT * FROM c"},
		{name: "cte feeding from-first insert", sql: "WITH c AS (SELECT 1) FROM c INSERT INTO t SELECT *"},

		// --- EXPLAIN over writes / EXPLAIN ANALYZE / nested EXPLAIN ---
		{name: "explain analyze select", sql: "EXPLAIN ANALYZE SELECT 1"},
		{name: "explain comment analyze select", sql: "EXPLAIN /* c */ ANALYZE SELECT 1"},
		{name: "explain bare", sql: "EXPLAIN"},
		{name: "explain formatted bare", sql: "EXPLAIN FORMATTED"},
		{name: "explain insert", sql: "EXPLAIN INSERT INTO t VALUES (1)"},
		{name: "explain formatted insert", sql: "EXPLAIN FORMATTED INSERT INTO t VALUES (1)"},
		{name: "explain extended update", sql: "EXPLAIN EXTENDED UPDATE t SET a = 1"},
		{name: "explain cost delete", sql: "EXPLAIN COST DELETE FROM t"},
		{name: "explain codegen drop", sql: "EXPLAIN CODEGEN DROP TABLE t"},
		{name: "explain merge", sql: "EXPLAIN MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE"},
		{name: "explain cte insert", sql: "EXPLAIN WITH c AS (SELECT 1) INSERT INTO t SELECT * FROM c"},
		{name: "explain call", sql: "EXPLAIN CALL app.proc()"},
		{name: "nested explain", sql: "EXPLAIN EXPLAIN SELECT 1"},
		{name: "describe query insert", sql: "DESCRIBE QUERY INSERT INTO t VALUES (1)"},
		{name: "describe insert", sql: "DESCRIBE INSERT INTO t VALUES (1)"},

		// --- Multi-statement & comment/raw-string injection attempts ---
		{name: "single-line select then drop", sql: "SELECT 1; DROP TABLE t"},
		{name: "single-line select then drop with semicolon", sql: "SELECT 1; DROP TABLE t;"},
		{name: "multi-line select then delete", sql: "SELECT 1;\nDELETE FROM t WHERE true;"},
		{name: "write then select", sql: "DROP TABLE t; SELECT 1"},
		{name: "select then set", sql: "SELECT 1; SET spark.sql.ansi.enabled = false"},
		{name: "select then use", sql: "SELECT 1; USE CATALOG other"},
		{name: "select then call", sql: "SELECT 1; CALL app.proc()"},
		{name: "raw string trailing backslash ends string before drop", sql: `SELECT r'foo\' ; DROP TABLE t`},
		{name: "raw uppercase double-quote trailing backslash ends string before drop", sql: `SELECT R"foo\" ; DROP TABLE t`},
		{name: "escaped backslash before closing quote then drop", sql: `SELECT 'foo\\' ; DROP TABLE t`},
		{name: "escaped quote then closing quote then drop", sql: `SELECT 'a\'' ; DROP TABLE t`},
		{name: "adjacent quotes then drop", sql: `SELECT 'a'''; DROP TABLE t`},
		{name: "nested comment hiding select before drop", sql: "/* outer /* inner */ SELECT 1 */ DROP TABLE t"},
		{name: "nested comment hiding explain before delete", sql: "/* outer /* inner */ EXPLAIN */ DELETE FROM t"},
		{name: "line comment terminated by CR before delete", sql: "-- comment\rDELETE FROM t"},
		{name: "line comment terminated by CRLF before delete", sql: "-- comment\r\nDELETE FROM t"},
		{name: "line comment backslash-LF hides select before delete", sql: "-- comment\\\nSELECT 1\nDELETE FROM t"},
		{name: "line comment backslash-CR does not continue past CR", sql: "-- comment\\\rDELETE FROM t"},

		// --- Unclassifiable / malformed statements (fail closed) ---
		{name: "unknown keyword", sql: "FOOBAR 1"},
		{name: "bare select", sql: "SELECT"},
		{name: "bare show", sql: "SHOW"},
		{name: "bare describe", sql: "DESCRIBE"},
		{name: "bare values", sql: "VALUES"},
		{name: "bare table", sql: "TABLE"},
		{name: "empty parens", sql: "()"},
		{name: "unbalanced open paren", sql: "(SELECT 1"},
		{name: "unbalanced close paren", sql: "SELECT 1)"},
		{name: "parenthesized drop", sql: "(DROP TABLE t)"},
		{name: "parenthesized select followed by drop", sql: "(SELECT 1) DROP TABLE t"},
		{name: "parenthesized union with drop rhs", sql: "(SELECT 1) UNION ALL (DROP TABLE t)"},
		{name: "delimiter directive", sql: "DELIMITER ;;"},

		// --- Lexical errors (fail closed with error) ---
		{name: "unterminated single string", sql: "SELECT 'unterminated", wantErr: true},
		{name: "unterminated double string", sql: `SELECT "unterminated`, wantErr: true},
		{name: "unterminated raw string", sql: "SELECT r'unterminated", wantErr: true},
		{name: "unterminated hex literal", sql: "SELECT X'1234", wantErr: true},
		{name: "unterminated string with trailing backslash", sql: `SELECT 'abc\`, wantErr: true},
		{name: "unterminated backtick identifier", sql: "SELECT `col FROM t", wantErr: true},
		{name: "unterminated block comment", sql: "/* unclosed", wantErr: true},
		{name: "unterminated nested block comment", sql: "/* outer /* inner */ SELECT 1", wantErr: true},
		{name: "unterminated block comment after select", sql: "SELECT 1; /* unclosed", wantErr: true},
		{name: "lexical error after write still returns error", sql: "DROP TABLE t; SELECT 'unclosed", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tc.sql)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ValidateQuery(%q) = (%v, %v, nil), want error", tc.sql, readOnly, returnsData)
				}
				if readOnly || returnsData {
					t.Errorf("ValidateQuery(%q) on error returned (%v, %v), want (false, false)", tc.sql, readOnly, returnsData)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateQuery(%q) unexpected error: %v", tc.sql, err)
			}
			if readOnly != tc.readOnly || returnsData != tc.returnsData {
				t.Errorf("ValidateQuery(%q) = (%v, %v), want (%v, %v)", tc.sql, readOnly, returnsData, tc.readOnly, tc.returnsData)
			}
		})
	}
}

func TestValidateQueryNestingAndScale(t *testing.T) {
	t.Run("deep nesting within limit stays linear", func(t *testing.T) {
		const depth = 5000
		sql := strings.Repeat("(", depth) + "SELECT 1" + strings.Repeat(")", depth)
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Fatalf("ValidateQuery(depth=%d) err = %v", depth, err)
		}
		if !readOnly || !returnsData {
			t.Errorf("ValidateQuery(depth=%d) = (%v, %v), want (true, true)", depth, readOnly, returnsData)
		}
	})

	t.Run("nesting past maxNestingDepth fails closed with error", func(t *testing.T) {
		const depth = maxNestingDepth + 1
		sql := strings.Repeat("(", depth) + "SELECT 1" + strings.Repeat(")", depth)
		readOnly, returnsData, err := ValidateQuery(sql)
		if !errors.Is(err, errNestingDepthExceeded) {
			t.Fatalf("ValidateQuery(depth=%d) err = %v, want %v", depth, err, errNestingDepthExceeded)
		}
		if readOnly || returnsData {
			t.Errorf("ValidateQuery(depth=%d) = (%v, %v), want (false, false)", depth, readOnly, returnsData)
		}
	})

	t.Run("long parenthesized set-op chain stays iterative", func(t *testing.T) {
		var b strings.Builder
		b.WriteString("(SELECT 1)")
		for i := 0; i < 15000; i++ {
			b.WriteString(" UNION ALL (SELECT 1)")
		}
		readOnly, returnsData, err := ValidateQuery(b.String())
		if err != nil {
			t.Fatalf("ValidateQuery(long union) err = %v", err)
		}
		if !readOnly || !returnsData {
			t.Errorf("ValidateQuery(long union) = (%v, %v), want (true, true)", readOnly, returnsData)
		}
	})
}

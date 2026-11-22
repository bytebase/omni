package parser

import (
	"strings"
	"testing"
)

func TestParseStatements(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		kind StmtKind
	}{
		// Queries, CTEs, joins, window, pipe, QUALIFY
		{
			name: "select with cte and window",
			sql: `WITH ranked AS (
				SELECT id, dept, salary,
				       ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC NULLS LAST) AS rn
				FROM main.hr.employees
			)
			SELECT /*+ BROADCAST(ranked) */ id, salary
			FROM ranked
			WHERE rn <= 3
			ORDER BY salary DESC
			LIMIT 10 OFFSET 5`,
			kind: StmtQuery,
		},
		{
			name: "select qualify extension",
			sql:  `SELECT id, ROW_NUMBER() OVER (PARTITION BY k ORDER BY ts DESC) AS rn FROM events QUALIFY rn = 1`,
			kind: StmtQuery,
		},
		{
			name: "select qualify after unaliased expression",
			sql:  `SELECT row_number() OVER () QUALIFY row_number() OVER () = 1`,
			kind: StmtQuery,
		},
		{
			name: "select qualify after group by",
			sql:  `SELECT k, count(*) FROM t GROUP BY k QUALIFY row_number() OVER (ORDER BY k) = 1`,
			kind: StmtQuery,
		},
		{
			name: "first last any_value with window and filter suffixes",
			sql:  `SELECT first(x) OVER (ORDER BY x), last(x) OVER (), any_value(x) FILTER (WHERE x > 0), first(x IGNORE NULLS) OVER (PARTITION BY k) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "case expression parenthesized value and case column",
			sql:  `SELECT CASE (x) WHEN 1 THEN 'a' ELSE 'b' END, case, case(1, 2) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "select all and distinct as identifier prefixes",
			sql:  `SELECT all.x, distinct.x, all[0], distinct::STRING FROM t AS all`,
			kind: StmtQuery,
		},
		{
			name: "trim with column names leading trailing both",
			sql:  `SELECT trim(leading), trim(trailing), trim(both), trim(leading, 'x'), trim(leading || 'a') FROM t`,
			kind: StmtQuery,
		},
		{
			name: "identifier dynamic function and table call",
			sql:  `SELECT IDENTIFIER('upper')('abc') FROM IDENTIFIER('main.t')(1)`,
			kind: StmtQuery,
		},
		{
			name: "digit-leading identifier with exponent prefix",
			sql:  `CREATE TABLE 1e2data (id INT)`,
			kind: StmtCreateTable,
		},
		{
			name: "select from digit-leading identifier with exponent prefix",
			sql:  `SELECT * FROM 1e2data`,
			kind: StmtQuery,
		},
		{
			name: "backtick-quoted hyphenated table name",
			sql:  "CREATE TABLE `my-table` (a INT)",
			kind: StmtCreateTable,
		},
		{
			name: "set spark timezone with path separator",
			sql:  `SET spark.sql.session.timeZone = America/Los_Angeles`,
			kind: StmtSet,
		},
		{
			name: "set spark uri config value",
			sql:  `SET fs.s3a.endpoint = https://s3.amazonaws.com/bucket`,
			kind: StmtSet,
		},
		{
			name: "select lateral view and pivot",
			sql:  `SELECT * FROM sales LATERAL VIEW EXPLODE(items) t AS item PIVOT (SUM(amount) FOR quarter IN ('Q1', 'Q2'))`,
			kind: StmtQuery,
		},
		{
			name: "select unpivot and tablesample",
			sql:  `SELECT * FROM t TABLESAMPLE (10 PERCENT) REPEATABLE (42) UNPIVOT EXCLUDE NULLS (val FOR col IN (c1, c2))`,
			kind: StmtQuery,
		},
		{
			name: "select set ops and pipe operator",
			sql:  `FROM orders |> WHERE status = 'OPEN' |> EXTEND total * 1.1 AS gross |> AGGREGATE SUM(gross) AS revenue GROUP BY region |> ORDER BY revenue DESC`,
			kind: StmtQuery,
		},
		{
			name: "select expressions lambdas casts intervals",
			sql: `SELECT
				TRANSFORM(arr, x -> x + 1),
				TRANSFORM_KEYS(m, (k, v) -> CAST(k AS STRING)),
				TRY_CAST(x AS DECIMAL(10, 2)),
				x::MAP<STRING, ARRAY<STRUCT<a: INT, b: TIMESTAMP_NTZ>>>,
				CASE WHEN EXISTS (SELECT 1 FROM t2) THEN DATE '2026-01-01' + INTERVAL 2 DAYS ELSE TIMESTAMP '2026-01-02 00:00:00' END,
				SUBSTRING(name FROM 1 FOR 3),
				TRIM(BOTH ' ' FROM name),
				OVERLAY(name PLACING 'x' FROM 2 FOR 1),
				EXTRACT(YEAR FROM created_at),
				COUNT(*) FILTER (WHERE active = TRUE),
				1L + 2S + 3Y + 1.5F + 2.5D + 3.14BD,
				R'raw\string'
			FROM tbl`,
			kind: StmtQuery,
		},
		{
			name: "select values and table",
			sql:  `VALUES (1, 'a'), (2, 'b') AS t(id, name)`,
			kind: StmtQuery,
		},
		{
			name: "select table shorthand",
			sql:  `TABLE main.default.users UNION ALL TABLE main.default.admins`,
			kind: StmtQuery,
		},

		// DML
		{
			name: "insert into by name",
			sql:  `INSERT INTO main.sales.orders BY NAME SELECT id, amount FROM staging`,
			kind: StmtInsert,
		},
		{
			name: "insert replace where",
			sql:  `INSERT INTO events REPLACE WHERE dt = DATE '2026-01-01' SELECT * FROM raw_events`,
			kind: StmtInsert,
		},
		{
			name: "update table",
			sql:  `UPDATE main.sales.orders AS o SET status = 'CLOSED', updated_at = CURRENT_TIMESTAMP WHERE o.id = 42`,
			kind: StmtUpdate,
		},
		{
			name: "delete from table",
			sql:  `DELETE FROM main.sales.orders o WHERE o.status = 'CANCELLED'`,
			kind: StmtDelete,
		},
		{
			name: "merge with schema evolution",
			sql: `MERGE WITH SCHEMA EVOLUTION INTO target t
				USING (SELECT id, val FROM source) s
				ON t.id = s.id
				WHEN MATCHED AND s.val IS NULL THEN DELETE
				WHEN MATCHED THEN UPDATE SET t.val = s.val
				WHEN NOT MATCHED BY TARGET THEN INSERT (id, val) VALUES (s.id, s.val)
				WHEN NOT MATCHED BY SOURCE THEN DELETE`,
			kind: StmtMerge,
		},
		{
			name: "load data and truncate",
			sql:  `LOAD DATA LOCAL INPATH '/tmp/data.csv' OVERWRITE INTO TABLE t PARTITION (dt = '2026-01-01')`,
			kind: StmtLoadData,
		},
		{
			name: "truncate table",
			sql:  `TRUNCATE TABLE main.default.events PARTITION (dt = '2026-01-01')`,
			kind: StmtTruncate,
		},

		// DDL
		{
			name: "create table complex types and identity",
			sql: `CREATE OR REPLACE TABLE main.analytics.events (
				id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1),
				payload VARIANT,
				tags ARRAY<MAP<STRING, ARRAY<INT>>>,
				meta STRUCT<k: STRING COLLATE UTF8_LCASE, v: DECIMAL(18, 4) NOT NULL COMMENT 'amount'>,
				created_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP
			) USING DELTA
			CLUSTER BY (id)
			COMMENT 'events table'
			TBLPROPERTIES ('delta.enableChangeDataFeed' = 'true')`,
			kind: StmtReplaceTable,
		},
		{
			name: "create schema",
			sql:  `CREATE SCHEMA IF NOT EXISTS main.analytics COMMENT 'analytics schema' LOCATION 's3://bucket/analytics' WITH DBPROPERTIES ('owner' = 'data')`,
			kind: StmtCreateNamespace,
		},
		{
			name: "create view with schema binding",
			sql:  `CREATE OR REPLACE VIEW v (a COMMENT 'col a') WITH SCHEMA COMPENSATION AS SELECT 1`,
			kind: StmtCreateView,
		},
		{
			name: "create sql udf",
			sql:  `CREATE OR REPLACE FUNCTION add_one(x INT) RETURNS INT LANGUAGE SQL DETERMINISTIC RETURN x + 1`,
			kind: StmtCreateFunction,
		},
		{
			name: "alter table alter column and cluster by",
			sql:  `ALTER TABLE t ALTER COLUMN c TYPE ARRAY<MAP<STRING, INT>>`,
			kind: StmtAlterTable,
		},
		{
			name: "drop table",
			sql:  `DROP TABLE IF EXISTS main.default.old_events PURGE`,
			kind: StmtDropTable,
		},
		{
			name: "declare variable",
			sql:  `DECLARE OR REPLACE VARIABLE v INT DEFAULT 42`,
			kind: StmtDeclare,
		},

		// Utility & Scripting
		{
			name: "use schema and set catalog",
			sql:  `SET CATALOG main`,
			kind: StmtSetCatalog,
		},
		{
			name: "set spark config",
			sql:  `SET spark.sql.shuffle.partitions = 200`,
			kind: StmtSet,
		},
		{
			name: "set spark string config",
			sql:  `SET spark.sql.catalogImplementation = hive`,
			kind: StmtSet,
		},
		{
			name: "reset spark config",
			sql:  `RESET spark.sql.shuffle.partitions`,
			kind: StmtReset,
		},
		{
			name: "show tables",
			sql:  `SHOW TABLES IN main.default LIKE 'ev*'`,
			kind: StmtShow,
		},
		{
			name: "describe table formatted",
			sql:  `DESCRIBE TABLE FORMATTED main.default.events`,
			kind: StmtDescribe,
		},
		{
			name: "explain formatted",
			sql:  `EXPLAIN FORMATTED SELECT 1`,
			kind: StmtExplain,
		},
		{
			name: "analyze table",
			sql:  `ANALYZE TABLE main.default.events COMPUTE STATISTICS FOR ALL COLUMNS`,
			kind: StmtAnalyze,
		},
		{
			name: "comment on table",
			sql:  `COMMENT ON TABLE main.default.events IS 'event log'`,
			kind: StmtComment,
		},
		{
			name: "refresh unquoted resource path",
			sql:  `REFRESH /mnt/data/events.parquet`,
			kind: StmtRefresh,
		},
		{
			name: "refresh quoted resource path",
			sql:  `REFRESH '/mnt/data/events.parquet'`,
			kind: StmtRefresh,
		},
		{
			name: "not identifier postfix forms",
			sql:  `SELECT not.x, not[0], not::STRING, NOT not.x FROM t AS not`,
			kind: StmtQuery,
		},
		{
			name: "count all and distinct star",
			sql:  `SELECT count(ALL *), count(DISTINCT *), f(all * 2) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "lateral as relation or schema name",
			sql:  `SELECT * FROM lateral JOIN lateral.t ON lateral.id = t.id JOIN lateral AS l2 USING (id)`,
			kind: StmtQuery,
		},
		{
			name: "from statement with transform higher-order function",
			sql:  `FROM t SELECT TRANSFORM(arr, x -> x + 1)`,
			kind: StmtQuery,
		},
		{
			name: "select transform hive script",
			sql:  `SELECT TRANSFORM(a, b) USING 'cat' AS (c INT, d STRING) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "select all and distinct multiplication vs star",
			sql:  `SELECT all * 2 FROM t UNION ALL SELECT distinct * 2 FROM t UNION ALL SELECT ALL * FROM t UNION ALL SELECT DISTINCT * EXCEPT (a) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "case identifier arithmetic and unary simple case",
			sql:  `SELECT case + 1, case - 2, case + (CASE x WHEN 1 THEN 2 END), CASE -1 WHEN -1 THEN 'neg' END, CASE case + case + 1 WHEN 1 THEN 'ok' END FROM t`,
			kind: StmtQuery,
		},
		{
			name: "parenthesized map and reduce function calls vs transform subquery",
			sql:  `SELECT (MAP('a', 1)), (REDUCE(arr, 0, (acc, x) -> acc + x)), (MAP a USING 'cat') FROM t`,
			kind: StmtQuery,
		},
		{
			name: "with recursive as cte name with column aliases",
			sql:  `WITH recursive(x) AS (SELECT 1) SELECT x FROM recursive`,
			kind: StmtQuery,
		},
		{
			name: "with recursive as cte name without aliases",
			sql:  `WITH recursive AS (SELECT 1) SELECT * FROM recursive`,
			kind: StmtQuery,
		},
		{
			name: "with recursive modifier on recursive cte name",
			sql:  `WITH RECURSIVE recursive(x) AS (SELECT 1) SELECT x FROM recursive`,
			kind: StmtQuery,
		},
		{
			name: "case operand with parenthesized expression continuations",
			sql:  `SELECT CASE (x) + 1 WHEN 2 THEN 'two' END, CASE (x)::INT WHEN 2 THEN 'two' END FROM t`,
			kind: StmtQuery,
		},
		{
			name: "nested simple case with unary operands",
			sql:  `SELECT CASE CASE -1 WHEN -1 THEN 2 END WHEN 2 THEN 3 END, CASE abs(CASE -1 WHEN -1 THEN 2 END) WHEN 2 THEN 3 END, CASE CASE(CASE(x)) + 1 WHEN 2 THEN 'two' END FROM t`,
			kind: StmtQuery,
		},
		{
			name: "qualified digit-leading column identifiers",
			sql:  `SELECT t.123abc, t.123_abc, t.1e2data, .5, .5e2, .5F, .5D, .5BD FROM t`,
			kind: StmtQuery,
		},
		{
			name: "values inline table comma separated with another relation",
			sql:  `SELECT * FROM VALUES (1) AS v, other_table`,
			kind: StmtQuery,
		},
		{
			name: "insert into table named table",
			sql:  `INSERT INTO table VALUES (1)`,
			kind: StmtInsert,
		},
		{
			name: "insert into schema named table",
			sql:  `INSERT INTO table.t VALUES (1)`,
			kind: StmtInsert,
		},
		{
			name: "insert overwrite table named table",
			sql:  `INSERT OVERWRITE table VALUES (1)`,
			kind: StmtInsert,
		},
		{
			name: "table function with subquery and table args",
			sql:  `SELECT * FROM f(TABLE(SELECT 1), TABLE(db.t) PARTITION BY (a), x => TABLE((SELECT 2)))`,
			kind: StmtQuery,
		},
		{
			name: "refresh cache repair",
			sql:  `MSCK REPAIR TABLE main.default.events SYNC PARTITIONS`,
			kind: StmtRepair,
		},
		{
			name: "call and execute immediate",
			sql:  `EXECUTE IMMEDIATE 'SELECT ? + ?' INTO x USING 1, 2`,
			kind: StmtExecuteImmediate,
		},
		{
			name: "compound script block",
			sql: `BEGIN
				DECLARE x INT DEFAULT 0;
				IF x = 0 THEN
					SET x = 1;
				ELSE
					SET x = 2;
				END IF;
			END`,
			kind: StmtCompound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			if res == nil || len(res.Statements) != 1 {
				t.Fatalf("Parse(%q) returned %v, want 1 statement", tc.sql, res)
			}
			ps, ok := res.Statements[0].(*ParsedStmt)
			if !ok {
				t.Fatalf("statement type = %T, want *ParsedStmt", res.Statements[0])
			}
			if ps.Kind != tc.kind {
				t.Errorf("Kind = %s, want %s", ps.Kind, tc.kind)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	bad := []struct {
		name string
		sql  string
	}{
		{"empty", ""},
		{"whitespace only", "   -- comment\n /* block */ "},
		{"unclosed string", "SELECT 'unterminated"},
		{"unclosed backtick", "SELECT `unterminated"},
		{"unclosed block comment", "SELECT /* never closed 1"},
		{"malformed empty hint", "SELECT /*+ */ 1"},
		{"malformed hint args", "SELECT /*+ BROADCAST(t */ 1"},
		{"trailing garbage", "SELECT 1 FROM t )"},
		{"incomplete select", "SELECT"},
		{"incomplete create index", "CREATE INDEX idx ON tbl("},
		{"incomplete set catalog", "SET CATALOG"},
		{"incomplete set variable", "SET VARIABLE (a = (SELECT 1)"},
		{"dangling comma in partition by", "FROM f(TABLE t PARTITION BY a,)"},
		{"unquoted hyphenated create table", "CREATE TABLE my-table (a INT)"},
		{"unquoted hyphenated select from", "SELECT * FROM my-table"},
		{"unpivot incomplete single alias", "SELECT * FROM t UNPIVOT (v FOR k IN (a AS))"},
		{"unpivot incomplete multi alias", "SELECT * FROM t UNPIVOT ((v1, v2) FOR k IN ((a, b) AS))"},
		{"pipe where with window", "FROM t |> WHERE x > 0 WINDOW w AS (PARTITION BY x)"},
		{"pipe pivot with unpivot", "FROM t |> PIVOT(SUM(a) FOR b IN (1)) UNPIVOT(v FOR k IN (c))"},
		{"pipe set dotted target", "FROM t |> SET a.b = 1"},
	}

	for _, tc := range bad {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) unexpectedly succeeded: %+v", tc.sql, res)
			}
		})
	}
}

func TestNestingLimit(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{
			name: "deep parens",
			sql:  "SELECT " + strings.Repeat("(", maxNestingDepth+50) + "1" + strings.Repeat(")", maxNestingDepth+50),
		},
		{
			name: "deep parenthesized relation",
			sql:  "SELECT * FROM " + strings.Repeat("(", maxNestingDepth+50) + "t" + strings.Repeat(")", maxNestingDepth+50),
		},
		{
			name: "deep grouping sets",
			sql:  "SELECT a FROM t GROUP BY " + strings.Repeat("GROUPING SETS(", maxNestingDepth+50) + "a" + strings.Repeat(")", maxNestingDepth+50),
		},
		{
			name: "deep unary minus",
			sql:  "SELECT " + strings.Repeat("-", maxNestingDepth+50) + "1",
		},
		{
			name: "deep unary not",
			sql:  "SELECT " + strings.Repeat("NOT ", maxNestingDepth+50) + "a",
		},
		{
			name: "deep unary bang",
			sql:  "SELECT " + strings.Repeat("!", maxNestingDepth+50) + "a",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.sql); err == nil {
				t.Fatalf("expected nesting depth error for %s", tc.name)
			}
		})
	}
}

func TestDeepNestedLinearParsing(t *testing.T) {
	// Deeply nested TRIM, OVERLAY, POSITION, and row constructors must parse in
	// linear time without exponential backtracking.
	const depth = 50
	cases := []struct {
		name string
		sql  string
	}{
		{
			name: "nested trim",
			sql:  "SELECT " + strings.Repeat("TRIM(", depth) + "x" + strings.Repeat(")", depth),
		},
		{
			name: "nested overlay",
			sql:  "SELECT " + strings.Repeat("OVERLAY(", depth) + "'abc'" + strings.Repeat(", 'x', 1, 1)", depth),
		},
		{
			name: "nested position",
			sql:  "SELECT " + strings.Repeat("POSITION(", depth) + "'a'" + strings.Repeat(", 'abc')", depth),
		},
		{
			name: "nested row constructor",
			sql:  "SELECT " + strings.Repeat("(", depth) + "1" + strings.Repeat(", 2)", depth),
		},
		{
			name: "nested identifier",
			sql:  "SELECT " + strings.Repeat("IDENTIFIER(", depth) + "'x'" + strings.Repeat(")", depth),
		},
		{
			name: "nested datediff",
			sql:  "SELECT " + strings.Repeat("DATEDIFF('2026-01-01', ", depth) + "d" + strings.Repeat(")", depth),
		},
		{
			name: "nested struct with named arg",
			sql:  "SELECT " + strings.Repeat("STRUCT(", depth) + "1" + strings.Repeat(", k => 1)", depth),
		},
		{
			name: "nested case function",
			sql:  "SELECT " + strings.Repeat("CASE(", depth) + "x" + strings.Repeat(")", depth),
		},
		{
			name: "nested select transform subquery",
			sql:  strings.Repeat("SELECT TRANSFORM((", 30) + "SELECT array(1)" + strings.Repeat("), x -> x)", 30),
		},
		{
			name: "nested from select transform subquery",
			sql:  "FROM t " + strings.Repeat("SELECT TRANSFORM((", 25) + "SELECT array(1)" + strings.Repeat("), x -> x)", 25),
		},
		{
			name: "repeated case identifier arithmetic",
			sql:  "SELECT " + strings.Repeat("case + ", 200) + "1",
		},
		{
			name: "repeated case identifier arithmetic inside simple case",
			sql:  "SELECT CASE " + strings.Repeat("case + ", 200) + "1 WHEN 1 THEN 'ok' END",
		},
		{
			name: "nested parenthesized map calls",
			sql:  "SELECT " + strings.Repeat("(MAP(1, ", 40) + "1" + strings.Repeat("))", 40),
		},
		{
			name: "nested parenthesized reduce calls",
			sql:  "SELECT " + strings.Repeat("(REDUCE(a, ", 40) + "0" + strings.Repeat(", (x, y) -> x + y))", 40),
		},
		{
			name: "nested parenthesized values over calls",
			sql:  "SELECT " + strings.Repeat("(VALUES(", 30) + "1" + strings.Repeat(") OVER ())", 30),
		},
		{
			name: "nested parenthesized values filter calls",
			sql:  "SELECT " + strings.Repeat("(VALUES(", 30) + "1" + strings.Repeat(") FILTER (WHERE true))", 30),
		},
		{
			name: "nested parenthesized values aliased row constructors",
			sql:  "SELECT " + strings.Repeat("(VALUES(", 30) + "1" + strings.Repeat(") AS a, 2)", 30),
		},
		{
			name: "nested values table function with named arg",
			sql:  strings.Repeat("SELECT * FROM VALUES((", 25) + "SELECT 1" + strings.Repeat("), k => 1)", 25),
		},
		{
			name: "nested table function call arg",
			sql:  strings.Repeat("SELECT * FROM f(TABLE((", 25) + "SELECT 1" + strings.Repeat("), k => 1))", 25),
		},
		{
			name: "nested named table function call arg",
			sql:  strings.Repeat("SELECT * FROM f(arg => TABLE((", 25) + "SELECT 1" + strings.Repeat("), k => 1))", 25),
		},
		{
			name: "nested parens",
			sql:  "SELECT " + strings.Repeat("(", depth) + "1" + strings.Repeat(")", depth),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.sql); err != nil {
				t.Fatalf("Parse %s failed: %v", tc.name, err)
			}
		})
	}
}

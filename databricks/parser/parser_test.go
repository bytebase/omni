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
			sql:  `SELECT first(x) OVER (ORDER BY x), last(x) OVER (), any_value(x) FILTER (WHERE x > 0), first(x) IGNORE NULLS OVER (PARTITION BY k), first(x IGNORE NULLS) FROM t`,
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
			sql:  `SELECT * FROM lateral JOIN lateral.t ON lateral.id = t.id JOIN lateral ON lateral.id = t.id`,
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
			name: "table argument partition and order by with parenthesized expression continuations",
			sql:  `SELECT * FROM f(TABLE(t) PARTITION BY (x) + 1 ORDER BY (y) + 1 DESC)`,
			kind: StmtQuery,
		},
		{
			name: "special-named function calls with ALL and DISTINCT set quantifiers",
			sql:  `SELECT trim(ALL x), substring(ALL x, 1, 2), overlay(ALL x, 'a', 1, 1), position(ALL 'a', x), extract(ALL x), cast(ALL x), struct(ALL x), datediff(ALL d1, d2) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "dynamic identifier table with temporal options and sample suffixes",
			sql:  `SELECT * FROM IDENTIFIER('t') VERSION AS OF 1 WITH ('x'='y') TABLESAMPLE (10 PERCENT) AS s`,
			kind: StmtQuery,
		},
		{
			name: "set backquoted config key with numeric boolean and string values",
			sql:  "SET `spark.sql.shuffle.partitions` = 8",
			kind: StmtSet,
		},
		{
			name: "set backquoted config key with boolean value",
			sql:  "SET `spark.sql.ansi.enabled` = true",
			kind: StmtSet,
		},
		{
			name: "qualified star on multi-part identifier",
			sql:  `SELECT t.*, t.x.*, db.t.* EXCEPT (a) FROM db.t`,
			kind: StmtQuery,
		},
		{
			name: "implicit filter and over aliases before FROM",
			sql:  `SELECT x filter FROM t UNION ALL SELECT x over FROM t UNION ALL SELECT count(*) filter FROM t UNION ALL SELECT count(*) over FROM t`,
			kind: StmtQuery,
		},
		{
			name: "non-reserved keywords as named window references",
			sql:  `SELECT count(*) OVER window, count(*) OVER order, count(*) OVER from, count(*) OVER group, count(*) OVER limit, count(*) OVER window c FROM t WINDOW window AS (), order AS (), from AS (), group AS (), limit AS ()`,
			kind: StmtQuery,
		},
		{
			name: "implicit over alias before window and order clauses",
			sql:  `SELECT count(*) over WINDOW w AS () UNION ALL SELECT count(*) over ORDER BY 1`,
			kind: StmtQuery,
		},
		{
			name: "ordinary named arguments on special-named functions",
			sql:  `SELECT CAST(k => 1), POSITION(a => 'x', b => 'xy'), SUBSTRING(str => s, pos => 1), OVERLAY(a => s, b => 'x', c => 1) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "parenthesized boolean operands in SQL-specific functions and unparenthesized in ordinary calls",
			sql:  `SELECT SUBSTRING((a AND b) FROM 1), TRIM((a AND b) FROM s), POSITION((a AND b) IN s), OVERLAY((a AND b) PLACING 'x' FROM 1), SUBSTRING(a AND b, 1, 2), TRIM(a AND b, s), POSITION(a AND b, s), OVERLAY(a AND b, 'x', 1) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "position IN with parenthesized valueExpression continuations and ordinary IN predicates",
			sql:  `SELECT POSITION('a' IN ('abc') || 'def'), POSITION('a' IN ('abc')::STRING), POSITION('a' IN ('a', 'b')[0]), POSITION(a IN (1, 2)), POSITION(a IN (1, 2), s), POSITION(a IN (1, 2) AND b, s), POSITION(a IN (1, 2)) OVER () FROM t`,
			kind: StmtQuery,
		},
		{
			name: "function table partition and order clauses followed by named and table args",
			sql:  `SELECT * FROM f(TABLE(t) PARTITION BY (a, b) ORDER BY c DESC, k => 1, TABLE(t2) PARTITION BY d)`,
			kind: StmtQuery,
		},
		{
			name: "parenthesized TABLE queries with multi-item clauses and set/pipe continuations",
			sql:  `SELECT (TABLE t ORDER BY a, b), (TABLE db.t UNION SELECT a, b FROM s), (TABLE IDENTIFIER('t') ORDER BY a, b), (TABLE t |> SELECT a, b), ((TABLE t ORDER BY a, b) LIMIT 1), ((TABLE t ORDER BY a, b), 1) FROM (TABLE t ORDER BY a, b) AS u`,
			kind: StmtQuery,
		},
		{
			name: "parenthesized TABLE query inside function table arg",
			sql:  `SELECT * FROM f(TABLE(TABLE t ORDER BY a, b))`,
			kind: StmtQuery,
		},
		{
			name: "parenthesized table identifier in row constructors, expressions, and relations",
			sql:  `SELECT (table a, b), (table AS a, b), (table AND x), (table AND limit), (table AND limit = 1), (table OR offset > 0), (table OR union = 1), (table IS NULL) FROM (table a JOIN t2 ON a.id = t2.id), (table AS t3)`,
			kind: StmtQuery,
		},
		{
			name: "describe relation named table",
			sql:  `DESCRIBE table`,
			kind: StmtDescribe,
		},
		{
			name: "describe qualified relation starting with table",
			sql:  `DESCRIBE table.x`,
			kind: StmtDescribe,
		},
		{
			name: "describe relation named table with partition and as json",
			sql:  `DESCRIBE EXTENDED table PARTITION (p = 1) AS JSON`,
			kind: StmtDescribe,
		},
		{
			name: "describe table keyword with relation named table",
			sql:  `DESCRIBE TABLE table.x`,
			kind: StmtDescribe,
		},
		{
			name: "describe table query with order by",
			sql:  `DESCRIBE TABLE t ORDER BY a, b`,
			kind: StmtDescribe,
		},
		{
			name: "lateral view function named outer",
			sql:  `SELECT * FROM t LATERAL VIEW outer(x) v AS x`,
			kind: StmtQuery,
		},
		{
			name: "lateral view qualified function starting with outer",
			sql:  `SELECT * FROM t LATERAL VIEW outer.f(x) v AS x`,
			kind: StmtQuery,
		},
		{
			name: "lateral view outer modifier with function named outer",
			sql:  `SELECT * FROM t LATERAL VIEW OUTER outer(x) v AS x LATERAL VIEW OUTER outer.f(x) v2 AS y`,
			kind: StmtQuery,
		},
		{
			name: "alter table alter column named column",
			sql:  `ALTER TABLE t ALTER column TYPE INT`,
			kind: StmtAlterTable,
		},
		{
			name: "alter table change column named column",
			sql:  `ALTER TABLE t CHANGE column column INT`,
			kind: StmtAlterTable,
		},
		{
			name: "ordering and frame keywords as implicit column aliases",
			sql:  `SELECT x rows FROM t UNION ALL SELECT x nulls FROM t UNION ALL SELECT x asc FROM t UNION ALL SELECT x desc, y range, z first, w last, u preceding, v following FROM t`,
			kind: StmtQuery,
		},
		{
			name: "struct with field aliases and struct with named arguments separately",
			sql:  `SELECT STRUCT(x AS a, y b), STRUCT(1, k => 2), STRUCT(k => 1, m => 2) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "tablesample bucket without on clause",
			sql:  "SELECT * FROM t TABLESAMPLE (BUCKET 1 OUT OF 2)",
			kind: StmtQuery,
		},
		{
			name: "table function named lateral vs lateral modifier before subquery or tvf",
			sql:  `SELECT * FROM lateral(1) JOIN lateral(1, 2) JOIN lateral() JOIN lateral(k => 1) JOIN lateral(TABLE(t) PARTITION BY a) JOIN LATERAL (SELECT 1) AS s JOIN LATERAL lateral(1) AS l`,
			kind: StmtQuery,
		},
		{
			name: "unquoted values table before where, join on, temporal, options, and pivot clauses",
			sql:  `SELECT * FROM values WHERE x = 1 UNION ALL SELECT * FROM t JOIN values ON t.x = values.x UNION ALL SELECT * FROM values VERSION AS OF 1 UNION ALL SELECT * FROM values FOR SYSTEM_VERSION AS OF 1 UNION ALL SELECT * FROM values WITH (k = 'v')`,
			kind: StmtQuery,
		},
		{
			name: "non-reserved clause keywords inside simple CASE operands",
			sql:  `SELECT CASE from WHEN 1 THEN 2 END, CASE order WHEN 1 THEN 2 END, CASE x + group WHEN 1 THEN 2 END, CASE from + 1 WHEN 1 THEN 2 END, CASE limit WHEN 1 THEN 2 END, CASE end WHEN 1 THEN 2 END FROM t WHERE CASE x WHEN 1 THEN true END`,
			kind: StmtQuery,
		},
		{
			name: "case column before FROM with CASE in WHERE",
			sql:  `SELECT case FROM t WHERE CASE x WHEN 1 THEN true END`,
			kind: StmtQuery,
		},
		{
			name: "non-reserved keywords as implicit table aliases before clauses",
			sql:  `SELECT * FROM t window WHERE x = 1 UNION ALL SELECT * FROM t order WHERE x = 1 UNION ALL SELECT * FROM t group WHERE x = 1 UNION ALL SELECT * FROM t pivot WHERE x = 1 UNION ALL SELECT * FROM t unpivot WHERE x = 1 UNION ALL SELECT * FROM t tablesample WHERE x = 1 UNION ALL SELECT * FROM t limit ORDER BY 1`,
			kind: StmtQuery,
		},
		{
			name: "trim with from column expression and SQL FROM separator forms",
			sql:  `SELECT TRIM(from, s), TRIM(from || x), TRIM(from.x), TRIM(from FROM s), TRIM(BOTH from FROM s), TRIM(BOTH from(s) FROM t), TRIM(BOTH from - 1 FROM t), TRIM(BOTH FROM (s)), TRIM(LEADING FROM (s)), TRIM(BOTH FROM -s), TRIM(FROM s), TRIM(FROM from), TRIM(FROM from || 'x'), TRIM(from(s), 'a'), TRIM(FROM (SELECT 'a')) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "parenthesized expressions and row constructors led by from, select, and with",
			sql:  `SELECT (from + 1), (from.x), (select, x), (select.x), (select * 2), (select AND x), (with a, b), (with AND x), (from(1) + 1) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "unary NOT before non-reserved clause keywords",
			sql:  `SELECT NOT order, NOT from, NOT where, NOT limit, NOT end, NOT as, NOT union FROM t WHERE NOT from.x AND NOT order AND NOT group`,
			kind: StmtQuery,
		},
		{
			name: "declare variables named var or variable with initializers and types",
			sql:  `DECLARE var = 1`,
			kind: StmtDeclare,
		},
		{
			name: "declare variable named variable with DEFAULT",
			sql:  `DECLARE variable DEFAULT 1`,
			kind: StmtDeclare,
		},
		{
			name: "declare variable named var with complex type",
			sql:  `DECLARE var VARCHAR(10) = 'a'`,
			kind: StmtDeclare,
		},
		{
			name: "declare var keyword with variable named default",
			sql:  `DECLARE VAR default = 1`,
			kind: StmtDeclare,
		},
		{
			name: "simple CASE with unquoted case column operand and nested searched CASE",
			sql:  `SELECT CASE case WHEN 1 THEN 2 END, CASE CASE WHEN true THEN 1 END WHEN 1 THEN 2 END, CASE CASE case WHEN 1 THEN 2 END WHEN 2 THEN 3 END, CASE x WHEN 1 THEN case WHEN 2 THEN 3 END FROM t`,
			kind: StmtQuery,
		},
		{
			name: "first, last, and any_value with all/distinct operands before IGNORE NULLS",
			sql:  `SELECT FIRST(all IGNORE NULLS), LAST(distinct IGNORE NULLS), ANY_VALUE(all IGNORE NULLS), FIRST(DISTINCT x), FIRST(DISTINCT ignore) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "grant with keyword tail",
			sql:  `GRANT SELECT ON TABLE t TO u`,
			kind: StmtHiveNative,
		},
		{
			name: "revoke all privileges with keyword tail",
			sql:  `REVOKE ALL PRIVILEGES ON SCHEMA s FROM u`,
			kind: StmtHiveNative,
		},
		{
			name: "set catalog as bare config key",
			sql:  `SET CATALOG`,
			kind: StmtSet,
		},
		{
			name: "case-sensitive CTE names in WITH INSERT",
			sql:  `WITH c AS (SELECT 1), C AS (SELECT 2) INSERT INTO t SELECT * FROM c`,
			kind: StmtInsert,
		},
		{
			name: "case-sensitive tblproperties keys",
			sql:  `CREATE TABLE t (x INT) TBLPROPERTIES ('k' = 'a', 'K' = 'b')`,
			kind: StmtCreateTable,
		},
		{
			name: "create function with valid resource types",
			sql:  `CREATE FUNCTION f AS 'com.example.F' USING JAR 's3://b/f.jar', FILE 's3://b/a.txt', ARCHIVE 's3://b/a.zip'`,
			kind: StmtCreateFunction,
		},
		{
			name: "create view with literal question mark string",
			sql:  `CREATE VIEW v AS SELECT '?' AS x`,
			kind: StmtCreateView,
		},
		{
			name: "scripting labels using ELSE, WHEN, END, UNTIL",
			sql:  `BEGIN END: LOOP SELECT 1; END LOOP; END`,
			kind: StmtCompound,
		},
		{
			name: "from query followed by select having column",
			sql:  `FROM a SELECT HAVING a`,
			kind: StmtQuery,
		},
		{
			name: "keyword-as-identifier edge cases matching Spark 4.0.5",
			sql:  `(SELECT 1 AS FROM) UNION ALL (SELECT ALL()) UNION ALL (SELECT NOT()) UNION ALL (SELECT DISTINCT()) UNION ALL (SELECT FALSE.*) UNION ALL (SELECT INTERVAL + zzz a) UNION ALL (SELECT ~ zzz ELSE) UNION ALL (SELECT ~ zzz WHEN) UNION ALL (SELECT ~ zzz END) UNION ALL (SELECT ~ zzz UNTIL)`,
			kind: StmtQuery,
		},
		{
			name: "chained unary operators before minimum-width negative literals",
			sql:  `SELECT +-128Y, ~-128Y, +-32768S, +-9223372036854775808L`,
			kind: StmtQuery,
		},
		{
			name: "timestamp literals with valid timezone suffixes",
			sql:  `SELECT TIMESTAMP '1997-01-31 09:26:56.66666666UTC+08:00', TIMESTAMP '2020-01-01 12:00:00Z', TIMESTAMP '2020-01-01 12:00:00+08:30:15', TIMESTAMP '2020-01-01 12:00:00 America/Los_Angeles', TIMESTAMP '2020-01-01 12:00:00 PST'`,
			kind: StmtQuery,
		},
		{
			name: "backticked dot in temporary view name",
			sql:  "CREATE TEMP VIEW `a.b` AS SELECT 1",
			kind: StmtCreateView,
		},
		{
			name: "backticked dot in temporary view using",
			sql:  "CREATE TEMP VIEW `a.b` USING parquet",
			kind: StmtCreateView,
		},
		{
			name: "backticked dot in temporary function name",
			sql:  "CREATE TEMP FUNCTION `a.b`() RETURNS INT RETURN 1",
			kind: StmtCreateFunction,
		},
		{
			name: "backticked dot in scripting label",
			sql:  "BEGIN `a.b`: LOOP LEAVE `a.b`; END LOOP `a.b`; END",
			kind: StmtCompound,
		},
		{
			name: "valid ordered interval data types",
			sql:  `CREATE TABLE t (a INTERVAL YEAR TO MONTH, b INTERVAL DAY TO HOUR, c INTERVAL DAY TO SECOND, d INTERVAL HOUR TO SECOND, e INTERVAL MINUTE TO SECOND)`,
			kind: StmtCreateTable,
		},
		{
			name: "explicit aliases named after non-reserved clause keywords",
			sql:  `(SELECT 1 AS limit) UNION ALL (SELECT 1 AS where FROM t) UNION ALL (SELECT 1 AS limit LIMIT 1) UNION ALL (SELECT 1 AS WHERE a = 1) UNION ALL (SELECT * FROM t AS limit) UNION ALL (SELECT * FROM t AS where WHERE a = 1) UNION ALL (SELECT * FROM t AS WHERE a = 1) UNION ALL (SELECT * FROM t AS JOIN s ON t.a = s.a)`,
			kind: StmtQuery,
		},
		{
			name: "built-in spark error condition in exception handler",
			sql:  `BEGIN DECLARE EXIT HANDLER FOR DIVIDE_BY_ZERO SELECT 0; SELECT 1 / 0; END`,
			kind: StmtCompound,
		},
		{
			name: "built-in spark error sub-condition in exception handler",
			sql:  `BEGIN DECLARE EXIT HANDLER FOR INVALID_SQL_SYNTAX.LATERAL_WITHOUT_SUBQUERY_OR_TABLE_VALUED_FUNC SELECT 0; SELECT 1; END`,
			kind: StmtCompound,
		},
		{
			name: "execute immediate using case-sensitive distinct argument aliases",
			sql:  `EXECUTE IMMEDIATE 'SELECT :x + :X' USING 1 AS x, 2 AS X`,
			kind: StmtExecuteImmediate,
		},
		{
			name: "500 nested parentheses succeed within 10000 nesting limit",
			sql:  "SELECT " + strings.Repeat("(", 500) + "1" + strings.Repeat(")", 500),
			kind: StmtQuery,
		},
		{
			name: "pipe aggregate with keyword-led expressions",
			sql:  `FROM t |> AGGREGATE order + sum(x), window, limit + 1 GROUP BY order, window, limit`,
			kind: StmtQuery,
		},
		{
			name: "pipe aggregate followed by query organization",
			sql:  `FROM t |> AGGREGATE sum(x) ORDER BY 1 LIMIT 10`,
			kind: StmtQuery,
		},
		{
			name: "undotted SET key with unquoted value",
			sql:  `SET foo = bar`,
			kind: StmtSet,
		},
		{
			name: "undotted SET compression with unquoted value",
			sql:  `SET compression = snappy`,
			kind: StmtSet,
		},
		{
			name: "insert into TABLE keyword with keyword-named targets",
			sql:  `INSERT INTO TABLE select SELECT 1`,
			kind: StmtInsert,
		},
		{
			name: "insert overwrite TABLE keyword with from target",
			sql:  `INSERT OVERWRITE TABLE from SELECT 1`,
			kind: StmtInsert,
		},
		{
			name: "insert into TABLE keyword with select target and column list",
			sql:  `INSERT INTO TABLE select (a, b) SELECT 1, 2`,
			kind: StmtInsert,
		},
		{
			name: "insert into TABLE keyword with with target and column list",
			sql:  `INSERT INTO TABLE with (a, b) SELECT 1, 2`,
			kind: StmtInsert,
		},
		{
			name: "insert into table named table with options clause",
			sql:  `INSERT INTO table WITH (k = 'v') SELECT 1`,
			kind: StmtInsert,
		},
		{
			name: "insert into TABLE keyword with parenthesized query",
			sql:  `INSERT INTO TABLE select (SELECT 1)`,
			kind: StmtInsert,
		},
		{
			name: "insert into table named table with inline TABLE query",
			sql:  `INSERT INTO table TABLE t`,
			kind: StmtInsert,
		},
		{
			name: "insert into table named table selecting keyword columns",
			sql:  `INSERT INTO table SELECT select, from FROM t`,
			kind: StmtInsert,
		},
		{
			name: "multi-insert with TABLE keyword and keyword target",
			sql:  `FROM s INSERT INTO TABLE from SELECT 1 INSERT INTO table SELECT from`,
			kind: StmtMultiInsert,
		},
		{
			name: "rollup with subquery and nested parenthesized row constructor",
			sql:  `SELECT 1 FROM t GROUP BY ROLLUP((SELECT 1), a), ROLLUP((((SELECT 1) AS a, 1)))`,
			kind: StmtQuery,
		},
		{
			name: "group by expression with grouping sets",
			sql:  `SELECT a FROM t GROUP BY a, GROUPING SETS ((b), ())`,
			kind: StmtQuery,
		},
		{
			name: "group by grouping sets with expression",
			sql:  `SELECT a FROM t GROUP BY GROUPING SETS ((b)), a`,
			kind: StmtQuery,
		},
		{
			name: "dynamic relation identifier with table column expression and function args",
			sql:  `SELECT * FROM IDENTIFIER(table)(1)`,
			kind: StmtQuery,
		},
		{
			name: "dynamic relation identifier with table qualified expression and function args",
			sql:  `SELECT * FROM IDENTIFIER(table.x)(1)`,
			kind: StmtQuery,
		},
		{
			name: "dynamic relation identifier with table column expression and temporal clause",
			sql:  `SELECT * FROM IDENTIFIER(table) VERSION AS OF 1`,
			kind: StmtQuery,
		},
		{
			name: "dynamic relation identifier with zero-arg table call and function args",
			sql:  `SELECT * FROM IDENTIFIER(table())(1)`,
			kind: StmtQuery,
		},
		{
			name: "table function argument with parenthesized partitioning and subsequent args",
			sql:  `SELECT * FROM f(TABLE(t) PARTITION BY (a, b) ORDER BY (a DESC, b ASC), c, k => 1, TABLE(u) PARTITION BY x, TABLE(SELECT 1))`,
			kind: StmtQuery,
		},
		{
			name: "merge with conditional non-final clauses and unconditional final clauses",
			sql:  `MERGE INTO t AS dst USING s AS src ON dst.id = src.id WHEN MATCHED AND src.del THEN DELETE WHEN MATCHED THEN UPDATE SET dst.v = src.v WHEN NOT MATCHED AND src.ok THEN INSERT * WHEN NOT MATCHED THEN INSERT (id) VALUES (src.id) WHEN NOT MATCHED BY SOURCE AND dst.old THEN DELETE WHEN NOT MATCHED BY SOURCE THEN UPDATE SET dst.v = 0`,
			kind: StmtMerge,
		},
		{
			name: "non-reserved clause keywords as implicit column aliases before FROM and clauses",
			sql:  `SELECT x window FROM t UNION ALL SELECT x order FROM t UNION ALL SELECT x group FROM t UNION ALL SELECT x limit FROM t UNION ALL SELECT x offset FROM t UNION ALL SELECT x where FROM t UNION ALL SELECT x having FROM t UNION ALL SELECT x from FROM t UNION ALL (SELECT x limit LIMIT 1) UNION ALL (SELECT x order ORDER BY 1) UNION ALL (SELECT x window WINDOW w AS ()) UNION ALL (SELECT x LIMIT from)`,
			kind: StmtQuery,
		},
		{
			name: "limit all vs limit expressions starting with all",
			sql:  `(SELECT 1 LIMIT ALL) UNION ALL (SELECT 1 LIMIT ALL OFFSET 2) UNION ALL (SELECT 1 LIMIT all + 1) UNION ALL (SELECT 1 LIMIT all.x) UNION ALL (SELECT 1 LIMIT all(1)) UNION ALL (SELECT 1 LIMIT all * 2 OFFSET 1)`,
			kind: StmtQuery,
		},
		{
			name: "create table named if",
			sql:  `CREATE TABLE if (x INT)`,
			kind: StmtCreateTable,
		},
		{
			name: "create schema named if",
			sql:  `CREATE SCHEMA if`,
			kind: StmtCreateNamespace,
		},
		{
			name: "create function named if",
			sql:  `CREATE FUNCTION if AS 'x'`,
			kind: StmtCreateFunction,
		},
		{
			name: "create view named if",
			sql:  `CREATE VIEW if AS SELECT 1`,
			kind: StmtCreateView,
		},
		{
			name: "uncache table named if",
			sql:  `UNCACHE TABLE if`,
			kind: StmtUncache,
		},
		{
			name: "backtick-quoted unicode identifiers",
			sql:  "CREATE TABLE `café` (`naïve` INT)",
			kind: StmtCreateTable,
		},
		{
			name: "nested WITH clause shadowing outer CTE name",
			sql:  `WITH c AS (WITH c AS (SELECT 1) SELECT * FROM c) SELECT * FROM c`,
			kind: StmtQuery,
		},
		{
			name: "property setters with values and unsetters with keys only",
			sql:  `ALTER TABLE t UNSET TBLPROPERTIES IF EXISTS ('k1', 'k2')`,
			kind: StmtAlterTable,
		},
		{
			name: "create table named identifier with column list",
			sql:  `CREATE TABLE identifier (a INT)`,
			kind: StmtCreateTable,
		},
		{
			name: "create function named identifier with parameter list",
			sql:  `CREATE FUNCTION identifier(x INT) RETURNS INT RETURN x`,
			kind: StmtCreateFunction,
		},
		{
			name: "create index on table named identifier",
			sql:  `CREATE INDEX idx ON identifier (a)`,
			kind: StmtCreateIndex,
		},
		{
			name: "call procedure named identifier",
			sql:  `CALL identifier(1)`,
			kind: StmtCall,
		},
		{
			name: "keyword-named implicit table aliases before WHERE, LIMIT, and JOIN",
			sql:  `SELECT * FROM t limit WHERE x = 1 UNION ALL SELECT * FROM t from WHERE x = 1 UNION ALL SELECT * FROM t where WHERE x = 1 UNION ALL SELECT * FROM t WHERE where = 1 UNION ALL SELECT * FROM t limit JOIN s ON limit.a = s.a UNION ALL (SELECT * FROM t from LIMIT 1) UNION ALL (SELECT * FROM t LIMIT where)`,
			kind: StmtQuery,
		},
		{
			name: "all and distinct arithmetic operands in SQL-specific function forms",
			sql:  `SELECT CAST(all + 1 AS INT), TRY_CAST(distinct - 1 AS INT), POSITION(distinct - 1 IN s), SUBSTRING(all + 1 FROM 2), SUBSTRING(all + 1, 2 FOR 3), TRIM(all + 1 FROM s), OVERLAY(all + 1 PLACING 'x' FROM 1), FIRST(all + 1 IGNORE NULLS), STRUCT(all + 1 AS x) FROM t`,
			kind: StmtQuery,
		},
		{
			name: "supported typed literal constructors",
			sql:  `SELECT DATE '2026-01-01', TIMESTAMP '2026-01-01 00:00:00', TIMESTAMP_NTZ '2026-01-01', TIMESTAMP_LTZ '2026-01-01', INTERVAL '1 day', X'deadbeef', X 'deadbeef'`,
			kind: StmtQuery,
		},
		{
			name: "query with distribute by and sort by",
			sql:  `SELECT x FROM t DISTRIBUTE BY x SORT BY x DESC`,
			kind: StmtQuery,
		},
		{
			name: "analyze table and tables with NOSCAN modifier",
			sql:  `ANALYZE TABLE t COMPUTE STATISTICS NOSCAN`,
			kind: StmtAnalyze,
		},
		{
			name: "analyze tables with lowercase noscan modifier",
			sql:  `ANALYZE TABLES COMPUTE STATISTICS noscan`,
			kind: StmtAnalyze,
		},
		{
			name: "where clause with keyword-named operands before group by",
			sql:  `SELECT * FROM t WHERE having GROUP BY y UNION ALL SELECT * FROM t WHERE limit GROUP BY y UNION ALL SELECT * FROM t WHERE offset ORDER BY y`,
			kind: StmtQuery,
		},
		{
			name: "compound keyword-led operands in simple CASE",
			sql:  `SELECT CASE and + 1 WHEN 1 THEN 2 END, CASE collate + 1 WHEN 1 THEN 2 END, CASE and(1) WHEN 1 THEN 2 END, CASE end + 1 WHEN 1 THEN 2 END FROM t`,
			kind: StmtQuery,
		},
		{
			name: "not column before sort and window frame modifiers",
			sql:  `SELECT sum(x) OVER (ORDER BY not ROWS 1 PRECEDING), sum(x) OVER (ORDER BY not RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t ORDER BY not NULLS FIRST`,
			kind: StmtQuery,
		},
		{
			name: "valid primitive type parameter arities",
			sql:  `CREATE TABLE t (a CHAR(10), b VARCHAR(20), c DECIMAL, d DECIMAL(10), e DECIMAL(10, 2))`,
			kind: StmtCreateTable,
		},
		{
			name: "insert overwrite static partition with if not exists",
			sql:  `INSERT OVERWRITE t PARTITION (p = 1) IF NOT EXISTS SELECT 1`,
			kind: StmtInsert,
		},
		{
			name: "ordinary natural join and cross join lateral",
			sql:  `SELECT * FROM t NATURAL LEFT JOIN s CROSS JOIN LATERAL (SELECT 1 AS x) l`,
			kind: StmtQuery,
		},
		{
			name: "set raw value with unrecognized characters",
			sql:  `SET spark.test = abc@def`,
			kind: StmtSet,
		},
		{
			name: "set raw value with variable interpolation syntax",
			sql:  `SET spark.test = ${env:FOO}`,
			kind: StmtSet,
		},
		{
			name: "double parenthesized row constructor with subquery",
			sql:  `SELECT (((SELECT 1), 1))`,
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
		{"set role all rejected by spark", "SET ROLE ALL"},
		{"incomplete set variable", "SET VARIABLE (a = (SELECT 1)"},
		{"dangling comma in partition by", "FROM f(TABLE t PARTITION BY a,)"},
		{"unquoted hyphenated create table", "CREATE TABLE my-table (a INT)"},
		{"unquoted hyphenated select from", "SELECT * FROM my-table"},
		{"unpivot incomplete single alias", "SELECT * FROM t UNPIVOT (v FOR k IN (a AS))"},
		{"unpivot incomplete multi alias", "SELECT * FROM t UNPIVOT ((v1, v2) FOR k IN ((a, b) AS))"},
		{"pipe where with window", "FROM t |> WHERE x > 0 WINDOW w AS (PARTITION BY x)"},
		{"pipe pivot with unpivot", "FROM t |> PIVOT(SUM(a) FOR b IN (1)) UNPIVOT(v FOR k IN (c))"},
		{"group by alias then grouping sets", "SELECT a FROM t GROUP BY a AS x, GROUPING SETS ((b))"},
		{"group by grouping sets then alias", "SELECT a FROM t GROUP BY GROUPING SETS ((b)), a AS x"},
		{"group by implicit alias then grouping sets", "SELECT a FROM t GROUP BY a x, GROUPING SETS ((b))"},
		{"group by grouping sets then implicit alias", "SELECT a FROM t GROUP BY GROUPING SETS ((b)), a x"},
		{"group by alias then rollup with empty set", "SELECT a FROM t GROUP BY a AS x, ROLLUP(())"},
		{"group by rollup with empty set then alias", "SELECT a FROM t GROUP BY ROLLUP(()), a AS x"},
		{"group by rollup with empty set then with rollup", "SELECT a FROM t GROUP BY ROLLUP(()) WITH ROLLUP"},
		{"table function argument missing parentheses", "SELECT * FROM f(TABLE t)"},
		{"table function argument unparenthesized multi partition by", "SELECT * FROM f(TABLE(t) PARTITION BY a, b)"},
		{"table function argument unparenthesized multi order by", "SELECT * FROM f(TABLE(t) WITH SINGLE PARTITION ORDER BY a, b)"},
		{"merge non-final unconditional matched clause", "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE WHEN MATCHED THEN UPDATE SET x = 1"},
		{"merge non-final unconditional not matched clause", "MERGE INTO t USING s ON t.id = s.id WHEN NOT MATCHED THEN INSERT * WHEN NOT MATCHED AND s.x > 0 THEN INSERT *"},
		{"merge non-final unconditional not matched by source clause", "MERGE INTO t USING s ON t.id = s.id WHEN NOT MATCHED BY SOURCE THEN DELETE WHEN NOT MATCHED BY SOURCE THEN UPDATE SET x = 1"},
		{"update table alias with column list", "UPDATE t AS x(a) SET a = 1"},
		{"update implicit table alias with column list", "UPDATE t x(a) SET a = 1"},
		{"delete table alias with column list", "DELETE FROM t AS x(a) WHERE a = 1"},
		{"merge target alias with column list", "MERGE INTO t AS x(a) USING s ON x.a = s.a WHEN MATCHED THEN DELETE"},
		{"merge source alias with column list", "MERGE INTO t USING s AS y(a) ON t.a = y.a WHEN MATCHED THEN DELETE"},
		{"duplicate default column descriptor", "CREATE TABLE t (a INT DEFAULT 1 DEFAULT 2)"},
		{"duplicate not null column descriptor", "CREATE TABLE t (a INT NOT NULL NOT NULL)"},
		{"duplicate comment column descriptor", "CREATE TABLE t (a INT COMMENT 'a' COMMENT 'b')"},
		{"duplicate generated column descriptor", "CREATE TABLE t (a INT GENERATED ALWAYS AS (1) GENERATED ALWAYS AS (2))"},
		{"alter table add columns duplicate default", "ALTER TABLE t ADD COLUMNS (a INT DEFAULT 1 DEFAULT 2)"},
		{"alter table add columns duplicate position", "ALTER TABLE t ADD COLUMNS (a INT FIRST AFTER b)"},
		{"unquoted unicode table name", "CREATE TABLE café (x INT)"},
		{"unquoted unicode column in select", "SELECT café FROM t"},
		{"unsupported primitive type in create table", "CREATE TABLE t (x foobar)"},
		{"unsupported primitive type in cast", "SELECT CAST(1 AS foo)"},
		{"duplicate CTE name in same WITH clause", "WITH c AS (SELECT 1), c AS (SELECT 2) SELECT * FROM c"},
		{"duplicate tblproperties key", "CREATE TABLE t (x INT) TBLPROPERTIES ('k'='a','k'='b')"},
		{"duplicate unset tblproperties key", "ALTER TABLE t UNSET TBLPROPERTIES ('k','k')"},
		{"lateral join using", "SELECT * FROM t JOIN LATERAL (SELECT 1 AS id) s USING (id)"},
		{"right join lateral", "SELECT * FROM t RIGHT JOIN LATERAL (SELECT 1 AS id) s ON t.id = s.id"},
		{"full join lateral", "SELECT * FROM t FULL OUTER JOIN LATERAL (SELECT 1 AS id) s ON t.id = s.id"},
		{"left semi join lateral", "SELECT * FROM t LEFT SEMI JOIN LATERAL (SELECT 1 AS id) s ON t.id = s.id"},
		{"join lateral bare table", "SELECT * FROM t JOIN LATERAL x"},
		{"comma lateral bare table", "SELECT * FROM t, LATERAL x"},
		{"qualified tvf lateral.f", "SELECT * FROM lateral.f(1)"},
		{"invalid sqlstate short", "BEGIN DECLARE bad CONDITION FOR SQLSTATE 'abc'; SELECT 1; END"},
		{"invalid sqlstate 00000", "BEGIN DECLARE bad CONDITION FOR SQLSTATE '00000'; SELECT 1; END"},
		{"invalid sqlstate in handler", "BEGIN DECLARE CONTINUE HANDLER FOR SQLSTATE '01000' SELECT 1; SELECT 1; END"},
		{"invalid function resource type", "CREATE FUNCTION f AS 'com.example.F' USING garbage 's3://bucket/f.jar'"},
		{"describe table as json without extended", "DESCRIBE TABLE t AS JSON"},
		{"describe table partition without value", "DESC SELECT PARTITION(a)"},
		{"create view positional parameter marker", "CREATE VIEW v AS SELECT ? AS x"},
		{"create view named parameter marker", "CREATE OR REPLACE TEMP VIEW v AS SELECT :param AS x"},
		{"declare var boolean missing default", "DECLARE VAR BOOLEAN"},
		{"declare var default unsupported type", "DECLARE VAR DEFAULT a"},
		{"declare var interval unsupported type", "DECLARE VAR INTERVAL YEAR"},
		{"alter table change column rename", "ALTER TABLE t CHANGE column c2 INT"},
		{"tablesample bucket on column", "SELECT * FROM t TABLESAMPLE (BUCKET 1 OUT OF 2 ON a)"},
		{"group by alias", "SELECT a FROM t GROUP BY a AS x"},
		{"group by rollup empty set", "SELECT 1 FROM t GROUP BY ROLLUP(())"},
		{"semi as table alias", "SELECT * FROM t semi WHERE x = 1"},
		{"anti as table alias", "SELECT * FROM t AS anti"},
		{"collate trim unsupported", "SELECT a COLLATE TRIM"},
		{"pipe aggregate group by all", "SELECT a |> AGGREGATE SUM(a) GROUP BY ALL"},
		{"unary minus then plus before 128Y out of range", "SELECT -+128Y"},
		{"timestamp invalid timezone suffix", "SELECT TIMESTAMP '2020-01-01 12:00:00+garbage'"},
		{"timestamp invalid named zone", "SELECT TIMESTAMP '2020-01-01 12:00:00 Invalid/Zone'"},
		{"timestamp out of range timezone offset", "SELECT TIMESTAMP '2020-01-01 12:00:00 UTC+25:00'"},
		{"interval type second to hour reversed", "CREATE TABLE t (i INTERVAL SECOND TO HOUR)"},
		{"interval type minute to hour reversed", "SELECT CAST(NULL AS INTERVAL MINUTE TO HOUR)"},
		{"interval type month to month invalid", "CREATE TABLE t (i INTERVAL MONTH TO MONTH)"},
		{"interval type hour to hour invalid", "CREATE TABLE t (i INTERVAL HOUR TO HOUR)"},
		{"grant with trailing garbage after principal", "GRANT SELECT ON TABLE t TO alice trailing garbage"},
		{"revoke with trailing garbage after principal", "REVOKE SELECT ON TABLE t FROM alice trailing garbage"},
		{"interval mixing year-month and day-time units", "SELECT INTERVAL 1 YEAR 1 DAY"},
		{"interval mixing month and hour string units", "SELECT INTERVAL '2' MONTH '3' HOUR"},
		{"interval mixing year and week units", "SELECT INTERVAL 1 YEAR 1 WEEK"},
		{"undeclared unknown error condition in handler", "BEGIN DECLARE EXIT HANDLER FOR NON_EXISTENT_ERROR_XYZ SELECT 0; SELECT 1; END"},
		{"execute immediate duplicate argument alias", "EXECUTE IMMEDIATE 'SELECT :x' USING 1 AS x, 2 AS x"},
		{"execute immediate parenthesized duplicate argument alias", "EXECUTE IMMEDIATE 'SELECT :x' USING (1 AS x, 2 AS x)"},
		{"execute immediate implicit duplicate argument alias", "EXECUTE IMMEDIATE 'SELECT :x' USING 1 x, 2 x"},
		{"create table tblproperties missing value", "CREATE TABLE t (x INT) TBLPROPERTIES ('k')"},
		{"alter table set tblproperties missing value", "ALTER TABLE t SET TBLPROPERTIES ('k')"},
		{"alter table unset tblproperties with value", "ALTER TABLE t UNSET TBLPROPERTIES ('k' = 'v')"},
		{"alter namespace unset dbproperties with value", "ALTER NAMESPACE ns UNSET DBPROPERTIES ('k' = 'v')"},
		{"alter view unset tblproperties with value", "ALTER VIEW v UNSET TBLPROPERTIES ('k' = 'v')"},
		{"unsupported typed literal constructor", "SELECT foobar 'abc'"},
		{"create table with schema and as select", "CREATE TABLE t (a INT) USING delta AS SELECT 1 AS a"},
		{"replace table with schema and as select", "REPLACE TABLE t (a INT) USING delta AS SELECT 1 AS a"},
		{"query order by combined with sort by", "SELECT x FROM t ORDER BY x SORT BY x"},
		{"query cluster by combined with distribute by", "SELECT x FROM t CLUSTER BY x DISTRIBUTE BY x"},
		{"query order by combined with cluster by", "SELECT x FROM t ORDER BY x CLUSTER BY x"},
		{"query order by combined with distribute by", "SELECT x FROM t ORDER BY x DISTRIBUTE BY x"},
		{"query cluster by combined with sort by", "SELECT x FROM t CLUSTER BY x SORT BY x"},
		{"alter column missing action", "ALTER TABLE t ALTER COLUMN a"},
		{"alter columns missing action", "ALTER TABLE t ALTER COLUMN a, b"},
		{"analyze table invalid statistics modifier", "ANALYZE TABLE t COMPUTE STATISTICS garbage"},
		{"analyze tables invalid statistics modifier", "ANALYZE TABLES COMPUTE STATISTICS garbage"},
		{"int type with parameter in create table", "CREATE TABLE t (x INT(10))"},
		{"int type with parameter in cast", "SELECT CAST(1 AS INT(10))"},
		{"decimal type with three parameters", "CREATE TABLE t (x DECIMAL(10, 2, 3))"},
		{"char type with two parameters", "CREATE TABLE t (x CHAR(10, 2))"},
		{"create table partitioned by and cluster by", "CREATE TABLE t (a INT) USING delta PARTITIONED BY (a) CLUSTER BY (a)"},
		{"create table cluster by and partitioned by", "CREATE TABLE t (a INT) USING delta CLUSTER BY (a) PARTITIONED BY (a)"},
		{"insert into with if not exists", "INSERT INTO t IF NOT EXISTS SELECT 1"},
		{"insert into partition with if not exists", "INSERT INTO t PARTITION (p = 1) IF NOT EXISTS SELECT 1"},
		{"duplicate partition key in insert", "INSERT INTO t PARTITION (p = 1, p = 2) SELECT 1"},
		{"duplicate partition key in alter table", "ALTER TABLE t ADD PARTITION (p = 1, P = 2)"},
		{"natural cross join", "SELECT * FROM t NATURAL CROSS JOIN s"},
		{"natural join lateral", "SELECT * FROM t NATURAL JOIN LATERAL (SELECT 1 AS x) s"},
		{"empty tablesample in select", "SELECT * FROM t TABLESAMPLE ()"},
		{"empty tablesample in pipe", "FROM t |> TABLESAMPLE ()"},
		{"unrecognized character at in select", "SELECT @x"},
		{"unrecognized character dollar in select", "SELECT $1"},
		{"pipe set dotted target", "FROM t |> SET a.b = 1"},
		{"qualified star on parenthesized expression", "SELECT (a + b).* FROM t"},
		{"qualified star on star", "SELECT *.* FROM t"},
		{"qualified star on qualified star", "SELECT t.*.* FROM t"},
		{"qualified star on function call", "SELECT f().* FROM t"},
		{"cast named arg with AS", "SELECT CAST(k => 1 AS INT)"},
		{"position named arg with IN", "SELECT POSITION(k => 1 IN s)"},
		{"substring named arg with FROM", "SELECT SUBSTRING(k => s FROM 1)"},
		{"overlay named arg with PLACING", "SELECT OVERLAY(k => s PLACING x FROM 1)"},
		{"trim named arg with FROM", "SELECT TRIM(k => ' ' FROM s)"},
		{"nested malformed dynamic identifier table", strings.Repeat("SELECT * FROM IDENTIFIER((", 25) + "SELECT 1 +" + strings.Repeat("))", 25)},
		{"nested malformed parenthesized query", strings.Repeat("SELECT ((", 25) + "SELECT 1 +" + strings.Repeat("))", 25)},
		{"nested malformed from subquery", strings.Repeat("SELECT * FROM ((", 25) + "SELECT 1 +" + strings.Repeat("))", 25)},
		{"nested malformed rollup grouping analytics", strings.Repeat("SELECT 1 FROM t GROUP BY ROLLUP(((", 25) + "SELECT 1 +" + strings.Repeat(")))", 25)},
		{"nested malformed grouping sets rollup analytics", strings.Repeat("SELECT 1 FROM t GROUP BY GROUPING SETS(ROLLUP(((", 25) + "SELECT 1 +" + strings.Repeat("))))", 25)},
		{"nested malformed table partition expression", strings.Repeat("SELECT * FROM f(TABLE(t) PARTITION BY x, (", 25) + "SELECT 1 +" + strings.Repeat("))", 25)},
		{"nested malformed table order expression", strings.Repeat("SELECT * FROM f(TABLE(t) PARTITION BY x ORDER BY y, (", 25) + "SELECT 1 +" + strings.Repeat("))", 25)},
		{"nested malformed values inline table", strings.Repeat("SELECT * FROM VALUES((", 25) + "SELECT 1 +" + strings.Repeat("))", 25)},
		{"position parenthesized continuation with OVER suffix", "SELECT POSITION('a' IN ('abc') || 'def') OVER () FROM t"},
		{"position IN with trailing IS NULL", "SELECT POSITION('a' IN (1) IS NULL)"},
		{"cast AS with OVER suffix", "SELECT CAST(x AS INT) OVER () FROM t"},
		{"trim FROM with OVER suffix", "SELECT TRIM(x FROM s) OVER () FROM t"},
		{"extract FROM with OVER suffix", "SELECT EXTRACT(YEAR FROM d) OVER () FROM t"},
		{"position IN with OVER suffix", "SELECT POSITION('a' IN s) OVER () FROM t"},
		{"substring FROM with OVER suffix", "SELECT SUBSTRING(s FROM 1) OVER () FROM t"},
		{"overlay PLACING with OVER suffix", "SELECT OVERLAY(s PLACING 'a' FROM 1) OVER () FROM t"},
		{"first IGNORE NULLS with OVER suffix", "SELECT FIRST(x IGNORE NULLS) OVER () FROM t"},
		{"substring boolean arg with FROM", "SELECT SUBSTRING(a AND b FROM 1)"},
		{"substring boolean first arg with FOR", "SELECT SUBSTRING(a AND b, 1 FOR 2)"},
		{"substring boolean second arg with FOR", "SELECT SUBSTRING(s, a AND b FOR 2)"},
		{"trim boolean arg with FROM", "SELECT TRIM(a AND b FROM s)"},
		{"trim not arg with FROM", "SELECT TRIM(NOT a FROM s)"},
		{"position boolean arg with IN", "SELECT POSITION(a AND b IN s)"},
		{"position predicate arg with IN", "SELECT POSITION(a IS NULL IN s)"},
		{"overlay boolean arg with PLACING", "SELECT OVERLAY(a AND b PLACING 'x' FROM 1)"},
		{"overlay not arg with PLACING", "SELECT OVERLAY(NOT a PLACING 'x' FROM 1)"},
		{"struct mixing field alias then named arg", "SELECT STRUCT(x AS a, k => 1)"},
		{"struct mixing implicit alias then named arg", "SELECT STRUCT(x a, k => 1)"},
		{"struct mixing named arg then field alias", "SELECT STRUCT(k => 1, x AS a)"},
		{"struct mixing named arg then implicit alias", "SELECT STRUCT(k => 1, x a)"},
		{"tablesample bucket qualified name without parens", "SELECT * FROM t TABLESAMPLE (BUCKET 1 OUT OF 2 ON a.b)"},
		{"alter table partition set tblproperties", "ALTER TABLE t PARTITION (p = 1) SET TBLPROPERTIES ('x' = 'y')"},
		{"arabic-indic digit in unquoted identifier", "SELECT a\u0663 FROM t"},
		{"nel whitespace rejected by spark", "SELECT\u00851"},
		{"paragraph separator rejected by spark", "SELECT\u20291"},
		{"interval from-to followed by multi-unit", "SELECT INTERVAL '1-2' YEAR TO MONTH 3 DAY"},
		{"interval from-to followed by from-to", "SELECT INTERVAL '1-2' YEAR TO MONTH '2-3' YEAR TO MONTH"},
		{"interval multi-unit followed by from-to", "SELECT INTERVAL 3 DAY '1-2' YEAR TO MONTH"},
		{"decimal with invalid L suffix", "SELECT 1.0L"},
		{"decimal with invalid Y suffix", "SELECT 1.0Y"},
		{"decimal with invalid S suffix", "SELECT 1.0S"},
		{"grouping sets followed by with rollup", "SELECT a FROM t GROUP BY GROUPING SETS ((a)) WITH ROLLUP"},
		{"grouping sets followed by with cube", "SELECT a FROM t GROUP BY GROUPING SETS ((a)) WITH CUBE"},
		{"grouping sets followed by grouping sets", "SELECT a FROM t GROUP BY GROUPING SETS ((a)) GROUPING SETS ((b))"},
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
			name: "deeply nested position with parenthesized in operand",
			sql:  "SELECT " + strings.Repeat("POSITION(x IN (", 8000) + "1" + strings.Repeat("))", 8000),
		},
		{
			name: "deeply nested parentheses",
			sql:  "SELECT " + strings.Repeat("(", 8000) + "1" + strings.Repeat(")", 8000),
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
			name: "deeply nested case function calls linear time",
			sql:  "SELECT " + strings.Repeat("CASE(", 4000) + "1" + strings.Repeat(")", 4000),
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
			name: "repeated case identifier arithmetic with parenthesized calls",
			sql:  "SELECT " + strings.Repeat("case + abs(1) + ", 5000) + "1",
		},
		{
			name: "repeated case identifier arithmetic with nested case in calls",
			sql:  "SELECT " + strings.Repeat("case + abs(case + 1) + ", 2000) + "1",
		},
		{
			name: "repeated case identifier arithmetic with calls inside nested simple cases",
			sql:  "SELECT CASE CASE " + strings.Repeat("case + abs(case + 1) + ", 1000) + "1 WHEN 1 THEN 2 END WHEN 2 THEN 3 END",
		},
		{
			name: "repeated sibling simple case expressions",
			sql:  "SELECT " + strings.Repeat("CASE x WHEN 1 THEN 2 END + ", 8000) + "1",
		},
		{
			name: "repeated sibling simple case expressions inside outer case",
			sql:  "SELECT CASE y WHEN 1 THEN " + strings.Repeat("CASE x WHEN 1 THEN 2 END + ", 4000) + "1 END",
		},
		{
			name: "repeated sibling simple case expressions in where clause",
			sql:  "SELECT 1 FROM t WHERE " + strings.Repeat("CASE a WHEN 1 THEN true END AND ", 4000) + "true",
		},
		{
			name: "deeply nested searched case expressions",
			sql:  "SELECT " + strings.Repeat("CASE WHEN true THEN ", 8000) + "1" + strings.Repeat(" END", 8000),
		},
		{
			name: "deeply nested simple case expressions",
			sql:  "SELECT " + strings.Repeat("CASE x WHEN 1 THEN ", 4000) + "1" + strings.Repeat(" END", 4000),
		},
		{
			name: "deeply nested searched case containing sibling cases",
			sql:  "SELECT " + strings.Repeat("CASE WHEN true THEN ", 2000) + strings.Repeat("CASE x WHEN 1 THEN 2 END + ", 2000) + "1" + strings.Repeat(" END", 2000),
		},
		{
			name: "sibling cases each containing nested cases inside outer case",
			sql:  "SELECT CASE WHEN true THEN " + strings.Repeat("CASE WHEN true THEN CASE WHEN true THEN 1 END END + ", 2000) + "1 END",
		},
		{
			name: "deeply nested parenthesized simple case operands",
			sql:  "SELECT " + strings.Repeat("CASE (", 2000) + "1" + strings.Repeat(") WHEN 1 THEN 2 END", 2000),
		},
		{
			name: "searched case with case column before next when",
			sql:  "SELECT CASE WHEN a THEN case WHEN b THEN 1 END",
		},
		{
			name: "sibling simple case with case column operand",
			sql:  "SELECT CASE WHEN true THEN CASE WHEN 1 THEN 2 END + CASE case WHEN 1 THEN 2 END END",
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
			name: "nested scalar subquery in row constructor",
			sql:  strings.Repeat("SELECT ((", 30) + "SELECT 1" + strings.Repeat("), 1)", 30),
		},
		{
			name: "nested scalar subquery in parenthesized arithmetic",
			sql:  strings.Repeat("SELECT ((", 30) + "SELECT 1" + strings.Repeat(") + 1)", 30),
		},
		{
			name: "nested rollup grouping-set parenthesized expression",
			sql:  strings.Repeat("SELECT 1 FROM t GROUP BY ROLLUP(((", 25) + "SELECT 1" + strings.Repeat(")) + 1)", 25),
		},
		{
			name: "nested rollup function call with arithmetic in group by",
			sql:  strings.Repeat("SELECT 1 FROM t GROUP BY ROLLUP((", 25) + "SELECT 1" + strings.Repeat(")) + 1", 25),
		},
		{
			name: "nested rollup named row constructor",
			sql:  strings.Repeat("SELECT 1 FROM t GROUP BY ROLLUP(((", 25) + "SELECT 1" + strings.Repeat(") AS a, 1))", 25),
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

func TestBeginBlockErrorPositions(t *testing.T) {
	cases := []struct {
		name   string
		sql    string
		minPos int
	}{
		{"bad select inside begin", "BEGIN SELECT (; END", 13},
		{"trailing token after end", "BEGIN SELECT 1; END zzz", 20},
		{"incomplete declare inside begin", "BEGIN DECLARE x INT DEFAULT; END", 27},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) unexpectedly succeeded", tc.sql)
			}
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("expected *ParseError, got %T: %v", err, err)
			}
			if pe.Position < tc.minPos {
				t.Fatalf("Parse(%q) error position = %d, want >= %d (%v)", tc.sql, pe.Position, tc.minPos, pe)
			}
		})
	}
}

func TestStatementSpansAndMultiStatement(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		want []struct {
			kind StmtKind
			loc  Loc
		}
	}{
		{
			name: "select with trailing semicolon excludes semicolon",
			sql:  "SELECT 1;",
			want: []struct {
				kind StmtKind
				loc  Loc
			}{{StmtQuery, Loc{Start: 0, End: 8}}},
		},
		{
			name: "begin block with trailing semicolon excludes semicolon",
			sql:  "BEGIN SELECT 1; END;",
			want: []struct {
				kind StmtKind
				loc  Loc
			}{{StmtCompound, Loc{Start: 0, End: 19}}},
		},
		{
			name: "two selects separated by semicolon",
			sql:  "SELECT 1; SELECT 2",
			want: []struct {
				kind StmtKind
				loc  Loc
			}{
				{StmtQuery, Loc{Start: 0, End: 8}},
				{StmtQuery, Loc{Start: 10, End: 18}},
			},
		},
		{
			name: "refresh followed by select",
			sql:  "REFRESH a; SELECT 1",
			want: []struct {
				kind StmtKind
				loc  Loc
			}{
				{StmtRefresh, Loc{Start: 0, End: 9}},
				{StmtQuery, Loc{Start: 11, End: 19}},
			},
		},
		{
			name: "begin block followed by select",
			sql:  "BEGIN SELECT 1; END; SELECT 2",
			want: []struct {
				kind StmtKind
				loc  Loc
			}{
				{StmtCompound, Loc{Start: 0, End: 19}},
				{StmtQuery, Loc{Start: 21, End: 29}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			if len(res.Statements) != len(tc.want) {
				t.Fatalf("Parse(%q) got %d statements, want %d", tc.sql, len(res.Statements), len(tc.want))
			}
			for i, w := range tc.want {
				ps, ok := res.Statements[i].(*ParsedStmt)
				if !ok {
					t.Fatalf("stmt[%d] = %T, want *ParsedStmt", i, res.Statements[i])
				}
				if ps.Kind != w.kind || ps.Loc != w.loc {
					t.Errorf("stmt[%d] = {%s, %+v}, want {%s, %+v}", i, ps.Kind, ps.Loc, w.kind, w.loc)
				}
			}
		})
	}
}

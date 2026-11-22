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
	deep := "SELECT " + strings.Repeat("(", maxNestingDepth+50) + "1" + strings.Repeat(")", maxNestingDepth+50)
	if _, err := Parse(deep); err == nil {
		t.Fatalf("expected nesting depth error for deeply nested expression")
	}
}

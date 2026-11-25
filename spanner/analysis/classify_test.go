package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/spanner/ast"
)

func TestClassifySQL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want QueryType
	}{
		// --- SELECT family ---
		{"plain select", "SELECT 1", Select},
		{"select from", "SELECT a, b FROM t", Select},
		{"select star", "SELECT * FROM t", Select},
		{"with cte", "WITH c AS (SELECT 1 AS n) SELECT n FROM c", Select},
		{"set op union", "SELECT a FROM t UNION ALL SELECT a FROM u", Select},
		{"parenthesized query", "(SELECT 1) UNION ALL (SELECT 2)", Select},
		// The `TABLE t` query primary (query_primary: table_clause_reserved) is a
		// read of that table, exactly like `SELECT * FROM t`.
		{"table clause", "TABLE t", Select},
		{"table clause with filter", "TABLE t WHERE a > 1", Select},
		{"table clause tvf", "TABLE ds.f(a)", Select},

		// --- DML ---
		{"insert", "INSERT INTO t (a) VALUES (1)", DML},
		{"insert select", "INSERT INTO t SELECT a FROM u", DML},
		{"insert or update (spanner upsert)", "INSERT OR UPDATE INTO t (a) VALUES (1)", DML},
		{"update", "UPDATE t SET a = 1 WHERE b = 2", DML},
		{"delete", "DELETE FROM t WHERE a = 1", DML},
		{"merge", "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", DML},

		// --- DDL ---
		{"create table", "CREATE TABLE t (a INT64) PRIMARY KEY (a)", DDL},
		{"create view", "CREATE VIEW v AS SELECT 1", DDL},
		{"create index", "CREATE INDEX idx ON t (a)", DDL},
		{"create schema", "CREATE SCHEMA s", DDL},
		{"create database", "CREATE DATABASE d", DDL},
		{"alter table", "ALTER TABLE t ADD COLUMN b INT64", DDL},
		{"drop table", "DROP TABLE t", DDL},
		{"drop view", "DROP VIEW v", DDL},

		// --- DDL — BigQuery-only objects (parser-ddl-bigquery node) ---
		{"create function", "CREATE FUNCTION ds.f(x INT64) RETURNS INT64 AS (x + 1)", DDL},
		{"create vector index", "CREATE VECTOR INDEX i ON ds.t(c) OPTIONS(distance_type='COSINE')", DDL},
		{"drop function", "DROP FUNCTION ds.f", DDL},

		// --- DDL — Spanner-only objects (parser-ddl-spanner node). These now PARSE
		// (previously over-rejected), so classification must report DDL, not Unknown. ---
		{"create change stream", "CREATE CHANGE STREAM s FOR ALL", DDL},
		{"alter change stream", "ALTER CHANGE STREAM s SET FOR ALL", DDL},
		{"drop change stream", "DROP CHANGE STREAM s", DDL},
		{"create sequence", "CREATE SEQUENCE q", DDL},
		{"alter sequence", "ALTER SEQUENCE q SET OPTIONS (start_with_counter=1)", DDL},
		{"drop sequence", "DROP SEQUENCE q", DDL},
		{"create role", "CREATE ROLE analyst", DDL},
		{"drop role", "DROP ROLE analyst", DDL},
		{"create locality group", "CREATE LOCALITY GROUP g OPTIONS (storage='ssd')", DDL},
		{"alter locality group", "ALTER LOCALITY GROUP g SET OPTIONS (storage='hdd')", DDL},
		{"drop locality group", "DROP LOCALITY GROUP g", DDL},
		{"create proto bundle", "CREATE PROTO BUNDLE (`a.b.C`)", DDL},
		{"alter proto bundle", "ALTER PROTO BUNDLE INSERT (`a.b.C`)", DDL},
		{"grant to role", "GRANT SELECT ON TABLE t TO ROLE r", DDL},
		{"undrop schema", "UNDROP SCHEMA s", DDL},
		{"grant role to role", "GRANT ROLE a TO ROLE b", DDL},
		{"revoke role from role", "REVOKE ROLE a FROM ROLE b", DDL},

		// --- DDL — CREATE PROPERTY GRAPH (parser-gql node). Every CREATE form is
		// DDL; the legacy queryTypeListener lists create_property_graph_statement
		// directly among the DDL alternatives. ---
		{"create property graph", "CREATE PROPERTY GRAPH g NODE TABLES (Person)", DDL},
		{"create or replace property graph", "CREATE OR REPLACE PROPERTY GRAPH g NODE TABLES (Person) EDGE TABLES (Knows SOURCE KEY (a) REFERENCES Person DESTINATION KEY (b) REFERENCES Person)", DDL},

		// --- Unknown / empty ---
		{"empty", "", Unknown},
		{"whitespace only", "   \n  ", Unknown},
		{"comment only", "-- just a comment\n", Unknown},
		{"garbage", "this is not sql", Unknown},

		// --- GQL graph query (parser-gql node). A bare `GRAPH … MATCH … RETURN …`
		// read classifies UNKNOWN, NOT Select: the legacy queryTypeListener has no
		// gql_statement case and falls through to `default: QueryTypeUnknown` in both
		// dialects. This is deliberate parity (see Classify's GQLStmt case + the
		// divergence ledger), so it lives with the Unknown group. ---
		{"gql graph query", "GRAPH g MATCH (n) RETURN n", Unknown},
		{"gql graph query with ops", "GRAPH g MATCH (a)-[e]->(b) WHERE a.id = 1 RETURN a, b", Unknown},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClassifySQL(tc.sql)
			if got != tc.want {
				t.Errorf("ClassifySQL(%q) = %v, want %v", tc.sql, got, tc.want)
			}
		})
	}
}

// TestClassify_InfoSchema covers the SelectInfoSchema promotion, which is
// dialect-specific: BigQuery promotes only on INFORMATION_SCHEMA; Spanner
// promotes on INFORMATION_SCHEMA or SPANNER_SYS. The promotion fires only for
// SELECT-family statements that read EXCLUSIVELY system tables.
func TestClassify_InfoSchema(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want QueryType
	}{
		{
			name: "information_schema both",
			sql:  "SELECT * FROM INFORMATION_SCHEMA.TABLES",
			want: SelectInfoSchema,
		},
		{
			name: "dataset qualified information_schema",
			sql:  "SELECT * FROM mydataset.INFORMATION_SCHEMA.COLUMNS",
			want: SelectInfoSchema,
		},
		{
			name: "spanner_sys only spanner",
			sql:  "SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE", // BigQuery has no SPANNER_SYS notion
			want: SelectInfoSchema,
		},
		{
			name: "user table stays select",
			sql:  "SELECT * FROM mydataset.users",
			want: Select,
		},
		{
			name: "non-select referencing info schema is not promoted",
			sql:  "CREATE VIEW v AS SELECT * FROM INFORMATION_SCHEMA.TABLES",
			want: DDL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ClassifySQL(tc.sql); got != tc.want {
				t.Errorf("ClassifySQL(%q) = %v, want %v", tc.sql, got, tc.want)
			}
		})
	}
}

// TestClassifyRecordsTheSpanTables: Classify's pass skips lineage and
// projections (spanWalker.tablesOnly), so it must still record exactly the
// tables, in order and with their locations, that the full span pass records.
// The statements reach tables from every place the skipped code sits next to:
// select items, star modifiers, joins of every kind, UNNEST and correlated
// paths, set operations, recursive CTEs, and DML.
func TestClassifyRecordsTheSpanTables(t *testing.T) {
	for _, sql := range []string{
		"SELECT (SELECT a FROM s), * REPLACE ((SELECT b FROM r) AS b) FROM t",
		"SELECT * EXCEPT (a) FROM INFORMATION_SCHEMA.TABLES WHERE a IN (SELECT a FROM u)",
		"SELECT * FROM t JOIN u USING (a) JOIN (SELECT a FROM v UNION ALL SELECT * FROM w) s USING (a)",
		"SELECT * FROM t NATURAL JOIN u, tvf((SELECT 1 FROM v)) JOIN w ON w.a = (SELECT 1 FROM x)",
		"SELECT * FROM t AS a0, a0.arr AS a1, a1.child, UNNEST(a1.list) AS e, e.f, ds.t2",
		"SELECT * FROM (t, u) JOIN (v JOIN (w JOIN x USING (k)) USING (k)) USING (k)",
		"SELECT f(a).b, (SELECT 1 FROM s).c, a[OFFSET(0)].d FROM t",
		"WITH RECURSIVE r AS (SELECT a FROM t UNION ALL SELECT a FROM r JOIN u USING (a)) SELECT * FROM r, SPANNER_SYS.X",
		"WITH c AS (SELECT * FROM t) SELECT * FROM c JOIN c AS d USING (a) UNION ALL BY NAME SELECT * FROM u",
		"SELECT a FROM t UNION ALL (SELECT a FROM u EXCEPT DISTINCT SELECT a FROM v) UNION ALL SELECT * FROM w",
		"INSERT INTO t SELECT * FROM u JOIN v USING (a)",
		"UPDATE t SET a = (SELECT 1 FROM s) FROM u, u.arr AS x WHERE x.b = 1",
		"SELECT x FROM t |> JOIN u ON true |> WHERE x IN (SELECT x FROM v)",
		"|> UNION ALL TABLE u, (SELECT 1 FROM v)",
		"|> WITH c AS (SELECT 1 FROM v) |> JOIN u",
		"SELECT 1 |> INSERT t (x)",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> CALL f(TABLE secret)",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> CALL f(arg => TABLE secret)",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> JOIN INFORMATION_SCHEMA.COLUMNS TABLESAMPLE SYSTEM (10 PERCENT PARTITION BY (SELECT x FROM secret)) ON true",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> JOIN f(TABLE secret) ON true",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> JOIN UNNEST([1]) TABLESAMPLE SYSTEM (10 PERCENT PARTITION BY (SELECT x FROM secret)) ON true",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> CALL f(TABLE g((SELECT x FROM secret)))",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> SELECT t.* REPLACE ((SELECT x FROM secret) AS c)",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> WINDOW t.* REPLACE ((SELECT x FROM secret) AS c)",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> AGGREGATE COUNT(*) GROUP BY ROLLUP((SELECT x FROM secret))",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> UNION ALL TABLE f(TABLE secret)",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> PIVOT(SUM((SELECT x FROM secret)) FOR c IN (1))",
		"SELECT * FROM INFORMATION_SCHEMA.TABLES |> SELECT SUM(x) OVER w WINDOW w AS (PARTITION BY (SELECT y FROM secret))",
	} {
		file, err := parseFile(sql)
		if err != nil {
			t.Fatalf("parse %q: %v", sql, err)
		}
		full := &QuerySpan{}
		newSpanWalker(full).analyzeStmt(file.Stmts[0])
		got := collectAccessTables(file.Stmts[0])
		if !reflect.DeepEqual(got, full.AccessTables) {
			t.Errorf("%q: Classify records %+v, the span records %+v", sql, got, full.AccessTables)
		}
	}
}

func TestClassify_Node(t *testing.T) {
	// nil node is Unknown.
	if got := Classify(nil); got != Unknown {
		t.Errorf("Classify(nil) = %v, want Unknown", got)
	}
	// An unrecognized node type is Unknown.
	if got := Classify(&ast.File{}); got != Unknown {
		t.Errorf("Classify(*ast.File) = %v, want Unknown", got)
	}
}

func TestQueryType_String(t *testing.T) {
	tests := []struct {
		qt   QueryType
		want string
	}{
		{Unknown, "UNKNOWN"},
		{Select, "SELECT"},
		{Explain, "EXPLAIN"},
		{SelectInfoSchema, "SELECT_INFO_SCHEMA"},
		{DDL, "DDL"},
		{DML, "DML"},
	}
	for _, tc := range tests {
		if got := tc.qt.String(); got != tc.want {
			t.Errorf("QueryType(%d).String() = %q, want %q", tc.qt, got, tc.want)
		}
	}
}

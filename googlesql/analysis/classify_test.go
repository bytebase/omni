package analysis

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
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

		// --- DML ---
		{"insert", "INSERT INTO t (a) VALUES (1)", DML},
		{"insert select", "INSERT INTO t SELECT a FROM u", DML},
		{"insert or update (spanner upsert)", "INSERT OR UPDATE INTO t (a) VALUES (1)", DML},
		{"update", "UPDATE t SET a = 1 WHERE b = 2", DML},
		{"delete", "DELETE FROM t WHERE a = 1", DML},
		{"merge", "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", DML},
		{"truncate", "TRUNCATE TABLE t", DML},

		// --- DDL ---
		{"create table", "CREATE TABLE t (a INT64) PRIMARY KEY (a)", DDL},
		{"create view", "CREATE VIEW v AS SELECT 1", DDL},
		{"create index", "CREATE INDEX idx ON t (a)", DDL},
		{"create schema", "CREATE SCHEMA s", DDL},
		{"create database", "CREATE DATABASE d", DDL},
		{"alter table", "ALTER TABLE t ADD COLUMN b INT64", DDL},
		{"drop table", "DROP TABLE t", DDL},
		{"drop view", "DROP VIEW v", DDL},
		{"grant", "GRANT SELECT ON t TO 'user@example.com'", DDL},
		{"revoke", "REVOKE SELECT ON t FROM 'user@example.com'", DDL},

		// --- Unknown / empty ---
		{"empty", "", Unknown},
		{"whitespace only", "   \n  ", Unknown},
		{"comment only", "-- just a comment\n", Unknown},
		{"garbage", "this is not sql", Unknown},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClassifySQL(tc.sql, DialectBigQuery)
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
		name        string
		sql         string
		wantBQ      QueryType
		wantSpanner QueryType
	}{
		{
			name:        "information_schema both",
			sql:         "SELECT * FROM INFORMATION_SCHEMA.TABLES",
			wantBQ:      SelectInfoSchema,
			wantSpanner: SelectInfoSchema,
		},
		{
			name:        "dataset qualified information_schema",
			sql:         "SELECT * FROM mydataset.INFORMATION_SCHEMA.COLUMNS",
			wantBQ:      SelectInfoSchema,
			wantSpanner: SelectInfoSchema,
		},
		{
			name:        "spanner_sys only spanner",
			sql:         "SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE",
			wantBQ:      Select, // BigQuery has no SPANNER_SYS notion
			wantSpanner: SelectInfoSchema,
		},
		{
			name:        "user table stays select",
			sql:         "SELECT * FROM mydataset.users",
			wantBQ:      Select,
			wantSpanner: Select,
		},
		{
			name:        "non-select referencing info schema is not promoted",
			sql:         "CREATE VIEW v AS SELECT * FROM INFORMATION_SCHEMA.TABLES",
			wantBQ:      DDL,
			wantSpanner: DDL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ClassifySQL(tc.sql, DialectBigQuery); got != tc.wantBQ {
				t.Errorf("BigQuery ClassifySQL(%q) = %v, want %v", tc.sql, got, tc.wantBQ)
			}
			if got := ClassifySQL(tc.sql, DialectSpanner); got != tc.wantSpanner {
				t.Errorf("Spanner ClassifySQL(%q) = %v, want %v", tc.sql, got, tc.wantSpanner)
			}
		})
	}
}

func TestClassify_Node(t *testing.T) {
	// nil node is Unknown.
	if got := Classify(nil, DialectBigQuery); got != Unknown {
		t.Errorf("Classify(nil) = %v, want Unknown", got)
	}
	// An unrecognized node type is Unknown.
	if got := Classify(&ast.File{}, DialectBigQuery); got != Unknown {
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

package analysis

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want QueryType
	}{
		// ---- SELECT family ----
		{"plain select", "SELECT 1", Select},
		{"select from table", "SELECT a FROM t", Select},
		{"with cte select", "WITH c AS (SELECT 1) SELECT * FROM c", Select},
		{"table query", "TABLE orders", Select},
		{"values query", "VALUES (1), (2)", Select},
		{"parenthesized select", "(SELECT 1)", Select},
		{"lowercase select", "select 1", Select},
		{"leading whitespace and comment", "  -- hi\n SELECT 1", Select},

		// ---- SelectInfoSchema: SELECT touching a system/info schema ----
		{"select information_schema", "SELECT * FROM information_schema.tables", SelectInfoSchema},
		{"select system schema", "SELECT * FROM system.runtime.nodes", SelectInfoSchema},
		{"select metadata schema", "SELECT * FROM metadata.foo", SelectInfoSchema},
		{"with cte over information_schema", "WITH c AS (SELECT * FROM information_schema.columns) SELECT * FROM c", SelectInfoSchema},

		// ---- SHOW / DESCRIBE family → SelectInfoSchema ----
		{"show tables", "SHOW TABLES", SelectInfoSchema},
		{"show schemas", "SHOW SCHEMAS", SelectInfoSchema},
		{"show columns", "SHOW COLUMNS FROM t", SelectInfoSchema},
		{"show catalogs", "SHOW CATALOGS", SelectInfoSchema},
		{"describe", "DESCRIBE t", SelectInfoSchema},
		{"desc", "DESC t", SelectInfoSchema},

		// ---- EXPLAIN family ----
		{"explain select", "EXPLAIN SELECT 1", Explain},
		{"explain analyze select", "EXPLAIN ANALYZE SELECT a FROM t", Explain},

		// ---- DML ----
		{"insert", "INSERT INTO t VALUES (1)", DML},
		{"insert select", "INSERT INTO t SELECT * FROM s", DML},
		{"update", "UPDATE t SET a = 1", DML},
		{"delete", "DELETE FROM t WHERE a = 1", DML},
		{"merge", "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", DML},
		{"truncate", "TRUNCATE TABLE t", DML},
		{"call", "CALL system.runtime.kill_query('x')", DML},

		// ---- DDL ----
		{"create table", "CREATE TABLE t (id INT)", DDL},
		{"create table as", "CREATE TABLE t AS SELECT 1", DDL},
		{"create view", "CREATE VIEW v AS SELECT 1", DDL},
		{"create schema", "CREATE SCHEMA s", DDL},
		{"alter table", "ALTER TABLE t ADD COLUMN c INT", DDL},
		{"drop table", "DROP TABLE t", DDL},
		{"drop view", "DROP VIEW v", DDL},
		{"drop schema", "DROP SCHEMA s", DDL},
		{"comment on", "COMMENT ON TABLE t IS 'x'", DDL},
		{"analyze", "ANALYZE t", DDL},
		{"grant", "GRANT SELECT ON t TO u", DDL},
		{"revoke", "REVOKE SELECT ON t FROM u", DDL},
		{"deny", "DENY SELECT ON t TO u", DDL},
		{"create role", "CREATE ROLE r", DDL},
		{"drop role", "DROP ROLE r", DDL},
		{"refresh materialized view", "REFRESH MATERIALIZED VIEW mv", DDL},

		// ---- Session / admin: read-only-ish, classified Select per legacy ----
		{"set session", "SET SESSION x = 1", Select},
		{"reset session", "RESET SESSION x", Select},

		// ---- Unknown / empty ----
		{"empty", "", Unknown},
		{"whitespace only", "   \n  ", Unknown},
		{"comment only", "-- just a comment", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Classify(tt.sql)
			if got != tt.want {
				t.Errorf("Classify(%q) = %v, want %v", tt.sql, got, tt.want)
			}
		})
	}
}

func TestQueryType_String(t *testing.T) {
	tests := []struct {
		qt   QueryType
		want string
	}{
		{Select, "SELECT"},
		{SelectInfoSchema, "SELECT_INFO_SCHEMA"},
		{Explain, "EXPLAIN"},
		{DML, "DML"},
		{DDL, "DDL"},
		{Unknown, "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := tt.qt.String(); got != tt.want {
			t.Errorf("QueryType(%d).String() = %q, want %q", tt.qt, got, tt.want)
		}
	}
}

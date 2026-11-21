package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/trino/parser"
)

func queryTypeOfText(t *testing.T, statement string) (QueryType, bool) {
	t.Helper()
	file, err := parser.Parse(statement)
	require.NoErrorf(t, err, "%q", statement)
	require.Len(t, file.Stmts, 1)
	return QueryTypeOf(file.Stmts[0])
}

// TestQueryTypeClassification verifies the statement classification (the
// query_type contract) and the read-only / data-changing / schema-changing
// helpers against the AST. It replaces the adapter test that asserted an
// ANTLR-specific StatementType enum; the user-facing behavior (QueryType plus
// the Is*Statement predicates) is preserved here.
func TestQueryTypeClassification(t *testing.T) {
	tests := []struct {
		name          string
		sql           string
		wantQueryType QueryType
		wantIsAnalyze bool
	}{
		{name: "SELECT statement", sql: "SELECT * FROM users", wantQueryType: Select},
		{name: "WITH FUNCTION query", sql: "WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(10)", wantQueryType: Select},
		{name: "EXPLAIN statement", sql: "EXPLAIN SELECT * FROM users", wantQueryType: Explain},
		{name: "EXPLAIN ANALYZE statement", sql: "EXPLAIN ANALYZE SELECT * FROM users", wantQueryType: Select, wantIsAnalyze: true},
		{name: "INSERT statement", sql: "INSERT INTO users (id, name) VALUES (1, 'John')", wantQueryType: DML},
		{name: "UPDATE statement", sql: "UPDATE users SET name = 'John' WHERE id = 1", wantQueryType: DML},
		{name: "DELETE statement", sql: "DELETE FROM users WHERE id = 1", wantQueryType: DML},
		{name: "CREATE TABLE statement", sql: "CREATE TABLE users (id INT, name VARCHAR)", wantQueryType: DDL},
		{name: "CREATE TABLE AS SELECT statement", sql: "CREATE TABLE new_users AS SELECT * FROM users", wantQueryType: DDL},
		{name: "CREATE VIEW statement", sql: "CREATE VIEW user_view AS SELECT id, name FROM users", wantQueryType: DDL},
		{name: "ALTER TABLE ADD COLUMN", sql: "ALTER TABLE users ADD COLUMN email VARCHAR", wantQueryType: DDL},
		{name: "ALTER TABLE DROP COLUMN", sql: "ALTER TABLE users DROP COLUMN email", wantQueryType: DDL},
		{name: "ALTER TABLE RENAME COLUMN", sql: "ALTER TABLE users RENAME COLUMN old_name TO new_name", wantQueryType: DDL},
		{name: "ALTER TABLE SET COLUMN TYPE", sql: "ALTER TABLE users ALTER COLUMN name SET DATA TYPE VARCHAR(100)", wantQueryType: DDL},
		{name: "DROP TABLE statement", sql: "DROP TABLE users", wantQueryType: DDL},
		{name: "DROP VIEW statement", sql: "DROP VIEW user_view", wantQueryType: DDL},
		{name: "CREATE SCHEMA statement", sql: "CREATE SCHEMA new_schema", wantQueryType: DDL},
		{name: "DROP SCHEMA statement", sql: "DROP SCHEMA old_schema", wantQueryType: DDL},
		{name: "RENAME TABLE statement", sql: "ALTER TABLE users RENAME TO new_users", wantQueryType: DDL},
		{name: "RENAME SCHEMA statement", sql: "ALTER SCHEMA old_schema RENAME TO new_schema", wantQueryType: DDL},
		{name: "SHOW TABLES statement", sql: "SHOW TABLES", wantQueryType: SelectInfoSchema},
		{name: "SHOW SCHEMAS statement", sql: "SHOW SCHEMAS", wantQueryType: SelectInfoSchema},
		{name: "SHOW COLUMNS statement", sql: "SHOW COLUMNS FROM users", wantQueryType: SelectInfoSchema},
		{name: "SET SESSION statement", sql: "SET SESSION optimize_hash_generation = true", wantQueryType: Select},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.Parse(tt.sql)
			require.NoError(t, err)
			require.Len(t, file.Stmts, 1)
			node := file.Stmts[0]

			queryType, isAnalyze := QueryTypeOf(node)
			assert.Equal(t, tt.wantQueryType, queryType, "Query type mismatch")
			assert.Equal(t, tt.wantIsAnalyze, isAnalyze, "isAnalyze mismatch")

			assert.Equal(t, tt.wantQueryType == Select || tt.wantQueryType == Explain || tt.wantQueryType == SelectInfoSchema, IsReadOnlyStatement(node), "IsReadOnlyStatement mismatch")
			assert.Equal(t, tt.wantQueryType == DML, IsDataChangingStatement(node), "IsDataChangingStatement mismatch")
			assert.Equal(t, tt.wantQueryType == DDL, IsSchemaChangingStatement(node), "IsSchemaChangingStatement mismatch")
		})
	}
}

// TestQueryTypeOfEXPLAINANALYZEWrite pins the flag that makes the editor refuse
// an EXPLAIN ANALYZE of a write: the inner type is DML and analyze is true, not
// a read-only Explain.
func TestQueryTypeOfEXPLAINANALYZEWrite(t *testing.T) {
	queryType, isAnalyze := queryTypeOfText(t, "EXPLAIN ANALYZE UPDATE users SET name = 'John'")
	require.Equal(t, DML, queryType)
	require.True(t, isAnalyze)
}

// TestQueryTypeFromText pins the statement-text system-schema promotion that
// only a text entry point can apply.
func TestQueryTypeFromText(t *testing.T) {
	tests := []struct {
		statement string
		want      QueryType
	}{
		{"SELECT * FROM users", Select},
		{"SELECT * FROM system.runtime.nodes", SelectInfoSchema},
		{"SELECT * FROM information_schema.tables", SelectInfoSchema},
		{"SELECT * FROM catalog.schemas", SelectInfoSchema},
		{"SELECT * FROM metadata.schemas", SelectInfoSchema},
		{"EXPLAIN SELECT * FROM system.runtime.nodes", SelectInfoSchema},
		{"INSERT INTO users VALUES (1)", DML},
		{"", Unknown},
		{"SELECT * FROM", Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.statement, func(t *testing.T) {
			require.Equal(t, tt.want, queryTypeFromText(tt.statement))
		})
	}
}

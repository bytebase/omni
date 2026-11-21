package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantReadOnly bool
		wantData     bool
		wantErr      bool
	}{
		{name: "Simple SELECT", sql: "SELECT * FROM users", wantReadOnly: true, wantData: true},
		{name: "INSERT statement", sql: "INSERT INTO users (id, name) VALUES (1, 'John')"},
		{name: "UPDATE statement", sql: "UPDATE users SET name = 'John' WHERE id = 1"},
		{name: "DELETE statement", sql: "DELETE FROM users WHERE id = 1"},
		{name: "CREATE TABLE statement", sql: "CREATE TABLE users (id INT, name VARCHAR)"},
		{name: "EXPLAIN statement", sql: "EXPLAIN SELECT * FROM users", wantReadOnly: true, wantData: true},
		{name: "EXPLAIN ANALYZE SELECT", sql: "EXPLAIN ANALYZE SELECT * FROM users", wantReadOnly: true, wantData: true},
		// EXPLAIN ANALYZE EXECUTES the inner UPDATE, so it is NOT read-only.
		{name: "EXPLAIN ANALYZE UPDATE", sql: "EXPLAIN ANALYZE UPDATE users SET name = 'John' WHERE id = 1"},
		{name: "WITH FUNCTION SELECT", sql: "WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", wantReadOnly: true, wantData: true},
		{name: "EXPLAIN WITH FUNCTION SELECT", sql: "EXPLAIN WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", wantReadOnly: true, wantData: true},
		{name: "EXPLAIN FORMAT JSON WITH FUNCTION SELECT", sql: "EXPLAIN (FORMAT JSON) WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", wantReadOnly: true, wantData: true},
		{name: "EXPLAIN ANALYZE WITH FUNCTION SELECT", sql: "EXPLAIN ANALYZE WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", wantReadOnly: true, wantData: true},
		{name: "SHOW TABLES", sql: "SHOW TABLES", wantReadOnly: true, wantData: true},
		// Session statements rebind the connection; they return no data.
		{name: "USE", sql: "USE hive.default", wantReadOnly: true},
		{name: "SET SESSION", sql: "SET SESSION optimize_hash_generation = true", wantReadOnly: true},
		{name: "RESET SESSION", sql: "RESET SESSION optimize_hash_generation", wantReadOnly: true},
		{name: "multi-statement read", sql: "SELECT * FROM users; SHOW TABLES", wantReadOnly: true, wantData: true},
		{name: "multi-statement with a write", sql: "SELECT * FROM users; DELETE FROM users"},
		{name: "Invalid SQL", sql: "SELECT * FFROM users", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tt.sql)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantReadOnly, readOnly, "readOnly")
			assert.Equal(t, tt.wantData, returnsData, "returnsData")
		})
	}
}

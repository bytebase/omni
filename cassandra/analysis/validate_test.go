package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cassandra/parser"
)

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name        string
		cql         string
		readOnly    bool
		returnsData bool
		wantErr     bool
	}{
		// --- Valid read-only CQL SELECT statements ---
		{name: "simple select", cql: "SELECT * FROM users", readOnly: true, returnsData: true},
		{name: "select star no space", cql: "select* from users", readOnly: true, returnsData: true},
		{name: "select qualified keyspace", cql: "SELECT id, name FROM ks.users WHERE id = 1;", readOnly: true, returnsData: true},
		{name: "select allow filtering", cql: "SELECT * FROM ks.users WHERE age > 20 ALLOW FILTERING", readOnly: true, returnsData: true},
		{name: "select distinct", cql: "SELECT DISTINCT partition_key FROM ks.events", readOnly: true, returnsData: true},
		{name: "select json", cql: "SELECT JSON id, name FROM ks.users WHERE id = 1", readOnly: true, returnsData: true},
		{name: "select with limit and per partition limit", cql: "SELECT * FROM ks.events PER PARTITION LIMIT 5 LIMIT 20", readOnly: true, returnsData: true},
		{name: "select with order by and group by", cql: "SELECT k, COUNT(*) FROM ks.events GROUP BY k ORDER BY c DESC", readOnly: true, returnsData: true},
		{name: "select multi-statement reads", cql: "SELECT * FROM ks.t1; SELECT id FROM ks.t2 WHERE id = 1;", readOnly: true, returnsData: true},

		// --- Strings, quoted identifiers, and comments ---
		{name: "semicolon and dml keywords in string literal", cql: "SELECT * FROM users WHERE name = 'a; DROP TABLE users; -- DELETE'", readOnly: true, returnsData: true},
		{name: "escaped single quote in string literal", cql: "SELECT * FROM users WHERE name = 'it''s; DROP TABLE users'", readOnly: true, returnsData: true},
		{name: "quoted identifiers named after write keywords", cql: `SELECT "insert", "update", "delete" FROM "drop"`, readOnly: true, returnsData: true},
		{name: "dash, slash, and block comments", cql: "-- DROP TABLE t\n// DELETE FROM t\n/* UPDATE t SET a = 1 */\nSELECT * FROM users", readOnly: true, returnsData: true},
		{name: "trailing comments after select", cql: "SELECT * FROM users; -- done\n/* block */", readOnly: true, returnsData: true},

		// --- Empty / whitespace / comment-only ---
		{name: "empty input", cql: "", readOnly: true, returnsData: true},
		{name: "whitespace only", cql: "   \n\t  ", readOnly: true, returnsData: true},
		{name: "dash comment only", cql: "-- just a comment\n", readOnly: true, returnsData: true},
		{name: "slash comment only", cql: "// just a comment\n", readOnly: true, returnsData: true},
		{name: "block comment only", cql: "/* just a comment */", readOnly: true, returnsData: true},
		{name: "semicolons only", cql: ";;;", readOnly: true, returnsData: true},

		// --- DML writes ---
		{name: "insert", cql: "INSERT INTO users (id, name) VALUES (1, 'alice')"},
		{name: "insert json", cql: `INSERT INTO users JSON '{"id": 1, "name": "alice"}'`},
		{name: "update", cql: "UPDATE users SET name = 'bob' WHERE id = 1"},
		{name: "delete", cql: "DELETE FROM users WHERE id = 1"},
		{name: "truncate", cql: "TRUNCATE TABLE users"},
		{name: "truncate without table keyword", cql: "TRUNCATE users"},
		{name: "logged batch", cql: "BEGIN BATCH INSERT INTO users (id) VALUES (1); APPLY BATCH"},
		{name: "unlogged batch", cql: "BEGIN UNLOGGED BATCH UPDATE users SET name = 'x' WHERE id = 1; APPLY BATCH"},
		{name: "counter batch", cql: "BEGIN COUNTER BATCH UPDATE stats SET hits = hits + 1 WHERE id = 1; APPLY BATCH"},

		// --- DDL writes ---
		{name: "create keyspace", cql: "CREATE KEYSPACE ks WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}"},
		{name: "alter keyspace", cql: "ALTER KEYSPACE ks WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3}"},
		{name: "drop keyspace", cql: "DROP KEYSPACE ks"},
		{name: "create table", cql: "CREATE TABLE users (id int PRIMARY KEY, name text)"},
		{name: "alter table", cql: "ALTER TABLE users ADD email text"},
		{name: "drop table", cql: "DROP TABLE users"},
		{name: "create index", cql: "CREATE INDEX idx_name ON users (name)"},
		{name: "drop index", cql: "DROP INDEX idx_name"},
		{name: "create type", cql: "CREATE TYPE address (street text, zip int)"},
		{name: "drop type", cql: "DROP TYPE address"},
		{name: "create materialized view", cql: "CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM users WHERE id IS NOT NULL AND name IS NOT NULL PRIMARY KEY (name, id)"},
		{name: "drop materialized view", cql: "DROP MATERIALIZED VIEW mv"},

		// --- DCL, session, and non-SELECT metadata commands ---
		{name: "use keyspace", cql: "USE ks"},
		{name: "grant permission", cql: "GRANT SELECT ON KEYSPACE ks TO analyst"},
		{name: "revoke permission", cql: "REVOKE SELECT ON KEYSPACE ks FROM analyst"},
		{name: "grant role", cql: "GRANT analyst TO bob"},
		{name: "revoke role", cql: "REVOKE analyst FROM bob"},
		{name: "create role", cql: "CREATE ROLE analyst WITH LOGIN = true"},
		{name: "drop role", cql: "DROP ROLE analyst"},
		{name: "create user", cql: "CREATE USER bob WITH PASSWORD 'secret'"},
		{name: "drop user", cql: "DROP USER bob"},
		{name: "describe tables", cql: "DESCRIBE TABLES"},
		{name: "list roles", cql: "LIST ROLES"},
		{name: "list users", cql: "LIST USERS"},
		{name: "list permissions", cql: "LIST ALL PERMISSIONS"},

		// --- Multi-statement injection attempts ---
		{name: "select then drop", cql: "SELECT * FROM employee; DROP TABLE employee"},
		{name: "select then insert", cql: "SELECT * FROM users; INSERT INTO users (id) VALUES (1)"},
		{name: "select then update", cql: "SELECT * FROM users; UPDATE users SET name = 'x' WHERE id = 1"},
		{name: "select then delete", cql: "SELECT * FROM users; DELETE FROM users WHERE id = 1"},
		{name: "select then batch", cql: "SELECT * FROM users; BEGIN BATCH INSERT INTO users (id) VALUES (1); APPLY BATCH"},
		{name: "select then use", cql: "SELECT * FROM users; USE other_ks"},
		{name: "use then select", cql: "USE other_ks; SELECT * FROM users"},
		{name: "write then select", cql: "DROP TABLE users; SELECT * FROM users"},

		// --- Syntax and lex errors ---
		{name: "truncated select", cql: "SELECT * FROM", wantErr: true},
		{name: "bare select", cql: "SELECT", wantErr: true},
		{name: "unknown keyword", cql: "EXPLAIN SELECT * FROM users", wantErr: true},
		{name: "unterminated string", cql: "SELECT * FROM users WHERE name = 'unterminated", wantErr: true},
		{name: "unterminated quoted ident", cql: `SELECT "col FROM users`, wantErr: true},
		{name: "unterminated block comment", cql: "/*", wantErr: true},
		{name: "unterminated block comment with text", cql: "/* unterminated", wantErr: true},
		{name: "unterminated block comment after select", cql: "SELECT * FROM users; /*", wantErr: true},
		{name: "select then syntax error", cql: "SELECT * FROM users; SELECT FROM", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tc.cql)
			if tc.wantErr {
				require.Error(t, err, "cql: %s", tc.cql)
				require.False(t, readOnly)
				require.False(t, returnsData)
				var perr *parser.ParseError
				require.ErrorAs(t, err, &perr)
				return
			}
			require.NoError(t, err, "cql: %s", tc.cql)
			require.Equal(t, tc.readOnly, readOnly, "readOnly for %q", tc.cql)
			require.Equal(t, tc.returnsData, returnsData, "returnsData for %q", tc.cql)
		})
	}
}

package cassandra

import (
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
)

func TestParseEmpty(t *testing.T) {
	stmts, err := Parse("")
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 0 {
		t.Fatalf("expected 0 statements, got %d", len(stmts))
	}
}

func TestParseBlank(t *testing.T) {
	stmts, err := Parse("  \n\t  ")
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 0 {
		t.Fatalf("expected 0 statements, got %d", len(stmts))
	}
}

func TestParseSelect(t *testing.T) {
	tests := []struct {
		input string
		check func(t *testing.T, s Statement)
	}{
		{
			input: "SELECT * FROM users",
			check: func(t *testing.T, s Statement) {
				sel, ok := s.AST.(*ast.SelectStmt)
				if !ok {
					t.Fatalf("expected *ast.SelectStmt, got %T", s.AST)
				}
				if sel.From == nil || len(sel.From.Parts) != 1 || sel.From.Parts[0].Name != "users" {
					t.Fatal("expected FROM users")
				}
			},
		},
		{
			input: "SELECT DISTINCT name, age FROM ks.users WHERE id = 1 LIMIT 10 ALLOW FILTERING",
			check: func(t *testing.T, s Statement) {
				sel := s.AST.(*ast.SelectStmt)
				if !sel.Distinct {
					t.Fatal("expected DISTINCT")
				}
				if len(sel.Elements) != 2 {
					t.Fatalf("expected 2 select elements, got %d", len(sel.Elements))
				}
				if sel.From == nil || len(sel.From.Parts) != 2 {
					t.Fatal("expected qualified table name ks.users")
				}
				if len(sel.Where) != 1 {
					t.Fatalf("expected 1 WHERE condition, got %d", len(sel.Where))
				}
				if sel.Limit == nil {
					t.Fatal("expected LIMIT")
				}
				if !sel.AllowFiltering {
					t.Fatal("expected ALLOW FILTERING")
				}
			},
		},
		{
			input: "SELECT JSON name AS n FROM users",
			check: func(t *testing.T, s Statement) {
				sel := s.AST.(*ast.SelectStmt)
				if !sel.JSON {
					t.Fatal("expected JSON")
				}
				if len(sel.Elements) != 1 {
					t.Fatalf("expected 1 element, got %d", len(sel.Elements))
				}
				if sel.Elements[0].Alias == nil || sel.Elements[0].Alias.Name != "n" {
					t.Fatal("expected alias 'n'")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmts, err := Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			tt.check(t, stmts[0])
		})
	}
}

func TestParseInsert(t *testing.T) {
	tests := []string{
		"INSERT INTO users (id, name) VALUES (1, 'Alice')",
		"INSERT INTO ks.users (id, name) VALUES (1, 'Bob') IF NOT EXISTS",
		"INSERT INTO users (id, name) VALUES (1, 'Charlie') USING TTL 86400",
		"INSERT INTO users JSON '{\"id\": 1, \"name\": \"Dave\"}'",
		"INSERT INTO users JSON '{\"id\": 1}' DEFAULT UNSET",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			stmts, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			if _, ok := stmts[0].AST.(*ast.InsertStmt); !ok {
				t.Fatalf("expected *ast.InsertStmt, got %T", stmts[0].AST)
			}
		})
	}
}

func TestParseUpdate(t *testing.T) {
	tests := []string{
		"UPDATE users SET name = 'Alice' WHERE id = 1",
		"UPDATE ks.users USING TTL 3600 SET name = 'Bob' WHERE id = 2",
		"UPDATE users SET name = 'Charlie' WHERE id = 3 IF EXISTS",
		"UPDATE users SET name = 'Dave' WHERE id = 4 IF name = 'old'",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			stmts, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			if _, ok := stmts[0].AST.(*ast.UpdateStmt); !ok {
				t.Fatalf("expected *ast.UpdateStmt, got %T", stmts[0].AST)
			}
		})
	}
}

func TestParseDelete(t *testing.T) {
	tests := []string{
		"DELETE FROM users WHERE id = 1",
		"DELETE name FROM ks.users WHERE id = 2",
		"DELETE FROM users WHERE id = 3 IF EXISTS",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			stmts, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			if _, ok := stmts[0].AST.(*ast.DeleteStmt); !ok {
				t.Fatalf("expected *ast.DeleteStmt, got %T", stmts[0].AST)
			}
		})
	}
}

func TestParseBatch(t *testing.T) {
	input := `BEGIN BATCH
		INSERT INTO users (id, name) VALUES (1, 'Alice');
		UPDATE users SET name = 'Bob' WHERE id = 2;
		DELETE FROM users WHERE id = 3;
	APPLY BATCH`
	stmts, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	batch, ok := stmts[0].AST.(*ast.BatchStmt)
	if !ok {
		t.Fatalf("expected *ast.BatchStmt, got %T", stmts[0].AST)
	}
	if len(batch.Statements) != 3 {
		t.Fatalf("expected 3 inner statements, got %d", len(batch.Statements))
	}
}

func TestParseDDL(t *testing.T) {
	tests := []struct {
		input    string
		nodeType string
	}{
		{"CREATE KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'}", "CreateKeyspaceStmt"},
		{"ALTER KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '3'}", "AlterKeyspaceStmt"},
		{"DROP KEYSPACE IF EXISTS ks", "DropKeyspaceStmt"},
		{"CREATE TABLE users (id int, name text, PRIMARY KEY (id))", "CreateTableStmt"},
		{"ALTER TABLE users ADD email text", "AlterTableStmt"},
		{"DROP TABLE IF EXISTS users", "DropTableStmt"},
		{"CREATE INDEX ON users (name)", "CreateIndexStmt"},
		{"DROP INDEX IF EXISTS users_name_idx", "DropIndexStmt"},
		{"CREATE TYPE address (street text, city text)", "CreateTypeStmt"},
		{"ALTER TYPE address ADD zip text", "AlterTypeStmt"},
		{"DROP TYPE IF EXISTS address", "DropTypeStmt"},
		{"TRUNCATE users", "TruncateStmt"},
		{"TRUNCATE TABLE ks.users", "TruncateStmt"},
		{"USE my_keyspace", "UseStmt"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmts, err := Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
		})
	}
}

func TestParseAuth(t *testing.T) {
	tests := []string{
		"GRANT SELECT ON TABLE users TO reader",
		"REVOKE ALL ON ALL KEYSPACES FROM admin",
		"LIST ALL PERMISSIONS OF admin",
		"LIST ROLES",
		"LIST ROLES OF admin NORECURSIVE",
		"CREATE ROLE myrole WITH PASSWORD = 'secret' AND LOGIN = true",
		"ALTER ROLE myrole WITH PASSWORD = 'newsecret'",
		"DROP ROLE IF EXISTS myrole",
		"CREATE USER myuser WITH PASSWORD 'secret' SUPERUSER",
		"DROP USER IF EXISTS myuser",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			stmts, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
		})
	}
}

func TestParseMultipleStatements(t *testing.T) {
	input := "SELECT * FROM users; INSERT INTO users (id) VALUES (1); USE ks"
	stmts, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(stmts))
	}
	if _, ok := stmts[0].AST.(*ast.SelectStmt); !ok {
		t.Fatalf("stmt 0: expected SelectStmt, got %T", stmts[0].AST)
	}
	if _, ok := stmts[1].AST.(*ast.InsertStmt); !ok {
		t.Fatalf("stmt 1: expected InsertStmt, got %T", stmts[1].AST)
	}
	if _, ok := stmts[2].AST.(*ast.UseStmt); !ok {
		t.Fatalf("stmt 2: expected UseStmt, got %T", stmts[2].AST)
	}
}

func TestParsePositions(t *testing.T) {
	input := "SELECT * FROM users"
	stmts, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	s := stmts[0]
	if s.ByteStart != 0 {
		t.Errorf("ByteStart = %d, want 0", s.ByteStart)
	}
	if s.ByteEnd != 19 {
		t.Errorf("ByteEnd = %d, want 19", s.ByteEnd)
	}
	if s.Start.Line != 1 || s.Start.Column != 1 {
		t.Errorf("Start = %+v, want {1 1}", s.Start)
	}
	if s.Text != "SELECT * FROM users" {
		t.Errorf("Text = %q", s.Text)
	}
}

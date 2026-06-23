package cassandra

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/cassandra/parser"
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
		{"CREATE TABLE t (id int PRIMARY KEY, v vector<float, 3>)", "CreateTableStmt"},
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

func TestParseMV(t *testing.T) {
	tests := []struct {
		input string
		check func(t *testing.T, s Statement)
	}{
		{
			input: "CREATE MATERIALIZED VIEW mv AS SELECT * FROM users WHERE id IS NOT NULL PRIMARY KEY (id)",
			check: func(t *testing.T, s Statement) {
				mv, ok := s.AST.(*ast.CreateMVStmt)
				if !ok {
					t.Fatalf("expected *ast.CreateMVStmt, got %T", s.AST)
				}
				if !mv.SelectAll {
					t.Fatal("expected SelectAll = true")
				}
				if len(mv.SelectColumns) != 0 {
					t.Fatalf("expected 0 SelectColumns, got %d", len(mv.SelectColumns))
				}
				if len(mv.WhereNotNull) != 1 || mv.WhereNotNull[0].Name != "id" {
					t.Fatal("expected WHERE id IS NOT NULL")
				}
			},
		},
		{
			input: "CREATE MATERIALIZED VIEW mv AS SELECT col1, col2 FROM users WHERE col1 IS NOT NULL AND col2 IS NOT NULL PRIMARY KEY (col1, col2)",
			check: func(t *testing.T, s Statement) {
				mv := s.AST.(*ast.CreateMVStmt)
				if mv.SelectAll {
					t.Fatal("expected SelectAll = false")
				}
				if len(mv.SelectColumns) != 2 {
					t.Fatalf("expected 2 SelectColumns, got %d", len(mv.SelectColumns))
				}
				if len(mv.WhereNotNull) != 2 {
					t.Fatalf("expected 2 WhereNotNull, got %d", len(mv.WhereNotNull))
				}
			},
		},
		{
			input: "CREATE MATERIALIZED VIEW IF NOT EXISTS ks.mv AS SELECT * FROM users WHERE id IS NOT NULL PRIMARY KEY (id) WITH comment = 'test'",
			check: func(t *testing.T, s Statement) {
				mv := s.AST.(*ast.CreateMVStmt)
				if !mv.IfNotExists {
					t.Fatal("expected IfNotExists")
				}
				if len(mv.Name.Parts) != 2 {
					t.Fatal("expected qualified name ks.mv")
				}
				if len(mv.Options) != 1 {
					t.Fatalf("expected 1 option, got %d", len(mv.Options))
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

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// IF NOT EXISTS strict validation
		{"IF NOT GARBAGE", "INSERT INTO users (id) VALUES (1) IF NOT GARBAGE"},
		{"IF NOT without EXISTS in CREATE", "CREATE TABLE IF NOT GARBAGE users (id int PRIMARY KEY)"},

		// Truncated/malformed DML
		{"truncated SELECT", "SELECT"},
		{"truncated SELECT FROM", "SELECT * FROM"},
		{"truncated INSERT", "INSERT INTO"},
		{"truncated INSERT no VALUES", "INSERT INTO users (id)"},
		{"truncated UPDATE", "UPDATE"},
		{"truncated UPDATE no SET", "UPDATE users"},
		{"truncated DELETE", "DELETE FROM"},
		{"truncated BATCH", "BEGIN BATCH"},

		// Truncated/malformed DDL
		{"truncated CREATE TABLE", "CREATE TABLE"},
		{"truncated CREATE KEYSPACE", "CREATE KEYSPACE"},
		{"truncated DROP", "DROP"},
		{"truncated ALTER", "ALTER"},
		{"CREATE without object", "CREATE"},

		// Invalid tokens
		{"bare operator", "< >"},
		{"invalid statement start", "123"},

		// MV IS NOT NULL with wrong tokens
		{"MV IS LOL NOPE", "CREATE MATERIALIZED VIEW mv AS SELECT * FROM t WHERE id IS LOL NOPE PRIMARY KEY (id)"},

		// Type generic validation
		{"vector missing element type", "CREATE TABLE t (id int PRIMARY KEY, v vector<3>)"},
		{"vector extra dimension", "CREATE TABLE t (id int PRIMARY KEY, v vector<float, 3, 4>)"},
		{"map with integer param", "CREATE TABLE t (id int PRIMARY KEY, m map<text, 3>)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tt.input)
			}
		})
	}
}

func TestParseNoPanic(t *testing.T) {
	inputs := []string{
		"",
		" ",
		"\t\n",
		";",
		";;;",
		"SELECT",
		"INSERT",
		"CREATE",
		"DROP TABLE",
		"ALTER TABLE users",
		"BEGIN BATCH APPLY BATCH",
		"SELECT * FROM users WHERE",
		"CREATE TABLE t (",
		"UPDATE users SET name =",
		"'unterminated string",
		`"unterminated ident`,
		"$$unterminated code block",
		"/* unterminated block comment",
		"SELECT 1e",
		"SELECT 1e+",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Parse(%q) panicked: %v", input, r)
				}
			}()
			Parse(input) // error is fine, panic is not
		})
	}
}

func TestParseLocWalker(t *testing.T) {
	tests := []string{
		"SELECT * FROM users",
		"SELECT name, age FROM ks.users WHERE id = 1",
		"INSERT INTO users (id, name) VALUES (1, 'Alice')",
		"UPDATE users SET name = 'Bob' WHERE id = 2",
		"DELETE FROM users WHERE id = 3",
		"CREATE TABLE t (id int, name text, PRIMARY KEY (id))",
		"CREATE KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'}",
		"DROP TABLE IF EXISTS users",
		"USE my_keyspace",
		"TRUNCATE TABLE users",
		"CREATE MATERIALIZED VIEW mv AS SELECT * FROM users WHERE id IS NOT NULL PRIMARY KEY (id)",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			stmts, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range stmts {
				if s.ByteStart < 0 {
					t.Errorf("ByteStart = %d, want >= 0", s.ByteStart)
				}
				if s.ByteEnd <= s.ByteStart {
					t.Errorf("ByteEnd = %d <= ByteStart = %d", s.ByteEnd, s.ByteStart)
				}
				if s.ByteEnd > len(input) {
					t.Errorf("ByteEnd = %d > len(input) = %d", s.ByteEnd, len(input))
				}
				text := input[s.ByteStart:s.ByteEnd]
				if text != s.Text {
					t.Errorf("input[%d:%d] = %q, s.Text = %q", s.ByteStart, s.ByteEnd, text, s.Text)
				}
				if s.Start.Line < 1 || s.Start.Column < 1 {
					t.Errorf("Start = %+v, want line >= 1 and column >= 1", s.Start)
				}
				loc := s.AST.GetLoc()
				if loc.Start < 0 {
					t.Errorf("AST Loc.Start = %d, want >= 0", loc.Start)
				}
				if loc.End <= loc.Start {
					t.Errorf("AST Loc.End = %d <= Loc.Start = %d", loc.End, loc.Start)
				}
			}
		})
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

// ---------------------------------------------------------------------------
// Phase 4: L3 Error Quality
// ---------------------------------------------------------------------------

func TestErrorLineColumn(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantLine int
		wantCol  int
		wantNear string
	}{
		{
			name:     "single line error",
			sql:      "SELECT * FORM users",
			wantLine: 1,
			wantCol:  10,
			wantNear: "FORM",
		},
		{
			name:     "second line error",
			sql:      "SELECT *\nFORM users",
			wantLine: 2,
			wantCol:  1,
			wantNear: "FORM",
		},
		{
			name:     "deep in statement",
			sql:      "INSERT INTO users (id) VALUE (1)",
			wantLine: 1,
			wantCol:  24,
			wantNear: "VALUE",
		},
		{
			name:     "unterminated string",
			sql:      "SELECT 'abc",
			wantLine: 1,
			wantCol:  8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if err == nil {
				t.Fatal("expected error")
			}
			var pe *parser.ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("expected *parser.ParseError, got %T: %v", err, err)
			}
			if pe.Line != tt.wantLine {
				t.Errorf("Line = %d, want %d (error: %s)", pe.Line, tt.wantLine, pe.Error())
			}
			if pe.Column != tt.wantCol {
				t.Errorf("Column = %d, want %d (error: %s)", pe.Column, tt.wantCol, pe.Error())
			}
			if tt.wantNear != "" && pe.Near != tt.wantNear {
				t.Errorf("Near = %q, want %q (error: %s)", pe.Near, tt.wantNear, pe.Error())
			}
			if !strings.Contains(pe.Error(), "line") {
				t.Errorf("error message missing 'line': %s", pe.Error())
			}
			if !strings.Contains(pe.Error(), "column") {
				t.Errorf("error message missing 'column': %s", pe.Error())
			}
		})
	}
}

func TestErrorAtOrNear(t *testing.T) {
	_, err := Parse("SELECT * FORM users")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "at or near") {
		t.Errorf("error message missing 'at or near': %s", msg)
	}
	if !strings.Contains(msg, "FORM") {
		t.Errorf("error message missing token text 'FORM': %s", msg)
	}
}

func TestTruncationFuzz(t *testing.T) {
	validSQL := []string{
		"SELECT * FROM users WHERE id = 1 ORDER BY name ASC LIMIT 10",
		"INSERT INTO users (id, name) VALUES (1, 'Alice') IF NOT EXISTS USING TTL 86400",
		"UPDATE users USING TTL 3600 SET name = 'Bob' WHERE id = 2 IF name = 'old'",
		"DELETE name FROM ks.users WHERE id = 2 IF EXISTS",
		"BEGIN UNLOGGED BATCH USING TIMESTAMP 12345 INSERT INTO t (id) VALUES (1); DELETE FROM t WHERE id = 2; APPLY BATCH",
		"CREATE TABLE t (id int, name text, age int, PRIMARY KEY ((id, name), age)) WITH CLUSTERING ORDER BY (age DESC) AND comment = 'test'",
		"CREATE KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'} AND DURABLE_WRITES = true",
		"CREATE MATERIALIZED VIEW mv AS SELECT col1, col2 FROM t WHERE col1 IS NOT NULL AND col2 IS NOT NULL PRIMARY KEY (col1, col2)",
		"CREATE FUNCTION ks.f(input text) CALLED ON NULL INPUT RETURNS text LANGUAGE java AS $$return input;$$",
		"CREATE AGGREGATE ks.agg(int) SFUNC plus STYPE int FINALFUNC fin INITCOND 0",
		"GRANT SELECT ON TABLE users TO reader",
		"CREATE ROLE myrole WITH PASSWORD = 'secret' AND LOGIN = true AND SUPERUSER = false",
	}
	for _, sql := range validSQL {
		for i := 0; i <= len(sql); i++ {
			truncated := sql[:i]
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("Parse(%q) panicked (truncated from %q at byte %d): %v", truncated, sql, i, r)
					}
				}()
				Parse(truncated)
			}()
		}
	}
}

func TestBinaryInputNoPanic(t *testing.T) {
	inputs := []string{
		"\x00\x01\x02\x03",
		string([]byte{0xFF, 0xFE, 0xFD}),
		"\x00SELECT * FROM users",
		"SELECT\x00FROM\x00users",
		string(make([]byte, 1024)),
	}
	for _, input := range inputs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Parse(binary) panicked: %v", r)
				}
			}()
			Parse(input)
		}()
	}
}

func TestParseAlterIfExists(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		check func(t *testing.T, s Statement)
	}{
		{
			name: "ALTER KEYSPACE IF EXISTS",
			sql:  "ALTER KEYSPACE IF EXISTS ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'}",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterKeyspaceStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
			},
		},
		{
			name: "ALTER TABLE IF EXISTS",
			sql:  "ALTER TABLE IF EXISTS t ADD col text",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTableStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
			},
		},
		{
			name: "ALTER TABLE ADD IF NOT EXISTS",
			sql:  "ALTER TABLE t ADD IF NOT EXISTS col text",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTableStmt)
				if !stmt.AddIfNotExists {
					t.Fatal("expected AddIfNotExists=true")
				}
			},
		},
		{
			name: "ALTER TABLE DROP IF EXISTS",
			sql:  "ALTER TABLE t DROP IF EXISTS col",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTableStmt)
				if !stmt.DropIfExists {
					t.Fatal("expected DropIfExists=true")
				}
			},
		},
		{
			name: "ALTER TABLE RENAME IF EXISTS",
			sql:  "ALTER TABLE t RENAME IF EXISTS a TO b",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTableStmt)
				if !stmt.RenameIfExists {
					t.Fatal("expected RenameIfExists=true")
				}
			},
		},
		{
			name: "ALTER TABLE IF EXISTS + ADD IF NOT EXISTS combined",
			sql:  "ALTER TABLE IF EXISTS t ADD IF NOT EXISTS col text",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTableStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
				if !stmt.AddIfNotExists {
					t.Fatal("expected AddIfNotExists=true")
				}
			},
		},
		{
			name: "ALTER TYPE IF EXISTS",
			sql:  "ALTER TYPE IF EXISTS mytype ADD f2 int",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTypeStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
			},
		},
		{
			name: "ALTER TYPE ADD IF NOT EXISTS",
			sql:  "ALTER TYPE mytype ADD IF NOT EXISTS f2 int",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTypeStmt)
				if !stmt.AddIfNotExists {
					t.Fatal("expected AddIfNotExists=true")
				}
			},
		},
		{
			name: "ALTER TYPE RENAME IF EXISTS",
			sql:  "ALTER TYPE mytype RENAME IF EXISTS f1 TO field1",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTypeStmt)
				if !stmt.RenameIfExists {
					t.Fatal("expected RenameIfExists=true")
				}
			},
		},
		{
			name: "ALTER MATERIALIZED VIEW IF EXISTS",
			sql:  "ALTER MATERIALIZED VIEW IF EXISTS mv WITH comment = 'test'",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterMVStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
			},
		},
		{
			name: "ALTER ROLE IF EXISTS",
			sql:  "ALTER ROLE IF EXISTS r WITH PASSWORD = 'x'",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterRoleStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
			},
		},
		{
			name: "ALTER USER IF EXISTS",
			sql:  "ALTER USER IF EXISTS u WITH PASSWORD 'y'",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterUserStmt)
				if !stmt.IfExists {
					t.Fatal("expected IfExists=true")
				}
			},
		},
		{
			name: "ALTER TABLE without IF EXISTS",
			sql:  "ALTER TABLE t ADD col text",
			check: func(t *testing.T, s Statement) {
				stmt := s.AST.(*ast.AlterTableStmt)
				if stmt.IfExists {
					t.Fatal("expected IfExists=false")
				}
				if stmt.AddIfNotExists {
					t.Fatal("expected AddIfNotExists=false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := Parse(tt.sql)
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

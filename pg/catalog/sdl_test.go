package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

func TestSDLValidation(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string // if non-empty, expect error containing this
	}{
		// ---- Accepted statements ----
		{
			name: "empty string returns empty catalog",
			sql:  "",
		},
		{
			name: "valid CREATE TABLE",
			sql:  "CREATE TABLE t (id int);",
		},
		{
			name: "valid CREATE VIEW",
			sql:  "CREATE TABLE t (id int); CREATE VIEW v AS SELECT id FROM t;",
		},
		{
			name: "valid CREATE FUNCTION",
			sql:  "CREATE FUNCTION f() RETURNS void AS $$ BEGIN END; $$ LANGUAGE plpgsql;",
		},
		{
			name: "valid CREATE INDEX",
			sql:  "CREATE TABLE t (id int); CREATE INDEX idx ON t (id);",
		},
		{
			name: "valid CREATE SEQUENCE",
			sql:  "CREATE SEQUENCE s;",
		},
		{
			name: "valid CREATE SCHEMA",
			sql:  "CREATE SCHEMA myschema;",
		},
		{
			name: "valid CREATE TYPE enum",
			sql:  "CREATE TYPE mood AS ENUM ('sad', 'happy');",
		},
		{
			name: "valid CREATE DOMAIN",
			sql:  "CREATE DOMAIN posint AS integer CHECK (VALUE > 0);",
		},
		{
			name: "valid CREATE TYPE composite",
			sql:  "CREATE TYPE pair AS (x int, y int);",
		},
		{
			name: "valid CREATE TYPE range",
			sql:  "CREATE TYPE floatrange AS RANGE (subtype = float8);",
		},
		{
			name: "valid CREATE EXTENSION",
			sql:  "CREATE EXTENSION IF NOT EXISTS pgcrypto;",
		},
		{
			name: "valid CREATE TRIGGER",
			sql: `CREATE TABLE t (id int);
CREATE FUNCTION trg_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();`,
		},
		{
			name: "valid CREATE POLICY",
			sql:  "CREATE TABLE t (id int); ALTER TABLE t ENABLE ROW LEVEL SECURITY; CREATE POLICY p ON t USING (true);",
		},
		{
			name: "valid CREATE MATERIALIZED VIEW",
			sql:  "CREATE TABLE t (id int); CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;",
		},
		{
			name: "valid CREATE CAST",
			sql:  "CREATE FUNCTION int4_to_text(integer) RETURNS text AS $$ SELECT $1::text; $$ LANGUAGE SQL; CREATE CAST (integer AS text) WITH FUNCTION int4_to_text(integer);",
		},
		{
			name: "valid CREATE FOREIGN TABLE",
			sql:  "CREATE FOREIGN TABLE ft (id int) SERVER myserver;",
		},
		{
			name: "valid COMMENT ON",
			sql:  "CREATE TABLE t (id int); COMMENT ON TABLE t IS 'my table';",
		},
		{
			name: "valid GRANT",
			sql:  "CREATE TABLE t (id int); GRANT SELECT ON t TO PUBLIC;",
		},
		{
			name: "ALTER SEQUENCE OWNED BY",
			sql:  "CREATE TABLE t (id int); CREATE SEQUENCE s; ALTER SEQUENCE s OWNED BY t.id;",
		},
		{
			name: "ALTER TABLE ENABLE ROW LEVEL SECURITY",
			sql:  "CREATE TABLE t (id int); ALTER TABLE t ENABLE ROW LEVEL SECURITY;",
		},
		{
			name: "ALTER TYPE ADD VALUE",
			sql:  "CREATE TYPE mood AS ENUM ('sad'); ALTER TYPE mood ADD VALUE 'happy';",
		},

		// ---- Rejected statements ----
		{
			name:    "INSERT rejected",
			sql:     "INSERT INTO t VALUES (1);",
			wantErr: "SDL does not allow INSERT statements",
		},
		{
			name:    "UPDATE rejected",
			sql:     "UPDATE t SET x = 1;",
			wantErr: "SDL does not allow UPDATE statements",
		},
		{
			name:    "DELETE rejected",
			sql:     "DELETE FROM t;",
			wantErr: "SDL does not allow DELETE statements",
		},
		{
			name:    "DROP TABLE rejected",
			sql:     "DROP TABLE t;",
			wantErr: "SDL does not allow DROP statements",
		},
		{
			name:    "ALTER TABLE ADD COLUMN rejected",
			sql:     "ALTER TABLE t ADD COLUMN x int;",
			wantErr: "SDL does not allow ALTER TABLE ADD/DROP COLUMN",
		},
		{
			name:    "ALTER TABLE DROP COLUMN rejected",
			sql:     "ALTER TABLE t DROP COLUMN x;",
			wantErr: "SDL does not allow ALTER TABLE ADD/DROP COLUMN",
		},
		{
			name:    "TRUNCATE rejected",
			sql:     "TRUNCATE t;",
			wantErr: "SDL does not allow TRUNCATE statements",
		},
		{
			name:    "DO block rejected",
			sql:     "DO $$ BEGIN END; $$;",
			wantErr: "SDL does not allow DO statements",
		},
		{
			name:    "parse error returns error",
			sql:     "CREATE TABL t (id int);",
			wantErr: "", // non-empty error from parser, tested separately below
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Special case: parse error test
			if tt.name == "parse error returns error" {
				c, err := LoadSDL(tt.sql)
				if err == nil {
					t.Fatal("expected parse error, got nil")
				}
				_ = c
				return
			}

			c, err := LoadSDL(tt.sql)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if c == nil {
				t.Fatal("expected non-nil catalog")
			}
		})
	}
}

// TestSDLNameResolution tests section 1.3: declared object collection and name resolution.
func TestSDLNameResolution(t *testing.T) {
	// Helper: parse SQL into bare statements.
	parseStmts := func(t *testing.T, sql string) []nodes.Node {
		t.Helper()
		list, err := pgparser.Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		stmts := make([]nodes.Node, 0, len(list.Items))
		for _, item := range list.Items {
			if raw, ok := item.(*nodes.RawStmt); ok {
				stmts = append(stmts, raw.Stmt)
			} else {
				stmts = append(stmts, item)
			}
		}
		return stmts
	}

	t.Run("unqualified names default to public schema", func(t *testing.T) {
		stmts := parseStmts(t, "CREATE TABLE users (id int);")
		declared := collectDeclaredObjects(stmts)
		if !declared["public.users"] {
			t.Fatalf("expected public.users in declared set, got %v", declared)
		}
	})

	t.Run("schema-qualified names preserved", func(t *testing.T) {
		stmts := parseStmts(t, "CREATE SCHEMA myapp; CREATE TABLE myapp.users (id int);")
		declared := collectDeclaredObjects(stmts)
		if !declared["myapp.users"] {
			t.Fatalf("expected myapp.users in declared set, got %v", declared)
		}
	})

	t.Run("same name in different schemas treated as distinct", func(t *testing.T) {
		stmts := parseStmts(t, `
			CREATE SCHEMA a;
			CREATE SCHEMA b;
			CREATE TABLE a.t (id int);
			CREATE TABLE b.t (id int);
		`)
		declared := collectDeclaredObjects(stmts)
		if !declared["a.t"] {
			t.Fatalf("expected a.t in declared set, got %v", declared)
		}
		if !declared["b.t"] {
			t.Fatalf("expected b.t in declared set, got %v", declared)
		}
		// They should be distinct entries.
		if declared["a.t"] == declared["b.t"] {
			// Both true, which is correct — they are different keys.
		}
	})

	t.Run("function identity includes argument types", func(t *testing.T) {
		stmts := parseStmts(t, `
			CREATE FUNCTION myfunc(integer, text) RETURNS void AS $$ BEGIN END; $$ LANGUAGE plpgsql;
			CREATE FUNCTION myfunc(text) RETURNS void AS $$ BEGIN END; $$ LANGUAGE plpgsql;
		`)
		declared := collectDeclaredObjects(stmts)
		if !declared["public.myfunc(int4,text)"] {
			t.Fatalf("expected public.myfunc(int4,text) in declared set, got %v", declared)
		}
		if !declared["public.myfunc(text)"] {
			t.Fatalf("expected public.myfunc(text) in declared set, got %v", declared)
		}
	})

	t.Run("built-in function not in declared set", func(t *testing.T) {
		// now() and count() are built-in functions — they should NOT appear in declared set.
		stmts := parseStmts(t, "CREATE TABLE t (id int, created_at timestamptz DEFAULT now());")
		declared := collectDeclaredObjects(stmts)
		for k := range declared {
			if strings.Contains(k, "now") || strings.Contains(k, "count") {
				t.Fatalf("built-in function should not be in declared set: %s", k)
			}
		}
	})

	t.Run("built-in type not in declared set", func(t *testing.T) {
		// integer, text are built-in types — should NOT appear in declared set.
		stmts := parseStmts(t, "CREATE TABLE t (id integer, name text);")
		declared := collectDeclaredObjects(stmts)
		for k := range declared {
			if k == "pg_catalog.int4" || k == "pg_catalog.text" || k == "public.integer" || k == "public.text" {
				t.Fatalf("built-in type should not be in declared set: %s", k)
			}
		}
	})

	t.Run("undeclared reference passes through to ProcessUtility", func(t *testing.T) {
		// A table references an undeclared type — should not error in collectDeclaredObjects.
		// The error should come from ProcessUtility if the type doesn't exist.
		stmts := parseStmts(t, `
			CREATE TABLE t (id integer, status mood);
		`)
		declared := collectDeclaredObjects(stmts)
		// "mood" is not declared, so it should not be in the set.
		if declared["public.mood"] {
			t.Fatalf("undeclared type 'mood' should not be in declared set")
		}
		// Only the table should be declared.
		if !declared["public.t"] {
			t.Fatalf("expected public.t in declared set")
		}
	})
}

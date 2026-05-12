package catalog

import (
	"strings"
	"testing"
)

// TestSDLExprDepsColumnDefault tests section 2.1: Column DEFAULT and CHECK expression dependencies.
func TestSDLExprDepsColumnDefault(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "DEFAULT nextval creates dependency on sequence",
			sql: `
				CREATE TABLE t (
					id integer DEFAULT nextval('my_seq')
				);
				CREATE SEQUENCE my_seq;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "DEFAULT my_function creates dependency on function",
			sql: `
				CREATE TABLE t (
					val integer DEFAULT gen_id()
				);
				CREATE FUNCTION gen_id() RETURNS integer AS $$ SELECT 1; $$ LANGUAGE sql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "CHECK with validate_func creates dependency on function",
			sql: `
				CREATE TABLE t (
					val integer CHECK (is_positive(val))
				);
				CREATE FUNCTION is_positive(integer) RETURNS boolean AS $$ SELECT $1 > 0; $$ LANGUAGE sql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "CHECK with subquery is rejected like PostgreSQL",
			sql: `
				CREATE TABLE t (
					val integer,
					CONSTRAINT chk CHECK (val IN (SELECT id FROM lookup))
				);
				CREATE TABLE lookup (id integer);
			`,
			wantErr: "cannot use subquery in check constraint",
		},
		{
			name: "DEFAULT with type cast to user type creates dependency",
			sql: `
				CREATE TABLE t (
					status text DEFAULT 'active'::mood
				);
				CREATE TYPE mood AS ENUM ('active', 'inactive');
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "multiple defaults referencing different functions",
			sql: `
				CREATE TABLE t (
					a integer DEFAULT gen_a(),
					b integer DEFAULT gen_b()
				);
				CREATE FUNCTION gen_a() RETURNS integer AS $$ SELECT 1; $$ LANGUAGE sql;
				CREATE FUNCTION gen_b() RETURNS integer AS $$ SELECT 2; $$ LANGUAGE sql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "column default with built-in function does NOT create dependency",
			sql: `
				CREATE TABLE t (
					id integer,
					created_at timestamptz DEFAULT now()
				);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

// TestSDLExprDepsViewQuery tests section 2.2: View query dependencies.
func TestSDLExprDepsViewQuery(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "simple SELECT FROM table creates dependency",
			sql: `
				CREATE VIEW v AS SELECT id FROM t;
				CREATE TABLE t (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "JOIN across two tables creates dependencies on both",
			sql: `
				CREATE VIEW v AS SELECT a.id, b.name FROM a JOIN b ON a.id = b.id;
				CREATE TABLE a (id integer);
				CREATE TABLE b (id integer, name text);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "subquery in WHERE referencing another table",
			sql: `
				CREATE VIEW v AS SELECT id FROM a WHERE id IN (SELECT id FROM b);
				CREATE TABLE a (id integer);
				CREATE TABLE b (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "function call in SELECT list creates dependency",
			sql: `
				CREATE VIEW v AS SELECT format_name(name) FROM t;
				CREATE TABLE t (name text);
				CREATE FUNCTION format_name(text) RETURNS text AS $$ SELECT upper($1); $$ LANGUAGE sql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "CTE name NOT treated as external dependency",
			sql: `
				CREATE VIEW v AS
					WITH cte AS (SELECT id FROM t)
					SELECT id FROM cte;
				CREATE TABLE t (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "CTE body referencing real table creates dependency",
			sql: `
				CREATE VIEW v AS
					WITH cte AS (SELECT id FROM real_table)
					SELECT id FROM cte;
				CREATE TABLE real_table (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
				if c.GetRelation("public", "real_table") == nil {
					t.Fatal("table real_table not found")
				}
			},
		},
		{
			name: "nested subquery in FROM dependencies extracted",
			sql: `
				CREATE VIEW v AS SELECT x.id FROM (SELECT id FROM t) x;
				CREATE TABLE t (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "view referencing another view - dependency resolved",
			sql: `
				CREATE VIEW v2 AS SELECT id FROM v1;
				CREATE VIEW v1 AS SELECT id FROM t;
				CREATE TABLE t (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v2") == nil {
					t.Fatal("view v2 not found")
				}
				if c.GetRelation("public", "v1") == nil {
					t.Fatal("view v1 not found")
				}
			},
		},
		{
			name: "UNION across tables - all dependencies extracted",
			sql: `
				CREATE VIEW v AS SELECT id FROM a UNION SELECT id FROM b;
				CREATE TABLE a (id integer);
				CREATE TABLE b (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "view with type cast to user type creates dependency",
			sql: `
				CREATE VIEW v AS SELECT 'active'::mood AS status FROM t;
				CREATE TABLE t (id integer);
				CREATE TYPE mood AS ENUM ('active', 'inactive');
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

// TestSDLExprDepsIndexTriggerPolicy tests section 2.3: Index, Trigger, Policy, Domain expression dependencies.
func TestSDLExprDepsIndexTriggerPolicy(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "expression index with function creates dependency",
			sql: `
				CREATE TABLE t (name text);
				CREATE INDEX idx ON t (normalize_name(name));
				CREATE FUNCTION normalize_name(text) RETURNS text AS $$ SELECT lower($1); $$ LANGUAGE sql IMMUTABLE;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "partial index WHERE with function creates dependency",
			sql: `
				CREATE TABLE t (id integer, active boolean);
				CREATE INDEX idx ON t (id) WHERE is_active(active);
				CREATE FUNCTION is_active(boolean) RETURNS boolean AS $$ SELECT $1; $$ LANGUAGE sql IMMUTABLE;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "trigger WHEN clause with function creates dependency",
			sql: `
				CREATE TABLE t (id integer, val text);
				CREATE FUNCTION should_fire(text) RETURNS boolean AS $$ SELECT $1 IS NOT NULL; $$ LANGUAGE sql IMMUTABLE;
				CREATE FUNCTION trg_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
				CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW WHEN (should_fire(NEW.val)) EXECUTE FUNCTION trg_fn();
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "policy USING with function creates dependency",
			sql: `
				CREATE TABLE t (id integer, owner_id integer);
				ALTER TABLE t ENABLE ROW LEVEL SECURITY;
				CREATE FUNCTION current_uid() RETURNS integer AS $$ SELECT 1; $$ LANGUAGE sql;
				CREATE POLICY p ON t USING (owner_id = current_uid());
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "policy WITH CHECK with function creates dependency",
			sql: `
				CREATE TABLE t (id integer, owner_id integer);
				ALTER TABLE t ENABLE ROW LEVEL SECURITY;
				CREATE FUNCTION is_owner(integer) RETURNS boolean AS $$ SELECT $1 = 1; $$ LANGUAGE sql;
				CREATE POLICY p ON t FOR INSERT WITH CHECK (is_owner(owner_id));
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "domain CHECK with function creates dependency",
			sql: `
				CREATE DOMAIN positive_int AS integer CHECK (check_positive(VALUE));
				CREATE FUNCTION check_positive(integer) RETURNS boolean AS $$ SELECT $1 > 0; $$ LANGUAGE sql;
			`,
			check: func(t *testing.T, c *Catalog) {
				// Domain should be in the catalog (it was loaded successfully).
				// We just verify no error occurred, meaning dependencies resolved.
			},
		},
		{
			name: "trigger EXECUTE FUNCTION creates dependency on function",
			sql: `
				CREATE TABLE t (id integer);
				CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION my_trigger_fn();
				CREATE FUNCTION my_trigger_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "policy on table creates structural dependency",
			sql: `
				CREATE POLICY p ON t USING (true);
				CREATE TABLE t (id integer);
				ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "ALTER SEQUENCE OWNED BY creates dependency on table",
			sql: `
				CREATE SEQUENCE my_seq;
				ALTER SEQUENCE my_seq OWNED BY t.id;
				CREATE TABLE t (id integer);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

// TestSDLExprDepsFunctionType tests section 2.4: Function and Type dependencies.
func TestSDLExprDepsFunctionType(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "function parameter type referencing user type creates dependency",
			sql: `
				CREATE FUNCTION process(val mood) RETURNS void AS $$ BEGIN END; $$ LANGUAGE plpgsql;
				CREATE TYPE mood AS ENUM ('happy', 'sad');
			`,
			check: func(t *testing.T, c *Catalog) {
				// Function should be loaded successfully after enum is created first.
			},
		},
		{
			name: "function RETURNS user type creates dependency",
			sql: `
				CREATE FUNCTION get_mood() RETURNS mood AS $$ SELECT 'happy'::mood; $$ LANGUAGE sql;
				CREATE TYPE mood AS ENUM ('happy', 'sad');
			`,
			check: func(t *testing.T, c *Catalog) {
				// Function should be loaded after the enum.
			},
		},
		{
			name: "function RETURNS SETOF table creates dependency",
			sql: `
				CREATE FUNCTION get_all_users() RETURNS SETOF users AS $$ SELECT * FROM users; $$ LANGUAGE sql;
				CREATE TABLE users (id integer, name text);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "users") == nil {
					t.Fatal("table users not found")
				}
			},
		},
		{
			name: "function parameter DEFAULT with function creates dependency",
			sql: `
				CREATE FUNCTION process(val integer DEFAULT get_default()) RETURNS void AS $$ BEGIN END; $$ LANGUAGE plpgsql;
				CREATE FUNCTION get_default() RETURNS integer AS $$ SELECT 42; $$ LANGUAGE sql;
			`,
			check: func(t *testing.T, c *Catalog) {
				// Both functions should be loaded.
			},
		},
		{
			name: "composite type column referencing another user type",
			sql: `
				CREATE TYPE address AS (city text, zip zip_code);
				CREATE DOMAIN zip_code AS text CHECK (length(VALUE) = 5);
			`,
			check: func(t *testing.T, c *Catalog) {
				// Composite type should be loaded after domain.
			},
		},
		{
			name: "range type with user-defined subtype creates dependency",
			sql: `
				CREATE TYPE my_range AS RANGE (subtype = my_float);
				CREATE DOMAIN my_float AS float8 CHECK (VALUE >= 0);
			`,
			check: func(t *testing.T, c *Catalog) {
				// Range type should be loaded after domain.
			},
		},
		{
			name: "domain based on another domain creates dependency",
			sql: `
				CREATE DOMAIN positive_int AS base_int CHECK (VALUE > 0);
				CREATE DOMAIN base_int AS integer CHECK (VALUE IS NOT NULL);
			`,
			check: func(t *testing.T, c *Catalog) {
				// positive_int should be loaded after base_int.
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

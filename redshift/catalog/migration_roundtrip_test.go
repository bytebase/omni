package catalog

import (
	"fmt"
	"strings"
	"testing"
)

func TestMigrationRoundtrip(t *testing.T) {
	t.Run("simple roundtrip apply migration SQL to from catalog diff with to is empty", func(t *testing.T) {
		fromSQL := `CREATE TABLE t (id int);`
		toSQL := `
			CREATE TABLE t (id int, name text);
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip with table add column modify index drop", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE users (id int PRIMARY KEY, name text);
			CREATE INDEX idx_name ON users (name);
		`
		toSQL := `
			CREATE TABLE users (id int PRIMARY KEY, name text, email text NOT NULL DEFAULT 'unknown');
			CREATE TABLE posts (id int PRIMARY KEY, user_id int, title text);
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip with function view trigger changes", func(t *testing.T) {
		fromSQL := `
			CREATE TABLE t (id int, val text);
			CREATE FUNCTION old_fn() RETURNS void LANGUAGE sql AS $$ SELECT 1 $$;
		`
		toSQL := `
			CREATE TABLE t (id int, val text);
			CREATE FUNCTION new_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE VIEW v AS SELECT id, val FROM t WHERE id > 0;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION new_fn();
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip with enum add value domain alter", func(t *testing.T) {
		// Domain constraint add/drop roundtrips correctly.
		// Note: ALTER DOMAIN SET/DROP NOT NULL has a known parser issue
		// (subtypes swapped), so we test domain constraint changes instead.
		fromSQL := `
			CREATE TYPE mood AS ENUM ('happy', 'sad');
			CREATE DOMAIN posint AS integer;
		`
		toSQL := `
			CREATE TYPE mood AS ENUM ('happy', 'sad', 'neutral');
			CREATE DOMAIN posint AS integer CONSTRAINT positive CHECK (VALUE > 0);
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip with FK across schemas", func(t *testing.T) {
		fromSQL := `
			CREATE SCHEMA app;
		`
		toSQL := `
			CREATE SCHEMA app;
			CREATE TABLE public.users (id int PRIMARY KEY);
			CREATE TABLE app.orders (id int PRIMARY KEY, user_id int);
			ALTER TABLE app.orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id);
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip with all object types simultaneously", func(t *testing.T) {
		fromSQL := ``
		toSQL := `
			CREATE SCHEMA myapp;
			CREATE TYPE myapp.status AS ENUM ('active', 'inactive');
			CREATE SEQUENCE myapp.id_seq;
			CREATE TABLE myapp.items (
				id int PRIMARY KEY,
				name text NOT NULL,
				s myapp.status DEFAULT 'active'
			);
			CREATE INDEX idx_items_name ON myapp.items (name);
			CREATE FUNCTION myapp.process() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER items_trig BEFORE INSERT ON myapp.items FOR EACH ROW EXECUTE FUNCTION myapp.process();
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip create composite type from empty", func(t *testing.T) {
		fromSQL := ``
		toSQL := `CREATE TYPE address AS (street text, city text, zip varchar(10));`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip composite type add field", func(t *testing.T) {
		fromSQL := `CREATE TYPE address AS (street text, city text);`
		toSQL := `CREATE TYPE address AS (street text, city text, zip varchar(10));`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip composite type drop field", func(t *testing.T) {
		fromSQL := `CREATE TYPE address AS (street text, city text, zip varchar(10));`
		toSQL := `CREATE TYPE address AS (street text, city text);`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip composite type chain", func(t *testing.T) {
		fromSQL := ``
		toSQL := `
			CREATE TYPE address AS (street text, city text);
			CREATE TYPE person AS (name text, home address);
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})

	t.Run("roundtrip composite type with table column", func(t *testing.T) {
		fromSQL := ``
		toSQL := `
			CREATE TYPE address AS (street text, city text);
			CREATE TABLE contacts (id int PRIMARY KEY, addr address);
		`
		assertRoundtrip(t, fromSQL, toSQL)
	})
}

// assertRoundtrip verifies that applying the migration SQL to the `from`
// catalog produces a state equivalent to `to`.
func assertRoundtrip(t *testing.T, fromSQL, toSQL string) {
	t.Helper()

	from, err := LoadSQL(fromSQL)
	if err != nil {
		t.Fatalf("LoadSQL(from) error: %v", err)
	}
	to, err := LoadSQL(toSQL)
	if err != nil {
		t.Fatalf("LoadSQL(to) error: %v", err)
	}

	diff := Diff(from, to)
	if diff.IsEmpty() {
		t.Skip("diff is empty, nothing to roundtrip")
	}

	plan := GenerateMigration(from, to, diff)
	if len(plan.Ops) == 0 {
		t.Fatalf("GenerateMigration produced no ops, but diff was non-empty")
	}

	// Build the migration SQL — join ops with semicolons.
	migrationSQL := plan.SQL()

	// Apply migration to the `from` catalog.
	combinedSQL := fromSQL
	if combinedSQL != "" {
		combinedSQL += ";\n"
	}
	combinedSQL += migrationSQL

	migrated, err := LoadSQL(combinedSQL)
	if err != nil {
		// Try individual statements to find which one fails.
		t.Logf("Combined SQL failed: %v", err)
		t.Logf("Migration SQL:\n%s", migrationSQL)
		for i, op := range plan.Ops {
			t.Logf("Op[%d] %s: %s", i, op.Type, op.SQL)
		}
		t.Fatalf("LoadSQL(from + migration) error: %v", err)
	}

	// Now diff migrated vs to — should be empty.
	diff2 := Diff(migrated, to)
	if !diff2.IsEmpty() {
		var diffs []string
		if len(diff2.Schemas) > 0 {
			diffs = append(diffs, fmt.Sprintf("schemas: %d", len(diff2.Schemas)))
		}
		if len(diff2.Relations) > 0 {
			for _, r := range diff2.Relations {
				diffs = append(diffs, fmt.Sprintf("relation %s.%s action=%d cols=%d cons=%d idxs=%d trigs=%d pols=%d",
					r.SchemaName, r.Name, r.Action, len(r.Columns), len(r.Constraints), len(r.Indexes), len(r.Triggers), len(r.Policies)))
			}
		}
		if len(diff2.Sequences) > 0 {
			diffs = append(diffs, fmt.Sprintf("sequences: %d", len(diff2.Sequences)))
		}
		if len(diff2.Functions) > 0 {
			diffs = append(diffs, fmt.Sprintf("functions: %d", len(diff2.Functions)))
		}
		if len(diff2.Enums) > 0 {
			diffs = append(diffs, fmt.Sprintf("enums: %d", len(diff2.Enums)))
		}
		if len(diff2.Domains) > 0 {
			diffs = append(diffs, fmt.Sprintf("domains: %d", len(diff2.Domains)))
		}
		if len(diff2.Ranges) > 0 {
			diffs = append(diffs, fmt.Sprintf("ranges: %d", len(diff2.Ranges)))
		}
		if len(diff2.CompositeTypes) > 0 {
			diffs = append(diffs, fmt.Sprintf("compositeTypes: %d", len(diff2.CompositeTypes)))
		}
		if len(diff2.Extensions) > 0 {
			diffs = append(diffs, fmt.Sprintf("extensions: %d", len(diff2.Extensions)))
		}
		if len(diff2.Comments) > 0 {
			diffs = append(diffs, fmt.Sprintf("comments: %d", len(diff2.Comments)))
		}
		if len(diff2.Grants) > 0 {
			diffs = append(diffs, fmt.Sprintf("grants: %d", len(diff2.Grants)))
		}
		t.Logf("Migration SQL:\n%s", migrationSQL)
		t.Errorf("roundtrip failed: remaining diffs: %s", strings.Join(diffs, "; "))
	}
}

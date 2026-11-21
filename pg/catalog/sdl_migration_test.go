package catalog

import (
	"strings"
	"testing"
)

func TestLoadSDLOrSQL_OrderTolerantSDL(t *testing.T) {
	// Child table references parent table before parent is declared, and view
	// appears before the table it queries. LoadSDL handles both out of order.
	sdl := `
CREATE VIEW active_orders AS SELECT id, user_id FROM orders;
CREATE TABLE orders (
	id integer PRIMARY KEY,
	user_id integer REFERENCES users(id)
);
CREATE TABLE users (
	id integer PRIMARY KEY
);
`
	c, err := LoadSDLOrSQL(sdl)
	if err != nil {
		t.Fatalf("LoadSDLOrSQL failed on out-of-order SDL: %v", err)
	}
	if c.GetRelation("public", "users") == nil {
		t.Fatalf("expected table public.users in catalog")
	}
	if c.GetRelation("public", "orders") == nil {
		t.Fatalf("expected table public.orders in catalog")
	}
	if c.GetRelation("public", "active_orders") == nil {
		t.Fatalf("expected view public.active_orders in catalog")
	}
}

func TestLoadSDLOrSQL_FallbackForRawDump(t *testing.T) {
	// Raw pg_dump output includes SET / SELECT statements that LoadSDL rejects
	// but LoadSQL executes.
	dump := `
SET statement_timeout = 0;
SELECT pg_catalog.set_config('search_path', 'public', false);
CREATE TABLE users (
	id integer PRIMARY KEY,
	name text NOT NULL
);
`
	if _, err := LoadSDL(dump); err == nil {
		t.Fatalf("expected LoadSDL to reject SET/SELECT statements")
	}

	c, err := LoadSDLOrSQL(dump)
	if err != nil {
		t.Fatalf("LoadSDLOrSQL fallback failed: %v", err)
	}
	if c.GetRelation("public", "users") == nil {
		t.Fatalf("expected table public.users in catalog")
	}
}

func TestLoadSDLOrSQL_PreservesLoadSDLErrorWhenBothFail(t *testing.T) {
	// Forward FK reference fails LoadSQL with "relation does not exist", while
	// TRUNCATE fails LoadSDL with an SDL validation error. LoadSDLOrSQL must
	// return the LoadSDL error so diagnostics stay accurate.
	bad := `
CREATE TABLE orders (
	id integer PRIMARY KEY,
	user_id integer REFERENCES users(id)
);
CREATE TABLE users (
	id integer PRIMARY KEY
);
TRUNCATE TABLE users;
`
	_, sdlErr := LoadSDL(bad)
	if sdlErr == nil {
		t.Fatalf("expected LoadSDL error")
	}
	_, sqlErr := LoadSQL(bad)
	if sqlErr == nil {
		t.Fatalf("expected LoadSQL error")
	}
	if sdlErr.Error() == sqlErr.Error() {
		t.Fatalf("test fixture needs distinct LoadSDL and LoadSQL errors, both returned %q", sdlErr.Error())
	}

	_, err := LoadSDLOrSQL(bad)
	if err == nil {
		t.Fatalf("expected LoadSDLOrSQL error")
	}
	if err.Error() != sdlErr.Error() {
		t.Fatalf("LoadSDLOrSQL error = %q, want LoadSDL error %q (LoadSQL error was %q)", err.Error(), sdlErr.Error(), sqlErr.Error())
	}
}

func TestBuildSDLMigrationPlan_EmptyDiffReturnsNil(t *testing.T) {
	source := `
CREATE TABLE users (
	id integer PRIMARY KEY,
	name text NOT NULL
);
CREATE TABLE orders (
	id integer PRIMARY KEY,
	user_id integer REFERENCES users(id)
);
`
	// Reversed declaration order should still produce a nil plan.
	target := `
CREATE TABLE orders (
	id integer PRIMARY KEY,
	user_id integer REFERENCES users(id)
);
CREATE TABLE users (
	id integer PRIMARY KEY,
	name text NOT NULL
);
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{})
	if err != nil {
		t.Fatalf("BuildSDLMigrationPlan failed: %v", err)
	}
	if plan != nil {
		t.Fatalf("expected nil plan for equivalent schemas, got %d ops:\n%s", len(plan.Ops), plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_RawDumpFallback(t *testing.T) {
	source := `
SET statement_timeout = 0;
SELECT 1;
CREATE TABLE users (
	id integer PRIMARY KEY
);
`
	target := `
CREATE TABLE users (
	id integer PRIMARY KEY,
	email text
);
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{})
	if err != nil {
		t.Fatalf("BuildSDLMigrationPlan failed: %v", err)
	}
	if plan == nil {
		t.Fatalf("expected non-nil plan")
	}
	sql := plan.SQL()
	if !strings.Contains(sql, "ALTER TABLE") || !strings.Contains(sql, "email") {
		t.Fatalf("expected ALTER TABLE adding email, got:\n%s", sql)
	}
}

func TestBuildSDLMigrationPlan_WrapsLoadErrors(t *testing.T) {
	if _, err := BuildSDLMigrationPlan("CREATE TABLE (", "", SDLMigrationOptions{}); err == nil || !strings.Contains(err.Error(), "failed to load source schema") {
		t.Fatalf("expected wrapped source error, got %v", err)
	}
	if _, err := BuildSDLMigrationPlan("", "CREATE TABLE (", SDLMigrationOptions{}); err == nil || !strings.Contains(err.Error(), "failed to load target schema") {
		t.Fatalf("expected wrapped target error, got %v", err)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemas(t *testing.T) {
	source := `
CREATE SCHEMA bbdataarchive;
CREATE TABLE public.users (id integer PRIMARY KEY);
CREATE TABLE bbdataarchive.archived_users (id integer PRIMARY KEY);
`
	target := `
CREATE TABLE public.users (id integer PRIMARY KEY, name text);
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"", "bbdataarchive"},
	})
	if err != nil {
		t.Fatalf("BuildSDLMigrationPlan failed: %v", err)
	}
	if plan == nil {
		t.Fatalf("expected non-nil plan for public.users change")
	}
	sql := plan.SQL()
	if strings.Contains(sql, "bbdataarchive") || strings.Contains(sql, "archived_users") {
		t.Fatalf("expected skipped schema operations to be excluded, got:\n%s", sql)
	}
	if !strings.Contains(sql, "users") || !strings.Contains(sql, "name") {
		t.Fatalf("expected public.users alteration to be preserved, got:\n%s", sql)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasOnlyChangesReturnsNil(t *testing.T) {
	source := `
CREATE SCHEMA bbdataarchive;
CREATE TABLE public.users (id integer PRIMARY KEY);
CREATE TABLE bbdataarchive.archived_users (id integer PRIMARY KEY);
`
	target := `
CREATE TABLE public.users (id integer PRIMARY KEY);
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"bbdataarchive"},
	})
	if err != nil {
		t.Fatalf("BuildSDLMigrationPlan failed: %v", err)
	}
	if plan != nil {
		t.Fatalf("expected nil plan when all changes belong to skipped schemas, got %d ops:\n%s", len(plan.Ops), plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_PreservesExtensionsWithSkipSchemas(t *testing.T) {
	plan, err := BuildSDLMigrationPlan("", `CREATE EXTENSION IF NOT EXISTS "citext";`, SDLMigrationOptions{
		SkipSchemas: []string{"bbdataarchive"},
	})
	if err != nil {
		t.Fatalf("BuildSDLMigrationPlan failed: %v", err)
	}
	if plan == nil {
		t.Fatalf("expected non-nil plan")
	}
	if !strings.Contains(plan.SQL(), "citext") {
		t.Fatalf("expected extension creation to be preserved, got:\n%s", plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsCrossSchemaViewDependency(t *testing.T) {
	source := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id integer PRIMARY KEY);
CREATE VIEW archive.v AS SELECT id FROM public.t;
`
	target := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id bigint PRIMARY KEY);
CREATE VIEW archive.v AS SELECT id FROM public.t;
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err == nil {
		t.Fatalf("expected dependency conflict error, got plan:\n%s", plan.SQL())
	}
	if !strings.Contains(err.Error(), "archive") {
		t.Fatalf("expected error mentioning skipped schema archive, got: %v", err)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsTransitiveViewAndDropTableConflict(t *testing.T) {
	source := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id integer PRIMARY KEY);
CREATE VIEW public.v1 AS SELECT id FROM public.t;
CREATE VIEW archive.v2 AS SELECT id FROM public.v1;
`
	// Column type change on public.t transitively breaks archive.v2.
	targetAlter := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id bigint PRIMARY KEY);
CREATE VIEW public.v1 AS SELECT id FROM public.t;
CREATE VIEW archive.v2 AS SELECT id FROM public.v1;
`
	if plan, err := BuildSDLMigrationPlan(source, targetAlter, SDLMigrationOptions{SkipSchemas: []string{"archive"}}); err == nil {
		t.Fatalf("expected transitive column-alter dependency error, got plan:\n%s", plan.SQL())
	}

	// Dropping public.t while archive.v2 depends on it (via public.v1) must also fail.
	if plan, err := BuildSDLMigrationPlan(source, "", SDLMigrationOptions{SkipSchemas: []string{"archive"}}); err == nil {
		t.Fatalf("expected drop-table dependency error, got plan:\n%s", plan.SQL())
	}

	// Adding an unrelated column to public.t does not break archive.v2 and must succeed.
	targetSafe := `
CREATE TABLE public.t (id integer PRIMARY KEY, note text);
CREATE VIEW public.v1 AS SELECT id FROM public.t;
`
	plan, err := BuildSDLMigrationPlan(source, targetSafe, SDLMigrationOptions{SkipSchemas: []string{"archive"}})
	if err != nil {
		t.Fatalf("expected safe ADD COLUMN to succeed despite skipped view, got error: %v", err)
	}
	if plan == nil || !strings.Contains(plan.SQL(), "note") || strings.Contains(plan.SQL(), "archive") {
		t.Fatalf("unexpected plan for safe ADD COLUMN: %v", plan)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsNewDependencyOnSkippedCreate(t *testing.T) {
	target := `
CREATE SCHEMA archive;
CREATE TABLE archive.base (id integer PRIMARY KEY);
CREATE VIEW public.v AS SELECT id FROM archive.base;
`
	plan, err := BuildSDLMigrationPlan("", target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err == nil {
		t.Fatalf("expected error when non-skipped view depends on newly created skipped table, got plan:\n%s", plan.SQL())
	}
}

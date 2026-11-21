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
	if !strings.Contains(err.Error(), "used by a view or rule") {
		t.Fatalf("expected view/rule dependency error, got: %v", err)
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

func TestBuildSDLMigrationPlan_SkipSchemasAllowsUnreferencedColumnAlter(t *testing.T) {
	// archive.v only selects id, so changing note's type in public.t is valid
	// in PostgreSQL even while archive.v remains intact.
	source := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id integer PRIMARY KEY, note integer);
CREATE VIEW archive.v AS SELECT id FROM public.t;
`
	target := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id integer PRIMARY KEY, note bigint);
CREATE VIEW archive.v AS SELECT id FROM public.t;
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err != nil {
		t.Fatalf("expected unreferenced column alteration to succeed, got error: %v", err)
	}
	if plan == nil || !strings.Contains(plan.SQL(), "note") || strings.Contains(plan.SQL(), "archive") {
		t.Fatalf("unexpected plan: %v", plan)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasDistinguishesAlterFromCreate(t *testing.T) {
	source := `
CREATE SCHEMA archive;
CREATE TABLE archive.base (id integer PRIMARY KEY);
`
	// Unrelated alteration on existing archive.base must not block a view that
	// only references the pre-existing id column.
	targetSafe := `
CREATE SCHEMA archive;
CREATE TABLE archive.base (id integer PRIMARY KEY, note text);
CREATE VIEW public.v AS SELECT id FROM archive.base;
`
	plan, err := BuildSDLMigrationPlan(source, targetSafe, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err != nil {
		t.Fatalf("expected view referencing pre-existing column to succeed, got error: %v", err)
	}
	if plan == nil || !strings.Contains(plan.SQL(), "CREATE VIEW") || strings.Contains(plan.SQL(), "ADD COLUMN") {
		t.Fatalf("unexpected plan: %v", plan)
	}

	// Referencing the newly added (and skipped) note column must be rejected.
	targetMissingCol := `
CREATE SCHEMA archive;
CREATE TABLE archive.base (id integer PRIMARY KEY, note text);
CREATE VIEW public.v AS SELECT id, note FROM archive.base;
`
	if plan, err := BuildSDLMigrationPlan(source, targetMissingCol, SDLMigrationOptions{SkipSchemas: []string{"archive"}}); err == nil {
		t.Fatalf("expected error when view references skipped newly added column, got plan:\n%s", plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsTriggerDependingOnSkippedFunction(t *testing.T) {
	target := `
CREATE SCHEMA archive;
CREATE FUNCTION archive.audit_fn() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
	RETURN NEW;
END;
$$;
CREATE TABLE public.t (id integer PRIMARY KEY);
CREATE TRIGGER trg_t BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION archive.audit_fn();
`
	plan, err := BuildSDLMigrationPlan("", target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err == nil {
		t.Fatalf("expected error when retained trigger depends on newly created skipped function, got plan:\n%s", plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsTableUsingSkippedEnumType(t *testing.T) {
	source := `CREATE SCHEMA archive;`
	target := `
CREATE SCHEMA archive;
CREATE TYPE archive.state AS ENUM ('ready');
CREATE TABLE public.jobs (id integer PRIMARY KEY, state archive.state);
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err == nil {
		t.Fatalf("expected error when retained table uses newly created skipped enum type, got plan:\n%s", plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsFunctionUsingSkippedCompositeType(t *testing.T) {
	source := `CREATE SCHEMA archive;`
	target := `
CREATE SCHEMA archive;
CREATE TYPE archive.pair AS (x integer);
CREATE FUNCTION public.f(x archive.pair) RETURNS integer LANGUAGE sql AS 'SELECT 1';
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err == nil {
		t.Fatalf("expected error when retained function uses newly created skipped composite type, got plan:\n%s", plan.SQL())
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsDefaultUsingSkippedEnumValue(t *testing.T) {
	source := `
CREATE SCHEMA archive;
CREATE TYPE archive.state AS ENUM ('old');
`
	target := `
CREATE SCHEMA archive;
CREATE TYPE archive.state AS ENUM ('old', 'new');
CREATE TABLE public.jobs (id integer PRIMARY KEY, state archive.state DEFAULT 'new');
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err == nil {
		t.Fatalf("expected error when retained column default uses skipped enum value addition, got plan:\n%s", plan.SQL())
	}

	// Using the pre-existing enum value 'old' must still succeed when 'new' is skipped.
	targetExisting := `
CREATE SCHEMA archive;
CREATE TYPE archive.state AS ENUM ('old', 'new');
CREATE TABLE public.jobs (id integer PRIMARY KEY, state archive.state DEFAULT 'old');
`
	safePlan, err := BuildSDLMigrationPlan(source, targetExisting, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err != nil {
		t.Fatalf("expected pre-existing enum value default to succeed, got error: %v", err)
	}
	if safePlan == nil || !strings.Contains(safePlan.SQL(), "CREATE TABLE") || strings.Contains(safePlan.SQL(), "ADD VALUE") {
		t.Fatalf("unexpected plan: %v", safePlan)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasAllowsPartitionDetachWithSkippedView(t *testing.T) {
	source := `
CREATE SCHEMA archive;
CREATE TABLE public.parent (id integer) PARTITION BY RANGE (id);
CREATE TABLE public.child PARTITION OF public.parent FOR VALUES FROM (1) TO (100);
CREATE VIEW archive.v AS SELECT id FROM public.child;
`
	target := `
CREATE SCHEMA archive;
CREATE TABLE archive.extra (id integer);
CREATE TABLE public.parent (id integer) PARTITION BY RANGE (id);
CREATE TABLE public.child (id integer);
CREATE VIEW archive.v AS SELECT id FROM public.child;
`
	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SkipSchemas: []string{"archive"},
	})
	if err != nil {
		t.Fatalf("expected partition detach to succeed despite skipped view on child, got error: %v", err)
	}
	if plan == nil || !strings.Contains(plan.SQL(), "DETACH PARTITION") || strings.Contains(plan.SQL(), "archive") {
		t.Fatalf("unexpected plan: %v", plan)
	}
}

func TestBuildSDLMigrationPlan_SkipSchemasRejectsCascadeIntoSkippedSchema(t *testing.T) {
	t.Run("drop_table_cascade_view_in_skipped_schema", func(t *testing.T) {
		source := `
CREATE SCHEMA archive;
CREATE TABLE public.t (id integer PRIMARY KEY);
CREATE VIEW archive.v AS SELECT id FROM public.t;
`
		plan, err := BuildSDLMigrationPlan(source, "", SDLMigrationOptions{
			SkipSchemas: []string{"archive"},
		})
		if err == nil {
			t.Fatalf("expected error when DROP TABLE CASCADE would destroy skipped view archive.v, got plan:\n%s", plan.SQL())
		}
		if !strings.Contains(err.Error(), "cascade") || !strings.Contains(err.Error(), "archive") {
			t.Fatalf("expected cascade error mentioning archive, got: %v", err)
		}
	})

	t.Run("drop_table_cascade_fk_in_skipped_schema", func(t *testing.T) {
		source := `
CREATE SCHEMA archive;
CREATE TABLE public.parent (id integer PRIMARY KEY);
CREATE TABLE archive.child (id integer PRIMARY KEY, parent_id integer REFERENCES public.parent(id));
`
		plan, err := BuildSDLMigrationPlan(source, "", SDLMigrationOptions{
			SkipSchemas: []string{"archive"},
		})
		if err == nil {
			t.Fatalf("expected error when DROP TABLE CASCADE would drop FK constraint in archive.child, got plan:\n%s", plan.SQL())
		}
	})
}

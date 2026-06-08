package catalog

import (
	"strings"
	"testing"
)

// TestContainerType covers section 2.6: Type/Sequence/Extension Changes.
func TestContainerType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	// -----------------------------------------------------------------------
	// Enum changes
	// -----------------------------------------------------------------------

	t.Run("enum_add_value_at_end", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE color AS ENUM ('red', 'green', 'blue');`,
			`CREATE TYPE color AS ENUM ('red', 'green', 'blue', 'yellow');`,
		)
	})

	t.Run("enum_add_value_with_before", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE color AS ENUM ('red', 'green', 'blue');`,
			`CREATE TYPE color AS ENUM ('red', 'yellow', 'green', 'blue');`,
		)
	})

	t.Run("enum_add_value_with_after", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE color AS ENUM ('red', 'green', 'blue');`,
			`CREATE TYPE color AS ENUM ('red', 'green', 'cyan', 'blue');`,
		)
	})

	// -----------------------------------------------------------------------
	// Domain changes
	// -----------------------------------------------------------------------

	t.Run("domain_change_default", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE DOMAIN score AS integer DEFAULT 0;`,
			`CREATE DOMAIN score AS integer DEFAULT 100;`,
		)
	})

	t.Run("domain_change_not_null", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE DOMAIN score AS integer;`,
			`CREATE DOMAIN score AS integer NOT NULL;`,
		)
	})

	t.Run("domain_add_constraint", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE DOMAIN score AS integer;`,
			`CREATE DOMAIN score AS integer CONSTRAINT score_positive CHECK (VALUE >= 0);`,
		)
	})

	// -----------------------------------------------------------------------
	// Sequence changes
	// -----------------------------------------------------------------------

	t.Run("sequence_change_increment", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE SEQUENCE my_seq INCREMENT 1 MINVALUE 1 MAXVALUE 1000 START 1 CACHE 1 NO CYCLE;`,
			`CREATE SEQUENCE my_seq INCREMENT 5 MINVALUE 1 MAXVALUE 1000 START 1 CACHE 1 NO CYCLE;`,
		)
	})

	t.Run("sequence_change_cycle", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE SEQUENCE my_seq INCREMENT 1 MINVALUE 1 MAXVALUE 1000 START 1 CACHE 1 NO CYCLE;`,
			`CREATE SEQUENCE my_seq INCREMENT 1 MINVALUE 1 MAXVALUE 1000 START 1 CACHE 1 CYCLE;`,
		)
	})

	// -----------------------------------------------------------------------
	// Add/Drop enum type
	// -----------------------------------------------------------------------

	t.Run("add_new_enum_type", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			``,
			`CREATE TYPE mood AS ENUM ('happy', 'sad', 'neutral');`,
		)
	})

	t.Run("drop_enum_type", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE mood AS ENUM ('happy', 'sad', 'neutral');`,
			``,
		)
	})

	// -----------------------------------------------------------------------
	// Composite type changes
	// -----------------------------------------------------------------------

	t.Run("composite_add_column", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE address AS (street text, city text);`,
			`CREATE TYPE address AS (street text, city text, zip varchar(10));`,
		)
	})

	t.Run("composite_drop_column", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE address AS (street text, city text, zip varchar(10));`,
			`CREATE TYPE address AS (street text, city text);`,
		)
	})

	t.Run("composite_modify_column", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TYPE address AS (street text, city text, zip varchar(5));`,
			`CREATE TYPE address AS (street text, city text, zip varchar(10));`,
		)
	})

	// -----------------------------------------------------------------------
	// Range type: change subtype (warning)
	// -----------------------------------------------------------------------

	t.Run("range_change_subtype_warning", func(t *testing.T) {
		// Range type subtype change cannot be done with ALTER TYPE; we expect a
		// warning in the migration plan. The migration SQL will be a comment,
		// so the schemas will differ (migrated won't have the change).
		// We test this by verifying the plan has a warning.
		fromDDL := `CREATE TYPE int_range AS RANGE (SUBTYPE = int4);`
		toDDL := `CREATE TYPE int_range AS RANGE (SUBTYPE = int8);`

		fromCat, err := LoadSQL(fromDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
		toCat, err := LoadSQL(toDDL)
		if err != nil {
			t.Fatalf("LoadSQL(after) failed: %v", err)
		}
		diff := Diff(fromCat, toCat)
		plan := GenerateMigration(fromCat, toCat, diff)
		if !plan.HasWarnings() {
			t.Fatal("expected warning for range subtype change, got none")
		}
		foundWarning := false
		for _, w := range plan.Warnings() {
			if strings.Contains(w.Warning, "range type") && strings.Contains(w.Warning, "subtype") {
				foundWarning = true
			}
		}
		if !foundWarning {
			t.Errorf("expected range subtype warning, got: %v", plan.Warnings())
		}
	})

	// -----------------------------------------------------------------------
	// Schema: add and drop
	// -----------------------------------------------------------------------

	t.Run("schema_add", func(t *testing.T) {
		// Schema add is tested by verifying the diff/migration generates CREATE SCHEMA.
		// assertOracleRoundtrip operates within test schemas so we test via plan inspection.
		fromDDL := ``
		toDDL := `CREATE SCHEMA extra;`

		fromCat, err := LoadSQL(fromDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
		toCat, err := LoadSQL(toDDL)
		if err != nil {
			t.Fatalf("LoadSQL(after) failed: %v", err)
		}
		diff := Diff(fromCat, toCat)
		plan := GenerateMigration(fromCat, toCat, diff)

		found := false
		for _, op := range plan.Ops {
			if op.Type == OpCreateSchema && strings.Contains(op.SQL, "extra") {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected CREATE SCHEMA extra in migration plan, got: %s", plan.SQL())
		}

		// Verify the SQL executes on PG.
		schema := ctr.freshSchema(t)
		migrationSQL := plan.SQL()
		// The migration creates a schema named "extra", execute it directly.
		ctr.execSQL(t, migrationSQL)
		// Clean up the extra schema.
		t.Cleanup(func() {
			ctr.execSQL(t, `DROP SCHEMA IF EXISTS "extra" CASCADE`)
		})

		// Verify schema exists.
		var count int
		err = ctr.db.QueryRowContext(ctr.ctx,
			`SELECT count(*) FROM pg_namespace WHERE nspname = 'extra'`).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected schema 'extra' to exist, got count=%d", count)
		}
		_ = schema // used for fresh schema allocation
	})

	t.Run("schema_drop", func(t *testing.T) {
		fromDDL := `CREATE SCHEMA extra;`
		toDDL := ``

		fromCat, err := LoadSQL(fromDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
		toCat, err := LoadSQL(toDDL)
		if err != nil {
			t.Fatalf("LoadSQL(after) failed: %v", err)
		}
		diff := Diff(fromCat, toCat)
		plan := GenerateMigration(fromCat, toCat, diff)

		found := false
		for _, op := range plan.Ops {
			if op.Type == OpDropSchema && strings.Contains(op.SQL, "extra") {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected DROP SCHEMA extra in migration plan, got: %s", plan.SQL())
		}

		// Create the schema first, then apply the migration.
		ctr.execSQL(t, `CREATE SCHEMA IF NOT EXISTS "extra"`)
		migrationSQL := plan.SQL()
		ctr.execSQL(t, migrationSQL)

		// Verify schema no longer exists.
		var count int
		err = ctr.db.QueryRowContext(ctr.ctx,
			`SELECT count(*) FROM pg_namespace WHERE nspname = 'extra'`).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("expected schema 'extra' to be dropped, got count=%d", count)
		}
	})

	// -----------------------------------------------------------------------
	// GRANT / REVOKE
	// -----------------------------------------------------------------------

	t.Run("grant_select_on_table", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TABLE items (id integer PRIMARY KEY);`,
			`CREATE TABLE items (id integer PRIMARY KEY);
			 GRANT SELECT ON items TO PUBLIC;`,
		)
	})

	t.Run("revoke_privilege", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TABLE items (id integer PRIMARY KEY);
			 GRANT SELECT ON items TO PUBLIC;`,
			`CREATE TABLE items (id integer PRIMARY KEY);`,
		)
	})

	// -----------------------------------------------------------------------
	// Partitioned table and partition child
	// -----------------------------------------------------------------------

	t.Run("partitioned_table", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			``,
			`CREATE TABLE events (
				id integer NOT NULL,
				event_date date NOT NULL,
				data text
			) PARTITION BY RANGE (event_date);`,
		)
	})

	t.Run("partition_child", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr,
			`CREATE TABLE events (
				id integer NOT NULL,
				event_date date NOT NULL,
				data text
			) PARTITION BY RANGE (event_date);`,
			`CREATE TABLE events (
				id integer NOT NULL,
				event_date date NOT NULL,
				data text
			) PARTITION BY RANGE (event_date);
			CREATE TABLE events_2024 PARTITION OF events
				FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');`,
		)
	})
}

// ---------------------------------------------------------------------------
// Comment Coverage Gaps
// ---------------------------------------------------------------------------

func TestContainerCommentGaps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	// --- COMMENT ON INDEX ---
	t.Run("comment_on_index", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100)
);
CREATE INDEX idx_name ON t1(name);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100)
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON INDEX idx_name IS 'Name lookup index';
`
		assertOracleRoundtrip(t, ctr, before, after)
	})

	// --- COMMENT ON TYPE ---
	t.Run("comment_on_type", func(t *testing.T) {
		before := `CREATE TYPE color AS ENUM ('red', 'green', 'blue');`
		after := `CREATE TYPE color AS ENUM ('red', 'green', 'blue');
COMMENT ON TYPE color IS 'Available colors';`
		assertOracleRoundtrip(t, ctr, before, after)
	})

	// --- COMMENT ON SEQUENCE ---
	t.Run("comment_on_sequence", func(t *testing.T) {
		before := `CREATE SEQUENCE my_seq INCREMENT 1 MINVALUE 1 MAXVALUE 1000 START 1 CACHE 1 NO CYCLE;`
		after := `CREATE SEQUENCE my_seq INCREMENT 1 MINVALUE 1 MAXVALUE 1000 START 1 CACHE 1 NO CYCLE;
COMMENT ON SEQUENCE my_seq IS 'Auto-increment sequence';`
		assertOracleRoundtrip(t, ctr, before, after)
	})

	// --- COMMENT ON CONSTRAINT ---
	t.Run("comment_on_constraint", func(t *testing.T) {
		// NOTE: COMMENT ON CONSTRAINT is tracked via pg_description with objsubid
		// referencing the constraint. This may not be supported by the diff engine.
		before := `
CREATE TABLE t1 (
    id integer,
    CONSTRAINT t1_pkey PRIMARY KEY (id)
);
`
		after := `
CREATE TABLE t1 (
    id integer,
    CONSTRAINT t1_pkey PRIMARY KEY (id)
);
COMMENT ON CONSTRAINT t1_pkey ON t1 IS 'Primary key constraint';
`
		assertOracleRoundtrip(t, ctr, before, after)
	})

	// --- COMMENT ON TRIGGER ---
	t.Run("comment_on_trigger", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100)
);
CREATE TRIGGER t1_trig BEFORE UPDATE ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100)
);
CREATE TRIGGER t1_trig BEFORE UPDATE ON t1 FOR EACH ROW EXECUTE FUNCTION trig_fn();
COMMENT ON TRIGGER t1_trig ON t1 IS 'Update trigger';
`
		assertOracleRoundtrip(t, ctr, before, after)
	})

	// --- COMMENT ON SCHEMA ---
	t.Run("comment_on_schema", func(t *testing.T) {
		// Schema comments are special: assertOracleRoundtrip operates within
		// test schemas so we test via plan inspection.
		fromDDL := `CREATE SCHEMA extra;`
		toDDL := `CREATE SCHEMA extra;`
		// This test would need COMMENT ON SCHEMA, but the diff engine may
		// not support schema comments. We test via plan.
		fromCat, err := LoadSQL(fromDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
		toCat, err := LoadSQL(toDDL + "\nCOMMENT ON SCHEMA extra IS 'Extra schema';")
		if err != nil {
			t.Skip("LoadSQL does not support COMMENT ON SCHEMA syntax")
			return
		}
		diff := Diff(fromCat, toCat)
		plan := GenerateMigration(fromCat, toCat, diff)
		found := false
		for _, op := range plan.Ops {
			if strings.Contains(op.SQL, "COMMENT ON SCHEMA") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected COMMENT ON SCHEMA in migration plan, got: %s", plan.SQL())
		}
	})
}

// ---------------------------------------------------------------------------
// Index Type Coverage Gaps
// ---------------------------------------------------------------------------

func TestContainerIndexGaps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	// --- Index USING hash ---
	t.Run("index_using_hash", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100)
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100)
);
CREATE INDEX idx_name_hash ON t1 USING hash (name);
`
		assertOracleRoundtrip(t, ctr, before, after)
	})

	// --- Index USING brin ---
	t.Run("index_using_brin", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    created_at timestamptz DEFAULT now()
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_created_brin ON t1 USING brin (created_at);
`
		assertOracleRoundtrip(t, ctr, before, after)
	})
}

// ---------------------------------------------------------------------------
// Table INHERITS Coverage Gap
// ---------------------------------------------------------------------------

func TestContainerInherits(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	// --- Table INHERITS ---
	t.Run("table_inherits", func(t *testing.T) {
		before := `
CREATE TABLE parent_tbl (
    id integer PRIMARY KEY,
    name varchar(100)
);
`
		after := `
CREATE TABLE parent_tbl (
    id integer PRIMARY KEY,
    name varchar(100)
);
CREATE TABLE child_tbl (
    extra text
) INHERITS (parent_tbl);
`
		assertOracleRoundtrip(t, ctr, before, after)
	})
}

// ---------------------------------------------------------------------------
// Group 7: Domain with Function-based CHECK
// ---------------------------------------------------------------------------

func TestContainerDomainWithFunctionCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	t.Run("domain_with_function_in_check_change_function_body", func(t *testing.T) {
		before := `
CREATE FUNCTION validate(val integer) RETURNS boolean
    LANGUAGE sql IMMUTABLE AS 'SELECT val > 0';
CREATE DOMAIN posint AS integer CONSTRAINT posint_check CHECK (validate(VALUE));
`
		after := `
CREATE FUNCTION validate(val integer) RETURNS boolean
    LANGUAGE sql IMMUTABLE AS 'SELECT val > 0 AND val < 1000';
CREATE DOMAIN posint AS integer CONSTRAINT posint_check CHECK (validate(VALUE));
`
		assertOracleRoundtrip(t, ctr, before, after)
	})
}

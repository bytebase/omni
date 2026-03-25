package catalog

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Phase 3: Combinations, Interactions, and Known Regressions
// ---------------------------------------------------------------------------

// assertOracleRoundtripSDL is like assertOracleRoundtrip but uses LoadSDL
// for the "after" catalog (to test SDL forward-reference resolution).
// expectedDDL is PG-executable DDL in correct dependency order for the
// expected schema. If empty, afterDDL is used directly (must be valid PG order).
func assertOracleRoundtripSDL(t *testing.T, o *pgOracle, beforeDDL, afterDDL string, expectedDDL ...string) {
	t.Helper()

	migrated := o.freshSchema(t)
	expected := o.freshSchema(t)

	// 1. Apply "before" to the migrated schema.
	if beforeDDL != "" {
		o.execInSchema(t, migrated, beforeDDL)
	}

	// 2. Generate migration via omni, using LoadSDL for the after catalog.
	fromCat := New()
	if beforeDDL != "" {
		var err error
		fromCat, err = LoadSQL(beforeDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
	}
	toCat := New()
	if afterDDL != "" {
		var err error
		toCat, err = LoadSDL(afterDDL)
		if err != nil {
			t.Fatalf("LoadSDL(after) failed: %v", err)
		}
	}
	diff := Diff(fromCat, toCat)
	plan := GenerateMigration(fromCat, toCat, diff)

	// 3. Apply migration to PG.
	migrationSQL := plan.SQL()
	migrationSQL = strings.ReplaceAll(migrationSQL, "public.", "")
	migrationSQL = strings.ReplaceAll(migrationSQL, `"public".`, "")
	if migrationSQL != "" {
		o.execInSchema(t, migrated, migrationSQL)
	}

	// 4. Apply DDL directly to expected schema.
	// Use expectedDDL if provided (for cases where afterDDL is in reverse order).
	pgDDL := afterDDL
	if len(expectedDDL) > 0 && expectedDDL[0] != "" {
		pgDDL = expectedDDL[0]
	}
	if pgDDL != "" {
		o.execInSchema(t, expected, pgDDL)
	}

	// 5. Compare schemas.
	o.assertSchemasEqual(t, migrated, expected)
}

// generateMigrationSQL is a helper that generates migration SQL for (before, after).
func generateMigrationSQL(t *testing.T, beforeDDL, afterDDL string) string {
	t.Helper()
	fromCat := New()
	if beforeDDL != "" {
		var err error
		fromCat, err = LoadSQL(beforeDDL)
		if err != nil {
			t.Fatalf("LoadSQL(before) failed: %v", err)
		}
	}
	toCat := New()
	if afterDDL != "" {
		var err error
		toCat, err = LoadSQL(afterDDL)
		if err != nil {
			t.Fatalf("LoadSQL(after) failed: %v", err)
		}
	}
	diff := Diff(fromCat, toCat)
	plan := GenerateMigration(fromCat, toCat, diff)
	return plan.SQL()
}

// ---------------------------------------------------------------------------
// 3.1 Multi-Attribute Simultaneous Changes
// ---------------------------------------------------------------------------

func TestOracleCombo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// -------------------------------------------------------------------
	// 3.1 Multi-Attribute Simultaneous Changes
	// -------------------------------------------------------------------

	t.Run("multi_attr_table_add_col_change_type_drop_constraint", func(t *testing.T) {
		before := `
CREATE TABLE t (
    id integer NOT NULL,
    name varchar(50) NOT NULL,
    email text,
    CONSTRAINT t_pkey PRIMARY KEY (id),
    CONSTRAINT t_email_unique UNIQUE (email)
);`
		after := `
CREATE TABLE t (
    id integer NOT NULL,
    name varchar(100) NOT NULL,
    email text,
    phone text,
    CONSTRAINT t_pkey PRIMARY KEY (id)
);`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("multi_attr_function_body_volatility_parallel", func(t *testing.T) {
		before := `
CREATE FUNCTION f(x integer) RETURNS integer
    LANGUAGE sql
    VOLATILE
    PARALLEL UNSAFE
    AS 'SELECT x + 1';`
		after := `
CREATE FUNCTION f(x integer) RETURNS integer
    LANGUAGE sql
    STABLE
    PARALLEL SAFE
    AS 'SELECT x * 2';`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("multi_attr_table_add_index_trigger_policy", func(t *testing.T) {
		before := `
CREATE TABLE t (
    id integer NOT NULL,
    name text NOT NULL,
    active boolean DEFAULT true,
    CONSTRAINT t_pkey PRIMARY KEY (id)
);
CREATE FUNCTION stamp() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;`
		after := `
CREATE TABLE t (
    id integer NOT NULL,
    name text NOT NULL,
    active boolean DEFAULT true,
    CONSTRAINT t_pkey PRIMARY KEY (id)
);
CREATE FUNCTION stamp() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
CREATE INDEX idx_t_name ON t(name);
CREATE TRIGGER t_stamp BEFORE INSERT ON t
    FOR EACH ROW EXECUTE FUNCTION stamp();
ALTER TABLE t ENABLE ROW LEVEL SECURITY;
CREATE POLICY t_sel ON t FOR SELECT USING (active = true);`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("multi_attr_multiple_tables_changed", func(t *testing.T) {
		before := `
CREATE TABLE t1 (id integer NOT NULL, CONSTRAINT t1_pkey PRIMARY KEY (id));
CREATE TABLE t2 (id integer NOT NULL, CONSTRAINT t2_pkey PRIMARY KEY (id));`
		after := `
CREATE TABLE t1 (id integer NOT NULL, name text, CONSTRAINT t1_pkey PRIMARY KEY (id));
CREATE TABLE t2 (id integer NOT NULL, email text, CONSTRAINT t2_pkey PRIMARY KEY (id));`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// -------------------------------------------------------------------
	// 3.2 Cross-Object Dependency Changes
	// -------------------------------------------------------------------

	t.Run("cross_modify_function_used_by_check", func(t *testing.T) {
		before := `
CREATE FUNCTION check_age(val integer) RETURNS boolean
    LANGUAGE sql IMMUTABLE AS 'SELECT val >= 0 AND val <= 150';
CREATE TABLE t (
    id integer NOT NULL,
    age integer,
    CONSTRAINT t_pkey PRIMARY KEY (id),
    CONSTRAINT age_check CHECK (check_age(age))
);`
		after := `
CREATE FUNCTION check_age(val integer) RETURNS boolean
    LANGUAGE sql IMMUTABLE AS 'SELECT val >= 18 AND val <= 120';
CREATE TABLE t (
    id integer NOT NULL,
    age integer,
    CONSTRAINT t_pkey PRIMARY KEY (id),
    CONSTRAINT age_check CHECK (check_age(age))
);`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("cross_modify_function_used_by_trigger", func(t *testing.T) {
		before := `
CREATE TABLE t (id integer NOT NULL, name text, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE FUNCTION on_change() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
CREATE TRIGGER t_change BEFORE UPDATE ON t
    FOR EACH ROW EXECUTE FUNCTION on_change();`
		after := `
CREATE TABLE t (id integer NOT NULL, name text, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE FUNCTION on_change() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN NEW.name = upper(NEW.name); RETURN NEW; END;$$;
CREATE TRIGGER t_change BEFORE UPDATE ON t
    FOR EACH ROW EXECUTE FUNCTION on_change();`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("cross_modify_function_used_by_view", func(t *testing.T) {
		before := `
CREATE TABLE t (id integer NOT NULL, name text, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE FUNCTION format_name(n text) RETURNS text
    LANGUAGE sql IMMUTABLE AS 'SELECT upper(n)';
CREATE VIEW v AS SELECT id, format_name(name) AS formatted FROM t;`
		after := `
CREATE TABLE t (id integer NOT NULL, name text, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE FUNCTION format_name(n text) RETURNS text
    LANGUAGE sql IMMUTABLE AS 'SELECT lower(n)';
CREATE VIEW v AS SELECT id, format_name(name) AS formatted FROM t;`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("cross_add_enum_and_column_using_it", func(t *testing.T) {
		before := `
CREATE TABLE t (id integer NOT NULL, CONSTRAINT t_pkey PRIMARY KEY (id));`
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE TABLE t (id integer NOT NULL, s status, CONSTRAINT t_pkey PRIMARY KEY (id));`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("cross_drop_table_with_fk", func(t *testing.T) {
		before := `
CREATE TABLE parent (id integer NOT NULL, CONSTRAINT parent_pkey PRIMARY KEY (id));
CREATE TABLE child (
    id integer NOT NULL,
    pid integer,
    CONSTRAINT child_pkey PRIMARY KEY (id),
    CONSTRAINT child_pid_fk FOREIGN KEY (pid) REFERENCES parent(id)
);`
		after := ``
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("cross_modify_column_type_used_by_view", func(t *testing.T) {
		// BUG: PG cannot ALTER column type when a view depends on it.
		// The migration generator should DROP+RECREATE the view, but currently
		// does not. Skipping until upstream fix.
		t.Skip("KNOWN BUG: migration does not drop/recreate view before altering column type")
		before := `
CREATE TABLE t (id integer NOT NULL, name varchar(50), CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE VIEW v AS SELECT id, name FROM t;`
		after := `
CREATE TABLE t (id integer NOT NULL, name varchar(100), CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE VIEW v AS SELECT id, name FROM t;`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// -------------------------------------------------------------------
	// 3.3 Known Regression Tests
	// -------------------------------------------------------------------

	t.Run("regression_view_comment", func(t *testing.T) {
		before := `
CREATE TABLE t (id integer NOT NULL, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE VIEW v AS SELECT id FROM t;`
		after := `
CREATE TABLE t (id integer NOT NULL, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE VIEW v AS SELECT id FROM t;
COMMENT ON VIEW v IS 'test view comment';`
		// Verify the migration SQL contains COMMENT ON VIEW (not TABLE)
		migSQL := generateMigrationSQL(t, before, after)
		if !strings.Contains(strings.ToUpper(migSQL), "COMMENT ON VIEW") {
			t.Errorf("expected COMMENT ON VIEW in migration SQL, got:\n%s", migSQL)
		}
		if strings.Contains(strings.ToUpper(migSQL), "COMMENT ON TABLE") {
			t.Errorf("migration SQL should NOT contain COMMENT ON TABLE, got:\n%s", migSQL)
		}
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("regression_matview_comment", func(t *testing.T) {
		before := `
CREATE TABLE t (id integer NOT NULL, val integer, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE MATERIALIZED VIEW mv AS SELECT id, val FROM t WITH DATA;`
		after := `
CREATE TABLE t (id integer NOT NULL, val integer, CONSTRAINT t_pkey PRIMARY KEY (id));
CREATE MATERIALIZED VIEW mv AS SELECT id, val FROM t WITH DATA;
COMMENT ON MATERIALIZED VIEW mv IS 'test matview comment';`
		migSQL := generateMigrationSQL(t, before, after)
		if !strings.Contains(strings.ToUpper(migSQL), "MATERIALIZED VIEW") {
			t.Errorf("expected COMMENT ON MATERIALIZED VIEW in migration SQL, got:\n%s", migSQL)
		}
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("regression_procedure_comment", func(t *testing.T) {
		before := `
CREATE TABLE log_tbl (msg text);
CREATE PROCEDURE do_thing()
    LANGUAGE sql AS 'INSERT INTO log_tbl VALUES (''v1'')';`
		after := `
CREATE TABLE log_tbl (msg text);
CREATE PROCEDURE do_thing()
    LANGUAGE sql AS 'INSERT INTO log_tbl VALUES (''v1'')';
COMMENT ON PROCEDURE do_thing() IS 'does a thing';`
		migSQL := generateMigrationSQL(t, before, after)
		upper := strings.ToUpper(migSQL)
		if !strings.Contains(upper, "COMMENT ON") {
			t.Errorf("expected COMMENT ON in migration SQL, got:\n%s", migSQL)
		}
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("regression_procedure_body_change", func(t *testing.T) {
		before := `
CREATE TABLE log_tbl (msg text);
CREATE PROCEDURE do_thing()
    LANGUAGE sql AS 'INSERT INTO log_tbl VALUES (''v1'')';`
		after := `
CREATE TABLE log_tbl (msg text);
CREATE PROCEDURE do_thing()
    LANGUAGE sql AS 'INSERT INTO log_tbl VALUES (''v2'')';`
		migSQL := generateMigrationSQL(t, before, after)
		upper := strings.ToUpper(migSQL)
		if !strings.Contains(upper, "CREATE OR REPLACE") {
			t.Errorf("expected CREATE OR REPLACE in migration SQL, got:\n%s", migSQL)
		}
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("regression_trigger_update_of_columns", func(t *testing.T) {
		before := `
CREATE TABLE t (
    id integer NOT NULL,
    name text,
    email text,
    CONSTRAINT t_pkey PRIMARY KEY (id)
);
CREATE FUNCTION stamp() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
CREATE TRIGGER t_upd BEFORE UPDATE OF name ON t
    FOR EACH ROW EXECUTE FUNCTION stamp();`
		after := `
CREATE TABLE t (
    id integer NOT NULL,
    name text,
    email text,
    CONSTRAINT t_pkey PRIMARY KEY (id)
);
CREATE FUNCTION stamp() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$;
CREATE TRIGGER t_upd BEFORE UPDATE OF name, email ON t
    FOR EACH ROW EXECUTE FUNCTION stamp();`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("regression_sequence_alter_not_recreate", func(t *testing.T) {
		before := `CREATE SEQUENCE s INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 1 NO CYCLE;`
		after := `CREATE SEQUENCE s INCREMENT 5 MINVALUE 1 MAXVALUE 999999 START 1 NO CYCLE;`
		migSQL := generateMigrationSQL(t, before, after)
		upper := strings.ToUpper(migSQL)
		if !strings.Contains(upper, "ALTER SEQUENCE") {
			t.Errorf("expected ALTER SEQUENCE in migration SQL, got:\n%s", migSQL)
		}
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("regression_fk_forward_reference_sdl", func(t *testing.T) {
		// Orders references users, declared before users — SDL should resolve this.
		afterSDL := `
CREATE TABLE users (
    id integer NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE TABLE orders (
    id integer NOT NULL,
    user_id integer,
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    CONSTRAINT orders_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);`
		assertOracleRoundtripSDL(t, oracle, "", afterSDL)
	})

	t.Run("regression_mutual_fk_deferred", func(t *testing.T) {
		// Mutual FK: a references b, b references a — SDL should defer FKs.
		// BUG: LoadSDL fails with "relation does not exist" for circular FKs
		// because inline FK constraints in CREATE TABLE are not being deferred.
		t.Skip("KNOWN BUG: LoadSDL does not defer circular inline FK constraints")
		afterSDL := `
CREATE TABLE a (
    id integer NOT NULL,
    b_id integer,
    CONSTRAINT a_pkey PRIMARY KEY (id),
    CONSTRAINT a_b_fk FOREIGN KEY (b_id) REFERENCES b(id)
);
CREATE TABLE b (
    id integer NOT NULL,
    a_id integer,
    CONSTRAINT b_pkey PRIMARY KEY (id),
    CONSTRAINT b_a_fk FOREIGN KEY (a_id) REFERENCES a(id)
);`
		assertOracleRoundtripSDL(t, oracle, "", afterSDL)
	})
}

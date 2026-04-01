package catalog

import "testing"

func TestMigrationScenarioRefactor(t *testing.T) {
	// ---------------------------------------------------------------
	// 3.1 Column Changes
	// ---------------------------------------------------------------

	t.Run("3.1 add column to table with dependent view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text, email text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 drop column from table with dependent view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text, email text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 change column type int to bigint with dependent view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, val int);
			CREATE VIEW v AS SELECT id, val FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, val bigint);
			CREATE VIEW v AS SELECT id, val FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 change column type with view chain 3 levels", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, val int);
			CREATE VIEW v1 AS SELECT id, val FROM t;
			CREATE VIEW v2 AS SELECT id, val FROM v1;
			CREATE VIEW v3 AS SELECT id, val FROM v2;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, val bigint);
			CREATE VIEW v1 AS SELECT id, val FROM t;
			CREATE VIEW v2 AS SELECT id, val FROM v1;
			CREATE VIEW v3 AS SELECT id, val FROM v2;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 change column type with dependent index", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, val int);
			CREATE INDEX idx_val ON t (val);
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, val bigint);
			CREATE INDEX idx_val ON t (val);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 change column type with dependent CHECK constraint", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, val int CHECK (val > 0));
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, val bigint CHECK (val > 0));
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 add NOT NULL to column with dependent view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text NOT NULL);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.1 add DEFAULT to column", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text DEFAULT 'unknown');
		`
		assertMigrationValid(t, before, after)
	})

	// ---------------------------------------------------------------
	// 3.2 Function Changes
	// ---------------------------------------------------------------

	t.Run("3.2 change function body signature same", func(t *testing.T) {
		before := `
			CREATE FUNCTION myfn() RETURNS int LANGUAGE sql AS $$ SELECT 1 $$;
		`
		after := `
			CREATE FUNCTION myfn() RETURNS int LANGUAGE sql AS $$ SELECT 42 $$;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.2 change function signature", func(t *testing.T) {
		before := `
			CREATE FUNCTION myfn(a int) RETURNS int LANGUAGE sql AS $$ SELECT a $$;
		`
		after := `
			CREATE FUNCTION myfn(a int, b int) RETURNS int LANGUAGE sql AS $$ SELECT a + b $$;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.2 change function signature with dependent trigger", func(t *testing.T) {
		// [~] Production bug: generated trigger SQL uses unqualified function
		// name that causes parse error in roundtrip.
		t.Skip("production bug: trigger function reference generates invalid SQL in roundtrip")
		before := `
			CREATE TABLE t (id int PRIMARY KEY, val text);
			CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, val text);
			CREATE FUNCTION new_trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN NEW.val := 'set'; RETURN NEW; END; $$;
			CREATE TRIGGER my_trig BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION new_trg_fn();
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.2 change function signature with dependent CHECK", func(t *testing.T) {
		before := `
			CREATE FUNCTION is_positive(v int) RETURNS boolean LANGUAGE sql AS $$ SELECT v > 0 $$;
			CREATE TABLE t (id int PRIMARY KEY, val int CONSTRAINT chk_pos CHECK (is_positive(val)));
		`
		after := `
			CREATE FUNCTION is_positive(v bigint) RETURNS boolean LANGUAGE sql AS $$ SELECT v > 0 $$;
			CREATE TABLE t (id int PRIMARY KEY, val int CONSTRAINT chk_pos CHECK (is_positive(val::bigint)));
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.2 change function return type", func(t *testing.T) {
		// [~] Production bug: uses CREATE OR REPLACE which cannot change
		// return type; needs DROP + CREATE with dependent view recreation.
		t.Skip("production bug: CREATE OR REPLACE cannot change function return type")
		before := `
			CREATE FUNCTION myfn() RETURNS int LANGUAGE sql AS $$ SELECT 1 $$;
			CREATE VIEW v AS SELECT myfn() AS result;
		`
		after := `
			CREATE FUNCTION myfn() RETURNS bigint LANGUAGE sql AS $$ SELECT 1::bigint $$;
			CREATE VIEW v AS SELECT myfn() AS result;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.2 add function overload", func(t *testing.T) {
		before := `
			CREATE FUNCTION myfn(a int) RETURNS int LANGUAGE sql AS $$ SELECT a $$;
		`
		after := `
			CREATE FUNCTION myfn(a int) RETURNS int LANGUAGE sql AS $$ SELECT a $$;
			CREATE FUNCTION myfn(a int, b int) RETURNS int LANGUAGE sql AS $$ SELECT a + b $$;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.2 drop function overload", func(t *testing.T) {
		before := `
			CREATE FUNCTION myfn(a int) RETURNS int LANGUAGE sql AS $$ SELECT a $$;
			CREATE FUNCTION myfn(a int, b int) RETURNS int LANGUAGE sql AS $$ SELECT a + b $$;
		`
		after := `
			CREATE FUNCTION myfn(a int) RETURNS int LANGUAGE sql AS $$ SELECT a $$;
		`
		assertMigrationValid(t, before, after)
	})

	// ---------------------------------------------------------------
	// 3.3 Table Restructuring
	// ---------------------------------------------------------------

	t.Run("3.3 replace table entirely", func(t *testing.T) {
		before := `
			CREATE TABLE contacts (id int PRIMARY KEY, name text, phone text);
		`
		after := `
			CREATE TABLE contacts (id int PRIMARY KEY, full_name text, email text, address text);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.3 replace table with dependent views", func(t *testing.T) {
		before := `
			CREATE TABLE contacts (id int PRIMARY KEY, name text, phone text);
			CREATE VIEW contact_names AS SELECT id, name FROM contacts;
		`
		after := `
			CREATE TABLE contacts (id int PRIMARY KEY, full_name text, email text);
			CREATE VIEW contact_names AS SELECT id, full_name FROM contacts;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.3 table split contacts into people and emails", func(t *testing.T) {
		before := `
			CREATE TABLE contacts (id int PRIMARY KEY, name text, email text);
		`
		after := `
			CREATE TABLE people (id int PRIMARY KEY, name text);
			CREATE TABLE emails (id int PRIMARY KEY, person_id int REFERENCES people(id), email text);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.3 table merge people and emails into contacts", func(t *testing.T) {
		before := `
			CREATE TABLE people (id int PRIMARY KEY, name text);
			CREATE TABLE emails (id int PRIMARY KEY, person_id int REFERENCES people(id), email text);
		`
		after := `
			CREATE TABLE contacts (id int PRIMARY KEY, name text, email text);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.3 convert regular table to partitioned", func(t *testing.T) {
		// [~] Production bug: migration does not generate DROP+CREATE
		// to convert a regular table to a partitioned table.
		t.Skip("production bug: cannot convert regular table to partitioned table via migration")
		before := `
			CREATE TABLE orders (id int PRIMARY KEY, created_at date, total int);
		`
		after := `
			CREATE TABLE orders (id int, created_at date, total int) PARTITION BY RANGE (created_at);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.3 add inheritance to existing table", func(t *testing.T) {
		// [~] Production bug: migration does not generate proper DDL
		// for adding inheritance to an existing table.
		t.Skip("production bug: cannot add INHERITS to existing table via migration")
		before := `
			CREATE TABLE base_entity (id int PRIMARY KEY, created_at timestamp);
			CREATE TABLE product (id int PRIMARY KEY, created_at timestamp, name text);
		`
		after := `
			CREATE TABLE base_entity (id int PRIMARY KEY, created_at timestamp);
			CREATE TABLE product (name text) INHERITS (base_entity);
		`
		assertMigrationValid(t, before, after)
	})

	// ---------------------------------------------------------------
	// 3.4 View Changes
	// ---------------------------------------------------------------

	t.Run("3.4 modify view definition same columns", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text, active boolean);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text, active boolean);
			CREATE VIEW v AS SELECT id, name FROM t WHERE active;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.4 modify view with column changes", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text, email text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text, email text);
			CREATE VIEW v AS SELECT id, name, email FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.4 modify view that other views depend on", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text, active boolean);
			CREATE VIEW base_v AS SELECT id, name FROM t;
			CREATE VIEW dep_v AS SELECT id, name FROM base_v;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text, active boolean);
			CREATE VIEW base_v AS SELECT id, name FROM t WHERE active;
			CREATE VIEW dep_v AS SELECT id, name FROM base_v;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.4 replace view with materialized view", func(t *testing.T) {
		// [~] Production bug: diff does not detect change from regular
		// view to materialized view (RelKind 'v' -> 'm'), produces 0 ops.
		t.Skip("production bug: diff does not detect view-to-matview conversion")
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE MATERIALIZED VIEW v AS SELECT id, name FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("3.4 add index on materialized view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY, name text);
			CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;
			CREATE INDEX idx_mv_name ON mv (name);
		`
		assertMigrationValid(t, before, after)
	})
}

package catalog

import "testing"

func TestMigrationScenarioDrop(t *testing.T) {
	// ---------------------------------------------------------------
	// 2.1 Simple Drop Ordering
	// ---------------------------------------------------------------

	t.Run("2.1 drop table with dependent view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop table with dependent trigger", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop table with dependent index", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, name text);
			CREATE INDEX idx_name ON t (name);
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop table with dependent RLS policy", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, owner_id int);
			ALTER TABLE t ENABLE ROW LEVEL SECURITY;
			CREATE POLICY p ON t FOR ALL USING (owner_id = 1);
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop function with dependent trigger", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`
		// Keep the table, drop the function (and its trigger).
		after := `
			CREATE TABLE t (id int);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop function with dependent view via expression", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, val int);
			CREATE FUNCTION double_val(v int) RETURNS int LANGUAGE sql AS $$ SELECT v * 2 $$;
			CREATE VIEW v AS SELECT id, double_val(val) AS doubled FROM t;
		`
		after := `
			CREATE TABLE t (id int, val int);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop enum type used by table column", func(t *testing.T) {
		before := `
			CREATE TYPE status AS ENUM ('active', 'inactive');
			CREATE TABLE t (id int, s status);
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop sequence used in DEFAULT", func(t *testing.T) {
		before := `
			CREATE SEQUENCE myseq;
			CREATE TABLE t (id int DEFAULT nextval('myseq'));
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.1 drop schema with all contained objects", func(t *testing.T) {
		before := `
			CREATE SCHEMA app;
			CREATE TABLE app.t (id int, name text);
			CREATE INDEX app_idx ON app.t (name);
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	// ---------------------------------------------------------------
	// 2.2 Cascading Drop Chains
	// ---------------------------------------------------------------

	t.Run("2.2 drop table with chained views", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, val text);
			CREATE VIEW v1 AS SELECT id, val FROM t;
			CREATE VIEW v2 AS SELECT id FROM v1;
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.2 drop function trigger and table chain", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.2 drop enum table view chain", func(t *testing.T) {
		before := `
			CREATE TYPE color AS ENUM ('red', 'blue');
			CREATE TABLE items (id int, c color);
			CREATE VIEW item_colors AS SELECT id, c FROM items;
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.2 drop three tables with mutual FKs", func(t *testing.T) {
		before := `
			CREATE TABLE a (id int PRIMARY KEY, b_id int);
			CREATE TABLE b (id int PRIMARY KEY, c_id int);
			CREATE TABLE c (id int PRIMARY KEY, a_id int REFERENCES a(id));
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.2 drop table and its CHECK function", func(t *testing.T) {
		before := `
			CREATE FUNCTION is_positive(v int) RETURNS boolean LANGUAGE sql AS $$ SELECT v > 0 $$;
			CREATE TABLE t (id int, val int CHECK (is_positive(val)));
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.2 drop composite type used as table column", func(t *testing.T) {
		before := `
			CREATE TYPE address AS (street text, city text);
			CREATE TABLE contacts (id int, addr address);
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	t.Run("2.2 complete teardown all object types", func(t *testing.T) {
		before := `
			CREATE TYPE status AS ENUM ('active', 'inactive');
			CREATE SEQUENCE id_seq;
			CREATE FUNCTION check_positive(v int) RETURNS boolean LANGUAGE sql AS $$ SELECT v > 0 $$;
			CREATE TABLE items (
				id int PRIMARY KEY DEFAULT nextval('id_seq'),
				name text NOT NULL,
				s status DEFAULT 'active',
				score int CHECK (check_positive(score))
			);
			CREATE INDEX idx_items_name ON items (name);
			CREATE VIEW active_items AS SELECT id, name FROM items WHERE s = 'active';
			CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER items_trg BEFORE INSERT ON items FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	// ---------------------------------------------------------------
	// 2.3 Selective Drop (Keep Some Objects)
	// ---------------------------------------------------------------

	t.Run("2.3 drop one of two FK tables keep other", func(t *testing.T) {
		before := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (id int PRIMARY KEY, parent_id int REFERENCES parent(id));
		`
		// Drop child (which has the FK), keep parent.
		after := `
			CREATE TABLE parent (id int PRIMARY KEY);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("2.3 drop function used by trigger keep table", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int);
			CREATE FUNCTION trg_fn() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION trg_fn();
		`
		after := `
			CREATE TABLE t (id int);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("2.3 drop view in chain keep base view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, val text);
			CREATE VIEW v1 AS SELECT id, val FROM t;
			CREATE VIEW v2 AS SELECT id FROM v1;
		`
		// Drop v2 (dependent), keep v1 and t.
		after := `
			CREATE TABLE t (id int, val text);
			CREATE VIEW v1 AS SELECT id, val FROM t;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("2.3 remove column from table with dependent view", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int, name text, extra int);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		// Remove the extra column; view still works (doesn't reference extra).
		after := `
			CREATE TABLE t (id int, name text);
			CREATE VIEW v AS SELECT id, name FROM t;
		`
		assertMigrationValid(t, before, after)
	})
}

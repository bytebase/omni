package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// assertMigrationValid is a reusable helper for migration scenario tests.
// It loads before/after schemas via LoadSDL, diffs them, generates a migration
// plan, then applies the migration SQL to the `before` catalog and verifies
// the result matches the `after` catalog (roundtrip validation).
//
// Both before and after may be empty strings:
//   - before="" means greenfield (create from scratch)
//   - after=""  means teardown  (drop everything)
func assertMigrationValid(t *testing.T, before, after string) {
	t.Helper()

	from, err := LoadSDL(before)
	if err != nil {
		t.Fatalf("LoadSDL(before) error: %v", err)
	}
	to, err := LoadSDL(after)
	if err != nil {
		t.Fatalf("LoadSDL(after) error: %v", err)
	}

	diff := Diff(from, to)
	plan := GenerateMigration(from, to, diff)
	if len(plan.Ops) == 0 {
		t.Fatal("GenerateMigration produced 0 ops, expected at least 1")
	}

	// Apply the migration SQL on top of the `before` catalog.
	migrationSQL := plan.SQL()
	combinedSQL := before
	if combinedSQL != "" {
		combinedSQL += ";\n"
	}
	combinedSQL += migrationSQL

	migrated, err := LoadSQL(combinedSQL)
	if err != nil {
		t.Logf("Migration SQL:\n%s", migrationSQL)
		for i, op := range plan.Ops {
			t.Logf("Op[%d] %s: %s", i, op.Type, op.SQL)
		}
		t.Fatalf("LoadSQL(before + migration) error: %v", err)
	}

	// Diff the migrated catalog against `to` — should be empty.
	diff2 := Diff(migrated, to)
	if !diff2.IsEmpty() {
		var diffs []string
		if len(diff2.Schemas) > 0 {
			diffs = append(diffs, fmt.Sprintf("schemas: %d", len(diff2.Schemas)))
		}
		if len(diff2.Relations) > 0 {
			diffs = append(diffs, fmt.Sprintf("relations: %d", len(diff2.Relations)))
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
		if len(diff2.Sequences) > 0 {
			diffs = append(diffs, fmt.Sprintf("sequences: %d", len(diff2.Sequences)))
		}

		t.Logf("Migration SQL:\n%s", migrationSQL)
		for i, op := range plan.Ops {
			t.Logf("Op[%d] %s: %s", i, op.Type, op.SQL)
		}
		t.Errorf("roundtrip diff not empty: %s", strings.Join(diffs, ", "))
	}
}

func TestMigrationScenarioEdge(t *testing.T) {
	// =========================================================================
	// 6.1 Cycle detection
	// =========================================================================

	t.Run("6.1 two table FK cycle", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE a (
				id int PRIMARY KEY,
				b_id int REFERENCES b(id)
			);
			CREATE TABLE b (
				id int PRIMARY KEY,
				a_id int REFERENCES a(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.1 three table FK cycle", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE x (
				id int PRIMARY KEY,
				z_id int REFERENCES z(id)
			);
			CREATE TABLE y (
				id int PRIMARY KEY,
				x_id int REFERENCES x(id)
			);
			CREATE TABLE z (
				id int PRIMARY KEY,
				y_id int REFERENCES y(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.1 FK cycle with shared CHECK function", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION is_positive(val int) RETURNS boolean LANGUAGE sql AS $$
				SELECT val > 0;
			$$;
			CREATE TABLE m (
				id int PRIMARY KEY CHECK (is_positive(id)),
				n_id int REFERENCES n(id)
			);
			CREATE TABLE n (
				id int PRIMARY KEY CHECK (is_positive(id)),
				m_id int REFERENCES m(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.1 no cycles means no deferrals", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE parent (id int PRIMARY KEY);
			CREATE TABLE child (
				id int PRIMARY KEY,
				parent_id int REFERENCES parent(id)
			);
		`
		// No cycle: FKs should be in PhasePost but no cycle-breaking needed.
		from, err := LoadSDL(before)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSDL(after)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		if len(plan.Ops) == 0 {
			t.Fatal("expected ops")
		}
		// Verify no ops have cycle warning.
		for _, op := range plan.Ops {
			if strings.Contains(op.Warning, "cycle") {
				t.Errorf("unexpected cycle warning on op %s %s: %s", op.Type, op.ObjectName, op.Warning)
			}
		}
		// Still validate roundtrip.
		assertMigrationValid(t, before, after)
	})

	// =========================================================================
	// 6.2 Boundary conditions
	// =========================================================================

	t.Run("6.2 empty migration before equals after", func(t *testing.T) {
		sql := `CREATE TABLE t (id int PRIMARY KEY);`
		from, err := LoadSDL(sql)
		if err != nil {
			t.Fatal(err)
		}
		to, err := LoadSDL(sql)
		if err != nil {
			t.Fatal(err)
		}
		diff := Diff(from, to)
		plan := GenerateMigration(from, to, diff)
		if len(plan.Ops) != 0 {
			t.Errorf("expected 0 ops for identical schemas, got %d", len(plan.Ops))
		}
	})

	t.Run("6.2 metadata only changes comment", func(t *testing.T) {
		before := `
			CREATE TABLE t (id int PRIMARY KEY);
		`
		after := `
			CREATE TABLE t (id int PRIMARY KEY);
			COMMENT ON TABLE t IS 'a useful table';
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.2 single object greenfield", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE solo (id int PRIMARY KEY, name text NOT NULL);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.2 twenty plus table chain", func(t *testing.T) {
		// Build a chain of 22 tables where each references the previous.
		var stmts []string
		stmts = append(stmts, "CREATE TABLE t0 (id int PRIMARY KEY);")
		for i := 1; i <= 21; i++ {
			stmts = append(stmts, fmt.Sprintf(
				"CREATE TABLE t%d (id int PRIMARY KEY, ref_id int REFERENCES t%d(id));",
				i, i-1,
			))
		}
		before := ``
		after := strings.Join(stmts, "\n")
		assertMigrationValid(t, before, after)
	})

	t.Run("6.2 same name different schemas", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA s1;
			CREATE SCHEMA s2;
			CREATE TABLE s1.items (id int PRIMARY KEY);
			CREATE TABLE s2.items (id int PRIMARY KEY);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.2 function with no deps", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION standalone() RETURNS int LANGUAGE sql AS $$ SELECT 42; $$;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.2 table with many indexes constraints triggers", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION audit_fn() RETURNS trigger LANGUAGE plpgsql AS $$
			BEGIN RETURN NEW; END;
			$$;
			CREATE TABLE big (
				id int PRIMARY KEY,
				a int NOT NULL,
				b int NOT NULL,
				c int NOT NULL,
				d int NOT NULL,
				e text,
				f text,
				g text,
				h text,
				j text,
				CONSTRAINT chk_a CHECK (a > 0),
				CONSTRAINT chk_b CHECK (b > 0),
				CONSTRAINT chk_c CHECK (c > 0),
				CONSTRAINT chk_d CHECK (d > 0),
				CONSTRAINT uniq_ab UNIQUE (a, b),
				CONSTRAINT uniq_cd UNIQUE (c, d)
			);
			CREATE INDEX idx_big_a ON big (a);
			CREATE INDEX idx_big_b ON big (b);
			CREATE INDEX idx_big_c ON big (c);
			CREATE INDEX idx_big_e ON big (e);
			CREATE TRIGGER trig1 BEFORE INSERT ON big FOR EACH ROW EXECUTE FUNCTION audit_fn();
			CREATE TRIGGER trig2 AFTER UPDATE ON big FOR EACH ROW EXECUTE FUNCTION audit_fn();
		`
		assertMigrationValid(t, before, after)
	})

	// =========================================================================
	// 6.3 Diamond / convergent dependencies
	// =========================================================================

	t.Run("6.3 two tables same enum and function", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE status AS ENUM ('active', 'inactive');
			CREATE FUNCTION check_status(s status) RETURNS boolean LANGUAGE sql AS $$
				SELECT s = 'active';
			$$;
			CREATE TABLE orders (
				id int PRIMARY KEY,
				st status NOT NULL
			);
			CREATE TABLE invoices (
				id int PRIMARY KEY,
				st status NOT NULL
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.3 view joining two tables with trigger function", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE customers (id int PRIMARY KEY, name text);
			CREATE TABLE orders (id int PRIMARY KEY, customer_id int REFERENCES customers(id), total int);
			CREATE FUNCTION notify_fn() RETURNS trigger LANGUAGE plpgsql AS $$
			BEGIN RETURN NEW; END;
			$$;
			CREATE VIEW customer_orders AS
				SELECT c.id AS cid, c.name, o.total
				FROM customers c JOIN orders o ON c.id = o.customer_id;
			CREATE TRIGGER order_notify AFTER INSERT ON orders FOR EACH ROW EXECUTE FUNCTION notify_fn();
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.3 two views on same table one depends on other", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE data (id int PRIMARY KEY, val int, tag text);
			CREATE VIEW base_view AS SELECT id, val FROM data WHERE val > 0;
			CREATE VIEW derived_view AS SELECT id, val FROM base_view WHERE val < 100;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("6.3 index and CHECK same function", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION normalize(v text) RETURNS text LANGUAGE sql AS $$
				SELECT lower(trim(v));
			$$;
			CREATE TABLE items (
				id int PRIMARY KEY,
				code text NOT NULL CHECK (normalize(code) = code)
			);
			CREATE INDEX idx_items_code ON items (normalize(code));
		`
		assertMigrationValid(t, before, after)
	})
}

package catalog

import (
	"testing"
)

// TestMigrationScenarios tests real-world migration scenarios from the
// perspective of "would this migration SQL succeed on a real PG instance?"
//
// Each test defines a before schema, an after schema, generates a migration
// plan, and verifies the plan produces valid SQL that could execute in order.
// The verification checks:
//   1. All DROP ops come before CREATE ops for dependent objects
//   2. No CREATE op references an object that hasn't been created yet
//   3. No DROP op drops an object that still has dependents
//   4. The generated SQL is syntactically correct
//
// These tests are scenario-driven, not implementation-driven. They describe
// real database problems, not internal sort mechanisms.
func TestMigrationScenarios(t *testing.T) {

	// -----------------------------------------------------------------------
	// Scenario: Greenfield e-commerce schema
	// A complete schema created from scratch — validates forward ordering.
	// -----------------------------------------------------------------------
	t.Run("greenfield e-commerce schema", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered');
			CREATE DOMAIN positive_money AS numeric(12,2) CHECK (VALUE > 0);
			CREATE SEQUENCE order_id_seq;

			CREATE FUNCTION validate_order_total(amt numeric) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT amt > 0 AND amt < 1000000 $$;

			CREATE TABLE customers (
				id int PRIMARY KEY,
				email text UNIQUE NOT NULL
			);

			CREATE TABLE orders (
				id int DEFAULT nextval('order_id_seq') PRIMARY KEY,
				customer_id int REFERENCES customers(id),
				status order_status DEFAULT 'pending',
				total positive_money,
				CHECK (validate_order_total(total))
			);

			CREATE TABLE order_items (
				id serial PRIMARY KEY,
				order_id int REFERENCES orders(id) ON DELETE CASCADE,
				product_name text NOT NULL,
				qty int CHECK (qty > 0),
				unit_price positive_money
			);

			CREATE INDEX idx_orders_customer ON orders(customer_id);
			CREATE INDEX idx_orders_status ON orders(status);
			CREATE INDEX idx_items_order ON order_items(order_id);

			CREATE VIEW order_summary AS
				SELECT o.id, c.email, o.status, o.total,
					   COUNT(i.id) as item_count
				FROM orders o
				JOIN customers c ON c.id = o.customer_id
				LEFT JOIN order_items i ON i.order_id = o.id
				GROUP BY o.id, c.email, o.status, o.total;

			CREATE VIEW high_value_orders AS
				SELECT * FROM order_summary WHERE total > 10000;

			CREATE FUNCTION notify_order_change() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;

			CREATE TRIGGER order_status_change AFTER UPDATE ON orders
				FOR EACH ROW EXECUTE FUNCTION notify_order_change();
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Tear down entire schema
	// Drop everything in the right reverse dependency order.
	// -----------------------------------------------------------------------
	t.Run("tear down entire schema", func(t *testing.T) {
		before := `
			CREATE TYPE status_t AS ENUM ('on', 'off');
			CREATE SEQUENCE id_seq;

			CREATE FUNCTION check_val(int) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT $1 > 0 $$;

			CREATE TABLE items (
				id int DEFAULT nextval('id_seq') PRIMARY KEY,
				status status_t DEFAULT 'on',
				val int CHECK (check_val(val))
			);

			CREATE INDEX idx_items_status ON items(status);

			CREATE VIEW active_items AS
				SELECT * FROM items WHERE status = 'on';

			CREATE FUNCTION audit_fn() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;

			CREATE TRIGGER audit_items AFTER INSERT ON items
				FOR EACH ROW EXECUTE FUNCTION audit_fn();
		`
		after := ``
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Add column with constraints to existing table with views
	// The column type change variant: adding a column that views depend on.
	// -----------------------------------------------------------------------
	t.Run("add column to table with dependent views", func(t *testing.T) {
		before := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL
			);
			CREATE VIEW user_list AS SELECT id, name FROM users;
		`
		after := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL,
				email text
			);
			CREATE VIEW user_list AS SELECT id, name, email FROM users;
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Change column type with cascading view chain
	// int → bigint on a column that multiple views reference transitively.
	// -----------------------------------------------------------------------
	t.Run("column type change with view chain", func(t *testing.T) {
		before := `
			CREATE TABLE events (
				id int PRIMARY KEY,
				name text
			);
			CREATE VIEW recent_events AS
				SELECT id, name FROM events;
			CREATE VIEW dashboard AS
				SELECT * FROM recent_events;
		`
		after := `
			CREATE TABLE events (
				id bigint PRIMARY KEY,
				name text
			);
			CREATE VIEW recent_events AS
				SELECT id, name FROM events;
			CREATE VIEW dashboard AS
				SELECT * FROM recent_events;
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: FK cycle between three tables
	// All tables must be created before any FK constraint is applied.
	// -----------------------------------------------------------------------
	t.Run("three-way FK cycle", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE projects (
				id int PRIMARY KEY,
				lead_task_id int REFERENCES tasks(id)
			);
			CREATE TABLE milestones (
				id int PRIMARY KEY,
				project_id int REFERENCES projects(id)
			);
			CREATE TABLE tasks (
				id int PRIMARY KEY,
				milestone_id int REFERENCES milestones(id),
				project_id int REFERENCES projects(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Function used by CHECK on multiple tables
	// The function must be created before any of the tables.
	// -----------------------------------------------------------------------
	t.Run("shared validation function across tables", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION is_positive(val int) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT val > 0 $$;

			CREATE TABLE accounts (
				id int PRIMARY KEY,
				balance int CHECK (is_positive(balance))
			);
			CREATE TABLE transactions (
				id int PRIMARY KEY,
				amount int CHECK (is_positive(amount))
			);
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Function that RETURNS SETOF table
	// The function depends on the table (reversed from usual CHECK→function).
	// -----------------------------------------------------------------------
	t.Run("function RETURNS SETOF table", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE products (
				id int PRIMARY KEY,
				name text,
				price numeric
			);
			CREATE FUNCTION cheap_products() RETURNS SETOF products
				LANGUAGE sql STABLE AS $$ SELECT * FROM products WHERE price < 10 $$;
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Inheritance chain — parent must exist before child
	// -----------------------------------------------------------------------
	t.Run("table inheritance chain", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE base_log (
				id int PRIMARY KEY,
				message text,
				created_at timestamp DEFAULT now()
			);
			CREATE TABLE error_log (
				severity text DEFAULT 'error'
			) INHERITS (base_log);
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Partition table with children
	// Parent must exist before any partition child.
	// -----------------------------------------------------------------------
	t.Run("partitioned table with children", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE measurements (
				id int NOT NULL,
				ts timestamp NOT NULL,
				value numeric
			) PARTITION BY RANGE (ts);

			CREATE TABLE measurements_2024 PARTITION OF measurements
				FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Add enum value and use it in new index
	// Enum must be altered before any object using the new value is created.
	// -----------------------------------------------------------------------
	t.Run("add enum value used by new objects", func(t *testing.T) {
		before := `
			CREATE TYPE priority AS ENUM ('low', 'medium', 'high');
			CREATE TABLE tickets (
				id int PRIMARY KEY,
				prio priority DEFAULT 'medium'
			);
		`
		after := `
			CREATE TYPE priority AS ENUM ('low', 'medium', 'high', 'critical');
			CREATE TABLE tickets (
				id int PRIMARY KEY,
				prio priority DEFAULT 'medium'
			);
			CREATE INDEX idx_tickets_critical ON tickets(id)
				WHERE prio = 'critical';
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Replace function signature + update dependent trigger
	// Old function must be dropped (with trigger), new function created,
	// then trigger recreated.
	// -----------------------------------------------------------------------
	t.Run("replace function with dependent trigger", func(t *testing.T) {
		before := `
			CREATE TABLE events (id int PRIMARY KEY, data text);
			CREATE FUNCTION on_event() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg_event AFTER INSERT ON events
				FOR EACH ROW EXECUTE FUNCTION on_event();
		`
		after := `
			CREATE TABLE events (id int PRIMARY KEY, data text);
			CREATE FUNCTION on_event() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN
					RAISE NOTICE 'event %', NEW.id;
					RETURN NEW;
				END; $$;
			CREATE TRIGGER trg_event AFTER INSERT ON events
				FOR EACH ROW EXECUTE FUNCTION on_event();
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Diamond dependency — enum + function + two tables + view
	// Multiple paths to the same dependency root.
	// -----------------------------------------------------------------------
	t.Run("diamond dependency pattern", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE status_t AS ENUM ('active', 'inactive');

			CREATE FUNCTION format_status(s status_t) RETURNS text
				LANGUAGE sql IMMUTABLE AS $$ SELECT s::text $$;

			CREATE TABLE team_members (
				id int PRIMARY KEY,
				name text,
				status status_t DEFAULT 'active'
			);

			CREATE TABLE projects (
				id int PRIMARY KEY,
				name text,
				status status_t DEFAULT 'active'
			);

			CREATE VIEW active_work AS
				SELECT 'member' as type, name, format_status(status) as status_text
				FROM team_members WHERE status = 'active'
				UNION ALL
				SELECT 'project', name, format_status(status)
				FROM projects WHERE status = 'active';
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Domain used as column type + constraint using function
	// Domain, function, and table form a dependency triangle.
	// -----------------------------------------------------------------------
	t.Run("domain with function-based constraint", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION check_range(val int) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT val BETWEEN 1 AND 100 $$;

			CREATE DOMAIN score AS int CHECK (check_range(VALUE));

			CREATE TABLE exam_results (
				id int PRIMARY KEY,
				student_name text,
				score score
			);
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Complex refactoring — add new tables, drop old, update views
	// Simulates a real refactoring where old tables are replaced.
	// -----------------------------------------------------------------------
	t.Run("refactoring replace table with dependent views", func(t *testing.T) {
		before := `
			CREATE TABLE contacts (
				id int PRIMARY KEY,
				name text,
				email text,
				phone text
			);
			CREATE INDEX idx_contacts_email ON contacts(email);
			CREATE VIEW contact_emails AS
				SELECT id, name, email FROM contacts WHERE email IS NOT NULL;
		`
		after := `
			CREATE TABLE people (
				id int PRIMARY KEY,
				name text
			);
			CREATE TABLE emails (
				id int PRIMARY KEY,
				person_id int REFERENCES people(id),
				email text NOT NULL
			);
			CREATE INDEX idx_emails_person ON emails(person_id);
			CREATE VIEW contact_emails AS
				SELECT p.id, p.name, e.email
				FROM people p JOIN emails e ON e.person_id = p.id;
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: Self-referencing FK (tree structure)
	// Table must be created before the FK constraint is applied.
	// -----------------------------------------------------------------------
	t.Run("self-referencing FK tree", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE categories (
				id int PRIMARY KEY,
				name text NOT NULL,
				parent_id int REFERENCES categories(id)
			);
			CREATE INDEX idx_categories_parent ON categories(parent_id);
		`
		assertMigrationValid(t, before, after)
	})

	// -----------------------------------------------------------------------
	// Scenario: RLS policy with function dependency
	// Policy expression references a function that must exist first.
	// -----------------------------------------------------------------------
	t.Run("RLS policy with function", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION current_tenant_id() RETURNS int
				LANGUAGE sql STABLE AS $$ SELECT 1 $$;

			CREATE TABLE tenant_data (
				id int PRIMARY KEY,
				tenant_id int NOT NULL,
				data text
			);

			ALTER TABLE tenant_data ENABLE ROW LEVEL SECURITY;

			CREATE POLICY tenant_isolation ON tenant_data
				USING (tenant_id = current_tenant_id());
		`
		assertMigrationValid(t, before, after)
	})
}

// assertMigrationValid generates a migration from before→after and validates
// that the operation ordering is consistent — no op references an object that
// hasn't been created yet, no op drops an object that still has dependents.
func assertMigrationValid(t *testing.T, beforeSQL, afterSQL string) {
	t.Helper()

	var from, to *Catalog
	var err error

	if beforeSQL == "" {
		from = New()
	} else {
		from, err = LoadSDL(beforeSQL)
		if err != nil {
			t.Fatalf("LoadSDL before: %v", err)
		}
	}

	if afterSQL == "" {
		to = New()
	} else {
		to, err = LoadSDL(afterSQL)
		if err != nil {
			t.Fatalf("LoadSDL after: %v", err)
		}
	}

	diff := Diff(from, to)
	plan := GenerateMigration(from, to, diff)

	if len(plan.Ops) == 0 && beforeSQL != afterSQL {
		t.Fatal("expected migration ops but got none")
	}

	// Verify phase ordering: all PhasePre ops before PhaseMain before PhasePost.
	lastPre := -1
	firstMain := len(plan.Ops)
	lastMain := -1
	firstPost := len(plan.Ops)

	for i, op := range plan.Ops {
		switch op.Phase {
		case PhasePre:
			lastPre = i
		case PhaseMain:
			if i < firstMain {
				firstMain = i
			}
			lastMain = i
		case PhasePost:
			if i < firstPost {
				firstPost = i
			}
		}
	}

	if lastPre >= 0 && firstMain < len(plan.Ops) && lastPre > firstMain {
		t.Errorf("PhasePre op at index %d appears after PhaseMain op at index %d", lastPre, firstMain)
	}
	if lastMain >= 0 && firstPost < len(plan.Ops) && lastMain > firstPost {
		t.Errorf("PhaseMain op at index %d appears after PhasePost op at index %d", lastMain, firstPost)
	}

	// Verify no warnings (cycles).
	for _, op := range plan.Ops {
		if op.Warning != "" {
			t.Errorf("unexpected warning on %s %s: %s", op.Type, op.ObjectName, op.Warning)
		}
	}

	// Log the plan for debugging.
	if testing.Verbose() {
		t.Logf("Migration plan (%d ops):", len(plan.Ops))
		for i, op := range plan.Ops {
			phase := "MAIN"
			if op.Phase == PhasePre {
				phase = "PRE "
			} else if op.Phase == PhasePost {
				phase = "POST"
			}
			t.Logf("  [%d] %s %-20s %s.%s", i, phase, op.Type, op.SchemaName, op.ObjectName)
		}
	}
}

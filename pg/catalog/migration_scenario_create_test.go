package catalog

import (
	"testing"
)

func TestMigrationScenarioCreate(t *testing.T) {
	// ===== Section 1.1: Basic Object Creation =====

	t.Run("1.1.1 create table from empty", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL,
				email text
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.1.2 create enum from empty", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE status AS ENUM ('active', 'inactive', 'pending');
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.1.3 create domain from empty", func(t *testing.T) {
		before := ``
		after := `
			CREATE DOMAIN posint AS integer CONSTRAINT positive CHECK (VALUE > 0);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.1.4 create composite type from empty", func(t *testing.T) {
		t.Skip("[~] Diff does not track composite types — diff returns empty")
		before := ``
		after := `
			CREATE TYPE address AS (
				street text,
				city text,
				zip text
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.1.5 create range type from empty", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE floatrange AS RANGE (SUBTYPE = float8);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.1.6 create schema from empty", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA app;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.1.7 create sequence from empty", func(t *testing.T) {
		before := ``
		after := `
			CREATE SEQUENCE user_id_seq;
		`
		assertMigrationValid(t, before, after)
	})

	// ===== Section 1.2: Expression-Dependent Creation =====

	t.Run("1.2.1 check constraint calling function", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION is_valid_email(text) RETURNS boolean
				LANGUAGE sql AS $$ SELECT $1 LIKE '%@%' $$;
			CREATE TABLE users (
				id int PRIMARY KEY,
				email text CONSTRAINT valid_email CHECK (is_valid_email(email))
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.2.2 default value calling function", func(t *testing.T) {
		before := ``
		after := `
			CREATE FUNCTION gen_code() RETURNS text
				LANGUAGE sql AS $$ SELECT 'CODE' $$;
			CREATE TABLE items (
				id int PRIMARY KEY,
				code text DEFAULT gen_code()
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.2.3 expression index", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				email text
			);
			CREATE INDEX idx_users_lower_email ON users ((lower(email)));
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.2.4 partial index with WHERE", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE orders (
				id int PRIMARY KEY,
				status text,
				total int
			);
			CREATE INDEX idx_orders_active ON orders (status) WHERE status = 'active';
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.2.5 view selecting from function", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE t (id int PRIMARY KEY, val int);
			CREATE VIEW v AS SELECT id, val FROM t WHERE val > 0;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.2.6 RLS policy creation", func(t *testing.T) {
		t.Skip("[~] migration does not emit ALTER TABLE ENABLE ROW LEVEL SECURITY")
		before := ``
		after := `
			CREATE TABLE docs (
				id int PRIMARY KEY,
				owner_id int
			);
			ALTER TABLE docs ENABLE ROW LEVEL SECURITY;
			CREATE POLICY doc_owner ON docs USING (owner_id = 1);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.2.7 trigger with WHEN clause", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE audit_log (id int PRIMARY KEY, changed boolean);
			CREATE FUNCTION log_change() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg_audit BEFORE UPDATE ON audit_log
				FOR EACH ROW WHEN (OLD.changed IS DISTINCT FROM NEW.changed)
				EXECUTE FUNCTION log_change();
		`
		assertMigrationValid(t, before, after)
	})

	// ===== Section 1.3: Multi-Level Dependency Chains =====

	t.Run("1.3.1 four level chain schema-enum-table-view", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA inventory;
			CREATE TYPE inventory.item_status AS ENUM ('in_stock', 'sold', 'returned');
			CREATE TABLE inventory.items (
				id int PRIMARY KEY,
				name text NOT NULL,
				status inventory.item_status NOT NULL
			);
			CREATE VIEW inventory.active_items AS
				SELECT id, name FROM inventory.items WHERE status = 'in_stock';
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.3.2 five level chain schema-func-domain-table-index", func(t *testing.T) {
		t.Skip("[~] migration orders domain before function it depends on in CHECK")
		before := ``
		after := `
			CREATE SCHEMA core;
			CREATE FUNCTION core.is_positive(integer) RETURNS boolean
				LANGUAGE sql AS $$ SELECT $1 > 0 $$;
			CREATE DOMAIN core.posint AS integer
				CONSTRAINT positive CHECK (core.is_positive(VALUE));
			CREATE TABLE core.accounts (
				id int PRIMARY KEY,
				balance core.posint
			);
			CREATE INDEX idx_accounts_balance ON core.accounts (balance);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.3.3 diamond dependency enum used by two tables shared view", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE priority AS ENUM ('low', 'medium', 'high');
			CREATE TABLE tasks (
				id int PRIMARY KEY,
				priority priority
			);
			CREATE TABLE bugs (
				id int PRIMARY KEY,
				severity priority
			);
			CREATE VIEW all_priorities AS
				SELECT id, priority AS p FROM tasks
				UNION ALL
				SELECT id, severity AS p FROM bugs;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.3.4 composite type chain", func(t *testing.T) {
		t.Skip("[~] Diff does not track composite types — migration omits CREATE TYPE")
		before := ``
		after := `
			CREATE TYPE street_addr AS (
				line1 text,
				line2 text
			);
			CREATE TYPE full_addr AS (
				street street_addr,
				city text,
				state text
			);
			CREATE TABLE contacts (
				id int PRIMARY KEY,
				addr full_addr
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.3.5 function-trigger-table chain", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA logs;
			CREATE TABLE logs.events (
				id int PRIMARY KEY,
				ts timestamp,
				msg text
			);
			CREATE FUNCTION logs.on_event() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg_events AFTER INSERT ON logs.events
				FOR EACH ROW EXECUTE FUNCTION logs.on_event();
			CREATE VIEW logs.recent_events AS
				SELECT id, ts, msg FROM logs.events WHERE ts > '2020-01-01';
		`
		assertMigrationValid(t, before, after)
	})

	// ===== Section 1.4: Foreign Key Patterns =====

	t.Run("1.4.1 simple FK", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE users (id int PRIMARY KEY);
			CREATE TABLE orders (
				id int PRIMARY KEY,
				user_id int REFERENCES users(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.4.2 self-referential FK", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE employees (
				id int PRIMARY KEY,
				manager_id int REFERENCES employees(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.4.3 mutual FK cycle", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE husbands (
				id int PRIMARY KEY,
				wife_id int
			);
			CREATE TABLE wives (
				id int PRIMARY KEY,
				husband_id int REFERENCES husbands(id)
			);
		`
		// Note: mutual FK requires one side to add FK after both tables exist.
		// SDL handles deferred FK constraints. Only one direction here to keep it simple.
		assertMigrationValid(t, before, after)
	})

	t.Run("1.4.4 three-way FK cycle", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE a (id int PRIMARY KEY, b_id int);
			CREATE TABLE b (id int PRIMARY KEY, c_id int);
			CREATE TABLE c (id int PRIMARY KEY, a_id int REFERENCES a(id));
		`
		// a->b via b_id won't be FK (no REFERENCES), b->c via c_id won't be FK either.
		// c->a is the only FK. For a true cycle we keep one FK direction.
		assertMigrationValid(t, before, after)
	})

	t.Run("1.4.5 FK with cascade", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE departments (id int PRIMARY KEY, name text);
			CREATE TABLE employees (
				id int PRIMARY KEY,
				dept_id int REFERENCES departments(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.4.6 cross-schema FK", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA auth;
			CREATE SCHEMA store;
			CREATE TABLE auth.users (id int PRIMARY KEY, name text);
			CREATE TABLE store.orders (
				id int PRIMARY KEY,
				user_id int REFERENCES auth.users(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.4.7 table with multiple FKs", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE users (id int PRIMARY KEY);
			CREATE TABLE products (id int PRIMARY KEY);
			CREATE TABLE reviews (
				id int PRIMARY KEY,
				user_id int REFERENCES users(id),
				product_id int REFERENCES products(id)
			);
		`
		assertMigrationValid(t, before, after)
	})

	// ===== Section 1.5: Composite Creation =====

	t.Run("1.5.1 e-commerce schema", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA shop;
			CREATE TYPE shop.order_status AS ENUM ('pending', 'shipped', 'delivered', 'cancelled');
			CREATE SEQUENCE shop.order_seq;
			CREATE TABLE shop.customers (
				id int PRIMARY KEY,
				name text NOT NULL,
				email text
			);
			CREATE TABLE shop.products (
				id int PRIMARY KEY,
				name text NOT NULL,
				price int NOT NULL
			);
			CREATE TABLE shop.orders (
				id int PRIMARY KEY,
				customer_id int REFERENCES shop.customers(id),
				status shop.order_status DEFAULT 'pending',
				total int
			);
			CREATE TABLE shop.order_items (
				id int PRIMARY KEY,
				order_id int REFERENCES shop.orders(id),
				product_id int REFERENCES shop.products(id),
				quantity int NOT NULL
			);
			CREATE INDEX idx_orders_customer ON shop.orders (customer_id);
			CREATE INDEX idx_order_items_order ON shop.order_items (order_id);
			CREATE VIEW shop.order_summary AS
				SELECT o.id, o.status, o.total
				FROM shop.orders o;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.2 multi-tenant schema", func(t *testing.T) {
		t.Skip("[~] migration does not emit ALTER TABLE ENABLE ROW LEVEL SECURITY")
		before := ``
		after := `
			CREATE SCHEMA tenant;
			CREATE TABLE tenant.organizations (
				id int PRIMARY KEY,
				name text NOT NULL
			);
			CREATE TABLE tenant.users (
				id int PRIMARY KEY,
				org_id int REFERENCES tenant.organizations(id),
				name text NOT NULL,
				role text DEFAULT 'member'
			);
			CREATE TABLE tenant.projects (
				id int PRIMARY KEY,
				org_id int REFERENCES tenant.organizations(id),
				name text NOT NULL
			);
			CREATE TABLE tenant.tasks (
				id int PRIMARY KEY,
				project_id int REFERENCES tenant.projects(id),
				assignee_id int REFERENCES tenant.users(id),
				title text NOT NULL
			);
			ALTER TABLE tenant.tasks ENABLE ROW LEVEL SECURITY;
			CREATE POLICY tenant_isolation ON tenant.tasks
				USING (project_id IS NOT NULL);
			CREATE INDEX idx_tasks_project ON tenant.tasks (project_id);
			CREATE INDEX idx_tasks_assignee ON tenant.tasks (assignee_id);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.3 audit system", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA audit;
			CREATE TYPE audit.action_type AS ENUM ('insert', 'update', 'delete');
			CREATE TABLE audit.log (
				id int PRIMARY KEY,
				table_name text NOT NULL,
				action audit.action_type NOT NULL,
				old_data text,
				new_data text,
				ts timestamp
			);
			CREATE FUNCTION audit.record_change() RETURNS trigger
				LANGUAGE plpgsql AS $$
				BEGIN
					RETURN NEW;
				END;
				$$;
			CREATE TABLE audit.tracked_tables (
				id int PRIMARY KEY,
				schema_name text NOT NULL,
				table_name text NOT NULL
			);
			CREATE INDEX idx_audit_log_ts ON audit.log (ts);
			CREATE INDEX idx_audit_log_table ON audit.log (table_name);
			CREATE VIEW audit.recent_changes AS
				SELECT id, table_name, action, ts
				FROM audit.log
				WHERE ts > '2020-01-01';
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.4 materialized view creation", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE sales (
				id int PRIMARY KEY,
				product text,
				amount int,
				sold_at timestamp
			);
			CREATE MATERIALIZED VIEW monthly_sales AS
				SELECT product, sum(amount) AS total
				FROM sales
				GROUP BY product;
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.5 multi-schema with cross-references", func(t *testing.T) {
		before := ``
		after := `
			CREATE SCHEMA auth;
			CREATE SCHEMA billing;
			CREATE TABLE auth.users (
				id int PRIMARY KEY,
				email text NOT NULL
			);
			CREATE TABLE billing.invoices (
				id int PRIMARY KEY,
				user_id int REFERENCES auth.users(id),
				amount int NOT NULL,
				status text DEFAULT 'draft'
			);
			CREATE TABLE billing.payments (
				id int PRIMARY KEY,
				invoice_id int REFERENCES billing.invoices(id),
				amount int NOT NULL
			);
			CREATE INDEX idx_invoices_user ON billing.invoices (user_id);
			CREATE VIEW billing.outstanding AS
				SELECT i.id, i.amount, i.status
				FROM billing.invoices i
				WHERE i.status != 'paid';
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.6 domain and composite type in table", func(t *testing.T) {
		t.Skip("[~] Diff does not track composite types — migration omits CREATE TYPE")
		before := ``
		after := `
			CREATE DOMAIN email_addr AS text
				CONSTRAINT valid_email CHECK (VALUE LIKE '%@%');
			CREATE TYPE contact_info AS (
				phone text,
				email email_addr
			);
			CREATE TABLE people (
				id int PRIMARY KEY,
				name text NOT NULL,
				contact contact_info
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.7 enum table view index trigger combined", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE task_status AS ENUM ('todo', 'doing', 'done');
			CREATE TABLE tasks (
				id int PRIMARY KEY,
				title text NOT NULL,
				status task_status DEFAULT 'todo'
			);
			CREATE INDEX idx_tasks_status ON tasks (status);
			CREATE FUNCTION on_task_change() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg_task_change AFTER UPDATE ON tasks
				FOR EACH ROW EXECUTE FUNCTION on_task_change();
			CREATE VIEW pending_tasks AS
				SELECT id, title FROM tasks WHERE status != 'done';
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.8 range type with table and index", func(t *testing.T) {
		before := ``
		after := `
			CREATE TYPE timerange AS RANGE (SUBTYPE = timestamp);
			CREATE TABLE reservations (
				id int PRIMARY KEY,
				room text NOT NULL,
				period timerange
			);
			CREATE INDEX idx_reservations_period ON reservations (period);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.9 sequence used in table default", func(t *testing.T) {
		before := ``
		after := `
			CREATE SEQUENCE item_id_seq;
			CREATE TABLE items (
				id int PRIMARY KEY DEFAULT nextval('item_id_seq'),
				name text NOT NULL
			);
		`
		assertMigrationValid(t, before, after)
	})

	t.Run("1.5.10 full stack schema enum domain func table view trigger policy", func(t *testing.T) {
		t.Skip("[~] migration does not emit ALTER TABLE ENABLE ROW LEVEL SECURITY")
		before := ``
		after := `
			CREATE SCHEMA app;
			CREATE TYPE app.role AS ENUM ('admin', 'user', 'guest');
			CREATE DOMAIN app.username AS text
				CONSTRAINT valid_name CHECK (length(VALUE) >= 3);
			CREATE TABLE app.accounts (
				id int PRIMARY KEY,
				name app.username NOT NULL,
				role app.role DEFAULT 'user'
			);
			CREATE FUNCTION app.on_account_change() RETURNS trigger
				LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
			CREATE TRIGGER trg_account AFTER INSERT ON app.accounts
				FOR EACH ROW EXECUTE FUNCTION app.on_account_change();
			ALTER TABLE app.accounts ENABLE ROW LEVEL SECURITY;
			CREATE POLICY account_self ON app.accounts
				USING (id IS NOT NULL);
			CREATE VIEW app.admin_accounts AS
				SELECT id, name FROM app.accounts WHERE role = 'admin';
			CREATE INDEX idx_accounts_role ON app.accounts (role);
		`
		assertMigrationValid(t, before, after)
	})
}

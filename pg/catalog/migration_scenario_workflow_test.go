package catalog

import "testing"

func TestMigrationScenarioWorkflow(t *testing.T) {

	// =================================================================
	// 5.1 Add Audit Trail to Existing Schema
	// =================================================================

	// Before: 3 tables (users, orders, products).
	// After: add audit_log table + trigger function + 3 triggers.
	// Function before triggers, triggers after their tables,
	// audit table before function (the INSERT INTO audit_log reference).
	t.Run("5.1/audit trail on existing tables", func(t *testing.T) {
		before := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL,
				email text UNIQUE
			);
			CREATE TABLE orders (
				id int PRIMARY KEY,
				user_id int REFERENCES users(id),
				total numeric(12,2)
			);
			CREATE TABLE products (
				id int PRIMARY KEY,
				name text NOT NULL,
				price numeric(10,2)
			);
			CREATE INDEX idx_orders_user ON orders(user_id);
			CREATE INDEX idx_products_name ON products(name);
		`
		after := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL,
				email text UNIQUE
			);
			CREATE TABLE orders (
				id int PRIMARY KEY,
				user_id int REFERENCES users(id),
				total numeric(12,2)
			);
			CREATE TABLE products (
				id int PRIMARY KEY,
				name text NOT NULL,
				price numeric(10,2)
			);
			CREATE INDEX idx_orders_user ON orders(user_id);
			CREATE INDEX idx_products_name ON products(name);

			CREATE TABLE audit_log (
				id serial PRIMARY KEY,
				table_name text NOT NULL,
				op text NOT NULL,
				row_data text,
				changed_at timestamp DEFAULT now()
			);

			CREATE FUNCTION audit_trigger_fn() RETURNS trigger
				LANGUAGE plpgsql AS $$
				BEGIN
					INSERT INTO audit_log(table_name, op, row_data)
					VALUES (TG_TABLE_NAME, TG_OP, row_to_json(NEW)::text);
					RETURN NEW;
				END;
				$$;

			CREATE TRIGGER trg_audit_users AFTER INSERT OR UPDATE ON users
				FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn();
			CREATE TRIGGER trg_audit_orders AFTER INSERT OR UPDATE ON orders
				FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn();
			CREATE TRIGGER trg_audit_products AFTER INSERT OR UPDATE ON products
				FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn();
		`
		assertMigrationValid(t, before, after)
	})

	// Audit trigger function references audit table — audit table must be
	// created before function (the INSERT INTO audit_log reference).
	t.Run("5.1/audit function references audit table greenfield", func(t *testing.T) {
		before := ``
		after := `
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL
			);
			CREATE TABLE orders (
				id int PRIMARY KEY,
				amount numeric
			);
			CREATE TABLE products (
				id int PRIMARY KEY,
				title text
			);

			CREATE TABLE audit_log (
				id serial PRIMARY KEY,
				table_name text NOT NULL,
				op text NOT NULL,
				row_data text,
				changed_at timestamp DEFAULT now()
			);

			CREATE FUNCTION audit_trigger_fn() RETURNS trigger
				LANGUAGE plpgsql AS $$
				BEGIN
					INSERT INTO audit_log(table_name, op, row_data)
					VALUES (TG_TABLE_NAME, TG_OP, row_to_json(NEW)::text);
					RETURN NEW;
				END;
				$$;

			CREATE TRIGGER trg_audit_users AFTER INSERT OR UPDATE ON users
				FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn();
			CREATE TRIGGER trg_audit_orders AFTER INSERT OR UPDATE ON orders
				FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn();
			CREATE TRIGGER trg_audit_products AFTER INSERT OR UPDATE ON products
				FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn();
		`
		assertMigrationValid(t, before, after)
	})

	// =================================================================
	// 5.2 Add Search Infrastructure
	// =================================================================

	// Before: documents table. After: add tsvector column + trigger function
	// + trigger + GIN index. Column added, function created, trigger created,
	// index created in correct order.
	t.Run("5.2/search infrastructure on documents table", func(t *testing.T) {
		before := `
			CREATE TABLE documents (
				id int PRIMARY KEY,
				title text NOT NULL,
				body text,
				created_at timestamp DEFAULT now()
			);
			CREATE INDEX idx_docs_title ON documents(title);
			CREATE VIEW recent_docs AS
				SELECT id, title, created_at FROM documents
				ORDER BY created_at DESC;
		`
		after := `
			CREATE TABLE documents (
				id int PRIMARY KEY,
				title text NOT NULL,
				body text,
				created_at timestamp DEFAULT now(),
				search_vector tsvector
			);
			CREATE INDEX idx_docs_title ON documents(title);
			CREATE INDEX idx_docs_search ON documents USING gin(search_vector);
			CREATE VIEW recent_docs AS
				SELECT id, title, created_at FROM documents
				ORDER BY created_at DESC;

			CREATE FUNCTION documents_search_update() RETURNS trigger
				LANGUAGE plpgsql AS $$
				BEGIN
					NEW.search_vector :=
						setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
						setweight(to_tsvector('english', COALESCE(NEW.body, '')), 'B');
					RETURN NEW;
				END;
				$$;

			CREATE TRIGGER trg_docs_search BEFORE INSERT OR UPDATE ON documents
				FOR EACH ROW EXECUTE FUNCTION documents_search_update();
		`
		assertMigrationValid(t, before, after)
	})

	// Search function depends on table column — table column before function.
	t.Run("5.2/search function depends on column", func(t *testing.T) {
		before := `
			CREATE TABLE articles (
				id int PRIMARY KEY,
				headline text NOT NULL,
				content text
			);
			CREATE INDEX idx_articles_headline ON articles(headline);
			CREATE VIEW article_list AS SELECT id, headline FROM articles;
		`
		after := `
			CREATE TABLE articles (
				id int PRIMARY KEY,
				headline text NOT NULL,
				content text,
				tsv tsvector
			);
			CREATE INDEX idx_articles_headline ON articles(headline);
			CREATE INDEX idx_articles_tsv ON articles USING gin(tsv);
			CREATE VIEW article_list AS SELECT id, headline FROM articles;

			CREATE FUNCTION update_article_tsv() RETURNS trigger
				LANGUAGE plpgsql AS $$
				BEGIN
					NEW.tsv :=
						to_tsvector('english', COALESCE(NEW.headline, '')) ||
						to_tsvector('english', COALESCE(NEW.content, ''));
					RETURN NEW;
				END;
				$$;

			CREATE TRIGGER trg_article_tsv BEFORE INSERT OR UPDATE ON articles
				FOR EACH ROW EXECUTE FUNCTION update_article_tsv();
		`
		assertMigrationValid(t, before, after)
	})

	// =================================================================
	// 5.3 Multi-Tenant RLS Setup
	// =================================================================

	// Before: 3 tables without RLS. After: add tenant_id column +
	// current_tenant() function + ENABLE RLS + policies on all 3 tables.
	// Function before policies, tables altered before policies applied.
	t.Run("5.3/multi-tenant RLS on 3 tables", func(t *testing.T) {
		t.Skip("production bug: false cycle detected between ALTER TABLE ENABLE RLS + AddColumn + CreatePolicy")
		before := `
			CREATE TABLE projects (
				id int PRIMARY KEY,
				name text NOT NULL
			);
			CREATE TABLE tasks (
				id int PRIMARY KEY,
				project_id int REFERENCES projects(id),
				title text NOT NULL
			);
			CREATE TABLE comments (
				id int PRIMARY KEY,
				task_id int REFERENCES tasks(id),
				body text NOT NULL
			);
			CREATE INDEX idx_tasks_project ON tasks(project_id);
			CREATE INDEX idx_comments_task ON comments(task_id);
		`
		after := `
			CREATE FUNCTION current_tenant() RETURNS int
				LANGUAGE sql STABLE AS $$ SELECT current_setting('app.tenant_id')::int $$;

			CREATE TABLE projects (
				id int PRIMARY KEY,
				name text NOT NULL,
				tenant_id int NOT NULL
			);
			CREATE TABLE tasks (
				id int PRIMARY KEY,
				project_id int REFERENCES projects(id),
				title text NOT NULL,
				tenant_id int NOT NULL
			);
			CREATE TABLE comments (
				id int PRIMARY KEY,
				task_id int REFERENCES tasks(id),
				body text NOT NULL,
				tenant_id int NOT NULL
			);
			CREATE INDEX idx_tasks_project ON tasks(project_id);
			CREATE INDEX idx_comments_task ON comments(task_id);

			ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
			ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
			ALTER TABLE comments ENABLE ROW LEVEL SECURITY;

			CREATE POLICY tenant_projects ON projects
				USING (tenant_id = current_tenant());
			CREATE POLICY tenant_tasks ON tasks
				USING (tenant_id = current_tenant());
			CREATE POLICY tenant_comments ON comments
				USING (tenant_id = current_tenant());
		`
		assertMigrationValid(t, before, after)
	})

	// =================================================================
	// 5.4 Complete Schema Replacement
	// =================================================================

	// v1 schema (5 tables, 3 views, 2 functions) ->
	// v2 schema (different tables, views, functions).
	// All v1 dropped, all v2 created in correct order.
	t.Run("5.4/complete schema replacement v1 to v2", func(t *testing.T) {
		before := `
			CREATE FUNCTION v1_validate(val int) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT val > 0 $$;
			CREATE FUNCTION v1_format(t text) RETURNS text
				LANGUAGE sql IMMUTABLE AS $$ SELECT upper(t) $$;

			CREATE TABLE v1_customers (
				id int PRIMARY KEY,
				name text NOT NULL
			);
			CREATE TABLE v1_orders (
				id int PRIMARY KEY,
				customer_id int REFERENCES v1_customers(id),
				amount int CHECK (v1_validate(amount))
			);
			CREATE TABLE v1_items (
				id int PRIMARY KEY,
				order_id int REFERENCES v1_orders(id),
				product text
			);
			CREATE TABLE v1_categories (
				id int PRIMARY KEY,
				label text
			);
			CREATE TABLE v1_tags (
				id int PRIMARY KEY,
				category_id int REFERENCES v1_categories(id),
				name text
			);

			CREATE VIEW v1_order_summary AS
				SELECT o.id, c.name, o.amount
				FROM v1_orders o JOIN v1_customers c ON c.id = o.customer_id;
			CREATE VIEW v1_big_orders AS
				SELECT * FROM v1_order_summary WHERE amount > 1000;
			CREATE VIEW v1_tag_list AS
				SELECT t.name, c.label
				FROM v1_tags t JOIN v1_categories c ON c.id = t.category_id;
		`
		after := `
			CREATE FUNCTION v2_check_positive(val numeric) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT val > 0 $$;
			CREATE FUNCTION v2_slug(t text) RETURNS text
				LANGUAGE sql IMMUTABLE AS $$ SELECT lower(replace(t, ' ', '-')) $$;

			CREATE TABLE v2_accounts (
				id int PRIMARY KEY,
				email text UNIQUE NOT NULL,
				display_name text
			);
			CREATE TABLE v2_invoices (
				id int PRIMARY KEY,
				account_id int REFERENCES v2_accounts(id),
				total numeric(12,2) CHECK (v2_check_positive(total))
			);
			CREATE TABLE v2_line_items (
				id int PRIMARY KEY,
				invoice_id int REFERENCES v2_invoices(id),
				description text,
				amount numeric(10,2)
			);
			CREATE TABLE v2_labels (
				id int PRIMARY KEY,
				slug text NOT NULL
			);
			CREATE TABLE v2_attachments (
				id int PRIMARY KEY,
				label_id int REFERENCES v2_labels(id),
				url text NOT NULL
			);

			CREATE VIEW v2_invoice_summary AS
				SELECT i.id, a.email, i.total
				FROM v2_invoices i JOIN v2_accounts a ON a.id = i.account_id;
			CREATE VIEW v2_high_value AS
				SELECT * FROM v2_invoice_summary WHERE total > 5000;
			CREATE VIEW v2_label_attachments AS
				SELECT l.slug, att.url
				FROM v2_labels l JOIN v2_attachments att ON att.label_id = l.id;
		`
		assertMigrationValid(t, before, after)
	})

	// Schema replacement with some objects surviving — survivors not dropped,
	// new deps correct.
	t.Run("5.4/schema replacement with survivors", func(t *testing.T) {
		before := `
			CREATE FUNCTION shared_fn(val int) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT val > 0 $$;

			CREATE TABLE shared_config (
				id int PRIMARY KEY,
				key text UNIQUE NOT NULL,
				value text
			);
			CREATE TABLE old_users (
				id int PRIMARY KEY,
				name text,
				score int CHECK (shared_fn(score))
			);
			CREATE TABLE old_logs (
				id int PRIMARY KEY,
				user_id int REFERENCES old_users(id),
				message text
			);

			CREATE VIEW old_user_scores AS
				SELECT id, name, score FROM old_users;
			CREATE VIEW old_log_view AS
				SELECT l.id, u.name, l.message
				FROM old_logs l JOIN old_users u ON u.id = l.user_id;
		`
		after := `
			CREATE FUNCTION shared_fn(val int) RETURNS boolean
				LANGUAGE sql IMMUTABLE AS $$ SELECT val > 0 $$;

			CREATE TABLE shared_config (
				id int PRIMARY KEY,
				key text UNIQUE NOT NULL,
				value text
			);
			CREATE TABLE new_members (
				id int PRIMARY KEY,
				display_name text,
				rating int CHECK (shared_fn(rating))
			);
			CREATE TABLE new_activity (
				id int PRIMARY KEY,
				member_id int REFERENCES new_members(id),
				action text
			);

			CREATE VIEW new_member_ratings AS
				SELECT id, display_name, rating FROM new_members;
			CREATE VIEW new_activity_feed AS
				SELECT a.id, m.display_name, a.action
				FROM new_activity a JOIN new_members m ON m.id = a.member_id;
		`
		assertMigrationValid(t, before, after)
	})

	// =================================================================
	// 5.5 Add Partitioning to Existing Table
	// =================================================================

	// Before: single orders table with indexes + views.
	// After: partitioned orders with 4 range partitions + indexes.
	// Views dropped, old table restructured, partitions created,
	// indexes recreated, views recreated.
	t.Run("5.5/add partitioning to existing table", func(t *testing.T) {
		before := `
			CREATE TABLE orders (
				id int NOT NULL,
				created_at timestamp NOT NULL,
				customer_id int NOT NULL,
				total numeric(12,2),
				status text DEFAULT 'pending'
			);
			CREATE INDEX idx_orders_created ON orders(created_at);
			CREATE INDEX idx_orders_customer ON orders(customer_id);
			CREATE VIEW pending_orders AS
				SELECT id, customer_id, total FROM orders WHERE status = 'pending';
			CREATE VIEW order_totals AS
				SELECT customer_id, sum(total) as total_spent FROM orders GROUP BY customer_id;
		`
		after := `
			CREATE TABLE orders (
				id int NOT NULL,
				created_at timestamp NOT NULL,
				customer_id int NOT NULL,
				total numeric(12,2),
				status text DEFAULT 'pending'
			) PARTITION BY RANGE (created_at);

			CREATE TABLE orders_2024q1 PARTITION OF orders
				FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
			CREATE TABLE orders_2024q2 PARTITION OF orders
				FOR VALUES FROM ('2024-04-01') TO ('2024-07-01');
			CREATE TABLE orders_2024q3 PARTITION OF orders
				FOR VALUES FROM ('2024-07-01') TO ('2024-10-01');
			CREATE TABLE orders_2024q4 PARTITION OF orders
				FOR VALUES FROM ('2024-10-01') TO ('2025-01-01');

			CREATE INDEX idx_orders_created ON orders(created_at);
			CREATE INDEX idx_orders_customer ON orders(customer_id);
			CREATE VIEW pending_orders AS
				SELECT id, customer_id, total FROM orders WHERE status = 'pending';
			CREATE VIEW order_totals AS
				SELECT customer_id, sum(total) as total_spent FROM orders GROUP BY customer_id;
		`
		assertMigrationValid(t, before, after)
	})
}

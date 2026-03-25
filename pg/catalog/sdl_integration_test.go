package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Section 3.4: Complex Multi-Object SDL
// ---------------------------------------------------------------------------

func TestSDLComplexMultiObject(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "full schema with 10+ object types loaded from shuffled SDL",
			sql: `
				-- Shuffled: triggers, views, indexes, grants, comments, functions,
				-- types, sequences, tables, schemas, policies — all out of order.
				CREATE TRIGGER audit_trg BEFORE INSERT ON orders EXECUTE FUNCTION audit_fn();
				CREATE VIEW order_summary AS SELECT o.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id;
				CREATE INDEX idx_orders_cust ON orders (customer_id);
				COMMENT ON TABLE orders IS 'order records';
				GRANT SELECT ON customers TO PUBLIC;
				CREATE FUNCTION audit_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
				CREATE TYPE order_status AS ENUM ('pending', 'shipped', 'delivered');
				CREATE SEQUENCE order_seq;
				CREATE TABLE orders (
					id int PRIMARY KEY DEFAULT nextval('order_seq'),
					customer_id int REFERENCES customers(id),
					status order_status
				);
				CREATE TABLE customers (id int PRIMARY KEY, name text);
				CREATE SCHEMA reports;
			`,
			check: func(t *testing.T, c *Catalog) {
				for _, name := range []string{"orders", "customers"} {
					if c.GetRelation("public", name) == nil {
						t.Fatalf("table %s not found", name)
					}
				}
				if c.GetRelation("public", "order_summary") == nil {
					t.Fatal("view order_summary not found")
				}
			},
		},
		{
			name: "function used in CHECK, view, and trigger — single function resolved for all",
			sql: `
				CREATE TRIGGER trg BEFORE INSERT ON t EXECUTE FUNCTION audit_fn();
				CREATE VIEW v AS SELECT validate() AS ok;
				CREATE TABLE t (
					id int,
					val int CHECK (validate())
				);
				CREATE FUNCTION validate() RETURNS boolean AS $$ BEGIN RETURN true; END; $$ LANGUAGE plpgsql;
				CREATE FUNCTION audit_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
			},
		},
		{
			name: "table with SERIAL column — implicit sequence created",
			sql: `
				CREATE TABLE t (id SERIAL PRIMARY KEY, name text);
			`,
			check: func(t *testing.T, c *Catalog) {
				rel := c.GetRelation("public", "t")
				if rel == nil {
					t.Fatal("table t not found")
				}
				// SERIAL creates an implicit sequence — verify column exists.
				found := false
				for _, col := range rel.Columns {
					if col.Name == "id" {
						found = true
						break
					}
				}
				if !found {
					t.Fatal("column id not found in table t")
				}
			},
		},
		{
			name: "materialized view with indexes — matview created before its indexes",
			sql: `
				CREATE INDEX idx_mv ON mv (id);
				CREATE MATERIALIZED VIEW mv AS SELECT id, name FROM t;
				CREATE TABLE t (id int, name text);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
				if c.GetRelation("public", "mv") == nil {
					t.Fatal("materialized view mv not found")
				}
			},
		},
		{
			name: "multiple schemas with cross-schema references",
			sql: `
				CREATE VIEW app.user_view AS SELECT id, email FROM app.users;
				CREATE TABLE app.users (id int PRIMARY KEY, email text);
				CREATE TABLE app.orders (id int PRIMARY KEY, user_id int REFERENCES app.users(id));
				CREATE SCHEMA app;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("app", "users") == nil {
					t.Fatal("table app.users not found")
				}
				if c.GetRelation("app", "orders") == nil {
					t.Fatal("table app.orders not found")
				}
				if c.GetRelation("app", "user_view") == nil {
					t.Fatal("view app.user_view not found")
				}
			},
		},
		{
			name: "SDL producing identical catalog to LoadSQL with same DDL in correct order",
			sql: `
				CREATE INDEX idx_orders_customer ON orders (customer_id);
				COMMENT ON TABLE orders IS 'all orders';
				GRANT SELECT ON customers TO PUBLIC;
				CREATE TRIGGER trg BEFORE INSERT ON orders EXECUTE FUNCTION audit_fn();
				CREATE VIEW order_view AS SELECT o.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id;
				CREATE FUNCTION audit_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
				CREATE TYPE status AS ENUM ('active', 'inactive');
				CREATE SEQUENCE order_seq;
				CREATE TABLE orders (
					id int PRIMARY KEY DEFAULT nextval('order_seq'),
					customer_id int REFERENCES customers(id),
					s status
				);
				CREATE TABLE customers (id int PRIMARY KEY, name text);
			`,
			check: func(t *testing.T, sdlCatalog *Catalog) {
				// Equivalent DDL in correct dependency order.
				correctOrder := `
					CREATE TYPE status AS ENUM ('active', 'inactive');
					CREATE SEQUENCE order_seq;
					CREATE TABLE customers (id int PRIMARY KEY, name text);
					CREATE TABLE orders (
						id int PRIMARY KEY DEFAULT nextval('order_seq'),
						customer_id int REFERENCES customers(id),
						s status
					);
					CREATE FUNCTION audit_fn() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
					CREATE VIEW order_view AS SELECT o.id, c.name FROM orders o JOIN customers c ON o.customer_id = c.id;
					CREATE INDEX idx_orders_customer ON orders (customer_id);
					CREATE TRIGGER trg BEFORE INSERT ON orders EXECUTE FUNCTION audit_fn();
					COMMENT ON TABLE orders IS 'all orders';
					GRANT SELECT ON customers TO PUBLIC;
				`
				sqlCatalog, err := LoadSQL(correctOrder)
				if err != nil {
					t.Fatalf("LoadSQL failed: %v", err)
				}
				diff := Diff(sdlCatalog, sqlCatalog)
				if !diff.IsEmpty() {
					var parts []string
					for _, r := range diff.Relations {
						parts = append(parts, fmt.Sprintf("relation %s.%s: action=%d", r.SchemaName, r.Name, r.Action))
					}
					for _, e := range diff.Enums {
						parts = append(parts, fmt.Sprintf("enum %s.%s: action=%d", e.SchemaName, e.Name, e.Action))
					}
					for _, s := range diff.Sequences {
						parts = append(parts, fmt.Sprintf("sequence %s.%s: action=%d", s.SchemaName, s.Name, s.Action))
					}
					for _, f := range diff.Functions {
						parts = append(parts, fmt.Sprintf("function %s: action=%d", f.Identity, f.Action))
					}
					for _, c := range diff.Comments {
						parts = append(parts, fmt.Sprintf("comment %s: action=%d", c.ObjDescription, c.Action))
					}
					for _, g := range diff.Grants {
						parts = append(parts, fmt.Sprintf("grant: action=%d", g.Action))
					}
					t.Errorf("catalogs differ:\n%s", strings.Join(parts, "\n"))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := LoadSDL(tt.sql)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

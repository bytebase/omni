package catalog

import (
	"fmt"
	"strings"
	"testing"
)

// TestSDLResolve tests section 1.2: basic dependency resolution.
func TestSDLResolve(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "two tables with no dependencies",
			sql: `
				CREATE TABLE a (id int);
				CREATE TABLE b (id int);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "a") == nil {
					t.Fatal("table a not found")
				}
				if c.GetRelation("public", "b") == nil {
					t.Fatal("table b not found")
				}
			},
		},
		{
			name: "FK forward reference - FK table defined first",
			sql: `
				CREATE TABLE orders (
					id int PRIMARY KEY,
					customer_id int REFERENCES customers(id)
				);
				CREATE TABLE customers (id int PRIMARY KEY);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "orders") == nil {
					t.Fatal("table orders not found")
				}
				if c.GetRelation("public", "customers") == nil {
					t.Fatal("table customers not found")
				}
			},
		},
		{
			name: "FK forward reference - FK table defined last",
			sql: `
				CREATE TABLE customers (id int PRIMARY KEY);
				CREATE TABLE orders (
					id int PRIMARY KEY,
					customer_id int REFERENCES customers(id)
				);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "orders") == nil {
					t.Fatal("table orders not found")
				}
				if c.GetRelation("public", "customers") == nil {
					t.Fatal("table customers not found")
				}
			},
		},
		{
			name: "enum type forward reference",
			sql: `
				CREATE TABLE t (id int, status mood);
				CREATE TYPE mood AS ENUM ('happy', 'sad');
			`,
			check: func(t *testing.T, c *Catalog) {
				rel := c.GetRelation("public", "t")
				if rel == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "domain type forward reference",
			sql: `
				CREATE TABLE t (id int, age posint);
				CREATE DOMAIN posint AS integer CHECK (VALUE > 0);
			`,
			check: func(t *testing.T, c *Catalog) {
				rel := c.GetRelation("public", "t")
				if rel == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "composite type forward reference",
			sql: `
				CREATE TABLE t (id int, loc point2d);
				CREATE TYPE point2d AS (x int, y int);
			`,
			check: func(t *testing.T, c *Catalog) {
				rel := c.GetRelation("public", "t")
				if rel == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "view forward reference - view before table",
			sql: `
				CREATE VIEW v AS SELECT id FROM t;
				CREATE TABLE t (id int);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v") == nil {
					t.Fatal("view v not found")
				}
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "index forward reference - index before table",
			sql: `
				CREATE INDEX idx ON t (id);
				CREATE TABLE t (id int);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "result catalog matches LoadSQL with correctly-ordered equivalent",
			sql: `
				CREATE TABLE orders (
					id int PRIMARY KEY,
					customer_id int REFERENCES customers(id)
				);
				CREATE TYPE mood AS ENUM ('happy', 'sad');
				CREATE TABLE customers (id int PRIMARY KEY, feeling mood);
				CREATE VIEW customer_view AS SELECT id FROM customers;
				CREATE INDEX idx_customers_id ON customers (id);
			`,
			check: func(t *testing.T, sdlCatalog *Catalog) {
				// LoadSQL with correctly ordered DDL.
				correctOrder := `
					CREATE TYPE mood AS ENUM ('happy', 'sad');
					CREATE TABLE customers (id int PRIMARY KEY, feeling mood);
					CREATE TABLE orders (
						id int PRIMARY KEY,
						customer_id int REFERENCES customers(id)
					);
					CREATE VIEW customer_view AS SELECT id FROM customers;
					CREATE INDEX idx_customers_id ON customers (id);
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

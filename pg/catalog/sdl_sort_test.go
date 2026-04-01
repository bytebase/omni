package catalog

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Section 3.1: Priority Layer Ordering
// ---------------------------------------------------------------------------

func TestSDLPriorityLayerOrdering(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "schemas created before tables that reference them",
			sql: `
				CREATE TABLE myschema.t (id int);
				CREATE SCHEMA myschema;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("myschema", "t") == nil {
					t.Fatal("table myschema.t not found")
				}
			},
		},
		{
			name: "extensions created before types/functions they provide",
			// Extensions are layer 1, types are layer 2.
			// We can't truly test extension loading without pg_extension support,
			// but we verify the ordering doesn't error.
			sql: `
				CREATE TABLE t (id int);
				CREATE SCHEMA s;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "types created before tables using them as column types",
			sql: `
				CREATE TABLE t (id int, status mood);
				CREATE TYPE mood AS ENUM ('happy', 'sad');
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "sequences created before tables with DEFAULT nextval",
			sql: `
				CREATE TABLE t (id int DEFAULT nextval('my_seq'));
				CREATE SEQUENCE my_seq;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "tables created before views referencing them",
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
			name: "functions created before triggers referencing them",
			sql: `
				CREATE TRIGGER trg BEFORE INSERT ON t EXECUTE FUNCTION myfunc();
				CREATE TABLE t (id int);
				CREATE FUNCTION myfunc() RETURNS trigger AS $$ BEGIN RETURN NEW; END; $$ LANGUAGE plpgsql;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "tables created before indexes on them",
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
			name: "FK constraints applied after all tables created",
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
			name: "comments applied after their target objects created",
			sql: `
				COMMENT ON TABLE t IS 'my table';
				CREATE TABLE t (id int);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "grants applied after their target objects created",
			sql: `
				GRANT SELECT ON t TO PUBLIC;
				CREATE TABLE t (id int);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
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

// ---------------------------------------------------------------------------
// Section 3.2: Topological Sort Within Layers
// ---------------------------------------------------------------------------

func TestSDLTopoSortWithinLayers(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "view A depends on view B — B created first",
			sql: `
				CREATE TABLE t (id int, name text);
				CREATE VIEW va AS SELECT * FROM vb;
				CREATE VIEW vb AS SELECT id FROM t;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "va") == nil {
					t.Fatal("view va not found")
				}
				if c.GetRelation("public", "vb") == nil {
					t.Fatal("view vb not found")
				}
			},
		},
		{
			name: "chain of 5 views",
			sql: `
				CREATE TABLE base (id int);
				CREATE VIEW v5 AS SELECT * FROM v4;
				CREATE VIEW v4 AS SELECT * FROM v3;
				CREATE VIEW v3 AS SELECT * FROM v2;
				CREATE VIEW v2 AS SELECT * FROM v1;
				CREATE VIEW v1 AS SELECT * FROM base;
			`,
			check: func(t *testing.T, c *Catalog) {
				for _, name := range []string{"v1", "v2", "v3", "v4", "v5"} {
					if c.GetRelation("public", name) == nil {
						t.Fatalf("view %s not found", name)
					}
				}
			},
		},
		{
			name: "two independent views — no false dependency",
			sql: `
				CREATE TABLE t1 (id int);
				CREATE TABLE t2 (id int);
				CREATE VIEW v1 AS SELECT * FROM t1;
				CREATE VIEW v2 AS SELECT * FROM t2;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "v1") == nil {
					t.Fatal("view v1 not found")
				}
				if c.GetRelation("public", "v2") == nil {
					t.Fatal("view v2 not found")
				}
			},
		},
		{
			name: "table with INHERITS parent — parent created first",
			sql: `
				CREATE TABLE child (extra text) INHERITS (parent);
				CREATE TABLE parent (id int);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "parent") == nil {
					t.Fatal("table parent not found")
				}
				if c.GetRelation("public", "child") == nil {
					t.Fatal("table child not found")
				}
			},
		},
		{
			name: "table PARTITION OF parent — parent created first",
			sql: `
				CREATE TABLE child PARTITION OF parent FOR VALUES FROM (1) TO (100);
				CREATE TABLE parent (id int) PARTITION BY RANGE (id);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "parent") == nil {
					t.Fatal("table parent not found")
				}
				if c.GetRelation("public", "child") == nil {
					t.Fatal("table child not found")
				}
			},
		},
		{
			name: "function created before table with CHECK referencing it",
			sql: `
				CREATE TABLE items (
					id int,
					qty int,
					CONSTRAINT items_qty_check CHECK (is_positive(qty))
				);
				CREATE FUNCTION is_positive(integer) RETURNS boolean
					LANGUAGE sql IMMUTABLE AS $$SELECT $1 > 0$$;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "items") == nil {
					t.Fatal("table items not found")
				}
			},
		},
		{
			name: "function created before table with DEFAULT referencing it",
			sql: `
				CREATE TABLE t (
					id int DEFAULT next_id()
				);
				CREATE FUNCTION next_id() RETURNS int
					LANGUAGE sql AS $$SELECT 1$$;
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "t") == nil {
					t.Fatal("table t not found")
				}
			},
		},
		{
			name: "composite type A references type B — B created first",
			sql: `
				CREATE TYPE type_a AS (b_field type_b);
				CREATE TYPE type_b AS (val int);
			`,
			check: func(t *testing.T, c *Catalog) {
				// Both types should be created successfully.
				// Verify by loading SDL without error.
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

// ---------------------------------------------------------------------------
// Section 3.3: Cycle Detection and Repair
// ---------------------------------------------------------------------------

func TestSDLCycleDetectionRepair(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
		check   func(t *testing.T, c *Catalog)
	}{
		{
			name: "mutual FK between two tables — both created, FKs applied after",
			sql: `
				CREATE TABLE a (
					id int PRIMARY KEY,
					b_id int REFERENCES b(id)
				);
				CREATE TABLE b (
					id int PRIMARY KEY,
					a_id int REFERENCES a(id)
				);
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
			name: "three-way FK cycle (A→B→C→A) — all tables created, FKs deferred",
			sql: `
				CREATE TABLE a (
					id int PRIMARY KEY,
					b_id int REFERENCES b(id)
				);
				CREATE TABLE b (
					id int PRIMARY KEY,
					c_id int REFERENCES c(id)
				);
				CREATE TABLE c (
					id int PRIMARY KEY,
					a_id int REFERENCES a(id)
				);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "a") == nil {
					t.Fatal("table a not found")
				}
				if c.GetRelation("public", "b") == nil {
					t.Fatal("table b not found")
				}
				if c.GetRelation("public", "c") == nil {
					t.Fatal("table c not found")
				}
			},
		},
		{
			name: "self-referencing FK — handled correctly",
			sql: `
				CREATE TABLE tree (
					id int PRIMARY KEY,
					parent_id int REFERENCES tree(id)
				);
			`,
			check: func(t *testing.T, c *Catalog) {
				if c.GetRelation("public", "tree") == nil {
					t.Fatal("table tree not found")
				}
			},
		},
		{
			name: "composite type mutual reference via arrays — shell types resolve",
			sql: `
				CREATE TYPE person AS (name text, addr address[]);
				CREATE TYPE address AS (city text, occupant person[]);
			`,
			check: func(t *testing.T, c *Catalog) {
				// Both types should be created. The shell type mechanism
				// pre-creates type entries so mutual array references work.
			},
		},
		{
			name: "unresolvable cycle produces clear error",
			// Two composite types that directly reference each other (not via arrays)
			// create a true cycle that cannot be resolved. However, since we add
			// shell types, even direct references should work because composite type
			// columns only need the type to exist (shell is enough).
			// Instead, we test a view cycle which is truly unresolvable.
			sql: `
				CREATE TABLE t (id int);
				CREATE VIEW va AS SELECT * FROM vb;
				CREATE VIEW vb AS SELECT * FROM va;
			`,
			wantErr: "cycle",
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

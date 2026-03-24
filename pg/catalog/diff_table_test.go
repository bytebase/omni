package catalog

import (
	"testing"
)

func TestDiffTable(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "table added",
			fromSQL: "",
			toSQL:   "CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "t1" {
					t.Fatalf("expected name t1, got %s", e.Name)
				}
				if e.SchemaName != "public" {
					t.Fatalf("expected schema public, got %s", e.SchemaName)
				}
				if e.From != nil {
					t.Fatalf("expected From to be nil")
				}
				if e.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
			},
		},
		{
			name:    "table dropped",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "t1" {
					t.Fatalf("expected name t1, got %s", e.Name)
				}
				if e.From == nil {
					t.Fatalf("expected From to be non-nil")
				}
				if e.To != nil {
					t.Fatalf("expected To to be nil")
				}
			},
		},
		{
			name:    "table unchanged",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 0 {
					t.Fatalf("expected 0 relation diffs, got %d", len(d.Relations))
				}
			},
		},
		{
			name:    "multiple tables across different schemas",
			fromSQL: "CREATE SCHEMA s1; CREATE TABLE s1.a (id int);",
			toSQL:   "CREATE SCHEMA s1; CREATE TABLE s1.a (id int); CREATE TABLE s1.b (id int); CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 2 {
					t.Fatalf("expected 2 relation diffs, got %d", len(d.Relations))
				}
				// Sorted: public.t1 and s1.b (both adds)
				for _, e := range d.Relations {
					if e.Action != DiffAdd {
						t.Fatalf("expected DiffAdd for %s.%s, got %d", e.SchemaName, e.Name, e.Action)
					}
				}
			},
		},
		{
			name:    "table persistence changed",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE UNLOGGED TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From == nil || e.To == nil {
					t.Fatalf("expected both From and To to be non-nil")
				}
				if e.From.Persistence == e.To.Persistence {
					t.Fatalf("expected persistence to differ")
				}
			},
		},
		{
			name:    "table in non-public schema added",
			fromSQL: "CREATE SCHEMA sales;",
			toSQL:   "CREATE SCHEMA sales; CREATE TABLE sales.orders (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.SchemaName != "sales" {
					t.Fatalf("expected schema sales, got %s", e.SchemaName)
				}
				if e.Name != "orders" {
					t.Fatalf("expected name orders, got %s", e.Name)
				}
			},
		},
		{
			name:    "table in non-public schema dropped",
			fromSQL: "CREATE SCHEMA sales; CREATE TABLE sales.orders (id int);",
			toSQL:   "CREATE SCHEMA sales;",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.SchemaName != "sales" {
					t.Fatalf("expected schema sales, got %s", e.SchemaName)
				}
			},
		},
		{
			name: "same name in different schemas treated as distinct",
			fromSQL: `
				CREATE SCHEMA s1;
				CREATE SCHEMA s2;
				CREATE TABLE s1.t1 (id int);
			`,
			toSQL: `
				CREATE SCHEMA s1;
				CREATE SCHEMA s2;
				CREATE TABLE s2.t1 (id int);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 2 {
					t.Fatalf("expected 2 relation diffs, got %d", len(d.Relations))
				}
				var hasAdd, hasDrop bool
				for _, e := range d.Relations {
					switch e.Action {
					case DiffAdd:
						hasAdd = true
						if e.SchemaName != "s2" {
							t.Fatalf("expected add in s2, got %s", e.SchemaName)
						}
					case DiffDrop:
						hasDrop = true
						if e.SchemaName != "s1" {
							t.Fatalf("expected drop in s1, got %s", e.SchemaName)
						}
					default:
						t.Fatalf("unexpected action %d", e.Action)
					}
				}
				if !hasAdd || !hasDrop {
					t.Fatalf("expected one add and one drop")
				}
			},
		},
		{
			name:    "table renamed detected as drop old add new",
			fromSQL: "CREATE TABLE old_name (id int);",
			toSQL:   "CREATE TABLE new_name (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 2 {
					t.Fatalf("expected 2 relation diffs (drop+add), got %d", len(d.Relations))
				}
				var hasAdd, hasDrop bool
				for _, e := range d.Relations {
					switch e.Action {
					case DiffAdd:
						hasAdd = true
						if e.Name != "new_name" {
							t.Fatalf("expected add new_name, got %s", e.Name)
						}
					case DiffDrop:
						hasDrop = true
						if e.Name != "old_name" {
							t.Fatalf("expected drop old_name, got %s", e.Name)
						}
					}
				}
				if !hasAdd || !hasDrop {
					t.Fatalf("expected one add and one drop")
				}
			},
		},
		{
			name:    "table ReplicaIdentity changed",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int); ALTER TABLE t1 REPLICA IDENTITY FULL;",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.ReplicaIdentity == e.To.ReplicaIdentity {
					t.Fatalf("expected ReplicaIdentity to differ")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, err := LoadSQL(tt.fromSQL)
			if err != nil {
				t.Fatal(err)
			}
			to, err := LoadSQL(tt.toSQL)
			if err != nil {
				t.Fatal(err)
			}
			d := Diff(from, to)
			tt.check(t, d)
		})
	}
}

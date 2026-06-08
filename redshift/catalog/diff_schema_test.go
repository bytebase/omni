package catalog

import (
	"testing"
)

func TestDiffSchema(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "types compile",
			fromSQL: "",
			toSQL:   "",
			check: func(t *testing.T, d *SchemaDiff) {
				// SchemaDiff struct and all diff entry types compile with no errors.
				// This test passing means compilation succeeded.
				_ = SchemaDiffEntry{}
				_ = RelationDiffEntry{}
				_ = ColumnDiffEntry{}
				_ = ConstraintDiffEntry{}
				_ = IndexDiffEntry{}
				_ = SequenceDiffEntry{}
				_ = FunctionDiffEntry{}
				_ = TriggerDiffEntry{}
				_ = EnumDiffEntry{}
				_ = DomainDiffEntry{}
				_ = RangeDiffEntry{}
				_ = ExtensionDiffEntry{}
				_ = PolicyDiffEntry{}
				_ = CommentDiffEntry{}
				_ = GrantDiffEntry{}
			},
		},
		{
			name:    "identical empty catalogs",
			fromSQL: "",
			toSQL:   "",
			check: func(t *testing.T, d *SchemaDiff) {
				if !d.IsEmpty() {
					t.Fatalf("expected empty diff, got %d schema entries", len(d.Schemas))
				}
			},
		},
		{
			name:    "identical non-empty catalogs",
			fromSQL: "CREATE SCHEMA sales; CREATE TABLE sales.orders (id int);",
			toSQL:   "CREATE SCHEMA sales; CREATE TABLE sales.orders (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Schemas) != 0 {
					t.Fatalf("expected 0 schema diffs, got %d", len(d.Schemas))
				}
			},
		},
		{
			name:    "schema added",
			fromSQL: "",
			toSQL:   "CREATE SCHEMA analytics;",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Schemas) != 1 {
					t.Fatalf("expected 1 schema diff, got %d", len(d.Schemas))
				}
				e := d.Schemas[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "analytics" {
					t.Fatalf("expected name 'analytics', got %q", e.Name)
				}
				if e.From != nil {
					t.Fatal("expected From to be nil for add")
				}
				if e.To == nil {
					t.Fatal("expected To to be non-nil for add")
				}
			},
		},
		{
			name:    "schema dropped",
			fromSQL: "CREATE SCHEMA analytics;",
			toSQL:   "",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Schemas) != 1 {
					t.Fatalf("expected 1 schema diff, got %d", len(d.Schemas))
				}
				e := d.Schemas[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "analytics" {
					t.Fatalf("expected name 'analytics', got %q", e.Name)
				}
				if e.From == nil {
					t.Fatal("expected From to be non-nil for drop")
				}
				if e.To != nil {
					t.Fatal("expected To to be nil for drop")
				}
			},
		},
		{
			name:    "schema modified owner changed",
			fromSQL: "CREATE SCHEMA sales AUTHORIZATION alice;",
			toSQL:   "CREATE SCHEMA sales AUTHORIZATION bob;",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Schemas) != 1 {
					t.Fatalf("expected 1 schema diff, got %d", len(d.Schemas))
				}
				e := d.Schemas[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "sales" {
					t.Fatalf("expected name 'sales', got %q", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil for modify")
				}
				if e.From.Owner != "alice" {
					t.Fatalf("expected From.Owner 'alice', got %q", e.From.Owner)
				}
				if e.To.Owner != "bob" {
					t.Fatalf("expected To.Owner 'bob', got %q", e.To.Owner)
				}
			},
		},
		{
			name:    "multiple schemas added dropped modified",
			fromSQL: "CREATE SCHEMA alpha AUTHORIZATION alice; CREATE SCHEMA beta; CREATE SCHEMA gamma;",
			toSQL:   "CREATE SCHEMA alpha AUTHORIZATION bob; CREATE SCHEMA delta; CREATE SCHEMA gamma;",
			check: func(t *testing.T, d *SchemaDiff) {
				// alpha: modified (owner changed)
				// beta: dropped
				// delta: added
				// gamma: unchanged
				if len(d.Schemas) != 3 {
					t.Fatalf("expected 3 schema diffs, got %d", len(d.Schemas))
				}

				// Sorted by name: alpha, beta, delta
				e0 := d.Schemas[0]
				if e0.Name != "alpha" || e0.Action != DiffModify {
					t.Fatalf("expected alpha/DiffModify, got %s/%d", e0.Name, e0.Action)
				}
				e1 := d.Schemas[1]
				if e1.Name != "beta" || e1.Action != DiffDrop {
					t.Fatalf("expected beta/DiffDrop, got %s/%d", e1.Name, e1.Action)
				}
				e2 := d.Schemas[2]
				if e2.Name != "delta" || e2.Action != DiffAdd {
					t.Fatalf("expected delta/DiffAdd, got %s/%d", e2.Name, e2.Action)
				}
			},
		},
		{
			name:    "pg_catalog and pg_toast excluded",
			fromSQL: "",
			toSQL:   "",
			check: func(t *testing.T, d *SchemaDiff) {
				// Both catalogs have pg_catalog and pg_toast built-in.
				// They should never appear in diff results.
				for _, e := range d.Schemas {
					if e.Name == "pg_catalog" || e.Name == "pg_toast" {
						t.Fatalf("system schema %q should be excluded from diff", e.Name)
					}
				}
				if !d.IsEmpty() {
					t.Fatalf("expected empty diff for two default catalogs, got %d schema entries", len(d.Schemas))
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

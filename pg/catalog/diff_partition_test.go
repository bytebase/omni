package catalog

import (
	"testing"
)

func TestDiffPartition(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "partitioned table added with PARTITION BY",
			fromSQL: "",
			toSQL:   "CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "orders" {
					t.Fatalf("expected name orders, got %s", e.Name)
				}
				if e.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
				if e.To.RelKind != 'p' {
					t.Fatalf("expected RelKind 'p' (partitioned), got %q", e.To.RelKind)
				}
				if e.To.PartitionInfo == nil {
					t.Fatalf("expected PartitionInfo to be non-nil")
				}
				if e.To.PartitionInfo.Strategy != 'l' {
					t.Fatalf("expected list strategy 'l', got %q", e.To.PartitionInfo.Strategy)
				}
			},
		},
		{
			name: "partition child added (PARTITION OF)",
			fromSQL: `
				CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
			`,
			toSQL: `
				CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
				CREATE TABLE orders_us PARTITION OF orders FOR VALUES IN ('US');
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff (child added), got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "orders_us" {
					t.Fatalf("expected name orders_us, got %s", e.Name)
				}
				if e.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
				if e.To.PartitionBound == nil {
					t.Fatalf("expected PartitionBound to be non-nil")
				}
				if e.To.PartitionOf == 0 {
					t.Fatalf("expected PartitionOf to be non-zero (partition child)")
				}
			},
		},
		{
			name: "partition strategy detected in diff (list/range/hash)",
			fromSQL: `
				CREATE TABLE t_list (id int, region text) PARTITION BY LIST (region);
			`,
			toSQL: `
				CREATE TABLE t_list (id int, region text) PARTITION BY LIST (region);
				CREATE TABLE t_range (id int, created_at date) PARTITION BY RANGE (created_at);
				CREATE TABLE t_hash (id int, val text) PARTITION BY HASH (id);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 2 {
					t.Fatalf("expected 2 relation diffs (range+hash added), got %d", len(d.Relations))
				}
				strategies := make(map[string]byte)
				for _, e := range d.Relations {
					if e.Action != DiffAdd {
						t.Fatalf("expected DiffAdd for %s, got %d", e.Name, e.Action)
					}
					if e.To.PartitionInfo == nil {
						t.Fatalf("expected PartitionInfo non-nil for %s", e.Name)
					}
					strategies[e.Name] = e.To.PartitionInfo.Strategy
				}
				if s, ok := strategies["t_range"]; !ok || s != 'r' {
					t.Fatalf("expected t_range strategy 'r', got %q", s)
				}
				if s, ok := strategies["t_hash"]; !ok || s != 'h' {
					t.Fatalf("expected t_hash strategy 'h', got %q", s)
				}
			},
		},
		{
			name: "partition bound values captured in diff",
			fromSQL: `
				CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
			`,
			toSQL: `
				CREATE TABLE orders (id int, region text) PARTITION BY LIST (region);
				CREATE TABLE orders_us PARTITION OF orders FOR VALUES IN ('US', 'CA');
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
				pb := e.To.PartitionBound
				if pb == nil {
					t.Fatalf("expected PartitionBound to be non-nil")
				}
				if pb.Strategy != 'l' {
					t.Fatalf("expected list strategy 'l', got %q", pb.Strategy)
				}
				if len(pb.ListValues) != 2 {
					t.Fatalf("expected 2 list values, got %d: %v", len(pb.ListValues), pb.ListValues)
				}
			},
		},
		{
			name: "table inheritance changed (INHERITS list modified)",
			fromSQL: `
				CREATE TABLE parent1 (id int);
				CREATE TABLE parent2 (id int);
				CREATE TABLE child (id int) INHERITS (parent1);
			`,
			toSQL: `
				CREATE TABLE parent1 (id int);
				CREATE TABLE parent2 (id int);
				CREATE TABLE child (id int) INHERITS (parent2);
			`,
			check: func(t *testing.T, d *SchemaDiff) {
				// child table should show as modified because InhParents changed
				var found bool
				for _, e := range d.Relations {
					if e.Name == "child" {
						found = true
						if e.Action != DiffModify {
							t.Fatalf("expected DiffModify for child, got %d", e.Action)
						}
						if e.From == nil || e.To == nil {
							t.Fatalf("expected both From and To to be non-nil")
						}
						if len(e.From.InhParents) != 1 {
							t.Fatalf("expected 1 InhParent in From, got %d", len(e.From.InhParents))
						}
						if len(e.To.InhParents) != 1 {
							t.Fatalf("expected 1 InhParent in To, got %d", len(e.To.InhParents))
						}
						// The OIDs will differ but they should resolve to different parent names
						break
					}
				}
				if !found {
					t.Fatalf("expected child table in diff results")
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

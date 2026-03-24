package catalog

import (
	"testing"
)

func TestDiffRange(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []RangeDiffEntry)
	}{
		{
			name:    "range type added",
			fromSQL: "",
			toSQL:   "CREATE TYPE floatrange AS RANGE (subtype = float8);",
			check: func(t *testing.T, entries []RangeDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "floatrange" {
					t.Fatalf("expected name 'floatrange', got %q", e.Name)
				}
				if e.SchemaName != "public" {
					t.Fatalf("expected schema 'public', got %q", e.SchemaName)
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
			name:    "range type dropped",
			fromSQL: "CREATE TYPE floatrange AS RANGE (subtype = float8);",
			toSQL:   "",
			check: func(t *testing.T, entries []RangeDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "floatrange" {
					t.Fatalf("expected name 'floatrange', got %q", e.Name)
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
			name:    "range subtype changed",
			fromSQL: "CREATE TYPE myrange AS RANGE (subtype = int4);",
			toSQL:   "CREATE TYPE myrange AS RANGE (subtype = int8);",
			check: func(t *testing.T, entries []RangeDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "myrange" {
					t.Fatalf("expected name 'myrange', got %q", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil for modify")
				}
			},
		},
		{
			name: "range type identity by schema and name",
			fromSQL: `CREATE SCHEMA s1;
				CREATE TYPE s1.myrange AS RANGE (subtype = int4);
				CREATE TYPE myrange AS RANGE (subtype = int4);`,
			toSQL: `CREATE SCHEMA s1;
				CREATE TYPE s1.myrange AS RANGE (subtype = int8);
				CREATE TYPE myrange AS RANGE (subtype = int4);`,
			check: func(t *testing.T, entries []RangeDiffEntry) {
				// Only the s1.myrange should show as modified; public.myrange is unchanged.
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.SchemaName != "s1" {
					t.Fatalf("expected schema 's1', got %q", e.SchemaName)
				}
				if e.Name != "myrange" {
					t.Fatalf("expected name 'myrange', got %q", e.Name)
				}
			},
		},
		{
			name:    "range unchanged produces no entry",
			fromSQL: "CREATE TYPE floatrange AS RANGE (subtype = float8);",
			toSQL:   "CREATE TYPE floatrange AS RANGE (subtype = float8);",
			check: func(t *testing.T, entries []RangeDiffEntry) {
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries, got %d", len(entries))
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
			entries := diffRanges(from, to)
			tt.check(t, entries)
		})
	}
}

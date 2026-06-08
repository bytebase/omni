package catalog

import (
	"testing"
)

func TestDiffDomain(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []DomainDiffEntry)
	}{
		{
			name:    "domain added",
			fromSQL: "",
			toSQL:   "CREATE DOMAIN posint AS integer CHECK (VALUE > 0);",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "posint" {
					t.Fatalf("expected name 'posint', got %q", e.Name)
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
			name:    "domain dropped",
			fromSQL: "CREATE DOMAIN posint AS integer CHECK (VALUE > 0);",
			toSQL:   "",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "posint" {
					t.Fatalf("expected name 'posint', got %q", e.Name)
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
			name:    "domain base type changed",
			fromSQL: "CREATE DOMAIN mydom AS integer;",
			toSQL:   "CREATE DOMAIN mydom AS bigint;",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "mydom" {
					t.Fatalf("expected name 'mydom', got %q", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil for modify")
				}
			},
		},
		{
			name:    "domain NOT NULL changed",
			fromSQL: "CREATE DOMAIN mydom AS integer;",
			toSQL:   "CREATE DOMAIN mydom AS integer NOT NULL;",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.NotNull != false {
					t.Fatal("expected From.NotNull to be false")
				}
				if e.To.NotNull != true {
					t.Fatal("expected To.NotNull to be true")
				}
			},
		},
		{
			name:    "domain default changed",
			fromSQL: "CREATE DOMAIN mydom AS integer DEFAULT 0;",
			toSQL:   "CREATE DOMAIN mydom AS integer DEFAULT 42;",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
			},
		},
		{
			name:    "domain constraint added",
			fromSQL: "CREATE DOMAIN mydom AS integer;",
			toSQL:   "CREATE DOMAIN mydom AS integer CONSTRAINT pos CHECK (VALUE > 0);",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.From.Constraints) != 0 {
					t.Fatalf("expected 0 from constraints, got %d", len(e.From.Constraints))
				}
				if len(e.To.Constraints) != 1 {
					t.Fatalf("expected 1 to constraint, got %d", len(e.To.Constraints))
				}
			},
		},
		{
			name:    "domain constraint dropped",
			fromSQL: "CREATE DOMAIN mydom AS integer CONSTRAINT pos CHECK (VALUE > 0);",
			toSQL:   "CREATE DOMAIN mydom AS integer;",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.From.Constraints) != 1 {
					t.Fatalf("expected 1 from constraint, got %d", len(e.From.Constraints))
				}
				if len(e.To.Constraints) != 0 {
					t.Fatalf("expected 0 to constraints, got %d", len(e.To.Constraints))
				}
			},
		},
		{
			name:    "domain constraint expression changed",
			fromSQL: "CREATE DOMAIN mydom AS integer CONSTRAINT pos CHECK (VALUE > 0);",
			toSQL:   "CREATE DOMAIN mydom AS integer CONSTRAINT pos CHECK (VALUE > 10);",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
			},
		},
		{
			name:    "domain unchanged produces no entry",
			fromSQL: "CREATE DOMAIN mydom AS integer NOT NULL DEFAULT 0 CONSTRAINT pos CHECK (VALUE >= 0);",
			toSQL:   "CREATE DOMAIN mydom AS integer NOT NULL DEFAULT 0 CONSTRAINT pos CHECK (VALUE >= 0);",
			check: func(t *testing.T, entries []DomainDiffEntry) {
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries for unchanged domain, got %d", len(entries))
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
			entries := diffDomains(from, to)
			tt.check(t, entries)
		})
	}
}

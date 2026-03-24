package catalog

import (
	"testing"
)

func TestDiffExtension(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (from, to *Catalog)
		check func(t *testing.T, entries []ExtensionDiffEntry)
	}{
		{
			name: "extension added",
			setup: func(t *testing.T) (*Catalog, *Catalog) {
				from := New()
				to := New()
				execSQL(t, to, `CREATE EXTENSION citext;`)
				return from, to
			},
			check: func(t *testing.T, entries []ExtensionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "citext" {
					t.Fatalf("expected name 'citext', got %q", e.Name)
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
			name: "extension dropped",
			setup: func(t *testing.T) (*Catalog, *Catalog) {
				from := New()
				execSQL(t, from, `CREATE EXTENSION citext;`)
				to := New()
				return from, to
			},
			check: func(t *testing.T, entries []ExtensionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "citext" {
					t.Fatalf("expected name 'citext', got %q", e.Name)
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
			name: "extension schema changed",
			setup: func(t *testing.T) (*Catalog, *Catalog) {
				from := New()
				execSQL(t, from, `CREATE EXTENSION citext;`)

				to := New()
				execSQL(t, to, `CREATE SCHEMA utils;`)
				execSQL(t, to, `CREATE EXTENSION citext SCHEMA utils;`)
				return from, to
			},
			check: func(t *testing.T, entries []ExtensionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "citext" {
					t.Fatalf("expected name 'citext', got %q", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil for modify")
				}
			},
		},
		{
			name: "extension relocatable flag changed",
			setup: func(t *testing.T) (*Catalog, *Catalog) {
				// Directly construct catalogs with different relocatable flags
				// since there's no DDL to change this property.
				from := New()
				pubOID := from.schemaByName["public"].OID
				fromExt := &Extension{
					OID:         from.oidGen.Next(),
					Name:        "myext",
					SchemaOID:   pubOID,
					Relocatable: true,
				}
				from.extensions[fromExt.OID] = fromExt
				from.extByName[fromExt.Name] = fromExt

				to := New()
				pubOID2 := to.schemaByName["public"].OID
				toExt := &Extension{
					OID:         to.oidGen.Next(),
					Name:        "myext",
					SchemaOID:   pubOID2,
					Relocatable: false,
				}
				to.extensions[toExt.OID] = toExt
				to.extByName[toExt.Name] = toExt

				return from, to
			},
			check: func(t *testing.T, entries []ExtensionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "myext" {
					t.Fatalf("expected name 'myext', got %q", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil for modify")
				}
				if !e.From.Relocatable {
					t.Fatal("expected From.Relocatable to be true")
				}
				if e.To.Relocatable {
					t.Fatal("expected To.Relocatable to be false")
				}
			},
		},
		{
			name: "extension unchanged produces no entry",
			setup: func(t *testing.T) (*Catalog, *Catalog) {
				from := New()
				execSQL(t, from, `CREATE EXTENSION citext;`)
				to := New()
				execSQL(t, to, `CREATE EXTENSION citext;`)
				return from, to
			},
			check: func(t *testing.T, entries []ExtensionDiffEntry) {
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries, got %d", len(entries))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to := tt.setup(t)
			entries := diffExtensions(from, to)
			tt.check(t, entries)
		})
	}
}

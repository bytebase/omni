package catalog

import (
	"testing"
)

func TestDiffView(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []RelationDiffEntry)
	}{
		{
			name:    "view added",
			fromSQL: "CREATE TABLE t1 (id int, name text);",
			toSQL:   "CREATE TABLE t1 (id int, name text); CREATE VIEW v1 AS SELECT id FROM t1;",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "v1" {
					t.Fatalf("expected name v1, got %s", e.Name)
				}
				if e.SchemaName != "public" {
					t.Fatalf("expected schema public, got %s", e.SchemaName)
				}
				if e.From != nil {
					t.Fatal("expected From to be nil")
				}
				if e.To == nil {
					t.Fatal("expected To to be non-nil")
				}
				if e.To.RelKind != 'v' {
					t.Fatalf("expected RelKind 'v', got %q", e.To.RelKind)
				}
			},
		},
		{
			name:    "view dropped",
			fromSQL: "CREATE TABLE t1 (id int, name text); CREATE VIEW v1 AS SELECT id FROM t1;",
			toSQL:   "CREATE TABLE t1 (id int, name text);",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "v1" {
					t.Fatalf("expected name v1, got %s", e.Name)
				}
				if e.From == nil {
					t.Fatal("expected From to be non-nil")
				}
				if e.To != nil {
					t.Fatal("expected To to be nil")
				}
			},
		},
		{
			name:    "view definition changed",
			fromSQL: "CREATE TABLE t1 (id int, name text); CREATE VIEW v1 AS SELECT id FROM t1;",
			toSQL:   "CREATE TABLE t1 (id int, name text); CREATE VIEW v1 AS SELECT id, name FROM t1;",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "v1" {
					t.Fatalf("expected name v1, got %s", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil")
				}
			},
		},
		{
			name:    "view check option changed local to cascaded",
			fromSQL: "CREATE TABLE t1 (id int); CREATE VIEW v1 AS SELECT id FROM t1 WITH LOCAL CHECK OPTION;",
			toSQL:   "CREATE TABLE t1 (id int); CREATE VIEW v1 AS SELECT id FROM t1 WITH CASCADED CHECK OPTION;",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.From.CheckOption != 'l' {
					t.Fatalf("expected From CheckOption 'l', got %q", e.From.CheckOption)
				}
				if e.To.CheckOption != 'c' {
					t.Fatalf("expected To CheckOption 'c', got %q", e.To.CheckOption)
				}
			},
		},
		{
			name:    "materialized view added",
			fromSQL: "CREATE TABLE t1 (id int, name text);",
			toSQL:   "CREATE TABLE t1 (id int, name text); CREATE MATERIALIZED VIEW mv1 AS SELECT id FROM t1;",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.Name != "mv1" {
					t.Fatalf("expected name mv1, got %s", e.Name)
				}
				if e.To == nil {
					t.Fatal("expected To to be non-nil")
				}
				if e.To.RelKind != 'm' {
					t.Fatalf("expected RelKind 'm', got %q", e.To.RelKind)
				}
			},
		},
		{
			name:    "materialized view dropped",
			fromSQL: "CREATE TABLE t1 (id int, name text); CREATE MATERIALIZED VIEW mv1 AS SELECT id FROM t1;",
			toSQL:   "CREATE TABLE t1 (id int, name text);",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.Name != "mv1" {
					t.Fatalf("expected name mv1, got %s", e.Name)
				}
				if e.From == nil {
					t.Fatal("expected From to be non-nil")
				}
				if e.From.RelKind != 'm' {
					t.Fatalf("expected RelKind 'm', got %q", e.From.RelKind)
				}
			},
		},
		{
			name:    "materialized view definition changed",
			fromSQL: "CREATE TABLE t1 (id int, name text); CREATE MATERIALIZED VIEW mv1 AS SELECT id FROM t1;",
			toSQL:   "CREATE TABLE t1 (id int, name text); CREATE MATERIALIZED VIEW mv1 AS SELECT id, name FROM t1;",
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.Name != "mv1" {
					t.Fatalf("expected name mv1, got %s", e.Name)
				}
				if e.From == nil || e.To == nil {
					t.Fatal("expected both From and To to be non-nil")
				}
				if e.From.RelKind != 'm' || e.To.RelKind != 'm' {
					t.Fatalf("expected RelKind 'm', got from=%q to=%q", e.From.RelKind, e.To.RelKind)
				}
			},
		},
		{
			name:    "view matview treated as relation with appropriate relkind",
			fromSQL: "CREATE TABLE t1 (id int, name text);",
			toSQL: `CREATE TABLE t1 (id int, name text);
				CREATE VIEW v1 AS SELECT id FROM t1;
				CREATE MATERIALIZED VIEW mv1 AS SELECT name FROM t1;`,
			check: func(t *testing.T, entries []RelationDiffEntry) {
				if len(entries) != 2 {
					t.Fatalf("expected 2 entries, got %d", len(entries))
				}
				// Sorted by name: mv1, v1
				mv := entries[0]
				v := entries[1]
				if mv.Name != "mv1" {
					t.Fatalf("expected first entry mv1, got %s", mv.Name)
				}
				if v.Name != "v1" {
					t.Fatalf("expected second entry v1, got %s", v.Name)
				}
				if mv.To.RelKind != 'm' {
					t.Fatalf("expected mv1 RelKind 'm', got %q", mv.To.RelKind)
				}
				if v.To.RelKind != 'v' {
					t.Fatalf("expected v1 RelKind 'v', got %q", v.To.RelKind)
				}
				if mv.Action != DiffAdd || v.Action != DiffAdd {
					t.Fatalf("expected both DiffAdd, got mv=%d v=%d", mv.Action, v.Action)
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
			entries := diffViews(from, to)
			tt.check(t, entries)
		})
	}
}

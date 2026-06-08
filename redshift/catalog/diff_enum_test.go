package catalog

import (
	"testing"
)

func TestDiffEnum(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []EnumDiffEntry)
	}{
		{
			name:    "enum type added",
			fromSQL: "",
			toSQL:   "CREATE TYPE status AS ENUM ('active', 'inactive');",
			check: func(t *testing.T, entries []EnumDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", e.Action)
				}
				if e.SchemaName != "public" || e.Name != "status" {
					t.Fatalf("expected public.status, got %s.%s", e.SchemaName, e.Name)
				}
				if len(e.FromValues) != 0 {
					t.Fatalf("expected no FromValues, got %v", e.FromValues)
				}
				if len(e.ToValues) != 2 || e.ToValues[0] != "active" || e.ToValues[1] != "inactive" {
					t.Fatalf("expected ToValues [active inactive], got %v", e.ToValues)
				}
			},
		},
		{
			name:    "enum type dropped",
			fromSQL: "CREATE TYPE status AS ENUM ('active', 'inactive');",
			toSQL:   "",
			check: func(t *testing.T, entries []EnumDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", e.Action)
				}
				if e.SchemaName != "public" || e.Name != "status" {
					t.Fatalf("expected public.status, got %s.%s", e.SchemaName, e.Name)
				}
				if len(e.FromValues) != 2 || e.FromValues[0] != "active" || e.FromValues[1] != "inactive" {
					t.Fatalf("expected FromValues [active inactive], got %v", e.FromValues)
				}
			},
		},
		{
			name:    "enum value added appended",
			fromSQL: "CREATE TYPE status AS ENUM ('active', 'inactive');",
			toSQL:   "CREATE TYPE status AS ENUM ('active', 'inactive'); ALTER TYPE status ADD VALUE 'pending';",
			check: func(t *testing.T, entries []EnumDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.FromValues) != 2 {
					t.Fatalf("expected 2 FromValues, got %d: %v", len(e.FromValues), e.FromValues)
				}
				if len(e.ToValues) != 3 {
					t.Fatalf("expected 3 ToValues, got %d: %v", len(e.ToValues), e.ToValues)
				}
				if e.ToValues[2] != "pending" {
					t.Fatalf("expected last ToValue to be 'pending', got %q", e.ToValues[2])
				}
			},
		},
		{
			name:    "enum value added with BEFORE positioning",
			fromSQL: "CREATE TYPE status AS ENUM ('active', 'inactive');",
			toSQL:   "CREATE TYPE status AS ENUM ('active', 'inactive'); ALTER TYPE status ADD VALUE 'pending' BEFORE 'active';",
			check: func(t *testing.T, entries []EnumDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.ToValues) != 3 {
					t.Fatalf("expected 3 ToValues, got %d: %v", len(e.ToValues), e.ToValues)
				}
				if e.ToValues[0] != "pending" {
					t.Fatalf("expected first ToValue to be 'pending', got %q", e.ToValues[0])
				}
			},
		},
		{
			name:    "enum values reordered detected",
			fromSQL: "CREATE TYPE color AS ENUM ('red', 'green', 'blue');",
			toSQL:   "CREATE TYPE color AS ENUM ('blue', 'green', 'red');",
			check: func(t *testing.T, entries []EnumDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if e.FromValues[0] != "red" || e.FromValues[2] != "blue" {
					t.Fatalf("expected FromValues [red green blue], got %v", e.FromValues)
				}
				if e.ToValues[0] != "blue" || e.ToValues[2] != "red" {
					t.Fatalf("expected ToValues [blue green red], got %v", e.ToValues)
				}
			},
		},
		{
			name:    "enum type identity by schema and name",
			fromSQL: "CREATE SCHEMA s1; CREATE TYPE s1.status AS ENUM ('a'); CREATE TYPE status AS ENUM ('x');",
			toSQL:   "CREATE SCHEMA s1; CREATE TYPE s1.status AS ENUM ('a', 'b'); CREATE TYPE status AS ENUM ('x');",
			check: func(t *testing.T, entries []EnumDiffEntry) {
				// Only s1.status should be modified; public.status is unchanged.
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.SchemaName != "s1" || e.Name != "status" {
					t.Fatalf("expected s1.status, got %s.%s", e.SchemaName, e.Name)
				}
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
			},
		},
		{
			name:    "enum unchanged produces no entry",
			fromSQL: "CREATE TYPE status AS ENUM ('active', 'inactive');",
			toSQL:   "CREATE TYPE status AS ENUM ('active', 'inactive');",
			check: func(t *testing.T, entries []EnumDiffEntry) {
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
			entries := diffEnums(from, to)
			tt.check(t, entries)
		})
	}
}

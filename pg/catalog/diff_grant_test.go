package catalog

import (
	"testing"
)

func TestDiffGrant(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []GrantDiffEntry)
	}{
		{
			name:    "grant on table added",
			fromSQL: `CREATE TABLE t1 (id int);`,
			toSQL:   `CREATE TABLE t1 (id int); GRANT SELECT ON t1 TO some_role;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.To.Privilege != "select" {
					t.Errorf("expected select privilege, got %s", e.To.Privilege)
				}
				if e.To.Grantee != "some_role" {
					t.Errorf("expected grantee some_role, got %s", e.To.Grantee)
				}
				if e.To.ObjType != 'r' {
					t.Errorf("expected ObjType 'r', got %c", e.To.ObjType)
				}
			},
		},
		{
			name:    "grant on table revoked",
			fromSQL: `CREATE TABLE t1 (id int); GRANT SELECT ON t1 TO some_role;`,
			toSQL:   `CREATE TABLE t1 (id int);`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Errorf("expected DiffDrop, got %d", e.Action)
				}
				if e.From.Privilege != "select" {
					t.Errorf("expected select privilege, got %s", e.From.Privilege)
				}
				if e.From.Grantee != "some_role" {
					t.Errorf("expected grantee some_role, got %s", e.From.Grantee)
				}
			},
		},
		{
			name: "grant privilege changed SELECT to ALL",
			fromSQL: `CREATE TABLE t1 (id int);
				GRANT SELECT ON t1 TO some_role;`,
			toSQL: `CREATE TABLE t1 (id int);
				GRANT ALL ON t1 TO some_role;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				// SELECT removed, ALL added = 2 entries (drop select, add ALL)
				if len(entries) != 2 {
					t.Fatalf("expected 2 entries, got %d", len(entries))
				}
				var hasAdd, hasDrop bool
				for _, e := range entries {
					switch e.Action {
					case DiffAdd:
						hasAdd = true
						if e.To.Privilege != "ALL" {
							t.Errorf("expected ALL privilege for add, got %s", e.To.Privilege)
						}
					case DiffDrop:
						hasDrop = true
						if e.From.Privilege != "select" {
							t.Errorf("expected select privilege for drop, got %s", e.From.Privilege)
						}
					}
				}
				if !hasAdd {
					t.Error("expected a DiffAdd entry")
				}
				if !hasDrop {
					t.Error("expected a DiffDrop entry")
				}
			},
		},
		{
			name: "grant WITH GRANT OPTION changed",
			fromSQL: `CREATE TABLE t1 (id int);
				GRANT SELECT ON t1 TO some_role;`,
			toSQL: `CREATE TABLE t1 (id int);
				GRANT SELECT ON t1 TO some_role WITH GRANT OPTION;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffModify {
					t.Errorf("expected DiffModify, got %d", e.Action)
				}
				if e.From.WithGrant {
					t.Error("expected From.WithGrant=false")
				}
				if !e.To.WithGrant {
					t.Error("expected To.WithGrant=true")
				}
			},
		},
		{
			name:    "column-level grant added",
			fromSQL: `CREATE TABLE t1 (id int, name text);`,
			toSQL: `CREATE TABLE t1 (id int, name text);
				GRANT SELECT (name) ON t1 TO some_role;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if len(e.To.Columns) != 1 || e.To.Columns[0] != "name" {
					t.Errorf("expected columns [name], got %v", e.To.Columns)
				}
			},
		},
		{
			name: "column-level grant revoked",
			fromSQL: `CREATE TABLE t1 (id int, name text);
				GRANT SELECT (name) ON t1 TO some_role;`,
			toSQL: `CREATE TABLE t1 (id int, name text);`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffDrop {
					t.Errorf("expected DiffDrop, got %d", e.Action)
				}
				if len(e.From.Columns) != 1 || e.From.Columns[0] != "name" {
					t.Errorf("expected columns [name], got %v", e.From.Columns)
				}
			},
		},
		{
			name:    "grant on function added and revoked",
			fromSQL: `CREATE FUNCTION f1() RETURNS void LANGUAGE sql AS 'SELECT 1';`,
			toSQL: `CREATE FUNCTION f1() RETURNS void LANGUAGE sql AS 'SELECT 1';
				GRANT EXECUTE ON FUNCTION f1() TO some_role;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.To.ObjType != 'f' {
					t.Errorf("expected ObjType 'f', got %c", e.To.ObjType)
				}
				if e.To.Privilege != "execute" {
					t.Errorf("expected execute privilege, got %s", e.To.Privilege)
				}
			},
		},
		{
			name:    "grant on schema added and revoked",
			fromSQL: ``,
			toSQL: `GRANT USAGE ON SCHEMA public TO some_role;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.To.ObjType != 'n' {
					t.Errorf("expected ObjType 'n', got %c", e.To.ObjType)
				}
				if e.To.Privilege != "usage" {
					t.Errorf("expected usage privilege, got %s", e.To.Privilege)
				}
			},
		},
		{
			name:    "grant on sequence added and revoked",
			fromSQL: `CREATE SEQUENCE seq1;`,
			toSQL: `CREATE SEQUENCE seq1;
				GRANT USAGE ON SEQUENCE seq1 TO some_role;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				e := entries[0]
				if e.Action != DiffAdd {
					t.Errorf("expected DiffAdd, got %d", e.Action)
				}
				if e.To.ObjType != 's' {
					t.Errorf("expected ObjType 's', got %c", e.To.ObjType)
				}
				if e.To.Privilege != "usage" {
					t.Errorf("expected usage privilege, got %s", e.To.Privilege)
				}
			},
		},
		{
			name: "grant identity is objType objName grantee privilege tuple",
			fromSQL: `CREATE TABLE t1 (id int);
				GRANT SELECT ON t1 TO role_a;
				GRANT INSERT ON t1 TO role_a;
				GRANT SELECT ON t1 TO role_b;`,
			toSQL: `CREATE TABLE t1 (id int);
				GRANT SELECT ON t1 TO role_a;
				GRANT INSERT ON t1 TO role_a;
				GRANT SELECT ON t1 TO role_b;`,
			check: func(t *testing.T, entries []GrantDiffEntry) {
				// Identical grants → no diff
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries for identical grants, got %d", len(entries))
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
			entries := diffGrants(from, to)
			tt.check(t, entries)
		})
	}
}

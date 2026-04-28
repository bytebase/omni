package completion

import (
	"testing"

	"github.com/bytebase/omni/mysql/catalog"
)

func TestExtractTableRefs(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		offset     int
		wantTables []string          // expected table names
		wantAlias  map[string]string // table -> alias (optional)
		wantDB     map[string]string // table -> database (optional)
		wantAbsent []string          // tables that should NOT appear
	}{
		{
			name:       "simple_select",
			sql:        "SELECT * FROM t WHERE ",
			offset:     22,
			wantTables: []string{"t"},
		},
		{
			name:       "with_alias",
			sql:        "SELECT * FROM t AS x WHERE ",
			offset:     27,
			wantTables: []string{"t"},
			wantAlias:  map[string]string{"t": "x"},
		},
		{
			name:       "join",
			sql:        "SELECT * FROM t1 JOIN t2 ON t1.a = t2.a WHERE ",
			offset:     46,
			wantTables: []string{"t1", "t2"},
		},
		{
			name:       "parenthesized_table_reference_list",
			sql:        "SELECT * FROM (t1, t2) WHERE ",
			offset:     29,
			wantTables: []string{"t1", "t2"},
		},
		{
			name:       "database_qualified",
			sql:        "SELECT * FROM db.t WHERE ",
			offset:     25,
			wantTables: []string{"t"},
			wantDB:     map[string]string{"t": "db"},
		},
		{
			name:       "subquery_no_leak",
			sql:        "SELECT * FROM outer_t WHERE a IN (SELECT b FROM inner_t) AND ",
			offset:     61,
			wantTables: []string{"outer_t"},
			wantAbsent: []string{"inner_t"},
		},
		{
			name:       "update",
			sql:        "UPDATE t SET a = 1 WHERE ",
			offset:     25,
			wantTables: []string{"t"},
		},
		{
			name:       "insert_into",
			sql:        "INSERT INTO t (a, b) VALUES ",
			offset:     28,
			wantTables: []string{"t"},
		},
		{
			name:       "delete_from",
			sql:        "DELETE FROM t WHERE ",
			offset:     20,
			wantTables: []string{"t"},
		},
		{
			name:       "lexer_fallback_incomplete",
			sql:        "SELECT * FROM t WHERE a = ",
			offset:     26,
			wantTables: []string{"t"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs := extractTableRefs(tt.sql, tt.offset)
			got := make(map[string]bool)
			refMap := make(map[string]TableRef)
			for _, r := range refs {
				got[r.Table] = true
				refMap[r.Table] = r
			}
			for _, w := range tt.wantTables {
				if !got[w] {
					t.Errorf("extractTableRefs(%q, %d): missing table %q, got refs=%+v", tt.sql, tt.offset, w, refs)
				}
			}
			for _, a := range tt.wantAbsent {
				if got[a] {
					t.Errorf("extractTableRefs(%q, %d): should not contain %q, got refs=%+v", tt.sql, tt.offset, a, refs)
				}
			}
			if tt.wantAlias != nil {
				for tbl, alias := range tt.wantAlias {
					if r, ok := refMap[tbl]; ok {
						if r.Alias != alias {
							t.Errorf("extractTableRefs(%q, %d): table %q alias want %q got %q", tt.sql, tt.offset, tbl, alias, r.Alias)
						}
					}
				}
			}
			if tt.wantDB != nil {
				for tbl, db := range tt.wantDB {
					if r, ok := refMap[tbl]; ok {
						if r.Database != db {
							t.Errorf("extractTableRefs(%q, %d): table %q database want %q got %q", tt.sql, tt.offset, tbl, db, r.Database)
						}
					}
				}
			}
		})
	}
}

// TestResolveColumnRefScoped tests that resolveColumnRefScoped returns
// columns only from tables referenced in the SQL.
func TestResolveColumnRefScoped(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE DATABASE test", nil)
	cat.SetCurrentDatabase("test")
	cat.Exec("CREATE TABLE t1 (a INT, b INT)", nil)
	cat.Exec("CREATE TABLE t2 (c INT, d INT)", nil)
	cat.Exec("CREATE TABLE t3 (e INT, f INT)", nil)

	tests := []struct {
		name       string
		sql        string
		cursor     int
		wantCols   []string // columns that should appear
		wantAbsent []string // columns that should NOT appear
	}{
		{
			name:       "scoped_single_table",
			sql:        "SELECT * FROM t1 WHERE ",
			cursor:     23,
			wantCols:   []string{"a", "b"},
			wantAbsent: []string{"c", "d", "e", "f"},
		},
		{
			name:       "scoped_join",
			sql:        "SELECT * FROM t1 JOIN t2 ON t1.a = t2.c WHERE ",
			cursor:     46,
			wantCols:   []string{"a", "b", "c", "d"},
			wantAbsent: []string{"e", "f"},
		},
		{
			name:       "scoped_update",
			sql:        "UPDATE t2 SET ",
			cursor:     14,
			wantCols:   []string{"c", "d"},
			wantAbsent: []string{"a", "b", "e", "f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates := resolveColumnRefScoped(cat, tt.sql, tt.cursor)
			for _, col := range tt.wantCols {
				if !containsCandidate(candidates, col, CandidateColumn) {
					t.Errorf("resolveColumnRefScoped(%q, %d): missing column %q", tt.sql, tt.cursor, col)
				}
			}
			for _, col := range tt.wantAbsent {
				if containsCandidate(candidates, col, CandidateColumn) {
					t.Errorf("resolveColumnRefScoped(%q, %d): should not contain column %q", tt.sql, tt.cursor, col)
				}
			}
		})
	}
}

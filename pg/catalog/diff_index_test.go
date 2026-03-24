package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestDiffIndex(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "standalone index added",
			fromSQL: "CREATE TABLE t1 (id int, name text);",
			toSQL:   "CREATE TABLE t1 (id int, name text); CREATE INDEX idx_name ON t1 (name);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", idx.Action)
				}
				if idx.Name != "idx_name" {
					t.Fatalf("expected idx_name, got %s", idx.Name)
				}
				if idx.From != nil {
					t.Fatalf("expected From to be nil")
				}
				if idx.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
			},
		},
		{
			name:    "standalone index dropped",
			fromSQL: "CREATE TABLE t1 (id int, name text); CREATE INDEX idx_name ON t1 (name);",
			toSQL:   "CREATE TABLE t1 (id int, name text);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", idx.Action)
				}
				if idx.Name != "idx_name" {
					t.Fatalf("expected idx_name, got %s", idx.Name)
				}
				if idx.From == nil {
					t.Fatalf("expected From to be non-nil")
				}
				if idx.To != nil {
					t.Fatalf("expected To to be nil")
				}
			},
		},
		{
			name: "index columns changed",
			fromSQL: `CREATE TABLE t1 (a int, b int);
				CREATE INDEX idx_t1 ON t1 (a);`,
			toSQL: `CREATE TABLE t1 (a int, b int);
				CREATE INDEX idx_t1 ON t1 (b);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				if idx.Name != "idx_t1" {
					t.Fatalf("expected idx_t1, got %s", idx.Name)
				}
				if idx.From == nil || idx.To == nil {
					t.Fatalf("expected both From and To to be non-nil")
				}
			},
		},
		{
			name: "index access method changed",
			fromSQL: `CREATE TABLE t1 (id int);
				CREATE INDEX idx_t1 ON t1 USING btree (id);`,
			toSQL: `CREATE TABLE t1 (id int);
				CREATE INDEX idx_t1 ON t1 USING hash (id);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				if idx.From.AccessMethod != "btree" {
					t.Fatalf("expected From access method btree, got %s", idx.From.AccessMethod)
				}
				if idx.To.AccessMethod != "hash" {
					t.Fatalf("expected To access method hash, got %s", idx.To.AccessMethod)
				}
			},
		},
		{
			name: "index uniqueness changed",
			fromSQL: `CREATE TABLE t1 (id int);
				CREATE INDEX idx_t1 ON t1 (id);`,
			toSQL: `CREATE TABLE t1 (id int);
				CREATE UNIQUE INDEX idx_t1 ON t1 (id);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				if idx.From.IsUnique {
					t.Fatalf("expected From.IsUnique to be false")
				}
				if !idx.To.IsUnique {
					t.Fatalf("expected To.IsUnique to be true")
				}
			},
		},
		{
			name: "partial index WHERE clause changed",
			fromSQL: `CREATE TABLE t1 (id int, active bool);
				CREATE INDEX idx_t1 ON t1 (id) WHERE active;`,
			toSQL: `CREATE TABLE t1 (id int, active bool);
				CREATE INDEX idx_t1 ON t1 (id) WHERE (NOT active);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				if idx.From.WhereClause == idx.To.WhereClause {
					t.Fatalf("expected WHERE clauses to differ")
				}
			},
		},
		{
			name: "index INCLUDE columns changed",
			fromSQL: `CREATE TABLE t1 (a int, b int, c int);
				CREATE INDEX idx_t1 ON t1 (a) INCLUDE (b);`,
			toSQL: `CREATE TABLE t1 (a int, b int, c int);
				CREATE INDEX idx_t1 ON t1 (a) INCLUDE (c);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				// NKeyColumns should be the same (1 key column), but total Columns differ.
				if idx.From.NKeyColumns != idx.To.NKeyColumns {
					t.Fatalf("expected same NKeyColumns, got %d vs %d", idx.From.NKeyColumns, idx.To.NKeyColumns)
				}
			},
		},
		{
			name: "index sort options changed",
			fromSQL: `CREATE TABLE t1 (id int);
				CREATE INDEX idx_t1 ON t1 (id ASC);`,
			toSQL: `CREATE TABLE t1 (id int);
				CREATE INDEX idx_t1 ON t1 (id DESC NULLS FIRST);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				// From should have ASC (IndOption[0] = 0), To should have DESC+NULLS_FIRST (IndOption[0] = 3).
				if idx.From.IndOption[0] != 0 {
					t.Fatalf("expected From IndOption[0]=0, got %d", idx.From.IndOption[0])
				}
				if idx.To.IndOption[0] != 3 {
					t.Fatalf("expected To IndOption[0]=3 (DESC|NULLS_FIRST), got %d", idx.To.IndOption[0])
				}
			},
		},
		{
			// NULLS NOT DISTINCT: the SQL parser does not yet support this syntax
			// on CREATE INDEX, so we use programmatic IndexStmt with Nulls_not_distinct.
			name: "NULLS NOT DISTINCT changed",
			fromSQL: `CREATE TABLE t1 (id int);
				CREATE UNIQUE INDEX idx_t1 ON t1 (id);`,
			toSQL: "CREATE TABLE t1 (id int);", // index added programmatically below
			check: func(t *testing.T, d *SchemaDiff) {
				// Overridden by TestDiffIndexNullsNotDistinct below.
			},
		},
		{
			name: "expression index expressions changed",
			fromSQL: `CREATE TABLE t1 (name text);
				CREATE INDEX idx_t1 ON t1 (lower(name));`,
			toSQL: `CREATE TABLE t1 (name text);
				CREATE INDEX idx_t1 ON t1 (upper(name));`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", idx.Action)
				}
				if idx.From == nil || idx.To == nil {
					t.Fatalf("expected both From and To to be non-nil")
				}
				// The Exprs should differ.
				if len(idx.From.Exprs) == 0 || len(idx.To.Exprs) == 0 {
					t.Fatalf("expected non-empty Exprs in both From and To")
				}
				if idx.From.Exprs[0] == idx.To.Exprs[0] {
					t.Fatalf("expected Exprs to differ, both are %q", idx.From.Exprs[0])
				}
			},
		},
		{
			name: "PK/UNIQUE backing index skipped",
			fromSQL: `CREATE TABLE t1 (id int PRIMARY KEY);`,
			toSQL: `CREATE TABLE t1 (id int PRIMARY KEY); CREATE INDEX idx_extra ON t1 (id);`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				// Should have exactly 1 index diff (the standalone idx_extra), not the PK backing index.
				if len(e.Indexes) != 1 {
					t.Fatalf("expected 1 index diff (standalone only), got %d", len(e.Indexes))
				}
				idx := e.Indexes[0]
				if idx.Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", idx.Action)
				}
				if idx.Name != "idx_extra" {
					t.Fatalf("expected idx_extra, got %s", idx.Name)
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
			// Skip the placeholder for NULLS NOT DISTINCT — tested separately.
			if tt.name == "NULLS NOT DISTINCT changed" {
				return
			}
			tt.check(t, d)
		})
	}
}

// TestDiffIndexNullsNotDistinct tests NULLS NOT DISTINCT detection using
// programmatic IndexStmt construction, since the SQL parser does not yet
// support this syntax on CREATE INDEX.
func TestDiffIndexNullsNotDistinct(t *testing.T) {
	// Build 'from' catalog: table + unique index (NullsNotDistinct=false).
	from, err := LoadSQL(`CREATE TABLE t1 (id int); CREATE UNIQUE INDEX idx_t1 ON t1 (id);`)
	if err != nil {
		t.Fatal(err)
	}

	// Build 'to' catalog: same table + unique index with NullsNotDistinct=true.
	to, err := LoadSQL(`CREATE TABLE t1 (id int);`)
	if err != nil {
		t.Fatal(err)
	}
	idxStmt := &nodes.IndexStmt{
		Idxname:            "idx_t1",
		Relation:           &nodes.RangeVar{Relname: "t1"},
		IndexParams:        &nodes.List{Items: []nodes.Node{&nodes.IndexElem{Name: "id"}}},
		Unique:             true,
		Nulls_not_distinct: true,
	}
	if err := to.DefineIndex(idxStmt); err != nil {
		t.Fatal(err)
	}

	d := Diff(from, to)
	if len(d.Relations) != 1 {
		t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
	}
	e := d.Relations[0]
	if len(e.Indexes) != 1 {
		t.Fatalf("expected 1 index diff, got %d", len(e.Indexes))
	}
	idx := e.Indexes[0]
	if idx.Action != DiffModify {
		t.Fatalf("expected DiffModify, got %d", idx.Action)
	}
	if idx.From.NullsNotDistinct {
		t.Fatalf("expected From.NullsNotDistinct to be false")
	}
	if !idx.To.NullsNotDistinct {
		t.Fatalf("expected To.NullsNotDistinct to be true")
	}
}

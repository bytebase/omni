package catalog

import (
	"testing"
)

func TestDiffColumn(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, d *SchemaDiff)
	}{
		{
			name:    "column added to existing table",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int, name text);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffAdd {
					t.Fatalf("expected column DiffAdd, got %d", c.Action)
				}
				if c.Name != "name" {
					t.Fatalf("expected column name 'name', got %s", c.Name)
				}
				if c.To == nil {
					t.Fatalf("expected To to be non-nil")
				}
				if c.From != nil {
					t.Fatalf("expected From to be nil")
				}
			},
		},
		{
			name:    "column dropped from existing table",
			fromSQL: "CREATE TABLE t1 (id int, name text);",
			toSQL:   "CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffDrop {
					t.Fatalf("expected column DiffDrop, got %d", c.Action)
				}
				if c.Name != "name" {
					t.Fatalf("expected column name 'name', got %s", c.Name)
				}
				if c.From == nil {
					t.Fatalf("expected From to be non-nil")
				}
				if c.To != nil {
					t.Fatalf("expected To to be nil")
				}
			},
		},
		{
			name:    "column type changed",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id bigint);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.Name != "id" {
					t.Fatalf("expected column name 'id', got %s", c.Name)
				}
				if c.From == nil || c.To == nil {
					t.Fatalf("expected both From and To to be non-nil")
				}
			},
		},
		{
			name:    "column nullability changed",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int NOT NULL);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.NotNull == c.To.NotNull {
					t.Fatalf("expected NotNull to differ")
				}
			},
		},
		{
			name:    "column default added",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int DEFAULT 42);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.HasDefault {
					t.Fatalf("expected From.HasDefault to be false")
				}
				if !c.To.HasDefault {
					t.Fatalf("expected To.HasDefault to be true")
				}
			},
		},
		{
			name:    "column default removed",
			fromSQL: "CREATE TABLE t1 (id int DEFAULT 42);",
			toSQL:   "CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if !c.From.HasDefault {
					t.Fatalf("expected From.HasDefault to be true")
				}
				if c.To.HasDefault {
					t.Fatalf("expected To.HasDefault to be false")
				}
			},
		},
		{
			name:    "column default value changed",
			fromSQL: "CREATE TABLE t1 (id int DEFAULT 42);",
			toSQL:   "CREATE TABLE t1 (id int DEFAULT 99);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.Default == c.To.Default {
					t.Fatalf("expected Default values to differ")
				}
			},
		},
		{
			name:    "column identity added",
			fromSQL: "CREATE TABLE t1 (id int);",
			toSQL:   "CREATE TABLE t1 (id int GENERATED ALWAYS AS IDENTITY);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.Identity != 0 {
					t.Fatalf("expected From.Identity to be 0, got %c", c.From.Identity)
				}
				if c.To.Identity != 'a' {
					t.Fatalf("expected To.Identity to be 'a', got %c", c.To.Identity)
				}
			},
		},
		{
			name:    "column identity removed",
			fromSQL: "CREATE TABLE t1 (id int GENERATED ALWAYS AS IDENTITY);",
			toSQL:   "CREATE TABLE t1 (id int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.Identity != 'a' {
					t.Fatalf("expected From.Identity to be 'a', got %c", c.From.Identity)
				}
				if c.To.Identity != 0 {
					t.Fatalf("expected To.Identity to be 0, got %c", c.To.Identity)
				}
			},
		},
		{
			name:    "generated column added",
			fromSQL: "CREATE TABLE t1 (a int, b int);",
			toSQL:   "CREATE TABLE t1 (a int, b int GENERATED ALWAYS AS (a * 2) STORED);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.Name != "b" {
					t.Fatalf("expected column name 'b', got %s", c.Name)
				}
				if c.From.Generated != 0 {
					t.Fatalf("expected From.Generated to be 0")
				}
				if c.To.Generated != 's' {
					t.Fatalf("expected To.Generated to be 's', got %c", c.To.Generated)
				}
			},
		},
		{
			name:    "generated column expression changed",
			fromSQL: "CREATE TABLE t1 (a int, b int GENERATED ALWAYS AS (a * 2) STORED);",
			toSQL:   "CREATE TABLE t1 (a int, b int GENERATED ALWAYS AS (a + 10) STORED);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				// Generation expression is stored in Default field by the catalog.
				if c.From.Default == c.To.Default {
					t.Fatalf("expected Default (generation expr) to differ")
				}
			},
		},
		{
			name:    "generated column removed",
			fromSQL: "CREATE TABLE t1 (a int, b int GENERATED ALWAYS AS (a * 2) STORED);",
			toSQL:   "CREATE TABLE t1 (a int, b int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.Generated != 's' {
					t.Fatalf("expected From.Generated to be 's'")
				}
				if c.To.Generated != 0 {
					t.Fatalf("expected To.Generated to be 0")
				}
			},
		},
		{
			name:    "column collation changed",
			fromSQL: `CREATE TABLE t1 (name text COLLATE "C");`,
			toSQL:   `CREATE TABLE t1 (name text COLLATE "POSIX");`,
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if len(e.Columns) != 1 {
					t.Fatalf("expected 1 column diff, got %d", len(e.Columns))
				}
				c := e.Columns[0]
				if c.Action != DiffModify {
					t.Fatalf("expected column DiffModify, got %d", c.Action)
				}
				if c.From.CollationName == c.To.CollationName {
					t.Fatalf("expected CollationName to differ")
				}
			},
		},
		{
			name:    "multiple column changes on same table",
			fromSQL: "CREATE TABLE t1 (a int, b text, c int);",
			toSQL:   "CREATE TABLE t1 (a bigint, b text, d int);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 1 {
					t.Fatalf("expected 1 relation diff, got %d", len(d.Relations))
				}
				e := d.Relations[0]
				if e.Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", e.Action)
				}
				// a: type changed (modify), c: dropped, d: added = 3 column diffs
				if len(e.Columns) != 3 {
					t.Fatalf("expected 3 column diffs, got %d", len(e.Columns))
				}
				// Results should be sorted by name: a, c, d
				if e.Columns[0].Name != "a" || e.Columns[0].Action != DiffModify {
					t.Fatalf("expected first column diff: modify a, got %s action=%d", e.Columns[0].Name, e.Columns[0].Action)
				}
				if e.Columns[1].Name != "c" || e.Columns[1].Action != DiffDrop {
					t.Fatalf("expected second column diff: drop c, got %s action=%d", e.Columns[1].Name, e.Columns[1].Action)
				}
				if e.Columns[2].Name != "d" || e.Columns[2].Action != DiffAdd {
					t.Fatalf("expected third column diff: add d, got %s action=%d", e.Columns[2].Name, e.Columns[2].Action)
				}
			},
		},
		{
			name:    "column unchanged",
			fromSQL: "CREATE TABLE t1 (id int NOT NULL, name text);",
			toSQL:   "CREATE TABLE t1 (id int NOT NULL, name text);",
			check: func(t *testing.T, d *SchemaDiff) {
				if len(d.Relations) != 0 {
					t.Fatalf("expected 0 relation diffs, got %d", len(d.Relations))
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

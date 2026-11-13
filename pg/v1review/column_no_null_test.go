package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

func TestColumnNoNull(t *testing.T) {
	schema := v1review.Schema{
		ObjectCaseSensitive: true,
		Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadata.TableMetadata{{Name: "t", Indexes: []*metadata.IndexMetadata{{Name: "t_idx", Expressions: []string{"b", "c"}}}}},
		}}},
	}
	// nullable is a finding as compared: the anchored text, empty for none.
	type nullable struct {
		Schema, Table, Column string
		Statement             int
		Text                  string
	}
	for _, tt := range []struct {
		name  string
		stmts []string
		want  []nullable
	}{
		{
			name:  "not null, primary key and NOT NULL constraint are excused",
			stmts: []string{"CREATE TABLE t(a INT NOT NULL, b INT PRIMARY KEY, c INT CONSTRAINT n NOT NULL, d INT, PRIMARY KEY (d));"},
		},
		{
			name:  "the search can match inside an earlier word",
			stmts: []string{"create table t (\n  e int\n);"},
			want:  []nullable{{"public", "t", "e", 0, "e"}},
		},
		{
			name:  "each search starts after the previous match",
			stmts: []string{"CREATE TABLE tbl (\n  x INT,\n  b INT\n);"},
			want:  []nullable{{"public", "tbl", "b", 0, "b"}, {"public", "tbl", "x", 0, "x"}},
		},
		{
			name:  "keyed by schema, sorted by schema, table and column",
			stmts: []string{"CREATE TABLE s.t(z INT);", "CREATE TABLE t(a INT);"},
			want:  []nullable{{"public", "t", "a", 1, "a"}, {"s", "t", "z", 0, "z"}},
		},
		{
			name:  "set not null forgets, drop not null records without a range",
			stmts: []string{"CREATE TABLE t(a INT, b INT);", "ALTER TABLE t ALTER COLUMN a SET NOT NULL;", "ALTER TABLE t ALTER COLUMN b DROP NOT NULL;"},
			want:  []nullable{{"public", "t", "b", 2, ""}},
		},
		{
			name:  "primary key using index forgets the index's expressions",
			stmts: []string{"CREATE TABLE t(a INT, b INT, c INT);", "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY USING INDEX t_idx;"},
			want:  []nullable{{"public", "t", "a", 0, "a"}},
		},
		{
			name:  "index names are case sensitive",
			stmts: []string{"CREATE TABLE t(b INT);", "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY USING INDEX \"T_IDX\";"},
			want:  []nullable{{"public", "t", "b", 0, "b"}},
		},
		{
			name:  "a later recording replaces the anchor",
			stmts: []string{"CREATE TABLE t(a INT);", "ALTER TABLE t ADD COLUMN a INT;"},
			want:  []nullable{{"public", "t", "a", 1, "a"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var nodes []ast.Node
			for _, sql := range tt.stmts {
				list, err := parser.Parse(sql)
				if err != nil {
					t.Fatalf("Parse(%q): %v", sql, err)
				}
				nodes = append(nodes, list.Items[0])
			}
			var got []nullable
			for _, f := range ColumnNoNull(nodes, tt.stmts, schema) {
				if f.Code != v1review.ColumnNullable {
					t.Errorf("code %q", f.Code)
				}
				got = append(got, nullable{f.Value, f.Table, f.Name, f.Statement, tt.stmts[f.Statement][f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

func TestDatabaseDropEmptyDatabase(t *testing.T) {
	snapshot := func(name string, tables ...string) v1review.Schema {
		schema := &metadata.SchemaMetadata{}
		for _, table := range tables {
			schema.Tables = append(schema.Tables, &metadata.TableMetadata{Name: table})
		}
		return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Name: name, Schemas: []*metadata.SchemaMetadata{schema}}}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		want   []systemFinding
	}{
		{name: "the current database, empty", sql: "DROP DATABASE shop;", schema: snapshot("shop")},
		{
			name:   "the current database, holding a table",
			sql:    "DROP DATABASE IF EXISTS shop;",
			schema: snapshot("shop", "t"),
			want:   []systemFinding{{Code: v1review.DatabaseNotEmpty, Name: "shop", Text: "DROP DATABASE IF EXISTS shop"}},
		},
		{
			name:   "another database, whatever it holds",
			sql:    "DROP DATABASE other;",
			schema: snapshot("shop"),
			want:   []systemFinding{{Code: v1review.DatabaseNotCurrent, Name: "other", Text: "DROP DATABASE other"}},
		},
		{
			name:   "the name is compared as written",
			sql:    "DROP SCHEMA `Shop`;",
			schema: snapshot("shop"),
			want:   []systemFinding{{Code: v1review.DatabaseNotCurrent, Name: "Shop", Text: "DROP SCHEMA `Shop`"}},
		},
		{
			name: "no snapshot",
			sql:  "DROP DATABASE shop;",
			want: []systemFinding{{Code: v1review.DatabaseNotCurrent, Name: "shop", Text: "DROP DATABASE shop"}},
		},
		{name: "other statements", sql: "DROP TABLE t; CREATE DATABASE other;", schema: snapshot("shop", "t")},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN DROP DATABASE other; END;", schema: snapshot("shop")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return DatabaseDropEmptyDatabase(stmts, tt.schema)
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

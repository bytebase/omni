package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableRequirePK(t *testing.T) {
	pk := &metadata.IndexMetadata{Name: "pkey", Primary: true, Expressions: []string{"id"}}
	table := func(name string, keyed bool) *metadata.TableMetadata {
		t := &metadata.TableMetadata{Name: name}
		if keyed {
			t.Indexes = []*metadata.IndexMetadata{pk}
		}
		return t
	}
	schema := func(tables ...*metadata.TableMetadata) *v1review.Schema {
		return &v1review.Schema{ObjectCaseSensitive: true, Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Name: "public", Tables: tables}}}}
	}
	type got struct {
		Schema, Table string
		Statement     int
	}
	for _, tt := range []struct {
		name            string
		sqls            []string
		lines           []int
		original, final *v1review.Schema
		want            []got
	}{
		{name: "created without a key", sqls: []string{"CREATE TABLE t (id int)"}, original: schema(), final: schema(table("t", false)), want: []got{{"public", "t", 0}}},
		{name: "created with a key", sqls: []string{"CREATE TABLE t (id int PRIMARY KEY)"}, original: schema(), final: schema(table("t", true))},
		{name: "key dropped from a keyed table", sqls: []string{"ALTER TABLE t DROP CONSTRAINT pkey"}, original: schema(table("t", true)), final: schema(table("t", false)), want: []got{{"public", "t", 0}}},
		{name: "keyed table unmentioned still reported", original: schema(table("t", true)), final: schema(table("t", false)), want: []got{{"public", "t", -1}}},
		{name: "unkeyed table out of scope", sqls: []string{"ALTER TABLE t ADD COLUMN v int"}, original: schema(table("t", false)), final: schema(table("t", false))},
		{name: "unkeyed table dropped and recreated", sqls: []string{"DROP TABLE t", "CREATE TABLE t (id int)"}, original: schema(table("t", false)), final: schema(table("t", false)), want: []got{{"public", "t", 1}}},
		{name: "new table not created by the change", sqls: []string{"SELECT 1"}, original: schema(), final: schema(table("t", false))},
		{name: "renamed keyed table looked up by its old name", sqls: []string{"ALTER TABLE t RENAME TO u", "ALTER TABLE u DROP CONSTRAINT pkey"}, original: schema(table("t", true)), final: schema(table("u", false)), want: []got{{"public", "u", 1}}},
		{name: "rename keeps created", sqls: []string{"CREATE TABLE t (id int)", "ALTER TABLE t RENAME TO u"}, original: schema(), final: schema(table("u", false)), want: []got{{"public", "u", 1}}},
		{
			name: "unqualified mention wins only on a later line", sqls: []string{"CREATE TABLE public.t (id int)", "ALTER TABLE t ADD COLUMN v int"}, lines: []int{1, 1},
			original: schema(), final: schema(table("t", false)), want: []got{{"public", "t", 0}},
		},
		{
			name: "unqualified mention on a later line", sqls: []string{"CREATE TABLE public.t (id int)", "ALTER TABLE t ADD COLUMN v int"}, lines: []int{1, 2},
			original: schema(), final: schema(table("t", false)), want: []got{{"public", "t", 1}},
		},
		{name: "nil final reports nothing", sqls: []string{"CREATE TABLE t (id int)"}, original: schema()},
		{name: "nil original holds nothing", sqls: []string{"CREATE TABLE t (id int)"}, final: schema(table("t", false)), want: []got{{"public", "t", 0}}},
		{name: "function body not read", sqls: []string{"CREATE FUNCTION f() RETURNS void AS $$ CREATE TABLE t (id int) $$ LANGUAGE sql"}, original: schema(), final: schema(table("t", false))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stmts []ast.Node
			for _, sql := range tt.sqls {
				list, err := parser.Parse(sql)
				if err != nil {
					t.Fatal(err)
				}
				stmts = append(stmts, list.Items...)
			}
			lines := tt.lines
			if lines == nil {
				for i := range stmts {
					lines = append(lines, i+1)
				}
			}
			var findings []got
			for _, f := range TableRequirePK(stmts, lines, tt.original, tt.final) {
				if f.Code != v1review.TablePrimaryKeyMissing {
					t.Fatalf("code %v", f.Code)
				}
				findings = append(findings, got{f.Value, f.Name, f.Statement})
			}
			if diff := cmp.Diff(tt.want, findings); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

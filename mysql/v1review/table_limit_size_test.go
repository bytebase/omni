package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableLimitSize(t *testing.T) {
	snapshot := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{Tables: []*metadata.TableMetadata{
			{Name: "big", RowCount: 5000},
			{Name: "small", RowCount: 3},
			{Name: "small", RowCount: 99999},
		}},
		{Name: "other", Tables: []*metadata.TableMetadata{{Name: "far", RowCount: 1 << 30}}},
	}}
	type got struct {
		Name, Value string
		Statement   int
	}
	for _, tt := range []struct {
		name     string
		sqls     []string
		snapshot *metadata.DatabaseSchemaMetadata
		maxRows  int
		want     []got
	}{
		{name: "alter on a large table", sqls: []string{"ALTER TABLE big ADD COLUMN x INT"}, snapshot: snapshot, maxRows: 5000, want: []got{{"big", "5000", 0}}},
		{name: "below the limit", sqls: []string{"ALTER TABLE big ADD COLUMN x INT"}, snapshot: snapshot, maxRows: 5001},
		{name: "first table of a name counts", sqls: []string{"TRUNCATE small"}, snapshot: snapshot, maxRows: 4},
		{name: "names compared exactly", sqls: []string{"DROP TABLE BIG"}, snapshot: snapshot, maxRows: 1},
		{name: "only the first schema", sqls: []string{"ALTER TABLE far ADD COLUMN x INT"}, snapshot: snapshot, maxRows: 1},
		{
			name: "missing table counts 0 rows", sqls: []string{"DROP TABLE missing, big"}, snapshot: snapshot, maxRows: 0,
			want: []got{{"missing", "0", 0}, {"big", "5000", 0}},
		},
		{
			name: "every finding anchors the last statement", sqls: []string{"ALTER TABLE big ADD COLUMN x INT", "TRUNCATE TABLE big", "SELECT 1"}, snapshot: snapshot, maxRows: 10,
			want: []got{{"big", "5000", 2}, {"big", "5000", 2}},
		},
		{name: "nil snapshot", sqls: []string{"DROP TABLE big"}, maxRows: 0},
		{name: "no schema", sqls: []string{"DROP TABLE big"}, snapshot: &metadata.DatabaseSchemaMetadata{}, maxRows: 0},
		{name: "create and rename not judged", sqls: []string{"CREATE TABLE big (a INT)", "RENAME TABLE big TO b2"}, snapshot: snapshot, maxRows: 0},
		{name: "routine body not judged", sqls: []string{"CREATE PROCEDURE p() BEGIN DROP TABLE big; END"}, snapshot: snapshot, maxRows: 0},
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
			var findings []got
			for _, f := range TableLimitSize(stmts, tt.snapshot, tt.maxRows) {
				if f.Code != v1review.TableTooLarge || f.Start != 0 || f.End != 0 {
					t.Fatalf("finding %+v", f)
				}
				findings = append(findings, got{f.Name, f.Value, f.Statement})
			}
			if diff := cmp.Diff(tt.want, findings); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

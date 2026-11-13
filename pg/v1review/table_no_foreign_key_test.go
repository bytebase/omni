package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableNoForeignKey(t *testing.T) {
	type fk struct{ Schema, Table, Text string }
	for _, tt := range []struct {
		name string
		sql  string
		want []fk
	}{
		{name: "no foreign key", sql: "CREATE TABLE t (a int PRIMARY KEY);"},
		{
			name: "table constraint, schema defaults to public",
			sql:  "CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES u (a));",
			want: []fk{{"public", "t", "CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES u (a))"}},
		},
		{
			name: "column references, once per statement",
			sql:  "CREATE TABLE s.\"T\" (a int REFERENCES u, b int REFERENCES v, FOREIGN KEY (b) REFERENCES v (b));",
			want: []fk{{"s", "T", "CREATE TABLE s.\"T\" (a int REFERENCES u, b int REFERENCES v, FOREIGN KEY (b) REFERENCES v (b))"}},
		},
		{
			name: "alter add constraint",
			sql:  "ALTER TABLE Public.T ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a), ADD CONSTRAINT fk2 FOREIGN KEY (b) REFERENCES u (b);",
			want: []fk{{"public", "t", "ALTER TABLE Public.T ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a), ADD CONSTRAINT fk2 FOREIGN KEY (b) REFERENCES u (b)"}},
		},
		{name: "alter add column references not judged", sql: "ALTER TABLE t ADD COLUMN b int REFERENCES u;"},
		{name: "alter add check not judged", sql: "ALTER TABLE t ADD CONSTRAINT c CHECK (a > 0);"},
		{name: "create table as not judged", sql: "CREATE TABLE t AS SELECT 1 AS a;"},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void AS $$ CREATE TABLE t (a int REFERENCES u) $$ LANGUAGE sql;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []fk
			for _, f := range TableNoForeignKey(list.Items) {
				if f.Code != v1review.TableForeignKey || f.Name != "" {
					t.Errorf("finding %+v", f)
				}
				got = append(got, fk{f.Value, f.Table, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

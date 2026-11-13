package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

// foreignKeyFinding is a TableNoForeignKey finding as the tests compare it,
// with Text the source its range anchors.
type foreignKeyFinding struct {
	Code  v1review.Code
	Table string
	Text  string
}

func TestTableNoForeignKey(t *testing.T) {
	fk := func(table, text string) foreignKeyFinding {
		return foreignKeyFinding{v1review.TableForeignKey, table, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []foreignKeyFinding
	}{
		{name: "no foreign key", sql: "CREATE TABLE t (a INT, PRIMARY KEY (a), UNIQUE KEY uk (a));"},
		{
			name: "create anchors each constraint, name without database",
			sql:  "CREATE TABLE db.t (\n  a INT,\n  CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES u(b),\n  FOREIGN KEY (a) REFERENCES u(c)\n);",
			want: []foreignKeyFinding{
				fk("t", "CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES u(b)"),
				fk("t", "FOREIGN KEY (a) REFERENCES u(c)"),
			},
		},
		{
			name: "alter anchors each constraint",
			sql:  "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u(b), ADD FOREIGN KEY (c) REFERENCES u(d);",
			want: []foreignKeyFinding{
				fk("t", "CONSTRAINT fk FOREIGN KEY (a) REFERENCES u(b)"),
				fk("t", "FOREIGN KEY (c) REFERENCES u(d)"),
			},
		},
		{name: "column references not judged", sql: "CREATE TABLE t (a INT REFERENCES u(b));"},
		{name: "drop foreign key not judged", sql: "ALTER TABLE t DROP FOREIGN KEY fk;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES u(b)); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			var got []foreignKeyFinding
			for _, f := range TableNoForeignKey(list.Items) {
				got = append(got, foreignKeyFinding{f.Code, f.Table, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

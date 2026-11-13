package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableNoForeignKey(t *testing.T) {
	f := func(name, text string) finding { return finding{v1review.TableForeignKey, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{
			name: "create anchors each constraint",
			sql:  "CREATE TABLE db.t (a INT, b INT,\n  FOREIGN KEY (a) REFERENCES u (a),\n  CONSTRAINT fk FOREIGN KEY (b) REFERENCES u (b));",
			want: []finding{f("t", "FOREIGN KEY (a) REFERENCES u (a)"), f("t", "CONSTRAINT fk FOREIGN KEY (b) REFERENCES u (b)")},
		},
		{name: "column references not judged", sql: "CREATE TABLE t (a INT REFERENCES u (a));"},
		{name: "other constraints not judged", sql: "CREATE TABLE t (a INT, PRIMARY KEY (a), UNIQUE (a));"},
		{
			name: "alter add",
			sql:  "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a), ADD INDEX (a);",
			want: []finding{f("t", "CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (a)")},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES u (a)); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, TableNoForeignKey)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableRequireCollation(t *testing.T) {
	f := func(name, text string) finding { return finding{v1review.TableCollationMissing, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "collate", sql: "CREATE TABLE t (a INT) COLLATE utf8mb4_bin;"},
		{name: "default collate", sql: "CREATE TABLE t (a INT) default collate = utf8mb4_bin;"},
		{name: "missing", sql: "CREATE TABLE db.t (a INT);", want: []finding{f("t", "CREATE TABLE db.t (a INT)")}},
		{name: "charset does not satisfy", sql: "CREATE TABLE t (a INT) CHARSET utf8mb4;", want: []finding{f("t", "CREATE TABLE t (a INT) CHARSET utf8mb4")}},
		{name: "column collate does not satisfy", sql: "CREATE TABLE t (a TEXT COLLATE utf8mb4_bin);", want: []finding{f("t", "CREATE TABLE t (a TEXT COLLATE utf8mb4_bin)")}},
		{name: "as select judged", sql: "CREATE TABLE t AS SELECT 1;", want: []finding{f("t", "CREATE TABLE t AS SELECT 1")}},
		{name: "alter not judged", sql: "ALTER TABLE t COLLATE utf8mb4_bin;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, TableRequireCollation)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

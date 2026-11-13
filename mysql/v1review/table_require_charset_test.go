package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableRequireCharset(t *testing.T) {
	f := func(name, text string) finding { return finding{v1review.TableCharsetMissing, name, text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "charset", sql: "CREATE TABLE t (a INT) CHARSET utf8mb4;"},
		{name: "default character set", sql: "CREATE TABLE t (a INT) default character set = utf8mb4;"},
		{name: "missing", sql: "CREATE TABLE db.t (a INT);", want: []finding{f("t", "CREATE TABLE db.t (a INT)")}},
		{name: "collate does not satisfy", sql: "CREATE TABLE t (a INT) COLLATE utf8mb4_bin;", want: []finding{f("t", "CREATE TABLE t (a INT) COLLATE utf8mb4_bin")}},
		{name: "like judged", sql: "CREATE TABLE t LIKE u;", want: []finding{f("t", "CREATE TABLE t LIKE u")}},
		{name: "alter not judged", sql: "ALTER TABLE t ADD COLUMN b INT;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, TableRequireCharset)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

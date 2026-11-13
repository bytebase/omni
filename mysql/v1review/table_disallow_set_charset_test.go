package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableDisallowSetCharset(t *testing.T) {
	f := func(text string) finding { return finding{v1review.TableCharsetSet, "", text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []finding
	}{
		{name: "charset", sql: "CREATE TABLE t (a INT) CHARSET utf8mb4;", want: []finding{f("CREATE TABLE t (a INT) CHARSET utf8mb4")}},
		{
			name: "one finding per option",
			sql:  "CREATE TABLE t (a INT) DEFAULT CHARSET=utf8 character set utf8;",
			want: []finding{f("CREATE TABLE t (a INT) DEFAULT CHARSET=utf8 character set utf8"), f("CREATE TABLE t (a INT) DEFAULT CHARSET=utf8 character set utf8")},
		},
		{name: "collate not judged", sql: "CREATE TABLE t (a INT) COLLATE utf8mb4_bin;"},
		{name: "column charset not judged", sql: "CREATE TABLE t (a VARCHAR(10) CHARSET utf8);"},
		{
			name: "convert to",
			sql:  "ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4;",
			want: []finding{f("ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4")},
		},
		{name: "alter table option not judged", sql: "ALTER TABLE t DEFAULT CHARSET utf8mb4;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a INT) CHARSET utf8; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tableFindings(t, tt.sql, TableDisallowSetCharset)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

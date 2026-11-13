package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnDisallowChangeType(t *testing.T) {
	changed := func(table, name, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnTypeChanged, Table: table, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []ruleFinding
	}{
		{name: "other actions ok", sql: "ALTER TABLE t ADD COLUMN a int, ALTER COLUMN b SET DEFAULT 1, DROP COLUMN c;"},
		{
			name: "one finding per statement, first column",
			sql:  "ALTER TABLE s.t ALTER COLUMN a TYPE bigint, ALTER b SET DATA TYPE text;\nALTER TABLE u ALTER COLUMN c TYPE int USING c::int;",
			want: []ruleFinding{
				changed("t", "a", "ALTER TABLE s.t ALTER COLUMN a TYPE bigint, ALTER b SET DATA TYPE text"),
				changed("u", "c", "ALTER TABLE u ALTER COLUMN c TYPE int USING c::int"),
			},
		},
		{
			name: "any relation",
			sql:  "ALTER FOREIGN TABLE f ALTER COLUMN a TYPE text;",
			want: []ruleFinding{changed("f", "a", "ALTER FOREIGN TABLE f ALTER COLUMN a TYPE text")},
		},
		{name: "alter type attribute judged too", sql: "ALTER TYPE typ ALTER ATTRIBUTE a TYPE text;", want: []ruleFinding{changed("typ", "a", "ALTER TYPE typ ALTER ATTRIBUTE a TYPE text")}},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC ALTER TABLE t ALTER COLUMN a TYPE text; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, ColumnDisallowChangeType)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnDisallowChangeType (-want +got):\n%s", diff)
			}
		})
	}
}

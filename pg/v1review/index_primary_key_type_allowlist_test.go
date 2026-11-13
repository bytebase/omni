package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestIndexPrimaryKeyTypeAllowlist(t *testing.T) {
	notAllowed := func(name, columnType, text string) indexFinding {
		return indexFinding{Code: v1review.PrimaryKeyTypeNotAllowed, Table: "t", Name: name, Type: columnType, Search: name, Text: text}
	}
	allow := []string{"serial", "bigint", "character varying(20)"}
	for _, tt := range []struct {
		name string
		sql  string
		want []indexFinding
	}{
		{name: "allowed through equivalence", sql: "CREATE TABLE t (a int PRIMARY KEY, b int8 PRIMARY KEY, c varchar(20) PRIMARY KEY, d serial PRIMARY KEY)"},
		{
			name: "only serial, bigserial and smallserial are aliases",
			sql:  "CREATE TABLE t (d serial4 PRIMARY KEY)",
			want: []indexFinding{notAllowed("d", "serial4", "CREATE TABLE t (d serial4 PRIMARY KEY)")},
		},
		{
			name: "column constraints first, then the table's, each constraint once",
			sql:  "CREATE TABLE s.t (a text, b varchar(10) PRIMARY KEY PRIMARY KEY, PRIMARY KEY (a, missing, b))",
			want: []indexFinding{
				notAllowed("b", "character varying(10)", "CREATE TABLE s.t (a text, b varchar(10) PRIMARY KEY PRIMARY KEY, PRIMARY KEY (a, missing, b))"),
				notAllowed("b", "character varying(10)", "CREATE TABLE s.t (a text, b varchar(10) PRIMARY KEY PRIMARY KEY, PRIMARY KEY (a, missing, b))"),
				notAllowed("a", "text", "CREATE TABLE s.t (a text, b varchar(10) PRIMARY KEY PRIMARY KEY, PRIMARY KEY (a, missing, b))"),
				notAllowed("b", "character varying(10)", "CREATE TABLE s.t (a text, b varchar(10) PRIMARY KEY PRIMARY KEY, PRIMARY KEY (a, missing, b))"),
			},
		},
		{
			name: "the type is the last name part with integer modifiers and array",
			sql:  "CREATE TABLE t (a pg_catalog.int4[] PRIMARY KEY, b numeric(10,2) PRIMARY KEY, c timestamp(3) with time zone PRIMARY KEY)",
			want: []indexFinding{
				notAllowed("a", "int4[]", "CREATE TABLE t (a pg_catalog.int4[] PRIMARY KEY, b numeric(10,2) PRIMARY KEY, c timestamp(3) with time zone PRIMARY KEY)"),
				notAllowed("b", "numeric(10,2)", "CREATE TABLE t (a pg_catalog.int4[] PRIMARY KEY, b numeric(10,2) PRIMARY KEY, c timestamp(3) with time zone PRIMARY KEY)"),
				notAllowed("c", "timestamptz(3)", "CREATE TABLE t (a pg_catalog.int4[] PRIMARY KEY, b numeric(10,2) PRIMARY KEY, c timestamp(3) with time zone PRIMARY KEY)"),
			},
		},
		{
			name: "a listed column takes its last definition's type",
			sql:  "CREATE TABLE t (a bigint, a text, PRIMARY KEY (a))",
			want: []indexFinding{notAllowed("a", "text", "CREATE TABLE t (a bigint, a text, PRIMARY KEY (a))")},
		},
		{
			name: "a column without a type has the empty type",
			sql:  "CREATE TABLE t OF ty (a PRIMARY KEY)",
			want: []indexFinding{notAllowed("a", "", "CREATE TABLE t OF ty (a PRIMARY KEY)")},
		},
		{name: "alter table is not judged", sql: "ALTER TABLE t ADD COLUMN a text PRIMARY KEY"},
		{name: "create schema elements are not judged", sql: "CREATE SCHEMA s CREATE TABLE t (a text PRIMARY KEY)"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sqls := []string{tt.sql}
			got := indexFindings(sqls, IndexPrimaryKeyTypeAllowlist(parseEach(t, sqls), allow))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
	if got := IndexPrimaryKeyTypeAllowlist(parseEach(t, []string{"CREATE TABLE t (a OF_NO_TYPE PRIMARY KEY)"}), []string{""}); len(got) != 1 {
		t.Errorf("an empty entry allows nothing but the empty type: got %v", got)
	}
}

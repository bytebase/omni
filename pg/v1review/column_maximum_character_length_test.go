package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnMaximumCharacterLength(t *testing.T) {
	tooLong := func(schema, table, name, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnCharTooLong, Table: table, Name: name, Value: schema, Search: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []ruleFinding
	}{
		{name: "within maximum", sql: "CREATE TABLE t (a char(20), b character(20));"},
		{name: "varchar and text not judged", sql: "CREATE TABLE t (a varchar(255), b character varying(300), c text, d \"char\");"},
		{
			name: "first violation per statement only",
			sql:  "CREATE TABLE t (a char(10), b char(30), c character(40));\nCREATE TABLE s.u (d national character(21));",
			want: []ruleFinding{
				tooLong("", "t", "b", "CREATE TABLE t (a char(10), b char(30), c character(40))"),
				tooLong("s", "u", "d", "CREATE TABLE s.u (d national character(21))"),
			},
		},
		{
			name: "alter add column and alter column type",
			sql:  "ALTER TABLE t ADD COLUMN a int, ALTER COLUMN b TYPE char(50), ADD COLUMN c char(60);",
			want: []ruleFinding{tooLong("", "t", "b", "ALTER TABLE t ADD COLUMN a int, ALTER COLUMN b TYPE char(50), ADD COLUMN c char(60)")},
		},
		{
			name: "catalog qualifier dropped",
			sql:  `ALTER TABLE d."S"."T" ADD COLUMN "Ç" nchar(25);`,
			want: []ruleFinding{tooLong("S", "T", "Ç", `ALTER TABLE d."S"."T" ADD COLUMN "Ç" nchar(25)`)},
		},
		{name: "length of a type named generically is not read", sql: "CREATE TABLE t (a bpchar(30), b pg_catalog.bpchar(30));"},
		{name: "char without length", sql: "CREATE TABLE t (a char, b character);"},
		{name: "array of char judged", sql: "CREATE TABLE t (a char(30)[]);", want: []ruleFinding{tooLong("", "t", "a", "CREATE TABLE t (a char(30)[])")}},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC CREATE TABLE t (a char(30)); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return ColumnMaximumCharacterLength(stmts, 20) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnMaximumCharacterLength (-want +got):\n%s", diff)
			}
		})
	}
}

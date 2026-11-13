package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnTypeDisallowList(t *testing.T) {
	disallowed := func(table, name, typ, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnTypeDisallowed, Table: table, Name: name, Value: typ, Search: name, Text: text}
	}
	for _, tt := range []struct {
		name       string
		sql        string
		disallowed []string
		want       []ruleFinding
	}{
		{name: "allowed types", sql: "CREATE TABLE t (a int, b text);", disallowed: []string{"JSON", "bigint"}},
		{
			name: "create, add column and alter type, in order", sql: "CREATE TABLE s.t (a json, b int8);\nALTER TABLE t ADD COLUMN c JSON, ALTER COLUMN d TYPE bigint;", disallowed: []string{"JSON", "BIGINT"},
			want: []ruleFinding{
				disallowed("t", "a", "json", "CREATE TABLE s.t (a json, b int8)"),
				disallowed("t", "b", "bigint", "CREATE TABLE s.t (a json, b int8)"),
				disallowed("t", "c", "json", "ALTER TABLE t ADD COLUMN c JSON, ALTER COLUMN d TYPE bigint"),
				disallowed("t", "d", "bigint", "ALTER TABLE t ADD COLUMN c JSON, ALTER COLUMN d TYPE bigint"),
			},
		},
		{
			name: "families and serial", sql: "CREATE TABLE t (a integer, b serial, c bool, d float8, e real);", disallowed: []string{"int", "boolean", "double precision", "float4"},
			want: []ruleFinding{
				disallowed("t", "a", "int", "CREATE TABLE t (a integer, b serial, c bool, d float8, e real)"),
				disallowed("t", "b", "int", "CREATE TABLE t (a integer, b serial, c bool, d float8, e real)"),
				disallowed("t", "c", "boolean", "CREATE TABLE t (a integer, b serial, c bool, d float8, e real)"),
				disallowed("t", "d", "double precision", "CREATE TABLE t (a integer, b serial, c bool, d float8, e real)"),
				disallowed("t", "e", "float4", "CREATE TABLE t (a integer, b serial, c bool, d float8, e real)"),
			},
		},
		{name: "serial is compared as written", sql: "CREATE TABLE t (a serial);", disallowed: []string{"serial"}, want: []ruleFinding{disallowed("t", "a", "serial", "CREATE TABLE t (a serial)")}},
		{
			name: "modifiers are not compared", sql: "CREATE TABLE t (a varchar(20), b varchar(30), c int[]);", disallowed: []string{"varchar(20)", "character varying", "integer"},
			want: []ruleFinding{
				disallowed("t", "a", "character varying", "CREATE TABLE t (a varchar(20), b varchar(30), c int[])"),
				disallowed("t", "b", "character varying", "CREATE TABLE t (a varchar(20), b varchar(30), c int[])"),
				disallowed("t", "c", "integer", "CREATE TABLE t (a varchar(20), b varchar(30), c int[])"),
			},
		},
		{
			name: "first match in list order", sql: "CREATE TABLE t (a int);", disallowed: []string{"INTEGER", "int4", "int"},
			want: []ruleFinding{disallowed("t", "a", "integer", "CREATE TABLE t (a int)")},
		},
		{name: "column without type not judged", sql: `CREATE TABLE t OF typ (a WITH OPTIONS NOT NULL);`, disallowed: []string{""}},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC CREATE TABLE t (a json); END;", disallowed: []string{"json"}},
		{name: "empty list", sql: "CREATE TABLE t (a json);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return ColumnTypeDisallowList(stmts, tt.disallowed) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnTypeDisallowList (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTypesEquivalent(t *testing.T) {
	for _, tt := range []struct {
		a, b string
		want bool
	}{
		{"int4", "INTEGER", true},
		{"varchar", "character varying", true},
		{"varchar(20)", "character  varying ( 20 )", true},
		{"char(5)", "character(5)", true},
		{"char(5)", "character(6)", false},
		{"bigserial", "int8", true},
		{"smallserial", "int2", true},
		{"float8", "double   precision", true},
		{"text", "varchar", false},
		{"", "", true},
		{"int", "", false},
	} {
		if got := typesEquivalent(tt.a, tt.b); got != tt.want {
			t.Errorf("typesEquivalent(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		if got := typesEquivalent(tt.b, tt.a); got != tt.want {
			t.Errorf("typesEquivalent(%q, %q) = %v, want %v", tt.b, tt.a, got, tt.want)
		}
	}
}

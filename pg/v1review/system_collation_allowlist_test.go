package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

func TestSystemCollationAllowlist(t *testing.T) {
	bad := func(name, text string) systemFinding {
		return systemFinding{Code: v1review.CollationNotAllowed, Name: name, Text: text}
	}
	allowed := []string{"C", "en_US"}
	for _, tt := range []struct {
		name    string
		sql     string
		allowed []string
		want    []systemFinding
	}{
		{name: "allowed", sql: `CREATE TABLE t (a text COLLATE "C", b text COLLATE "en_US", c int);`, allowed: allowed},
		{name: "no collation", sql: "CREATE TABLE t (a text);", allowed: allowed},
		{
			name:    "compared exactly",
			sql:     `CREATE TABLE t (a text COLLATE "c", b text COLLATE "EN_US");`,
			allowed: allowed,
			want: []systemFinding{
				bad("c", `CREATE TABLE t (a text COLLATE "c", b text COLLATE "EN_US")`),
				bad("EN_US", `CREATE TABLE t (a text COLLATE "c", b text COLLATE "EN_US")`),
			},
		},
		{
			name:    "an unquoted name is folded to lower case",
			sql:     "CREATE TABLE t (a text COLLATE En_US);",
			allowed: allowed,
			want:    []systemFinding{bad("en_us", "CREATE TABLE t (a text COLLATE En_US)")},
		},
		{
			name:    "only the last part of a qualified name is judged",
			sql:     `CREATE TABLE t (a text COLLATE pg_catalog."C", b text COLLATE public."x");`,
			allowed: allowed,
			want:    []systemFinding{bad("x", `CREATE TABLE t (a text COLLATE pg_catalog."C", b text COLLATE public."x")`)},
		},
		{
			name:    "alter table add column and alter column type",
			sql:     `ALTER TABLE t ADD COLUMN a text COLLATE "x", ALTER COLUMN b TYPE text COLLATE "y", ALTER COLUMN c SET DEFAULT 'z';`,
			allowed: allowed,
			want: []systemFinding{
				bad("x", `ALTER TABLE t ADD COLUMN a text COLLATE "x", ALTER COLUMN b TYPE text COLLATE "y", ALTER COLUMN c SET DEFAULT 'z'`),
				bad("y", `ALTER TABLE t ADD COLUMN a text COLLATE "x", ALTER COLUMN b TYPE text COLLATE "y", ALTER COLUMN c SET DEFAULT 'z'`),
			},
		},
		{
			name:    "empty list allows nothing",
			sql:     `CREATE TABLE t (a text COLLATE "C");`,
			allowed: nil,
			want:    []systemFinding{bad("C", `CREATE TABLE t (a text COLLATE "C")`)},
		},
		{
			name:    "a collation outside a column definition is not judged",
			sql:     `SELECT a COLLATE "x" FROM t; CREATE INDEX i ON t (a COLLATE "x"); CREATE DOMAIN d AS text COLLATE "x"; CREATE DATABASE d LC_COLLATE 'x'; CREATE TABLE u AS SELECT 'a' COLLATE "x" AS a;`,
			allowed: allowed,
		},
		{name: "a routine body is not judged", sql: `CREATE FUNCTION f() RETURNS void LANGUAGE sql AS $$ CREATE TABLE t (a text COLLATE "x") $$;`, allowed: allowed},
		{
			name:    "each statement in order",
			sql:     "CREATE TABLE t (a text COLLATE \"x\");\nALTER TABLE t ADD COLUMN b text;\nALTER TABLE t ADD COLUMN c text COLLATE \"y\";",
			allowed: allowed,
			want: []systemFinding{
				bad("x", `CREATE TABLE t (a text COLLATE "x")`),
				bad("y", `ALTER TABLE t ADD COLUMN c text COLLATE "y"`),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemCollationAllowlist(stmts, tt.allowed)
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

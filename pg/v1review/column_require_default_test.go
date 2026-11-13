package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnRequireDefault(t *testing.T) {
	noDefault := func(table, name, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnNoDefault, Table: table, Name: name, Search: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []ruleFinding
	}{
		{name: "defaults and serials ok", sql: "CREATE TABLE t (a serial, b BIGSERIAL, c serial8, d smallserial, e serial2, f serial4, g int DEFAULT 1, h text DEFAULT now());"},
		{
			name: "every column without default in order",
			sql:  "CREATE TABLE s.t (a int, b int DEFAULT 1, c int NOT NULL, d int GENERATED ALWAYS AS IDENTITY, e int GENERATED ALWAYS AS (a) STORED);",
			want: []ruleFinding{
				noDefault("t", "a", "CREATE TABLE s.t (a int, b int DEFAULT 1, c int NOT NULL, d int GENERATED ALWAYS AS IDENTITY, e int GENERATED ALWAYS AS (a) STORED)"),
				noDefault("t", "c", "CREATE TABLE s.t (a int, b int DEFAULT 1, c int NOT NULL, d int GENERATED ALWAYS AS IDENTITY, e int GENERATED ALWAYS AS (a) STORED)"),
				noDefault("t", "d", "CREATE TABLE s.t (a int, b int DEFAULT 1, c int NOT NULL, d int GENERATED ALWAYS AS IDENTITY, e int GENERATED ALWAYS AS (a) STORED)"),
				noDefault("t", "e", "CREATE TABLE s.t (a int, b int DEFAULT 1, c int NOT NULL, d int GENERATED ALWAYS AS IDENTITY, e int GENERATED ALWAYS AS (a) STORED)"),
			},
		},
		{
			name: "alter add column on any relation",
			sql:  "ALTER TABLE t ADD COLUMN a int DEFAULT 0, ADD b int, ALTER COLUMN c DROP DEFAULT; ALTER FOREIGN TABLE f ADD COLUMN x int;",
			want: []ruleFinding{
				noDefault("t", "b", "ALTER TABLE t ADD COLUMN a int DEFAULT 0, ADD b int, ALTER COLUMN c DROP DEFAULT"),
				noDefault("f", "x", "ALTER FOREIGN TABLE f ADD COLUMN x int"),
			},
		},
		{
			name: "typed table column without type",
			sql:  `CREATE TABLE t OF typ ("A" WITH OPTIONS NOT NULL);`,
			want: []ruleFinding{noDefault("t", "A", `CREATE TABLE t OF typ ("A" WITH OPTIONS NOT NULL)`)},
		},
		{name: "schema-qualified serial is still serial", sql: "CREATE TABLE t (a pg_catalog.serial);"},
		{name: "create table as and routine body not judged", sql: "CREATE TABLE t AS SELECT 1 AS a; CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC CREATE TABLE u (a int); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, ColumnRequireDefault)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnRequireDefault (-want +got):\n%s", diff)
			}
		})
	}
}

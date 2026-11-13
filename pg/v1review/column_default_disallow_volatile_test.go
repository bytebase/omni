package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnDefaultDisallowVolatile(t *testing.T) {
	volatile := func(table, name, text string) ruleFinding {
		return ruleFinding{Code: v1review.ColumnDefaultVolatile, Table: table, Name: name, Search: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []ruleFinding
	}{
		{name: "constant default ok", sql: "ALTER TABLE t ADD COLUMN a int DEFAULT 1, ADD COLUMN b text DEFAULT 'x'::text;"},
		{name: "sql value function is no call", sql: "ALTER TABLE t ADD COLUMN a timestamp DEFAULT CURRENT_TIMESTAMP, ADD COLUMN b text DEFAULT current_user;"},
		{
			name: "any call counts, nested too",
			sql:  "ALTER TABLE s.t ADD COLUMN a timestamp DEFAULT now(), ADD COLUMN b int, ADD c text DEFAULT 'x' || lower('Y');",
			want: []ruleFinding{
				volatile("t", "a", "ALTER TABLE s.t ADD COLUMN a timestamp DEFAULT now(), ADD COLUMN b int, ADD c text DEFAULT 'x' || lower('Y')"),
				volatile("t", "c", "ALTER TABLE s.t ADD COLUMN a timestamp DEFAULT now(), ADD COLUMN b int, ADD c text DEFAULT 'x' || lower('Y')"),
			},
		},
		{name: "create table not judged", sql: "CREATE TABLE t (a timestamp DEFAULT now());"},
		{name: "set default not judged", sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT now();"},
		{
			name: "any relation",
			sql:  "ALTER FOREIGN TABLE f ADD COLUMN a int DEFAULT random();",
			want: []ruleFinding{volatile("f", "a", "ALTER FOREIGN TABLE f ADD COLUMN a int DEFAULT random()")},
		},
		{name: "generated expression is no default", sql: "ALTER TABLE t ADD COLUMN a int GENERATED ALWAYS AS (abs(b)) STORED;"},
		{name: "routine body not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC ALTER TABLE t ADD COLUMN a int DEFAULT now(); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runColumnRule(t, tt.sql, ColumnDefaultDisallowVolatile)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ColumnDefaultDisallowVolatile (-want +got):\n%s", diff)
			}
		})
	}
}

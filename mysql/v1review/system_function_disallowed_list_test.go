package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestSystemFunctionDisallowedList(t *testing.T) {
	call := func(written string) systemFinding {
		return systemFinding{Code: v1review.FunctionDisallowed, Name: strings.ToUpper(written), Text: written + "()"}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []systemFinding
	}{
		{name: "allowed", sql: "SELECT NOW(), COUNT(*) FROM t;"},
		{name: "select list in any case", sql: "SELECT Rand();", want: []systemFinding{call("Rand")}},
		{name: "schema qualifier ignored", sql: "SELECT db.rand();", want: []systemFinding{{Code: v1review.FunctionDisallowed, Name: "RAND", Text: "db.rand()"}}},
		{
			name: "where, group by, having, order by, nested",
			sql:  "SELECT a FROM t WHERE a = ABS(RAND()) GROUP BY UUID() HAVING SLEEP() ORDER BY RAND();",
			want: []systemFinding{call("RAND"), call("UUID"), call("SLEEP"), call("RAND")},
		},
		{
			name: "ctes, set operands and parenthesized bodies",
			sql:  "WITH c AS (SELECT RAND()) SELECT UUID() UNION (SELECT SLEEP());",
			want: []systemFinding{call("RAND"), call("UUID"), call("SLEEP")},
		},
		{
			name: "a subquery in a searched expression is searched whole",
			sql:  "SELECT (SELECT a FROM (SELECT RAND() AS a) d LIMIT 1);",
			want: []systemFinding{call("RAND")},
		},
		{name: "from, limit and window not searched", sql: "SELECT a FROM (SELECT RAND() AS a) d, JSON_TABLE(UUID(), '$' COLUMNS (x INT PATH '$')) j;"},
		{
			name: "insert values and select, not on duplicate key update",
			sql:  "INSERT INTO t VALUES (RAND()) ON DUPLICATE KEY UPDATE a = UUID(); INSERT INTO t SELECT SLEEP();",
			want: []systemFinding{call("RAND"), call("SLEEP")},
		},
		{
			name: "update set and where, not order by",
			sql:  "UPDATE t SET a = RAND() WHERE b = UUID() ORDER BY SLEEP();",
			want: []systemFinding{call("RAND"), call("UUID")},
		},
		{name: "delete where", sql: "DELETE FROM t WHERE a = RAND();", want: []systemFinding{call("RAND")}},
		{
			name: "create table defaults, not generated or as select",
			sql:  "CREATE TABLE t (a INT DEFAULT (RAND()), b CHAR(36) DEFAULT (UUID()), c INT AS (SLEEP())) AS SELECT RAND();",
			want: []systemFinding{call("RAND"), call("UUID")},
		},
		{
			name: "alter table columns and set default",
			sql:  "ALTER TABLE t ADD COLUMN a INT DEFAULT (RAND()), MODIFY b INT DEFAULT (UUID()), ALTER COLUMN c SET DEFAULT (SLEEP());",
			want: []systemFinding{call("RAND"), call("UUID"), call("SLEEP")},
		},
		{name: "explain and routine bodies not searched", sql: "EXPLAIN SELECT RAND(); CREATE PROCEDURE p() BEGIN SELECT RAND(); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemFunctionDisallowedList(stmts, []string{"rand", "Uuid", "SLEEP"})
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

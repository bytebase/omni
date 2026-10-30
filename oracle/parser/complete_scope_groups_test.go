package parser

import (
	"strings"
	"testing"
)

// TestCompletionScopeKeepsToOpenGroups: the scope at the caret is that of
// the query blocks whose parentheses are open at it. Three counts kept by
// depth outlived the group they were taken in. A set operator in an earlier
// group at the same depth started the arm of a SELECT in a later group, so
// (SELECT a FROM t UNION SELECT b FROM u) x, (SELECT | FROM v) y read FROM
// u and missed FROM v. The last SELECT before the caret was its query block
// even when its group had closed, so SELECT a, (SELECT max(b) FROM u), |
// FROM t saw u and not t. And the last SELECT at each depth was an outer
// query block even when it sat in a closed group beside the caret's.
func TestCompletionScopeKeepsToOpenGroups(t *testing.T) {
	tests := []struct {
		name  string
		input string
		local []string
		outer [][]string
	}{
		// A set operator in an earlier group.
		{
			name:  "derived table after a UNION group",
			input: "SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u) x, (SELECT | FROM v) y",
			local: []string{"v"},
			outer: [][]string{{"x", "y"}},
		},
		{
			name:  "WHERE of a derived table after a UNION group",
			input: "SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u) x, (SELECT c FROM v WHERE |) y",
			local: []string{"v"},
			outer: [][]string{{"x", "y"}},
		},
		{
			name:  "joined after a MINUS group",
			input: "SELECT * FROM (SELECT a FROM t MINUS SELECT b FROM u) x JOIN (SELECT | FROM v) y ON x.a = y.c",
			local: []string{"v"},
			outer: [][]string{{"x", "y"}},
		},
		{
			name:  "cross joined after an INTERSECT group",
			input: "SELECT * FROM (SELECT a FROM t INTERSECT SELECT b FROM u) x CROSS JOIN (SELECT | FROM v) y",
			local: []string{"v"},
			outer: [][]string{{"x", "y"}},
		},
		{
			name:  "third derived table",
			input: "SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u) x, (SELECT c FROM v) y, (SELECT | FROM z) w",
			local: []string{"z"},
			outer: [][]string{{"x", "y", "w"}},
		},
		{
			name:  "lateral after a UNION group",
			input: "SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u) x, LATERAL (SELECT | FROM v WHERE v.id = x.a) l",
			local: []string{"v"},
			outer: [][]string{{"x"}},
		},
		{
			name:  "CTE after a CTE with a UNION",
			input: "WITH x AS (SELECT a FROM t UNION SELECT b FROM u), y AS (SELECT | FROM v) SELECT * FROM y",
			local: []string{"v"},
		},
		{
			name:  "IN subquery after an IN subquery with a UNION",
			input: "SELECT * FROM t WHERE a IN (SELECT b FROM u UNION SELECT c FROM w) AND d IN (SELECT | FROM v)",
			local: []string{"v"},
			outer: [][]string{{"t"}},
		},
		{
			name:  "derived table after nested set operations",
			input: "SELECT * FROM ((SELECT a FROM t UNION SELECT b FROM u) MINUS SELECT c FROM w) x, (SELECT | FROM v) y",
			local: []string{"v"},
			outer: [][]string{{"x", "y"}},
		},
		{
			name:  "IN subquery after a parenthesized UNION operand",
			input: "SELECT * FROM t WHERE a IN ((SELECT b FROM u) UNION SELECT c FROM w) AND d IN (SELECT | FROM v)",
			local: []string{"v"},
			outer: [][]string{{"t"}},
		},
		{
			name:  "derived table in the query after a CTE with a UNION",
			input: "WITH x AS (SELECT a FROM t UNION SELECT b FROM u) SELECT * FROM x, (SELECT | FROM v) y",
			local: []string{"v"},
			outer: [][]string{{"x", "y"}},
		},
		{
			name:  "scalar subquery after one with a UNION",
			input: "SELECT (SELECT b FROM u UNION SELECT c FROM w), (SELECT | FROM v) FROM t",
			local: []string{"v"},
			outer: [][]string{{"t"}},
		},
		{
			name:  "EXISTS after a nested UNION",
			input: "SELECT * FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.id IN (SELECT id FROM w UNION SELECT id FROM z) AND EXISTS (SELECT 1 FROM v WHERE |))",
			local: []string{"v"},
			outer: [][]string{{"u"}, {"t"}},
		},

		// A SELECT in a group closed before the caret.
		{
			name:  "select list after a scalar subquery",
			input: "SELECT a, (SELECT max(b) FROM u), | FROM t",
			local: []string{"t"},
		},
		{
			name:  "WHERE after an IN subquery",
			input: "SELECT a FROM t WHERE x IN (SELECT b FROM u) AND |",
			local: []string{"t"},
		},
		{
			name:  "WHERE after a derived table",
			input: "SELECT * FROM (SELECT a FROM t) x WHERE |",
			local: []string{"x"},
		},
		{
			name:  "ORDER BY after derived tables with a UNION",
			input: "SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u) x, (SELECT c FROM v) y ORDER BY |",
			local: []string{"x", "y"},
		},
		{
			name:  "ON after a joined derived table",
			input: "SELECT * FROM t1 JOIN (SELECT a FROM t UNION SELECT b FROM u) x ON x.a = t1.| JOIN v ON 1 = 1",
			local: []string{"t1", "x", "v"},
		},
		{
			name:  "HAVING after a scalar subquery",
			input: "SELECT a FROM t GROUP BY a HAVING count(*) > (SELECT count(*) FROM u) AND |",
			local: []string{"t"},
		},
		{
			name:  "CONNECT BY after a START WITH subquery",
			input: "SELECT * FROM t START WITH id IN (SELECT id FROM u) CONNECT BY PRIOR |",
			local: []string{"t"},
		},
		{
			name:  "nested derived tables",
			input: "SELECT * FROM (SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u) x, (SELECT c FROM v) y) z WHERE z.|",
			local: []string{"z"},
		},
		{
			name:  "correlated subquery after an IN subquery",
			input: "SELECT * FROM t WHERE a IN (SELECT b FROM u) AND nvl((SELECT 1 FROM w WHERE |), 0) = 1",
			local: []string{"w"},
			outer: [][]string{{"t"}},
		},
		{
			name:  "correlated subquery after an IN subquery in a parenthesized condition",
			input: "SELECT * FROM t WHERE (a IN (SELECT b FROM u) OR (EXISTS (SELECT 1 FROM w WHERE |)))",
			local: []string{"w"},
			outer: [][]string{{"t"}},
		},

		// No SELECT around the caret: the DML statement's own tables.
		{
			name:  "UPDATE WHERE after a SET subquery",
			input: "UPDATE t SET a = (SELECT max(b) FROM u) WHERE |",
			local: []string{"t"},
		},
		{
			name:  "DELETE WHERE after an IN subquery",
			input: "DELETE FROM t WHERE x IN (SELECT b FROM u) AND |",
			local: []string{"t"},
		},
		{
			name:  "INSERT VALUES after a scalar subquery",
			input: "INSERT INTO t (c1, c2) VALUES ((SELECT max(a) FROM u), |)",
			local: []string{"t"},
		},
		{
			name:  "MERGE SET after a USING subquery with a UNION",
			input: "MERGE INTO t m USING (SELECT a FROM s UNION SELECT b FROM u) src ON (m.id = src.a) WHEN MATCHED THEN UPDATE SET m.c = |",
			local: []string{"m", "src"},
		},
		{
			name:  "DELETE without FROM",
			input: "DELETE t WHERE a IN (SELECT b FROM u) AND |",
			local: []string{"t"},
		},
		{
			name:  "MERGE source still in parentheses",
			input: "MERGE INTO t m USING ((SELECT a FROM s JOIN u USING (id) CROSS JOIN w)|",
			local: []string{"m"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := collectCompletionMarked(t, tt.input).Scope
			if got := oracleRefNames(scope.LocalReferences); !sameStringSlice(got, tt.local) {
				t.Errorf("local refs = %v, want %v", got, tt.local)
			}
			var outer [][]string
			for _, level := range scope.OuterReferences {
				outer = append(outer, oracleRefNames(level))
			}
			if got, want := joinLevels(outer), joinLevels(tt.outer); got != want {
				t.Errorf("outer refs = %s, want %s", got, want)
			}
		})
	}
}

func joinLevels(levels [][]string) string {
	var parts []string
	for _, level := range levels {
		parts = append(parts, "["+strings.Join(level, " ")+"]")
	}
	return strings.Join(parts, "")
}

package parser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// tableSampleCases covers TABLESAMPLE on a named table. It follows the alias
// and precedes the table hints, so SQL Server rejects an alias after it; it
// takes an undelimited SYSTEM and REPEATABLE; and it needs its sample number
// and seed. An AS with no alias after it, which parsing the alias first would
// otherwise let through as t AS TABLESAMPLE (...), is an error. shape is the AliasedTableRef of an accepted statement, printed by
// tableSampleShape.
var tableSampleCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	// Where TABLESAMPLE goes.
	{"no_alias", "SELECT * FROM t TABLESAMPLE (10 PERCENT)", true, "t TABLESAMPLE (PERCENT)"},
	{"alias_bare", "SELECT * FROM t s TABLESAMPLE (10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"alias_as", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"alias_bracketed", "SELECT * FROM t AS [s] TABLESAMPLE (10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"qualified_table", "SELECT * FROM dbo.t AS s TABLESAMPLE (10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"no_alias_hint", "SELECT * FROM t TABLESAMPLE (10 PERCENT) WITH (NOLOCK)", true, "t TABLESAMPLE (PERCENT) WITH 1"},
	{"alias_hint", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) WITH (NOLOCK)", true, "t AS s TABLESAMPLE (PERCENT) WITH 1"},
	{"alias_as_hint", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) WITH (NOLOCK)", true, "t AS s TABLESAMPLE (PERCENT) WITH 1"},
	{"then_where", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) WHERE s.a = 1", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"then_join", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) JOIN t1 ON s.a = t1.a", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"no_alias_then_join", "SELECT * FROM t TABLESAMPLE (10 PERCENT) JOIN t1 ON t.a = t1.a", true, "t TABLESAMPLE (PERCENT)"},
	{"hint_then_join", "SELECT * FROM t TABLESAMPLE (10 PERCENT) WITH (NOLOCK) JOIN t1 ON t.a = t1.a", true, "t TABLESAMPLE (PERCENT) WITH 1"},
	{"then_comma", "SELECT * FROM t TABLESAMPLE (10 PERCENT), t1", true, "t TABLESAMPLE (PERCENT)"},
	{"then_cross_apply", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) CROSS APPLY (SELECT 1 AS x) ca", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"then_pivot", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) PIVOT (MAX(col) FOR a IN ([1])) p", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"update_from", "UPDATE s SET a = 1 FROM t AS s TABLESAMPLE (10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"delete_from", "DELETE s FROM t AS s TABLESAMPLE (10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"merge_using", "MERGE INTO dbo.t USING t1 AS s TABLESAMPLE (10 PERCENT) ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t1 AS s TABLESAMPLE (PERCENT)"},
	{"alias_after", "SELECT * FROM t TABLESAMPLE (10 PERCENT) s", false, ""},
	{"as_alias_after", "SELECT * FROM t TABLESAMPLE (10 PERCENT) AS s", false, ""},
	{"alias_after_hint", "SELECT * FROM t TABLESAMPLE (10 PERCENT) WITH (NOLOCK) s", false, ""},
	{"alias_before_and_after", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) s2", false, ""},
	{"as_after", "SELECT * FROM t TABLESAMPLE (10 PERCENT) AS", false, ""},
	{"as_without_alias", "SELECT * FROM t AS TABLESAMPLE (10 PERCENT)", false, ""},
	{"hint_before", "SELECT * FROM t AS s WITH (NOLOCK) TABLESAMPLE (10 PERCENT)", false, ""},
	{"hint_before_no_alias", "SELECT * FROM t WITH (NOLOCK) TABLESAMPLE (10 PERCENT)", false, ""},
	{"hint_twice", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) WITH (NOLOCK) WITH (NOLOCK)", false, ""},
	{"holdlock_after", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) HOLDLOCK", false, ""},
	{"twice", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) TABLESAMPLE (10 PERCENT)", false, ""},
	{"derived_table", "SELECT * FROM (SELECT 1 AS a) d TABLESAMPLE (10 PERCENT)", false, ""},
	{"bracketed_keyword", "SELECT * FROM t [TABLESAMPLE] (10 PERCENT)", false, ""},
	{"alias_as_window", "SELECT * FROM t AS window TABLESAMPLE (10 PERCENT)", true, "t AS window TABLESAMPLE (PERCENT)"},

	// An AS that no alias follows is an error, not an empty alias.
	{"dangling_as_end", "SELECT * FROM t AS", false, ""},
	{"dangling_as_where", "SELECT * FROM t AS WHERE 1 = 1", false, ""},
	{"dangling_as_join", "SELECT * FROM t AS JOIN t1 ON t.a = t1.a", false, ""},
	{"dangling_as_derived_table", "SELECT * FROM (SELECT 1 AS a) AS WHERE 1 = 1", false, ""},
	{"dangling_as_rowset", "SELECT * FROM OPENJSON('[]') AS WHERE 1 = 1", false, ""},
	{"dangling_as_merge_source", "MERGE INTO dbo.t USING t1 AS ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},

	// The clause.
	{"system", "SELECT * FROM t AS s TABLESAMPLE SYSTEM (10 PERCENT)", true, "t AS s TABLESAMPLE SYSTEM (PERCENT)"},
	{"system_lower_case", "SELECT * FROM t AS s TABLESAMPLE system (10 PERCENT)", true, "t AS s TABLESAMPLE SYSTEM (PERCENT)"},
	{"system_no_alias", "SELECT * FROM t TABLESAMPLE SYSTEM (10 PERCENT)", true, "t TABLESAMPLE SYSTEM (PERCENT)"},
	{"rows", "SELECT * FROM t AS s TABLESAMPLE (10 ROWS)", true, "t AS s TABLESAMPLE (ROWS)"},
	{"no_unit", "SELECT * FROM t AS s TABLESAMPLE (10)", true, "t AS s TABLESAMPLE ()"},
	{"unit_lower_case", "SELECT * FROM t AS s TABLESAMPLE (10 percent)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"size_decimal", "SELECT * FROM t AS s TABLESAMPLE (10.5 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"size_expression", "SELECT * FROM t AS s TABLESAMPLE (10 + 1 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"size_negative", "SELECT * FROM t AS s TABLESAMPLE (-10 PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"size_parenthesized", "SELECT * FROM t AS s TABLESAMPLE ((10) PERCENT)", true, "t AS s TABLESAMPLE (PERCENT)"},
	{"size_variable", "DECLARE @n INT = 10; SELECT * FROM t AS s TABLESAMPLE (@n PERCENT) REPEATABLE (@n)", true, "t AS s TABLESAMPLE (PERCENT) REPEATABLE"},
	{"repeatable", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5)", true, "t AS s TABLESAMPLE (PERCENT) REPEATABLE"},
	{"repeatable_lower_case", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) repeatable (5)", true, "t AS s TABLESAMPLE (PERCENT) REPEATABLE"},
	{"repeatable_expression", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5 + 1)", true, "t AS s TABLESAMPLE (PERCENT) REPEATABLE"},
	{"everything", "SELECT * FROM t AS s TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5) WITH (NOLOCK)", true, "t AS s TABLESAMPLE SYSTEM (PERCENT) REPEATABLE WITH 1"},
	{"keyword_alone", "SELECT * FROM t s TABLESAMPLE", false, ""},
	{"size_missing", "SELECT * FROM t AS s TABLESAMPLE ()", false, ""},
	{"size_missing_with_unit", "SELECT * FROM t AS s TABLESAMPLE (PERCENT)", false, ""},
	{"size_unparenthesized", "SELECT * FROM t AS s TABLESAMPLE 10 PERCENT", false, ""},
	{"unit_unknown", "SELECT * FROM t AS s TABLESAMPLE (10 ROW)", false, ""},
	{"unit_twice", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT ROWS)", false, ""},
	{"system_bracketed", "SELECT * FROM t AS s TABLESAMPLE [SYSTEM] (10 PERCENT)", false, ""},
	{"system_twice", "SELECT * FROM t AS s TABLESAMPLE SYSTEM SYSTEM (10 PERCENT)", false, ""},
	{"system_alone", "SELECT * FROM t AS s TABLESAMPLE SYSTEM", false, ""},
	{"method_unknown", "SELECT * FROM t AS s TABLESAMPLE BERNOULLI (10 PERCENT)", false, ""},
	{"repeatable_alone", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE", false, ""},
	{"repeatable_empty", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE ()", false, ""},
	{"repeatable_unparenthesized", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE 5", false, ""},
	{"repeatable_bracketed", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) [REPEATABLE] (5)", false, ""},
	{"repeatable_twice", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5) REPEATABLE (6)", false, ""},
	{"repeatable_two_seeds", "SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5, 6)", false, ""},
}

// tableSampleShape prints the table, alias, clause options, and hint count
// of a sampled table.
func tableSampleShape(ref *ast.AliasedTableRef) string {
	s := ref.Table.(*ast.TableRef).Object
	if ref.Alias != "" {
		s += " AS " + ref.Alias
	}
	s += " TABLESAMPLE "
	ts := ref.TableSample
	if ts.System {
		s += "SYSTEM "
	}
	s += "(" + ts.Unit + ")"
	if ts.Repeatable != nil {
		s += " REPEATABLE"
	}
	if ref.Hints != nil {
		s += fmt.Sprintf(" WITH %d", len(ref.Hints.Items))
	}
	return s
}

// TestTableSample checks that omni accepts the TABLESAMPLE forms SQL Server
// accepts, rejects the others, and keeps every Loc inside its parent.
func TestTableSample(t *testing.T) {
	for _, tc := range tableSampleCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error:\n%s", tc.sql, ast.NodeToString(list))
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			ref, ok := nthNodeByType(list, "AliasedTableRef", 0).(*ast.AliasedTableRef)
			if !ok || ref.TableSample == nil {
				t.Fatalf("no AliasedTableRef with TABLESAMPLE in %s", ast.NodeToString(list))
			}
			if got := tableSampleShape(ref); got != tc.shape {
				t.Errorf("shape = %s, want %s", got, tc.shape)
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestTableSampleOracle checks accept/reject parity with SQL Server for every
// case in tableSampleCases.
func TestTableSampleOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range tableSampleCases {
		t.Run(tc.name, func(t *testing.T) {
			ssErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			ssAccepts := ssErr == nil
			if ssAccepts != tc.accepts {
				t.Errorf("table says SQL Server %s but it %s (%v):\n  SQL: %s",
					boolToAcceptReject(tc.accepts), boolToAcceptReject(ssAccepts), ssErr, tc.sql)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil
			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s (%v), omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), ssErr, boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}

// TestTableSampleLoc checks that a sampled table spans its alias, clause, and
// hints, and that the table and the clause span their own text.
func TestTableSampleLoc(t *testing.T) {
	tests := []struct {
		sql  string
		node string
		want string
	}{
		{"SELECT * FROM dbo.t AS s TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5) WITH (NOLOCK) WHERE 1 = 1", "AliasedTableRef", "dbo.t AS s TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5) WITH (NOLOCK)"},
		{"SELECT * FROM dbo.t AS s TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5) WITH (NOLOCK) WHERE 1 = 1", "TableRef", "dbo.t AS s"},
		{"SELECT * FROM dbo.t AS s TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5) WITH (NOLOCK) WHERE 1 = 1", "TableSampleClause", "TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5)"},
		{"SELECT * FROM t TABLESAMPLE (10 PERCENT) JOIN t1 ON t.a = t1.a", "AliasedTableRef", "t TABLESAMPLE (10 PERCENT)"},
		{"SELECT * FROM t TABLESAMPLE (10 PERCENT) JOIN t1 ON t.a = t1.a", "JoinClause", "t TABLESAMPLE (10 PERCENT) JOIN t1 ON t.a = t1.a"},
	}
	for _, tc := range tests {
		t.Run(tc.node+"/"+tc.want, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			n := nthNodeByType(list, tc.node, 0)
			if n == nil {
				t.Fatalf("no %s in %s", tc.node, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(n)); got != tc.want {
				t.Errorf("%s text = %q, want %q", tc.node, got, tc.want)
			}
		})
	}
}

// TestTableSampleErrorPosition checks where a misplaced alias or an
// incomplete clause is reported.
func TestTableSampleErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{"SELECT * FROM t TABLESAMPLE (10 PERCENT) s", `syntax error at or near "s"`},
		{"SELECT * FROM t TABLESAMPLE (10 PERCENT) AS s", `syntax error at or near "AS"`},
		{"SELECT * FROM t AS s TABLESAMPLE ()", `syntax error at or near ")"`},
		{"SELECT * FROM t AS s TABLESAMPLE (PERCENT)", `syntax error at or near "PERCENT"`},
		{"SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE", "syntax error at end of input"},
		{"SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE ()", `syntax error at or near ")"`},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			if got := pe.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

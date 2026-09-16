package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// nestedWithCases lists every position where a SELECT can appear below
// statement level, each with a leading CTE list. SQL Server treats a nested
// SELECT as a query expression, which has no WITH clause, so the CTE is a
// syntax error at the WITH keyword. Only an object body whose grammar names
// select_statement (view, inline TVF RETURN, cursor) carries its own CTE list.
var nestedWithCases = []struct {
	name    string
	sql     string
	accepts bool
}{
	// Query-expression positions: rejected.
	{"derived_table", "SELECT * FROM (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c) d", false},
	{"merge_using_derived_table", "MERGE INTO dbo.t USING (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c) s ON dbo.t.a = s.a WHEN NOT MATCHED THEN INSERT (a) VALUES (s.a);", false},
	{"scalar_subquery_select_list", "SELECT (WITH c AS (SELECT a FROM dbo.t) SELECT MAX(a) FROM c) AS m", false},
	{"scalar_subquery_set", "DECLARE @x INT; SET @x = (WITH c AS (SELECT a FROM dbo.t) SELECT MAX(a) FROM c)", false},
	{"exists", "SELECT a FROM dbo.t WHERE EXISTS (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)", false},
	{"in", "SELECT a FROM dbo.t WHERE a IN (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)", false},
	{"quantified_any", "SELECT a FROM dbo.t WHERE a = ANY (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)", false},
	{"cte_body", "WITH c AS (WITH d AS (SELECT a FROM dbo.t) SELECT a FROM d) SELECT a FROM c", false},
	{"insert_source", "INSERT INTO dbo.t (a) WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", false},
	{"union_right_operand", "SELECT a FROM dbo.t UNION WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", false},
	{"union_right_operand_after_cte", "WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c UNION ALL WITH d AS (SELECT a FROM dbo.t) SELECT a FROM d", false},
	{"except_right_operand", "SELECT a FROM dbo.t EXCEPT WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", false},

	// Object bodies whose select_statement carries a CTE list: accepted.
	{"create_view", "CREATE VIEW dbo.v AS WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", true},
	{"alter_view", "ALTER VIEW dbo.v AS WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", true},
	{"inline_tvf_return_paren", "CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)", true},
	{"inline_tvf_return_no_paren", "CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", true},
	{"declare_cursor_for", "DECLARE cur CURSOR FOR WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c", true},

	// Controls: the same positions without a CTE are fine.
	{"derived_table_plain", "SELECT * FROM (SELECT a FROM dbo.t) d", true},
	{"union_plain", "SELECT a FROM dbo.t UNION SELECT a FROM dbo.t", true},
	{"insert_source_plain", "INSERT INTO dbo.t (a) SELECT a FROM dbo.t", true},
	{"inline_tvf_return_paren_plain", "CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN (SELECT a FROM dbo.t)", true},
	{"inline_tvf_return_no_paren_plain", "CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN SELECT a FROM dbo.t", true},
}

// TestWithClauseNested checks that omni accepts a CTE list only where T-SQL
// does, and that a rejected one fails at the WITH keyword.
func TestWithClauseNested(t *testing.T) {
	for _, tc := range nestedWithCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error at WITH", tc.sql)
				}
				if !strings.Contains(err.Error(), `syntax error at or near "WITH"`) {
					t.Errorf("Parse(%q) error = %q, want it at the WITH keyword", tc.sql, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			if !strings.Contains(tc.sql, "WITH c") {
				return
			}
			// The CTE list must be attached to the SELECT that owns it.
			found := false
			ast.Inspect(list, func(n ast.Node) bool {
				if s, ok := n.(*ast.SelectStmt); ok && s.WithClause != nil {
					found = true
				}
				return true
			})
			if !found {
				t.Errorf("no SelectStmt carries the WithClause: %s", ast.NodeToString(list))
			}
		})
	}
}

// TestWithClauseNestedOracle checks accept/reject parity with SQL Server for
// a CTE list in every nested SELECT position.
func TestWithClauseNestedOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range nestedWithCases {
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

package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/mssql/ast"
)

// parenQueryCases covers parenthesized query expressions, ( <query_expression> ),
// which T-SQL accepts as a query primary wherever a query expression goes. A
// statement's query (SELECT, INSERT, cursor) takes ORDER BY, OFFSET, FOR, and
// OPTION only after the whole expression; a nested one (derived table,
// subquery, CTE, view, inline function) takes them only inside a query
// specification. shape is the query structure of an accepted statement,
// printed by queryShape.
var parenQueryCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	// Derived tables.
	{"derived_paren", "SELECT * FROM ((SELECT 1 AS a)) x", true, "* {1 AS a}"},
	{"derived_paren_triple", "SELECT * FROM (((SELECT 1 AS a))) x", true, "* {1 AS a}"},
	{"derived_union_parens", "SELECT * FROM ((SELECT 1 AS a) UNION (SELECT 2)) x", true, "* {(1 AS a UNION 2)}"},
	{"derived_paren_then_bare", "SELECT * FROM ((SELECT 1 AS a) UNION SELECT 2) x", true, "* {(1 AS a UNION 2)}"},
	{"derived_bare_then_paren", "SELECT * FROM (SELECT 1 AS a UNION (SELECT 2)) x", true, "* {(1 AS a UNION 2)}"},
	{"derived_union_chain", "SELECT * FROM ((SELECT 1 AS a) UNION (SELECT 2) UNION (SELECT 3)) x", true, "* {(1 AS a UNION (2 UNION 3))}"},
	{"derived_union_all", "SELECT * FROM ((SELECT 1 AS a) UNION ALL (SELECT 2)) x", true, "* {(1 AS a UNION ALL 2)}"},
	{"derived_intersect", "SELECT * FROM ((SELECT 1 AS a) INTERSECT (SELECT 2)) x", true, "* {(1 AS a INTERSECT 2)}"},
	{"derived_except", "SELECT * FROM ((SELECT 1 AS a) EXCEPT (SELECT 2)) x", true, "* {(1 AS a EXCEPT 2)}"},
	{"derived_grouped_right", "SELECT * FROM (SELECT 1 AS a UNION ((SELECT 2) UNION (SELECT 3))) x", true, "* {(1 AS a UNION (2 UNION 3))}"},
	{"derived_column_list", "SELECT * FROM ((SELECT 1 AS a)) AS x(c)", true, "* {1 AS a}"},
	{"derived_comma", "SELECT * FROM ((SELECT 1 AS a)) x, ((SELECT 2 AS b)) y", true, "* {1 AS a, 2 AS b}"},
	{"derived_then_join", "SELECT * FROM ((SELECT 1 AS a) UNION (SELECT 2)) x JOIN t ON x.a = t.a", true, "* {(1 AS a UNION 2)}"},
	{"derived_join_right", "SELECT * FROM t JOIN ((SELECT 1 AS a)) x ON t.a = x.a", true, "* {1 AS a}"},
	{"derived_cross_apply", "SELECT * FROM ((SELECT 1 AS a)) x CROSS APPLY ((SELECT 2 AS b)) y", true, "* {1 AS a, 2 AS b}"},
	{"derived_outer_apply_union", "SELECT * FROM t OUTER APPLY ((SELECT 1 AS a) UNION (SELECT 2)) x", true, "* {(1 AS a UNION 2)}"},
	{"derived_pivot", "SELECT * FROM ((SELECT 1 AS a)) x PIVOT (MAX(a) FOR a IN ([1])) p", true, "* {1 AS a}"},
	{"derived_top_order_inside", "SELECT * FROM ((SELECT TOP 1 a FROM t ORDER BY a)) x", true, "* {a ORDER BY}"},
	{"derived_top_order_operand", "SELECT * FROM ((SELECT TOP 1 a FROM t ORDER BY a) UNION (SELECT 2)) x", true, "* {(a ORDER BY UNION 2)}"},
	{"derived_offset_inside", "SELECT * FROM ((SELECT a FROM t ORDER BY a OFFSET 0 ROWS)) x", true, "* {a ORDER BY OFFSET}"},
	{"derived_for_xml_inside", "SELECT * FROM ((SELECT a FROM t FOR XML PATH)) x", true, "* {a FOR}"},
	{"derived_last_operand_order", "SELECT * FROM ((SELECT 1 AS a) UNION SELECT TOP 1 a FROM t ORDER BY 1) x", true, "* {(1 AS a UNION a ORDER BY)}"},
	{"derived_last_operand_offset", "SELECT * FROM ((SELECT 1 AS a) UNION SELECT 2 ORDER BY 1 OFFSET 0 ROWS) x", true, "* {(1 AS a UNION 2 ORDER BY OFFSET)}"},
	{"merge_using", "MERGE INTO t USING ((SELECT 1 AS a)) s ON t.a = s.a WHEN MATCHED THEN DELETE;", true, "1 AS a"},
	{"update_from", "UPDATE t SET a = 1 FROM t JOIN ((SELECT 1 AS a)) x ON t.a = x.a", true, "1 AS a"},
	{"delete_from", "DELETE t FROM t JOIN ((SELECT 1 AS a) UNION (SELECT 2)) x ON t.a = x.a", true, "(1 AS a UNION 2)"},

	// A '((' in FROM may open a joined table whose first operand is a
	// derived table instead.
	{"join_operand_paren", "SELECT * FROM (((SELECT 1 AS a)) s JOIN t2 ON s.a = t2.a)", true, "* {1 AS a}"},
	{"join_operand_union", "SELECT * FROM (((SELECT 1 AS a) UNION (SELECT 2)) x JOIN t ON x.a = t.a)", true, "* {(1 AS a UNION 2)}"},
	{"join_operand_aliased", "SELECT * FROM ((SELECT 1 AS a) AS x)", false, ""},
	{"join_operand_union_aliased", "SELECT * FROM ((SELECT 1 AS a) UNION (SELECT 2) x)", false, ""},

	// Nothing follows the ')' of a nested query expression.
	{"derived_order_after_paren", "SELECT * FROM ((SELECT 1 AS a) ORDER BY a) x", false, ""},
	{"derived_order_after_paren_top", "SELECT * FROM ((SELECT TOP 1 a FROM t) ORDER BY a) x", false, ""},
	{"derived_offset_after_paren", "SELECT * FROM ((SELECT a FROM t) ORDER BY a OFFSET 0 ROWS) x", false, ""},
	{"derived_for_after_paren", "SELECT * FROM ((SELECT 1 AS a) FOR XML PATH) x", false, ""},
	{"derived_option_after_paren", "SELECT * FROM ((SELECT 1 AS a) OPTION (MAXDOP 1)) x", false, ""},
	{"derived_order_after_union", "SELECT * FROM ((SELECT 1 AS a) UNION (SELECT 2) ORDER BY 1) x", false, ""},
	{"derived_into_inside", "SELECT * FROM ((SELECT 1 AS a INTO #t)) x", false, ""},
	{"derived_values_paren", "SELECT * FROM ((VALUES (1))) v(a)", false, ""},
	{"derived_with_paren", "SELECT * FROM ((WITH c AS (SELECT 1 a) SELECT * FROM c)) x", false, ""},
	{"derived_hint", "SELECT * FROM ((SELECT 1 AS a)) WITH (NOLOCK) x", false, ""},

	// A statement.
	{"stmt_paren", "(SELECT 1)", true, "1"},
	{"stmt_paren_double", "((SELECT 1))", true, "1"},
	{"stmt_paren_quadruple", "((((SELECT 1))))", true, "1"},
	{"stmt_union_parens", "(SELECT 1) UNION (SELECT 2)", true, "(1 UNION 2)"},
	{"stmt_paren_then_bare", "(SELECT 1) UNION SELECT 2", true, "(1 UNION 2)"},
	{"stmt_bare_then_paren", "SELECT 1 UNION (SELECT 2)", true, "(1 UNION 2)"},
	{"stmt_bare_paren_bare", "SELECT 1 UNION (SELECT 2) UNION SELECT 3", true, "(1 UNION (2 UNION 3))"},
	{"stmt_union_all_double", "(SELECT 1) UNION ALL ((SELECT 2))", true, "(1 UNION ALL 2)"},
	{"stmt_except_intersect", "(SELECT 1) EXCEPT (SELECT 2) INTERSECT (SELECT 3)", true, "(1 EXCEPT (2 INTERSECT 3))"},
	{"stmt_grouped_right", "SELECT 1 UNION ((SELECT 2) UNION (SELECT 3))", true, "(1 UNION (2 UNION 3))"},
	{"stmt_grouped_left", "((SELECT 1) UNION (SELECT 2)) UNION SELECT 3", true, "((1 UNION 2) UNION 3)"},
	{"stmt_grouped_whole", "(( SELECT 1 ) UNION ( SELECT 2 ))", true, "(1 UNION 2)"},
	{"stmt_grouped_nested", "((((SELECT 1)) UNION ((SELECT 2))))", true, "(1 UNION 2)"},
	{"stmt_set_op_inside", "(SELECT 1 UNION SELECT 2)", true, "(1 UNION 2)"},
	{"stmt_full_spec", "(SELECT 1 FROM t WHERE a = 1 GROUP BY a HAVING a > 0)", true, "1"},
	{"stmt_distinct_top", "(SELECT DISTINCT TOP 3 a FROM t)", true, "a"},
	{"stmt_window", "(SELECT a FROM t WINDOW w AS (ORDER BY a))", true, "a"},
	{"stmt_semicolon", "(SELECT 1);", true, "1"},
	{"stmt_after_statement", "SELECT 1 (SELECT 2)", true, "1, 2"},
	{"stmt_after_semicolon", "SELECT 1; (SELECT 2)", true, "1, 2"},
	{"stmt_in_if", "IF 1 = 1 (SELECT 1) ELSE (SELECT 2)", true, "1, 2"},
	{"stmt_in_while", "WHILE 1 = 0 (SELECT 1)", true, "1"},
	{"stmt_in_block", "BEGIN (SELECT 1) END", true, "1"},

	// The clauses that close a statement follow the whole query
	// expression, on the last operand like an unparenthesized one.
	{"stmt_order", "(SELECT 1) ORDER BY 1", true, "1 ORDER BY"},
	{"stmt_order_union", "(SELECT 1 AS a) UNION (SELECT 2) ORDER BY a", true, "(1 AS a UNION 2 ORDER BY)"},
	{"stmt_order_bare_then_paren", "SELECT 1 UNION (SELECT 2) ORDER BY 1", true, "(1 UNION 2 ORDER BY)"},
	{"stmt_order_paren_then_bare", "(SELECT 1) UNION SELECT 2 ORDER BY 1", true, "(1 UNION 2 ORDER BY)"},
	{"stmt_order_chain", "(SELECT 1) UNION ALL (SELECT 2) UNION (SELECT 3) ORDER BY 1", true, "(1 UNION ALL (2 UNION 3 ORDER BY))"},
	{"stmt_order_grouped", "((SELECT 1) UNION (SELECT 2)) ORDER BY 1", true, "(1 UNION 2) ORDER BY"},
	{"stmt_order_option", "( (SELECT 1) ) ORDER BY 1 OPTION (MAXDOP 1)", true, "1 ORDER BY OPTION"},
	{"stmt_offset_fetch", "(SELECT 1 AS a) ORDER BY a OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY", true, "1 AS a ORDER BY OFFSET"},
	{"stmt_option", "(SELECT 1) OPTION (MAXDOP 1)", true, "1 OPTION"},
	{"stmt_for_xml", "(SELECT 1) FOR XML PATH", true, "1 FOR"},
	{"stmt_for_xml_option", "(SELECT 1) FOR XML PATH OPTION (MAXDOP 1)", true, "1 FOR OPTION"},
	{"stmt_order_for_option", "(SELECT 1) ORDER BY 1 FOR XML PATH OPTION (MAXDOP 1)", true, "1 ORDER BY FOR OPTION"},
	{"stmt_for_browse", "(SELECT a FROM t) FOR BROWSE", true, "a FOR"},
	{"stmt_union_for_xml", "(SELECT a FROM t) UNION (SELECT a FROM t) FOR XML PATH", true, "(a UNION a FOR)"},
	{"stmt_into", "(SELECT 1 INTO #t)", true, "1 INTO"},
	{"stmt_into_from", "(SELECT 1 AS a INTO #t FROM t)", true, "1 AS a INTO"},
	{"stmt_into_first_operand", "(SELECT 1 INTO #t) UNION (SELECT 2)", true, "(1 INTO UNION 2)"},
	{"stmt_into_second_operand", "(SELECT 1) UNION (SELECT 2 INTO #t)", true, "(1 UNION 2 INTO)"},
	{"stmt_into_double", "((SELECT 1 INTO #t))", true, "1 INTO"},
	{"with_paren", "WITH c AS (SELECT 1 a) (SELECT * FROM c)", true, "WITH[1 a] *"},
	{"with_paren_double", "WITH c AS (SELECT 1 a) ((SELECT * FROM c))", true, "WITH[1 a] *"},
	{"with_paren_union", "WITH c AS (SELECT 1 a) (SELECT * FROM c) UNION (SELECT 2)", true, "(WITH[1 a] * UNION 2)"},
	{"with_paren_order", "WITH c AS (SELECT 1 a) (SELECT * FROM c) ORDER BY 1", true, "WITH[1 a] * ORDER BY"},
	{"with_paren_union_bare_tail", "WITH c AS (SELECT 1 a) (SELECT * FROM c) UNION SELECT 2 ORDER BY 1 OPTION (MAXDOP 1)", true, "(WITH[1 a] * UNION 2 ORDER BY OPTION)"},

	// Inside a statement's parentheses a query specification takes none of
	// those clauses, and after them each comes once, in order.
	{"stmt_order_inside", "(SELECT 1 ORDER BY 1)", false, ""},
	{"stmt_top_order_inside", "(SELECT TOP 1 a FROM t ORDER BY a)", false, ""},
	{"stmt_top_order_inside_operand", "(SELECT TOP 1 a FROM t ORDER BY a) UNION (SELECT 2)", false, ""},
	{"stmt_top_order_inside_right", "SELECT 1 UNION (SELECT TOP 1 a FROM t ORDER BY a)", false, ""},
	{"stmt_order_inside_set_op", "(SELECT 1 UNION SELECT 2 ORDER BY 1)", false, ""},
	{"stmt_offset_inside", "(SELECT a FROM t ORDER BY a OFFSET 0 ROWS)", false, ""},
	{"stmt_for_inside", "(SELECT 1 FOR XML PATH)", false, ""},
	{"stmt_option_inside", "(SELECT 1 OPTION (MAXDOP 1))", false, ""},
	{"with_order_inside", "WITH c AS (SELECT 1 a) (SELECT * FROM c ORDER BY 1)", false, ""},
	{"stmt_order_before_union", "(SELECT 1) ORDER BY 1 UNION (SELECT 2)", false, ""},
	{"stmt_option_before_order", "(SELECT 1) OPTION (MAXDOP 1) ORDER BY 1", false, ""},
	{"stmt_order_twice", "(SELECT 1) ORDER BY 1 ORDER BY 1", false, ""},
	{"stmt_offset_without_order", "(SELECT 1) OFFSET 0 ROWS", false, ""},
	{"stmt_where_after_paren", "(SELECT 1) WHERE 1 = 1", false, ""},
	{"stmt_alias", "(SELECT 1) x", false, ""},
	{"stmt_alias_as", "(SELECT 1) AS x", false, ""},
	{"stmt_union_missing_operand", "(SELECT 1) UNION", false, ""},
	{"stmt_union_empty_parens", "(SELECT 1) UNION ()", false, ""},
	{"stmt_union_expression", "(SELECT 1) UNION (1)", false, ""},
	{"stmt_empty_parens", "()", false, ""},
	{"stmt_no_select_list", "(SELECT)", false, ""},
	{"stmt_unclosed", "(SELECT 1", false, ""},
	{"stmt_extra_close", "(SELECT 1))", false, ""},

	// CTE bodies are nested.
	{"cte_paren", "WITH c AS ((SELECT 1 a)) SELECT * FROM c", true, "WITH[1 a] *"},
	{"cte_union_parens", "WITH c AS ((SELECT 1 a) UNION (SELECT 2)) SELECT * FROM c", true, "WITH[(1 a UNION 2)] *"},
	{"cte_bare_then_paren", "WITH c AS (SELECT 1 a UNION (SELECT 2)) SELECT * FROM c", true, "WITH[(1 a UNION 2)] *"},
	{"cte_recursive", "WITH c(a) AS ((SELECT 1) UNION ALL (SELECT a + 1 FROM c WHERE a < 3)) SELECT * FROM c", true, "WITH[(1 UNION ALL a + 1)] *"},
	{"cte_top_order_inside", "WITH c AS ((SELECT TOP 1 a FROM t ORDER BY a)) SELECT * FROM c", true, "WITH[a ORDER BY] *"},
	{"cte_order_after_paren", "WITH c AS ((SELECT 1 a) ORDER BY a) SELECT * FROM c", false, ""},

	// Subqueries in expressions.
	{"scalar_paren_double", "SELECT ((SELECT 1))", true, "_ {1}"},
	{"scalar_union", "SELECT ((SELECT 1) UNION (SELECT 2))", true, "_ {(1 UNION 2)}"},
	{"scalar_union_alias", "SELECT ((SELECT 1) UNION (SELECT 2)) AS x, 3", true, "_ {(1 UNION 2)}"},
	{"scalar_grouped_nested", "SELECT (((SELECT 1)) UNION ALL ((SELECT 2)))", true, "_ {(1 UNION ALL 2)}"},
	{"scalar_bare_then_paren", "SELECT (SELECT 1 UNION (SELECT 2))", true, "_ {(1 UNION 2)}"},
	{"scalar_union_plus", "SELECT ((SELECT 1) UNION (SELECT 2)) + 1", true, "_ {(1 UNION 2)}"},
	{"scalar_in_case", "SELECT CASE WHEN ((SELECT 1) UNION (SELECT 2)) = 1 THEN 1 END", true, "_ {(1 UNION 2)}"},
	{"scalar_for_xml_inside", "SELECT ((SELECT a FROM t FOR XML PATH('')))", true, "_ {a FOR}"},
	{"scalar_expression", "SELECT ((SELECT 1) + 1)", true, "_ {1}"},
	{"scalar_then_union", "SELECT ((SELECT 1)) UNION (SELECT 2)", true, "(_ {1} UNION 2)"},
	{"scalar_order_after_union", "SELECT ((SELECT 1) UNION (SELECT 2) ORDER BY 1)", false, ""},
	{"where_union_is_null", "SELECT * FROM t WHERE ((SELECT 1) UNION (SELECT 2)) IS NULL", true, "* {(1 UNION 2)}"},
	{"where_union_compare", "SELECT * FROM t WHERE a = ((SELECT 1) UNION (SELECT 2))", true, "* {(1 UNION 2)}"},
	{"where_paren_conditions", "SELECT * FROM t WHERE ((SELECT 1) = 1 AND (SELECT 2) = 2)", true, "* {1, 2}"},
	{"declare_default", "DECLARE @x int = ((SELECT 1) UNION (SELECT 2))", true, "(1 UNION 2)"},
	{"update_set", "UPDATE t SET a = ((SELECT 1) UNION (SELECT 2))", true, "(1 UNION 2)"},

	// IN takes a subquery when its parentheses hold a query expression,
	// and a list otherwise.
	{"in_paren", "SELECT * FROM t WHERE a IN ((SELECT 1))", true, "* {IN 1}"},
	{"in_paren_triple", "SELECT * FROM t WHERE a IN (((SELECT 1)))", true, "* {IN 1}"},
	{"in_union", "SELECT * FROM t WHERE a IN ((SELECT 1) UNION (SELECT 2))", true, "* {IN (1 UNION 2)}"},
	{"in_bare_then_paren", "SELECT * FROM t WHERE a IN (SELECT 1 UNION (SELECT 2))", true, "* {IN (1 UNION 2)}"},
	{"not_in_union", "SELECT * FROM t WHERE a NOT IN ((SELECT 1) UNION (SELECT 2))", true, "* {IN (1 UNION 2)}"},
	{"in_list_subqueries", "SELECT * FROM t WHERE a IN ((SELECT 1), 2)", true, "* {1}"},
	{"in_list_nested_subquery", "SELECT * FROM t WHERE a IN (((SELECT 1)), 2)", true, "* {1}"},
	{"in_list_expression", "SELECT * FROM t WHERE a IN ((SELECT 1) + 1)", true, "* {1}"},
	{"in_both_sides", "SELECT * FROM t WHERE ((SELECT 1)) IN ((SELECT 1))", true, "* {1, IN 1}"},
	{"in_union_in_list", "SELECT * FROM t WHERE a IN ((SELECT 1) UNION (SELECT 2), 3)", false, ""},

	{"exists_paren", "SELECT * FROM t WHERE EXISTS ((SELECT 1))", true, "* {1}"},
	{"exists_union", "SELECT * FROM t WHERE EXISTS ((SELECT 1) UNION (SELECT 2))", true, "* {(1 UNION 2)}"},
	{"exists_bare_then_paren", "SELECT * FROM t WHERE EXISTS (SELECT 1 UNION (SELECT 2))", true, "* {(1 UNION 2)}"},
	{"exists_top_order_inside", "SELECT * FROM t WHERE EXISTS ((SELECT TOP 1 a FROM t ORDER BY a))", true, "* {a ORDER BY}"},
	{"not_exists_grouped", "SELECT * FROM t WHERE NOT EXISTS (((SELECT 1)) UNION ((SELECT 2)))", true, "* {(1 UNION 2)}"},
	{"all_paren", "SELECT * FROM t WHERE a = ALL ((SELECT 1))", true, "* {1}"},
	{"all_union", "SELECT * FROM t WHERE a = ALL ((SELECT 1) UNION (SELECT 2))", true, "* {(1 UNION 2)}"},
	{"some_union", "SELECT * FROM t WHERE a = SOME ((SELECT 1) UNION (SELECT 2))", true, "* {(1 UNION 2)}"},
	{"any_paren_triple", "SELECT * FROM t WHERE a = ANY (((SELECT 1)))", true, "* {1}"},
	{"any_bare_then_paren", "SELECT * FROM t WHERE a = ANY (SELECT 1 UNION (SELECT 2))", true, "* {(1 UNION 2)}"},

	// An INSERT source takes ORDER BY and OFFSET after the whole query
	// expression; OPTION belongs to the INSERT.
	{"insert_paren", "INSERT INTO t (a) (SELECT 1)", true, "1"},
	{"insert_paren_two_columns", "INSERT INTO t (a, b) (SELECT 1, 2)", true, "1, 2"},
	{"insert_paren_table_hint", "INSERT INTO t WITH (TABLOCK) (a) (SELECT 1)", true, "1"},
	{"insert_union_parens", "INSERT INTO t (a) (SELECT 1) UNION (SELECT 2)", true, "(1 UNION 2)"},
	{"insert_grouped", "INSERT INTO t (a) ((SELECT 1) UNION (SELECT 2))", true, "(1 UNION 2)"},
	{"insert_bare_then_paren", "INSERT INTO t SELECT 1 UNION (SELECT 2)", true, "(1 UNION 2)"},
	{"insert_order", "INSERT INTO t (a) (SELECT 1) ORDER BY 1", true, "1 ORDER BY"},
	{"insert_double_order", "INSERT INTO t (a) ((SELECT 1)) ORDER BY 1", true, "1 ORDER BY"},
	{"insert_order_offset", "INSERT INTO t (a) (SELECT 1) ORDER BY 1 OFFSET 0 ROWS", true, "1 ORDER BY OFFSET"},
	{"insert_order_paren_then_bare", "INSERT INTO t (a) (SELECT 1) UNION SELECT 2 ORDER BY 1", true, "(1 UNION 2 ORDER BY)"},
	{"insert_order_bare_then_paren", "INSERT INTO t (a) SELECT 1 UNION (SELECT 2) ORDER BY 1", true, "(1 UNION 2 ORDER BY)"},
	{"insert_option", "INSERT INTO t (a) (SELECT 1) OPTION (MAXDOP 1)", true, "1"},
	{"insert_output_alias", "INSERT INTO t (a) OUTPUT inserted.a AS z (SELECT 1)", true, "1"},
	{"insert_output_star", "INSERT INTO t (a) OUTPUT inserted.* (SELECT 1)", true, "1"},
	{"insert_no_columns", "INSERT INTO t (SELECT 1)", false, ""},
	{"insert_no_columns_double", "INSERT INTO t ((SELECT 1))", false, ""},
	{"insert_top_order_inside", "INSERT INTO t (a) (SELECT TOP 1 a FROM t ORDER BY a)", false, ""},
	{"insert_order_inside", "INSERT INTO t (a) (SELECT 1 ORDER BY 1)", false, ""},
	{"insert_for_after", "INSERT INTO t (a) (SELECT 1) FOR XML PATH", false, ""},
	{"insert_into_inside", "INSERT INTO t (a) (SELECT 1 INTO #x)", false, ""},
	{"insert_alias", "INSERT INTO t (a) (SELECT 1) x", false, ""},
	{"insert_values_paren", "INSERT INTO t (a) (VALUES (1))", false, ""},
	{"insert_exec_paren", "INSERT INTO t (a) (EXEC p)", false, ""},
	// SQL Server reads inserted.a (...) as a method call.
	{"insert_output_method_call", "INSERT INTO t (a) OUTPUT inserted.a (SELECT 1)", false, ""},

	// A view body is nested.
	{"view_paren", "CREATE VIEW v AS (SELECT 1 AS a)", true, "1 AS a"},
	{"view_paren_double", "CREATE VIEW v AS ((SELECT 1 AS a))", true, "1 AS a"},
	{"view_union", "CREATE VIEW v AS (SELECT 1 AS a) UNION (SELECT 2)", true, "(1 AS a UNION 2)"},
	{"view_grouped", "CREATE VIEW v AS ((SELECT 1 AS a) UNION (SELECT 2))", true, "(1 AS a UNION 2)"},
	{"view_top_order_inside", "CREATE VIEW v AS (SELECT TOP 1 a FROM t ORDER BY a)", true, "a ORDER BY"},
	{"view_top_order_operand", "CREATE VIEW v AS (SELECT TOP 1 a FROM t ORDER BY a) UNION (SELECT 2)", true, "(a ORDER BY UNION 2)"},
	{"view_check_option", "CREATE VIEW v AS (SELECT 1 AS a) WITH CHECK OPTION", true, "1 AS a"},
	{"view_cte_paren", "CREATE VIEW v AS WITH c AS (SELECT 1 a) (SELECT a FROM c)", true, "WITH[1 a] a"},
	{"alter_view", "ALTER VIEW v AS (SELECT 1 AS a)", true, "1 AS a"},
	{"create_or_alter_view", "CREATE OR ALTER VIEW v AS (SELECT 1 AS a) UNION (SELECT 2)", true, "(1 AS a UNION 2)"},
	{"view_order_after", "CREATE VIEW v AS (SELECT 1 AS a) ORDER BY a", false, ""},
	{"view_order_after_top", "CREATE VIEW v AS (SELECT TOP 1 a FROM t) ORDER BY a", false, ""},
	{"view_for_after", "CREATE VIEW v AS (SELECT 1 AS a) FOR XML PATH", false, ""},
	{"view_option_after", "CREATE VIEW v AS (SELECT 1 AS a) OPTION (MAXDOP 1)", false, ""},

	// A cursor's query is a statement's.
	{"cursor_paren", "DECLARE c CURSOR FOR (SELECT 1)", true, "1"},
	{"cursor_union", "DECLARE c CURSOR FOR (SELECT 1) UNION (SELECT 2)", true, "(1 UNION 2)"},
	{"cursor_options_union", "DECLARE c CURSOR LOCAL FAST_FORWARD FOR (SELECT 1) UNION (SELECT 2)", true, "(1 UNION 2)"},
	{"cursor_order", "DECLARE c CURSOR FOR (SELECT 1) ORDER BY 1", true, "1 ORDER BY"},
	{"cursor_option", "DECLARE c CURSOR FOR (SELECT 1) OPTION (MAXDOP 1)", true, "1 OPTION"},
	{"cursor_for_update", "DECLARE c CURSOR FOR (SELECT a FROM t) FOR UPDATE OF a", true, "a"},
	{"cursor_cte_paren", "DECLARE c CURSOR FOR WITH x AS (SELECT 1 a) (SELECT a FROM x)", true, "WITH[1 a] a"},
	{"cursor_top_order_inside", "DECLARE c CURSOR FOR (SELECT TOP 1 a FROM t ORDER BY a)", false, ""},
	{"cursor_into_inside", "DECLARE c CURSOR FOR (SELECT 1 INTO #t)", false, ""},
	{"cursor_into_operand", "DECLARE c CURSOR FOR (SELECT 1) UNION (SELECT 2 INTO #t)", false, ""},

	// An inline table-valued function's RETURN is nested; RETURN (WITH ...)
	// puts the CTE list inside the parentheses.
	{"tvf_paren", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT TOP 1 a FROM t ORDER BY a)", true, "a ORDER BY"},
	{"tvf_paren_double", "CREATE FUNCTION f() RETURNS TABLE AS RETURN ((SELECT TOP 1 a FROM t ORDER BY a))", true, "a ORDER BY"},
	{"tvf_union", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT 1 AS a) UNION (SELECT 2)", true, "(1 AS a UNION 2)"},
	{"tvf_union_top_order", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT TOP 1 a FROM t ORDER BY a) UNION (SELECT 2)", true, "(a ORDER BY UNION 2)"},
	{"tvf_grouped", "CREATE FUNCTION f() RETURNS TABLE AS RETURN ((SELECT 1 AS a) UNION (SELECT 2))", true, "(1 AS a UNION 2)"},
	{"tvf_bare_then_paren", "CREATE FUNCTION f() RETURNS TABLE AS RETURN SELECT 1 AS a UNION (SELECT 2)", true, "(1 AS a UNION 2)"},
	{"tvf_paren_bare_then_paren", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT 1 AS a UNION (SELECT 2))", true, "(1 AS a UNION 2)"},
	{"tvf_cte_in_parens", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (WITH c AS (SELECT 1 a) SELECT a FROM c)", true, "WITH[1 a] a"},
	{"tvf_cte_in_parens_paren_query", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (WITH c AS (SELECT 1 a) (SELECT a FROM c))", true, "WITH[1 a] a"},
	{"tvf_cte_paren_query", "CREATE FUNCTION f() RETURNS TABLE AS RETURN WITH c AS (SELECT 1 a) (SELECT a FROM c)", true, "WITH[1 a] a"},
	{"tvf_order_after", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT 1 AS a) ORDER BY a", false, ""},
	{"tvf_cte_in_parens_union", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (WITH c AS (SELECT 1 a) SELECT a FROM c) UNION (SELECT 2)", false, ""},
	{"tvf_cte_in_double_parens", "CREATE FUNCTION f() RETURNS TABLE AS RETURN ((WITH c AS (SELECT 1 a) SELECT a FROM c))", false, ""},
	{"tvf_unclosed", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT 1 AS a", false, ""},
	{"tvf_alias", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (SELECT 1 AS a) x", false, ""},
	{"tvf_expression", "CREATE FUNCTION f() RETURNS TABLE AS RETURN (1)", false, ""},
	{"tvf_missing_query", "CREATE FUNCTION f() RETURNS TABLE AS RETURN", false, ""},
}

var setOpNames = map[ast.SetOperation]string{
	ast.SetOpUnion:     "UNION",
	ast.SetOpIntersect: "INTERSECT",
	ast.SetOpExcept:    "EXCEPT",
}

// queryShape prints the query expressions in n, separated by ", ". A set
// operation prints as (L OP R) and a query specification as its select list,
// or as _ when the select list holds a subquery. The clauses a SelectStmt
// carries follow it (INTO, ORDER BY, OFFSET, FOR, OPTION), its CTE list
// precedes it as WITH[...], and the query expressions nested in a query
// specification follow its select list in braces, an IN subquery marked IN.
func queryShape(sql string, n ast.Node) string {
	var parts []string
	ast.Inspect(n, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SelectStmt:
			parts = append(parts, selectShape(sql, v))
			return false
		case *ast.InExpr:
			if v.Subquery == nil {
				return true
			}
			if left := queryShape(sql, v.Expr); left != "" {
				parts = append(parts, left)
			}
			parts = append(parts, "IN "+queryShape(sql, v.Subquery))
			return false
		}
		return true
	})
	return strings.Join(parts, ", ")
}

func selectShape(sql string, s *ast.SelectStmt) string {
	var b strings.Builder
	if s.WithClause != nil {
		fmt.Fprintf(&b, "WITH[%s] ", queryShape(sql, s.WithClause))
	}
	if s.Op != ast.SetOpNone {
		op := setOpNames[s.Op]
		if s.All {
			op += " ALL"
		}
		fmt.Fprintf(&b, "(%s %s %s)", selectShape(sql, s.Larg), op, selectShape(sql, s.Rarg))
	} else {
		label := locSliceForMessage(sql, ast.NodeLoc(s.TargetList))
		var nested []string
		if inTargets := queryShape(sql, s.TargetList); inTargets != "" {
			label = "_"
			nested = append(nested, inTargets)
		}
		for _, clause := range []ast.Node{s.FromClause, s.WhereClause, s.HavingClause} {
			if shape := queryShape(sql, clause); shape != "" {
				nested = append(nested, shape)
			}
		}
		b.WriteString(label)
		if len(nested) > 0 {
			fmt.Fprintf(&b, " {%s}", strings.Join(nested, ", "))
		}
	}
	for _, clause := range []struct {
		present bool
		name    string
	}{
		{s.IntoTable != nil, "INTO"},
		{s.OrderByClause != nil, "ORDER BY"},
		{s.OffsetClause != nil, "OFFSET"},
		{s.ForClause != nil, "FOR"},
		{s.OptionClause != nil, "OPTION"},
	} {
		if clause.present {
			b.WriteString(" " + clause.name)
		}
	}
	return b.String()
}

// TestParenQuery checks that omni accepts the parenthesized query
// expressions SQL Server accepts, builds the query they denote, and keeps
// every Loc in each query inside its parent. The check covers the queries
// only: an INSERT target's table hints and an inline function's RETURNS
// TABLE carry Locs outside their parents whatever the query is.
func TestParenQuery(t *testing.T) {
	for _, tc := range parenQueryCases {
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
			if got := queryShape(tc.sql, list); got != tc.shape {
				t.Errorf("query shape = %s, want %s", got, tc.shape)
			}
			ast.Inspect(list, func(n ast.Node) bool {
				if stmt, ok := n.(*ast.SelectStmt); ok {
					assertChildLocsWithinParent(t, tc.sql, stmt)
					return false
				}
				return true
			})
		})
	}
}

// TestParenQueryOracle checks accept/reject parity with SQL Server for every
// case in parenQueryCases.
func TestParenQueryOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range parenQueryCases {
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

// TestParenQueryLoc checks the span of each SelectStmt: a parenthesized query
// expression spans its parentheses and the clauses that follow them, and a
// set operation spans both operands.
func TestParenQueryLoc(t *testing.T) {
	tests := []struct {
		sql   string
		index int // SelectStmt in pre-order
		want  string
	}{
		{"(SELECT 1)", 0, "(SELECT 1)"},
		{"((SELECT 1))", 0, "((SELECT 1))"},
		{"(SELECT 1) ORDER BY 1 OPTION (MAXDOP 1)", 0, "(SELECT 1) ORDER BY 1 OPTION (MAXDOP 1)"},
		{"(SELECT 1) UNION (SELECT 2)", 0, "(SELECT 1) UNION (SELECT 2)"},
		{"(SELECT 1) UNION (SELECT 2)", 1, "(SELECT 1)"},
		{"(SELECT 1) UNION (SELECT 2) ORDER BY 1", 2, "(SELECT 2) ORDER BY 1"},
		{"((SELECT 1) UNION (SELECT 2)) ORDER BY 1", 0, "((SELECT 1) UNION (SELECT 2)) ORDER BY 1"},
		{"SELECT 1 UNION SELECT 2", 0, "SELECT 1 UNION SELECT 2"},
		{"SELECT 1 UNION SELECT 2", 1, "SELECT 1"},
		{"SELECT 1 UNION SELECT 2 UNION SELECT 3", 2, "SELECT 2 UNION SELECT 3"},
		{"WITH c AS (SELECT 1 a) (SELECT * FROM c) UNION (SELECT 2)", 0, "WITH c AS (SELECT 1 a) (SELECT * FROM c) UNION (SELECT 2)"},
		{"WITH c AS (SELECT 1 a) (SELECT * FROM c) UNION (SELECT 2)", 1, "WITH c AS (SELECT 1 a) (SELECT * FROM c)"},
		{"SELECT * FROM ((SELECT 1 AS a) UNION (SELECT 2)) x", 1, "(SELECT 1 AS a) UNION (SELECT 2)"},
		{"SELECT * FROM (((SELECT 1 AS a))) x", 1, "((SELECT 1 AS a))"},
		{"INSERT INTO t (a) (SELECT 1) ORDER BY 1", 0, "(SELECT 1) ORDER BY 1"},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s#%d", tc.sql, tc.index), func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			stmt := nthNodeByType(list, "SelectStmt", tc.index)
			if stmt == nil {
				t.Fatalf("no SelectStmt #%d in %s", tc.index, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(stmt)); got != tc.want {
				t.Errorf("SelectStmt #%d text = %q, want %q", tc.index, got, tc.want)
			}
		})
	}
}

// TestParenQueryErrorPosition checks that a clause in the wrong place is
// reported at its keyword.
func TestParenQueryErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{"(SELECT TOP 1 a FROM t ORDER BY a)", `syntax error at or near "ORDER"`},
		{"(SELECT 1) ORDER BY 1 UNION (SELECT 2)", `syntax error at or near "UNION"`},
		{"SELECT * FROM ((SELECT 1 AS a) ORDER BY a) x", `syntax error at or near "ORDER"`},
		{"INSERT INTO t (a) (SELECT 1 INTO #x)", `syntax error at or near "INTO"`},
		{"(SELECT 1) UNION (1)", `syntax error at or near "1"`},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want %s", tc.sql, tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Parse(%q) error = %q, want %s", tc.sql, err.Error(), tc.want)
			}
		})
	}
}

// TestParenQueryCompletion checks the candidates offered after a '(' that
// opens a query expression and after one that closes it.
func TestParenQueryCompletion(t *testing.T) {
	tests := []struct {
		sql       string
		wantToks  []int
		wantRules []string
	}{
		{"(", []int{kwSELECT}, nil},
		{"SELECT 1 UNION (", []int{kwSELECT}, nil},
		{"WITH c AS (SELECT 1 a) (", []int{kwSELECT}, nil},
		{"INSERT INTO t (a) (", []int{kwSELECT}, nil},
		{"SELECT * FROM t WHERE EXISTS ((", []int{kwSELECT}, nil},
		{"(SELECT 1) ORDER BY ", nil, []string{"columnref"}},
	}
	for _, tc := range tests {
		t.Run(tc.sql+"|", func(t *testing.T) {
			cs := Collect(tc.sql, len(tc.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, tok := range tc.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s", TokenName(tok))
				}
			}
			for _, r := range tc.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
		})
	}
}

// TestParenQueryStaysLinear guards parenIsQuery: the parser asks whether a
// '(' opens a query expression at every nesting level as it descends, and
// answering by rescanning would make n nested parens cost O(n²). Each shape
// below takes milliseconds when linear and minutes otherwise.
func TestParenQueryStaysLinear(t *testing.T) {
	const depth = 10000
	nest := func(open, mid, closing string) string {
		return strings.Repeat(open, depth) + mid + strings.Repeat(closing, depth)
	}
	for _, tc := range []struct {
		name string
		sql  string
	}{
		{"statement", nest("(", "SELECT 1", ")")},
		{"statement set operations", nest("(", "SELECT 1", ") UNION (SELECT 2)")},
		{"derived table", "SELECT * FROM " + nest("(", "SELECT 1 AS a", ")") + " x"},
		{"derived table set operations", "SELECT * FROM (" + nest("(", "SELECT 1 AS a", ") UNION (SELECT 2)") + ") x"},
		{"joined table", "SELECT * FROM " + nest("(", "(SELECT 1 AS a) s JOIN t ON 1 = 1", ")")},
		{"scalar subquery", "SELECT " + nest("(", "SELECT 1", ")")},
		{"scalar set operations", "SELECT (" + nest("(", "SELECT 1", ") UNION (SELECT 2)") + ")"},
		{"expression", "SELECT " + nest("(", "(SELECT 1) + 1", ")")},
		{"IN subquery", "SELECT * FROM t WHERE a IN " + nest("(", "SELECT 1", ")")},
	} {
		start := time.Now()
		if _, err := Parse(tc.sql); err != nil {
			t.Fatalf("%s: Parse failed: %v", tc.name, err)
		}
		if elapsed := time.Since(start); elapsed > 2*time.Second {
			t.Errorf("%s: Parse through %d nested parens took %v", tc.name, depth, elapsed)
		}
	}
}

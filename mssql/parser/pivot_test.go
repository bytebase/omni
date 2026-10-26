package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// pivotAggregateNames are the unqualified aggregates SQL Server 2022 takes in
// a PIVOT. pivotBodyCases accepts each of them and rejects the near misses.
var pivotAggregateNames = []string{
	"APPROX_COUNT_DISTINCT", "APPROX_PERCENTILE_CONT", "APPROX_PERCENTILE_DISC",
	"AVG", "CHECKSUM_AGG", "COUNT", "COUNT_BIG", "MAX", "MIN", "STDEV", "STDEVP",
	"STRING_AGG", "SUM", "VAR", "VARP",
}

// pivotBodyCases covers the body of PIVOT and UNPIVOT. Every part is
// required; every column is a name of one to four parts; the PIVOT
// aggregate is a built-in aggregate by its undelimited name, or a qualified
// user-defined one, over columns only; and the alias may be WINDOW. shape is
// the outermost PivotExpr or UnpivotExpr of an accepted statement, printed by
// pivotShape.
var pivotBodyCases = func() []struct {
	name    string
	sql     string
	accepts bool
	shape   string
} {
	cases := []struct {
		name    string
		sql     string
		accepts bool
		shape   string
	}{
		// The aggregate.
		{"agg_lower_case", "SELECT * FROM t PIVOT (max(col) FOR a IN ([1])) p", true, "PIVOT max(col) FOR a IN (1) AS p"},
		{"agg_mixed_case", "SELECT * FROM t PIVOT (Sum(col) FOR a IN ([1])) p", true, "PIVOT Sum(col) FOR a IN (1) AS p"},
		{"agg_space_before_paren", "SELECT * FROM t PIVOT (MAX (col) FOR a IN ([1])) p", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"agg_schema_qualified", "SELECT * FROM t PIVOT (dbo.agg(col) FOR a IN ([1])) p", true, "PIVOT dbo.agg(col) FOR a IN (1) AS p"},
		{"agg_database_qualified", "SELECT * FROM t PIVOT (db.dbo.agg(col) FOR a IN ([1])) p", true, "PIVOT db.dbo.agg(col) FOR a IN (1) AS p"},
		{"agg_qualified_delimited", "SELECT * FROM t PIVOT ([dbo].[agg](col) FOR a IN ([1])) p", true, "PIVOT dbo.agg(col) FOR a IN (1) AS p"},
		{"agg_qualified_builtin_name", "SELECT * FROM t PIVOT ([dbo].[MAX](col) FOR a IN ([1])) p", true, "PIVOT dbo.MAX(col) FOR a IN (1) AS p"},
		{"agg_qualified_window", "SELECT * FROM t PIVOT (dbo.window(col) FOR a IN ([1])) p", true, "PIVOT dbo.window(col) FOR a IN (1) AS p"},
		{"agg_two_columns", "SELECT * FROM t PIVOT (MAX(col, a) FOR a IN ([1])) p", true, "PIVOT MAX(col, a) FOR a IN (1) AS p"},
		{"agg_string_agg_columns", "SELECT * FROM t PIVOT (STRING_AGG(col, a) FOR a IN ([1])) p", true, "PIVOT STRING_AGG(col, a) FOR a IN (1) AS p"},
		{"agg_qualified_column", "SELECT * FROM t PIVOT (MAX(t.col) FOR a IN ([1])) p", true, "PIVOT MAX(t.col) FOR a IN (1) AS p"},
		{"agg_four_part_column", "SELECT * FROM t PIVOT (MAX(db.dbo.t.col) FOR a IN ([1])) p", true, "PIVOT MAX(db.dbo.t.col) FOR a IN (1) AS p"},
		{"agg_delimited_column", "SELECT * FROM t PIVOT (MAX([col]) FOR a IN ([1])) p", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"agg_window_column", "SELECT * FROM t PIVOT (MAX(window) FOR a IN ([1])) p", true, "PIVOT MAX(window) FOR a IN (1) AS p"},
		{"agg_star", "SELECT * FROM t PIVOT (COUNT(*) FOR a IN ([1])) p", false, ""},
		{"agg_qualified_star", "SELECT * FROM t PIVOT (MAX(t.*) FOR a IN ([1])) p", false, ""},
		{"agg_expression", "SELECT * FROM t PIVOT (MAX(col + 1) FOR a IN ([1])) p", false, ""},
		{"agg_distinct", "SELECT * FROM t PIVOT (MAX(DISTINCT col) FOR a IN ([1])) p", false, ""},
		{"agg_all", "SELECT * FROM t PIVOT (MAX(ALL col) FOR a IN ([1])) p", false, ""},
		{"agg_literal", "SELECT * FROM t PIVOT (MAX(1) FOR a IN ([1])) p", false, ""},
		{"agg_variable", "SELECT * FROM t PIVOT (MAX(@v) FOR a IN ([1])) p", false, ""},
		{"agg_paren_column", "SELECT * FROM t PIVOT (MAX((col)) FOR a IN ([1])) p", false, ""},
		{"agg_paren_call", "SELECT * FROM t PIVOT ((MAX(col)) FOR a IN ([1])) p", false, ""},
		{"agg_collate", "SELECT * FROM t PIVOT (MAX(col COLLATE Latin1_General_CI_AS) FOR a IN ([1])) p", false, ""},
		{"agg_literal_second_argument", "SELECT * FROM t PIVOT (STRING_AGG(col, ',') FOR a IN ([1])) p", false, ""},
		{"agg_no_arguments", "SELECT * FROM t PIVOT (MAX() FOR a IN ([1])) p", false, ""},
		{"agg_not_a_call", "SELECT * FROM t PIVOT (col FOR a IN ([1])) p", false, ""},
		{"agg_over", "SELECT * FROM t PIVOT (MAX(col) OVER () FOR a IN ([1])) p", false, ""},
		{"agg_five_part_column", "SELECT * FROM t PIVOT (MAX(x.db.dbo.t.col) FOR a IN ([1])) p", false, ""},
		{"agg_empty_column_part", "SELECT * FROM t PIVOT (MAX(db..col) FOR a IN ([1])) p", false, ""},
		{"agg_leading_dot_column", "SELECT * FROM t PIVOT (MAX(.col) FOR a IN ([1])) p", false, ""},
		{"agg_identitycol", "SELECT * FROM t PIVOT (MAX(IDENTITYCOL) FOR a IN ([1])) p", false, ""},
		{"agg_pseudo_column", "SELECT * FROM t PIVOT (MAX($IDENTITY) FOR a IN ([1])) p", false, ""},
		{"agg_action_pseudo_column", "SELECT * FROM t PIVOT (MAX($action) FOR a IN ([1])) p", false, ""},
		{"agg_bracketed_builtin", "SELECT * FROM t PIVOT ([MAX](col) FOR a IN ([1])) p", false, ""},
		{"agg_quoted_builtin", `SELECT * FROM t PIVOT ("MAX"(col) FOR a IN ([1])) p`, false, ""},
		{"agg_unknown", "SELECT * FROM t PIVOT (myagg(col) FOR a IN ([1])) p", false, ""},
		{"agg_grouping", "SELECT * FROM t PIVOT (GROUPING(col) FOR a IN ([1])) p", false, ""},
		{"agg_grouping_id", "SELECT * FROM t PIVOT (GROUPING_ID(col, a) FOR a IN ([1])) p", false, ""},
		{"agg_scalar_function", "SELECT * FROM t PIVOT (ABS(col) FOR a IN ([1])) p", false, ""},
		{"agg_window_function", "SELECT * FROM t PIVOT (ROW_NUMBER(col) FOR a IN ([1])) p", false, ""},
		{"agg_keyword_function", "SELECT * FROM t PIVOT (LEFT(col) FOR a IN ([1])) p", false, ""},
		{"agg_coalesce", "SELECT * FROM t PIVOT (COALESCE(col) FOR a IN ([1])) p", false, ""},
		{"agg_four_part_name", "SELECT * FROM t PIVOT (srv.db.dbo.agg(col) FOR a IN ([1])) p", false, ""},
		{"agg_empty_name_part", "SELECT * FROM t PIVOT (db..agg(col) FOR a IN ([1])) p", false, ""},
		{"agg_leading_dot_name", "SELECT * FROM t PIVOT (.agg(col) FOR a IN ([1])) p", false, ""},
		{"agg_reserved_name_part", "SELECT * FROM t PIVOT (dbo.select(col) FOR a IN ([1])) p", false, ""},

		// FOR and the pivot column.
		{"for_qualified", "SELECT * FROM t PIVOT (MAX(col) FOR t.a IN ([1])) p", true, "PIVOT MAX(col) FOR t.a IN (1) AS p"},
		{"for_four_part", "SELECT * FROM t PIVOT (MAX(col) FOR db.dbo.t.a IN ([1])) p", true, "PIVOT MAX(col) FOR db.dbo.t.a IN (1) AS p"},
		{"for_delimited", "SELECT * FROM t PIVOT (MAX(col) FOR [a] IN ([1])) p", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"for_window", "SELECT * FROM t PIVOT (MAX(col) FOR window IN ([1])) p", true, "PIVOT MAX(col) FOR window IN (1) AS p"},
		{"for_missing", "SELECT * FROM t PIVOT (MAX(col) IN ([1])) p", false, ""},
		{"for_without_column", "SELECT * FROM t PIVOT (MAX(col) FOR IN ([1])) p", false, ""},
		{"for_five_part", "SELECT * FROM t PIVOT (MAX(col) FOR x.db.dbo.t.a IN ([1])) p", false, ""},
		{"for_empty_part", "SELECT * FROM t PIVOT (MAX(col) FOR db..a IN ([1])) p", false, ""},
		{"for_trailing_dot", "SELECT * FROM t PIVOT (MAX(col) FOR a. IN ([1])) p", false, ""},
		{"for_parenthesized", "SELECT * FROM t PIVOT (MAX(col) FOR (a) IN ([1])) p", false, ""},
		{"for_pseudo_column", "SELECT * FROM t PIVOT (MAX(col) FOR $IDENTITY IN ([1])) p", false, ""},

		// IN and the output columns. SQL Server parses a prefix here and
		// rejects it only when it binds the query.
		{"in_two", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1], [2])) p", true, "PIVOT MAX(col) FOR a IN (1, 2) AS p"},
		{"in_regular", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (x, y)) p", true, "PIVOT MAX(col) FOR a IN (x, y) AS p"},
		{"in_quoted", `SELECT * FROM t PIVOT (MAX(col) FOR a IN ("x")) p`, true, "PIVOT MAX(col) FOR a IN (x) AS p"},
		{"in_duplicate", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1], [1])) p", true, "PIVOT MAX(col) FOR a IN (1, 1) AS p"},
		{"in_window", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (window)) p", true, "PIVOT MAX(col) FOR a IN (window) AS p"},
		{"in_qualified", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (t.x)) p", true, "PIVOT MAX(col) FOR a IN (t.x) AS p"},
		{"in_four_part", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (db.dbo.t.x)) p", true, "PIVOT MAX(col) FOR a IN (db.dbo.t.x) AS p"},
		{"in_space_before_close", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]) ) p", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"in_missing", "SELECT * FROM t PIVOT (MAX(col) FOR a) p", false, ""},
		{"in_without_list", "SELECT * FROM t PIVOT (MAX(col) FOR a IN) p", false, ""},
		{"in_unparenthesized", "SELECT * FROM t PIVOT (MAX(col) FOR a IN [1]) p", false, ""},
		{"in_empty", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ()) p", false, ""},
		{"in_trailing_comma", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1],)) p", false, ""},
		{"in_missing_comma", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1] [2])) p", false, ""},
		{"in_number", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (1)) p", false, ""},
		{"in_string", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ('x')) p", false, ""},
		{"in_null", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (NULL)) p", false, ""},
		{"in_variable", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (@v)) p", false, ""},
		{"in_reserved", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (select)) p", false, ""},
		{"in_alias", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1] AS x)) p", false, ""},
		{"in_five_part", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (y.db.dbo.t.x)) p", false, ""},
		{"in_empty_part", "SELECT * FROM t PIVOT (MAX(col) FOR a IN (db..x)) p", false, ""},
		{"in_extra_token", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]) extra) p", false, ""},

		// The alias.
		{"alias_as", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS p", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"alias_bracketed", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS [p]", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"alias_quoted", `SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) "p"`, true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"alias_window", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) window", true, "PIVOT MAX(col) FOR a IN (1) AS window"},
		{"alias_as_window", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS window", true, "PIVOT MAX(col) FOR a IN (1) AS window"},
		{"alias_then_where", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS p WHERE p.a = 1", true, "PIVOT MAX(col) FOR a IN (1) AS p"},
		{"alias_missing", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]))", false, ""},
		{"alias_missing_before_where", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) WHERE 1 = 1", false, ""},
		{"alias_missing_before_order_by", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) ORDER BY 1", false, ""},
		{"alias_as_only", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS", false, ""},
		{"alias_reserved", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS select", false, ""},
		{"alias_string", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS 'p'", false, ""},
		{"alias_bare_string", "SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) 'p'", false, ""},

		// UNPIVOT. The value and pivot columns are new, but SQL Server parses
		// a prefix there and rejects it only when it binds the query.
		{"unpivot_basic", "SELECT * FROM t UNPIVOT (v FOR c IN (a, col)) u", true, "UNPIVOT v FOR c IN (a, col) AS u"},
		{"unpivot_as_window", "SELECT * FROM t UNPIVOT (v FOR c IN (a)) AS window", true, "UNPIVOT v FOR c IN (a) AS window"},
		{"unpivot_delimited", "SELECT * FROM t UNPIVOT ([v] FOR [c] IN ([a])) u", true, "UNPIVOT v FOR c IN (a) AS u"},
		{"unpivot_window_value", "SELECT * FROM t UNPIVOT (window FOR c IN (a)) u", true, "UNPIVOT window FOR c IN (a) AS u"},
		{"unpivot_window_in", "SELECT * FROM t UNPIVOT (v FOR c IN (window)) u", true, "UNPIVOT v FOR c IN (window) AS u"},
		{"unpivot_qualified_value", "SELECT * FROM t UNPIVOT (t.v FOR c IN (a)) u", true, "UNPIVOT t.v FOR c IN (a) AS u"},
		{"unpivot_qualified_for", "SELECT * FROM t UNPIVOT (v FOR t.c IN (a)) u", true, "UNPIVOT v FOR t.c IN (a) AS u"},
		{"unpivot_qualified_in", "SELECT * FROM t UNPIVOT (v FOR c IN (t.a, t.col)) u", true, "UNPIVOT v FOR c IN (t.a, t.col) AS u"},
		{"unpivot_four_part_in", "SELECT * FROM t UNPIVOT (v FOR c IN (db.dbo.t.a)) u", true, "UNPIVOT v FOR c IN (db.dbo.t.a) AS u"},
		{"unpivot_missing_value", "SELECT * FROM t UNPIVOT (FOR c IN (a)) u", false, ""},
		{"unpivot_missing_for", "SELECT * FROM t UNPIVOT (v c IN (a)) u", false, ""},
		{"unpivot_for_without_column", "SELECT * FROM t UNPIVOT (v FOR IN (a)) u", false, ""},
		{"unpivot_missing_in", "SELECT * FROM t UNPIVOT (v FOR c (a)) u", false, ""},
		{"unpivot_in_unparenthesized", "SELECT * FROM t UNPIVOT (v FOR c IN a) u", false, ""},
		{"unpivot_in_empty", "SELECT * FROM t UNPIVOT (v FOR c IN ()) u", false, ""},
		{"unpivot_in_string", "SELECT * FROM t UNPIVOT (v FOR c IN ('a')) u", false, ""},
		{"unpivot_in_number", "SELECT * FROM t UNPIVOT (v FOR c IN (1)) u", false, ""},
		{"unpivot_in_alias", "SELECT * FROM t UNPIVOT (v FOR c IN (a AS x)) u", false, ""},
		{"unpivot_in_five_part", "SELECT * FROM t UNPIVOT (v FOR c IN (x.db.dbo.t.a)) u", false, ""},
		{"unpivot_in_empty_part", "SELECT * FROM t UNPIVOT (v FOR c IN (db..a)) u", false, ""},
		{"unpivot_value_five_part", "SELECT * FROM t UNPIVOT (x.db.dbo.t.v FOR c IN (a)) u", false, ""},
		{"unpivot_for_five_part", "SELECT * FROM t UNPIVOT (v FOR x.db.dbo.t.c IN (a)) u", false, ""},
		{"unpivot_value_call", "SELECT * FROM t UNPIVOT (MAX(v) FOR c IN (a)) u", false, ""},
		{"unpivot_value_parenthesized", "SELECT * FROM t UNPIVOT ((v) FOR c IN (a)) u", false, ""},
		{"unpivot_value_variable", "SELECT * FROM t UNPIVOT (@v FOR c IN (a)) u", false, ""},
		{"unpivot_alias_missing", "SELECT * FROM t UNPIVOT (v FOR c IN (a))", false, ""},
		{"unpivot_alias_as_only", "SELECT * FROM t UNPIVOT (v FOR c IN (a)) AS", false, ""},
	}
	for _, name := range pivotAggregateNames {
		cases = append(cases, struct {
			name    string
			sql     string
			accepts bool
			shape   string
		}{
			"agg_builtin_" + strings.ToLower(name),
			"SELECT * FROM t PIVOT (" + name + "(col) FOR a IN ([1])) p",
			true,
			"PIVOT " + name + "(col) FOR a IN (1) AS p",
		})
	}
	return cases
}()

// pivotColumnText prints a PIVOT or UNPIVOT column as its dotted name.
func pivotColumnText(ref *ast.ColumnRef) string {
	var parts []string
	for _, part := range []string{ref.Server, ref.Database, ref.Schema, ref.Table, ref.Column} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

// pivotColumnListText prints a list of PIVOT or UNPIVOT columns.
func pivotColumnListText(list *ast.List) string {
	var cols []string
	for _, item := range list.Items {
		cols = append(cols, pivotColumnText(item.(*ast.ColumnRef)))
	}
	return "(" + strings.Join(cols, ", ") + ")"
}

// pivotShape prints the body of a PivotExpr or UnpivotExpr.
func pivotShape(t *testing.T, n ast.Node) string {
	t.Helper()
	switch v := n.(type) {
	case *ast.PivotExpr:
		agg, ok := v.AggFunc.(*ast.FuncCallExpr)
		if !ok {
			t.Fatalf("PivotExpr.AggFunc is %T, want *ast.FuncCallExpr", v.AggFunc)
		}
		name := agg.Name.Object
		if agg.Name.Schema != "" {
			name = agg.Name.Schema + "." + name
		}
		if agg.Name.Database != "" {
			name = agg.Name.Database + "." + name
		}
		return "PIVOT " + name + pivotColumnListText(agg.Args) +
			" FOR " + pivotColumnText(v.ForCol) +
			" IN " + pivotColumnListText(v.InValues) + " AS " + v.Alias
	case *ast.UnpivotExpr:
		return "UNPIVOT " + pivotColumnText(v.ValueCol) +
			" FOR " + pivotColumnText(v.ForCol) +
			" IN " + pivotColumnListText(v.InCols) + " AS " + v.Alias
	}
	t.Fatalf("pivotShape(%T)", n)
	return ""
}

// TestPivotBody checks that omni accepts the PIVOT and UNPIVOT bodies SQL
// Server accepts, rejects the others, and keeps every Loc inside its parent.
func TestPivotBody(t *testing.T) {
	for _, tc := range pivotBodyCases {
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
			n := outermostJoin(list)
			if n == nil {
				t.Fatalf("no PivotExpr or UnpivotExpr in %s", ast.NodeToString(list))
			}
			if got := pivotShape(t, n); got != tc.shape {
				t.Errorf("shape = %s, want %s", got, tc.shape)
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestPivotBodyOracle checks accept/reject parity with SQL Server for every
// case in pivotBodyCases.
func TestPivotBodyOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range pivotBodyCases {
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

// TestPivotNameParts checks how a multi-part name fills a ColumnRef and the
// aggregate's TableRef.
func TestPivotNameParts(t *testing.T) {
	list, err := Parse("SELECT * FROM t PIVOT (db.dbo.agg(db.dbo.t.col) FOR s.t.a IN (t.x)) p")
	if err != nil {
		t.Fatal(err)
	}
	pivot := nthNodeByType(list, "PivotExpr", 0).(*ast.PivotExpr)
	agg := pivot.AggFunc.(*ast.FuncCallExpr)
	if got := *agg.Name; got.Database != "db" || got.Schema != "dbo" || got.Object != "agg" {
		t.Errorf("aggregate name = %+v, want db.dbo.agg", got)
	}
	arg := agg.Args.Items[0].(*ast.ColumnRef)
	if arg.Database != "db" || arg.Schema != "dbo" || arg.Table != "t" || arg.Column != "col" {
		t.Errorf("aggregate argument = %+v, want db.dbo.t.col", *arg)
	}
	if f := pivot.ForCol; f.Schema != "s" || f.Table != "t" || f.Column != "a" {
		t.Errorf("FOR column = %+v, want s.t.a", *f)
	}
	if in := pivot.InValues.Items[0].(*ast.ColumnRef); in.Table != "t" || in.Column != "x" {
		t.Errorf("IN column = %+v, want t.x", *in)
	}
}

// TestPivotBodyLoc checks the span of each part of a PIVOT and UNPIVOT body.
func TestPivotBodyLoc(t *testing.T) {
	tests := []struct {
		sql   string
		node  string
		index int // node of that type in pre-order
		want  string
	}{
		{"SELECT * FROM t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p WHERE 1 = 1", "PivotExpr", 0, "t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p"},
		{"SELECT * FROM t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p WHERE 1 = 1", "FuncCallExpr", 0, "dbo.agg(t.col)"},
		{"SELECT * FROM t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p WHERE 1 = 1", "ColumnRef", 0, "t.col"},
		{"SELECT * FROM t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p WHERE 1 = 1", "ColumnRef", 1, "db.s.t.a"},
		{"SELECT * FROM t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p WHERE 1 = 1", "ColumnRef", 2, "[1]"},
		{"SELECT * FROM t PIVOT (dbo.agg(t.col) FOR db.s.t.a IN ([1], x)) AS p WHERE 1 = 1", "ColumnRef", 3, "x"},
		{"SELECT * FROM t UNPIVOT ([v] FOR t.c IN (a, db.s.t.b)) u", "UnpivotExpr", 0, "t UNPIVOT ([v] FOR t.c IN (a, db.s.t.b)) u"},
		{"SELECT * FROM t UNPIVOT ([v] FOR t.c IN (a, db.s.t.b)) u", "ColumnRef", 0, "[v]"},
		{"SELECT * FROM t UNPIVOT ([v] FOR t.c IN (a, db.s.t.b)) u", "ColumnRef", 1, "t.c"},
		{"SELECT * FROM t UNPIVOT ([v] FOR t.c IN (a, db.s.t.b)) u", "ColumnRef", 3, "db.s.t.b"},
	}
	for _, tc := range tests {
		t.Run(tc.node+"/"+tc.want, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			n := nthNodeByType(list, tc.node, tc.index)
			if n == nil {
				t.Fatalf("no %s #%d in %s", tc.node, tc.index, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(n)); got != tc.want {
				t.Errorf("%s #%d text = %q, want %q", tc.node, tc.index, got, tc.want)
			}
		})
	}
}

// TestPivotBodyErrorPosition checks where a malformed PIVOT or UNPIVOT body
// is reported, and the message an unknown aggregate carries.
func TestPivotBodyErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		want string // Error()
		msg  string // Message, when it says more than Error()
		pos  int
	}{
		{"SELECT * FROM t PIVOT (myagg(col) FOR a IN ([1])) p", `syntax error at or near "myagg"`, "'myagg' is not a recognized aggregate function", 23},
		{"SELECT * FROM t PIVOT ([MAX](col) FOR a IN ([1])) p", `syntax error at or near "["`, "", 23},
		{"SELECT * FROM t PIVOT (COUNT(*) FOR a IN ([1])) p", `syntax error at or near "*"`, "", 29},
		{"SELECT * FROM t PIVOT (MAX(col) IN ([1])) p", `syntax error at or near "IN"`, "", 32},
		{"SELECT * FROM t PIVOT (MAX(col) FOR a) p", `syntax error at or near ")"`, "", 37},
		{"SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]))", "syntax error at end of input", "", 47},
		{"SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) WHERE 1 = 1", `syntax error at or near "WHERE"`, "", 48},
		{"SELECT * FROM t UNPIVOT (v FOR c IN (a)) AS", "syntax error at end of input", "", 43},
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
			if tc.msg != "" && pe.Message != tc.msg {
				t.Errorf("Message = %q, want %q", pe.Message, tc.msg)
			}
			if pe.Position != tc.pos {
				t.Errorf("Position = %d, want %d", pe.Position, tc.pos)
			}
		})
	}
}

// TestPivotBodyCompletion checks the candidates offered inside a PIVOT or
// UNPIVOT body.
func TestPivotBodyCompletion(t *testing.T) {
	tests := []struct {
		sql       string
		wantToks  []int
		wantRules []string
	}{
		{"SELECT * FROM t PIVOT (SUM(", nil, []string{"columnref"}},
		{"SELECT * FROM t PIVOT (dbo.", nil, []string{"func_name"}},
		{"SELECT * FROM t PIVOT (SUM(col) ", []int{kwFOR}, nil},
		{"SELECT * FROM t PIVOT (SUM(col) FOR ", nil, []string{"columnref"}},
		{"SELECT * FROM t PIVOT (SUM(col) FOR a ", []int{kwIN}, nil},
		{"SELECT * FROM t PIVOT (SUM(col) FOR a IN ([1])) ", []int{kwAS}, nil},
		{"SELECT * FROM t UNPIVOT (v ", []int{kwFOR}, nil},
		{"SELECT * FROM t UNPIVOT (v FOR c ", []int{kwIN}, nil},
		{"SELECT * FROM t UNPIVOT (v FOR c IN (a, ", nil, []string{"columnref"}},
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

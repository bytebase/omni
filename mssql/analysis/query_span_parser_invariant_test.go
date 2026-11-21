package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

func TestQuerySpanParserStructuralInvariants(t *testing.T) {
	type invariant struct {
		name  string
		sql   string
		check func(*testing.T, ast.Node)
	}

	invariants := []invariant{
		{
			name: "res_target_as_alias",
			sql:  "SELECT c AS a FROM t",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				rt := firstOf[*ast.ResTarget](t, sel.TargetList)
				if rt.Name != "a" {
					t.Errorf("want ResTarget.Name=a, got %q", rt.Name)
				}
			},
		},
		{
			name: "res_target_equals_alias",
			sql:  "SELECT a = c FROM t",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				rt := firstOf[*ast.ResTarget](t, sel.TargetList)
				if rt.Name != "a" {
					t.Errorf("ResTarget.Name: want a, got %q", rt.Name)
				}
			},
		},
		{
			name: "res_target_bare_column",
			sql:  "SELECT c FROM t",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				rt := firstOf[*ast.ResTarget](t, sel.TargetList)
				if rt.Name != "" {
					t.Errorf("ResTarget.Name: want empty, got %q", rt.Name)
				}
				if _, ok := rt.Val.(*ast.ColumnRef); !ok {
					t.Errorf("ResTarget.Val: want *ColumnRef, got %T", rt.Val)
				}
			},
		},
		{
			name: "star_unqualified",
			sql:  "SELECT * FROM t",
			check: func(t *testing.T, n ast.Node) {
				star := findFirst[*ast.StarExpr](mustSelect(t, n))
				if star == nil || star.Qualifier != "" {
					t.Errorf("got %+v", star)
				}
			},
		},
		{
			name: "star_qualified",
			sql:  "SELECT t1.* FROM t1",
			check: func(t *testing.T, n ast.Node) {
				star := findFirst[*ast.StarExpr](mustSelect(t, n))
				if star == nil || star.Qualifier != "t1" {
					t.Errorf("got %+v", star)
				}
			},
		},
		{
			name: "derived_table_column_list",
			sql:  "SELECT c1, c2 FROM (SELECT a, b FROM t) AS x(c1, c2)",
			check: func(t *testing.T, n ast.Node) {
				at := findFirst[*ast.AliasedTableRef](mustSelect(t, n).FromClause)
				if at == nil || at.Alias != "x" || at.Columns == nil || at.Columns.Len() != 2 {
					t.Errorf("got %+v", at)
				}
			},
		},
		{
			name: "cte_column_list",
			sql:  "WITH cte(c1, c2) AS (SELECT a, b FROM t) SELECT c1 FROM cte",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				c := sel.WithClause.CTEs.Items[0].(*ast.CommonTableExpr)
				if !strings.EqualFold(c.Name, "cte") || c.Columns == nil || c.Columns.Len() != 2 {
					t.Errorf("got %+v", c)
				}
			},
		},
		{
			name: "declare_table_variable_columns",
			sql:  "DECLARE @t TABLE (id INT, name NVARCHAR(50))",
			check: func(t *testing.T, n ast.Node) {
				d := n.(*ast.DeclareStmt)
				v := d.Variables.Items[0].(*ast.VariableDecl)
				if !v.IsTable || v.TableDef == nil || v.TableDef.Len() != 2 {
					t.Errorf("got %+v", v)
				}
			},
		},
		{
			name: "union_all_set_op",
			sql:  "SELECT a FROM t1 UNION ALL SELECT a FROM t2",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				if sel.Op != ast.SetOpUnion || !sel.All || sel.Larg == nil || sel.Rarg == nil {
					t.Errorf("got %+v", sel)
				}
			},
		},
		{
			name: "intersect_except",
			sql:  "SELECT a FROM t1 INTERSECT SELECT a FROM t2 EXCEPT SELECT a FROM t3",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				if sel.Op != ast.SetOpIntersect && sel.Op != ast.SetOpExcept {
					t.Errorf("want Op=Intersect or Except, got %v", sel.Op)
				}
			},
		},
		{
			name: "correlated_subquery_alias",
			sql:  "SELECT a FROM outerT AS o WHERE EXISTS (SELECT 1 FROM innerT i WHERE i.a = o.a)",
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				var foundExists, foundOuterRef bool
				ast.Inspect(sel.WhereClause, func(n ast.Node) bool {
					if _, ok := n.(*ast.ExistsExpr); ok {
						foundExists = true
					}
					if cr, ok := n.(*ast.ColumnRef); ok && strings.EqualFold(cr.Table, "o") {
						foundOuterRef = true
					}
					return true
				})
				if !foundExists || !foundOuterRef {
					t.Errorf("exists=%v outerRef=%v", foundExists, foundOuterRef)
				}
			},
		},
		{
			name: "cross_apply_tvf",
			sql:  "SELECT * FROM t1 CROSS APPLY fn(t1.a) AS x(v)",
			check: func(t *testing.T, n ast.Node) {
				j := findFirst[*ast.JoinClause](mustSelect(t, n).FromClause)
				if j == nil || j.Type != ast.JoinCrossApply {
					t.Errorf("got %+v", j)
				}
			},
		},
		{
			name: "pivot",
			sql:  "SELECT * FROM src PIVOT (SUM(amt) FOR cat IN ([a],[b])) AS p",
			check: func(t *testing.T, n ast.Node) {
				if findFirst[*ast.PivotExpr](mustSelect(t, n).FromClause) == nil {
					t.Error("no PivotExpr")
				}
			},
		},
		{
			name: "unpivot",
			sql:  "SELECT * FROM src UNPIVOT (v FOR c IN ([x],[y])) AS u",
			check: func(t *testing.T, n ast.Node) {
				if findFirst[*ast.UnpivotExpr](mustSelect(t, n).FromClause) == nil {
					t.Error("no UnpivotExpr")
				}
			},
		},
		{
			name: "values_table_constructor",
			sql:  "SELECT * FROM (VALUES (1,2),(3,4)) AS v(a,b)",
			check: func(t *testing.T, n ast.Node) {
				if findFirst[*ast.ValuesClause](mustSelect(t, n).FromClause) == nil {
					t.Error("no ValuesClause")
				}
			},
		},
		{
			name: "contains_and_freetext_predicates",
			sql:  "SELECT a FROM t WHERE CONTAINS(b, 'foo') AND FREETEXT(b, 'bar')",
			check: func(t *testing.T, n ast.Node) {
				if mustSelect(t, n).WhereClause == nil {
					t.Error("WHERE is nil")
				}
			},
		},
		{
			name: "bracketed_and_quoted_identifiers",
			sql:  `SELECT [col name], "col2" FROM [my table]`,
			check: func(t *testing.T, n ast.Node) {
				sel := mustSelect(t, n)
				col := findFirst[*ast.ColumnRef](sel)
				tr := findFirst[*ast.TableRef](sel)
				if col == nil || col.Column != "col name" || tr == nil || tr.Object != "my table" {
					t.Errorf("col=%+v tr=%+v", col, tr)
				}
			},
		},
		{
			name: "linked_server_reference",
			sql:  "SELECT c FROM srv1.db1.dbo.t",
			check: func(t *testing.T, n ast.Node) {
				tr := findFirst[*ast.TableRef](mustSelect(t, n).FromClause)
				if tr == nil || tr.Server != "srv1" || tr.Database != "db1" || tr.Schema != "dbo" || tr.Object != "t" {
					t.Errorf("got %+v", tr)
				}
			},
		},
		{
			name: "table_variable_in_from",
			sql:  "DECLARE @t TABLE(id INT); SELECT id FROM @t",
			check: func(t *testing.T, _ ast.Node) {
				stmts, err := mssql.Parse("DECLARE @t TABLE(id INT); SELECT id FROM @t")
				if err != nil || len(stmts) < 2 {
					t.Fatalf("parse: %v, len=%d", err, len(stmts))
				}
				tv := findFirst[*ast.TableVarRef](mustSelect(t, stmts[1].AST).FromClause)
				if tv == nil || tv.Name != "@t" {
					t.Errorf("got %+v", tv)
				}
			},
		},
		{
			name: "temp_table_hash_prefix",
			sql:  "SELECT c FROM #t",
			check: func(t *testing.T, n ast.Node) {
				tr := findFirst[*ast.TableRef](mustSelect(t, n).FromClause)
				if tr == nil || tr.Object != "#t" {
					t.Errorf("got %+v", tr)
				}
			},
		},
		{
			name: "xml_nodes_method",
			sql:  "SELECT T.c.value('@id', 'int') FROM @xml.nodes('/root/item') AS T(c)",
			check: func(t *testing.T, n ast.Node) {
				if findFirst[*ast.TableVarMethodCallRef](mustSelect(t, n).FromClause) == nil {
					t.Error("no TableVarMethodCallRef")
				}
			},
		},
		{
			name: "for_xml_clause",
			sql:  "SELECT a FROM t FOR XML AUTO",
			check: func(t *testing.T, n ast.Node) {
				if mustSelect(t, n).ForClause == nil {
					t.Error("ForClause is nil")
				}
			},
		},
		{
			name: "merge_output",
			sql:  "MERGE tgt USING src ON tgt.id = src.id WHEN MATCHED THEN UPDATE SET v = src.v;",
			check: func(t *testing.T, n ast.Node) {
				if _, ok := n.(*ast.MergeStmt); !ok {
					t.Errorf("want *MergeStmt, got %T", n)
				}
			},
		},
	}

	for _, p := range invariants {
		t.Run(p.name, func(t *testing.T) {
			stmts, err := mssql.Parse(p.sql)
			if err != nil {
				t.Fatalf("parse err: %v\nSQL: %s", err, p.sql)
			}
			if len(stmts) == 0 || stmts[0].AST == nil {
				t.Fatal("empty parse result")
			}
			p.check(t, stmts[0].AST)
		})
	}
}

func mustSelect(t *testing.T, n ast.Node) *ast.SelectStmt {
	t.Helper()
	sel, ok := n.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected *SelectStmt, got %T", n)
	}
	return sel
}

func findFirst[T ast.Node](root ast.Node) T {
	var zero T
	if root == nil {
		return zero
	}
	var found T
	ast.Inspect(root, func(n ast.Node) bool {
		if cast, ok := n.(T); ok {
			if any(found) == any(zero) {
				found = cast
			}
			return false
		}
		return true
	})
	return found
}

func firstOf[T ast.Node](t *testing.T, list *ast.List) T {
	t.Helper()
	var zero T
	if list == nil || list.Len() == 0 {
		t.Fatal("list empty")
	}
	cast, ok := list.Items[0].(T)
	if !ok {
		t.Fatalf("first item is %T, want %T", list.Items[0], zero)
	}
	return cast
}

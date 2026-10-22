package parser

import (
	"context"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestSortByUsingOperator checks that ORDER BY ... USING keeps its operator
// and accepts and rejects what PostgreSQL does.
//
// gram.y: sortby — a_expr USING qual_all_Op opt_nulls_order
func TestSortByUsingOperator(t *testing.T) {
	o := StartParenOracle(t)
	for _, tc := range []struct {
		sql string
		op  string // the UseOp names joined by ".", for accepted statements
	}{
		{sql: "SELECT 1 ORDER BY 1 USING >", op: ">"},
		{sql: "SELECT 1 ORDER BY 1 USING < NULLS LAST", op: "<"},
		{sql: "SELECT 1 ORDER BY 1 USING ~<~", op: "~<~"},
		{sql: "SELECT 1 ORDER BY 1 USING OPERATOR(pg_catalog.>)", op: "pg_catalog.>"},
		{sql: "SELECT array_agg(1 ORDER BY 1 USING >)", op: ">"},
		{sql: "SELECT 1 ORDER BY 1 USING"},
		{sql: "SELECT 1 ORDER BY 1 USING NULLS FIRST"},
		{sql: "SELECT 1 ORDER BY 1 USING OPERATOR(>"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			pgStatus, _, pgMsg := classifyPG(context.Background(), o, tc.sql)
			list, err := Parse(tc.sql)
			if (pgStatus == PGAccept) != (err == nil) {
				t.Fatalf("omni error %v, postgres %v %s", err, pgStatus, pgMsg)
			}
			if err != nil {
				return
			}
			sel := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
			var sb *nodes.SortBy
			if sel.SortClause != nil {
				sb = sel.SortClause.Items[0].(*nodes.SortBy)
			} else {
				fc := sel.TargetList.Items[0].(*nodes.ResTarget).Val.(*nodes.FuncCall)
				sb = fc.AggOrder.Items[0].(*nodes.SortBy)
			}
			if sb.SortbyDir != nodes.SORTBY_USING || sb.UseOp == nil {
				t.Fatalf("SortBy dir %v, UseOp %v", sb.SortbyDir, sb.UseOp)
			}
			var names []string
			for _, item := range sb.UseOp.Items {
				names = append(names, item.(*nodes.String).Str)
			}
			if got := strings.Join(names, "."); got != tc.op {
				t.Errorf("UseOp %q, want %q", got, tc.op)
			}
		})
	}
}

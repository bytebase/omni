package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// TestAlterSemanticViewMaterialization covers ALTER SEMANTIC VIEW's
// materialization actions.
// https://docs.snowflake.com/en/sql-reference/sql/alter-semantic-view
func TestAlterSemanticViewMaterialization(t *testing.T) {
	sql := `ALTER SEMANTIC VIEW revenue_analysis ADD MATERIALIZATION revenue_by_customer
  WAREHOUSE = my_wh
  REFRESH_MODE = INCREMENTAL
  IMMUTABLE WHERE (order_date < '2024-01-01')
  AS
    DIMENSIONS customers.customer_name, region
    METRICS orders.total_revenue
    WHERE (amount > 0)`
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	s := file.Stmts[0].(*ast.AlterSemanticViewStmt)
	if s.Action != ast.AlterSemanticViewAddMaterialization || s.Materialization.Normalize() != "REVENUE_BY_CUSTOMER" ||
		s.Warehouse.Normalize() != "MY_WH" || s.RefreshMode != "INCREMENTAL" || s.ImmutableWhere == nil ||
		len(s.Dimensions) != 2 || s.Dimensions[0].Normalize() != "CUSTOMERS.CUSTOMER_NAME" || len(s.Metrics) != 1 ||
		s.Where == nil || s.Loc.End != len(sql) {
		t.Fatalf("ADD: %+v", s)
	}

	for verb, action := range map[string]ast.AlterSemanticViewAction{
		"DROP":    ast.AlterSemanticViewDropMaterialization,
		"SUSPEND": ast.AlterSemanticViewSuspendMaterialization,
		"RESUME":  ast.AlterSemanticViewResumeMaterialization,
		"REFRESH": ast.AlterSemanticViewRefreshMaterialization,
	} {
		sql := "ALTER SEMANTIC VIEW db.s.sv " + verb + " MATERIALIZATION m1"
		file, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		s := file.Stmts[0].(*ast.AlterSemanticViewStmt)
		if s.Action != action || s.Materialization.Normalize() != "M1" || s.Loc.End != len(sql) {
			t.Errorf("%q: %+v", sql, s)
		}
	}

	for _, sql := range []string{
		// The page writes no IF EXISTS on these forms.
		"ALTER SEMANTIC VIEW IF EXISTS sv DROP MATERIALIZATION m",
		// WAREHOUSE and AS DIMENSIONS ... METRICS are required, in order.
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m AS DIMENSIONS d METRICS x",
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh AS METRICS x",
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh AS DIMENSIONS d",
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh IMMUTABLE WHERE (a) REFRESH_MODE = AUTO AS DIMENSIONS d METRICS x",
		// REFRESH_MODE is AUTO, FULL or INCREMENTAL.
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh REFRESH_MODE = SOMETIMES AS DIMENSIONS d METRICS x",
		// The conditions are parenthesized.
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh AS DIMENSIONS d METRICS x WHERE a > 0",
		// Names: one part for the materialization and warehouse, at most two for members.
		"ALTER SEMANTIC VIEW sv DROP MATERIALIZATION s.m",
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = db.wh AS DIMENSIONS d METRICS x",
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh AS DIMENSIONS a.b.c METRICS x",
		"ALTER SEMANTIC VIEW sv ADD MATERIALIZATION m WAREHOUSE = wh AS DIMENSIONS d, METRICS x",
		// Quoted words are not the keywords.
		`ALTER SEMANTIC VIEW sv "DROP" MATERIALIZATION m`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

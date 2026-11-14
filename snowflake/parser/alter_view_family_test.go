package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func requireStmt[T ast.Node](t *testing.T, sql string) T {
	t.Helper()
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): %d statements", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(T)
	if !ok {
		t.Fatalf("Parse(%q): got %T", sql, file.Stmts[0])
	}
	return stmt
}

func requireParseErrors(t *testing.T, sqls ...string) {
	t.Helper()
	for _, sql := range sqls {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

// ALTER VIEW (https://docs.snowflake.com/en/sql-reference/sql/alter-view).
func TestAlterView_DocumentedForms(t *testing.T) {
	cases := map[string]ast.AlterViewAction{
		"ALTER VIEW v SET CHANGE_TRACKING = TRUE":                                        ast.AlterViewSetChangeTracking,
		"ALTER VIEW v SET CONTACT steward = alice, support = bob":                        ast.AlterViewSetContact,
		"ALTER VIEW v UNSET CONTACT steward":                                             ast.AlterViewUnsetContact,
		"ALTER VIEW v SET TAG t1 = 'x', t2 = 'y'":                                        ast.AlterViewSetTag,
		"ALTER VIEW v UNSET TAG t1, t2":                                                  ast.AlterViewUnsetTag,
		"ALTER VIEW v SET DATA_METRIC_SCHEDULE = 'TRIGGER_ON_CHANGES'":                   ast.AlterViewDataMetric,
		"ALTER VIEW v UNSET DATA_METRIC_SCHEDULE":                                        ast.AlterViewDataMetric,
		"ALTER VIEW v ADD DATA METRIC FUNCTION m ON (a, TABLE r (k)) FILTER (a > 1)":     ast.AlterViewDataMetric,
		"ALTER VIEW v MODIFY DATA METRIC FUNCTION m ON (a) SUSPEND":                      ast.AlterViewDataMetric,
		"ALTER VIEW v DROP DATA METRIC FUNCTION m ON (a)":                                ast.AlterViewDataMetric,
		"ALTER VIEW v SET JOIN POLICY p FORCE":                                           ast.AlterViewSetJoinPolicy,
		"ALTER VIEW v ALTER c SET PROJECTION POLICY p FORCE":                             ast.AlterViewColumnSetProjectionPolicy,
		"ALTER VIEW v MODIFY COLUMN c UNSET PROJECTION POLICY":                           ast.AlterViewColumnUnsetProjectionPolicy,
		"ALTER VIEW v ADD COLUMN IF NOT EXISTS c INT WITH MASKING POLICY m USING (c, d)": ast.AlterViewAddColumn,
		"ALTER VIEW v ADD c VARCHAR PROJECTION POLICY p WITH TAG (t = 'x')":              ast.AlterViewAddColumn,
		"ALTER VIEW v MODIFY data SET MASKING POLICY m":                                  ast.AlterViewColumnSetMaskingPolicy,
		"ALTER VIEW v DROP ROW ACCESS POLICY p, ADD ROW ACCESS POLICY p2 ON (a)":         ast.AlterViewDropRowAccessPolicy,
		"ALTER VIEW v ADD ROW ACCESS POLICY p ON (a, b)":                                 ast.AlterViewAddRowAccessPolicy,
	}
	for sql, want := range cases {
		if got := requireStmt[*ast.AlterViewStmt](t, sql).Action; got != want {
			t.Errorf("%q: action = %v, want %v", sql, got, want)
		}
	}

	stmt := requireStmt[*ast.AlterViewStmt](t, "ALTER VIEW v ALTER a SET TAG t = 's', b UNSET TAG t, COLUMN c SET MASKING POLICY m")
	if len(stmt.ColumnActions) != 3 || stmt.ColumnActions[2].Column.Name != "c" {
		t.Errorf("column actions = %+v", stmt.ColumnActions)
	}
	stmt = requireStmt[*ast.AlterViewStmt](t, "ALTER VIEW v ADD c NUMBER(10, 2) WITH MASKING POLICY m")
	if stmt.Column.Name != "c" || stmt.ColumnType == nil || stmt.MaskingPolicy == nil {
		t.Errorf("add column = %+v", stmt)
	}

	requireParseErrors(t,
		"ALTER VIEW v SET CHANGE_TRACKING = 1",
		"ALTER VIEW v ADD c",                                         // a type is required
		"ALTER VIEW v ADD c INT MASKING POLICY m USING (c)",          // USING names two or more columns
		"ALTER VIEW v ADD c INT WITH TAG (t = 'x') MASKING POLICY m", // documented order
		"ALTER VIEW v ALTER c SET PROJECTION p",
	)
}

// ALTER MATERIALIZED VIEW
// (https://docs.snowflake.com/en/sql-reference/sql/alter-materialized-view).
func TestAlterMaterializedView_DocumentedForms(t *testing.T) {
	stmt := requireStmt[*ast.AlterMaterializedViewStmt](t, "ALTER MATERIALIZED VIEW mv SET SECURE CONTACT steward = alice COMMENT = 'c'")
	if stmt.Action != ast.AlterMVSetSecure || stmt.Contact != "steward = alice" || stmt.Comment == nil {
		t.Errorf("stmt = %+v", stmt)
	}
	cases := map[string]ast.AlterMaterializedViewAction{
		"ALTER MATERIALIZED VIEW mv SET CONTACT a = b, c = d":       ast.AlterMVSetContact,
		"ALTER MATERIALIZED VIEW mv UNSET CONTACT a":                ast.AlterMVUnsetContact,
		"ALTER MATERIALIZED VIEW mv SET DATA_METRIC_SCHEDULE = 's'": ast.AlterMVDataMetric,
		"ALTER MATERIALIZED VIEW mv UNSET DATA_METRIC_SCHEDULE":     ast.AlterMVDataMetric,
		"ALTER MATERIALIZED VIEW mv SET COMMENT = 'c'":              ast.AlterMVSetComment,
	}
	for sql, want := range cases {
		if got := requireStmt[*ast.AlterMaterializedViewStmt](t, sql).Action; got != want {
			t.Errorf("%q: action = %v, want %v", sql, got, want)
		}
	}
	requireParseErrors(t,
		"ALTER MATERIALIZED VIEW mv SET CONTACT",
		"ALTER MATERIALIZED VIEW mv SET COMMENT = 'c' SECURE",
	)
}

// ALTER DYNAMIC TABLE
// (https://docs.snowflake.com/en/sql-reference/sql/alter-dynamic-table).
func TestAlterDynamicTable_DocumentedForms(t *testing.T) {
	cases := map[string]ast.AlterDynamicTableAction{
		"ALTER DYNAMIC TABLE d SUSPEND RECLUSTER":                                        ast.AlterDynamicTableSuspendRecluster,
		"ALTER DYNAMIC TABLE d RESUME RECLUSTER":                                         ast.AlterDynamicTableResumeRecluster,
		"ALTER DYNAMIC TABLE d ADD ROW ACCESS POLICY p ON (a, b)":                        ast.AlterDynamicTableAddRowAccessPolicy,
		"ALTER DYNAMIC TABLE d DROP ROW ACCESS POLICY p, ADD ROW ACCESS POLICY q ON (a)": ast.AlterDynamicTableDropRowAccessPolicy,
		"ALTER DYNAMIC TABLE d DROP ALL ROW ACCESS POLICIES":                             ast.AlterDynamicTableDropAllRowAccessPolicies,
		"ALTER DYNAMIC TABLE d SET AGGREGATION POLICY p ENTITY KEY (a) FORCE":            ast.AlterDynamicTableSetAggregationPolicy,
		"ALTER DYNAMIC TABLE d UNSET AGGREGATION POLICY":                                 ast.AlterDynamicTableUnsetAggregationPolicy,
		"ALTER DYNAMIC TABLE d ALTER a SET TAG t = 's', COLUMN b UNSET MASKING POLICY":   ast.AlterDynamicTableColumn,
		"ALTER DYNAMIC TABLE d MODIFY a SET PROJECTION POLICY p":                         ast.AlterDynamicTableColumn,
		"ALTER DYNAMIC TABLE d ADD SEARCH OPTIMIZATION ON FULL_TEXT(a, ANALYZER => 'x')": ast.AlterDynamicTableAddSearchOpt,
		"ALTER DYNAMIC TABLE d DROP SEARCH OPTIMIZATION ON a, 1":                         ast.AlterDynamicTableDropSearchOpt,
		"ALTER DYNAMIC TABLE d SUSPEND SEARCH OPTIMIZATION ON a":                         ast.AlterDynamicTableSuspendSearchOpt,
		"ALTER DYNAMIC TABLE d RESUME SEARCH OPTIMIZATION":                               ast.AlterDynamicTableResumeSearchOpt,
		"ALTER DYNAMIC TABLE d ADD STORAGE LIFECYCLE POLICY p ON (a)":                    ast.AlterDynamicTableAddStorageLifecyclePolicy,
		"ALTER DYNAMIC TABLE d DROP STORAGE LIFECYCLE POLICY":                            ast.AlterDynamicTableDropStorageLifecyclePolicy,
		"ALTER DYNAMIC TABLE d, e REFRESH":                                               ast.AlterDynamicTableRefresh,
		"ALTER DYNAMIC TABLE d UNSET CONTACT a, FROZEN WHERE, EXECUTE AS USER":           ast.AlterDynamicTableUnset,
		"ALTER DYNAMIC TABLE d SET TAG t = 's'":                                          ast.AlterDynamicTableSetTag,
	}
	for sql, want := range cases {
		if got := requireStmt[*ast.AlterDynamicTableStmt](t, sql).Action; got != want {
			t.Errorf("%q: action = %v, want %v", sql, got, want)
		}
	}

	stmt := requireStmt[*ast.AlterDynamicTableStmt](t, "ALTER DYNAMIC TABLE d SET TARGET_LAG = '1 minute', "+
		"CONTACT (steward = alice), FROZEN WHERE (a > 1), EXECUTE AS USER u USE SECONDARY ROLES r1, r2, WAREHOUSE = wh")
	if stmt.Action != ast.AlterDynamicTableSet || len(stmt.Options) != 2 || stmt.Contact != "(steward = alice)" ||
		stmt.FrozenWhere == nil || stmt.ExecuteAsUser == nil || len(stmt.SecondaryRoles) != 2 {
		t.Errorf("stmt = %+v", stmt)
	}
	if stmt := requireStmt[*ast.AlterDynamicTableStmt](t, "ALTER DYNAMIC TABLE d UNSET CONTACT a, FROZEN WHERE"); len(stmt.UnsetProps) != 2 ||
		stmt.UnsetProps[0] != "CONTACT a" || stmt.UnsetProps[1] != "FROZEN WHERE" {
		t.Errorf("unset = %v", stmt.UnsetProps)
	}

	requireParseErrors(t,
		"ALTER DYNAMIC TABLE d, e SUSPEND", // several names only for REFRESH
		"ALTER DYNAMIC TABLE d SET TARGET_LAG = '1 minute',",
		"ALTER DYNAMIC TABLE d SET WAREHOUSE",
		"ALTER DYNAMIC TABLE d ADD SEARCH OPTIMIZATION ON a", // ADD takes a method
		"ALTER DYNAMIC TABLE d ADD c INT",
	)
}

func TestAlterViewFamilyLookahead(t *testing.T) {
	cases := []struct {
		sql  string
		f    func(p *Parser) bool
		want bool
	}{
		{"ADD DATA METRIC", (*Parser).verbStartsDataMetricFunction, true},
		{"MODIFY data SET", (*Parser).verbStartsDataMetricFunction, false},
		{"CONTACT (a = b)", (*Parser).startsDynamicTableSetItem, true},
		{"FROZEN WHERE", (*Parser).startsDynamicTableSetItem, true},
		{"EXECUTE AS USER", (*Parser).startsDynamicTableSetItem, true},
		{"r2", (*Parser).startsDynamicTableSetItem, false},
		{", WAREHOUSE = w", (*Parser).commaStartsDynamicTableSetItem, true},
		{", r2", (*Parser).commaStartsDynamicTableSetItem, false},
	}
	for _, c := range cases {
		p := &Parser{lexer: NewLexer(c.sql), input: c.sql}
		p.advance()
		before := p.cur
		if got := c.f(p); got != c.want {
			t.Errorf("%q: got %v, want %v", c.sql, got, c.want)
		}
		if p.cur != before {
			t.Errorf("%q: lookahead moved the parser", c.sql)
		}
	}
}

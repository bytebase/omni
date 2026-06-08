package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAlterMaterializedViewOptionsParse(t *testing.T) {
	tests := []struct {
		sql     string
		schema  string
		name    string
		cmdName string
		check   func(t *testing.T, cmd *nodes.AlterTableCmd)
	}{
		{
			sql:     "ALTER MATERIALIZED VIEW tickets_mv AUTO REFRESH YES;",
			name:    "tickets_mv",
			cmdName: "auto_refresh",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionBool(t, cmd, "enabled", true)
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW sales.revenue_mv AUTO REFRESH NO;",
			schema:  "sales",
			name:    "revenue_mv",
			cmdName: "auto_refresh",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionBool(t, cmd, "enabled", false)
			},
		},
		{
			sql:     `ALTER MATERIALIZED VIEW "Customer Orders" ALTER DISTKEY "Customer ID";`,
			name:    "Customer Orders",
			cmdName: "distkey",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionString(t, cmd, "column", "Customer ID")
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW transaction_mv ALTER DISTSTYLE KEY DISTKEY transaction_id;",
			name:    "transaction_mv",
			cmdName: "diststyle",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionString(t, cmd, "style", "key")
				assertAlterTableCmdOptionString(t, cmd, "distkey", "transaction_id")
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW analytics_mv ALTER DISTSTYLE AUTO;",
			name:    "analytics_mv",
			cmdName: "diststyle",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionString(t, cmd, "style", "auto")
			},
		},
		{
			sql:     `ALTER MATERIALIZED VIEW "Monthly-Reports" ALTER SORTKEY ("Report-Date", "Region-Code");`,
			name:    "Monthly-Reports",
			cmdName: "sortkey",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionList(t, cmd, "columns", 2)
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW automated_mv ALTER SORTKEY AUTO;",
			name:    "automated_mv",
			cmdName: "sortkey",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionString(t, cmd, "value", "auto")
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW customer_behavior_mv ALTER COMPOUND SORTKEY (customer_segment, activity_date);",
			name:    "customer_behavior_mv",
			cmdName: "sortkey",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionString(t, cmd, "style", "compound")
				assertAlterTableCmdOptionList(t, cmd, "columns", 2)
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW shared_insights_mv ROW LEVEL SECURITY ON CONJUNCTION TYPE OR FOR DATASHARES;",
			name:    "shared_insights_mv",
			cmdName: "row_level_security",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionBool(t, cmd, "enabled", true)
				assertAlterTableCmdOptionString(t, cmd, "conjunction_type", "or")
				assertAlterTableCmdOptionBool(t, cmd, "for_datashares", true)
			},
		},
		{
			sql:     "ALTER MATERIALIZED VIEW public_data_mv ROW LEVEL SECURITY OFF;",
			name:    "public_data_mv",
			cmdName: "row_level_security",
			check: func(t *testing.T, cmd *nodes.AlterTableCmd) {
				assertAlterTableCmdOptionBool(t, cmd, "enabled", false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstAlterTableStmt(t, tt.sql)
			if stmt.ObjType != int(nodes.OBJECT_MATVIEW) {
				t.Fatalf("expected materialized view object type, got %d", stmt.ObjType)
			}
			if stmt.Relation.Schemaname != tt.schema || stmt.Relation.Relname != tt.name {
				t.Fatalf("expected relation %q.%q, got %#v", tt.schema, tt.name, stmt.Relation)
			}
			if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
				t.Fatalf("expected one command, got %#v", stmt.Cmds)
			}
			cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
			if !ok {
				t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
			}
			if cmd.Subtype != int(nodes.AT_RedshiftAlterOption) {
				t.Fatalf("expected AT_RedshiftAlterOption, got %d", cmd.Subtype)
			}
			if cmd.Name != tt.cmdName {
				t.Fatalf("expected command name %q, got %q", tt.cmdName, cmd.Name)
			}
			tt.check(t, cmd)
		})
	}
}

func alterTableCmdOptions(t *testing.T, cmd *nodes.AlterTableCmd) *nodes.List {
	t.Helper()
	opts, ok := cmd.Def.(*nodes.List)
	if !ok {
		t.Fatalf("expected command options list, got %T", cmd.Def)
	}
	return opts
}

func assertAlterTableCmdOptionString(t *testing.T, cmd *nodes.AlterTableCmd, name string, want string) {
	t.Helper()
	assertDefElemString(t, alterTableCmdOptions(t, cmd), name, want)
}

func assertAlterTableCmdOptionBool(t *testing.T, cmd *nodes.AlterTableCmd, name string, want bool) {
	t.Helper()
	assertDefElemBool(t, alterTableCmdOptions(t, cmd), name, want)
}

func assertAlterTableCmdOptionList(t *testing.T, cmd *nodes.AlterTableCmd, name string, length int) {
	t.Helper()
	elem := findDefElem(alterTableCmdOptions(t, cmd), name)
	if elem == nil {
		t.Fatalf("expected option %q", name)
	}
	list, ok := elem.Arg.(*nodes.List)
	if !ok {
		t.Fatalf("expected option %q list arg, got %T", name, elem.Arg)
	}
	if len(list.Items) != length {
		t.Fatalf("expected option %q list length %d, got %d", name, length, len(list.Items))
	}
}

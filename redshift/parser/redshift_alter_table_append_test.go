package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAlterTableAppendParse(t *testing.T) {
	tests := []struct {
		sql          string
		targetSchema string
		targetName   string
		sourceSchema string
		sourceName   string
		option       string
	}{
		{
			sql:        "ALTER TABLE sales APPEND FROM sales_monthly;",
			targetName: "sales",
			sourceName: "sales_monthly",
		},
		{
			sql:          "ALTER TABLE public.sales APPEND FROM staging.sales_monthly;",
			targetSchema: "public",
			targetName:   "sales",
			sourceSchema: "staging",
			sourceName:   "sales_monthly",
		},
		{
			sql:        "ALTER TABLE sales APPEND FROM sales_listing IGNOREEXTRA;",
			targetName: "sales",
			sourceName: "sales_listing",
			option:     "ignoreextra",
		},
		{
			sql:        "ALTER TABLE sales_report APPEND FROM sales_month FILLTARGET;",
			targetName: "sales_report",
			sourceName: "sales_month",
			option:     "filltarget",
		},
		{
			sql:        `ALTER TABLE "Sales_Table" APPEND FROM "Monthly_Sales";`,
			targetName: "Sales_Table",
			sourceName: "Monthly_Sales",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstAlterTableStmt(t, tt.sql)
			if stmt.Relation.Schemaname != tt.targetSchema || stmt.Relation.Relname != tt.targetName {
				t.Fatalf("expected target %q.%q, got %#v", tt.targetSchema, tt.targetName, stmt.Relation)
			}
			if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
				t.Fatalf("expected one ALTER TABLE command, got %#v", stmt.Cmds)
			}
			cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
			if !ok {
				t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
			}
			if cmd.Subtype != int(nodes.AT_RedshiftAppend) {
				t.Fatalf("expected AT_RedshiftAppend, got %d", cmd.Subtype)
			}
			source, ok := cmd.Def.(*nodes.RangeVar)
			if !ok {
				t.Fatalf("expected source RangeVar, got %T", cmd.Def)
			}
			if source.Schemaname != tt.sourceSchema || source.Relname != tt.sourceName {
				t.Fatalf("expected source %q.%q, got %#v", tt.sourceSchema, tt.sourceName, source)
			}
			if tt.option == "" {
				if cmd.Name != "" {
					t.Fatalf("expected empty option, got %q", cmd.Name)
				}
				return
			}
			if cmd.Name != tt.option {
				t.Fatalf("expected option %q, got %q", tt.option, cmd.Name)
			}
		})
	}
}

func firstAlterTableStmt(t *testing.T, sql string) *nodes.AlterTableStmt {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected AlterTableStmt, got %T", raw.Stmt)
	}
	return stmt
}

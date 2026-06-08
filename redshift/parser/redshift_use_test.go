package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftUseParse(t *testing.T) {
	for _, tc := range []struct {
		sql      string
		database string
	}{
		{"USE mydb", "mydb"},
		{"USE dev_v2", "dev_v2"},
		{"USE database_with_very_long_name_but_still_valid", "database_with_very_long_name_but_still_valid"},
	} {
		stmt := singleStmt(t, tc.sql)
		use, ok := stmt.(*nodes.VariableSetStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected VariableSetStmt, got %T", tc.sql, stmt)
		}
		if use.Kind != nodes.VAR_SET_VALUE {
			t.Fatalf("Parse(%q): kind = %v, want VAR_SET_VALUE", tc.sql, use.Kind)
		}
		if use.Name != "database" {
			t.Fatalf("Parse(%q): name = %q, want database", tc.sql, use.Name)
		}
		assertSingleStringArg(t, use.Args, tc.database)
	}
}

func assertSingleStringArg(t *testing.T, args *nodes.List, want string) {
	t.Helper()
	if args == nil || len(args.Items) != 1 {
		t.Fatalf("args = %#v, want one string arg", args)
	}
	arg, ok := args.Items[0].(*nodes.A_Const)
	if !ok {
		t.Fatalf("arg = %T, want A_Const", args.Items[0])
	}
	str, ok := arg.Val.(*nodes.String)
	if !ok {
		t.Fatalf("arg value = %T, want String", arg.Val)
	}
	if str.Str != want {
		t.Fatalf("arg = %q, want %q", str.Str, want)
	}
}

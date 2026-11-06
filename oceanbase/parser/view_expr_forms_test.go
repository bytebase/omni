package parser

import (
	"testing"

	ast "github.com/bytebase/omni/oceanbase/ast"
)

// Expression forms tidb/parser was found to misread while writing TiDB's view
// restore, held here because this engine shares them: N'...', TRIM with a
// direction and no remstr, and CHAR ... USING.
//
// The verdicts were measured against oceanbase/oceanbase-ce:4.3.5-lts (4.3.5.6)
// through PREPARE on 2026-11-06. TestViewExprFormsOceanBaseOracle holds them to the server, and TestViewExprFormTrees
// holds the tree each accepted form builds.

type viewExprFormCase struct {
	sql    string
	accept bool
}

var viewExprFormCases = []viewExprFormCase{
	// N or n directly before a single quote is NATIONAL_LITERAL, and the lexer
	// folds adjacent quoted strings into it (sql_parser_mysql_mode.l:429).
	{`select N'nat'`, true},
	{`select n'nat'`, true},
	{`select N''`, true},
	{`select N'x' 'y'`, true},
	{`select N'x' 'y' 'z'`, true},
	{`select N'x' as y`, true},
	{`select N 'x' from t`, true},
	{`select N 'x' as y from t`, false},
	{"select `N`'x' from t", true},
	{`select N"x" from t`, true},
	{`select N"x" as y from t`, false},

	// TRIM takes a direction with or without remstr, and FROM after either.
	{`select trim(both from c) from t`, true},
	{`select trim(leading from c) from t`, true},
	{`select trim(trailing from c) from t`, true},
	{`select trim(both 'x' from c) from t`, true},
	{`select trim('x' from c) from t`, true},
	{`select trim(c) from t`, true},
	{`select trim(leading 'x')`, false},
	{`select trim(both)`, false},
	{`select trim(both from)`, false},
	{`select trim(both from 'x' from 'y')`, false},

	// CHAR takes an expression list and an optional USING charset_name.
	{`select char(65 using utf8mb4)`, true},
	{`select char(65, 66 using 'latin1')`, true},
	{`select char(65 using binary)`, true},
	{"select char(65 using `utf8mb4`)", true},
	{`select char(65)`, true},
	{`select char()`, false},
	{`select char(using utf8mb4)`, false},
	{`select char(65 using)`, false},
	{`select char(65 using utf8mb4 collate utf8mb4_bin)`, false},
	{`select char(65 using utf8mb4, 66)`, false},
	{`select char(65 using utf8mb4) over ()`, false},
	{"select `char`(65 using utf8mb4)", false},
}

// TestViewExprForms holds the parser to the verdicts above.
func TestViewExprForms(t *testing.T) {
	for _, c := range viewExprFormCases {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if got := err == nil; got != c.accept {
				t.Fatalf("Parse(%q) accepts = %v, the engine answered %v (err: %v)", c.sql, got, c.accept, err)
			}
			if err == nil && len(list.Items) != 1 {
				t.Errorf("Parse(%q) read %d statements, want 1", c.sql, len(list.Items))
			}
		})
	}
}

// TestViewExprFormTrees pins what each form carries.
func TestViewExprFormTrees(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{`select N'nat'`, `{SELECT :loc 0 :targets ({STRING_LIT :val "nat" :loc 7 :national true})}`},
		{`select N'x' 'y'`, `{SELECT :loc 0 :targets ({STRING_LIT :val "xy" :loc 7 :national true :first-segment "x"})}`},
		{`select N 'x'`, `{SELECT :loc 0 :targets ({RES_TARGET :loc 7 :name x :val {COLREF :loc 7 :col N}})}`},
		{`select trim(both from c)`, `{SELECT :loc 0 :targets ({FUNCCALL :name TRIM_BOTH :loc 7 :args ({COLREF :loc 22 :col c})})}`},
		{`select char(65, 66 using 'latin1')`, `{SELECT :loc 0 :targets ({FUNCCALL :name CHAR :loc 7 :using latin1 :args ({INT_LIT :val 65 :loc 12} {INT_LIT :val 66 :loc 16})})}`},
		{`select char(65)`, `{SELECT :loc 0 :targets ({FUNCCALL :name CHAR :loc 7 :args ({INT_LIT :val 65 :loc 12})})}`},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			if got := ast.NodeToString(list.Items[0]); got != tc.want {
				t.Errorf("tree\n got %s\nwant %s", got, tc.want)
			}
		})
	}
}

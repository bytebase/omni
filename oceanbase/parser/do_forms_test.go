package parser

import "testing"

// DO expression forms verified via PREPARE on pinned OceanBase CE
// 4.3.5.6-106000012026040916 (5.7.25-OceanBase_CE-v4.3.5.6).
//
// Ref: src/sql/parser/ob_parser.cpp (S_DO) and
// src/pl/parser/pl_parser_mysql_mode.y (do_sp_stmt, lines 554-606): OceanBase
// rejects any top-level T_COLUMN_REF item (including parenthesized column
// references and table.*) with 1064, while accepting non-column expressions
// (literals, operators, function calls, user/system variables, subqueries, and
// bare `*`, which fails after parsing with 1149).
var doFormCases = []parserFormCase{
	// Rejected by do_sp_stmt (1064): top-level column references (bare or parenthesized).
	{"DO a", false},
	{"DO `a`", false},
	{"DO (a)", false},
	{"DO ((a))", false},
	{"DO (/* + */ a)", false},
	{"DO 1, a", false},
	{"DO a, 1", false},
	{"DO t.a", false},
	{"DO (t.a)", false},
	{"DO a.b.c", false},
	{"DO a.b.c.d", false},
	{"DO t.*", false},
	{"DO a zzz", false},
	{"DO 1 zzz", false},
	{"DO", false},

	// Accepted by OceanBase's parser (either OK, 1327 undeclared variable, or 1149).
	{"DO 1", true},
	{"DO (1)", true},
	{"DO 1, 2", true},
	{"DO 's'", true},
	{"DO @v", true},
	{"DO @@version", true},
	{"DO NULL", true},
	{"DO TRUE", true},
	{"DO *", true},
	{"DO 1, *", true},
	{"DO a + 1", true},
	{"DO 1 + a", true},
	{"DO (a + 1)", true},
	{"DO -a", true},
	{"DO +a", true},
	{"DO (+a)", true},
	{"DO +((a))", true},
	{"DO NOT a", true},
	{"DO ~a", true},
	{"DO a IS NULL", true},
	{"DO a = 1", true},
	{"DO @v := a", true},
	{"DO COALESCE(a, 1)", true},
	{"DO (SELECT a)", true},
}

func TestDoForms(t *testing.T) {
	for _, c := range doFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

func TestDoFormsMultiStatement(t *testing.T) {
	cases := []parserFormCase{
		{"SELECT 1; DO +a", true},
		{"SELECT 1; DO (+a)", true},
		{"SELECT 1; DO +((a))", true},
		{"SELECT 1; DO a", false},
		{"SELECT 1; DO (a)", false},
		{"SELECT 1; DO (/* + */ a)", false},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			got := err == nil
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v: %v", c.sql, got, c.accept, err)
			}
		})
	}
}

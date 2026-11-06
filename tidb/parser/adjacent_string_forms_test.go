package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Adjacent string literals as TiDB v8.5.5 reads them, every verdict measured with
// PREPARE on pingcap/tidb:v8.5.5 on 2026-11-06.
//
// StringLiteral (parser.y:7660) is a run of bare strings, folded into one: 'a' 'b'
// is the string ab, wherever an expression stands. Its column is named after the
// first string (SetProjectionOffset). An introduced string takes one string
// ("UNDERSCORE_CHARSET" stringLit, parser.y:7599), so _utf8mb4'a' 'b' and N'a' 'b'
// are the string a aliased b, and a run cannot continue with an introducer or a hex
// literal. COMMENT and ENUM read a single stringLit, so they refuse a run.

type adjacentStringCase struct {
	sql    string
	accept bool
}

var adjacentStringCases = []adjacentStringCase{
	{"select 'a' 'b'", true},
	{"select 1 from t where a = 'a' 'b'", true},
	{"select ('a' 'b')", true},
	{"select 'a' 'b' 'c'", true},
	{"select 'a'\n'b'", true},
	{"select 'a' /* c */ 'b'", true},
	{"select 'a' \"b\"", true},
	{"select '' 'b'", true},
	{"select 'a' 'b' as c", true},
	{"select 'a' 'b' c", true},
	{"select concat('a' 'b', 'c')", true},
	{"insert into t values ('a' 'b')", true},
	{"select _utf8mb4'a' 'b'", true},
	{"select N'a' 'b'", true},
	{"select 'a' _utf8mb4'b'", false},
	{"select 'a' x'62'", false},
	{"select * from t limit 1 offset 'a' 'b'", false},
	{"create table t (a varchar(10) default 'a' 'b')", true},
	{"create table t (a int comment 'x' 'y')", false},
	{"create table t (a enum('x' 'y'))", false},
	{"create table t (a int) comment 'x' 'y'", false}, {"select 1 from t where 'a' like 'b' escape 'x' 'y'", false},
	{"select 1 from t where 'a' ilike 'b' escape 'x' 'y'", false},
	{"select 'a' like 'b' escape 'x' 'y'", true},
}

// TestAdjacentStringForms holds omni to the recorded verdicts without a server.
func TestAdjacentStringForms(t *testing.T) {
	for _, c := range adjacentStringCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbCorpusVersion, err)
			}
		})
	}
}

// TestAdjacentStringTrees checks a run folds into one literal holding the joined
// value and its first string, and that an introduced string does not start one.
func TestAdjacentStringTrees(t *testing.T) {
	item := func(t *testing.T, sql string) nodes.ExprNode {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0].(*nodes.SelectStmt).TargetList[0]
	}
	for _, c := range []struct {
		sql, value, first, text string
	}{
		{"select 'a' 'b'", "ab", "a", "'a' 'b'"},
		{"select 'a' 'b' 'c'", "abc", "a", "'a' 'b' 'c'"},
		{"select '' 'b'", "b", "", "'' 'b'"},
		{"select 'a'\n'b'", "ab", "a", "'a'\n'b'"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			lit, ok := item(t, c.sql).(*nodes.StringLit)
			if !ok {
				t.Fatalf("Parse(%q): item is not a StringLit", c.sql)
			}
			if lit.Value != c.value || !lit.Concatenated || lit.FirstSegment != c.first {
				t.Errorf("Parse(%q) = %q concatenated=%v first=%q, want %q first %q", c.sql, lit.Value, lit.Concatenated, lit.FirstSegment, c.value, c.first)
			}
			if got := c.sql[lit.Loc.Start:lit.Loc.End]; got != c.text {
				t.Errorf("Loc covers %q, want %q", got, c.text)
			}
		})
	}
	t.Run("one string", func(t *testing.T) {
		lit, ok := item(t, "select 'a'").(*nodes.StringLit)
		if !ok || lit.Concatenated || lit.FirstSegment != "" {
			t.Errorf("select 'a' = %#v, want an unfolded literal", lit)
		}
	})
	for _, sql := range []string{"select _utf8mb4'a' 'b'", "select N'a' 'b'"} {
		t.Run(sql, func(t *testing.T) {
			rt, ok := item(t, sql).(*nodes.ResTarget)
			if !ok || rt.Name != "b" {
				t.Fatalf("Parse(%q) = %#v, want the string a aliased b", sql, item(t, sql))
			}
			if lit, ok := rt.Val.(*nodes.StringLit); !ok || lit.Value != "a" || lit.Concatenated {
				t.Errorf("Parse(%q) value = %#v, want the string a", sql, rt.Val)
			}
		})
	}
}

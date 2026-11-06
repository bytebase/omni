package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// Adjacent string literals as MariaDB 11.8.8 reads them, every verdict measured
// with PREPARE on mariadb:11.8.8 on 2026-11-06.
//
// text_literal (sql/sql_yacc.yy) folds a run of quoted strings into one: 'a' 'b' is
// the string ab, wherever an expression stands. Unlike TiDB, a run may start with an
// introduced or national string (_utf8mb4'a' 'b' and N'a' 'b' are ab), and unlike
// MySQL, the column is named after the whole value (SELECT 'a' 'b' is column ab). A
// run continues only with bare strings, and COMMENT and ENUM take one string.

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
	{"create table t (a int) comment 'x' 'y'", false}, {"select 1 from t where 'a' like 'b' escape 'x' 'y'", true},
	{"select 'a' like 'b' escape 'x' 'y'", true},
}

// TestAdjacentStringForms holds omni to the recorded verdicts without a server.
func TestAdjacentStringForms(t *testing.T) {
	for _, c := range adjacentStringCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (MariaDB 11.8.8): %v", c.sql, got, c.accept, err)
			}
		})
	}
}

// TestAdjacentStringTrees checks a run folds into one literal, keeping the
// introducer or national marker it starts with.
func TestAdjacentStringTrees(t *testing.T) {
	for _, c := range []struct {
		sql, value, charset, text string
		national                  bool
	}{
		{"select 'a' 'b'", "ab", "", "'a' 'b'", false},
		{"select 'a' 'b' 'c'", "abc", "", "'a' 'b' 'c'", false},
		{"select _utf8mb4'a' 'b'", "ab", "_utf8mb4", "_utf8mb4'a' 'b'", false},
		{"select N'a' 'b'", "ab", "", "N'a' 'b'", true},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			lit, ok := list.Items[0].(*nodes.SelectStmt).TargetList[0].(*nodes.StringLit)
			if !ok {
				t.Fatalf("Parse(%q): item is not a StringLit", c.sql)
			}
			if lit.Value != c.value || lit.Charset != c.charset || lit.National != c.national {
				t.Errorf("Parse(%q) = %q charset %q national %v, want %q %q %v", c.sql, lit.Value, lit.Charset, lit.National, c.value, c.charset, c.national)
			}
			if got := c.sql[lit.Loc.Start:lit.Loc.End]; got != c.text {
				t.Errorf("Loc covers %q, want %q", got, c.text)
			}
		})
	}
}

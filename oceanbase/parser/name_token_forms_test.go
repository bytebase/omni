package parser

import (
	"fmt"
	"testing"
)

// Three ways a name position read more than a name, or less.
//
// parseIdent took every identifier token, and the lexer hands it two that are not names:
// a user variable, which it scans whole as @a for the expression reader, and a character
// set introducer such as _utf8, which is one word to it. OceanBase's lexer has a token for
// each, so `CREATE TABLE @a (b int)` and `CREATE TABLE _utf8 (b int)` are 1064s, while the
// backquoted spellings are names. Which words introduce is OceanBase's own list, not
// MySQL's: `_hkscs` introduces, and `_ucs2` and `_cp1250` are names.
//
// OceanBase's flex lexer takes the longest match, and reads a word, a dot and the
// identifier characters after it as one qualified name, so the digits after the dot are a
// name: `SELECT t.1 FROM t` and `CREATE TABLE d.1b (c int)` parse, and omni read `.1` as a
// number. The same rule makes an introducer that runs into `.t` a name, and a number spelt
// with a letter, such as 1e5 or 0x10, the first part of a qualified name.
//
// Every case was measured against OceanBase CE 4.3.5.6 (oceanbase/oceanbase-ce:4.3.5-lts)
// through PREPARE on 2026-11-05, with a database selected.
// TestNameTokenFormsOceanBaseOracle keeps it a recording. `CREATE TABLE a.1.5 (c int)`
// is left out: OceanBase reads a three-part table name as catalog.database.table, which
// this parser does not read with any spelling. So is `SELECT @@session.autocommit.1`: a
// name after @@scope.name is a 1064 on OceanBase, which this parser takes with any
// spelling.
var nameTokenFormCases = []parserFormCase{
	// User variables are not names.
	{"CREATE TABLE @a (b int)", false},
	{"SELECT * FROM @a", false},
	{"SELECT * FROM t.@a", false},
	{"SELECT * FROM @a.t", false},
	{"SELECT 1 AS @a", false},
	{"SELECT * FROM `@a`", true},
	{"CREATE TABLE `@a` (b int)", true},
	{"SELECT @a", true},
	{"SELECT @a.1", true},
	{"SET @a = 1", true},

	// Character set introducers are not names, unless a dot joins them to one.
	{"CREATE TABLE _utf8 (b int)", false},
	{"SELECT * FROM _utf8", false},
	{"SELECT * FROM _UTF8MB4", false},
	{"SELECT * FROM _binary", false},
	{"SELECT * FROM _gb18030", false},
	{"SELECT * FROM _hkscs", false},
	{"SELECT _utf8 FROM t", false},
	{"SELECT 1 AS _utf8", false},
	{"CREATE TABLE t (_latin1 int)", false},
	{"SELECT * FROM _utf8. t", false},
	{"SELECT * FROM _utf8 .t", false},
	{"SELECT * FROM d. _utf8", false},
	{"SELECT d. _utf8 FROM t", false},
	{"SELECT * FROM `_utf8`", true},
	{"CREATE TABLE `_utf8` (b int)", true},
	{"SELECT * FROM d._utf8", true},
	{"SELECT * FROM _utf8.t", true},
	{"SELECT * FROM _utf8.1", true},
	{"SELECT _utf8.t.c FROM t", true},
	{"SELECT * FROM t._utf8 x", true},
	{"SELECT a._utf8 FROM t", true},
	{"SELECT * FROM _ucs2", true},
	{"SELECT * FROM _cp1250", true},
	{"SELECT * FROM _filename", true},
	{"SELECT * FROM _utf8mb4_bin", true},
	{"SELECT * FROM _utf8_general_ci", true},
	{"SELECT * FROM _foo", true},
	{"SELECT _utf8mb4'x'", true},

	// A digit run after a separating dot is a name.
	{"SELECT a.1 FROM t", true},
	{"SELECT t.1 FROM t", true},
	{"SELECT t.1b FROM t", true},
	{"SELECT t.1e5 FROM t", true},
	{"SELECT t.1e+5 FROM t", true},
	{"SELECT t.1e FROM t", true},
	{"SELECT t.0x10 FROM t", true},
	{"SELECT t.0b1 FROM t", true},
	{"SELECT t.123456789012345678901234567890 FROM t", true},
	{"SELECT t.1.5 FROM t", true},
	{"SELECT t.1.e1 FROM t", true},
	{"SELECT d.t.1 FROM t", true},
	{"SELECT t.`1` FROM t", true},
	{"SELECT select.1 FROM t", true},
	{"SELECT 1a.1 FROM t", true},
	{"SELECT 12e.1 FROM t", true},
	{"SELECT a.$1 FROM t", true},
	{"SELECT a.1$ FROM t", true},
	{"SELECT t.1-1 FROM t", true},
	{"SELECT t.*, t.1 FROM t", true},
	{"SELECT t.1 'x' FROM t", true},
	{"SELECT t.1 AS x FROM t", true},
	{"SELECT COUNT(t.1) FROM t", true},
	{"SELECT 1 FROM t WHERE t.1 = 1", true},
	{"SELECT 1 FROM t ORDER BY t.1", true},
	{"SELECT t.1() FROM t", true},
	{"SELECT d.1(1)", true},
	{"UPDATE t SET t.1 = 1", true},
	{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.1 = 1", true},
	{"CREATE TABLE a.1b (c int)", true},
	{"CREATE TABLE a.1 (c int)", true},
	{"CREATE TABLE a.1e5 (c int)", true},
	{"CREATE TABLE a. 1b (c int)", true},
	{"SELECT * FROM 1db.1table", true},
	{"SELECT * FROM d.1", true},
	{"SELECT * FROM d.1 x", true},
	{"INSERT INTO d.1 VALUES (1)", true},
	{"SELECT 1.5", true},
	{"SELECT 1.a FROM t", true},

	// A number spelt with a letter is a word to flex, which runs into the dot.
	{"SELECT 1e5.1 FROM t", true},
	{"SELECT 1e5.a FROM t", true},
	{"SELECT 1e5.1.2 FROM t", true},
	{"SELECT 0x10.1 FROM t", true},
	{"SELECT 0b1.a FROM t", true},
	{"SELECT 0x.1 FROM t", true},
	{"SELECT 1e.1 FROM t", true},
	{"SELECT 1 FROM t WHERE a=1e5.1", true},
	{"SELECT 1 FROM t WHERE a=0x1F.b", true},
	{"SELECT 1e5 .a FROM t", false},
	{"SELECT 1e5.`a` FROM t", false},
	{"SELECT 1e+5.1 FROM t", false},
	{"SELECT 12.5.a FROM t", false},

	// A quoted name, a space, a number, a fourth part or @@ does not.
	{"SELECT `t`.1 FROM t", false},
	{"SELECT `t`.1b FROM t", false},
	{"SELECT t .1 FROM t", false},
	{"SELECT t. 1 FROM t", false},
	{"SELECT t . 1 FROM t", false},
	{"SELECT t.1. FROM t", false},
	{"SELECT t.1.2.3 FROM t", false},
	{"SELECT 1 .5", false},
	{"SELECT 'a'.1", false},
	{"SELECT x'0a'.1", false},
	{"SELECT (t).1 FROM t", false},
	{"SELECT @`a`.1", false},
	{"SELECT @@global.1", false},
	{"SELECT @@autocommit.1", false},
	{"SELECT @@global.autocommit", true},
	{"CREATE TABLE a .1b (c int)", false},
	{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW .1 = 1", false},
	{"CREATE TABLE `a`.1b (c int)", false},
}

// parserFormCase is one statement text and whether the engine's parser accepts it.
type parserFormCase struct {
	sql    string
	accept bool
}

// TestNameTokenForms holds omni to nameTokenFormCases.
func TestNameTokenForms(t *testing.T) {
	for _, tc := range nameTokenFormCases {
		_, err := Parse(tc.sql)
		if got := err == nil; got != tc.accept {
			t.Errorf("Parse(%q) accepts=%v, OceanBase 4.3.5.6 accepts=%v: %v", tc.sql, got, tc.accept, err)
		}
	}
}

// TestDotSeparatorTokens pins the tokens a dot after a word leaves, which a verdict does
// not always show: `1e5.1` is refused whichever way it is cut. A dot separates only after a
// word, a keyword among them, and only then is what follows it a name.
func TestDotSeparatorTokens(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []int
	}{
		{"t.1", []int{tokIDENT, '.', tokIDENT}},
		{"t.1.5", []int{tokIDENT, '.', tokIDENT, '.', tokIDENT}},
		{"t.1e5", []int{tokIDENT, '.', tokIDENT}},
		{"NEW.1", []int{kwNEW, '.', tokIDENT}},
		{"1a.1", []int{tokIDENT, '.', tokIDENT}},
		{"`t`.1", []int{tokIDENT, tokFCONST}},
		{"t .1", []int{tokIDENT, tokFCONST}},
		{"t. 1", []int{tokIDENT, '.', tokICONST}},
		{"'a'.1", []int{tokSCONST, tokFCONST}},
		{"1e5.1", []int{tokIDENT, '.', tokIDENT}},
		{"0x10.1", []int{tokIDENT, '.', tokIDENT}},
		{"1e+5.1", []int{tokFCONST, tokFCONST}},
		{"12.5.a", []int{tokFCONST, '.', tokIDENT}},
		{"@@global.1", []int{tokAt2, kwGLOBAL, tokFCONST}},
		{"@@autocommit.1", []int{tokAt2, tokIDENT, tokFCONST}},
	} {
		l := NewLexer(c.sql)
		var got []int
		for tok := l.NextToken(); tok.Type != tokEOF; tok = l.NextToken() {
			got = append(got, tok.Type)
		}
		if fmt.Sprint(got) != fmt.Sprint(c.want) {
			t.Errorf("tokens of %q = %v, want %v", c.sql, got, c.want)
		}
	}
}

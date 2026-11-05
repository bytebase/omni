package parser

import (
	"fmt"
	"testing"
)

// Three ways a name position read more than a name, or less.
//
// parseIdent took every identifier token, and the lexer hands it two that are not names:
// a user variable, which it scans whole as @a for the expression reader, and a character
// set introducer such as _utf8, which is one word to it. The server's lexer gives neither
// to its ident rule, so `CREATE TABLE @a (b int)` and `CREATE TABLE _utf8 (b int)` are
// 1064s, while the backquoted spellings are names. Which words introduce is the server's
// list of character sets: `_filename`, collation names and `_utf8mb4_bin` are names.
//
// A word that runs straight into a dot and an identifier character makes the dot a
// separator, and the server reads the identifier characters after it as a name even when
// they start with a digit. omni read `.1` as a number, so `SELECT t.1 FROM t` and `CREATE
// TABLE d.1b (c int)` failed.
//
// Every case was measured against MySQL 8.0 (mysql:8.0, 8.0.46) through PREPARE on
// 2026-11-05, with a database selected. TestNameTokenFormsOracle keeps it a recording.
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

	// Character set introducers are not names, before a dot either.
	{"CREATE TABLE _utf8 (b int)", false},
	{"SELECT * FROM _utf8", false},
	{"SELECT * FROM _UTF8MB4", false},
	{"SELECT * FROM _binary", false},
	{"SELECT * FROM _utf8.t", false},
	{"SELECT _utf8 FROM t", false},
	{"SELECT 1 AS _utf8", false},
	{"CREATE TABLE t (_latin1 int)", false},
	{"SELECT * FROM `_utf8`", true},
	{"CREATE TABLE `_utf8` (b int)", true},
	{"SELECT * FROM d._utf8", true},
	{"SELECT a._utf8 FROM t", true},
	{"SELECT * FROM d. _utf8", false},
	{"SELECT d. _utf8 FROM t", false},
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
	{"SELECT @@global.1", true},             // 1193, a variable looked up by name
	{"SELECT @@autocommit.1", true},         // 1193
	{"SELECT @@session.autocommit.1", true}, // 1193
	{"CREATE TABLE a.1b (c int)", true},
	{"CREATE TABLE a.1 (c int)", true},
	{"CREATE TABLE a.1e5 (c int)", true},
	{"CREATE TABLE a. 1b (c int)", true},
	{"SELECT * FROM 1db.1table", true},
	{"SELECT * FROM d.1", true},
	{"SELECT * FROM d.1 x", true},
	{"INSERT INTO d.1 VALUES (1)", true},
	{"SELECT 1.5", true},

	// A quoted name, a space, a number or a fourth part does not.
	{"SELECT `t`.1 FROM t", false},
	{"SELECT `t`.1b FROM t", false},
	{"SELECT t .1 FROM t", false},
	{"SELECT t. 1 FROM t", false},
	{"SELECT t . 1 FROM t", false},
	{"SELECT t.1. FROM t", false},
	{"SELECT t.1.2.3 FROM t", false},
	{"SELECT 1e5.1 FROM t", false},
	{"SELECT 1 .5", false},
	{"SELECT 'a'.1", false},
	{"SELECT x'0a'.1", false},
	{"SELECT (t).1 FROM t", false},
	{"SELECT @`a`.1", false},
	{"CREATE TABLE a.1.5 (c int)", false},
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
			t.Errorf("Parse(%q) accepts=%v, MySQL 8.0 accepts=%v: %v", tc.sql, got, tc.accept, err)
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
		{"1e5.1", []int{tokFCONST, tokFCONST}},
		{"0x10.1", []int{tokXCONST, tokFCONST}},
		{"@@global.1", []int{tokAt2, kwGLOBAL, '.', tokIDENT}},
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

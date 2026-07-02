package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/ast"
)

// MySQL concatenates adjacent quoted string literals into a single literal:
// 'a' 'b' → 'ab' (manual 9.1.1: "Quoted strings placed next to each other are
// concatenated to a single string"). The rule lives on text_literal only —
// hex/bit/temporal literals never join a run, a charset introducer is valid on
// the FIRST segment only, and the implicit output column name derives from the
// first segment ('SELECT 'a' 'b'' → column "a", value "ab").
//
// Oracle evidence: MySQL 8.0.32 and 5.7.25 agree on every case below
// (probed live; see PR body). The real-world motivation is the stock MySQL 8.0
// sys schema: ps_trace_thread's body contains an adjacent-literal run, so a
// canonical dump of sys was unparseable before this support.

// firstSelectItem parses sql and returns the single select item, which may be a
// bare expression or a *ast.ResTarget wrapper when an alias was attached.
func firstSelectItem(t *testing.T, sql string) ast.ExprNode {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): expected 1 statement, got %d", sql, len(list.Items))
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): expected *ast.SelectStmt, got %T", sql, list.Items[0])
	}
	if len(sel.TargetList) != 1 {
		t.Fatalf("Parse(%q): expected 1 select item, got %d", sql, len(sel.TargetList))
	}
	return sel.TargetList[0]
}

func TestAdjacentStringLiteralConcat(t *testing.T) {
	cases := []struct {
		name         string
		sql          string
		wantValue    string
		wantCharset  string
		wantFirstSeg string // FirstSegment when the run concatenated; "" for aliasless single
	}{
		{name: "two literals", sql: `SELECT 'a' 'b'`, wantValue: "ab", wantFirstSeg: "a"},
		{name: "three literals", sql: `SELECT 'a' 'b' 'c'`, wantValue: "abc", wantFirstSeg: "a"},
		{name: "newline between", sql: "SELECT 'a'\n'b'", wantValue: "ab", wantFirstSeg: "a"},
		{name: "block comment between", sql: `SELECT 'a' /* c */ 'b'`, wantValue: "ab", wantFirstSeg: "a"},
		{name: "line comment between", sql: "SELECT 'a' -- c\n'b'", wantValue: "ab", wantFirstSeg: "a"},
		{name: "double-quote mix", sql: `SELECT 'a' "b"`, wantValue: "ab", wantFirstSeg: "a"},
		{name: "double-quote first", sql: `SELECT "a" 'b'`, wantValue: "ab", wantFirstSeg: "a"},
		{name: "empty middle segment", sql: `SELECT 'a' '' 'b'`, wantValue: "ab", wantFirstSeg: "a"},
		{name: "empty first segment", sql: `SELECT '' 'b'`, wantValue: "b", wantFirstSeg: ""},
		{name: "escaped quote in segment", sql: `SELECT 'a''x' 'b'`, wantValue: "a'xb", wantFirstSeg: "a'x"},
		{name: "introducer on first", sql: `SELECT _utf8mb4'a' 'b'`, wantValue: "ab", wantCharset: "_utf8mb4", wantFirstSeg: "a"},
		{name: "single literal unchanged", sql: `SELECT 'a'`, wantValue: "a"},
		{name: "single introducer unchanged", sql: `SELECT _utf8mb4'a'`, wantValue: "a", wantCharset: "_utf8mb4"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := firstSelectItem(t, tc.sql)
			lit, ok := item.(*ast.StringLit)
			if !ok {
				t.Fatalf("expected *ast.StringLit select item, got %T", item)
			}
			if lit.Value != tc.wantValue {
				t.Errorf("Value = %q, want %q", lit.Value, tc.wantValue)
			}
			if lit.Charset != tc.wantCharset {
				t.Errorf("Charset = %q, want %q", lit.Charset, tc.wantCharset)
			}
			wantConcat := strings.Count(tc.sql, "'")+strings.Count(tc.sql, `"`) > 2
			if lit.Concatenated != wantConcat {
				t.Errorf("Concatenated = %v, want %v", lit.Concatenated, wantConcat)
			}
			if lit.FirstSegment != tc.wantFirstSeg {
				t.Errorf("FirstSegment = %q, want %q", lit.FirstSegment, tc.wantFirstSeg)
			}
		})
	}
}

// TestAdjacentStringLiteralExprContexts proves the run folds inside every
// expression position the oracle accepts it in: function args, comparisons,
// IN lists, ORDER BY, and INSERT values.
func TestAdjacentStringLiteralExprContexts(t *testing.T) {
	cases := []string{
		`SELECT CONCAT('a' 'b', 'c')`,
		`SELECT 'a' 'b' = 'ab'`,
		`SELECT * FROM t WHERE c = 'a' 'b'`,
		`SELECT * FROM t WHERE c IN ('a' 'b', 'c')`,
		`SELECT c FROM t ORDER BY 'a' 'b'`,
		`INSERT INTO t VALUES ('a' 'b')`,
		`SET @v = 'a' 'b'`,
		`CREATE TABLE t (a varchar(10) DEFAULT 'x' 'y')`,
		`CREATE TABLE t (a varchar(10) DEFAULT 'x' 'y' 'z' NOT NULL)`,
		`CREATE TABLE t (a varchar(10) DEFAULT _utf8mb4'x' 'y')`,
		`CREATE TABLE t (a varchar(20) DEFAULT ('x' 'y'))`,
		`CREATE TABLE t (a int, g varchar(10) GENERATED ALWAYS AS ('x' 'y') STORED)`,
		"CREATE TABLE tp (c varchar(10) NOT NULL) PARTITION BY LIST COLUMNS(c) (PARTITION p0 VALUES IN ('a' 'b'), PARTITION p1 VALUES IN ('z'))",
		`CREATE VIEW v AS SELECT 'a' 'b' AS c`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

// TestAdjacentStringLiteralDefaultValue pins the folded value on the column
// default: MySQL stores DEFAULT 'x' 'y' as DEFAULT 'xy' (oracle: SHOW CREATE
// TABLE on 8.0.32 and 5.7.25).
func TestAdjacentStringLiteralDefaultValue(t *testing.T) {
	list, err := Parse(`CREATE TABLE t (a varchar(10) DEFAULT 'x' 'y' 'z')`)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	ct, ok := list.Items[0].(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", list.Items[0])
	}
	lit, ok := ct.Columns[0].DefaultValue.(*ast.StringLit)
	if !ok {
		t.Fatalf("expected *ast.StringLit default, got %T", ct.Columns[0].DefaultValue)
	}
	if lit.Value != "xyz" {
		t.Errorf("default Value = %q, want %q", lit.Value, "xyz")
	}
}

// TestAdjacentStringLiteralAliasInterplay pins the boundary between adjacency
// and MySQL's implicit string alias: a quoted string after a TEXT literal joins
// the run (greedy, matching MySQL), while a quoted string after any other
// literal stays an alias.
func TestAdjacentStringLiteralAliasInterplay(t *testing.T) {
	t.Run("run swallows would-be alias", func(t *testing.T) {
		item := firstSelectItem(t, `SELECT 'a' 'b'`)
		if _, ok := item.(*ast.ResTarget); ok {
			t.Fatalf("SELECT 'a' 'b' parsed as aliased item; MySQL concatenates instead")
		}
	})
	t.Run("explicit AS breaks the run", func(t *testing.T) {
		item := firstSelectItem(t, `SELECT 'a' AS 'b'`)
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			t.Fatalf("expected *ast.ResTarget, got %T", item)
		}
		if rt.Name != "b" {
			t.Errorf("alias = %q, want %q", rt.Name, "b")
		}
		if lit, ok := rt.Val.(*ast.StringLit); !ok || lit.Value != "a" {
			t.Errorf("value = %#v, want StringLit 'a'", rt.Val)
		}
	})
	t.Run("run then explicit AS alias", func(t *testing.T) {
		item := firstSelectItem(t, `SELECT 'a' 'b' AS c`)
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			t.Fatalf("expected *ast.ResTarget, got %T", item)
		}
		if lit, ok := rt.Val.(*ast.StringLit); !ok || lit.Value != "ab" {
			t.Errorf("value = %#v, want StringLit 'ab'", rt.Val)
		}
	})
	t.Run("int literal keeps string alias", func(t *testing.T) {
		item := firstSelectItem(t, `SELECT 1 'x'`)
		rt, ok := item.(*ast.ResTarget)
		if !ok || rt.Name != "x" {
			t.Fatalf("SELECT 1 'x': expected alias 'x', got %#v", item)
		}
	})
	// Hex, bit, and temporal literals never join a run: MySQL parses the
	// trailing string as the item alias (oracle: SELECT x'41' 'b' yields value
	// 'A' named 'b', and CONCAT(x'41' 'b') is not a concatenation).
	for _, tc := range []struct {
		sql  string
		want string
	}{
		{`SELECT x'41' 'b'`, "b"},
		{`SELECT b'01000001' 'c'`, "c"},
		{`SELECT DATE '2024-01-01' 'd'`, "d"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			item := firstSelectItem(t, tc.sql)
			rt, ok := item.(*ast.ResTarget)
			if !ok {
				t.Fatalf("expected aliased item, got %T", item)
			}
			if rt.Name != tc.want {
				t.Errorf("alias = %q, want %q", rt.Name, tc.want)
			}
			if _, isStr := rt.Val.(*ast.StringLit); isStr {
				t.Errorf("value folded into a StringLit; hex/bit/temporal must not concatenate")
			}
		})
	}
}

// TestAdjacentStringLiteralRejects pins the forms MySQL rejects (1064 on both
// 8.0.32 and 5.7.25): a continuation segment cannot carry a charset introducer,
// a string cannot continue a hex literal inside an expression argument, and the
// TEXT_STRING_sys contexts (COMMENT, ENUM members) never concatenate.
func TestAdjacentStringLiteralRejects(t *testing.T) {
	cases := []string{
		`SELECT 'a' _utf8mb4'b'`,
		`SELECT CONCAT('a' x'41')`,
		`CREATE TABLE t (a int) COMMENT 'x' 'y'`,
		`CREATE TABLE t (a int COMMENT 'x' 'y')`,
		`CREATE TABLE t (e ENUM('a' 'b'))`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) succeeded; MySQL rejects this form", sql)
			}
		})
	}
}

// TestAdjacentStringLiteralSysProcedureRepro is the real-world shape: the stock
// MySQL 8.0 sys schema's ps_trace_thread body concatenates '\n' with the next
// segment across a line break. The body must parse and BodyText must round-trip
// the source bytes verbatim (routine bodies are stored verbatim by MySQL).
func TestAdjacentStringLiteralSysProcedureRepro(t *testing.T) {
	body := "BEGIN\n    SELECT CONCAT('tmp disk tables: ', 3, '\\n'\n                  'select scan: ', 4, '\\n');\nEND"
	sql := "CREATE PROCEDURE p()\n" + body
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	fn, ok := list.Items[0].(*ast.CreateFunctionStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateFunctionStmt, got %T", list.Items[0])
	}
	if fn.BodyText != body {
		t.Errorf("BodyText not verbatim:\n got: %q\nwant: %q", fn.BodyText, body)
	}

	// Control: the same body with a comma instead of adjacency parses too.
	control := "CREATE PROCEDURE p()\nBEGIN\n    SELECT CONCAT('tmp disk tables: ', 3, '\\n',\n                  'select scan: ', 4, '\\n');\nEND"
	if _, err := Parse(control); err != nil {
		t.Fatalf("control Parse error: %v", err)
	}

	// Trigger bodies share the expression path.
	trigger := "CREATE TRIGGER trg1 BEFORE INSERT ON t1 FOR EACH ROW\nBEGIN\n    SET NEW.a = CONCAT('x' 'y');\nEND"
	if _, err := Parse(trigger); err != nil {
		t.Fatalf("trigger Parse error: %v", err)
	}
}

package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// execCommentSyntax matches the bytes of executable-comment syntax in the test
// inputs below, none of which holds "*/" inside a string or comment.
var execCommentSyntax = regexp.MustCompile(`/\*![0-9]*|\*/`)

// blankExecComments replaces every "/*!NNNNN" and "*/" in sql with as many
// spaces, keeping every other byte at its offset.
func blankExecComments(sql string) string {
	return execCommentSyntax.ReplaceAllStringFunc(sql, func(m string) string {
		return strings.Repeat(" ", len(m))
	})
}

// nodeLocs parses sql and returns "Type [start,end)" for every node that
// carries a Loc, in walk order.
func nodeLocs(t *testing.T, sql string) []string {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	var locs []string
	for _, item := range list.Items {
		ast.Inspect(item, func(n ast.Node) bool {
			v := reflect.Indirect(reflect.ValueOf(n))
			if v.Kind() != reflect.Struct {
				return true
			}
			if f := v.FieldByName("Loc"); f.IsValid() && f.Type() == reflect.TypeOf(ast.Loc{}) {
				loc := f.Interface().(ast.Loc)
				locs = append(locs, fmt.Sprintf("%s [%d,%d)", v.Type().Name(), loc.Start, loc.End))
			}
			return true
		})
	}
	return locs
}

// nodeTexts parses sql and returns "Type: text" for every node whose Loc is a
// valid range of sql, text being the slice of sql the Loc names.
func nodeTexts(t *testing.T, sql string) []string {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	var texts []string
	for _, item := range list.Items {
		ast.Inspect(item, func(n ast.Node) bool {
			v := reflect.Indirect(reflect.ValueOf(n))
			if v.Kind() != reflect.Struct {
				return true
			}
			if f := v.FieldByName("Loc"); f.IsValid() && f.Type() == reflect.TypeOf(ast.Loc{}) {
				loc := f.Interface().(ast.Loc)
				if loc.Start >= 0 && loc.Start <= loc.End && loc.End <= len(sql) {
					texts = append(texts, v.Type().Name()+": "+sql[loc.Start:loc.End])
				}
			}
			return true
		})
	}
	return texts
}

var execCommentLocCases = []string{
	"/*!80000 SELECT * FROM t WHERE n = NULL */",
	"/*! SELECT a FROM t WHERE b = 1 */",
	"SELECT /*!40001 SQL_NO_CACHE */ a FROM t WHERE b = NULL",
	"SELECT a FROM t WHERE /*!50000 b = NULL */ AND c = 2",
	"/*!50000 /*!50100 SELECT a FROM t WHERE b = NULL */ */",
	"/*!*/SELECT a FROM t WHERE b = NULL",
	"SELECT 1;\n/*!80000 SELECT a FROM t WHERE n = NULL */",
	"CREATE TABLE t (a INT) /*!50100 PARTITION BY HASH (a) PARTITIONS 4 */",
	"CREATE TABLE t (a INT /*!80023 INVISIBLE */, b INT)",
	"CREATE DATABASE d /*!40100 DEFAULT CHARACTER SET utf8mb4 */ /*!80016 DEFAULT ENCRYPTION='N' */",
	"/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NULL */",
	"/*!50001 CREATE ALGORITHM=UNDEFINED */ /*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */ /*!50001 VIEW v AS select a from t where b = 1 */",
	"/*!50003 CREATE*/ /*!50020 DEFINER=`root`@`localhost`*/ /*!50003 PROCEDURE p() BEGIN SELECT a FROM t WHERE b = NULL; END */",
}

// TestExecCommentLoc_SameAsWhitespace checks that executable-comment syntax
// shifts no location: every node parsed from inside /*! ... */ carries the
// same Loc it carries when each "/*!NNNNN" and "*/" is replaced by spaces of
// the same width, so every Loc is an offset into the text given to the parser.
func TestExecCommentLoc_SameAsWhitespace(t *testing.T) {
	for _, sql := range execCommentLocCases {
		t.Run(sql, func(t *testing.T) {
			blank := blankExecComments(sql)
			got, want := nodeLocs(t, sql), nodeLocs(t, blank)
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Errorf("locations differ from the same text with the comment syntax blanked (%q)\n got: %v\nwant: %v", blank, got, want)
			}
		})
	}
}

// TestExecCommentLoc_SlicesSource checks that the Loc of a node read from
// inside an executable comment slices the input to that node's source text.
func TestExecCommentLoc_SlicesSource(t *testing.T) {
	cases := []struct {
		sql  string
		want []string
	}{
		{
			sql:  "/*!80000 SELECT * FROM t WHERE n = NULL */",
			want: []string{"SelectStmt: SELECT * FROM t WHERE n = NULL", "TableRef: t"},
		},
		{
			sql:  "/*! SELECT a FROM t1 JOIN t2 ON t1.id = t2.id */",
			want: []string{"SelectStmt: SELECT a FROM t1 JOIN t2 ON t1.id = t2.id", "TableRef: t1", "TableRef: t2"},
		},
		{
			sql:  "/*!50000 /*!50100 SELECT a FROM t WHERE b = ABS(c) */ */",
			want: []string{"SelectStmt: SELECT a FROM t WHERE b = ABS(c)", "FuncCallExpr: ABS(c)"},
		},
		{
			sql:  "SELECT 1;\n/*!80000 SELECT a FROM t WHERE n = NULL */",
			want: []string{"SelectStmt: SELECT 1", "SelectStmt: SELECT a FROM t WHERE n = NULL"},
		},
		{
			sql:  "CREATE TABLE t (a INT) /*!50100 PARTITION BY HASH (a) PARTITIONS 4 */",
			want: []string{"PartitionClause: PARTITION BY HASH (a) PARTITIONS 4", "ColumnRef: a", "ColumnDef: a INT"},
		},
		{
			sql:  "/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER trg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = NULL */",
			want: []string{"TableRef: t", "SetStmt: SET NEW.a = NULL"},
		},
		{
			sql:  "/*!50001 CREATE ALGORITHM=UNDEFINED */ /*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */ /*!50001 VIEW v AS select a from t where b = 1 */",
			want: []string{"TableRef: v", "SelectStmt: select a from t where b = 1", "TableRef: t"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			got := nodeTexts(t, tc.sql)
			have := make(map[string]bool, len(got))
			for _, s := range got {
				have[s] = true
			}
			for _, w := range tc.want {
				if !have[w] {
					t.Errorf("no node slices to %q; node texts: %q", w, got)
				}
			}
		})
	}
}

// TestExecCommentLoc_Tokens checks token offsets directly: each token read
// from inside an executable comment slices the input to its own text, and the
// comment syntax moves no token.
func TestExecCommentLoc_Tokens(t *testing.T) {
	sql := "/*!80000 SELECT a, 'x' FROM t */ WHERE /*! b = NULL */"
	want := []string{"SELECT", "a", ",", "'x'", "FROM", "t", "WHERE", "b", "=", "NULL"}
	toks := Tokenize(sql)
	var got []string
	for _, tok := range toks {
		got = append(got, sql[tok.Loc:tok.End])
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("token texts = %q, want %q", got, want)
	}
	blank := Tokenize(blankExecComments(sql))
	if len(blank) != len(toks) {
		t.Fatalf("token count %d, want %d", len(toks), len(blank))
	}
	for i := range toks {
		if toks[i] != blank[i] {
			t.Errorf("token %d = %+v, want %+v", i, toks[i], blank[i])
		}
	}
}

// TestExecCommentLoc_BodyText pins the body text captured across executable
// comments, as mysqldump writes them: it holds the SQL inside the comments,
// without the comment syntax.
func TestExecCommentLoc_BodyText(t *testing.T) {
	view := "/*!50001 CREATE ALGORITHM=UNDEFINED */\n/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */\n/*!50001 VIEW `v` AS select 1 AS `1` */"
	trigger := "/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `trg` BEFORE INSERT ON `t` FOR EACH ROW SET NEW.a = 1 */"
	procedure := "/*!50003 CREATE*/ /*!50020 DEFINER=`root`@`localhost`*/ /*!50003 PROCEDURE p() BEGIN SELECT 1; END */"
	event := "/*!50106 CREATE*/ /*!50117 DEFINER=`root`@`localhost`*/ /*!50106 EVENT e ON SCHEDULE EVERY 1 HOUR DO SELECT 1 */"

	parseOne := func(sql string) ast.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		if len(list.Items) != 1 {
			t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(list.Items))
		}
		return list.Items[0]
	}

	if got := parseOne(view).(*ast.CreateViewStmt).SelectText; got != "select 1 AS `1`" {
		t.Errorf("view SelectText = %q", got)
	}
	if got := parseOne(trigger).(*ast.CreateTriggerStmt).BodyText; got != "SET NEW.a = 1 " {
		t.Errorf("trigger BodyText = %q", got)
	}
	if got := parseOne(procedure).(*ast.CreateFunctionStmt).BodyText; got != "BEGIN SELECT 1; END " {
		t.Errorf("procedure BodyText = %q", got)
	}
	ev := parseOne(event).(*ast.CreateEventStmt)
	if ev.BodyText != "SELECT 1 " || ev.Schedule.RawText != "EVERY 1 HOUR" {
		t.Errorf("event BodyText = %q, schedule RawText = %q", ev.BodyText, ev.Schedule.RawText)
	}
}

// TestExecCommentLoc_ErrorPosition checks that errors inside an executable
// comment report offsets into the input.
func TestExecCommentLoc_ErrorPosition(t *testing.T) {
	cases := []struct {
		sql  string
		msg  string
		want int
	}{
		// A syntax error inside the comment body names the token's offset.
		{"/*!50000 SELECT FROM t */", "syntax error", 16},
		{"SELECT 1; /*!80000 SELECT a FROM WHERE */", "expected identifier", 33},
		// A string or a line comment that swallows the comment's closing */
		// leaves the comment open, which MySQL rejects; the error names the
		// comment's start.
		{"SELECT 'x' /*! , 'a*/'", "unterminated comment", 11},
		{"SELECT 1 /*!50000 + 1 -- c */\n", "unterminated comment", 9},
		{"/*!50000 SELECT 'a*/' */", "unterminated comment", 0},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			if !strings.Contains(pe.Message, tc.msg) || pe.Position != tc.want {
				t.Errorf("Parse(%q) = %q at %d, want %q at %d", tc.sql, pe.Message, pe.Position, tc.msg, tc.want)
			}
		})
	}
}

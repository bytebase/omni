package analysis

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/clickhouse/parser"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		offsets [][2]int
		empty   []bool
	}{
		{
			name:   "empty string",
			script: "",
			want:   nil,
		},
		{
			name:    "single statement",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			offsets: [][2]int{{0, 9}},
			empty:   []bool{false},
		},
		{
			name:    "multi-line statement",
			script:  "SELECT\n  1;",
			want:    []string{"SELECT\n  1;"},
			offsets: [][2]int{{0, 11}},
			empty:   []bool{false},
		},
		{
			name:    "two statements on one line",
			script:  "SELECT 1; SELECT 2;",
			want:    []string{"SELECT 1;", " SELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 19}},
			empty:   []bool{false, false},
		},
		{
			name:    "multiple statements with comments and standalone semicolon",
			script:  "SELECT 1;\nSELECT 2\n;\n-- a comment\nSELECT 3",
			want:    []string{"SELECT 1;", "\nSELECT 2\n;", "\n-- a comment\nSELECT 3"},
			offsets: [][2]int{{0, 9}, {9, 20}, {20, 42}},
			empty:   []bool{false, false, false},
		},
		{
			name:    "multi-byte UTF-8",
			script:  "SELECT '中';\nSELECT '文';",
			want:    []string{"SELECT '中';", "\nSELECT '文';"},
			offsets: [][2]int{{0, 13}, {13, 27}},
			empty:   []bool{false, false},
		},
		{
			name:    "leading newlines preserved in range",
			script:  "\n  \n  SELECT 1;",
			want:    []string{"\n  \n  SELECT 1;"},
			offsets: [][2]int{{0, 15}},
			empty:   []bool{false},
		},
		{
			name:    "CRLF multiple statements",
			script:  "SELECT 1;\r\nSELECT 2;",
			want:    []string{"SELECT 1;", "\r\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 20}},
			empty:   []bool{false, false},
		},
		{
			name:    "comment only script is preserved as empty",
			script:  "\n-- only a comment\n",
			want:    []string{"\n-- only a comment\n"},
			offsets: [][2]int{{0, 19}},
			empty:   []bool{true},
		},
		{
			name:    "lone semicolon is preserved as empty",
			script:  ";\n",
			want:    []string{";", "\n"},
			offsets: [][2]int{{0, 1}, {1, 2}},
			empty:   []bool{true, true},
		},
		{
			name:    "inline block comment before statement",
			script:  "/* a comment */ SELECT 1;",
			want:    []string{"/* a comment */ SELECT 1;"},
			offsets: [][2]int{{0, 25}},
			empty:   []bool{false},
		},
		{
			name:    "multiline block comment closing mid-line",
			script:  "/* a\ncomment */ SELECT 1;",
			want:    []string{"/* a\ncomment */ SELECT 1;"},
			offsets: [][2]int{{0, 25}},
			empty:   []bool{false},
		},
		{
			name:    "nested block comments with semicolons inside",
			script:  "/* outer /* inner ; */ outer ; */ SELECT 1; SELECT 2;",
			want:    []string{"/* outer /* inner ; */ outer ; */ SELECT 1;", " SELECT 2;"},
			offsets: [][2]int{{0, 43}, {43, 53}},
			empty:   []bool{false, false},
		},
		{
			name:    "dollar-quoted strings containing semicolons and keywords",
			script:  "SELECT $$a; DROP TABLE t;$$, $tag$x;y$tag$; SELECT 2;",
			want:    []string{"SELECT $$a; DROP TABLE t;$$, $tag$x;y$tag$;", " SELECT 2;"},
			offsets: [][2]int{{0, 43}, {43, 53}},
			empty:   []bool{false, false},
		},
		{
			name:    "single, double, and backtick quotes containing semicolons",
			script:  `SELECT 'a;b', 'it\'s;ok', "c;d", ` + "`e;f`" + `; SELECT 2;`,
			want:    []string{`SELECT 'a;b', 'it\'s;ok', "c;d", ` + "`e;f`" + `;`, " SELECT 2;"},
			offsets: [][2]int{{0, 39}, {39, 49}},
			empty:   []bool{false, false},
		},
		{
			name:    "same-line trailing comment stays attached after semicolon",
			script:  "SELECT 1; -- note\nSELECT 2; # hash",
			want:    []string{"SELECT 1; -- note", "\nSELECT 2; # hash"},
			offsets: [][2]int{{0, 17}, {17, 34}},
			empty:   []bool{false, false},
		},
		{
			name:    "unterminated string stays non-empty in its segment",
			script:  "INSERT INTO q (id) VALUES (5);\nSELECT 'unterminated",
			want:    []string{"INSERT INTO q (id) VALUES (5);", "\nSELECT 'unterminated"},
			offsets: [][2]int{{0, 30}, {30, 51}},
			empty:   []bool{false, false},
		},
		{
			name:    "insert format csv preserves semicolons in inline payload",
			script:  "INSERT INTO t FORMAT CSV\nhello;world\n",
			want:    []string{"INSERT INTO t FORMAT CSV\nhello;world\n"},
			offsets: [][2]int{{0, 37}},
			empty:   []bool{false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Split(tc.script)
			if err != nil {
				t.Fatalf("Split: %v", err)
			}
			var texts []string
			prevEnd := 0
			for i, s := range got {
				texts = append(texts, s.Text)
				if s.Start != prevEnd {
					t.Errorf("statement %d Start = %d, want contiguous %d", i, s.Start, prevEnd)
				}
				if s.Text != tc.script[s.Start:s.End] {
					t.Errorf("statement text %q != script[%d:%d] (%q)", s.Text, s.Start, s.End, tc.script[s.Start:s.End])
				}
				prevEnd = s.End
			}
			if prevEnd != len(tc.script) {
				t.Errorf("final End = %d, want len(script) = %d", prevEnd, len(tc.script))
			}
			if !reflect.DeepEqual(texts, tc.want) {
				t.Fatalf("statements = %q, want %q", texts, tc.want)
			}
			for i, off := range tc.offsets {
				if got[i].Start != off[0] || got[i].End != off[1] {
					t.Errorf("statement %d offsets = [%d,%d), want [%d,%d)", i, got[i].Start, got[i].End, off[0], off[1])
				}
			}
			for i, wantEmpty := range tc.empty {
				if got[i].Empty != wantEmpty {
					t.Errorf("statement %d Empty = %v, want %v", i, got[i].Empty, wantEmpty)
				}
			}
		})
	}
}

func TestExplainStatement(t *testing.T) {
	for _, tc := range []struct {
		statement, format, want string
	}{
		{"SELECT 1", "", "EXPLAIN SELECT 1"},
		{"SELECT 1", "TEXT", "EXPLAIN SELECT 1"},
		{"SELECT 1", "XML", "EXPLAIN SELECT 1"},
		{"SELECT 1", "JSON", "EXPLAIN json = 1, indexes = 1 SELECT 1"},
		{"SELECT 1", "json", "EXPLAIN json = 1, indexes = 1 SELECT 1"},
		{"EXPLAIN SELECT 1", "JSON", "EXPLAIN SELECT 1"},
		{"EXPLAIN PIPELINE SELECT 1", "JSON", "EXPLAIN PIPELINE SELECT 1"},
		{"/* c */ explain json = 1 SELECT 1", "JSON", "/* c */ explain json = 1 SELECT 1"},
		{"/* outer /* inner */ outer */ EXPLAIN SELECT 1", "JSON", "/* outer /* inner */ outer */ EXPLAIN SELECT 1"},
		{"# comment\nEXPLAIN SELECT 1", "JSON", "# comment\nEXPLAIN SELECT 1"},
		{"explained_view_query()", "", "EXPLAIN explained_view_query()"},
		{"`EXPLAIN` SELECT 1", "JSON", "EXPLAIN json = 1, indexes = 1 `EXPLAIN` SELECT 1"},
		{"/* a /* b */ EXPLAIN */ SELECT 1", "JSON", "EXPLAIN json = 1, indexes = 1 /* a /* b */ EXPLAIN */ SELECT 1"},
		{"/* unclosed EXPLAIN SELECT 1", "JSON", "EXPLAIN json = 1, indexes = 1 /* unclosed EXPLAIN SELECT 1"},
	} {
		if got := ExplainStatement(tc.statement, tc.format); got != tc.want {
			t.Errorf("ExplainStatement(%q, %q) = %q, want %q", tc.statement, tc.format, got, tc.want)
		}
	}
}

func TestDescribeExplain(t *testing.T) {
	for _, tc := range []struct {
		statement, format string
		executed, ok      bool
	}{
		{statement: "SELECT 1"},
		{statement: "  -- leading comment\nSELECT 1"},
		{statement: "/* a /* b */ EXPLAIN */ SELECT 1"},
		{statement: "/* unclosed EXPLAIN SELECT 1"},
		{statement: "`EXPLAIN` SELECT 1"},
		{statement: "EXPLAIN SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN SELECT 1 AS json", format: "TEXT", ok: true},
		{statement: "EXPLAIN SELECT json = 1 FROM t", format: "TEXT", ok: true},
		{statement: "EXPLAIN AST SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN SYNTAX SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN PIPELINE SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN ESTIMATE SELECT 1 FROM t", format: "TEXT", ok: true},
		{statement: "EXPLAIN TABLE OVERRIDE t", format: "TEXT", ok: true},
		{statement: "EXPLAIN PLAN json = 1 SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN json = 0 SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = false SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = off SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = no SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = 1, json = 0 SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = 1, indexes = 1 SELECT 1", format: "JSON", ok: true},
		{statement: "/* outer /* inner */ */ EXPLAIN PLAN json = 1 SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN QUERY TREE SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN JSON SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN ANALYZE SELECT 1", format: "TEXT", executed: true, ok: true},
		{statement: "EXPLAIN ANALYZE json = 1 SELECT 1", format: "JSON", executed: true, ok: true},
		{statement: "EXPLAIN ANALYZE json = 0 SELECT 1", format: "TEXT", executed: true, ok: true},
	} {
		format, executed, ok := DescribeExplain(tc.statement)
		if ok != tc.ok || format != tc.format || executed != tc.executed {
			t.Errorf("DescribeExplain(%q) = (%q, %v, %v), want (%q, %v, %v)",
				tc.statement, format, executed, ok, tc.format, tc.executed, tc.ok)
		}
		if got := StartsWithExplain(tc.statement); got != tc.ok {
			t.Errorf("StartsWithExplain(%q) = %v, want %v", tc.statement, got, tc.ok)
		}
	}
}

func TestStatementWithResultLimit(t *testing.T) {
	if got, want := StatementWithResultLimit("SELECT 1;  ", 5), "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"; got != want {
		t.Errorf("limit = %q, want %q", got, want)
	}
	for _, tc := range []struct{ statement, want string }{
		{"SELECT 1 -- note", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; -- note", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; /* note */", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; /* outer /* inner */ outer */", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"/* outer /* inner */ */ SELECT 1; -- tail", "WITH result AS (/* outer /* inner */ */ SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT '你好' -- note", "WITH result AS (SELECT '你好') SELECT * FROM result LIMIT 5;"},
		{"SELECT 1 # note", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; # note", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT $doc$--$doc$ -- note", "WITH result AS (SELECT $doc$--$doc$) SELECT * FROM result LIMIT 5;"},
		{"SELECT $$--$$ # note", "WITH result AS (SELECT $$--$$) SELECT * FROM result LIMIT 5;"},
		{"SELECT `a--b`", "WITH result AS (SELECT `a--b`) SELECT * FROM result LIMIT 5;"},
		{"WITH x AS (SELECT 1) SELECT * FROM x", "WITH result AS (WITH x AS (SELECT 1) SELECT * FROM x) SELECT * FROM result LIMIT 5;"},
		{"(SELECT 1) UNION ALL (SELECT 2)", "WITH result AS ((SELECT 1) UNION ALL (SELECT 2)) SELECT * FROM result LIMIT 5;"},
		{"FROM t SELECT a", "WITH result AS (FROM t SELECT a) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1 SETTINGS max_threads = 2", "WITH result AS (SELECT 1 SETTINGS max_threads = 2) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1 FORMAT JSON", "SELECT 1 FORMAT JSON"},
		{"SELECT 1 INTO OUTFILE '/tmp/out.csv'", "SELECT 1 INTO OUTFILE '/tmp/out.csv'"},
		{"WITH x AS (SELECT 1) INSERT INTO t SELECT * FROM x", "WITH x AS (SELECT 1) INSERT INTO t SELECT * FROM x"},
		{"SELECT 1; SELECT 2;", "SELECT 1; SELECT 2;"},
		{"INSERT INTO t VALUES (1)", "INSERT INTO t VALUES (1)"},
		{"SHOW TABLES", "SHOW TABLES"},
		{"DESCRIBE TABLE t", "DESCRIBE TABLE t"},
		{"EXPLAIN SELECT 1", "EXPLAIN SELECT 1"},
		{"DELETE FROM t WHERE a = 1", "DELETE FROM t WHERE a = 1"},
	} {
		if got := StatementWithResultLimit(tc.statement, 5); got != tc.want {
			t.Errorf("StatementWithResultLimit(%q, 5) = %q, want %q", tc.statement, got, tc.want)
		}
	}
	for _, statement := range []string{"SELECT '--'", `SELECT "--"`, "SELECT `a--b`", "SELECT 1 /* -- */", "SELECT 1 -- note\n"} {
		if got := StatementWithResultLimit(statement, 5); !strings.HasSuffix(got, "LIMIT 5;") {
			t.Errorf("StatementWithResultLimit(%q, 5) = %q", statement, got)
		}
	}
}

func TestWrapResultLimit(t *testing.T) {
	for _, tc := range []struct {
		statement string
		limit     int
		want      string
	}{
		{"SELECT * FROM t", 10, "WITH result AS (SELECT * FROM t) SELECT * FROM result LIMIT 10;"},
		{"SELECT * FROM t;", 10, "WITH result AS (SELECT * FROM t) SELECT * FROM result LIMIT 10;"},
		{"DELETE FROM t", 5, "WITH result AS (DELETE FROM t) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; -- note", 5, "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; /* outer /* inner */ */", 5, "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
	} {
		if got := WrapResultLimit(tc.statement, tc.limit); got != tc.want {
			t.Errorf("WrapResultLimit(%q, %d) = %q, want %q", tc.statement, tc.limit, got, tc.want)
		}
	}
}

func TestValidateQuery(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{"SELECT 1", true},
		{" select * from t", true},
		{"select* from t", true},
		{"SELECT(1)", true},
		{"SELECT*FROM q", true},
		{"SELECT 'x'", true},
		{"(SELECT 1)", true},
		{"(SELECT 1) UNION ALL (SELECT 2)", true},
		{"FROM t SELECT a, b WHERE a > 0", true},
		{"SELECT 1 FORMAT JSON", true},
		{"SELECT 1 SETTINGS max_threads = 2", true},
		{"SELECT 1 INTO OUTFILE '/tmp/out.csv'", false},
		{"EXPLAIN SELECT 1", true},
		{"explain select * from t;", true},
		{"EXPLAIN AST SELECT 1", true},
		{"EXPLAIN AST DROP TABLE t", true},
		{"EXPLAIN SYNTAX SELECT 1", true},
		{"EXPLAIN PLAN json = 1 SELECT 1", true},
		{"EXPLAIN PIPELINE SELECT 1", true},
		{"EXPLAIN ESTIMATE SELECT 1 FROM t", true},
		{"EXPLAIN QUERY TREE SELECT 1", true},
		{"EXPLAIN TABLE OVERRIDE t", true},
		{"EXPLAIN ANALYZE SELECT 1", false},
		{"explain    analyze select * from t", false},
		{"EXPLAIN ANALYZE json = 1 SELECT 1", false},
		{"WITH x AS (SELECT 1) SELECT * FROM x", true},
		{"WITH RECURSIVE x AS (SELECT 1) SELECT * FROM x", true},
		{"WITH 1 AS `x` SELECT 2", true},
		{"WITH 1 AS x SELECT 'foo' AS y", true},
		{"WITH 1 AS x SELECT [] AS y", true},
		{"WITH 1 AS update SELECT update", true},
		{"WITH 1 AS insert SELECT insert", true},
		{"WITH 1 AS delete SELECT delete", true},
		{"WITH 1 AS select SELECT select", true},
		{"INSERT INTO t VALUES (1)", false},
		{"INSERT INTO t FORMAT CSV\nhello;world\n\"unclosed", false},
		{"UPDATE t SET a = 1 WHERE b = 2", false},
		{"DELETE FROM t WHERE a = 1", false},
		{"create table t (a int);", false},
		{"create table t (a int) ENGINE = Memory;", false},
		{"ALTER TABLE t ADD COLUMN c Int32", false},
		{"DROP TABLE t", false},
		{"TRUNCATE TABLE t", false},
		{"RENAME TABLE a TO b", false},
		{"USE mydb", false},
		{"SET max_threads = 4", false},
		{"SYSTEM FLUSH LOGS", false},
		{"KILL QUERY WHERE query_id = '1'", false},
		{"OPTIMIZE TABLE t FINAL", false},
		{"GRANT SELECT ON t TO u", false},
		{"SELECT 'DROP TABLE t' AS x", true},
		{"select * from t where a = 'klasjdfkljsa$tag$; -- lkjdlkfajslkdfj'", true},
		{"With t as (select * from t1 where a = 'insert'), tx as (select * from \"delete\") /* UPDATE */ select \"update\" from t;", true},
		{"-- SELECT\nDROP TABLE t", false},
		{"SELECT 1; DROP TABLE t", false},
		{"SELECT 1; SELECT 2", false},
		{"/* /* */ */ SELECT 1", true},
		{"SHOW TABLES", true},
		{"SHOW DATABASES", true},
		{"SHOW CREATE TABLE t", true},
		{"SHOW TABLES INTO OUTFILE '/tmp/out'", false},
		{"DESCRIBE TABLE q", true},
		{"DESC q", true},
		{"EXISTS TABLE q", true},
		{"CHECK TABLE q", true},
		{"SHOW TABLES; DROP TABLE q", false},
		{"SELECT `a--b` FROM t", true},
		{`SELECT "a\"b" FROM q`, true},
		{"SELECT `a;b` FROM t", true},
		{"# note\nSELECT 1", true},
		{"SELECT 1 # note", true},
		{"SELECT $doc$it's text$doc$", true},
		{"SELECT $$a; DROP TABLE t;$$", true},
		{"", false},
		{"-- comment only", false},
		{";", false},
	} {
		readOnly, returnsData, err := ValidateQuery(tc.statement)
		if err != nil {
			t.Errorf("ValidateQuery(%q): %v", tc.statement, err)
			continue
		}
		if readOnly != tc.want || returnsData != tc.want {
			t.Errorf("ValidateQuery(%q) = (%v, %v), want (%v, %v)", tc.statement, readOnly, returnsData, tc.want, tc.want)
		}
		if got := CheckStatementWithoutQuotedTextAndComment(tc.statement); got != tc.want {
			t.Errorf("CheckStatementWithoutQuotedTextAndComment(%q) = %v, want %v", tc.statement, got, tc.want)
		}
	}

	for _, bad := range []struct {
		name    string
		sql     string
		wantPos int
		wantSub string
	}{
		{"unterminated single-quoted string", "SELECT 'unterminated", 7, "invalid string"},
		{"unterminated dollar-quoted string", "SELECT $doc$unterminated", 7, "unterminated dollar-quoted"},
		{"unterminated backtick ident", "SELECT `unterminated", 7, "unterminated quoted identifier"},
		{"unterminated double-quoted ident", `SELECT "unterminated`, 7, "unterminated quoted identifier"},
		{"unterminated block comment", "/* unclosed SELECT 1", 0, "unterminated block comment"},
		{"unterminated nested block comment", "/* /* */ SELECT 1", 0, "unterminated block comment"},
		{"with insert rejected as syntax error", "WITH x AS (SELECT 1) INSERT INTO t SELECT * FROM x", 21, "unexpected token"},
		{"with update rejected as syntax error", "WITH x AS (SELECT 1) UPDATE t SET a = 1", 21, "unexpected token"},
		{"with delete rejected as syntax error", "WITH x AS (SELECT 1) DELETE FROM t", 21, "unexpected token"},
		{"syntax error propagates parse error", "SELECT 1 FROM", 13, "expected"},
		{"syntax error after leading semicolon keeps script offset", "; SELECT 1 FROM", 15, "expected"},
		{"unparsed select propagates parse error", "SELECT 1e1 2e1", 11, "unexpected token"},
	} {
		t.Run(bad.name, func(t *testing.T) {
			ro, rd, err := ValidateQuery(bad.sql)
			if err == nil {
				t.Fatalf("ValidateQuery(%q) = (%v, %v, nil), want error", bad.sql, ro, rd)
			}
			if ro || rd {
				t.Errorf("ValidateQuery(%q) returned (%v, %v) on error", bad.sql, ro, rd)
			}
			if !strings.Contains(err.Error(), bad.wantSub) {
				t.Errorf("ValidateQuery(%q) error = %q, want substring %q", bad.sql, err.Error(), bad.wantSub)
			}
			var pe *parser.ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("ValidateQuery(%q) error type = %T, want *parser.ParseError", bad.sql, err)
			}
			if pe.Position != bad.wantPos {
				t.Errorf("ValidateQuery(%q) Position = %d, want %d", bad.sql, pe.Position, bad.wantPos)
			}
			if got := CheckStatementWithoutQuotedTextAndComment(bad.sql); got {
				t.Errorf("CheckStatementWithoutQuotedTextAndComment(%q) = true, want false", bad.sql)
			}
		})
	}
}

func TestStatementReadsRows(t *testing.T) {
	for _, tc := range []struct {
		statement string
		rows      bool
	}{
		{"SELECT 1", true},
		{"SHOW TABLES", true},
		{"DESCRIBE TABLE q", true},
		{"EXISTS TABLE q", true},
		{"CHECK TABLE q", true},
		{"EXPLAIN SELECT 1", true},
		{"EXPLAIN AST SELECT 1", true},
		{"EXPLAIN ANALYZE SELECT 1", true},
		{"-- c\nEXPLAIN ANALYZE SELECT 1", true},
		{"/* outer /* inner */ */ EXPLAIN ANALYZE SELECT 1", true},
		{"EXPLAIN ANALYZE json = 1 SELECT 1", true},
		{"EXPLAIN ANALYZE SELECT 1; DROP TABLE t", false},
		{"INSERT INTO t VALUES (1)", false},
		{"INSERT INTO t FORMAT CSV\nhello;world\n\"unclosed", false},
		{"SELECT 1; DROP TABLE t", false},
		{"USE mydb", false},
		{"SYSTEM FLUSH LOGS", false},
	} {
		got, err := StatementReadsRows(tc.statement)
		if err != nil {
			t.Errorf("StatementReadsRows(%q): %v", tc.statement, err)
			continue
		}
		if got != tc.rows {
			t.Errorf("StatementReadsRows(%q) = %v, want %v", tc.statement, got, tc.rows)
		}
	}
	for _, bad := range []struct {
		sql     string
		wantPos int
	}{
		{"SELECT 'unterminated", 7},
		{"/* /* */ SELECT 1", 0},
		{"WITH x AS (SELECT 1) DELETE FROM t", 21},
		{"SELECT 1e1 2e1", 11},
		{"; SELECT 1 FROM", 15},
	} {
		got, err := StatementReadsRows(bad.sql)
		if err == nil {
			t.Errorf("StatementReadsRows(%q) = (%v, nil), want error", bad.sql, got)
			continue
		}
		if got {
			t.Errorf("StatementReadsRows(%q) returned true on error", bad.sql)
		}
		var pe *parser.ParseError
		if !errors.As(err, &pe) || pe.Position != bad.wantPos {
			t.Errorf("StatementReadsRows(%q) err = %v, want *parser.ParseError at %d", bad.sql, err, bad.wantPos)
		}
	}
}

package analysis

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		offsets [][2]int
		empty   []bool
		err     string
	}{
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
			name:    "multiple statements with comments and standalone semicolon",
			script:  "SELECT 1;\nSELECT 2\n;\n-- a comment\nSELECT 3",
			want:    []string{"SELECT 1;", "\nSELECT 2\n;", "\n-- a comment\nSELECT 3"},
			offsets: [][2]int{{0, 9}, {9, 20}, {20, 42}},
		},
		{
			name:    "multi-byte UTF-8",
			script:  "SELECT '中';\nSELECT '文';",
			want:    []string{"SELECT '中';", "\nSELECT '文';"},
			offsets: [][2]int{{0, 13}, {13, 27}},
		},
		{
			name:    "leading newlines preserved in range",
			script:  "\n  \n  SELECT 1;",
			want:    []string{"\n  \n  SELECT 1;"},
			offsets: [][2]int{{0, 15}},
		},
		{
			name:    "CRLF multiple statements",
			script:  "SELECT 1;\r\nSELECT 2;",
			want:    []string{"SELECT 1;", "\r\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 20}},
		},
		{
			name:   "comment only produces no statements",
			script: "\n-- only a comment\n",
			want:   nil,
		},
		{
			name:   "lone semicolon survives as content",
			script: ";\n",
			want:   []string{";"},
			empty:  []bool{false},
		},
		{
			name:   "delimiter block",
			script: "DELIMITER ;;\nCREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\nEND;\nDELIMITER ;",
			want:   []string{"DELIMITER ;;\nCREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\nEND;"},
		},
		{
			name:   "block comment closing mid-line errors",
			script: "/* a comment */ SELECT 1;",
			err:    "`*/` must be the end of the line",
		},
		{
			name:   "multiline block comment closing mid-line errors",
			script: "/* a\ncomment */ SELECT 1;",
			err:    "`*/` must be the end of the line",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Split(tc.script)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Split error = %v, want %q", err, tc.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Split: %v", err)
			}
			var texts []string
			for _, s := range got {
				texts = append(texts, s.Text)
				if s.Text != tc.script[s.Start:s.End] {
					t.Errorf("statement text %q != script[%d:%d] (%q)", s.Text, s.Start, s.End, tc.script[s.Start:s.End])
				}
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
		{"# comment\nEXPLAIN SELECT 1", "JSON", "# comment\nEXPLAIN SELECT 1"},
		{"explained_view_query()", "", "EXPLAIN explained_view_query()"},
		{"/* a /* b */ EXPLAIN */ SELECT 1", "JSON", "EXPLAIN json = 1, indexes = 1 /* a /* b */ EXPLAIN */ SELECT 1"},
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
		{statement: "EXPLAIN SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN SELECT 1 AS json", format: "TEXT", ok: true},
		{statement: "EXPLAIN PLAN json = 1 SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN json = 0 SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN json = 1, indexes = 1 SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN QUERY TREE SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN JSON SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN ANALYZE SELECT 1", format: "TEXT", executed: true, ok: true},
		{statement: "EXPLAIN ANALYZE json = 1 SELECT 1", format: "JSON", executed: true, ok: true},
	} {
		format, executed, ok := DescribeExplain(tc.statement)
		if ok != tc.ok || format != tc.format || executed != tc.executed {
			t.Errorf("DescribeExplain(%q) = (%q, %v, %v), want (%q, %v, %v)",
				tc.statement, format, executed, ok, tc.format, tc.executed, tc.ok)
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
		{"SELECT '你好' -- note", "WITH result AS (SELECT '你好') SELECT * FROM result LIMIT 5;"},
		{"SELECT 1 # note", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT 1; # note", "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 5;"},
		{"SELECT $doc$--$doc$ -- note", "WITH result AS (SELECT $doc$--$doc$) SELECT * FROM result LIMIT 5;"},
		{"SELECT `a--b`", "WITH result AS (SELECT `a--b`) SELECT * FROM result LIMIT 5;"},
		{"INSERT INTO t VALUES (1)", "INSERT INTO t VALUES (1)"},
		{"SHOW TABLES", "SHOW TABLES"},
		{"DELETE FROM t", "DELETE FROM t"},
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
		{"EXPLAIN SELECT 1", true},
		{"explain select * from t;", true},
		{"EXPLAIN ANALYZE SELECT 1", false},
		{"explain    analyze select * from t", false},
		{"WITH x AS (SELECT 1) SELECT * FROM x", true},
		{"WITH 1 AS `x` SELECT 2", true},
		{"WITH 1 AS x SELECT 'foo' AS y", true},
		{"WITH 1 AS x SELECT [] AS y", true},
		{"WITH x AS (SELECT 1) INSERT INTO t SELECT * FROM x", false},
		{"WITH x AS (SELECT 1) UPDATE t SET a = 1", false},
		{"WITH x AS (SELECT 1) DELETE FROM t", false},
		{"WITH 1 AS update SELECT update", true},
		{"WITH 1 AS insert SELECT insert", true},
		{"WITH 1 AS delete SELECT delete", true},
		{"WITH 1 AS select SELECT select", true},
		{"INSERT INTO t VALUES (1)", false},
		{"create table t (a int);", false},
		{"SELECT 'DROP TABLE t' AS x", true},
		{"select * from t where a = 'klasjdfkljsa$tag$; -- lkjdlkfajslkdfj'", true},
		{"With t as (select * from t1 where a = 'insert'), tx as (select * from \"delete\") /* UPDATE */ select \"update\" from t;", true},
		{"-- SELECT\nDROP TABLE t", false},
		{"SELECT 1; DROP TABLE t", false},
		{"/* /* */ SELECT 1", true},
		{"SHOW TABLES", true},
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
	} {
		readOnly, returnsData, err := ValidateQuery(tc.statement)
		if err != nil {
			t.Errorf("ValidateQuery(%q): %v", tc.statement, err)
			continue
		}
		if readOnly != tc.want || returnsData != tc.want {
			t.Errorf("ValidateQuery(%q) = (%v, %v), want (%v, %v)", tc.statement, readOnly, returnsData, tc.want, tc.want)
		}
	}

	if _, _, err := ValidateQuery("SELECT 'unterminated"); err == nil {
		t.Error("ValidateQuery of an unterminated string returned no error")
	}
}

func TestStatementReadsRows(t *testing.T) {
	for _, tc := range []struct {
		statement string
		rows      bool
	}{
		{"SELECT 1", true},
		{"SHOW TABLES", true},
		{"EXPLAIN SELECT 1", true},
		{"EXPLAIN ANALYZE SELECT 1", true},
		{"-- c\nEXPLAIN ANALYZE SELECT 1", true},
		{"EXPLAIN ANALYZE json = 1 SELECT 1", true},
		{"INSERT INTO t VALUES (1)", false},
		{"WITH x AS (SELECT 1) DELETE FROM t", false},
		{"SELECT 1; DROP TABLE t", false},
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
	if _, err := StatementReadsRows("SELECT 'unterminated"); err == nil {
		t.Error("StatementReadsRows of an unterminated literal returned no error")
	}
}

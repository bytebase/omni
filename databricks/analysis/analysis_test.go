package analysis

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		offsets [][2]int
		err     string
	}{
		{
			name:    "one statement",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			offsets: [][2]int{{0, 9}},
		},
		{
			name:    "two statements",
			script:  "SELECT 1;\nSELECT 2;",
			want:    []string{"SELECT 1;", "\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 19}},
		},
		{
			name:   "no terminator",
			script: "SELECT 1",
			want:   []string{"SELECT 1"},
		},
		{
			name:   "a comment line belongs to the statement after it",
			script: "-- a comment\nSELECT 1;",
			want:   []string{"-- a comment\nSELECT 1;"},
		},
		{
			name:   "a comment after the last statement belongs to none",
			script: "SELECT 1;\n-- trailing\n",
			want:   []string{"SELECT 1;"},
		},
		{
			name:   "an empty script",
			script: "\n\n",
			want:   []string{},
		},
		{
			name:   "a semicolon inside a multiline string",
			script: "SELECT 'first;\nsecond' AS s;",
			want:   []string{"SELECT 'first;\nsecond' AS s;"},
		},
		{
			name:   "a semicolon inside a multiline string, then another statement",
			script: "SELECT 'first;\nsecond' AS s;\nSELECT 2;",
			want:   []string{"SELECT 'first;\nsecond' AS s;", "\nSELECT 2;"},
		},
		{
			name:   "a semicolon inside a trailing line comment",
			script: "SELECT 1 -- note;\n + 2;",
			want:   []string{"SELECT 1 -- note;\n + 2;"},
		},
		{
			name:   "a semicolon inside a multiline block comment after SQL",
			script: "SELECT 1 /*\n;\n*/ + 2;",
			want:   []string{"SELECT 1 /*\n;\n*/ + 2;"},
		},
		{
			name:   "a trailing comment after a terminator does not merge the next statement",
			script: "SELECT 1; -- note\nSELECT 2;",
			want:   []string{"SELECT 1; -- note", "\nSELECT 2;"},
		},
		{
			name:   "a terminating semicolon before a multiline block comment keeps its boundary",
			script: "SELECT 1; /* note\ncontinued */\nSELECT 2;",
			want:   []string{"SELECT 1; /* note\ncontinued */", "\nSELECT 2;"},
		},
		{
			name:   "a trailing block comment after a terminator does not merge the next statement",
			script: "SELECT 1; /* note */\nSELECT 2;",
			want:   []string{"SELECT 1; /* note */", "\nSELECT 2;"},
		},
		{
			name:   "a block comment that ends its line",
			script: "/* a comment */\nSELECT 1;",
			want:   []string{"/* a comment */\nSELECT 1;"},
		},
		{
			name:   "a block comment whose closing does not end its line",
			script: "/* a comment */ SELECT 1;",
			err:    "`*/` must be the end of the line",
		},
		{
			name:    "CRLF multiple statements",
			script:  "SELECT 1;\r\nSELECT 2;",
			want:    []string{"SELECT 1;", "\r\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 20}},
		},
		{
			name:    "multi-byte UTF-8",
			script:  "SELECT '中';\nSELECT '文';",
			want:    []string{"SELECT '中';", "\nSELECT '文';"},
			offsets: [][2]int{{0, 13}, {13, 27}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			statements, err := Split(tc.script)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Split error = %v, want %q", err, tc.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Split: %v", err)
			}
			texts := make([]string, 0, len(statements))
			for _, stmt := range statements {
				texts = append(texts, stmt.Text)
				if stmt.Text != tc.script[stmt.Start:stmt.End] {
					t.Errorf("statement text %q does not match its offsets", stmt.Text)
				}
			}
			if len(texts) != len(tc.want) {
				t.Fatalf("got %d statements %q, want %d", len(texts), texts, len(tc.want))
			}
			for i := range texts {
				if texts[i] != tc.want[i] {
					t.Errorf("statement %d = %q, want %q", i, texts[i], tc.want[i])
				}
			}
			for i, offsets := range tc.offsets {
				if statements[i].Start != offsets[0] || statements[i].End != offsets[1] {
					t.Errorf("statement %d offsets = [%d,%d), want [%d,%d)",
						i, statements[i].Start, statements[i].End, offsets[0], offsets[1])
				}
			}
		})
	}
}

func TestClassifyStatement(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{statement: "SELECT 1", want: true},
		{statement: "SELECT * FROM t", want: true},
		{statement: "WITH c AS (SELECT 1 AS n) SELECT n FROM c", want: true},
		{statement: "SHOW TABLES", want: true},
		{statement: "DESCRIBE t", want: true},
		{statement: "DESC t", want: true},
		{statement: "EXPLAIN SELECT 1", want: true},
		{statement: "TABLE t", want: true},
		{statement: "VALUES (1)", want: true},
		{statement: "  (SELECT 1)", want: true},
		{statement: "SELECT 'INSERT INTO t VALUES (1)'", want: true},
		{statement: "SELECT 1 -- INSERT\n", want: true},
		{statement: "WITH c AS (INSERT INTO t VALUES (1) RETURNING *) SELECT * FROM c", want: false},
		{statement: "WITH c AS (UPDATE t SET a = 1) SELECT * FROM c", want: false},
		{statement: "WITH c AS (DELETE FROM t) SELECT * FROM c", want: false},
		{statement: "SET spark.sql.ansi.enabled = true", want: false},
		{statement: "INSERT INTO t VALUES (1)", want: false},
		{statement: "UPDATE t SET a = 1", want: false},
		{statement: "DELETE FROM t", want: false},
		{statement: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", want: false},
		{statement: "CREATE TABLE t (a int)", want: false},
		{statement: "", want: false},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			if got := ClassifyStatement(tc.statement); got != tc.want {
				t.Errorf("ClassifyStatement(%q) = %v, want %v", tc.statement, got, tc.want)
			}
		})
	}
}

func TestSanitizeSQL(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      string
		wantErr   bool
	}{
		{
			name:      "a quoted keyword is not a keyword",
			statement: "SELECT 'FROM' AS word",
			want:      "SELECT  x  AS word",
		},
		{
			name:      "a quoted identifier is opaque",
			statement: "SELECT `INSERT` FROM t",
			want:      "SELECT  x  FROM t",
		},
		{
			name:      "comments are removed",
			statement: "-- SELECT FROM t\nSELECT 1 /* WHERE */",
			want:      " \nSELECT 1  ",
		},
		{
			name:      "an unterminated string is reported",
			statement: "SELECT 'x",
			wantErr:   true,
		},
		{
			name:      "an unterminated comment is reported",
			statement: "SELECT 1 /* x",
			wantErr:   true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := sanitizeSQL(tc.statement)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("sanitizeSQL(%q) accepted the statement", tc.statement)
				}
				return
			}
			if err != nil {
				t.Fatalf("sanitizeSQL: %v", err)
			}
			if got != tc.want {
				t.Errorf("sanitizeSQL(%q) = %q, want %q", tc.statement, got, tc.want)
			}
		})
	}
}

func TestExplainStatement(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statement string
		want      string
	}{
		{name: "a query is planned", statement: "SELECT 1", want: "EXPLAIN FORMATTED SELECT 1"},
		{name: "a query with terminator is planned", statement: "SELECT 1;", want: "EXPLAIN FORMATTED SELECT 1;"},
		{name: "a write is planned", statement: "DELETE FROM t WHERE id = 1", want: "EXPLAIN FORMATTED DELETE FROM t WHERE id = 1"},
		{name: "a plan is not planned twice", statement: "EXPLAIN SELECT 1", want: "EXPLAIN SELECT 1"},
		{name: "a typed plan keeps its mode", statement: "EXPLAIN EXTENDED SELECT 1", want: "EXPLAIN EXTENDED SELECT 1"},
		{name: "a comment before a plan does not hide it", statement: "-- plan\nEXPLAIN SELECT 1", want: "-- plan\nEXPLAIN SELECT 1"},
		{name: "a plan that executes keeps its mode", statement: "EXPLAIN ANALYZE SELECT 1", want: "EXPLAIN ANALYZE SELECT 1"},
		{name: "nested comment hiding EXPLAIN before DELETE", statement: "/* outer /* inner */\nEXPLAIN\n*/\nDELETE FROM t", want: "EXPLAIN FORMATTED /* outer /* inner */\nEXPLAIN\n*/\nDELETE FROM t"},
		{name: "nested comment hiding EXPLAIN before SELECT", statement: "/* a /* b */ EXPLAIN */ SELECT 1", want: "EXPLAIN FORMATTED /* a /* b */ EXPLAIN */ SELECT 1"},
		{name: "nested comment closed before EXPLAIN refuses to guess", statement: "/* a /* b */ */ EXPLAIN SELECT 1", want: "EXPLAIN FORMATTED /* a /* b */ */ EXPLAIN SELECT 1"},
		{name: "unclosed nested comment", statement: "/* a /* b */ EXPLAIN SELECT 1", want: "EXPLAIN FORMATTED /* a /* b */ EXPLAIN SELECT 1"},
		{name: "line comment before non-nested block comment", statement: "-- /* a\n/* b */ EXPLAIN SELECT 1", want: "-- /* a\n/* b */ EXPLAIN SELECT 1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, format := range []string{"", "TEXT", "JSON", "XML"} {
				if got := ExplainStatement(tc.statement, format); got != tc.want {
					t.Errorf("ExplainStatement(%q, %q) = %q, want %q", tc.statement, format, got, tc.want)
				}
			}
		})
	}
}

func TestPlanExecuted(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      bool
	}{
		{statement: "EXPLAIN SELECT 1", want: false},
		{statement: "EXPLAIN ANALYZE SELECT 1", want: true},
		{statement: "-- plan\nEXPLAIN ANALYZE SELECT 1", want: true},
		{statement: "EXPLAIN /* c */ ANALYZE SELECT 1", want: true},
		{statement: "SELECT 1", want: false},
		{statement: "/* a /* b */ EXPLAIN ANALYZE */ DELETE FROM t", want: false},
	} {
		if got := PlanExecuted(tc.statement); got != tc.want {
			t.Errorf("PlanExecuted(%q) = %v, want %v", tc.statement, got, tc.want)
		}
	}
}

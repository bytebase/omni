package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/hive/parser"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		empty   []bool
		offsets [][2]int
	}{
		{
			name:    "one statement",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 9}},
		},
		{
			name:    "two statements",
			script:  "SELECT 1;\nSELECT 2;",
			want:    []string{"SELECT 1;", "\nSELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 19}},
		},
		{
			name:    "two statements on one line",
			script:  "SELECT 1; SELECT 2;",
			want:    []string{"SELECT 1;", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 19}},
		},
		{
			name:   "no terminator",
			script: "SELECT 1",
			want:   []string{"SELECT 1"},
			empty:  []bool{false},
		},
		{
			name:   "a comment line belongs to the statement after it",
			script: "-- a comment\nSELECT 1;",
			want:   []string{"-- a comment\nSELECT 1;"},
			empty:  []bool{false},
		},
		{
			name:    "a trailing comment after the last statement is marked empty",
			script:  "SELECT 1;\n-- trailing\n",
			want:    []string{"SELECT 1;", "\n-- trailing\n"},
			empty:   []bool{false, true},
			offsets: [][2]int{{0, 9}, {9, 22}},
		},
		{
			name:   "empty string",
			script: "",
			want:   []string{},
		},
		{
			name:   "whitespace only script is marked empty",
			script: "\n\n",
			want:   []string{"\n\n"},
			empty:  []bool{true},
		},
		{
			name:   "comments only script is marked empty",
			script: "-- line\n/* block /* nested */ */\n",
			want:   []string{"-- line\n/* block /* nested */ */\n"},
			empty:  []bool{true},
		},
		{
			name:   "a block comment that ends its line",
			script: "/* a comment */\nSELECT 1;",
			want:   []string{"/* a comment */\nSELECT 1;"},
			empty:  []bool{false},
		},
		{
			name:   "an inline block comment before a statement",
			script: "/* a comment */ SELECT 1;",
			want:   []string{"/* a comment */ SELECT 1;"},
			empty:  []bool{false},
		},
		{
			name:   "nested block comments with semicolon inside",
			script: "/* outer /* inner ; */ outer */ SELECT 1;",
			want:   []string{"/* outer /* inner ; */ outer */ SELECT 1;"},
			empty:  []bool{false},
		},
		{
			name:   "semicolon inside single and double quoted strings and escapes",
			script: `SELECT 'a;b', 'it\'s;ok', "c;d"; SELECT 'e;\nf';`,
			want:   []string{`SELECT 'a;b', 'it\'s;ok', "c;d";`, ` SELECT 'e;\nf';`},
			empty:  []bool{false, false},
		},
		{
			name:   "semicolon inside multiline string literal",
			script: "SELECT 'line1;\nline2';\nSELECT 2;",
			want:   []string{"SELECT 'line1;\nline2';", "\nSELECT 2;"},
			empty:  []bool{false, false},
		},
		{
			name:   "semicolon inside backtick identifier with escaped backtick",
			script: "SELECT `a;b`, `c``;d` FROM t; SELECT 2;",
			want:   []string{"SELECT `a;b`, `c``;d` FROM t;", " SELECT 2;"},
			empty:  []bool{false, false},
		},
		{
			name:   "query hint inside statement",
			script: "SELECT /*+ MAPJOIN(a) */ * FROM a; SELECT 2;",
			want:   []string{"SELECT /*+ MAPJOIN(a) */ * FROM a;", " SELECT 2;"},
			empty:  []bool{false, false},
		},
		{
			name:   "inline comment after semicolon",
			script: "SELECT 1; -- comment\nSELECT 2;",
			want:   []string{"SELECT 1;", " -- comment\nSELECT 2;"},
			empty:  []bool{false, false},
		},
		{
			name:    "empty semicolons are preserved and marked empty",
			script:  ";SELECT 1;;SELECT 2;;",
			want:    []string{";", "SELECT 1;", ";", "SELECT 2;", ";"},
			empty:   []bool{true, false, true, false, true},
			offsets: [][2]int{{0, 1}, {1, 10}, {10, 11}, {11, 20}, {20, 21}},
		},
		{
			name:   "unclosed block comment stays non-empty in its segment",
			script: "SELECT 1; /* unclosed",
			want:   []string{"SELECT 1;", " /* unclosed"},
			empty:  []bool{false, false},
		},
		{
			name:   "unclosed string literal consumes to EOF",
			script: "SELECT 'unclosed; SELECT 2;",
			want:   []string{"SELECT 'unclosed; SELECT 2;"},
			empty:  []bool{false},
		},
		{
			name:    "CRLF multiple statements",
			script:  "SELECT 1;\r\nSELECT 2;",
			want:    []string{"SELECT 1;", "\r\nSELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 20}},
		},
		{
			name:    "multi-byte UTF-8",
			script:  "SELECT '中';\nSELECT '文';",
			want:    []string{"SELECT '中';", "\nSELECT '文';"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 13}, {13, 27}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			statements, err := Split(tc.script)
			if err != nil {
				t.Fatalf("Split: %v", err)
			}
			texts := make([]string, 0, len(statements))
			for _, stmt := range statements {
				texts = append(texts, stmt.Text)
				if stmt.Text != tc.script[stmt.Start:stmt.End] {
					t.Errorf("statement text %q does not match its offsets [%d,%d)", stmt.Text, stmt.Start, stmt.End)
				}
			}
			if len(texts) != len(tc.want) {
				t.Fatalf("got %d statements %q, want %d", len(texts), texts, len(tc.want))
			}
			for i := range texts {
				if texts[i] != tc.want[i] {
					t.Errorf("statement %d = %q, want %q", i, texts[i], tc.want[i])
				}
				if i < len(tc.empty) && statements[i].Empty != tc.empty[i] {
					t.Errorf("statement %d Empty = %v, want %v", i, statements[i].Empty, tc.empty[i])
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

func TestExplainStatement(t *testing.T) {
	cases := []struct {
		statement string
		format    string
		want      string
	}{
		{"SELECT 1", "", "EXPLAIN SELECT 1"},
		{"SELECT 1;", "", "EXPLAIN SELECT 1;"},
		{"SELECT 1", "TEXT", "EXPLAIN SELECT 1"},
		{"SELECT 1", "XML", "EXPLAIN SELECT 1"},
		{"SELECT 1", "JSON", "EXPLAIN FORMATTED SELECT 1"},
		{"select 1", "json", "EXPLAIN FORMATTED select 1"},
		{"EXPLAIN SELECT 1", "JSON", "EXPLAIN SELECT 1"},
		{"EXPLAIN EXTENDED SELECT 1", "JSON", "EXPLAIN EXTENDED SELECT 1"},
		{"EXPLAIN DEPENDENCY SELECT 1", "JSON", "EXPLAIN DEPENDENCY SELECT 1"},
		{"EXPLAIN AUTHORIZATION SELECT 1", "JSON", "EXPLAIN AUTHORIZATION SELECT 1"},
		{"explain formatted select 1", "", "explain formatted select 1"},
		{"/* a comment */\nEXPLAIN SELECT 1", "JSON", "/* a comment */\nEXPLAIN SELECT 1"},
		{"/* outer /* inner */ outer */\nEXPLAIN SELECT 1", "JSON", "/* outer /* inner */ outer */\nEXPLAIN SELECT 1"},
		{"-- plan it\nEXPLAIN SELECT 1", "", "-- plan it\nEXPLAIN SELECT 1"},
		{"EXPLAINED", "", "EXPLAIN EXPLAINED"},
		{"explained_view_query()", "", "EXPLAIN explained_view_query()"},
		{"`EXPLAIN` SELECT 1", "JSON", "EXPLAIN FORMATTED `EXPLAIN` SELECT 1"},
		{"/* /* */ EXPLAIN */ SELECT 1", "JSON", "EXPLAIN FORMATTED /* /* */ EXPLAIN */ SELECT 1"},
		{"/* unclosed EXPLAIN SELECT 1", "JSON", "EXPLAIN FORMATTED /* unclosed EXPLAIN SELECT 1"},
	}
	for _, tc := range cases {
		if got := ExplainStatement(tc.statement, tc.format); got != tc.want {
			t.Errorf("ExplainStatement(%q, %q) = %q, want %q", tc.statement, tc.format, got, tc.want)
		}
	}
}

func TestPlanFormatAndExecuted(t *testing.T) {
	cases := []struct {
		statement string
		format    string
		ok        bool
		executed  bool
	}{
		{"EXPLAIN SELECT 1", "TEXT", true, false},
		{"EXPLAIN SELECT 1;", "TEXT", true, false},
		{"/* c */\nEXPLAIN FORMATTED SELECT 1", "JSON", true, false},
		{"/* outer /* inner */ outer */\nEXPLAIN FORMATTED SELECT 1", "JSON", true, false},
		{"EXPLAIN /* c */ FORMATTED SELECT 1", "JSON", true, false},
		{"EXPLAIN EXTENDED SELECT 1", "TEXT", true, false},
		{"EXPLAIN DEPENDENCY SELECT 1", "TEXT", true, false},
		{"EXPLAIN AUTHORIZATION SELECT 1", "TEXT", true, false},
		{"EXPLAIN CBO COST FORMATTED SELECT 1", "JSON", true, false},
		{"EXPLAIN VECTORIZATION ONLY DETAIL FORMATTED SELECT 1", "JSON", true, false},
		{"EXPLAIN REWRITE SELECT 1", "TEXT", true, false},
		{"EXPLAIN ANALYZE SELECT 1", "TEXT", true, true},
		{"EXPLAIN ANALYZE FORMATTED SELECT 1", "JSON", true, true},
		{"EXPLAIN FORMATTED ANALYZE SELECT 1", "JSON", true, true},
		{"EXPLAIN ANALYZE TABLE t COMPUTE STATISTICS", "TEXT", true, false},
		{"EXPLAIN DESCRIBE FORMATTED t", "TEXT", true, false},
		{"EXPLAINED", "", false, false},
		{"`EXPLAIN` SELECT 1", "", false, false},
		{"/*+ MAPJOIN(t) */ EXPLAIN SELECT 1", "", false, false},
		{"/* unclosed EXPLAIN SELECT 1", "", false, false},
		{"SELECT 1", "", false, false},
		{"-- only a comment", "", false, false},
	}
	for _, tc := range cases {
		format, ok := PlanFormat(tc.statement)
		if ok != tc.ok || format != tc.format {
			t.Errorf("PlanFormat(%q) = (%q, %v), want (%q, %v)", tc.statement, format, ok, tc.format, tc.ok)
		}
		if got := StartsWithExplain(tc.statement); got != tc.ok {
			t.Errorf("StartsWithExplain(%q) = %v, want %v", tc.statement, got, tc.ok)
		}
		if got := PlanExecuted(tc.statement); got != tc.executed {
			t.Errorf("PlanExecuted(%q) = %v, want %v", tc.statement, got, tc.executed)
		}
	}
}

func TestValidateQuery(t *testing.T) {
	cases := []struct {
		name        string
		script      string
		readOnly    bool
		returnsData bool
		wantErr     bool
		wantErrPos  int
	}{
		{name: "empty", script: "", readOnly: true, returnsData: true},
		{name: "comments only", script: "-- a\n/* b */", readOnly: true, returnsData: true},
		{name: "simple select", script: "SELECT 1", readOnly: true, returnsData: true},
		{name: "multi select", script: "SELECT 1; SELECT 2;", readOnly: true, returnsData: true},
		{name: "cte select", script: "WITH c AS (SELECT 1 AS x) SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "from-first select", script: "FROM t SELECT a, b", readOnly: true, returnsData: true},
		{name: "values", script: "VALUES (1, 2)", readOnly: true, returnsData: true},
		{name: "show tables", script: "SHOW TABLES IN default", readOnly: true, returnsData: true},
		{name: "describe formatted", script: "DESCRIBE FORMATTED t", readOnly: true, returnsData: true},
		{name: "desc table", script: "DESC t", readOnly: true, returnsData: true},
		{name: "explain select", script: "EXPLAIN SELECT 1", readOnly: true, returnsData: true},
		{name: "explain formatted", script: "EXPLAIN FORMATTED SELECT 1", readOnly: true, returnsData: true},
		{name: "explain extended", script: "EXPLAIN EXTENDED SELECT 1", readOnly: true, returnsData: true},
		{name: "explain dependency", script: "EXPLAIN DEPENDENCY SELECT 1", readOnly: true, returnsData: true},
		{name: "explain authorization", script: "EXPLAIN AUTHORIZATION SELECT 1", readOnly: true, returnsData: true},
		{name: "explain rewrite", script: "EXPLAIN REWRITE SELECT 1", readOnly: true, returnsData: true},
		{name: "explain analyze table is non-executing plan", script: "EXPLAIN ANALYZE TABLE t COMPUTE STATISTICS", readOnly: true, returnsData: true},
		{name: "explain analyze executes", script: "EXPLAIN ANALYZE SELECT 1", readOnly: false, returnsData: false},
		{name: "explain formatted analyze executes", script: "EXPLAIN FORMATTED ANALYZE SELECT 1", readOnly: false, returnsData: false},
		{name: "insert", script: "INSERT INTO t VALUES (1)", readOnly: false, returnsData: false},
		{name: "cte insert rejected", script: "WITH c AS (SELECT 1 AS x) INSERT INTO t SELECT * FROM c", readOnly: false, returnsData: false},
		{name: "from-first insert rejected", script: "FROM src INSERT OVERWRITE TABLE dst SELECT a", readOnly: false, returnsData: false},
		{name: "from-first insert directory rejected", script: "FROM src INSERT OVERWRITE LOCAL DIRECTORY '/tmp/out' SELECT a", readOnly: false, returnsData: false},
		{name: "select transform rejected", script: "SELECT TRANSFORM(a) USING 'cat' FROM t", readOnly: false, returnsData: false},
		{name: "map rejected", script: "MAP a USING 'cat' FROM t", readOnly: false, returnsData: false},
		{name: "reduce rejected", script: "REDUCE a USING 'cat' FROM t", readOnly: false, returnsData: false},
		{name: "update rejected", script: "UPDATE t SET a = 1", readOnly: false, returnsData: false},
		{name: "delete rejected", script: "DELETE FROM t WHERE a = 1", readOnly: false, returnsData: false},
		{name: "ddl create rejected", script: "CREATE TABLE t (id INT)", readOnly: false, returnsData: false},
		{name: "ddl drop rejected", script: "DROP TABLE t", readOnly: false, returnsData: false},
		{name: "use db rejected", script: "USE other_db", readOnly: false, returnsData: false},
		{name: "select followed by drop rejected", script: "SELECT 1; DROP TABLE t;", readOnly: false, returnsData: false},
		{name: "syntax error", script: "SELECT FROM", wantErr: true, wantErrPos: 7},
		{name: "syntax error in second statement has script offset", script: "SELECT 1; SELECT FROM", wantErr: true, wantErrPos: 17},
		{name: "unclosed block comment fails closed", script: "/* unclosed", wantErr: true, wantErrPos: 0},
		{name: "unclosed string fails closed", script: "SELECT 'unclosed", wantErr: true, wantErrPos: 7},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ro, rd, err := ValidateQuery(tc.script)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ValidateQuery(%q) = (%v, %v, nil), want error", tc.script, ro, rd)
				}
				if ro || rd {
					t.Errorf("ValidateQuery(%q) returned (%v, %v) on error", tc.script, ro, rd)
				}
				var pe *parser.ParseError
				if !errors.As(err, &pe) {
					t.Fatalf("expected *parser.ParseError, got %T: %v", err, err)
				}
				if pe.Position != tc.wantErrPos {
					t.Errorf("ParseError.Position = %d, want %d (err: %v)", pe.Position, tc.wantErrPos, pe)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateQuery(%q) unexpected error: %v", tc.script, err)
			}
			if ro != tc.readOnly || rd != tc.returnsData {
				t.Errorf("ValidateQuery(%q) = (%v, %v), want (%v, %v)", tc.script, ro, rd, tc.readOnly, tc.returnsData)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		statement string
		want      QueryType
		wantErr   bool
	}{
		{"", analysis.QueryTypeUnknown, false},
		{"-- comment only", analysis.QueryTypeUnknown, false},
		{"SELECT 1;", analysis.Select, false},
		{"WITH c AS (SELECT 1) SELECT * FROM c", analysis.Select, false},
		{"FROM t SELECT a", analysis.Select, false},
		{"EXPLAIN FORMATTED SELECT 1", analysis.Explain, false},
		{"SHOW TABLES", analysis.SelectInfoSchema, false},
		{"DESCRIBE FORMATTED t", analysis.SelectInfoSchema, false},
		{"DESC t", analysis.SelectInfoSchema, false},
		{"CREATE TABLE t (id INT)", analysis.DDL, false},
		{"ALTER TABLE t ADD COLUMNS (b STRING)", analysis.DDL, false},
		{"DROP TABLE t", analysis.DDL, false},
		{"INSERT INTO t VALUES (1)", analysis.DML, false},
		{"WITH c AS (SELECT 1) INSERT INTO t SELECT * FROM c", analysis.DML, false},
		{"FROM src INSERT OVERWRITE TABLE dst SELECT a", analysis.DML, false},
		{"UPDATE t SET a = 1", analysis.DML, false},
		{"DELETE FROM t WHERE a = 1", analysis.DML, false},
		{"LOAD DATA LOCAL INPATH '/tmp/f' INTO TABLE t", analysis.DML, false},
		{"COMMIT", analysis.QueryTypeUnknown, false},
		{"SELECT 1; SELECT 2;", analysis.QueryTypeUnknown, true},
		{"SELECT FROM", analysis.QueryTypeUnknown, true},
		{"/* unclosed", analysis.QueryTypeUnknown, true},
	}
	for _, tc := range cases {
		got, err := Classify(tc.statement)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Classify(%q) = (%v, nil), want error", tc.statement, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Classify(%q) unexpected error: %v", tc.statement, err)
			continue
		}
		if got != tc.want {
			t.Errorf("Classify(%q) = %v, want %v", tc.statement, got, tc.want)
		}
	}
}

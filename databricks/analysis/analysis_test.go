package analysis

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/databricks/parser"
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
			name:   "a comment after the last statement is preserved as empty",
			script: "SELECT 1;\n-- trailing\n",
			want:   []string{"SELECT 1;", "\n-- trailing\n"},
		},
		{
			name:   "an empty script is preserved as empty",
			script: "\n\n",
			want:   []string{"\n\n"},
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
			want:   []string{"/* a comment */ SELECT 1;"},
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
		{name: "nested comment closed before EXPLAIN is recognized", statement: "/* a /* b */ */ EXPLAIN SELECT 1", want: "/* a /* b */ */ EXPLAIN SELECT 1"},
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
		{statement: "EXPLAIN ANALYZE TABLE t COMPUTE STATISTICS", want: false},
		{statement: "EXPLAIN ANALYZE TABLES COMPUTE STATISTICS", want: false},
		{statement: "/* a /* b */ */ EXPLAIN ANALYZE SELECT 1", want: true},
	} {
		if got := PlanExecuted(tc.statement); got != tc.want {
			t.Errorf("PlanExecuted(%q) = %v, want %v", tc.statement, got, tc.want)
		}
	}
}

func TestClassifyAndStatementReadsRows(t *testing.T) {
	cases := []struct {
		name      string
		statement string
		wantType  QueryType
		wantRows  bool
		wantErr   bool
	}{
		{name: "select", statement: "SELECT 1", wantType: analysis.Select, wantRows: true},
		{name: "cte select", statement: "WITH c AS (SELECT 1) SELECT * FROM c", wantType: analysis.Select, wantRows: true},
		{name: "values", statement: "VALUES (1, 'a')", wantType: analysis.Select, wantRows: true},
		{name: "table", statement: "TABLE main.app.users", wantType: analysis.Select, wantRows: true},
		{name: "from pipe", statement: "FROM t |> WHERE x > 1 |> SELECT *", wantType: analysis.Select, wantRows: true},
		{name: "show", statement: "SHOW TABLES", wantType: analysis.SelectInfoSchema, wantRows: true},
		{name: "describe", statement: "DESCRIBE EXTENDED t", wantType: analysis.SelectInfoSchema, wantRows: true},
		{name: "explain", statement: "EXPLAIN FORMATTED SELECT 1", wantType: analysis.Explain, wantRows: true},
		{name: "explain analyze", statement: "EXPLAIN ANALYZE SELECT 1", wantType: analysis.Explain, wantRows: true},
		{name: "insert", statement: "INSERT INTO t VALUES (1)", wantType: analysis.DML, wantRows: false},
		{name: "cte insert", statement: "WITH c AS (SELECT 1) INSERT INTO t SELECT * FROM c", wantType: analysis.DML, wantRows: false},
		{name: "from multi-insert", statement: "FROM s INSERT INTO t SELECT *", wantType: analysis.DML, wantRows: false},
		{name: "update", statement: "UPDATE t SET a = 1", wantType: analysis.DML, wantRows: false},
		{name: "delete", statement: "DELETE FROM t WHERE id = 1", wantType: analysis.DML, wantRows: false},
		{name: "merge", statement: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE", wantType: analysis.DML, wantRows: false},
		{name: "create table", statement: "CREATE TABLE t (id INT)", wantType: analysis.DDL, wantRows: false},
		{name: "create view", statement: "CREATE VIEW v AS SELECT 1", wantType: analysis.DDL, wantRows: false},
		{name: "create materialized view", statement: "CREATE MATERIALIZED VIEW mv AS SELECT 1", wantType: analysis.DDL, wantRows: false},
		{name: "alter table", statement: "ALTER TABLE t ADD COLUMNS (b STRING)", wantType: analysis.DDL, wantRows: false},
		{name: "drop table", statement: "DROP TABLE t", wantType: analysis.DDL, wantRows: false},
		{name: "truncate", statement: "TRUNCATE TABLE t", wantType: analysis.DDL, wantRows: false},
		{name: "set", statement: "SET spark.sql.ansi.enabled = true", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "use catalog", statement: "USE CATALOG main", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "use schema", statement: "USE SCHEMA app", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "set catalog", statement: "SET CATALOG main", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "call", statement: "CALL app.my_proc(1)", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "grant", statement: "GRANT SELECT ON TABLE t TO `user@example.com`", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "compound block", statement: "BEGIN SELECT 1; SELECT 2; END;", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "empty", statement: "   -- comment\n; ", wantType: analysis.QueryTypeUnknown, wantRows: false},
		{name: "syntax error", statement: "SELECT 1 +", wantErr: true},
		{name: "unclosed comment", statement: "/* unclosed SELECT 1", wantErr: true},
		{name: "unclosed string", statement: "SELECT 'abc", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotType, err := Classify(tc.statement)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Classify(%q) = %v, nil; want error", tc.statement, gotType)
				}
				var pe *parser.ParseError
				if !errors.As(err, &pe) {
					t.Errorf("Classify(%q) err = %T (%v), want *parser.ParseError", tc.statement, err, err)
				}
				if _, rerr := StatementReadsRows(tc.statement); rerr == nil {
					t.Errorf("StatementReadsRows(%q) returned nil error, want error", tc.statement)
				}
				if ClassifyStatement(tc.statement) {
					t.Errorf("ClassifyStatement(%q) = true, want false on error", tc.statement)
				}
				return
			}
			if err != nil {
				t.Fatalf("Classify(%q) unexpected error: %v", tc.statement, err)
			}
			if gotType != tc.wantType {
				t.Errorf("Classify(%q) = %v, want %v", tc.statement, gotType, tc.wantType)
			}
			gotRows, rerr := StatementReadsRows(tc.statement)
			if rerr != nil {
				t.Fatalf("StatementReadsRows(%q) unexpected error: %v", tc.statement, rerr)
			}
			if gotRows != tc.wantRows {
				t.Errorf("StatementReadsRows(%q) = %v, want %v", tc.statement, gotRows, tc.wantRows)
			}
			if got := ClassifyStatement(tc.statement); got != tc.wantRows {
				t.Errorf("ClassifyStatement(%q) = %v, want %v", tc.statement, got, tc.wantRows)
			}
		})
	}

	t.Run("multi-statement Classify returns error", func(t *testing.T) {
		if _, err := Classify("SELECT 1; SELECT 2;"); err == nil {
			t.Fatal("Classify accepted multi-statement input")
		}
		if reads, err := StatementReadsRows("SELECT 1; SELECT 2;"); err != nil || reads {
			t.Errorf("StatementReadsRows(multi) = (%v, %v), want (false, nil)", reads, err)
		}
	})
}

func TestSplitCompoundBlock(t *testing.T) {
	script := "BEGIN\n  SELECT 1;\n  SELECT 2;\nEND;\nSELECT 3;"
	stmts, err := Split(script)
	if err != nil {
		t.Fatalf("Split: %v", err)
	}
	if len(stmts) != 2 {
		t.Fatalf("got %d statements, want 2: %+v", len(stmts), stmts)
	}
	if stmts[0].Empty || stmts[1].Empty {
		t.Errorf("unexpected empty flag: %+v", stmts)
	}
	if stmts[0].Text != "BEGIN\n  SELECT 1;\n  SELECT 2;\nEND;" {
		t.Errorf("stmt[0] = %q", stmts[0].Text)
	}
	if stmts[1].Text != "\nSELECT 3;" {
		t.Errorf("stmt[1] = %q", stmts[1].Text)
	}
}

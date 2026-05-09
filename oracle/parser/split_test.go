package parser

import "testing"

func TestSplitOrdinarySQL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "empty input",
			sql:  "",
			want: nil,
		},
		{
			name: "two semicolon statements",
			sql:  "SELECT 1 FROM dual; SELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", " SELECT 2 FROM dual"},
		},
		{
			name: "trailing statement without semicolon",
			sql:  "SELECT 1 FROM dual; SELECT 2 FROM dual",
			want: []string{"SELECT 1 FROM dual", " SELECT 2 FROM dual"},
		},
		{
			name: "filters empty statements",
			sql:  ";\n -- comment only\n ; SELECT 1 FROM dual;",
			want: []string{" SELECT 1 FROM dual"},
		},
		{
			name: "semicolon inside single quoted string",
			sql:  "SELECT 'a;b' FROM dual; SELECT 1 FROM dual;",
			want: []string{"SELECT 'a;b' FROM dual", " SELECT 1 FROM dual"},
		},
		{
			name: "semicolon inside q quote",
			sql:  "SELECT q'[a;b]' FROM dual; SELECT 1 FROM dual;",
			want: []string{"SELECT q'[a;b]' FROM dual", " SELECT 1 FROM dual"},
		},
		{
			name: "semicolon inside double quoted identifier",
			sql:  `SELECT "a;b" FROM dual; SELECT 1 FROM dual;`,
			want: []string{`SELECT "a;b" FROM dual`, " SELECT 1 FROM dual"},
		},
		{
			name: "semicolon inside comments",
			sql:  "SELECT 1 /* ; */ FROM dual; -- ;\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 /* ; */ FROM dual", " -- ;\nSELECT 2 FROM dual"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTexts(Split(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitSlashDelimiter(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "slash on own line separates SQL",
			sql:  "SELECT 1 FROM dual\n/\nSELECT 2 FROM dual\n/",
			want: []string{"SELECT 1 FROM dual", "\nSELECT 2 FROM dual"},
		},
		{
			name: "slash line allows surrounding whitespace",
			sql:  "SELECT 1 FROM dual\n  /  \nSELECT 2 FROM dual",
			want: []string{"SELECT 1 FROM dual", "\nSELECT 2 FROM dual"},
		},
		{
			name: "slash in division expression is not delimiter",
			sql:  "SELECT 10 / 2 FROM dual; SELECT 3 FROM dual;",
			want: []string{"SELECT 10 / 2 FROM dual", " SELECT 3 FROM dual"},
		},
		{
			name: "slash in string and comment is not delimiter",
			sql:  "SELECT '/' FROM dual; /*\n/\n*/ SELECT 2 FROM dual;",
			want: []string{"SELECT '/' FROM dual", " /*\n/\n*/ SELECT 2 FROM dual"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTexts(Split(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitRanges(t *testing.T) {
	sql := "\nSELECT 1 FROM dual;\n  /  \nSELECT 2 FROM dual;"
	got := Split(sql)
	if len(got) != 2 {
		t.Fatalf("got %d segments, want 2: %#v", len(got), got)
	}
	for _, seg := range got {
		if seg.Text != sql[seg.ByteStart:seg.ByteEnd] {
			t.Fatalf("segment range [%d,%d] extracts %q, want Text %q", seg.ByteStart, seg.ByteEnd, sql[seg.ByteStart:seg.ByteEnd], seg.Text)
		}
	}
	if got[0].Text != "\nSELECT 1 FROM dual" {
		t.Fatalf("first Text = %q", got[0].Text)
	}
	if got[1].Text != "\nSELECT 2 FROM dual" {
		t.Fatalf("second Text = %q", got[1].Text)
	}
}

func TestSplitPLSQLBlocks(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "anonymous begin end block",
			sql:  "BEGIN NULL; END;\n/\nSELECT 1 FROM dual;",
			want: []string{"BEGIN NULL; END;", "\nSELECT 1 FROM dual"},
		},
		{
			name: "declare block with nested if",
			sql: "DECLARE\n" +
				"  v NUMBER := 1;\n" +
				"BEGIN\n" +
				"  IF v > 0 THEN\n" +
				"    v := v + 1;\n" +
				"  END IF;\n" +
				"END;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"DECLARE\n  v NUMBER := 1;\nBEGIN\n  IF v > 0 THEN\n    v := v + 1;\n  END IF;\nEND;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "create procedure with internal statements",
			sql: "CREATE OR REPLACE PROCEDURE p IS\n" +
				"BEGIN\n" +
				"  NULL;\n" +
				"  NULL;\n" +
				"END;\n" +
				"/\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE OR REPLACE PROCEDURE p IS\nBEGIN\n  NULL;\n  NULL;\nEND;",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "create editionable procedure",
			sql: "CREATE OR REPLACE EDITIONABLE PROCEDURE p IS\n" +
				"BEGIN\n" +
				"  NULL;\n" +
				"END;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE OR REPLACE EDITIONABLE PROCEDURE p IS\nBEGIN\n  NULL;\nEND;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "create type body",
			sql: "CREATE TYPE BODY typ AS\n" +
				"  MEMBER FUNCTION f RETURN NUMBER IS\n" +
				"  BEGIN\n" +
				"    RETURN 1;\n" +
				"  END;\n" +
				"END;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE TYPE BODY typ AS\n  MEMBER FUNCTION f RETURN NUMBER IS\n  BEGIN\n    RETURN 1;\n  END;\nEND;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "two package units with slash separators",
			sql: "CREATE PACKAGE pkg IS\n" +
				"  PROCEDURE p;\n" +
				"END pkg;\n" +
				"/\n" +
				"CREATE PACKAGE BODY pkg IS\n" +
				"  PROCEDURE p IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END;\n" +
				"END pkg;\n" +
				"/",
			want: []string{
				"CREATE PACKAGE pkg IS\n  PROCEDURE p;\nEND pkg;",
				"\nCREATE PACKAGE BODY pkg IS\n  PROCEDURE p IS\n  BEGIN\n    NULL;\n  END;\nEND pkg;",
			},
		},
		{
			name: "create function with division expression",
			sql: "CREATE FUNCTION f RETURN NUMBER IS\n" +
				"BEGIN\n" +
				"  RETURN 10 / 2;\n" +
				"END;\n" +
				"/\n",
			want: []string{
				"CREATE FUNCTION f RETURN NUMBER IS\nBEGIN\n  RETURN 10 / 2;\nEND;",
			},
		},
		{
			name: "compound trigger",
			sql: "CREATE TRIGGER trg\n" +
				"FOR INSERT ON t\n" +
				"COMPOUND TRIGGER\n" +
				"  BEFORE EACH ROW IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END BEFORE EACH ROW;\n" +
				"END trg;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE TRIGGER trg\nFOR INSERT ON t\nCOMPOUND TRIGGER\n  BEFORE EACH ROW IS\n  BEGIN\n    NULL;\n  END BEFORE EACH ROW;\nEND trg;",
				"\nSELECT 1 FROM dual",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTexts(Split(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitSQLPlusCommands(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "environment commands are skipped",
			sql: "SET DEFINE OFF\n" +
				"SET SERVEROUTPUT ON;\n" +
				"PROMPT creating table;\n" +
				"SPOOL install.log\n" +
				"SELECT 1 FROM dual;\n" +
				"SPOOL OFF\n" +
				"SELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", "SELECT 2 FROM dual"},
		},
		{
			name: "script and session commands are skipped",
			sql: "CONNECT scott/tiger@db\n" +
				"@preflight.sql\n" +
				"@@nested/install.sql arg1 arg2\n" +
				"START post.sql\n" +
				"WHENEVER SQLERROR EXIT SQL.SQLCODE ROLLBACK\n" +
				"SELECT 1 FROM dual;\n" +
				"EXIT SUCCESS",
			want: []string{"SELECT 1 FROM dual"},
		},
		{
			name: "remark and host commands are skipped",
			sql: "REM this ; is a SQL*Plus comment\n" +
				"REMARK this is also ignored\n" +
				"HOST echo before\n" +
				"! echo shell\n" +
				"SELECT 1 FROM dual;",
			want: []string{"SELECT 1 FROM dual"},
		},
		{
			name: "formatting and variable commands are skipped",
			sql: "COLUMN c FORMAT A20\n" +
				"BREAK ON report\n" +
				"COMPUTE SUM OF sal ON report\n" +
				"TTITLE left 'Report'\n" +
				"BTITLE off\n" +
				"DEFINE schema_name = HR\n" +
				"UNDEFINE schema_name\n" +
				"ACCEPT v PROMPT 'Value: '\n" +
				"VARIABLE rc NUMBER\n" +
				"PRINT rc\n" +
				"SELECT 1 FROM dual;",
			want: []string{"SELECT 1 FROM dual"},
		},
		{
			name: "run flushes current SQL buffer",
			sql:  "SELECT 1 FROM dual\nRUN\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", "\nSELECT 2 FROM dual"},
		},
		{
			name: "sqlplus command words inside SQL are not skipped",
			sql:  "SELECT 'SET DEFINE OFF' AS txt FROM dual; SELECT prompt FROM t;",
			want: []string{"SELECT 'SET DEFINE OFF' AS txt FROM dual", " SELECT prompt FROM t"},
		},
		{
			name: "sqlplus command words in plsql are not skipped",
			sql: "BEGIN\n" +
				"  prompt := 'not a command';\n" +
				"  NULL;\n" +
				"END;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{"BEGIN\n  prompt := 'not a command';\n  NULL;\nEND;", "\nSELECT 1 FROM dual"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitTexts(Split(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func splitTexts(segs []Segment) []string {
	if len(segs) == 0 {
		return nil
	}
	out := make([]string, 0, len(segs))
	for _, seg := range segs {
		out = append(out, seg.Text)
	}
	return out
}

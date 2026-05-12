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
		{
			name: "create sequence multiline start with is not SQLPlus",
			sql: "CREATE SEQUENCE audit_log_seq\n" +
				"  START WITH 1000\n" +
				"  INCREMENT BY 1\n" +
				"  CACHE 20;",
			want: []string{"CREATE SEQUENCE audit_log_seq\n  START WITH 1000\n  INCREMENT BY 1\n  CACHE 20"},
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
			name: "create function with declarations without slash separator",
			sql: "CREATE FUNCTION calc_bonus(p_start_date DATE)\n" +
				"RETURN DATE\n" +
				"IS\n" +
				"  v_current_date DATE := p_start_date;\n" +
				"BEGIN\n" +
				"  RETURN v_current_date;\n" +
				"END calc_bonus;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE FUNCTION calc_bonus(p_start_date DATE)\nRETURN DATE\nIS\n  v_current_date DATE := p_start_date;\nBEGIN\n  RETURN v_current_date;\nEND calc_bonus;",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "create procedure with declarations without slash separator",
			sql: "CREATE PROCEDURE update_salary(p_employee_id NUMBER)\n" +
				"IS\n" +
				"  v_delta NUMBER := 1;\n" +
				"BEGIN\n" +
				"  UPDATE employees SET salary = salary + v_delta WHERE id = p_employee_id;\n" +
				"END update_salary;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE PROCEDURE update_salary(p_employee_id NUMBER)\nIS\n  v_delta NUMBER := 1;\nBEGIN\n  UPDATE employees SET salary = salary + v_delta WHERE id = p_employee_id;\nEND update_salary;",
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
			name: "two package units without slash separators",
			sql: "CREATE PACKAGE pkg IS\n" +
				"  PROCEDURE p;\n" +
				"END pkg;\n" +
				"CREATE PACKAGE BODY pkg IS\n" +
				"  PROCEDURE p IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END;\n" +
				"END pkg;",
			want: []string{
				"CREATE PACKAGE pkg IS\n  PROCEDURE p;\nEND pkg;",
				"\nCREATE PACKAGE BODY pkg IS\n  PROCEDURE p IS\n  BEGIN\n    NULL;\n  END;\nEND pkg;",
			},
		},
		{
			name: "package body initialization without slash separator",
			sql: "CREATE PACKAGE BODY pkg IS\n" +
				"  PROCEDURE p IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END;\n" +
				"BEGIN\n" +
				"  p;\n" +
				"END pkg;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE PACKAGE BODY pkg IS\n  PROCEDURE p IS\n  BEGIN\n    NULL;\n  END;\nBEGIN\n  p;\nEND pkg;",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "package spec case expression without slash separator",
			sql: "CREATE PACKAGE pkg IS\n" +
				"  c CONSTANT NUMBER := CASE WHEN 1 = 1 THEN 1 ELSE 0 END;\n" +
				"END pkg;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE PACKAGE pkg IS\n  c CONSTANT NUMBER := CASE WHEN 1 = 1 THEN 1 ELSE 0 END;\nEND pkg;",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "type body without slash separator",
			sql: "CREATE TYPE BODY typ AS\n" +
				"  MEMBER FUNCTION f RETURN NUMBER IS\n" +
				"  BEGIN\n" +
				"    RETURN 1;\n" +
				"  END;\n" +
				"END;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE TYPE BODY typ AS\n  MEMBER FUNCTION f RETURN NUMBER IS\n  BEGIN\n    RETURN 1;\n  END;\nEND;",
				"\nCREATE TABLE t (id NUMBER)",
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
			name: "procedure local subprogram without slash separator",
			sql: "CREATE PROCEDURE p IS\n" +
				"  PROCEDURE q IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END q;\n" +
				"BEGIN\n" +
				"  q;\n" +
				"END;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE PROCEDURE p IS\n  PROCEDURE q IS\n  BEGIN\n    NULL;\n  END q;\nBEGIN\n  q;\nEND;",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "declare local subprogram without slash separator",
			sql: "DECLARE\n" +
				"  PROCEDURE q IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END q;\n" +
				"BEGIN\n" +
				"  q;\n" +
				"END;\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"DECLARE\n  PROCEDURE q IS\n  BEGIN\n    NULL;\n  END q;\nBEGIN\n  q;\nEND;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "ordinary begin backup is not plsql",
			sql:  "ALTER DATABASE BEGIN BACKUP;\nALTER DATABASE END BACKUP;",
			want: []string{"ALTER DATABASE BEGIN BACKUP", "\nALTER DATABASE END BACKUP"},
		},
		{
			name: "procedure call spec without end",
			sql: "CREATE PROCEDURE p AS LANGUAGE JAVA NAME 'Pkg.p()';\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE PROCEDURE p AS LANGUAGE JAVA NAME 'Pkg.p()'",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "trigger call body without end",
			sql: "CREATE TRIGGER trg BEFORE INSERT ON t CALL proc();\n" +
				"CREATE TABLE u (id NUMBER);",
			want: []string{
				"CREATE TRIGGER trg BEFORE INSERT ON t CALL proc()",
				"\nCREATE TABLE u (id NUMBER)",
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
		{
			name: "nested case expression without slash separator",
			sql: "BEGIN\n" +
				"  x := CASE WHEN a = 1 THEN CASE WHEN b = 1 THEN 1 ELSE 2 END ELSE 3 END;\n" +
				"END;\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"BEGIN\n  x := CASE WHEN a = 1 THEN CASE WHEN b = 1 THEN 1 ELSE 2 END ELSE 3 END;\nEND;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "compound trigger without slash separator",
			sql: "CREATE TRIGGER trg\n" +
				"FOR INSERT ON t\n" +
				"COMPOUND TRIGGER\n" +
				"  BEFORE EACH ROW IS\n" +
				"  BEGIN\n" +
				"    NULL;\n" +
				"  END BEFORE EACH ROW;\n" +
				"END trg;\n" +
				"CREATE TABLE t2 (id NUMBER);",
			want: []string{
				"CREATE TRIGGER trg\nFOR INSERT ON t\nCOMPOUND TRIGGER\n  BEFORE EACH ROW IS\n  BEGIN\n    NULL;\n  END BEFORE EACH ROW;\nEND trg;",
				"\nCREATE TABLE t2 (id NUMBER)",
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
		name      string
		sql       string
		want      []string
		wantKinds []SegmentKind
	}{
		{
			name: "environment commands are returned as line segments",
			sql: "SET DEFINE OFF\n" +
				"SET SERVEROUTPUT ON;\n" +
				"PROMPT creating table;\n" +
				"SPOOL install.log\n" +
				"SELECT 1 FROM dual;\n" +
				"SPOOL OFF\n" +
				"SELECT 2 FROM dual;",
			want: []string{
				"SET DEFINE OFF",
				"SET SERVEROUTPUT ON;",
				"PROMPT creating table;",
				"SPOOL install.log",
				"SELECT 1 FROM dual",
				"\nSPOOL OFF",
				"SELECT 2 FROM dual",
			},
			wantKinds: []SegmentKind{
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQL,
				SegmentSQLPlusCommand,
				SegmentSQL,
			},
		},
		{
			name: "script and session commands are returned as line segments",
			sql: "CONNECT scott/tiger@db\n" +
				"@preflight.sql\n" +
				"@@nested/install.sql arg1 arg2\n" +
				"START post.sql\n" +
				"WHENEVER SQLERROR EXIT SQL.SQLCODE ROLLBACK\n" +
				"SELECT 1 FROM dual;\n" +
				"EXIT SUCCESS",
			want: []string{
				"CONNECT scott/tiger@db",
				"@preflight.sql",
				"@@nested/install.sql arg1 arg2",
				"START post.sql",
				"WHENEVER SQLERROR EXIT SQL.SQLCODE ROLLBACK",
				"SELECT 1 FROM dual",
				"\nEXIT SUCCESS",
			},
			wantKinds: []SegmentKind{
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQL,
				SegmentSQLPlusCommand,
			},
		},
		{
			name: "remark and host commands are returned as line segments",
			sql: "REM this ; is a SQL*Plus comment\n" +
				"REMARK this is also ignored\n" +
				"HOST echo before\n" +
				"! echo shell\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"REM this ; is a SQL*Plus comment",
				"REMARK this is also ignored",
				"HOST echo before",
				"! echo shell",
				"SELECT 1 FROM dual",
			},
			wantKinds: []SegmentKind{
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQL,
			},
		},
		{
			name: "formatting and variable commands are returned as line segments",
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
			want: []string{
				"COLUMN c FORMAT A20",
				"BREAK ON report",
				"COMPUTE SUM OF sal ON report",
				"TTITLE left 'Report'",
				"BTITLE off",
				"DEFINE schema_name = HR",
				"UNDEFINE schema_name",
				"ACCEPT v PROMPT 'Value: '",
				"VARIABLE rc NUMBER",
				"PRINT rc",
				"SELECT 1 FROM dual",
			},
			wantKinds: []SegmentKind{
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQL,
			},
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
			if len(tt.wantKinds) > 0 {
				segs := Split(tt.sql)
				for i, want := range tt.wantKinds {
					if segs[i].Kind != want {
						t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, segs[i].Kind, segs[i].Text, want)
					}
				}
			}
		})
	}
}

func TestSplitClassifiesSQLPlusCommands(t *testing.T) {
	sql := "SET DEFINE OFF\nSELECT 1 FROM dual;\nPROMPT done\nSELECT 2 FROM dual;"
	got := Split(sql)
	wantKinds := []SegmentKind{
		SegmentSQLPlusCommand,
		SegmentSQL,
		SegmentSQLPlusCommand,
		SegmentSQL,
	}
	if len(got) != len(wantKinds) {
		t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(wantKinds))
	}
	for i, want := range wantKinds {
		if got[i].Kind != want {
			t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, got[i].Kind, got[i].Text, want)
		}
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

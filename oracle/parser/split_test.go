package parser

import (
	"os"
	"path/filepath"
	"testing"
)

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
			name: "create function with slash operator on its own line",
			sql: "CREATE OR REPLACE FUNCTION f RETURN NUMBER IS\n" +
				"BEGIN\n" +
				"  RETURN 10\n" +
				"  /\n" +
				"  2;\n" +
				"END;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE OR REPLACE FUNCTION f RETURN NUMBER IS\nBEGIN\n  RETURN 10\n  /\n  2;\nEND;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "wrapped procedure with slash separator",
			sql: "CREATE OR REPLACE PROCEDURE wrapped_proc WRAPPED\n" +
				"a000000\n" +
				"abcd\n" +
				"/\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE OR REPLACE PROCEDURE wrapped_proc WRAPPED\na000000\nabcd",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "wrapped procedure with semicolon separator",
			sql: "CREATE OR REPLACE PROCEDURE wrapped_proc WRAPPED\n" +
				"a000000\n" +
				"abcd;\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE OR REPLACE PROCEDURE wrapped_proc WRAPPED\na000000\nabcd",
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
			// The ';' belongs to the call spec: Oracle compiles the unit with
			// PLS-00103 (INVALID) when the text sent ends without it.
			name: "procedure call spec without end",
			sql: "CREATE PROCEDURE p AS LANGUAGE JAVA NAME 'Pkg.p()';\n" +
				"CREATE TABLE t (id NUMBER);",
			want: []string{
				"CREATE PROCEDURE p AS LANGUAGE JAVA NAME 'Pkg.p()';",
				"\nCREATE TABLE t (id NUMBER)",
			},
		},
		{
			name: "function call spec with repeated slash separators",
			sql: "CREATE OR REPLACE EDITIONABLE FUNCTION \"S\".\"F\" (p IN BLOB) RETURN VARCHAR2\n" +
				"AS LANGUAGE JAVA\n" +
				"NAME 'Signer.check(oracle.sql.BLOB) return java.lang.String';\n" +
				"/\n" +
				"/\n" +
				"\n" +
				"  CREATE OR REPLACE FUNCTION g (p IN BINARY_INTEGER) RETURN BINARY_INTEGER\n" +
				"AS EXTERNAL NAME \"c_g\" LIBRARY lib;\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE OR REPLACE EDITIONABLE FUNCTION \"S\".\"F\" (p IN BLOB) RETURN VARCHAR2\nAS LANGUAGE JAVA\nNAME 'Signer.check(oracle.sql.BLOB) return java.lang.String';",
				"\n\n  CREATE OR REPLACE FUNCTION g (p IN BINARY_INTEGER) RETURN BINARY_INTEGER\nAS EXTERNAL NAME \"c_g\" LIBRARY lib;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "package with call spec members",
			sql: "CREATE OR REPLACE PACKAGE BODY pk AS\n" +
				"  FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';\n" +
				"  PROCEDURE p IS BEGIN NULL; END;\n" +
				"  PROCEDURE q AS LANGUAGE C LIBRARY lib;\n" +
				"END pk;\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE OR REPLACE PACKAGE BODY pk AS\n  FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';\n  PROCEDURE p IS BEGIN NULL; END;\n  PROCEDURE q AS LANGUAGE C LIBRARY lib;\nEND pk;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			name: "nested call spec in a declare section",
			sql: "DECLARE\n" +
				"  FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';\n" +
				"BEGIN NULL; END;\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"DECLARE\n  FUNCTION f RETURN NUMBER AS LANGUAGE JAVA NAME 'X.f() return int';\nBEGIN NULL; END;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			// Every clause of the EXTERNAL form can open it; the unit ends at
			// its ';' instead of swallowing the statements after it.
			name: "external call specs with legacy clauses",
			sql: "CREATE OR REPLACE PROCEDURE p1 AS EXTERNAL LANGUAGE C NAME \"c_p\" LIBRARY lib;\n" +
				"CREATE OR REPLACE PROCEDURE p2 AS EXTERNAL CALLING STANDARD PASCAL LIBRARY lib;\n" +
				"CREATE OR REPLACE PROCEDURE p3(a IN BINARY_INTEGER) AS EXTERNAL PARAMETERS (a INT) LIBRARY lib;\n" +
				"CREATE OR REPLACE PROCEDURE p4 AS EXTERNAL WITH CONTEXT LIBRARY lib;\n" +
				"CREATE OR REPLACE PROCEDURE p5(a IN BINARY_INTEGER) AS EXTERNAL AGENT IN (a) LIBRARY lib;\n" +
				"CREATE TABLE t2 (a NUMBER);\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE OR REPLACE PROCEDURE p1 AS EXTERNAL LANGUAGE C NAME \"c_p\" LIBRARY lib;",
				"\nCREATE OR REPLACE PROCEDURE p2 AS EXTERNAL CALLING STANDARD PASCAL LIBRARY lib;",
				"\nCREATE OR REPLACE PROCEDURE p3(a IN BINARY_INTEGER) AS EXTERNAL PARAMETERS (a INT) LIBRARY lib;",
				"\nCREATE OR REPLACE PROCEDURE p4 AS EXTERNAL WITH CONTEXT LIBRARY lib;",
				"\nCREATE OR REPLACE PROCEDURE p5(a IN BINARY_INTEGER) AS EXTERNAL AGENT IN (a) LIBRARY lib;",
				"\nCREATE TABLE t2 (a NUMBER)",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			// Oracle 23ai compiles this VALID: only the word right after
			// IS|AS opens a call spec.
			name: "parameter and variable named like call spec words",
			sql: "CREATE FUNCTION f(language IN VARCHAR2, wrapped NUMBER) RETURN NUMBER IS\n" +
				"  x NUMBER;\n" +
				"  external NUMBER;\n" +
				"  mle NUMBER;\n" +
				"BEGIN RETURN 1; END;\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"CREATE FUNCTION f(language IN VARCHAR2, wrapped NUMBER) RETURN NUMBER IS\n  x NUMBER;\n  external NUMBER;\n  mle NUMBER;\nBEGIN RETURN 1; END;",
				"\nSELECT 1 FROM dual",
			},
		},
		{
			// The word after IS|AS commits to a call spec, as it does for
			// Oracle and the parser, so a malformed one still ends at its
			// ';' instead of swallowing the next statement.
			name: "malformed call specs end at their semicolon",
			sql: "CREATE PROCEDURE p AS LANGUAGE PYTHON;\n" +
				"SELECT 1 FROM dual;\n" +
				"CREATE PROCEDURE q AS EXTERNAL;\n" +
				"SELECT 2 FROM dual;\n" +
				"CREATE FUNCTION f RETURN NUMBER IS mle NUMBER;\n" +
				"SELECT 3 FROM dual;",
			want: []string{
				"CREATE PROCEDURE p AS LANGUAGE PYTHON;",
				"\nSELECT 1 FROM dual",
				"\nCREATE PROCEDURE q AS EXTERNAL;",
				"\nSELECT 2 FROM dual",
				"\nCREATE FUNCTION f RETURN NUMBER IS mle NUMBER;",
				"\nSELECT 3 FROM dual",
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

// TestSplitEmbeddedSource pins CREATE JAVA ... AS source_char and CREATE MLE
// MODULE ... AS module_text as SQL*Plus reads them: the text after AS runs to
// a line holding only "/", whatever it contains. Verified with SQL*Plus
// against Oracle 23ai: a SELECT after such source without a "/" line is sent
// as part of the statement. The delimited code of an inline MLE call spec
// ends at its closing delimiter instead, and the unit at the ';' after it.
func TestSplitEmbeddedSource(t *testing.T) {
	const sqlKind, srcKind = SegmentSQL, SegmentEmbeddedSource
	tests := []struct {
		name      string
		sql       string
		want      []string
		wantKinds []SegmentKind
	}{
		{
			name: "java sources with internal semicolons and repeated slashes",
			sql: "   CREATE JAVA SOURCE NAMED \"S\".\"Signer\" AS\n" +
				"import java.io.InputStream;import java.util.Base64;public final class Signer {    private Signer() { }    static int f() { int a = 1; return a; } }\n" +
				"/\n" +
				"/\n" +
				" \n" +
				"   CREATE JAVA SOURCE NAMED \"S\".\"Normalizer\" AS\n" +
				"import java.text.Normalizer;public final class Normalizer { }\n" +
				"/\n" +
				"/\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"   CREATE JAVA SOURCE NAMED \"S\".\"Signer\" AS\nimport java.io.InputStream;import java.util.Base64;public final class Signer {    private Signer() { }    static int f() { int a = 1; return a; } }",
				"\n \n   CREATE JAVA SOURCE NAMED \"S\".\"Normalizer\" AS\nimport java.text.Normalizer;public final class Normalizer { }",
				"\nSELECT 1 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, srcKind, sqlKind},
		},
		{
			name: "java text the SQL lexer would misread",
			sql: "CREATE OR REPLACE AND COMPILE JAVA SOURCE NAMED \"Q\" AS\n" +
				"public class Q {\n" +
				"  // don't split (here; a stray quote and paren\n" +
				"  static char c = '\\'';\n" +
				"  static String s = \"it's /\";\n" +
				"  static int g(int i) { i--; return i / 2; }\n" +
				"  /* SELECT 1 FROM dual; */\n" +
				"}\n" +
				"/\n" +
				"SELECT 2 FROM dual;",
			want: []string{
				"CREATE OR REPLACE AND COMPILE JAVA SOURCE NAMED \"Q\" AS\npublic class Q {\n  // don't split (here; a stray quote and paren\n  static char c = '\\'';\n  static String s = \"it's /\";\n  static int g(int i) { i--; return i / 2; }\n  /* SELECT 1 FROM dual; */\n}",
				"\nSELECT 2 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, sqlKind},
		},
		{
			name: "and resolve noforce java source without a slash runs to the end",
			sql: "CREATE OR REPLACE AND RESOLVE NOFORCE JAVA SOURCE NAMED r AS public class R { int a; }\n" +
				"SELECT 3 FROM dual;",
			want: []string{
				"CREATE OR REPLACE AND RESOLVE NOFORCE JAVA SOURCE NAMED r AS public class R { int a; }\nSELECT 3 FROM dual;",
			},
			wantKinds: []SegmentKind{srcKind},
		},
		{
			// Only the line break before the "/" line is framing; trailing
			// blanks and blank lines are part of source_char.
			name: "trailing blanks of the java source are kept",
			sql:  "CREATE JAVA SOURCE NAMED t AS\r\npublic class T { }  \t\r\n\r\n  /\r\nSELECT 5 FROM dual;",
			want: []string{
				"CREATE JAVA SOURCE NAMED t AS\r\npublic class T { }  \t\r\n",
				"\r\nSELECT 5 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, sqlKind},
		},
		{
			name: "java statements without source text end at semicolons",
			sql: "CREATE JAVA CLASS USING BFILE (dir, 'X.class');\n" +
				"CREATE JAVA SOURCE NAMED s USING CLOB (SELECT src AS text FROM t);\n" +
				"SELECT 4 FROM dual;",
			want: []string{
				"CREATE JAVA CLASS USING BFILE (dir, 'X.class')",
				"\nCREATE JAVA SOURCE NAMED s USING CLOB (SELECT src AS text FROM t)",
				"\nSELECT 4 FROM dual",
			},
			wantKinds: []SegmentKind{sqlKind, sqlKind, sqlKind},
		},
		{
			name: "mle module source runs to the slash line",
			sql: "CREATE OR REPLACE MLE MODULE m LANGUAGE JAVASCRIPT AS\n" +
				"export function f(a) {\n" +
				"  // don't split; here\n" +
				"  let s = \"it's /\";\n" +
				"  return a--;\n" +
				"}\n" +
				"/\n" +
				"SELECT 6 FROM dual;",
			want: []string{
				"CREATE OR REPLACE MLE MODULE m LANGUAGE JAVASCRIPT AS\nexport function f(a) {\n  // don't split; here\n  let s = \"it's /\";\n  return a--;\n}",
				"\nSELECT 6 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, sqlKind},
		},
		{
			name: "mle module without source text ends at its semicolon",
			sql: "CREATE MLE MODULE m LANGUAGE JAVASCRIPT USING BFILE (dir, 'm.js');\n" +
				"SELECT 7 FROM dual;",
			want: []string{
				"CREATE MLE MODULE m LANGUAGE JAVASCRIPT USING BFILE (dir, 'm.js')",
				"\nSELECT 7 FROM dual",
			},
			wantKinds: []SegmentKind{sqlKind, sqlKind},
		},
		{
			name: "inline mle call spec ends at the semicolon after its code",
			sql: "CREATE OR REPLACE FUNCTION f(a IN NUMBER) RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{\n" +
				"  // don't split; here\n" +
				"  let s = \"it's\"; return a--;\n" +
				"}};\n" +
				"SELECT 8 FROM dual;",
			want: []string{
				"CREATE OR REPLACE FUNCTION f(a IN NUMBER) RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{\n  // don't split; here\n  let s = \"it's\"; return a--;\n}};",
				"\nSELECT 8 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, sqlKind},
		},
		{
			name: "inline mle call spec with a quote delimiter",
			sql: "CREATE PROCEDURE p AS MLE LANGUAGE JAVASCRIPT PURE 'q console.log(1); 'q;\n" +
				"SELECT 9 FROM dual;",
			want: []string{
				"CREATE PROCEDURE p AS MLE LANGUAGE JAVASCRIPT PURE 'q console.log(1); 'q;",
				"\nSELECT 9 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, sqlKind},
		},
		{
			name: "inline mle call spec in a package body",
			sql: "CREATE OR REPLACE PACKAGE BODY pk AS\n" +
				"  FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT ## return \"a;b\"; -- x ##;\n" +
				"  PROCEDURE z IS BEGIN NULL; END;\n" +
				"END;\n" +
				"/\n" +
				"SELECT 10 FROM dual;",
			want: []string{
				"CREATE OR REPLACE PACKAGE BODY pk AS\n  FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT ## return \"a;b\"; -- x ##;\n  PROCEDURE z IS BEGIN NULL; END;\nEND;",
				"\nSELECT 10 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, sqlKind},
		},
		{
			// A type specification ends at its ';' like SQL, but not at one
			// inside the inline code of a method's MLE call spec.
			name: "inline mle call spec in an object type spec",
			sql: "CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ s = \"--\"; return 1; }}, MEMBER PROCEDURE p);\n" +
				"CREATE TYPE u AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }})\n" +
				"/\n" +
				"SELECT 12 FROM dual;",
			want: []string{
				"CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ s = \"--\"; return 1; }}, MEMBER PROCEDURE p)",
				"\nCREATE TYPE u AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER AS MLE LANGUAGE JAVASCRIPT {{ return 1; }})",
				"\nSELECT 12 FROM dual",
			},
			wantKinds: []SegmentKind{srcKind, srcKind, sqlKind},
		},
		{
			name: "mle module call spec holds no embedded code",
			sql: "CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m ENV e SIGNATURE 'f()';\n" +
				"SELECT 11 FROM dual;",
			want: []string{
				"CREATE FUNCTION f RETURN NUMBER AS MLE MODULE m ENV e SIGNATURE 'f()';",
				"\nSELECT 11 FROM dual",
			},
			wantKinds: []SegmentKind{sqlKind, sqlKind},
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
			for i, seg := range Split(tt.sql) {
				if seg.Kind != tt.wantKinds[i] {
					t.Fatalf("segment[%d] kind = %v, want %v", i, seg.Kind, tt.wantKinds[i])
				}
				if _, err := ParseRange(tt.sql, seg.ByteStart, seg.ByteEnd); err != nil {
					t.Fatalf("ParseRange(%q): %v", seg.Text, err)
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
			name: "run does not flush unterminated SQL buffer",
			sql:  "SELECT 1 FROM dual\nRUN\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual\nRUN\nSELECT 2 FROM dual"},
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
		{
			name: "sqlplus commands after plsql slash delimiter",
			sql: "BEGIN\n" +
				"  NULL;\n" +
				"END;\n" +
				"/\n" +
				"PROMPT block complete\n" +
				"SPOOL install.log\n" +
				"SELECT 1 FROM dual;",
			want: []string{
				"BEGIN\n  NULL;\nEND;",
				"\nPROMPT block complete",
				"SPOOL install.log",
				"SELECT 1 FROM dual",
			},
			wantKinds: []SegmentKind{
				SegmentSQL,
				SegmentSQLPlusCommand,
				SegmentSQLPlusCommand,
				SegmentSQL,
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

func TestSplitClassifiesOracleSetStatementsAsSQL(t *testing.T) {
	sql := "SET DEFINE OFF\n" +
		"SET TRANSACTION READ ONLY;\n" +
		"SET ROLE app_role;\n" +
		"SET CONSTRAINTS ALL IMMEDIATE;"
	got := Split(sql)
	wantTexts := []string{
		"SET DEFINE OFF",
		"SET TRANSACTION READ ONLY",
		"\nSET ROLE app_role",
		"\nSET CONSTRAINTS ALL IMMEDIATE",
	}
	wantKinds := []SegmentKind{
		SegmentSQLPlusCommand,
		SegmentSQL,
		SegmentSQL,
		SegmentSQL,
	}
	if len(got) != len(wantKinds) {
		t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(wantKinds))
	}
	for i := range wantKinds {
		if got[i].Text != wantTexts[i] {
			t.Fatalf("segment[%d] Text = %q, want %q", i, got[i].Text, wantTexts[i])
		}
		if got[i].Kind != wantKinds[i] {
			t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, got[i].Kind, got[i].Text, wantKinds[i])
		}
	}
}

func TestSplitDoesNotClassifySQLContinuationLinesAsSQLPlus(t *testing.T) {
	sql := "SELECT employee_id\n" +
		"FROM employees\n" +
		"START WITH manager_id IS NULL\n" +
		"CONNECT BY PRIOR employee_id = manager_id;\n" +
		"CREATE DATABASE LINK remote_db\n" +
		"CONNECT TO remote_user IDENTIFIED BY remote_pass\n" +
		"USING 'remote_tns';\n" +
		"CREATE DATABASE mydb\n" +
		"SET DEFAULT BIGFILE TABLESPACE;"
	got := Split(sql)
	wantTexts := []string{
		"SELECT employee_id\nFROM employees\nSTART WITH manager_id IS NULL\nCONNECT BY PRIOR employee_id = manager_id",
		"\nCREATE DATABASE LINK remote_db\nCONNECT TO remote_user IDENTIFIED BY remote_pass\nUSING 'remote_tns'",
		"\nCREATE DATABASE mydb\nSET DEFAULT BIGFILE TABLESPACE",
	}
	if len(got) != len(wantTexts) {
		t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(wantTexts))
	}
	for i := range wantTexts {
		if got[i].Text != wantTexts[i] {
			t.Fatalf("segment[%d] Text = %q, want %q", i, got[i].Text, wantTexts[i])
		}
		if got[i].Kind != SegmentSQL {
			t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, got[i].Kind, got[i].Text, SegmentSQL)
		}
	}
}

func TestSplitRequiresExplicitSQLTerminatorBeforeSQLPlusCommand(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "table named like remark command after from",
			sql:  "SELECT * FROM\nremark",
			want: []string{"SELECT * FROM\nremark"},
		},
		{
			name: "column named like prompt command in select list",
			sql:  "SELECT\nprompt\nFROM t",
			want: []string{"SELECT\nprompt\nFROM t"},
		},
		{
			name: "sqlplus command after semicolon",
			sql:  "SELECT 1 FROM dual;\nPROMPT done\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", "\nPROMPT done", "SELECT 2 FROM dual"},
		},
		{
			name: "slash delimiter explicitly terminates sql before sqlplus command",
			sql:  "SELECT 1 FROM dual\n/\nPROMPT done\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", "\nPROMPT done", "SELECT 2 FROM dual"},
		},
		{
			name: "line comment after semicolon can precede sqlplus command",
			sql:  "SELECT 1 FROM dual;\n-- setup output\nPROMPT done\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", "\n-- setup output\nPROMPT done", "SELECT 2 FROM dual"},
		},
		{
			name: "block comment after semicolon can precede sqlplus command",
			sql:  "SELECT 1 FROM dual;\n/* setup output */\nSPOOL out.log\nSELECT 2 FROM dual;",
			want: []string{"SELECT 1 FROM dual", "\n/* setup output */\nSPOOL out.log", "SELECT 2 FROM dual"},
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

func TestSplitDoesNotClassifySQLPlusWordsInsideUnterminatedSQL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "update set define column",
			sql: "UPDATE accounts\n" +
				"SET DEFINE = 1\n" +
				"WHERE id = 1;",
		},
		{
			name: "merge update set prompt column",
			sql: "MERGE INTO dst d USING src s ON (d.id = s.id)\n" +
				"WHEN MATCHED THEN UPDATE\n" +
				"SET PROMPT = s.prompt;",
		},
		{
			name: "match recognize define clause",
			sql: "SELECT *\n" +
				"FROM trades\n" +
				"MATCH_RECOGNIZE (\n" +
				"  ORDER BY ts\n" +
				"  MEASURES A.price AS start_price\n" +
				"  PATTERN (A B)\n" +
				"  DEFINE\n" +
				"    B AS B.price > A.price\n" +
				") mr;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.sql)
			if len(got) != 1 {
				t.Fatalf("got %d segments %q, want 1", len(got), splitTexts(got))
			}
			want := tt.sql[:len(tt.sql)-1]
			if got[0].Text != want {
				t.Fatalf("segment Text = %q, want %q", got[0].Text, want)
			}
			if got[0].Kind != SegmentSQL {
				t.Fatalf("segment Kind = %v, want %v", got[0].Kind, SegmentSQL)
			}
		})
	}
}

func TestSplitOrdinarySQLStatementsWithSQLPlusLikeLines(t *testing.T) {
	sql := "CREATE VIEW v_sqlplus_words AS\n" +
		"SELECT\n" +
		"  prompt,\n" +
		"  remark,\n" +
		"  spool\n" +
		"FROM source_table;\n" +
		"SELECT employee_id\n" +
		"FROM employees\n" +
		"START WITH manager_id IS NULL\n" +
		"CONNECT BY PRIOR employee_id = manager_id;\n" +
		"UPDATE accounts\n" +
		"SET DEFINE = 1,\n" +
		"    PROMPT = 'on'\n" +
		"WHERE id = 1;\n" +
		"MERGE INTO dst d USING src s ON (d.id = s.id)\n" +
		"WHEN MATCHED THEN UPDATE\n" +
		"SET SPOOL = s.spool;\n" +
		"SELECT *\n" +
		"FROM trades\n" +
		"MATCH_RECOGNIZE (\n" +
		"  ORDER BY ts\n" +
		"  MEASURES A.price AS start_price\n" +
		"  PATTERN (A B)\n" +
		"  DEFINE\n" +
		"    B AS B.price > A.price\n" +
		") mr;"
	got := Split(sql)
	wantTexts := []string{
		"CREATE VIEW v_sqlplus_words AS\nSELECT\n  prompt,\n  remark,\n  spool\nFROM source_table",
		"\nSELECT employee_id\nFROM employees\nSTART WITH manager_id IS NULL\nCONNECT BY PRIOR employee_id = manager_id",
		"\nUPDATE accounts\nSET DEFINE = 1,\n    PROMPT = 'on'\nWHERE id = 1",
		"\nMERGE INTO dst d USING src s ON (d.id = s.id)\nWHEN MATCHED THEN UPDATE\nSET SPOOL = s.spool",
		"\nSELECT *\nFROM trades\nMATCH_RECOGNIZE (\n  ORDER BY ts\n  MEASURES A.price AS start_price\n  PATTERN (A B)\n  DEFINE\n    B AS B.price > A.price\n) mr",
	}
	if len(got) != len(wantTexts) {
		t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(wantTexts))
	}
	for i := range wantTexts {
		if got[i].Text != wantTexts[i] {
			t.Fatalf("segment[%d] Text = %q, want %q", i, got[i].Text, wantTexts[i])
		}
		if got[i].Kind != SegmentSQL {
			t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, got[i].Kind, got[i].Text, SegmentSQL)
		}
	}
}

func TestSplitDoesNotClassifyRemarkColumnAsSQLPlus(t *testing.T) {
	sql := `create table NCF_DI_PO_SOURCE_TEMP
(
    data_id             number not null
        constraint NC_DI_pk
            primary key,
    purchase_receive_id number,
    po_number           varchar2(60),
    sku_code            varchar2(100),
    qty                 number,
    receive_number      varchar2(100),
    ETD                 date,
    ETA                 date,
    ATD                 date,
    SO_NUMBER           varchar2(100),
    SO_line_amount      number(10,2),
    source              varchar2(100),
    remark              varchar2(100),
    prompt              varchar2(100),
    spool               varchar2(100),
    attribute1          varchar2(100),
    attribute2          varchar2(100)
)`
	got := Split(sql)
	if len(got) != 1 {
		t.Fatalf("got %d segments %q, want 1", len(got), splitTexts(got))
	}
	if got[0].Text != sql {
		t.Fatalf("segment Text = %q, want original SQL", got[0].Text)
	}
	if got[0].Kind != SegmentSQL {
		t.Fatalf("segment Kind = %v, want %v", got[0].Kind, SegmentSQL)
	}
}

func TestSplitDoesNotTreatQualifiedRAtLineStartAsRun(t *testing.T) {
	sql := "CREATE VIEW v AS SELECT\n" +
		"  r.id,\n" +
		"  r rows,\n" +
		"  r.name\n" +
		"FROM records r;"
	got := Split(sql)
	wantTexts := []string{
		"CREATE VIEW v AS SELECT\n  r.id,\n  r rows,\n  r.name\nFROM records r",
	}
	if len(got) != len(wantTexts) {
		t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(wantTexts))
	}
	for i := range wantTexts {
		if got[i].Text != wantTexts[i] {
			t.Fatalf("segment[%d] Text = %q, want %q", i, got[i].Text, wantTexts[i])
		}
		if got[i].Kind != SegmentSQL {
			t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, got[i].Kind, got[i].Text, SegmentSQL)
		}
	}
}

func TestSplitClassifiesSQLPlusCommandsAfterExplicitSQLTerminators(t *testing.T) {
	sql := "SELECT 1 FROM dual;\n" +
		"PROMPT running next query\n" +
		"SPOOL install.log\n" +
		"SELECT 2 FROM dual;\n" +
		"SET DEFINE OFF\n" +
		"SELECT 3 FROM dual;\n" +
		"SET DEF OFF\n" +
		"SELECT 4 FROM dual;\n" +
		"SET SERVEROUT ON\n" +
		"SELECT 3 FROM dual;\n" +
		"CONNECT scott/tiger@db\n" +
		"SELECT 2 FROM dual;"
	got := Split(sql)
	wantTexts := []string{
		"SELECT 1 FROM dual",
		"\nPROMPT running next query",
		"SPOOL install.log",
		"SELECT 2 FROM dual",
		"\nSET DEFINE OFF",
		"SELECT 3 FROM dual",
		"\nSET DEF OFF",
		"SELECT 4 FROM dual",
		"\nSET SERVEROUT ON",
		"SELECT 3 FROM dual",
		"\nCONNECT scott/tiger@db",
		"SELECT 2 FROM dual",
	}
	wantKinds := []SegmentKind{
		SegmentSQL,
		SegmentSQLPlusCommand,
		SegmentSQLPlusCommand,
		SegmentSQL,
		SegmentSQLPlusCommand,
		SegmentSQL,
		SegmentSQLPlusCommand,
		SegmentSQL,
		SegmentSQLPlusCommand,
		SegmentSQL,
		SegmentSQLPlusCommand,
		SegmentSQL,
	}
	if len(got) != len(wantTexts) {
		t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(wantTexts))
	}
	for i := range wantTexts {
		if got[i].Text != wantTexts[i] {
			t.Fatalf("segment[%d] Text = %q, want %q", i, got[i].Text, wantTexts[i])
		}
		if got[i].Kind != wantKinds[i] {
			t.Fatalf("segment[%d] Kind = %v for %q, want %v", i, got[i].Kind, got[i].Text, wantKinds[i])
		}
	}
}

func TestSplitDoesNotClassifyValidCorpusStatementsAsSQLPlus(t *testing.T) {
	corpusDir := filepath.Join("..", "quality", "corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		corpusDir = filepath.Join("oracle", "quality", "corpus")
		entries, err = os.ReadDir(corpusDir)
		if err != nil {
			t.Fatalf("cannot read corpus directory: %v", err)
		}
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		path := filepath.Join(corpusDir, entry.Name())
		for _, stmt := range loadCorpusFile(t, path) {
			if stmt.valid != "true" {
				continue
			}
			for _, seg := range Split(stmt.sql) {
				if seg.Kind == SegmentSQLPlusCommand {
					t.Fatalf("%s/%s classified valid SQL as SQL*Plus command: %q", entry.Name(), stmt.name, seg.Text)
				}
			}
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

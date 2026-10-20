package parser

import (
	"errors"
	"strings"
	"testing"
)

// TestPLSQLMissingSemicolonMessages pins the error for PL/SQL text that
// lacks a required ';'. Oracle 23.26 compiles each case with PLS-00103 at the
// same token, expecting ';', and SQL*Plus sends the text without adding the
// ';', so a unit written without it is INVALID there too. The message names
// the missing ';' where that is the likely mistake: after END, a call spec,
// or a SUBTYPE declaration, where ';' is the only continuation, and between
// statements or declarations when the next token closes the list or starts a
// new line with a word. Elsewhere the item may have stopped early on syntax
// the parser does not take, and the message stays plain.
func TestPLSQLMissingSemicolonMessages(t *testing.T) {
	const atEnd = ""
	tests := []struct {
		sql  string
		at   string // text at the error position; atEnd for the end of input
		want string
	}{
		{"BEGIN NULL; END", atEnd, "syntax error at end of input: expected ';' after END"},
		{"CREATE OR REPLACE PROCEDURE p AS BEGIN NULL; END p", atEnd, "syntax error at end of input: expected ';' after END"},
		{"CREATE OR REPLACE PACKAGE pk AS PROCEDURE q; END", atEnd, "syntax error at end of input: expected ';' after END"},
		{"CREATE OR REPLACE TYPE BODY t AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END", atEnd, "syntax error at end of input: expected ';' after END"},
		{"CREATE TRIGGER trg FOR INSERT ON t COMPOUND TRIGGER BEFORE STATEMENT IS BEGIN NULL; END BEFORE STATEMENT; END trg", atEnd, "syntax error at end of input: expected ';' after END"},
		{"BEGIN IF a = 1 THEN NULL; END IF END;", "END;", `syntax error at or near "END": expected ';' after END IF`},
		{"BEGIN CASE WHEN a = 1 THEN NULL; END CASE END;", "END;", `syntax error at or near "END": expected ';' after END CASE`},
		{"CREATE PROCEDURE p AS LANGUAGE JAVA NAME 'X.y()'", atEnd, "syntax error at end of input: expected ';' after the call spec"},
		{"DECLARE SUBTYPE s IS NUMBER NOT NULL\nBEGIN NULL; END;", "BEGIN NULL", `syntax error at or near "BEGIN": expected ';' after the SUBTYPE declaration`},
		// Between statements: the next token closes the list, or starts a
		// new line with a word.
		{"BEGIN NULL END;", "END;", `syntax error at or near "END": expected ';' after the previous statement`},
		{"BEGIN\n  NULL\nEND;", "END;", `syntax error at or near "END": expected ';' after the previous statement`},
		{"BEGIN\n  NULL\n  NULL;\nEND;", "NULL;", `syntax error at or near "NULL": expected ';' after the previous statement`},
		{"BEGIN\n  x := 1\n  y := 2;\nEND;", "y :=", `syntax error at or near "y": expected ';' after the previous statement`},
		{"BEGIN\n  UPDATE t SET a = 1\n  COMMIT;\nEND;", "COMMIT", `syntax error at or near "COMMIT": expected ';' after the previous statement`},
		{"BEGIN\n  FOR i IN 1 .. 2 LOOP NULL; END LOOP\n  NULL;\nEND;", "NULL;\nEND", `syntax error at or near "NULL": expected ';' after the previous statement`},
		// Between declarations.
		{"DECLARE x NUMBER BEGIN NULL; END;", "BEGIN", `syntax error at or near "BEGIN": expected ';' after the previous declaration`},
		{"DECLARE\n  x NUMBER\n  y NUMBER;\nBEGIN NULL; END;", "y NUMBER", `syntax error at or near "y": expected ';' after the previous declaration`},
		{"CREATE OR REPLACE PACKAGE BODY pk AS\n  x NUMBER\n  PROCEDURE q IS BEGIN NULL; END;\nEND;", "PROCEDURE", `syntax error at or near "PROCEDURE": expected ';' after the previous declaration`},
		// The same line, or a token that does not start a word: plain.
		{"BEGIN x := 1 y := 2; END;", "y :=", `syntax error at or near "y"`},
		{"DECLARE x NUMBER y NUMBER; BEGIN NULL; END;", "y NUMBER", `syntax error at or near "y"`},
		{"DECLARE TYPE b IS (FALSE, TRUE); BEGIN NULL; END;", "(FALSE", `syntax error at or near "("`},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := Parse(tt.sql)
			var perr *ParseError
			if !errors.As(err, &perr) {
				t.Fatalf("Parse error = %v, want a *ParseError", err)
			}
			wantPos := len(tt.sql)
			if tt.at != atEnd {
				wantPos = strings.Index(tt.sql, tt.at)
			}
			if perr.Message != tt.want || perr.Position != wantPos {
				t.Fatalf("got %q at %d, want %q at %d", perr.Message, perr.Position, tt.want, wantPos)
			}
		})
	}
}

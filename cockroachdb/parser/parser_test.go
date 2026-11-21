package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParseErrorContract(t *testing.T) {
	cases := []struct {
		name     string
		sql      string
		wantPos  int
		wantCode string
		wantMsg  string
	}{
		{
			name:     "syntax error mid-statement",
			sql:      "SELECT FROM FROM",
			wantPos:  12,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "syntax error at byte zero",
			sql:      "FROM",
			wantPos:  0,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "syntax error at EOF",
			sql:      "SELECT (",
			wantPos:  8,
			wantCode: "42601",
			wantMsg:  "EOF",
		},
		{
			name:     "multi-byte UTF-8 byte offset",
			sql:      "SELECT 'é', 'ü' FROM FROM",
			wantPos:  23,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "leading whitespace and comment",
			sql:      "  -- comment\n  SELECT FROM FROM",
			wantPos:  27,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "multi-line SQL",
			sql:      "SELECT 1\nFROM FROM",
			wantPos:  14,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "multi-statement failure in second statement",
			sql:      "SELECT 1; SELECT FROM",
			wantPos:  21,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "multi-statement failure before trailing valid statement",
			sql:      "SELECT 1; SELECT FROM; SELECT 2",
			wantPos:  21,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "lexical unterminated string",
			sql:      "SELECT 'unterminated",
			wantPos:  7,
			wantCode: "42601",
			wantMsg:  "unterminated string",
		},
		{
			name:     "lexical unterminated comment",
			sql:      "SELECT 1 /* unterminated",
			wantPos:  9,
			wantCode: "42601",
			wantMsg:  "unterminated comment",
		},
		{
			name:     "nested PL/pgSQL DO block error",
			sql:      "DO $$ BEGIN INVALID SYNTAX; END $$",
			wantPos:  12,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
		{
			name:     "unimplemented grammar action",
			sql:      "CREATE TABLE a (b INT, CONSTRAINT c CHECK (b > 0) DEFERRABLE)",
			wantPos:  60,
			wantCode: "0A000",
			wantMsg:  "unimplemented",
		},
		{
			name:     "help token has unknown position",
			sql:      "??",
			wantPos:  -1,
			wantCode: "42601",
			wantMsg:  "help token",
		},
		{
			name:     "binary input returns ParseError without panicking",
			sql:      "\x01",
			wantPos:  0,
			wantCode: "42601",
			wantMsg:  "syntax error",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) unexpectedly succeeded with %d statements", tc.sql, len(stmts))
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %T, want *ParseError", tc.sql, err)
			}
			if pe.Position != tc.wantPos {
				t.Errorf("Position = %d, want %d (err=%q)", pe.Position, tc.wantPos, pe.Error())
			}
			if pe.Code != tc.wantCode {
				t.Errorf("Code = %q, want %q", pe.Code, tc.wantCode)
			}
			if !strings.Contains(pe.Message, tc.wantMsg) {
				t.Errorf("Message = %q, want substring %q", pe.Message, tc.wantMsg)
			}
			if !strings.Contains(pe.Error(), pe.Message) {
				t.Errorf("Error() = %q, want it to contain Message %q", pe.Error(), pe.Message)
			}
		})
	}
}

func TestParseSuccessAndStatementText(t *testing.T) {
	stmts, err := Parse("SELECT 1; ; CREATE TABLE t (id INT)")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(stmts) != 2 {
		t.Fatalf("len(stmts) = %d, want 2", len(stmts))
	}
	if stmts[0].Text() != "SELECT 1" {
		t.Errorf("stmts[0].Text() = %q, want %q", stmts[0].Text(), "SELECT 1")
	}
	if stmts[1].Text() != "CREATE TABLE t (id INT)" {
		t.Errorf("stmts[1].Text() = %q, want %q", stmts[1].Text(), "CREATE TABLE t (id INT)")
	}
	for _, empty := range []string{"", "   ", ";", ";;", "\x00"} {
		stmts, err := Parse(empty)
		if err != nil {
			t.Errorf("Parse(%q): unexpected error %v", empty, err)
		}
		if len(stmts) != 0 {
			t.Errorf("Parse(%q): len = %d, want 0", empty, len(stmts))
		}
	}
}

package parser

import (
	"strings"
	"testing"
)

// Oracle block comments do not nest: the first */ terminates the comment
// regardless of /* sequences inside it. Engine-verified on 11gR2 in both
// directions (BYT-9963):
//
//	SELECT 1 /* a /* b */ FROM dual        -- succeeds (comment ends at first */)
//	SELECT 2 /* a /* b */ */ FROM dual     -- ORA-00936 (legal only under nesting)
//
// PostgreSQL nests block comments; that is a genuine per-engine divergence,
// so these fences are Oracle-only and the PG lexer keeps its counting.

func TestLexerNonNestedBlockComment(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"innerOpenIgnored", "SELECT 1 /* a /* b */ FROM dual", false},
		{"strayCloseAfterComment", "SELECT 2 /* a /* b */ */ FROM dual", true},
		{"plainComment", "SELECT 1 /* plain */ FROM dual", false},
		{"unterminated", "SELECT 1 /* never closed", true},
		{"unterminatedWithInnerOpen", "SELECT 1 /* a /* b", true},
		{"hintEndsAtFirstClose", "SELECT /*+ FULL(t) /* x */ c FROM t", false},
		{"division", "SELECT 4 / 2 FROM dual", false},
		{"openMarkerInString", "SELECT '/*' FROM dual", false},
		{"closeMarkerInString", "SELECT '*/' FROM dual", false},
		{"openMarkerInLineComment", "SELECT 1 FROM dual -- /*\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if tt.wantErr && err == nil {
				t.Fatalf("want error, got none")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want success, got %v", err)
			}
		})
	}
}

// The customer shape (BYT-9963, MBBank): a commented-out code region that
// itself contains a /* marker, inside a package body terminated by a
// standalone SQL*Plus slash. The slash must be recognized as a delimiter,
// not retained inside the SQL segment.
func TestSplitNonNestedCommentSlashDelimiter(t *testing.T) {
	sql := "create or replace procedure p as\nbegin\n  /* outer /* inner */\n  null;\nend;\n/\n"
	segs := Split(sql)

	// The standalone slash is a pure delimiter: it must not appear in any
	// segment (see TestSplitSlashDelimiter for the established contract).
	var sqlSegs int
	for _, s := range segs {
		if s.Kind == SegmentSQLPlusCommand {
			continue
		}
		sqlSegs++
		if strings.Contains(s.Text, "\n/") || strings.HasSuffix(strings.TrimSpace(s.Text), "/") {
			t.Fatalf("SQL segment retains slash delimiter: %q", s.Text)
		}
	}
	if sqlSegs != 1 {
		t.Fatalf("got %d SQL segments, want 1 (all: %#v)", sqlSegs, segs)
	}
}

// Two-unit script in the customer file's structure: spec, slash, body whose
// nested-marker comment must not swallow the terminator, slash.
func TestSplitCustomerShapePackagePair(t *testing.T) {
	sql := strings.Join([]string{
		"create or replace package pkg as\n  procedure q;\nend pkg;",
		"/",
		"create or replace package body pkg as\n  procedure q as\n  begin\n    /*old := 1;\n      /*\n      if x then\n        y;\n      end if;*/\n    null;\n  end;\nend pkg;",
		"/",
		"",
	}, "\n")
	segs := Split(sql)

	var sqlTexts []string
	for _, s := range segs {
		if s.Kind == SegmentSQLPlusCommand {
			continue
		}
		sqlTexts = append(sqlTexts, s.Text)
	}
	if len(sqlTexts) != 2 {
		t.Fatalf("got %d SQL segments, want 2", len(sqlTexts))
	}
	for i, txt := range sqlTexts {
		if strings.HasSuffix(strings.TrimSpace(txt), "/") {
			t.Fatalf("SQL segment %d retains trailing slash", i)
		}
	}
	if !strings.Contains(sqlTexts[1], "end pkg;") {
		t.Fatalf("body segment truncated: %q", sqlTexts[1][:80])
	}
}

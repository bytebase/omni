package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

func TestParseEmpty(t *testing.T) {
	file, errs := Parse("")
	if file == nil {
		t.Fatal("expected non-nil File")
	}
	if len(file.Stmts) != 0 {
		t.Errorf("expected 0 stmts, got %d", len(file.Stmts))
	}
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d", len(errs))
	}
}

func TestParseUnsupported(t *testing.T) {
	// All statement types are stubbed as unsupported in F4.
	file, errs := Parse("SELECT 1")
	if file == nil {
		t.Fatal("expected non-nil File")
	}
	// unsupported returns nil node, so Stmts should be empty.
	if len(file.Stmts) != 0 {
		t.Errorf("expected 0 stmts (all unsupported), got %d", len(file.Stmts))
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Msg != "SELECT statement parsing is not yet supported" {
		t.Errorf("unexpected error: %q", errs[0].Msg)
	}
}

func TestParseMultipleUnsupported(t *testing.T) {
	file, errs := Parse("SELECT 1; INSERT INTO t VALUES (1); CREATE TABLE t (id INT)")
	if file == nil {
		t.Fatal("expected non-nil File")
	}
	if len(file.Stmts) != 0 {
		t.Errorf("expected 0 stmts, got %d", len(file.Stmts))
	}
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d: %v", len(errs), errs)
	}
}

func TestParseUnknownStatement(t *testing.T) {
	_, errs := Parse("FOOBAR something")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Msg != "unknown or unsupported statement starting with FOOBAR" {
		t.Errorf("unexpected error: %q", errs[0].Msg)
	}
}

func TestParseFileLocCoversInput(t *testing.T) {
	input := "SELECT 1; SELECT 2"
	file, _ := Parse(input)
	if file.Loc != (ast.Loc{Start: 0, End: len(input)}) {
		t.Errorf("File.Loc = %v, want {0, %d}", file.Loc, len(input))
	}
}

func TestParserAdvanceAndPeek(t *testing.T) {
	p := &Parser{
		lexer: NewLexer("SELECT 1"),
	}
	p.advance() // prime cur

	// cur should be SELECT
	if p.cur.Kind != kwSELECT {
		t.Fatalf("cur = %d, want kwSELECT", p.cur.Kind)
	}

	// peek should return SELECT without consuming
	if p.peek().Kind != kwSELECT {
		t.Errorf("peek() = %d, want kwSELECT", p.peek().Kind)
	}

	// peekNext should return 1 (tokInt) without consuming
	next := p.peekNext()
	if next.Kind != tokInt {
		t.Errorf("peekNext() = %d, want tokInt", next.Kind)
	}
	// cur should still be SELECT
	if p.cur.Kind != kwSELECT {
		t.Errorf("cur after peekNext = %d, want kwSELECT", p.cur.Kind)
	}

	// advance should return SELECT and move to 1
	prev := p.advance()
	if prev.Kind != kwSELECT {
		t.Errorf("advance() returned %d, want kwSELECT", prev.Kind)
	}
	if p.cur.Kind != tokInt {
		t.Errorf("cur after advance = %d, want tokInt", p.cur.Kind)
	}
}

func TestParserMatch(t *testing.T) {
	p := &Parser{
		lexer: NewLexer("SELECT 1"),
	}
	p.advance()

	// match with wrong type should not consume
	_, ok := p.match(kwINSERT, kwUPDATE)
	if ok {
		t.Error("match should not have matched")
	}

	// match with correct type should consume
	tok, ok := p.match(kwSELECT)
	if !ok {
		t.Error("match should have matched kwSELECT")
	}
	if tok.Kind != kwSELECT {
		t.Errorf("matched token = %d, want kwSELECT", tok.Kind)
	}
}

func TestParserExpect(t *testing.T) {
	p := &Parser{
		lexer: NewLexer("SELECT 1"),
	}
	p.advance()

	// expect correct type
	tok, err := p.expect(kwSELECT)
	if err != nil {
		t.Fatalf("expect(kwSELECT) error: %v", err)
	}
	if tok.Kind != kwSELECT {
		t.Errorf("expected kwSELECT, got %d", tok.Kind)
	}

	// expect wrong type
	_, err = p.expect(kwFROM)
	if err == nil {
		t.Error("expected error for wrong token type")
	}
}

func TestParserSkipToNextStatement(t *testing.T) {
	p := &Parser{
		lexer: NewLexer("SELECT 1; INSERT"),
	}
	p.advance()

	p.skipToNextStatement()
	// Should have skipped past SELECT 1 ; and now be at INSERT
	if p.cur.Kind != kwINSERT {
		t.Errorf("after skip, cur = %d (%s), want kwINSERT", p.cur.Kind, TokenName(p.cur.Kind))
	}
}

func TestParseAllDispatchCategories(t *testing.T) {
	// Verify each dispatch category produces the expected "not yet supported" error.
	tests := []struct {
		input   string
		wantMsg string
	}{
		{"CREATE TABLE t (id INT)", "CREATE"},
		{"ALTER TABLE t ADD COLUMN c INT", "ALTER"},
		{"DROP TABLE t", "DROP"},
		{"TRUNCATE TABLE t", "TRUNCATE"},
		{"INSERT INTO t VALUES (1)", "INSERT"},
		{"UPDATE t SET c=1", "UPDATE"},
		{"DELETE FROM t", "DELETE"},
		{"MERGE INTO t USING s ON t.id=s.id WHEN MATCHED THEN DELETE", "MERGE"},
		{"LOAD DATA INFILE 'f' INTO TABLE t", "LOAD"},
		{"EXPORT TABLE t", "EXPORT"},
		{"BEGIN", "BEGIN"},
		{"COMMIT", "COMMIT"},
		{"ROLLBACK", "ROLLBACK"},
		{"GRANT SELECT ON t TO u", "GRANT"},
		{"REVOKE SELECT ON t FROM u", "REVOKE"},
		{"SHOW TABLES", "SHOW"},
		{"DESCRIBE t", "DESCRIBE"},
		{"EXPLAIN SELECT 1", "EXPLAIN"},
		{"USE db", "USE"},
		{"SET x = 1", "SET"},
		{"UNSET x", "UNSET"},
		{"ADMIN SHOW REPLICA", "ADMIN"},
		{"KILL 123", "KILL"},
		{"BACKUP SNAPSHOT s", "BACKUP"},
		{"RESTORE SNAPSHOT s", "RESTORE"},
		{"RECOVER DATABASE db", "RECOVER"},
		{"REFRESH CATALOG c", "REFRESH"},
		{"CANCEL LOAD", "CANCEL"},
		{"ANALYZE TABLE t", "ANALYZE"},
		{"CLEAN ALL PROFILE", "CLEAN"},
	}
	for _, tt := range tests {
		_, errs := Parse(tt.input)
		if len(errs) == 0 {
			t.Errorf("Parse(%q): expected error", tt.input)
			continue
		}
		want := tt.wantMsg + " statement parsing is not yet supported"
		if errs[0].Msg != want {
			t.Errorf("Parse(%q): got %q, want %q", tt.input, errs[0].Msg, want)
		}
	}
}

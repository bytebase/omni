package pg

import "testing"

// TestParseRejectsUnterminatedConstructAtStatementStart: the public entry
// point must not swallow a lexer failure that begins a statement.
func TestParseRejectsUnterminatedConstructAtStatementStart(t *testing.T) {
	for _, sql := range []string{"'unterminated", "SELECT 1; /* unterminated", `SELECT 1; "abc`, "SELECT 1; $$abc"} {
		stmts, err := Parse(sql)
		if err == nil {
			t.Errorf("Parse(%q) returned %d statements and no error", sql, len(stmts))
		}
	}
}

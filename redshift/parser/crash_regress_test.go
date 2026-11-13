package parser

import (
	"errors"
	"testing"
)

// TestStatementsThatCrashedAreRejected pins nine statements the parser
// panicked or looped forever on, as pg/parser did. PostgreSQL 17 rejects
// each (pg/parser's TestRoutineBodyStatementsMatchPG and
// TestMissingNamesMatchPG hold that to the server).
func TestStatementsThatCrashedAreRejected(t *testing.T) {
	for _, sql := range []string{
		// routine_body_stmt_list: every statement ends in a semicolon.
		"CREATE FUNCTION a() BEGIN ATOMIC; zzz END",
		"CREATE FUNCTION a() BEGIN ATOMIC zzz END",
		"CREATE FUNCTION a() BEGIN ATOMIC zzz; END",
		"CREATE FUNCTION a() BEGIN ATOMIC a; END",
		"CREATE FUNCTION a() BEGIN ATOMIC RETURN a zzz; END",
		"CREATE FUNCTION a() BEGIN ATOMIC RETURNING a; END",
		// A name the grammar requires is missing.
		"ALTER EXTENSION a ADD OPERATOR CLASS USING a",
		"ALTER EXTENSION a ADD OPERATOR FAMILY USING a",
		"CREATE FUNCTION a(a.% TYPE)",
		"ALTER SYSTEM SET TO TRUE",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a *ParseError", sql, err)
		}
	}
	for _, sql := range []string{
		"CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 1; RETURN 2; END",
		"CREATE FUNCTION f() LANGUAGE sql BEGIN ATOMIC ; END",
		"ALTER EXTENSION a ADD OPERATOR CLASS b USING c",
		"CREATE FUNCTION a(b c.d%TYPE) RETURNS int LANGUAGE sql RETURN 1",
		"ALTER SYSTEM SET a TO TRUE",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

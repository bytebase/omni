package parser

import (
	"strings"
	"testing"
)

// KB-2d bug 2: DROP FUNCTION / PROCEDURE / ROUTINE / AGGREGATE / OPERATOR
// previously silently discarded the syntax error raised by parseFuncName /
// parseAnyOperator when the name was missing. As a result, inputs like
//
//	drop function ()
//	drop aggregate ()
//	drop routine ()
//	drop operator (int4, int4)
//
// parsed as DropStmt + orphan `()` — the `(`/`)` were left as residual tokens
// that Parse()'s top-level multi-stmt fallback split into a second bogus
// RawStmt. PG rejects each of these as `syntax error at or near "("`.
//
// These tests pin the parser to propagate the parseFuncName / parseAnyOperator
// error so that the whole DROP statement rejects with a syntax error rather
// than silently double-parsing.
func TestDropEmptyNameRejects(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		// expected substring of error.Error() so we don't pin the
		// *exact* message (the parser's error formatting can evolve).
		wantErrSub string
	}{
		{"drop function ()", "drop function ()", `"("`},
		{"drop procedure ()", "drop procedure ()", `"("`},
		{"drop routine ()", "drop routine ()", `"("`},
		{"drop aggregate ()", "drop aggregate ()", `"("`},
		{"drop operator ()", "drop operator ()", `"("`},
		{"drop operator (int4, int4)", "drop operator (int4, int4)", `"("`},
		{"drop aggregate (missing name)", "drop aggregate", "syntax error"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", tc.sql)
			}
			if !strings.Contains(err.Error(), tc.wantErrSub) {
				t.Errorf("Parse(%q): expected error containing %q, got %v",
					tc.sql, tc.wantErrSub, err)
			}
		})
	}
}

// Guard: the happy-path DROPs with explicit names still parse as single
// statements — we must not over-reject.
func TestDropWithNameStillParses(t *testing.T) {
	cases := []string{
		"drop function foo()",
		"drop function foo(int, text)",
		"drop procedure foo()",
		"drop routine foo()",
		"drop aggregate foo(int)",
		"drop aggregate foo(*)",
		"drop operator === (int4, int4)",
		"drop operator + (int, int)",
		// Multi-item lists still work.
		"drop function f1(), f2()",
		"drop aggregate a1(int), a2(text)",
		"drop operator === (int, int), !== (text, text)",
		// IF EXISTS variants.
		"drop function if exists foo()",
		"drop aggregate if exists foo(int)",
	}
	for _, sql := range cases {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
			if list == nil || len(list.Items) != 1 {
				t.Fatalf("Parse(%q): expected 1 stmt, got %d", sql, len(list.Items))
			}
		})
	}
}

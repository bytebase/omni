package parser

import (
	"testing"
)

// TestSnapshotRestoreIdentity verifies that snapshotTokenStream followed
// by restoreTokenStream rewinds the parser to a state where subsequent
// advance() calls produce the same token sequence as if no snapshot/
// restore had happened.
//
// This is the foundational test for the backtrack helper. If it fails,
// the helper is missing some piece of state that affects token emission,
// and every caller of snapshot/restore is unsound.
func TestSnapshotRestoreIdentity(t *testing.T) {
	corpora := []string{
		// Plain SQL with mixed token types
		`SELECT 1, 'a', $1, x.y, +1, -1 FROM t WHERE c = 'hello'`,
		// CREATE FUNCTION with multi-word type
		`CREATE FUNCTION f(p double precision) RETURNS int AS 'select 1' LANGUAGE sql`,
		// Qualified type
		`CREATE FUNCTION f(x pg_catalog.int4) RETURNS int AS 'select 1' LANGUAGE sql`,
		// Statements with comments and whitespace
		`SELECT /* comment */ 1, -- line comment
		2 FROM t`,
		// Dollar-quoted string
		`SELECT $body$hello world$body$ FROM t`,
	}

	for _, sql := range corpora {
		t.Run(truncate(sql, 40), func(t *testing.T) {
			// Parse 1: walk the entire token stream from scratch and record it
			expected := walkTokens(t, sql)

			// Parse 2: walk a few tokens, snapshot, walk a few more, restore,
			// then walk the rest. Compare to expected.
			p := newParserForTest(t, sql)
			var actual []Token

			// Walk 3 tokens
			for i := 0; i < 3 && p.cur.Type != lex_EOF; i++ {
				actual = append(actual, p.cur)
				p.advance()
			}

			// Snapshot
			snap := p.snapshotTokenStream()
			snapPos := len(actual)

			// Walk 5 more tokens (these will be discarded)
			for i := 0; i < 5 && p.cur.Type != lex_EOF; i++ {
				p.advance()
			}

			// Restore
			p.restoreTokenStream(snap)

			// Walk rest of stream
			for p.cur.Type != lex_EOF {
				actual = append(actual, p.cur)
				p.advance()
			}

			// Compare
			if len(actual) != len(expected) {
				t.Fatalf("token count mismatch after restore: got %d, want %d (snapshotted at index %d)",
					len(actual), len(expected), snapPos)
			}
			for i := range actual {
				if !tokensEqual(actual[i], expected[i]) {
					t.Errorf("token %d mismatch after restore: got %+v, want %+v", i, actual[i], expected[i])
				}
			}
		})
	}
}

// walkTokens parses sql and returns the full token sequence (excluding EOF).
func walkTokens(t *testing.T, sql string) []Token {
	t.Helper()
	p := newParserForTest(t, sql)
	var tokens []Token
	for p.cur.Type != lex_EOF {
		tokens = append(tokens, p.cur)
		p.advance()
	}
	return tokens
}

// newParserForTest constructs a Parser ready to walk sql. It mirrors
// what Parse() does internally, but exposes the Parser struct so tests
// can call snapshot/restore directly.
func newParserForTest(t *testing.T, sql string) *Parser {
	t.Helper()
	p := &Parser{
		lexer: NewLexer(sql),
	}
	p.advance() // prime cur with the first token
	return p
}

// tokensEqual compares two tokens for equality on Type and Str (the
// fields that affect parser behavior). Loc is excluded because the test
// only cares about token-stream equivalence, not position metadata.
func tokensEqual(a, b Token) bool {
	return a.Type == b.Type && a.Str == b.Str && a.Ival == b.Ival
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

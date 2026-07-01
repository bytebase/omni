package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// A lexer must never panic on malformed input; a truncated token that runs to
// EOF without its closing delimiter must surface a clean parse error. These are
// regression guards for the `slice bounds out of range` panic in
// skipWhitespaceAndComments on an unterminated /*! ... */ executable comment
// (and the adjacent unterminated-string / -identifier scan sites).

// parseNoPanic parses sql, converting any panic into a test failure, and returns
// the resulting error (nil on success).
func parseNoPanic(t *testing.T, sql string) (err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Parse panicked on %q: %v", sql, r)
		}
	}()
	_, err = Parse(sql)
	return err
}

func TestLexer_UnterminatedComment_NoPanic(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		// The confirmed real-world repro: versioned executable comment, no closing */.
		{"versioned_exec_comment", "CREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB /*!50100 PARTITION BY HASH (`id`)"},
		// Bare executable comment, no version, no closing */.
		{"bare_exec_comment", "CREATE TABLE `t` (`id` int) /*! PARTITION BY HASH (`id`)"},
		// Executable comment ending on a lone '*' at EOF (no trailing '/').
		{"exec_comment_star_at_eof", "SELECT 1 /*!50100 abc*"},
		// Plain (non-executable) block comment, no closing */.
		{"plain_block_comment", "CREATE TABLE `t` (`id` int) /* unterminated"},
		// Nested block comment where the inner one closes but the outer does not.
		{"nested_block_comment", "SELECT 1 /* outer /* inner */"},
		// Executable comment opened right at EOF.
		{"exec_comment_only", "/*!50100"},
		// Bare plain block comment as the whole input: Split must NOT drop it as an
		// empty segment — the lexer has to surface the error (regression guard).
		{"bare_block_comment_only", "/*"},
		{"bare_block_comment_ws", "   /*   "},
		// Trailing unterminated comment after a complete statement + delimiter: the
		// trailing segment must not be silently dropped by Split.
		{"trailing_after_stmt", "SELECT 1; /*"},
		{"trailing_exec_after_stmt", "SELECT 1; /*!50100"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseNoPanic(t, tc.sql)
			if err == nil {
				t.Fatalf("expected a parse error for unterminated comment, got nil (sql: %q)", tc.sql)
			}
			if !strings.Contains(err.Error(), "unterminated comment") {
				t.Errorf("error %q does not mention 'unterminated comment' (sql: %q)", err.Error(), tc.sql)
			}
		})
	}
}

func TestLexer_UnterminatedString_NoPanic(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{"single_quote", "SELECT 'abc"},
		{"double_quote", "SELECT \"abc"},
		// Trailing backslash at EOF: the escape has no following char and there is
		// no closing quote — still unterminated, must not read past the buffer.
		{"trailing_backslash", "SELECT 'abc\\"},
		{"empty_open_quote", "SELECT '"},
		{"string_in_ddl", "CREATE TABLE `t` (`id` int) COMMENT='oops"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseNoPanic(t, tc.sql)
			if err == nil {
				t.Fatalf("expected a parse error for unterminated string, got nil (sql: %q)", tc.sql)
			}
			if !strings.Contains(err.Error(), "unterminated string literal") {
				t.Errorf("error %q does not mention 'unterminated string literal' (sql: %q)", err.Error(), tc.sql)
			}
		})
	}
}

func TestLexer_UnterminatedBacktickIdent_NoPanic(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{"select_ident", "SELECT `abc"},
		{"empty_open_backtick", "SELECT `"},
		// Doubled backtick escape immediately before EOF — still unterminated.
		{"escaped_backtick_eof", "SELECT `a``"},
		{"ddl_table_name", "CREATE TABLE `t"},
		// Backtick-quoted user variable @`... scanned by scanVariable (a separate
		// scan site from scanBacktickIdent) — must also report untermination.
		{"backtick_user_variable", "SELECT @`abc"},
		{"backtick_user_variable_empty", "SELECT @`"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseNoPanic(t, tc.sql)
			if err == nil {
				t.Fatalf("expected a parse error for unterminated identifier, got nil (sql: %q)", tc.sql)
			}
			if !strings.Contains(err.Error(), "unterminated quoted identifier") {
				t.Errorf("error %q does not mention 'unterminated quoted identifier' (sql: %q)", err.Error(), tc.sql)
			}
		})
	}
}

// TestLexer_ValidCommentsAndDelimiters_StillParse is the regression counterpart:
// the bound-checks must not break well-formed comments, strings, or identifiers.
func TestLexer_ValidCommentsAndDelimiters_StillParse(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{"plain_block_comment", "SELECT 1 /* a normal comment */ + 2"},
		{"exec_comment", "SELECT 1 /*! + 2 */"},
		{"versioned_exec_comment", "SELECT 1 /*!50100 + 2 */"},
		{"versioned_partition", "CREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB /*!50100 PARTITION BY HASH (`id`) */"},
		{"nested_block_comment", "SELECT 1 /* outer /* inner */ still-outer */ + 2"},
		{"line_comment", "SELECT 1 -- trailing\n+ 2"},
		{"hash_comment", "SELECT 1 # trailing\n+ 2"},
		{"single_quote_string", "SELECT 'hello world'"},
		{"double_quote_string", "SELECT \"hello world\""},
		{"escaped_string", "SELECT 'it''s fine'"},
		{"backtick_ident", "SELECT `weird name`"},
		{"escaped_backtick_ident", "SELECT `a``b`"},
		{"backslash_escape_string", "SELECT 'line1\\nline2'"},
		{"backtick_user_variable", "SELECT @`v`"},
		// A valid, terminated executable comment followed by more SQL that also
		// happens to be valid — exercises the splice path without any error.
		{"exec_comment_then_string", "/*!  */ SELECT 'abc'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseNoPanic(t, tc.sql)
			if err != nil {
				t.Errorf("valid SQL failed to parse: %v (sql: %q)", err, tc.sql)
			}
		})
	}
}

// TestLexer_ErrorPosition_AfterSplice guards that a lexer error carries the correct
// position in ORIGINAL-input coordinates after an executable-comment splice shortens
// the buffer. Two distinct cases must both be right: an error INSIDE the retained
// comment body (shifted only by the leading "/*!NNNNN") and one AFTER the whole
// comment (shifted by both the leading run and the trailing "*/"). The `want` in each
// case is the byte offset of the opening single quote in the ORIGINAL string.
func TestLexer_ErrorPosition_AfterSplice(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		want int // expected ParseError.Position (0-based byte offset of the opening quote)
	}{
		// After the comment body.
		{"after_body_bare", "/*!  */ SELECT 'abc", 15},
		{"after_body_versioned", "/*!50100  */ SELECT 'x", 20},
		{"after_body_double_splice", "/*!  */ /*!  */ SELECT 'x", 23},
		// Inside the retained body of a terminated executable comment: the trailing
		// */ is removed but sits AFTER the error, so it must not shift the position.
		{"inside_body_bare", "/*! SELECT 'abc */", 11},
		{"inside_body_versioned", "/*!50100 SELECT 'x */", 16},
		{"inside_body_nested_comment", "/*! /* c */ SELECT 'x */", 19},
		// Nested EXECUTABLE comment: the inner /*! ... */ splices inside the outer
		// body, shortening the buffer to the left of the outer comment's trailing
		// */. The outer trailing gap must be re-mapped or the position drifts 2 bytes
		// too early (regression guard for the nested-splice coordinate-space bug).
		{"nested_exec_then_string", "/*! /*! SELECT 1 */ */'x", 22},
		{"nested_exec_then_string_ws", "/*! /*! SELECT 1 */ */ 'x", 23},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("expected an unterminated-string error for %q", tc.sql)
			}
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("expected *ParseError, got %T", err)
			}
			if pe.Position != tc.want {
				t.Errorf("%q: error Position = %d, want %d (the opening quote in the original SQL)", tc.sql, pe.Position, tc.want)
			}
		})
	}
}

// TestLexer_VersionedPartitionComment_ParsesPartition proves the executable-comment
// splice still exposes the inner SQL for a valid /*!50100 ... */ partition clause —
// i.e. the fix hardened the unterminated case without dropping the terminated one.
// The PARTITION BY wrapped in /*!50100 ... */ must be parsed as real SQL, so the
// resulting CREATE TABLE carries a partition clause.
func TestLexer_VersionedPartitionComment_ParsesPartition(t *testing.T) {
	withComment := "CREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB /*!50100 PARTITION BY HASH (`id`) PARTITIONS 4 */"
	got, err := Parse(withComment)
	if err != nil {
		t.Fatalf("versioned-comment partition failed to parse: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(got.Items))
	}
	ct, ok := got.Items[0].(*nodes.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *CreateTableStmt, got %T", got.Items[0])
	}
	if ct.Partitions == nil {
		t.Errorf("PARTITION BY inside /*!50100 ... */ was not parsed: CreateTableStmt.Partitions is nil")
	}
}

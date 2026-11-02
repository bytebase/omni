package parser

import (
	"strings"
	"testing"
)

// unterminatedCase is a statement whose last token runs to end of input without
// its closing delimiter, and the message the lexer must report for it.
type unterminatedCase struct {
	name string
	sql  string
	want string
}

// unterminatedCases covers the tokens the lexer scans to a closing delimiter
// and that TiDB rejects when they run off the end of the input. A ';' inside the
// open token is deliberate: Split must not read it as a statement boundary and
// hand the parser a segment that looks complete. An executable comment is not
// here, because TiDB drops an opener that never closes: see unclosedExecCases.
var unterminatedCases = []unterminatedCase{
	{"string literal", "SELECT 'oops;", "unterminated string literal"},
	{"double-quoted string", `SELECT "oops;`, "unterminated string literal"},
	{"charset-prefixed string", "SELECT _utf8'oops;", "unterminated string literal"},
	{"quoted identifier", "SELECT `oops;", "unterminated quoted identifier"},
	{"quoted variable name", "SELECT @`oops;", "unterminated quoted identifier"},
	{"block comment", "SELECT 1 /* oops;", "unterminated comment"},
	{"block comment alone", "/* oops", "unterminated comment"},
	// A feature TiDB does not know makes the comment an ordinary one, which must
	// still close.
	{"tidb unsupported feature", "SELECT 1 /*T![nonexistent_feature] oops", "unterminated comment"},
	// A feature list TiDB cannot read leaves "/*T!" an executable comment whose
	// body starts at the "[", and TiDB v8.5.5 rejects it there (column 14).
	{"tidb feature bracket", "SELECT 1 /*T![oops", `syntax error at or near "["`},
}

// unclosedExecCases are executable comments that never close. TiDB v8.5.5 drops
// the opener and reads the rest of the input as ordinary SQL, so each of these
// parses as the statement left behind: `SELECT 1 /*T! oops` is `SELECT 1 oops`,
// an aliased column. MySQL and MariaDB reject every one as an unterminated
// comment, which is why only this fork carries the list.
var unclosedExecCases = []string{
	"SELECT 1 /*!40000 oops",
	"SELECT 1 /*T! oops",
	"SELECT 1 /*T![clustered_index] oops",
	"/*T! SELECT 1",
	"SELECT /*T! 1",
	"SELECT /*T![ttl] 1",
	"SELECT /*T![auto_rand] col",
	"SELECT 1 /*! ",
}

// terminatedCounterparts are the same statements with their delimiters closed.
// They must still parse, so the guards reject only what runs off the end of the
// input. A comment whose feature TiDB does not support is skipped whole, so its
// body need not be SQL.
var terminatedCounterparts = []string{
	"SELECT 'oops';",
	`SELECT "oops";`,
	"SELECT _utf8'oops';",
	"SELECT `oops`;",
	"SELECT @`oops`;",
	"SELECT 1 /* oops */;",
	"/* oops */",
	"SELECT 1 /*!40000 + 1 */",
	"SELECT 1 /*T! + 1 */",
	"SELECT 1 /*T![clustered_index] + 1 */",
	"SELECT 1 /*T![nonexistent_feature] not sql at all */",
}

// TestUnterminatedLiteralRejected pins that an unterminated token is a parse
// error naming what was left open, rather than a half-open token the parser
// accepts. Scanning stops at EOF and records the error, so a statement that
// parsed before the bad token still fails.
func TestUnterminatedLiteralRejected(t *testing.T) {
	for _, tc := range unterminatedCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) accepted an unterminated token", tc.sql)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Parse(%q) error = %q, want it to mention %q", tc.sql, err, tc.want)
			}
		})
	}
}

func TestTerminatedCounterpartsAccepted(t *testing.T) {
	for _, sql := range terminatedCounterparts {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

// TestUnclosedExecCommentDropsTheOpener pins TiDB's handling of an executable
// comment that never closes: the opener is dropped and the rest of the input is
// read as ordinary SQL, so the statement stands or falls on what is left.
func TestUnclosedExecCommentDropsTheOpener(t *testing.T) {
	for _, sql := range unclosedExecCases {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

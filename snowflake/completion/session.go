package completion

import (
	"cmp"
	"context"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/snowflake/parser"
)

// SessionConfig configures session-aware Snowflake completion.
type SessionConfig struct {
	// DefaultDatabase is the current database, passed to LoadDatabase.
	DefaultDatabase string
	// DefaultSchema is the current schema, which resolves names with no
	// schema. PUBLIC when empty, as in a Snowflake session.
	DefaultSchema string
	// LoadDatabase fetches a database's metadata: DefaultDatabase's, and
	// another database's when a name in the statement is qualified by it.
	// When nil or failing, completion runs without catalog objects and
	// returns keywords, built-in functions and statement-local names.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
}

// CompleteSession completes the statement of a script at the 1-based line
// and 0-based UTF-16 code-unit caretOffset, as an editor reports a caret.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	cat := newCatalog(ctx, cfg)
	stmt, pos := statementAtCaret(statement, lineOffsetToBytePos(statement, caretLine, caretOffset))
	got := complete(stmt, pos, cat)

	type key struct {
		typ  CandidateType
		text string
	}
	seen := make(map[key]bool, len(got))
	result := make([]Candidate, 0, len(got))
	for _, c := range got {
		k := key{c.Type, c.Text}
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, c)
	}
	slices.SortFunc(result, func(a, b Candidate) int {
		return cmp.Or(cmp.Compare(a.Type, b.Type), cmp.Compare(a.Text, b.Text))
	})
	return result
}

// bareIdentifier is an identifier Snowflake reads back unchanged without
// quotes: unquoted identifiers resolve upper case.
var bareIdentifier = regexp.MustCompile(`^[A-Z_][A-Z0-9_$]*$`)

// QuoteIdentifierIfNeeded double-quotes name unless Snowflake reads it back
// bare: an unquoted identifier resolves upper case, so a name with lower
// case or special characters, or a reserved keyword, needs quotes. Inside an
// open double quote only embedded quotes are doubled.
func QuoteIdentifierIfNeeded(name string, caretInQuotes bool) string {
	if caretInQuotes {
		return strings.ReplaceAll(name, `"`, `""`)
	}
	if bareIdentifier.MatchString(name) && !parser.IsReservedKeyword(name) {
		return name
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// lineOffsetToBytePos converts a 1-based line and 0-based UTF-16 code-unit
// column into a byte offset into s.
func lineOffsetToBytePos(s string, line, offset int) int {
	currentLine := 1
	col := 0
	for i := 0; i < len(s); {
		if currentLine == line && col >= offset {
			return i
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\n' {
			if currentLine == line {
				return i
			}
			currentLine++
			col = 0
		} else if r <= 0xFFFF {
			col++
		} else {
			col += 2
		}
		i += size
	}
	return len(s)
}

// statementAtCaret isolates the statement holding bytePos, so neither
// earlier statements nor broken later ones affect completion or its cost.
// Split leaves each semicolon between two segments, so a caret right after
// one is in the next statement, and a caret past the last one, where Split
// may return no segment, starts a statement.
func statementAtCaret(script string, bytePos int) (string, int) {
	bytePos = max(0, min(bytePos, len(script)))
	for _, seg := range parser.Split(script) {
		if bytePos >= seg.ByteStart && bytePos <= seg.ByteEnd {
			return seg.Text, bytePos - seg.ByteStart
		}
	}
	return "", 0
}

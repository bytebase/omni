package completion

import (
	"cmp"
	"context"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/starrocks/parser"
)

// SessionConfig configures session-aware StarRocks completion.
type SessionConfig struct {
	// DefaultDatabase is the current database, which resolves names with no
	// database.
	DefaultDatabase string
	// LoadDatabase fetches a database's metadata: DefaultDatabase's, and
	// another database's when a name in the statement is qualified by it.
	// When nil or failing, completion offers keywords only.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
	// ListDatabaseNames lists the instance's databases, offered where a
	// table name fits. Optional.
	ListDatabaseNames func(ctx context.Context) ([]string, error)
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

// bareIdentifier is an identifier StarRocks's lexer reads as one unquoted word.
var bareIdentifier = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

// QuoteIdentifierIfNeeded backtick-quotes name unless StarRocks reads it back
// bare: a word of identifier characters that is not a reserved keyword.
// Inside an open backtick only embedded backticks are doubled.
func QuoteIdentifierIfNeeded(name string, caretInBackticks bool) string {
	if caretInBackticks {
		return strings.ReplaceAll(name, "`", "``")
	}
	if bareIdentifier.MatchString(name) {
		if kind, ok := parser.KeywordToken(name); !ok || !parser.IsReserved(kind) {
			return name
		}
	}
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
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

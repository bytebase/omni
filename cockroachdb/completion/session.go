package completion

import (
	"bytes"
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
)

// SessionConfig configures session-aware CockroachDB completion.
type SessionConfig struct {
	// DefaultDatabase is the active database name passed to LoadDatabase.
	DefaultDatabase string
	// DefaultSchema resolves unqualified names before "public" when
	// non-empty.
	DefaultSchema string
	// LoadDatabase fetches metadata for DefaultDatabase. When nil or failing,
	// completion runs without catalog objects and returns only keywords and
	// statement-local names such as CTEs.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
}

// CompleteSession completes the statement of a script at the 1-based line
// and 0-based UTF-16 code-unit caretOffset, as an editor reports a caret.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	var cat *catalog
	if cfg.LoadDatabase != nil && cfg.DefaultDatabase != "" {
		if meta, err := cfg.LoadDatabase(ctx, cfg.DefaultDatabase); err == nil && meta != nil {
			cat = newCatalog(meta, cfg.DefaultSchema)
		}
	}

	stmt, pos := statementAtCaret(statement, lineOffsetToBytePos(statement, caretLine, caretOffset))
	got := complete(ctx, stmt, pos, cat)

	type key struct {
		typ  CandidateType
		text string
	}
	seen := make(map[key]bool, len(got.candidates))
	result := make([]Candidate, 0, len(got.candidates))
	for _, c := range got.candidates {
		c.Text = formatCandidateText(c, got.quoted, got.call)
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

// QuoteIdentifierIfNeeded double-quotes name when CockroachDB would not read
// it back bare: it has upper case or special characters, or is a reserved
// keyword. Inside an open double quote only embedded quotes are doubled.
func QuoteIdentifierIfNeeded(name string, caretInQuotes bool) string {
	if caretInQuotes {
		return strings.ReplaceAll(name, `"`, `""`)
	}
	var buf bytes.Buffer
	lexbase.EncodeRestrictedSQLIdent(&buf, name, lexbase.EncNoFlags)
	return buf.String()
}

func formatCandidateText(c Candidate, quoted, call bool) string {
	switch c.Type {
	case CandidateKeyword:
		return c.Text
	case CandidateFunction:
		text := QuoteIdentifierIfNeeded(c.Text, quoted)
		if call && !quoted {
			text += "()"
		}
		return text
	default:
		return QuoteIdentifierIfNeeded(c.Text, quoted)
	}
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
// earlier statements nor broken later ones affect completion or its cost. A
// caret right after a statement's semicolon starts the next statement.
func statementAtCaret(script string, bytePos int) (string, int) {
	bytePos = max(0, min(bytePos, len(script)))
	start := 0
	for {
		n, ok := splitFirst(script[start:])
		if !ok {
			return script[start:], bytePos - start
		}
		end := start + n
		if bytePos < end {
			return script[start:end], bytePos - start
		}
		start = end
		if bytePos == end && strings.TrimSpace(script[end:]) == "" {
			return "", 0
		}
	}
}

// splitFirst returns the length of the first statement of script through
// its semicolon, or false when no semicolon ends one.
func splitFirst(script string) (n int, ok bool) {
	defer func() {
		if recover() != nil {
			n, ok = 0, false
		}
	}()
	return crparser.SplitFirstStatement(script)
}

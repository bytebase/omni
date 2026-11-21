package completion

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	omniredshift "github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/catalog"
	"github.com/bytebase/omni/redshift/parser"
)

// SessionConfig configures session-aware Redshift completion.
type SessionConfig struct {
	// DefaultDatabase is the active database name passed to LoadDatabase.
	DefaultDatabase string
	// DefaultSchema is prepended to the catalog search_path before "public"
	// when non-empty.
	DefaultSchema string
	// LoadDatabase fetches metadata for DefaultDatabase. When nil or failing,
	// completion runs without catalog objects and returns only keywords and
	// statement-local scopes.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
}

// CompleteSession runs session-aware Redshift auto-completion over statement
// at the 1-based line and 0-based UTF-16 code-unit caretOffset.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	cat := buildSessionCatalog(ctx, cfg)

	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	stmt, pos := statementAtCaret(statement, bytePos)
	caretInQuotes := CaretInsideDoubleQuotes(stmt, pos)
	raw := Complete(stmt, pos, cat)

	type candKey struct {
		typ        CandidateType
		text       string
		definition string
	}
	seen := make(map[candKey]bool, len(raw))
	result := make([]Candidate, 0, len(raw))
	for _, c := range raw {
		candidate := Candidate{
			Text:       formatSessionCandidateText(c, caretInQuotes),
			Type:       c.Type,
			Definition: c.Definition,
			Comment:    c.Comment,
		}
		k := candKey{typ: candidate.Type, text: candidate.Text, definition: candidate.Definition}
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, candidate)
	}
	slices.SortFunc(result, compareSessionCandidates)
	return result
}

// CaretInsideDoubleQuotes reports whether byteOffset in sql sits inside an
// unclosed double-quoted identifier, respecting Redshift comments, string
// literals, and dollar-quoted strings.
func CaretInsideDoubleQuotes(sql string, byteOffset int) bool {
	return insideDoubleQuotes(sql, byteOffset)
}

func insideDoubleQuotes(sql string, offset int) bool {
	if offset <= 0 {
		return false
	}
	if offset > len(sql) {
		offset = len(sql)
	}
	lex := parser.NewLexerRange(sql, 0, offset)
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
	}
	return lex.Err != nil && lex.Err.Error() == "unterminated quoted identifier"
}

// QuoteIdentifierIfNeeded double-quotes a Redshift identifier when it is not a
// valid lowercase unquoted identifier or collides with a reserved keyword.
// When caretInQuotes is true, outer delimiters are omitted while embedded
// double quotes are still doubled.
func QuoteIdentifierIfNeeded(name string, caretInQuotes bool) string {
	if caretInQuotes {
		return strings.ReplaceAll(name, `"`, `""`)
	}
	if strings.ToLower(name) != name || parser.IsReservedKeyword(strings.ToUpper(name)) || !isValidUnquotedIdentifier(name) {
		return quoteIdent(name)
	}
	return name
}

// LineOffsetToBytePos converts a 1-based line and 0-based UTF-16 code-unit
// column offset into a byte offset into s, matching LSP/Monaco positions.
func LineOffsetToBytePos(s string, line, offset int) int {
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

// statementAtCaret isolates the single statement containing bytePos so neither
// earlier statements nor broken trailing statements affect completion or cost.
func statementAtCaret(statement string, bytePos int) (string, int) {
	bytePos = max(0, min(bytePos, len(statement)))
	segs := omniredshift.Split(statement)
	if len(segs) == 0 {
		return statement, bytePos
	}
	for i, seg := range segs {
		if bytePos < seg.ByteStart {
			continue
		}
		if bytePos == seg.ByteEnd {
			if i+1 < len(segs) {
				continue
			}
			if strings.HasSuffix(seg.Text, ";") {
				return "", 0
			}
		}
		if bytePos <= seg.ByteEnd {
			return seg.Text, bytePos - seg.ByteStart
		}
	}
	last := segs[len(segs)-1]
	return last.Text, bytePos - last.ByteStart
}

func buildSessionCatalog(ctx context.Context, cfg SessionConfig) *catalog.Catalog {
	if cfg.LoadDatabase == nil || cfg.DefaultDatabase == "" {
		return nil
	}
	protoMeta, err := cfg.LoadDatabase(ctx, cfg.DefaultDatabase)
	if err != nil || protoMeta == nil {
		return nil
	}

	cat := catalog.New()
	cat.LoadMetadata(protoMeta)

	searchPath := []string{"public"}
	if cfg.DefaultSchema != "" && cfg.DefaultSchema != "public" {
		searchPath = []string{cfg.DefaultSchema, "public"}
	}
	cat.SetSearchPath(searchPath)
	return cat
}

func formatSessionCandidateText(c Candidate, caretInQuotes bool) string {
	switch c.Type {
	case CandidateKeyword, CandidateType_:
		return c.Text
	default:
		return QuoteIdentifierIfNeeded(c.Text, caretInQuotes)
	}
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func isValidUnquotedIdentifier(s string) bool {
	if s == "" || !utf8.ValidString(s) {
		return false
	}
	first, size := utf8.DecodeRuneInString(s)
	if !unicode.IsLetter(first) && first != '_' {
		return false
	}
	for _, ch := range s[size:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}
	return true
}

func compareSessionCandidates(a, b Candidate) int {
	return cmp.Or(
		cmp.Compare(a.Type, b.Type),
		cmp.Compare(a.Text, b.Text),
		cmp.Compare(a.Definition, b.Definition),
		cmp.Compare(a.Comment, b.Comment),
	)
}

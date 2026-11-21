package completion

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/mssql/parser"
)

// SessionConfig configures multi-database, session-aware MSSQL completion.
type SessionConfig struct {
	// DefaultDatabase is the active session database, if any.
	DefaultDatabase string
	// DefaultSchema is the active session schema (catalog defaults to "dbo" when empty).
	DefaultSchema string
	// LoadDatabase fetches metadata for a database name on demand.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
	// ListDatabaseNames enumerates known database names on the instance.
	ListDatabaseNames func(ctx context.Context) ([]string, error)
}

// CompleteSession runs session-aware MSSQL auto-completion over statement at
// the 1-based line and 0-based UTF-16 code-unit caretOffset.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	stmt, pos := isolateStatementAtCursor(statement, bytePos)

	cat := buildSessionCatalog(ctx, cfg, stmt)
	caretQuoted := CaretTokenIsQuoted(stmt, pos)
	raw := Complete(stmt, pos, cat)

	type candKey struct {
		typ        CandidateType
		text       string
		definition string
	}
	seen := make(map[candKey]Candidate, len(raw))
	for _, c := range raw {
		candidate := Candidate{
			Type:       c.Type,
			Text:       formatSessionCandidateText(c, caretQuoted),
			Definition: c.Definition,
			Comment:    c.Comment,
		}
		k := candKey{typ: candidate.Type, text: candidate.Text, definition: candidate.Definition}
		if _, exists := seen[k]; !exists {
			seen[k] = candidate
		}
	}

	result := make([]Candidate, 0, len(seen))
	for _, c := range seen {
		result = append(result, c)
	}
	slices.SortFunc(result, compareSessionCandidates)
	return result
}

// CaretTokenIsQuoted reports whether cursorByteOffset in sql sits inside an
// open bracketed ([...]) or double-quoted ("...") T-SQL identifier.
func CaretTokenIsQuoted(sql string, cursorByteOffset int) bool {
	if cursorByteOffset <= 0 {
		return false
	}
	if cursorByteOffset > len(sql) {
		cursorByteOffset = len(sql)
	}
	lex := parser.NewLexer(sql[:cursorByteOffset])
	var lastErr error
	for {
		lex.Err = nil
		tok := lex.NextToken()
		if lex.Err != nil {
			lastErr = lex.Err
		} else if tok.Type != 0 {
			lastErr = nil
		}
		if tok.Type == 0 {
			break
		}
	}
	if lastErr == nil {
		return false
	}
	msg := lastErr.Error()
	return msg == "unterminated bracketed identifier" || msg == "unterminated quoted identifier"
}

// QuoteIdentifierIfNeeded bracket-quotes a T-SQL identifier when it is not a
// valid regular identifier or collides with a reserved keyword. When
// caretQuoted is true, outer delimiters are omitted.
func QuoteIdentifierIfNeeded(name string, caretQuoted bool) string {
	if caretQuoted || isRegularIdentifier(name) {
		return name
	}
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
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

func buildSessionCatalog(ctx context.Context, cfg SessionConfig, statement string) *catalog.Catalog {
	if cfg.LoadDatabase == nil && cfg.ListDatabaseNames == nil {
		return nil
	}

	var allNames []string
	if cfg.ListDatabaseNames != nil {
		if names, err := cfg.ListDatabaseNames(ctx); err == nil {
			allNames = names
		}
	}

	defaultDB := cfg.DefaultDatabase
	for _, name := range allNames {
		if strings.EqualFold(name, defaultDB) {
			defaultDB = name
			break
		}
	}

	cat := catalog.New()
	if defaultDB != "" {
		loadDatabaseIntoCatalog(ctx, cfg, cat, defaultDB, true)
	}
	if cfg.DefaultSchema != "" {
		cat.SetDefaultSchema(cfg.DefaultSchema)
	}

	lowerStmt := strings.ToLower(statement)
	for _, name := range allNames {
		if name == "" || strings.EqualFold(name, defaultDB) {
			continue
		}
		// Names containing ] or " are escaped as ]] / "" when quoted, so a raw
		// substring check can miss them; load those unconditionally.
		needFullLoad := strings.Contains(lowerStmt, strings.ToLower(name)) ||
			strings.ContainsAny(name, `]"`)
		loadDatabaseIntoCatalog(ctx, cfg, cat, name, needFullLoad)
	}
	return cat
}

func loadDatabaseIntoCatalog(ctx context.Context, cfg SessionConfig, cat *catalog.Catalog, dbName string, full bool) {
	if full && cfg.LoadDatabase != nil {
		if meta, err := cfg.LoadDatabase(ctx, dbName); err == nil && meta != nil {
			_ = cat.LoadMetadataInto(ctx, meta, dbName)
			return
		}
	}
	_ = cat.LoadMetadataInto(ctx, nil, dbName)
}

func formatSessionCandidateText(c Candidate, caretQuoted bool) string {
	switch {
	case c.Type == CandidateKeyword || c.Type == CandidateType_:
		return c.Text
	case c.Type == CandidateFunction:
		if strings.HasSuffix(c.Text, "()") {
			return c.Text
		}
		return c.Text + "()"
	case c.Text == "$action":
		return c.Text
	default:
		return QuoteIdentifierIfNeeded(c.Text, caretQuoted)
	}
}

func isRegularIdentifier(identifier string) bool {
	if identifier == "" || !utf8.ValidString(identifier) {
		return false
	}

	first, size := utf8.DecodeRuneInString(identifier)
	if !unicode.IsLetter(first) && first != '_' && first != '@' && first != '#' {
		return false
	}
	for _, r := range identifier[size:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '@' && r != '$' && r != '#' && r != '_' {
			return false
		}
	}
	if parser.IsReservedKeyword(identifier) {
		return false
	}
	for _, r := range identifier {
		if r == ' ' || !unicode.IsPrint(r) {
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

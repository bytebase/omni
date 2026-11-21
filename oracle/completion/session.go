package completion

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/parser"
)

// SessionConfig configures session-aware Oracle SQL completion.
type SessionConfig struct {
	// ReadOnly restricts completion to read-only query statements (SELECT,
	// WITH, EXPLAIN). At statement start only query-initiating keywords are
	// offered; inside non-query statements completion returns nil.
	ReadOnly bool
	// DefaultSchema is the active Oracle schema used for unqualified object
	// resolution.
	DefaultSchema string
	// LoadSchema lazily fetches metadata for an Oracle schema name on demand.
	// Results are cached for the duration of a CompleteSession call.
	LoadSchema func(ctx context.Context, schema string) (*metadatapb.DatabaseSchemaMetadata, error)
	// ListSchemaNames enumerates known Oracle schema names.
	ListSchemaNames func(ctx context.Context) ([]string, error)
}

var oracleQuerySceneStartKeywords = map[string]bool{
	"EXPLAIN": true,
	"SELECT":  true,
	"WITH":    true,
}

// CompleteSession runs session-aware Oracle auto-completion over statement at
// the 1-based line and 0-based UTF-16 code-unit caretOffset.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	sql, pos := isolateStatementAtCursor(statement, bytePos)

	var tokens []parser.Token
	if cfg.ReadOnly {
		tokens = parser.Tokenize(sql)
		if querySceneDisallowsCompletion(tokens, pos) {
			return nil
		}
	}

	opts := buildSessionOptions(ctx, cfg)
	raw := Complete(sql, pos, nil, opts...)
	caretInQuotes := CaretInsideDoubleQuotes(sql, pos)

	type candKey struct {
		typ        CandidateType
		text       string
		definition string
	}
	seen := make(map[candKey]Candidate, len(raw))
	for _, c := range raw {
		candidate := Candidate{
			Type:       c.Type,
			Text:       formatSessionCandidateText(c, caretInQuotes),
			Definition: c.Definition,
			Comment:    c.Comment,
		}
		k := candKey{typ: candidate.Type, text: candidate.Text, definition: candidate.Definition}
		if _, exists := seen[k]; !exists {
			seen[k] = candidate
		}
	}

	result := make([]Candidate, 0, len(seen))
	for _, candidate := range seen {
		result = append(result, candidate)
	}
	slices.SortFunc(result, compareSessionCandidates)

	if cfg.ReadOnly {
		result = filterCandidatesByReadOnly(sql, tokens, pos, result)
	}
	return result
}

// CaretInsideDoubleQuotes reports whether byteOffset in sql sits inside an
// unclosed double-quoted identifier, respecting Oracle q/nq-quoted string
// literals, single-quoted strings, and SQL comments.
func CaretInsideDoubleQuotes(sql string, byteOffset int) bool {
	_, _, inDQuote := scanSQLState(sql, byteOffset)
	return inDQuote
}

// QuoteIdentifierIfNeeded double-quotes an Oracle identifier when it would
// otherwise be invalid or case-folded as a bare identifier. When caretInQuotes
// is true, outer delimiters are omitted while embedded double quotes are still
// doubled.
func QuoteIdentifierIfNeeded(name string, caretInQuotes bool) string {
	if caretInQuotes {
		return strings.ReplaceAll(name, `"`, `""`)
	}
	return quoteIdentifierForSQL(name)
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

func buildSessionOptions(ctx context.Context, cfg SessionConfig) []Option {
	var opts []Option
	if cfg.DefaultSchema != "" {
		opts = append(opts, WithDefaultSchema(cfg.DefaultSchema))
	}
	if cfg.ListSchemaNames != nil {
		if names, err := cfg.ListSchemaNames(ctx); err == nil && len(names) > 0 {
			opts = append(opts, WithSchemaNames(names))
		}
	}
	if cfg.LoadSchema != nil {
		cache := make(map[string]*metadatapb.DatabaseSchemaMetadata)
		loadSchema := func(schema string) *metadatapb.DatabaseSchemaMetadata {
			if schema == "" {
				return nil
			}
			if cached, ok := cache[schema]; ok {
				return cached
			}
			meta, err := cfg.LoadSchema(ctx, schema)
			if err != nil || meta == nil {
				cache[schema] = nil
				return nil
			}
			cache[schema] = meta
			return meta
		}
		if cfg.DefaultSchema != "" {
			if defaultMeta := loadSchema(cfg.DefaultSchema); defaultMeta != nil {
				opts = append(opts, WithSchema(cfg.DefaultSchema, defaultMeta))
			}
		}
		opts = append(opts, WithSchemaLoader(loadSchema))
	}
	return opts
}

var oracleBuiltinDataTypeSet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(oracleDataTypes))
	for _, typ := range oracleDataTypes {
		set[typ] = struct{}{}
	}
	return set
}()

func formatSessionCandidateText(c Candidate, caretInQuotes bool) string {
	switch c.Type {
	case CandidateKeyword:
		return c.Text
	case CandidateType_:
		if !caretInQuotes {
			if _, builtin := oracleBuiltinDataTypeSet[c.Text]; builtin {
				return c.Text
			}
		}
		return QuoteIdentifierIfNeeded(c.Text, caretInQuotes)
	case CandidateFunction:
		name := strings.TrimSuffix(c.Text, "()")
		if caretInQuotes {
			return QuoteIdentifierIfNeeded(name, true)
		}
		if !oracleBuiltinFunctionKeywords[strings.ToUpper(name)] {
			name = QuoteIdentifierIfNeeded(name, false)
		}
		return name + "()"
	default:
		return QuoteIdentifierIfNeeded(c.Text, caretInQuotes)
	}
}

func compareSessionCandidates(a, b Candidate) int {
	return cmp.Or(
		cmp.Compare(a.Type, b.Type),
		cmp.Compare(a.Text, b.Text),
		cmp.Compare(a.Definition, b.Definition),
		cmp.Compare(a.Comment, b.Comment),
	)
}

func filterCandidatesByReadOnly(sql string, tokens []parser.Token, byteOffset int, candidates []Candidate) []Candidate {
	if !isAtStatementStartCompletion(sql, tokens, byteOffset) {
		return candidates
	}
	result := make([]Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if (candidate.Type == CandidateKeyword || candidate.Type == CandidateType_) &&
			!oracleQuerySceneStartKeywords[strings.ToUpper(candidate.Text)] {
			continue
		}
		result = append(result, candidate)
	}
	return result
}

func querySceneDisallowsCompletion(tokens []parser.Token, byteOffset int) bool {
	idx := currentStatementMainTokenIndex(tokens, byteOffset)
	if idx < 0 || tokens[idx].End >= byteOffset {
		return false
	}
	if strings.EqualFold(tokens[idx].Str, "EXPLAIN") {
		return false
	}
	switch tokens[idx].Type {
	case parser.SELECT, parser.WITH:
		return false
	default:
		return true
	}
}

func isAtStatementStartCompletion(sql string, tokens []parser.Token, byteOffset int) bool {
	collectOffset := max(0, min(byteOffset, len(sql)))
	for collectOffset > 0 && isOracleIdentByte(sql[collectOffset-1]) {
		collectOffset--
	}
	start := currentStatementStartTokenIndex(tokens, collectOffset)
	if start >= len(tokens) {
		return true
	}
	token := tokens[start]
	return token.Type == ';' || token.Loc >= collectOffset
}

func currentStatementStartTokenIndex(tokens []parser.Token, offset int) int {
	start := 0
	for i, token := range tokens {
		if token.Loc >= offset {
			break
		}
		if token.Type == ';' {
			start = i + 1
		}
	}
	return start
}

func currentStatementMainTokenIndex(tokens []parser.Token, byteOffset int) int {
	start := currentStatementStartTokenIndex(tokens, byteOffset)
	if start >= len(tokens) {
		return -1
	}
	token := tokens[start]
	if token.Loc >= byteOffset || token.Type == ';' {
		return -1
	}
	if token.Type == parser.WITH {
		depth := 0
		for i := start + 1; i < len(tokens); i++ {
			t := tokens[i]
			if t.Loc >= byteOffset || t.Type == ';' {
				break
			}
			switch t.Type {
			case '(':
				depth++
			case ')':
				if depth > 0 {
					depth--
				}
			default:
				if depth == 0 && (t.Type == parser.SELECT || t.Type == parser.INSERT ||
					t.Type == parser.UPDATE || t.Type == parser.DELETE || t.Type == parser.MERGE) {
					return i
				}
			}
		}
	}
	return start
}

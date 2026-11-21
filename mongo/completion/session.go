package completion

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mongo/catalog"
	"github.com/bytebase/omni/mongo/parser"
)

// SessionConfig configures session-aware MongoDB shell completion.
type SessionConfig struct {
	// DefaultDatabase is the active database name passed to LoadDatabase.
	DefaultDatabase string
	// LoadDatabase fetches metadata for DefaultDatabase. When nil or failing,
	// completion runs without catalog collections and returns only built-in
	// keywords, methods, and operators.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
}

// CompleteSession runs session-aware MongoDB shell auto-completion over
// statement at the 1-based line and 0-based UTF-16 code-unit caretOffset.
//
// Collection candidates are formatted for the active access syntax:
//   - dot access (`db.|`): valid JavaScript identifiers stay bare; others use `["name"]`.
//   - bracket access (`db[|`): formatted as `"name"]`.
//   - inside an open bracket quote (`db["|` or `db['|`): escaped content without outer delimiters.
//
// Results are deduplicated and sorted deterministically.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	cat := buildSessionCatalog(ctx, cfg)
	bytePos := LineOffsetToBytePos(statement, caretLine, caretOffset)
	raw := Complete(statement, bytePos, cat)
	inBracket, quote := bracketContextAt(statement, bytePos)

	type candKey struct {
		typ        CandidateType
		text       string
		definition string
	}
	seen := make(map[candKey]bool, len(raw))
	result := make([]Candidate, 0, len(raw))
	for _, c := range raw {
		text := c.Text
		if c.Type == CandidateCollection {
			text = formatCollectionCandidate(c.Text, inBracket, quote)
		}
		candidate := Candidate{
			Text:       text,
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

func buildSessionCatalog(ctx context.Context, cfg SessionConfig) *catalog.Catalog {
	cat := catalog.New()
	if cfg.DefaultDatabase == "" || cfg.LoadDatabase == nil {
		return cat
	}
	meta, err := cfg.LoadDatabase(ctx, cfg.DefaultDatabase)
	if err != nil || meta == nil {
		return cat
	}
	cat.LoadMetadata(meta)
	return cat
}

// bracketContextAt reports whether the completion context at cursorOffset is
// bracket collection access (`db[...`) and, if the user has already opened a
// string quote inside the brackets (`db["...` or `db['...`), which quote char.
func bracketContextAt(input string, cursorOffset int) (inBracket bool, quote byte) {
	cursorOffset = max(0, min(cursorOffset, len(input)))
	prefix := extractPrefix(input, cursorOffset)
	tokens := tokenize(input, cursorOffset-len(prefix))
	if detectContext(tokens) != contextAfterBracket {
		return false, 0
	}
	last := tokens[len(tokens)-1]
	if last.Type == parser.TokString && last.Loc >= 0 && last.Loc < len(input) {
		if q := input[last.Loc]; q == '"' || q == '\'' {
			return true, q
		}
	}
	return true, 0
}

func formatCollectionCandidate(name string, inBracket bool, quote byte) string {
	if inBracket {
		if quote == '"' || quote == '\'' {
			return escapeJSString(name, quote)
		}
		return `"` + escapeJSString(name, '"') + `"]`
	}
	if !isValidJSIdentifier(name) {
		return `["` + escapeJSString(name, '"') + `"]`
	}
	return name
}

func escapeJSString(s string, quote byte) string {
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	q := string(quote)
	return strings.ReplaceAll(escaped, q, `\`+q)
}

func isValidJSIdentifier(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		isStart := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '$'
		if i == 0 {
			if !isStart {
				return false
			}
		} else if !isStart && (c < '0' || c > '9') {
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

package completion

import (
	"context"
	"strings"

	"github.com/bytebase/omni/pg/yacc"
)

// Complete performs SQL auto-completion at the cursor position.
// Use '|' as a cursor marker in the SQL string, or pass cursorOffset >= 0.
func Complete(ctx context.Context, sql string, cursorOffset int, meta MetadataProvider) ([]Candidate, error) {
	// Resolve cursor position
	cleanSQL := sql
	actualOffset := cursorOffset
	if idx := findCursorMarker(sql); idx >= 0 {
		cleanSQL = sql[:idx] + sql[idx+1:]
		actualOffset = idx
	}
	if actualOffset < 0 {
		actualOffset = len(cleanSQL)
	}

	// Extract prefix (partial identifier at cursor) before simulation
	prefix := extractPrefix(cleanSQL, actualOffset)

	// Preprocess: isolate statement, tokenize up to cursor
	pp := preprocessFromClean(cleanSQL, actualOffset)

	// If there's a prefix being typed, the last token IS the prefix.
	// Exclude it from simulation so the state reflects the position
	// before the partial token, allowing correct valid-token collection.
	simTokens := pp.tokens
	if prefix != "" && len(simTokens) > 0 {
		lastTok := simTokens[len(simTokens)-1]
		if strings.EqualFold(lastTok.tok.Str, prefix) {
			simTokens = simTokens[:len(simTokens)-1]
		}
	}

	// Simulate parse up to cursor (excluding prefix token)
	state := simulateParse(simTokens)

	// Collect valid tokens at this position
	validTokens := collectValidTokens(state.stack)

	// Infer grammar contexts
	contexts := inferContexts(state.stack)

	hint := completionHint{
		validTokens: validTokens,
		contexts:    contexts,
	}

	// Expand hints to candidates using the full SQL (not just the cursor statement)
	candidates := expand(ctx, hint, cleanSQL, actualOffset, meta)

	// Filter by prefix
	if prefix != "" {
		candidates = filterByPrefix(candidates, prefix)
	}

	return candidates, nil
}

// preprocessFromClean preprocesses already-cleaned SQL (no cursor marker).
func preprocessFromClean(sql string, cursorOffset int) preprocessResult {
	allTokens := tokenize(sql)
	stmtTokens, cursorIdx := isolateStatement(allTokens, cursorOffset)
	return preprocessResult{
		tokens:   stmtTokens,
		cursorAt: cursorIdx,
	}
}

// filterByPrefix filters candidates whose text starts with the given prefix.
// For quoted identifiers like `"Order"`, matches against both the quoted and unquoted forms.
func filterByPrefix(candidates []Candidate, prefix string) []Candidate {
	var filtered []Candidate
	for _, c := range candidates {
		text := strings.ToLower(c.Text)
		if strings.HasPrefix(text, prefix) {
			filtered = append(filtered, c)
			continue
		}
		// Also match against the unquoted name for quoted identifiers
		if len(c.Text) >= 2 && c.Text[0] == '"' && c.Text[len(c.Text)-1] == '"' {
			unquoted := strings.ToLower(c.Text[1 : len(c.Text)-1])
			if strings.HasPrefix(unquoted, prefix) {
				filtered = append(filtered, c)
			}
		}
	}
	return filtered
}

// TokenNameForID returns the human-readable name for a goyacc token ID.
func TokenNameForID(tokID int) string {
	toknames := yacc.TokNames()
	if tokID >= 1 && tokID-1 < len(toknames) {
		name := toknames[tokID-1]
		if name != "" {
			return name
		}
	}
	return ""
}

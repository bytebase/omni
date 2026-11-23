// Package completion provides auto-complete for MongoDB shell (mongosh) commands.
package completion

import (
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mongo/catalog"
	"github.com/bytebase/omni/mongo/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType = analysis.CandidateType

const (
	CandidateKeyword       = analysis.CandidateKeyword       // top-level keywords (db, rs, sh, show, ...)
	CandidateCollection    = analysis.CandidateCollection    // collection name from catalog
	CandidateMethod        = analysis.CandidateMethod        // collection method (find, insertOne, ...)
	CandidateCursorMethod  = analysis.CandidateCursorMethod  // cursor modifier (sort, limit, ...)
	CandidateAggStage      = analysis.CandidateAggStage      // aggregation stage ($match, $group, ...)
	CandidateQueryOperator = analysis.CandidateQueryOperator // query operator ($gt, $in, ...)
	CandidateBSONHelper    = analysis.CandidateBSONHelper    // BSON constructor (ObjectId, NumberLong, ...)
	CandidateShowTarget    = analysis.CandidateShowTarget    // show command target (dbs, collections, ...)
	CandidateDbMethod      = analysis.CandidateDbMethod      // database method (getName, runCommand, ...)
	CandidateRsMethod      = analysis.CandidateRsMethod      // replica set method (status, conf, ...)
	CandidateShMethod      = analysis.CandidateShMethod      // sharding method (addShard, status, ...)
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition/signature
	Comment    string        // optional doc comment
}

// Complete returns completion candidates for the given mongosh input at the cursor offset.
// cat may be nil if no catalog context is available.
func Complete(input string, cursorOffset int, cat *catalog.Catalog) []Candidate {
	if cursorOffset > len(input) {
		cursorOffset = len(input)
	}

	prefix := extractPrefix(input, cursorOffset)
	tokens := tokenize(input, cursorOffset-len(prefix))
	ctx := detectContext(tokens)
	candidates := candidatesForContext(ctx, cat)

	return filterByPrefix(candidates, prefix)
}

// tokenize lexes input up to the given byte offset and returns all tokens.
func tokenize(input string, limit int) []parser.Token {
	if limit > len(input) {
		limit = len(input)
	}
	if limit < 0 {
		limit = 0
	}
	lex := parser.NewLexer(input[:limit])
	var tokens []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == parser.TokEOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

// extractPrefix returns the partial token the user is typing at cursorOffset.
// Includes $ as a valid prefix character (for $match, $gt, etc.).
func extractPrefix(input string, cursorOffset int) string {
	if cursorOffset > len(input) {
		cursorOffset = len(input)
	}
	i := cursorOffset
	for i > 0 {
		c := input[i-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' || c == '$' {
			i--
		} else {
			break
		}
	}
	return input[i:cursorOffset]
}

// filterByPrefix filters candidates whose Text starts with prefix.
// Matching is case-sensitive (mongosh is case-sensitive).
func filterByPrefix(candidates []Candidate, prefix string) []Candidate {
	if prefix == "" {
		return candidates
	}
	var result []Candidate
	for _, c := range candidates {
		if strings.HasPrefix(c.Text, prefix) {
			result = append(result, c)
		}
	}
	return result
}

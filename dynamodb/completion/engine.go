package completion

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/dynamodb/catalog"
	"github.com/bytebase/omni/metadata"
)

// Engine is DynamoDB's analysis.Completer.
var Engine engine

type engine struct{}

var _ analysis.Completer = Engine

// Complete completes from the default database's tables only. A read-only
// request leaves out the keywords only a change uses.
func (engine) Complete(ctx context.Context, request analysis.CompletionRequest, script string, caretLine, caretOffset int) []analysis.Candidate {
	cat := catalog.New()
	var meta *metadata.DatabaseSchemaMetadata
	if load := request.Catalog.MetadataLoader(); load != nil && request.DefaultDatabase != "" {
		if m, err := load(ctx, request.DefaultDatabase); err == nil && m != nil {
			meta = m
			cat.LoadMetadata(meta)
		}
	}
	candidates := CompleteWithMetadata(script, lineOffsetToBytePos(script, caretLine, caretOffset), cat, meta)
	var result []analysis.Candidate
	for _, c := range candidates {
		if request.ReadOnly && c.Type == CandidateKeyword && isDMLKeyword(c.Text) {
			continue
		}
		result = append(result, analysis.Candidate{Text: c.Text, Type: c.Type})
	}
	return result
}

// lineOffsetToBytePos converts a 1-based line and a 0-based offset counted in
// runes, not UTF-16 code units, to a byte position in s.
func lineOffsetToBytePos(s string, line, offset int) int {
	pos := 0
	currentLine := 1
	for pos < len(s) && currentLine < line {
		if s[pos] == '\n' {
			currentLine++
		}
		pos++
	}
	for i := 0; i < offset && pos < len(s); i++ {
		_, size := utf8.DecodeRuneInString(s[pos:])
		pos += size
	}
	return pos
}

func isDMLKeyword(w string) bool {
	switch strings.ToUpper(w) {
	case "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "REPLACE",
		"UPSERT", "REMOVE", "SET":
		return true
	}
	return false
}

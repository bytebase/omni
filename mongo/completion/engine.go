package completion

import (
	"context"

	"github.com/bytebase/omni/analysis"
)

// Engine is MongoDB's analysis.Completer.
var Engine engine

type engine struct{}

var _ analysis.Completer = Engine

// Complete completes against the default database; MongoDB has no schemas.
func (engine) Complete(ctx context.Context, request analysis.CompletionRequest, script string, caretLine, caretOffset int) []analysis.Candidate {
	return convertCandidates(CompleteSession(ctx, SessionConfig{
		DefaultDatabase: request.DefaultDatabase,
		LoadDatabase:    request.Catalog.MetadataLoader(),
	}, script, caretLine, caretOffset))
}

// convertCandidates offers an untyped candidate as a function.
func convertCandidates(candidates []Candidate) []analysis.Candidate {
	if candidates == nil {
		return nil
	}
	converted := make([]analysis.Candidate, 0, len(candidates))
	for _, c := range candidates {
		typ := c.Type
		if typ == analysis.CandidateTypeNone {
			typ = analysis.CandidateFunction
		}
		converted = append(converted, analysis.Candidate{Text: c.Text, Type: typ, Definition: c.Definition, Comment: c.Comment})
	}
	return converted
}

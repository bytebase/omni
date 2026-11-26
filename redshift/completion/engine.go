package completion

import (
	"context"

	"github.com/bytebase/omni/analysis"
)

// Engine is Redshift's analysis.Completer.
var Engine engine

type engine struct{}

var _ analysis.Completer = Engine

func (engine) Complete(ctx context.Context, request analysis.CompletionRequest, script string, caretLine, caretOffset int) []analysis.Candidate {
	return convertCandidates(CompleteSession(ctx, SessionConfig{
		DefaultDatabase: request.DefaultDatabase,
		DefaultSchema:   request.DefaultSchema,
		LoadDatabase:    request.Catalog.MetadataLoader(),
	}, script, caretLine, caretOffset))
}

func convertCandidates(candidates []Candidate) []analysis.Candidate {
	if candidates == nil {
		return nil
	}
	converted := make([]analysis.Candidate, 0, len(candidates))
	for _, c := range candidates {
		converted = append(converted, analysis.Candidate{Text: c.Text, Type: c.Type, Definition: c.Definition, Comment: c.Comment})
	}
	return converted
}

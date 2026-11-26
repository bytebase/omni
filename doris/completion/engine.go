package completion

import (
	"context"

	"github.com/bytebase/omni/analysis"
)

// Engine is Doris's analysis.Completer.
var Engine engine

type engine struct{}

var _ analysis.Completer = Engine

func (engine) Complete(ctx context.Context, request analysis.CompletionRequest, script string, caretLine, caretOffset int) []analysis.Candidate {
	candidates := CompleteSession(ctx, SessionConfig{
		DefaultDatabase:   request.DefaultDatabase,
		LoadDatabase:      request.Catalog.MetadataLoader(),
		ListDatabaseNames: request.Catalog.ListDatabaseNames,
	}, script, caretLine, caretOffset)
	converted := make([]analysis.Candidate, 0, len(candidates))
	for _, c := range candidates {
		converted = append(converted, analysis.Candidate{Text: c.Text, Type: c.Type, Definition: c.Definition, Comment: c.Comment})
	}
	return converted
}

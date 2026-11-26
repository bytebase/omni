package completion

import (
	"context"

	omnianalysis "github.com/bytebase/omni/analysis"
)

// Engine is Oracle's analysis.Completer.
var Engine engine

type engine struct{}

var _ omnianalysis.Completer = Engine

// Complete treats the request's default database as the default schema, as an
// Oracle database is a schema.
func (engine) Complete(ctx context.Context, request omnianalysis.CompletionRequest, script string, caretLine, caretOffset int) []omnianalysis.Candidate {
	return convertCandidates(CompleteSession(ctx, SessionConfig{
		ReadOnly:        request.ReadOnly,
		DefaultSchema:   request.DefaultDatabase,
		LoadSchema:      request.Catalog.MetadataLoader(),
		ListSchemaNames: request.Catalog.ListDatabaseNames,
	}, script, caretLine, caretOffset))
}

// convertCandidates reports package variables as user variables, not server
// system variables.
func convertCandidates(candidates []Candidate) []omnianalysis.Candidate {
	if candidates == nil {
		return nil
	}
	converted := make([]omnianalysis.Candidate, 0, len(candidates))
	for _, c := range candidates {
		typ := c.Type
		if typ == CandidateVariable {
			typ = omnianalysis.CandidateUserVariable
		}
		converted = append(converted, omnianalysis.Candidate{Text: c.Text, Type: typ, Definition: c.Definition, Comment: c.Comment})
	}
	return converted
}

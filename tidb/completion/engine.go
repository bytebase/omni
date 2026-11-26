package completion

import (
	"context"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
)

// Engine is TiDB's analysis.Completer.
var Engine engine

type engine struct{}

var _ analysis.Completer = Engine

func (engine) Complete(ctx context.Context, request analysis.CompletionRequest, script string, caretLine, caretOffset int) []analysis.Candidate {
	return convertCandidates(CompleteSession(ctx, SessionConfig{
		ReadOnly:          request.ReadOnly,
		DefaultDatabase:   request.DefaultDatabase,
		LoadDatabase:      databaseLoader(request.Catalog),
		ListDatabaseNames: request.Catalog.ListDatabaseNames,
	}, script, caretLine, caretOffset))
}

// databaseLoader adapts the catalog to SessionConfig.LoadDatabase, or returns
// nil when the catalog cannot look databases up.
func databaseLoader(catalog analysis.Catalog) func(ctx context.Context, name string) (string, *metadata.DatabaseSchemaMetadata, bool, bool) {
	if catalog.GetDatabase == nil {
		return nil
	}
	return func(ctx context.Context, name string) (string, *metadata.DatabaseSchemaMetadata, bool, bool) {
		database, err := catalog.GetDatabase(ctx, name)
		if err != nil || database == nil {
			return "", nil, false, false
		}
		return database.Name, database.Metadata, database.ObjectCaseSensitive, true
	}
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

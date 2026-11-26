package elasticsearch

import "github.com/bytebase/omni/analysis"

// Engine is Elasticsearch's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var _ analysis.StatementRanger = Engine

// StatementRanges returns an empty, non-nil list for a script with no request.
func (engine) StatementRanges(script string) []analysis.LSPRange {
	// GetStatementRanges reports a script it cannot parse as having no
	// requests; its error is always nil.
	ranges, _ := GetStatementRanges(script)
	converted := make([]analysis.LSPRange, 0, len(ranges))
	for _, r := range ranges {
		converted = append(converted, analysis.LSPRange{
			Start: analysis.LSPPosition{Line: r.Start.Line, Character: r.Start.Character},
			End:   analysis.LSPPosition{Line: r.End.Line, Character: r.End.Character},
		})
	}
	return converted
}

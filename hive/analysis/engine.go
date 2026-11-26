package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/hive/parser"
)

// Engine is Hive's implementation of the omni/analysis capability interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.QueryValidator = Engine
	_ analysis.Explainer      = Engine
)

func (engine) Split(script string) []analysis.Statement {
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, seg := range segments {
		statements = append(statements, analysis.Statement{
			Text:  seg.Text,
			Empty: seg.Empty(),
			Start: seg.ByteStart,
			End:   seg.ByteEnd,
		})
	}
	return statements
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

// Explain plans with EXPLAIN FORMATTED when asked for JSON and with EXPLAIN
// otherwise.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	name := ""
	if format == analysis.ExplainFormatJSON {
		name = "JSON"
	}
	return ExplainStatement(statement, name), nil
}

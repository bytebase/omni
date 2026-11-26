package analysis

import (
	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/databricks/parser"
)

// Engine is Databricks's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter  = Engine
	_ analysis.Explainer = Engine
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

// Explain plans with EXPLAIN FORMATTED whatever format the caller asks for.
func (engine) Explain(statement string, _ analysis.ExplainFormat) (string, error) {
	return ExplainStatement(statement, ""), nil
}

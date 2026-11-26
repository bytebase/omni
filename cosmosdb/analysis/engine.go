package analysis

import "github.com/bytebase/omni/analysis"

// Engine is Cosmos DB's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var _ analysis.QueryValidator = Engine

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

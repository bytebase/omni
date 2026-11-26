package analysis

import "github.com/bytebase/omni/analysis"

// Engine is Redis's implementation of the omni/analysis capability interfaces.
var Engine engine

type engine struct{}

var _ analysis.QueryValidator = Engine

// ValidateQuery returns the splitter's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

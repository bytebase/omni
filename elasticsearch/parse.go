// Package elasticsearch provides a parser for Kibana Dev Console-style
// Elasticsearch REST request blocks.
package elasticsearch

// ParseResult holds parsed requests and any syntax errors encountered.
type ParseResult struct {
	Requests []*Request     `yaml:"requests"`
	Errors   []*SyntaxError `yaml:"errors,omitempty"`
}

// Request is a single parsed Kibana Dev Console REST request.
type Request struct {
	Method      string   `yaml:"method"`
	URL         string   `yaml:"url"`
	Data        []string `yaml:"data,omitempty"`
	StartOffset int      `yaml:"startoffset"`
	EndOffset   int      `yaml:"endoffset"`
}

// SyntaxError represents a parse error at a specific source position.
type SyntaxError struct {
	Position   Position `yaml:"position"`
	Message    string   `yaml:"message"`
	RawMessage string   `yaml:"rawmessage"`
}

// Position represents a location in source text.
type Position struct {
	Line   int `yaml:"line"`
	Column int `yaml:"column"`
}

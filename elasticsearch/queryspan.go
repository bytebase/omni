package elasticsearch

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/elasticsearch/analysis"
)

// QuerySpan is the result of analysing a single Elasticsearch REST statement.
// It contains the query classification, the masking analysis, and the set of
// predicate field paths expressed as PathAST chains.
type QuerySpan struct {
	// Type is the high-level classification of the request.
	Type QueryType
	// ElasticsearchAnalysis holds the body-level masking analysis result.
	ElasticsearchAnalysis *analysis.RequestAnalysis
	// PredicatePaths maps each dot-notation predicate field name to its PathAST
	// representation.  Nil when the request has no predicate fields.
	PredicatePaths map[string]*PathAST
}

// PathAST is a linked chain of path segments that represents a dot-delimited
// field reference such as "user.address.city".  Each node is an ItemSelector
// that holds a single segment name and a pointer to the next segment.
type PathAST struct {
	// Root is the first segment of the path.
	Root *ItemSelector
}

// ItemSelector is a single named segment in a PathAST chain.
type ItemSelector struct {
	// Name is the path segment (e.g. "user", "address", "city").
	Name string
	// Next is the following segment, or nil for the last segment.
	Next *ItemSelector
}

// String returns the dot-joined representation of the full PathAST chain
// starting at this node, useful for debugging.
func (s *ItemSelector) String() string {
	if s.Next == nil {
		return s.Name
	}
	return fmt.Sprintf("%s.%s", s.Name, s.Next.String())
}

// GetQuerySpan parses a single Elasticsearch REST statement and returns a
// QuerySpan describing its type, masking analysis, and predicate field paths.
func GetQuerySpan(statement string) (*QuerySpan, error) {
	parseResult, err := ParseElasticsearchREST(statement)
	if err != nil {
		return nil, err
	}

	if parseResult == nil {
		return &QuerySpan{Type: QueryTypeUnknown}, nil
	}

	if len(parseResult.Errors) > 0 {
		firstErr := parseResult.Errors[0]
		return nil, fmt.Errorf("syntax error at line %d, column %d: %s",
			firstErr.Position.Line, firstErr.Position.Column, firstErr.Message)
	}

	if len(parseResult.Requests) == 0 {
		return &QuerySpan{Type: QueryTypeUnknown}, nil
	}

	// After splitting, each statement should contain a single request.
	// Use the first request for classification.
	req := parseResult.Requests[0]
	queryType := ClassifyRequest(req.Method, req.URL)

	span := &QuerySpan{Type: queryType}

	requestAnalysis := analysis.AnalyzeRequest(req.Method, req.URL, strings.Join(req.Data, "\n"))
	span.ElasticsearchAnalysis = requestAnalysis

	if len(requestAnalysis.PredicateFields) > 0 {
		span.PredicatePaths = make(map[string]*PathAST, len(requestAnalysis.PredicateFields))
		for _, field := range requestAnalysis.PredicateFields {
			span.PredicatePaths[field] = DotPathToPathAST(field)
		}
	}

	return span, nil
}

// DotPathToPathAST converts a dot-delimited field path (e.g. "contact.phone")
// into a PathAST with linked ItemSelector nodes.
func DotPathToPathAST(dotPath string) *PathAST {
	parts := strings.Split(dotPath, ".")
	if len(parts) == 0 {
		return nil
	}
	root := &ItemSelector{Name: parts[0]}
	current := root
	for _, part := range parts[1:] {
		next := &ItemSelector{Name: part}
		current.Next = next
		current = next
	}
	return &PathAST{Root: root}
}

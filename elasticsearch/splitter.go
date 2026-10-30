package elasticsearch

import (
	"strings"

	"github.com/bytebase/omni/review"
)

// Statement is a single parsed Elasticsearch REST API request with position and byte range info.
type Statement struct {
	Text  string   `yaml:"text,omitempty"`
	Empty bool     `yaml:"empty,omitempty"`
	Start Position `yaml:"start,omitempty"`
	End   Position `yaml:"end,omitempty"`
	Range Range    `yaml:"range,omitempty"`
}

// Range holds the byte start and end offsets of a statement.
type Range struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

// SplitMultiSQL splits the input into individual Elasticsearch REST API requests.
// A request whose body fails to parse still becomes a statement, the last
// request of the script too: it runs through the last non-blank line before
// the next request, the lines Kibana's console reads its body from.
func SplitMultiSQL(statement string) ([]Statement, error) {
	parseResult, ends, err := parseREST(statement, false)
	if err != nil {
		return nil, err
	}

	if parseResult == nil || len(parseResult.Requests) == 0 {
		if len(strings.TrimSpace(statement)) == 0 {
			return nil, nil
		}
		return []Statement{{Text: statement}}, nil
	}

	index := review.Index(statement)
	var statements []Statement
	for i, req := range parseResult.Requests {
		if req == nil {
			continue
		}
		end := ends[i]
		text := statement[req.StartOffset:end]
		empty := len(strings.TrimSpace(text)) == 0

		startLine, startColumn := index.Position(req.StartOffset)
		endLine, endColumn := index.Position(end)

		statements = append(statements, Statement{
			Text:  text,
			Empty: empty,
			Start: Position{
				Line:   startLine,
				Column: startColumn,
			},
			End: Position{
				Line:   endLine,
				Column: endColumn,
			},
			Range: Range{
				Start: req.StartOffset,
				End:   end,
			},
		})
	}
	return statements, nil
}

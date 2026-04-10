package elasticsearch

import "strings"

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
func SplitMultiSQL(statement string) ([]Statement, error) {
	parseResult, err := ParseElasticsearchREST(statement)
	if err != nil {
		return nil, err
	}

	if parseResult == nil || len(parseResult.Requests) == 0 {
		if len(strings.TrimSpace(statement)) == 0 {
			return nil, nil
		}
		return []Statement{{Text: statement}}, nil
	}

	var statements []Statement
	for _, req := range parseResult.Requests {
		if req == nil {
			continue
		}
		text := statement[req.StartOffset:req.EndOffset]
		empty := len(strings.TrimSpace(text)) == 0

		startLine, startColumn := byteOffsetToPosition(statement, req.StartOffset)
		endLine, endColumn := byteOffsetToPosition(statement, req.EndOffset)

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
				End:   req.EndOffset,
			},
		})
	}
	return statements, nil
}

// byteOffsetToPosition converts a byte offset to 1-based line and column numbers.
// Column is measured in Unicode code points (runes), not bytes.
func byteOffsetToPosition(text string, byteOffset int) (line, column int) {
	line = 1
	column = 1
	currentByte := 0

	for _, r := range text {
		if currentByte >= byteOffset {
			break
		}
		if r == '\n' {
			line++
			column = 1
		} else {
			column++
		}
		currentByte += len(string(r))
	}
	return line, column
}

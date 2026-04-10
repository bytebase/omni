package elasticsearch

import (
	"slices"

	"github.com/bytebase/omni/elasticsearch/parser"
)

// ParseElasticsearchREST parses the Elasticsearch REST API request text into
// structured requests and syntax errors. This is the primary entry point for
// the Elasticsearch parser.
func ParseElasticsearchREST(text string) (*ParseResult, error) {
	raw, err := parser.Parse(text)
	if err != nil {
		return nil, err
	}

	var requests []*Request
	for i, pr := range raw.Requests {
		// See https://sourcegraph.com/github.com/elastic/kibana/-/blob/src/platform/plugins/shared/console/public/application/containers/editor/monaco_editor_actions_provider.ts?L261.
		var nextRequest *parser.ParsedRequest
		if i < len(raw.Requests)-1 {
			nextRequest = &raw.Requests[i+1]
		}
		adjustedOffset := parser.GetAdjustedParsedRequest(pr, text, nextRequest)
		if adjustedOffset.StartLineNumber > adjustedOffset.EndLineNumber {
			continue
		}
		editorRequest := parser.GetEditorRequest(text, adjustedOffset)
		if editorRequest == nil {
			continue
		}
		if len(editorRequest.Data) > 0 {
			for j, v := range editorRequest.Data {
				if parser.ContainsComments(v) {
					// parse and stringify to remove comments.
					editorRequest.Data[j] = parser.IndentData(v)
				}
				editorRequest.Data[j] = parser.CollapseLiteralString(editorRequest.Data[j])
			}
		}
		requests = append(requests, &Request{
			Method:      editorRequest.Method,
			URL:         editorRequest.URL,
			Data:        editorRequest.Data,
			StartOffset: pr.StartOffset,
			EndOffset:   pr.EndOffset,
		})
	}

	// Convert errors
	var syntaxErrors []*SyntaxError

	slices.SortFunc(raw.Errors, func(a, b parser.SyntaxError) int {
		if a.ByteOffset < b.ByteOffset {
			return -1
		}
		if a.ByteOffset > b.ByteOffset {
			return 1
		}
		return 0
	})

	for i, rawErr := range raw.Errors {
		line := 0
		column := 0
		pos := 0
		if i > 0 {
			line = syntaxErrors[i-1].Position.Line
			column = syntaxErrors[i-1].Position.Column
			pos = raw.Errors[i-1].ByteOffset
		}
		boundary := rawErr.ByteOffset
		if boundary >= len(text) {
			boundary = len(text) - 1
		}

		for j := pos; j <= boundary; j++ {
			if text[j] == '\n' {
				if j == boundary {
					// Decorate the \n position instead of the next line.
					column++
					continue
				}
				line++
				column = 0
			} else {
				column++
			}
		}
		syntaxErrors = append(syntaxErrors, &SyntaxError{
			Position: Position{
				Line:   line,
				Column: column,
			},
			Message: rawErr.Message,
		})
	}

	return &ParseResult{
		Requests: requests,
		Errors:   syntaxErrors,
	}, nil
}

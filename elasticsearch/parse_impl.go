package elasticsearch

import (
	"slices"
	"strings"
	"unicode"

	"github.com/bytebase/omni/elasticsearch/parser"
)

// ParseElasticsearchREST parses the Elasticsearch REST API request text into
// structured requests and syntax errors. This is the primary entry point for
// the Elasticsearch parser.
func ParseElasticsearchREST(text string) (*ParseResult, error) {
	result, _, err := parseREST(text, true)
	return result, err
}

// parseREST is ParseElasticsearchREST that also returns, for each request,
// the byte offset where its text ends (see requestEnd). Without cleanData it
// leaves each document of Data as the request text holds it, for callers that
// do not read Data: rewriting a document that holds comments as indented JSON
// writes indentation that grows with its depth, 200MB for a 110KB document
// nested 10,000 deep.
func parseREST(text string, cleanData bool) (*ParseResult, []int, error) {
	raw, err := parser.Parse(text)
	if err != nil {
		return nil, nil, err
	}

	var requests []*Request
	var ends []int
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
		if cleanData && len(editorRequest.Data) > 0 {
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
		ends = append(ends, requestEnd(text, pr, nextRequest, adjustedOffset.EndLineNumber))
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
	}, ends, nil
}

// requestEnd returns the byte offset where the text of request r ends. A
// request the parser finished ends at its EndOffset. A request whose body
// failed to parse keeps EndOffset -1; Kibana's console takes it to run through
// endLine, the last non-blank line before the next request, and reads its body
// from those lines. Its text therefore ends at the end of endLine, without
// trailing whitespace, and never past the start of the next request. A
// request that holds only whitespace keeps it, and splits as an empty
// statement.
func requestEnd(text string, r parser.ParsedRequest, next *parser.ParsedRequest, endLine int) int {
	if r.EndOffset >= 0 {
		return r.EndOffset
	}
	end := len(text)
	for i, line := 0, 0; i < len(text); i++ {
		if text[i] == '\n' {
			if line == endLine {
				end = i
				break
			}
			line++
		}
	}
	if next != nil && next.StartOffset < end {
		end = next.StartOffset
	}
	end = max(end, r.StartOffset)
	if trimmed := strings.TrimRightFunc(text[r.StartOffset:end], unicode.IsSpace); trimmed != "" {
		return r.StartOffset + len(trimmed)
	}
	return end
}

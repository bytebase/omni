package elasticsearch

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytebase/omni/elasticsearch/parser"
	"github.com/bytebase/omni/review"
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

	lines := parser.NewLines(text)
	var requests []*Request
	var ends []int
	for i, pr := range raw.Requests {
		// See https://sourcegraph.com/github.com/elastic/kibana/-/blob/src/platform/plugins/shared/console/public/application/containers/editor/monaco_editor_actions_provider.ts?L261.
		var nextRequest *parser.ParsedRequest
		if i < len(raw.Requests)-1 {
			nextRequest = &raw.Requests[i+1]
		}
		adjustedOffset := lines.AdjustedParsedRequest(pr, nextRequest)
		if adjustedOffset.StartLineNumber > adjustedOffset.EndLineNumber {
			continue
		}
		editorRequest := lines.EditorRequest(adjustedOffset)
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
		ends = append(ends, requestEnd(text, lines, pr, nextRequest, adjustedOffset.EndLineNumber))
	}

	// Convert errors. A stable sort keeps the order in which errors at one
	// offset were recorded: the "Syntax error" for text left over after the
	// last request comes after the error that stopped it.
	var syntaxErrors []*SyntaxError

	slices.SortStableFunc(raw.Errors, func(a, b parser.SyntaxError) int {
		if a.ByteOffset < b.ByteOffset {
			return -1
		}
		if a.ByteOffset > b.ByteOffset {
			return 1
		}
		return 0
	})

	index := review.Index(text)
	for _, rawErr := range raw.Errors {
		syntaxErrors = append(syntaxErrors, &SyntaxError{
			Position: errorPosition(text, index, rawErr.ByteOffset),
			Message:  rawErr.Message,
		})
	}

	return &ParseResult{
		Requests: requests,
		Errors:   syntaxErrors,
	}, ends, nil
}

// errorPosition returns the position of a syntax error at byte offset: the
// line review.Index gives, made 0-based as Diagnose documents, and its column,
// 1-based and counted in code points, so a byte that is not UTF-8 counts as
// one. An offset at or past the end of the text stands for the last
// character. The column was counted in bytes, and an error on the line of the
// error before it counted that error's byte a second time.
func errorPosition(text string, index *review.Text, offset int) Position {
	if offset >= len(text) {
		_, size := utf8.DecodeLastRuneInString(text)
		offset = len(text) - size
	}
	line, column := index.Position(offset)
	return Position{Line: line - 1, Column: column}
}

// requestEnd returns the byte offset where the text of request r ends. A
// request the parser finished ends at its EndOffset. A request whose body
// failed to parse keeps EndOffset -1; Kibana's console takes it to run through
// endLine, the last non-blank line before the next request, and reads its body
// from those lines. Its text therefore ends at the end of endLine, without
// trailing whitespace, and never past the start of the next request. A
// request that holds only whitespace keeps it, and splits as an empty
// statement.
func requestEnd(text string, lines *parser.Lines, r parser.ParsedRequest, next *parser.ParsedRequest, endLine int) int {
	if r.EndOffset >= 0 {
		return r.EndOffset
	}
	end := lines.LineEnd(endLine)
	if next != nil && next.StartOffset < end {
		end = next.StartOffset
	}
	end = max(end, r.StartOffset)
	if trimmed := strings.TrimRightFunc(text[r.StartOffset:end], unicode.IsSpace); trimmed != "" {
		return r.StartOffset + len(trimmed)
	}
	return end
}

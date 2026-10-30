package elasticsearch

import (
	"strings"
	"unicode/utf8"
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
// A request whose body fails to parse still becomes a statement: it runs
// through the last non-blank line before the next request, the lines Kibana's
// console reads its body from.
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

	positions := newPositioner(statement)
	var statements []Statement
	for i, req := range parseResult.Requests {
		if req == nil {
			continue
		}
		end := ends[i]
		text := statement[req.StartOffset:end]
		empty := len(strings.TrimSpace(text)) == 0

		startLine, startColumn := positions.position(req.StartOffset)
		endLine, endColumn := positions.position(end)

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

// positioner converts byte offsets to 1-based line and column numbers.
// Column is measured in Unicode code points (runes), not bytes. It resumes
// from the offset it was last asked for instead of scanning the text from
// its start for every offset, which made splitting a script of n requests
// take time quadratic in n.
type positioner struct {
	text string
	// next is the byte offset of the next rune to read, and consumed the
	// bytes read so far as the scan counts them: an invalid byte reads as
	// U+FFFD and counts as its three bytes.
	next, consumed int
	line, column   int
	// last is the offset last asked for.
	last int
}

func newPositioner(text string) *positioner {
	return &positioner{text: text, line: 1, column: 1}
}

// position returns the line and column at byteOffset.
func (p *positioner) position(byteOffset int) (line, column int) {
	if byteOffset < p.last {
		// Offsets come in order; start over for one that does not.
		*p = *newPositioner(p.text)
	}
	p.last = byteOffset
	for p.next < len(p.text) && p.consumed < byteOffset {
		r, size := utf8.DecodeRuneInString(p.text[p.next:])
		if r == '\n' {
			p.line++
			p.column = 1
		} else {
			p.column++
		}
		p.consumed += utf8.RuneLen(r)
		p.next += size
	}
	return p.line, p.column
}

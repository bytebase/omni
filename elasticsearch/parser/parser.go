// Package parser implements a hand-written recursive-descent parser for
// Kibana Dev Console-style Elasticsearch REST request blocks.
//
// This is a faithful port of the parser from
// bytebase/bytebase/backend/plugin/parser/elasticsearch/parser.go.
// The state machine logic, byte-offset accounting, comment handling, HJSON
// normalization, error recovery, and multi-document body splitting are
// preserved exactly as in the original.
package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	hjson "github.com/hjson/hjson-go/v4"
)

// ParseResult holds the raw parse output: offset-tracked request ranges and
// syntax errors with byte offsets. The public API in the parent package
// converts these into the omni-level types.
type ParseResult struct {
	Requests []ParsedRequest
	Errors   []SyntaxError
}

// ParsedRequest is the byte range of a single request (left inclusive, right exclusive).
type ParsedRequest struct {
	// StartOffset is the byte offset of the first character of the request.
	StartOffset int
	// EndOffset is the byte offset of the end position.
	EndOffset int
}

// SyntaxError records a parse error at a byte offset.
type SyntaxError struct {
	// ByteOffset is the byte offset of the first character where the error occurred, starting at 0.
	ByteOffset int
	Message    string
}

// EditorRequest is a single parsed request with method, URL, and body data.
type EditorRequest struct {
	Method string
	URL    string
	Data   []string
}

// AdjustedParsedRequest maps byte offsets to 0-based line numbers.
type AdjustedParsedRequest struct {
	// StartLineNumber is the line number of the first character of the request, starting from 0.
	StartLineNumber int
	// EndLineNumber is the line number of the end position, starting from 0.
	EndLineNumber int
}

// Lines indexes the lines of a text once, for the functions that map a
// request to its lines. GetAdjustedParsedRequest and GetEditorRequest split
// the whole text and count its newlines from the start on every call, so
// calling them for each request of a script took time quadratic in the
// number of requests.
type Lines struct {
	text  string
	lines []string
	// starts holds the byte offset where each line starts.
	starts []int
}

// NewLines indexes the lines of text.
func NewLines(text string) *Lines {
	lines := strings.Split(text, "\n")
	starts := make([]int, len(lines))
	offset := 0
	for i, line := range lines {
		starts[i] = offset
		offset += len(line) + 1
	}
	return &Lines{text: text, lines: lines, starts: starts}
}

// LineEnd returns the byte offset where 0-based line n ends: the offset of
// its newline, or the length of the text for the last line.
func (l *Lines) LineEnd(n int) int {
	return l.starts[n] + len(l.lines[n])
}

// newlinesBefore counts the newlines in the text before byte offset. The
// newline that ends line i is at starts[i+1]-1, so they are as many as the
// lines after the first that start at or before offset.
func (l *Lines) newlinesBefore(offset int) int {
	return max(sort.SearchInts(l.starts, offset+1)-1, 0)
}

// newlinesBetween counts the newlines at byte offsets from start up to, not
// including, end.
func (l *Lines) newlinesBetween(start, end int) int {
	if start >= end {
		return 0
	}
	return l.newlinesBefore(end) - l.newlinesBefore(start)
}

// startsWithMethodRegex matches a line that starts a request.
var startsWithMethodRegex = regexp.MustCompile(`(?i)^\s*(GET|POST|PUT|PATCH|DELETE)`)

// GetAdjustedParsedRequest is Lines.AdjustedParsedRequest for one request of
// text. Callers mapping every request of a text index it once with NewLines.
func GetAdjustedParsedRequest(r ParsedRequest, text string, nextRequest *ParsedRequest) AdjustedParsedRequest {
	return NewLines(text).AdjustedParsedRequest(r, nextRequest)
}

// AdjustedParsedRequest returns the lines request r spans.
// See https://sourcegraph.com/github.com/elastic/kibana/-/blob/src/platform/plugins/shared/console/public/application/containers/editor/utils/requests_utils.ts?L76.
// Combine getRequestStartLineNumber and getRequestEndLineNumber.
func (l *Lines) AdjustedParsedRequest(r ParsedRequest, nextRequest *ParsedRequest) AdjustedParsedRequest {
	size := len(l.text)
	endLineNumber := 0
	startOffset := r.StartOffset
	// The startOffset is out of range, returning the end of document like
	// what the model.getPositionAt does.
	if r.StartOffset >= size {
		startOffset = size - 1
	}
	startLineNumber := l.newlinesBefore(r.StartOffset)

	if r.EndOffset >= 0 {
		// if the parser set an end offset for this request, then find the line number for it.
		endOffset := r.EndOffset
		if endOffset >= size {
			endOffset = size - 1
		}
		endLineNumber = startLineNumber + l.newlinesBetween(startOffset, endOffset)
	} else {
		// if no end offset, try to find the line before the next request starts.
		if nextRequest != nil {
			nextRequestOffset := nextRequest.StartOffset
			if nextRequestOffset >= size {
				nextRequestOffset = size - 1
			}
			nextRequestStartLine := l.newlinesBefore(nextRequestOffset)
			if nextRequestStartLine > startLineNumber {
				endLineNumber = nextRequestStartLine - 1
			} else {
				endLineNumber = startLineNumber
			}
		} else {
			// If there is no next request, find the end of the text or the line that starts with a method.
			nextLineNumber := startLineNumber + 1
			for nextLineNumber < len(l.lines) {
				content := strings.TrimSpace(l.lines[nextLineNumber])
				if startsWithMethodRegex.MatchString(content) {
					break
				}
				nextLineNumber++
			}
			// nextLineNumber is now either the line with a method or 1 line after the end of the text
			// set the end line for this request to the line before nextLineNumber
			if nextLineNumber > startLineNumber {
				endLineNumber = nextLineNumber - 1
			} else {
				endLineNumber = startLineNumber
			}
		}
	}
	// if the end is empty, go up to find the first non-empty line.
	for endLineNumber >= 0 && strings.TrimSpace(l.lines[endLineNumber]) == "" {
		endLineNumber--
	}
	return AdjustedParsedRequest{
		StartLineNumber: startLineNumber,
		EndLineNumber:   endLineNumber,
	}
}

type state struct {
	// at is the byte offset just past the current rune, ch. Reading past the
	// end of the text moves at past len(text) and sets ch to 0, as the
	// upstream's text.charAt moves at past the end and returns ''; eof tells
	// that state apart from a NUL byte in the text.
	at       int
	ch       rune
	escapee  map[rune]string
	text     string
	errors   []SyntaxError
	requests []ParsedRequest
	// requestStartOffset is the byte offset of the first character of the request.
	requestStartOffset int
	// requestEndOffset is the byte offset of the first character of the next request.
	requestEndOffset int
	// depth is the nesting depth of the body values being parsed (see enterNesting).
	depth int
}

func newState(text string) *state {
	return &state{
		at: 0,
		ch: 0,
		escapee: map[rune]string{
			'"':  `"`,
			'\\': `\`,
			'/':  `/`,
			'b':  "\b",
			'f':  "\f",
			'n':  "\n",
			'r':  "\r",
			't':  "\t",
		},
		text: text,
	}
}

// eof reports whether the parser has read past the last rune of the text.
// The upstream tests its current character for the empty string, which only
// the end of the text gives: a NUL byte is an ordinary character there,
// whitespace between tokens and part of a URL, string or comment. Testing ch
// for 0 instead stopped the parse at a NUL byte as if the text ended there
// and dropped the requests after it without an error.
func (s *state) eof() bool {
	return s.at > len(s.text)
}

// Parse parses the input text and returns the raw parse result.
func Parse(text string) (*ParseResult, error) {
	s := newState(text)
	requests, err := s.parse()
	if err != nil {
		return nil, err
	}
	return &ParseResult{
		Requests: requests,
		Errors:   s.errors,
	}, nil
}

// See https://sourcegraph.com/github.com/elastic/kibana/-/blob/src/platform/packages/shared/kbn-monaco/src/languages/console/parser.js.
func (s *state) parse() ([]ParsedRequest, error) {
	if _, err := s.nextEmptyInput(); err != nil {
		return nil, err
	}
	if err := s.multiRequest(); err != nil {
		return nil, err
	}

	// Text is left over only when multiRequest stopped at an error with no
	// request after it to resume at. The upstream records the leftover as one
	// more error of the script, "Syntax error", rather than throwing, so the
	// requests before it, and the one that failed, are still returned. A white
	// that fails stops at a byte it cannot decode, which leaves text over too.
	if err := s.white(); err != nil || !s.eof() {
		s.errors = append(s.errors, SyntaxError{ByteOffset: s.at, Message: "Syntax error"})
	}
	return s.requests, nil
}

// requestStartRegex finds the next line that starts a request, where the
// parser resumes after an error.
var requestStartRegex = regexp.MustCompile(`(?m)^(POST|HEAD|GET|PUT|DELETE|PATCH)`)

func (s *state) multiRequest() error {
	catch := func(e error) int {
		s.errors = append(s.errors, SyntaxError{
			ByteOffset: s.at,
			Message:    e.Error(),
		})
		if s.at >= len(s.text) {
			return -1
		}
		remain := s.text[s.at:]
		match := requestStartRegex.FindStringIndex(remain)
		if match != nil {
			return match[0] + s.at
		}
		return -1
	}
	for !s.eof() {
		if err := s.white(); err != nil {
			if next := catch(err); next >= 0 {
				if err := s.reset(next); err != nil {
					return err
				}
			} else {
				return nil
			}
		}
		if s.eof() {
			continue
		}
		if err := s.comment(); err != nil {
			if next := catch(err); next >= 0 {
				if err := s.reset(next); err != nil {
					return err
				}
			} else {
				return nil
			}
		}
		if err := s.white(); err != nil {
			if next := catch(err); next >= 0 {
				if err := s.reset(next); err != nil {
					return err
				}
			} else {
				return nil
			}
		}
		if s.eof() {
			continue
		}
		if err := s.request(); err != nil {
			if next := catch(err); next >= 0 {
				if err := s.reset(next); err != nil {
					return err
				}
			} else {
				return nil
			}
		}
		if err := s.white(); err != nil {
			if next := catch(err); next >= 0 {
				if err := s.reset(next); err != nil {
					return err
				}
			} else {
				return nil
			}
		}
	}

	return nil
}

func (s *state) request() error {
	if err := s.white(); err != nil {
		return err
	}
	if err := s.addRequestStart(); err != nil {
		return err
	}
	if _, err := s.method(); err != nil {
		return err
	}
	if err := s.updateRequestEnd(); err != nil {
		return err
	}
	if err := s.strictWhite(); err != nil {
		return err
	}
	if _, err := s.url(); err != nil {
		return err
	}
	if err := s.updateRequestEnd(); err != nil {
		return err
	}
	// advance to one new line
	if err := s.strictWhite(); err != nil {
		return err
	}
	if err := s.newLine(); err != nil {
		return err
	}
	if err := s.strictWhite(); err != nil {
		return err
	}
	if s.ch == '{' {
		if _, err := s.object(); err != nil {
			return err
		}
		if err := s.updateRequestEnd(); err != nil {
			return err
		}
	}
	// multi doc request
	// advance to one new line
	if err := s.strictWhite(); err != nil {
		return err
	}
	if err := s.newLine(); err != nil {
		return err
	}
	if err := s.strictWhite(); err != nil {
		return err
	}
	for s.ch == '{' {
		// another object
		if _, err := s.object(); err != nil {
			return err
		}
		if err := s.updateRequestEnd(); err != nil {
			return err
		}
		if err := s.strictWhite(); err != nil {
			return err
		}
		if err := s.newLine(); err != nil {
			return err
		}
		if err := s.strictWhite(); err != nil {
			return err
		}
	}
	return s.addRequestEnd()
}

func (s *state) url() (string, error) {
	var url strings.Builder
	for !s.eof() && s.ch != '\n' {
		url.WriteRune(s.ch)
		if _, err := s.nextEmptyInput(); err != nil {
			return "", err
		}
	}
	if url.Len() == 0 {
		return "", fmt.Errorf("Missing url")
	}
	return url.String(), nil
}

func (s *state) object() (map[string]any, error) {
	key := ""
	object := make(map[string]any)

	if s.ch == '{' {
		if err := s.next('{'); err != nil {
			return nil, err
		}
		if err := s.white(); err != nil {
			return nil, err
		}
		if s.ch == '}' {
			if err := s.next('}'); err != nil {
				return nil, err
			}
			// empty object
			return object, nil
		}
		for !s.eof() {
			var err error
			key, err = s.string()
			if err != nil {
				return nil, err
			}
			if err := s.white(); err != nil {
				return nil, err
			}
			if err := s.next(':'); err != nil {
				return nil, err
			}
			if _, ok := object[key]; ok {
				return nil, fmt.Errorf("duplicate key '%s'", key)
			}
			v, err := s.value()
			if err != nil {
				return nil, err
			}
			object[key] = v
			if err := s.white(); err != nil {
				return nil, err
			}
			if s.ch == '}' {
				if err := s.next('}'); err != nil {
					return nil, err
				}
				return object, nil
			}
			if err := s.next(','); err != nil {
				return nil, err
			}
			if err := s.white(); err != nil {
				return nil, err
			}
		}
	}
	return nil, fmt.Errorf("bad object")
}

func (s *state) value() (any, error) {
	if err := s.enterNesting(); err != nil {
		return nil, err
	}
	defer s.leaveNesting()
	if err := s.white(); err != nil {
		return nil, err
	}
	switch s.ch {
	case '{':
		return s.object()
	case '[':
		return s.array()
	case '"':
		return s.string()
	case '-':
		return s.number()
	default:
		if s.ch >= '0' && s.ch <= '9' {
			return s.number()
		}
		return s.word()
	}
}

func (s *state) word() (any, error) {
	switch s.ch {
	case 't':
		if err := s.next('t'); err != nil {
			return nil, err
		}
		if err := s.next('r'); err != nil {
			return nil, err
		}
		if err := s.next('u'); err != nil {
			return nil, err
		}
		if err := s.next('e'); err != nil {
			return nil, err
		}
		return true, nil
	case 'f':
		if err := s.next('f'); err != nil {
			return nil, err
		}
		if err := s.next('a'); err != nil {
			return nil, err
		}
		if err := s.next('l'); err != nil {
			return nil, err
		}
		if err := s.next('s'); err != nil {
			return nil, err
		}
		if err := s.next('e'); err != nil {
			return nil, err
		}
		return false, nil
	case 'n':
		if err := s.next('n'); err != nil {
			return nil, err
		}
		if err := s.next('u'); err != nil {
			return nil, err
		}
		if err := s.next('l'); err != nil {
			return nil, err
		}
		if err := s.next('l'); err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected '%c'", s.ch)
	}
}

func (s *state) number() (float64, error) {
	var str strings.Builder
	if s.ch == '-' {
		str.WriteByte('-')
		if err := s.next('-'); err != nil {
			return 0, err
		}
	}
	for s.ch >= '0' && s.ch <= '9' {
		str.WriteRune(s.ch)
		if _, err := s.nextEmptyInput(); err != nil {
			return 0, err
		}
	}
	if s.ch == '.' {
		str.WriteByte('.')
		for {
			if _, err := s.nextEmptyInput(); err != nil {
				return 0, err
			}
			if s.ch >= '0' && s.ch <= '9' {
				str.WriteRune(s.ch)
			} else {
				break
			}
		}
	}
	if s.ch == 'e' || s.ch == 'E' {
		str.WriteRune(s.ch)
		if _, err := s.nextEmptyInput(); err != nil {
			return 0, err
		}
		if s.ch == '+' || s.ch == '-' {
			str.WriteRune(s.ch)
			if _, err := s.nextEmptyInput(); err != nil {
				return 0, err
			}
		}
		for s.ch >= '0' && s.ch <= '9' {
			str.WriteRune(s.ch)
			if _, err := s.nextEmptyInput(); err != nil {
				return 0, err
			}
		}
	}
	num, err := strconv.ParseFloat(str.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("bad number")
	}
	return num, nil
}

func (s *state) array() ([]any, error) {
	var array []any
	if s.ch == '[' {
		if err := s.next('['); err != nil {
			return nil, err
		}
		if err := s.white(); err != nil {
			return nil, err
		}
		if s.ch == ']' {
			if err := s.next(']'); err != nil {
				return nil, err
			}
			// empty array
			return array, nil
		}
		for !s.eof() {
			v, err := s.value()
			if err != nil {
				return nil, err
			}
			array = append(array, v)
			if err := s.white(); err != nil {
				return nil, err
			}
			if s.ch == ']' {
				if err := s.next(']'); err != nil {
					return nil, err
				}
				return array, nil
			}
			if err := s.next(','); err != nil {
				return nil, err
			}
			if err := s.white(); err != nil {
				return nil, err
			}
		}
	}
	return nil, fmt.Errorf("bad array")
}

func (s *state) string() (string, error) {
	var str strings.Builder
	var uffff int32
	if s.ch == '"' {
		if s.peek(0) == '"' && s.peek(1) == '"' {
			// literal
			if err := s.next('"'); err != nil {
				return "", err
			}
			if err := s.next('"'); err != nil {
				return "", err
			}
			return s.nextUpTo(`"""`, `failed to find closing '"""'`)
		}
		for {
			if _, err := s.nextEmptyInput(); err != nil {
				return "", err
			}
			if s.eof() {
				break
			}
			if s.ch == '"' {
				if _, err := s.nextEmptyInput(); err != nil {
					return "", err
				}
				return str.String(), nil
			} else if s.ch == '\\' {
				if _, err := s.nextEmptyInput(); err != nil {
					return "", err
				}
				if s.ch == 'u' {
					uffff = 0
					for i := 0; i < 4; i++ {
						nextRune, err := s.nextEmptyInput()
						if err != nil {
							return "", err
						}
						// Parse next rune into hex.
						hex, err := strconv.ParseUint(string(nextRune), 16, 32)
						if err != nil {
							break
						}
						uffff = (uffff << 4) | int32(hex&0xF)
					}
					// Treat uffff as UTF-16 encoded rune.
					str.WriteRune(rune(uffff))
				} else if v, ok := s.escapee[s.ch]; ok {
					str.WriteString(v)
				} else {
					break
				}
			} else {
				str.WriteRune(s.ch)
			}
		}
	}

	return "", fmt.Errorf("bad string")
}

func (s *state) nextUpTo(upTo string, errorMessage string) (string, error) {
	currentAt := s.at
	i := strings.Index(s.text[s.at:], upTo)
	if i < 0 {
		if errorMessage != "" {
			return "", fmt.Errorf("%s", errorMessage)
		}
		return "", fmt.Errorf("expected '%s'", upTo)
	}
	i += currentAt
	if err := s.reset(i + len(upTo)); err != nil {
		return "", err
	}
	return s.text[currentAt:i], nil
}

func (s *state) reset(newAt int) error {
	ch, sz := utf8.DecodeRuneInString(s.text[newAt:])
	if ch == utf8.RuneError {
		if sz == 0 {
			return fmt.Errorf("unexpected empty input")
		}
		if sz == 1 {
			return fmt.Errorf("invalid UTF-8 character")
		}
		return fmt.Errorf("unknown decoding rune error")
	}
	s.ch = ch
	s.at = newAt + sz
	return nil
}

func (s *state) newLine() error {
	if s.ch == '\n' {
		if _, err := s.nextEmptyInput(); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) strictWhite() error {
	for s.ch == ' ' || s.ch == '\t' {
		if _, err := s.nextEmptyInput(); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) nextOneOf(rs []rune) error {
	if !includes(rs, s.ch) {
		return fmt.Errorf("expected one of %+v instead of '%c'", rs, s.ch)
	}
	if s.at >= len(s.text) {
		// EOF, just increase the at by 1 and set the ch to 0.
		s.at++
		s.ch = 0
		return nil
	}
	ch, sz := utf8.DecodeRuneInString(s.text[s.at:])
	if ch == utf8.RuneError {
		if sz == 0 {
			return fmt.Errorf("unexpected empty input")
		}
		if sz == 1 {
			return fmt.Errorf("invalid UTF-8 character")
		}
		return fmt.Errorf("unknown decoding rune error")
	}
	s.ch = ch
	s.at += sz
	return nil
}

func (s *state) method() (string, error) {
	uppercase := strings.ToUpper(string(s.ch))
	switch uppercase {
	case "G":
		if err := s.nextOneOf([]rune{'G', 'g'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'E', 'e'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'T', 't'}); err != nil {
			return "", err
		}
		return "GET", nil
	case "H":
		if err := s.nextOneOf([]rune{'H', 'h'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'E', 'e'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'A', 'a'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'D', 'd'}); err != nil {
			return "", err
		}
		return "HEAD", nil
	case "D":
		if err := s.nextOneOf([]rune{'D', 'd'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'E', 'e'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'L', 'l'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'E', 'e'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'T', 't'}); err != nil {
			return "", err
		}
		if err := s.nextOneOf([]rune{'E', 'e'}); err != nil {
			return "", err
		}
		return "DELETE", nil
	case "P":
		if err := s.nextOneOf([]rune{'P', 'p'}); err != nil {
			return "", err
		}
		nextUppercase := strings.ToUpper(string(s.ch))
		switch nextUppercase {
		case "A":
			if err := s.nextOneOf([]rune{'A', 'a'}); err != nil {
				return "", err
			}
			if err := s.nextOneOf([]rune{'T', 't'}); err != nil {
				return "", err
			}
			if err := s.nextOneOf([]rune{'C', 'c'}); err != nil {
				return "", err
			}
			if err := s.nextOneOf([]rune{'H', 'h'}); err != nil {
				return "", err
			}
			return "PATCH", nil
		case "U":
			if err := s.nextOneOf([]rune{'U', 'u'}); err != nil {
				return "", err
			}
			if err := s.nextOneOf([]rune{'T', 't'}); err != nil {
				return "", err
			}
			return "PUT", nil
		case "O":
			if err := s.nextOneOf([]rune{'O', 'o'}); err != nil {
				return "", err
			}
			if err := s.nextOneOf([]rune{'S', 's'}); err != nil {
				return "", err
			}
			if err := s.nextOneOf([]rune{'T', 't'}); err != nil {
				return "", err
			}
			return "POST", nil
		default:
			return "", fmt.Errorf("unexpected '%c'", s.ch)
		}
	default:
		return "", fmt.Errorf("expected one of GET/POST/PUT/DELETE/HEAD/PATCH")
	}
}

func (s *state) addRequestStart() error {
	previousRune, sz := utf8.DecodeLastRuneInString(s.text[:s.at])
	if previousRune == utf8.RuneError {
		if sz == 0 {
			return fmt.Errorf("unexpected empty input")
		}
		if sz == 1 {
			return fmt.Errorf("invalid UTF-8 character")
		}
		return fmt.Errorf("unknown decoding rune error")
	}
	s.requestStartOffset = s.at - sz
	s.requests = append(s.requests, ParsedRequest{
		StartOffset: s.requestStartOffset,
		EndOffset:   -1,
	})
	return nil
}

func (s *state) addRequestEnd() error {
	if len(s.requests) == 0 {
		return fmt.Errorf("unexpected empty requests")
	}
	s.requests[len(s.requests)-1].EndOffset = s.requestEndOffset
	return nil
}

// updateRequestEnd ends the request before the current character, as the
// upstream's at - 1 does. Past the end of the text, that is the end of the
// text. The test used to be at >= len(text), which also held while the
// current character was the last one, and a request followed by a
// multi-byte last character then ended inside it.
func (s *state) updateRequestEnd() error {
	if s.eof() {
		s.requestEndOffset = s.at - 1
		return nil
	}
	previousRune, sz := utf8.DecodeLastRuneInString(s.text[:s.at])
	if previousRune == utf8.RuneError {
		if sz == 0 {
			return fmt.Errorf("unexpected empty input")
		}
		if sz == 1 {
			return fmt.Errorf("invalid UTF-8 character")
		}
		return fmt.Errorf("unknown decoding rune error")
	}
	s.requestEndOffset = s.at - sz
	return nil
}

func (s *state) comment() error {
	for s.ch == '#' {
		for !s.eof() && s.ch != '\n' {
			if _, err := s.nextEmptyInput(); err != nil {
				return err
			}
		}
		if err := s.white(); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) peek(offset uint) rune {
	if s.at >= len(s.text) {
		return 0
	}
	tempAt := s.at
	var peekCh rune
	var sz int
	for i := uint(0); i <= offset; i++ {
		if tempAt >= len(s.text) {
			return 0
		}
		peekCh, sz = utf8.DecodeRuneInString(s.text[tempAt:])
		if peekCh == utf8.RuneError {
			return 0
		}
		tempAt += sz
	}
	return peekCh
}

func (s *state) white() error {
	for !s.eof() {
		// Skip whitespace, NUL included.
		for !s.eof() && s.ch <= ' ' {
			if _, err := s.nextEmptyInput(); err != nil {
				return err
			}
		}

		// if the current rune in iteration is '#' or the rune and the next rune is equal to '//'
		// we are on the single line comment.
		if s.ch == '#' || (s.ch == '/' && s.peek(0) == '/') {
			// Until we are on the new line, skip to the next char.
			for !s.eof() && s.ch != '\n' {
				if _, err := s.nextEmptyInput(); err != nil {
					return err
				}
			}
		} else if s.ch == '/' && s.peek(0) == '*' {
			// If the chars starts with '/*', we are on the multiline comment.
			if err := s.nNextEmptyInput(2); err != nil {
				return err
			}
			for !s.eof() && (s.ch != '*' || s.peek(0) != '/') {
				// Until we have closing tags '*', skip to the next char.
				if _, err := s.nextEmptyInput(); err != nil {
					return err
				}
			}
			if !s.eof() {
				if err := s.nNextEmptyInput(2); err != nil {
					return err
				}
			}
		} else {
			break
		}
	}

	return nil
}

func (s *state) next(c rune) error {
	if c != s.ch {
		return fmt.Errorf("expected '%c' instead of '%c'", c, s.ch)
	}

	_, err := s.nextEmptyInput()
	return err
}

func (s *state) nextEmptyInput() (rune, error) {
	if s.at >= len(s.text) {
		// EOF
		s.at++
		s.ch = 0
		return 0, nil
	}
	nextCh, sz := utf8.DecodeRuneInString(s.text[s.at:])
	if nextCh == utf8.RuneError {
		if sz == 0 {
			return 0, fmt.Errorf("unexpected empty input")
		}
		if sz == 1 {
			return 0, fmt.Errorf("invalid UTF-8 character")
		}
		return 0, fmt.Errorf("unknown decoding rune error")
	}
	s.ch = nextCh

	s.at += sz
	return nextCh, nil
}

// call s.nextEmptyInput n times.
func (s *state) nNextEmptyInput(n int) error {
	var err error
	for i := 0; i < n; i++ {
		if _, err = s.nextEmptyInput(); err != nil {
			return err
		}
	}
	return nil
}

func includes[T rune](sl []T, e T) bool {
	for _, v := range sl {
		if v == e {
			return true
		}
	}
	return false
}

// GetEditorRequest is Lines.EditorRequest for one request of text. Callers
// extracting every request of a text index it once with NewLines.
func GetEditorRequest(text string, a AdjustedParsedRequest) *EditorRequest {
	return NewLines(text).EditorRequest(a)
}

// EditorRequest extracts a single editor request from the text using adjusted line numbers.
// See https://sourcegraph.com/github.com/elastic/kibana/-/blob/src/platform/plugins/shared/console/public/application/containers/editor/utils/requests_utils.ts?L204.
func (l *Lines) EditorRequest(a AdjustedParsedRequest) *EditorRequest {
	e := &EditorRequest{}
	lines := l.lines
	methodURLLine := strings.TrimSpace(lines[a.StartLineNumber])
	if methodURLLine == "" {
		return nil
	}
	method, url := ParseLine(methodURLLine)
	if method == "" || url == "" {
		return nil
	}
	e.Method = method
	e.URL = url

	if a.EndLineNumber <= a.StartLineNumber {
		return e
	}

	dataString := ""
	if a.StartLineNumber < len(lines)-1 {
		// The lines after the method line through EndLineNumber, joined by
		// their newlines.
		dataString = strings.TrimSpace(l.text[l.starts[a.StartLineNumber+1]:l.LineEnd(a.EndLineNumber)])
	}

	data := SplitDataIntoJSONObjects(dataString)
	return &EditorRequest{
		Method: method,
		URL:    url,
		Data:   data,
	}
}

// SplitDataIntoJSONObjects splits a concatenated string of JSON objects into individual JSON objects.
// This function takes a string containing one or more JSON objects concatenated together,
// separated by optional whitespace, and splits them into an array of individual JSON strings.
// It ensures that nested objects and strings containing braces do not interfere with the splitting logic.
//
// Example inputs:
// - '{ "query": "test"} { "query": "test" }' -> ['{ "query": "test"}', '{ "query": "test" }']
// - '{ "query": "test"}' -> ['{ "query": "test"}']
// - '{ "query": "{a} {b}"}' -> ['{ "query": "{a} {b}"}'].
func SplitDataIntoJSONObjects(s string) []string {
	var jsonObjects []string
	// Track the depth of nested braces
	depth := 0
	// Holds the current JSON object as we iterate. A builder, because
	// appending each character to a string copied the whole object every time,
	// which made a long body take time quadratic in its length.
	var currentObject strings.Builder
	// Tracks whether the current position is inside a string
	insideString := false
	// The previous character, to tell an escaped double quote
	prev := rune(0)
	// Iterate through each character in the input string. An invalid UTF-8
	// byte reads as utf8.RuneError and is written as U+FFFD.
	for _, r := range s {
		// Append the character to the current JSON object string
		currentObject.WriteRune(r)

		// If the character is a double quote and it is not escaped, toggle the `insideString` state
		if r == '"' && prev != '\\' {
			insideString = !insideString
		} else if !insideString {
			// Only modify depth if not inside a string
			switch r {
			case '{':
				depth++
			case '}':
				depth--
			default:
				// Other characters don't affect depth
			}

			if depth == 0 {
				jsonObjects = append(jsonObjects, strings.TrimSpace(currentObject.String()))
				currentObject.Reset()
			}
		}
		prev = r
	}

	// If there's remaining data in currentObject, add it as the last JSON object.
	if rest := strings.TrimSpace(currentObject.String()); rest != "" {
		jsonObjects = append(jsonObjects, rest)
	}

	// Filter out any empty strings from the result
	var result []string
	for i := 0; i < len(jsonObjects); i++ {
		if jsonObjects[i] != "" {
			result = append(result, jsonObjects[i])
		}
	}
	return result
}

// ParseLine extracts the method and URL from a request line.
func ParseLine(line string) (method string, url string) {
	line = strings.TrimSpace(line)
	firstWhitespaceIndex := strings.Index(line, " ")
	if firstWhitespaceIndex < 0 {
		// There is no url, only method
		return line, ""
	}

	// 1st part is the method
	method = strings.ToUpper(strings.TrimSpace(line[0:firstWhitespaceIndex]))
	// 2nd part is the url
	url = RemoveTrailingWhitespace(strings.TrimSpace(line[firstWhitespaceIndex:]))
	return method, url
}

// RemoveTrailingWhitespace removes any trailing comments from a URL string.
// For example: "_search // comment" -> "_search"
// Ideally the parser would do that, but currently they are included in the url.
func RemoveTrailingWhitespace(s string) string {
	index := 0
	whitespaceIndex := -1
	isQueryParam := false
	for {
		r, sz := utf8.DecodeRuneInString(s[index:])
		if r == utf8.RuneError {
			break
		}
		if r == '"' {
			isQueryParam = !isQueryParam
		} else if r == ' ' && !isQueryParam {
			whitespaceIndex = index
			break
		}
		index += sz
	}
	if whitespaceIndex > 0 {
		return s[:whitespaceIndex]
	}
	return s
}

// CollapseLiteralString collapses triple-quoted literal strings in the given string.
func CollapseLiteralString(s string) string {
	splitData := strings.Split(s, `"""`)
	for idx := 1; idx < len(splitData)-1; idx += 2 {
		v, err := json.Marshal(splitData[idx])
		if err != nil {
			continue
		}
		splitData[idx] = string(v)
	}
	return strings.Join(splitData, "")
}

// IndentData round-trips a string through HJSON to strip comments, then
// re-marshals as indented JSON.
//
// It returns what json.MarshalIndent gives for the map hjson.Unmarshal
// decodes s into, but does not let hjson.Unmarshal build it: that converts
// the HJSON tree to JSON by marshaling every object once for each object it
// is nested in, which made a deeply nested body take time quadratic in its
// depth. The tree is decoded as an hjson.Node instead, which takes the same
// parse, and jsonObject converts it in one pass.
func IndentData(s string) string {
	var node hjson.Node
	// UseJSONNumber keeps numbers as hjson.Unmarshal does when it converts
	// to JSON: as json.Number, so that jsonValue can validate them the same way.
	if err := hjson.UnmarshalWithOptions([]byte(s), &node, hjson.DecoderOptions{UseJSONNumber: true}); err != nil {
		return s
	}
	var v map[string]any
	switch root := node.Value.(type) {
	case nil:
		// Decoding null into the map leaves it nil, which marshals as null.
	case *hjson.OrderedMap:
		var ok bool
		if v, ok = jsonObject(root); !ok {
			return s
		}
	default:
		// No other JSON value decodes into a map.
		return s
	}
	m, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return s
	}
	return string(m)
}

// jsonObject converts an object of an hjson.Node tree to the map
// json.Unmarshal decodes from the JSON hjson.Unmarshal would produce for it.
// It reports false where producing or decoding that JSON fails.
func jsonObject(o *hjson.OrderedMap) (map[string]any, bool) {
	m := make(map[string]any, len(o.Keys))
	for _, key := range o.Keys {
		child, ok := o.Map[key].(*hjson.Node)
		if !ok {
			return nil, false
		}
		v, ok := jsonValue(child)
		if !ok {
			return nil, false
		}
		// Keys that become equal once made valid UTF-8 decode as one key, the
		// last one written.
		m[jsonString(key)] = v
	}
	return m, true
}

// jsonValue is jsonObject for any value of an hjson.Node tree.
func jsonValue(n *hjson.Node) (any, bool) {
	switch v := n.Value.(type) {
	case *hjson.OrderedMap:
		return jsonObject(v)
	case []any:
		a := make([]any, 0, len(v))
		for _, e := range v {
			child, ok := e.(*hjson.Node)
			if !ok {
				return nil, false
			}
			ev, ok := jsonValue(child)
			if !ok {
				return nil, false
			}
			a = append(a, ev)
		}
		return a, true
	case string:
		return jsonString(v), true
	case json.Number:
		// json.Marshal rejects a number HJSON accepts and JSON does not, such
		// as "1.", and json.Unmarshal one outside the range of float64.
		if _, err := json.Marshal(v); err != nil {
			return nil, false
		}
		literal := string(v)
		if literal == "" {
			// json.Marshal writes the zero Number as 0.
			literal = "0"
		}
		f, err := strconv.ParseFloat(literal, 64)
		if err != nil {
			return nil, false
		}
		return f, true
	case bool, nil:
		return v, true
	default:
		return nil, false
	}
}

// jsonString returns s as it decodes after json.Marshal, which writes each
// byte of s that is not valid UTF-8 as U+FFFD.
func jsonString(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	return string([]rune(s))
}

// ContainsComments checks whether a string contains JSON/JS-style comments
// outside of string literals.
func ContainsComments(s string) bool {
	insideString := false
	var prevR rune
	rs := []rune(s)
	for i, r := range rs {
		nextR := rune(0)
		if i+1 < len(rs) {
			nextR = rs[i+1]
		}

		if !insideString && r == '"' {
			insideString = true
		} else if insideString && r == '"' && prevR != '\\' {
			insideString = false
		} else if !insideString {
			if r == '/' && (nextR == '/' || nextR == '*') {
				return true
			}
		}
		prevR = r
	}
	return false
}

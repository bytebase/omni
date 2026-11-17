package analysis

import "github.com/bytebase/omni/redshift"

// Range is a byte range in the script, start inclusive and end exclusive.
type Range struct {
	Start int
	End   int
}

// StatementRanges returns the byte range of each statement in a script,
// including its trailing ';'. Empty segments are skipped. Positions leave as
// byte offsets, and a caller reporting LSP positions converts them.
//
// A script that parses gives each statement's first token as its start, so the
// comments and whitespace before it are not in the range. A script the parser
// rejects falls back to the splitter's own boundaries, which keep them.
func StatementRanges(script string) []Range {
	if _, err := redshift.Parse(script); err != nil {
		return lexicalStatementRanges(script)
	}
	return parsedStatementRanges(script)
}

func parsedStatementRanges(script string) []Range {
	var ranges []Range
	for _, segment := range redshift.Split(script) {
		if segment.Empty() {
			continue
		}
		ranges = append(ranges, Range{
			Start: segment.ByteStart + leadingTriviaLen(segment.Text),
			End:   segment.ByteEnd,
		})
	}
	return ranges
}

func lexicalStatementRanges(script string) []Range {
	var ranges []Range
	for _, segment := range redshift.Split(script) {
		if segment.Empty() {
			continue
		}
		ranges = append(ranges, Range{Start: segment.ByteStart, End: segment.ByteEnd})
	}
	return ranges
}

// leadingTriviaLen returns the length of the whitespace, semicolons and
// comments a segment holds before its first statement token.
func leadingTriviaLen(statement string) int {
	i := 0
	for i < len(statement) {
		switch {
		case isTriviaSpace(statement[i]) || statement[i] == ';':
			i++
		case statement[i] == '-' && i+1 < len(statement) && statement[i+1] == '-':
			i += 2
			for i < len(statement) && statement[i] != '\n' {
				i++
			}
		case statement[i] == '/' && i+1 < len(statement) && statement[i+1] == '*':
			i += 2
			depth := 1
			for i < len(statement) && depth > 0 {
				if statement[i] == '/' && i+1 < len(statement) && statement[i+1] == '*' {
					depth++
					i += 2
				} else if statement[i] == '*' && i+1 < len(statement) && statement[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
		default:
			return i
		}
	}
	return i
}

func isTriviaSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

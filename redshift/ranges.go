package redshift

import (
	"unicode/utf16"
	"unicode/utf8"
)

// StatementRange is the source range of a SQL statement.
type StatementRange = LSPRange

// StatementRanges returns LSP UTF-16 ranges for non-empty SQL statements.
func StatementRanges(sql string) ([]StatementRange, error) {
	if _, err := Parse(sql); err != nil {
		return nil, err
	}

	var ranges []StatementRange
	for _, segment := range Split(sql) {
		if segment.Empty() {
			continue
		}
		start := segment.ByteStart + leadingTriviaLen(segment.Text)
		end := segment.ByteEnd
		ranges = append(ranges, StatementRange{
			Start: offsetToLSPPosition(sql, start),
			End:   offsetToLSPPosition(sql, end),
		})
	}
	return ranges, nil
}

func offsetToLSPPosition(s string, offset int) LSPPosition {
	if offset < 0 {
		offset = 0
	}
	if offset > len(s) {
		offset = len(s)
	}
	line := 0
	lineStart := 0
	for i := 0; i < offset; {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\n' {
			line++
			lineStart = i + size
		}
		i += size
	}
	return LSPPosition{
		Line:      line,
		Character: utf16Len(s[lineStart:offset]),
	}
}

func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		n += len(utf16.Encode([]rune{r}))
	}
	return n
}

func leadingTriviaLen(s string) int {
	i := 0
	for i < len(s) {
		switch {
		case isSpace(s[i]) || s[i] == ';':
			i++
		case s[i] == '-' && i+1 < len(s) && s[i+1] == '-':
			i += 2
			for i < len(s) && s[i] != '\n' {
				i++
			}
		case s[i] == '/' && i+1 < len(s) && s[i+1] == '*':
			i += 2
			depth := 1
			for i < len(s) && depth > 0 {
				if s[i] == '/' && i+1 < len(s) && s[i+1] == '*' {
					depth++
					i += 2
				} else if s[i] == '*' && i+1 < len(s) && s[i+1] == '/' {
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

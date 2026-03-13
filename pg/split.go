package pg

import "strings"

// Segment represents a portion of SQL text delimited by top-level semicolons.
type Segment struct {
	Text      string // the raw text of this segment
	ByteStart int    // byte offset of start in original sql
	ByteEnd   int    // byte offset of end (exclusive) in original sql
}

// Empty returns true if the segment contains only whitespace, comments, and semicolons.
func (s Segment) Empty() bool {
	t := s.Text
	i := 0
	for i < len(t) {
		b := t[i]
		// Skip whitespace and semicolons.
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == ';' {
			i++
			continue
		}
		// Skip line comments.
		if b == '-' && i+1 < len(t) && t[i+1] == '-' {
			i += 2
			for i < len(t) && t[i] != '\n' {
				i++
			}
			continue
		}
		// Skip block comments.
		if b == '/' && i+1 < len(t) && t[i+1] == '*' {
			i += 2
			depth := 1
			for i < len(t) && depth > 0 {
				if t[i] == '/' && i+1 < len(t) && t[i+1] == '*' {
					depth++
					i += 2
				} else if t[i] == '*' && i+1 < len(t) && t[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth > 0 {
				return false // unterminated block comment is not empty
			}
			continue
		}
		// Found a non-whitespace, non-comment, non-semicolon character.
		return false
	}
	return true
}

// Split splits SQL text into segments at top-level semicolons.
// It is a pure lexical scanner that does not parse SQL, so it works
// on both valid and invalid SQL. Each returned segment includes
// the terminating semicolon (if present). Segments are returned
// with their byte offsets in the original string.
func Split(sql string) []Segment {
	if len(sql) == 0 {
		return nil
	}

	var segments []Segment
	start := 0
	i := 0

	for i < len(sql) {
		b := sql[i]

		switch {
		// Single-quoted string.
		case b == '\'':
			i = skipSingleQuote(sql, i)

		// Double-quoted identifier.
		case b == '"':
			i = skipDoubleQuote(sql, i)

		// Dollar-quoted string.
		case b == '$' && isDollarQuoteStart(sql, i):
			i = skipDollarQuote(sql, i)

		// Block comment.
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockComment(sql, i)

		// Line comment.
		case b == '-' && i+1 < len(sql) && sql[i+1] == '-':
			i = skipLineComment(sql, i)

		// BEGIN ATOMIC block.
		case (b == 'b' || b == 'B') && matchKeyword(sql, i, "BEGIN") && isFollowedByAtomic(sql, i+5):
			i = skipBeginAtomic(sql, i)

		// Top-level semicolon — split here.
		case b == ';':
			i++
			segments = append(segments, Segment{
				Text:      sql[start:i],
				ByteStart: start,
				ByteEnd:   i,
			})
			start = i

		default:
			i++
		}
	}

	// Trailing content after the last semicolon.
	if start < len(sql) {
		segments = append(segments, Segment{
			Text:      sql[start:],
			ByteStart: start,
			ByteEnd:   len(sql),
		})
	}

	return segments
}

// isIdentChar returns true for [a-zA-Z0-9_].
func isIdentChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// matchKeyword checks if the keyword kw (uppercase) appears at position i
// with proper word boundaries. kw must be uppercase ASCII.
func matchKeyword(sql string, i int, kw string) bool {
	n := len(kw)
	if i+n > len(sql) {
		return false
	}
	for j := range n {
		c := sql[i+j]
		// Convert to uppercase for comparison.
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		if c != kw[j] {
			return false
		}
	}
	// Check word boundaries.
	if i > 0 && isIdentChar(sql[i-1]) {
		return false
	}
	if i+n < len(sql) && isIdentChar(sql[i+n]) {
		return false
	}
	return true
}

// skipSingleQuote skips a single-quoted string starting at position i.
// Handles '' escape. Returns position after the closing quote (or end of input).
func skipSingleQuote(sql string, i int) int {
	i++ // skip opening '
	for i < len(sql) {
		if sql[i] == '\'' {
			i++
			if i < len(sql) && sql[i] == '\'' {
				i++ // escaped ''
				continue
			}
			return i
		}
		i++
	}
	return i // unterminated
}

// skipDoubleQuote skips a double-quoted identifier starting at position i.
// Handles "" escape. Returns position after the closing quote (or end of input).
func skipDoubleQuote(sql string, i int) int {
	i++ // skip opening "
	for i < len(sql) {
		if sql[i] == '"' {
			i++
			if i < len(sql) && sql[i] == '"' {
				i++ // escaped ""
				continue
			}
			return i
		}
		i++
	}
	return i // unterminated
}

// isDollarQuoteStart checks if position i starts a valid dollar-quote tag.
// A dollar-quote is $$ or $tag$ where tag is [a-zA-Z_][a-zA-Z0-9_]*.
func isDollarQuoteStart(sql string, i int) bool {
	if i >= len(sql) || sql[i] != '$' {
		return false
	}
	j := i + 1
	if j >= len(sql) {
		return false
	}
	// $$ case
	if sql[j] == '$' {
		return true
	}
	// $tag$ case — tag must start with letter or underscore.
	if !((sql[j] >= 'a' && sql[j] <= 'z') || (sql[j] >= 'A' && sql[j] <= 'Z') || sql[j] == '_') {
		return false
	}
	j++
	for j < len(sql) && isIdentChar(sql[j]) {
		j++
	}
	return j < len(sql) && sql[j] == '$'
}

// skipDollarQuote skips a dollar-quoted string starting at position i.
// Returns position after the closing tag (or end of input).
func skipDollarQuote(sql string, i int) int {
	// Extract the tag including the $ delimiters.
	j := i + 1
	for j < len(sql) && sql[j] != '$' {
		j++
	}
	if j >= len(sql) {
		return len(sql) // unterminated (shouldn't happen if isDollarQuoteStart was true)
	}
	j++ // include closing $
	tag := sql[i:j]

	// Search for the closing tag.
	idx := strings.Index(sql[j:], tag)
	if idx < 0 {
		return len(sql) // unterminated
	}
	return j + idx + len(tag)
}

// skipBlockComment skips a block comment starting at position i.
// Supports nesting. Returns position after the closing */ (or end of input).
func skipBlockComment(sql string, i int) int {
	i += 2 // skip /*
	depth := 1
	for i < len(sql) && depth > 0 {
		if sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*' {
			depth++
			i += 2
		} else if sql[i] == '*' && i+1 < len(sql) && sql[i+1] == '/' {
			depth--
			i += 2
		} else {
			i++
		}
	}
	return i
}

// skipLineComment skips a line comment starting at position i.
// Returns position after the newline (or end of input).
func skipLineComment(sql string, i int) int {
	i += 2 // skip --
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	if i < len(sql) {
		i++ // skip the \n
	}
	return i
}

// isFollowedByAtomic checks if ATOMIC (case-insensitive, with word boundaries)
// follows after position i, skipping whitespace and comments.
func isFollowedByAtomic(sql string, i int) bool {
	i = skipWhitespaceAndComments(sql, i)
	return matchKeyword(sql, i, "ATOMIC")
}

// skipWhitespaceAndComments skips whitespace, line comments, and block comments.
func skipWhitespaceAndComments(sql string, i int) int {
	for i < len(sql) {
		b := sql[i]
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			i++
		} else if b == '-' && i+1 < len(sql) && sql[i+1] == '-' {
			i = skipLineComment(sql, i)
		} else if b == '/' && i+1 < len(sql) && sql[i+1] == '*' {
			i = skipBlockComment(sql, i)
		} else {
			break
		}
	}
	return i
}

// skipBeginAtomic skips a BEGIN ATOMIC ... END block starting at position i.
// Tracks CASE/END depth. Returns position after the closing END.
func skipBeginAtomic(sql string, i int) int {
	// Skip past "BEGIN"
	i += 5
	// Skip whitespace/comments to get past "ATOMIC"
	i = skipWhitespaceAndComments(sql, i)
	i += 6 // skip "ATOMIC"

	depth := 1
	for i < len(sql) && depth > 0 {
		b := sql[i]

		switch {
		case b == '\'':
			i = skipSingleQuote(sql, i)
		case b == '"':
			i = skipDoubleQuote(sql, i)
		case b == '$' && isDollarQuoteStart(sql, i):
			i = skipDollarQuote(sql, i)
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockComment(sql, i)
		case b == '-' && i+1 < len(sql) && sql[i+1] == '-':
			i = skipLineComment(sql, i)
		case (b == 'c' || b == 'C') && matchKeyword(sql, i, "CASE"):
			depth++
			i += 4
		case (b == 'b' || b == 'B') && matchKeyword(sql, i, "BEGIN"):
			depth++
			i += 5
		case (b == 'e' || b == 'E') && matchKeyword(sql, i, "END"):
			depth--
			i += 3
		default:
			i++
		}
	}
	return i
}

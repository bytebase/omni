package parser

// Segment represents a portion of SQL text delimited by top-level semicolons.
type Segment struct {
	Text      string // the raw text of this segment (without trailing semicolon)
	ByteStart int    // byte offset of start in original sql
	ByteEnd   int    // byte offset of end (exclusive) in original sql
}

// Empty returns true if the segment contains only whitespace and comments.
func (s Segment) Empty() bool {
	t := s.Text
	i := 0
	for i < len(t) {
		b := t[i]
		// Skip whitespace.
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			i++
			continue
		}
		// Skip -- line comments (MySQL requires space/tab/newline after --).
		if isDashComment(t, i) {
			i = skipDashComment(t, i)
			continue
		}
		// Skip # line comments.
		if b == '#' {
			i = skipHashComment(t, i)
			continue
		}
		// Skip block comments (including /*!...*/).
		if b == '/' && i+1 < len(t) && t[i+1] == '*' {
			prev := i
			i = skipBlockCommentMySQL(t, i)
			if i == prev {
				// Shouldn't happen, but guard against infinite loop.
				return false
			}
			continue
		}
		// Found a non-whitespace, non-comment character.
		return false
	}
	return true
}

// Split splits SQL text into segments at top-level semicolons.
// It is a pure lexical scanner that does not parse SQL, so it works
// on both valid and invalid SQL. Segments do NOT include the trailing
// semicolon delimiter. Empty segments (whitespace/comments only) are
// filtered out. Returns nil for empty input.
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
			i = skipSingleQuoteMySQL(sql, i)

		// Double-quoted string (MySQL treats as string literal).
		case b == '"':
			i = skipDoubleQuoteMySQL(sql, i)

		// Backtick-quoted identifier.
		case b == '`':
			i = skipBacktick(sql, i)

		// Block comment (including conditional comments /*!...*/)).
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockCommentMySQL(sql, i)

		// Dash line comment (-- followed by space/tab/newline/EOF).
		case isDashComment(sql, i):
			i = skipDashComment(sql, i)

		// Hash line comment.
		case b == '#':
			i = skipHashComment(sql, i)

		// Top-level semicolon — split here.
		case b == ';':
			seg := Segment{
				Text:      sql[start:i],
				ByteStart: start,
				ByteEnd:   i,
			}
			if !seg.Empty() {
				segments = append(segments, seg)
			}
			i++
			start = i

		default:
			i++
		}
	}

	// Trailing content after the last semicolon.
	if start < len(sql) {
		seg := Segment{
			Text:      sql[start:],
			ByteStart: start,
			ByteEnd:   len(sql),
		}
		if !seg.Empty() {
			segments = append(segments, seg)
		}
	}

	if len(segments) == 0 {
		return nil
	}
	return segments
}

// skipSingleQuoteMySQL skips a single-quoted string starting at position i.
// Handles '' escape and \ backslash escape. Returns position after the closing
// quote (or end of input if unterminated).
func skipSingleQuoteMySQL(sql string, i int) int {
	i++ // skip opening '
	for i < len(sql) {
		ch := sql[i]
		if ch == '\\' {
			i += 2 // skip backslash and escaped char
			continue
		}
		if ch == '\'' {
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

// skipDoubleQuoteMySQL skips a double-quoted string starting at position i.
// MySQL treats double-quoted strings as string literals (not identifiers).
// Handles "" escape and \ backslash escape. Returns position after the closing
// quote (or end of input if unterminated).
func skipDoubleQuoteMySQL(sql string, i int) int {
	i++ // skip opening "
	for i < len(sql) {
		ch := sql[i]
		if ch == '\\' {
			i += 2 // skip backslash and escaped char
			continue
		}
		if ch == '"' {
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

// skipBacktick skips a backtick-quoted identifier starting at position i.
// Only `` (double backtick) is an escape. Returns position after the closing
// backtick (or end of input if unterminated).
func skipBacktick(sql string, i int) int {
	i++ // skip opening `
	for i < len(sql) {
		if sql[i] == '`' {
			i++
			if i < len(sql) && sql[i] == '`' {
				i++ // escaped ``
				continue
			}
			return i
		}
		i++
	}
	return i // unterminated
}

// isDashComment returns true if position i starts a MySQL -- comment.
// MySQL requires -- to be followed by a space, tab, newline, or end-of-input.
func isDashComment(sql string, i int) bool {
	if sql[i] != '-' || i+1 >= len(sql) || sql[i+1] != '-' {
		return false
	}
	// Must be followed by whitespace or end of input.
	if i+2 >= len(sql) {
		return true
	}
	c := sql[i+2]
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// skipDashComment skips a -- line comment starting at position i.
// Returns position after the newline (or end of input).
func skipDashComment(sql string, i int) int {
	i += 2 // skip --
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	if i < len(sql) {
		i++ // skip the \n
	}
	return i
}

// skipHashComment skips a # line comment starting at position i.
// Returns position after the newline (or end of input).
func skipHashComment(sql string, i int) int {
	i++ // skip #
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	if i < len(sql) {
		i++ // skip the \n
	}
	return i
}

// skipBlockCommentMySQL skips a block comment starting at position i.
// Supports nesting. Handles both regular /* ... */ and conditional /*!...*/
// comments (for Split purposes, the entire construct is skipped).
// Returns position after the closing */ (or end of input).
func skipBlockCommentMySQL(sql string, i int) int {
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

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
		// Skip block comments, but NOT conditional comments /*!...*/.
		// Conditional comments contain executable SQL — segment is not empty.
		if b == '/' && i+1 < len(t) && t[i+1] == '*' {
			if i+2 < len(t) && t[i+2] == '!' {
				return false
			}
			prev := i
			i = skipBlockCommentMySQL(t, i)
			if i == prev {
				return false
			}
			continue
		}
		// Found a non-whitespace, non-comment character.
		return false
	}
	return true
}

// isIdentByte returns true if b is a valid identifier byte [a-zA-Z0-9_].
func isIdentByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// matchWord performs a case-insensitive keyword match at position i in sql,
// ensuring word boundaries on both sides. kw must be uppercase.
func matchWord(sql string, i int, kw string) bool {
	// Check left boundary: must be start of string or non-ident byte.
	if i > 0 && isIdentByte(sql[i-1]) {
		return false
	}
	// Check length.
	if i+len(kw) > len(sql) {
		return false
	}
	// Case-insensitive compare.
	for j := 0; j < len(kw); j++ {
		c := sql[i+j]
		// Uppercase the byte.
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		if c != kw[j] {
			return false
		}
	}
	// Check right boundary: must be end of string or non-ident byte.
	if i+len(kw) < len(sql) && isIdentByte(sql[i+len(kw)]) {
		return false
	}
	return true
}

// skipToEndOfWord advances past identifier bytes starting at position i.
func skipToEndOfWord(sql string, i int) int {
	for i < len(sql) && isIdentByte(sql[i]) {
		i++
	}
	return i
}

// skipWhitespace skips spaces, tabs, and newlines (NOT comments).
func skipWhitespace(sql string, i int) int {
	for i < len(sql) {
		b := sql[i]
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			i++
		} else {
			break
		}
	}
	return i
}

// nextWordAfter skips whitespace after pos, reads the next word, and returns it uppercase.
// Returns "" if no word follows.
func nextWordAfter(sql string, pos int) string {
	j := skipWhitespace(sql, pos)
	if j >= len(sql) || !isIdentByte(sql[j]) {
		return ""
	}
	start := j
	j = skipToEndOfWord(sql, j)
	// Uppercase the word.
	word := make([]byte, j-start)
	for k := start; k < j; k++ {
		c := sql[k]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		word[k-start] = c
	}
	return string(word)
}

// prevWord finds the last word before position i (skipping trailing whitespace backwards)
// and returns it uppercase. Returns "" if no word is found.
func prevWord(sql string, i int) string {
	// Skip whitespace backwards.
	j := i - 1
	for j >= 0 {
		b := sql[j]
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			j--
		} else {
			break
		}
	}
	if j < 0 || !isIdentByte(sql[j]) {
		return ""
	}
	end := j + 1
	for j >= 0 && isIdentByte(sql[j]) {
		j--
	}
	start := j + 1
	word := make([]byte, end-start)
	for k := start; k < end; k++ {
		c := sql[k]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		word[k-start] = c
	}
	return string(word)
}

// matchDelimiter returns true if sql[i:] starts with the given delimiter string.
func matchDelimiter(sql string, i int, delim string) bool {
	if i+len(delim) > len(sql) {
		return false
	}
	return sql[i:i+len(delim)] == delim
}

// nextNonSpaceChar returns the next non-whitespace character after pos, or 0 if none.
func nextNonSpaceChar(sql string, pos int) byte {
	j := skipWhitespace(sql, pos)
	if j >= len(sql) {
		return 0
	}
	return sql[j]
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
	stmtStart := 0
	i := 0
	depth := 0  // compound block nesting depth
	delim := ";" // active delimiter

	for i < len(sql) {
		// Check for DELIMITER directive.
		if matchWord(sql, i, "DELIMITER") {
			j := skipToEndOfWord(sql, i)
			j = skipWhitespace(sql, j)
			delimStart := j
			for j < len(sql) && sql[j] != '\n' && sql[j] != '\r' {
				j++
			}
			delimEnd := j
			for delimEnd > delimStart && (sql[delimEnd-1] == ' ' || sql[delimEnd-1] == '\t') {
				delimEnd--
			}
			if delimEnd > delimStart {
				delim = sql[delimStart:delimEnd]
			}
			if j < len(sql) && sql[j] == '\r' {
				j++
			}
			if j < len(sql) && sql[j] == '\n' {
				j++
			}
			i = j
			stmtStart = i
			continue
		}

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

		// BEGIN — increment depth unless it's a transaction (BEGIN WORK, BEGIN alone, XA BEGIN).
		// Note: "BEGIN" without semicolon followed by another statement is invalid MySQL
		// syntax, so we don't need to handle that case.
		case (b == 'b' || b == 'B') && matchWord(sql, i, "BEGIN"):
			endOfWord := skipToEndOfWord(sql, i)
			next := nextWordAfter(sql, endOfWord)
			prev := prevWord(sql, i)
			// BEGIN WORK => transaction, not compound.
			// BEGIN at EOF or followed by ; => transaction.
			// XA BEGIN => transaction.
			if next != "WORK" && prev != "XA" && nextNonSpaceChar(sql, endOfWord) != ';' && nextNonSpaceChar(sql, endOfWord) != 0 {
				depth++
			}
			i = endOfWord

		// IF — increment depth unless preceded by END, followed by EXISTS, or followed by '('.
		case (b == 'i' || b == 'I') && matchWord(sql, i, "IF"):
			endOfWord := skipToEndOfWord(sql, i)
			prev := prevWord(sql, i)
			next := nextWordAfter(sql, endOfWord)
			// Treat "IF [NOT] EXISTS" as a DDL modifier only when EXISTS is
			// followed by a bare identifier (e.g. DROP TABLE IF EXISTS t).
			// When EXISTS is followed by '(' (e.g. IF EXISTS (subquery) THEN),
			// it is the EXISTS subquery predicate used inside a compound IF.
			isExistsClause := false
			if next == "EXISTS" {
				existsEnd := skipToEndOfWord(sql, skipWhitespace(sql, endOfWord))
				if nextNonSpaceChar(sql, existsEnd) != '(' {
					isExistsClause = true
				}
			} else if next == "NOT" {
				notEnd := skipToEndOfWord(sql, skipWhitespace(sql, endOfWord))
				if nextWordAfter(sql, notEnd) == "EXISTS" {
					existsEnd := skipToEndOfWord(sql, skipWhitespace(sql, notEnd))
					if nextNonSpaceChar(sql, existsEnd) != '(' {
						isExistsClause = true
					}
				}
			}
			if prev != "END" && !isExistsClause && nextNonSpaceChar(sql, endOfWord) != '(' {
				depth++
			}
			i = endOfWord

		// CASE — increment depth unless preceded by END.
		case (b == 'c' || b == 'C') && matchWord(sql, i, "CASE"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			i = endOfWord

		// WHILE — increment depth unless preceded by END.
		case (b == 'w' || b == 'W') && matchWord(sql, i, "WHILE"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			i = endOfWord

		// LOOP — increment depth unless preceded by END.
		case (b == 'l' || b == 'L') && matchWord(sql, i, "LOOP"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			i = endOfWord

		// REPEAT — increment depth unless preceded by END or followed by '('.
		case (b == 'r' || b == 'R') && matchWord(sql, i, "REPEAT"):
			endOfWord := skipToEndOfWord(sql, i)
			prev := prevWord(sql, i)
			if prev != "END" && nextNonSpaceChar(sql, endOfWord) != '(' {
				depth++
			}
			i = endOfWord

		// END — decrement depth (if > 0), skip if preceded by XA.
		case (b == 'e' || b == 'E') && matchWord(sql, i, "END"):
			endOfWord := skipToEndOfWord(sql, i)
			if depth > 0 && prevWord(sql, i) != "XA" {
				depth--
			}
			i = endOfWord

		default:
			// Check for delimiter match at current position (only at top level).
			if depth == 0 && matchDelimiter(sql, i, delim) {
				seg := Segment{
					Text:      sql[stmtStart:i],
					ByteStart: stmtStart,
					ByteEnd:   i,
				}
				if !seg.Empty() {
					segments = append(segments, seg)
				}
				i += len(delim)
				stmtStart = i
			} else {
				i++
			}
		}
	}

	// Trailing content after the last delimiter.
	if stmtStart < len(sql) {
		seg := Segment{
			Text:      sql[stmtStart:],
			ByteStart: stmtStart,
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

// findCompoundBodyEnd scans sql starting at start and returns the byte offset
// where a routine/trigger/event body ends: the first top-level ';' (at
// compound-block depth 0) or end of input. It tracks nested compound scopes
// opened by BEGIN/IF/CASE/WHILE/LOOP/REPEAT and closed by END, using the same
// heuristics as Split so that END IF, END CASE, etc. correctly balance against
// their matching opener instead of prematurely closing an outer BEGIN...END.
//
// This is used by parseCreateFunctionStmt/parseCreateTriggerStmt/
// parseCreateEventStmt to delimit the raw body text while the compound body
// itself is not yet parsed into an AST.
func findCompoundBodyEnd(sql string, start int) int {
	i := start
	depth := 0
	for i < len(sql) {
		b := sql[i]
		switch {
		case b == '\'':
			i = skipSingleQuoteMySQL(sql, i)
		case b == '"':
			i = skipDoubleQuoteMySQL(sql, i)
		case b == '`':
			i = skipBacktick(sql, i)
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockCommentMySQL(sql, i)
		case isDashComment(sql, i):
			i = skipDashComment(sql, i)
		case b == '#':
			i = skipHashComment(sql, i)

		case (b == 'b' || b == 'B') && matchWord(sql, i, "BEGIN"):
			endOfWord := skipToEndOfWord(sql, i)
			next := nextWordAfter(sql, endOfWord)
			prev := prevWord(sql, i)
			// BEGIN WORK / XA BEGIN / lone BEGIN => transaction, not compound.
			if next != "WORK" && prev != "XA" && nextNonSpaceChar(sql, endOfWord) != ';' && nextNonSpaceChar(sql, endOfWord) != 0 {
				depth++
			}
			i = endOfWord

		case (b == 'i' || b == 'I') && matchWord(sql, i, "IF"):
			endOfWord := skipToEndOfWord(sql, i)
			prev := prevWord(sql, i)
			// Unlike Split's top-level scanner, we intentionally do NOT skip
			// IF when followed by EXISTS / NOT EXISTS. Inside a routine body
			// "IF EXISTS (subquery) THEN ... END IF" is valid flow control,
			// and any DDL IF EXISTS inside the body is wrapped by BEGIN/END
			// so depth stays balanced regardless of whether we track this IF.
			if prev != "END" && nextNonSpaceChar(sql, endOfWord) != '(' {
				depth++
			}
			i = endOfWord

		case (b == 'c' || b == 'C') && matchWord(sql, i, "CASE"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			i = endOfWord

		case (b == 'w' || b == 'W') && matchWord(sql, i, "WHILE"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			i = endOfWord

		case (b == 'l' || b == 'L') && matchWord(sql, i, "LOOP"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			i = endOfWord

		case (b == 'r' || b == 'R') && matchWord(sql, i, "REPEAT"):
			endOfWord := skipToEndOfWord(sql, i)
			prev := prevWord(sql, i)
			if prev != "END" && nextNonSpaceChar(sql, endOfWord) != '(' {
				depth++
			}
			i = endOfWord

		case (b == 'e' || b == 'E') && matchWord(sql, i, "END"):
			endOfWord := skipToEndOfWord(sql, i)
			if depth > 0 && prevWord(sql, i) != "XA" {
				depth--
			}
			i = endOfWord

		default:
			if depth == 0 && b == ';' {
				return i
			}
			i++
		}
	}
	return i
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

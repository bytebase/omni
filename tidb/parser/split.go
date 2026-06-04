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
			// TiDB-specific comments also contain executable SQL.
			if i+3 < len(t) && t[i+2] == 'T' && t[i+3] == '!' {
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

// nextSignificantChar returns the next character after pos, skipping whitespace
// and comments, or 0 if none. Used for keyword lookahead (e.g. distinguishing
// the IF(...) / REPEAT(...) functions from the compound IF / REPEAT statements).
func nextSignificantChar(sql string, pos int) byte {
	j := skipSpaceAndComments(sql, pos)
	if j >= len(sql) {
		return 0
	}
	return sql[j]
}

// skipSpaceAndComments advances past whitespace and comments (block /* */, line
// -- and #) starting at position i. A non-empty conditional comment
// (/*!NNNNN sql */, /*T! sql */) holds executable SQL and is significant — the
// scan stops at it (matching Segment.Empty). An EMPTY conditional comment
// (/*!50000*/, a bare version gate) is a no-op and is skipped like whitespace,
// so e.g. `IF /*!50000*/ (...)` is still recognised as the IF(...) function.
// Conditional comments remain opaque to Split's main scan; this only affects
// keyword lookahead.
func skipSpaceAndComments(sql string, i int) int {
	for i < len(sql) {
		switch {
		case sql[i] == ' ' || sql[i] == '\t' || sql[i] == '\n' || sql[i] == '\r':
			i++
		case sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*':
			if end, empty, ok := conditionalCommentEnd(sql, i); ok {
				if !empty {
					return i // executable SQL — significant.
				}
				i = end // bare version gate — skip like whitespace.
			} else {
				i = skipBlockCommentMySQL(sql, i)
			}
		case isDashComment(sql, i):
			i = skipDashComment(sql, i)
		case sql[i] == '#':
			i = skipHashComment(sql, i)
		default:
			return i
		}
	}
	return i
}

// conditionalCommentEnd reports whether a conditional comment (/*!NNNNN...*/ or
// /*T!...*/) begins at i. If so it returns end = the index just past the closing
// */ and empty = whether the inner content (after the /*!NNNNN or /*T! marker)
// is only whitespace. ok is false for a regular block comment.
func conditionalCommentEnd(sql string, i int) (end int, empty bool, ok bool) {
	if i+2 >= len(sql) || sql[i] != '/' || sql[i+1] != '*' {
		return 0, false, false
	}
	var innerStart int
	switch {
	case sql[i+2] == '!':
		innerStart = i + 3
		for innerStart < len(sql) && sql[innerStart] >= '0' && sql[innerStart] <= '9' {
			innerStart++
		}
	case i+3 < len(sql) && sql[i+2] == 'T' && sql[i+3] == '!':
		innerStart = i + 4
	default:
		return 0, false, false
	}
	end = skipBlockCommentMySQL(sql, i)
	innerEnd := end
	if innerEnd-2 >= innerStart && sql[innerEnd-2] == '*' && sql[innerEnd-1] == '/' {
		innerEnd -= 2
	}
	for j := innerStart; j < innerEnd; j++ {
		if sql[j] != ' ' && sql[j] != '\t' && sql[j] != '\n' && sql[j] != '\r' {
			return end, false, true
		}
	}
	return end, true, true
}

// isIfExistsModifier reports whether the IF whose word ends at ifEnd is the DDL
// modifier "IF [NOT] EXISTS <identifier>" (e.g. DROP TABLE IF EXISTS t) rather
// than a compound IF statement. It is comment-aware: "IF EXISTS /*c*/ (subquery)
// THEN" is a compound IF (the EXISTS subquery predicate), not a DDL modifier,
// because EXISTS is followed by '('. Shared by Split and findCompoundBodyEnd so
// they classify IF the same way.
func isIfExistsModifier(sql string, ifEnd int) bool {
	p := skipSpaceAndComments(sql, ifEnd)
	if matchWord(sql, p, "NOT") {
		p = skipSpaceAndComments(sql, skipToEndOfWord(sql, p))
	}
	if !matchWord(sql, p, "EXISTS") {
		return false
	}
	p = skipSpaceAndComments(sql, skipToEndOfWord(sql, p))
	return p >= len(sql) || sql[p] != '('
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
			if next != "WORK" && prev != "XA" && nextSignificantChar(sql, endOfWord) != ';' && nextSignificantChar(sql, endOfWord) != 0 {
				depth++
			}
			i = endOfWord

		// IF — increment depth unless preceded by END, an IF [NOT] EXISTS DDL
		// modifier, or an IF(...) function call.
		case (b == 'i' || b == 'I') && matchWord(sql, i, "IF"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" && !isIfExistsModifier(sql, endOfWord) && nextSignificantChar(sql, endOfWord) != '(' {
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
			if prev != "END" && nextSignificantChar(sql, endOfWord) != '(' {
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

// findCompoundBodyEnd scans a routine/trigger/event body starting at byte offset
// start and returns the offset of the first top-level (depth 0) ';', or len(sql)
// if none is found. It balances every compound opener (BEGIN/IF/CASE/WHILE/LOOP/
// REPEAT) against END using the same heuristics as Split, so an inner block
// terminator like END IF / END CASE / END WHILE / END LOOP / END REPEAT does not
// prematurely close the enclosing BEGIN block. This is the body-capture analogue
// of Split's depth tracking; a naive BEGIN/END-only counter mis-reads the END of
// an inner block as closing the routine and truncates the body.
func findCompoundBodyEnd(sql string, start int) int {
	i := start
	depth := 0
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

		// Block comment.
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockCommentMySQL(sql, i)

		// Dash line comment.
		case isDashComment(sql, i):
			i = skipDashComment(sql, i)

		// Hash line comment.
		case b == '#':
			i = skipHashComment(sql, i)

		// BEGIN — increment depth unless it's a transaction (BEGIN WORK, lone
		// BEGIN at ';'/EOF, XA BEGIN).
		case (b == 'b' || b == 'B') && matchWord(sql, i, "BEGIN"):
			endOfWord := skipToEndOfWord(sql, i)
			next := nextWordAfter(sql, endOfWord)
			prev := prevWord(sql, i)
			if next != "WORK" && prev != "XA" && nextSignificantChar(sql, endOfWord) != ';' && nextSignificantChar(sql, endOfWord) != 0 {
				depth++
			}
			i = endOfWord

		// IF — same classification as Split: count it as a compound opener
		// unless preceded by END, used as the IF(...) function (next char '('),
		// or a DDL "IF [NOT] EXISTS <ident>" modifier. "IF EXISTS (subquery)
		// THEN ... END IF" flow control is counted (EXISTS followed by '('); a
		// body-internal "DROP ... IF EXISTS t" is not (so it does not push depth
		// and swallow a following statement in a custom-delimiter segment).
		case (b == 'i' || b == 'I') && matchWord(sql, i, "IF"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" && !isIfExistsModifier(sql, endOfWord) && nextSignificantChar(sql, endOfWord) != '(' {
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
			if prev != "END" && nextSignificantChar(sql, endOfWord) != '(' {
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

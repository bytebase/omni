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
	depth := 0   // compound block nesting depth
	delim := ";" // active delimiter
	// atStmtStart tracks whether the next keyword/token would appear at a
	// statement-start position (after ';', after BEGIN/THEN/ELSE/ELSEIF/DO/
	// LOOP/REPEAT/UNTIL, after a label ':', at start of input, or across
	// whitespace/comments). Used to disambiguate IF-as-compound from
	// IF-as-function-call / DDL-modifier — mirrors MySQL's yacc grammar-
	// position disambiguation at the text-scan level.
	atStmtStart := true

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
			atStmtStart = true
			continue
		}

		b := sql[i]

		switch {
		// Single-quoted string — non-stmt-start after.
		case b == '\'':
			i = skipSingleQuoteMySQL(sql, i)
			atStmtStart = false

		// Double-quoted string (MySQL treats as string literal).
		case b == '"':
			i = skipDoubleQuoteMySQL(sql, i)
			atStmtStart = false

		// Backtick-quoted identifier.
		case b == '`':
			i = skipBacktick(sql, i)
			atStmtStart = false

		// Block comment: transparent to atStmtStart.
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockCommentMySQL(sql, i)

		// Dash line comment: transparent to atStmtStart.
		case isDashComment(sql, i):
			i = skipDashComment(sql, i)

		// Hash line comment: transparent to atStmtStart.
		case b == '#':
			i = skipHashComment(sql, i)

		// BEGIN — compound block opener (depth++) unless it's a transaction
		// (BEGIN WORK, lone BEGIN, XA BEGIN). After a compound BEGIN, the
		// next token is at stmt-start; after a transaction BEGIN, it's not.
		case (b == 'b' || b == 'B') && matchWord(sql, i, "BEGIN"):
			endOfWord := skipToEndOfWord(sql, i)
			next := nextWordAfter(sql, endOfWord)
			prev := prevWord(sql, i)
			isCompound := next != "WORK" && prev != "XA" && nextNonSpaceChar(sql, endOfWord) != ';' && nextNonSpaceChar(sql, endOfWord) != 0
			if isCompound {
				depth++
				atStmtStart = true
			} else {
				atStmtStart = false
			}
			i = endOfWord

		// IF — compound flow-control iff at stmt-start and not preceded by END.
		// atStmtStart encodes MySQL's yacc grammar-position disambiguation:
		// IF in expression context (after '=', ',', '(', keywords like DROP/
		// CREATE/SELECT/SET) is a function or DDL modifier, not a compound.
		case (b == 'i' || b == 'I') && matchWord(sql, i, "IF"):
			endOfWord := skipToEndOfWord(sql, i)
			prev := prevWord(sql, i)
			if prev != "END" && atStmtStart {
				depth++
			}
			// After IF the next tokens form either the condition expression
			// (compound) or the argument list (function) — either way, the
			// immediately following content is expression context.
			atStmtStart = false
			i = endOfWord

		// CASE — always paired with END (expr form) or END CASE (stmt form).
		case (b == 'c' || b == 'C') && matchWord(sql, i, "CASE"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			atStmtStart = false
			i = endOfWord

		// WHILE / LOOP — opener unless preceded by END (end-WHILE / end-LOOP).
		case (b == 'w' || b == 'W') && matchWord(sql, i, "WHILE"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			atStmtStart = false
			i = endOfWord

		case (b == 'l' || b == 'L') && matchWord(sql, i, "LOOP"):
			endOfWord := skipToEndOfWord(sql, i)
			if prevWord(sql, i) != "END" {
				depth++
			}
			// LOOP keyword introduces the loop body (stmt-start), unless
			// this is END LOOP (closer) or the `LOOP(...)` function (rare).
			if prevWord(sql, i) != "END" && nextNonSpaceChar(sql, endOfWord) != '(' {
				atStmtStart = true
			} else {
				atStmtStart = false
			}
			i = endOfWord

		// REPEAT — opener unless preceded by END or followed by '(' (function).
		case (b == 'r' || b == 'R') && matchWord(sql, i, "REPEAT"):
			endOfWord := skipToEndOfWord(sql, i)
			prev := prevWord(sql, i)
			if prev != "END" && nextNonSpaceChar(sql, endOfWord) != '(' {
				depth++
				atStmtStart = true
			} else {
				atStmtStart = false
			}
			i = endOfWord

		// END — close a compound opener. atStmtStart stays false because the
		// token after END is either an end-label, IF/CASE/WHILE/LOOP/REPEAT
		// (the compound-close suffix), or the next token in a statement.
		case (b == 'e' || b == 'E') && matchWord(sql, i, "END"):
			endOfWord := skipToEndOfWord(sql, i)
			if depth > 0 && prevWord(sql, i) != "XA" {
				depth--
			}
			atStmtStart = false
			i = endOfWord

		// THEN / ELSE / ELSEIF / DO / UNTIL — these keywords introduce
		// statement lists (IF branches, WHILE loop body, REPEAT predicate).
		// Explicitly mark atStmtStart=true after consuming.
		case (b == 't' || b == 'T') && matchWord(sql, i, "THEN"):
			i = skipToEndOfWord(sql, i)
			atStmtStart = true
		case (b == 'e' || b == 'E') && matchWord(sql, i, "ELSE"):
			i = skipToEndOfWord(sql, i)
			atStmtStart = true
		case (b == 'e' || b == 'E') && matchWord(sql, i, "ELSEIF"):
			i = skipToEndOfWord(sql, i)
			atStmtStart = false // ELSEIF is followed by its own condition expr
		case (b == 'd' || b == 'D') && matchWord(sql, i, "DO"):
			i = skipToEndOfWord(sql, i)
			atStmtStart = true
		case (b == 'u' || b == 'U') && matchWord(sql, i, "UNTIL"):
			i = skipToEndOfWord(sql, i)
			atStmtStart = false // UNTIL is followed by its predicate expr

		default:
			// Delimiter match — segment split at top level.
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
				atStmtStart = true
				continue
			}
			// Single character classification.
			switch b {
			case ';':
				atStmtStart = true
			case ':':
				atStmtStart = true
			case ' ', '\t', '\r', '\n':
				// whitespace — transparent to atStmtStart
			default:
				atStmtStart = false
			}
			i++
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

package pg

import (
	"strings"

	"github.com/bytebase/omni/pg/internal/copyscan"
	"github.com/bytebase/omni/pg/internal/metacmd"
)

// Segment represents a portion of SQL text delimited by top-level semicolons.
type Segment struct {
	Text      string // the raw text of this segment
	ByteStart int    // byte offset of start in original sql
	ByteEnd   int    // byte offset of end (exclusive) in original sql
}

// Empty returns true if the segment contains only whitespace, comments,
// semicolons, and psql metacommand lines (which are not SQL statements).
func (s Segment) Empty() bool {
	t := s.Text
	i := 0
	if strings.HasPrefix(t, "\xEF\xBB\xBF") {
		i = 3
	}
	for i < len(t) {
		b := t[i]
		// psql metacommand line: statement-less.
		if metacmd.IsMetaCommand(t, i) {
			i = metacmd.SkipLine(t, i)
			continue
		}
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
	// A UTF-8 BOM at the start of the script is trivia: psql strips it
	// before scanning (mainloop.c). The bytes stay in the first segment.
	if strings.HasPrefix(sql, "\xEF\xBB\xBF") {
		i = 3
	}
	// Statement-separating semicolons only occur at parenthesis depth zero:
	// PG's grammar nests semicolons inside parentheses (CREATE RULE's
	// multi-action list), and psqlscan likewise never splits inside parens.
	parenDepth := 0

	for i < len(sql) {
		b := sql[i]

		switch {
		// psql metacommand line (\restrict, \connect, ...): not SQL — the
		// line is consumed like a comment, so its content never affects
		// boundaries, and it rides along as trivia of the segment it sits
		// in. pg_dump plain format has emitted \restrict/\unrestrict since
		// the CVE-2025-8714 point releases. Recognition requires a true
		// preceding newline: a segment that merely begins with a backslash
		// (after a mid-line semicolon boundary) is not reinterpreted, which
		// keeps re-splitting stable. A metacommand at byte zero of a script
		// is glued here but skipped by the parser, so the pipeline still
		// works on header-stripped dumps.
		case metacmd.IsMetaCommand(sql, i):
			i = metacmd.SkipLine(sql, i)

		// Single-quoted string. E'...' (escape string) processes
		// backslash escapes; plain '...' does not (PG default
		// standard_conforming_strings=on).
		case b == '\'':
			if isEscapeStringQuote(sql, i) {
				i = skipEscapeString(sql, i)
			} else {
				i = skipSingleQuote(sql, i)
			}

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

		case b == '(':
			parenDepth++
			i++

		case b == ')':
			// Clamp at zero: a stray ')' does not open a "negative" group,
			// matching psqlscan (paren_depth only decremented when > 0).
			if parenDepth > 0 {
				parenDepth--
			}
			i++

		// Top-level semicolon — split here. Inside parentheses the
		// semicolon is part of the statement; an unclosed '(' therefore
		// leaves the remainder as a single segment, same as psql
		// buffering to end of input.
		case b == ';' && parenDepth == 0:
			i++
			// COPY ... FROM STDIN is followed by inline data lines that end
			// at a line containing only "\.": psql scripts and pg_dump
			// plain-format output carry the data in the SQL stream, and the
			// data may contain semicolons that are not boundaries. Keep the
			// statement, its data, and the terminator as one segment.
			if isCopyFromStdin(sql[start:i]) && copyscan.RestOfLineBlank(sql, i) {
				i = copyscan.SkipData(sql, i)
			}
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
	// Check word boundaries with scan.l identifier bytes: '$' and
	// multibyte bytes continue identifiers (x$BEGIN is one identifier),
	// so they block a keyword match just like letters do.
	if i > 0 && isIdentByte(sql[i-1]) {
		return false
	}
	if i+n < len(sql) && isIdentByte(sql[i+n]) {
		return false
	}
	return true
}

// skipSingleQuote skips a single-quoted string starting at position i.
// Handles the doubled-quote escape. Returns position after the closing quote (or end of input).
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

// isEscapeStringQuote reports whether the quote at position i opens an
// E'...' escape string: the quote is directly preceded by a lone E/e
// that is not the tail of an identifier (scan.l: {xestart} = [eE]{quote},
// only recognized when E starts a token).
func isEscapeStringQuote(sql string, i int) bool {
	if i == 0 {
		return false
	}
	prev := sql[i-1]
	if prev != 'e' && prev != 'E' {
		return false
	}
	if i >= 2 && isIdentByte(sql[i-2]) {
		return false // E is the tail of an identifier like abcE'...'
	}
	return true
}

// isIdentByte reports whether b can appear in a PG identifier
// (letters, digits, underscore, dollar, or any multibyte byte).
func isIdentByte(b byte) bool {
	return b == '_' || b == '$' ||
		(b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') || b >= 0x80
}

// skipEscapeString skips an E'...' string body starting at the opening
// quote. Unlike plain strings, backslash escapes the next character
// (including \' and \\); a doubled single quote is also honored. Returns the
// position after the closing quote, or len(sql) if unterminated.
func skipEscapeString(sql string, i int) int {
	i++ // skip opening '
	for i < len(sql) {
		switch sql[i] {
		case '\\':
			if i+1 >= len(sql) {
				return len(sql) // trailing backslash at EOF
			}
			i += 2
		case '\'':
			i++
			if i < len(sql) && sql[i] == '\'' {
				i++ // escaped ''
				continue
			}
			return i
		default:
			i++
		}
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
// A dollar-quote is $$ or $tag$ where tag is [a-zA-Z_][a-zA-Z0-9_]*, and the
// '$' must not continue a preceding identifier.
func isDollarQuoteStart(sql string, i int) bool {
	if i >= len(sql) || sql[i] != '$' {
		return false
	}
	// scan.l: '$' is an identifier-continuation byte ({ident_cont} includes
	// '$'), so abc$tag$y is a single identifier, not "abc" followed by a
	// dollar-quote. An opening delimiter is only recognized when the '$'
	// starts a new token. (Known documented divergence: PG lexes tokens, so
	// 123$t$...$t$ is number+string there but stays unsplit here; adjacent
	// number/param + string is never valid grammar, so the difference only
	// moves the boundary of an error — accepted in the splitter audit.)
	if i > 0 && isIdentByte(sql[i-1]) {
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
	// Keyword occurrences only count when they appear in statement position:
	// dot-qualified references (t.end, t . case) and AS aliases (AS end) use
	// reserved words as column labels, and bare unreserved BEGIN can be a
	// plain column reference — none of them open or close a block
	// (engine-verified on PostgreSQL 17). Track just enough left context to
	// tell the difference; whitespace and comments preserve it.
	afterDot := false
	afterAS := false
	for i < len(sql) && depth > 0 {
		b := sql[i]

		switch {
		// Metacommand lines are trivia here too: psql executes them
		// mid-buffer and the body continues, so their words (\echo END)
		// must not count as block delimiters. Context is preserved, same
		// as comments.
		case metacmd.IsMetaCommand(sql, i):
			i = metacmd.SkipLine(sql, i)
		case b == ' ' || b == '\t' || b == '\n' || b == '\r':
			i++
		case b == '\'':
			if isEscapeStringQuote(sql, i) {
				i = skipEscapeString(sql, i)
			} else {
				i = skipSingleQuote(sql, i)
			}
			afterDot, afterAS = false, false
		case b == '"':
			i = skipDoubleQuote(sql, i)
			afterDot, afterAS = false, false
		case b == '$' && isDollarQuoteStart(sql, i):
			i = skipDollarQuote(sql, i)
			afterDot, afterAS = false, false
		case b == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i = skipBlockComment(sql, i)
		case b == '-' && i+1 < len(sql) && sql[i+1] == '-':
			i = skipLineComment(sql, i)
		case b == '.':
			afterDot = true
			afterAS = false
			i++
		case isIdentByte(b):
			j := i
			for j < len(sql) && isIdentByte(sql[j]) {
				j++
			}
			word := strings.ToUpper(sql[i:j])
			if !afterDot && !afterAS {
				switch word {
				case "CASE":
					depth++
				case "BEGIN":
					// Only a nested BEGIN ATOMIC opens a block; bare BEGIN
					// is an unreserved keyword usable as an identifier.
					if isFollowedByAtomic(sql, j) {
						depth++
					}
				case "END":
					depth--
				}
			}
			afterAS = word == "AS"
			afterDot = false
			i = j
		default:
			afterDot, afterAS = false, false
			i++
		}
	}
	return i
}

// isCopyFromStdin reports whether the statement text is a COPY command
// reading inline data: the first word is COPY and the words FROM STDIN
// appear in sequence at parenthesis depth zero (so a relation named
// "stdin" inside a COPY (query) form does not match). Word scanning skips
// strings, comments, and dollar-quotes with the same helpers as Split.
func isCopyFromStdin(stmt string) bool {
	i := 0
	if strings.HasPrefix(stmt, "\xEF\xBB\xBF") {
		i = 3 // leading UTF-8 BOM is trivia, same as in Split
	}
	depth := 0
	sawCopy := false
	prevWord := ""
	for i < len(stmt) {
		b := stmt[i]
		switch {
		// Metacommand lines are trivia (same as in Split): a statement
		// segment may carry leading \connect/\restrict lines — pg_dumpall
		// emits \connect right before each database's statements — and they
		// must not shadow COPY as the first word. The line-start rule is the
		// same strict one the scan loop uses (a true preceding newline), so
		// splitting a segment's text again reproduces the same decision.
		case metacmd.IsMetaCommand(stmt, i):
			i = metacmd.SkipLine(stmt, i)
		case b == '\'':
			if isEscapeStringQuote(stmt, i) {
				i = skipEscapeString(stmt, i)
			} else {
				i = skipSingleQuote(stmt, i)
			}
		case b == '"':
			i = skipDoubleQuote(stmt, i)
		case b == '$' && isDollarQuoteStart(stmt, i):
			i = skipDollarQuote(stmt, i)
		case b == '/' && i+1 < len(stmt) && stmt[i+1] == '*':
			i = skipBlockComment(stmt, i)
		case b == '-' && i+1 < len(stmt) && stmt[i+1] == '-':
			i = skipLineComment(stmt, i)
		case b == '(':
			depth++
			i++
		case b == ')':
			if depth > 0 {
				depth--
			}
			i++
		case isIdentByte(b):
			j := i
			for j < len(stmt) && isIdentByte(stmt[j]) {
				j++
			}
			word := strings.ToUpper(stmt[i:j])
			if !sawCopy {
				// The statement must start with COPY.
				if word != "COPY" {
					return false
				}
				sawCopy = true
			} else if depth == 0 && (word == "STDIN" || word == "STDOUT") && prevWord == "FROM" {
				// gram.y maps both STDIN and STDOUT to a NULL filename, and
				// the server runs copy-in for either spelling of COPY FROM
				// (engine-verified: inline data loads via FROM STDOUT).
				return true
			}
			if depth == 0 {
				prevWord = word
			}
			i = j
		default:
			i++
		}
	}
	return false
}

// StripTopLevelMetacommands returns text with psql metacommand lines
// removed, but only where a backslash sits in top-level scan state — not
// inside a string, comment, dollar-quote, or COPY data. It walks with
// the same construct-skipping logic as Split, so a backslash that is
// really string content (e.g. 'a \restrict b') is preserved. Used to
// derive the SQL a segment would actually send to the server, since the
// database has no notion of psql backslash commands.
func StripTopLevelMetacommands(text string) string {
	var b strings.Builder
	i := 0
	start := 0
	for i < len(text) {
		c := text[i]
		switch {
		case metacmd.IsMetaCommand(text, i):
			b.WriteString(text[start:i])
			i = metacmd.SkipLine(text, i)
			start = i
		case c == '\'':
			if isEscapeStringQuote(text, i) {
				i = skipEscapeString(text, i)
			} else {
				i = skipSingleQuote(text, i)
			}
		case c == '"':
			i = skipDoubleQuote(text, i)
		case c == '$' && isDollarQuoteStart(text, i):
			i = skipDollarQuote(text, i)
		case c == '/' && i+1 < len(text) && text[i+1] == '*':
			i = skipBlockComment(text, i)
		case c == '-' && i+1 < len(text) && text[i+1] == '-':
			i = skipLineComment(text, i)
		case c == ';':
			i++
			// Same data-mode gate as Split: same-line SQL after the COPY
			// statement keeps data mode off, so stripping and splitting
			// walk the bytes identically.
			if isCopyFromStdin(text[start:i]) && copyscan.RestOfLineBlank(text, i) {
				i = copyscan.SkipData(text, i)
			}
		default:
			i++
		}
	}
	b.WriteString(text[start:])
	return b.String()
}

package parser

import "strings"

// embeddedSourceHead recognizes the head of CREATE ... JAVA ... AS source_char
// and of CREATE ... MLE MODULE ... AS module_text. The text after that AS is
// Java or JavaScript source, not SQL: it holds ';', quotes, and comments the
// SQL lexer would misread. The splitter and the parser take it verbatim
// instead of lexing it, as SQL*Plus does: SQL*Plus buffers either statement
// until a line holding only "/" (checked on Oracle 23ai for both).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/sqlrf/CREATE-JAVA.html
//
//	CREATE [ OR REPLACE ] [ AND { RESOLVE | COMPILE } ] [ NOFORCE ]
//	    JAVA { { SOURCE | RESOURCE } [ IF NOT EXISTS ] NAMED [ schema. ] primary_name
//	         | CLASS [ IF NOT EXISTS ] [ SCHEMA schema ] }
//	    [ SHARING = { METADATA | NONE } ]
//	    [ invoker_rights_clause ]
//	    [ RESOLVER ( ( match_string [,] { schema_name | - } )... ) ]
//	    { USING { BFILE ( directory_object_name, server_file_name )
//	            | { CLOB | BLOB | BFILE } subquery
//	            | 'key_for_BLOB' }
//	    | AS source_char }
type embeddedSourceHead struct {
	state int
	depth int // parenthesis depth inside the JAVA head
}

const (
	javaHeadStart  = iota // at a statement start
	javaHeadCreate        // after CREATE and its modifiers
	javaHeadMLE           // after CREATE ... MLE, before MODULE
	javaHeadJava          // after JAVA or MLE MODULE, before USING or AS
	javaHeadNone          // not CREATE JAVA ... AS or CREATE MLE MODULE ... AS
)

// observe feeds the next token of the statement and reports whether it is
// the AS after which the Java source starts.
func (h *embeddedSourceHead) observe(tok Token) bool {
	switch h.state {
	case javaHeadStart:
		h.state = javaHeadNone
		if tok.Type == kwCREATE {
			h.state = javaHeadCreate
		}
	case javaHeadCreate:
		switch {
		case tok.Type == kwJAVA:
			h.state = javaHeadJava
		case tok.Type == tokIDENT && tok.Str == "MLE":
			h.state = javaHeadMLE
		case tok.Type == kwOR, tok.Type == kwREPLACE, tok.Type == kwAND,
			tok.Type == kwIF, tok.Type == kwNOT, tok.Type == kwEXISTS:
		case tok.Type == tokIDENT && (tok.Str == "RESOLVE" || tok.Str == "COMPILE" || tok.Str == "NOFORCE"):
		default:
			h.state = javaHeadNone
		}
	case javaHeadMLE:
		h.state = javaHeadNone
		if tok.Type == tokIDENT && tok.Str == "MODULE" {
			h.state = javaHeadJava
		}
	case javaHeadJava:
		switch tok.Type {
		case '(':
			h.depth++
		case ')':
			h.depth--
		case kwUSING:
			if h.depth == 0 {
				h.state = javaHeadNone
			}
		case kwAS:
			if h.depth == 0 {
				h.state = javaHeadNone
				return true
			}
		}
	}
	return false
}

// javaSourceEnd returns where Java source starting after offset pos ends:
// before the line break that precedes the next line holding only "/", or at
// the end of sql. next is the offset just past that "/", or len(sql) without
// one. Only the delimiter's framing is dropped: trailing blanks and blank
// lines before it are source_char, kept verbatim.
func javaSourceEnd(sql string, pos int) (end, next int) {
	for line := lineEndAfterBreak(sql, pos); line < len(sql); line = lineEndAfterBreak(sql, line) {
		i := skipHorizontalSpace(sql, line)
		if i < len(sql) && sql[i] == '/' && isSlashDelimiterLine(sql, i, i+1) {
			return trimLineBreakBefore(sql, line, pos), i + 1
		}
	}
	return len(sql), len(sql)
}

// trimLineBreakBefore returns line, a line's start offset, moved back over
// the one line break ("\n", "\r\n", or "\r") that ends the previous line,
// without going below floor.
func trimLineBreakBefore(sql string, line, floor int) int {
	end := line
	if end > floor && sql[end-1] == '\n' {
		end--
	}
	if end > floor && sql[end-1] == '\r' {
		end--
	}
	return end
}

// mleInlineCodeEnd finds the inline code of an MLE call spec, which follows
// MLE LANGUAGE language_name: an optional PURE, then the code wrapped in
// delimiters. pos is the offset just past the language name and limit the end
// of the text to search. It returns where the opening delimiter starts, the
// offset just past the closing delimiter, and whether PURE was given; ok is
// false when there is no delimiter or no closing delimiter.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/mlejs/creating-inline-mle-call-specification.html
//
// Oracle 23ai reads the opening delimiter as the run of non-blank characters
// that follows (so "{{return" is one delimiter) and finds the closing one by a
// plain search of the text after it, string literals included. A delimiter of
// opening brackets closes with its mirror reversed ("{<" with ">}"); any other
// delimiter, "))", "ab", or "'q" alike, closes with itself. A lone "." is not a
// delimiter (PLS-00881).
func mleInlineCodeEnd(sql string, pos, limit int) (codeStart, end int, pure, ok bool) {
	if limit > len(sql) {
		limit = len(sql)
	}
	i := skipSpace(sql, pos, limit)
	if word := blankRun(sql, i, limit); strings.EqualFold(word, "PURE") && i+len(word) < limit {
		pure = true
		i = skipSpace(sql, i+len(word), limit)
	}
	open := blankRun(sql, i, limit)
	if open == "" || open == "." {
		return i, limit, pure, false
	}
	closing := mleClosingDelimiter(open)
	idx := strings.Index(sql[i+len(open):limit], closing)
	if idx < 0 {
		return i, limit, pure, false
	}
	return i, i + len(open) + idx + len(closing), pure, true
}

// mleClosingDelimiter returns the delimiter that closes an inline MLE code
// opened by open.
func mleClosingDelimiter(open string) string {
	mirror := map[byte]byte{'{': '}', '(': ')', '[': ']', '<': '>'}
	closing := make([]byte, len(open))
	for k := 0; k < len(open); k++ {
		m, isBracket := mirror[open[k]]
		if !isBracket {
			return open
		}
		closing[len(open)-1-k] = m
	}
	return string(closing)
}

// blankRun returns the run of non-blank bytes that starts at pos.
func blankRun(sql string, pos, limit int) string {
	end := pos
	for end < limit && !isBlankByte(sql[end]) {
		end++
	}
	return sql[pos:end]
}

func isBlankByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

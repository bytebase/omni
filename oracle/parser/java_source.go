package parser

// javaSourceHead recognizes the head of CREATE ... JAVA ... AS source_char.
// The text after that AS is Java source, not SQL: it holds ';', quotes, and
// comments the SQL lexer would misread. The splitter and the parser take it
// verbatim instead of lexing it, as SQL*Plus does: SQL*Plus buffers a CREATE
// JAVA statement until a line holding only "/".
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/sqlrf/CREATE-JAVA.html
//
//	CREATE [ OR REPLACE ] [ AND { RESOLVE | COMPILE } ] [ NOFORCE ]
//	    JAVA { { SOURCE | RESOURCE } NAMED [ schema. ] primary_name
//	         | CLASS [ SCHEMA schema ] }
//	    [ SHARING = { METADATA | NONE } ]
//	    [ invoker_rights_clause ]
//	    [ RESOLVER ( ( match_string [,] { schema_name | - } )... ) ]
//	    { USING { BFILE ( directory_object_name, server_file_name )
//	            | { CLOB | BLOB | BFILE } subquery
//	            | 'key_for_BLOB' }
//	    | AS source_char }
type javaSourceHead struct {
	state int
	depth int // parenthesis depth inside the JAVA head
}

const (
	javaHeadStart  = iota // at a statement start
	javaHeadCreate        // after CREATE and its modifiers
	javaHeadJava          // after JAVA, before USING or AS
	javaHeadNone          // not CREATE JAVA ... AS
)

// observe feeds the next token of the statement and reports whether it is
// the AS after which the Java source starts.
func (h *javaSourceHead) observe(tok Token) bool {
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
		case tok.Type == kwOR, tok.Type == kwREPLACE, tok.Type == kwAND,
			tok.Type == kwIF, tok.Type == kwNOT, tok.Type == kwEXISTS:
		case tok.Type == tokIDENT && (tok.Str == "RESOLVE" || tok.Str == "COMPILE" || tok.Str == "NOFORCE"):
		default:
			h.state = javaHeadNone
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

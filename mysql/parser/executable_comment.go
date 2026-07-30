package parser

// HasExecutableComment reports whether the lexer has encountered a MySQL
// executable comment (/*! ... */). Callers using the streaming lexer should
// consume the input through EOF before consulting this method.
//
// Unknown version prefixes and unterminated executable comments are reported
// as present so security-sensitive callers can fail closed.
func (l *Lexer) HasExecutableComment() bool {
	return l.hasExecutableComment
}

// ContainsExecutableComment reports whether sql contains a MySQL executable
// comment (/*! ... */). Recognition is performed by the MySQL lexer, so text
// inside string literals, quoted identifiers, ordinary comments, optimizer
// hints, and line comments is not reported.
//
// Unknown version prefixes and unterminated executable comments return true.
func ContainsExecutableComment(sql string) bool {
	lexer := NewLexer(sql)
	for lexer.NextToken().Type != tokEOF {
	}
	return lexer.HasExecutableComment()
}

package parser

// ExecutableCommentOptions configures SQL modes that affect lexical comment
// recognition.
type ExecutableCommentOptions struct {
	// NoBackslashEscapes must match the session's NO_BACKSLASH_ESCAPES SQL
	// mode. When true, backslashes in string literals are ordinary bytes.
	NoBackslashEscapes bool
}

// ContainsExecutableComment reports whether sql contains a MySQL executable
// comment (/*! ... */). Recognition is performed by the MySQL lexer, so text
// inside string literals, quoted identifiers, ordinary comments, optimizer
// hints, and line comments is not reported.
//
// Unknown version prefixes and unterminated executable comments return true.
func ContainsExecutableComment(sql string, options ExecutableCommentOptions) bool {
	lexer := NewLexer(sql)
	lexer.noBackslashEscapes = options.NoBackslashEscapes
	lexer.stopAtExecutableComment = true

	for {
		tok := lexer.NextToken()
		if lexer.hasExecutableComment {
			return true
		}
		if tok.Type == tokEOF {
			return false
		}
	}
}

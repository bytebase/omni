package completion

import (
	"strings"

	"github.com/bytebase/omni/pg/yacc"
)

// tokenInfo holds a token and its parser-internal token ID.
type tokenInfo struct {
	tok         yacc.Token
	parserToken int // the goyacc-internal token ID (after pgTok1/pgTok2/pgTok3 mapping)
}

// preprocessResult holds the result of preprocessing: the isolated statement
// tokens and the index within those tokens where the cursor is.
type preprocessResult struct {
	tokens   []tokenInfo
	cursorAt int // index in tokens where cursor is (tokens before cursor)
}

// preprocess takes SQL text and a cursor byte offset, isolates the statement
// containing the cursor, tokenizes it, and returns the tokens before the cursor.
func preprocess(sql string, cursorOffset int) preprocessResult {
	// Handle pipe-delimited cursor marker (for testing)
	if idx := findCursorMarker(sql); idx >= 0 {
		cursorOffset = idx
		sql = sql[:idx] + sql[idx+1:] // remove the | marker
	}

	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}

	// Tokenize the full input
	allTokens := tokenize(sql)

	// Split by semicolons, find the statement containing the cursor
	stmtTokens, cursorIdx := isolateStatement(allTokens, cursorOffset)

	return preprocessResult{
		tokens:   stmtTokens,
		cursorAt: cursorIdx,
	}
}

// findCursorMarker finds the position of '|' used as cursor marker in test input.
// Returns -1 if no marker found or if '|' is inside a string literal.
func findCursorMarker(sql string) int {
	inSingle := false
	inDouble := false
	for i := 0; i < len(sql); i++ {
		switch sql[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '|':
			if !inSingle && !inDouble {
				return i
			}
		}
	}
	return -1
}

// tokenize converts SQL text into a slice of tokenInfo.
func tokenize(sql string) []tokenInfo {
	lexer := yacc.NewLexer(sql)
	var tokens []tokenInfo

	for {
		tok := lexer.NextToken()
		if tok.Type == 0 { // EOF
			break
		}
		parserTok := mapLexerToken(tok)
		tokens = append(tokens, tokenInfo{
			tok:         tok,
			parserToken: parserTok,
		})
	}
	return tokens
}

// mapLexerToken maps a lexer token to its goyacc-internal token ID.
// Two-step process:
// 1. Convert lex_* internal types to parser constants (like parserLexer.mapTokenType)
// 2. Run through pglex1's tok1/tok2/tok3 tables
func mapLexerToken(tok yacc.Token) int {
	// Step 1: convert lex_* types to parser token constants
	char := mapLexerToParserType(tok)

	// Step 2: convert parser constant to internal goyacc token ID (pglex1 logic)
	if char <= 0 {
		tok1 := yacc.Tok1()
		return int(tok1[0])
	}
	tok1 := yacc.Tok1()
	if char < len(tok1) {
		return int(tok1[char])
	}
	priv := yacc.Private()
	tok2 := yacc.Tok2()
	if char >= priv {
		if char < priv+len(tok2) {
			return int(tok2[char-priv])
		}
	}
	tok3 := yacc.Tok3()
	for i := 0; i < len(tok3); i += 2 {
		if int(tok3[i]) == char {
			return int(tok3[i+1])
		}
	}
	return int(tok2[1]) // unknown char
}

// mapLexerToParserType converts raw lexer token types to parser token constants.
// Replicates parserLexer.mapTokenType from parse.go.
func mapLexerToParserType(tok yacc.Token) int {
	if tok.Type == 0 {
		return 0
	}
	// Single-character tokens map directly
	if tok.Type > 0 && tok.Type < 256 {
		return tok.Type
	}
	// Non-keyword tokens (lex_* constants starting at 800)
	if tok.Type >= 800 && tok.Type < 900 {
		offset := tok.Type - 800
		switch offset {
		case 0: // lex_ICONST
			return yacc.ICONST
		case 1: // lex_FCONST
			return yacc.FCONST
		case 2: // lex_SCONST
			return yacc.SCONST
		case 3: // lex_BCONST
			return yacc.BCONST
		case 4: // lex_XCONST
			return yacc.XCONST
		case 5: // lex_USCONST
			return yacc.SCONST
		case 6: // lex_IDENT
			return yacc.IDENT
		case 7: // lex_UIDENT
			return yacc.IDENT
		case 8: // lex_TYPECAST
			return yacc.TYPECAST
		case 9: // lex_DOT_DOT
			return yacc.DOT_DOT
		case 10: // lex_COLON_EQUALS
			return yacc.COLON_EQUALS
		case 11: // lex_EQUALS_GREATER
			return yacc.EQUALS_GREATER
		case 12: // lex_LESS_EQUALS
			return yacc.LESS_EQUALS
		case 13: // lex_GREATER_EQUALS
			return yacc.GREATER_EQUALS
		case 14: // lex_NOT_EQUALS
			return yacc.NOT_EQUALS
		case 15: // lex_PARAM
			return yacc.PARAM
		case 16: // lex_Op
			return yacc.Op
		}
		return 0
	}
	// Keywords and other parser constants pass through directly
	return tok.Type
}

// isolateStatement finds the statement containing the cursor offset and returns
// only the tokens before the cursor.
func isolateStatement(tokens []tokenInfo, cursorOffset int) ([]tokenInfo, int) {
	if len(tokens) == 0 {
		return nil, 0
	}

	// Find semicolons to split statements
	// The cursor is at a byte offset; tokens before the cursor belong to the current statement
	stmtStart := 0
	cursorIdx := len(tokens) // default: cursor is after all tokens

	for i, t := range tokens {
		// Is this token at or past the cursor?
		if t.tok.Loc >= cursorOffset {
			cursorIdx = i
			break
		}
		// Semicollon: start a new statement after this
		if isSemicolon(t) {
			stmtStart = i + 1
		}
	}

	// Also check if there are semicolons after cursorIdx
	stmtEnd := cursorIdx
	for i := cursorIdx; i < len(tokens); i++ {
		if isSemicolon(tokens[i]) {
			break
		}
		stmtEnd = i + 1
	}
	_ = stmtEnd // we don't need tokens after cursor for completion

	stmtTokens := tokens[stmtStart:cursorIdx]
	return stmtTokens, len(stmtTokens)
}

// isSemicolon checks if a token is a semicolon.
func isSemicolon(t tokenInfo) bool {
	return t.tok.Type == ';'
}

// extractPrefix extracts the text prefix at the cursor position for filtering.
// If the cursor is right after "tab|", prefix is "tab".
func extractPrefix(sql string, cursorOffset int) string {
	if cursorOffset <= 0 || cursorOffset > len(sql) {
		return ""
	}
	// Walk backwards to find the start of the current identifier
	end := cursorOffset
	start := end
	for start > 0 {
		r := rune(sql[start-1])
		if isIdentChar(r) {
			start--
		} else {
			break
		}
	}
	return strings.ToLower(sql[start:end])
}

func isIdentChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') || r == '_' || r == '$'
}

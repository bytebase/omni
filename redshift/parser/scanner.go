package parser

// Scanner provides a cursor-based token stream for navigating over lexed tokens.
// It mirrors the API surface of Bytebase's base.Scanner, allowing Bytebase's
// scanner-based analysis functions (determineQualifiedName, determineColumnRef,
// collectTableReferences) to work with omni's token stream.
type Scanner struct {
	tokens     []Token
	index      int
	tokenStack []int
}

// NewScanner creates a Scanner from a SQL string by fully lexing it.
// Token types are mapped from lexer-internal constants to parser constants.
func NewScanner(sql string) *Scanner {
	l := NewLexer(sql)
	var tokens []Token
	for {
		tok := l.NextToken()
		tok.Type = mapTokenType(tok.Type)
		tokens = append(tokens, tok)
		if tok.Type == 0 { // EOF
			break
		}
	}
	return &Scanner{
		tokens: tokens,
		index:  0,
	}
}

// NewScannerFromTokens creates a Scanner from pre-lexed tokens.
func NewScannerFromTokens(tokens []Token) *Scanner {
	return &Scanner{
		tokens: tokens,
		index:  0,
	}
}

// GetIndex returns the current token index.
func (s *Scanner) GetIndex() int {
	return s.index
}

// GetTokenType returns the type of the current token.
func (s *Scanner) GetTokenType() int {
	if s.index >= len(s.tokens) {
		return 0 // EOF
	}
	return s.tokens[s.index].Type
}

// GetTokenText returns the string value of the current token.
func (s *Scanner) GetTokenText() string {
	if s.index >= len(s.tokens) {
		return ""
	}
	return s.tokens[s.index].Str
}

// GetTokenOffset returns the byte offset of the current token.
func (s *Scanner) GetTokenOffset() int {
	if s.index >= len(s.tokens) {
		return -1
	}
	return s.tokens[s.index].Loc
}

// GetToken returns the current token.
func (s *Scanner) GetToken() Token {
	if s.index >= len(s.tokens) {
		return Token{}
	}
	return s.tokens[s.index]
}

// Forward moves the index forward by one token, skipping whitespace/comments.
// Returns false if already at EOF.
func (s *Scanner) Forward() bool {
	if s.index >= len(s.tokens)-1 {
		return false
	}
	s.index++
	return true
}

// Backward moves the index backward by one token.
// Returns false if already at the beginning.
func (s *Scanner) Backward() bool {
	if s.index <= 0 {
		return false
	}
	s.index--
	return true
}

// Push saves the current index on the stack.
func (s *Scanner) Push() {
	s.tokenStack = append(s.tokenStack, s.index)
}

// PopAndRestore restores the index from the stack.
// Returns false if the stack is empty.
func (s *Scanner) PopAndRestore() bool {
	if len(s.tokenStack) == 0 {
		return false
	}
	s.index = s.tokenStack[len(s.tokenStack)-1]
	s.tokenStack = s.tokenStack[:len(s.tokenStack)-1]
	return true
}

// PopAndDiscard discards the saved index from the stack without restoring.
func (s *Scanner) PopAndDiscard() bool {
	if len(s.tokenStack) == 0 {
		return false
	}
	s.tokenStack = s.tokenStack[:len(s.tokenStack)-1]
	return true
}

// GetPreviousTokenType returns the type of the previous token without
// changing the current index.
func (s *Scanner) GetPreviousTokenType() int {
	if s.index <= 0 {
		return 0
	}
	return s.tokens[s.index-1].Type
}

// GetPreviousTokenText returns the text of the previous token without
// changing the current index.
func (s *Scanner) GetPreviousTokenText() string {
	if s.index <= 0 {
		return ""
	}
	return s.tokens[s.index-1].Str
}

// IsTokenType checks if the current token matches the given type.
func (s *Scanner) IsTokenType(tokenType int) bool {
	return s.GetTokenType() == tokenType
}

// IsIdentifier returns true if the given token type is an identifier
// (IDENT or a non-reserved keyword that can be used as an identifier).
func (s *Scanner) IsIdentifier(tokenType int) bool {
	if tokenType == IDENT {
		return true
	}
	for i := range Keywords {
		if Keywords[i].Token == tokenType {
			cat := Keywords[i].Category
			return cat == UnreservedKeyword || cat == ColNameKeyword || cat == TypeFuncNameKeyword
		}
	}
	return false
}

// IsCurrentIdentifier returns true if the current token is an identifier.
func (s *Scanner) IsCurrentIdentifier() bool {
	return s.IsIdentifier(s.GetTokenType())
}

// SeekIndex moves the scanner to the given token index.
func (s *Scanner) SeekIndex(index int) {
	if index < 0 {
		index = 0
	}
	if index >= len(s.tokens) {
		index = len(s.tokens) - 1
	}
	s.index = index
}

// SeekOffset moves the scanner to the token at or after the given byte offset.
func (s *Scanner) SeekOffset(offset int) {
	for i, tok := range s.tokens {
		if tok.Loc >= offset {
			s.index = i
			return
		}
	}
	s.index = len(s.tokens) - 1
}

// GetFollowingText returns the source text from the current token position
// to the end of the input. Requires the original SQL string.
func (s *Scanner) GetFollowingText(sql string) string {
	off := s.GetTokenOffset()
	if off < 0 || off >= len(sql) {
		return ""
	}
	return sql[off:]
}

// GetFollowingTextAfter returns the source text from the given byte offset
// to the end of the input.
func (s *Scanner) GetFollowingTextAfter(sql string, offset int) string {
	if offset < 0 || offset >= len(sql) {
		return ""
	}
	return sql[offset:]
}

// Size returns the number of tokens in the scanner.
func (s *Scanner) Size() int {
	return len(s.tokens)
}

// Tokens returns all tokens in the scanner.
func (s *Scanner) Tokens() []Token {
	return s.tokens
}

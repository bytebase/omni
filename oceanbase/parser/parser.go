package parser

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// Parser is a recursive descent parser for MySQL SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid

	castShiftHalf bool // the inner of ARRAY<ARRAY<...>> has closed half of a >> (see parseCastArrayElement)

	// reservedColumnRef is the last column reference read in a form column_ref has and
	// column_definition_ref lacks: with a leading dot, or with a reserved word for a
	// name. The JSON arrows refuse it as their left operand.
	reservedColumnRef *nodes.ColumnRef

	// parenRunStart and parenRunPast cache the last firstTokenPastParens scan:
	// the offset of the run's first '(' and the first token past the run.
	parenRunStart int
	parenRunPast  Token

	// parenHasComma holds parenGroupHasComma's verdict for each '(' of the
	// groups it scanned: one scan answers for every group it passes, so a
	// nested INTERVAL (INTERVAL (...)) costs one pass and not one per level.
	parenHasComma map[int]bool

	// parenQuery holds parenOpensQuery's verdict for each '(' of the runs it
	// classified, and parenAfter, for each '(' of the blocks it scanned, the
	// type of the token after the matching ')'.
	parenQuery map[int]bool
	parenAfter map[int]int

	// parenQueryFail holds every parseSelectWithParens failure by the offset
	// of its '('. parseSubqueryHead speculates a query at each level of a
	// nested `((((SELECT 1) = 0) = 0) = 0)` and rewinds when it fails, and the
	// next level down fails on the same tokens; recording the failures keeps
	// that linear.
	parenQueryFail map[int]parenQueryFailRecord

	// parenQueryReplay holds the same failures in completion mode, where a
	// failed parse also leaves candidates behind (parenQueryFailure).
	parenQueryReplay map[parenQueryKey]*parenQueryFailure

	depth int // nesting depth of the productions being parsed (see enterNesting)

	noBareDefault bool // inside an expression that takes DEFAULT(col) but no bare DEFAULT

	newlines []int // offsets of the input's newlines, found on first use (see lineCol)

	// Completion mode fields.
	completing bool          // true when collecting completion candidates
	cursorOff  int           // byte offset of the cursor in source
	candidates *CandidateSet // collected candidates
	collecting bool          // true once cursor position is reached
	maxCollect int           // max exploration depth
}

// Parse parses a SQL string into an AST list.
// It uses Split() internally to handle DELIMITER directives and compound blocks,
// then parses each segment individually.
func Parse(sql string) (*nodes.List, error) {
	segments := Split(sql)
	if len(segments) == 0 {
		return &nodes.List{}, nil
	}
	var allStmts []nodes.Node
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		list, err := parseSingle(seg.Text, seg.ByteStart)
		if err != nil {
			// Re-enrich error with full SQL for correct absolute line/col.
			if pe, ok := err.(*ParseError); ok {
				pe.Line, pe.Column = lineColStatic(sql, pe.Position)
				pe.RelatedText = extractLineStatic(sql, pe.Position)
			}
			return nil, err
		}
		allStmts = append(allStmts, list.Items...)
	}
	if len(allStmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: allStmts}, nil
}

// parseSingle parses a single SQL segment (without DELIMITER directives) into an AST list.
// baseOffset is added to all token positions so Loc values are absolute within the original input.
//
// A statement that does not parse with its executable comments' bodies in it
// is read again with every executable comment as a plain comment, as OceanBase
// reads it (exec_fallback.go): "SELECT 1 /*! SELECT 2 */" is SELECT 1.
func parseSingle(sql string, baseOffset int) (*nodes.List, error) {
	p := &Parser{
		lexer: NewLexerWithOffset(sql, baseOffset),
	}
	p.advance()

	var stmts []nodes.Node
	// from is where the statement being read starts in sql: past the last ';',
	// so a retry reads an executable comment that opens the statement.
	// opener is the offset of the first "/*!" at or after from, once looked up,
	// or len(sql) when there is none.
	from, opener := 0, -1
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			from = p.cur.End - baseOffset
			p.advance()
			continue
		}
		stmt, err := p.readStatement()
		// The nesting limit is this parser's, not a reading the engine fails:
		// read again without the comments, the statement would lose the
		// expression that went too deep, and "SELECT 0 /*! OR 1 OR 1 ... */"
		// would pass as SELECT 0.
		if isNestingError(err) {
			return nil, err
		}
		// A statement that parsed and left an executable comment open has no
		// retry; lexerError reports the comment after the loop.
		if err != nil || p.lexer.hasErr {
			if opener < from {
				opener = len(sql)
				if i := strings.Index(sql[from:], "/*!"); i >= 0 {
					opener = from + i
				}
			}
			if opener == len(sql) || !retriesAsSQL(sql[from:]) {
				if err == nil {
					err = p.enrichError(p.lexerError())
				}
				return nil, err
			}
			plain := p.restart(from, true)
			first := plain.cur
			if stmt, err = plain.readStatement(); err != nil {
				return nil, err
			}
			if lexErr := plain.lexerError(); lexErr != nil {
				return nil, plain.enrichError(lexErr)
			}
			if stmt != nil && !sqlParserReads(stmt) {
				return nil, plain.syntaxErrorAtTok(first)
			}
			// What follows the statement is read with its executable comments
			// in it again.
			p = plain.restart(plain.cur.Loc-baseOffset, false)
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	// A malformed token that ran to EOF without its closing delimiter is recorded
	// on the lexer rather than panicking; surface it as a clean parse error even
	// when the statements otherwise parsed successfully.
	if lexErr := p.lexerError(); lexErr != nil {
		return nil, p.enrichError(lexErr)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// restart returns a parser over p's input that starts at pos, reading every
// executable comment as a plain comment when plainExec is set. It keeps what
// p knows of the input that does not depend on where reading starts: the
// lexer's close-scan memo (restartAt) and the offsets of the newlines, which a
// syntax error looks up. Rebuilding those for each statement read again would
// make a segment of n such statements quadratic.
func (p *Parser) restart(pos int, plainExec bool) *Parser {
	q := &Parser{lexer: p.lexer.restartAt(pos, plainExec), newlines: p.newlines}
	q.advance()
	return q
}

// readStatement reads the statement at the current token, which ends at a ';'
// or at the end of the text. It returns no statement and no error when the
// text holds none, as when every token was in a comment.
func (p *Parser) readStatement() (nodes.Node, error) {
	if p.cur.Type == ';' || p.cur.Type == tokEOF {
		return nil, nil
	}
	stmtStart := p.cur.Loc
	stmt, err := p.parseStmt()
	if err != nil {
		// A lexing error (unterminated comment/string/identifier) is the more
		// specific cause; report it in preference to the derived syntax error.
		if lexErr := p.lexerError(); lexErr != nil {
			return nil, p.enrichError(lexErr)
		}
		return nil, p.enrichError(err)
	}
	if stmt == nil {
		// parseStmt returned nil without error — unconsumed tokens remain.
		if lexErr := p.lexerError(); lexErr != nil {
			return nil, p.enrichError(lexErr)
		}
		if p.cur.Type != tokEOF {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}
	if treeTooHigh(stmt, p.prev.End-stmtStart) {
		return nil, p.enrichError(nestedTooDeeply(stmtStart))
	}
	// A statement ends at a ';' or at the end of the text. Any other token
	// here starts a second statement with no separator in front of it,
	// which the server refuses with 1064.
	//
	// Ref: OceanBase CE v4.3.5_CE_BP1 src/sql/parser/sql_parser_mysql_mode.y:591
	// (stmt_list), which follows a stmt with DELIMITER or END_P and nothing
	// else.
	if p.cur.Type != ';' && p.cur.Type != tokEOF {
		if lexErr := p.lexerError(); lexErr != nil {
			return nil, p.enrichError(lexErr)
		}
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// lexerError converts a recorded lexer error, if any, into a ParseError with an
// absolute position. It returns nil when the lexer recorded no error.
//
// An executable comment that never closed is reported last, because a token
// inside its body that ran off the end is the more specific cause.
func (p *Parser) lexerError() *ParseError {
	if p.lexer.hasErr {
		return &ParseError{
			Message:  p.lexer.errMsg,
			Position: p.lexer.errPos + p.lexer.baseOffset,
		}
	}
	if p.lexer.execUnclosed {
		return &ParseError{
			Message:  "unterminated comment",
			Position: p.lexer.execUnclosedPos + p.lexer.baseOffset,
		}
	}
	return nil
}

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// match checks if the current token type matches any of the given types.
// If it matches, the token is consumed and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	if p.cur.Type == tokEOF {
		return Token{}, p.syntaxErrorAtCur()
	}
	return Token{}, &ParseError{
		Message:  "unexpected token",
		Position: p.cur.Loc,
	}
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// inputText returns the input between start and end byte positions, without the
// executable-comment syntax (/*!NNNNN, */) the lexer stepped over in that range.
// start and end are absolute (include baseOffset), so we subtract it to index into the segment text.
func (p *Parser) inputText(start, end int) string {
	start -= p.lexer.baseOffset
	end -= p.lexer.baseOffset
	if start < 0 {
		start = 0
	}
	if end > len(p.lexer.input) {
		end = len(p.lexer.input)
	}
	if start >= end {
		return ""
	}
	return p.lexer.expandedText(start, end)
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message     string
	Position    int
	Line        int    // 1-based line number
	Column      int    // 1-based column number
	RelatedText string // source text of the line where the error occurred
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		if e.RelatedText != "" {
			return fmt.Sprintf("%s (line %d, column %d)\nrelated text: %s", e.Message, e.Line, e.Column, e.RelatedText)
		}
		return fmt.Sprintf("%s (line %d, column %d)", e.Message, e.Line, e.Column)
	}
	return e.Message
}

// syntaxErrorAtCur returns a ParseError describing a syntax error at the current token.
// If the current token is EOF, the message says "at end of input".
// Otherwise it says "at or near" with the token text.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	return p.syntaxErrorAtTok(p.cur)
}

// syntaxErrorAtTok returns a ParseError describing a syntax error at the given token.
func (p *Parser) syntaxErrorAtTok(tok Token) *ParseError {
	line, col := p.lineCol(tok.Loc)
	var msg string
	if tok.Type == tokEOF {
		msg = "syntax error at end of input"
	} else {
		text := tok.Str
		if text == "" {
			text = string(rune(tok.Type))
		}
		msg = fmt.Sprintf("syntax error at or near %q", text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
		Line:     line,
		Column:   col,
	}
}

// lineCol computes the 1-based line and column for a byte offset in the input.
// The offset is absolute (includes baseOffset), so we subtract it to index into segment text.
//
// Every syntax error asks, and completion resumes after each statement that
// fails, so the newlines are found once and searched: counting them from the
// start of the input for each error made completion over a run of n failing
// tokens quadratic.
func (p *Parser) lineCol(offset int) (int, int) {
	offset -= p.lexer.baseOffset
	input := p.lexer.input
	if offset < 0 {
		offset = 0
	}
	if offset > len(input) {
		offset = len(input)
	}
	newlines := p.newlineOffsets()
	// The newlines before offset.
	n := sort.SearchInts(newlines, offset)
	lastNewline := -1
	if n > 0 {
		lastNewline = newlines[n-1]
	}
	return 1 + n, offset - lastNewline
}

// newlineOffsets returns the offsets of the input's newlines, found on first
// use and kept for every later error (see lineCol and extractLine).
func (p *Parser) newlineOffsets() []int {
	if p.newlines == nil {
		input := p.lexer.input
		p.newlines = make([]int, 0, strings.Count(input, "\n"))
		for i := 0; i < len(input); i++ {
			if input[i] == '\n' {
				p.newlines = append(p.newlines, i)
			}
		}
	}
	return p.newlines
}

// lineColStatic computes line and column from raw input (for use outside a Parser).
func lineColStatic(input string, offset int) (int, int) {
	if offset > len(input) {
		offset = len(input)
	}
	line := 1
	lastNewline := -1
	for i := 0; i < offset; i++ {
		if input[i] == '\n' {
			line++
			lastNewline = i
		}
	}
	col := offset - lastNewline
	return line, col
}

// enrichError fills in Line/Column and RelatedText on a ParseError using the parser's input.
func (p *Parser) enrichError(err error) error {
	if pe, ok := err.(*ParseError); ok {
		if pe.Line == 0 {
			pe.Line, pe.Column = p.lineCol(pe.Position)
		}
		if pe.RelatedText == "" {
			pe.RelatedText = p.extractLine(pe.Position)
		}
	}
	return err
}

// extractLine returns the source line containing the given byte offset.
// The offset is absolute (includes baseOffset), so we subtract it to index into segment text.
//
// The line's ends come from the newline offsets, not from a scan out from
// offset: every statement that fails its first reading and is read again
// (parseSingle) builds an error, and scanning a line that holds every
// statement, as one under a client delimiter can, made that quadratic.
func (p *Parser) extractLine(offset int) string {
	offset -= p.lexer.baseOffset
	input := p.lexer.input
	if offset < 0 || offset > len(input) {
		return ""
	}
	newlines := p.newlineOffsets()
	n := sort.SearchInts(newlines, offset)
	start, end := 0, len(input)
	if n > 0 {
		start = newlines[n-1] + 1
	}
	if n < len(newlines) {
		end = newlines[n]
	}
	line := input[start:end]
	// Trim to reasonable length
	if len(line) > 200 {
		line = line[:200] + "..."
	}
	return line
}

// extractLineStatic returns the source line containing the given byte offset (for use outside a Parser).
func extractLineStatic(input string, offset int) string {
	if offset < 0 || offset > len(input) {
		return ""
	}
	start := offset
	for start > 0 && input[start-1] != '\n' {
		start--
	}
	end := offset
	for end < len(input) && input[end] != '\n' {
		end++
	}
	line := input[start:end]
	if len(line) > 200 {
		line = line[:200] + "..."
	}
	return line
}

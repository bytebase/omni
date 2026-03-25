// Package parser implements a recursive descent PL/pgSQL parser.
package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/pg/plpgsql/ast"

	pgparser "github.com/bytebase/omni/pg/parser"
)

// Parser is a recursive descent parser for PL/pgSQL function bodies.
type Parser struct {
	lexer   *pgparser.Lexer
	source  string
	cur     pgparser.Token
	prev    pgparser.Token
	nextBuf pgparser.Token
	hasNext bool
}

// newParser creates a new PL/pgSQL parser for the given source text.
func newParser(source string) *Parser {
	p := &Parser{
		lexer:  pgparser.NewLexer(source),
		source: source,
	}
	p.advance()
	return p
}

// advance consumes the current token and moves to the next one.
// Token types are mapped from lexer-internal types to parser constants.
func (p *Parser) advance() pgparser.Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		tok := p.lexer.NextToken()
		tok.Type = mapLexTokenType(tok.Type)
		p.cur = tok
	}
	return p.prev
}

// mapLexTokenType maps lexer-internal token types (lex_* constants at 800+)
// to the parser token constants used in pg/parser/tokens.go.
// Single-char tokens (0-255) and keyword tokens pass through unchanged.
func mapLexTokenType(typ int) int {
	if typ == 0 || (typ > 0 && typ < 256) {
		return typ
	}
	// Non-keyword token base is 800 in the lexer
	const nonKwBase = 800
	if typ >= nonKwBase && typ < nonKwBase+100 {
		offset := typ - nonKwBase
		switch offset {
		case 0: // lex_ICONST
			return pgparser.ICONST
		case 1: // lex_FCONST
			return pgparser.FCONST
		case 2: // lex_SCONST
			return pgparser.SCONST
		case 3: // lex_BCONST
			return pgparser.BCONST
		case 4: // lex_XCONST
			return pgparser.XCONST
		case 5: // lex_USCONST
			return pgparser.SCONST
		case 6: // lex_IDENT
			return pgparser.IDENT
		case 7: // lex_UIDENT
			return pgparser.IDENT
		case 8: // lex_TYPECAST
			return pgparser.TYPECAST
		case 9: // lex_DOT_DOT
			return pgparser.DOT_DOT
		case 10: // lex_COLON_EQUALS
			return pgparser.COLON_EQUALS
		case 11: // lex_EQUALS_GREATER
			return pgparser.EQUALS_GREATER
		case 12: // lex_LESS_EQUALS
			return pgparser.LESS_EQUALS
		case 13: // lex_GREATER_EQUALS
			return pgparser.GREATER_EQUALS
		case 14: // lex_NOT_EQUALS
			return pgparser.NOT_EQUALS
		case 15: // lex_PARAM
			return pgparser.PARAM
		case 16: // lex_Op
			return pgparser.Op
		}
		return 0
	}
	return typ
}

// peek returns the current token without consuming it.
func (p *Parser) peek() pgparser.Token {
	return p.cur
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() pgparser.Token {
	if !p.hasNext {
		tok := p.lexer.NextToken()
		tok.Type = mapLexTokenType(tok.Type)
		p.nextBuf = tok
		p.hasNext = true
	}
	return p.nextBuf
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (pgparser.Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return pgparser.Token{}, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// tokenText returns a human-readable text representation of a token.
func (p *Parser) tokenText(tok pgparser.Token) string {
	if tok.Type == 0 {
		return "end of input"
	}
	if tok.Str != "" {
		return tok.Str
	}
	if tok.Type > 0 && tok.Type < 256 {
		return string(rune(tok.Type))
	}
	return p.source[tok.Loc:tok.End]
}

// errorf returns a formatted parse error with position information.
func (p *Parser) errorf(format string, args ...any) error {
	return &ParseError{
		Message:  fmt.Sprintf(format, args...),
		Position: p.cur.Loc,
	}
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("at position %d: %s", e.Position, e.Message)
}

// isIdent checks if the current token is an IDENT token.
func (p *Parser) isIdent() bool {
	return p.cur.Type == pgparser.IDENT
}

// isKeyword checks if the current token matches a PL/pgSQL keyword (case-insensitive).
// PL/pgSQL keywords come through as IDENT from the SQL lexer, or as SQL keywords
// that happen to share the same name (e.g., BEGIN, END, NOT, NULL, FOR, etc.).
func (p *Parser) isKeyword(name string) bool {
	lower := strings.ToLower(name)
	// Check IDENT tokens against PL/pgSQL keyword table
	if p.cur.Type == pgparser.IDENT {
		return strings.EqualFold(p.cur.Str, name)
	}
	// Some PL/pgSQL keywords are also SQL keywords and come through as their SQL token type.
	// We need to check the token text against the keyword name.
	switch lower {
	case "begin":
		return p.cur.Type == pgparser.BEGIN_P
	case "end":
		return p.cur.Type == pgparser.END_P
	case "not":
		return p.cur.Type == pgparser.NOT
	case "null":
		return p.cur.Type == pgparser.NULL_P
	case "for":
		return p.cur.Type == pgparser.FOR
	case "case":
		return p.cur.Type == pgparser.CASE
	case "if":
		return p.cur.Type == pgparser.IF_P
	case "else":
		return p.cur.Type == pgparser.ELSE
	case "in":
		return p.cur.Type == pgparser.IN_P
	case "into":
		return p.cur.Type == pgparser.INTO
	case "from":
		return p.cur.Type == pgparser.FROM
	case "to":
		return p.cur.Type == pgparser.TO
	case "or":
		return p.cur.Type == pgparser.OR
	case "when":
		return p.cur.Type == pgparser.WHEN
	case "all":
		return p.cur.Type == pgparser.ALL
	case "by":
		return p.cur.Type == pgparser.BY
	case "execute":
		return p.cur.Type == pgparser.EXECUTE
	case "then":
		return p.cur.Type == pgparser.THEN
	case "using":
		return p.cur.Type == pgparser.USING
	case "declare":
		return p.cur.Type == pgparser.DECLARE
	case "collate":
		return p.cur.Type == pgparser.COLLATE
	case "default":
		return p.cur.Type == pgparser.DEFAULT
	case "cursor":
		return p.cur.Type == pgparser.CURSOR
	case "scroll":
		return p.cur.Type == pgparser.SCROLL
	case "no":
		return p.cur.Type == pgparser.NO
	case "is":
		return p.cur.Type == pgparser.IS
	case "type":
		return p.cur.Type == pgparser.TYPE_P
	case "select":
		return p.cur.Type == pgparser.SELECT
	case "with":
		return p.cur.Type == pgparser.WITH
	case "array":
		return p.cur.Type == pgparser.ARRAY
	case "close":
		return p.cur.Type == pgparser.CLOSE
	case "continue":
		return p.cur.Type == pgparser.CONTINUE_P
	case "current":
		return p.cur.Type == pgparser.CURRENT_P
	case "return":
		return p.cur.Type == pgparser.RETURN
	case "next":
		return p.cur.Type == pgparser.NEXT
	case "strict":
		return p.cur.Type == pgparser.STRICT_P
	case "insert":
		return p.cur.Type == pgparser.INSERT
	case "update":
		return p.cur.Type == pgparser.UPDATE
	case "delete":
		return p.cur.Type == pgparser.DELETE_P
	case "merge":
		return p.cur.Type == pgparser.MERGE
	case "import":
		return p.cur.Type == pgparser.IMPORT_P
	case "fetch":
		return p.cur.Type == pgparser.FETCH
	case "move":
		return p.cur.Type == pgparser.MOVE
	case "absolute":
		return p.cur.Type == pgparser.ABSOLUTE_P
	case "backward":
		return p.cur.Type == pgparser.BACKWARD
	case "forward":
		return p.cur.Type == pgparser.FORWARD
	case "first":
		return p.cur.Type == pgparser.FIRST_P
	case "last":
		return p.cur.Type == pgparser.LAST_P
	case "prior":
		return p.cur.Type == pgparser.PRIOR
	case "relative":
		return p.cur.Type == pgparser.RELATIVE_P
	case "column":
		return p.cur.Type == pgparser.COLUMN
	case "constraint":
		return p.cur.Type == pgparser.CONSTRAINT
	case "table":
		return p.cur.Type == pgparser.TABLE
	case "schema":
		return p.cur.Type == pgparser.SCHEMA
	case "and":
		return p.cur.Type == pgparser.AND
	case "call":
		return p.cur.Type == pgparser.CALL
	case "commit":
		return p.cur.Type == pgparser.COMMIT
	case "rollback":
		return p.cur.Type == pgparser.ROLLBACK
	case "do":
		return p.cur.Type == pgparser.DO
	case "chain":
		return p.cur.Type == pgparser.CHAIN
	}
	return false
}

// consumeKeyword consumes the current token if it matches the given PL/pgSQL keyword.
// Returns true if consumed.
func (p *Parser) consumeKeyword(name string) bool {
	if p.isKeyword(name) {
		p.advance()
		return true
	}
	return false
}

// expectKeyword consumes the current token if it matches the given keyword,
// otherwise returns an error.
func (p *Parser) expectKeyword(name string) error {
	if p.isKeyword(name) {
		p.advance()
		return nil
	}
	return p.errorf("syntax error at or near %q, expected %s", p.tokenText(p.cur), strings.ToUpper(name))
}

// isEOF checks if the current token is end of input.
func (p *Parser) isEOF() bool {
	return p.cur.Type == 0
}

// --------------------------------------------------------------------------
// Entry point
// --------------------------------------------------------------------------

// parse parses the source as a PL/pgSQL function body and returns the root block.
func (p *Parser) parse() (*ast.PLBlock, error) {
	// Parse optional compiler directives (#option, #print_strict_params, #variable_conflict)
	options, err := p.parseOptions()
	if err != nil {
		return nil, err
	}

	// Parse optional label
	label, err := p.parseOptionalLabel()
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock(label)
	if err != nil {
		return nil, err
	}

	block.Options = options

	// After the block, allow optional trailing semicolon then expect EOF
	if p.cur.Type == ';' {
		p.advance()
	}
	if !p.isEOF() {
		return nil, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
	}

	return block, nil
}

// --------------------------------------------------------------------------
// Label parsing
// --------------------------------------------------------------------------

// parseOptionalLabel parses an optional <<label>> syntax.
// Returns the label name (empty if none) and any error.
// The SQL lexer tokenizes << and >> as Op tokens with Str "<<" and ">>".
func (p *Parser) parseOptionalLabel() (string, error) {
	// Check for << operator token
	if !p.isOp("<<") {
		return "", nil
	}
	p.advance() // consume <<

	// Expect identifier for label name
	if !p.isIdent() && !p.isAnyKeywordAsIdent() {
		return "", p.errorf("syntax error at or near %q, expected label name", p.tokenText(p.cur))
	}
	label := p.identText()
	p.advance()

	// Expect >> operator token
	if !p.isOp(">>") {
		return "", p.errorf("syntax error at or near %q, expected >>", p.tokenText(p.cur))
	}
	p.advance()

	return label, nil
}

// isOp checks if the current token is an Op with the given string value.
func (p *Parser) isOp(op string) bool {
	return p.cur.Type == pgparser.Op && p.cur.Str == op
}

// isAnyKeywordAsIdent checks if the current token is a SQL keyword that can be
// used as a PL/pgSQL identifier (label name, variable name, etc.).
// Any SQL keyword (type > 256 with non-empty Str) can serve as a PL/pgSQL identifier.
func (p *Parser) isAnyKeywordAsIdent() bool {
	return p.cur.Type > 256 && p.cur.Str != "" && p.cur.Type != pgparser.Op
}

// identText returns the text of the current token as a lowercase identifier.
func (p *Parser) identText() string {
	if p.cur.Type == pgparser.IDENT {
		return strings.ToLower(p.cur.Str)
	}
	// For SQL keyword tokens used as identifiers, extract from source
	return strings.ToLower(p.source[p.cur.Loc:p.cur.End])
}

// --------------------------------------------------------------------------
// Block parsing
// --------------------------------------------------------------------------

// parseBlock parses [DECLARE ...] BEGIN stmts END [label]
func (p *Parser) parseBlock(label string) (*ast.PLBlock, error) {
	startPos := p.pos()
	if label != "" {
		// The label was already parsed before the block, so start from the label position
		// We don't have the label start pos here, so use current position
	}

	block := &ast.PLBlock{
		Label: label,
	}

	// Parse optional DECLARE section
	if p.isKeyword("DECLARE") {
		p.advance()
		decls, err := p.parseDeclareSection()
		if err != nil {
			return nil, err
		}
		block.Declarations = decls
	}

	// Expect BEGIN
	if err := p.expectKeyword("BEGIN"); err != nil {
		return nil, err
	}

	// Parse statement list
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	block.Body = body

	// Parse optional EXCEPTION section (before END)
	if p.isKeyword("EXCEPTION") {
		exceptions, err := p.parseExceptionBlock()
		if err != nil {
			return nil, err
		}
		block.Exceptions = exceptions
	}

	// Expect END
	if !p.isKeyword("END") {
		return nil, p.errorf("syntax error at or near %q, expected END", p.tokenText(p.cur))
	}
	p.advance()

	// Check for mismatched END keyword (e.g., END IF, END LOOP, END CASE)
	if p.isKeyword("IF") || p.isKeyword("LOOP") || p.isKeyword("CASE") {
		return nil, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
	}

	// Check for optional trailing label after END.
	// The label can be any identifier (IDENT) or SQL keyword used as a PL/pgSQL identifier.
	// But it must not be a PL/pgSQL keyword that starts a new construct (IF, LOOP, CASE).
	if (p.isIdent() || p.isAnyKeywordAsIdent()) && !p.isKeyword("IF") && !p.isKeyword("LOOP") && !p.isKeyword("CASE") {
		endLabel := p.identText()
		if label == "" || !strings.EqualFold(endLabel, label) {
			if label == "" {
				return nil, p.errorf("end label %q specified for unlabeled block", endLabel)
			}
			return nil, p.errorf("end label %q does not match block label %q", endLabel, label)
		}
		p.advance()
	}

	// Consume optional semicolon at end of block (for nested blocks)
	// The caller handles the trailing semicolon for the top-level block.

	block.Loc = ast.Loc{Start: startPos, End: p.prev.End}
	return block, nil
}

// --------------------------------------------------------------------------
// Statement list parsing
// --------------------------------------------------------------------------

// parseStmtList parses statements until END, EXCEPTION, ELSE, ELSIF, WHEN.
func (p *Parser) parseStmtList() ([]ast.Node, error) {
	var stmts []ast.Node
	for {
		// Skip standalone semicolons
		for p.cur.Type == ';' {
			p.advance()
		}
		// Check for terminators
		if p.isEOF() || p.isKeyword("END") || p.isKeyword("EXCEPTION") ||
			p.isKeyword("ELSE") || p.isKeyword("ELSIF") || p.isKeyword("ELSEIF") ||
			p.isKeyword("WHEN") {
			break
		}
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// --------------------------------------------------------------------------
// Statement dispatch
// --------------------------------------------------------------------------

// parseStmt parses a single PL/pgSQL statement.
func (p *Parser) parseStmt() (ast.Node, error) {
	// Check for label before statement (for labeled nested blocks/loops)
	label, err := p.parseOptionalLabel()
	if err != nil {
		return nil, err
	}

	// Check for nested block
	if p.isKeyword("DECLARE") || p.isKeyword("BEGIN") {
		block, err := p.parseBlock(label)
		if err != nil {
			return nil, err
		}
		// Consume trailing semicolon for nested blocks
		if p.cur.Type == ';' {
			p.advance()
		}
		return block, nil
	}

	// Labeled loops
	if label != "" {
		switch {
		case p.isKeyword("LOOP"):
			return p.parseLoop(label)
		case p.isKeyword("WHILE"):
			return p.parseWhile(label)
		case p.isKeyword("FOR"):
			return p.parseFor(label)
		case p.isKeyword("FOREACH"):
			return p.parseForEach(label)
		default:
			return nil, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
		}
	}

	// Unlabeled statement dispatch
	switch {
	case p.isKeyword("IF"):
		return p.parseIf()
	case p.isKeyword("CASE"):
		return p.parseCase()
	case p.isKeyword("LOOP"):
		return p.parseLoop("")
	case p.isKeyword("WHILE"):
		return p.parseWhile("")
	case p.isKeyword("FOR"):
		return p.parseFor("")
	case p.isKeyword("FOREACH"):
		return p.parseForEach("")
	case p.isKeyword("EXIT"):
		return p.parseExit()
	case p.isKeyword("CONTINUE"):
		return p.parseContinue()
	case p.isKeyword("CLOSE"):
		return p.parseClose()
	case p.isKeyword("GET"):
		return p.parseGetDiag()
	case p.isKeyword("NULL"):
		return p.parseNull()
	case p.isKeyword("ASSERT"):
		return p.parseAssert()
	case p.isKeyword("CALL"):
		return p.parseCall()
	case p.isKeyword("DO"):
		return p.parseDo()
	case p.isKeyword("COMMIT"):
		return p.parseCommit()
	case p.isKeyword("ROLLBACK"):
		return p.parseRollback()
	case p.isKeyword("RETURN"):
		return p.parseReturn()
	case p.isKeyword("PERFORM"):
		return p.parsePerform()
	case p.isKeyword("EXECUTE"):
		return p.parseDynExecute()
	case p.isKeyword("RAISE"):
		return p.parseRaise()
	case p.isKeyword("OPEN"):
		return p.parseOpen()
	case p.isKeyword("FETCH"):
		return p.parseFetch()
	case p.isKeyword("MOVE"):
		return p.parseMove()
	case p.cur.Type == pgparser.INSERT || p.cur.Type == pgparser.UPDATE ||
		p.cur.Type == pgparser.DELETE_P || p.cur.Type == pgparser.SELECT ||
		p.cur.Type == pgparser.MERGE || p.isKeyword("IMPORT"):
		return p.parseExecSQL()
	}

	// Check for assignment: identifier followed by := or = or [ or .
	if p.isIdent() || p.isAnyKeywordAsIdent() {
		return p.parseAssignOrCall()
	}

	return nil, p.errorf("syntax error at or near %q", p.tokenText(p.cur))
}

// parseIf, parseCase are implemented in control.go
// parseLoop, parseWhile, parseFor, parseForEach are implemented in loop.go
// parseRaise is implemented in raise.go
// parseExit, parseContinue are implemented in exit.go
// parseReturn, parsePerform, parseExecSQL are implemented in stmt.go
// parseDynExecute is implemented in execute.go

// parseNull parses: NULL ;
func (p *Parser) parseNull() (ast.Node, error) {
	startPos := p.pos()
	p.advance() // consume NULL
	if p.cur.Type == ';' {
		p.advance()
	}
	return &ast.PLNull{
		Loc: ast.Loc{Start: startPos, End: p.prev.End},
	}, nil
}

// Parse parses a PL/pgSQL function body and returns the root block.
func Parse(body string) (*ast.PLBlock, error) {
	p := newParser(body)
	return p.parse()
}

// parseAssignOrCall, parseReturn, parsePerform, parseExecSQL are implemented in stmt.go
// parseDynExecute is implemented in execute.go

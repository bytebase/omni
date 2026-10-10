// Package parser implements a recursive descent SQL parser for Oracle PL/SQL.
//
// This parser produces AST nodes from the oracle/ast package.
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// Parser is a recursive descent parser for Oracle SQL/PL/SQL.
//
// Oracle parser helpers follow the project-wide dual-return contract: required
// parse functions return (Node, error), optional probes use a tryParse or
// parseOptional name and may return (nil, nil), and nested parse errors must be
// returned to the caller instead of stored or discarded.
type Parser struct {
	lexer   *Lexer
	source  string // original SQL input
	cur     Token  // current token
	prev    Token  // previous token (for error reporting)
	nextBuf Token  // buffered next token for 2-token lookahead
	hasNext bool   // whether nextBuf is valid
}

// Parse parses a SQL string into an AST list.
// Each statement is wrapped in a *RawStmt.
func Parse(sql string) (*nodes.List, error) {
	return ParseRange(sql, 0, len(sql))
}

// ParseRange parses source[start:end] in place. Every Loc and error Position
// is an offset into source, so a caller that split a script can parse each
// segment with absolute positions without copying or padding the text.
func ParseRange(source string, start, end int) (*nodes.List, error) {
	if err := validateBalancedDelimiters(source, start, end); err != nil {
		return nil, err
	}

	p := &Parser{
		lexer:  NewLexerRange(source, start, end),
		source: source,
	}
	p.advance()

	var stmts []nodes.Node
	needSeparator := false
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			needSeparator = false
			continue
		}
		if needSeparator {
			return nil, p.syntaxErrorAtCur()
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}

		stmtLoc := p.pos()
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}
		if p.cur.Type == tokEOF && (p.isIncompleteStatementEnd(stmt) || isIncompleteParsedStmt(stmt)) {
			return nil, p.syntaxErrorAtCur()
		}
		if stmt == nil {
			return nil, p.syntaxErrorAtCur()
		}

		raw := &nodes.RawStmt{
			Stmt: stmt,
			Loc:  nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}
		stmts = append(stmts, raw)
		needSeparator = true
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

func validateBalancedDelimiters(sql string, start, end int) error {
	lexer := NewLexerRange(sql, start, end)
	depth := 0
	var javaHead embeddedSourceHead
	var prev1, prev2, prev3 Token
	for {
		tok := lexer.NextToken()
		if lexer.Err != nil {
			return nil
		}
		if javaHead.observe(tok) {
			// The rest of the range is Java or JavaScript source, not SQL.
			return nil
		}
		if (prev3.Type == kwIS || prev3.Type == kwAS) &&
			prev2.Type == tokIDENT && prev2.Str == "MLE" && prev1.Type == tokIDENT && prev1.Str == "LANGUAGE" {
			// tok names the language of an inline MLE call spec; its
			// delimited code is not SQL.
			if _, codeEnd, _, ok := mleInlineCodeEnd(sql, tok.End, end); ok {
				lexer.pos = codeEnd
			}
		}
		prev3, prev2, prev1 = prev2, prev1, tok
		switch tok.Type {
		case ';':
			javaHead = embeddedSourceHead{}
		case tokEOF:
			if depth > 0 {
				return &ParseError{Message: "syntax error at end of input", Position: tok.Loc}
			}
			return nil
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				text := ")"
				return &ParseError{Message: fmt.Sprintf("syntax error at or near %q", text), Position: tok.Loc}
			}
		}
	}
}

func (p *Parser) isIncompleteStatementEnd(stmt nodes.StmtNode) bool {
	switch p.prev.Type {
	case '(', ',', '.', '=', '+', '-', '*', '/', tokASSIGN, tokASSOC, tokCONCAT:
		return true
	}

	// Keywords Oracle does not reserve (CONTENT, JOIN, USING) are left out: a
	// statement can end with a column of that name, and a clause they open
	// without its operand ("FROM t JOIN") fails in the clause's own parser.
	switch stmt.(type) {
	case *nodes.SelectStmt:
		switch p.prev.Type {
		case kwAND, kwAS, kwBY, kwFROM, kwGROUP, kwHAVING,
			kwIS, kwNOT, kwON, kwOR, kwORDER, kwTHEN, kwUNION,
			kwWHERE:
			return true
		}
	case *nodes.InsertStmt:
		return p.prev.Type == kwINTO || p.prev.Type == kwVALUES
	case *nodes.UpdateStmt:
		return p.prev.Type == kwUPDATE || p.prev.Type == kwSET || p.prev.Type == kwWHERE
	case *nodes.DeleteStmt:
		return p.prev.Type == kwFROM || p.prev.Type == kwWHERE
	case *nodes.MergeStmt:
		switch p.prev.Type {
		case kwINTO, kwON, kwSET, kwTHEN, kwWHERE:
			return true
		}
	case *nodes.CreateTableStmt:
		return p.prev.Type == kwTABLE || p.prev.Type == kwDEFAULT
	case *nodes.CreateIndexStmt:
		return p.prev.Type == kwINDEX
	case *nodes.CreateViewStmt:
		return p.prev.Type == kwAS
	case *nodes.AlterTableStmt:
		return p.prev.Type == kwTABLE || p.prev.Type == kwADD || p.prev.Type == kwMODIFY
	case *nodes.CreateProcedureStmt, *nodes.CreateFunctionStmt, *nodes.CreatePackageStmt,
		*nodes.CreateTriggerStmt, *nodes.PLSQLBlock:
		return p.prev.Type == kwBEGIN || p.prev.Type == kwTHEN
	case *nodes.CreateUserStmt:
		return p.prev.Type == kwUSER
	case *nodes.GrantStmt:
		return p.prev.Type == kwON || p.prev.Type == kwTO
	case *nodes.DropStmt:
		return p.prev.Type == kwDROP
	}

	return false
}

func isIncompleteParsedStmt(stmt nodes.StmtNode) bool {
	switch s := stmt.(type) {
	case *nodes.SelectStmt:
		return (s.FromClause == nil || s.FromClause.Len() == 0) && selectTargetListContainsStar(s.TargetList)
	case *nodes.DropStmt:
		if s.ObjectType == nodes.OBJECT_DATABASE && (s.Names == nil || s.Names.Len() == 0) {
			return false
		}
		if s.Names == nil || s.Names.Len() == 0 {
			return true
		}
		for _, item := range s.Names.Items {
			name, ok := item.(*nodes.ObjectName)
			if ok && name.Name == "" {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func selectTargetListContainsStar(list *nodes.List) bool {
	if list == nil {
		return false
	}
	for _, item := range list.Items {
		target, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if _, ok := target.Expr.(*nodes.Star); ok {
			return true
		}
	}
	return false
}

// parseStmt dispatches to statement-specific parsers.
// Statement parsers must return syntax and lexer errors directly so Parse can
// fail at the earliest invalid token while preserving partial-node Loc data.
func (p *Parser) parseStmt() (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwSELECT, kwWITH:
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwMERGE:
		return p.parseMergeStmt()
	case kwDROP:
		return p.parseDropStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSAVEPOINT:
		return p.parseSavepointStmt()
	case kwSET:
		next := p.peekNext()
		if next.Type == kwTRANSACTION {
			start := p.pos()
			p.advance() // consume SET
			p.advance() // consume TRANSACTION
			stmt, parseErr828 := p.parseSetTransactionStmt()
			if parseErr828 != nil {
				return nil, parseErr828
			}
			stmt.(*nodes.SetTransactionStmt).Loc.Start = start
			return stmt, nil
		}
		if next.Type == kwROLE {
			return p.parseSetRoleStmt()
		}
		if next.Type == kwCONSTRAINT || next.Type == kwCONSTRAINTS {
			return p.parseSetConstraintsStmt()
		}
		return nil, p.syntaxErrorAtCur()
	case kwAUDIT:
		return p.parseAuditStmt()
	case kwNOAUDIT:
		return p.parseNoauditStmt()
	case kwASSOCIATE:
		return p.parseAssociateStatisticsStmt()
	case kwDISASSOCIATE:
		return p.parseDisassociateStatisticsStmt()
	case kwCOMMENT:
		return p.parseCommentStmt()
	case kwTRUNCATE:
		return p.parseTruncateStmt()
	case kwANALYZE:
		return p.parseAnalyzeStmt()
	case kwEXPLAIN:
		return p.parseExplainPlanStmt()
	case kwFLASHBACK:
		next := p.peekNext()
		if next.Type == kwDATABASE ||
			(next.Type == tokIDENT && (next.Str == "STANDBY" || next.Str == "PLUGGABLE")) {
			return p.parseFlashbackDatabaseStmt()
		}
		return p.parseFlashbackTableStmt()
	case kwPURGE:
		return p.parsePurgeStmt()
	case kwLOCK:
		return p.parseLockTableStmt()
	case kwCALL:
		return p.parseCallStmt()
	case kwRENAME:
		return p.parseRenameStmt()
	case kwALTER:
		return p.parseAlterStmt()
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()
	case kwCREATE:
		return p.parseCreateStmt()
	case kwDECLARE, kwBEGIN:
		return p.parsePLSQLBlock()
	case tokLABELOPEN:
		return p.parsePLSQLBlock()
	default:
		// Handle identifier-based statements
		if p.isIdentLike() {
			switch p.cur.Str {
			case "ADMINISTER":
				return p.parseAdministerKeyManagementStmt()
			}
		}
		return nil, p.syntaxErrorAtCur()
	}
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

// peekAhead returns the n tokens after cur without consuming them; tokens
// past the end of input are tokEOF. It lexes a copy of the lexer, so the
// parser's position is untouched.
func (p *Parser) peekAhead(n int) []Token {
	toks := make([]Token, 0, n)
	if n == 0 {
		return toks
	}
	toks = append(toks, p.peekNext())
	lexer := *p.lexer
	for len(toks) < n {
		toks = append(toks, lexer.NextToken())
	}
	return toks
}

// parserState is a snapshot of everything the parser advances, so a probe can
// parse ahead and then put the parser back where it was.
type parserState struct {
	lexer   Lexer
	cur     Token
	prev    Token
	nextBuf Token
	hasNext bool
}

func (p *Parser) saveState() parserState {
	return parserState{lexer: *p.lexer, cur: p.cur, prev: p.prev, nextBuf: p.nextBuf, hasNext: p.hasNext}
}

func (p *Parser) restoreState(s parserState) {
	*p.lexer = s.lexer
	p.cur, p.prev, p.nextBuf, p.hasNext = s.cur, s.prev, s.nextBuf, s.hasNext
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
	return Token{}, p.syntaxErrorAtCur()
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// isKeyword checks whether the current token is a specific keyword.
func (p *Parser) isKeyword(kw int) bool {
	return p.cur.Type == kw
}

// matchKeyword consumes the token if it is the given keyword.
func (p *Parser) matchKeyword(kw int) (Token, bool) {
	if p.cur.Type == kw {
		return p.advance(), true
	}
	return Token{}, false
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Severity string // e.g., "ERROR", "WARNING"; defaults to "ERROR"
	Code     string // SQLSTATE code; defaults to "42601"
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	sev := e.Severity
	if sev == "" {
		sev = "ERROR"
	}
	code := e.Code
	if code == "" {
		code = "42601"
	}
	return fmt.Sprintf("%s: %s (SQLSTATE %s)", sev, e.Message, code)
}

// tokenText returns the original source text covered by tok.
func (p *Parser) tokenText(tok Token) string {
	if tok.Type == tokEOF {
		return ""
	}
	if tok.Loc >= 0 && tok.End >= tok.Loc && tok.End <= len(p.source) {
		return p.source[tok.Loc:tok.End]
	}
	if tok.Str != "" {
		return tok.Str
	}
	if tok.Type > 0 && tok.Type < 256 {
		return string(rune(tok.Type))
	}
	return ""
}

// syntaxErrorAtCur returns a syntax error for the current token.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	if p.lexer.Err != nil {
		return p.lexerError()
	}
	return p.syntaxErrorAtTok(p.cur)
}

// syntaxErrorAtTok returns a PG-style syntax error for tok.
func (p *Parser) syntaxErrorAtTok(tok Token) *ParseError {
	text := p.tokenText(tok)
	msg := "syntax error at end of input"
	if text != "" {
		msg = fmt.Sprintf("syntax error at or near %q", text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
	}
}

// lexerError returns the lexer's error with token context.
func (p *Parser) lexerError() *ParseError {
	tok := p.cur
	if tok.Type == tokEOF && p.prev.Type != 0 {
		tok = p.prev
	}
	text := p.tokenText(tok)
	msg := p.lexer.Err.Error()
	if text != "" {
		msg = fmt.Sprintf("%s at or near %q", msg, text)
	}
	return &ParseError{
		Message:  msg,
		Position: tok.Loc,
	}
}

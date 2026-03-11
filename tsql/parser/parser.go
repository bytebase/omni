// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
//
// This parser reuses the lexer and keyword definitions from this package
// and produces AST nodes from the tsql/ast package.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// Parser is a recursive descent parser for T-SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid
}

// Parse parses a T-SQL string into an AST list.
// Currently supports basic infrastructure; statement dispatch will be
// implemented incrementally across batches.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{
		lexer: NewLexer(sql),
	}
	p.advance()

	var stmts []nodes.Node
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		stmt := p.parseStmt()
		if stmt == nil {
			if p.cur.Type != tokEOF {
				return nil, &ParseError{
					Message:  "unexpected token in statement",
					Position: p.cur.Loc,
				}
			}
			break
		}
		stmts = append(stmts, stmt)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// parseStmt dispatches to statement-specific parsers.
func (p *Parser) parseStmt() nodes.StmtNode {
	switch p.cur.Type {
	case kwSELECT:
		return p.parseSelectStmt()
	case kwWITH:
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwUPDATE:
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwMERGE:
		return p.parseMergeStmt()
	case kwCREATE:
		return p.parseCreateStmt()
	case kwALTER:
		return p.parseAlterStmt()
	case kwDROP:
		return p.parseDropStmt()
	case kwTRUNCATE:
		return p.parseTruncateStmt()
	case kwIF:
		return p.parseIfStmt()
	case kwWHILE:
		return p.parseWhileStmt()
	case kwBEGIN:
		return p.parseBeginStmt()
	case kwRETURN:
		return p.parseReturnStmt()
	case kwBREAK:
		return p.parseBreakStmt()
	case kwCONTINUE:
		return p.parseContinueStmt()
	case kwGOTO:
		return p.parseGotoStmt()
	case kwWAITFOR:
		return p.parseWaitForStmt()
	case kwDECLARE:
		return p.parseDeclareStmt()
	case kwSET:
		return p.parseSetStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSAVE:
		return p.parseSaveTransStmt()
	case kwEXEC, kwEXECUTE:
		return p.parseExecStmt()
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()
	case kwDENY:
		return p.parseDenyStmt()
	case kwUSE:
		return p.parseUseStmt()
	case kwPRINT:
		return p.parsePrintStmt()
	case kwRAISERROR:
		return p.parseRaiseErrorStmt()
	case kwTHROW:
		return p.parseThrowStmt()
	case kwGO:
		return p.parseGoStmt()
	default:
		// Check for label: identifier followed by ':'
		if p.isIdentLike() {
			next := p.peekNext()
			if next.Type == ':' {
				return p.parseLabelStmt()
			}
		}
		return nil
	}
}

// parseCreateStmt dispatches CREATE to the appropriate sub-parser.
func (p *Parser) parseCreateStmt() nodes.StmtNode {
	loc := p.pos()
	p.advance() // consume CREATE

	// OR ALTER
	orAlter := false
	if p.cur.Type == kwOR {
		next := p.peekNext()
		if next.Type == kwALTER {
			p.advance() // OR
			p.advance() // ALTER
			orAlter = true
		}
	}

	// UNIQUE for CREATE UNIQUE INDEX
	unique := false
	if p.cur.Type == kwUNIQUE {
		unique = true
		p.advance()
	}

	switch p.cur.Type {
	case kwTABLE:
		p.advance() // consume TABLE
		stmt := p.parseCreateTableStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwINDEX, kwCLUSTERED, kwNONCLUSTERED, kwCOLUMNSTORE:
		stmt := p.parseCreateIndexStmt(unique)
		stmt.Loc.Start = loc
		return stmt
	case kwVIEW:
		p.advance() // consume VIEW
		stmt := p.parseCreateViewStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwPROCEDURE, kwPROC:
		p.advance() // consume PROCEDURE/PROC
		stmt := p.parseCreateProcedureStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwFUNCTION:
		p.advance() // consume FUNCTION
		stmt := p.parseCreateFunctionStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwDATABASE:
		p.advance() // consume DATABASE
		stmt := p.parseCreateDatabaseStmt()
		stmt.Loc.Start = loc
		return stmt
	default:
		return nil
	}
}

// parseAlterStmt dispatches ALTER to the appropriate sub-parser.
func (p *Parser) parseAlterStmt() nodes.StmtNode {
	loc := p.pos()
	p.advance() // consume ALTER

	switch p.cur.Type {
	case kwTABLE:
		p.advance() // consume TABLE
		stmt := p.parseAlterTableStmt()
		stmt.Loc.Start = loc
		return stmt
	default:
		return nil
	}
}

// parseLabelStmt parses a label: statement.
func (p *Parser) parseLabelStmt() *nodes.LabelStmt {
	loc := p.pos()
	label := p.cur.Str
	p.advance() // consume identifier
	p.advance() // consume :
	return &nodes.LabelStmt{
		Label: label,
		Loc:   nodes.Loc{Start: loc, End: p.pos()},
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
	return Token{}, &ParseError{
		Message:  "unexpected token",
		Position: p.cur.Loc,
	}
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}

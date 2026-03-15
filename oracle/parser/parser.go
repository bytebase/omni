// Package parser implements a recursive descent SQL parser for Oracle PL/SQL.
//
// This parser produces AST nodes from the oracle/ast package.
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// Parser is a recursive descent parser for Oracle SQL/PL/SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid
}

// Parse parses a SQL string into an AST list.
// Each statement is wrapped in a *RawStmt.
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

		stmtLoc := p.pos()
		stmt := p.parseStmt()
		if stmt == nil {
			if p.cur.Type != tokEOF {
				return nil, &ParseError{
					Message:  fmt.Sprintf("unexpected token %q at position %d", p.cur.Str, p.cur.Loc),
					Position: p.cur.Loc,
				}
			}
			break
		}

		raw := &nodes.RawStmt{
			Stmt:         stmt,
			StmtLocation: stmtLoc,
			StmtLen:      p.pos() - stmtLoc,
		}
		stmts = append(stmts, raw)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// parseStmt dispatches to statement-specific parsers.
// Each batch wires in its statement parsers incrementally.
func (p *Parser) parseStmt() nodes.StmtNode {
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
			stmt := p.parseSetTransactionStmt()
			stmt.(*nodes.SetTransactionStmt).Loc.Start = start
			return stmt
		}
		if next.Type == kwROLE {
			return p.parseSetRoleStmt()
		}
		if next.Type == kwCONSTRAINT || next.Type == kwCONSTRAINTS {
			return p.parseSetConstraintsStmt()
		}
		return nil
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
		return nil
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
	return Token{}, &ParseError{
		Message:  fmt.Sprintf("expected token type %d, got %d (%q)", tokenType, p.cur.Type, p.cur.Str),
		Position: p.cur.Loc,
	}
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
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}

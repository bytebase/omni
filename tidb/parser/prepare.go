package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parsePrepareStmt parses a PREPARE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/prepare.html
//
//	PREPARE stmt_name FROM preparable_stmt
func (p *Parser) parsePrepareStmt() (*nodes.PrepareStmt, error) {
	start := p.pos()
	p.advance() // consume PREPARE

	stmt := &nodes.PrepareStmt{Loc: nodes.Loc{Start: start}}

	// Statement name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// FROM
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// SQL string or variable — mandatory.
	//
	// Ref: PreparedStmt in TiDB v8.5.5 pkg/parser/parser.y is
	// `"PREPARE" Identifier "FROM" PrepareSQL`, and PrepareSQL is only
	// `stringLit | UserVariable`; FROM with no value is invalid.
	if p.cur.Type == tokSCONST {
		stmt.Stmt = p.cur.Str
		p.advance()
	} else if p.isVariableRef() {
		vref, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = "@" + vref.Name
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseExecuteStmt parses an EXECUTE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/execute.html
//
//	EXECUTE stmt_name [USING @var_name [, @var_name] ...]
func (p *Parser) parseExecuteStmt() (*nodes.ExecuteStmt, error) {
	start := p.pos()
	p.advance() // consume EXECUTE

	stmt := &nodes.ExecuteStmt{Loc: nodes.Loc{Start: start}}

	// Statement name
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional USING
	// USING takes user variables and nothing else. Ref: ExecuteStmt's
	// "USING" UserVariableList in tidb/testdata/parser.y.
	if p.cur.Type == kwUSING {
		p.advance()
		for {
			if !p.atUsingVariable() {
				return nil, p.syntaxErrorAtCur()
			}
			ref, err := p.parseVariableRef()
			if err != nil {
				return nil, err
			}
			stmt.Params = append(stmt.Params, ref)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDeallocateStmt parses a DEALLOCATE PREPARE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/deallocate-prepare.html
//
//	{DEALLOCATE | DROP} PREPARE stmt_name
func (p *Parser) parseDeallocateStmt() (*nodes.DeallocateStmt, error) {
	start := p.pos()
	p.advance() // consume DEALLOCATE

	// PREPARE is mandatory.
	//
	// Ref: DeallocateStmt in TiDB v8.5.5 pkg/parser/parser.y is
	// `DeallocateSym "PREPARE" Identifier`.
	if _, err := p.expect(kwPREPARE); err != nil {
		return nil, err
	}

	stmt := &nodes.DeallocateStmt{Loc: nodes.Loc{Start: start}}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// atUsingVariable reports whether p.cur is a user variable as singleAtIdentifier
// spells one, the unnamed one included: a bare @, and an @ before two quotes or
// two backquotes, which the lexer folds into a token whose text is "@" alone. The
// source text, which starts at the @, tells it from a backquoted name.
func (p *Parser) atUsingVariable() bool {
	return p.cur.Type == tokIDENT && strings.HasPrefix(p.cur.Str, "@") &&
		strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "@")
}

package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// Pipes (StarRocks.g4:1913-1939). PIPE is nonReserved, so DESC PIPE is told
// from a table named pipe by the name that follows it.

// parsePipeName reads a pipe's name. The AST builder refuses a name of more
// than a database and a pipe ("Invalid pipe name 'a.b.c'"), so a third part is
// a syntax error here.
func (p *Parser) parsePipeName() (*ast.ObjectName, error) {
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	if len(name.Parts) > 2 {
		return nil, &ParseError{
			Position: name.Loc.Start, End: name.Loc.End,
			Message: "a pipe's name is a database and a name",
		}
	}
	return name, nil
}

// parseCreatePipe parses createPipeStatement (StarRocks.g4:1913):
//
//	CREATE [OR REPLACE] PIPE [IF NOT EXISTS] qualifiedName properties? AS insertStatement
//
// On entry CREATE has been consumed; cur is PIPE.
func (p *Parser) parseCreatePipe(startLoc ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume PIPE

	stmt := &ast.CreatePipeStmt{OrReplace: orReplace}
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	name, err := p.parsePipeName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	endLoc := name.Loc

	// The AST builder refuses the pair: "Conflicted options OR REPLACE and
	// IF NOT EXISTS".
	if stmt.OrReplace && stmt.IfNotExists {
		return nil, &ParseError{
			Position: startLoc.Start, End: endLoc.End,
			Message: "OR REPLACE and IF NOT EXISTS conflict",
		}
	}

	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		endLoc = p.prev.Loc // the closing ')'
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	// The body is insertStatement (:2190), which takes an explainDesc prefix:
	// `CREATE PIPE p AS EXPLAIN INSERT INTO t SELECT 1` is accepted on the
	// pinned server. Whatever the prefix, the body must wrap an insert.
	var insert ast.Node
	switch p.cur.Kind {
	case kwEXPLAIN, kwDESC, kwDESCRIBE:
		explained, err := p.parseExplain()
		if err != nil {
			return nil, err
		}
		exp, ok := explained.(*ast.ExplainStmt)
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		if _, ok := exp.Query.(*ast.InsertStmt); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		insert = exp
	case kwINSERT:
		node, err := p.parseInsert()
		if err != nil {
			return nil, err
		}
		insert = node
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Insert = insert
	endLoc = ast.NodeLoc(insert)

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseDropPipe parses dropPipeStatement (StarRocks.g4:1919):
//
//	DROP PIPE [IF EXISTS] qualifiedName
//
// On entry DROP has been consumed; cur is PIPE.
func (p *Parser) parseDropPipe(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PIPE
	stmt := &ast.DropPipeStmt{}
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}
	name, err := p.parsePipeName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc = startLoc.Merge(name.Loc)
	return stmt, nil
}

// parseAlterPipe parses alterPipeStatement (StarRocks.g4:1923) with its
// alterPipeClause (:1927):
//
//	ALTER PIPE qualifiedName
//	    { SUSPEND | RESUME | RETRY ALL | RETRY FILE string | SET propertyList }
//
// On entry ALTER has been consumed; cur is PIPE.
func (p *Parser) parseAlterPipe(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PIPE
	name, err := p.parsePipeName()
	if err != nil {
		return nil, err
	}
	stmt := &ast.AlterPipeStmt{Name: name}
	endLoc := name.Loc

	switch p.cur.Kind {
	case kwSUSPEND:
		stmt.Action = "SUSPEND"
		endLoc = p.cur.Loc
		p.advance()
	case kwRESUME:
		stmt.Action = "RESUME"
		endLoc = p.cur.Loc
		p.advance()
	case kwRETRY:
		p.advance() // consume RETRY
		switch p.cur.Kind {
		case kwALL:
			stmt.Action = "RETRY ALL"
			endLoc = p.cur.Loc
			p.advance()
		case kwFILE:
			p.advance() // consume FILE
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = "RETRY FILE"
			stmt.File = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwSET:
		p.advance() // consume SET
		props, err := p.parsePropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Action = "SET"
		stmt.Properties = props
		endLoc = p.prev.Loc // the closing ')'
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// atDescPipe reports whether a DESC/DESCRIBE at cur describes a pipe: PIPE
// followed by the name, since `DESC pipe` alone describes a table named pipe.
func (p *Parser) atDescPipe() bool {
	if p.peekNext().Kind != kwPIPE {
		return false
	}
	mark := p.save()
	p.advance() // consume DESC/DESCRIBE
	p.advance() // consume PIPE
	ok := isIdentifierToken(p.cur.Kind)
	p.restore(mark)
	return ok
}

// parseDescPipe parses descPipeStatement (StarRocks.g4:1935):
//
//	(DESC | DESCRIBE) PIPE qualifiedName
//
// On entry the verb has been consumed; cur is PIPE.
func (p *Parser) parseDescPipe(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PIPE
	name, err := p.parsePipeName()
	if err != nil {
		return nil, err
	}
	return &ast.DescPipeStmt{Name: name, Loc: startLoc.Merge(name.Loc)}, nil
}

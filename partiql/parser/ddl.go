package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parseCreateCommand parses the two CREATE variants:
//
//	CREATE TABLE symbolPrimitive
//	CREATE INDEX ON symbolPrimitive ( pathSimple (, pathSimple)* )
//
// Grammar: createCommand (PartiQLParser.g4 lines 78-81).
func (p *Parser) parseCreateCommand() (ast.StmtNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume CREATE

	switch p.cur.Type {
	case tokTABLE:
		p.advance() // consume TABLE
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.CreateTableStmt{
			Name: &ast.VarRef{
				Name:          name,
				CaseSensitive: caseSensitive,
				Loc:           nameLoc,
			},
			Loc: ast.Loc{Start: start, End: nameLoc.End},
		}, nil

	case tokINDEX:
		p.advance() // consume INDEX
		if _, err := p.expect(tokON); err != nil {
			return nil, err
		}
		tableName, tableCaseSensitive, tableLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		tableRef := &ast.VarRef{
			Name:          tableName,
			CaseSensitive: tableCaseSensitive,
			Loc:           tableLoc,
		}
		if _, err := p.expect(tokPAREN_LEFT); err != nil {
			return nil, &ParseError{
				Message: "expected PAREN_LEFT after table name in CREATE INDEX",
				Loc:     p.cur.Loc,
			}
		}
		var paths []*ast.PathExpr
		path, err := p.parsePathSimple()
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
		for p.cur.Type == tokCOMMA {
			p.advance() // consume ,
			path, err := p.parsePathSimple()
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
		}
		rp, err := p.expect(tokPAREN_RIGHT)
		if err != nil {
			return nil, err
		}
		return &ast.CreateIndexStmt{
			Table: tableRef,
			Paths: paths,
			Loc:   ast.Loc{Start: start, End: rp.Loc.End},
		}, nil

	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected TABLE or INDEX after CREATE, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
}

// parseDropCommand parses the two DROP variants:
//
//	DROP TABLE symbolPrimitive
//	DROP INDEX symbolPrimitive ON symbolPrimitive
//
// Grammar: dropCommand (PartiQLParser.g4 lines 83-86).
func (p *Parser) parseDropCommand() (ast.StmtNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume DROP

	switch p.cur.Type {
	case tokTABLE:
		p.advance() // consume TABLE
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.DropTableStmt{
			Name: &ast.VarRef{
				Name:          name,
				CaseSensitive: caseSensitive,
				Loc:           nameLoc,
			},
			Loc: ast.Loc{Start: start, End: nameLoc.End},
		}, nil

	case tokINDEX:
		p.advance() // consume INDEX
		idxName, idxCaseSensitive, idxLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokON); err != nil {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected ON after index name in DROP INDEX, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		tableName, tableCaseSensitive, tableLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.DropIndexStmt{
			Index: &ast.VarRef{
				Name:          idxName,
				CaseSensitive: idxCaseSensitive,
				Loc:           idxLoc,
			},
			Table: &ast.VarRef{
				Name:          tableName,
				CaseSensitive: tableCaseSensitive,
				Loc:           tableLoc,
			},
			Loc: ast.Loc{Start: start, End: tableLoc.End},
		}, nil

	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected TABLE or INDEX after DROP, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
}

// parsePathSimple parses a simple path used in CREATE INDEX column lists:
//
//	symbolPrimitive pathSimpleSteps*
//
// where pathSimpleSteps are:
//
//	BRACKET_LEFT key=literal BRACKET_RIGHT
//	BRACKET_LEFT key=symbolPrimitive BRACKET_RIGHT
//	PERIOD key=symbolPrimitive
//
// Grammar: pathSimple, pathSimpleSteps (PartiQLParser.g4 lines 110-117).
func (p *Parser) parsePathSimple() (*ast.PathExpr, error) {
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	root := &ast.VarRef{
		Name:          name,
		CaseSensitive: caseSensitive,
		Loc:           nameLoc,
	}

	var steps []ast.PathStep
	end := nameLoc.End

loop:
	for {
		switch p.cur.Type {
		case tokPERIOD:
			stepStart := p.cur.Loc.Start
			p.advance() // consume .
			fieldName, fieldCS, fieldLoc, err := p.parseSymbolPrimitive()
			if err != nil {
				return nil, err
			}
			steps = append(steps, &ast.DotStep{
				Field:         fieldName,
				CaseSensitive: fieldCS,
				Loc:           ast.Loc{Start: stepStart, End: fieldLoc.End},
			})
			end = fieldLoc.End

		case tokBRACKET_LEFT:
			stepStart := p.cur.Loc.Start
			p.advance() // consume [
			// Try literal first, then symbolPrimitive.
			var idx ast.ExprNode
			if p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
				lit, err := p.parseLiteral()
				if err != nil {
					return nil, err
				}
				idx = lit
			} else {
				// symbolPrimitive as identifier reference inside bracket.
				keyName, keyCS, keyLoc, err := p.parseSymbolPrimitive()
				if err != nil {
					return nil, err
				}
				idx = &ast.VarRef{
					Name:          keyName,
					CaseSensitive: keyCS,
					Loc:           keyLoc,
				}
			}
			rb, err := p.expect(tokBRACKET_RIGHT)
			if err != nil {
				return nil, err
			}
			steps = append(steps, &ast.IndexStep{
				Index: idx,
				Loc:   ast.Loc{Start: stepStart, End: rb.Loc.End},
			})
			end = rb.Loc.End

		default:
			break loop
		}
	}

	return &ast.PathExpr{
		Root:  root,
		Steps: steps,
		Loc:   ast.Loc{Start: nameLoc.Start, End: end},
	}, nil
}

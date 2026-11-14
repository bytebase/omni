package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parseCreateCommand parses the two CREATE variants:
//
//	CREATE TABLE qualifiedName ( PAREN_LEFT tableDef PAREN_RIGHT )?
//	CREATE INDEX ON symbolPrimitive ( pathSimple (, pathSimple)* )
//
// Grammar: createCommand (PartiQLParser.g4:87-90).
func (p *Parser) parseCreateCommand() (ast.StmtNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume CREATE

	switch p.cur.Type {
	case tokTABLE:
		p.advance() // consume TABLE
		qualifier, name, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		stmt := &ast.CreateTableStmt{
			Name:      name,
			Qualifier: qualifier,
			Loc:       ast.Loc{Start: start, End: name.Loc.End},
		}
		if p.cur.Type == tokPAREN_LEFT {
			cols, end, err := p.parseTableDef()
			if err != nil {
				return nil, err
			}
			stmt.Columns = cols
			stmt.Loc.End = end
		}
		return stmt, nil

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
				Message:  "expected PAREN_LEFT after table name in CREATE INDEX",
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
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
			Message:  fmt.Sprintf("expected TABLE or INDEX after CREATE, got %q", p.cur.Str),
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		}
	}
}

// parseDropCommand parses the two DROP variants:
//
//	DROP TABLE qualifiedName
//	DROP INDEX symbolPrimitive ON symbolPrimitive
//
// Grammar: dropCommand (PartiQLParser.g4:92-95).
func (p *Parser) parseDropCommand() (ast.StmtNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume DROP

	switch p.cur.Type {
	case tokTABLE:
		p.advance() // consume TABLE
		qualifier, name, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		return &ast.DropTableStmt{
			Name:      name,
			Qualifier: qualifier,
			Loc:       ast.Loc{Start: start, End: name.Loc.End},
		}, nil

	case tokINDEX:
		p.advance() // consume INDEX
		idxName, idxCaseSensitive, idxLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokON); err != nil {
			return nil, &ParseError{
				Message:  fmt.Sprintf("expected ON after index name in DROP INDEX, got %q", p.cur.Str),
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
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
			Message:  fmt.Sprintf("expected TABLE or INDEX after DROP, got %q", p.cur.Str),
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
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
			// pathSimpleSteps has two bracket alternatives
			// (PartiQLParser.g4:141-142): `[ key=literal ]` and
			// `[ key=symbolPrimitive ]`. Every literal-start token goes
			// through parseLiteral, and a bare identifier is a symbol. The
			// grammar takes any literal, and the reference parser's visitor
			// takes a number or a string and refuses any other
			// (PartiQLParserDefault.kt:859-876, "Expected an integer or
			// string literal").
			var idx ast.ExprNode
			if isLiteralStart(p.cur.Type) {
				lit, err := p.parseLiteral()
				if err != nil {
					return nil, err
				}
				switch lit.(type) {
				case *ast.NumberLit, *ast.StringLit:
				default:
					return nil, &ParseError{
						Message:  "expected an integer or string literal in a path step",
						Position: lit.GetLoc().Start, End: lit.GetLoc().End,
					}
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

// isLiteralStart reports whether tok can begin a `literal`
// (PartiQLParser.g4:775-787): NULL, MISSING, TRUE, FALSE, string, integer,
// decimal, Ion, or a DATE/TIME/TIMESTAMP literal. It mirrors the literal
// cases of parsePrimaryBase's dispatch (exprprimary.go) so the bracket-key
// path step in parsePathSimple reads the grammar's `literal` set.
func isLiteralStart(tok int) bool {
	switch tok {
	case tokNULL, tokMISSING, tokTRUE, tokFALSE,
		tokSCONST, tokICONST, tokFCONST, tokION_LITERAL,
		tokDATE, tokTIME, tokTIMESTAMP:
		return true
	default:
		return false
	}
}

// parseQualifiedName parses
//
//	qualifiedName : (qualifier+=symbolPrimitive PERIOD)* name=symbolPrimitive
//
// (PartiQLParser.g4:75) and returns the qualifiers and the name.
func (p *Parser) parseQualifiedName() ([]*ast.VarRef, *ast.VarRef, error) {
	var qualifier []*ast.VarRef
	for {
		name, caseSensitive, loc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, nil, err
		}
		ref := &ast.VarRef{Name: name, CaseSensitive: caseSensitive, Loc: loc}
		if p.cur.Type != tokPERIOD {
			return qualifier, ref, nil
		}
		p.advance() // consume .
		qualifier = append(qualifier, ref)
	}
}

// parseTableDef parses the column declarations of CREATE TABLE and returns
// them with the end of the closing parenthesis. p.cur is the opening
// parenthesis.
//
//	PAREN_LEFT tableDef PAREN_RIGHT
//	tableDef         : tableDefPart ( COMMA tableDefPart )*
//	tableDefPart     : columnName type columnConstraint*     # ColumnDeclaration
//	columnConstraint : ( CONSTRAINT columnConstraintName )? columnConstraintDef
//	columnConstraintDef : NOT NULL | NULL
//
// Grammar: PartiQLParser.g4:88, 97-112.
func (p *Parser) parseTableDef() ([]*ast.ColumnDef, int, error) {
	p.advance() // consume (
	var cols []*ast.ColumnDef
	for {
		col, err := p.parseColumnDef()
		if err != nil {
			return nil, 0, err
		}
		cols = append(cols, col)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, 0, err
	}
	return cols, rp.Loc.End, nil
}

// parseColumnDef parses `columnName type columnConstraint*`.
func (p *Parser) parseColumnDef() (*ast.ColumnDef, error) {
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	col := &ast.ColumnDef{
		Name: &ast.VarRef{Name: name, CaseSensitive: caseSensitive, Loc: nameLoc},
		Type: typ,
		Loc:  ast.Loc{Start: nameLoc.Start, End: typ.Loc.End},
	}
	for p.cur.Type == tokCONSTRAINT || p.cur.Type == tokNOT || p.cur.Type == tokNULL {
		c, err := p.parseColumnConstraint()
		if err != nil {
			return nil, err
		}
		col.Constraints = append(col.Constraints, c)
		col.Loc.End = c.Loc.End
	}
	return col, nil
}

// parseColumnConstraint parses
// `( CONSTRAINT columnConstraintName )? ( NOT NULL | NULL )`.
func (p *Parser) parseColumnConstraint() (*ast.ColumnConstraint, error) {
	c := &ast.ColumnConstraint{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	if p.cur.Type == tokCONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, caseSensitive, loc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		c.Name = &ast.VarRef{Name: name, CaseSensitive: caseSensitive, Loc: loc}
	}
	if p.cur.Type == tokNOT {
		p.advance() // consume NOT
		c.NotNull = true
	}
	end, err := p.expect(tokNULL)
	if err != nil {
		return nil, err
	}
	c.Loc.End = end.Loc.End
	return c, nil
}

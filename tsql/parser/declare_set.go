// Package parser - declare_set.go implements T-SQL DECLARE and SET statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseDeclareStmt parses a DECLARE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/declare-local-variable-transact-sql
//
//	DECLARE @var type [= expr], ...
//	DECLARE @var TABLE (col_def, ...)
func (p *Parser) parseDeclareStmt() *nodes.DeclareStmt {
	loc := p.pos()
	p.advance() // consume DECLARE

	stmt := &nodes.DeclareStmt{
		Loc: nodes.Loc{Start: loc},
	}

	var vars []nodes.Node
	for {
		vd := p.parseVariableDecl()
		if vd == nil {
			break
		}
		vars = append(vars, vd)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Variables = &nodes.List{Items: vars}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseVariableDecl parses a single variable declaration.
//
//	variable_decl = @name type [= expr]
//	             | @name TABLE (col_def, ...)
//	             | @name CURSOR
func (p *Parser) parseVariableDecl() *nodes.VariableDecl {
	if p.cur.Type != tokVARIABLE {
		return nil
	}

	loc := p.pos()
	vd := &nodes.VariableDecl{
		Name: p.cur.Str,
		Loc:  nodes.Loc{Start: loc},
	}
	p.advance() // consume @var

	// TABLE type
	if p.cur.Type == kwTABLE {
		p.advance()
		vd.IsTable = true
		if p.cur.Type == '(' {
			p.advance()
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				col := p.parseColumnDef()
				if col != nil {
					cols = append(cols, col)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			_, _ = p.expect(')')
			vd.TableDef = &nodes.List{Items: cols}
		}
		vd.Loc.End = p.pos()
		return vd
	}

	// CURSOR type
	if p.cur.Type == kwCURSOR {
		p.advance()
		vd.IsCursor = true
		vd.Loc.End = p.pos()
		return vd
	}

	// Data type
	vd.DataType = p.parseDataType()

	// Optional default value
	if p.cur.Type == '=' {
		p.advance()
		vd.Default = p.parseExpr()
	}

	vd.Loc.End = p.pos()
	return vd
}

// parseSetStmt parses a SET statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/set-local-variable-transact-sql
//
//	SET @var = expr
//	SET option ON|OFF
func (p *Parser) parseSetStmt() *nodes.SetStmt {
	loc := p.pos()
	p.advance() // consume SET

	stmt := &nodes.SetStmt{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == tokVARIABLE {
		// SET @var = expr
		stmt.Variable = p.cur.Str
		p.advance()
		if _, err := p.expect('='); err == nil {
			stmt.Value = p.parseExpr()
		}
	} else {
		// SET option ON|OFF (e.g., SET NOCOUNT ON, SET XACT_ABORT ON)
		if p.isIdentLike() || p.cur.Type == kwNOCOUNT || p.cur.Type == kwXACT_ABORT ||
			p.cur.Type == kwROWCOUNT || p.cur.Type == kwTEXTSIZE ||
			p.cur.Type == kwSTATISTICS || p.cur.Type == kwIDENTITY_INSERT {
			stmt.Variable = strings.ToUpper(p.cur.Str)
			p.advance()

			// Handle SET IDENTITY_INSERT table ON|OFF
			// For now, just consume ON/OFF
			if p.cur.Type == kwON {
				// Store "ON" as a literal
				onLoc := p.pos()
				p.advance()
				stmt.Value = &nodes.ColumnRef{Column: "ON", Loc: nodes.Loc{Start: onLoc}}
			} else if p.cur.Type == kwOFF {
				offLoc := p.pos()
				p.advance()
				stmt.Value = &nodes.ColumnRef{Column: "OFF", Loc: nodes.Loc{Start: offLoc}}
			} else {
				// Could be SET ROWCOUNT n or other SET options with values
				stmt.Value = p.parseExpr()
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// Package parser - create_type.go implements CREATE TYPE parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateTypeStmt parses a CREATE TYPE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-type-transact-sql
//
// BNF:
//
//	CREATE TYPE [ schema_name. ] type_name
//	{
//	    FROM base_type [ ( precision [ , scale ] ) ] [ NULL | NOT NULL ]
//	  | EXTERNAL NAME assembly_name [ .class_name ]
//	  | AS TABLE ( <column_definition> [ , ...n ] [ <table_constraint> ] [ , ...n ] )
//	}
func (p *Parser) parseCreateTypeStmt() *nodes.CreateTypeStmt {
	stmt := &nodes.CreateTypeStmt{}

	stmt.Name = p.parseTableRef()

	switch {
	case p.cur.Type == kwFROM:
		p.advance()
		stmt.BaseType = p.parseDataType()
		if p.cur.Type == kwNULL {
			b := true
			stmt.Nullable = &b
			p.advance()
		} else if p.cur.Type == kwNOT {
			next := p.peekNext()
			if next.Type == kwNULL {
				b := false
				stmt.Nullable = &b
				p.advance()
				p.advance()
			}
		}

	case p.cur.Type == kwEXTERNAL:
		p.advance()
		p.matchIdentCI("NAME")
		name, _ := p.parseIdentifier()
		extName := name
		if p.cur.Type == '.' {
			p.advance()
			className, _ := p.parseIdentifier()
			extName += "." + className
		}
		stmt.ExternalName = extName

	case p.cur.Type == kwAS:
		p.advance()
		p.match(kwTABLE)
		if _, err := p.expect('('); err == nil {
			var elements []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY ||
					p.cur.Type == kwUNIQUE || p.cur.Type == kwCHECK ||
					p.cur.Type == kwFOREIGN {
					constraint := p.parseTableConstraint()
					if constraint != nil {
						elements = append(elements, constraint)
					}
				} else if p.cur.Type == kwINDEX {
					p.skipToNextCommaOrParen()
				} else {
					col := p.parseColumnDef()
					if col != nil {
						elements = append(elements, col)
					}
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(elements) > 0 {
				stmt.TableDef = &nodes.List{Items: elements}
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// skipToNextCommaOrParen consumes tokens until a comma, closing paren, or EOF.
func (p *Parser) skipToNextCommaOrParen() {
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			if depth == 0 {
				return
			}
			depth--
		} else if p.cur.Type == ',' && depth == 0 {
			return
		}
		p.advance()
	}
}

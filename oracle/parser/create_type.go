package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateTypeStmt parses a CREATE [OR REPLACE] TYPE statement.
// The CREATE keyword has already been consumed. The caller has already parsed
// OR REPLACE if present and passes orReplace.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TYPE.html
//
//	CREATE [ OR REPLACE ] TYPE [ schema. ] type_name AS OBJECT (
//	    attribute_name datatype [, ...]
//	)
//	CREATE [ OR REPLACE ] TYPE [ schema. ] type_name AS TABLE OF datatype
//	CREATE [ OR REPLACE ] TYPE [ schema. ] type_name AS VARRAY ( n ) OF datatype
//	CREATE [ OR REPLACE ] TYPE BODY [ schema. ] type_name IS|AS ...
func (p *Parser) parseCreateTypeStmt(start int, orReplace bool) *nodes.CreateTypeStmt {
	stmt := &nodes.CreateTypeStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}

	// TYPE keyword
	if p.cur.Type == kwTYPE {
		p.advance()
	}

	// Check for TYPE BODY
	if p.cur.Type == kwBODY {
		stmt.IsBody = true
		p.advance()
	}

	// Type name
	stmt.Name = p.parseObjectName()

	// AS or IS
	if p.cur.Type == kwAS || p.cur.Type == kwIS {
		p.advance()
	}

	// Determine what kind of type:
	// - OBJECT ( ... )
	// - TABLE OF type
	// - VARRAY ( n ) OF type
	switch {
	case p.isIdentLikeStr("OBJECT"):
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			stmt.Attributes = p.parseTypeAttributeList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case p.cur.Type == kwTABLE:
		p.advance()
		if p.cur.Type == kwOF {
			p.advance()
		}
		stmt.AsTable = p.parseTypeName()

	case p.cur.Type == kwVARRAY || p.isIdentLikeStr("VARYING"):
		p.advance()
		// Handle VARYING ARRAY
		if p.isIdentLikeStr("ARRAY") {
			p.advance()
		}
		// ( size_limit )
		if p.cur.Type == '(' {
			p.advance()
			stmt.VarraySize = p.parseExpr()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		if p.cur.Type == kwOF {
			p.advance()
		}
		stmt.AsVarray = p.parseTypeName()

	default:
		// For TYPE BODY or other forms, skip the body for now.
		// If it's a body with IS/AS, we consume until we see a matching END.
		if stmt.IsBody {
			p.skipToEndBlock()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTypeAttributeList parses a comma-separated list of type attributes
// (attribute_name datatype).
func (p *Parser) parseTypeAttributeList() *nodes.List {
	list := &nodes.List{}
	for {
		if p.cur.Type == ')' || p.cur.Type == tokEOF {
			break
		}

		start := p.pos()
		name := p.parseIdentifier()
		if name == "" {
			break
		}

		typeName := p.parseTypeName()

		colDef := &nodes.ColumnDef{
			Name:     name,
			TypeName: typeName,
			Loc:      nodes.Loc{Start: start, End: p.pos()},
		}
		list.Items = append(list.Items, colDef)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list
}

// skipToEndBlock skips tokens until we find END; for TYPE BODY parsing.
// This is a placeholder for full PL/SQL body parsing.
func (p *Parser) skipToEndBlock() {
	depth := 1
	for p.cur.Type != tokEOF && depth > 0 {
		if p.cur.Type == kwBEGIN {
			depth++
		} else if p.cur.Type == kwEND {
			depth--
			if depth == 0 {
				p.advance() // consume END
				// consume optional type name after END
				if p.isIdentLike() {
					p.advance()
				}
				return
			}
		}
		p.advance()
	}
}

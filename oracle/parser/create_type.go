package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateTypeStmt parses a CREATE [OR REPLACE] TYPE statement.
// The CREATE keyword has already been consumed. The caller has already parsed
// OR REPLACE if present and passes orReplace.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TYPE.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-TYPE-BODY-statement.html
//
//	CREATE [ OR REPLACE ] TYPE [ schema. ] type_name AS OBJECT (
//	    attribute_name datatype [, ...]
//	)
//	CREATE [ OR REPLACE ] TYPE [ schema. ] type_name AS TABLE OF datatype
//	CREATE [ OR REPLACE ] TYPE [ schema. ] type_name AS VARRAY ( n ) OF datatype
//	CREATE [ OR REPLACE ] TYPE BODY [ schema. ] type_name { IS | AS }
//	  { { MEMBER | STATIC } { procedure_definition | function_definition }
//	    | MAP MEMBER function_definition
//	    | ORDER MEMBER function_definition
//	    | CONSTRUCTOR FUNCTION type_name
//	      [ ( [ SELF IN OUT type_name , ] parameter [, ...] ) ]
//	      RETURN SELF AS RESULT
//	      { IS | AS } { [ declare_section ] BEGIN statement... [EXCEPTION ...] END [name] ; }
//	  } ...
//	END [ type_name ] ;
func (p *Parser) parseCreateTypeStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) *nodes.CreateTypeStmt {
	stmt := &nodes.CreateTypeStmt{
		OrReplace:      orReplace,
		IfNotExists:    ifNotExists,
		Editionable:    editionable,
		NonEditionable: nonEditionable,
		Loc:            nodes.Loc{Start: start},
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
	// - TYPE BODY members
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
		// For TYPE BODY, parse structured members.
		if stmt.IsBody {
			stmt.Body = p.parseTypeBodyMembers()

			// END [type_name] ;
			if p.cur.Type == kwEND {
				p.advance()
			}
			// Optional type name after END
			if p.isIdentLike() && p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
			if p.cur.Type == ';' {
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTypeBodyMembers parses the member definitions inside a CREATE TYPE BODY.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-TYPE-BODY-statement.html
//
//	type_body_member:
//	  { MEMBER | STATIC } { procedure_definition | function_definition }
//	  | MAP MEMBER function_definition
//	  | ORDER MEMBER function_definition
//	  | CONSTRUCTOR FUNCTION type_name
//	    [ ( [ SELF IN OUT type_name , ] parameter [, ...] ) ]
//	    RETURN SELF AS RESULT
//	    { IS | AS } plsql_block
func (p *Parser) parseTypeBodyMembers() *nodes.List {
	members := &nodes.List{}

	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		// Skip standalone semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}

		member := p.parseTypeBodyMember()
		if member == nil {
			break
		}
		members.Items = append(members.Items, member)
	}

	return members
}

// parseTypeBodyMember parses a single type body member definition.
//
//	type_body_member:
//	  { MEMBER | STATIC } { PROCEDURE proc_name [(params)] IS|AS plsql_block
//	                       | FUNCTION func_name [(params)] RETURN type IS|AS plsql_block }
//	  | MAP MEMBER FUNCTION func_name [(params)] RETURN type IS|AS plsql_block
//	  | ORDER MEMBER FUNCTION func_name [(params)] RETURN type IS|AS plsql_block
//	  | CONSTRUCTOR FUNCTION type_name
//	    [ ( [ SELF IN OUT [NOCOPY] type_name , ] parameter [, ...] ) ]
//	    RETURN SELF AS RESULT IS|AS plsql_block
func (p *Parser) parseTypeBodyMember() *nodes.TypeBodyMember {
	start := p.pos()
	member := &nodes.TypeBodyMember{
		Loc: nodes.Loc{Start: start},
	}

	// Determine the kind prefix
	switch {
	case p.isIdentLikeStr("MEMBER"):
		member.Kind = nodes.TYPE_BODY_MEMBER
		p.advance() // consume MEMBER

	case p.isIdentLikeStr("STATIC"):
		member.Kind = nodes.TYPE_BODY_STATIC
		p.advance() // consume STATIC

	case p.isIdentLikeStr("MAP"):
		member.Kind = nodes.TYPE_BODY_MAP
		p.advance() // consume MAP
		// Expect MEMBER
		if p.isIdentLikeStr("MEMBER") {
			p.advance()
		}

	case p.cur.Type == kwORDER:
		member.Kind = nodes.TYPE_BODY_ORDER
		p.advance() // consume ORDER
		// Expect MEMBER
		if p.isIdentLikeStr("MEMBER") {
			p.advance()
		}

	case p.isIdentLikeStr("CONSTRUCTOR"):
		member.Kind = nodes.TYPE_BODY_CONSTRUCTOR
		p.advance() // consume CONSTRUCTOR

	default:
		return nil
	}

	// Parse the subprogram (PROCEDURE or FUNCTION)
	switch {
	case p.cur.Type == kwPROCEDURE:
		member.Subprog = p.parseTypeBodyProcedure()
	case p.cur.Type == kwFUNCTION:
		member.Subprog = p.parseTypeBodyFunction(member.Kind == nodes.TYPE_BODY_CONSTRUCTOR)
	default:
		return nil
	}

	member.Loc.End = p.pos()
	return member
}

// parseTypeBodyProcedure parses a PROCEDURE definition inside a type body.
//
//	PROCEDURE proc_name [ ( parameter [, ...] ) ]
//	  { IS | AS }
//	  [ declare_section ] BEGIN statements [ EXCEPTION handlers ] END [ name ] ;
func (p *Parser) parseTypeBodyProcedure() *nodes.CreateProcedureStmt {
	start := p.pos()
	p.advance() // consume PROCEDURE

	stmt := &nodes.CreateProcedureStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Optional parameter list
	if p.cur.Type == '(' {
		stmt.Parameters = p.parseParameterList()
	}

	// IS | AS
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}

	// PL/SQL block body
	stmt.Body = p.parsePLSQLBlock()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTypeBodyFunction parses a FUNCTION definition inside a type body.
// If isConstructor is true, it handles the RETURN SELF AS RESULT clause.
//
//	FUNCTION func_name [ ( parameter [, ...] ) ]
//	  RETURN datatype
//	  [ DETERMINISTIC ] [ PIPELINED ] [ PARALLEL_ENABLE ] [ RESULT_CACHE ]
//	  { IS | AS }
//	  [ declare_section ] BEGIN statements [ EXCEPTION handlers ] END [ name ] ;
//
//	constructor_function:
//	  FUNCTION type_name [ ( [ SELF IN OUT [NOCOPY] type_name , ] parameter [, ...] ) ]
//	  RETURN SELF AS RESULT
//	  { IS | AS }
//	  [ declare_section ] BEGIN statements [ EXCEPTION handlers ] END [ name ] ;
func (p *Parser) parseTypeBodyFunction(isConstructor bool) *nodes.CreateFunctionStmt {
	start := p.pos()
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Optional parameter list
	if p.cur.Type == '(' {
		stmt.Parameters = p.parseParameterList()
	}

	// RETURN type or RETURN SELF AS RESULT
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		if isConstructor && p.isIdentLikeStr("SELF") {
			// RETURN SELF AS RESULT
			p.advance() // consume SELF
			if p.cur.Type == kwAS {
				p.advance() // consume AS
			}
			if p.isIdentLikeStr("RESULT") {
				p.advance() // consume RESULT
			}
			// Set return type to indicate SELF AS RESULT
			stmt.ReturnType = &nodes.TypeName{
				Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELF AS RESULT"}}},
			}
		} else {
			stmt.ReturnType = p.parseTypeName()
		}
	}

	// Optional function properties
	p.parseFunctionProperties(stmt)

	// IS | AS
	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}

	// PL/SQL block body
	stmt.Body = p.parsePLSQLBlock()

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

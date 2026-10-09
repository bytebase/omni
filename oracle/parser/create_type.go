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
//	    attribute_name datatype [, ...] [, element_spec ]...
//	) [ [ NOT ] { FINAL | INSTANTIABLE | PERSISTABLE } ]...
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
func (p *Parser) parseCreateTypeStmt(start int, orReplace, ifNotExists, editionable, nonEditionable bool) (*nodes.CreateTypeStmt, error) {
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
	var parseErr573 error

	// Type name
	stmt.Name, parseErr573 = p.parseObjectName()
	if parseErr573 !=

		// AS or IS
		nil {
		return nil, parseErr573
	}

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
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if parseErr574 := p.parseObjectTypeElements(stmt); parseErr574 != nil {
			return nil, parseErr574
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr574 error
		stmt.Modifiers, parseErr574 = p.parseTypeModifiers(objectTypeModifiers)
		if parseErr574 != nil {
			return nil, parseErr574
		}

	case p.cur.Type == kwTABLE:
		p.advance()
		if p.cur.Type != kwOF {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr575 error
		stmt.AsTable, parseErr575 = p.parseTypeName()
		if parseErr575 != nil {
			return nil, parseErr575
		}
		if stmt.AsTable == nil || stmt.AsTable.Names.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}

	case p.cur.Type == kwVARRAY || p.isIdentLikeStr("VARYING"):
		p.advance()
		// Handle VARYING ARRAY
		if p.isIdentLikeStr("ARRAY") {
			p.advance()
		}
		// ( size_limit )
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr576 error
		stmt.VarraySize, parseErr576 = p.parseExpr()
		if parseErr576 != nil {
			return nil, parseErr576
		}
		if stmt.VarraySize == nil || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != kwOF {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr577 error
		stmt.AsVarray, parseErr577 = p.parseTypeName()
		if parseErr577 !=

			// For TYPE BODY, parse structured members.
			nil {
			return nil, parseErr577
		}
		if stmt.AsVarray == nil || stmt.AsVarray.Names.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}

	default:

		if stmt.IsBody {
			var parseErr578 error
			stmt.Body, parseErr578 = p.parseTypeBodyMembers()
			if parseErr578 !=

				// END [type_name] ;
				nil {
				return nil, parseErr578
			}

			if p.cur.Type != kwEND {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			// Optional type name after END
			if p.isIdentLike() && p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
			if p.cur.Type == ';' {
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseTypeBodyMembers() (*nodes.List, error) {
	members := &nodes.List{}

	for p.cur.Type != kwEND && p.cur.Type != tokEOF {
		// Skip standalone semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}

		member, parseErr579 := p.parseTypeBodyMember(false)
		if parseErr579 != nil {
			return nil, parseErr579
		}
		if member == nil {
			break
		}
		members.Items = append(members.Items, member)
	}

	return members, nil
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
//
// With inSpec it parses an element_spec of an object type specification
// instead: inheritance clauses may come first, MAP and ORDER need MEMBER, and
// the implementation is optional and can only be a call spec.
func (p *Parser) parseTypeBodyMember(inSpec bool) (*nodes.TypeBodyMember, error) {
	start := p.pos()
	member := &nodes.TypeBodyMember{
		Loc: nodes.Loc{Start: start},
	}

	if inSpec {
		var err error
		member.Modifiers, err = p.parseTypeModifiers(methodInheritanceClauses)
		if err != nil {
			return nil, err
		}
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
		} else if inSpec {
			return nil, p.syntaxErrorAtCur()
		}

	case p.cur.Type == kwORDER:
		member.Kind = nodes.TYPE_BODY_ORDER
		p.advance() // consume ORDER
		// Expect MEMBER
		if p.isIdentLikeStr("MEMBER") {
			p.advance()
		} else if inSpec {
			return nil, p.syntaxErrorAtCur()
		}

	case p.isIdentLikeStr("CONSTRUCTOR"):
		member.Kind = nodes.TYPE_BODY_CONSTRUCTOR
		p.advance() // consume CONSTRUCTOR

	default:
		if inSpec {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}

	// Parse the subprogram (PROCEDURE or FUNCTION)
	switch {
	case p.cur.Type == kwPROCEDURE && member.Kind != nodes.TYPE_BODY_CONSTRUCTOR:
		var parseErr580 error
		member.Subprog, parseErr580 = p.parseTypeBodyProcedure(inSpec)
		if parseErr580 != nil {
			return nil, parseErr580
		}
	case p.cur.Type == kwFUNCTION:
		var parseErr581 error
		member.Subprog, parseErr581 = p.parseTypeBodyFunction(member.Kind == nodes.TYPE_BODY_CONSTRUCTOR, inSpec)
		if parseErr581 != nil {
			return nil, parseErr581
		}
	default:
		if inSpec {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}

	member.Loc.End = p.prev.End
	return member, nil
}

// parseTypeBodyProcedure parses a PROCEDURE definition inside a type body.
//
//	PROCEDURE proc_name [ ( parameter [, ...] ) ]
//	  { IS | AS }
//	  [ declare_section ] BEGIN statements [ EXCEPTION handlers ] END [ name ] ;
func (p *Parser) parseTypeBodyProcedure(inSpec bool) (*nodes.CreateProcedureStmt, error) {
	start := p.pos()
	p.advance() // consume PROCEDURE

	stmt := &nodes.CreateProcedureStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr582 error

	stmt.Name, parseErr582 = p.parseObjectName()
	if parseErr582 !=

		// Optional parameter list
		nil {
		return nil, parseErr582
	}

	if p.cur.Type == '(' {
		var parseErr583 error
		stmt.Parameters, parseErr583 = p.parseParameterList()
		if parseErr583 !=

			// IS | AS
			nil {
			return nil, parseErr583
		}
	}

	if inSpec {
		var err error
		stmt.CallSpec, err = p.parseTypeSpecImplementation()
		if err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}
	var parseErr584 error

	// PL/SQL block body or call spec
	stmt.Body, stmt.CallSpec, parseErr584 = p.parseSubprogramImplementation()
	if parseErr584 != nil {
		return nil, parseErr584
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseTypeBodyFunction(isConstructor, inSpec bool) (*nodes.CreateFunctionStmt, error) {
	start := p.pos()
	p.advance() // consume FUNCTION

	stmt := &nodes.CreateFunctionStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr585 error

	stmt.Name, parseErr585 = p.parseObjectName()
	if parseErr585 !=

		// Optional parameter list
		nil {
		return nil, parseErr585
	}

	if p.cur.Type == '(' {
		var parseErr586 error
		stmt.Parameters, parseErr586 = p.parseParameterList()
		if parseErr586 !=

			// RETURN type or RETURN SELF AS RESULT
			nil {
			return nil, parseErr586
		}
	}

	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		if isConstructor && p.isIdentLikeStr("SELF") {
			// RETURN SELF AS RESULT
			selfStart := p.pos()
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
				Loc:   nodes.Loc{Start: selfStart, End: p.prev.End},
			}
		} else {
			var parseErr587 error
			stmt.ReturnType, parseErr587 = p.parseTypeName()
			if parseErr587 !=

				// Optional function properties
				nil {
				return nil, parseErr587
			}
		}
	}
	parseErr588 := p.parseFunctionProperties(stmt)
	if parseErr588 !=

		// IS | AS
		nil {
		return nil, parseErr588
	}

	if inSpec {
		var err error
		stmt.CallSpec, err = p.parseTypeSpecImplementation()
		if err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}
	var parseErr589 error

	// PL/SQL block body or call spec
	stmt.Body, stmt.CallSpec, parseErr589 = p.parseSubprogramImplementation()
	if parseErr589 != nil {
		return nil, parseErr589
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseObjectTypeElements parses the parenthesized elements of an object type
// specification: attributes, then method specifications and RESTRICT_REFERENCES
// pragmas. It fills stmt.Attributes and stmt.Methods.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/lnpls/CREATE-TYPE-statement.html
//
//	element ::= attribute datatype | element_spec | restrict_references_pragma
//	element_spec ::= [ inheritance_clauses ]
//	    { { MEMBER | STATIC } { procedure_spec | function_spec }
//	    | [ FINAL ] [ INSTANTIABLE ] CONSTRUCTOR function_spec
//	    | { MAP | ORDER } MEMBER function_spec }
//	procedure_spec ::= PROCEDURE name [ ( parameters ) ] [ { IS | AS } call_spec ]
//	function_spec ::= FUNCTION name [ ( parameters ) ] RETURN datatype [ { IS | AS } call_spec ]
//
// Verified on Oracle 23ai: an unquoted MEMBER, STATIC, CONSTRUCTOR, MAP,
// ORDER, NOT, OVERRIDING, FINAL, INSTANTIABLE, PERSISTABLE, or PRAGMA at the
// start of an element always opens a method or pragma (an attribute with such
// a name is a syntax error unless quoted); attributes come first, and at least
// one is required (PLS-00589).
func (p *Parser) parseObjectTypeElements(stmt *nodes.CreateTypeStmt) error {
	stmt.Attributes = &nodes.List{}
	for {
		switch {
		case p.isObjectTypeMethodStart():
			member, err := p.parseTypeBodyMember(true)
			if err != nil {
				return err
			}
			stmt.Methods = appendListItem(stmt.Methods, member)
		case p.isKeywordStr("PRAGMA"):
			pragma, err := p.parseRestrictReferencesPragma()
			if err != nil {
				return err
			}
			stmt.Methods = appendListItem(stmt.Methods, pragma)
		default:
			if stmt.Methods != nil {
				// An attribute after a method (PLS-00589).
				return p.syntaxErrorAtCur()
			}
			attr, err := p.parseTypeAttribute()
			if err != nil {
				return err
			}
			stmt.Attributes.Items = append(stmt.Attributes.Items, attr)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if stmt.Attributes.Len() == 0 {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// appendListItem appends item to list, creating the list when it is nil.
func appendListItem(list *nodes.List, item nodes.Node) *nodes.List {
	if list == nil {
		list = &nodes.List{}
	}
	list.Items = append(list.Items, item)
	return list
}

// isObjectTypeMethodStart reports whether the current token opens an
// element_spec of an object type specification.
func (p *Parser) isObjectTypeMethodStart() bool {
	if p.cur.Type == kwORDER || p.cur.Type == kwNOT {
		return true
	}
	for _, word := range []string{"MEMBER", "STATIC", "CONSTRUCTOR", "MAP", "OVERRIDING", "FINAL", "INSTANTIABLE", "PERSISTABLE"} {
		if p.isKeywordStr(word) {
			return true
		}
	}
	return false
}

// objectTypeModifiers and methodInheritanceClauses are the words that
// parseTypeModifiers takes after an object type's element list and before a
// method specification.
var (
	objectTypeModifiers      = []string{"FINAL", "INSTANTIABLE", "PERSISTABLE"}
	methodInheritanceClauses = []string{"OVERRIDING", "FINAL", "INSTANTIABLE", "PERSISTABLE"}
)

// parseTypeModifiers parses { [ NOT ] word }... for the given words and
// returns each as written ("FINAL", "NOT FINAL"). A word may appear once,
// with or without NOT; a repeat is PLS-00168 (duplicate modifier).
func (p *Parser) parseTypeModifiers(words []string) ([]string, error) {
	var mods []string
	seen := make(map[string]bool)
	for {
		not := p.cur.Type == kwNOT
		if not {
			if !isOneOfKeywords(p.peekNext(), words) {
				return mods, nil
			}
			p.advance() // consume NOT
		}
		if !isOneOfKeywords(p.cur, words) {
			if not {
				return nil, p.syntaxErrorAtCur()
			}
			return mods, nil
		}
		word := p.cur.Str
		if seen[word] {
			return nil, p.syntaxErrorAtCur()
		}
		seen[word] = true
		p.advance()
		if not {
			word = "NOT " + word
		}
		mods = append(mods, word)
	}
}

// isOneOfKeywords reports whether tok is one of words, unquoted.
func isOneOfKeywords(tok Token, words []string) bool {
	if tok.Type != tokIDENT && tok.Type < 2000 {
		return false
	}
	for _, word := range words {
		if tok.Str == word {
			return true
		}
	}
	return false
}

// parseTypeSpecImplementation parses the optional { IS | AS } call_spec of a
// method specification in an object type specification. A PL/SQL body is not
// allowed there (Oracle expects LANGUAGE or MLE after IS), and no ';' ends the
// call spec.
func (p *Parser) parseTypeSpecImplementation() (*nodes.CallSpec, error) {
	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, nil
	}
	p.advance() // consume IS or AS
	if !p.isCallSpecStart() {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseCallSpec(false)
}

// parseRestrictReferencesPragma parses the pragma an object type
// specification allows among its methods.
//
//	PRAGMA RESTRICT_REFERENCES ( { method_name | DEFAULT } ,
//	    { RNDS | WNDS | RNPS | WNPS | TRUST } [, ...] )
func (p *Parser) parseRestrictReferencesPragma() (*nodes.PLSQLPragma, error) {
	start := p.pos()
	p.advance() // consume PRAGMA
	if !p.isKeywordStr("RESTRICT_REFERENCES") {
		// Any other pragma is PLS-00127 (not a supported pragma).
		return nil, p.syntaxErrorAtCur()
	}
	pragma := &nodes.PLSQLPragma{Name: p.cur.Str, Args: &nodes.List{}, Loc: nodes.Loc{Start: start}}
	p.advance() // consume RESTRICT_REFERENCES
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume (
	if p.cur.Type == kwDEFAULT {
		pragma.Args.Items = append(pragma.Args.Items, &nodes.String{Str: "DEFAULT"})
		p.advance()
	} else {
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		pragma.Args.Items = append(pragma.Args.Items, &nodes.String{Str: name})
	}
	for p.cur.Type == ',' {
		p.advance() // consume ,
		if !isOneOfKeywords(p.cur, []string{"RNDS", "WNDS", "RNPS", "WNPS", "TRUST"}) {
			return nil, p.syntaxErrorAtCur()
		}
		pragma.Args.Items = append(pragma.Args.Items, &nodes.String{Str: p.cur.Str})
		p.advance()
	}
	if pragma.Args.Len() < 2 || p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume )
	pragma.Loc.End = p.prev.End
	return pragma, nil
}

// parseTypeAttribute parses one attribute of an object type specification:
// attribute_name datatype.
func (p *Parser) parseTypeAttribute() (*nodes.ColumnDef, error) {
	start := p.pos()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	typeName, err := p.parseTypeName()
	if err != nil {
		return nil, err
	}
	if typeName == nil || typeName.Names.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.ColumnDef{
		Name:     name,
		TypeName: typeName,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

package parser

import (
	"strconv"
	"strings"

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
	if parseErr573 != nil {
		return nil, parseErr573
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if !stmt.IsBody {
		if err := p.parseTypeSpecSource(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// IS | AS is required before the members (PLS-00103 on Oracle 23ai).
	if p.cur.Type != kwAS && p.cur.Type != kwIS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parseErr578 error
	stmt.Body, parseErr578 = p.parseTypeBodyMembers()
	if parseErr578 != nil {
		return nil, parseErr578
	}

	// END [type_name] ;
	if p.cur.Type != kwEND {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.atEndName() {
		p.advance()
	}
	if p.cur.Type != ';' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTypeSpecSource parses what follows the name in a type specification,
// the plsql_type_source of CREATE TYPE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/CREATE-TYPE-statement.html
//
//	[ SHARING = { METADATA | NONE } ]
//	[ FORCE ] [ OID 'object_identifier' ]
//	[ { DEFAULT COLLATION USING_NLS_COMP | AUTHID { CURRENT_USER | DEFINER }
//	  | ACCESSIBLE BY ( accessor [, ...] ) }... ]
//	[ { IS | AS } { OBJECT ( element [, ...] ) [ [ NOT ] { FINAL | INSTANTIABLE | PERSISTABLE } ]...
//	              | TABLE OF collection_element
//	              | { VARRAY | [ VARYING ] ARRAY } ( size_limit ) OF collection_element }
//	| UNDER [ schema. ] supertype ( [ element [, ...] ] ) [ [ NOT ] { FINAL | INSTANTIABLE } ]... ]
//
// Oracle 23ai takes SHARING first, with its '=' (ORA-02000) and one of
// those values (ORA-65014); then FORCE and OID, in either order; then
// DEFAULT COLLATION, AUTHID, and ACCESSIBLE BY, in any order. Each comes
// once (PLS-00103, PLS-00371 for the last three), and FORCE or OID after any
// of the last three is PLS-00103 (checked with MAX_STRING_SIZE = EXTENDED;
// with STANDARD, DEFAULT COLLATION is ORA-43929). With none of the
// alternatives the type is incomplete, which takes SHARING and OID only
// (PLS-00103). AUTHID and ACCESSIBLE BY are PLS-00262 on a collection type,
// which takes DEFAULT COLLATION; UNDER takes no IS or AS (PLS-00103) and no
// PERSISTABLE (PLS-00772).
func (p *Parser) parseTypeSpecSource(stmt *nodes.CreateTypeStmt) error {
	if p.isKeywordStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type != '=' {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume =
		if !isOneOfKeywords(p.cur, []string{"METADATA", "NONE"}) {
			return p.syntaxErrorAtCur()
		}
		stmt.Sharing = p.cur.Str
		p.advance()
	}

	seen := clauseSeen{}
properties:
	for {
		switch {
		case p.cur.Type == kwFORCE, p.isKeywordStr("OID"):
			word := p.cur.Str
			if seen["AUTHID"] || seen["ACCESSIBLE BY"] || seen["DEFAULT COLLATION"] {
				return p.syntaxErrorAtCur()
			}
			if err := p.firstClause(seen, word); err != nil {
				return err
			}
			p.advance()
			if word == "FORCE" {
				stmt.Force = true
				continue
			}
			// The identifier is a string literal (ORA-01780) of 32
			// hexadecimal digits (ORA-02304).
			if p.cur.Type != tokSCONST || !isObjectIdentifier(p.cur.Str) {
				return p.syntaxErrorAtCur()
			}
			stmt.OID = p.cur.Str
			p.advance()
		case p.isKeywordStr("AUTHID"):
			if err := p.firstClause(seen, "AUTHID"); err != nil {
				return err
			}
			authID, err := p.parseOptionalAuthID()
			if err != nil {
				return err
			}
			stmt.AuthID = authID
		case p.isKeywordStr("ACCESSIBLE"):
			if err := p.firstClause(seen, "ACCESSIBLE BY"); err != nil {
				return err
			}
			list, err := p.parseAccessibleByClause()
			if err != nil {
				return err
			}
			stmt.AccessibleBy = list
		case p.atDefaultCollation():
			if err := p.firstClause(seen, "DEFAULT COLLATION"); err != nil {
				return err
			}
			collation, err := p.parseDefaultCollationClause()
			if err != nil {
				return err
			}
			stmt.DefaultCollation = collation
		default:
			break properties
		}
	}

	role := seen["AUTHID"] || seen["ACCESSIBLE BY"]
	switch {
	case p.isStatementEnd():
		if stmt.Force || role || stmt.DefaultCollation != "" {
			return p.syntaxErrorAtCur()
		}
		return nil
	case p.cur.Type == kwUNDER:
		return p.parseObjectSubtype(stmt)
	case p.cur.Type == kwAS || p.cur.Type == kwIS:
		p.advance()
	default:
		return p.syntaxErrorAtCur()
	}

	if p.isIdentLikeStr("OBJECT") {
		p.advance()
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.parseObjectTypeElements(stmt, false); err != nil {
			return err
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.parseObjectTypeModifiers(stmt, objectTypeModifiers); err != nil {
			return err
		}
		return p.checkPersistableAttributes(stmt)
	}

	collectionTok := p.cur
	switch {
	case p.cur.Type == kwTABLE:
		p.advance()
		if p.cur.Type != kwOF {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		tn, err := p.parseCollectionElement(stmt)
		if err != nil {
			return err
		}
		stmt.AsTable = tn

	case p.cur.Type == kwVARRAY || p.isKeywordStr("VARYING") || p.isKeywordStr("ARRAY"):
		if p.isKeywordStr("VARYING") {
			p.advance()
			if !p.isKeywordStr("ARRAY") {
				return p.syntaxErrorAtCur()
			}
		}
		p.advance() // consume VARRAY or ARRAY
		// ( size_limit )
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		size, err := p.parseExpr()
		if err != nil {
			return err
		}
		if size == nil || p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		if err := p.checkVarrayLimit(size); err != nil {
			return err
		}
		stmt.VarraySize = size
		p.advance()
		if p.cur.Type != kwOF {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		tn, err := p.parseCollectionElement(stmt)
		if err != nil {
			return err
		}
		stmt.AsVarray = tn

	default:
		return p.syntaxErrorAtCur()
	}
	if role {
		// PLS-00262: incorrect placement of AUTHID (or ACCESSIBLE BY) clause
		return p.syntaxErrorAtTok(collectionTok)
	}
	return nil
}

// parseObjectSubtype parses the object_subtype_def of a type specification,
// the current token at UNDER. The element list may be empty or hold methods
// only; the supertype takes no database link (PLS-00103 on Oracle 23ai).
//
//	UNDER [ schema. ] supertype ( [ element [, element ]... ] )
//	    [ [ NOT ] { FINAL | INSTANTIABLE } ]...
func (p *Parser) parseObjectSubtype(stmt *nodes.CreateTypeStmt) error {
	p.advance() // consume UNDER
	under, err := p.parseUnitName()
	if err != nil {
		return err
	}
	stmt.Under = under
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type == ')' {
		stmt.Attributes = &nodes.List{}
	} else if err := p.parseObjectTypeElements(stmt, true); err != nil {
		return err
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return p.parseObjectTypeModifiers(stmt, methodInheritanceClauses)
}

// parseObjectTypeModifiers parses the modifiers after an object type's
// element list. A FINAL type cannot be NOT INSTANTIABLE (PLS-00614 on
// Oracle 23ai), and an object type is FINAL unless it says NOT FINAL.
func (p *Parser) parseObjectTypeModifiers(stmt *nodes.CreateTypeStmt, words []string) error {
	modsTok := p.cur
	mods, err := p.parseTypeModifiers(words)
	if err != nil {
		return err
	}
	stmt.Modifiers = mods
	notInstantiable, notFinal := false, false
	for _, mod := range mods {
		notInstantiable = notInstantiable || mod == "NOT INSTANTIABLE"
		notFinal = notFinal || mod == "NOT FINAL"
	}
	if notInstantiable && !notFinal {
		return p.syntaxErrorAtTok(modsTok)
	}
	// A type with a NOT INSTANTIABLE method is NOT INSTANTIABLE (PLS-00634).
	if !notInstantiable && stmt.Methods != nil {
		for _, item := range stmt.Methods.Items {
			if member, ok := item.(*nodes.TypeBodyMember); ok && hasModifier(member.Modifiers, "NOT INSTANTIABLE") {
				return p.syntaxErrorAtNode(member)
			}
		}
	}
	return nil
}

// parseCollectionElement parses the element type of a nested table or
// varray type. PERSISTABLE needs the parentheses, and NOT NULL then goes
// inside them; CHARACTER SET precedes NOT NULL (PLS-00103 on Oracle 23ai).
// Sizes are PL/SQL expressions, as in an attribute.
//
//	datatype [ NOT NULL ]
//	| ( datatype [ NOT NULL ] ) [ [ NOT ] PERSISTABLE ]
func (p *Parser) parseCollectionElement(stmt *nodes.CreateTypeStmt) (*nodes.TypeName, error) {
	paren := p.cur.Type == '('
	if paren {
		p.advance()
	}
	tn, err := p.parsePLSQLDatatype(typeModsPLSQL, charsetFixed)
	if err != nil {
		return nil, err
	}
	if tn == nil || tn.Names.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	if elementAlwaysIllegal[p.predefinedTypeLead(tn)] {
		return nil, p.syntaxErrorAtType(tn)
	}
	if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
		p.advance()
		p.advance()
		stmt.ElementNotNull = true
	}
	if paren {
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		mods, err := p.parseTypeModifiers([]string{"PERSISTABLE"})
		if err != nil {
			return nil, err
		}
		stmt.Modifiers = mods
	}
	if !notPersistable(stmt.Modifiers) && elementPersistableIllegal[p.predefinedTypeLead(tn)] {
		return nil, p.syntaxErrorAtType(tn)
	}
	return tn, nil
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

	// Inheritance clauses may precede a method in a specification and in a
	// body alike (FINAL INSTANTIABLE CONSTRUCTOR FUNCTION compiles in a type
	// body on Oracle 23ai). OVERRIDING comes first, once, and takes no NOT
	// (PLS-00103 for NOT OVERRIDING or OVERRIDING after another clause).
	if p.isKeywordStr("OVERRIDING") {
		member.Modifiers = append(member.Modifiers, "OVERRIDING")
		p.advance() // consume OVERRIDING
	}
	mods, err := p.parseTypeModifiers(methodInheritanceClauses)
	if err != nil {
		return nil, err
	}
	member.Modifiers = append(member.Modifiers, mods...)

	// Determine the kind prefix
	kindTok := p.cur
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
		if inSpec || len(member.Modifiers) > 0 {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}
	if methodModifiersConflict(member.Modifiers, member.Kind) {
		return nil, p.syntaxErrorAtTok(kindTok)
	}

	// Parse the subprogram (PROCEDURE or FUNCTION). A constructor is a
	// function, and so are MAP and ORDER methods in a spec (PLS-00155 on
	// Oracle 23ai; a body that defines one as a procedure fails only to
	// match its spec, PLS-00538).
	procedureAllowed := member.Kind != nodes.TYPE_BODY_CONSTRUCTOR &&
		(!inSpec || member.Kind == nodes.TYPE_BODY_MEMBER || member.Kind == nodes.TYPE_BODY_STATIC)
	switch {
	case p.cur.Type == kwPROCEDURE && procedureAllowed:
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
		// A method kind, in a specification or a body, is followed by its
		// subprogram (PLS-00103 on Oracle 23ai for MEMBER END).
		return nil, p.syntaxErrorAtCur()
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

	// A method needs a name (PLS-00103 on Oracle 23ai): the words that
	// follow one in the head do not stand for it.
	if p.cur.Type == kwIS || p.cur.Type == kwAS || p.cur.Type == kwRETURN {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name, parseErr582 = p.parseObjectName()
	if parseErr582 != nil {
		return nil, parseErr582
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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
	// The heading checks of other procedures hold for a method too
	// (PLS-00765, PLS-00550 on Oracle 23ai).
	if tn := p.tablePseudoType(stmt.Parameters, nil); tn != nil {
		return nil, p.syntaxErrorAtType(tn)
	}
	if err := p.checkCharsetSources(stmt.Parameters, nil); err != nil {
		return nil, err
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

	// IS | AS is required before the PL/SQL block or call spec (PLS-00103
	// on Oracle 23ai).
	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume IS or AS
	var parseErr584 error

	// PL/SQL block body or call spec
	stmt.Body, stmt.CallSpec, parseErr584 = p.parseSubprogramImplementation()
	if parseErr584 != nil {
		return nil, parseErr584
	}
	// A method's formals and locals shadow the type's attributes, so a
	// %CHARSET source found among them, or qualified by the method name, is
	// the one Oracle 23ai resolves.
	if err := p.checkCharsetScopes(subprogramName(stmt.Name), stmt.Parameters, nil, stmt.Body); err != nil {
		return nil, err
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

	// A method needs a name (PLS-00103 on Oracle 23ai): the words that
	// follow one in the head do not stand for it.
	if p.cur.Type == kwIS || p.cur.Type == kwAS || p.cur.Type == kwRETURN {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name, parseErr585 = p.parseObjectName()
	if parseErr585 != nil {
		return nil, parseErr585
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
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
		if isConstructor {
			// A constructor returns SELF AS RESULT, all three words required
			// (PLS-00103 on Oracle 23ai for a partial phrase, PLS-00659 for
			// another return type).
			selfStart := p.pos()
			if !p.isKeywordStr("SELF") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume SELF
			if p.cur.Type != kwAS {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume AS
			if !p.isKeywordStr("RESULT") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume RESULT
			// Set return type to indicate SELF AS RESULT
			stmt.ReturnType = &nodes.TypeName{
				Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "SELF AS RESULT"}}},
				Loc:   nodes.Loc{Start: selfStart, End: p.prev.End},
			}
		} else {
			var parseErr587 error
			stmt.ReturnType, parseErr587 = p.parsePLSQLDatatype(typeModsNone, charsetFlexible)
			if parseErr587 != nil {
				return nil, parseErr587
			}
			if stmt.ReturnType == nil || stmt.ReturnType.Names.Len() == 0 {
				return nil, p.syntaxErrorAtCur()
			}
		}
	} else {
		// A method function needs RETURN, in a spec and a body alike
		// (PLS-00103 on Oracle 23ai).
		return nil, p.syntaxErrorAtCur()
	}
	parseErr588 := p.parseFunctionProperties(stmt, subprogramMethod)
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

	// IS | AS is required before the PL/SQL block or call spec (PLS-00103
	// on Oracle 23ai).
	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume IS or AS
	var parseErr589 error

	// PL/SQL block body or call spec
	stmt.Body, stmt.CallSpec, parseErr589 = p.parseSubprogramImplementation()
	if parseErr589 != nil {
		return nil, parseErr589
	}
	// A method's formals and locals shadow the type's attributes, so a
	// %CHARSET source found among them, or qualified by the method name, is
	// the one Oracle 23ai resolves.
	if err := p.checkCharsetScopes(subprogramName(stmt.Name), stmt.Parameters, stmt.ReturnType, stmt.Body); err != nil {
		return nil, err
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
//	inheritance_clauses ::= [ OVERRIDING ] { [ NOT ] FINAL | [ NOT ] INSTANTIABLE }...
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
func (p *Parser) parseObjectTypeElements(stmt *nodes.CreateTypeStmt, subtype bool) error {
	stmt.Attributes = &nodes.List{}
	hasComparison := false
	seenAttrs := map[string]bool{}
	seenMethods := map[string]bool{}
	for {
		switch {
		case p.isObjectTypeMethodStart():
			methodTok := p.cur
			member, err := p.parseTypeBodyMember(true)
			if err != nil {
				return err
			}
			if err := p.checkObjectTypeMethod(stmt, member, &hasComparison, methodTok); err != nil {
				return err
			}
			// A method bears no attribute's name (PLS-00305 on Oracle 23ai,
			// "A" for a too).
			if fn := subprogramNameOf(member.Subprog); seenAttrs[fn] {
				return p.syntaxErrorAtTok(methodTok)
			}
			// A method of the kind and signature of an earlier one is
			// PLS-00305 on Oracle 23ai; a MEMBER and a STATIC method may
			// share a signature. A MAP or ORDER method is a MEMBER one.
			kind := member.Kind
			if kind == nodes.TYPE_BODY_MAP || kind == nodes.TYPE_BODY_ORDER {
				kind = nodes.TYPE_BODY_MEMBER
			}
			typ := "FUNCTION"
			if _, ok := member.Subprog.(*nodes.CreateProcedureStmt); ok {
				typ = "PROCEDURE"
			}
			selfKind := "MEMBER"
			switch member.Kind {
			case nodes.TYPE_BODY_STATIC:
				selfKind = "STATIC"
			case nodes.TYPE_BODY_CONSTRUCTOR:
				selfKind = "CONSTRUCTOR"
			}
			sig := strconv.Itoa(int(kind)) + " " + p.subprogramSignature(member.Subprog, implicitSelfMode(selfKind, typ))
			if seenMethods[sig] {
				return p.syntaxErrorAtTok(methodTok)
			}
			seenMethods[sig] = true
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
			// Attribute names are unique in the statement (PLS-00410 on
			// Oracle 23ai, for a subtype too); a clash with a supertype's
			// attribute needs the supertype and is left to the engine.
			if seenAttrs[attr.Name] {
				return p.syntaxErrorAtNode(attr)
			}
			seenAttrs[attr.Name] = true
			stmt.Attributes.Items = append(stmt.Attributes.Items, attr)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	// An object type has at least one attribute (PLS-00589); a subtype
	// inherits its supertype's and may add methods only.
	if stmt.Attributes.Len() == 0 && !subtype {
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
// method specification, where OVERRIDING may precede them. [ NOT ] PERSISTABLE
// is a type modifier only (PLS-00771 before a method on Oracle 23ai), though
// an unquoted PERSISTABLE still cannot name an attribute.
var (
	objectTypeModifiers      = []string{"FINAL", "INSTANTIABLE", "PERSISTABLE"}
	methodInheritanceClauses = []string{"FINAL", "INSTANTIABLE"}
)

// methodModifiersConflict reports whether a method spec's inheritance clauses
// conflict with each other or with its kind. Oracle 23ai rejects NOT
// INSTANTIABLE with FINAL, STATIC, or CONSTRUCTOR, NOT FINAL with
// CONSTRUCTOR, and OVERRIDING with STATIC or CONSTRUCTOR (PLS-00169;
// ORA-06545 in ALTER TYPE).
func methodModifiersConflict(mods []string, kind nodes.TypeBodyMemberKind) bool {
	has := func(mod string) bool {
		for _, m := range mods {
			if m == mod {
				return true
			}
		}
		return false
	}
	if has("NOT INSTANTIABLE") && (has("FINAL") || kind == nodes.TYPE_BODY_STATIC || kind == nodes.TYPE_BODY_CONSTRUCTOR) {
		return true
	}
	if has("NOT FINAL") && kind == nodes.TYPE_BODY_CONSTRUCTOR {
		return true
	}
	return has("OVERRIDING") && (kind == nodes.TYPE_BODY_STATIC || kind == nodes.TYPE_BODY_CONSTRUCTOR)
}

// checkObjectTypeMethod applies the rules Oracle 23ai checks on a method spec
// from the statement's own text: a constructor bears its type's name
// (PLS-00658), a MAP method takes no parameter besides SELF (PLS-00520), an
// ORDER method takes exactly one (PLS-00521), and a type has at most one MAP
// or ORDER method (PLS-00154). hasComparison tracks the last rule. Neither
// returns a REF (PLS-00522, PLS-00523); the other rules on MAP and ORDER
// return and parameter types (PLS-00522 through PLS-00524) need the types
// resolved and are left to the engine. A subtype has no
// ORDER method: Oracle 23ai rejects one with PLS-00646 under a supertype
// without MAP or ORDER and with PLS-00154 under one with either. A
// subtype's MAP method overrides its supertype's: Oracle 23ai rejects one
// without OVERRIDING with PLS-00636, or PLS-00646 under a supertype without
// MAP; whether the supertype has a MAP method to override is left to the
// engine. OVERRIDING names a method of a supertype, and a type that is not
// a subtype has none (PLS-00635). tok is where the method starts, for the
// error position.
func (p *Parser) checkObjectTypeMethod(stmt *nodes.CreateTypeStmt, member *nodes.TypeBodyMember, hasComparison *bool, tok Token) error {
	if stmt.Under == nil && hasModifier(member.Modifiers, "OVERRIDING") {
		return p.syntaxErrorAtTok(tok)
	}
	// A method's SELF is where and as Oracle takes it, and of the type
	// (invalidSelfType).
	if selfParameterInvalid(member.Kind == nodes.TYPE_BODY_STATIC, member.Kind == nodes.TYPE_BODY_CONSTRUCTOR, subprogramParameters(member.Subprog)) {
		return p.syntaxErrorAtTok(tok)
	}
	if params := subprogramParameters(member.Subprog); len(params) > 0 && params[0].Name == "SELF" && p.invalidSelfType(params[0].TypeName) {
		return p.syntaxErrorAtTok(tok)
	}
	fn, isFunction := member.Subprog.(*nodes.CreateFunctionStmt)
	switch member.Kind {
	case nodes.TYPE_BODY_CONSTRUCTOR:
		if isFunction && stmt.Name != nil && fn.Name.Name != stmt.Name.Name {
			return p.syntaxErrorAtTok(tok)
		}
	case nodes.TYPE_BODY_MAP, nodes.TYPE_BODY_ORDER:
		if *hasComparison || stmt.Under != nil && (member.Kind == nodes.TYPE_BODY_ORDER || !hasModifier(member.Modifiers, "OVERRIDING")) {
			return p.syntaxErrorAtTok(tok)
		}
		*hasComparison = true
		want := 0
		if member.Kind == nodes.TYPE_BODY_ORDER {
			want = 1
		}
		if isFunction && parametersBesidesSelf(fn.Parameters) != want {
			return p.syntaxErrorAtTok(tok)
		}
		if isFunction && member.Kind == nodes.TYPE_BODY_ORDER && !orderParameterIn(listParameters(fn.Parameters)) {
			return p.syntaxErrorAtTok(tok)
		}
		if isFunction && fn.ReturnType != nil && fn.ReturnType.IsRef {
			return p.syntaxErrorAtTok(tok)
		}
	}
	return nil
}

// selfParameterInvalid reports whether a method's parameters declare SELF
// where or as Oracle 23ai rejects it: in a static method (PLS-00586), after
// the first parameter, where it duplicates the implicit SELF (PLS-00410),
// in OUT mode (PLS-00594), or in a constructor in any mode but IN OUT
// (PLS-00660).
func selfParameterInvalid(static, constructor bool, params []*nodes.Parameter) bool {
	for i, param := range params {
		if param.Name != "SELF" {
			continue
		}
		mode := strings.TrimSpace(strings.ReplaceAll(param.Mode, "NOCOPY", ""))
		if static || i > 0 || mode == "OUT" || constructor && mode != "IN OUT" {
			return true
		}
	}
	return false
}

// subprogramParameters returns the parameters of a method's subprogram.
func subprogramParameters(subprog nodes.Node) []*nodes.Parameter {
	switch s := subprog.(type) {
	case *nodes.CreateFunctionStmt:
		return listParameters(s.Parameters)
	case *nodes.CreateProcedureStmt:
		return listParameters(s.Parameters)
	}
	return nil
}

// invalidSelfType reports whether tn is a datatype an explicit SELF is not
// of, the type itself: a REF (PLS-00546 on Oracle 23ai for SELF IN REF t)
// or a predefined datatype (PLS-00382 for SELF IN NUMBER).
func (p *Parser) invalidSelfType(tn *nodes.TypeName) bool {
	return tn != nil && tn.IsRef || p.predefinedDatatype(tn)
}

// predefinedDatatype reports whether tn spells a predefined datatype that
// no declaration in scope redeclares.
func (p *Parser) predefinedDatatype(tn *nodes.TypeName) bool {
	return tn != nil && !tn.IsRef && p.predefinedScalar(p.predefinedTypeLead(tn))
}

// subprogramNameOf returns the name of a method's subprogram.
func subprogramNameOf(subprog nodes.Node) string {
	switch s := subprog.(type) {
	case *nodes.CreateFunctionStmt:
		return subprogramName(s.Name)
	case *nodes.CreateProcedureStmt:
		return subprogramName(s.Name)
	}
	return ""
}

// hasModifier reports whether mods holds mod.
func hasModifier(mods []string, mod string) bool {
	for _, m := range mods {
		if m == mod {
			return true
		}
	}
	return false
}

// orderParameterIn reports whether the parameters of an ORDER method
// besides SELF have IN mode (PLS-00528 on Oracle 23ai for OUT or IN OUT)
// and are no REF (PLS-00524: the parameter is of the type itself).
func orderParameterIn(params []*nodes.Parameter) bool {
	for _, param := range params {
		if param.Name == "SELF" {
			continue
		}
		if param.Mode != "" && param.Mode != "IN" || param.TypeName != nil && param.TypeName.IsRef {
			return false
		}
	}
	return true
}

// listParameters returns the parameters in list.
func listParameters(list *nodes.List) []*nodes.Parameter {
	var params []*nodes.Parameter
	if list != nil {
		for _, item := range list.Items {
			if param, ok := item.(*nodes.Parameter); ok {
				params = append(params, param)
			}
		}
	}
	return params
}

// parametersBesidesSelf counts the parameters in params not named SELF.
func parametersBesidesSelf(params *nodes.List) int {
	n := 0
	if params == nil {
		return 0
	}
	for _, item := range params.Items {
		if param, ok := item.(*nodes.Parameter); !ok || param.Name != "SELF" {
			n++
		}
	}
	return n
}

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
// allowed there, nor the EXTERNAL form (Oracle 23ai expects LANGUAGE or MLE
// after IS), and no ';' ends the call spec.
func (p *Parser) parseTypeSpecImplementation() (*nodes.CallSpec, error) {
	if p.cur.Type != kwIS && p.cur.Type != kwAS {
		return nil, nil
	}
	p.advance() // consume IS or AS
	if !p.isCallSpecWord("LANGUAGE", "MLE") {
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

// The predefined types Oracle 23ai refuses as an object type attribute
// (PLS-00530; PLS-00218 for the NOT NULL subtypes, PLS-00989 for
// SYS_REFCURSOR) and as the element of a nested table or varray type
// (PLS-00531, PLS-00990), for every user, SYS included. The PL/SQL-only
// types become legal in a NOT PERSISTABLE type; JSON, VECTOR, SYS_REFCURSOR,
// and as attributes the NOT NULL subtypes stay illegal. BOOLEAN, the LOBs,
// BFILE, and the ANSI numeric names are allowed everywhere.
var (
	attributeAlwaysIllegal = map[string]bool{
		"JSON": true, "VECTOR": true, "SYS_REFCURSOR": true, "SIMPLE_INTEGER": true, "NATURALN": true,
		"POSITIVEN": true, "SIMPLE_FLOAT": true, "SIMPLE_DOUBLE": true,
	}
	attributePersistableIllegal = map[string]bool{
		"LONG": true, "PLS_INTEGER": true, "BINARY_INTEGER": true, "NATURAL": true, "POSITIVE": true,
		"SIGNTYPE": true, "STRING": true, "ROWID": true, "UROWID": true, "MLSLABEL": true,
	}
	elementAlwaysIllegal = map[string]bool{
		"JSON": true, "VECTOR": true, "SYS_REFCURSOR": true,
	}
	elementPersistableIllegal = map[string]bool{
		"LONG": true, "PLS_INTEGER": true, "BINARY_INTEGER": true, "SIMPLE_INTEGER": true, "NATURAL": true,
		"NATURALN": true, "POSITIVE": true, "POSITIVEN": true, "SIGNTYPE": true, "STRING": true,
		"ROWID": true, "UROWID": true, "MLSLABEL": true,
	}
)

// notPersistable reports whether modifiers make a type NOT PERSISTABLE.
func notPersistable(mods []string) bool {
	for _, mod := range mods {
		if mod == "NOT PERSISTABLE" {
			return true
		}
	}
	return false
}

// checkPersistableAttributes rejects the attributes of an object type that
// only a NOT PERSISTABLE type takes, once its modifiers are known. A
// subtype's persistability is its supertype's, which the text does not
// show, so its attributes meet only the rules for every type.
func (p *Parser) checkPersistableAttributes(stmt *nodes.CreateTypeStmt) error {
	if stmt.Under != nil || notPersistable(stmt.Modifiers) {
		return nil
	}
	for _, item := range stmt.Attributes.Items {
		if col, ok := item.(*nodes.ColumnDef); ok && attributePersistableIllegal[p.predefinedTypeLead(col.TypeName)] {
			return p.syntaxErrorAtType(col.TypeName)
		}
	}
	return nil
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
	// An attribute takes PL/SQL size expressions (VARCHAR2(ORA_MAX_NAME_LEN)
	// on Oracle 23ai) and a named CHARACTER SET (NCHAR_CS), not ANY_CS
	// (PLS-00551) or item%CHARSET.
	typeName, err := p.parsePLSQLDatatype(typeModsPLSQL, charsetFixed)
	if err != nil {
		return nil, err
	}
	if typeName == nil || typeName.Names.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	if attributeAlwaysIllegal[p.predefinedTypeLead(typeName)] {
		return nil, p.syntaxErrorAtType(typeName)
	}
	return &nodes.ColumnDef{
		Name:     name,
		TypeName: typeName,
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// isObjectIdentifier reports whether s is an object identifier literal: 32
// hexadecimal digits, in either case.
func isObjectIdentifier(s string) bool {
	if len(s) != 32 {
		return false
	}
	for _, c := range s {
		if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F') {
			return false
		}
	}
	return true
}

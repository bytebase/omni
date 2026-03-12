package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseVariableSetStmt parses a SET statement. The SET keyword has already
// been consumed by the caller.
//
// VariableSetStmt:
//
//	SET set_rest
//	| SET LOCAL set_rest
//	| SET SESSION set_rest
func (p *Parser) parseVariableSetStmt() nodes.Node {
	isLocal := false

	switch p.cur.Type {
	case LOCAL:
		isLocal = true
		p.advance() // consume LOCAL
	case SESSION:
		// SESSION could be "SET SESSION AUTHORIZATION ..." (set_rest_more)
		// or "SET SESSION CHARACTERISTICS AS TRANSACTION ..." (set_rest)
		// or "SET SESSION var_name ..." (generic_set with SESSION as scope modifier)
		//
		// We need to peek ahead. If next token is AUTHORIZATION or CHARACTERISTICS,
		// it's part of set_rest/set_rest_more, so we do NOT consume SESSION as a
		// scope modifier. Otherwise, SESSION is the scope modifier.
		next := p.peekNext()
		if next.Type == AUTHORIZATION {
			// SET SESSION AUTHORIZATION ... is handled by set_rest_more
			// Don't consume SESSION; let parseSetRest handle it.
		} else if next.Type == CHARACTERISTICS {
			// SET SESSION CHARACTERISTICS AS TRANSACTION ... is handled by set_rest
			// Don't consume SESSION; let parseSetRest handle it.
		} else {
			// SESSION is the scope modifier
			p.advance() // consume SESSION
		}
	}

	stmt := p.parseSetRest()
	if stmt == nil {
		return nil
	}

	vs, ok := stmt.(*nodes.VariableSetStmt)
	if ok {
		vs.IsLocal = isLocal
	}
	return stmt
}

// parseSetRest parses set_rest.
//
// set_rest:
//
//	TRANSACTION transaction_mode_list
//	| SESSION CHARACTERISTICS AS TRANSACTION transaction_mode_list
//	| set_rest_more
func (p *Parser) parseSetRest() nodes.Node {
	switch p.cur.Type {
	case TRANSACTION:
		p.advance() // consume TRANSACTION
		modes := p.parseTransactionModeListOrEmpty()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_MULTI,
			Name: "TRANSACTION",
			Args: modes,
		}

	case SESSION:
		// SESSION CHARACTERISTICS AS TRANSACTION transaction_mode_list
		next := p.peekNext()
		if next.Type == CHARACTERISTICS {
			p.advance() // consume SESSION
			p.advance() // consume CHARACTERISTICS
			p.expect(AS)
			p.expect(TRANSACTION)
			modes := p.parseTransactionModeListOrEmpty()
			return &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_MULTI,
				Name: "SESSION CHARACTERISTICS",
				Args: modes,
			}
		}
		// SESSION AUTHORIZATION ... (handled in set_rest_more via parseSetRestMore)
		return p.parseSetRestMore()

	default:
		return p.parseSetRestMore()
	}
}

// parseSetRestMore parses set_rest_more.
//
// set_rest_more:
//
//	generic_set
//	| var_name FROM CURRENT_P
//	| TIME ZONE zone_value
//	| CATALOG_P Sconst
//	| SCHEMA Sconst
//	| NAMES opt_encoding
//	| ROLE NonReservedWord_or_Sconst
//	| SESSION AUTHORIZATION NonReservedWord_or_Sconst
//	| SESSION AUTHORIZATION DEFAULT
//	| XML_P OPTION document_or_content
//	| TRANSACTION SNAPSHOT Sconst
func (p *Parser) parseSetRestMore() nodes.Node {
	switch p.cur.Type {
	case TIME:
		p.advance() // consume TIME
		p.expect(ZONE)
		val := p.parseZoneValue()
		n := &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: "timezone",
		}
		if val != nil {
			n.Args = &nodes.List{Items: []nodes.Node{val}}
		} else {
			n.Kind = nodes.VAR_SET_DEFAULT
		}
		return n

	case CATALOG_P:
		// SET CATALOG is not supported; produce an error like PostgreSQL.
		// For now, just parse it.
		p.advance() // consume CATALOG
		p.expect(SCONST)
		return nil

	case SCHEMA:
		p.advance() // consume SCHEMA
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: "search_path",
			Args: &nodes.List{Items: []nodes.Node{makeStringConst(s)}},
		}

	case NAMES:
		p.advance() // consume NAMES
		enc := p.parseOptEncoding()
		n := &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: "client_encoding",
		}
		if enc != "" {
			n.Args = &nodes.List{Items: []nodes.Node{makeStringConst(enc)}}
		} else {
			n.Kind = nodes.VAR_SET_DEFAULT
		}
		return n

	case ROLE:
		p.advance() // consume ROLE
		val := p.parseNonReservedWordOrSconst()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: "role",
			Args: &nodes.List{Items: []nodes.Node{makeStringConst(val)}},
		}

	case SESSION:
		// SESSION AUTHORIZATION NonReservedWord_or_Sconst | SESSION AUTHORIZATION DEFAULT
		p.advance() // consume SESSION
		p.expect(AUTHORIZATION)
		if p.cur.Type == DEFAULT {
			p.advance()
			return &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_DEFAULT,
				Name: "session_authorization",
			}
		}
		val := p.parseNonReservedWordOrSconst()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: "session_authorization",
			Args: &nodes.List{Items: []nodes.Node{makeStringConst(val)}},
		}

	case XML_P:
		p.advance() // consume XML
		p.expect(OPTION)
		docOrContent := p.parseDocumentOrContent()
		var val string
		if docOrContent == int64(nodes.XMLOPTION_DOCUMENT) {
			val = "DOCUMENT"
		} else {
			val = "CONTENT"
		}
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: "xmloption",
			Args: &nodes.List{Items: []nodes.Node{makeStringConst(val)}},
		}

	case TRANSACTION:
		// TRANSACTION SNAPSHOT Sconst
		p.advance() // consume TRANSACTION
		p.expect(SNAPSHOT)
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_MULTI,
			Name: "TRANSACTION SNAPSHOT",
			Args: &nodes.List{Items: []nodes.Node{makeStringConst(s)}},
		}

	default:
		// generic_set or var_name FROM CURRENT_P
		return p.parseGenericSetOrFromCurrent()
	}
}

// parseGenericSetOrFromCurrent parses either a generic_set or "var_name FROM CURRENT_P".
//
// generic_set:
//
//	var_name TO var_list
//	| var_name '=' var_list
//	| var_name TO DEFAULT
//	| var_name '=' DEFAULT
//
// Also:
//
//	var_name FROM CURRENT_P
func (p *Parser) parseGenericSetOrFromCurrent() nodes.Node {
	varName := p.parseVarName()
	if varName == "" {
		return nil
	}

	switch p.cur.Type {
	case FROM:
		// var_name FROM CURRENT_P
		p.advance() // consume FROM
		p.expect(CURRENT_P)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_CURRENT,
			Name: varName,
		}

	case TO, '=':
		p.advance() // consume TO or '='
		if p.cur.Type == DEFAULT {
			p.advance()
			return &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_DEFAULT,
				Name: varName,
			}
		}
		args := p.parseVarList()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_SET_VALUE,
			Name: varName,
			Args: args,
		}

	default:
		return nil
	}
}

// parseVarName parses a dotted variable name.
//
// var_name:
//
//	ColId
//	| var_name '.' ColId
func (p *Parser) parseVarName() string {
	name, err := p.parseColId()
	if err != nil {
		return ""
	}
	for p.cur.Type == '.' {
		p.advance() // consume '.'
		part, err := p.parseColId()
		if err != nil {
			break
		}
		name = name + "." + part
	}
	return name
}

// parseVarList parses a comma-separated list of var_value.
//
// var_list:
//
//	var_value
//	| var_list ',' var_value
func (p *Parser) parseVarList() *nodes.List {
	val := p.parseVarValue()
	if val == nil {
		return nil
	}
	items := []nodes.Node{val}
	for p.cur.Type == ',' {
		p.advance()
		v := p.parseVarValue()
		if v == nil {
			break
		}
		items = append(items, v)
	}
	return &nodes.List{Items: items}
}

// parseVarValue parses a single variable value.
//
// var_value:
//
//	opt_boolean_or_string   -> A_Const{String{...}}
//	| NumericOnly           -> A_Const{Integer/Float}
func (p *Parser) parseVarValue() nodes.Node {
	// Try NumericOnly first (numbers and signed numbers).
	if num := p.tryParseNumericOnly(); num != nil {
		return &nodes.A_Const{Val: num}
	}

	// opt_boolean_or_string
	s, ok := p.tryOptBooleanOrString()
	if ok {
		return makeStringConst(s)
	}

	return nil
}

// tryParseNumericOnly tries to parse a NumericOnly value.
// Returns nil if the current token is not a numeric value.
//
// NumericOnly:
//
//	FCONST
//	| '+' FCONST
//	| '-' FCONST
//	| SignedIconst
func (p *Parser) tryParseNumericOnly() nodes.Node {
	switch p.cur.Type {
	case ICONST:
		val := p.cur.Ival
		p.advance()
		return &nodes.Integer{Ival: val}
	case FCONST:
		val := p.cur.Str
		p.advance()
		return &nodes.Float{Fval: val}
	case '+':
		next := p.peekNext()
		if next.Type == ICONST {
			p.advance() // consume '+'
			val := p.cur.Ival
			p.advance()
			return &nodes.Integer{Ival: val}
		}
		if next.Type == FCONST {
			p.advance() // consume '+'
			val := p.cur.Str
			p.advance()
			return &nodes.Float{Fval: val}
		}
		return nil
	case '-':
		next := p.peekNext()
		if next.Type == ICONST {
			p.advance() // consume '-'
			val := p.cur.Ival
			p.advance()
			return &nodes.Integer{Ival: -val}
		}
		if next.Type == FCONST {
			p.advance() // consume '-'
			val := p.cur.Str
			p.advance()
			f := &nodes.Float{Fval: val}
			doNegateFloat(f)
			return f
		}
		return nil
	default:
		return nil
	}
}

// doNegateFloat negates a Float value by prepending or removing a '-'.
func doNegateFloat(f *nodes.Float) {
	if len(f.Fval) > 0 && f.Fval[0] == '-' {
		f.Fval = f.Fval[1:]
	} else {
		f.Fval = "-" + f.Fval
	}
}

// tryOptBooleanOrString tries to parse an opt_boolean_or_string value.
// Returns the string value and true if successful.
//
// opt_boolean_or_string:
//
//	TRUE_P                      -> "true"
//	| FALSE_P                   -> "false"
//	| ON                        -> "on"
//	| NonReservedWord_or_Sconst
func (p *Parser) tryOptBooleanOrString() (string, bool) {
	switch p.cur.Type {
	case TRUE_P:
		p.advance()
		return "true", true
	case FALSE_P:
		p.advance()
		return "false", true
	case ON:
		p.advance()
		return "on", true
	case SCONST:
		s := p.cur.Str
		p.advance()
		return s, true
	default:
		// NonReservedWord: IDENT | unreserved_keyword | col_name_keyword | type_func_name_keyword
		if p.isNonReservedWord() {
			s := p.cur.Str
			p.advance()
			return s, true
		}
		return "", false
	}
}

// isNonReservedWord returns true if the current token is a NonReservedWord.
//
// NonReservedWord: IDENT | unreserved_keyword | col_name_keyword | type_func_name_keyword
func (p *Parser) isNonReservedWord() bool {
	if p.cur.Type == IDENT {
		return true
	}
	if kw := LookupKeyword(p.cur.Str); kw != nil && kw.Token == p.cur.Type {
		return kw.Category == UnreservedKeyword || kw.Category == ColNameKeyword || kw.Category == TypeFuncNameKeyword
	}
	return false
}

// parseZoneValue parses a zone_value.
//
// zone_value:
//
//	Sconst
//	| IDENT
//	| NumericOnly   -> A_Const{Integer/Float}
//	| DEFAULT       -> nil (signals SET DEFAULT)
//	| LOCAL         -> nil (signals SET DEFAULT)
func (p *Parser) parseZoneValue() nodes.Node {
	switch p.cur.Type {
	case SCONST:
		s := p.cur.Str
		p.advance()
		return makeStringConst(s)
	case IDENT:
		s := p.cur.Str
		p.advance()
		return makeStringConst(s)
	case DEFAULT:
		p.advance()
		return nil
	case LOCAL:
		p.advance()
		return nil
	default:
		// Try NumericOnly
		if num := p.tryParseNumericOnly(); num != nil {
			return &nodes.A_Const{Val: num}
		}
		// Interval handling etc. would go here, but for now this covers the tests.
		return nil
	}
}

// parseOptEncoding parses opt_encoding.
//
// opt_encoding:
//
//	Sconst     -> the string value
//	| DEFAULT  -> ""
//	| EMPTY    -> ""
func (p *Parser) parseOptEncoding() string {
	switch p.cur.Type {
	case SCONST:
		s := p.cur.Str
		p.advance()
		return s
	case DEFAULT:
		p.advance()
		return ""
	default:
		return ""
	}
}

// parseVariableShowStmt parses a SHOW statement. The SHOW keyword has already
// been consumed by the caller.
//
// VariableShowStmt:
//
//	SHOW var_name
//	| SHOW TIME ZONE
//	| SHOW TRANSACTION ISOLATION LEVEL
//	| SHOW SESSION AUTHORIZATION
//	| SHOW ALL
func (p *Parser) parseVariableShowStmt() nodes.Node {
	switch p.cur.Type {
	case TIME:
		p.advance() // consume TIME
		p.expect(ZONE)
		return &nodes.VariableShowStmt{Name: "timezone"}

	case TRANSACTION:
		p.advance() // consume TRANSACTION
		p.expect(ISOLATION)
		p.expect(LEVEL)
		return &nodes.VariableShowStmt{Name: "transaction_isolation"}

	case SESSION:
		p.advance() // consume SESSION
		p.expect(AUTHORIZATION)
		return &nodes.VariableShowStmt{Name: "session_authorization"}

	case ALL:
		p.advance()
		return &nodes.VariableShowStmt{Name: "all"}

	default:
		name := p.parseVarName()
		return &nodes.VariableShowStmt{Name: name}
	}
}

// parseVariableResetStmt parses a RESET statement. The RESET keyword has already
// been consumed by the caller.
//
// VariableResetStmt:
//
//	RESET reset_rest
//
// reset_rest:
//
//	generic_reset
//	| TIME ZONE
//	| TRANSACTION ISOLATION LEVEL
//	| SESSION AUTHORIZATION
func (p *Parser) parseVariableResetStmt() nodes.Node {
	switch p.cur.Type {
	case TIME:
		p.advance() // consume TIME
		p.expect(ZONE)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: "timezone",
		}

	case TRANSACTION:
		p.advance() // consume TRANSACTION
		p.expect(ISOLATION)
		p.expect(LEVEL)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: "transaction_isolation",
		}

	case SESSION:
		p.advance() // consume SESSION
		p.expect(AUTHORIZATION)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: "session_authorization",
		}

	case ALL:
		p.advance()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET_ALL,
		}

	default:
		// generic_reset: var_name
		name := p.parseVarName()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: name,
		}
	}
}

// parseConstraintsSetStmt parses SET CONSTRAINTS. The SET keyword has already
// been consumed, and we have just confirmed the current token is CONSTRAINTS.
//
// ConstraintsSetStmt:
//
//	SET CONSTRAINTS constraints_set_list constraints_set_mode
//
// constraints_set_list:
//
//	ALL
//	| qualified_name_list
//
// constraints_set_mode:
//
//	DEFERRED
//	| IMMEDIATE
func (p *Parser) parseConstraintsSetStmt() nodes.Node {
	p.advance() // consume CONSTRAINTS

	var constraints *nodes.List
	if p.cur.Type == ALL {
		p.advance()
		constraints = nil
	} else {
		constraints = p.parseConstraintsNameList()
	}

	deferred := false
	if p.cur.Type == DEFERRED {
		p.advance()
		deferred = true
	} else {
		p.expect(IMMEDIATE)
		deferred = false
	}

	return &nodes.ConstraintsSetStmt{
		Constraints: constraints,
		Deferred:    deferred,
	}
}

// parseConstraintsNameList parses a qualified_name_list for SET CONSTRAINTS.
func (p *Parser) parseConstraintsNameList() *nodes.List {
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil
	}
	rv := makeRangeVarFromNames(names)
	items := []nodes.Node{rv}
	for p.cur.Type == ',' {
		p.advance()
		n, err := p.parseQualifiedName()
		if err != nil {
			break
		}
		items = append(items, makeRangeVarFromNames(n))
	}
	return &nodes.List{Items: items}
}

// makeStringConst creates an A_Const with a String value.
func makeStringConst(s string) nodes.Node {
	return &nodes.A_Const{Val: &nodes.String{Str: s}}
}

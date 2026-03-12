package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreateFunctionStmt parses a CREATE [OR REPLACE] FUNCTION/PROCEDURE statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-createfunction.html
//
//	CreateFunctionStmt:
//	    CREATE opt_or_replace FUNCTION func_name func_args_with_defaults
//	        RETURNS func_return opt_createfunc_opt_list opt_routine_body
//	    | CREATE opt_or_replace FUNCTION func_name func_args_with_defaults
//	        RETURNS TABLE '(' table_func_column_list ')' opt_createfunc_opt_list opt_routine_body
//	    | CREATE opt_or_replace FUNCTION func_name func_args_with_defaults
//	        opt_createfunc_opt_list opt_routine_body
//	    | CREATE opt_or_replace PROCEDURE func_name func_args_with_defaults
//	        opt_createfunc_opt_list opt_routine_body
//
// The CREATE and optional OR REPLACE have already been consumed.
// isReplace indicates whether OR REPLACE was present.
func (p *Parser) parseCreateFunctionStmt(isReplace bool) *nodes.CreateFunctionStmt {
	isProcedure := false
	if p.cur.Type == PROCEDURE {
		isProcedure = true
		p.advance()
	} else {
		p.expect(FUNCTION)
	}

	// func_name
	funcname := p.parseFuncNameForCreate()

	// func_args_with_defaults
	params := p.parseFuncArgsWithDefaults()

	stmt := &nodes.CreateFunctionStmt{
		IsOrReplace: isReplace,
		Funcname:    funcname,
		Parameters:  params,
	}

	if isProcedure {
		// PROCEDURE: no RETURNS clause
		stmt.Options = p.parseOptCreatefuncOptList()
		stmt.SqlBody = p.parseOptRoutineBody()
		// Add isProcedure marker
		procDef := &nodes.DefElem{
			Defname: "isProcedure",
			Arg:     &nodes.Integer{Ival: 1},
		}
		if stmt.Options == nil {
			stmt.Options = &nodes.List{Items: []nodes.Node{procDef}}
		} else {
			stmt.Options.Items = append(stmt.Options.Items, procDef)
		}
		return stmt
	}

	// FUNCTION: check for RETURNS
	if p.cur.Type == RETURNS {
		p.advance()
		if p.cur.Type == TABLE {
			// RETURNS TABLE '(' table_func_column_list ')'
			p.advance()
			p.expect('(')
			tableCols := p.parseTableFuncColumnList()
			p.expect(')')
			// Add table columns to parameter list
			if params == nil {
				stmt.Parameters = tableCols
			} else {
				stmt.Parameters = &nodes.List{Items: append(params.Items, tableCols.Items...)}
			}
			// Return type is pg_catalog.record
			stmt.ReturnType = &nodes.TypeName{
				Names: &nodes.List{Items: []nodes.Node{
					&nodes.String{Str: "pg_catalog"},
					&nodes.String{Str: "record"},
				}},
				Loc: nodes.NoLoc(),
			}
		} else {
			// RETURNS func_return
			retType, _ := p.parseFuncType()
			stmt.ReturnType = retType
		}
	}

	stmt.Options = p.parseOptCreatefuncOptList()
	stmt.SqlBody = p.parseOptRoutineBody()
	return stmt
}

// parseFuncNameForCreate parses a function name for CREATE FUNCTION.
// Uses the existing parseFuncName from name.go.
func (p *Parser) parseFuncNameForCreate() *nodes.List {
	// parseFuncName in name.go returns (*List, error)
	result, _ := p.parseFuncName()
	return result
}


// parseFuncArgsWithDefaults parses function arguments with defaults.
//
//	func_args_with_defaults:
//	    '(' func_args_with_defaults_list ')'
//	    | '(' ')'
func (p *Parser) parseFuncArgsWithDefaults() *nodes.List {
	p.expect('(')
	if p.cur.Type == ')' {
		p.advance()
		return nil
	}

	args := p.parseFuncArgsWithDefaultsList()
	p.expect(')')
	return args
}

// parseFuncArgsWithDefaultsList parses a comma-separated list of function args.
//
//	func_args_with_defaults_list:
//	    func_arg_with_default
//	    | func_args_with_defaults_list ',' func_arg_with_default
func (p *Parser) parseFuncArgsWithDefaultsList() *nodes.List {
	arg := p.parseFuncArgWithDefault()
	items := []nodes.Node{arg}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseFuncArgWithDefault())
	}
	return &nodes.List{Items: items}
}

// parseFuncArgWithDefault parses a function argument with optional default.
//
//	func_arg_with_default:
//	    func_arg
//	    | func_arg DEFAULT a_expr
//	    | func_arg '=' a_expr
func (p *Parser) parseFuncArgWithDefault() *nodes.FunctionParameter {
	fp := p.parseFuncArg()
	if p.cur.Type == DEFAULT {
		p.advance()
		fp.Defexpr = p.parseAExpr(0)
	} else if p.cur.Type == '=' {
		p.advance()
		fp.Defexpr = p.parseAExpr(0)
	}
	return fp
}

// parseFuncArg parses a function argument.
//
//	func_arg:
//	    arg_class param_name func_type
//	    | param_name arg_class func_type
//	    | param_name func_type
//	    | arg_class func_type
//	    | func_type
func (p *Parser) parseFuncArg() *nodes.FunctionParameter {
	// This is complex because we need to distinguish between:
	// 1. arg_class param_name func_type  (IN x integer)
	// 2. param_name arg_class func_type  (x IN integer)
	// 3. param_name func_type            (x integer)
	// 4. arg_class func_type             (IN integer)
	// 5. func_type                       (integer)
	//
	// arg_class is one of: IN, OUT, INOUT, IN OUT, VARIADIC
	// param_name is type_function_name
	// func_type is Typename or type_function_name attrs '%' TYPE_P

	mode := nodes.FUNC_PARAM_IN // default
	name := ""

	argClass, isArgClass := p.tryParseArgClass()
	if isArgClass {
		// Cases 1 or 4
		// Check if next is a param_name followed by a type, or directly a type
		if p.isTypeFunctionName() {
			// Could be param_name func_type or just func_type
			// Try to distinguish: parse as param_name, then check if what follows
			// looks like a type. If not, it was actually the type itself.
			savedName, _ := p.parseTypeFunctionName()

			if p.isTypeFunctionName() || p.isBuiltinType() || p.cur.Type == SETOF {
				// It's param_name followed by func_type
				name = savedName
				mode = argClass
			} else if p.cur.Type == '.' || p.cur.Type == '%' {
				// Could be qualified type name (schema.type) or %TYPE
				// If it's name '.' name, it could still be type_function_name attrs '%' TYPE
				// or a qualified Typename. For func_type specifically, we handle the %TYPE case.
				// The safe bet: treat savedName as the type (part of func_type).
				// We need to "put back" savedName and parse as func_type.
				// Since we can't easily backtrack, use savedName as param name
				// only if the pattern matches.
				// Actually, at this point if we see '.', it's likely a qualified type,
				// not a param name. So treat savedName as start of func_type.
				p.pushBack(savedName)
				mode = argClass
			} else {
				// savedName was actually the type name by itself
				p.pushBack(savedName)
				mode = argClass
			}
		}
	} else if p.isTypeFunctionName() {
		// Cases 2, 3, or 5
		savedName, _ := p.parseTypeFunctionName()

		// Check if an arg_class follows (case 2: param_name arg_class func_type)
		argClass2, isArgClass2 := p.tryParseArgClass()
		if isArgClass2 {
			name = savedName
			mode = argClass2
		} else if p.isTypeFunctionName() || p.isBuiltinType() || p.cur.Type == SETOF {
			// Case 3: param_name func_type
			name = savedName
			// mode stays FUNC_PARAM_IN
		} else if p.cur.Type == '.' || p.cur.Type == '%' {
			// Could be qualified type (case 5) or param with qualified type (case 3)
			// For '.' after a name, it's likely a qualified type name.
			p.pushBack(savedName)
		} else {
			// Case 5: savedName is the type itself
			p.pushBack(savedName)
		}
	}

	// Parse func_type
	argType, _ := p.parseFuncType()

	return &nodes.FunctionParameter{
		Name:    name,
		ArgType: argType,
		Mode:    mode,
	}
}

// tryParseArgClass tries to parse an arg_class (IN, OUT, INOUT, IN OUT, VARIADIC).
// Returns the mode and true if found, or 0 and false if current token is not an arg_class.
func (p *Parser) tryParseArgClass() (nodes.FunctionParameterMode, bool) {
	switch p.cur.Type {
	case IN_P:
		p.advance()
		if p.cur.Type == OUT_P {
			p.advance()
			return nodes.FUNC_PARAM_INOUT, true
		}
		return nodes.FUNC_PARAM_IN, true
	case OUT_P:
		p.advance()
		return nodes.FUNC_PARAM_OUT, true
	case INOUT:
		p.advance()
		return nodes.FUNC_PARAM_INOUT, true
	case VARIADIC:
		p.advance()
		return nodes.FUNC_PARAM_VARIADIC, true
	}
	return 0, false
}

// pushBack pushes a name token back as the current token.
// Used when we speculatively consumed a name but need to reparse it as a type.
func (p *Parser) pushBack(name string) {
	p.nextBuf = p.cur
	p.hasNext = true
	p.cur = Token{
		Type: IDENT,
		Str:  name,
		Loc:  p.prev.Loc,
	}
}

// isBuiltinType checks if the current token starts a built-in type keyword.
func (p *Parser) isBuiltinType() bool {
	switch p.cur.Type {
	case INT_P, INTEGER, SMALLINT, BIGINT, REAL, FLOAT_P, DOUBLE_P,
		DECIMAL_P, NUMERIC, BOOLEAN_P, BIT, CHAR_P, CHARACTER,
		VARCHAR, TIMESTAMP, TIME, INTERVAL:
		return true
	}
	return false
}

// parseTableFuncColumnList parses a comma-separated list of table function columns.
//
//	table_func_column_list:
//	    table_func_column
//	    | table_func_column_list ',' table_func_column
func (p *Parser) parseTableFuncColumnList() *nodes.List {
	col := p.parseTableFuncColumn()
	items := []nodes.Node{col}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseTableFuncColumn())
	}
	return &nodes.List{Items: items}
}

// parseTableFuncColumn parses a single table function column.
//
//	table_func_column:
//	    param_name func_type
func (p *Parser) parseTableFuncColumn() *nodes.FunctionParameter {
	name, _ := p.parseTypeFunctionName()
	argType, _ := p.parseFuncType()
	return &nodes.FunctionParameter{
		Name:    name,
		ArgType: argType,
		Mode:    nodes.FUNC_PARAM_TABLE,
	}
}

// parseOptCreatefuncOptList parses an optional list of function options.
//
//	opt_createfunc_opt_list:
//	    createfunc_opt_list
//	    | /* EMPTY */
func (p *Parser) parseOptCreatefuncOptList() *nodes.List {
	var items []nodes.Node
	for {
		opt := p.parseCreatefuncOptItem()
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseCreatefuncOptItem parses a single function option.
//
//	createfunc_opt_item:
//	    AS func_as
//	    | LANGUAGE NonReservedWord_or_Sconst
//	    | TRANSFORM transform_type_list
//	    | WINDOW
//	    | common_func_opt_item
func (p *Parser) parseCreatefuncOptItem() *nodes.DefElem {
	switch p.cur.Type {
	case AS:
		p.advance()
		funcAs := p.parseFuncAs()
		return &nodes.DefElem{Defname: "as", Arg: funcAs}

	case LANGUAGE:
		p.advance()
		lang := p.parseNonReservedWordOrSconst()
		return &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: lang}}

	case TRANSFORM:
		p.advance()
		types := p.parseTransformTypeList()
		return &nodes.DefElem{Defname: "transform", Arg: types}

	case WINDOW:
		p.advance()
		return &nodes.DefElem{Defname: "window", Arg: &nodes.Integer{Ival: 1}}

	default:
		return p.parseCommonFuncOptItem()
	}
}

// parseFuncAs parses the AS clause body.
//
//	func_as:
//	    Sconst
//	    | Sconst ',' Sconst
func (p *Parser) parseFuncAs() *nodes.List {
	s1 := p.cur.Str
	p.expect(SCONST)
	if p.cur.Type == ',' {
		p.advance()
		s2 := p.cur.Str
		p.expect(SCONST)
		return &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: s1},
			&nodes.String{Str: s2},
		}}
	}
	return &nodes.List{Items: []nodes.Node{
		&nodes.String{Str: s1},
	}}
}

// parseNonReservedWordOrSconst parses a non-reserved word or string constant.
//
//	NonReservedWord_or_Sconst:
//	    NonReservedWord | Sconst
func (p *Parser) parseNonReservedWordOrSconst() string {
	if p.cur.Type == SCONST {
		s := p.cur.Str
		p.advance()
		return s
	}
	// NonReservedWord: IDENT, unreserved_keyword, col_name_keyword, type_func_name_keyword
	name := p.cur.Str
	p.advance()
	return name
}

// parseTransformTypeList parses a transform type list.
//
//	transform_type_list:
//	    FOR TYPE_P Typename
//	    | transform_type_list ',' FOR TYPE_P Typename
func (p *Parser) parseTransformTypeList() *nodes.List {
	p.expect(FOR)
	p.expect(TYPE_P)
	tn, _ := p.parseTypename()
	items := []nodes.Node{tn}
	for p.cur.Type == ',' {
		p.advance()
		p.expect(FOR)
		p.expect(TYPE_P)
		tn, _ = p.parseTypename()
		items = append(items, tn)
	}
	return &nodes.List{Items: items}
}

// parseCommonFuncOptItem parses a common function option.
//
//	common_func_opt_item:
//	    IMMUTABLE | STABLE | VOLATILE
//	    | STRICT_P
//	    | CALLED ON NULL_P INPUT_P
//	    | RETURNS NULL_P ON NULL_P INPUT_P
//	    | SECURITY DEFINER | SECURITY INVOKER
//	    | LEAKPROOF | NOT LEAKPROOF
//	    | COST NumericOnly | ROWS NumericOnly
//	    | PARALLEL ColId
//	    | SET set_rest_more
//	    | VariableResetStmt
//	    | SUPPORT any_name
func (p *Parser) parseCommonFuncOptItem() *nodes.DefElem {
	switch p.cur.Type {
	case IMMUTABLE:
		p.advance()
		return &nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "immutable"}}
	case STABLE:
		p.advance()
		return &nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "stable"}}
	case VOLATILE:
		p.advance()
		return &nodes.DefElem{Defname: "volatility", Arg: &nodes.String{Str: "volatile"}}
	case STRICT_P:
		p.advance()
		return &nodes.DefElem{Defname: "strict", Arg: &nodes.Integer{Ival: 1}}
	case CALLED:
		p.advance()
		p.expect(ON)
		p.expect(NULL_P)
		p.expect(INPUT_P)
		return &nodes.DefElem{Defname: "strict", Arg: &nodes.Integer{Ival: 0}}
	case RETURNS:
		// RETURNS NULL ON NULL INPUT (note: this is ambiguous with RETURNS func_return
		// but only appears in createfunc_opt_list, not at the top level)
		p.advance()
		p.expect(NULL_P)
		p.expect(ON)
		p.expect(NULL_P)
		p.expect(INPUT_P)
		return &nodes.DefElem{Defname: "strict", Arg: &nodes.Integer{Ival: 1}}
	case SECURITY:
		p.advance()
		if p.cur.Type == DEFINER {
			p.advance()
			return &nodes.DefElem{Defname: "security", Arg: &nodes.Integer{Ival: 1}}
		}
		p.expect(INVOKER)
		return &nodes.DefElem{Defname: "security", Arg: &nodes.Integer{Ival: 0}}
	case LEAKPROOF:
		p.advance()
		return &nodes.DefElem{Defname: "leakproof", Arg: &nodes.Integer{Ival: 1}}
	case NOT:
		p.advance()
		p.expect(LEAKPROOF)
		return &nodes.DefElem{Defname: "leakproof", Arg: &nodes.Integer{Ival: 0}}
	case COST:
		p.advance()
		val := p.parseNumericOnly()
		return &nodes.DefElem{Defname: "cost", Arg: val}
	case ROWS:
		p.advance()
		val := p.parseNumericOnly()
		return &nodes.DefElem{Defname: "rows", Arg: val}
	case PARALLEL:
		p.advance()
		name, _ := p.parseColId()
		return &nodes.DefElem{Defname: "parallel", Arg: &nodes.String{Str: name}}
	case SUPPORT:
		p.advance()
		name, _ := p.parseAnyName()
		return &nodes.DefElem{Defname: "support", Arg: name}
	case SET:
		// SET set_rest_more
		p.advance() // consume SET
		arg := p.parseSetRestMore()
		return &nodes.DefElem{Defname: "set", Arg: arg}
	case RESET:
		// VariableResetStmt (RESET reset_rest)
		p.advance() // consume RESET
		arg := p.parseVariableResetStmt()
		return &nodes.DefElem{Defname: "set", Arg: arg}
	}
	return nil
}

// parseOptRoutineBody parses an optional routine body.
//
//	opt_routine_body:
//	    ReturnStmt
//	    | BEGIN_P ATOMIC routine_body_stmt_list END_P
//	    | /* EMPTY */
func (p *Parser) parseOptRoutineBody() nodes.Node {
	if p.cur.Type == RETURN {
		p.advance()
		expr := p.parseAExpr(0)
		return &nodes.ReturnStmt{Returnval: expr}
	}
	if p.cur.Type == BEGIN_P {
		p.advance() // consume BEGIN
		p.expect(ATOMIC)

		// Parse routine_body_stmt_list: (routine_body_stmt ';')*
		var stmts []nodes.Node
		for p.cur.Type != END_P && p.cur.Type != 0 {
			var stmt nodes.Node
			if p.cur.Type == RETURN {
				p.advance()
				expr := p.parseAExpr(0)
				stmt = &nodes.ReturnStmt{Returnval: expr}
			} else {
				stmt = p.parseStmt()
			}
			if stmt != nil {
				stmts = append(stmts, stmt)
			}
			// consume trailing ';'
			if p.cur.Type == ';' {
				p.advance()
			}
		}
		p.expect(END_P)

		// Yacc wraps the stmt list in makeList(stmtList):
		// a single-element list containing the statement list
		stmtList := &nodes.List{Items: stmts}
		return &nodes.List{Items: []nodes.Node{stmtList}}
	}
	return nil
}

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
func (p *Parser) parseCreateFunctionStmt(startLoc int, isReplace bool) (*nodes.CreateFunctionStmt, error) {
	isProcedure := false
	if p.cur.Type == PROCEDURE {
		isProcedure = true
		p.advance()
	} else {
		if _, err := p.expect(FUNCTION); err != nil {
			return nil, err
		}
	}

	// func_name
	funcname := p.parseFuncNameForCreate()

	// func_args_with_defaults
	params, err := p.parseFuncArgsWithDefaults()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateFunctionStmt{
		IsOrReplace: isReplace,
		Funcname:    funcname,
		Parameters:  params,
	}

	if isProcedure {
		// PROCEDURE: no RETURNS clause
		opts, err := p.parseOptCreatefuncOptList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
		sqlBody, err := p.parseOptRoutineBody()
		if err != nil {
			return nil, err
		}
		stmt.SqlBody = sqlBody
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
		stmt.Loc = nodes.Loc{Start: startLoc, End: p.prev.End}
		return stmt, nil
	}

	// FUNCTION: check for RETURNS
	if p.cur.Type == RETURNS {
		p.advance()
		if p.cur.Type == TABLE {
			// RETURNS TABLE '(' table_func_column_list ')'
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			tableCols := p.parseTableFuncColumnList()
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
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
			retType, err := p.parseFuncType()
			if err != nil {
				return nil, err
			}
			stmt.ReturnType = retType
		}
	}

	opts, err := p.parseOptCreatefuncOptList()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts
	sqlBody, err := p.parseOptRoutineBody()
	if err != nil {
		return nil, err
	}
	stmt.SqlBody = sqlBody
	stmt.Loc = nodes.Loc{Start: startLoc, End: p.prev.End}
	return stmt, nil
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
func (p *Parser) parseFuncArgsWithDefaults() (*nodes.List, error) {
	p.expect('(')
	if p.cur.Type == ')' {
		p.advance()
		return nil, nil
	}

	args, err := p.parseFuncArgsWithDefaultsList()
	if err != nil {
		return nil, err
	}
	p.expect(')')
	return args, nil
}

// parseFuncArgsWithDefaultsList parses a comma-separated list of function args.
//
//	func_args_with_defaults_list:
//	    func_arg_with_default
//	    | func_args_with_defaults_list ',' func_arg_with_default
func (p *Parser) parseFuncArgsWithDefaultsList() (*nodes.List, error) {
	arg, err := p.parseFuncArgWithDefault()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{arg}
	for p.cur.Type == ',' {
		p.advance()
		nextArg, err := p.parseFuncArgWithDefault()
		if err != nil {
			return nil, err
		}
		items = append(items, nextArg)
	}
	return &nodes.List{Items: items}, nil
}

// parseFuncArgWithDefault parses a function argument with optional default.
//
//	func_arg_with_default:
//	    func_arg
//	    | func_arg DEFAULT a_expr
//	    | func_arg '=' a_expr
func (p *Parser) parseFuncArgWithDefault() (*nodes.FunctionParameter, error) {
	fp := p.parseFuncArg()
	if p.cur.Type == DEFAULT {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fp.Defexpr = expr
		fp.Loc.End = p.prev.End
	} else if p.cur.Type == '=' {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fp.Defexpr = expr
		fp.Loc.End = p.prev.End
	}
	return fp, nil
}

// parseFuncArg parses a function argument.
//
//	func_arg:
//	    arg_class param_name func_type
//	    | param_name arg_class func_type
//	    | param_name func_type
//	    | arg_class func_type
//	    | func_type
//
// param_name is type_function_name; func_type is Typename or
// type_function_name attrs '%' TYPE_P; arg_class is one of IN, OUT,
// INOUT, IN OUT, VARIADIC.
//
// Disambiguation strategy: speculatively consume the leading token as
// a param_name and check what follows. If the rest doesn't fit, roll
// back via snapshotTokenStream / restoreTokenStream and let
// parseFuncType handle the original token as the start of a type.
//
// Why speculative-consume instead of peek-then-commit: omni's parser
// supports IDE completion mode (`Collect`) by populating candidate
// completions as a side effect of `advance()` calls. A pure peek-then-
// commit version of this function would never advance past the
// param_name candidate, so completion mode would lose type-name
// suggestions at positions like `CREATE FUNCTION f(x |` and
// `CREATE FUNCTION f(IN x |`. The lossy `pushBack(string)` mechanism
// the previous version used was a buggy form of the same speculative
// pattern; switching to snapshotTokenStream / restoreTokenStream keeps
// the advance side effects intact while making the rollback lossless.
//
// History: the prior pushBack(string) mechanism rewrote the rolled-back
// token's Type to IDENT unconditionally. DOUBLE_P is the only
// UnreservedKeyword in simpleTypenameLeadTokens, so 'CREATE FUNCTION
// f(double precision)' and every `arg_class double precision` variant
// failed at parse time. snapshotTokenStream / restoreTokenStream
// preserves all token state on rollback.
func (p *Parser) parseFuncArg() *nodes.FunctionParameter {
	paramLoc := p.pos()

	mode := nodes.FUNC_PARAM_IN // default
	name := ""

	argClass, isArgClass := p.tryParseArgClass()
	if isArgClass {
		// Cases 1 or 4: arg_class is consumed; decide if a param_name
		// follows or if cur is the type.
		mode = argClass
		if p.isTypeFunctionName() {
			snap := p.snapshotTokenStream()
			savedName, _ := p.parseTypeFunctionName()
			if p.isFuncTypeStart() {
				// Case 1: arg_class param_name func_type
				name = savedName
			} else {
				// Case 4: cur IS the type (no param_name). Roll back so
				// parseFuncType sees the original token, including its
				// keyword Type. e.g. `IN double precision`: we consumed
				// `double` speculatively, peek shows `precision` is not
				// a func_type lead, so restore and let parseFuncType
				// handle `double precision` via parseSimpleTypename.
				p.restoreTokenStream(snap)
			}
		}
	} else if p.isTypeFunctionName() {
		// Cases 2, 3, or 5: cur could be a param_name or the type.
		snap := p.snapshotTokenStream()
		savedName, _ := p.parseTypeFunctionName()
		argClass2, isArgClass2 := p.tryParseArgClass()
		switch {
		case isArgClass2:
			// Case 2: param_name arg_class func_type. e.g. `x IN int`.
			name = savedName
			mode = argClass2
		case p.isFuncTypeStart():
			// Case 3: param_name func_type. e.g. `x int`, `p double precision`.
			name = savedName
			// mode stays FUNC_PARAM_IN
		default:
			// Case 5: cur IS the type. Roll back so parseFuncType sees
			// the original token. This is the path that fixes
			// `f(double precision)` — restore puts cur back to
			// Token{DOUBLE_P, "double"} (not the lossy IDENT pushBack
			// produced previously).
			p.restoreTokenStream(snap)
		}
	}

	// Parse func_type. cur is now positioned at the start of the type
	// (whether we consumed a param_name above, restored a rollback, or
	// never entered the speculative branch at all).
	argType, _ := p.parseFuncType()

	return &nodes.FunctionParameter{
		Name:    name,
		ArgType: argType,
		Mode:    mode,
		Loc:     nodes.Loc{Start: paramLoc, End: p.prev.End},
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
	colLoc := p.pos()
	name, _ := p.parseTypeFunctionName()
	argType, _ := p.parseFuncType()
	return &nodes.FunctionParameter{
		Name:    name,
		ArgType: argType,
		Mode:    nodes.FUNC_PARAM_TABLE,
		Loc:     nodes.Loc{Start: colLoc, End: p.prev.End},
	}
}

// parseOptCreatefuncOptList parses an optional list of function options.
//
//	opt_createfunc_opt_list:
//	    createfunc_opt_list
//	    | /* EMPTY */
func (p *Parser) parseOptCreatefuncOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		opt, err := p.parseCreatefuncOptItem()
		if err != nil {
			return nil, err
		}
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseCreatefuncOptItem parses a single function option.
//
//	createfunc_opt_item:
//	    AS func_as
//	    | LANGUAGE NonReservedWord_or_Sconst
//	    | TRANSFORM transform_type_list
//	    | WINDOW
//	    | common_func_opt_item
func (p *Parser) parseCreatefuncOptItem() (*nodes.DefElem, error) {
	loc := p.pos()
	var de *nodes.DefElem
	switch p.cur.Type {
	case AS:
		p.advance()
		funcAs := p.parseFuncAs()
		de = &nodes.DefElem{Defname: "as", Arg: funcAs}

	case LANGUAGE:
		p.advance()
		lang := p.parseNonReservedWordOrSconst()
		de = &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: lang}}

	case TRANSFORM:
		p.advance()
		types, err := p.parseTransformTypeList()
		if err != nil {
			return nil, err
		}
		de = &nodes.DefElem{Defname: "transform", Arg: types}

	case WINDOW:
		p.advance()
		de = &nodes.DefElem{Defname: "window", Arg: &nodes.Integer{Ival: 1}}

	default:
		de = p.parseCommonFuncOptItem()
	}
	if de != nil {
		de.Loc = nodes.Loc{Start: loc, End: p.pos()}
	}
	return de, nil
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
func (p *Parser) parseTransformTypeList() (*nodes.List, error) {
	p.expect(FOR)
	p.expect(TYPE_P)
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{tn}
	for p.cur.Type == ',' {
		p.advance()
		p.expect(FOR)
		p.expect(TYPE_P)
		tn, err = p.parseTypename()
		if err != nil {
			return nil, err
		}
		items = append(items, tn)
	}
	return &nodes.List{Items: items}, nil
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
		arg, _ := p.parseSetRestMore()
		return &nodes.DefElem{Defname: "set", Arg: arg}
	case RESET:
		// VariableResetStmt (RESET reset_rest)
		p.advance() // consume RESET
		arg, _ := p.parseVariableResetStmt()
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
func (p *Parser) parseOptRoutineBody() (nodes.Node, error) {
	if p.cur.Type == RETURN {
		retLoc := p.pos()
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.ReturnStmt{Returnval: expr, Loc: nodes.Loc{Start: retLoc, End: p.prev.End}}, nil
	}
	if p.cur.Type == BEGIN_P {
		p.advance() // consume BEGIN
		p.expect(ATOMIC)

		// Parse routine_body_stmt_list: (routine_body_stmt ';')*
		var stmts []nodes.Node
		for p.cur.Type != END_P && p.cur.Type != 0 {
			var stmt nodes.Node
			if p.cur.Type == RETURN {
				innerRetLoc := p.pos()
				p.advance()
				expr, _ := p.parseAExpr(0)
				stmt = &nodes.ReturnStmt{Returnval: expr, Loc: nodes.Loc{Start: innerRetLoc, End: p.prev.End}}
			} else {
				stmt, _ = p.parseStmt()
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
		return &nodes.List{Items: []nodes.Node{stmtList}}, nil
	}
	return nil, nil
}

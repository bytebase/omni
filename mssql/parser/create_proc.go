// Package parser - create_proc.go implements T-SQL CREATE PROCEDURE/FUNCTION statement parsing.
package parser

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateProcedureStmt parses a CREATE [OR ALTER] PROCEDURE statement.
//
// BNF: mssql/parser/bnf/create-procedure-transact-sql.bnf
//
//	CREATE [ OR ALTER ] { PROC | PROCEDURE }
//	    [ schema_name. ] procedure_name [ ; number ]
//	    [ { @parameter_name [ type_schema_name. ] data_type }
//	      [ VARYING ] [ NULL ] [ = default ] [ OUT | OUTPUT | READONLY ]
//	    ] [ ,...n ]
//	[ WITH <procedure_option> [ ,...n ] ]
//	[ FOR REPLICATION ]
//	AS { [ BEGIN ] sql_statement [;] [ ...n ] [ END ] }
//	[;]
//
//	-- CLR syntax:
//	CREATE [ OR ALTER ] { PROC | PROCEDURE }
//	    [ schema_name. ] procedure_name [ ; number ]
//	    [ { @parameter_name [ type_schema_name. ] data_type }
//	      [ = default ] [ OUT | OUTPUT ] [ READONLY ]
//	    ] [ ,...n ]
//	[ WITH EXECUTE AS Clause ]
//	AS { EXTERNAL NAME assembly_name.class_name.method_name }
//	[;]
//
//	<procedure_option> ::=
//	    [ ENCRYPTION ]
//	    [ RECOMPILE ]
//	    [ EXECUTE AS { CALLER | SELF | OWNER | 'user_name' } ]
//	    [ NATIVE_COMPILATION ]
//	    [ SCHEMABINDING ]
//
// The grammar is procedureStatementBody (TSql160.g:24945): a parameter list
// in parentheses must not be empty, AS is required, and the body may be
// empty.
func (p *Parser) parseCreateProcedureStmt(orAlter bool) (*nodes.CreateProcedureStmt, error) {
	loc := p.pos()

	// Completion: after CREATE/ALTER PROCEDURE → identifier or proc name
	if p.collectMode() {
		if orAlter {
			p.addRuleCandidate("proc_name")
		} else {
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	stmt := &nodes.CreateProcedureStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// Procedure name
	nameRule := "identifier"
	if orAlter {
		nameRule = "proc_ref"
	}
	name, err := p.routineParseName(nameRule)
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional procedure number: ; number
	if p.cur.Type == ';' {
		next := p.peekNext()
		if next.Type == tokICONST {
			p.advance() // consume ;
			if n, err := strconv.Atoi(p.cur.Str); err == nil {
				stmt.Number = n
			}
			p.advance() // consume number
		}
	}

	// Parameters, in parentheses or not; a parenthesized list is not empty.
	if p.cur.Type == '(' {
		p.advance()
		params, err := p.routineParseParamList(true)
		if err != nil {
			return nil, err
		}
		stmt.Params = params
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == tokVARIABLE {
		params, err := p.routineParseParamList(true)
		if err != nil {
			return nil, err
		}
		stmt.Params = params
	}

	// Completion: after CREATE PROCEDURE p [params] → AS/WITH
	if p.collectMode() {
		p.addTokenCandidate(kwAS)
		p.addTokenCandidate(kwWITH)
		return nil, errCollecting
	}

	// WITH <procedure_option> [,...n]
	// In CREATE PROC context, WITH always introduces procedure options (never a CTE).
	var opts routineOptions
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts, err = p.routineParseOptions(routineProcOptions)
		if err != nil {
			return nil, err
		}
		stmt.Options = opts.list()
	}

	// FOR REPLICATION
	if p.cur.Type == kwFOR {
		next := p.peekNext()
		if next.Type != kwREPLICATION {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume FOR
		p.advance() // consume REPLICATION
		stmt.ForReplication = true
	}

	// A natively compiled procedure takes neither ENCRYPTION, RECOMPILE
	// nor FOR REPLICATION (Msg 10794).
	if opts.native && (opts.has["ENCRYPTION"] || opts.has["RECOMPILE"] || stmt.ForReplication) {
		return nil, p.syntaxErrorAtCur()
	}

	// AS
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Completion: after CREATE PROCEDURE p AS → statement keywords
	if p.collectMode() {
		for _, kw := range topLevelKeywords {
			p.addTokenCandidate(kw)
		}
		return nil, errCollecting
	}

	// RETURN takes a value inside a routine's body (procRoutineDepth).
	p.procRoutineDepth++
	defer func() { p.procRoutineDepth-- }()
	switch {
	case opts.native:
		// The body of a natively compiled procedure is an ATOMIC block,
		// never EXTERNAL NAME (Msg 10794).
		body, err := p.routineParseAtomicBody()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	case p.cur.Type == kwEXTERNAL:
		// EXTERNAL NAME assembly.class.method (CLR)
		ext, err := p.routineParseMethodSpecifier()
		if err != nil {
			return nil, err
		}
		stmt.ExternalName = ext
	default:
		// Body (BEGIN...END block or single statement); it may be empty.
		body, err := p.routineParseStmtBody(false)
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateFunctionStmt parses a CREATE [OR ALTER] FUNCTION statement.
//
// BNF: mssql/parser/bnf/create-function-transact-sql.bnf
//
//	-- Scalar Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] [ READONLY ] }
//	    [ , ...n ]
//	  ] )
//	RETURNS return_data_type
//	    [ WITH <function_option> [ , ...n ] ]
//	    [ AS ]
//	    BEGIN
//	        function_body
//	        RETURN scalar_expression
//	    END
//	[ ; ]
//
//	-- Inline Table-Valued Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] [ READONLY ] }
//	    [ , ...n ]
//	  ] )
//	RETURNS TABLE
//	    [ WITH <function_option> [ , ...n ] ]
//	    [ AS ]
//	    RETURN [ ( ] select_stmt [ ) ]
//	[ ; ]
//
//	-- Multi-Statement Table-Valued Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( [ { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] [ READONLY ] }
//	    [ , ...n ]
//	  ] )
//	RETURNS @return_variable TABLE <table_type_definition>
//	    [ WITH <function_option> [ , ...n ] ]
//	    [ AS ]
//	    BEGIN
//	        function_body
//	        RETURN
//	    END
//	[ ; ]
//
//	-- CLR Scalar Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] }
//	    [ , ...n ]
//	)
//	RETURNS { return_data_type }
//	    [ WITH <clr_function_option> [ , ...n ] ]
//	    [ AS ] EXTERNAL NAME <method_specifier>
//	[ ; ]
//
//	-- CLR Table-Valued Function:
//	CREATE [ OR ALTER ] FUNCTION [ schema_name. ] function_name
//	( { @parameter_name [ AS ] [ type_schema_name. ] parameter_data_type [ NULL ]
//	    [ = default ] }
//	    [ , ...n ]
//	)
//	RETURNS TABLE <clr_table_type_definition>
//	    [ WITH <clr_function_option> [ , ...n ] ]
//	    [ ORDER ( <order_clause> ) ]
//	    [ AS ] EXTERNAL NAME <method_specifier>
//	[ ; ]
//
//	<function_option> ::=
//	{
//	    [ ENCRYPTION ]
//	  | [ SCHEMABINDING ]
//	  | [ RETURNS NULL ON NULL INPUT | CALLED ON NULL INPUT ]
//	  | [ EXECUTE_AS_Clause ]
//	  | [ INLINE = { ON | OFF } ]
//	}
//
//	<method_specifier> ::= assembly_name.class_name.method_name
//
// The grammar is functionStatementBody and functionReturnTypeAndBody
// (TSql160.g:12984, TSql160.g:13029).
func (p *Parser) parseCreateFunctionStmt(orAlter bool) (*nodes.CreateFunctionStmt, error) {
	loc := p.pos()

	// Completion: after CREATE/ALTER FUNCTION → identifier or func name
	if p.collectMode() {
		if orAlter {
			p.addRuleCandidate("func_name")
		} else {
			p.addRuleCandidate("identifier")
		}
		return nil, errCollecting
	}

	stmt := &nodes.CreateFunctionStmt{
		OrAlter: orAlter,
		Loc:     nodes.Loc{Start: loc, End: -1},
	}

	// Function name
	nameRule := "identifier"
	if orAlter {
		nameRule = "func_name"
	}
	name, err := p.routineParseName(nameRule)
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Parameters in parentheses, which are required.
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Completion: inside function params → @ variable
	if p.collectMode() {
		p.addRuleCandidate("variable")
		return nil, errCollecting
	}
	if p.cur.Type != ')' {
		params, err := p.routineParseParamList(false)
		if err != nil {
			return nil, err
		}
		stmt.Params = params
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// RETURNS type or RETURNS TABLE
	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, err
	}
	// Completion: after RETURNS → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		p.addTokenCandidate(kwTABLE)
		return nil, errCollecting
	}

	// kind is the function kind, which decides the options and the body.
	var kind routineFunctionKind
	switch p.cur.Type {
	case kwTABLE:
		p.advance()
		stmt.ReturnsTable = &nodes.ReturnsTableDef{
			Loc: nodes.Loc{Start: p.pos(), End: -1},
		}
		kind = routineInlineTable
		if p.cur.Type == '(' {
			// RETURNS TABLE ( ... ): a CLR table-valued function.
			cols, err := p.routineParseTableDefinition()
			if err != nil {
				return nil, err
			}
			stmt.ReturnsTable.Columns = cols
			kind = routineCLRTable
		}
		stmt.ReturnsTable.Loc.End = p.prevEnd()
	case tokVARIABLE:
		// RETURNS @var TABLE ( ... ): a multi-statement table-valued function.
		// SQL Server refuses AS between the variable and TABLE (Msg 156).
		varName := p.cur.Str
		p.advance() // consume @var
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		stmt.ReturnsTable = &nodes.ReturnsTableDef{
			Variable: varName,
			Loc:      nodes.Loc{Start: p.pos(), End: -1},
		}
		cols, err := p.routineParseTableDefinition()
		if err != nil {
			return nil, err
		}
		stmt.ReturnsTable.Columns = cols
		stmt.ReturnsTable.Loc.End = p.prevEnd()
		kind = routineMultiTable
	default:
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		if dt == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.ReturnType = dt
		kind = routineScalar
	}

	// WITH <function_option> [,...n]
	var opts routineOptions
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		opts, err = p.routineParseOptions(routineFunctionOptions)
		if err != nil {
			return nil, err
		}
		stmt.Options = opts.list()
	}

	// ORDER ( column [ ASC | DESC ] [,...n] ) for CLR table-valued functions
	// (clrTableValuedFunctionOrderHint, TSql160.g:13220).
	if kind == routineCLRTable && p.cur.Type == kwORDER {
		order, err := p.routineParseOrderHint()
		if err != nil {
			return nil, err
		}
		stmt.OrderClause = order
	}

	// AS
	p.match(kwAS)

	// A scalar function becomes a CLR one with EXTERNAL NAME.
	if kind == routineScalar && p.cur.Type == kwEXTERNAL {
		kind = routineCLRScalar
	}
	// A natively compiled function takes neither ENCRYPTION nor RETURNS
	// NULL ON NULL INPUT (Msg 10794).
	if !opts.allowedFor(kind) || opts.native && (opts.has["ENCRYPTION"] || opts.has["RETURNS NULL ON NULL INPUT"]) {
		return nil, p.newParseError(loc, "an option is not valid for this kind of function")
	}

	switch kind {
	case routineCLRScalar, routineCLRTable:
		if p.cur.Type != kwEXTERNAL {
			return nil, p.syntaxErrorAtCur()
		}
		ext, err := p.routineParseMethodSpecifier()
		if err != nil {
			return nil, err
		}
		stmt.ExternalName = ext
	case routineInlineTable:
		// Inline table-valued function: RETURN [ ( ] select_stmt [ ) ]
		// The select_statement may carry its own CTE list, with or without
		// the parentheses (SQL Server 2022 oracle).
		if p.cur.Type != kwRETURN {
			return nil, p.syntaxErrorAtCur()
		}
		retLoc := p.pos()
		p.advance() // consume RETURN
		hasParen := false
		if p.cur.Type == '(' {
			hasParen = true
			p.advance()
		}
		selectStmt, err := p.parseSelectStmtWithCTE()
		if err != nil {
			return nil, err
		}
		if selectStmt == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if hasParen {
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		stmt.Body = &nodes.ReturnStmt{
			Value: &nodes.SubqueryExpr{Query: selectStmt, Loc: nodes.Loc{Start: retLoc, End: p.prevEnd()}},
			Loc:   nodes.Loc{Start: retLoc, End: p.prevEnd()},
		}
	default:
		// Scalar and multi-statement table-valued functions: BEGIN ... END,
		// an ATOMIC block when natively compiled. RETURN takes a value here
		// (procRoutineDepth).
		p.procRoutineDepth++
		defer func() { p.procRoutineDepth-- }()
		var body nodes.StmtNode
		if opts.native {
			body, err = p.routineParseAtomicBody()
		} else {
			body, err = p.routineParseBeginEnd()
		}
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// routineIsWord reports whether the current token is the option word word,
// written in upper case: its keyword token, or an identifier spelled so that
// is not bracketed or quoted. SQL Server refuses a delimited option word
// (`WITH [RECOMPILE]`, Msg 102).
func (p *Parser) routineIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// routineMatchWord consumes the current token if it is the option word word.
func (p *Parser) routineMatchWord(word string) bool {
	if p.routineIsWord(word) {
		p.advance()
		return true
	}
	return false
}

// routineExpectWord consumes the option word word or fails.
func (p *Parser) routineExpectWord(word string) error {
	if !p.routineMatchWord(word) {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// routineParseName parses the name of a procedure, function or trigger,
// which is required. SQL Server refuses a database prefix on a module it
// creates (Msg 166), though schemaObjectThreePartName reads one; Parse still
// takes it, since mssql/review replays what the server does with it.
func (p *Parser) routineParseName(completionRule string) (*nodes.TableRef, error) {
	name, err := p.parseObjectRef(completionRule)
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return name, nil
}

// routineParseParamList parses procedureParameterList (TSql160.g:25028) or
// functionParameterList (TSql160.g:13002): at least one parameter, separated
// by commas, no name twice (Msg 134). Names are compared as written: whether
// @a and @A are the same variable depends on the collation, which Parse does
// not know.
func (p *Parser) routineParseParamList(isProc bool) (*nodes.List, error) {
	var params []nodes.Node
	seen := map[string]bool{}
	for {
		param, err := p.routineParseParam(isProc)
		if err != nil {
			return nil, err
		}
		key := param.Name
		if seen[key] {
			return nil, p.newParseError(param.Loc.Start, "the variable name "+param.Name+" has already been declared")
		}
		seen[key] = true
		params = append(params, param)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: params}, nil
}

// routineParseParam parses one parameter of a procedure or a function.
//
//	procedureParameter (TSql160.g:25044):
//	    @name [ AS ] { CURSOR VARYING { OUTPUT | OUT }
//	                 | data_type [ NULL | NOT NULL ] [ = default ] [ READONLY | OUTPUT | OUT ] }
//	functionParameter (TSql160.g:13018):
//	    @name [ AS ] data_type [ NULL | NOT NULL ] [ = default ] [ READONLY ]
//
// A cursor parameter must be VARYING and OUTPUT, in that order (Msg 1051);
// a function refuses CURSOR (Msg 5306) and OUTPUT (Msg 181). VARYING follows
// only CURSOR.
func (p *Parser) routineParseParam(isProc bool) (*nodes.ParamDef, error) {
	if p.cur.Type != tokVARIABLE {
		return nil, p.syntaxErrorAtCur()
	}

	loc := p.pos()
	param := &nodes.ParamDef{
		Name: p.cur.Str,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}
	p.advance() // consume @param

	p.match(kwAS)

	// CURSOR VARYING { OUTPUT | OUT } (cursorProcedureParameter, TSql160.g:25108).
	if p.cur.Type == kwCURSOR {
		if !isProc {
			return nil, p.syntaxErrorAtCur()
		}
		param.DataType = &nodes.DataType{Name: "CURSOR", Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}}
		p.advance()
		if _, err := p.expect(kwVARYING); err != nil {
			return nil, err
		}
		param.Varying = true
		if p.cur.Type != kwOUTPUT && p.cur.Type != kwOUT {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		param.Output = true
		param.Loc.End = p.prevEnd()
		return param, nil
	}

	// scalarProcedureParameter (TSql160.g:25176). VARYING names no type: it
	// follows CURSOR only (Msg 102).
	if p.cur.Type == kwVARYING {
		return nil, p.syntaxErrorAtCur()
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.syntaxErrorAtCur()
	}
	param.DataType = dt

	// [ NULL | NOT NULL ] (nullableConstraint)
	switch {
	case p.cur.Type == kwNULL:
		param.Null = true
		p.advance()
	case p.cur.Type == kwNOT && p.peekNext().Type == kwNULL:
		param.NotNull = true
		p.advance()
		p.advance()
	}

	// [ = default ]
	if p.cur.Type == '=' {
		p.advance()
		def, err := p.routineParseParamDefault()
		if err != nil {
			return nil, err
		}
		param.Default = def
	}

	// [ READONLY | OUTPUT | OUT ]
	switch p.cur.Type {
	case kwREADONLY:
		param.ReadOnly = true
		p.advance()
	case kwOUTPUT, kwOUT:
		if !isProc {
			return nil, p.syntaxErrorAtCur()
		}
		param.Output = true
		p.advance()
	}

	param.Loc.End = p.prevEnd()
	return param, nil
}

// routineParseParamDefault parses the default of a parameter:
// possibleNegativeConstantOrIdentifierWithDefault (TSql160.g:25155), a
// literal, a negative number, a bare identifier, NULL or DEFAULT. SQL Server
// refuses an expression (`= 1+1`) and a variable (`= @b`).
func (p *Parser) routineParseParamDefault() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case '-':
		switch p.peekNext().Type {
		case tokICONST, tokFCONST, tokMONEY:
		default:
			return nil, p.syntaxErrorAtCur()
		}
		loc := p.pos()
		p.advance() // consume -
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{
			Op:      nodes.UnaryMinus,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST, tokMONEY, kwNULL, kwDEFAULT:
		return p.parsePrimary()
	}
	if p.isIdentLike() {
		tok := p.advance()
		return &nodes.ColumnRef{
			Column: tok.Str,
			Loc:    nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// routineParseTableDefinition parses ( tableDefinition ), the columns of
// RETURNS TABLE: column definitions, unnamed PRIMARY KEY, UNIQUE and CHECK
// constraints, and indexes, at least one of them (declareTableBodyMain,
// TSql160.g:16716).
func (p *Parser) routineParseTableDefinition() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.Node
	for {
		var item nodes.Node
		switch p.cur.Type {
		case kwPRIMARY, kwUNIQUE, kwCHECK:
			c, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}
			if c != nil {
				item = c
			}
		case kwINDEX:
			idx, err := p.parseInlineTableIndex()
			if err != nil {
				return nil, err
			}
			if idx != nil {
				item = idx
			}
		default:
			col, err := p.parseTypedColumnDef()
			if err != nil {
				return nil, err
			}
			if col != nil {
				item = col
			}
		}
		if item == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, item)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// routineParseOrderHint parses ORDER ( column [ ASC | DESC ] [ ,...n ] ), the
// order hint of a CLR table-valued function (bulkInsertSortOrderOption,
// TSql160.g:6347).
func (p *Parser) routineParseOrderHint() (*nodes.List, error) {
	p.advance() // consume ORDER
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.Node
	for {
		loc := p.pos()
		col, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		item := &nodes.OrderByItem{
			Expr: &nodes.ColumnRef{Column: col, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}},
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		switch p.cur.Type {
		case kwASC:
			item.SortDir = nodes.SortAsc
			p.advance()
		case kwDESC:
			item.SortDir = nodes.SortDesc
			p.advance()
		}
		item.Loc.End = p.prevEnd()
		items = append(items, item)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// routineParseMethodSpecifier parses EXTERNAL NAME assembly_name.class_name.
// method_name (methodSpecifier, TSql160.g:5117): exactly three identifiers.
// The current token is EXTERNAL.
func (p *Parser) routineParseMethodSpecifier() (string, error) {
	p.advance() // consume EXTERNAL
	if err := p.routineExpectWord("NAME"); err != nil {
		return "", err
	}
	parts := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		if i > 0 {
			if _, err := p.expect('.'); err != nil {
				return "", err
			}
		}
		part, ok := p.parseIdentifier()
		if !ok {
			return "", p.syntaxErrorAtCur()
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "."), nil
}

// routineParseStmtBody parses the body of a procedure or a trigger: optional
// semicolons, then a statement. The statements after the first stay
// statements of the batch. A procedure's body may be empty
// (procedureStatementBody, TSql160.g:24945); a trigger's may not
// (triggerStatementBody, TSql160.g:25231).
func (p *Parser) routineParseStmtBody(required bool) (nodes.StmtNode, error) {
	for p.cur.Type == ';' {
		p.advance()
	}
	body, err := p.parseStmt()
	if err != nil {
		return nil, err
	}
	if isNilNode(body) {
		if required || p.cur.Type != tokEOF {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}
	return body, nil
}

// routineParseBeginEnd parses BEGIN ... END, the body of a scalar or a
// multi-statement table-valued function (functionReturnTypeAndBody,
// TSql160.g:13029).
func (p *Parser) routineParseBeginEnd() (nodes.StmtNode, error) {
	if p.cur.Type != kwBEGIN {
		return nil, p.syntaxErrorAtCur()
	}
	loc := p.pos()
	body, err := p.parseBeginStmt()
	if err != nil {
		return nil, err
	}
	if _, ok := body.(*nodes.BeginEndStmt); !ok {
		return nil, p.newParseError(loc, "the body of a function must be a BEGIN ... END block")
	}
	return body, nil
}

// routineParseAtomicBody parses the body of a natively compiled module, which
// must be an ATOMIC block (Msg 10783):
//
//	BEGIN ATOMIC WITH ( <atomic_option> [ ,...n ] ) [;] sql_statement [;] [ ...n ] END
//
//	<atomic_option> ::=
//	    TRANSACTION ISOLATION LEVEL = { SNAPSHOT | REPEATABLE READ | SERIALIZABLE }
//	  | LANGUAGE = { 'language' | language }
//	  | DATEFIRST = number
//	  | DATEFORMAT = 'format'
//	  | DELAYED_DURABILITY = { ON | OFF }
//
// beginEndAtomicBlock and atomicBlockOption (TSql160.g:16190,
// TSql160.g:16229). TRANSACTION ISOLATION LEVEL and LANGUAGE are required
// (Msg 10784), no option comes twice (Msg 1039), and SQL Server refuses
// READ COMMITTED and READ UNCOMMITTED in a natively compiled module
// (Msg 10794).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-procedure-transact-sql
func (p *Parser) routineParseAtomicBody() (*nodes.BeginEndStmt, error) {
	loc := p.pos()
	if _, err := p.expect(kwBEGIN); err != nil {
		return nil, err
	}
	if err := p.routineExpectWord("ATOMIC"); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		var name, value string
		switch {
		case p.cur.Type == kwTRANSACTION:
			p.advance()
			if err := p.routineExpectWord("ISOLATION"); err != nil {
				return nil, err
			}
			if err := p.routineExpectWord("LEVEL"); err != nil {
				return nil, err
			}
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			name = "TRANSACTION ISOLATION LEVEL"
			switch {
			case p.routineMatchWord("SNAPSHOT"):
				value = "SNAPSHOT"
			case p.routineMatchWord("SERIALIZABLE"):
				value = "SERIALIZABLE"
			case p.routineMatchWord("REPEATABLE"):
				if err := p.routineExpectWord("READ"); err != nil {
					return nil, err
				}
				value = "REPEATABLE READ"
			default:
				return nil, p.syntaxErrorAtCur()
			}
		case p.routineIsWord("LANGUAGE"):
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			name = "LANGUAGE"
			switch {
			case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
				value = p.cur.Str
				p.advance()
			case p.isIdentLike():
				value = p.cur.Str
				p.advance()
			default:
				return nil, p.syntaxErrorAtCur()
			}
		case p.routineIsWord("DATEFIRST"):
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			name = "DATEFIRST"
			if p.cur.Type == '-' {
				value = "-"
				p.advance()
			}
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			value += p.cur.Str
			p.advance()
		case p.routineIsWord("DATEFORMAT"):
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			name = "DATEFORMAT"
			if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			value = p.cur.Str
			p.advance()
		case p.routineIsWord("DELAYED_DURABILITY"):
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			name = "DELAYED_DURABILITY"
			switch p.cur.Type {
			case kwON:
				value = "ON"
			case kwOFF:
				value = "OFF"
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if seen[name] {
			return nil, p.newParseError(optLoc, "option "+name+" is specified more than once")
		}
		seen[name] = true
		opts = append(opts, &nodes.AtomicBlockOption{Name: name, Value: value, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if !seen["TRANSACTION ISOLATION LEVEL"] || !seen["LANGUAGE"] {
		return nil, p.newParseError(loc, "BEGIN ATOMIC requires TRANSACTION ISOLATION LEVEL and LANGUAGE")
	}

	var stmts []nodes.Node
	for {
		for p.cur.Type == ';' {
			p.advance()
		}
		if p.cur.Type == kwEND {
			break
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if isNilNode(s) {
			return nil, p.syntaxErrorAtCur()
		}
		stmts = append(stmts, s)
	}
	if len(stmts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume END
	return &nodes.BeginEndStmt{
		Stmts:         &nodes.List{Items: stmts},
		AtomicOptions: &nodes.List{Items: opts},
		Loc:           nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
}

// routineOptionKind is a bit set naming the statements an option belongs to.
type routineOptionKind int

const (
	routineProcOptions routineOptionKind = 1 << iota
	routineTriggerOptions
	routineFunctionOptions
)

// routineFunctionKind is the kind of a function, which decides what options
// it takes and what its body is.
type routineFunctionKind int

const (
	routineScalar routineFunctionKind = iota
	routineCLRScalar
	routineInlineTable
	routineMultiTable
	routineCLRTable
)

// routineFunctionOptionKinds lists the kinds of function each function
// option is valid for. SQL Server refuses the others with Msg 487 ("An
// invalid option was specified for the statement CREATE/ALTER FUNCTION"):
// an inline table-valued function takes only ENCRYPTION and SCHEMABINDING,
// a CLR function neither of them, and NATIVE_COMPILATION goes only with a
// T-SQL scalar function.
var routineFunctionOptionKinds = map[string][]routineFunctionKind{
	"ENCRYPTION":                 {routineScalar, routineInlineTable, routineMultiTable},
	"SCHEMABINDING":              {routineScalar, routineInlineTable, routineMultiTable},
	"NATIVE_COMPILATION":         {routineScalar},
	"RETURNS NULL ON NULL INPUT": {routineScalar, routineCLRScalar},
	"CALLED ON NULL INPUT":       {routineScalar, routineCLRScalar},
	"INLINE":                     {routineScalar},
	"EXECUTE AS":                 {routineScalar, routineCLRScalar, routineMultiTable, routineCLRTable},
}

// routineOptions is a parsed WITH list of a procedure, function or trigger.
type routineOptions struct {
	items  []nodes.Node
	has    map[string]bool
	native bool // NATIVE_COMPILATION
}

func (o routineOptions) list() *nodes.List {
	if len(o.items) == 0 {
		return nil
	}
	return &nodes.List{Items: o.items}
}

// allowedFor reports whether every option in o is valid for a function of
// kind k.
func (o routineOptions) allowedFor(k routineFunctionKind) bool {
	for name := range o.has {
		ok := false
		for _, allowed := range routineFunctionOptionKinds[name] {
			if allowed == k {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// routineParseOptions parses the options after WITH of a procedure
// (procedureOption, TSql160.g:24878), a trigger (triggerOption,
// TSql160.g:24837) or a function (functionAttribute, TSql160.g:13203):
//
//	procedure: ENCRYPTION | RECOMPILE | NATIVE_COMPILATION | SCHEMABINDING | EXECUTE AS ...
//	trigger:   ENCRYPTION | NATIVE_COMPILATION | SCHEMABINDING | EXECUTE AS ...
//	function:  ENCRYPTION | SCHEMABINDING | NATIVE_COMPILATION | EXECUTE AS ...
//	           | RETURNS NULL ON NULL INPUT | CALLED ON NULL INPUT | INLINE = { ON | OFF }
//
// ScriptDOM reads the one-word options as an identifier; SQL Server knows
// only these words (Msg 195 "is not a recognized option"; VIEW_METADATA and
// RECOMPILE are refused where they do not belong, Msg 487). No option comes
// twice (Msg 1039), RETURNS NULL ON NULL INPUT and CALLED ON NULL INPUT
// conflict (Msg 1052), and a procedure or a trigger takes SCHEMABINDING
// exactly when it takes NATIVE_COMPILATION, which a function requires too
// (Msg 10796).
//
// Procedures and functions keep each option as a nodes.String
// ("EXECUTE AS CALLER", "INLINE = ON"); triggers as a nodes.TriggerOption.
func (p *Parser) routineParseOptions(kind routineOptionKind) (routineOptions, error) {
	opts := routineOptions{has: map[string]bool{}}
	for {
		optLoc := p.pos()
		var name, value string
		switch {
		case p.cur.Type == kwEXEC || p.cur.Type == kwEXECUTE:
			v, err := p.routineParseExecuteAs()
			if err != nil {
				return opts, err
			}
			name, value = "EXECUTE AS", v
		case p.routineIsWord("ENCRYPTION"):
			name = "ENCRYPTION"
			p.advance()
		case p.routineIsWord("SCHEMABINDING"):
			name = "SCHEMABINDING"
			p.advance()
		case p.routineIsWord("NATIVE_COMPILATION"):
			name = "NATIVE_COMPILATION"
			p.advance()
		case kind == routineProcOptions && p.routineIsWord("RECOMPILE"):
			name = "RECOMPILE"
			p.advance()
		case kind == routineFunctionOptions && p.cur.Type == kwRETURNS:
			p.advance()
			if _, err := p.expect(kwNULL); err != nil {
				return opts, err
			}
			if _, err := p.expect(kwON); err != nil {
				return opts, err
			}
			if _, err := p.expect(kwNULL); err != nil {
				return opts, err
			}
			if err := p.routineExpectWord("INPUT"); err != nil {
				return opts, err
			}
			name = "RETURNS NULL ON NULL INPUT"
		case kind == routineFunctionOptions && p.routineIsWord("CALLED"):
			p.advance()
			if _, err := p.expect(kwON); err != nil {
				return opts, err
			}
			if _, err := p.expect(kwNULL); err != nil {
				return opts, err
			}
			if err := p.routineExpectWord("INPUT"); err != nil {
				return opts, err
			}
			name = "CALLED ON NULL INPUT"
		case kind == routineFunctionOptions && p.routineIsWord("INLINE"):
			// INLINE = { ON | OFF } (inlineFunctionOption, TSql160.g:13181)
			p.advance()
			if _, err := p.expect('='); err != nil {
				return opts, err
			}
			switch p.cur.Type {
			case kwON:
				value = "ON"
			case kwOFF:
				value = "OFF"
			default:
				return opts, p.syntaxErrorAtCur()
			}
			p.advance()
			name = "INLINE"
		default:
			return opts, p.syntaxErrorAtCur()
		}
		if opts.has[name] {
			return opts, p.newParseError(optLoc, "option "+name+" is specified more than once")
		}
		opts.has[name] = true
		end := p.prevEnd()
		if kind == routineTriggerOptions {
			opts.items = append(opts.items, &nodes.TriggerOption{Name: name, Value: value, Loc: nodes.Loc{Start: optLoc, End: end}})
		} else {
			str := name
			switch name {
			case "EXECUTE AS":
				str = "EXECUTE AS " + value
			case "INLINE":
				str = "INLINE = " + value
			}
			opts.items = append(opts.items, &nodes.String{Str: str})
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if opts.has["RETURNS NULL ON NULL INPUT"] && opts.has["CALLED ON NULL INPUT"] {
		return opts, p.syntaxErrorAtCur()
	}
	opts.native = opts.has["NATIVE_COMPILATION"]
	if opts.native && !opts.has["SCHEMABINDING"] ||
		kind != routineFunctionOptions && opts.has["SCHEMABINDING"] && !opts.native {
		return opts, p.newParseError(p.pos(), "SCHEMABINDING is supported only for natively compiled modules, and is required for those modules")
	}
	return opts, nil
}

// routineParseExecuteAs parses EXECUTE AS { CALLER | SELF | OWNER | 'user_name' }
// (executeAsClause, TSql160.g:2380; the action takes CALLER, SELF and OWNER
// as words), with EXEC or EXECUTE the current token. It returns CALLER, SELF,
// OWNER, or the user name.
func (p *Parser) routineParseExecuteAs() (string, error) {
	p.advance() // consume EXEC/EXECUTE
	if _, err := p.expect(kwAS); err != nil {
		return "", err
	}
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		v := p.cur.Str
		p.advance()
		return v, nil
	}
	for _, w := range []string{"CALLER", "SELF", "OWNER"} {
		if p.routineMatchWord(w) {
			return w, nil
		}
	}
	return "", p.syntaxErrorAtCur()
}

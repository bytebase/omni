// Package parser - execute.go implements T-SQL EXEC/EXECUTE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseExecStmt parses an EXEC/EXECUTE statement.
//
// BNF: mssql/parser/bnf/execute-transact-sql.bnf
//
//	-- Execute a stored procedure or function:
//	[ { EXEC | EXECUTE } ]
//	    {
//	      [ @return_status = ]
//	      { module_name [ ;number ] | @module_name_var }
//	        [ [ @parameter = ] { value
//	                           | @variable [ OUTPUT ]
//	                           | [ DEFAULT ]
//	                           }
//	        ]
//	      [ ,...n ]
//	      [ WITH <execute_option> [ ,...n ] ]
//	    }
//	[ ; ]
//
//	-- Execute a character string:
//	{ EXEC | EXECUTE }
//	    ( { @string_variable | [ N ]'tsql_string' } [ + ...n ] )
//	    [ AS { LOGIN | USER } = ' name ' ]
//	[ ; ]
//
//	-- Execute a pass-through command against a linked server:
//	{ EXEC | EXECUTE }
//	    ( { @string_variable | [ N ] 'command_string [ ? ]' } [ + ...n ]
//	        [ { , { value | @variable [ OUTPUT ] } } [ ...n ] ]
//	    )
//	    [ AS { LOGIN | USER } = ' name ' ]
//	    [ AT linked_server_name ]
//	    [ AT DATA_SOURCE data_source_name ]
//	[ ; ]
//
//	<execute_option> ::=
//	{
//	        RECOMPILE
//	    | { RESULT SETS UNDEFINED }
//	    | { RESULT SETS NONE }
//	    | { RESULT SETS ( <result_sets_definition> [,...n ] ) }
//	}
//
//	<result_sets_definition> ::=
//	{
//	    (
//	         { column_name
//	           data_type
//	         [ COLLATE collation_name ]
//	         [ NULL | NOT NULL ] }
//	         [,...n ]
//	    )
//	    | AS OBJECT
//	        [ db_name . [ schema_name ] . | schema_name . ]
//	        {table_name | view_name | table_valued_function_name }
//	    | AS TYPE [ schema_name.]table_type_name
//	    | AS FOR XML
//	}
func (p *Parser) parseExecStmt() (*nodes.ExecStmt, error) {
	loc := p.pos()
	p.advance() // consume EXEC or EXECUTE

	// Completion: after EXEC/EXECUTE → proc_ref
	if p.collectMode() {
		p.addRuleCandidate("proc_ref")
		return nil, errCollecting
	}

	stmt := &nodes.ExecStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Check for EXEC ('string') syntax -- execute character string
	if p.cur.Type == '(' {
		if err := p.parseExecString(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// Check for @return_var = proc_name
	if p.cur.Type == tokVARIABLE {
		next := p.peekNext()
		if next.Type == '=' {
			stmt.ReturnVar = p.cur.Str
			p.advance() // consume @var
			p.advance() // consume =
		}
	}

	// Procedure name, or a variable holding it (varObjectReference,
	// TSql160.g:26342), which a return variable may not precede: execTypes
	// (TSql160.g:26231) takes `variable = execProcEx`, and execProcEx's
	// variable form is execProc's.
	if p.cur.Type == tokVARIABLE {
		stmt.ModuleVar = p.cur.Str
		p.advance()
	} else {
		name, err := p.parseObjectRef("proc_ref")
		if err != nil {
			return nil, err
		}
		if name == nil {
			return nil, p.newParseError(p.cur.Loc, "expected procedure name after EXEC")
		}
		stmt.Name = name

		// Optional procedure number: ;number (EXEC sp_test;1), procNumOpt
		// (TSql160.g:26505).
		if p.cur.Type == ';' {
			next := p.peekNext()
			if next.Type == tokICONST {
				p.advance() // consume ;
				p.advance() // consume number (stored in Name context)
			}
		}
	}

	// Arguments
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != ')' &&
		p.cur.Type != kwGO && p.cur.Type != kwEND &&
		!p.isStatementStart() && !p.isExecWithOption() {
		args, err := p.parseExecArgList(true)
		if err != nil {
			return nil, err
		}
		if len(args) > 0 {
			stmt.Args = &nodes.List{Items: args}
		}
	}

	// WITH <execute_option> [,...n]
	if p.isExecWithOption() {
		p.advance() // consume WITH
		opts, err := p.parseExecOptionList()
		if err != nil {
			return nil, err
		}
		stmt.WithOptions = opts
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// execIsWord reports whether the current token is the word w (upper case):
// the keyword w, or an undelimited identifier spelled w when w is not a
// keyword.
func (p *Parser) execIsWord(w string) bool {
	if kw := lookupKeyword(w); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == w
}

// execExpectWord consumes the word w or fails at the current token.
func (p *Parser) execExpectWord(w string) error {
	if !p.execIsWord(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseExecString parses the character string form after EXEC:
//
//	( execSqlList [ , setParamList ] ) [ AS { LOGIN | USER } = 'name' ]
//	    [ AT linked_server_name | AT DATA_SOURCE data_source_name ]
//
// execTypes, execStrTypes and execSqlList (TSql160.g:26231, 26287, 26371):
// the string is string literals and variables joined by +, which is all
// SQL Server takes there (`EXEC (1)` and `EXEC (('a'))` are Msg 102). The
// principal of AS is a string (execCtxtStmtPrincipal, TSql160.g:26283),
// which SQL Server requires although the grammar makes it optional. Pass-
// through parameters need AT: SQL Server refuses `EXEC ('a', 1)` (Msg 102).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/execute-transact-sql
func (p *Parser) parseExecString(stmt *nodes.ExecStmt) error {
	p.advance() // consume '('
	expr, err := p.parseExecSqlList()
	if err != nil {
		return err
	}
	stmt.ExecString = expr
	hasParams := false
	if p.cur.Type == ',' {
		p.advance() // consume ','
		args, err := p.parseExecArgList(false)
		if err != nil {
			return err
		}
		stmt.Args = &nodes.List{Items: args}
		hasParams = true
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}

	// AS { LOGIN | USER } = 'name'
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		isLogin := p.cur.Type == kwLOGIN
		if !isLogin && p.cur.Type != kwUSER {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume LOGIN or USER
		if _, err := p.expect('='); err != nil {
			return err
		}
		if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
			return p.syntaxErrorAtCur()
		}
		if isLogin {
			stmt.AsLogin = p.cur.Str
		} else {
			stmt.AsUser = p.cur.Str
		}
		p.advance()
	}

	// AT linked_server_name (linkedServer_AT, TSql160.g:9567) | AT
	// DATA_SOURCE data_source_name.
	if p.cur.Type == kwAT {
		p.advance() // consume AT
		if p.cur.Type == kwDATA_SOURCE && p.peekNext().Type != tokEOF && p.peekNext().Type != ';' {
			p.advance() // consume DATA_SOURCE
			name, ok := p.parseIdentifier()
			if !ok {
				return p.syntaxErrorAtCur()
			}
			stmt.AtDataSource = name
		} else {
			name, ok := p.parseIdentifier()
			if !ok {
				return p.syntaxErrorAtCur()
			}
			stmt.AtServer = name
		}
	} else if hasParams {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseExecSqlList parses execSqlList (TSql160.g:26371): a string literal or
// variable, then more joined by +.
func (p *Parser) parseExecSqlList() (nodes.ExprNode, error) {
	start := p.pos()
	var expr nodes.ExprNode
	for {
		switch p.cur.Type {
		case tokSCONST, tokNSCONST, tokVARIABLE, tokSYSVARIABLE:
		default:
			return nil, p.syntaxErrorAtCur()
		}
		item, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			expr = item
		} else {
			expr = &nodes.BinaryExpr{
				Op:    nodes.BinOpAdd,
				Left:  expr,
				Right: item,
				Loc:   nodes.Loc{Start: start, End: p.prevEnd()},
			}
		}
		if p.cur.Type != '+' {
			return expr, nil
		}
		p.advance() // consume '+'
	}
}

// isExecWithOption returns true if current token is WITH followed by an exec option.
func (p *Parser) isExecWithOption() bool {
	if p.cur.Type != kwWITH {
		return false
	}
	next := p.peekNext()
	if next.Str != "" {
		s := strings.ToUpper(next.Str)
		if s == "RECOMPILE" || s == "RESULT" {
			return true
		}
	}
	return false
}

// parseExecOptionList parses WITH execute_option [, execute_option ...]
// (executeStatement, TSql160.g:26055). SQL Server refuses an option given
// twice (`WITH RECOMPILE, RECOMPILE` and two RESULT SETS are Msg 102).
func (p *Parser) parseExecOptionList() (*nodes.List, error) {
	var items []nodes.Node
	seen := map[string]bool{}
	for {
		opt, err := p.parseExecOption()
		if err != nil {
			return nil, err
		}
		key := opt.Str
		if strings.HasPrefix(key, "RESULT SETS") {
			key = "RESULT SETS"
		}
		if seen[key] {
			return nil, p.newParseError(p.prev.Loc, "duplicate EXECUTE option "+key)
		}
		seen[key] = true
		items = append(items, opt)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}, nil
}

// parseExecOption parses a single execute option (executeOption,
// TSql160.g:26084; executeOptionResultSets, TSql160.g:26099).
//
//	RECOMPILE | RESULT SETS UNDEFINED | RESULT SETS NONE
//	| RESULT SETS ( <result_sets_definition> [,...n] )
func (p *Parser) parseExecOption() (*nodes.String, error) {
	if p.cur.Type == kwRECOMPILE {
		p.advance()
		return &nodes.String{Str: "RECOMPILE"}, nil
	}
	if p.cur.Type != kwRESULT {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume RESULT
	if _, err := p.expect(kwSETS); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case kwUNDEFINED:
		p.advance()
		return &nodes.String{Str: "RESULT SETS UNDEFINED"}, nil
	case kwNONE:
		p.advance()
		return &nodes.String{Str: "RESULT SETS NONE"}, nil
	case '(':
		p.advance() // consume '('
		for {
			if err := p.parseExecResultSetDefinition(); err != nil {
				return nil, err
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.String{Str: "RESULT SETS DEFINED"}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseExecResultSetDefinition parses one result set definition
// (resultSetDefinition, TSql160.g:26143):
//
//	( column_name data_type [ COLLATE collation ] [ NULL | NOT NULL ] [ ,...n ] )
//	| AS FOR XML
//	| AS { OBJECT | TYPE } [ [ server . ] database . ] [ schema . ] name
func (p *Parser) parseExecResultSetDefinition() error {
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			return p.execExpectWord("XML")
		}
		if !p.execIsWord("OBJECT") && !p.execIsWord("TYPE") {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume OBJECT or TYPE
		name, err := p.parseObjectRef("table_ref")
		if err != nil {
			return err
		}
		if name == nil {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		// resultColumnDefinition (TSql160.g:26169): columnDefinitionBasic
		// (TSql160.g:29519) and an optional NULL or NOT NULL.
		if _, ok := p.parseIdentifier(); !ok {
			return p.syntaxErrorAtCur()
		}
		dt, err := p.parseDataType()
		if err != nil {
			return err
		}
		if dt == nil {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwCOLLATE {
			p.advance() // consume COLLATE
			if _, err := p.parseCollationName(); err != nil {
				return err
			}
		}
		if p.cur.Type == kwNOT {
			p.advance() // consume NOT
			if _, err := p.expect(kwNULL); err != nil {
				return err
			}
		} else if p.cur.Type == kwNULL {
			p.advance()
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, err := p.expect(')')
	return err
}

// parseExecArgList parses setParamList (TSql160.g:26422): parameters
// separated by commas, each one required. When procedure is set the list
// follows a procedure name, where SQL Server refuses a positional parameter
// after a named one (Msg 119; the action at TSql160.g:26454).
func (p *Parser) parseExecArgList(procedure bool) ([]nodes.Node, error) {
	var args []nodes.Node
	named := false
	for {
		arg, err := p.parseExecArg()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			if p.collectMode() {
				return args, nil
			}
			return nil, p.syntaxErrorAtCur()
		}
		if arg.Name != "" {
			named = true
		} else if named && procedure {
			return nil, p.newParseError(arg.Loc.Start, "parameters after a named parameter must be passed as '@name = value'")
		}
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			return args, nil
		}
	}
}

// parseExecArg parses a single EXEC argument (setParam, TSql160.g:26439).
//
//	exec_arg = [@param =] { value [OUTPUT|OUT] | DEFAULT }
//
// A value is a constant, a negative number, a variable or an identifier
// (possibleNegativeConstantOrIdentifier, TSql160.g:5186), not an
// expression: SQL Server refuses `EXEC a 1+1` and `EXEC a GETDATE()`
// (Msg 102). OUTPUT follows only a local variable (Msg 179 otherwise).
func (p *Parser) parseExecArg() (*nodes.ExecArg, error) {
	loc := p.pos()

	arg := &nodes.ExecArg{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Check for DEFAULT keyword
	if p.cur.Type == kwDEFAULT {
		arg.IsDefault = true
		p.advance()
		arg.Loc.End = p.prevEnd()
		return arg, nil
	}

	// Check for named argument: @param = value
	if p.cur.Type == tokVARIABLE {
		next := p.peekNext()
		if next.Type == '=' {
			arg.Name = p.cur.Str
			p.advance() // consume @param
			p.advance() // consume =
			// Completion: after EXEC p @param = → expression context
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				p.addRuleCandidate("func_name")
				return nil, nil
			}
			// Check for DEFAULT after @param =
			if p.cur.Type == kwDEFAULT {
				arg.IsDefault = true
				p.advance()
				arg.Loc.End = p.prevEnd()
				return arg, nil
			}
		}
	}

	value, err := p.parseExecArgValue()
	if err != nil {
		return nil, err
	}
	if value == nil {
		if arg.Name != "" {
			return nil, p.syntaxErrorAtCur()
		}
		return nil, nil
	}
	arg.Value = value

	// Check for OUTPUT/OUT
	if p.cur.Type == kwOUTPUT || p.cur.Type == kwOUT {
		if v, ok := value.(*nodes.VariableRef); !ok || strings.HasPrefix(v.Name, "@@") {
			return nil, p.newParseError(p.cur.Loc, "cannot use the OUTPUT option when passing a constant to a stored procedure")
		}
		arg.Output = true
		p.advance()
	}

	arg.Loc.End = p.prevEnd()
	return arg, nil
}

// parseExecArgValue parses possibleNegativeConstantOrIdentifier
// (TSql160.g:5186): a literal (TSql160.g:6987), a minus sign before a number
// or money constant (negativeConstant, TSql160.g:5210), or an identifier,
// which passes its name as a string. It returns nil when the current token
// starts none of them.
func (p *Parser) parseExecArgValue() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST, tokMONEY,
		kwNULL, tokVARIABLE, tokSYSVARIABLE:
		return p.parsePrimary()
	case '-':
		loc := p.pos()
		p.advance() // consume '-'
		switch p.cur.Type {
		case tokICONST, tokFCONST, tokMONEY:
			if execIsBinary(p.cur) {
				return nil, p.syntaxErrorAtCur()
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{
			Op:      nodes.UnaryMinus,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	if p.isIdentLike() {
		tok := p.advance()
		return &nodes.ColumnRef{
			Column: tok.Str,
			Loc:    nodes.Loc{Start: tok.Loc, End: p.prevEnd()},
		}, nil
	}
	return nil, nil
}

// execIsBinary reports whether a numeric token is a binary constant (0x...).
func execIsBinary(tok Token) bool {
	return len(tok.Str) > 1 && tok.Str[0] == '0' && (tok.Str[1] == 'x' || tok.Str[1] == 'X')
}

// isStatementStart returns true if the current token starts a new statement.
func (p *Parser) isStatementStart() bool {
	switch p.cur.Type {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwMERGE,
		kwCREATE, kwALTER, kwDROP, kwTRUNCATE,
		kwDECLARE, kwSET, kwIF, kwWHILE, kwBEGIN,
		kwRETURN, kwBREAK, kwCONTINUE, kwGOTO,
		kwEXEC, kwEXECUTE, kwPRINT, kwRAISERROR, kwTHROW,
		kwGRANT, kwREVOKE, kwDENY, kwUSE, kwWAITFOR,
		kwWITH, kwGO:
		return true
	}
	return false
}

package parser

import (
	"github.com/bytebase/omni/spanner/ast"
)

// ---------------------------------------------------------------------------
// BigQuery DDL — CREATE [AGGREGATE] FUNCTION / TABLE FUNCTION / PROCEDURE
// (parser-ddl-bigquery node)
// ---------------------------------------------------------------------------
//
// Ports the legacy ANTLR create_function_statement,
// create_table_function_statement, and create_procedure_statement
// (GoogleSQLParser.g4 §2.2) — a hand-port of Google's open-source ZetaSQL
// reference. These object kinds are routed here from the core parseCreateStmt
// fan-out (they were stubbed by the parser-ddl node).
//
//	create_function_statement:
//	  CREATE opt_or_replace? opt_create_scope? opt_aggregate? FUNCTION opt_if_not_exists?
//	    function_declaration opt_function_returns? opt_sql_security_clause?
//	    opt_determinism_level? opt_language_or_remote_with_connection?
//	    unordered_options_body?
//	create_table_function_statement:
//	  CREATE opt_or_replace? opt_create_scope? TABLE FUNCTION opt_if_not_exists?
//	    path_expression opt_function_parameters? opt_returns? opt_sql_security_clause?
//	    unordered_language_options? opt_as_query_or_string?
//	create_procedure_statement:
//	  CREATE opt_or_replace? opt_create_scope? PROCEDURE opt_if_not_exists?
//	    path_expression procedure_parameters opt_external_security_clause?
//	    with_connection_clause? opt_options_list? begin_end_block_or_language_as_code
//
// ORACLE NOTE — BigQuery-only at the union level (oracle.md). The Spanner
// emulator parses a NARROWER CREATE FUNCTION (a bare `CREATE FUNCTION f(x T) AS
// (expr)` accepts) but REJECTS TEMP / AGGREGATE / LANGUAGE / DETERMINISTIC /
// REMOTE and the whole TABLE FUNCTION / PROCEDURE forms (probed 2026-06-05). So
// the Spanner verdict is non-authoritative for these; they are triangulated
// against the legacy .g4 + BigQuery truth1 (DDL-015..020) and proven by the unit
// tests in bq_function_procedure_test.go.

// parseCreateFunction parses a CREATE [AGGREGATE] FUNCTION or CREATE TABLE
// FUNCTION statement. The shared CREATE prefix (OR REPLACE / scope) has been
// consumed by parseCreateStmt; aggregate reports whether an AGGREGATE qualifier
// preceded FUNCTION, and isTableFunc whether this is the TABLE FUNCTION form
// (cur is at FUNCTION in both cases). create is the CREATE token (for Loc).
func (p *Parser) parseCreateFunction(create Token, orReplace bool, scope string, aggregate, isTableFunc bool) (ast.Node, error) {
	if scope != "" || aggregate || isTableFunc {
		return nil, p.syntaxErrorAtCur()
	}
	if isTableFunc {
		// TABLE FUNCTION: cur is at TABLE (the dispatch matched kwTABLE and peeked
		// FUNCTION). Consume both.
		p.advance() // TABLE
	}
	p.advance() // consume FUNCTION

	stmt := &ast.CreateFunctionStmt{
		OrReplace:   orReplace,
		Scope:       scope,
		Aggregate:   aggregate,
		IsTableFunc: isTableFunc,
	}
	stmt.Loc.Start = create.Loc.Start

	// opt_if_not_exists? — note the legacy create_function_statement marks this as
	// REQUIRED in the rule text (opt_if_not_exists without `?` on a rule that is
	// itself optional), but it is genuinely optional in practice; parseIfNotExists
	// returns false when absent.
	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// function_declaration: path_expression function_parameters.
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// function_parameters / opt_function_parameters — a parenthesized parameter
	// list. For a scalar/aggregate FUNCTION the list is required (function_declaration
	// requires it); for a TABLE FUNCTION it is opt_function_parameters? but the
	// grammar's action emits "Expected (" when missing — so both effectively require
	// the '('. We require it for both (a missing '(' is a syntax error either way).
	if p.cur.Type != int('(') {
		return nil, p.syntaxErrorAtCur()
	}
	params, err := p.parseFunctionParameters(false /*procedure*/)
	if err != nil {
		return nil, err
	}
	stmt.Params = params

	// opt_returns? — RETURNS type_or_tvf_schema. For a TVF this is RETURNS TABLE<…>.
	// Spanner scalar SQL UDFs require RETURNS <type>.
	if p.cur.Type == kwRETURNS {
		if err := p.parseFunctionReturns(stmt); err != nil {
			return nil, err
		}
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	// opt_sql_security_clause? — SQL SECURITY {INVOKER|DEFINER}.
	if p.cur.Type == kwSQL {
		sec, err := p.parseSQLSecurityClause()
		if err != nil {
			return nil, err
		}
		stmt.SQLSecurity = sec
	}

	if isTableFunc {
		// TABLE FUNCTION tail: unordered_language_options? opt_as_query_or_string?
		//   unordered_language_options: language opt_options_list? | opt_options_list language?
		if err := p.parseUnorderedLanguageOptions(stmt); err != nil {
			return nil, err
		}
		// opt_as_query_or_string?: AS query | AS string_literal.
		if p.cur.Type == kwAS {
			if err := p.parseAsQueryOrString(stmt); err != nil {
				return nil, err
			}
		}
	} else {
		// scalar / aggregate FUNCTION tail:
		//   opt_determinism_level? opt_language_or_remote_with_connection?
		//     unordered_options_body?
		// Spanner remote functions allow `[NOT DETERMINISTIC] LANGUAGE REMOTE OPTIONS (...)`.
		if p.cur.Type == kwNOT && p.peekNext().Type == kwDETERMINISTIC {
			p.advance() // NOT
			p.advance() // DETERMINISTIC
			stmt.Determinism = "NOT DETERMINISTIC"
		}
		isSpannerRemote := false
		if p.cur.Type == kwLANGUAGE && p.peekNext().Type == kwREMOTE {
			p.advance()            // LANGUAGE
			langTok := p.advance() // REMOTE
			stmt.Language = p.tokenSource(langTok)
			isSpannerRemote = true
		}
		if err := p.parseUnorderedOptionsBody(stmt); err != nil {
			return nil, err
		}
		if isSpannerRemote {
			if len(stmt.Options) == 0 || stmt.Body != nil || stmt.HasBodyString {
				return nil, p.syntaxErrorAtCur()
			}
		} else {
			if stmt.Determinism != "" || stmt.Body == nil {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	// The function body / TVF query may embed subqueries captured as RawText by the
	// frozen expressions node; fill them so the query-span consumer sees a complete
	// tree (mirrors parseStmtWithSubqueries; idempotent).
	p.fillSubqueries(stmt)
	return stmt, nil
}

// parseFunctionReturns parses a RETURNS type_or_tvf_schema into stmt. cur is at
// RETURNS. A TVF schema (TABLE < col, … >) sets ReturnsTable + ReturnColumns; a
// plain type sets Returns.
func (p *Parser) parseFunctionReturns(stmt *ast.CreateFunctionStmt) error {
	p.advance() // RETURNS
	if p.cur.Type == kwTABLE && p.peekNext().Type == int('<') {
		return p.syntaxErrorAtCur()
	}
	dt, err := p.parseType()
	if err != nil {
		return err
	}
	stmt.Returns = &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}
	return nil
}

// parseTVFSchema parses a tvf_schema: `TABLE < tvf_schema_column (, …) >`, where
// each tvf_schema_column is `identifier type | type`. cur is at TABLE; the next
// token is the template-open '<'. Returns the columns as ColumnDefs (Name "" for
// a bare-type column).
func (p *Parser) parseTVFSchema() ([]*ast.ColumnDef, error) {
	p.advance() // TABLE
	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}
	var cols []*ast.ColumnDef
	for {
		col, err := p.parseTVFSchemaColumn()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expectTemplateClose(); err != nil {
		return nil, err
	}
	return cols, nil
}

// beginsTypeOrTVFFollowingName reports whether tokenType can begin a
// type_or_tvf_schema when it follows a candidate parameter/column NAME. It
// extends beginsTypeFollowingName (the struct-field disambiguator) with the two
// type-starts that only appear in the type_or_tvf_schema position: ANY (a
// templated_parameter_type `ANY <kind>`) and TABLE (a tvf_schema `TABLE<…>`).
// Both are needed so `p ANY TYPE` / `t TABLE<x INT64>` are read as named
// parameters, not a bare type with trailing junk.
func beginsTypeOrTVFFollowingName(tokenType int) bool {
	if tokenType == kwANY || tokenType == kwTABLE {
		return true
	}
	return beginsTypeFollowingName(tokenType)
}

// parseTVFSchemaColumn parses one tvf_schema_column: `identifier type | type`.
// It is `identifier type` when the current token is an identifier-start AND the
// following token begins a type (so a bare type name like `INT64` is not mistaken
// for a column name with a missing type).
func (p *Parser) parseTVFSchemaColumn() (*ast.ColumnDef, error) {
	start := p.cur.Loc
	if isIdentifierStart(p.cur.Type) && beginsTypeOrTVFFollowingName(p.peekNext().Type) {
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		dt, err := p.parseType()
		if err != nil {
			return nil, err
		}
		return &ast.ColumnDef{
			Name: name,
			Type: &ast.TypeRef{Text: dt.String(), Loc: dt.Loc},
			Loc:  ast.Loc{Start: start.Start, End: p.prev.Loc.End},
		}, nil
	}
	// Bare type column (no name).
	dt, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &ast.ColumnDef{
		Type: &ast.TypeRef{Text: dt.String(), Loc: dt.Loc},
		Loc:  ast.Loc{Start: start.Start, End: p.prev.Loc.End},
	}, nil
}

// parseFunctionParameters parses a parenthesized parameter list. cur is at the
// opening '('. When isProcedure is true the procedure_parameter grammar applies
// (an optional IN/OUT/INOUT mode then `identifier type_or_tvf_schema`); otherwise
// the function_parameter grammar applies (`identifier type_or_tvf_schema [AS
// alias] [DEFAULT expr] [NOT AGGREGATE]` or a bare `type_or_tvf_schema`). An empty
// `()` list is legal.
func (p *Parser) parseFunctionParameters(isProcedure bool) ([]*ast.FunctionParam, error) {
	p.advance() // '('
	var params []*ast.FunctionParam
	if p.cur.Type == int(')') {
		p.advance()
		return params, nil
	}
	for {
		var (
			param *ast.FunctionParam
			err   error
		)
		if isProcedure {
			param, err = p.parseProcedureParameter()
		} else {
			param, err = p.parseFunctionParameter()
		}
		if err != nil {
			return nil, err
		}
		params = append(params, param)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return params, nil
}

// parseFunctionParameter parses one function_parameter:
//
//	identifier type_or_tvf_schema [AS alias] [DEFAULT expr] [NOT AGGREGATE]
//	| type_or_tvf_schema [AS alias] [NOT AGGREGATE]
//
// The named form applies when cur is an identifier-start AND the following token
// begins a type; otherwise the bare-type form applies (so `INT64` is a bare type,
// not a name with a missing type).
func (p *Parser) parseFunctionParameter() (*ast.FunctionParam, error) {
	start := p.cur.Loc
	param := &ast.FunctionParam{}

	if isIdentifierStart(p.cur.Type) && beginsTypeOrTVFFollowingName(p.peekNext().Type) {
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		param.Name = name
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	// type_or_tvf_schema: a TVF `TABLE<…>` schema, ANY TYPE/<kind>, or a plain type.
	tref, err := p.parseTypeOrTVFSchema()
	if err != nil {
		return nil, err
	}
	param.Type = tref

	// opt_as_alias_with_required_as? — AS identifier.
	if p.cur.Type == kwAS {
		return nil, p.syntaxErrorAtCur()
	}

	// opt_default_expression? — DEFAULT expr (named form only in the grammar; the
	// bare-type form has no DEFAULT, but parsing it here only affects an already-
	// invalid bare-type-with-DEFAULT input, which is rare).
	if param.Name != "" && p.cur.Type == kwDEFAULT {
		p.advance() // DEFAULT
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		param.Default = expr
	}

	// opt_not_aggregate? — NOT AGGREGATE.
	if p.cur.Type == kwNOT && p.peekNext().Type == kwAGGREGATE {
		return nil, p.syntaxErrorAtCur()
	}

	param.Loc = ast.Loc{Start: start.Start, End: p.prev.Loc.End}
	return param, nil
}

// parseProcedureParameter parses one procedure_parameter:
//
//	[IN|OUT|INOUT] identifier type_or_tvf_schema
//
// The error alternative (a parameter with no type) surfaces naturally as a
// syntax error from parseTypeOrTVFSchema at the ')' or ',' terminator.
func (p *Parser) parseProcedureParameter() (*ast.FunctionParam, error) {
	start := p.cur.Loc
	param := &ast.FunctionParam{}

	// opt_procedure_parameter_mode?: IN | OUT | INOUT. These are non-reserved and
	// may legitimately be a parameter NAME — disambiguate by lookahead: a mode is a
	// mode only when an identifier (the real name) follows it; `IN INT64` means the
	// parameter named `IN` of type `INT64`.
	switch p.cur.Type {
	case kwIN:
		if isIdentifierStart(p.peekNext().Type) {
			p.advance()
			param.Mode = ast.ParamModeIn
		}
	case kwOUT:
		if isIdentifierStart(p.peekNext().Type) {
			p.advance()
			param.Mode = ast.ParamModeOut
		}
	case kwINOUT:
		if isIdentifierStart(p.peekNext().Type) {
			p.advance()
			param.Mode = ast.ParamModeInout
		}
	}

	// identifier (the parameter name) — required.
	nameTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	param.Name = name

	// type_or_tvf_schema — required.
	tref, err := p.parseTypeOrTVFSchema()
	if err != nil {
		return nil, err
	}
	param.Type = tref

	param.Loc = ast.Loc{Start: start.Start, End: p.prev.Loc.End}
	return param, nil
}

// parseTypeOrTVFSchema parses a type_or_tvf_schema (type | templated_parameter_type
// | tvf_schema) and returns it as a rendered TypeRef. A `TABLE<…>` tvf_schema and
// an `ANY <kind>` templated type are rendered to text via parseType (which handles
// both — ANY is a type_name path and TABLE< is the function/array template path),
// so this is a thin wrapper that records the source span. The bare ANY TYPE form
// (templated_parameter_kind) is handled here because parseType treats ANY as a
// plain type name.
func (p *Parser) parseTypeOrTVFSchema() (*ast.TypeRef, error) {
	start := p.cur.Loc

	// tvf_schema: TABLE < … >. parseType does not model a multi-column TABLE<…>
	// (its FUNCTION< / ARRAY< template paths are different), so handle it here and
	// render to text.
	if p.cur.Type == kwTABLE && p.peekNext().Type == int('<') {
		cols, err := p.parseTVFSchema()
		if err != nil {
			return nil, err
		}
		return &ast.TypeRef{Text: renderTVFSchema(cols), Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil
	}

	// templated_parameter_type: ANY {PROTO|ENUM|STRUCT|ARRAY|identifier}.
	if p.cur.Type == kwANY {
		anyTok := p.advance() // ANY
		kindTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		kind := p.tokenSource(kindTok)
		return &ast.TypeRef{
			Text: "ANY " + kind,
			Loc:  ast.Loc{Start: anyTok.Loc.Start, End: kindTok.Loc.End},
		}, nil
	}

	dt, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}, nil
}

// renderTVFSchema renders a TVF schema column list back to `TABLE<col TYPE, …>`
// text for the TypeRef. (The structured columns are kept on the statement node;
// this is only for the parameter-position TypeRef text.)
func renderTVFSchema(cols []*ast.ColumnDef) string {
	out := "TABLE<"
	for i, c := range cols {
		if i > 0 {
			out += ", "
		}
		if c.Name != "" {
			out += c.Name + " "
		}
		if c.Type != nil {
			out += c.Type.Text
		}
	}
	out += ">"
	return out
}

// parseConnectionClause parses a with_connection_clause `WITH CONNECTION
// connection_clause` and returns the connection name text. connection_clause is a
// path_expression or the DEFAULT keyword (BigQuery default connection). cur is at
// WITH.
func (p *Parser) parseConnectionClause() (string, error) {
	p.advance() // WITH
	if _, err := p.expect(kwCONNECTION); err != nil {
		return "", err
	}
	// connection_clause: DEFAULT | path_expression. A dotted path_expression, not
	// maybe_dashed_path_expression: a dashed project name must be backtick-quoted
	// (`my-project.us.c`) like any other identifier.
	if p.cur.Type == kwDEFAULT {
		p.advance()
		return "DEFAULT", nil
	}
	path, err := p.parsePathExpr()
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

// parseUnorderedOptionsBody parses unordered_options_body into stmt:
//
//	opt_options_list as_sql_function_body_or_string?
//	| as_sql_function_body_or_string opt_options_list?
//
// where as_sql_function_body_or_string is `AS sql_function_body | AS
// string_literal`, and sql_function_body is `( expression )` (the `( SELECT … )`
// alternative is the legacy error path — a query body must be a scalar subquery,
// so a leading SELECT inside the parens is rejected).
func (p *Parser) parseUnorderedOptionsBody(stmt *ast.CreateFunctionStmt) error {
	// Leading OPTIONS, then an optional AS body.
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseDDLOptionsList()
		if err != nil {
			return err
		}
		stmt.Options = opts
		if p.cur.Type == kwAS {
			return p.parseSQLFunctionBodyOrString(stmt)
		}
		return nil
	}
	// Leading AS body, then an optional OPTIONS.
	if p.cur.Type == kwAS {
		if err := p.parseSQLFunctionBodyOrString(stmt); err != nil {
			return err
		}
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseDDLOptionsList()
			if err != nil {
				return err
			}
			stmt.Options = opts
		}
	}
	return nil
}

// parseSQLFunctionBodyOrString parses an as_sql_function_body_or_string into stmt.
// cur is at AS. It is `AS ( expression )` (sql_function_body) or `AS string`.
func (p *Parser) parseSQLFunctionBodyOrString(stmt *ast.CreateFunctionStmt) error {
	p.advance() // AS
	if p.cur.Type == int('(') {
		p.advance() // '('
		// sql_function_body error alt: `( SELECT …` is rejected (must be a scalar
		// subquery `( ( SELECT … ) )`). This reproduces the legacy ZetaSQL
		// diagnostic.
		if p.cur.Type == kwSELECT {
			return &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "The body of each CREATE FUNCTION statement is an expression, not a query; to use a query as an expression, the query must be wrapped with additional parentheses to make it a scalar subquery expression",
			}
		}
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
		stmt.Body = expr
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseUnorderedLanguageOptions parses unordered_language_options into stmt (TVF
// form):
//
//	language opt_options_list?
//	| opt_options_list language?
//
// where language is `LANGUAGE identifier`.
func (p *Parser) parseUnorderedLanguageOptions(stmt *ast.CreateFunctionStmt) error {
	if p.cur.Type == kwLANGUAGE {
		p.advance() // LANGUAGE
		langTok, err := p.expectIdentifier()
		if err != nil {
			return err
		}
		stmt.Language = p.tokenSource(langTok)
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseOptionsList()
			if err != nil {
				return err
			}
			stmt.Options = opts
		}
		return nil
	}
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return err
		}
		stmt.Options = opts
		if p.cur.Type == kwLANGUAGE {
			p.advance() // LANGUAGE
			langTok, err := p.expectIdentifier()
			if err != nil {
				return err
			}
			stmt.Language = p.tokenSource(langTok)
		}
	}
	return nil
}

// parseAsQueryOrString parses an opt_as_query_or_string into stmt (TVF form):
// `AS query | AS string_literal`. cur is at AS.
func (p *Parser) parseAsQueryOrString(stmt *ast.CreateFunctionStmt) error {
	p.advance() // AS
	if p.cur.Type == tokString {
		strTok := p.advance()
		stmt.BodyString = strTok.Str
		stmt.HasBodyString = true
		return nil
	}
	q, err := p.parseQuery()
	if err != nil {
		return err
	}
	stmt.AsQuery = q
	return nil
}

// ---------------------------------------------------------------------------
// DROP FUNCTION / TABLE FUNCTION / PROCEDURE
// ---------------------------------------------------------------------------

// parseBQDropFunction parses a DROP FUNCTION or DROP TABLE FUNCTION statement. The
// DROP keyword has been consumed by parseDropStmt; cur is at FUNCTION (DROP
// FUNCTION) or TABLE (DROP TABLE FUNCTION). Grammar (drop_statement
// table_or_table_function alt): `DROP table_or_table_function opt_if_exists?
// maybe_dashed_path opt_function_parameters?`, where table_or_table_function is
// `TABLE FUNCTION?`. The plain `DROP TABLE` is owned by the core node; this is
// reached only for FUNCTION / TABLE FUNCTION.
func (p *Parser) parseBQDropFunction(drop Token) (ast.Node, error) {
	kind := ast.BQDropFunction
	if p.cur.Type == kwTABLE {
		p.advance() // TABLE
		kind = ast.BQDropTableFunction
	}
	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}
	return p.finishBQDropObject(drop, kind, true /*allowFunctionParams*/)
}

// finishBQDropObject parses the shared tail `opt_if_exists? maybe_dashed_path
// [opt_function_parameters]` for a simple BigQuery DROP of an object whose kind is
// already determined and whose object keyword has been consumed. When
// allowFunctionParams is true, a trailing `( … )` parameter list (used by the
// FUNCTION / TABLE FUNCTION / PROCEDURE drop alternatives to disambiguate
// overloads) is consumed structurally and discarded.
func (p *Parser) finishBQDropObject(drop Token, kind ast.BQDropObjectKind, allowFunctionParams bool) (ast.Node, error) {
	stmt := &ast.BQDropStmt{Object: kind}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// opt_function_parameters? — a parenthesized parameter list disambiguating an
	// overloaded function/procedure. Consumed structurally and discarded.
	if allowFunctionParams && p.cur.Type == int('(') {
		if err := p.skipBalancedParens(); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

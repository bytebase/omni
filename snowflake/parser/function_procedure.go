package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// Routine DDL — CREATE / ALTER FUNCTION & PROCEDURE (T4.5)
// ---------------------------------------------------------------------------
//
// A UDF / UDTF / stored procedure / external function carries a large,
// version-growing vocabulary of property clauses (LANGUAGE, RUNTIME_VERSION,
// HANDLER, PACKAGES, IMPORTS, TARGET_PATH, EXTERNAL_ACCESS_INTEGRATIONS,
// SECRETS, COMMENT, the volatility / null-handling modifiers, and the
// external-function cloud params API_INTEGRATION / HEADERS / CONTEXT_HEADERS /
// MAX_BATCH_ROWS / COMPRESSION / REQUEST_TRANSLATOR / RESPONSE_TRANSLATOR).
// Rather than mirror the legacy ANTLR grammar's finite, already-stale
// enumerations (its create_function rule lacks TARGET_PATH,
// EXTERNAL_ACCESS_INTEGRATIONS, SECRETS, the SCALA runtime, ... all in the docs
// corpus), every property is captured as an open-ended `KEY = <value>` pair or a
// bare modifier word (ast.CopyOption), reusing the merged STAGE (T4.1) / FILE
// FORMAT (T4.2) / PIPE (T4.3) / COPY (T5.2) machinery. Only the structural
// anchors are modeled explicitly: the argument list, the RETURNS clause (a
// scalar data type or a TABLE (...) column list), the EXECUTE AS {CALLER|OWNER}
// clause, and the AS <body> clause.
//
// Body handling is the delicate part. The body after AS is either single-quoted
// ('...') or dollar-quoted ($$...$$). A multi-line single-quoted body — very
// common for SQL / JavaScript / Python handlers — cannot be taken from the lexer
// token stream: omni's scanString does NOT span newlines (it emits an
// "unterminated string literal" on the first embedded '\n'). parseRoutineBody
// therefore RAW-SCANS the body directly from the source text, mirroring the
// p.base / p.input raw-source-slice + lexer re-sync pattern copy.go established
// for unquoted file:// paths, and returns the verbatim body (delimiters
// included). The body is OPAQUE — its contents (SQL-scripting / JS / Python /
// Java / Scala) are never parsed.

// ---------------------------------------------------------------------------
// CREATE FUNCTION / EXTERNAL FUNCTION / PROCEDURE
// ---------------------------------------------------------------------------

// parseCreateFunctionStmt parses CREATE [OR REPLACE] [SECURE] [TEMP|TEMPORARY]
// FUNCTION ... . The CREATE keyword and the OR REPLACE / SECURE / TEMPORARY
// modifiers have already been consumed by parseCreateStmt; start is the Loc of
// the CREATE token and cur is the FUNCTION keyword.
func (p *Parser) parseCreateFunctionStmt(start ast.Loc, orReplace, secure, temporary bool) (ast.Node, error) {
	p.advance() // consume FUNCTION
	return p.parseCreateRoutineBody(start, ast.RoutineFunction, orReplace, secure, temporary)
}

// parseCreateExternalFunctionStmt parses CREATE [OR REPLACE] [SECURE] EXTERNAL
// FUNCTION ... . The CREATE / OR REPLACE / SECURE modifiers and the EXTERNAL
// keyword have already been consumed by parseCreateStmt; start is the Loc of the
// CREATE token and cur is the FUNCTION keyword. An external function is exactly
// a function whose API_INTEGRATION property and quoted-URL body the
// catalog/semantic layer can detect; structurally it parses through the same
// path (its API_INTEGRATION / HEADERS / CONTEXT_HEADERS / ... clauses are
// captured open-ended like any other property, and its `AS '<url>'` body is the
// single-quoted form).
func (p *Parser) parseCreateExternalFunctionStmt(start ast.Loc, orReplace, secure bool) (ast.Node, error) {
	p.advance() // consume FUNCTION
	return p.parseCreateRoutineBody(start, ast.RoutineExternalFunction, orReplace, secure, false)
}

// parseCreateProcedureStmt parses CREATE [OR REPLACE] [SECURE] PROCEDURE ... .
// The CREATE / OR REPLACE / SECURE modifiers have already been consumed by
// parseCreateStmt; start is the Loc of the CREATE token and cur is the PROCEDURE
// keyword.
func (p *Parser) parseCreateProcedureStmt(start ast.Loc, orReplace, secure bool) (ast.Node, error) {
	p.advance() // consume PROCEDURE
	return p.parseCreateRoutineBody(start, ast.RoutineProcedure, orReplace, secure, false)
}

// parseCreateRoutineBody parses the shared tail common to CREATE FUNCTION /
// EXTERNAL FUNCTION / PROCEDURE, after the object-type keyword has been consumed:
//
//	[ IF NOT EXISTS ] <name> ( [ <arg> [, ...] ] )
//	  RETURNS { <data_type> | TABLE ( <col> <type> [, ...] ) } [ [ NOT ] NULL ]
//	  [ <property clauses> ] [ EXECUTE AS { CALLER | OWNER } ] [ <property clauses> ]
//	  AS { '<body>' | $$<body>$$ }
//
// IF NOT EXISTS is documented for FUNCTION only; it is accepted uniformly and
// harmlessly for all three (an unsupported IF NOT EXISTS on a procedure is a
// semantic, not syntactic, concern). The body is optional: the docs permit a
// function defined entirely by its handler (no AS clause), in which case Body
// stays "".
func (p *Parser) parseCreateRoutineBody(start ast.Loc, kind ast.RoutineKind, orReplace, secure, temporary bool) (ast.Node, error) {
	stmt := &ast.CreateRoutineStmt{
		Kind:      kind,
		OrReplace: orReplace,
		Secure:    secure,
		Temporary: temporary,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Argument list: ( [ <arg_name> <data_type> [ DEFAULT <expr> ] [, ...] ] ).
	args, err := p.parseRoutineArgList()
	if err != nil {
		return nil, err
	}
	stmt.Args = args

	// RETURNS { <data_type> | TABLE ( ... ) } [ [ NOT ] NULL ].
	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, err
	}
	if err := p.parseRoutineReturns(stmt); err != nil {
		return nil, err
	}

	// Property clauses, the optional EXECUTE AS clause, and the AS body may
	// appear in any documented order; AS is the terminal anchor. The loop reads
	// open-ended property clauses (KEY = value pairs and bare modifier words such
	// as VOLATILE / IMMUTABLE / STRICT / MEMOIZABLE / SECURE and the
	// CALLED ON NULL INPUT / RETURNS NULL ON NULL INPUT phrases), stopping at AS
	// or at EXECUTE (which introduces EXECUTE AS, not a property).
	for {
		if p.cur.Type == kwAS {
			break
		}
		if p.curIsWord("EXECUTE") {
			if err := p.parseRoutineExecuteAs(&stmt.ExecuteAs); err != nil {
				return nil, err
			}
			continue
		}
		if p.startsCopyOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		break
	}

	// Optional AS <body>. A function may have no body (handler-only); when AS is
	// present the body is mandatory.
	if p.cur.Type == kwAS {
		body, dollar, end, err := p.parseRoutineBody()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
		stmt.BodyDollar = dollar
		stmt.Loc.End = end
		return stmt, nil
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRoutineArgList parses the parenthesized argument list:
//
//	( [ <arg_name> <data_type> [ DEFAULT <expr> ] [ , ... ] ] )
//
// An empty list — () — yields nil args. cur must be '('.
func (p *Parser) parseRoutineArgList() ([]*ast.RoutineArg, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Empty arg list.
	if p.cur.Type == ')' {
		p.advance() // consume ')'
		return nil, nil
	}

	var args []*ast.RoutineArg
	for {
		arg, err := p.parseRoutineArg()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return args, nil
}

// parseRoutineArg parses one argument: <arg_name> <data_type> [ DEFAULT <expr> ].
func (p *Parser) parseRoutineArg() (*ast.RoutineArg, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	typ, err := p.parseDataType()
	if err != nil {
		return nil, err
	}

	arg := &ast.RoutineArg{
		Name: name,
		Type: typ,
		Loc:  ast.Loc{Start: name.Loc.Start, End: typ.Loc.End},
	}

	// Optional DEFAULT <expr>.
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		def, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		arg.Default = def
		arg.Loc.End = p.prev.Loc.End
	}

	return arg, nil
}

// parseRoutineReturns parses the RETURNS clause body (RETURNS already consumed):
//
//	<data_type> [ [ NOT ] NULL ]
//	TABLE ( <column_name> <data_type> [ , ... ] )
//
// The TABLE form (a UDTF / table stored procedure) records its columns in
// ReturnTable; the scalar form records its type in ReturnType. A trailing
// [ NOT ] NULL nullability modifier is accepted on the scalar form (the docs and
// the corpus show `RETURNS FLOAT NOT NULL`).
func (p *Parser) parseRoutineReturns(stmt *ast.CreateRoutineStmt) error {
	if p.cur.Type == kwTABLE {
		p.advance() // consume TABLE
		cols, err := p.parseRoutineTableColumns()
		if err != nil {
			return err
		}
		// Non-nil (possibly empty) slice records the TABLE form explicitly.
		if cols == nil {
			cols = []*ast.RoutineTableColumn{}
		}
		stmt.ReturnTable = cols
		return nil
	}

	typ, err := p.parseDataType()
	if err != nil {
		return err
	}
	stmt.ReturnType = typ

	// Optional [ NOT ] NULL nullability modifier.
	switch p.cur.Type {
	case kwNOT:
		if p.peekNext().Type == kwNULL {
			p.advance() // consume NOT
			p.advance() // consume NULL
			stmt.ReturnNotNull = true
		}
	case kwNULL:
		p.advance() // consume NULL
		stmt.ReturnNull = true
	}
	return nil
}

// parseRoutineTableColumns parses the RETURNS TABLE ( ... ) column list:
//
//	( [ <column_name> <data_type> [ , ... ] ] )
//
// An empty TABLE () yields a non-nil empty slice (recorded by the caller). cur
// must be '('.
func (p *Parser) parseRoutineTableColumns() ([]*ast.RoutineTableColumn, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if p.cur.Type == ')' {
		p.advance() // consume ')'
		return []*ast.RoutineTableColumn{}, nil
	}

	var cols []*ast.RoutineTableColumn
	for {
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		typ, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		cols = append(cols, &ast.RoutineTableColumn{
			Name: name,
			Type: typ,
			Loc:  ast.Loc{Start: name.Loc.Start, End: typ.Loc.End},
		})

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseRoutineExecuteAs parses an EXECUTE AS { CALLER | OWNER } clause. cur is
// the EXECUTE word (an identifier — EXECUTE is not an omni keyword). EXECUTE AS
// is documented for procedures; it is parsed uniformly so the option loop never
// mistakes the EXECUTE word for a bare property and then swallows the body AS.
func (p *Parser) parseRoutineExecuteAs(dst *ast.ExecuteAs) error {
	p.advance() // consume EXECUTE
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	switch p.cur.Type {
	case kwCALLER:
		p.advance()
		*dst = ast.ExecuteAsCaller
	case kwOWNER:
		p.advance()
		*dst = ast.ExecuteAsOwner
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseRoutineBody parses the AS <body> clause. On entry cur is the AS keyword
// (NOT yet consumed). The body is one of:
//
//	'<...>'            single-quoted (may span newlines; '' is an escaped quote)
//	$$<...>$$          dollar-quoted (always spans newlines)
//
// and is returned VERBATIM, delimiters included, together with whether the
// dollar-quoted form was used and the absolute end offset of the body.
//
// The body is read by RAW-SCANNING the source text rather than from the lexer
// token stream. This is mandatory for the single-quoted form: omni's scanString
// rejects an embedded newline ("unterminated string literal"), so a multi-line
// single-quoted body (a SQL / JS / Python handler) never lexes to a usable
// token. To stay uniform — and to never depend on the lexer for an opaque body —
// the dollar-quoted form is raw-scanned too. After the body, the lexer is
// re-synced past it (pos set to the byte after the closing delimiter, the
// buffered lookahead cleared) and cur is re-primed; any spurious lex error the
// (now-discarded) token scan recorded for a multi-line single-quoted body is
// truncated away.
func (p *Parser) parseRoutineBody() (body string, dollar bool, end int, err error) {
	// Snapshot the lexer's error count BEFORE consuming AS's following token, so
	// a spurious "unterminated string literal" from a multi-line single-quoted
	// body (which the about-to-run scan triggers) can be dropped.
	errLenBefore := len(p.lexer.errors)

	asTok := p.advance() // consume AS; cur is now the body's first token
	_ = asTok

	if p.cur.Type == tokEOF {
		// AS with no body is a syntax error.
		return "", false, 0, p.syntaxErrorAtCur()
	}

	// Local (input-relative) scan cursor at the body's first byte. Token Locs are
	// absolute (input-relative + base); p.input is the input-relative segment.
	startLocal := p.cur.Loc.Start - p.base
	if startLocal < 0 || startLocal >= len(p.input) {
		return "", false, 0, p.syntaxErrorAtCur()
	}

	var endLocal int
	switch p.input[startLocal] {
	case '$':
		// Dollar-quoted body: $$ ... $$. Snowflake's dollar-quoted body is a
		// single $$-delimited run (no nesting, no escapes).
		if startLocal+1 >= len(p.input) || p.input[startLocal+1] != '$' {
			// A single '$' is not a body opener.
			return "", false, 0, p.syntaxErrorAtCur()
		}
		j := startLocal + 2
		for j+1 < len(p.input) && !(p.input[j] == '$' && p.input[j+1] == '$') {
			j++
		}
		if j+1 >= len(p.input) {
			// Unterminated $$ body.
			return "", false, 0, &ParseError{
				Loc: p.cur.Loc,
				Msg: "unterminated dollar-quoted routine body",
			}
		}
		endLocal = j + 2 // include the closing $$
		dollar = true

	case '\'':
		// Single-quoted body: ' ... ' where '' is an escaped quote and the body
		// may span newlines (unlike a lexer string literal). Scan to the matching
		// unescaped closing quote.
		j := startLocal + 1
		closed := false
		for j < len(p.input) {
			if p.input[j] == '\'' {
				if j+1 < len(p.input) && p.input[j+1] == '\'' {
					j += 2 // '' escaped quote
					continue
				}
				j++ // consume closing quote
				closed = true
				break
			}
			j++
		}
		if !closed {
			return "", false, 0, &ParseError{
				Loc: p.cur.Loc,
				Msg: "unterminated routine body string literal",
			}
		}
		endLocal = j
		dollar = false

	default:
		return "", false, 0, p.syntaxErrorAtCur()
	}

	body = p.input[startLocal:endLocal]
	end = endLocal + p.base

	// Drop any spurious lex error the discarded body token produced (a multi-line
	// single-quoted body makes scanString append an unterminated-string error).
	if len(p.lexer.errors) > errLenBefore {
		p.lexer.errors = p.lexer.errors[:errLenBefore]
	}

	// Re-sync the lexer past the body and re-prime cur. Because the body bypassed
	// the token stream, the buffered lookahead and current token must be
	// discarded and re-read from endLocal.
	p.lexer.pos = endLocal
	p.hasNext = false
	p.advance()

	return body, dollar, end, nil
}

// ---------------------------------------------------------------------------
// ALTER FUNCTION / PROCEDURE
// ---------------------------------------------------------------------------

// parseAlterFunctionStmt parses ALTER FUNCTION [IF EXISTS] <name> ( [argtypes] )
// <action>. The ALTER keyword has already been consumed; cur is the FUNCTION
// keyword.
func (p *Parser) parseAlterFunctionStmt() (ast.Node, error) {
	p.advance() // consume FUNCTION
	return p.parseAlterRoutineBody(false)
}

// parseAlterProcedureStmt parses ALTER PROCEDURE [IF EXISTS] <name>
// ( [argtypes] ) <action>. The ALTER keyword has already been consumed; cur is
// the PROCEDURE keyword.
func (p *Parser) parseAlterProcedureStmt() (ast.Node, error) {
	p.advance() // consume PROCEDURE
	return p.parseAlterRoutineBody(true)
}

// parseAlterRoutineBody parses the shared ALTER FUNCTION / PROCEDURE tail after
// the object-type keyword has been consumed:
//
//	[ IF EXISTS ] <name> ( [ <argtype> [ , ... ] ] ) <action>
//
//	RENAME TO <new_name>
//	SET <options>                              -- open-ended (SECURE / COMMENT / LOG_LEVEL / ...)
//	UNSET <property> [ , ... ]                  -- SECURE / COMMENT / ...
//	EXECUTE AS { CALLER | OWNER }               -- procedure only
//
// The ALTER signature carries argument TYPES only (no names). SET SECURE /
// UNSET SECURE fall out of the open-ended SET options / UNSET property list (a
// bare SECURE word), so they need no special case.
func (p *Parser) parseAlterRoutineBody(procedure bool) (ast.Node, error) {
	stmt := &ast.AlterRoutineStmt{
		Procedure: procedure,
		Loc:       ast.Loc{Start: p.prev.Loc.Start},
	}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// ( [ <argtype> [ , ... ] ] ) signature.
	argTypes, err := p.parseRoutineArgTypeList()
	if err != nil {
		return nil, err
	}
	stmt.ArgTypes = argTypes

	// Action.
	switch p.cur.Type {
	case kwRENAME:
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoutineRename
		stmt.NewName = newName

	case kwSET:
		p.advance() // consume SET
		opts, err := p.parseRequiredOptions()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoutineSet
		stmt.Options = opts

	case kwUNSET:
		p.advance() // consume UNSET
		props, err := p.parseUnsetPropertyList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoutineUnset
		stmt.UnsetProps = props

	default:
		// EXECUTE AS { CALLER | OWNER } (procedure only). EXECUTE is an identifier,
		// not a keyword, so it is matched by text.
		if p.curIsWord("EXECUTE") {
			if err := p.parseRoutineExecuteAs(&stmt.ExecuteAs); err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterRoutineExecuteAs
			break
		}
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseRoutineArgTypeList parses the ALTER signature's argument-type list:
//
//	( [ <data_type> [ , ... ] ] )
//
// Unlike the CREATE arg list this carries types only (no names). An empty
// list — () — yields nil. cur must be '('.
func (p *Parser) parseRoutineArgTypeList() ([]*ast.TypeName, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	if p.cur.Type == ')' {
		p.advance() // consume ')'
		return nil, nil
	}

	var types []*ast.TypeName
	for {
		typ, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		types = append(types, typ)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return types, nil
}

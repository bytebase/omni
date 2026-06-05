package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-utility` DAG node. It implements GoogleSQL's
// metadata / utility statements (GoogleSQLParser.g4 §2.10 + the §2.3
// rename_statement), a hand-port of Google's open-source ZetaSQL reference and
// the grammar bytebase consumes:
//
//	assert_statement:    ASSERT expression opt_description?    (opt_description: AS string_literal)
//	analyze_statement:   ANALYZE opt_options_list? table_and_column_info_list?
//	table_and_column_info: maybe_dashed_path_expression column_list?
//	describe_statement:  describe_keyword describe_info
//	describe_info:       identifier? maybe_slashed_or_dashed_path_expression opt_from_path_expression?
//	rename_statement:    RENAME identifier path_expression TO path_expression
//	call_statement:      CALL path_expression '(' (tvf_argument (',' tvf_argument)*)? ')'
//
// CORRECTNESS (correctness-protocol.md). Adjudication per oracle.md:
//   - CALL is parsed in FULL by the live Spanner emulator (it validates the arg
//     list — `CALL p(1 => 2)` and `CALL p(,)` reject), so CALL is oracle-
//     authoritative (utility_oracle_test.go).
//   - RENAME and bare ANALYZE go through the emulator's real DDL parser
//     (authoritative): `RENAME TABLE a TO b` accepts, `RENAME a TO b` rejects
//     ("Expecting 'TABLE'"), bare `ANALYZE` accepts.
//   - ASSERT / DESCRIBE are accepted by the emulator's SHALLOW recognizer:
//     leading-form ACCEPT is authoritative, but it swallows arbitrary trailing
//     tokens (it accepts `ASSERT 1 AS notstring`, `DESCRIBE @#$ !!!`), so the
//     precise grammar follows the ZetaSQL .g4 + the canonical ZetaSQL corpus
//     (assert.sql / describe.sql) — these tails are flagged divergences.
//   - ANALYZE with targets and RENAME with a non-TABLE object kind are Spanner
//     narrowings vs the union grammar (the emulator rejects them); they follow
//     the .g4 + corpus (analyze.sql) and are flagged divergences, not diffed.

// parseAssertStmt parses `ASSERT expression opt_description?`. ASSERT is the
// current token.
func (p *Parser) parseAssertStmt() (ast.Node, error) {
	assert := p.advance() // ASSERT
	stmt := &ast.AssertStmt{Loc: assert.Loc}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.Expr = expr
	stmt.Loc.End = nodeLoc(expr).End

	// opt_description: AS string_literal. The .g4 requires a STRING literal here
	// (the emulator's shallow recognizer accepts a non-string, but that is the
	// non-authoritative trailing-token behavior — see the file header).
	if p.cur.Type == kwAS {
		p.advance() // AS
		if p.cur.Type != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		desc, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		stmt.Description = desc
		stmt.Loc.End = nodeLoc(desc).End
	}
	return stmt, nil
}

// parseAnalyzeStmt parses `ANALYZE opt_options_list? table_and_column_info_list?`.
// ANALYZE is the current token.
//
// The leading `OPTIONS(` is ambiguous: it is the opt_options_list ONLY when the
// parentheses hold `key <assign-op> value` entries (or are empty, `OPTIONS()`);
// otherwise `OPTIONS` is a table-name target carrying a column_list
// (table_and_column_info — corpus analyze.sql: `ANALYZE OPTIONS(a, b, c)` and
// `ANALYZE OPTIONS(a = b) Options(a, b, c)`). analyzeLeadingIsOptionsList does
// the multi-token lookahead to decide.
func (p *Parser) parseAnalyzeStmt() (ast.Node, error) {
	analyze := p.advance() // ANALYZE
	stmt := &ast.AnalyzeStmt{Loc: analyze.Loc}
	stmt.Loc.End = p.prev.Loc.End

	// opt_options_list?
	if p.cur.Type == kwOPTIONS && p.analyzeLeadingIsOptionsList() {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		// Distinguish a present-but-empty OPTIONS() from an absent clause: keep a
		// non-nil (possibly zero-length) slice so the clause's presence round-trips.
		if opts == nil {
			opts = []*ast.OptionsEntry{}
		}
		stmt.Options = opts
		stmt.Loc.End = p.prev.Loc.End
	}

	// table_and_column_info_list?  — present iff a table path follows.
	if isIdentifierStart(p.cur.Type) {
		targets, err := p.parseTableAndColumnInfoList()
		if err != nil {
			return nil, err
		}
		stmt.Targets = targets
		stmt.Loc.End = targets[len(targets)-1].Loc.End
	}
	return stmt, nil
}

// analyzeLeadingIsOptionsList reports whether the leading `OPTIONS` token starts
// an opt_options_list (vs a table named OPTIONS). The current token is OPTIONS.
// It scans a fresh lexer from OPTIONS: it is an options-list iff `OPTIONS` is
// followed by `(` and then either `)` (empty list) or `identifier <assign-op>`
// (a real options_entry). Anything else (`OPTIONS(a, b)`, `OPTIONS a`, bare
// `OPTIONS` at end) makes OPTIONS a table target.
func (p *Parser) analyzeLeadingIsOptionsList() bool {
	startIdx := absIndex(p, p.cur.Loc.Start)
	lx := NewLexerWithOffset(p.input[startIdx:], 0)
	if lx.NextToken().Type != kwOPTIONS {
		return false
	}
	if lx.NextToken().Type != int('(') {
		return false // bare OPTIONS (no parens) is a table name
	}
	first := lx.NextToken()
	if first.Type == int(')') {
		return true // OPTIONS() — empty options-list
	}
	if !isAnyKeywordIdentifier(first.Type) {
		return false
	}
	switch lx.NextToken().Type {
	case int('='), tokPlusEqual, tokMinusEqual:
		return true // identifier <assign-op> => options_entry
	}
	return false // identifier not followed by an assign-op => column_list
}

// parseTableAndColumnInfoList parses table_and_column_info_list:
// table_and_column_info (',' table_and_column_info)*. The first token is a table
// path start. A trailing comma is a syntax error.
func (p *Parser) parseTableAndColumnInfoList() ([]*ast.TableAndColumnInfo, error) {
	var targets []*ast.TableAndColumnInfo
	for {
		target, err := p.parseTableAndColumnInfo()
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return targets, nil
}

// parseTableAndColumnInfo parses table_and_column_info:
// maybe_dashed_path_expression column_list?.
func (p *Parser) parseTableAndColumnInfo() (*ast.TableAndColumnInfo, error) {
	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	info := &ast.TableAndColumnInfo{Path: path, Loc: path.Loc}
	if p.cur.Type == int('(') {
		cols, err := p.parseColumnList()
		if err != nil {
			return nil, err
		}
		info.Columns = cols
		info.Loc.End = p.prev.Loc.End
	}
	return info, nil
}

// parseDescribeStmt parses `describe_keyword describe_info`:
//
//	{ DESCRIBE | DESC } [<object-type-identifier>] <path> [FROM <path>]
//
// The describe keyword (DESCRIBE / DESC) is the current token. describe_info's
// leading `identifier?` is the object-type word (e.g. INDEX / TYPE / `FUNCTION`);
// it is present when an identifier is followed by ANOTHER name-start (i.e. there
// is still a path after it). With only a single name token, that token is the
// path itself and there is no object-type.
//
// NOTE the path is the .g4's maybe_slashed_or_dashed_path_expression; this parser
// handles the dotted + dashed forms via parseTablePath. The third alternative,
// slashed_path_expression (`/span/nonprod/…`, an identifier run separated by '/'
// with optional embedded floating-point segments), is a DELIBERATE flagged gap:
//   - it appears in NO BigQuery or Spanner documentation (truth1) and in NO
//     describe.sql corpus case — it is a ZetaSQL-internal path form;
//   - the Spanner emulator's verdict on it is NON-authoritative anyway — DESCRIBE
//     goes through the emulator's shallow recognizer, which "accepts" any
//     trailing tokens (it accepts `DESCRIBE @#$ !!!` too), so its accept of
//     `DESCRIBE /span/foo` proves nothing about the precise grammar;
//   - implementing the full slashed_path_expression grammar (with the
//     floating-point-literal-embedded separator cases) is disproportionate for a
//     form no consumer of this parser (BigQuery / Spanner via bytebase) emits.
//
// A leading '/' therefore surfaces as a syntax error here. Recorded in the
// divergence ledger as a known truth2-only gap; promote to implemented if a real
// corpus ever needs it.
func (p *Parser) parseDescribeStmt() (ast.Node, error) {
	kw := p.advance() // DESCRIBE | DESC
	stmt := &ast.DescribeStmt{IsDesc: kw.Type == kwDESC, Loc: kw.Loc}

	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}

	// Optional leading object-type identifier: taken only if a SECOND name follows
	// (so the first is the type and the second begins the path). FROM does not
	// count as a path continuation — `DESCRIBE foo FROM s` has foo as the path.
	if isAnyKeywordIdentifier(p.peekNext().Type) && p.peekNext().Type != kwFROM {
		typeTok := p.advance()
		objType, err := p.identifierText(typeTok)
		if err != nil {
			return nil, err
		}
		stmt.ObjectType = objType
	}

	path, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Path = path
	stmt.Loc.End = path.Loc.End

	// opt_from_path_expression: FROM <path>.
	if p.cur.Type == kwFROM {
		p.advance() // FROM
		from, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.FromPath = from
		stmt.Loc.End = from.Loc.End
	}
	return stmt, nil
}

// parseRenameStmt parses `RENAME identifier path_expression TO path_expression`.
// RENAME is the current token. The object-kind identifier is required (the .g4
// allows any identifier; Spanner narrows it to TABLE — flagged divergence).
func (p *Parser) parseRenameStmt() (ast.Node, error) {
	rename := p.advance() // RENAME

	typeTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	objType, err := p.identifierText(typeTok)
	if err != nil {
		return nil, err
	}
	stmt := &ast.RenameStmt{ObjectType: objType, Loc: rename.Loc}

	from, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	stmt.From = from

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}

	to, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	stmt.To = to
	stmt.Loc.End = to.Loc.End
	return stmt, nil
}

// parseCallStmt parses `CALL path_expression '(' (tvf_argument (',' tvf_argument)*)? ')'`.
// CALL is the current token.
func (p *Parser) parseCallStmt() (ast.Node, error) {
	call := p.advance() // CALL

	proc, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	stmt := &ast.CallStmt{Proc: proc, Loc: call.Loc}

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	// Possibly-empty argument list.
	if p.cur.Type != int(')') {
		for {
			arg, err := p.parseCallArg()
			if err != nil {
				return nil, err
			}
			stmt.Args = append(stmt.Args, arg)
			if _, ok := p.match(int(',')); ok {
				continue
			}
			break
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = closeTok.Loc.End
	return stmt, nil
}

// parseCallArg parses one tvf_argument of a CALL:
//
//	expression | descriptor_argument | table_clause | model_clause
//	| connection_clause | named_argument
//
// A plain expression / named argument is returned as the expression / *NamedArg
// node directly; TABLE / MODEL / CONNECTION / DESCRIPTOR are wrapped in *CallArg.
//
// The leading TABLE/MODEL/CONNECTION/DESCRIPTOR keyword is NOT always a clause:
// these words are all non-reserved, so `model`, `table.x`, `connection`, and
// `descriptor` are also valid identifier / field-access EXPRESSIONS (oracle: the
// Spanner emulator — authoritative for CALL — accepts `CALL p(model)`,
// `CALL p(table.x)`, `CALL p(descriptor)`). The keyword opens its clause ONLY
// when the next token begins the clause body:
//   - TABLE/MODEL  → followed by a path start (`MODEL my.model`, `TABLE x`); a
//     bare keyword or a `.`-continuation (`model`, `model.col`) is an expression.
//   - CONNECTION   → followed by a path start or DEFAULT.
//   - DESCRIPTOR   → followed by '(' (`DESCRIPTOR(a, b)`); bare `descriptor` is an
//     expression (oracle: `CALL p(descriptor)` accepts, `CALL p(descriptor x)`
//     rejects on the stray `x`).
//
// (The .g4's parenthesized-clause and bare-SELECT/WITH error alternatives — which
// only emit a tailored "must not be parenthesized" / "wrap in parens" message —
// are not reproduced; a parenthesized `(SELECT …)` is a normal subquery
// expression, which parseExpr handles.)
func (p *Parser) parseCallArg() (ast.Node, error) {
	switch p.cur.Type {
	case kwTABLE:
		if isIdentifierStart(p.peekNext().Type) {
			return p.parseCallTableArg()
		}
	case kwMODEL:
		if isIdentifierStart(p.peekNext().Type) {
			return p.parseCallModelArg()
		}
	case kwCONNECTION:
		if isIdentifierStart(p.peekNext().Type) || p.peekNext().Type == kwDEFAULT {
			return p.parseCallConnectionArg()
		}
	case kwDESCRIPTOR:
		if p.peekNext().Type == int('(') {
			return p.parseCallDescriptorArg()
		}
	}

	// Named argument: identifier '=>' value. The value may itself be a lambda —
	// named_argument: identifier '=>' (expression | lambda_argument) — so
	// `f => x -> x` and `f => (x) -> x` are valid (oracle: both accept). Reuses
	// the expressions node's lambda-aware argument-value parser so CALL and TVF
	// argument lists agree. (A bare POSITIONAL lambda `CALL p(x -> x)` is NOT a
	// valid tvf_argument — the positional alternative is a plain expression — and
	// the oracle rejects it; the positional branch below correctly uses parseExpr.)
	if isIdentifierStart(p.cur.Type) && p.peekNext().Type == tokFatArrow {
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		p.advance() // '=>'
		val, err := p.parseArgValueMaybeLambda()
		if err != nil {
			return nil, err
		}
		return &ast.NamedArg{Name: name, Value: val, Loc: ast.Loc{Start: nameTok.Loc.Start, End: nodeLoc(val).End}}, nil
	}

	// Positional expression (tvf_argument's `expression` alternative — NOT a
	// lambda; a bare positional lambda is rejected by the grammar and the oracle).
	return p.parseExpr()
}

// parseCallTableArg parses table_clause as a CALL argument: `TABLE path_expression`
// (the bare `TABLE tvf_with_suffixes` form's suffixes — pivot/hint/alias — are
// not part of a CALL argument in practice and are not modeled here; a trailing
// suffix would be left for the call's ')' check). TABLE is the current token.
func (p *Parser) parseCallTableArg() (ast.Node, error) {
	tab := p.advance() // TABLE
	path, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	return &ast.CallArg{Kind: ast.CallArgTable, Path: path, Loc: ast.Loc{Start: tab.Loc.Start, End: path.Loc.End}}, nil
}

// parseCallModelArg parses model_clause: `MODEL path_expression`. MODEL is the
// current token.
func (p *Parser) parseCallModelArg() (ast.Node, error) {
	model := p.advance() // MODEL
	path, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	return &ast.CallArg{Kind: ast.CallArgModel, Path: path, Loc: ast.Loc{Start: model.Loc.Start, End: path.Loc.End}}, nil
}

// parseCallConnectionArg parses connection_clause: `CONNECTION (path_expression |
// DEFAULT)`. CONNECTION is the current token.
func (p *Parser) parseCallConnectionArg() (ast.Node, error) {
	conn := p.advance() // CONNECTION
	arg := &ast.CallArg{Kind: ast.CallArgConnection, Loc: conn.Loc}
	if _, ok := p.match(kwDEFAULT); ok {
		arg.Default = true
		arg.Loc.End = p.prev.Loc.End
		return arg, nil
	}
	path, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	arg.Path = path
	arg.Loc.End = path.Loc.End
	return arg, nil
}

// parseCallDescriptorArg parses descriptor_argument: `DESCRIPTOR '(' column (','
// column)* ')'` where each column is an identifier. DESCRIPTOR is the current
// token.
func (p *Parser) parseCallDescriptorArg() (ast.Node, error) {
	desc := p.advance() // DESCRIPTOR
	arg := &ast.CallArg{Kind: ast.CallArgDescriptor, Loc: desc.Loc}
	cols, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	arg.Columns = cols
	arg.Loc.End = p.prev.Loc.End
	return arg, nil
}

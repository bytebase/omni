package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is the `parser-routines` DAG node's statement shell (with
// routine_body.go, which holds the control-flow body language): it implements
// Trino's SQL-routine definition statements — the legacy TrinoParser.g4
// `statement` alternatives labelled createFunction and dropFunction, plus the
// `withFunction` inline-routine prefix on a top-level query (the rootQuery rule).
//
// Legacy grammar (truth2):
//
//	createFunction      : CREATE (OR REPLACE)? functionSpecification
//	dropFunction        : DROP FUNCTION (IF EXISTS)? functionDeclaration
//	withFunction        : WITH functionSpecification (, functionSpecification)*   (precedes a query)
//	functionSpecification : FUNCTION functionDeclaration returnsClause routineCharacteristic* controlStatement
//	functionDeclaration   : qualifiedName ( (parameterDeclaration (, parameterDeclaration)*)? )
//	parameterDeclaration  : identifier? type
//	returnsClause         : RETURNS type
//	routineCharacteristic :
//	      LANGUAGE identifier              # languageCharacteristic
//	    | NOT? DETERMINISTIC               # deterministicCharacteristic
//	    | RETURNS NULL ON NULL INPUT       # returnsNullOnNullInputCharacteristic
//	    | CALLED ON NULL INPUT             # calledOnNullInputCharacteristic
//	    | SECURITY (DEFINER | INVOKER)     # securityCharacteristic
//	    | COMMENT string                   # commentCharacteristic
//
// As with the other statement nodes (show.go / session.go / grant_revoke.go),
// the three top-level statement structs are PARSER-PACKAGE types that satisfy
// ast.Node via tags declared in trino/ast/nodetags.go (T_CreateFunctionStmt,
// T_DropFunctionStmt, T_WithFunctionStmt). The functionSpecification and its
// child clauses are parser-local helper structs (no ast.NodeTag) because they
// are only ever embedded inside one of those three shells. The control-flow body
// (RoutineStatement) lives in routine_body.go.
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed facts baked in
// (see oracle_routines_test.go for the differential corpus):
//
//	F1 (the whole functionSpecification parses, even though the memory connector
//	   cannot STORE a UDF). `CREATE FUNCTION f() RETURNS bigint RETURN 42` and the
//	   full-characteristic / compound-body variants are oracle-ACCEPTED: Trino's
//	   parser+analyzer succeed and only the connector rejects with NOT_SUPPORTED
//	   (a non-SYNTAX_ERROR ⇒ accepted by the differential oracle). So this node
//	   must ACCEPT every well-formed functionSpecification.
//	F2 (parameter name is optional). `CREATE FUNCTION f(bigint) …` (type only),
//	   `f(x bigint, varchar, z double)` (mixed), and named forms are all accepted
//	   — parameterDeclaration is `identifier? type`.
//	F3 (characteristics are an unordered, repeatable list). `routineCharacteristic*`
//	   accepts any order and even repeats: `… DETERMINISTIC DETERMINISTIC …` and
//	   `… COMMENT 'x' LANGUAGE SQL DETERMINISTIC …` are accepted. parseFunctionSpec
//	   loops, appending each characteristic, with no ordering/uniqueness check.
//	F4 (WITH ( property = expr, … ) is a real but legacy-absent characteristic).
//	   The Trino 481 docs list a `WITH ( property_name = expression, … )` clause
//	   in CREATE FUNCTION (used e.g. for LANGUAGE PYTHON handler properties); the
//	   oracle ACCEPTS it (`CREATE FUNCTION f() RETURNS int WITH (handler = 'x')
//	   RETURN 1`). It is ABSENT from the legacy ANTLR routineCharacteristic rule.
//	   Because Diagnose must not false-reject valid Trino 481, this node
//	   IMPLEMENTS it as a propertiesCharacteristic and FLAGS it as a confirmed
//	   docs-ahead-of-legacy divergence (see the divergence ledger / PR body). An
//	   empty `WITH ()` is rejected (at least one property required).
//	F5 (DROP FUNCTION reuses functionDeclaration). `DROP FUNCTION f(bigint)`,
//	   `DROP FUNCTION IF EXISTS my.func(integer, varchar)`, and `DROP FUNCTION f()`
//	   are accepted; the parenthesised parameter list (types, optionally named) is
//	   mandatory even when empty (it disambiguates overloads).
//	F6 (WITH FUNCTION … query is the inline-routine prefix). `WITH FUNCTION
//	   f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(10)` is accepted; several
//	   comma-separated specifications may precede the query, and a CTE `WITH` may
//	   still follow (`WITH FUNCTION … WITH t AS (…) SELECT …`). This is the
//	   rootQuery rule; parseQuery (select.go) routes a leading `WITH FUNCTION` here.

// ---------------------------------------------------------------------------
// functionSpecification and its child clauses (parser-local helpers)
// ---------------------------------------------------------------------------

// ParameterDeclaration is one `identifier? type` parameter of a function
// signature (parameterDeclaration). Name is nil for the type-only form (F2).
type ParameterDeclaration struct {
	Name *ast.Identifier // nil when only a type was given
	Type *DataType
	Loc  ast.Loc
}

// FunctionDeclaration is `qualifiedName ( (parameterDeclaration, …)? )` — the
// function signature shared by CREATE FUNCTION (via functionSpecification) and
// DROP FUNCTION (F5). The parameter parens are always present (possibly empty).
type FunctionDeclaration struct {
	Name       *ast.QualifiedName
	Parameters []ParameterDeclaration
	Loc        ast.Loc
}

// RoutineCharacteristicKind identifies which routineCharacteristic alternative a
// RoutineCharacteristic represents.
type RoutineCharacteristicKind int

const (
	// CharacteristicLanguage is `LANGUAGE identifier`.
	CharacteristicLanguage RoutineCharacteristicKind = iota
	// CharacteristicDeterministic is `DETERMINISTIC`.
	CharacteristicDeterministic
	// CharacteristicNotDeterministic is `NOT DETERMINISTIC`.
	CharacteristicNotDeterministic
	// CharacteristicReturnsNullOnNullInput is `RETURNS NULL ON NULL INPUT`.
	CharacteristicReturnsNullOnNullInput
	// CharacteristicCalledOnNullInput is `CALLED ON NULL INPUT`.
	CharacteristicCalledOnNullInput
	// CharacteristicSecurityDefiner is `SECURITY DEFINER`.
	CharacteristicSecurityDefiner
	// CharacteristicSecurityInvoker is `SECURITY INVOKER`.
	CharacteristicSecurityInvoker
	// CharacteristicComment is `COMMENT string`.
	CharacteristicComment
	// CharacteristicProperties is `WITH ( property = expr, … )` (F4 — a
	// docs-ahead-of-legacy extension flagged in the divergence ledger).
	CharacteristicProperties
)

// RoutineProperty is one `name = expression` entry of a WITH (...) properties
// characteristic (F4). Name is a single identifier (quoted allowed); Value is a
// full expression.
type RoutineProperty struct {
	Name  *ast.Identifier
	Value Expr
	Loc   ast.Loc
}

// RoutineCharacteristic is one routineCharacteristic. Only the fields relevant
// to Kind are populated: Language (CharacteristicLanguage), Comment
// (CharacteristicComment), or Properties (CharacteristicProperties). The
// remaining kinds are keyword-only and carry no payload.
type RoutineCharacteristic struct {
	Kind       RoutineCharacteristicKind
	Language   *ast.Identifier   // CharacteristicLanguage
	Comment    Expr              // CharacteristicComment (a string literal expression)
	Properties []RoutineProperty // CharacteristicProperties
	Loc        ast.Loc
}

// FunctionSpecification is
// `FUNCTION functionDeclaration returnsClause routineCharacteristic* controlStatement`
// — the function body shared by CREATE FUNCTION and WITH FUNCTION.
type FunctionSpecification struct {
	Declaration     FunctionDeclaration
	ReturnType      *DataType               // the RETURNS type
	Characteristics []RoutineCharacteristic // nil when none
	Body            RoutineStatement        // the trailing controlStatement (RETURN / BEGIN / …)
	Loc             ast.Loc
}

// ---------------------------------------------------------------------------
// CREATE FUNCTION
// ---------------------------------------------------------------------------

// CreateFunctionStmt is `CREATE [OR REPLACE] functionSpecification`
// (createFunction). OrReplace is true for the OR REPLACE form.
type CreateFunctionStmt struct {
	OrReplace bool
	Spec      FunctionSpecification
	Loc       ast.Loc
}

// Tag implements ast.Node.
func (s *CreateFunctionStmt) Tag() ast.NodeTag { return ast.T_CreateFunctionStmt }

// Span returns the source byte range.
func (s *CreateFunctionStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*CreateFunctionStmt)(nil)

// createOrReplaceFunctionFollows reports whether a `CREATE OR REPLACE …` is the
// FUNCTION variant (`CREATE OR REPLACE FUNCTION …`) rather than a parser-ddl
// OR-REPLACE object (TABLE / VIEW / MATERIALIZED VIEW). On entry cur is CREATE
// and peekNext is OR. It speculatively consumes CREATE OR REPLACE via a
// checkpoint — the parser's two-token window cannot otherwise see the keyword
// after OR REPLACE — and rewinds, leaving the cursor unchanged. A malformed
// `CREATE OR <not-REPLACE>` returns false so the dispatcher falls through (the
// real error surfaces later from whatever claims it).
func (p *Parser) createOrReplaceFunctionFollows() bool {
	cp := p.checkpoint()
	defer p.restore(cp)
	p.advance() // consume CREATE
	if p.cur.Kind != kwOR {
		return false
	}
	p.advance() // consume OR
	if p.cur.Kind != kwREPLACE {
		return false
	}
	p.advance() // consume REPLACE
	return p.cur.Kind == kwFUNCTION
}

// parseCreateFunctionStmt parses `CREATE [OR REPLACE] functionSpecification`.
// The CREATE keyword has already been consumed by the dispatcher; createStart is
// its start offset, and the current token is either OR or FUNCTION.
func (p *Parser) parseCreateFunctionStmt(createStart int) (ast.Node, error) {
	s := &CreateFunctionStmt{Loc: ast.Loc{Start: createStart}}
	if p.cur.Kind == kwOR {
		p.advance() // consume OR
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		s.OrReplace = true
	}
	spec, err := p.parseFunctionSpecification()
	if err != nil {
		return nil, err
	}
	s.Spec = spec
	s.Loc.End = spec.Loc.End
	return s, nil
}

// ---------------------------------------------------------------------------
// DROP FUNCTION
// ---------------------------------------------------------------------------

// DropFunctionStmt is `DROP FUNCTION [IF EXISTS] functionDeclaration`
// (dropFunction, F5). IfExists is true for the IF EXISTS form.
type DropFunctionStmt struct {
	IfExists    bool
	Declaration FunctionDeclaration
	Loc         ast.Loc
}

// Tag implements ast.Node.
func (s *DropFunctionStmt) Tag() ast.NodeTag { return ast.T_DropFunctionStmt }

// Span returns the source byte range.
func (s *DropFunctionStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*DropFunctionStmt)(nil)

// parseDropFunctionStmt parses `DROP FUNCTION [IF EXISTS] functionDeclaration`.
// The DROP and FUNCTION keywords have already been consumed by the dispatcher;
// dropStart is DROP's start offset and the current token is IF or the
// declaration's qualified name.
func (p *Parser) parseDropFunctionStmt(dropStart int) (ast.Node, error) {
	s := &DropFunctionStmt{Loc: ast.Loc{Start: dropStart}}
	if p.cur.Kind == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		s.IfExists = true
	}
	decl, err := p.parseFunctionDeclaration()
	if err != nil {
		return nil, err
	}
	s.Declaration = decl
	s.Loc.End = decl.Loc.End
	return s, nil
}

// ---------------------------------------------------------------------------
// WITH FUNCTION (inline routine prefix on a query)
// ---------------------------------------------------------------------------

// WithFunctionStmt is `WITH functionSpecification (, functionSpecification)* query`
// (the withFunction rule on rootQuery, F6). Functions is the non-empty list of
// inline routine definitions; Query is the trailing query that may call them.
type WithFunctionStmt struct {
	Functions []FunctionSpecification
	Query     *Query
	Loc       ast.Loc
}

// Tag implements ast.Node.
func (s *WithFunctionStmt) Tag() ast.NodeTag { return ast.T_WithFunctionStmt }

// Span returns the source byte range.
func (s *WithFunctionStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*WithFunctionStmt)(nil)

// parseWithFunctionQuery parses `WITH functionSpecification (, functionSpecification)* query`
// (withFunction). WITH is the current token and the caller (parseQuery in
// select.go) has already confirmed FUNCTION follows it. Returns a
// *WithFunctionStmt wrapping the inline routines and the trailing query.
func (p *Parser) parseWithFunctionQuery() (ast.Node, error) {
	withTok := p.advance() // consume WITH
	s := &WithFunctionStmt{Loc: ast.Loc{Start: withTok.Loc.Start}}

	first, err := p.parseFunctionSpecification()
	if err != nil {
		return nil, err
	}
	s.Functions = append(s.Functions, first)
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseFunctionSpecification()
		if err != nil {
			return nil, err
		}
		s.Functions = append(s.Functions, next)
	}

	// The trailing query (which may itself begin with a CTE `WITH`, F6).
	q, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	s.Query = q
	s.Loc.End = q.Loc.End
	return s, nil
}

// ---------------------------------------------------------------------------
// functionSpecification / functionDeclaration / characteristics
// ---------------------------------------------------------------------------

// parseFunctionSpecification parses
// `FUNCTION functionDeclaration returnsClause routineCharacteristic* controlStatement`.
// FUNCTION is the current token.
func (p *Parser) parseFunctionSpecification() (FunctionSpecification, error) {
	funcTok, err := p.expect(kwFUNCTION)
	if err != nil {
		return FunctionSpecification{}, err
	}
	spec := FunctionSpecification{Loc: ast.Loc{Start: funcTok.Loc.Start, End: funcTok.Loc.End}}

	decl, err := p.parseFunctionDeclaration()
	if err != nil {
		return FunctionSpecification{}, err
	}
	spec.Declaration = decl

	// returnsClause: RETURNS type.
	if _, err := p.expect(kwRETURNS); err != nil {
		return FunctionSpecification{}, err
	}
	retType, err := p.parseType()
	if err != nil {
		return FunctionSpecification{}, err
	}
	spec.ReturnType = retType

	// routineCharacteristic* (F3 — unordered, repeatable).
	for {
		ch, ok, err := p.tryParseRoutineCharacteristic()
		if err != nil {
			return FunctionSpecification{}, err
		}
		if !ok {
			break
		}
		spec.Characteristics = append(spec.Characteristics, ch)
	}

	// controlStatement (the trailing body — RETURN / BEGIN / IF / …).
	body, err := p.parseRoutineStatement()
	if err != nil {
		return FunctionSpecification{}, err
	}
	spec.Body = body
	spec.Loc.End = body.Span().End
	return spec, nil
}

// parseFunctionDeclaration parses
// `qualifiedName ( (parameterDeclaration (, parameterDeclaration)*)? )` (F2, F5).
// The qualified name is the current token; the parameter parens are mandatory
// (possibly empty).
func (p *Parser) parseFunctionDeclaration() (FunctionDeclaration, error) {
	name, err := p.parseQualifiedName()
	if err != nil {
		return FunctionDeclaration{}, err
	}
	decl := FunctionDeclaration{Name: name, Loc: ast.Loc{Start: name.Loc.Start, End: name.Loc.End}}

	if _, err := p.expect(int('(')); err != nil {
		return FunctionDeclaration{}, err
	}
	if p.cur.Kind != int(')') {
		first, err := p.parseParameterDeclaration()
		if err != nil {
			return FunctionDeclaration{}, err
		}
		decl.Parameters = append(decl.Parameters, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseParameterDeclaration()
			if err != nil {
				return FunctionDeclaration{}, err
			}
			decl.Parameters = append(decl.Parameters, next)
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return FunctionDeclaration{}, err
	}
	decl.Loc.End = closeTok.Loc.End
	return decl, nil
}

// parseParameterDeclaration parses one `identifier? type` parameter (F2). The
// grammar is ambiguous between an unnamed `type` and a named `identifier type`:
// a parameter name is itself an identifier (and may be a non-reserved type
// keyword — `f(comment bigint)` names a parameter `comment`), while a type can
// be multi-token (DOUBLE PRECISION, TIME WITHOUT TIME ZONE, INTERVAL DAY TO
// SECOND, t ARRAY[5], ARRAY(int)). Two tokens of lookahead cannot resolve every
// case (`f(double precision)` is one unnamed DOUBLE PRECISION type, but
// `f(x double)` names a parameter `x` of type DOUBLE — the deciding token comes
// after the second), so this mirrors parseRowField's speculative resolution.
//
// Resolution follows ANTLR's alternative order. The `identifier?` is optional,
// so the no-name reading is tried first: parse a single type and keep it if it
// lands exactly on a parameter boundary (',' or ')'). Only if that does not fit
// is the named form `identifier type` tried (consume a one-token name, then a
// full type, again requiring a clean boundary). Each attempt rewinds on failure
// via a parser/lexer checkpoint. This accepts every Trino-481 form — type-only
// (`bigint`, `double precision`, `interval day to second`, `ROW(a int)`),
// named (`x bigint`, `comment bigint`, `x double precision`), and quoted-name
// (`"my x" bigint`) — matching the oracle.
func (p *Parser) parseParameterDeclaration() (ParameterDeclaration, error) {
	start := p.cur.Loc.Start

	// Attempt 1 — unnamed `type` (the `identifier?` is absent).
	cp := p.checkpoint()
	if typ, err := p.parseType(); err == nil && p.atParameterBoundary() {
		return ParameterDeclaration{Type: typ, Loc: ast.Loc{Start: start, End: typ.Loc.End}}, nil
	}
	p.restore(cp)

	// Attempt 2 — named `identifier type`. The name is a single identifier
	// token; the type follows and must reach a parameter boundary.
	if isIdentifierStart(p.cur.Kind) {
		name := identFromToken(p.advance())
		if typ, err := p.parseType(); err == nil && p.atParameterBoundary() {
			return ParameterDeclaration{
				Name: name,
				Type: typ,
				Loc:  ast.Loc{Start: start, End: typ.Loc.End},
			}, nil
		}
		p.restore(cp)
	}

	// Neither reading fits: report the error at the parameter's first token.
	return ParameterDeclaration{}, p.syntaxErrorAtCur()
}

// atParameterBoundary reports whether the cursor sits at the end of a
// parameterDeclaration (a ',' before the next parameter or the closing ')').
func (p *Parser) atParameterBoundary() bool {
	return p.cur.Kind == int(',') || p.cur.Kind == int(')')
}

// tryParseRoutineCharacteristic parses one routineCharacteristic if the current
// token begins one, returning ok == false (and consuming nothing) otherwise so
// the caller's loop can stop and fall through to the controlStatement body.
func (p *Parser) tryParseRoutineCharacteristic() (RoutineCharacteristic, bool, error) {
	start := p.cur.Loc.Start
	switch p.cur.Kind {
	case kwLANGUAGE:
		p.advance() // consume LANGUAGE
		lang, err := p.parseIdentifier()
		if err != nil {
			return RoutineCharacteristic{}, false, err
		}
		return RoutineCharacteristic{
			Kind:     CharacteristicLanguage,
			Language: lang,
			Loc:      ast.Loc{Start: start, End: lang.Loc.End},
		}, true, nil

	case kwNOT:
		p.advance() // consume NOT
		detTok, err := p.expect(kwDETERMINISTIC)
		if err != nil {
			return RoutineCharacteristic{}, false, err
		}
		return RoutineCharacteristic{
			Kind: CharacteristicNotDeterministic,
			Loc:  ast.Loc{Start: start, End: detTok.Loc.End},
		}, true, nil

	case kwDETERMINISTIC:
		detTok := p.advance() // consume DETERMINISTIC
		return RoutineCharacteristic{
			Kind: CharacteristicDeterministic,
			Loc:  ast.Loc{Start: start, End: detTok.Loc.End},
		}, true, nil

	case kwRETURNS:
		// RETURNS NULL ON NULL INPUT. (The functionSpecification's own RETURNS
		// type clause has already been consumed before the characteristic loop,
		// so a RETURNS here is unambiguously the null-input characteristic.)
		p.advance() // consume RETURNS
		if _, err := p.expect(kwNULL); err != nil {
			return RoutineCharacteristic{}, false, err
		}
		if _, err := p.expect(kwON); err != nil {
			return RoutineCharacteristic{}, false, err
		}
		if _, err := p.expect(kwNULL); err != nil {
			return RoutineCharacteristic{}, false, err
		}
		inputTok, err := p.expect(kwINPUT)
		if err != nil {
			return RoutineCharacteristic{}, false, err
		}
		return RoutineCharacteristic{
			Kind: CharacteristicReturnsNullOnNullInput,
			Loc:  ast.Loc{Start: start, End: inputTok.Loc.End},
		}, true, nil

	case kwCALLED:
		p.advance() // consume CALLED
		if _, err := p.expect(kwON); err != nil {
			return RoutineCharacteristic{}, false, err
		}
		if _, err := p.expect(kwNULL); err != nil {
			return RoutineCharacteristic{}, false, err
		}
		inputTok, err := p.expect(kwINPUT)
		if err != nil {
			return RoutineCharacteristic{}, false, err
		}
		return RoutineCharacteristic{
			Kind: CharacteristicCalledOnNullInput,
			Loc:  ast.Loc{Start: start, End: inputTok.Loc.End},
		}, true, nil

	case kwSECURITY:
		p.advance() // consume SECURITY
		switch p.cur.Kind {
		case kwDEFINER:
			defTok := p.advance()
			return RoutineCharacteristic{
				Kind: CharacteristicSecurityDefiner,
				Loc:  ast.Loc{Start: start, End: defTok.Loc.End},
			}, true, nil
		case kwINVOKER:
			invTok := p.advance()
			return RoutineCharacteristic{
				Kind: CharacteristicSecurityInvoker,
				Loc:  ast.Loc{Start: start, End: invTok.Loc.End},
			}, true, nil
		default:
			return RoutineCharacteristic{}, false, p.syntaxErrorAtCur()
		}

	case kwCOMMENT:
		p.advance() // consume COMMENT
		// The grammar's commentCharacteristic is `COMMENT string_` — a STRING
		// literal specifically, NOT an arbitrary expression. Trino 481 rejects a
		// non-string COMMENT (`COMMENT 42`, `COMMENT foo`) with a SYNTAX_ERROR, so
		// only a basic or unicode string token is accepted here.
		if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
			return RoutineCharacteristic{}, false, p.syntaxErrorAtCur()
		}
		comment, err := p.parseStringLiteral()
		if err != nil {
			return RoutineCharacteristic{}, false, err
		}
		return RoutineCharacteristic{
			Kind:    CharacteristicComment,
			Comment: comment,
			Loc:     ast.Loc{Start: start, End: comment.Span().End},
		}, true, nil

	case kwWITH:
		// F4: WITH ( property = expr, … ) — a docs-ahead-of-legacy properties
		// characteristic (flagged divergence). At least one property is required
		// (`WITH ()` is rejected by Trino 481).
		props, end, err := p.parseRoutineProperties()
		if err != nil {
			return RoutineCharacteristic{}, false, err
		}
		return RoutineCharacteristic{
			Kind:       CharacteristicProperties,
			Properties: props,
			Loc:        ast.Loc{Start: start, End: end},
		}, true, nil

	default:
		return RoutineCharacteristic{}, false, nil
	}
}

// parseRoutineProperties parses `WITH ( property (, property)* )` where each
// property is `name = expression` (F4). WITH is the current token. Returns the
// property list and the end offset (the closing ')'). An empty list is a syntax
// error (Trino 481 rejects `WITH ()`).
func (p *Parser) parseRoutineProperties() ([]RoutineProperty, int, error) {
	p.advance() // consume WITH
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}
	first, err := p.parseRoutineProperty()
	if err != nil {
		return nil, 0, err
	}
	props := []RoutineProperty{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseRoutineProperty()
		if err != nil {
			return nil, 0, err
		}
		props = append(props, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}
	return props, closeTok.Loc.End, nil
}

// parseRoutineProperty parses one `name = expression` property entry (F4). Name
// is a single identifier (quoted allowed); value is a full expression.
func (p *Parser) parseRoutineProperty() (RoutineProperty, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return RoutineProperty{}, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return RoutineProperty{}, err
	}
	val, err := p.parseExpr()
	if err != nil {
		return RoutineProperty{}, err
	}
	return RoutineProperty{
		Name:  name,
		Value: val,
		Loc:   ast.Loc{Start: name.Loc.Start, End: val.Span().End},
	}, nil
}

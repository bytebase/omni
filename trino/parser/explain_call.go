package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the `parser-utility` DAG node (with show.go and
// session.go): it implements Trino's EXPLAIN / EXPLAIN ANALYZE and CALL
// statements — the legacy TrinoParser.g4 `statement` alternatives labelled
// explain, explainAnalyze, and call.
//
// As in show.go / session.go, the statement nodes are PARSER-PACKAGE types
// satisfying ast.Node via tags declared in trino/ast/nodetags.go.
//
// The legacy alternatives implemented here:
//
//	EXPLAIN (LPAREN explainOption (, explainOption)* RPAREN)? statement  # explain
//	EXPLAIN ANALYZE VERBOSE? statement                                   # explainAnalyze
//	CALL qualifiedName LPAREN (callArgument (, callArgument)*)? RPAREN    # call
//
//	explainOption : FORMAT (TEXT|GRAPHVIZ|JSON) | TYPE (LOGICAL|DISTRIBUTED|VALIDATE|IO) ;
//	callArgument  : expression | identifier RDOUBLEARROW expression ;
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed facts baked
// in:
//
//	E1 (EXPLAIN wraps the full statement grammar). EXPLAIN / EXPLAIN ANALYZE may
//	   prefix ANY statement (SELECT, INSERT, CREATE TABLE AS, even SHOW TABLES);
//	   all are oracle-ACCEPTED. The inner statement is parsed by recursing into
//	   the top-level dispatch (parseStmt), so EXPLAIN gains every statement form
//	   automatically as the other DAG nodes land. A stubbed inner statement
//	   (e.g. SELECT before parser-select merges) surfaces the inner node's
//	   "not yet supported" error, which is correct foundation behaviour.
//	E2 (explainOption values are a closed enum). `EXPLAIN (TYPE FOO) …` and
//	   `EXPLAIN (FORMAT BAR) …` are SYNTAX_ERRORs; the FORMAT/TYPE values are the
//	   fixed keyword sets below, not arbitrary identifiers. An empty option list
//	   `EXPLAIN () …` is also rejected (at least one option required inside the
//	   parens).
//	E3 (CALL name is at most 3 parts). `CALL c.s.p()` is accepted but
//	   `CALL c.s.p.x()` is a SYNTAX_ERROR — a procedure name is
//	   [catalog.[schema.]]procedure. parseCallStmt enforces the ≤3-part limit.
//	   The parens are mandatory (`CALL f` rejects); arguments may be empty,
//	   positional, named (`name => expr`), or mixed in any order (Trino 481
//	   accepts `CALL f(name => 1, 2)`).

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

// ExplainFormat is the FORMAT value of an EXPLAIN option.
type ExplainFormat int

const (
	// ExplainFormatUnset means no FORMAT option was given.
	ExplainFormatUnset ExplainFormat = iota
	// ExplainFormatText is FORMAT TEXT.
	ExplainFormatText
	// ExplainFormatGraphviz is FORMAT GRAPHVIZ.
	ExplainFormatGraphviz
	// ExplainFormatJSON is FORMAT JSON.
	ExplainFormatJSON
)

// ExplainType is the TYPE value of an EXPLAIN option.
type ExplainType int

const (
	// ExplainTypeUnset means no TYPE option was given.
	ExplainTypeUnset ExplainType = iota
	// ExplainTypeLogical is TYPE LOGICAL.
	ExplainTypeLogical
	// ExplainTypeDistributed is TYPE DISTRIBUTED.
	ExplainTypeDistributed
	// ExplainTypeValidate is TYPE VALIDATE.
	ExplainTypeValidate
	// ExplainTypeIO is TYPE IO.
	ExplainTypeIO
)

// ExplainStmt is EXPLAIN [ (option, …) ] statement or
// EXPLAIN ANALYZE [VERBOSE] statement. Analyze and the option fields are
// mutually exclusive in valid Trino (ANALYZE has no parenthesised options), but
// the node represents whichever the source used.
type ExplainStmt struct {
	// Analyze is true for EXPLAIN ANALYZE.
	Analyze bool
	// Verbose is true for EXPLAIN ANALYZE VERBOSE.
	Verbose bool
	// Format / Type are the parenthesised options (non-analyze form); Unset
	// when the corresponding option was absent.
	Format ExplainFormat
	Type   ExplainType
	// Statement is the explained inner statement, parsed by recursing into the
	// top-level dispatch. Nil only if the inner statement failed to parse.
	Statement ast.Node
	Loc       ast.Loc
}

// Tag implements ast.Node.
func (s *ExplainStmt) Tag() ast.NodeTag { return ast.T_ExplainStmt }

// Span returns the source byte range.
func (s *ExplainStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*ExplainStmt)(nil)

// parseExplainStmt parses EXPLAIN / EXPLAIN ANALYZE. EXPLAIN is the current
// token.
func (p *Parser) parseExplainStmt() (ast.Node, error) {
	explainTok := p.advance() // consume EXPLAIN
	s := &ExplainStmt{Loc: ast.Loc{Start: explainTok.Loc.Start, End: explainTok.Loc.End}}

	if _, ok := p.match(kwANALYZE); ok {
		s.Analyze = true
		if _, ok := p.match(kwVERBOSE); ok {
			s.Verbose = true
		}
	} else if p.cur.Kind == int('(') {
		// EXPLAIN ( option (, option)* ) — at least one option required.
		if err := p.parseExplainOptions(s); err != nil {
			return nil, err
		}
	}

	stmt, err := p.parseStmt()
	// Even when the inner statement errors — most importantly when it is a
	// still-stubbed "not yet supported" form (SELECT/INSERT/CREATE before
	// parser-select / dml / ddl land) — return the ExplainStmt node alongside
	// the error. parseSingle records the error AND keeps the node, so Diagnose
	// reports the inner gap while downstream consumers still see an EXPLAIN node
	// (and the result is classifiable as pending-inner rather than a hard
	// rejection of EXPLAIN's own grammar). Once the inner node lands the error
	// disappears and the same EXPLAIN parse becomes clean.
	s.Statement = stmt
	if sp, ok := stmt.(interface{ Span() ast.Loc }); ok && stmt != nil {
		s.Loc.End = sp.Span().End
	} else if end := p.prev.Loc.End; end > s.Loc.End {
		// Inner statement produced no measurable node (e.g. a still-stubbed form
		// whose unsupported() helper consumed to the segment boundary). Extend
		// the EXPLAIN span to the last token the inner parser consumed so the
		// node still covers the explained statement text rather than ending at
		// the EXPLAIN keyword.
		s.Loc.End = end
	}
	return s, err
}

// parseExplainOptions parses `( explainOption (, explainOption)* )` into s. The
// opening '(' is the current token. At least one option is required (an empty
// list is a syntax error, matching Trino 481).
func (p *Parser) parseExplainOptions(s *ExplainStmt) error {
	p.advance() // consume '('
	if err := p.parseExplainOption(s); err != nil {
		return err
	}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		if err := p.parseExplainOption(s); err != nil {
			return err
		}
	}
	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	return nil
}

// parseExplainOption parses one `FORMAT (TEXT|GRAPHVIZ|JSON)` or
// `TYPE (LOGICAL|DISTRIBUTED|VALIDATE|IO)` option into s. The value keyword set
// is closed (E2).
func (p *Parser) parseExplainOption(s *ExplainStmt) error {
	switch p.cur.Kind {
	case kwFORMAT:
		p.advance() // consume FORMAT
		switch p.cur.Kind {
		case kwTEXT:
			s.Format = ExplainFormatText
		case kwGRAPHVIZ:
			s.Format = ExplainFormatGraphviz
		case kwJSON:
			s.Format = ExplainFormatJSON
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume the value keyword
		return nil
	case kwTYPE:
		p.advance() // consume TYPE
		switch p.cur.Kind {
		case kwLOGICAL:
			s.Type = ExplainTypeLogical
		case kwDISTRIBUTED:
			s.Type = ExplainTypeDistributed
		case kwVALIDATE:
			s.Type = ExplainTypeValidate
		case kwIO:
			s.Type = ExplainTypeIO
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume the value keyword
		return nil
	default:
		return p.syntaxErrorAtCur()
	}
}

// ---------------------------------------------------------------------------
// CALL
// ---------------------------------------------------------------------------

// CallArgument is one argument of a CALL: positional (`expression`) or named
// (`name => expression`). Name is nil for the positional form.
type CallArgument struct {
	Name  *ast.Identifier
	Value Expr
	Loc   ast.Loc
}

// CallStmt is CALL name ( [name =>] expression, … ). Name has 1-3 parts
// ([catalog.[schema.]]procedure).
type CallStmt struct {
	Name      *ast.QualifiedName
	Arguments []CallArgument
	Loc       ast.Loc
}

// Tag implements ast.Node.
func (s *CallStmt) Tag() ast.NodeTag { return ast.T_CallStmt }

// Span returns the source byte range.
func (s *CallStmt) Span() ast.Loc { return s.Loc }

var _ ast.Node = (*CallStmt)(nil)

// parseCallStmt parses CALL qualifiedName ( (callArgument (, callArgument)*)? ).
// CALL is the current token.
func (p *Parser) parseCallStmt() (ast.Node, error) {
	callTok := p.advance() // consume CALL
	// E3: a procedure name is at most catalog.schema.procedure (3 parts);
	// `CALL c.s.p.x()` is a SYNTAX_ERROR in Trino 481. parseBoundedQualifiedName
	// (show.go) enforces the limit, matching the other name positions.
	name, err := p.parseBoundedQualifiedName(3, "procedure name")
	if err != nil {
		return nil, err
	}
	s := &CallStmt{Name: name, Loc: ast.Loc{Start: callTok.Loc.Start, End: name.Loc.End}}

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	if p.cur.Kind != int(')') {
		first, err := p.parseCallArgument()
		if err != nil {
			return nil, err
		}
		s.Arguments = append(s.Arguments, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseCallArgument()
			if err != nil {
				return nil, err
			}
			s.Arguments = append(s.Arguments, next)
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	s.Loc.End = closeTok.Loc.End
	return s, nil
}

// parseCallArgument parses one callArgument: a named `identifier => expression`
// or a positional `expression`. The named form is detected by an identifier
// followed by the => (tokDoubleArrow) token; otherwise the whole thing is a
// positional expression (which may itself begin with an identifier, e.g. a
// column or function call).
func (p *Parser) parseCallArgument() (CallArgument, error) {
	if isIdentifierStart(p.cur.Kind) && p.peekNext().Kind == tokDoubleArrow {
		nameTok := p.advance() // consume the name
		name := identFromToken(nameTok)
		p.advance() // consume =>
		value, err := p.parseExpr()
		if err != nil {
			return CallArgument{}, err
		}
		return CallArgument{
			Name:  name,
			Value: value,
			Loc:   ast.Loc{Start: name.Loc.Start, End: value.Span().End},
		}, nil
	}
	value, err := p.parseExpr()
	if err != nil {
		return CallArgument{}, err
	}
	return CallArgument{Value: value, Loc: value.Span()}, nil
}

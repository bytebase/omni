package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// parsePrimary produces a primary expression and attaches any trailing
// path steps. It's the entry point for the exprPrimary + pathStep+
// combination in the grammar (line 528):
//
//	exprPrimary pathStep+      # ExprPrimaryPath
//
// Grammar: exprPrimary (lines 514-534).
func (p *Parser) parsePrimary() (ast.ExprNode, error) {
	base, err := p.parsePrimaryBase()
	if err != nil {
		return nil, err
	}
	if !isPathStepStart(p.cur.Type) {
		return base, nil
	}
	return p.parsePathSteps(base)
}

// parsePrimaryBase dispatches on the current token to produce one of
// the 16 primary-expression alternatives. Foundation handles exprTerm
// alternatives directly; every other alternative calls deferredFeature
// to return a stub error pointing at the owning DAG node.
//
// Grammar: exprPrimary base alternatives (lines 514-534) + exprTerm
// alternatives (lines 542-549).
func (p *Parser) parsePrimaryBase() (ast.ExprNode, error) {
	switch p.cur.Type {
	// ------------------------------------------------------------------
	// Real: literal primary forms (delegates to literals.go).
	// ------------------------------------------------------------------
	case tokNULL, tokMISSING, tokTRUE, tokFALSE,
		tokSCONST, tokICONST, tokFCONST, tokION_LITERAL,
		tokDATE, tokTIME:
		return p.parseLiteral()

	// ------------------------------------------------------------------
	// Real: parenthesized expr (valueList and (SELECT ...) are stubbed
	// inside parseParenExpr).
	// ------------------------------------------------------------------
	case tokPAREN_LEFT:
		return p.parseParenExpr()

	// ------------------------------------------------------------------
	// Real: collection literals.
	// ------------------------------------------------------------------
	case tokBRACKET_LEFT:
		return p.parseArrayLit()
	case tokANGLE_DOUBLE_LEFT:
		return p.parseBagLit()
	case tokBRACE_LEFT:
		return p.parseTupleLit()

	// ------------------------------------------------------------------
	// Real: parameter and varRef.
	// ------------------------------------------------------------------
	case tokQUESTION_MARK:
		return p.parseParamRef()
	case tokAT_SIGN, tokIDENT, tokIDENT_QUOTED:
		return p.parseVarRef()

	// ------------------------------------------------------------------
	// Stub: VALUES row list (no AST node yet).
	// ------------------------------------------------------------------
	case tokVALUES:
		return nil, p.deferredFeature("VALUES", "parser-dml (DAG node 6)")

	// ------------------------------------------------------------------
	// Real: CAST family — cast / canCast / canLosslessCast (DAG node 15b).
	// ------------------------------------------------------------------
	case tokCAST:
		return p.parseCastExpr(ast.CastKindCast)
	case tokCAN_CAST:
		return p.parseCastExpr(ast.CastKindCanCast)
	case tokCAN_LOSSLESS_CAST:
		return p.parseCastExpr(ast.CastKindCanLosslessCast)

	// ------------------------------------------------------------------
	// Real: CASE expression (DAG node 15b.3).
	// ------------------------------------------------------------------
	case tokCASE:
		return p.parseCaseExpr()

	// ------------------------------------------------------------------
	// Real: keyword-bearing builtin functions (DAG node 15b).
	//
	// COALESCE/NULLIF/SUBSTRING/TRIM/EXTRACT get dedicated typed nodes
	// because their argument syntax uses keywords (AS/FROM/FOR) or a
	// fixed arity. DATE_ADD/DATE_DIFF have ordinary comma-separated
	// argument syntax (the date-part is an IDENTIFIER argument) so they
	// use the generic FuncCall via parseKeywordFuncCall, matching the
	// AST comment that "Built-ins with ordinary name(arg, …) syntax
	// (DATE_ADD, DATE_DIFF, …) use plain FuncCall."
	// ------------------------------------------------------------------
	case tokCOALESCE:
		return p.parseCoalesceExpr()
	case tokNULLIF:
		return p.parseNullIfExpr()
	case tokSUBSTRING:
		return p.parseSubstringExpr()
	case tokTRIM:
		return p.parseTrimExpr()
	case tokEXTRACT:
		return p.parseExtractExpr()
	case tokDATE_ADD, tokDATE_DIFF:
		return p.parseDateFunction()

	// ------------------------------------------------------------------
	// Reserved-name scalar functions — FunctionCallReserved
	// (PartiQLParser.g4:611-614). Argument syntax is identical to the
	// generic IDENT(args) form; parseReservedFuncCall delegates to
	// parseFuncCallArgs for the comma-list. COUNT is excluded here
	// because it belongs to DAG node 14 (parser-aggregates).
	// ------------------------------------------------------------------
	case tokCHAR_LENGTH, tokCHARACTER_LENGTH, tokOCTET_LENGTH, tokBIT_LENGTH,
		tokUPPER, tokLOWER, tokSIZE, tokEXISTS:
		return p.parseReservedFuncCall()

	// ------------------------------------------------------------------
	// Real: sequenceConstructor (LIST/SEXP) (DAG node 15b).
	//
	// sequenceConstructor: (LIST|SEXP) '(' (expr (',' expr)*)? ')'. The
	// argument list is an ordinary comma-list, so these reuse the generic
	// FuncCall via parseKeywordFuncCall (Name=LIST/SEXP). A bare LIST/SEXP
	// without parens is a *type* name, not a constructor — that form is
	// handled by parseType in DDL/CAST contexts, never here.
	// ------------------------------------------------------------------
	case tokLIST, tokSEXP:
		return p.parseKeywordFuncCall()

	// ------------------------------------------------------------------
	// Stub: aggregates → parser-aggregates
	// ------------------------------------------------------------------
	case tokCOUNT:
		return nil, p.deferredFeature("COUNT() aggregate", "parser-aggregates (DAG node 14)")
	case tokMAX:
		return nil, p.deferredFeature("MAX() aggregate", "parser-aggregates (DAG node 14)")
	case tokMIN:
		return nil, p.deferredFeature("MIN() aggregate", "parser-aggregates (DAG node 14)")
	case tokSUM:
		return nil, p.deferredFeature("SUM() aggregate", "parser-aggregates (DAG node 14)")
	case tokAVG:
		return nil, p.deferredFeature("AVG() aggregate", "parser-aggregates (DAG node 14)")

	// ------------------------------------------------------------------
	// Stub: window functions → parser-window
	// ------------------------------------------------------------------
	case tokLAG:
		return nil, p.deferredFeature("LAG() window", "parser-window (DAG node 13)")
	case tokLEAD:
		return nil, p.deferredFeature("LEAD() window", "parser-window (DAG node 13)")
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("unexpected token %q in expression", p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// parseParenExpr handles tokPAREN_LEFT dispatch for primary expressions:
//
//	(expr)             → plain parenthesized expression, returns inner expr
//	(SELECT ...)       → routes through parseSelectExpr, which returns the
//	                     parser-select (DAG node 5) stub error
//	(expr, expr, ...)  → STUB: valueList deferred to parser-dml (DAG node 6)
//	(expr MATCH ...)   → GPML graph match, parsed by parseGraphMatch (graph.go)
//
// Grammar: exprTerm#ExprTermWrappedQuery (line 543) + exprGraphMatchMany
// (lines 625-626).
func (p *Parser) parseParenExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume (

	// Empty parens `()` are not valid in any PartiQL primary position.
	if p.cur.Type == tokPAREN_RIGHT {
		return nil, &ParseError{
			Message: "empty parenthesized expression",
			Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
		}
	}

	first, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	// Graph match: (expr MATCH ...) — GPML. parseGraphMatch (graph.go) owns
	// MATCH through the pattern list and stops on the closing paren, which the
	// PAREN_RIGHT handling below consumes.
	if p.cur.Type == tokMATCH {
		match, err := p.parseGraphMatch(first)
		if err != nil {
			return nil, err
		}
		first = match
	}

	// valueList: (expr, expr, ...) — deferred to parser-dml.
	if p.cur.Type == tokCOMMA {
		return nil, p.deferredFeature("valueList", "parser-dml (DAG node 6)")
	}

	// Plain (expr): consume the closing paren and return the inner expr.
	// Note: the returned expression does NOT get a new wrapping node —
	// parentheses are purely syntactic. The inner expression's Loc is
	// preserved as-is (we don't extend it to cover the outer parens).
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}

	// If the inner expression is a SubLink (parenthesized SELECT),
	// update its Loc to span the outer parens and return it directly.
	// The parens are syntactic but the SubLink carries the query.
	if sub, ok := first.(*ast.SubLink); ok {
		sub.Loc = ast.Loc{Start: start, End: rp.Loc.End}
		return sub, nil
	}

	return first, nil
}

// parseArrayLit parses `[expr, expr, ...]`. Empty brackets `[]` are
// valid (empty array literal).
//
// Grammar: array (line 649-650).
func (p *Parser) parseArrayLit() (*ast.ListLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume [
	var items []ast.ExprNode
	if p.cur.Type != tokBRACKET_RIGHT {
		for {
			item, err := p.parseExprTop()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokBRACKET_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.ListLit{
		Items: items,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseBagLit parses `<<expr, expr, ...>>`. Empty `<<>>` is valid.
//
// Grammar: bag (line 652-653).
func (p *Parser) parseBagLit() (*ast.BagLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume <<
	var items []ast.ExprNode
	if p.cur.Type != tokANGLE_DOUBLE_RIGHT {
		for {
			item, err := p.parseExprTop()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokANGLE_DOUBLE_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.BagLit{
		Items: items,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTupleLit parses `{key: value, key: value, ...}`. Empty `{}` is
// valid.
//
// Grammar: tuple (line 655-656) + pair (line 658-659).
func (p *Parser) parseTupleLit() (*ast.TupleLit, error) {
	start := p.cur.Loc.Start
	p.advance() // consume {
	var pairs []*ast.TuplePair
	if p.cur.Type != tokBRACE_RIGHT {
		for {
			pair, err := p.parseTuplePair()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, pair)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
	}
	rp, err := p.expect(tokBRACE_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.TupleLit{
		Pairs: pairs,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTuplePair parses one `key: value` entry inside a tuple literal.
func (p *Parser) parseTuplePair() (*ast.TuplePair, error) {
	key, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokCOLON); err != nil {
		return nil, err
	}
	value, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	return &ast.TuplePair{
		Key:   key,
		Value: value,
		Loc:   ast.Loc{Start: key.GetLoc().Start, End: value.GetLoc().End},
	}, nil
}

// parseParamRef consumes a QUESTION_MARK and returns an ast.ParamRef.
//
// Grammar: parameter (line 632-633).
func (p *Parser) parseParamRef() (*ast.ParamRef, error) {
	loc := p.cur.Loc
	p.advance()
	return &ast.ParamRef{Loc: loc}, nil
}

// parseCaseExpr parses a CASE expression — PartiQLParser.g4's caseExpr
// rule. The current token must be tokCASE. Handles both forms:
//
//   - Searched: CASE WHEN cond THEN result [WHEN ...]... [ELSE r] END
//   - Simple:   CASE operand WHEN val THEN result [WHEN ...]... [ELSE r] END
//
// Whether it's searched or simple is decided by lookahead after CASE:
// if the next token is WHEN, it's searched; otherwise we parse an
// Operand expression first. At least one WHEN clause is required.
//
// Owned by DAG node 15b.3.
func (p *Parser) parseCaseExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume CASE

	var operand ast.ExprNode
	if p.cur.Type != tokWHEN {
		// Simple form: parse the operand.
		op, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		operand = op
	}

	// Require at least one WHEN clause.
	if p.cur.Type != tokWHEN {
		return nil, &ParseError{
			Message: "CASE expression requires at least one WHEN clause",
			Loc:     p.cur.Loc,
		}
	}

	var whens []*ast.CaseWhen
	for p.cur.Type == tokWHEN {
		whenStart := p.cur.Loc.Start
		p.advance() // consume WHEN
		whenExpr, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokTHEN); err != nil {
			return nil, err
		}
		thenExpr, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		whens = append(whens, &ast.CaseWhen{
			When: whenExpr,
			Then: thenExpr,
			Loc:  ast.Loc{Start: whenStart, End: thenExpr.GetLoc().End},
		})
	}

	var elseExpr ast.ExprNode
	if p.cur.Type == tokELSE {
		p.advance() // consume ELSE
		e, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		elseExpr = e
	}

	endOff := p.cur.Loc.End
	if _, err := p.expect(tokEND); err != nil {
		return nil, err
	}

	return &ast.CaseExpr{
		Operand: operand,
		Whens:   whens,
		Else:    elseExpr,
		Loc:     ast.Loc{Start: start, End: endOff},
	}, nil
}

// parseReservedFuncCall parses a reserved-name function call —
// PartiQLParser.g4:611-614 FunctionCallReserved. The current token must
// be the reserved keyword (e.g. tokSIZE, tokEXISTS). It uppercases the
// keyword name, advances past it, and delegates to parseFuncCallArgs.
//
// Owned by DAG node 15b (parser-builtins-typed); 15b.1 covers SIZE and
// EXISTS, with the rest of the FunctionCallReserved keywords landing
// in subsequent 15b sub-PRs.
func (p *Parser) parseReservedFuncCall() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	name := strings.ToUpper(p.cur.Str)
	p.advance() // consume the reserved keyword
	args, endOff, err := p.parseFuncCallArgs()
	if err != nil {
		return nil, err
	}
	return &ast.FuncCall{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: endOff},
	}, nil
}

// parseKeywordFuncCall parses a keyword-named function whose argument
// syntax is the ordinary comma-separated list — i.e. structurally
// identical to FunctionCallIdent but spelled with a dedicated keyword
// token. It uppercases the keyword name, advances past it, and delegates
// to parseFuncCallArgs.
//
// Used for the LIST/SEXP sequence constructors (PartiQLParser.g4:569-570
// sequenceConstructor). Both produce an *ast.FuncCall, matching the AST
// guidance that built-ins with ordinary name(arg, …) syntax use the
// generic node.
//
// Note: arity is intentionally NOT enforced here — LIST/SEXP accept any
// count, including zero. DATE_ADD/DATE_DIFF do NOT use this path: their
// first argument is a constrained date-part identifier rather than a
// general expression, so they have a dedicated parser (parseDateFunction).
func (p *Parser) parseKeywordFuncCall() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	name := strings.ToUpper(p.cur.Str)
	p.advance() // consume the keyword
	args, endOff, err := p.parseFuncCallArgs()
	if err != nil {
		return nil, err
	}
	return &ast.FuncCall{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: endOff},
	}, nil
}

// parseDateFunction parses DATE_ADD / DATE_DIFF (PartiQLParser.g4:608-609):
//
//	dateFunction
//	    : func=(DATE_ADD|DATE_DIFF)
//	        PAREN_LEFT dt=IDENTIFIER COMMA expr COMMA expr PAREN_RIGHT;
//
// Unlike an ordinary function call, the FIRST argument is a constrained
// date-part — `dt=IDENTIFIER`, a bare unquoted identifier (year, month,
// day, hour, minute, second, timezone_hour, timezone_minute). It is NOT a
// general expression: a quoted string, a double-quoted identifier, an
// arithmetic/path expression, or a numeric literal in that position is a
// parse error. Arguments two and three remain general expressions. This
// mirrors the EXTRACT field handling (see parseExtractExpr) where the
// grammar likewise pins the leading position to an unquoted IDENTIFIER.
//
// truth1 (PartiQL spec) and the partiql-lang-kotlin reference impl both
// document the part as a bare keyword — e.g. `DATE_ADD(year, 5, ts)` —
// corroborating the grammar.
//
// The AST shape is the generic *ast.FuncCall (Name=DATE_ADD/DATE_DIFF,
// Args=[partRef, expr, expr]); the part lands as a VarRef so the node
// stays identical to its previously-shipped form — only the now-rejected
// over-permissive inputs change behavior.
func (p *Parser) parseDateFunction() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	name := strings.ToUpper(p.cur.Str)
	p.advance() // consume DATE_ADD / DATE_DIFF
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}

	// First argument: the date-part. dt=IDENTIFIER — a bare unquoted
	// identifier only. Reject quoted identifiers, string literals,
	// expressions, paths, and any other non-identifier token here.
	if p.cur.Type != tokIDENT {
		return nil, &ParseError{
			Message: fmt.Sprintf("%s date part must be an unquoted identifier, got %q", name, p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	part := &ast.VarRef{Name: p.cur.Str, Loc: p.cur.Loc}
	p.advance() // consume the date-part identifier
	args := []ast.ExprNode{part}

	// Arguments two and three: general expressions, comma-separated.
	for i := 0; i < 2; i++ {
		if _, err := p.expect(tokCOMMA); err != nil {
			return nil, err
		}
		arg, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.FuncCall{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseCastExpr parses CAST / CAN_CAST / CAN_LOSSLESS_CAST — the three
// cast-family rules (PartiQLParser.g4:593-600):
//
//	cast            : CAST              '(' expr AS type ')'
//	canCast         : CAN_CAST          '(' expr AS type ')'
//	canLosslessCast : CAN_LOSSLESS_CAST '(' expr AS type ')'
//
// The current token must be the cast keyword; kind selects which of the
// three. All three share the `'(' expr AS type ')'` body. The target
// type is parsed by the shared parseType helper.
func (p *Parser) parseCastExpr(kind ast.CastKind) (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume CAST / CAN_CAST / CAN_LOSSLESS_CAST
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	expr, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokAS); err != nil {
		return nil, err
	}
	asType, err := p.parseType()
	if err != nil {
		return nil, err
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.CastExpr{
		Kind:   kind,
		Expr:   expr,
		AsType: asType,
		Loc:    ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseExtractExpr parses EXTRACT (PartiQLParser.g4:602-603):
//
//	extract : EXTRACT '(' IDENTIFIER FROM rhs=expr ')'
//
// The datetime field is an IDENTIFIER (e.g. YEAR, MONTH, TIMEZONE_HOUR),
// not a reserved keyword — the grammar and the official spec both leave
// the field unrestricted at the syntax level (the valid-field check is a
// semantic concern, not a parse concern), so we accept any unquoted
// identifier and store its raw text. A double-quoted identifier is
// rejected: the grammar's IDENTIFIER alternative does not include the
// IDENTIFIER_QUOTED token.
func (p *Parser) parseExtractExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume EXTRACT
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	if p.cur.Type == tokIDENT_QUOTED {
		return nil, &ParseError{
			Message: "EXTRACT field must be an unquoted identifier",
			Loc:     p.cur.Loc,
		}
	}
	if p.cur.Type != tokIDENT {
		return nil, &ParseError{
			Message: fmt.Sprintf("expected identifier for EXTRACT field, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	field := p.cur.Str
	p.advance() // consume the field identifier
	if _, err := p.expect(tokFROM); err != nil {
		return nil, err
	}
	from, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.ExtractExpr{
		Field: field,
		From:  from,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseTrimExpr parses TRIM (PartiQLParser.g4:605-606):
//
//	trimFunction : TRIM '(' ( mod=IDENTIFIER? sub=expr? FROM )? target=expr ')'
//
// The optional `( mod? sub? FROM )` prefix is matched ONLY when a FROM
// keyword appears at the top level of the argument list (before TRIM's
// matching close paren). Everything between the parens that is NOT
// terminated by such a FROM is the bare `target=expr`. `mod` is an
// IDENTIFIER constrained semantically to LEADING / TRAILING / BOTH (the
// official spec header is `TRIM([[LEADING|TRAILING|BOTH] removalChar FROM]
// str)`); we map it to the TrimSpec enum. The grammar's broader
// `mod=IDENTIFIER` cannot be faithfully represented by the TrimSpec enum
// for arbitrary identifiers, and treating non-LEADING/TRAILING/BOTH
// identifiers as part of the target is the only spec-consistent reading.
//
// Form disambiguation uses a BOUNDED, LINEAR token lookahead
// (trimHasTopLevelFrom): from the open paren we scan forward tracking
// paren depth and report whether a depth-0 FROM precedes the matching
// close paren. The scan touches each token at most once and parses no
// sub-expression, so it is O(n) in the argument length; an earlier
// revision instead speculatively parsed the WHOLE argument via
// parseExprTop, discarded it, then re-parsed it, which made nested bare
// targets (TRIM(TRIM(...TRIM(x)...))) re-parse at every level — O(2^n),
// a DoS-class blowup on live editor input. With the lookahead each part
// (mod / sub / target) is parsed EXACTLY ONCE per TRIM.
//
// A top-level FROM is exactly equivalent to the prefix form for every
// VALID input: a bare target `parseExprTop` consumes the entire argument
// and lands on the close paren only when no top-level FROM exists (FROM is
// never an infix operator, so it always terminates expression parsing).
// For INVALID inputs both readings reject; the verdict (accept/reject) is
// preserved either way.
//
//	TRIM(both(x))     => target is the function call both(x)
//	TRIM(leading)     => target is the var ref `leading`
//	TRIM(trailing+1)  => target is the expression trailing + 1
//
// parse correctly even though BOTH / LEADING / TRAILING are non-reserved
// identifiers: a single-token lookahead that an earlier revision used
// could not tell `BOTH FROM` (modifier) from `both(x)` (call target).
//
// Forms handled:
//
//	TRIM(s)                  -> {From:s}
//	TRIM(FROM s)             -> {From:s}
//	TRIM(LEADING FROM s)     -> {Spec:LEADING From:s}
//	TRIM(' ' FROM s)         -> {Sub:' ' From:s}
//	TRIM(BOTH ' ' FROM s)    -> {Spec:BOTH Sub:' ' From:s}
//	TRIM(both)               -> {From:both}      (no FROM => 'both' is target)
//	TRIM(both(x))            -> {From:both(x)}   (no FROM => call is target)
func (p *Parser) parseTrimExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume TRIM
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}

	// Decide the form with a bounded linear lookahead rather than a
	// speculative full-expression parse: a top-level FROM selects the
	// `( mod? sub? FROM ) target` prefix form, its absence the bare
	// `target` form.
	if p.trimHasTopLevelFrom() {
		spec := ast.TrimSpecNone
		// Optional modifier: a LEADING / TRAILING / BOTH keyword. ANTLR
		// binds `mod=IDENTIFIER` greedily, so consume it when present.
		if s, isSpec := trimSpecFor(p.cur); isSpec {
			spec = s
			p.advance() // consume the modifier identifier
		}
		// Optional removal char `sub=expr`, present iff the next token is
		// not yet the FROM. Parsed exactly once. (When no modifier was
		// consumed this also covers the `( sub FROM )` form.)
		var sub ast.ExprNode
		if p.cur.Type != tokFROM {
			var err error
			sub, err = p.parseExprTop()
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(tokFROM); err != nil {
			return nil, err
		}
		target, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		return p.finishTrim(start, spec, sub, target)
	}

	// No top-level FROM: the entire argument is the target expression
	// (e.g. TRIM(s), TRIM(both), TRIM(both(x)), TRIM(leading || x)).
	// Spec is None and there is no removal char in this form.
	target, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	return p.finishTrim(start, ast.TrimSpecNone, nil, target)
}

// trimHasTopLevelFrom reports whether TRIM's argument list contains a FROM
// keyword at the top level — paren depth 0 relative to the open paren the
// caller has already consumed — before that paren's matching close. It is
// the linear disambiguator for the optional `( mod? sub? FROM )` prefix
// (PartiQLParser.g4:606).
//
// The scan is bounded and side-effect free: it snapshots the parse
// position, advances token-by-token tracking the nesting of every paired
// delimiter — () [] {} <<>> — and returns at the first depth-0 FROM, or
// when the matching close paren / EOF ends the argument list. It then
// restores the snapshot, leaving cur exactly where the caller left it.
// Because it consumes each token at most once and never recurses into a
// sub-expression, it is O(n) in the number of tokens before the matching
// close paren.
//
// Depth tracking means a FROM nested inside a subquery or collection
// literal — e.g. TRIM((SELECT a FROM t) FROM s) — is NOT mistaken for the
// prefix delimiter: only the depth-0 FROM (the one separating sub from
// target) selects the prefix form.
func (p *Parser) trimHasTopLevelFrom() bool {
	snapshot := p.save()
	defer p.restore(snapshot)

	depth := 0
	for {
		switch p.cur.Type {
		case tokEOF:
			// Unterminated argument list (also the lexer-error sentinel):
			// no top-level FROM was found. The bare-target parse below
			// surfaces the real syntax/lexer error.
			return false
		case tokFROM:
			if depth == 0 {
				return true
			}
		case tokPAREN_LEFT, tokBRACKET_LEFT, tokBRACE_LEFT, tokANGLE_DOUBLE_LEFT:
			depth++
		case tokPAREN_RIGHT:
			if depth == 0 {
				// TRIM's own matching close paren: argument list ended
				// without a top-level FROM.
				return false
			}
			depth--
		case tokBRACKET_RIGHT, tokBRACE_RIGHT, tokANGLE_DOUBLE_RIGHT:
			if depth > 0 {
				depth--
			}
		}
		p.advance()
	}
}

// finishTrim consumes the closing paren and builds the TrimExpr.
func (p *Parser) finishTrim(start int, spec ast.TrimSpec, sub, target ast.ExprNode) (ast.ExprNode, error) {
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.TrimExpr{
		Spec: spec,
		Sub:  sub,
		From: target,
		Loc:  ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// trimSpecFor maps a token to a TrimSpec if it is an unquoted identifier
// spelling LEADING / TRAILING / BOTH (case-insensitive). The second
// return is false for any other token.
func trimSpecFor(tok Token) (ast.TrimSpec, bool) {
	if tok.Type != tokIDENT {
		return ast.TrimSpecNone, false
	}
	switch strings.ToUpper(tok.Str) {
	case "LEADING":
		return ast.TrimSpecLeading, true
	case "TRAILING":
		return ast.TrimSpecTrailing, true
	case "BOTH":
		return ast.TrimSpecBoth, true
	default:
		return ast.TrimSpecNone, false
	}
}

// parseSubstringExpr parses SUBSTRING (PartiQLParser.g4:572-575):
//
//	substring
//	  : SUBSTRING '(' expr ( ',' expr ( ',' expr )? )? ')'
//	  | SUBSTRING '(' expr ( FROM expr ( FOR expr )? )? ')'
//
// Two argument styles: the comma form `SUBSTRING(str, start[, len])` and
// the keyword form `SUBSTRING(str FROM start [FOR len])`. The two styles
// are mutually exclusive — the first separator (COMMA vs FROM) selects
// the form. Mixing them (e.g. `SUBSTRING(s, 2 FOR 3)`) is rejected
// because the trailing FOR/extra token is left for the PAREN_RIGHT
// expectation, which fails.
//
// In both forms the start/length operands are optional in the grammar
// (the `(...)?` groups), so `SUBSTRING(s)` is accepted as
// {Expr:s From:<nil>}. The From field is left nil in that degenerate
// case rather than fabricating a default.
func (p *Parser) parseSubstringExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume SUBSTRING
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	expr, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	var from, forLen ast.ExprNode
	switch p.cur.Type {
	case tokCOMMA:
		// Comma form: , start ( , len )?
		p.advance() // consume ,
		from, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == tokCOMMA {
			p.advance() // consume ,
			forLen, err = p.parseExprTop()
			if err != nil {
				return nil, err
			}
		}
	case tokFROM:
		// Keyword form: FROM start ( FOR len )?
		p.advance() // consume FROM
		from, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == tokFOR {
			p.advance() // consume FOR
			forLen, err = p.parseExprTop()
			if err != nil {
				return nil, err
			}
		}
	}

	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.SubstringExpr{
		Expr: expr,
		From: from,
		For:  forLen,
		Loc:  ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseCoalesceExpr parses COALESCE (PartiQLParser.g4:554-555):
//
//	coalesce : COALESCE '(' expr ( ',' expr )* ')'
//
// At least one argument is required (the leading `expr` is not optional).
// Empty COALESCE() is rejected.
func (p *Parser) parseCoalesceExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume COALESCE
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	var args []ast.ExprNode
	for {
		arg, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.CoalesceExpr{
		Args: args,
		Loc:  ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseNullIfExpr parses NULLIF (PartiQLParser.g4:551-552):
//
//	nullIf : NULLIF '(' expr ',' expr ')'
//
// Exactly two arguments are required.
func (p *Parser) parseNullIfExpr() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume NULLIF
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	left, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokCOMMA); err != nil {
		return nil, err
	}
	right, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.NullIfExpr{
		Left:  left,
		Right: right,
		Loc:   ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

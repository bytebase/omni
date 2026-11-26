package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// Binding power constants for Pratt parsing (low to high).
// Doris operator precedence from DorisParser.g4:
//
//	OR / || < XOR < AND / && < NOT (prefix) < predicates/comparison <
//	bitwise | < bitwise & < shifts << >> < add +- < mul */% DIV MOD <
//	bitwise ^ < unary (-, +, ~) < primary
const (
	bpNone    = 0
	bpOr      = 10 // OR, ||
	bpAnd     = 20 // AND, &&
	bpNot     = 25 // NOT (prefix)
	bpPred    = 30 // IS, BETWEEN, IN, LIKE, REGEXP/RLIKE
	bpCmp     = 35 // =, <, >, <=, >=, <>, !=, <=>
	bpBitOr   = 40 // |
	bpBitAnd  = 45 // &
	bpShift   = 50 // <<, >>
	bpAdd     = 55 // +, -
	bpMul     = 60 // *, /, %, DIV
	bpBitXor  = 65 // ^
	bpMatch   = 68 // MATCH: below the unary operators, above every binary one
	bpUnary   = 70 // unary -, +, ~
	bpPrimary = 80 // atoms
)

// ---------------------------------------------------------------------------
// Core Pratt loop
// ---------------------------------------------------------------------------

// parseExpr is the entry point for expression parsing.
func (p *Parser) parseExpr() (ast.Node, error) {
	return p.parseExprPrec(bpNone + 1)
}

// parseExprPrec parses an expression with the given minimum binding power.
func (p *Parser) parseExprPrec(minBP int) (ast.Node, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}

	for {
		// primaryExpression NOT? MATCH primaryExpression (StarRocks.g4:2425)
		// binds before any binary operator, but after the unary ones that
		// come before it in primaryExpression (:2408-2409): !a MATCH b is
		// (!a) MATCH b.
		if minBP <= bpMatch && (p.cur.Kind == kwMATCH || (p.cur.Kind == kwNOT && p.peekNext().Kind == kwMATCH)) {
			// MATCH is nonReserved (StarRocks.g4:2930), so where no primary
			// expression follows it the engine reads it as a column label:
			// `SELECT 1 MATCH` is one item aliased MATCH. The alternative
			// cannot hold then, so the token is left to the label reader.
			mark := p.save()
			matched, matchErr := p.parseMatchExpr(left, p.cur.Kind == kwNOT)
			if matchErr != nil {
				p.restore(mark)
				break
			}
			left = matched
			continue
		}
		// primaryExpression COLLATE (StarRocks.g4:2403) binds before any operator.
		if p.cur.Kind == kwCOLLATE {
			if left, err = p.parseCollateExpr(left); err != nil {
				return nil, err
			}
			// The result is a primaryExpression again, which a subscript
			// may follow (StarRocks.g4:2419).
			if left, err = p.parsePostfixSuffixes(left); err != nil {
				return nil, err
			}
			continue
		}
		bp, op, special, ok := p.infixBindingPower()
		if !ok || bp < minBP {
			break
		}

		if special != "" {
			switch special {
			case "IS":
				left, err = p.parseIsExpr(left)
			case "BETWEEN":
				left, err = p.parseBetweenExpr(left, false)
			case "NOT_BETWEEN":
				left, err = p.parseBetweenExpr(left, true)
			case "IN":
				left, err = p.parseInExpr(left, false)
			case "NOT_IN":
				left, err = p.parseInExpr(left, true)
			case "LIKE":
				left, err = p.parseLikeExpr(left, false)
			case "NOT_LIKE":
				left, err = p.parseLikeExpr(left, true)
			case "REGEXP":
				left, err = p.parseRegexpExpr(left, false)
			case "NOT_REGEXP":
				left, err = p.parseRegexpExpr(left, true)
			}
			if err != nil {
				return nil, err
			}
			continue
		}

		// Regular binary operator: consume the operator, parse right side.
		p.advance()
		right, err := p.parseExprPrec(bp + 1)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(right).End},
		}
	}

	return left, nil
}

// infixBindingPower returns (bp, op, special, ok) for the current token.
// If special is non-empty, the caller must dispatch to the named handler
// instead of building a BinaryExpr.
func (p *Parser) infixBindingPower() (int, ast.BinaryOp, string, bool) {
	switch p.cur.Kind {
	// Logical
	case kwOR, tokDoublePipes:
		return bpOr, ast.BinOr, "", true
	case kwAND, tokLogicalAnd:
		return bpAnd, ast.BinAnd, "", true

	// NOT as infix prefix for BETWEEN, IN, LIKE, REGEXP/RLIKE
	case kwNOT:
		next := p.peekNext()
		switch next.Kind {
		case kwBETWEEN:
			return bpPred, 0, "NOT_BETWEEN", true
		case kwIN:
			return bpPred, 0, "NOT_IN", true
		case kwLIKE:
			return bpPred, 0, "NOT_LIKE", true
		case kwREGEXP, kwRLIKE:
			return bpPred, 0, "NOT_REGEXP", true
		}
		return 0, 0, "", false

	// Predicates
	case kwIS:
		return bpPred, 0, "IS", true
	case kwBETWEEN:
		return bpPred, 0, "BETWEEN", true
	case kwIN:
		return bpPred, 0, "IN", true
	case kwLIKE:
		return bpPred, 0, "LIKE", true
	case kwREGEXP, kwRLIKE:
		return bpPred, 0, "REGEXP", true

	// Comparison operators
	case int('='):
		return bpCmp, ast.BinEq, "", true
	case tokNotEq:
		return bpCmp, ast.BinNe, "", true
	case int('<'):
		return bpCmp, ast.BinLt, "", true
	case int('>'):
		return bpCmp, ast.BinGt, "", true
	case tokLessEq:
		return bpCmp, ast.BinLe, "", true
	case tokGreaterEq:
		return bpCmp, ast.BinGe, "", true
	case tokNullSafeEq:
		return bpCmp, ast.BinNullSafeEq, "", true

	// Bitwise OR
	case int('|'):
		return bpBitOr, ast.BinBitOr, "", true

	// Bitwise AND
	case int('&'):
		return bpBitAnd, ast.BinBitAnd, "", true

	// Shifts
	case tokShiftLeft:
		return bpShift, ast.BinShiftLeft, "", true
	case tokShiftRight:
		return bpShift, ast.BinShiftRight, "", true

	// Addition / subtraction
	case int('+'):
		return bpAdd, ast.BinAdd, "", true
	case int('-'):
		return bpAdd, ast.BinSub, "", true

	// Multiplication / division / modulo
	case int('*'):
		return bpMul, ast.BinMul, "", true
	case int('/'):
		return bpMul, ast.BinDiv, "", true
	case int('%'):
		return bpMul, ast.BinMod, "", true
	case kwDIV:
		return bpMul, ast.BinIntDiv, "", true
	case kwMOD:
		return bpMul, ast.BinMod, "", true
	// The shift words of StarRocks.g4:2392-2394, which the lexer spells
	// BITSHIFTLEFT and so on. They bind looser than | (:2391), each looser
	// than the one before.
	case kwBITSHIFTLEFT:
		return bpBitOr - 1, ast.BinShiftLeft, "", true
	case kwBITSHIFTRIGHT:
		return bpBitOr - 2, ast.BinShiftRight, "", true
	case kwBITSHIFTRIGHTLOGICAL:
		return bpBitOr - 3, ast.BinShiftRightLogical, "", true

	// Bitwise XOR
	case int('^'):
		return bpBitXor, ast.BinBitXor, "", true
	}

	return 0, 0, "", false
}

// ---------------------------------------------------------------------------
// Prefix dispatch
// ---------------------------------------------------------------------------

// parsePrefixExpr handles prefix unary operators (-, +, ~, NOT) and delegates
// everything else to parsePrimaryExpr.
func (p *Parser) parsePrefixExpr() (ast.Node, error) {
	switch p.cur.Kind {
	case int('-'):
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryMinus,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case int('+'):
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryPlus,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case int('~'):
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryBitNot,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case kwNOT:
		// NOT as a prefix unary operator. Only if the next token does NOT form
		// an infix predicate (NOT BETWEEN, NOT IN, NOT LIKE, NOT REGEXP/RLIKE).
		// In the prefix case we are at the start of an expression, so NOT always
		// acts as prefix.
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpNot)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryNot,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case kwBINARY:
		// BINARY <expr> cast-to-binary operator. StarRocks's grammar is
		// `(BINARY)? booleanExpression`, so the operand is a booleanExpression —
		// it spans comparisons/predicates but stops below AND/OR/NOT:
		// BINARY name = 'abc' => BINARY(name = 'abc'); BINARY a AND b => (BINARY a) AND b.
		// BINARY is also nonReserved, so where no operand follows it names a
		// column: the pinned server takes `SELECT BINARY` and `SELECT BINARY
		// zzz`, and refuses the operator with nothing to apply it to.
		mark := p.save()
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpNot + 1)
		if err != nil {
			if isNestingError(err) {
				return nil, err
			}
			p.restore(mark)
			return p.parseIdentExpr()
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryBinary,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case int('!'):
		// LOGICAL NOT (StarRocks.g4:2409, before a primaryExpression).
		start := p.cur.Loc
		p.advance()
		operand, err := p.parseExprPrec(bpUnary)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Op:   ast.UnaryNot,
			Expr: operand,
			Loc:  ast.Loc{Start: start.Start, End: ast.NodeLoc(operand).End},
		}, nil

	case kwEXISTS:
		return p.parseExistsExpr()

	case kwINTERVAL:
		// interval (StarRocks.g4:2764) is INTERVAL value from=unitIdentifier,
		// with the unit required. INTERVAL is also nonReserved, so a form with
		// no unit — `SELECT INTERVAL`, `SELECT INTERVAL zzz` — is the column,
		// which is what the pinned server answers.
		mark := p.save()
		node, err := p.parseIntervalExpr()
		if err == nil {
			return node, nil
		}
		if isNestingError(err) {
			return nil, err
		}
		p.restore(mark)
		return p.parseIdentExpr()

	default:
		node, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		return p.parsePostfixSuffixes(node)
	}
}

// parsePostfixSuffixes parses any chain of suffixes that binds tighter than
// every operator: element access (m['k'], arr[1][2]) and the JSON path
// operator (j->'$.a'). The engine chains them in either order —
// `j->'$.b'['c']` and `j['b']->'$.c'` both parse — so one loop takes both.
//
// A subscript index is a valueExpression — comparisons are allowed but
// AND/OR/NOT are not, hence bpNot+1. Unlike Doris, the slice form
// `arr[begin:end]` is engine-rejected here, so a ':' inside the brackets
// stays a syntax error.
func (p *Parser) parsePostfixSuffixes(node ast.Node) (ast.Node, error) {
	for {
		switch {
		case p.cur.Kind == int('['):
			p.advance() // consume '['
			index, err := p.parseExprPrec(bpNot + 1)
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(int(']'))
			if err != nil {
				return nil, err
			}
			node = &ast.ElementAtExpr{
				Value: node,
				Index: index,
				Loc:   ast.Loc{Start: ast.NodeLoc(node).Start, End: closeTok.Loc.End},
			}

		case p.cur.Kind == tokArrow && p.peekNext().Kind == tokString:
			arrow, err := p.parseJSONArrow(node)
			if err != nil {
				return nil, err
			}
			node = arrow

		case p.cur.Kind == tokDotIdent:
			// primaryExpression DOT_IDENTIFIER (StarRocks.g4:2406,
			// StarRocksLex.g4:574): a field name that opens with digits, which
			// the lexer reads as one token, as '.1a' is.
			field := p.advance()
			node = &ast.DereferenceExpr{
				Expr:  node,
				Field: field.Str,
				Loc:   ast.Loc{Start: ast.NodeLoc(node).Start, End: field.Loc.End},
			}

		case p.cur.Kind == int('.') && isIdentifierTokenQualified(p.peekNext().Kind):
			// primaryExpression '.' identifier (StarRocks.g4:2406). A column
			// name's parts were read with it, so this is a field of another
			// expression.
			p.advance() // consume '.'
			field := p.advance()
			node = &ast.DereferenceExpr{
				Expr:  node,
				Field: field.Str,
				Loc:   ast.Loc{Start: ast.NodeLoc(node).Start, End: field.Loc.End},
			}

		default:
			return node, nil
		}
	}
}

// parseJSONArrow parses the `-> '<path>'` suffix of the JSON path operator,
// given the already-parsed expression to its left.
//
// The one-token lookahead in the caller is the engine's own disambiguation
// rule, not a heuristic: StarRocks 3.4.10 accepts `X -> Y` exactly when Y is
// a bare string literal (the JSON arrow, X any primary expression) or X is an
// identifier or identifier list (a lambda, Y any expression). `SELECT 1 -> 1`
// satisfies neither and is a syntax error. So a `->` followed by a string is
// the arrow wherever it appears, including in the higher-order argument
// position a lambda also lives in: the engine reads `array_map(j -> '$.a',
// [1,2])` as the arrow and then reports "There are no lambda functions in
// high-order function array_map". parseFuncArgList's lambda check runs after
// this one for the same reason.
//
// That ordering is also what lineage needs. The arrow reading keeps the left
// operand an ordinary operand of a BinaryExpr, so ast.Walk reaches a ColumnRef
// there and query_span reports it as a source column. Reading it as a lambda
// instead would bind `j` as a parameter, and the collector drops bound names —
// `SELECT secret->'$.a' FROM t` would lose the column that masking keys on.
//
// The path is a bare string token, matching the grammar: the engine does not
// take a parenthesized or computed path here (`j -> ('$.a')` and
// `j -> CONCAT('$.','a')` fall through to its lambda branch), and it does not
// concatenate adjacent string literals. The caller's lookahead decides whether
// to enter this production at all; the path is still taken with expect, so the
// production reports a missing one the way every other one does.
func (p *Parser) parseJSONArrow(left ast.Node) (ast.Node, error) {
	p.advance() // consume '->'
	pathTok, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	return &ast.BinaryExpr{
		Op:   ast.BinJSONArrow,
		Left: left,
		Right: &ast.Literal{
			Kind:  ast.LitString,
			Value: pathTok.Str,
			Loc:   pathTok.Loc,
		},
		Loc: ast.Loc{Start: ast.NodeLoc(left).Start, End: pathTok.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// Primary expression dispatch
// ---------------------------------------------------------------------------

// parsePrimaryExpr parses atomic expressions: literals, identifiers,
// function calls, parenthesized expressions, CASE, CAST, etc.
func (p *Parser) parsePrimaryExpr() (ast.Node, error) {
	switch p.cur.Kind {
	// KEY ident — Doris encryption-key reference (used in AES_ENCRYPT/AES_DECRYPT
	// function arguments). Represent as a ColumnRef with a "KEY/" prefix so the
	// AST is round-trippable.
	case kwKEY:
		keyLoc := p.advance().Loc
		if isIdentifierToken(p.cur.Kind) {
			name, err := p.parseMultipartIdentifier()
			if err != nil {
				return nil, err
			}
			return &ast.ColumnRef{
				Name: &ast.ObjectName{Parts: append([]string{"KEY"}, name.Parts...), Loc: ast.Loc{Start: keyLoc.Start, End: name.Loc.End}},
				Loc:  ast.Loc{Start: keyLoc.Start, End: name.Loc.End},
			}, nil
		}
		// Bare KEY — return as literal
		return &ast.Literal{Kind: ast.LitKeyword, Value: "KEY", Loc: keyLoc}, nil

	case tokInt:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitInt,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case tokFloat:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitFloat,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case tokString:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitString,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case tokHexLiteral:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitHex,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case tokBitLiteral:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitBit,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case kwTRUE:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitBool,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case kwFALSE:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitBool,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case kwNULL:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitNull,
			Value: tok.Str,
			Loc:   tok.Loc,
		}, nil

	case int('('):
		return p.parseParenExprOrSubquery()

	case kwMAP:
		// map<k,v>{...} (typed) or MAP{...} (untyped) map constructor; otherwise
		// MAP is an ordinary identifier (it is a non-reserved keyword, so e.g.
		// `map < 5` is a column comparison).
		switch p.peekNext().Kind {
		case int('{'):
			// Untyped MAP{...} — unambiguous.
			start := p.cur.Loc.Start
			p.advance() // consume MAP
			return p.finishMapLiteral(start, nil)
		case int('<'):
			// Ambiguous with `map < expr`: speculatively parse the type and
			// commit only if '{' follows, else backtrack to the identifier path.
			saved := p.save()
			start := p.cur.Loc.Start
			p.speculating++
			mapType, err := p.parseDataType()
			p.speculating--
			if err == nil && p.cur.Kind == int('{') {
				return p.finishMapLiteral(start, mapType)
			}
			if isNestingError(err) {
				return nil, err
			}
			p.restore(saved)
		}
		return p.parseIdentExpr()

	case kwARRAY:
		// array<t>[...] typed array constructor; otherwise an ordinary identifier
		// (ARRAY is non-reserved too, so `array < 5` is a column comparison).
		if p.peekNext().Kind == int('<') {
			saved := p.save()
			start := p.cur.Loc.Start
			p.speculating++
			elemType, err := p.parseDataType()
			p.speculating--
			if err == nil && p.cur.Kind == int('[') {
				return p.finishArrayLiteral(start, elemType)
			}
			if isNestingError(err) {
				return nil, err
			}
			p.restore(saved)
		}
		return p.parseIdentExpr()

	case int('['):
		// ARRAY_ELEMENT, the lexer's '[*]' (StarRocksLex.g4:506), is a
		// nonReserved identifier, so `SELECT [*]` names a column. The token is
		// matched by its three characters touching, as the lexer rule does.
		if node, ok := p.parseArrayElementIdent(); ok {
			return node, nil
		}
		// Untyped array constructor: [ e, ... ].
		return p.finishArrayLiteral(p.cur.Loc.Start, nil)

	case tokDotDotDot:
		// DOTDOTDOT, the lexer's '...' (StarRocksLex.g4:604), is a nonReserved
		// identifier, so `SELECT ...` names a column named "...".
		tok := p.advance()
		loc := tok.Loc
		return &ast.ColumnRef{
			Name: &ast.ObjectName{Parts: []string{"..."}, Loc: loc},
			Loc:  loc,
		}, nil

	case kwCASE:
		return p.parseCaseExpr()

	case kwDATE, kwDATETIME:
		// dateLiteral (StarRocks.g4:2432) is (DATE | DATETIME) string. Both
		// words are nonReserved, but before a string they are this literal,
		// which the pinned server refuses when the text is no date
		// ("Invalid date literal"). It is kept as a cast of the string.
		if p.peekNext().Kind == tokString {
			return p.parseDateLiteral()
		}
		return p.parseIdentExpr()

	case kwCAST:
		// cast (StarRocks.g4:2413) is CAST '(' expression AS type ')'. CAST is
		// also nonReserved, so unapplied it names a column: the pinned server
		// takes `SELECT CAST` and `SELECT CAST zzz`.
		if p.peekNext().Kind == int('(') {
			return p.parseCastExpr()
		}
		return p.parseIdentExpr()

	case kwCONVERT:
		// CONVERT is non-reserved, so only treat it specially when applied.
		if p.peekNext().Kind == int('(') {
			return p.parseConvertExpr()
		}
		return p.parseIdentExpr()

	case kwEXTRACT:
		// extract (StarRocks.g4:2441) is EXTRACT '(' identifier FROM
		// valueExpression ')'. EXTRACT is also nonReserved, so unapplied it
		// names a column: the pinned server takes `SELECT EXTRACT` and
		// `SELECT EXTRACT zzz`.
		if p.peekNext().Kind == int('(') {
			return p.parseExtractExpr()
		}
		return p.parseIdentExpr()

	case kwCEIL, kwFLOOR:
		// Both are nonReserved, and both are also written applied only:
		// specialFunctionExpression has CEIL '(' expression ')' and the same
		// for FLOOR (StarRocks.g4:2515). The pinned server refuses the bare
		// word in a name position (`SELECT CEIL`, `SELECT FLOOR`) and takes
		// `CEIL(1)`, so the name needs the argument list.
		if p.peekNext().Kind != int('(') {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parseIdentExpr()

	case kwTIMESTAMPADD, kwTIMESTAMPDIFF:
		// Their first argument is a unitIdentifier, not an expression, which an
		// argument list cannot read now that MICROSECOND and MILLISECOND are
		// reserved; as bare words they are non-reserved identifiers.
		if p.peekNext().Kind == int('(') {
			return p.parseTimestampUnitCall()
		}
		return p.parseIdentExpr()

	// Nilary functions: CURRENT_DATE, CURRENT_TIME, CURRENT_TIMESTAMP,
	// CURRENT_USER, CURRENT_ROLE, LOCALTIME, LOCALTIMESTAMP.
	case kwCURRENT_DATE, kwCURRENT_TIME, kwCURRENT_TIMESTAMP,
		kwCURRENT_USER, kwCURRENT_ROLE, kwLOCALTIME, kwLOCALTIMESTAMP:
		return p.parseNilaryFunction()

	case tokDoubleAt:
		return p.parseVariableRef(true)

	case int('@'):
		return p.parseVariableRef(false)

	case tokPlaceholder:
		tok := p.advance()
		return &ast.Literal{
			Kind:  ast.LitNull, // placeholder treated as a null-ish literal
			Value: "?",
			Loc:   tok.Loc,
		}, nil

	default:
		// A reserved keyword that the grammar still allows as a function name,
		// but only when it is immediately applied: `TRIM(s)` is a call,
		// `LEFT JOIN` is not.
		if isFunctionNameKeyword(p.cur.Kind) && p.peekNext().Kind == int('(') {
			tok := p.advance()
			name := &ast.ObjectName{Parts: []string{strings.ToUpper(tok.Str)}, Loc: tok.Loc}
			return p.parseFuncCall(name)
		}

		// Identifier-based: column ref or function call.
		if p.isExprIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// isFunctionNameKeyword reports whether kind is a reserved keyword that is
// nonetheless usable as a function name, because StarRocks.g4 gives it a
// production of its own (informationFunctionExpression,
// specialFunctionExpression). Without it `IF(a,1,2)`, `LEFT(s,2)` and
// `DATABASE()` all fail to parse even though the engine accepts them. TRIM,
// CONNECTION_ID, CURRENT_CATALOG and SESSION_USER are not keywords in
// StarRocksLex.g4 at all, so StarRocks reads them as plain function names.
func isFunctionNameKeyword(kind TokenKind) bool {
	switch kind {
	// informationFunctionExpression (StarRocks.g4:2477),
	// specialFunctionExpression (:2493), the GROUPING forms (:2442) and
	// windowFunction (:2521), restricted to the words StarRocks.g4's
	// nonReserved does not already make identifiers: CATALOG, DAY, HOUR,
	// MINUTE, MONTH, QUARTER, SECOND, YEAR, FLOOR, CEIL, RANK, CUME_DIST,
	// PERCENT_RANK and TIMESTAMPADD reach them as identifiers.
	case kwCHAR, kwCURRENT_ROLE, kwCURRENT_USER, kwDATABASE, kwDENSE_RANK,
		kwFIRST_VALUE, kwGROUPING, kwGROUPING_ID, kwIF, kwLAG, kwLAST_VALUE,
		kwLEAD, kwLEFT, kwLIKE, kwMOD, kwNTILE, kwPASSWORD, kwREGEXP,
		kwREPLACE, kwRIGHT, kwRLIKE, kwROW_NUMBER, kwSCHEMA, kwUSER:
		return true
	}
	return false
}

// isExprIdentToken reports whether the current token is usable as an identifier
// in expression context: tokIdent, tokQuotedIdent, or a non-reserved keyword.
func (p *Parser) isExprIdentToken() bool {
	switch p.cur.Kind {
	case tokIdent, tokQuotedIdent:
		return true
	default:
		return p.cur.Kind >= 700 && !IsReserved(p.cur.Kind)
	}
}

// ---------------------------------------------------------------------------
// Primary expression helpers
// ---------------------------------------------------------------------------

// parseArrayElementIdent reads ARRAY_ELEMENT — the lexer rule '[*]'
// (StarRocksLex.g4:506), which StarRocks.g4:2951 makes a nonReserved
// identifier — as a column name. The three characters must touch, as the
// lexer rule writes them together; otherwise '[ * ]' is an array constructor
// holding '*', which no expression rule derives.
func (p *Parser) parseArrayElementIdent() (ast.Node, bool) {
	if p.peekNext().Kind != int('*') {
		return nil, false
	}
	mark := p.save()
	open := p.advance() // '['
	star := p.advance() // '*'
	close := p.cur      // ']'
	if close.Kind == int(']') &&
		open.Loc.End == star.Loc.Start && star.Loc.End == close.Loc.Start {
		p.advance()
		loc := ast.Loc{Start: open.Loc.Start, End: close.Loc.End}
		return &ast.ColumnRef{
			Name: &ast.ObjectName{Parts: []string{"[*]"}, Loc: loc},
			Loc:  loc,
		}, true
	}
	p.restore(mark)
	return nil, false
}

// parseIdentExpr parses an identifier that could be a column ref or function
// call (name followed by '(').
func (p *Parser) parseIdentExpr() (ast.Node, error) {
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for function call: name(
	if p.cur.Kind == int('(') {
		return p.parseFuncCall(name)
	}

	// Simple or qualified column reference.
	col := &ast.ColumnRef{
		Name: name,
		Loc:  name.Loc,
	}
	// identifier '->' expression is a lambda wherever a primaryExpression
	// goes (StarRocks.g4:2423); '->' before a string is the JSON arrow.
	if len(name.Parts) == 1 && p.cur.Kind == tokArrow && p.peekNext().Kind != tokString {
		return p.parseLambdaExpr(col)
	}
	return col, nil
}

// finishMapLiteral parses the { key:val, ... } body of a map constructor. The
// optional map<k,v> type (nil for the untyped MAP{...} form) and the MAP keyword
// have already been consumed; the current token is '{'.
func (p *Parser) finishMapLiteral(start int, mapType *ast.TypeName) (ast.Node, error) {
	lit := &ast.MapLiteral{MapType: mapType}

	if _, err := p.expect(int('{')); err != nil {
		return nil, err
	}
	if p.cur.Kind != int('}') {
		entry, err := p.parseMapEntry()
		if err != nil {
			return nil, err
		}
		lit.Entries = append(lit.Entries, entry)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			entry, err = p.parseMapEntry()
			if err != nil {
				return nil, err
			}
			lit.Entries = append(lit.Entries, entry)
		}
	}
	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	lit.Loc = ast.Loc{Start: start, End: closeTok.Loc.End}
	return lit, nil
}

// parseMapEntry parses one key:value pair of a map constructor.
func (p *Parser) parseMapEntry() (*ast.MapEntry, error) {
	key, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(':')); err != nil {
		return nil, err
	}
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.MapEntry{
		Key:   key,
		Value: value,
		Loc:   ast.Loc{Start: ast.NodeLoc(key).Start, End: ast.NodeLoc(value).End},
	}, nil
}

// finishArrayLiteral parses the [ e, ... ] body of an array constructor. The
// optional array<t> type (nil for the untyped [...] form) has already been
// consumed; the current token is '['.
func (p *Parser) finishArrayLiteral(start int, elemType *ast.TypeName) (ast.Node, error) {
	lit := &ast.ArrayLiteral{ElemType: elemType}

	if _, err := p.expect(int('[')); err != nil {
		return nil, err
	}
	if p.cur.Kind != int(']') {
		el, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		lit.Elements = append(lit.Elements, el)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			el, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
			lit.Elements = append(lit.Elements, el)
		}
	}
	closeTok, err := p.expect(int(']'))
	if err != nil {
		return nil, err
	}
	lit.Loc = ast.Loc{Start: start, End: closeTok.Loc.End}
	return lit, nil
}

// parseFuncCall parses a function call starting at the opening '('.
// The function name has already been parsed as an ObjectName.
func (p *Parser) parseFuncCall(name *ast.ObjectName) (*ast.FuncCallExpr, error) {
	p.advance() // consume '('

	funcName := strings.ToUpper(name.Parts[len(name.Parts)-1])
	fc := &ast.FuncCallExpr{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	// COUNT(*) special handling
	if funcName == "COUNT" && p.cur.Kind == int('*') {
		p.advance()
		fc.Star = true
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeTok.Loc.End
		return fc, nil
	}

	// DISTINCT handling for aggregates like COUNT(DISTINCT x). ALL, the
	// default, is the other setQuantifier, which only the aggregates of
	// StarRocks.g4:2453-2462 take; all but COUNT then need an argument.
	needArg := false
	quantified := false
	if p.cur.Kind == kwDISTINCT {
		p.advance()
		fc.Distinct = true
		quantified = true
	} else if p.cur.Kind == kwALL && len(name.Parts) == 1 {
		switch funcName {
		case "COUNT":
			p.advance()
			quantified = true
		case "AVG", "MAX", "MIN", "SUM", "ARRAY_AGG", "GROUP_CONCAT":
			p.advance()
			needArg = true
		}
	}
	// COUNT '(' (setQuantifier bracketHint?)? ... (StarRocks.g4:2455): the
	// hint guides the planner and, like a /*+ */ hint, is not kept. A '['
	// that does not read as a hint opens an array argument (COUNT(DISTINCT
	// [])).
	if quantified && funcName == "COUNT" && len(name.Parts) == 1 && p.cur.Kind == int('[') {
		cp := p.save()
		if err := p.skipBracketHint(); err != nil {
			if isNestingError(err) {
				return nil, err
			}
			p.restore(cp)
		}
	}

	if fc.Distinct && p.cur.Kind == int(')') {
		// A setQuantifier quantifies arguments, and DISTINCT needs one:
		// `COUNT(DISTINCT)` is "Incorrect number of arguments in expr
		// 'count'" on the pinned server, while `COUNT()` and `COUNT(ALL)`
		// are accepted.
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Kind != int(')') || needArg {
		args, err := p.parseFuncArgList()
		if err != nil {
			return nil, err
		}
		fc.Args = args

		// Optional IGNORE NULLS after the FIRST argument (window null-treatment,
		// e.g. FIRST_VALUE(v IGNORE NULLS), LAG(v IGNORE NULLS, 1)). The grammar
		// attaches it only after the first expression, so we require exactly one
		// arg so far; further args may follow it.
		if len(fc.Args) == 1 && p.consumeIgnoreNulls() {
			fc.IgnoreNulls = true
			for p.cur.Kind == int(',') {
				p.advance() // consume ','
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				fc.Args = append(fc.Args, arg)
			}
		}

		// Optional ORDER BY within aggregate functions (e.g., GROUP_CONCAT)
		if p.cur.Kind == kwORDER {
			p.advance() // consume ORDER
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			orderItems, err := p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			fc.OrderBy = orderItems
		}

		// Optional SEPARATOR for GROUP_CONCAT
		// (SEPARATOR expression)? (StarRocks.g4:2462).
		if p.cur.Kind == kwSEPARATOR {
			p.advance() // consume SEPARATOR
			if p.cur.Kind == tokString && p.peekNext().Kind == int(')') {
				fc.Separator = p.advance().Str
			} else {
				sep, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				fc.SeparatorExpr = sep
			}
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	fc.Loc.End = closeTok.Loc.End

	// A named function's production fixes how many arguments it takes
	// (StarRocks.g4:2477-2542); the parser refuses the others.
	if err := p.checkFuncArity(name, fc); err != nil {
		return nil, err
	}

	// Optional IGNORE NULLS after the closing paren, e.g.
	// LEAD(v, 1) IGNORE NULLS OVER (...).
	if !fc.IgnoreNulls && p.consumeIgnoreNulls() {
		fc.IgnoreNulls = true
		fc.Loc.End = p.prev.Loc.End
	}

	// Optional OVER (...) window specification — turns this call into a window
	// function. Without consuming OVER here it is left dangling, which breaks
	// parsing wherever a specific following token is expected (e.g. the ')'
	// closing a CTE body: WITH x AS (SELECT f() OVER (...) ... ) ...).
	if p.cur.Kind == kwOVER {
		over, err := p.parseWindowSpec()
		if err != nil {
			return nil, err
		}
		fc.Over = over
		fc.Loc.End = over.Loc.End
	}

	// IGNORE NULLS is valid only on the four window value functions and requires
	// an OVER clause (StarRocks grammar: windowFunction is one of LEAD/LAG/
	// FIRST_VALUE/LAST_VALUE, used as `windowFunction over`). Reject otherwise.
	if fc.IgnoreNulls && (fc.Over == nil || !isIgnoreNullsWindowFunc(funcName)) {
		return nil, p.syntaxErrorAtCur()
	}

	return fc, nil
}

// funcArity is how many arguments the named function production of
// StarRocks.g4 takes, for the names that only that production reads:
// informationFunctionExpression (:2477), specialFunctionExpression (:2493),
// the GROUPING forms (:2442) and windowFunction (:2521). -1 as the maximum
// means unbounded.
//
// A name nonReserved makes an identifier is not here: it is also an ordinary
// function call, which takes any arguments, so the engine accepts `YEAR()`,
// `PASSWORD()` and `FLOOR()` while it refuses `CHAR()` and `LEFT()`
// (StarRocks.g4:2951 and the pinned server).
var funcArity = map[string][2]int{
	// informationFunctionExpression.
	"DATABASE": {0, 0}, "SCHEMA": {0, 0},
	"CURRENT_USER": {0, 0}, "CURRENT_ROLE": {0, 0},
	// specialFunctionExpression.
	"CHAR": {1, 1}, "LEFT": {2, 2}, "RIGHT": {2, 2}, "LIKE": {2, 2},
	"RLIKE": {2, 2}, "REGEXP": {2, 2}, "MOD": {2, 2},
	"IF": {0, -1}, "REPLACE": {0, -1},
	// groupingOperation.
	"GROUPING": {0, -1}, "GROUPING_ID": {0, -1},
	// windowFunction.
	"ROW_NUMBER": {0, 0}, "DENSE_RANK": {0, 0}, "NTILE": {0, 1},
	"LEAD": {0, -1}, "LAG": {0, -1}, "FIRST_VALUE": {0, -1},
	"LAST_VALUE": {0, -1},
}

// checkFuncArity refuses a call to a named function with an argument count its
// production does not take. Only an unqualified reserved name is held to it: a
// qualified call is an ordinary function, a non-reserved name is an identifier
// the grammar reads as one, and a backticked name is an identifier too, so
// `SELECT `LEFT`(1)` is an ordinary call and `SELECT LEFT(1)` is not.
func (p *Parser) checkFuncArity(name *ast.ObjectName, fc *ast.FuncCallExpr) error {
	if len(name.Parts) != 1 || fc.Star || p.isQuotedName(name) {
		return nil
	}
	upper := strings.ToUpper(name.Parts[0])
	arity, ok := funcArity[upper]
	if !ok {
		return nil
	}
	if kind, isKeyword := KeywordToken(upper); !isKeyword || !IsReserved(kind) {
		return nil
	}
	n := len(fc.Args)
	if n < arity[0] || arity[1] >= 0 && n > arity[1] {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// isQuotedName reports whether the name was written with backticks, which the
// token keeps only in the source it came from.
func (p *Parser) isQuotedName(name *ast.ObjectName) bool {
	start := name.Loc.Start - p.baseOffset
	return start >= 0 && start < len(p.input) && p.input[start] == '`'
}

// isIgnoreNullsWindowFunc reports whether name is one of the window value
// functions that accept an IGNORE NULLS null-treatment clause.
func isIgnoreNullsWindowFunc(name string) bool {
	switch name {
	case "FIRST_VALUE", "LAST_VALUE", "LEAD", "LAG":
		return true
	default:
		return false
	}
}

// consumeIgnoreNulls consumes an optional IGNORE NULLS null-treatment clause and
// reports whether it was present. (RESPECT NULLS is the default and does not
// exist as a keyword in StarRocks 3.4.)
func (p *Parser) consumeIgnoreNulls() bool {
	if p.cur.Kind == kwIGNORE && p.peekNext().Kind == kwNULLS {
		p.advance() // consume IGNORE
		p.advance() // consume NULLS
		return true
	}
	return false
}

// parseWindowSpec parses an OVER clause window specification. The OVER keyword
// is the current token (not yet consumed):
//
//	OVER ( [PARTITION BY expr, ...] [ORDER BY item, ...] [frame] )
//	OVER window_name
func (p *Parser) parseWindowSpec() (*ast.WindowSpec, error) {
	overTok := p.advance() // consume OVER
	spec := &ast.WindowSpec{Loc: ast.Loc{Start: overTok.Loc.Start}}

	// OVER window_name — a named-window reference (no parentheses).
	if p.cur.Kind != int('(') {
		name, nameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		spec.Name = name
		spec.Loc.End = nameLoc.End
		return spec, nil
	}

	p.advance() // consume '('

	// Optional PARTITION BY.
	if p.cur.Kind == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		spec.PartitionBy = exprs
	}

	// Optional ORDER BY.
	if p.cur.Kind == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		items, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		spec.OrderBy = items
	}

	// Optional frame clause.
	if p.cur.Kind == kwROWS || p.cur.Kind == kwRANGE {
		frame, err := p.parseWindowFrame()
		if err != nil {
			return nil, err
		}
		spec.Frame = frame
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	spec.Loc.End = closeTok.Loc.End
	return spec, nil
}

// parseWindowFrame parses a window frame clause. The ROWS/RANGE keyword is the
// current token (not yet consumed):
//
//	(ROWS | RANGE) frame_bound
//	(ROWS | RANGE) BETWEEN frame_bound AND frame_bound
func (p *Parser) parseWindowFrame() (*ast.WindowFrame, error) {
	unitTok := p.advance() // consume ROWS or RANGE
	frame := &ast.WindowFrame{
		Unit: strings.ToUpper(unitTok.Str),
		Loc:  ast.Loc{Start: unitTok.Loc.Start},
	}

	if p.cur.Kind == kwBETWEEN {
		p.advance()
		startType, startExpr, err := p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
		frame.StartType, frame.StartExpr = startType, startExpr
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		endType, endExpr, err := p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
		frame.EndType, frame.EndExpr = endType, endExpr
	} else {
		startType, startExpr, err := p.parseWindowFrameBound()
		if err != nil {
			return nil, err
		}
		frame.StartType, frame.StartExpr = startType, startExpr
	}

	frame.Loc.End = p.prev.Loc.End
	return frame, nil
}

// parseWindowFrameBound parses one window frame bound:
//
//	UNBOUNDED PRECEDING | UNBOUNDED FOLLOWING | CURRENT ROW
//	<expr> PRECEDING | <expr> FOLLOWING
//
// It returns the bound kind and, for the "<expr> PRECEDING/FOLLOWING" forms,
// the offset expression (nil otherwise).
func (p *Parser) parseWindowFrameBound() (string, ast.Node, error) {
	switch p.cur.Kind {
	case kwUNBOUNDED:
		p.advance()
		switch p.cur.Kind {
		case kwPRECEDING:
			p.advance()
			return "UNBOUNDED PRECEDING", nil, nil
		case kwFOLLOWING:
			p.advance()
			return "UNBOUNDED FOLLOWING", nil, nil
		default:
			return "", nil, &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "expected PRECEDING or FOLLOWING after UNBOUNDED"}
		}
	case kwCURRENT:
		p.advance()
		if _, err := p.expect(kwROW); err != nil {
			return "", nil, err
		}
		return "CURRENT ROW", nil, nil
	default:
		expr, err := p.parseExpr()
		if err != nil {
			return "", nil, err
		}
		switch p.cur.Kind {
		case kwPRECEDING:
			p.advance()
			return "PRECEDING", expr, nil
		case kwFOLLOWING:
			p.advance()
			return "FOLLOWING", expr, nil
		default:
			return "", nil, &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "expected PRECEDING or FOLLOWING in window frame bound"}
		}
	}
}

// tryParseParenLambda speculatively parses the parenthesized multi-parameter
// lambda form `(x, y) -> body`, rolling the parser back and reporting ok=false
// when the input turns out to be an ordinary parenthesized expression.
//
// The grammar requires at least two parameters in this form; a single
// parameter must be written bare (`x -> ...`), and the engine rejects
// `(x) -> ...`. That is why len(params) < 2 backtracks rather than accepting.
func (p *Parser) tryParseParenLambda() (ast.Node, bool, error) {
	if p.cur.Kind != int('(') {
		return nil, false, nil
	}
	saved := p.save()
	start := p.cur.Loc.Start
	p.advance() // consume '('

	var params []string
	for {
		if !p.isExprIdentToken() {
			p.restore(saved)
			return nil, false, nil
		}
		params = append(params, p.advance().Str)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if len(params) < 1 || p.cur.Kind != int(')') {
		p.restore(saved)
		return nil, false, nil
	}
	p.advance() // consume ')'
	if p.cur.Kind != tokArrow {
		p.restore(saved)
		return nil, false, nil
	}
	// A '->' followed by a string literal is the JSON arrow, which the engine
	// reads even when the left side is a one-name identifier list:
	// `array_map((x) -> '$.a', [1,2])` is the arrow (see parseJSONArrow).
	if len(params) == 1 && p.peekNext().Kind == tokString {
		p.restore(saved)
		return nil, false, nil
	}
	p.advance() // consume '->'

	// identifierList '->' '(' (expressionList)? ')' (StarRocks.g4:2424) takes an
	// empty body, which the pinned server accepts (`SELECT (a) -> ()`) and
	// omni keeps as a row of no expression.
	if p.cur.Kind == int('(') && p.peekNext().Kind == int(')') {
		openTok := p.advance() // consume '('
		closeTok := p.advance()
		empty := &ast.RowExpr{Loc: ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}}
		return &ast.LambdaExpr{
			Params: params,
			Body:   empty,
			Loc:    ast.Loc{Start: start, End: closeTok.Loc.End},
		}, true, nil
	}

	// Past the arrow the shape is unambiguous, so a bad body is a real error
	// rather than a reason to backtrack.
	body, err := p.parseExpr()
	if err != nil {
		return nil, false, err
	}
	return &ast.LambdaExpr{
		Params: params,
		Body:   body,
		Loc:    ast.Loc{Start: start, End: ast.NodeLoc(body).End},
	}, true, nil
}

// parseFuncArgList parses a function-call argument list, which may hold a
// lambda (`x -> body`, `(x, y) -> body`). parseExpr reads a lambda wherever
// a primaryExpression goes (StarRocks.g4:2423); the checks here remain for an
// argument whose left of '->' is not a bare name or name list, which
// parseLambdaExpr then refuses.
//
// A `->` followed by a string literal has already been taken as the JSON path
// operator by then (see parseJSONArrow), which is how the engine resolves that
// overlap too.
func (p *Parser) parseFuncArgList() ([]ast.Node, error) {
	var list []ast.Node
	for {
		lam, ok, err := p.tryParseParenLambda()
		if err != nil {
			return nil, err
		}
		if !ok {
			// A bare-name lambda is read here rather than inside parseExpr,
			// so a lambda argument costs one nesting level, not two.
			lam, ok, err = p.tryParseNameLambda()
			if err != nil {
				return nil, err
			}
		}
		if !ok {
			lam, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
			if p.cur.Kind == tokArrow {
				lam, err = p.parseLambdaExpr(lam)
				if err != nil {
					return nil, err
				}
			}
		}
		list = append(list, lam)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() ([]ast.Node, error) {
	var list []ast.Node
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
}

// parseOrderByList parses a comma-separated list of ORDER BY items.
func (p *Parser) parseOrderByList() ([]*ast.OrderByItem, error) {
	var items []*ast.OrderByItem
	for {
		item, err := p.parseOrderByItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return items, nil
}

// parseOrderByItem parses one ORDER BY item: expr [ASC|DESC] [NULLS FIRST|LAST].
func (p *Parser) parseOrderByItem() (*ast.OrderByItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	item := &ast.OrderByItem{
		Expr: expr,
		Loc:  ast.NodeLoc(expr),
	}

	// Optional ASC/DESC
	if p.cur.Kind == kwASC {
		p.advance()
		item.Desc = false
		item.Loc.End = p.prev.Loc.End
	} else if p.cur.Kind == kwDESC {
		p.advance()
		item.Desc = true
		item.Loc.End = p.prev.Loc.End
	}

	// Optional NULLS FIRST/LAST
	if p.cur.Kind == kwNULLS {
		p.advance() // consume NULLS
		switch p.cur.Kind {
		case kwFIRST:
			p.advance()
			b := true
			item.NullsFirst = &b
			item.Loc.End = p.prev.Loc.End
		case kwLAST:
			p.advance()
			b := false
			item.NullsFirst = &b
			item.Loc.End = p.prev.Loc.End
		default:
			return nil, &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "expected FIRST or LAST after NULLS",
			}
		}
	}

	return item, nil
}

// parseParenExprOrSubquery parses a parenthesized expression or subquery.
// For subqueries (SELECT/WITH after open paren), we create a placeholder
// SubqueryExpr that stores the raw text.
func (p *Parser) parseParenExprOrSubquery() (ast.Node, error) {
	// identifierList '->' expression (StarRocks.g4:2423), the lambda form
	// tryParseParenLambda reads.
	if lam, ok, err := p.tryParseParenLambda(); err != nil || ok {
		return lam, err
	}
	openTok := p.advance() // consume '('

	// Subquery: (SELECT ...) or (WITH ... SELECT ...)
	if p.cur.Kind == kwSELECT || p.cur.Kind == kwWITH {
		return p.parseSubqueryPlaceholder(openTok.Loc.Start)
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Kind == int(',') {
		return p.parseTupleInSubquery(openTok.Loc.Start, expr)
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	return &ast.ParenExpr{
		Expr: expr,
		Loc:  ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseSubqueryPlaceholder consumes tokens until the matching closing ')' for
// a subquery expression and returns a SubqueryExpr with the raw text.
// startOffset is the byte offset of the opening '(' token.
//
// Strict mode validates every body (validateRawQueries). It first parses the
// body in place: a body that parses as a query ending at the matching ')',
// recording no error on the way, is valid and is not parsed again. Validating
// it afterwards instead reparsed the text of every nested body once per
// level, so n nested subqueries cost O(n²). A body that fails here is scanned
// raw and validated afterwards as before, so its errors are unchanged.
func (p *Parser) parseSubqueryPlaceholder(startOffset int) (*ast.SubqueryExpr, error) {
	subStart := p.cur.Loc.Start
	if p.strictTrailing {
		cp := p.save()
		n, err := p.parseStmt()
		if err == nil && isQueryBody(n) && p.cur.Kind == int(')') &&
			len(p.errors) == cp.errLen && len(p.lexer.errors) == cp.lexErrLen {
			sub := p.finishSubqueryPlaceholder(startOffset, subStart)
			sub.Query = n
			return sub, nil
		}
		p.restore(cp)
	}
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		if p.cur.Kind == int('(') {
			p.skipParenGroup()
			if p.cur.Kind == tokEOF {
				break
			}
		}
		p.advance()
	}
	p.recordParenGroupEnd(startOffset)
	if p.cur.Kind == tokEOF {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "unterminated subquery",
		}
	}
	sub := p.finishSubqueryPlaceholder(startOffset, subStart)
	p.rawQueries = append(p.rawQueries, rawQuery{text: sub.RawText, start: sub.TextStart, depth: p.depth, sub: sub})
	return sub, nil
}

// skipParenGroup moves from the '(' at p.cur to its matching ')', left as
// p.cur, or to EOF when it never closes. A group scanned before is jumped
// over (see parenGroupEnd).
//
// The scan keeps the groups it is inside on a stack of its own rather than
// recursing per nested group, which overflowed the Go stack on a few hundred
// kilobytes of '('. It records the same group ends in the same order.
func (p *Parser) skipParenGroup() {
	var open []int // the '(' offset of each group being scanned, innermost last
	for {
		// p.cur is the '(' of a group to skip: jump over it when it was
		// scanned before, enter it otherwise.
		if r, ok := p.parenGroupEnd[p.cur.Loc.Start]; ok {
			p.replayParen(r)
		} else {
			open = append(open, p.cur.Loc.Start)
			p.advance() // consume '('
			p.advanceToParen()
			if p.cur.Kind == int('(') {
				continue
			}
			p.recordParenGroupEnd(open[len(open)-1])
			open = open[:len(open)-1]
		}
		// The group skipped last ended at p.cur, its ')' or EOF. Resume the
		// groups around it, innermost first, until one meets another '('.
		for len(open) > 0 {
			if p.cur.Kind != tokEOF {
				p.advance() // consume the ')' that closed the inner group
				p.advanceToParen()
				if p.cur.Kind == int('(') {
					break
				}
			}
			p.recordParenGroupEnd(open[len(open)-1])
			open = open[:len(open)-1]
		}
		if len(open) == 0 {
			return
		}
	}
}

// advanceToParen moves to the next '(' or ')', or to EOF.
func (p *Parser) advanceToParen() {
	for p.cur.Kind != int('(') && p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		p.advance()
	}
}

// recordParenGroupEnd records the current token, the matching ')' or EOF, as
// the end of the '(' group at absolute offset open.
func (p *Parser) recordParenGroupEnd(open int) {
	if p.parenGroupEnd == nil {
		p.parenGroupEnd = make(map[int]parenReplay)
	}
	p.parenGroupEnd[open] = p.recordParenReplay(open)
}

// finishSubqueryPlaceholder consumes the ')' that closes a subquery body
// starting at subStart and returns its placeholder.
func (p *Parser) finishSubqueryPlaceholder(startOffset, subStart int) *ast.SubqueryExpr {
	rawText := p.input[subStart-p.baseOffset : p.cur.Loc.Start-p.baseOffset]
	closeTok := p.advance() // consume ')'

	trimmedFromStart := len(rawText) - len(strings.TrimLeft(rawText, " \t\r\n"))
	return &ast.SubqueryExpr{
		RawText:   strings.TrimSpace(rawText),
		TextStart: subStart + trimmedFromStart,
		Loc:       ast.Loc{Start: startOffset, End: closeTok.Loc.End},
	}
}

// parseCaseExpr parses a CASE expression (simple or searched).
func (p *Parser) parseCaseExpr() (*ast.CaseExpr, error) {
	caseTok := p.advance() // consume CASE

	ce := &ast.CaseExpr{
		Loc: ast.Loc{Start: caseTok.Loc.Start},
	}

	// Determine if this is a simple CASE (CASE expr WHEN ...) or
	// searched CASE (CASE WHEN ...).
	if p.cur.Kind != kwWHEN {
		// Simple CASE: parse the operand expression.
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Kind = ast.CaseSimple
		ce.Operand = operand
	} else {
		ce.Kind = ast.CaseSearched
	}

	// Parse WHEN clauses.
	for p.cur.Kind == kwWHEN {
		whenTok := p.advance() // consume WHEN
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		result, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Whens = append(ce.Whens, &ast.WhenClause{
			Cond:   cond,
			Result: result,
			Loc:    ast.Loc{Start: whenTok.Loc.Start, End: ast.NodeLoc(result).End},
		})
	}

	if len(ce.Whens) == 0 {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "expected WHEN in CASE expression",
		}
	}

	// Optional ELSE.
	if p.cur.Kind == kwELSE {
		p.advance() // consume ELSE
		elseExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Else = elseExpr
	}

	endTok, err := p.expect(kwEND)
	if err != nil {
		return nil, err
	}
	ce.Loc.End = endTok.Loc.End

	return ce, nil
}

// parseCastExpr parses CAST(expr AS type). StarRocksLex.g4 has no TRY_CAST
// token, and the 3.4.10 server refuses TRY_CAST(x AS t) at the form's AS: the
// name is an ordinary function there, so `TRY_CAST(x)` parses and this does
// not.
func (p *Parser) parseCastExpr() (*ast.CastExpr, error) {
	castTok := p.advance() // consume CAST
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	typeName, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ast.CastExpr{
		Expr:     expr,
		TypeName: typeName,
		Loc:      ast.Loc{Start: castTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseExtractExpr parses EXTRACT(unit FROM expr).
//
// EXTRACT is a reserved keyword, so without a production for it the whole
// statement fails to parse at the EXTRACT token.
func (p *Parser) parseExtractExpr() (ast.Node, error) {
	extractTok := p.advance() // consume EXTRACT
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	unit, err := p.parseExtractUnit()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ast.ExtractExpr{
		Expr: expr,
		Unit: unit,
		Loc:  ast.Loc{Start: extractTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseExtractUnit parses the unit of EXTRACT(unit FROM expr).
//
// The grammar takes an identifier, which covers plain identifiers (DAYOFWEEK,
// DOY, ...) and the non-reserved unit keywords (YEAR, MONTH, WEEK, ...).
// Doris's compound units DAY_HOUR, DAY_SECOND and MINUTE_SECOND are keywords
// here only in Doris: StarRocks.g4::2776 names no such unit, and the engine
// reads each as an identifier, which this default takes.
func (p *Parser) parseExtractUnit() (string, error) {
	if p.isExprIdentToken() {
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	}
	return "", &ParseError{
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: fmt.Sprintf("expected EXTRACT unit, got %q", p.cur.Str),
	}
}

// parseLambdaExpr parses `params -> body`, given the already-parsed expression
// to the left of the arrow.
func (p *Parser) parseLambdaExpr(left ast.Node) (ast.Node, error) {
	arrowTok := p.advance() // consume '->'

	params, ok := lambdaParams(left)
	if !ok {
		return nil, &ParseError{
			Position: arrowTok.Loc.Start, End: arrowTok.Loc.End,
			Message: "invalid lambda parameter list before '->'",
		}
	}
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.LambdaExpr{
		Params: params,
		Body:   body,
		Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(body).End},
	}, nil
}

// lambdaParams extracts parameter names from the expression left of '->'. The
// parameters parse as ordinary column references (`x`, or `(x, y)` as a
// parenthesized row), so they are reinterpreted here rather than lexed
// specially.
func lambdaParams(n ast.Node) ([]string, bool) {
	switch v := n.(type) {
	case *ast.ColumnRef:
		if v.Name == nil || len(v.Name.Parts) != 1 {
			return nil, false
		}
		return []string{v.Name.Parts[0]}, true
	case *ast.ParenExpr:
		// `(x) -> ...`. Unlike Doris, this engine accepts the parenthesized
		// form with a single parameter as well as the bare one.
		return lambdaParams(v.Expr)
	}
	return nil, false
}

// parseGroupingElementCall parses a reserved-keyword grouping element such as
// CUBE(a, b) as a function call, keeping GROUP BY items uniform with ROLLUP and
// GROUPING SETS (both non-reserved, so they already take the expression path).
func (p *Parser) parseGroupingElementCall() (*ast.FuncCallExpr, error) {
	tok := p.advance()
	name := &ast.ObjectName{Parts: []string{strings.ToUpper(tok.Str)}, Loc: tok.Loc}
	return p.parseFuncCall(name)
}

// parseConvertExpr parses CONVERT(expr, type), which is CAST spelled
// differently and so becomes a CastExpr. The transcoding form
// CONVERT(expr USING charset) is not accepted by this engine.
func (p *Parser) parseConvertExpr() (ast.Node, error) {
	convertTok := p.advance() // consume CONVERT
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(int(',')); err != nil {
		return nil, err
	}
	typeName, err := p.parseConvertTargetType()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ast.CastExpr{
		Expr:     expr,
		TypeName: typeName,
		Loc:      ast.Loc{Start: convertTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseConvertTargetType parses the target type of CONVERT(expr, type). It
// additionally accepts SIGNED / UNSIGNED [INTEGER], which are cast-only type
// names rather than column types and so are unknown to parseDataType.
func (p *Parser) parseConvertTargetType() (*ast.TypeName, error) {
	switch p.cur.Kind {
	case kwSIGNED, kwUNSIGNED:
		tok := p.advance()
		tn := &ast.TypeName{Name: strings.ToUpper(tok.Str), Loc: tok.Loc}
		if p.cur.Kind == kwINTEGER || p.cur.Kind == kwINT {
			end := p.advance()
			tn.Loc.End = end.Loc.End
		}
		return tn, nil
	}
	return p.parseDataType()
}

// parseVariableRef parses a system variable (@@name, @@scope.name) or a
// user-defined variable (@name).
func (p *Parser) parseVariableRef(system bool) (ast.Node, error) {
	atTok := p.advance() // consume '@@' or '@'

	// A user variable name may also be a quoted string — @'name' / @"name"
	// (grammar: AT identifierOrText; engine-verified accept).
	if !system && p.cur.Kind == tokString {
		tok := p.advance()
		return &ast.VariableRef{
			System: false,
			Name:   tok.Str,
			Loc:    ast.Loc{Start: atTok.Loc.Start, End: tok.Loc.End},
		}, nil
	}

	// systemVariable is AT AT (varType '.')? identifier (StarRocks.g4:2469),
	// and identifier accepts only non-reserved keywords. A reserved word is
	// not a variable name, so @@GRANT.a is refused.
	if !isIdentifierToken(p.cur.Kind) {
		return nil, p.syntaxErrorAtCur()
	}
	first := p.advance()
	v := &ast.VariableRef{
		System: system,
		Name:   first.Str,
		Loc:    ast.Loc{Start: atTok.Loc.Start, End: first.Loc.End},
	}
	// @@global.x / @@local.x / @@session.x — varType is the only scope.
	if system && p.cur.Kind == int('.') &&
		(first.Kind == kwGLOBAL || first.Kind == kwLOCAL || first.Kind == kwSESSION) {
		p.advance() // consume '.'
		if !isIdentifierToken(p.cur.Kind) {
			return nil, p.syntaxErrorAtCur()
		}
		second := p.advance()
		v.Scope = strings.ToUpper(v.Name)
		v.Name = second.Str
		v.Loc.End = second.Loc.End
	}
	return v, nil
}

// identOrKeywordToken accepts any identifier- or keyword-shaped token, for
// positions where there is no ambiguity: after an @ / @@ prefix, or after
// USING. Scope names (SESSION, GLOBAL), charset names and many variable names
// are themselves keywords.
func (p *Parser) identOrKeywordToken() (Token, bool) {
	switch {
	case p.cur.Kind == tokIdent || p.cur.Kind == tokQuotedIdent:
		return p.advance(), true
	case p.cur.Kind >= 700: // any keyword
		return p.advance(), true
	}
	return Token{}, false
}

// parseNilaryFunction parses zero-argument keyword functions like
// CURRENT_DATE, CURRENT_TIMESTAMP, etc.
func (p *Parser) parseNilaryFunction() (ast.Node, error) {
	tok := p.advance()
	name := strings.ToUpper(tok.Str)
	objName := &ast.ObjectName{
		Parts: []string{name},
		Loc:   tok.Loc,
	}

	// Some of these may optionally have empty parens: CURRENT_TIMESTAMP()
	fc := &ast.FuncCallExpr{
		Name: objName,
		Loc:  tok.Loc,
	}
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		// Only CURRENT_TIMESTAMP takes a fraction of a second of precision,
		// and only an integer: `CURRENT_TIMESTAMP ( '(' (INTEGER_VALUE)? ')' )?`
		// (StarRocks.g4:2540). The pinned server refuses CURRENT_TIMESTAMP(zzz),
		// and CURRENT_TIME(3), LOCALTIME(3), CURRENT_DATE(1) and
		// CURRENT_USER(1, 2) alike.
		if p.cur.Kind != int(')') {
			if name != "CURRENT_TIMESTAMP" || p.cur.Kind != tokInt {
				return nil, p.syntaxErrorAtCur()
			}
			arg := p.advance()
			fc.Args = []ast.Node{&ast.Literal{Kind: ast.LitInt, Value: arg.Str, Loc: arg.Loc}}
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		fc.Loc.End = closeTok.Loc.End
	}
	return fc, nil
}

// parseExistsExpr parses EXISTS (subquery).
func (p *Parser) parseExistsExpr() (ast.Node, error) {
	existsTok := p.advance() // consume EXISTS
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	subq, err := p.parseSubqueryPlaceholder(existsTok.Loc.Start)
	if err != nil {
		return nil, err
	}

	return &ast.ExistsExpr{
		Subquery: subq,
		Loc:      ast.Loc{Start: existsTok.Loc.Start, End: subq.Loc.End},
	}, nil
}

// parseIntervalExpr parses INTERVAL expr unit.
func (p *Parser) parseIntervalExpr() (ast.Node, error) {
	intervalTok := p.advance() // consume INTERVAL

	value, err := p.parseExprPrec(bpPrimary)
	if err != nil {
		return nil, err
	}

	// Parse interval unit keyword
	unit, err := p.parseIntervalUnit()
	if err != nil {
		return nil, err
	}

	return &ast.IntervalExpr{
		Value: value,
		Unit:  unit,
		Loc:   ast.Loc{Start: intervalTok.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseTimestampUnitCall parses TIMESTAMPADD(unit, expr, expr) and
// TIMESTAMPDIFF(unit, expr, expr) (StarRocks.g4:2510-2511). The first argument
// is a unitIdentifier and not an expression, which matters now that
// MICROSECOND and MILLISECOND are the reserved keywords the grammar makes
// them: an argument list cannot read one.
func (p *Parser) parseTimestampUnitCall() (ast.Node, error) {
	tok := p.advance() // consume TIMESTAMPADD or TIMESTAMPDIFF
	fc := &ast.FuncCallExpr{
		Name: &ast.ObjectName{Parts: []string{strings.ToUpper(tok.Str)}, Loc: tok.Loc},
		Loc:  tok.Loc,
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	unitTok := p.cur
	unit, err := p.parseIntervalUnit()
	if err != nil {
		return nil, err
	}
	fc.Args = append(fc.Args, &ast.ColumnRef{
		Name: &ast.ObjectName{Parts: []string{unit}, Loc: unitTok.Loc},
		Loc:  unitTok.Loc,
	})
	for range 2 {
		if _, err := p.expect(int(',')); err != nil {
			return nil, err
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		fc.Args = append(fc.Args, arg)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	fc.Loc.End = closeTok.Loc.End
	return fc, nil
}

// parseIntervalUnit parses and returns the interval unit as a string: exactly
// the tokens of unitIdentifier (StarRocks.g4:2776). An identifier is not a
// unit, and the engine refuses one: INTERVAL 1 DAY_HOUR is a syntax error at
// DAY_HOUR.
func (p *Parser) parseIntervalUnit() (string, error) {
	switch p.cur.Kind {
	case kwYEAR, kwMONTH, kwWEEK, kwDAY, kwHOUR, kwMINUTE, kwSECOND, kwQUARTER,
		kwMILLISECOND, kwMICROSECOND:
		tok := p.advance()
		return strings.ToUpper(tok.Str), nil
	default:
		return "", &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: fmt.Sprintf("expected interval unit, got %q", p.cur.Str),
		}
	}
}

// ---------------------------------------------------------------------------
// Infix/postfix handlers
// ---------------------------------------------------------------------------

// parseIsExpr parses IS [NOT] NULL/TRUE/FALSE.
func (p *Parser) parseIsExpr(left ast.Node) (ast.Node, error) {
	p.advance() // consume IS

	notFlag := false
	if p.cur.Kind == kwNOT {
		p.advance()
		notFlag = true
	}

	var isWhat string
	switch p.cur.Kind {
	case kwNULL:
		isWhat = "NULL"
	case kwTRUE:
		isWhat = "TRUE"
	case kwFALSE:
		isWhat = "FALSE"
	default:
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "expected NULL, TRUE, or FALSE after IS [NOT]",
		}
	}
	endTok := p.advance()

	return &ast.IsExpr{
		Expr:   left,
		Not:    notFlag,
		IsWhat: isWhat,
		Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: endTok.Loc.End},
	}, nil
}

// parseBetweenExpr parses [NOT] BETWEEN low AND high.
func (p *Parser) parseBetweenExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume BETWEEN

	// Parse low bound at bpPred+1 to avoid re-consuming AND as an
	// infix binary operator.
	low, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAND); err != nil {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "expected AND in BETWEEN expression",
		}
	}

	high, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	return &ast.BetweenExpr{
		Expr: left,
		Low:  low,
		High: high,
		Not:  notFlag,
		Loc:  ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(high).End},
	}, nil
}

// parseInExpr parses [NOT] IN (values...) or [NOT] IN (subquery).
func (p *Parser) parseInExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume IN

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Subquery: IN (SELECT ...) or IN (WITH ... SELECT ...)
	if p.cur.Kind == kwSELECT || p.cur.Kind == kwWITH {
		subq, err := p.parseSubqueryPlaceholder(ast.NodeLoc(left).Start)
		if err != nil {
			return nil, err
		}
		return &ast.InExpr{
			Expr:   left,
			Values: []ast.Node{subq},
			Not:    notFlag,
			Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: subq.Loc.End},
		}, nil
	}

	values, err := p.parseExprList()
	if err != nil {
		return nil, err
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}

	return &ast.InExpr{
		Expr:   left,
		Values: values,
		Not:    notFlag,
		Loc:    ast.Loc{Start: ast.NodeLoc(left).Start, End: closeTok.Loc.End},
	}, nil
}

// parseLikeExpr parses [NOT] LIKE pattern [ESCAPE esc].
func (p *Parser) parseLikeExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume LIKE

	pattern, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	le := &ast.LikeExpr{
		Expr:    left,
		Pattern: pattern,
		Not:     notFlag,
		Loc:     ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(pattern).End},
	}

	// Optional ESCAPE clause
	if p.cur.Kind == kwESCAPE {
		p.advance() // consume ESCAPE
		esc, err := p.parseExprPrec(bpPred + 1)
		if err != nil {
			return nil, err
		}
		le.Escape = esc
		le.Loc.End = ast.NodeLoc(esc).End
	}

	return le, nil
}

// parseRegexpExpr parses [NOT] REGEXP|RLIKE pattern.
func (p *Parser) parseRegexpExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	p.advance() // consume REGEXP or RLIKE

	pattern, err := p.parseExprPrec(bpPred + 1)
	if err != nil {
		return nil, err
	}

	return &ast.RegexpExpr{
		Expr:    left,
		Pattern: pattern,
		Not:     notFlag,
		Loc:     ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(pattern).End},
	}, nil
}

// parseCollateExpr reads COLLATE (identifier | string) after left (StarRocks.g4:2403).
func (p *Parser) parseCollateExpr(left ast.Node) (ast.Node, error) {
	p.advance() // consume COLLATE
	var name string
	switch {
	case p.cur.Kind == tokString:
		name = p.cur.Str
		p.advance()
	default:
		ident, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		name = ident
	}
	return &ast.CollateExpr{
		Expr:      left,
		Collation: name,
		Loc:       ast.Loc{Start: ast.NodeLoc(left).Start, End: p.prev.Loc.End},
	}, nil
}

// parseMatchExpr reads [NOT] MATCH primaryExpression after left (StarRocks.g4:2425).
func (p *Parser) parseMatchExpr(left ast.Node, notFlag bool) (ast.Node, error) {
	if notFlag {
		p.advance() // consume NOT
	}
	op := strings.ToUpper(p.advance().Str) // consume MATCH
	pattern, err := p.parseMatchOperand()
	if err != nil {
		return nil, err
	}
	return &ast.MatchExpr{
		Expr:    left,
		Op:      op,
		Pattern: pattern,
		Not:     notFlag,
		Loc:     ast.Loc{Start: ast.NodeLoc(left).Start, End: ast.NodeLoc(pattern).End},
	}, nil
}

// parseMatchOperand reads the primaryExpression right of MATCH, with the
// subscripts and COLLATE that bind as tightly.
func (p *Parser) parseMatchOperand() (ast.Node, error) {
	// A primaryExpression may open with unary operators (StarRocks.g4:2408-
	// 2409); a loop reads them, so a long run of them cannot deepen the stack.
	type prefix struct {
		op    ast.UnaryOp
		start int
	}
	var prefixes []prefix
	for {
		var op ast.UnaryOp
		switch p.cur.Kind {
		case int('-'):
			op = ast.UnaryMinus
		case int('+'):
			op = ast.UnaryPlus
		case int('~'):
			op = ast.UnaryBitNot
		case int('!'):
			op = ast.UnaryNot
		default:
			goto operand
		}
		prefixes = append(prefixes, prefix{op, p.advance().Loc.Start})
	}
operand:
	node, err := p.parseMatchPrimary()
	if err != nil {
		return nil, err
	}
	for i := len(prefixes) - 1; i >= 0; i-- {
		node = &ast.UnaryExpr{
			Op:   prefixes[i].op,
			Expr: node,
			Loc:  ast.Loc{Start: prefixes[i].start, End: ast.NodeLoc(node).End},
		}
	}
	return node, nil
}

// parseMatchPrimary reads a primary expression with the suffixes that bind
// as tightly.
func (p *Parser) parseMatchPrimary() (ast.Node, error) {
	node, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}
	node, err = p.parsePostfixSuffixes(node)
	if err != nil {
		return nil, err
	}
	for p.cur.Kind == kwCOLLATE {
		if node, err = p.parseCollateExpr(node); err != nil {
			return nil, err
		}
		if node, err = p.parsePostfixSuffixes(node); err != nil {
			return nil, err
		}
	}
	return node, nil
}

// parseTupleInSubquery reads the rest of tupleInSubquery (StarRocks.g4:2367),
// '(' expression (',' expression)+ ')' NOT? IN '(' queryRelation ')', after its
// first expression; cur is the ','. A tuple is only the left of this
// predicate, so it is read whole here.
func (p *Parser) parseTupleInSubquery(start int, first ast.Node) (ast.Node, error) {
	row := &ast.RowExpr{Items: []ast.Node{first}}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		item, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		row.Items = append(row.Items, item)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	row.Loc = ast.Loc{Start: start, End: closeTok.Loc.End}

	not := false
	if p.cur.Kind == kwNOT {
		p.advance() // consume NOT
		not = true
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	open, err := p.expect(int('('))
	if err != nil {
		return nil, err
	}
	if p.cur.Kind != kwSELECT && p.cur.Kind != kwWITH {
		return nil, p.syntaxErrorAtCur()
	}
	sub, err := p.parseSubqueryPlaceholder(open.Loc.Start)
	if err != nil {
		return nil, err
	}
	return &ast.InExpr{
		Expr:   row,
		Values: []ast.Node{sub},
		Not:    not,
		Loc:    ast.Loc{Start: start, End: sub.Loc.End},
	}, nil
}

// tryParseNameLambda reads identifier '->' expression (StarRocks.g4:2423) at
// the current token, or leaves the parser where it was and reports ok=false.
func (p *Parser) tryParseNameLambda() (ast.Node, bool, error) {
	if !p.isExprIdentToken() || p.peekNext().Kind != tokArrow {
		return nil, false, nil
	}
	saved := p.save()
	tok := p.advance()
	if p.peekNext().Kind == tokString {
		p.restore(saved)
		return nil, false, nil
	}
	name := &ast.ObjectName{Parts: []string{tok.Str}, Loc: tok.Loc}
	lam, err := p.parseLambdaExpr(&ast.ColumnRef{Name: name, Loc: tok.Loc})
	return lam, err == nil, err
}

// skipBracketHint reads bracketHint (StarRocks.g4:2249):
// '[' identifier (',' identifier)* ']' or
// '[' identifier '|' primaryExpression literalExpressionList ']'.
func (p *Parser) skipBracketHint() error {
	// The hint's primaryExpression can hold another COUNT with a hint, so
	// this is a recursion cycle and passes through the nesting limit.
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	p.advance() // consume '['
	if _, _, err := p.parseIdentifier(); err != nil {
		return err
	}
	if p.cur.Kind == int('|') {
		p.advance() // consume '|'
		if _, err := p.parsePrimaryExpr(); err != nil {
			return err
		}
		if _, err := p.expect(int('(')); err != nil {
			return err
		}
		for {
			if _, err := p.parseExpr(); err != nil {
				return err
			}
			if p.cur.Kind != int(',') {
				break
			}
			p.advance() // consume ','
		}
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
	} else {
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			if _, _, err := p.parseIdentifier(); err != nil {
				return err
			}
		}
	}
	_, err := p.expect(int(']'))
	return err
}

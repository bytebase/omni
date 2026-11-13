package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// Precedence levels for Pratt parsing (low to high).
//
// The levels follow the rules in sql_yacc.yy, whose shape a flat precedence
// table cannot say on its own:
//
//	expr:      expr OR expr | expr XOR expr | expr AND expr | NOT expr
//	         | bool_pri IS [NOT] {TRUE|FALSE|UNKNOWN} | bool_pri
//	bool_pri:  bool_pri IS [NOT] NULL | bool_pri comp_op predicate
//	         | bool_pri comp_op {ANY|ALL} (subquery) | predicate
//	predicate: bit_expr [NOT] IN (...) | bit_expr [NOT] BETWEEN bit_expr AND predicate
//	         | bit_expr SOUNDS LIKE bit_expr | bit_expr [NOT] LIKE simple_expr
//	         | bit_expr [NOT] REGEXP bit_expr | bit_expr MEMBER [OF] (simple_expr)
//	         | bit_expr
//
// So a predicate binds tighter than a comparison (`a = b IN (1)` is `a = (b IN
// (1))`), its left operand is a bit_expr and so never another predicate, IS TRUE
// is an expr and takes no comparison after it, and NOT reaches only as far down
// as an expr goes: `1 + NOT 2` is a syntax error and `NOT a = b` is `NOT (a = b)`.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/operator-precedence.html
const (
	precNone       = 0
	precAssign     = 1  // the lowest level: a whole expression
	precOr         = 2  // OR, ||
	precXor        = 3  // XOR
	precAnd        = 4  // AND, &&
	precNot        = 5  // NOT (prefix)
	precComparison = 6  // =, <=>, >=, >, <=, <, <>, !=, IS
	precPredicate  = 7  // [NOT] IN, BETWEEN, LIKE, REGEXP, SOUNDS LIKE, MEMBER OF
	precBitOr      = 8  // | (the bottom of a bit_expr)
	precBitAnd     = 9  // &
	precShift      = 10 // <<, >>
	precAdd        = 11 // +, -
	precMul        = 12 // *, /, DIV, %, MOD
	precBitXor     = 13 // ^
	precUnary      = 14 // -, ~, !
	precCollate    = 15 // COLLATE (the bottom of a simple_expr)
	precJsonAccess = 16 // ->, ->>
)

// parseExpr parses an expression using Pratt parsing / precedence climbing.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/expressions.html
func (p *Parser) parseExpr() (nodes.ExprNode, error) {
	return p.parseExprPrec(precAssign)
}

// parseBitExpr reads a bit_expr: arithmetic, bitwise and unary operators over
// primaries, and none of NOT, AND, OR, XOR, a comparison, IS, IN, BETWEEN, LIKE or
// REGEXP outside parentheses. A partition expression and a partition bound are
// each one, which is why `PARTITION BY HASH (a OR 1)` and `VALUES LESS THAN (NOT
// 1)` are 1064 while `HASH ((a OR 1))` parses.
func (p *Parser) parseBitExpr() (nodes.ExprNode, error) {
	return p.parseExprPrec(precBitOr)
}

// parseExprPrec parses an expression with the given minimum precedence.
func (p *Parser) parseExprPrec(minPrec int) (nodes.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	start := p.pos()
	// NOT is an expr, not a primary: it stands where an expr may stand and its
	// operand runs down to the next AND, so it is read here where the level the
	// caller asked for is known.
	if p.cur.Type == kwNOT {
		if err := p.offerExprCandidates(); err != nil {
			return nil, err
		}
		if minPrec > precNot {
			return nil, p.syntaxErrorAtCur()
		}
		left, err := p.parseNotExpr(start)
		if err != nil {
			return nil, err
		}
		return p.parseInfixExprPrec(left, start, minPrec, precNot)
	}
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}
	return p.parseInfixExprPrec(left, start, minPrec, precJsonAccess)
}

// parseInfixExprPrec continues Pratt parsing with an already-parsed left
// operand, folding infix operators of at least minPrec. start is the offset
// of the left operand's first token, where every folded expression starts.
// Split from
// parseExprPrec so callers that produce the leading primary themselves — a
// parenthesized query expression that turns out to be the left operand of a
// scalar expression, `((SELECT ...) = 0)` — can resume the operator loop.
//
// level is the highest operator that may take left as its left operand; it
// starts at the top for a primary and drops as operators fold. That is what
// says the grammar's non-associative rules: nothing above a predicate follows
// a predicate, and nothing above AND follows IS TRUE.
func (p *Parser) parseInfixExprPrec(left nodes.ExprNode, start int, minPrec int, level int) (nodes.ExprNode, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	var err error
	for {
		prec, binOp, ok := p.infixPrecedence()
		if !ok || prec < minPrec || prec > level {
			break
		}

		// Handle special infix operators
		switch {
		case p.cur.Type == kwMEMBER:
			left, err = p.parseMemberOfExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			continue

		case p.cur.Type == kwIS:
			left, err = p.parseIsExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			if is, ok := left.(*nodes.IsExpr); ok && is.Test != nodes.IsNull {
				// IS [NOT] TRUE/FALSE/UNKNOWN is an expr, so only AND, XOR and
				// OR may follow it, where IS [NOT] NULL is still a bool_pri.
				level = precNot
			}
			continue

		case p.cur.Type == kwBETWEEN || (p.cur.Type == kwNOT && p.peekNext().Type == kwBETWEEN):
			left, err = p.parseBetweenExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			continue

		case p.cur.Type == kwIN || (p.cur.Type == kwNOT && p.peekNext().Type == kwIN):
			left, err = p.parseInExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			continue

		case p.cur.Type == kwLIKE || (p.cur.Type == kwNOT && p.peekNext().Type == kwLIKE):
			left, err = p.parseLikeExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			continue

		case p.cur.Type == kwSOUNDS:
			left, err = p.parseSoundsLikeExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			continue

		case p.cur.Type == kwREGEXP || p.cur.Type == kwRLIKE ||
			(p.cur.Type == kwNOT && (p.peekNext().Type == kwREGEXP || p.peekNext().Type == kwRLIKE)):
			left, err = p.parseRegexpExpr(left, start)
			if err != nil {
				return nil, err
			}
			level = precComparison
			continue

		case p.cur.Type == kwCOLLATE:
			p.advance() // consume COLLATE
			// The expression COLLATE takes an identifier collation name only:
			// MySQL 8.0 answers 1064 for `x COLLATE BINARY` (the reserved
			// word), unlike the column option's collation_name, which has a
			// BINARY_SYM arm. Measured on mysql:8.0.
			collation, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			left = &nodes.CollateExpr{
				Loc:       nodes.Loc{Start: start, End: p.prev.End},
				Expr:      left,
				Collation: collation,
			}
			continue
		}

		// Regular binary operator
		// Capture original operator text for auto-alias when it differs from canonical form.
		var originalOp string
		switch {
		case p.cur.Type == kwMOD:
			originalOp = "MOD"
		case p.cur.Type == tokNotEq && p.cur.Str == "!=":
			originalOp = "!="
		case p.cur.Type == tokOrOr:
			originalOp = "||"
		case p.cur.Type == tokAndAnd:
			originalOp = "&&"
		}
		isComparison := prec == precComparison
		p.advance() // consume operator

		// Quantified-subquery RHS: `expr comp_op {ANY|SOME|ALL} (subquery)`.
		// Only after a comparison operator, and <=> is not a comp_op: `a <=> ANY
		// (q)` is 1064. The quantifier is recorded on the SubqueryExpr so
		// consumers can distinguish from a scalar subquery comparison.
		if isComparison && binOp != nodes.BinOpNullSafeEq && (p.cur.Type == kwALL || ((p.cur.Type == kwANY || p.cur.Type == kwSOME) && p.peekNext().Type == '(')) {
			// ANY and SOME are not reserved: `1 = any` compares with a column.
			quantTok := p.advance()
			quantifier := strings.ToUpper(quantTok.Str)
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			if !p.isQueryExpressionStart() {
				return nil, p.syntaxErrorAtCur()
			}
			subStart := p.pos()
			sub, err := p.parseSubqueryExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			sub.Loc.Start = subStart
			sub.Loc.End = p.prev.End
			sub.Quantifier = quantifier
			left = &nodes.BinaryExpr{
				Loc:        nodes.Loc{Start: start, End: p.prev.End},
				Op:         binOp,
				Left:       left,
				Right:      sub,
				OriginalOp: originalOp,
			}
			level = precComparison
			continue
		}

		right, err := p.parseExprPrec(prec + 1)
		if err != nil {
			return nil, err
		}

		left = &nodes.BinaryExpr{
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
			Op:         binOp,
			Left:       left,
			Right:      right,
			OriginalOp: originalOp,
		}
		level = prec
	}

	return left, nil
}

// infixPrecedence returns the precedence and operator for the current token
// if it is an infix operator.
func (p *Parser) infixPrecedence() (int, nodes.BinaryOp, bool) {
	switch p.cur.Type {
	case kwOR, tokOrOr:
		return precOr, nodes.BinOpOr, true
	case kwXOR:
		return precXor, nodes.BinOpXor, true
	case kwAND, tokAndAnd:
		return precAnd, nodes.BinOpAnd, true

	// NOT as infix: NOT IN, NOT LIKE, NOT BETWEEN, NOT REGEXP
	case kwNOT:
		switch p.peekNext().Type {
		case kwIN, kwLIKE, kwBETWEEN, kwREGEXP, kwRLIKE:
			return precPredicate, 0, true // handled specially
		}
		return 0, 0, false

	// Comparison operators
	case '=':
		return precComparison, nodes.BinOpEq, true
	case tokNullSafeEq:
		return precComparison, nodes.BinOpNullSafeEq, true
	case tokNotEq:
		return precComparison, nodes.BinOpNe, true
	case '<':
		return precComparison, nodes.BinOpLt, true
	case '>':
		return precComparison, nodes.BinOpGt, true
	case tokLessEq:
		return precComparison, nodes.BinOpLe, true
	case tokGreaterEq:
		return precComparison, nodes.BinOpGe, true
	case kwIS:
		return precComparison, 0, true // handled specially
	case kwLIKE, kwREGEXP, kwRLIKE, kwIN, kwBETWEEN, kwSOUNDS, kwMEMBER:
		return precPredicate, 0, true

	// Bit operators
	case '|':
		return precBitOr, nodes.BinOpBitOr, true
	case '&':
		return precBitAnd, nodes.BinOpBitAnd, true
	case tokShiftLeft:
		return precShift, nodes.BinOpShiftLeft, true
	case tokShiftRight:
		return precShift, nodes.BinOpShiftRight, true

	// Arithmetic
	case '+':
		return precAdd, nodes.BinOpAdd, true
	case '-':
		return precAdd, nodes.BinOpSub, true
	case '*':
		return precMul, nodes.BinOpMul, true
	case '/':
		return precMul, nodes.BinOpDiv, true
	case kwDIV:
		return precMul, nodes.BinOpDivInt, true
	case '%':
		return precMul, nodes.BinOpMod, true
	case kwMOD:
		return precMul, nodes.BinOpMod, true

	case '^':
		return precBitXor, nodes.BinOpBitXor, true

	case kwCOLLATE:
		return precCollate, 0, true

	// JSON column-path operators
	case tokJsonExtract:
		return precJsonAccess, nodes.BinOpJsonExtract, true
	case tokJsonUnquote:
		return precJsonAccess, nodes.BinOpJsonUnquote, true
	}

	return 0, 0, false
}

// offerExprCandidates is the completion hook at the start of any expression: it
// offers column and function candidates when the cursor is here.
func (p *Parser) offerExprCandidates() error {
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return &ParseError{Message: "collecting"}
	}
	return nil
}

// parseNotExpr reads `NOT expr`, whose operand runs down to the next AND, XOR or
// OR. `!` is not this: it is a simple_expr and binds to the next primary.
func (p *Parser) parseNotExpr(start int) (nodes.ExprNode, error) {
	p.advance() // consume NOT
	operand, err := p.parseExprPrec(precNot)
	if err != nil {
		return nil, err
	}
	if operand == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.UnaryExpr{
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
		Op:      nodes.UnaryNot,
		Operand: operand,
	}, nil
}

// parsePrefixExpr parses prefix/primary expressions.
func (p *Parser) parsePrefixExpr() (nodes.ExprNode, error) {
	if err := p.offerExprCandidates(); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case '-':
		return p.parseUnaryExpr(nodes.UnaryMinus)
	case '+':
		// Unary plus — just parse the operand. It recurses here rather than
		// through parseExprPrec, so it counts its own nesting.
		if err := p.enterNesting(); err != nil {
			return nil, err
		}
		defer p.leaveNesting()
		p.advance()
		return p.parsePrefixExpr()
	case '~':
		return p.parseUnaryExpr(nodes.UnaryBitNot)
	case '!':
		expr, err := p.parseUnaryExpr(nodes.UnaryNot)
		if err == nil {
			if ue, ok := expr.(*nodes.UnaryExpr); ok {
				ue.OriginalOp = "!"
			}
		}
		return expr, err
	case kwEXISTS_KW:
		return p.parseExistsExpr()
	case kwCASE:
		return p.parseCaseExpr()
	case kwINTERVAL:
		return p.parseIntervalExpr()
	case '{':
		return p.parseODBCEscape()
	case kwMATCH:
		return p.parseMatchExpr()
	case kwBINARY:
		return p.parseUnaryExpr(nodes.UnaryBinary)
	default:
		return p.parsePrimaryExpr()
	}
}

// atUserVarAssign reports whether the next two tokens start `@name := ...`. A
// user variable is the only thing := assigns to here: `a := 1`, `1 := 2` and
// `(@a) := 1` are syntax errors, and so is `@@a := 1`, since the lexer gives a
// system variable as tokAt2.
func (p *Parser) atUserVarAssign() bool {
	if p.atQuotedUserVar() {
		lookahead := *p.lexer
		return lookahead.NextToken().Type == tokAssign
	}
	return p.cur.Type == tokIDENT && len(p.cur.Str) > 1 && p.cur.Str[0] == '@' &&
		p.peekNext().Type == tokAssign
}

// parseUserVarAssign reads `@name := expr`, with the right side read at rhsPrec.
func (p *Parser) parseUserVarAssign(rhsPrec int) (nodes.ExprNode, error) {
	start := p.pos()
	ref, err := p.parseVariableRef()
	if err != nil {
		return nil, err
	}
	p.advance() // consume :=
	right, err := p.parseExprPrec(rhsPrec)
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Op:    nodes.BinOpAssign,
		Left:  ref,
		Right: right,
	}, nil
}

// parseUnaryExpr parses a unary expression.
func (p *Parser) parseUnaryExpr(op nodes.UnaryOp) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume operator

	operand, err := p.parseExprPrec(precUnary)
	if err != nil {
		return nil, err
	}

	return &nodes.UnaryExpr{
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
		Op:      op,
		Operand: operand,
	}, nil
}

// foldAdjacentStringLits folds MySQL's adjacent string-literal concatenation
// into lit: quoted strings placed next to each other are a single literal whose
// value is the concatenation of the segments ('a' 'b' → 'ab', manual 9.1.1).
// The rule lives on text_literal only — continuation segments must be bare
// quoted strings (never a charset introducer, hex/bit literal, or temporal
// literal; MySQL 8.0.32 and 5.7.25 both reject or alias those forms), and the
// character set of the run is that of the first segment. FirstSegment records
// the first segment's value because the implicit output column name derives
// from it, not from the folded value (SELECT 'a' 'b' → column "a", value "ab").
func (p *Parser) foldAdjacentStringLits(lit *nodes.StringLit) {
	if p.cur.Type != tokSCONST {
		return
	}
	var sb strings.Builder
	sb.WriteString(lit.Value)
	lit.FirstSegment = lit.Value
	lit.Concatenated = true
	for p.cur.Type == tokSCONST {
		next := p.advance()
		sb.WriteString(next.Str)
	}
	lit.Value = sb.String()
	lit.Loc.End = p.prev.End
}

// parsePrimaryExpr parses atoms: literals, column refs, subqueries, func calls, parenthesized exprs.
func (p *Parser) parsePrimaryExpr() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST:
		tok := p.advance()
		// Carry the exact source digits in Text: tok.Ival is clamped to math.MaxInt64 for
		// an unsigned literal beyond int64 range (e.g. BIGINT UNSIGNED DEFAULT
		// 18446744073709551615), so consumers that must round-trip the literal exactly
		// (column defaults) read Text instead of the truncated Value.
		return &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Ival, Text: tok.Str}, nil

	case tokFCONST:
		tok := p.advance()
		if strings.Trim(tok.Str, "0123456789") == "" {
			// A digit run past 64 bits is a decimal to the grammar and an integer to
			// read: it keeps the node it has as a column default, with its digits in Text.
			return &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Ival, Text: tok.Str}, nil
		}
		return &nodes.FloatLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Str}, nil

	case tokSCONST:
		tok := p.advance()
		lit := &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Str}
		p.foldAdjacentStringLits(lit)
		return lit, nil

	case tokXCONST:
		tok := p.advance()
		return &nodes.HexLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Str}, nil

	case tokBCONST:
		tok := p.advance()
		return &nodes.BitLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Str}, nil

	case kwTRUE:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: true}, nil

	case kwFALSE:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: false}, nil

	case kwOFF:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: false}, nil

	case kwNULL:
		tok := p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}}, nil

	case kwDEFAULT:
		return p.parseDefaultExpr()

	case kwCURRENT_DATE, kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwCURRENT_USER, kwLOCALTIME, kwLOCALTIMESTAMP,
		kwUTC_DATE, kwUTC_TIME, kwUTC_TIMESTAMP:
		// The refusals read the keyword before parseTemporalKeyword consumes it.
		if err := p.refusePartitionFunction(strings.ToUpper(p.cur.Str), false, p.cur.Loc); err != nil {
			return nil, err
		}
		return p.parseTemporalKeyword()

	// A `?` parameter marker is a simple_expr alternative in MySQL 8.0
	// (sql_yacc.yy), so it is a valid operand anywhere an expression is: a bare
	// `DO ?`, the offset of LEAD/LAG (read as an ordinary argument), and a
	// window-frame bound (read through parseFrameBound's expr). MySQL accepts
	// `?` only under PREPARE; the catalog analyzer boundary rejects any
	// statement carrying an unbound marker.
	case '?':
		tok := p.advance()
		return &nodes.ParamMarkerExpr{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}}, nil

	case '(':
		return p.parseParenExpr()

	case kwCAST:
		if p.funcWordIsName() {
			return p.parseIdentExpr()
		}
		return p.parseCastExpr()

	case kwEXTRACT:
		if p.funcWordIsName() {
			return p.parseIdentExpr()
		}
		return p.parseExtractExpr()

	case kwCONVERT:
		return p.parseConvertExpr()

	case kwROW:
		return p.parseRowConstructor()

	case kwVALUES:
		return p.parseValuesFunc()

	// Words the parser reads as function calls by token. Unless `(` follows they are names, as
	// in `SELECT count FROM t`, where the server's lexer reads them as identifiers; parseIdentExpr
	// then refuses the words the server reserves.
	case kwRANK, kwDENSE_RANK, kwROW_NUMBER, kwNTILE,
		kwLAG, kwLEAD, kwFIRST_VALUE, kwLAST_VALUE, kwNTH_VALUE,
		kwPERCENT_RANK, kwCUME_DIST, kwGROUPING,
		// Reserved function-name keywords registered in Phase 2
		kwNOW, kwCURDATE, kwCURTIME, kwSYSDATE,
		kwDATE_ADD, kwDATE_SUB, kwMID, kwSUBSTR,
		kwSTD, kwSTDDEV, kwSTDDEV_POP, kwSTDDEV_SAMP,
		kwBIT_AND, kwBIT_OR, kwBIT_XOR,
		kwJSON_ARRAYAGG, kwJSON_OBJECTAGG, kwJSON_DUALITY_OBJECT,
		kwVAR_POP, kwVAR_SAMP, kwVARIANCE,
		// Aggregate/scalar function keywords classified reserved in Phase 3
		kwCOUNT, kwSUM, kwMIN, kwMAX, kwGROUP_CONCAT,
		kwSUBSTRING, kwTRIM, kwPOSITION,
		// IF(cond, val, val) — reserved keyword also used as function
		kwIF,
		// Reserved keywords that are also function names (function_call_keyword / function_call_conflict)
		kwCHAR, kwINSERT, kwLEFT, kwRIGHT, kwDATABASE,
		kwMOD, kwREPEAT, kwREPLACE:
		if p.funcWordIsName() {
			return p.parseIdentExpr()
		}
		start := p.pos()
		tok := p.advance()
		return p.parseFuncCall(start, "", strings.ToUpper(tok.Str))

	default:
		// `@name := expr` is a primary here (variable_aux in sql_yacc.yy), so it
		// stands wherever a primary does and takes all of the expression to its
		// right: `1 + @a := 2 + 3` is 6, with @a set to 5.
		if p.atUserVarAssign() {
			return p.parseUserVarAssign(precAssign)
		}

		// Variable reference
		if p.isVariableRef() {
			return p.parseVariableRef()
		}

		// Temporal literals: DATE '2024-01-01', TIME '12:00:00', TIMESTAMP '2024-01-01 12:00:00'
		if (p.cur.Type == kwDATE || p.cur.Type == kwTIME || p.cur.Type == kwTIMESTAMP) && p.peekNext().Type == tokSCONST {
			typeTok := p.advance() // consume DATE/TIME/TIMESTAMP keyword
			valTok := p.advance()  // consume the string literal
			return &nodes.TemporalLit{
				Loc:   nodes.Loc{Start: typeTok.Loc, End: valTok.End},
				Type:  strings.ToUpper(typeTok.Str),
				Value: valTok.Str,
			}, nil
		}

		// N'text' is a string in the national character set (NCHAR_STRING in
		// sql_yacc.yy's text_literal): the lexer reads N or n directly before
		// a quote as one token, so N 'x' is still the column N with an alias.
		// Adjacent strings fold into it as into any text_literal: N'x' 'y' is
		// 'xy' (mysql:8.0.46).
		if p.atNationalString() {
			nTok := p.advance() // consume N
			strTok := p.advance()
			lit := &nodes.StringLit{
				Loc:      nodes.Loc{Start: nTok.Loc, End: strTok.End},
				Value:    strTok.Str,
				National: true,
			}
			p.foldAdjacentStringLits(lit)
			return lit, nil
		}

		// Charset introducer: _utf8mb4'hello', _latin1'world'. Only a real charset name,
		// unquoted, forms one: `_latin1`'x' is a column with an alias (mysql:8.0.46 answers 1054
		// for it after parsing), and `_latin1` X'61' is a 1064. Only a real
		// charset name forms an introducer — MySQL lexes _typo'abc' as an
		// identifier followed by a separate string literal (oracle 8.0.32 +
		// 5.7.25: SELECT _typo'abc' → 1054 unknown column, not a syntax
		// error), so an unknown name falls through to identifier parsing.
		if p.cur.Type == tokIDENT && isCharsetIntroducer(p.cur.Str) && !p.backtickQuotedAt(p.cur.Loc) && p.peekNext().Type == tokSCONST {
			charsetTok := p.advance() // consume the charset identifier
			strTok := p.advance()     // consume the string literal
			lit := &nodes.StringLit{
				Loc:     nodes.Loc{Start: charsetTok.Loc, End: strTok.End},
				Value:   strTok.Str,
				Charset: charsetTok.Str,
			}
			p.foldAdjacentStringLits(lit)
			return lit, nil
		}

		// An introducer takes a hex or bit literal too, and makes it a string of that charset:
		// mysql:8.0.46 stores _latin1 X'61' as _latin1'a'.
		if p.cur.Type == tokIDENT && isCharsetIntroducer(p.cur.Str) && !p.backtickQuotedAt(p.cur.Loc) {
			if next := p.peekNext().Type; next == tokXCONST || next == tokBCONST {
				charsetTok := p.advance()
				litTok := p.advance()
				// The digits are held to the literal's own rules: X'GG', b'2' and X'6' are 1064.
				value, err := p.decodeByteLiteral(litTok)
				if err != nil {
					return nil, err
				}
				return &nodes.StringLit{
					Loc:     nodes.Loc{Start: charsetTok.Loc, End: litTok.End},
					Value:   value,
					Charset: charsetTok.Str,
				}, nil
			}
		}

		// Identifier — could be column ref or function call
		if p.isIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// funcWordIsName reports whether the function word at p.cur is a name here, which is when no
// `(` follows.
func (p *Parser) funcWordIsName() bool {
	return p.peekNext().Type != '('
}

// parseIdentExpr parses an identifier that could be a column ref or function call.
func (p *Parser) parseIdentExpr() (nodes.ExprNode, error) {
	if p.refusesFunctionName() {
		return nil, p.syntaxErrorAtTok(p.peekNext())
	}
	start := p.pos()
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for function call: name(
	if p.cur.Type == '(' {
		// But first check if this could be a schema-qualified function: name.name(
		// No — MySQL doesn't use schema.func() syntax in expressions (use schema.func())
		return p.parseFuncCall(start, "", name)
	}

	// Check for qualified name: name.name or name.name.name
	if p.cur.Type == '.' {
		p.advance()

		// Completion: after "table.", offer columnref
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}

		name2, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		// name.name( — schema.func()
		if p.cur.Type == '(' {
			return p.parseFuncCall(start, name, name2)
		}

		// Check for third part: schema.table.col or schema.table.*
		if p.cur.Type == '.' {
			p.advance()

			// Completion: after "schema.table.", offer columnref
			p.checkCursor()
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, &ParseError{Message: "collecting"}
			}

			// A three-part name holds its database to check_db_name while it parses:
			// simple_ident_q in MySQL 8.0 sql/sql_yacc.yy.
			if err := refuseDatabaseName(name, start); err != nil {
				return nil, err
			}

			name3, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			return &nodes.ColumnRef{
				Loc:    nodes.Loc{Start: start, End: p.prev.End},
				Schema: name,
				Table:  name2,
				Column: name3,
			}, nil
		}

		// A trigger's NEW and OLD answer to its event as they are read.
		if err := p.refuseTriggerRowRead(name, start); err != nil {
			return nil, err
		}

		// table.col
		return &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Table:  name,
			Column: name2,
		}, nil
	}

	// Plain column ref
	return &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Column: name,
	}, nil
}

// backtickQuotedAt reports whether the token at absolute offset off is a
// backtick-quoted identifier.
func (p *Parser) backtickQuotedAt(off int) bool {
	return p.byteAt(off) == '`'
}

// byteAt returns the input byte at offset off, or 0 outside the input.
func (p *Parser) byteAt(off int) byte {
	if i := off - p.lexer.baseOffset; i >= 0 && i < len(p.lexer.input) {
		return p.lexer.input[i]
	}
	return 0
}

// parseFuncCall parses a function call after the name has been consumed.
func (p *Parser) parseFuncCall(start int, schema, name string) (nodes.ExprNode, error) {
	// The builtins below take their special argument syntax (COUNT(*),
	// DISTINCT, TRIM's FROM, SUBSTRING's FROM and FOR, GROUP_CONCAT's ORDER BY
	// and SEPARATOR, the keyword arguments, OVER) from the grammar's keyword for
	// the name, and only the unqualified, unquoted name is that keyword. A
	// schema-qualified name is a stored function and a backtick-quoted name an
	// identifier: either way the call is an ordinary one, with ordinary
	// expressions for arguments and no window. MySQL 8.0 parses `count`(a) and
	// test.trim(b, 'x') and answers 1630 for the missing stored function, and
	// answers 1064 for `count`(*), `trim`(LEADING 'x' FROM b),
	// test.substring(b FROM 1), `sum`(a) OVER () and the other special forms
	// under a quoted or qualified name. p.prev is the name's token.
	// spacedGeneric records that a gated builtin keyword became a generic call
	// only because a space sits before its '(' (see spacedFuncWords); deparse
	// keeps the space so the call does not re-lex as the builtin.
	spacedGeneric := schema == "" && !p.backtickQuotedAt(p.prev.Loc) && p.spacedFuncWordIsName(strings.ToUpper(name))
	builtin := schema == "" && !p.backtickQuotedAt(p.prev.Loc) && !spacedGeneric
	// A spaced gated keyword is a generic stored-function call, which a partition
	// expression forbids like any qualified or quoted one: mysql:8.0 accepts
	// HASH(sum(a)) but rejects HASH(sum (a)) and HASH(`sum`(a)).
	if err := p.refusePartitionFunction(strings.ToUpper(name), schema != "" || nameQuotedFunction(p, schema) || spacedGeneric, start); err != nil {
		return nil, err
	}
	p.advance() // consume '('

	fc := &nodes.FuncCallExpr{
		Loc:         nodes.Loc{Start: start},
		Schema:      schema,
		Name:        strings.ToUpper(name),
		HasParens:   true,
		SpacedParen: spacedGeneric,
	}

	// Handle special function forms
	upperName := strings.ToUpper(name)

	if builtin && upperName == "JSON_VALUE" {
		return p.parseJSONValueFunc(fc)
	}

	if builtin {
		switch upperName {
		case "COUNT":
			// COUNT(*), COUNT(DISTINCT ...) and the rest of its own
			// productions.
			return p.parseCountFunc(fc)
		case "POSITION":
			// POSITION(substr IN str) — the keyword is the only separator.
			return p.parsePositionFunc(start)
		case "TRIM":
			// TRIM([LEADING|TRAILING|BOTH] [remstr FROM] str)
			return p.parseTrimFunc(fc)
		case "SUBSTRING", "SUBSTR":
			// SUBSTRING(str, pos, len) or SUBSTRING(str FROM pos FOR len)
			return p.parseSubstringFunc(fc)
		// Keyword-argument builtins: their first argument (or a trailing
		// clause) is a bare keyword the engine lexes as a token, never an
		// identifier — a backtick-quoted unit is a syntax error (oracle 8.0.32
		// + 5.7.25).
		case "TIMESTAMPDIFF", "TIMESTAMPADD":
			return p.parseTimestampFunc(fc)
		case "GET_FORMAT":
			return p.parseGetFormatFunc(fc)
		case "WEIGHT_STRING":
			return p.parseWeightStringFunc(fc)
		}

		// An aggregate's arguments may carry a quantifier prefix, and which
		// prefixes and how many arguments is per-aggregate (see funcargs.go).
		// A builtin that is not an aggregate takes neither: mysql:8.0 answers 1064
		// for `abs(distinct c1)`.
		if spec, ok := aggArgSpecs[upperName]; ok {
			// An aggregate whose grammar puts the quantifier before each
			// argument reads it there instead, so ALL appears at most once per
			// argument.
			if !spec.allPerArg {
				q, err := p.parseAggQuantifier(spec.quant)
				if err != nil {
					return nil, err
				}
				fc.Quantifier = q
			}
			if upperName == "GROUP_CONCAT" {
				return p.parseGroupConcatFunc(fc)
			}
			if err := p.parseAggArgs(fc, spec); err != nil {
				return nil, err
			}
			return p.finishFuncCall(fc, builtin)
		}
	}

	// Completion: at the start of the function argument list (including
	// empty parens), offer columnref/func_name candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}

	if builtin {
		if done, err := p.parseTemporalCall(fc); done {
			if err != nil {
				return nil, err
			}
			return fc, nil
		}
	}

	// Empty argument list. CHAR's expr_list is never empty: mysql:8.0.46 answers 1064 for
	// CHAR().
	if p.cur.Type == ')' {
		if builtin && upperName == "CHAR" {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		return p.finishFuncCall(fc, builtin)
	}

	// Regular argument list
	for {
		// A call's arguments are expressions, which neither * nor a bare
		// INTERVAL is: mysql:8.0 answers 1064 for abs(*), `count`(*) and
		// `date_add`(d, INTERVAL 1 DAY), and parses `f`(d + INTERVAL 1 DAY).
		// COUNT is the one builtin whose grammar has a star, in its own
		// production (see parseCountFunc).
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if isBareInterval(arg) {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// CHAR(n, ... USING charset_name). Ref: CHAR_SYM '(' expr_list USING
	// charset_name ')' in MySQL 8.0 sql/sql_yacc.yy, and
	// https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_char.
	if builtin && upperName == "CHAR" && p.cur.Type == kwUSING {
		p.advance()
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("charset")
			return nil, &ParseError{Message: "collecting"}
		}
		cs, err := p.parseCharsetName()
		if err != nil {
			return nil, err
		}
		fc.Using = cs
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		// The production has no window.
		fc.Loc.End = p.prev.End
		return fc, nil
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return p.finishFuncCall(fc, builtin)
}

// finishFuncCall reads the optional OVER clause of a builtin call and closes
// the call's Loc. A quoted or schema-qualified name is a stored function and
// takes no window, which is what builtin says.
func (p *Parser) finishFuncCall(fc *nodes.FuncCallExpr, builtin bool) (nodes.ExprNode, error) {
	// OVER clause for window functions
	//
	// Ref: https://dev.mysql.com/doc/refman/8.0/en/window-functions-usage.html
	//
	//   over_clause:
	//       {OVER (window_spec) | OVER window_name}
	//
	//   window_spec:
	//       [window_name] [partition_clause] [order_clause] [frame_clause]
	//
	//   partition_clause:
	//       PARTITION BY expr [, expr] ...
	//
	//   order_clause:
	//       ORDER BY expr [ASC|DESC] [, expr [ASC|DESC]] ...
	//
	//   frame_clause:
	//       frame_units frame_extent
	//
	//   frame_units:
	//       {ROWS | RANGE | GROUPS}
	//
	//   frame_extent:
	//       {frame_start | frame_between}
	//
	//   frame_between:
	//       BETWEEN frame_start AND frame_end
	//
	//   frame_start, frame_end: {
	//       CURRENT ROW
	//     | UNBOUNDED PRECEDING
	//     | UNBOUNDED FOLLOWING
	//     | expr PRECEDING
	//     | expr FOLLOWING
	//   }
	// A dedicated window function takes a fixed argument count in its
	// WindowFuncCall arm (sql_yacc.yy window_func_call), so a call outside
	// those bounds is a 1064 before OVER is even reached: NTILE(), NTILE(a, b),
	// LEAD() and FIRST_VALUE() are all rejected on mysql:8.0.
	//
	// Its windowing_clause is no option, so the call without OVER is a 1064 too:
	// `SELECT ROW_NUMBER()` and `DO LAG(a) b`.
	if builtin {
		if bounds, ok := windowFuncArity[fc.Name]; ok {
			if len(fc.Args) < bounds[0] || len(fc.Args) > bounds[1] {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type != kwOVER {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}
	// OVER is grammatical only on a dedicated window function or an aggregate
	// (overEligible). On any other builtin, including INTERVAL, or a
	// stored/quoted call (builtin is false), OVER is a 1064: mysql:8.0 rejects
	// `abs(c1) over ()`, `interval(c1) over ()` and `` `sum`(a) over () ``.
	if p.cur.Type == kwOVER {
		if !builtin || !overEligible(fc.Name) {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		wd, err := p.parseOverClause()
		if err != nil {
			return nil, err
		}
		fc.Over = wd
	}

	fc.Loc.End = p.prev.End
	return fc, nil
}

// windowFuncArity is the inclusive [min, max] argument count each dedicated
// window function takes in its window_func_call arm (MySQL 8.0 sql_yacc.yy):
// the rank family takes none, NTILE/FIRST_VALUE/LAST_VALUE one, NTH_VALUE two,
// and LEAD/LAG one to three (expr opt_lead_lag_info, where opt_lead_lag_info
// adds `, NUM_literal` and then `, expr`). mysql:8.0 answers 1064 for a call
// outside these bounds, e.g. NTILE(), NTILE(a, b), LEAD() and FIRST_VALUE().
var windowFuncArity = map[string][2]int{
	"ROW_NUMBER":   {0, 0},
	"RANK":         {0, 0},
	"DENSE_RANK":   {0, 0},
	"CUME_DIST":    {0, 0},
	"PERCENT_RANK": {0, 0},
	"NTILE":        {1, 1},
	"LEAD":         {1, 3},
	"LAG":          {1, 3},
	"FIRST_VALUE":  {1, 1},
	"LAST_VALUE":   {1, 1},
	"NTH_VALUE":    {2, 2},
}

// overEligible reports whether an OVER clause may follow a builtin call. Only
// the dedicated window functions (window_func_call) and the aggregate
// functions (sum_expr: each ends in opt_windowing_clause) carry a windowing
// clause. Any other builtin with OVER — INTERVAL, GROUPING, a generic builtin —
// is a 1064 on mysql:8.0, e.g. `abs(c1) over ()` and `interval(c1) over ()`.
func overEligible(name string) bool {
	if _, ok := windowFuncArity[name]; ok {
		return true
	}
	if _, ok := aggArgSpecs[name]; ok {
		return true
	}
	return name == "COUNT"
}

// parseTrimFunc parses TRIM([LEADING|TRAILING|BOTH] [remstr FROM] str).
// parseMemberOfExpr parses value MEMBER OF(json_array).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html
//
//	value MEMBER OF(json_array)
func (p *Parser) parseMemberOfExpr(value nodes.ExprNode, start int) (nodes.ExprNode, error) {
	p.advance() // consume MEMBER

	// OF is optional: `1 MEMBER ('[1]')` parses (MEMBER_SYM opt_of).
	p.match(kwOF)

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	array, err := p.parseExprPrec(precCollate) // a simple_expr: `MEMBER OF (1+1)` is 1064
	if err != nil {
		return nil, err
	}
	if array == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.MemberOfExpr{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Value: value,
		Array: array,
	}, nil
}

// parseOverClause parses OVER (window_spec) or OVER window_name.
func (p *Parser) parseOverClause() (*nodes.WindowDef, error) {
	start := p.pos()

	// OVER window_name (identifier reference)
	if p.cur.Type != '(' {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return &nodes.WindowDef{
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
			RefName: name,
		}, nil
	}

	// OVER (window_spec)
	p.advance() // consume '('

	// Completion: inside OVER (, offer PARTITION BY, ORDER BY, and frame keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwPARTITION)
		p.addTokenCandidate(kwORDER)
		p.addTokenCandidate(kwROWS)
		p.addTokenCandidate(kwRANGE)
		p.addTokenCandidate(kwGROUPS)
		return nil, &ParseError{Message: "collecting"}
	}

	wd := &nodes.WindowDef{Loc: nodes.Loc{Start: start}}

	// Optional window_name reference
	if p.isIdentToken() && p.cur.Type != kwPARTITION && p.cur.Type != kwORDER &&
		p.cur.Type != kwROWS && p.cur.Type != kwRANGE && p.cur.Type != kwGROUPS {
		var err error
		wd.RefName, _, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
	}

	// PARTITION BY
	if p.cur.Type == kwPARTITION {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		wd.PartitionBy = exprs
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		wd.OrderBy = orderBy
	}

	// frame_clause: {ROWS | RANGE | GROUPS} frame_extent
	if p.cur.Type == kwROWS || p.cur.Type == kwRANGE || p.cur.Type == kwGROUPS {
		frame, err := p.parseFrameClause()
		if err != nil {
			return nil, err
		}
		wd.Frame = frame
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	wd.Loc.End = p.prev.End
	return wd, nil
}

// parseFrameClause parses a window frame clause.
func (p *Parser) parseFrameClause() (*nodes.WindowFrame, error) {
	start := p.pos()
	frame := &nodes.WindowFrame{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwROWS:
		frame.Type = nodes.FrameRows
	case kwRANGE:
		frame.Type = nodes.FrameRange
	case kwGROUPS:
		frame.Type = nodes.FrameGroups
	}
	p.advance()

	// Completion: after ROWS/RANGE/GROUPS, offer frame extent keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwBETWEEN)
		p.addTokenCandidate(kwUNBOUNDED)
		p.addTokenCandidate(kwCURRENT)
		return nil, &ParseError{Message: "collecting"}
	}

	// frame_extent: frame_start | BETWEEN frame_start AND frame_end
	if _, ok := p.match(kwBETWEEN); ok {
		startBound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.Start = startBound
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		endBound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.End = endBound
	} else {
		// Single frame_start
		startBound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.Start = startBound
	}

	frame.Loc.End = p.prev.End
	return frame, nil
}

// parseFrameBound parses a window frame bound.
func (p *Parser) parseFrameBound() (*nodes.WindowFrameBound, error) {
	start := p.pos()
	bound := &nodes.WindowFrameBound{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwCURRENT {
		p.advance()
		p.match(kwROW)
		bound.Type = nodes.BoundCurrentRow
	} else if p.cur.Type == kwUNBOUNDED {
		p.advance()
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundUnboundedPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundUnboundedFollowing
		} else {
			return nil, &ParseError{
				Message:  "expected PRECEDING or FOLLOWING after UNBOUNDED",
				Position: p.cur.Loc,
			}
		}
	} else {
		// expr PRECEDING | expr FOLLOWING
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		bound.Offset = expr
		if _, ok := p.match(kwPRECEDING); ok {
			bound.Type = nodes.BoundPreceding
		} else if _, ok := p.match(kwFOLLOWING); ok {
			bound.Type = nodes.BoundFollowing
		} else {
			return nil, &ParseError{
				Message:  "expected PRECEDING or FOLLOWING",
				Position: p.cur.Loc,
			}
		}
	}

	bound.Loc.End = p.prev.End
	return bound, nil
}

// parseTrimFunc parses TRIM's argument list, after the '('. A direction is
// carried in the name (TRIM_LEADING, TRIM_TRAILING, TRIM_BOTH) and Args are
// [remstr, str], or [str] when no remstr is written, which trims spaces.
//
// Ref: trim_operands in MySQL 8.0 sql/sql_yacc.yy, and
// https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_trim:
//
//	TRIM([{BOTH | LEADING | TRAILING} [remstr] FROM] str), TRIM([remstr FROM] str)
//
// A direction always takes FROM: mysql:8.0.46 answers 1064 for TRIM(LEADING 'x').
func (p *Parser) parseTrimFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	direction := false
	switch p.cur.Type {
	case kwLEADING, kwTRAILING, kwBOTH:
		fc.Name = "TRIM_" + strings.ToUpper(p.cur.Str)
		p.advance()
		direction = true
		if p.cur.Type == kwFROM {
			// TRIM(LEADING FROM str)
			p.advance()
			str, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if str == nil {
				return nil, p.syntaxErrorAtCur()
			}
			fc.Args = []nodes.ExprNode{str}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			fc.Loc.End = p.prev.End
			return fc, nil
		}
	}

	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type == kwFROM {
		p.advance()
		// remstr FROM str
		str, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if str == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if fc.Name == "TRIM" {
			// No direction keyword: TRIM(remstr FROM str). The engine keeps
			// this form as-is in stored bodies — trim(' x' from 'axa'), not
			// a comma-argument call (oracle 8.0.32 + 5.7.25) — so mark it
			// for the deparser the same way the directional forms are.
			fc.Name = "TRIM_FROM"
		}
		fc.Args = []nodes.ExprNode{arg, str}
	} else if direction {
		return nil, p.syntaxErrorAtCur()
	} else {
		fc.Args = []nodes.ExprNode{arg}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseSubstringFunc parses SUBSTRING(str, pos, len) or SUBSTRING(str FROM pos [FOR len]).
func (p *Parser) parseSubstringFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}
	fc.Args = append(fc.Args, arg)

	if p.cur.Type == kwFROM {
		p.advance()
		pos, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if pos == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, pos)

		if p.cur.Type == kwFOR {
			p.advance()
			length, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if length == nil {
				return nil, p.syntaxErrorAtCur()
			}
			fc.Args = append(fc.Args, length)
		}
	} else if p.cur.Type == ',' {
		p.advance()
		pos, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if pos == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, pos)

		if p.cur.Type == ',' {
			p.advance()
			length, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if length == nil {
				return nil, p.syntaxErrorAtCur()
			}
			fc.Args = append(fc.Args, length)
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// sqlTSIUnits maps the SQL_TSI_-prefixed unit synonyms MySQL's lexer
// registers to their canonical unit. Only these eight exist —
// SQL_TSI_MICROSECOND is not a token and the engine rejects it
// (oracle 8.0.32 + 5.7.25: error 1064).
var sqlTSIUnits = map[string]string{
	"SQL_TSI_SECOND":  "SECOND",
	"SQL_TSI_MINUTE":  "MINUTE",
	"SQL_TSI_HOUR":    "HOUR",
	"SQL_TSI_DAY":     "DAY",
	"SQL_TSI_WEEK":    "WEEK",
	"SQL_TSI_MONTH":   "MONTH",
	"SQL_TSI_QUARTER": "QUARTER",
	"SQL_TSI_YEAR":    "YEAR",
}

// normalizeTimeUnit uppercases unit and folds a SQL_TSI_ synonym to its
// canonical name. This mirrors the engine's stored form: SHOW CREATE VIEW
// renders timestampdiff(sql_tsi_second,...) as timestampdiff(SECOND,...)
// (oracle 8.0.32 + 5.7.25).
func normalizeTimeUnit(unit string) string {
	u := strings.ToUpper(unit)
	if canon, ok := sqlTSIUnits[u]; ok {
		return canon
	}
	return u
}

// simpleTimeUnits is the unit set TIMESTAMPDIFF/TIMESTAMPADD accept
// (interval_time_stamp in sql_yacc.yy): the nine simple units only.
// Compound units are a syntax error there (oracle 8.0.32 + 5.7.25:
// TIMESTAMPDIFF(DAY_HOUR,...) → 1064).
var simpleTimeUnits = map[string]bool{
	"MICROSECOND": true,
	"SECOND":      true,
	"MINUTE":      true,
	"HOUR":        true,
	"DAY":         true,
	"WEEK":        true,
	"MONTH":       true,
	"QUARTER":     true,
	"YEAR":        true,
}

// parseTimeUnitKeyword parses a bare temporal-unit keyword and returns it as
// a normalized KeywordArg. The engine lexes the unit as a keyword token —
// quoted identifiers and string literals are syntax errors in this position
// (oracle 8.0.32 + 5.7.25) — so only keyword tokens are accepted; all valid
// units are registered keywords. simpleOnly restricts to the
// TIMESTAMPDIFF/TIMESTAMPADD set; otherwise the full interval-unit set
// (including compound units) is allowed.
func (p *Parser) parseTimeUnitKeyword(simpleOnly bool) (*nodes.KeywordArg, error) {
	if p.cur.Type < 700 {
		return nil, &ParseError{Message: "expected time unit keyword", Position: p.cur.Loc}
	}
	tok := p.advance()
	unit := normalizeTimeUnit(tok.Str)
	valid := isValidIntervalUnit(unit)
	if simpleOnly {
		valid = simpleTimeUnits[unit]
	}
	if !valid {
		return nil, &ParseError{Message: "invalid time unit: " + tok.Str, Position: tok.Loc}
	}
	return &nodes.KeywordArg{
		Loc:     nodes.Loc{Start: tok.Loc, End: p.prev.End},
		Keyword: unit,
	}, nil
}

// parseTimestampFunc parses TIMESTAMPDIFF(unit, expr, expr) and
// TIMESTAMPADD(unit, expr, expr) — the unit is a bare keyword argument.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_timestampdiff
func (p *Parser) parseTimestampFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Completion parity with the generic argument path.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}
	unit, err := p.parseTimeUnitKeyword(true)
	if err != nil {
		return nil, err
	}
	fc.Args = append(fc.Args, unit)
	for range 2 {
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, arg)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseGetFormatFunc parses GET_FORMAT({DATE|TIME|DATETIME|TIMESTAMP}, expr).
// The first argument is a bare keyword; TIMESTAMP is a synonym the engine
// stores as DATETIME, and anything else — including YEAR and string
// literals — is a syntax error (oracle 8.0.32 + 5.7.25).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_get-format
func (p *Parser) parseGetFormatFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Completion parity with the generic argument path.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}
	if p.cur.Type < 700 {
		return nil, &ParseError{Message: "expected DATE, TIME, DATETIME, or TIMESTAMP", Position: p.cur.Loc}
	}
	tok := p.advance()
	kind := strings.ToUpper(tok.Str)
	switch kind {
	case "DATE", "TIME", "DATETIME":
	case "TIMESTAMP":
		kind = "DATETIME"
	default:
		return nil, &ParseError{Message: "invalid GET_FORMAT type: " + tok.Str, Position: tok.Loc}
	}
	fc.Args = append(fc.Args, &nodes.KeywordArg{
		Loc:     nodes.Loc{Start: tok.Loc, End: p.prev.End},
		Keyword: kind,
	})
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}
	fc.Args = append(fc.Args, arg)
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseWeightStringFunc parses
//
//	WEIGHT_STRING(expr [AS {CHAR|BINARY}(n)])
//
// The plain no-suffix form stays a generic FuncCallExpr (its round-trip
// predates this parser and is pinned by tests); the AS forms build a
// WeightStringExpr. AS BINARY(n) is not stored by the engine — it desugars
// to weight_string(cast(expr as char(n) charset binary)) (oracle 8.0.32),
// so the parser produces that shape directly.
//
// MySQL 5.7 also took a LEVEL clause here (opt_ws_levels). 8.0 removed it
// from the grammar, so LEVEL is a 1064 at the LEVEL token, which is where
// the closing parenthesis is expected below (measured on mysql:8.0).
//
// Ref: the WEIGHT_STRING_SYM arms of function_call_conflict in mysql-server
// 8.0 sql/sql_yacc.yy;
// https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_weight-string
func (p *Parser) parseWeightStringFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Completion parity with the generic argument path.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}
	argStart := p.cur.Loc
	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}

	ws := &nodes.WeightStringExpr{Loc: fc.Loc, Expr: arg}
	asBinary := false

	if p.cur.Type == kwAS {
		p.advance()
		var typeName string
		switch p.cur.Type {
		case kwCHAR, kwCHARACTER:
			// CHARACTER is a spelling of CHAR_SYM (sql_yacc.yy), so
			// WEIGHT_STRING('x' AS CHARACTER(1)) parses on MySQL 8.0 as
			// WEIGHT_STRING('x' AS CHAR(1)) does.
			typeName = "CHAR"
		case kwBINARY:
			typeName = "BINARY"
		default:
			return nil, &ParseError{Message: "expected CHAR or BINARY", Position: p.cur.Loc}
		}
		typeTok := p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if p.cur.Type != tokICONST {
			return nil, &ParseError{Message: "expected length", Position: p.cur.Loc}
		}
		if p.cur.Ival < 1 {
			// The engine requires N >= 1: WEIGHT_STRING('x' AS CHAR(0)) is a
			// 1064 syntax error (oracle 8.0.32 + 5.7.25).
			return nil, &ParseError{Message: "WEIGHT_STRING AS length must be at least 1", Position: p.cur.Loc}
		}
		lenTok := p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		dt := &nodes.DataType{
			Loc:    nodes.Loc{Start: typeTok.Loc, End: p.prev.End},
			Name:   "CHAR",
			Length: int(lenTok.Ival),
		}
		if typeName == "BINARY" {
			// Engine-stored form of AS BINARY(n):
			// weight_string(cast(expr as char(n) charset binary)).
			dt.Charset = "binary"
			asBinary = true
			// The synthetic cast wraps arg, so its Loc must start at arg, not
			// at the BINARY token, to stay inside WeightStringExpr and keep arg
			// inside it (the Loc contract: a child lies within its parent).
			ws.Expr = &nodes.CastExpr{
				Loc:      nodes.Loc{Start: argStart, End: p.prev.End},
				Expr:     arg,
				TypeName: dt,
			}
		} else {
			ws.AsChar = dt
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if ws.AsChar == nil && !asBinary {
		// Plain WEIGHT_STRING(expr) — keep the generic function-call shape.
		fc.Args = append(fc.Args, arg)
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	ws.Loc.End = p.prev.End
	return ws, nil
}

// parseGroupConcatFunc parses GROUP_CONCAT([DISTINCT] expr [, expr ...] [ORDER BY ...] [SEPARATOR str]).
func (p *Parser) parseGroupConcatFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Parse arguments. The grammar takes an expr_list, so group_concat() is a
	// syntax error (mysql:8.0: 1064).
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, arg)

		if p.cur.Type == ',' {
			// Peek to see if next is ORDER or SEPARATOR — if so, don't consume comma
			next := p.peekNext()
			if next.Type == kwORDER || next.Type == kwSEPARATOR {
				break
			}
			p.advance()
		} else {
			break
		}
	}

	// ORDER BY
	if p.cur.Type == kwORDER {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		fc.OrderBy = orderBy
	}

	// SEPARATOR
	if p.cur.Type == kwSEPARATOR {
		p.advance()
		sep, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if sep == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Separator = sep
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	// The production ends in an optional windowing clause, so the call takes a
	// window: mysql:8.0 parses `group_concat(c1) over ()` and answers 1235, the
	// unsupported-feature error, after it.
	return p.finishFuncCall(fc, true)
}

// subqueryHeadOutcome classifies what a '(' run that looked like a query
// expression (isQueryExpressionStart) actually turned out to be.
type subqueryHeadOutcome int

const (
	// subqueryHeadClosed: the query expression parsed and the next token is
	// ')' — plain subquery parens. The caller consumes the ')' and owns the
	// Loc fixups.
	subqueryHeadClosed subqueryHeadOutcome = iota
	// subqueryHeadOperand: the query expression parsed but more tokens follow
	// inside the parens — the subquery heads a scalar expression
	// (`((SELECT ...) = 0)`) or list. The caller resumes the infix loop.
	subqueryHeadOperand
	// subqueryHeadScalar: the query parse failed on a head nested behind
	// extra parens — a scalar paren layer wraps a subquery-headed expression
	// (`(((SELECT 1) = 0))`). The scan has been rewound; the caller reparses
	// as a scalar expression. The discarded speculative error accompanies
	// this outcome (see parseSubqueryHead) for furthest-error reporting.
	subqueryHeadScalar
)

// furthestErr picks the error that made it furthest into the input — the
// standard heuristic when two alternative parses of the same region both
// fail. Used to keep error positions truthful when a rewound speculative
// query parse failed DEEPER than the scalar reparse does (e.g.
// `((SELECT 1) UNION (SELECT 2 FROM))`: the speculation reports the real
// mistake after FROM, while the scalar reparse only gets as far as UNION).
func furthestErr(reparse, speculative error) error {
	rp, ok1 := reparse.(*ParseError)
	sp, ok2 := speculative.(*ParseError)
	if ok1 && ok2 && sp.Position > rp.Position {
		return speculative
	}
	return reparse
}

// parseSubqueryHead speculatively parses the query expression behind a '('
// run that isQueryExpressionStart matched, and classifies the result. Only
// tokens PAST the query expression distinguish `((SELECT 1))` from
// `((SELECT 1) = 0)` from `(((SELECT 1) = 0))`, so the parse must be
// speculative when the query-primary head is nested behind more parens:
// on failure there, the scan (and the completion collect latch — see
// scanState) is rewound and the caller re-reads the run as a scalar
// expression, peeling one paren per recursion level. When the head token
// itself is the query primary there is nothing to unwrap, so a parse error
// is a real subquery syntax error and is returned at its true position.
func (p *Parser) parseSubqueryHead() (*nodes.SubqueryExpr, subqueryHeadOutcome, error) {
	nested := p.cur.Type == '('
	var save scanState
	if nested {
		save = p.saveScan() // only the nested-head failure path rewinds
	}
	sub, err := p.parseSubqueryExpr()
	switch {
	case err == nil && p.cur.Type == ')':
		return sub, subqueryHeadClosed, nil
	case err == nil && nested && p.prev.Type == ')':
		// The query expression ended at a closing paren and more tokens
		// follow: the parenthesized subquery is a scalar operand.
		sub.Loc.End = p.prev.End
		return sub, subqueryHeadOperand, nil
	case err == nil:
		// Trailing tokens after a query expression that did NOT end at its
		// own ')': a BARE head — `(TABLE t = 1)`, `(VALUES ROW(1) = 1)` — or
		// unparenthesized trailing clauses — `((SELECT 1) LIMIT 1 = 1)`,
		// `((SELECT 1) ORDER BY 1 DESC = 1)`. MySQL rejects these forms
		// (1064 on 8.0.32 and 5.7.25; the legal spellings are
		// `((SELECT 1)) = 1` and `(((SELECT 1) LIMIT 1) = 1)`). Keep the
		// pre-fix hard ')' expectation and error position.
		//
		// The operand gate above (prev == ')') is a deliberate harmless
		// SUPERSET: it also matches the closing paren of a right set-op ARM,
		// so `((SELECT 1) UNION (SELECT 2) = 1)` parses here even though the
		// engine 1064s it (both versions). The accepted parse deparses to the
		// legal extra-wrapped spelling, and readbacks never contain the raw
		// form, so nothing invalid round-trips.
		_, errP := p.expect(')')
		return nil, 0, errP
	case isNestingError(err):
		// A nesting error ends the parse. Retrying the scalar reading would
		// repeat both readings at every enclosing '(', which takes time
		// exponential in the depth, since parseSelectWithParens does not
		// record a nesting error.
		return nil, 0, err
	case nested:
		// Rewind and let the caller reparse as a scalar expression. The
		// discarded error rides along with the scalar outcome (it is NOT a
		// failure): if the reparse also fails, the caller reports whichever
		// error lies further into the input.
		p.restoreScan(save)
		return nil, subqueryHeadScalar, err
	default:
		return nil, 0, err
	}
}

// parseParenExpr parses a parenthesized expression or subquery.
func (p *Parser) parseParenExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume '('

	var expr nodes.ExprNode
	var specErr error // discarded speculative error (subqueryHeadScalar only)
	// Check for subquery. isQueryExpressionStart scans past a leading run of
	// '(' to the innermost head, so `((SELECT ...` lands here whether the
	// parens wrap a query expression — `((SELECT 1))`, `((SELECT 1) UNION
	// (SELECT 2))` — or a scalar expression whose first primary is a
	// parenthesized subquery — `((SELECT 1) = 0)`, `(((SELECT 1) = 0))`.
	if p.isQueryExpressionStart() {
		subStart := p.pos()
		sub, outcome, err := p.parseSubqueryHead()
		if outcome == subqueryHeadScalar {
			specErr = err
		} else if err != nil {
			return nil, err
		}
		switch outcome {
		case subqueryHeadClosed:
			// Plain parenthesized query expression `(SELECT ...)`.
			p.advance()
			sub.Loc.Start = start
			sub.Loc.End = p.prev.End
			return sub, nil
		case subqueryHeadOperand:
			// The query expression is the LEFT operand of a scalar expression
			// inside the parens — `((SELECT ...) = 0)`, `(((SELECT ...)) + 1)`
			// — or the first item of a row constructor `((SELECT ...), 2)`.
			// MySQL 8.0.32 and 5.7.25 both accept every binary operator here
			// (oracle-probed; the stock sys.metrics view stores this shape).
			// Resume the infix loop with the subquery as the parsed primary
			// and fall through to the shared row-constructor / ')' handling.
			expr, err = p.parseInfixExprPrec(sub, subStart, precAssign, precJsonAccess)
			if err != nil {
				return nil, err
			}
		default: // subqueryHeadScalar
			// The paren run was not a query expression after all: a scalar
			// paren layer wraps a subquery-headed expression, e.g. the
			// `(((SELECT 1) = 0))` shape. The scan was rewound; expr stays
			// nil and the shared scalar parse below re-reads the content —
			// the recursion (parsePrimaryExpr → this function) peels exactly
			// one paren per level and re-decides, so it terminates, and the
			// innermost level reports real subquery syntax errors at their
			// true position.
		}
	}
	if expr == nil {
		var err error
		expr, err = p.parseExpr()
		if err != nil {
			return nil, furthestErr(err, specErr)
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type == ',' {
		if err := checkRowItem(expr, start+1); err != nil {
			return nil, err
		}
		row := &nodes.RowExpr{Loc: nodes.Loc{Start: start}, Items: []nodes.ExprNode{expr}}
		for p.cur.Type == ',' {
			p.advance()
			item, err := p.parseRowItem()
			if err != nil {
				return nil, furthestErr(err, specErr)
			}
			row.Items = append(row.Items, item)
		}
		if _, err := p.expect(')'); err != nil {
			return nil, furthestErr(err, specErr)
		}
		row.Loc.End = p.prev.End
		return row, nil
	}

	if _, err := p.expect(')'); err != nil {
		return nil, furthestErr(err, specErr)
	}

	return &nodes.ParenExpr{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Expr: expr,
	}, nil
}

// parseCastExpr parses CAST(expr AS type).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html
//
//	CAST(expr AS data_type)
func (p *Parser) parseCastExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CAST

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	dt, err := p.parseCastDataType()
	if err != nil {
		return nil, err
	}

	// Optional trailing ARRAY for multi-valued functional indexes (8.0.17+):
	// CAST(expr AS type ARRAY). Maps to the grammar's opt_array_cast. ARRAY is a
	// non-reserved keyword, consumed only in this cast-target position.
	isArray := false
	if _, ok := p.match(kwARRAY); ok {
		isArray = true
		if err := p.castArrayRefused(dt); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.CastExpr{
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
		Expr:     expr,
		TypeName: dt,
		Array:    isArray,
	}, nil
}

// parseExtractExpr parses EXTRACT(unit FROM expr).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_extract
//
//	EXTRACT(unit FROM expr)
func (p *Parser) parseExtractExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume EXTRACT

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Parse the unit keyword. EXTRACT accepts the full interval-unit set,
	// including compound units (DAY_HOUR, YEAR_MONTH, ...), which are
	// reserved keywords — so this must accept keyword tokens, not
	// identifiers (oracle 8.0.32 + 5.7.25: extract(day_hour from ...) is
	// the engine-stored form).
	unit, err := p.parseTimeUnitKeyword(false)
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
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.ExtractExpr{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Unit: unit.Keyword,
		Expr: expr,
	}, nil
}

// parseConvertExpr parses CONVERT(expr, type) or CONVERT(expr USING charset).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/cast-functions.html
//
//	CONVERT(expr, type)
//	CONVERT(expr USING transcoding_name)
func (p *Parser) parseConvertExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CONVERT

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	conv := &nodes.ConvertExpr{
		Loc:  nodes.Loc{Start: start},
		Expr: expr,
	}

	if p.cur.Type == kwUSING {
		p.advance()
		// Completion: after CONVERT(... USING, offer charset candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("charset")
			return nil, &ParseError{Message: "collecting"}
		}
		// charset_name is ident_or_text | BINARY_SYM (sql_yacc.yy): an
		// identifier, a string literal, or the reserved word BINARY, so
		// CONVERT(x USING BINARY) parses on MySQL 8.0.
		charset, err := p.parseCharsetName()
		if err != nil {
			return nil, err
		}
		conv.Charset = charset
	} else if p.cur.Type == ',' {
		p.advance()
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		conv.TypeName = dt
	} else {
		// CONVERT(expr) has no one-argument form: every engine answers 1064 for it.
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	conv.Loc.End = p.prev.End
	return conv, nil
}

// parseIsExpr parses IS [NOT] NULL / TRUE / FALSE / UNKNOWN.
func (p *Parser) parseIsExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	p.advance() // consume IS

	// Completion: after IS, offer NULL/NOT/TRUE/FALSE/UNKNOWN.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwNULL)
		p.addTokenCandidate(kwNOT)
		p.addTokenCandidate(kwTRUE)
		p.addTokenCandidate(kwFALSE)
		return nil, &ParseError{Message: "collecting"}
	}

	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}

	is := &nodes.IsExpr{
		Loc:  nodes.Loc{Start: start},
		Not:  not,
		Expr: left,
	}

	switch p.cur.Type {
	case kwNULL:
		is.Test = nodes.IsNull
		p.advance()
	case kwTRUE:
		is.Test = nodes.IsTrue
		p.advance()
	case kwFALSE:
		is.Test = nodes.IsFalse
		p.advance()
	default:
		// IS [NOT] UNKNOWN — check for identifier "unknown"
		if p.cur.Type == kwUNKNOWN {
			is.Test = nodes.IsUnknown
			p.advance()
		} else if p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()
		} else {
			return nil, &ParseError{
				Message:  "expected NULL, TRUE, FALSE, or UNKNOWN after IS",
				Position: p.cur.Loc,
			}
		}
	}

	is.Loc.End = p.prev.End
	return is, nil
}

// parseBetweenExpr parses [NOT] BETWEEN low AND high.
func (p *Parser) parseBetweenExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume BETWEEN

	// Each side is read above AND, which is what keeps this AND the separator.
	low, err := p.parseExprPrec(precBitOr)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAND); err != nil {
		return nil, err
	}

	high, err := p.parseExprPrec(precPredicate)
	if err != nil {
		return nil, err
	}

	return &nodes.BetweenExpr{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Not:  not,
		Expr: left,
		Low:  low,
		High: high,
	}, nil
}

// parseInExpr parses [NOT] IN (value_list) or [NOT] IN (subquery).
func (p *Parser) parseInExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume IN

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	inExpr := &nodes.InExpr{
		Loc:  nodes.Loc{Start: start},
		Not:  not,
		Expr: left,
	}

	// Check for subquery. Speculative for the same reason as parseParenExpr:
	// `IN ((SELECT ...))` is a table subquery, but the same leading tokens may
	// open a VALUE list instead — `IN ((SELECT 1) = 1)`, `IN ((SELECT 1), 2)`,
	// `IN (((SELECT 1) = 1))` (all accepted by MySQL 8.0.32 and 5.7.25).
	var specErr error // discarded speculative error (subqueryHeadScalar only)
	if p.isQueryExpressionStart() {
		subStart := p.pos()
		sub, outcome, err := p.parseSubqueryHead()
		if outcome == subqueryHeadScalar {
			specErr = err
		} else if err != nil {
			return nil, err
		}
		switch outcome {
		case subqueryHeadClosed:
			// Table subquery: `x IN (SELECT ...)` / `x IN ((SELECT ...))`.
			p.advance()
			inExpr.Select = sub.Select
			inExpr.Loc.End = p.prev.End
			return inExpr, nil
		case subqueryHeadOperand:
			// The query expression starts the first VALUE-LIST element:
			// resume the infix loop for that element and seed the shared
			// value-list loop below with it.
			first, err := p.parseInfixExprPrec(sub, subStart, precAssign, precJsonAccess)
			if err != nil {
				return nil, err
			}
			inExpr.List = append(inExpr.List, first)
		default: // subqueryHeadScalar
			// A scalar paren layer wraps the subquery-headed first element,
			// e.g. `IN (((SELECT 1) = 1))`; the scan was rewound — parse as
			// an ordinary value list below.
		}
	}

	// Value list; may already be seeded with a subquery-headed first element.
	for len(inExpr.List) == 0 || p.cur.Type == ',' {
		if len(inExpr.List) > 0 {
			p.advance() // consume ','
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, furthestErr(err, specErr)
		}
		if val == nil {
			return nil, p.syntaxErrorAtCur()
		}
		inExpr.List = append(inExpr.List, val)
	}

	if _, err := p.expect(')'); err != nil {
		return nil, furthestErr(err, specErr)
	}

	inExpr.Loc.End = p.prev.End
	return inExpr, nil
}

// parseLikeExpr parses [NOT] LIKE pattern [ESCAPE escape_char].
func (p *Parser) parseLikeExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume LIKE

	pattern, err := p.parseExprPrec(precCollate)
	if err != nil {
		return nil, err
	}

	like := &nodes.LikeExpr{
		Loc:     nodes.Loc{Start: start},
		Not:     not,
		Expr:    left,
		Pattern: pattern,
	}

	if _, ok := p.match(kwESCAPE); ok {
		esc, err := p.parseExprPrec(precCollate)
		if err != nil {
			return nil, err
		}
		like.Escape = esc
	}

	like.Loc.End = p.prev.End
	return like, nil
}

// parseRegexpExpr parses [NOT] REGEXP pattern.
func (p *Parser) parseRegexpExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume REGEXP or RLIKE

	pattern, err := p.parseExprPrec(precBitOr)
	if err != nil {
		return nil, err
	}

	// Represent as BinaryExpr with BinOpRegexp
	expr := &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Op:    nodes.BinOpRegexp,
		Left:  left,
		Right: pattern,
	}

	if not {
		return &nodes.UnaryExpr{
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
			Op:      nodes.UnaryNot,
			Operand: expr,
		}, nil
	}

	return expr, nil
}

// parseCaseExpr parses a CASE expression.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/flow-control-functions.html#operator_case
//
//	CASE [operand]
//	    WHEN condition THEN result
//	    [WHEN condition THEN result ...]
//	    [ELSE result]
//	END
func (p *Parser) parseCaseExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume CASE

	ce := &nodes.CaseExpr{Loc: nodes.Loc{Start: start}}

	// Simple CASE: CASE operand WHEN ...
	if p.cur.Type != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if operand == nil {
			return nil, p.syntaxErrorAtCur()
		}
		ce.Operand = operand
	}

	// WHEN clauses
	for p.cur.Type == kwWHEN {
		whenStart := p.pos()
		p.advance() // consume WHEN
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if cond == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		result, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if result == nil {
			return nil, p.syntaxErrorAtCur()
		}
		ce.Whens = append(ce.Whens, &nodes.CaseWhen{
			Loc:    nodes.Loc{Start: whenStart, End: p.prev.End},
			Cond:   cond,
			Result: result,
		})
	}

	// ELSE
	if _, ok := p.match(kwELSE); ok {
		def, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if def == nil {
			return nil, p.syntaxErrorAtCur()
		}
		ce.Default = def
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}

	ce.Loc.End = p.prev.End
	return ce, nil
}

// parseExistsExpr parses EXISTS (subquery).
func (p *Parser) parseExistsExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume EXISTS

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: after EXISTS (, offer SELECT keyword.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		return nil, &ParseError{Message: "collecting"}
	}

	if !p.isQueryExpressionStart() {
		return nil, p.syntaxErrorAtCur()
	}
	sub, err := p.parseSubqueryExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.ExistsExpr{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Select: sub.Select,
	}, nil
}

// parseIntervalExpr parses INTERVAL value unit.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/expressions.html#temporal-intervals
//
//	INTERVAL expr unit
//
// The value is a full expression (bit_expr in MySQL yacc), so arithmetic like
// `INTERVAL DAY(d) - 1 DAY` and `INTERVAL @n + 1 HOUR` is valid. Unit
// keywords (DAY, HOUR, ...) are not infix operators, so Pratt parsing halts
// naturally before the unit.
func (p *Parser) parseIntervalExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume INTERVAL

	// INTERVAL(expr, expr, ...) is a function, not the interval literal. The
	// grammar has both `INTERVAL_SYM '(' expr ',' expr_list ')'` and
	// `INTERVAL_SYM expr interval`, whose value may be parenthesized, so the
	// comma tells them apart: mysql:8.0 parses `select a + interval (1) day from t`
	// as the literal and answers 1064 for `select interval(1) day`.
	if p.cur.Type == '(' && p.parenGroupHasComma() {
		return p.parseFuncCall(start, "", "INTERVAL")
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Parse unit: DAY, HOUR, MINUTE, SECOND, MONTH, YEAR, etc.
	// Use parseKeywordOrIdent because compound interval units (DAY_HOUR, etc.)
	// are registered as reserved keywords. SQL_TSI_ synonyms fold to their
	// canonical unit, matching the engine's stored form (oracle 8.0.32 +
	// 5.7.25: INTERVAL 1 SQL_TSI_DAY stores as interval 1 day).
	unit, _, err := p.parseKeywordOrIdent()
	if err != nil {
		return nil, err
	}

	upper := normalizeTimeUnit(unit)
	if !isValidIntervalUnit(upper) {
		return nil, &ParseError{
			Position: p.pos(),
			Message:  "invalid INTERVAL unit: " + unit,
		}
	}

	return &nodes.IntervalExpr{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Value: val,
		Unit:  upper,
	}, nil
}

// validIntervalUnits is the set of valid MySQL INTERVAL unit keywords.
var validIntervalUnits = map[string]bool{
	"MICROSECOND":        true,
	"SECOND":             true,
	"MINUTE":             true,
	"HOUR":               true,
	"DAY":                true,
	"WEEK":               true,
	"MONTH":              true,
	"QUARTER":            true,
	"YEAR":               true,
	"SECOND_MICROSECOND": true,
	"MINUTE_MICROSECOND": true,
	"MINUTE_SECOND":      true,
	"HOUR_MICROSECOND":   true,
	"HOUR_SECOND":        true,
	"HOUR_MINUTE":        true,
	"DAY_MICROSECOND":    true,
	"DAY_SECOND":         true,
	"DAY_MINUTE":         true,
	"DAY_HOUR":           true,
	"YEAR_MONTH":         true,
}

// isValidIntervalUnit returns true if unit (already uppercased) is a valid MySQL interval unit.
func isValidIntervalUnit(unit string) bool {
	return validIntervalUnits[unit]
}

// parseMatchExpr parses MATCH (col_list) AGAINST (expr [modifier]).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/fulltext-search.html
//
//	MATCH (col1 [, col2, ...]) AGAINST (expr [search_modifier])
//
//	search_modifier:
//	  {
//	       IN NATURAL LANGUAGE MODE
//	     | IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION
//	     | IN BOOLEAN MODE
//	     | WITH QUERY EXPANSION
//	  }
func (p *Parser) parseMatchExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume MATCH

	// ident_list_arg: the column list, in parentheses or not.
	//
	// Ref: simple_expr's `MATCH ident_list_arg AGAINST` and ident_list_arg in
	// MySQL 8.0 sql/sql_yacc.yy.
	_, parenthesized := p.match('(')

	me := &nodes.MatchExpr{Loc: nodes.Loc{Start: start}}

	// Column list
	for {
		ref, err := p.parseColumnRef()
		if err != nil {
			return nil, err
		}
		me.Columns = append(me.Columns, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if parenthesized {
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(kwAGAINST); err != nil {
		return nil, err
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Parse against expression with high enough precedence to avoid consuming
	// IN (which is precComparison) as part of the expression — IN here starts
	// the search modifier, not an IN-list comparison.
	against, err := p.parseExprPrec(precBitOr)
	if err != nil {
		return nil, err
	}
	me.Against = against

	// Search modifier
	if p.cur.Type == kwIN {
		p.advance()
		if p.cur.Type == kwNATURAL {
			// IN NATURAL LANGUAGE MODE [WITH QUERY EXPANSION]
			p.advance() // consume NATURAL
			if _, err := p.expect(kwLANGUAGE); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwMODE); err != nil {
				return nil, err
			}
			// Check for optional WITH QUERY EXPANSION
			if p.cur.Type == kwWITH && p.peekNext().Type == kwQUERY {
				p.advance() // consume WITH
				p.advance() // consume QUERY
				if _, err := p.expect(kwEXPANSION); err != nil {
					return nil, err
				}
				me.Modifier = "IN NATURAL LANGUAGE MODE WITH QUERY EXPANSION"
			} else {
				me.Modifier = "IN NATURAL LANGUAGE MODE"
			}
		} else {
			// IN BOOLEAN MODE
			if _, err := p.expect(kwBOOLEAN); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwMODE); err != nil {
				return nil, err
			}
			me.Modifier = "IN BOOLEAN MODE"
		}
	} else if p.cur.Type == kwWITH && p.peekNext().Type == kwQUERY {
		// WITH QUERY EXPANSION
		p.advance() // consume WITH
		p.advance() // consume QUERY
		if _, err := p.expect(kwEXPANSION); err != nil {
			return nil, err
		}
		me.Modifier = "WITH QUERY EXPANSION"
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	me.Loc.End = p.prev.End
	return me, nil
}

// parseSoundsLikeExpr parses expr SOUNDS LIKE expr.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html
//
//	expr1 SOUNDS LIKE expr2
//
// This is equivalent to SOUNDEX(expr1) = SOUNDEX(expr2).
func (p *Parser) parseSoundsLikeExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	p.advance() // consume SOUNDS

	if _, err := p.expect(kwLIKE); err != nil {
		return nil, err
	}

	right, err := p.parseExprPrec(precBitOr)
	if err != nil {
		return nil, err
	}

	return &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Op:    nodes.BinOpSoundsLike,
		Left:  left,
		Right: right,
	}, nil
}

// parseRowConstructor parses a ROW(expr, expr, ...) constructor.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/row-subqueries.html
//
//	ROW(val1, val2, ..., valN)
func (p *Parser) parseRowConstructor() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume ROW

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	row := &nodes.RowExpr{Loc: nodes.Loc{Start: start}}

	for {
		expr, err := p.parseRowItem()
		if err != nil {
			return nil, err
		}
		row.Items = append(row.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// A row has at least two elements: ROW(1) is a 1064 (measured).
	if len(row.Items) < 2 {
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	row.Loc.End = p.prev.End
	return row, nil
}

// parseRowItem parses one element of a row constructor. It is an expression, so
// not the bare DEFAULT that parseExpr also reads, and not *.
//
// Ref: MySQL 8.0 sql/sql_yacc.yy: simple_expr's ROW_SYM '(' expr ',' expr_list ')'.
func (p *Parser) parseRowItem() (nodes.ExprNode, error) {
	start := p.pos()
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if err := checkRowItem(e, start); err != nil {
		return nil, err
	}
	return e, nil
}

// checkRowItem refuses what parseExpr reads and a row does not hold, at start, where
// the element began.
func checkRowItem(e nodes.ExprNode, start int) error {
	if d, ok := e.(*nodes.DefaultExpr); ok && d.Column == "" {
		return &ParseError{Message: `syntax error at or near "DEFAULT"`, Position: start}
	}
	return nil
}

// parseDefaultExpr parses DEFAULT or DEFAULT(col_name).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_default
//
//	DEFAULT
//	DEFAULT(col_name)
func (p *Parser) parseDefaultExpr() (nodes.ExprNode, error) {
	if p.noBareDefault && p.peekNext().Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	start := p.pos()
	p.advance() // consume DEFAULT

	de := &nodes.DefaultExpr{Loc: nodes.Loc{Start: start}}

	// Check for DEFAULT(col_name) form
	if p.cur.Type == '(' {
		p.advance() // consume '('
		col, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		de.Column = col
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	de.Loc.End = p.prev.End
	return de, nil
}

// parseValuesFunc parses VALUES(col_name) used in INSERT ... ON DUPLICATE KEY UPDATE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_values
//
//	VALUES(col_name)
func (p *Parser) parseValuesFunc() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume VALUES

	// VALUES without '(' is handled elsewhere (e.g., VALUES row constructor in INSERT)
	// Here we handle the function form VALUES(col_name)
	if p.cur.Type != '(' {
		// Not a function call — this shouldn't be reached from parsePrimaryExpr
		// since VALUES without '(' is not a valid expression
		return nil, &ParseError{
			Message:  "expected '(' after VALUES in expression context",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume '('

	col, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	fc := &nodes.FuncCallExpr{
		Loc:       nodes.Loc{Start: start, End: p.prev.End},
		Name:      "VALUES",
		HasParens: true,
		Args:      []nodes.ExprNode{col},
	}
	return fc, nil
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() ([]nodes.ExprNode, error) {
	var list []nodes.ExprNode
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		list = append(list, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// nameQuotedFunction reports a call whose name was backquoted, which is a stored function's
// and no built-in's where a partition clause is concerned.
func nameQuotedFunction(p *Parser, schema string) bool {
	return schema == "" && p.backtickQuotedAt(p.prev.Loc)
}

// parseODBCEscape reads `{ident expr}`, ODBC's escape. With d, t or ts and a
// string it is that DATE, TIME or TIMESTAMP literal; anything else is the
// expression itself.
//
// Ref: `'{' ident expr '}'` and PTI_odbc_date in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseODBCEscape() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume '{'
	word, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect('}'); err != nil {
		return nil, err
	}
	if lit, ok := expr.(*nodes.StringLit); ok {
		kind := map[string]string{"D": "DATE", "T": "TIME", "TS": "TIMESTAMP"}[strings.ToUpper(word)]
		if kind != "" {
			return &nodes.TemporalLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Type: kind, Value: lit.Value}, nil
		}
	}
	return expr, nil
}

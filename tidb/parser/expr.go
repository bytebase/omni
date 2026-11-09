package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Precedence levels for Pratt parsing (low to high).
//
// The levels follow pkg/parser/parser.y (TiDB v8.5.5), whose shape a flat
// precedence table cannot say on its own:
//
//	Expression:    Expression logOr Expression | Expression "XOR" Expression
//	             | Expression logAnd Expression | "NOT" Expression
//	             | singleAtIdentifier assignmentEq Expression
//	             | BoolPri IsOrNotOp {trueKwd|falseKwd|"UNKNOWN"} | BoolPri     (parser.y:6098)
//	BoolPri:       BoolPri IsOrNotOp "NULL" | BoolPri CompareOp PredicateExpr
//	             | BoolPri CompareOp AnyOrAll SubSelect
//	             | BoolPri CompareOp singleAtIdentifier assignmentEq PredicateExpr
//	             | PredicateExpr                                                (parser.y:6248)
//	PredicateExpr: BitExpr InOrNotOp (...) | BitExpr BetweenOrNotOp BitExpr "AND" PredicateExpr
//	             | BitExpr LikeOrNotOp SimpleExpr | BitExpr RegexpOrNotOp SimpleExpr
//	             | BitExpr memberof '(' SimpleExpr ')' | BitExpr                (parser.y:6385)
//
// So a predicate binds tighter than a comparison (`a = b IN (1)` is `a = (b IN
// (1))`), its left operand is a BitExpr and so never another predicate, IS TRUE
// is an Expression and takes no comparison after it, and NOT and := reach only
// as far down as an Expression goes: `1 + NOT 2` and `1 + @a := 2` are 1064.
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
// start is the offset of the left operand's first token, where every folded
// expression starts.
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
	// `@name := expr` is an expr here too, which is why `not @a := 1` and `1 or
	// @a := 1` parse and `1 + @a := 1` does not.
	if minPrec <= precNot && p.atUserVarAssign() {
		left, err := p.parseUserVarAssign(precAssign)
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
// operand, folding infix operators of at least minPrec. start is the offset of
// the left operand's first token, where every folded expression starts.
//
// level is the highest operator that may take left as its left operand; it
// starts at the top for a primary and drops as operators fold. That is what
// says the grammar's non-associative rules.
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

		case p.cur.Type == kwLIKE || p.cur.Type == kwILIKE ||
			(p.cur.Type == kwNOT && (p.peekNext().Type == kwLIKE || p.peekNext().Type == kwILIKE)):
			left, err = p.parseLikeExpr(left, start)
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
			// CollationName is StringName | binaryType (parser.y SimpleExpr
			// "COLLATE" CollationName), so the reserved word BINARY and a
			// string literal are names here as well as an identifier:
			// `x COLLATE BINARY` parses on TiDB v8.5.5.
			collation, err := p.parseCollationName()
			if err != nil {
				return nil, err
			}
			left = &nodes.CollateExpr{
				Loc:       nodes.Loc{Start: start, End: p.prev.End},
				Expr:      left,
				Collation: collation,
			}
			continue

		case p.cur.Type == tokJsonExtract || p.cur.Type == tokJsonUnquote:
			// -> and ->> are not general operators: SimpleIdent jss|juss stringLit
			// (parser.y:8052) takes a column reference on the left and a string
			// literal on the right. `'{}'->'$.a'`, `a->3` and `a->b` are 1064s on
			// v8.5.5.
			op := binOp
			if ref, ok := left.(*nodes.ColumnRef); !ok || ref.Star {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume -> or ->>
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			strTok := p.advance()
			left = &nodes.BinaryExpr{
				Loc:   nodes.Loc{Start: start, End: strTok.End},
				Op:    op,
				Left:  left,
				Right: &nodes.StringLit{Loc: nodes.Loc{Start: strTok.Loc, End: strTok.End}, Value: strTok.Str},
			}
			level = precJsonAccess
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

		// Quantified-subquery RHS: BoolPri CompareOp AnyOrAll SubSelect
		// (parser.y:6248). The quantifier is recorded on the SubqueryExpr so
		// consumers can tell it from a scalar subquery comparison.
		if isComparison && (p.cur.Type == kwALL || ((p.cur.Type == kwANY || p.cur.Type == kwSOME) && p.peekNext().Type == '(')) {
			// ANY and SOME are not reserved: `1 = any` compares with a column.
			quantTok := p.advance()
			if p.cur.Type != '(' {
				return nil, p.syntaxErrorAtCur()
			}
			subStart := p.pos()
			p.advance() // consume '('
			if !p.atQueryPrimaryStart() && (p.cur.Type != '(' || !p.lookaheadParenContentIsSubquery()) {
				return nil, p.syntaxErrorAtCur()
			}
			sub, err := p.parseSubqueryExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			sub.Loc.Start = subStart
			sub.Loc.End = p.prev.End
			sub.Quantifier = strings.ToUpper(quantTok.Str)
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

		var right nodes.ExprNode
		if isComparison && p.atUserVarAssign() {
			// `1 = @a := 2`: BoolPri CompareOp singleAtIdentifier assignmentEq
			// PredicateExpr (parser.y:6248), the one place below an Expression
			// that takes one.
			right, err = p.parseUserVarAssign(precPredicate)
		} else {
			right, err = p.parseExprPrec(prec + 1)
		}
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
		case kwIN, kwLIKE, kwILIKE, kwBETWEEN, kwREGEXP, kwRLIKE:
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
	case kwLIKE, kwILIKE, kwREGEXP, kwRLIKE, kwIN, kwBETWEEN:
		return precPredicate, 0, true
	case kwMEMBER:
		// MEMBER is not reserved here: `select 1 member` is an alias, and only
		// MEMBER OF is the operator (memberof, parser.y).
		if p.peekNext().Type == kwOF {
			return precPredicate, 0, true
		}
		return 0, 0, false

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
		// Unary plus is kept as an operator, as TiDB's grammar keeps it
		// ('+' SimpleExpr %prec neg, an UnaryOperationExpr): a view stores
		// SELECT +a as +`a`, named "+a" (verified on v8.5.5).
		return p.parseUnaryExpr(nodes.UnaryPlus)
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

	// A prefix operator takes a SimpleExpr (parser.y:7892-7965: `'!' SimpleExpr`,
	// `'~' SimpleExpr`, `'-' SimpleExpr`, `'+' SimpleExpr`, `"BINARY" SimpleExpr`,
	// all %prec neg), and a bare INTERVAL literal is no SimpleExpr — it is only an
	// operand of the binary +/- date-arithmetic productions (parser.y:7793, 7804,
	// 7815). TiDB v8.5.5 answers 1064 for `DO - INTERVAL 1 DAY` and
	// `DO + INTERVAL 1 DAY`, and parses `DO - INTERVAL 1 DAY + d`, where the
	// INTERVAL is the right operand of the binary minus, not the unary one.
	if isBareInterval(operand) {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.UnaryExpr{
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
		Op:      op,
		Operand: operand,
	}, nil
}

// foldAdjacentStringLits folds a run of adjacent quoted strings into lit, as
// TiDB's StringLiteral does (parser.y:7660, v8.5.5): 'a' 'b' is one literal with
// the value ab. Only bare strings join: an introduced string takes one
// ("UNDERSCORE_CHARSET" stringLit, parser.y:7599), so _utf8mb4'a' 'b' and N'a' 'b'
// stay the string a aliased b, and a string followed by an introducer or a hex
// literal is a 1064. FirstSegment keeps the first string, which names the column
// (SetProjectionOffset in the same action): SELECT 'a' 'b' is column a, value ab.
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
		return &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Ival}, nil

	case tokFCONST:
		tok := p.advance()
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

	case kwNULL:
		tok := p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}}, nil

	case kwDEFAULT:
		return p.parseDefaultExpr()

	case kwCURRENT_DATE, kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwCURRENT_USER, kwLOCALTIME, kwLOCALTIMESTAMP,
		kwUTC_DATE, kwUTC_TIME, kwUTC_TIMESTAMP:
		return p.parseTemporalKeyword()

	case '*':
		tok := p.advance()
		return &nodes.StarExpr{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}}, nil

	// paramMarker is a SimpleExpr alternative (parser.y:7886), so a `?`
	// placeholder is a valid operand anywhere an expression is: a bare `DO ?`,
	// the offset of LEAD/LAG (OptLeadLagInfo, parser.y:9827, read as an ordinary
	// argument), and a window-frame bound (WindowFrameStart/Bound, parser.y:9700
	// and 9729, read through parseFrameBound's expr). The catalog analyzer
	// boundary rejects any statement carrying an unbound marker.
	case '?':
		tok := p.advance()
		return &nodes.ParamMarkerExpr{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}}, nil

	case '(':
		return p.parseParenExpr()

	case '{':
		return p.parseODBCEscape()

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

	// NEXT VALUE FOR, NextValueForSequence in parser.y. VALUE after NEXT commits to this form:
	// `SELECT next value FROM t` is a 1064 rather than a column with an alias.
	case kwNEXT:
		if p.peekNext().Type == kwVALUE {
			return p.parseNextValueFor()
		}
		return p.parseIdentExpr()

	// Words the parser reads as function calls by token. Unless `(` follows they are
	// names, as in `SELECT count FROM t` and `SELECT now.a FROM t`, where TiDB does not
	// reserve them.
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
		kwCHAR, kwINSERT, kwLEFT, kwRIGHT, kwDATABASE, kwSCHEMA,
		kwMOD, kwREPEAT, kwREPLACE:
		if p.funcWordIsName() {
			return p.parseIdentExpr()
		}
		start := p.pos()
		tok := p.advance()
		return p.parseFuncCall(start, "", strings.ToUpper(tok.Str))

	default:
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

		// N'text' is a string in the national character set, utf8. TiDB's lexer
		// reads N or n directly before a quote as the introducer _utf8
		// (startWithNn, pkg/parser/lexer.go), so N 'x' is still the column N
		// with an alias, and N'x' 'y' is the string with the alias y, as
		// _utf8'x' 'y' is: the introducer takes one string (Literal,
		// "UNDERSCORE_CHARSET" stringLit, parser.y).
		if p.atNationalString() {
			nTok := p.advance() // consume N
			strTok := p.advance()
			return &nodes.StringLit{
				Loc:     nodes.Loc{Start: nTok.Loc, End: strTok.End},
				Value:   strTok.Str,
				Charset: "_utf8",
			}, nil
		}

		// Charset introducer: _utf8mb4'hello', _latin1'world'
		if p.cur.Type == tokIDENT && strings.HasPrefix(p.cur.Str, "_") && p.peekNext().Type == tokSCONST {
			charsetTok := p.advance() // consume the charset identifier
			strTok := p.advance()     // consume the string literal
			return &nodes.StringLit{
				Loc:     nodes.Loc{Start: charsetTok.Loc, End: strTok.End},
				Value:   strTok.Str,
				Charset: charsetTok.Str,
			}, nil
		}

		// Charset introducer before a hex or bit literal: _utf8 0x61, _utf8 X'61',
		// _utf8 0b1, _utf8 B'1' ("UNDERSCORE_CHARSET" hexLit | bitLit,
		// parser.y:7625 and :7642). The introducer must name a charset TiDB's lexer
		// knows, which is how the lexer tells it from an identifier: _notacharset is
		// a column, and `_notacharset 0x61` is a 1064 (a column then a bare literal),
		// so isCharsetIntroducerWord gates it. The 1115 an action raises for a known
		// but unsupported charset is a name looked up while parsing, which this
		// parser does not reproduce.
		if isCharsetIntroducerWord(strings.ToLower(p.cur.Str)) && !p.backtickQuotedAt(p.cur.Loc) &&
			(p.peekNext().Type == tokXCONST || p.peekNext().Type == tokBCONST) {
			csTok := p.advance()
			litTok := p.advance()
			loc := nodes.Loc{Start: csTok.Loc, End: litTok.End}
			if litTok.Type == tokXCONST {
				return &nodes.HexLit{Loc: loc, Value: litTok.Str, Charset: csTok.Str}, nil
			}
			return &nodes.BitLit{Loc: loc, Value: litTok.Str, Charset: csTok.Str}, nil
		}

		// Identifier — could be column ref or function call
		if p.isIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// atNationalString reports whether the current token is the N of N'text': an
// unquoted N or n with a string literal directly after it.
func (p *Parser) atNationalString() bool {
	if p.cur.Type != tokIDENT || (p.cur.Str != "N" && p.cur.Str != "n") || p.backtickQuotedAt(p.cur.Loc) {
		return false
	}
	next := p.peekNext()
	return next.Type == tokSCONST && next.Loc == p.cur.End && p.lexer.inputAt(next.Loc) == '\''
}

// isSequenceNameArg reports whether a call has the one TableName argument that NEXTVAL
// and LASTVAL each take (FunctionNameSequence, parser.y:8437): one name, or two with a
// dot. `NEXTVAL(1)`, `LASTVAL()`, `LASTVAL(a, b)` and `NEXTVAL(a.b.c)` are 1064s.
func isSequenceNameArg(call nodes.ExprNode) bool {
	fn, ok := call.(*nodes.FuncCallExpr)
	if !ok || len(fn.Args) != 1 {
		return false
	}
	arg, ok := fn.Args[0].(*nodes.ColumnRef)
	return ok && arg.Schema == ""
}

// parseNextValueFor reads NEXT VALUE FOR TableName. It is the tree NEXTVAL(TableName) gives,
// as upstream's is. p.cur is NEXT.
func (p *Parser) parseNextValueFor() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume NEXT
	p.advance() // consume VALUE
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	argStart := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	seq := &nodes.ColumnRef{Column: name}
	if p.cur.Type == '.' {
		p.advance()
		qualified, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		seq.Table, seq.Column = name, qualified
	}
	seq.Loc = nodes.Loc{Start: argStart, End: p.prev.End}
	return &nodes.FuncCallExpr{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Name: "NEXTVAL",
		Args: []nodes.ExprNode{seq},
	}, nil
}

// funcWordIsName reports whether the function word at p.cur is a name here, which is
// when no `(` follows. parseIdentExpr then refuses the words TiDB reserves.
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
		call, err := p.parseFuncCall(start, "", name)
		if err != nil {
			return nil, err
		}
		// A quoted name is an ordinary call, which takes any arguments: `nextval`(1) and
		// `lastval`() parse. NEXTVAL and LASTVAL both take one TableName (parser.y:8437).
		if (strings.EqualFold(name, "nextval") || strings.EqualFold(name, "lastval")) &&
			!p.backtickQuotedAt(start) && !isSequenceNameArg(call) {
			return nil, &ParseError{Message: "sequence function takes a sequence name", Position: start}
		}
		return call, nil
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

		// table.*
		if p.cur.Type == '*' {
			p.advance()
			return &nodes.ColumnRef{
				Loc:   nodes.Loc{Start: start, End: p.prev.End},
				Table: name,
				Star:  true,
			}, nil
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

			if p.cur.Type == '*' {
				p.advance()
				return &nodes.ColumnRef{
					Loc:    nodes.Loc{Start: start, End: p.prev.End},
					Schema: name,
					Table:  name2,
					Star:   true,
				}, nil
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
	i := off - p.lexer.baseOffset
	return i >= 0 && i < len(p.lexer.input) && p.lexer.input[i] == '`'
}

// parseFuncCall parses a function call after the name has been consumed.
func (p *Parser) parseFuncCall(start int, schema, name string) (nodes.ExprNode, error) {
	// The builtins below take their special argument syntax (COUNT(*),
	// DISTINCT, TRIM's FROM, SUBSTRING's FROM and FOR, GROUP_CONCAT's ORDER BY
	// and SEPARATOR, the keyword arguments, OVER) from the grammar's keyword for
	// the name, and only the unqualified, unquoted name is that keyword. A
	// schema-qualified name is a stored function and a backtick-quoted name an
	// identifier: either way the call is an ordinary one, with ordinary
	// expressions for arguments and no window. TiDB v8.5.5 runs `trim`(b, 'x')
	// and `timestampdiff`(1, 2, 3), answers 1305 for `count`(a) and
	// test.trim(b, 'x'), and answers 1064 for `count`(*),
	// `trim`(LEADING 'x' FROM b), test.substring(b FROM 1), `sum`(a) OVER ()
	// and the other special forms under a quoted or qualified name. p.prev is
	// the name's token.
	builtin := schema == "" && !p.backtickQuotedAt(p.prev.Loc)
	if builtin {
		ordinary, err := p.checkSpacedFuncWord(strings.ToUpper(name))
		if err != nil {
			return nil, err
		}
		builtin = !ordinary
	}
	p.advance() // consume '('

	fc := &nodes.FuncCallExpr{
		Loc:       nodes.Loc{Start: start},
		Schema:    schema,
		Name:      strings.ToUpper(name),
		HasParens: true,
	}

	// Handle special function forms
	upperName := strings.ToUpper(name)

	if builtin && upperName == "SETVAL" {
		return p.parseSetvalArgs(fc)
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
		// identifier — a backtick-quoted unit is a syntax error (oracle mysql
		// 8.0).
		case "TIMESTAMPDIFF", "TIMESTAMPADD":
			return p.parseTimestampFunc(fc)
		case "GET_FORMAT":
			return p.parseGetFormatFunc(fc)
		case "MOD":
			// MOD '(' BitExpr ',' BitExpr ')' (parser.y:8221): exactly two
			// bit_expr arguments, so `MOD(10, 1 AND 1)` and
			// `MOD(10, 2 BETWEEN 0 AND 5)` are 1064s, where a general argument
			// list would take them.
			return p.parseModFunc(fc)
		case "WEIGHT_STRING":
			// WEIGHT_STRING(expr [AS CHAR(n) | AS BINARY(n)]) — the AS clause
			// is weight_string's own grammar, not a general cast.
			return p.parseWeightStringFunc(fc)
		}

		// An aggregate's arguments may carry a quantifier prefix, and which
		// prefixes and how many arguments is per-aggregate (see funcargs.go).
		// A builtin that is not an aggregate takes neither: TiDB v8.5.5
		// answers 1064 for `abs(distinct c1)`.
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
		done, err := p.parseTemporalCall(fc)
		if done {
			if err != nil {
				return nil, err
			}
			return fc, nil
		}
		// NOW reads an expression list, so it comes down to the ordinary reader
		// below, but its production has no window.
		if upperName == "NOW" {
			builtin = false
		}
	}

	// Empty argument list. CHAR's production takes an ExpressionList, which
	// is never empty: CHAR() is a 1064 (parser.y:8184).
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
		// INTERVAL is: TiDB v8.5.5 answers 1064 for abs(*), `count`(*) and
		// `date_add`(d, INTERVAL 1 DAY), and parses `f`(d + INTERVAL 1 DAY).
		// COUNT is the one builtin whose grammar has a star, in its own
		// production (see parseCountFunc).
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil || isStarArg(arg) {
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

	// CHAR(n, ... USING charset_name). Ref: TiDB v8.5.5 parser.y:8193,
	// "CHAR" '(' ExpressionList "USING" CharsetName ')'.
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
	// A builtin whose dedicated production fixes a minimum argument count is a 1064
	// below it: APPROX_COUNT_DISTINCT and APPROX_PERCENTILE take ExpressionList, which
	// is never empty (parser.y:8470,8474), so a zero-argument call is rejected on TiDB
	// v8.5.5. A quoted or schema-qualified name is an ordinary call (builtin false),
	// which carries no such floor. Reached through the generic argument path, so
	// unlike windowFuncArity this covers only the empty-list case.
	if builtin {
		if min, ok := builtinMinArgs[fc.Name]; ok && len(fc.Args) < min {
			return nil, p.syntaxErrorAtCur()
		}
	}

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
	// The INTERVAL function is the one builtin reached here whose production
	// carries no windowing clause, so a window after it is a syntax error:
	// TiDB v8.5.5 answers 1064 for `interval(c1) over ()`.
	// A window value function reads its null treatment, and NTH_VALUE its FROM
	// FIRST/LAST, between the argument list and OVER (parser.y:9802-9813). Each
	// belongs to that WindowFuncCall arm, so OVER must follow, which
	// requireOverAfterWindowOpt enforces.
	if builtin {
		if err := p.parseWindowFuncOpts(fc); err != nil {
			return nil, err
		}
		// A dedicated window function takes a fixed argument count in its
		// WindowFuncCall arm (parser.y:9761-9813), so a call outside those
		// bounds is a 1064 before OVER is even reached: NTILE(), NTILE(a, b),
		// LEAD() and FIRST_VALUE() are all rejected on TiDB v8.5.5.
		if bounds, ok := windowFuncArity[fc.Name]; ok {
			if len(fc.Args) < bounds[0] || len(fc.Args) > bounds[1] {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}
	// OVER is grammatical only on a dedicated window function or an aggregate
	// (overEligible). On any other builtin, or a stored/quoted call (builtin is
	// false), OVER is a 1064: TiDB v8.5.5 rejects `rate_limit() over w`,
	// `interval(c1) over ()` and `` `sum`(a) over () ``.
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

// parseTrimFunc parses TRIM([LEADING|TRAILING|BOTH] [remstr FROM] str).
// parseMemberOfExpr parses value MEMBER OF(json_array).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html
//
//	value MEMBER OF(json_array)
func (p *Parser) parseMemberOfExpr(value nodes.ExprNode, start int) (nodes.ExprNode, error) {
	p.advance() // consume MEMBER

	p.advance() // consume OF; infixPrecedence only answers for MEMBER OF

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	array, err := p.parseExprPrec(precCollate) // a SimpleExpr: `MEMBER OF (1+1)` is 1064
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

// nullTreatmentFuncs are the window functions that read OptNullTreatment
// (RESPECT|IGNORE NULLS, parser.y:9845): the five value functions. FROM
// FIRST/LAST is NTH_VALUE's alone.
var nullTreatmentFuncs = map[string]bool{
	"LEAD": true, "LAG": true, "FIRST_VALUE": true, "LAST_VALUE": true, "NTH_VALUE": true,
}

// windowFuncArity is the inclusive [min, max] argument count each dedicated
// window function takes in its WindowFuncCall arm (parser.y:9761-9813, tag
// v8.5.5): the rank family takes none, NTILE/FIRST_VALUE/LAST_VALUE one,
// NTH_VALUE two, and LEAD/LAG one to three (Expression OptLeadLagInfo, where
// OptLeadLagInfo adds `, NumLiteral` and then `, Expression`). TiDB v8.5.5
// answers 1064 for a call outside these bounds, e.g. NTILE(), NTILE(a, b),
// LEAD() and FIRST_VALUE().
// builtinMinArgs is the minimum argument count a non-window builtin's dedicated
// production requires, for the ones reached through the generic argument path (the
// window functions are in windowFuncArity and the aggregates have their own readers).
// Each takes ExpressionList, never empty (parser.y SumExpr), so finishFuncCall rejects
// a zero-argument call. TiDB v8.5.5 answers 1064 for APPROX_COUNT_DISTINCT() and
// APPROX_PERCENTILE().
var builtinMinArgs = map[string]int{
	"APPROX_COUNT_DISTINCT": 1,
	"APPROX_PERCENTILE":     1,
}

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
// the dedicated window functions (WindowFuncCall, parser.y:9761) and the
// aggregate functions (SumExpr: every arm ends in OptWindowingClause,
// parser.y:8461) carry a windowing clause. Any other builtin with OVER —
// INTERVAL, GROUPING, a generic builtin — is a 1064 on TiDB v8.5.5, e.g.
// `rate_limit() over w` and `interval(c1) over ()`.
func overEligible(name string) bool {
	if _, ok := windowFuncArity[name]; ok {
		return true
	}
	if _, ok := aggArgSpecs[name]; ok {
		return true
	}
	return name == "COUNT"
}

// parseWindowFuncOpts reads a window value function's FROM FIRST/LAST (NTH_VALUE
// only) and null treatment, in that order, both before OVER. A function that takes
// neither reads neither, so `COUNT(x) IGNORE NULLS` stays a 1064. Once either is
// read OVER must follow, since these live only on the WindowFuncCall arms.
func (p *Parser) parseWindowFuncOpts(fc *nodes.FuncCallExpr) error {
	if !nullTreatmentFuncs[fc.Name] {
		return nil
	}
	// NTH_VALUE takes exactly two arguments (parser.y:9810): `NTH_VALUE(v)` and
	// `NTH_VALUE(v, 1, 2)` are 1064s on v8.5.5.
	if fc.Name == "NTH_VALUE" && len(fc.Args) != 2 {
		return p.syntaxErrorAtCur()
	}
	read := false
	if fc.Name == "NTH_VALUE" && p.cur.Type == kwFROM {
		switch p.peekNext().Type {
		case kwFIRST:
			p.advance()
			p.advance()
			read = true
		case kwLAST:
			p.advance()
			p.advance()
			fc.FromLast = true
			read = true
		}
	}
	switch p.cur.Type {
	case kwRESPECT, kwIGNORE:
		if p.peekNext().Type != kwNULLS {
			return p.syntaxErrorAtTok(p.peekNext())
		}
		fc.IgnoreNull = p.cur.Type == kwIGNORE
		p.advance()
		p.advance()
		read = true
	}
	if read && p.cur.Type != kwOVER {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseModFunc parses MOD '(' BitExpr ',' BitExpr ')'. p.cur is just past '('.
func (p *Parser) parseModFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	a, err := p.parseBitExpr()
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	b, err := p.parseBitExpr()
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, p.syntaxErrorAtCur()
	}
	fc.Args = []nodes.ExprNode{a, b}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return p.finishFuncCall(fc, true)
}

// parseWeightStringFunc parses weight_string '(' Expression
// [ "AS" Char FieldLen | "AS" "BINARY" FieldLen ] ')' (FunctionCallNonKeyword,
// parser.y:8367,8374,8381). Char is CHAR | CHARACTER and FieldLen is
// '(' LengthNum ')'. The AS clause is optional and names no general cast type,
// so `WEIGHT_STRING(c AS DECIMAL(1))` and a bare `WEIGHT_STRING(c AS CHAR)`
// without the length are 1064s on TiDB v8.5.5. p.cur is just past '('.
func (p *Parser) parseWeightStringFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Completion parity with the generic argument path this intercepts.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}
	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}
	fc.Args = []nodes.ExprNode{arg}

	if p.cur.Type == kwAS {
		p.advance() // consume AS
		switch p.cur.Type {
		case kwCHAR, kwCHARACTER:
			fc.WeightAs = "CHAR"
		case kwBINARY:
			fc.WeightAs = "BINARY"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume CHAR/CHARACTER/BINARY
		// FieldLen is required: it is '(' LengthNum ')', not OptFieldLen.
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		n, err := p.parseLengthNum()
		if err != nil {
			return nil, err
		}
		fc.WeightAsLen = n
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return p.finishFuncCall(fc, true)
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
		// CURRENT ROW: ROW is mandatory, so `ROWS CURRENT` is a 1064
		// (parser.y:9709-9712, `"CURRENT" "ROW"`).
		p.advance()
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
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
		// The offset is exactly a NumLiteral (intLit | floatLit | decLit), a
		// paramMarker, or INTERVAL Expression TimeUnit — never an arbitrary
		// expression, so a bare keyword/identifier offset like `UNCOMMITTED
		// PRECEDING` is a 1064 (parser.y:9691-9707 WindowFrameStart,
		// 9719-9734 WindowFrameBound). omni's lexer has no separate decimal
		// token: a decimal lexes as tokFCONST and an integer (oversized ones
		// included) as tokICONST, so those two cover NumLiteral.
		switch p.cur.Type {
		case kwINTERVAL:
			expr, err := p.parseIntervalExpr()
			if err != nil {
				return nil, err
			}
			bound.Offset = expr
		case tokICONST:
			tok := p.advance()
			bound.Offset = &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Ival}
		case tokFCONST:
			tok := p.advance()
			bound.Offset = &nodes.FloatLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Str}
		case '?':
			tok := p.advance()
			bound.Offset = &nodes.ParamMarkerExpr{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}}
		default:
			return nil, p.syntaxErrorAtCur()
		}
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
// Ref: TiDB v8.5.5 parser.y:8336-8366:
//
//	builtinTrim '(' Expression ')'
//	builtinTrim '(' Expression "FROM" Expression ')'
//	builtinTrim '(' TrimDirection "FROM" Expression ')'
//	builtinTrim '(' TrimDirection Expression "FROM" Expression ')'
//
// A direction always takes FROM: TRIM(LEADING 'x') is a 1064.
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

// parseGroupConcatFunc parses GROUP_CONCAT([DISTINCT] expr [, expr ...] [ORDER BY ...] [SEPARATOR str]).
func (p *Parser) parseGroupConcatFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Parse arguments. The grammar takes an ExpressionList, so group_concat()
	// is a syntax error (TiDB v8.5.5: 1064).

	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil || isStarArg(arg) {
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
	// The production ends in OptWindowingClause (parser.y), so the call takes
	// a window: TiDB v8.5.5, MySQL 8.0 and MariaDB 11.8 all parse
	// `group_concat(c1) over ()` and answer 1235, the unsupported-feature
	// error, after it.
	return p.finishFuncCall(fc, true)
}

// lookaheadParenContentIsSubquery reports whether the content just inside an
// already-consumed '(' (p.cur is at the content's leading '(') is a query
// expression (subquery) rather than an expression / value list. It classifies
// WITHOUT consuming. Only invoked when p.cur == '(' — simple (SELECT ...)/(WITH
// ...) bodies are caught by the cheaper peek before this. Ported from
// pg/parser/expr.go lookaheadParenContentIsSubquery.
//
// The deciding token is the first one at depth 0 (relative to the consumed
// outer '(') past the balanced groups the content opens with. A ',' commits to
// a value list; a set-op or ORDER/LIMIT commits to a subquery; the matching ')'
// means a single element, a subquery iff the first token past its leading run
// of '(' starts a SELECT.
func (p *Parser) lookaheadParenContentIsSubquery() bool {
	if p.cur.Type != '(' {
		return false
	}
	next := p.cur
	for next.Type == '(' {
		after, ok := p.parenAfter[next.Loc]
		if !ok {
			next = p.scanParenGroups()
			break
		}
		next = after
	}
	switch next.Type {
	case ')':
		// Single element: a subquery iff that element's block opens one, so
		// `(((SELECT 1) + 0))` is an expression. FOR is excluded (TiDB rejects
		// `IN ((SELECT 1) FOR UPDATE)`).
		return p.parenBeginsSubquery()
	case kwUNION, kwINTERSECT, kwEXCEPT, kwORDER, kwLIMIT, kwFETCH:
		return true // query-expression continuation
	default:
		// A value list, e.g. IN ((SELECT 1), (SELECT 2)); an ordinary
		// expression, e.g. ((SELECT 1) + 1); or EOF before the matching ')',
		// where the value-list branch reports the error.
		return false
	}
}

// scanParenGroups returns the first token at depth 0 past the balanced groups
// that start at p.cur, without consuming input. On the way it records in
// parenAfter, for every '(' it passes, the token after its matching ')' (EOF
// for a '(' never closed).
func (p *Parser) scanParenGroups() Token {
	snap := p.snapshotTokenStream()
	defer p.restoreTokenStream(snap)
	if p.parenAfter == nil {
		p.parenAfter = make(map[int]Token)
	}
	var open []int
	for {
		switch {
		case p.cur.Type == '(':
			open = append(open, p.cur.Loc)
		case p.cur.Type == tokEOF:
			for _, loc := range open {
				p.parenAfter[loc] = p.cur
			}
			return p.cur
		case len(open) == 0:
			return p.cur
		case p.cur.Type == ')':
			loc := open[len(open)-1]
			open = open[:len(open)-1]
			p.advance()
			p.parenAfter[loc] = p.cur
			continue
		}
		p.advance()
	}
}

// parseParenExpr parses a parenthesized expression or subquery.
func (p *Parser) parseParenExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume '('

	// Check for subquery (incl. a parenthesized query expression like
	// ((SELECT 1) UNION (SELECT 2)) — classified by a backtracking scan).
	if p.atQueryPrimaryStart() ||
		(p.cur.Type == '(' && p.lookaheadParenContentIsSubquery()) {
		sub, err := p.parseSubqueryExpr()
		if err != nil {
			return nil, err
		}
		if _, errP := p.expect(')'); errP != nil {
			return nil, errP
		}
		sub.Loc.Start = start
		sub.Loc.End = p.prev.End
		return sub, nil
	}

	itemStart := p.pos()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// (a, b) is a row, as the corpus row Select (1, 1) > (1, 1) needs; every engine
	// in this family reads it, and TiDB v8.5.5 does (measured).
	if p.cur.Type == ',' {
		if err := checkRowItem(expr, itemStart); err != nil {
			return nil, err
		}
		row := &nodes.RowExpr{Loc: nodes.Loc{Start: start}, Items: []nodes.ExprNode{expr}}
		for p.cur.Type == ',' {
			p.advance()
			item, err := p.parseRowItem()
			if err != nil {
				return nil, err
			}
			row.Items = append(row.Items, item)
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		row.Loc.End = p.prev.End
		return row, nil
	}

	// A parenthesized expression is `'(' Expression ')'` (parser.y:7917); a bare
	// `*` or `t.*` is a Field wildcard (parser.y:6464-6473), never an Expression,
	// so TiDB v8.5.5 answers 1064 for `DO (*)` and `DO (t.*)`.
	if err := checkNotStar(expr); err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
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

	// ARRAY, the multi-valued index form: `builtinCast '(' Expression "AS"
	// CastType ArrayKwdOpt ')'` (pingcap/tidb pkg/parser/parser.y). Only CAST
	// takes it — TiDB v8.5.5 answers 1064 for convert(b, signed array) and
	// for a second ARRAY — and the server accepts the cast itself only inside
	// a functional index, which it reports as 1235 after parsing.
	array := false
	if p.cur.Type == kwARRAY {
		p.advance()
		array = true
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.CastExpr{
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
		Expr:     expr,
		TypeName: dt,
		Array:    array,
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

	// Parse the unit keyword. EXTRACT takes the full interval-unit set,
	// compound units (DAY_HOUR, YEAR_MONTH, ...) included — `builtinExtract
	// '(' TimeUnit "FROM" Expression ')'` (pingcap/tidb pkg/parser/parser.y).
	// Those are reserved keywords, so this reads a keyword token and not an
	// identifier: TiDB v8.5.5 answers 1064 for extract(`day_hour` from ...)
	// and for a unit that is not one, where an identifier would make the unit
	// a column reference and the statement mean something else.
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
		// CharsetName is StringName | binaryType (parser.y CONVERT "("
		// Expression "USING" CharsetName ")"): a string, an identifier or the
		// reserved word BINARY. CONVERT(x USING 'utf8mb4') is the form TiDB
		// writes into a view it stores, so it has to read back, and
		// CONVERT(x USING BINARY) parses on TiDB v8.5.5.
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

	// Check for subquery (incl. a parenthesized query expression like
	// ((SELECT 1) UNION (SELECT 2)) — classified by a backtracking scan; a
	// top-level comma instead routes to the value-list branch below).
	if p.atQueryPrimaryStart() ||
		(p.cur.Type == '(' && p.lookaheadParenContentIsSubquery()) {
		sub, err := p.parseSubqueryExpr()
		if err != nil {
			return nil, err
		}
		if _, errP := p.expect(')'); errP != nil {
			return nil, errP
		}
		inExpr.Select = sub.Select
		inExpr.Loc.End = p.prev.End
		return inExpr, nil
	}

	// Value list
	for {
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, p.syntaxErrorAtCur()
		}
		// The value list is an ExpressionList (parser.y:6385, `BitExpr InOrNotOp
		// '(' ExpressionList ')'`); a bare `*` or `t.*` is a Field wildcard
		// (parser.y:6464-6473), not an Expression, so TiDB v8.5.5 answers 1064 for
		// `a IN (*)`.
		if err := checkNotStar(val); err != nil {
			return nil, err
		}
		inExpr.List = append(inExpr.List, val)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	inExpr.Loc.End = p.prev.End
	return inExpr, nil
}

// parseLikeExpr parses [NOT] {LIKE | ILIKE} pattern [ESCAPE escape_char]. ILIKE is TiDB's, with
// the same shape as LIKE in PredicateExpr and a word of its own that v8.5.5 reserves.
func (p *Parser) parseLikeExpr(left nodes.ExprNode, start int) (nodes.ExprNode, error) {
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	ilike := p.cur.Type == kwILIKE
	p.advance() // consume LIKE or ILIKE

	pattern, err := p.parseExprPrec(precCollate)
	if err != nil {
		return nil, err
	}

	like := &nodes.LikeExpr{
		Loc:     nodes.Loc{Start: start},
		Not:     not,
		ILike:   ilike,
		Expr:    left,
		Pattern: pattern,
	}

	if _, ok := p.match(kwESCAPE); ok {
		// ESCAPE takes a string literal and nothing else (LikeOrIlikeEscapeOpt,
		// parser.y:6453): `ESCAPE concat('c')` is a 1064.
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		// PredicateExpr's action answers 1210, ErrWrongArguments, for more than one
		// byte: `ESCAPE '||'` and `ESCAPE '€'` are refused, `ESCAPE ''` is not.
		if len(p.cur.Str) > 1 {
			return nil, &ParseError{Message: "Incorrect arguments to ESCAPE", Position: p.cur.Loc}
		}
		// One stringLit, not a StringLiteral run: `ESCAPE 'x' 'y'` is a 1064 at
		// 'y', where parsePrimaryExpr would fold the run into the two-byte xy.
		tok := p.advance()
		like.Escape = &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: tok.Str}
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

	pattern, err := p.parseExprPrec(precCollate)
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
func (p *Parser) parseIntervalExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume INTERVAL

	// INTERVAL(expr, expr, ...) is a function, not the interval literal: TiDB
	// reads the name through FunctionNameConflict (pingcap/tidb
	// pkg/parser/parser.y), so a '(' here always opens its argument list and
	// no unit follows. TiDB v8.5.5 answers 1582 — the wrong-argument-count
	// error, past parsing — for `select a + interval (1) day from t`, where
	// MySQL 8.0 reads that as the literal. The arity is the function's own
	// business for the same reason: `interval()` parses too.
	if p.cur.Type == '(' {
		return p.parseFuncCall(start, "", "INTERVAL")
	}

	// The value may be another interval literal, so this cycle —
	// parseIntervalExpr, parsePrefixExpr, parseIntervalExpr — counts its own
	// nesting. Without it a deep enough chain overflows the Go stack, which
	// recover cannot catch, before the post-parse tree-height check runs.
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()

	// The leading operand is a primary, a signed numeric literal, or a nested
	// INTERVAL — the narrow set the value has always allowed, each of which
	// deparse writes back out bare. TiDB reads `INTERVAL (x) unit` as the
	// INTERVAL function (error 1582), never the literal, so a value that deparse
	// would wrap in parentheses, drop an operator from, or lead with a '('
	// (CASE, EXISTS, BINARY, COLLATE, a sign on anything but a number) stays
	// unread here even though the engine parses it. mysql, mariadb and oceanbase
	// take the whole expression, because their engines read a parenthesized
	// value as the literal.
	valStart := p.pos()
	var left nodes.ExprNode
	var err error
	switch p.cur.Type {
	case kwINTERVAL:
		left, err = p.parseIntervalExpr()
	case '-', '+':
		// A sign directly on a numeric literal is part of the value: TiDB v8.5.5
		// runs INTERVAL -1 DAY and INTERVAL -2.2 MINUTE_MICROSECOND. A signed
		// literal deparses without parentheses (`-2.2`, or `1` for a dropped unary
		// plus), so it stays the interval literal. A sign on anything else stays
		// narrowed: `interval (1 + 2) day` is the INTERVAL function to TiDB.
		//
		// The operand is read with parsePrimaryExpr, not parseUnaryExpr, so a
		// suffix that binds tighter than the sign (COLLATE) is not folded into it:
		// `INTERVAL +1 COLLATE binary DAY` leaves COLLATE for the unit, which
		// refuses it.
		if nt := p.peekNext().Type; nt != tokICONST && nt != tokFCONST {
			return nil, p.syntaxErrorAtCur()
		}
		op := nodes.UnaryMinus
		if p.cur.Type == '+' {
			op = nodes.UnaryPlus
		}
		opStart := p.pos()
		p.advance() // consume the sign
		var operand nodes.ExprNode
		operand, err = p.parsePrimaryExpr()
		if err == nil {
			left = &nodes.UnaryExpr{Loc: nodes.Loc{Start: opStart, End: p.prev.End}, Op: op, Operand: operand}
		}
	default:
		left, err = p.parsePrimaryExpr()
	}
	if err != nil {
		return nil, err
	}

	// Fold the bit_expr binary operators (|, &, <<, >>, +, -, *, /, %, DIV, MOD,
	// ^) that follow, so `INTERVAL 1 << 37 SECOND` and `INTERVAL 1 + 2 DAY` read
	// as the arithmetic value TiDB v8.5.5 stores (container-verified) and deparse
	// writes back out bare. Folding stops below COLLATE and the JSON accessors
	// and never reaches a comparison or logical operator, so the leading token
	// decides what the value may be; a later operand it folds in is never the
	// token the engine reads as `INTERVAL (`, so it may be richer.
	val, err := p.parseInfixExprPrec(left, valStart, precBitOr, precBitXor)
	if err != nil {
		return nil, err
	}

	// Parse unit: DAY, HOUR, MINUTE, SECOND, MONTH, YEAR, etc.
	// Use parseKeywordOrIdent because compound interval units (DAY_HOUR, etc.)
	// are registered as reserved keywords. SQL_TSI_ synonyms fold to their
	// canonical unit, matching the engine's stored form.
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

// sqlTSIUnits maps the SQL_TSI_-prefixed unit synonyms MySQL's lexer
// registers to their canonical unit. Only these eight exist —
// SQL_TSI_MICROSECOND is not a token and the engine rejects it
// (oracle mysql 8.0: error 1064).
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
// (oracle mysql 8.0).
func normalizeTimeUnit(unit string) string {
	u := strings.ToUpper(unit)
	if canon, ok := sqlTSIUnits[u]; ok {
		return canon
	}
	return u
}

// simpleTimeUnits is the unit set TIMESTAMPDIFF/TIMESTAMPADD accept
// (interval_time_stamp in sql_yacc.yy): the nine simple units only.
// Compound units are a syntax error there (oracle mysql 8.0:
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
// (oracle mysql 8.0) — so only keyword tokens are accepted; all valid
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
// literals — is a syntax error (oracle mysql 8.0).
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

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

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

	if _, err := p.expect(')'); err != nil {
		return nil, err
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
// Ref: TiDB parser.y (v8.5.5) line 7930: "ROW" '(' ExpressionList ',' Expression ')'.
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
	// A row element is an Expression (parser.y:7918, `'(' ExpressionList ','
	// Expression ')'`), and a bare INTERVAL literal is no Expression — it is only
	// an operand of the binary +/- date-arithmetic productions. TiDB v8.5.5
	// answers 1064 for `DO (d, INTERVAL 1 DAY)` and parses `DO (d, d + INTERVAL 1 DAY)`.
	if isBareInterval(e) {
		return &ParseError{Message: `syntax error at or near "INTERVAL"`, Position: start}
	}
	return checkNotStar(e)
}

// parseDefaultExpr parses DEFAULT or DEFAULT(col_name).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html#function_default
//
//	DEFAULT
//	DEFAULT(col_name)
func (p *Parser) parseDefaultExpr() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume DEFAULT

	de := &nodes.DefaultExpr{Loc: nodes.Loc{Start: start}}

	// Check for DEFAULT(col_name) form. The name may be qualified — TiDB v8.5.5
	// runs DEFAULT(t1.a) and DEFAULT(db.t1.a) — so read up to schema.table.column.
	if p.cur.Type == '(' {
		p.advance() // consume '('
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		de.Column = name
		for i := 0; i < 2 && p.cur.Type == '.'; i++ {
			p.advance() // consume '.'
			next, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			de.Schema, de.Table, de.Column = de.Table, de.Column, next
		}
		if p.cur.Type == '.' {
			// A fourth part is one qualifier too many.
			return nil, p.syntaxErrorAtCur()
		}
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

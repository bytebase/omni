package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// Precedence levels for Pratt parsing (low to high, matching MySQL).
const (
	precNone       = 0
	precAssign     = 1  // :=
	precOr         = 2  // OR, ||
	precXor        = 3  // XOR
	precAnd        = 4  // AND, &&
	precNot        = 5  // NOT (prefix)
	precComparison = 6  // =, <=>, >=, >, <=, <, <>, !=, IS, LIKE, REGEXP, IN, BETWEEN
	precBitOr      = 7  // |
	precBitAnd     = 8  // &
	precShift      = 9  // <<, >>
	precAdd        = 10 // +, -
	precMul        = 11 // *, /, DIV, %, MOD
	precBitXor     = 12 // ^
	precUnary      = 13 // -, ~, !
	precCollate    = 14 // COLLATE
	precJsonAccess = 15 // ->, ->>
)

// parseExpr parses an expression using Pratt parsing / precedence climbing.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/expressions.html
func (p *Parser) parseExpr() (nodes.ExprNode, error) {
	return p.parseExprPrec(precAssign)
}

// parseExprPrec parses an expression with the given minimum precedence.
func (p *Parser) parseExprPrec(minPrec int) (nodes.ExprNode, error) {
	left, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}
	return p.parseInfixExprPrec(left, minPrec)
}

// parseInfixExprPrec continues Pratt parsing with an already-parsed left
// operand, folding infix operators of at least minPrec. Split from
// parseExprPrec so callers that produce the leading primary themselves — a
// parenthesized query expression that turns out to be the left operand of a
// scalar expression, `((SELECT ...) = 0)` — can resume the operator loop.
func (p *Parser) parseInfixExprPrec(left nodes.ExprNode, minPrec int) (nodes.ExprNode, error) {
	var err error
	for {
		// MEMBER OF special handling (not in infixPrecedence since it's keyword-based)
		if p.cur.Type == kwMEMBER {
			left, err = p.parseMemberOfExpr(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		prec, binOp, ok := p.infixPrecedence()
		if !ok || prec < minPrec {
			break
		}

		// Handle special infix operators
		switch {
		case p.cur.Type == kwIS:
			left, err = p.parseIsExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwBETWEEN || (p.cur.Type == kwNOT && p.peekNext().Type == kwBETWEEN):
			left, err = p.parseBetweenExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwIN || (p.cur.Type == kwNOT && p.peekNext().Type == kwIN):
			left, err = p.parseInExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwLIKE || (p.cur.Type == kwNOT && p.peekNext().Type == kwLIKE):
			left, err = p.parseLikeExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwSOUNDS:
			left, err = p.parseSoundsLikeExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwREGEXP || p.cur.Type == kwRLIKE ||
			(p.cur.Type == kwNOT && (p.peekNext().Type == kwREGEXP || p.peekNext().Type == kwRLIKE)):
			left, err = p.parseRegexpExpr(left)
			if err != nil {
				return nil, err
			}
			continue

		case p.cur.Type == kwCOLLATE:
			colStart := p.pos()
			p.advance() // consume COLLATE
			collation, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			left = &nodes.CollateExpr{
				Loc:       nodes.Loc{Start: colStart, End: p.pos()},
				Expr:      left,
				Collation: collation,
			}
			continue
		}

		// Regular binary operator
		opStart := p.pos()
		// Capture original operator text for auto-alias when it differs from canonical form.
		var originalOp string
		switch {
		case p.cur.Type == kwMOD:
			originalOp = "MOD"
		case p.cur.Type == tokNotEq && p.cur.Str == "!=":
			originalOp = "!="
		}
		isComparison := prec == precComparison
		p.advance() // consume operator

		// Quantified-subquery RHS: `expr comp_op {ANY|SOME|ALL} (subquery)`.
		// Only after a comparison operator. The quantifier is recorded on
		// the SubqueryExpr so consumers can distinguish from a scalar
		// subquery comparison.
		if isComparison && (p.cur.Type == kwANY || p.cur.Type == kwSOME || p.cur.Type == kwALL) {
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
				Loc:        nodes.Loc{Start: opStart, End: p.pos()},
				Op:         binOp,
				Left:       left,
				Right:      sub,
				OriginalOp: originalOp,
			}
			continue
		}

		// Right-associative for assignment
		nextPrec := prec + 1
		if binOp == nodes.BinOpAssign {
			nextPrec = prec
		}

		right, err := p.parseExprPrec(nextPrec)
		if err != nil {
			return nil, err
		}

		left = &nodes.BinaryExpr{
			Loc:        nodes.Loc{Start: opStart, End: p.pos()},
			Op:         binOp,
			Left:       left,
			Right:      right,
			OriginalOp: originalOp,
		}
	}

	return left, nil
}

// infixPrecedence returns the precedence and operator for the current token
// if it is an infix operator.
func (p *Parser) infixPrecedence() (int, nodes.BinaryOp, bool) {
	switch p.cur.Type {
	case tokAssign:
		return precAssign, nodes.BinOpAssign, true
	case kwOR:
		return precOr, nodes.BinOpOr, true
	case kwXOR:
		return precXor, nodes.BinOpXor, true
	case kwAND:
		return precAnd, nodes.BinOpAnd, true

	// NOT as infix: NOT IN, NOT LIKE, NOT BETWEEN, NOT REGEXP
	case kwNOT:
		switch p.peekNext().Type {
		case kwIN, kwLIKE, kwBETWEEN, kwREGEXP, kwRLIKE:
			return precComparison, 0, true // handled specially
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
	case kwLIKE:
		return precComparison, 0, true
	case kwREGEXP, kwRLIKE:
		return precComparison, 0, true
	case kwIN:
		return precComparison, 0, true
	case kwBETWEEN:
		return precComparison, 0, true
	case kwSOUNDS:
		return precComparison, 0, true

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

// parsePrefixExpr parses prefix/primary expressions.
func (p *Parser) parsePrefixExpr() (nodes.ExprNode, error) {
	// Completion: at the start of any expression, offer column/function candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case '-':
		return p.parseUnaryExpr(nodes.UnaryMinus)
	case '+':
		// Unary plus — just parse the operand
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
	case kwNOT:
		return p.parseUnaryExpr(nodes.UnaryNot)
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

// parseUnaryExpr parses a unary expression.
func (p *Parser) parseUnaryExpr(op nodes.UnaryOp) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume operator

	operand, err := p.parseExprPrec(precUnary)
	if err != nil {
		return nil, err
	}

	return &nodes.UnaryExpr{
		Loc:     nodes.Loc{Start: start, End: p.pos()},
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
	lit.Loc.End = p.pos()
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
		return &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}, Value: tok.Ival, Text: tok.Str}, nil

	case tokFCONST:
		tok := p.advance()
		return &nodes.FloatLit{Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}, Value: tok.Str}, nil

	case tokSCONST:
		tok := p.advance()
		lit := &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}, Value: tok.Str}
		p.foldAdjacentStringLits(lit)
		return lit, nil

	case tokXCONST:
		tok := p.advance()
		return &nodes.HexLit{Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}, Value: tok.Str}, nil

	case tokBCONST:
		tok := p.advance()
		return &nodes.BitLit{Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}, Value: tok.Str}, nil

	case kwTRUE:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc}, Value: true}, nil

	case kwFALSE:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc}, Value: false}, nil

	case kwON:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc}, Value: true}, nil

	case kwOFF:
		tok := p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: tok.Loc}, Value: false}, nil

	case kwNULL:
		tok := p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: tok.Loc}}, nil

	case kwDEFAULT:
		return p.parseDefaultExpr()

	case kwCURRENT_DATE, kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwCURRENT_USER, kwLOCALTIME, kwLOCALTIMESTAMP,
		kwUTC_DATE, kwUTC_TIME, kwUTC_TIMESTAMP:
		tok := p.advance()
		name := strings.ToUpper(tok.Str)
		fc := &nodes.FuncCallExpr{Loc: nodes.Loc{Start: tok.Loc}, Name: name}
		// Optional parentheses with optional fsp argument
		if p.cur.Type == '(' {
			fc.HasParens = true
			p.advance()
			if p.cur.Type != ')' {
				// Parse arguments (e.g., fractional seconds precision)
				for {
					arg, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					if arg == nil {
						return nil, p.syntaxErrorAtCur()
					}
					fc.Args = append(fc.Args, arg)
					if p.cur.Type != ',' {
						break
					}
					p.advance() // consume ','
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		fc.Loc.End = p.prev.End
		return fc, nil

	case '*':
		tok := p.advance()
		return &nodes.StarExpr{Loc: nodes.Loc{Start: tok.Loc, End: p.pos()}}, nil

	case '(':
		return p.parseParenExpr()

	case kwCAST:
		return p.parseCastExpr()

	case kwEXTRACT:
		return p.parseExtractExpr()

	case kwCONVERT:
		return p.parseConvertExpr()

	case kwROW:
		return p.parseRowConstructor()

	case kwVALUES:
		return p.parseValuesFunc()

	// Window functions and GROUPING — reserved keywords used as function calls.
	// After registration as reserved keywords, they no longer match isIdentToken(),
	// so we dispatch them here to parseFuncCall explicitly.
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
		start := p.pos()
		tok := p.advance()
		name := strings.ToUpper(tok.Str)
		if p.cur.Type == '(' {
			return p.parseFuncCall(start, "", name)
		}
		// Not followed by '(' — syntax error (these are reserved and can't be identifiers)
		return nil, p.syntaxErrorAtCur()

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
				Loc:   nodes.Loc{Start: typeTok.Loc, End: valTok.Loc + len(valTok.Str) + 2},
				Type:  strings.ToUpper(typeTok.Str),
				Value: valTok.Str,
			}, nil
		}

		// Charset introducer: _utf8mb4'hello', _latin1'world'. Only a real
		// charset name forms an introducer — MySQL lexes _typo'abc' as an
		// identifier followed by a separate string literal (oracle 8.0.32 +
		// 5.7.25: SELECT _typo'abc' → 1054 unknown column, not a syntax
		// error), so an unknown name falls through to identifier parsing.
		if p.cur.Type == tokIDENT && isCharsetIntroducer(p.cur.Str) && p.peekNext().Type == tokSCONST {
			charsetTok := p.advance() // consume the charset identifier
			strTok := p.advance()     // consume the string literal
			lit := &nodes.StringLit{
				Loc:     nodes.Loc{Start: charsetTok.Loc, End: strTok.Loc + len(strTok.Str) + 2},
				Value:   strTok.Str,
				Charset: charsetTok.Str,
			}
			p.foldAdjacentStringLits(lit)
			return lit, nil
		}

		// Identifier — could be column ref or function call
		if p.isIdentToken() {
			return p.parseIdentExpr()
		}

		return nil, p.syntaxErrorAtCur()
	}
}

// parseIdentExpr parses an identifier that could be a column ref or function call.
func (p *Parser) parseIdentExpr() (nodes.ExprNode, error) {
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

		// table.*
		if p.cur.Type == '*' {
			p.advance()
			return &nodes.ColumnRef{
				Loc:   nodes.Loc{Start: start, End: p.pos()},
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
					Loc:    nodes.Loc{Start: start, End: p.pos()},
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
				Loc:    nodes.Loc{Start: start, End: p.pos()},
				Schema: name,
				Table:  name2,
				Column: name3,
			}, nil
		}

		// table.col
		return &nodes.ColumnRef{
			Loc:    nodes.Loc{Start: start, End: p.pos()},
			Table:  name,
			Column: name2,
		}, nil
	}

	// Plain column ref
	return &nodes.ColumnRef{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Column: name,
	}, nil
}

// parseFuncCall parses a function call after the name has been consumed.
func (p *Parser) parseFuncCall(start int, schema, name string) (nodes.ExprNode, error) {
	p.advance() // consume '('

	fc := &nodes.FuncCallExpr{
		Loc:       nodes.Loc{Start: start},
		Schema:    schema,
		Name:      strings.ToUpper(name),
		HasParens: true,
	}

	// Handle special function forms
	upperName := strings.ToUpper(name)

	// COUNT(*)
	if upperName == "COUNT" && p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		goto overCheck
	}

	// Aggregate with DISTINCT: COUNT(DISTINCT ...), SUM(DISTINCT ...), etc.
	if p.cur.Type == kwDISTINCT {
		fc.Distinct = true
		p.advance()
	}

	// TRIM special syntax: TRIM([LEADING|TRAILING|BOTH] [remstr FROM] str)
	if upperName == "TRIM" {
		return p.parseTrimFunc(fc)
	}

	// SUBSTRING special syntax: SUBSTRING(str, pos, len) or SUBSTRING(str FROM pos FOR len)
	if upperName == "SUBSTRING" || upperName == "SUBSTR" {
		return p.parseSubstringFunc(fc)
	}

	// GROUP_CONCAT special syntax
	if upperName == "GROUP_CONCAT" {
		return p.parseGroupConcatFunc(fc)
	}

	// Keyword-argument builtins: their first argument (or a trailing clause)
	// is a bare keyword the engine lexes as a token, never an identifier —
	// a backtick-quoted unit is a syntax error (oracle 8.0.32 + 5.7.25).
	// Only the unqualified builtin has this shape; a schema-qualified name is
	// a stored function whose arguments are ordinary expressions.
	if schema == "" {
		switch upperName {
		case "TIMESTAMPDIFF", "TIMESTAMPADD":
			return p.parseTimestampFunc(fc)
		case "GET_FORMAT":
			return p.parseGetFormatFunc(fc)
		case "WEIGHT_STRING":
			return p.parseWeightStringFunc(fc)
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

	// Empty argument list
	if p.cur.Type == ')' {
		p.advance()
		goto overCheck
	}

	// Regular argument list
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

overCheck:

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
	if p.cur.Type == kwOVER {
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
func (p *Parser) parseMemberOfExpr(value nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume MEMBER

	// Expect OF
	if p.cur.Type != kwOF {
		return nil, &ParseError{
			Message:  "expected OF after MEMBER",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume OF

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	array, err := p.parseExpr()
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
		Loc:   nodes.Loc{Start: start, End: p.pos()},
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
			Loc:     nodes.Loc{Start: start, End: p.pos()},
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

func (p *Parser) parseTrimFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Check for LEADING, TRAILING, BOTH
	switch p.cur.Type {
	case kwLEADING, kwTRAILING, kwBOTH:
		// Include as first "argument" indicator via the function name
		fc.Name = "TRIM_" + strings.ToUpper(p.cur.Str)
		p.advance()
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
		Loc:     nodes.Loc{Start: tok.Loc, End: p.pos()},
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
		Loc:     nodes.Loc{Start: tok.Loc, End: p.pos()},
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
//	WEIGHT_STRING(expr [AS {CHAR|BINARY}(n)] [LEVEL n [ASC|DESC|REVERSE] [, ...] | LEVEL n - n])
//
// The plain no-suffix form stays a generic FuncCallExpr (its round-trip
// predates this parser and is pinned by tests); AS/LEVEL forms build a
// WeightStringExpr. AS BINARY(n) is not stored by the engine — it desugars
// to weight_string(cast(expr as char(n) charset binary)) (oracle 8.0.32 +
// 5.7.25), so the parser produces that shape directly. LEVEL is 5.7-only
// syntax (8.0 rejects it); 5.7 stores a plain level list such as
// "level 1 desc", with ASC dropped and ranges collapsed.
//
// Ref: https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_weight-string
func (p *Parser) parseWeightStringFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	// Completion parity with the generic argument path.
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

	ws := &nodes.WeightStringExpr{Loc: fc.Loc, Expr: arg}
	asBinary := false

	if p.cur.Type == kwAS {
		p.advance()
		var typeName string
		switch p.cur.Type {
		case kwCHAR:
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
			Loc:    nodes.Loc{Start: typeTok.Loc, End: p.pos()},
			Name:   "CHAR",
			Length: int(lenTok.Ival),
		}
		if typeName == "BINARY" {
			// Engine-stored form of AS BINARY(n):
			// weight_string(cast(expr as char(n) charset binary)).
			dt.Charset = "binary"
			asBinary = true
			ws.Expr = &nodes.CastExpr{
				Loc:      nodes.Loc{Start: typeTok.Loc, End: p.pos()},
				Expr:     arg,
				TypeName: dt,
			}
		} else {
			ws.AsChar = dt
		}
	}

	if p.cur.Type == kwLEVEL {
		p.advance()
		for {
			if p.cur.Type != tokICONST {
				return nil, &ParseError{Message: "expected level number", Position: p.cur.Loc}
			}
			lvTok := p.advance()
			// The engine clamps out-of-range levels rather than rejecting:
			// LEVEL 0 stores as level 1 (oracle 5.7.25). The collation-max
			// upper clamp (LEVEL 7 also stored level 1 on a single-level
			// collation) is collation-dependent and stays a flagged
			// user-form residual.
			lvNum := int(lvTok.Ival)
			if lvNum < 1 {
				lvNum = 1
			}
			level := nodes.WeightStringLevel{Level: lvNum}
			// Grammar: n [ASC|DESC] [REVERSE] — direction first, then an
			// optional REVERSE (LEVEL 1 DESC REVERSE stores as
			// "level 1 desc reverse"; REVERSE DESC is a 1064 and ASC
			// normalizes away — oracle 5.7.25).
			switch p.cur.Type {
			case kwASC:
				p.advance() // ASC is the default; normalize away
			case kwDESC:
				p.advance()
				level.Desc = true
			}
			if p.cur.Type == kwREVERSE {
				p.advance()
				level.Reverse = true
			}
			ws.Levels = append(ws.Levels, level)
			// Range form LEVEL n - m: expand to the full list. Weight levels
			// are 1..6 per the MySQL reference (no collation has more), so
			// the expansion is capped — an absurd range like LEVEL 1 - 1e18
			// parses without materializing entries beyond the domain.
			if len(ws.Levels) == 1 && !level.Desc && !level.Reverse && p.cur.Type == '-' {
				p.advance()
				if p.cur.Type != tokICONST {
					return nil, &ParseError{Message: "expected level number", Position: p.cur.Loc}
				}
				hiTok := p.advance()
				const maxWeightLevels = 6
				for l := level.Level + 1; l <= int(hiTok.Ival) && len(ws.Levels) < maxWeightLevels; l++ {
					ws.Levels = append(ws.Levels, nodes.WeightStringLevel{Level: l})
				}
				break
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if ws.AsChar == nil && !asBinary && len(ws.Levels) == 0 {
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
	if p.cur.Type == ')' {
		p.advance()
		fc.Loc.End = p.prev.End
		return fc, nil
	}

	// Parse arguments
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
	fc.Loc.End = p.prev.End
	return fc, nil
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
			expr, err = p.parseInfixExprPrec(sub, precAssign)
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
		row := &nodes.RowExpr{Loc: nodes.Loc{Start: start}, Items: []nodes.ExprNode{expr}}
		for p.cur.Type == ',' {
			p.advance()
			item, err := p.parseExpr()
			if err != nil {
				return nil, furthestErr(err, specErr)
			}
			if item == nil {
				return nil, p.syntaxErrorAtCur()
			}
			row.Items = append(row.Items, item)
		}
		if _, err := p.expect(')'); err != nil {
			return nil, furthestErr(err, specErr)
		}
		row.Loc.End = p.pos()
		return row, nil
	}

	if _, err := p.expect(')'); err != nil {
		return nil, furthestErr(err, specErr)
	}

	return &nodes.ParenExpr{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
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
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.CastExpr{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
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
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Unit: unit.Keyword,
		Expr: expr,
	}, nil
}

// parseCastDataType parses data types valid in CAST expressions.
// CAST supports a subset of types: BINARY, CHAR, DATE, DATETIME, DECIMAL, SIGNED, UNSIGNED, TIME, JSON, etc.
func (p *Parser) parseCastDataType() (*nodes.DataType, error) {
	// Completion: offer CAST/CONVERT type candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	// SIGNED [INTEGER]
	if p.cur.Type == kwSIGNED {
		p.advance()
		p.match(kwINTEGER)
		p.match(kwINT)
		return &nodes.DataType{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "SIGNED"}, nil
	}

	// UNSIGNED [INTEGER]
	if p.cur.Type == kwUNSIGNED {
		p.advance()
		p.match(kwINTEGER)
		p.match(kwINT)
		return &nodes.DataType{Loc: nodes.Loc{Start: start, End: p.pos()}, Name: "UNSIGNED"}, nil
	}

	return p.parseDataType()
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
		charset, _, err := p.parseIdentifier()
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
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	conv.Loc.End = p.prev.End
	return conv, nil
}

// parseIsExpr parses IS [NOT] NULL / TRUE / FALSE / UNKNOWN.
func (p *Parser) parseIsExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
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
func (p *Parser) parseBetweenExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume BETWEEN

	low, err := p.parseExprPrec(precAdd) // parse at higher precedence to avoid AND confusion
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(kwAND); err != nil {
		return nil, err
	}

	high, err := p.parseExprPrec(precAdd)
	if err != nil {
		return nil, err
	}

	return &nodes.BetweenExpr{
		Loc:  nodes.Loc{Start: start, End: p.pos()},
		Not:  not,
		Expr: left,
		Low:  low,
		High: high,
	}, nil
}

// parseInExpr parses [NOT] IN (value_list) or [NOT] IN (subquery).
func (p *Parser) parseInExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
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
			first, err := p.parseInfixExprPrec(sub, precAssign)
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
func (p *Parser) parseLikeExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume LIKE

	pattern, err := p.parseExprPrec(precComparison + 1)
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
		esc, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		like.Escape = esc
	}

	like.Loc.End = p.prev.End
	return like, nil
}

// parseRegexpExpr parses [NOT] REGEXP pattern.
func (p *Parser) parseRegexpExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	not := false
	if _, ok := p.match(kwNOT); ok {
		not = true
	}
	p.advance() // consume REGEXP or RLIKE

	pattern, err := p.parseExprPrec(precComparison + 1)
	if err != nil {
		return nil, err
	}

	// Represent as BinaryExpr with BinOpRegexp
	expr := &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
		Op:    nodes.BinOpRegexp,
		Left:  left,
		Right: pattern,
	}

	if not {
		return &nodes.UnaryExpr{
			Loc:     nodes.Loc{Start: start, End: p.pos()},
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
			Loc:    nodes.Loc{Start: whenStart, End: p.pos()},
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
		Loc:    nodes.Loc{Start: start, End: p.pos()},
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
		Loc:   nodes.Loc{Start: start, End: p.pos()},
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
	against, err := p.parseExprPrec(precComparison + 1)
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
func (p *Parser) parseSoundsLikeExpr(left nodes.ExprNode) (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume SOUNDS

	if _, err := p.expect(kwLIKE); err != nil {
		return nil, err
	}

	right, err := p.parseExprPrec(precComparison + 1)
	if err != nil {
		return nil, err
	}

	return &nodes.BinaryExpr{
		Loc:   nodes.Loc{Start: start, End: p.pos()},
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
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		row.Items = append(row.Items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	row.Loc.End = p.prev.End
	return row, nil
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
		Loc:       nodes.Loc{Start: start, End: p.pos()},
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

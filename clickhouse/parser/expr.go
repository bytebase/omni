package parser

import (
	"fmt"
	"strings"
)

// Precedence levels (lowest to highest), matching ClickHouse's expression parser.
const (
	precNone    = 0
	precTernary = 1 // ? :
	precOr      = 2 // OR
	precAnd     = 3 // AND
	precNot     = 4 // NOT
	precCompare = 5 // =, ==, !=, <>, <, <=, >, >=, IS [NOT] NULL, [GLOBAL] [NOT] IN, [NOT] [I]LIKE, [NOT] BETWEEN
	precAdd     = 6 // +, -, ||
	precMul     = 7 // *, /, %
	precUnary   = 8 // unary -, unary +
	precPostfix = 9 // ::, [], .member
)

// parseColumnExprList parses a comma-separated list of column expressions
// (including `* [EXCEPT ...]` and subqueries), allowing optional aliases and
// ClickHouse's optional trailing comma before clause/list terminators.
func (p *parser) parseColumnExprList() ([]Expr, error) {
	var list []Expr
	for {
		expr, err := p.parseColumnsExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
		if !p.match(TokComma) {
			break
		}
		if p.isTrailingCommaTerminator() {
			break
		}
	}
	return list, nil
}

func (p *parser) isTrailingCommaTerminator() bool {
	switch p.cur().Type {
	case TokEOF, TokSemicolon, TokRParen, TokRBracket,
		KwFrom, KwWhere, KwPrewhere, KwGroup, KwHaving, KwQualify,
		KwOrder, KwLimit, KwSettings, KwFormat, KwInto, KwUnion,
		KwIntersect, KwExcept, KwWindow, KwArray, KwLeft, KwInner:
		return true
	default:
		return false
	}
}

// parseColumnsExpr parses one item of columnExprList (columnsExpr in ClickHouseParser.g4).
func (p *parser) parseColumnsExpr() (Expr, error) {
	// Bare `* [EXCEPT ...]`
	if p.cur().Type == TokStar {
		return p.parseStarExpr(nil, p.cur().Pos)
	}
	// `a.* [EXCEPT ...]` or `db.tbl.* [EXCEPT ...]`
	if parts, startPos, ok := p.tryMatchQualifiedStar(); ok {
		return p.parseStarExpr(parts, startPos)
	}
	return p.parseExprWithAlias()
}

func (p *parser) tryMatchQualifiedStar() ([]string, int, bool) {
	save := p.pos
	if !isIdentifier(p.cur()) {
		return nil, 0, false
	}
	startPos := p.cur().Pos
	first := p.advance().Val
	if !p.match(TokDot) {
		p.pos = save
		return nil, 0, false
	}
	if p.cur().Type == TokStar {
		return []string{first}, startPos, true
	}
	if isIdentifier(p.cur()) {
		second := p.advance().Val
		if p.match(TokDot) && p.cur().Type == TokStar {
			return []string{first, second}, startPos, true
		}
	}
	p.pos = save
	return nil, 0, false
}

func (p *parser) parseStarExpr(qualifier []string, startPos int) (Expr, error) {
	tok, err := p.expect(TokStar)
	if err != nil {
		return nil, err
	}
	if len(qualifier) == 0 {
		startPos = tok.Pos
	}
	star := &StarExpr{
		span:      span{Pos: startPos, End: tok.End},
		Qualifier: qualifier,
	}
	for {
		switch {
		case p.cur().Type == KwExcept && !p.isSetOpExceptAfterStar():
			trStart := p.advance().Pos // EXCEPT
			tr := &StarTransformer{span: span{Pos: trStart}, Kind: "EXCEPT"}
			if p.isWord("STRICT") {
				p.advance()
				tr.Strict = true
			}
			if p.match(TokLParen) {
				if isStringTok(p.cur().Type) {
					lit, err := p.parseStringLiteral()
					if err != nil {
						return nil, err
					}
					tr.Regex = lit.Val
					star.ExceptRe = lit.Val
				} else {
					for {
						id, err := p.parseIdentifier()
						if err != nil {
							return nil, err
						}
						tr.Columns = append(tr.Columns, id)
						star.Except = append(star.Except, id)
						if !p.match(TokComma) {
							break
						}
					}
				}
				rp, err := p.expect(TokRParen)
				if err != nil {
					return nil, err
				}
				tr.End = rp.End
			} else if isStringTok(p.cur().Type) {
				lit, err := p.parseStringLiteral()
				if err != nil {
					return nil, err
				}
				tr.Regex = lit.Val
				star.ExceptRe = lit.Val
				tr.End = lit.End
			} else {
				id, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				tr.Columns = []string{id}
				star.Except = append(star.Except, id)
				tr.End = p.peek(-1).End
			}
			star.Transformers = append(star.Transformers, tr)
			star.End = tr.End

		case p.cur().Type == KwReplace && (p.peek(1).Type == TokLParen || (p.isWordAt(1, "STRICT") && p.peek(2).Type == TokLParen)):
			trStart := p.advance().Pos // REPLACE
			tr := &StarTransformer{span: span{Pos: trStart}, Kind: "REPLACE"}
			if p.isWord("STRICT") {
				p.advance()
				tr.Strict = true
			}
			if _, err := p.expect(TokLParen); err != nil {
				return nil, err
			}
			for {
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(KwAs); err != nil {
					return nil, err
				}
				alias, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				tr.Replaces = append(tr.Replaces, &AliasExpr{
					span:  span{Pos: expr.StartPos(), End: p.peek(-1).End},
					Expr:  expr,
					Alias: alias,
				})
				if !p.match(TokComma) {
					break
				}
			}
			rp, err := p.expect(TokRParen)
			if err != nil {
				return nil, err
			}
			tr.End = rp.End
			star.Transformers = append(star.Transformers, tr)
			star.End = tr.End

		case p.isWord("APPLY") && (p.peek(1).Type == TokLParen || isBareAlias(p.peek(1))):
			trStart := p.advance().Pos // APPLY
			tr := &StarTransformer{span: span{Pos: trStart}, Kind: "APPLY"}
			if p.match(TokLParen) {
				if lam, ok, err := p.tryParseLambda(); err != nil {
					return nil, err
				} else if ok {
					tr.ApplyExpr = lam
				} else {
					fnStart := p.cur()
					fnName, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					if p.cur().Type == TokLParen {
						fc, err := p.parseFuncCall(fnStart)
						if err != nil {
							return nil, err
						}
						tr.ApplyExpr = fc
					} else {
						tr.ApplyFunc = fnName
					}
				}
				rp, err := p.expect(TokRParen)
				if err != nil {
					return nil, err
				}
				tr.End = rp.End
			} else {
				fnName, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				tr.ApplyFunc = fnName
				tr.End = p.peek(-1).End
			}
			star.Transformers = append(star.Transformers, tr)
			star.End = tr.End

		default:
			return star, nil
		}
	}
}

func (p *parser) isSetOpExceptAfterStar() bool {
	next := p.peek(1).Type
	if next == KwSelect || next == KwWith || next == KwFrom {
		return true
	}
	if next == KwAll || next == KwDistinct {
		n2 := p.peek(2).Type
		return n2 == KwSelect || n2 == KwWith || n2 == KwFrom || n2 == TokLParen
	}
	return false
}

// parseExprWithAlias parses an expression with an optional `AS alias` or bare alias.
func (p *parser) parseExprWithAlias() (Expr, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return p.parseOptionalAlias(expr)
}

// parseExprWithExplicitAlias parses an expression with an optional explicit
// `AS identifier` alias (without consuming bare implicit aliases).
func (p *parser) parseExprWithExplicitAlias() (Expr, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.match(KwAs) {
		alias, err := p.parseAnyWordIdent()
		if err != nil {
			return nil, err
		}
		return &AliasExpr{
			span:  span{Pos: expr.StartPos(), End: p.peek(-1).End},
			Expr:  expr,
			Alias: alias,
		}, nil
	}
	return expr, nil
}

func (p *parser) parseOptionalAlias(expr Expr) (Expr, error) {
	if p.match(KwAs) {
		alias, err := p.parseAnyWordIdent()
		if err != nil {
			return nil, err
		}
		return &AliasExpr{
			span:  span{Pos: expr.StartPos(), End: p.peek(-1).End},
			Expr:  expr,
			Alias: alias,
		}, nil
	}
	if isBareAlias(p.cur()) && !p.isFetchClauseStart() && !p.isTrailingCommentClause() {
		tok := p.advance()
		return &AliasExpr{
			span:  span{Pos: expr.StartPos(), End: tok.End},
			Expr:  expr,
			Alias: tok.Val,
		}, nil
	}
	return expr, nil
}

// parseExpr parses an expression without a trailing alias.
func (p *parser) parseExpr() (Expr, error) {
	return p.parseExprPrec(precNone)
}

// parseExprNoAlias is an explicit alias for parseExpr used in clauses where
// trailing keywords or clauses follow directly.
func (p *parser) parseExprPrec(minPrec int) (Expr, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	left, err := p.parseUnaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.cur()

		// Postfix: :: DataType, [idx], .member
		if precPostfix > minPrec {
			if tok.Type == TokDoubleColon {
				p.advance()
				dt, err := p.parseDataType()
				if err != nil {
					return nil, err
				}
				left = &CastExpr{
					span: span{Pos: left.StartPos(), End: dt.EndPos()},
					Expr: left,
					Type: dt,
				}
				continue
			}
			if tok.Type == TokLBracket {
				p.advance()
				idx, err := p.parseExprWithAlias()
				if err != nil {
					return nil, err
				}
				rb, err := p.expect(TokRBracket)
				if err != nil {
					return nil, err
				}
				left = &IndexExpr{
					span:  span{Pos: left.StartPos(), End: rb.End},
					Expr:  left,
					Index: idx,
				}
				continue
			}
			if tok.Type == TokDot {
				p.advance()
				next := p.cur()
				switch {
				case next.Type == TokInt:
					p.advance()
					left = &TupleAccessExpr{
						span:   span{Pos: left.StartPos(), End: next.End},
						Expr:   left,
						Member: next.Raw,
					}
				case next.Type == TokFloat && isSimpleDottedInts(next.Raw):
					// Handle nested tuple access like `t.1.2` where lexer produced `1.2`.
					p.advance()
					parts := strings.Split(next.Raw, ".")
					memberEnd := next.Pos
					for i, part := range parts {
						if i > 0 {
							memberEnd++ // account for '.'
						}
						memberEnd += len(part)
						left = &TupleAccessExpr{
							span:   span{Pos: left.StartPos(), End: memberEnd},
							Expr:   left,
							Member: part,
						}
					}
				case isStringTok(next.Type):
					lit, err := p.parseStringLiteral()
					if err != nil {
						return nil, err
					}
					left = &TupleAccessExpr{
						span:   span{Pos: left.StartPos(), End: lit.End},
						Expr:   left,
						Member: lit.Val,
					}
				case isIdentifier(next):
					p.advance()
					left = &TupleAccessExpr{
						span:   span{Pos: left.StartPos(), End: next.End},
						Expr:   left,
						Member: next.Val,
					}
				default:
					return nil, p.errAt(next, fmt.Sprintf("expected tuple element or identifier after '.', got %q", next.Raw))
				}
				continue
			}
		}

		// Multiplicative: *, /, %
		if precMul > minPrec && (tok.Type == TokStar || tok.Type == TokSlash || tok.Type == TokPercent) {
			op := p.advance().Raw
			right, err := p.parseExprPrec(precMul)
			if err != nil {
				return nil, err
			}
			left = &BinaryExpr{span: span{Pos: left.StartPos(), End: right.EndPos()}, Op: op, Left: left, Right: right}
			continue
		}

		// Additive & concat: +, -, ||
		if precAdd > minPrec && (tok.Type == TokPlus || tok.Type == TokMinus || tok.Type == TokConcat) {
			op := p.advance().Raw
			right, err := p.parseExprPrec(precAdd)
			if err != nil {
				return nil, err
			}
			left = &BinaryExpr{span: span{Pos: left.StartPos(), End: right.EndPos()}, Op: op, Left: left, Right: right}
			continue
		}

		// Comparisons, IS [NOT] NULL, [NOT] BETWEEN, [GLOBAL] [NOT] IN, [NOT] [I]LIKE
		if precCompare > minPrec {
			if tok.Type == TokEq || tok.Type == TokEqDouble || tok.Type == TokNotEq ||
				tok.Type == TokLT || tok.Type == TokLE || tok.Type == TokGT || tok.Type == TokGE {
				op := p.advance().Raw
				right, err := p.parseExprPrec(precCompare)
				if err != nil {
					return nil, err
				}
				left = &BinaryExpr{span: span{Pos: left.StartPos(), End: right.EndPos()}, Op: op, Left: left, Right: right}
				continue
			}
			if tok.Type == KwIs {
				if !p.isIsPredicateStart() {
					break
				}
				p.advance()
				not := p.match(KwNot)
				if p.cur().Type == KwNull {
					nt := p.advance()
					left = &IsNullExpr{span: span{Pos: left.StartPos(), End: nt.End}, Not: not, Expr: left}
					continue
				}
				if p.cur().Type == KwTrue || p.cur().Type == KwFalse {
					lit := p.advance()
					op := "IS"
					if not {
						op = "IS NOT"
					}
					left = &BinaryExpr{
						span:  span{Pos: left.StartPos(), End: lit.End},
						Op:    op,
						Left:  left,
						Right: &LiteralExpr{span: span{Pos: lit.Pos, End: lit.End}, Kind: "BOOL", Raw: lit.Raw, Val: strings.ToLower(lit.Val)},
					}
					continue
				}
				if p.cur().Type == KwDistinct && p.peek(1).Type == KwFrom {
					p.advance()
					p.advance()
					right, err := p.parseExprPrec(precCompare)
					if err != nil {
						return nil, err
					}
					op := "IS DISTINCT FROM"
					if not {
						op = "IS NOT DISTINCT FROM"
					}
					left = &BinaryExpr{span: span{Pos: left.StartPos(), End: right.EndPos()}, Op: op, Left: left, Right: right}
					continue
				}
				return nil, p.errAt(p.cur(), fmt.Sprintf("expected NULL after IS, got %q", p.cur().Raw))
			}
			if tok.Type == KwBetween || (tok.Type == KwNot && p.peek(1).Type == KwBetween) {
				not := p.match(KwNot)
				p.advance() // BETWEEN
				low, err := p.parseExprPrec(precCompare)
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(KwAnd); err != nil {
					return nil, err
				}
				high, err := p.parseExprPrec(precCompare)
				if err != nil {
					return nil, err
				}
				left = &BetweenExpr{
					span: span{Pos: left.StartPos(), End: high.EndPos()},
					Not:  not,
					Expr: left,
					Low:  low,
					High: high,
				}
				continue
			}
			if (tok.Type == KwIn && p.peek(1).Type != KwPartition) ||
				(tok.Type == KwNot && p.peek(1).Type == KwIn) ||
				(tok.Type == KwGlobal && (p.peek(1).Type == KwIn || (p.peek(1).Type == KwNot && p.peek(2).Type == KwIn))) {
				global := p.match(KwGlobal)
				not := p.match(KwNot)
				p.advance() // IN
				right, err := p.parseExprPrec(precCompare)
				if err != nil {
					return nil, err
				}
				left = &InExpr{
					span:   span{Pos: left.StartPos(), End: right.EndPos()},
					Global: global,
					Not:    not,
					Left:   left,
					Right:  right,
				}
				continue
			}
			if tok.Type == KwLike || tok.Type == KwIlike ||
				(tok.Type == KwNot && (p.peek(1).Type == KwLike || p.peek(1).Type == KwIlike)) {
				not := p.match(KwNot)
				ilike := p.advance().Type == KwIlike
				right, err := p.parseExprPrec(precCompare)
				if err != nil {
					return nil, err
				}
				left = &LikeExpr{
					span:  span{Pos: left.StartPos(), End: right.EndPos()},
					Not:   not,
					Ilike: ilike,
					Left:  left,
					Right: right,
				}
				continue
			}
		}

		// AND
		if precAnd > minPrec && tok.Type == KwAnd {
			p.advance()
			right, err := p.parseExprPrec(precAnd)
			if err != nil {
				return nil, err
			}
			left = &BinaryExpr{span: span{Pos: left.StartPos(), End: right.EndPos()}, Op: "AND", Left: left, Right: right}
			continue
		}

		// OR
		if precOr > minPrec && tok.Type == KwOr {
			p.advance()
			right, err := p.parseExprPrec(precOr)
			if err != nil {
				return nil, err
			}
			left = &BinaryExpr{span: span{Pos: left.StartPos(), End: right.EndPos()}, Op: "OR", Left: left, Right: right}
			continue
		}

		// Right-associative ternary: cond ? then : else
		if precTernary >= minPrec && tok.Type == TokQuestion && minPrec <= precTernary {
			p.advance()
			thenExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TokColon); err != nil {
				return nil, err
			}
			elseExpr, err := p.parseExprPrec(precTernary)
			if err != nil {
				return nil, err
			}
			left = &TernaryExpr{
				span: span{Pos: left.StartPos(), End: elseExpr.EndPos()},
				Cond: left,
				Then: thenExpr,
				Else: elseExpr,
			}
			continue
		}

		break
	}

	return left, nil
}

// isIsPredicateStart reports whether the current IS keyword begins a ClickHouse
// IS [NOT] {NULL|TRUE|FALSE|DISTINCT FROM} predicate. When it does not (e.g.
// `SELECT 1e1 IS`), the expression stops so the caller may consume IS as a bare
// alias, matching ClickHouse's expression parser.
func (p *parser) isIsPredicateStart() bool {
	next := p.peek(1)
	switch next.Type {
	case KwNull, KwTrue, KwFalse, KwNot:
		return true
	case KwDistinct:
		return p.peek(2).Type == KwFrom
	default:
		return false
	}
}

func isSimpleDottedInts(raw string) bool {
	if strings.ContainsAny(raw, "eEpPxXbB+-_") {
		return false
	}
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
		for i := 0; i < len(part); i++ {
			if !isDigit(part[i]) {
				return false
			}
		}
	}
	return true
}

func (p *parser) parseUnaryExpr() (Expr, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	tok := p.cur()
	switch tok.Type {
	case KwNot:
		p.advance()
		operand, err := p.parseExprPrec(precNot)
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{span: span{Pos: tok.Pos, End: operand.EndPos()}, Op: "NOT", Expr: operand}, nil
	case TokMinus:
		// If followed directly by a number literal (without unary recursion), parse signed literal or negate.
		p.advance()
		operand, err := p.parseExprPrec(precUnary)
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{span: span{Pos: tok.Pos, End: operand.EndPos()}, Op: "-", Expr: operand}, nil
	case TokPlus:
		// ClickHouse supports unary + before numeric literals (numberLiteral: (PLUS | DASH)? ...)
		if isNumberStart(p.peek(1), p.peek(2)) && p.peek(1).Type != TokPlus && p.peek(1).Type != TokMinus {
			return p.parseNumberLiteral()
		}
	}
	return p.parsePrimaryExpr()
}

func (p *parser) parsePrimaryExpr() (Expr, error) {
	tok := p.cur()

	switch tok.Type {
	case KwCase:
		if p.peek(1).Type == KwWhen {
			return p.parseCaseExpr()
		}
		if p.canStartCaseOperand(p.peek(1)) && !p.noCaseAt[p.pos] {
			save := p.pos
			if c, ok, err := p.tryParseCaseWithOperand(); err != nil {
				return nil, err
			} else if ok {
				return c, nil
			}
			if p.noCaseAt == nil {
				p.noCaseAt = make(map[int]bool)
			}
			p.noCaseAt[save] = true
			p.pos = save
		}
	case KwCast:
		if p.peek(1).Type == TokLParen {
			return p.parseCastExpr()
		}
	case KwDate:
		if isStringTok(p.peek(1).Type) {
			p.advance()
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			return &CastExpr{
				span: span{Pos: tok.Pos, End: lit.End},
				Expr: lit,
				Type: &DataType{span: span{Pos: tok.Pos, End: tok.End}, Name: "Date"},
			}, nil
		}
	case KwTimestamp:
		if isStringTok(p.peek(1).Type) {
			p.advance()
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			return &CastExpr{
				span: span{Pos: tok.Pos, End: lit.End},
				Expr: lit,
				Type: &DataType{span: span{Pos: tok.Pos, End: tok.End}, Name: "DateTime"},
			}, nil
		}
	case KwExtract:
		if p.peek(1).Type == TokLParen {
			return p.parseExtractExpr()
		}
	case KwInterval:
		// INTERVAL expr unit (only when not used as column name, e.g. followed by expr and unit)
		if p.peek(1).Type != TokEOF && p.peek(1).Type != TokComma && p.peek(1).Type != TokRParen && p.peek(1).Type != TokSemicolon {
			if expr, ok, err := p.tryParseIntervalExpr(); err != nil {
				return nil, err
			} else if ok {
				return expr, nil
			}
		}
	case KwSubstring:
		if p.peek(1).Type == TokLParen {
			return p.parseSubstringOrFuncCall()
		}
	case KwTrim:
		if p.peek(1).Type == TokLParen {
			if expr, ok, err := p.tryParseTrimExpr(); err != nil {
				return nil, err
			} else if ok {
				return expr, nil
			}
		}
	case TokInt, TokFloat, TokString, TokHexString, TokBinString, KwTrue, KwFalse, KwNull, KwInf, KwNan:
		return p.parseLiteral()
	case TokStar:
		return p.parseStarExpr(nil, tok.Pos)
	case TokLBracket:
		return p.parseArrayExpr()
	case TokLParen:
		return p.parseParenOrSubqueryOrTuple()
	case TokLBrace:
		return p.parseQueryParam()
	}

	// Identifier, qualified identifier, or function call
	if isIdentifier(tok) {
		return p.parseIdentOrFuncCall()
	}

	return nil, p.errAt(tok, fmt.Sprintf("unexpected token %q in expression", tok.Raw))
}

func (p *parser) canStartCaseOperand(tok Token) bool {
	switch tok.Type {
	case TokEOF, TokComma, TokRParen, TokRBracket, TokRBrace, TokSemicolon,
		TokEq, TokEqDouble, TokNotEq, TokLT, TokLE, TokGT, TokGE,
		TokStar, TokSlash, TokPercent, TokConcat, TokQuestion, TokColon,
		TokDoubleColon, TokDot,
		KwAs, KwFrom, KwWhere, KwPrewhere, KwGroup, KwHaving, KwQualify,
		KwOrder, KwLimit, KwSettings, KwFormat, KwInto, KwUnion,
		KwIntersect, KwExcept, KwWindow, KwThen, KwElse, KwEnd,
		KwAnd, KwOr, KwIs, KwIn, KwBetween, KwLike, KwIlike, KwGlobal,
		KwJoin, KwLeft, KwRight, KwInner, KwFull, KwCross, KwOn, KwUsing:
		return false
	default:
		return true
	}
}

func (p *parser) tryParseCaseWithOperand() (Expr, bool, error) {
	start := p.advance().Pos // CASE
	op, err := p.parseExprWithAlias()
	if err != nil {
		return nil, false, nil
	}
	if p.cur().Type != KwWhen {
		return nil, false, nil
	}
	c := &CaseExpr{span: span{Pos: start}, Operand: op}
	for p.match(KwWhen) {
		cond, err := p.parseExprWithAlias()
		if err != nil {
			return nil, true, err
		}
		if _, err := p.expect(KwThen); err != nil {
			return nil, true, err
		}
		then, err := p.parseExprWithAlias()
		if err != nil {
			return nil, true, err
		}
		c.Whens = append(c.Whens, CaseWhen{Cond: cond, Then: then})
	}
	if p.match(KwElse) {
		el, err := p.parseExprWithAlias()
		if err != nil {
			return nil, true, err
		}
		c.Else = el
	}
	endTok, err := p.expect(KwEnd)
	if err != nil {
		return nil, true, err
	}
	c.End = endTok.End
	return c, true, nil
}

func (p *parser) parseCaseExpr() (Expr, error) {
	start := p.advance().Pos // CASE
	c := &CaseExpr{span: span{Pos: start}}
	if p.cur().Type != KwWhen {
		op, err := p.parseExprWithAlias()
		if err != nil {
			return nil, err
		}
		c.Operand = op
	}
	if p.cur().Type != KwWhen {
		return nil, p.errAt(p.cur(), "expected WHEN in CASE expression")
	}
	for p.match(KwWhen) {
		cond, err := p.parseExprWithAlias()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwThen); err != nil {
			return nil, err
		}
		then, err := p.parseExprWithAlias()
		if err != nil {
			return nil, err
		}
		c.Whens = append(c.Whens, CaseWhen{Cond: cond, Then: then})
	}
	if p.match(KwElse) {
		el, err := p.parseExprWithAlias()
		if err != nil {
			return nil, err
		}
		c.Else = el
	}
	endTok, err := p.expect(KwEnd)
	if err != nil {
		return nil, err
	}
	c.End = endTok.End
	return c, nil
}

func (p *parser) parseCastExpr() (Expr, error) {
	start := p.advance().Pos // CAST
	p.advance()              // (
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.match(KwAs) {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		return &CastExpr{span: span{Pos: start, End: rp.End}, Expr: expr, Type: dt}, nil
	}
	// ClickHouse also supports CAST(expr, 'Type') as a regular 2-arg function call!
	if p.match(TokComma) {
		arg2, err := p.parseExprWithAlias()
		if err != nil {
			return nil, err
		}
		args := []Expr{expr, arg2}
		for p.match(TokComma) {
			arg, err := p.parseExprWithAlias()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		return &FuncCallExpr{span: span{Pos: start, End: rp.End}, Name: "CAST", Args: args}, nil
	}
	return nil, p.errAt(p.cur(), "expected AS or ',' in CAST expression")
}

func (p *parser) parseExtractExpr() (Expr, error) {
	start := p.advance().Pos // EXTRACT
	p.advance()              // (
	if !isIntervalUnit(p.cur().Type) {
		return nil, p.errAt(p.cur(), fmt.Sprintf("expected interval unit in EXTRACT, got %q", p.cur().Raw))
	}
	part := strings.ToUpper(p.advance().Val)
	if _, err := p.expect(KwFrom); err != nil {
		return nil, err
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	return &ExtractExpr{span: span{Pos: start, End: rp.End}, Part: part, From: from}, nil
}

func (p *parser) tryParseIntervalExpr() (Expr, bool, error) {
	if p.noIntervalAt[p.pos] {
		return nil, false, nil
	}
	save := p.pos
	start := p.advance().Pos // INTERVAL
	val, err := p.parseExpr()
	if err != nil || !isIntervalUnit(p.cur().Type) {
		if p.noIntervalAt == nil {
			p.noIntervalAt = make(map[int]bool)
		}
		p.noIntervalAt[save] = true
		p.pos = save
		return nil, false, nil
	}
	unitTok := p.advance()
	return &IntervalExpr{
		span:  span{Pos: start, End: unitTok.End},
		Value: val,
		Unit:  strings.ToUpper(unitTok.Val),
	}, true, nil
}

// parseSubstringOrFuncCall parses SUBSTRING(...) in a single pass without
// reparsing the first argument, handling both SQL-standard
// `SUBSTRING(expr FROM start [FOR len])` and function-call `SUBSTRING(a, b, c)`.
func (p *parser) parseSubstringOrFuncCall() (Expr, error) {
	nameTok := p.advance() // SUBSTRING
	p.advance()            // (

	if p.cur().Type == TokRParen {
		rp := p.advance()
		return p.parseFuncCallSuffix(&FuncCallExpr{
			span: span{Pos: nameTok.Pos, End: rp.End},
			Name: nameTok.Val,
		})
	}

	distinct := p.match(KwDistinct)
	first, err := p.parseFuncArg()
	if err != nil {
		return nil, err
	}

	// SQL-standard SUBSTRING(expr [alias] FROM start [alias] [FOR len [alias]])
	if !distinct {
		// Allow an optional bare alias before FROM as ClickHouse permits.
		if p.cur().Type != KwFrom && isBareAlias(p.cur()) && p.peek(1).Type == KwFrom {
			aliasTok := p.advance()
			first = &AliasExpr{
				span:  span{Pos: first.StartPos(), End: aliasTok.End},
				Expr:  first,
				Alias: aliasTok.Val,
			}
		}
		if p.match(KwFrom) {
			from, err := p.parseExprWithAlias()
			if err != nil {
				return nil, err
			}
			var forExpr Expr
			if p.match(KwFor) {
				forExpr, err = p.parseExprWithAlias()
				if err != nil {
					return nil, err
				}
			}
			rp, err := p.expect(TokRParen)
			if err != nil {
				return nil, err
			}
			return &SubstringExpr{
				span: span{Pos: nameTok.Pos, End: rp.End},
				Expr: first,
				From: from,
				For:  forExpr,
			}, nil
		}
	}

	// Comma-separated ordinary function call: reuse the already-parsed first argument.
	args := []Expr{first}
	for p.match(TokComma) {
		arg, err := p.parseFuncArg()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	fc := &FuncCallExpr{
		span:     span{Pos: nameTok.Pos, End: rp.End},
		Name:     nameTok.Val,
		Distinct: distinct,
		Args:     args,
	}
	return p.parseFuncCallSuffix(fc)
}

func (p *parser) parseFuncCallSuffix(fc *FuncCallExpr) (Expr, error) {
	if p.match(KwOver) {
		if p.cur().Type == TokLParen {
			spec, err := p.parseWindowSpec()
			if err != nil {
				return nil, err
			}
			fc.Over = spec
			fc.End = spec.End
		} else {
			winName, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			fc.OverName = winName
			fc.End = p.peek(-1).End
		}
	}
	return fc, nil
}

func (p *parser) tryParseTrimExpr() (Expr, bool, error) {
	save := p.pos
	start := p.advance().Pos // TRIM
	p.advance()              // (
	if p.cur().Type != KwBoth && p.cur().Type != KwLeading && p.cur().Type != KwTrailing {
		p.pos = save
		return nil, false, nil
	}
	mode := strings.ToUpper(p.advance().Val)
	c, err := p.parseStringLiteral()
	if err != nil {
		return nil, true, err
	}
	if _, err := p.expect(KwFrom); err != nil {
		return nil, true, err
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, true, err
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, true, err
	}
	return &TrimExpr{
		span:  span{Pos: start, End: rp.End},
		Mode:  mode,
		Chars: c,
		From:  from,
	}, true, nil
}

func (p *parser) parseArrayExpr() (Expr, error) {
	start := p.advance().Pos // [
	var items []Expr
	if p.cur().Type != TokRBracket {
		for {
			item, err := p.parseColumnsExpr()
			if err != nil {
				return nil, err
			}
			items = append(items, item)
			if !p.match(TokComma) {
				break
			}
		}
	}
	rb, err := p.expect(TokRBracket)
	if err != nil {
		return nil, err
	}
	return &ArrayExpr{span: span{Pos: start, End: rb.End}, Items: items}, nil
}

func (p *parser) parseParenOrSubqueryOrTuple() (Expr, error) {
	start := p.advance().Pos // (
	if p.cur().Type == TokRParen {
		rp := p.advance()
		return &TupleExpr{span: span{Pos: start, End: rp.End}}, nil
	}
	// Subquery: (SELECT ...) or (WITH ... SELECT ...)
	// Note: `(SELECT)` with no projection is the parenthesized identifier `SELECT`.
	if (p.cur().Type == KwSelect && p.peek(1).Type != TokRParen) || p.cur().Type == KwWith || p.cur().Type == KwFrom {
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		return &SubqueryExpr{span: span{Pos: start, End: rp.End}, Query: sub}, nil
	}
	// Parenthesized select union: ((SELECT ...))
	if p.cur().Type == TokLParen && p.looksLikeSelect() {
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		return &SubqueryExpr{span: span{Pos: start, End: rp.End}, Query: sub}, nil
	}

	first, err := p.parseColumnsExpr()
	if err != nil {
		return nil, err
	}
	if p.cur().Type == TokRParen {
		rp := p.advance()
		return &ParenExpr{span: span{Pos: start, End: rp.End}, Expr: first}, nil
	}
	items := []Expr{first}
	for p.match(TokComma) {
		item, err := p.parseColumnsExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	return &TupleExpr{span: span{Pos: start, End: rp.End}, Items: items}, nil
}

func (p *parser) parseQueryParam() (Expr, error) {
	return p.parseQueryParamNode()
}

func (p *parser) parseQueryParamNode() (*ParamExpr, error) {
	start := p.advance().Pos // {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TokColon); err != nil {
		return nil, err
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	rb, err := p.expect(TokRBrace)
	if err != nil {
		return nil, err
	}
	return &ParamExpr{span: span{Pos: start, End: rb.End}, Name: name, Type: dt}, nil
}

func (p *parser) parseIdentOrFuncCall() (Expr, error) {
	startTok := p.advance()
	// Function call: ident(...)
	if p.cur().Type == TokLParen {
		return p.parseFuncCall(startTok)
	}
	return &IdentExpr{
		span:  span{Pos: startTok.Pos, End: startTok.End},
		Parts: []string{startTok.Val},
	}, nil
}

func (p *parser) parseFuncCall(nameTok Token) (Expr, error) {
	fc := &FuncCallExpr{
		span: span{Pos: nameTok.Pos},
		Name: nameTok.Val,
	}
	p.advance() // (

	// Parse first argument list
	distinct := p.match(KwDistinct)
	var firstArgs []Expr
	if p.cur().Type != TokRParen {
		args, err := p.parseFuncArgList()
		if err != nil {
			return nil, err
		}
		firstArgs = args
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	fc.End = rp.End

	// Parametric aggregate function: f(params)(args) — params must be non-empty in ClickHouse.
	if !distinct && len(firstArgs) > 0 && p.cur().Type == TokLParen {
		fc.Params = firstArgs
		p.advance() // (
		fc.Distinct = p.match(KwDistinct)
		if p.cur().Type != TokRParen {
			args, err := p.parseFuncArgList()
			if err != nil {
				return nil, err
			}
			fc.Args = args
		}
		rp2, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		fc.End = rp2.End
	} else {
		fc.Distinct = distinct
		fc.Args = firstArgs
	}

	// Optional window OVER (...) or OVER window_name
	if p.match(KwOver) {
		if p.cur().Type == TokLParen {
			spec, err := p.parseWindowSpec()
			if err != nil {
				return nil, err
			}
			fc.Over = spec
			fc.End = spec.End
		} else {
			winName, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			fc.OverName = winName
			fc.End = p.peek(-1).End
		}
	}
	return fc, nil
}

func (p *parser) parseFuncArgList() ([]Expr, error) {
	var args []Expr
	for {
		arg, err := p.parseFuncArg()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if !p.match(TokComma) {
			break
		}
	}
	return args, nil
}

func (p *parser) parseFuncArg() (Expr, error) {
	// Lambda expression: `(x, y) -> expr`, `() -> expr`, `x -> expr`, or `x, y -> expr`
	if lam, ok, err := p.tryParseLambda(); err != nil {
		return nil, err
	} else if ok {
		return lam, nil
	}
	if p.cur().Type == TokStar {
		return p.parseStarExpr(nil, p.cur().Pos)
	}
	if parts, startPos, ok := p.tryMatchQualifiedStar(); ok {
		return p.parseStarExpr(parts, startPos)
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	// ClickHouse allows explicit `AS alias` inside function calls (e.g. `f(1 AS a)`),
	// but rejects bare implicit aliases (e.g. `f(1 a)`).
	if p.match(KwAs) {
		alias, err := p.parseAnyWordIdent()
		if err != nil {
			return nil, err
		}
		return &AliasExpr{
			span:  span{Pos: expr.StartPos(), End: p.peek(-1).End},
			Expr:  expr,
			Alias: alias,
		}, nil
	}
	return expr, nil
}

func (p *parser) tryParseLambda() (Expr, bool, error) {
	save := p.pos
	start := p.cur().Pos
	var params []string
	if p.match(TokLParen) {
		if p.cur().Type != TokRParen {
			if !isIdentifier(p.cur()) {
				p.pos = save
				return nil, false, nil
			}
			for {
				id, err := p.parseIdentifier()
				if err != nil {
					p.pos = save
					return nil, false, nil
				}
				params = append(params, id)
				if !p.match(TokComma) {
					break
				}
			}
		}
		if !p.match(TokRParen) || !p.match(TokArrow) {
			p.pos = save
			return nil, false, nil
		}
	} else if isIdentifier(p.cur()) {
		// Unparenthesized lambda: `x -> expr` or `x, y -> expr`.
		// If a previous lookahead already proved no `->` ends this comma-separated
		// identifier run, skip rescanning so wide argument lists stay O(N).
		if p.pos < p.noLambdaBefore {
			return nil, false, nil
		}
		for {
			id, err := p.parseIdentifier()
			if err != nil {
				if p.pos > p.noLambdaBefore {
					p.noLambdaBefore = p.pos
				}
				p.pos = save
				return nil, false, nil
			}
			params = append(params, id)
			if !p.match(TokComma) {
				break
			}
		}
		if !p.match(TokArrow) {
			if p.pos > p.noLambdaBefore {
				p.noLambdaBefore = p.pos
			}
			p.pos = save
			return nil, false, nil
		}
	} else {
		return nil, false, nil
	}
	body, err := p.parseExpr()
	if err != nil {
		return nil, true, err
	}
	return &LambdaExpr{
		span:   span{Pos: start, End: body.EndPos()},
		Params: params,
		Body:   body,
	}, true, nil
}

// parseDataType parses a ClickHouse columnTypeExpr:
//   - Simple: UInt64, String, DateTime
//   - Nested: Nested(a UInt64, b String)
//   - Enum: Enum8('a' = 1, 'b' = 2)
//   - Complex: Array(Nullable(UInt64)), Tuple(String, Int32), Tuple(a String, b Int32), Map(String, UInt64)
//   - Parametric: FixedString(16), Decimal(18, 4), DateTime64(3, 'UTC')
func (p *parser) parseDataType() (*DataType, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	start := p.cur().Pos
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	dt := &DataType{
		span: span{Pos: start, End: p.peek(-1).End},
		Name: name,
	}
	if !p.match(TokLParen) {
		return dt, nil
	}
	if p.cur().Type == TokRParen {
		rp := p.advance()
		dt.End = rp.End
		return dt, nil
	}

	// Enum values: stringLiteral '=' numberLiteral (',' ...)*
	if isStringTok(p.cur().Type) && p.peek(1).Type == TokEq {
		for {
			lit, err := p.parseStringLiteral()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TokEq); err != nil {
				return nil, err
			}
			num, err := p.parseNumberLiteral()
			if err != nil {
				return nil, err
			}
			dt.EnumValues = append(dt.EnumValues, EnumValue{Name: lit.Val, Value: num.Raw})
			if !p.match(TokComma) {
				break
			}
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		dt.End = rp.End
		return dt, nil
	}

	// Named elements (Nested / named Tuple): ident dataType (',' ident dataType)*
	// In ClickHouse, named elements in a type expression are only valid for Nested and Tuple.
	upperName := strings.ToUpper(name)
	if (upperName == "NESTED" || upperName == "TUPLE") && isIdentifier(p.cur()) && isIdentifier(p.peek(1)) {
		save := p.pos
		var cols []NestedCol
		ok := true
		for {
			colName, err := p.parseIdentifier()
			if err != nil {
				ok = false
				break
			}
			colType, err := p.parseDataType()
			if err != nil {
				ok = false
				break
			}
			cols = append(cols, NestedCol{Name: colName, Type: colType})
			if !p.match(TokComma) {
				break
			}
		}
		if ok && p.cur().Type == TokRParen {
			rp := p.advance()
			dt.NestedCols = cols
			dt.End = rp.End
			return dt, nil
		}
		p.pos = save
	}

	// General arguments: mix of data types and literal/expr params (no bare `*`, optional trailing comma)
	for {
		arg, typeArg, err := p.parseTypeArg()
		if err != nil {
			return nil, err
		}
		dt.Args = append(dt.Args, arg)
		if typeArg != nil {
			dt.TypeArgs = append(dt.TypeArgs, typeArg)
		}
		if !p.match(TokComma) {
			break
		}
		if p.cur().Type == TokRParen {
			break
		}
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	dt.End = rp.End
	return dt, nil
}

func (p *parser) parseTypeArg() (Expr, *DataType, error) {
	// Nested data type (e.g. Nullable(UInt64) or bare UInt64 inside Array(UInt64))
	if isIdentifier(p.cur()) && p.peek(1).Type == TokLParen && !p.noTypeArgAt[p.pos] {
		save := p.pos
		if dt, err := p.parseDataType(); err == nil && (p.cur().Type == TokComma || p.cur().Type == TokRParen) {
			return &TypeExpr{span: dt.span, Type: dt}, dt, nil
		}
		if p.noTypeArgAt == nil {
			p.noTypeArgAt = make(map[int]bool)
		}
		p.noTypeArgAt[save] = true
		p.pos = save
	}
	if isIdentifier(p.cur()) && (p.peek(1).Type == TokComma || p.peek(1).Type == TokRParen) {
		tok := p.advance()
		dt := &DataType{span: span{Pos: tok.Pos, End: tok.End}, Name: tok.Val}
		return &TypeExpr{span: dt.span, Type: dt}, dt, nil
	}
	if p.cur().Type == TokStar {
		return nil, nil, p.errAt(p.cur(), "unexpected '*' in data type argument")
	}
	expr, err := p.parseExpr()
	return expr, nil, err
}

package parser

import (
	"fmt"
	"strings"
)

func (p *parser) looksLikeSelect() bool {
	offset := 0
	for p.peek(offset).Type == TokLParen {
		offset++
		if offset > maxNestingDepth {
			return false
		}
	}
	t := p.peek(offset).Type
	return t == KwSelect || t == KwWith || t == KwFrom
}

func (p *parser) parseSelectUnionStmt() (*SelectUnionStmt, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	start := p.cur().Pos
	first, err := p.parseSelectStmtWithParens()
	if err != nil {
		return nil, err
	}
	u := &SelectUnionStmt{
		span:  span{Pos: start, End: first.EndPos()},
		Items: []*SelectStmt{first},
	}
	for p.cur().Type == KwUnion || p.cur().Type == KwIntersect || p.cur().Type == KwExcept {
		op := strings.ToUpper(p.advance().Val)
		if p.match(KwAll) {
			op += " ALL"
		} else if p.match(KwDistinct) {
			op += " DISTINCT"
		}
		next, err := p.parseSelectStmtWithParens()
		if err != nil {
			return nil, err
		}
		u.Ops = append(u.Ops, op)
		u.Items = append(u.Items, next)
		u.End = next.EndPos()
	}
	return u, nil
}

func (p *parser) parseSelectStmtWithParens() (*SelectStmt, error) {
	if p.cur().Type == TokLParen {
		start := p.advance().Pos
		inner, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		return &SelectStmt{
			span:   span{Pos: start, End: rp.End},
			Nested: inner,
		}, nil
	}
	return p.parseSelectStmt()
}

func (p *parser) parseSelectStmt() (*SelectStmt, error) {
	start := p.cur().Pos
	stmt := &SelectStmt{span: span{Pos: start}}

	// WITH clause
	if p.match(KwWith) {
		withs, err := p.parseWithExprList()
		if err != nil {
			return nil, err
		}
		stmt.With = withs
	}

	// ClickHouse supports FROM-first queries (`FROM tbl [WHERE ...] [SELECT ...]`).
	if p.cur().Type == KwFrom {
		p.advance()
		from, err := p.parseJoinExpr()
		if err != nil {
			return nil, err
		}
		stmt.From = from
		if p.match(KwSelect) {
			if err := p.parseSelectProjection(stmt); err != nil {
				return nil, err
			}
		}
	} else {
		if _, err := p.expect(KwSelect); err != nil {
			return nil, err
		}
		if err := p.parseSelectProjection(stmt); err != nil {
			return nil, err
		}
		// FROM clause
		if p.match(KwFrom) {
			from, err := p.parseJoinExpr()
			if err != nil {
				return nil, err
			}
			stmt.From = from
		}
	}

	// [LEFT | INNER] ARRAY JOIN (requires a FROM clause in ClickHouse)
	if stmt.From != nil && p.isArrayJoinStart() {
		aj, err := p.parseArrayJoinClause()
		if err != nil {
			return nil, err
		}
		stmt.ArrayJoin = aj
	}

	// WINDOW clause (may appear before or after WHERE/GROUP BY/HAVING in ClickHouse)
	if p.cur().Type == KwWindow {
		wins, err := p.parseWindowClause()
		if err != nil {
			return nil, err
		}
		stmt.Windows = append(stmt.Windows, wins...)
	}

	// PREWHERE
	if p.match(KwPrewhere) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Prewhere = expr
	}

	// WHERE
	if p.match(KwWhere) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Where = expr
	}

	// GROUP BY ... [WITH CUBE | WITH ROLLUP] [WITH TOTALS]
	if p.cur().Type == KwGroup && p.peek(1).Type == KwBy {
		gb, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = gb
	}
	if p.cur().Type == KwWith && (p.peek(1).Type == KwCube || p.peek(1).Type == KwRollup) {
		p.advance()
		mod := p.advance().Type
		if stmt.GroupBy == nil {
			stmt.GroupBy = &GroupByClause{span: span{Pos: p.peek(-2).Pos, End: p.peek(-1).End}}
		}
		if mod == KwCube {
			stmt.GroupBy.WithCube = true
		} else {
			stmt.GroupBy.WithRollup = true
		}
	}
	if p.cur().Type == KwWith && p.peek(1).Type == KwTotals {
		p.advance()
		p.advance()
		if stmt.GroupBy == nil {
			stmt.GroupBy = &GroupByClause{span: span{Pos: p.peek(-2).Pos, End: p.peek(-1).End}}
		}
		stmt.GroupBy.WithTotals = true
	}

	// HAVING
	if p.match(KwHaving) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Having = expr
	}

	// WINDOW (standard SQL placement after HAVING)
	if p.cur().Type == KwWindow {
		wins, err := p.parseWindowClause()
		if err != nil {
			return nil, err
		}
		stmt.Windows = append(stmt.Windows, wins...)
	}

	// QUALIFY
	if p.match(KwQualify) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Qualify = expr
	}

	// WINDOW may also follow QUALIFY
	if p.cur().Type == KwWindow {
		wins, err := p.parseWindowClause()
		if err != nil {
			return nil, err
		}
		stmt.Windows = append(stmt.Windows, wins...)
	}

	// ORDER BY
	if p.cur().Type == KwOrder && p.peek(1).Type == KwBy {
		p.advance()
		p.advance()
		items, err := p.parseOrderExprList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = items
	}

	// INTERPOLATE (only valid after ORDER BY in ClickHouse):
	// INTERPOLATE ( col [AS expr] [, ...] )
	if len(stmt.OrderBy) > 0 && p.match(KwInterpolate) {
		if p.match(TokLParen) {
			if p.cur().Type != TokRParen {
				for {
					col, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					if p.match(KwAs) {
						val, err := p.parseExpr()
						if err != nil {
							return nil, err
						}
						col = &BinaryExpr{span: span{Pos: col.StartPos(), End: val.EndPos()}, Op: "AS", Left: col, Right: val}
					}
					stmt.Interpolate = append(stmt.Interpolate, col)
					if !p.match(TokComma) {
						break
					}
				}
			}
			if _, err := p.expect(TokRParen); err != nil {
				return nil, err
			}
		}
	}

	// LIMIT ... BY ... and/or LIMIT ... [WITH TIES]
	if p.cur().Type == KwLimit {
		limStart := p.advance().Pos
		limExpr, offExpr, err := p.parseLimitExpr()
		if err != nil {
			return nil, err
		}
		if p.match(KwBy) {
			byExprs, err := p.parseColumnExprList()
			if err != nil {
				return nil, err
			}
			stmt.LimitBy = &LimitByClause{
				span:   span{Pos: limStart, End: p.peek(-1).End},
				Limit:  limExpr,
				Offset: offExpr,
				By:     byExprs,
			}
			// Optional second LIMIT clause after LIMIT ... BY ...
			if p.match(KwLimit) {
				lim2Start := p.peek(-1).Pos
				lim2, off2, err := p.parseLimitExpr()
				if err != nil {
					return nil, err
				}
				lc := &LimitClause{span: span{Pos: lim2Start, End: p.peek(-1).End}, Limit: lim2, Offset: off2}
				if p.cur().Type == KwWith && p.peek(1).Type == KwTies {
					p.advance()
					ties := p.advance()
					lc.WithTies = true
					lc.End = ties.End
				}
				stmt.Limit = lc
			}
		} else {
			lc := &LimitClause{span: span{Pos: limStart, End: p.peek(-1).End}, Limit: limExpr, Offset: offExpr}
			if p.cur().Type == KwWith && p.peek(1).Type == KwTies {
				p.advance()
				ties := p.advance()
				lc.WithTies = true
				lc.End = ties.End
			}
			stmt.Limit = lc
		}
	}

	// ClickHouse enforces inside ParserSelectQuery that WITH TIES requires ORDER BY.
	if (stmt.Top != nil && stmt.Top.WithTies) || (stmt.Limit != nil && stmt.Limit.WithTies) {
		if len(stmt.OrderBy) == 0 {
			return nil, p.errAt(p.cur(), "cannot use WITH TIES without ORDER BY")
		}
	}

	// SETTINGS
	if p.cur().Type == KwSettings {
		settings, err := p.parseSettingsClause()
		if err != nil {
			return nil, err
		}
		stmt.Settings = settings
	}

	stmt.End = p.peek(-1).End
	return stmt, nil
}

func (p *parser) parseSelectProjection(stmt *SelectStmt) error {
	// DISTINCT [ON (...)]
	if p.match(KwDistinct) {
		stmt.Distinct = true
		if p.cur().Type == KwOn && p.peek(1).Type == TokLParen {
			p.advance() // ON
			p.advance() // (
			exprs, err := p.parseColumnExprList()
			if err != nil {
				return err
			}
			if _, err := p.expect(TokRParen); err != nil {
				return err
			}
			stmt.DistinctOn = exprs
		}
	}

	// TOP n [WITH TIES]
	if p.match(KwTop) {
		topStart := p.peek(-1).Pos
		numTok, err := p.expect(TokInt)
		if err != nil {
			return err
		}
		top := &TopClause{span: span{Pos: topStart, End: numTok.End}, Count: numTok.Raw}
		if p.cur().Type == KwWith && p.peek(1).Type == KwTies {
			p.advance()
			ties := p.advance()
			top.WithTies = true
			top.End = ties.End
		}
		stmt.Top = top
	}

	// Projection list
	cols, err := p.parseColumnExprList()
	if err != nil {
		return err
	}
	stmt.Columns = cols
	return nil
}

func (p *parser) parseWithExprList() ([]*WithExpr, error) {
	var list []*WithExpr
	for {
		w, err := p.parseWithExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, w)
		if !p.match(TokComma) {
			break
		}
	}
	return list, nil
}

func (p *parser) parseWithExpr() (*WithExpr, error) {
	start := p.cur().Pos
	// Subquery CTE: [RECURSIVE] ident AS ( selectUnionStmt )
	if p.cur().Type == KwRecursive || (isIdentifier(p.cur()) && p.peek(1).Type == KwAs && p.peek(2).Type == TokLParen) {
		save := p.pos
		rec := p.match(KwRecursive)
		if isIdentifier(p.cur()) && p.peek(1).Type == KwAs && p.peek(2).Type == TokLParen {
			name, _ := p.parseIdentifier()
			p.advance() // AS
			p.advance() // (
			if p.looksLikeSelect() {
				sub, err := p.parseSelectUnionStmt()
				if err != nil {
					return nil, err
				}
				rp, err := p.expect(TokRParen)
				if err != nil {
					return nil, err
				}
				return &WithExpr{
					span:      span{Pos: start, End: rp.End},
					Recursive: rec,
					Alias:     name,
					Subquery:  sub,
				}, nil
			}
		}
		if rec {
			return nil, p.errAt(p.cur(), "expected recursive CTE definition")
		}
		p.pos = save
	}

	// Expression CTE: columnExpr AS identifier
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
	return &WithExpr{
		span:  span{Pos: start, End: p.peek(-1).End},
		Alias: alias,
		Expr:  expr,
	}, nil
}

func (p *parser) isArrayJoinStart() bool {
	t := p.cur().Type
	if t == KwArray && p.peek(1).Type == KwJoin {
		return true
	}
	if (t == KwLeft || t == KwInner) && p.peek(1).Type == KwArray && p.peek(2).Type == KwJoin {
		return true
	}
	return false
}

func (p *parser) parseArrayJoinClause() (*ArrayJoinClause, error) {
	start := p.cur().Pos
	aj := &ArrayJoinClause{span: span{Pos: start}}
	if p.match(KwLeft) {
		aj.Left = true
	} else if p.match(KwInner) {
		aj.Inner = true
	}
	if _, err := p.expect(KwArray); err != nil {
		return nil, err
	}
	if _, err := p.expect(KwJoin); err != nil {
		return nil, err
	}
	exprs, err := p.parseColumnExprList()
	if err != nil {
		return nil, err
	}
	aj.Exprs = exprs
	aj.End = p.peek(-1).End
	return aj, nil
}

func (p *parser) parseWindowClause() ([]*WindowDef, error) {
	if _, err := p.expect(KwWindow); err != nil {
		return nil, err
	}
	var wins []*WindowDef
	for {
		start := p.cur().Pos
		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(KwAs); err != nil {
			return nil, err
		}
		if _, err := p.expect(TokLParen); err != nil {
			return nil, err
		}
		spec, err := p.parseWindowSpec()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		wins = append(wins, &WindowDef{
			span: span{Pos: start, End: rp.End},
			Name: name,
			Spec: spec,
		})
		if !p.match(TokComma) {
			break
		}
	}
	return wins, nil
}

func (p *parser) parseWindowSpec() (*WindowSpec, error) {
	spec := &WindowSpec{span: span{Pos: p.cur().Pos}}
	// Optional base window name: OVER (existing_window_name [PARTITION BY ...] ...)
	if isIdentifier(p.cur()) &&
		!(p.cur().Type == KwPartition && p.peek(1).Type == KwBy) &&
		!(p.cur().Type == KwOrder && p.peek(1).Type == KwBy) &&
		p.cur().Type != KwRows && p.cur().Type != KwRange {
		p.advance()
	}
	if p.cur().Type == KwPartition && p.peek(1).Type == KwBy {
		p.advance()
		p.advance()
		// Parse partition keys without consuming ROWS/RANGE as an implicit alias.
		var exprs []Expr
		for {
			var expr Expr
			var err error
			if p.cur().Type == TokStar {
				expr, err = p.parseStarExpr(nil, p.cur().Pos)
			} else if parts, startPos, ok := p.tryMatchQualifiedStar(); ok {
				expr, err = p.parseStarExpr(parts, startPos)
			} else {
				expr, err = p.parseExpr()
			}
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, expr)
			if !p.match(TokComma) {
				break
			}
		}
		spec.PartitionBy = exprs
	}
	if p.cur().Type == KwOrder && p.peek(1).Type == KwBy {
		p.advance()
		p.advance()
		items, err := p.parseOrderExprList()
		if err != nil {
			return nil, err
		}
		spec.OrderBy = items
	}
	if p.cur().Type == KwRows || p.cur().Type == KwRange {
		spec.FrameUnits = strings.ToUpper(p.advance().Val)
		if p.match(KwBetween) {
			startBound, err := p.parseWindowFrameBound(true)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(KwAnd); err != nil {
				return nil, err
			}
			endBound, err := p.parseWindowFrameBound(false)
			if err != nil {
				return nil, err
			}
			spec.FrameStart = startBound
			spec.FrameEnd = endBound
		} else {
			startBound, err := p.parseWindowFrameBound(true)
			if err != nil {
				return nil, err
			}
			spec.FrameStart = startBound
		}
	}
	spec.End = p.peek(-1).End
	return spec, nil
}

func (p *parser) parseWindowFrameBound(isStart bool) (string, error) {
	startTok := p.cur()
	if p.cur().Type == KwCurrent && p.peek(1).Type == KwRow {
		p.advance()
		p.advance()
		return "CURRENT ROW", nil
	}
	if p.match(KwUnbounded) {
		if p.match(KwPreceding) {
			if !isStart {
				return "", p.errAt(startTok, "frame end cannot be UNBOUNDED PRECEDING")
			}
			return "UNBOUNDED PRECEDING", nil
		}
		if p.match(KwFollowing) {
			if isStart {
				return "", p.errAt(startTok, "frame start cannot be UNBOUNDED FOLLOWING")
			}
			return "UNBOUNDED FOLLOWING", nil
		}
		return "", p.errAt(p.cur(), "expected PRECEDING or FOLLOWING after UNBOUNDED")
	}
	num, err := p.parseNumberLiteral()
	if err != nil {
		return "", err
	}
	if p.match(KwPreceding) {
		return num.Raw + " PRECEDING", nil
	}
	if p.match(KwFollowing) {
		return num.Raw + " FOLLOWING", nil
	}
	return "", p.errAt(p.cur(), "expected PRECEDING or FOLLOWING in window frame bound")
}

func (p *parser) parseGroupByClause() (*GroupByClause, error) {
	start := p.advance().Pos // GROUP
	p.advance()              // BY
	gb := &GroupByClause{span: span{Pos: start}}

	if p.cur().Type == KwAll && p.peek(1).Type != TokLParen {
		tok := p.advance()
		gb.All = true
		gb.End = tok.End
		return gb, nil
	}
	if (p.cur().Type == KwCube || p.cur().Type == KwRollup) && p.peek(1).Type == TokLParen {
		gb.Modifier = strings.ToUpper(p.advance().Val)
		p.advance() // (
		exprs, err := p.parseColumnExprList()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		gb.Exprs = exprs
		gb.End = rp.End
		return gb, nil
	}
	if p.cur().Type == KwGrouping && p.peek(1).Type == KwSets && p.peek(2).Type == TokLParen {
		p.advance()
		p.advance()
		p.advance() // (
		gb.Modifier = "GROUPING SETS"
		exprs, err := p.parseColumnExprList()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		gb.Exprs = exprs
		gb.End = rp.End
		return gb, nil
	}
	exprs, err := p.parseColumnExprList()
	if err != nil {
		return nil, err
	}
	gb.Exprs = exprs
	gb.End = p.peek(-1).End
	return gb, nil
}

func (p *parser) parseOrderExprList() ([]*OrderByItem, error) {
	var items []*OrderByItem
	for {
		item, err := p.parseOrderExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if !p.match(TokComma) {
			break
		}
	}
	return items, nil
}

func (p *parser) parseOrderExpr() (*OrderByItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	item := &OrderByItem{span: span{Pos: expr.StartPos(), End: expr.EndPos()}, Expr: expr}
	if p.match(KwAsc) {
		item.Desc = false
	} else if p.match(KwDesc) || p.match(KwDescending) {
		item.Desc = true
	}
	if p.match(KwNulls) {
		if p.match(KwFirst) {
			item.NullsPos = "FIRST"
		} else if p.match(KwLast) {
			item.NullsPos = "LAST"
		} else {
			return nil, p.errAt(p.cur(), "expected FIRST or LAST after NULLS")
		}
	}
	if p.match(KwCollate) {
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		item.Collate = lit.Val
	}
	if p.cur().Type == KwWith && p.peek(1).Type == KwFill {
		p.advance()
		p.advance()
		item.WithFill = true
		if p.match(KwFrom) {
			fromExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			item.FillFrom = fromExpr
		}
		if p.match(KwTo) {
			toExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			item.FillTo = toExpr
		}
		if p.match(KwStep) {
			stepExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			item.FillStep = stepExpr
		}
	}
	item.End = p.peek(-1).End
	return item, nil
}

// parseLimitExpr parses `count`, `count OFFSET offset`, or `offset, count`
// and always returns (limit, offset, err).
func (p *parser) parseLimitExpr() (Expr, Expr, error) {
	first, err := p.parseExpr()
	if err != nil {
		return nil, nil, err
	}
	if p.match(TokComma) {
		count, err := p.parseExpr()
		if err != nil {
			return nil, nil, err
		}
		return count, first, nil
	}
	if p.match(KwOffset) {
		offset, err := p.parseExpr()
		if err != nil {
			return nil, nil, err
		}
		return first, offset, nil
	}
	return first, nil, nil
}

// --- FROM / JOIN parsing ---

func (p *parser) parseJoinExpr() (*JoinExpr, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	left, err := p.parseJoinPrimary()
	if err != nil {
		return nil, err
	}
	for {
		// Comma cross join
		if p.match(TokComma) {
			right, err := p.parseJoinPrimary()
			if err != nil {
				return nil, err
			}
			left = &JoinExpr{
				span:  span{Pos: left.Pos, End: right.End},
				Left:  left,
				Right: right,
				Kind:  "COMMA",
			}
			continue
		}
		// Explicit JOIN
		if ok, global, local, kind, strictness, isCross := p.tryParseJoinOp(); ok {
			right, err := p.parseJoinPrimary()
			if err != nil {
				return nil, err
			}
			je := &JoinExpr{
				span:       span{Pos: left.Pos, End: right.End},
				Left:       left,
				Right:      right,
				Global:     global,
				Local:      local,
				Kind:       kind,
				Strictness: strictness,
			}
			if !isCross {
				if p.match(KwOn) {
					onExprs, err := p.parseColumnExprList()
					if err != nil {
						return nil, err
					}
					je.On = onExprs
					je.End = p.peek(-1).End
				} else if p.match(KwUsing) {
					if p.match(TokLParen) {
						usingExprs, err := p.parseColumnExprList()
						if err != nil {
							return nil, err
						}
						rp, err := p.expect(TokRParen)
						if err != nil {
							return nil, err
						}
						je.Using = usingExprs
						je.End = rp.End
					} else {
						usingExprs, err := p.parseColumnExprList()
						if err != nil {
							return nil, err
						}
						je.Using = usingExprs
						je.End = p.peek(-1).End
					}
				} else {
					return nil, p.errAt(p.cur(), "expected ON or USING after JOIN")
				}
			}
			left = je
			continue
		}
		break
	}
	return left, nil
}

func (p *parser) tryParseJoinOp() (ok, global, local bool, kind, strictness string, isCross bool) {
	save := p.pos
	if p.match(KwGlobal) {
		global = true
	} else if p.match(KwLocal) {
		local = true
	}

	if p.cur().Type == KwCross && p.peek(1).Type == KwJoin {
		p.advance()
		p.advance()
		return true, global, local, "CROSS", "", true
	}

	// Parse join kind + strictness in either order (ClickHouse permits both
	// `ALL LEFT JOIN` and `LEFT ALL JOIN`, `LEFT OUTER ALL JOIN`, etc.)
	readStrictness := func() bool {
		switch p.cur().Type {
		case KwAll, KwAny, KwAsof, KwSemi, KwAnti:
			if strictness != "" {
				return false
			}
			strictness = strings.ToUpper(p.advance().Val)
			return true
		default:
			return false
		}
	}
	readKind := func() bool {
		switch p.cur().Type {
		case KwInner, KwLeft, KwRight, KwFull:
			if kind != "" {
				return false
			}
			kind = strings.ToUpper(p.advance().Val)
			if (kind == "LEFT" || kind == "RIGHT" || kind == "FULL") && p.cur().Type == KwOuter {
				p.advance()
			}
			return true
		default:
			return false
		}
	}

	readStrictness()
	readKind()
	readStrictness()

	if p.match(KwJoin) {
		if kind == "" {
			kind = "INNER"
		}
		// Validate strictness compatibility with kind:
		// SEMI / ANTI are only valid for LEFT / RIGHT joins; ASOF is not valid for FULL.
		if (strictness == "SEMI" || strictness == "ANTI") && kind != "LEFT" && kind != "RIGHT" {
			p.pos = save
			return false, false, false, "", "", false
		}
		if strictness == "ASOF" && kind == "FULL" {
			p.pos = save
			return false, false, false, "", "", false
		}
		return true, global, local, kind, strictness, false
	}

	p.pos = save
	return false, false, false, "", "", false
}

func (p *parser) parseJoinPrimary() (*JoinExpr, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	if p.cur().Type == TokLParen && !p.looksLikeSelect() {
		start := p.advance().Pos
		inner, err := p.parseJoinExpr()
		if err != nil {
			return nil, err
		}
		// ClickHouse rejects parenthesized single table names like `FROM (a)`:
		// parenthesized FROM items must be either a subquery or a JOIN expression.
		if inner.Left == nil {
			return nil, p.errAt(p.cur(), "expected JOIN inside parenthesized FROM expression")
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		inner.Pos = start
		inner.End = rp.End
		return inner, nil
	}
	te, err := p.parseTableExpr()
	if err != nil {
		return nil, err
	}
	je := &JoinExpr{
		span:  te.span,
		Table: te,
	}
	if p.match(KwFinal) {
		je.Final = true
		je.End = p.peek(-1).End
	}
	if p.match(KwSample) {
		r, err := p.parseRatioExpr()
		if err != nil {
			return nil, err
		}
		je.Sample = r
		if p.match(KwOffset) {
			off, err := p.parseRatioExpr()
			if err != nil {
				return nil, err
			}
			je.SampleOff = off
		}
		je.End = p.peek(-1).End
	}
	return je, nil
}

func (p *parser) parseRatioExpr() (string, error) {
	num, err := p.parseNumberLiteral()
	if err != nil {
		return "", err
	}
	if p.match(TokSlash) {
		den, err := p.parseNumberLiteral()
		if err != nil {
			return "", err
		}
		return num.Raw + "/" + den.Raw, nil
	}
	return num.Raw, nil
}

func (p *parser) parseTableExpr() (*TableExpr, error) {
	start := p.cur().Pos
	te := &TableExpr{span: span{Pos: start}}
	if p.cur().Type == TokLParen {
		p.advance()
		sub, err := p.parseSelectUnionStmt()
		if err != nil {
			return nil, err
		}
		rp, err := p.expect(TokRParen)
		if err != nil {
			return nil, err
		}
		te.Subquery = sub
		te.End = rp.End
	} else if isIdentifier(p.cur()) && p.peek(1).Type == TokLParen {
		fn, err := p.parseTableFunctionExpr()
		if err != nil {
			return nil, err
		}
		te.Func = fn
		te.End = fn.End
	} else {
		db, tbl, err := p.parseTableIdentifier()
		if err != nil {
			return nil, err
		}
		te.Database = db
		te.Table = tbl
		te.End = p.peek(-1).End
	}

	// Optional alias: AS ident | bare alias (not FINAL / SAMPLE / JOIN keywords)
	if p.match(KwAs) {
		alias, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		te.Alias = alias
		te.End = p.peek(-1).End
	} else if isBareAlias(p.cur()) {
		tok := p.advance()
		te.Alias = tok.Val
		te.End = tok.End
	}
	return te, nil
}

func (p *parser) parseTableFunctionExpr() (*FuncCallExpr, error) {
	if err := p.enter(); err != nil {
		return nil, err
	}
	defer p.leave()

	startTok := p.advance()
	if _, err := p.expect(TokLParen); err != nil {
		return nil, err
	}
	fc := &FuncCallExpr{span: span{Pos: startTok.Pos}, Name: startTok.Val}
	if p.cur().Type != TokRParen {
		for {
			arg, err := p.parseTableArgExpr()
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, arg)
			if !p.match(TokComma) {
				break
			}
		}
	}
	rp, err := p.expect(TokRParen)
	if err != nil {
		return nil, err
	}
	fc.End = rp.End
	return fc, nil
}

func (p *parser) parseTableArgExpr() (Expr, error) {
	if isIdentifier(p.cur()) && p.peek(1).Type == TokLParen {
		return p.parseTableFunctionExpr()
	}
	if isLiteralStart(p.cur(), p.peek(1)) {
		return p.parseLiteral()
	}
	// nestedIdentifier or general expression
	return p.parseExpr()
}

// Helper for error formatting in select.go
var _ = fmt.Sprintf

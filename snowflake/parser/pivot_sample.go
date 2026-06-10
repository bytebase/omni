package parser

// PIVOT / UNPIVOT and SAMPLE / TABLESAMPLE table clauses (T5.3).
//
//	PIVOT ( <agg>(<col>) [ [AS] alias ]
//	        FOR <col> IN ( <values> | ANY [ORDER BY …] | <subquery> )
//	        [ DEFAULT ON NULL (<expr>) ] ) [ [AS] alias ]
//
//	UNPIVOT [ {INCLUDE|EXCLUDE} NULLS ]
//	  ( <value_col> FOR <name_col> IN ( <col> [AS alias], … ) ) [alias]
//
//	{ SAMPLE | TABLESAMPLE } [ {BERNOULLI|ROW|SYSTEM|BLOCK} ]
//	  ( <prob> | <n> ROWS ) [ {SEED|REPEATABLE} (<n>) ]
//
// Triangulation: the legacy SnowflakeParser.g4 restricts the PIVOT aggregate
// to `id_(id_)` and the IN list to bare literals/ANY/subquery, but the
// Snowflake docs and the official corpus show an aliased aggregate
// (`SUM(amount) AS total`) and aliased IN values (`'2023_Q1' AS q1`). The docs
// win.

import "github.com/bytebase/omni/snowflake/ast"

// parsePivotTrailingAlias parses the optional trailing alias of a PIVOT /
// UNPIVOT clause. Unlike parseOptionalAlias it refuses to treat the head of a
// chained PIVOT/UNPIVOT clause — the (non-reserved) keyword followed by
// '(' — as an implicit alias, so `t PIVOT(...) PIVOT(...)` stays parseable as
// a chain instead of the second clause being eaten as alias "PIVOT" and its
// body discarded. An explicit `AS pivot` alias is still accepted.
func (p *Parser) parsePivotTrailingAlias() (ast.Ident, bool) {
	if (p.cur.Type == kwPIVOT || p.cur.Type == kwUNPIVOT) && p.peekNext().Type == '(' {
		return ast.Ident{}, false
	}
	return p.parseOptionalAlias()
}

// parsePivotClause parses a PIVOT (...) clause. The caller has verified that
// p.cur is kwPIVOT.
func (p *Parser) parsePivotClause() (*ast.PivotClause, error) {
	pivotTok := p.advance() // consume PIVOT
	clause := &ast.PivotClause{
		Loc: ast.Loc{Start: pivotTok.Loc.Start},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Aggregate function: <agg>(<col>) [ [AS] alias ]
	aggExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	agg, ok := aggExpr.(*ast.FuncCallExpr)
	if !ok {
		return nil, &ParseError{
			Loc: ast.NodeLoc(aggExpr),
			Msg: "expected an aggregate function call in PIVOT",
		}
	}
	clause.Agg = agg

	// Optional aggregate alias: AS total | total
	if alias, has := p.parseOptionalAlias(); has {
		clause.AggAlias = alias
	}

	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// FOR pivot column. A column reference (possibly qualified).
	forCol, err := p.parsePivotColumnRef()
	if err != nil {
		return nil, err
	}
	clause.ForColumn = forCol

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	in, err := p.parsePivotInClause()
	if err != nil {
		return nil, err
	}
	clause.In = in
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// Optional DEFAULT ON NULL ( <expr> )
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwNULL); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		dflt, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		clause.DefaultVal = dflt
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	clause.Loc.End = closeTok.Loc.End

	// Optional trailing [AS] alias for the pivot result.
	if alias, has := p.parsePivotTrailingAlias(); has {
		clause.Alias = alias
		clause.Loc.End = p.prev.Loc.End
	}

	return clause, nil
}

// parsePivotColumnRef parses the FOR <col> column reference: one or more
// dotted identifiers. PIVOT names a single column, but we accept a qualified
// name for robustness.
func (p *Parser) parsePivotColumnRef() (*ast.ColumnRef, error) {
	first, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	ref := &ast.ColumnRef{
		Parts: []ast.Ident{first},
		Loc:   first.Loc,
	}
	for p.cur.Type == '.' {
		p.advance() // consume '.'
		part, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ref.Parts = append(ref.Parts, part)
		ref.Loc.End = part.Loc.End
	}
	return ref, nil
}

// parsePivotInClause parses the body of PIVOT … IN ( … ): one of
//   - ANY [ ORDER BY … ]
//   - <subquery>
//   - <value> [AS alias] (, <value> [AS alias])*
//
// The caller has consumed the opening '(' and will consume the closing ')'.
func (p *Parser) parsePivotInClause() (*ast.PivotInClause, error) {
	in := &ast.PivotInClause{
		Loc: ast.Loc{Start: p.cur.Loc.Start},
	}

	switch {
	case p.cur.Type == kwANY:
		p.advance() // consume ANY
		in.Kind = ast.PivotInAny
		if p.cur.Type == kwORDER {
			p.advance() // consume ORDER
			if _, err := p.expect(kwBY); err != nil {
				return nil, err
			}
			items, err := p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			in.OrderBy = items
		}
		in.Loc.End = p.prev.Loc.End
		return in, nil

	case p.cur.Type == kwSELECT || p.cur.Type == kwWITH:
		in.Kind = ast.PivotInSubquery
		var query ast.Node
		var err error
		if p.cur.Type == kwWITH {
			query, err = p.parseWithQueryExpr()
		} else {
			query, err = p.parseQueryExpr()
		}
		if err != nil {
			return nil, err
		}
		in.Subquery = query
		in.Loc.End = ast.NodeLoc(query).End
		return in, nil

	default:
		in.Kind = ast.PivotInValues
		for {
			val, err := p.parsePivotValue()
			if err != nil {
				return nil, err
			}
			in.Values = append(in.Values, val)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		in.Loc.End = p.prev.Loc.End
		return in, nil
	}
}

// parsePivotValue parses one value in a PIVOT … IN ( v [AS alias], … ) list.
func (p *Parser) parsePivotValue() (*ast.PivotValue, error) {
	startLoc := p.cur.Loc
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	val := &ast.PivotValue{
		Value: expr,
		Loc:   ast.Loc{Start: startLoc.Start, End: ast.NodeLoc(expr).End},
	}
	if alias, has := p.parseOptionalAlias(); has {
		val.Alias = alias
		val.Loc.End = p.prev.Loc.End
	}
	return val, nil
}

// parseUnpivotClause parses an UNPIVOT (...) clause. The caller has verified
// that p.cur is kwUNPIVOT.
func (p *Parser) parseUnpivotClause() (*ast.UnpivotClause, error) {
	unpivotTok := p.advance() // consume UNPIVOT
	clause := &ast.UnpivotClause{
		Loc: ast.Loc{Start: unpivotTok.Loc.Start},
	}

	// Optional { INCLUDE | EXCLUDE } NULLS
	switch p.cur.Type {
	case kwINCLUDE:
		p.advance()
		if _, err := p.expect(kwNULLS); err != nil {
			return nil, err
		}
		clause.NullsMode = ast.UnpivotIncludeNulls
	case kwEXCLUDE:
		p.advance()
		if _, err := p.expect(kwNULLS); err != nil {
			return nil, err
		}
		clause.NullsMode = ast.UnpivotExcludeNulls
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// <value_col> FOR <name_col> IN ( <col> [AS alias], … )
	valueCol, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	clause.ValueColumn = valueCol

	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	nameCol, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	clause.NameColumn = nameCol

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		col, err := p.parseUnpivotColumn()
		if err != nil {
			return nil, err
		}
		clause.Columns = append(clause.Columns, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	clause.Loc.End = closeTok.Loc.End

	// Optional trailing alias.
	if alias, has := p.parsePivotTrailingAlias(); has {
		clause.Alias = alias
		clause.Loc.End = p.prev.Loc.End
	}

	return clause, nil
}

// parseUnpivotColumn parses one column in UNPIVOT … IN ( col [AS alias], … ).
func (p *Parser) parseUnpivotColumn() (*ast.UnpivotColumn, error) {
	col, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	uc := &ast.UnpivotColumn{
		Column: col,
		Loc:    col.Loc,
	}
	if alias, has := p.parseOptionalAlias(); has {
		uc.Alias = alias
		uc.Loc.End = p.prev.Loc.End
	}
	return uc, nil
}

// parseSampleClause parses a { SAMPLE | TABLESAMPLE } clause. The caller has
// verified that p.cur is kwSAMPLE or kwTABLESAMPLE.
func (p *Parser) parseSampleClause() (*ast.SampleClause, error) {
	kwTok := p.advance() // consume SAMPLE or TABLESAMPLE
	clause := &ast.SampleClause{
		Loc: ast.Loc{Start: kwTok.Loc.Start},
	}
	if kwTok.Type == kwTABLESAMPLE {
		clause.Keyword = ast.SampleKwTablesample
	} else {
		clause.Keyword = ast.SampleKwSample
	}

	// Optional sampling method.
	switch p.cur.Type {
	case kwBERNOULLI:
		p.advance()
		clause.Method = ast.SampleMethodBernoulli
	case kwROW:
		p.advance()
		clause.Method = ast.SampleMethodRow
	case kwSYSTEM:
		p.advance()
		clause.Method = ast.SampleMethodSystem
	case kwBLOCK:
		p.advance()
		clause.Method = ast.SampleMethodBlock
	}

	// ( <prob> | <n> ROWS )
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	qty, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	clause.Quantity = qty
	if p.cur.Type == kwROWS {
		p.advance() // consume ROWS
		clause.Rows = true
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	clause.Loc.End = p.prev.Loc.End

	// Optional { SEED | REPEATABLE } ( <n> )
	switch p.cur.Type {
	case kwSEED:
		p.advance()
		clause.SeedKind = ast.SampleSeedSeed
	case kwREPEATABLE:
		p.advance()
		clause.SeedKind = ast.SampleSeedRepeatable
	}
	if clause.SeedKind != ast.SampleSeedNone {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		seed, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		clause.Seed = seed
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		clause.Loc.End = closeTok.Loc.End
	}

	return clause, nil
}

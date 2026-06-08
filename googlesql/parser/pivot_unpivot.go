package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// This file is part of the `parser-query-clauses` DAG node. It implements
// GoogleSQL's PIVOT and UNPIVOT table-source operators (GoogleSQLParser.g4
// §2.14 pivot_clause / unpivot_clause and the opt_pivot_or_unpivot_clause_and_
// alias / pivot_or_unpivot_clause_and_aliases productions), a hand-port of
// Google's open-source ZetaSQL reference.
//
// PIVOT/UNPIVOT attach to a FROM source (a path table, a parenthesized
// subquery, or a table-valued-function call — the grammar's table_path_
// expression, table_subquery, and tvf_with_suffixes all carry the operator).
// The operator binds AFTER the source's own `[AS] alias` and BEFORE WITH OFFSET
// / FOR SYSTEM_TIME, matching the grammar's ordering of
// `opt_pivot_or_unpivot_clause_and_alias`. At most one PIVOT or UNPIVOT applies
// to a single source (oracle: `t PIVOT(...) UNPIVOT(...)` rejects).
//
// DIALECT NOTE. PIVOT/UNPIVOT are documented as BigQuery-only (truth1/bigquery
// QUERY-013/014), but the live Cluster Spanner emulator's ZetaSQL ACCEPTS them
// (parsing them, then feature-rejecting with "PIVOT is not supported") — so the
// oracle is authoritative for these forms after all (both polarities), and the
// pivot_unpivot_oracle differential drives the PROVE gate. This is the union
// grammar bytebase needs for BigQuery + Spanner.

// atPivotOrUnpivot reports whether the current token begins a PIVOT or UNPIVOT
// operator. Both PIVOT and UNPIVOT are non-reserved keywords, so this guards
// against eating them as bare table aliases (see atTableAliasStop in
// from_join.go).
func (p *Parser) atPivotOrUnpivot() bool {
	return p.cur.Type == kwPIVOT || p.cur.Type == kwUNPIVOT
}

// parsePivotOrUnpivot parses an optional PIVOT(...) or UNPIVOT(...) operator and
// its trailing `[AS] alias`, storing the result on the source's Pivot / Unpivot
// field (TableExpr or UnnestExpr). The current token is the source's first
// post-alias token. If it is not PIVOT or UNPIVOT this is a no-op (the source
// has no pivot operator). At most one operator is consumed.
func (p *Parser) parsePivotOrUnpivot(src ast.Node) error {
	switch p.cur.Type {
	case kwPIVOT:
		pv, err := p.parsePivotClause()
		if err != nil {
			return err
		}
		attachPivot(src, pv)
	case kwUNPIVOT:
		up, err := p.parseUnpivotClause()
		if err != nil {
			return err
		}
		attachUnpivot(src, up)
	default:
		return nil
	}
	return nil
}

// attachPivot stores a parsed PIVOT clause on a table source (TableExpr /
// UnnestExpr), updating its end Loc.
func attachPivot(src ast.Node, pv *ast.PivotClause) {
	switch s := src.(type) {
	case *ast.TableExpr:
		s.Pivot = pv
		s.Loc.End = pv.Loc.End
	case *ast.UnnestExpr:
		s.Pivot = pv
		s.Loc.End = pv.Loc.End
	}
}

// attachUnpivot stores a parsed UNPIVOT clause on a table source (TableExpr /
// UnnestExpr), updating its end Loc.
func attachUnpivot(src ast.Node, up *ast.UnpivotClause) {
	switch s := src.(type) {
	case *ast.TableExpr:
		s.Unpivot = up
		s.Loc.End = up.Loc.End
	case *ast.UnnestExpr:
		s.Unpivot = up
		s.Loc.End = up.Loc.End
	}
}

// parsePivotClause parses a PIVOT operator (pivot_clause + trailing as_alias):
//
//	PIVOT ( <agg>[, …] FOR <for_expr> IN ( <value>[, …] ) ) [ [AS] alias ]
//
// PIVOT is the current token. The FOR expression is parsed at higher-than-AND
// precedence (expression_higher_prec_than_and) so the following `IN (...)` is
// read as the pivot value list rather than folded into an `expr IN (...)`
// predicate.
func (p *Parser) parsePivotClause() (*ast.PivotClause, error) {
	pivotTok := p.advance() // PIVOT
	pv := &ast.PivotClause{Loc: pivotTok.Loc}

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Aggregate list: pivot_expression (, pivot_expression)* — each
	// `expression [[AS] alias]`. At least one is required (oracle:
	// `PIVOT(FOR q IN ...)` rejects on the missing aggregate).
	aggs, err := p.parsePivotExprList()
	if err != nil {
		return nil, err
	}
	pv.Aggregates = aggs

	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}

	// FOR input column/expression. parseBinaryExpr(bpBitOr) is the
	// expression_higher_prec_than_and band — it stops before IN/comparison, so
	// the trailing IN starts the value list (NOT a predicate on this expression).
	forExpr, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, err
	}
	pv.For = forExpr

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	// Value list: pivot_value (, pivot_value)* — `expression [[AS] alias]`.
	vals, err := p.parsePivotExprList()
	if err != nil {
		return nil, err
	}
	pv.Values = vals
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	pv.Loc.End = closeTok.Loc.End

	// Trailing `[AS] alias` for the pivot result (as_alias?).
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		pv.Alias = alias
		pv.Loc.End = p.prev.Loc.End
	}
	return pv, nil
}

// parsePivotExprList parses a comma-separated list of `expression [[AS] alias]`
// (pivot_expression_list / pivot_value_list). At least one item is required.
// The current token is the first expression; parsing stops at the closing token
// of the enclosing parens (FOR for the aggregate list, ')' for the value list),
// which the caller consumes.
func (p *Parser) parsePivotExprList() ([]*ast.PivotExpr, error) {
	var items []*ast.PivotExpr
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		pe := &ast.PivotExpr{Expr: expr, Loc: nodeLoc(expr)}
		// Optional `[AS] alias` (as_alias).
		if alias, ok, err := p.tryParsePivotAlias(); err != nil {
			return nil, err
		} else if ok {
			pe.Alias = alias
			pe.Loc.End = p.prev.Loc.End
		}
		items = append(items, pe)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ','
	}
	return items, nil
}

// tryParsePivotAlias parses an optional `[AS] alias` for a pivot aggregate /
// value (as_alias). Unlike tryParseTableAlias, a bare (AS-less) identifier is
// accepted as an alias only when it does NOT begin a list/clause continuation
// (FOR / IN / ',' / ')'). An explicit AS requires a following identifier.
func (p *Parser) tryParsePivotAlias() (string, bool, error) {
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok, err := p.expectIdentifier()
		if err != nil {
			return "", false, err
		}
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return "", false, err
		}
		return alias, true, nil
	}
	// A bare identifier that is not FOR/IN (those terminate the aggregate list)
	// and not a separator is an implicit alias. isIdentifierStart excludes
	// reserved keywords; FOR and IN are reserved, so they already stop here.
	if isIdentifierStart(p.cur.Type) {
		aliasTok := p.advance()
		alias, err := p.identifierText(aliasTok)
		if err != nil {
			return "", false, err
		}
		return alias, true, nil
	}
	return "", false, nil
}

// parseUnpivotClause parses an UNPIVOT operator (unpivot_clause + trailing
// as_alias):
//
//	UNPIVOT [ {INCLUDE|EXCLUDE} NULLS ]
//	  ( <value_col(s)> FOR <name_col> IN ( <in_item>[, …] ) ) [ [AS] alias ]
//
// UNPIVOT is the current token. The value-column list is `path_expression_list_
// with_opt_parens` — a bare column for single-column UNPIVOT or a parenthesized
// group `(c1, c2)` for multi-column UNPIVOT.
func (p *Parser) parseUnpivotClause() (*ast.UnpivotClause, error) {
	unpivotTok := p.advance() // UNPIVOT
	up := &ast.UnpivotClause{Loc: unpivotTok.Loc}

	// Optional INCLUDE NULLS | EXCLUDE NULLS (unpivot_nulls_filter).
	switch p.cur.Type {
	case kwINCLUDE:
		p.advance()
		if _, err := p.expect(kwNULLS); err != nil {
			return nil, err
		}
		up.NullsMode = ast.UnpivotIncludeNulls
	case kwEXCLUDE:
		p.advance()
		if _, err := p.expect(kwNULLS); err != nil {
			return nil, err
		}
		up.NullsMode = ast.UnpivotExcludeNulls
	}

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Value column(s): a single path or a parenthesized list (multi-column).
	cols, err := p.parsePathListWithOptParens()
	if err != nil {
		return nil, err
	}
	up.ValueColumns = cols

	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	// Name column (a single path_expression).
	nameCol, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	up.NameColumn = nameCol

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	// IN ( unpivot_in_item (, unpivot_in_item)* ).
	items, err := p.parseUnpivotInItemList()
	if err != nil {
		return nil, err
	}
	up.Items = items

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	up.Loc.End = closeTok.Loc.End

	// Trailing `[AS] alias` (unpivot_alias).
	if alias, ok, err := p.tryParseTableAlias(); err != nil {
		return nil, err
	} else if ok {
		up.Alias = alias
		up.Loc.End = p.prev.Loc.End
	}
	return up, nil
}

// parseUnpivotInItemList parses `( unpivot_in_item (, unpivot_in_item)* )`
// (unpivot_in_item_list). The current token is the opening '('. At least one
// item is required (oracle: `UNPIVOT(... IN ())` rejects). The closing ')' is
// consumed.
func (p *Parser) parseUnpivotInItemList() ([]*ast.UnpivotInItem, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var items []*ast.UnpivotInItem
	for {
		item, err := p.parseUnpivotInItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return items, nil
}

// parseUnpivotInItem parses one IN-list entry (unpivot_in_item):
// `path_expression_list_with_opt_parens [ [AS] string|int ]` — a column or
// parenthesized column group, with an optional string/integer row-value alias.
func (p *Parser) parseUnpivotInItem() (*ast.UnpivotInItem, error) {
	cols, err := p.parsePathListWithOptParens()
	if err != nil {
		return nil, err
	}
	item := &ast.UnpivotInItem{Columns: cols}
	if len(cols) > 0 {
		item.Loc = ast.Loc{Start: cols[0].Loc.Start, End: cols[len(cols)-1].Loc.End}
	}

	// Optional `[AS] string|int` (opt_as_string_or_integer). AS is optional; the
	// alias is a string OR integer literal.
	if p.cur.Type == kwAS {
		p.advance()
		if err := p.parseUnpivotInItemAlias(item); err != nil {
			return nil, err
		}
	} else if p.cur.Type == tokString || p.cur.Type == tokInteger {
		// AS-less alias is only the literal forms (a following column would be a
		// separate IN item, separated by a comma).
		if err := p.parseUnpivotInItemAlias(item); err != nil {
			return nil, err
		}
	}
	return item, nil
}

// parseUnpivotInItemAlias consumes the string- or integer-literal alias of an
// unpivot IN item (opt_as_string_or_integer body) and records it on item. The
// current token must be a string or integer literal.
func (p *Parser) parseUnpivotInItemAlias(item *ast.UnpivotInItem) error {
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		item.HasAlias = true
		item.AliasIsInt = false
		item.AliasString = tok.Str
		item.Loc.End = tok.Loc.End
		return nil
	case tokInteger:
		tok := p.advance()
		item.HasAlias = true
		item.AliasIsInt = true
		item.AliasInt = tok.Str
		item.Loc.End = tok.Loc.End
		return nil
	default:
		return p.syntaxErrorAtCur()
	}
}

// parsePathListWithOptParens parses path_expression_list_with_opt_parens: either
// `( path (, path)* )` (a parenthesized group, used for multi-column UNPIVOT) or
// a bare comma-less single `path`. (Note: outside parens the grammar's
// path_expression_list allows commas, but in the UNPIVOT value-column and
// in-item contexts the comma at the top level separates IN items / belongs to
// the enclosing list, so a bare path is single here; multi-column groups are
// always parenthesized.)
func (p *Parser) parsePathListWithOptParens() ([]*ast.PathExpr, error) {
	if p.cur.Type == int('(') {
		p.advance() // '('
		var paths []*ast.PathExpr
		for {
			path, err := p.parsePathExpr()
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
			if p.cur.Type != int(',') {
				break
			}
			p.advance() // ','
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		return paths, nil
	}
	path, err := p.parsePathExpr()
	if err != nil {
		return nil, err
	}
	return []*ast.PathExpr{path}, nil
}

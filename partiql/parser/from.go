// from.go implements FROM clause parsing for the PartiQL parser.
//
// This file handles the FROM clause, table references (including joins),
// UNPIVOT, and aliasing (AS/AT/BY).
//
// Grammar references cite PartiQLParser.g4 line numbers.
package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parseFromClause parses FROM tableReference.
//
// Grammar: fromClause (lines 297-298):
//
//	FROM tableReference
func (p *Parser) parseFromClause() (ast.TableExpr, error) {
	p.advance() // consume FROM
	return p.parseTableReference()
}

// parseTableReference parses a table reference, handling joins as
// left-associative iteration. This covers:
//   - base table references (expr + aliases)
//   - UNPIVOT
//   - comma joins (implicit cross join)
//   - explicit JOINs (CROSS, INNER, LEFT, RIGHT, FULL, OUTER)
//   - parenthesized table references
//
// Grammar: tableReference (lines 389-395):
//
//	lhs=tableReference joinType? CROSS JOIN rhs=joinRhs   # TableCrossJoin
//	lhs=tableReference COMMA rhs=joinRhs                  # TableCrossJoin
//	lhs=tableReference joinType? JOIN rhs=joinRhs joinSpec # TableQualifiedJoin
//	tableNonJoin                                           # TableRefBase
//	PAREN_LEFT tableReference PAREN_RIGHT                  # TableWrapped
func (p *Parser) parseTableReference() (ast.TableExpr, error) {
	left, err := p.parseTablePrimary()
	if err != nil {
		return nil, err
	}

	// Left-associative join loop.
	for {
		startLoc := left.GetLoc().Start

		// Comma join: table1, table2 -> CROSS JOIN.
		if p.cur.Type == tokCOMMA {
			p.advance() // consume ,
			right, err := p.parseJoinRhs()
			if err != nil {
				return nil, err
			}
			left = &ast.JoinExpr{
				Kind:  ast.JoinKindCross,
				Left:  left,
				Right: right,
				Loc:   ast.Loc{Start: startLoc, End: right.GetLoc().End},
			}
			continue
		}

		// Explicit join: [joinType] JOIN/CROSS JOIN.
		joinKind, hasJoin := p.tryParseJoinType()
		if !hasJoin {
			break
		}

		right, err := p.parseJoinRhs()
		if err != nil {
			return nil, err
		}

		var on ast.ExprNode
		end := right.GetLoc().End

		// ON clause is required for non-CROSS joins.
		if joinKind != ast.JoinKindCross {
			if p.cur.Type != tokON {
				return nil, &ParseError{
					Message: fmt.Sprintf("expected ON after JOIN, got %q", p.cur.Str),
					Loc:     p.cur.Loc,
				}
			}
			p.advance() // consume ON
			on, err = p.parseExprTop()
			if err != nil {
				return nil, err
			}
			end = on.GetLoc().End
		}

		left = &ast.JoinExpr{
			Kind:  joinKind,
			Left:  left,
			Right: right,
			On:    on,
			Loc:   ast.Loc{Start: startLoc, End: end},
		}
	}

	return left, nil
}

// parseTablePrimary parses a non-join table reference: a base table
// reference, UNPIVOT, or a parenthesized table reference.
//
// A leading '(' is grammatically ambiguous between two productions:
//
//	tableReference#TableRefBase  whose source is a parenthesized
//	    value expression / subquery — e.g. (t1), (SELECT ...), ((t1)).
//	    This reading DOES accept a trailing AS/AT/BY alias chain, so it
//	    is what the bare-/aliased cases resolve to.
//	tableReference#TableWrapped  — `( tableReference )` (PartiQLParser.g4
//	    line 394), where the parenthesized content is a join, comma
//	    cross-join, or inner-aliased table that the value-expression
//	    grammar cannot consume up to the matching ')'. The parentheses
//	    are purely structural (the inner table expression is returned
//	    directly), and this reading does NOT accept a trailing (outer)
//	    alias.
//
// The generated ANTLR parser resolves this with unbounded lookahead,
// preferring TableRefBase whenever the value-expression reading reaches
// the closing ')'. We reproduce that faithfully by trying the base
// reference first and, only if it fails, rewinding and parsing the
// TableWrapped form. Bug B6: an earlier version always routed '(' through
// the base/value-expression path, so a wrapped join such as
// "(t1 CROSS JOIN t2)" was rejected at the first join keyword.
func (p *Parser) parseTablePrimary() (ast.TableExpr, error) {
	if p.cur.Type == tokPAREN_LEFT {
		// First attempt: '(' begins a value expression / subquery used as
		// a base table reference (handles (expr), (SELECT...), and a
		// trailing AS/AT/BY alias). parseTableBaseReference -> parseSelectExpr
		// -> the expression-level parseParenExpr.
		start := p.save()
		base, err := p.parseTableBaseReference()
		if err == nil {
			return base, nil
		}
		// The value-expression reading failed (e.g. the content is a join
		// or carries an inner alias). Rewind and try the wrapped
		// table-reference form `( tableReference )`.
		p.restore(start)
		wrapped, werr := p.parseWrappedTableReference()
		if werr != nil {
			// Neither reading applies. Surface the wrapped-form error: at
			// this point we know the content was not a plain value
			// expression, so the table-reference-oriented diagnostic is the
			// more relevant one.
			return nil, werr
		}
		return wrapped, nil
	}

	// UNPIVOT.
	if p.cur.Type == tokUNPIVOT {
		return p.parseTableUnpivot()
	}

	// Base table reference: source expression + optional aliases.
	return p.parseTableBaseReference()
}

// parseWrappedTableReference parses the TableWrapped form
// `( tableReference )` (PartiQLParser.g4 line 394). The parentheses are
// purely structural: the inner table expression is returned directly with
// no wrapper node and no outer alias (matching how parseJoinRhs unwraps the
// JoinRhsTableJoined form). A trailing alias is deliberately left
// unconsumed so the caller rejects it at the next-token / EOF check, which
// is what the generated ANTLR grammar does ("(t1 CROSS JOIN t2) AS x" is a
// syntax error).
func (p *Parser) parseWrappedTableReference() (ast.TableExpr, error) {
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	inner, err := p.parseTableReference()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, err
	}
	return inner, nil
}

// parseJoinRhs parses the right-hand side of a JOIN: either a
// non-join table reference or a parenthesized table reference.
//
// Grammar: joinRhs (lines 411-414):
//
//	tableNonJoin                           # JoinRhsBase
//	PAREN_LEFT tableReference PAREN_RIGHT  # JoinRhsTableJoined
func (p *Parser) parseJoinRhs() (ast.TableExpr, error) {
	if p.cur.Type == tokPAREN_LEFT {
		p.advance() // consume (
		inner, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokPAREN_RIGHT); err != nil {
			return nil, err
		}
		return inner, nil
	}

	// UNPIVOT.
	if p.cur.Type == tokUNPIVOT {
		return p.parseTableUnpivot()
	}

	return p.parseTableBaseReference()
}

// parseTableBaseReference parses a source expression with optional
// AS/AT/BY aliases.
//
// Grammar: tableBaseReference (lines 402-406):
//
//	source=exprSelect symbolPrimitive                   # TableBaseRefSymbol
//	source=exprSelect asIdent? atIdent? byIdent?        # TableBaseRefClauses
//	source=exprGraphMatchOne asIdent? atIdent? byIdent? # TableBaseRefMatch
func (p *Parser) parseTableBaseReference() (ast.TableExpr, error) {
	source, err := p.parseSelectExpr()
	if err != nil {
		return nil, err
	}

	// FROM-clause graph match (tableBaseReference#TableBaseRefMatch, g4:405):
	//   source=exprGraphMatchOne asIdent? atIdent? byIdent?
	// where exprGraphMatchOne is `exprPrimary MATCH gpmlPattern` (g4:628-629) —
	// the unparenthesised single-pattern form, e.g. `FROM g MATCH (a)-[e]->(b)`.
	// parseSelectExpr above parses just the `exprPrimary` (the graph source) and
	// leaves MATCH unconsumed, so we dispatch here. The resulting MatchExpr is a
	// TableExpr, so it flows through the alias handling below like any other FROM
	// source (the parenthesised many-pattern form is instead parsed inside
	// parseParenExpr -> parseGraphMatch and arrives here already as a MatchExpr).
	if p.cur.Type == tokMATCH {
		match, matchErr := p.parseGraphMatchOne(source)
		if matchErr != nil {
			return nil, matchErr
		}
		source = match
	}

	// Try to parse optional aliases. We detect bare alias (implicit AS)
	// by checking if the current token is an identifier that is NOT a
	// keyword that starts a clause. The grammar rule
	// tableBaseReference#TableBaseRefSymbol allows a bare symbolPrimitive
	// alias without AS.
	as, at, by, endLoc, err := p.parseTableAliases()
	if err != nil {
		return nil, err
	}

	// If no aliases, the source expression can serve directly as a
	// TableExpr if it implements the interface (VarRef, PathExpr, SubLink).
	if as == nil && at == nil && by == nil {
		if te, ok := source.(ast.TableExpr); ok {
			return te, nil
		}
		// Wrap in AliasedSource with no aliases for non-TableExpr sources.
		return &ast.AliasedSource{
			Source: exprAsTableExpr(source),
			Loc:    source.GetLoc(),
		}, nil
	}

	return &ast.AliasedSource{
		Source: exprAsTableExpr(source),
		As:     as,
		At:     at,
		By:     by,
		Loc:    ast.Loc{Start: source.GetLoc().Start, End: endLoc},
	}, nil
}

// exprAsTableExpr converts an ExprNode to a TableExpr. If the expr
// already implements TableExpr (VarRef, PathExpr, SubLink), it returns
// it directly. Otherwise it wraps it in an AliasedSource with no
// aliases (the expression itself becomes the source).
func exprAsTableExpr(expr ast.ExprNode) ast.TableExpr {
	if te, ok := expr.(ast.TableExpr); ok {
		return te
	}
	// For expressions that are not directly TableExpr (e.g., function
	// calls, literals used as table sources), wrap in AliasedSource.
	return &ast.AliasedSource{
		Source: nil, // will be overwritten by caller
		Loc:    expr.GetLoc(),
	}
}

// parseTableAliases parses the optional AS/AT/BY alias chain on a
// table reference.
//
// Grammar fragments:
//
//	asIdent: AS symbolPrimitive
//	atIdent: AT symbolPrimitive
//	byIdent: BY symbolPrimitive
//
// Also handles the implicit AS form (bare symbolPrimitive) per
// TableBaseRefSymbol (line 403) and TableBaseRefClauses (line 404).
func (p *Parser) parseTableAliases() (as, at, by *string, endLoc int, err error) {
	// AS alias (explicit or implicit).
	if p.cur.Type == tokAS {
		p.advance() // consume AS
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, nil, nil, 0, parseErr
		}
		as = &name
		endLoc = nameLoc.End
	} else if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		// Implicit AS: bare identifier. Only if it does not look like a
		// keyword that starts the next clause. The lexer classifies
		// keywords as their specific tok constants, not tokIDENT, so if
		// cur.Type is tokIDENT it is safe to treat as an alias.
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, nil, nil, 0, parseErr
		}
		as = &name
		endLoc = nameLoc.End
	}

	// AT alias.
	if p.cur.Type == tokAT {
		p.advance() // consume AT
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, nil, nil, 0, parseErr
		}
		at = &name
		endLoc = nameLoc.End
	}

	// BY alias.
	if p.cur.Type == tokBY {
		p.advance() // consume BY
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, nil, nil, 0, parseErr
		}
		by = &name
		endLoc = nameLoc.End
	}

	return as, at, by, endLoc, nil
}

// parseTableUnpivot parses UNPIVOT expr [AS alias] [AT pos] [BY key].
//
// Grammar: tableUnpivot (lines 408-409):
//
//	UNPIVOT expr asIdent? atIdent? byIdent?
func (p *Parser) parseTableUnpivot() (*ast.UnpivotExpr, error) {
	start := p.cur.Loc.Start
	p.advance() // consume UNPIVOT

	source, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	end := source.GetLoc().End

	// Optional aliases.
	var as, at, by *string

	if p.cur.Type == tokAS {
		p.advance()
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, parseErr
		}
		as = &name
		end = nameLoc.End
	}

	if p.cur.Type == tokAT {
		p.advance()
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, parseErr
		}
		at = &name
		end = nameLoc.End
	}

	if p.cur.Type == tokBY {
		p.advance()
		name, _, nameLoc, parseErr := p.parseSymbolPrimitive()
		if parseErr != nil {
			return nil, parseErr
		}
		by = &name
		end = nameLoc.End
	}

	return &ast.UnpivotExpr{
		Source: source,
		As:     as,
		At:     at,
		By:     by,
		Loc:    ast.Loc{Start: start, End: end},
	}, nil
}

// tryParseJoinType attempts to parse a join type. Returns the JoinKind
// and true if a complete join keyword sequence (a modifier and/or CROSS
// followed by JOIN, or a bare JOIN) was found, or (Invalid, false) if the
// current token does not start a join.
//
// IMPORTANT: when a join *modifier* (LEFT/RIGHT/FULL/INNER/OUTER, with an
// optional trailing OUTER, or CROSS) is present but is NOT followed by the
// required JOIN keyword, this is NOT a join — and the consumed modifier
// token(s) are rewound (save/restore) so they remain in the token stream.
// The caller then breaks out of the join loop with the modifier still
// pending, which makes a dangling modifier ("FROM t1 LEFT") fail at the
// next-token / EOF check rather than being silently swallowed. Bug B5: an
// earlier version returned (Invalid, false) here WITHOUT restoring, which
// dropped the keyword and wrongly accepted dangling-modifier inputs that
// the generated ANTLR grammar rejects.
//
// Grammar: joinType (lines 419-425):
//
//	mod=INNER
//	mod=LEFT OUTER?
//	mod=RIGHT OUTER?
//	mod=FULL OUTER?
//	mod=OUTER
//
// Plus the CROSS JOIN form (line 390).
func (p *Parser) tryParseJoinType() (ast.JoinKind, bool) {
	// Snapshot so a modifier not followed by JOIN can be fully rewound.
	start := p.save()

	switch p.cur.Type {
	case tokCROSS:
		p.advance() // consume CROSS
		if p.cur.Type != tokJOIN {
			// CROSS must be followed by JOIN.
			p.restore(start)
			return ast.JoinKindInvalid, false
		}
		p.advance() // consume JOIN
		return ast.JoinKindCross, true

	case tokINNER:
		p.advance() // consume INNER
		if p.cur.Type != tokJOIN {
			p.restore(start)
			return ast.JoinKindInvalid, false
		}
		p.advance() // consume JOIN
		return ast.JoinKindInner, true

	case tokLEFT:
		p.advance() // consume LEFT
		p.match(tokOUTER)
		if p.cur.Type != tokJOIN {
			p.restore(start)
			return ast.JoinKindInvalid, false
		}
		p.advance() // consume JOIN
		return ast.JoinKindLeft, true

	case tokRIGHT:
		p.advance() // consume RIGHT
		p.match(tokOUTER)
		if p.cur.Type != tokJOIN {
			p.restore(start)
			return ast.JoinKindInvalid, false
		}
		p.advance() // consume JOIN
		return ast.JoinKindRight, true

	case tokFULL:
		p.advance() // consume FULL
		p.match(tokOUTER)
		if p.cur.Type != tokJOIN {
			p.restore(start)
			return ast.JoinKindInvalid, false
		}
		p.advance() // consume JOIN
		return ast.JoinKindFull, true

	case tokOUTER:
		p.advance() // consume OUTER
		if p.cur.Type != tokJOIN {
			p.restore(start)
			return ast.JoinKindInvalid, false
		}
		p.advance() // consume JOIN
		return ast.JoinKindOuter, true

	case tokJOIN:
		// Bare JOIN => defaults to INNER.
		p.advance() // consume JOIN
		return ast.JoinKindInner, true

	default:
		return ast.JoinKindInvalid, false
	}
}

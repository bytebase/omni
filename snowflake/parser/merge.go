package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// MERGE statement parser
// ---------------------------------------------------------------------------

// parseMergeStmt parses a MERGE statement:
//
//	MERGE INTO target [alias] USING source [alias] ON cond
//	  [WHEN MATCHED [AND cond] THEN {UPDATE SET ... | DELETE}]...
//	  [WHEN NOT MATCHED [BY TARGET|SOURCE] [AND cond] THEN {INSERT ...}]...
func (p *Parser) parseMergeStmt() (*ast.MergeStmt, error) {
	mergeTok, err := p.expect(kwMERGE)
	if err != nil {
		return nil, err
	}

	// INTO is required
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	target, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.MergeStmt{
		Target: target,
		Loc:    ast.Loc{Start: mergeTok.Loc.Start},
	}

	// Optional target alias
	targetAlias, hasTargetAlias := p.parseOptionalAlias()
	if hasTargetAlias {
		stmt.TargetAlias = targetAlias
	}

	// USING source
	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}

	source, err := p.parseMergeSource()
	if err != nil {
		return nil, err
	}
	stmt.Source = source

	// Optional source alias
	sourceAlias, hasSourceAlias := p.parseOptionalAlias()
	if hasSourceAlias {
		stmt.SourceAlias = sourceAlias
	}

	// ON condition
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	on, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.On = on

	// WHEN clauses (one or more)
	for p.cur.Type == kwWHEN {
		when, err := p.parseMergeWhen()
		if err != nil {
			return nil, err
		}
		stmt.Whens = append(stmt.Whens, when)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseMergeSource parses the USING source: either a table reference or
// a parenthesized subquery. It does NOT parse the alias; the caller parses
// the optional alias after this function returns.
func (p *Parser) parseMergeSource() (ast.Node, error) {
	if p.cur.Type == '(' {
		startLoc := p.cur.Loc
		next := p.peekNext()
		if next.Type == kwSELECT || next.Type == kwWITH {
			p.advance() // consume '('
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
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			ref := &ast.TableRef{
				Subquery: query,
				Loc:      ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
			}
			return ref, nil
		}
	}
	// Simple table reference: parse name only, NOT alias.
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	return &ast.TableRef{
		Name: name,
		Loc:  name.Loc,
	}, nil
}

// parseMergeWhen parses one WHEN clause:
//
//	WHEN MATCHED [AND cond] THEN {UPDATE SET ... | DELETE}
//	WHEN NOT MATCHED [BY TARGET|SOURCE] [AND cond] THEN INSERT ...
func (p *Parser) parseMergeWhen() (*ast.MergeWhen, error) {
	whenTok, err := p.expect(kwWHEN)
	if err != nil {
		return nil, err
	}

	when := &ast.MergeWhen{
		Loc: ast.Loc{Start: whenTok.Loc.Start},
	}

	if p.cur.Type == kwMATCHED {
		p.advance() // consume MATCHED
		when.Matched = true
	} else if p.cur.Type == kwNOT {
		p.advance() // consume NOT
		if _, err := p.expect(kwMATCHED); err != nil {
			return nil, err
		}
		when.Matched = false
		// Optional BY TARGET|SOURCE
		if p.cur.Type == kwBY {
			p.advance() // consume BY
			switch {
			case p.cur.Type == kwSOURCE:
				p.advance() // consume SOURCE
				when.BySource = true
			case p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "TARGET":
				p.advance() // consume TARGET (identifier)
				when.ByTarget = true
			default:
				// Unexpected token after BY
				return nil, p.syntaxErrorAtCur()
			}
		} else {
			// Plain WHEN NOT MATCHED defaults to BY TARGET
			when.ByTarget = true
		}
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	// Optional AND condition
	if p.cur.Type == kwAND {
		p.advance() // consume AND
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		when.AndCond = cond
	}

	// THEN action
	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwUPDATE:
		p.advance() // consume UPDATE
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		sets, err := p.parseUpdateSetList()
		if err != nil {
			return nil, err
		}
		when.Action = ast.MergeActionUpdate
		when.Sets = sets

	case kwDELETE:
		p.advance() // consume DELETE
		when.Action = ast.MergeActionDelete

	case kwINSERT:
		p.advance() // consume INSERT
		when.Action = ast.MergeActionInsert

		// Optional column list
		if p.cur.Type == '(' {
			// Peek to see if it's a column list or VALUES
			next := p.peekNext()
			if next.Type != kwVALUES && next.Type != kwDEFAULT {
				p.advance() // consume '('
				cols, err := p.parseIdentList()
				if err != nil {
					return nil, err
				}
				when.InsertCols = cols
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}
		}

		// VALUES (exprs) or VALUES DEFAULT
		if _, err := p.expect(kwVALUES); err != nil {
			return nil, err
		}

		if p.cur.Type == kwDEFAULT {
			p.advance() // consume DEFAULT
			when.InsertDefault = true
		} else {
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			vals, err := p.parseExprList()
			if err != nil {
				return nil, err
			}
			when.InsertVals = vals
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	when.Loc.End = p.prev.Loc.End
	return when, nil
}

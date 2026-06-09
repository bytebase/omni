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
		when.Action = ast.MergeActionUpdate
		// `UPDATE ALL BY NAME` is an alternative to `UPDATE SET col = val, ...`.
		if p.matchAllByName() {
			when.AllByName = true
			break
		}
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		sets, err := p.parseUpdateSetList()
		if err != nil {
			return nil, err
		}
		when.Sets = sets

	case kwDELETE:
		p.advance() // consume DELETE
		when.Action = ast.MergeActionDelete

	case kwINSERT:
		p.advance() // consume INSERT
		when.Action = ast.MergeActionInsert

		// `INSERT ALL BY NAME` is an alternative to `INSERT (cols) VALUES (vals)`.
		if p.matchAllByName() {
			when.AllByName = true
			break
		}

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

// matchAllByName consumes the `ALL BY NAME` token sequence that introduces the
// MERGE column-list-by-name action forms (`UPDATE ALL BY NAME`,
// `INSERT ALL BY NAME`) and reports whether it was present. It only advances
// when the full three-token sequence matches, so a partial run is left intact
// for the caller's normal SET / VALUES handling. ALL, BY, and NAME are all
// non-reserved keywords (kwALL/kwBY/kwNAME).
func (p *Parser) matchAllByName() bool {
	if p.cur.Type != kwALL {
		return false
	}
	if p.peekNext().Type != kwBY {
		return false
	}
	// Confirm `NAME` follows `BY` before consuming anything. Snapshot the
	// token-stream state, advance past ALL/BY, peek for NAME, and restore if
	// the sequence is incomplete (no allocation, bounded by three tokens).
	savedCur := p.cur
	savedPrev := p.prev
	savedNextBuf := p.nextBuf
	savedHasNext := p.hasNext
	savedLexPos := p.lexer.pos

	p.advance() // consume ALL (cur is now BY)
	p.advance() // consume BY  (cur is the token after BY)
	if p.cur.Type == kwNAME {
		p.advance() // consume NAME
		return true
	}

	// Not `ALL BY NAME` — restore and report no match.
	p.cur = savedCur
	p.prev = savedPrev
	p.nextBuf = savedNextBuf
	p.hasNext = savedHasNext
	p.lexer.pos = savedLexPos
	return false
}

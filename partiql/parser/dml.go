package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parseInsertStmt parses both legacy and RFC 0011 INSERT forms:
//
//	INSERT INTO pathSimple VALUE value=expr ( AT pos=expr )? onConflictClause? returningClause?
//	INSERT INTO symbolPrimitive asIdent? value=expr onConflictClause?
//
// Disambiguation: parse pathSimple, then if the next token is VALUE
// keyword, it is the legacy form; otherwise it is the RFC 0011 form
// (which only allows a symbolPrimitive, not a multi-step path).
//
// Grammar: insertCommand#InsertLegacy, insertCommand#Insert,
// insertCommandReturning (PartiQLParser.g4 lines 130-137).
func (p *Parser) parseInsertStmt() (*ast.InsertStmt, error) {
	start := p.cur.Loc.Start
	p.advance() // consume INSERT
	if _, err := p.expect(tokINTO); err != nil {
		return nil, err
	}

	// Parse pathSimple (works for both forms since symbolPrimitive is
	// a prefix of pathSimple).
	path, err := p.parsePathSimple()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == tokVALUE {
		// Legacy form: INSERT INTO pathSimple VALUE expr [AT pos] [ON CONFLICT ...] [RETURNING ...]
		return p.parseInsertLegacy(start, path)
	}

	// RFC 0011 form: INSERT INTO symbolPrimitive [AS alias] expr [ON CONFLICT ...]
	// The target must be a bare symbolPrimitive (no path steps).
	if len(path.Steps) > 0 {
		return nil, &ParseError{
			Message: "INSERT INTO (RFC 0011 form) target must be a simple identifier, not a path",
			Loc:     path.Loc,
		}
	}
	return p.parseInsertRFC0011(start, path.Root.(*ast.VarRef))
}

// parseInsertLegacy finishes parsing:
//
//	... VALUE expr [AT pos] [ON CONFLICT ...] [RETURNING ...]
func (p *Parser) parseInsertLegacy(start int, target *ast.PathExpr) (*ast.InsertStmt, error) {
	p.advance() // consume VALUE

	value, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	var pos ast.ExprNode
	if p.cur.Type == tokAT {
		p.advance() // consume AT
		pos, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
	}

	var oc *ast.OnConflict
	if p.cur.Type == tokON {
		oc, err = p.parseOnConflict()
		if err != nil {
			return nil, err
		}
	}

	var ret *ast.ReturningClause
	if p.cur.Type == tokRETURNING {
		ret, err = p.parseReturningClause()
		if err != nil {
			return nil, err
		}
	}

	end := p.prev.Loc.End
	return &ast.InsertStmt{
		Target:     target,
		Value:      value,
		Pos:        pos,
		OnConflict: oc,
		Returning:  ret,
		Loc:        ast.Loc{Start: start, End: end},
	}, nil
}

// parseInsertRFC0011 finishes parsing:
//
//	... [AS alias] expr [ON CONFLICT ...]
func (p *Parser) parseInsertRFC0011(start int, target *ast.VarRef) (*ast.InsertStmt, error) {
	var alias *string
	if p.cur.Type == tokAS {
		p.advance() // consume AS
		name, _, _, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		alias = &name
	}

	value, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	var oc *ast.OnConflict
	if p.cur.Type == tokON {
		oc, err = p.parseOnConflict()
		if err != nil {
			return nil, err
		}
	}

	end := p.prev.Loc.End
	return &ast.InsertStmt{
		Target:     target,
		AsAlias:    alias,
		Value:      value,
		OnConflict: oc,
		Loc:        ast.Loc{Start: start, End: end},
	}, nil
}

// parseDeleteStmt parses:
//
//	DELETE fromClauseSimple [WHERE expr] [RETURNING ...]
//
// Grammar: deleteCommand (PartiQLParser.g4 line 191-192).
func (p *Parser) parseDeleteStmt() (*ast.DeleteStmt, error) {
	start := p.cur.Loc.Start
	p.advance() // consume DELETE

	source, err := p.parseFromClauseSimple()
	if err != nil {
		return nil, err
	}

	var where ast.ExprNode
	if p.cur.Type == tokWHERE {
		p.advance() // consume WHERE
		where, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
	}

	var ret *ast.ReturningClause
	if p.cur.Type == tokRETURNING {
		ret, err = p.parseReturningClause()
		if err != nil {
			return nil, err
		}
	}

	end := p.prev.Loc.End
	return &ast.DeleteStmt{
		Source:    source,
		Where:     where,
		Returning: ret,
		Loc:       ast.Loc{Start: start, End: end},
	}, nil
}

// parseUpdateStmt parses the UPDATE form of dml#DmlBaseWrapper (1st alt):
//
//	updateClause dmlBaseCommand+ whereClause? returningClause?
//
// where updateClause is `UPDATE tableBaseReference` and dmlBaseCommand is
// one of setCommand / insertCommand / replaceCommand / removeCommand /
// upsertCommand. All of `UPDATE t SET a=1`, `UPDATE t REMOVE a.b`,
// `UPDATE t SET a=1 REMOVE b`, `UPDATE t INSERT INTO a VALUE 1`, … are
// accepted by the grammar (confirmed against the generated ANTLR parser).
//
// Note the clause order: in this alt the optional whereClause follows the
// commands (`UPDATE t SET a=1 WHERE x=1`), unlike the FROM-led alt where
// WHERE precedes them.
//
// The source is parsed as a pathSimple with optional aliases since this
// branch does not use the full FROM parser. For cases like
// `UPDATE "Music" SET ...`, pathSimple + aliases is sufficient.
//
// Grammar: updateClause + dml#DmlBaseWrapper (PartiQLParser.g4 lines 94-96, 182-183).
func (p *Parser) parseUpdateStmt() (*ast.UpdateStmt, error) {
	start := p.cur.Loc.Start
	p.advance() // consume UPDATE

	// Parse the source table reference. In the grammar this is
	// tableBaseReference which is exprSelect with optional aliases.
	// Since we do not have SELECT parsing on this branch, we parse
	// a pathSimple with optional AS/AT/BY aliases.
	source, err := p.parseUpdateSource()
	if err != nil {
		return nil, err
	}

	// dmlBaseCommand+ : at least one command is required.
	sets, commands, err := p.parseDmlBaseCommands()
	if err != nil {
		return nil, err
	}

	// whereClause? — AFTER the commands in this alternative.
	var where ast.ExprNode
	if p.cur.Type == tokWHERE {
		p.advance() // consume WHERE
		where, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
	}

	// returningClause?
	var ret *ast.ReturningClause
	if p.cur.Type == tokRETURNING {
		ret, err = p.parseReturningClause()
		if err != nil {
			return nil, err
		}
	}

	end := p.prev.Loc.End
	return &ast.UpdateStmt{
		Source:    source,
		Sets:      sets,
		Commands:  commands,
		Where:     where,
		Returning: ret,
		Loc:       ast.Loc{Start: start, End: end},
	}, nil
}

// parseFromLedDml parses the FROM-led form of dml#DmlBaseWrapper (2nd alt):
//
//	fromClause whereClause? dmlBaseCommand+ returningClause?
//
// e.g. `FROM t SET a=1`, `FROM t WHERE x=1 SET a=1`, `FROM t REMOVE a.b`,
// `FROM a, b SET x=1` (joins). The source is a full tableReference (handled
// by parseFromClause), so joins / paths / parenthesised sources all work.
//
// Note the clause order differs from the UPDATE alt: here the optional
// whereClause PRECEDES the commands, and there is no trailing whereClause
// (a trailing WHERE is a syntax error, matching the ANTLR oracle).
//
// The result is represented as an UpdateStmt with From=true.
//
// Grammar: fromClause + dml#DmlBaseWrapper (PartiQLParser.g4 lines 94-97, 297-298).
func (p *Parser) parseFromLedDml() (*ast.UpdateStmt, error) {
	start := p.cur.Loc.Start

	// fromClause : FROM tableReference (parseFromClause consumes FROM).
	source, err := p.parseFromClause()
	if err != nil {
		return nil, err
	}

	// whereClause? — BEFORE the commands in this alternative.
	var where ast.ExprNode
	if p.cur.Type == tokWHERE {
		p.advance() // consume WHERE
		where, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
	}

	// dmlBaseCommand+ : at least one command is required.
	sets, commands, err := p.parseDmlBaseCommands()
	if err != nil {
		return nil, err
	}

	// returningClause?
	var ret *ast.ReturningClause
	if p.cur.Type == tokRETURNING {
		ret, err = p.parseReturningClause()
		if err != nil {
			return nil, err
		}
	}

	end := p.prev.Loc.End
	return &ast.UpdateStmt{
		Source:    source,
		From:      true,
		Sets:      sets,
		Commands:  commands,
		Where:     where,
		Returning: ret,
		Loc:       ast.Loc{Start: start, End: end},
	}, nil
}

// parseDmlBaseCommands parses the shared `dmlBaseCommand+` sequence used by
// both alternatives of dml#DmlBaseWrapper. It requires at least one command
// and stops at the first token that does not start a dmlBaseCommand.
//
// It returns:
//   - sets: every SET assignment, flattened in order (the dominant case);
//   - commands: the non-SET commands (*InsertStmt/*ReplaceStmt/*RemoveStmt/
//     *UpsertStmt) in order, nil when the statement only has SET assignments.
//
// At least one command (SET or otherwise) must be present, so a statement
// that has neither is a syntax error — the ANTLR oracle rejects such input.
//
// Grammar: dmlBaseCommand (g4 lines 102-108) =
// insertCommand | setCommand | replaceCommand | removeCommand | upsertCommand.
func (p *Parser) parseDmlBaseCommands() (sets []*ast.SetAssignment, commands []ast.Node, err error) {
	seen := false
	for {
		switch p.cur.Type {
		case tokSET:
			s, serr := p.parseSetCommand()
			if serr != nil {
				return nil, nil, serr
			}
			sets = append(sets, s...)
			seen = true
		case tokINSERT:
			ins, ierr := p.parseInsertStmt()
			if ierr != nil {
				return nil, nil, ierr
			}
			commands = append(commands, ins)
			seen = true
		case tokREPLACE:
			rep, rerr := p.parseReplaceStmt()
			if rerr != nil {
				return nil, nil, rerr
			}
			commands = append(commands, rep)
			seen = true
		case tokREMOVE:
			rem, rerr := p.parseRemoveStmt()
			if rerr != nil {
				return nil, nil, rerr
			}
			commands = append(commands, rem)
			seen = true
		case tokUPSERT:
			ups, uerr := p.parseUpsertStmt()
			if uerr != nil {
				return nil, nil, uerr
			}
			commands = append(commands, ups)
			seen = true
		default:
			if !seen {
				return nil, nil, &ParseError{
					// "expected SET" leads the message because SET is the
					// dominant DML command; the parenthetical lists the rest.
					Message: fmt.Sprintf("expected SET or another DML command (INSERT, REPLACE, REMOVE, UPSERT), got %q", p.cur.Str),
					Loc:     p.cur.Loc,
				}
			}
			return sets, commands, nil
		}
	}
}

// parseUpdateSource parses the table reference after UPDATE. This is
// `tableBaseReference` in the grammar, which is `exprSelect [AS a] [AT a] [BY a]`.
// Since we do not have SELECT on this branch, we parse a symbolPrimitive
// (possibly quoted) with optional AS alias. This covers the common cases:
//
//	UPDATE t SET ...
//	UPDATE "Music" SET ...
func (p *Parser) parseUpdateSource() (ast.TableExpr, error) {
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	ref := &ast.VarRef{
		Name:          name,
		CaseSensitive: caseSensitive,
		Loc:           nameLoc,
	}

	// Check for optional AS alias, AT alias, BY alias.
	var asAlias, atAlias, byAlias *string
	end := nameLoc.End

	if p.cur.Type == tokAS {
		p.advance() // consume AS
		a, _, aLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		asAlias = &a
		end = aLoc.End
	} else if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		// Implicit alias (symbolPrimitive without AS keyword).
		// Only applies if the next token looks like an alias, not SET.
		// We need to be careful: `UPDATE t SET ...` means t is the table and SET is the command.
		// An implicit alias is only for `tableBaseRefSymbol` which is `exprSelect symbolPrimitive`.
		// We skip implicit alias parsing since it conflicts with SET.
	}

	if p.cur.Type == tokAT {
		p.advance() // consume AT
		a, _, aLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		atAlias = &a
		end = aLoc.End
	}
	if p.cur.Type == tokBY {
		p.advance() // consume BY
		a, _, aLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		byAlias = &a
		end = aLoc.End
	}

	if asAlias != nil || atAlias != nil || byAlias != nil {
		return &ast.AliasedSource{
			Source: ref,
			As:     asAlias,
			At:     atAlias,
			By:     byAlias,
			Loc:    ast.Loc{Start: nameLoc.Start, End: end},
		}, nil
	}
	return ref, nil
}

// parseReplaceStmt parses:
//
//	REPLACE INTO symbolPrimitive [AS alias] value=expr
//
// Grammar: replaceCommand (PartiQLParser.g4 line 120-121).
func (p *Parser) parseReplaceStmt() (*ast.ReplaceStmt, error) {
	start := p.cur.Loc.Start
	p.advance() // consume REPLACE
	if _, err := p.expect(tokINTO); err != nil {
		return nil, err
	}

	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	target := &ast.VarRef{
		Name:          name,
		CaseSensitive: caseSensitive,
		Loc:           nameLoc,
	}

	var alias *string
	if p.cur.Type == tokAS {
		p.advance() // consume AS
		a, _, _, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		alias = &a
	}

	value, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	end := p.prev.Loc.End
	return &ast.ReplaceStmt{
		Target:  target,
		AsAlias: alias,
		Value:   value,
		Loc:     ast.Loc{Start: start, End: end},
	}, nil
}

// parseUpsertStmt parses:
//
//	UPSERT INTO symbolPrimitive [AS alias] value=expr
//
// Grammar: upsertCommand (PartiQLParser.g4 lines 124-125).
func (p *Parser) parseUpsertStmt() (*ast.UpsertStmt, error) {
	start := p.cur.Loc.Start
	p.advance() // consume UPSERT
	if _, err := p.expect(tokINTO); err != nil {
		return nil, err
	}

	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	target := &ast.VarRef{
		Name:          name,
		CaseSensitive: caseSensitive,
		Loc:           nameLoc,
	}

	var alias *string
	if p.cur.Type == tokAS {
		p.advance() // consume AS
		a, _, _, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		alias = &a
	}

	value, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}

	end := p.prev.Loc.End
	return &ast.UpsertStmt{
		Target:  target,
		AsAlias: alias,
		Value:   value,
		Loc:     ast.Loc{Start: start, End: end},
	}, nil
}

// parseRemoveStmt parses:
//
//	REMOVE pathSimple
//
// Grammar: removeCommand (PartiQLParser.g4 lines 127-128).
func (p *Parser) parseRemoveStmt() (*ast.RemoveStmt, error) {
	start := p.cur.Loc.Start
	p.advance() // consume REMOVE

	path, err := p.parsePathSimple()
	if err != nil {
		return nil, err
	}

	return &ast.RemoveStmt{
		Path: path,
		Loc:  ast.Loc{Start: start, End: path.Loc.End},
	}, nil
}

// parseSetCommand parses:
//
//	SET pathSimple = expr (, pathSimple = expr)*
//
// Grammar: setCommand (PartiQLParser.g4 lines 185-186).
func (p *Parser) parseSetCommand() ([]*ast.SetAssignment, error) {
	if _, err := p.expect(tokSET); err != nil {
		return nil, err
	}

	var assignments []*ast.SetAssignment
	for {
		a, err := p.parseSetAssignment()
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	return assignments, nil
}

// parseSetAssignment parses:
//
//	pathSimple = expr
//
// Grammar: setAssignment (PartiQLParser.g4 line 188-189).
func (p *Parser) parseSetAssignment() (*ast.SetAssignment, error) {
	target, err := p.parsePathSimple()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokEQ); err != nil {
		return nil, err
	}
	value, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	return &ast.SetAssignment{
		Target: target,
		Value:  value,
		Loc:    ast.Loc{Start: target.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseOnConflict parses:
//
//	ON CONFLICT WHERE expr DO NOTHING          (legacy form)
//	ON CONFLICT [conflictTarget] conflictAction (new form)
//
// Grammar: onConflictClause (PartiQLParser.g4 lines 139-142).
func (p *Parser) parseOnConflict() (*ast.OnConflict, error) {
	start := p.cur.Loc.Start
	p.advance() // consume ON
	if _, err := p.expect(tokCONFLICT); err != nil {
		return nil, err
	}

	// Legacy form: ON CONFLICT WHERE expr DO NOTHING
	if p.cur.Type == tokWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExprTop()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokDO); err != nil {
			return nil, err
		}
		if _, err := p.expect(tokNOTHING); err != nil {
			return nil, err
		}
		return &ast.OnConflict{
			Action: ast.OnConflictDoNothing,
			Where:  where,
			Loc:    ast.Loc{Start: start, End: p.prev.Loc.End},
		}, nil
	}

	// New form: ON CONFLICT [conflictTarget] conflictAction
	var target *ast.OnConflictTarget
	if p.cur.Type == tokPAREN_LEFT {
		// Conflict target: ( symbolPrimitive, ... )
		target, err := p.parseOnConflictTargetCols()
		if err != nil {
			return nil, err
		}
		_ = target
		// Re-assign properly
		action, err := p.parseConflictAction()
		if err != nil {
			return nil, err
		}
		return &ast.OnConflict{
			Target: target,
			Action: action,
			Loc:    ast.Loc{Start: start, End: p.prev.Loc.End},
		}, nil
	}
	if p.cur.Type == tokON {
		// ON CONSTRAINT constraintName
		target, err := p.parseOnConflictTargetConstraint()
		if err != nil {
			return nil, err
		}
		_ = target
		action, err := p.parseConflictAction()
		if err != nil {
			return nil, err
		}
		return &ast.OnConflict{
			Target: target,
			Action: action,
			Loc:    ast.Loc{Start: start, End: p.prev.Loc.End},
		}, nil
	}

	// No conflict target, just the action.
	action, err := p.parseConflictAction()
	if err != nil {
		return nil, err
	}
	return &ast.OnConflict{
		Target: target,
		Action: action,
		Loc:    ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}

// parseOnConflictTargetCols parses:
//
//	( symbolPrimitive, symbolPrimitive, ... )
func (p *Parser) parseOnConflictTargetCols() (*ast.OnConflictTarget, error) {
	start := p.cur.Loc.Start
	p.advance() // consume (

	var cols []*ast.VarRef
	for {
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		cols = append(cols, &ast.VarRef{
			Name:          name,
			CaseSensitive: caseSensitive,
			Loc:           nameLoc,
		})
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	return &ast.OnConflictTarget{
		Cols: cols,
		Loc:  ast.Loc{Start: start, End: rp.Loc.End},
	}, nil
}

// parseOnConflictTargetConstraint parses:
//
//	ON CONSTRAINT symbolPrimitive
func (p *Parser) parseOnConflictTargetConstraint() (*ast.OnConflictTarget, error) {
	start := p.cur.Loc.Start
	p.advance() // consume ON
	if _, err := p.expect(tokCONSTRAINT); err != nil {
		return nil, err
	}
	name, _, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	return &ast.OnConflictTarget{
		ConstraintName: name,
		Loc:            ast.Loc{Start: start, End: nameLoc.End},
	}, nil
}

// parseConflictAction parses:
//
//	DO NOTHING
//	DO REPLACE EXCLUDED
//	DO UPDATE EXCLUDED
func (p *Parser) parseConflictAction() (ast.OnConflictAction, error) {
	if _, err := p.expect(tokDO); err != nil {
		return ast.OnConflictInvalid, err
	}
	switch p.cur.Type {
	case tokNOTHING:
		p.advance()
		return ast.OnConflictDoNothing, nil
	case tokREPLACE:
		p.advance() // consume REPLACE
		if _, err := p.expect(tokEXCLUDED); err != nil {
			return ast.OnConflictInvalid, err
		}
		return ast.OnConflictDoReplaceExcluded, nil
	case tokUPDATE:
		p.advance() // consume UPDATE
		if _, err := p.expect(tokEXCLUDED); err != nil {
			return ast.OnConflictInvalid, err
		}
		return ast.OnConflictDoUpdateExcluded, nil
	default:
		return ast.OnConflictInvalid, &ParseError{
			Message: fmt.Sprintf("expected NOTHING, REPLACE, or UPDATE after DO, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
}

// parseReturningClause parses:
//
//	RETURNING returningColumn (, returningColumn)*
//
// Grammar: returningClause (PartiQLParser.g4 lines 194-195).
func (p *Parser) parseReturningClause() (*ast.ReturningClause, error) {
	start := p.cur.Loc.Start
	p.advance() // consume RETURNING

	var items []*ast.ReturningItem
	for {
		item, err := p.parseReturningItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		if p.cur.Type != tokCOMMA {
			break
		}
		p.advance() // consume ,
	}
	return &ast.ReturningClause{
		Items: items,
		Loc:   ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}

// parseReturningItem parses:
//
//	(MODIFIED|ALL) (OLD|NEW) (* | expr)
//
// Grammar: returningColumn (PartiQLParser.g4 lines 198-200).
func (p *Parser) parseReturningItem() (*ast.ReturningItem, error) {
	start := p.cur.Loc.Start

	var status ast.ReturningStatus
	switch p.cur.Type {
	case tokMODIFIED:
		status = ast.ReturningStatusModified
		p.advance()
	case tokALL:
		status = ast.ReturningStatusAll
		p.advance()
	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected MODIFIED or ALL in RETURNING clause, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}

	var mapping ast.ReturningMapping
	switch p.cur.Type {
	case tokOLD:
		mapping = ast.ReturningMappingOld
		p.advance()
	case tokNEW:
		mapping = ast.ReturningMappingNew
		p.advance()
	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected OLD or NEW in RETURNING clause, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}

	if p.cur.Type == tokASTERISK {
		end := p.cur.Loc.End
		p.advance()
		return &ast.ReturningItem{
			Status:  status,
			Mapping: mapping,
			Star:    true,
			Loc:     ast.Loc{Start: start, End: end},
		}, nil
	}

	expr, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	return &ast.ReturningItem{
		Status:  status,
		Mapping: mapping,
		Expr:    expr,
		Loc:     ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}

// parseFromClauseSimple parses:
//
//	FROM pathSimple [AS alias] [AT alias] [BY alias]
//	FROM pathSimple symbolPrimitive                   (implicit alias)
//
// Grammar: fromClauseSimple (PartiQLParser.g4 lines 202-205).
func (p *Parser) parseFromClauseSimple() (ast.TableExpr, error) {
	if _, err := p.expect(tokFROM); err != nil {
		return nil, err
	}

	path, err := p.parsePathSimple()
	if err != nil {
		return nil, err
	}

	// Check for aliases. The explicit form uses AS/AT/BY keywords;
	// the implicit form has a bare symbolPrimitive as alias.
	var asAlias, atAlias, byAlias *string
	end := path.Loc.End

	if p.cur.Type == tokAS {
		p.advance() // consume AS
		a, _, aLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		asAlias = &a
		end = aLoc.End
	} else if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		// Implicit alias form: FROM pathSimple symbolPrimitive
		// But only if the next token is an identifier, not a keyword
		// that could start a WHERE or RETURNING clause.
		// We check it is not a contextual keyword.
		if !isClauseKeyword(p.cur.Type) {
			a, _, aLoc, err := p.parseSymbolPrimitive()
			if err != nil {
				return nil, err
			}
			asAlias = &a
			end = aLoc.End
			// Implicit alias form does not support AT/BY.
			return &ast.AliasedSource{
				Source: path,
				As:     asAlias,
				Loc:    ast.Loc{Start: path.Loc.Start, End: end},
			}, nil
		}
	}

	if p.cur.Type == tokAT {
		p.advance() // consume AT
		a, _, aLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		atAlias = &a
		end = aLoc.End
	}
	if p.cur.Type == tokBY {
		p.advance() // consume BY
		a, _, aLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		byAlias = &a
		end = aLoc.End
	}

	if asAlias != nil || atAlias != nil || byAlias != nil {
		return &ast.AliasedSource{
			Source: path,
			As:     asAlias,
			At:     atAlias,
			By:     byAlias,
			Loc:    ast.Loc{Start: path.Loc.Start, End: end},
		}, nil
	}
	return path, nil
}

// isClauseKeyword returns true if the token is a keyword that starts a
// clause (WHERE, RETURNING, SET, etc.) and thus cannot be an implicit alias.
func isClauseKeyword(tokType int) bool {
	switch tokType {
	case tokWHERE, tokRETURNING, tokSET, tokON, tokORDER, tokGROUP,
		tokHAVING, tokLIMIT, tokOFFSET, tokUNION, tokINTERSECT, tokEXCEPT:
		return true
	}
	return false
}

package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/update.html
//
// Single-table:
//
//	UPDATE [LOW_PRIORITY | HIGH_PRIORITY | DELAYED] [IGNORE] table_reference
//	    SET assignment_list
//	    [WHERE where_condition]
//	    [ORDER BY ...]
//	    [LIMIT row_count]
//
// Multi-table:
//
//	UPDATE [LOW_PRIORITY | HIGH_PRIORITY | DELAYED] [IGNORE] table_references
//	    SET assignment_list
//	    [WHERE where_condition]
func (p *Parser) parseUpdateStmt() (*nodes.UpdateStmt, error) {
	start := p.pos()
	p.advance() // consume UPDATE

	// Completion: after UPDATE keyword, offer table_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		p.addRuleCandidate("database_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.UpdateStmt{Loc: nodes.Loc{Start: start}}

	// UpdateStmtNoWith (parser.y:13626, v8.5.5) reads PriorityOpt, then IGNORE.
	stmt.Priority = p.parsePriorityOpt()
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	// Parse table references
	tables, err := p.parseTableReferenceList()
	if err != nil {
		return nil, err
	}
	stmt.Tables = tables

	// SET clause (required)
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}

	setList, err := p.parseAssignmentList()
	if err != nil {
		return nil, err
	}
	stmt.SetList = setList

	// Completion: after SET clause, offer WHERE/ORDER/LIMIT.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwWHERE, kwORDER, kwLIMIT} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// WHERE clause (optional)
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && where == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Where = where
	}

	// ORDER BY and LIMIT are the single-table form's alone: UpdateStmtNoWith
	// reads them after one TableRef (which may be a join), while the
	// comma-separated TableRefs form has neither (parser.y:13626, :13654,
	// v8.5.5). A comma list is len(Tables) > 1. TiDB v8.5.5 answers 1064 for
	// `UPDATE t1, t2 SET ... LIMIT 1`.
	multiTable := len(stmt.Tables) > 1

	// ORDER BY clause (optional, single-table only)
	if p.cur.Type == kwORDER {
		if multiTable {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// LIMIT clause (optional, single-table only)
	if p.cur.Type == kwLIMIT {
		if multiTable {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDeleteStmt parses a DELETE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/delete.html
//
// Single-table:
//
//	DELETE [LOW_PRIORITY | HIGH_PRIORITY | DELAYED] [QUICK] [IGNORE] FROM tbl_name
//	    [WHERE where_condition]
//	    [ORDER BY ...]
//	    [LIMIT row_count]
//
// Multi-table (syntax 1):
//
//	DELETE [LOW_PRIORITY | HIGH_PRIORITY | DELAYED] [QUICK] [IGNORE]
//	    tbl_name[.*] [, tbl_name[.*]] ...
//	    FROM table_references
//	    [WHERE where_condition]
//
// Multi-table (syntax 2):
//
//	DELETE [LOW_PRIORITY | HIGH_PRIORITY | DELAYED] [QUICK] [IGNORE]
//	    FROM tbl_name[.*] [, tbl_name[.*]] ...
//	    USING table_references
//	    [WHERE where_condition]
func (p *Parser) parseDeleteStmt() (*nodes.DeleteStmt, error) {
	return p.parseDeleteStmtIn(deleteAnywhere)
}

// deletePlace says which DELETE forms a position takes.
type deletePlace int

const (
	// deleteAnywhere takes every form, as a statement of its own does.
	deleteAnywhere deletePlace = iota
	// deleteNoUsing leaves out the forms with a USING clause, which is
	// what a SQL plan binding takes: TiDB binds DeleteWithoutUsingStmt,
	// and in `DROP BINDING FOR DELETE FROM t1 USING t1, t2` the USING is
	// the binding's, so v8.5.5 fails at `t1, t2` rather than reading a
	// multiple-table DELETE. The FROM-first form then names one table:
	// `DROP BINDING FOR DELETE FROM t1, t2` is a 1064 at the comma, as
	// `DELETE FROM t1, t2` is on its own.
	deleteNoUsing
)

// parseDeleteStmtNoUsing parses the DELETE forms a SQL plan binding takes.
func (p *Parser) parseDeleteStmtNoUsing() (*nodes.DeleteStmt, error) {
	return p.parseDeleteStmtIn(deleteNoUsing)
}

func (p *Parser) parseDeleteStmtIn(place deletePlace) (*nodes.DeleteStmt, error) {
	start := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{Loc: nodes.Loc{Start: start}}
	multiTable := false

	// DeleteWithoutUsingStmt and DeleteWithUsingStmt (parser.y:5231 and :5281,
	// v8.5.5) read PriorityOpt, then QUICK, then IGNORE.
	stmt.Priority = p.parsePriorityOpt()
	if _, ok := p.match(kwQUICK); ok {
		stmt.Quick = true
	}
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	if _, ok := p.match(kwFROM); ok {
		// Could be single-table: DELETE FROM tbl WHERE ...
		// or multi-table syntax 2: DELETE FROM t1, t2 USING ...
		tables, sawWildcard, err := p.parseDeleteTableListIn(place)
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables

		// Single-table: optional [AS] alias and PARTITION
		if len(tables) == 1 {
			if tblRef, ok := tables[0].(*nodes.TableRef); ok {
				// Optional [AS] alias
				if _, ok := p.match(kwAS); ok {
					alias, _, err := p.parseIdent()
					if err != nil {
						return nil, err
					}
					tblRef.Alias = alias
					tblRef.Loc.End = p.prev.End
				} else if p.isIdentToken() {
					// TableAsNameOpt reads an Identifier, which is any
					// unreserved keyword as well as a plain identifier, so
					// `DELETE FROM t config` names an alias `config`.
					// Testing for tokIDENT alone let any word registered as
					// an unreserved keyword fall through unconsumed.
					alias, _, err := p.parseIdent()
					if err != nil {
						return nil, err
					}
					tblRef.Alias = alias
					tblRef.Loc.End = p.prev.End
				}
				// Optional PARTITION (p0, p1, ...)
				if p.cur.Type == kwPARTITION {
					p.advance()
					parts, err := p.parseParenIdentList()
					if err != nil {
						return nil, err
					}
					tblRef.Partitions = parts
					tblRef.Loc.End = p.prev.End
				}
			}
		}

		// Completion: after DELETE FROM table, offer WHERE/ORDER/LIMIT/USING.
		p.checkCursor()
		if p.collectMode() {
			for _, tok := range []int{kwWHERE, kwORDER, kwLIMIT, kwUSING} {
				p.addTokenCandidate(tok)
			}
			return nil, &ParseError{Message: "collecting"}
		}

		// Check for USING (multi-table syntax 2). A binding's FOR or USING
		// position takes no such DELETE, and the USING there belongs to the
		// binding, so this leaves the token where it is -- the place is
		// tested before the match, which would consume it either way.
		if place == deleteAnywhere {
			if _, ok := p.match(kwUSING); ok {
				using, err := p.parseTableReferenceList()
				if err != nil {
					return nil, err
				}
				stmt.Using = using
				multiTable = true
			}
		}
		// Without USING this is the single-table form, DELETE FROM TableName
		// (parser.y:5231): one plain table, no `.*`. `DELETE FROM t1, t2` and
		// `DELETE FROM t1.*` are 1064s on v8.5.5; the `.*` and the comma list
		// belong to the multi-table forms, which carry USING or name the
		// targets before FROM.
		if !multiTable && (len(stmt.Tables) != 1 || sawWildcard) {
			return nil, p.syntaxErrorAtCur()
		}
	} else {
		// Multi-table syntax 1: DELETE t1, t2 FROM table_references ...
		tables, _, err := p.parseDeleteTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables

		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}

		using, err := p.parseTableReferenceList()
		if err != nil {
			return nil, err
		}
		stmt.Using = using
		multiTable = true
	}

	// WHERE clause (optional)
	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && where == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Where = where
	}

	// ORDER BY and LIMIT are the single-table form's alone (parser.y:5231); the
	// multi-table forms (targets before FROM, or FROM ... USING) have neither.
	// TiDB v8.5.5 answers 1064 for `DELETE t1 FROM t1, t2 ORDER BY t1.a`.
	if p.cur.Type == kwORDER {
		if multiTable {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		orderBy, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	if p.cur.Type == kwLIMIT {
		if multiTable {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAssignmentList parses a comma-separated list of column=value assignments.
//
//	assignment_list:
//	    col_name = value [, col_name = value] ...
func (p *Parser) parseAssignmentList() ([]*nodes.Assignment, error) {
	var list []*nodes.Assignment

	for {
		asgn, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}
		list = append(list, asgn)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseAssignment parses a single col_name = value assignment.
func (p *Parser) parseAssignment() (*nodes.Assignment, error) {
	// Completion: at the start of an assignment, offer columnref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

	start := p.pos()

	col, err := p.parseNamedColumn()
	if err != nil {
		return nil, err
	}

	if err := p.expectEqual(); err != nil {
		return nil, err
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Column: col,
		Value:  val,
	}, nil
}

// parseNamedColumn parses the column an assignment or an INSERT column list
// names. parseColumnRef also reads a user variable, which is no column
// (simple_ident_nospvar), and t.*, which names none; `@a` in backquotes is a
// column like any other.
func (p *Parser) parseNamedColumn() (*nodes.ColumnRef, error) {
	if p.isUserVariable() {
		return nil, p.syntaxErrorAtCur()
	}
	ref, err := p.parseColumnRef()
	if err != nil {
		return nil, err
	}
	if ref.Star {
		return nil, &ParseError{Message: `syntax error at or near "*"`, Position: ref.Loc.End - 1}
	}
	return ref, nil
}

// expectEqual consumes the = or := between a column and its value in an
// assignment list.
//
// Ref: TiDB parser.y (v8.5.5) line 3415: Assignment is ColumnName
// EqOrAssignmentEq ExprOrDefault. LOAD DATA and IMPORT INTO take = alone
// (LoadDataSetItem, line 15069), so their SET lists do not come here.
func (p *Parser) expectEqual() error {
	if p.cur.Type == tokAssign {
		p.advance()
		return nil
	}
	_, err := p.expect('=')
	return err
}

// parseDeleteTableList parses a comma-separated list of table names
// for DELETE statements (with optional .* suffix).
func (p *Parser) parseDeleteTableList() ([]nodes.TableExpr, bool, error) {
	return p.parseDeleteTableListIn(deleteAnywhere)
}

// parseDeleteTableListIn parses the list, stopping at the first name where
// the place takes only one.
func (p *Parser) parseDeleteTableListIn(place deletePlace) ([]nodes.TableExpr, bool, error) {
	var tables []nodes.TableExpr
	sawWildcard := false

	for {
		// Completion: at the start of each table name, offer table_ref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, false, &ParseError{Message: "collecting"}
		}

		ref, wild, err := p.parseDeleteTableName()
		if err != nil {
			return nil, false, err
		}
		sawWildcard = sawWildcard || wild

		tables = append(tables, ref)
		if p.cur.Type != ',' {
			break
		}
		if place == deleteNoUsing {
			// The FROM-first form of the DELETE a binding takes names one
			// table, so the comma is where v8.5.5 fails.
			return nil, false, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	return tables, sawWildcard, nil
}

// parseDeleteTableName parses a single table name for DELETE,
// handling the optional .* suffix and schema.table.* form.
func (p *Parser) parseDeleteTableName() (*nodes.TableRef, bool, error) {
	start := p.pos()

	// `*.t1`, TiDB's wildcard database, which reaches DELETE's table list
	// through the table-name rule as it reaches every other one. It takes
	// no `.*` suffix: `DELETE FROM *.t1.*` is a 1064.
	if p.matchWildcardDatabase() {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, false, err
		}
		return &nodes.TableRef{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Schema: wildcardDatabase,
			Name:   name,
		}, false, nil
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, false, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Check for dot: could be schema.table or name.*
	if p.cur.Type == '.' {
		next := p.peekNext()
		if next.Type == '*' {
			// name.* -- skip the .* suffix, report it
			p.advance() // consume '.'
			p.advance() // consume '*'
			ref.Loc.End = p.prev.End
			return ref, true, nil
		}
		// schema.table or schema.table.*
		p.advance() // consume '.'
		// Completion: after "db.", offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, false, &ParseError{Message: "collecting"}
		}
		if err := refuseTableNameDatabase(name, start); err != nil {
			return nil, false, err
		}
		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, false, err
		}
		ref.Schema = name
		ref.Name = name2

		// Check for .* after schema.table
		if p.cur.Type == '.' {
			next := p.peekNext()
			if next.Type == '*' {
				p.advance() // consume '.'
				p.advance() // consume '*'
				ref.Loc.End = p.prev.End
				return ref, true, nil
			}
		}
	}

	ref.Loc.End = p.prev.End
	return ref, false, nil
}

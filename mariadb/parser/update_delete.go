package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/update.html
//
// Single-table:
//
//	UPDATE [LOW_PRIORITY] [IGNORE] table_reference
//	    SET assignment_list
//	    [WHERE where_condition]
//	    [ORDER BY ...]
//	    [LIMIT row_count]
//
// Multi-table:
//
//	UPDATE [LOW_PRIORITY] [IGNORE] table_references
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

	// Optional modifiers
	if _, ok := p.match(kwLOW_PRIORITY); ok {
		stmt.LowPriority = true
	}
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	// Parse table references
	tables, err := p.parseTableReferenceList()
	if err != nil {
		return nil, err
	}
	if err := refuseFromDuplicates(tables); err != nil {
		return nil, err
	}
	stmt.Tables = tables

	// FOR PORTION OF <period> FROM x TO y (application-time). Valid only on a
	// single bare base table (no join, no alias before the clause); for UPDATE
	// the alias, if any, follows the clause:
	//   UPDATE t FOR PORTION OF p FROM x TO y [AS] alias SET ...
	if p.cur.Type == kwFOR {
		// A FOR SYSTEM_TIME-bearing target cannot also take FOR PORTION OF
		// (MariaDB 11.8.8 rejects the combination); leave the trailing FOR to
		// the expect(kwSET) below, which reports 1064.
		if tref := singleBaseTable(stmt.Tables); tref != nil && tref.Alias == "" && tref.SystemTime == nil {
			fp, err := p.parseForPortionOf()
			if err != nil {
				return nil, err
			}
			stmt.ForPortionOf = fp
			if err := p.parseOptionalTableAlias(tref); err != nil {
				return nil, err
			}
		}
		// Otherwise the misplaced FOR is rejected by expect(kwSET) below (1064).
	}

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

	// ORDER BY clause (optional, single-table only)
	if p.cur.Type == kwORDER {
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
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseForPortionOf parses the application-time clause
//
//	FOR PORTION OF <period> FROM <expr> TO <expr>
//
// on a single-table UPDATE / DELETE. FOR is the current token. The period name
// follows the identifier rule (non-reserved keywords allowed); SYSTEM_TIME is
// rejected (1064), matching MariaDB. PORTION is a non-reserved keyword.
func (p *Parser) parseForPortionOf() (*nodes.ForPortionOf, error) {
	start := p.pos()
	p.advance() // FOR
	if _, err := p.expect(kwPORTION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwOF); err != nil {
		return nil, err
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(name, "SYSTEM_TIME") {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && from == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	to, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.collectMode() && to == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.ForPortionOf{
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
		PeriodName: name,
		From:       from,
		To:         to,
	}, nil
}

// singleBaseTable returns the table reference when tables is a single bare base
// table (not a join or comma-list), else nil.
func singleBaseTable(tables []nodes.TableExpr) *nodes.TableRef {
	if len(tables) != 1 {
		return nil
	}
	tref, _ := tables[0].(*nodes.TableRef)
	return tref
}

// parseOptionalTableAlias consumes an optional [AS] alias and records it on tref.
func (p *Parser) parseOptionalTableAlias(tref *nodes.TableRef) error {
	_, hasAS := p.match(kwAS)
	if !hasAS && !p.isIdentToken() {
		return nil
	}
	// WINDOW names a table anywhere but as its alias: `UPDATE t window SET a = 1` is a 1064.
	if p.cur.Type == kwWINDOW {
		return p.syntaxErrorAtCur()
	}
	alias, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	tref.Alias = alias
	tref.Loc.End = p.prev.End
	return nil
}

// parseDeleteStmt parses a DELETE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/delete.html
//
// Single-table:
//
//	DELETE [LOW_PRIORITY] [QUICK] [IGNORE] FROM tbl_name
//	    [WHERE where_condition]
//	    [ORDER BY ...]
//	    [LIMIT row_count]
//
// Multi-table (syntax 1):
//
//	DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
//	    tbl_name[.*] [, tbl_name[.*]] ...
//	    FROM table_references
//	    [WHERE where_condition]
//
// Multi-table (syntax 2):
//
//	DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
//	    FROM tbl_name[.*] [, tbl_name[.*]] ...
//	    USING table_references
//	    [WHERE where_condition]
func (p *Parser) parseDeleteStmt() (*nodes.DeleteStmt, error) {
	start := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{Loc: nodes.Loc{Start: start}}

	// Optional modifiers
	if _, ok := p.match(kwLOW_PRIORITY); ok {
		stmt.LowPriority = true
	}
	if _, ok := p.match(kwQUICK); ok {
		stmt.Quick = true
	}
	if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	if _, ok := p.match(kwFROM); ok {
		// Could be single-table: DELETE FROM tbl WHERE ...
		// or multi-table syntax 2: DELETE FROM t1, t2 USING ...
		tables, err := p.parseDeleteTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables

		// Single-table: optional [AS] alias and PARTITION
		if len(tables) == 1 {
			if tblRef, ok := tables[0].(*nodes.TableRef); ok {
				// Optional [AS] alias
				if _, ok := p.match(kwAS); ok {
					if p.cur.Type == kwWINDOW {
						return nil, p.syntaxErrorAtCur()
					}
					alias, _, err := p.parseIdent()
					if err != nil {
						return nil, err
					}
					tblRef.Alias = alias
					tblRef.Loc.End = p.prev.End
				} else if p.isIdentToken() && p.cur.Type != kwWINDOW {
					// table_alias reads an ident, which is any unreserved
					// keyword as well as a plain identifier, so
					// `DELETE FROM t wait WHERE wait.a = 1` names an alias
					// `wait` (container-verified, as are the `locked` and
					// `nowait` forms). Testing for tokIDENT alone left any
					// word registered as an unreserved keyword unconsumed.
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

		// Check for USING (multi-table syntax 2)
		if _, ok := p.match(kwUSING); ok {
			using, err := p.parseTableReferenceList()
			if err != nil {
				return nil, err
			}
			stmt.Using = using
		}
	} else {
		// Multi-table syntax 1: DELETE t1, t2 FROM table_references ...
		tables, err := p.parseDeleteTableList()
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
	}

	if err := refuseFromDuplicates(stmt.Using); err != nil {
		return nil, err
	}
	if err := refuseDeleteTargets(stmt.Tables, stmt.Using); err != nil {
		return nil, err
	}

	// FOR PORTION OF <period> FROM x TO y (application-time, single-table only).
	if len(stmt.Tables) == 1 && len(stmt.Using) == 0 && p.cur.Type == kwFOR {
		fp, err := p.parseForPortionOf()
		if err != nil {
			return nil, err
		}
		stmt.ForPortionOf = fp
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

	// ORDER BY clause (optional, single-table only)
	if p.cur.Type == kwORDER {
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
	if _, ok := p.match(kwLIMIT); ok {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	// RETURNING (MariaDB) — single-table DELETE only. Multi-table delete (either
	// syntax) has >1 table or USING; MariaDB rejects RETURNING there (1064), so
	// the clause is not reachable on those paths.
	if len(stmt.Tables) == 1 && stmt.Using == nil && p.cur.Type == kwRETURNING {
		ret, err := p.parseReturningClause()
		if err != nil {
			return nil, err
		}
		stmt.Returning = ret
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

	if ignore, ok := p.parseIgnoreValue(); ok {
		return &nodes.Assignment{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Column: col,
			Value:  ignore,
		}, nil
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

// parseNamedColumn parses the column an assignment names. parseColumnRef
// also reads a user variable, which is no column (simple_ident_nospvar), and
// t.*, which names none; `@a` in backquotes is a column like any other.
func (p *Parser) parseNamedColumn() (*nodes.ColumnRef, error) {
	ref, err := p.parseColumnNotVariable()
	if err != nil {
		return nil, err
	}
	if ref.Star {
		return nil, &ParseError{Message: `syntax error at or near "*"`, Position: ref.Loc.End - 1}
	}
	return ref, nil
}

// parseColumnNotVariable parses a column reference that is not a user variable.
func (p *Parser) parseColumnNotVariable() (*nodes.ColumnRef, error) {
	if p.isUserVariable() {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseColumnRef()
}

// matchEqual consumes the = or := of an option, where the grammar makes it
// optional, and reports whether it found one. A position the grammar writes
// with a literal = takes no :=, such as a view's ALGORITHM, DEFINER, FORMAT, and
// every CHANGE REPLICATION option, so those read = alone.
//
// Ref: MariaDB 11.8 sql/sql_yacc.yy: opt_equal is empty or equal, and opt_with and
// opt_by read opt_equal too.
func (p *Parser) matchEqual() bool {
	if p.cur.Type == '=' || p.cur.Type == tokAssign {
		p.advance()
		return true
	}
	return false
}

// expectEqual consumes the = or := between a column and its value in an
// assignment list.
//
// Ref: MariaDB 11.8 sql/sql_yacc.yy: update_elem, insert_update_elem, ident_eq_value
// and load_data_set_elem all read a column, equal, and a value, and equal is
// '=' or SET_VAR.
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
func (p *Parser) parseDeleteTableList() ([]nodes.TableExpr, error) {
	var tables []nodes.TableExpr

	for {
		// Completion: at the start of each table name, offer table_ref candidates.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			p.addRuleCandidate("database_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		ref, err := p.parseDeleteTableName()
		if err != nil {
			return nil, err
		}

		tables = append(tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return tables, nil
}

// parseDeleteTableName parses a single table name for DELETE,
// handling the optional .* suffix and schema.table.* form.
func (p *Parser) parseDeleteTableName() (*nodes.TableRef, error) {
	start := p.pos()
	dotted := p.peekNext().Type == '.'
	ref, err := p.parseUncheckedDeleteTableName()
	if err != nil {
		return nil, err
	}
	// name.* names no database, where `name`.t does, and an empty one is no Schema.
	qualified := ref.Schema != "" || (dotted && !strings.HasSuffix(p.inputText(start, p.prev.End), ".*"))
	return ref, refuseTableRefName(ref, qualified, start)
}

func (p *Parser) parseUncheckedDeleteTableName() (*nodes.TableRef, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	ref := &nodes.TableRef{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Check for dot: could be schema.table or name.*
	if p.cur.Type == '.' {
		next := p.peekNext()
		if next.Type == '*' {
			// name.* -- skip the .* suffix
			p.advance() // consume '.'
			p.advance() // consume '*'
			ref.Loc.End = p.prev.End
			return ref, nil
		}
		// schema.table or schema.table.*
		p.advance() // consume '.'
		// Completion: after "db.", offer table_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		name2, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ref.Schema = name
		ref.Name = name2

		// Check for .* after schema.table
		if p.cur.Type == '.' {
			next := p.peekNext()
			if next.Type == '*' {
				p.advance() // consume '.'
				p.advance() // consume '*'
			}
		}
	}

	ref.Loc.End = p.prev.End
	return ref, nil
}

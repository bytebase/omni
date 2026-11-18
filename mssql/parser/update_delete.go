// Package parser - update_delete.go implements T-SQL UPDATE and DELETE statement parsing.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// BNF: mssql/parser/bnf/update-transact-sql.bnf
//
//	[ WITH <common_table_expression> [...n] ]
//	UPDATE
//	    [ TOP ( expression ) [ PERCENT ] ]
//	    { { table_alias | <object> | rowset_function_limited
//	         [ WITH ( <Table_Hint_Limited> [ ...n ] ) ]
//	      }
//	      | @table_variable
//	    }
//	    SET
//	        { column_name = { expression | DEFAULT | NULL }
//	          | { udt_column_name.{ { property_name = expression
//	                                | field_name = expression }
//	                                | method_name ( argument [ ,...n ] )
//	                              }
//	          }
//	          | column_name { .WRITE ( expression , @Offset , @Length ) }
//	          | @variable = expression
//	          | @variable = column = expression
//	          | column_name { += | -= | *= | /= | %= | &= | ^= | |= } expression
//	          | @variable { += | -= | *= | /= | %= | &= | ^= | |= } expression
//	          | @variable = column { += | -= | *= | /= | %= | &= | ^= | |= } expression
//	        } [ ,...n ]
//
//	    [ <OUTPUT Clause> ]
//	    [ FROM{ <table_source> } [ ,...n ] ]
//	    [ WHERE { <search_condition>
//	            | { [ CURRENT OF
//	                  { { [ GLOBAL ] cursor_name }
//	                      | cursor_variable_name
//	                  }
//	                ]
//	              }
//	            }
//	    ]
//	    [ OPTION ( <query_hint> [ ,...n ] ) ]
//	[ ; ]
func (p *Parser) parseUpdateStmt() (*nodes.UpdateStmt, error) {
	loc := p.pos()
	p.advance() // consume UPDATE

	stmt := &nodes.UpdateStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		top, err := p.parseTopClause()
		if err != nil {
			return nil, err
		}
		stmt.Top = top
	}

	// Completion: after UPDATE → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	// Table name, rowset function or @table_variable
	rel, err := p.parseDmlTarget()
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, p.newParseError(p.cur.Loc, "expected table name after UPDATE")
	}
	stmt.Relation = rel

	// SET clause, which updateSpecification (TSql160.g:19808) requires.
	if _, err := p.expect(kwSET); err != nil {
		if p.collectMode() {
			p.addTokenCandidate(kwSET)
			return nil, errCollecting
		}
		return nil, err
	}
	// Completion: after SET → columnref
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	setList, err := p.parseSetClauseList()
	if err != nil {
		return nil, err
	}
	if setList == nil || len(setList.Items) == 0 {
		return nil, p.newParseError(p.cur.Loc, "expected assignment after SET")
	}
	stmt.SetClause = setList

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		oc, err := p.parseOutputClause()
		if err != nil {
			return nil, err
		}
		stmt.OutputClause = oc
	}

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.FromClause = from
	}

	// WHERE clause (includes CURRENT OF cursor support)
	if _, ok := p.match(kwWHERE); ok {
		// Completion: after WHERE in UPDATE → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		wc, err := p.parseWhereClauseBody()
		if err != nil {
			return nil, err
		}
		stmt.WhereClause = wc
	}

	// OPTION clause
	if p.cur.Type == kwOPTION {
		oc, err := p.parseOptionClause()
		if err != nil {
			return nil, err
		}
		stmt.OptionClause = oc
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDeleteStmt parses a DELETE statement.
//
// BNF: mssql/parser/bnf/delete-transact-sql.bnf
//
//	[ WITH <common_table_expression> [ ,...n ] ]
//	DELETE
//	    [ TOP ( expression ) [ PERCENT ] ]
//	    [ FROM ]
//	    { { table_alias
//	      | <object>
//	      | rowset_function_limited
//	      [ WITH ( table_hint_limited [ ...n ] ) ] }
//	      | @table_variable
//	    }
//	    [ <OUTPUT Clause> ]
//	    [ FROM table_source [ ,...n ] ]
//	    [ WHERE { <search_condition>
//	            | { [ CURRENT OF
//	                   { { [ GLOBAL ] cursor_name }
//	                       | cursor_variable_name
//	                   }
//	                ]
//	              }
//	            }
//	    ]
//	    [ OPTION ( <Query Hint> [ ,...n ] ) ]
//	[; ]
func (p *Parser) parseDeleteStmt() (*nodes.DeleteStmt, error) {
	loc := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		top, err := p.parseTopClause()
		if err != nil {
			return nil, err
		}
		stmt.Top = top
	}

	// Optional FROM before table name
	p.match(kwFROM)

	// Completion: after DELETE [FROM] → table_ref
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, errCollecting
	}

	// Table name, rowset function or @table_variable
	rel, err := p.parseDmlTarget()
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, p.newParseError(p.cur.Loc, "expected table name after DELETE")
	}
	stmt.Relation = rel

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		oc, err := p.parseOutputClause()
		if err != nil {
			return nil, err
		}
		stmt.OutputClause = oc
	}

	// FROM clause (second FROM for join)
	if _, ok := p.match(kwFROM); ok {
		from, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.FromClause = from
	}

	// WHERE clause (includes CURRENT OF cursor support)
	if _, ok := p.match(kwWHERE); ok {
		// Completion: after WHERE in DELETE → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("func_name")
			return nil, errCollecting
		}
		wc, err := p.parseWhereClauseBody()
		if err != nil {
			return nil, err
		}
		stmt.WhereClause = wc
	}

	// OPTION clause
	if p.cur.Type == kwOPTION {
		oc, err := p.parseOptionClause()
		if err != nil {
			return nil, err
		}
		stmt.OptionClause = oc
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDmlTarget parses the target of an UPDATE or DELETE:
//
//	dmlTarget (TSql160.g:20084)
//	    : schemaObjectDmlTarget | openRowset | variableDmlTarget
//
// SQL Server takes the rowset functions here as rowset functions in FROM
// (measured): OPENQUERY ( linked_server , 'query' ), OPENDATASOURCE ( ... )
// . table, and a name with an argument list, as in DELETE a(1). It takes no
// alias after OPENQUERY (Msg 102), no WITH ( hint ) after the call (Msg 319),
// and no OPENROWSET or OPENXML (Msg 7222 and 156), so a name without an
// argument list keeps the hints a named table takes.
func (p *Parser) parseDmlTarget() (nodes.TableExpr, error) {
	if p.cur.Type == tokVARIABLE {
		return p.parseVariableDmlTarget(), nil
	}
	switch p.cur.Type {
	case kwOPENDATASOURCE:
		return p.parseRowsetFunction()
	case kwOPENQUERY:
		ref, err := p.parseRowsetFunction()
		if err != nil {
			return nil, err
		}
		if a, ok := ref.(*nodes.AliasedTableRef); ok && a.Alias != "" {
			return nil, p.newParseError(a.Loc.End-1, "OPENQUERY takes no alias as a DML target")
		}
		return ref, nil
	}

	rel, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, nil
	}

	// A name with an argument list is a function: schemaObjectFunctionDmlTarget's
	// parenthesizedOptExpressionWithDefaultList.
	if p.cur.Type == '(' {
		args, err := p.parseDmlTargetArgs()
		if err != nil {
			return nil, err
		}
		fc := &nodes.FuncCallExpr{
			Name: rel,
			Args: args,
			Loc:  nodes.Loc{Start: rel.Loc.Start, End: p.prevEnd()},
		}
		return fc, nil
	}

	// Optional WITH ( <Table_Hint_Limited> ) on target — only allowed on named tables.
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		hints, err := p.parseDmlTargetHints()
		if err != nil {
			return nil, err
		}
		rel.Hints = hints
	}
	return rel, nil
}

// parseDmlTargetArg parses one argument of a DML target's call. SQL Server
// takes a constant, a variable, DEFAULT or a negated number there and no
// expression of its own: `DELETE a(zzz)` and `DELETE a(1+1)` are Msg 102,
// `DELETE a(1)`, a(@x), a('a'), a(NULL), a(-1) and a(DEFAULT) are not
// (measured).
func (p *Parser) parseDmlTargetArg() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST, tokMONEY, kwNULL, kwDEFAULT,
		tokVARIABLE, tokSYSVARIABLE:
		return p.parsePrimary()
	case '-':
		loc := p.pos()
		p.advance() // consume '-'
		switch p.cur.Type {
		case tokICONST, tokFCONST, tokMONEY:
			if execIsBinary(p.cur) {
				return nil, p.syntaxErrorAtCur()
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &nodes.UnaryExpr{
			Op:      nodes.UnaryMinus,
			Operand: operand,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	return nil, nil
}

// parseDmlTargetArgs parses the argument list a DML target's function call
// takes: parenthesizedOptExpressionWithDefaultList (TSql160.g), whose
// expressions may not carry DEFAULT there (Msg 156, measured).
func (p *Parser) parseDmlTargetArgs() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, ok := p.match(')'); ok {
		return nil, nil
	}
	var items []nodes.Node
	for {
		arg, err := p.parseDmlTargetArg()
		if err != nil {
			return nil, err
		}
		if arg == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, arg)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseWhereClauseBody parses the body of a WHERE clause, handling both
// normal search conditions and CURRENT OF cursor_name.
//
//	WHERE { <search_condition>
//	      | CURRENT OF { { [ GLOBAL ] cursor_name } | cursor_variable_name }
//	      }
func (p *Parser) parseWhereClauseBody() (nodes.ExprNode, error) {
	if p.cur.Type == kwCURRENT {
		next := p.peekNext()
		if next.Type == kwOF {
			loc := p.pos()
			p.advance() // consume CURRENT
			p.advance() // consume OF

			// CURRENT OF [ GLOBAL ] cursor_name | @cursor_variable; the
			// cursor is required (cursorId, TSql160.g:14652).
			global := p.matchIdentCI("GLOBAL")
			var cursorName string
			if p.cur.Type == tokVARIABLE && !global {
				cursorName = p.cur.Str
				p.advance()
			} else if name, ok := p.parseIdentifier(); ok {
				cursorName = name
			} else {
				return nil, p.syntaxErrorAtCur()
			}
			return &nodes.CurrentOfExpr{
				CursorName: cursorName,
				Global:     global,
				Loc:        nodes.Loc{Start: loc, End: p.prevEnd()},
			}, nil
		}
	}
	p.enterSearchCondition()
	defer p.leaveSearchCondition()
	return p.parseExpr()
}

// parseDmlTargetHints parses WITH ( table_hint_limited [ ...n ] ) after the
// target of an INSERT, UPDATE or DELETE. The list may not be empty
// (tableHints, TSql160.g:21027), and SQL Server refuses on a target what it
// takes in a FROM clause: an index hint (Msg 1069), FORCESEEK (Msg 10724),
// FORCESCAN (Msg 10745), NOLOCK and READUNCOMMITTED (Msg 1065).
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/hints-transact-sql-table
func (p *Parser) parseDmlTargetHints() (*nodes.List, error) {
	hints, err := p.parseTableHints()
	if err != nil {
		return nil, err
	}
	if hints == nil {
		return nil, p.newParseError(p.prev.Loc, "syntax error at or near \")\"")
	}
	for _, item := range hints.Items {
		hint, ok := item.(*nodes.TableHint)
		if !ok {
			continue
		}
		switch strings.ToUpper(hint.Name) {
		case "INDEX", "FORCESEEK", "FORCESCAN", "NOLOCK", "READUNCOMMITTED":
			return nil, p.newParseError(hint.Loc.Start, fmt.Sprintf("the %s hint is not allowed on the target of a data modification statement", strings.ToUpper(hint.Name)))
		}
	}
	return hints, nil
}

// isCompoundAssign returns the operator string if the current token is a compound assignment operator,
// or empty string if not.
func (p *Parser) isCompoundAssign() string {
	switch p.cur.Type {
	case tokPLUSEQUAL:
		return "+="
	case tokMINUSEQUAL:
		return "-="
	case tokMULEQUAL:
		return "*="
	case tokDIVEQUAL:
		return "/="
	case tokMODEQUAL:
		return "%="
	case tokANDEQUAL:
		return "&="
	case tokOREQUAL:
		return "|="
	case tokXOREQUAL:
		return "^="
	default:
		return ""
	}
}

// parseSetClauseList parses a comma-separated list of SET assignments.
//
//	set_clause_list = set_clause { ',' set_clause }
func (p *Parser) parseSetClauseList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		// Completion: at start of SET clause item → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, errCollecting
		}
		item, err := p.parseSetClause()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
		if _, ok := p.match(','); !ok {
			break
		}
		// Completion: after comma in SET clause → columnref
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, errCollecting
		}
	}
	return &nodes.List{Items: items}, nil
}

// parseSetClause parses a single SET assignment.
//
//	column_name = { expression | DEFAULT | NULL }
//	column_name { .WRITE ( expression , @Offset , @Length ) }
//	@variable = expression
//	@variable = column = expression
//	column_name { += | -= | *= | /= | %= | &= | ^= | |= } expression
//	@variable { += | -= | *= | /= | %= | &= | ^= | |= } expression
//	@variable = column { += | -= | *= | /= | %= | &= | ^= | |= } expression
func (p *Parser) parseSetClause() (*nodes.SetExpr, error) {
	loc := p.pos()

	se := &nodes.SetExpr{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == tokVARIABLE {
		se.Variable = p.cur.Str
		p.advance()
	} else {
		// Parse column reference (not a full expression - just the column name)
		target, err := p.parseSetTarget()
		if err != nil {
			return nil, err
		}
		se.Column = target
		if se.Column == nil {
			return nil, nil
		}

		// Check for .WRITE(expression, @Offset, @Length) form
		if p.cur.Type == '.' && p.peekNext().Type == kwWRITE {
			next := p.peekNext()
			if next.Type == kwWRITE {
				p.advance() // consume .
				writeLoc := p.pos()
				p.advance() // consume WRITE
				if p.cur.Type == '(' {
					se.WriteMethod = true
					// Parse as a FuncCallExpr to store the args
					fc, err := p.parseFuncCall("WRITE", writeLoc)
					if err != nil {
						return nil, err
					}
					se.Value = fc
					se.Loc.End = p.prevEnd()
					return se, nil
				}
			}
		}
	}

	// Check for compound assignment operators (+=, -=, *=, /=, %=, &=, ^=, |=) or simple =
	if op := p.isCompoundAssign(); op != "" {
		se.Operator = op
		p.advance()
	} else if _, err := p.expect('='); err != nil {
		return nil, nil
	}
	if p.collectMode() {
		p.addExpressionCandidates()
		return nil, errCollecting
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.newParseError(p.cur.Loc, "expected expression after =")
	}
	se.Value = val
	se.Loc.End = p.prevEnd()
	return se, nil
}

// parseSetTarget parses the left side of a SET assignment (qualified column
// name). Bare known pseudo-columns are accepted like the engine does —
// mutability and context are binding-time errors there (UPDATE e SET
// $from_id = ... is Msg 271, SET $identity is Msg 8102), while unknown $foo
// is Msg 126. The table-qualified form (e.$from_id) is a syntax error in the
// engine (Msg 102) and stays rejected.
func (p *Parser) parseSetTarget() (*nodes.ColumnRef, error) {
	loc := p.pos()

	if p.cur.Type == tokPSEUDOCOL {
		if !knownPseudoColumns[strings.ToLower(p.cur.Str)] {
			return nil, p.newParseError(p.cur.Loc, fmt.Sprintf("invalid pseudocolumn %q", p.cur.Str))
		}
		ref := &nodes.ColumnRef{
			Column: p.cur.Str,
			Loc:    nodes.Loc{Start: loc, End: p.cur.End},
		}
		p.advance()
		return ref, nil
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, nil
	}

	ref := &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// Check for qualified name: table.column
	// But don't consume dot if next is .WRITE
	if p.cur.Type == '.' {
		next := p.peekNext()
		// Don't consume if it's .WRITE (handled separately)
		if next.Type == kwWRITE {
			return ref, nil
		}
		p.advance()
		col, ok := p.parseIdentifier()
		if ok {
			ref.Table = name
			ref.Column = col
		}
	}

	return ref, nil
}

// Package parser - update_delete.go implements T-SQL UPDATE and DELETE statement parsing.
package parser

import (
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
func (p *Parser) parseUpdateStmt() *nodes.UpdateStmt {
	loc := p.pos()
	p.advance() // consume UPDATE

	stmt := &nodes.UpdateStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		stmt.Top, _ = p.parseTopClause()
	}

	// Table name or @table_variable
	stmt.Relation, _ = p.parseTableRef()

	// Optional WITH ( <Table_Hint_Limited> ) on target
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		stmt.Relation.Hints, _ = p.parseTableHints()
	}

	// SET clause
	if _, err := p.expect(kwSET); err == nil {
		stmt.SetClause = p.parseSetClauseList()
	}

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		stmt.OutputClause = p.parseOutputClause()
	}

	// FROM clause
	if _, ok := p.match(kwFROM); ok {
		stmt.FromClause, _ = p.parseFromClause()
	}

	// WHERE clause (includes CURRENT OF cursor support)
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseWhereClauseBody()
	}

	// OPTION clause
	if p.cur.Type == kwOPTION {
		stmt.OptionClause, _ = p.parseOptionClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseDeleteStmt() *nodes.DeleteStmt {
	loc := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Optional TOP
	if p.cur.Type == kwTOP {
		stmt.Top, _ = p.parseTopClause()
	}

	// Optional FROM before table name
	p.match(kwFROM)

	// Table name or @table_variable
	stmt.Relation, _ = p.parseTableRef()

	// Optional WITH ( <Table_Hint_Limited> ) on target
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		stmt.Relation.Hints, _ = p.parseTableHints()
	}

	// OUTPUT clause
	if p.cur.Type == kwOUTPUT {
		stmt.OutputClause = p.parseOutputClause()
	}

	// FROM clause (second FROM for join)
	if _, ok := p.match(kwFROM); ok {
		stmt.FromClause, _ = p.parseFromClause()
	}

	// WHERE clause (includes CURRENT OF cursor support)
	if _, ok := p.match(kwWHERE); ok {
		stmt.WhereClause = p.parseWhereClauseBody()
	}

	// OPTION clause
	if p.cur.Type == kwOPTION {
		stmt.OptionClause, _ = p.parseOptionClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseWhereClauseBody parses the body of a WHERE clause, handling both
// normal search conditions and CURRENT OF cursor_name.
//
//	WHERE { <search_condition>
//	      | CURRENT OF { { [ GLOBAL ] cursor_name } | cursor_variable_name }
//	      }
func (p *Parser) parseWhereClauseBody() nodes.ExprNode {
	if p.cur.Type == kwCURRENT {
		next := p.peekNext()
		if next.Type == kwOF {
			loc := p.pos()
			p.advance() // consume CURRENT
			p.advance() // consume OF

			// CURRENT OF [ GLOBAL ] cursor_name | @cursor_variable
			global := p.matchIdentCI("GLOBAL")
			var cursorName string
			if p.cur.Type == tokVARIABLE {
				cursorName = p.cur.Str
				p.advance()
			} else if name, ok := p.parseIdentifier(); ok {
				cursorName = name
			}
			return &nodes.CurrentOfExpr{
				CursorName: cursorName,
				Global:     global,
				Loc:        nodes.Loc{Start: loc, End: p.pos()},
			}
		}
	}
	result, _ := p.parseExpr()
	return result
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
func (p *Parser) parseSetClauseList() *nodes.List {
	var items []nodes.Node
	for {
		item := p.parseSetClause()
		if item == nil {
			break
		}
		items = append(items, item)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: items}
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
func (p *Parser) parseSetClause() *nodes.SetExpr {
	loc := p.pos()

	se := &nodes.SetExpr{
		Loc: nodes.Loc{Start: loc},
	}

	if p.cur.Type == tokVARIABLE {
		se.Variable = p.cur.Str
		p.advance()
	} else {
		// Parse column reference (not a full expression - just the column name)
		se.Column = p.parseSetTarget()
		if se.Column == nil {
			return nil
		}

		// Check for .WRITE(expression, @Offset, @Length) form
		if p.cur.Type == '.' && p.peekNext().Type == tokIDENT {
			next := p.peekNext()
			if strings.EqualFold(next.Str, "WRITE") {
				p.advance() // consume .
				writeLoc := p.pos()
				p.advance() // consume WRITE
				if p.cur.Type == '(' {
					se.WriteMethod = true
					// Parse as a FuncCallExpr to store the args
					se.Value, _ = p.parseFuncCall("WRITE", writeLoc)
					se.Loc.End = p.pos()
					return se
				}
			}
		}
	}

	// Check for compound assignment operators (+=, -=, *=, /=, %=, &=, ^=, |=) or simple =
	if op := p.isCompoundAssign(); op != "" {
		se.Operator = op
		p.advance()
	} else if _, err := p.expect('='); err != nil {
		return nil
	}

	se.Value, _ = p.parseExpr()
	se.Loc.End = p.pos()
	return se
}

// parseSetTarget parses the left side of a SET assignment (qualified column name).
func (p *Parser) parseSetTarget() *nodes.ColumnRef {
	loc := p.pos()

	name, ok := p.parseIdentifier()
	if !ok {
		return nil
	}

	ref := &nodes.ColumnRef{
		Column: name,
		Loc:    nodes.Loc{Start: loc},
	}

	// Check for qualified name: table.column
	// But don't consume dot if next is .WRITE
	if p.cur.Type == '.' {
		next := p.peekNext()
		// Don't consume if it's .WRITE (handled separately)
		if next.Type == tokIDENT && strings.EqualFold(next.Str, "WRITE") {
			return ref
		}
		p.advance()
		col, ok := p.parseIdentifier()
		if ok {
			ref.Table = name
			ref.Column = col
		}
	}

	return ref
}

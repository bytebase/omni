package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseUpdateStmt parses an UPDATE statement.
//
// BNF: oracle/parser/bnf/UPDATE.bnf
//
//	UPDATE [ hint ]
//	    { { [ schema. ] { table | view | materialized_view }
//	        [ t_alias ]
//	        [ partition_extension_clause ]
//	        [ @ dblink ]
//	      }
//	    | ( subquery [ subquery_restriction_clause ] )
//	    | table_collection_expression
//	    }
//	    update_set_clause
//	    [ from_using_clause ]
//	    [ where_clause ]
//	    [ returning_clause ]
//	    [ error_logging_clause ] ;
//
//	partition_extension_clause:
//	    { PARTITION ( partition_name )
//	    | PARTITION FOR ( partition_key_value [, partition_key_value]... )
//	    | SUBPARTITION ( subpartition_name )
//	    | SUBPARTITION FOR ( subpartition_key_value [, subpartition_key_value]... )
//	    }
//
//	subquery_restriction_clause:
//	    WITH { READ ONLY | CHECK OPTION [ CONSTRAINT constraint_name ] }
//
//	table_collection_expression:
//	    TABLE ( collection_expression ) [ (+) ]
//
//	update_set_clause:
//	    SET
//	    { { column = { expr | ( subquery ) | DEFAULT }
//	      | ( column [, column]... ) = ( subquery )
//	      } [, { column = { expr | ( subquery ) | DEFAULT }
//	           | ( column [, column]... ) = ( subquery )
//	           } ]...
//	    | VALUE ( t_alias ) = { expr | ( subquery ) }
//	    }
//
//	from_using_clause:
//	    FROM { table_reference | join_clause }
//	         [, { table_reference | join_clause } ]...
//
//	where_clause:
//	    WHERE condition
//
//	order_by_clause:
//	    ORDER [ SIBLINGS ] BY
//	    { expr | position | c_alias }
//	    [ ASC | DESC ]
//	    [ NULLS FIRST | NULLS LAST ]
//	    [, { expr | position | c_alias }
//	       [ ASC | DESC ]
//	       [ NULLS FIRST | NULLS LAST ]
//	    ]...
//
//	returning_clause:
//	    { RETURN | RETURNING }
//	    { expr | OLD column | NEW column }
//	    [, { expr | OLD column | NEW column } ]...
//	    INTO data_item [, data_item]...
//
//	error_logging_clause:
//	    LOG ERRORS
//	    [ INTO [ schema. ] table_name ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
func (p *Parser) parseUpdateStmt() *nodes.UpdateStmt {
	start := p.pos()
	p.advance() // consume UPDATE

	stmt := &nodes.UpdateStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Hints
	if p.cur.Type == tokHINT {
		stmt.Hints = &nodes.List{}
		stmt.Hints.Items = append(stmt.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// Table name
	stmt.Table = p.parseObjectName()

	// Optional alias (before partition clause, matching BNF order for UPDATE)
	if p.cur.Type == kwAS {
		p.advance()
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// Partition extension clause
	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		stmt.PartitionExt = p.parsePartitionExtClause()
	}

	// @dblink
	if p.cur.Type == '@' {
		p.advance()
		stmt.Dblink = p.parseIdentifier()
	}

	// SET
	if p.cur.Type == kwSET {
		p.advance()
		stmt.SetClauses = p.parseSetClauses()
	}

	// FROM clause (Oracle 23c+)
	if p.cur.Type == kwFROM {
		p.advance()
		stmt.FromClause = p.parseFromList()
	}

	// WHERE
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.WhereClause = p.parseExpr()
	}

	// RETURNING / RETURN ... INTO ...
	if p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN {
		stmt.Returning = p.parseReturningClause()
	}

	// LOG ERRORS
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFromList parses a comma-separated list of table references for FROM clause.
func (p *Parser) parseFromList() *nodes.List {
	list := &nodes.List{}
	for {
		tr := p.parseTableRef()
		if tr == nil {
			break
		}
		list.Items = append(list.Items, tr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}

// parseSetClauses parses comma-separated SET clauses.
func (p *Parser) parseSetClauses() *nodes.List {
	list := &nodes.List{}

	for {
		sc := p.parseSetClause()
		if sc == nil {
			break
		}
		list.Items = append(list.Items, sc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list
}

// parseSetClause parses a single SET clause: col = expr | (col1,col2) = (subquery).
func (p *Parser) parseSetClause() *nodes.SetClause {
	start := p.pos()

	// Multi-column: (col1, col2) = (subquery)
	if p.cur.Type == '(' {
		p.advance()
		sc := &nodes.SetClause{
			Columns: &nodes.List{},
			Loc:     nodes.Loc{Start: start},
		}
		for {
			col := p.parseColumnRef()
			if col != nil {
				sc.Columns.Items = append(sc.Columns.Items, col)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		if p.cur.Type == '=' {
			p.advance()
		}
		sc.Value = p.parseExpr()
		sc.Loc.End = p.pos()
		return sc
	}

	// Single column: col = expr
	if !p.isIdentLike() {
		return nil
	}

	col := p.parseColumnRef()
	sc := &nodes.SetClause{
		Column: col,
		Loc:    nodes.Loc{Start: start},
	}

	if p.cur.Type == '=' {
		p.advance()
	}

	// DEFAULT keyword or expression
	if p.cur.Type == kwDEFAULT {
		sc.Value = &nodes.ColumnRef{
			Column: "DEFAULT",
			Loc:    nodes.Loc{Start: p.pos(), End: p.pos()},
		}
		p.advance()
	} else {
		sc.Value = p.parseExpr()
	}
	sc.Loc.End = p.pos()
	return sc
}

// parseDeleteStmt parses a DELETE statement.
//
// BNF: oracle/parser/bnf/DELETE.bnf
//
//	DELETE [ hint ]
//	    [ FROM ]
//	    dml_table_expression_clause
//	    [ from_using_clause ]
//	    [ where_clause ]
//	    [ returning_clause ]
//	    [ error_logging_clause ] ;
//
//	dml_table_expression_clause:
//	    { [ schema. ] { table | view | materialized_view }
//	        [ partition_extension_clause ]
//	        [ @ dblink ]
//	        [ t_alias ]
//	    | ( subquery [ subquery_restriction_clause ] ) [ t_alias ]
//	    | TABLE ( collection_expression ) [ (+) ] [ t_alias ]
//	    }
//
//	partition_extension_clause:
//	    { PARTITION ( partition_name )
//	    | SUBPARTITION ( subpartition_name )
//	    }
//
//	subquery_restriction_clause:
//	    { WITH READ ONLY
//	    | WITH CHECK OPTION [ CONSTRAINT constraint_name ]
//	    }
//
//	from_using_clause:
//	    FROM { table_reference
//	         | [ LATERAL ] inline_view [ join_clause ]
//	         | join_clause
//	         }
//
//	where_clause:
//	    WHERE condition
//
//	returning_clause:
//	    RETURNING expr [, expr ]...
//	    INTO data_item [, data_item ]...
//
//	error_logging_clause:
//	    LOG ERRORS [ INTO [ schema. ] table_name ]
//	    [ ( simple_expression ) ]
//	    [ REJECT LIMIT { integer | UNLIMITED } ]
func (p *Parser) parseDeleteStmt() *nodes.DeleteStmt {
	start := p.pos()
	p.advance() // consume DELETE

	stmt := &nodes.DeleteStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Hints
	if p.cur.Type == tokHINT {
		stmt.Hints = &nodes.List{}
		stmt.Hints.Items = append(stmt.Hints.Items, &nodes.Hint{
			Text: p.cur.Str,
			Loc:  nodes.Loc{Start: p.pos(), End: p.pos()},
		})
		p.advance()
	}

	// Optional FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

	// Table name
	stmt.Table = p.parseObjectName()

	// Partition extension clause
	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		stmt.PartitionExt = p.parsePartitionExtClause()
	}

	// @dblink
	if p.cur.Type == '@' {
		p.advance()
		stmt.Dblink = p.parseIdentifier()
	}

	// Optional alias
	if p.cur.Type == kwAS {
		p.advance()
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	} else if p.isTableAliasCandidate() {
		stmt.Alias = &nodes.Alias{Name: p.parseIdentifier()}
	}

	// WHERE
	if p.cur.Type == kwWHERE {
		p.advance()
		stmt.WhereClause = p.parseExpr()
	}

	// RETURNING / RETURN ... INTO ...
	if p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN {
		stmt.Returning = p.parseReturningClause()
	}

	// LOG ERRORS
	if p.cur.Type == kwLOG {
		stmt.ErrorLog = p.parseErrorLogClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

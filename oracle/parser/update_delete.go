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
func (p *Parser) parseUpdateStmt() (*nodes.UpdateStmt, error) {
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
			Loc:  nodes.Loc{Start: p.pos(), End: p.prev.End},
		})
		p.advance()
	}
	var parseErr1127 error

	// Table name
	stmt.Table, parseErr1127 = p.parseObjectName()
	if parseErr1127 !=

		// Optional alias (before partition clause, matching BNF order for UPDATE)
		nil {
		return nil, parseErr1127
	}

	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1128 error
		stmt.Alias, parseErr1128 = p.parseAlias()
		if parseErr1128 != nil {
			return nil, parseErr1128
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1129 error
		stmt.Alias, parseErr1129 = p.parseAlias()
		if parseErr1129 !=

			// Partition extension clause
			nil {
			return nil, parseErr1129
		}
	}

	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		var parseErr1130 error
		stmt.PartitionExt, parseErr1130 = p.parsePartitionExtClause()
		if parseErr1130 !=

			// @dblink
			nil {
			return nil, parseErr1130
		}
	}

	if p.cur.Type == '@' {
		p.advance()
		var parseErr1131 error
		stmt.Dblink, parseErr1131 = p.parseIdentifier()
		if parseErr1131 !=

			// SET
			nil {
			return nil, parseErr1131
		}
	}

	if p.cur.Type == kwSET {
		p.advance()
		var parseErr1132 error
		stmt.SetClauses, parseErr1132 = p.parseSetClauses()
		if parseErr1132 !=

			// FROM clause (Oracle 23c+)
			nil {
			return nil, parseErr1132
		}
	}

	if p.cur.Type == kwFROM {
		p.advance()
		var parseErr1133 error
		stmt.FromClause, parseErr1133 = p.parseFromList()
		if parseErr1133 !=

			// WHERE
			nil {
			return nil, parseErr1133
		}
	}

	if p.cur.Type == kwWHERE {
		p.advance()
		var parseErr1134 error
		stmt.WhereClause, parseErr1134 = p.parseExpr()
		if parseErr1134 !=

			// RETURNING / RETURN ... INTO ...
			nil {
			return nil, parseErr1134
		}
	}

	if p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN {
		var parseErr1135 error
		stmt.Returning, parseErr1135 = p.parseReturningClause()
		if parseErr1135 !=

			// LOG ERRORS
			nil {
			return nil, parseErr1135
		}
	}

	if p.cur.Type == kwLOG {
		var parseErr1136 error
		stmt.ErrorLog, parseErr1136 = p.parseErrorLogClause()
		if parseErr1136 != nil {
			return nil, parseErr1136
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFromList parses a comma-separated list of table references for FROM clause.
func (p *Parser) parseFromList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		tr, parseErr1137 := p.parseTableRef()
		if parseErr1137 != nil {
			return nil, parseErr1137
		}
		if tr == nil {
			break
		}
		list.Items = append(list.Items, tr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list, nil
}

// parseSetClauses parses comma-separated SET clauses.
func (p *Parser) parseSetClauses() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		sc, parseErr1138 := p.parseSetClause()
		if parseErr1138 != nil {
			return nil, parseErr1138
		}
		if sc == nil {
			break
		}
		list.Items = append(list.Items, sc)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	return list, nil
}

// parseSetClause parses a single SET clause: col = expr | (col1,col2) = (subquery).
func (p *Parser) parseSetClause() (*nodes.SetClause, error) {
	start := p.pos()

	// Multi-column: (col1, col2) = (subquery)
	if p.cur.Type == '(' {
		p.advance()
		sc := &nodes.SetClause{
			Columns: &nodes.List{},
			Loc:     nodes.Loc{Start: start},
		}
		for {
			col, parseErr1139 := p.parseColumnRef()
			if parseErr1139 != nil {
				return nil, parseErr1139
			}
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
		var parseErr1140 error
		sc.Value, parseErr1140 = p.parseExpr()
		if parseErr1140 != nil {
			return nil, parseErr1140
		}
		sc.Loc.End = p.prev.End
		return sc, nil
	}

	// Single column: col = expr
	if !p.isIdentLike() {
		return nil, nil
	}

	col, parseErr1141 := p.parseColumnRef()
	if parseErr1141 != nil {
		return nil, parseErr1141
	}
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
			Loc:    nodes.Loc{Start: p.pos(), End: p.prev.End},
		}
		p.advance()
	} else {
		var parseErr1142 error
		sc.Value, parseErr1142 = p.parseExpr()
		if parseErr1142 != nil {
			return nil, parseErr1142
		}
	}
	sc.Loc.End = p.prev.End
	return sc, nil
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
func (p *Parser) parseDeleteStmt() (*nodes.DeleteStmt, error) {
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
			Loc:  nodes.Loc{Start: p.pos(), End: p.prev.End},
		})
		p.advance()
	}

	// Optional FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}
	var parseErr1143 error

	// Table name
	stmt.Table, parseErr1143 = p.parseObjectName()
	if parseErr1143 !=

		// Partition extension clause
		nil {
		return nil, parseErr1143
	}

	if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
		var parseErr1144 error
		stmt.PartitionExt, parseErr1144 = p.parsePartitionExtClause()
		if parseErr1144 !=

			// @dblink
			nil {
			return nil, parseErr1144
		}
	}

	if p.cur.Type == '@' {
		p.advance()
		var parseErr1145 error
		stmt.Dblink, parseErr1145 = p.parseIdentifier()
		if parseErr1145 !=

			// Optional alias
			nil {
			return nil, parseErr1145
		}
	}

	if p.cur.Type == kwAS {
		p.advance()
		var parseErr1146 error
		stmt.Alias, parseErr1146 = p.parseAlias()
		if parseErr1146 != nil {
			return nil, parseErr1146
		}
	} else if p.isTableAliasCandidate() {
		var parseErr1147 error
		stmt.Alias, parseErr1147 = p.parseAlias()
		if parseErr1147 !=

			// WHERE
			nil {
			return nil, parseErr1147
		}
	}

	if p.cur.Type == kwWHERE {
		p.advance()
		var parseErr1148 error
		stmt.WhereClause, parseErr1148 = p.parseExpr()
		if parseErr1148 !=

			// RETURNING / RETURN ... INTO ...
			nil {
			return nil, parseErr1148
		}
	}

	if p.cur.Type == kwRETURNING || p.cur.Type == kwRETURN {
		var parseErr1149 error
		stmt.Returning, parseErr1149 = p.parseReturningClause()
		if parseErr1149 !=

			// LOG ERRORS
			nil {
			return nil, parseErr1149
		}
	}

	if p.cur.Type == kwLOG {
		var parseErr1150 error
		stmt.ErrorLog, parseErr1150 = p.parseErrorLogClause()
		if parseErr1150 != nil {
			return nil, parseErr1150
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCommentStmt parses a COMMENT ON statement.
//
// BNF: oracle/parser/bnf/COMMENT.bnf
//
//	COMMENT ON
//	    { AUDIT POLICY policy_name
//	    | COLUMN [ schema. ] { table | view | materialized_view } . column
//	    | EDITION edition_name
//	    | INDEXTYPE [ schema. ] indextype
//	    | MATERIALIZED VIEW [ schema. ] materialized_view
//	    | MINING MODEL [ schema. ] mining_model
//	    | OPERATOR [ schema. ] operator
//	    | PROPERTY GRAPH [ schema. ] property_graph
//	    | TABLE [ schema. ] { table | view }
//	    }
//	    IS string ;
func (p *Parser) parseCommentStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume COMMENT

	// Expect ON
	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	stmt := &nodes.CommentStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr196 error

	// Parse object type
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	case kwCOLUMN:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
		// For COLUMN, parse a dotted name: [schema.]table.column
		name, parseErr194 := p.parseObjectName()
		if parseErr194 != nil {
			return nil, parseErr194
		}

		// The last part of the name is the column name.
		// parseObjectName gives us schema.name, but for COMMENT ON COLUMN
		// we may have schema.table.column or table.column.
		// We need to check if there was an additional ".column" part.
		// parseObjectName handles schema.name, so let's see if there's another dot.
		if p.cur.Type == '.' {
			// schema.table.column case: name has schema=schema, name=table
			// and we need to parse the column part
			p.advance()
			var parseErr195 error
			stmt.Column, parseErr195 = p.parseIdentifier()
			if parseErr195 != nil {
				return nil, parseErr195
			}
			stmt.Object = name
		} else {
			// table.column case: name has schema="", name=table or schema=table, name=column
			if name.Schema != "" {
				// schema.table was parsed, but that was actually table.column
				stmt.Object = &nodes.ObjectName{
					Name: name.Schema,
					Loc:  name.Loc,
				}
				stmt.Column = name.Name
			} else {
				// Just "name" with no dots - odd, but handle it
				stmt.Object = name
			}
		}
		goto parseIs
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
	case kwMATERIALIZED:
		stmt.ObjectType = nodes.OBJECT_MATERIALIZED_VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			p.advance()
		}
	default:
		// Could be other object types (EDITION, INDEXTYPE, etc.) - treat as TABLE
		stmt.ObjectType = nodes.OBJECT_TABLE
		if p.isIdentLike() {
			p.advance()
		}
	}

	// Parse object name for non-COLUMN cases
	stmt.Object, parseErr196 = p.parseObjectName()
	if parseErr196 != nil {
		return nil, parseErr196
	}
	if stmt.Object == nil || stmt.Object.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

parseIs:

	if p.cur.Type != kwIS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != tokSCONST && p.cur.Type != kwNULL {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == tokSCONST {
		stmt.Comment = p.cur.Str
	}
	p.advance()

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

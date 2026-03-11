package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCommentStmt parses a COMMENT ON statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/COMMENT.html
//
//	COMMENT ON { TABLE schema.table | COLUMN schema.table.column |
//	             INDEX schema.index | MATERIALIZED VIEW schema.mv |
//	             ... }
//	IS 'comment_text'
func (p *Parser) parseCommentStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume COMMENT

	// Expect ON
	if p.cur.Type == kwON {
		p.advance()
	}

	stmt := &nodes.CommentStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse object type
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	case kwCOLUMN:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
		// For COLUMN, parse a dotted name: [schema.]table.column
		name := p.parseObjectName()
		// The last part of the name is the column name.
		// parseObjectName gives us schema.name, but for COMMENT ON COLUMN
		// we may have schema.table.column or table.column.
		// We need to check if there was an additional ".column" part.
		// parseObjectName handles schema.name, so let's see if there's another dot.
		if p.cur.Type == '.' {
			// schema.table.column case: name has schema=schema, name=table
			// and we need to parse the column part
			p.advance()
			stmt.Column = p.parseIdentifier()
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
	stmt.Object = p.parseObjectName()

parseIs:
	// IS 'comment_text'
	if p.cur.Type == kwIS {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.Comment = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

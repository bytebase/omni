// Package parser - alter_table.go implements T-SQL ALTER TABLE statement parsing.
package parser

import (
	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseAlterTableStmt parses an ALTER TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-table-transact-sql
//
//	ALTER TABLE name { ADD col_def | DROP COLUMN name | ALTER COLUMN ... | ADD CONSTRAINT ... | DROP CONSTRAINT ... }
func (p *Parser) parseAlterTableStmt() *nodes.AlterTableStmt {
	loc := p.pos()

	stmt := &nodes.AlterTableStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Table name
	stmt.Name = p.parseTableRef()

	// Parse actions
	var actions []nodes.Node
	action := p.parseAlterTableAction()
	if action != nil {
		actions = append(actions, action)
	}
	stmt.Actions = &nodes.List{Items: actions}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterTableAction parses a single ALTER TABLE action.
func (p *Parser) parseAlterTableAction() *nodes.AlterTableAction {
	loc := p.pos()
	action := &nodes.AlterTableAction{
		Loc: nodes.Loc{Start: loc},
	}

	switch p.cur.Type {
	case kwADD:
		p.advance()
		// ADD CONSTRAINT or ADD column
		if p.cur.Type == kwCONSTRAINT {
			action.Type = nodes.ATAddConstraint
			action.Constraint = p.parseTableConstraint()
		} else {
			action.Type = nodes.ATAddColumn
			action.Column = p.parseColumnDef()
		}
	case kwDROP:
		p.advance()
		if p.cur.Type == kwCOLUMN {
			p.advance()
			action.Type = nodes.ATDropColumn
			name, _ := p.parseIdentifier()
			action.ColName = name
		} else if p.cur.Type == kwCONSTRAINT {
			p.advance()
			action.Type = nodes.ATDropConstraint
			action.Constraint = &nodes.ConstraintDef{
				Loc: nodes.Loc{Start: p.pos()},
			}
			name, _ := p.parseIdentifier()
			action.Constraint.Name = name
			action.Constraint.Loc.End = p.pos()
		} else {
			// Could be DROP column without COLUMN keyword
			action.Type = nodes.ATDropColumn
			name, _ := p.parseIdentifier()
			action.ColName = name
		}
	case kwALTER:
		p.advance()
		if p.cur.Type == kwCOLUMN {
			p.advance()
		}
		action.Type = nodes.ATAlterColumn
		name, _ := p.parseIdentifier()
		action.ColName = name
		action.DataType = p.parseDataType()
		// NULL / NOT NULL
		if p.cur.Type == kwNULL {
			p.advance()
		} else if p.cur.Type == kwNOT {
			next := p.peekNext()
			if next.Type == kwNULL {
				p.advance() // NOT
				p.advance() // NULL
			}
		}
	default:
		return nil
	}

	action.Loc.End = p.pos()
	return action
}

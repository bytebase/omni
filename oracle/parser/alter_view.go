package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseAlterViewStmt parses an ALTER VIEW statement. It is called with VIEW
// as the current token.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-VIEW.html
//
//	ALTER VIEW [ IF EXISTS ] [ schema. ] view
//	    { COMPILE
//	    | { READ ONLY | READ WRITE }
//	    | ADD out_of_line_constraint
//	    | MODIFY { CONSTRAINT constraint_name | PRIMARY KEY | UNIQUE ( column [, column ]... ) }
//	        { RELY | NORELY }
//	    | DROP { CONSTRAINT constraint_name | PRIMARY KEY | UNIQUE ( column [, column ]... ) }
//	    | { EDITIONABLE | NONEDITIONABLE }
//	    | annotations_clause
//	    }
//
//	annotations_clause:
//	    ANNOTATIONS ( annotation_item [, annotation_item ]... )
//
//	annotation_item:
//	    [ ADD [ IF NOT EXISTS | OR REPLACE ] | DROP [ IF EXISTS ] | REPLACE ]
//	    annotation_name [ 'annotation_value' ]
//
// The page names constraints only by name and writes the annotations clause
// as ANNOTATIONS ( { ADD OR REPLACE | DROP } ANNOTATIONS ( ... ) ); Oracle
// 23ai takes PRIMARY KEY and UNIQUE as well, and reads the annotations as
// CREATE TABLE does (ORA-11548 for a reserved name, ORA-11549 for a value
// that is not a string, ORA-11551 for an empty list). MODIFY needs RELY or
// NORELY (ORA-00921); a constraint name is not reserved (ORA-02250); the
// view takes no database link (ORA-00922 at the '@').
func (p *Parser) parseAlterViewStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume VIEW
	stmt := &nodes.AlterViewStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance() // consume IF
		p.advance() // consume EXISTS
	}
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.isKeywordStr("COMPILE"), p.isKeywordStr("RECOMPILE"):
		stmt.Action = "COMPILE"
		p.advance()

	case p.cur.Type == kwADD:
		stmt.Action = "ADD_CONSTRAINT"
		p.advance() // consume ADD
		tc, err := p.parseTableConstraintBody(false)
		if err != nil {
			return nil, err
		}
		if err := p.parseViewConstraintState(); err != nil {
			return nil, err
		}
		// DISABLE / ENABLE and NOVALIDATE / VALIDATE after the constraint.
		for p.cur.Type == kwDISABLE || p.cur.Type == kwENABLE || p.isKeywordStr("NOVALIDATE") || p.isKeywordStr("VALIDATE") {
			p.advance()
		}
		tc.Loc.End = p.prev.End
		stmt.Constraint = tc

	case p.cur.Type == kwMODIFY:
		p.advance() // consume MODIFY
		stmt.Action = "MODIFY_CONSTRAINT"
		if err := p.parseAlterViewConstraintRef(stmt); err != nil {
			return nil, err
		}
		switch {
		case p.cur.Type == kwRELY:
			stmt.Rely = true
		case p.isKeywordStr("NORELY"):
			stmt.NoRely = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume RELY or NORELY

	case p.cur.Type == kwDROP:
		p.advance() // consume DROP
		stmt.Action = "DROP_CONSTRAINT"
		if err := p.parseAlterViewConstraintRef(stmt); err != nil {
			return nil, err
		}

	case p.isKeywordStr("ANNOTATIONS"):
		stmt.Action = "ANNOTATIONS"
		if stmt.Annotations, err = p.parseAlterViewAnnotations(); err != nil {
			return nil, err
		}

	case p.cur.Type == kwREAD:
		p.advance() // consume READ
		switch {
		case p.isKeywordStr("ONLY"):
			stmt.Action = "READ_ONLY"
		case p.cur.Type == kwWRITE:
			stmt.Action = "READ_WRITE"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

	case p.isKeywordStr("EDITIONABLE"), p.isKeywordStr("NONEDITIONABLE"):
		stmt.Action = p.cur.Str
		p.advance()

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterViewConstraintRef parses the constraint that ALTER VIEW MODIFY or
// DROP names: CONSTRAINT constraint_name, PRIMARY KEY, or UNIQUE ( column
// [, column ]... ). A constraint named by its kind changes stmt's action to
// MODIFY_PRIMARY_KEY, MODIFY_UNIQUE, DROP_PRIMARY_KEY, or DROP_UNIQUE.
func (p *Parser) parseAlterViewConstraintRef(stmt *nodes.AlterViewStmt) error {
	prefix := "DROP_"
	if stmt.Action == "MODIFY_CONSTRAINT" {
		prefix = "MODIFY_"
	}
	switch {
	case p.cur.Type == kwCONSTRAINT:
		p.advance() // consume CONSTRAINT
		if !p.isCreateTableNameStart() {
			return p.syntaxErrorAtCur()
		}
		stmt.ConstraintName = p.cur.Str
		p.advance()
	case p.cur.Type == kwPRIMARY:
		p.advance() // consume PRIMARY
		if err := p.expectCreateTableToken(kwKEY); err != nil {
			return err
		}
		stmt.Action = prefix + "PRIMARY_KEY"
	case p.cur.Type == kwUNIQUE:
		p.advance() // consume UNIQUE
		if err := p.parseCreateTableColumnNames(); err != nil {
			return err
		}
		stmt.Action = prefix + "UNIQUE"
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterViewAnnotations parses the annotations_clause of ALTER VIEW, from
// its ANNOTATIONS, into one option per annotation: the key is the name, the
// value the annotation's value, and the items hold the operation (ADD, ADD IF
// NOT EXISTS, ADD OR REPLACE, DROP, DROP IF EXISTS, or REPLACE) when one is
// written.
func (p *Parser) parseAlterViewAnnotations() (*nodes.List, error) {
	p.advance() // consume ANNOTATIONS
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	result := &nodes.List{}
	for {
		itemStart := p.pos()
		opStart := p.pos()
		var op string
		switch {
		case p.cur.Type == kwADD:
			op = "ADD"
			p.advance() // consume ADD
			switch {
			case p.cur.Type == kwIF:
				if p.peekNext().Type != kwNOT || p.peekAhead(2)[1].Type != kwEXISTS {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume IF
				p.advance() // consume NOT
				p.advance() // consume EXISTS
				op = "ADD IF NOT EXISTS"
			case p.cur.Type == kwOR:
				if p.peekNext().Type != kwREPLACE {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume OR
				p.advance() // consume REPLACE
				op = "ADD OR REPLACE"
			}
		case p.cur.Type == kwDROP:
			op = "DROP"
			p.advance() // consume DROP
			if p.cur.Type == kwIF {
				if p.peekNext().Type != kwEXISTS {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume IF
				p.advance() // consume EXISTS
				op = "DROP IF EXISTS"
			}
		case p.cur.Type == kwREPLACE:
			op = "REPLACE"
			p.advance() // consume REPLACE
		}
		opEnd := p.prev.End
		if !p.isCreateTableNameStart() {
			return nil, p.syntaxErrorAtCur()
		}
		opt := &nodes.DDLOption{Key: p.cur.Str}
		p.advance()
		if p.cur.Type == tokSCONST || p.cur.Type == tokNCHARLIT {
			opt.Value = p.cur.Str
			p.advance()
		}
		if op != "" {
			opt.Items = &nodes.List{Items: []nodes.Node{
				&nodes.DDLOption{Key: op, Loc: nodes.Loc{Start: opStart, End: opEnd}},
			}}
		}
		opt.Loc = nodes.Loc{Start: itemStart, End: p.prev.End}
		result.Items = append(result.Items, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	return result, nil
}

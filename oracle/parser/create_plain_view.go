package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// finishCreatePlainViewStmt finishes a CREATE VIEW statement, not a
// materialized one, after VIEW.
//
// BNF: oracle/parser/bnf/CREATE-VIEW.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-VIEW.html
//
//	VIEW [ IF NOT EXISTS ] [ schema. ] view
//	[ SHARING = { METADATA | DATA | EXTENDED DATA | NONE } ]
//	[ ( { alias [ VISIBLE | INVISIBLE ] [ inline_constraint ]... } [, ...]
//	    [, out_of_line_constraint ]... ) ]
//	[ object_view_clause ]
//	[ annotations_clause ]
//	[ DEFAULT COLLATION collation_name ]
//	[ BEQUEATH { CURRENT_USER | DEFINER } ]
//	AS subquery
//	[ subquery_restriction_clause ]
//	[ CONTAINER_MAP | CONTAINERS_DEFAULT ] ;
//
//	object_view_clause:
//	    OF [ schema. ] type_name
//	    { WITH OBJECT IDENTIFIER { DEFAULT | ( attribute [, attribute ]... ) }
//	    | UNDER [ schema. ] superview }
//	subquery_restriction_clause:
//	    { WITH READ ONLY | WITH CHECK OPTION } [ CONSTRAINT constraint_name ]
//
// A view's constraints are declarative: Oracle 23ai takes UNIQUE, PRIMARY
// KEY and REFERENCES, each with [ RELY | NORELY ] DISABLE [ NOVALIDATE ],
// and refuses NOT NULL, NULL, CHECK and ON DELETE on one (ORA-00922 at
// them). An XMLType view (OF XMLTYPE) is not read.
func (p *Parser) finishCreatePlainViewStmt(stmt *nodes.CreateViewStmt) (*nodes.CreateViewStmt, error) {
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		ifTok := p.cur
		p.advance() // IF
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance() // EXISTS
		stmt.IfNotExists = true
	}

	var err error
	stmt.Name, err = p.parseReservedCheckedObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.isKeywordStr("SHARING") {
		p.advance() // SHARING
		if err := p.expectCreateTableToken('='); err != nil {
			return nil, err
		}
		switch {
		case p.isKeywordStr("METADATA"), p.isKeywordStr("DATA"), p.isKeywordStr("NONE"):
			stmt.Sharing = p.cur.Str
			p.advance()
		case p.isKeywordStr("EXTENDED"):
			p.advance()
			if err := p.expectCreateTableWord("DATA"); err != nil {
				return nil, err
			}
			stmt.Sharing = "EXTENDED DATA"
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type == '(' {
		if err := p.parsePlainViewAliases(stmt); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwOF {
		if err := p.parseObjectViewClause(); err != nil {
			return nil, err
		}
	}

	if p.isKeywordStr("ANNOTATIONS") {
		annotations, err := p.parseDomainAnnotations(false)
		if err != nil {
			return nil, err
		}
		if annotations != nil {
			seenAnns := map[string]bool{}
			for _, item := range annotations.Items {
				opt, ok := item.(*nodes.DDLOption)
				if !ok {
					continue
				}
				if seenAnns[opt.Key] {
					return nil, p.syntaxErrorAtCur()
				}
				seenAnns[opt.Key] = true
			}
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance() // DEFAULT
		if err := p.expectCreateTableWord("COLLATION"); err != nil {
			return nil, err
		}
		if p.isStatementEnd() {
			// Oracle checks DEFAULT COLLATION against MAX_STRING_SIZE
			// first (ORA-43929 on this server), and compat.tsv records the
			// statement cut short here as accepted; omni keeps its earlier
			// reading of it.
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
		stmt.DefaultCollation = p.cur.Str
		if err := p.parseCreateTableName(); err != nil {
			return nil, err
		}
	}

	if p.isKeywordStr("BEQUEATH") {
		p.advance()
		switch {
		case p.isKeywordStr("CURRENT_USER"), p.isKeywordStr("DEFINER"):
			stmt.Bequeath = p.cur.Str
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	if err := p.expectCreateTableToken(kwAS); err != nil {
		return nil, err
	}
	if stmt.Query, err = p.parseSelectStmt(); err != nil {
		return nil, err
	}
	if !selectStmtHasTargets(stmt.Query) {
		return nil, p.syntaxErrorAtCur()
	}
	if stmt.Editioning {
		if err := p.checkEditioningViewTargets(stmt.Query); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwWITH {
		p.advance() // WITH
		switch {
		case p.cur.Type == kwCHECK:
			p.advance()
			if err := p.expectCreateTableToken(kwOPTION); err != nil {
				return nil, err
			}
			stmt.WithCheckOpt = true
		case p.cur.Type == kwREAD:
			p.advance()
			if err := p.expectCreateTableToken(kwONLY); err != nil {
				return nil, err
			}
			stmt.WithReadOnly = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwCONSTRAINTS {
			p.advance()
			stmt.ConstraintName = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return nil, err
			}
		}
	}

	switch {
	case p.isKeywordStr("CONTAINER_MAP"):
		stmt.ContainerMap = true
		p.advance()
	case p.isKeywordStr("CONTAINERS_DEFAULT"):
		stmt.ContainersDefault = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePlainViewAliases reads a view's ( alias ... ) list.
func (p *Parser) parsePlainViewAliases(stmt *nodes.CreateViewStmt) error {
	p.advance() // (
	stmt.Columns = &nodes.List{}
	seenCols := map[string]bool{}
	for {
		if p.isViewConstraintStart() || p.cur.Type == kwUNIQUE {
			if (p.cur.Type == kwCONSTRAINT || p.cur.Type == kwCONSTRAINTS) &&
				p.peekNext().Type != kwUNIQUE && p.peekNext().Type != kwPRIMARY && p.peekNext().Type != kwFOREIGN {
				ahead := p.peekAhead(3)
				if len(ahead) == 3 && ahead[2].Type == kwUNIQUE {
					p.advance() // CONSTRAINT
					p.advance() // name
					p.advance() // unexpected word
					p.advance() // UNIQUE
					return p.syntaxErrorAtCur()
				}
			}
			viewConstraint, err := p.parseTableConstraintBody(false)
			if err != nil {
				return err
			}
			if viewConstraint.Type == nodes.CONSTRAINT_CHECK {
				return p.syntaxErrorAtNode(viewConstraint)
			}
			if err := p.parseViewConstraintState(); err != nil {
				return err
			}
			viewConstraint.Loc.End = p.prev.End
			if stmt.Constraints == nil {
				stmt.Constraints = &nodes.List{}
			}
			stmt.Constraints.Items = append(stmt.Constraints.Items, viewConstraint)
		} else {
			// Oracle 23ai takes an alias after an out-of-line constraint too.
			aliasTok := p.cur
			name := aliasTok.Str
			if !p.isCreateTableNameStart() {
				return p.syntaxErrorAtCur()
			}
			if seenCols[name] {
				return p.syntaxErrorAtTok(aliasTok)
			}
			seenCols[name] = true
			p.advance()
			stmt.Columns.Items = append(stmt.Columns.Items, &nodes.String{Str: name})
			if p.isKeywordStr("VISIBLE") || p.isKeywordStr("INVISIBLE") {
				p.advance()
			}
			for p.isPlainViewInlineConstraintStart() {
				if err := p.parsePlainViewInlineConstraint(); err != nil {
					return err
				}
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ,
	}
	if stmt.Columns.Len() == 0 {
		return p.syntaxErrorAtCur()
	}
	return p.expectCreateTableToken(')')
}

// isPlainViewInlineConstraintStart reports whether an alias's
// inline_constraint starts here.
func (p *Parser) isPlainViewInlineConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT, kwCONSTRAINTS, kwUNIQUE, kwPRIMARY, kwREFERENCES:
		return true
	}
	return false
}

// parsePlainViewInlineConstraint reads a view alias's inline_constraint.
//
//	[ CONSTRAINT constraint_name ]
//	{ UNIQUE | PRIMARY KEY | REFERENCES [ schema. ] table [ ( column ) ] }
//	[ RELY | NORELY ] DISABLE [ NOVALIDATE ]
func (p *Parser) parsePlainViewInlineConstraint() error {
	var constTok Token
	hasConst := false
	if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwCONSTRAINTS {
		constTok = p.advance()
		hasConst = true
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}
	switch p.cur.Type {
	case kwUNIQUE:
		p.advance()
	case kwPRIMARY:
		p.advance()
		if err := p.expectCreateTableToken(kwKEY); err != nil {
			return err
		}
	case kwREFERENCES:
		p.advance()
		if !p.isCreateTableNameStart() {
			return p.syntaxErrorAtCur()
		}
		table, err := p.parseObjectName()
		if err != nil {
			return err
		}
		_ = table
		if p.cur.Type == '(' {
			p.advance()
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
			if err := p.expectCreateTableToken(')'); err != nil {
				return err
			}
		}
	default:
		if hasConst {
			return p.syntaxErrorAtTok(constTok)
		}
		return p.syntaxErrorAtCur()
	}
	return p.parseViewConstraintState()
}

// parseObjectViewClause reads an object_view_clause.
//
//	OF [ schema. ] type_name
//	{ WITH OBJECT IDENTIFIER { DEFAULT | ( attribute [, attribute ]... ) }
//	| UNDER [ schema. ] superview }
func (p *Parser) parseObjectViewClause() error {
	ofTok := p.cur
	p.advance() // OF
	if p.isKeywordStr("XMLTYPE") {
		return p.syntaxErrorAtTok(ofTok)
	}
	if !p.isCreateTableNameStart() {
		return p.syntaxErrorAtCur()
	}
	typeName, err := p.parseObjectName()
	if err != nil {
		return err
	}
	_ = typeName
	switch {
	case p.cur.Type == kwWITH:
		p.advance()
		if err := p.expectCreateTableWord("OBJECT"); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("IDENTIFIER"); err != nil {
			return err
		}
		if p.cur.Type == kwDEFAULT {
			p.advance()
			return nil
		}
		attrs, err := p.parseGraphNameList()
		if err != nil {
			return err
		}
		_ = attrs
	case p.isKeywordStr("UNDER"):
		p.advance()
		if !p.isCreateTableNameStart() {
			return p.syntaxErrorAtCur()
		}
		superview, err := p.parseObjectName()
		if err != nil {
			return err
		}
		_ = superview
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// checkEditioningViewTargets reports the first select item of an editioning
// view that is not a column: Oracle 23ai reads only columns there
// (ORA-00936 at the item).
func (p *Parser) checkEditioningViewTargets(query nodes.StmtNode) error {
	sel, ok := query.(*nodes.SelectStmt)
	if !ok || sel.TargetList == nil {
		return nil
	}
	for _, item := range sel.TargetList.Items {
		target, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		switch target.Expr.(type) {
		case *nodes.ColumnRef, *nodes.Star:
		default:
			return p.syntaxErrorAtNode(target.Expr)
		}
	}
	return nil
}

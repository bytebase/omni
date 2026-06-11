package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// CREATE VIEW
// ---------------------------------------------------------------------------

// parseCreateViewStmt parses the body of a CREATE [OR REPLACE] [SECURE]
// [RECURSIVE] VIEW statement. The CREATE keyword and OR REPLACE / SECURE /
// RECURSIVE modifiers have already been consumed; start is the Loc of the
// CREATE token.
func (p *Parser) parseCreateViewStmt(start ast.Loc, orReplace, orAlter, secure, recursive bool) (ast.Node, error) {
	p.advance() // consume VIEW

	stmt := &ast.CreateViewStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Secure:    secure,
		Recursive: recursive,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// View name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional ( column_list_with_comment )
	if p.cur.Type == '(' {
		cols, err := p.parseViewColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// Optional view_col* — column-level masking/tag bindings (outside parens)
	viewCols, err := p.parseViewCols()
	if err != nil {
		return nil, err
	}
	stmt.ViewCols = viewCols

	// Optional WITH ROW ACCESS POLICY / WITH TAG / COPY GRANTS / COMMENT = '...'
	// These can appear in any order per Snowflake docs; the legacy grammar puts them
	// before AS. We consume them in a loop until we hit AS.
	if err := p.parseViewProperties(stmt); err != nil {
		return nil, err
	}

	// AS query_statement
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	query, err := p.parseViewQuery()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseViewQuery parses the query body after AS in CREATE VIEW / CREATE MATERIALIZED VIEW.
// Handles both plain SELECT and WITH ... SELECT (CTE) forms.
func (p *Parser) parseViewQuery() (ast.Node, error) {
	if p.cur.Type == kwWITH {
		return p.parseWithQueryExpr()
	}
	return p.parseQueryExpr()
}

// parseViewColumnList parses ( col_name [COMMENT 'text'], ... ).
// Used by both CREATE VIEW and CREATE MATERIALIZED VIEW.
func (p *Parser) parseViewColumnList() ([]*ast.ViewColumn, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var cols []*ast.ViewColumn
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		col := &ast.ViewColumn{Loc: p.cur.Loc}
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		col.Name = name

		// Optional COMMENT 'text'
		if p.cur.Type == kwCOMMENT {
			p.advance() // consume COMMENT
			tok, err := p.expect(tokString)
			if err != nil {
				return nil, err
			}
			s := tok.Str
			col.Comment = &s
		}

		col.Loc.End = p.prev.Loc.End
		cols = append(cols, col)

		if p.cur.Type == ',' {
			p.advance() // consume ','
		} else {
			break
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseViewCols parses zero or more view_col entries:
//
//	column_name WITH MASKING POLICY id [USING (col, ...)] [WITH TAG (...)]
//
// Each entry starts with an identifier that is NOT followed by a WITH that
// leads to ROW or TAG at the statement level — we distinguish by checking
// if the ident is followed by WITH MASKING POLICY.
func (p *Parser) parseViewCols() ([]*ast.ViewColumn, error) {
	var cols []*ast.ViewColumn

	for {
		// A view_col starts with an identifier (column name) followed by
		// WITH MASKING POLICY or WITH TAG. We need two tokens of lookahead
		// to detect this pattern. Use the heuristic: if cur is an identifier
		// and next is WITH, and the token after WITH is MASKING or TAG —
		// consume a view_col. Otherwise stop.
		if p.cur.Type != tokIdent {
			break
		}
		next := p.peekNext()
		if next.Type != kwWITH {
			break
		}

		// We have ident WITH — save the ident start, then check what follows WITH.
		// We must consume to find the token after WITH; use a two-step look-ahead
		// trick: peek ahead to see if the token after WITH is MASKING or TAG.
		// Because we only have 1-token lookahead, consume optimistically and rollback
		// is not possible. Instead, check: if peekNext() == kwWITH, advance past the
		// column name and WITH, then check cur for MASKING or TAG.
		colStartLoc := p.cur.Loc
		colName, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		// cur is now WITH
		p.advance() // consume WITH

		// Now cur should be MASKING or TAG to be a valid view_col.
		if p.cur.Type != kwMASKING && p.cur.Type != kwTAG {
			// This is not a view_col — we over-consumed. This shouldn't happen in
			// practice because well-formed SQL has WITH [ROW|TAG|COPY|MASKING] here,
			// and ROW/COPY lead to statement-level properties, not column bindings.
			// If cur is ROW or COPY or similar, we need to put back the column name
			// and WITH. Since we can't un-advance, treat this as a parse issue and
			// stop view_col parsing. The outer loop will handle WITH at statement level.
			// Return what we have.
			break
		}

		col := &ast.ViewColumn{
			Name: colName,
			Loc:  ast.Loc{Start: colStartLoc.Start},
		}

		// Parse the WITH MASKING POLICY / WITH TAG chain.
		if err := p.parseViewColChain(col); err != nil {
			return nil, err
		}

		col.Loc.End = p.prev.Loc.End
		cols = append(cols, col)
	}

	return cols, nil
}

// parseViewColChain parses the WITH MASKING POLICY ... and/or WITH TAG (...)
// clauses that follow a column name in a view_col. Called after WITH has been
// consumed and cur is at MASKING or TAG.
func (p *Parser) parseViewColChain(col *ast.ViewColumn) error {
	for {
		switch p.cur.Type {
		case kwMASKING:
			// MASKING POLICY name [USING (col, ...)]
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return err
			}
			policyName, err := p.parseObjectName()
			if err != nil {
				return err
			}
			col.MaskingPolicy = policyName

			// Optional USING (col [, col ...])
			if p.cur.Type == kwUSING {
				p.advance() // consume USING
				if _, err := p.expect('('); err != nil {
					return err
				}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					id, err := p.parseIdent()
					if err != nil {
						return err
					}
					col.MaskingUsing = append(col.MaskingUsing, id)
					if p.cur.Type == ',' {
						p.advance() // consume ','
					} else {
						break
					}
				}
				if _, err := p.expect(')'); err != nil {
					return err
				}
			}

			// Check for optional WITH TAG following
			if p.cur.Type == kwWITH && p.peekNext().Type == kwTAG {
				p.advance() // consume WITH
				// fall through to kwTAG case on next iteration
				continue
			}
			return nil

		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			col.Tags = append(col.Tags, tags...)

			// Check for optional WITH MASKING following
			if p.cur.Type == kwWITH && p.peekNext().Type == kwMASKING {
				p.advance() // consume WITH
				continue
			}
			return nil

		default:
			return nil
		}
	}
}

// parseViewProperties parses the optional clauses that appear before AS in
// a CREATE VIEW / CREATE MATERIALIZED VIEW: WITH ROW ACCESS POLICY,
// WITH TAG, COPY GRANTS, COMMENT = '...'. Stops when it hits AS or EOF.
//
// The stmt parameter is an interface — use type switch to set fields on either
// CreateViewStmt or CreateMaterializedViewStmt.
func (p *Parser) parseViewProperties(stmt interface{}) error {
	for {
		switch p.cur.Type {
		case kwWITH:
			next := p.peekNext()
			switch next.Type {
			case kwROW:
				// WITH ROW ACCESS POLICY name ON (cols)
				p.advance() // consume WITH
				rp, err := p.parseRowAccessPolicyClause()
				if err != nil {
					return err
				}
				switch s := stmt.(type) {
				case *ast.CreateViewStmt:
					s.RowPolicy = rp
				case *ast.CreateMaterializedViewStmt:
					s.RowPolicy = rp
				}
			case kwTAG:
				// WITH TAG (name = 'val', ...)
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				switch s := stmt.(type) {
				case *ast.CreateViewStmt:
					s.Tags = append(s.Tags, tags...)
				case *ast.CreateMaterializedViewStmt:
					s.Tags = append(s.Tags, tags...)
				}
			default:
				return nil
			}

		case kwTAG:
			// TAG (...) without WITH prefix
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			switch s := stmt.(type) {
			case *ast.CreateViewStmt:
				s.Tags = append(s.Tags, tags...)
			case *ast.CreateMaterializedViewStmt:
				s.Tags = append(s.Tags, tags...)
			}

		case kwCOPY:
			if p.peekNext().Type == kwGRANTS {
				p.advance() // consume COPY
				p.advance() // consume GRANTS
				switch s := stmt.(type) {
				case *ast.CreateViewStmt:
					s.CopyGrants = true
				case *ast.CreateMaterializedViewStmt:
					s.CopyGrants = true
				}
			} else {
				return nil
			}

		case kwCOMMENT:
			// COMMENT = 'text'
			p.advance() // consume COMMENT
			if p.cur.Type == '=' {
				p.advance() // consume '='
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			switch v := stmt.(type) {
			case *ast.CreateViewStmt:
				v.Comment = &s
			case *ast.CreateMaterializedViewStmt:
				v.Comment = &s
			}

		case kwCHANGE_TRACKING:
			// CHANGE_TRACKING [=] { TRUE | FALSE } — consume and discard, mirroring
			// the CREATE TABLE handling in create_table.go. Documented for CREATE
			// VIEW (official create-view example_11 / example_13).
			p.advance() // consume CHANGE_TRACKING
			if p.cur.Type == '=' {
				p.advance() // consume optional '='
			}
			if p.cur.Type == kwTRUE || p.cur.Type == kwFALSE {
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
			}

		default:
			return nil
		}
	}
}

// parseRowAccessPolicyClause parses:
//
//	ROW ACCESS POLICY policy_name ON (col [, col ...])
//
// Called after WITH has been consumed and cur is at ROW.
func (p *Parser) parseRowAccessPolicyClause() (*ast.RowAccessPolicy, error) {
	p.advance() // consume ROW
	if _, err := p.expect(kwACCESS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	policyName, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	rp := &ast.RowAccessPolicy{PolicyName: policyName}

	// ON (col [, col ...])
	if p.cur.Type == kwON {
		p.advance() // consume ON
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			id, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			rp.Columns = append(rp.Columns, id)
			if p.cur.Type == ',' {
				p.advance() // consume ','
			} else {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	return rp, nil
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// parseCreateMaterializedViewStmt parses the body of a CREATE [OR REPLACE]
// [SECURE] [INTERACTIVE] MATERIALIZED VIEW statement. The CREATE keyword and OR
// REPLACE / SECURE / INTERACTIVE modifiers have already been consumed;
// MATERIALIZED has also been consumed. start is the Loc of the CREATE token.
func (p *Parser) parseCreateMaterializedViewStmt(start ast.Loc, orReplace, orAlter, secure, interactive bool) (ast.Node, error) {
	p.advance() // consume VIEW

	stmt := &ast.CreateMaterializedViewStmt{
		OrReplace:   orReplace,
		OrAlter:     orAlter,
		Secure:      secure,
		Interactive: interactive,
		Loc:         ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// View name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional ( column_list_with_comment )
	if p.cur.Type == '(' {
		cols, err := p.parseViewColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// Optional view_col*
	viewCols, err := p.parseViewCols()
	if err != nil {
		return nil, err
	}
	stmt.ViewCols = viewCols

	// Optional WITH ROW ACCESS POLICY / WITH TAG / COPY GRANTS / COMMENT
	if err := p.parseViewProperties(stmt); err != nil {
		return nil, err
	}

	// Optional CLUSTER BY [LINEAR] (exprs)
	if p.cur.Type == kwCLUSTER {
		p.advance() // consume CLUSTER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if p.cur.Type == kwLINEAR {
			p.advance() // consume LINEAR
			stmt.Linear = true
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.ClusterBy = exprs
	}

	// AS select_statement
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	query, err := p.parseViewQuery()
	if err != nil {
		return nil, err
	}
	stmt.Query = query

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER VIEW
// ---------------------------------------------------------------------------

// parseAlterViewStmt parses ALTER VIEW ... (all action variants).
// The ALTER keyword has already been consumed; cur is at VIEW.
func (p *Parser) parseAlterViewStmt() (ast.Node, error) {
	altTok := p.advance() // consume VIEW
	stmt := &ast.AlterViewStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// View name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO new_name
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterViewRename
		stmt.NewName = newName

	case kwSET:
		p.advance() // consume SET
		switch p.cur.Type {
		case kwCOMMENT:
			// SET COMMENT = '...'
			p.advance() // consume COMMENT
			if p.cur.Type == '=' {
				p.advance() // consume '='
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return nil, err
			}
			s := tok.Str
			stmt.Action = ast.AlterViewSetComment
			stmt.Comment = &s

		case kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterViewSetSecure
			stmt.Secure = true

		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterViewSetTag
			stmt.Tags = tags

		case kwJOIN:
			// SET JOIN POLICY name [ { ALLOWED | ENFORCED } JOIN KEYS ( col, ... ) ]
			if err := p.parseAlterViewSetJoinPolicy(stmt); err != nil {
				return nil, err
			}

		default:
			// SET AGGREGATION POLICY name [ ENTITY KEY ( col, ... ) ] [ FORCE ].
			// AGGREGATION is not a reserved keyword, so it lexes as an identifier.
			if p.curIsWord("AGGREGATION") {
				if err := p.parseAlterViewSetAggregationPolicy(stmt); err != nil {
					return nil, err
				}
				break
			}
			return nil, p.syntaxErrorAtCur()
		}

	case kwUNSET:
		p.advance() // consume UNSET
		switch p.cur.Type {
		case kwCOMMENT:
			p.advance() // consume COMMENT
			stmt.Action = ast.AlterViewUnsetComment

		case kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterViewUnsetSecure

		case kwTAG:
			names, err := p.parseUnsetTagList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterViewUnsetTag
			stmt.UnsetTags = names

		case kwJOIN:
			// UNSET JOIN POLICY
			p.advance() // consume JOIN
			if _, err := p.expect(kwPOLICY); err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterViewUnsetJoinPolicy

		default:
			// UNSET AGGREGATION POLICY (AGGREGATION lexes as an identifier).
			if p.curIsWord("AGGREGATION") {
				p.advance() // consume AGGREGATION
				if _, err := p.expect(kwPOLICY); err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewUnsetAggregationPolicy
				break
			}
			return nil, p.syntaxErrorAtCur()
		}

	case kwADD:
		// ADD ROW ACCESS POLICY policy_name ON (cols)
		p.advance() // consume ADD
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		policyName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterViewAddRowAccessPolicy
		stmt.PolicyName = policyName
		stmt.PolicyCols = cols

	case kwDROP:
		p.advance() // consume DROP
		switch p.cur.Type {
		case kwROW:
			// DROP ROW ACCESS POLICY policy_name
			// OR DROP ALL ROW ACCESS POLICIES
			p.advance() // consume ROW
			if _, err := p.expect(kwACCESS); err != nil {
				return nil, err
			}
			if p.cur.Type == kwPOLICY {
				p.advance() // consume POLICY
				policyName, err := p.parseObjectName()
				if err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewDropRowAccessPolicy
				stmt.PolicyName = policyName

				// Documented drop-and-add combination:
				// DROP ROW ACCESS POLICY p1, ADD ROW ACCESS POLICY p2 ON (cols)
				if p.cur.Type == ',' && p.peekNext().Type == kwADD {
					p.advance() // consume ','
					p.advance() // consume ADD
					if _, err := p.expect(kwROW); err != nil {
						return nil, err
					}
					if _, err := p.expect(kwACCESS); err != nil {
						return nil, err
					}
					if _, err := p.expect(kwPOLICY); err != nil {
						return nil, err
					}
					addName, err := p.parseObjectName()
					if err != nil {
						return nil, err
					}
					if _, err := p.expect(kwON); err != nil {
						return nil, err
					}
					addCols, err := p.parseIdentListInParens()
					if err != nil {
						return nil, err
					}
					stmt.AddPolicyName = addName
					stmt.AddPolicyCols = addCols
				}
			} else if p.cur.Type == kwPOLICIES {
				p.advance() // consume POLICIES
				stmt.Action = ast.AlterViewDropAllRowAccessPolicies
			} else {
				return nil, p.syntaxErrorAtCur()
			}
		case kwALL:
			// DROP ALL ROW ACCESS POLICIES
			p.advance() // consume ALL
			if _, err := p.expect(kwROW); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwACCESS); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwPOLICIES); err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterViewDropAllRowAccessPolicies
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwALTER, kwMODIFY:
		// ALTER|MODIFY [COLUMN] col SET MASKING POLICY / UNSET MASKING POLICY
		//                          | SET TAG / UNSET TAG
		// — comma-separated list form included:
		// MODIFY COLUMN a SET MASKING POLICY p1, COLUMN b SET MASKING POLICY p2
		p.advance() // consume ALTER or MODIFY
		for {
			action, err := p.parseAlterViewColumnAction()
			if err != nil {
				return nil, err
			}
			stmt.ColumnActions = append(stmt.ColumnActions, action)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		// Mirror the first action into the legacy single-action fields.
		first := stmt.ColumnActions[0]
		stmt.Action = first.Action
		stmt.Column = first.Column
		stmt.MaskingPolicy = first.MaskingPolicy
		stmt.MaskingUsing = first.MaskingUsing
		stmt.Tags = first.Tags
		stmt.UnsetTags = first.UnsetTags

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterViewColumnAction parses ONE element of an ALTER VIEW
// MODIFY/ALTER column action list:
//
//	[COLUMN] <col> { SET MASKING POLICY <p> [USING (col, ...)] [FORCE]
//	               | UNSET MASKING POLICY
//	               | SET TAG <t> = '<v>' [, ...]
//	               | UNSET TAG <t> [, ...] }
//
// The leading ALTER|MODIFY keyword (or the list comma) has already been
// consumed.
func (p *Parser) parseAlterViewColumnAction() (*ast.AlterViewColumnAction, error) {
	action := &ast.AlterViewColumnAction{Loc: ast.Loc{Start: p.cur.Loc.Start}}

	// Optional COLUMN keyword
	if p.cur.Type == kwCOLUMN {
		p.advance() // consume COLUMN
	}
	colName, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	action.Column = colName

	switch p.cur.Type {
	case kwSET:
		p.advance() // consume SET
		switch p.cur.Type {
		case kwMASKING:
			// SET MASKING POLICY policy_name [USING (col, ...)] [FORCE]
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return nil, err
			}
			policyName, err := p.parseObjectName()
			if err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnSetMaskingPolicy
			action.MaskingPolicy = policyName

			// Optional USING (col, ...)
			if p.cur.Type == kwUSING {
				p.advance() // consume USING
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					id, err := p.parseIdent()
					if err != nil {
						return nil, err
					}
					action.MaskingUsing = append(action.MaskingUsing, id)
					if p.cur.Type == ',' {
						p.advance()
					} else {
						break
					}
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}

			// Optional FORCE
			if p.cur.Type == kwFORCE {
				p.advance() // consume FORCE
				action.Force = true
			}

		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnSetTag
			action.Tags = tags

		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwUNSET:
		p.advance() // consume UNSET
		switch p.cur.Type {
		case kwMASKING:
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnUnsetMaskingPolicy

		case kwTAG:
			names, err := p.parseUnsetTagList()
			if err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnUnsetTag
			action.UnsetTags = names

		default:
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// parseAlterViewSetJoinPolicy parses
//
//	SET JOIN POLICY <name> [ { ALLOWED | ENFORCED } JOIN KEYS ( col, ... ) ]
//
// On entry cur is the JOIN keyword. ALLOWED lexes as an identifier; ENFORCED is
// a reserved keyword.
func (p *Parser) parseAlterViewSetJoinPolicy(stmt *ast.AlterViewStmt) error {
	p.advance() // consume JOIN
	if _, err := p.expect(kwPOLICY); err != nil {
		return err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterViewSetJoinPolicy
	stmt.PolicyName = name

	// Optional { ALLOWED | ENFORCED } JOIN KEYS ( col, ... ).
	var kind string
	if p.curIsWord("ALLOWED") {
		kind = "ALLOWED"
		p.advance() // consume ALLOWED
	} else if p.cur.Type == kwENFORCED {
		kind = "ENFORCED"
		p.advance() // consume ENFORCED
	}
	if kind != "" {
		if _, err := p.expect(kwJOIN); err != nil {
			return err
		}
		if _, err := p.expect(kwKEYS); err != nil {
			return err
		}
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return err
		}
		stmt.PolicyKeyKind = kind
		stmt.PolicyKeyCols = cols
	}
	return nil
}

// parseAlterViewSetAggregationPolicy parses
//
//	SET AGGREGATION POLICY <name> [ ENTITY KEY ( col, ... ) ] [ FORCE ]
//
// On entry cur is the AGGREGATION identifier. ENTITY lexes as an identifier.
func (p *Parser) parseAlterViewSetAggregationPolicy(stmt *ast.AlterViewStmt) error {
	p.advance() // consume AGGREGATION
	if _, err := p.expect(kwPOLICY); err != nil {
		return err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterViewSetAggregationPolicy
	stmt.PolicyName = name

	// Optional ENTITY KEY ( col, ... ).
	if p.curIsWord("ENTITY") {
		p.advance() // consume ENTITY
		if _, err := p.expect(kwKEY); err != nil {
			return err
		}
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return err
		}
		stmt.PolicyKeyCols = cols
	}

	// Optional FORCE.
	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		stmt.PolicyForce = true
	}
	return nil
}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// parseAlterMaterializedViewStmt parses ALTER MATERIALIZED VIEW ... (all action variants).
// The ALTER keyword has already been consumed; MATERIALIZED has also been consumed;
// cur is at VIEW.
func (p *Parser) parseAlterMaterializedViewStmt() (ast.Node, error) {
	altTok := p.advance() // consume VIEW
	stmt := &ast.AlterMaterializedViewStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Note: the legacy grammar does NOT have IF EXISTS for ALTER MATERIALIZED VIEW.
	// The grammar uses plain `id_` (not `if_exists? object_name`).

	// View name (1-part only per legacy grammar, but we use parseObjectName for generality)
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO new_name
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterMVRename
		stmt.NewName = newName

	case kwCLUSTER:
		// CLUSTER BY (exprs)
		p.advance() // consume CLUSTER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if p.cur.Type == kwLINEAR {
			p.advance() // consume LINEAR
			stmt.Linear = true
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterMVClusterBy
		stmt.ClusterBy = exprs

	case kwDROP:
		// DROP CLUSTERING KEY
		p.advance() // consume DROP
		if _, err := p.expect(kwCLUSTERING); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterMVDropClusteringKey

	case kwSUSPEND:
		// SUSPEND [RECLUSTER]
		p.advance() // consume SUSPEND
		if p.cur.Type == kwRECLUSTER {
			p.advance() // consume RECLUSTER
			stmt.Action = ast.AlterMVSuspendRecluster
		} else {
			stmt.Action = ast.AlterMVSuspend
		}

	case kwRESUME:
		// RESUME [RECLUSTER]
		p.advance() // consume RESUME
		if p.cur.Type == kwRECLUSTER {
			p.advance() // consume RECLUSTER
			stmt.Action = ast.AlterMVResumeRecluster
		} else {
			stmt.Action = ast.AlterMVResume
		}

	case kwSET:
		// SET [SECURE] [COMMENT = '...']
		// Can be SET SECURE, SET COMMENT = '...', or SET SECURE COMMENT = '...'
		p.advance() // consume SET
		switch p.cur.Type {
		case kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterMVSetSecure
			stmt.Secure = true
			// Check if COMMENT follows (SET SECURE COMMENT = '...')
			if p.cur.Type == kwCOMMENT {
				p.advance() // consume COMMENT
				if p.cur.Type == '=' {
					p.advance()
				}
				tok, err := p.expect(tokString)
				if err != nil {
					return nil, err
				}
				s := tok.Str
				stmt.Comment = &s
			}
		case kwCOMMENT:
			p.advance() // consume COMMENT
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return nil, err
			}
			s := tok.Str
			stmt.Action = ast.AlterMVSetComment
			stmt.Comment = &s
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwUNSET:
		// UNSET SECURE | UNSET COMMENT
		p.advance() // consume UNSET
		switch p.cur.Type {
		case kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterMVUnsetSecure
		case kwCOMMENT:
			p.advance() // consume COMMENT
			stmt.Action = ast.AlterMVUnsetComment
		default:
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

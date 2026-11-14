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
		setTok := p.advance() // consume SET
		switch p.cur.Type {
		case kwCHANGE_TRACKING:
			// SET CHANGE_TRACKING = { TRUE | FALSE }
			p.advance() // consume CHANGE_TRACKING
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if p.cur.Type != kwTRUE && p.cur.Type != kwFALSE {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.ChangeTracking = p.advance().Type == kwTRUE
			stmt.Action = ast.AlterViewSetChangeTracking

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
			// SET TAG name = 'val' [, ...] | SET TAG (name = 'val', ...)
			tags, err := p.parseSetTagClause()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterViewSetTag
			stmt.Tags = tags

		case kwJOIN:
			// SET JOIN POLICY name [ FORCE ] [ { ALLOWED | ENFORCED } JOIN KEYS ( col, ... ) ]
			if err := p.parseAlterViewSetJoinPolicy(stmt); err != nil {
				return nil, err
			}

		default:
			switch {
			case p.startsPolicyWord("AGGREGATION"):
				// SET AGGREGATION POLICY name [ ENTITY KEY ( col, ... ) ] [ FORCE ].
				// AGGREGATION is not a reserved keyword, so it lexes as an identifier.
				name, keys, force, err := p.parseSetAggregationPolicy()
				if err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewSetAggregationPolicy
				stmt.PolicyName = name
				stmt.PolicyKeyCols = keys
				stmt.PolicyForce = force
			case p.curIsWord("CONTACT"):
				// SET CONTACT purpose = contact [ , ... ]
				contact, err := p.parseContactValue()
				if err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewSetContact
				stmt.Contact = contact
			case p.startsDataMetricSchedule():
				dm, err := p.parseDataMetricSchedule(true, setTok.Loc.Start)
				if err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewDataMetric
				stmt.DataMetric = dm
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}

	case kwUNSET:
		unsetTok := p.advance() // consume UNSET
		switch p.cur.Type {
		case kwCOMMENT:
			p.advance() // consume COMMENT
			stmt.Action = ast.AlterViewUnsetComment

		case kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterViewUnsetSecure

		case kwTAG:
			// UNSET TAG name [, ...] | UNSET TAG (name, ...)
			names, err := p.parseUnsetTagClause()
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
			switch {
			case p.curIsWord("AGGREGATION"):
				// UNSET AGGREGATION POLICY (AGGREGATION lexes as an identifier).
				p.advance() // consume AGGREGATION
				if _, err := p.expect(kwPOLICY); err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewUnsetAggregationPolicy
			case p.curIsWord("CONTACT"):
				// UNSET CONTACT purpose
				p.advance() // consume CONTACT
				purpose, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewUnsetContact
				stmt.Contact = purpose.String()
			case p.startsDataMetricSchedule():
				dm, err := p.parseDataMetricSchedule(false, unsetTok.Loc.Start)
				if err != nil {
					return nil, err
				}
				stmt.Action = ast.AlterViewDataMetric
				stmt.DataMetric = dm
			default:
				return nil, p.syntaxErrorAtCur()
			}
		}

	case kwADD, kwDROP, kwMODIFY, kwALTER:
		verb := p.cur
		switch {
		case verb.Type != kwALTER && p.verbStartsDataMetricFunction():
			// { ADD | DROP | MODIFY } DATA METRIC FUNCTION ...
			p.advance() // consume the verb
			kind := ast.DataMetricAdd
			switch verb.Type {
			case kwDROP:
				kind = ast.DataMetricDrop
			case kwMODIFY:
				kind = ast.DataMetricModify
			}
			dm, err := p.parseDataMetricFunctionAction(kind, verb.Loc.Start)
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterViewDataMetric
			stmt.DataMetric = dm
		case verb.Type == kwALTER || verb.Type == kwMODIFY:
			p.advance() // consume ALTER or MODIFY
			if err := p.parseAlterViewColumnActions(stmt); err != nil {
				return nil, err
			}
		case verb.Type == kwDROP:
			p.advance() // consume DROP
			if err := p.parseAlterViewDropRowAccessPolicy(stmt); err != nil {
				return nil, err
			}
		case p.peekNext().Type == kwROW:
			p.advance() // consume ADD
			if err := p.parseAlterViewAddRowAccessPolicy(stmt); err != nil {
				return nil, err
			}
		default:
			p.advance() // consume ADD
			if err := p.parseAlterViewAddColumn(stmt); err != nil {
				return nil, err
			}
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// verbStartsDataMetricFunction reports whether the ADD / DROP / MODIFY at cur
// is followed by DATA METRIC FUNCTION.
func (p *Parser) verbStartsDataMetricFunction() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // consume the verb
	return p.startsDataMetricFunction()
}

// parseAlterViewAddRowAccessPolicy parses ROW ACCESS POLICY policy_name ON
// (cols) after ADD.
func (p *Parser) parseAlterViewAddRowAccessPolicy(stmt *ast.AlterViewStmt) error {
	name, cols, err := p.parseAddRowAccessPolicyTail()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterViewAddRowAccessPolicy
	stmt.PolicyName = name
	stmt.PolicyCols = cols
	return nil
}

// parseAlterViewDropRowAccessPolicy parses, after DROP,
//
//	ROW ACCESS POLICY <p> [ , ADD ROW ACCESS POLICY <p2> ON ( <col> [ , ... ] ) ]
//	ALL ROW ACCESS POLICIES
func (p *Parser) parseAlterViewDropRowAccessPolicy(stmt *ast.AlterViewStmt) error {
	all, name, addName, addCols, err := p.parseDropRowAccessPolicyTail()
	if err != nil {
		return err
	}
	if all {
		stmt.Action = ast.AlterViewDropAllRowAccessPolicies
		return nil
	}
	stmt.Action = ast.AlterViewDropRowAccessPolicy
	stmt.PolicyName = name
	stmt.AddPolicyName = addName
	stmt.AddPolicyCols = addCols
	return nil
}

// parseAlterViewColumnActions parses the comma-separated column action list
// after ALTER / MODIFY:
//
//	MODIFY COLUMN a SET MASKING POLICY p1, COLUMN b SET MASKING POLICY p2
func (p *Parser) parseAlterViewColumnActions(stmt *ast.AlterViewStmt) error {
	actions, err := p.parseColumnPolicyActions()
	if err != nil {
		return err
	}
	stmt.ColumnActions = actions
	// Mirror the first action into the legacy single-action fields.
	first := actions[0]
	stmt.Action = first.Action
	stmt.Column = first.Column
	stmt.MaskingPolicy = first.MaskingPolicy
	stmt.MaskingUsing = first.MaskingUsing
	stmt.ProjectionPolicy = first.ProjectionPolicy
	stmt.Tags = first.Tags
	stmt.UnsetTags = first.UnsetTags
	return nil
}

// parseAlterViewAddColumn parses, after ADD,
//
//	[ COLUMN ] [ IF NOT EXISTS ] <col> <type>
//	  [ [ WITH ] MASKING POLICY <p> [ USING ( <col> , <col> [ , ... ] ) ] ]
//	  [ [ WITH ] PROJECTION POLICY <p> ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
func (p *Parser) parseAlterViewAddColumn(stmt *ast.AlterViewStmt) error {
	if p.cur.Type == kwCOLUMN {
		p.advance() // consume COLUMN
	}
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwNOT); err != nil {
			return err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return err
		}
		stmt.IfNotExists = true
	}
	col, err := p.parseIdent()
	if err != nil {
		return err
	}
	typ, err := p.parseDataType()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterViewAddColumn
	stmt.Column = col
	stmt.ColumnType = typ

	// [ WITH ] precedes each clause; the clauses come in the documented order.
	withThen := func(word string) bool {
		if p.cur.Type == kwWITH {
			return p.peekIsWord(word)
		}
		return p.curIsWord(word)
	}
	if withThen("MASKING") {
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		p.advance() // consume MASKING
		if _, err := p.expect(kwPOLICY); err != nil {
			return err
		}
		if stmt.MaskingPolicy, err = p.parseObjectName(); err != nil {
			return err
		}
		if p.cur.Type == kwUSING {
			p.advance() // consume USING
			if stmt.MaskingUsing, err = p.parseIdentListInParens(); err != nil {
				return err
			}
			if len(stmt.MaskingUsing) < 2 {
				// USING names the masked column and at least one condition column.
				return p.syntaxErrorAtCur()
			}
		}
	}
	if withThen("PROJECTION") {
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		if err := p.expectPolicyWord(); err != nil {
			return err
		}
		if stmt.ProjectionPolicy, err = p.parseObjectName(); err != nil {
			return err
		}
	}
	if withThen("TAG") {
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		if stmt.Tags, err = p.parseTagAssignments(); err != nil {
			return err
		}
	}
	return nil
}

// parseAlterViewColumnAction parses ONE element of an ALTER VIEW
// MODIFY/ALTER column action list:
//
//	[COLUMN] <col> { SET MASKING POLICY <p> [USING (col, ...)] [FORCE]
//	               | UNSET MASKING POLICY
//	               | SET PROJECTION POLICY <p> [FORCE]
//	               | UNSET PROJECTION POLICY
//	               | SET TAG <t> = '<v>' [, ...]
//	               | UNSET TAG <t> [, ...] }
//
// The leading ALTER|MODIFY keyword (or the list comma) has already been
// consumed. ALTER DYNAMIC TABLE's column actions share this grammar.
//
//	https://docs.snowflake.com/en/sql-reference/sql/alter-view
//	https://docs.snowflake.com/en/sql-reference/sql/alter-dynamic-table
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
			tags, err := p.parseSetTagClause()
			if err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnSetTag
			action.Tags = tags

		default:
			if !p.startsPolicyWord("PROJECTION") {
				return nil, p.syntaxErrorAtCur()
			}
			// SET PROJECTION POLICY policy_name [FORCE]
			if err := p.expectPolicyWord(); err != nil {
				return nil, err
			}
			policyName, err := p.parseObjectName()
			if err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnSetProjectionPolicy
			action.ProjectionPolicy = policyName
			if p.cur.Type == kwFORCE {
				p.advance() // consume FORCE
				action.Force = true
			}
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
			names, err := p.parseUnsetTagClause()
			if err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnUnsetTag
			action.UnsetTags = names

		default:
			if !p.startsPolicyWord("PROJECTION") {
				return nil, p.syntaxErrorAtCur()
			}
			if err := p.expectPolicyWord(); err != nil {
				return nil, err
			}
			action.Action = ast.AlterViewColumnUnsetProjectionPolicy
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// parseAlterViewSetJoinPolicy parses
//
//	SET JOIN POLICY <name> [ FORCE ] [ { ALLOWED | ENFORCED } JOIN KEYS ( col, ... ) ]
//
// On entry cur is the JOIN keyword. ALLOWED lexes as an identifier; ENFORCED is
// a reserved keyword. The page writes [ FORCE ]; the JOIN KEYS clause is what
// omni has always read.
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
	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		stmt.PolicyForce = true
	}

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
		// SET SECURE [ CONTACT ... ] [ COMMENT = '...' ]
		// SET CONTACT purpose = contact [ , ... ] [ COMMENT = '...' ]
		// SET COMMENT = '...'
		// SET DATA_METRIC_SCHEDULE = '...'
		setTok := p.advance() // consume SET
		switch {
		case p.cur.Type == kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterMVSetSecure
			stmt.Secure = true
			if err := p.parseAlterMVContactComment(stmt); err != nil {
				return nil, err
			}
		case p.curIsWord("CONTACT"):
			stmt.Action = ast.AlterMVSetContact
			if err := p.parseAlterMVContactComment(stmt); err != nil {
				return nil, err
			}
		case p.cur.Type == kwCOMMENT:
			stmt.Action = ast.AlterMVSetComment
			if err := p.parseAlterMVContactComment(stmt); err != nil {
				return nil, err
			}
		case p.startsDataMetricSchedule():
			dm, err := p.parseDataMetricSchedule(true, setTok.Loc.Start)
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterMVDataMetric
			stmt.DataMetric = dm
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwUNSET:
		// UNSET { SECURE | CONTACT purpose | COMMENT | DATA_METRIC_SCHEDULE }
		unsetTok := p.advance() // consume UNSET
		switch {
		case p.cur.Type == kwSECURE:
			p.advance() // consume SECURE
			stmt.Action = ast.AlterMVUnsetSecure
		case p.cur.Type == kwCOMMENT:
			p.advance() // consume COMMENT
			stmt.Action = ast.AlterMVUnsetComment
		case p.curIsWord("CONTACT"):
			p.advance() // consume CONTACT
			purpose, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterMVUnsetContact
			stmt.Contact = purpose.String()
		case p.startsDataMetricSchedule():
			dm, err := p.parseDataMetricSchedule(false, unsetTok.Loc.Start)
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterMVDataMetric
			stmt.DataMetric = dm
		default:
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterMVContactComment parses the optional
// [ CONTACT purpose = contact [ , ... ] ] [ COMMENT = '...' ] tail of ALTER
// MATERIALIZED VIEW ... SET. A SET COMMENT or SET CONTACT action requires its
// clause.
func (p *Parser) parseAlterMVContactComment(stmt *ast.AlterMaterializedViewStmt) error {
	if p.curIsWord("CONTACT") {
		contact, err := p.parseContactValue()
		if err != nil {
			return err
		}
		stmt.Contact = contact
	} else if stmt.Action == ast.AlterMVSetContact {
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwCOMMENT {
		p.advance()            // consume COMMENT
		if p.cur.Type == '=' { // omni has always read the '=' as optional here
			p.advance()
		}
		tok, err := p.expect(tokString)
		if err != nil {
			return err
		}
		s := tok.Str
		stmt.Comment = &s
	} else if stmt.Action == ast.AlterMVSetComment {
		return p.syntaxErrorAtCur()
	}
	return nil
}

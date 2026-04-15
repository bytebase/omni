package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// ALTER TABLE statement parser
// ---------------------------------------------------------------------------

// parseAlterTableStmt parses ALTER TABLE [IF EXISTS] name action [, action ...].
// The ALTER keyword has already been consumed.
func (p *Parser) parseAlterTableStmt() (ast.Node, error) {
	altTok := p.advance() // consume TABLE
	stmt := &ast.AlterTableStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Table name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Parse one or more comma-separated actions.
	// Not all action types allow comma-separated chaining, but we handle it
	// generically at this level (the grammar allows multiple actions for ADD
	// COLUMN, DROP COLUMN, etc.).
	for {
		action, err := p.parseAlterTableAction()
		if err != nil {
			return nil, err
		}
		stmt.Actions = append(stmt.Actions, action)

		// Some action forms allow comma-chaining (ADD COLUMN col, col; DROP COLUMN a, b...)
		// but the top-level multi-action form (action , action) is handled here.
		// We stop if the next token is not a comma that introduces a new top-level action.
		// For ADD COLUMN and DROP COLUMN we allow comma continuations inside the action
		// itself; at this level we stop on anything other than another action keyword.
		if p.cur.Type != ',' {
			break
		}
		// Peek ahead: if the comma is followed by an action keyword we continue;
		// otherwise stop (it might be a semi-colon, EOF, etc.).
		next := p.peekNext()
		if !p.isAlterTableActionKeyword(next.Type) {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// isAlterTableActionKeyword returns true if tok starts a new ALTER TABLE action.
func (p *Parser) isAlterTableActionKeyword(tokType int) bool {
	switch tokType {
	case kwRENAME, kwSWAP, kwADD, kwDROP, kwALTER, kwMODIFY,
		kwCLUSTER, kwSUSPEND, kwRESUME, kwSET, kwUNSET, kwRECLUSTER:
		return true
	}
	return false
}

// parseAlterTableAction parses a single ALTER TABLE action and returns it.
func (p *Parser) parseAlterTableAction() (*ast.AlterTableAction, error) {
	action := &ast.AlterTableAction{Loc: p.cur.Loc}

	switch p.cur.Type {
	case kwRENAME:
		p.advance() // consume RENAME
		switch p.cur.Type {
		case kwTO:
			// RENAME TO new_name
			p.advance() // consume TO
			newName, err := p.parseObjectName()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableRename
			action.NewName = newName

		case kwCOLUMN:
			// RENAME COLUMN old TO new
			p.advance() // consume COLUMN
			oldName, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			newCol, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableRenameColumn
			action.OldName = oldName
			action.NewColName = newCol

		case kwCONSTRAINT:
			// RENAME CONSTRAINT old TO new
			p.advance() // consume CONSTRAINT
			oldName, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			newName, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableRenameConstraint
			action.ConstraintName = oldName
			action.NewConstraintName = newName

		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwSWAP:
		// SWAP WITH other_table
		p.advance() // consume SWAP
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		otherName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		action.Kind = ast.AlterTableSwapWith
		action.NewName = otherName

	case kwADD:
		p.advance() // consume ADD
		if err := p.parseAlterTableAddAction(action); err != nil {
			return nil, err
		}

	case kwDROP:
		p.advance() // consume DROP
		if err := p.parseAlterTableDropAction(action); err != nil {
			return nil, err
		}

	case kwALTER, kwMODIFY:
		p.advance() // consume ALTER or MODIFY
		if err := p.parseAlterTableModifyAction(action); err != nil {
			return nil, err
		}

	case kwCLUSTER:
		// CLUSTER BY [LINEAR] (exprs)
		p.advance() // consume CLUSTER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		if p.cur.Type == kwLINEAR {
			p.advance() // consume LINEAR
			action.Linear = true
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
		action.Kind = ast.AlterTableClusterBy
		action.ClusterBy = exprs

	case kwRECLUSTER:
		// RECLUSTER [MAX_SIZE = n] [WHERE expr]
		p.advance() // consume RECLUSTER
		action.Kind = ast.AlterTableRecluster
		if p.cur.Type == kwMAX_SIZE {
			p.advance() // consume MAX_SIZE
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokInt)
			if err != nil {
				return nil, err
			}
			v := tok.Ival
			action.ReclusterMaxSize = &v
		}
		if p.cur.Type == kwWHERE {
			p.advance() // consume WHERE
			whereExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			action.ReclusterWhere = whereExpr
		}

	case kwSUSPEND:
		// SUSPEND RECLUSTER
		p.advance() // consume SUSPEND
		if p.cur.Type == kwRECLUSTER {
			p.advance() // consume RECLUSTER
			action.Kind = ast.AlterTableSuspendRecluster
		} else {
			// plain SUSPEND (table suspend)
			action.Kind = ast.AlterTableSuspendRecluster
		}

	case kwRESUME:
		// RESUME RECLUSTER
		p.advance() // consume RESUME
		if p.cur.Type == kwRECLUSTER {
			p.advance() // consume RECLUSTER
			action.Kind = ast.AlterTableResumeRecluster
		} else {
			action.Kind = ast.AlterTableResumeRecluster
		}

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			// SET TAG (name = 'val', ...)
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableSetTag
			action.Tags = tags
		} else {
			// SET table properties
			props, err := p.parseTableSetProps()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableSet
			action.Props = props
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			// UNSET TAG (name, ...)
			names, err := p.parseUnsetTagList()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableUnsetTag
			action.UnsetTags = names
		} else {
			// UNSET property_name [, ...]
			props, err := p.parseTableUnsetProps()
			if err != nil {
				return nil, err
			}
			action.Kind = ast.AlterTableUnset
			action.UnsetProps = props
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	action.Loc.End = p.prev.Loc.End
	return action, nil
}

// ---------------------------------------------------------------------------
// ADD sub-actions
// ---------------------------------------------------------------------------

// parseAlterTableAddAction parses the portion after ADD in an ALTER TABLE action.
func (p *Parser) parseAlterTableAddAction(action *ast.AlterTableAction) error {
	switch p.cur.Type {
	case kwCOLUMN, kwIF:
		// ADD [COLUMN] [IF NOT EXISTS] col_def [, col_def ...]
		// Optionally consume COLUMN keyword
		if p.cur.Type == kwCOLUMN {
			p.advance()
		}
		if p.cur.Type == kwIF {
			if p.peekNext().Type == kwNOT {
				p.advance() // consume IF
				p.advance() // consume NOT
				if _, err := p.expect(kwEXISTS); err != nil {
					return err
				}
				action.IfNotExists = true
			}
		}
		action.Kind = ast.AlterTableAddColumn
		return p.parseAddColumnList(action)

	case kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN:
		// ADD [CONSTRAINT name] PK/UK/FK
		con, err := p.parseOutOfLineConstraint()
		if err != nil {
			return err
		}
		action.Kind = ast.AlterTableAddConstraint
		action.Constraint = con
		return nil

	case kwROW:
		// ADD ROW ACCESS POLICY name ON (cols)
		p.advance() // consume ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return err
		}
		policyName, err := p.parseObjectName()
		if err != nil {
			return err
		}
		if _, err := p.expect(kwON); err != nil {
			return err
		}
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return err
		}
		action.Kind = ast.AlterTableAddRowAccessPolicy
		action.PolicyName = policyName
		action.PolicyCols = cols
		return nil

	case kwSEARCH:
		// ADD SEARCH OPTIMIZATION [ON search_method(target), ...]
		p.advance() // consume SEARCH
		if _, err := p.expect(kwOPTIMIZATION); err != nil {
			return err
		}
		action.Kind = ast.AlterTableAddSearchOpt
		if p.cur.Type == kwON {
			p.advance() // consume ON
			if err := p.parseSearchOptTargets(action); err != nil {
				return err
			}
		}
		return nil

	default:
		// ADD col_def (no COLUMN keyword, no IF NOT EXISTS)
		// Check whether current token looks like the start of a column name
		if p.cur.Type == tokIdent || p.isNonReservedKeyword(p.cur.Type) {
			action.Kind = ast.AlterTableAddColumn
			return p.parseAddColumnList(action)
		}
		return p.syntaxErrorAtCur()
	}
}

// parseAddColumnList parses one or more column definitions for ADD COLUMN.
// Called after ADD [COLUMN] [IF NOT EXISTS] has already been consumed.
func (p *Parser) parseAddColumnList(action *ast.AlterTableAction) error {
	col, err := p.parseColumnDef()
	if err != nil {
		return err
	}
	action.Columns = append(action.Columns, col)

	for p.cur.Type == ',' {
		// Peek ahead: if next is another column name (ident / non-reserved kw),
		// keep going; otherwise stop (could be next top-level action).
		next := p.peekNext()
		if !p.isColumnNameToken(next.Type) {
			break
		}
		p.advance() // consume ','
		col, err = p.parseColumnDef()
		if err != nil {
			return err
		}
		action.Columns = append(action.Columns, col)
	}
	return nil
}

// parseSearchOptTargets parses the ON method(target) list for SEARCH OPTIMIZATION.
func (p *Parser) parseSearchOptTargets(action *ast.AlterTableAction) error {
	for {
		// method can be EQUALITY, SUBSTRING, GEO, or an identifier
		var method string
		switch p.cur.Type {
		case kwEQUALITY:
			p.advance()
			method = "EQUALITY"
		case kwSUBSTRING:
			p.advance()
			method = "SUBSTRING"
		case kwGEO:
			p.advance()
			method = "GEO"
		case tokIdent:
			tok := p.advance()
			method = tok.Str
		default:
			return p.syntaxErrorAtCur()
		}
		// (STAR | expr)
		if _, err := p.expect('('); err != nil {
			return err
		}
		var target string
		if p.cur.Type == '*' {
			p.advance()
			target = "*"
		} else {
			// skip expression tokens until matching ')'
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			_ = expr
			target = "expr"
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
		action.SearchOptOn = append(action.SearchOptOn, method+"("+target+")")
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return nil
}

// ---------------------------------------------------------------------------
// DROP sub-actions
// ---------------------------------------------------------------------------

// parseAlterTableDropAction parses the portion after DROP in an ALTER TABLE action.
func (p *Parser) parseAlterTableDropAction(action *ast.AlterTableAction) error {
	switch p.cur.Type {
	case kwCOLUMN, kwIF:
		// DROP [COLUMN] [IF EXISTS] col [, col ...]
		if p.cur.Type == kwCOLUMN {
			p.advance()
		}
		if p.cur.Type == kwIF {
			if p.peekNext().Type == kwEXISTS {
				p.advance() // consume IF
				p.advance() // consume EXISTS
				action.IfExists = true
			}
		}
		action.Kind = ast.AlterTableDropColumn
		return p.parseDropColumnList(action)

	case kwCONSTRAINT:
		// DROP CONSTRAINT name [CASCADE|RESTRICT]
		p.advance() // consume CONSTRAINT
		cname, err := p.parseIdent()
		if err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropConstraint
		action.ConstraintName = cname
		p.parseDropCascadeRestrict(action)
		return nil

	case kwPRIMARY:
		// DROP PRIMARY KEY [CASCADE|RESTRICT]
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropConstraint
		action.IsPrimaryKey = true
		p.parseDropCascadeRestrict(action)
		return nil

	case kwUNIQUE:
		// DROP UNIQUE (cols)? [CASCADE|RESTRICT]
		p.advance() // consume UNIQUE
		action.Kind = ast.AlterTableDropConstraint
		action.DropUnique = true
		// optional column list
		if p.cur.Type == '(' {
			if err := p.skipParenthesized(); err != nil {
				return err
			}
		}
		p.parseDropCascadeRestrict(action)
		return nil

	case kwFOREIGN:
		// DROP FOREIGN KEY (cols)? [CASCADE|RESTRICT]
		p.advance() // consume FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropConstraint
		action.DropForeignKey = true
		if p.cur.Type == '(' {
			if err := p.skipParenthesized(); err != nil {
				return err
			}
		}
		p.parseDropCascadeRestrict(action)
		return nil

	case kwCLUSTERING:
		// DROP CLUSTERING KEY
		p.advance() // consume CLUSTERING
		if _, err := p.expect(kwKEY); err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropClusterKey
		return nil

	case kwROW:
		// DROP ROW ACCESS POLICY name
		p.advance() // consume ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return err
		}
		policyName, err := p.parseObjectName()
		if err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropRowAccessPolicy
		action.PolicyName = policyName

		// Optional: , ADD ROW ACCESS POLICY name ON (cols)
		// (combined drop+add form — treat as a separate grammar note,
		// we just consume and discard the extra ADD for now)
		if p.cur.Type == ',' {
			next := p.peekNext()
			if next.Type == kwADD {
				p.advance() // consume ','
				p.advance() // consume ADD
				p.advance() // consume ROW
				if _, err := p.expect(kwACCESS); err != nil {
					return err
				}
				if _, err := p.expect(kwPOLICY); err != nil {
					return err
				}
				if _, err := p.parseObjectName(); err != nil {
					return err
				}
				if _, err := p.expect(kwON); err != nil {
					return err
				}
				if err := p.skipParenthesized(); err != nil {
					return err
				}
			}
		}
		return nil

	case kwALL:
		// DROP ALL ROW ACCESS POLICIES
		p.advance() // consume ALL
		p.advance() // consume ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return err
		}
		if _, err := p.expect(kwPOLICIES); err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropAllRowAccessPolicies
		return nil

	case kwSEARCH:
		// DROP SEARCH OPTIMIZATION [ON search_method(target), ...]
		p.advance() // consume SEARCH
		if _, err := p.expect(kwOPTIMIZATION); err != nil {
			return err
		}
		action.Kind = ast.AlterTableDropSearchOpt
		if p.cur.Type == kwON {
			p.advance() // consume ON
			if err := p.parseSearchOptTargets(action); err != nil {
				return err
			}
		}
		return nil

	default:
		// DROP col [, col ...] — column name directly without COLUMN keyword
		if p.cur.Type == tokIdent || p.isNonReservedKeyword(p.cur.Type) {
			action.Kind = ast.AlterTableDropColumn
			return p.parseDropColumnList(action)
		}
		return p.syntaxErrorAtCur()
	}
}

// parseDropColumnList parses a comma-separated list of column names after
// DROP [COLUMN] [IF EXISTS].
func (p *Parser) parseDropColumnList(action *ast.AlterTableAction) error {
	col, err := p.parseIdent()
	if err != nil {
		return err
	}
	action.DropColumnNames = append(action.DropColumnNames, col)

	for p.cur.Type == ',' {
		next := p.peekNext()
		if !p.isColumnNameToken(next.Type) {
			break
		}
		p.advance() // consume ','
		col, err = p.parseIdent()
		if err != nil {
			return err
		}
		action.DropColumnNames = append(action.DropColumnNames, col)
	}
	return nil
}

// parseDropCascadeRestrict optionally consumes CASCADE or RESTRICT.
func (p *Parser) parseDropCascadeRestrict(action *ast.AlterTableAction) {
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		action.Cascade = true
	case kwRESTRICT:
		p.advance()
		action.Restrict = true
	}
}

// ---------------------------------------------------------------------------
// ALTER/MODIFY sub-actions
// ---------------------------------------------------------------------------

// parseAlterTableModifyAction parses the portion after ALTER/MODIFY in an
// ALTER TABLE action (column modifications, masking policies, tag ops).
func (p *Parser) parseAlterTableModifyAction(action *ast.AlterTableAction) error {
	switch p.cur.Type {
	case kwCOLUMN:
		// ALTER/MODIFY COLUMN col <specific-op>
		p.advance() // consume COLUMN
		colName, err := p.parseIdent()
		if err != nil {
			return err
		}
		return p.parseColumnSpecificOp(action, colName)

	default:
		// ALTER/MODIFY col <specific-op> or
		// ALTER/MODIFY (col_alter, col_alter, ...) or
		// ALTER/MODIFY col_alter [, col_alter ...]
		if p.cur.Type == '(' {
			// ALTER MODIFY (alter_column_clause, ...)
			p.advance() // consume '('
			if err := p.parseAlterColumnClauses(action); err != nil {
				return err
			}
			if _, err := p.expect(')'); err != nil {
				return err
			}
			return nil
		}
		// check if it looks like a column name (ident or non-reserved kw)
		if p.cur.Type == tokIdent || p.isNonReservedKeyword(p.cur.Type) {
			colName, err := p.parseIdent()
			if err != nil {
				return err
			}
			return p.parseColumnSpecificOp(action, colName)
		}
		return p.syntaxErrorAtCur()
	}
}

// parseColumnSpecificOp handles ALTER/MODIFY COLUMN col <op>.
// <op> is one of: SET MASKING POLICY, UNSET MASKING POLICY,
//
//	SET TAG, UNSET TAG, or the standard alter_column_opts.
func (p *Parser) parseColumnSpecificOp(action *ast.AlterTableAction, colName ast.Ident) error {
	switch p.cur.Type {
	case kwSET:
		next := p.peekNext()
		switch next.Type {
		case kwMASKING:
			// SET MASKING POLICY name [USING (...)] [FORCE]
			p.advance() // consume SET
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return err
			}
			policyName, err := p.parseObjectName()
			if err != nil {
				return err
			}
			if p.cur.Type == kwUSING {
				p.advance() // consume USING
				if err := p.skipParenthesized(); err != nil {
					return err
				}
			}
			if p.cur.Type == kwFORCE {
				p.advance() // consume FORCE
			}
			action.Kind = ast.AlterTableSetMaskingPolicy
			action.MaskColumn = colName
			action.MaskingPolicy = policyName
			return nil

		case kwTAG:
			// SET TAG (name = 'val', ...)
			p.advance() // consume SET
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			action.Kind = ast.AlterTableSetColumnTag
			action.TagColumn = colName
			action.Tags = tags
			return nil

		default:
			// SET DEFAULT, SET NOT NULL, SET DATA TYPE handled below
		}

	case kwUNSET:
		next := p.peekNext()
		switch next.Type {
		case kwMASKING:
			// UNSET MASKING POLICY
			p.advance() // consume UNSET
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return err
			}
			action.Kind = ast.AlterTableUnsetMaskingPolicy
			action.MaskColumn = colName
			return nil

		case kwTAG:
			// UNSET TAG (name, ...)
			p.advance() // consume UNSET
			names, err := p.parseUnsetTagList()
			if err != nil {
				return err
			}
			action.Kind = ast.AlterTableUnsetColumnTag
			action.TagColumn = colName
			action.UnsetTags = names
			return nil

		default:
			// UNSET COMMENT handled below
		}
	}

	// Standard alter_column_opts (inline, non-parenthesized form)
	ca, err := p.parseAlterColumnOpts(colName)
	if err != nil {
		return err
	}
	action.Kind = ast.AlterTableAlterColumn
	action.ColumnAlters = append(action.ColumnAlters, ca)

	// More column alters may follow, comma-separated (no parentheses)
	for p.cur.Type == ',' {
		next := p.peekNext()
		if !p.isAlterColumnStart(next.Type) {
			break
		}
		p.advance() // consume ','
		// Optional COLUMN keyword
		if p.cur.Type == kwCOLUMN {
			p.advance()
		}
		nextCol, err := p.parseIdent()
		if err != nil {
			return err
		}
		ca, err = p.parseAlterColumnOpts(nextCol)
		if err != nil {
			return err
		}
		action.ColumnAlters = append(action.ColumnAlters, ca)
	}
	return nil
}

// parseAlterColumnClauses parses ALTER/MODIFY ( col_alter, ... ) — the
// parenthesized multi-column form. Expects the opening '(' already consumed.
func (p *Parser) parseAlterColumnClauses(action *ast.AlterTableAction) error {
	action.Kind = ast.AlterTableAlterColumn

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Optional COLUMN keyword
		if p.cur.Type == kwCOLUMN {
			p.advance()
		}
		colName, err := p.parseIdent()
		if err != nil {
			return err
		}
		ca, err := p.parseAlterColumnOpts(colName)
		if err != nil {
			return err
		}
		action.ColumnAlters = append(action.ColumnAlters, ca)

		if p.cur.Type == ',' {
			p.advance() // consume ','
		} else {
			break
		}
	}
	return nil
}

// parseAlterColumnOpts parses the alter_column_opts production for one column:
//
//	DROP DEFAULT
//	SET DEFAULT obj.NEXTVAL
//	SET NOT NULL | DROP NOT NULL
//	(SET DATA)? TYPE? data_type  | just data_type
//	COMMENT 'str'
//	UNSET COMMENT
func (p *Parser) parseAlterColumnOpts(colName ast.Ident) (*ast.ColumnAlter, error) {
	ca := &ast.ColumnAlter{Column: colName}

	switch p.cur.Type {
	case kwDROP:
		p.advance() // consume DROP
		switch p.cur.Type {
		case kwDEFAULT:
			p.advance()
			ca.Kind = ast.ColumnAlterDropDefault

		case kwNOT:
			p.advance() // consume NOT
			if _, err := p.expect(kwNULL); err != nil {
				return nil, err
			}
			ca.Kind = ast.ColumnAlterDropNotNull

		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwSET:
		p.advance() // consume SET
		switch p.cur.Type {
		case kwDEFAULT:
			// SET DEFAULT obj.NEXTVAL
			p.advance() // consume DEFAULT
			// Parse obj.NEXTVAL as an expression
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			ca.Kind = ast.ColumnAlterSetDefault
			ca.DefaultExpr = expr

		case kwNOT:
			// SET NOT NULL
			p.advance() // consume NOT
			if _, err := p.expect(kwNULL); err != nil {
				return nil, err
			}
			ca.Kind = ast.ColumnAlterSetNotNull

		case kwDATA:
			// SET DATA TYPE data_type
			p.advance() // consume DATA
			if _, err := p.expect(kwTYPE); err != nil {
				return nil, err
			}
			dt, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			ca.Kind = ast.ColumnAlterSetDataType
			ca.DataType = dt

		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwNOT:
		// NOT NULL (shorthand, without SET)
		p.advance() // consume NOT
		if _, err := p.expect(kwNULL); err != nil {
			return nil, err
		}
		ca.Kind = ast.ColumnAlterSetNotNull

	case kwTYPE:
		// TYPE data_type
		p.advance() // consume TYPE
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		ca.Kind = ast.ColumnAlterSetDataType
		ca.DataType = dt

	case kwCOMMENT:
		// COMMENT 'str'
		p.advance() // consume COMMENT
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		s := tok.Str
		ca.Kind = ast.ColumnAlterSetComment
		ca.Comment = &s

	case kwUNSET:
		// UNSET COMMENT
		p.advance() // consume UNSET
		if _, err := p.expect(kwCOMMENT); err != nil {
			return nil, err
		}
		ca.Kind = ast.ColumnAlterUnsetComment

	default:
		// data_type directly (no SET DATA TYPE prefix)
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		ca.Kind = ast.ColumnAlterSetDataType
		ca.DataType = dt
	}

	return ca, nil
}

// isAlterColumnStart returns true if tok can start a column alter inside
// a multi-column ALTER MODIFY list.
func (p *Parser) isAlterColumnStart(tokType int) bool {
	switch tokType {
	case tokIdent, kwCOLUMN:
		return true
	}
	return p.isNonReservedKeyword(tokType)
}

// ---------------------------------------------------------------------------
// SET / UNSET table properties
// ---------------------------------------------------------------------------

// parseTableSetProps parses the SET properties after ALTER TABLE ... SET.
// Returns a slice of TableProp, consuming all recognized properties.
func (p *Parser) parseTableSetProps() ([]*ast.TableProp, error) {
	var props []*ast.TableProp
	for {
		prop, ok, err := p.parseOneTableProp()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		props = append(props, prop)
	}
	return props, nil
}

// parseOneTableProp attempts to parse a single known SET property.
// Returns (prop, true, nil) on success, (nil, false, nil) if no property found.
func (p *Parser) parseOneTableProp() (*ast.TableProp, bool, error) {
	switch p.cur.Type {
	case kwDATA_RETENTION_TIME_IN_DAYS:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, false, err
		}
		return &ast.TableProp{Name: "DATA_RETENTION_TIME_IN_DAYS", Value: tok.Str}, true, nil

	case kwMAX_DATA_EXTENSION_TIME_IN_DAYS:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, false, err
		}
		return &ast.TableProp{Name: "MAX_DATA_EXTENSION_TIME_IN_DAYS", Value: tok.Str}, true, nil

	case kwCHANGE_TRACKING:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == kwTRUE || p.cur.Type == kwFALSE {
			tok := p.advance()
			return &ast.TableProp{Name: "CHANGE_TRACKING", Value: tok.Str}, true, nil
		}
		return nil, false, p.syntaxErrorAtCur()

	case kwDEFAULT_DDL_COLLATION:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, false, err
		}
		return &ast.TableProp{Name: "DEFAULT_DDL_COLLATION", Value: tok.Str}, true, nil

	case kwCOMMENT:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, false, err
		}
		return &ast.TableProp{Name: "COMMENT", Value: tok.Str}, true, nil

	case kwSTAGE_FILE_FORMAT:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if err := p.skipParenthesized(); err != nil {
			return nil, false, err
		}
		return &ast.TableProp{Name: "STAGE_FILE_FORMAT", Value: "(...)"}, true, nil

	case kwSTAGE_COPY_OPTIONS:
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if err := p.skipParenthesized(); err != nil {
			return nil, false, err
		}
		return &ast.TableProp{Name: "STAGE_COPY_OPTIONS", Value: "(...)"}, true, nil

	default:
		return nil, false, nil
	}
}

// parseTableUnsetProps parses the UNSET property names (one or more, comma-sep).
func (p *Parser) parseTableUnsetProps() ([]string, error) {
	name, err := p.parseTableUnsetPropName()
	if err != nil {
		return nil, err
	}
	names := []string{name}
	for p.cur.Type == ',' {
		next := p.peekNext()
		if !p.isTableUnsetPropToken(next.Type) {
			break
		}
		p.advance() // consume ','
		name, err = p.parseTableUnsetPropName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

// parseTableUnsetPropName parses a single UNSET property name keyword.
func (p *Parser) parseTableUnsetPropName() (string, error) {
	switch p.cur.Type {
	case kwDATA_RETENTION_TIME_IN_DAYS:
		p.advance()
		return "DATA_RETENTION_TIME_IN_DAYS", nil
	case kwMAX_DATA_EXTENSION_TIME_IN_DAYS:
		p.advance()
		return "MAX_DATA_EXTENSION_TIME_IN_DAYS", nil
	case kwCHANGE_TRACKING:
		p.advance()
		return "CHANGE_TRACKING", nil
	case kwDEFAULT_DDL_COLLATION:
		p.advance()
		return "DEFAULT_DDL_COLLATION", nil
	case kwCOMMENT:
		p.advance()
		return "COMMENT", nil
	case tokIdent:
		tok := p.advance()
		return tok.Str, nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// isTableUnsetPropToken returns true if tok can be a TABLE UNSET prop name.
func (p *Parser) isTableUnsetPropToken(tokType int) bool {
	switch tokType {
	case kwDATA_RETENTION_TIME_IN_DAYS, kwMAX_DATA_EXTENSION_TIME_IN_DAYS,
		kwCHANGE_TRACKING, kwDEFAULT_DDL_COLLATION, kwCOMMENT, tokIdent:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// isColumnNameToken returns true if tok could be a column name.
func (p *Parser) isColumnNameToken(tokType int) bool {
	if tokType == tokIdent {
		return true
	}
	return p.isNonReservedKeyword(tokType)
}

// isNonReservedKeyword returns true if a keyword token is non-reserved
// (usable as an identifier without quoting) in Snowflake. This is used to
// disambiguate comma-separated lists where the next token could be a column
// name that happens to be a keyword.
func (p *Parser) isNonReservedKeyword(tokType int) bool {
	// keyword tokens start at kwABORT (700); anything below is not a keyword
	if tokType < kwABORT {
		return false
	}
	return !keywordReserved[tokType]
}

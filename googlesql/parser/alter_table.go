package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Core DDL — ALTER (parser-ddl node)
// ---------------------------------------------------------------------------
//
// Ports the table-like-object alternatives of the legacy ANTLR alter_statement
// + the alter_action union (GoogleSQLParser.g4 §2.3/§2.5):
//
//	alter_statement:
//	    ALTER table_or_table_function opt_if_exists? maybe_dashed_path alter_action_list
//	  | ALTER schema_object_kind     opt_if_exists? path_expression   alter_action_list
//	  | … (generic entity / privilege restriction / row-access-policy — not owned here)
//	alter_action_list: alter_action (, alter_action)*
//
// This node owns ALTER for TABLE / VIEW / INDEX / SCHEMA / DATABASE and the
// alter_action variants those objects use (SET OPTIONS, ADD/DROP/RENAME/ALTER
// COLUMN, ADD/DROP/ALTER CONSTRAINT, DROP PRIMARY KEY, RENAME TO, SET DEFAULT
// COLLATE, the ROW DELETION POLICY actions, the Spanner SET ON DELETE and
// ALTER COLUMN actions, and ALTER INDEX ADD/DROP STORED COLUMN). The generic-
// entity / row-access-policy / privilege-restriction alternatives belong to the
// dialect-specific DDL nodes.
//
// ORACLE NOTE — Spanner SUBSET (oracle.md). The Spanner emulator accepts the
// Spanner action set (ADD/DROP/ALTER COLUMN single-action, ADD/DROP CONSTRAINT,
// SET ON DELETE, the ROW DELETION POLICY actions, RENAME TO) and REJECTS the
// BigQuery-only actions (SET OPTIONS on a table, comma-separated multi-action,
// ALTER COLUMN SET DATA TYPE, RENAME COLUMN, ALTER COLUMN DROP NOT NULL, and
// `ALTER COLUMN … SET NOT NULL` which is not a grammar form at all) — those
// rejects are NON-AUTHORITATIVE for the union and are triangulated from the .g4
// + BigQuery truth1. `ALTER DATABASE … SET OPTIONS` accepts when the option NAME
// is one the emulator knows; an unknown name is an over-reject (oracle.md), so
// the parser accepts arbitrary option names. Verified in alter_table_oracle_test.go.

// parseAlterStmt parses an ALTER statement. The ALTER keyword has NOT yet been
// consumed (parseStmt peeks it). It dispatches on the object keyword.
func (p *Parser) parseAlterStmt() (ast.Node, error) {
	alter := p.advance() // consume ALTER

	var kind ast.AlterObjectKind
	switch p.cur.Type {
	case kwTABLE:
		// ALTER TABLE — but ALTER TABLE FUNCTION is dialect-node territory.
		if p.peekNext().Type == kwFUNCTION {
			return p.unsupported("ALTER TABLE FUNCTION")
		}
		kind = ast.AlterTable
	case kwVIEW:
		kind = ast.AlterView
	case kwINDEX:
		kind = ast.AlterIndex
	case kwSCHEMA:
		kind = ast.AlterSchema
	case kwDATABASE:
		kind = ast.AlterDatabase
	// --- BigQuery-only ALTER objects (parser-ddl-bigquery node) ---
	case kwMATERIALIZED, kwAPPROX:
		// ALTER MATERIALIZED|APPROX VIEW SET OPTIONS (bq_materialized_view.go).
		return p.parseBQAlterMaterializedView(alter)
	case kwVECTOR:
		// ALTER VECTOR INDEX … ON … REBUILD (bq_search_vector_index.go).
		return p.parseBQAlterVectorIndex(alter)
	default:
		// ALTER <generic-entity> (CAPACITY / RESERVATION / ASSIGNMENT — the keyword
		// lexes as an identifier) routes to the generic-entity alter
		// (bq_capacity.go). ALTER FUNCTION / PROCEDURE / MODEL / PRIVILEGE
		// RESTRICTION / ROW ACCESS POLICY / ALL ROW ACCESS POLICIES are not owned
		// here.
		if isGenericEntityType(p.cur.Type) {
			return p.parseBQAlterEntity(alter)
		}
		return p.unsupported("ALTER")
	}
	p.advance() // consume the object keyword

	stmt := &ast.AlterStmt{Object: kind}
	stmt.Loc.Start = alter.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// alter_action_list: alter_action (, alter_action)* — requires at least one.
	for {
		action, err := p.parseAlterAction()
		if err != nil {
			return nil, err
		}
		stmt.Actions = append(stmt.Actions, action)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterAction parses one alter_action, dispatching on its leading keyword.
func (p *Parser) parseAlterAction() (*ast.AlterAction, error) {
	start := p.cur.Loc
	switch p.cur.Type {
	case kwSET:
		return p.parseAlterSetAction(start)
	case kwADD:
		return p.parseAlterAddAction(start)
	case kwDROP:
		return p.parseAlterDropAction(start)
	case kwALTER:
		return p.parseAlterAlterAction(start)
	case kwRENAME:
		return p.parseAlterRenameAction(start)
	case kwREPLACE:
		return p.parseAlterReplaceAction(start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterSetAction parses the SET-led alter_action forms:
//
//	SET OPTIONS options_list
//	SET DEFAULT collate_clause
//	SET ON DELETE foreign_key_action   (Spanner spanner_set_on_delete_action)
func (p *Parser) parseAlterSetAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // SET
	switch p.cur.Type {
	case kwOPTIONS:
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterSetOptions, Options: opts, Loc: locTo(start, p.prev)}, nil
	case kwDEFAULT:
		p.advance() // DEFAULT
		coll, _, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterSetDefaultCollate, Collate: coll, Loc: locTo(start, p.prev)}, nil
	case kwON:
		p.advance() // ON
		if _, err := p.expect(kwDELETE); err != nil {
			return nil, err
		}
		act, err := p.parseForeignKeyAction()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterSetOnDelete, OnDelete: act, Loc: locTo(start, p.prev)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterAddAction parses the ADD-led alter_action forms:
//
//	ADD table_constraint_spec                                  (FK | CHECK)
//	ADD primary_key_spec                                       (PRIMARY KEY)
//	ADD CONSTRAINT opt_if_not_exists? identifier primary_key_or_table_constraint_spec
//	ADD COLUMN opt_if_not_exists? table_column_definition column_position? fill_using_expression?
//	ADD ROW DELETION POLICY opt_if_not_exists? '(' expression ')'
//	ADD STORED COLUMN identifier                               (ALTER INDEX)
func (p *Parser) parseAlterAddAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // ADD
	switch p.cur.Type {
	case kwCOLUMN:
		return p.parseAddColumnAction(start)
	case kwCONSTRAINT:
		return p.parseAddConstraintNamedAction(start)
	case kwPRIMARY:
		con, err := p.parsePrimaryKeySpec(p.cur.Loc)
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAddConstraint, Constraint: con, Loc: locTo(start, p.prev)}, nil
	case kwFOREIGN:
		con, err := p.parseForeignKeyConstraintSpec(p.cur.Loc)
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAddConstraint, Constraint: con, Loc: locTo(start, p.prev)}, nil
	case kwCHECK:
		con, err := p.parseCheckConstraintSpec(p.cur.Loc)
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAddConstraint, Constraint: con, Loc: locTo(start, p.prev)}, nil
	case kwROW:
		expr, ifne, err := p.parseAddRowDeletionPolicy()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAddRowDeletionPolicy, RowDeletion: expr, IfNotExists: ifne, Loc: locTo(start, p.prev)}, nil
	case kwSTORED:
		// ALTER INDEX … ADD STORED COLUMN identifier.
		name, err := p.parseStoredColumn()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAddStoredColumn, StoredColumn: name, Loc: locTo(start, p.prev)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAddColumnAction parses `ADD COLUMN opt_if_not_exists?
// table_column_definition column_position? fill_using_expression?` (the ADD
// keyword and COLUMN are at/after cur). The COLUMN keyword is the current token.
func (p *Parser) parseAddColumnAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // COLUMN
	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	col, err := p.parseColumnDefinition()
	if err != nil {
		return nil, err
	}
	action := &ast.AlterAction{Kind: ast.AlterAddColumn, IfNotExists: ifNotExists, Column: col}

	// column_position?  — PRECEDING id | FOLLOWING id.
	switch p.cur.Type {
	case kwPRECEDING:
		p.advance()
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		action.Position = "PRECEDING " + tok.Str
	case kwFOLLOWING:
		p.advance()
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		action.Position = "FOLLOWING " + tok.Str
	}

	// fill_using_expression?  — FILL USING expr.
	if p.cur.Type == kwFILL {
		p.advance() // FILL
		if _, err := p.expect(kwUSING); err != nil {
			return nil, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		action.FillUsing = expr
	}

	action.Loc = locTo(start, p.prev)
	return action, nil
}

// parseAddConstraintNamedAction parses `ADD CONSTRAINT opt_if_not_exists?
// identifier primary_key_or_table_constraint_spec` (the CONSTRAINT keyword is
// the current token).
func (p *Parser) parseAddConstraintNamedAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // CONSTRAINT
	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	nameTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	// primary_key_or_table_constraint_spec.
	var con *ast.TableConstraint
	switch p.cur.Type {
	case kwPRIMARY:
		con, err = p.parsePrimaryKeySpec(p.cur.Loc)
	case kwFOREIGN:
		con, err = p.parseForeignKeyConstraintSpec(p.cur.Loc)
	case kwCHECK:
		con, err = p.parseCheckConstraintSpec(p.cur.Loc)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	con.Name = nameTok.Str
	return &ast.AlterAction{
		Kind:           ast.AlterAddConstraint,
		IfNotExists:    ifNotExists,
		ConstraintName: nameTok.Str,
		Constraint:     con,
		Loc:            locTo(start, p.prev),
	}, nil
}

// parseAlterDropAction parses the DROP-led alter_action forms:
//
//	DROP CONSTRAINT opt_if_exists? identifier
//	DROP PRIMARY KEY opt_if_exists?
//	DROP COLUMN opt_if_exists? identifier
//	DROP ROW DELETION POLICY opt_if_exists?
//	DROP STORED COLUMN identifier   (ALTER INDEX)
func (p *Parser) parseAlterDropAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // DROP
	switch p.cur.Type {
	case kwCONSTRAINT:
		p.advance() // CONSTRAINT
		ifExists, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterDropConstraint, IfExists: ifExists, ConstraintName: nameTok.Str, Loc: locTo(start, p.prev)}, nil
	case kwPRIMARY:
		p.advance() // PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		ifExists, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterDropPrimaryKey, IfExists: ifExists, Loc: locTo(start, p.prev)}, nil
	case kwCOLUMN:
		p.advance() // COLUMN
		ifExists, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(nameTok)
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterDropColumn, IfExists: ifExists, ColumnName: name, Loc: locTo(start, p.prev)}, nil
	case kwROW:
		// DROP ROW DELETION POLICY opt_if_exists?.
		p.advance() // ROW
		if _, err := p.expect(kwDELETION); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		ifExists, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterDropRowDeletionPolicy, IfExists: ifExists, Loc: locTo(start, p.prev)}, nil
	case kwSTORED:
		name, err := p.parseStoredColumn()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterDropStoredColumn, StoredColumn: name, Loc: locTo(start, p.prev)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterAlterAction parses the ALTER-led alter_action forms:
//
//	ALTER CONSTRAINT opt_if_exists? identifier constraint_enforcement
//	ALTER CONSTRAINT opt_if_exists? identifier SET OPTIONS options_list
//	ALTER COLUMN  …  (BigQuery SET DATA TYPE / SET OPTIONS / SET DEFAULT / DROP
//	                  DEFAULT / DROP NOT NULL / DROP GENERATED, OR the Spanner
//	                  ALTER COLUMN <schema> form)
func (p *Parser) parseAlterAlterAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // ALTER
	switch p.cur.Type {
	case kwCONSTRAINT:
		return p.parseAlterConstraintAction(start)
	case kwCOLUMN:
		return p.parseAlterColumnAction(start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterConstraintAction parses `ALTER CONSTRAINT opt_if_exists? identifier
// {constraint_enforcement | SET OPTIONS options_list}` (CONSTRAINT is cur).
func (p *Parser) parseAlterConstraintAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // CONSTRAINT
	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	nameTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == kwSET {
		p.advance() // SET
		if _, err := p.expect(kwOPTIONS); err == nil {
			// Already consumed OPTIONS via expect; parse the options_list body.
			// Re-enter parseOptionsList expects the OPTIONS keyword, so build the
			// list inline here instead.
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			opts, err := p.parseOptionsListBody()
			if err != nil {
				return nil, err
			}
			return &ast.AlterAction{Kind: ast.AlterAlterConstraintOptions, IfExists: ifExists, ConstraintName: nameTok.Str, Options: opts, Loc: locTo(start, p.prev)}, nil
		}
		return nil, p.syntaxErrorAtCur()
	}
	// constraint_enforcement: [NOT] ENFORCED (required in this alt).
	enf, ok := p.tryParseEnforcement()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	return &ast.AlterAction{Kind: ast.AlterAlterConstraintEnforce, IfExists: ifExists, ConstraintName: nameTok.Str, Enforced: enf, Loc: locTo(start, p.prev)}, nil
}

// parseAlterColumnAction parses the `ALTER COLUMN …` alter_action forms (COLUMN
// is cur), covering both the BigQuery `opt_if_exists? identifier {SET DATA TYPE
// field_schema | SET OPTIONS … | SET DEFAULT expr | DROP DEFAULT | DROP NOT NULL
// | DROP GENERATED}` forms and the Spanner spanner_alter_column_action
// `opt_if_exists? identifier column_schema_inner not_null? generated/default?
// OPTIONS?`.
//
// Disambiguation: after `ALTER COLUMN [IF EXISTS] <name>`, a leading SET/DROP
// keyword selects a BigQuery action; anything else (a type token) is the Spanner
// in-place column redefinition.
func (p *Parser) parseAlterColumnAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // COLUMN
	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	nameTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwSET:
		return p.parseAlterColumnSet(start, ifExists, name)
	case kwDROP:
		return p.parseAlterColumnDrop(start, ifExists, name)
	default:
		// Spanner spanner_alter_column_action: <name> column_schema_inner
		// not_null_column_attribute? spanner_generated_or_default? opt_options_list?.
		return p.parseSpannerAlterColumn(start, ifExists, name)
	}
}

// parseAlterColumnSet parses the `SET …` tail of an ALTER COLUMN action (SET is
// cur): SET DATA TYPE field_schema | SET OPTIONS options_list | SET DEFAULT expr.
func (p *Parser) parseAlterColumnSet(start ast.Loc, ifExists bool, name string) (*ast.AlterAction, error) {
	p.advance() // SET
	switch p.cur.Type {
	case kwDATA:
		p.advance() // DATA
		if _, err := p.expect(kwTYPE); err != nil {
			return nil, err
		}
		// field_schema: column_schema_inner collate_clause? not_null? OPTIONS? —
		// the core captures the type (and a trailing NOT NULL/collate are absorbed
		// by parseType's collate handling + the optional NOT NULL below).
		dt, err := p.parseType()
		if err != nil {
			return nil, err
		}
		action := &ast.AlterAction{Kind: ast.AlterAlterColumnType, IfExists: ifExists, ColumnName: name, NewType: &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}}
		if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
			p.advance()
			p.advance()
			action.NotNull = true
		}
		action.Loc = locTo(start, p.prev)
		return action, nil
	case kwOPTIONS:
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAlterColumnOptions, IfExists: ifExists, ColumnName: name, Options: opts, Loc: locTo(start, p.prev)}, nil
	case kwDEFAULT:
		p.advance() // DEFAULT
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAlterColumnSetDefault, IfExists: ifExists, ColumnName: name, Default: expr, Loc: locTo(start, p.prev)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterColumnDrop parses the `DROP …` tail of an ALTER COLUMN action (DROP
// is cur): DROP DEFAULT | DROP NOT NULL | DROP GENERATED.
func (p *Parser) parseAlterColumnDrop(start ast.Loc, ifExists bool, name string) (*ast.AlterAction, error) {
	p.advance() // DROP
	switch p.cur.Type {
	case kwDEFAULT:
		p.advance()
		return &ast.AlterAction{Kind: ast.AlterAlterColumnDropDefault, IfExists: ifExists, ColumnName: name, Loc: locTo(start, p.prev)}, nil
	case kwNOT:
		p.advance() // NOT
		if _, err := p.expect(kwNULL); err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterAlterColumnDropNotNull, IfExists: ifExists, ColumnName: name, Loc: locTo(start, p.prev)}, nil
	case kwGENERATED:
		p.advance()
		return &ast.AlterAction{Kind: ast.AlterAlterColumnDropGenerated, IfExists: ifExists, ColumnName: name, Loc: locTo(start, p.prev)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseSpannerAlterColumn parses the Spanner spanner_alter_column_action tail
// after `ALTER COLUMN [IF EXISTS] <name>`: `column_schema_inner
// not_null_column_attribute? spanner_generated_or_default? opt_options_list?`.
// cur is at the type token.
func (p *Parser) parseSpannerAlterColumn(start ast.Loc, ifExists bool, name string) (*ast.AlterAction, error) {
	dt, err := p.parseType()
	if err != nil {
		return nil, err
	}
	action := &ast.AlterAction{
		Kind:       ast.AlterSpannerAlterColumn,
		IfExists:   ifExists,
		ColumnName: name,
		NewType:    &ast.TypeRef{Text: dt.String(), Loc: dt.Loc},
	}

	// not_null_column_attribute?  — NOT NULL.
	if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
		p.advance()
		p.advance()
		action.NotNull = true
	}

	// spanner_generated_or_default?  — AS '(' expression ')' STORED.
	if p.cur.Type == kwAS {
		p.advance() // AS
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSTORED); err != nil {
			return nil, err
		}
		action.Generated = expr
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		action.Options = opts
	}

	action.Loc = locTo(start, p.prev)
	return action, nil
}

// parseAlterRenameAction parses the RENAME-led alter_action forms:
//
//	RENAME TO path_expression
//	RENAME COLUMN opt_if_exists? identifier TO identifier
func (p *Parser) parseAlterRenameAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // RENAME
	switch p.cur.Type {
	case kwTO:
		p.advance() // TO
		path, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterRenameTo, RenameTo: path, Loc: locTo(start, p.prev)}, nil
	case kwCOLUMN:
		p.advance() // COLUMN
		ifExists, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		fromTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		from, err := p.identifierText(fromTok)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		toTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		to, err := p.identifierText(toTok)
		if err != nil {
			return nil, err
		}
		return &ast.AlterAction{Kind: ast.AlterRenameColumn, IfExists: ifExists, ColumnName: from, NewName: to, Loc: locTo(start, p.prev)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterReplaceAction parses the REPLACE-led alter_action form:
//
//	REPLACE ROW DELETION POLICY opt_if_exists? '(' expression ')'
func (p *Parser) parseAlterReplaceAction(start ast.Loc) (*ast.AlterAction, error) {
	p.advance() // REPLACE
	if _, err := p.expect(kwROW); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDELETION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}
	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return &ast.AlterAction{Kind: ast.AlterReplaceRowDeletionPolicy, IfExists: ifExists, RowDeletion: expr, Loc: locTo(start, p.prev)}, nil
}

// parseAddRowDeletionPolicy parses `ROW DELETION POLICY opt_if_not_exists? '('
// expression ')'` (ROW is cur). Returns the policy expression and the IF NOT
// EXISTS flag.
func (p *Parser) parseAddRowDeletionPolicy() (ast.Node, bool, error) {
	p.advance() // ROW
	if _, err := p.expect(kwDELETION); err != nil {
		return nil, false, err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, false, err
	}
	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, false, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, false, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, false, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, false, err
	}
	return expr, ifNotExists, nil
}

// parseStoredColumn parses `STORED COLUMN identifier` (STORED is cur) for the
// ALTER INDEX ADD/DROP STORED COLUMN actions. Returns the column name.
func (p *Parser) parseStoredColumn() (string, error) {
	p.advance() // STORED
	if _, err := p.expect(kwCOLUMN); err != nil {
		return "", err
	}
	tok, err := p.expectIdentifier()
	if err != nil {
		return "", err
	}
	return p.identifierText(tok)
}

// parseOptionsListBody parses the inner `(options_entry (, options_entry)*)? )`
// of an options_list when the OPTIONS keyword and the opening '(' have already
// been consumed. (Used by the ALTER CONSTRAINT SET OPTIONS path, which consumes
// OPTIONS before discovering it needs the list.)
func (p *Parser) parseOptionsListBody() ([]*ast.OptionsEntry, error) {
	var entries []*ast.OptionsEntry
	if p.cur.Type == int(')') {
		p.advance()
		return entries, nil
	}
	for {
		entry, err := p.parseOptionsEntry()
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return entries, nil
}

// locTo builds a Loc from a start Loc through the end of an end token.
func locTo(start ast.Loc, endTok Token) ast.Loc {
	return ast.Loc{Start: start.Start, End: endTok.Loc.End}
}

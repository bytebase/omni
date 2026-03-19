package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseAlterTableStmt parses an ALTER TABLE statement and its variants
// (ALTER INDEX, ALTER SEQUENCE, ALTER VIEW, ALTER MATERIALIZED VIEW, ALTER FOREIGN TABLE).
//
// Ref: https://www.postgresql.org/docs/17/sql-altertable.html
//
//	ALTER TABLE [ IF EXISTS ] [ ONLY ] name [ * ]
//	    action [, ... ]
//	ALTER TABLE [ IF EXISTS ] [ ONLY ] name [ * ]
//	    RENAME [ COLUMN ] column_name TO new_column_name
//	ALTER TABLE [ IF EXISTS ] [ ONLY ] name [ * ]
//	    RENAME CONSTRAINT constraint_name TO new_constraint_name
//	ALTER TABLE [ IF EXISTS ] name
//	    RENAME TO new_name
//	ALTER TABLE ALL IN TABLESPACE name [ OWNED BY role_name [, ... ] ]
//	    SET TABLESPACE new_tablespace [ NOWAIT ]
func (p *Parser) parseAlterTableStmt() nodes.Node {
	// Already consumed ALTER; current token is TABLE, INDEX, SEQUENCE, VIEW, MATERIALIZED, or FOREIGN.
	switch p.cur.Type {
	case TABLE:
		return p.parseAlterTable()
	case INDEX:
		return p.parseAlterIndex()
	case SEQUENCE:
		return p.parseAlterSequence()
	case VIEW:
		return p.parseAlterView()
	case MATERIALIZED:
		return p.parseAlterMaterializedView()
	case FOREIGN:
		return p.parseAlterForeignTable()
	case EVENT:
		p.advance() // consume EVENT
		n, _ := p.parseAlterEventTrigStmt()
		return n
	case EXTENSION:
		return p.parseAlterExtensionStmt()
	default:
		return nil
	}
}

// parseAlterTable handles ALTER TABLE ...
func (p *Parser) parseAlterTable() nodes.Node {
	p.advance() // consume TABLE

	// ALTER TABLE ALL IN TABLESPACE ...
	if p.cur.Type == ALL {
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_TABLE))
	}

	// IF EXISTS
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}

	// relation_expr
	rel, _ := p.parseRelationExpr()

	// Check for RENAME (produces RenameStmt, not AlterTableStmt)
	if p.cur.Type == RENAME {
		return p.parseAlterTableRename(rel, missingOk)
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_TABLE,
			Relation:   rel,
			Newschema:  newschema,
			MissingOk:  missingOk,
		}
	}

	// alter_table_cmds
	cmds := p.parseAlterTableCmds()
	return &nodes.AlterTableStmt{
		Relation:   rel,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_TABLE),
		Missing_ok: missingOk,
	}
}

// parseAlterIndex handles ALTER INDEX ...
func (p *Parser) parseAlterIndex() nodes.Node {
	p.advance() // consume INDEX

	// ALTER INDEX ALL IN TABLESPACE ...
	if p.cur.Type == ALL {
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_INDEX))
	}

	// IF EXISTS
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}

	// qualified_name
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	// Check for RENAME
	if p.cur.Type == RENAME {
		p.advance() // consume RENAME
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_INDEX,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
		}
	}

	// ALTER INDEX name ATTACH PARTITION name
	if p.cur.Type == ATTACH {
		p.advance() // ATTACH
		p.expect(PARTITION)
		partNames, _ := p.parseQualifiedName()
		partRv := makeRangeVarFromAnyName(partNames)
		cmd := &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AttachPartition),
			Def: &nodes.PartitionCmd{
				Name: partRv,
			},
		}
		return &nodes.AlterTableStmt{
			Relation:   rv,
			Cmds:       &nodes.List{Items: []nodes.Node{cmd}},
			ObjType:    int(nodes.OBJECT_INDEX),
			Missing_ok: missingOk,
		}
	}

	// ALTER INDEX name [NO] DEPENDS ON EXTENSION ext_name
	if p.cur.Type == DEPENDS || (p.cur.Type == NO && p.peekNext().Type == DEPENDS) {
		remove := p.parseOptNo()
		p.advance() // consume DEPENDS
		p.expect(ON)
		p.expect(EXTENSION)
		extname, _ := p.parseName()
		return &nodes.AlterObjectDependsStmt{
			ObjectType: nodes.OBJECT_INDEX,
			Relation:   rv,
			Extname:    &nodes.String{Str: extname},
			Remove:     remove,
		}
	}

	cmds := p.parseAlterTableCmds()
	return &nodes.AlterTableStmt{
		Relation:   rv,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_INDEX),
		Missing_ok: missingOk,
	}
}

// parseAlterSequence handles ALTER SEQUENCE ...
func (p *Parser) parseAlterSequence() nodes.Node {
	p.advance() // consume SEQUENCE

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}

	names, _ := p.parseQualifiedName()

	// Check for RENAME (produces RenameStmt, uses makeRangeVarFromAnyName)
	if p.cur.Type == RENAME {
		rv := makeRangeVarFromAnyName(names)
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_SEQUENCE,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
		}
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		rv := makeRangeVarFromAnyName(names)
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_SEQUENCE,
			Relation:   rv,
			Newschema:  newschema,
			MissingOk:  missingOk,
		}
	}

	// OWNER TO goes through alter_table_cmds in the yacc grammar
	if p.cur.Type == OWNER {
		rv := makeRangeVarFromAnyName(names)
		cmds := p.parseAlterTableCmds()
		return &nodes.AlterTableStmt{
			Relation:   rv,
			Cmds:       cmds,
			ObjType:    int(nodes.OBJECT_SEQUENCE),
			Missing_ok: missingOk,
		}
	}

	// AlterSeqStmt: ALTER SEQUENCE name SeqOptList
	rv := makeRangeVarFromNames(names)
	options := p.parseSeqOptList()
	return &nodes.AlterSeqStmt{
		Sequence:  rv,
		Options:   options,
		MissingOk: missingOk,
	}
}

// parseAlterView handles ALTER VIEW ...
func (p *Parser) parseAlterView() nodes.Node {
	p.advance() // consume VIEW

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}

	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	// Check for RENAME
	if p.cur.Type == RENAME {
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_VIEW,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
		}
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_VIEW,
			Relation:   rv,
			Newschema:  newschema,
			MissingOk:  missingOk,
		}
	}

	cmds := p.parseAlterTableCmds()
	return &nodes.AlterTableStmt{
		Relation:   rv,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_VIEW),
		Missing_ok: missingOk,
	}
}

// parseAlterMaterializedView handles ALTER MATERIALIZED VIEW ...
func (p *Parser) parseAlterMaterializedView() nodes.Node {
	p.advance() // consume MATERIALIZED
	p.expect(VIEW)

	// ALTER MATERIALIZED VIEW ALL IN TABLESPACE ...
	if p.cur.Type == ALL {
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_MATVIEW))
	}

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}

	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	// Check for RENAME
	if p.cur.Type == RENAME {
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_MATVIEW,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
		}
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_MATVIEW,
			Relation:   rv,
			Newschema:  newschema,
			MissingOk:  missingOk,
		}
	}

	// ALTER MATERIALIZED VIEW name [NO] DEPENDS ON EXTENSION ext_name
	if p.cur.Type == DEPENDS || (p.cur.Type == NO && p.peekNext().Type == DEPENDS) {
		remove := p.parseOptNo()
		p.advance() // consume DEPENDS
		p.expect(ON)
		p.expect(EXTENSION)
		extname, _ := p.parseName()
		return &nodes.AlterObjectDependsStmt{
			ObjectType: nodes.OBJECT_MATVIEW,
			Relation:   rv,
			Extname:    &nodes.String{Str: extname},
			Remove:     remove,
		}
	}

	cmds := p.parseAlterTableCmds()
	return &nodes.AlterTableStmt{
		Relation:   rv,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_MATVIEW),
		Missing_ok: missingOk,
	}
}

// parseAlterForeignTable handles ALTER FOREIGN TABLE ... and ALTER FOREIGN DATA WRAPPER ...
func (p *Parser) parseAlterForeignTable() nodes.Node {
	p.advance() // consume FOREIGN

	// ALTER FOREIGN DATA WRAPPER ...
	if p.cur.Type == DATA_P {
		return p.parseAlterFdwStmt()
	}

	p.expect(TABLE)

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}

	rel, _ := p.parseRelationExpr()

	// Check for RENAME
	if p.cur.Type == RENAME {
		return p.parseAlterTableRename(rel, missingOk)
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_FOREIGN_TABLE,
			Relation:   rel,
			Newschema:  newschema,
			MissingOk:  missingOk,
		}
	}

	cmds := p.parseAlterTableCmds()
	return &nodes.AlterTableStmt{
		Relation:   rel,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_FOREIGN_TABLE),
		Missing_ok: missingOk,
	}
}

// parseAlterTableMoveAll parses ALTER TABLE/INDEX/MATERIALIZED VIEW ALL IN TABLESPACE ...
//
//	ALTER TABLE ALL IN TABLESPACE name [OWNED BY role_list] SET TABLESPACE name [NOWAIT]
func (p *Parser) parseAlterTableMoveAll(objType int) *nodes.AlterTableMoveAllStmt {
	p.advance() // consume ALL
	p.expect(IN_P)
	p.expect(TABLESPACE)
	origTs, _ := p.parseName()

	var roles *nodes.List
	if p.cur.Type == OWNED {
		p.advance() // OWNED
		p.expect(BY)
		roles = p.parseRoleList()
	}

	p.expect(SET)
	p.expect(TABLESPACE)
	newTs, _ := p.parseName()

	nowait := false
	if p.cur.Type == NOWAIT {
		p.advance()
		nowait = true
	}

	return &nodes.AlterTableMoveAllStmt{
		OrigTablespacename: origTs,
		ObjType:            objType,
		Roles:              roles,
		NewTablespacename:  newTs,
		Nowait:             nowait,
	}
}

// parseAlterTableRename parses ALTER TABLE ... RENAME ...
func (p *Parser) parseAlterTableRename(rel *nodes.RangeVar, missingOk bool) *nodes.RenameStmt {
	p.advance() // consume RENAME

	switch p.cur.Type {
	case TO:
		// RENAME TO name
		p.advance()
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TABLE,
			Relation:   rel,
			Newname:    newname,
			MissingOk:  missingOk,
		}
	case COLUMN:
		// RENAME COLUMN colname TO newname
		p.advance() // consume COLUMN
		oldname, _ := p.parseColId()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType:   nodes.OBJECT_COLUMN,
			RelationType: nodes.OBJECT_TABLE,
			Relation:     rel,
			Subname:      oldname,
			Newname:      newname,
			MissingOk:    missingOk,
		}
	case CONSTRAINT:
		// RENAME CONSTRAINT oldname TO newname
		p.advance() // consume CONSTRAINT
		oldname, _ := p.parseName()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType:   nodes.OBJECT_TABCONSTRAINT,
			RelationType: nodes.OBJECT_TABLE,
			Relation:     rel,
			Subname:      oldname,
			Newname:      newname,
			MissingOk:    missingOk,
		}
	default:
		// RENAME colname TO newname (implicit column rename, no COLUMN keyword)
		oldname, _ := p.parseColId()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType:   nodes.OBJECT_COLUMN,
			RelationType: nodes.OBJECT_TABLE,
			Relation:     rel,
			Subname:      oldname,
			Newname:      newname,
			MissingOk:    missingOk,
		}
	}
}

// parseAlterTableCmds parses a comma-separated list of alter_table_cmd.
//
//	alter_table_cmds: alter_table_cmd | alter_table_cmds ',' alter_table_cmd
func (p *Parser) parseAlterTableCmds() *nodes.List {
	cmd := p.parseAlterTableCmd()
	items := []nodes.Node{cmd}
	for p.cur.Type == ',' {
		p.advance()
		cmd = p.parseAlterTableCmd()
		items = append(items, cmd)
	}
	return &nodes.List{Items: items}
}

// parseAlterTableCmd parses a single alter_table_cmd.
func (p *Parser) parseAlterTableCmd() *nodes.AlterTableCmd {
	if p.collectMode() {
		p.cachedCollect("parseAlterTableCmd", func() {
			for _, t := range []int{
				ADD_P, DROP, ALTER, OWNER, VALIDATE,
				INHERIT, NO, ATTACH, DETACH,
				ENABLE_P, DISABLE_P, FORCE, CLUSTER,
				SET, RESET, REPLICA, OF, NOT, OPTIONS,
				RENAME,
			} {
				p.addTokenCandidate(t)
			}
		})
		return nil
	}
	switch p.cur.Type {
	case ADD_P:
		return p.parseAlterTableAdd()
	case DROP:
		return p.parseAlterTableDrop()
	case ALTER:
		return p.parseAlterTableAlter()
	case OWNER:
		return p.parseAlterTableOwner()
	case VALIDATE:
		return p.parseAlterTableValidate()
	case INHERIT:
		return p.parseAlterTableInherit()
	case NO:
		return p.parseAlterTableNo()
	case ATTACH:
		return p.parseAlterTableAttach()
	case DETACH:
		return p.parseAlterTableDetach()
	case ENABLE_P:
		return p.parseAlterTableEnable()
	case DISABLE_P:
		return p.parseAlterTableDisable()
	case FORCE:
		return p.parseAlterTableForce()
	case CLUSTER:
		return p.parseAlterTableCluster()
	case SET:
		return p.parseAlterTableSet()
	case RESET:
		return p.parseAlterTableReset()
	case REPLICA:
		return p.parseAlterTableReplica()
	case OF:
		return p.parseAlterTableOf()
	case NOT:
		return p.parseAlterTableNot()
	case OPTIONS:
		return p.parseAlterTableOptions()
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterTableAdd handles ADD ... subcommands.
func (p *Parser) parseAlterTableAdd() *nodes.AlterTableCmd {
	p.advance() // consume ADD

	if p.collectMode() {
		p.cachedCollect("parseAlterTableAdd", func() {
			for _, t := range []int{
				COLUMN, IF_P, CONSTRAINT, CHECK,
				UNIQUE, PRIMARY, EXCLUDE, FOREIGN,
			} {
				p.addTokenCandidate(t)
			}
			// Also valid: column name (for ADD column_def without COLUMN keyword)
			p.addRuleCandidate("columnref")
		})
		return nil
	}

	switch p.cur.Type {
	case COLUMN:
		p.advance() // consume COLUMN
		// IF NOT EXISTS
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(NOT)
			p.expect(EXISTS)
			missingOk = true
		}
		coldef := p.parseColumnDef()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_AddColumn),
			Def:        coldef,
			Missing_ok: missingOk,
		}
	case IF_P:
		// ADD IF NOT EXISTS columnDef (without COLUMN keyword)
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		coldef := p.parseColumnDef()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_AddColumn),
			Def:        coldef,
			Missing_ok: true,
		}
	case CONSTRAINT, CHECK, UNIQUE, PRIMARY, EXCLUDE, FOREIGN:
		// ADD TableConstraint
		constr := p.parseTableConstraint()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AddConstraint),
			Def:     constr,
		}
	default:
		// ADD columnDef (without COLUMN keyword)
		// Try to distinguish between column def and constraint.
		// If current token looks like a column name followed by a type, it's a column def.
		coldef := p.parseColumnDef()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AddColumn),
			Def:     coldef,
		}
	}
}

// parseAlterTableDrop handles DROP ... subcommands.
func (p *Parser) parseAlterTableDrop() *nodes.AlterTableCmd {
	p.advance() // consume DROP

	if p.collectMode() {
		p.cachedCollect("parseAlterTableDrop", func() {
			for _, t := range []int{COLUMN, CONSTRAINT, IF_P} {
				p.addTokenCandidate(t)
			}
			// Also valid: column name directly (DROP col_name)
			p.addRuleCandidate("columnref")
		})
		return nil
	}

	switch p.cur.Type {
	case COLUMN:
		p.advance() // consume COLUMN
		// IF EXISTS
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(EXISTS)
			missingOk = true
		}
		name, _ := p.parseColId()
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropColumn),
			Name:       name,
			Behavior:   behavior,
			Missing_ok: missingOk,
		}
	case CONSTRAINT:
		p.advance() // consume CONSTRAINT
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(EXISTS)
			missingOk = true
		}
		name, _ := p.parseName()
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropConstraint),
			Name:       name,
			Behavior:   behavior,
			Missing_ok: missingOk,
		}
	case IF_P:
		// DROP IF EXISTS colid opt_drop_behavior (without COLUMN keyword)
		p.advance()
		p.expect(EXISTS)
		name, _ := p.parseColId()
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropColumn),
			Name:       name,
			Behavior:   behavior,
			Missing_ok: true,
		}
	default:
		// DROP colid opt_drop_behavior (without COLUMN keyword)
		name, _ := p.parseColId()
		behavior := p.parseOptDropBehavior()
		return &nodes.AlterTableCmd{
			Subtype:  int(nodes.AT_DropColumn),
			Name:     name,
			Behavior: behavior,
		}
	}
}

// parseAlterTableAlter handles ALTER [COLUMN] ... subcommands.
func (p *Parser) parseAlterTableAlter() *nodes.AlterTableCmd {
	p.advance() // consume ALTER

	if p.collectMode() {
		p.addTokenCandidate(COLUMN)
		p.addTokenCandidate(CONSTRAINT)
		p.addRuleCandidate("columnref")
		return nil
	}

	hasColumnKeyword := false
	if p.cur.Type == COLUMN {
		p.advance()
		hasColumnKeyword = true
	}

	// Check for ALTER CONSTRAINT (no column involved)
	if !hasColumnKeyword && p.cur.Type == CONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, _ := p.parseName()
		// ConstraintAttributeSpec (we consume but don't store — matches yacc behavior)
		p.parseConstraintAttributeSpec()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AlterConstraint),
			Name:    name,
		}
	}

	// Check for ALTER COLUMN Iconst (numeric column reference) SET STATISTICS
	if hasColumnKeyword && p.cur.Type == ICONST {
		num := int16(p.cur.Ival)
		p.advance()
		p.expect(SET)
		p.expect(STATISTICS)
		val := p.parseSignedIconst()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStatistics),
			Num:     num,
			Def:     makeIntConst(val),
		}
	}
	// ALTER Iconst (without COLUMN keyword) SET STATISTICS
	if !hasColumnKeyword && p.cur.Type == ICONST {
		num := int16(p.cur.Ival)
		p.advance()
		p.expect(SET)
		p.expect(STATISTICS)
		val := p.parseSignedIconst()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStatistics),
			Num:     num,
			Def:     makeIntConst(val),
		}
	}

	// Regular column name
	colname, _ := p.parseColId()

	return p.parseAlterColumnAction(colname)
}

// parseAlterColumnAction parses the action after ALTER [COLUMN] colname.
func (p *Parser) parseAlterColumnAction(colname string) *nodes.AlterTableCmd {
	switch p.cur.Type {
	case SET:
		p.advance() // consume SET
		return p.parseAlterColumnSet(colname)
	case DROP:
		p.advance() // consume DROP
		return p.parseAlterColumnDrop(colname)
	case TYPE_P:
		return p.parseAlterColumnType(colname, false)
	case ADD_P:
		p.advance() // consume ADD
		return p.parseAlterColumnAddGenerated(colname)
	case RESET:
		// ALTER COLUMN ColId RESET '(' def_list ')'
		p.advance() // consume RESET
		p.expect('(')
		defs, _ := p.parseDefList()
		p.expect(')')
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ResetOptions),
			Name:    colname,
			Def:     defs,
		}
	case OPTIONS:
		// alter_generic_options
		opts := p.parseAlterGenericOptions()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AlterColumnGenericOptions),
			Name:    colname,
			Def:     opts,
		}
	default:
		// ALTER COLUMN ColId alter_identity_column_option_list
		// The first token should be RESTART, SET, or a keyword that starts an identity option.
		if p.cur.Type == RESTART || (p.cur.Type == SET && p.isAlterIdentityOption()) {
			opts := p.parseAlterIdentityColumnOptionList()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetIdentity),
				Name:    colname,
				Def:     opts,
			}
		}
		return &nodes.AlterTableCmd{}
	}
}

// isAlterIdentityOption checks if the current SET is the start of an identity column option
// (SET GENERATED ...) rather than other SET subcommands. We peek ahead.
func (p *Parser) isAlterIdentityOption() bool {
	next := p.peekNext()
	return next.Type == GENERATED
}

// parseAlterColumnSet handles ALTER COLUMN colname SET ...
func (p *Parser) parseAlterColumnSet(colname string) *nodes.AlterTableCmd {
	switch p.cur.Type {
	case DEFAULT:
		// SET DEFAULT a_expr
		p.advance()
		expr, _ := p.parseAExpr(0)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ColumnDefault),
			Name:    colname,
			Def:     expr,
		}
	case NOT:
		// SET NOT NULL
		p.advance() // NOT
		p.expect(NULL_P)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetNotNull),
			Name:    colname,
		}
	case DATA_P:
		// SET DATA TYPE Typename ...
		p.advance() // DATA
		p.expect(TYPE_P)
		return p.parseAlterColumnTypeInner(colname)
	case STATISTICS:
		// SET STATISTICS SignedIconst
		p.advance()
		val := p.parseSignedIconst()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStatistics),
			Name:    colname,
			Def:     makeIntConst(val),
		}
	case STORAGE:
		// SET STORAGE ColId | SET STORAGE DEFAULT
		p.advance()
		var storageVal string
		if p.cur.Type == DEFAULT {
			p.advance()
			storageVal = "default"
		} else {
			storageVal, _ = p.parseColId()
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStorage),
			Name:    colname,
			Def:     &nodes.String{Str: storageVal},
		}
	case COMPRESSION:
		// SET COMPRESSION ColId
		p.advance()
		compVal, _ := p.parseColId()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetCompression),
			Name:    colname,
			Def:     &nodes.String{Str: compVal},
		}
	case EXPRESSION:
		// SET EXPRESSION AS '(' a_expr ')'
		p.advance()
		p.expect(AS)
		p.expect('(')
		expr, _ := p.parseAExpr(0)
		p.expect(')')
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetExpression),
			Name:    colname,
			Def:     expr,
		}
	case GENERATED:
		// SET GENERATED generated_when
		p.advance()
		gw := p.parseGeneratedWhen()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetIdentity),
			Name:    colname,
			Def: &nodes.List{Items: []nodes.Node{
				makeDefElem("generated", makeIntConst(int64(gw))),
			}},
		}
	case '(':
		// SET (def_list)
		p.advance()
		defs, _ := p.parseDefList()
		p.expect(')')
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetOptions),
			Name:    colname,
			Def:     defs,
		}
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterColumnDrop handles ALTER COLUMN colname DROP ...
func (p *Parser) parseAlterColumnDrop(colname string) *nodes.AlterTableCmd {
	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ColumnDefault),
			Name:    colname,
		}
	case NOT:
		p.advance() // NOT
		p.expect(NULL_P)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DropNotNull),
			Name:    colname,
		}
	case EXPRESSION:
		p.advance()
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(EXISTS)
			missingOk = true
		}
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropExpression),
			Name:       colname,
			Missing_ok: missingOk,
		}
	case IDENTITY_P:
		p.advance()
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			p.expect(EXISTS)
			missingOk = true
		}
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_DropIdentity),
			Name:       colname,
			Missing_ok: missingOk,
		}
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterColumnType handles ALTER COLUMN colname TYPE Typename ...
func (p *Parser) parseAlterColumnType(colname string, hasSetData bool) *nodes.AlterTableCmd {
	if !hasSetData {
		p.advance() // consume TYPE
	}
	return p.parseAlterColumnTypeInner(colname)
}

// parseAlterColumnTypeInner is the common code for TYPE typename [COLLATE] [USING]
func (p *Parser) parseAlterColumnTypeInner(colname string) *nodes.AlterTableCmd {
	tn, _ := p.parseTypename()
	coldef := &nodes.ColumnDef{TypeName: tn}

	// opt_collate_clause
	if p.cur.Type == COLLATE {
		p.advance()
		collname, _ := p.parseAnyName()
		coldef.CollClause = &nodes.CollateClause{
			Collname: collname,
			Loc: nodes.NoLoc(),
		}
	}

	// USING a_expr (consumed but not stored — matches yacc behavior where it goes to
	// transform analysis, not the raw parse tree)
	if p.cur.Type == USING {
		p.advance()
		_, _ = p.parseAExpr(0)
	}

	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AlterColumnType),
		Name:    colname,
		Def:     coldef,
	}
}

// parseAlterColumnAddGenerated handles ALTER COLUMN colname ADD GENERATED ...
func (p *Parser) parseAlterColumnAddGenerated(colname string) *nodes.AlterTableCmd {
	p.expect(GENERATED)
	gw := p.parseGeneratedWhen()
	p.expect(AS)
	p.expect(IDENTITY_P)
	opts := p.parseOptParenthesizedSeqOptList()
	c := &nodes.Constraint{
		Contype:       nodes.CONSTR_IDENTITY,
		GeneratedWhen: byte(gw),
		Options:       opts,
		Loc: nodes.NoLoc(),
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddIdentity),
		Name:    colname,
		Def:     c,
	}
}

// parseAlterTableOwner handles OWNER TO RoleSpec.
func (p *Parser) parseAlterTableOwner() *nodes.AlterTableCmd {
	p.advance() // consume OWNER
	p.expect(TO)
	roleSpec := p.parseRoleSpec()
	return &nodes.AlterTableCmd{
		Subtype:  int(nodes.AT_ChangeOwner),
		Newowner: roleSpec,
	}
}

// parseAlterTableValidate handles VALIDATE CONSTRAINT name.
func (p *Parser) parseAlterTableValidate() *nodes.AlterTableCmd {
	p.advance() // consume VALIDATE
	p.expect(CONSTRAINT)
	name, _ := p.parseName()
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ValidateConstraint),
		Name:    name,
	}
}

// parseAlterTableInherit handles INHERIT qualified_name.
func (p *Parser) parseAlterTableInherit() *nodes.AlterTableCmd {
	p.advance() // consume INHERIT
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddInherit),
		Def:     rv,
	}
}

// parseAlterTableNo handles NO INHERIT and NO FORCE ROW LEVEL SECURITY.
func (p *Parser) parseAlterTableNo() *nodes.AlterTableCmd {
	p.advance() // consume NO

	switch p.cur.Type {
	case INHERIT:
		p.advance()
		names, _ := p.parseQualifiedName()
		rv := makeRangeVarFromAnyName(names)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DropInherit),
			Def:     rv,
		}
	case FORCE:
		p.advance() // FORCE
		p.expect(ROW)
		p.expect(LEVEL)
		p.expect(SECURITY)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_NoForceRowSecurity),
		}
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterTableAttach handles ATTACH PARTITION qualified_name ForValues.
func (p *Parser) parseAlterTableAttach() *nodes.AlterTableCmd {
	p.advance() // consume ATTACH
	p.expect(PARTITION)
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)
	bound := p.parseForValues()
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AttachPartition),
		Def: &nodes.PartitionCmd{
			Name:  rv,
			Bound: bound.(*nodes.PartitionBoundSpec),
		},
	}
}

// parseAlterTableDetach handles DETACH PARTITION ...
func (p *Parser) parseAlterTableDetach() *nodes.AlterTableCmd {
	p.advance() // consume DETACH
	p.expect(PARTITION)
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	if p.cur.Type == CONCURRENTLY {
		p.advance()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DetachPartition),
			Def: &nodes.PartitionCmd{
				Name:       rv,
				Concurrent: true,
			},
		}
	}
	if p.cur.Type == FINALIZE {
		p.advance()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DetachPartitionFinalize),
			Def: &nodes.PartitionCmd{
				Name: rv,
			},
		}
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_DetachPartition),
		Def: &nodes.PartitionCmd{
			Name: rv,
		},
	}
}

// parseAlterTableEnable handles ENABLE TRIGGER/RULE/ROW LEVEL SECURITY.
func (p *Parser) parseAlterTableEnable() *nodes.AlterTableCmd {
	p.advance() // consume ENABLE

	switch p.cur.Type {
	case TRIGGER:
		p.advance()
		switch p.cur.Type {
		case ALL:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableTrigAll)}
		case USER:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableTrigUser)}
		default:
			name, _ := p.parseName()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableTrig),
				Name:    name,
			}
		}
	case ALWAYS:
		p.advance() // consume ALWAYS
		if p.cur.Type == TRIGGER {
			p.advance()
			name, _ := p.parseName()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableAlwaysTrig),
				Name:    name,
			}
		}
		// ALWAYS RULE
		p.expect(RULE)
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_EnableAlwaysRule),
			Name:    name,
		}
	case REPLICA:
		p.advance() // consume REPLICA
		if p.cur.Type == TRIGGER {
			p.advance()
			name, _ := p.parseName()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableReplicaTrig),
				Name:    name,
			}
		}
		// REPLICA RULE
		p.expect(RULE)
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_EnableReplicaRule),
			Name:    name,
		}
	case RULE:
		p.advance()
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_EnableRule),
			Name:    name,
		}
	case ROW:
		p.advance() // ROW
		p.expect(LEVEL)
		p.expect(SECURITY)
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableRowSecurity)}
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterTableDisable handles DISABLE TRIGGER/RULE/ROW LEVEL SECURITY.
func (p *Parser) parseAlterTableDisable() *nodes.AlterTableCmd {
	p.advance() // consume DISABLE

	switch p.cur.Type {
	case TRIGGER:
		p.advance()
		switch p.cur.Type {
		case ALL:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableTrigAll)}
		case USER:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableTrigUser)}
		default:
			name, _ := p.parseName()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DisableTrig),
				Name:    name,
			}
		}
	case RULE:
		p.advance()
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DisableRule),
			Name:    name,
		}
	case ROW:
		p.advance() // ROW
		p.expect(LEVEL)
		p.expect(SECURITY)
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableRowSecurity)}
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterTableForce handles FORCE ROW LEVEL SECURITY.
func (p *Parser) parseAlterTableForce() *nodes.AlterTableCmd {
	p.advance() // consume FORCE
	p.expect(ROW)
	p.expect(LEVEL)
	p.expect(SECURITY)
	return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ForceRowSecurity)}
}

// parseAlterTableCluster handles CLUSTER ON name.
func (p *Parser) parseAlterTableCluster() *nodes.AlterTableCmd {
	p.advance() // consume CLUSTER
	p.expect(ON)
	name, _ := p.parseName()
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ClusterOn),
		Name:    name,
	}
}

// parseAlterTableSet handles SET ... subcommands at the table level.
func (p *Parser) parseAlterTableSet() *nodes.AlterTableCmd {
	p.advance() // consume SET

	switch p.cur.Type {
	case WITHOUT:
		p.advance() // consume WITHOUT
		switch p.cur.Type {
		case CLUSTER:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropCluster)}
		case OIDS:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOids)}
		default:
			return &nodes.AlterTableCmd{}
		}
	case LOGGED:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetLogged)}
	case UNLOGGED:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetUnLogged)}
	case ACCESS:
		p.advance() // consume ACCESS
		p.expect(METHOD)
		if p.cur.Type == DEFAULT {
			p.advance()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetAccessMethod),
				Name:    "",
			}
		}
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetAccessMethod),
			Name:    name,
		}
	case '(':
		// SET reloptions (starts with '(' because reloptions = '(' reloption_list ')')
		opts := p.parseReloptions()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetRelOptions),
			Def:     opts,
		}
	case WITH:
		p.advance() // consume WITH
		if p.cur.Type == OIDS {
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOids)}
		}
		return &nodes.AlterTableCmd{}
	case TABLESPACE:
		p.advance()
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetTableSpace),
			Name:    name,
		}
	default:
		return &nodes.AlterTableCmd{}
	}
}

// parseAlterTableReset handles RESET reloptions.
func (p *Parser) parseAlterTableReset() *nodes.AlterTableCmd {
	p.advance() // consume RESET
	opts := p.parseReloptions()
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ResetRelOptions),
		Def:     opts,
	}
}

// parseAlterTableReplica handles REPLICA IDENTITY ...
func (p *Parser) parseAlterTableReplica() *nodes.AlterTableCmd {
	p.advance() // consume REPLICA
	p.expect(IDENTITY_P)

	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ReplicaIdentity)}
	case FULL:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ReplicaIdentity)}
	case NOTHING:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ReplicaIdentity)}
	case USING:
		p.advance() // USING
		p.expect(INDEX)
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ReplicaIdentity),
			Name:    name,
		}
	default:
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ReplicaIdentity)}
	}
}

// parseAlterTableOf handles OF any_name.
func (p *Parser) parseAlterTableOf() *nodes.AlterTableCmd {
	p.advance() // consume OF
	names, _ := p.parseAnyName()
	tn := makeTypeNameFromNameList(names)
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddOf),
		Def:     tn,
	}
}

// parseAlterTableNot handles NOT OF.
func (p *Parser) parseAlterTableNot() *nodes.AlterTableCmd {
	p.advance() // consume NOT
	p.expect(OF)
	return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOf)}
}

// parseAlterTableOptions handles bare alter_generic_options (for foreign tables).
func (p *Parser) parseAlterTableOptions() *nodes.AlterTableCmd {
	opts := p.parseAlterGenericOptions()
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_GenericOptions),
		Def:     opts,
	}
}

// ---------------------------------------------------------------------------
// Helper parsers
// ---------------------------------------------------------------------------

// parseOptDropBehavior parses opt_drop_behavior.
//
//	opt_drop_behavior: CASCADE | RESTRICT | /* EMPTY */
func (p *Parser) parseOptDropBehavior() int {
	switch p.cur.Type {
	case CASCADE:
		p.advance()
		return int(nodes.DROP_CASCADE)
	case RESTRICT:
		p.advance()
		return int(nodes.DROP_RESTRICT)
	default:
		return int(nodes.DROP_RESTRICT)
	}
}

// parseRoleSpec parses a RoleSpec.
//
//	RoleSpec: ColId | CURRENT_ROLE | CURRENT_USER | SESSION_USER
func (p *Parser) parseRoleSpec() *nodes.RoleSpec {
	loc := p.pos()
	var rs *nodes.RoleSpec
	switch p.cur.Type {
	case CURRENT_ROLE:
		p.advance()
		rs = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CURRENT_ROLE)}
	case CURRENT_USER:
		p.advance()
		rs = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_CURRENT_USER)}
	case SESSION_USER:
		p.advance()
		rs = &nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_SESSION_USER)}
	default:
		name, _ := p.parseColId()
		rs = &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: name,
		}
	}
	rs.Loc = nodes.Loc{Start: loc, End: p.pos()}
	return rs
}

// parseRoleList parses role_list: RoleSpec [, RoleSpec ...]
func (p *Parser) parseRoleList() *nodes.List {
	rs := p.parseRoleSpec()
	items := []nodes.Node{rs}
	for p.cur.Type == ',' {
		p.advance()
		rs = p.parseRoleSpec()
		items = append(items, rs)
	}
	return &nodes.List{Items: items}
}

// parseSignedIconst parses SignedIconst: Iconst | '+' Iconst | '-' Iconst
func (p *Parser) parseSignedIconst() int64 {
	neg := false
	if p.cur.Type == '+' {
		p.advance()
	} else if p.cur.Type == '-' {
		p.advance()
		neg = true
	}
	val := p.cur.Ival
	p.advance()
	if neg {
		return -val
	}
	return val
}

// parseGeneratedWhen parses generated_when: ALWAYS | BY DEFAULT
func (p *Parser) parseGeneratedWhen() byte {
	if p.cur.Type == ALWAYS {
		p.advance()
		return 'a'
	}
	// BY DEFAULT
	p.expect(BY)
	p.expect(DEFAULT)
	return 'd'
}

// parseAlterIdentityColumnOptionList parses alter_identity_column_option_list.
func (p *Parser) parseAlterIdentityColumnOptionList() *nodes.List {
	var items []nodes.Node
	for {
		opt := p.parseAlterIdentityColumnOption()
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseAlterIdentityColumnOption parses alter_identity_column_option.
//
//	RESTART | RESTART opt_with NumericOnly | SET SeqOptElem | SET GENERATED generated_when
func (p *Parser) parseAlterIdentityColumnOption() nodes.Node {
	switch p.cur.Type {
	case RESTART:
		p.advance()
		// optional: WITH NumericOnly or just NumericOnly
		if p.cur.Type == WITH_LA || p.cur.Type == WITH {
			p.advance()
		}
		if p.cur.Type == ICONST || p.cur.Type == FCONST || p.cur.Type == '+' || p.cur.Type == '-' {
			val := p.parseNumericOnly()
			return makeDefElem("restart", val)
		}
		return makeDefElem("restart", nil)
	case SET:
		p.advance()
		if p.cur.Type == GENERATED {
			p.advance()
			gw := p.parseGeneratedWhen()
			return makeDefElem("generated", makeIntConst(int64(gw)))
		}
		// SET SeqOptElem
		return p.parseOneSeqOptElem()
	default:
		return nil
	}
}

// parseAlterGenericOptions parses alter_generic_options:
//
//	OPTIONS '(' alter_generic_option_list ')'
func (p *Parser) parseAlterGenericOptions() *nodes.List {
	p.expect(OPTIONS)
	p.expect('(')
	list := p.parseAlterGenericOptionList()
	p.expect(')')
	return list
}

// parseAlterGenericOptionList parses alter_generic_option_list.
func (p *Parser) parseAlterGenericOptionList() *nodes.List {
	elem := p.parseAlterGenericOptionElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseAlterGenericOptionElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseAlterGenericOptionElem parses alter_generic_option_elem.
//
//	generic_option_elem
//	| SET generic_option_elem
//	| ADD_P generic_option_elem
//	| DROP generic_option_name
func (p *Parser) parseAlterGenericOptionElem() *nodes.DefElem {
	switch p.cur.Type {
	case SET:
		p.advance()
		n := p.parseGenericOptionElem()
		n.Defaction = int(nodes.DEFELEM_SET)
		return n
	case ADD_P:
		p.advance()
		n := p.parseGenericOptionElem()
		n.Defaction = int(nodes.DEFELEM_ADD)
		return n
	case DROP:
		p.advance()
		name, _ := p.parseColLabel()
		return &nodes.DefElem{
			Defname:   name,
			Defaction: int(nodes.DEFELEM_DROP),
			Loc: nodes.NoLoc(),
		}
	default:
		return p.parseGenericOptionElem()
	}
}

// parseGenericOptionElem parses generic_option_elem: generic_option_name generic_option_arg
func (p *Parser) parseGenericOptionElem() *nodes.DefElem {
	name, _ := p.parseColLabel()
	// generic_option_arg is Sconst
	arg := p.cur.Str
	p.advance()
	return &nodes.DefElem{
		Defname:  name,
		Arg:      &nodes.String{Str: arg},
		Loc: nodes.NoLoc(),
	}
}

// parseOneSeqOptElem parses a single SeqOptElem using the existing parseSeqOptList.
func (p *Parser) parseOneSeqOptElem() *nodes.DefElem {
	list := p.parseSeqOptList()
	if list != nil && len(list.Items) > 0 {
		if de, ok := list.Items[0].(*nodes.DefElem); ok {
			return de
		}
	}
	return &nodes.DefElem{Loc: nodes.NoLoc()}
}

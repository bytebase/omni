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
func (p *Parser) parseAlterTableStmt() (nodes.Node, error) {
	// Already consumed ALTER; current token is TABLE, INDEX, SEQUENCE, VIEW, MATERIALIZED, or FOREIGN.
	alterLoc := p.prev.Loc // position of ALTER keyword
	switch p.cur.Type {
	case TABLE:
		return p.parseAlterTable(alterLoc)
	case INDEX:
		return p.parseAlterIndex(alterLoc)
	case SEQUENCE:
		return p.parseAlterSequence(alterLoc)
	case VIEW:
		return p.parseAlterView(alterLoc)
	case MATERIALIZED:
		return p.parseAlterMaterializedView(alterLoc)
	case FOREIGN:
		return p.parseAlterForeignTable(alterLoc)
	case EVENT:
		p.advance() // consume EVENT
		return p.parseAlterEventTrigStmt()
	case EXTENSION:
		return p.parseAlterExtensionStmt()
	// exhaustive: gram.y:2081 AlterTableStmt — caller handles nil via outer error
	default:
		return nil, nil
	}
}

// parseAlterTable handles ALTER TABLE ...
func (p *Parser) parseAlterTable(alterLoc int) (nodes.Node, error) {
	p.advance() // consume TABLE

	// ALTER TABLE ALL IN TABLESPACE ...
	if p.cur.Type == ALL {
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_TABLE), alterLoc), nil
	}

	// IF EXISTS
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	// relation_expr
	rel, _ := p.parseRelationExpr()

	// Check for RENAME (produces RenameStmt, not AlterTableStmt)
	if p.cur.Type == RENAME {
		return p.parseAlterTableRename(rel, missingOk, alterLoc)
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
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}

	// alter_table_cmds
	cmds, err := p.parseAlterTableCmds()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableStmt{
		Relation:   rel,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_TABLE),
		Missing_ok: missingOk,
		Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseAlterIndex handles ALTER INDEX ...
func (p *Parser) parseAlterIndex(alterLoc int) (nodes.Node, error) {
	p.advance() // consume INDEX

	// ALTER INDEX ALL IN TABLESPACE ...
	if p.cur.Type == ALL {
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_INDEX), alterLoc), nil
	}

	// IF EXISTS
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	// qualified_name
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	// Check for RENAME
	if p.cur.Type == RENAME {
		p.advance() // consume RENAME
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_INDEX,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}

	// ALTER INDEX name ATTACH PARTITION name
	if p.cur.Type == ATTACH {
		attachLoc := p.pos()
		p.advance() // ATTACH
		if _, err := p.expect(PARTITION); err != nil {
			return nil, err
		}
		partNames, _ := p.parseQualifiedName()
		partRv := makeRangeVarFromAnyName(partNames)
		cmd := &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AttachPartition),
			Def: &nodes.PartitionCmd{
				Name: partRv,
				Loc:  nodes.Loc{Start: attachLoc, End: p.prev.End},
			},
			Loc: nodes.Loc{Start: attachLoc, End: p.prev.End},
		}
		return &nodes.AlterTableStmt{
			Relation:   rv,
			Cmds:       &nodes.List{Items: []nodes.Node{cmd}},
			ObjType:    int(nodes.OBJECT_INDEX),
			Missing_ok: missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}

	// ALTER INDEX name [NO] DEPENDS ON EXTENSION ext_name
	if p.cur.Type == DEPENDS || (p.cur.Type == NO && p.peekNext().Type == DEPENDS) {
		remove := p.parseOptNo()
		p.advance() // consume DEPENDS
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXTENSION); err != nil {
			return nil, err
		}
		extname, _ := p.parseName()
		return &nodes.AlterObjectDependsStmt{
			ObjectType: nodes.OBJECT_INDEX,
			Relation:   rv,
			Extname:    &nodes.String{Str: extname},
			Remove:     remove,
		}, nil
	}

	cmds, err := p.parseAlterTableCmds()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableStmt{
		Relation:   rv,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_INDEX),
		Missing_ok: missingOk,
		Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseAlterSequence handles ALTER SEQUENCE ...
func (p *Parser) parseAlterSequence(alterLoc int) (nodes.Node, error) {
	p.advance() // consume SEQUENCE

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	names, _ := p.parseQualifiedName()

	// Check for RENAME (produces RenameStmt, uses makeRangeVarFromAnyName)
	if p.cur.Type == RENAME {
		rv := makeRangeVarFromAnyName(names)
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_SEQUENCE,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
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
		}, nil
	}

	// PG 15+: ALTER SEQUENCE name SET LOGGED | SET UNLOGGED
	// These appear in PG's SeqOptElem grammar but require the SET keyword (unlike
	// CREATE SEQUENCE where bare LOGGED/UNLOGGED is accepted). Produce the same
	// DefElem shape as the bare LOGGED/UNLOGGED SeqOptElem (Defname="logged",
	// Arg=Boolean{true/false}) so downstream consumers stay uniform.
	if p.cur.Type == SET && (p.peekNext().Type == LOGGED || p.peekNext().Type == UNLOGGED) {
		rv := makeRangeVarFromNames(names)
		p.advance() // consume SET
		logged := p.cur.Type == LOGGED
		p.advance() // consume LOGGED or UNLOGGED
		opt := makeDefElem("logged", &nodes.Boolean{Boolval: logged})
		return &nodes.AlterSeqStmt{
			Sequence:  rv,
			Options:   &nodes.List{Items: []nodes.Node{opt}},
			MissingOk: missingOk,
			Loc:       nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}

	// OWNER TO goes through alter_table_cmds in the yacc grammar
	if p.cur.Type == OWNER {
		rv := makeRangeVarFromAnyName(names)
		cmds, err := p.parseAlterTableCmds()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableStmt{
			Relation:   rv,
			Cmds:       cmds,
			ObjType:    int(nodes.OBJECT_SEQUENCE),
			Missing_ok: missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}

	// AlterSeqStmt: ALTER SEQUENCE name SeqOptList
	rv := makeRangeVarFromNames(names)
	options, err := p.parseSeqOptList()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterSeqStmt{
		Sequence:  rv,
		Options:   options,
		MissingOk: missingOk,
		Loc:       nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseAlterView handles ALTER VIEW ...
func (p *Parser) parseAlterView(alterLoc int) (nodes.Node, error) {
	p.advance() // consume VIEW

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	// Check for RENAME
	if p.cur.Type == RENAME {
		p.advance()
		// RENAME COLUMN oldname TO newname
		if p.cur.Type == COLUMN {
			p.advance()
			oldname, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			newname, _ := p.parseName()
			return &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_VIEW,
				Relation:     rv,
				Subname:      oldname,
				Newname:      newname,
				MissingOk:    missingOk,
				Loc:          nodes.Loc{Start: alterLoc, End: p.prev.End},
			}, nil
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_VIEW,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
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
		}, nil
	}

	cmds, err := p.parseAlterTableCmds()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableStmt{
		Relation:   rv,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_VIEW),
		Missing_ok: missingOk,
		Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseAlterMaterializedView handles ALTER MATERIALIZED VIEW ...
func (p *Parser) parseAlterMaterializedView(alterLoc int) (nodes.Node, error) {
	p.advance() // consume MATERIALIZED
	if _, err := p.expect(VIEW); err != nil {
		return nil, err
	}

	// ALTER MATERIALIZED VIEW ALL IN TABLESPACE ...
	if p.cur.Type == ALL {
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_MATVIEW), alterLoc), nil
	}

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names)

	// Check for RENAME
	if p.cur.Type == RENAME {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_MATVIEW,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
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
		}, nil
	}

	// ALTER MATERIALIZED VIEW name [NO] DEPENDS ON EXTENSION ext_name
	if p.cur.Type == DEPENDS || (p.cur.Type == NO && p.peekNext().Type == DEPENDS) {
		remove := p.parseOptNo()
		p.advance() // consume DEPENDS
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXTENSION); err != nil {
			return nil, err
		}
		extname, _ := p.parseName()
		return &nodes.AlterObjectDependsStmt{
			ObjectType: nodes.OBJECT_MATVIEW,
			Relation:   rv,
			Extname:    &nodes.String{Str: extname},
			Remove:     remove,
		}, nil
	}

	cmds, err := p.parseAlterTableCmds()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableStmt{
		Relation:   rv,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_MATVIEW),
		Missing_ok: missingOk,
		Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseAlterForeignTable handles ALTER FOREIGN TABLE ... and ALTER FOREIGN DATA WRAPPER ...
func (p *Parser) parseAlterForeignTable(alterLoc int) (nodes.Node, error) {
	p.advance() // consume FOREIGN

	// ALTER FOREIGN DATA WRAPPER ...
	if p.cur.Type == DATA_P {
		return p.parseAlterFdwStmt()
	}

	if _, err := p.expect(TABLE); err != nil {
		return nil, err
	}

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	rel, _ := p.parseRelationExpr()

	// Check for RENAME
	if p.cur.Type == RENAME {
		return p.parseAlterTableRename(rel, missingOk, alterLoc)
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
		}, nil
	}

	cmds, err := p.parseAlterTableCmds()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableStmt{
		Relation:   rel,
		Cmds:       cmds,
		ObjType:    int(nodes.OBJECT_FOREIGN_TABLE),
		Missing_ok: missingOk,
		Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
	}, nil
}

// parseAlterTableMoveAll parses ALTER TABLE/INDEX/MATERIALIZED VIEW ALL IN TABLESPACE ...
//
//	ALTER TABLE ALL IN TABLESPACE name [OWNED BY role_list] SET TABLESPACE name [NOWAIT]
func (p *Parser) parseAlterTableMoveAll(objType int, alterLoc int) *nodes.AlterTableMoveAllStmt {
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
		Loc:                nodes.Loc{Start: alterLoc, End: p.prev.End},
	}
}

// parseAlterTableRename parses ALTER TABLE ... RENAME ...
func (p *Parser) parseAlterTableRename(rel *nodes.RangeVar, missingOk bool, alterLoc int) (*nodes.RenameStmt, error) {
	p.advance() // consume RENAME

	if p.collectMode() {
		p.addTokenCandidate(TO)
		p.addTokenCandidate(COLUMN)
		p.addTokenCandidate(CONSTRAINT)
		// Also valid: column name directly (RENAME col_name TO ...)
		p.addRuleCandidate("columnref")
		return nil, nil
	}

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
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	case COLUMN:
		// RENAME COLUMN colname TO newname
		p.advance() // consume COLUMN
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, nil
		}
		oldname, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		if oldname == "" {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType:   nodes.OBJECT_COLUMN,
			RelationType: nodes.OBJECT_TABLE,
			Relation:     rel,
			Subname:      oldname,
			Newname:      newname,
			MissingOk:    missingOk,
			Loc:          nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	case CONSTRAINT:
		// RENAME CONSTRAINT oldname TO newname
		p.advance() // consume CONSTRAINT
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			return nil, nil
		}
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
			Loc:          nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
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
			Loc:          nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}
}

// parseAlterTableCmds parses a comma-separated list of alter_table_cmd.
//
//	alter_table_cmds: alter_table_cmd | alter_table_cmds ',' alter_table_cmd
func (p *Parser) parseAlterTableCmds() (*nodes.List, error) {
	cmd, err := p.parseAlterTableCmd()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{cmd}
	for p.cur.Type == ',' {
		p.advance()
		cmd, err = p.parseAlterTableCmd()
		if err != nil {
			return nil, err
		}
		items = append(items, cmd)
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterTableCmd parses a single alter_table_cmd.
func (p *Parser) parseAlterTableCmd() (*nodes.AlterTableCmd, error) {
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
		return nil, nil
	}
	cmdLoc := p.pos()
	var cmd *nodes.AlterTableCmd
	var err error
	switch p.cur.Type {
	case ADD_P:
		cmd, err = p.parseAlterTableAdd(), nil
	case DROP:
		cmd, err = p.parseAlterTableDrop(), nil
	case ALTER:
		cmd, err = p.parseAlterTableAlter()
	case OWNER:
		cmd, err = p.parseAlterTableOwner(), nil
	case VALIDATE:
		cmd, err = p.parseAlterTableValidate(), nil
	case INHERIT:
		cmd, err = p.parseAlterTableInherit()
	case NO:
		cmd, err = p.parseAlterTableNo(), nil
	case ATTACH:
		cmd, err = p.parseAlterTableAttach(), nil
	case DETACH:
		cmd, err = p.parseAlterTableDetach(), nil
	case ENABLE_P:
		cmd, err = p.parseAlterTableEnable(), nil
	case DISABLE_P:
		cmd, err = p.parseAlterTableDisable(), nil
	case FORCE:
		cmd, err = p.parseAlterTableForce(), nil
	case CLUSTER:
		cmd, err = p.parseAlterTableCluster(), nil
	case SET:
		cmd, err = p.parseAlterTableSet()
	case RESET:
		cmd, err = p.parseAlterTableReset(), nil
	case REPLICA:
		cmd, err = p.parseAlterTableReplica(), nil
	case OF:
		cmd, err = p.parseAlterTableOf()
	case NOT:
		cmd, err = p.parseAlterTableNot(), nil
	case OPTIONS:
		cmd, err = p.parseAlterTableOptions(), nil
	default:
		// exhaustive-gap fix: unknown alter_table_cmd lead keyword.
		// Previously returned a zero-value AlterTableCmd (Subtype=0 == AT_AddColumn),
		// injecting a phantom AddColumn command into the Cmds list.
		return nil, p.syntaxErrorAtCur()
	}
	if cmd != nil {
		cmd.Loc = nodes.Loc{Start: cmdLoc, End: p.prev.End}
	}
	return cmd, err
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

	if p.isTableConstraintStart() {
		// ADD TableConstraint
		constr := p.parseTableConstraint()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AddConstraint),
			Def:     constr,
		}
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
		coldef, _ := p.parseColumnDef()
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
		coldef, _ := p.parseColumnDef()
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_AddColumn),
			Def:        coldef,
			Missing_ok: true,
		}
	default:
		// ADD columnDef (without COLUMN keyword)
		// Try to distinguish between column def and constraint.
		// If current token looks like a column name followed by a type, it's a column def.
		coldef, _ := p.parseColumnDef()
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
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addTokenCandidate(IF_P)
			return nil
		}
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
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			p.addTokenCandidate(IF_P)
			return nil
		}
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
func (p *Parser) parseAlterTableAlter() (*nodes.AlterTableCmd, error) {
	p.advance() // consume ALTER

	if p.collectMode() {
		p.addTokenCandidate(COLUMN)
		p.addTokenCandidate(CONSTRAINT)
		p.addRuleCandidate("columnref")
		return nil, nil
	}

	hasColumnKeyword := false
	if p.cur.Type == COLUMN {
		p.advance()
		hasColumnKeyword = true
	}

	if hasColumnKeyword && p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, nil
	}

	// Check for ALTER CONSTRAINT (no column involved)
	if !hasColumnKeyword && p.cur.Type == CONSTRAINT {
		p.advance() // consume CONSTRAINT
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			return nil, nil
		}
		name, _ := p.parseName()
		// ConstraintAttributeSpec (we consume but don't store — matches yacc behavior)
		p.parseConstraintAttributeSpec()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AlterConstraint),
			Name:    name,
		}, nil
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
		}, nil
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
		}, nil
	}

	// Regular column name
	colname, _ := p.parseColId()

	return p.parseAlterColumnAction(colname)
}

// parseAlterColumnAction parses the action after ALTER [COLUMN] colname.
func (p *Parser) parseAlterColumnAction(colname string) (*nodes.AlterTableCmd, error) {
	switch p.cur.Type {
	case SET:
		if p.isAlterIdentityOption() {
			opts, err := p.parseAlterIdentityColumnOptionList()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetIdentity),
				Name:    colname,
				Def:     opts,
			}, nil
		}
		p.advance() // consume SET
		return p.parseAlterColumnSet(colname)
	case DROP:
		p.advance() // consume DROP
		return p.parseAlterColumnDrop(colname), nil
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
		}, nil
	case OPTIONS:
		// alter_generic_options
		opts := p.parseAlterGenericOptions()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AlterColumnGenericOptions),
			Name:    colname,
			Def:     opts,
		}, nil
	default:
		// ALTER COLUMN ColId alter_identity_column_option_list
		// The first token should be RESTART, SET, or a keyword that starts an identity option.
		if p.cur.Type == RESTART || (p.cur.Type == SET && p.isAlterIdentityOption()) {
			opts, err := p.parseAlterIdentityColumnOptionList()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetIdentity),
				Name:    colname,
				Def:     opts,
			}, nil
		}
		// exhaustive-gap fix: unknown sub-action after ALTER [COLUMN] name.
		// Previously returned a zero-value AlterTableCmd (AT_AddColumn ghost).
		return nil, p.syntaxErrorAtCur()
	}
}

// isAlterIdentityOption checks if the current SET is the start of an identity column option
// (SET GENERATED ...) rather than other SET subcommands. We peek ahead.
func (p *Parser) isAlterIdentityOption() bool {
	next := p.peekNext()
	switch next.Type {
	case GENERATED, CACHE, CYCLE, NO, INCREMENT, MAXVALUE, MINVALUE, START:
		return true
	}
	return false
}

// parseAlterColumnSet handles ALTER COLUMN colname SET ...
func (p *Parser) parseAlterColumnSet(colname string) (*nodes.AlterTableCmd, error) {
	switch p.cur.Type {
	case DEFAULT:
		// SET DEFAULT a_expr
		p.advance()
		expr, _ := p.parseAExpr(0)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ColumnDefault),
			Name:    colname,
			Def:     expr,
		}, nil
	case NOT:
		// SET NOT NULL
		p.advance() // NOT
		p.expect(NULL_P)
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetNotNull),
			Name:    colname,
		}, nil
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
		}, nil
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
		}, nil
	case COMPRESSION:
		// SET COMPRESSION ColId | SET COMPRESSION DEFAULT
		p.advance()
		var compVal string
		if p.cur.Type == DEFAULT {
			p.advance()
			compVal = "default"
		} else {
			compVal, _ = p.parseColId()
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetCompression),
			Name:    colname,
			Def:     &nodes.String{Str: compVal},
		}, nil
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
		}, nil
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
		}, nil
	case '(':
		// SET (def_list)
		p.advance()
		defs, _ := p.parseDefList()
		p.expect(')')
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetOptions),
			Name:    colname,
			Def:     defs,
		}, nil
	default:
		// exhaustive-gap fix: unknown keyword after ALTER [COLUMN] name SET.
		// Previously returned a zero-value AlterTableCmd (AT_AddColumn ghost).
		return nil, p.syntaxErrorAtCur()
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
func (p *Parser) parseAlterColumnType(colname string, hasSetData bool) (*nodes.AlterTableCmd, error) {
	if !hasSetData {
		p.advance() // consume TYPE
	}
	return p.parseAlterColumnTypeInner(colname)
}

// parseAlterColumnTypeInner is the common code for TYPE typename [COLLATE] [USING]
func (p *Parser) parseAlterColumnTypeInner(colname string) (*nodes.AlterTableCmd, error) {
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if tn == nil {
		return nil, p.syntaxErrorAtCur()
	}
	coldef := &nodes.ColumnDef{TypeName: tn}

	// opt_collate_clause
	if p.cur.Type == COLLATE {
		p.advance()
		collname, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		coldef.CollClause = &nodes.CollateClause{
			Collname: collname,
			Loc:      nodes.NoLoc(),
		}
	}

	// USING a_expr (consumed but not stored — matches yacc behavior where it goes to
	// transform analysis, not the raw parse tree)
	if p.cur.Type == USING {
		p.advance()
		_, err = p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
	}

	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AlterColumnType),
		Name:    colname,
		Def:     coldef,
	}, nil
}

// parseAlterColumnAddGenerated handles ALTER COLUMN colname ADD GENERATED ...
func (p *Parser) parseAlterColumnAddGenerated(colname string) (*nodes.AlterTableCmd, error) {
	p.expect(GENERATED)
	gw := p.parseGeneratedWhen()
	p.expect(AS)
	p.expect(IDENTITY_P)
	opts, err := p.parseOptParenthesizedSeqOptList()
	if err != nil {
		return nil, err
	}
	c := &nodes.Constraint{
		Contype:       nodes.CONSTR_IDENTITY,
		GeneratedWhen: byte(gw),
		Options:       opts,
		Loc:           nodes.NoLoc(),
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddIdentity),
		Name:    colname,
		Def:     c,
	}, nil
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
func (p *Parser) parseAlterTableInherit() (*nodes.AlterTableCmd, error) {
	p.advance() // consume INHERIT
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if names == nil {
		return nil, p.syntaxErrorAtCur()
	}
	rv := makeRangeVarFromAnyName(names)
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddInherit),
		Def:     rv,
	}, nil
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
	partLoc := p.pos()
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
			Loc:   nodes.Loc{Start: partLoc, End: p.prev.End},
		},
	}
}

// parseAlterTableDetach handles DETACH PARTITION ...
func (p *Parser) parseAlterTableDetach() *nodes.AlterTableCmd {
	partLoc := p.pos()
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
				Loc:        nodes.Loc{Start: partLoc, End: p.prev.End},
			},
		}
	}
	if p.cur.Type == FINALIZE {
		p.advance()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DetachPartitionFinalize),
			Def: &nodes.PartitionCmd{
				Name: rv,
				Loc:  nodes.Loc{Start: partLoc, End: p.prev.End},
			},
		}
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_DetachPartition),
		Def: &nodes.PartitionCmd{
			Name: rv,
			Loc:  nodes.Loc{Start: partLoc, End: p.prev.End},
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
func (p *Parser) parseAlterTableSet() (*nodes.AlterTableCmd, error) {
	p.advance() // consume SET

	switch p.cur.Type {
	case WITHOUT:
		p.advance() // consume WITHOUT
		switch p.cur.Type {
		case CLUSTER:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropCluster)}, nil
		case OIDS:
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOids)}, nil
		default:
			// exhaustive-gap fix: SET WITHOUT <unknown>. Only CLUSTER / OIDS legal.
			return nil, p.syntaxErrorAtCur()
		}
	case LOGGED:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetLogged)}, nil
	case UNLOGGED:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_SetUnLogged)}, nil
	case ACCESS:
		p.advance() // consume ACCESS
		p.expect(METHOD)
		if p.cur.Type == DEFAULT {
			p.advance()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetAccessMethod),
				Name:    "",
			}, nil
		}
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetAccessMethod),
			Name:    name,
		}, nil
	case '(':
		// SET reloptions (starts with '(' because reloptions = '(' reloption_list ')')
		opts := p.parseReloptions()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetRelOptions),
			Def:     opts,
		}, nil
	case WITH:
		p.advance() // consume WITH
		if p.cur.Type == OIDS {
			p.advance()
			return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOids)}, nil
		}
		// exhaustive-gap fix: SET WITH <unknown>. Only OIDS legal.
		return nil, p.syntaxErrorAtCur()
	case TABLESPACE:
		p.advance()
		name, _ := p.parseName()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetTableSpace),
			Name:    name,
		}, nil
	default:
		// exhaustive-gap fix: SET <unknown>. No legal SET alternative matched.
		return nil, p.syntaxErrorAtCur()
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
func (p *Parser) parseAlterTableOf() (*nodes.AlterTableCmd, error) {
	p.advance() // consume OF
	names, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	if names == nil {
		return nil, p.syntaxErrorAtCur()
	}
	tn := makeTypeNameFromNameList(names)
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddOf),
		Def:     tn,
	}, nil
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
func (p *Parser) parseAlterIdentityColumnOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		opt, err := p.parseAlterIdentityColumnOption()
		if err != nil {
			return nil, err
		}
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterIdentityColumnOption parses alter_identity_column_option.
//
//	RESTART | RESTART opt_with NumericOnly | SET SeqOptElem | SET GENERATED generated_when
func (p *Parser) parseAlterIdentityColumnOption() (nodes.Node, error) {
	switch p.cur.Type {
	case RESTART:
		p.advance()
		// optional: WITH NumericOnly or just NumericOnly
		if p.cur.Type == WITH_LA || p.cur.Type == WITH {
			p.advance()
		}
		if p.cur.Type == ICONST || p.cur.Type == FCONST || p.cur.Type == '+' || p.cur.Type == '-' {
			val := p.parseNumericOnly()
			return makeDefElem("restart", val), nil
		}
		return makeDefElem("restart", nil), nil
	case SET:
		p.advance()
		if p.cur.Type == GENERATED {
			p.advance()
			gw := p.parseGeneratedWhen()
			return makeDefElem("generated", makeIntConst(int64(gw))), nil
		}
		// SET SeqOptElem
		return p.parseOneSeqOptElem()
	// optional-probe: alter_identity_column_option list — list caller breaks on nil after first match
	default:
		return nil, nil
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
			Loc:       nodes.NoLoc(),
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
		Defname: name,
		Arg:     &nodes.String{Str: arg},
		Loc:     nodes.NoLoc(),
	}
}

// parseOneSeqOptElem parses a single SeqOptElem using the existing parseSeqOptList.
func (p *Parser) parseOneSeqOptElem() (*nodes.DefElem, error) {
	list, err := p.parseSeqOptList()
	if err != nil {
		return nil, err
	}
	if list != nil && len(list.Items) > 0 {
		if de, ok := list.Items[0].(*nodes.DefElem); ok {
			return de, nil
		}
	}
	// exhaustive-gap fix: no SeqOptElem matched. Previously returned an empty
	// DefElem which downstream treated as an un-named sequence option.
	return nil, p.syntaxErrorAtCur()
}

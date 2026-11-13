package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
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
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_TABLE), alterLoc)
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
	rel, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}

	// Check for RENAME (produces RenameStmt, not AlterTableStmt)
	if p.cur.Type == RENAME {
		return p.parseAlterTableRename(rel, nodes.OBJECT_TABLE, missingOk, alterLoc)
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
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
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_INDEX), alterLoc)
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
	nameStart := p.pos()
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})

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
		partStart := p.pos()
		partNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && partNames == nil {
			return nil, p.syntaxErrorAtCur()
		}
		partRv := makeRangeVarFromAnyName(partNames, nodes.Loc{Start: partStart, End: p.prev.End})
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

	nameStart := p.pos()
	names, _ := p.parseQualifiedName()
	nameLoc := nodes.Loc{Start: nameStart, End: p.prev.End}

	// Check for RENAME (produces RenameStmt, uses makeRangeVarFromAnyName)
	if p.cur.Type == RENAME {
		rv := makeRangeVarFromAnyName(names, nameLoc)
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
		rv := makeRangeVarFromAnyName(names, nameLoc)
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
		rv := makeRangeVarFromNames(names, nameLoc)
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
		rv := makeRangeVarFromAnyName(names, nameLoc)
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
	rv := makeRangeVarFromNames(names, nameLoc)
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

	nameStart := p.pos()
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	// Check for RENAME
	if p.cur.Type == RENAME {
		return p.parseAlterViewRename(rv, missingOk, alterLoc, nodes.OBJECT_VIEW)
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
		return p.parseAlterTableMoveAll(int(nodes.OBJECT_MATVIEW), alterLoc)
	}

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	nameStart := p.pos()
	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	// Check for RENAME
	if p.cur.Type == RENAME {
		return p.parseAlterViewRename(rv, missingOk, alterLoc, nodes.OBJECT_MATVIEW)
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

// parseAlterViewRename parses the RENAME forms shared by ALTER VIEW and
// ALTER MATERIALIZED VIEW (gram.y RenameStmt):
//
//	ALTER [MATERIALIZED] VIEW qualified_name RENAME TO name
//	ALTER [MATERIALIZED] VIEW qualified_name RENAME opt_column name TO name
//
// COLUMN is optional in the column form, exactly as for ALTER TABLE.
// relType is OBJECT_VIEW or OBJECT_MATVIEW.
func (p *Parser) parseAlterViewRename(rv *nodes.RangeVar, missingOk bool, alterLoc int, relType nodes.ObjectType) (nodes.Node, error) {
	p.advance() // consume RENAME

	if p.collectMode() {
		p.addTokenCandidate(TO)
		p.addTokenCandidate(COLUMN)
		// Also valid: column name directly (RENAME col_name TO ...)
		p.addRuleCandidate("columnref")
		return nil, nil
	}

	// RENAME TO newname
	if p.cur.Type == TO {
		p.advance()
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: relType,
			Relation:   rv,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	}

	// RENAME [COLUMN] oldname TO newname
	if p.cur.Type == COLUMN {
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, nil
		}
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
	newname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.RenameStmt{
		RenameType:   nodes.OBJECT_COLUMN,
		RelationType: relType,
		Relation:     rv,
		Subname:      oldname,
		Newname:      newname,
		MissingOk:    missingOk,
		Loc:          nodes.Loc{Start: alterLoc, End: p.prev.End},
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

	rel, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}

	// Check for RENAME
	if p.cur.Type == RENAME {
		return p.parseAlterTableRename(rel, nodes.OBJECT_FOREIGN_TABLE, missingOk, alterLoc)
	}

	// Check for SET SCHEMA (produces AlterObjectSchemaStmt)
	if p.cur.Type == SET && p.peekNext().Type == SCHEMA {
		p.advance() // consume SET
		p.advance() // consume SCHEMA
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_FOREIGN_TABLE,
			Relation:   rel,
			Newschema:  newschema,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
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
func (p *Parser) parseAlterTableMoveAll(objType int, alterLoc int) (*nodes.AlterTableMoveAllStmt, error) {
	p.advance() // consume ALL
	if _, err := p.expect(IN_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(TABLESPACE); err != nil {
		return nil, err
	}
	origTs, err := p.parseName()
	if err != nil {
		return nil, err
	}

	var roles *nodes.List
	if p.cur.Type == OWNED {
		p.advance() // OWNED
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		roles, err = p.parseRoleList()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(SET); err != nil {
		return nil, err
	}
	if _, err := p.expect(TABLESPACE); err != nil {
		return nil, err
	}
	newTs, err := p.parseName()
	if err != nil {
		return nil, err
	}

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
	}, nil
}

// parseAlterTableRename parses the RENAME forms of ALTER TABLE and ALTER
// FOREIGN TABLE. relType is OBJECT_TABLE or OBJECT_FOREIGN_TABLE.
//
//	ALTER TABLE relation_expr RENAME TO name                         -- gram.y:9480
//	ALTER FOREIGN TABLE relation_expr RENAME TO name                 -- gram.y:9590
//	ALTER TABLE relation_expr RENAME opt_column name TO name         -- gram.y:9612
//	ALTER TABLE relation_expr RENAME CONSTRAINT name TO name         -- gram.y:9684
//	ALTER FOREIGN TABLE relation_expr RENAME opt_column name TO name -- gram.y:9706
//
// Each production has an IF_P EXISTS twin. The foreign-table productions
// set renameType (RENAME TO) or relationType (RENAME COLUMN) to
// OBJECT_FOREIGN_TABLE, and there is no RENAME CONSTRAINT production for a
// foreign table.
func (p *Parser) parseAlterTableRename(rel *nodes.RangeVar, relType nodes.ObjectType, missingOk bool, alterLoc int) (nodes.Node, error) {
	p.advance() // consume RENAME

	if p.collectMode() {
		p.addTokenCandidate(TO)
		p.addTokenCandidate(COLUMN)
		if relType == nodes.OBJECT_TABLE {
			p.addTokenCandidate(CONSTRAINT)
		}
		// Also valid: column name directly (RENAME col_name TO ...)
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}

	switch p.cur.Type {
	case TO:
		// RENAME TO name
		p.advance()
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: relType,
			Relation:   rel,
			Newname:    newname,
			MissingOk:  missingOk,
			Loc:        nodes.Loc{Start: alterLoc, End: p.prev.End},
		}, nil
	case CONSTRAINT:
		// RENAME CONSTRAINT name TO name
		if relType != nodes.OBJECT_TABLE {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume CONSTRAINT
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			return nil, errCollecting
		}
		oldname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
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
		// RENAME opt_column name TO name
		if p.cur.Type == COLUMN {
			p.advance() // consume COLUMN
			if p.collectMode() {
				p.addRuleCandidate("columnref")
				return nil, errCollecting
			}
		}
		oldname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType:   nodes.OBJECT_COLUMN,
			RelationType: relType,
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
				RENAME, IDENT,
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
		cmd, err = p.parseAlterTableAdd()
	case DROP:
		cmd, err = p.parseAlterTableDrop()
	case ALTER:
		cmd, err = p.parseAlterTableAlter()
	case OWNER:
		cmd, err = p.parseAlterTableOwner()
	case VALIDATE:
		cmd, err = p.parseAlterTableValidate()
	case INHERIT:
		cmd, err = p.parseAlterTableInherit()
	case NO:
		cmd, err = p.parseAlterTableNo()
	case ATTACH:
		cmd, err = p.parseAlterTableAttach()
	case DETACH:
		cmd, err = p.parseAlterTableDetach()
	case ENABLE_P:
		cmd, err = p.parseAlterTableEnable()
	case DISABLE_P:
		cmd, err = p.parseAlterTableDisable()
	case FORCE:
		cmd, err = p.parseAlterTableForce()
	case CLUSTER:
		cmd, err = p.parseAlterTableCluster()
	case SET:
		cmd, err = p.parseAlterTableSet()
	case RESET:
		cmd, err = p.parseAlterTableReset()
	case REPLICA:
		cmd, err = p.parseAlterTableReplica()
	case OF:
		cmd, err = p.parseAlterTableOf()
	case NOT:
		cmd, err = p.parseAlterTableNot()
	case OPTIONS:
		cmd, err = p.parseAlterTableOptions()
	default:
		if p.redshiftWordEqual("append") {
			cmd, err = p.parseRedshiftAlterTableAppend()
			break
		}
		if p.redshiftWordEqual("auto") {
			cmd, err = p.parseRedshiftAutoRefresh()
			break
		}
		if p.cur.Type == ROW || p.redshiftWordEqual("row") {
			cmd, err = p.parseRedshiftRowLevelSecurity()
			break
		}
		if p.redshiftWordEqual("masking") {
			cmd, err = p.parseRedshiftMaskingForDatashares()
			break
		}
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

// parseRedshiftAlterTableAppend parses Redshift's
// ALTER TABLE target APPEND FROM source [IGNOREEXTRA | FILLTARGET].
func (p *Parser) parseRedshiftAlterTableAppend() (*nodes.AlterTableCmd, error) {
	p.advance() // consume APPEND
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	source, err := p.parseRelationExpr()
	if err != nil {
		return nil, err
	}
	option := ""
	if p.redshiftWordEqual("ignoreextra") || p.redshiftWordEqual("filltarget") {
		option = strings.ToLower(p.advance().Str)
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_RedshiftAppend),
		Name:    option,
		Def:     source,
	}, nil
}

func (p *Parser) parseRedshiftAutoRefresh() (*nodes.AlterTableCmd, error) {
	p.advance() // consume AUTO
	if _, err := p.expect(REFRESH); err != nil {
		return nil, err
	}
	enabled, err := p.parseRedshiftYesNo()
	if err != nil {
		return nil, err
	}
	return makeRedshiftAlterOptionCmd("auto_refresh", []nodes.Node{
		makeDefElem("enabled", &nodes.Boolean{Boolval: enabled}),
	}), nil
}

func (p *Parser) parseRedshiftRowLevelSecurity() (*nodes.AlterTableCmd, error) {
	p.advance() // consume ROW
	if _, err := p.expect(LEVEL); err != nil {
		return nil, err
	}
	if _, err := p.expect(SECURITY); err != nil {
		return nil, err
	}
	enabled, err := p.parseRedshiftOnOff()
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{makeDefElem("enabled", &nodes.Boolean{Boolval: enabled})}
	if p.redshiftWordEqual("conjunction") {
		p.advance()
		if _, err := p.expect(TYPE_P); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case AND:
			p.advance()
			opts = append(opts, makeDefElem("conjunction_type", &nodes.String{Str: "and"}))
		case OR:
			p.advance()
			opts = append(opts, makeDefElem("conjunction_type", &nodes.String{Str: "or"}))
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == FOR {
		p.advance()
		if !p.consumeRedshiftWord("datashares") {
			return nil, p.syntaxErrorAtCur()
		}
		opts = append(opts, makeDefElem("for_datashares", &nodes.Boolean{Boolval: true}))
	}
	return makeRedshiftAlterOptionCmd("row_level_security", opts), nil
}

func (p *Parser) parseRedshiftMaskingForDatashares() (*nodes.AlterTableCmd, error) {
	p.advance() // consume MASKING
	enabled, err := p.parseRedshiftOnOff()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	if !p.consumeRedshiftWord("datashares") {
		return nil, p.syntaxErrorAtCur()
	}
	return makeRedshiftAlterOptionCmd("masking", []nodes.Node{
		makeDefElem("enabled", &nodes.Boolean{Boolval: enabled}),
		makeDefElem("for_datashares", &nodes.Boolean{Boolval: true}),
	}), nil
}

func (p *Parser) parseRedshiftAlterTableAddPartition(missingOk bool) (*nodes.AlterTableCmd, error) {
	p.advance() // consume PARTITION
	spec, err := p.parseRedshiftPartitionSpec()
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{makeDefElem("partition", spec)}
	if missingOk {
		opts = append(opts, makeDefElem("missing_ok", &nodes.Boolean{Boolval: true}))
	}
	opts, err = p.parseRedshiftAlterTableLocationInto(opts)
	if err != nil {
		return nil, err
	}
	return makeRedshiftAlterOptionCmd("add_partition", opts), nil
}

func (p *Parser) parseRedshiftAlterTableDropPartition() (*nodes.AlterTableCmd, error) {
	p.advance() // consume PARTITION
	spec, err := p.parseRedshiftPartitionSpec()
	if err != nil {
		return nil, err
	}
	return makeRedshiftAlterOptionCmd("drop_partition", []nodes.Node{
		makeDefElem("partition", spec),
	}), nil
}

func (p *Parser) parseRedshiftPartitionSpec() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	spec, err := p.parseDefList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return spec, nil
}

func (p *Parser) parseRedshiftAlterTableLocationInto(opts []nodes.Node) ([]nodes.Node, error) {
	if p.cur.Type != LOCATION {
		return opts, nil
	}
	p.advance()
	location, err := p.parseRedshiftStringLiteral()
	if err != nil {
		return nil, err
	}
	return append(opts, makeDefElem("location", &nodes.String{Str: location})), nil
}

func (p *Parser) parseRedshiftStringLiteral() (string, error) {
	if p.cur.Type != SCONST {
		return "", p.syntaxErrorAtCur()
	}
	value := p.cur.Str
	p.advance()
	return value, nil
}

func (p *Parser) parseRedshiftStringOptionList() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.Node
	for {
		name, err := p.parseRedshiftStringLiteral()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		value, err := p.parseRedshiftStringLiteral()
		if err != nil {
			return nil, err
		}
		items = append(items, makeDefElem(name, &nodes.String{Str: value}))
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseRedshiftYesNo() (bool, error) {
	switch {
	case p.cur.Type == YES_P || p.redshiftWordEqual("yes"):
		p.advance()
		return true, nil
	case p.cur.Type == NO || p.redshiftWordEqual("no"):
		p.advance()
		return false, nil
	default:
		return false, p.syntaxErrorAtCur()
	}
}

func (p *Parser) parseRedshiftOnOff() (bool, error) {
	switch {
	case p.cur.Type == ON || p.redshiftWordEqual("on"):
		p.advance()
		return true, nil
	case p.cur.Type == OFF || p.redshiftWordEqual("off"):
		p.advance()
		return false, nil
	default:
		return false, p.syntaxErrorAtCur()
	}
}

func makeRedshiftAlterOptionCmd(name string, opts []nodes.Node) *nodes.AlterTableCmd {
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_RedshiftAlterOption),
		Name:    name,
		Def:     &nodes.List{Items: opts},
	}
}

// parseAlterTableAdd handles ADD ... subcommands.
func (p *Parser) parseAlterTableAdd() (*nodes.AlterTableCmd, error) {
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
		return nil, nil
	}

	if p.cur.Type == PARTITION {
		return p.parseRedshiftAlterTableAddPartition(false)
	}

	if p.startsTableConstraint() {
		// ADD TableConstraint
		constr, err := p.parseTableConstraint()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AddConstraint),
			Def:     constr,
		}, nil
	}
	switch p.cur.Type {
	case COLUMN:
		p.advance() // consume COLUMN
		// IF NOT EXISTS
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			if err := p.expectNotExists(); err != nil {
				return nil, err
			}
			missingOk = true
		}
		coldef, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_AddColumn),
			Def:        coldef,
			Missing_ok: missingOk,
		}, nil
	case IF_P:
		p.advance()
		if err := p.expectNotExists(); err != nil {
			return nil, err
		}
		if p.cur.Type == PARTITION {
			return p.parseRedshiftAlterTableAddPartition(true)
		}
		// ADD IF NOT EXISTS columnDef (without COLUMN keyword)
		coldef, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype:    int(nodes.AT_AddColumn),
			Def:        coldef,
			Missing_ok: true,
		}, nil
	default:
		// ADD columnDef (without COLUMN keyword)
		// Try to distinguish between column def and constraint.
		// If current token looks like a column name followed by a type, it's a column def.
		coldef, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AddColumn),
			Def:     coldef,
		}, nil
	}
}

// parseAlterTableDrop handles DROP ... subcommands.
//
//	DROP opt_column IF_P EXISTS ColId opt_drop_behavior
//	DROP opt_column ColId opt_drop_behavior
//	DROP CONSTRAINT IF_P EXISTS name opt_drop_behavior
//	DROP CONSTRAINT name opt_drop_behavior
func (p *Parser) parseAlterTableDrop() (*nodes.AlterTableCmd, error) {
	p.advance() // consume DROP

	if p.collectMode() {
		p.cachedCollect("parseAlterTableDrop", func() {
			for _, t := range []int{COLUMN, CONSTRAINT, IF_P} {
				p.addTokenCandidate(t)
			}
			// Also valid: column name directly (DROP col_name)
			p.addRuleCandidate("columnref")
		})
		return nil, nil
	}

	if p.cur.Type == PARTITION {
		return p.parseRedshiftAlterTableDropPartition()
	}

	subtype := nodes.AT_DropColumn
	switch p.cur.Type {
	case COLUMN:
		p.advance() // consume COLUMN
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addTokenCandidate(IF_P)
			return nil, nil
		}
	case CONSTRAINT:
		p.advance() // consume CONSTRAINT
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
			p.addTokenCandidate(IF_P)
			return nil, nil
		}
		subtype = nodes.AT_DropConstraint
	}

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}
	// ColId and name are the same production.
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype:    int(subtype),
		Name:       name,
		Behavior:   p.parseOptDropBehavior(),
		Missing_ok: missingOk,
	}, nil
}

// parseAlterTableAlter handles ALTER [COLUMN] ... subcommands.
func (p *Parser) parseAlterTableAlter() (*nodes.AlterTableCmd, error) {
	p.advance() // consume ALTER

	if cmd, err := p.parseRedshiftAlterTableStorageOption(); cmd != nil || err != nil {
		return cmd, err
	}

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
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		// ConstraintAttributeSpec (we consume but don't store — matches
		// yacc behavior). Only a foreign key's deferrability changes.
		if _, err := p.parseKindAttributes("FOREIGN KEY", true, false, false); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AlterConstraint),
			Name:    name,
		}, nil
	}

	// ALTER opt_column Iconst SET STATISTICS set_statistics_value
	if p.cur.Type == ICONST {
		num := int16(p.cur.Ival)
		p.advance()
		if _, err := p.expect(SET); err != nil {
			return nil, err
		}
		if _, err := p.expect(STATISTICS); err != nil {
			return nil, err
		}
		val, err := p.parseSetStatisticsValue()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStatistics),
			Num:     num,
			Def:     val,
		}, nil
	}

	// Regular column name
	colname, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	return p.parseAlterColumnAction(colname)
}

func (p *Parser) parseRedshiftAlterTableStorageOption() (*nodes.AlterTableCmd, error) {
	switch {
	case p.redshiftWordEqual("distkey"):
		p.advance()
		column, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return makeRedshiftAlterOptionCmd("distkey", []nodes.Node{
			makeDefElem("column", &nodes.String{Str: column}),
		}), nil
	case p.redshiftWordEqual("diststyle"):
		p.advance()
		style, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		style = strings.ToLower(style)
		opts := []nodes.Node{makeDefElem("style", &nodes.String{Str: style})}
		if style == "key" && p.consumeRedshiftWord("distkey") {
			column, err := p.parseColId()
			if err != nil {
				return nil, err
			}
			opts = append(opts, makeDefElem("distkey", &nodes.String{Str: column}))
		}
		return makeRedshiftAlterOptionCmd("diststyle", opts), nil
	case p.redshiftWordEqual("encode"):
		p.advance()
		value, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		return makeRedshiftAlterOptionCmd("encode", []nodes.Node{
			makeDefElem("value", &nodes.String{Str: strings.ToLower(value)}),
		}), nil
	case p.redshiftWordEqual("compound"), p.redshiftWordEqual("interleaved"), p.redshiftWordEqual("sortkey"):
		return p.parseRedshiftAlterSortkey()
	// optional-probe: Redshift ALTER TABLE storage options are optional before ALTER [COLUMN] fallback.
	default:
		return nil, nil
	}
}

func (p *Parser) parseRedshiftAlterSortkey() (*nodes.AlterTableCmd, error) {
	opts := []nodes.Node{}
	if p.redshiftWordEqual("compound") || p.redshiftWordEqual("interleaved") {
		style := strings.ToLower(p.advance().Str)
		opts = append(opts, makeDefElem("style", &nodes.String{Str: style}))
	}
	if !p.consumeRedshiftWord("sortkey") {
		return nil, p.syntaxErrorAtCur()
	}
	if p.redshiftWordEqual("auto") || p.redshiftWordEqual("none") {
		value := strings.ToLower(p.advance().Str)
		opts = append(opts, makeDefElem("value", &nodes.String{Str: value}))
		return makeRedshiftAlterOptionCmd("sortkey", opts), nil
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	cols, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opts = append(opts, makeDefElem("columns", cols))
	return makeRedshiftAlterOptionCmd("sortkey", opts), nil
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
		return p.parseAlterColumnDrop(colname)
	case TYPE_P:
		return p.parseAlterColumnType(colname, false)
	case ADD_P:
		p.advance() // consume ADD
		return p.parseAlterColumnAddGenerated(colname)
	case RESET:
		// ALTER opt_column ColId RESET reloptions
		p.advance() // consume RESET
		defs, err := p.parseReloptions()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ResetOptions),
			Name:    colname,
			Def:     defs,
		}, nil
	case OPTIONS:
		// alter_generic_options
		opts, err := p.parseAlterGenericOptions()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_AlterColumnGenericOptions),
			Name:    colname,
			Def:     opts,
		}, nil
	default:
		if p.consumeRedshiftWord("encode") {
			encoding, err := p.parseRedshiftOptionWord()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetCompression),
				Name:    colname,
				Def:     &nodes.String{Str: strings.ToLower(encoding)},
			}, nil
		}
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
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ColumnDefault),
			Name:    colname,
			Def:     expr,
		}, nil
	case NOT:
		// SET NOT NULL
		p.advance() // NOT
		if _, err := p.expect(NULL_P); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetNotNull),
			Name:    colname,
		}, nil
	case DATA_P:
		// SET DATA TYPE Typename ...
		p.advance() // DATA
		if _, err := p.expect(TYPE_P); err != nil {
			return nil, err
		}
		return p.parseAlterColumnTypeInner(colname)
	case STATISTICS:
		// SET STATISTICS set_statistics_value
		p.advance()
		val, err := p.parseSetStatisticsValue()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStatistics),
			Name:    colname,
			Def:     val,
		}, nil
	case STORAGE:
		// SET column_storage: STORAGE ColId | STORAGE DEFAULT
		p.advance()
		storageVal, err := p.parseColIdOrDefault()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetStorage),
			Name:    colname,
			Def:     &nodes.String{Str: storageVal},
		}, nil
	case COMPRESSION:
		// SET column_compression: COMPRESSION ColId | COMPRESSION DEFAULT
		p.advance()
		compVal, err := p.parseColIdOrDefault()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetCompression),
			Name:    colname,
			Def:     &nodes.String{Str: compVal},
		}, nil
	case EXPRESSION:
		// SET EXPRESSION AS '(' a_expr ')'
		p.advance()
		if _, err := p.expect(AS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetExpression),
			Name:    colname,
			Def:     expr,
		}, nil
	case GENERATED:
		// SET GENERATED generated_when
		p.advance()
		gw, err := p.parseGeneratedWhen()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetIdentity),
			Name:    colname,
			Def: &nodes.List{Items: []nodes.Node{
				makeDefElem("generated", makeIntConst(int64(gw))),
			}},
		}, nil
	case '(':
		// SET reloptions
		defs, err := p.parseReloptions()
		if err != nil {
			return nil, err
		}
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
//
//	ALTER opt_column ColId DROP DEFAULT
//	ALTER opt_column ColId DROP NOT NULL_P
//	ALTER opt_column ColId DROP EXPRESSION [IF_P EXISTS]
//	ALTER opt_column ColId DROP IDENTITY_P [IF_P EXISTS]
func (p *Parser) parseAlterColumnDrop(colname string) (*nodes.AlterTableCmd, error) {
	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_ColumnDefault),
			Name:    colname,
		}, nil
	case NOT:
		p.advance() // NOT
		if _, err := p.expect(NULL_P); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DropNotNull),
			Name:    colname,
		}, nil
	case EXPRESSION, IDENTITY_P:
		subtype := nodes.AT_DropExpression
		if p.cur.Type == IDENTITY_P {
			subtype = nodes.AT_DropIdentity
		}
		p.advance()
		missingOk := false
		if p.cur.Type == IF_P {
			p.advance()
			if _, err := p.expect(EXISTS); err != nil {
				return nil, err
			}
			missingOk = true
		}
		return &nodes.AlterTableCmd{
			Subtype:    int(subtype),
			Name:       colname,
			Missing_ok: missingOk,
		}, nil
	default:
		return nil, p.syntaxErrorAtCur()
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
	if _, err := p.expect(GENERATED); err != nil {
		return nil, err
	}
	gw, err := p.parseGeneratedWhen()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	if _, err := p.expect(IDENTITY_P); err != nil {
		return nil, err
	}
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
func (p *Parser) parseAlterTableOwner() (*nodes.AlterTableCmd, error) {
	p.advance() // consume OWNER
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	roleSpec, err := p.parseRoleSpec()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype:  int(nodes.AT_ChangeOwner),
		Newowner: roleSpec,
	}, nil
}

// parseAlterTableValidate handles VALIDATE CONSTRAINT name.
func (p *Parser) parseAlterTableValidate() (*nodes.AlterTableCmd, error) {
	p.advance() // consume VALIDATE
	if _, err := p.expect(CONSTRAINT); err != nil {
		return nil, err
	}
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ValidateConstraint),
		Name:    name,
	}, nil
}

// parseAlterTableInherit handles INHERIT qualified_name.
func (p *Parser) parseAlterTableInherit() (*nodes.AlterTableCmd, error) {
	p.advance() // consume INHERIT
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if names == nil {
		return nil, p.syntaxErrorAtCur()
	}
	rv := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AddInherit),
		Def:     rv,
	}, nil
}

// parseAlterTableNo handles NO INHERIT and NO FORCE ROW LEVEL SECURITY.
//
//	NO INHERIT qualified_name
//	NO FORCE ROW LEVEL SECURITY
func (p *Parser) parseAlterTableNo() (*nodes.AlterTableCmd, error) {
	p.advance() // consume NO

	switch p.cur.Type {
	case INHERIT:
		p.advance()
		rv, err := p.parseAlterTableQualifiedName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_DropInherit),
			Def:     rv,
		}, nil
	case FORCE:
		p.advance() // FORCE
		if err := p.expectRowLevelSecurity(); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_NoForceRowSecurity),
		}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterTableAttach handles ATTACH PARTITION qualified_name ForValues.
func (p *Parser) parseAlterTableAttach() (*nodes.AlterTableCmd, error) {
	partLoc := p.pos()
	p.advance() // consume ATTACH
	if _, err := p.expect(PARTITION); err != nil {
		return nil, err
	}
	rv, err := p.parseAlterTableQualifiedName()
	if err != nil {
		return nil, err
	}
	bound, err := p.parseForValues()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_AttachPartition),
		Def: &nodes.PartitionCmd{
			Name:  rv,
			Bound: bound,
			Loc:   nodes.Loc{Start: partLoc, End: p.prev.End},
		},
	}, nil
}

// parseAlterTableDetach handles DETACH PARTITION ...
//
//	DETACH PARTITION qualified_name opt_concurrently
//	DETACH PARTITION qualified_name FINALIZE
func (p *Parser) parseAlterTableDetach() (*nodes.AlterTableCmd, error) {
	partLoc := p.pos()
	p.advance() // consume DETACH
	if _, err := p.expect(PARTITION); err != nil {
		return nil, err
	}
	rv, err := p.parseAlterTableQualifiedName()
	if err != nil {
		return nil, err
	}

	subtype := nodes.AT_DetachPartition
	concurrent := false
	switch p.cur.Type {
	case CONCURRENTLY:
		p.advance()
		concurrent = true
	case FINALIZE:
		p.advance()
		subtype = nodes.AT_DetachPartitionFinalize
	}
	return &nodes.AlterTableCmd{
		Subtype: int(subtype),
		Def: &nodes.PartitionCmd{
			Name:       rv,
			Concurrent: concurrent,
			Loc:        nodes.Loc{Start: partLoc, End: p.prev.End},
		},
	}, nil
}

// parseAlterTableQualifiedName parses the qualified_name of INHERIT, NO
// INHERIT, ATTACH PARTITION and DETACH PARTITION into a RangeVar.
func (p *Parser) parseAlterTableQualifiedName() (*nodes.RangeVar, error) {
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End}), nil
}

// parseAlterTableEnable handles ENABLE TRIGGER/RULE/ROW LEVEL SECURITY.
//
//	ENABLE_P TRIGGER name | ENABLE_P TRIGGER ALL | ENABLE_P TRIGGER USER
//	ENABLE_P ALWAYS TRIGGER name | ENABLE_P REPLICA TRIGGER name
//	ENABLE_P RULE name | ENABLE_P ALWAYS RULE name | ENABLE_P REPLICA RULE name
//	ENABLE_P ROW LEVEL SECURITY
func (p *Parser) parseAlterTableEnable() (*nodes.AlterTableCmd, error) {
	p.advance() // consume ENABLE

	switch p.cur.Type {
	case TRIGGER:
		p.advance()
		return p.parseAlterTableTriggerTarget(nodes.AT_EnableTrig, nodes.AT_EnableTrigAll, nodes.AT_EnableTrigUser)
	case ALWAYS, REPLICA:
		trig, rule := nodes.AT_EnableAlwaysTrig, nodes.AT_EnableAlwaysRule
		if p.cur.Type == REPLICA {
			trig, rule = nodes.AT_EnableReplicaTrig, nodes.AT_EnableReplicaRule
		}
		p.advance() // consume ALWAYS or REPLICA
		subtype := rule
		if p.cur.Type == TRIGGER {
			subtype = trig
			p.advance()
		} else if _, err := p.expect(RULE); err != nil {
			return nil, err
		}
		return p.parseAlterTableNamedCmd(subtype)
	case RULE:
		p.advance()
		return p.parseAlterTableNamedCmd(nodes.AT_EnableRule)
	case ROW:
		p.advance() // ROW
		if err := p.expectLevelSecurity(); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_EnableRowSecurity)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterTableDisable handles DISABLE TRIGGER/RULE/ROW LEVEL SECURITY.
//
//	DISABLE_P TRIGGER name | DISABLE_P TRIGGER ALL | DISABLE_P TRIGGER USER
//	DISABLE_P RULE name
//	DISABLE_P ROW LEVEL SECURITY
func (p *Parser) parseAlterTableDisable() (*nodes.AlterTableCmd, error) {
	p.advance() // consume DISABLE

	switch p.cur.Type {
	case TRIGGER:
		p.advance()
		return p.parseAlterTableTriggerTarget(nodes.AT_DisableTrig, nodes.AT_DisableTrigAll, nodes.AT_DisableTrigUser)
	case RULE:
		p.advance()
		return p.parseAlterTableNamedCmd(nodes.AT_DisableRule)
	case ROW:
		p.advance() // ROW
		if err := p.expectLevelSecurity(); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DisableRowSecurity)}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterTableTriggerTarget parses what follows ENABLE_P TRIGGER or
// DISABLE_P TRIGGER: ALL, USER, or a trigger name.
func (p *Parser) parseAlterTableTriggerTarget(named, all, user nodes.AlterTableType) (*nodes.AlterTableCmd, error) {
	switch p.cur.Type {
	case ALL:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(all)}, nil
	case USER:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(user)}, nil
	default:
		return p.parseAlterTableNamedCmd(named)
	}
}

// parseAlterTableNamedCmd parses the name that ends a trigger, rule or
// constraint command and returns the command of the given subtype.
func (p *Parser) parseAlterTableNamedCmd(subtype nodes.AlterTableType) (*nodes.AlterTableCmd, error) {
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype: int(subtype),
		Name:    name,
	}, nil
}

// parseAlterTableForce handles FORCE ROW LEVEL SECURITY.
func (p *Parser) parseAlterTableForce() (*nodes.AlterTableCmd, error) {
	p.advance() // consume FORCE
	if err := p.expectRowLevelSecurity(); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ForceRowSecurity)}, nil
}

// expectRowLevelSecurity consumes ROW LEVEL SECURITY.
func (p *Parser) expectRowLevelSecurity() error {
	if _, err := p.expect(ROW); err != nil {
		return err
	}
	return p.expectLevelSecurity()
}

// expectLevelSecurity consumes LEVEL SECURITY, the rest of ROW LEVEL SECURITY.
func (p *Parser) expectLevelSecurity() error {
	if _, err := p.expect(LEVEL); err != nil {
		return err
	}
	_, err := p.expect(SECURITY)
	return err
}

// parseAlterTableCluster handles CLUSTER ON name.
func (p *Parser) parseAlterTableCluster() (*nodes.AlterTableCmd, error) {
	p.advance() // consume CLUSTER
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	return p.parseAlterTableNamedCmd(nodes.AT_ClusterOn)
}

// parseAlterTableSet handles SET ... subcommands at the table level.
func (p *Parser) parseAlterTableSet() (*nodes.AlterTableCmd, error) {
	p.advance() // consume SET

	switch p.cur.Type {
	case TABLE:
		p.advance()
		if !p.consumeRedshiftWord("properties") {
			return nil, p.syntaxErrorAtCur()
		}
		props, err := p.parseRedshiftStringOptionList()
		if err != nil {
			return nil, err
		}
		return makeRedshiftAlterOptionCmd("table_properties", []nodes.Node{
			makeDefElem("properties", props),
		}), nil
	case LOCATION:
		p.advance()
		location, err := p.parseRedshiftStringLiteral()
		if err != nil {
			return nil, err
		}
		return makeRedshiftAlterOptionCmd("location", []nodes.Node{
			makeDefElem("location", &nodes.String{Str: location}),
		}), nil
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
		// SET ACCESS METHOD set_access_method_name: ColId | DEFAULT
		p.advance() // consume ACCESS
		if _, err := p.expect(METHOD); err != nil {
			return nil, err
		}
		if p.cur.Type == DEFAULT {
			p.advance()
			return &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetAccessMethod),
				Name:    "",
			}, nil
		}
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Subtype: int(nodes.AT_SetAccessMethod),
			Name:    name,
		}, nil
	case '(':
		// SET reloptions (starts with '(' because reloptions = '(' reloption_list ')')
		opts, err := p.parseReloptions()
		if err != nil {
			return nil, err
		}
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
		return p.parseAlterTableNamedCmd(nodes.AT_SetTableSpace)
	default:
		if p.consumeRedshiftWord("file") {
			if _, err := p.expect(FORMAT); err != nil {
				return nil, err
			}
			format, err := p.parseRedshiftOptionWord()
			if err != nil {
				return nil, err
			}
			return makeRedshiftAlterOptionCmd("file_format", []nodes.Node{
				makeDefElem("format", &nodes.String{Str: strings.ToLower(format)}),
			}), nil
		}
		// exhaustive-gap fix: SET <unknown>. No legal SET alternative matched.
		return nil, p.syntaxErrorAtCur()
	}
}

// parseAlterTableReset handles RESET reloptions.
func (p *Parser) parseAlterTableReset() (*nodes.AlterTableCmd, error) {
	p.advance() // consume RESET
	opts, err := p.parseReloptions()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_ResetRelOptions),
		Def:     opts,
	}, nil
}

// parseAlterTableReplica handles REPLICA IDENTITY_P replica_identity.
//
//	replica_identity: NOTHING | FULL | DEFAULT | USING INDEX name
func (p *Parser) parseAlterTableReplica() (*nodes.AlterTableCmd, error) {
	p.advance() // consume REPLICA
	if _, err := p.expect(IDENTITY_P); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case DEFAULT, FULL, NOTHING:
		p.advance()
		return &nodes.AlterTableCmd{Subtype: int(nodes.AT_ReplicaIdentity)}, nil
	case USING:
		p.advance() // USING
		if _, err := p.expect(INDEX); err != nil {
			return nil, err
		}
		return p.parseAlterTableNamedCmd(nodes.AT_ReplicaIdentity)
	default:
		return nil, p.syntaxErrorAtCur()
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
func (p *Parser) parseAlterTableNot() (*nodes.AlterTableCmd, error) {
	p.advance() // consume NOT
	if _, err := p.expect(OF); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{Subtype: int(nodes.AT_DropOf)}, nil
}

// parseAlterTableOptions handles bare alter_generic_options (for foreign tables).
func (p *Parser) parseAlterTableOptions() (*nodes.AlterTableCmd, error) {
	opts, err := p.parseAlterGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Subtype: int(nodes.AT_GenericOptions),
		Def:     opts,
	}, nil
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
//	RoleSpec: NonReservedWord | CURRENT_ROLE | CURRENT_USER | SESSION_USER
func (p *Parser) parseRoleSpec() (*nodes.RoleSpec, error) {
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
		name, err := p.parseNonReservedWord()
		if err != nil {
			return nil, err
		}
		rs = &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CSTRING),
			Rolename: name,
		}
	}
	rs.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return rs, nil
}

// parseRoleList parses role_list: RoleSpec [, RoleSpec ...]
func (p *Parser) parseRoleList() (*nodes.List, error) {
	rs, err := p.parseRoleSpec()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{rs}
	for p.cur.Type == ',' {
		p.advance()
		rs, err = p.parseRoleSpec()
		if err != nil {
			return nil, err
		}
		items = append(items, rs)
	}
	return &nodes.List{Items: items}, nil
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
func (p *Parser) parseGeneratedWhen() (byte, error) {
	if p.cur.Type == ALWAYS {
		p.advance()
		return 'a', nil
	}
	if _, err := p.expect(BY); err != nil {
		return 0, err
	}
	if _, err := p.expect(DEFAULT); err != nil {
		return 0, err
	}
	return 'd', nil
}

// parseSetStatisticsValue parses set_statistics_value. DEFAULT yields a nil
// node, as in gram.y.
//
//	set_statistics_value: SignedIconst | DEFAULT
func (p *Parser) parseSetStatisticsValue() (nodes.Node, error) {
	if p.cur.Type == DEFAULT {
		p.advance()
		return nil, nil
	}
	// SignedIconst: Iconst | '+' Iconst | '-' Iconst
	sign := p.cur.Type
	if sign == '+' || sign == '-' {
		p.advance()
	}
	num, err := p.expect(ICONST)
	if err != nil {
		return nil, err
	}
	if sign == '-' {
		return makeIntConst(-num.Ival), nil
	}
	return makeIntConst(num.Ival), nil
}

// parseColIdOrDefault parses the argument of column_storage and
// column_compression: ColId, or DEFAULT read as "default".
func (p *Parser) parseColIdOrDefault() (string, error) {
	if p.cur.Type == DEFAULT {
		p.advance()
		return "default", nil
	}
	return p.parseColId()
}

// expectNotExists consumes NOT EXISTS, the rest of IF NOT EXISTS after IF.
func (p *Parser) expectNotExists() error {
	if _, err := p.expect(NOT); err != nil {
		return err
	}
	_, err := p.expect(EXISTS)
	return err
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
		return nil, p.syntaxErrorAtCur()
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
		// opt_with NumericOnly: after WITH the number is required.
		hasWith := p.cur.Type == WITH_LA || p.cur.Type == WITH
		if hasWith {
			p.advance()
		}
		if p.cur.Type == ICONST || p.cur.Type == FCONST || p.cur.Type == '+' || p.cur.Type == '-' {
			val := p.parseNumericOnly()
			return makeDefElem("restart", val), nil
		}
		if hasWith {
			return nil, p.syntaxErrorAtCur()
		}
		return makeDefElem("restart", nil), nil
	case SET:
		p.advance()
		if p.cur.Type == GENERATED {
			p.advance()
			gw, err := p.parseGeneratedWhen()
			if err != nil {
				return nil, err
			}
			return makeDefElem("generated", makeIntConst(int64(gw))), nil
		}
		// SET SeqOptElem. gram.y rejects AS, RESTART and OWNED BY here
		// with a syntax error at the option.
		optTok := p.cur
		elem, err := p.parseOneSeqOptElem()
		if err != nil {
			return nil, err
		}
		switch elem.Defname {
		case "as", "restart", "owned_by":
			return nil, p.syntaxErrorAtTok(optTok)
		}
		return elem, nil
	// optional-probe: alter_identity_column_option list — list caller breaks on nil after first match
	default:
		return nil, nil
	}
}

// parseAlterGenericOptions parses alter_generic_options:
//
//	OPTIONS '(' alter_generic_option_list ')'
func (p *Parser) parseAlterGenericOptions() (*nodes.List, error) {
	if _, err := p.expect(OPTIONS); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseAlterGenericOptionList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseAlterGenericOptionList parses alter_generic_option_list.
func (p *Parser) parseAlterGenericOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		elem, err := p.parseAlterGenericOptionElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseAlterGenericOptionElem parses alter_generic_option_elem.
//
//	generic_option_elem
//	| SET generic_option_elem
//	| ADD_P generic_option_elem
//	| DROP generic_option_name
func (p *Parser) parseAlterGenericOptionElem() (*nodes.DefElem, error) {
	// SET, ADD and DROP are also ColLabels: followed by the Sconst of
	// generic_option_arg, the keyword is the option's name.
	if p.peekNext().Type == SCONST {
		return p.parseGenericOptionElem()
	}
	switch p.cur.Type {
	case SET, ADD_P:
		action := nodes.DEFELEM_SET
		if p.cur.Type == ADD_P {
			action = nodes.DEFELEM_ADD
		}
		p.advance()
		n, err := p.parseGenericOptionElem()
		if err != nil {
			return nil, err
		}
		n.Defaction = int(action)
		return n, nil
	case DROP:
		p.advance()
		name, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		return &nodes.DefElem{
			Defname:   name,
			Defaction: int(nodes.DEFELEM_DROP),
			Loc:       nodes.NoLoc(),
		}, nil
	default:
		return p.parseGenericOptionElem()
	}
}

// parseGenericOptionElem parses generic_option_elem.
//
//	generic_option_elem: generic_option_name generic_option_arg
//	generic_option_name: ColLabel
//	generic_option_arg: Sconst
func (p *Parser) parseGenericOptionElem() (*nodes.DefElem, error) {
	name, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	arg, err := p.expect(SCONST)
	if err != nil {
		return nil, err
	}
	return &nodes.DefElem{
		Defname: name,
		Arg:     &nodes.String{Str: arg.Str},
		Loc:     nodes.NoLoc(),
	}, nil
}

// parseOneSeqOptElem parses the single SeqOptElem after SET in
// alter_identity_column_option. A following option is the next
// alter_identity_column_option, so it is not read here.
func (p *Parser) parseOneSeqOptElem() (*nodes.DefElem, error) {
	elem, err := p.parseSeqOptElem()
	if err != nil {
		return nil, err
	}
	if elem == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return elem, nil
}

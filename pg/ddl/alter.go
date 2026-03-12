package ddl

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// atCmdPass maps AlterTableType to execution pass.
// PG executes ALTER TABLE subcommands in pass order, not statement order.
// Mirrors PG's AT_PASS_* constants:
//
//	0: AT_PASS_DROP      — drops (AT_DropColumn, AT_DropConstraint)
//	1: AT_PASS_ALTER_TYPE — type changes (AT_AlterColumnType)
//	3: AT_PASS_ADD_COL   — add column
//	5: AT_PASS_ADD_CONSTR — add constraint
//	6: AT_PASS_COL_ATTRS  — column attributes (NOT NULL, DEFAULT, etc.)
//	7: AT_PASS_ADD_INDEXCONSTR — add index constraint (PK/UNIQUE via index)
//	8: AT_PASS_ADD_INDEX — add index
//	10: AT_PASS_MISC      — miscellaneous
//
// pg: src/backend/commands/tablecmds.c — ATRewriteTables pass ordering
func atCmdPass(subtype nodes.AlterTableType) int {
	switch subtype {
	case nodes.AT_DropColumn, nodes.AT_DropConstraint:
		return 0 // AT_PASS_DROP
	case nodes.AT_AlterColumnType:
		return 1 // AT_PASS_ALTER_TYPE
	case nodes.AT_AddColumn:
		return 3 // AT_PASS_ADD_COL
	case nodes.AT_AddConstraint:
		return 5 // AT_PASS_ADD_CONSTR
	case nodes.AT_SetNotNull, nodes.AT_DropNotNull, nodes.AT_ColumnDefault:
		return 6 // AT_PASS_COL_ATTRS
	case nodes.AT_AddIndexConstraint:
		return 7 // AT_PASS_ADD_INDEXCONSTR
	case nodes.AT_AddIndex:
		return 8 // AT_PASS_ADD_INDEX
	default:
		return 10 // AT_PASS_MISC
	}
}

// pendingATCmd holds an ALTER TABLE subcommand with its execution pass and original order.
type pendingATCmd struct {
	pass     int
	origIdx  int
	cmd      *nodes.AlterTableCmd
}

// AlterTableStmt applies ALTER TABLE commands from a pgparser AST.
// Commands are sorted into passes before execution (DROP first, then TYPE,
// then ADD, then attribute changes), mirroring PG's multi-pass approach.
//
// pg: src/backend/commands/tablecmds.c — AlterTable
func (c *Catalog) AlterTableStmt(stmt *nodes.AlterTableStmt) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	schema, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		return err
	}
	// Validate the object type matches the relation kind.
	// pg: src/backend/commands/tablecmds.c — AlterTableGetRelIdAndLock (relkind check)
	objType := nodes.ObjectType(stmt.ObjType)
	switch objType {
	case nodes.OBJECT_TABLE:
		if rel.RelKind != 'r' && rel.RelKind != 'p' && rel.RelKind != 'f' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a table", relName)}
		}
	case nodes.OBJECT_VIEW:
		if rel.RelKind != 'v' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a view", relName)}
		}
	case nodes.OBJECT_MATVIEW:
		if rel.RelKind != 'm' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a materialized view", relName)}
		}
	case nodes.OBJECT_FOREIGN_TABLE:
		if rel.RelKind != 'f' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a foreign table", relName)}
		}
	case nodes.OBJECT_TYPE:
		if rel.RelKind != 'c' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a composite type", relName)}
		}
	default:
		// Legacy/unspecified: allow regular tables and partitioned tables.
		if rel.RelKind != 'r' && rel.RelKind != 'p' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a table", relName)}
		}
	}

	if stmt.Cmds == nil {
		return nil
	}

	// Collect commands with pass numbers.
	var pending []pendingATCmd
	for i, item := range stmt.Cmds.Items {
		atc, ok := item.(*nodes.AlterTableCmd)
		if !ok {
			continue
		}
		pending = append(pending, pendingATCmd{
			pass:    atCmdPass(nodes.AlterTableType(atc.Subtype)),
			origIdx: i,
			cmd:     atc,
		})
	}

	// Sort by pass, preserving original order within each pass.
	sort.SliceStable(pending, func(i, j int) bool {
		return pending[i].pass < pending[j].pass
	})

	// Check for duplicate ALTER TYPE on the same column.
	// pg: src/backend/commands/tablecmds.c:13197-13203
	typeChangeCols := make(map[string]bool)
	for _, p := range pending {
		if nodes.AlterTableType(p.cmd.Subtype) == nodes.AT_AlterColumnType {
			if typeChangeCols[p.cmd.Name] {
				return &Error{Code: CodeInvalidObjectDefinition,
					Message: fmt.Sprintf("cannot alter type of column \"%s\" twice", p.cmd.Name)}
			}
			typeChangeCols[p.cmd.Name] = true
		}
	}

	// Execute in pass order.
	for _, p := range pending {
		if err := c.execAlterTableCmd(schema, rel, relName, p.cmd); err != nil {
			return err
		}
	}

	// After executing on the main relation, propagate to children if applicable.
	// pg: src/backend/commands/tablecmds.c — ATRewriteTables (recurse to children)
	children := c.findChildRelations(rel.OID)
	if len(children) > 0 {
		for _, p := range pending {
			if shouldRecurse(nodes.AlterTableType(p.cmd.Subtype)) {
				for _, childOID := range children {
					childRel := c.relationByOID[childOID]
					if childRel == nil {
						continue
					}
					_ = c.execAlterTableCmd(childRel.Schema, childRel, childRel.Name, p.cmd)
				}
			}
		}
	}

	return nil
}

// execAlterTableCmd executes a single ALTER TABLE subcommand.
func (c *Catalog) execAlterTableCmd(schema *Schema, rel *Relation, relName string, atc *nodes.AlterTableCmd) error {
	cascade := atc.Behavior == int(nodes.DROP_CASCADE)

	switch nodes.AlterTableType(atc.Subtype) {
	case nodes.AT_AddColumn:
		cd, ok := atc.Def.(*nodes.ColumnDef)
		if !ok {
			return fmt.Errorf("AT_AddColumn: expected ColumnDef")
		}
		colDef, inlineCons, err := c.convertColumnDef(cd, relName, schema)
		if err != nil {
			return err
		}
		if err := c.atAddColumn(rel, colDef); err != nil {
			return err
		}
		// Analyze column default after column is added to the relation.
		// Coerce to column type to match PG's cookDefault (COERCE_IMPLICIT_CAST format).
		if colDef.RawDefault != nil {
			col := rel.Columns[len(rel.Columns)-1]
			if col.HasDefault {
				if analyzed, err := c.AnalyzeStandaloneExpr(colDef.RawDefault, rel); err == nil && analyzed != nil {
					if coerced, cerr := c.coerceToTargetType(analyzed, analyzed.exprType(), col.TypeOID, 'i'); cerr == nil && coerced != nil {
						analyzed = coerced
					}
					col.DefaultAnalyzed = analyzed
				}
			}
		}
		for _, con := range inlineCons {
			if err := c.addConstraint(schema, rel, con); err != nil {
				return err
			}
		}
		return nil

	case nodes.AT_ColumnDefault:
		if atc.Def != nil {
			defStr := deparseExprNode(atc.Def)
			rawExpr := atc.Def
			if con, ok := atc.Def.(*nodes.Constraint); ok {
				if con.CookedExpr != "" {
					defStr = con.CookedExpr
				}
				if con.RawExpr != nil {
					rawExpr = con.RawExpr
				}
			}
			if err := c.atSetDefault(rel, atc.Name, defStr); err != nil {
				return err
			}
			// Analyze the default expression and coerce to column type.
			// pg: cookDefault uses COERCE_IMPLICIT_CAST ('i') as display format
			if rawExpr != nil {
				if analyzed, err := c.AnalyzeStandaloneExpr(rawExpr, rel); err == nil && analyzed != nil {
					if idx, exists := rel.colByName[atc.Name]; exists {
						if coerced, cerr := c.coerceToTargetType(analyzed, analyzed.exprType(), rel.Columns[idx].TypeOID, 'i'); cerr == nil && coerced != nil {
							analyzed = coerced
						}
						rel.Columns[idx].DefaultAnalyzed = analyzed
					}
				}
			}
			return nil
		}
		return c.atDropDefault(rel, atc.Name)

	case nodes.AT_DropNotNull:
		return c.atDropNotNull(rel, atc.Name)

	case nodes.AT_SetNotNull:
		return c.atSetNotNull(rel, atc.Name)

	case nodes.AT_DropColumn:
		return c.atDropColumn(schema, rel, atc.Name, cascade, atc.Missing_ok)

	case nodes.AT_AddConstraint:
		con, ok := atc.Def.(*nodes.Constraint)
		if !ok {
			return fmt.Errorf("AT_AddConstraint: expected Constraint")
		}
		cdef, valid := convertConstraintNode(con)
		if !valid {
			return nil // Unsupported constraint type; skip.
		}
		return c.addConstraint(schema, rel, cdef)

	case nodes.AT_DropConstraint:
		return c.atDropConstraint(schema, rel, atc.Name, cascade, atc.Missing_ok)

	case nodes.AT_AlterColumnType:
		cd, ok := atc.Def.(*nodes.ColumnDef)
		if !ok {
			return fmt.Errorf("AT_AlterColumnType: expected ColumnDef")
		}
		newType := convertTypeNameToInternal(cd.TypeName)
		return c.atAlterColumnType(rel, atc.Name, newType)

	case nodes.AT_AddIdentity:
		cd, ok := atc.Def.(*nodes.ColumnDef)
		if !ok {
			return fmt.Errorf("AT_AddIdentity: expected ColumnDef")
		}
		return c.atAddIdentity(schema, rel, relName, atc.Name, cd)

	case nodes.AT_DropIdentity:
		return c.atDropIdentity(rel, atc.Name, atc.Missing_ok)

	case nodes.AT_SetIdentity:
		return nil // SET IDENTITY: no-op for pgddl (just sequence options)

	case nodes.AT_SetStatistics:
		return nil // SET STATISTICS: no-op for pgddl (planner hint only)

	case nodes.AT_SetStorage:
		// pgparser doesn't pass the storage type in Def, so this is a no-op.
		return nil

	case nodes.AT_ChangeOwner:
		return nil // OWNER TO: no-op for pgddl (no ACL)

	case nodes.AT_SetRelOptions, nodes.AT_ResetRelOptions, nodes.AT_ReplaceRelOptions:
		return nil // SET/RESET storage parameters: no-op for pgddl

	case nodes.AT_EnableTrig:
		return c.atEnableDisableTrigger(rel, atc.Name, 'O')

	case nodes.AT_EnableAlwaysTrig:
		return c.atEnableDisableTrigger(rel, atc.Name, 'A')

	case nodes.AT_EnableReplicaTrig:
		return c.atEnableDisableTrigger(rel, atc.Name, 'R')

	case nodes.AT_DisableTrig:
		return c.atEnableDisableTrigger(rel, atc.Name, 'D')

	case nodes.AT_EnableTrigAll:
		return c.atEnableDisableAllTriggers(rel, 'O')

	case nodes.AT_DisableTrigAll:
		return c.atEnableDisableAllTriggers(rel, 'D')

	case nodes.AT_EnableTrigUser:
		return c.atEnableDisableUserTriggers(rel, 'O')

	case nodes.AT_DisableTrigUser:
		return c.atEnableDisableUserTriggers(rel, 'D')

	case nodes.AT_ValidateConstraint:
		return c.atValidateConstraint(rel, atc.Name)

	case nodes.AT_EnableRowSecurity:
		rel.RowSecurity = true
		return nil

	case nodes.AT_DisableRowSecurity:
		rel.RowSecurity = false
		return nil

	case nodes.AT_ForceRowSecurity:
		rel.ForceRowSecurity = true
		return nil

	case nodes.AT_NoForceRowSecurity:
		rel.ForceRowSecurity = false
		return nil

	case nodes.AT_SetOptions, nodes.AT_ResetOptions:
		return nil // column-level options: no-op for pgddl

	case nodes.AT_AlterColumnGenericOptions:
		return nil // ALTER COLUMN SET/RESET options: no-op for pgddl

	case nodes.AT_AddIndex:
		return nil // Used internally by pg_dump restores: no-op

	case nodes.AT_AddIndexConstraint:
		// Used by pg_dump to add PK/UNIQUE constraints using an existing index.
		// pg: src/backend/commands/tablecmds.c — ATExecAddIndexConstraint
		idxStmt, ok := atc.Def.(*nodes.IndexStmt)
		if !ok {
			return fmt.Errorf("AT_AddIndexConstraint: expected IndexStmt")
		}
		ctype := ConstraintUnique
		if idxStmt.Primary {
			ctype = ConstraintPK
		}
		var colNames []string
		if idxStmt.IndexParams != nil {
			for _, item := range idxStmt.IndexParams.Items {
				if elem, ok := item.(*nodes.IndexElem); ok && elem.Name != "" {
					colNames = append(colNames, elem.Name)
				}
			}
		}
		def := ConstraintDef{
			Name:    idxStmt.Idxname,
			Type:    ctype,
			Columns: colNames,
		}
		return c.addConstraint(schema, rel, def)

	case nodes.AT_AddInherit:
		rv, ok := atc.Def.(*nodes.RangeVar)
		if !ok {
			return fmt.Errorf("AT_AddInherit: expected RangeVar")
		}
		return c.atAddInherit(rel, rv)

	case nodes.AT_DropInherit:
		rv, ok := atc.Def.(*nodes.RangeVar)
		if !ok {
			return fmt.Errorf("AT_DropInherit: expected RangeVar")
		}
		return c.atDropInherit(rel, rv)

	case nodes.AT_AlterConstraint:
		// pg: src/backend/commands/tablecmds.c — ATExecAlterConstraint
		con, ok := atc.Def.(*nodes.Constraint)
		if !ok {
			return nil // no constraint info: no-op
		}
		return c.atAlterConstraint(rel, con)

	case nodes.AT_SetExpression:
		return c.atSetExpression(rel, atc.Name, atc.Def)

	case nodes.AT_DropExpression:
		return c.atDropExpression(rel, atc.Name, atc.Missing_ok)

	case nodes.AT_SetLogged:
		rel.Persistence = 'p'
		return nil

	case nodes.AT_SetUnLogged:
		rel.Persistence = 'u'
		return nil

	case nodes.AT_ClusterOn:
		return c.atClusterOn(rel, atc.Name)

	case nodes.AT_DropCluster:
		return c.atDropCluster(rel)

	case nodes.AT_ReplicaIdentity:
		return c.atReplicaIdentity(rel, atc)

	case nodes.AT_SetCompression:
		method := byte(0)
		if s, ok := atc.Def.(*nodes.String); ok && s.Str != "" {
			switch s.Str {
			case "pglz":
				method = 'p'
			case "lz4":
				method = 'l'
			}
		}
		return c.atSetCompression(rel, atc.Name, method)

	case nodes.AT_AddOf:
		return c.atAddOf(rel, atc)

	case nodes.AT_DropOf:
		return c.atDropOf(rel)

	case nodes.AT_EnableRule, nodes.AT_DisableRule, nodes.AT_EnableAlwaysRule, nodes.AT_EnableReplicaRule:
		return nil // Rules: no-op for pgddl

	case nodes.AT_AttachPartition:
		pc, ok := atc.Def.(*nodes.PartitionCmd)
		if !ok {
			return fmt.Errorf("AT_AttachPartition: expected PartitionCmd")
		}
		return c.atExecAttachPartition(schema, rel, pc)

	case nodes.AT_DetachPartition:
		pc, ok := atc.Def.(*nodes.PartitionCmd)
		if !ok {
			return fmt.Errorf("AT_DetachPartition: expected PartitionCmd")
		}
		return c.atExecDetachPartition(rel, pc)

	default:
		return &Error{Code: CodeFeatureNotSupported, Message: fmt.Sprintf("unsupported ALTER TABLE subcommand %d", atc.Subtype)}
	}
}

// ExecRenameStmt handles RENAME operations from a pgparser AST.
//
// pg: src/backend/commands/alter.c — ExecRenameStmt
func (c *Catalog) ExecRenameStmt(stmt *nodes.RenameStmt) error {
	// Non-relation object types are handled first (don't need Relation lookup).
	switch stmt.RenameType {
	case nodes.OBJECT_SCHEMA:
		return c.renameSchema(stmt)
	case nodes.OBJECT_TYPE, nodes.OBJECT_DOMAIN:
		return c.renameType(stmt)
	case nodes.OBJECT_INDEX:
		return c.renameIndex(stmt)
	case nodes.OBJECT_SEQUENCE:
		return c.renameSequence(stmt)
	case nodes.OBJECT_TRIGGER:
		return c.renameTrigger(stmt)
	case nodes.OBJECT_TABCONSTRAINT:
		return c.renameConstraint(stmt)
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE:
		return c.renameFunction(stmt)
	case nodes.OBJECT_POLICY:
		return c.renamePolicy(stmt)
	case nodes.OBJECT_DOMCONSTRAINT:
		return nil // Domain constraints: no-op for pgddl
	case nodes.OBJECT_ATTRIBUTE:
		return c.renameAttribute(stmt)
	}

	// Relation-based renames (TABLE, VIEW, MATVIEW, COLUMN).
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	schema, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		return err
	}

	switch stmt.RenameType {
	case nodes.OBJECT_TABLE:
		if rel.RelKind != 'r' && rel.RelKind != 'p' {
			return errWrongObjectType(relName, "a table")
		}
		return c.renameRelation(schema, rel, stmt.Newname)

	case nodes.OBJECT_VIEW:
		if rel.RelKind != 'v' {
			return errWrongObjectType(relName, "a view")
		}
		return c.renameRelation(schema, rel, stmt.Newname)

	case nodes.OBJECT_MATVIEW:
		if rel.RelKind != 'm' {
			return errWrongObjectType(relName, "a materialized view")
		}
		return c.renameRelation(schema, rel, stmt.Newname)

	case nodes.OBJECT_FOREIGN_TABLE:
		if rel.RelKind != 'f' {
			return errWrongObjectType(relName, "a foreign table")
		}
		return c.renameRelation(schema, rel, stmt.Newname)

	case nodes.OBJECT_COLUMN:
		return c.atRenameColumn(rel, stmt.Subname, stmt.Newname)

	default:
		return &Error{Code: CodeFeatureNotSupported, Message: fmt.Sprintf("unsupported RENAME type %d", stmt.RenameType)}
	}
}

// renameRelation renames a relation (table or view) and updates associated types.
//
// (pgddl helper — generalizes PG's RenameRelationInternal)
func (c *Catalog) renameRelation(schema *Schema, rel *Relation, newName string) error {
	if _, exists := schema.Relations[newName]; exists {
		return errDuplicateTable(newName)
	}

	oldName := rel.Name
	delete(schema.Relations, oldName)
	rel.Name = newName
	schema.Relations[newName] = rel

	// Update row type name.
	if rt := c.typeByOID[rel.RowTypeOID]; rt != nil {
		delete(c.typeByName, typeKey{ns: rt.Namespace, name: rt.TypeName})
		rt.TypeName = newName
		c.typeByName[typeKey{ns: rt.Namespace, name: newName}] = rt
	}

	// Update array type name.
	if at := c.typeByOID[rel.ArrayOID]; at != nil {
		delete(c.typeByName, typeKey{ns: at.Namespace, name: at.TypeName})
		at.TypeName = "_" + newName
		c.typeByName[typeKey{ns: at.Namespace, name: at.TypeName}] = at
	}

	return nil
}

// renameSchema renames a schema.
//
// pg: src/backend/commands/alter.c — ExecRenameStmt (OBJECT_SCHEMA case)
func (c *Catalog) renameSchema(stmt *nodes.RenameStmt) error {
	oldName := ""
	if stmt.Object != nil {
		oldName = stringVal(stmt.Object)
	}
	// Fallback: some parser versions put schema name in Subname.
	if oldName == "" {
		oldName = stmt.Subname
	}
	newName := stmt.Newname

	s := c.schemaByName[oldName]
	if s == nil {
		return errUndefinedSchema(oldName)
	}

	// Prevent renaming built-in schemas.
	if s.OID == PGCatalogNamespace || s.OID == PGToastNamespace {
		return errInvalidParameterValue("cannot rename schema " + oldName)
	}

	// Reject reserved schema name prefix for the new name.
	// pg: src/backend/commands/schemacmds.c — RenameSchema (pg_ prefix check)
	if strings.HasPrefix(newName, "pg_") {
		return &Error{
			Code:    CodeReservedName,
			Message: "unacceptable schema name \"" + newName + "\"\nDetail: The prefix \"pg_\" is reserved for system schemas.",
		}
	}

	if _, exists := c.schemaByName[newName]; exists {
		return errDuplicateSchema(newName)
	}

	delete(c.schemaByName, oldName)
	s.Name = newName
	c.schemaByName[newName] = s
	return nil
}

// renameType renames a user-defined type (enum or domain).
//
// pg: src/backend/commands/alter.c — ExecRenameStmt (OBJECT_TYPE/OBJECT_DOMAIN case)
func (c *Catalog) renameType(stmt *nodes.RenameStmt) error {
	var schemaName, typeName string
	if stmt.Object != nil {
		// Object is a qualified name list.
		if list, ok := stmt.Object.(*nodes.List); ok {
			schemaName, typeName = qualifiedName(list)
		} else {
			typeName = stringVal(stmt.Object)
		}
	}
	newName := stmt.Newname

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	bt := c.typeByName[typeKey{ns: schema.OID, name: typeName}]
	if bt == nil {
		return errUndefinedType(typeName)
	}

	// Validate type kind.
	if stmt.RenameType == nodes.OBJECT_DOMAIN && bt.Type != 'd' {
		return errWrongObjectType(typeName, "a domain")
	}
	if stmt.RenameType == nodes.OBJECT_TYPE && bt.Type != 'e' && bt.Type != 'd' && bt.Type != 'r' && bt.Type != 'c' {
		return errWrongObjectType(typeName, "a type")
	}

	// Check for duplicate name.
	if c.typeByName[typeKey{ns: schema.OID, name: newName}] != nil {
		return errDuplicateObject("type", newName)
	}

	// For composite types, rename the backing relation too.
	if bt.Type == 'c' && bt.RelID != 0 {
		if rel := c.relationByOID[bt.RelID]; rel != nil && rel.RelKind == 'c' {
			return c.renameRelation(rel.Schema, rel, newName)
		}
	}

	// Update type name.
	delete(c.typeByName, typeKey{ns: schema.OID, name: typeName})
	bt.TypeName = newName
	c.typeByName[typeKey{ns: schema.OID, name: newName}] = bt

	// Update array type name.
	if bt.Array != 0 {
		if at := c.typeByOID[bt.Array]; at != nil {
			delete(c.typeByName, typeKey{ns: at.Namespace, name: at.TypeName})
			at.TypeName = "_" + newName
			c.typeByName[typeKey{ns: at.Namespace, name: at.TypeName}] = at
		}
	}

	return nil
}

// renameIndex renames an index.
//
// pg: src/backend/commands/alter.c — ExecRenameStmt (OBJECT_INDEX case)
func (c *Catalog) renameIndex(stmt *nodes.RenameStmt) error {
	var schemaName, idxName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		idxName = stmt.Relation.Relname
	}
	newName := stmt.Newname

	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	idx := schema.Indexes[idxName]
	if idx == nil {
		return errUndefinedObject("index", idxName)
	}

	// Check for duplicate name (index or relation).
	if _, exists := schema.Indexes[newName]; exists {
		return errDuplicateObject("index", newName)
	}
	if _, exists := schema.Relations[newName]; exists {
		return errDuplicateTable(newName)
	}

	delete(schema.Indexes, idxName)
	idx.Name = newName
	schema.Indexes[newName] = idx
	return nil
}

// renameSequence renames a sequence.
//
// pg: src/backend/commands/alter.c — ExecRenameStmt (OBJECT_SEQUENCE case)
func (c *Catalog) renameSequence(stmt *nodes.RenameStmt) error {
	var schemaName, seqName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		seqName = stmt.Relation.Relname
	}
	newName := stmt.Newname

	seq, err := c.findSequence(schemaName, seqName)
	if err != nil {
		return err
	}

	// Check for duplicate name.
	if _, exists := seq.Schema.Sequences[newName]; exists {
		return errDuplicateObject("sequence", newName)
	}

	delete(seq.Schema.Sequences, seqName)
	seq.Name = newName
	seq.Schema.Sequences[newName] = seq
	return nil
}

// renameTrigger renames a trigger on a relation.
//
// pg: src/backend/commands/alter.c — ExecRenameStmt (OBJECT_TRIGGER case)
func (c *Catalog) renameTrigger(stmt *nodes.RenameStmt) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	_, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		return err
	}

	oldName := stmt.Subname
	newName := stmt.Newname

	// Check for duplicate trigger name on this relation.
	for _, trig := range c.triggersByRel[rel.OID] {
		if trig.Name == newName {
			return errDuplicateObject("trigger", newName)
		}
	}

	// Find and rename.
	for _, trig := range c.triggersByRel[rel.OID] {
		if trig.Name == oldName {
			trig.Name = newName
			return nil
		}
	}
	return errUndefinedTrigger(oldName, relName)
}

// renameConstraint renames a constraint on a relation.
//
// pg: src/backend/commands/tablecmds.c — RenameConstraint
func (c *Catalog) renameConstraint(stmt *nodes.RenameStmt) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	_, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		return err
	}

	oldName := stmt.Subname
	newName := stmt.Newname

	// Check for duplicate constraint name on this relation.
	for _, con := range c.consByRel[rel.OID] {
		if con.Name == newName {
			return errDuplicateObject("constraint", newName)
		}
	}

	// Find and rename.
	for _, con := range c.consByRel[rel.OID] {
		if con.Name == oldName {
			con.Name = newName
			return nil
		}
	}
	return &Error{
		Code:    CodeUndefinedObject,
		Message: fmt.Sprintf("constraint %q of relation %q does not exist", oldName, relName),
	}
}

// renameFunction renames a user-defined function.
//
// pg: src/backend/commands/alter.c — ExecRenameStmt (OBJECT_FUNCTION case)
func (c *Catalog) renameFunction(stmt *nodes.RenameStmt) error {
	schemaName, funcName := extractObjectName(stmt.Object)
	schema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
	}

	procs := c.findUserProcsByName(schema, funcName)
	if len(procs) == 0 {
		return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("function %s does not exist", funcName)}
	}

	newName := stmt.Newname
	for _, up := range procs {
		bp := c.procByOID[up.OID]
		if bp != nil {
			// Remove from old name list.
			list := c.procByName[funcName]
			for i, p := range list {
				if p.OID == bp.OID {
					c.procByName[funcName] = append(list[:i], list[i+1:]...)
					break
				}
			}
			if len(c.procByName[funcName]) == 0 {
				delete(c.procByName, funcName)
			}
			bp.Name = newName
			c.procByName[newName] = append(c.procByName[newName], bp)
		}
		up.Name = newName
	}
	return nil
}

// renamePolicy renames a policy on a relation.
//
// pg: src/backend/commands/policy.c — rename_policy
func (c *Catalog) renamePolicy(stmt *nodes.RenameStmt) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	_, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		return err
	}

	oldName := stmt.Subname
	newName := stmt.Newname

	// Check for duplicate.
	for _, p := range c.policiesByRel[rel.OID] {
		if p.Name == newName {
			return errDuplicateObject("policy", newName)
		}
	}

	for _, p := range c.policiesByRel[rel.OID] {
		if p.Name == oldName {
			p.Name = newName
			return nil
		}
	}
	return errUndefinedObject("policy", oldName)
}

// renameAttribute renames a composite type attribute.
//
// pg: src/backend/commands/tablecmds.c — renameatt
func (c *Catalog) renameAttribute(stmt *nodes.RenameStmt) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	_, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		return err
	}

	// For OBJECT_ATTRIBUTE, the relation should be a composite type.
	if rel.RelKind != 'c' {
		return errWrongObjectType(relName, "a composite type")
	}

	return c.atRenameColumn(rel, stmt.Subname, stmt.Newname)
}

func (c *Catalog) atAddColumn(rel *Relation, colDef ColumnDef) error {
	// pg: src/backend/commands/tablecmds.c — ATExecAddColumn (MaxHeapAttributeNumber check)
	if len(rel.Columns) >= MaxHeapAttributeNumber {
		return &Error{Code: CodeTooManyColumns, Message: fmt.Sprintf("tables can have at most %d columns", MaxHeapAttributeNumber)}
	}

	if _, exists := rel.colByName[colDef.Name]; exists {
		return errDuplicateColumn(colDef.Name)
	}

	typeOID, typmod, err := c.ResolveType(colDef.Type)
	if err != nil {
		return err
	}
	typ := c.typeByOID[typeOID]

	col := &Column{
		AttNum:    int16(len(rel.Columns) + 1),
		Name:      colDef.Name,
		TypeOID:   typeOID,
		TypeMod:   typmod,
		NotNull:   colDef.NotNull,
		Len:       typ.Len,
		ByVal:     typ.ByVal,
		Align:     typ.Align,
		Storage:   typ.Storage,
		Collation: typ.Collation,
	}
	if colDef.Default != "" {
		col.HasDefault = true
		col.Default = colDef.Default
	}

	rel.Columns = append(rel.Columns, col)
	rel.colByName[col.Name] = len(rel.Columns) - 1
	return nil
}

// atDropColumn drops a column from a relation.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropColumn
func (c *Catalog) atDropColumn(schema *Schema, rel *Relation, colName string, cascade, ifExists bool) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		if ifExists {
			return nil
		}
		return errUndefinedColumn(colName)
	}

	// Block dropping inherited columns that are not locally defined.
	col := rel.Columns[idx]
	if col.InhCount > 0 && !col.IsLocal {
		return &Error{
			Code:    CodeInvalidObjectDefinition,
			Message: fmt.Sprintf("cannot drop inherited column %q", colName),
		}
	}

	// Cannot drop a column named in partition key.
	// pg: src/backend/commands/tablecmds.c — ATExecDropColumn (partition key check)
	if rel.PartitionInfo != nil {
		for _, keyAttNum := range rel.PartitionInfo.KeyAttNums {
			if keyAttNum == col.AttNum {
				return &Error{
					Code:    CodeInvalidObjectDefinition,
					Message: fmt.Sprintf("cannot drop column named in partition key"),
				}
			}
		}
	}

	// Check for dependent views.
	if deps := c.findNormalDependents('r', rel.OID); len(deps) > 0 {
		if !cascade {
			return errDependentObjects("column", colName)
		}
		c.dropDependents('r', rel.OID)
	}

	// Remove the column.
	rel.Columns = append(rel.Columns[:idx], rel.Columns[idx+1:]...)
	c.renumberColumns(rel)
	return nil
}

func (c *Catalog) atDropConstraint(schema *Schema, rel *Relation, conName string, cascade, ifExists bool) error {
	for _, con := range c.consByRel[rel.OID] {
		if con.Name == conName {
			c.removeConstraint(schema, con)
			return nil
		}
	}
	if ifExists {
		return nil
	}
	return &Error{
		Code:    CodeUndefinedObject,
		Message: fmt.Sprintf("constraint %q of relation %q does not exist", conName, rel.Name),
	}
}

func (c *Catalog) atRenameColumn(rel *Relation, oldName, newName string) error {
	idx, exists := rel.colByName[oldName]
	if !exists {
		return errUndefinedColumn(oldName)
	}
	if _, dup := rel.colByName[newName]; dup {
		return errDuplicateColumn(newName)
	}

	rel.Columns[idx].Name = newName
	delete(rel.colByName, oldName)
	rel.colByName[newName] = idx
	return nil
}


func (c *Catalog) atSetNotNull(rel *Relation, colName string) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	rel.Columns[idx].NotNull = true
	return nil
}

func (c *Catalog) atDropNotNull(rel *Relation, colName string) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}

	// Cannot drop NOT NULL on identity columns.
	// pg: src/backend/commands/tablecmds.c — ATExecDropNotNull (identity check)
	col := rel.Columns[idx]
	if col.Identity != 0 {
		return &Error{
			Code:    CodeInvalidObjectDefinition,
			Message: fmt.Sprintf("column %q of relation %q is an identity column", colName, rel.Name),
		}
	}

	// Cannot drop NOT NULL on PK columns.
	pk := c.findPKConstraint(rel.OID)
	if pk != nil {
		for _, attnum := range pk.Columns {
			if attnum == col.AttNum {
				return &Error{
					Code:    CodeDependentObjects,
					Message: fmt.Sprintf("column %q is in a primary key", colName),
				}
			}
		}
	}

	col.NotNull = false
	return nil
}

// atSetDefault sets a column's default expression.
//
// pg: src/backend/commands/tablecmds.c — ATExecColumnDefault
func (c *Catalog) atSetDefault(rel *Relation, colName, expr string) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]

	// Cannot set default on identity columns (use SET IDENTITY instead).
	if col.Identity != 0 {
		return errInvalidObjectDefinition(fmt.Sprintf(
			"column %q of relation %q is an identity column", colName, rel.Name))
	}

	// Cannot set default on generated columns.
	if col.Generated != 0 {
		return errInvalidObjectDefinition(fmt.Sprintf(
			"column %q of relation %q is a generated column", colName, rel.Name))
	}

	col.HasDefault = true
	col.Default = expr
	return nil
}

// atDropDefault removes a column's default expression.
//
// pg: src/backend/commands/tablecmds.c — ATExecColumnDefault
func (c *Catalog) atDropDefault(rel *Relation, colName string) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]

	// Cannot drop default on identity columns (use DROP IDENTITY instead).
	if col.Identity != 0 {
		return errInvalidObjectDefinition(fmt.Sprintf(
			"column %q of relation %q is an identity column", colName, rel.Name))
	}

	// Cannot drop default on generated columns (use DROP EXPRESSION instead).
	if col.Generated != 0 {
		return errInvalidObjectDefinition(fmt.Sprintf(
			"column %q of relation %q is a generated column", colName, rel.Name))
	}

	col.HasDefault = false
	col.Default = ""
	return nil
}

// atAlterColumnType changes the type of a column.
//
// pg: src/backend/commands/tablecmds.c — ATExecAlterColumnType
func (c *Catalog) atAlterColumnType(rel *Relation, colName string, newType TypeName) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}

	col := rel.Columns[idx]

	// Block if column is an identity column.
	if col.Identity != 0 {
		return errInvalidObjectDefinition(fmt.Sprintf(
			"cannot alter type of a column used as identity"))
	}

	// Block if any view depends on this relation.
	if deps := c.findNormalDependents('r', rel.OID); len(deps) > 0 {
		return &Error{
			Code:    CodeFeatureNotSupported,
			Message: "cannot alter type of a column used by a view or rule",
		}
	}

	newOID, newMod, err := c.ResolveType(newType)
	if err != nil {
		return err
	}

	if col.TypeOID != newOID && !c.CanCoerce(col.TypeOID, newOID, 'a') {
		oldType := c.typeByOID[col.TypeOID]
		newTyp := c.typeByOID[newOID]
		oldName, newName := "unknown", "unknown"
		if oldType != nil {
			oldName = oldType.TypeName
		}
		if newTyp != nil {
			newName = newTyp.TypeName
		}
		return errDatatypeMismatch(fmt.Sprintf(
			"column %q cannot be cast automatically to type %s from %s",
			colName, newName, oldName,
		))
	}

	col.TypeOID = newOID
	col.TypeMod = newMod

	typ := c.typeByOID[newOID]
	if typ != nil {
		col.Len = typ.Len
		col.ByVal = typ.ByVal
		col.Align = typ.Align
		col.Storage = typ.Storage
	}

	return nil
}

// atAddIdentity converts a column to an identity column.
//
// pg: src/backend/commands/tablecmds.c — ATExecAddIdentity
func (c *Catalog) atAddIdentity(schema *Schema, rel *Relation, relName, colName string, cd *nodes.ColumnDef) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]

	if col.Identity != 0 {
		return errInvalidObjectDefinition(fmt.Sprintf("column %q of relation %q is already an identity column", colName, rel.Name))
	}

	identity := cd.Identity
	if identity == 0 {
		identity = 'a' // default to ALWAYS
	}
	col.Identity = identity

	// Create identity sequence.
	seqName := fmt.Sprintf("%s_%s_seq", relName, colName)
	if _, seqExists := schema.Sequences[seqName]; !seqExists {
		seq := c.createSequenceInternal(schema, seqName, col.TypeOID)
		seq.OwnerRelOID = rel.OID
		seq.OwnerAttNum = col.AttNum
		c.recordDependency('s', seq.OID, 0, 'r', rel.OID, int32(col.AttNum), DepAuto)

		col.HasDefault = true
		col.Default = fmt.Sprintf("nextval('%s.%s'::regclass)", schema.Name, seqName)
	}

	return nil
}

// atDropIdentity removes identity from a column.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropIdentity
func (c *Catalog) atDropIdentity(rel *Relation, colName string, ifExists bool) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]

	if col.Identity == 0 {
		if ifExists {
			return nil
		}
		return errInvalidObjectDefinition(fmt.Sprintf("column %q of relation %q is not an identity column", colName, rel.Name))
	}

	col.Identity = 0
	// Note: In PG, dropping identity also drops the sequence.
	// For pgddl, we keep the default but clear the identity flag.
	return nil
}


// atAlterConstraint modifies a constraint's deferrable/deferred flags.
//
// pg: src/backend/commands/tablecmds.c — ATExecAlterConstraint
func (c *Catalog) atAlterConstraint(rel *Relation, con *nodes.Constraint) error {
	conName := con.Conname
	for _, existing := range c.consByRel[rel.OID] {
		if existing.Name == conName {
			existing.Deferrable = con.Deferrable
			existing.Deferred = con.Initdeferred
			return nil
		}
	}
	return &Error{
		Code:    CodeUndefinedObject,
		Message: fmt.Sprintf("constraint %q of relation %q does not exist", conName, rel.Name),
	}
}

// atEnableDisableTrigger sets the firing mode of a named trigger.
//
// pg: src/backend/commands/tablecmds.c — ATExecEnableDisableTrigger
func (c *Catalog) atEnableDisableTrigger(rel *Relation, trigName string, mode byte) error {
	for _, trig := range c.triggersByRel[rel.OID] {
		if trig.Name == trigName {
			trig.Enabled = mode
			return nil
		}
	}
	return errUndefinedTrigger(trigName, rel.Name)
}

// atEnableDisableAllTriggers sets the firing mode of all triggers on a relation.
//
// pg: src/backend/commands/tablecmds.c — ATExecEnableDisableTrigger (ALL)
func (c *Catalog) atEnableDisableAllTriggers(rel *Relation, mode byte) error {
	for _, trig := range c.triggersByRel[rel.OID] {
		trig.Enabled = mode
	}
	return nil
}

// atEnableDisableUserTriggers sets the firing mode of user-defined triggers on a relation.
//
// pg: src/backend/commands/tablecmds.c — ATExecEnableDisableTrigger (USER)
func (c *Catalog) atEnableDisableUserTriggers(rel *Relation, mode byte) error {
	for _, trig := range c.triggersByRel[rel.OID] {
		trig.Enabled = mode
	}
	return nil
}

// atValidateConstraint marks a constraint as validated.
//
// pg: src/backend/commands/tablecmds.c — ATExecValidateConstraint
func (c *Catalog) atValidateConstraint(rel *Relation, conName string) error {
	for _, con := range c.consByRel[rel.OID] {
		if con.Name == conName {
			con.Validated = true
			return nil
		}
	}
	return errUndefinedObject("constraint", conName)
}

// atExecAttachPartition attaches a partition to a partitioned table.
//
// pg: src/backend/commands/tablecmds.c — ATExecAttachPartition
func (c *Catalog) atExecAttachPartition(schema *Schema, parent *Relation, pc *nodes.PartitionCmd) error {
	if parent.PartitionInfo == nil {
		return errInvalidObjectDefinition(fmt.Sprintf("table %q is not partitioned", parent.Name))
	}
	if pc.Name == nil {
		return errInvalidParameterValue("ATTACH PARTITION requires a partition name")
	}

	_, part, err := c.findRelation(pc.Name.Schemaname, pc.Name.Relname)
	if err != nil {
		return err
	}

	// Store partition bound.
	if pc.Bound != nil {
		part.PartitionBound = &PartitionBound{
			Strategy:   pc.Bound.Strategy,
			IsDefault:  pc.Bound.IsDefault,
			ListValues: deparseDatumList(pc.Bound.Listdatums),
			LowerBound: deparseDatumList(pc.Bound.Lowerdatums),
			UpperBound: deparseDatumList(pc.Bound.Upperdatums),
			Modulus:    pc.Bound.Modulus,
			Remainder:  pc.Bound.Remainder,
		}
	}
	part.PartitionOf = parent.OID
	part.InhParents = []uint32{parent.OID}
	part.InhCount = 1

	// Record inheritance entry.
	c.inhEntries = append(c.inhEntries, InhEntry{
		InhRelID:  part.OID,
		InhParent: parent.OID,
		InhSeqNo:  1,
	})

	// Record auto-dependency (partition auto-dropped with parent).
	c.recordDependency('r', part.OID, 0, 'r', parent.OID, 0, DepAuto)

	// Clone parent's indexes, triggers, and FK constraints to the attached partition.
	// pg: src/backend/commands/tablecmds.c — ATExecAttachPartition
	c.cloneIndexesToPartition(schema, parent, part)
	c.cloneRowTriggersToPartition(parent, part)
	c.cloneForeignKeyConstraints(schema, parent, part)

	return nil
}

// atExecDetachPartition detaches a partition from a partitioned table.
//
// pg: src/backend/commands/tablecmds.c — ATExecDetachPartition
func (c *Catalog) atExecDetachPartition(parent *Relation, pc *nodes.PartitionCmd) error {
	if parent.PartitionInfo == nil {
		return errInvalidObjectDefinition(fmt.Sprintf("table %q is not partitioned", parent.Name))
	}
	if pc.Name == nil {
		return errInvalidParameterValue("DETACH PARTITION requires a partition name")
	}

	_, part, err := c.findRelation(pc.Name.Schemaname, pc.Name.Relname)
	if err != nil {
		return err
	}

	if part.PartitionOf != parent.OID {
		return errInvalidObjectDefinition(fmt.Sprintf("relation %q is not a partition of %q", part.Name, parent.Name))
	}

	// Clear partition info.
	part.PartitionBound = nil
	part.PartitionOf = 0
	part.InhParents = nil
	part.InhCount = 0

	// Remove inheritance entry.
	c.removeInhEntries(part.OID)

	// Remove dependency.
	c.removeDepsOf('r', part.OID)

	return nil
}

// renumberColumns rebuilds attnums and colByName after a column removal.
func (c *Catalog) renumberColumns(rel *Relation) {
	rel.colByName = make(map[string]int, len(rel.Columns))
	for i, col := range rel.Columns {
		col.AttNum = int16(i + 1)
		rel.colByName[col.Name] = i
	}
}

// ExecAlterObjectSchemaStmt moves an object to a different schema (SET SCHEMA).
//
// pg: src/backend/commands/alter.c — ExecAlterObjectSchemaStmt
func (c *Catalog) ExecAlterObjectSchemaStmt(stmt *nodes.AlterObjectSchemaStmt) error {
	newSchema := c.schemaByName[stmt.Newschema]
	if newSchema == nil {
		return errUndefinedSchema(stmt.Newschema)
	}

	switch stmt.ObjectType {
	case nodes.OBJECT_TABLE, nodes.OBJECT_VIEW, nodes.OBJECT_MATVIEW, nodes.OBJECT_SEQUENCE, nodes.OBJECT_FOREIGN_TABLE:
		return c.alterRelationNamespace(stmt, newSchema)
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE:
		return c.alterFunctionNamespace(stmt, newSchema)
	case nodes.OBJECT_TYPE, nodes.OBJECT_DOMAIN:
		return c.alterTypeNamespace(stmt, newSchema)
	default:
		return nil // no-op for unsupported object types
	}
}

// alterRelationNamespace moves a relation (table/view/matview/sequence) to a new schema.
//
// pg: src/backend/commands/tablecmds.c — AlterTableNamespace
func (c *Catalog) alterRelationNamespace(stmt *nodes.AlterObjectSchemaStmt, newSchema *Schema) error {
	var schemaName, relName string
	if stmt.Relation != nil {
		schemaName = stmt.Relation.Schemaname
		relName = stmt.Relation.Relname
	}

	// Handle sequences.
	if stmt.ObjectType == nodes.OBJECT_SEQUENCE {
		seq, err := c.findSequence(schemaName, relName)
		if err != nil {
			if stmt.MissingOk {
				return nil
			}
			return err
		}
		if _, exists := newSchema.Sequences[seq.Name]; exists {
			return errDuplicateObject("sequence", seq.Name)
		}
		delete(seq.Schema.Sequences, seq.Name)
		newSchema.Sequences[seq.Name] = seq
		seq.Schema = newSchema
		return nil
	}

	oldSchema, rel, err := c.findRelation(schemaName, relName)
	if err != nil {
		if stmt.MissingOk {
			return nil
		}
		return err
	}

	// Check for name conflict in target schema.
	if _, exists := newSchema.Relations[rel.Name]; exists {
		return errDuplicateTable(rel.Name)
	}

	// Move relation.
	delete(oldSchema.Relations, rel.Name)
	newSchema.Relations[rel.Name] = rel
	rel.Schema = newSchema

	// Move row type and array type namespace.
	if rt := c.typeByOID[rel.RowTypeOID]; rt != nil {
		delete(c.typeByName, typeKey{ns: oldSchema.OID, name: rt.TypeName})
		rt.Namespace = newSchema.OID
		c.typeByName[typeKey{ns: newSchema.OID, name: rt.TypeName}] = rt
	}
	if at := c.typeByOID[rel.ArrayOID]; at != nil {
		delete(c.typeByName, typeKey{ns: oldSchema.OID, name: at.TypeName})
		at.Namespace = newSchema.OID
		c.typeByName[typeKey{ns: newSchema.OID, name: at.TypeName}] = at
	}

	// Move indexes.
	for _, idx := range c.indexesByRel[rel.OID] {
		if idx.Schema == oldSchema {
			delete(oldSchema.Indexes, idx.Name)
			newSchema.Indexes[idx.Name] = idx
			idx.Schema = newSchema
		}
	}

	return nil
}

// alterFunctionNamespace moves a function to a new schema.
//
// pg: src/backend/commands/functioncmds.c — AlterFunctionNamespace
func (c *Catalog) alterFunctionNamespace(stmt *nodes.AlterObjectSchemaStmt, newSchema *Schema) error {
	schemaName, funcName := extractObjectName(stmt.Object)
	oldSchema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		if stmt.MissingOk {
			return nil
		}
		return err
	}

	// Find all matching user procs in the old schema.
	procs := c.findUserProcsByName(oldSchema, funcName)
	if len(procs) == 0 {
		if stmt.MissingOk {
			return nil
		}
		return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("function %s does not exist", funcName)}
	}

	for _, up := range procs {
		up.Schema = newSchema
	}
	return nil
}

// alterTypeNamespace moves a type to a new schema.
//
// pg: src/backend/commands/typecmds.c — AlterTypeNamespace
func (c *Catalog) alterTypeNamespace(stmt *nodes.AlterObjectSchemaStmt, newSchema *Schema) error {
	schemaName, typeName := extractObjectName(stmt.Object)
	oldSchema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		if stmt.MissingOk {
			return nil
		}
		return err
	}

	bt := c.typeByName[typeKey{ns: oldSchema.OID, name: typeName}]
	if bt == nil {
		if stmt.MissingOk {
			return nil
		}
		return errUndefinedType(typeName)
	}

	// Check for name conflict.
	if c.typeByName[typeKey{ns: newSchema.OID, name: typeName}] != nil {
		return errDuplicateObject("type", typeName)
	}

	// Move type.
	delete(c.typeByName, typeKey{ns: oldSchema.OID, name: typeName})
	bt.Namespace = newSchema.OID
	c.typeByName[typeKey{ns: newSchema.OID, name: typeName}] = bt

	// Move array type.
	if bt.Array != 0 {
		if at := c.typeByOID[bt.Array]; at != nil {
			delete(c.typeByName, typeKey{ns: oldSchema.OID, name: at.TypeName})
			at.Namespace = newSchema.OID
			c.typeByName[typeKey{ns: newSchema.OID, name: at.TypeName}] = at
		}
	}

	// For composite types, move the backing relation too.
	if bt.Type == 'c' && bt.RelID != 0 {
		if rel := c.relationByOID[bt.RelID]; rel != nil && rel.Schema == oldSchema {
			delete(oldSchema.Relations, rel.Name)
			newSchema.Relations[rel.Name] = rel
			rel.Schema = newSchema
		}
	}

	return nil
}

// extractObjectName extracts schema and name from an Object node (typically a *nodes.List or *nodes.String).
//
// (pgddl helper)
func extractObjectName(obj nodes.Node) (schema, name string) {
	switch n := obj.(type) {
	case *nodes.List:
		return qualifiedName(n)
	case *nodes.String:
		return "", n.Str
	default:
		return "", ""
	}
}

// atAddInherit adds a parent to a table's inheritance list.
//
// pg: src/backend/commands/tablecmds.c — ATExecAddInherit
func (c *Catalog) atAddInherit(rel *Relation, rv *nodes.RangeVar) error {
	_, parent, err := c.findRelation(rv.Schemaname, rv.Relname)
	if err != nil {
		return err
	}
	if parent.RelKind != 'r' && parent.RelKind != 'p' {
		return errWrongObjectType(rv.Relname, "a table")
	}

	// Circularity check: walk parent's inheritance tree.
	if c.isAncestor(parent.OID, rel.OID) {
		return errInvalidObjectDefinition(fmt.Sprintf("circular inheritance not allowed"))
	}

	// Merge parent columns into child (type must match, increment InhCount).
	for _, pcol := range parent.Columns {
		if existIdx, exists := rel.colByName[pcol.Name]; exists {
			childCol := rel.Columns[existIdx]
			if childCol.TypeOID != pcol.TypeOID {
				return errDatatypeMismatch(fmt.Sprintf(
					"column %q has a type conflict", pcol.Name,
				))
			}
			if pcol.NotNull {
				childCol.NotNull = true
			}
			childCol.InhCount++
		} else {
			// Add inherited column.
			col := &Column{
				AttNum:    int16(len(rel.Columns) + 1),
				Name:      pcol.Name,
				TypeOID:   pcol.TypeOID,
				TypeMod:   pcol.TypeMod,
				NotNull:   pcol.NotNull,
				Len:       pcol.Len,
				ByVal:     pcol.ByVal,
				Align:     pcol.Align,
				Storage:   pcol.Storage,
				Collation: pcol.Collation,
				IsLocal:   false,
				InhCount:  1,
			}
			if pcol.HasDefault {
				col.HasDefault = true
				col.Default = pcol.Default
			}
			rel.Columns = append(rel.Columns, col)
			rel.colByName[col.Name] = len(rel.Columns) - 1
		}
	}

	// Record inheritance entry.
	seqNo := int32(len(rel.InhParents) + 1)
	rel.InhParents = append(rel.InhParents, parent.OID)
	rel.InhCount++
	c.inhEntries = append(c.inhEntries, InhEntry{
		InhRelID:  rel.OID,
		InhParent: parent.OID,
		InhSeqNo:  seqNo,
	})
	c.recordDependency('r', rel.OID, 0, 'r', parent.OID, 0, DepNormal)

	return nil
}

// atDropInherit removes a parent from a table's inheritance list.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropInherit
func (c *Catalog) atDropInherit(rel *Relation, rv *nodes.RangeVar) error {
	// Cannot use on partitions (must DETACH instead).
	if rel.PartitionOf != 0 {
		return errInvalidObjectDefinition("cannot drop inheritance from a partition; use DETACH PARTITION instead")
	}

	_, parent, err := c.findRelation(rv.Schemaname, rv.Relname)
	if err != nil {
		return err
	}

	// Verify this is actually a parent.
	found := false
	for i, pid := range rel.InhParents {
		if pid == parent.OID {
			rel.InhParents = append(rel.InhParents[:i], rel.InhParents[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errInvalidObjectDefinition(fmt.Sprintf("relation %q is not a parent of %q", parent.Name, rel.Name))
	}
	rel.InhCount--

	// Decrement InhCount on inherited columns; if InhCount→0, set IsLocal=true.
	for _, pcol := range parent.Columns {
		if existIdx, exists := rel.colByName[pcol.Name]; exists {
			childCol := rel.Columns[existIdx]
			childCol.InhCount--
			if childCol.InhCount <= 0 {
				childCol.IsLocal = true
				childCol.InhCount = 0
			}
		}
	}

	// Remove inheritance entry.
	n := 0
	for _, e := range c.inhEntries {
		if e.InhRelID == rel.OID && e.InhParent == parent.OID {
			continue
		}
		c.inhEntries[n] = e
		n++
	}
	c.inhEntries = c.inhEntries[:n]

	// Remove dependency.
	c.removeDep('r', rel.OID, 'r', parent.OID)

	return nil
}

// isAncestor checks if ancestorOID is an ancestor of relOID in the inheritance tree.
func (c *Catalog) isAncestor(ancestorOID, relOID uint32) bool {
	for _, e := range c.inhEntries {
		if e.InhRelID == ancestorOID && e.InhParent == relOID {
			return true
		}
		if e.InhRelID == ancestorOID {
			if c.isAncestor(e.InhParent, relOID) {
				return true
			}
		}
	}
	return false
}

// removeDep removes a specific dependency entry.
func (c *Catalog) removeDep(objType byte, objOID uint32, refType byte, refOID uint32) {
	n := 0
	for _, dep := range c.deps {
		if dep.ObjType == objType && dep.ObjOID == objOID && dep.RefType == refType && dep.RefOID == refOID {
			continue
		}
		c.deps[n] = dep
		n++
	}
	c.deps = c.deps[:n]
}

// atClusterOn sets the clustered index on a relation.
//
// pg: src/backend/commands/tablecmds.c — ATExecClusterOn
func (c *Catalog) atClusterOn(rel *Relation, indexName string) error {
	// Clear IsClustered on all indexes first.
	for _, idx := range c.indexesByRel[rel.OID] {
		idx.IsClustered = false
	}
	// Find and set the named index.
	for _, idx := range c.indexesByRel[rel.OID] {
		if idx.Name == indexName {
			idx.IsClustered = true
			return nil
		}
	}
	return errUndefinedObject("index", indexName)
}

// atDropCluster clears the clustered index on a relation.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropCluster
func (c *Catalog) atDropCluster(rel *Relation) error {
	for _, idx := range c.indexesByRel[rel.OID] {
		idx.IsClustered = false
	}
	return nil
}

// atReplicaIdentity sets the replica identity for a relation.
//
// pg: src/backend/commands/tablecmds.c — ATExecReplicaIdentity
func (c *Catalog) atReplicaIdentity(rel *Relation, atc *nodes.AlterTableCmd) error {
	// Clear IsReplicaIdent on all indexes.
	for _, idx := range c.indexesByRel[rel.OID] {
		idx.IsReplicaIdent = false
	}

	if atc.Name != "" {
		// USING INDEX variant.
		rel.ReplicaIdentity = 'i'
		for _, idx := range c.indexesByRel[rel.OID] {
			if idx.Name == atc.Name {
				if !idx.IsUnique {
					return errInvalidObjectDefinition(fmt.Sprintf("index %q is not a unique index", atc.Name))
				}
				idx.IsReplicaIdent = true
				return nil
			}
		}
		return errUndefinedObject("index", atc.Name)
	}

	// DEFAULT/FULL/NOTHING — pgparser doesn't distinguish, default to 'd'.
	rel.ReplicaIdentity = 'd'
	return nil
}

// atSetCompression sets the compression method for a column.
//
// pg: src/backend/commands/tablecmds.c — ATExecSetCompression
func (c *Catalog) atSetCompression(rel *Relation, colName string, method byte) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]
	// Validate column type is variable-length (Len == -1).
	if col.Len != -1 {
		return errInvalidObjectDefinition(fmt.Sprintf("column data type %s does not support compression", c.typeByOID[col.TypeOID].TypeName))
	}
	col.Compression = method
	return nil
}

// atAddOf sets the typed table association (OF type).
//
// pg: src/backend/commands/tablecmds.c — ATExecAddOf
func (c *Catalog) atAddOf(rel *Relation, atc *nodes.AlterTableCmd) error {
	tn, ok := atc.Def.(*nodes.TypeName)
	if !ok {
		return errInvalidParameterValue("OF requires a type name")
	}
	typName := convertTypeNameToInternal(tn)
	typeOID, _, err := c.ResolveType(typName)
	if err != nil {
		return err
	}

	bt := c.typeByOID[typeOID]
	if bt == nil || bt.Type != 'c' {
		return errWrongObjectType(typName.Name, "a composite type")
	}

	// Verify columns match 1:1 with the composite type's relation.
	if bt.RelID != 0 {
		typeRel := c.relationByOID[bt.RelID]
		if typeRel != nil {
			if len(rel.Columns) != len(typeRel.Columns) {
				return errDatatypeMismatch(fmt.Sprintf("table has %d columns but type has %d", len(rel.Columns), len(typeRel.Columns)))
			}
			for i, tcol := range typeRel.Columns {
				rcol := rel.Columns[i]
				if rcol.Name != tcol.Name || rcol.TypeOID != tcol.TypeOID {
					return errDatatypeMismatch(fmt.Sprintf("table column %q does not match type column %q", rcol.Name, tcol.Name))
				}
			}
		}
	}

	rel.OfTypeOID = typeOID
	c.recordDependency('r', rel.OID, 0, 't', typeOID, 0, DepNormal)
	return nil
}

// atDropOf removes the typed table association.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropOf
func (c *Catalog) atDropOf(rel *Relation) error {
	if rel.OfTypeOID != 0 {
		c.removeDep('r', rel.OID, 't', rel.OfTypeOID)
		rel.OfTypeOID = 0
	}
	return nil
}

// atSetExpression updates a generated column's expression.
//
// pg: src/backend/commands/tablecmds.c — ATExecSetExpression
func (c *Catalog) atSetExpression(rel *Relation, colName string, expr nodes.Node) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]
	if col.Generated != 's' {
		return errInvalidObjectDefinition(fmt.Sprintf("column %q of relation %q is not a generated column", colName, rel.Name))
	}
	col.GenerationExpr = deparseExprNode(expr)
	col.Default = col.GenerationExpr
	return nil
}

// atDropExpression drops the generated column expression.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropExpression
func (c *Catalog) atDropExpression(rel *Relation, colName string, ifExists bool) error {
	idx, exists := rel.colByName[colName]
	if !exists {
		return errUndefinedColumn(colName)
	}
	col := rel.Columns[idx]
	if col.Generated != 's' {
		if ifExists {
			return nil
		}
		return errInvalidObjectDefinition(fmt.Sprintf("column %q of relation %q is not a generated column", colName, rel.Name))
	}
	col.Generated = 0
	col.GenerationExpr = ""
	col.HasDefault = false
	col.Default = ""
	return nil
}

// findChildRelations returns OIDs of all direct children (inheritance or partition) of a relation.
//
// (pgddl helper — scans pg_inherits entries)
func (c *Catalog) findChildRelations(parentOID uint32) []uint32 {
	var children []uint32
	for _, e := range c.inhEntries {
		if e.InhParent == parentOID {
			children = append(children, e.InhRelID)
		}
	}
	return children
}

// shouldRecurse returns true if a given ALTER TABLE subcommand should be propagated
// to inheritance children.
//
// pg: src/backend/commands/tablecmds.c — ATSimpleRecursion
func shouldRecurse(subtype nodes.AlterTableType) bool {
	switch subtype {
	case nodes.AT_SetNotNull, nodes.AT_DropNotNull, nodes.AT_AddConstraint,
		nodes.AT_DropConstraint, nodes.AT_AddColumn, nodes.AT_DropColumn:
		return true
	}
	return false
}

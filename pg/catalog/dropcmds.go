package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// RemoveObjects dispatches DROP statements for non-relation object types.
// Handles OBJECT_INDEX, OBJECT_SEQUENCE, OBJECT_TYPE, OBJECT_FUNCTION,
// OBJECT_PROCEDURE, and OBJECT_TRIGGER.
//
// For OBJECT_TABLE and OBJECT_VIEW, use RemoveRelations instead.
// For OBJECT_SCHEMA, use RemoveSchemas instead.
//
// pg: src/backend/commands/dropcmds.c — RemoveObjects
func (c *Catalog) RemoveObjects(stmt *nodes.DropStmt) error {
	switch nodes.ObjectType(stmt.RemoveType) {
	case nodes.OBJECT_INDEX:
		return c.removeIndexObjects(stmt)
	case nodes.OBJECT_SEQUENCE:
		return c.removeSequenceObjects(stmt)
	case nodes.OBJECT_TYPE, nodes.OBJECT_DOMAIN:
		return c.removeTypeObjects(stmt)
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE:
		return c.removeFunctionObjects(stmt)
	case nodes.OBJECT_AGGREGATE:
		// pg: src/backend/commands/dropcmds.c — RemoveObjects (OBJECT_AGGREGATE)
		// For pgddl, aggregates are stored as functions with prokind='a'.
		// Treat DROP AGGREGATE as DROP FUNCTION since pgddl doesn't distinguish them.
		return c.removeFunctionObjects(stmt)
	case nodes.OBJECT_TRIGGER:
		return c.removeTriggerObjects(stmt)
	case nodes.OBJECT_EVENT_TRIGGER:
		return c.removeEventTriggerObjects(stmt)
	case nodes.OBJECT_POLICY:
		return c.removePolicyObjects(stmt)
	case nodes.OBJECT_TABLE, nodes.OBJECT_VIEW, nodes.OBJECT_MATVIEW:
		return c.RemoveRelations(stmt)
	case nodes.OBJECT_FOREIGN_TABLE:
		// DROP FOREIGN TABLE: route to RemoveRelations (real handling).
		// pg: src/backend/commands/tablecmds.c — RemoveRelations
		return c.RemoveRelations(stmt)
	case nodes.OBJECT_SCHEMA:
		return c.RemoveSchemas(stmt)

	// No-op object types: these objects are not tracked by pgddl.
	case nodes.OBJECT_COLLATION,
		nodes.OBJECT_CONVERSION,
		nodes.OBJECT_OPERATOR,
		nodes.OBJECT_OPCLASS,
		nodes.OBJECT_OPFAMILY,
		nodes.OBJECT_LANGUAGE,
		nodes.OBJECT_FDW,
		nodes.OBJECT_FOREIGN_SERVER,
		nodes.OBJECT_ACCESS_METHOD,
		nodes.OBJECT_PUBLICATION,
		nodes.OBJECT_RULE,
		nodes.OBJECT_CAST,
		nodes.OBJECT_TRANSFORM,
		nodes.OBJECT_STATISTIC_EXT,
		nodes.OBJECT_TSPARSER,
		nodes.OBJECT_TSDICTIONARY,
		nodes.OBJECT_TSTEMPLATE,
		nodes.OBJECT_TSCONFIGURATION,
		nodes.OBJECT_EXTENSION:
		return nil

	default:
		return errInvalidParameterValue("unsupported object type for DROP")
	}
}

// removeIndexObjects handles DROP INDEX for one or more indexes.
func (c *Catalog) removeIndexObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		schemaName, idxName := extractDropObjectName(obj)

		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("index %q does not exist, skipping", idxName))
					continue
				}
			}
			return err
		}

		idx, exists := schema.Indexes[idxName]
		if !exists {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("index %q does not exist, skipping", idxName))
				continue
			}
			return errUndefinedObject("index", idxName)
		}

		if idx.ConstraintOID != 0 {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop index %q because constraint requires it", idxName),
			}
		}

		c.removeIndex(schema, idxName, idx)
	}
	return nil
}

// removeSequenceObjects handles DROP SEQUENCE for one or more sequences.
func (c *Catalog) removeSequenceObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		schemaName, seqName := extractDropObjectName(obj)

		seq, err := c.findSequence(schemaName, seqName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("sequence %q does not exist, skipping", seqName))
					continue
				}
			}
			return err
		}

		// Check dependents.
		if deps := c.findNormalDependents('s', seq.OID); len(deps) > 0 {
			if !cascade {
				return errDependentObjects("sequence", seqName)
			}
			c.dropDependents('s', seq.OID)
		}

		c.removeSequence(seq.Schema, seq)
	}
	return nil
}

// removeTypeObjects handles DROP TYPE for one or more types (enum or domain).
func (c *Catalog) removeTypeObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		schemaName, typeName := extractDropObjectName(obj)

		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("type %q does not exist, skipping", typeName))
					continue
				}
			}
			return err
		}

		bt := c.typeByName[typeKey{ns: schema.OID, name: typeName}]
		if bt == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("type %q does not exist, skipping", typeName))
				continue
			}
			return errUndefinedType(typeName)
		}

		if bt.Type != 'e' && bt.Type != 'd' && bt.Type != 'r' && !(bt.Type == 'c' && bt.RelID != 0) {
			return errWrongObjectType(typeName, "a droppable type")
		}

		// Composite type: drop the backing relation (relkind='c') which also removes the type.
		if bt.Type == 'c' && bt.RelID != 0 {
			rel := c.relationByOID[bt.RelID]
			if rel != nil && rel.RelKind == 'c' {
				// Check for dependents.
				relDeps := c.findTypeDependents(bt.OID)
				recordedDeps := c.findNormalDependents('t', bt.OID)
				if len(relDeps) > 0 || len(recordedDeps) > 0 {
					if !cascade {
						return errDependentObjects("type", typeName)
					}
					for _, depRelOID := range relDeps {
						depRel := c.relationByOID[depRelOID]
						if depRel == nil {
							continue
						}
						c.dropDependents('r', depRel.OID)
						c.removeRelation(depRel.Schema, depRel.Name, depRel)
					}
					c.dropDependents('t', bt.OID)
				}
				c.removeRelation(rel.Schema, rel.Name, rel)
				continue
			}
		}

		// Check for dependents: relations using this type, and recorded deps.
		relDeps := c.findTypeDependents(bt.OID)
		recordedDeps := c.findNormalDependents('t', bt.OID)
		if len(relDeps) > 0 || len(recordedDeps) > 0 {
			if !cascade {
				return errDependentObjects("type", typeName)
			}
			for _, relOID := range relDeps {
				rel := c.relationByOID[relOID]
				if rel == nil {
					continue
				}
				c.dropDependents('r', rel.OID)
				c.removeRelation(rel.Schema, rel.Name, rel)
			}
			c.dropDependents('t', bt.OID)
		}

		// Remove enum/domain/range metadata.
		if bt.Type == 'e' {
			delete(c.enumTypes, bt.OID)
		} else if bt.Type == 'd' {
			delete(c.domainTypes, bt.OID)
		} else if bt.Type == 'r' {
			delete(c.rangeTypes, bt.OID)
		}

		// Remove array type.
		if bt.Array != 0 {
			if at := c.typeByOID[bt.Array]; at != nil {
				delete(c.typeByName, typeKey{ns: at.Namespace, name: at.TypeName})
			}
			delete(c.typeByOID, bt.Array)
		}

		// Remove the type itself.
		delete(c.typeByName, typeKey{ns: bt.Namespace, name: bt.TypeName})
		delete(c.typeByOID, bt.OID)
		c.removeComments('t', bt.OID)
		c.removeDepsOf('t', bt.OID)
		c.removeDepsOn('t', bt.OID)
	}
	return nil
}

// removeFunctionObjects handles DROP FUNCTION/PROCEDURE for one or more functions.
func (c *Catalog) removeFunctionObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	isProcedure := nodes.ObjectType(stmt.RemoveType) == nodes.OBJECT_PROCEDURE
	dropKindStr := "function"
	if isProcedure {
		dropKindStr = "procedure"
	}
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		owa, ok := obj.(*nodes.ObjectWithArgs)
		if !ok {
			continue
		}

		schemaName, funcName := qualifiedName(owa.Objname)
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("%s %s does not exist, skipping", dropKindStr, funcName))
					continue
				}
			}
			return err
		}

		var bp *BuiltinProc

		if !owa.ArgsUnspecified && owa.Objargs != nil {
			// Resolve arg types and find exact match.
			argOIDs := make([]uint32, 0, len(owa.Objargs.Items))
			for _, item := range owa.Objargs.Items {
				tn, ok := item.(*nodes.TypeName)
				if !ok {
					continue
				}
				oid, _, err := c.resolveTypeName(tn)
				if err != nil {
					return err
				}
				argOIDs = append(argOIDs, oid)
			}
			bp = c.findExactProc(schema, funcName, argOIDs)
		} else {
			// No arg types: find by name. If ambiguous, error.
			candidates := c.findUserProcsByName(schema, funcName)
			if len(candidates) > 1 {
				return errAmbiguousFunction(funcName)
			}
			if len(candidates) == 1 {
				bp = c.procByOID[candidates[0].OID]
			}
		}

		if bp == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("%s %s does not exist, skipping", dropKindStr, funcName))
				continue
			}
			// pg: src/backend/commands/dropcmds.c — uses ERRCODE_UNDEFINED_FUNCTION
			return &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", dropKindStr, funcName)}
		}

		// Verify it's user-defined.
		if _, ok := c.userProcs[bp.OID]; !ok {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("function %q is not user-defined", funcName)}
		}

		// Verify DROP FUNCTION doesn't drop a procedure, and vice versa.
		// pg: src/backend/commands/functioncmds.c — check_routine_type
		if isProcedure && bp.Kind != 'p' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a procedure", funcName)}
		}
		if !isProcedure && nodes.ObjectType(stmt.RemoveType) == nodes.OBJECT_FUNCTION && bp.Kind == 'p' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is not a function", funcName)}
		}

		// DROP FUNCTION on an aggregate should hint to use DROP AGGREGATE instead.
		// pg: src/backend/commands/dropcmds.c — RemoveObjects (line 91-99)
		if !isProcedure && nodes.ObjectType(stmt.RemoveType) == nodes.OBJECT_FUNCTION && bp.Kind == 'a' {
			return &Error{Code: CodeWrongObjectType,
				Message: fmt.Sprintf("%q is an aggregate function\nHint: Use DROP AGGREGATE to drop aggregate functions.", funcName)}
		}

		// Check dependents (triggers).
		if deps := c.findNormalDependents('f', bp.OID); len(deps) > 0 {
			if !cascade {
				kindStr := "function"
				if isProcedure {
					kindStr = "procedure"
				}
				return errDependentObjects(kindStr, funcName)
			}
			c.dropDependents('f', bp.OID)
		}

		c.removeFunction(bp.OID, funcName)
	}
	return nil
}

// removeTriggerObjects handles DROP TRIGGER for one or more triggers.
// Triggers in a DropStmt are somewhat special because they require the table name.
// The parser represents them as a list with the trigger name; the table is on the
// DropStmt itself for single-object drops, but for the general case we handle
// the Objects list which contains qualified names.
func (c *Catalog) removeTriggerObjects(stmt *nodes.DropStmt) error {
	// DROP TRIGGER is unusual: the parser does not emit OBJECT_TRIGGER in DropStmt.
	// Instead, triggers are dropped via a dedicated DropTriggerStmt-like representation.
	// However, pgparser uses DropStmt with RemoveType=OBJECT_TRIGGER.
	// The Objects list contains pairs: [table-name, trigger-name] or just [trigger-name].
	// We need to handle this carefully.
	//
	// In practice, pgparser represents DROP TRIGGER name ON table as:
	//   Objects: [[table, trigger]] or [[trigger]] depending on parser version.
	// We'll handle the common case.
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		// For triggers, the object list contains a sublist: [tablename, triggername].
		list, ok := obj.(*nodes.List)
		if !ok {
			continue
		}

		var schemaName, tableName, trigName string
		items := list.Items
		switch len(items) {
		case 2:
			// [table, trigger]
			tableName = stringVal(items[0])
			trigName = stringVal(items[1])
		case 3:
			// [schema, table, trigger]
			schemaName = stringVal(items[0])
			tableName = stringVal(items[1])
			trigName = stringVal(items[2])
		default:
			continue
		}

		_, rel, err := c.findRelation(schemaName, tableName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("trigger %q for relation %q does not exist, skipping", trigName, tableName))
					continue
				}
			}
			return err
		}

		var found *Trigger
		for _, trig := range c.triggersByRel[rel.OID] {
			if trig.Name == trigName {
				found = trig
				break
			}
		}

		if found == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("trigger %q for relation %q does not exist, skipping", trigName, tableName))
				continue
			}
			return errUndefinedTrigger(trigName, tableName)
		}

		if cascade {
			c.dropDependents('g', found.OID)
		}

		c.removeTrigger(found)
	}
	return nil
}

// removePolicyObjects handles DROP POLICY for one or more policies.
// DROP POLICY uses the same pattern as DROP TRIGGER: Objects contains
// sublists of [table, policy] or [schema, table, policy].
//
// pg: src/backend/commands/policy.c — RemoveRoleFromObjectPolicy (drop path)
func (c *Catalog) removePolicyObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		list, ok := obj.(*nodes.List)
		if !ok {
			continue
		}

		var schemaName, tableName, policyName string
		items := list.Items
		switch len(items) {
		case 2:
			tableName = stringVal(items[0])
			policyName = stringVal(items[1])
		case 3:
			schemaName = stringVal(items[0])
			tableName = stringVal(items[1])
			policyName = stringVal(items[2])
		default:
			continue
		}

		_, rel, err := c.findRelation(schemaName, tableName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("policy %q for table %q does not exist, skipping", policyName, tableName))
					continue
				}
			}
			return err
		}

		var found *Policy
		for _, p := range c.policiesByRel[rel.OID] {
			if p.Name == policyName {
				found = p
				break
			}
		}

		if found == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("policy %q for table %q does not exist, skipping", policyName, tableName))
				continue
			}
			return &Error{
				Code:    CodeUndefinedObject,
				Message: fmt.Sprintf("policy %q for table %q does not exist", policyName, tableName),
			}
		}

		delete(c.policies, found.OID)
		c.removeDepsOf('p', found.OID)
		c.removeComments('p', found.OID)
		// Remove from policiesByRel slice.
		plist := c.policiesByRel[rel.OID]
		for i, p := range plist {
			if p.OID == found.OID {
				c.policiesByRel[rel.OID] = append(plist[:i], plist[i+1:]...)
				break
			}
		}
	}
	return nil
}

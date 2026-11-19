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
	type target struct {
		schema *Schema
		name   string
		idx    *Index
	}
	var targets []target
	for _, obj := range stmt.Objects.Items {
		schemaName, idxName := extractDropObjectName(obj)

		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("index %q does not exist, skipping", idxName))
				continue
			}
			return err
		}

		idx, err := c.findNamedIndex(schema, idxName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("index %q does not exist, skipping", idxName))
				continue
			}
			return err
		}

		if idx.ConstraintOID != 0 {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop index %q because constraint requires it", idxName),
			}
		}
		targets = append(targets, target{schema: schema, name: idxName, idx: idx})
	}
	for _, t := range targets {
		c.removeIndex(t.schema, t.name, t.idx)
	}
	return nil
}

// removeSequenceObjects handles DROP SEQUENCE for one or more sequences.
func (c *Catalog) removeSequenceObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}
	type target struct {
		seq  *Sequence
		name string
	}
	var targets []target
	for _, obj := range stmt.Objects.Items {
		schemaName, seqName := extractDropObjectName(obj)

		seq, err := c.findSequence(schemaName, seqName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("sequence %q does not exist, skipping", seqName))
				continue
			}
			return err
		}
		targets = append(targets, target{seq: seq, name: seqName})
	}
	for _, t := range targets {
		// Check dependents.
		if deps := c.findNormalDependents('s', t.seq.OID); len(deps) > 0 {
			if !cascade {
				return errDependentObjects("sequence", t.name)
			}
			c.dropDependents('s', t.seq.OID)
		}

		c.removeSequence(t.seq.Schema, t.seq)
	}
	return nil
}

// removeTypeObjects handles DROP TYPE for one or more types (enum or domain).
func (c *Catalog) removeTypeObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}
	type target struct {
		bt       *BuiltinType
		typeName string
		compRel  *Relation
	}
	var targets []target
	dropTypeOIDs := map[uint32]bool{}
	dropRelOIDs := map[uint32]bool{}
	dropDomainConOIDs := map[uint32]bool{}
	for _, obj := range stmt.Objects.Items {
		// The parser gives a plain name as a list of strings, and a type name
		// with modifiers, array bounds or SETOF as a TypeName. SETOF names the
		// type itself.
		// pg: src/backend/catalog/objectaddress.c — get_object_address_type
		var tn *nodes.TypeName
		var schemaName, typeName string
		if t, ok := obj.(*nodes.TypeName); ok {
			tn = t
			schemaName, typeName = qualifiedName(t.Names)
		} else {
			schemaName, typeName = extractDropObjectName(obj)
		}
		shownName := typeName
		if tn != nil && tn.ArrayBounds != nil {
			shownName += "[]"
		}

		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("type %q does not exist, skipping", shownName))
					continue
				}
			}
			return err
		}

		bt := c.typeByName[typeKey{ns: schema.OID, name: typeName}]
		if bt == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("type %q does not exist, skipping", shownName))
				continue
			}
			return errUndefinedType(shownName)
		}

		isUserBaseOrShell := bt.OID >= FirstNormalObjectId && bt.Elem == 0 && (bt.Type == 'b' || !bt.IsDefined)
		if bt.Type != 'e' && bt.Type != 'd' && bt.Type != 'r' && !(bt.Type == 'c' && bt.RelID != 0) && !isUserBaseOrShell {
			return errWrongObjectType(typeName, "a droppable type")
		}

		if tn != nil && tn.Typmods != nil {
			// pg: src/backend/parser/parse_type.c — typenameTypeMod
			return &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("type modifier is not allowed for type %q", typeName)}
		}
		if tn != nil && tn.ArrayBounds != nil {
			// The array type is the object named; its element type owns it.
			if nodes.ObjectType(stmt.RemoveType) == nodes.OBJECT_DOMAIN {
				return errWrongObjectType(shownName, "a domain")
			}
			return &Error{Code: CodeDependentObjects, Message: fmt.Sprintf("cannot drop type %s because type %s requires it", shownName, typeName)}
		}

		var compRel *Relation
		if bt.Type == 'c' && bt.RelID != 0 {
			if rel := c.relationByOID[bt.RelID]; rel != nil && rel.RelKind == 'c' {
				compRel = rel
				dropRelOIDs[rel.OID] = true
			}
		}
		for _, oid := range c.typeFamilyOIDs(bt) {
			dropTypeOIDs[oid] = true
		}
		if dt := c.domainTypes[bt.OID]; dt != nil {
			for _, dc := range dt.Constraints {
				dropDomainConOIDs[dc.OID] = true
			}
		}
		targets = append(targets, target{bt: bt, typeName: typeName, compRel: compRel})
	}

	if !cascade {
		for _, t := range targets {
			for _, relOID := range c.typeRelationDependents(t.bt) {
				if !dropRelOIDs[relOID] {
					return errDependentObjects("type", t.typeName)
				}
			}
			for _, dep := range c.typeNormalDependents(t.bt) {
				if (dep.ObjType == 't' && dropTypeOIDs[dep.ObjOID]) ||
					(dep.ObjType == 'r' && dropRelOIDs[dep.ObjOID]) ||
					(dep.ObjType == 'c' && dropDomainConOIDs[dep.ObjOID]) {
					continue
				}
				return errDependentObjects("type", t.typeName)
			}
		}
	}

	for _, t := range targets {
		if c.typeByOID[t.bt.OID] == nil {
			continue
		}
		if cascade {
			for _, relOID := range c.typeRelationDependents(t.bt) {
				if dropRelOIDs[relOID] {
					continue
				}
				if depRel := c.relationByOID[relOID]; depRel != nil {
					c.dropDependents('r', depRel.OID)
					c.removeRelation(depRel.Schema, depRel.Name, depRel)
				}
			}
			c.dropTypeDependents(t.bt)
		}
		if t.compRel != nil {
			if c.relationByOID[t.compRel.OID] != nil {
				c.removeRelation(t.compRel.Schema, t.compRel.Name, t.compRel)
			}
			continue
		}
		c.deleteUserType(t.bt)
	}
	return nil
}

// typeFamilyOIDs returns the OIDs that belong to bt as one logical type: the
// type itself, its array type, and for a range type its companion multirange
// and multirange array.
func (c *Catalog) typeFamilyOIDs(bt *BuiltinType) []uint32 {
	if bt == nil {
		return nil
	}
	oids := []uint32{bt.OID}
	if bt.Array != 0 {
		oids = append(oids, bt.Array)
	}
	if bt.Type == 'r' {
		if rt := c.rangeTypes[bt.OID]; rt != nil && rt.MultirangeOID != 0 {
			oids = append(oids, rt.MultirangeOID)
			if mr := c.typeByOID[rt.MultirangeOID]; mr != nil && mr.Array != 0 {
				oids = append(oids, mr.Array)
			}
		}
	}
	return oids
}

func (c *Catalog) typeRelationDependents(bt *BuiltinType) []uint32 {
	var out []uint32
	seen := map[uint32]bool{}
	for _, oid := range c.typeFamilyOIDs(bt) {
		for _, relOID := range c.findTypeDependents(oid) {
			if !seen[relOID] {
				seen[relOID] = true
				out = append(out, relOID)
			}
		}
	}
	return out
}

func (c *Catalog) typeNormalDependents(bt *BuiltinType) []DepEntry {
	var out []DepEntry
	family := map[uint32]bool{}
	for _, oid := range c.typeFamilyOIDs(bt) {
		family[oid] = true
	}
	for _, oid := range c.typeFamilyOIDs(bt) {
		for _, dep := range c.findNormalDependents('t', oid) {
			if dep.ObjType == 't' && family[dep.ObjOID] {
				continue
			}
			out = append(out, dep)
		}
	}
	return out
}

func (c *Catalog) dropTypeDependents(bt *BuiltinType) {
	for _, oid := range c.typeFamilyOIDs(bt) {
		c.dropDependents('t', oid)
	}
}

// deleteUserType removes a non-composite user-defined type and its companion
// array/multirange entries, ownership records, grants, comments, and deps.
func (c *Catalog) deleteUserType(bt *BuiltinType) {
	if bt == nil {
		return
	}
	if dt := c.domainTypes[bt.OID]; dt != nil {
		for _, dc := range dt.Constraints {
			c.removeDepsOf('c', dc.OID)
		}
	}
	if bt.Type == 'r' {
		if rt := c.rangeTypes[bt.OID]; rt != nil && rt.MultirangeOID != 0 {
			if mr := c.typeByOID[rt.MultirangeOID]; mr != nil {
				if mr.Array != 0 {
					if mra := c.typeByOID[mr.Array]; mra != nil {
						delete(c.typeByName, typeKey{ns: mra.Namespace, name: mra.TypeName})
					}
					delete(c.typeByOID, mr.Array)
					delete(c.typeOwners, mr.Array)
					c.removeGrantsForObject('T', mr.Array)
					c.removeComments('t', mr.Array)
					c.removeDepsOf('t', mr.Array)
					c.removeDepsOn('t', mr.Array)
				}
				delete(c.typeByName, typeKey{ns: mr.Namespace, name: mr.TypeName})
			}
			delete(c.typeByOID, rt.MultirangeOID)
			delete(c.typeOwners, rt.MultirangeOID)
			c.removeGrantsForObject('T', rt.MultirangeOID)
			c.removeComments('t', rt.MultirangeOID)
			c.removeDepsOf('t', rt.MultirangeOID)
			c.removeDepsOn('t', rt.MultirangeOID)
		}
	}

	switch bt.Type {
	case 'e':
		delete(c.enumTypes, bt.OID)
	case 'd':
		delete(c.domainTypes, bt.OID)
	case 'r':
		delete(c.rangeTypes, bt.OID)
	}

	if bt.Array != 0 {
		if at := c.typeByOID[bt.Array]; at != nil {
			delete(c.typeByName, typeKey{ns: at.Namespace, name: at.TypeName})
		}
		delete(c.typeByOID, bt.Array)
		delete(c.typeOwners, bt.Array)
		c.removeGrantsForObject('T', bt.Array)
		c.removeComments('t', bt.Array)
		c.removeDepsOf('t', bt.Array)
		c.removeDepsOn('t', bt.Array)
	}

	delete(c.typeByName, typeKey{ns: bt.Namespace, name: bt.TypeName})
	delete(c.typeByOID, bt.OID)
	delete(c.typeOwners, bt.OID)
	c.removeGrantsForObject('T', bt.OID)
	c.removeComments('t', bt.OID)
	c.removeDepsOf('t', bt.OID)
	c.removeDepsOn('t', bt.OID)
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

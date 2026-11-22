package catalog

import (
	"fmt"
	"sort"

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
	case nodes.OBJECT_EXTENSION:
		return c.removeExtensionObjects(stmt)

	case nodes.OBJECT_COLLATION:
		return c.removeCollationObjects(stmt)
	case nodes.OBJECT_RULE:
		return c.removeRuleObjects(stmt)
	case nodes.OBJECT_OPERATOR:
		return c.removeOperatorObjects(stmt)
	case nodes.OBJECT_ACCESS_METHOD:
		return c.removeAccessMethodObjects(stmt)
	case nodes.OBJECT_FDW:
		return c.removeFdwObjects(stmt)
	case nodes.OBJECT_FOREIGN_SERVER:
		return c.removeForeignServerObjects(stmt)
	case nodes.OBJECT_OPCLASS:
		return c.removeOpClassObjects(stmt)
	case nodes.OBJECT_CONVERSION,
		nodes.OBJECT_OPFAMILY,
		nodes.OBJECT_LANGUAGE,
		nodes.OBJECT_PUBLICATION,
		nodes.OBJECT_CAST,
		nodes.OBJECT_TRANSFORM,
		nodes.OBJECT_STATISTIC_EXT,
		nodes.OBJECT_TSPARSER,
		nodes.OBJECT_TSDICTIONARY,
		nodes.OBJECT_TSTEMPLATE,
		nodes.OBJECT_TSCONFIGURATION:
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
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	type target struct {
		schema *Schema
		name   string
		idx    *Index
	}
	var targets []target
	dropIdxOIDs := map[uint32]bool{}
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

		if c.isExtensionMember(idx.OID) {
			return errDependentObjects("extension", c.extensionMembers[idx.OID])
		}
		if idx.ConstraintOID != 0 {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop index %q because constraint requires it", idxName),
			}
		}
		if idx.ParentIndexOID != 0 {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop index %q because parent index requires it", idxName),
			}
		}
		if !dropIdxOIDs[idx.OID] {
			dropIdxOIDs[idx.OID] = true
			targets = append(targets, target{schema: schema, name: idxName, idx: idx})
		}
	}
	for changed := true; changed; {
		changed = false
		var added []target
		for oid, child := range c.indexes {
			if child != nil && child.Schema != nil && child.ParentIndexOID != 0 && dropIdxOIDs[child.ParentIndexOID] && !dropIdxOIDs[oid] {
				dropIdxOIDs[oid] = true
				added = append(added, target{schema: child.Schema, name: child.Name, idx: child})
				changed = true
			}
		}
		sort.Slice(added, func(i, j int) bool { return added[i].idx.OID < added[j].idx.OID })
		targets = append(targets, added...)
	}
	if !cascade {
		for _, t := range targets {
			if t.idx.ConstraintOID != 0 {
				return &Error{
					Code:    CodeDependentObjects,
					Message: fmt.Sprintf("cannot drop index %q because constraint requires it", t.name),
				}
			}
			for _, refClass := range []byte{'r', 'i'} {
				for _, dep := range c.findNormalDependents(refClass, t.idx.OID) {
					if dep.ObjType == 'i' && dropIdxOIDs[dep.ObjOID] {
						continue
					}
					return errDependentObjects("index", t.name)
				}
			}
		}
	}
	for _, t := range targets {
		if c.indexes[t.idx.OID] != t.idx {
			continue
		}
		if cascade {
			c.dropDependentObject('i', t.idx.OID)
		} else {
			c.removeIndex(t.schema, t.name, t.idx)
		}
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
	dropSeqOIDs := map[uint32]bool{}
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
		if c.isExtensionMember(seq.OID) {
			return errDependentObjects("extension", c.extensionMembers[seq.OID])
		}
		if !dropSeqOIDs[seq.OID] {
			dropSeqOIDs[seq.OID] = true
			targets = append(targets, target{seq: seq, name: seqName})
		}
	}
	if !cascade {
		for _, t := range targets {
			for _, dep := range c.findNormalDependents('s', t.seq.OID) {
				if dep.ObjType == 's' && dropSeqOIDs[dep.ObjOID] {
					continue
				}
				return errDependentObjects("sequence", t.name)
			}
		}
	}
	for _, t := range targets {
		if c.sequenceByOID[t.seq.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('s', t.seq.OID)) > 0 {
			c.dropDependents('s', t.seq.OID)
		}
		if c.sequenceByOID[t.seq.OID] != nil {
			c.removeSequence(t.seq.Schema, t.seq)
		}
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

		if tn != nil && tn.Typmods != nil {
			// pg: src/backend/parser/parse_type.c — typenameTypeMod
			if bt.ModIn == 0 {
				return &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("type modifier is not allowed for type %q", typeName)}
			}
			if _, _, modErr := c.resolveTypeName(tn); modErr != nil {
				return modErr
			}
		}
		if nodes.ObjectType(stmt.RemoveType) == nodes.OBJECT_DOMAIN {
			if (tn != nil && tn.ArrayBounds != nil) || bt.Type != 'd' {
				return errWrongObjectType(shownName, "a domain")
			}
		}
		if tn != nil && tn.ArrayBounds != nil {
			// The array type is the object named; its element type owns it.
			return &Error{Code: CodeDependentObjects, Message: fmt.Sprintf("cannot drop type %s because type %s requires it", shownName, typeName)}
		}

		// Built-in / system types are pinned and cannot be dropped (2BP01).
		// pg: src/backend/catalog/dependency.c — findDependentObjects (DEPENDENCY_PIN)
		if bt.OID < FirstNormalObjectId {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop type %s because it is required by the database system", bt.TypeName),
			}
		}

		isUserBaseOrShell := bt.OID >= FirstNormalObjectId && bt.Elem == 0 && (bt.Type == 'b' || !bt.IsDefined)
		if bt.Type != 'e' && bt.Type != 'd' && bt.Type != 'r' && !(bt.Type == 'c' && bt.RelID != 0) && !isUserBaseOrShell {
			return errWrongObjectType(typeName, "a droppable type")
		}

		if c.isExtensionMember(bt.OID) {
			return errDependentObjects("extension", c.extensionMembers[bt.OID])
		}

		var compRel *Relation
		if bt.Type == 'c' && bt.RelID != 0 {
			if rel := c.relationByOID[bt.RelID]; rel != nil {
				if rel.RelKind == 'c' {
					compRel = rel
					dropRelOIDs[rel.OID] = true
				} else {
					// Table/view/matview/sequence row type has an internal dependency on its relation;
					// PostgreSQL refuses directly dropping it even with CASCADE.
					return &Error{
						Code:    CodeDependentObjects,
						Message: fmt.Sprintf("cannot drop type %s because table %s requires it", typeName, rel.Name),
					}
				}
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
			if t.compRel != nil {
				for _, dep := range c.relationDependents(t.compRel.OID, 0) {
					if dep.Normal && !(dep.ObjType == 'r' && dropRelOIDs[dep.ObjOID]) {
						return errDependentObjects("type", t.typeName)
					}
				}
			}
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
			famOIDs := map[uint32]bool{}
			for _, oid := range c.typeFamilyOIDs(t.bt) {
				famOIDs[oid] = true
			}
			for _, cast := range c.castIndex {
				if famOIDs[cast.Source] || famOIDs[cast.Target] {
					return errDependentObjects("type", t.typeName)
				}
			}
		}
	}

	for _, t := range targets {
		if c.typeByOID[t.bt.OID] == nil {
			continue
		}
		if cascade {
			done, _ := c.beginDrop('t', t.bt.OID, c.typeFamilyOIDs(t.bt)...)
			for _, relOID := range c.typeRelationDependents(t.bt) {
				if dropRelOIDs[relOID] {
					continue
				}
				if depRel := c.relationByOID[relOID]; depRel != nil {
					c.cascadeDropRelation(depRel)
				}
			}
			c.dropTypeDependents(t.bt)
			if done != nil {
				done()
			}
		}
		if t.compRel != nil {
			if c.relationByOID[t.compRel.OID] != nil {
				if cascade {
					c.cascadeDropRelation(t.compRel)
				} else {
					c.removeRelation(t.compRel.Schema, t.compRel.Name, t.compRel)
				}
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
	famOIDs := map[uint32]bool{}
	for _, oid := range c.typeFamilyOIDs(bt) {
		famOIDs[oid] = true
	}
	for key, cast := range c.castIndex {
		if famOIDs[cast.Source] || famOIDs[cast.Target] {
			delete(c.castIndex, key)
			delete(c.extensionCasts, key)
		}
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
					delete(c.extensionMembers, mr.Array)
					c.removeGrantsForObject('T', mr.Array)
					c.removeComments('t', mr.Array)
					c.removeDepsOf('t', mr.Array)
					c.removeDepsOn('t', mr.Array)
				}
				delete(c.typeByName, typeKey{ns: mr.Namespace, name: mr.TypeName})
			}
			delete(c.typeByOID, rt.MultirangeOID)
			delete(c.typeOwners, rt.MultirangeOID)
			delete(c.extensionMembers, rt.MultirangeOID)
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
		delete(c.extensionMembers, bt.Array)
		c.removeGrantsForObject('T', bt.Array)
		c.removeComments('t', bt.Array)
		c.removeDepsOf('t', bt.Array)
		c.removeDepsOn('t', bt.Array)
	}

	delete(c.typeByName, typeKey{ns: bt.Namespace, name: bt.TypeName})
	delete(c.typeByOID, bt.OID)
	delete(c.typeOwners, bt.OID)
	delete(c.extensionMembers, bt.OID)
	c.removeGrantsForObject('T', bt.OID)
	c.removeComments('t', bt.OID)
	c.removeDepsOf('t', bt.OID)
	c.removeDepsOn('t', bt.OID)
}

// removeFunctionObjects handles DROP FUNCTION/PROCEDURE/ROUTINE/AGGREGATE for one or more routines.
func (c *Catalog) removeFunctionObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	objType := nodes.ObjectType(stmt.RemoveType)
	dropKindStr := "function"
	switch objType {
	case nodes.OBJECT_PROCEDURE:
		dropKindStr = "procedure"
	case nodes.OBJECT_AGGREGATE:
		dropKindStr = "aggregate"
	case nodes.OBJECT_ROUTINE:
		dropKindStr = "routine"
	}
	if stmt.Objects == nil {
		return nil
	}
	type target struct {
		bp   *BuiltinProc
		name string
	}
	var targets []target
	dropProcOIDs := map[uint32]bool{}
	for _, obj := range stmt.Objects.Items {
		owa, ok := obj.(*nodes.ObjectWithArgs)
		if !ok {
			continue
		}
		_, funcName := qualifiedName(owa.Objname)
		resolveType := objType
		if resolveType == nodes.OBJECT_AGGREGATE {
			resolveType = nodes.OBJECT_ROUTINE
		}
		_, bp, _, err := c.resolveRoutineTarget(resolveType, owa)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("%s %s does not exist, skipping", dropKindStr, funcName))
				continue
			}
			return err
		}
		if objType == nodes.OBJECT_AGGREGATE && bp.Kind == 'p' {
			return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not an aggregate", funcName)}
		}
		if c.isExtensionMember(bp.OID) {
			return errDependentObjects("extension", c.extensionMembers[bp.OID])
		}
		if !dropProcOIDs[bp.OID] {
			dropProcOIDs[bp.OID] = true
			targets = append(targets, target{bp: bp, name: funcName})
		}
	}

	if !cascade {
		for _, t := range targets {
			for _, dep := range c.findNormalDependents('f', t.bp.OID) {
				if dep.ObjType == 'f' && dropProcOIDs[dep.ObjOID] {
					continue
				}
				return errDependentObjects(dropKindStr, t.name)
			}
			for _, cast := range c.castIndex {
				if cast.Func == t.bp.OID {
					return errDependentObjects(dropKindStr, t.name)
				}
			}
		}
	}

	for _, t := range targets {
		if c.userProcs[t.bp.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('f', t.bp.OID)) > 0 {
			done, _ := c.beginDrop('f', t.bp.OID)
			c.dropDependents('f', t.bp.OID)
			if done != nil {
				done()
			}
		}
		if c.userProcs[t.bp.OID] != nil {
			c.removeFunction(t.bp.OID, t.name)
		}
	}
	return nil
}

// removeTriggerObjects handles DROP TRIGGER for one or more triggers.
func (c *Catalog) removeTriggerObjects(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	if stmt.Objects == nil {
		return nil
	}
	var targets []*Trigger
	dropTrigOIDs := map[uint32]bool{}
	for _, obj := range stmt.Objects.Items {
		list, ok := obj.(*nodes.List)
		if !ok {
			continue
		}

		var schemaName, tableName, trigName string
		items := list.Items
		switch len(items) {
		case 2:
			tableName = stringVal(items[0])
			trigName = stringVal(items[1])
		case 3:
			schemaName = stringVal(items[0])
			tableName = stringVal(items[1])
			trigName = stringVal(items[2])
		default:
			continue
		}

		rel, isSeq, err := c.openTableLikeRelation(schemaName, tableName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("trigger %q for relation %q does not exist, skipping", trigName, tableName))
				continue
			}
			return err
		}

		var found *Trigger
		if !isSeq && rel != nil {
			for _, trig := range c.triggersByRel[rel.OID] {
				if trig.Name == trigName {
					found = trig
					break
				}
			}
		}

		if found == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("trigger %q for relation %q does not exist, skipping", trigName, tableName))
				continue
			}
			return errUndefinedTrigger(trigName, tableName)
		}
		if c.isExtensionMember(found.OID) {
			return errDependentObjects("extension", c.extensionMembers[found.OID])
		}
		if found.ParentTriggerOID != 0 {
			return errDependentObjects("trigger", found.Name)
		}
		if !dropTrigOIDs[found.OID] {
			dropTrigOIDs[found.OID] = true
			targets = append(targets, found)
		}
	}
	if !cascade {
		for _, trig := range targets {
			for _, dep := range c.findNormalDependents('g', trig.OID) {
				if dep.ObjType == 'g' && dropTrigOIDs[dep.ObjOID] {
					continue
				}
				return errDependentObjects("trigger", trig.Name)
			}
		}
	}
	for _, trig := range targets {
		if c.triggers[trig.OID] == nil {
			continue
		}
		if cascade {
			c.dropDependents('g', trig.OID)
		}
		c.removeTrigger(trig)
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
	var targets []*Policy
	seen := map[uint32]bool{}
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

		rel, isSeq, err := c.openTableLikeRelation(schemaName, tableName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("policy %q for table %q does not exist, skipping", policyName, tableName))
				continue
			}
			return err
		}

		var found *Policy
		if !isSeq && rel != nil {
			for _, p := range c.policiesByRel[rel.OID] {
				if p.Name == policyName {
					found = p
					break
				}
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
		if c.isExtensionMember(found.OID) {
			return errDependentObjects("extension", c.extensionMembers[found.OID])
		}
		if !seen[found.OID] {
			seen[found.OID] = true
			targets = append(targets, found)
		}
	}
	for _, pol := range targets {
		c.removePolicy(pol)
	}
	return nil
}

func (c *Catalog) removeCollationObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	type collTarget struct {
		schema *Schema
		name   string
		oid    uint32
	}
	var targets []collTarget
	seen := map[uint32]bool{}

	for _, obj := range stmt.Objects.Items {
		schemaName, name := extractDropObjectName(obj)
		var targetSchema *Schema
		var collOID uint32
		if schemaName != "" {
			s := c.schemaByName[schemaName]
			if s == nil {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("collation %q does not exist, skipping", name))
					continue
				}
				return errUndefinedSchema(schemaName)
			}
			if s.OID == PGCatalogNamespace && isBuiltinCollationName(name) {
				return errDependentObjects("system", name)
			}
			targetSchema = s
			if uc := s.Collations[name]; uc != nil {
				collOID = uc.OID
			}
		} else {
			for _, nsOID := range c.searchPathWithCatalog() {
				s := c.schemas[nsOID]
				if s == nil {
					continue
				}
				if nsOID == PGCatalogNamespace && isBuiltinCollationName(name) {
					return errDependentObjects("system", name)
				}
				if uc := s.Collations[name]; uc != nil {
					targetSchema = s
					collOID = uc.OID
					break
				}
			}
		}
		if collOID == 0 {
			if stmt.Missing_ok || c.snapshot {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("collation %q does not exist, skipping", name))
				}
				continue
			}
			return errUndefinedObject("collation", name)
		}
		if c.isExtensionMember(collOID) {
			return errDependentObjects("extension", c.extensionMembers[collOID])
		}
		if !seen[collOID] {
			seen[collOID] = true
			targets = append(targets, collTarget{schema: targetSchema, name: name, oid: collOID})
		}
	}

	if !cascade {
		for _, t := range targets {
			if len(c.findNormalDependents('L', t.oid)) > 0 {
				return errDependentObjects("collation", t.name)
			}
		}
	}

	restore := func() {}
	if cascade && len(targets) > 0 {
		restore = c.snapshotForAlterRollback()
	}

	for _, t := range targets {
		uc := t.schema.Collations[t.name]
		if uc == nil || uc.OID != t.oid {
			continue
		}
		if cascade {
			if err := c.cascadeDropCollation(t.schema, t.name, t.oid); err != nil {
				restore()
				return err
			}
		} else {
			delete(t.schema.Collations, t.name)
			c.removeDepsOn('L', t.oid)
		}
	}
	return nil
}

func (c *Catalog) cascadeDropCollation(schema *Schema, name string, collOID uint32) error {
	deps := c.findNormalDependents('L', collOID)

	type colTarget struct {
		relOID  uint32
		colName string
	}
	var colDrops []colTarget
	seenCol := map[colTarget]bool{}

	for _, dep := range deps {
		if dep.ObjType != 'r' || dep.ObjSubID <= 0 {
			continue
		}
		rel := c.relationByOID[dep.ObjOID]
		if rel == nil || rel.RelKind == 'c' {
			continue
		}
		for _, col := range rel.Columns {
			if col.AttNum == int16(dep.ObjSubID) {
				if col.UserCollationOID == collOID || col.Generated == 's' {
					ct := colTarget{relOID: rel.OID, colName: col.Name}
					if !seenCol[ct] {
						seenCol[ct] = true
						colDrops = append(colDrops, ct)
					}
				} else if col.HasDefault && col.Generated == 0 {
					c.clearColumnDefaultDeps(rel, col)
					col.HasDefault = false
					col.Default = ""
					col.DefaultAnalyzed = nil
				}
				break
			}
		}
	}

	// Drop root column definitions (InhCount == 0) first and recurse down
	// inheritance/partition hierarchies, repeating until a fixed point so
	// merged child columns whose InhCount drops to 0 are also removed.
	for progressed := true; progressed; {
		progressed = false
		for _, ct := range colDrops {
			rel := c.relationByOID[ct.relOID]
			if rel == nil {
				continue
			}
			idx, exists := rel.colByName[ct.colName]
			if !exists {
				continue
			}
			col := rel.Columns[idx]
			if col.InhCount != 0 {
				continue
			}
			if err := c.checkChildColumnDrops(rel, col.Name, true); err != nil {
				return err
			}
			if err := c.atDropColumn(rel.Schema, rel, col.Name, true, true, false); err != nil {
				return err
			}
			for _, childOID := range c.findChildRelations(rel.OID) {
				if child := c.relationByOID[childOID]; child != nil {
					if err := c.atDropColumn(child.Schema, child, ct.colName, true, true, true); err != nil {
						return err
					}
				}
			}
			progressed = true
		}
	}

	// If any dependent column still remains (e.g. an inherited column whose
	// parent did not depend on collOID), attempt the drop so PostgreSQL's
	// refusal (42P16) aborts DROP COLLATION before deleting the collation.
	for _, ct := range colDrops {
		rel := c.relationByOID[ct.relOID]
		if rel == nil {
			continue
		}
		if _, exists := rel.colByName[ct.colName]; exists {
			if err := c.atDropColumn(rel.Schema, rel, ct.colName, true, false, false); err != nil {
				return err
			}
		}
	}

	// Drop remaining non-column dependents (composites, domains, ranges, views, indexes, etc.).
	for _, dep := range c.findNormalDependents('L', collOID) {
		if dep.ObjType == 'r' && dep.ObjSubID > 0 {
			if rel := c.relationByOID[dep.ObjOID]; rel != nil && rel.RelKind == 'c' {
				c.cascadeDropRelation(rel)
			}
			continue
		}
		c.dropDependentObject(dep.ObjType, dep.ObjOID)
	}

	if schema != nil {
		delete(schema.Collations, name)
	}
	c.removeDepsOn('L', collOID)
	return nil
}

// removeRuleObjects handles DROP RULE [IF EXISTS] name ON rel [CASCADE|RESTRICT].
//
// pg: src/backend/rewrite/rewriteRemove.c — get_rewrite_oid, RemoveRewriteRuleById
func (c *Catalog) removeRuleObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	type ruleTarget struct {
		rel  *Relation
		rule *Rule
	}
	var targets []ruleTarget
	seenRule := map[uint32]bool{}

	for _, obj := range stmt.Objects.Items {
		list, ok := obj.(*nodes.List)
		if !ok || len(list.Items) < 2 {
			continue
		}
		var schemaName, relName string
		if len(list.Items) == 2 {
			relName = stringVal(list.Items[0])
		} else {
			schemaName = stringVal(list.Items[len(list.Items)-3])
			relName = stringVal(list.Items[len(list.Items)-2])
		}
		ruleName := stringVal(list.Items[len(list.Items)-1])

		rel, otherRelExists, err := c.findRuleTargetRelation(schemaName, relName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("rule %q for relation %q does not exist, skipping", ruleName, relName))
				continue
			}
			return err
		}
		if otherRelExists {
			if stmt.Missing_ok || c.snapshot {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("rule %q for relation %q does not exist, skipping", ruleName, relName))
				}
				continue
			}
			return &Error{
				Code:    CodeUndefinedObject,
				Message: fmt.Sprintf("rule %q for relation %q does not exist", ruleName, relName),
			}
		}

		if rel.RelKind == 'v' && ruleName == "_RETURN" {
			if c.isExtensionMember(rel.OID) {
				return errDependentObjects("extension", c.extensionMembers[rel.OID])
			}
			return errDependentObjects("view", rel.Name)
		}

		var rule *Rule
		if rel.Rules != nil {
			rule = rel.Rules[ruleName]
		}
		if rule == nil {
			if stmt.Missing_ok || c.snapshot {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("rule %q for relation %q does not exist, skipping", ruleName, rel.Name))
				}
				continue
			}
			return &Error{
				Code:    CodeUndefinedObject,
				Message: fmt.Sprintf("rule %q for relation %q does not exist", ruleName, rel.Name),
			}
		}
		if c.isExtensionMember(rule.OID) {
			return errDependentObjects("extension", c.extensionMembers[rule.OID])
		}
		if !seenRule[rule.OID] {
			seenRule[rule.OID] = true
			targets = append(targets, ruleTarget{rel: rel, rule: rule})
		}
	}

	if !cascade {
		for _, t := range targets {
			if t.rule != nil && len(c.findNormalDependents('R', t.rule.OID)) > 0 {
				return errDependentObjects("rule", t.rule.Name)
			}
		}
	}

	for _, t := range targets {
		if t.rule == nil || c.rules[t.rule.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('R', t.rule.OID)) > 0 {
			c.dropDependents('R', t.rule.OID)
		}
		if c.rules[t.rule.OID] != nil {
			c.removeRule(t.rel, t.rule)
		}
	}
	return nil
}

func (c *Catalog) findRuleTargetRelation(schemaName, relName string) (*Relation, bool, error) {
	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, false, errUndefinedSchema(schemaName)
		}
		if r := s.Relations[relName]; r != nil && r.RelKind != 'c' {
			return r, false, nil
		}
		if s.Sequences[relName] != nil || s.Indexes[relName] != nil {
			return nil, true, nil
		}
		return nil, false, errUndefinedTable(relName)
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		if r := s.Relations[relName]; r != nil && r.RelKind != 'c' {
			return r, false, nil
		}
		if s.Sequences[relName] != nil || s.Indexes[relName] != nil {
			return nil, true, nil
		}
	}
	return nil, false, errUndefinedTable(relName)
}

func (c *Catalog) removeOperatorObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	var targets []*BuiltinOperator
	dropOpOIDs := make(map[uint32]bool)
	for _, item := range stmt.Objects.Items {
		owa, ok := item.(*nodes.ObjectWithArgs)
		if !ok || owa == nil {
			continue
		}
		op, err := c.lookupUserOperator(owa)
		if err != nil {
			if (stmt.Missing_ok || c.snapshot) && missingObject(err) {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, "operator does not exist, skipping")
				}
				continue
			}
			return err
		}
		if c.isExtensionMember(op.OID) {
			return errDependentObjects("extension", c.extensionMembers[op.OID])
		}
		if !dropOpOIDs[op.OID] {
			dropOpOIDs[op.OID] = true
			targets = append(targets, op)
		}
	}
	if !cascade {
		for _, op := range targets {
			for _, dep := range c.findNormalDependents('O', op.OID) {
				if !(dep.ObjType == 'O' && dropOpOIDs[dep.ObjOID]) {
					return errDependentObjects("operator", op.Name)
				}
			}
		}
	}
	for _, op := range targets {
		c.detachOperatorFromOpFamilies(op)
		if cascade {
			c.dropDependents('O', op.OID)
		}
		c.removeOperator(op.OID)
	}
	return nil
}

// detachOperatorFromOpFamilies removes a dropped operator from the families
// that list it as a loose member (ALTER OPERATOR FAMILY ... ADD), which
// depend on it automatically: the family stays.
func (c *Catalog) detachOperatorFromOpFamilies(op *BuiltinOperator) {
	var fams []uint32
	for _, d := range c.deps {
		if d.ObjType == 'Q' && d.RefType == 'O' && d.RefOID == op.OID && d.DepType == DepAuto {
			fams = append(fams, d.ObjOID)
		}
	}
	for _, famOID := range fams {
		c.removeDep('Q', famOID, 'O', op.OID)
		fam := c.opFamilies[famOID]
		if fam == nil {
			continue
		}
		// Another overload of the same name may remain a member.
		stillListed := false
		for _, d := range c.deps {
			if d.ObjType == 'Q' && d.ObjOID == famOID && d.RefType == 'O' && d.RefOID != op.OID {
				if other := c.operByOID[d.RefOID]; other != nil && other.Name == op.Name {
					stillListed = true
					break
				}
			}
		}
		if !stillListed {
			delete(fam.Operators, op.Name)
		}
	}
}

func (c *Catalog) removeAccessMethodObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	var targets []*AccessMethod
	seen := make(map[uint32]bool)
	for _, item := range stmt.Objects.Items {
		name := stringVal(item)
		if name == "" {
			if list, ok := item.(*nodes.List); ok && len(list.Items) > 0 {
				name = stringVal(list.Items[len(list.Items)-1])
			}
		}
		am := c.accessMethodByName[name]
		if am == nil {
			if stmt.Missing_ok || c.snapshot {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("access method %q does not exist, skipping", name))
				}
				continue
			}
			return errUndefinedObject("access method", name)
		}
		if am.OID < FirstNormalObjectId {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop access method %s because it is required by the database system", am.Name),
			}
		}
		if c.isExtensionMember(am.OID) {
			return errDependentObjects("extension", c.extensionMembers[am.OID])
		}
		if !seen[am.OID] {
			seen[am.OID] = true
			targets = append(targets, am)
		}
	}
	if !cascade {
		for _, am := range targets {
			if c.accessMethodNeedsCascade(am.OID) {
				return errDependentObjects("access method", am.Name)
			}
		}
	}
	for _, am := range targets {
		c.cascadeDropAccessMethod(am.OID)
	}
	return nil
}

func (c *Catalog) removeOpClassObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	var targets []*OpClass
	seen := make(map[uint32]bool)
	for _, item := range stmt.Objects.Items {
		amName, schemaName, className := extractOpClassObject(item)
		am, err := c.lookupIndexAccessMethod(amName)
		if err != nil {
			if (stmt.Missing_ok || c.snapshot) && missingObject(err) {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("access method %q does not exist, skipping", amName))
				}
				continue
			}
			return err
		}
		opc, err := c.lookupOpClass(am, schemaName, className)
		if err != nil {
			if (stmt.Missing_ok || c.snapshot) && missingObject(err) {
				if stmt.Missing_ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("operator class %q does not exist, skipping", className))
				}
				continue
			}
			return err
		}
		if c.isExtensionMember(opc.OID) {
			return errDependentObjects("extension", c.extensionMembers[opc.OID])
		}
		if !seen[opc.OID] {
			seen[opc.OID] = true
			targets = append(targets, opc)
		}
	}
	if !cascade {
		for _, opc := range targets {
			for _, refClass := range []byte{'o', 'q'} {
				if len(c.findNormalDependents(refClass, opc.OID)) > 0 {
					return errDependentObjects("operator class", opc.Name)
				}
			}
		}
	}
	for _, opc := range targets {
		if c.opClasses[opc.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropOpClass(opc.OID)
		} else {
			c.removeOpClass(opc.OID)
		}
	}
	return nil
}

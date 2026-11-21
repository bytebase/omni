package catalog

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// Extension represents a registered extension.
//
// pg: src/backend/commands/extension.c — ExtensionControlFile
type Extension struct {
	OID         uint32
	Name        string
	SchemaOID   uint32
	OwnerName   string
	Relocatable bool
	// Version is the VERSION the extension was created with, or "" when the
	// statement named none.
	Version string
}

// CreateExtension loads and executes a bundled extension DDL script.
//
// pg: src/backend/commands/extension.c — CreateExtension
func (c *Catalog) CreateExtension(stmt *nodes.CreateExtensionStmt) error {
	// IF NOT EXISTS check.
	if existing := c.extByName[stmt.Extname]; existing != nil {
		if stmt.IfNotExists || (c.installingSnapshot && existing.OID < FirstNormalObjectId) {
			if stmt.IfNotExists {
				c.addWarning(CodeWarningSkip,
					fmt.Sprintf("extension %q already exists, skipping", stmt.Extname))
			}
			return nil
		}
		return errDuplicateObject("extension", stmt.Extname)
	}

	// Parse Options for "schema" DefElem.
	// pg: src/backend/commands/extension.c — CreateExtension (option parsing)
	var schemaName, version string
	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			sv, ok := de.Arg.(*nodes.String)
			if !ok {
				continue
			}
			switch de.Defname {
			case "schema":
				schemaName = sv.Str
			case "new_version":
				version = sv.Str
			}
		}
	}

	// Resolve target schema.
	// Default: first in search path (typically "public"), or pg_catalog for plpgsql.
	var targetSchema *Schema
	var err error
	if schemaName == "" && stmt.Extname == "plpgsql" {
		targetSchema = c.schemas[PGCatalogNamespace]
	} else {
		targetSchema, err = c.resolveTargetSchema(schemaName)
		if err != nil {
			return err
		}
	}

	// An extension without a bundled script is registered with no objects.
	sql, ok := extensionScripts[stmt.Extname]
	if !ok {
		c.addWarning(CodeWarning,
			fmt.Sprintf("extension %q is not bundled; its objects are unknown", stmt.Extname))
	}

	// Parse the extension DDL script.
	list, err := pgparser.Parse(sql)
	if err != nil {
		return fmt.Errorf("extension %q: parse error: %w", stmt.Extname, err)
	}
	firstMember, firstGrant := c.oidGen.next, len(c.grants)
	commentsBefore := maps.Clone(c.comments)
	castsBefore := make(map[castKey]bool, len(c.castIndex))
	for key := range c.castIndex {
		castsBefore[key] = true
	}

	// Snapshot catalog state before executing script so a mid-script error rolls back cleanly.
	saved := c.Clone()

	// Temporarily prepend target schema to search path so unqualified types
	// and functions land in the extension's schema.
	// pg: src/backend/commands/extension.c — execute_extension_script (line 973)
	origSearchPath := c.searchPath
	if targetSchema.Name != "" {
		found := false
		for _, s := range c.searchPath {
			if s == targetSchema.Name {
				found = true
				break
			}
		}
		if !found {
			c.searchPath = append([]string{targetSchema.Name}, c.searchPath...)
		}
	}

	// Execute each statement through ProcessUtility.
	if list != nil {
		for _, item := range list.Items {
			var node nodes.Node
			if raw, ok := item.(*nodes.RawStmt); ok {
				node = raw.Stmt.(nodes.Node)
			} else {
				node = item.(nodes.Node)
			}
			if err := c.ProcessUtility(node); err != nil {
				*c = *saved
				return fmt.Errorf("extension %q: %w", stmt.Extname, err)
			}
		}
	}

	// Restore original search path.
	c.searchPath = origSearchPath
	c.markExtensionMembers(firstMember, c.oidGen.next, stmt.Extname)
	for i := firstGrant; i < len(c.grants); i++ {
		c.grants[i].FromExtension = true
	}
	for key, text := range c.comments {
		if before, ok := commentsBefore[key]; !ok || before != text {
			c.extensionComments[key] = text
		}
	}
	for key := range c.castIndex {
		if !castsBefore[key] {
			c.extensionCasts[key] = stmt.Extname
		}
	}

	// Register extension metadata.
	ext := &Extension{
		OID:         c.oidGen.Next(),
		Name:        stmt.Extname,
		SchemaOID:   targetSchema.OID,
		OwnerName:   c.currentOwner(),
		Relocatable: stmt.Extname != "plpgsql",
		Version:     version,
	}
	c.extensions[ext.OID] = ext
	c.extByName[ext.Name] = ext

	return nil
}

// removeExtensionObjects handles DROP EXTENSION for one or more extensions,
// validating external dependencies across all member objects before removal.
func (c *Catalog) removeExtensionObjects(stmt *nodes.DropStmt) error {
	if stmt.Objects == nil {
		return nil
	}
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	var exts []*Extension
	dropExtNames := map[string]bool{}
	dropExtOIDs := map[uint32]bool{}
	for _, obj := range stmt.Objects.Items {
		_, extName := extractDropObjectName(obj)
		ext := c.extByName[extName]
		if ext == nil {
			if stmt.Missing_ok {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("extension %q does not exist, skipping", extName))
				continue
			}
			if c.snapshot {
				continue
			}
			return errUndefinedObject("extension", extName)
		}
		if !dropExtNames[ext.Name] {
			dropExtNames[ext.Name] = true
			dropExtOIDs[ext.OID] = true
			exts = append(exts, ext)
		}
	}
	if len(exts) == 0 {
		return nil
	}

	autoExtOIDs := map[uint32]bool{}
	for changed := true; changed; {
		changed = false
		for _, d := range c.deps {
			if d.RefType != 'e' || !dropExtOIDs[d.RefOID] || d.DepType != DepAutoExtension {
				continue
			}
			if d.ObjOID < FirstNormalObjectId {
				return &Error{
					Code:    CodeDependentObjects,
					Message: "cannot drop system object because it is required by the database system",
				}
			}
			autoExtOIDs[d.ObjOID] = true
			if ownerExtName := c.extensionMembers[d.ObjOID]; ownerExtName != "" && !dropExtNames[ownerExtName] {
				if !cascade {
					return errDependentObjects("extension", exts[0].Name)
				}
				if ownerExt := c.extByName[ownerExtName]; ownerExt != nil {
					dropExtNames[ownerExtName] = true
					dropExtOIDs[ownerExt.OID] = true
					exts = append(exts, ownerExt)
					changed = true
				}
			}
		}
	}

	inExtDrop := func(oid uint32) bool {
		extName := c.extensionMembers[oid]
		return extName != "" && dropExtNames[extName]
	}
	type memberRel struct {
		schema *Schema
		rel    *Relation
	}
	var rels []memberRel
	relOIDs := map[uint32]bool{}
	type memberIdx struct {
		schema *Schema
		idx    *Index
	}
	var idxs []memberIdx
	idxOIDs := map[uint32]bool{}
	type memberSeq struct {
		schema *Schema
		seq    *Sequence
	}
	var seqs []memberSeq
	seqOIDs := map[uint32]bool{}
	var cascadeIdxs []*Index
	dropIdxOIDs := map[uint32]bool{}
	var extSchemas []*Schema
	extSchemaOIDs := map[uint32]bool{}
	for _, s := range c.schemas {
		if inExtDrop(s.OID) && s.OID != PGCatalogNamespace && s.OID != PGToastNamespace {
			extSchemas = append(extSchemas, s)
			extSchemaOIDs[s.OID] = true
		}
		for _, rel := range s.Relations {
			if inExtDrop(rel.OID) || autoExtOIDs[rel.OID] {
				rels = append(rels, memberRel{schema: s, rel: rel})
				relOIDs[rel.OID] = true
			}
		}
		for _, idx := range s.Indexes {
			if inExtDrop(idx.OID) || autoExtOIDs[idx.OID] {
				idxs = append(idxs, memberIdx{schema: s, idx: idx})
				idxOIDs[idx.OID] = true
			}
		}
		for _, seq := range s.Sequences {
			if inExtDrop(seq.OID) {
				seqs = append(seqs, memberSeq{schema: s, seq: seq})
				seqOIDs[seq.OID] = true
			}
		}
	}
	sort.Slice(rels, func(i, j int) bool { return rels[i].rel.OID < rels[j].rel.OID })
	sort.Slice(idxs, func(i, j int) bool { return idxs[i].idx.OID < idxs[j].idx.OID })
	sort.Slice(seqs, func(i, j int) bool { return seqs[i].seq.OID < seqs[j].seq.OID })
	for changed := true; changed; {
		changed = false
		for oid, child := range c.indexes {
			if child != nil && child.Schema != nil && child.ParentIndexOID != 0 && idxOIDs[child.ParentIndexOID] && !idxOIDs[oid] {
				idxOIDs[oid] = true
				idxs = append(idxs, memberIdx{schema: child.Schema, idx: child})
				changed = true
			}
		}
	}
	sort.Slice(idxs, func(i, j int) bool { return idxs[i].idx.OID < idxs[j].idx.OID })
	var trigs []*Trigger
	trigOIDs := map[uint32]bool{}
	for oid, trig := range c.triggers {
		if trig != nil && (inExtDrop(oid) || autoExtOIDs[oid]) {
			trigs = append(trigs, trig)
			trigOIDs[oid] = true
		}
	}
	for changed := true; changed; {
		changed = false
		for oid, child := range c.triggers {
			if child != nil && child.ParentTriggerOID != 0 && trigOIDs[child.ParentTriggerOID] && !trigOIDs[oid] {
				trigOIDs[oid] = true
				trigs = append(trigs, child)
				changed = true
			}
		}
	}
	sort.Slice(trigs, func(i, j int) bool { return trigs[i].OID < trigs[j].OID })
	for _, x := range idxs {
		if x.idx.ParentIndexOID != 0 && !idxOIDs[x.idx.ParentIndexOID] && !relOIDs[x.idx.RelOID] {
			return &Error{
				Code:    CodeDependentObjects,
				Message: fmt.Sprintf("cannot drop index %q because parent index requires it", x.idx.Name),
			}
		}
	}
	for _, trig := range trigs {
		if trig.ParentTriggerOID != 0 && !trigOIDs[trig.ParentTriggerOID] && !relOIDs[trig.RelOID] {
			return errDependentObjects("trigger", trig.Name)
		}
	}
	var procs []*UserProc
	procOIDs := map[uint32]bool{}
	for _, proc := range c.userProcs {
		if inExtDrop(proc.OID) || autoExtOIDs[proc.OID] {
			procs = append(procs, proc)
			procOIDs[proc.OID] = true
		}
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i].OID < procs[j].OID })
	var types []*BuiltinType
	typeOIDs := map[uint32]bool{}
	seenPrimary := map[uint32]bool{}
	for _, bt := range c.typeByOID {
		if !inExtDrop(bt.OID) || bt.OID < FirstNormalObjectId {
			continue
		}
		typeOIDs[bt.OID] = true
		if bt.Elem != 0 || bt.Type == 'm' || bt.Type == 'c' {
			continue
		}
		if !seenPrimary[bt.OID] {
			seenPrimary[bt.OID] = true
			types = append(types, bt)
			for _, famOID := range c.typeFamilyOIDs(bt) {
				typeOIDs[famOID] = true
			}
		}
	}
	for _, x := range rels {
		if x.rel.RowTypeOID != 0 {
			typeOIDs[x.rel.RowTypeOID] = true
		}
		if x.rel.ArrayOID != 0 {
			typeOIDs[x.rel.ArrayOID] = true
		}
	}
	domConOIDs := map[uint32]bool{}
	for oid := range typeOIDs {
		if dt := c.domainTypes[oid]; dt != nil {
			for _, dc := range dt.Constraints {
				domConOIDs[dc.OID] = true
			}
		}
	}
	// Access methods and operators the extension created can have external
	// users: a table or index built ON a member access method, or a user
	// operator implemented by a member function.
	var extAccessMethodOIDs, extFunctionOIDs []uint32
	for oid, extName := range c.extensionMembers {
		if !dropExtNames[extName] {
			continue
		}
		if c.accessMethods[oid] != nil {
			extAccessMethodOIDs = append(extAccessMethodOIDs, oid)
		}
		if c.userProcs[oid] != nil {
			extFunctionOIDs = append(extFunctionOIDs, oid)
		}
	}
	for oid := range procOIDs {
		if !inExtDrop(oid) {
			extFunctionOIDs = append(extFunctionOIDs, oid)
		}
	}
	var extFamilyOIDs []uint32
	for oid, extName := range c.extensionMembers {
		if dropExtNames[extName] && c.opFamilies[oid] != nil {
			extFamilyOIDs = append(extFamilyOIDs, oid)
		}
	}
	sort.Slice(extFamilyOIDs, func(i, j int) bool { return extFamilyOIDs[i] < extFamilyOIDs[j] })
	sort.Slice(extAccessMethodOIDs, func(i, j int) bool { return extAccessMethodOIDs[i] < extAccessMethodOIDs[j] })
	sort.Slice(extFunctionOIDs, func(i, j int) bool { return extFunctionOIDs[i] < extFunctionOIDs[j] })
	var amRelOIDs, amIdxOIDs []uint32
	var amIdx []*Index
	amRelSeen := map[uint32]bool{}
	amIdxSeen := map[uint32]bool{}
	var opDepOIDs []uint32
	opDepSeen := map[uint32]bool{}
	var extAMUsers []DepEntry
	amUserSeen := map[[2]uint32]bool{}
	for _, famOID := range extFamilyOIDs {
		for _, dep := range c.findNormalDependents('Q', famOID) {
			if dep.ObjType != 'q' || inExtDrop(dep.ObjOID) {
				continue
			}
			key := [2]uint32{uint32(dep.ObjType), dep.ObjOID}
			if !amUserSeen[key] {
				amUserSeen[key] = true
				extAMUsers = append(extAMUsers, dep)
			}
		}
	}
	for _, amOID := range extAccessMethodOIDs {
		for _, dep := range c.findNormalDependents('m', amOID) {
			switch dep.ObjType {
			case 'm', 'q', 'Q':
				if !inExtDrop(dep.ObjOID) {
					key := [2]uint32{uint32(dep.ObjType), dep.ObjOID}
					if !amUserSeen[key] {
						amUserSeen[key] = true
						extAMUsers = append(extAMUsers, dep)
					}
				}
			case 'r':
				if !amRelSeen[dep.ObjOID] {
					amRelSeen[dep.ObjOID] = true
					amRelOIDs = append(amRelOIDs, dep.ObjOID)
				}
			case 'i':
				if idx := c.indexes[dep.ObjOID]; idx != nil && !amIdxSeen[idx.OID] {
					amIdxSeen[idx.OID] = true
					amIdxOIDs = append(amIdxOIDs, idx.OID)
					amIdx = append(amIdx, idx)
				}
			}
		}
	}
	for _, fnOID := range extFunctionOIDs {
		for _, dep := range c.findNormalDependents('f', fnOID) {
			switch dep.ObjType {
			case 'O':
				if !opDepSeen[dep.ObjOID] && !inExtDrop(dep.ObjOID) {
					opDepSeen[dep.ObjOID] = true
					opDepOIDs = append(opDepOIDs, dep.ObjOID)
				}
			case 'm', 'q', 'Q':
				key := [2]uint32{uint32(dep.ObjType), dep.ObjOID}
				if !amUserSeen[key] && !inExtDrop(dep.ObjOID) {
					amUserSeen[key] = true
					extAMUsers = append(extAMUsers, dep)
				}
			}
		}
	}
	sort.Slice(amRelOIDs, func(i, j int) bool { return amRelOIDs[i] < amRelOIDs[j] })
	sort.Slice(amIdxOIDs, func(i, j int) bool { return amIdxOIDs[i] < amIdxOIDs[j] })
	sort.Slice(opDepOIDs, func(i, j int) bool { return opDepOIDs[i] < opDepOIDs[j] })
	sort.Slice(extAMUsers, func(i, j int) bool {
		if extAMUsers[i].ObjType != extAMUsers[j].ObjType {
			return extAMUsers[i].ObjType < extAMUsers[j].ObjType
		}
		return extAMUsers[i].ObjOID < extAMUsers[j].ObjOID
	})
	// A user-created cast that references a member type or function is an
	// external dependent too, though it carries no OID of its own.
	var externalCasts []castKey
	for key, cast := range c.castIndex {
		if owner := c.extensionCasts[key]; owner != "" && dropExtNames[owner] {
			continue // the cast goes with the extension that created it
		}
		if inExtDrop(cast.Source) || typeOIDs[cast.Source] || inExtDrop(cast.Target) || typeOIDs[cast.Target] ||
			(cast.Func != 0 && (inExtDrop(cast.Func) || procOIDs[cast.Func])) {
			externalCasts = append(externalCasts, key)
		}
	}
	sort.Slice(externalCasts, func(i, j int) bool {
		if externalCasts[i].source != externalCasts[j].source {
			return externalCasts[i].source < externalCasts[j].source
		}
		return externalCasts[i].target < externalCasts[j].target
	})

	inDropSet := func(objType byte, objOID uint32) bool {
		if inExtDrop(objOID) {
			return true
		}
		switch objType {
		case 'r':
			return relOIDs[objOID]
		case 'i':
			return idxOIDs[objOID] || (c.indexes[objOID] != nil && relOIDs[c.indexes[objOID].RelOID])
		case 's':
			return seqOIDs[objOID] || (c.sequenceByOID[objOID] != nil && relOIDs[c.sequenceByOID[objOID].OwnerRelOID])
		case 'f':
			return procOIDs[objOID]
		case 't':
			return typeOIDs[objOID]
		case 'c':
			return domConOIDs[objOID] || (c.constraints[objOID] != nil && relOIDs[c.constraints[objOID].RelOID])
		case 'g':
			return trigOIDs[objOID] || (c.triggers[objOID] != nil && relOIDs[c.triggers[objOID].RelOID])
		case 'p':
			return c.policies[objOID] != nil && relOIDs[c.policies[objOID].RelOID]
		case 'R':
			return c.rules[objOID] != nil && relOIDs[c.rules[objOID].OwnerRelOID]
		case 'm':
			return inExtDrop(objOID)
		case 'q':
			return inExtDrop(objOID)
		case 'Q':
			return inExtDrop(objOID)
		case 'O':
			return inExtDrop(objOID)
		case 'n':
			return extSchemaOIDs[objOID]
		}
		return false
	}
	sort.Slice(extSchemas, func(i, j int) bool { return extSchemas[i].OID < extSchemas[j].OID })

	if !cascade {
		extLabel := exts[0].Name
		for _, x := range rels {
			for _, dep := range c.relationDependents(x.rel.OID, 0) {
				if dep.Normal && !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("extension", extLabel)
				}
			}
			if x.rel.RowTypeOID != 0 {
				if bt := c.typeByOID[x.rel.RowTypeOID]; bt != nil {
					for _, depRelOID := range c.typeRelationDependents(bt) {
						if !relOIDs[depRelOID] {
							return errDependentObjects("extension", extLabel)
						}
					}
					for _, dep := range c.typeNormalDependents(bt) {
						if !inDropSet(dep.ObjType, dep.ObjOID) {
							return errDependentObjects("extension", extLabel)
						}
					}
				}
			}
		}
		for _, x := range idxs {
			if x.idx.ConstraintOID != 0 && !inDropSet('c', x.idx.ConstraintOID) {
				return errDependentObjects("extension", extLabel)
			}
			for _, refClass := range []byte{'r', 'i'} {
				for _, dep := range c.findNormalDependents(refClass, x.idx.OID) {
					if !inDropSet(dep.ObjType, dep.ObjOID) {
						return errDependentObjects("extension", extLabel)
					}
				}
			}
		}
		for _, trig := range trigs {
			for _, dep := range c.findNormalDependents('g', trig.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("extension", extLabel)
				}
			}
		}
		for _, x := range seqs {
			for _, dep := range c.findNormalDependents('s', x.seq.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("extension", extLabel)
				}
			}
		}
		for _, proc := range procs {
			for _, dep := range c.findNormalDependents('f', proc.OID) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("extension", extLabel)
				}
			}
		}
		for _, relOID := range amRelOIDs {
			if !relOIDs[relOID] {
				return errDependentObjects("extension", extLabel)
			}
		}
		for _, idxOID := range amIdxOIDs {
			if !idxOIDs[idxOID] {
				return errDependentObjects("extension", extLabel)
			}
		}
		for _, opOID := range opDepOIDs {
			if !inDropSet('O', opOID) {
				return errDependentObjects("extension", extLabel)
			}
		}
		if len(externalCasts) > 0 {
			return errDependentObjects("extension", extLabel)
		}
		if len(extAMUsers) > 0 {
			return errDependentObjects("extension", extLabel)
		}
		for _, bt := range types {
			for _, depRelOID := range c.typeRelationDependents(bt) {
				if !relOIDs[depRelOID] {
					return errDependentObjects("extension", extLabel)
				}
			}
			for _, dep := range c.typeNormalDependents(bt) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
					return errDependentObjects("extension", extLabel)
				}
			}
		}
		for oid, extName := range c.extensionMembers {
			if !dropExtNames[extName] {
				continue
			}
			for _, refClass := range []byte{'o', 'q', 'Q', 'O', 'n'} {
				for _, dep := range c.findNormalDependents(refClass, oid) {
					if !inDropSet(dep.ObjType, dep.ObjOID) {
						return errDependentObjects("extension", extLabel)
					}
				}
			}
		}
		for _, s := range extSchemas {
			for _, rel := range s.Relations {
				if !inDropSet('r', rel.OID) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for _, idx := range s.Indexes {
				if !inDropSet('i', idx.OID) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for _, seq := range s.Sequences {
				if !inDropSet('s', seq.OID) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for _, bt := range c.typeByOID {
				if bt.Namespace == s.OID && bt.OID >= FirstNormalObjectId && bt.Elem == 0 && bt.Type != 'c' && !inDropSet('t', bt.OID) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for _, up := range c.userProcs {
				if up.Schema.OID == s.OID && !inDropSet('f', up.OID) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for _, ext2 := range c.extensions {
				if ext2.SchemaOID == s.OID && !dropExtNames[ext2.Name] {
					return errDependentObjects("extension", extLabel)
				}
			}
			for oid, nsOID := range c.operatorSchemas {
				if nsOID == s.OID && !inDropSet('O', oid) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for oid, opc := range c.opClasses {
				if opc.Schema == s.OID && !inDropSet('q', oid) {
					return errDependentObjects("extension", extLabel)
				}
			}
			for oid, fam := range c.opFamilies {
				if fam.Schema == s.OID && !inDropSet('Q', oid) {
					return errDependentObjects("extension", extLabel)
				}
			}
		}
	}

	// Operator classes belong to the extension, so an index that references
	// one goes with a CASCADE drop too.
	if cascade {
		for oid, extName := range c.extensionMembers {
			if !dropExtNames[extName] {
				continue
			}
			for _, dep := range c.findNormalDependents('o', oid) {
				if dep.ObjType != 'i' || dropIdxOIDs[dep.ObjOID] {
					continue
				}
				idx := c.indexes[dep.ObjOID]
				if idx == nil {
					continue
				}
				cascadeIdxs = append(cascadeIdxs, idx)
				dropIdxOIDs[idx.OID] = true
			}
		}
	}

	for _, x := range rels {
		if c.relationByOID[x.rel.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropRelation(x.rel)
			continue
		}
		c.removeRelation(x.schema, x.rel.Name, x.rel)
	}
	for _, x := range idxs {
		if x.schema.Indexes[x.idx.Name] == x.idx {
			if cascade {
				c.dropDependentObject('i', x.idx.OID)
			} else {
				c.removeIndex(x.schema, x.idx.Name, x.idx)
			}
		}
	}
	if cascade {
		for _, relOID := range amRelOIDs {
			if rel := c.relationByOID[relOID]; rel != nil {
				c.cascadeDropRelation(rel)
			}
		}
		for _, idx := range amIdx {
			if c.indexes[idx.OID] == idx {
				c.dropDependentObject('i', idx.OID)
			}
		}
		for _, opOID := range opDepOIDs {
			c.dropDependents('O', opOID)
			c.removeOperator(opOID)
		}
		for _, key := range externalCasts {
			delete(c.castIndex, key)
			delete(c.extensionCasts, key)
		}
		for _, dep := range extAMUsers {
			switch dep.ObjType {
			case 'm':
				c.cascadeDropAccessMethod(dep.ObjOID)
			case 'q':
				c.cascadeDropOpClass(dep.ObjOID)
			case 'Q':
				c.cascadeDropOpFamily(dep.ObjOID)
			}
		}
		var extOpOIDs []uint32
		for oid, extName := range c.extensionMembers {
			if dropExtNames[extName] && c.operByOID[oid] != nil && oid >= FirstNormalObjectId {
				extOpOIDs = append(extOpOIDs, oid)
			}
		}
		sort.Slice(extOpOIDs, func(i, j int) bool { return extOpOIDs[i] < extOpOIDs[j] })
		for _, opOID := range extOpOIDs {
			c.dropDependents('O', opOID)
		}
	}
	for _, idx := range cascadeIdxs {
		if c.indexes[idx.OID] == idx {
			c.dropDependentObject('i', idx.OID)
		}
	}
	for _, x := range seqs {
		if c.sequenceByOID[x.seq.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('s', x.seq.OID)) > 0 {
			c.dropDependents('s', x.seq.OID)
		}
		if c.sequenceByOID[x.seq.OID] != nil {
			c.removeSequence(x.schema, x.seq)
		}
	}
	for _, proc := range procs {
		if c.userProcs[proc.OID] == nil {
			continue
		}
		if cascade && len(c.findNormalDependents('f', proc.OID)) > 0 {
			c.dropDependents('f', proc.OID)
		}
		if c.userProcs[proc.OID] != nil {
			c.removeFunction(proc.OID, proc.Name)
		}
	}
	for _, bt := range types {
		if c.typeByOID[bt.OID] == nil {
			continue
		}
		if cascade {
			for _, depRelOID := range c.typeRelationDependents(bt) {
				if depRel := c.relationByOID[depRelOID]; depRel != nil {
					c.cascadeDropRelation(depRel)
				}
			}
			c.dropTypeDependents(bt)
		}
		if c.typeByOID[bt.OID] != nil {
			c.deleteUserType(bt)
		}
	}
	if cascade {
		for _, s := range extSchemas {
			c.cascadeDropSchemaContents(s)
		}
	}
	for _, trig := range trigs {
		if c.triggers[trig.OID] == nil {
			continue
		}
		if cascade {
			c.dropDependents('g', trig.OID)
		}
		c.removeTrigger(trig)
	}
	for _, ext := range exts {
		c.removeExtensionMetadata(ext)
	}
	return nil
}

// removeExtensionMetadata unregisters ext and cleans up its extension-member
// tracking, operator classes/families, access methods, comments, and grants.
func (c *Catalog) removeExtensionMetadata(ext *Extension) {
	if ext == nil {
		return
	}
	delete(c.extensions, ext.OID)
	delete(c.extByName, ext.Name)
	c.removeComments('e', ext.OID)
	c.removeDepsOn('e', ext.OID)
	for key, extName := range c.extensionCasts {
		if extName != ext.Name {
			continue
		}
		delete(c.castIndex, key)
		delete(c.extensionCasts, key)
	}
	memberOIDs := map[uint32]bool{}
	for oid, extName := range c.extensionMembers {
		if extName == ext.Name {
			memberOIDs[oid] = true
			delete(c.extensionMembers, oid)
			if evt := c.eventTriggers[oid]; evt != nil {
				c.removeEventTrigger(evt)
			}
			if trig := c.triggers[oid]; trig != nil {
				c.removeTrigger(trig)
			}
			if pol := c.policies[oid]; pol != nil {
				c.removePolicy(pol)
			}
			if con := c.constraints[oid]; con != nil {
				if ownerRel := c.relationByOID[con.RelOID]; ownerRel != nil {
					c.removeConstraint(ownerRel.Schema, con)
				}
			}
			if c.opClasses[oid] != nil {
				c.removeOpClass(oid)
			}
			if c.opFamilies[oid] != nil {
				c.removeOpFamily(oid)
			}
			if c.operByOID[oid] != nil {
				c.removeOperator(oid)
			}
			if c.accessMethods[oid] != nil && oid >= FirstNormalObjectId {
				c.removeAccessMethod(oid)
			}
			if s := c.schemas[oid]; s != nil && oid != PGCatalogNamespace && oid != PGToastNamespace {
				c.deleteSchemaRecord(s)
			}
		}
	}
	if len(memberOIDs) > 0 {
		for key := range c.extensionComments {
			if memberOIDs[key.ObjOID] {
				delete(c.extensionComments, key)
				delete(c.comments, key)
			}
		}
		kept := c.grants[:0]
		for _, g := range c.grants {
			if g.FromExtension && memberOIDs[g.ObjOID] {
				continue
			}
			kept = append(kept, g)
		}
		c.grants = kept
	}
}

// markExtensionMembers records the objects whose OIDs fall in [first, end)
// as belonging to the extension ext, "" when unknown, as pg_depend's 'e'
// entries do. Every object an extension script creates takes its OID while
// the script runs.
func (c *Catalog) markExtensionMembers(first, end uint32, ext string) {
	for oid := first; oid < end; oid++ {
		c.extensionMembers[oid] = ext
	}
}

// isExtensionMember reports whether the object with this OID belongs to an
// extension. Such objects come and go with their extension, so Diff leaves
// them out.
func (c *Catalog) isExtensionMember(oid uint32) bool {
	_, ok := c.extensionMembers[oid]
	return ok
}

// diffsMemberSetting reports whether Diff compares a comment or grant on the
// object with this OID in c, the from catalog when isFrom, against other.
// On an extension's object, everything is compared while both catalogs have
// the extension. When only from has it, the settings go with the extension;
// when only the target has it, those beyond the script's baseline are the
// user's and are compared. On an object of an unknown extension, which a
// snapshot marks skip_dump, there is no baseline and nothing is compared.
func (c *Catalog) diffsMemberSetting(oid uint32, other *Catalog, isFrom, baseline bool) bool {
	ext, ok := c.extensionMembers[oid]
	switch {
	case !ok:
		return true
	case ext == "":
		return false
	case other.extByName[ext] != nil:
		return true
	default:
		return !isFrom && !baseline
	}
}

// RegisterExtensionSQL allows external code to register extension DDL scripts.
// Once registered, CREATE EXTENSION <name> will execute the provided SQL.
func RegisterExtensionSQL(name, sql string) {
	extensionScripts[name] = sql
}

// alterExtensionNamespace moves an extension and its relocatable members to
// another schema.
//
// pg: src/backend/commands/extension.c — AlterExtensionNamespace
func (c *Catalog) alterExtensionNamespace(stmt *nodes.AlterObjectSchemaStmt) error {
	name := extractSimpleObjectName(stmt.Object)
	ext := c.extByName[name]
	if ext == nil {
		if stmt.MissingOk {
			return nil
		}
		return errUndefinedObject("extension", name)
	}
	newSchema := c.schemaByName[stmt.Newschema]
	if newSchema == nil {
		return errUndefinedSchema(stmt.Newschema)
	}
	if !ext.Relocatable {
		return &Error{Code: CodeFeatureNotSupported,
			Message: fmt.Sprintf("extension %q does not support SET SCHEMA", ext.Name)}
	}
	oldOID := ext.SchemaOID
	if oldOID == newSchema.OID {
		return nil
	}
	oldSchema := c.schemas[oldOID]

	// Only this extension's members move.
	var members []uint32
	for oid, extName := range c.extensionMembers {
		if extName == ext.Name {
			members = append(members, oid)
		}
	}
	sort.Slice(members, func(i, j int) bool { return members[i] < members[j] })

	// Preflight the target relation namespace before mutating either schema.
	for _, oid := range members {
		rel := c.relationByOID[oid]
		if rel == nil || oldSchema == nil || rel.Schema != oldSchema {
			continue
		}
		if held := newSchema.Relations[rel.Name]; held != nil && held.OID != rel.OID {
			return errDuplicateTable(rel.Name)
		}
		if newSchema.Indexes[rel.Name] != nil || newSchema.Sequences[rel.Name] != nil {
			return errDuplicateTable(rel.Name)
		}
		if newSchema.Relations[rel.Name] == nil {
			if _, ok := relationNamespaceKind(newSchema, rel.Name); ok {
				return errDuplicateTable(rel.Name)
			}
		}
	}
	for _, oid := range members {
		if idx := c.indexes[oid]; idx != nil && oldSchema != nil && idx.Schema == oldSchema {
			if held := newSchema.Indexes[idx.Name]; held != nil && held.OID != idx.OID {
				return errDuplicateObject("index", idx.Name)
			}
			if _, ok := relationNamespaceKind(newSchema, idx.Name); ok && newSchema.Indexes[idx.Name] == nil {
				return errDuplicateTable(idx.Name)
			}
		}
		if seq := c.sequenceByOID[oid]; seq != nil && oldSchema != nil && seq.Schema == oldSchema {
			if held := newSchema.Sequences[seq.Name]; held != nil && held.OID != seq.OID {
				return errDuplicateObject("sequence", seq.Name)
			}
			if _, ok := relationNamespaceKind(newSchema, seq.Name); ok && newSchema.Sequences[seq.Name] == nil {
				return errDuplicateTable(seq.Name)
			}
		}
		if bt := c.typeByOID[oid]; bt != nil && bt.Namespace == oldOID {
			if held := c.typeByName[typeKey{ns: newSchema.OID, name: bt.TypeName}]; held != nil && held.OID != bt.OID {
				return errDuplicateObject("type", bt.TypeName)
			}
		}
		if up := c.userProcs[oid]; up != nil && oldSchema != nil && up.Schema == oldSchema {
			// pg_proc's uniqueness ignores whether the routine is a function,
			// procedure, or aggregate.
			for _, other := range c.findUserProcsByName(newSchema, up.Name) {
				if other.OID == up.OID {
					continue
				}
				if slices.Equal(c.procByOID[other.OID].ArgTypes, c.procByOID[up.OID].ArgTypes) {
					return errDuplicateFunction(up.Name)
				}
			}
		}
		if opc := c.opClasses[oid]; opc != nil && opc.Schema == oldOID {
			for _, other := range c.opClasses {
				if other.OID == opc.OID || other.Schema != newSchema.OID || other.AMOID != opc.AMOID {
					continue
				}
				if other.Name == opc.Name {
					return errDuplicateObject("operator class", opc.Name)
				}
			}
		}
		if fam := c.opFamilies[oid]; fam != nil && fam.Schema == oldOID {
			for _, other := range c.opFamilies {
				if other.OID == fam.OID || other.Schema != newSchema.OID || other.AMOID != fam.AMOID {
					continue
				}
				if other.Name == fam.Name {
					return errDuplicateObject("operator family", fam.Name)
				}
			}
		}
		if op := c.operByOID[oid]; op != nil && c.operatorSchemas[oid] == oldOID {
			for _, other := range c.operByKey[operKey{name: op.Name, left: op.Left, right: op.Right}] {
				otherNS := uint32(PGCatalogNamespace)
				if other.OID >= FirstNormalObjectId {
					otherNS = c.operatorSchemas[other.OID]
				}
				if other.OID != op.OID && otherNS == newSchema.OID {
					return errDuplicateObject("operator", op.Name)
				}
			}
		}
	}

	for _, oid := range members {
		if rel := c.relationByOID[oid]; rel != nil && oldSchema != nil && rel.Schema == oldSchema {
			delete(oldSchema.Relations, rel.Name)
			newSchema.Relations[rel.Name] = rel
			rel.Schema = newSchema
			for _, con := range c.consByRel[rel.OID] {
				con.Namespace = newSchema.OID
			}
		}
		if idx := c.indexes[oid]; idx != nil && oldSchema != nil && idx.Schema == oldSchema {
			delete(oldSchema.Indexes, idx.Name)
			newSchema.Indexes[idx.Name] = idx
			idx.Schema = newSchema
		}
		if seq := c.sequenceByOID[oid]; seq != nil && oldSchema != nil && seq.Schema == oldSchema {
			delete(oldSchema.Sequences, seq.Name)
			newSchema.Sequences[seq.Name] = seq
			seq.Schema = newSchema
		}
		if bt := c.typeByOID[oid]; bt != nil && bt.Namespace == oldOID {
			delete(c.typeByName, typeKey{ns: oldOID, name: bt.TypeName})
			bt.Namespace = newSchema.OID
			c.typeByName[typeKey{ns: newSchema.OID, name: bt.TypeName}] = bt
			if rt := c.rangeTypes[oid]; rt != nil {
				rt.Namespace = newSchema.OID
			}
		}
		if up := c.userProcs[oid]; up != nil && oldSchema != nil && up.Schema == oldSchema {
			up.Schema = newSchema
		}
		if c.operatorSchemas[oid] == oldOID {
			c.operatorSchemas[oid] = newSchema.OID
		}
		if opc := c.opClasses[oid]; opc != nil && opc.Schema == oldOID {
			opc.Schema = newSchema.OID
		}
		if fam := c.opFamilies[oid]; fam != nil && fam.Schema == oldOID {
			fam.Schema = newSchema.OID
		}
		if evt := c.eventTriggers[oid]; evt != nil {
			_ = evt
		}
	}
	ext.SchemaOID = newSchema.OID
	return nil
}

// ExecAlterObjectDependsStmt handles ALTER ... [NO] DEPENDS ON EXTENSION.
//
// pg: src/backend/commands/alter.c — ExecAlterObjectDependsStmt
func (c *Catalog) ExecAlterObjectDependsStmt(stmt *nodes.AlterObjectDependsStmt) error {
	var objType byte
	var objOID uint32

	switch stmt.ObjectType {
	case nodes.OBJECT_MATVIEW:
		var schemaName, relName string
		if stmt.Relation != nil {
			schemaName = stmt.Relation.Schemaname
			relName = stmt.Relation.Relname
		}
		_, rel, err := c.findRelation(schemaName, relName)
		if err != nil {
			return err
		}
		if rel.RelKind != 'm' {
			return errWrongObjectType(relName, "a materialized view")
		}
		objType, objOID = 'r', rel.OID

	case nodes.OBJECT_INDEX:
		var schemaName, relName string
		if stmt.Relation != nil {
			schemaName = stmt.Relation.Schemaname
			relName = stmt.Relation.Relname
		}
		_, idx, err := c.findIndexRelation(schemaName, relName)
		if err != nil {
			return err
		}
		objType, objOID = 'i', idx.OID

	case nodes.OBJECT_TRIGGER:
		var schemaName, relName string
		if stmt.Relation != nil {
			schemaName = stmt.Relation.Schemaname
			relName = stmt.Relation.Relname
		}
		rel, err := c.findTriggerTargetRelation(schemaName, relName)
		if err != nil {
			return err
		}
		_, trigName := extractObjectName(stmt.Object)
		var found *Trigger
		if rel != nil {
			for _, trig := range c.triggersByRel[rel.OID] {
				if trig.Name == trigName {
					found = trig
					break
				}
			}
		}
		if found == nil {
			return errUndefinedTrigger(trigName, relName)
		}
		objType, objOID = 'g', found.OID

	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE:
		owa, ok := stmt.Object.(*nodes.ObjectWithArgs)
		if !ok {
			return errInvalidParameterValue("ALTER ... DEPENDS ON EXTENSION requires a routine specification")
		}
		bp, err := c.lookupFuncWithArgs(stmt.ObjectType, owa)
		if err != nil {
			return err
		}
		objType, objOID = 'f', bp.OID

	default:
		return errInvalidParameterValue(fmt.Sprintf("unsupported AlterObjectDependsStmt object type: %d", stmt.ObjectType))
	}

	_, extName := extractObjectName(stmt.Extname)
	ext := c.extByName[extName]
	if ext == nil {
		return errUndefinedObject("extension", extName)
	}

	if stmt.Remove {
		n := 0
		for _, d := range c.deps {
			if d.ObjType == objType && d.ObjOID == objOID && d.ObjSubID == 0 &&
				d.RefType == 'e' && d.RefOID == ext.OID && d.DepType == DepAutoExtension {
				continue
			}
			c.deps[n] = d
			n++
		}
		c.deps = c.deps[:n]
		return nil
	}

	for _, d := range c.deps {
		if d.ObjType == objType && d.ObjOID == objOID && d.ObjSubID == 0 &&
			d.RefType == 'e' && d.RefOID == ext.OID && d.DepType == DepAutoExtension {
			return nil
		}
	}
	c.recordDependency(objType, objOID, 0, 'e', ext.OID, 0, DepAutoExtension)
	return nil
}

// findTriggerTargetRelation implements PostgreSQL's table_openrv check in
// get_object_address_relobject: missing schema is 3F000, missing relation is
// 42P01, indexes and composite types are 42809, and sequences pass table_open
// with nil *Relation so get_trigger_oid reports 42704.
//
// pg: src/backend/catalog/objectaddress.c — get_object_address_relobject
// pg: src/backend/access/table/table.c — validate_relation_kind
func (c *Catalog) findTriggerTargetRelation(schemaName, relName string) (*Relation, error) {
	checkInSchema := func(s *Schema) (*Relation, bool, error) {
		if s.Indexes[relName] != nil {
			return nil, true, errWrongObjectType(relName, "a table")
		}
		if r := s.Relations[relName]; r != nil {
			if r.RelKind == 'c' {
				return nil, true, errWrongObjectType(relName, "a table")
			}
			return r, true, nil
		}
		if s.Sequences[relName] != nil {
			return nil, true, nil
		}
		return nil, false, nil
	}

	if schemaName != "" {
		s := c.schemaByName[schemaName]
		if s == nil {
			return nil, errUndefinedSchema(schemaName)
		}
		rel, found, err := checkInSchema(s)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, errUndefinedTable(relName)
		}
		return rel, nil
	}

	for _, nsOID := range c.searchPathWithCatalog() {
		s := c.schemas[nsOID]
		if s == nil {
			continue
		}
		rel, found, err := checkInSchema(s)
		if err != nil {
			return nil, err
		}
		if found {
			return rel, nil
		}
	}
	return nil, errUndefinedTable(relName)
}

// lookupFuncWithArgs resolves an ObjectWithArgs across user-defined and
// built-in routines using PostgreSQL's LookupFuncWithArgs rules.
//
// pg: src/backend/parser/parse_func.c — LookupFuncWithArgs
func (c *Catalog) lookupFuncWithArgs(objType nodes.ObjectType, owa *nodes.ObjectWithArgs) (*BuiltinProc, error) {
	kindStr := "function"
	switch objType {
	case nodes.OBJECT_PROCEDURE:
		kindStr = "procedure"
	case nodes.OBJECT_ROUTINE:
		kindStr = "routine"
	case nodes.OBJECT_AGGREGATE:
		kindStr = "aggregate"
	}

	schemaName, funcName := qualifiedName(owa.Objname)
	argOIDs, err := c.routineArgOIDs(owa)
	if err != nil {
		return nil, err
	}

	var nsOIDs []uint32
	if schemaName != "" {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return nil, err
		}
		nsOIDs = []uint32{schema.OID}
	} else {
		nsOIDs = c.routineSearchPathWithCatalog()
	}

	if argOIDs != nil {
		var bp *BuiltinProc
		for _, nsOID := range nsOIDs {
			for _, p := range c.procByName[funcName] {
				if c.procInNamespace(p, nsOID) && slices.Equal(p.ArgTypes, argOIDs) {
					bp = p
					break
				}
			}
			if bp != nil {
				break
			}
		}
		if bp == nil {
			return nil, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
		}
		switch objType {
		case nodes.OBJECT_PROCEDURE:
			if bp.Kind != 'p' {
				return nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a procedure", funcName)}
			}
		case nodes.OBJECT_FUNCTION:
			if bp.Kind == 'p' {
				return nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a function", funcName)}
			}
		case nodes.OBJECT_AGGREGATE:
			if bp.Kind != PROKIND_AGGREGATE {
				return nil, &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not an aggregate", funcName)}
			}
		}
		return bp, nil
	}

	// ArgsUnspecified: FuncnameGetCandidates collects visible candidates across
	// nsOIDs, masking identical argument-type signatures from later schemas, and
	// LookupFuncWithArgs filters by objType.
	var visible []*BuiltinProc
	for _, nsOID := range nsOIDs {
		var inNS []*BuiltinProc
		for _, p := range c.procByName[funcName] {
			if !c.procInNamespace(p, nsOID) {
				continue
			}
			masked := false
			for _, prev := range visible {
				if slices.Equal(prev.ArgTypes, p.ArgTypes) {
					masked = true
					break
				}
			}
			if !masked {
				inNS = append(inNS, p)
			}
		}
		visible = append(visible, inNS...)
	}

	var matched []*BuiltinProc
	for _, p := range visible {
		switch objType {
		case nodes.OBJECT_PROCEDURE:
			if p.Kind == 'p' {
				matched = append(matched, p)
			}
		case nodes.OBJECT_FUNCTION:
			if p.Kind != 'p' {
				matched = append(matched, p)
			}
		case nodes.OBJECT_AGGREGATE:
			if p.Kind == PROKIND_AGGREGATE {
				matched = append(matched, p)
			}
		default:
			matched = append(matched, p)
		}
	}

	if len(matched) == 0 {
		return nil, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("could not find a %s named %q", kindStr, funcName)}
	}
	if len(matched) > 1 {
		return nil, &Error{Code: CodeAmbiguousFunction, Message: fmt.Sprintf("%s name %q is not unique", kindStr, funcName)}
	}
	return matched[0], nil
}

func (c *Catalog) procInNamespace(p *BuiltinProc, nsOID uint32) bool {
	if p == nil {
		return false
	}
	if up := c.userProcs[p.OID]; up != nil {
		return up.Schema != nil && up.Schema.OID == nsOID
	}
	return nsOID == PGCatalogNamespace
}

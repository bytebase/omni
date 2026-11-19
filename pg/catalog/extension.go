package catalog

import (
	"fmt"
	"maps"
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
	if c.extByName[stmt.Extname] != nil {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip,
				fmt.Sprintf("extension %q already exists, skipping", stmt.Extname))
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
	// Default: first in search path (typically "public").
	targetSchema, err := c.resolveTargetSchema(schemaName)
	if err != nil {
		return err
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
				c.searchPath = origSearchPath
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
		Relocatable: true,
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
			exts = append(exts, ext)
		}
	}
	if len(exts) == 0 {
		return nil
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
	for _, s := range c.schemas {
		for _, rel := range s.Relations {
			if inExtDrop(rel.OID) {
				rels = append(rels, memberRel{schema: s, rel: rel})
				relOIDs[rel.OID] = true
			}
		}
		for _, idx := range s.Indexes {
			if inExtDrop(idx.OID) {
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
	var procs []*UserProc
	procOIDs := map[uint32]bool{}
	for _, proc := range c.userProcs {
		if inExtDrop(proc.OID) {
			procs = append(procs, proc)
			procOIDs[proc.OID] = true
		}
	}
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
	sort.Slice(extAccessMethodOIDs, func(i, j int) bool { return extAccessMethodOIDs[i] < extAccessMethodOIDs[j] })
	sort.Slice(extFunctionOIDs, func(i, j int) bool { return extFunctionOIDs[i] < extFunctionOIDs[j] })
	var amRelOIDs, amIdxOIDs []uint32
	var amIdx []*Index
	amRelSeen := map[uint32]bool{}
	amIdxSeen := map[uint32]bool{}
	var opDepOIDs []uint32
	opDepSeen := map[uint32]bool{}
	for _, amOID := range extAccessMethodOIDs {
		for _, dep := range c.findNormalDependents('m', amOID) {
			switch dep.ObjType {
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
			if dep.ObjType != 'O' || opDepSeen[dep.ObjOID] || inExtDrop(dep.ObjOID) {
				continue
			}
			opDepSeen[dep.ObjOID] = true
			opDepOIDs = append(opDepOIDs, dep.ObjOID)
		}
	}
	sort.Slice(amRelOIDs, func(i, j int) bool { return amRelOIDs[i] < amRelOIDs[j] })
	sort.Slice(amIdxOIDs, func(i, j int) bool { return amIdxOIDs[i] < amIdxOIDs[j] })
	sort.Slice(opDepOIDs, func(i, j int) bool { return opDepOIDs[i] < opDepOIDs[j] })

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
			return c.triggers[objOID] != nil && relOIDs[c.triggers[objOID].RelOID]
		case 'p':
			return c.policies[objOID] != nil && relOIDs[c.policies[objOID].RelOID]
		case 'm':
			return inExtDrop(objOID)
		case 'O':
			return inExtDrop(objOID)
		}
		return false
	}

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
			for _, dep := range c.findNormalDependents('o', oid) {
				if !inDropSet(dep.ObjType, dep.ObjOID) {
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
			c.removeIndex(x.schema, x.idx.Name, x.idx)
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
				c.removeIndex(idx.Schema, idx.Name, idx)
			}
		}
		for _, opOID := range opDepOIDs {
			c.removeOperator(opOID)
		}
	}
	for _, idx := range cascadeIdxs {
		if c.indexes[idx.OID] == idx {
			c.removeIndex(idx.Schema, idx.Name, idx)
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
			if opc := c.opClasses[oid]; opc != nil {
				if opc.IsDefault {
					delete(c.opClassByKey, opClassKey{amOID: opc.AMOID, typeOID: opc.TypeOID})
				}
				delete(c.opClasses, oid)
			}
			delete(c.opFamilies, oid)
			if c.operByOID[oid] != nil {
				c.removeOperator(oid)
			}
			if am := c.accessMethods[oid]; am != nil && oid >= FirstNormalObjectId {
				delete(c.accessMethods, oid)
				delete(c.accessMethodByName, am.Name)
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

package catalog

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"

	metadatapb "github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// ApplyDiffToMetadata applies a SchemaDiff (from user DDL execution) to a
// clone of the original metadata proto, producing the post-DDL metadata.
// The original proto is NOT modified.
func ApplyDiffToMetadata(original *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, diff *SchemaDiff) *metadatapb.DatabaseSchemaMetadata {
	if (diff == nil || diff.IsEmpty()) && !hasRelationOwnerChange(catBefore, catAfter) && !hasSchemaOwnerChange(catBefore, catAfter) && !hasDatabaseOwnerChange(catBefore, catAfter) && !hasExtensionMemberChange(catBefore, catAfter) {
		return proto.CloneOf(original)
	}
	if diff == nil {
		diff = &SchemaDiff{}
	}

	result, ok := proto.Clone(original).(*metadatapb.DatabaseSchemaMetadata)
	if !ok {
		return proto.CloneOf(original)
	}

	renamedSchemas := applySchemaDiffs(result, diff.Schemas)
	canonicalSchema := func(name string) string {
		if renamed, ok := renamedSchemas[name]; ok {
			return renamed
		}
		return name
	}

	renamed := renamePartitions(result, catBefore, catAfter, diff.Relations, canonicalSchema)
	for _, rel := range diff.Relations {
		if (rel.From != nil && renamed[rel.From.OID]) || (rel.To != nil && renamed[rel.To.OID]) {
			continue
		}
		schemaName := canonicalSchema(rel.SchemaName)
		if isTempRelation(rel.SchemaName, rel.From) || isTempRelation(rel.SchemaName, rel.To) || isTempSchemaName(schemaName) {
			continue
		}
		if rel.Action == DiffDrop {
			dropRelation(result, schemaName, rel)
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch rel.Action {
		case DiffAdd:
			addRelation(result, sm, catAfter, rel)
		case DiffModify:
			modifyRelation(result, sm, catBefore, catAfter, rel)
		default:
		}
	}

	applySequenceDiffs(result, catAfter, diff.Sequences, renamedSchemas, canonicalSchema)

	for _, e := range diff.Enums {
		schemaName := canonicalSchema(e.SchemaName)
		if e.Action == DiffDrop {
			if renamedSchemas[e.SchemaName] != "" && hasMatchingEnumAdd(diff.Enums, renamedSchemas[e.SchemaName], e.Name) {
				continue
			}
			if sm := findSchema(result, schemaName); sm != nil {
				sm.EnumTypes = removeEnumByName(sm.EnumTypes, e.Name)
			}
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch e.Action {
		case DiffAdd:
			upsertEnum(sm, e.Name, e.ToValues)
		case DiffModify:
			upsertEnum(sm, e.Name, e.ToValues)
		default:
		}
	}

	for _, ct := range diff.CompositeTypes {
		schemaName := canonicalSchema(ct.SchemaName)
		if ct.Action == DiffDrop {
			if renamedSchemas[ct.SchemaName] != "" && hasMatchingCompositeAdd(diff.CompositeTypes, renamedSchemas[ct.SchemaName], ct.Name) {
				continue
			}
			if sm := findSchema(result, schemaName); sm != nil {
				sm.CompositeTypes = removeCompositeTypeByName(sm.CompositeTypes, ct.Name)
			}
			continue
		}
		if ct.To == nil {
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch ct.Action {
		case DiffAdd, DiffModify:
			idx := -1
			for i, existing := range sm.CompositeTypes {
				if existing.Name == ct.Name {
					idx = i
					break
				}
			}
			if idx >= 0 {
				sm.CompositeTypes[idx] = compositeTypeToProto(catAfter, ct.Name, ct.From, ct.To, sm.CompositeTypes[idx])
			} else {
				sm.CompositeTypes = append(sm.CompositeTypes, compositeTypeToProto(catAfter, ct.Name, nil, ct.To, nil))
			}
		default:
		}
	}

	applyRangeDiffs(result, catAfter, diff.Ranges, canonicalSchema)

	applyCollationDiffs(result, catAfter, diff.Collations, renamedSchemas, canonicalSchema)

	applyFunctionDiffs(result, catAfter, diff.Functions, canonicalSchema)

	applyEventTriggerDiffs(result, catBefore, catAfter, diff.EventTriggers)
	applyForeignDataWrapperDiffs(result, diff.ForeignDataWrappers)
	applyForeignServerDiffs(result, diff.ForeignServers)
	applyExtensionDiffs(result, catBefore, catAfter, diff.Extensions, canonicalSchema)
	applyExtensionMemberChanges(result, original, catBefore, catAfter, canonicalSchema)
	applyForeignTableDiffs(result, catBefore, catAfter, diff.ForeignTables, canonicalSchema)
	applyCommentDiffs(result, catBefore, catAfter, diff.Comments)
	applyRelationOwnerChanges(result, catBefore, catAfter)
	applySchemaOwnerChanges(result, catBefore, catAfter)
	if hasDatabaseOwnerChange(catBefore, catAfter) {
		result.Owner = effectiveDatabaseOwner(catAfter)
	}

	return result
}

func effectiveDatabaseOwner(cat *Catalog) string {
	if cat == nil {
		return ""
	}
	if cat.databaseName != "" {
		return cat.databaseOwners[cat.databaseName]
	}
	return cat.databaseOwners["\x00current"]
}

func hasDatabaseOwnerChange(catBefore, catAfter *Catalog) bool {
	if catBefore == nil || catAfter == nil {
		return false
	}
	return effectiveDatabaseOwner(catBefore) != effectiveDatabaseOwner(catAfter)
}

func applyForeignDataWrapperDiffs(result *metadatapb.DatabaseSchemaMetadata, entries []ForeignDataWrapperDiffEntry) {
	if result == nil || len(entries) == 0 {
		return
	}
	dropped := make(map[uint32]ForeignDataWrapperDiffEntry)
	added := make(map[uint32]ForeignDataWrapperDiffEntry)
	for _, e := range entries {
		if e.Action == DiffDrop && e.From != nil {
			dropped[e.From.OID] = e
		} else if e.Action == DiffAdd && e.To != nil {
			added[e.To.OID] = e
		}
	}
	toProto := func(fdw *ForeignDataWrapper) *metadatapb.ForeignDataWrapperMetadata {
		return &metadatapb.ForeignDataWrapperMetadata{
			Name:      fdw.Name,
			Handler:   fdw.Handler,
			Validator: fdw.Validator,
			Options:   slices.Clone(fdw.Options),
		}
	}
	pop := func(name string) {
		result.ForeignDataWrappers = slices.DeleteFunc(result.ForeignDataWrappers, func(f *metadatapb.ForeignDataWrapperMetadata) bool {
			return f.Name == name
		})
	}
	upsert := func(updated *metadatapb.ForeignDataWrapperMetadata) {
		for i, existing := range result.ForeignDataWrappers {
			if existing.Name == updated.Name {
				result.ForeignDataWrappers[i] = updated
				return
			}
		}
		result.ForeignDataWrappers = append(result.ForeignDataWrappers, updated)
	}
	handledDrop := make(map[string]bool)
	handledAdd := make(map[string]bool)
	for oid, dropEntry := range dropped {
		addEntry, ok := added[oid]
		if !ok {
			continue
		}
		handledDrop[dropEntry.Name] = true
		handledAdd[addEntry.Name] = true
		pop(dropEntry.Name)
		upsert(toProto(addEntry.To))
	}
	for _, e := range entries {
		switch e.Action {
		case DiffDrop:
			if handledDrop[e.Name] {
				continue
			}
			pop(e.Name)
		case DiffAdd:
			if handledAdd[e.Name] || e.To == nil {
				continue
			}
			upsert(toProto(e.To))
		case DiffModify:
			if e.To == nil {
				continue
			}
			upsert(toProto(e.To))
		}
	}
}

func applyForeignServerDiffs(result *metadatapb.DatabaseSchemaMetadata, entries []ForeignServerDiffEntry) {
	if result == nil || len(entries) == 0 {
		return
	}
	dropped := make(map[uint32]ForeignServerDiffEntry)
	added := make(map[uint32]ForeignServerDiffEntry)
	for _, e := range entries {
		if e.Action == DiffDrop && e.From != nil {
			dropped[e.From.OID] = e
		} else if e.Action == DiffAdd && e.To != nil {
			added[e.To.OID] = e
		}
	}
	toProto := func(srv *ForeignServer) *metadatapb.ForeignServerMetadata {
		return &metadatapb.ForeignServerMetadata{
			Name:    srv.Name,
			Fdw:     srv.FdwName,
			Options: slices.Clone(srv.Options),
			Type:    cloneStringPtr(srv.Type),
			Version: cloneStringPtr(srv.Version),
		}
	}
	pop := func(name string) {
		result.ForeignServers = slices.DeleteFunc(result.ForeignServers, func(s *metadatapb.ForeignServerMetadata) bool {
			return s.Name == name
		})
	}
	upsert := func(updated *metadatapb.ForeignServerMetadata) {
		for i, existing := range result.ForeignServers {
			if existing.Name == updated.Name {
				result.ForeignServers[i] = updated
				return
			}
		}
		result.ForeignServers = append(result.ForeignServers, updated)
	}
	handledDrop := make(map[string]bool)
	handledAdd := make(map[string]bool)
	for oid, dropEntry := range dropped {
		addEntry, ok := added[oid]
		if !ok {
			continue
		}
		handledDrop[dropEntry.Name] = true
		handledAdd[addEntry.Name] = true
		pop(dropEntry.Name)
		upsert(toProto(addEntry.To))
	}
	for _, e := range entries {
		switch e.Action {
		case DiffDrop:
			if handledDrop[e.Name] {
				continue
			}
			pop(e.Name)
		case DiffAdd:
			if handledAdd[e.Name] || e.To == nil {
				continue
			}
			upsert(toProto(e.To))
		case DiffModify:
			if e.To == nil {
				continue
			}
			upsert(toProto(e.To))
		}
	}
}

func hasSchemaOwnerChange(catBefore, catAfter *Catalog) bool {
	if catBefore == nil || catAfter == nil {
		return false
	}
	for _, s := range catAfter.UserSchemas() {
		if isTempSchemaName(s.Name) {
			continue
		}
		prev := catBefore.schemas[s.OID]
		if prev == nil {
			prev = catBefore.schemaByName[s.Name]
		}
		if prev != nil && prev.effectiveOwner() != s.effectiveOwner() {
			return true
		}
	}
	return false
}

func applySchemaOwnerChanges(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog) {
	if catBefore == nil || catAfter == nil {
		return
	}
	for _, s := range catAfter.UserSchemas() {
		if isTempSchemaName(s.Name) {
			continue
		}
		prev := catBefore.schemas[s.OID]
		if prev == nil {
			prev = catBefore.schemaByName[s.Name]
		}
		if prev == nil || prev.effectiveOwner() == s.effectiveOwner() {
			continue
		}
		if sm := findSchema(result, s.Name); sm != nil {
			sm.Owner = s.effectiveOwner()
		}
	}
}

func hasRelationOwnerChange(catBefore, catAfter *Catalog) bool {
	if catBefore == nil || catAfter == nil {
		return false
	}
	for _, s := range catAfter.UserSchemas() {
		if isTempSchemaName(s.Name) {
			continue
		}
		for _, rel := range s.Relations {
			if rel.Persistence == 't' || (rel.RelKind != 'r' && rel.RelKind != 'p') {
				continue
			}
			prev := catBefore.GetRelationByOID(rel.OID)
			if prev == nil {
				prev = catBefore.GetRelation(s.Name, rel.Name)
			}
			if prev != nil && prev.OwnerName != rel.OwnerName {
				return true
			}
		}
	}
	return false
}

func applyRelationOwnerChanges(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog) {
	if catBefore == nil || catAfter == nil {
		return
	}
	for _, s := range catAfter.UserSchemas() {
		if isTempSchemaName(s.Name) {
			continue
		}
		for _, rel := range s.Relations {
			if rel.Persistence == 't' || (rel.RelKind != 'r' && rel.RelKind != 'p') {
				continue
			}
			prev := catBefore.GetRelationByOID(rel.OID)
			if prev == nil {
				prev = catBefore.GetRelation(s.Name, rel.Name)
			}
			if prev == nil || prev.OwnerName == rel.OwnerName {
				continue
			}
			if sm := findSchema(result, s.Name); sm != nil {
				if tbl := findTable(sm, rel.Name); tbl != nil {
					tbl.Owner = rel.OwnerName
				}
			}
		}
	}
}

func applyFunctionDiffs(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, functions []FunctionDiffEntry, canonicalSchema func(string) string) {
	droppedByOID := make(map[uint32]FunctionDiffEntry)
	addedByOID := make(map[uint32]FunctionDiffEntry)
	for _, f := range functions {
		if isTempSchemaName(f.SchemaName) || isTempSchemaName(canonicalSchema(f.SchemaName)) {
			continue
		}
		if f.Action == DiffDrop && f.From != nil && f.From.OID != 0 {
			droppedByOID[f.From.OID] = f
		} else if f.Action == DiffAdd && f.To != nil && f.To.OID != 0 {
			addedByOID[f.To.OID] = f
		}
	}
	handledOIDs := make(map[uint32]bool)
	for oid, dropEntry := range droppedByOID {
		addEntry, ok := addedByOID[oid]
		if !ok {
			continue
		}
		handledOIDs[oid] = true
		srcSchema := canonicalSchema(dropEntry.SchemaName)
		dstSchema := canonicalSchema(addEntry.SchemaName)
		srcSM := findSchema(result, srcSchema)
		dstSM := findOrCreateSchema(result, dstSchema)
		upsertRoutineFromDiff(srcSM, dstSM, catAfter, dropEntry.SchemaName, dropEntry.Identity, addEntry.Identity, dropEntry.From, addEntry.To)
	}
	for _, f := range functions {
		if isTempSchemaName(f.SchemaName) || isTempSchemaName(canonicalSchema(f.SchemaName)) {
			continue
		}
		if (f.From != nil && f.From.OID != 0 && handledOIDs[f.From.OID]) ||
			(f.To != nil && f.To.OID != 0 && handledOIDs[f.To.OID]) {
			continue
		}
		schemaName := canonicalSchema(f.SchemaName)
		switch f.Action {
		case DiffDrop:
			if sm := findSchema(result, schemaName); sm != nil {
				if f.From != nil && f.From.Kind == 'p' {
					_, sm.Procedures = popProcedureMatchingIdentity(sm.Procedures, f.SchemaName, f.Identity)
				} else {
					_, sm.Functions = popFunctionMatchingIdentity(sm.Functions, f.SchemaName, f.Identity)
				}
			}
		case DiffAdd:
			if f.To != nil {
				sm := findOrCreateSchema(result, schemaName)
				upsertRoutineFromDiff(nil, sm, catAfter, f.SchemaName, "", f.Identity, nil, f.To)
			}
		case DiffModify:
			if f.To != nil {
				sm := findOrCreateSchema(result, schemaName)
				upsertRoutineFromDiff(sm, sm, catAfter, f.SchemaName, f.Identity, f.Identity, f.From, f.To)
			}
		default:
		}
	}
}

func upsertRoutineFromDiff(srcSM, dstSM *metadatapb.SchemaMetadata, catAfter *Catalog, srcSchema, oldIdentity, newIdentity string, from, to *UserProc) {
	if to == nil || dstSM == nil {
		return
	}
	if to.Kind == 'p' {
		var prev *metadatapb.ProcedureMetadata
		if srcSM != nil && oldIdentity != "" {
			prev, srcSM.Procedures = popProcedureMatchingIdentity(srcSM.Procedures, srcSchema, oldIdentity)
		}
		p := procedureToProto(catAfter, to, newIdentity)
		if prev != nil {
			p.Comment = prev.Comment
			p.SkipDump = prev.SkipDump
		}
		dstSM.Procedures = append(dstSM.Procedures, p)
		return
	}
	var prev *metadatapb.FunctionMetadata
	if srcSM != nil && oldIdentity != "" {
		prev, srcSM.Functions = popFunctionMatchingIdentity(srcSM.Functions, srcSchema, oldIdentity)
	}
	fn := functionToProto(catAfter, to, newIdentity)
	bodyChanged := from == nil || from.Body != to.Body || from.Language != to.Language || from.HasSqlBody != to.HasSqlBody
	if bodyChanged {
		fn.DependencyTables = extractFunctionDependencyTables(catAfter, to)
	} else if prev != nil {
		fn.DependencyTables = prev.DependencyTables
	}
	if prev != nil {
		fn.Comment = prev.Comment
		fn.SkipDump = prev.SkipDump
	}
	dstSM.Functions = append(dstSM.Functions, fn)
}

func applySchemaDiffs(result *metadatapb.DatabaseSchemaMetadata, schemas []SchemaDiffEntry) map[string]string {
	dropped := make(map[uint32]SchemaDiffEntry)
	added := make(map[uint32]SchemaDiffEntry)
	for _, s := range schemas {
		if s.Action == DiffDrop && s.From != nil {
			dropped[s.From.OID] = s
		} else if s.Action == DiffAdd && s.To != nil {
			added[s.To.OID] = s
		}
	}

	renamed := make(map[string]string)
	handledDrop := make(map[string]bool)
	handledAdd := make(map[string]bool)
	for oid, dropEntry := range dropped {
		addEntry, ok := added[oid]
		if !ok {
			continue
		}
		oldName := dropEntry.Name
		newName := addEntry.Name
		renamed[oldName] = newName
		handledDrop[oldName] = true
		handledAdd[newName] = true
		if sm := findSchema(result, oldName); sm != nil {
			if owner := addEntry.To.effectiveOwner(); owner != "" && (dropEntry.From == nil || dropEntry.From.effectiveOwner() != owner) {
				sm.Owner = owner
			}
			sm.Name = newName
		} else {
			sm := findOrCreateSchema(result, newName)
			if owner := addEntry.To.effectiveOwner(); owner != "" {
				sm.Owner = owner
			}
		}
		rewritePartitionSchemaReferences(result, oldName, newName)
	}

	for _, s := range schemas {
		switch s.Action {
		case DiffAdd:
			if handledAdd[s.Name] {
				continue
			}
			sm := findOrCreateSchema(result, s.Name)
			if s.To != nil {
				if owner := s.To.effectiveOwner(); owner != "" {
					sm.Owner = owner
				}
			}
		case DiffDrop:
			if handledDrop[s.Name] {
				continue
			}
			result.Schemas = removeSchema(result.Schemas, s.Name)
		case DiffModify:
			if s.To != nil {
				sm := findOrCreateSchema(result, s.Name)
				sm.Owner = s.To.Owner
			}
		default:
		}
	}
	return renamed
}

func rewritePartitionSchemaReferences(meta *metadatapb.DatabaseSchemaMetadata, oldName, newName string) {
	var walk func(partitions []*metadatapb.TablePartitionMetadata)
	walk = func(partitions []*metadatapb.TablePartitionMetadata) {
		for _, p := range partitions {
			if p.Schema == oldName {
				p.Schema = newName
			}
			walk(p.Subpartitions)
		}
	}
	for _, s := range meta.Schemas {
		for _, tbl := range s.Tables {
			walk(tbl.Partitions)
		}
	}
	for _, ext := range meta.Extensions {
		if ext.Schema == oldName {
			ext.Schema = newName
		}
	}
	rewriteDependencyReferences(meta, oldName, "", newName, "")
}

func rewriteTableDependencyReferences(meta *metadatapb.DatabaseSchemaMetadata, oldSchema, oldTable, newSchema, newTable string) {
	if oldSchema == newSchema && oldTable == newTable {
		return
	}
	rewriteDependencyReferences(meta, oldSchema, oldTable, newSchema, newTable)
}

func rewriteDependencyReferences(meta *metadatapb.DatabaseSchemaMetadata, oldSchema, oldTable, newSchema, newTable string) {
	for _, s := range meta.Schemas {
		for _, v := range s.Views {
			for _, dep := range v.DependencyColumns {
				if dep.Schema == oldSchema && (oldTable == "" || dep.Table == oldTable) {
					dep.Schema = newSchema
					if oldTable != "" {
						dep.Table = newTable
					}
				}
			}
		}
		for _, mv := range s.MaterializedViews {
			for _, dep := range mv.DependencyColumns {
				if dep.Schema == oldSchema && (oldTable == "" || dep.Table == oldTable) {
					dep.Schema = newSchema
					if oldTable != "" {
						dep.Table = newTable
					}
				}
			}
		}
		for _, fn := range s.Functions {
			for _, dep := range fn.DependencyTables {
				if dep.Schema == oldSchema && (oldTable == "" || dep.Table == oldTable) {
					dep.Schema = newSchema
					if oldTable != "" {
						dep.Table = newTable
					}
				}
			}
		}
	}
}

func rewriteColumnDependencyReferences(meta *metadatapb.DatabaseSchemaMetadata, schema, table, oldCol, newCol string) {
	if meta == nil || oldCol == newCol {
		return
	}
	for _, s := range meta.Schemas {
		for _, v := range s.Views {
			for _, dep := range v.DependencyColumns {
				if dep.Schema == schema && dep.Table == table && dep.Column == oldCol {
					dep.Column = newCol
				}
			}
		}
		for _, mv := range s.MaterializedViews {
			for _, dep := range mv.DependencyColumns {
				if dep.Schema == schema && dep.Table == table && dep.Column == oldCol {
					dep.Column = newCol
				}
			}
		}
	}
	if sm := findSchema(meta, schema); sm != nil {
		for _, seq := range sm.Sequences {
			if seq.OwnerTable == table && seq.OwnerColumn == oldCol {
				seq.OwnerColumn = newCol
			}
		}
	}
}

func rewriteColumnRenamesForRelation(meta *metadatapb.DatabaseSchemaMetadata, schema, table string, fromRel, toRel *Relation) {
	if meta == nil || fromRel == nil || toRel == nil {
		return
	}
	fromByAttNum := make(map[int16]string, len(fromRel.Columns))
	for _, col := range fromRel.Columns {
		if col != nil && col.AttNum > 0 {
			fromByAttNum[col.AttNum] = col.Name
		}
	}
	for _, col := range toRel.Columns {
		if col == nil || col.AttNum <= 0 {
			continue
		}
		if oldName, ok := fromByAttNum[col.AttNum]; ok && oldName != col.Name {
			rewriteColumnDependencyReferences(meta, schema, table, oldName, col.Name)
		}
	}
}

func rewritePartitionSchemasOnParentMove(partitions []*metadatapb.TablePartitionMetadata, oldSchema, newSchema string) {
	if oldSchema == newSchema {
		return
	}
	for _, p := range partitions {
		switch p.Schema {
		case "":
			p.Schema = oldSchema
		case newSchema:
			p.Schema = ""
		}
	}
}

func rewriteParentIndexSchema(meta *metadatapb.DatabaseSchemaMetadata, oldSchema, newSchema string, indexes []*metadatapb.IndexMetadata) {
	if meta == nil || oldSchema == newSchema || len(indexes) == 0 {
		return
	}
	indexNames := make(map[string]bool, len(indexes))
	for _, idx := range indexes {
		if idx.Name != "" {
			indexNames[idx.Name] = true
		}
	}
	if len(indexNames) == 0 {
		return
	}
	var walkPartitions func(partitions []*metadatapb.TablePartitionMetadata)
	walkPartitions = func(partitions []*metadatapb.TablePartitionMetadata) {
		for _, p := range partitions {
			for _, idx := range p.Indexes {
				if idx.ParentIndexSchema == oldSchema && indexNames[idx.ParentIndexName] {
					idx.ParentIndexSchema = newSchema
				}
			}
			walkPartitions(p.Subpartitions)
		}
	}
	for _, s := range meta.Schemas {
		for _, tbl := range s.Tables {
			for _, idx := range tbl.Indexes {
				if idx.ParentIndexSchema == oldSchema && indexNames[idx.ParentIndexName] {
					idx.ParentIndexSchema = newSchema
				}
			}
			walkPartitions(tbl.Partitions)
		}
	}
}

func extSchemaName(cat *Catalog, ext *Extension) string {
	if cat == nil || ext == nil || ext.SchemaOID == 0 {
		return ""
	}
	if s := cat.schemas[ext.SchemaOID]; s != nil {
		return s.Name
	}
	return ""
}

func applyExtensionDiffs(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, extensions []ExtensionDiffEntry, canonicalSchema func(string) string) {
	for _, e := range extensions {
		switch e.Action {
		case DiffDrop:
			result.Extensions = slices.DeleteFunc(result.Extensions, func(ext *metadatapb.ExtensionMetadata) bool {
				return ext.Name == e.Name
			})
			removeExtensionMembers(result, catBefore, e.Name, canonicalSchema)
		case DiffAdd:
			if e.To == nil {
				continue
			}
			updated := &metadatapb.ExtensionMetadata{
				Name:    e.Name,
				Schema:  extSchemaName(catAfter, e.To),
				Version: e.To.Version,
			}
			replaced := false
			for i, existing := range result.Extensions {
				if existing.Name == e.Name {
					updated.Description = existing.Description
					result.Extensions[i] = updated
					replaced = true
					break
				}
			}
			if !replaced {
				result.Extensions = append(result.Extensions, updated)
			}
		case DiffModify:
			if e.To == nil {
				continue
			}
			for _, existing := range result.Extensions {
				if existing.Name == e.Name {
					if schemaName := extSchemaName(catAfter, e.To); schemaName != "" {
						existing.Schema = schemaName
					}
					if e.To.Version != "" {
						existing.Version = e.To.Version
					}
					break
				}
			}
		default:
		}
	}
}

func removeExtensionMembers(result *metadatapb.DatabaseSchemaMetadata, catBefore *Catalog, extName string, canonicalSchema func(string) string) {
	if result == nil || catBefore == nil || extName == "" {
		return
	}
	var dropSchemas []string
	for oid, memberExt := range catBefore.extensionMembers {
		if memberExt != extName {
			continue
		}
		if s := catBefore.schemas[oid]; s != nil && oid != PGCatalogNamespace && oid != PGToastNamespace {
			dropSchemas = append(dropSchemas, canonicalSchema(s.Name))
			continue
		}
		if rel := catBefore.relationByOID[oid]; rel != nil && rel.Schema != nil {
			schemaName := canonicalSchema(rel.Schema.Name)
			if sm := findSchema(result, schemaName); sm != nil {
				switch rel.RelKind {
				case 'v':
					sm.Views = removeViewByName(sm.Views, rel.Name)
				case 'm':
					sm.MaterializedViews = removeMatViewByName(sm.MaterializedViews, rel.Name)
				case 'f':
					sm.ExternalTables = slices.DeleteFunc(sm.ExternalTables, func(ext *metadatapb.ExternalTableMetadata) bool {
						return ext.Name == rel.Name
					})
				default:
					sm.Tables = removeTableByName(sm.Tables, rel.Name)
					_ = popPartitionInSchema(result, schemaName, rel.Name)
				}
			}
		}
		if seq := catBefore.sequenceByOID[oid]; seq != nil && seq.Schema != nil {
			if sm := findSchema(result, canonicalSchema(seq.Schema.Name)); sm != nil {
				sm.Sequences = removeSequenceByName(sm.Sequences, seq.Name)
			}
		}
		if up := catBefore.userProcs[oid]; up != nil && up.Schema != nil {
			if sm := findSchema(result, canonicalSchema(up.Schema.Name)); sm != nil {
				identity := funcIdentity(catBefore, up)
				if up.Kind == 'p' {
					_, sm.Procedures = popProcedureMatchingIdentity(sm.Procedures, up.Schema.Name, identity)
				} else {
					_, sm.Functions = popFunctionMatchingIdentity(sm.Functions, up.Schema.Name, identity)
				}
			}
		}
		if bt := catBefore.typeByOID[oid]; bt != nil {
			if ns := catBefore.schemas[bt.Namespace]; ns != nil {
				if sm := findSchema(result, canonicalSchema(ns.Name)); sm != nil {
					switch bt.Type {
					case 'e':
						sm.EnumTypes = removeEnumByName(sm.EnumTypes, bt.TypeName)
					case 'c':
						sm.CompositeTypes = removeCompositeTypeByName(sm.CompositeTypes, bt.TypeName)
					case 'r':
						_, sm.RangeTypes = popRangeTypeByName(sm.RangeTypes, bt.TypeName)
					}
				}
			}
		}
		if evt := catBefore.eventTriggers[oid]; evt != nil {
			_, result.EventTriggers = popEventTriggerByName(result.EventTriggers, evt.Name)
		}
	}
	for _, schemaName := range dropSchemas {
		result.Schemas = slices.DeleteFunc(result.Schemas, func(s *metadatapb.SchemaMetadata) bool {
			return s.Name == schemaName
		})
	}
}

func hasMatchingSequenceAdd(entries []SequenceDiffEntry, newSchema, name string) bool {
	for _, e := range entries {
		if e.Action == DiffAdd && e.SchemaName == newSchema && e.Name == name {
			return true
		}
	}
	return false
}

func hasMatchingEnumAdd(entries []EnumDiffEntry, newSchema, name string) bool {
	for _, e := range entries {
		if e.Action == DiffAdd && e.SchemaName == newSchema && e.Name == name {
			return true
		}
	}
	return false
}

func hasMatchingCompositeAdd(entries []CompositeTypeDiffEntry, newSchema, name string) bool {
	for _, e := range entries {
		if e.Action == DiffAdd && e.SchemaName == newSchema && e.Name == name {
			return true
		}
	}
	return false
}

func applyRangeDiffs(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, ranges []RangeDiffEntry, canonicalSchema func(string) string) {
	if len(ranges) == 0 {
		return
	}
	addedByOID := make(map[uint32]RangeDiffEntry)
	for _, r := range ranges {
		if r.Action == DiffAdd && r.To != nil && r.To.OID != 0 {
			addedByOID[r.To.OID] = r
		}
	}

	handledDropOIDs := make(map[uint32]bool)
	handledAddOIDs := make(map[uint32]bool)
	for _, r := range ranges {
		if r.Action != DiffDrop || r.From == nil || r.From.OID == 0 {
			continue
		}
		add, ok := addedByOID[r.From.OID]
		if !ok || add.To == nil {
			continue
		}
		var prev *metadatapb.RangeTypeMetadata
		if srcSM := findSchema(result, canonicalSchema(r.SchemaName)); srcSM != nil {
			prev, srcSM.RangeTypes = popRangeTypeByName(srcSM.RangeTypes, r.Name)
		}
		dstSM := findOrCreateSchema(result, add.SchemaName)
		upsertRangeType(dstSM, rangeTypeToProto(catAfter, add.To, prev))
		handledDropOIDs[r.From.OID] = true
		handledAddOIDs[add.To.OID] = true
	}

	for _, r := range ranges {
		if r.Action == DiffDrop && r.From != nil && handledDropOIDs[r.From.OID] {
			continue
		}
		if r.Action == DiffAdd && r.To != nil && handledAddOIDs[r.To.OID] {
			continue
		}
		schemaName := canonicalSchema(r.SchemaName)
		if r.Action == DiffDrop {
			if sm := findSchema(result, schemaName); sm != nil {
				_, sm.RangeTypes = popRangeTypeByName(sm.RangeTypes, r.Name)
			}
			continue
		}
		if r.To == nil {
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch r.Action {
		case DiffAdd, DiffModify:
			upsertRangeType(sm, rangeTypeToProto(catAfter, r.To, findRangeType(sm, r.Name)))
		default:
		}
	}
}

func findRangeType(sm *metadatapb.SchemaMetadata, name string) *metadatapb.RangeTypeMetadata {
	for _, rt := range sm.RangeTypes {
		if rt.Name == name {
			return rt
		}
	}
	return nil
}

func upsertRangeType(sm *metadatapb.SchemaMetadata, updated *metadatapb.RangeTypeMetadata) {
	for i, existing := range sm.RangeTypes {
		if existing.Name == updated.Name {
			sm.RangeTypes[i] = updated
			return
		}
	}
	sm.RangeTypes = append(sm.RangeTypes, updated)
}

func applyEventTriggerDiffs(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, entries []EventTriggerDiffEntry) {
	if len(entries) == 0 {
		return
	}
	addedByOID := make(map[uint32]EventTriggerDiffEntry)
	for _, e := range entries {
		if e.Action == DiffAdd && e.To != nil && e.To.OID != 0 {
			addedByOID[e.To.OID] = e
		}
	}

	handledDropOIDs := make(map[uint32]bool)
	handledAddOIDs := make(map[uint32]bool)
	for _, e := range entries {
		if e.Action != DiffDrop || e.From == nil || e.From.OID == 0 {
			continue
		}
		add, ok := addedByOID[e.From.OID]
		if !ok || add.To == nil {
			continue
		}
		var prev *metadatapb.EventTriggerMetadata
		prev, result.EventTriggers = popEventTriggerByName(result.EventTriggers, e.Name)
		upsertEventTrigger(result, eventTriggerToProto(catBefore, catAfter, e.From, add.To, prev))
		handledDropOIDs[e.From.OID] = true
		handledAddOIDs[add.To.OID] = true
	}

	for _, e := range entries {
		if e.Action == DiffDrop && e.From != nil && handledDropOIDs[e.From.OID] {
			continue
		}
		if e.Action == DiffAdd && e.To != nil && handledAddOIDs[e.To.OID] {
			continue
		}
		switch e.Action {
		case DiffDrop:
			_, result.EventTriggers = popEventTriggerByName(result.EventTriggers, e.Name)
		case DiffAdd, DiffModify:
			if e.To == nil {
				continue
			}
			upsertEventTrigger(result, eventTriggerToProto(catBefore, catAfter, e.From, e.To, findEventTrigger(result, e.Name)))
		default:
		}
	}
}

func eventTriggerToProto(catBefore, catAfter *Catalog, fromEvt, toEvt *EventTrigger, previous *metadatapb.EventTriggerMetadata) *metadatapb.EventTriggerMetadata {
	out := &metadatapb.EventTriggerMetadata{
		Name:    toEvt.Name,
		Event:   toEvt.EventName,
		Tags:    slices.Clone(toEvt.Tags),
		Enabled: toEvt.Enabled != 'D',
	}
	if catAfter != nil {
		if up := catAfter.userProcs[toEvt.FuncOID]; up != nil {
			if up.Schema != nil {
				out.FunctionSchema = up.Schema.Name
			}
			out.FunctionName = up.Name
		} else if p := catAfter.procByOID[toEvt.FuncOID]; p != nil {
			out.FunctionName = p.Name
		}
	}
	if previous != nil {
		out.Comment = previous.Comment
		out.SkipDump = previous.SkipDump
		if previous.Definition != "" && fromEvt != nil && fromEvt.Name == toEvt.Name && !structuralEventTriggerChange(catBefore, catAfter, fromEvt, toEvt) {
			out.Definition = previous.Definition
		}
	}
	if out.Definition == "" && catAfter != nil {
		if ops := buildCreateEventTriggerOps(catAfter, toEvt); len(ops) > 0 {
			out.Definition = ops[0].SQL
		}
	}
	return out
}

func findEventTrigger(meta *metadatapb.DatabaseSchemaMetadata, name string) *metadatapb.EventTriggerMetadata {
	for _, et := range meta.EventTriggers {
		if et.Name == name {
			return et
		}
	}
	return nil
}

func popEventTriggerByName(triggers []*metadatapb.EventTriggerMetadata, name string) (*metadatapb.EventTriggerMetadata, []*metadatapb.EventTriggerMetadata) {
	var popped *metadatapb.EventTriggerMetadata
	out := make([]*metadatapb.EventTriggerMetadata, 0, len(triggers))
	for _, et := range triggers {
		if et.Name == name && popped == nil {
			popped = et
			continue
		}
		out = append(out, et)
	}
	return popped, out
}

func upsertEventTrigger(meta *metadatapb.DatabaseSchemaMetadata, updated *metadatapb.EventTriggerMetadata) {
	for i, existing := range meta.EventTriggers {
		if existing.Name == updated.Name {
			meta.EventTriggers[i] = updated
			return
		}
	}
	meta.EventTriggers = append(meta.EventTriggers, updated)
}

func popSequenceByName(seqs []*metadatapb.SequenceMetadata, name string) (*metadatapb.SequenceMetadata, []*metadatapb.SequenceMetadata) {
	var popped *metadatapb.SequenceMetadata
	out := make([]*metadatapb.SequenceMetadata, 0, len(seqs))
	for _, s := range seqs {
		if s.Name == name && popped == nil {
			popped = s
			continue
		}
		out = append(out, s)
	}
	return popped, out
}

func mergeSequenceMetadata(updated, existing *metadatapb.SequenceMetadata) {
	if existing == nil {
		return
	}
	if updated.Comment == "" {
		updated.Comment = existing.Comment
	}
	if !updated.SkipDump {
		updated.SkipDump = existing.SkipDump
	}
	if updated.LastValue == "" {
		updated.LastValue = existing.LastValue
	}
}

func upsertSequence(sm *metadatapb.SchemaMetadata, updated *metadatapb.SequenceMetadata) {
	for i, existing := range sm.Sequences {
		if existing.Name == updated.Name {
			mergeSequenceMetadata(updated, existing)
			sm.Sequences[i] = updated
			return
		}
	}
	sm.Sequences = append(sm.Sequences, updated)
}

func applySequenceDiffs(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, sequences []SequenceDiffEntry, renamedSchemas map[string]string, canonicalSchema func(string) string) {
	droppedByOID := make(map[uint32]SequenceDiffEntry)
	addedByOID := make(map[uint32]SequenceDiffEntry)
	for _, seq := range sequences {
		if isTempSequenceDiff(seq, canonicalSchema) {
			continue
		}
		if seq.Action == DiffDrop && seq.From != nil && seq.From.OID != 0 {
			droppedByOID[seq.From.OID] = seq
		} else if seq.Action == DiffAdd && seq.To != nil && seq.To.OID != 0 {
			addedByOID[seq.To.OID] = seq
		}
	}
	handledOIDs := make(map[uint32]bool)
	for oid, dropEntry := range droppedByOID {
		addEntry, ok := addedByOID[oid]
		if !ok {
			continue
		}
		handledOIDs[oid] = true
		srcSchema := canonicalSchema(dropEntry.SchemaName)
		dstSchema := canonicalSchema(addEntry.SchemaName)
		var prevSeq *metadatapb.SequenceMetadata
		if srcSM := findSchema(result, srcSchema); srcSM != nil {
			prevSeq, srcSM.Sequences = popSequenceByName(srcSM.Sequences, dropEntry.Name)
		}
		updated := sequenceToProto(catAfter, addEntry.To)
		mergeSequenceMetadata(updated, prevSeq)
		upsertSequence(findOrCreateSchema(result, dstSchema), updated)
	}
	for _, seq := range sequences {
		if isTempSequenceDiff(seq, canonicalSchema) {
			continue
		}
		if (seq.From != nil && seq.From.OID != 0 && handledOIDs[seq.From.OID]) ||
			(seq.To != nil && seq.To.OID != 0 && handledOIDs[seq.To.OID]) {
			continue
		}
		schemaName := canonicalSchema(seq.SchemaName)
		if seq.Action == DiffDrop {
			if renamedSchemas[seq.SchemaName] != "" && hasMatchingSequenceAdd(sequences, renamedSchemas[seq.SchemaName], seq.Name) {
				continue
			}
			if sm := findSchema(result, schemaName); sm != nil {
				sm.Sequences = removeSequenceByName(sm.Sequences, seq.Name)
			}
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		if (seq.Action == DiffAdd || seq.Action == DiffModify) && seq.To != nil {
			upsertSequence(sm, sequenceToProto(catAfter, seq.To))
		}
	}
}

func upsertEnum(sm *metadatapb.SchemaMetadata, name string, values []string) {
	for _, existing := range sm.EnumTypes {
		if existing.Name == name {
			existing.Values = slices.Clone(values)
			return
		}
	}
	sm.EnumTypes = append(sm.EnumTypes, &metadatapb.EnumTypeMetadata{
		Name:   name,
		Values: values,
	})
}

func findSchema(meta *metadatapb.DatabaseSchemaMetadata, name string) *metadatapb.SchemaMetadata {
	for _, s := range meta.Schemas {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func findOrCreateSchema(meta *metadatapb.DatabaseSchemaMetadata, name string) *metadatapb.SchemaMetadata {
	for _, s := range meta.Schemas {
		if s.Name == name {
			return s
		}
	}
	s := &metadatapb.SchemaMetadata{Name: name}
	meta.Schemas = append(meta.Schemas, s)
	return s
}

func removeSchema(schemas []*metadatapb.SchemaMetadata, name string) []*metadatapb.SchemaMetadata {
	out := make([]*metadatapb.SchemaMetadata, 0, len(schemas))
	for _, s := range schemas {
		if s.Name != name {
			out = append(out, s)
		}
	}
	return out
}

func addRelation(result *metadatapb.DatabaseSchemaMetadata, sm *metadatapb.SchemaMetadata, cat *Catalog, rel RelationDiffEntry) {
	if rel.To == nil {
		return
	}
	switch rel.To.RelKind {
	case 'r', 'p':
		if rel.To.PartitionOf != 0 {
			if appendPartitionToParent(result, sm, cat, rel) {
				return
			}
		}
		sm.Tables = append(sm.Tables, relationToTableProto(cat, rel.To))
	case 'v':
		sm.Views = append(sm.Views, relationToViewProto(cat, rel.To))
	case 'm':
		sm.MaterializedViews = append(sm.MaterializedViews, relationToMatViewProto(cat, rel.To))
	default:
	}
}

func appendPartitionToParent(result *metadatapb.DatabaseSchemaMetadata, fallbackSchema *metadatapb.SchemaMetadata, cat *Catalog, rel RelationDiffEntry) bool {
	parentRel := cat.GetRelationByOID(rel.To.PartitionOf)
	if parentRel == nil {
		return false
	}
	parentSchemaName := rel.SchemaName
	if parentRel.Schema != nil && parentRel.Schema.Name != "" {
		parentSchemaName = parentRel.Schema.Name
	}
	part := relationToPartitionProto(cat, rel.To, parentRel, parentSchemaName)
	return attachPartitionProto(result, fallbackSchema, parentRel, parentSchemaName, part)
}

func attachPartitionProto(result *metadatapb.DatabaseSchemaMetadata, fallbackSchema *metadatapb.SchemaMetadata, parentRel *Relation, parentSchemaName string, part *metadatapb.TablePartitionMetadata) bool {
	if parentRel == nil {
		return false
	}
	parentSM := findSchema(result, parentSchemaName)
	if parentSM == nil {
		parentSM = fallbackSchema
	}
	if parentSM == nil {
		return false
	}
	if parentTbl := findTable(parentSM, parentRel.Name); parentTbl != nil {
		parentTbl.Partitions = append(parentTbl.Partitions, part)
		return true
	}
	if parentPart := findPartitionInSchema(result, parentSchemaName, parentRel.Name); parentPart != nil {
		parentPart.Subpartitions = append(parentPart.Subpartitions, part)
		return true
	}
	return false
}

func relationToPartitionProto(c *Catalog, rel, parentRel *Relation, parentSchema string) *metadatapb.TablePartitionMetadata {
	part := &metadatapb.TablePartitionMetadata{
		Name:  rel.Name,
		Value: formatPartitionBound(rel.PartitionBound),
	}
	if rel.Schema != nil && rel.Schema.Name != "" && rel.Schema.Name != parentSchema {
		part.Schema = rel.Schema.Name
	}
	var strategy byte
	if parentRel != nil && parentRel.PartitionInfo != nil {
		strategy = parentRel.PartitionInfo.Strategy
		part.Expression = strings.TrimPrefix(formatPartitionByClause(parentRel), " PARTITION BY ")
	} else if rel.PartitionBound != nil {
		strategy = rel.PartitionBound.Strategy
	}
	switch strategy {
	case 'r':
		part.Type = metadatapb.TablePartitionMetadata_RANGE
	case 'l':
		part.Type = metadatapb.TablePartitionMetadata_LIST
	case 'h':
		part.Type = metadatapb.TablePartitionMetadata_HASH
	}
	tbl := relationToTableProto(c, rel)
	part.Indexes = tbl.Indexes
	part.CheckConstraints = tbl.CheckConstraints
	part.ExcludeConstraints = tbl.ExcludeConstraints
	return part
}

func dropRelation(result *metadatapb.DatabaseSchemaMetadata, schemaName string, rel RelationDiffEntry) {
	name := rel.Name
	if sm := findSchema(result, schemaName); sm != nil {
		sm.Tables = removeTableByName(sm.Tables, name)
		sm.Views = removeViewByName(sm.Views, name)
		sm.MaterializedViews = removeMatViewByName(sm.MaterializedViews, name)
	}
	_ = popPartitionInSchema(result, schemaName, name)
}

func popTableByName(tables []*metadatapb.TableMetadata, name string) (*metadatapb.TableMetadata, []*metadatapb.TableMetadata) {
	var popped *metadatapb.TableMetadata
	out := make([]*metadatapb.TableMetadata, 0, len(tables))
	for _, t := range tables {
		if t.Name == name && popped == nil {
			popped = t
			continue
		}
		out = append(out, t)
	}
	return popped, out
}

func upsertTableInSchema(sm *metadatapb.SchemaMetadata, tbl *metadatapb.TableMetadata) {
	for i, existing := range sm.Tables {
		if existing.Name == tbl.Name {
			sm.Tables[i] = tbl
			return
		}
	}
	sm.Tables = append(sm.Tables, tbl)
}

func findPartitionInSchema(meta *metadatapb.DatabaseSchemaMetadata, targetSchema, name string) *metadatapb.TablePartitionMetadata {
	var walk func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) *metadatapb.TablePartitionMetadata
	walk = func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) *metadatapb.TablePartitionMetadata {
		for _, p := range partitions {
			partSchema := partitionSchemaName(parentSchema, p)
			if p.Name == name && partSchema == targetSchema {
				return p
			}
			if found := walk(partSchema, p.Subpartitions); found != nil {
				return found
			}
		}
		return nil
	}
	for _, s := range meta.Schemas {
		for _, t := range s.Tables {
			if found := walk(s.Name, t.Partitions); found != nil {
				return found
			}
		}
	}
	return nil
}

func popPartitionInSchema(meta *metadatapb.DatabaseSchemaMetadata, targetSchema, name string) *metadatapb.TablePartitionMetadata {
	var popped *metadatapb.TablePartitionMetadata
	var walk func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) []*metadatapb.TablePartitionMetadata
	walk = func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) []*metadatapb.TablePartitionMetadata {
		if len(partitions) == 0 {
			return partitions
		}
		out := make([]*metadatapb.TablePartitionMetadata, 0, len(partitions))
		for _, p := range partitions {
			partSchema := partitionSchemaName(parentSchema, p)
			if p.Name == name && partSchema == targetSchema && popped == nil {
				popped = p
				continue
			}
			p.Subpartitions = walk(partSchema, p.Subpartitions)
			out = append(out, p)
		}
		return out
	}
	for _, s := range meta.Schemas {
		for _, t := range s.Tables {
			t.Partitions = walk(s.Name, t.Partitions)
		}
	}
	return popped
}

func samePartitionParent(catBefore, catAfter *Catalog, fromRel, toRel *Relation) bool {
	if fromRel == nil || toRel == nil {
		return false
	}
	if fromRel.PartitionOf == 0 && toRel.PartitionOf == 0 {
		return true
	}
	if (fromRel.PartitionOf == 0) != (toRel.PartitionOf == 0) {
		return false
	}
	fromParent := catBefore.GetRelationByOID(fromRel.PartitionOf)
	toParent := catAfter.GetRelationByOID(toRel.PartitionOf)
	if fromParent == nil || toParent == nil {
		return fromRel.PartitionOf == toRel.PartitionOf
	}
	fromSchema, toSchema := "", ""
	if fromParent.Schema != nil {
		fromSchema = fromParent.Schema.Name
	}
	if toParent.Schema != nil {
		toSchema = toParent.Schema.Name
	}
	return fromSchema == toSchema && fromParent.Name == toParent.Name
}

func modifyRelation(result *metadatapb.DatabaseSchemaMetadata, sm *metadatapb.SchemaMetadata, catBefore, cat *Catalog, rel RelationDiffEntry) {
	if rel.To == nil {
		return
	}
	switch rel.To.RelKind {
	case 'r', 'p':
		wasPartition := rel.From != nil && rel.From.PartitionOf != 0
		isPartition := rel.To.PartitionOf != 0
		if wasPartition && !isPartition {
			prevPart := popPartitionInSchema(result, sm.Name, rel.Name)
			tbl := relationToTableProto(cat, rel.To)
			if prevPart != nil {
				tbl.Comment = prevPart.Comment
				preserveIndexMetadata(tbl.Indexes, prevPart.Indexes)
				preserveConstraintComments(tbl, &metadatapb.TableMetadata{CheckConstraints: prevPart.CheckConstraints})
			}
			upsertTableInSchema(sm, tbl)
			return
		}
		if isPartition && (!wasPartition || !samePartitionParent(catBefore, cat, rel.From, rel.To)) {
			var prevComment string
			var prevIndexes []*metadatapb.IndexMetadata
			var prevChecks []*metadatapb.CheckConstraintMetadata
			if prevTbl, remaining := popTableByName(sm.Tables, rel.Name); prevTbl != nil {
				sm.Tables = remaining
				prevComment = prevTbl.Comment
				prevIndexes = prevTbl.Indexes
				prevChecks = prevTbl.CheckConstraints
			} else if prevPart := popPartitionInSchema(result, sm.Name, rel.Name); prevPart != nil {
				prevComment = prevPart.Comment
				prevIndexes = prevPart.Indexes
				prevChecks = prevPart.CheckConstraints
			}
			parentRel := cat.GetRelationByOID(rel.To.PartitionOf)
			parentSchemaName := rel.SchemaName
			if parentRel != nil && parentRel.Schema != nil && parentRel.Schema.Name != "" {
				parentSchemaName = parentRel.Schema.Name
			}
			part := relationToPartitionProto(cat, rel.To, parentRel, parentSchemaName)
			if prevComment != "" {
				part.Comment = prevComment
			}
			preserveIndexMetadata(part.Indexes, prevIndexes)
			preserveConstraintComments(&metadatapb.TableMetadata{CheckConstraints: part.CheckConstraints}, &metadatapb.TableMetadata{CheckConstraints: prevChecks})
			if !attachPartitionProto(result, sm, parentRel, parentSchemaName, part) {
				upsertTableInSchema(sm, relationToTableProto(cat, rel.To))
			}
			return
		}
		tbl := findTable(sm, rel.Name)
		if tbl == nil {
			if p := findPartitionInSchema(result, sm.Name, rel.Name); p != nil {
				applyPartitionDiffs(p, catBefore, cat, rel)
				return
			}
			sm.Tables = append(sm.Tables, relationToTableProto(cat, rel.To))
			return
		}
		applyColumnDiffs(result, sm.Name, rel.Name, &tbl.Columns, cat, rel)
		applyIndexDiffs(tbl, cat, rel)
		applyConstraintDiffs(tbl, catBefore, cat, rel)
		applyTriggerDiffs(&tbl.Triggers, cat, rel)
		applyRuleDiffs(&tbl.Rules, cat, rel)
		if (rel.From != nil && rel.From.OwnerName != rel.To.OwnerName) || (rel.To.OwnerName != "" && tbl.Owner != rel.To.OwnerName) {
			tbl.Owner = rel.To.OwnerName
		}
	case 'v':
		for i, existing := range sm.Views {
			if existing.Name == rel.Name {
				sm.Views[i] = mergeViewProto(cat, existing, rel, rel.Name)
				return
			}
		}
		sm.Views = append(sm.Views, relationToViewProto(cat, rel.To))
	case 'm':
		for i, existing := range sm.MaterializedViews {
			if existing.Name == rel.Name {
				sm.MaterializedViews[i] = mergeMatViewProto(cat, existing, rel.To, rel.Name)
				return
			}
		}
		sm.MaterializedViews = append(sm.MaterializedViews, relationToMatViewProto(cat, rel.To))
	default:
	}
}

func mergeViewProto(cat *Catalog, existing *metadatapb.ViewMetadata, relDiff RelationDiffEntry, name string) *metadatapb.ViewMetadata {
	updated := relationToViewProto(cat, relDiff.To)
	updated.Name = name
	updated.Comment = existing.Comment
	preserveColumnMetadata(updated.Columns, existing.Columns, relDiff.From, relDiff.To)
	if len(updated.Columns) == 0 {
		updated.Columns = existing.Columns
	}
	if len(relDiff.Rules) > 0 {
		updated.Rules = slices.Clone(existing.Rules)
		applyRuleDiffs(&updated.Rules, cat, relDiff)
	} else if len(updated.Rules) == 0 {
		updated.Rules = existing.Rules
	}
	updated.SkipDump = existing.SkipDump
	preserveTriggerComments(updated.Triggers, existing.Triggers)
	return updated
}

func mergeMatViewProto(cat *Catalog, existing *metadatapb.MaterializedViewMetadata, rel *Relation, name string) *metadatapb.MaterializedViewMetadata {
	updated := relationToMatViewProto(cat, rel)
	updated.Name = name
	updated.Comment = existing.Comment
	updated.SkipDump = existing.SkipDump
	preserveTriggerComments(updated.Triggers, existing.Triggers)
	preserveIndexMetadata(updated.Indexes, existing.Indexes)
	return updated
}

func findExternalTable(sm *metadatapb.SchemaMetadata, name string) *metadatapb.ExternalTableMetadata {
	for _, ext := range sm.ExternalTables {
		if ext.Name == name {
			return ext
		}
	}
	return nil
}

func relationToExternalTableProto(cat *Catalog, rel *Relation) *metadatapb.ExternalTableMetadata {
	tbl := relationToTableProto(cat, rel)
	return &metadatapb.ExternalTableMetadata{
		Name:               rel.Name,
		ExternalServerName: rel.ForeignServer,
		Columns:            tbl.Columns,
		Options:            slices.Clone(rel.ForeignOptions),
		CheckConstraints:   tbl.CheckConstraints,
		Triggers:           tbl.Triggers,
	}
}

func applyForeignTableDiffs(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, entries []RelationDiffEntry, canonicalSchema func(string) string) {
	if len(entries) == 0 {
		return
	}
	renamed := renameForeignTables(result, catAfter, entries, canonicalSchema)
	for _, rel := range entries {
		if (rel.From != nil && renamed[rel.From.OID]) || (rel.To != nil && renamed[rel.To.OID]) {
			continue
		}
		schemaName := canonicalSchema(rel.SchemaName)
		if isTempRelation(rel.SchemaName, rel.From) || isTempRelation(rel.SchemaName, rel.To) || isTempSchemaName(schemaName) {
			continue
		}
		if rel.Action == DiffDrop {
			if sm := findSchema(result, schemaName); sm != nil {
				sm.ExternalTables = slices.DeleteFunc(sm.ExternalTables, func(ext *metadatapb.ExternalTableMetadata) bool {
					return ext.Name == rel.Name
				})
			}
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch rel.Action {
		case DiffAdd:
			if rel.To != nil {
				sm.ExternalTables = append(sm.ExternalTables, relationToExternalTableProto(catAfter, rel.To))
			}
		case DiffModify:
			if rel.To == nil {
				continue
			}
			ext := findExternalTable(sm, rel.Name)
			if ext == nil {
				sm.ExternalTables = append(sm.ExternalTables, relationToExternalTableProto(catAfter, rel.To))
				continue
			}
			applyColumnDiffs(result, sm.Name, rel.Name, &ext.Columns, catAfter, rel)
			tbl := &metadatapb.TableMetadata{CheckConstraints: ext.CheckConstraints}
			applyConstraintDiffs(tbl, catBefore, catAfter, rel)
			ext.CheckConstraints = tbl.CheckConstraints
			applyTriggerDiffs(&ext.Triggers, catAfter, rel)
			if rel.To.ForeignServer != "" {
				ext.ExternalServerName = rel.To.ForeignServer
			}
			if len(rel.To.ForeignOptions) > 0 || (rel.From != nil && len(rel.From.ForeignOptions) > 0) {
				ext.Options = slices.Clone(rel.To.ForeignOptions)
			}
		default:
		}
	}
}

func renameForeignTables(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, entries []RelationDiffEntry, canonicalSchema func(string) string) map[uint32]bool {
	added := make(map[uint32]RelationDiffEntry)
	for _, rel := range entries {
		if rel.Action == DiffAdd && rel.To != nil {
			added[rel.To.OID] = rel
		}
	}
	renamed := make(map[uint32]bool)
	for _, rel := range entries {
		if rel.Action != DiffDrop || rel.From == nil {
			continue
		}
		add, ok := added[rel.From.OID]
		srcSchema := canonicalSchema(rel.SchemaName)
		sm := findSchema(result, srcSchema)
		if !ok || sm == nil {
			continue
		}
		ext := findExternalTable(sm, rel.Name)
		if ext == nil {
			continue
		}
		rebuilt := relationToExternalTableProto(catAfter, add.To)
		preserveColumnMetadata(rebuilt.Columns, ext.Columns, rel.From, add.To)
		preserveCheckConstraintComments(rebuilt.CheckConstraints, ext.CheckConstraints)
		preserveTriggerComments(rebuilt.Triggers, ext.Triggers)
		ext.Name = add.Name
		ext.Columns = rebuilt.Columns
		ext.CheckConstraints = rebuilt.CheckConstraints
		ext.Triggers = rebuilt.Triggers
		if add.To.ForeignServer != "" {
			ext.ExternalServerName = add.To.ForeignServer
		}
		if len(add.To.ForeignOptions) > 0 || len(rel.From.ForeignOptions) > 0 {
			ext.Options = slices.Clone(add.To.ForeignOptions)
		}
		if add.SchemaName != srcSchema {
			sm.ExternalTables = slices.DeleteFunc(sm.ExternalTables, func(e *metadatapb.ExternalTableMetadata) bool {
				return e.Name == add.Name
			})
			findOrCreateSchema(result, add.SchemaName).ExternalTables = append(findOrCreateSchema(result, add.SchemaName).ExternalTables, ext)
		}
		rewriteColumnRenamesForRelation(result, srcSchema, rel.Name, rel.From, add.To)
		rewriteTableDependencyReferences(result, srcSchema, rel.Name, add.SchemaName, add.Name)
		renamed[rel.From.OID] = true
	}
	return renamed
}

func preserveTriggerComments(dst, src []*metadatapb.TriggerMetadata) {
	if len(dst) == 0 || len(src) == 0 {
		return
	}
	comments := make(map[string]string, len(src))
	for _, tr := range src {
		if tr.Comment != "" {
			comments[tr.Name] = tr.Comment
		}
	}
	for _, tr := range dst {
		if c, ok := comments[tr.Name]; ok && tr.Comment == "" {
			tr.Comment = c
		}
	}
}

func preserveCheckConstraintComments(dst, src []*metadatapb.CheckConstraintMetadata) {
	if len(dst) == 0 || len(src) == 0 {
		return
	}
	comments := make(map[string]string, len(src))
	for _, chk := range src {
		if chk.Comment != "" {
			comments[chk.Name] = chk.Comment
		}
	}
	for _, chk := range dst {
		if c, ok := comments[chk.Name]; ok && chk.Comment == "" {
			chk.Comment = c
		}
	}
}

func preserveConstraintComments(dst, src *metadatapb.TableMetadata) {
	if len(dst.ForeignKeys) > 0 && len(src.ForeignKeys) > 0 {
		fkComments := make(map[string]string, len(src.ForeignKeys))
		for _, fk := range src.ForeignKeys {
			if fk.Comment != "" {
				fkComments[fk.Name] = fk.Comment
			}
		}
		for _, fk := range dst.ForeignKeys {
			if c, ok := fkComments[fk.Name]; ok && fk.Comment == "" {
				fk.Comment = c
			}
		}
	}
	preserveCheckConstraintComments(dst.CheckConstraints, src.CheckConstraints)
}

func preserveIndexMetadata(dst, src []*metadatapb.IndexMetadata) {
	if len(dst) == 0 || len(src) == 0 {
		return
	}
	byName := make(map[string]*metadatapb.IndexMetadata, len(src))
	for _, idx := range src {
		byName[idx.Name] = idx
	}
	for _, idx := range dst {
		prev, ok := byName[idx.Name]
		if !ok {
			continue
		}
		if idx.Comment == "" {
			idx.Comment = prev.Comment
		}
		if len(idx.Constraints) == 0 {
			idx.Constraints = prev.Constraints
		}
	}
}

// renamePartitions renames each partition or table the diff renames where it
// stands, preserving nested partitions and metadata the catalog does not track.
// The diff shows a rename as a drop and an add of the same relation, which would
// otherwise turn a partition into a table of its own or drop a table's
// partitions. It returns the OIDs it handled.
func renamePartitions(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, relations []RelationDiffEntry, canonicalSchema func(string) string) map[uint32]bool {
	added := make(map[uint32]RelationDiffEntry)
	for _, rel := range relations {
		if rel.Action == DiffAdd && rel.To != nil {
			added[rel.To.OID] = rel
		}
	}
	renamed := make(map[uint32]bool)
	for _, rel := range relations {
		if rel.Action != DiffDrop || rel.From == nil {
			continue
		}
		add, ok := added[rel.From.OID]
		srcSchema := canonicalSchema(rel.SchemaName)
		sm := findSchema(result, srcSchema)
		if !ok {
			continue
		}
		if p := findPartitionInSchema(result, srcSchema, rel.Name); p != nil {
			tbl := relationToTableProto(catAfter, add.To)
			preserveIndexMetadata(tbl.Indexes, p.Indexes)
			preserveCheckConstraintComments(tbl.CheckConstraints, p.CheckConstraints)
			p.Name = add.Name
			p.Indexes, p.CheckConstraints, p.ExcludeConstraints = tbl.Indexes, tbl.CheckConstraints, tbl.ExcludeConstraints
			if add.SchemaName != srcSchema {
				rewritePartitionSchemasOnParentMove(p.Subpartitions, srcSchema, add.SchemaName)
				rewriteParentIndexSchema(result, srcSchema, add.SchemaName, tbl.Indexes)
				parentRel := catAfter.GetRelationByOID(add.To.PartitionOf)
				if parentRel != nil && parentRel.Schema != nil && parentRel.Schema.Name == add.SchemaName {
					p.Schema = ""
				} else {
					p.Schema = add.SchemaName
				}
			}
			rewriteTableDependencyReferences(result, srcSchema, rel.Name, add.SchemaName, add.Name)
			renamed[rel.From.OID] = true
			continue
		}
		if sm == nil {
			continue
		}
		if rel.From.RelKind == 'v' {
			for i, existing := range sm.Views {
				if existing.Name == rel.Name {
					renamedDiff := rel
					renamedDiff.SchemaName = add.SchemaName
					renamedDiff.Name = add.Name
					renamedDiff.To = add.To
					renamedDiff.Rules = diffRules(catBefore, catAfter, rel.From, add.To)
					updated := mergeViewProto(catAfter, existing, renamedDiff, add.Name)
					if add.SchemaName != srcSchema {
						sm.Views = removeViewByName(sm.Views, rel.Name)
						findOrCreateSchema(result, add.SchemaName).Views = append(findOrCreateSchema(result, add.SchemaName).Views, updated)
					} else {
						sm.Views[i] = updated
					}
					rewriteTableDependencyReferences(result, srcSchema, rel.Name, add.SchemaName, add.Name)
					renamed[rel.From.OID] = true
					break
				}
			}
			continue
		}
		if rel.From.RelKind == 'm' {
			for i, existing := range sm.MaterializedViews {
				if existing.Name == rel.Name {
					updated := mergeMatViewProto(catAfter, existing, add.To, add.Name)
					if add.SchemaName != srcSchema {
						sm.MaterializedViews = removeMatViewByName(sm.MaterializedViews, rel.Name)
						findOrCreateSchema(result, add.SchemaName).MaterializedViews = append(findOrCreateSchema(result, add.SchemaName).MaterializedViews, updated)
					} else {
						sm.MaterializedViews[i] = updated
					}
					rewriteTableDependencyReferences(result, srcSchema, rel.Name, add.SchemaName, add.Name)
					renamed[rel.From.OID] = true
					break
				}
			}
			continue
		}
		if tbl := findTable(sm, rel.Name); tbl != nil {
			rebuilt := relationToTableProto(catAfter, add.To)
			preserveColumnMetadata(rebuilt.Columns, tbl.Columns, rel.From, add.To)
			preserveIndexMetadata(rebuilt.Indexes, tbl.Indexes)
			preserveConstraintComments(rebuilt, tbl)
			preserveTriggerComments(rebuilt.Triggers, tbl.Triggers)
			tbl.Name = add.Name
			tbl.Columns = rebuilt.Columns
			tbl.Indexes = rebuilt.Indexes
			tbl.ForeignKeys = rebuilt.ForeignKeys
			tbl.CheckConstraints = rebuilt.CheckConstraints
			tbl.ExcludeConstraints = rebuilt.ExcludeConstraints
			tbl.Triggers = rebuilt.Triggers
			if len(rebuilt.Rules) > 0 || (rel.From != nil && len(rel.From.Rules) > 0) {
				tbl.Rules = rebuilt.Rules
			}
			if rel.From.OwnerName != add.To.OwnerName {
				tbl.Owner = add.To.OwnerName
			}
			if add.SchemaName != srcSchema {
				rewritePartitionSchemasOnParentMove(tbl.Partitions, srcSchema, add.SchemaName)
				rewriteParentIndexSchema(result, srcSchema, add.SchemaName, tbl.Indexes)
				sm.Tables = removeTableByName(sm.Tables, rel.Name)
				targetSchema := findOrCreateSchema(result, add.SchemaName)
				targetSchema.Tables = append(targetSchema.Tables, tbl)
			} else if add.Name != rel.Name {
				for _, seq := range sm.Sequences {
					if seq.OwnerTable == rel.Name {
						seq.OwnerTable = add.Name
					}
				}
			}
			rewriteColumnRenamesForRelation(result, srcSchema, rel.Name, rel.From, add.To)
			rewriteTableDependencyReferences(result, srcSchema, rel.Name, add.SchemaName, add.Name)
			renamed[rel.From.OID] = true
		}
	}
	return renamed
}

// applyPartitionDiffs applies a partition's bound, index, and constraint
// changes. A partition's metadata holds only those; its columns belong to its
// parent table.
func applyPartitionDiffs(p *metadatapb.TablePartitionMetadata, catBefore, cat *Catalog, rel RelationDiffEntry) {
	if rel.To != nil && rel.To.PartitionBound != nil && rel.To.PartitionOf != 0 {
		parentRel := cat.GetRelationByOID(rel.To.PartitionOf)
		parentSchema := rel.SchemaName
		if parentRel != nil && parentRel.Schema != nil && parentRel.Schema.Name != "" {
			parentSchema = parentRel.Schema.Name
		}
		meta := relationToPartitionProto(cat, rel.To, parentRel, parentSchema)
		p.Value = meta.Value
		p.Type = meta.Type
		p.Expression = meta.Expression
		p.Schema = meta.Schema
	}
	tbl := &metadatapb.TableMetadata{Indexes: p.Indexes, CheckConstraints: p.CheckConstraints, ExcludeConstraints: p.ExcludeConstraints}
	applyIndexDiffs(tbl, cat, rel)
	applyConstraintDiffs(tbl, catBefore, cat, rel)
	p.Indexes, p.CheckConstraints, p.ExcludeConstraints = tbl.Indexes, tbl.CheckConstraints, tbl.ExcludeConstraints
}

func mergeColumnProto(cat *Catalog, existing *metadatapb.ColumnMetadata, fromCol, toCol *Column) *metadatapb.ColumnMetadata {
	updated := columnToProto(cat, toCol)
	updated.Comment = existing.Comment
	updated.CharacterSet = existing.CharacterSet
	if len(toCol.ForeignOptions) == 0 && (fromCol == nil || len(fromCol.ForeignOptions) == 0) && len(existing.Options) > 0 {
		updated.Options = slices.Clone(existing.Options)
	}
	if fromCol != nil &&
		fromCol.CollationName == toCol.CollationName &&
		columnCollationProtoSchema(cat, fromCol) == columnCollationProtoSchema(cat, toCol) {
		updated.Collation = existing.Collation
		updated.CollationSchema = existing.CollationSchema
	}
	return updated
}

func preserveColumnMetadata(dst, src []*metadatapb.ColumnMetadata, fromRel, toRel *Relation) {
	if len(dst) == 0 || len(src) == 0 {
		return
	}
	srcByName := make(map[string]*metadatapb.ColumnMetadata, len(src))
	for _, col := range src {
		srcByName[col.Name] = col
	}
	fromByAttNum := make(map[int16]*Column)
	if fromRel != nil {
		for _, c := range fromRel.Columns {
			if c != nil {
				fromByAttNum[c.AttNum] = c
			}
		}
	}
	toByName := make(map[string]*Column)
	if toRel != nil {
		for _, c := range toRel.Columns {
			if c != nil {
				toByName[c.Name] = c
			}
		}
	}
	for i, col := range dst {
		if tc := toByName[col.Name]; tc != nil {
			if fc := fromByAttNum[tc.AttNum]; fc != nil {
				if existing := srcByName[fc.Name]; existing != nil {
					if col.Comment == "" {
						col.Comment = existing.Comment
					}
					if col.CharacterSet == "" {
						col.CharacterSet = existing.CharacterSet
					}
					if len(col.Options) == 0 && len(fc.ForeignOptions) == 0 && len(tc.ForeignOptions) == 0 {
						col.Options = existing.Options
					}
					if fc.CollationName == tc.CollationName {
						col.Collation = existing.Collation
						col.CollationSchema = existing.CollationSchema
					}
					dst[i] = col
				}
			}
		}
	}
}

func applyColumnDiffs(result *metadatapb.DatabaseSchemaMetadata, schemaName, tableName string, cols *[]*metadatapb.ColumnMetadata, cat *Catalog, rel RelationDiffEntry) {
	addedByAttNum := make(map[int16]ColumnDiffEntry)
	for _, cd := range rel.Columns {
		if cd.Action == DiffAdd && cd.To != nil && cd.To.AttNum > 0 {
			addedByAttNum[cd.To.AttNum] = cd
		}
	}
	handledAttNum := make(map[int16]bool)
	for _, cd := range rel.Columns {
		if cd.Action != DiffDrop || cd.From == nil || cd.From.AttNum <= 0 {
			continue
		}
		add, ok := addedByAttNum[cd.From.AttNum]
		if !ok {
			continue
		}
		for i, col := range *cols {
			if col.Name == cd.Name {
				(*cols)[i] = mergeColumnProto(cat, col, cd.From, add.To)
				handledAttNum[cd.From.AttNum] = true
				rewriteColumnDependencyReferences(result, schemaName, tableName, cd.Name, add.To.Name)
				break
			}
		}
	}
	for _, cd := range rel.Columns {
		switch cd.Action {
		case DiffAdd:
			if cd.To == nil || (cd.To.AttNum > 0 && handledAttNum[cd.To.AttNum]) {
				continue
			}
			*cols = append(*cols, columnToProto(cat, cd.To))
		case DiffDrop:
			if cd.From != nil && cd.From.AttNum > 0 && handledAttNum[cd.From.AttNum] {
				continue
			}
			*cols = removeColumnByName(*cols, cd.Name)
		case DiffModify:
			if cd.To != nil {
				for i, col := range *cols {
					if col.Name == cd.Name {
						(*cols)[i] = mergeColumnProto(cat, col, cd.From, cd.To)
						break
					}
				}
			}
		default:
		}
	}
}

func applyIndexDiffs(tbl *metadatapb.TableMetadata, cat *Catalog, rel RelationDiffEntry) {
	for _, id := range rel.Indexes {
		switch id.Action {
		case DiffAdd:
			if id.To != nil {
				tbl.Indexes = append(tbl.Indexes, indexToProto(cat, rel.To, id.To))
			}
		case DiffDrop:
			tbl.Indexes = removeIndexByName(tbl.Indexes, id.Name)
		case DiffModify:
			if id.To != nil {
				for i, idx := range tbl.Indexes {
					if idx.Name == id.Name {
						updated := indexToProto(cat, rel.To, id.To)
						if updated.Comment == "" {
							updated.Comment = idx.Comment
						}
						if len(updated.Constraints) == 0 {
							updated.Constraints = idx.Constraints
						}
						tbl.Indexes[i] = updated
						break
					}
				}
			}
		default:
		}
	}
}

func applyConstraintDiffs(tbl *metadatapb.TableMetadata, catBefore, catAfter *Catalog, rel RelationDiffEntry) {
	for _, cd := range rel.Constraints {
		switch cd.Action {
		case DiffAdd:
			if cd.To != nil {
				addConstraintToTable(tbl, catAfter, rel, cd.To)
			}
		case DiffDrop:
			if cd.From != nil {
				removeConstraintFromTable(tbl, catBefore, cd.From)
			}
		case DiffModify:
			var prevComment string
			if cd.From != nil {
				prevComment = constraintCommentOnTable(tbl, catBefore, cd.From)
				removeConstraintFromTable(tbl, catBefore, cd.From)
			}
			if cd.To != nil {
				addConstraintToTable(tbl, catAfter, rel, cd.To)
				if prevComment != "" {
					setAddedConstraintCommentOnTable(tbl, catAfter, cd.To, prevComment)
				}
			}
		default:
		}
	}
}

func addConstraintToTable(tbl *metadatapb.TableMetadata, cat *Catalog, rel RelationDiffEntry, con *Constraint) {
	switch con.Type {
	case ConstraintPK, ConstraintUnique:
		idx := cat.GetIndexByOID(con.IndexOID)
		if idx != nil && rel.To != nil {
			tbl.Indexes = append(tbl.Indexes, indexToProto(cat, rel.To, idx))
		}
	case ConstraintFK:
		tbl.ForeignKeys = append(tbl.ForeignKeys, constraintToFKProto(cat, rel.To, con))
	case ConstraintCheck:
		tbl.CheckConstraints = append(tbl.CheckConstraints, &metadatapb.CheckConstraintMetadata{
			Name:       con.Name,
			Expression: con.CheckExpr,
		})
	case ConstraintExclude:
		tbl.ExcludeConstraints = append(tbl.ExcludeConstraints, &metadatapb.ExcludeConstraintMetadata{
			Name:       con.Name,
			Expression: con.CheckExpr,
		})
	default:
	}
}

func removeConstraintFromTable(tbl *metadatapb.TableMetadata, cat *Catalog, con *Constraint) {
	switch con.Type {
	case ConstraintPK, ConstraintUnique:
		// Constraint name and backing index name may differ; resolve via IndexOID.
		idxName := con.Name
		if idx := cat.GetIndexByOID(con.IndexOID); idx != nil {
			idxName = idx.Name
		}
		tbl.Indexes = removeIndexByName(tbl.Indexes, idxName)
	case ConstraintFK:
		out := make([]*metadatapb.ForeignKeyMetadata, 0, len(tbl.ForeignKeys))
		for _, fk := range tbl.ForeignKeys {
			if fk.Name != con.Name {
				out = append(out, fk)
			}
		}
		tbl.ForeignKeys = out
	case ConstraintCheck:
		out := make([]*metadatapb.CheckConstraintMetadata, 0, len(tbl.CheckConstraints))
		for _, c := range tbl.CheckConstraints {
			if c.Name != con.Name {
				out = append(out, c)
			}
		}
		tbl.CheckConstraints = out
	case ConstraintExclude:
		out := make([]*metadatapb.ExcludeConstraintMetadata, 0, len(tbl.ExcludeConstraints))
		for _, c := range tbl.ExcludeConstraints {
			if c.Name != con.Name {
				out = append(out, c)
			}
		}
		tbl.ExcludeConstraints = out
	default:
	}
}

// --- Conversion helpers: omni types → proto types ---

func columnToProto(cat *Catalog, col *Column) *metadatapb.ColumnMetadata {
	cm := &metadatapb.ColumnMetadata{
		Name:     col.Name,
		Position: int32(col.AttNum),
		Type:     cat.FormatType(col.TypeOID, col.TypeMod),
		Nullable: !col.NotNull,
		Default:  col.Default,
		Options:  slices.Clone(col.ForeignOptions),
	}
	if col.CollationName != "" {
		cm.Collation = col.CollationName
		cm.CollationSchema = columnCollationProtoSchema(cat, col)
	}
	if col.Generated == 's' {
		cm.Generation = &metadatapb.GenerationMetadata{
			Type:       metadatapb.GenerationMetadata_TYPE_STORED,
			Expression: col.GenerationExpr,
		}
	}
	if col.Identity != 0 {
		cm.IsIdentity = true
		switch col.Identity {
		case 'a':
			cm.IdentityGeneration = metadatapb.ColumnMetadata_ALWAYS
		case 'd':
			cm.IdentityGeneration = metadatapb.ColumnMetadata_BY_DEFAULT
		default:
		}
	}
	return cm
}

func indexToProto(cat *Catalog, rel *Relation, idx *Index) *metadatapb.IndexMetadata {
	im := &metadatapb.IndexMetadata{
		Name:         idx.Name,
		Type:         idx.AccessMethod,
		Unique:       idx.IsUnique,
		Primary:      idx.IsPrimary,
		IsConstraint: idx.ConstraintOID != 0,
	}
	if cat != nil && idx.ParentIndexOID != 0 {
		if parentIdx := cat.GetIndexByOID(idx.ParentIndexOID); parentIdx != nil {
			im.ParentIndexName = parentIdx.Name
			if parentRel := cat.GetRelationByOID(parentIdx.RelOID); parentRel != nil && parentRel.Schema != nil {
				im.ParentIndexSchema = parentRel.Schema.Name
			}
		}
	}

	nKey := idx.NKeyColumns
	if nKey <= 0 || nKey > len(idx.Columns) {
		nKey = len(idx.Columns)
	}

	exprIdx := 0
	for i := 0; i < nKey; i++ {
		attnum := idx.Columns[i]
		if attnum == 0 {
			if exprIdx < len(idx.Exprs) {
				im.Expressions = append(im.Expressions, idx.Exprs[exprIdx])
				exprIdx++
			}
		} else {
			colName := ""
			if rel != nil {
				for _, col := range rel.Columns {
					if col.AttNum == attnum {
						colName = col.Name
						break
					}
				}
			}
			im.Expressions = append(im.Expressions, colName)
		}

		if i < len(idx.IndOption) {
			im.Descending = append(im.Descending, idx.IndOption[i]&1 != 0)
		}
	}

	if len(idx.OpClasses) > 0 {
		im.OpclassNames = make([]string, len(idx.OpClasses))
		im.OpclassDefaults = make([]bool, len(idx.OpClasses))
		for i, oc := range idx.OpClasses {
			im.OpclassNames[i] = oc
			im.OpclassDefaults[i] = oc == ""
		}
	}

	if !idx.IsPrimary && idx.ConstraintOID == 0 && rel != nil && rel.Schema != nil &&
		(idx.WhereClause != "" || nKey < len(idx.Columns) || len(idx.OpClasses) > 0 || idx.NullsNotDistinct || hasUnmodeledIndexOptions(idx)) {
		im.Definition = buildCreateIndexSQL(idx, rel, rel.Schema.Name)
	}

	return im
}

func hasUnmodeledIndexOptions(idx *Index) bool {
	if idx == nil {
		return false
	}
	if slices.ContainsFunc(idx.Collations, func(parts []string) bool { return len(parts) > 0 }) {
		return true
	}
	for _, opt := range idx.IndOption {
		desc := opt&1 != 0
		nullsFirst := opt&2 != 0
		if desc != nullsFirst {
			return true
		}
	}
	return false
}

func columnCollationProtoSchema(cat *Catalog, col *Column) string {
	if col == nil || col.CollationName == "" {
		return ""
	}
	if cat != nil && col.UserCollationOID != 0 {
		for _, s := range cat.schemas {
			if s == nil {
				continue
			}
			if uc := s.Collations[col.CollationName]; uc != nil && uc.OID == col.UserCollationOID {
				if s.OID == PGCatalogNamespace {
					return ""
				}
				return s.Name
			}
		}
	}
	if col.CollationSchema != "" && col.CollationSchema != "pg_catalog" {
		return col.CollationSchema
	}
	return ""
}

func constraintToFKProto(cat *Catalog, rel *Relation, con *Constraint) *metadatapb.ForeignKeyMetadata {
	fk := &metadatapb.ForeignKeyMetadata{
		Name: con.Name,
	}

	for _, attnum := range con.Columns {
		if rel != nil {
			for _, col := range rel.Columns {
				if col.AttNum == attnum {
					fk.Columns = append(fk.Columns, col.Name)
					break
				}
			}
		}
	}

	fRel := cat.GetRelationByOID(con.FRelOID)
	if fRel != nil {
		if fRel.Schema != nil {
			fk.ReferencedSchema = fRel.Schema.Name
		}
		fk.ReferencedTable = fRel.Name
		for _, attnum := range con.FColumns {
			for _, col := range fRel.Columns {
				if col.AttNum == attnum {
					fk.ReferencedColumns = append(fk.ReferencedColumns, col.Name)
					break
				}
			}
		}
	}

	fk.OnUpdate = wtFKActionToString(con.FKUpdAction)
	fk.OnDelete = wtFKActionToString(con.FKDelAction)
	fk.MatchType = wtFKMatchToString(con.FKMatchType)
	return fk
}

func sequenceToProto(cat *Catalog, seq *Sequence) *metadatapb.SequenceMetadata {
	sm := &metadatapb.SequenceMetadata{
		Name:      seq.Name,
		DataType:  cat.FormatType(seq.TypeOID, -1),
		Start:     fmt.Sprintf("%d", seq.Start),
		MinValue:  fmt.Sprintf("%d", seq.MinValue),
		MaxValue:  fmt.Sprintf("%d", seq.MaxValue),
		Increment: fmt.Sprintf("%d", seq.Increment),
		Cycle:     seq.Cycle,
		CacheSize: fmt.Sprintf("%d", seq.CacheValue),
	}
	if seq.OwnerRelOID != 0 {
		if ownerRel := cat.GetRelationByOID(seq.OwnerRelOID); ownerRel != nil {
			sm.OwnerTable = ownerRel.Name
			for _, col := range ownerRel.Columns {
				if col.AttNum == seq.OwnerAttNum {
					sm.OwnerColumn = col.Name
					break
				}
			}
		}
	}
	return sm
}

func wtFKActionToString(action byte) string {
	switch action {
	case 'r':
		return "RESTRICT"
	case 'c':
		return "CASCADE"
	case 'n':
		return "SET NULL"
	case 'd':
		return "SET DEFAULT"
	default:
		return "NO ACTION"
	}
}

func wtFKMatchToString(match byte) string {
	switch match {
	case 'f':
		return "FULL"
	case 'p':
		return "PARTIAL"
	default:
		return "SIMPLE"
	}
}

// --- Slice helpers ---

func findTable(sm *metadatapb.SchemaMetadata, name string) *metadatapb.TableMetadata {
	for _, t := range sm.Tables {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func removeTableByName(tables []*metadatapb.TableMetadata, name string) []*metadatapb.TableMetadata {
	out := make([]*metadatapb.TableMetadata, 0, len(tables))
	for _, t := range tables {
		if t.Name != name {
			out = append(out, t)
		}
	}
	return out
}

func removeViewByName(views []*metadatapb.ViewMetadata, name string) []*metadatapb.ViewMetadata {
	out := make([]*metadatapb.ViewMetadata, 0, len(views))
	for _, v := range views {
		if v.Name != name {
			out = append(out, v)
		}
	}
	return out
}

func removeMatViewByName(mvs []*metadatapb.MaterializedViewMetadata, name string) []*metadatapb.MaterializedViewMetadata {
	out := make([]*metadatapb.MaterializedViewMetadata, 0, len(mvs))
	for _, m := range mvs {
		if m.Name != name {
			out = append(out, m)
		}
	}
	return out
}

func removeSequenceByName(seqs []*metadatapb.SequenceMetadata, name string) []*metadatapb.SequenceMetadata {
	out := make([]*metadatapb.SequenceMetadata, 0, len(seqs))
	for _, s := range seqs {
		if s.Name != name {
			out = append(out, s)
		}
	}
	return out
}

func removeColumnByName(cols []*metadatapb.ColumnMetadata, name string) []*metadatapb.ColumnMetadata {
	out := make([]*metadatapb.ColumnMetadata, 0, len(cols))
	for _, c := range cols {
		if c.Name != name {
			out = append(out, c)
		}
	}
	return out
}

func removeIndexByName(indexes []*metadatapb.IndexMetadata, name string) []*metadatapb.IndexMetadata {
	out := make([]*metadatapb.IndexMetadata, 0, len(indexes))
	for _, i := range indexes {
		if i.Name != name {
			out = append(out, i)
		}
	}
	return out
}

func removeEnumByName(enums []*metadatapb.EnumTypeMetadata, name string) []*metadatapb.EnumTypeMetadata {
	out := make([]*metadatapb.EnumTypeMetadata, 0, len(enums))
	for _, e := range enums {
		if e.Name != name {
			out = append(out, e)
		}
	}
	return out
}

func functionToProto(cat *Catalog, up *UserProc, identity string) *metadatapb.FunctionMetadata {
	return &metadatapb.FunctionMetadata{
		Name:       up.Name,
		Definition: buildUserProcDDL(cat, up),
		Signature:  identity,
	}
}

func procedureToProto(cat *Catalog, up *UserProc, identity string) *metadatapb.ProcedureMetadata {
	return &metadatapb.ProcedureMetadata{
		Name:       up.Name,
		Definition: buildUserProcDDL(cat, up),
		Signature:  identity,
	}
}

func buildUserProcDDL(cat *Catalog, up *UserProc) string {
	var b strings.Builder
	if up.Kind == 'p' {
		b.WriteString("CREATE OR REPLACE PROCEDURE ")
	} else {
		b.WriteString("CREATE OR REPLACE FUNCTION ")
	}
	if up.Schema != nil && up.Schema.Name != "" {
		b.WriteString(quoteWalkThroughIdent(up.Schema.Name))
		b.WriteByte('.')
	}
	b.WriteString(quoteWalkThroughIdent(up.Name))
	b.WriteByte('(')
	b.WriteString(formatFuncParams(cat, up))
	b.WriteString(")\n")
	if up.Kind != 'p' {
		b.WriteString(" RETURNS ")
		if hasTableReturn(up) {
			b.WriteString(formatReturnsTable(cat, up))
		} else {
			if up.RetSet {
				b.WriteString("SETOF ")
			}
			b.WriteString(cat.FormatType(up.RetType, -1))
		}
		b.WriteByte('\n')
	}
	b.WriteString(" LANGUAGE ")
	b.WriteString(up.Language)
	b.WriteByte('\n')
	switch up.Volatile {
	case 'i':
		b.WriteString(" IMMUTABLE\n")
	case 's':
		b.WriteString(" STABLE\n")
	default:
	}
	if up.IsStrict {
		b.WriteString(" STRICT\n")
	}
	if up.SecDef {
		b.WriteString(" SECURITY DEFINER\n")
	}
	if up.LeakProof {
		b.WriteString(" LEAKPROOF\n")
	}
	switch up.Parallel {
	case 's':
		b.WriteString(" PARALLEL SAFE\n")
	case 'r':
		b.WriteString(" PARALLEL RESTRICTED\n")
	default:
	}
	if up.IsWindow {
		b.WriteString(" WINDOW\n")
	}
	if up.Cost != "" {
		b.WriteString(" COST ")
		b.WriteString(up.Cost)
		b.WriteByte('\n')
	}
	if up.Rows != "" {
		b.WriteString(" ROWS ")
		b.WriteString(up.Rows)
		b.WriteByte('\n')
	}
	for _, sc := range up.SetConfig {
		b.WriteByte(' ')
		b.WriteString(sc)
		b.WriteByte('\n')
	}
	if up.HasSqlBody {
		b.WriteString(up.Body)
		b.WriteByte('\n')
	} else {
		b.WriteString("AS ")
		b.WriteString(dollarQuoteWithDefault(up.Body, "$function$"))
		b.WriteByte('\n')
	}
	return b.String()
}

func isTempRelation(schemaName string, rel *Relation) bool {
	if isTempSchemaName(schemaName) {
		return true
	}
	if rel == nil {
		return false
	}
	return rel.Persistence == 't' || (rel.Schema != nil && isTempSchemaName(rel.Schema.Name))
}

func isTempSequenceDiff(seq SequenceDiffEntry, canonicalSchema func(string) string) bool {
	if isTempSchemaName(seq.SchemaName) || isTempSchemaName(canonicalSchema(seq.SchemaName)) {
		return true
	}
	if seq.From != nil && (seq.From.Persistence == 't' || (seq.From.Schema != nil && isTempSchemaName(seq.From.Schema.Name))) {
		return true
	}
	if seq.To != nil && (seq.To.Persistence == 't' || (seq.To.Schema != nil && isTempSchemaName(seq.To.Schema.Name))) {
		return true
	}
	return false
}

func extractFunctionDependencyTables(cat *Catalog, up *UserProc) []*metadatapb.DependencyTable {
	if cat == nil || up == nil || (!up.HasSqlBody && !strings.EqualFold(up.Language, "sql")) {
		return nil
	}
	sqlText := functionSQLBodyForDeps(up)
	if sqlText == "" {
		return nil
	}
	list, err := pgparser.Parse(sqlText)
	if err != nil || list == nil {
		return nil
	}
	cteNames := make(map[string]bool)
	nodes.Inspect(list, func(n nodes.Node) bool {
		if cte, ok := n.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
			cteNames[cte.Ctename] = true
		}
		return true
	})
	seen := make(map[string]bool)
	var deps []*metadatapb.DependencyTable
	nodes.Inspect(list, func(n nodes.Node) bool {
		rv, ok := n.(*nodes.RangeVar)
		if !ok || rv.Relname == "" || (rv.Schemaname == "" && cteNames[rv.Relname]) {
			return true
		}
		var rel *Relation
		if rv.Schemaname != "" {
			rel = cat.GetRelation(rv.Schemaname, rv.Relname)
		} else {
			rel = cat.GetRelation("", rv.Relname)
		}
		if rel == nil || rel.Persistence == 't' || rel.Schema == nil {
			return true
		}
		if rel.Schema.OID == PGCatalogNamespace || rel.Schema.OID == PGToastNamespace || isTempSchemaName(rel.Schema.Name) {
			return true
		}
		key := rel.Schema.Name + "." + rel.Name
		if !seen[key] {
			seen[key] = true
			deps = append(deps, &metadatapb.DependencyTable{
				Schema: rel.Schema.Name,
				Table:  rel.Name,
			})
		}
		return true
	})
	slices.SortFunc(deps, func(a, b *metadatapb.DependencyTable) int {
		if a.Schema != b.Schema {
			return strings.Compare(a.Schema, b.Schema)
		}
		return strings.Compare(a.Table, b.Table)
	})
	return deps
}

func functionSQLBodyForDeps(up *UserProc) string {
	trimmed := strings.TrimSpace(up.Body)
	if trimmed == "" {
		return ""
	}
	if !up.HasSqlBody {
		return trimmed
	}
	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "RETURN ") {
		return "SELECT " + strings.TrimSpace(trimmed[len("RETURN "):])
	}
	if strings.HasPrefix(upper, "BEGIN ATOMIC") && strings.HasSuffix(upper, "END") {
		inner := strings.TrimSpace(trimmed[len("BEGIN ATOMIC") : len(trimmed)-len("END")])
		// Convert any RETURN <expr> lines inside BEGIN ATOMIC into SELECT <expr> so pgparser parses them.
		var rewritten []string
		for _, stmt := range strings.Split(inner, ";") {
			s := strings.TrimSpace(stmt)
			if s == "" {
				continue
			}
			if strings.HasPrefix(strings.ToUpper(s), "RETURN ") {
				s = "SELECT " + strings.TrimSpace(s[len("RETURN "):])
			}
			rewritten = append(rewritten, s)
		}
		return strings.Join(rewritten, "; ")
	}
	return trimmed
}

func popProcedureMatchingIdentity(procs []*metadatapb.ProcedureMetadata, schemaName, identity string) (*metadatapb.ProcedureMetadata, []*metadatapb.ProcedureMetadata) {
	popped, rest := popProcedureByIdentity(procs, identity)
	if popped == nil && schemaName != "" {
		if shortID := strings.TrimPrefix(identity, schemaName+"."); shortID != identity {
			popped, rest = popProcedureByIdentity(rest, shortID)
		}
	}
	return popped, rest
}

func popFunctionMatchingIdentity(funcs []*metadatapb.FunctionMetadata, schemaName, identity string) (*metadatapb.FunctionMetadata, []*metadatapb.FunctionMetadata) {
	popped, rest := popFunctionByIdentity(funcs, identity)
	if popped == nil && schemaName != "" {
		if shortID := strings.TrimPrefix(identity, schemaName+"."); shortID != identity {
			popped, rest = popFunctionByIdentity(rest, shortID)
		}
	}
	return popped, rest
}

func popProcedureByIdentity(procs []*metadatapb.ProcedureMetadata, identity string) (*metadatapb.ProcedureMetadata, []*metadatapb.ProcedureMetadata) {
	var matched *metadatapb.ProcedureMetadata
	out := make([]*metadatapb.ProcedureMetadata, 0, len(procs))
	for _, p := range procs {
		if p.Signature == identity {
			matched = p
			continue
		}
		out = append(out, p)
	}
	return matched, out
}

func popFunctionByIdentity(funcs []*metadatapb.FunctionMetadata, identity string) (*metadatapb.FunctionMetadata, []*metadatapb.FunctionMetadata) {
	var matched *metadatapb.FunctionMetadata
	out := make([]*metadatapb.FunctionMetadata, 0, len(funcs))
	for _, f := range funcs {
		if f.Signature == identity {
			matched = f
			continue
		}
		out = append(out, f)
	}
	return matched, out
}

// compositeTypeToProto builds metadata from the catalog relation, preserving
// the type comment and per-attribute comments/collations from the previous
// metadata (the catalog does not track comments, and it only keeps the bare
// collation name) for attributes that survive by name. Attributes unchanged
// between from and to keep their previous metadata verbatim, so a composite
// degraded to a text-backed stand-in does not rewrite untouched
// attributes to text — mirroring the per-column granularity of table diffs.
//
// Within a degraded composite this is a deliberate trade-off: an attribute
// genuinely altered TO text is indistinguishable from an unchanged one (both
// read text in from and to), so that rare change is missed in favor of not
// corrupting every untouched attribute. Walk-through metadata is advisory;
// post-execution sync restores ground truth. The root fix is modeling the
// types LoadMetadata cannot install today (e.g. domains).
func compositeTypeToProto(cat *Catalog, name string, from, to *Relation, previous *metadatapb.CompositeTypeMetadata) *metadatapb.CompositeTypeMetadata {
	previousAttributes := make(map[string]*metadatapb.CompositeTypeAttribute)
	composite := &metadatapb.CompositeTypeMetadata{Name: name}
	if previous != nil {
		composite.Comment = previous.Comment
		composite.SkipDump = previous.SkipDump
		for _, attribute := range previous.Attributes {
			previousAttributes[attribute.Name] = attribute
		}
	}
	fromByAttNum := make(map[int16]*Column)
	if from != nil {
		for _, c := range from.Columns {
			if c != nil {
				fromByAttNum[c.AttNum] = c
			}
		}
	}
	for _, c := range to.Columns {
		if c == nil || c.Name == "" {
			continue
		}
		// The attribute number is the stable column identity: it survives
		// RENAME ATTRIBUTE and is never reused, so a dropped-and-re-added
		// attribute (same name, new attnum) correctly reads as a new column
		// while a renamed survivor still carries its previous metadata.
		fc := fromByAttNum[c.AttNum]
		if fc != nil &&
			fc.TypeOID == c.TypeOID && fc.TypeMod == c.TypeMod && fc.CollationName == c.CollationName &&
			columnCollationProtoSchema(cat, fc) == columnCollationProtoSchema(cat, c) {
			if prev := previousAttributes[fc.Name]; prev != nil {
				composite.Attributes = append(composite.Attributes, &metadatapb.CompositeTypeAttribute{
					Name:      c.Name,
					Type:      prev.Type,
					Collation: prev.Collation,
					Comment:   prev.Comment,
				})
				continue
			}
		}
		attribute := &metadatapb.CompositeTypeAttribute{
			Name: c.Name,
			Type: qualifyWalkThroughAttributeType(cat, c.TypeOID, cat.FormatType(c.TypeOID, c.TypeMod)),
		}
		// Previous metadata is only consulted for the column identity that
		// survived (same attnum): a dropped-and-re-added attribute must not
		// inherit the old attribute's collation reference or comment.
		var surviving *metadatapb.CompositeTypeAttribute
		if fc != nil {
			surviving = previousAttributes[fc.Name]
		}
		if c.CollationName != "" {
			collSchema := columnCollationProtoSchema(cat, c)
			if collSchema != "" {
				attribute.Collation = migrationQualifiedName(collSchema, c.CollationName)
			} else {
				attribute.Collation = quoteIdentAlways(c.CollationName)
			}
			// Substitute the previous emit-ready (possibly schema-qualified)
			// reference only when it names the same collation in the same schema —
			// a changed collation must reflect the catalog's current value.
			if surviving != nil && collationReferenceBareName(surviving.Collation) == c.CollationName &&
				collationReferenceSchema(surviving.Collation) == collSchema {
				attribute.Collation = surviving.Collation
			}
		}
		if surviving != nil {
			attribute.Comment = surviving.Comment
		}
		composite.Attributes = append(composite.Attributes, attribute)
	}
	return composite
}

func removeCompositeTypeByName(composites []*metadatapb.CompositeTypeMetadata, name string) []*metadatapb.CompositeTypeMetadata {
	out := make([]*metadatapb.CompositeTypeMetadata, 0, len(composites))
	for _, composite := range composites {
		if composite.Name != name {
			out = append(out, composite)
		}
	}
	return out
}

// collationReferenceBareName extracts the bare collation name from a stored
// emit-ready reference (e.g. `"C"` -> C, `locale.en_us` -> en_us) for
// comparison with the catalog's collation name. The comparison cannot see
// schema qualifiers — the omni catalog only records the bare collation name —
// so switching an attribute between two same-named collations in different
// schemas keeps the previous reference. Known, accepted limitation.
func collationReferenceBareName(reference string) string {
	if _, name, ok := parseQualifiedTypeIdent(reference); ok {
		return name
	}
	name := strings.Trim(reference, `"`)
	return strings.ReplaceAll(name, `""`, `"`)
}

func collationReferenceSchema(reference string) string {
	if schema, _, ok := parseQualifiedTypeIdent(reference); ok && schema != "pg_catalog" {
		return schema
	}
	return ""
}

// qualifyWalkThroughAttributeType enforces the CompositeTypeAttribute.Type
// contract: user-defined types are always schema-qualified. FormatType
// renders search-path-visible types unqualified, so the (element) type's
// namespace is resolved and prepended when missing. System-schema types come
// back unqualified from QueryPgNamespace and stay untouched.
func qualifyWalkThroughAttributeType(cat *Catalog, typeOID uint32, formatted string) string {
	t := cat.TypeByOID(typeOID)
	if t == nil {
		return formatted
	}
	element := t
	if t.Elem != 0 && t.Len == -1 {
		if el := cat.TypeByOID(t.Elem); el != nil {
			element = el
		}
	}
	schemaName := ""
	for _, row := range cat.QueryPgNamespace() {
		if row.OID == element.Namespace {
			schemaName = row.NspName
			break
		}
	}
	// QueryPgNamespace already omits builtin-OID schemas; the explicit check
	// keeps the invariant independent of the catalog's OID layout and covers
	// information_schema, mirroring the sync query's exclusions.
	if schemaName == "" || wtIsSystemSchema(schemaName) {
		return formatted
	}
	prefix := quoteWalkThroughIdent(schemaName) + "."
	quotedPrefix := `"` + strings.ReplaceAll(schemaName, `"`, `""`) + `".`
	if strings.HasPrefix(formatted, prefix) || strings.HasPrefix(formatted, quotedPrefix) {
		return formatted
	}
	return prefix + formatted
}

// quoteWalkThroughIdent mirrors quote_ident: quote when the name needs it,
// including reserved keywords.
func quoteWalkThroughIdent(name string) string {
	simple := name != ""
	for i, r := range name {
		if (r >= 'a' && r <= 'z') || r == '_' || (i > 0 && ((r >= '0' && r <= '9') || r == '$')) {
			continue
		}
		simple = false
		break
	}
	if simple && !pgparser.IsReservedKeyword(name) {
		return name
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func wtIsSystemSchema(s string) bool {
	switch s {
	case "pg_catalog", "pg_toast", "information_schema":
		return true
	}
	return strings.HasPrefix(s, "pg_")
}

// relationToTableProto converts an omni Relation to a metadatapb.TableMetadata.
func relationToTableProto(c *Catalog, rel *Relation) *metadatapb.TableMetadata {
	table := &metadatapb.TableMetadata{
		Name:  rel.Name,
		Owner: rel.OwnerName,
	}

	for _, col := range rel.Columns {
		table.Columns = append(table.Columns, columnToProto(c, col))
	}

	for _, idx := range c.IndexesOf(rel.OID) {
		table.Indexes = append(table.Indexes, indexToProto(c, rel, idx))
	}

	for _, con := range c.ConstraintsOf(rel.OID) {
		switch con.Type {
		case 'f':
			table.ForeignKeys = append(table.ForeignKeys, constraintToFKProto(c, rel, con))
		case 'c':
			table.CheckConstraints = append(table.CheckConstraints, &metadatapb.CheckConstraintMetadata{
				Name:       con.Name,
				Expression: con.CheckExpr,
			})
		case 'x':
			table.ExcludeConstraints = append(table.ExcludeConstraints, &metadatapb.ExcludeConstraintMetadata{
				Name:       con.Name,
				Expression: con.CheckExpr,
			})
		default:
		}
	}

	if rel.Schema != nil {
		for _, trig := range c.TriggersOf(rel.OID) {
			if trig.ParentTriggerOID == 0 {
				table.Triggers = append(table.Triggers, triggerToProto(c, rel.Schema.Name, rel.Name, trig))
			}
		}
		table.Rules = rulesToProto(c, rel.Schema.Name, rel)
	}

	return table
}

// relationToViewProto converts an omni Relation (view) to a metadatapb.ViewMetadata.
func relationToViewProto(c *Catalog, rel *Relation) *metadatapb.ViewMetadata {
	view := &metadatapb.ViewMetadata{
		Name:              rel.Name,
		DependencyColumns: relationDependencyColumns(c, rel),
	}
	for _, col := range rel.Columns {
		view.Columns = append(view.Columns, columnToProto(c, col))
	}
	if rel.Schema != nil {
		if def, err := c.GetViewDefinition(rel.Schema.Name, rel.Name); err == nil {
			view.Definition = def
		}
		for _, trig := range c.TriggersOf(rel.OID) {
			if trig.ParentTriggerOID == 0 {
				view.Triggers = append(view.Triggers, triggerToProto(c, rel.Schema.Name, rel.Name, trig))
			}
		}
		view.Rules = rulesToProto(c, rel.Schema.Name, rel)
	}
	return view
}

// relationToMatViewProto converts an omni Relation (materialized view) to a metadatapb.MaterializedViewMetadata.
func relationToMatViewProto(c *Catalog, rel *Relation) *metadatapb.MaterializedViewMetadata {
	mv := &metadatapb.MaterializedViewMetadata{
		Name:              rel.Name,
		DependencyColumns: relationDependencyColumns(c, rel),
	}
	if rel.Schema != nil {
		if def, err := c.GetMatViewDefinition(rel.Schema.Name, rel.Name); err == nil {
			mv.Definition = def
		}
		for _, trig := range c.TriggersOf(rel.OID) {
			if trig.ParentTriggerOID == 0 {
				mv.Triggers = append(mv.Triggers, triggerToProto(c, rel.Schema.Name, rel.Name, trig))
			}
		}
	}

	for _, idx := range c.IndexesOf(rel.OID) {
		mv.Indexes = append(mv.Indexes, indexToProto(c, rel, idx))
	}

	return mv
}

func applyTriggerDiffs(triggers *[]*metadatapb.TriggerMetadata, cat *Catalog, rel RelationDiffEntry) {
	if len(rel.Triggers) == 0 {
		return
	}
	schemaName := rel.SchemaName
	tableName := rel.Name
	upsert := func(trig *Trigger) {
		updated := triggerToProto(cat, schemaName, tableName, trig)
		for i, existing := range *triggers {
			if existing.Name == trig.Name {
				updated.Comment = existing.Comment
				(*triggers)[i] = updated
				return
			}
		}
		*triggers = append(*triggers, updated)
	}
	for _, td := range rel.Triggers {
		switch td.Action {
		case DiffAdd, DiffModify:
			if td.To != nil {
				upsert(td.To)
			} else if td.Action == DiffModify {
				*triggers = removeTriggerByName(*triggers, td.Name)
			}
		case DiffDrop:
			*triggers = removeTriggerByName(*triggers, td.Name)
		default:
		}
	}
}

func triggerToProto(cat *Catalog, schemaName, tableName string, trig *Trigger) *metadatapb.TriggerMetadata {
	timing := ""
	switch trig.Timing {
	case TriggerBefore:
		timing = "BEFORE"
	case TriggerAfter:
		timing = "AFTER"
	case TriggerInsteadOf:
		timing = "INSTEAD OF"
	}
	var events []string
	if trig.Events&TriggerEventInsert != 0 {
		events = append(events, "INSERT")
	}
	if trig.Events&TriggerEventUpdate != 0 {
		events = append(events, "UPDATE")
	}
	if trig.Events&TriggerEventDelete != 0 {
		events = append(events, "DELETE")
	}
	if trig.Events&TriggerEventTruncate != 0 {
		events = append(events, "TRUNCATE")
	}
	return &metadatapb.TriggerMetadata{
		Name:   trig.Name,
		Timing: timing,
		Event:  strings.Join(events, " OR "),
		Body:   buildCreateTriggerOp(cat, schemaName, tableName, trig).SQL,
	}
}

func removeTriggerByName(triggers []*metadatapb.TriggerMetadata, name string) []*metadatapb.TriggerMetadata {
	out := make([]*metadatapb.TriggerMetadata, 0, len(triggers))
	for _, tr := range triggers {
		if tr.Name != name {
			out = append(out, tr)
		}
	}
	return out
}

func rangeTypeToProto(cat *Catalog, rt *RangeType, previous *metadatapb.RangeTypeMetadata) *metadatapb.RangeTypeMetadata {
	name := rt.Name
	if bt := cat.TypeByOID(rt.OID); bt != nil && bt.TypeName != "" {
		name = bt.TypeName
	}
	out := &metadatapb.RangeTypeMetadata{
		Name:           name,
		Subtype:        qualifyWalkThroughAttributeType(cat, rt.SubTypeOID, cat.FormatType(rt.SubTypeOID, -1)),
		SubtypeOpclass: rt.SubtypeOpclass,
		Collation:      rangeCollationProto(cat, rt),
		Canonical:      rt.Canonical,
		SubtypeDiff:    rt.SubtypeDiff,
	}
	if rt.MultirangeOID != 0 {
		out.MultirangeName = qualifyWalkThroughAttributeType(cat, rt.MultirangeOID, cat.FormatType(rt.MultirangeOID, -1))
	}
	if previous != nil {
		if out.SubtypeOpclass == "" {
			out.SubtypeOpclass = previous.SubtypeOpclass
		}
		if out.Collation == "" {
			out.Collation = previous.Collation
		}
		if out.Canonical == "" {
			out.Canonical = previous.Canonical
		}
		if out.SubtypeDiff == "" {
			out.SubtypeDiff = previous.SubtypeDiff
		}
		out.Comment = previous.Comment
		out.SkipDump = previous.SkipDump
	}
	return out
}

func rangeCollationProto(cat *Catalog, rt *RangeType) string {
	if rt == nil || rt.CollationName == "" {
		return ""
	}
	collSchema := rt.CollationSchema
	collName := rt.CollationName
	if cat != nil && rt.UserCollationOID != 0 {
		if sc, uc := collationByOID(cat, rt.UserCollationOID); uc != nil {
			collSchema = sc
			collName = uc.Name
		}
	}
	if collSchema != "" && collSchema != "pg_catalog" {
		return migrationQualifiedName(collSchema, collName)
	}
	return quoteIdentAlways(collName)
}

func popRangeTypeByName(ranges []*metadatapb.RangeTypeMetadata, name string) (*metadatapb.RangeTypeMetadata, []*metadatapb.RangeTypeMetadata) {
	var popped *metadatapb.RangeTypeMetadata
	out := make([]*metadatapb.RangeTypeMetadata, 0, len(ranges))
	for _, r := range ranges {
		if r.Name == name && popped == nil {
			popped = r
			continue
		}
		out = append(out, r)
	}
	return popped, out
}

func commentParts2(c CommentDiffEntry) (string, string, bool) {
	if len(c.objParts) >= 2 {
		return c.objParts[0], c.objParts[1], true
	}
	return strings.Cut(c.ObjDescription, ".")
}

func commentParts3(c CommentDiffEntry) (string, string, string, bool) {
	if len(c.objParts) >= 3 {
		return c.objParts[0], c.objParts[1], c.objParts[2], true
	}
	parts := strings.SplitN(c.ObjDescription, ".", 3)
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func applyCommentDiffs(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, comments []CommentDiffEntry) {
	for _, c := range comments {
		switch c.ObjType {
		case 'n':
			schemaName := c.ObjDescription
			if len(c.objParts) > 0 {
				schemaName = c.objParts[0]
			}
			if sm := findSchema(result, schemaName); sm != nil {
				sm.Comment = c.To
			}
		case 'r':
			schemaName, relName, ok := commentParts2(c)
			if !ok {
				continue
			}
			sm := findSchema(result, schemaName)
			if c.SubID != 0 {
				if sm == nil {
					continue
				}
				colName := resolveCommentColumnName(catAfter, catBefore, schemaName, relName, c.SubID)
				if colName == "" {
					continue
				}
				if tbl := findTable(sm, relName); tbl != nil {
					setColumnComment(tbl.Columns, colName, c.To)
				} else if ext := findExternalTable(sm, relName); ext != nil {
					setColumnComment(ext.Columns, colName, c.To)
				} else if v := findView(sm, relName); v != nil {
					setColumnComment(v.Columns, colName, c.To)
				} else {
					for _, ct := range sm.CompositeTypes {
						if ct.Name == relName {
							for _, attr := range ct.Attributes {
								if attr.Name == colName {
									attr.Comment = c.To
									break
								}
							}
							break
						}
					}
				}
				continue
			}
			if sm != nil {
				if tbl := findTable(sm, relName); tbl != nil {
					tbl.Comment = c.To
					continue
				}
			}
			if p := findPartitionInSchema(result, schemaName, relName); p != nil {
				p.Comment = c.To
			} else if sm != nil {
				for _, v := range sm.Views {
					if v.Name == relName {
						v.Comment = c.To
						break
					}
				}
				for _, mv := range sm.MaterializedViews {
					if mv.Name == relName {
						mv.Comment = c.To
						break
					}
				}
			}
		case 'i':
			schemaName, idxName, ok := commentParts2(c)
			if !ok {
				continue
			}
			setIndexComment(result, schemaName, idxName, c.To)
		case 't':
			schemaName, typeName, ok := commentParts2(c)
			if !ok {
				continue
			}
			sm := findSchema(result, schemaName)
			if sm == nil {
				continue
			}
			if c.SubID != 0 {
				attrName := resolveCommentColumnName(catAfter, catBefore, schemaName, typeName, c.SubID)
				if attrName == "" {
					continue
				}
				for _, ct := range sm.CompositeTypes {
					if ct.Name == typeName {
						for _, attr := range ct.Attributes {
							if attr.Name == attrName {
								attr.Comment = c.To
								break
							}
						}
						break
					}
				}
				continue
			}
			for _, ct := range sm.CompositeTypes {
				if ct.Name == typeName {
					ct.Comment = c.To
					break
				}
			}
			for _, et := range sm.EnumTypes {
				if et.Name == typeName {
					et.Comment = c.To
					break
				}
			}
			for _, rt := range sm.RangeTypes {
				if rt.Name == typeName {
					rt.Comment = c.To
					break
				}
			}
		case 's':
			schemaName, seqName, ok := commentParts2(c)
			if !ok {
				continue
			}
			if sm := findSchema(result, schemaName); sm != nil {
				for _, seq := range sm.Sequences {
					if seq.Name == seqName {
						seq.Comment = c.To
						break
					}
				}
			}
		case 'c':
			schemaName, tableName, conName, ok := commentParts3(c)
			if !ok {
				continue
			}
			if c.Action == DiffDrop && isRecreatedConstraint(catBefore, catAfter, schemaName, tableName, conName) {
				continue
			}
			sm := findSchema(result, schemaName)
			if sm == nil {
				continue
			}
			if tbl := findTable(sm, tableName); tbl != nil {
				setConstraintCommentOnTable(tbl, conName, c.To)
			} else if ext := findExternalTable(sm, tableName); ext != nil {
				for _, chk := range ext.CheckConstraints {
					if chk.Name == conName {
						chk.Comment = c.To
						break
					}
				}
			}
		case 'g':
			schemaName, tableName, trigName, ok := commentParts3(c)
			if !ok {
				continue
			}
			sm := findSchema(result, schemaName)
			if sm == nil {
				continue
			}
			setTriggerCommentInSchema(sm, tableName, trigName, c.To)
		case 'f':
			schemaName, identity, ok := commentParts2(c)
			if !ok {
				continue
			}
			if sm := findSchema(result, schemaName); sm != nil {
				for _, fn := range sm.Functions {
					if fn.Signature == identity || fn.Signature == c.ObjDescription {
						fn.Comment = c.To
						break
					}
				}
				for _, p := range sm.Procedures {
					if p.Signature == identity || p.Signature == c.ObjDescription {
						p.Comment = c.To
						break
					}
				}
			}
		case 'E':
			etName := c.ObjDescription
			if len(c.objParts) > 0 {
				etName = c.objParts[0]
			}
			for _, et := range result.EventTriggers {
				if et.Name == etName {
					et.Comment = c.To
					break
				}
			}
		default:
		}
	}
}

func resolveCommentColumnName(catAfter, catBefore *Catalog, schemaName, relName string, attNum int16) string {
	for _, cat := range []*Catalog{catAfter, catBefore} {
		if cat == nil {
			continue
		}
		if rel := cat.GetRelation(schemaName, relName); rel != nil {
			for _, col := range rel.Columns {
				if col.AttNum == attNum {
					return col.Name
				}
			}
		}
	}
	return ""
}

func setIndexComment(meta *metadatapb.DatabaseSchemaMetadata, targetSchema, idxName, comment string) {
	var setInPartitions func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) bool
	setInPartitions = func(parentSchema string, partitions []*metadatapb.TablePartitionMetadata) bool {
		for _, p := range partitions {
			partSchema := partitionSchemaName(parentSchema, p)
			if partSchema == targetSchema {
				for _, idx := range p.Indexes {
					if idx.Name == idxName {
						idx.Comment = comment
						return true
					}
				}
			}
			if setInPartitions(partSchema, p.Subpartitions) {
				return true
			}
		}
		return false
	}
	if sm := findSchema(meta, targetSchema); sm != nil {
		for _, tbl := range sm.Tables {
			for _, idx := range tbl.Indexes {
				if idx.Name == idxName {
					idx.Comment = comment
					return
				}
			}
		}
		for _, mv := range sm.MaterializedViews {
			for _, idx := range mv.Indexes {
				if idx.Name == idxName {
					idx.Comment = comment
					return
				}
			}
		}
	}
	for _, s := range meta.Schemas {
		for _, tbl := range s.Tables {
			if setInPartitions(s.Name, tbl.Partitions) {
				return
			}
		}
	}
}

func constraintCommentSlot(tbl *metadatapb.TableMetadata, cat *Catalog, con *Constraint) *string {
	switch con.Type {
	case ConstraintPK, ConstraintUnique:
		idxName := con.Name
		if idx := cat.GetIndexByOID(con.IndexOID); idx != nil {
			idxName = idx.Name
		}
		for _, idx := range tbl.Indexes {
			if idx.Name == idxName {
				return &idx.Comment
			}
		}
	case ConstraintFK:
		for _, fk := range tbl.ForeignKeys {
			if fk.Name == con.Name {
				return &fk.Comment
			}
		}
	case ConstraintCheck:
		for _, chk := range tbl.CheckConstraints {
			if chk.Name == con.Name {
				return &chk.Comment
			}
		}
	default:
	}
	return nil
}

func constraintCommentOnTable(tbl *metadatapb.TableMetadata, cat *Catalog, con *Constraint) string {
	if slot := constraintCommentSlot(tbl, cat, con); slot != nil {
		return *slot
	}
	return ""
}

func setAddedConstraintCommentOnTable(tbl *metadatapb.TableMetadata, cat *Catalog, con *Constraint, comment string) {
	if slot := constraintCommentSlot(tbl, cat, con); slot != nil {
		*slot = comment
	}
}

func setConstraintCommentOnTable(tbl *metadatapb.TableMetadata, conName, comment string) {
	for _, fk := range tbl.ForeignKeys {
		if fk.Name == conName {
			fk.Comment = comment
			return
		}
	}
	for _, chk := range tbl.CheckConstraints {
		if chk.Name == conName {
			chk.Comment = comment
			return
		}
	}
	for _, idx := range tbl.Indexes {
		if idx.Name == conName && idx.IsConstraint {
			idx.Comment = comment
			return
		}
	}
}

func setTriggerCommentInList(triggers []*metadatapb.TriggerMetadata, trigName, comment string) bool {
	for _, tr := range triggers {
		if tr.Name == trigName {
			tr.Comment = comment
			return true
		}
	}
	return false
}

func setTriggerCommentInSchema(sm *metadatapb.SchemaMetadata, relName, trigName, comment string) {
	if tbl := findTable(sm, relName); tbl != nil && setTriggerCommentInList(tbl.Triggers, trigName, comment) {
		return
	}
	if ext := findExternalTable(sm, relName); ext != nil && setTriggerCommentInList(ext.Triggers, trigName, comment) {
		return
	}
	for _, v := range sm.Views {
		if v.Name == relName && setTriggerCommentInList(v.Triggers, trigName, comment) {
			return
		}
	}
	for _, mv := range sm.MaterializedViews {
		if mv.Name == relName && setTriggerCommentInList(mv.Triggers, trigName, comment) {
			return
		}
	}
}

func findView(sm *metadatapb.SchemaMetadata, name string) *metadatapb.ViewMetadata {
	if sm == nil {
		return nil
	}
	for _, v := range sm.Views {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func setColumnComment(cols []*metadatapb.ColumnMetadata, colName, comment string) {
	for _, col := range cols {
		if col.Name == colName {
			col.Comment = comment
			return
		}
	}
}

func relationDependencyColumns(c *Catalog, rel *Relation) []*metadatapb.DependencyColumn {
	if c == nil || rel == nil {
		return nil
	}
	type relKeyPair struct {
		schema string
		table  string
	}
	type depKey struct {
		schema string
		table  string
		column string
	}
	seen := make(map[depKey]bool)
	hasColDeps := make(map[relKeyPair]bool)
	wholeRels := make(map[relKeyPair]bool)
	var out []*metadatapb.DependencyColumn

	for _, d := range c.deps {
		if d.ObjType != 'r' || d.ObjOID != rel.OID || d.RefType != 'r' {
			continue
		}
		refRel := c.GetRelationByOID(d.RefOID)
		if refRel == nil || refRel.Schema == nil || wtIsSystemSchema(refRel.Schema.Name) {
			continue
		}
		if d.RefSubID == 0 {
			wholeRels[relKeyPair{schema: refRel.Schema.Name, table: refRel.Name}] = true
			continue
		}
		colName := ""
		for _, col := range refRel.Columns {
			if col != nil && col.AttNum == int16(d.RefSubID) {
				colName = col.Name
				break
			}
		}
		if colName == "" {
			continue
		}
		hasColDeps[relKeyPair{schema: refRel.Schema.Name, table: refRel.Name}] = true
		k := depKey{schema: refRel.Schema.Name, table: refRel.Name, column: colName}
		if !seen[k] {
			seen[k] = true
			out = append(out, &metadatapb.DependencyColumn{
				Schema: refRel.Schema.Name,
				Table:  refRel.Name,
				Column: colName,
			})
		}
	}
	for pair := range wholeRels {
		if hasColDeps[pair] {
			continue
		}
		k := depKey{schema: pair.schema, table: pair.table}
		if !seen[k] {
			seen[k] = true
			out = append(out, &metadatapb.DependencyColumn{
				Schema: pair.schema,
				Table:  pair.table,
			})
		}
	}
	slices.SortFunc(out, func(a, b *metadatapb.DependencyColumn) int {
		return cmp.Or(
			cmp.Compare(a.Schema, b.Schema),
			cmp.Compare(a.Table, b.Table),
			cmp.Compare(a.Column, b.Column),
		)
	})
	return out
}

func rulesToProto(c *Catalog, schemaName string, rel *Relation) []*metadatapb.RuleMetadata {
	if c == nil || rel == nil || len(rel.Rules) == 0 {
		return nil
	}
	names := make([]string, 0, len(rel.Rules))
	for name, r := range rel.Rules {
		if r != nil {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	out := make([]*metadatapb.RuleMetadata, 0, len(names))
	for _, name := range names {
		out = append(out, ruleToProto(c, schemaName, rel.Name, rel.Rules[name]))
	}
	return out
}

func ruleToProto(cat *Catalog, schemaName, relName string, r *Rule) *metadatapb.RuleMetadata {
	actionsSQL := cat.deparseRuleActions(r)
	var action string
	switch len(actionsSQL) {
	case 0:
		action = "NOTHING"
	case 1:
		action = actionsSQL[0]
	default:
		action = "(" + strings.Join(actionsSQL, "; ") + ")"
	}
	return &metadatapb.RuleMetadata{
		Name:       r.Name,
		Event:      ruleEventKeyword(r.Event),
		Condition:  cat.deparseRuleWhere(r),
		Action:     action,
		IsInstead:  r.Instead,
		IsEnabled:  normalizeRuleMode(r.Mode) != 'D',
		Definition: buildCreateRuleOpWithReplace(cat, schemaName, relName, r, false).SQL,
	}
}

func applyRuleDiffs(rules *[]*metadatapb.RuleMetadata, cat *Catalog, rel RelationDiffEntry) {
	if len(rel.Rules) == 0 {
		return
	}
	upsert := func(r *Rule) {
		updated := ruleToProto(cat, rel.SchemaName, rel.Name, r)
		for i, existing := range *rules {
			if existing.Name == r.Name {
				(*rules)[i] = updated
				return
			}
		}
		*rules = append(*rules, updated)
	}
	for _, rd := range rel.Rules {
		switch rd.Action {
		case DiffAdd, DiffModify:
			if rd.To != nil {
				upsert(rd.To)
			} else if rd.Action == DiffModify {
				*rules = removeRuleByName(*rules, rd.Name)
			}
		case DiffDrop:
			*rules = removeRuleByName(*rules, rd.Name)
		default:
		}
	}
}

func removeRuleByName(rules []*metadatapb.RuleMetadata, name string) []*metadatapb.RuleMetadata {
	out := make([]*metadatapb.RuleMetadata, 0, len(rules))
	for _, r := range rules {
		if r.Name != name {
			out = append(out, r)
		}
	}
	return out
}

func applyCollationDiffs(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, collations []CollationDiffEntry, renamedSchemas map[string]string, canonicalSchema func(string) string) {
	if len(collations) == 0 {
		return
	}
	addedByOID := make(map[uint32]CollationDiffEntry)
	for _, col := range collations {
		if col.Action == DiffAdd && col.To != nil && col.To.OID != 0 {
			addedByOID[col.To.OID] = col
		}
	}

	handledDropOIDs := make(map[uint32]bool)
	handledAddOIDs := make(map[uint32]bool)
	for _, col := range collations {
		if col.Action != DiffDrop || col.From == nil || col.From.OID == 0 {
			continue
		}
		add, ok := addedByOID[col.From.OID]
		if !ok || add.To == nil {
			continue
		}
		var prev *metadatapb.CollationMetadata
		if srcSM := findSchema(result, canonicalSchema(col.SchemaName)); srcSM != nil {
			prev, srcSM.Collations = popCollationByName(srcSM.Collations, col.Name)
		}
		dstSchema := add.SchemaName
		if renamed, ok := renamedSchemas[col.SchemaName]; ok && col.SchemaName == add.SchemaName {
			dstSchema = renamed
		}
		dstSM := findOrCreateSchema(result, dstSchema)
		upsertCollation(dstSM, collationToProto(catAfter, add.To, prev))
		handledDropOIDs[col.From.OID] = true
		handledAddOIDs[add.To.OID] = true
	}

	for _, col := range collations {
		if col.Action == DiffDrop && col.From != nil && handledDropOIDs[col.From.OID] {
			continue
		}
		if col.Action == DiffAdd && col.To != nil && handledAddOIDs[col.To.OID] {
			continue
		}
		schemaName := canonicalSchema(col.SchemaName)
		if col.Action == DiffDrop {
			if sm := findSchema(result, schemaName); sm != nil {
				_, sm.Collations = popCollationByName(sm.Collations, col.Name)
			}
			continue
		}
		if col.To == nil {
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch col.Action {
		case DiffAdd, DiffModify:
			upsertCollation(sm, collationToProto(catAfter, col.To, findCollationInSchema(sm, col.Name)))
		default:
		}
	}
}

func collationToProto(cat *Catalog, uc *UserCollation, previous *metadatapb.CollationMetadata) *metadatapb.CollationMetadata {
	out := &metadatapb.CollationMetadata{Name: uc.Name}
	if previous != nil {
		out.Comment = previous.Comment
		out.SkipDump = previous.SkipDump
		out.Provider = previous.Provider
		out.Locale = previous.Locale
		out.LcCollate = previous.LcCollate
		out.LcCtype = previous.LcCtype
		out.Rules = previous.Rules
		if previous.Deterministic != nil {
			d := *previous.Deterministic
			out.Deterministic = &d
		}
	}
	applyCollationDefClause(cat, out, uc.DefClause)
	return out
}

func applyCollationDefClause(cat *Catalog, cm *metadatapb.CollationMetadata, defClause string) {
	resolved := resolveEffectiveCollationDefClause(cat, defClause, 0)
	trimmed := strings.TrimSpace(resolved)
	if !strings.HasPrefix(trimmed, "(") || !strings.HasSuffix(trimmed, ")") {
		applyBuiltinCollationDef(cm, trimmed)
		return
	}
	inner := trimmed[1 : len(trimmed)-1]
	// Clear fields derived from DefClause before parsing key=value pairs.
	cm.Provider = ""
	cm.Locale = ""
	cm.LcCollate = ""
	cm.LcCtype = ""
	cm.Rules = ""
	cm.Deterministic = nil

	for _, entry := range splitCollationDefItems(inner) {
		k, v, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(k))
		val := unquoteCollationValue(strings.TrimSpace(v))
		switch key {
		case "PROVIDER":
			switch strings.ToLower(val) {
			case "icu", "i":
				cm.Provider = "i"
			case "libc", "c":
				cm.Provider = "c"
			case "builtin", "b":
				cm.Provider = "b"
			case "default", "d":
				cm.Provider = "d"
			}
		case "LOCALE":
			cm.Locale = val
		case "LC_COLLATE":
			cm.LcCollate = val
		case "LC_CTYPE":
			cm.LcCtype = val
		case "RULES":
			cm.Rules = val
		case "DETERMINISTIC":
			b := strings.EqualFold(val, "true")
			cm.Deterministic = &b
		}
	}
	if cm.Provider == "" && (cm.LcCollate != "" || cm.LcCtype != "") {
		cm.Provider = "c"
	}
}

func resolveEffectiveCollationDefClause(cat *Catalog, defClause string, depth int) string {
	trimmed := strings.TrimSpace(defClause)
	if len(trimmed) < 5 || !strings.EqualFold(trimmed[:4], "FROM") || (trimmed[4] != ' ' && trimmed[4] != '\t' && trimmed[4] != '\n') {
		return trimmed
	}
	target := strings.TrimSpace(trimmed[4:])
	if depth >= 16 || cat == nil {
		return target
	}
	schemaName, collName, ok := splitQualifiedIdent(target)
	if !ok {
		collName = unquoteIdent(target)
	}
	if collName == "" {
		return target
	}
	if schemaName != "" {
		if s := cat.schemaByName[schemaName]; s != nil {
			if uc := s.Collations[collName]; uc != nil && uc.DefClause != "" {
				return resolveEffectiveCollationDefClause(cat, uc.DefClause, depth+1)
			}
		}
		return target
	}
	for _, nsOID := range cat.searchPathWithCatalog() {
		if s := cat.schemas[nsOID]; s != nil {
			if uc := s.Collations[collName]; uc != nil && uc.DefClause != "" {
				return resolveEffectiveCollationDefClause(cat, uc.DefClause, depth+1)
			}
		}
	}
	return target
}

func applyBuiltinCollationDef(cm *metadatapb.CollationMetadata, fromTarget string) {
	if cm == nil || fromTarget == "" {
		return
	}
	schemaName, collName, ok := splitQualifiedIdent(fromTarget)
	if !ok {
		collName = unquoteIdent(fromTarget)
	}
	if collName == "" || (schemaName != "" && schemaName != "pg_catalog") {
		return
	}
	det := true
	switch collName {
	case "C", "POSIX", "C.UTF-8", "C.utf8":
		cm.Provider = "c"
		cm.Locale = ""
		cm.LcCollate = collName
		cm.LcCtype = collName
		cm.Rules = ""
		cm.Deterministic = &det
	case "ucs_basic":
		cm.Provider = "c"
		cm.Locale = ""
		cm.LcCollate = "C"
		cm.LcCtype = "C"
		cm.Rules = ""
		cm.Deterministic = &det
	case "pg_c_utf8":
		cm.Provider = "b"
		cm.Locale = "C.UTF-8"
		cm.LcCollate = ""
		cm.LcCtype = ""
		cm.Rules = ""
		cm.Deterministic = &det
	case "unicode":
		cm.Provider = "b"
		cm.Locale = "PG_UNICODE_FAST"
		cm.LcCollate = ""
		cm.LcCtype = ""
		cm.Rules = ""
		cm.Deterministic = &det
	default:
		if strings.HasSuffix(strings.ToLower(collName), "-x-icu") {
			locale := collName[:len(collName)-len("-x-icu")]
			cm.Provider = "i"
			cm.Locale = locale
			cm.LcCollate = ""
			cm.LcCtype = ""
			cm.Rules = ""
			cm.Deterministic = &det
		} else if isBuiltinCollationName(collName) {
			cm.Provider = "c"
			cm.Locale = ""
			cm.LcCollate = collName
			cm.LcCtype = collName
			cm.Rules = ""
			cm.Deterministic = &det
		}
	}
}

func splitCollationDefItems(s string) []string {
	var out []string
	var cur strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inQuote {
			cur.WriteByte(ch)
			if ch == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					cur.WriteByte(s[i+1])
					i++
					continue
				}
				inQuote = false
			}
			continue
		}
		if ch == '\'' {
			inQuote = true
			cur.WriteByte(ch)
			continue
		}
		if ch == ',' {
			out = append(out, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteByte(ch)
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

func unquoteCollationValue(v string) string {
	if len(v) >= 2 && v[0] == '\'' && v[len(v)-1] == '\'' {
		return strings.ReplaceAll(v[1:len(v)-1], "''", "'")
	}
	return v
}

func findCollationInSchema(sm *metadatapb.SchemaMetadata, name string) *metadatapb.CollationMetadata {
	if sm == nil {
		return nil
	}
	for _, c := range sm.Collations {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func popCollationByName(collations []*metadatapb.CollationMetadata, name string) (*metadatapb.CollationMetadata, []*metadatapb.CollationMetadata) {
	var popped *metadatapb.CollationMetadata
	out := make([]*metadatapb.CollationMetadata, 0, len(collations))
	for _, c := range collations {
		if c.Name == name && popped == nil {
			popped = c
			continue
		}
		out = append(out, c)
	}
	return popped, out
}

func upsertCollation(sm *metadatapb.SchemaMetadata, updated *metadatapb.CollationMetadata) {
	for i, existing := range sm.Collations {
		if existing.Name == updated.Name {
			sm.Collations[i] = updated
			return
		}
	}
	sm.Collations = append(sm.Collations, updated)
}

func collectExtensionMemberOIDs(catBefore, catAfter *Catalog) []uint32 {
	if catBefore == nil && catAfter == nil {
		return nil
	}
	seen := make(map[uint32]bool)
	var oids []uint32
	if catBefore != nil {
		for oid := range catBefore.extensionMembers {
			if !seen[oid] {
				seen[oid] = true
				oids = append(oids, oid)
			}
		}
	}
	if catAfter != nil {
		for oid := range catAfter.extensionMembers {
			if !seen[oid] {
				seen[oid] = true
				oids = append(oids, oid)
			}
		}
	}
	slices.Sort(oids)
	return oids
}

func collationByOID(c *Catalog, oid uint32) (string, *UserCollation) {
	if c == nil {
		return "", nil
	}
	for _, s := range c.schemas {
		if s.OID == PGCatalogNamespace || s.OID == PGToastNamespace {
			continue
		}
		for _, uc := range s.Collations {
			if uc != nil && uc.OID == oid {
				return s.Name, uc
			}
		}
	}
	return "", nil
}

func hasExtensionMemberChange(catBefore, catAfter *Catalog) bool {
	if catBefore == nil || catAfter == nil {
		return false
	}
	for _, oid := range collectExtensionMemberOIDs(catBefore, catAfter) {
		if catBefore.extensionMembers[oid] != catAfter.extensionMembers[oid] {
			return true
		}
		if !catAfter.isExtensionMember(oid) {
			continue
		}
		if fromRel, toRel := catBefore.relationByOID[oid], catAfter.relationByOID[oid]; fromRel != nil || toRel != nil {
			if fromRel == nil || toRel == nil {
				return true
			}
			if toRel.RelKind != 'i' && toRel.RelKind != 'c' {
				if extensionRelationChanged(catBefore, fromRel, catAfter, toRel) {
					return true
				}
			}
		}
		if fromSeq, toSeq := catBefore.sequenceByOID[oid], catAfter.sequenceByOID[oid]; fromSeq != nil || toSeq != nil {
			if fromSeq == nil || toSeq == nil {
				return true
			}
			if fromSeq.Schema == nil || toSeq.Schema == nil || fromSeq.Schema.Name != toSeq.Schema.Name ||
				fromSeq.Name != toSeq.Name || fromSeq.OwnerName != toSeq.OwnerName || sequencesChanged(catBefore, catAfter, fromSeq, toSeq) {
				return true
			}
		}
		if fromUP, toUP := catBefore.userProcs[oid], catAfter.userProcs[oid]; fromUP != nil || toUP != nil {
			if fromUP == nil || toUP == nil {
				return true
			}
			if funcIdentity(catBefore, fromUP) != funcIdentity(catAfter, toUP) || funcChanged(catBefore, catAfter, fromUP, toUP) {
				return true
			}
		}
		if fromBT, toBT := catBefore.typeByOID[oid], catAfter.typeByOID[oid]; fromBT != nil || toBT != nil {
			if extensionTypeChanged(catBefore, fromBT, catAfter, toBT) {
				return true
			}
		}
		fromCollSchema, fromColl := collationByOID(catBefore, oid)
		toCollSchema, toColl := collationByOID(catAfter, oid)
		if fromColl != nil || toColl != nil {
			if fromColl == nil || toColl == nil || fromCollSchema != toCollSchema || fromColl.Name != toColl.Name ||
				fromColl.DefClause != toColl.DefClause || fromColl.OwnerName != toColl.OwnerName {
				return true
			}
		}
	}
	return false
}

func extensionRelationChanged(catBefore *Catalog, fromRel *Relation, catAfter *Catalog, toRel *Relation) bool {
	if fromRel.Schema == nil || toRel.Schema == nil {
		return true
	}
	if fromRel.Schema.Name != toRel.Schema.Name || fromRel.Name != toRel.Name || fromRel.OwnerName != toRel.OwnerName {
		return true
	}
	key := relKey{schema: toRel.Schema.Name, name: toRel.Name}
	switch toRel.RelKind {
	case 'v', 'm':
		_, changed := compareView(catBefore, catAfter, key, fromRel, toRel)
		return changed
	default:
		_, changed := compareRelation(catBefore, catAfter, key, fromRel, toRel)
		return changed
	}
}

func extensionTypeChanged(catBefore *Catalog, fromBT *BuiltinType, catAfter *Catalog, toBT *BuiltinType) bool {
	if fromBT == nil || toBT == nil {
		return true
	}
	fromNS, toNS := catBefore.schemas[fromBT.Namespace], catAfter.schemas[toBT.Namespace]
	if fromNS == nil || toNS == nil || fromNS.Name != toNS.Name || fromBT.TypeName != toBT.TypeName || fromBT.Type != toBT.Type {
		return true
	}
	switch toBT.Type {
	case 'e':
		return !slices.Equal(enumLabels(catBefore, fromBT.OID), enumLabels(catAfter, toBT.OID))
	case 'c':
		fromRel, toRel := catBefore.relationByOID[fromBT.RelID], catAfter.relationByOID[toBT.RelID]
		if fromRel == nil || toRel == nil {
			return fromRel != toRel
		}
		return len(diffColumns(catBefore, catAfter, fromRel, toRel)) > 0
	case 'r':
		fromRT, toRT := catBefore.rangeTypes[fromBT.OID], catAfter.rangeTypes[toBT.OID]
		if fromRT == nil || toRT == nil {
			return fromRT != toRT
		}
		return rangeChanged(catBefore, fromRT, catAfter, toRT)
	}
	return false
}

func applyExtensionMemberChanges(result, original *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, canonicalSchema func(string) string) {
	if catBefore == nil || catAfter == nil {
		return
	}
	var relEntries []RelationDiffEntry
	var ftEntries []RelationDiffEntry
	var seqEntries []SequenceDiffEntry
	var fnEntries []FunctionDiffEntry
	var rtEntries []RangeDiffEntry
	var collEntries []CollationDiffEntry

	for _, oid := range collectExtensionMemberOIDs(catBefore, catAfter) {
		inBefore := catBefore.isExtensionMember(oid)
		inAfter := catAfter.isExtensionMember(oid)
		if !inBefore && !inAfter {
			continue
		}
		// Skip whole schemas handled by schema diffs.
		if s := catAfter.schemas[oid]; s != nil && oid != PGCatalogNamespace && oid != PGToastNamespace {
			if sm := findSchema(result, s.Name); sm != nil {
				sm.SkipDump = inAfter
			}
			continue
		}
		// Relations (tables, partitions, views, matviews, foreign tables).
		if fromRel, toRel := catBefore.relationByOID[oid], catAfter.relationByOID[oid]; fromRel != nil || toRel != nil {
			switch {
			case toRel == nil && inBefore:
				if fromRel.Schema != nil {
					dropRelation(result, canonicalSchema(fromRel.Schema.Name), RelationDiffEntry{
						Action:     DiffDrop,
						SchemaName: fromRel.Schema.Name,
						Name:       fromRel.Name,
						From:       fromRel,
					})
				}
			case toRel != nil && toRel.RelKind != 'i' && toRel.RelKind != 'c':
				toSchema := ""
				if toRel.Schema != nil {
					toSchema = toRel.Schema.Name
				}
				fromSchema := ""
				if fromRel != nil && fromRel.Schema != nil {
					fromSchema = fromRel.Schema.Name
				}
				if inAfter {
					if !inBefore && fromRel != nil && original != nil {
						restoreOriginalRelation(result, original, canonicalSchema(fromSchema), fromRel.Name, fromRel.RelKind)
					}
					if fromRel == nil {
						entry := RelationDiffEntry{Action: DiffAdd, SchemaName: toSchema, Name: toRel.Name, To: toRel}
						if toRel.RelKind == 'f' {
							ftEntries = append(ftEntries, entry)
						} else {
							relEntries = append(relEntries, entry)
						}
					} else if extensionRelationChanged(catBefore, fromRel, catAfter, toRel) {
						if fromSchema != toSchema || fromRel.Name != toRel.Name {
							dropE := RelationDiffEntry{Action: DiffDrop, SchemaName: fromSchema, Name: fromRel.Name, From: fromRel}
							addE := RelationDiffEntry{Action: DiffAdd, SchemaName: toSchema, Name: toRel.Name, To: toRel}
							if toRel.RelKind == 'f' {
								ftEntries = append(ftEntries, dropE, addE)
							} else {
								relEntries = append(relEntries, dropE, addE)
							}
						} else {
							var entry RelationDiffEntry
							key := relKey{schema: toSchema, name: toRel.Name}
							if toRel.RelKind == 'v' || toRel.RelKind == 'm' {
								entry, _ = compareView(catBefore, catAfter, key, fromRel, toRel)
							} else {
								entry, _ = compareRelation(catBefore, catAfter, key, fromRel, toRel)
								if entry.Action == 0 {
									entry = RelationDiffEntry{Action: DiffModify, SchemaName: toSchema, Name: toRel.Name, From: fromRel, To: toRel}
								}
							}
							if toRel.RelKind == 'f' {
								ftEntries = append(ftEntries, entry)
							} else {
								relEntries = append(relEntries, entry)
							}
						}
					}
				}
			}
		}

		if fromSeq, toSeq := catBefore.sequenceByOID[oid], catAfter.sequenceByOID[oid]; fromSeq != nil || toSeq != nil {
			if toSeq != nil && inAfter {
				fromSchema, toSchema := "", ""
				if fromSeq != nil && fromSeq.Schema != nil {
					fromSchema = fromSeq.Schema.Name
				}
				if toSeq.Schema != nil {
					toSchema = toSeq.Schema.Name
				}
				if !inBefore && fromSeq != nil && original != nil {
					if srcSM := findSchema(original, canonicalSchema(fromSchema)); srcSM != nil {
						if prev := findSequenceInSchema(srcSM, fromSeq.Name); prev != nil {
							upsertSequence(findOrCreateSchema(result, canonicalSchema(fromSchema)), proto.CloneOf(prev))
						}
					}
				}
				if fromSeq == nil {
					seqEntries = append(seqEntries, SequenceDiffEntry{Action: DiffAdd, SchemaName: toSchema, Name: toSeq.Name, To: toSeq})
				} else if fromSchema != toSchema || fromSeq.Name != toSeq.Name {
					seqEntries = append(seqEntries,
						SequenceDiffEntry{Action: DiffDrop, SchemaName: fromSchema, Name: fromSeq.Name, From: fromSeq},
						SequenceDiffEntry{Action: DiffAdd, SchemaName: toSchema, Name: toSeq.Name, To: toSeq},
					)
				} else if fromSeq.OwnerName != toSeq.OwnerName || sequencesChanged(catBefore, catAfter, fromSeq, toSeq) {
					seqEntries = append(seqEntries, SequenceDiffEntry{Action: DiffModify, SchemaName: toSchema, Name: toSeq.Name, From: fromSeq, To: toSeq})
				}
			}
		}

		if fromUP, toUP := catBefore.userProcs[oid], catAfter.userProcs[oid]; fromUP != nil || toUP != nil {
			if toUP != nil && inAfter {
				toSchema := ""
				if toUP.Schema != nil {
					toSchema = toUP.Schema.Name
				}
				toID := funcIdentity(catAfter, toUP)
				fromID := ""
				fromSchema := ""
				if fromUP != nil {
					if fromUP.Schema != nil {
						fromSchema = fromUP.Schema.Name
					}
					fromID = funcIdentity(catBefore, fromUP)
				}
				if !inBefore && fromUP != nil && original != nil {
					if srcSM := findSchema(original, canonicalSchema(fromSchema)); srcSM != nil {
						if fromUP.Kind == 'p' {
							if p, _ := popProcedureByIdentity(srcSM.Procedures, fromID); p != nil {
								dstSM := findOrCreateSchema(result, canonicalSchema(fromSchema))
								if _, rem := popProcedureByIdentity(dstSM.Procedures, fromID); len(rem) == len(dstSM.Procedures) {
									dstSM.Procedures = append(dstSM.Procedures, proto.CloneOf(p))
								}
							}
						} else {
							if fn, _ := popFunctionByIdentity(srcSM.Functions, fromID); fn != nil {
								dstSM := findOrCreateSchema(result, canonicalSchema(fromSchema))
								if _, rem := popFunctionByIdentity(dstSM.Functions, fromID); len(rem) == len(dstSM.Functions) {
									dstSM.Functions = append(dstSM.Functions, proto.CloneOf(fn))
								}
							}
						}
					}
				}
				if fromUP == nil {
					fnEntries = append(fnEntries, FunctionDiffEntry{Action: DiffAdd, SchemaName: toSchema, Identity: toID, To: toUP})
				} else if fromID != toID {
					fnEntries = append(fnEntries,
						FunctionDiffEntry{Action: DiffDrop, SchemaName: fromSchema, Identity: fromID, From: fromUP},
						FunctionDiffEntry{Action: DiffAdd, SchemaName: toSchema, Identity: toID, To: toUP},
					)
				} else if funcChanged(catBefore, catAfter, fromUP, toUP) {
					fnEntries = append(fnEntries, FunctionDiffEntry{Action: DiffModify, SchemaName: toSchema, Identity: toID, From: fromUP, To: toUP})
				}
			}
		}

		if fromBT, toBT := catBefore.typeByOID[oid], catAfter.typeByOID[oid]; fromBT != nil || toBT != nil {
			if toBT != nil && inAfter {
				toNS := catAfter.schemas[toBT.Namespace]
				if toNS != nil {
					var fromNS *Schema
					if fromBT != nil {
						fromNS = catBefore.schemas[fromBT.Namespace]
					}
					if !inBefore && fromBT != nil && fromNS != nil && original != nil {
						restoreOriginalType(result, original, canonicalSchema(fromNS.Name), fromBT.TypeName, fromBT.Type)
					}
					switch toBT.Type {
					case 'e':
						if fromBT == nil || extensionTypeChanged(catBefore, fromBT, catAfter, toBT) {
							var prev *metadatapb.EnumTypeMetadata
							if fromBT != nil && fromNS != nil {
								if srcSM := findSchema(result, canonicalSchema(fromNS.Name)); srcSM != nil {
									for _, e := range srcSM.EnumTypes {
										if e.Name == fromBT.TypeName {
											prev = proto.CloneOf(e)
											break
										}
									}
									if fromNS.Name != toNS.Name || fromBT.TypeName != toBT.TypeName {
										srcSM.EnumTypes = removeEnumByName(srcSM.EnumTypes, fromBT.TypeName)
									}
								}
							}
							dstSM := findOrCreateSchema(result, toNS.Name)
							upsertEnum(dstSM, toBT.TypeName, enumLabels(catAfter, toBT.OID))
							if prev != nil {
								for _, e := range dstSM.EnumTypes {
									if e.Name == toBT.TypeName {
										e.Comment = prev.Comment
										e.SkipDump = prev.SkipDump
									}
								}
							}
						}
					case 'c':
						if fromBT == nil || extensionTypeChanged(catBefore, fromBT, catAfter, toBT) {
							var prev *metadatapb.CompositeTypeMetadata
							if fromBT != nil && fromNS != nil {
								if srcSM := findSchema(result, canonicalSchema(fromNS.Name)); srcSM != nil {
									for _, ct := range srcSM.CompositeTypes {
										if ct.Name == fromBT.TypeName {
											prev = proto.CloneOf(ct)
											break
										}
									}
									if fromNS.Name != toNS.Name || fromBT.TypeName != toBT.TypeName {
										srcSM.CompositeTypes = removeCompositeTypeByName(srcSM.CompositeTypes, fromBT.TypeName)
									}
								}
							}
							var fromCompRel *Relation
							if fromBT != nil {
								fromCompRel = catBefore.relationByOID[fromBT.RelID]
							}
							toCompRel := catAfter.relationByOID[toBT.RelID]
							if toCompRel != nil {
								dstSM := findOrCreateSchema(result, toNS.Name)
								updated := compositeTypeToProto(catAfter, toBT.TypeName, fromCompRel, toCompRel, prev)
								replaced := false
								for i, ct := range dstSM.CompositeTypes {
									if ct.Name == toBT.TypeName {
										dstSM.CompositeTypes[i] = updated
										replaced = true
										break
									}
								}
								if !replaced {
									dstSM.CompositeTypes = append(dstSM.CompositeTypes, updated)
								}
							}
						}
					case 'r':
						fromRT := catBefore.rangeTypes[oid]
						toRT := catAfter.rangeTypes[oid]
						if toRT != nil && (fromRT == nil || fromNS == nil || fromNS.Name != toNS.Name || fromBT.TypeName != toBT.TypeName || rangeChanged(catBefore, fromRT, catAfter, toRT)) {
							if fromRT == nil || fromNS == nil {
								rtEntries = append(rtEntries, RangeDiffEntry{Action: DiffAdd, SchemaName: toNS.Name, Name: toBT.TypeName, To: toRT})
							} else if fromNS.Name != toNS.Name || fromBT.TypeName != toBT.TypeName {
								rtEntries = append(rtEntries,
									RangeDiffEntry{Action: DiffDrop, SchemaName: fromNS.Name, Name: fromBT.TypeName, From: fromRT},
									RangeDiffEntry{Action: DiffAdd, SchemaName: toNS.Name, Name: toBT.TypeName, To: toRT},
								)
							} else {
								rtEntries = append(rtEntries, RangeDiffEntry{Action: DiffModify, SchemaName: toNS.Name, Name: toBT.TypeName, From: fromRT, To: toRT})
							}
						}
					}
				}
			}
		}

		if fromCollSchema, fromColl := collationByOID(catBefore, oid); fromColl != nil {
			toCollSchema, toColl := collationByOID(catAfter, oid)
			if toColl != nil && inAfter {
				if !inBefore && original != nil {
					if srcSM := findSchema(original, canonicalSchema(fromCollSchema)); srcSM != nil {
						if prev := findCollationInSchema(srcSM, fromColl.Name); prev != nil {
							upsertCollation(findOrCreateSchema(result, canonicalSchema(fromCollSchema)), proto.CloneOf(prev))
						}
					}
				}
				if fromCollSchema != toCollSchema || fromColl.Name != toColl.Name {
					collEntries = append(collEntries,
						CollationDiffEntry{Action: DiffDrop, SchemaName: fromCollSchema, Name: fromColl.Name, From: fromColl},
						CollationDiffEntry{Action: DiffAdd, SchemaName: toCollSchema, Name: toColl.Name, To: toColl},
					)
				} else if fromColl.DefClause != toColl.DefClause || fromColl.OwnerName != toColl.OwnerName {
					collEntries = append(collEntries, CollationDiffEntry{Action: DiffModify, SchemaName: toCollSchema, Name: toColl.Name, From: fromColl, To: toColl})
				}
			}
		} else if toCollSchema, toColl := collationByOID(catAfter, oid); toColl != nil && inAfter {
			collEntries = append(collEntries, CollationDiffEntry{Action: DiffAdd, SchemaName: toCollSchema, Name: toColl.Name, To: toColl})
		}
	}

	if len(relEntries) > 0 {
		renamed := renamePartitions(result, catBefore, catAfter, relEntries, canonicalSchema)
		for _, rel := range relEntries {
			if (rel.From != nil && renamed[rel.From.OID]) || (rel.To != nil && renamed[rel.To.OID]) {
				continue
			}
			sm := findOrCreateSchema(result, canonicalSchema(rel.SchemaName))
			switch rel.Action {
			case DiffAdd:
				addRelation(result, sm, catAfter, rel)
			case DiffModify:
				modifyRelation(result, sm, catBefore, catAfter, rel)
			}
		}
	}
	applySequenceDiffs(result, catAfter, seqEntries, nil, canonicalSchema)
	applyFunctionDiffs(result, catAfter, fnEntries, canonicalSchema)
	applyRangeDiffs(result, catAfter, rtEntries, canonicalSchema)
	applyCollationDiffs(result, catAfter, collEntries, nil, canonicalSchema)
	applyForeignTableDiffs(result, catBefore, catAfter, ftEntries, canonicalSchema)

	// Synchronize SkipDump flags with final extension membership in catAfter.
	for _, oid := range collectExtensionMemberOIDs(catBefore, catAfter) {
		inAfter := catAfter.isExtensionMember(oid)
		setExtensionObjectSkipDump(result, catAfter, oid, inAfter)
	}
}

func restoreOriginalRelation(result, original *metadatapb.DatabaseSchemaMetadata, schemaName, name string, relKind byte) {
	origSM := findSchema(original, schemaName)
	if origSM == nil {
		return
	}
	dstSM := findOrCreateSchema(result, schemaName)
	switch relKind {
	case 'v':
		if v := findView(origSM, name); v != nil && findView(dstSM, name) == nil {
			dstSM.Views = append(dstSM.Views, proto.CloneOf(v))
		}
	case 'm':
		for _, mv := range origSM.MaterializedViews {
			if mv.Name == name {
				dstSM.MaterializedViews = append(dstSM.MaterializedViews, proto.CloneOf(mv))
				return
			}
		}
	case 'f':
		if ext := findExternalTable(origSM, name); ext != nil && findExternalTable(dstSM, name) == nil {
			dstSM.ExternalTables = append(dstSM.ExternalTables, proto.CloneOf(ext))
		}
	default:
		if tbl := findTable(origSM, name); tbl != nil && findTable(dstSM, name) == nil {
			dstSM.Tables = append(dstSM.Tables, proto.CloneOf(tbl))
		}
	}
}

func restoreOriginalType(result, original *metadatapb.DatabaseSchemaMetadata, schemaName, name string, typ byte) {
	origSM := findSchema(original, schemaName)
	if origSM == nil {
		return
	}
	dstSM := findOrCreateSchema(result, schemaName)
	switch typ {
	case 'e':
		for _, e := range origSM.EnumTypes {
			if e.Name == name {
				dstSM.EnumTypes = append(dstSM.EnumTypes, proto.CloneOf(e))
				return
			}
		}
	case 'c':
		for _, ct := range origSM.CompositeTypes {
			if ct.Name == name {
				dstSM.CompositeTypes = append(dstSM.CompositeTypes, proto.CloneOf(ct))
				return
			}
		}
	case 'r':
		if rt := findRangeType(origSM, name); rt != nil && findRangeType(dstSM, name) == nil {
			dstSM.RangeTypes = append(dstSM.RangeTypes, proto.CloneOf(rt))
		}
	}
}

func setExtensionObjectSkipDump(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, oid uint32, skipDump bool) {
	if result == nil || catAfter == nil {
		return
	}
	if rel := catAfter.relationByOID[oid]; rel != nil && rel.Schema != nil {
		if sm := findSchema(result, rel.Schema.Name); sm != nil {
			switch rel.RelKind {
			case 'v':
				if v := findView(sm, rel.Name); v != nil {
					v.SkipDump = skipDump
				}
			case 'm':
				for _, mv := range sm.MaterializedViews {
					if mv.Name == rel.Name {
						mv.SkipDump = skipDump
						break
					}
				}
			case 'f':
				if ext := findExternalTable(sm, rel.Name); ext != nil {
					ext.SkipDump = skipDump
				}
			default:
				if tbl := findTable(sm, rel.Name); tbl != nil {
					tbl.SkipDump = skipDump
				}
			}
		}
	}
	if seq := catAfter.sequenceByOID[oid]; seq != nil && seq.Schema != nil {
		if sm := findSchema(result, seq.Schema.Name); sm != nil {
			if s := findSequenceInSchema(sm, seq.Name); s != nil {
				s.SkipDump = skipDump
			}
		}
	}
	if up := catAfter.userProcs[oid]; up != nil && up.Schema != nil {
		if sm := findSchema(result, up.Schema.Name); sm != nil {
			id := funcIdentity(catAfter, up)
			shortID := strings.TrimPrefix(id, up.Schema.Name+".")
			if up.Kind == 'p' {
				for _, p := range sm.Procedures {
					if p.Signature == id || p.Signature == shortID {
						p.SkipDump = skipDump
					}
				}
			} else {
				for _, fn := range sm.Functions {
					if fn.Signature == id || fn.Signature == shortID {
						fn.SkipDump = skipDump
					}
				}
			}
		}
	}
	if bt := catAfter.typeByOID[oid]; bt != nil {
		if ns := catAfter.schemas[bt.Namespace]; ns != nil {
			if sm := findSchema(result, ns.Name); sm != nil {
				switch bt.Type {
				case 'e':
					for _, e := range sm.EnumTypes {
						if e.Name == bt.TypeName {
							e.SkipDump = skipDump
						}
					}
				case 'c':
					for _, ct := range sm.CompositeTypes {
						if ct.Name == bt.TypeName {
							ct.SkipDump = skipDump
						}
					}
				case 'r':
					if rt := findRangeType(sm, bt.TypeName); rt != nil {
						rt.SkipDump = skipDump
					}
				}
			}
		}
	}
	if collSchema, uc := collationByOID(catAfter, oid); uc != nil {
		if sm := findSchema(result, collSchema); sm != nil {
			if c := findCollationInSchema(sm, uc.Name); c != nil {
				c.SkipDump = skipDump
			}
		}
	}
}

func enumLabels(c *Catalog, typeOID uint32) []string {
	if c == nil {
		return nil
	}
	et := c.enumTypes[typeOID]
	if et == nil {
		return nil
	}
	labels := make([]string, len(et.Values))
	for i, ev := range et.Values {
		labels[i] = ev.Label
	}
	return labels
}

func rangeChanged(catBefore *Catalog, fromRT *RangeType, catAfter *Catalog, toRT *RangeType) bool {
	if fromRT == nil || toRT == nil {
		return fromRT != toRT
	}
	return catBefore.FormatType(fromRT.SubTypeOID, -1) != catAfter.FormatType(toRT.SubTypeOID, -1)
}

func findSequenceInSchema(sm *metadatapb.SchemaMetadata, name string) *metadatapb.SequenceMetadata {
	if sm == nil {
		return nil
	}
	for _, s := range sm.Sequences {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func isRecreatedConstraint(catBefore, catAfter *Catalog, schemaName, tableName, conName string) bool {
	if catBefore == nil || catAfter == nil {
		return false
	}
	fromRel := catBefore.GetRelation(schemaName, tableName)
	toRel := catAfter.GetRelation(schemaName, tableName)
	if fromRel == nil || toRel == nil {
		return false
	}
	var fromOID, toOID uint32
	for _, con := range catBefore.ConstraintsOf(fromRel.OID) {
		if con != nil && con.Name == conName {
			fromOID = con.OID
			break
		}
	}
	for _, con := range catAfter.ConstraintsOf(toRel.OID) {
		if con != nil && con.Name == conName {
			toOID = con.OID
			break
		}
	}
	return fromOID != 0 && toOID != 0 && fromOID != toOID
}

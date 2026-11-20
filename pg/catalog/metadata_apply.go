package catalog

import (
	"fmt"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"

	metadatapb "github.com/bytebase/omni/metadata"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// ApplyDiffToMetadata applies a SchemaDiff (from user DDL execution) to a
// clone of the original metadata proto, producing the post-DDL metadata.
// The original proto is NOT modified.
func ApplyDiffToMetadata(original *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, diff *SchemaDiff) *metadatapb.DatabaseSchemaMetadata {
	if (diff == nil || diff.IsEmpty()) && !hasRelationOwnerChange(catBefore, catAfter) {
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

	renamed := renamePartitions(result, catAfter, diff.Relations, canonicalSchema)
	for _, rel := range diff.Relations {
		if (rel.From != nil && renamed[rel.From.OID]) || (rel.To != nil && renamed[rel.To.OID]) {
			continue
		}
		schemaName := canonicalSchema(rel.SchemaName)
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

	for _, seq := range diff.Sequences {
		schemaName := canonicalSchema(seq.SchemaName)
		if seq.Action == DiffDrop {
			if renamedSchemas[seq.SchemaName] != "" && hasMatchingSequenceAdd(diff.Sequences, renamedSchemas[seq.SchemaName], seq.Name) {
				continue
			}
			if sm := findSchema(result, schemaName); sm != nil {
				sm.Sequences = removeSequenceByName(sm.Sequences, seq.Name)
			}
			continue
		}
		sm := findOrCreateSchema(result, schemaName)
		switch seq.Action {
		case DiffAdd:
			if seq.To != nil {
				upsertSequence(sm, sequenceToProto(catAfter, seq.To))
			}
		case DiffModify:
			if seq.To != nil {
				upsertSequence(sm, sequenceToProto(catAfter, seq.To))
			}
		default:
		}
	}

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

	applyFunctionDiffs(result, catAfter, diff.Functions, canonicalSchema)

	applyEventTriggerDiffs(result, catAfter, diff.EventTriggers)
	applyCommentDiffs(result, catBefore, catAfter, diff.Comments)
	applyExtensionDiffs(result, catAfter, diff.Extensions)
	applyForeignTableDiffs(result, catBefore, catAfter, diff.ForeignTables, canonicalSchema)
	applyRelationOwnerChanges(result, catBefore, catAfter)

	return result
}

func hasRelationOwnerChange(catBefore, catAfter *Catalog) bool {
	if catBefore == nil || catAfter == nil {
		return false
	}
	for _, s := range catAfter.UserSchemas() {
		for _, rel := range s.Relations {
			if rel.RelKind != 'r' && rel.RelKind != 'p' {
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
		for _, rel := range s.Relations {
			if rel.RelKind != 'r' && rel.RelKind != 'p' {
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
		upsertRoutineFromDiff(srcSM, dstSM, catAfter, dropEntry.Identity, addEntry.Identity, addEntry.To)
	}
	for _, f := range functions {
		if (f.From != nil && f.From.OID != 0 && handledOIDs[f.From.OID]) ||
			(f.To != nil && f.To.OID != 0 && handledOIDs[f.To.OID]) {
			continue
		}
		schemaName := canonicalSchema(f.SchemaName)
		switch f.Action {
		case DiffDrop:
			if sm := findSchema(result, schemaName); sm != nil {
				if f.From != nil && f.From.Kind == 'p' {
					_, sm.Procedures = popProcedureByIdentity(sm.Procedures, f.Identity)
				} else {
					_, sm.Functions = popFunctionByIdentity(sm.Functions, f.Identity)
				}
			}
		case DiffAdd:
			if f.To != nil {
				sm := findOrCreateSchema(result, schemaName)
				upsertRoutineFromDiff(nil, sm, catAfter, "", f.Identity, f.To)
			}
		case DiffModify:
			if f.To != nil {
				sm := findOrCreateSchema(result, schemaName)
				upsertRoutineFromDiff(sm, sm, catAfter, f.Identity, f.Identity, f.To)
			}
		default:
		}
	}
}

func upsertRoutineFromDiff(srcSM, dstSM *metadatapb.SchemaMetadata, catAfter *Catalog, oldIdentity, newIdentity string, to *UserProc) {
	if to == nil || dstSM == nil {
		return
	}
	if to.Kind == 'p' {
		var prev *metadatapb.ProcedureMetadata
		if srcSM != nil && oldIdentity != "" {
			prev, srcSM.Procedures = popProcedureByIdentity(srcSM.Procedures, oldIdentity)
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
		prev, srcSM.Functions = popFunctionByIdentity(srcSM.Functions, oldIdentity)
	}
	fn := functionToProto(catAfter, to, newIdentity)
	if prev != nil {
		fn.Comment = prev.Comment
		fn.SkipDump = prev.SkipDump
		fn.DependencyTables = prev.DependencyTables
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

func extSchemaName(cat *Catalog, ext *Extension) string {
	if cat == nil || ext == nil || ext.SchemaOID == 0 {
		return ""
	}
	if s := cat.schemas[ext.SchemaOID]; s != nil {
		return s.Name
	}
	return ""
}

func applyExtensionDiffs(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, extensions []ExtensionDiffEntry) {
	for _, e := range extensions {
		switch e.Action {
		case DiffDrop:
			result.Extensions = slices.DeleteFunc(result.Extensions, func(ext *metadatapb.ExtensionMetadata) bool {
				return ext.Name == e.Name
			})
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

func applyEventTriggerDiffs(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, entries []EventTriggerDiffEntry) {
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
		upsertEventTrigger(result, eventTriggerToProto(catAfter, add.To, prev))
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
			upsertEventTrigger(result, eventTriggerToProto(catAfter, e.To, findEventTrigger(result, e.Name)))
		default:
		}
	}
}

func eventTriggerToProto(cat *Catalog, evt *EventTrigger, previous *metadatapb.EventTriggerMetadata) *metadatapb.EventTriggerMetadata {
	out := &metadatapb.EventTriggerMetadata{
		Name:    evt.Name,
		Event:   evt.EventName,
		Tags:    slices.Clone(evt.Tags),
		Enabled: evt.Enabled != 'D',
	}
	if cat != nil {
		if up := cat.userProcs[evt.FuncOID]; up != nil {
			if up.Schema != nil {
				out.FunctionSchema = up.Schema.Name
			}
			out.FunctionName = up.Name
		} else if p := cat.procByOID[evt.FuncOID]; p != nil {
			out.FunctionName = p.Name
		}
	}
	if previous != nil {
		out.Comment = previous.Comment
		out.SkipDump = previous.SkipDump
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

func upsertSequence(sm *metadatapb.SchemaMetadata, updated *metadatapb.SequenceMetadata) {
	for i, existing := range sm.Sequences {
		if existing.Name == updated.Name {
			updated.Comment = existing.Comment
			updated.SkipDump = existing.SkipDump
			updated.LastValue = existing.LastValue
			sm.Sequences[i] = updated
			return
		}
	}
	sm.Sequences = append(sm.Sequences, updated)
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
		applyColumnDiffs(&tbl.Columns, cat, rel)
		applyIndexDiffs(tbl, cat, rel)
		applyConstraintDiffs(tbl, catBefore, cat, rel)
		applyTriggerDiffs(&tbl.Triggers, cat, rel)
		if (rel.From != nil && rel.From.OwnerName != rel.To.OwnerName) || (rel.To.OwnerName != "" && tbl.Owner != rel.To.OwnerName) {
			tbl.Owner = rel.To.OwnerName
		}
	case 'v':
		for i, existing := range sm.Views {
			if existing.Name == rel.Name {
				sm.Views[i] = mergeViewProto(cat, existing, rel.To, rel.Name)
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

func mergeViewProto(cat *Catalog, existing *metadatapb.ViewMetadata, rel *Relation, name string) *metadatapb.ViewMetadata {
	updated := relationToViewProto(cat, rel)
	updated.Name = name
	updated.Comment = existing.Comment
	updated.DependencyColumns = existing.DependencyColumns
	if len(updated.Columns) == 0 {
		updated.Columns = existing.Columns
	}
	updated.Rules = existing.Rules
	updated.SkipDump = existing.SkipDump
	preserveTriggerComments(updated.Triggers, existing.Triggers)
	return updated
}

func mergeMatViewProto(cat *Catalog, existing *metadatapb.MaterializedViewMetadata, rel *Relation, name string) *metadatapb.MaterializedViewMetadata {
	updated := relationToMatViewProto(cat, rel)
	updated.Name = name
	updated.Comment = existing.Comment
	updated.DependencyColumns = existing.DependencyColumns
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
			applyColumnDiffs(&ext.Columns, catAfter, rel)
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
	if len(dst.CheckConstraints) > 0 && len(src.CheckConstraints) > 0 {
		chkComments := make(map[string]string, len(src.CheckConstraints))
		for _, chk := range src.CheckConstraints {
			if chk.Comment != "" {
				chkComments[chk.Name] = chk.Comment
			}
		}
		for _, chk := range dst.CheckConstraints {
			if c, ok := chkComments[chk.Name]; ok && chk.Comment == "" {
				chk.Comment = c
			}
		}
	}
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
func renamePartitions(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, relations []RelationDiffEntry, canonicalSchema func(string) string) map[uint32]bool {
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
			p.Name = add.Name
			p.Indexes, p.CheckConstraints, p.ExcludeConstraints = tbl.Indexes, tbl.CheckConstraints, tbl.ExcludeConstraints
			if add.SchemaName != srcSchema {
				parentRel := catAfter.GetRelationByOID(add.To.PartitionOf)
				if parentRel != nil && parentRel.Schema != nil && parentRel.Schema.Name == add.SchemaName {
					p.Schema = ""
				} else {
					p.Schema = add.SchemaName
				}
			}
			renamed[rel.From.OID] = true
			continue
		}
		if sm == nil {
			continue
		}
		if rel.From.RelKind == 'v' {
			for i, existing := range sm.Views {
				if existing.Name == rel.Name {
					updated := mergeViewProto(catAfter, existing, add.To, add.Name)
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
			if rel.From.OwnerName != add.To.OwnerName {
				tbl.Owner = add.To.OwnerName
			}
			if add.SchemaName != srcSchema {
				sm.Tables = removeTableByName(sm.Tables, rel.Name)
				targetSchema := findOrCreateSchema(result, add.SchemaName)
				targetSchema.Tables = append(targetSchema.Tables, tbl)
			}
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
	updated.Options = existing.Options
	if fromCol != nil && fromCol.CollationName != toCol.CollationName {
		if toCol.CollationName != "" {
			updated.Collation = quoteIdentAlways(toCol.CollationName)
		}
	} else {
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
					if len(col.Options) == 0 {
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

func applyColumnDiffs(cols *[]*metadatapb.ColumnMetadata, cat *Catalog, rel RelationDiffEntry) {
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
		(idx.WhereClause != "" || nKey < len(idx.Columns) || len(idx.OpClasses) > 0 || idx.NullsNotDistinct) {
		im.Definition = buildCreateIndexSQL(idx, rel, rel.Schema.Name)
	}

	return im
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
	b.WriteString("AS $function$")
	b.WriteString(up.Body)
	b.WriteString("$function$\n")
	return b.String()
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
			fc.TypeOID == c.TypeOID && fc.TypeMod == c.TypeMod && fc.CollationName == c.CollationName {
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
			// The catalog keeps only the bare collation name, so a collation
			// introduced by walk-through DDL loses its schema qualifier here
			// (same limitation as below). Walk-through metadata is advisory —
			// the authoritative metadata is re-synced from the database after
			// execution, which stores the properly qualified reference.
			attribute.Collation = "\"" + c.CollationName + "\""
			// Substitute the previous emit-ready (possibly schema-qualified)
			// reference only when it names the same collation — a changed
			// collation must reflect the catalog's current value.
			if surviving != nil && collationReferenceBareName(surviving.Collation) == c.CollationName {
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
	}

	return table
}

// relationToViewProto converts an omni Relation (view) to a metadatapb.ViewMetadata.
func relationToViewProto(c *Catalog, rel *Relation) *metadatapb.ViewMetadata {
	view := &metadatapb.ViewMetadata{
		Name: rel.Name,
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
	}
	return view
}

// relationToMatViewProto converts an omni Relation (materialized view) to a metadatapb.MaterializedViewMetadata.
func relationToMatViewProto(c *Catalog, rel *Relation) *metadatapb.MaterializedViewMetadata {
	mv := &metadatapb.MaterializedViewMetadata{
		Name: rel.Name,
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
		Name:    name,
		Subtype: qualifyWalkThroughAttributeType(cat, rt.SubTypeOID, cat.FormatType(rt.SubTypeOID, -1)),
	}
	if rt.MultirangeOID != 0 {
		out.MultirangeName = qualifyWalkThroughAttributeType(cat, rt.MultirangeOID, cat.FormatType(rt.MultirangeOID, -1))
	}
	if previous != nil {
		out.SubtypeOpclass = previous.SubtypeOpclass
		out.Collation = previous.Collation
		out.Canonical = previous.Canonical
		out.SubtypeDiff = previous.SubtypeDiff
		out.Comment = previous.Comment
		out.SkipDump = previous.SkipDump
	}
	return out
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
					for _, col := range tbl.Columns {
						if col.Name == colName {
							col.Comment = c.To
							break
						}
					}
				} else if ext := findExternalTable(sm, relName); ext != nil {
					for _, col := range ext.Columns {
						if col.Name == colName {
							col.Comment = c.To
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
			sm := findSchema(result, schemaName)
			if sm == nil {
				continue
			}
			if tbl := findTable(sm, tableName); tbl != nil {
				setConstraintCommentOnTable(tbl, conName, c.To)
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

func setTriggerCommentInSchema(sm *metadatapb.SchemaMetadata, relName, trigName, comment string) {
	if tbl := findTable(sm, relName); tbl != nil {
		for _, tr := range tbl.Triggers {
			if tr.Name == trigName {
				tr.Comment = comment
				return
			}
		}
	}
	for _, v := range sm.Views {
		if v.Name == relName {
			for _, tr := range v.Triggers {
				if tr.Name == trigName {
					tr.Comment = comment
					return
				}
			}
		}
	}
	for _, mv := range sm.MaterializedViews {
		if mv.Name == relName {
			for _, tr := range mv.Triggers {
				if tr.Name == trigName {
					tr.Comment = comment
					return
				}
			}
		}
	}
}

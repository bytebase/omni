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
	if diff == nil || diff.IsEmpty() {
		return proto.CloneOf(original)
	}

	result, ok := proto.Clone(original).(*metadatapb.DatabaseSchemaMetadata)
	if !ok {
		return proto.CloneOf(original)
	}

	for _, s := range diff.Schemas {
		switch s.Action {
		case DiffAdd:
			result.Schemas = append(result.Schemas, &metadatapb.SchemaMetadata{Name: s.Name})
		case DiffDrop:
			result.Schemas = removeSchema(result.Schemas, s.Name)
		default:
		}
	}

	renamed := renamePartitions(result, catAfter, diff.Relations)
	for _, rel := range diff.Relations {
		if (rel.From != nil && renamed[rel.From.OID]) || (rel.To != nil && renamed[rel.To.OID]) {
			continue
		}
		if rel.Action == DiffDrop {
			if sm := findSchema(result, rel.SchemaName); sm != nil {
				dropRelation(sm, rel)
			}
			continue
		}
		sm := findOrCreateSchema(result, rel.SchemaName)
		switch rel.Action {
		case DiffAdd:
			addRelation(sm, catAfter, rel)
		case DiffModify:
			modifyRelation(sm, catBefore, catAfter, rel)
		default:
		}
	}

	for _, seq := range diff.Sequences {
		if seq.Action == DiffDrop {
			if sm := findSchema(result, seq.SchemaName); sm != nil {
				sm.Sequences = removeSequenceByName(sm.Sequences, seq.Name)
			}
			continue
		}
		sm := findOrCreateSchema(result, seq.SchemaName)
		switch seq.Action {
		case DiffAdd:
			if seq.To != nil {
				sm.Sequences = append(sm.Sequences, sequenceToProto(catAfter, seq.To))
			}
		case DiffModify:
			if seq.To != nil {
				for i, s := range sm.Sequences {
					if s.Name == seq.Name {
						sm.Sequences[i] = sequenceToProto(catAfter, seq.To)
						break
					}
				}
			}
		default:
		}
	}

	for _, e := range diff.Enums {
		if e.Action == DiffDrop {
			if sm := findSchema(result, e.SchemaName); sm != nil {
				sm.EnumTypes = removeEnumByName(sm.EnumTypes, e.Name)
			}
			continue
		}
		sm := findOrCreateSchema(result, e.SchemaName)
		switch e.Action {
		case DiffAdd:
			sm.EnumTypes = append(sm.EnumTypes, &metadatapb.EnumTypeMetadata{
				Name:   e.Name,
				Values: e.ToValues,
			})
		case DiffModify:
			for i, et := range sm.EnumTypes {
				if et.Name == e.Name {
					sm.EnumTypes[i].Values = e.ToValues
					break
				}
			}
		default:
		}
	}

	for _, ct := range diff.CompositeTypes {
		if ct.Action == DiffDrop {
			if sm := findSchema(result, ct.SchemaName); sm != nil {
				sm.CompositeTypes = removeCompositeTypeByName(sm.CompositeTypes, ct.Name)
			}
			continue
		}
		if ct.To == nil {
			continue
		}
		sm := findOrCreateSchema(result, ct.SchemaName)
		switch ct.Action {
		case DiffAdd:
			sm.CompositeTypes = append(sm.CompositeTypes, compositeTypeToProto(catAfter, ct.Name, nil, ct.To, nil))
		case DiffModify:
			for i, existing := range sm.CompositeTypes {
				if existing.Name == ct.Name {
					sm.CompositeTypes[i] = compositeTypeToProto(catAfter, ct.Name, ct.From, ct.To, existing)
					break
				}
			}
		default:
		}
	}

	for _, r := range diff.Ranges {
		if r.Action == DiffDrop {
			if sm := findSchema(result, r.SchemaName); sm != nil {
				sm.RangeTypes = removeRangeTypeByName(sm.RangeTypes, r.Name)
			}
			continue
		}
		if r.To == nil {
			continue
		}
		sm := findOrCreateSchema(result, r.SchemaName)
		switch r.Action {
		case DiffAdd:
			sm.RangeTypes = append(sm.RangeTypes, rangeTypeToProto(catAfter, r.To, nil))
		case DiffModify:
			for i, existing := range sm.RangeTypes {
				if existing.Name == r.Name {
					sm.RangeTypes[i] = rangeTypeToProto(catAfter, r.To, existing)
					break
				}
			}
		default:
		}
	}

	for _, f := range diff.Functions {
		isProcedure := f.To != nil && f.To.Kind == 'p'
		wasProcedure := f.From != nil && f.From.Kind == 'p'
		if f.Action == DiffDrop {
			if sm := findSchema(result, f.SchemaName); sm != nil {
				if wasProcedure {
					sm.Procedures = removeProcedureByIdentity(sm.Procedures, f.Identity)
				} else {
					sm.Functions = removeFunctionByIdentity(sm.Functions, f.Identity)
				}
			}
			continue
		}
		sm := findOrCreateSchema(result, f.SchemaName)
		switch f.Action {
		case DiffAdd:
			if f.To != nil {
				if isProcedure {
					sm.Procedures = append(sm.Procedures, procedureToProto(catAfter, f.To, f.Identity))
				} else {
					sm.Functions = append(sm.Functions, functionToProto(catAfter, f.To, f.Identity))
				}
			}
		case DiffModify:
			if wasProcedure {
				sm.Procedures = removeProcedureByIdentity(sm.Procedures, f.Identity)
			} else {
				sm.Functions = removeFunctionByIdentity(sm.Functions, f.Identity)
			}
			if f.To != nil {
				if isProcedure {
					sm.Procedures = append(sm.Procedures, procedureToProto(catAfter, f.To, f.Identity))
				} else {
					sm.Functions = append(sm.Functions, functionToProto(catAfter, f.To, f.Identity))
				}
			}
		default:
		}
	}

	applyCommentDiffs(result, catBefore, catAfter, diff.Comments)

	return result
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

func addRelation(sm *metadatapb.SchemaMetadata, cat *Catalog, rel RelationDiffEntry) {
	if rel.To == nil {
		return
	}
	switch rel.To.RelKind {
	case 'r', 'p', 'f':
		sm.Tables = append(sm.Tables, relationToTableProto(cat, rel.To))
	case 'v':
		sm.Views = append(sm.Views, relationToViewProto(cat, rel.To))
	case 'm':
		sm.MaterializedViews = append(sm.MaterializedViews, relationToMatViewProto(cat, rel.To))
	default:
	}
}

func dropRelation(sm *metadatapb.SchemaMetadata, rel RelationDiffEntry) {
	name := rel.Name
	sm.Tables = removeTableByName(sm.Tables, name)
	sm.Views = removeViewByName(sm.Views, name)
	sm.MaterializedViews = removeMatViewByName(sm.MaterializedViews, name)
	for _, t := range sm.Tables {
		t.Partitions = removePartitionByName(t.Partitions, name)
	}
}

func modifyRelation(sm *metadatapb.SchemaMetadata, catBefore, cat *Catalog, rel RelationDiffEntry) {
	if rel.To == nil {
		return
	}
	switch rel.To.RelKind {
	case 'r', 'p', 'f':
		tbl := findTable(sm, rel.Name)
		if tbl == nil {
			if p := findPartition(sm, rel.Name); p != nil {
				applyPartitionDiffs(p, catBefore, cat, rel)
				return
			}
			sm.Tables = append(sm.Tables, relationToTableProto(cat, rel.To))
			return
		}
		applyColumnDiffs(tbl, cat, rel)
		applyIndexDiffs(tbl, cat, rel)
		applyConstraintDiffs(tbl, catBefore, cat, rel)
		applyTriggerDiffs(tbl, cat, rel)
	case 'v':
		sm.Views = removeViewByName(sm.Views, rel.Name)
		sm.Views = append(sm.Views, relationToViewProto(cat, rel.To))
	case 'm':
		sm.MaterializedViews = removeMatViewByName(sm.MaterializedViews, rel.Name)
		sm.MaterializedViews = append(sm.MaterializedViews, relationToMatViewProto(cat, rel.To))
	default:
	}
}

// renamePartitions renames each partition or table the diff renames where it
// stands, preserving nested partitions and metadata the catalog does not track.
// The diff shows a rename as a drop and an add of the same relation, which would
// otherwise turn a partition into a table of its own or drop a table's
// partitions. It returns the OIDs it handled.
func renamePartitions(result *metadatapb.DatabaseSchemaMetadata, catAfter *Catalog, relations []RelationDiffEntry) map[uint32]bool {
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
		sm := findSchema(result, rel.SchemaName)
		if !ok || sm == nil {
			continue
		}
		if add.SchemaName == rel.SchemaName {
			if p := findPartition(sm, rel.Name); p != nil {
				tbl := relationToTableProto(catAfter, add.To)
				p.Name = add.Name
				p.Indexes, p.CheckConstraints, p.ExcludeConstraints = tbl.Indexes, tbl.CheckConstraints, tbl.ExcludeConstraints
				renamed[rel.From.OID] = true
				continue
			}
		}
		if tbl := findTable(sm, rel.Name); tbl != nil {
			rebuilt := relationToTableProto(catAfter, add.To)
			tbl.Name = add.Name
			tbl.Columns = rebuilt.Columns
			tbl.Indexes = rebuilt.Indexes
			tbl.ForeignKeys = rebuilt.ForeignKeys
			tbl.CheckConstraints = rebuilt.CheckConstraints
			tbl.ExcludeConstraints = rebuilt.ExcludeConstraints
			if add.SchemaName != rel.SchemaName {
				sm.Tables = removeTableByName(sm.Tables, rel.Name)
				targetSchema := findOrCreateSchema(result, add.SchemaName)
				targetSchema.Tables = append(targetSchema.Tables, tbl)
			}
			renamed[rel.From.OID] = true
		}
	}
	return renamed
}

// applyPartitionDiffs applies a partition's index and constraint changes. A
// partition's metadata holds only those; its columns are its parent table's.
func applyPartitionDiffs(p *metadatapb.TablePartitionMetadata, catBefore, cat *Catalog, rel RelationDiffEntry) {
	tbl := &metadatapb.TableMetadata{Indexes: p.Indexes, CheckConstraints: p.CheckConstraints, ExcludeConstraints: p.ExcludeConstraints}
	applyIndexDiffs(tbl, cat, rel)
	applyConstraintDiffs(tbl, catBefore, cat, rel)
	p.Indexes, p.CheckConstraints, p.ExcludeConstraints = tbl.Indexes, tbl.CheckConstraints, tbl.ExcludeConstraints
}

func applyColumnDiffs(tbl *metadatapb.TableMetadata, cat *Catalog, rel RelationDiffEntry) {
	for _, cd := range rel.Columns {
		switch cd.Action {
		case DiffAdd:
			if cd.To != nil {
				tbl.Columns = append(tbl.Columns, columnToProto(cat, cd.To))
			}
		case DiffDrop:
			tbl.Columns = removeColumnByName(tbl.Columns, cd.Name)
		case DiffModify:
			if cd.To != nil {
				for i, col := range tbl.Columns {
					if col.Name == cd.Name {
						updated := columnToProto(cat, cd.To)
						updated.Comment = col.Comment
						updated.CharacterSet = col.CharacterSet
						updated.Options = col.Options
						if cd.From != nil && cd.From.CollationName != cd.To.CollationName {
							if cd.To.CollationName != "" {
								updated.Collation = quoteIdentAlways(cd.To.CollationName)
							}
						} else {
							updated.Collation = col.Collation
							updated.CollationSchema = col.CollationSchema
						}
						tbl.Columns[i] = updated
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
						tbl.Indexes[i] = indexToProto(cat, rel.To, id.To)
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
			if cd.From != nil {
				removeConstraintFromTable(tbl, catBefore, cd.From)
			}
			if cd.To != nil {
				addConstraintToTable(tbl, catAfter, rel, cd.To)
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

func indexToProto(_ *Catalog, rel *Relation, idx *Index) *metadatapb.IndexMetadata {
	im := &metadatapb.IndexMetadata{
		Name:         idx.Name,
		Type:         idx.AccessMethod,
		Unique:       idx.IsUnique,
		Primary:      idx.IsPrimary,
		IsConstraint: idx.ConstraintOID != 0,
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

// findPartition returns the partition named name, at any depth, of a table in sm.
func findPartition(sm *metadatapb.SchemaMetadata, name string) *metadatapb.TablePartitionMetadata {
	var find func(partitions []*metadatapb.TablePartitionMetadata) *metadatapb.TablePartitionMetadata
	find = func(partitions []*metadatapb.TablePartitionMetadata) *metadatapb.TablePartitionMetadata {
		for _, p := range partitions {
			if p.Name == name {
				return p
			}
			if sub := find(p.Subpartitions); sub != nil {
				return sub
			}
		}
		return nil
	}
	for _, t := range sm.Tables {
		if p := find(t.Partitions); p != nil {
			return p
		}
	}
	return nil
}

// removePartitionByName removes the partition named name, at any depth.
func removePartitionByName(partitions []*metadatapb.TablePartitionMetadata, name string) []*metadatapb.TablePartitionMetadata {
	partitions = slices.DeleteFunc(partitions, func(p *metadatapb.TablePartitionMetadata) bool { return p.Name == name })
	for _, p := range partitions {
		p.Subpartitions = removePartitionByName(p.Subpartitions, name)
	}
	return partitions
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
		b.WriteString(up.Schema.Name)
		b.WriteByte('.')
	}
	b.WriteString(up.Name)
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

func removeProcedureByIdentity(procs []*metadatapb.ProcedureMetadata, identity string) []*metadatapb.ProcedureMetadata {
	out := make([]*metadatapb.ProcedureMetadata, 0, len(procs))
	for _, p := range procs {
		if p.Signature != identity {
			out = append(out, p)
		}
	}
	return out
}

func removeFunctionByIdentity(funcs []*metadatapb.FunctionMetadata, identity string) []*metadatapb.FunctionMetadata {
	out := make([]*metadatapb.FunctionMetadata, 0, len(funcs))
	for _, f := range funcs {
		if f.Signature != identity {
			out = append(out, f)
		}
	}
	return out
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
		Name: rel.Name,
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

func applyTriggerDiffs(tbl *metadatapb.TableMetadata, cat *Catalog, rel RelationDiffEntry) {
	if len(rel.Triggers) == 0 {
		return
	}
	schemaName := rel.SchemaName
	tableName := rel.Name
	for _, td := range rel.Triggers {
		switch td.Action {
		case DiffAdd:
			if td.To != nil {
				tbl.Triggers = append(tbl.Triggers, triggerToProto(cat, schemaName, tableName, td.To))
			}
		case DiffDrop:
			tbl.Triggers = removeTriggerByName(tbl.Triggers, td.Name)
		case DiffModify:
			tbl.Triggers = removeTriggerByName(tbl.Triggers, td.Name)
			if td.To != nil {
				tbl.Triggers = append(tbl.Triggers, triggerToProto(cat, schemaName, tableName, td.To))
			}
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
	out := &metadatapb.RangeTypeMetadata{
		Name:    rt.Name,
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

func removeRangeTypeByName(ranges []*metadatapb.RangeTypeMetadata, name string) []*metadatapb.RangeTypeMetadata {
	out := make([]*metadatapb.RangeTypeMetadata, 0, len(ranges))
	for _, r := range ranges {
		if r.Name != name {
			out = append(out, r)
		}
	}
	return out
}

func applyCommentDiffs(result *metadatapb.DatabaseSchemaMetadata, catBefore, catAfter *Catalog, comments []CommentDiffEntry) {
	for _, c := range comments {
		switch c.ObjType {
		case 'n':
			if sm := findSchema(result, c.ObjDescription); sm != nil {
				sm.Comment = c.To
			}
		case 'r':
			schemaName, relName, ok := strings.Cut(c.ObjDescription, ".")
			if !ok {
				continue
			}
			sm := findSchema(result, schemaName)
			if sm == nil {
				continue
			}
			if c.SubID != 0 {
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
				}
				continue
			}
			if tbl := findTable(sm, relName); tbl != nil {
				tbl.Comment = c.To
			} else if p := findPartition(sm, relName); p != nil {
				p.Comment = c.To
			} else {
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
			schemaName, idxName, ok := strings.Cut(c.ObjDescription, ".")
			if !ok {
				continue
			}
			sm := findSchema(result, schemaName)
			if sm == nil {
				continue
			}
			setIndexCommentInSchema(sm, idxName, c.To)
		case 't':
			schemaName, typeName, ok := strings.Cut(c.ObjDescription, ".")
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
			schemaName, seqName, ok := strings.Cut(c.ObjDescription, ".")
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
			parts := strings.SplitN(c.ObjDescription, ".", 3)
			if len(parts) != 3 {
				continue
			}
			sm := findSchema(result, parts[0])
			if sm == nil {
				continue
			}
			if tbl := findTable(sm, parts[1]); tbl != nil {
				setConstraintCommentOnTable(tbl, parts[2], c.To)
			}
		case 'g':
			parts := strings.SplitN(c.ObjDescription, ".", 3)
			if len(parts) != 3 {
				continue
			}
			sm := findSchema(result, parts[0])
			if sm == nil {
				continue
			}
			setTriggerCommentInSchema(sm, parts[1], parts[2], c.To)
		case 'f':
			schemaName, identity, ok := strings.Cut(c.ObjDescription, ".")
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
			for _, et := range result.EventTriggers {
				if et.Name == c.ObjDescription {
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

func setIndexCommentInSchema(sm *metadatapb.SchemaMetadata, idxName, comment string) {
	var setInPartitions func(partitions []*metadatapb.TablePartitionMetadata) bool
	setInPartitions = func(partitions []*metadatapb.TablePartitionMetadata) bool {
		for _, p := range partitions {
			for _, idx := range p.Indexes {
				if idx.Name == idxName {
					idx.Comment = comment
					return true
				}
			}
			if setInPartitions(p.Subpartitions) {
				return true
			}
		}
		return false
	}
	for _, tbl := range sm.Tables {
		for _, idx := range tbl.Indexes {
			if idx.Name == idxName {
				idx.Comment = comment
				return
			}
		}
		if setInPartitions(tbl.Partitions) {
			return
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

package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generatePartitionDDL produces DDL operations for partition-related features,
// replica identity changes, view check options, and table inheritance.
// It supplements generateTableDDL which handles basic CREATE/DROP TABLE.
func generatePartitionDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Relations {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}

			// Partition child: generate CREATE TABLE ... PARTITION OF ... FOR VALUES ...
			if entry.To.PartitionBound != nil && entry.To.PartitionOf != 0 {
				ops = append(ops, generatePartitionChildDDL(from, to, entry)...)
			}
			if (entry.To.RelKind == 'r' || entry.To.RelKind == 'p') && normalizeReplicaIdentity(entry.To.ReplicaIdentity) != 'd' {
				ops = append(ops, generateReplicaIdentityDDL(to, entry)...)
			}

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}

			// RelKind changed between table types (e.g., regular 'r' → partitioned 'p').
			// PG does not support ALTER TABLE to change RelKind. We must
			// DROP + CREATE the table with a warning about data loss.
			// Only handle table-type RelKinds; view/matview changes are handled by generateViewDDL.
			fromIsTable := entry.From.RelKind == 'r' || entry.From.RelKind == 'p'
			toIsTable := entry.To.RelKind == 'r' || entry.To.RelKind == 'p'
			// tableRebuilt says the same for the other generators.
			if fromIsTable && toIsTable && entry.From.RelKind != entry.To.RelKind {
				ops = append(ops, buildTableRecreateOps(from, to, diff, entry,
					fmt.Sprintf("converting table from relkind '%c' to '%c'", entry.From.RelKind, entry.To.RelKind))...)
				continue
			}

			// Partition attach/detach/reattach: generate ALTER TABLE ... DETACH/ATTACH
			// PARTITION rather than dropping and recreating the table.
			if fromIsTable && toIsTable && (isPartitionRelation(entry.From) || isPartitionRelation(entry.To)) {
				if !inhParentsEqual(from, to, entry.From.InhParents, entry.To.InhParents) || !partitionBoundEqual(entry.From.PartitionBound, entry.To.PartitionBound) {
					ops = append(ops, generatePartitionAttachmentOps(from, to, entry)...)
				}
			} else if fromIsTable && toIsTable && !inhParentsEqual(from, to, entry.From.InhParents, entry.To.InhParents) {
				// Non-partition inheritance changed: use ALTER TABLE ... NO INHERIT / INHERIT
				// instead of dropping and recreating the table.
				ops = append(ops, generateTableInheritOps(from, to, entry)...)
			}

			// Replica identity changed, or its index is rebuilt: a new index
			// does not inherit the replica identity of the one it replaces.
			if normalizeReplicaIdentity(entry.From.ReplicaIdentity) != normalizeReplicaIdentity(entry.To.ReplicaIdentity) ||
				replicaIdentityIndexName(from, entry.From) != replicaIdentityIndexName(to, entry.To) ||
				replicaIdentityIndexRebuilt(to, entry) {
				ops = append(ops, generateReplicaIdentityDDL(to, entry)...)
			}

			// Table access method changed (views/matviews are handled by generateViewDDL).
			if fromIsTable && toIsTable {
				if fromAM, toAM := normalizeRelationAccessMethod(entry.From), normalizeRelationAccessMethod(entry.To); fromAM != toAM {
					qn := formatRelationName(entry.SchemaName, entry.Name, entry.To)
					target := "DEFAULT"
					if !(entry.To.RelKind == 'p' && toAM == "") {
						target = quoteIdentifier(toAM)
					}
					phase := PhaseMain
					objOID := entry.To.OID
					if fromAM != "" && fromAM != "heap" && to.accessMethodByName[fromAM] == nil {
						phase = PhasePre
						objOID = entry.From.OID
					}
					ops = append(ops, MigrationOp{
						Type:          OpAlterTable,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("ALTER TABLE %s SET ACCESS METHOD %s", qn, target),
						Transactional: true,
						Phase:         phase,
						ObjType:       'r',
						ObjOID:        objOID,
						Priority:      PriorityTable,
					})
				}
			}
		}
	}

	rebuiltServers := rebuiltForeignServerNames(diff)
	diffFTKeys := make(map[string]bool, len(diff.ForeignTables))
	for _, entry := range diff.ForeignTables {
		diffFTKeys[entry.SchemaName+"\x00"+entry.Name] = true
		switch entry.Action {
		case DiffAdd:
			if entry.To != nil && entry.To.RelKind == 'f' && entry.To.PartitionBound != nil && entry.To.PartitionOf != 0 {
				if op, ok := buildForeignTableAttachOp(to, entry.SchemaName, entry.Name, entry.To); ok {
					ops = append(ops, op)
				}
			}
		case DiffModify:
			if entry.From == nil || entry.To == nil || entry.From.RelKind != 'f' || entry.To.RelKind != 'f' {
				continue
			}
			ftRebuilt := entry.From.ForeignServer != entry.To.ForeignServer || rebuiltServers[entry.From.ForeignServer]
			if ftRebuilt {
				if entry.To.PartitionOf != 0 {
					if op, ok := buildForeignTableAttachOp(to, entry.SchemaName, entry.Name, entry.To); ok {
						ops = append(ops, op)
					}
				}
				continue
			}
			if (isPartitionRelation(entry.From) || isPartitionRelation(entry.To)) &&
				(!inhParentsEqual(from, to, entry.From.InhParents, entry.To.InhParents) || !partitionBoundEqual(entry.From.PartitionBound, entry.To.PartitionBound)) {
				ops = append(ops, generatePartitionAttachmentOps(from, to, entry)...)
			}
		}
	}
	if from != nil && to != nil && len(rebuiltServers) > 0 {
		for _, s := range to.schemas {
			if s == nil {
				continue
			}
			for _, toRel := range s.Relations {
				if toRel == nil || toRel.RelKind != 'f' || toRel.PartitionOf == 0 || diffFTKeys[s.Name+"\x00"+toRel.Name] {
					continue
				}
				fromRel := from.GetRelation(s.Name, toRel.Name)
				if fromRel != nil && fromRel.RelKind == 'f' && rebuiltServers[fromRel.ForeignServer] {
					if op, ok := buildForeignTableAttachOp(to, s.Name, toRel.Name, toRel); ok {
						ops = append(ops, op)
					}
				}
			}
		}
	}

	// Deterministic ordering within Phase/Priority so table creation precedes post-create ALTERs.
	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Phase != ops[j].Phase {
			return ops[i].Phase < ops[j].Phase
		}
		if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}
		if ops[i].Type != ops[j].Type {
			return ops[i].Type < ops[j].Type
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

func isPartitionRelation(rel *Relation) bool {
	return rel != nil && (rel.IsPartition || rel.PartitionOf != 0)
}

func relationSchemaOrDefault(rel *Relation, fallback string) string {
	if rel != nil && rel.Schema != nil && rel.Schema.Name != "" {
		return rel.Schema.Name
	}
	return fallback
}

func generatePartitionAttachmentOps(from, to *Catalog, entry RelationDiffEntry) []MigrationOp {
	var ops []MigrationOp
	childQN := migrationQualifiedName(entry.SchemaName, entry.Name)

	if entry.From != nil && entry.From.PartitionOf != 0 {
		if parentRel := from.GetRelationByOID(entry.From.PartitionOf); parentRel != nil {
			parentSchema := relationSchemaOrDefault(parentRel, entry.SchemaName)
			ops = append(ops, MigrationOp{
				Type:            OpAlterTable,
				SchemaName:      parentSchema,
				ObjectName:      parentRel.Name,
				SQL:             fmt.Sprintf("ALTER TABLE %s DETACH PARTITION %s", migrationQualifiedName(parentSchema, parentRel.Name), childQN),
				Transactional:   true,
				Phase:           PhasePre,
				ObjType:         'r',
				ObjOID:          entry.From.OID,
				Priority:        PriorityTrigger,
				secondarySchema: entry.SchemaName,
			})
		}
	}

	if entry.To != nil && entry.To.PartitionOf != 0 {
		if op, ok := buildForeignTableAttachOp(to, entry.SchemaName, entry.Name, entry.To); ok {
			ops = append(ops, op)
		}
	}

	return ops
}

func buildForeignTableAttachOp(to *Catalog, schemaName, relName string, rel *Relation) (MigrationOp, bool) {
	if to == nil || rel == nil || rel.PartitionOf == 0 {
		return MigrationOp{}, false
	}
	parentRel := to.GetRelationByOID(rel.PartitionOf)
	if parentRel == nil {
		return MigrationOp{}, false
	}
	childQN := migrationQualifiedName(schemaName, relName)
	parentSchema := relationSchemaOrDefault(parentRel, schemaName)
	boundClause := formatPartitionBound(rel.PartitionBound)
	sql := fmt.Sprintf("ALTER TABLE %s ATTACH PARTITION %s", migrationQualifiedName(parentSchema, parentRel.Name), childQN)
	if boundClause != "" {
		sql += " " + boundClause
	}
	return MigrationOp{
		Type:            OpAlterTable,
		SchemaName:      parentSchema,
		ObjectName:      parentRel.Name,
		SQL:             sql,
		Transactional:   true,
		Phase:           PhaseMain,
		ObjType:         'r',
		ObjOID:          rel.OID,
		Priority:        PriorityIndex,
		secondarySchema: schemaName,
	}, true
}

// tableRebuilt reports whether the plan drops and recreates a table: its kind
// changes between plain and partitioned, which ALTER TABLE cannot do.
func tableRebuilt(from, to *Catalog, fromRel, toRel *Relation) bool {
	fromIsTable := fromRel.RelKind == 'r' || fromRel.RelKind == 'p'
	toIsTable := toRel.RelKind == 'r' || toRel.RelKind == 'p'
	return fromIsTable && toIsTable && fromRel.RelKind != toRel.RelKind
}

func generateTableInheritOps(from, to *Catalog, entry RelationDiffEntry) []MigrationOp {
	type inheritParent struct {
		schema string
		name   string
		qn     string
	}
	resolveParents := func(c *Catalog, oids []uint32) []inheritParent {
		parents := make([]inheritParent, 0, len(oids))
		for _, oid := range oids {
			rel := c.GetRelationByOID(oid)
			if rel == nil || rel.Schema == nil {
				continue
			}
			parents = append(parents, inheritParent{
				schema: rel.Schema.Name,
				name:   rel.Name,
				qn:     formatRelationTargetName(rel.Schema.Name, rel),
			})
		}
		return parents
	}

	fromParents := resolveParents(from, entry.From.InhParents)
	toParents := resolveParents(to, entry.To.InhParents)
	fromChildQN := formatRelationName(entry.SchemaName, entry.Name, entry.From)
	toChildQN := formatRelationName(entry.SchemaName, entry.Name, entry.To)

	var ops []MigrationOp
	j := 0
	for _, fp := range fromParents {
		if j < len(toParents) && fp.schema == toParents[j].schema && fp.name == toParents[j].name {
			j++
			continue
		}
		ops = append(ops, MigrationOp{
			Type:            OpAlterTable,
			SchemaName:      entry.SchemaName,
			ObjectName:      entry.Name,
			SQL:             fmt.Sprintf("ALTER TABLE %s NO INHERIT %s", fromChildQN, fp.qn),
			Transactional:   true,
			Phase:           PhasePre,
			ObjType:         'r',
			ObjOID:          entry.From.OID,
			Priority:        PriorityTrigger,
			secondarySchema: fp.schema,
		})
	}
	for k := j; k < len(toParents); k++ {
		tp := toParents[k]
		ops = append(ops, MigrationOp{
			Type:            OpAlterTable,
			SchemaName:      entry.SchemaName,
			ObjectName:      entry.Name,
			SQL:             fmt.Sprintf("ALTER TABLE %s INHERIT %s", toChildQN, tp.qn),
			Transactional:   true,
			Phase:           PhaseMain,
			ObjType:         'r',
			ObjOID:          entry.To.OID,
			Priority:        PriorityIndex,
			secondarySchema: tp.schema,
		})
	}
	return ops
}

// generatePartitionChildDDL generates CREATE TABLE ... PARTITION OF ... FOR VALUES ...
func generatePartitionChildDDL(from, c *Catalog, entry RelationDiffEntry) []MigrationOp {
	rel := entry.To
	parentRel := c.GetRelationByOID(rel.PartitionOf)
	if parentRel == nil {
		return nil
	}

	qn := formatRelationName(entry.SchemaName, rel.Name, rel)

	var parentSchema string
	if parentRel.Schema != nil {
		parentSchema = parentRel.Schema.Name
	} else {
		parentSchema = entry.SchemaName
	}
	parentQN := formatRelationName(parentSchema, parentRel.Name, parentRel)

	var b strings.Builder
	switch rel.Persistence {
	case 'u':
		b.WriteString("CREATE UNLOGGED TABLE ")
	case 't':
		b.WriteString("CREATE TEMP TABLE ")
	default:
		b.WriteString("CREATE TABLE ")
	}
	b.WriteString(qn)
	b.WriteString(" PARTITION OF ")
	b.WriteString(parentQN)
	if boundClause := formatPartitionBound(rel.PartitionBound); boundClause != "" {
		b.WriteString(" ")
		b.WriteString(boundClause)
	}
	// A partition with subpartitions is itself partitioned.
	b.WriteString(formatPartitionByClause(rel))
	if rel.RelKind == 'p' {
		if rel.AccessMethod != "" {
			b.WriteString(" USING ")
			b.WriteString(quoteIdentifier(rel.AccessMethod))
		}
	} else if am := normalizeRelationAccessMethod(rel); am != "" {
		if (parentRel.AccessMethod != "" && am != parentRel.AccessMethod) ||
			(parentRel.AccessMethod == "" && (am != "heap" || rel.AccessMethod != "" || catalogHasNonDefaultTableAM(from) || catalogHasNonDefaultTableAM(c))) {
			b.WriteString(" USING ")
			b.WriteString(quoteIdentifier(am))
		}
	}

	ops := []MigrationOp{{
		Type:            OpCreateTable,
		SchemaName:      entry.SchemaName,
		ObjectName:      entry.Name,
		SQL:             b.String(),
		Transactional:   true,
		Phase:           PhaseMain,
		ObjType:         'r',
		ObjOID:          rel.OID,
		Priority:        PriorityTable,
		secondarySchema: parentSchema,
	}}
	if ownerOp, ok := buildTableOwnerOp(entry.SchemaName, entry.Name, rel); ok {
		ops = append(ops, ownerOp)
	}
	return ops
}

// formatPartitionValue formats a single partition bound value.
// Values that look like special keywords (MINVALUE, MAXVALUE, NULL) are passed through.
// String values are single-quoted; numeric values are passed through.
func formatPartitionValue(v string) string {
	upper := strings.ToUpper(v)
	if upper == "MINVALUE" || upper == "MAXVALUE" || upper == "NULL" {
		return upper
	}
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

func formatPartitionDatum(v string, idx int, keywords []bool) string {
	if idx < len(keywords) {
		if keywords[idx] {
			return strings.ToUpper(v)
		}
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	}
	return formatPartitionValue(v)
}

// generateReplicaIdentityDDL generates ALTER TABLE ... REPLICA IDENTITY ...
// replicaIdentityIndexRebuilt reports whether the migration drops and
// recreates the index, or the constraint owning it, that the target table
// uses as its replica identity.
func replicaIdentityIndexRebuilt(to *Catalog, entry RelationDiffEntry) bool {
	name := replicaIdentityIndexName(to, entry.To)
	if name == "" {
		return false
	}
	for _, ie := range entry.Indexes {
		if ie.Name == name {
			return true
		}
	}
	for _, ce := range entry.Constraints {
		if ce.Name == name {
			return true
		}
	}
	return false
}

func generateReplicaIdentityDDL(c *Catalog, entry RelationDiffEntry) []MigrationOp {
	qn := formatRelationName(entry.SchemaName, entry.Name, entry.To)

	var clause string
	switch entry.To.ReplicaIdentity {
	case 'f':
		clause = "FULL"
	case 'n':
		clause = "NOTHING"
	case 'i':
		clause = "USING INDEX" // would need index name, handled as best-effort
		// Try to find the replica identity index.
		for _, idx := range c.IndexesOf(entry.To.OID) {
			if idx.IsReplicaIdent {
				clause = "USING INDEX " + quoteIdentAlways(idx.Name)
				break
			}
		}
	default: // 'd' or anything else
		clause = "DEFAULT"
	}

	// USING INDEX names an index this migration may create or rebuild, so it
	// waits for every index to exist.
	phase := PhaseMain
	if entry.To.ReplicaIdentity == 'i' {
		phase = PhasePost
	}
	return []MigrationOp{{
		Type:          OpAlterTable,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           fmt.Sprintf("ALTER TABLE %s REPLICA IDENTITY %s", qn, clause),
		Transactional: true,
		Phase:         phase,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityIndex + 1,
	}}
}

// formatPartitionByClause formats a PARTITION BY clause for a partitioned table.
// This is appended to the CREATE TABLE statement generated by FormatCreateTable.
func formatPartitionByClause(rel *Relation) string {
	if rel.PartitionInfo == nil {
		return ""
	}

	pi := rel.PartitionInfo
	var strategy string
	switch pi.Strategy {
	case 'r':
		strategy = "RANGE"
	case 'l':
		strategy = "LIST"
	case 'h':
		strategy = "HASH"
	default:
		return ""
	}

	// Resolve key column names from attnums.
	keyCols := columnAttnumsToQuotedNames(rel, pi.KeyAttNums)
	if len(keyCols) == 0 {
		return ""
	}

	return fmt.Sprintf(" PARTITION BY %s (%s)", strategy, strings.Join(keyCols, ", "))
}

// formatInheritsClause formats an INHERITS clause for table inheritance.
func formatInheritsClause(c *Catalog, rel *Relation) string {
	if len(rel.InhParents) == 0 {
		return ""
	}

	// Only emit INHERITS for regular inheritance (not partition children).
	if rel.IsPartition || rel.PartitionOf != 0 {
		return ""
	}

	var parents []string
	for _, parentOID := range rel.InhParents {
		parentRel := c.GetRelationByOID(parentOID)
		if parentRel == nil {
			continue
		}
		var schema string
		if parentRel.Schema != nil {
			schema = parentRel.Schema.Name
		}
		if schema != "" {
			parents = append(parents, formatRelationName(schema, parentRel.Name, parentRel))
		} else {
			parents = append(parents, quoteIdentAlways(parentRel.Name))
		}
	}

	if len(parents) == 0 {
		return ""
	}

	return " INHERITS (" + strings.Join(parents, ", ") + ")"
}

// buildTableRecreateOps generates DROP TABLE CASCADE + CREATE TABLE + CREATE INDEX ops
// and restores unchanged non-inline dependents (FKs, EXCLUDE constraints, triggers, RLS,
// policies, grants, comments, managed sequence metadata) for a table that must be recreated
// (RelKind change). Dependent views are handled by wrapColumnTypeChangesWithViewOps
// and rules by generateRuleDDL.
func buildTableRecreateOps(from, to *Catalog, diff *SchemaDiff, entry RelationDiffEntry, reason string) []MigrationOp {
	var ops []MigrationOp
	qn := formatRelationName(entry.SchemaName, entry.Name, entry.From)

	ops = append(ops, MigrationOp{
		Type:          OpDropTable,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           fmt.Sprintf("DROP TABLE %s CASCADE", qn),
		Warning:       fmt.Sprintf("%s requires DROP + CREATE (data loss): %s", qn, reason),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'r',
		ObjOID:        entry.From.OID,
		Priority:      PriorityTable,
	})

	ddl := formatCreateTable(from, to, entry.SchemaName, entry.From, entry.To, false)
	ops = append(ops, MigrationOp{
		Type:          OpCreateTable,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           ddl,
		Warning:       fmt.Sprintf("table %s recreated — verify data migration", qn),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityTable,
	})

	// Recreate standalone indexes (not constraint-backed) that exist in the target.
	indexes := to.IndexesOf(entry.To.OID)
	for _, idx := range indexes {
		if idx.ConstraintOID != 0 {
			continue // managed by constraint DDL
		}
		sql := buildCreateIndexSQL(idx, entry.To, entry.SchemaName)
		ops = append(ops, MigrationOp{
			Type:          OpCreateIndex,
			SchemaName:    entry.SchemaName,
			ObjectName:    idx.Name,
			SQL:           sql,
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       'i',
			ObjOID:        idx.OID,
			Priority:      PriorityIndex,
		})
	}
	if normalizeReplicaIdentity(entry.To.ReplicaIdentity) != 'd' {
		ops = append(ops, generateReplicaIdentityDDL(to, entry)...)
	}

	ops = append(ops, restoreRebuiltTableDependents(from, to, diff, entry)...)

	return ops
}

func restoreRebuiltTableDependents(from, to *Catalog, diff *SchemaDiff, entry RelationDiffEntry) []MigrationOp {
	var ops []MigrationOp
	rel := entry.To
	if rel == nil {
		return nil
	}

	handledCon := make(map[string]bool)
	for _, ce := range entry.Constraints {
		if ce.Action == DiffAdd || ce.Action == DiffModify {
			handledCon[ce.Name] = true
		}
	}
	for _, con := range to.ConstraintsOf(rel.OID) {
		if (con.Type != ConstraintFK && con.Type != ConstraintExclude) || con.ConParentID != 0 || handledCon[con.Name] {
			continue
		}
		ops = append(ops, buildAddConstraintOp(to, entry.SchemaName, entry.Name, con))
	}

	// Inbound FKs from kept, non-rebuilt relations are also dropped by DROP TABLE ... CASCADE.
	relDiffByKey := make(map[string]RelationDiffEntry)
	if diff != nil {
		for _, re := range diff.Relations {
			relDiffByKey[re.SchemaName+"."+re.Name] = re
		}
	}
	var schemaNames []string
	for _, sc := range to.UserSchemas() {
		if sc != nil && !isTempSchemaName(sc.Name) {
			schemaNames = append(schemaNames, sc.Name)
		}
	}
	sort.Strings(schemaNames)
	for _, scName := range schemaNames {
		sc := to.schemaByName[scName]
		if sc == nil {
			continue
		}
		relNames := make([]string, 0, len(sc.Relations))
		for rName := range sc.Relations {
			relNames = append(relNames, rName)
		}
		sort.Strings(relNames)
		for _, rName := range relNames {
			other := sc.Relations[rName]
			if other == nil || other.OID == rel.OID || (other.RelKind != 'r' && other.RelKind != 'p') {
				continue
			}
			otherKey := scName + "." + rName
			otherDiff, hasDiff := relDiffByKey[otherKey]
			if hasDiff && otherDiff.Action == DiffAdd {
				continue
			}
			if hasDiff && otherDiff.From != nil && otherDiff.To != nil && tableRebuilt(from, to, otherDiff.From, otherDiff.To) {
				continue
			}
			otherHandled := make(map[string]bool)
			if hasDiff {
				for _, ce := range otherDiff.Constraints {
					if ce.Action == DiffAdd || ce.Action == DiffModify {
						otherHandled[ce.Name] = true
					}
				}
			}
			for _, con := range to.ConstraintsOf(other.OID) {
				if con.Type == ConstraintFK && con.ConParentID == 0 && con.FRelOID == rel.OID && !otherHandled[con.Name] {
					ops = append(ops, buildAddConstraintOp(to, scName, rName, con))
				}
			}
		}
	}

	handledTrig := make(map[string]bool)
	for _, te := range entry.Triggers {
		if te.Action == DiffAdd || te.Action == DiffModify {
			handledTrig[te.Name] = true
		}
	}
	for _, trig := range to.TriggersOf(rel.OID) {
		if trig.ParentTriggerOID == 0 && !handledTrig[trig.Name] {
			ops = append(ops, buildCreateAndModeTriggerOps(to, entry.SchemaName, entry.Name, trig)...)
		}
	}

	targetQN := formatRelationName(entry.SchemaName, entry.Name, rel)
	if !entry.RLSChanged {
		if rel.RowSecurity {
			ops = append(ops, MigrationOp{
				Type:          OpAlterTable,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", targetQN),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'r',
				ObjOID:        rel.OID,
				Priority:      PriorityColumn,
			})
		}
		if rel.ForceRowSecurity {
			ops = append(ops, MigrationOp{
				Type:          OpAlterTable,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER TABLE %s FORCE ROW LEVEL SECURITY", targetQN),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'r',
				ObjOID:        rel.OID,
				Priority:      PriorityColumn,
			})
		}
	}

	handledPol := make(map[string]bool)
	for _, pe := range entry.Policies {
		if pe.Action == DiffAdd || pe.Action == DiffModify {
			handledPol[pe.Name] = true
		}
	}
	for _, pol := range to.policiesByRel[rel.OID] {
		if !handledPol[pol.Name] {
			ops = append(ops, buildCreatePolicyOpForRel(entry.SchemaName, entry.Name, rel, pol))
		}
	}

	handledGrantSQL := make(map[string]bool)
	if diff != nil {
		for _, ge := range diff.Grants {
			if (ge.Action == DiffAdd || ge.Action == DiffModify) &&
				((ge.To.ObjType == 'r' && ge.To.ObjOID == rel.OID) || ge.To.ObjType == 's') {
				handledGrantSQL[formatGrantSQL(ge.To, resolveGrantObjName(to, ge.To.ObjType, ge.To.ObjOID))] = true
			}
		}
	}
	for _, gop := range restoreGrantOps(to, 'r', rel.OID) {
		if !handledGrantSQL[gop.SQL] {
			ops = append(ops, gop)
		}
	}

	handledComments := make(map[string]bool)
	if diff != nil {
		for _, ce := range diff.Comments {
			if ce.Action == DiffAdd || ce.Action == DiffModify {
				if sql := formatCommentEntrySQL(to, ce, ce.To); sql != "" {
					handledComments[sql] = true
				}
			}
		}
	}
	appendComment := func(objType byte, oid uint32, subID int16, desc string, parts []string) {
		if txt, ok := to.GetComment(objType, oid, subID); ok && txt != "" {
			if sql := formatCommentEntrySQL(to, CommentDiffEntry{
				ObjType:        objType,
				ObjDescription: desc,
				SubID:          subID,
				objParts:       parts,
			}, txt); sql != "" && !handledComments[sql] {
				handledComments[sql] = true
				ops = append(ops, MigrationOp{
					Type:          OpComment,
					ObjectName:    desc,
					SQL:           sql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       objType,
					Priority:      PriorityMetadata,
				})
			}
		}
	}

	if relSc := to.schemaByName[entry.SchemaName]; relSc != nil {
		colByAttNum := make(map[int16]*Column, len(rel.Columns))
		for _, col := range rel.Columns {
			colByAttNum[col.AttNum] = col
		}
		var seqNames []string
		for sName, seq := range relSc.Sequences {
			if seq != nil && seq.OwnerRelOID == rel.OID && sequenceManagedByColumn(to, seq) {
				seqNames = append(seqNames, sName)
			}
		}
		sort.Strings(seqNames)
		for _, sName := range seqNames {
			seq := relSc.Sequences[sName]
			ownerCol := colByAttNum[seq.OwnerAttNum]
			if ownerCol == nil {
				continue
			}
			defTypeOID := ownerCol.TypeOID
			bounds, hasBounds := seqBounds[defTypeOID]
			if !hasBounds {
				defTypeOID = seq.TypeOID
				bounds = seqBounds[defTypeOID]
			}
			defaultSeqName := rel.Name + "_" + ownerCol.Name + "_seq"
			if ownerCol.Identity != 0 && seq.Name != defaultSeqName {
				defaultTarget := migrationQualifiedName(entry.SchemaName, defaultSeqName)
				if rel.Persistence == 't' || isTempSchemaName(entry.SchemaName) {
					defaultTarget = quoteIdentAlways(defaultSeqName)
				}
				ops = append(ops, MigrationOp{
					Type:          OpAlterSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    seq.Name,
					SQL:           fmt.Sprintf("ALTER SEQUENCE %s RENAME TO %s", defaultTarget, quoteIdentAlways(seq.Name)),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'S',
					ObjOID:        seq.OID,
					Priority:      PriorityColumn,
				})
			}
			defaultSeq := Sequence{
				Name:        seq.Name,
				Schema:      seq.Schema,
				Persistence: normalizeRelPersistence(rel.Persistence),
				TypeOID:     defTypeOID,
				Increment:   1,
				MinValue:    1,
				MaxValue:    bounds[1],
				Start:       1,
				CacheValue:  1,
				OwnerRelOID: rel.OID,
				OwnerAttNum: ownerCol.AttNum,
			}
			if alterSQL := formatAlterSequence(to, to, entry.SchemaName, &defaultSeq, seq); alterSQL != "" {
				ops = append(ops, MigrationOp{
					Type:          OpAlterSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    seq.Name,
					SQL:           alterSQL,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'S',
					ObjOID:        seq.OID,
					Priority:      PriorityConstraint,
				})
			}
			seqPers := normalizeRelPersistence(seq.Persistence)
			relPers := normalizeRelPersistence(rel.Persistence)
			if seqPers != relPers && relPers != 't' && (seqPers == 'u' || seqPers == 'p') {
				mode := "LOGGED"
				if seqPers == 'u' {
					mode = "UNLOGGED"
				}
				ops = append(ops, MigrationOp{
					Type:          OpAlterSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    seq.Name,
					SQL:           fmt.Sprintf("ALTER SEQUENCE %s SET %s", formatSequenceTargetName(entry.SchemaName, seq), mode),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'S',
					ObjOID:        seq.OID,
					Priority:      PriorityConstraint,
				})
			}
			for _, gop := range restoreGrantOps(to, 's', seq.OID) {
				if !handledGrantSQL[gop.SQL] {
					ops = append(ops, gop)
				}
			}
			appendComment('s', seq.OID, 0, entry.SchemaName+"."+seq.Name, []string{entry.SchemaName, seq.Name})
		}
	}

	relDesc := entry.SchemaName + "." + entry.Name
	appendComment('r', rel.OID, 0, relDesc, []string{entry.SchemaName, entry.Name})
	for _, col := range rel.Columns {
		appendComment('r', rel.OID, col.AttNum, relDesc, []string{entry.SchemaName, entry.Name})
	}
	for _, idx := range to.IndexesOf(rel.OID) {
		appendComment('i', idx.OID, 0, entry.SchemaName+"."+idx.Name, []string{entry.SchemaName, idx.Name})
	}
	for _, con := range to.ConstraintsOf(rel.OID) {
		appendComment('c', con.OID, 0, relDesc+"."+con.Name, []string{entry.SchemaName, entry.Name, con.Name})
	}
	for _, trig := range to.TriggersOf(rel.OID) {
		if trig.ParentTriggerOID == 0 {
			appendComment('g', trig.OID, 0, relDesc+"."+trig.Name, []string{entry.SchemaName, entry.Name, trig.Name})
		}
	}
	for _, pol := range to.policiesByRel[rel.OID] {
		appendComment('p', pol.OID, 0, relDesc+"."+pol.Name, []string{entry.SchemaName, entry.Name, pol.Name})
	}
	for _, r := range to.rules {
		if r != nil && r.OwnerRelOID == rel.OID {
			appendComment('R', r.OID, 0, relDesc+"."+r.Name, []string{entry.SchemaName, entry.Name, r.Name})
		}
	}

	return ops
}

// formatPartitionBound formats a partition bound clause (e.g. "FOR VALUES IN (...)", "DEFAULT").
func formatPartitionBound(bound *PartitionBound) string {
	if bound == nil {
		return ""
	}
	if bound.IsDefault {
		return "DEFAULT"
	}
	var b strings.Builder
	switch bound.Strategy {
	case 'l':
		b.WriteString("FOR VALUES IN (")
		for i, v := range bound.ListValues {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatPartitionDatum(v, i, bound.ListKeywords))
		}
		b.WriteString(")")
	case 'r':
		b.WriteString("FOR VALUES FROM (")
		for i, v := range bound.LowerBound {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatPartitionDatum(v, i, bound.LowerKeywords))
		}
		b.WriteString(") TO (")
		for i, v := range bound.UpperBound {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatPartitionDatum(v, i, bound.UpperKeywords))
		}
		b.WriteString(")")
	case 'h':
		fmt.Fprintf(&b, "FOR VALUES WITH (MODULUS %d, REMAINDER %d)", bound.Modulus, bound.Remainder)
	}
	return b.String()
}

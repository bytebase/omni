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
				ops = append(ops, generatePartitionChildDDL(to, entry)...)
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
			if fromIsTable && toIsTable && entry.From.RelKind != entry.To.RelKind {
				ops = append(ops, buildTableRecreateOps(from, to, entry,
					fmt.Sprintf("converting table from relkind '%c' to '%c'", entry.From.RelKind, entry.To.RelKind))...)
				continue
			}

			// Inheritance changed: PG does not support ALTER TABLE ... INHERITS.
			// We must DROP + CREATE the table. Only for actual tables.
			if fromIsTable && toIsTable && !inhParentsEqual(from, to, entry.From.InhParents, entry.To.InhParents) {
				ops = append(ops, buildTableRecreateOps(from, to, entry, "changing table inheritance")...)
				continue
			}

			// Replica identity changed.
			if entry.From.ReplicaIdentity != entry.To.ReplicaIdentity {
				ops = append(ops, generateReplicaIdentityDDL(to, entry)...)
			}
		}
	}

	// Deterministic ordering.
	sort.Slice(ops, func(i, j int) bool {
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

// generatePartitionChildDDL generates CREATE TABLE ... PARTITION OF ... FOR VALUES ...
func generatePartitionChildDDL(c *Catalog, entry RelationDiffEntry) []MigrationOp {
	rel := entry.To
	parentRel := c.GetRelationByOID(rel.PartitionOf)
	if parentRel == nil {
		return nil
	}

	qn := migrationQualifiedName(entry.SchemaName, rel.Name)

	var parentSchema string
	if parentRel.Schema != nil {
		parentSchema = parentRel.Schema.Name
	} else {
		parentSchema = entry.SchemaName
	}
	parentQN := migrationQualifiedName(parentSchema, parentRel.Name)

	var b strings.Builder
	b.WriteString("CREATE TABLE ")
	b.WriteString(qn)
	b.WriteString(" PARTITION OF ")
	b.WriteString(parentQN)

	bound := rel.PartitionBound
	if bound.IsDefault {
		b.WriteString(" DEFAULT")
	} else {
		switch bound.Strategy {
		case 'l': // LIST
			b.WriteString(" FOR VALUES IN (")
			for i, v := range bound.ListValues {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(formatPartitionValue(v))
			}
			b.WriteString(")")
		case 'r': // RANGE
			b.WriteString(" FOR VALUES FROM (")
			for i, v := range bound.LowerBound {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(formatPartitionValue(v))
			}
			b.WriteString(") TO (")
			for i, v := range bound.UpperBound {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(formatPartitionValue(v))
			}
			b.WriteString(")")
		case 'h': // HASH
			b.WriteString(fmt.Sprintf(" FOR VALUES WITH (MODULUS %d, REMAINDER %d)", bound.Modulus, bound.Remainder))
		}
	}

	return []MigrationOp{{
		Type:          OpCreateTable,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        rel.OID,
		Priority:      PriorityTable,
	}}
}

// formatPartitionValue formats a single partition bound value.
// Values that look like special keywords (MINVALUE, MAXVALUE) are passed through.
// String values are single-quoted; numeric values are passed through.
func formatPartitionValue(v string) string {
	upper := strings.ToUpper(v)
	if upper == "MINVALUE" || upper == "MAXVALUE" {
		return upper
	}
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

// generateReplicaIdentityDDL generates ALTER TABLE ... REPLICA IDENTITY ...
func generateReplicaIdentityDDL(c *Catalog, entry RelationDiffEntry) []MigrationOp {
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)

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

	return []MigrationOp{{
		Type:          OpAlterTable,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           fmt.Sprintf("ALTER TABLE %s REPLICA IDENTITY %s", qn, clause),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityTable,
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
			parents = append(parents, migrationQualifiedName(schema, parentRel.Name))
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
// for a table that must be recreated (e.g., RelKind change, inheritance change).
// DROP TABLE CASCADE also destroys dependent views and indexes, so indexes must
// be regenerated. Dependent views are handled by wrapColumnTypeChangesWithViewOps.
func buildTableRecreateOps(from, to *Catalog, entry RelationDiffEntry, reason string) []MigrationOp {
	var ops []MigrationOp
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)

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

	ddl := FormatCreateTable(to, entry.SchemaName, entry.To)
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

	return ops
}

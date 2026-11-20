package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateSequenceDDL produces CREATE SEQUENCE, DROP SEQUENCE, and ALTER SEQUENCE
// operations from the diff. Sequences owned by a column (SERIAL) are skipped
// because the diff layer already filters them out (OwnerRelOID != 0).
func generateSequenceDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Sequences {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			// Skip creating sequences that are owned by identity columns.
			// ADD GENERATED AS IDENTITY automatically creates the sequence.
			if isIdentitySequence(to, entry.SchemaName, entry.Name) {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpCreateSequence,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           formatCreateSequence(to, entry.SchemaName, entry.To),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'S',
				ObjOID:        entry.To.OID,
				Priority:      PrioritySequence,
			})
			if owner := sequenceOwner(to, entry.To); owner != "" {
				ops = append(ops, sequenceOwnedByOp(entry.SchemaName, entry.To, owner, PhasePost))
			}

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			// Skip dropping sequences that are owned by identity columns.
			// DROP IDENTITY automatically drops the identity sequence.
			if isIdentitySequence(from, entry.SchemaName, entry.Name) {
				continue
			}
			// Dropping the owning table or column drops the sequence, and must:
			// the owner's default may still draw from it.
			if ownerDropped(from, to, entry.From) {
				continue
			}
			qn := formatSequenceTargetName(entry.SchemaName, entry.From)
			phase := PhasePre
			if sequenceDrawnByKeptColumn(from, to, entry.From) {
				// The ALTER COLUMN that changes that default runs first.
				phase = PhaseLateDrop
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropSequence,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP SEQUENCE %s", qn),
				Transactional: true,
				Phase:         phase,
				ObjType:       preSequenceObjType,
				ObjOID:        entry.From.OID,
				Priority:      PrioritySequence,
			})

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			fromPers := normalizeRelPersistence(entry.From.Persistence)
			toPers := normalizeRelPersistence(entry.To.Persistence)
			fromOwner, toOwner := sequenceOwner(from, entry.From), sequenceOwner(to, entry.To)
			if fromPers != toPers && (fromPers == 't' || toPers == 't') {
				if fromOwner != "" && !ownerDropped(from, to, entry.From) {
					ops = append(ops, sequenceOwnedByOp(entry.SchemaName, entry.From, "NONE", PhasePre))
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("DROP SEQUENCE %s", formatSequenceTargetName(entry.SchemaName, entry.From)),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       preSequenceObjType,
					ObjOID:        entry.From.OID,
					Priority:      PrioritySequence,
				})
				ops = append(ops, MigrationOp{
					Type:          OpCreateSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           formatCreateSequence(to, entry.SchemaName, entry.To),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'S',
					ObjOID:        entry.To.OID,
					Priority:      PrioritySequence,
				})
				if toOwner != "" {
					ops = append(ops, sequenceOwnedByOp(entry.SchemaName, entry.To, toOwner, PhasePost))
				}
				// Unchanged grants make no diff, so give them back.
				ops = append(ops, restoreGrantOps(to, 's', entry.To.OID)...)
				continue
			}
			sql := formatAlterSequence(from, to, entry.SchemaName, entry.From, entry.To)
			if fromOwner != toOwner || fromOwner != "" && ownerDropped(from, to, entry.From) {
				// Release the old owner first: dropping it would take the sequence.
				if fromOwner != "" {
					ops = append(ops, sequenceOwnedByOp(entry.SchemaName, entry.From, "NONE", PhasePre))
				}
				if toOwner != "" {
					ops = append(ops, sequenceOwnedByOp(entry.SchemaName, entry.To, toOwner, PhasePost))
				}
			}
			if fromPers != toPers {
				mode := "LOGGED"
				if toPers == 'u' {
					mode = "UNLOGGED"
				}
				ops = append(ops, MigrationOp{
					Type:          OpAlterSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER SEQUENCE %s SET %s", formatSequenceTargetName(entry.SchemaName, entry.To), mode),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'S',
					ObjOID:        entry.To.OID,
					Priority:      PrioritySequence,
				})
			}
			if sql != "" {
				ops = append(ops, MigrationOp{
					Type:          OpAlterSequence,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           sql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'S',
					ObjOID:        entry.To.OID,
					Priority:      PrioritySequence,
				})
			}
		}
	}

	// Deterministic ordering: drops first, then creates, then alters; within each group by (schema, name).
	sort.Slice(ops, func(i, j int) bool {
		oi, oj := seqOpOrder(ops[i].Type), seqOpOrder(ops[j].Type)
		if oi != oj {
			return oi < oj
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// sequenceDrawnByKeptColumn reports whether a column default in from draws
// on seq and to still has that column.
func sequenceDrawnByKeptColumn(from, to *Catalog, seq *Sequence) bool {
	for _, d := range from.deps {
		if d.RefType != 's' || d.RefOID != seq.OID || d.ObjType != 'r' || d.ObjSubID == 0 {
			continue
		}
		rel := from.GetRelationByOID(d.ObjOID)
		if rel == nil || rel.Schema == nil {
			continue
		}
		if toRel := to.GetRelation(rel.Schema.Name, rel.Name); toRel != nil {
			if _, kept := toRel.colByName[resolveAttNum(rel, int16(d.ObjSubID))]; kept {
				return true
			}
		}
	}
	return false
}

// ownerDropped reports whether the plan drops the column that owns seq in
// from: it is gone in to, or its table is rebuilt.
func ownerDropped(from, to *Catalog, seq *Sequence) bool {
	rel := from.GetRelationByOID(seq.OwnerRelOID)
	if rel == nil || rel.Schema == nil {
		return false
	}
	toRel := to.GetRelation(rel.Schema.Name, rel.Name)
	if toRel == nil || tableRebuilt(from, to, rel, toRel) {
		return true
	}
	_, kept := toRel.colByName[resolveAttNum(rel, seq.OwnerAttNum)]
	return !kept
}

// preSequenceObjType is the ObjType of a sequence op among the drops: 's',
// the letter the catalog's dependencies use, so a drop or OWNED BY NONE sorts
// before dropping the table that owns the sequence, which would take it along.
// Creates keep 'S' and their order by priority: a sequence depends on its
// owner while the owner's default depends on it.
const preSequenceObjType = 's'

func formatSequenceTargetName(schemaName string, seq *Sequence) string {
	if seq != nil && seq.Persistence == 't' && !isTempSchemaName(schemaName) {
		return quoteIdentAlways(seq.Name)
	}
	return migrationQualifiedName(schemaName, seq.Name)
}

// sequenceOwnedByOp sets the column that owns seq, or NONE. A new owner is
// set after every table exists; NONE runs among the drops.
func sequenceOwnedByOp(schemaName string, seq *Sequence, owner string, phase MigrationPhase) MigrationOp {
	objType := byte('S')
	if phase == PhasePre {
		objType = preSequenceObjType
	}
	return MigrationOp{
		Type:          OpAlterSequence,
		SchemaName:    schemaName,
		ObjectName:    seq.Name,
		SQL:           fmt.Sprintf("ALTER SEQUENCE %s OWNED BY %s", formatSequenceTargetName(schemaName, seq), owner),
		Transactional: true,
		Phase:         phase,
		ObjType:       objType,
		ObjOID:        seq.OID,
		Priority:      PrioritySequence,
	}
}

// seqOpOrder returns a sort key: drops < creates < alters.
func seqOpOrder(t MigrationOpType) int {
	switch t {
	case OpDropSequence:
		return 0
	case OpCreateSequence:
		return 1
	case OpAlterSequence:
		return 2
	default:
		return 3
	}
}

// isIdentitySequence checks whether a sequence is the implicit backing sequence
// for an identity column. Identity columns create sequences named "{table}_{column}_seq"
// and the column's Default references nextval('{seqname}'::regclass).
func isIdentitySequence(c *Catalog, schemaName, seqName string) bool {
	if c == nil {
		return false
	}
	schema := c.schemaByName[schemaName]
	if schema == nil {
		return false
	}
	for _, rel := range schema.Relations {
		for _, col := range rel.Columns {
			if col.Identity != 0 {
				expectedSeqName := fmt.Sprintf("%s_%s_seq", rel.Name, col.Name)
				if expectedSeqName == seqName {
					return true
				}
			}
		}
	}
	return false
}

// formatCreateSequence builds a CREATE SEQUENCE statement with all options.
// The type is written unless it is the default, bigint.
func formatCreateSequence(c *Catalog, schemaName string, seq *Sequence) string {
	var b strings.Builder
	switch seq.Persistence {
	case 'u':
		b.WriteString("CREATE UNLOGGED SEQUENCE ")
	case 't':
		b.WriteString("CREATE TEMP SEQUENCE ")
	default:
		b.WriteString("CREATE SEQUENCE ")
	}
	b.WriteString(formatSequenceTargetName(schemaName, seq))
	if seq.TypeOID != 0 && seq.TypeOID != INT8OID {
		b.WriteString("\n    AS " + c.FormatType(seq.TypeOID, -1))
	}

	b.WriteString(fmt.Sprintf("\n    INCREMENT BY %d", seq.Increment))
	b.WriteString(fmt.Sprintf("\n    MINVALUE %d", seq.MinValue))
	b.WriteString(fmt.Sprintf("\n    MAXVALUE %d", seq.MaxValue))
	b.WriteString(fmt.Sprintf("\n    START WITH %d", seq.Start))
	b.WriteString(fmt.Sprintf("\n    CACHE %d", seq.CacheValue))
	if seq.Cycle {
		b.WriteString("\n    CYCLE")
	} else {
		b.WriteString("\n    NO CYCLE")
	}

	return b.String()
}

// formatAlterSequence builds an ALTER SEQUENCE statement for changed properties.
func formatAlterSequence(fromCat, toCat *Catalog, schemaName string, from, to *Sequence) string {
	var clauses []string

	if toType := toCat.FormatType(to.TypeOID, -1); fromCat.FormatType(from.TypeOID, -1) != toType {
		clauses = append(clauses, "AS "+toType)
	}

	if from.Increment != to.Increment {
		clauses = append(clauses, fmt.Sprintf("INCREMENT BY %d", to.Increment))
	}
	if from.MinValue != to.MinValue {
		clauses = append(clauses, fmt.Sprintf("MINVALUE %d", to.MinValue))
	}
	if from.MaxValue != to.MaxValue {
		clauses = append(clauses, fmt.Sprintf("MAXVALUE %d", to.MaxValue))
	}
	if from.Start != to.Start {
		clauses = append(clauses, fmt.Sprintf("START WITH %d", to.Start))
	}
	if from.CacheValue != to.CacheValue {
		clauses = append(clauses, fmt.Sprintf("CACHE %d", to.CacheValue))
	}
	if from.Cycle != to.Cycle {
		if to.Cycle {
			clauses = append(clauses, "CYCLE")
		} else {
			clauses = append(clauses, "NO CYCLE")
		}
	}

	if len(clauses) == 0 {
		return ""
	}

	return fmt.Sprintf("ALTER SEQUENCE %s\n    %s",
		formatSequenceTargetName(schemaName, to),
		strings.Join(clauses, "\n    "))
}

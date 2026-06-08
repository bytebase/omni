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
				SQL:           formatCreateSequence(entry.SchemaName, entry.To),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'S',
				ObjOID:        entry.To.OID,
				Priority:      PrioritySequence,
			})

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			// Skip dropping sequences that are owned by identity columns.
			// DROP IDENTITY automatically drops the identity sequence.
			if isIdentitySequence(from, entry.SchemaName, entry.Name) {
				continue
			}
			qn := migrationQualifiedName(entry.SchemaName, entry.Name)
			ops = append(ops, MigrationOp{
				Type:          OpDropSequence,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP SEQUENCE %s", qn),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'S',
				ObjOID:        entry.From.OID,
				Priority:      PrioritySequence,
			})

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			sql := formatAlterSequence(entry.SchemaName, entry.From, entry.To)
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
func formatCreateSequence(schemaName string, seq *Sequence) string {
	var b strings.Builder
	b.WriteString("CREATE SEQUENCE ")
	b.WriteString(migrationQualifiedName(schemaName, seq.Name))

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
func formatAlterSequence(schemaName string, from, to *Sequence) string {
	var clauses []string

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
		migrationQualifiedName(schemaName, to.Name),
		strings.Join(clauses, "\n    "))
}

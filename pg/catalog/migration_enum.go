package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateEnumDDL produces CREATE TYPE, DROP TYPE, and ALTER TYPE ADD VALUE
// operations for enum type changes.
func generateEnumDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Enums {
		switch entry.Action {
		case DiffAdd:
			qn := migrationQualifiedName(entry.SchemaName, entry.Name)
			quotedVals := make([]string, len(entry.ToValues))
			for i, v := range entry.ToValues {
				quotedVals[i] = singleQuote(v)
			}
			sql := fmt.Sprintf("CREATE TYPE %s AS ENUM (%s)", qn, strings.Join(quotedVals, ", "))
			ops = append(ops, MigrationOp{
				Type:          OpCreateType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           sql,
				Transactional: true,
			})

		case DiffDrop:
			qn := migrationQualifiedName(entry.SchemaName, entry.Name)
			ops = append(ops, MigrationOp{
				Type:          OpDropType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP TYPE %s", qn),
				Transactional: true,
			})

		case DiffModify:
			ops = append(ops, generateEnumModifyOps(entry)...)
		}
	}

	// Deterministic ordering: drops first, then creates, then alters.
	sort.Slice(ops, func(i, j int) bool {
		oi, oj := enumOpOrder(ops[i].Type), enumOpOrder(ops[j].Type)
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

// enumOpOrder returns a sort key for enum operation ordering.
func enumOpOrder(t MigrationOpType) int {
	switch t {
	case OpDropType:
		return 0
	case OpCreateType:
		return 1
	default:
		return 2
	}
}

// generateEnumModifyOps produces ALTER TYPE ADD VALUE ops for new enum values
// and warning ops for removed values.
func generateEnumModifyOps(entry EnumDiffEntry) []MigrationOp {
	var ops []MigrationOp
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)

	// Build set of old values for quick lookup.
	fromSet := make(map[string]bool, len(entry.FromValues))
	for _, v := range entry.FromValues {
		fromSet[v] = true
	}

	// Build set of new values for quick lookup.
	toSet := make(map[string]bool, len(entry.ToValues))
	for _, v := range entry.ToValues {
		toSet[v] = true
	}

	// Values removed from the enum → warning (PG cannot remove enum values).
	for _, v := range entry.FromValues {
		if !toSet[v] {
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("-- cannot remove enum value %s from %s", singleQuote(v), qn),
				Warning:       fmt.Sprintf("PostgreSQL does not support removing enum value %s from type %s", singleQuote(v), qn),
				Transactional: true,
			})
		}
	}

	// New values not in FromValues → ADD VALUE with position.
	// We need to figure out the position relative to existing values in ToValues.
	for i, v := range entry.ToValues {
		if fromSet[v] {
			continue // already exists
		}

		// Determine position clause.
		posClause := enumPositionClause(entry.ToValues, i, fromSet)
		sql := fmt.Sprintf("ALTER TYPE %s ADD VALUE %s%s", qn, singleQuote(v), posClause)

		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           sql,
			Transactional: false,
		})
	}

	return ops
}

// enumPositionClause determines the BEFORE/AFTER clause for an ADD VALUE statement.
// idx is the index of the new value in ToValues. fromSet contains the existing values.
func enumPositionClause(toValues []string, idx int, fromSet map[string]bool) string {
	// Find the previous existing value (scanning left from idx).
	prevExisting := ""
	for j := idx - 1; j >= 0; j-- {
		if fromSet[toValues[j]] {
			prevExisting = toValues[j]
			break
		}
	}

	// Find the next existing value (scanning right from idx).
	nextExisting := ""
	for j := idx + 1; j < len(toValues); j++ {
		if fromSet[toValues[j]] {
			nextExisting = toValues[j]
			break
		}
	}

	// If there's no previous existing value, use BEFORE the next existing value.
	if prevExisting == "" && nextExisting != "" {
		return " BEFORE " + singleQuote(nextExisting)
	}

	// If there's a previous existing value, use AFTER it.
	if prevExisting != "" {
		return " AFTER " + singleQuote(prevExisting)
	}

	// No existing neighbors (all values are new or only new values) — no position clause needed.
	return ""
}

// singleQuote wraps a string in single quotes, escaping embedded single quotes.
func singleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

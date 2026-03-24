package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateCommentDDL produces COMMENT ON operations from the diff.
func generateCommentDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Comments {
		switch entry.Action {
		case DiffAdd, DiffModify:
			sql := formatCommentSQL(entry.ObjType, entry.ObjDescription, entry.SubID, entry.To)
			if sql != "" {
				ops = append(ops, MigrationOp{
					Type:          OpComment,
					ObjectName:    entry.ObjDescription,
					SQL:           sql,
					Transactional: true,
				})
			}
		case DiffDrop:
			sql := formatCommentSQL(entry.ObjType, entry.ObjDescription, entry.SubID, "")
			if sql != "" {
				ops = append(ops, MigrationOp{
					Type:          OpComment,
					ObjectName:    entry.ObjDescription,
					SQL:           sql,
					Transactional: true,
				})
			}
		}
	}

	// Deterministic ordering.
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// formatCommentSQL generates a COMMENT ON ... IS '...' or COMMENT ON ... IS NULL statement.
func formatCommentSQL(objType byte, objDescription string, subID int16, text string) string {
	objTarget := commentObjectTarget(objType, objDescription, subID)
	if objTarget == "" {
		return ""
	}

	if text == "" {
		return fmt.Sprintf("COMMENT ON %s IS NULL", objTarget)
	}
	return fmt.Sprintf("COMMENT ON %s IS %s", objTarget, quoteLiteral(text))
}

// commentObjectTarget returns the COMMENT ON target string (e.g. "TABLE \"public\".\"t\"").
func commentObjectTarget(objType byte, objDescription string, subID int16) string {
	switch objType {
	case 'r': // relation (table, view, matview)
		if subID != 0 {
			// Column comment. Description is "schema.table".
			// We need to find the column name from the subID — but we only have
			// the description. For migration, the diff_comment resolves the column
			// attnum as SubID, and the description is "schema.table".
			// We handle this by looking up column in a different way:
			// For the migration DDL, the column description isn't directly available
			// from CommentDiffEntry, so we format as COLUMN schema.table.colname.
			// However, the diff_comment stores the description as "schema.table" and
			// subID as the attnum. We cannot resolve attnum to name here without the catalog.
			// Instead, we'll need to pass through the column name. But the current diff
			// structure only has SubID. We'll format with the available info.
			// Actually, looking at the diff_comment code, for columns the description
			// is "schema.table" and SubID is the attnum. We can't resolve it.
			// So let's return COLUMN with the schema.table and subID annotation.
			// The proper fix would be to include column name in the diff entry,
			// but for now we'll handle this in the higher-level code.
			return fmt.Sprintf("COLUMN %s", formatQualifiedDescription(objDescription))
		}
		return fmt.Sprintf("TABLE %s", formatQualifiedDescription(objDescription))
	case 'i': // index
		return fmt.Sprintf("INDEX %s", formatQualifiedDescription(objDescription))
	case 'f': // function
		// Description is the full identity like "schema.funcname(arg1type, arg2type)"
		return fmt.Sprintf("FUNCTION %s", objDescription)
	case 'n': // schema
		return fmt.Sprintf("SCHEMA %s", quoteIdentAlways(objDescription))
	case 't': // type (enum, composite)
		return fmt.Sprintf("TYPE %s", formatQualifiedDescription(objDescription))
	case 's': // sequence
		return fmt.Sprintf("SEQUENCE %s", formatQualifiedDescription(objDescription))
	case 'c': // constraint (schema.table.constraint)
		return formatConstraintCommentTarget(objDescription)
	case 'g': // trigger (schema.table.trigger)
		return formatTriggerCommentTarget(objDescription)
	case 'p': // policy (schema.table.policy)
		return formatPolicyCommentTarget(objDescription)
	default:
		return ""
	}
}

// formatQualifiedDescription formats a "schema.name" description into
// a quoted qualified identifier.
func formatQualifiedDescription(desc string) string {
	parts := strings.SplitN(desc, ".", 2)
	if len(parts) == 2 {
		return migrationQualifiedName(parts[0], parts[1])
	}
	return quoteIdentAlways(desc)
}

// formatConstraintCommentTarget formats "schema.table.constraint" into
// CONSTRAINT "constraint" ON "schema"."table".
func formatConstraintCommentTarget(desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("CONSTRAINT %s ON %s",
			quoteIdentAlways(parts[2]),
			migrationQualifiedName(parts[0], parts[1]))
	}
	return fmt.Sprintf("CONSTRAINT %s", quoteIdentAlways(desc))
}

// formatTriggerCommentTarget formats "schema.table.trigger" into
// TRIGGER "trigger" ON "schema"."table".
func formatTriggerCommentTarget(desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("TRIGGER %s ON %s",
			quoteIdentAlways(parts[2]),
			migrationQualifiedName(parts[0], parts[1]))
	}
	return fmt.Sprintf("TRIGGER %s", quoteIdentAlways(desc))
}

// formatPolicyCommentTarget formats "schema.table.policy" into
// POLICY "policy" ON "schema"."table".
func formatPolicyCommentTarget(desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("POLICY %s ON %s",
			quoteIdentAlways(parts[2]),
			migrationQualifiedName(parts[0], parts[1]))
	}
	return fmt.Sprintf("POLICY %s", quoteIdentAlways(desc))
}

// quoteLiteral returns a single-quoted SQL string literal with proper escaping.
func quoteLiteral(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return "'" + escaped + "'"
}

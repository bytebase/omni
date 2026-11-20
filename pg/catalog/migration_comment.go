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
			sql := formatCommentSQL(to, entry.ObjType, entry.ObjDescription, entry.SubID, entry.To)
			if sql != "" {
				ops = append(ops, MigrationOp{
					Type:          OpComment,
					ObjectName:    entry.ObjDescription,
					SQL:           sql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       entry.ObjType,
					Priority:      PriorityMetadata,
				})
			}
		case DiffDrop:
			sql := formatCommentSQL(from, entry.ObjType, entry.ObjDescription, entry.SubID, "")
			if sql != "" {
				ops = append(ops, MigrationOp{
					Type:          OpComment,
					ObjectName:    entry.ObjDescription,
					SQL:           sql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       entry.ObjType,
					Priority:      PriorityMetadata,
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
func formatCommentSQL(c *Catalog, objType byte, objDescription string, subID int16, text string) string {
	objTarget := commentObjectTarget(c, objType, objDescription, subID)
	if objTarget == "" {
		return ""
	}

	if text == "" {
		return fmt.Sprintf("COMMENT ON %s IS NULL", objTarget)
	}
	return fmt.Sprintf("COMMENT ON %s IS %s", objTarget, quoteLiteral(text))
}

// commentObjectTarget returns the COMMENT ON target string (e.g. "TABLE \"public\".\"t\"").
// The catalog c is used to resolve relation kind (view, matview) and function kind (procedure).
func commentObjectTarget(c *Catalog, objType byte, objDescription string, subID int16) string {
	switch objType {
	case 'r': // relation (table, view, matview)
		var rel *Relation
		parts := strings.SplitN(objDescription, ".", 2)
		if len(parts) == 2 && c != nil {
			rel = c.GetRelation(parts[0], parts[1])
		}
		targetRelName := formatQualifiedDescription(objDescription)
		if len(parts) == 2 {
			targetRelName = formatRelationName(parts[0], parts[1], rel)
		}
		if subID != 0 {
			// Resolve the column name from the catalog using the attnum (subID).
			colName := ""
			if rel != nil {
				for _, col := range rel.Columns {
					if col.AttNum == subID {
						colName = col.Name
						break
					}
				}
			}
			if colName != "" {
				return fmt.Sprintf("COLUMN %s.%s", targetRelName, quoteIdentAlways(colName))
			}
			return fmt.Sprintf("COLUMN %s", targetRelName)
		}
		// Determine relation kind from the catalog.
		relKind := resolveRelKindFromDescription(c, objDescription)
		switch relKind {
		case 'v':
			return fmt.Sprintf("VIEW %s", targetRelName)
		case 'm':
			return fmt.Sprintf("MATERIALIZED VIEW %s", targetRelName)
		default:
			return fmt.Sprintf("TABLE %s", targetRelName)
		}
	case 'i': // index
		parts := strings.SplitN(objDescription, ".", 2)
		if len(parts) == 2 && c != nil {
			if sc := c.GetSchema(parts[0]); sc != nil {
				if idx := sc.Indexes[parts[1]]; idx != nil {
					if rel := c.GetRelationByOID(idx.RelOID); rel != nil && rel.Persistence == 't' && !isTempSchemaName(parts[0]) {
						return fmt.Sprintf("INDEX %s", quoteIdentAlways(parts[1]))
					}
				}
			}
		}
		return fmt.Sprintf("INDEX %s", formatQualifiedDescription(objDescription))
	case 'f': // function/procedure
		// Determine if this is a procedure from the catalog.
		if resolveIsProcedure(c, objDescription) {
			return fmt.Sprintf("PROCEDURE %s", objDescription)
		}
		return fmt.Sprintf("FUNCTION %s", objDescription)
	case 'n': // schema
		return fmt.Sprintf("SCHEMA %s", quoteIdentAlways(objDescription))
	case 't': // type (enum, composite)
		return fmt.Sprintf("TYPE %s", formatQualifiedDescription(objDescription))
	case 's': // sequence
		parts := strings.SplitN(objDescription, ".", 2)
		if len(parts) == 2 && c != nil {
			if sc := c.GetSchema(parts[0]); sc != nil {
				if seq := sc.Sequences[parts[1]]; seq != nil && seq.Persistence == 't' && !isTempSchemaName(parts[0]) {
					return fmt.Sprintf("SEQUENCE %s", quoteIdentAlways(parts[1]))
				}
			}
		}
		return fmt.Sprintf("SEQUENCE %s", formatQualifiedDescription(objDescription))
	case 'c': // constraint (schema.table.constraint)
		return formatConstraintCommentTarget(c, objDescription)
	case 'd': // domain constraint (schema.domain.constraint)
		return formatDomainConstraintCommentTarget(objDescription)
	case 'g': // trigger (schema.table.trigger)
		return formatTriggerCommentTarget(c, objDescription)
	case 'E': // event trigger
		return fmt.Sprintf("EVENT TRIGGER %s", quoteIdentAlways(objDescription))
	case 'p': // policy (schema.table.policy)
		return formatPolicyCommentTarget(c, objDescription)
	case 'R': // rule (schema.table.rule)
		return formatRuleCommentTarget(c, objDescription)
	default:
		return ""
	}
}

func formatRuleCommentTarget(c *Catalog, desc string) string {
	sep := "."
	if strings.Contains(desc, "\x00") {
		sep = "\x00"
	}
	parts := strings.SplitN(desc, sep, 3)
	if len(parts) == 3 {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(parts[0], parts[1])
		}
		return fmt.Sprintf("RULE %s ON %s",
			quoteIdentAlways(parts[2]),
			formatRelationName(parts[0], parts[1], rel))
	}
	return fmt.Sprintf("RULE %s", quoteIdentAlways(desc))
}

// resolveRelKindFromDescription looks up a relation by "schema.name" description
// and returns its RelKind. Returns 0 if not found.
func resolveRelKindFromDescription(c *Catalog, desc string) byte {
	if c == nil {
		return 0
	}
	parts := strings.SplitN(desc, ".", 2)
	if len(parts) != 2 {
		return 0
	}
	rel := c.GetRelation(parts[0], parts[1])
	if rel == nil {
		return 0
	}
	return rel.RelKind
}

// resolveIsProcedure checks if a function identity (e.g. "schema.name(argtypes)")
// refers to a procedure (Kind='p') in the catalog.
func resolveIsProcedure(c *Catalog, identity string) bool {
	if c == nil {
		return false
	}
	for _, up := range c.userProcs {
		if up != nil && up.Kind == 'p' && funcIdentity(c, up) == identity {
			return true
		}
	}
	return false
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
func formatConstraintCommentTarget(c *Catalog, desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(parts[0], parts[1])
		}
		return fmt.Sprintf("CONSTRAINT %s ON %s",
			quoteIdentAlways(parts[2]),
			formatRelationName(parts[0], parts[1], rel))
	}
	return fmt.Sprintf("CONSTRAINT %s", quoteIdentAlways(desc))
}

func formatDomainConstraintCommentTarget(desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("CONSTRAINT %s ON DOMAIN %s",
			quoteIdentAlways(parts[2]),
			migrationQualifiedName(parts[0], parts[1]))
	}
	return fmt.Sprintf("CONSTRAINT %s", quoteIdentAlways(desc))
}

// formatTriggerCommentTarget formats "schema.table.trigger" into
// TRIGGER "trigger" ON "schema"."table".
func formatTriggerCommentTarget(c *Catalog, desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(parts[0], parts[1])
		}
		return fmt.Sprintf("TRIGGER %s ON %s",
			quoteIdentAlways(parts[2]),
			formatRelationName(parts[0], parts[1], rel))
	}
	return fmt.Sprintf("TRIGGER %s", quoteIdentAlways(desc))
}

// formatPolicyCommentTarget formats "schema.table.policy" into
// POLICY "policy" ON "schema"."table".
func formatPolicyCommentTarget(c *Catalog, desc string) string {
	parts := strings.SplitN(desc, ".", 3)
	if len(parts) == 3 {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(parts[0], parts[1])
		}
		return fmt.Sprintf("POLICY %s ON %s",
			quoteIdentAlways(parts[2]),
			formatRelationName(parts[0], parts[1], rel))
	}
	return fmt.Sprintf("POLICY %s", quoteIdentAlways(desc))
}

// quoteLiteral returns a single-quoted SQL string literal with proper escaping.
func quoteLiteral(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return "'" + escaped + "'"
}

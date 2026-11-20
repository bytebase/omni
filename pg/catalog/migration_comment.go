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
			sql := formatCommentSQL(to, entry, entry.To)
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
			sql := formatCommentSQL(from, entry, "")
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
func formatCommentSQL(c *Catalog, entry CommentDiffEntry, text string) string {
	objTarget := commentObjectTarget(c, entry)
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
func commentObjectTarget(c *Catalog, entry CommentDiffEntry) string {
	qualifiedTarget := func() string {
		if schema, name, ok := commentParts2(entry); ok {
			return migrationQualifiedName(schema, name)
		}
		return formatQualifiedDescription(entry.ObjDescription)
	}

	switch entry.ObjType {
	case 'r': // relation (table, view, matview)
		if entry.SubID != 0 {
			colName := ""
			if schema, relName, ok := commentParts2(entry); ok && c != nil {
				if rel := c.GetRelation(schema, relName); rel != nil {
					for _, col := range rel.Columns {
						if col.AttNum == entry.SubID {
							colName = col.Name
							break
						}
					}
				}
			}
			if colName != "" {
				return fmt.Sprintf("COLUMN %s.%s", qualifiedTarget(), quoteIdentAlways(colName))
			}
			return fmt.Sprintf("COLUMN %s", qualifiedTarget())
		}
		switch resolveRelKindFromEntry(c, entry) {
		case 'v':
			return fmt.Sprintf("VIEW %s", qualifiedTarget())
		case 'm':
			return fmt.Sprintf("MATERIALIZED VIEW %s", qualifiedTarget())
		default:
			return fmt.Sprintf("TABLE %s", qualifiedTarget())
		}
	case 'i': // index
		return fmt.Sprintf("INDEX %s", qualifiedTarget())
	case 'f': // function/procedure
		target := entry.ObjDescription
		if schema, obj, ok := commentParts2(entry); ok {
			if name, args, hasParen := strings.Cut(obj, "("); hasParen && strings.ContainsAny(schema+name, `. "`) {
				target = migrationQualifiedName(schema, name) + "(" + args
			}
		}
		if resolveIsProcedure(c, entry.ObjDescription) {
			return fmt.Sprintf("PROCEDURE %s", target)
		}
		return fmt.Sprintf("FUNCTION %s", target)
	case 'n': // schema
		name := entry.ObjDescription
		if len(entry.objParts) > 0 {
			name = entry.objParts[0]
		}
		return fmt.Sprintf("SCHEMA %s", quoteIdentAlways(name))
	case 't': // type (enum, composite)
		return fmt.Sprintf("TYPE %s", qualifiedTarget())
	case 's': // sequence
		return fmt.Sprintf("SEQUENCE %s", qualifiedTarget())
	case 'c': // constraint (schema.table.constraint)
		return formatSubObjectCommentTarget("CONSTRAINT", entry)
	case 'g': // trigger (schema.table.trigger)
		return formatSubObjectCommentTarget("TRIGGER", entry)
	case 'E': // event trigger
		name := entry.ObjDescription
		if len(entry.objParts) > 0 {
			name = entry.objParts[0]
		}
		return fmt.Sprintf("EVENT TRIGGER %s", quoteIdentAlways(name))
	case 'p': // policy (schema.table.policy)
		return formatSubObjectCommentTarget("POLICY", entry)
	default:
		return ""
	}
}

// resolveRelKindFromEntry looks up a relation from entry's structured parts or
// flattened description and returns its RelKind. Returns 0 if not found.
func resolveRelKindFromEntry(c *Catalog, entry CommentDiffEntry) byte {
	if c == nil {
		return 0
	}
	schema, relName, ok := commentParts2(entry)
	if !ok {
		return 0
	}
	rel := c.GetRelation(schema, relName)
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

func formatSubObjectCommentTarget(keyword string, entry CommentDiffEntry) string {
	if schema, table, sub, ok := commentParts3(entry); ok {
		return fmt.Sprintf("%s %s ON %s", keyword, quoteIdentAlways(sub), migrationQualifiedName(schema, table))
	}
	return fmt.Sprintf("%s %s", keyword, quoteIdentAlways(entry.ObjDescription))
}

// quoteLiteral returns a single-quoted SQL string literal with proper escaping.
func quoteLiteral(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return "'" + escaped + "'"
}

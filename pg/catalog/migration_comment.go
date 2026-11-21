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
			sql := formatCommentEntrySQL(to, entry, entry.To)
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
			if commentDroppedWithForeignTable(from, diff, entry) {
				continue
			}
			sql := formatCommentEntrySQL(from, entry, "")
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

func commentDroppedWithForeignTable(from *Catalog, diff *SchemaDiff, entry CommentDiffEntry) bool {
	if diff == nil || (len(diff.ForeignTables) == 0 && len(diff.ForeignServers) == 0) {
		return false
	}
	schemaName, relName, ok := commentParts2(entry)
	if !ok {
		return false
	}
	rebuiltServers := rebuiltForeignServerNames(diff)
	relKey := schemaName + "\x00" + relName
	rebuiltOrDroppedFT := make(map[string]bool)
	droppedFTCols := make(map[string]map[int16]bool)
	droppedFTCons := make(map[string]map[string]bool)
	droppedFTTrigs := make(map[string]map[string]bool)

	for _, ft := range diff.ForeignTables {
		k := ft.SchemaName + "\x00" + ft.Name
		switch ft.Action {
		case DiffDrop:
			rebuiltOrDroppedFT[k] = true
		case DiffModify:
			if ft.From != nil && ft.To != nil &&
				(ft.From.ForeignServer != ft.To.ForeignServer || rebuiltServers[ft.From.ForeignServer]) {
				rebuiltOrDroppedFT[k] = true
				continue
			}
			for _, col := range ft.Columns {
				if col.Action == DiffDrop && col.From != nil {
					if droppedFTCols[k] == nil {
						droppedFTCols[k] = make(map[int16]bool)
					}
					droppedFTCols[k][col.From.AttNum] = true
				}
			}
			for _, ce := range ft.Constraints {
				if (ce.Action == DiffDrop || ce.Action == DiffModify) && ce.From != nil {
					if droppedFTCons[k] == nil {
						droppedFTCons[k] = make(map[string]bool)
					}
					droppedFTCons[k][ce.From.Name] = true
				}
			}
			for _, te := range ft.Triggers {
				if te.Action == DiffDrop && te.From != nil {
					if droppedFTTrigs[k] == nil {
						droppedFTTrigs[k] = make(map[string]bool)
					}
					droppedFTTrigs[k][te.From.Name] = true
				}
			}
		}
	}
	if from != nil && len(rebuiltServers) > 0 {
		if rel := from.GetRelation(schemaName, relName); rel != nil && rel.RelKind == 'f' && rebuiltServers[rel.ForeignServer] {
			rebuiltOrDroppedFT[relKey] = true
		}
	}

	switch entry.ObjType {
	case 'r':
		if rebuiltOrDroppedFT[relKey] {
			return true
		}
		if entry.SubID != 0 {
			return droppedFTCols[relKey][entry.SubID]
		}
	case 'c':
		if rebuiltOrDroppedFT[relKey] {
			return true
		}
		if len(entry.objParts) >= 3 {
			return droppedFTCons[relKey][entry.objParts[2]]
		}
	case 'g':
		if rebuiltOrDroppedFT[relKey] {
			return true
		}
		if len(entry.objParts) >= 3 {
			return droppedFTTrigs[relKey][entry.objParts[2]]
		}
	case 'R', 'p':
		return rebuiltOrDroppedFT[relKey]
	}
	return false
}

// formatCommentSQL generates a COMMENT ON ... IS '...' or COMMENT ON ... IS NULL statement.
func formatCommentSQL(c *Catalog, objType byte, objDescription string, subID int16, text string) string {
	return formatCommentEntrySQL(c, CommentDiffEntry{
		ObjType:        objType,
		ObjDescription: objDescription,
		SubID:          subID,
	}, text)
}

func formatCommentEntrySQL(c *Catalog, entry CommentDiffEntry, text string) string {
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
		var rel *Relation
		targetRelName := qualifiedTarget()
		if schema, relName, ok := commentParts2(entry); ok {
			if c != nil {
				rel = c.GetRelation(schema, relName)
			}
			targetRelName = formatRelationName(schema, relName, rel)
		}
		if entry.SubID != 0 {
			colName := ""
			if rel != nil {
				for _, col := range rel.Columns {
					if col.AttNum == entry.SubID {
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
		var relKind byte
		if rel != nil {
			relKind = rel.RelKind
		}
		switch relKind {
		case 'v':
			return fmt.Sprintf("VIEW %s", targetRelName)
		case 'm':
			return fmt.Sprintf("MATERIALIZED VIEW %s", targetRelName)
		case 'f':
			return fmt.Sprintf("FOREIGN TABLE %s", targetRelName)
		default:
			return fmt.Sprintf("TABLE %s", targetRelName)
		}
	case 'i': // index
		if schema, name, ok := commentParts2(entry); ok && c != nil {
			if sc := c.GetSchema(schema); sc != nil {
				if idx := sc.Indexes[name]; idx != nil {
					if rel := c.GetRelationByOID(idx.RelOID); rel != nil && rel.Persistence == 't' && !isTempSchemaName(schema) {
						return fmt.Sprintf("INDEX %s", quoteIdentAlways(name))
					}
				}
			}
		}
		return fmt.Sprintf("INDEX %s", qualifiedTarget())
	case 'f': // function/procedure
		target := entry.ObjDescription
		if schema, obj, ok := commentParts2(entry); ok {
			if name, args, hasParen := splitRoutineIdentityNameAndArgs(obj); hasParen {
				target = quoteIdentifier(schema) + "." + quoteIdentifier(name) + "(" + args
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
		if schema, name, ok := commentParts2(entry); ok && c != nil {
			if sc := c.GetSchema(schema); sc != nil {
				if seq := sc.Sequences[name]; seq != nil && seq.Persistence == 't' && !isTempSchemaName(schema) {
					return fmt.Sprintf("SEQUENCE %s", quoteIdentAlways(name))
				}
			}
		}
		return fmt.Sprintf("SEQUENCE %s", qualifiedTarget())
	case 'c': // constraint (schema.table.constraint)
		return formatSubObjectCommentTarget(c, "CONSTRAINT", entry)
	case 'd': // domain constraint (schema.domain.constraint)
		if schema, domain, sub, ok := commentParts3(entry); ok {
			return fmt.Sprintf("CONSTRAINT %s ON DOMAIN %s",
				quoteIdentAlways(sub),
				migrationQualifiedName(schema, domain))
		}
		return fmt.Sprintf("CONSTRAINT %s", quoteIdentAlways(entry.ObjDescription))
	case 'g': // trigger (schema.table.trigger)
		return formatSubObjectCommentTarget(c, "TRIGGER", entry)
	case 'E': // event trigger
		name := entry.ObjDescription
		if len(entry.objParts) > 0 {
			name = entry.objParts[0]
		}
		return fmt.Sprintf("EVENT TRIGGER %s", quoteIdentAlways(name))
	case 'p': // policy (schema.table.policy)
		return formatSubObjectCommentTarget(c, "POLICY", entry)
	case 'R': // rule (schema.table.rule)
		return formatRuleCommentTarget(c, entry)
	default:
		return ""
	}
}

func formatRuleCommentTarget(c *Catalog, entry CommentDiffEntry) string {
	if len(entry.objParts) >= 3 {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(entry.objParts[0], entry.objParts[1])
		}
		return fmt.Sprintf("RULE %s ON %s",
			quoteIdentAlways(entry.objParts[2]),
			formatRelationName(entry.objParts[0], entry.objParts[1], rel))
	}
	sep := "."
	if strings.Contains(entry.ObjDescription, "\x00") {
		sep = "\x00"
	}
	parts := strings.SplitN(entry.ObjDescription, sep, 3)
	if len(parts) == 3 {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(parts[0], parts[1])
		}
		return fmt.Sprintf("RULE %s ON %s",
			quoteIdentAlways(parts[2]),
			formatRelationName(parts[0], parts[1], rel))
	}
	return fmt.Sprintf("RULE %s", quoteIdentAlways(entry.ObjDescription))
}

// splitRoutineIdentityNameAndArgs splits "funcname(arg1, numeric(10, 2))" at
// the opening paren of its parameter list even when funcname itself contains '('.
func splitRoutineIdentityNameAndArgs(obj string) (name, args string, ok bool) {
	if !strings.HasSuffix(obj, ")") {
		return strings.Cut(obj, "(")
	}
	depth := 0
	for i := len(obj) - 1; i >= 0; i-- {
		switch obj[i] {
		case ')':
			depth++
		case '(':
			depth--
			if depth == 0 {
				return obj[:i], obj[i+1:], true
			}
		}
	}
	return "", "", false
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

func formatSubObjectCommentTarget(c *Catalog, keyword string, entry CommentDiffEntry) string {
	if schema, table, sub, ok := commentParts3(entry); ok {
		var rel *Relation
		if c != nil {
			rel = c.GetRelation(schema, table)
		}
		return fmt.Sprintf("%s %s ON %s", keyword, quoteIdentAlways(sub), formatRelationName(schema, table, rel))
	}
	return fmt.Sprintf("%s %s", keyword, quoteIdentAlways(entry.ObjDescription))
}

// quoteLiteral returns a single-quoted SQL string literal with proper escaping.
func quoteLiteral(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return "'" + escaped + "'"
}

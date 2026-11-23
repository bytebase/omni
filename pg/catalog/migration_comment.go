package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateCommentDDL produces COMMENT ON operations from the diff.
func generateCommentDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	emittedRenames := make(map[string]bool)

	for _, entry := range diff.Comments {
		switch entry.Action {
		case DiffAdd, DiffModify:
			if renameOp, ok := renameNotNullForComment(from, to, entry); ok && !emittedRenames[renameOp.SQL] {
				emittedRenames[renameOp.SQL] = true
				ops = append(ops, renameOp)
			}
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
			if commentDroppedWithForeignTable(from, diff, entry) || commentDroppedWithNotNullConstraint(from, to, entry) {
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

// commentDroppedWithNotNullConstraint reports whether entry removes a comment
// on a PG 18 cataloged NOT NULL constraint that no longer exists in to (e.g.
// because the table/column was dropped or ALTER COLUMN ... DROP NOT NULL
// implicitly removed the constraint).
func commentDroppedWithNotNullConstraint(from, to *Catalog, entry CommentDiffEntry) bool {
	if from == nil || entry.ObjType != 'c' || len(entry.objParts) < 3 {
		return false
	}
	schemaName, relName, conName := entry.objParts[0], entry.objParts[1], entry.objParts[2]
	fromRel := from.GetRelation(schemaName, relName)
	if fromRel == nil {
		return false
	}
	var isNN bool
	for _, con := range from.consByRel[fromRel.OID] {
		if con.Name == conName && con.Type == ConstraintNotNull {
			isNN = true
			break
		}
	}
	if !isNN {
		return false
	}
	if to == nil {
		return true
	}
	toRel := to.GetRelation(schemaName, relName)
	if toRel == nil {
		return true
	}
	for _, con := range to.consByRel[toRel.OID] {
		if con.Name == conName && con.Type == ConstraintNotNull {
			return false
		}
	}
	return true
}

// renameNotNullForComment returns an ALTER TABLE ... RENAME CONSTRAINT op when
// a comment targets a PG 18 cataloged NOT NULL constraint whose column was
// already NOT NULL in from under a different constraint name.
func renameNotNullForComment(from, to *Catalog, entry CommentDiffEntry) (MigrationOp, bool) {
	if from == nil || to == nil || from.Version() < PG18 || to.Version() < PG18 || entry.ObjType != 'c' || len(entry.objParts) < 3 {
		return MigrationOp{}, false
	}
	schemaName, relName, conName := entry.objParts[0], entry.objParts[1], entry.objParts[2]
	toRel := to.GetRelation(schemaName, relName)
	fromRel := from.GetRelation(schemaName, relName)
	if toRel == nil || fromRel == nil || (toRel.RelKind != 'r' && toRel.RelKind != 'p') {
		return MigrationOp{}, false
	}
	var toColName string
	for _, con := range to.consByRel[toRel.OID] {
		if con.Name == conName && con.Type == ConstraintNotNull && len(con.Columns) == 1 {
			if ci := indexOfAttNum(toRel, con.Columns[0]); ci >= 0 {
				toColName = toRel.Columns[ci].Name
			}
			break
		}
	}
	if toColName == "" {
		return MigrationOp{}, false
	}
	fi, ok := fromRel.colByName[toColName]
	if !ok || !fromRel.Columns[fi].NotNull {
		return MigrationOp{}, false
	}
	var fromCon *Constraint
	for _, con := range from.consByRel[fromRel.OID] {
		if con.Type == ConstraintNotNull && len(con.Columns) == 1 && con.Columns[0] == fromRel.Columns[fi].AttNum {
			fromCon = con
			break
		}
	}
	if fromCon == nil || fromCon.Name == "" || fromCon.Name == conName {
		return MigrationOp{}, false
	}
	fromNN := fromCon.Name

	// If the NOT NULL constraint on fromRel is inherited (ConInhCount > 0),
	// PostgreSQL rejects renaming it directly on the child/partition; resolve
	// the owning ancestor so recursive RENAME CONSTRAINT updates the hierarchy.
	targetRel := fromRel
	if fromCon.ConInhCount > 0 {
		visited := make(map[uint32]bool)
		var findOwner func(r *Relation) *Relation
		findOwner = func(r *Relation) *Relation {
			if r == nil || visited[r.OID] {
				return nil
			}
			visited[r.OID] = true
			var parents []uint32
			if r.PartitionOf != 0 {
				parents = append(parents, r.PartitionOf)
			}
			parents = append(parents, r.InhParents...)
			for _, poid := range parents {
				p := from.relationByOID[poid]
				if p == nil {
					continue
				}
				pi, ok := p.colByName[toColName]
				if !ok {
					continue
				}
				for _, pcon := range from.consByRel[p.OID] {
					if pcon.Type == ConstraintNotNull && len(pcon.Columns) == 1 && pcon.Columns[0] == p.Columns[pi].AttNum && pcon.Name == fromNN {
						if pcon.ConInhCount == 0 {
							return p
						}
						if owner := findOwner(p); owner != nil {
							return owner
						}
					}
				}
			}
			return nil
		}
		if owner := findOwner(fromRel); owner != nil {
			targetRel = owner
		}
	}
	targetSchema := schemaName
	if targetRel.Schema != nil {
		targetSchema = targetRel.Schema.Name
	}
	targetObjOID := toRel.OID
	if toOwner := to.GetRelation(targetSchema, targetRel.Name); toOwner != nil {
		targetObjOID = toOwner.OID
	}
	tableName := fmt.Sprintf("%s.%s", quoteIdentifier(targetSchema), quoteIdentifier(targetRel.Name))
	if targetRel.Persistence == 't' || isTempSchemaName(targetSchema) {
		tableName = quoteIdentifier(targetRel.Name)
	}
	return MigrationOp{
		Type:          OpAlterColumn,
		SchemaName:    targetSchema,
		ObjectName:    targetRel.Name,
		SQL:           fmt.Sprintf("ALTER TABLE %s RENAME CONSTRAINT %s TO %s", tableName, quoteIdentifier(fromNN), quoteIdentifier(conName)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        targetObjOID,
		Priority:      PriorityColumn,
	}, true
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

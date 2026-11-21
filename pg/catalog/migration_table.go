package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// quoteIdentAlways unconditionally double-quotes an identifier.
// Migration DDL always quotes identifiers to avoid keyword collisions.
func quoteIdentAlways(name string) string {
	return "\"" + strings.ReplaceAll(name, "\"", "\"\"") + "\""
}

// migrationQualifiedName returns a schema-qualified, always-quoted identifier for migration DDL.
func migrationQualifiedName(schema, name string) string {
	return quoteIdentAlways(schema) + "." + quoteIdentAlways(name)
}

func formatRelationTargetName(schemaName string, rel *Relation) string {
	if rel != nil {
		return formatRelationName(schemaName, rel.Name, rel)
	}
	return formatRelationName(schemaName, "", nil)
}

func formatRelationName(schemaName, relName string, rel *Relation) string {
	if (rel != nil && rel.Persistence == 't') || isTempSchemaName(schemaName) {
		return quoteIdentAlways(relName)
	}
	return migrationQualifiedName(schemaName, relName)
}

func buildTableOwnerOp(schemaName, tableName string, rel *Relation) (MigrationOp, bool) {
	if rel == nil || rel.OwnerName == "" || (!rel.OwnerExplicit && rel.OwnerName == "postgres") {
		return MigrationOp{}, false
	}
	qn := formatRelationTargetName(schemaName, rel)
	return MigrationOp{
		Type:          OpAlterTable,
		SchemaName:    schemaName,
		ObjectName:    tableName,
		SQL:           fmt.Sprintf("ALTER TABLE %s OWNER TO %s", qn, quoteIdentifier(rel.OwnerName)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        rel.OID,
		Priority:      PriorityMetadata,
	}, true
}

// generateTableDDL produces CREATE TABLE and DROP TABLE operations from the diff.
func generateTableDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Relations {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			rk := entry.To.RelKind
			if rk != 'r' && rk != 'p' {
				continue
			}
			// Skip partition children — handled by generatePartitionDDL.
			if entry.To.PartitionBound != nil && entry.To.PartitionOf != 0 {
				continue
			}
			ddl := formatCreateTable(from, to, entry.SchemaName, entry.From, entry.To, false)
			ops = append(ops, MigrationOp{
				Type:          OpCreateTable,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           ddl,
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'r',
				ObjOID:        entry.To.OID,
				Priority:      PriorityTable,
			})
			if ownerOp, ok := buildTableOwnerOp(entry.SchemaName, entry.Name, entry.To); ok {
				ops = append(ops, ownerOp)
			}

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			rk := entry.To.RelKind
			if rk != 'r' && rk != 'p' {
				continue
			}
			if effectiveTableOwner(entry.From) != effectiveTableOwner(entry.To) || tableRebuilt(from, to, entry.From, entry.To) {
				if ownerOp, ok := buildTableOwnerOp(entry.SchemaName, entry.Name, entry.To); ok {
					ops = append(ops, ownerOp)
				}
			}

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			rk := entry.From.RelKind
			if rk != 'r' && rk != 'p' {
				continue
			}
			qn := formatRelationTargetName(entry.SchemaName, entry.From)
			ops = append(ops, MigrationOp{
				Type:          OpDropTable,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP TABLE %s CASCADE", qn),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'r',
				ObjOID:        entry.From.OID,
				Priority:      PriorityTable,
			})
		}
	}

	// Deterministic ordering: drops first (sorted), then creates (sorted).
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			// Drops before creates.
			if ops[i].Type == OpDropTable {
				return true
			}
			if ops[j].Type == OpDropTable {
				return false
			}
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// FormatCreateTable generates a full CREATE TABLE DDL statement for a relation.
// Exported so that section 1.7 (partitions) can reuse it.
func FormatCreateTable(c *Catalog, schemaName string, rel *Relation) string {
	return formatCreateTable(nil, c, schemaName, nil, rel, false)
}

func formatCreateTable(from, c *Catalog, schemaName string, fromRel, rel *Relation, forceExplicitHeap bool) string {
	var b strings.Builder

	// Persistence prefix.
	persistence := rel.Persistence
	if persistence == 0 {
		persistence = 'p'
	}
	switch persistence {
	case 'u':
		b.WriteString("CREATE UNLOGGED TABLE ")
	case 't':
		b.WriteString("CREATE TEMP TABLE ")
	default:
		b.WriteString("CREATE TABLE ")
	}

	if persistence == 't' && !isTempSchemaName(schemaName) {
		b.WriteString(quoteIdentAlways(rel.Name))
	} else {
		b.WriteString(migrationQualifiedName(schemaName, rel.Name))
	}
	b.WriteString(" (\n")

	// Column definitions.
	for i, col := range rel.Columns {
		if i > 0 {
			b.WriteString(",\n")
		}
		b.WriteString("    ")
		b.WriteString(formatColumnDef(c, rel, col))
	}

	// Inline constraints: PK, UNIQUE, CHECK (not FK).
	constraints := c.ConstraintsOf(rel.OID)
	var inlineConstraints []string
	for _, con := range constraints {
		switch con.Type {
		case ConstraintPK:
			inlineConstraints = append(inlineConstraints, formatPKConstraint(rel, con))
		case ConstraintUnique:
			inlineConstraints = append(inlineConstraints, formatUniqueConstraint(c, rel, con))
		case ConstraintCheck:
			inlineConstraints = append(inlineConstraints, formatCheckConstraint(con))
		}
	}

	// Sort inline constraints for determinism.
	sort.Strings(inlineConstraints)

	for _, cs := range inlineConstraints {
		b.WriteString(",\n    ")
		b.WriteString(cs)
	}

	b.WriteString("\n)")

	// Append PARTITION BY clause for partitioned tables.
	b.WriteString(formatPartitionByClause(rel))

	// Append INHERITS clause for table inheritance.
	b.WriteString(formatInheritsClause(c, rel))

	// Append USING clause for explicit or non-default table access method.
	b.WriteString(formatRelationUsingClause(from, c, fromRel, rel, forceExplicitHeap))

	return b.String()
}

func formatRelationUsingClause(from, to *Catalog, fromRel, toRel *Relation, forceExplicitHeap bool) string {
	if toRel == nil {
		return ""
	}
	if toRel.RelKind == 'p' {
		if toRel.AccessMethod != "" {
			return " USING " + quoteIdentifier(toRel.AccessMethod)
		}
		return ""
	}
	am := normalizeRelationAccessMethod(toRel)
	if am == "" {
		return ""
	}
	if am != "heap" ||
		toRel.AccessMethod != "" ||
		(fromRel != nil && fromRel.AccessMethod != "") ||
		catalogHasNonDefaultTableAM(from) ||
		catalogHasNonDefaultTableAM(to) ||
		forceExplicitHeap {
		return " USING " + quoteIdentifier(am)
	}
	return ""
}

func catalogHasNonDefaultTableAM(c *Catalog) bool {
	if c == nil {
		return false
	}
	if c.defaultTableAccessMethod != "" {
		return true
	}
	for _, am := range c.accessMethods {
		if am != nil && am.Type == 't' && am.OID >= FirstNormalObjectId {
			return true
		}
	}
	return false
}

// formatColumnDef formats a single column definition for CREATE TABLE.
func formatColumnDef(c *Catalog, rel *Relation, col *Column) string {
	var parts []string
	parts = append(parts, quoteIdentAlways(col.Name))

	// Detect serial pattern: integer/bigint/smallint with nextval default owned by this column.
	// Output "serial"/"bigserial"/"smallserial" instead of "integer DEFAULT nextval(...)".
	if serialType := detectSerialType(c, rel, col); serialType != "" {
		parts = append(parts, serialType)
		if col.NotNull {
			parts = append(parts, "NOT NULL")
		}
		return strings.Join(parts, " ")
	}

	parts = append(parts, c.FormatType(col.TypeOID, col.TypeMod)+formatColumnCollate(c, col))

	if stor := columnExplicitStorage(c, col); stor != "" {
		parts = append(parts, "STORAGE "+stor)
	}
	if comp := columnExplicitCompression(col); comp != "" {
		parts = append(parts, "COMPRESSION "+comp)
	}

	if col.NotNull {
		parts = append(parts, "NOT NULL")
	}

	if col.HasDefault && col.Default != "" && col.Identity == 0 && col.Generated == 0 {
		parts = append(parts, "DEFAULT "+col.Default)
	}

	if col.Identity == 'a' {
		parts = append(parts, "GENERATED ALWAYS AS IDENTITY")
	} else if col.Identity == 'd' {
		parts = append(parts, "GENERATED BY DEFAULT AS IDENTITY")
	}

	if col.Generated == 's' {
		expr := col.GenerationExpr
		if expr == "" {
			expr = col.Default
		}
		if expr != "" {
			parts = append(parts, "GENERATED ALWAYS AS ("+expr+") STORED")
		}
	}

	return strings.Join(parts, " ")
}

// detectSerialType checks if a column should be output as serial/bigserial/smallserial.
// A serial column has a nextval default on the sequence it owns, named as
// serial would name it (sequenceManagedByColumn).
func detectSerialType(c *Catalog, rel *Relation, col *Column) string {
	if col.Identity != 0 || col.Generated != 0 {
		return ""
	}
	if !col.HasDefault || col.Default == "" {
		return ""
	}
	if !strings.Contains(col.Default, "nextval(") {
		return ""
	}

	// Check the type: serial = int4, bigserial = int8, smallserial = int2
	typeName := c.FormatType(col.TypeOID, col.TypeMod)
	var serialType string
	switch typeName {
	case "integer":
		serialType = "serial"
	case "bigint":
		serialType = "bigserial"
	case "smallint":
		serialType = "smallserial"
	default:
		return ""
	}

	// Extract sequence name from nextval('seqname'::regclass) and check if it's owned
	seqName := extractSeqNameFromDefault(col.Default)
	if seqName == "" {
		return ""
	}

	if rel.Schema == nil {
		return ""
	}
	seq := rel.Schema.Sequences[seqName]
	if seq == nil || seq.OwnerRelOID != rel.OID || seq.OwnerAttNum != col.AttNum || !sequenceManagedByColumn(c, seq) {
		return ""
	}
	return serialType
}

// extractSeqNameFromDefault extracts a sequence name from "nextval('seqname'::regclass)".
func extractSeqNameFromDefault(def string) string {
	// Pattern: nextval('seqname'::regclass)
	start := strings.Index(def, "nextval('")
	if start < 0 {
		return ""
	}
	start += len("nextval('")
	end := strings.Index(def[start:], "'")
	if end < 0 {
		return ""
	}
	name := def[start : start+end]
	// Handle schema-qualified: public.seqname -> seqname
	if dot := strings.LastIndex(name, "."); dot >= 0 {
		name = name[dot+1:]
	}
	return name
}

// formatCollate returns the COLLATE clause for a column or attribute with an
// explicit collation, or "" for the type's default.
func formatCollate(collation string) string {
	if collation == "" {
		return ""
	}
	return " COLLATE " + quoteIdentAlways(collation)
}

func columnCollationResolvedNamespace(c *Catalog, col *Column) string {
	if col == nil || col.CollationName == "" {
		return ""
	}
	if c != nil && col.UserCollationOID != 0 {
		for _, s := range c.schemas {
			if s != nil {
				if uc := s.Collations[col.CollationName]; uc != nil && uc.OID == col.UserCollationOID {
					return s.Name
				}
			}
		}
	}
	if col.CollationSchema != "" {
		return col.CollationSchema
	}
	if isBuiltinCollationName(col.CollationName) {
		return "pg_catalog"
	}
	return "public"
}

func columnCollationSchema(c *Catalog, col *Column) string {
	if col == nil || col.CollationName == "" {
		return ""
	}
	ns := columnCollationResolvedNamespace(c, col)
	switch ns {
	case "", "pg_catalog":
		if col.CollationSchema == "pg_catalog" {
			return "pg_catalog"
		}
		if c != nil {
			for _, s := range c.schemas {
				if s != nil && s.OID != PGCatalogNamespace && s.Collations[col.CollationName] != nil {
					return "pg_catalog"
				}
			}
		}
		return ""
	case "public":
		if col.UserCollationOID != 0 && (isBuiltinCollationName(col.CollationName) || col.CollationSchema == "public") {
			return "public"
		}
		if col.CollationSchema == "public" {
			return "public"
		}
		return ""
	default:
		return ns
	}
}

func columnCollationKey(c *Catalog, col *Column) string {
	if col == nil || col.CollationName == "" {
		return ""
	}
	if ns := columnCollationResolvedNamespace(c, col); ns != "" {
		return migrationQualifiedName(ns, col.CollationName)
	}
	return quoteIdentAlways(col.CollationName)
}

func formatColumnCollate(c *Catalog, col *Column) string {
	if col == nil || col.CollationName == "" {
		return ""
	}
	if sc := columnCollationSchema(c, col); sc != "" {
		return " COLLATE " + migrationQualifiedName(sc, col.CollationName)
	}
	return formatCollate(col.CollationName)
}

func columnExplicitStorage(c *Catalog, col *Column) string {
	if col == nil || col.Storage == 0 {
		return ""
	}
	if c != nil {
		if typ := c.typeByOID[col.TypeOID]; typ != nil && col.Storage == typ.Storage {
			return ""
		}
	}
	return storageModeKeyword(col.Storage)
}

func storageModeKeyword(stor byte) string {
	switch stor {
	case 'p':
		return "PLAIN"
	case 'e':
		return "EXTERNAL"
	case 'x':
		return "EXTENDED"
	case 'm':
		return "MAIN"
	}
	return ""
}

func columnExplicitCompression(col *Column) string {
	if col == nil {
		return ""
	}
	switch col.Compression {
	case 'p':
		return "pglz"
	case 'l':
		return "lz4"
	}
	return ""
}

// formatPKConstraint formats a PRIMARY KEY constraint clause.
func formatPKConstraint(rel *Relation, con *Constraint) string {
	cols := columnAttnumsToQuotedNames(rel, con.Columns)
	clause := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)",
		quoteIdentAlways(con.Name), strings.Join(cols, ", "))
	if len(con.IncludeColumns) > 0 {
		inc := columnAttnumsToQuotedNames(rel, con.IncludeColumns)
		clause += fmt.Sprintf(" INCLUDE (%s)", strings.Join(inc, ", "))
	}
	return clause
}

// formatUniqueConstraint formats a UNIQUE constraint clause.
func formatUniqueConstraint(c *Catalog, rel *Relation, con *Constraint) string {
	cols := columnAttnumsToQuotedNames(rel, con.Columns)
	nulls := ""
	if constraintNullsNotDistinct(c, con) {
		nulls = " NULLS NOT DISTINCT"
	}
	clause := fmt.Sprintf("CONSTRAINT %s UNIQUE%s (%s)",
		quoteIdentAlways(con.Name), nulls, strings.Join(cols, ", "))
	if len(con.IncludeColumns) > 0 {
		inc := columnAttnumsToQuotedNames(rel, con.IncludeColumns)
		clause += fmt.Sprintf(" INCLUDE (%s)", strings.Join(inc, ", "))
	}
	return clause
}

// formatCheckConstraint formats a CHECK constraint clause.
func formatCheckConstraint(con *Constraint) string {
	clause := fmt.Sprintf("CONSTRAINT %s CHECK (%s)",
		quoteIdentAlways(con.Name), con.CheckExpr)
	if con.ConNoInherit {
		clause += " NO INHERIT"
	}
	return clause
}

// resolveAttnumsToNames converts attnums to quoted column names.
func columnAttnumsToQuotedNames(rel *Relation, attnums []int16) []string {
	names := make([]string, 0, len(attnums))
	for _, attnum := range attnums {
		for _, col := range rel.Columns {
			if col.AttNum == attnum {
				names = append(names, quoteIdentAlways(col.Name))
				break
			}
		}
	}
	return names
}

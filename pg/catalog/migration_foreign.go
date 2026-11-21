package catalog

import (
	"fmt"
	"sort"
	"strings"
)

func generateForeignDataWrapperDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.ForeignDataWrappers {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			var b strings.Builder
			b.WriteString("CREATE FOREIGN DATA WRAPPER ")
			b.WriteString(quoteIdentifier(entry.To.Name))
			if entry.To.Handler != "" {
				b.WriteString(" HANDLER ")
				b.WriteString(formatCatalogIdentRef(entry.To.Handler))
			}
			if entry.To.Validator != "" {
				b.WriteString(" VALIDATOR ")
				b.WriteString(formatCatalogIdentRef(entry.To.Validator))
			}
			b.WriteString(formatForeignOptionsCreateClause(entry.To.Options))
			ops = append(ops, MigrationOp{
				Type:          OpCreateTable,
				ObjectName:    entry.Name,
				SQL:           b.String(),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'F',
				ObjOID:        entry.To.OID,
				Priority:      PriorityTable,
			})
			owner := effectiveForeignOwner(entry.To.OwnerName)
			ops = append(ops, MigrationOp{
				Type:          OpAlterTable,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER FOREIGN DATA WRAPPER %s OWNER TO %s", quoteIdentifier(entry.To.Name), quoteIdentifier(owner)),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'F',
				ObjOID:        entry.To.OID,
				Priority:      PriorityMetadata,
			})
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropTable,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP FOREIGN DATA WRAPPER %s CASCADE", quoteIdentifier(entry.From.Name)),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'F',
				ObjOID:        entry.From.OID,
				Priority:      PriorityFunction,
			})
		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			if entry.From.Name != entry.To.Name {
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER FOREIGN DATA WRAPPER %s RENAME TO %s", quoteIdentifier(entry.From.Name), quoteIdentifier(entry.To.Name)),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'F',
					ObjOID:        entry.To.OID,
					Priority:      PriorityFunction,
				})
			}
			var clauses []string
			if entry.From.Handler != entry.To.Handler {
				if entry.To.Handler == "" {
					clauses = append(clauses, "NO HANDLER")
				} else {
					clauses = append(clauses, "HANDLER "+formatCatalogIdentRef(entry.To.Handler))
				}
			}
			if entry.From.Validator != entry.To.Validator {
				if entry.To.Validator == "" {
					clauses = append(clauses, "NO VALIDATOR")
				} else {
					clauses = append(clauses, "VALIDATOR "+formatCatalogIdentRef(entry.To.Validator))
				}
			}
			if !foreignOptionsEqual(entry.From.Options, entry.To.Options) {
				if optClause := formatForeignOptionsAlterClause(entry.From.Options, entry.To.Options); optClause != "" {
					clauses = append(clauses, optClause)
				}
			}
			if len(clauses) > 0 {
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER FOREIGN DATA WRAPPER %s %s", quoteIdentifier(entry.To.Name), strings.Join(clauses, " ")),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'F',
					ObjOID:        entry.To.OID,
					Priority:      PriorityTable,
				})
			}
			if effectiveForeignOwner(entry.From.OwnerName) != effectiveForeignOwner(entry.To.OwnerName) {
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER FOREIGN DATA WRAPPER %s OWNER TO %s", quoteIdentifier(entry.To.Name), quoteIdentifier(effectiveForeignOwner(entry.To.OwnerName))),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'F',
					ObjOID:        entry.To.OID,
					Priority:      PriorityMetadata,
				})
			}
		}
	}
	return ops
}

func renamedForeignServerNames(diff *SchemaDiff) map[string]string {
	renamed := make(map[string]string)
	if diff == nil {
		return renamed
	}
	for _, srv := range diff.ForeignServers {
		if srv.Action == DiffModify && srv.From != nil && srv.To != nil && srv.From.Name != srv.To.Name {
			renamed[srv.From.Name] = srv.To.Name
		}
	}
	return renamed
}

func generateForeignServerDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	renamedFDWs := make(map[string]string)
	for _, fdw := range diff.ForeignDataWrappers {
		if fdw.Action == DiffModify && fdw.From != nil && fdw.To != nil && fdw.From.Name != fdw.To.Name {
			renamedFDWs[fdw.From.Name] = fdw.To.Name
		}
	}
	var ops []MigrationOp
	for _, entry := range diff.ForeignServers {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			ops = append(ops, buildCreateForeignServerOps(entry.To)...)
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropTable,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP SERVER %s CASCADE", quoteIdentifier(entry.From.Name)),
				Transactional: true,
				Phase:         PhasePre,
				Priority:      PriorityTable,
			})
		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			fdwChanged := entry.From.FdwName != entry.To.FdwName && renamedFDWs[entry.From.FdwName] != entry.To.FdwName
			if fdwChanged || !stringPtrEqual(entry.From.Type, entry.To.Type) {
				// DROP SERVER ... CASCADE would destroy unmodeled user mappings outside the catalog/snapshot.
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("-- foreign server %s FDW or TYPE changed; cannot alter in place without losing dependent user mappings", quoteIdentifier(entry.To.Name)),
					Warning:       fmt.Sprintf("foreign server %s cannot be ALTERed in place: FDW or TYPE changed and DROP SERVER CASCADE would remove dependent user mappings", quoteIdentifier(entry.To.Name)),
					Transactional: true,
					Phase:         PhaseMain,
					Priority:      PriorityTable,
				})
				continue
			}
			if entry.From.Name != entry.To.Name {
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER SERVER %s RENAME TO %s", quoteIdentifier(entry.From.Name), quoteIdentifier(entry.To.Name)),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'r',
					ObjOID:        entry.To.OID,
					Priority:      PriorityFunction,
				})
			}
			var clauses []string
			if !stringPtrEqual(entry.From.Version, entry.To.Version) {
				if entry.To.Version == nil {
					clauses = append(clauses, "VERSION NULL")
				} else {
					clauses = append(clauses, "VERSION "+quoteLiteral(*entry.To.Version))
				}
			}
			if !foreignOptionsEqual(entry.From.Options, entry.To.Options) {
				if optClause := formatForeignOptionsAlterClause(entry.From.Options, entry.To.Options); optClause != "" {
					clauses = append(clauses, optClause)
				}
			}
			if len(clauses) > 0 {
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER SERVER %s %s", quoteIdentifier(entry.To.Name), strings.Join(clauses, " ")),
					Transactional: true,
					Phase:         PhaseMain,
					Priority:      PriorityTable,
				})
			}
			if effectiveForeignOwner(entry.From.OwnerName) != effectiveForeignOwner(entry.To.OwnerName) {
				ops = append(ops, MigrationOp{
					Type:          OpAlterTable,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER SERVER %s OWNER TO %s", quoteIdentifier(entry.To.Name), quoteIdentifier(effectiveForeignOwner(entry.To.OwnerName))),
					Transactional: true,
					Phase:         PhaseMain,
					Priority:      PriorityMetadata,
				})
			}
		}
	}
	return ops
}

func buildCreateForeignServerOps(srv *ForeignServer) []MigrationOp {
	var b strings.Builder
	b.WriteString("CREATE SERVER ")
	b.WriteString(quoteIdentifier(srv.Name))
	if srv.Type != nil {
		b.WriteString(" TYPE ")
		b.WriteString(quoteLiteral(*srv.Type))
	}
	if srv.Version != nil {
		b.WriteString(" VERSION ")
		b.WriteString(quoteLiteral(*srv.Version))
	}
	b.WriteString(" FOREIGN DATA WRAPPER ")
	b.WriteString(quoteIdentifier(srv.FdwName))
	b.WriteString(formatForeignOptionsCreateClause(srv.Options))
	owner := effectiveForeignOwner(srv.OwnerName)
	ops := []MigrationOp{{
		Type:          OpCreateTable,
		ObjectName:    srv.Name,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		Priority:      PriorityTable,
	}, {
		Type:          OpAlterTable,
		ObjectName:    srv.Name,
		SQL:           fmt.Sprintf("ALTER SERVER %s OWNER TO %s", quoteIdentifier(srv.Name), quoteIdentifier(owner)),
		Transactional: true,
		Phase:         PhaseMain,
		Priority:      PriorityMetadata,
	}}
	return ops
}

// rebuiltForeignServerNames returns servers replaced via DROP + CREATE during
// migration. Foreign server FDW/TYPE changes refuse CASCADE rebuild to preserve
// unmodeled user mappings, so no server is rebuilt here.
func rebuiltForeignServerNames(_ *SchemaDiff) map[string]bool {
	return nil
}

func formatCatalogIdentRef(ref string) string {
	schemaName, objName := parseRoutineRefName(ref)
	if schemaName != "" {
		return quoteIdentifier(schemaName) + "." + quoteIdentifier(objName)
	}
	return quoteIdentifier(objName)
}

func generateForeignTableDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	renamedServers := renamedForeignServerNames(diff)
	var ops []MigrationOp

	for _, entry := range diff.ForeignTables {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil || entry.To.RelKind != 'f' {
				continue
			}
			ops = append(ops, buildCreateForeignTableOps(to, entry.SchemaName, entry.Name, entry.To)...)

		case DiffDrop:
			if entry.From == nil || entry.From.RelKind != 'f' {
				continue
			}
			qn := formatRelationTargetName(entry.SchemaName, entry.From)
			ops = append(ops, MigrationOp{
				Type:          OpDropTable,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP FOREIGN TABLE %s CASCADE", qn),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'r',
				ObjOID:        entry.From.OID,
				Priority:      PriorityColumn,
			})

		case DiffModify:
			if entry.From == nil || entry.To == nil || entry.To.RelKind != 'f' {
				continue
			}
			qn := formatRelationTargetName(entry.SchemaName, entry.To)
			if entry.From.ForeignServer != entry.To.ForeignServer && renamedServers[entry.From.ForeignServer] != entry.To.ForeignServer {
				fromQN := formatRelationTargetName(entry.SchemaName, entry.From)
				ops = append(ops, MigrationOp{
					Type:          OpDropTable,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("DROP FOREIGN TABLE %s CASCADE", fromQN),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        entry.From.OID,
					Priority:      PriorityColumn,
				})
				ops = append(ops, buildCreateForeignTableOps(to, entry.SchemaName, entry.Name, entry.To)...)
				ops = append(ops, restoreForeignTableDependents(to, diff, entry)...)
				continue
			}

			skipCheckDrops := make(map[string]bool)
			for _, ce := range entry.Constraints {
				if (ce.Action == DiffDrop || ce.Action == DiffModify) && ce.From != nil && ce.From.Type == ConstraintCheck {
					skipCheckDrops[ce.From.Name] = true
				}
			}

			if !foreignOptionsEqual(entry.From.ForeignOptions, entry.To.ForeignOptions) {
				if clause := formatForeignOptionsAlterClause(entry.From.ForeignOptions, entry.To.ForeignOptions); clause != "" {
					ops = append(ops, MigrationOp{
						Type:          OpAlterTable,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s %s", qn, clause),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        entry.To.OID,
						Priority:      PriorityTable,
					})
				}
			}

			for _, col := range entry.Columns {
				switch col.Action {
				case DiffAdd:
					if col.To == nil {
						continue
					}
					ops = append(ops, MigrationOp{
						Type:          OpAddColumn,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s ADD COLUMN %s", qn, formatForeignColumnDef(to, col.To)),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        entry.To.OID,
						Priority:      PriorityColumn,
					})
				case DiffDrop:
					if col.From == nil {
						continue
					}
					ops = append(ops, MigrationOp{
						Type:          OpDropColumn,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s DROP COLUMN %s", qn, quoteIdentifier(col.Name)),
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       'r',
						ObjOID:        entry.From.OID,
						Priority:      PriorityColumn,
					})
				case DiffModify:
					if col.From == nil || col.To == nil {
						continue
					}
					modOps := columnModifyOpsForStmt(from, to, "ALTER FOREIGN TABLE", qn, entry.SchemaName, entry.Name, col, skipCheckDrops)
					for i := range modOps {
						modOps[i].Phase = PhaseMain
						modOps[i].ObjType = 'r'
						modOps[i].ObjOID = entry.To.OID
						modOps[i].Priority = PriorityColumn
					}
					ops = append(ops, modOps...)
					if !foreignOptionsEqual(col.From.ForeignOptions, col.To.ForeignOptions) {
						if clause := formatForeignOptionsAlterClause(col.From.ForeignOptions, col.To.ForeignOptions); clause != "" {
							ops = append(ops, MigrationOp{
								Type:          OpAlterColumn,
								SchemaName:    entry.SchemaName,
								ObjectName:    entry.Name,
								SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s ALTER COLUMN %s %s", qn, quoteIdentifier(col.To.Name), clause),
								Transactional: true,
								Phase:         PhaseMain,
								ObjType:       'r',
								ObjOID:        entry.To.OID,
								Priority:      PriorityColumn,
							})
						}
					}
				}
			}

			for _, ce := range entry.Constraints {
				switch ce.Action {
				case DiffDrop:
					if ce.From != nil {
						ops = append(ops, MigrationOp{
							Type:          OpDropConstraint,
							SchemaName:    entry.SchemaName,
							ObjectName:    entry.Name,
							SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s DROP CONSTRAINT %s", qn, quoteIdentifier(ce.From.Name)),
							Transactional: true,
							Phase:         PhasePre,
							ObjType:       'r',
							ObjOID:        entry.From.OID,
							Priority:      PriorityConstraint,
						})
					}
				case DiffAdd:
					if ce.To != nil && ce.To.Type == ConstraintCheck {
						ops = append(ops, MigrationOp{
							Type:          OpAddConstraint,
							SchemaName:    entry.SchemaName,
							ObjectName:    entry.Name,
							SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s ADD %s", qn, formatForeignCheckConstraint(ce.To)),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        entry.To.OID,
							Priority:      PriorityConstraint,
						})
					}
				case DiffModify:
					if ce.From != nil {
						ops = append(ops, MigrationOp{
							Type:          OpDropConstraint,
							SchemaName:    entry.SchemaName,
							ObjectName:    entry.Name,
							SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s DROP CONSTRAINT %s", qn, quoteIdentifier(ce.From.Name)),
							Transactional: true,
							Phase:         PhasePre,
							ObjType:       'r',
							ObjOID:        entry.From.OID,
							Priority:      PriorityConstraint,
						})
					}
					if ce.To != nil && ce.To.Type == ConstraintCheck {
						ops = append(ops, MigrationOp{
							Type:          OpAddConstraint,
							SchemaName:    entry.SchemaName,
							ObjectName:    entry.Name,
							SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s ADD %s", qn, formatForeignCheckConstraint(ce.To)),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        entry.To.OID,
							Priority:      PriorityConstraint,
						})
					}
				}
			}

			if effectiveTableOwner(entry.From) != effectiveTableOwner(entry.To) {
				if ownerOp, ok := buildForeignTableOwnerOp(entry.SchemaName, entry.Name, entry.To); ok {
					ops = append(ops, ownerOp)
				}
			}

			ops = append(ops, triggerOpsForRelation(to, entry.SchemaName, entry.Name, entry.From, entry.Triggers)...)
		}
	}

	return ops
}

func formatForeignCheckConstraint(con *Constraint) string {
	sql := formatCheckConstraint(con)
	if con != nil && !con.Validated {
		sql += " NOT VALID"
	}
	return sql
}

func buildCreateForeignTableOps(to *Catalog, schemaName, tableName string, rel *Relation) []MigrationOp {
	qn := formatRelationName(schemaName, tableName, rel)
	var b strings.Builder
	b.WriteString("CREATE FOREIGN TABLE ")
	b.WriteString(qn)
	b.WriteString(" (")

	var checks []string
	if to != nil {
		for _, con := range to.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintCheck {
				checks = append(checks, formatForeignCheckConstraint(con))
			}
		}
	}
	sort.Strings(checks)

	if len(rel.Columns) > 0 || len(checks) > 0 {
		b.WriteByte('\n')
		first := true
		for _, col := range rel.Columns {
			if !first {
				b.WriteString(",\n")
			}
			first = false
			b.WriteString("    ")
			b.WriteString(formatForeignColumnDef(to, col))
		}
		for _, chk := range checks {
			if !first {
				b.WriteString(",\n")
			}
			first = false
			b.WriteString("    ")
			b.WriteString(chk)
		}
		b.WriteByte('\n')
	}
	b.WriteString(")")
	if rel.ForeignServer != "" {
		b.WriteString(" SERVER ")
		b.WriteString(quoteIdentifier(rel.ForeignServer))
	}
	b.WriteString(formatForeignOptionsCreateClause(rel.ForeignOptions))

	ops := []MigrationOp{{
		Type:          OpCreateTable,
		SchemaName:    schemaName,
		ObjectName:    tableName,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        rel.OID,
		Priority:      PriorityTable,
	}}
	if ownerOp, ok := buildForeignTableOwnerOp(schemaName, tableName, rel); ok {
		ops = append(ops, ownerOp)
	}
	if to != nil {
		for _, trig := range to.TriggersOf(rel.OID) {
			if trig.ParentTriggerOID == 0 {
				ops = append(ops, buildCreateAndModeTriggerOps(to, schemaName, tableName, trig)...)
			}
		}
	}
	return ops
}

func buildForeignTableOwnerOp(schemaName, tableName string, rel *Relation) (MigrationOp, bool) {
	if rel == nil || rel.OwnerName == "" || (!rel.OwnerExplicit && rel.OwnerName == "postgres") {
		return MigrationOp{}, false
	}
	qn := formatRelationName(schemaName, tableName, rel)
	return MigrationOp{
		Type:          OpAlterTable,
		SchemaName:    schemaName,
		ObjectName:    tableName,
		SQL:           fmt.Sprintf("ALTER FOREIGN TABLE %s OWNER TO %s", qn, quoteIdentifier(rel.OwnerName)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        rel.OID,
		Priority:      PriorityMetadata,
	}, true
}

func restoreForeignTableDependents(to *Catalog, diff *SchemaDiff, entry RelationDiffEntry) []MigrationOp {
	var ops []MigrationOp
	rel := entry.To
	if rel == nil {
		return nil
	}

	handledGrantSQL := make(map[string]bool)
	if diff != nil {
		for _, ge := range diff.Grants {
			if (ge.Action == DiffAdd || ge.Action == DiffModify) && ge.To.ObjType == 'r' && ge.To.ObjOID == rel.OID {
				handledGrantSQL[formatGrantSQL(ge.To, resolveGrantObjName(to, 'r', rel.OID))] = true
			}
		}
	}
	for _, gop := range restoreGrantOps(to, 'r', rel.OID) {
		if !handledGrantSQL[gop.SQL] {
			ops = append(ops, gop)
		}
	}

	handledComments := make(map[string]bool)
	if diff != nil {
		for _, ce := range diff.Comments {
			if ce.Action == DiffAdd || ce.Action == DiffModify {
				if sql := formatCommentEntrySQL(to, ce, ce.To); sql != "" {
					handledComments[sql] = true
				}
			}
		}
	}
	appendComment := func(objType byte, oid uint32, subID int16, desc string, parts []string) {
		if txt, ok := to.GetComment(objType, oid, subID); ok && txt != "" {
			if sql := formatCommentEntrySQL(to, CommentDiffEntry{
				ObjType:        objType,
				ObjDescription: desc,
				SubID:          subID,
				objParts:       parts,
			}, txt); sql != "" && !handledComments[sql] {
				handledComments[sql] = true
				ops = append(ops, MigrationOp{
					Type:          OpComment,
					ObjectName:    desc,
					SQL:           sql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       objType,
					Priority:      PriorityMetadata,
				})
			}
		}
	}
	relDesc := entry.SchemaName + "." + entry.Name
	appendComment('r', rel.OID, 0, relDesc, []string{entry.SchemaName, entry.Name})
	for _, col := range rel.Columns {
		appendComment('r', rel.OID, col.AttNum, relDesc, []string{entry.SchemaName, entry.Name})
	}
	for _, con := range to.ConstraintsOf(rel.OID) {
		appendComment('c', con.OID, 0, relDesc+"."+con.Name, []string{entry.SchemaName, entry.Name, con.Name})
	}
	for _, trig := range to.TriggersOf(rel.OID) {
		if trig.ParentTriggerOID == 0 {
			appendComment('g', trig.OID, 0, relDesc+"."+trig.Name, []string{entry.SchemaName, entry.Name, trig.Name})
		}
	}
	return ops
}

func formatForeignColumnDef(c *Catalog, col *Column) string {
	var parts []string
	parts = append(parts, quoteIdentAlways(col.Name))
	parts = append(parts, c.FormatType(col.TypeOID, col.TypeMod)+formatColumnCollate(c, col))
	if opts := strings.TrimPrefix(formatForeignOptionsCreateClause(col.ForeignOptions), " "); opts != "" {
		parts = append(parts, opts)
	}
	if col.NotNull {
		parts = append(parts, "NOT NULL")
	}
	if col.Generated == 's' {
		expr := col.GenerationExpr
		if expr == "" {
			expr = col.Default
		}
		if expr != "" {
			parts = append(parts, "GENERATED ALWAYS AS ("+expr+") STORED")
		}
	} else if col.HasDefault && col.Default != "" {
		parts = append(parts, "DEFAULT "+col.Default)
	}
	return strings.Join(parts, " ")
}

func splitForeignOption(entry string) (string, string) {
	k, v, ok := strings.Cut(entry, "=")
	if !ok {
		return entry, ""
	}
	return k, v
}

func formatForeignOptionsCreateClause(options []string) string {
	if len(options) == 0 {
		return ""
	}
	parts := make([]string, 0, len(options))
	for _, opt := range options {
		k, v := splitForeignOption(opt)
		parts = append(parts, fmt.Sprintf("%s %s", quoteIdentifier(k), quoteLiteral(v)))
	}
	return " OPTIONS (" + strings.Join(parts, ", ") + ")"
}

func formatForeignOptionsAlterClause(oldOpts, newOpts []string) string {
	oldMap := make(map[string]string, len(oldOpts))
	var oldOrder []string
	for _, opt := range oldOpts {
		k, v := splitForeignOption(opt)
		oldMap[k] = v
		oldOrder = append(oldOrder, k)
	}
	newMap := make(map[string]string, len(newOpts))
	var newOrder []string
	for _, opt := range newOpts {
		k, v := splitForeignOption(opt)
		newMap[k] = v
		newOrder = append(newOrder, k)
	}

	var actions []string
	for _, k := range oldOrder {
		if _, ok := newMap[k]; !ok {
			actions = append(actions, "DROP "+quoteIdentifier(k))
		}
	}
	for _, k := range newOrder {
		nv := newMap[k]
		if ov, ok := oldMap[k]; !ok {
			actions = append(actions, fmt.Sprintf("ADD %s %s", quoteIdentifier(k), quoteLiteral(nv)))
		} else if ov != nv {
			actions = append(actions, fmt.Sprintf("SET %s %s", quoteIdentifier(k), quoteLiteral(nv)))
		}
	}
	if len(actions) == 0 {
		return ""
	}
	return "OPTIONS (" + strings.Join(actions, ", ") + ")"
}

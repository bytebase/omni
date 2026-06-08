package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateConstraintDDL produces DDL operations for constraint changes.
func generateConstraintDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, rel := range diff.Relations {
		switch rel.Action {
		case DiffModify:
			ops = append(ops, constraintOpsForRelation(to, rel)...)
		case DiffAdd:
			// For newly added relations, inline constraints (PK, UNIQUE, CHECK)
			// are handled by generateTableDDL. FK constraints should be deferred
			// here so they are emitted after all tables are created.
			if rel.To != nil {
				// First check diff entry's constraint slice (may be empty for new tables).
				fkFound := false
				for _, ce := range rel.Constraints {
					if ce.Action == DiffAdd && ce.To != nil && ce.To.Type == ConstraintFK {
						ops = append(ops, buildAddConstraintOp(to, rel.SchemaName, rel.Name, ce.To))
						fkFound = true
					}
				}
				// If no FKs in diff entry, look up from the catalog directly.
				if !fkFound {
					for _, con := range to.ConstraintsOf(rel.To.OID) {
						if con.Type == ConstraintFK {
							ops = append(ops, buildAddConstraintOp(to, rel.SchemaName, rel.Name, con))
						}
					}
				}
			}
		}
	}

	// Sort for determinism: by schema, table, then constraint name.
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		if ops[i].ParentObject != ops[j].ParentObject {
			return ops[i].ParentObject < ops[j].ParentObject
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// constraintOpsForRelation generates constraint DDL for a modified relation.
func constraintOpsForRelation(to *Catalog, rel RelationDiffEntry) []MigrationOp {
	var ops []MigrationOp

	for _, ce := range rel.Constraints {
		// Skip constraint-trigger constraints — they are handled by trigger DDL.
		if (ce.To != nil && ce.To.Type == ConstraintTrigger) ||
			(ce.From != nil && ce.From.Type == ConstraintTrigger) {
			continue
		}
		switch ce.Action {
		case DiffAdd:
			if ce.To == nil {
				continue
			}
			ops = append(ops, buildAddConstraintOp(to, rel.SchemaName, rel.Name, ce.To))
		case DiffDrop:
			if ce.From == nil {
				continue
			}
			ops = append(ops, buildDropConstraintOp(rel.SchemaName, rel.Name, ce.From.Name))
		case DiffModify:
			// Modified constraint = DROP old + ADD new.
			if ce.From != nil {
				ops = append(ops, buildDropConstraintOp(rel.SchemaName, rel.Name, ce.From.Name))
			}
			if ce.To != nil {
				ops = append(ops, buildAddConstraintOp(to, rel.SchemaName, rel.Name, ce.To))
			}
		}
	}

	return ops
}

// buildDropConstraintOp creates a DROP CONSTRAINT operation.
func buildDropConstraintOp(schemaName, tableName, conName string) MigrationOp {
	qualifiedTable := quoteIdentifier(schemaName) + "." + quoteIdentifier(tableName)
	return MigrationOp{
		Type:         OpDropConstraint,
		SchemaName:   schemaName,
		ObjectName:   conName,
		ParentObject: tableName,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", qualifiedTable, quoteIdentifier(conName)),
		Phase:        PhasePre,
		ObjType:      'c',
		Priority:     PriorityConstraint,
	}
}

// buildAddConstraintOp creates an ADD CONSTRAINT operation.
func buildAddConstraintOp(to *Catalog, schemaName, tableName string, con *Constraint) MigrationOp {
	qualifiedTable := quoteIdentifier(schemaName) + "." + quoteIdentifier(tableName)
	var sql string

	switch con.Type {
	case ConstraintPK:
		colNames := resolveAttnumsToNames(to, con.RelOID, con.Columns)
		clause := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY (%s)",
			qualifiedTable, quoteIdentifier(con.Name), joinQuoted(colNames))
		if con.Deferrable {
			clause += " DEFERRABLE"
			if con.Deferred {
				clause += " INITIALLY DEFERRED"
			} else {
				clause += " INITIALLY IMMEDIATE"
			}
		}
		sql = clause

	case ConstraintUnique:
		colNames := resolveAttnumsToNames(to, con.RelOID, con.Columns)
		clause := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)",
			qualifiedTable, quoteIdentifier(con.Name), joinQuoted(colNames))
		if con.Deferrable {
			clause += " DEFERRABLE"
			if con.Deferred {
				clause += " INITIALLY DEFERRED"
			} else {
				clause += " INITIALLY IMMEDIATE"
			}
		}
		sql = clause

	case ConstraintCheck:
		sql = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s CHECK (%s)",
			qualifiedTable, quoteIdentifier(con.Name), con.CheckExpr)

	case ConstraintFK:
		colNames := resolveAttnumsToNames(to, con.RelOID, con.Columns)
		refColNames := resolveAttnumsToNames(to, con.FRelOID, con.FColumns)
		refTable := resolveRelationName(to, con.FRelOID)

		fkClause := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			qualifiedTable, quoteIdentifier(con.Name),
			joinQuoted(colNames), refTable, joinQuoted(refColNames))

		// ON DELETE action.
		if action := fkActionString(con.FKDelAction); action != "" {
			fkClause += " ON DELETE " + action
		}
		// ON UPDATE action.
		if action := fkActionString(con.FKUpdAction); action != "" {
			fkClause += " ON UPDATE " + action
		}
		// DEFERRABLE.
		if con.Deferrable {
			if con.Deferred {
				fkClause += " DEFERRABLE INITIALLY DEFERRED"
			} else {
				fkClause += " DEFERRABLE INITIALLY IMMEDIATE"
			}
		}

		sql = fkClause

	case ConstraintExclude:
		colNames := resolveAttnumsToNames(to, con.RelOID, con.Columns)
		var parts []string
		for i, colName := range colNames {
			op := "="
			if i < len(con.ExclOps) {
				op = con.ExclOps[i]
			}
			parts = append(parts, fmt.Sprintf("%s WITH %s", quoteIdentifier(colName), op))
		}
		// Determine the index access method from the backing index.
		method := "gist" // default for EXCLUDE constraints
		if con.IndexOID != 0 {
			if idx := to.GetIndexByOID(con.IndexOID); idx != nil && idx.AccessMethod != "" {
				method = idx.AccessMethod
			}
		}
		clause := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s EXCLUDE USING %s (%s)",
			qualifiedTable, quoteIdentifier(con.Name), method, strings.Join(parts, ", "))
		if con.Deferrable {
			clause += " DEFERRABLE"
			if con.Deferred {
				clause += " INITIALLY DEFERRED"
			} else {
				clause += " INITIALLY IMMEDIATE"
			}
		}
		sql = clause

	default:
		sql = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s /* unsupported type %c */",
			qualifiedTable, quoteIdentifier(con.Name), byte(con.Type))
	}

	// NOT VALID for unvalidated constraints.
	if !con.Validated {
		sql += " NOT VALID"
	}

	// FK constraints are deferred to PhasePost; all others are PhaseMain.
	phase := PhaseMain
	priority := PriorityConstraint
	if con.Type == ConstraintFK {
		phase = PhasePost
		priority = PriorityFKDeferred
	}

	return MigrationOp{
		Type:         OpAddConstraint,
		SchemaName:   schemaName,
		ObjectName:   con.Name,
		ParentObject: tableName,
		SQL:          sql,
		Phase:        phase,
		ObjType:      'c',
		ObjOID:       con.OID,
		Priority:     priority,
	}
}

// resolveAttnumsToNames converts attnum slice to column name slice using the catalog.
func resolveAttnumsToNames(c *Catalog, relOID uint32, attnums []int16) []string {
	rel := c.GetRelationByOID(relOID)
	if rel == nil {
		// Fallback: return attnum-based names.
		names := make([]string, len(attnums))
		for i, n := range attnums {
			names[i] = fmt.Sprintf("col%d", n)
		}
		return names
	}
	names := make([]string, len(attnums))
	for i, attnum := range attnums {
		found := false
		for _, col := range rel.Columns {
			if col.AttNum == attnum {
				names[i] = col.Name
				found = true
				break
			}
		}
		if !found {
			names[i] = fmt.Sprintf("col%d", attnum)
		}
	}
	return names
}

// resolveRelationName returns a schema-qualified, quoted table name for a given OID.
func resolveRelationName(c *Catalog, relOID uint32) string {
	rel := c.GetRelationByOID(relOID)
	if rel == nil {
		return fmt.Sprintf("/* unknown OID %d */", relOID)
	}
	schemaName := "public"
	if rel.Schema != nil {
		schemaName = rel.Schema.Name
	}
	return quoteIdentifier(schemaName) + "." + quoteIdentifier(rel.Name)
}

// fkActionString maps a FK action byte to SQL clause text.
// Returns empty string for NO ACTION (the default), since it can be omitted.
func fkActionString(action byte) string {
	switch action {
	case 'r':
		return "RESTRICT"
	case 'c':
		return "CASCADE"
	case 'n':
		return "SET NULL"
	case 'd':
		return "SET DEFAULT"
	case 'a', 0:
		return "" // NO ACTION is the default, omit
	default:
		return "NO ACTION"
	}
}

// joinQuoted joins column names with double-quoting.
func joinQuoted(names []string) string {
	quoted := make([]string, len(names))
	for i, n := range names {
		quoted[i] = quoteIdentifier(n)
	}
	return strings.Join(quoted, ", ")
}

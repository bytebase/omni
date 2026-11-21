package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generatePolicyDDL produces CREATE POLICY, DROP POLICY, ALTER POLICY, and
// RLS ALTER TABLE operations from the diff.
func generatePolicyDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, rel := range diff.Relations {
		targetRel := rel.To
		if targetRel == nil {
			targetRel = rel.From
		}
		qualifiedTable := formatRelationName(rel.SchemaName, rel.Name, targetRel)

		// RLS flag changes (ENABLE/DISABLE ROW LEVEL SECURITY, FORCE/NO FORCE).
		if rel.RLSChanged {
			if rel.From != nil && rel.To != nil {
				if rel.From.RowSecurity != rel.To.RowSecurity {
					if rel.To.RowSecurity {
						ops = append(ops, MigrationOp{
							Type:          OpAlterTable,
							SchemaName:    rel.SchemaName,
							ObjectName:    rel.Name,
							SQL:           fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", qualifiedTable),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        rel.To.OID,
							Priority:      PriorityTable,
						})
					} else {
						ops = append(ops, MigrationOp{
							Type:          OpAlterTable,
							SchemaName:    rel.SchemaName,
							ObjectName:    rel.Name,
							SQL:           fmt.Sprintf("ALTER TABLE %s DISABLE ROW LEVEL SECURITY", qualifiedTable),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        rel.To.OID,
							Priority:      PriorityTable,
						})
					}
				}
				if rel.From.ForceRowSecurity != rel.To.ForceRowSecurity {
					if rel.To.ForceRowSecurity {
						ops = append(ops, MigrationOp{
							Type:          OpAlterTable,
							SchemaName:    rel.SchemaName,
							ObjectName:    rel.Name,
							SQL:           fmt.Sprintf("ALTER TABLE %s FORCE ROW LEVEL SECURITY", qualifiedTable),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        rel.To.OID,
							Priority:      PriorityTable,
						})
					} else {
						ops = append(ops, MigrationOp{
							Type:          OpAlterTable,
							SchemaName:    rel.SchemaName,
							ObjectName:    rel.Name,
							SQL:           fmt.Sprintf("ALTER TABLE %s NO FORCE ROW LEVEL SECURITY", qualifiedTable),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        rel.To.OID,
							Priority:      PriorityTable,
						})
					}
				}
			}
		}

		// Policy add/drop/modify within a modified or added relation.
		switch rel.Action {
		case DiffModify:
			if rel.From != nil && rel.To != nil && tableRebuilt(from, to, rel.From, rel.To) {
				// DROP + CREATE rebuilds the table; restoreRebuiltTableDependents recreates its target policies.
				continue
			}
			ops = append(ops, policyOpsForRelation(rel.SchemaName, rel.Name, targetRel, rel.Policies)...)
		case DiffAdd:
			// For newly added relations, emit ENABLE ROW LEVEL SECURITY / FORCE ROW LEVEL SECURITY
			// if the target table has RLS enabled.
			if rel.To != nil {
				if rel.To.RowSecurity {
					ops = append(ops, MigrationOp{
						Type:          OpAlterTable,
						SchemaName:    rel.SchemaName,
						ObjectName:    rel.Name,
						SQL:           fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", qualifiedTable),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        rel.To.OID,
						Priority:      PriorityTable,
					})
				}
				if rel.To.ForceRowSecurity {
					ops = append(ops, MigrationOp{
						Type:          OpAlterTable,
						SchemaName:    rel.SchemaName,
						ObjectName:    rel.Name,
						SQL:           fmt.Sprintf("ALTER TABLE %s FORCE ROW LEVEL SECURITY", qualifiedTable),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        rel.To.OID,
						Priority:      PriorityTable,
					})
				}

				// Policies come from the target catalog.
				policies := to.policiesByRel[rel.To.OID]
				for _, p := range policies {
					ops = append(ops, buildCreatePolicyOpForRel(rel.SchemaName, rel.Name, rel.To, p))
				}
			}
		}
	}

	// Deterministic ordering.
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

// policyOpsForRelation generates policy DDL for a modified relation.
func policyOpsForRelation(schemaName, tableName string, targetRel *Relation, policies []PolicyDiffEntry) []MigrationOp {
	var ops []MigrationOp

	for _, pe := range policies {
		switch pe.Action {
		case DiffAdd:
			if pe.To == nil {
				continue
			}
			ops = append(ops, buildCreatePolicyOpForRel(schemaName, tableName, targetRel, pe.To))
		case DiffDrop:
			if pe.From == nil {
				continue
			}
			ops = append(ops, buildDropPolicyOpForRel(schemaName, tableName, targetRel, pe.From.Name, pe.From.OID))
		case DiffModify:
			if pe.From == nil || pe.To == nil {
				continue
			}
			// If CmdType or Permissive changed, it's a complex change: DROP + CREATE.
			if pe.From.CmdType != pe.To.CmdType || pe.From.Permissive != pe.To.Permissive {
				ops = append(ops, buildDropPolicyOpForRel(schemaName, tableName, targetRel, pe.From.Name, pe.From.OID))
				ops = append(ops, buildCreatePolicyOpForRel(schemaName, tableName, targetRel, pe.To))
			} else {
				// Simple change: ALTER POLICY.
				ops = append(ops, buildAlterPolicyOpForRel(schemaName, tableName, targetRel, pe.From, pe.To))
			}
		}
	}

	return ops
}

// buildCreatePolicyOp creates a CREATE POLICY operation.
func buildCreatePolicyOp(schemaName, tableName string, p *Policy) MigrationOp {
	return buildCreatePolicyOpForRel(schemaName, tableName, nil, p)
}

func buildCreatePolicyOpForRel(schemaName, tableName string, rel *Relation, p *Policy) MigrationOp {
	qualifiedTable := formatRelationName(schemaName, tableName, rel)
	var b strings.Builder
	b.WriteString("CREATE POLICY ")
	b.WriteString(quoteIdentAlways(p.Name))
	b.WriteString(" ON ")
	b.WriteString(qualifiedTable)

	// AS PERMISSIVE/RESTRICTIVE
	if !p.Permissive {
		b.WriteString(" AS RESTRICTIVE")
	}

	// FOR command
	if p.CmdType != "" && p.CmdType != "all" {
		b.WriteString(" FOR ")
		b.WriteString(strings.ToUpper(p.CmdType))
	}

	// TO roles
	if len(p.Roles) > 0 {
		b.WriteString(" TO ")
		b.WriteString(formatPolicyRoles(p.Roles))
	}

	// USING
	if p.UsingExpr != "" {
		b.WriteString(" USING (")
		b.WriteString(p.UsingExpr)
		b.WriteString(")")
	}

	// WITH CHECK
	if p.CheckExpr != "" {
		b.WriteString(" WITH CHECK (")
		b.WriteString(p.CheckExpr)
		b.WriteString(")")
	}

	return MigrationOp{
		Type:          OpCreatePolicy,
		SchemaName:    schemaName,
		ObjectName:    p.Name,
		ParentObject:  tableName,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'p',
		ObjOID:        p.OID,
		Priority:      PriorityPolicy,
	}
}

// buildDropPolicyOp creates a DROP POLICY operation.
func buildDropPolicyOp(schemaName, tableName, policyName string, policyOID uint32) MigrationOp {
	return buildDropPolicyOpForRel(schemaName, tableName, nil, policyName, policyOID)
}

func buildDropPolicyOpForRel(schemaName, tableName string, rel *Relation, policyName string, policyOID uint32) MigrationOp {
	qualifiedTable := formatRelationName(schemaName, tableName, rel)
	return MigrationOp{
		Type:          OpDropPolicy,
		SchemaName:    schemaName,
		ObjectName:    policyName,
		ParentObject:  tableName,
		SQL:           fmt.Sprintf("DROP POLICY %s ON %s", quoteIdentAlways(policyName), qualifiedTable),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'p',
		ObjOID:        policyOID,
		Priority:      PriorityPolicy,
	}
}

// buildAlterPolicyOp creates an ALTER POLICY operation for simple changes.
func buildAlterPolicyOp(schemaName, tableName string, from, to *Policy) MigrationOp {
	return buildAlterPolicyOpForRel(schemaName, tableName, nil, from, to)
}

func buildAlterPolicyOpForRel(schemaName, tableName string, rel *Relation, from, to *Policy) MigrationOp {
	qualifiedTable := formatRelationName(schemaName, tableName, rel)
	var b strings.Builder
	b.WriteString("ALTER POLICY ")
	b.WriteString(quoteIdentAlways(to.Name))
	b.WriteString(" ON ")
	b.WriteString(qualifiedTable)

	// Roles changed?
	if !stringSliceEqualSorted(from.Roles, to.Roles) {
		b.WriteString(" TO ")
		b.WriteString(formatPolicyRoles(to.Roles))
	}

	// USING changed?
	if from.UsingExpr != to.UsingExpr {
		b.WriteString(" USING (")
		b.WriteString(to.UsingExpr)
		b.WriteString(")")
	}

	// WITH CHECK changed?
	if from.CheckExpr != to.CheckExpr {
		b.WriteString(" WITH CHECK (")
		b.WriteString(to.CheckExpr)
		b.WriteString(")")
	}

	return MigrationOp{
		Type:          OpAlterPolicy,
		SchemaName:    schemaName,
		ObjectName:    to.Name,
		ParentObject:  tableName,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'p',
		ObjOID:        to.OID,
		Priority:      PriorityPolicy,
	}
}

// formatPolicyRoles formats a list of role names for a POLICY DDL clause.
func formatPolicyRoles(roles []string) string {
	quoted := make([]string, len(roles))
	for i, r := range roles {
		if r == "public" || r == "" {
			quoted[i] = "PUBLIC"
		} else {
			quoted[i] = quoteIdentAlways(r)
		}
	}
	return strings.Join(quoted, ", ")
}

// stringSliceEqualSorted compares two string slices after sorting copies.
func stringSliceEqualSorted(a, b []string) bool {
	ac := make([]string, len(a))
	copy(ac, a)
	sort.Strings(ac)
	bc := make([]string, len(b))
	copy(bc, b)
	sort.Strings(bc)
	return stringSliceEqual(ac, bc)
}

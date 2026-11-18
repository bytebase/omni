package catalog

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// generateGrantDDL produces GRANT and REVOKE operations from the diff.
func generateGrantDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Grants {
		switch entry.Action {
		case DiffAdd:
			g := entry.To
			objName := resolveGrantObjName(to, g.ObjType, g.ObjOID)
			sql := formatGrantSQL(g, objName)
			ops = append(ops, MigrationOp{
				Type:          OpGrant,
				ObjectName:    objName,
				SQL:           sql,
				Transactional: grantTransactional(g),
				Phase:         PhaseMain,
				ObjType:       g.ObjType,
				ObjOID:        g.ObjOID,
				Priority:      PriorityMetadata,
			})

		case DiffDrop:
			g := entry.From
			objName := resolveGrantObjName(from, g.ObjType, g.ObjOID)
			sql := formatRevokeSQL(g, objName)
			ops = append(ops, MigrationOp{
				Type:          OpRevoke,
				ObjectName:    objName,
				SQL:           sql,
				Transactional: grantTransactional(g),
				Phase:         PhaseMain,
				ObjType:       g.ObjType,
				ObjOID:        g.ObjOID,
				Priority:      PriorityMetadata,
			})

		case DiffModify:
			// Modified grant (WithGrant changed): REVOKE old + GRANT new.
			gFrom := entry.From
			gTo := entry.To
			objNameFrom := resolveGrantObjName(from, gFrom.ObjType, gFrom.ObjOID)
			objNameTo := resolveGrantObjName(to, gTo.ObjType, gTo.ObjOID)
			ops = append(ops, MigrationOp{
				Type:          OpRevoke,
				ObjectName:    objNameFrom,
				SQL:           formatRevokeSQL(gFrom, objNameFrom),
				Transactional: grantTransactional(gFrom),
				Phase:         PhaseMain,
				ObjType:       gFrom.ObjType,
				ObjOID:        gFrom.ObjOID,
				Priority:      PriorityMetadata,
			})
			ops = append(ops, MigrationOp{
				Type:          OpGrant,
				ObjectName:    objNameTo,
				SQL:           formatGrantSQL(gTo, objNameTo),
				Transactional: grantTransactional(gTo),
				Phase:         PhaseMain,
				ObjType:       gTo.ObjType,
				ObjOID:        gTo.ObjOID,
				Priority:      PriorityMetadata,
			})
		}
	}

	// Deterministic ordering: revokes first, then grants; within each by name.
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			if ops[i].Type == OpRevoke {
				return true
			}
			if ops[j].Type == OpRevoke {
				return false
			}
		}
		return ops[i].SQL < ops[j].SQL
	})

	return ops
}

// grantTransactional reports whether a grant runs inside a transaction block. A
// GRANT on an external resource cannot, the docs say
// (https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html).
func grantTransactional(g Grant) bool {
	return g.Target != nodes.OBJECT_FOREIGN_TABLE && g.Target != nodes.OBJECT_EXTERNAL_SCHEMA
}

// formatGrantSQL generates a GRANT statement.
func formatGrantSQL(g Grant, objName string) string {
	var b strings.Builder
	b.WriteString("GRANT ")
	b.WriteString(formatGrantPrivilege(g))
	b.WriteString(" ON ")
	b.WriteString(formatGrantObjectRef(g, objName))
	b.WriteString(" TO ")
	b.WriteString(formatGrantee(g))
	if g.WithGrant {
		b.WriteString(" WITH GRANT OPTION")
	}
	return b.String()
}

// formatRevokeSQL generates a REVOKE statement.
func formatRevokeSQL(g Grant, objName string) string {
	var b strings.Builder
	b.WriteString("REVOKE ")
	b.WriteString(formatGrantPrivilege(g))
	b.WriteString(" ON ")
	b.WriteString(formatGrantObjectRef(g, objName))
	b.WriteString(" FROM ")
	b.WriteString(formatGrantee(g))
	return b.String()
}

// formatGrantPrivilege formats the privilege portion.
func formatGrantPrivilege(g Grant) string {
	priv := strings.ToUpper(g.Privilege)
	if len(g.Columns) > 0 {
		cols := make([]string, len(g.Columns))
		for i, c := range g.Columns {
			cols[i] = quoteIdentAlways(c)
		}
		return fmt.Sprintf("%s (%s)", priv, strings.Join(cols, ", "))
	}
	return priv
}

// formatGrantObjectRef formats the ON clause for a grant/revoke, including the
// Redshift objects whose code alone does not say how they are written.
func formatGrantObjectRef(g Grant, objName string) string {
	switch g.Target {
	case nodes.OBJECT_PROCEDURE:
		return fmt.Sprintf("PROCEDURE %s", objName)
	case nodes.OBJECT_FOREIGN_TABLE:
		return fmt.Sprintf("EXTERNAL TABLE %s", formatQualifiedDescription(objName))
	case nodes.OBJECT_EXTERNAL_SCHEMA:
		return fmt.Sprintf("EXTERNAL SCHEMA %s", quoteIdentAlways(objName))
	}
	switch g.ObjType {
	case 'r':
		return fmt.Sprintf("TABLE %s", formatQualifiedDescription(objName))
	case 's':
		return fmt.Sprintf("SEQUENCE %s", formatQualifiedDescription(objName))
	case 'f':
		return fmt.Sprintf("FUNCTION %s", objName)
	case 'n':
		return fmt.Sprintf("SCHEMA %s", quoteIdentAlways(objName))
	case 'T':
		return fmt.Sprintf("TYPE %s", formatQualifiedDescription(objName))
	default:
		return quoteIdentAlways(objName)
	}
}

// formatGrantee formats the grantee for GRANT/REVOKE, keeping the alternative
// the statement was written with: an IAM role, a namespace or account, the
// admins of an account, a datashare, or a role name.
func formatGrantee(g Grant) string {
	name := g.GranteeName
	switch g.GranteeKind {
	case "iam_role":
		return fmt.Sprintf("IAM_ROLE '%s'", strings.ReplaceAll(name, "'", "''"))
	case "namespace":
		return fmt.Sprintf("NAMESPACE '%s'", strings.ReplaceAll(name, "'", "''"))
	case "account":
		return fmt.Sprintf("ACCOUNT '%s'", strings.ReplaceAll(name, "'", "''"))
	case "account_via_data_catalog":
		return fmt.Sprintf("ACCOUNT '%s' VIA DATA CATALOG", strings.ReplaceAll(name, "'", "''"))
	case "admins_of_account":
		return fmt.Sprintf("ADMINS OF ACCOUNT '%s'", strings.ReplaceAll(name, "'", "''"))
	case "datashare":
		return "DATASHARE " + quoteIdentAlways(name)
	case "rls_policy":
		return "RLS POLICY " + quoteIdentAlways(name)
	}
	if g.Grantee == "" {
		return "PUBLIC"
	}
	return quoteIdentAlways(g.Grantee)
}

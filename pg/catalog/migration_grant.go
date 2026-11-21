package catalog

import (
	"fmt"
	"sort"
	"strings"
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
				Transactional: true,
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
				Transactional: true,
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
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       gFrom.ObjType,
				ObjOID:        gFrom.ObjOID,
				Priority:      PriorityMetadata,
			})
			ops = append(ops, MigrationOp{
				Type:          OpGrant,
				ObjectName:    objNameTo,
				SQL:           formatGrantSQL(gTo, objNameTo),
				Transactional: true,
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

// restoreGrantOps grants again every privilege an object holds in c. A
// migration that drops and recreates the object while its grants stay the
// same has no grant diff, so it must restore them itself.
func restoreGrantOps(c *Catalog, objType byte, oid uint32) []MigrationOp {
	objName := resolveGrantObjName(c, objType, oid)
	// A temporary relation or sequence is created unqualified in pg_temp.
	switch objType {
	case 'r':
		if rel := c.GetRelationByOID(oid); rel != nil && (rel.Persistence == 't' || (rel.Schema != nil && isTempSchemaName(rel.Schema.Name))) {
			objName = rel.Name
		}
	case 's':
		if seq := c.GetSequenceByOID(oid); seq != nil && (seq.Persistence == 't' || (seq.Schema != nil && isTempSchemaName(seq.Schema.Name))) {
			objName = seq.Name
		}
	}
	var ops []MigrationOp
	for _, g := range c.grants {
		if g.ObjType != objType || g.ObjOID != oid || g.FromExtension {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpGrant,
			ObjectName:    objName,
			SQL:           formatGrantSQL(g, objName),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       objType,
			ObjOID:        oid,
			Priority:      PriorityMetadata,
		})
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i].SQL < ops[j].SQL })
	return ops
}

// formatGrantSQL generates a GRANT statement.
func formatGrantSQL(g Grant, objName string) string {
	var b strings.Builder
	b.WriteString("GRANT ")
	b.WriteString(formatGrantPrivilege(g))
	b.WriteString(" ON ")
	b.WriteString(formatGrantObjectRef(g.ObjType, objName))
	b.WriteString(" TO ")
	b.WriteString(formatGrantee(g.Grantee))
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
	b.WriteString(formatGrantObjectRef(g.ObjType, objName))
	b.WriteString(" FROM ")
	b.WriteString(formatGrantee(g.Grantee))
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

// formatGrantObjectRef formats the ON clause for a grant/revoke.
func formatGrantObjectRef(objType byte, objName string) string {
	switch objType {
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

// formatGrantee formats the grantee name for GRANT/REVOKE.
func formatGrantee(grantee string) string {
	if grantee == "" {
		return "PUBLIC"
	}
	return quoteIdentAlways(grantee)
}

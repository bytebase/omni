package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowDropObject reports every DROP: the DROP statement of any
// object kind, DROP DATABASE, DROP ROLE, DROP TABLESPACE, DROP OWNED,
// DROP SUBSCRIPTION, DROP USER MAPPING, and ALTER TABLE ... DROP COLUMN,
// which anchors the subcommand. IF EXISTS does not exempt a statement.
// Dropping a constraint belongs to DisallowDropConstraint.
func checkDisallowDropObject(s *statement, r *reporter) {
	report := func(rng review.Range, kind string, names []string) {
		msg := "drops " + kind
		if len(names) > 0 {
			msg += " " + strings.Join(names, ", ")
		}
		r.report(review.DisallowDropObject, s.index, rng, msg)
	}
	switch v := s.node.(type) {
	case *ast.DropStmt:
		kind := ast.ObjectType(v.RemoveType)
		onRelation := kind == ast.OBJECT_TRIGGER || kind == ast.OBJECT_POLICY || kind == ast.OBJECT_RULE
		var names []string
		if v.Objects != nil {
			for _, obj := range v.Objects.Items {
				if name := objectName(obj, onRelation); name != "" {
					names = append(names, name)
				}
			}
		}
		report(rangeOf(v.Loc), objectKind(kind), names)
	case *ast.DropdbStmt:
		report(rangeOf(v.Loc), "database", []string{ident(v.Dbname)})
	case *ast.DropRoleStmt:
		report(rangeOf(v.Loc), "role", roleNames(v.Roles))
	case *ast.DropTableSpaceStmt:
		report(rangeOf(v.Loc), "tablespace", []string{ident(v.Tablespacename)})
	case *ast.DropSubscriptionStmt:
		report(rangeOf(v.Loc), "subscription", []string{ident(v.Subname)})
	case *ast.DropUserMappingStmt:
		report(rangeOf(v.Loc), "user mapping", []string{"for " + roleName(v.User) + " on server " + ident(v.Servername)})
	case *ast.DropOwnedStmt:
		report(rangeOf(v.Loc), "objects owned by", roleNames(v.Roles))
	case *ast.AlterTableStmt:
		if v.Cmds == nil {
			return
		}
		for _, item := range v.Cmds.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok || ast.AlterTableType(cmd.Subtype) != ast.AT_DropColumn {
				continue
			}
			report(rangeOf(cmd.Loc), "column", []string{ident(cmd.Name) + " of " + relation(v.Relation)})
		}
	}
}

package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowDropObject reports every DROP: the DROP statement of any
// object kind, DROP DATABASE, DROP ROLE, DROP TABLESPACE, DROP OWNED,
// DROP SUBSCRIPTION, DROP USER MAPPING, and ALTER TABLE ... DROP COLUMN,
// which anchors the subcommand. IF EXISTS does not exempt a statement,
// and neither does a SQL-standard routine body around it: the statement
// tree is walked, so a DROP inside BEGIN ATOMIC is reported where it is.
// Dropping a constraint belongs to DisallowDropConstraint.
func checkDisallowDropObject(s *statement, r *reporter) {
	report := func(rng review.Range, kind string, names []string) {
		msg := "drops " + kind
		if len(names) > 0 {
			msg += " " + strings.Join(names, ", ")
		}
		r.report(review.DisallowDropObject, s.index, rng, msg)
	}
	ast.Inspect(s.node, func(n ast.Node) bool {
		checkDropNode(n, report)
		return true
	})
}

func checkDropNode(n ast.Node, report func(review.Range, string, []string)) {
	switch v := n.(type) {
	case *ast.DropStmt:
		kind := ast.ObjectType(v.RemoveType)
		var names []string
		if v.Objects != nil {
			for _, obj := range v.Objects.Items {
				if name := objectName(kind, obj); name != "" {
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
		// ALTER TYPE ... DROP ATTRIBUTE is the same node with the type
		// as the relation.
		kind, of := "column", ""
		if ast.ObjectType(v.ObjType) == ast.OBJECT_TYPE {
			kind, of = "attribute", "type "
		}
		for _, item := range v.Cmds.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok || ast.AlterTableType(cmd.Subtype) != ast.AT_DropColumn {
				continue
			}
			report(rangeOf(cmd.Loc), kind, []string{ident(cmd.Name) + " of " + of + relation(v.Relation)})
		}
	}
}

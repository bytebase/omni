package review

import (
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
	report := func(kind string, rng review.Range, msg string) {
		r.report(review.DropCode(kind), s.index, rng, msg)
	}
	ast.Inspect(s.node, func(n ast.Node) bool {
		checkDropNode(n, report)
		return true
	})
}

// checkDropNode reports each drop with the kind word its message names
// the object by, which gives the finding its code (review.DropCode). A
// message that names it by a phrase passes the phrase's kind: cast,
// transform, user mapping, objects owned by.
func checkDropNode(n ast.Node, report func(kind string, rng review.Range, msg string)) {
	switch v := n.(type) {
	case *ast.DropStmt:
		kind := ast.ObjectType(v.RemoveType)
		var targets []string
		if v.Objects != nil {
			for _, obj := range v.Objects.Items {
				if target := dropTarget(kind, obj); target != "" {
					targets = append(targets, target)
				}
			}
		}
		report(objectKind(kind), rangeOf(v.Loc), dropsMessage(kind, targets))
	case *ast.DropdbStmt:
		report("database", rangeOf(v.Loc), drops("database", []string{ident(v.Dbname)}))
	case *ast.DropRoleStmt:
		report("role", rangeOf(v.Loc), drops("role", roleNames(v.Roles)))
	case *ast.DropTableSpaceStmt:
		report("tablespace", rangeOf(v.Loc), drops("tablespace", []string{ident(v.Tablespacename)}))
	case *ast.DropSubscriptionStmt:
		report("subscription", rangeOf(v.Loc), drops("subscription", []string{ident(v.Subname)}))
	case *ast.DropUserMappingStmt:
		report("user mapping", rangeOf(v.Loc), "Drops the user mapping for "+review.Code(roleName(v.User))+" on server "+review.Code(ident(v.Servername))+".")
	case *ast.DropOwnedStmt:
		report("objects owned by", rangeOf(v.Loc), "Drops the objects owned by "+review.Names(roleNames(v.Roles)...)+".")
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
			report(kind, rangeOf(cmd.Loc), "Drops "+kind+" "+review.Code(ident(cmd.Name))+" of "+of+review.Code(relation(v.Relation))+".")
		}
	}
}

// drops writes the message of a DROP of names of one kind: "Drops table
// `t`.", "Drops tables `t` and `p`.", or "Drops the database." when the
// statement names nothing the review can write.
func drops(kind string, names []string) string {
	targets := make([]string, len(names))
	for i, name := range names {
		targets[i] = review.Code(name)
	}
	return dropsTargets(kind, targets)
}

// dropsTargets is drops for targets already written as Markdown.
func dropsTargets(kind string, targets []string) string {
	if len(targets) == 0 {
		return "Drops the " + kind + "."
	}
	return "Drops " + review.Plural(kind, len(targets)) + " " + review.Join(targets...) + "."
}

// dropsMessage writes the message of a DROP statement. A cast and a
// transform have no name, only the types or the language they are for, so
// their targets carry their own words and take no kind word before them.
func dropsMessage(kind ast.ObjectType, targets []string) string {
	if (kind == ast.OBJECT_CAST || kind == ast.OBJECT_TRANSFORM) && len(targets) > 0 {
		return "Drops " + review.Join(targets...) + "."
	}
	return dropsTargets(objectKind(kind), targets)
}

// dropTarget writes one object a DROP statement names as Markdown, only
// the names in code spans. Each object kind has its own list shape: a
// cast is two types ("the cast from `integer` to `text`"), a transform a
// type and a language ("the transform for `integer` language
// `plpgsql`"), an operator class or family an access method then a name
// ("`c` using `btree`"), and a trigger, policy, or rule a relation then a
// name ("`tr` on `t`"). Any other object is its name, in one code span
// with the argument types of a routine or an operator.
func dropTarget(kind ast.ObjectType, obj ast.Node) string {
	switch kind {
	case ast.OBJECT_CAST:
		if types := typeNames(listOf(obj)); len(types) == 2 {
			return "the cast from " + review.Code(types[0]) + " to " + review.Code(types[1])
		}
	case ast.OBJECT_TRANSFORM:
		if l := listOf(obj); l != nil && len(l.Items) == 2 {
			tn, _ := l.Items[0].(*ast.TypeName)
			lang, _ := l.Items[1].(*ast.String)
			if tn != nil && lang != nil {
				return "the transform for " + review.Code(typeName(tn)) + " language " + review.Code(ident(lang.Str))
			}
		}
	case ast.OBJECT_OPCLASS, ast.OBJECT_OPFAMILY:
		if parts := nameParts(listOf(obj)); len(parts) > 1 {
			return review.Code(qualified(parts[1:])) + " using " + review.Code(ident(parts[0]))
		}
	case ast.OBJECT_TRIGGER, ast.OBJECT_POLICY, ast.OBJECT_RULE:
		if parts := nameParts(listOf(obj)); len(parts) > 1 {
			return review.Code(ident(parts[len(parts)-1])) + " on " + review.Code(qualified(parts[:len(parts)-1]))
		}
	}
	if name := objectName(kind, obj); name != "" {
		return review.Code(name)
	}
	return ""
}

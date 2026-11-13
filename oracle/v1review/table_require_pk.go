package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequirePK judges, after reading every statement, which tables
// top-level statements leave without a primary key. A CREATE TABLE sets
// whether its table has one, from a PRIMARY KEY column or table constraint.
// An ALTER TABLE ADD CONSTRAINT ... PRIMARY KEY gives the table one, and
// DROP PRIMARY KEY takes it away; either records the table even if no
// CREATE TABLE did. DROP CONSTRAINT never takes it away, whatever the
// constraint's name: V1 also matched a command whose ColumnName is
// PRIMARY, which the AST never sets for a dropped constraint. A DROP
// statement of any object kind forgets each name it lists.
//
// Each table left without one is a TablePrimaryKeyMissing finding, in the
// order the tables were first recorded; V1 reported them in Go map order.
// Name is the table's qualified name (see qualifiedTable). A finding
// anchors the table's last CREATE TABLE statement or DROP PRIMARY KEY
// command, whichever came later, even across a DROP that forgot the table;
// Statement is that statement's index in stmts. Statements nested in
// CREATE SCHEMA or a PL/SQL block are not judged. A statement may be passed
// bare or wrapped in an *ast.RawStmt. Ported from Bytebase's V1 rule
// TABLE_REQUIRE_PK.
func TableRequirePK(stmts []ast.Node, currentSchema string) []v1review.Finding {
	var names []string
	hasPK := make(map[string]bool)
	anchors := make(map[string]tableAnchor)
	set := func(name string, has bool) {
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
		hasPK[name] = has
	}
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			name := qualifiedTable(n.Name, currentSchema)
			set(name, createTableHasConstraint(n, ast.CONSTRAINT_PRIMARY))
			anchors[name] = tableAnchor{i, n.Loc}
		case *ast.AlterTableStmt:
			if n.Actions == nil {
				continue
			}
			name := qualifiedTable(n.Name, currentSchema)
			for _, item := range n.Actions.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch cmd.Action {
				case ast.AT_ADD_CONSTRAINT:
					if cmd.Constraint != nil && cmd.Constraint.Type == ast.CONSTRAINT_PRIMARY {
						set(name, true)
					}
				case ast.AT_DROP_CONSTRAINT:
					if strings.Contains(strings.ToUpper(cmd.Subtype), "PRIMARY") || strings.EqualFold(cmd.ColumnName, "PRIMARY") {
						set(name, false)
						anchors[name] = tableAnchor{i, cmd.Loc}
					}
				}
			}
		case *ast.DropStmt:
			if n.Names == nil {
				continue
			}
			for _, item := range n.Names.Items {
				if name, ok := item.(*ast.ObjectName); ok {
					delete(hasPK, qualifiedTable(name, currentSchema))
				}
			}
		}
	}
	var findings []v1review.Finding
	for _, name := range names {
		if has, ok := hasPK[name]; ok && !has {
			findings = append(findings, anchors[name].finding(v1review.TablePrimaryKeyMissing, name))
		}
	}
	return findings
}

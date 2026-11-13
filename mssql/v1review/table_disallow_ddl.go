package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowDDL judges the tables top-level CREATE TABLE, ALTER TABLE,
// DROP TABLE and TRUNCATE TABLE statements name; statements nested in a
// routine body or an IF or BEGIN block are not judged. A table is judged as
// written with the qualifiers it was written with, joined by dots and
// without brackets or double quotes (dbo.t, db..t as db.t, [T] as T); a
// table whose name equals an entry of disallow exactly (case-sensitively)
// is TableDisallowed, with Name that joined name. DROP TABLE judges each
// table it lists, in order. Every finding anchors the whole statement.
// Ported from Bytebase's V1 rule TABLE_DISALLOW_DDL.
func TableDisallowDDL(stmts []ast.Node, disallow []string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var refs []*ast.TableRef
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			loc = n.Loc
			if n.Name != nil {
				refs = append(refs, n.Name)
			}
		case *ast.AlterTableStmt:
			loc = n.Loc
			if n.Name != nil {
				refs = append(refs, n.Name)
			}
		case *ast.DropStmt:
			loc = n.Loc
			if n.ObjectType != ast.DropTable || n.Names == nil {
				continue
			}
			for _, item := range n.Names.Items {
				if ref, ok := item.(*ast.TableRef); ok && ref != nil {
					refs = append(refs, ref)
				}
			}
		case *ast.TruncateStmt:
			loc = n.Loc
			if n.Table != nil {
				refs = append(refs, n.Table)
			}
		default:
			continue
		}
		start, end := span(loc)
		for _, ref := range refs {
			if name := joinTableRef(ref, ""); slices.Contains(disallow, name) {
				findings = append(findings, v1review.Finding{Code: v1review.TableDisallowed, Name: name, Start: start, End: end})
			}
		}
	}
	return findings
}

// joinTableRef joins a table reference's non-empty database, schema and
// object names with dots; defaultSchema stands in for a missing schema.
func joinTableRef(ref *ast.TableRef, defaultSchema string) string {
	if ref == nil {
		return ""
	}
	schema := ref.Schema
	if schema == "" {
		schema = defaultSchema
	}
	var parts []string
	for _, part := range []string{ref.Database, schema, ref.Object} {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

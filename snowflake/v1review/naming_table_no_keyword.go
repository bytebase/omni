package v1review

import (
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTableNoKeyword judges the table names top-level statements
// introduce against the V1 keyword table (v1Keywords): the table a plain
// CREATE TABLE creates, and the new name an ALTER TABLE ... RENAME TO
// gives. CREATE TABLE ... AS SELECT, LIKE and CLONE are not judged, nor
// are statements inside a Snowflake Scripting block, a task's body, or
// under EXPLAIN.
//
// The name is judged without its database or schema qualifier, as
// Snowflake resolves it and then upper-cased, so a quoted name matches
// whatever its case: "connect" is reported though Snowflake keeps it
// distinct from the keyword. Finding.Name is the resolved name, before
// that upper-casing: an unquoted name upper-cased, a quoted name as its
// content. An IDENTIFIER(...) name is judged as written, keyword and
// parentheses included, so it never matches.
//
// A finding has no range: the V1 advice carries no position.
// NameIsKeyword is the only code. Ported from Bytebase's V1 rule
// NAMING_TABLE_NO_KEYWORD.
func NamingTableNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	check := func(ident ast.Ident) {
		name := ident.Normalize()
		if isV1Keyword(name) {
			findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: name})
		}
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.AsSelect != nil || n.Like != nil || n.Clone != nil || n.Name == nil {
				continue
			}
			check(n.Name.Name)
		case *ast.AlterTableStmt:
			for _, a := range n.Actions {
				if a == nil || a.Kind != ast.AlterTableRename || a.NewName == nil {
					continue
				}
				check(a.NewName.Name)
			}
		}
	}
	return findings
}

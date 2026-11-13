package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCollationAllowlist judges the collations top-level CREATE and ALTER
// DATABASE and CREATE and ALTER TABLE statements name against the allowed
// list, compared without regard to letter case. Statements nested in a
// routine body are not judged.
//
// A CREATE or ALTER DATABASE has only its first COLLATE option judged. A
// CREATE TABLE reports at most one collation: its first COLLATE table option
// when that is not allowed, otherwise the first column whose type names one
// that is not. An ALTER TABLE reports at most one too, the last not allowed
// among its COLLATE table options, the COLLATE of its CONVERT TO CHARACTER
// SET, and the column types of its MODIFY and CHANGE COLUMN actions and of
// the first such column each ADD COLUMN gives.
//
// Name is the collation lower-cased. A finding from a CREATE TABLE column
// anchors that column's definition; every other finding anchors the whole
// statement. Ported from Bytebase's V1 rule SYSTEM_COLLATION_ALLOWLIST.
func SystemCollationAllowlist(stmts []ast.Node, allowed []string) []v1review.Finding {
	return judgeAllowlist(stmts, allowed, allowlistRule{
		code:    v1review.CollationNotAllowed,
		options: []string{"COLLATE"},
		column: func(col *ast.ColumnDef) string {
			if col == nil || col.TypeName == nil {
				return ""
			}
			return col.TypeName.Collate
		},
		convert: func(cmd *ast.AlterTableCmd) string { return cmd.NewName },
	})
}

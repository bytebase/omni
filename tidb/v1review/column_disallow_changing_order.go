package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowChangingOrder reports each top-level ALTER TABLE statement
// holding a CHANGE or MODIFY COLUMN action with FIRST or AFTER, once however
// many it holds, anchored at the whole statement. ADD COLUMN with FIRST or
// AFTER is not judged. Ported from Bytebase's V1 rule
// COLUMN_DISALLOW_CHANGING_ORDER.
func ColumnDisallowChangingOrder(stmts []ast.Node) []v1review.Finding {
	return alterRule(stmts, v1review.ColumnOrderChanged, func(cmd *ast.AlterTableCmd) bool {
		return (cmd.Type == ast.ATChangeColumn || cmd.Type == ast.ATModifyColumn) && (cmd.First || cmd.After != "")
	})
}

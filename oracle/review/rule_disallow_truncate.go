package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowTruncate reports TRUNCATE TABLE and TRUNCATE CLUSTER, and
// ALTER TABLE ... TRUNCATE PARTITION or SUBPARTITION, which anchors its
// clause. A TRUNCATE runs inside PL/SQL only as the text of an EXECUTE
// IMMEDIATE, which is not read.
func checkDisallowTruncate(sql string, s *statement, r *reporter) {
	switch v := s.node.(type) {
	case *ast.TruncateStmt:
		what := "table "
		if v.Cluster {
			what = "cluster "
		}
		r.report(review.DisallowTruncate, s.index, rangeOf(v.Loc), "truncates "+what+objectName(v.Table))
	case *ast.AlterTableStmt:
		if v.Actions == nil {
			return
		}
		for _, item := range v.Actions.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok || cmd.Action != ast.AT_TRUNCATE_PARTITION {
				continue
			}
			kind := "partition"
			if strings.EqualFold(cmd.Subtype, "SUBPARTITION") {
				kind = "subpartition"
			}
			r.report(review.DisallowTruncate, s.index, rangeOf(cmd.Loc), "truncates "+partitionWhat(kind, sql, cmd)+" of "+objectName(v.Name))
		}
	}
}

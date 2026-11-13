package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// liveTarget returns the statement a live-connection rule judges: the
// statement itself, or the one an EXPLAIN ANALYZE runs, as Bytebase's
// UnwrapExplainAnalyze picks it. An EXPLAIN is an EXPLAIN ANALYZE when its
// last ANALYZE option is not false, off or 0.
func liveTarget(stmt ast.Node) ast.Node {
	stmt = unwrapRaw(stmt)
	explain, ok := stmt.(*ast.ExplainStmt)
	if !ok || explain.Options == nil {
		return stmt
	}
	analyze := false
	for _, item := range explain.Options.Items {
		option, ok := item.(*ast.DefElem)
		if !ok || !strings.EqualFold(option.Defname, "analyze") {
			continue
		}
		switch arg := option.Arg.(type) {
		case *ast.String:
			analyze = !strings.EqualFold(arg.Str, "false") && !strings.EqualFold(arg.Str, "off")
		case *ast.Integer:
			analyze = arg.Ival != 0
		default:
			analyze = true
		}
	}
	if !analyze {
		return stmt
	}
	return explain.Query
}

// hasDataModifyingCTE reports whether a WITH clause holds an INSERT, UPDATE,
// DELETE or MERGE.
func hasDataModifyingCTE(with *ast.WithClause) bool {
	if with == nil || with.Ctes == nil {
		return false
	}
	for _, item := range with.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		switch cte.Ctequery.(type) {
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			return true
		}
	}
	return false
}

// rowLimitFinding judges one statement's row estimate against maxRows.
func rowLimitFinding(probe v1review.Probe, maxRows, statement int) (v1review.Finding, bool) {
	switch {
	case probe.Error != "":
		return v1review.Finding{Code: v1review.DryRunFailed, Detail: probe.Error, Statement: statement}, true
	case probe.CountError != "":
		return v1review.Finding{Code: v1review.RowCountFailed, Detail: probe.CountError, Statement: statement}, true
	case probe.Rows > int64(maxRows):
		return v1review.Finding{Code: v1review.RowLimitExceeded, Count: int(probe.Rows), Statement: statement}, true
	}
	return v1review.Finding{}, false
}

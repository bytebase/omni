package analysis

import (
	"errors"

	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
)

// ErrExplainAnalyzeExecutes is ExplainStatement's error for EXPLAIN ANALYZE,
// which plans by running what it plans.
var ErrExplainAnalyzeExecutes = errors.New("EXPLAIN ANALYZE executes the statement it plans, which an explain request never does: run it as a query, or drop ANALYZE to plan the statement")

// ExplainStatement returns the statement that plans statement. Redshift returns
// one plan format here whatever the caller asked for: its EXPLAIN has no
// machine-readable form.
//
// An EXPLAIN is rebuilt around the statement it plans, so EXPLAIN VERBOSE
// SELECT 1 stays one plan rather than a plan of a plan.
func ExplainStatement(statement string) (string, error) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", err
	}
	planned := statement
	if explain, ok := node.(*ast.ExplainStmt); ok {
		if hasExplainAnalyze(explain) {
			return "", ErrExplainAnalyzeExecutes
		}
		if inner, ok := explainedText(explain, statement); ok {
			planned = inner
		} else {
			// Nothing to rebuild around: run what the user wrote rather than
			// plan a plan.
			return statement, nil
		}
	}
	return "EXPLAIN " + planned, nil
}

// singleStatement returns the one statement in text. An explain request plans one
// statement at a time: both the drivers and the API gate split a multi-statement
// request first.
func singleStatement(text string) (ast.Node, error) {
	stmts, err := redshift.Parse(text)
	if err != nil {
		return nil, err
	}
	var node ast.Node
	for _, stmt := range stmts {
		if stmt.AST == nil {
			continue
		}
		if node != nil {
			return nil, errors.New("EXPLAIN plans a single statement")
		}
		node = stmt.AST
	}
	if node == nil {
		return nil, errors.New("EXPLAIN needs a statement to plan")
	}
	return node, nil
}

// explainedText returns the text within text of the statement explain plans, and
// whether it could be located.
func explainedText(explain *ast.ExplainStmt, text string) (string, bool) {
	// The statement runs from its WITH clause, which the location of a SELECT
	// leaves out, to the end of the EXPLAIN, which its location can also leave
	// out, as for ORDER BY.
	start := ast.NodeLoc(explain.Query).Start
	if with := explainedWithClause(explain.Query); with != nil && with.Loc.Start >= 0 && with.Loc.Start < start {
		start = with.Loc.Start
	}
	if start < 0 || start >= explain.Loc.End || explain.Loc.End > len(text) {
		return "", false
	}
	return text[start:explain.Loc.End], true
}

// explainedWithClause returns the WITH clause of a statement EXPLAIN can plan.
func explainedWithClause(node ast.Node) *ast.WithClause {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return n.WithClause
	case *ast.InsertStmt:
		return n.WithClause
	case *ast.UpdateStmt:
		return n.WithClause
	case *ast.DeleteStmt:
		return n.WithClause
	default:
		return nil
	}
}

package analysis

import (
	"errors"

	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
)

// ExplainStatement returns the statement that plans statement. StarRocks returns
// one plan format here whatever the caller asked for: its other EXPLAIN forms are
// plans of their own shape that nothing downstream reads (the driver capability
// table offers this engine the default plan only).
//
// An EXPLAIN is rebuilt around the statement it plans, so EXPLAIN SELECT 1
// stays one plan rather than a plan of a plan.
func ExplainStatement(statement string) (string, error) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", err
	}
	planned := statement
	switch n := node.(type) {
	case *ast.ExplainStmt:
		if inner, ok := wrappedText(n.Loc, n.Query, statement); ok {
			planned = inner
		} else {
			// Nothing to rebuild around: run what the user wrote rather than
			// plan a plan.
			return statement, nil
		}
	case *ast.OptimizerTraceStmt:
		// TRACE is already a plan, so plan the query it wraps.
		if inner, ok := wrappedText(n.Loc, n.Query, statement); ok {
			planned = inner
		} else {
			return statement, nil
		}
	}
	return "EXPLAIN " + planned, nil
}

// DescribeExplain returns the plan format statement asks for and whether
// producing it runs the statement it plans (EXPLAIN ANALYZE). It is not an
// explain at all unless it plans a statement: EXPLAIN tbl_name parses here too
// and describes a table's columns, which is a result like any other.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return "", false, false
	}
	switch node := file.Stmts[0].(type) {
	case *ast.ExplainStmt:
		if node.Query == nil {
			return "", false, false
		}
		format = node.Type
		if format == "" {
			format = "TEXT"
		}
		return format, node.Type == "ANALYZE", true
	case *ast.OptimizerTraceStmt:
		if node.Query == nil {
			return "", false, false
		}
		// TRACE is a plan too, and it never executes the traced query.
		return "TRACE", false, true
	}
	return "", false, false
}

// singleStatement returns the one statement in text. An explain request plans one
// statement at a time: both the drivers and the API gate split a multi-statement
// request first.
func singleStatement(text string) (ast.Node, error) {
	file, err := parser.Parse(text)
	if err != nil {
		// The engine's entry points report the first parse error, not the
		// list's summary.
		return nil, parser.FirstError(err)
	}
	if file == nil || len(file.Stmts) == 0 || file.Stmts[0] == nil {
		return nil, errors.New("EXPLAIN needs a statement to plan")
	}
	if len(file.Stmts) > 1 {
		return nil, errors.New("EXPLAIN plans a single statement")
	}
	return file.Stmts[0], nil
}

// wrappedText returns the text within text of the statement a plan wrapper
// plans, and whether it could be located.
func wrappedText(wrapper ast.Loc, child ast.Node, text string) (string, bool) {
	inner := ast.NodeLoc(child)
	// The plan ends where the wrapper does, which is past clauses that the
	// planned statement's own location can leave out.
	if inner.Start < 0 || inner.Start >= wrapper.End || wrapper.End > len(text) {
		return "", false
	}
	return text[inner.Start:wrapper.End], true
}

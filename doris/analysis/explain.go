package analysis

import (
	"errors"

	"github.com/bytebase/omni/doris/ast"
	"github.com/bytebase/omni/doris/parser"
)

// ExplainStatement returns the statement that plans statement. Doris returns
// one plan format here whatever the caller asked for: its other EXPLAIN forms
// are plans of their own shape that nothing downstream reads (the driver
// capability table offers this engine the default plan only).
//
// An EXPLAIN is rebuilt around the statement it plans, so EXPLAIN SELECT 1
// stays one plan rather than a plan of a plan.
func ExplainStatement(statement string) (string, error) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", err
	}
	planned := statement
	if explain, ok := node.(*ast.ExplainStmt); ok {
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

// DescribeExplain returns the plan format statement asks for and whether
// producing it runs the statement it plans (EXPLAIN ANALYZE). It is not an
// explain at all unless it plans a statement: EXPLAIN tbl_name parses here too
// and describes a table's columns, which is a result like any other.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return "", false, false
	}
	explain, ok := file.Stmts[0].(*ast.ExplainStmt)
	if !ok || explain.Query == nil {
		return "", false, false
	}
	format = explain.Type
	if format == "" {
		format = "TEXT"
	}
	return format, explain.Type == "ANALYZE", true
}

// singleStatement returns the one statement in text. An explain request plans one
// statement at a time: both the drivers and the API gate split a multi-statement
// request first.
func singleStatement(text string) (ast.Node, error) {
	file, err := parser.Parse(text)
	if err != nil {
		return nil, err
	}
	if file == nil || len(file.Stmts) == 0 || file.Stmts[0] == nil {
		return nil, errors.New("EXPLAIN needs a statement to plan")
	}
	if len(file.Stmts) > 1 {
		return nil, errors.New("EXPLAIN plans a single statement")
	}
	return file.Stmts[0], nil
}

// explainedText returns the text within text of the statement explain plans, and
// whether it could be located.
func explainedText(explain *ast.ExplainStmt, text string) (string, bool) {
	inner := ast.NodeLoc(explain.Query)
	// The plan ends where the EXPLAIN does, which is past clauses that the
	// planned statement's own location can leave out.
	if inner.Start < 0 || inner.Start >= explain.Loc.End || explain.Loc.End > len(text) {
		return "", false
	}
	return text[inner.Start:explain.Loc.End], true
}

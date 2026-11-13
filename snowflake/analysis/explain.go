package analysis

import (
	"errors"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// ExplainStatement returns the statement that plans statement: EXPLAIN USING
// JSON when json is set, otherwise EXPLAIN, which returns a table. An EXPLAIN
// is rebuilt around the statement it plans; one whose planned statement cannot
// be located is returned as written.
func ExplainStatement(statement string, json bool) (string, error) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", err
	}
	planned := statement
	if explain, ok := node.(*ast.ExplainStmt); ok {
		inner, ok := explainedText(explain, statement)
		if !ok {
			return statement, nil
		}
		planned = inner
	}
	if json {
		return "EXPLAIN USING JSON " + planned, nil
	}
	return "EXPLAIN " + planned, nil
}

// DescribeExplain returns the output format of an EXPLAIN statement, "json"
// or "text", and whether it executes the statement it plans. Snowflake's
// EXPLAIN compiles a statement without running it, so executed is always
// false. ok is false when statement is not one EXPLAIN.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", false, false
	}
	explain, ok := node.(*ast.ExplainStmt)
	if !ok {
		return "", false, false
	}
	if explain.Format == ast.ExplainJSON {
		return "json", false, true
	}
	// The tabular and text formats are both tables of operations.
	return "text", false, true
}

var errNoStatement = errors.New("EXPLAIN needs a statement to plan")

// singleStatement returns the one statement in text: an explain request plans
// one statement at a time.
func singleStatement(text string) (ast.Node, error) {
	var stmts []Statement
	for _, stmt := range SplitStatements(text) {
		if !stmt.Empty {
			stmts = append(stmts, stmt)
		}
	}
	if len(stmts) == 0 {
		return nil, errNoStatement
	}
	if len(stmts) > 1 {
		return nil, errors.New("EXPLAIN plans a single statement")
	}
	file, err := parser.Parse(stmts[0].Text)
	if err != nil {
		return nil, err
	}
	if file == nil || len(file.Stmts) == 0 {
		return nil, errNoStatement
	}
	return file.Stmts[0], nil
}

// explainedText returns the text of the statement explain plans, and whether
// it could be located.
func explainedText(explain *ast.ExplainStmt, text string) (string, bool) {
	inner := ast.NodeLoc(explain.Stmt)
	// The plan ends where the EXPLAIN does, past clauses the planned
	// statement's own location can leave out.
	if inner.Start < 0 || inner.Start >= explain.Loc.End || explain.Loc.End > len(text) {
		return "", false
	}
	return text[inner.Start:explain.Loc.End], true
}

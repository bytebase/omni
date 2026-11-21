package analysis

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/analysis"
	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// ErrExplainAnalyzeExecutes is the explain builders' error for EXPLAIN
// ANALYZE, which runs the statement it plans.
var ErrExplainAnalyzeExecutes = analysis.ErrExplainAnalyzeExecutes

// plannedStatement returns the statement an EXPLAIN should be built around, and
// whether one can be: a statement that already asks for a plan is rebuilt around
// the statement it plans, and one whose planned statement cannot be located runs
// as the caller wrote it rather than becoming the plan of a plan.
func plannedStatement(statement string) (planned string, rebuild bool, err error) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", false, err
	}
	explain, ok := node.(*ast.ExplainStmt)
	if !ok {
		return statement, true, nil
	}
	if isExplainAnalyzeStmt(explain) {
		return "", false, ErrExplainAnalyzeExecutes
	}
	inner, ok := explainedText(explain, statement)
	if !ok {
		return statement, false, nil
	}
	return inner, true, nil
}

// ExplainStatement returns the statement that plans statement in format:
// "json", "xml" or "yaml", or the engine's default text plan for any other
// value. An EXPLAIN is rebuilt around the statement it plans, and one whose
// planned statement cannot be located is returned as written. EXPLAIN ANALYZE
// returns ErrExplainAnalyzeExecutes, and a statement that does not parse
// returns the parser's error.
func ExplainStatement(statement string, format string) (string, error) {
	planned, rebuild, err := plannedStatement(statement)
	if err != nil {
		return "", err
	}
	if !rebuild {
		return planned, nil
	}
	switch format {
	case "json":
		return "EXPLAIN (FORMAT JSON) " + planned, nil
	case "xml":
		return "EXPLAIN (FORMAT XML) " + planned, nil
	case "yaml":
		return "EXPLAIN (FORMAT YAML) " + planned, nil
	default:
		return "EXPLAIN " + planned, nil
	}
}

// DescribeExplain returns the output format of an EXPLAIN statement and
// whether it executes the statement it plans.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	node, err := singleStatement(statement)
	if err != nil {
		return "", false, false
	}
	explain, ok := node.(*ast.ExplainStmt)
	if !ok {
		return "", false, false
	}
	return explainFormat(explain), isExplainAnalyzeStmt(explain), true
}

func explainFormat(explain *ast.ExplainStmt) string {
	format := "text"
	if explain.Options == nil {
		return format
	}
	for _, item := range explain.Options.Items {
		option, ok := item.(*ast.DefElem)
		if !ok || !strings.EqualFold(option.Defname, "format") {
			continue
		}
		format = ""
		if arg, ok := option.Arg.(*ast.String); ok {
			format = strings.ToLower(arg.Str)
		}
	}
	return format
}

// CockroachExplainStatement is ExplainStatement for CockroachDB, which Bytebase
// parses with the PostgreSQL grammar but whose EXPLAIN takes its own option list rather than PostgreSQL's FORMAT and
// has no option that returns a plan as a document.
//
// What it does have is DISTSQL, which returns the text plan as always and a
// link to a diagram of the physical flow beside it; the payload in that link is
// the one machine-readable plan CockroachDB publishes, and the driver decodes
// it. The statement does not execute, as EXPLAIN ANALYZE (DISTSQL) would.
func CockroachExplainStatement(statement string, json bool) (string, error) {
	planned, rebuild, err := plannedStatement(statement)
	if err != nil {
		return "", err
	}
	if !rebuild {
		return planned, nil
	}
	if json {
		return "EXPLAIN (DISTSQL) " + planned, nil
	}
	return "EXPLAIN " + planned, nil
}

// singleStatement returns the one statement in text. An explain request plans one
// statement at a time: both the drivers and the API gate split a multi-statement
// request first.
func singleStatement(text string) (ast.Node, error) {
	stmts, err := omnipg.Parse(text)
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

package analysis

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	"github.com/bytebase/omni/analysis"
)

// ErrExplainAnalyzeExecutes is ExplainStatement's error for EXPLAIN ANALYZE,
// matching pg/analysis and the other analysis packages.
var ErrExplainAnalyzeExecutes = analysis.ErrExplainAnalyzeExecutes

// ExplainStatement builds the EXPLAIN statement for statement in CockroachDB's
// grammar, unwrapping an existing EXPLAIN so a plan never plans another plan.
// EXPLAIN ANALYZE returns ErrExplainAnalyzeExecutes.
func ExplainStatement(statement string, json bool) (string, error) {
	parsed, err := Parse(statement)
	if err != nil {
		return "", err
	}
	if len(parsed) != 1 {
		return "", fmt.Errorf("EXPLAIN plans a single statement")
	}
	planned := statement
	switch node := parsed[0].AST.(type) {
	case *tree.ExplainAnalyze:
		return "", ErrExplainAnalyzeExecutes
	case *tree.Explain:
		planned = node.Statement.String()
	}
	if json {
		return "EXPLAIN (DISTSQL) " + planned, nil
	}
	return "EXPLAIN " + planned, nil
}

// DescribeExplain reports a directly submitted EXPLAIN statement's format
// ("JSON" or "TEXT") and whether it executes the statement it plans
// (EXPLAIN ANALYZE). ok is false when statement is not a single EXPLAIN.
func DescribeExplain(statement string) (format string, executed bool, ok bool) {
	parsed, err := Parse(statement)
	if err != nil || len(parsed) != 1 {
		return "", false, false
	}
	switch node := parsed[0].AST.(type) {
	case *tree.Explain:
		return explainFormat(node.ExplainOptions), false, true
	case *tree.ExplainAnalyze:
		return explainFormat(node.ExplainOptions), true, true
	}
	return "", false, false
}

func explainFormat(opts tree.ExplainOptions) string {
	if opts.Flags[tree.ExplainFlagJSON] {
		return "JSON"
	}
	return "TEXT"
}

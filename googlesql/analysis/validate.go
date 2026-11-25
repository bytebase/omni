package analysis

import (
	"github.com/bytebase/omni/googlesql/ast"
	"github.com/bytebase/omni/googlesql/parser"
)

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns rows:
// (readOnly, returnsData, error).
//
// Decision is AST-based and parsed with the requested dialect's strict grammar
// rules (DialectBigQuery or DialectSpanner). Every non-empty top-level
// statement in a script must be read-only (SELECT/WITH/TABLE/set-op queries,
// read-only pipe queries, GQL graph reads, SHOW, DESCRIBE, or EXPLAIN over a
// read-only statement). Locking reads (FOR UPDATE), writing or exporting pipe
// operators (|> INSERT, |> CREATE TABLE, |> EXPORT DATA), EXPLAIN ANALYZE,
// EXPORT/LOAD/CLONE, DML, DDL, DCL, SET, transactions, CALL, EXECUTE IMMEDIATE,
// and procedural scripting blocks are rejected. Empty, whitespace-only, or
// comment-only input carries no statements to refuse and returns
// (true, true, nil). Syntax and lex errors are returned so callers fail closed.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	file, err := parser.Parse(statement)
	if err != nil {
		return false, false, err
	}
	if file == nil {
		return true, true, nil
	}
	for _, node := range file.Stmts {
		if !isReadOnlyNode(node) {
			return false, false, nil
		}
	}
	return true, true, nil
}

func isReadOnlyNode(node ast.Node) bool {
	if node == nil {
		return false
	}
	switch n := node.(type) {
	case *ast.QueryStmt, *ast.SelectStmt, *ast.SetOperation, *ast.GQLStmt:
		return isReadOnlyQuery(n)
	case *ast.SubpipelineStmt:
		return n.Pipeline != nil && isReadOnlyQuery(n)
	case *ast.PipeStatement:
		return n.Head != nil && isReadOnlyNode(n.Head) && n.Pipeline != nil && isReadOnlyQuery(n.Pipeline)
	case *ast.ShowStmt:
		return n.Target != ""
	case *ast.DescribeStmt:
		return n.Path != nil
	case *ast.ExplainStmt:
		if n.Stmt == nil {
			return false
		}
		if _, nested := n.Stmt.(*ast.ExplainStmt); nested {
			return false
		}
		return isReadOnlyNode(n.Stmt)
	default:
		return false
	}
}

// isReadOnlyQuery walks a query or pipeline subtree and rejects locking reads
// (FOR UPDATE), writing/exporting pipe operators, or any embedded write node.
func isReadOnlyQuery(node ast.Node) bool {
	if node == nil {
		return false
	}
	readOnly := true
	ast.Inspect(node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.QueryStmt:
			if v.ForUpdate {
				readOnly = false
				return false
			}
		case *ast.PipeInsert, *ast.PipeCreateTable, *ast.PipeExportData:
			readOnly = false
			return false
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt, *ast.TruncateStmt,
			*ast.ExportDataStmt, *ast.ExportModelStmt, *ast.ExportMetadataStmt,
			*ast.LoadDataStmt, *ast.CloneDataStmt,
			*ast.CallStmt, *ast.ExecuteImmediateStmt:
			readOnly = false
			return false
		}
		return readOnly
	})
	return readOnly
}

package analysis

import (
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// QueryTypeOf returns the QueryType of one parsed Trino statement and whether
// it is an EXPLAIN ANALYZE.
//
// EXPLAIN ANALYZE actually executes the inner query, so its read-only-ness is
// the inner statement's: the returned type is the inner type and analyze is
// true. A plain EXPLAIN is Explain and never executes. A caller that gates
// writes on read-only-ness must refuse an analyzed non-SELECT.
//
// A nil node yields (Unknown, false).
func QueryTypeOf(node ast.Node) (QueryType, bool) {
	if node == nil {
		return Unknown, false
	}
	switch n := node.(type) {
	case *parser.QueryStmt, *parser.WithFunctionStmt:
		return Select, false
	case *parser.ExplainStmt:
		if n.Analyze {
			// EXPLAIN ANALYZE EXECUTES the inner statement (oracle-confirmed:
			// Trino 481 runs it), so its read-only-ness is the inner statement's,
			// not unconditionally read-only. Report the inner type and flag
			// analyze so ValidateQuery rejects EXPLAIN ANALYZE of a non-SELECT
			// (e.g. EXPLAIN ANALYZE UPDATE, which would otherwise run a write
			// through the read-only SQL-editor gate).
			innerType, _ := QueryTypeOf(n.Statement)
			return innerType, true
		}
		return Explain, false
	case *parser.ShowStmt,
		*parser.DescribeInputStmt, *parser.DescribeOutputStmt:
		return SelectInfoSchema, false
	case *parser.SetSessionStmt, *parser.ResetSessionStmt,
		*parser.SetSessionAuthorizationStmt, *parser.ResetSessionAuthorizationStmt,
		*parser.SetPathStmt, *parser.SetRoleStmt, *parser.SetTimeZoneStmt,
		*parser.UseStmt:
		// Session-state statements were classified read-only (Select) by the
		// legacy listener (EnterSetSession / EnterResetSession). USE is also a
		// session-scoped no-op for read purposes. Keep them read-only.
		return Select, false
	case *parser.InsertStmt, *parser.UpdateStmt, *parser.DeleteStmt,
		*parser.MergeStmt, *parser.TruncateStmt, *parser.CallStmt:
		return DML, false
	}
	// Everything else (CREATE/ALTER/DROP/COMMENT/GRANT/REVOKE/ANALYZE/
	// transaction control/prepared-statement admin/...) is DDL. This is the safe
	// upper bound for ACL purposes and matches the legacy listener's
	// "everything else is DDL" default.
	return DDL, false
}

// queryTypeFromText parses statement and classifies its single top-level
// statement, applying the legacy statement-text system-schema promotion: a
// Select that references a Trino system/metadata schema becomes
// SelectInfoSchema. EXPLAIN over such a query is likewise promoted.
//
// On parse failure or empty/multi input it returns Unknown.
func queryTypeFromText(statement string) QueryType {
	file, err := parser.Parse(statement)
	if err != nil || file == nil || len(file.Stmts) == 0 {
		return Unknown
	}
	qt, _ := QueryTypeOf(file.Stmts[0])
	switch qt {
	case Select, Explain:
		if hasSystemSchema(statement) {
			return SelectInfoSchema
		}
	default:
	}
	return qt
}

// IsReadOnlyStatement reports whether the given AST node is a read-only Trino
// statement (SELECT family, EXPLAIN, SHOW/DESCRIBE, session-state).
func IsReadOnlyStatement(node ast.Node) bool {
	queryType, _ := QueryTypeOf(node)
	return queryType == Select ||
		queryType == Explain ||
		queryType == SelectInfoSchema
}

// IsDataChangingStatement reports whether the node is a DML statement.
func IsDataChangingStatement(node ast.Node) bool {
	queryType, _ := QueryTypeOf(node)
	return queryType == DML
}

// IsSchemaChangingStatement reports whether the node is a DDL statement.
func IsSchemaChangingStatement(node ast.Node) bool {
	queryType, _ := QueryTypeOf(node)
	return queryType == DDL
}

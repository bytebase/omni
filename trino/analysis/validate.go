package analysis

import (
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns data:
// (readOnly, returnsData, error).
//
// EXPLAIN ANALYZE is special-cased: because it executes the inner query, it is
// only accepted when QueryTypeOf reports the inner statement as a read-only
// SELECT; it then counts as read-only and data-returning. A session statement
// (USE, SET/RESET SESSION and the authorization forms) is read-only but returns
// no data, because it rebinds the connection rather than reading from it.
func ValidateQuery(statement string) (bool, bool, error) {
	file, err := parser.Parse(statement)
	if err != nil {
		return false, false, err
	}

	allReadOnly := true
	allReturnData := true

	if file != nil {
		for _, node := range file.Stmts {
			queryType, isAnalyze := QueryTypeOf(node)

			if isAnalyze {
				// EXPLAIN ANALYZE runs the query. Only allow it when the type is a
				// read-only SELECT; in that case it is read-only and returns data.
				if queryType != Select {
					return false, false, nil
				}
				continue
			}

			readOnly := queryType == Select ||
				queryType == Explain ||
				queryType == SelectInfoSchema
			// USE and the SET family rebind the connection's catalog, schema, role
			// or identity, so a later statement resolves somewhere the caller did
			// not name. Reported through the same bool as every other engine's SET.
			returnsData := readOnly && !changesSession(node)

			if !readOnly {
				allReadOnly = false
			}
			if !returnsData {
				allReturnData = false
			}
			if !allReadOnly {
				break
			}
		}
	}

	return allReadOnly, allReturnData, nil
}

// changesSession reports whether a statement rebinds the connection rather than
// reading from it.
func changesSession(node ast.Node) bool {
	switch node.(type) {
	case *parser.SetSessionStmt, *parser.ResetSessionStmt,
		*parser.SetSessionAuthorizationStmt, *parser.ResetSessionAuthorizationStmt,
		*parser.SetPathStmt, *parser.SetRoleStmt, *parser.SetTimeZoneStmt,
		*parser.UseStmt:
		return true
	default:
		return false
	}
}

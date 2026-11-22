package analysis

import (
	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/cassandra/ast"
)

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns rows:
// (readOnly, returnsData, error).
//
// Decision is AST-based: every non-empty statement in statement must parse as a
// CQL SELECT statement. DML (INSERT, UPDATE, DELETE, BATCH, TRUNCATE), DDL, DCL
// (GRANT, REVOKE, role/user/identity management), USE, DESCRIBE, and LIST are
// rejected. Empty, whitespace-only, comment-only, or semicolon-only input
// carries no statements to refuse and returns (true, true, nil). Syntax and lex
// errors are returned so callers fail closed before execution.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	stmts, err := cassandra.Parse(statement)
	if err != nil {
		return false, false, err
	}
	for _, stmt := range stmts {
		if stmt.Empty() {
			continue
		}
		if _, ok := stmt.AST.(*ast.SelectStmt); !ok {
			return false, false, nil
		}
	}
	return true, true, nil
}

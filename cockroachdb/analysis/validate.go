package analysis

import (
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// ValidateQuery reports whether a script may run in the SQL editor's
// read-only mode, and whether every statement returns data rather than set
// session state. A script that does not parse returns the parser's error.
//
// SELECT, VALUES and TABLE are allowed unless a CTE, statement source or
// subquery writes; EXPLAIN is allowed, and EXPLAIN ANALYZE only of such a
// read; SHOW is allowed; SET and RESET are allowed but return no data.
// Everything else, DDL and DML included, is refused.
func ValidateQuery(script string) (readOnly bool, returnsData bool, err error) {
	parsed, err := Parse(script)
	if err != nil {
		return false, false, err
	}
	var setsState bool
	for _, stmt := range parsed {
		switch n := stmt.AST.(type) {
		case *tree.Select, *tree.ParenSelect:
			if !isRead(n) {
				return false, false, nil
			}
		case *tree.Explain:
		case *tree.ExplainAnalyze:
			// EXPLAIN ANALYZE runs the statement it explains.
			switch n.Statement.(type) {
			case *tree.Select, *tree.ParenSelect:
				if !isRead(n.Statement) {
					return false, false, nil
				}
			default:
				return false, false, nil
			}
		case *tree.SetVar, *tree.SetTransaction, *tree.SetSessionCharacteristics, *tree.SetSessionAuthorizationDefault, *tree.SetTracing:
			setsState = true
		default:
			if !strings.HasPrefix(n.StatementTag(), "SHOW") {
				return false, false, nil
			}
		}
	}
	return true, !setsState, nil
}

func isRead(stmt tree.Statement) bool {
	return len(getMutationTypes(stmt)) == 0
}

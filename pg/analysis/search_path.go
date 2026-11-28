package analysis

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
)

func SearchPathForSelectedSchema(selectedSchema string, defaultSearchPath []string) []string {
	if selectedSchema == "" {
		return defaultSearchPath
	}
	if selectedSchema == "public" || strings.ContainsRune(selectedSchema, '\x00') {
		return []string{selectedSchema}
	}
	return []string{selectedSchema, "public"}
}

// searchPathPrefix starts the statement WithSearchPath puts before a DML
// statement.
const searchPathPrefix = "SET LOCAL search_path TO "

// WithSearchPath returns statement preceded by a SET LOCAL search_path
// statement for schemas.
func WithSearchPath(statement string, schemas []string) string {
	quoted := make([]string, len(schemas))
	for i, schema := range schemas {
		quoted[i] = `"` + strings.ReplaceAll(schema, `"`, `""`) + `"`
	}
	return searchPathPrefix + strings.Join(quoted, ", ") + ";\n" + statement
}

// SplitSearchPath returns the SET LOCAL search_path statement that
// WithSearchPath put before statement, or "" when there is none, and the
// statement after it.
func SplitSearchPath(statement string) (string, string) {
	if !strings.HasPrefix(statement, searchPathPrefix) {
		return "", statement
	}
	setup, rest, ok := strings.Cut(statement, ";\n")
	if !ok {
		return "", statement
	}
	return setup, rest
}

// searchPathTracker replays the search path changes of a script's statements
// in order: SET, RESET, DISCARD ALL and the transactions that scope SET LOCAL.
type searchPathTracker struct {
	initial []string
	// current is the search path the next statement runs under.
	current []string
	// session leaves out SET LOCAL, which lasts until the transaction ends, and
	// transaction is the session search path that ROLLBACK restores.
	session, transaction []string
	inTransaction        bool
}

func newSearchPathTracker(initial []string) *searchPathTracker {
	return &searchPathTracker{initial: initial, current: initial, session: initial, transaction: initial}
}

// apply records the search path change node makes, if any.
func (t *searchPathTracker) apply(node ast.Node) {
	switch n := node.(type) {
	case *ast.TransactionStmt:
		switch n.Kind {
		case ast.TRANS_STMT_BEGIN, ast.TRANS_STMT_START:
			// A BEGIN inside a transaction only warns.
			if !t.inTransaction {
				t.transaction, t.inTransaction = t.session, true
			}
		case ast.TRANS_STMT_COMMIT, ast.TRANS_STMT_ROLLBACK:
			// A ROLLBACK outside a transaction only warns.
			if n.Kind == ast.TRANS_STMT_ROLLBACK && t.inTransaction {
				t.session = t.transaction
			}
			// COMMIT AND CHAIN and ROLLBACK AND CHAIN start the next transaction from here.
			t.current, t.transaction, t.inTransaction = t.session, t.session, n.Chain
		default:
		}

	case *ast.DiscardStmt:
		if n.Target == ast.DISCARD_ALL {
			t.current, t.session = t.initial, t.initial
		}

	case *ast.VariableSetStmt:
		var searchPath []string
		if n.Kind == ast.VAR_RESET_ALL || (strings.EqualFold(n.Name, "search_path") && (n.Kind == ast.VAR_RESET || n.Kind == ast.VAR_SET_DEFAULT)) {
			searchPath = t.initial
		} else if strings.EqualFold(n.Name, "search_path") && n.Kind == ast.VAR_SET_CURRENT {
			// SET ... FROM CURRENT makes the current search path, which a SET LOCAL may have set, the session's.
			searchPath = t.current
		} else if strings.EqualFold(n.Name, "search_path") && n.Args != nil {
			for _, arg := range n.Args.Items {
				if ac, ok := arg.(*ast.A_Const); ok {
					if s, ok := ac.Val.(*ast.String); ok {
						searchPath = append(searchPath, s.Str)
					}
				}
			}
		}
		if len(searchPath) > 0 {
			t.current = searchPath
			if !n.IsLocal {
				t.session = searchPath
			}
		}
	default:
	}
}

package analysis

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// ValidateQuery reports whether every statement in a script may run in a
// read-only SQL editor, and whether every statement returns data. The error is
// the parse error of the first statement that does not parse.
//
// SELECT, set operations, SHOW, DESCRIBE, EXPLAIN and SET are read-only; SET
// returns no data. EXPLAIN is accepted whatever it plans. USE is rejected
// although QueryTypeOf calls it Select. A result pipe is read-only only when
// its source statement and its query both are, and returns its query's data.
// Statements that parse to nothing are skipped.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	returnsData = true
	for _, stmt := range SplitStatements(statement) {
		if stmt.Empty || strings.TrimSpace(stmt.Text) == "" {
			continue
		}

		file, perr := parser.Parse(stmt.Text)
		if perr != nil {
			return false, false, perr
		}
		if file == nil || len(file.Stmts) == 0 {
			continue
		}

		ok, isSet := classifyForEditor(QueryTypeOf(file.Stmts[0]), file.Stmts[0])
		if !ok {
			return false, false, nil
		}
		if isSet {
			returnsData = false
		}
	}

	return true, returnsData, nil
}

// classifyForEditor returns whether a statement may run in the read-only
// editor, and whether it is a SET (read-only, but returns no data).
func classifyForEditor(qt QueryType, node ast.Node) (bool, bool) {
	if sh, ok := node.(*ast.ShowStmt); ok && sh.Pipe != nil {
		return classifyForEditor(QueryTypeOf(sh.Pipe), sh.Pipe)
	}
	if rs, ok := node.(*ast.ResultScanStmt); ok {
		srcOK, _ := classifyForEditor(QueryTypeOf(rs.Source), rs.Source)
		queryOK, querySet := classifyForEditor(QueryTypeOf(rs.Query), rs.Query)
		return srcOK && queryOK, querySet
	}
	if _, ok := node.(*ast.UseStmt); ok {
		return false, false
	}
	if _, ok := node.(*ast.SetStmt); ok {
		return true, true
	}
	switch qt {
	case Select, SelectInfoSchema, Explain:
		return true, false
	default:
		return false, false
	}
}

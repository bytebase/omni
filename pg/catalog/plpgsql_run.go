package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	plast "github.com/bytebase/omni/pg/plpgsql/ast"
	plparser "github.com/bytebase/omni/pg/plpgsql/parser"
)

// maxPLRunDepth bounds how deeply a run body may CALL or DO another. Past
// it the inner body is not run, where PostgreSQL would run out of stack.
const maxPLRunDepth = 16

// runPLpgSQL runs a compiled PL/pgSQL body, a DO block or a procedure
// CALL runs, when the body is a straight list of statements: no variables,
// no control flow, no exception handlers. Each statement then runs as
// PL/pgSQL runs it: a utility statement executes, a DML statement is
// analyzed against the catalog (its rows are not simulated), and RAISE
// EXCEPTION fails the block. A failed block leaves the catalog as it was.
// A body with anything else is not run, as its path depends on data.
//
// fn holds the parameters a procedure's statements may name, nil for DO.
//
// pg: src/pl/plpgsql/src/pl_exec.c — exec_stmt_block, exec_stmt_execsql, exec_stmt_raise
func (c *Catalog) runPLpgSQL(body string, fn *sqlFnParams) error {
	if c.plRunDepth >= maxPLRunDepth {
		return nil
	}
	block, err := plparser.Parse(body)
	if err != nil {
		return nil
	}
	var stmts []plast.Node
	if !plStraightLine(block, &stmts) {
		return nil
	}
	c.plRunDepth++
	defer func() { c.plRunDepth-- }()
	var saved *Catalog
	fail := func(err error) error {
		if saved != nil {
			warnings, depth := c.warnings, c.plRunDepth
			*c = *saved
			c.warnings, c.plRunDepth = warnings, depth
		}
		return err
	}
	for _, n := range stmts {
		text := ""
		switch v := n.(type) {
		case *plast.PLRaise:
			if v.Level == plast.RaiseLevelNone || v.Level == plast.RaiseLevelError {
				code := v.SQLState
				if code == "" {
					code = CodeRaiseException
				}
				return fail(&Error{Code: code, Message: plRaiseMessage(v.Message)})
			}
			continue
		case *plast.PLExecSQL:
			text = v.SQLText
		case *plast.PLCall:
			text = v.SQLText
		}
		list, err := parser.Parse(text)
		if err != nil || list == nil {
			return fail(nil)
		}
		for _, item := range list.Items {
			stmt := unwrapRawStmt(item)
			if isDML(stmt) {
				if _, err := c.bindSQLFunctionBody([]nodes.Node{stmt}, fn); err != nil {
					return fail(err)
				}
				continue
			}
			if saved == nil {
				saved = c.Clone()
			}
			if err := c.ProcessUtility(stmt); err != nil {
				return fail(err)
			}
		}
	}
	return nil
}

// plStraightLine collects the statements of a block that runs every one of
// them in order, whatever the data: nested blocks without declarations or
// handlers, NULL, RAISE, CALL and DO, and SQL statements that return no
// rows (one that does needs INTO).
func plStraightLine(b *plast.PLBlock, out *[]plast.Node) bool {
	if len(b.Declarations) > 0 || len(b.Exceptions) > 0 {
		return false
	}
	for _, n := range b.Body {
		switch v := n.(type) {
		case *plast.PLBlock:
			if !plStraightLine(v, out) {
				return false
			}
		case *plast.PLNull:
		case *plast.PLCall:
			*out = append(*out, v)
		case *plast.PLRaise:
			if v.CondName != "" || len(v.Options) > 0 || len(v.Params) > 0 {
				return false
			}
			*out = append(*out, v)
		case *plast.PLExecSQL:
			if v.Into != nil || plReturnsRows(v.SQLText) {
				return false
			}
			*out = append(*out, v)
		default:
			return false
		}
	}
	return true
}

// plReturnsRows reports whether a statement returns rows: a SELECT, VALUES
// or TABLE, or DML with RETURNING. PL/pgSQL refuses one without INTO.
//
// pg: src/pl/plpgsql/src/pl_exec.c — exec_stmt_execsql ("query has no destination for result data")
func plReturnsRows(text string) bool {
	list, err := parser.Parse(text)
	if err != nil || list == nil {
		return true
	}
	for _, item := range list.Items {
		switch s := unwrapRawStmt(item).(type) {
		case *nodes.SelectStmt:
			if s.IntoClause == nil {
				return true
			}
		case *nodes.InsertStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return true
			}
		case *nodes.UpdateStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return true
			}
		case *nodes.DeleteStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return true
			}
		case *nodes.MergeStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return true
			}
		}
	}
	return false
}

// plRaiseMessage is a RAISE format with no parameters as the message it
// raises: %% is %.
func plRaiseMessage(format string) string {
	return strings.ReplaceAll(format, "%%", "%")
}

// procParams is the parameters a procedure's statements may name: its IN,
// INOUT and VARIADIC ones.
func procParams(up *UserProc) *sqlFnParams {
	fn := &sqlFnParams{name: up.Name}
	for i, name := range up.ArgNames {
		if up.ArgModes != nil && i < len(up.ArgModes) && (up.ArgModes[i] == 'o' || up.ArgModes[i] == 't') {
			continue
		}
		fn.names = append(fn.names, name)
	}
	fn.types = up.ArgTypes
	for len(fn.names) < len(fn.types) {
		fn.names = append(fn.names, "")
	}
	return fn
}

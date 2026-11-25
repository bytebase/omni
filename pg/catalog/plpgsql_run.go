package catalog

import (
	"errors"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	plast "github.com/bytebase/omni/pg/plpgsql/ast"
	plparser "github.com/bytebase/omni/pg/plpgsql/parser"
)

// maxPLRunDepth bounds how deeply a run body may CALL or DO another. Past
// it the inner body is not run, where PostgreSQL would run out of stack.
const maxPLRunDepth = 16

// errNotRun reports a body the catalog does not run. A body that calls one
// is not run either, as it cannot know whether the call returns.
var errNotRun = errors.New("PL/pgSQL body not run")

// runPLpgSQL runs a compiled PL/pgSQL body, a DO block or a procedure
// CALL runs, when the body is a straight list of statements: no variables,
// no control flow, no exception handlers. Each statement then runs as
// PL/pgSQL runs it: a utility statement executes, a DML statement is
// analyzed against the catalog (its rows are not simulated), and RAISE
// EXCEPTION fails the block. A failed block leaves the catalog as it was.
// A body with anything else returns errNotRun, as its path depends on data.
//
// fn holds the parameters a procedure's statements may name, nil for DO.
//
// pg: src/pl/plpgsql/src/pl_exec.c — exec_stmt_block, exec_stmt_execsql, exec_stmt_raise
func (c *Catalog) runPLpgSQL(body string, fn *sqlFnParams) error {
	if c.plRunDepth >= maxPLRunDepth {
		return errNotRun
	}
	block, err := plparser.Parse(body)
	if err != nil {
		return errNotRun
	}
	var stmts []plast.Node
	if !plStraightLine(block, &stmts) {
		return errNotRun
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
			return fail(errNotRun)
		}
		for _, item := range list.Items {
			stmt := unwrapRawStmt(item)
			if isDML(stmt) {
				if _, err := c.bindSQLFunctionBody([]nodes.Node{stmt}, fn); err != nil {
					return fail(err)
				}
				continue
			}
			if err := plRefusedUtility(stmt); err != nil {
				return fail(err)
			}
			if saved == nil {
				saved = c.Clone()
			}
			switch s := stmt.(type) {
			case *nodes.CallStmt:
				err = c.execCall(s.Funccall, fn)
			case *nodes.DoStmt:
				err = c.execDo(s)
			default:
				err = c.ProcessUtility(stmt)
			}
			if err != nil {
				return fail(err)
			}
		}
	}
	return nil
}

// plRefusedUtility is the error PostgreSQL raises for a statement a function
// may not run, as it cannot run inside a transaction block; errNotRun for a
// transaction statement, which PL/pgSQL handles apart.
//
// pg: src/backend/tcop/utility.c — standard_ProcessUtility, ProcessUtilitySlow, ExecDropStmt (PreventInTransactionBlock)
// pg: src/backend/commands/vacuum.c — vacuum
// pg: src/backend/access/transam/xact.c — PreventInTransactionBlock
func plRefusedUtility(stmt nodes.Node) error {
	what := ""
	switch s := stmt.(type) {
	case *nodes.TransactionStmt:
		return errNotRun
	case *nodes.CreateTableSpaceStmt:
		what = "CREATE TABLESPACE"
	case *nodes.DropTableSpaceStmt:
		what = "DROP TABLESPACE"
	case *nodes.CreatedbStmt:
		what = "CREATE DATABASE"
	case *nodes.DropdbStmt:
		what = "DROP DATABASE"
	case *nodes.AlterSystemStmt:
		what = "ALTER SYSTEM"
	case *nodes.VacuumStmt:
		if s.IsVacuumCmd {
			what = "VACUUM"
		}
	case *nodes.IndexStmt:
		if s.Concurrent {
			what = "CREATE INDEX CONCURRENTLY"
		}
	case *nodes.DropStmt:
		if s.Concurrent {
			what = "DROP INDEX CONCURRENTLY"
		}
	case *nodes.ReindexStmt:
		// REINDEX CONCURRENTLY, DATABASE and SYSTEM refuse a transaction
		// block too, each with its own message; the catalog does not run them.
		return errNotRun
	case *nodes.AlterTableStmt:
		if s.Cmds != nil {
			for _, item := range s.Cmds.Items {
				if cmd, ok := item.(*nodes.AlterTableCmd); ok && cmd.Subtype == int(nodes.AT_DetachPartition) {
					if pc, ok := cmd.Def.(*nodes.PartitionCmd); ok && pc.Concurrent {
						what = "ALTER TABLE ... DETACH CONCURRENTLY"
					}
				}
			}
		}
	}
	if what == "" {
		return nil
	}
	return &Error{Code: CodeActiveSQLTransaction, Message: what + " cannot be executed from a function"}
}

// execCall resolves a CALL's procedure, its arguments naming fn's
// parameters, and runs a PL/pgSQL body when runPLpgSQL can. A procedure
// with SET options runs under settings the catalog does not apply, so it
// is not run.
//
// pg: src/backend/commands/functioncmds.c — ExecuteCallStmt
func (c *Catalog) execCall(fc *nodes.FuncCall, fn *sqlFnParams) error {
	b := newRuleBindings(0, 0)
	b.fn = fn
	if err := c.checkCall(b, fc, &ruleEvalContext{}); err != nil {
		return err
	}
	// checkCall records the procedure after what its arguments use.
	n := len(b.funcDeps)
	if n == 0 {
		return errNotRun
	}
	up := c.userProcs[b.funcDeps[n-1]]
	if up == nil || up.Kind != 'p' || !strings.EqualFold(up.Language, "plpgsql") || len(up.SetConfig) > 0 {
		return errNotRun
	}
	return c.runPLpgSQL(up.Body, procParams(up))
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

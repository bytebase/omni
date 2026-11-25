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
// pg: src/pl/plpgsql/src/pl_exec.c — exec_stmt_block, exec_stmt_execsql, exec_stmt_raise
func (c *Catalog) runPLpgSQL(body string) error {
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
				if _, err := c.bindSQLFunctionBody([]nodes.Node{stmt}, nil); err != nil {
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
				err = c.execCall(s.Funccall)
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
// may not run, as it cannot run inside a transaction block, and errNotRun
// for a utility statement that defines no object: session state, locks,
// notifications, maintenance and transaction control, each with its own
// rules inside a function, which the catalog does not model. Only DDL, CALL
// and DO run.
//
// pg: src/backend/tcop/utility.c — standard_ProcessUtility, ProcessUtilitySlow, ExecDropStmt (PreventInTransactionBlock)
// pg: src/backend/commands/vacuum.c — vacuum
// pg: src/backend/access/transam/xact.c — PreventInTransactionBlock
func plRefusedUtility(stmt nodes.Node) error {
	what := ""
	switch s := stmt.(type) {
	case *nodes.TransactionStmt, *nodes.VariableSetStmt, *nodes.VariableShowStmt,
		*nodes.DiscardStmt, *nodes.LockStmt, *nodes.ListenStmt, *nodes.UnlistenStmt,
		*nodes.NotifyStmt, *nodes.LoadStmt, *nodes.ConstraintsSetStmt, *nodes.CheckPointStmt,
		*nodes.ClusterStmt, *nodes.CopyStmt, *nodes.ExplainStmt, *nodes.PrepareStmt,
		*nodes.ExecuteStmt, *nodes.DeallocateStmt:
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
		if !s.IsVacuumCmd {
			return errNotRun
		}
		what = "VACUUM"
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

// execCall resolves a CALL's procedure and runs its PL/pgSQL body when
// runPLpgSQL can. Only a procedure without parameters, SET options or
// SECURITY DEFINER runs: the catalog binds no PL/pgSQL variable, applies no
// setting and switches no role.
//
// pg: src/backend/commands/functioncmds.c — ExecuteCallStmt
func (c *Catalog) execCall(fc *nodes.FuncCall) error {
	b := newRuleBindings(0, 0)
	if err := c.checkCall(b, fc, &ruleEvalContext{}); err != nil {
		return err
	}
	// checkCall records the procedure after what its arguments use.
	n := len(b.funcDeps)
	if n == 0 {
		return errNotRun
	}
	up := c.userProcs[b.funcDeps[n-1]]
	if up == nil || up.Kind != 'p' || !strings.EqualFold(up.Language, "plpgsql") || len(up.ArgNames) > 0 || len(up.ArgTypes) > 0 || len(up.SetConfig) > 0 || up.SecDef {
		return errNotRun
	}
	return c.runPLpgSQL(up.Body)
}

// plStraightLine collects the statements of a block that runs every one of
// them in order, whatever the data: nested blocks without declarations or
// handlers, NULL, RAISE, CALL and DO, and the SQL statements
// plRunnableSQL takes.
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
			if v.Into != nil || !plRunnableSQL(v.SQLText) {
				return false
			}
			*out = append(*out, v)
		default:
			return false
		}
	}
	return true
}

// plRunnableSQL reports whether runPLpgSQL can run a SQL statement: it
// returns no rows (PL/pgSQL refuses one without INTO: a SELECT, VALUES or
// TABLE, DML with RETURNING, SHOW, EXPLAIN, FETCH, EXECUTE), is no COPY, and names no
// FOUND, the variable every PL/pgSQL block has, which the catalog does not
// bind.
//
// pg: src/pl/plpgsql/src/pl_exec.c — exec_stmt_execsql ("query has no destination for result data")
// pg: src/pl/plpgsql/src/pl_comp.c — plpgsql_build_variable ("found")
func plRunnableSQL(text string) bool {
	list, err := parser.Parse(text)
	if err != nil || list == nil {
		return false
	}
	for _, item := range list.Items {
		switch s := unwrapRawStmt(item).(type) {
		case *nodes.SelectStmt:
			if s.IntoClause == nil {
				return false
			}
		case *nodes.InsertStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return false
			}
		case *nodes.UpdateStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return false
			}
		case *nodes.DeleteStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return false
			}
		case *nodes.MergeStmt:
			if s.ReturningList != nil && len(s.ReturningList.Items) > 0 {
				return false
			}
		case *nodes.VariableShowStmt, *nodes.ExplainStmt, *nodes.ExecuteStmt, *nodes.CopyStmt:
			// COPY to or from the client fails in PL/pgSQL, and the catalog
			// runs no COPY.
			return false
		case *nodes.FetchStmt:
			if !s.Ismove {
				return false
			}
		}
	}
	found := false
	nodes.Inspect(list, func(n nodes.Node) bool {
		if cr, ok := n.(*nodes.ColumnRef); ok && cr.Fields != nil && len(cr.Fields.Items) > 0 {
			if str, ok := cr.Fields.Items[0].(*nodes.String); ok && strings.EqualFold(str.Str, "found") {
				found = true
			}
		}
		return !found
	})
	return !found
}

// plRaiseMessage is a RAISE format with no parameters as the message it
// raises: %% is %.
func plRaiseMessage(format string) string {
	return strings.ReplaceAll(format, "%%", "%")
}

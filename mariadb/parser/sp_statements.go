package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB 11.8 refuses of a stored program's statement while it parses, by what the
// statement is.
//
// In any stored program: USE, LOCK and UNLOCK TABLES, ALTER VIEW, HELP, CHECK TABLE, LOAD
// DATA and LOAD XML, and HANDLER are 1314, ER_SP_BADSTATEMENT (the HELP, CHECK, LOCK,
// UNLOCK and HANDLER rules of sql/sql_yacc.yy; LEX::add_alter_view and LEX::set_command_with_check
// and the USE and LOAD actions); CREATE PROCEDURE, FUNCTION and TRIGGER are 1303,
// ER_SP_NO_RECURSIVE_CREATE (LEX::make_sp_head); ALTER and DROP of a PROCEDURE or FUNCTION
// 1357, ER_SP_NO_DROP_SP (LEX::stmt_drop_routine, stmt_alter_*); CREATE EVENT, and ALTER
// EVENT with DO, 1576, ER_EVENT_RECURSION_FORBIDDEN (event_tail).
//
// In a function or a trigger, sp_head::is_not_allowed_in_function (sql/sp_head.h) refuses
// the flags sp_get_flags_for_command (sql/sp_head.cc) gave the body's statements: dynamic
// SQL and RESET and FLUSH are 1336, ER_STMT_NOT_ALLOWED_IN_SF_OR_TRG; a result set 1415
// (refuseResultSet), a RETURNING clause's included; SET autocommit and SET PASSWORD 1445,
// ER_SP_CANT_SET_AUTOCOMMIT; and a statement that commits 1422,
// ER_COMMIT_NOT_ALLOWED_IN_SF_OR_TRG. MariaDB's list of those is shorter than MySQL's: ALTER
// DATABASE, the server and replication statements do not commit there, and sequences do.
// SET STATEMENT ... FOR is held to what its statement is.

// refuseSPStatement holds a stored program's statement to what the program may hold.
func (p *Parser) refuseSPStatement(stmt nodes.Node, pos int) error {
	if !p.spChecks() {
		return nil
	}
	if s, ok := stmt.(*nodes.SetStatementStmt); ok && s.Stmt != nil {
		if err := p.refuseResultSet(s.Stmt, pos); err != nil {
			return err
		}
		return p.refuseSPStatement(s.Stmt, pos)
	}
	if msg := spBadStatement(stmt); msg != "" {
		return &ParseError{Message: msg, Position: pos}
	}
	if !p.routine.noResultSet {
		return nil
	}
	if set, ok := stmt.(*nodes.SetStmt); ok && p.setsAutocommit(set) {
		return &ParseError{Message: "Not allowed to set autocommit from a stored function or trigger", Position: pos}
	}
	if msg := spNotInFunction(stmt); msg != "" {
		return &ParseError{Message: msg, Position: pos}
	}
	return nil
}

// spBadStatement is what any stored program refuses, "" for what it takes.
func spBadStatement(stmt nodes.Node) string {
	bad := ""
	switch s := stmt.(type) {
	case *nodes.UseStmt:
		bad = "USE"
	case *nodes.LockTablesStmt:
		bad = "LOCK"
	case *nodes.UnlockTablesStmt:
		bad = "UNLOCK"
	case *nodes.AlterViewStmt:
		bad = "ALTER VIEW"
	case *nodes.HelpStmt:
		bad = "HELP"
	case *nodes.CheckTableStmt:
		bad = "CHECK"
	case *nodes.LoadDataStmt:
		bad = "LOAD DATA"
		if s.IsXML {
			bad = "LOAD XML"
		}
	case *nodes.HandlerOpenStmt, *nodes.HandlerReadStmt, *nodes.HandlerCloseStmt:
		bad = "HANDLER"
	case *nodes.CreateFunctionStmt:
		if s.Body == nil {
			return "" // a loadable function
		}
		kind := "FUNCTION"
		if s.IsProcedure {
			kind = "PROCEDURE"
		}
		return "Can't create a " + kind + " from within another stored routine"
	case *nodes.CreateTriggerStmt:
		return "Can't create a TRIGGER from within another stored routine"
	case *nodes.DropRoutineStmt:
		return spNoDropSP(s.IsProcedure)
	case *nodes.AlterRoutineStmt:
		return spNoDropSP(s.IsProcedure)
	case *nodes.CreateEventStmt:
		return "Recursion of EVENT DDL statements is forbidden when body is present"
	case *nodes.AlterEventStmt:
		if s.Body != nil {
			return "Recursion of EVENT DDL statements is forbidden when body is present"
		}
	}
	if bad == "" {
		return ""
	}
	return bad + " is not allowed in stored procedures"
}

func spNoDropSP(procedure bool) string {
	if procedure {
		return "Can't drop or alter a PROCEDURE from within another stored routine"
	}
	return "Can't drop or alter a FUNCTION from within another stored routine"
}

// spNotInFunction is what a function or a trigger refuses beyond result sets, "" for what
// it takes.
func spNotInFunction(stmt nodes.Node) string {
	switch s := stmt.(type) {
	case *nodes.PrepareStmt, *nodes.ExecuteStmt, *nodes.DeallocateStmt:
		return "Dynamic SQL is not allowed in stored function or trigger"
	case *nodes.ResetMasterStmt, *nodes.ResetReplicaStmt:
		return "RESET is not allowed in stored function or trigger"
	case *nodes.FlushStmt:
		return "FLUSH is not allowed in stored function or trigger"
	case *nodes.SetPasswordStmt:
		return "Not allowed to set autocommit from a stored function or trigger"
	case *nodes.CreateTableStmt:
		if s.Temporary {
			return ""
		}
	case *nodes.DropTableStmt:
		if s.Temporary {
			return ""
		}
	case *nodes.CreateSequenceStmt:
		if s.Temporary {
			return ""
		}
	case *nodes.DropSequenceStmt:
		if s.Temporary {
			return ""
		}
	case *nodes.RollbackStmt:
		if s.Savepoint != "" {
			return "" // ROLLBACK TO SAVEPOINT is SQLCOM_ROLLBACK_TO_SAVEPOINT
		}
	case *nodes.AlterEventStmt:
		// One with DO is refused in any stored program; one without commits.
	case *nodes.CommitStmt, *nodes.BeginStmt, *nodes.AlterTableStmt, *nodes.AlterSequenceStmt,
		*nodes.TruncateStmt, *nodes.CreateIndexStmt, *nodes.DropIndexStmt, *nodes.RenameTableStmt,
		*nodes.CreateViewStmt, *nodes.DropViewStmt, *nodes.DropTriggerStmt, *nodes.DropEventStmt,
		*nodes.CreateDatabaseStmt, *nodes.DropDatabaseStmt,
		*nodes.CreateUserStmt, *nodes.AlterUserStmt, *nodes.DropUserStmt, *nodes.RenameUserStmt,
		*nodes.GrantStmt, *nodes.RevokeStmt, *nodes.GrantRoleStmt, *nodes.RevokeRoleStmt,
		*nodes.CreateRoleStmt, *nodes.DropRoleStmt,
		*nodes.InstallPluginStmt, *nodes.UninstallPluginStmt:
	default:
		return ""
	}
	return "Explicit or implicit commit is not allowed in stored function or trigger."
}

// setsAutocommit reports a SET of the autocommit system variable, in any scope and however
// written: autocommit, @@autocommit, @@session.autocommit, GLOBAL autocommit. A bare name the
// stored program declared is its variable, and a scope word holds for the assignments after
// it.
func (p *Parser) setsAutocommit(set *nodes.SetStmt) bool {
	scoped := false
	for _, a := range set.Assignments {
		if a == nil || a.Column == nil || a.CharsetForm {
			continue
		}
		if a.Scope != "" {
			scoped = true
		}
		name := a.Column.Column
		switch {
		case strings.HasPrefix(name, "@@"):
			name = name[2:]
			if i := strings.LastIndexByte(name, '.'); i >= 0 {
				name = name[i+1:]
			}
		case strings.HasPrefix(name, "@"), a.Column.Table != "":
			continue
		case !scoped && p.routine.sp.findVar(name):
			continue
		}
		if strings.EqualFold(name, "autocommit") {
			return true
		}
	}
	return false
}

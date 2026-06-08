package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCreateUserOptionsParse(t *testing.T) {
	stmt := singleStmt(t, `CREATE USER "myco_aad:bob" EXTERNALID "ABC123" PASSWORD DISABLE SYSLOG ACCESS RESTRICTED SESSION TIMEOUT 120 CONNECTION LIMIT UNLIMITED`)
	createUser, ok := stmt.(*nodes.CreateRoleStmt)
	if !ok {
		t.Fatalf("expected CreateRoleStmt, got %T", stmt)
	}
	if createUser.StmtType != nodes.ROLESTMT_USER {
		t.Fatalf("stmt type = %v, want ROLESTMT_USER", createUser.StmtType)
	}
	if createUser.Role != "myco_aad:bob" {
		t.Fatalf("role = %q, want myco_aad:bob", createUser.Role)
	}
	assertDefElemString(t, createUser.Options, "externalid", "ABC123")
	assertDefElemBool(t, createUser.Options, "password_disabled", true)
	assertDefElemString(t, createUser.Options, "syslog_access", "restricted")
	assertDefElemString(t, createUser.Options, "connectionlimit", "unlimited")
	assertDefElemInt(t, createUser.Options, "session_timeout", 120)
}

func TestRedshiftCreateUserExternalNameParse(t *testing.T) {
	stmt := singleStmt(t, `CREATE USER myco_aad:bob EXTERNALID "ABC123" PASSWORD DISABLE`)
	createUser, ok := stmt.(*nodes.CreateRoleStmt)
	if !ok {
		t.Fatalf("expected CreateRoleStmt, got %T", stmt)
	}
	if createUser.Role != "myco_aad:bob" {
		t.Fatalf("role = %q, want myco_aad:bob", createUser.Role)
	}
	assertDefElemString(t, createUser.Options, "externalid", "ABC123")
	assertDefElemBool(t, createUser.Options, "password_disabled", true)
}

func TestRedshiftAlterUserOptionsParse(t *testing.T) {
	stmt := singleStmt(t, `ALTER USER myco_aad:bob EXTERNALID "ABC123" PASSWORD DISABLE SYSLOG ACCESS UNRESTRICTED SESSION TIMEOUT 300 CONNECTION LIMIT UNLIMITED`)
	alterUser, ok := stmt.(*nodes.AlterRoleStmt)
	if !ok {
		t.Fatalf("expected AlterRoleStmt, got %T", stmt)
	}
	if alterUser.Role == nil || alterUser.Role.Rolename != "myco_aad:bob" {
		t.Fatalf("role = %#v, want myco_aad:bob", alterUser.Role)
	}
	assertDefElemString(t, alterUser.Options, "externalid", "ABC123")
	assertDefElemBool(t, alterUser.Options, "password_disabled", true)
	assertDefElemString(t, alterUser.Options, "syslog_access", "unrestricted")
	assertDefElemString(t, alterUser.Options, "connectionlimit", "unlimited")
	assertDefElemInt(t, alterUser.Options, "session_timeout", 300)
}

func TestRedshiftAlterUserResetSessionTimeoutParse(t *testing.T) {
	stmt := singleStmt(t, `ALTER USER dbuser RESET SESSION TIMEOUT`)
	alterUser, ok := stmt.(*nodes.AlterRoleStmt)
	if !ok {
		t.Fatalf("expected AlterRoleStmt, got %T", stmt)
	}
	if alterUser.Role == nil || alterUser.Role.Rolename != "dbuser" {
		t.Fatalf("role = %#v, want dbuser", alterUser.Role)
	}
	assertDefElemBool(t, alterUser.Options, "reset_session_timeout", true)
}

func TestRedshiftAlterUserSetAndResetParametersParse(t *testing.T) {
	for _, sql := range []string{
		`ALTER USER admin SET search_path TO myschema, public`,
		`ALTER USER dbuser RESET search_path`,
		`ALTER USER admin SET search_path TO schema1, public SET statement_timeout = 120000`,
	} {
		stmt := singleStmt(t, sql)
		alterUser, ok := stmt.(*nodes.AlterRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected AlterRoleStmt, got %T", sql, stmt)
		}
		if alterUser.Options == nil || len(alterUser.Options.Items) == 0 {
			t.Fatalf("Parse(%q): expected SET/RESET options", sql)
		}
	}
}

func assertDefElemInt(t *testing.T, list *nodes.List, name string, want int64) {
	t.Helper()
	elem := findDefElem(list, name)
	if elem == nil {
		t.Fatalf("expected option %q in %#v", name, list)
	}
	arg, ok := elem.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected option %q integer arg, got %T", name, elem.Arg)
	}
	if arg.Ival != want {
		t.Fatalf("expected option %q=%d, got %d", name, want, arg.Ival)
	}
}

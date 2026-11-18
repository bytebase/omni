package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseGrantObjectPrivilege tests GRANT SELECT ON table TO user.
func TestParseGrantObjectPrivilege(t *testing.T) {
	result := ParseAndCheck(t, "GRANT SELECT ON employees TO hr_user")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("expected GrantStmt, got %T", raw.Stmt)
	}
	if stmt.AllPriv {
		t.Error("expected AllPriv=false")
	}
	if stmt.Privileges == nil || stmt.Privileges.Len() != 1 {
		t.Fatalf("expected 1 privilege, got %v", stmt.Privileges)
	}
	priv := stmt.Privileges.Items[0].(*ast.String)
	if priv.Str != "SELECT" {
		t.Errorf("expected privilege SELECT, got %q", priv.Str)
	}
	if stmt.OnObject == nil || stmt.OnObject.Name != "EMPLOYEES" {
		t.Errorf("expected ON object EMPLOYEES, got %v", stmt.OnObject)
	}
	if stmt.Grantees == nil || stmt.Grantees.Len() != 1 {
		t.Fatalf("expected 1 grantee, got %v", stmt.Grantees)
	}
	grantee := stmt.Grantees.Items[0].(*ast.String)
	if grantee.Str != "HR_USER" {
		t.Errorf("expected grantee HR_USER, got %q", grantee.Str)
	}
	if stmt.WithGrant {
		t.Error("expected WithGrant=false")
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
}

// TestParseGrantMultiplePrivileges tests GRANT with multiple privileges.
func TestParseGrantMultiplePrivileges(t *testing.T) {
	result := ParseAndCheck(t, "GRANT SELECT, INSERT, UPDATE ON hr.employees TO app_user, admin_user")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("expected GrantStmt, got %T", raw.Stmt)
	}
	if stmt.Privileges.Len() != 3 {
		t.Fatalf("expected 3 privileges, got %d", stmt.Privileges.Len())
	}
	if stmt.OnObject == nil || stmt.OnObject.Schema != "HR" || stmt.OnObject.Name != "EMPLOYEES" {
		t.Errorf("expected ON HR.EMPLOYEES, got %v", stmt.OnObject)
	}
	if stmt.Grantees.Len() != 2 {
		t.Fatalf("expected 2 grantees, got %d", stmt.Grantees.Len())
	}
}

// TestParseGrantAllPrivileges tests GRANT ALL PRIVILEGES.
func TestParseGrantAllPrivileges(t *testing.T) {
	result := ParseAndCheck(t, "GRANT ALL PRIVILEGES ON schema1.t1 TO admin_user WITH GRANT OPTION")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("expected GrantStmt, got %T", raw.Stmt)
	}
	if !stmt.AllPriv {
		t.Error("expected AllPriv=true")
	}
	if !stmt.WithGrant {
		t.Error("expected WithGrant=true")
	}
	if stmt.OnObject == nil || stmt.OnObject.Name != "T1" {
		t.Errorf("expected ON object T1, got %v", stmt.OnObject)
	}
}

// TestParseGrantRole tests GRANT role TO user (no ON clause).
func TestParseGrantRole(t *testing.T) {
	result := ParseAndCheck(t, "GRANT dba_role TO scott")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("expected GrantStmt, got %T", raw.Stmt)
	}
	if stmt.OnObject != nil {
		t.Error("expected nil OnObject for role grant")
	}
	if stmt.Privileges == nil || stmt.Privileges.Len() != 1 {
		t.Fatalf("expected 1 privilege (role name), got %v", stmt.Privileges)
	}
	priv := stmt.Privileges.Items[0].(*ast.String)
	if priv.Str != "DBA_ROLE" {
		t.Errorf("expected DBA_ROLE, got %q", priv.Str)
	}
	if stmt.Grantees.Len() != 1 {
		t.Fatalf("expected 1 grantee, got %d", stmt.Grantees.Len())
	}
}

// TestParseGrantWithAdminOption tests GRANT with WITH ADMIN OPTION.
func TestParseGrantWithAdminOption(t *testing.T) {
	result := ParseAndCheck(t, "GRANT manager_role TO scott WITH ADMIN OPTION")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("expected GrantStmt, got %T", raw.Stmt)
	}
	if !stmt.WithAdmin {
		t.Error("expected WithAdmin=true")
	}
}

// TestParseRevokeObjectPrivilege tests REVOKE SELECT ON table FROM user.
func TestParseRevokeObjectPrivilege(t *testing.T) {
	result := ParseAndCheck(t, "REVOKE SELECT ON employees FROM hr_user")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("expected RevokeStmt, got %T", raw.Stmt)
	}
	if stmt.AllPriv {
		t.Error("expected AllPriv=false")
	}
	if stmt.Privileges == nil || stmt.Privileges.Len() != 1 {
		t.Fatalf("expected 1 privilege, got %v", stmt.Privileges)
	}
	priv := stmt.Privileges.Items[0].(*ast.String)
	if priv.Str != "SELECT" {
		t.Errorf("expected privilege SELECT, got %q", priv.Str)
	}
	if stmt.OnObject == nil || stmt.OnObject.Name != "EMPLOYEES" {
		t.Errorf("expected ON object EMPLOYEES, got %v", stmt.OnObject)
	}
	if stmt.Grantees == nil || stmt.Grantees.Len() != 1 {
		t.Fatalf("expected 1 grantee, got %v", stmt.Grantees)
	}
	grantee := stmt.Grantees.Items[0].(*ast.String)
	if grantee.Str != "HR_USER" {
		t.Errorf("expected grantee HR_USER, got %q", grantee.Str)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
}

// TestParseRevokeAllPrivileges tests REVOKE ALL PRIVILEGES.
func TestParseRevokeAllPrivileges(t *testing.T) {
	result := ParseAndCheck(t, "REVOKE ALL PRIVILEGES ON hr.employees FROM app_user")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("expected RevokeStmt, got %T", raw.Stmt)
	}
	if !stmt.AllPriv {
		t.Error("expected AllPriv=true")
	}
	if stmt.OnObject == nil || stmt.OnObject.Schema != "HR" || stmt.OnObject.Name != "EMPLOYEES" {
		t.Errorf("expected ON HR.EMPLOYEES, got %v", stmt.OnObject)
	}
}

// TestParseRevokeRole tests REVOKE role FROM user.
func TestParseRevokeRole(t *testing.T) {
	result := ParseAndCheck(t, "REVOKE dba_role FROM scott")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("expected RevokeStmt, got %T", raw.Stmt)
	}
	if stmt.OnObject != nil {
		t.Error("expected nil OnObject for role revoke")
	}
	if stmt.Privileges.Len() != 1 {
		t.Fatalf("expected 1 privilege, got %d", stmt.Privileges.Len())
	}
}

// TestParseRevokeMultiplePrivileges tests REVOKE with multiple privileges.
func TestParseRevokeMultiplePrivileges(t *testing.T) {
	result := ParseAndCheck(t, "REVOKE INSERT, DELETE ON orders FROM clerk1, clerk2")
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("expected RevokeStmt, got %T", raw.Stmt)
	}
	if stmt.Privileges.Len() != 2 {
		t.Fatalf("expected 2 privileges, got %d", stmt.Privileges.Len())
	}
	if stmt.Grantees.Len() != 2 {
		t.Fatalf("expected 2 grantees, got %d", stmt.Grantees.Len())
	}
}

// TestParseGrantLocSet tests that Loc is set on GrantStmt.
func TestParseGrantLocSet(t *testing.T) {
	result := ParseAndCheck(t, "GRANT SELECT ON t TO u")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.GrantStmt)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got End=%d", stmt.Loc.End)
	}
}

// TestParseRevokeLocSet tests that Loc is set on RevokeStmt.
func TestParseRevokeLocSet(t *testing.T) {
	result := ParseAndCheck(t, "REVOKE SELECT ON t FROM u")
	raw := result.Items[0].(*ast.RawStmt)
	stmt := raw.Stmt.(*ast.RevokeStmt)
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got End=%d", stmt.Loc.End)
	}
}

// TestParseGrantGrammar tests the forms Oracle 23ai parses and refuses in
// GRANT, REVOKE and traditional NOAUDIT.
func TestParseGrantGrammar(t *testing.T) {
	for _, sql := range []string{
		"GRANT CONNECT, RESOURCE TO u",
		"GRANT CREATE SESSION, r TO x, y WITH ADMIN OPTION CONTAINER = CURRENT",
		"GRANT UPDATE (a, b), INSERT (a) ON t TO x WITH HIERARCHY OPTION",
		"GRANT ON COMMIT REFRESH, QUERY REWRITE ON s.t TO x",
		"GRANT INHERIT PRIVILEGES ON USER u TO x",
		"GRANT READ ON DIRECTORY d TO x",
		"GRANT SELECT ANY TABLE ON SCHEMA s TO x",
		"GRANT r TO FUNCTION f, PROCEDURE s.p, app.PACKAGE k",
		"REVOKE SELECT ON t FROM x CASCADE CONSTRAINTS",
		"REVOKE r FROM PACKAGE k, FUNCTION s.f",
		// The page's revoke_roles_from_programs branch carries the
		// statement's REVOKE and its program unit's kind again.
		"REVOKE REVOKE r FROM PROCEDURE p",
		"REVOKE REVOKE r FROM PROCEDURE FUNCTION p",
		"REVOKE REVOKE r FROM PROCEDURE s.FUNCTION p",
		"REVOKE REVOKE FROM, r FROM PROCEDURE p",
		"REVOKE REVOKE, r FROM PROCEDURE p",
		"REVOKE REVOKE r FROM PROCEDURE FUNCTION p CONTAINER = CURRENT",
		"REVOKE REVOKE r FROM PROCEDURE FUNCTION CONTAINER = CURRENT",
		"REVOKE REVOKE ALL FROM PROCEDURE p",
		"NOAUDIT SELECT, UPDATE ON t WHENEVER NOT SUCCESSFUL",
		"NOAUDIT TABLE, CREATE SESSION BY u, v CONTAINER = ALL",
		"NOAUDIT CONTEXT",
		`GRANT "app role" TO u`,
		`REVOKE "app role" FROM u`,
		// Oracle 23ai parses a REVOKE of columns and refuses it after the
		// parse (ORA-01750).
		"REVOKE UPDATE (a) ON t FROM u",
	} {
		ParseAndCheck(t, sql)
	}
	for _, sql := range []string{
		"GRANT a ON t TO x",
		"GRANT a b TO x",
		"GRANT SELECT (a) ON t TO x",
		"GRANT SELECT ON TABLE t TO x",
		"GRANT ALL ON SCHEMA s TO x",
		"GRANT SELECT ON t TO x WITH ADMIN OPTION",
		"GRANT r TO FUNCTION f WITH ADMIN OPTION",
		"REVOKE SELECT ON t FROM x CASCADE CONSTRAINTS FORCE",
		"REVOKE r FROM PROCEDURE FUNCTION p",
		"REVOKE REVOKE r FROM PROCEDURE",
		"REVOKE REVOKE REVOKE r FROM PROCEDURE p",
		"REVOKE REVOKE r FROM PROCEDURE FUNCTION s.p",
		"REVOKE REVOKE r FROM s.PROCEDURE FUNCTION p",
		"REVOKE REVOKE r FROM PROCEDURE SELECT.FUNCTION p",
		"REVOKE REVOKE ALL PRIVILEGES FROM PROCEDURE p",
		"REVOKE REVOKE CREATE SESSION FROM PROCEDURE p",
		"NOAUDIT a",
		"NOAUDIT SELECT",
		"NOAUDIT TABLE ON t",
		"NOAUDIT SELECT ON t BY u",
		"NOAUDIT TABLE BY ACCESS",
		"GRANT r TO FUNCTION f, s garbage PROCEDURE p",
		"REVOKE r FROM FUNCTION f, s garbage PROCEDURE p",
	} {
		ParseShouldFail(t, sql)
	}
}

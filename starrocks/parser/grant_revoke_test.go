package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseGrantStmt(t *testing.T, sql string) *ast.GrantStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("expected *ast.GrantStmt, got %T", n)
	}
	return stmt
}

func parseRevokeStmt(t *testing.T, sql string) *ast.RevokeStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("expected *ast.RevokeStmt, got %T", n)
	}
	return stmt
}

// ---------------------------------------------------------------------------
// GRANT tests
// ---------------------------------------------------------------------------

// TestGrantSelectOnAllObjects verifies:
//
//	GRANT SELECT_PRIV ON *.*.* TO 'jack'@'%'
func TestGrantSelectOnAllObjects(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT_PRIV ON *.*.* TO 'jack'@'%'")

	if len(stmt.Privileges) != 1 || stmt.Privileges[0] != "SELECT_PRIV" {
		t.Errorf("privileges = %v, want [SELECT_PRIV]", stmt.Privileges)
	}
	if stmt.ObjectType != "TABLE" {
		t.Errorf("object type = %q, want TABLE", stmt.ObjectType)
	}
	if stmt.Object == nil || len(stmt.Object.Parts) != 3 {
		t.Errorf("object parts = %v, want [*, *, *]", stmt.Object)
	}
	if len(stmt.Grantees) != 1 || stmt.Grantees[0] != "jack@%" {
		t.Errorf("grantees = %v, want [jack@%%]", stmt.Grantees)
	}
}

// TestGrantMultiplePrivsOnTable verifies:
//
//	GRANT SELECT_PRIV, ALTER_PRIV, LOAD_PRIV ON ctl1.db1.tbl1 TO 'jack'@'192.8.%'
func TestGrantMultiplePrivsOnTable(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT_PRIV, ALTER_PRIV, LOAD_PRIV ON ctl1.db1.tbl1 TO 'jack'@'192.8.%'")

	if len(stmt.Privileges) != 3 {
		t.Fatalf("want 3 privileges, got %d: %v", len(stmt.Privileges), stmt.Privileges)
	}
	if stmt.Privileges[0] != "SELECT_PRIV" {
		t.Errorf("priv[0] = %q, want SELECT_PRIV", stmt.Privileges[0])
	}
	if stmt.Privileges[1] != "ALTER_PRIV" {
		t.Errorf("priv[1] = %q, want ALTER_PRIV", stmt.Privileges[1])
	}
	if stmt.Privileges[2] != "LOAD_PRIV" {
		t.Errorf("priv[2] = %q, want LOAD_PRIV", stmt.Privileges[2])
	}
	if stmt.Object == nil || len(stmt.Object.Parts) != 3 {
		t.Errorf("object parts count = %d, want 3", len(stmt.Object.Parts))
	}
	if stmt.Grantees[0] != "jack@192.8.%" {
		t.Errorf("grantee = %q, want jack@192.8.%%", stmt.Grantees[0])
	}
}

// TestGrantLoadPrivOnDbStar verifies:
//
//	GRANT LOAD_PRIV ON ctl1.db1.* TO ROLE 'my_role'
func TestGrantLoadPrivOnDbStar(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT LOAD_PRIV ON ctl1.db1.* TO ROLE 'my_role'")

	if len(stmt.Privileges) != 1 || stmt.Privileges[0] != "LOAD_PRIV" {
		t.Errorf("privileges = %v, want [LOAD_PRIV]", stmt.Privileges)
	}
	if stmt.ToType != "ROLE" {
		t.Errorf("toType = %q, want ROLE", stmt.ToType)
	}
	if len(stmt.Grantees) != 1 || stmt.Grantees[0] != "my_role" {
		t.Errorf("grantees = %v, want [my_role]", stmt.Grantees)
	}
}

// TestGrantUsageOnResourceStar verifies:
//
//	GRANT USAGE_PRIV ON RESOURCE * TO 'jack'@'%'
func TestGrantUsageOnResourceStar(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON RESOURCE * TO 'jack'@'%'")

	if stmt.ObjectType != "RESOURCE" {
		t.Errorf("object type = %q, want RESOURCE", stmt.ObjectType)
	}
	if stmt.Object == nil || len(stmt.Object.Parts) != 1 || stmt.Object.Parts[0] != "*" {
		t.Errorf("object = %v, want [*]", stmt.Object)
	}
}

// TestGrantUsageOnNamedResource verifies:
//
//	GRANT USAGE_PRIV ON RESOURCE 'spark_resource' TO 'jack'@'%'
func TestGrantUsageOnNamedResource(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON RESOURCE 'spark_resource' TO 'jack'@'%'")

	if stmt.ObjectType != "RESOURCE" {
		t.Errorf("object type = %q, want RESOURCE", stmt.ObjectType)
	}
	if stmt.Object == nil || stmt.Object.Parts[0] != "spark_resource" {
		t.Errorf("object = %v, want [spark_resource]", stmt.Object)
	}
}

// TestGrantUsageOnResourceToRole verifies:
//
//	GRANT USAGE_PRIV ON RESOURCE 'spark_resource' TO ROLE 'my_role'
func TestGrantUsageOnResourceToRole(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON RESOURCE 'spark_resource' TO ROLE 'my_role'")

	if stmt.ToType != "ROLE" {
		t.Errorf("toType = %q, want ROLE", stmt.ToType)
	}
	if stmt.Grantees[0] != "my_role" {
		t.Errorf("grantee = %q, want my_role", stmt.Grantees[0])
	}
}

// TestGrantRoleToUser verifies:
//
//	GRANT 'role1', 'role2' TO 'jack'@'%'
func TestGrantRoleToUser(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT 'role1', 'role2' TO 'jack'@'%'")

	if len(stmt.Roles) != 2 {
		t.Fatalf("roles = %v, want [role1, role2]", stmt.Roles)
	}
	if stmt.Roles[0] != "role1" || stmt.Roles[1] != "role2" {
		t.Errorf("roles = %v, want [role1, role2]", stmt.Roles)
	}
	if len(stmt.Grantees) != 1 || stmt.Grantees[0] != "jack@%" {
		t.Errorf("grantees = %v, want [jack@%%]", stmt.Grantees)
	}
}

// TestGrantUsageOnWorkloadGroup verifies:
//
//	GRANT USAGE_PRIV ON WORKLOAD GROUP 'g1' TO 'jack'@'%'
func TestGrantUsageOnWorkloadGroup(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON WORKLOAD GROUP 'g1' TO 'jack'@'%'")

	if stmt.ObjectType != "WORKLOAD GROUP" {
		t.Errorf("object type = %q, want WORKLOAD GROUP", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "g1" {
		t.Errorf("object = %v, want [g1]", stmt.Object.Parts)
	}
}

// TestGrantUsageOnWorkloadGroupWildcard verifies:
//
//	GRANT USAGE_PRIV ON WORKLOAD GROUP '%' TO 'jack'@'%'
func TestGrantUsageOnWorkloadGroupWildcard(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON WORKLOAD GROUP '%' TO 'jack'@'%'")

	if stmt.ObjectType != "WORKLOAD GROUP" {
		t.Errorf("object type = %q, want WORKLOAD GROUP", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "%" {
		t.Errorf("object part = %q, want %%", stmt.Object.Parts[0])
	}
}

// TestGrantUsageOnWorkloadGroupToRole verifies:
//
//	GRANT USAGE_PRIV ON WORKLOAD GROUP 'g1' TO ROLE 'my_role'
func TestGrantUsageOnWorkloadGroupToRole(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON WORKLOAD GROUP 'g1' TO ROLE 'my_role'")

	if stmt.ObjectType != "WORKLOAD GROUP" {
		t.Errorf("object type = %q, want WORKLOAD GROUP", stmt.ObjectType)
	}
	if stmt.ToType != "ROLE" {
		t.Errorf("toType = %q, want ROLE", stmt.ToType)
	}
}

// TestGrantShowViewPriv verifies:
//
//	GRANT SHOW_VIEW_PRIV ON db1.view1 TO 'jack'@'%'
func TestGrantShowViewPriv(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SHOW_VIEW_PRIV ON db1.view1 TO 'jack'@'%'")

	if len(stmt.Privileges) != 1 || stmt.Privileges[0] != "SHOW_VIEW_PRIV" {
		t.Errorf("privileges = %v, want [SHOW_VIEW_PRIV]", stmt.Privileges)
	}
	if stmt.Object == nil || len(stmt.Object.Parts) != 2 {
		t.Errorf("object parts = %v, want [db1, view1]", stmt.Object.Parts)
	}
}

// TestGrantUsageOnComputeGroup verifies:
//
//	GRANT USAGE_PRIV ON COMPUTE GROUP 'group1' TO 'jack'@'%'
func TestGrantUsageOnComputeGroup(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON COMPUTE GROUP 'group1' TO 'jack'@'%'")

	if stmt.ObjectType != "COMPUTE GROUP" {
		t.Errorf("object type = %q, want COMPUTE GROUP", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "group1" {
		t.Errorf("object part = %q, want group1", stmt.Object.Parts[0])
	}
}

// TestGrantUsageOnComputeGroupToRole verifies:
//
//	GRANT USAGE_PRIV ON COMPUTE GROUP 'group1' TO ROLE 'my_role'
func TestGrantUsageOnComputeGroupToRole(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON COMPUTE GROUP 'group1' TO ROLE 'my_role'")

	if stmt.ObjectType != "COMPUTE GROUP" {
		t.Errorf("object type = %q, want COMPUTE GROUP", stmt.ObjectType)
	}
	if stmt.ToType != "ROLE" {
		t.Errorf("toType = %q, want ROLE", stmt.ToType)
	}
}

// TestGrantUsageOnComputeGroupWildcard verifies:
//
//	GRANT USAGE_PRIV ON COMPUTE GROUP '*' TO 'jack'@'%'
func TestGrantUsageOnComputeGroupWildcard(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON COMPUTE GROUP '*' TO 'jack'@'%'")

	if stmt.ObjectType != "COMPUTE GROUP" {
		t.Errorf("object type = %q, want COMPUTE GROUP", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "*" {
		t.Errorf("object part = %q, want *", stmt.Object.Parts[0])
	}
}

// TestGrantUsageOnStorageVault verifies:
//
//	GRANT USAGE_PRIV ON STORAGE VAULT 'vault1' TO 'jack'@'%'
func TestGrantUsageOnStorageVault(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON STORAGE VAULT 'vault1' TO 'jack'@'%'")

	if stmt.ObjectType != "STORAGE VAULT" {
		t.Errorf("object type = %q, want STORAGE VAULT", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "vault1" {
		t.Errorf("object part = %q, want vault1", stmt.Object.Parts[0])
	}
}

// TestGrantUsageOnStorageVaultToRole verifies:
//
//	GRANT USAGE_PRIV ON STORAGE VAULT 'vault1' TO ROLE 'my_role'
func TestGrantUsageOnStorageVaultToRole(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON STORAGE VAULT 'vault1' TO ROLE 'my_role'")

	if stmt.ObjectType != "STORAGE VAULT" {
		t.Errorf("object type = %q, want STORAGE VAULT", stmt.ObjectType)
	}
	if stmt.ToType != "ROLE" {
		t.Errorf("toType = %q, want ROLE", stmt.ToType)
	}
	if stmt.Grantees[0] != "my_role" {
		t.Errorf("grantee = %q, want my_role", stmt.Grantees[0])
	}
}

// TestGrantUsageOnStorageVaultWildcard verifies:
//
//	GRANT USAGE_PRIV ON STORAGE VAULT '*' TO 'jack'@'%'
func TestGrantUsageOnStorageVaultWildcard(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE_PRIV ON STORAGE VAULT '*' TO 'jack'@'%'")

	if stmt.ObjectType != "STORAGE VAULT" {
		t.Errorf("object type = %q, want STORAGE VAULT", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "*" {
		t.Errorf("object part = %q, want *", stmt.Object.Parts[0])
	}
}

// ---------------------------------------------------------------------------
// REVOKE tests
// ---------------------------------------------------------------------------

// TestRevokeSelectOnDb verifies:
//
//	REVOKE SELECT_PRIV ON db1.* FROM 'jack'@'192.%'
func TestRevokeSelectOnDb(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE SELECT_PRIV ON db1.* FROM 'jack'@'192.%'")

	if len(stmt.Privileges) != 1 || stmt.Privileges[0] != "SELECT_PRIV" {
		t.Errorf("privileges = %v, want [SELECT_PRIV]", stmt.Privileges)
	}
	if stmt.ObjectType != "TABLE" {
		t.Errorf("object type = %q, want TABLE", stmt.ObjectType)
	}
	if stmt.Object == nil || len(stmt.Object.Parts) != 2 {
		t.Errorf("object parts = %v, want [db1, *]", stmt.Object.Parts)
	}
	if len(stmt.Revokees) != 1 || stmt.Revokees[0] != "jack@192.%" {
		t.Errorf("revokees = %v, want [jack@192.%%]", stmt.Revokees)
	}
}

// TestRevokeUsageOnResource verifies:
//
//	REVOKE USAGE_PRIV ON RESOURCE 'spark_resource' FROM 'jack'@'192.%'
func TestRevokeUsageOnResource(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON RESOURCE 'spark_resource' FROM 'jack'@'192.%'")

	if stmt.ObjectType != "RESOURCE" {
		t.Errorf("object type = %q, want RESOURCE", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "spark_resource" {
		t.Errorf("object part = %q, want spark_resource", stmt.Object.Parts[0])
	}
	if stmt.Revokees[0] != "jack@192.%" {
		t.Errorf("revokee = %q, want jack@192.%%", stmt.Revokees[0])
	}
}

// TestRevokeRolesFromUser verifies:
//
//	REVOKE 'role1', 'role2' FROM 'jack'@'192.%'
func TestRevokeRolesFromUser(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE 'role1', 'role2' FROM 'jack'@'192.%'")

	if len(stmt.Roles) != 2 {
		t.Fatalf("roles = %v, want [role1, role2]", stmt.Roles)
	}
	if stmt.Roles[0] != "role1" || stmt.Roles[1] != "role2" {
		t.Errorf("roles = %v, want [role1, role2]", stmt.Roles)
	}
	if stmt.Revokees[0] != "jack@192.%" {
		t.Errorf("revokee = %q, want jack@192.%%", stmt.Revokees[0])
	}
}

// TestRevokeUsageOnWorkloadGroup verifies:
//
//	REVOKE USAGE_PRIV ON WORKLOAD GROUP 'g1' FROM 'jack'@'%'
func TestRevokeUsageOnWorkloadGroup(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON WORKLOAD GROUP 'g1' FROM 'jack'@'%'")

	if stmt.ObjectType != "WORKLOAD GROUP" {
		t.Errorf("object type = %q, want WORKLOAD GROUP", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "g1" {
		t.Errorf("object part = %q, want g1", stmt.Object.Parts[0])
	}
}

// TestRevokeUsageOnWorkloadGroupWildcard verifies:
//
//	REVOKE USAGE_PRIV ON WORKLOAD GROUP '%' FROM 'jack'@'%'
func TestRevokeUsageOnWorkloadGroupWildcard(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON WORKLOAD GROUP '%' FROM 'jack'@'%'")

	if stmt.ObjectType != "WORKLOAD GROUP" {
		t.Errorf("object type = %q, want WORKLOAD GROUP", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "%" {
		t.Errorf("object part = %q, want %%", stmt.Object.Parts[0])
	}
}

// TestRevokeUsageOnComputeGroup verifies:
//
//	REVOKE USAGE_PRIV ON COMPUTE GROUP 'group1' FROM 'jack'@'%'
func TestRevokeUsageOnComputeGroup(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON COMPUTE GROUP 'group1' FROM 'jack'@'%'")

	if stmt.ObjectType != "COMPUTE GROUP" {
		t.Errorf("object type = %q, want COMPUTE GROUP", stmt.ObjectType)
	}
}

// TestRevokeUsageOnComputeGroupFromRole verifies:
//
//	REVOKE USAGE_PRIV ON COMPUTE GROUP 'group1' FROM ROLE 'my_role'
func TestRevokeUsageOnComputeGroupFromRole(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON COMPUTE GROUP 'group1' FROM ROLE 'my_role'")

	if stmt.FromType != "ROLE" {
		t.Errorf("fromType = %q, want ROLE", stmt.FromType)
	}
	if stmt.Revokees[0] != "my_role" {
		t.Errorf("revokee = %q, want my_role", stmt.Revokees[0])
	}
}

// TestRevokeUsageOnStorageVault verifies:
//
//	REVOKE USAGE_PRIV ON STORAGE VAULT 'vault1' FROM 'jack'@'%'
func TestRevokeUsageOnStorageVault(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON STORAGE VAULT 'vault1' FROM 'jack'@'%'")

	if stmt.ObjectType != "STORAGE VAULT" {
		t.Errorf("object type = %q, want STORAGE VAULT", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "vault1" {
		t.Errorf("object part = %q, want vault1", stmt.Object.Parts[0])
	}
}

// TestRevokeUsageOnStorageVaultFromRole verifies:
//
//	REVOKE USAGE_PRIV ON STORAGE VAULT 'vault1' FROM ROLE 'my_role'
func TestRevokeUsageOnStorageVaultFromRole(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON STORAGE VAULT 'vault1' FROM ROLE 'my_role'")

	if stmt.ObjectType != "STORAGE VAULT" {
		t.Errorf("object type = %q, want STORAGE VAULT", stmt.ObjectType)
	}
	if stmt.FromType != "ROLE" {
		t.Errorf("fromType = %q, want ROLE", stmt.FromType)
	}
}

// TestRevokeUsageOnStorageVaultWildcard verifies:
//
//	REVOKE USAGE_PRIV ON STORAGE VAULT '%' FROM 'jack'@'%'
func TestRevokeUsageOnStorageVaultWildcard(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE USAGE_PRIV ON STORAGE VAULT '%' FROM 'jack'@'%'")

	if stmt.ObjectType != "STORAGE VAULT" {
		t.Errorf("object type = %q, want STORAGE VAULT", stmt.ObjectType)
	}
	if stmt.Object.Parts[0] != "%" {
		t.Errorf("object part = %q, want %%", stmt.Object.Parts[0])
	}
}

// ---------------------------------------------------------------------------
// AST tag tests
// ---------------------------------------------------------------------------

func TestGrantRevokeNodeTags(t *testing.T) {
	grantNode := parseOne(t, "GRANT SELECT_PRIV ON *.* TO 'user'@'%'")
	if grantNode.Tag() != ast.T_GrantStmt {
		t.Errorf("grant tag = %v, want T_GrantStmt", grantNode.Tag())
	}

	revokeNode := parseOne(t, "REVOKE SELECT_PRIV ON db.* FROM 'user'@'%'")
	if revokeNode.Tag() != ast.T_RevokeStmt {
		t.Errorf("revoke tag = %v, want T_RevokeStmt", revokeNode.Tag())
	}
}

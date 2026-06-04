package parser

import (
	"testing"
)

// This file is the parser-dcl-tcl node's structural correctness gate for the
// DCL statements (CREATE/DROP ROLE, GRANT/REVOKE roles & privileges, DENY). The
// authoritative accept/reject differential against the live Trino 481 oracle
// lives in oracle_dcl_tcl_test.go; here we pin the AST shape and the
// role-vs-privilege disambiguation.

func TestCreateDropRole_Structure(t *testing.T) {
	t.Run("create_role", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "CREATE ROLE admin").(*CreateRoleStmt)
		if !ok {
			t.Fatalf("got %T, want *CreateRoleStmt", stmt)
		}
		if stmt.Name.Normalize() != "admin" {
			t.Errorf("Name = %q", stmt.Name.Normalize())
		}
		if stmt.Admin != nil || stmt.Catalog != nil {
			t.Errorf("unexpected Admin=%v Catalog=%v", stmt.Admin, stmt.Catalog)
		}
	})

	t.Run("create_role_with_admin_user", func(t *testing.T) {
		stmt := dclParseOne(t, "CREATE ROLE moderator WITH ADMIN USER bob").(*CreateRoleStmt)
		if stmt.Admin == nil || stmt.Admin.Kind != GrantorPrincipal {
			t.Fatalf("Admin = %+v, want a principal grantor", stmt.Admin)
		}
		if stmt.Admin.Principal.Kind != PrincipalUser || stmt.Admin.Principal.Name.Normalize() != "bob" {
			t.Errorf("Admin principal = %+v, want USER bob", stmt.Admin.Principal)
		}
	})

	t.Run("create_role_in_catalog", func(t *testing.T) {
		stmt := dclParseOne(t, "CREATE ROLE admin IN hive").(*CreateRoleStmt)
		if stmt.Catalog == nil || stmt.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v, want hive", stmt.Catalog)
		}
	})

	t.Run("create_role_with_admin_current_role_in_catalog", func(t *testing.T) {
		stmt := dclParseOne(t, "CREATE ROLE r WITH ADMIN CURRENT_ROLE IN hive").(*CreateRoleStmt)
		if stmt.Admin == nil || stmt.Admin.Kind != GrantorCurrentRole {
			t.Errorf("Admin = %+v, want CURRENT_ROLE", stmt.Admin)
		}
		if stmt.Catalog == nil || stmt.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v, want hive", stmt.Catalog)
		}
	})

	t.Run("drop_role", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "DROP ROLE admin").(*DropRoleStmt)
		if !ok {
			t.Fatalf("got %T, want *DropRoleStmt", stmt)
		}
		if stmt.IfExists {
			t.Errorf("IfExists = true, want false")
		}
		if stmt.Name.Normalize() != "admin" {
			t.Errorf("Name = %q", stmt.Name.Normalize())
		}
	})

	t.Run("drop_role_if_exists_in_catalog", func(t *testing.T) {
		stmt := dclParseOne(t, "DROP ROLE IF EXISTS analyst IN hive").(*DropRoleStmt)
		if !stmt.IfExists {
			t.Errorf("IfExists = false, want true")
		}
		if stmt.Name.Normalize() != "analyst" {
			t.Errorf("Name = %q", stmt.Name.Normalize())
		}
		if stmt.Catalog == nil || stmt.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v, want hive", stmt.Catalog)
		}
	})
}

func TestGrantRoles_Structure(t *testing.T) {
	t.Run("grant_role_to_user", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "GRANT bar TO USER foo").(*GrantRolesStmt)
		if !ok {
			t.Fatalf("got %T, want *GrantRolesStmt", stmt)
		}
		if len(stmt.Roles) != 1 || stmt.Roles[0].Normalize() != "bar" {
			t.Errorf("Roles = %v, want [bar]", stmt.Roles)
		}
		if len(stmt.Grantees) != 1 || stmt.Grantees[0].Kind != PrincipalUser || stmt.Grantees[0].Name.Normalize() != "foo" {
			t.Errorf("Grantees = %+v, want [USER foo]", stmt.Grantees)
		}
		if stmt.AdminOption {
			t.Errorf("AdminOption = true, want false")
		}
	})

	t.Run("grant_roles_multi_with_admin_option", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT bar, foo TO USER baz, ROLE qux WITH ADMIN OPTION").(*GrantRolesStmt)
		if len(stmt.Roles) != 2 {
			t.Errorf("got %d roles, want 2", len(stmt.Roles))
		}
		if len(stmt.Grantees) != 2 {
			t.Fatalf("got %d grantees, want 2", len(stmt.Grantees))
		}
		if stmt.Grantees[0].Kind != PrincipalUser || stmt.Grantees[1].Kind != PrincipalRole {
			t.Errorf("grantee kinds = %v,%v want USER,ROLE", stmt.Grantees[0].Kind, stmt.Grantees[1].Kind)
		}
		if !stmt.AdminOption {
			t.Errorf("AdminOption = false, want true")
		}
	})

	t.Run("grant_role_granted_by_in", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT bar TO USER foo GRANTED BY CURRENT_USER IN hive").(*GrantRolesStmt)
		if stmt.GrantedBy == nil || stmt.GrantedBy.Kind != GrantorCurrentUser {
			t.Errorf("GrantedBy = %+v, want CURRENT_USER", stmt.GrantedBy)
		}
		if stmt.Catalog == nil || stmt.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v, want hive", stmt.Catalog)
		}
	})

	t.Run("grant_role_bare_grantee", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT bar TO foo").(*GrantRolesStmt)
		if stmt.Grantees[0].Kind != PrincipalUnspecified {
			t.Errorf("grantee kind = %v, want Unspecified", stmt.Grantees[0].Kind)
		}
	})

	t.Run("grant_role_all_as_rolename", func(t *testing.T) {
		// `ALL` here is a role name (non-reserved), not all-privileges, because
		// it is followed by TO (role form), not ON.
		stmt := dclParseOne(t, "GRANT ALL TO foo").(*GrantRolesStmt)
		if len(stmt.Roles) != 1 || stmt.Roles[0].Normalize() != "all" {
			t.Errorf("Roles = %v, want [all]", stmt.Roles)
		}
	})
}

func TestRevokeRoles_Structure(t *testing.T) {
	t.Run("revoke_role", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "REVOKE bar FROM USER foo").(*RevokeRolesStmt)
		if !ok {
			t.Fatalf("got %T, want *RevokeRolesStmt", stmt)
		}
		if stmt.AdminOptionFor {
			t.Errorf("AdminOptionFor = true, want false")
		}
		if len(stmt.Roles) != 1 || stmt.Roles[0].Normalize() != "bar" {
			t.Errorf("Roles = %v", stmt.Roles)
		}
	})

	t.Run("revoke_admin_option_for_multi", func(t *testing.T) {
		stmt := dclParseOne(t, "REVOKE ADMIN OPTION FOR bar, foo FROM USER baz, ROLE qux").(*RevokeRolesStmt)
		if !stmt.AdminOptionFor {
			t.Errorf("AdminOptionFor = false, want true")
		}
		if len(stmt.Roles) != 2 || len(stmt.Grantees) != 2 {
			t.Errorf("got %d roles / %d grantees, want 2/2", len(stmt.Roles), len(stmt.Grantees))
		}
	})

	t.Run("revoke_role_granted_by_role_in", func(t *testing.T) {
		stmt := dclParseOne(t, "REVOKE bar FROM foo GRANTED BY CURRENT_ROLE IN hive").(*RevokeRolesStmt)
		if stmt.GrantedBy == nil || stmt.GrantedBy.Kind != GrantorCurrentRole {
			t.Errorf("GrantedBy = %+v, want CURRENT_ROLE", stmt.GrantedBy)
		}
		if stmt.Catalog == nil || stmt.Catalog.Normalize() != "hive" {
			t.Errorf("Catalog = %v", stmt.Catalog)
		}
	})
}

func TestGrantPrivileges_Structure(t *testing.T) {
	t.Run("grant_select_insert_on_table", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "GRANT INSERT, SELECT ON orders TO alice").(*GrantPrivStmt)
		if !ok {
			t.Fatalf("got %T, want *GrantPrivStmt", stmt)
		}
		if stmt.AllPrivileges {
			t.Errorf("AllPrivileges = true, want false")
		}
		if len(stmt.Privileges) != 2 {
			t.Fatalf("got %d privileges, want 2", len(stmt.Privileges))
		}
		if stmt.On.Qualifier != nil || stmt.On.Object.Normalize() != "orders" {
			t.Errorf("On = %+v, want unqualified orders", stmt.On)
		}
		if stmt.On.Branch != nil {
			t.Errorf("Branch = %v, want nil", stmt.On.Branch)
		}
		if stmt.Grantee.Kind != PrincipalUnspecified || stmt.Grantee.Name.Normalize() != "alice" {
			t.Errorf("Grantee = %+v, want bare alice", stmt.Grantee)
		}
	})

	t.Run("grant_delete_on_schema", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT DELETE ON SCHEMA finance TO bob").(*GrantPrivStmt)
		if !stmt.On.IsSchema() || stmt.On.Object.Normalize() != "finance" {
			t.Errorf("On = %+v, want SCHEMA finance", stmt.On)
		}
	})

	t.Run("grant_select_on_table_kw_with_grant_option", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT SELECT ON TABLE orders TO alice WITH GRANT OPTION").(*GrantPrivStmt)
		if !stmt.On.IsTable() {
			t.Errorf("On.Qualifier = %v, want TABLE", stmt.On.Qualifier)
		}
		if !stmt.GrantOption {
			t.Errorf("GrantOption = false, want true")
		}
	})

	t.Run("grant_select_to_role_public", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT SELECT ON orders TO ROLE PUBLIC").(*GrantPrivStmt)
		if stmt.Grantee.Kind != PrincipalRole || stmt.Grantee.Name.Normalize() != "public" {
			t.Errorf("Grantee = %+v, want ROLE PUBLIC", stmt.Grantee)
		}
	})

	t.Run("grant_all_privileges", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT ALL PRIVILEGES ON test TO alice").(*GrantPrivStmt)
		if !stmt.AllPrivileges || len(stmt.Privileges) != 0 {
			t.Errorf("got AllPrivileges=%v Privileges=%v, want true/empty", stmt.AllPrivileges, stmt.Privileges)
		}
	})

	t.Run("grant_bare_all", func(t *testing.T) {
		// `GRANT ALL ON test` — bare ALL followed by ON is all-privileges (D3).
		stmt := dclParseOne(t, "GRANT ALL ON test TO alice").(*GrantPrivStmt)
		if !stmt.AllPrivileges {
			t.Errorf("AllPrivileges = false, want true")
		}
	})

	t.Run("grant_on_branch", func(t *testing.T) {
		// D2 docs-ahead extension: ON BRANCH <branch> IN <table>.
		stmt := dclParseOne(t, "GRANT INSERT ON BRANCH audit IN orders TO alice").(*GrantPrivStmt)
		if stmt.On.Branch == nil || stmt.On.Branch.Normalize() != "audit" {
			t.Errorf("Branch = %v, want audit", stmt.On.Branch)
		}
		if stmt.On.Object.Normalize() != "orders" {
			t.Errorf("Object = %v, want orders", stmt.On.Object)
		}
	})

	t.Run("grant_open_privilege_vocab", func(t *testing.T) {
		// D1: a non-reserved identifier is a valid privilege.
		stmt := dclParseOne(t, "GRANT mypriv ON test TO alice").(*GrantPrivStmt)
		if len(stmt.Privileges) != 1 || stmt.Privileges[0].Normalize() != "mypriv" {
			t.Errorf("Privileges = %v, want [mypriv]", stmt.Privileges)
		}
	})

	t.Run("grant_privilege_named_to", func(t *testing.T) {
		// D1 edge: TO is non-reserved, so a privilege literally named "TO" is
		// valid; the role-vs-privilege scan must not mistake this first TO for
		// the role-form terminator.
		stmt := dclParseOne(t, "GRANT TO ON t TO alice").(*GrantPrivStmt)
		if len(stmt.Privileges) != 1 || stmt.Privileges[0].Normalize() != "to" {
			t.Errorf("Privileges = %v, want [to]", stmt.Privileges)
		}
	})

	t.Run("grant_bare_user_principal", func(t *testing.T) {
		// USER is non-reserved: a bare USER not followed by a name is itself the
		// (unspecified) grantee name.
		stmt := dclParseOne(t, "GRANT SELECT ON t TO USER").(*GrantPrivStmt)
		if stmt.Grantee.Kind != PrincipalUnspecified || stmt.Grantee.Name.Normalize() != "user" {
			t.Errorf("Grantee = %+v, want bare name 'user'", stmt.Grantee)
		}
	})

	t.Run("grant_qualified_target", func(t *testing.T) {
		stmt := dclParseOne(t, "GRANT SELECT ON cat.sch.tbl TO alice").(*GrantPrivStmt)
		if stmt.On.Object.Normalize() != "cat.sch.tbl" {
			t.Errorf("Object = %v, want cat.sch.tbl", stmt.On.Object)
		}
	})

	t.Run("grant_open_entity_kind", func(t *testing.T) {
		// D5: an arbitrary single word is a valid entity-kind qualifier.
		stmt := dclParseOne(t, "GRANT SELECT ON VIEW v TO alice").(*GrantPrivStmt)
		if stmt.On.Qualifier == nil || stmt.On.Qualifier.Normalize() != "view" {
			t.Errorf("Qualifier = %v, want view", stmt.On.Qualifier)
		}
		if stmt.On.Object.Normalize() != "v" {
			t.Errorf("Object = %v, want v", stmt.On.Object)
		}
	})

	t.Run("grant_object_named_branch", func(t *testing.T) {
		// `ON branch` is an ordinary object named "branch", NOT the BRANCH
		// clause (which requires `BRANCH <ident> IN`). The trailing TO is a
		// non-reserved keyword, which previously fooled the branch heuristic.
		stmt := dclParseOne(t, "GRANT SELECT ON branch TO alice").(*GrantPrivStmt)
		if stmt.On.Branch != nil {
			t.Errorf("Branch = %v, want nil (object named branch)", stmt.On.Branch)
		}
		if stmt.On.Object.Normalize() != "branch" {
			t.Errorf("Object = %v, want branch", stmt.On.Object)
		}
	})

	t.Run("grant_entity_qualifier_named_branch", func(t *testing.T) {
		// `ON branch orders` — "branch" is the entity-kind qualifier (no IN
		// follows), "orders" is the object.
		stmt := dclParseOne(t, "GRANT SELECT ON branch orders TO alice").(*GrantPrivStmt)
		if stmt.On.Branch != nil {
			t.Errorf("Branch = %v, want nil", stmt.On.Branch)
		}
		if stmt.On.Qualifier == nil || stmt.On.Qualifier.Normalize() != "branch" {
			t.Errorf("Qualifier = %v, want branch", stmt.On.Qualifier)
		}
		if stmt.On.Object.Normalize() != "orders" {
			t.Errorf("Object = %v, want orders", stmt.On.Object)
		}
	})

	t.Run("grant_branch_in_qualified", func(t *testing.T) {
		// `ON BRANCH x IN orders` — real branch form; the IN target may itself
		// carry an entity-kind qualifier.
		stmt := dclParseOne(t, "GRANT INSERT ON BRANCH audit IN TABLE orders TO alice").(*GrantPrivStmt)
		if stmt.On.Branch == nil || stmt.On.Branch.Normalize() != "audit" {
			t.Errorf("Branch = %v, want audit", stmt.On.Branch)
		}
		if !stmt.On.IsTable() || stmt.On.Object.Normalize() != "orders" {
			t.Errorf("IN target = qualifier %v object %v, want TABLE orders", stmt.On.Qualifier, stmt.On.Object)
		}
	})
}

func TestRevokePrivileges_Structure(t *testing.T) {
	t.Run("revoke_grant_option_for", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "REVOKE GRANT OPTION FOR SELECT ON nation FROM alice").(*RevokePrivStmt)
		if !ok {
			t.Fatalf("got %T, want *RevokePrivStmt", stmt)
		}
		if !stmt.GrantOptionFor {
			t.Errorf("GrantOptionFor = false, want true")
		}
		if len(stmt.Privileges) != 1 {
			t.Errorf("got %d privileges, want 1", len(stmt.Privileges))
		}
	})

	t.Run("revoke_all_privileges", func(t *testing.T) {
		stmt := dclParseOne(t, "REVOKE ALL PRIVILEGES ON test FROM alice").(*RevokePrivStmt)
		if !stmt.AllPrivileges {
			t.Errorf("AllPrivileges = false, want true")
		}
	})

	t.Run("revoke_on_branch", func(t *testing.T) {
		stmt := dclParseOne(t, "REVOKE INSERT ON BRANCH audit IN orders FROM alice").(*RevokePrivStmt)
		if stmt.On.Branch == nil || stmt.On.Branch.Normalize() != "audit" {
			t.Errorf("Branch = %v, want audit", stmt.On.Branch)
		}
	})
}

func TestDeny_Structure(t *testing.T) {
	t.Run("deny_privs_on_table", func(t *testing.T) {
		stmt, ok := dclParseOne(t, "DENY INSERT, SELECT ON orders TO alice").(*DenyStmt)
		if !ok {
			t.Fatalf("got %T, want *DenyStmt", stmt)
		}
		if len(stmt.Privileges) != 2 {
			t.Errorf("got %d privileges, want 2", len(stmt.Privileges))
		}
	})

	t.Run("deny_on_schema_to_role_public", func(t *testing.T) {
		stmt := dclParseOne(t, "DENY SELECT ON orders TO ROLE PUBLIC").(*DenyStmt)
		if stmt.Grantee.Kind != PrincipalRole {
			t.Errorf("Grantee kind = %v, want ROLE", stmt.Grantee.Kind)
		}
	})

	t.Run("deny_on_branch", func(t *testing.T) {
		stmt := dclParseOne(t, "DENY INSERT ON BRANCH audit IN orders TO alice").(*DenyStmt)
		if stmt.On.Branch == nil || stmt.On.Branch.Normalize() != "audit" {
			t.Errorf("Branch = %v, want audit", stmt.On.Branch)
		}
	})
}

func TestDCL_Negative(t *testing.T) {
	// Oracle-confirmed (SYNTAX_ERROR in Trino 481) rejections:
	negatives := []string{
		"GRANT SELECT TO foo",                                   // privilege keyword with no ON, not a valid role
		"GRANT select TO foo",                                   // reserved word can't be a role
		"GRANT CREATE TO foo",                                   // reserved CREATE can't be a role
		"GRANT SELECT, INSERT TO foo",                           // reserved privilege list with no ON
		"GRANT ALL PRIVILEGES TO foo",                           // ALL PRIVILEGES is privilege-only (needs ON)
		"GRANT SELECT TO alice IN hive",                         // IN catalog only valid for role grant
		"GRANT bar TO foo WITH GRANT OPTION",                    // WITH GRANT OPTION is privilege-only
		"GRANT SELECT ON test TO alice WITH ADMIN OPTION",       // WITH ADMIN OPTION is role-only
		"GRANT SELECT ON test TO alice GRANTED BY CURRENT_USER", // GRANTED BY is role-only
		"REVOKE SELECT ON test FROM alice CASCADE",              // Trino has no CASCADE/RESTRICT
		"GRANT CREATE TABLE ON SCHEMA s TO alice",               // privileges are single tokens
		"CREATE ROLE",                                           // missing name
		"DROP ROLE",                                             // missing name
	}
	for _, sql := range negatives {
		dclParseErr(t, sql)
	}
}

package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// TestGrantRoleOptions covers the trailing options of GRANT [DATABASE] ROLE:
// WITH GRANT OPTION and the expiry of a role granted to a user.
// https://docs.snowflake.com/en/sql-reference/sql/grant-role
// https://docs.snowflake.com/en/sql-reference/sql/grant-database-role
func TestGrantRoleOptions(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT ROLE analyst TO ROLE data_steward WITH GRANT OPTION")
	if len(errs) > 0 || !stmt.GrantOption {
		t.Fatalf("WITH GRANT OPTION: errs=%v stmt=%+v", errs, stmt)
	}
	stmt, errs = testParseGrantStmt(t, "GRANT DATABASE ROLE mydb.analyst TO DATABASE ROLE mydb.dr2 WITH GRANT OPTION")
	if len(errs) > 0 || !stmt.GrantOption || stmt.RoleKind != ast.GrantedDatabaseRole {
		t.Fatalf("database role WITH GRANT OPTION: errs=%v stmt=%+v", errs, stmt)
	}
	stmt, errs = testParseGrantStmt(t, "GRANT ROLE analyst TO USER alice FOR INTERVAL '48 HOURS' RESET EXPIRATION")
	if len(errs) > 0 || stmt.ExpireInterval != "48 HOURS" || stmt.ExpireUntil != "" || !stmt.ResetExpiration {
		t.Fatalf("FOR INTERVAL: errs=%v stmt=%+v", errs, stmt)
	}
	stmt, errs = testParseGrantStmt(t, "GRANT ROLE analyst TO USER alice UNTIL '2027-08-01T17:00:00Z'")
	if len(errs) > 0 || stmt.ExpireUntil != "2027-08-01T17:00:00Z" || stmt.ResetExpiration {
		t.Fatalf("UNTIL: errs=%v stmt=%+v", errs, stmt)
	}
	if stmt.Loc.End != len("GRANT ROLE analyst TO USER alice UNTIL '2027-08-01T17:00:00Z'") {
		t.Errorf("Loc.End = %d", stmt.Loc.End)
	}

	for _, sql := range []string{
		// WITH GRANT OPTION follows a grant to a role only.
		"GRANT ROLE analyst TO USER alice WITH GRANT OPTION",
		"GRANT DATABASE ROLE dr1 TO USER alice WITH GRANT OPTION",
		// The interval and the timestamp are quoted strings.
		"GRANT ROLE analyst TO USER alice FOR INTERVAL 1",
		"GRANT ROLE analyst TO USER alice UNTIL CURRENT_TIMESTAMP",
		// Expiry is for an account role granted to a user.
		"GRANT ROLE analyst TO ROLE r FOR INTERVAL '1 DAY'",
		"GRANT DATABASE ROLE dr1 TO USER alice FOR INTERVAL '1 DAY'",
		// RESET EXPIRATION follows FOR INTERVAL or UNTIL.
		"GRANT ROLE analyst TO USER alice RESET EXPIRATION",
		"GRANT ROLE analyst TO USER alice FOR INTERVAL '1 DAY' UNTIL '2027-01-01'",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

// TestRevokeRoleOptions covers GRANT OPTION FOR and RESTRICT | CASCADE on
// REVOKE [DATABASE] ROLE.
// https://docs.snowflake.com/en/sql-reference/sql/revoke-role
// https://docs.snowflake.com/en/sql-reference/sql/revoke-database-role
func TestRevokeRoleOptions(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE GRANT OPTION FOR ROLE analyst FROM ROLE data_steward")
	if len(errs) > 0 || !stmt.GrantOptionFor || stmt.RoleKind != ast.GrantedAccountRole {
		t.Fatalf("GRANT OPTION FOR ROLE: errs=%v stmt=%+v", errs, stmt)
	}
	stmt, errs = testParseRevokeStmt(t, "REVOKE ROLE analyst FROM ROLE data_steward CASCADE")
	if len(errs) > 0 || !stmt.Cascade {
		t.Fatalf("CASCADE: errs=%v stmt=%+v", errs, stmt)
	}
	stmt, errs = testParseRevokeStmt(t, "REVOKE GRANT OPTION FOR DATABASE ROLE mydb.analyst FROM DATABASE ROLE mydb.dr2 RESTRICT")
	if len(errs) > 0 || !stmt.GrantOptionFor || !stmt.Restrict || stmt.RoleKind != ast.GrantedDatabaseRole {
		t.Fatalf("database role: errs=%v stmt=%+v", errs, stmt)
	}
	// GRANT OPTION FOR still opens a privilege revoke.
	if _, errs := testParseRevokeStmt(t, "REVOKE GRANT OPTION FOR SELECT ON TABLE t FROM ROLE r"); len(errs) > 0 {
		t.Fatalf("privilege revoke: %v", errs)
	}

	for _, sql := range []string{
		"REVOKE GRANT OPTION FOR ROLE analyst FROM USER alice",
		"REVOKE ROLE analyst FROM USER alice CASCADE",
		"REVOKE GRANT OPTION FOR DATABASE ROLE dr1 FROM APPLICATION app",
		"REVOKE DATABASE ROLE dr1 FROM APPLICATION app RESTRICT",
		"REVOKE GRANT OPTION FOR SERVICE ROLE svc!r FROM ROLE r",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

// TestGrantServiceRole covers GRANT and REVOKE SERVICE ROLE, whose name is
// written <service>!<role>.
// https://docs.snowflake.com/en/sql-reference/sql/grant-service-role
// https://docs.snowflake.com/en/sql-reference/sql/revoke-service-role
func TestGrantServiceRole(t *testing.T) {
	sql := "GRANT SERVICE ROLE db.s.echo_service!echoendpoint_role TO ROLE service_function_user_role"
	g, errs := testParseGrantStmt(t, sql)
	if len(errs) > 0 {
		t.Fatalf("errs=%v", errs)
	}
	if g.Kind != ast.GrantRole || g.RoleKind != ast.GrantedServiceRole ||
		g.Role.Normalize() != "DB.S.ECHO_SERVICE" || g.ServiceRole.Normalize() != "ECHOENDPOINT_ROLE" ||
		g.Grantee.Kind != ast.GranteeRole || g.Loc.End != len(sql) {
		t.Fatalf("stmt=%+v", g)
	}
	for _, ok := range []string{
		"GRANT SERVICE ROLE svc!r TO APPLICATION ROLE app_role",
		"GRANT SERVICE ROLE svc!r TO DATABASE ROLE db.dr",
	} {
		if _, err := Parse(ok); err != nil {
			t.Errorf("Parse(%q): %v", ok, err)
		}
	}
	r, errs := testParseRevokeStmt(t, "REVOKE SERVICE ROLE echo_service!echoendpoint_role FROM DATABASE ROLE db.dr")
	if len(errs) > 0 || r.RoleKind != ast.GrantedServiceRole || r.ServiceRole.Normalize() != "ECHOENDPOINT_ROLE" ||
		r.Grantee.Kind != ast.GranteeDatabaseRole {
		t.Fatalf("revoke: errs=%v stmt=%+v", errs, r)
	}

	for _, sql := range []string{
		// The page requires service-name!service-role-name.
		"GRANT SERVICE ROLE svc TO ROLE r",
		"REVOKE SERVICE ROLE svc FROM ROLE r",
		"GRANT SERVICE ROLE svc! TO ROLE r",
		// The grantee is a ROLE, APPLICATION ROLE or DATABASE ROLE.
		"GRANT SERVICE ROLE svc!r TO r2",
		"GRANT SERVICE ROLE svc!r TO USER u",
		"GRANT SERVICE ROLE svc!r TO APPLICATION app",
		"REVOKE SERVICE ROLE svc!r FROM SHARE s",
		"GRANT SERVICE ROLE svc!r TO ROLE r WITH GRANT OPTION",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

// TestGrantInheritedCallerContainers covers the IN list of INHERITED CALLER
// grants, which also names the account, an application or an application
// package.
// https://docs.snowflake.com/en/sql-reference/sql/grant-caller
// https://docs.snowflake.com/en/sql-reference/sql/revoke-caller
func TestGrantInheritedCallerContainers(t *testing.T) {
	g, errs := testParseGrantStmt(t, "GRANT ALL INHERITED CALLER PRIVILEGES ON ALL SCHEMAS IN ACCOUNT TO ROLE owner_role")
	if len(errs) > 0 || g.On.Container != ast.GrantContainerAccount || g.On.ContainerName != nil {
		t.Fatalf("IN ACCOUNT: errs=%v on=%+v", errs, g.On)
	}
	g, errs = testParseGrantStmt(t, "GRANT INHERITED CALLER SELECT, INSERT ON ALL TABLES IN APPLICATION app TO ROLE r")
	if len(errs) > 0 || g.On.Container != ast.GrantContainerApplication || g.On.ContainerName.Normalize() != "APP" {
		t.Fatalf("IN APPLICATION: errs=%v on=%+v", errs, g.On)
	}
	r, errs := testParseRevokeStmt(t, "REVOKE ALL INHERITED CALLER PRIVILEGES ON ALL VIEWS IN APPLICATION PACKAGE pkg FROM ROLE r")
	if len(errs) > 0 || r.On.Container != ast.GrantContainerApplicationPackage || r.On.ContainerName.Normalize() != "PKG" {
		t.Fatalf("IN APPLICATION PACKAGE: errs=%v on=%+v", errs, r.On)
	}

	for _, sql := range []string{
		// Only INHERITED CALLER grants take the wider IN list.
		"GRANT SELECT ON ALL TABLES IN ACCOUNT TO ROLE r",
		"GRANT ALL CALLER PRIVILEGES ON ALL TABLES IN ACCOUNT TO ROLE r",
		"GRANT INHERITED CALLER PRIVILEGES ON ALL TABLES IN ACCOUNT TO ROLE r",
		"REVOKE SELECT ON ALL TABLES IN APPLICATION app FROM ROLE r",
		"GRANT OWNERSHIP ON ALL TABLES IN APPLICATION app TO ROLE r",
		// ... and only on ON ALL.
		"GRANT INHERITED CALLER SELECT ON FUTURE TABLES IN ACCOUNT TO ROLE r",
		// IN ACCOUNT takes no name; APPLICATION [PACKAGE] takes one part.
		"GRANT ALL INHERITED CALLER PRIVILEGES ON ALL TABLES IN ACCOUNT a TO ROLE r",
		"GRANT ALL INHERITED CALLER PRIVILEGES ON ALL TABLES IN APPLICATION PACKAGE TO ROLE r",
		"GRANT ALL INHERITED CALLER PRIVILEGES ON ALL TABLES IN APPLICATION db.app TO ROLE r",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

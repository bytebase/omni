package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func testParseGrantStmt(t *testing.T, input string) (*ast.GrantStmt, []ParseError) {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.GrantStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a GrantStmt"})
	}
	return stmt, result.Errors
}

func testParseRevokeStmt(t *testing.T, input string) (*ast.RevokeStmt, []ParseError) {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.RevokeStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a RevokeStmt"})
	}
	return stmt, result.Errors
}

func privNames(privs []*ast.Privilege) []string {
	out := make([]string, len(privs))
	for i, p := range privs {
		out[i] = p.Name
	}
	return out
}

// sigNames renders a function/procedure signature's argument data types to
// their uppercased source names for comparison.
func sigNames(sig []*ast.TypeName) []string {
	out := make([]string, len(sig))
	for i, t := range sig {
		out[i] = strings.ToUpper(t.Name)
	}
	return out
}

func eqStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// GRANT ROLE / GRANT DATABASE ROLE
// ---------------------------------------------------------------------------

func TestGrantRole_ToRole(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT ROLE analyst TO ROLE SYSADMIN")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.GrantRole {
		t.Errorf("kind = %v, want GrantRole", stmt.Kind)
	}
	if stmt.RoleKind != ast.GrantedAccountRole {
		t.Errorf("role kind = %v, want GrantedAccountRole", stmt.RoleKind)
	}
	if stmt.Role.Normalize() != "ANALYST" {
		t.Errorf("role = %q, want ANALYST", stmt.Role.Normalize())
	}
	if stmt.Grantee.Kind != ast.GranteeRole {
		t.Errorf("grantee kind = %v, want GranteeRole", stmt.Grantee.Kind)
	}
	if stmt.Grantee.Name.Normalize() != "SYSADMIN" {
		t.Errorf("grantee = %q, want SYSADMIN", stmt.Grantee.Name.Normalize())
	}
}

func TestGrantRole_ToUser(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT ROLE analyst TO USER user1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeUser {
		t.Errorf("grantee kind = %v, want GranteeUser", stmt.Grantee.Kind)
	}
	if stmt.Grantee.Name.Normalize() != "USER1" {
		t.Errorf("grantee = %q, want USER1", stmt.Grantee.Name.Normalize())
	}
}

func TestGrantDatabaseRole_ToDatabaseRole(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT DATABASE ROLE mydb.dr1 TO DATABASE ROLE mydb.dr2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.GrantRole {
		t.Errorf("kind = %v, want GrantRole", stmt.Kind)
	}
	if stmt.RoleKind != ast.GrantedDatabaseRole {
		t.Errorf("role kind = %v, want GrantedDatabaseRole", stmt.RoleKind)
	}
	if stmt.Role.Normalize() != "MYDB.DR1" {
		t.Errorf("role = %q, want MYDB.DR1", stmt.Role.Normalize())
	}
	if stmt.Grantee.Kind != ast.GranteeDatabaseRole {
		t.Errorf("grantee kind = %v, want GranteeDatabaseRole", stmt.Grantee.Kind)
	}
	if stmt.Grantee.Name.Normalize() != "MYDB.DR2" {
		t.Errorf("grantee = %q, want MYDB.DR2", stmt.Grantee.Name.Normalize())
	}
}

func TestGrantDatabaseRole_ToUser(t *testing.T) {
	// From official corpus grant-privilege/example_20.sql.
	stmt, errs := testParseGrantStmt(t, "GRANT DATABASE ROLE mydb.dr1 TO USER testuser")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.RoleKind != ast.GrantedDatabaseRole {
		t.Errorf("role kind = %v, want GrantedDatabaseRole", stmt.RoleKind)
	}
	if stmt.Grantee.Kind != ast.GranteeUser {
		t.Errorf("grantee kind = %v, want GranteeUser", stmt.Grantee.Kind)
	}
}

func TestGrantDatabaseRole_ToApplication(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT DATABASE ROLE mydb.dr1 TO APPLICATION app1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeApplication {
		t.Errorf("grantee kind = %v, want GranteeApplication", stmt.Grantee.Kind)
	}
}

func TestGrantApplicationRole_ToUser(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT APPLICATION ROLE app1.ar1 TO USER user1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.GrantRole {
		t.Errorf("kind = %v, want GrantRole", stmt.Kind)
	}
	if stmt.RoleKind != ast.GrantedApplicationRole {
		t.Errorf("role kind = %v, want GrantedApplicationRole", stmt.RoleKind)
	}
	if stmt.Role.Normalize() != "APP1.AR1" {
		t.Errorf("role = %q, want APP1.AR1", stmt.Role.Normalize())
	}
	if stmt.Grantee.Kind != ast.GranteeUser {
		t.Errorf("grantee kind = %v, want GranteeUser", stmt.Grantee.Kind)
	}
}

func TestGrantApplicationRole_ToApplicationRole(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT APPLICATION ROLE app1.ar1 TO APPLICATION ROLE app2.ar2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.RoleKind != ast.GrantedApplicationRole {
		t.Errorf("role kind = %v, want GrantedApplicationRole", stmt.RoleKind)
	}
	if stmt.Grantee.Kind != ast.GranteeApplicationRole {
		t.Errorf("grantee kind = %v, want GranteeApplicationRole", stmt.Grantee.Kind)
	}
}

func TestRevokeApplicationRole_FromApplicationRole(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE APPLICATION ROLE app1.ar1 FROM APPLICATION ROLE app2.ar2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.RoleKind != ast.GrantedApplicationRole {
		t.Errorf("role kind = %v, want GrantedApplicationRole", stmt.RoleKind)
	}
	if stmt.Grantee.Kind != ast.GranteeApplicationRole {
		t.Errorf("grantee kind = %v, want GranteeApplicationRole", stmt.Grantee.Kind)
	}
}

func TestRevokeRole_FromRole(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE ROLE public FROM ROLE public")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.RevokeRole {
		t.Errorf("kind = %v, want RevokeRole", stmt.Kind)
	}
	if stmt.Role.Normalize() != "PUBLIC" {
		t.Errorf("role = %q, want PUBLIC", stmt.Role.Normalize())
	}
	if stmt.Grantee.Kind != ast.GranteeRole {
		t.Errorf("grantee kind = %v, want GranteeRole", stmt.Grantee.Kind)
	}
}

func TestRevokeDatabaseRole_FromUser(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE DATABASE ROLE mydb.dr1 FROM USER bob")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.RoleKind != ast.GrantedDatabaseRole {
		t.Errorf("role kind = %v, want GrantedDatabaseRole", stmt.RoleKind)
	}
	if stmt.Grantee.Kind != ast.GranteeUser {
		t.Errorf("grantee kind = %v, want GranteeUser", stmt.Grantee.Kind)
	}
}

// ---------------------------------------------------------------------------
// GRANT <privileges> ON ...
// ---------------------------------------------------------------------------

func TestGrantPriv_OnWarehouse(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OPERATE ON WAREHOUSE report_wh TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.GrantPrivileges {
		t.Errorf("kind = %v, want GrantPrivileges", stmt.Kind)
	}
	if stmt.AllPrivileges {
		t.Error("expected AllPrivileges=false")
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"OPERATE"}) {
		t.Errorf("privileges = %v, want [OPERATE]", got)
	}
	if stmt.On.Kind != ast.GrantTargetObject {
		t.Errorf("target kind = %v, want GrantTargetObject", stmt.On.Kind)
	}
	if stmt.On.ObjectType != "WAREHOUSE" {
		t.Errorf("object type = %q, want WAREHOUSE", stmt.On.ObjectType)
	}
	if stmt.On.Name.Normalize() != "REPORT_WH" {
		t.Errorf("object name = %q, want REPORT_WH", stmt.On.Name.Normalize())
	}
	if stmt.Grantee.Kind != ast.GranteeRole || stmt.Grantee.Name.Normalize() != "ANALYST" {
		t.Errorf("grantee = %v %q", stmt.Grantee.Kind, stmt.Grantee.Name.Normalize())
	}
}

func TestGrantPriv_WithGrantOption(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OPERATE ON WAREHOUSE report_wh TO ROLE analyst WITH GRANT OPTION")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.GrantOption {
		t.Error("expected GrantOption=true")
	}
}

func TestGrantPriv_MultiPrivilegeList(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT, INSERT ON FUTURE TABLES IN SCHEMA mydb.myschema TO ROLE role1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"SELECT", "INSERT"}) {
		t.Errorf("privileges = %v, want [SELECT INSERT]", got)
	}
	if stmt.On.Kind != ast.GrantTargetFutureIn {
		t.Errorf("target kind = %v, want GrantTargetFutureIn", stmt.On.Kind)
	}
	if stmt.On.ObjectTypePlural != "TABLES" {
		t.Errorf("plural = %q, want TABLES", stmt.On.ObjectTypePlural)
	}
	if stmt.On.Container != ast.GrantContainerSchema {
		t.Errorf("container = %v, want GrantContainerSchema", stmt.On.Container)
	}
	if stmt.On.ContainerName.Normalize() != "MYDB.MYSCHEMA" {
		t.Errorf("container name = %q, want MYDB.MYSCHEMA", stmt.On.ContainerName.Normalize())
	}
}

func TestGrantPriv_AllPrivilegesOnDatabase(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT ALL PRIVILEGES ON DATABASE mydb TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.AllPrivileges {
		t.Error("expected AllPrivileges=true")
	}
	if len(stmt.Privileges) != 0 {
		t.Errorf("expected no explicit privileges, got %v", privNames(stmt.Privileges))
	}
	if stmt.On.ObjectType != "DATABASE" || stmt.On.Name.Normalize() != "MYDB" {
		t.Errorf("target = %q %q", stmt.On.ObjectType, stmt.On.Name.Normalize())
	}
}

func TestGrantPriv_AllBareNoPrivileges(t *testing.T) {
	// ALL without the PRIVILEGES keyword is also valid.
	stmt, errs := testParseGrantStmt(t, "GRANT ALL ON DATABASE mydb TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.AllPrivileges {
		t.Error("expected AllPrivileges=true")
	}
}

func TestGrantPriv_OnAccount(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT CREATE WAREHOUSE ON ACCOUNT TO ROLE myrole")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Kind != ast.GrantTargetAccount {
		t.Errorf("target kind = %v, want GrantTargetAccount", stmt.On.Kind)
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"CREATE WAREHOUSE"}) {
		t.Errorf("privileges = %v, want [CREATE WAREHOUSE]", got)
	}
}

func TestGrantPriv_MultiWordPrivilege(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT CREATE MATERIALIZED VIEW ON SCHEMA mydb.myschema TO ROLE myrole")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"CREATE MATERIALIZED VIEW"}) {
		t.Errorf("privileges = %v, want [CREATE MATERIALIZED VIEW]", got)
	}
	if stmt.On.ObjectType != "SCHEMA" {
		t.Errorf("object type = %q, want SCHEMA", stmt.On.ObjectType)
	}
}

func TestGrantPriv_OnFunctionWithSignature(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT ALL PRIVILEGES ON FUNCTION mydb.myschema.add5(number) TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.ObjectType != "FUNCTION" {
		t.Errorf("object type = %q, want FUNCTION", stmt.On.ObjectType)
	}
	if stmt.On.Name.Normalize() != "MYDB.MYSCHEMA.ADD5" {
		t.Errorf("name = %q, want MYDB.MYSCHEMA.ADD5", stmt.On.Name.Normalize())
	}
	if !eqStrings(sigNames(stmt.On.Signature), []string{"NUMBER"}) {
		t.Errorf("signature = %v, want [NUMBER]", sigNames(stmt.On.Signature))
	}
}

func TestGrantPriv_OnProcedureMultiArgSignature(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT ALL PRIVILEGES ON PROCEDURE clean_schema(string, string) TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.ObjectType != "PROCEDURE" {
		t.Errorf("object type = %q, want PROCEDURE", stmt.On.ObjectType)
	}
	if !eqStrings(sigNames(stmt.On.Signature), []string{"STRING", "STRING"}) {
		t.Errorf("signature = %v, want [STRING STRING]", sigNames(stmt.On.Signature))
	}
}

func TestGrantPriv_ParameterizedSignature(t *testing.T) {
	// Function arguments are full data types, including parameterized
	// (NUMBER(38,0)) and multi-word (DOUBLE PRECISION) forms.
	stmt, errs := testParseGrantStmt(t, "GRANT USAGE ON FUNCTION mydb.s.f(NUMBER(38,0), VARCHAR(100), DOUBLE PRECISION) TO ROLE r")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := sigNames(stmt.On.Signature); !eqStrings(got, []string{"NUMBER", "VARCHAR", "DOUBLE PRECISION"}) {
		// parseDataType captures the canonical type name (numeric params live
		// in TypeName.Params; DOUBLE PRECISION keeps its two-word Name).
		t.Errorf("signature names = %v, want [NUMBER VARCHAR DOUBLE PRECISION]", got)
	}
	if len(stmt.On.Signature) != 3 {
		t.Fatalf("expected 3 signature args, got %d", len(stmt.On.Signature))
	}
	// The first arg retains its numeric parameters.
	if got := stmt.On.Signature[0].Params; len(got) != 2 || got[0] != 38 || got[1] != 0 {
		t.Errorf("first arg params = %v, want [38 0]", got)
	}
}

func TestGrantPriv_AllTablesInSchema(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT ON ALL TABLES IN SCHEMA mydb.myschema TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Kind != ast.GrantTargetAllIn {
		t.Errorf("target kind = %v, want GrantTargetAllIn", stmt.On.Kind)
	}
	if stmt.On.ObjectTypePlural != "TABLES" {
		t.Errorf("plural = %q, want TABLES", stmt.On.ObjectTypePlural)
	}
	if stmt.On.Container != ast.GrantContainerSchema {
		t.Errorf("container = %v, want GrantContainerSchema", stmt.On.Container)
	}
}

func TestGrantPriv_FutureSchemasInDatabase(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT USAGE ON FUTURE SCHEMAS IN DATABASE mydb TO ROLE role1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Kind != ast.GrantTargetFutureIn {
		t.Errorf("target kind = %v, want GrantTargetFutureIn", stmt.On.Kind)
	}
	if stmt.On.ObjectTypePlural != "SCHEMAS" {
		t.Errorf("plural = %q, want SCHEMAS", stmt.On.ObjectTypePlural)
	}
	if stmt.On.Container != ast.GrantContainerDatabase {
		t.Errorf("container = %v, want GrantContainerDatabase", stmt.On.Container)
	}
	if stmt.On.ContainerName.Normalize() != "MYDB" {
		t.Errorf("container name = %q, want MYDB", stmt.On.ContainerName.Normalize())
	}
}

func TestGrantPriv_ToDatabaseRoleGrantee(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT ON ALL TABLES IN SCHEMA mydb.myschema TO DATABASE ROLE mydb.dr1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeDatabaseRole {
		t.Errorf("grantee kind = %v, want GranteeDatabaseRole", stmt.Grantee.Kind)
	}
	if stmt.Grantee.Name.Normalize() != "MYDB.DR1" {
		t.Errorf("grantee = %q, want MYDB.DR1", stmt.Grantee.Name.Normalize())
	}
}

func TestGrantPriv_BareRoleGrantee(t *testing.T) {
	// TO <role_name> with the ROLE keyword omitted.
	stmt, errs := testParseGrantStmt(t, "GRANT USAGE ON DATABASE mydb TO myrole")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeRole {
		t.Errorf("grantee kind = %v, want GranteeRole (bare)", stmt.Grantee.Kind)
	}
	if stmt.Grantee.Name.Normalize() != "MYROLE" {
		t.Errorf("grantee = %q, want MYROLE", stmt.Grantee.Name.Normalize())
	}
}

func TestGrantPriv_IdentifierObjectType(t *testing.T) {
	// NOTEBOOK is NOT a keyword in the lexer — it must be accepted as an
	// identifier object type (docs win over the stale legacy enum).
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON NOTEBOOK db_one.schema_one.mynotebook TO ROLE finance")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.ObjectType != "NOTEBOOK" {
		t.Errorf("object type = %q, want NOTEBOOK", stmt.On.ObjectType)
	}
}

func TestGrantPriv_MultiWordIdentifierPrivilege(t *testing.T) {
	// CREATE PROVISIONED THROUGHPUT — PROVISIONED/THROUGHPUT are identifiers.
	stmt, errs := testParseGrantStmt(t, "GRANT CREATE PROVISIONED THROUGHPUT ON ACCOUNT TO ROLE myrole")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"CREATE PROVISIONED THROUGHPUT"}) {
		t.Errorf("privileges = %v, want [CREATE PROVISIONED THROUGHPUT]", got)
	}
}

func TestGrantPriv_DottedIdentifierPrivilege(t *testing.T) {
	// CREATE SNOWFLAKE.CORE.BUDGET — a dotted class privilege.
	stmt, errs := testParseGrantStmt(t, "GRANT CREATE SNOWFLAKE.CORE.BUDGET ON SCHEMA budgets_db.budgets_schema TO ROLE budget_admin")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"CREATE SNOWFLAKE.CORE.BUDGET"}) {
		t.Errorf("privileges = %v, want [CREATE SNOWFLAKE.CORE.BUDGET]", got)
	}
}

// TestGrantPriv_MultiWordObjectTypes exercises the full vocabulary of
// multi-word object types, including those whose connector words the lexer
// emits as identifiers (COMPUTE POOL, EXTERNAL VOLUME, NETWORK RULE). The
// object name must never be swallowed into the type.
func TestGrantPriv_MultiWordObjectTypes(t *testing.T) {
	cases := []struct {
		sql      string
		wantType string
		wantName string
	}{
		{"GRANT MONITOR ON RESOURCE MONITOR rm TO ROLE r", "RESOURCE MONITOR", "RM"},
		{"GRANT SELECT ON EXTERNAL TABLE t TO ROLE r", "EXTERNAL TABLE", "T"},
		{"GRANT SELECT ON MATERIALIZED VIEW mv TO ROLE r", "MATERIALIZED VIEW", "MV"},
		{"GRANT SELECT ON DYNAMIC TABLE dt TO ROLE r", "DYNAMIC TABLE", "DT"},
		{"GRANT SELECT ON ICEBERG TABLE it TO ROLE r", "ICEBERG TABLE", "IT"},
		{"GRANT APPLY ON MASKING POLICY p TO ROLE r", "MASKING POLICY", "P"},
		{"GRANT APPLY ON ROW ACCESS POLICY p TO ROLE r", "ROW ACCESS POLICY", "P"},
		{"GRANT USAGE ON FILE FORMAT ff TO ROLE r", "FILE FORMAT", "FF"},
		{"GRANT OPERATE ON COMPUTE POOL cp TO ROLE r", "COMPUTE POOL", "CP"},
		{"GRANT USAGE ON EXTERNAL VOLUME v TO ROLE r", "EXTERNAL VOLUME", "V"},
		{"GRANT USAGE ON NETWORK RULE nr TO ROLE r", "NETWORK RULE", "NR"},
		{"GRANT MODIFY ON FAILOVER GROUP fg TO ROLE r", "FAILOVER GROUP", "FG"},
		{"GRANT MODIFY ON REPLICATION GROUP rg TO ROLE r", "REPLICATION GROUP", "RG"},
		{"GRANT USAGE ON SEMANTIC VIEW sv TO SHARE s", "SEMANTIC VIEW", "SV"},
		// Three-word object types and arbitrary new types (open-ended, no fixed
		// vocabulary): the object name is always the final space-separated unit.
		{"GRANT USAGE ON CORTEX SEARCH SERVICE mydb.s.css TO ROLE r", "CORTEX SEARCH SERVICE", "MYDB.S.CSS"},
		{"GRANT ALL ON DATA METRIC FUNCTION mydb.s.dmf TO ROLE r", "DATA METRIC FUNCTION", "MYDB.S.DMF"},
		{"GRANT APPLY ON STORAGE LIFECYCLE POLICY p TO ROLE r", "STORAGE LIFECYCLE POLICY", "P"},
		{"GRANT USAGE ON GIT REPOSITORY repo TO ROLE r", "GIT REPOSITORY", "REPO"},
		{"GRANT READ ON IMAGE REPOSITORY ir TO ROLE r", "IMAGE REPOSITORY", "IR"},
		{"GRANT USAGE ON MODEL m TO ROLE r", "MODEL", "M"},
		{"GRANT USAGE ON STREAMLIT st TO ROLE r", "STREAMLIT", "ST"},
		{"GRANT READ ON SECRET sec TO ROLE r", "SECRET", "SEC"},
		{"GRANT USAGE ON WORKSPACE mydb.s.ws TO ROLE r", "WORKSPACE", "MYDB.S.WS"},
	}
	for _, c := range cases {
		t.Run(c.wantType, func(t *testing.T) {
			stmt, errs := testParseGrantStmt(t, c.sql)
			if len(errs) > 0 {
				t.Fatalf("unexpected errors: %v", errs)
			}
			if stmt.On.ObjectType != c.wantType {
				t.Errorf("object type = %q, want %q", stmt.On.ObjectType, c.wantType)
			}
			if stmt.On.Name.Normalize() != c.wantName {
				t.Errorf("object name = %q, want %q", stmt.On.Name.Normalize(), c.wantName)
			}
		})
	}
}

func TestGrantPriv_ImportedPrivileges(t *testing.T) {
	// IMPORTED PRIVILEGES is a two-word privilege (not the ALL PRIVILEGES form).
	stmt, errs := testParseGrantStmt(t, "GRANT IMPORTED PRIVILEGES ON DATABASE shared_db TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.AllPrivileges {
		t.Error("expected AllPrivileges=false for IMPORTED PRIVILEGES")
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"IMPORTED PRIVILEGES"}) {
		t.Errorf("privileges = %v, want [IMPORTED PRIVILEGES]", got)
	}
}

func TestGrantPriv_ReadWriteOnTag(t *testing.T) {
	// READ and WRITE as a comma-separated privilege list.
	stmt, errs := testParseGrantStmt(t, "GRANT READ, WRITE ON TAG mydb.s.t TO ROLE r")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"READ", "WRITE"}) {
		t.Errorf("privileges = %v, want [READ WRITE]", got)
	}
	if stmt.On.ObjectType != "TAG" {
		t.Errorf("object type = %q, want TAG", stmt.On.ObjectType)
	}
}

func TestGrantPriv_OnObjectKeywordName(t *testing.T) {
	// A reserved keyword used as an (unquoted) object name part: the object
	// type must terminate so the keyword becomes the name, not the type.
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT ON TABLE mydb.public.t TO ROLE r")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.ObjectType != "TABLE" {
		t.Errorf("object type = %q, want TABLE", stmt.On.ObjectType)
	}
	if stmt.On.Name.Normalize() != "MYDB.PUBLIC.T" {
		t.Errorf("name = %q, want MYDB.PUBLIC.T", stmt.On.Name.Normalize())
	}
}

func TestGrantPriv_QuotedConnectorWordName(t *testing.T) {
	// A quoted object name that spells a multi-word-type connector word
	// (e.g. "VIEW") must NOT be absorbed into the object type. Quoted
	// identifiers are always names.
	stmt, errs := testParseGrantStmt(t, `GRANT SELECT ON TABLE "VIEW" TO ROLE r`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.ObjectType != "TABLE" {
		t.Errorf("object type = %q, want TABLE", stmt.On.ObjectType)
	}
	if stmt.On.Name.String() != `"VIEW"` {
		t.Errorf("object name = %q, want \"VIEW\"", stmt.On.Name.String())
	}
}

func TestGrantPriv_NoSignatureLeavesNil(t *testing.T) {
	// A non-function object must not carry a signature.
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT ON TABLE t TO ROLE r")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Signature != nil {
		t.Errorf("signature = %v, want nil", stmt.On.Signature)
	}
}

func TestGrantPriv_EmptySignature(t *testing.T) {
	// FUNCTION foo() — empty arg list, non-nil empty signature distinguishes
	// from "no parens".
	stmt, errs := testParseGrantStmt(t, "GRANT USAGE ON FUNCTION mydb.s.f() TO ROLE r")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Signature == nil {
		t.Fatal("signature = nil, want non-nil empty slice")
	}
	if len(stmt.On.Signature) != 0 {
		t.Errorf("signature = %v, want empty", stmt.On.Signature)
	}
}

// ---------------------------------------------------------------------------
// GRANT OWNERSHIP
// ---------------------------------------------------------------------------

func TestGrantOwnership_OnObject(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON DATABASE mydb TO ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.GrantOwnership {
		t.Errorf("kind = %v, want GrantOwnership", stmt.Kind)
	}
	if stmt.On.ObjectType != "DATABASE" || stmt.On.Name.Normalize() != "MYDB" {
		t.Errorf("target = %q %q", stmt.On.ObjectType, stmt.On.Name.Normalize())
	}
	if stmt.CurrentGrants != ast.CurrentGrantsNone {
		t.Errorf("current grants = %v, want None", stmt.CurrentGrants)
	}
}

func TestGrantOwnership_CopyCurrentGrants(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON TABLE mydb.public.mytable TO ROLE analyst COPY CURRENT GRANTS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.CurrentGrants != ast.CurrentGrantsCopy {
		t.Errorf("current grants = %v, want Copy", stmt.CurrentGrants)
	}
}

func TestGrantOwnership_AllTablesCopyCurrentGrants(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON ALL TABLES IN SCHEMA mydb.public TO ROLE analyst COPY CURRENT GRANTS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Kind != ast.GrantTargetAllIn {
		t.Errorf("target kind = %v, want GrantTargetAllIn", stmt.On.Kind)
	}
	if stmt.CurrentGrants != ast.CurrentGrantsCopy {
		t.Errorf("current grants = %v, want Copy", stmt.CurrentGrants)
	}
}

func TestGrantOwnership_RevokeCurrentGrants(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON DATABASE mydb TO ROLE r2 REVOKE CURRENT GRANTS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.CurrentGrants != ast.CurrentGrantsRevoke {
		t.Errorf("current grants = %v, want Revoke", stmt.CurrentGrants)
	}
}

func TestGrantOwnership_ToDatabaseRole(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON ALL TABLES IN SCHEMA mydb.public TO DATABASE ROLE mydb.dr1 COPY CURRENT GRANTS")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeDatabaseRole {
		t.Errorf("grantee kind = %v, want GranteeDatabaseRole", stmt.Grantee.Kind)
	}
}

func TestGrantOwnership_FutureInDatabase(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT OWNERSHIP ON FUTURE TABLES IN DATABASE mydb TO ROLE r2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Kind != ast.GrantTargetFutureIn {
		t.Errorf("target kind = %v, want GrantTargetFutureIn", stmt.On.Kind)
	}
}

// ---------------------------------------------------------------------------
// GRANT / REVOKE ... TO/FROM SHARE
// ---------------------------------------------------------------------------

func TestGrantPriv_ToShare(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT USAGE ON DATABASE mydb TO SHARE myshare")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeShare {
		t.Errorf("grantee kind = %v, want GranteeShare", stmt.Grantee.Kind)
	}
	if stmt.Grantee.Name.Normalize() != "MYSHARE" {
		t.Errorf("grantee = %q, want MYSHARE", stmt.Grantee.Name.Normalize())
	}
}

func TestGrantPriv_ToShareAllTablesInSchema(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT ON ALL TABLES IN SCHEMA mydb.public TO SHARE myshare")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.On.Kind != ast.GrantTargetAllIn {
		t.Errorf("target kind = %v, want GrantTargetAllIn", stmt.On.Kind)
	}
	if stmt.Grantee.Kind != ast.GranteeShare {
		t.Errorf("grantee kind = %v, want GranteeShare", stmt.Grantee.Kind)
	}
}

func TestRevokePriv_FromShare(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE USAGE ON DATABASE mydb FROM SHARE myshare")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.RevokePrivileges {
		t.Errorf("kind = %v, want RevokePrivileges", stmt.Kind)
	}
	if stmt.Grantee.Kind != ast.GranteeShare {
		t.Errorf("grantee kind = %v, want GranteeShare", stmt.Grantee.Kind)
	}
}

// ---------------------------------------------------------------------------
// REVOKE <privileges> ON ...
// ---------------------------------------------------------------------------

func TestRevokePriv_OnAccount(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE CREATE WAREHOUSE ON ACCOUNT FROM ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.RevokePrivileges {
		t.Errorf("kind = %v, want RevokePrivileges", stmt.Kind)
	}
	if stmt.On.Kind != ast.GrantTargetAccount {
		t.Errorf("target kind = %v, want GrantTargetAccount", stmt.On.Kind)
	}
	if stmt.GrantOptionFor {
		t.Error("expected GrantOptionFor=false")
	}
}

func TestRevokePriv_GrantOptionFor(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE GRANT OPTION FOR OPERATE ON WAREHOUSE report_wh FROM ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.GrantOptionFor {
		t.Error("expected GrantOptionFor=true")
	}
	if got := privNames(stmt.Privileges); !eqStrings(got, []string{"OPERATE"}) {
		t.Errorf("privileges = %v, want [OPERATE]", got)
	}
}

func TestRevokePriv_Cascade(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE OPERATE ON WAREHOUSE report_wh FROM ROLE analyst CASCADE")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Cascade {
		t.Error("expected Cascade=true")
	}
	if stmt.Restrict {
		t.Error("expected Restrict=false")
	}
}

func TestRevokePriv_Restrict(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE OPERATE ON WAREHOUSE report_wh FROM ROLE analyst RESTRICT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Restrict {
		t.Error("expected Restrict=true")
	}
}

func TestRevokePriv_AllPrivilegesOnFunction(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE ALL PRIVILEGES ON FUNCTION add5(number) FROM ROLE analyst")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.AllPrivileges {
		t.Error("expected AllPrivileges=true")
	}
	if stmt.On.ObjectType != "FUNCTION" {
		t.Errorf("object type = %q, want FUNCTION", stmt.On.ObjectType)
	}
	if !eqStrings(sigNames(stmt.On.Signature), []string{"NUMBER"}) {
		t.Errorf("signature = %v, want [NUMBER]", sigNames(stmt.On.Signature))
	}
}

func TestRevokePriv_FromDatabaseRole(t *testing.T) {
	stmt, errs := testParseRevokeStmt(t, "REVOKE SELECT ON ALL TABLES IN SCHEMA mydb.myschema FROM DATABASE ROLE mydb.dr1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Grantee.Kind != ast.GranteeDatabaseRole {
		t.Errorf("grantee kind = %v, want GranteeDatabaseRole", stmt.Grantee.Kind)
	}
}

// ---------------------------------------------------------------------------
// Position tracking
// ---------------------------------------------------------------------------

func TestGrant_Position(t *testing.T) {
	input := "GRANT SELECT ON TABLE t TO ROLE r"
	stmt, errs := testParseGrantStmt(t, input)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// Walker integration — the GrantStmt's children must be reachable.
// ---------------------------------------------------------------------------

func TestGrant_WalkerVisitsChildren(t *testing.T) {
	stmt, errs := testParseGrantStmt(t, "GRANT SELECT ON TABLE mydb.s.t TO ROLE r")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	var objectNames int
	ast.Inspect(stmt, func(n ast.Node) bool {
		if _, ok := n.(*ast.ObjectName); ok {
			objectNames++
		}
		return true
	})
	// Expect at least the ON target name and the grantee name.
	if objectNames < 2 {
		t.Errorf("walker visited %d ObjectName nodes, want >= 2", objectNames)
	}
}

// ---------------------------------------------------------------------------
// Negative tests — malformed GRANT/REVOKE must be rejected.
// ---------------------------------------------------------------------------

func TestGrant_Negative(t *testing.T) {
	cases := []string{
		"GRANT",                                   // nothing after GRANT
		"GRANT SELECT ON TABLE t",                 // missing TO grantee
		"GRANT SELECT ON TABLE t TO",              // TO with no grantee
		"GRANT ROLE r TO",                         // role grant missing grantee
		"GRANT ROLE r TO ROLE",                    // grantee keyword with no name
		"GRANT SELECT TABLE t TO ROLE r",          // missing ON
		"GRANT OWNERSHIP TO ROLE r",               // ownership missing ON
		"GRANT SELECT ON ALL TABLES TO ROLE r",    // ALL ... missing IN container
		"GRANT SELECT ON FUTURE TABLES TO ROLE r", // FUTURE ... missing IN container
		"GRANT SELECT ON ALL TABLES IN TO ROLE r", // IN with no DATABASE/SCHEMA
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			result := ParseBestEffort(c)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", c, len(result.File.Stmts))
			}
		})
	}
}

func TestRevoke_Negative(t *testing.T) {
	cases := []string{
		"REVOKE",                        // nothing after REVOKE
		"REVOKE SELECT ON TABLE t",      // missing FROM grantee
		"REVOKE SELECT ON TABLE t FROM", // FROM with no grantee
		"REVOKE ROLE r FROM",            // role revoke missing grantee
		"REVOKE GRANT OPTION FOR",       // GRANT OPTION FOR with nothing else
		"REVOKE SELECT FROM ROLE r",     // missing ON
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			result := ParseBestEffort(c)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", c, len(result.File.Stmts))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Official documentation corpus — every GRANT/REVOKE example must parse with
// zero errors. These corpora are the authoritative oracle (truth1).
// ---------------------------------------------------------------------------

// grantCorpusDirs are the official-docs corpora that consist solely of
// GRANT/REVOKE statements (plus benign USE/CREATE/INSERT/SELECT context lines,
// which are filtered out below since those nodes belong to other DAG nodes).
var grantCorpusDirs = []string{
	"testdata/official/grant-role",
	"testdata/official/grant-ownership",
	"testdata/official/grant-privilege",
	"testdata/official/revoke-privilege",
}

func TestGrantRevoke_OfficialCorpus(t *testing.T) {
	for _, dir := range grantCorpusDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertGrantRevokeStatementsParse(t, string(data))
			})
		}
	}
}

func TestGrantRevoke_LegacyCorpus(t *testing.T) {
	path := "testdata/legacy/grant.sql"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	assertGrantRevokeStatementsParse(t, string(data))
}

// assertGrantRevokeStatementsParse parses sql and asserts that every GRANT and
// REVOKE statement in it parses with no errors attributable to that statement.
// Non-GRANT/REVOKE context statements (USE, CREATE, INSERT, SELECT) are checked
// per-segment so that an unimplemented sibling statement type does not mask a
// real GRANT/REVOKE failure.
func assertGrantRevokeStatementsParse(t *testing.T, sql string) {
	t.Helper()
	segs := Split(sql)
	for _, seg := range segs {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)
		isGrant := strings.HasPrefix(upper, "GRANT")
		isRevoke := strings.HasPrefix(upper, "REVOKE")
		if !isGrant && !isRevoke {
			continue // context statement owned by another DAG node
		}
		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if len(errs) > 0 {
			t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
		}
		if isGrant {
			if _, ok := node.(*ast.GrantStmt); !ok {
				t.Errorf("statement %q did not parse to *ast.GrantStmt (got %T)", text, node)
			}
		}
		if isRevoke {
			if _, ok := node.(*ast.RevokeStmt); !ok {
				t.Errorf("statement %q did not parse to *ast.RevokeStmt (got %T)", text, node)
			}
		}
	}
}

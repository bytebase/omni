package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftGrantDatashareParse(t *testing.T) {
	for _, sql := range []string{
		"GRANT ALTER ON DATASHARE salesshare TO myuser WITH GRANT OPTION",
		"GRANT USAGE ON DATASHARE salesshare TO namespace 'a3f3ae8c-14e8-45ba-9eaa-b42b9b7ae635'",
	} {
		stmt := singleStmt(t, sql)
		grant, ok := stmt.(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt, got %T", sql, stmt)
		}
		if !grant.IsGrant {
			t.Fatalf("Parse(%q): expected grant", sql)
		}
		if grant.Objtype != nodes.OBJECT_DATASHARE {
			t.Fatalf("Parse(%q): objtype = %v, want OBJECT_DATASHARE", sql, grant.Objtype)
		}
		if grant.Objects == nil || len(grant.Objects.Items) != 1 {
			t.Fatalf("Parse(%q): objects = %#v, want one object", sql, grant.Objects)
		}
		if grant.Grantees == nil || len(grant.Grantees.Items) != 1 {
			t.Fatalf("Parse(%q): grantees = %#v, want one grantee", sql, grant.Grantees)
		}
	}
}

func TestRedshiftGrantRoleKeywordParse(t *testing.T) {
	for _, sql := range []string{
		"GRANT ROLE sample_role1 TO reguser",
		"GRANT ROLE admin_role TO user1 WITH ADMIN OPTION",
		"GRANT ROLE analyst_role TO ROLE senior_analyst",
	} {
		stmt := singleStmt(t, sql)
		grant, ok := stmt.(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt, got %T", sql, stmt)
		}
		if !grant.IsGrant {
			t.Fatalf("Parse(%q): expected grant role", sql)
		}
		if grant.GrantedRoles == nil || len(grant.GrantedRoles.Items) != 1 {
			t.Fatalf("Parse(%q): granted roles = %#v, want one", sql, grant.GrantedRoles)
		}
		if grant.GranteeRoles == nil || len(grant.GranteeRoles.Items) != 1 {
			t.Fatalf("Parse(%q): grantee roles = %#v, want one", sql, grant.GranteeRoles)
		}
	}
}

func TestRedshiftGrantModelParse(t *testing.T) {
	for _, sql := range []string{
		"GRANT EXECUTE ON MODEL customer_churn_model TO ml_user",
		"GRANT ALL PRIVILEGES ON MODEL recommendation_model TO data_scientist",
	} {
		grant, ok := singleStmt(t, sql).(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt", sql)
		}
		if !grant.IsGrant {
			t.Fatalf("Parse(%q): expected grant", sql)
		}
		if grant.Objtype != nodes.OBJECT_MODEL {
			t.Fatalf("Parse(%q): objtype = %v, want OBJECT_MODEL", sql, grant.Objtype)
		}
		if grant.Objects == nil || len(grant.Objects.Items) != 1 {
			t.Fatalf("Parse(%q): objects = %#v, want one object", sql, grant.Objects)
		}
	}

	// GRANT CREATE MODEL names no object: r_GRANT writes it as a privilege of the
	// warehouse, with the grantee's optional WITH GRANT OPTION.
	for _, tc := range []struct {
		sql      string
		withOpt  bool
		grantees int
	}{
		{"GRANT CREATE MODEL TO ml_admin", false, 1},
		{"GRANT CREATE MODEL TO a WITH GRANT OPTION", true, 1},
		{"GRANT CREATE MODEL TO ROLE a, GROUP g, PUBLIC", false, 3},
	} {
		grant, ok := singleStmt(t, tc.sql).(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt", tc.sql)
		}
		if !grant.IsGrant || grant.Objtype != nodes.OBJECT_MODEL {
			t.Fatalf("Parse(%q): isGrant/objtype = %v/%v, want grant on OBJECT_MODEL", tc.sql, grant.IsGrant, grant.Objtype)
		}
		if grant.Objects != nil {
			t.Fatalf("Parse(%q): objects = %#v, want none", tc.sql, grant.Objects)
		}
		if grant.GrantOption != tc.withOpt {
			t.Fatalf("Parse(%q): grant option = %v, want %v", tc.sql, grant.GrantOption, tc.withOpt)
		}
		if grant.Grantees == nil || len(grant.Grantees.Items) != tc.grantees {
			t.Fatalf("Parse(%q): grantees = %#v, want %d", tc.sql, grant.Grantees, tc.grantees)
		}
	}
}

func TestRedshiftGrantSystemPermissionsParse(t *testing.T) {
	for _, sql := range []string{
		"GRANT CREATE USER TO ROLE user_admin",
		"GRANT ALTER DEFAULT PRIVILEGES TO ROLE privilege_admin",
		"GRANT CREATE OR REPLACE EXTERNAL FUNCTION TO ROLE external_developer",
		"GRANT ACCESS CATALOG TO ROLE catalog_user",
		"GRANT EXPLAIN RLS TO ROLE security_analyst",
	} {
		stmt := singleStmt(t, sql)
		grant, ok := stmt.(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt, got %T", sql, stmt)
		}
		if !grant.IsGrant {
			t.Fatalf("Parse(%q): expected grant system permission", sql)
		}
		if grant.GrantedRoles == nil || len(grant.GrantedRoles.Items) != 1 {
			t.Fatalf("Parse(%q): granted permissions = %#v, want one", sql, grant.GrantedRoles)
		}
		if grant.GranteeRoles == nil || len(grant.GranteeRoles.Items) != 1 {
			t.Fatalf("Parse(%q): grantee roles = %#v, want one", sql, grant.GranteeRoles)
		}
	}
}

package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftRevokeColumnPrivilegesParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE SELECT (col1) ON TABLE t1 FROM PUBLIC",
		"REVOKE UPDATE (col1, col2) ON TABLE t1 FROM user1",
		"REVOKE ALL (order_id, order_date) ON orders FROM report_user",
	} {
		stmt := singleStmt(t, sql)
		revoke, ok := stmt.(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt, got %T", sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke", sql)
		}
		if revoke.Privileges == nil || len(revoke.Privileges.Items) != 1 {
			t.Fatalf("Parse(%q): privileges = %#v, want one privilege", sql, revoke.Privileges)
		}
		priv, ok := revoke.Privileges.Items[0].(*nodes.AccessPriv)
		if !ok {
			t.Fatalf("Parse(%q): privilege = %T, want AccessPriv", sql, revoke.Privileges.Items[0])
		}
		if priv.Cols == nil || len(priv.Cols.Items) == 0 {
			t.Fatalf("Parse(%q): privilege cols = %#v, want column list", sql, priv.Cols)
		}
	}
}

func TestRedshiftRevokeDatashareParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE ALTER ON DATASHARE salesshare FROM myuser",
		"REVOKE USAGE ON DATASHARE salesshare FROM namespace 'a3f3ae8c-14e8-45ba-9eaa-b42b9b7ae635'",
		"REVOKE USAGE ON DATASHARE customerdata FROM ACCOUNT '123456789012'",
		"REVOKE USAGE ON DATASHARE productdata FROM ACCOUNT '987654321098' VIA DATA_CATALOG",
		"REVOKE GRANT OPTION FOR ALTER ON DATASHARE salesshare FROM share_admin",
	} {
		stmt := singleStmt(t, sql)
		revoke, ok := stmt.(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt, got %T", sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke", sql)
		}
		if revoke.Objtype != nodes.OBJECT_DATASHARE {
			t.Fatalf("Parse(%q): objtype = %v, want OBJECT_DATASHARE", sql, revoke.Objtype)
		}
	}
}

func TestRedshiftRevokeRoleKeywordParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE ROLE sample_role1 FROM reguser",
		"REVOKE ROLE analyst_role FROM ROLE senior_analyst",
		"REVOKE ADMIN OPTION FOR ROLE admin_role FROM user1",
		"REVOKE ROLE junior_analyst, ROLE data_entry FROM new_hire",
	} {
		stmt := singleStmt(t, sql)
		revoke, ok := stmt.(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt, got %T", sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke role", sql)
		}
		if revoke.GrantedRoles == nil || len(revoke.GrantedRoles.Items) == 0 {
			t.Fatalf("Parse(%q): granted roles = %#v, want roles", sql, revoke.GrantedRoles)
		}
		if revoke.GranteeRoles == nil || len(revoke.GranteeRoles.Items) == 0 {
			t.Fatalf("Parse(%q): grantee roles = %#v, want grantees", sql, revoke.GranteeRoles)
		}
	}
}

func TestRedshiftRevokeModelParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE EXECUTE ON MODEL customer_churn_model FROM ml_user",
		"REVOKE ALL PRIVILEGES ON MODEL recommendation_model FROM data_scientist",
		"REVOKE CREATE MODEL FROM untrusted_user",
	} {
		stmt := singleStmt(t, sql)
		switch revoke := stmt.(type) {
		case *nodes.GrantStmt:
			if revoke.IsGrant {
				t.Fatalf("Parse(%q): expected revoke", sql)
			}
			if revoke.Objtype != nodes.OBJECT_MODEL {
				t.Fatalf("Parse(%q): objtype = %v, want OBJECT_MODEL", sql, revoke.Objtype)
			}
		case *nodes.GrantRoleStmt:
			if revoke.IsGrant {
				t.Fatalf("Parse(%q): expected revoke role/system privilege", sql)
			}
		default:
			t.Fatalf("Parse(%q): expected GrantStmt or GrantRoleStmt, got %T", sql, stmt)
		}
	}
}

func TestRedshiftRevokeSystemPermissionsParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE CREATE USER FROM ROLE user_admin",
		"REVOKE ALTER DEFAULT PRIVILEGES FROM ROLE privilege_admin",
		"REVOKE CREATE OR REPLACE EXTERNAL FUNCTION FROM ROLE external_developer",
		"REVOKE ACCESS CATALOG FROM ROLE catalog_user",
		"REVOKE TRUNCATE TABLE FROM ROLE data_manager",
		"REVOKE EXPLAIN MASKING FROM ROLE compliance_officer",
	} {
		stmt := singleStmt(t, sql)
		revoke, ok := stmt.(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt, got %T", sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke system permission", sql)
		}
		if revoke.GrantedRoles == nil || len(revoke.GrantedRoles.Items) != 1 {
			t.Fatalf("Parse(%q): granted permissions = %#v, want one", sql, revoke.GrantedRoles)
		}
		if revoke.GranteeRoles == nil || len(revoke.GranteeRoles.Items) != 1 {
			t.Fatalf("Parse(%q): grantee roles = %#v, want one", sql, revoke.GranteeRoles)
		}
	}
}

func TestRedshiftRevokeFromRLSPolicyParse(t *testing.T) {
	stmt := singleStmt(t, "REVOKE SELECT ON TABLE sensitive_data FROM RLS POLICY high_security_policy")
	revoke, ok := stmt.(*nodes.GrantStmt)
	if !ok {
		t.Fatalf("expected GrantStmt, got %T", stmt)
	}
	if revoke.IsGrant {
		t.Fatal("expected revoke")
	}
	if revoke.Grantees == nil || len(revoke.Grantees.Items) != 1 {
		t.Fatalf("grantees = %#v, want one grantee", revoke.Grantees)
	}
	grantee, ok := revoke.Grantees.Items[0].(*nodes.RoleSpec)
	if !ok {
		t.Fatalf("grantee = %T, want RoleSpec", revoke.Grantees.Items[0])
	}
	if grantee.Rolename != "rls_policy:high_security_policy" {
		t.Fatalf("grantee = %q, want rls_policy:high_security_policy", grantee.Rolename)
	}
}

func TestRedshiftRevokeAssumeRoleParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM sales_user FOR COPY",
		"REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM ml_user FOR CREATE MODEL",
		"REVOKE ASSUMEROLE DEFAULT FROM legacy_user FOR COPY",
		"REVOKE ASSUMEROLE ALL FROM untrusted_user FOR ALL",
		"REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/Role1', 'arn:aws:iam::123456789012:role/Role2' FROM multi_role_user FOR COPY",
	} {
		stmt := singleStmt(t, sql)
		revoke, ok := stmt.(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt, got %T", sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke assume role", sql)
		}
		if revoke.GrantedRoles == nil || len(revoke.GrantedRoles.Items) == 0 {
			t.Fatalf("Parse(%q): assume roles = %#v, want roles", sql, revoke.GrantedRoles)
		}
		if revoke.GranteeRoles == nil || len(revoke.GranteeRoles.Items) != 1 {
			t.Fatalf("Parse(%q): grantees = %#v, want one grantee", sql, revoke.GranteeRoles)
		}
		if findDefElem(revoke.Opt, "for") == nil {
			t.Fatalf("Parse(%q): opts = %#v, want FOR action", sql, revoke.Opt)
		}
	}
}

func TestRedshiftRevokeExternalObjectIAMRoleParse(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		objtype nodes.ObjectType
	}{
		{"REVOKE CREATE ON EXTERNAL SCHEMA spectrum_schema FROM IAM_ROLE 'spectrum_user'", nodes.OBJECT_SCHEMA},
		{"REVOKE SELECT ON EXTERNAL TABLE spectrum_schema.data FROM IAM_ROLE 'arn:aws:iam::123456789012:role/SpectrumRole'", nodes.OBJECT_FOREIGN_TABLE},
	} {
		stmt := singleStmt(t, tc.sql)
		revoke, ok := stmt.(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt, got %T", tc.sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke", tc.sql)
		}
		if revoke.Objtype != tc.objtype {
			t.Fatalf("Parse(%q): objtype = %v, want %v", tc.sql, revoke.Objtype, tc.objtype)
		}
		if revoke.Grantees == nil || len(revoke.Grantees.Items) != 1 {
			t.Fatalf("Parse(%q): grantees = %#v, want one grantee", tc.sql, revoke.Grantees)
		}
		grantee, ok := revoke.Grantees.Items[0].(*nodes.RoleSpec)
		if !ok {
			t.Fatalf("Parse(%q): grantee = %T, want RoleSpec", tc.sql, revoke.Grantees.Items[0])
		}
		if !strings.HasPrefix(grantee.Rolename, "iam_role:") {
			t.Fatalf("Parse(%q): grantee = %q, want iam_role prefix", tc.sql, grantee.Rolename)
		}
	}
}

func TestRedshiftRevokeCopyJobParse(t *testing.T) {
	for _, sql := range []string{
		"REVOKE ALTER ON COPY JOB my_copy_job FROM job_modifier",
		"REVOKE DROP ON COPY JOB old_copy_job FROM job_cleaner",
		"REVOKE ALTER, DROP ON COPY JOB job1, job2, job3 FROM job_admin",
	} {
		stmt := singleStmt(t, sql)
		revoke, ok := stmt.(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt, got %T", sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke", sql)
		}
		if revoke.Objtype != nodes.OBJECT_COPY_JOB {
			t.Fatalf("Parse(%q): objtype = %v, want OBJECT_COPY_JOB", sql, revoke.Objtype)
		}
		if revoke.Objects == nil || len(revoke.Objects.Items) == 0 {
			t.Fatalf("Parse(%q): objects = %#v, want copy jobs", sql, revoke.Objects)
		}
	}
}

func TestRedshiftRevokeScopedPermissionsParse(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		objtype nodes.ObjectType
	}{
		{"REVOKE CREATE FOR SCHEMAS IN DATABASE mydb FROM schema_creator", nodes.OBJECT_SCHEMA},
		{"REVOKE SELECT FOR TABLES IN SCHEMA reporting FROM analyst_role", nodes.OBJECT_TABLE},
		{"REVOKE SELECT FOR TABLES IN SCHEMA finance DATABASE prod_db FROM finance_reader", nodes.OBJECT_TABLE},
		{"REVOKE EXECUTE FOR FUNCTIONS IN SCHEMA utilities FROM function_user", nodes.OBJECT_FUNCTION},
		{"REVOKE EXECUTE FOR PROCEDURES IN SCHEMA etl FROM etl_runner", nodes.OBJECT_PROCEDURE},
		{"REVOKE USAGE FOR LANGUAGES IN DATABASE dev_db FROM developer", nodes.OBJECT_LANGUAGE},
		{"REVOKE ALTER FOR COPY JOBS IN SCHEMA data_import FROM job_modifier", nodes.OBJECT_COPY_JOB},
	} {
		stmt := singleStmt(t, tc.sql)
		revoke, ok := stmt.(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt, got %T", tc.sql, stmt)
		}
		if revoke.IsGrant {
			t.Fatalf("Parse(%q): expected revoke", tc.sql)
		}
		if revoke.Targtype != nodes.ACL_TARGET_SCOPED {
			t.Fatalf("Parse(%q): targtype = %v, want ACL_TARGET_SCOPED", tc.sql, revoke.Targtype)
		}
		if revoke.Objtype != tc.objtype {
			t.Fatalf("Parse(%q): objtype = %v, want %v", tc.sql, revoke.Objtype, tc.objtype)
		}
		if revoke.Objects == nil || len(revoke.Objects.Items) == 0 {
			t.Fatalf("Parse(%q): objects = %#v, want scoped options", tc.sql, revoke.Objects)
		}
	}
}

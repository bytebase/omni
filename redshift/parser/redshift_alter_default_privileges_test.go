package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAlterDefaultPrivilegesRoleAndGroupGranteesParse(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		wantGrantee string
		wantObjtype nodes.ObjectType
	}{
		{
			sql:         "ALTER DEFAULT PRIVILEGES GRANT SELECT, INSERT ON TABLES TO ROLE data_analyst",
			wantGrantee: "data_analyst",
			wantObjtype: nodes.OBJECT_TABLE,
		},
		{
			sql:         "ALTER DEFAULT PRIVILEGES FOR USER proc_owner GRANT EXECUTE ON PROCEDURES TO GROUP app_users",
			wantGrantee: "app_users",
			wantObjtype: nodes.OBJECT_PROCEDURE,
		},
		{
			sql:         "ALTER DEFAULT PRIVILEGES IN SCHEMA app_schema GRANT ALL PRIVILEGES ON PROCEDURES TO ROLE app_admin",
			wantGrantee: "app_admin",
			wantObjtype: nodes.OBJECT_PROCEDURE,
		},
	} {
		stmt := singleStmt(t, tc.sql)
		alter, ok := stmt.(*nodes.AlterDefaultPrivilegesStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected AlterDefaultPrivilegesStmt, got %T", tc.sql, stmt)
		}
		if alter.Action == nil {
			t.Fatalf("Parse(%q): expected action", tc.sql)
		}
		if !alter.Action.IsGrant {
			t.Fatalf("Parse(%q): expected grant action", tc.sql)
		}
		if alter.Action.Targtype != nodes.ACL_TARGET_DEFAULTS {
			t.Fatalf("Parse(%q): targtype = %v, want ACL_TARGET_DEFAULTS", tc.sql, alter.Action.Targtype)
		}
		if alter.Action.Objtype != tc.wantObjtype {
			t.Fatalf("Parse(%q): objtype = %v, want %v", tc.sql, alter.Action.Objtype, tc.wantObjtype)
		}
		if alter.Action.Grantees == nil || len(alter.Action.Grantees.Items) != 1 {
			t.Fatalf("Parse(%q): grantees = %#v, want one grantee", tc.sql, alter.Action.Grantees)
		}
		grantee, ok := alter.Action.Grantees.Items[0].(*nodes.RoleSpec)
		if !ok {
			t.Fatalf("Parse(%q): grantee = %T, want RoleSpec", tc.sql, alter.Action.Grantees.Items[0])
		}
		if grantee.Rolename != tc.wantGrantee {
			t.Fatalf("Parse(%q): grantee = %q, want %q", tc.sql, grantee.Rolename, tc.wantGrantee)
		}
	}
}

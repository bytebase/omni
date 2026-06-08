package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCreateRoleExternalIDParse(t *testing.T) {
	stmt := singleStmt(t, `CREATE ROLE sample_role1 EXTERNALID "ABC123"`)
	createRole, ok := stmt.(*nodes.CreateRoleStmt)
	if !ok {
		t.Fatalf("expected CreateRoleStmt, got %T", stmt)
	}
	if createRole.Role != "sample_role1" {
		t.Fatalf("role = %q, want sample_role1", createRole.Role)
	}
	assertDefElemString(t, createRole.Options, "externalid", "ABC123")
}

func TestRedshiftAlterRoleOwnerParse(t *testing.T) {
	for _, sql := range []string{
		`ALTER ROLE sample_role1 OWNER TO user1`,
		`ALTER ROLE sample_role1 WITH OWNER TO user1`,
	} {
		stmt := singleStmt(t, sql)
		alterOwner, ok := stmt.(*nodes.AlterOwnerStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected AlterOwnerStmt, got %T", sql, stmt)
		}
		if alterOwner.ObjectType != nodes.OBJECT_ROLE {
			t.Fatalf("Parse(%q): object type = %v, want OBJECT_ROLE", sql, alterOwner.ObjectType)
		}
		object, ok := alterOwner.Object.(*nodes.String)
		if !ok {
			t.Fatalf("Parse(%q): owner object = %T, want String", sql, alterOwner.Object)
		}
		if object.Str != "sample_role1" {
			t.Fatalf("Parse(%q): role object = %q, want sample_role1", sql, object.Str)
		}
		if alterOwner.Newowner == nil || alterOwner.Newowner.Rolename != "user1" {
			t.Fatalf("Parse(%q): new owner = %#v, want user1", sql, alterOwner.Newowner)
		}
	}
}

func TestRedshiftAlterRoleExternalIDParse(t *testing.T) {
	for _, sql := range []string{
		`ALTER ROLE sample_role1 EXTERNALID TO "XYZ456"`,
		`ALTER ROLE sample_role1 WITH EXTERNALID TO "XYZ456"`,
	} {
		stmt := singleStmt(t, sql)
		alterRole, ok := stmt.(*nodes.AlterRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected AlterRoleStmt, got %T", sql, stmt)
		}
		if alterRole.Role == nil || alterRole.Role.Rolename != "sample_role1" {
			t.Fatalf("Parse(%q): role = %#v, want sample_role1", sql, alterRole.Role)
		}
		assertDefElemString(t, alterRole.Options, "externalid", "XYZ456")
	}
}

func TestRedshiftAlterRoleWithRenameParse(t *testing.T) {
	stmt := singleStmt(t, `ALTER ROLE role1 WITH RENAME TO role2`)
	rename, ok := stmt.(*nodes.RenameStmt)
	if !ok {
		t.Fatalf("expected RenameStmt, got %T", stmt)
	}
	if rename.RenameType != nodes.OBJECT_ROLE {
		t.Fatalf("rename type = %v, want OBJECT_ROLE", rename.RenameType)
	}
	if rename.Subname != "role1" || rename.Newname != "role2" {
		t.Fatalf("rename %q -> %q, want role1 -> role2", rename.Subname, rename.Newname)
	}
}

func TestRedshiftDropRoleBehaviorParse(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		option string
	}{
		{`DROP ROLE sample_role FORCE`, "force"},
		{`DROP ROLE sample_role RESTRICT`, "restrict"},
	} {
		stmt := singleStmt(t, tc.sql)
		dropRole, ok := stmt.(*nodes.DropRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected DropRoleStmt, got %T", tc.sql, stmt)
		}
		if dropRole.Roles == nil || len(dropRole.Roles.Items) != 1 {
			t.Fatalf("Parse(%q): roles = %#v, want one role", tc.sql, dropRole.Roles)
		}
		if findDefElem(dropRole.Options, tc.option) == nil {
			t.Fatalf("Parse(%q): expected %q option", tc.sql, tc.option)
		}
	}
}

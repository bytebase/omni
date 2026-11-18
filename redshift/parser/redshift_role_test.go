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

// TestRedshiftAlterRoleClauseListParse pins r_ALTER_ROLE's comma-separated
// RENAME TO and OWNER TO clauses and the EXTERNALID TO tail. One clause keeps
// the node PostgreSQL gives it; a repeat or an EXTERNALID tail is an
// AlterRoleStmt's option list.
func TestRedshiftAlterRoleClauseListParse(t *testing.T) {
	rename, ok := singleStmt(t, "ALTER ROLE role1 RENAME TO role2").(*nodes.RenameStmt)
	if !ok {
		t.Fatalf("expected RenameStmt, got %T", singleStmt(t, "ALTER ROLE role1 RENAME TO role2"))
	}
	if rename.Subname != "role1" || rename.Newname != "role2" {
		t.Fatalf("rename %q -> %q, want role1 -> role2", rename.Subname, rename.Newname)
	}
	if rename.Loc.Start < 0 || rename.Loc.End <= rename.Loc.Start {
		t.Fatalf("Loc = %#v, want the statement's byte range", rename.Loc)
	}

	owner, ok := singleStmt(t, "ALTER ROLE role1 OWNER TO user1").(*nodes.AlterOwnerStmt)
	if !ok {
		t.Fatalf("expected AlterOwnerStmt, got %T", singleStmt(t, "ALTER ROLE role1 OWNER TO user1"))
	}
	if owner.ObjectType != nodes.OBJECT_ROLE || owner.Newowner.Rolename != "user1" {
		t.Fatalf("owner = %#v, want role1 owned by user1", owner)
	}

	for _, tc := range []struct {
		sql    string
		names  []string
		values []string
	}{
		{"ALTER ROLE a RENAME TO b, RENAME TO c", []string{"rename_to", "rename_to"}, []string{"b", "c"}},
		{"ALTER ROLE a RENAME TO b, OWNER TO c", []string{"rename_to", "owner_to"}, []string{"b", "c"}},
		{"ALTER ROLE a RENAME TO b EXTERNALID TO x", []string{"rename_to", "externalid"}, []string{"b", "x"}},
	} {
		alter, ok := singleStmt(t, tc.sql).(*nodes.AlterRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected AlterRoleStmt, got %T", tc.sql, singleStmt(t, tc.sql))
		}
		if alter.Role == nil || alter.Role.Rolename != "a" {
			t.Fatalf("Parse(%q): role = %#v, want a", tc.sql, alter.Role)
		}
		if alter.Options == nil || len(alter.Options.Items) != len(tc.names) {
			t.Fatalf("Parse(%q): options = %#v, want %d", tc.sql, alter.Options, len(tc.names))
		}
		for i, name := range tc.names {
			elem, ok := alter.Options.Items[i].(*nodes.DefElem)
			if !ok || elem.Defname != name {
				t.Fatalf("Parse(%q): option %d = %#v, want %q", tc.sql, i, alter.Options.Items[i], name)
			}
			switch arg := elem.Arg.(type) {
			case *nodes.String:
				if arg.Str != tc.values[i] {
					t.Fatalf("Parse(%q): option %d = %q, want %q", tc.sql, i, arg.Str, tc.values[i])
				}
			case *nodes.RoleSpec:
				if arg.Rolename != tc.values[i] {
					t.Fatalf("Parse(%q): option %d = %q, want %q", tc.sql, i, arg.Rolename, tc.values[i])
				}
			default:
				t.Fatalf("Parse(%q): option %d arg = %T", tc.sql, i, elem.Arg)
			}
		}
	}

	// EXTERNALID takes TO, and a comma takes a clause after it.
	for _, sql := range []string{
		"ALTER ROLE a RENAME TO b EXTERNALID b",
		"ALTER ROLE a RENAME TO b,",
		"ALTER ROLE a RENAME TO b, TO c",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftCreateGroupWithListParse pins r_CREATE_GROUP's WITH clause, which
// the docs repeat with commas between the user lists and whose USER list is
// optional.
func TestRedshiftCreateGroupWithListParse(t *testing.T) {
	for _, sql := range []string{
		"CREATE GROUP a WITH, WITH",
		"CREATE GROUP a WITH USER b, WITH USER c",
		"CREATE GROUP a WITH USER b, WITH",
	} {
		create, ok := singleStmt(t, sql).(*nodes.CreateRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected CreateRoleStmt, got %T", sql, singleStmt(t, sql))
		}
		if create.StmtType != nodes.ROLESTMT_GROUP || create.Role != "a" {
			t.Fatalf("Parse(%q): = %#v, want group a", sql, create)
		}
	}
	if _, err := Parse("CREATE GROUP a WITH,"); err == nil {
		t.Error("CREATE GROUP a WITH, parsed, want a syntax error")
	}
}

package parser

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

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

func TestGrantPrivilegesOnObjects(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT, ALTER, CREATE MATERIALIZED VIEW ON db1.*, db2.t TO USER 'jack'@'%' WITH GRANT OPTION")
	if !reflect.DeepEqual(stmt.Privileges, []string{"SELECT", "ALTER", "CREATE MATERIALIZED VIEW"}) {
		t.Errorf("Privileges = %v", stmt.Privileges)
	}
	if stmt.ObjectType != "" || len(stmt.Objects) != 2 || stmt.Object.String() != "db1.*" {
		t.Errorf("objects = %q %v", stmt.ObjectType, stmt.Objects)
	}
	if stmt.ToType != "USER" || !reflect.DeepEqual(stmt.Grantees, []string{"jack@%"}) || !stmt.WithGrantOption {
		t.Errorf("grantee = %q %v %v", stmt.ToType, stmt.Grantees, stmt.WithGrantOption)
	}
}

func TestGrantOnTypedObjects(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT USAGE ON RESOURCE GROUP * TO ROLE 'r'")
	if stmt.ObjectType != "RESOURCE GROUP" || stmt.ToType != "ROLE" || stmt.Grantees[0] != "r" {
		t.Errorf("stmt = %+v", stmt)
	}
	stmt = parseGrantStmt(t, "GRANT ALL PRIVILEGES ON ALL TABLES IN DATABASE d TO a")
	if stmt.ObjectType != "ALL TABLES" || stmt.InDatabase != "d" || stmt.Privileges[0] != "ALL PRIVILEGES" {
		t.Errorf("stmt = %+v", stmt)
	}
	stmt = parseGrantStmt(t, "GRANT USAGE ON GLOBAL FUNCTION f(INT), g(, ...) TO a")
	if stmt.ObjectType != "GLOBAL FUNCTION" || len(stmt.Objects) != 2 || len(stmt.FunctionArgs[0]) != 1 {
		t.Errorf("stmt = %+v", stmt)
	}
	stmt = parseGrantStmt(t, "GRANT IMPERSONATE ON USER a, 'b'@'%' TO c")
	if stmt.ObjectType != "USER" || len(stmt.Objects) != 2 {
		t.Errorf("stmt = %+v", stmt)
	}
	stmt = parseGrantStmt(t, "GRANT OPERATE ON SYSTEM TO a")
	if stmt.ObjectType != "SYSTEM" || len(stmt.Objects) != 0 {
		t.Errorf("stmt = %+v", stmt)
	}
}

func TestGrantRoles(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT r1, 'r2' TO USER u@'%'")
	if !reflect.DeepEqual(stmt.Roles, []string{"r1", "r2"}) || stmt.ToType != "USER" || stmt.Grantees[0] != "u@%" {
		t.Errorf("stmt = %+v", stmt)
	}
	rv := parseRevokeStmt(t, "REVOKE r1 FROM ROLE r2")
	if !reflect.DeepEqual(rv.Roles, []string{"r1"}) || rv.FromType != "ROLE" || rv.Revokees[0] != "r2" {
		t.Errorf("stmt = %+v", rv)
	}
}

func TestRevokePrivileges(t *testing.T) {
	stmt := parseRevokeStmt(t, "REVOKE INSERT ON ALL CATALOGS IN ALL DATABASES FROM a")
	if stmt.ObjectType != "ALL CATALOGS" || !stmt.InAllDatabases || stmt.Privileges[0] != "INSERT" {
		t.Errorf("stmt = %+v", stmt)
	}
}

// TestGrantRevokeGrammar holds GRANT and REVOKE to StarRocks.g4; the
// rejected forms are Doris's, which this file once pinned, and the near
// misses testdata/compat.tsv judged.
func TestGrantRevokeGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"GRANT ALL ON * TO ROLE", true}, // ROLE names the user
		{"GRANT a TO ROLE", true},
		{"GRANT ALL ON *, * TO a", true},
		{"GRANT ALL ON catalog TO a", true}, // CATALOG names the object
		{"GRANT ALL ON CATALOG * TO a", true},
		{"GRANT ALL ON STORAGE VOLUME * TO a", true},
		{"GRANT ALL ON ALL GLOBAL FUNCTIONS TO a", true},
		{"GRANT ALL ON * TO a@['h']", true},
		{"GRANT ALL ON ALL TO a", false},
		{"GRANT ALL ON STORAGE * TO a", false},
		{"GRANT RESOURCE ON * TO a", false},
		{"GRANT CREATE ON * TO a", false},
		{"GRANT ALL ON * TO a WITH GRANT", false},
		{"REVOKE ALL ON * FROM a WITH GRANT OPTION", false},
		{"GRANT SELECT_PRIV ON *.* TO 'jack'@'%'", false},
		{"GRANT SELECT ON *.*.* TO 'jack'@'%'", false},
		{"GRANT USAGE ON WORKLOAD GROUP 'g1' TO a", false},
		{"GRANT ROLE r TO a", false},
		{"GRANT IMPERSONATE, SELECT ON USER a TO b", false},
		{"GRANT SELECT ON * TO a, b", false},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}

func TestGrantRevokeNodeTags(t *testing.T) {
	if n := parseOne(t, "GRANT SELECT ON *.* TO 'user'@'%'"); n.Tag() != ast.T_GrantStmt {
		t.Errorf("grant tag = %v, want T_GrantStmt", n.Tag())
	}
	if n := parseOne(t, "REVOKE SELECT ON db.* FROM 'user'@'%'"); n.Tag() != ast.T_RevokeStmt {
		t.Errorf("revoke tag = %v, want T_RevokeStmt", n.Tag())
	}
}

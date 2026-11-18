package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestRedshiftGrantExternalParticipants pins GRANT and REVOKE on an external
// table, an external schema and a datashare, whose alternatives of r_GRANT and
// r_REVOKE list their own privileges and grantees
// (https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html).
func TestRedshiftGrantExternalParticipants(t *testing.T) {
	for _, tc := range []struct {
		sql      string
		objtype  nodes.ObjectType
		isGrant  bool
		grantees int
	}{
		{"GRANT SELECT(a) ON EXTERNAL TABLE s.t TO IAM_ROLE 'arn:aws:iam::1:role/r'", nodes.OBJECT_FOREIGN_TABLE, true, 1},
		{"GRANT SELECT(a) ON EXTERNAL TABLE s.t TO IAM_ROLE r", nodes.OBJECT_FOREIGN_TABLE, true, 1},
		{"GRANT ALL(a) ON EXTERNAL TABLE s.t TO IAM_ROLE 'r', IAM_ROLE 'r2' WITH GRANT OPTION", nodes.OBJECT_FOREIGN_TABLE, true, 2},
		{"GRANT SELECT, ALTER, DROP, DELETE, INSERT ON EXTERNAL TABLE s.t, s.u TO PUBLIC", nodes.OBJECT_FOREIGN_TABLE, true, 1},
		{"GRANT ALL PRIVILEGES ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'", nodes.OBJECT_FOREIGN_TABLE, true, 1},
		{"REVOKE SELECT(a) ON EXTERNAL TABLE s.t FROM IAM_ROLE 'r'", nodes.OBJECT_FOREIGN_TABLE, false, 1},
		{"REVOKE SELECT, INSERT ON EXTERNAL TABLE s.t FROM PUBLIC", nodes.OBJECT_FOREIGN_TABLE, false, 1},
		{"REVOKE ALL ON EXTERNAL TABLE s.t FROM IAM_ROLE 'r', IAM_ROLE 'r2'", nodes.OBJECT_FOREIGN_TABLE, false, 2},
		{"GRANT CREATE ON EXTERNAL SCHEMA s TO IAM_ROLE 'r'", nodes.OBJECT_EXTERNAL_SCHEMA, true, 1},
		{"GRANT CREATE, ALTER, DROP ON EXTERNAL SCHEMA s, s2 TO IAM_ROLE r WITH GRANT OPTION", nodes.OBJECT_EXTERNAL_SCHEMA, true, 1},
		{"GRANT ALL ON EXTERNAL SCHEMA s TO IAM_ROLE 'r'", nodes.OBJECT_EXTERNAL_SCHEMA, true, 1},
		{"REVOKE DROP ON EXTERNAL SCHEMA s FROM IAM_ROLE 'r', IAM_ROLE 'r2'", nodes.OBJECT_EXTERNAL_SCHEMA, false, 2},
		{"GRANT ALTER ON DATASHARE d TO a WITH GRANT OPTION", nodes.OBJECT_DATASHARE, true, 1},
		{"GRANT SHARE ON DATASHARE d TO ROLE r, GROUP g, PUBLIC", nodes.OBJECT_DATASHARE, true, 3},
		{"GRANT ALTER ON DATASHARE d TO role", nodes.OBJECT_DATASHARE, true, 1},
		{"GRANT USAGE ON DATASHARE d TO NAMESPACE 'n'", nodes.OBJECT_DATASHARE, true, 1},
		{"GRANT USAGE ON DATASHARE d TO ACCOUNT '1' VIA DATA CATALOG", nodes.OBJECT_DATASHARE, true, 1},
		{"REVOKE ALTER ON DATASHARE d FROM a WITH GRANT OPTION", nodes.OBJECT_DATASHARE, false, 1},
		{"REVOKE ALTER ON DATASHARE d FROM role", nodes.OBJECT_DATASHARE, false, 1},
		{"REVOKE USAGE ON DATASHARE d FROM NAMESPACE 'n', 'n2'", nodes.OBJECT_DATASHARE, false, 2},
		{"REVOKE USAGE ON DATASHARE d FROM ACCOUNT '1' VIA DATA CATALOG, '2'", nodes.OBJECT_DATASHARE, false, 2},
		{"GRANT EXECUTE ON PROCEDURE p() TO DATASHARE d", nodes.OBJECT_PROCEDURE, true, 1},
		{"REVOKE EXECUTE ON PROCEDURE p(INT) FROM DATASHARE d", nodes.OBJECT_PROCEDURE, false, 1},
	} {
		grant, ok := singleStmt(t, tc.sql).(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt", tc.sql)
		}
		if grant.IsGrant != tc.isGrant {
			t.Errorf("Parse(%q): IsGrant = %v, want %v", tc.sql, grant.IsGrant, tc.isGrant)
		}
		if grant.Objtype != tc.objtype {
			t.Errorf("Parse(%q): objtype = %v, want %v", tc.sql, grant.Objtype, tc.objtype)
		}
		if grant.Grantees == nil || len(grant.Grantees.Items) != tc.grantees {
			t.Errorf("Parse(%q): grantees = %#v, want %d", tc.sql, grant.Grantees, tc.grantees)
		}
	}

	// A datashare consumer keeps the VIA DATA CATALOG qualifier, per account.
	for sql, want := range map[string][]string{
		"GRANT USAGE ON DATASHARE d TO ACCOUNT '1' VIA DATA CATALOG":         {"account_via_data_catalog:1"},
		"GRANT USAGE ON DATASHARE d TO ACCOUNT '1'":                          {"account:1"},
		"GRANT USAGE ON DATASHARE d TO NAMESPACE 'n'":                        {"namespace:n"},
		"REVOKE USAGE ON DATASHARE d FROM ACCOUNT '1' VIA DATA CATALOG, '2'": {"account_via_data_catalog:1", "account:2"},
	} {
		grant, ok := singleStmt(t, sql).(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt", sql)
		}
		if len(grant.Grantees.Items) != len(want) {
			t.Fatalf("Parse(%q): grantees = %#v, want %v", sql, grant.Grantees, want)
		}
		for i, name := range want {
			rs := grant.Grantees.Items[i].(*nodes.RoleSpec)
			if rs.Rolename != name {
				t.Errorf("Parse(%q): grantee %d = %q, want %q", sql, i, rs.Rolename, name)
			}
		}
	}

	// r_REVOKE writes a producer-side grantee as ID ( WITH GRANT OPTION )?, so
	// the suffix revokes the grant option, as GRANT OPTION FOR does.
	revoke := singleStmt(t, "REVOKE ALTER ON DATASHARE d FROM a WITH GRANT OPTION").(*nodes.GrantStmt)
	if !revoke.GrantOption {
		t.Error("REVOKE ... WITH GRANT OPTION: GrantOption = false, want true")
	}

	// DATASHARE is unreserved, so a grant that names no datashare after it reads
	// it as the object's name, and the statement is a table grant.
	if grant := singleStmt(t, "GRANT ALTER ON DATASHARE TO a").(*nodes.GrantStmt); grant.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("ON DATASHARE TO: objtype = %v, want OBJECT_TABLE", grant.Objtype)
	}

	for _, sql := range []string{
		// Privileges each alternative does not list.
		"GRANT USAGE ON EXTERNAL SCHEMA s TO IAM_ROLE 'r'",
		"GRANT SELECT ON EXTERNAL SCHEMA s TO IAM_ROLE 'r'",
		"GRANT TRUNCATE ON EXTERNAL SCHEMA s TO IAM_ROLE 'r'",
		"GRANT CREATE ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'",
		"GRANT UPDATE ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'",
		"GRANT USAGE ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'",
		"GRANT SELECT ON DATASHARE d TO a",
		"GRANT ALTER, SHARE ON DATASHARE d TO a",
		"GRANT USAGE, ALTER ON DATASHARE d TO a",
		"GRANT ALL ON DATASHARE d TO a",
		"REVOKE TRUNCATE ON EXTERNAL SCHEMA s FROM IAM_ROLE 'r'",
		"REVOKE CREATE ON EXTERNAL TABLE s.t FROM IAM_ROLE 'r'",
		"REVOKE SHARE, ALTER ON DATASHARE d FROM a",
		// Grantees each alternative does not list.
		"GRANT SELECT ON EXTERNAL TABLE s.t TO a",
		"GRANT SELECT ON EXTERNAL TABLE s.t TO ROLE r",
		"GRANT SELECT ON EXTERNAL TABLE s.t TO PUBLICACCESSIBLE",
		"GRANT SELECT(a) ON EXTERNAL TABLE s.t TO PUBLIC",
		"GRANT SELECT(a), INSERT ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'",
		"GRANT SELECT(a), (b) ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'",
		"GRANT ALTER(a) ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'",
		"GRANT CREATE ON EXTERNAL SCHEMA s TO a",
		"GRANT CREATE ON EXTERNAL SCHEMA s TO PUBLIC",
		"REVOKE SELECT ON EXTERNAL TABLE s.t FROM a",
		"REVOKE CREATE ON EXTERNAL SCHEMA s FROM PUBLIC",
		"GRANT USAGE ON DATASHARE d TO a",
		"GRANT USAGE ON DATASHARE d TO ROLE a",
		"GRANT USAGE ON DATASHARE d TO ACCOUNT '1' VIA CATALOG",
		"GRANT USAGE ON DATASHARE d TO NAMESPACE 'n', 'n2'",
		"GRANT ALTER ON DATASHARE d TO DATASHARE e",
		// A privilege is a keyword, not a quoted name.
		`GRANT "SELECT" ON EXTERNAL TABLE s.t TO IAM_ROLE 'r'`,
		`GRANT SELECT ON EXTERNAL TABLE s.t TO "IAM_ROLE" 'r'`,
		// A datashare is one name, and the keyword is never a quoted one.
		"GRANT ALTER ON DATASHARE d, d2 TO a",
		`GRANT ALTER ON "DATASHARE" d TO a`,
		// The docs' one-token GRANT_OPTION is notation, not SQL.
		"REVOKE GRANT_OPTION CREATE FOR COPY JOBS IN DATABASE d FROM a",
		"REVOKE GRANT OPTION CREATE MODEL FROM a",
		// The procedure-to-datashare alternatives name one procedure with its
		// argument list and end at the datashare name.
		"GRANT EXECUTE ON PROCEDURE p() TO DATASHARE d WITH GRANT OPTION",
		"REVOKE EXECUTE ON PROCEDURE p() FROM DATASHARE d RESTRICT",
		"GRANT EXECUTE ON PROCEDURE p TO DATASHARE d",
		"GRANT EXECUTE ON PROCEDURE p(), q() TO DATASHARE d",
		"GRANT EXECUTE ON ALL PROCEDURES TO DATASHARE d",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftGrantFederatedIdentity pins CONNECT and DEBUG, whose grantees are
// AWS IAM Identity Center identities
// (https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html).
func TestRedshiftGrantFederatedIdentity(t *testing.T) {
	for _, tc := range []struct {
		sql      string
		grantees []string
		option   string
	}{
		{"GRANT CONNECT TO AWSIDC:alice", []string{"awsidc:alice"}, ""},
		{`GRANT CONNECT TO "AWSIDC:alice"`, []string{"AWSIDC:alice"}, ""},
		{"GRANT CONNECT TO ROLE AWSIDC:analysts", []string{"awsidc:analysts"}, ""},
		{`GRANT CONNECT TO ROLE "AWSIDC:analysts"`, []string{"AWSIDC:analysts"}, ""},
		{"GRANT CONNECT TO ROLE :analysts", []string{"role:analysts"}, ""},
		{"GRANT CONNECT TO PUBLIC", nil, ""},
		{"GRANT CONNECT ON WORKGROUP TO AWSIDC:alice", []string{"awsidc:alice"}, "on_workgroup"},
		{"GRANT DEBUG ON DATABASE d TO AWSIDC:alice", []string{"awsidc:alice"}, "on_database"},
		{`GRANT DEBUG ON DATABASE "sales_db@finance-catalog" TO "IAM:sales_auditor"`, []string{"IAM:sales_auditor"}, "on_database"},
		{"GRANT DEBUG ON DATABASE d TO ROLE AWSIDC:analysts", []string{"awsidc:analysts"}, "on_database"},
		{"GRANT DEBUG ON DATABASE d TO ADMINS OF ACCOUNT '1234'", []string{"admins_of_account:1234"}, "on_database"},
		{"REVOKE DEBUG ON DATABASE d FROM ADMINS OF ACCOUNT '1234'", []string{"admins_of_account:1234"}, "on_database"},
	} {
		grant, ok := singleStmt(t, tc.sql).(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt", tc.sql)
		}
		if grant.GrantedRoles == nil || len(grant.GrantedRoles.Items) != 1 {
			t.Fatalf("Parse(%q): granted permissions = %#v, want one", tc.sql, grant.GrantedRoles)
		}
		priv := grant.GrantedRoles.Items[0].(*nodes.AccessPriv)
		if want := strings.ToLower(strings.Fields(tc.sql)[1]); priv.PrivName != want {
			t.Errorf("Parse(%q): privilege = %q, want %q", tc.sql, priv.PrivName, want)
		}
		if priv.Loc.Start != strings.Index(strings.ToUpper(tc.sql), strings.ToUpper(priv.PrivName)) {
			t.Errorf("Parse(%q): privilege Loc = %+v, want its own token", tc.sql, priv.Loc)
		}
		if got := findDefElem(grant.Opt, tc.option) != nil; tc.option != "" && !got {
			t.Errorf("Parse(%q): Opt = %#v, want %s", tc.sql, grant.Opt, tc.option)
		}
		if len(tc.grantees) == 0 {
			rs, ok := grant.GranteeRoles.Items[0].(*nodes.RoleSpec)
			if !ok || nodes.RoleSpecType(rs.Roletype) != nodes.ROLESPEC_PUBLIC {
				t.Errorf("Parse(%q): grantee = %+v, want PUBLIC", tc.sql, grant.GranteeRoles.Items[0])
			}
			continue
		}
		if len(grant.GranteeRoles.Items) != 1 {
			t.Fatalf("Parse(%q): grantees = %#v, want one", tc.sql, grant.GranteeRoles)
		}
		rs, ok := grant.GranteeRoles.Items[0].(*nodes.RoleSpec)
		if !ok || rs.Rolename != tc.grantees[0] {
			t.Errorf("Parse(%q): grantee = %+v, want %q", tc.sql, grant.GranteeRoles.Items[0], tc.grantees[0])
		}
	}

	// ROLE names an Identity Center role, not a user of that name.
	role := singleStmt(t, "GRANT CONNECT TO ROLE AWSIDC:analysts").(*nodes.GrantRoleStmt)
	if rs, ok := role.GranteeRoles.Items[0].(*nodes.RoleSpec); !ok || nodes.RoleSpecType(rs.Roletype) != nodes.ROLESPEC_IDC_ROLE {
		t.Errorf("TO ROLE AWSIDC:analysts: grantee = %+v, want ROLESPEC_IDC_ROLE", role.GranteeRoles.Items[0])
	}
	user := singleStmt(t, "GRANT CONNECT TO AWSIDC:analysts").(*nodes.GrantRoleStmt)
	if rs, ok := user.GranteeRoles.Items[0].(*nodes.RoleSpec); !ok || nodes.RoleSpecType(rs.Roletype) != nodes.ROLESPEC_CSTRING {
		t.Errorf("TO AWSIDC:analysts: grantee = %+v, want ROLESPEC_CSTRING", user.GranteeRoles.Items[0])
	}

	for _, sql := range []string{
		// DEBUG's revoke starts at DEBUG, with no role-option prefix.
		"REVOKE GRANT OPTION FOR DEBUG ON DATABASE d FROM AWSIDC:alice",
		"REVOKE ADMIN OPTION FOR DEBUG ON DATABASE d FROM AWSIDC:alice",
		// A CONNECT and DEBUG grantee is <prefix>:<name>, ROLE <prefix>:<name>,
		// PUBLIC (CONNECT) or the admins of an account (DEBUG).
		"GRANT CONNECT TO alice",
		"GRANT CONNECT TO a:b:c",
		`GRANT CONNECT TO "alice"`,
		`GRANT CONNECT TO "a:b:c"`,
		`GRANT DEBUG ON DATABASE d TO "alice"`,
		`GRANT CONNECT TO ROLE "alice"`,
		"GRANT CONNECT TO ROLE alice",
		"GRANT CONNECT TO a, b",
		"GRANT CONNECT ON TABLE t TO a:b",
		"GRANT CONNECT TO ADMINS OF ACCOUNT '1'",
		"GRANT DEBUG ON DATABASE d TO alice",
		"GRANT DEBUG ON DATABASE d TO PUBLIC",
		"GRANT DEBUG ON DATABASE d TO ADMINS ACCOUNT '1'",
		"GRANT DEBUG ON DATABASE d.schema TO AWSIDC:alice",
		"GRANT DEBUG ON DATABASE d TO ADMINS OF ACCOUNT 1",
		"REVOKE CONNECT FROM a:b",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftGrantScopedTemplates pins the TEMPLATES scoped form and the
// privilege list that precedes FOR, which r_GRANT and r_REVOKE write.
func TestRedshiftGrantScopedTemplates(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		objtype nodes.ObjectType
	}{
		{"GRANT ALTER FOR TEMPLATES IN SCHEMA s TO a", nodes.OBJECT_TEMPLATE},
		{"GRANT ALTER, DROP, USAGE FOR TEMPLATES IN SCHEMA s DATABASE d TO a, ROLE r", nodes.OBJECT_TEMPLATE},
		{"GRANT ALL FOR TEMPLATES IN DATABASE d TO a", nodes.OBJECT_TEMPLATE},
		{"REVOKE ALTER FOR TEMPLATES IN SCHEMA s FROM a", nodes.OBJECT_TEMPLATE},
		{"GRANT SELECT, SELECT FOR TABLES IN SCHEMA s TO a", nodes.OBJECT_TABLE},
		{"REVOKE SELECT, INSERT FOR TABLES IN DATABASE d FROM a", nodes.OBJECT_TABLE},
		{"REVOKE GRANT OPTION CREATE FOR SCHEMAS IN DATABASE d FROM a", nodes.OBJECT_SCHEMA},
		{"REVOKE GRANT OPTION CREATE FOR COPY JOBS IN DATABASE d FROM a", nodes.OBJECT_COPY_JOB},
	} {
		grant, ok := singleStmt(t, tc.sql).(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt", tc.sql)
		}
		if grant.Targtype != nodes.ACL_TARGET_SCOPED {
			t.Errorf("Parse(%q): targtype = %v, want ACL_TARGET_SCOPED", tc.sql, grant.Targtype)
		}
		if grant.Objtype != tc.objtype {
			t.Errorf("Parse(%q): objtype = %v, want %v", tc.sql, grant.Objtype, tc.objtype)
		}
	}

	// r_GRANT writes username ( WITH GRANT OPTION )? per grantee; the AST keeps
	// one option for the statement, as PostgreSQL's does, so the trailing one is
	// what it reads.
	if grant := singleStmt(t, "GRANT CREATE FOR SCHEMAS IN DATABASE d TO a WITH GRANT OPTION").(*nodes.GrantStmt); !grant.GrantOption {
		t.Error("scoped grant: GrantOption = false, want true")
	}
	if grant := singleStmt(t, "GRANT CREATE FOR SCHEMAS IN DATABASE d TO a, b WITH GRANT OPTION").(*nodes.GrantStmt); !grant.GrantOption {
		t.Error("scoped grant to two grantees: GrantOption = false, want true")
	}
	if grant := singleStmt(t, "GRANT CREATE FOR SCHEMAS IN DATABASE d TO a, b").(*nodes.GrantStmt); grant.GrantOption {
		t.Error("scoped grant without an option: GrantOption = true, want false")
	}

	for _, sql := range []string{
		"GRANT ALTER FOR TEMPLATES IN a TO b",
		"GRANT ALTER FOR TEMPLATES IN SCHEMA TO a",
		"REVOKE ALTER FOR TEMPLATES IN SCHEMA a FROM",
		"REVOKE GRANT OPTION CREATE MODEL FROM a",
		"REVOKE GRANT_OPTION CREATE FOR SCHEMAS IN DATABASE d FROM a",
		// { SCHEMA name [ DATABASE name ] | DATABASE name } takes one scope.
		"GRANT ALTER FOR TEMPLATES IN DATABASE d DATABASE e TO a",
		"GRANT CREATE FOR SCHEMAS IN DATABASE d DATABASE e TO a",
		// Scoped grantees are a username or ROLE.
		"GRANT CREATE FOR SCHEMAS IN DATABASE d TO GROUP g",
		"GRANT CREATE FOR SCHEMAS IN DATABASE d TO PUBLIC",
		// The TEMPLATES alternative takes ALTER, DROP, USAGE or ALL.
		"GRANT SELECT FOR TEMPLATES IN SCHEMA s TO a",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftGrantRepeatedColumnLists pins r_GRANT's column-level alternative,
// which repeats a bare column list per privilege:
//
//	{ SELECT | UPDATE } ( column [, ...] ) [ , ( column [, ...] ) ... ]
func TestRedshiftGrantRepeatedColumnLists(t *testing.T) {
	for _, sql := range []string{
		"GRANT SELECT(a) ON t TO a",
		"GRANT SELECT(a), (b) ON t TO a",
		"GRANT SELECT(a), (b), (c) ON t TO a",
		"GRANT UPDATE(a, b), (c) ON t TO a",
		"REVOKE SELECT(a), (b) ON t FROM a",
		"GRANT ALL(a, b) ON t TO a",
	} {
		stmt, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		if stmt == nil || len(stmt.Items) != 1 {
			t.Fatalf("Parse(%q): statements = %#v, want one", sql, stmt)
		}
	}

	// Only SELECT and UPDATE repeat a bare list, and ALL takes one list.
	for _, sql := range []string{
		"GRANT INSERT(a), (b) ON t TO a",
		"GRANT ALL(a), (b) ON t TO a",
		"GRANT SELECT(a), () ON t TO a",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftFuncArgNameVersusType pins the docs' argname reading against the
// multiword types: `f(character varying)` is one parameter of that type, while
// a column-name keyword before a type names the parameter (`f(int int)`).
func TestRedshiftFuncArgNameVersusType(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		name string
	}{
		{"CREATE FUNCTION f(character varying) RETURNS int AS 's' LANGUAGE sql", ""},
		{"CREATE FUNCTION f(national character) RETURNS int AS 's' LANGUAGE sql", ""},
		{"CREATE FUNCTION f(bit varying) RETURNS int AS 's' LANGUAGE sql", ""},
		{"CREATE FUNCTION f(interval year) RETURNS int AS 's' LANGUAGE sql", ""},
		{"CREATE FUNCTION f(timestamp with time zone) RETURNS int AS 's' LANGUAGE sql", ""},
		{"CREATE FUNCTION f(IN character varying) RETURNS int AS 's' LANGUAGE sql", ""},
		{"CREATE FUNCTION f(IN int int) RETURNS int AS 's' LANGUAGE sql", "int"},
		{"CREATE FUNCTION f(int int) RETURNS int AS 's' LANGUAGE sql", "int"},
		{"CREATE FUNCTION f(a int, b varchar(10)) RETURNS int AS 's' LANGUAGE sql", "a"},
	} {
		fn, ok := singleStmt(t, tc.sql).(*nodes.CreateFunctionStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected CreateFunctionStmt", tc.sql)
		}
		fp := fn.Parameters.Items[0].(*nodes.FunctionParameter)
		if fp.Name != tc.name {
			t.Errorf("Parse(%q): first parameter name = %q, want %q", tc.sql, fp.Name, tc.name)
		}
	}
}

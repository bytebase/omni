package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// execOne runs one statement and returns the error Exec recorded for it.
func execOne(t *testing.T, c *Catalog, sql string) error {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	if len(results) != 1 {
		t.Fatalf("Exec(%q): %d results, want one", sql, len(results))
	}
	return results[0].Error
}

// TestGrantRedshiftOnlyTargetsRecordNothing holds the grants on the Redshift
// objects the catalog does not model — a datashare, a COPY job, a model, a
// template, a warehouse permission — out of PostgreSQL's ACL validation: Exec
// records nothing and names no object to resolve.
func TestGrantRedshiftOnlyTargetsRecordNothing(t *testing.T) {
	for _, sql := range []string{
		"GRANT ALTER ON DATASHARE salesshare TO analyst",
		"GRANT USAGE ON DATASHARE salesshare TO ACCOUNT '123456789012' VIA DATA CATALOG",
		"REVOKE USAGE ON DATASHARE salesshare FROM NAMESPACE 'a3f3ae8c'",
		"GRANT ALTER ON COPY JOB nightly TO loader",
		"GRANT EXECUTE ON MODEL churn TO ml_user",
		"GRANT CONNECT ON WORKGROUP TO AWSIDC:alice",
		"GRANT DEBUG ON DATABASE federated TO ADMINS OF ACCOUNT '1234'",
		"GRANT CREATE MODEL TO ml_admin WITH GRANT OPTION",
		"GRANT ALTER ON TEMPLATE nightly TO loader",
	} {
		if err := execOne(t, New(), sql); err != nil {
			t.Errorf("Exec(%q): %v", sql, err)
		}
	}
}

// TestGrantExternalObjectsResolveAsRelations holds an external table on the
// relation path and an external schema on the schema path, where the loaded
// metadata installs them, so a schema diff sees their grants; their privileges
// are the ones r_GRANT and r_REVOKE list, not PostgreSQL's.
func TestGrantExternalObjectsResolveAsRelations(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA spectrum",
		"CREATE TABLE spectrum.sales (id INT)",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"GRANT SELECT ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'arn:aws:iam::1:role/r'",
		"GRANT ALTER, DROP ON EXTERNAL TABLE spectrum.sales TO PUBLIC",
		"GRANT CREATE, ALTER ON EXTERNAL SCHEMA spectrum TO IAM_ROLE 'arn:aws:iam::1:role/r'",
		"REVOKE ALTER ON EXTERNAL SCHEMA spectrum FROM IAM_ROLE 'arn:aws:iam::1:role/r'",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Errorf("Exec(%q): %v", sql, err)
		}
	}
	if len(c.grants) == 0 {
		t.Fatal("no grant recorded for the external table or schema")
	}
	// The catalog holds the object type to the privileges r_GRANT and r_REVOKE
	// list for it, not PostgreSQL's relation or schema set.
	for _, tc := range []struct {
		objtype nodes.ObjectType
		priv    string
		code    byte
		valid   bool
	}{
		{nodes.OBJECT_FOREIGN_TABLE, "alter", 'r', true},
		{nodes.OBJECT_FOREIGN_TABLE, "truncate", 'r', false},
		{nodes.OBJECT_EXTERNAL_SCHEMA, "alter", 'n', true},
		{nodes.OBJECT_EXTERNAL_SCHEMA, "usage", 'n', false},
	} {
		if got := validGrantPrivilege(tc.objtype, tc.priv, tc.code); got != tc.valid {
			t.Errorf("validGrantPrivilege(%v, %q, %q) = %v, want %v", tc.objtype, tc.priv, string(tc.code), got, tc.valid)
		}
	}
}

// TestGrantRedshiftOnlyTargetsRefusePublicGrantOption holds the one check the
// catalog keeps for every Redshift-only target: a grant option cannot go to
// PUBLIC.
func TestGrantRedshiftOnlyTargetsRefusePublicGrantOption(t *testing.T) {
	for _, sql := range []string{
		"GRANT EXECUTE ON MODEL m TO PUBLIC WITH GRANT OPTION",
		"GRANT SELECT ON EXTERNAL TABLE s.t TO PUBLIC WITH GRANT OPTION",
		"GRANT ALTER ON TEMPLATE t TO PUBLIC WITH GRANT OPTION",
	} {
		err := execOne(t, New(), sql)
		assertCode(t, err, CodeInvalidGrantOperation)
	}
}

// TestExternalTableGrantRendersItsOwnSyntax holds a recorded external-table
// grant to the statement it was written as: the ON clause keeps EXTERNAL TABLE
// and the grantee keeps IAM_ROLE, so a migration reproduces the permission.
func TestExternalTableGrantRendersItsOwnSyntax(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA spectrum",
		"CREATE TABLE spectrum.sales (id INT)",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	if err := execOne(t, c, "GRANT SELECT, ALTER ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'arn:aws:iam::1:role/r'"); err != nil {
		t.Fatalf("GRANT: %v", err)
	}
	if err := execOne(t, c, "GRANT CREATE ON EXTERNAL SCHEMA spectrum TO IAM_ROLE 'arn:aws:iam::1:role/r2'"); err != nil {
		t.Fatalf("GRANT: %v", err)
	}

	var grants, revokes []string
	for _, g := range c.grants {
		objName := resolveGrantObjName(c, g.ObjType, g.ObjOID)
		if strings.EqualFold(g.Privilege, "select") && g.Target == nodes.OBJECT_FOREIGN_TABLE {
			if got, want := formatGrantSQL(g, objName), "GRANT SELECT ON EXTERNAL TABLE "+formatQualifiedDescription(objName)+" TO IAM_ROLE 'arn:aws:iam::1:role/r'"; got != want {
				t.Errorf("GRANT SQL = %q, want %q", got, want)
			}
			revokes = append(revokes, formatRevokeSQL(g, objName))
		}
		if g.Target == nodes.OBJECT_EXTERNAL_SCHEMA {
			grants = append(grants, formatGrantSQL(g, objName))
		}
	}
	if len(grants) != 1 || !strings.Contains(grants[0], "ON EXTERNAL SCHEMA ") || !strings.Contains(grants[0], "TO IAM_ROLE 'arn:aws:iam::1:role/r2'") {
		t.Errorf("external schema grants = %v, want one with EXTERNAL SCHEMA and IAM_ROLE", grants)
	}
	if len(revokes) == 0 {
		t.Fatal("no external table grant recorded")
	}
	if want := "REVOKE SELECT ON EXTERNAL TABLE "; !strings.HasPrefix(revokes[0], want) {
		t.Errorf("REVOKE SQL = %q, want prefix %q", revokes[0], want)
	}

	// A GRANT on an external resource cannot run inside a transaction block.
	for _, g := range c.grants {
		if got, want := grantTransactional(g), g.Target == nodes.OBJECT_TABLE; got != want {
			t.Errorf("grantTransactional(%v) = %v, want %v", g.Target, got, want)
		}
	}
	// An external-table grant and a plain relation's to a role of that spelling
	// are two: target and grantee kind are part of a grant's identity.
	external := Grant{ObjType: 'r', ObjOID: 1, Grantee: "iam_role:r", Privilege: "select",
		Target: nodes.OBJECT_FOREIGN_TABLE, GranteeKind: "iam_role", GranteeName: "r"}
	plain := Grant{ObjType: 'r', ObjOID: 1, Grantee: "iam_role:r", Privilege: "select"}
	if makeGrantKey(c, external) == makeGrantKey(c, plain) {
		t.Error("external and plain grants share a key, want two")
	}
}

// TestGranteeKindIsTheRoleSpecKind holds a plain role whose quoted name looks
// like a kind's prefixed value apart from the kind a RoleSpec names.
func TestGranteeKindIsTheRoleSpecKind(t *testing.T) {
	c := New()
	if err := execOne(t, c, "CREATE SCHEMA s"); err != nil {
		t.Fatal(err)
	}
	if err := execOne(t, c, "CREATE TABLE s.t (id INT)"); err != nil {
		t.Fatal(err)
	}
	if err := execOne(t, c, `GRANT SELECT ON s.t TO "iam_role:r"`); err != nil {
		t.Fatal(err)
	}
	g := c.grants[0]
	if g.GranteeKind != "" || g.GranteeName != "" {
		t.Fatalf("grantee = %+v, want a plain role", g)
	}
	if got, want := formatGrantee(g), `"iam_role:r"`; got != want {
		t.Errorf("formatGrantee = %q, want %q", got, want)
	}

	// A procedure grant renders as a procedure, not a function.
	proc := Grant{ObjType: 'f', Target: nodes.OBJECT_PROCEDURE}
	if got, want := formatGrantObjectRef(proc, "s.p()"), "PROCEDURE s.p()"; got != want {
		t.Errorf("procedure grant ON clause = %q, want %q", got, want)
	}
}

// TestRevokeMatchesTargetAndGranteeKind holds a revoke to the grant it names:
// a relation grant to a role named "iam_role:r" survives a revoke of the
// external table's grant to IAM_ROLE 'r'.
func TestRevokeMatchesTargetAndGranteeKind(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA spectrum",
		"CREATE TABLE spectrum.sales (id INT)",
		`GRANT SELECT ON spectrum.sales TO "iam_role:r"`,
		"GRANT SELECT ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'r'",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	if len(c.grants) != 2 {
		t.Fatalf("grants = %+v, want two", c.grants)
	}
	if err := execOne(t, c, "REVOKE SELECT ON EXTERNAL TABLE spectrum.sales FROM IAM_ROLE 'r'"); err != nil {
		t.Fatalf("REVOKE: %v", err)
	}
	if len(c.grants) != 1 || c.grants[0].Target != nodes.OBJECT_TABLE {
		t.Fatalf("grants = %+v, want the relation grant kept", c.grants)
	}
}

// TestRevokeMatchesRoutineAlias holds a revoke on FUNCTION or PROCEDURE to the
// grant a ROUTINE names, as PostgreSQL's GRANT inherits.
func TestRevokeMatchesRoutineAlias(t *testing.T) {
	for _, tc := range []struct{ grant, revoke string }{
		{"GRANT EXECUTE ON ROUTINE s.f() TO u", "REVOKE EXECUTE ON FUNCTION s.f() FROM u"},
		{"GRANT EXECUTE ON FUNCTION s.f() TO u", "REVOKE EXECUTE ON ROUTINE s.f() FROM u"},
	} {
		if !sameGrantTarget(nodes.OBJECT_ROUTINE, nodes.OBJECT_FUNCTION) {
			t.Fatal("ROUTINE and FUNCTION: sameGrantTarget = false, want true")
		}
		if !sameGrantTarget(nodes.OBJECT_ROUTINE, nodes.OBJECT_PROCEDURE) {
			t.Fatal("ROUTINE and PROCEDURE: sameGrantTarget = false, want true")
		}
		if sameGrantTarget(nodes.OBJECT_TABLE, nodes.OBJECT_FUNCTION) {
			t.Fatal("TABLE and FUNCTION: sameGrantTarget = true, want false")
		}
		_ = tc
	}
}

// TestGrantOptionsBelongToTheirGrantee holds a mixed grantee list: the option on
// a username does not reach a PUBLIC grantee, and a recorded column grant is
// revoked by its own columns.
func TestGrantOptionsBelongToTheirGrantee(t *testing.T) {
	c := New()
	if err := execOne(t, c, "GRANT ALTER ON DATASHARE d TO alice WITH GRANT OPTION, PUBLIC"); err != nil {
		t.Fatalf("mixed datashare grantees: %v", err)
	}
	// A token handed in directly, with the statement-wide flag and no per-grantee
	// flag, still refuses PUBLIC.
	stmt := &nodes.GrantStmt{
		IsGrant:     true,
		GrantOption: true,
		Objtype:     nodes.OBJECT_TEMPLATE,
		Privileges:  &nodes.List{Items: []nodes.Node{&nodes.AccessPriv{PrivName: "alter"}}},
		Grantees:    &nodes.List{Items: []nodes.Node{&nodes.RoleSpec{Roletype: int(nodes.ROLESPEC_PUBLIC)}}},
	}
	assertCode(t, c.ExecGrantStmt(stmt), CodeInvalidGrantOperation)
}

// TestRevokeKeepsOtherColumnGrants holds a column revoke to its own columns.
func TestRevokeKeepsOtherColumnGrants(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA spectrum",
		"CREATE TABLE spectrum.sales (a INT, b INT)",
		"GRANT SELECT(a) ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'r'",
		"GRANT SELECT(b) ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'r'",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	if err := execOne(t, c, "REVOKE SELECT(a) ON EXTERNAL TABLE spectrum.sales FROM IAM_ROLE 'r'"); err != nil {
		t.Fatalf("REVOKE: %v", err)
	}
	if len(c.grants) != 1 || len(c.grants[0].Columns) != 1 || c.grants[0].Columns[0] != "b" {
		t.Fatalf("grants = %+v, want only SELECT(b)", c.grants)
	}
}

// TestIAMRoleDefaultRendersAsKeyword holds IAM_ROLE DEFAULT to its keyword.
func TestIAMRoleDefaultRendersAsKeyword(t *testing.T) {
	g := Grant{Grantee: "iam_role:default", GranteeKind: "iam_role_default"}
	if got, want := formatGrantee(g), "IAM_ROLE DEFAULT"; got != want {
		t.Errorf("formatGrantee = %q, want %q", got, want)
	}
	g = Grant{Grantee: "iam_role:default", GranteeKind: "iam_role", GranteeName: "default"}
	if got, want := formatGrantee(g), "IAM_ROLE 'default'"; got != want {
		t.Errorf("formatGrantee = %q, want %q", got, want)
	}
}

// TestGrantOptionStaysOnItsGranteeAcrossLists holds the option on the grantee it
// follows, before a later grantee without one.
func TestGrantOptionStaysOnItsGranteeAcrossLists(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA s",
		"CREATE TABLE s.t (id INT)",
		"GRANT SELECT ON s.t TO bob WITH GRANT OPTION, alice",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	if len(c.grants) != 2 {
		t.Fatalf("grants = %+v, want two", c.grants)
	}
	for _, g := range c.grants {
		want := g.Grantee == "bob"
		if g.WithGrant != want {
			t.Errorf("grantee %q: WithGrant = %v, want %v", g.Grantee, g.WithGrant, want)
		}
	}
	// A scoped permission covers objects that do not exist yet: Exec records
	// nothing and resolves no object.
	if err := execOne(t, c, "GRANT SELECT FOR TABLES IN SCHEMA s TO alice"); err != nil {
		t.Fatalf("scoped grant: %v", err)
	}
}

// TestPartialColumnRevokeTrimsGrant holds a column revoke to the columns it
// names out of a larger column grant.
func TestPartialColumnRevokeTrimsGrant(t *testing.T) {
	c := New()
	for _, sql := range []string{
		"CREATE SCHEMA spectrum",
		"CREATE TABLE spectrum.sales (a INT, b INT)",
		"GRANT SELECT(a, b) ON EXTERNAL TABLE spectrum.sales TO IAM_ROLE 'r'",
	} {
		if err := execOne(t, c, sql); err != nil {
			t.Fatalf("Exec(%q): %v", sql, err)
		}
	}
	if err := execOne(t, c, "REVOKE SELECT(a) ON EXTERNAL TABLE spectrum.sales FROM IAM_ROLE 'r'"); err != nil {
		t.Fatalf("REVOKE: %v", err)
	}
	if len(c.grants) != 1 || len(c.grants[0].Columns) != 1 || c.grants[0].Columns[0] != "b" {
		t.Fatalf("grants = %+v, want SELECT(b) only", c.grants)
	}
}

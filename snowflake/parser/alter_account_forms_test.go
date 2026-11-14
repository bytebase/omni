package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func parseOneAlterAccount(t *testing.T, sql string) *ast.AlterAccountStmt {
	t.Helper()
	file, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterAccountStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.AlterAccountStmt", sql, file.Stmts[0])
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Parse(%q): Loc.End = %d, want %d", sql, stmt.Loc.End, len(sql))
	}
	return stmt
}

// TestAlterAccountForms covers the ALTER ACCOUNT forms added beside the
// generic SET: organization user groups, contacts, ingress and maintenance
// policies, and the parenthesized LOGIN_IDP_REDIRECT and DOMAIN_NAMES values.
// https://docs.snowflake.com/en/sql-reference/sql/alter-account
// https://docs.snowflake.com/en/sql-reference/sql/alter-account-set-ingress-policy
func TestAlterAccountForms(t *testing.T) {
	s := parseOneAlterAccount(t, "ALTER ACCOUNT ADD ORGANIZATION USER GROUP g1")
	if s.Action != ast.AlterAccountAddOrganizationUserGroup || s.UserGroup.Normalize() != "G1" || s.Name != nil {
		t.Errorf("ADD: %+v", s)
	}
	s = parseOneAlterAccount(t, "ALTER ACCOUNT REMOVE ORGANIZATION USER GROUP g1")
	if s.Action != ast.AlterAccountRemoveOrganizationUserGroup {
		t.Errorf("REMOVE: %+v", s)
	}
	s = parseOneAlterAccount(t, "ALTER ACCOUNT SET CONTACT security = sec_contact, support = sup_contact")
	if s.Action != ast.AlterAccountSetContact || len(s.Contacts) != 2 || s.Contacts[1].Name != "SUPPORT" ||
		s.Contacts[1].Values[0] != "SUP_CONTACT" {
		t.Errorf("SET CONTACT: %+v", s)
	}
	s = parseOneAlterAccount(t, "ALTER ACCOUNT SET INGRESS POLICY security_db.policies.corp_only FORCE")
	if s.Action != ast.AlterAccountSetPolicy || s.PolicyKind != "INGRESS POLICY" || !s.Force {
		t.Errorf("INGRESS: %+v", s)
	}
	s = parseOneAlterAccount(t, "ALTER ACCOUNT SET MAINTENANCE POLICY p FORCE FOR ALL APPLICATIONS")
	if s.PolicyKind != "MAINTENANCE POLICY" || !s.Force || s.PolicyScope != "FOR ALL APPLICATIONS" {
		t.Errorf("MAINTENANCE: %+v", s)
	}
	s = parseOneAlterAccount(t, "ALTER ACCOUNT SET LOGIN_IDP_REDIRECT = (SNOWFLAKE_UI = my_saml, SNOWFLAKE_CLI = other) OBJECT_VISIBILITY = PRIVILEGED")
	if len(s.Options) != 2 || len(s.Options[0].Pairs) != 2 || s.Options[0].Pairs[1].Values[0] != "OTHER" {
		t.Errorf("LOGIN_IDP_REDIRECT: %+v", s.Options)
	}
	s = parseOneAlterAccount(t, "ALTER ACCOUNT acct1 SET DOMAIN_NAMES = ('a.example.com', 'b.example.com')")
	if len(s.Options) != 1 || len(s.Options[0].List) != 2 || s.Options[0].List[1].Value != "b.example.com" {
		t.Errorf("DOMAIN_NAMES: %+v", s.Options)
	}
	// A contact is a schema-level object.
	s = parseOneAlterAccount(t, "ALTER ACCOUNT SET CONTACT support = governance.contacts.helpdesk")
	if len(s.Contacts) != 1 || s.Contacts[0].Values[0] != "GOVERNANCE.CONTACTS.HELPDESK" {
		t.Errorf("qualified contact: %+v", s.Contacts)
	}
	parseOneAlterAccount(t, "ALTER ACCOUNT UNSET MAINTENANCE POLICY FOR ALL APPLICATIONS")
	parseOneAlterAccount(t, "ALTER ACCOUNT UNSET SESSION POLICY")
	// CONTACT followed by '=' is a parameter of that name.
	if s = parseOneAlterAccount(t, "ALTER ACCOUNT SET CONTACT = a"); s.Action != ast.AlterAccountSet {
		t.Errorf("CONTACT = a: %+v", s)
	}

	for _, sql := range []string{
		"ALTER ACCOUNT ADD ORGANIZATION USER GROUP db.g",
		"ALTER ACCOUNT ADD ORGANIZATION GROUP g",
		"ALTER ACCOUNT SET CONTACT security = 'x'",
		"ALTER ACCOUNT SET CONTACT security = c,",
		// MAINTENANCE POLICY writes FORCE before its scope, which it requires.
		"ALTER ACCOUNT SET MAINTENANCE POLICY p FOR ALL APPLICATIONS FORCE",
		"ALTER ACCOUNT SET MAINTENANCE POLICY p FORCE",
		"ALTER ACCOUNT SET MAINTENANCE POLICY p",
		"ALTER ACCOUNT UNSET MAINTENANCE POLICY",
		"ALTER ACCOUNT SET MAINTENANCE POLICY p FOR ALL DATABASES",
		"ALTER ACCOUNT SET FEATURE POLICY p",
		"ALTER ACCOUNT UNSET FEATURE POLICY",
		"ALTER ACCOUNT SET INGRESS POLICY p FOR ALL APPLICATIONS",
		"ALTER ACCOUNT SET SESSION POLICY p FOR ALL APPLICATIONS",
		// The scopes are the page's.
		"ALTER ACCOUNT SET FEATURE POLICY p FOR ALL",
		"ALTER ACCOUNT SET FEATURE POLICY p FOR ALL APPLICATIONS zzz",
		"ALTER ACCOUNT UNSET AUTHENTICATION POLICY FOR ALL USERS",
		// The parenthesized values' shapes.
		"ALTER ACCOUNT SET LOGIN_IDP_REDIRECT = ()",
		"ALTER ACCOUNT SET LOGIN_IDP_REDIRECT = (SNOWFLAKE_UI = 'x')",
		"ALTER ACCOUNT acct1 SET DOMAIN_NAMES = (a)",
		"ALTER ACCOUNT acct1 SET DOMAIN_NAMES = ('a',)",
		// DOMAIN_NAMES is set on a named account.
		"ALTER ACCOUNT SET DOMAIN_NAMES = ('a')",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
}

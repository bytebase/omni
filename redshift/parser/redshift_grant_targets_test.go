package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestRedshiftGrantTemplateTarget pins GRANT and REVOKE on a COPY template
// (https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html): a GRANT may
// qualify the first name, a REVOKE names a one-part template, and either takes
// only ALTER, DROP, USAGE or ALL PRIVILEGES, with no column list and no drop
// behavior.
func TestRedshiftGrantTemplateTarget(t *testing.T) {
	for sql, wantGrant := range map[string]bool{
		"GRANT ALTER ON TEMPLATE t TO u":                          true,
		"GRANT USAGE, DROP ON TEMPLATE db.s.t, t2 TO ROLE r":      true,
		"GRANT ALL ON TEMPLATE t TO PUBLIC WITH GRANT OPTION":     true,
		"REVOKE ALTER ON TEMPLATE t FROM u":                       false,
		"REVOKE GRANT OPTION FOR USAGE ON TEMPLATE t FROM ROLE r": false,
	} {
		grant, ok := singleStmt(t, sql).(*nodes.GrantStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantStmt", sql)
		}
		if grant.IsGrant != wantGrant {
			t.Errorf("Parse(%q): IsGrant = %v, want %v", sql, grant.IsGrant, wantGrant)
		}
		if grant.Objtype != nodes.OBJECT_TEMPLATE {
			t.Errorf("Parse(%q): objtype = %v, want OBJECT_TEMPLATE", sql, grant.Objtype)
		}
	}

	for _, sql := range []string{
		"GRANT SELECT ON TEMPLATE t TO u",
		"REVOKE INSERT ON TEMPLATE t FROM u",
		"GRANT ALL(a) ON TEMPLATE t TO u",
		`GRANT "ALTER" ON TEMPLATE t TO u`,
		"REVOKE ALTER ON TEMPLATE t FROM u CASCADE",
		"REVOKE ALTER ON TEMPLATE t FROM u RESTRICT",
		"REVOKE ALTER ON TEMPLATE db.s.t FROM u",
		"GRANT ALTER ON TEMPLATE a.b.c.d TO u",
		"GRANT ALTER ON TEMPLATE t TO IAM_ROLE 'x'",
		"GRANT ALTER ON TEMPLATE t TO NAMESPACE 'x'",
		"GRANT ALTER ON TEMPLATE t TO RLS POLICY p",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}

	// A grant that names no template reads TEMPLATE as the object's name, as it
	// did before templates had a target of their own.
	if grant := singleStmt(t, "GRANT ALTER ON TEMPLATE TO u").(*nodes.GrantStmt); grant.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("objtype = %v, want OBJECT_TABLE", grant.Objtype)
	}

	// A quoted public is a user, not the PUBLIC pseudo-role.
	grant := singleStmt(t, `GRANT ALTER ON TEMPLATE t TO "public"`).(*nodes.GrantStmt)
	if rs, ok := grant.Grantees.Items[0].(*nodes.RoleSpec); !ok || nodes.RoleSpecType(rs.Roletype) != nodes.ROLESPEC_CSTRING || rs.Rolename != "public" {
		t.Errorf(`TO "public": grantee = %+v, want a user named public`, grant.Grantees.Items[0])
	}
}

// TestRedshiftAssumeRoleGrant pins ASSUMEROLE's ON clause and its actions, as
// https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html writes them: GRANT
// takes a list of actions, REVOKE one. The spelling without ON still parses.
func TestRedshiftAssumeRoleGrant(t *testing.T) {
	for _, sql := range []string{
		"GRANT ASSUMEROLE ON 'arn:aws:iam::123456789012:role/r' TO u FOR COPY",
		"GRANT ASSUMEROLE ON 'a', 'b' TO ROLE r, u FOR EXTERNAL FUNCTION",
		"GRANT ASSUMEROLE ON ALL TO u FOR ALL, UNLOAD, CREATE MODEL",
		"REVOKE ASSUMEROLE ON 'arn:aws:iam::123456789012:role/r', DEFAULT FROM u FOR COPY",
		"GRANT ASSUMEROLE 'arn:aws:iam::123456789012:role/r' TO u FOR COPY",
		"GRANT ASSUMEROLE ON 'r' TO PUBLIC FOR COPY",
	} {
		grant, ok := singleStmt(t, sql).(*nodes.GrantRoleStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected GrantRoleStmt", sql)
		}
		if grant.Opt == nil || len(grant.Opt.Items) == 0 {
			t.Errorf("Parse(%q): options = %#v, want the FOR actions", sql, grant.Opt)
		}
	}

	for _, sql := range []string{
		"GRANT ASSUMEROLE ON 'r' TO u",
		"GRANT ASSUMEROLE ON 'r' TO u FOR",
		"GRANT ASSUMEROLE ON 'r' TO u FOR SELECT",
		"GRANT ASSUMEROLE ON a TO u FOR COPY",
		// REVOKE's FOR takes one action; each action is one alternative, not a
		// run of its words.
		"REVOKE ASSUMEROLE ON 'r' FROM u FOR COPY, ALL",
		"GRANT ASSUMEROLE ON 'r' TO u FOR ALL COPY",
		"GRANT ASSUMEROLE ON 'r' TO u FOR COPY COPY",
		"GRANT ASSUMEROLE ON 'r' TO u FOR EXTERNAL",
		"GRANT ASSUMEROLE ON 'r' TO u FOR CREATE",
		"GRANT ASSUMEROLE ON 'r' TO u FOR CREATE FUNCTION",
		`GRANT ASSUMEROLE ON 'r' TO u FOR "COPY"`,
		// DEFAULT and ALL stand alone among the roles.
		"GRANT ASSUMEROLE ON DEFAULT, ALL TO u FOR COPY",
		"REVOKE ASSUMEROLE ON ALL, DEFAULT FROM u FOR COPY",
		// The grantees are a user, a role, a group or PUBLIC.
		"GRANT ASSUMEROLE ON 'r' TO IAM_ROLE 'x' FOR COPY",
		"GRANT ASSUMEROLE ON 'r' TO NAMESPACE 'x' FOR COPY",
		"GRANT ASSUMEROLE ON 'r' TO RLS POLICY p FOR COPY",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}

	// ROLE alone names a user called role.
	if grant, ok := singleStmt(t, "GRANT ASSUMEROLE ON 'r' TO role FOR COPY").(*nodes.GrantRoleStmt); ok {
		if rs, ok := grant.GranteeRoles.Items[0].(*nodes.RoleSpec); !ok || rs.Rolename != "role" {
			t.Errorf("TO role: grantee = %+v, want the user role", grant.GranteeRoles.Items[0])
		}
	}
}

package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestGranteeNamedRole pins ROLE before a grantee: ROLE role_name where a
// name follows, CASCADE and RESTRICT included, since both are unreserved,
// and a user named role where the grantee ends there.
func TestGranteeNamedRole(t *testing.T) {
	for sql, want := range map[string]string{
		"GRANT SELECT ON t TO ROLE":                   "role",
		"GRANT SELECT ON t TO ROLE, a":                "role",
		"GRANT SELECT ON t TO ROLE r":                 "r",
		"GRANT SELECT ON t TO ROLE cascade":           "cascade",
		"GRANT SELECT ON t TO ROLE restrict":          "restrict",
		"REVOKE SELECT ON t FROM ROLE cascade":        "cascade",
		"REVOKE SELECT ON t FROM ROLE restrict":       "restrict",
		"REVOKE SELECT ON t FROM ROLE":                "role",
		"GRANT SELECT ON t TO ROLE WITH GRANT OPTION": "role",
	} {
		stmts, err := Parse(sql)
		if err != nil {
			t.Errorf("%s: %v", sql, err)
			continue
		}
		g, ok := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.GrantStmt)
		if !ok {
			t.Errorf("%s: %T", sql, stmts.Items[0].(*nodes.RawStmt).Stmt)
			continue
		}
		if got := g.Grantees.Items[0].(*nodes.RoleSpec).Rolename; got != want {
			t.Errorf("%s: grantee %q, want %q", sql, got, want)
		}
	}
}

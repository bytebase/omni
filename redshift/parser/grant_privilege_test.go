package parser

import "testing"

// TestGrantRequiresPrivilege: a GRANT or REVOKE on an object names its
// privileges (gram.y privileges: privilege_list | ALL ...).
func TestGrantRequiresPrivilege(t *testing.T) {
	for _, sql := range []string{"GRANT ON t TO u", "REVOKE ON t FROM u", "GRANT SELECT, ON t TO u"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
	for _, sql := range []string{"GRANT SELECT ON t TO u", "GRANT ALL ON t TO u", "GRANT USAGE ON SCHEMA s TO u", "REVOKE SELECT ON t FROM u"} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

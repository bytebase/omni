package parser

import "testing"

// TestGrantRequiresPrivilege: a GRANT or REVOKE on an object names its
// privileges (gram.y privileges: privilege_list | ALL ...), and fails where
// its target or WITH GRANT OPTION is cut short.
func TestGrantRequiresPrivilege(t *testing.T) {
	for _, sql := range []string{"GRANT ON t TO u", "REVOKE ON t FROM u", "GRANT SELECT, ON t TO u", "GRANT SELECT ON t TO u WITH GRANT", "GRANT SELECT ON LARGE OBJECT TO u", "GRANT USAGE ON FOREIGN DATA a TO u", "GRANT USAGE ON FOREIGN a TO u", "GRANT SELECT ON ALL x IN SCHEMA s TO u", "GRANT SELECT ON ALL TABLES SCHEMA s TO u"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
	for _, sql := range []string{"GRANT SELECT ON t TO u", "GRANT ALL ON t TO u", "GRANT USAGE ON SCHEMA s TO u", "REVOKE SELECT ON t FROM u", "GRANT SELECT ON t TO u WITH GRANT OPTION", "GRANT SELECT ON LARGE OBJECT 1, 2 TO u", "GRANT SELECT ON ALL TABLES IN SCHEMA s TO u"} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

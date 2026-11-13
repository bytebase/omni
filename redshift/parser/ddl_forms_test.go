package parser

import "testing"

// TestCreatedbOptionNeedsValue: createdb_opt_item takes a NumericOnly,
// opt_boolean_or_string or DEFAULT after its name, so an option cut short
// fails, also ALTER DATABASE's OWNER without TO.
func TestCreatedbOptionNeedsValue(t *testing.T) {
	for _, sql := range []string{"ALTER DATABASE d OWNER", "ALTER DATABASE d OWNER =", "ALTER DATABASE d OWNER +", "CREATE DATABASE d ENCODING", "ALTER DATABASE d CONNECTION 5"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted", sql)
		}
	}
	for _, sql := range []string{"ALTER DATABASE d OWNER x", "ALTER DATABASE d OWNER = 1.5", "ALTER DATABASE d OWNER TO x", "CREATE DATABASE d OWNER = x ENCODING 'UTF8' CONNECTION LIMIT -1", "ALTER DATABASE d WITH ALLOW_CONNECTIONS true"} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

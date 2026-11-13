package parser

import (
	"slices"
	"testing"
)

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

// TestCollectAfterStrictForms: the strict role options and the %TYPE
// lookahead keep completion's candidates.
func TestCollectAfterStrictForms(t *testing.T) {
	if cs := Collect("ALTER ROLE r ", 13); len(cs.Tokens) == 0 {
		t.Error("Collect after ALTER ROLE r offers nothing")
	}
	for _, c := range []struct {
		sql string
		off int
		tok int
	}{{"ALTER SEQUENCE s OWNER TO r", 17, OWNER}, {"ALTER SEQUENCE s ", 17, CYCLE}, {"ALTER DATABASE d OWNER ", 23, TO}} {
		if cs := Collect(c.sql, c.off); !cs.HasToken(c.tok) {
			t.Errorf("Collect(%q, %d) does not offer token %d", c.sql, c.off, c.tok)
		}
	}
	q := "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'select 1'"
	if cs := Collect(q, 28); !slices.Contains(cs.Tokens, IDENT) {
		t.Errorf("Collect(%q, 28) offers no type name", q)
	}
}

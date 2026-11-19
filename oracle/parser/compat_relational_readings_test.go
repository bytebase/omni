package parser

import (
	"errors"
	"testing"
)

// TestSupplementalLoggingPropertyReading pins the reading this change adds
// at the bytes compat.tsv records the pinned Oracle refusing at.
//
// A supplemental_logging_props clause is one of relational_properties
// (CREATE-TABLE.bnf lines 58-65 and 183), so SUPPLEMENTAL LOG GROUP a(a)
// stands in a table's column list. Oracle reads it and then refuses the
// closing ')' of a list that holds no column, as it refuses one that holds
// constraints alone; a clause broken inside itself, and a token after it
// that is neither ',' nor ')', answer at their own byte.
func TestSupplementalLoggingPropertyReading(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// The list ends with no column in it.
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a))`, 42},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(ALL) COLUMNS)`, 49},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(PRIMARY KEY) COLUMNS)`, 57},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(ALL))`, 41},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a) ALWAYS)`, 49},
		// A clause broken inside itself.
		{`CREATE TABLE a(SUPPLEMENTAL LOG a(a))`, 32},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a())`, 40},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUPING a(a))`, 32},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a LOG))`, 42},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a zzz))`, 42},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a a))`, 42},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a,))`, 42},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a NO LOG))`, 49},
		{`CREATE TABLE a(SUPPLEMENTAL LOG(ALL) COLUMNS)`, 31},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA() COLUMNS)`, 37},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(PRIMARY) COLUMNS)`, 44},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(FOREIGN) COLUMNS)`, 44},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(UNITE) COLUMNS)`, 37},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(ALL ALL) COLUMNS)`, 41},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATABASE(ALL) COLUMNS)`, 32},
		// A token after the clause that is neither ',' nor ')'.
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a) zzz)`, 43},
		{`CREATE TABLE a(SUPPLEMENTAL LOG GROUP a(a) ALWAYS zzz)`, 50},
		{`CREATE TABLE a(SUPPLEMENTAL LOG DATA(ALL) COLUMNS zzz)`, 50},
	}
	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.want {
			t.Errorf("Parse(%q): refused at %d, want %d", tc.sql, pe.Position, tc.want)
		}
	}

	// The controls: the same clause after a column, which Oracle takes, and
	// an object table's object_properties, which took it all along.
	for _, sql := range []string{
		`CREATE TABLE a (b NUMBER, SUPPLEMENTAL LOG GROUP a(b))`,
		`CREATE TABLE a (b NUMBER, SUPPLEMENTAL LOG DATA(ALL) COLUMNS)`,
		`CREATE TABLE a OF a(SUPPLEMENTAL LOG GROUP a(a))`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}
}

// TestAlterSessionSettingNameReading pins the byte ALTER SESSION SET
// answers for a setting name no '=' follows.
//
// Oracle refuses a name alter_session_set_clause does not name at the name
// itself (ORA-02248, ALTER SESSION SET a 1 at a), so the name of a setting
// after another is refused there too (ALTER SESSION SET EDITION = a zzz at
// zzz), not at the end of the text where the '=' it wanted is missing. A
// setting the page names is refused at the token after it (ALTER SESSION
// SET EDITION a at a).
func TestAlterSessionSettingNameReading(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		{`ALTER SESSION SET a 1`, 18},
		{`ALTER SESSION SET EDITION = a zzz`, 30},
		{`ALTER SESSION SET CONTAINER = a zzz`, 32},
		{`ALTER SESSION SET CONSTRAINT = IMMEDIATE zzz`, 41},
		{`ALTER SESSION SET ROW ARCHIVAL VISIBILITY = ACTIVE zzz`, 51},
		{`ALTER SESSION SET CONTAINER = a SERVICE = a zzz`, 44},
		{`ALTER SESSION SET CURRENT_SCHEMA = a zzz`, 37},
		// A setting the page names reports the token after the name.
		{`ALTER SESSION SET EDITION a`, 26},
		{`ALTER SESSION SET CONTAINER a`, 28},
		{`ALTER SESSION SET ROW ARCHIVAL VISIBILITY ACTIVE`, 42},
		{`ALTER SESSION SET CONSTRAINT IMMEDIATE`, 29},
		{`ALTER SESSION SET CONTAINER = a SERVICE a`, 40},
	}
	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.want {
			t.Errorf("Parse(%q): refused at %d, want %d", tc.sql, pe.Position, tc.want)
		}
	}

	// The controls: the settings the page names, and a plain parameter,
	// parse.
	for _, sql := range []string{
		`ALTER SESSION SET EDITION = a`,
		`ALTER SESSION SET CONSTRAINT = IMMEDIATE`,
		`ALTER SESSION SET ROW ARCHIVAL VISIBILITY = ACTIVE`,
		`ALTER SESSION SET USE_PRIVATE_OUTLINES = TRUE`,
		`ALTER SESSION SET a = 1`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}
}

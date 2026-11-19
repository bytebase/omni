package parser

import (
	"errors"
	"slices"
	"testing"
)

// TestAlterParameterReading holds ALTER SESSION SET, ALTER SYSTEM SET and
// ALTER SYSTEM RESET to the parameter names V$PARAMETER lists on the pinned
// Oracle, and to the values V$PARAMETER_VALID_VALUES lists for them (see
// parameters.go for the queries and the server).
//
// A name the list does not hold is refused at the name: Oracle answers
// ORA-02248 there for ALTER SESSION SET and ORA-02065 for ALTER SYSTEM SET.
// A name it holds is read on, so a value where the '=' belongs is refused at
// the value (ALTER SESSION SET optimizer_mode ALL_ROWS, ORA-00927), and a
// value the parameter cannot take is refused at the value too. ALTER SYSTEM
// RESET answers at the start of the statement for a name it does not hold.
func TestAlterParameterReading(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// No parameter by these names.
		{`ALTER SESSION SET a = 1`, 18},
		{`ALTER SESSION SET a 1`, 18},
		{`ALTER SESSION SET a =`, 18},
		{`ALTER SESSION SET a = 1 zzz`, 18},
		{`ALTER SESSION SET a = 1 a = 1`, 18},
		{`ALTER SESSION SET EDITIONABLE = a`, 18},
		{`ALTER SESSION SET ARCHIVAL VISIBILITY = ACTIVE`, 18},
		{`ALTER SYSTEM SET a = 1`, 17},
		{`ALTER SYSTEM SET a 1`, 17},
		{`ALTER SYSTEM SET a = 1 SCOPE = MEMORY`, 17},
		{`ALTER SYSTEM SET a = 1 CONTAINER = ALL`, 17},
		{`ALTER SYSTEM RESET a`, 0},
		{`ALTER SYSTEM RESET`, 0},
		{`ALTER SYSTEM RESET a SCOPE = MEMORY`, 0},
		{`ALTER SYSTEM RESET a zzz`, 0},
		// A parameter, refused after its name.
		{`ALTER SESSION SET optimizer_mode ALL_ROWS`, 33},
		{`ALTER SYSTEM SET optimizer_mode ALL_ROWS`, 32},
		// A parameter's value outside V$PARAMETER_VALID_VALUES.
		{`ALTER SESSION SET optimizer_mode = BAD`, 35},
		{`ALTER SESSION SET cursor_sharing = FOO`, 35},
		{`ALTER SYSTEM SET optimizer_mode = BAD`, 34},
		{`ALTER SYSTEM SET cursor_sharing = FOO`, 34},
		{`ALTER SESSION SET optimizer_mode = 'BAD'`, 35},
		{`ALTER SYSTEM SET cursor_sharing = 'FOO'`, 34},
		{`ALTER SESSION SET "_NOSUCH_HIDDEN" = 1`, 18},
	}
	for _, c := range refusedAt {
		_, err := Parse(c.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", c.sql, err)
			continue
		}
		if pe.Position != c.want {
			t.Errorf("Parse(%q): refused at %d, want %d", c.sql, pe.Position, c.want)
		}
	}

	// A parameter's own value, in any case, and a free-text parameter's
	// value: V$PARAMETER_VALID_VALUES lists no words for db_unique_name.
	for _, sql := range []string{
		`ALTER SESSION SET optimizer_mode = ALL_ROWS`,
		`ALTER SESSION SET OPTIMIZER_MODE = all_rows`,
		`ALTER SESSION SET cursor_sharing = FORCE`,
		`ALTER SESSION SET db_unique_name = zzz`,
		`ALTER SESSION SET "optimizer_mode" = all_rows`,
		`ALTER SESSION SET optimizer_mode = 'ALL_ROWS'`,
		`ALTER SESSION SET cursor_sharing = 'force'`,
		`ALTER SYSTEM SET optimizer_mode = 'all_rows'`,
		`ALTER SESSION SET "_ORACLE_SCRIPT" = true`,
		`ALTER SYSTEM RESET "_ORACLE_SCRIPT"`,
		`ALTER SESSION SET nls_date_format = 'YYYY-MM-DD'`,
		`ALTER SESSION SET sort_area_size = 1`,
		`ALTER SYSTEM SET optimizer_mode = ALL_ROWS`,
		`ALTER SYSTEM SET db_unique_name = zzz`,
		`ALTER SYSTEM SET open_cursors = 300`,
		`ALTER SYSTEM RESET optimizer_mode`,
		`ALTER SYSTEM RESET optimizer_mode SCOPE = MEMORY`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}
}

// TestParameterLists holds the recorded lists to the shapes parameters.go
// documents: V$PARAMETER's 566 names, X$KSPPI's 7,137 hidden ones (7,700
// distinct), the two enumerations the compat rows meet, and a free-text
// parameter with no V$PARAMETER_VALID_VALUES row.
func TestParameterLists(t *testing.T) {
	if got := len(oracleParameters); got != 7700 {
		t.Errorf("oracleParameters holds %d names, want V$PARAMETER and X$KSPPI's 7,700", got)
	}
	for _, name := range []string{"OPTIMIZER_MODE", "CURSOR_SHARING", "NLS_DATE_FORMAT", "DB_UNIQUE_NAME", "_DMM_BLAS_LIBRARY", "_ORACLE_SCRIPT"} {
		if !isOracleParameter(name) {
			t.Errorf("isOracleParameter(%q) = false, want true", name)
		}
	}
	if isOracleParameter("a") {
		t.Error(`isOracleParameter("a") = true, want false`)
	}
	for name, words := range map[string][]string{
		"OPTIMIZER_MODE": {"RULE", "CHOOSE", "ALL_ROWS", "FIRST_ROWS_1000"},
		"CURSOR_SHARING": {"FORCE", "EXACT", "SIMILAR"},
	} {
		got := oracleParameterValues[name]
		for _, word := range words {
			if !slices.Contains(got, word) {
				t.Errorf("%s holds %v, want %q", name, got, word)
			}
		}
	}
	if _, ok := oracleParameterValues["DB_UNIQUE_NAME"]; ok {
		t.Error("DB_UNIQUE_NAME has a value list, want free text")
	}
}

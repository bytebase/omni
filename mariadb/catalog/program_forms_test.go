package catalog

import (
	"strings"
	"testing"
)

// TestProgramFormsKeepWhatWasWritten pins two forms the programs suite does not
// reach: a parameter DEFAULT is kept as written, so a plan made under
// NO_BACKSLASH_ESCAPES writes it back unchanged, and a FOR range whose bound
// starts with nested parentheses is a range, not the loop's query.
func TestProgramFormsKeepWhatWasWritten(t *testing.T) {
	c, err := LoadSQL("CREATE DATABASE d; USE d; SET sql_mode = 'NO_BACKSLASH_ESCAPES';\n" +
		"CREATE PROCEDURE p(x VARCHAR(10) DEFAULT 'it''s\\n') SELECT x;\n" +
		"CREATE PROCEDURE q() FOR i IN ((1))..3 DO SELECT i; END FOR;\n" +
		"CREATE PROCEDURE r() FOR i IN (SELECT 1 AS b) DO SELECT i.b; END FOR;")
	if err != nil {
		t.Fatal(err)
	}
	plan := renderCreateRoutine("d", c.GetDatabase("d").Procedures["p"], true)
	if want := "DEFAULT 'it''s\\n')"; !strings.Contains(plan, want) {
		t.Errorf("plan %q does not hold %q", plan, want)
	}
	for _, name := range []string{"q", "r"} {
		if c.GetDatabase("d").Procedures[name] == nil {
			t.Errorf("procedure %s was not created", name)
		}
	}
}

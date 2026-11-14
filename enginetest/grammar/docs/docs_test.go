package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testDialect classifies every placeholder as a name and starts a form at
// the DDL verbs. Snowflake's and Redshift's converters test their own
// dialects.
var testDialect = &Dialect{
	Grammar:  "TestDocs",
	Site:     "https://example.com",
	Classify: func(string) string { return ClassName },
	StatementWords: func(string) map[string]bool {
		return map[string]bool{"CREATE": true, "ALTER": true, "GRANT": true, "REVOKE": true}
	},
}

// TestAnEqualsSignAloneInAChoiceIsASign pins { TO | = }: a choice of the
// word or the sign, where KEY = with nothing after it is a value left out.
func TestAnEqualsSignAloneInAChoiceIsASign(t *testing.T) {
	got, err := Form(testDialect, "p", "SET <p> { TO | = } <v>")
	if err != nil || strings.Join(got, " | ") != "'SET' ID ('TO' | '=') ID" {
		t.Errorf("got %q, %v", got, err)
	}
	if _, err := Form(testDialect, "p", "SET <p> { A = | B }"); err == nil {
		t.Errorf("KEY = with no value in a choice was converted")
	}
}

// TestAValueNamingAFailedDefinitionFails pins a form whose <value> names a
// definition the page writes that did not convert: the form is not
// converted either, rather than reading the value as a name.
func TestAValueNamingAFailedDefinitionFails(t *testing.T) {
	pg := ConvertPage(testDialect, "alter-x", []string{"ALTER X <name> <option>\n\noption ::=\n  A | [ B\n"}, nil)
	if len(pg.Forms) != 0 || len(pg.Failed) != 2 || !strings.Contains(pg.Failed[1], `"option"`) {
		t.Errorf("forms %q, failed %q", pg.Forms, pg.Failed)
	}
}

// TestAVariantMaySpanLines pins a definition whose variants, between --
// lines, are whole statements written over several lines.
func TestAVariantMaySpanLines(t *testing.T) {
	pg := ConvertPage(testDialect, "alter-x", []string{"ALTER X <clause>\n\nclause ::=\nGRANT A\n  ON B\n--\nREVOKE A\n  ON B"}, nil)
	if len(pg.Frags) != 1 || strings.Join(pg.Frags[0].Alts, " | ") != "'GRANT' 'A' 'ON' 'B' | 'REVOKE' 'A' 'ON' 'B'" {
		t.Errorf("fragments %+v, failed %q", pg.Frags, pg.Failed)
	}
}

// TestAnOptionalRepeatKeepsItsParentheses pins [ X [ ... ] ]: (X+)?, since
// ANTLR reads X+? as a non-greedy X+.
func TestAnOptionalRepeatKeepsItsParentheses(t *testing.T) {
	got, err := Form(testDialect, "p", "CREATE X <name> [ <e> [ ... ] ]")
	if err != nil || strings.Join(got, " | ") != "'CREATE' 'X' ID (ID+)?" {
		t.Errorf("got %q, %v", got, err)
	}
}

// TestAListElementTakesItsOptionalValues pins [ <argname> ] [ <argmode> ]
// <argtype> [, ...]: each element after a comma may name its argument too.
func TestAListElementTakesItsOptionalValues(t *testing.T) {
	got, err := Form(testDialect, "p", "F ( [ <argname> ] [ <argmode> ] <argtype> [, ...] )")
	if err != nil || strings.Join(got, " | ") != "'F' '(' ID? ID? ID (',' ID? ID? ID)* ')'" {
		t.Errorf("got %q, %v", got, err)
	}
}

// TestAFragmentNamingAFailedDefinitionFails pins a form that names, in lower
// case, a definition the page writes that did not convert: the form is not
// converted either, rather than reading the word as a literal value.
func TestAFragmentNamingAFailedDefinitionFails(t *testing.T) {
	pg := ConvertPage(testDialect, "alter-x", []string{"ALTER X <name> { option | B }\n\noption ::=\n  A | [ B\n"}, nil)
	if len(pg.Forms) != 0 || len(pg.Failed) != 2 || !strings.Contains(pg.Failed[1], `"option"`) {
		t.Errorf("forms %q, failed %q", pg.Forms, pg.Failed)
	}
}

// TestAPinRoleKeepsTheOtherRows pins the row a converter with a PinRole
// writes: its own row, replaced or added, among the file's other lines, as
// an engine judged by its server keeps the server in default.
func TestAPinRoleKeepsTheOtherRows(t *testing.T) {
	file := filepath.Join(t.TempDir(), "pins.tsv")
	server := "# a header\ndefault\timage\tsha256:1\t1.0\n"
	if err := os.WriteFile(file, []byte(server), 0o644); err != nil {
		t.Fatal(err)
	}
	d := &Dialect{Command: "c", Site: "https://example.com/", Index: "ref/", PinRole: "docs"}
	row, err := pinFor(d, file, "sha256:2", "2026-01-01")
	if err != nil {
		t.Fatal(err)
	}
	want := "docs\thttps://example.com/ref/\tsha256:2\t2026-01-01"
	if b, _ := os.ReadFile(file); row != want || string(b) != server+want+"\n" {
		t.Fatalf("row %q, file %q", row, b)
	}
	// The same digest keeps the row; another replaces it in place.
	if row, _ := pinFor(d, file, "sha256:2", ""); row != want {
		t.Errorf("same digest: row %q", row)
	}
	if _, err := pinFor(d, file, "sha256:3", "2026-02-02"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(file); string(b) != server+"docs\thttps://example.com/ref/\tsha256:3\t2026-02-02\n" {
		t.Errorf("replaced: file %q", b)
	}
}

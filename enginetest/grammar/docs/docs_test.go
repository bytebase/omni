package docs

import (
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
	got, err := Form(testDialect, "SET <p> { TO | = } <v>")
	if err != nil || strings.Join(got, " | ") != "'SET' ID ('TO' | '=') ID" {
		t.Errorf("got %q, %v", got, err)
	}
	if _, err := Form(testDialect, "SET <p> { A = | B }"); err == nil {
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
	got, err := Form(testDialect, "CREATE X <name> [ <e> [ ... ] ]")
	if err != nil || strings.Join(got, " | ") != "'CREATE' 'X' ID (ID+)?" {
		t.Errorf("got %q, %v", got, err)
	}
}

// TestAListElementTakesItsOptionalValues pins [ <argname> ] [ <argmode> ]
// <argtype> [, ...]: each element after a comma may name its argument too.
func TestAListElementTakesItsOptionalValues(t *testing.T) {
	got, err := Form(testDialect, "F ( [ <argname> ] [ <argmode> ] <argtype> [, ...] )")
	if err != nil || strings.Join(got, " | ") != "'F' '(' ID? ID? ID (',' ID? ID? ID)* ')'" {
		t.Errorf("got %q, %v", got, err)
	}
}

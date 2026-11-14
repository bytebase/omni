package catalog

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/parser"
)

// TestShowCreateProcedureAnchoredParams pins the parameter list SHOW CREATE writes for
// MariaDB's ROW(...), TYPE OF and ROW TYPE OF: each name quoted, so a name with a dot or a
// space reads back as the one the definition gave.
func TestShowCreateProcedureAnchoredParams(t *testing.T) {
	c := newCatalogWithDB(t, "mydb")
	stmt := mustParseRoutine(t, "CREATE PROCEDURE p(a TYPE OF `t.name`.`c.name`, b ROW(`x y` INT), c ROW TYPE OF t) BEGIN END")
	if err := c.DefineProcedure(stmt); err != nil {
		t.Fatalf("DefineProcedure: %v", err)
	}
	got := c.ShowCreateProcedure("mydb", "p")
	want := "(a TYPE OF `t.name`.`c.name`, b ROW(`x y` INT), c ROW TYPE OF `t`)"
	if !strings.Contains(got, want) {
		t.Fatalf("ShowCreateProcedure = %q, want it to hold %q", got, want)
	}
	if _, err := parser.Parse(got); err != nil {
		t.Fatalf("SHOW CREATE output does not parse: %v\n%s", err, got)
	}
}

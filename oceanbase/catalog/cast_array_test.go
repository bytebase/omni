package catalog

import (
	"context"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
	"github.com/bytebase/omni/oceanbase/parser"
)

// castArrayForms are the casts to ARRAY the server keeps, and what it stores each as.
var castArrayForms = map[string]string{
	"ARRAY(INT)":                      "ARRAY(INT)",
	"ARRAY<INT>":                      "ARRAY(INT)",
	"ARRAY(BIGINT)":                   "ARRAY(BIGINT)",
	"ARRAY(BIGINT UNSIGNED)":          "ARRAY(BIGINT UNSIGNED)",
	"ARRAY(VARCHAR(3))":               "ARRAY(VARCHAR(3))",
	"ARRAY<VARCHAR(3)>":               "ARRAY(VARCHAR(3))",
	"ARRAY(VARCHAR(3) CHARSET ascii)": "ARRAY(VARCHAR(3))",
	"ARRAY(FLOAT)":                    "ARRAY(FLOAT)",
	"ARRAY<ARRAY<INT>>":               "ARRAY(ARRAY(INT))",
	"ARRAY(ARRAY(INT))":               "ARRAY(ARRAY(INT))",
	"ARRAY<ARRAY(INT)>":               "ARRAY(ARRAY(INT))",
	"ARRAY(ARRAY<INT>)":               "ARRAY(ARRAY(INT))",
	"ARRAY<ARRAY<ARRAY<INT>>>":        "ARRAY(ARRAY(ARRAY(INT)))",
	"ARRAY(TINYINT(1))":               "ARRAY(TINYINT)",
	"ARRAY(BOOL)":                     "ARRAY(TINYINT)",
	"ARRAY(BIGINT(5))":                "ARRAY(BIGINT)",
	"ARRAY(INT ZEROFILL)":             "ARRAY(INT UNSIGNED)",
	"ARRAY(INTEGER)":                  "ARRAY(INT)",
	"ARRAY(SMALLINT)":                 "ARRAY(SMALLINT)",
	"ARRAY(DOUBLE(5,2))":              "ARRAY(DOUBLE)",
	"ARRAY(FLOAT(5,2))":               "ARRAY(FLOAT)",
	"ARRAY(FLOAT(30))":                "ARRAY(FLOAT)",
	"ARRAY(REAL)":                     "ARRAY(DOUBLE)",
	"ARRAY(NVARCHAR(3))":              "ARRAY(VARCHAR(3))",
	"ARRAY(VECTOR(3))":                "ARRAY(VECTOR(3))",
	"ARRAY(VARCHAR(0))":               "ARRAY(VARCHAR(0))",
}

func deparsedCast(t *testing.T, sql string) string {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	return deparse.Deparse(list.Items[0].(*nodes.SelectStmt).TargetList[0])
}

// castArrayFormsTheServerRefuses are elements OceanBase 4.3.5.6 answers 1235 for. They still
// print so that the text reads back.
var castArrayFormsTheServerRefuses = map[string]string{
	"ARRAY(ENUM('a','b'))": "ARRAY(ENUM('a','b'))",
	"ARRAY(SET('a','b'))":  "ARRAY(SET('a','b'))",
	"ARRAY(ENUM('it''s'))": "ARRAY(ENUM('it\\'s'))",
	"ARRAY(DECIMAL(5,2))":  "ARRAY(DECIMAL(5,2))",
	"ARRAY(VARBINARY(0))":  "ARRAY(VARBINARY(0))",
}

// TestCastArrayKeepsItsElement: the element is what tells one cast from another, and what is
// written reads back as the same cast.
func TestCastArrayKeepsItsElement(t *testing.T) {
	forms := map[string]string{}
	for written, stored := range castArrayForms {
		forms[written] = stored
	}
	for written, stored := range castArrayFormsTheServerRefuses {
		forms[written] = stored
	}
	for written, stored := range forms {
		t.Run(written, func(t *testing.T) {
			first := deparsedCast(t, "SELECT CAST(a AS "+written+")")
			if want := "cast(`a` as " + stored + ")"; first != want {
				t.Errorf("wrote %q, want %q", first, want)
			}
			if again := deparsedCast(t, "SELECT "+first); again != first {
				t.Errorf("%q reads back as %q", first, again)
			}
		})
	}
}

// TestOracle_CastArrayStoredForms holds castArrayForms to the server.
func TestOracle_CastArrayStoredForms(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if _, err := storeAndRead(ctx, o, "zz_cast_array", "CREATE TABLE zz (a JSON)", "zz"); err != nil {
		t.Fatal(err)
	}
	for written, stored := range castArrayForms {
		t.Run(written, func(t *testing.T) {
			if _, err := o.db.ExecContext(ctx, "CREATE OR REPLACE VIEW v AS SELECT CAST(a AS "+written+") x FROM zz"); err != nil {
				t.Fatal(err)
			}
			var name, ddl, charset, collation string
			if err := o.db.QueryRowContext(ctx, "SHOW CREATE VIEW v").Scan(&name, &ddl, &charset, &collation); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(ddl, " as "+stored+")") {
				t.Errorf("the server stored %s, recorded %q", ddl, stored)
			}
		})
	}
}

// TestCastKeepsADecimalsPrecision: NUMERIC, FIXED, DEC and NUMBER are DECIMAL under other names, and
// a cast to any of them keeps its precision and scale.
func TestCastKeepsADecimalsPrecision(t *testing.T) {
	for _, name := range []string{"DECIMAL", "DEC", "NUMERIC", "FIXED", "NUMBER"} {
		t.Run(name, func(t *testing.T) {
			first := deparsedCast(t, "SELECT CAST(a AS "+name+"(3,2))")
			if !strings.Contains(first, "(3,2)") {
				t.Errorf("wrote %q, which has lost the precision and scale", first)
			}
			if again := deparsedCast(t, "SELECT "+first); again != first {
				t.Errorf("%q reads back as %q", first, again)
			}
		})
	}
}

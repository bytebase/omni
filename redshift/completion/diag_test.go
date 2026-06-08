package completion

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/redshift/catalog"
	"github.com/bytebase/omni/redshift/parser"
)

func TestDiag(t *testing.T) {
	cases := []struct {
		name string
		sql  string
	}{
		{"ALTER TABLE DROP COLUMN", "ALTER TABLE t1 DROP COLUMN "},
		{"ALTER TABLE ALTER COLUMN", "ALTER TABLE t1 ALTER COLUMN "},
		{"ALTER TABLE ALTER COLUMN mid", "ALTER TABLE t1 ALTER COLUMN  SET NOT NULL"},
		{"COMMENT ON COLUMN t1.", "COMMENT ON COLUMN t1."},
		{"COMMENT ON COLUMN public.t1.", "COMMENT ON COLUMN public.t1."},
		{"Paren UNION", "(SELECT c1 FROM t1) UNION (SELECT  FROM t2)"},
	}

	cat := catalog.New()
	cat.Exec("CREATE TABLE t1 (c1 int); CREATE TABLE t2 (c1 int, c2 int);", nil)

	for _, tc := range cases {
		// Find cursor position (where | would be or trailing space)
		offset := len(tc.sql)
		if tc.name == "ALTER TABLE ALTER COLUMN mid" {
			offset = 33 // after "ALTER TABLE t1 ALTER COLUMN "
		}
		if tc.name == "Paren UNION" {
			offset = 34 // after "(SELECT c1 FROM t1) UNION (SELECT "
		}

		cs := parser.Collect(tc.sql, offset)
		hasColumnref := cs != nil && cs.HasRule("columnref")
		hasAnyName := cs != nil && cs.HasRule("any_name")

		candidates := Complete(tc.sql, offset, cat)
		var colCandidates []string
		for _, c := range candidates {
			if c.Type == CandidateColumn {
				colCandidates = append(colCandidates, c.Text)
			}
		}

		fmt.Printf("%-35s | columnref=%-5v any_name=%-5v | columns=%v\n",
			tc.name, hasColumnref, hasAnyName, colCandidates)
	}
}

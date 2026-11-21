package catalog

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// querySimulationFixture is the starting schema for TestQuerySimulation: a
// table named ident with a column named ident, so both table references and
// unqualified column references in generated compat.tsv queries and
// expressions resolve against a real table on TiDB.
const querySimulationFixture = "one_table@column"

// querySimulationPrefixes are the leading words of compat.tsv query statements
// wrapped directly as CREATE VIEW _v AS <stmt>.
var querySimulationPrefixes = []string{"SELECT ", "WITH ", "TABLE ", "VALUES "}

// TestQuerySimulation holds the catalog to testdata/queries/simulate.tsv:
// every generated query and DO expression statement TiDB accepts in the
// parser's compat.tsv, wrapped in CREATE VIEW _v AS ... on one_table@column.
// TiDB's stored view text (SHOW CREATE VIEW) is the engine's own reading of
// the query and expression.
//
// TiDB's side is on testdata/tapes/TestQuerySimulation.tape, so the test
// starts no server. OMNI_RECORD=1 records every row on the pinned TiDB;
// -update-compat commits the rows a catalog fix brought to TiDB's outcome; and
// -simulate-show=<text> prints both schemas of the rows that hold the text.
func TestQuerySimulation(t *testing.T) {
	rows, err := querySimulationRows(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "queries", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

func querySimulationRows(compatPath string) ([]simulate.Generated, error) {
	f, err := compat.ReadFile(compatPath)
	if err != nil {
		return nil, err
	}
	var rows []simulate.Generated
	seen := make(map[string]bool)
	for _, r := range f.Rows {
		if r.Expected != compat.Accept || handLabel(r.Rule) {
			continue
		}
		stmt, ok := wrapQueryView(r.Statement)
		if !ok || seen[stmt] {
			continue
		}
		seen[stmt] = true
		rows = append(rows, simulate.Generated{
			Fixture:   querySimulationFixture,
			Rule:      r.Rule,
			Statement: stmt,
		})
	}
	return rows, nil
}

func wrapQueryView(statement string) (string, bool) {
	if slices.ContainsFunc(querySimulationPrefixes, func(p string) bool { return strings.HasPrefix(statement, p) }) {
		body := strings.TrimSpace(strings.TrimRight(statement, " \t\r\n;"))
		return "CREATE VIEW _v AS " + body, true
	}
	if expr, ok := strings.CutPrefix(statement, "DO "); ok {
		body := strings.TrimSpace(strings.TrimRight(expr, " \t\r\n;"))
		return "CREATE VIEW _v AS SELECT " + body + " FROM ident", true
	}
	return "", false
}

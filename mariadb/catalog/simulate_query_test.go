package catalog

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mariadb/internal/testcontainer"
)

// TestQuerySimulation holds the catalog to testdata/queries/simulate.tsv:
// every generated query (SELECT, WITH, TABLE, VALUES) and expression (DO)
// statement MariaDB accepts in the parser's compat.tsv, wrapped in
// CREATE VIEW _v AS ... on one_table@column, with what MariaDB did and what
// Exec does.
//
// MariaDB's side is on the test's tape (testdata/tapes/TestQuerySimulation.tape),
// so the test starts no server. OMNI_RECORD=1 runs every row on the pinned
// MariaDB again; -update-compat commits the rows a catalog fix brought to
// MariaDB's outcome; and -simulate-show=<text> prints both schemas of the rows
// that hold the text.
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
		Engine:   simulationEngineTablesOnly,
	})
}

const querySimulationFixture = "one_table@column"

var querySimulationPrefixes = []string{
	"SELECT ",
	"WITH ",
	"TABLE ",
	"VALUES ",
}

func querySimulationRows(compatPath string) ([]simulate.Generated, error) {
	f, err := compat.ReadFile(compatPath)
	if err != nil {
		return nil, err
	}
	var rows []simulate.Generated
	seen := map[string]bool{}
	for _, r := range f.Rows {
		if r.Expected != compat.Accept || isHandCompatRule(r.Rule) {
			continue
		}
		wrapped, ok := wrapQuerySimulationStmt(r.Statement)
		if !ok || seen[wrapped] {
			continue
		}
		seen[wrapped] = true
		rows = append(rows, simulate.Generated{
			Fixture:   querySimulationFixture,
			Rule:      r.Rule,
			Statement: wrapped,
		})
	}
	return rows, nil
}

func isHandCompatRule(rule string) bool {
	i := strings.IndexAny(rule, ":/")
	return i >= 0 && rule[i] == ':'
}

func wrapQuerySimulationStmt(stmt string) (string, bool) {
	trimmed := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(stmt), ";"))
	if slices.ContainsFunc(querySimulationPrefixes, func(prefix string) bool {
		return strings.HasPrefix(stmt, prefix)
	}) {
		return "CREATE VIEW _v AS " + trimmed, true
	}
	if strings.HasPrefix(stmt, "DO ") {
		expr := strings.TrimSpace(strings.TrimPrefix(trimmed, "DO "))
		return "CREATE VIEW _v AS SELECT " + expr + " FROM a", true
	}
	return "", false
}

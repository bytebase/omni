package catalog

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// TestSimulationMySQL84 holds the catalog to testdata/mysql84/simulate.tsv:
// the same table, index and view DDL statements and fixtures as TestSimulation,
// plus MySQL 8.4 LTS-specific hand-added cases, recorded against pinned
// MySQL 8.4 LTS (mysql84 pin in mysql/testdata/pins.tsv).
//
// MySQL's side is on the test's tape (testdata/tapes/TestSimulationMySQL84.tape),
// so the test starts no server. OMNI_RECORD=1 runs every row on the pinned
// MySQL 8.4 LTS server; -update-compat commits the rows a catalog fix brought
// to MySQL's outcome; and -simulate-show=<text> prints both schemas of the rows
// that hold the text.
func TestSimulationMySQL84(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "mysql84", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.MySQL84,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     append(rows, familyRows...),
		Omni:     simulationOmniAt(MySQL84),
		Engine:   simulationEngineOn(testcontainer.MySQL84),
	})
}

// TestQuerySimulationMySQL84 holds the catalog at MySQL84 to
// testdata/mysql84/queries/simulate.tsv: TestQuerySimulation's rows, recorded
// against the pinned MySQL 8.4 LTS on testdata/tapes/TestQuerySimulationMySQL84.tape.
func TestQuerySimulationMySQL84(t *testing.T) {
	rows, err := querySimulationRows(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "mysql84", "queries", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.MySQL84,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		Omni:     simulationOmniAt(MySQL84),
		Engine:   simulationEngineOn(testcontainer.MySQL84),
	})
}

// TestProgramSimulationMySQL84 holds the catalog at MySQL84 to
// testdata/mysql84/programs/simulate.tsv: TestProgramSimulation's rows,
// recorded against the pinned MySQL 8.4 LTS on
// testdata/tapes/TestProgramSimulationMySQL84.tape.
func TestProgramSimulationMySQL84(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), programSimulationFixtures, programSimulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "mysql84", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.MySQL84,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     append(rows, eventRows...),
		Omni:     programSimulationOmniAt(MySQL84),
		Engine:   programSimulationEngineOn(testcontainer.MySQL84),
	})
}

package walktest

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/review/walk"
)

// Simulation holds an engine's simulate.tsv and how its WALK_THROUGH runs
// on each starting fixture. CheckSimulation prepares each fixture once and
// checks every row offline without starting a server.
type Simulation struct {
	// Path is the engine's simulate.tsv.
	Path string
	// Fixtures is the directory holding each fixture's <name>.sql.
	Fixtures string
	// Setup prepares one fixture from its SQL, returning the fixture's
	// starting object inventory and the walk runner for statements on it.
	// Called once per fixture named in Path.
	Setup func(t testing.TB, fx simulate.Fixture) SimFixture
	// Codes lists the error codes a finding's message stands for on stmt.
	Codes func(message, stmt string) []string
	// Same reports whether a key the walk claims names a stored object's
	// key, as the engine folds names.
	Same func(claimed, stored walk.Key) bool
}

// SimFixture is one prepared simulation fixture: its starting object
// inventory and the walk entry points for statements run against it.
type SimFixture struct {
	// Objects lists every object the fixture holds, as the engine's Lister
	// lists them for Judge.
	Objects []Object
	// Report reviews one statement on the fixture and returns its
	// WALK_THROUGH finding, or nil.
	Report func(t testing.TB, stmt string) *Finding
	// Violation returns the walk's violation for a statement Report
	// reported on.
	Violation func(t testing.TB, stmt string) *walk.Violation
}

// SimulationResult counts the rows CheckSimulation checked.
type SimulationResult struct {
	Total   int
	OK      int
	Refused int
	Found   int
}

// CheckSimulation holds an engine's WALK_THROUGH to its simulate.tsv
// (§3 of docs/specs/2026-11-11-engine-compatibility-without-servers-design.md):
//
//   - On an ok:<hash> row, any WALK_THROUGH finding fails.
//   - On a refused row, silence passes; a finding must name the engine's error
//     code (via Codes) or claim only facts true of the fixture's Objects (via
//     Claims and Same, matching Judge).
//   - An unjudged row (?) fails.
func CheckSimulation(t testing.TB, s Simulation) SimulationResult {
	t.Helper()
	if s.Setup == nil {
		t.Fatalf("walktest: Simulation.Setup is required")
	}
	if s.Same == nil {
		t.Fatalf("walktest: Simulation.Same is required")
	}
	f, err := simulate.ReadFile(s.Path)
	if err != nil {
		t.Fatalf("walktest: read %s: %v", s.Path, err)
	}
	prepared := make(map[string]SimFixture)
	var res SimulationResult
	for _, r := range f.Rows {
		res.Total++
		where := simRowWhere(s.Path, r)
		if r.Expected == simulate.Unjudged {
			t.Errorf("%s (%s): unjudged row (%s):\n    %s", where, r.Rule, r.Expected, r.Statement)
			continue
		}
		fx, ok := prepared[r.Fixture]
		if !ok {
			fxPath := filepath.Join(s.Fixtures, r.Fixture+".sql")
			b, err := os.ReadFile(fxPath)
			if err != nil {
				t.Fatalf("%s names fixture %s: %v", where, r.Fixture, err)
			}
			fx = s.Setup(t, simulate.Fixture{Name: r.Fixture, SQL: strings.TrimSpace(string(b))})
			if fx.Report == nil {
				t.Fatalf("walktest: fixture %s Setup returned nil Report", r.Fixture)
			}
			prepared[r.Fixture] = fx
		}
		if r.Accepted() {
			res.OK++
			if finding := fx.Report(t, r.Statement); finding != nil {
				t.Errorf("%s (%s): WALK_THROUGH reported %q on a statement the engine accepts (%s):\n    %s",
					where, r.Rule, finding.Message, r.Expected, r.Statement)
			}
			continue
		}
		res.Refused++
		finding := fx.Report(t, r.Statement)
		if finding == nil {
			continue
		}
		res.Found++
		if s.Codes != nil && slices.Contains(s.Codes(finding.Message, r.Statement), r.Expected) {
			continue
		}
		if fx.Violation == nil {
			t.Fatalf("walktest: fixture %s Setup returned nil Violation", r.Fixture)
		}
		v := fx.Violation(t, r.Statement)
		if v == nil {
			t.Errorf("%s (%s): reported %q on a statement the engine rejects with %s, with no violation to verify:\n    %s",
				where, r.Rule, finding.Message, r.Expected, r.Statement)
			continue
		}
		if holds, why := listed(Claims(v), fx.Objects, s.Same); !holds {
			t.Errorf("%s (%s): reported %q on a statement the engine rejects with %s, and the claim does not hold: %s\n    %s",
				where, r.Rule, finding.Message, r.Expected, why, r.Statement)
		}
	}
	pct := 0.0
	if res.Refused > 0 {
		pct = 100 * float64(res.Found) / float64(res.Refused)
	}
	t.Logf("%s: %d rows (%d ok, %d refused); WALK_THROUGH found %d/%d refusals (%.1f%%)",
		s.Path, res.Total, res.OK, res.Refused, res.Found, res.Refused, pct)
	return res
}

func simRowWhere(path string, r simulate.Row) string {
	if r.Line() > 0 {
		return fmt.Sprintf("%s:%d [%s]", path, r.Line(), r.Fixture)
	}
	return fmt.Sprintf("%s [%s]", path, r.Fixture)
}

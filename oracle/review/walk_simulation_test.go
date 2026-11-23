package review

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkSimFixtures are the starting schemas in testdata/fixtures/ for
// TestWalkSimulation. Compat.tsv statements write `a` for every sampled
// identifier (and `t` as a trigger's base table), so what a statement finds
// depends on what object(s) the fixture names `a`.
var walkSimFixtures = []string{
	"one_table@missing",
	"one_table@table",
	"one_table@column",
	"objects@partitioned",
	"objects@index_trigger",
	"objects@view",
	"objects@matview",
	"objects@sequence",
	"objects@synonym",
	"objects@function",
	"objects@type",
}

// walkSimStatement matches Oracle DDL and DML statements relevant to
// oracle/review's WALK_THROUGH. Materialized view logs are excluded because
// Oracle 23 hangs indefinitely on several ALTER MATERIALIZED VIEW LOG forms
// (documented in oracle/parser/testdata/README.md) and WALK_THROUGH does not
// lower them.
var walkSimStatement = regexp.MustCompile(`(?i)^(?:` +
	`CREATE\s+(?:OR\s+REPLACE\s+)?(?:EDITIONABLE\s+|NONEDITIONABLE\s+)?(?:` +
	`(?:GLOBAL\s+TEMPORARY\s+|PRIVATE\s+TEMPORARY\s+|SHARDED\s+|DUPLICATED\s+|IMMUTABLE\s+|BLOCKCHAIN\s+|JSON\s+COLLECTION\s+)*TABLE\b` +
	`|(?:UNIQUE\s+|BITMAP\s+|MULTIVALUE\s+|VECTOR\s+)*INDEX\b` +
	`|(?:NO\s*FORCE\s+|FORCE\s+)?(?:EDITIONING\s+|EDITIONABLE\s+|NONEDITIONABLE\s+)?(?:JSON\s+(?:RELATIONAL\s+)?DUALITY\s+)?VIEW\b` +
	`|MATERIALIZED\s+VIEW\b` +
	`|SEQUENCE\b` +
	`|(?:PUBLIC\s+)?SYNONYM\b` +
	`|TYPE\b` +
	`|PACKAGE\b` +
	`|PROCEDURE\b` +
	`|FUNCTION\b` +
	`|TRIGGER\b` +
	`)` +
	`|ALTER\s+(?:TABLE|INDEX|VIEW|MATERIALIZED\s+VIEW|SEQUENCE|SYNONYM|TYPE|PACKAGE|PROCEDURE|FUNCTION|TRIGGER)\b` +
	`|DROP\s+(?:TABLE|INDEX|VIEW|MATERIALIZED\s+VIEW|SEQUENCE|(?:PUBLIC\s+)?SYNONYM|TYPE|PACKAGE|PROCEDURE|FUNCTION|TRIGGER)\b` +
	`|TRUNCATE\s+TABLE\b` +
	`|RENAME\b` +
	`|COMMENT\s+ON\b` +
	`|(?:INSERT|UPDATE|DELETE|MERGE)\b` +
	`)`)

var walkSimMVLog = regexp.MustCompile(`(?i)^(?:CREATE|ALTER|DROP)\s+MATERIALIZED\s+VIEW\s+LOG\b`)

var walkSimDML = regexp.MustCompile(`(?i)^(?:INSERT|UPDATE|DELETE|MERGE)\b`)

func isWalkSimStatement(stmt string) bool {
	s := strings.TrimSpace(stmt)
	return walkSimStatement.MatchString(s) && !walkSimMVLog.MatchString(s)
}

var walkSimShowMissed = flag.Bool("walk-sim-missed", false, "log top missed refusal categories in TestWalkSimulation")

const walkSimMinTotalRecallPct = 85.5

var walkSimMinFixtureRecallPct = map[string]float64{
	"one_table@missing":     88.5,
	"one_table@table":       82.0,
	"one_table@column":      80.5,
	"objects@partitioned":   87.5,
	"objects@index_trigger": 91.5,
	"objects@view":          91.5,
	"objects@matview":       56.5,
	"objects@sequence":      93.0,
	"objects@synonym":       87.0,
	"objects@function":      93.5,
	"objects@type":          85.0,
}

const (
	walkSimVerdictsPath = "testdata/walk_simulation_verdicts.tsv"
	walkSimVerdictsHdr  = "# columns: expected, fixture, rule, statement"
	walkSimPinPrefix    = "# pin: "
)

type walkSimVerdictRow struct {
	Expected  string // "ok" or 5-digit Oracle ORA error code (e.g. "00942")
	Fixture   string
	Rule      string
	Statement string
	line      int
}

type walkSimVerdictFile struct {
	Pin  string
	Rows []walkSimVerdictRow
}

func readWalkSimVerdicts(path string) (*walkSimVerdictFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := &walkSimVerdictFile{}
	for i, raw := range strings.Split(string(b), "\n") {
		line := i + 1
		switch {
		case raw == "":
			continue
		case strings.HasPrefix(raw, "#"):
			if strings.HasPrefix(raw, walkSimPinPrefix) {
				f.Pin = strings.TrimPrefix(raw, walkSimPinPrefix)
			}
			continue
		}
		cols := strings.SplitN(raw, "\t", 4)
		if len(cols) != 4 {
			return nil, fmt.Errorf("%s:%d: want 4 tab-separated columns (expected, fixture, rule, statement), got %d", path, line, len(cols))
		}
		f.Rows = append(f.Rows, walkSimVerdictRow{
			Expected:  cols[0],
			Fixture:   cols[1],
			Rule:      cols[2],
			Statement: cols[3],
			line:      line,
		})
	}
	return f, nil
}

func writeWalkSimVerdicts(path string, f *walkSimVerdictFile) error {
	mode, err := enginetest.RunMode()
	if err != nil {
		return err
	}
	if mode == enginetest.ModeVerify {
		return fmt.Errorf("walk_simulation: cannot write %s with %s=verify: verification writes nothing", path, enginetest.RecordEnv)
	}
	var b strings.Builder
	b.WriteString(walkSimPinPrefix + f.Pin + "\n")
	b.WriteString(walkSimVerdictsHdr + "\n")
	for _, r := range f.Rows {
		if strings.Contains(r.Statement, "\n") {
			return fmt.Errorf("%s: statement holds a newline: %q", r.Rule, r.Statement)
		}
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", r.Expected, r.Fixture, r.Rule, r.Statement)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// TestWalkSimulation holds Oracle's WALK_THROUGH to every DDL and DML statement
// Oracle 23 accepts in oracle/parser/testdata/compat.tsv across representative
// schema fixtures synced from pinned Oracle 23.
//
//   - Offline replay (default, zero Docker): replays testdata/tapes/TestWalkSimulation.tape
//     and checks testdata/walk_simulation_verdicts.tsv.
//   - Recording (OMNI_RECORD=1): runs live against pinned Oracle 23 and writes
//     both the tape and testdata/walk_simulation_verdicts.tsv.
//   - Verification (OMNI_RECORD=verify): runs live against pinned Oracle 23,
//     verifies every verdict against both committed files, and writes nothing.
//
// Invariant:
//   - On every row Oracle 23 accepts ("ok"), WALK_THROUGH MUST report 0 findings.
//   - On every row Oracle 23 refuses, silence passes; any finding must carry a
//     valid code and either map to Oracle's error code or claim only facts that
//     hold in the fixture's real Oracle inventory.
func TestWalkSimulation(t *testing.T) {
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatal(err)
	}
	pin, err := enginetest.LoadPin(testcontainer.Pins, testcontainer.Default)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), walkSimFixtures, isWalkSimStatement)
	if err != nil {
		t.Fatal(err)
	}
	if len(generated) == 0 {
		t.Fatal("no statements selected from oracle/parser/testdata/compat.tsv")
	}

	recording := mode == enginetest.ModeRecord
	verifying := mode == enginetest.ModeVerify
	recordCmd := enginetest.RecordCommand(t)

	var committed *walkSimVerdictFile
	committedByKey := map[[3]string]string{}
	if !recording {
		committed, err = readWalkSimVerdicts(walkSimVerdictsPath)
		if err != nil {
			t.Fatalf("read %s: %v (record with: %s)", walkSimVerdictsPath, err, recordCmd)
		}
		if committed.Pin != pin.String() {
			t.Fatalf("%s was recorded at pin %q, but %s pins %q. Record it:\n  %s",
				walkSimVerdictsPath, committed.Pin, testcontainer.Pins, pin, recordCmd)
		}
		if len(committed.Rows) != len(generated) {
			t.Fatalf("%s has %d rows, but compat.tsv generates %d rows. Record it:\n  %s",
				walkSimVerdictsPath, len(committed.Rows), len(generated), recordCmd)
		}
		for i, g := range generated {
			r := committed.Rows[i]
			if r.Fixture != g.Fixture || r.Rule != g.Rule || r.Statement != g.Statement {
				t.Fatalf("%s:%d does not match generated row %d (%s / %s). Record it:\n  %s",
					walkSimVerdictsPath, r.line, i, g.Fixture, g.Rule, recordCmd)
			}
			committedByKey[[3]string{r.Fixture, r.Rule, r.Statement}] = r.Expected
		}
	}

	// Group generated rows by fixture in fixture order.
	rowsByFixture := make(map[string][]simulate.Generated, len(walkSimFixtures))
	var dmlTexts []string
	for _, g := range generated {
		rowsByFixture[g.Fixture] = append(rowsByFixture[g.Fixture], g)
		if walkSimDML.MatchString(strings.TrimSpace(g.Statement)) {
			dmlTexts = append(dmlTexts, g.Statement)
		}
	}
	// Only DML statements resolve public synonyms (spec §3-§4). Extracting
	// their target names keeps the Oracle string literal well under 4000 bytes.
	allNames := simPublicSynonymNames(dmlTexts)

	const shop, other = "SIMSHOP", "SIMOTHER"
	w := resetWalk(t, shop, other, nil)
	// Disable the recycle bin so DROP TABLE without PURGE never leaves
	// randomly named BIN$... objects in dba_tab_cols/dba_indexes across runs.
	ctx, cancel := walkContext()
	if _, err := w.adminConn.ExecContext(ctx, "ALTER SESSION SET RECYCLEBIN = OFF"); err != nil {
		cancel()
		t.Fatalf("disable recyclebin: %v", err)
	}
	cancel()
	installSimReset(t, w.adminConn)

	recorded := &walkSimVerdictFile{
		Pin:  pin.String(),
		Rows: make([]walkSimVerdictRow, 0, len(generated)),
	}

	var totalRows, totalOK, totalRefused, totalFound, totalFP int
	var mismatches []string
	var falsePositives []string
	var badClaims []string
	missedCounts := map[string]int{}
	missedSample := map[string]string{}

	for _, fxName := range walkSimFixtures {
		fxRows := rowsByFixture[fxName]
		fxPath := filepath.Join("testdata", "fixtures", fxName+".sql")
		rawSQL, err := os.ReadFile(fxPath)
		if err != nil {
			t.Fatalf("read fixture %s: %v", fxPath, err)
		}
		fxStmts := splitFixtureStatements(t, string(rawSQL))
		resetSimFixture(t, w.adminConn, shop, fxStmts, true)

		e := &oracleEngine{
			w:      w,
			conn:   w.adminConn,
			target: review.Target{Schema: walkSync(t, w.adminConn, shop, allNames)},
		}
		inventory := simNormalizeInventory(e.Inventory(t, allNames))

		var fxOK, fxRefused, fxFound int
		for _, g := range fxRows {
			totalRows++
			verdict := execSimStatement(t, w.adminConn, g.Statement)
			if verdict == "ok" {
				if walkSimDML.MatchString(strings.TrimSpace(g.Statement)) {
					ctx, cancel := walkContext()
					_, rerr := w.adminConn.ExecContext(ctx, "ROLLBACK")
					cancel()
					if rerr != nil {
						t.Fatalf("ROLLBACK after %q: %v", g.Statement, rerr)
					}
				} else {
					dropPub := strings.Contains(strings.ToUpper(g.Statement), "PUBLIC")
					resetSimFixture(t, w.adminConn, shop, fxStmts, dropPub)
				}
			}

			recorded.Rows = append(recorded.Rows, walkSimVerdictRow{
				Expected:  verdict,
				Fixture:   g.Fixture,
				Rule:      g.Rule,
				Statement: g.Statement,
			})

			if !recording {
				want := committedByKey[[3]string{g.Fixture, g.Rule, g.Statement}]
				if verdict != want {
					mismatches = append(mismatches, fmt.Sprintf("[%s] %s: got %s, committed %s\n    %s",
						g.Fixture, g.Rule, verdict, want, g.Statement))
				}
			}

			finding := walkSimReport(t, e, g.Statement)
			if verdict == "ok" {
				totalOK++
				fxOK++
				if finding != nil {
					totalFP++
					falsePositives = append(falsePositives, fmt.Sprintf("[%s] %s: WALK_THROUGH reported %q on Oracle-accepted statement:\n    %s",
						g.Fixture, g.Rule, finding.Message, g.Statement))
				}
				continue
			}

			totalRefused++
			fxRefused++
			if finding == nil {
				if *walkSimShowMissed {
					words := strings.Fields(g.Statement)
					prefix := ""
					if len(words) >= 2 {
						prefix = strings.ToUpper(words[0] + " " + words[1])
					} else if len(words) == 1 {
						prefix = strings.ToUpper(words[0])
					}
					mk := fmt.Sprintf("ORA-%s | %-22s | %s", verdict, prefix, g.Fixture)
					missedCounts[mk]++
					if _, ok := missedSample[mk]; !ok {
						missedSample[mk] = g.Statement
					}
				}
				continue
			}
			totalFound++
			fxFound++
			if slices.Contains(e.Codes(finding.Message), verdict) {
				continue
			}
			v := e.Violation(t, g.Statement)
			if holds, why := simClaimsHold(walktest.Claims(v), inventory, e.Same); !holds {
				badClaims = append(badClaims, fmt.Sprintf("[%s] %s: reported %q on ORA-%s, and claim does not hold (%s):\n    %s",
					g.Fixture, g.Rule, finding.Message, verdict, why, g.Statement))
			}
		}
		pct := 0.0
		if fxRefused > 0 {
			pct = 100 * float64(fxFound) / float64(fxRefused)
		}
		t.Logf("fixture %-22s: %4d rows (%4d ok, %4d refused), caught %4d/%4d (%.1f%%)",
			fxName, len(fxRows), fxOK, fxRefused, fxFound, fxRefused, pct)
		if minPct, ok := walkSimMinFixtureRecallPct[fxName]; ok && pct < minPct {
			t.Errorf("fixture %s recall %.1f%% (%d/%d) dropped below floor %.1f%%", fxName, pct, fxFound, fxRefused, minPct)
		}
	}

	if recording {
		if err := writeWalkSimVerdicts(walkSimVerdictsPath, recorded); err != nil {
			t.Fatalf("write %s: %v", walkSimVerdictsPath, err)
		}
		t.Logf("wrote %s (%d rows)", walkSimVerdictsPath, len(recorded.Rows))
	} else if verifying && len(mismatches) > 0 {
		t.Errorf("OMNI_RECORD=verify: %d verdicts differ from %s (showing up to 25):\n%s",
			len(mismatches), walkSimVerdictsPath, strings.Join(mismatches[:min(25, len(mismatches))], "\n"))
	} else if len(mismatches) > 0 {
		t.Errorf("tape and %s disagree on %d rows (showing up to 25):\n%s",
			walkSimVerdictsPath, len(mismatches), strings.Join(mismatches[:min(25, len(mismatches))], "\n"))
	}

	if len(falsePositives) > 0 {
		t.Errorf("FALSE POSITIVES: %d Oracle-accepted statements reported by WALK_THROUGH (showing up to 50):\n%s",
			len(falsePositives), strings.Join(falsePositives[:min(50, len(falsePositives))], "\n"))
	}
	if len(badClaims) > 0 {
		t.Errorf("INVALID CLAIMS: %d refused statements with unverified claims (showing up to 25):\n%s",
			len(badClaims), strings.Join(badClaims[:min(25, len(badClaims))], "\n"))
	}

	pct := 0.0
	if totalRefused > 0 {
		pct = 100 * float64(totalFound) / float64(totalRefused)
	}
	t.Logf("TOTAL: %d rows (%d ok, %d refused); false positives: %d; caught %d/%d refusals (%.1f%%)",
		totalRows, totalOK, totalRefused, totalFP, totalFound, totalRefused, pct)
	if pct < walkSimMinTotalRecallPct {
		t.Errorf("total WALK_THROUGH refusal recall %.1f%% (%d/%d) dropped below floor %.1f%%", pct, totalFound, totalRefused, walkSimMinTotalRecallPct)
	}
	t.Run("queries", func(t *testing.T) {
		walkSimSuite(t, w, shop, allNames, "queries", walkSimQueriesSuite, walkSimQueriesMinRecallPct)
	})
	t.Run("programs", func(t *testing.T) {
		walkSimSuite(t, w, shop, allNames, "programs", walkSimProgramsSuite, walkSimProgramsMinRecallPct)
	})
	if *walkSimShowMissed {
		type kv struct {
			k string
			n int
		}
		var list []kv
		for k, n := range missedCounts {
			list = append(list, kv{k, n})
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].n != list[j].n {
				return list[i].n > list[j].n
			}
			return list[i].k < list[j].k
		})
		for i, item := range list {
			if i >= 60 {
				break
			}
			t.Logf("MISSED %4d | %s | e.g. %s", item.n, item.k, missedSample[item.k])
		}
	}
}

// walkSimQueriesSuite is oracle/catalog's query suite: each query Oracle
// accepts in compat.tsv, and a corpus of Oracle's query forms, as CREATE
// VIEW v AS <query>, on the catalog's fixtures. The walk reads it for its
// statements and fixtures, and asks the server for each verdict itself.
const walkSimQueriesSuite = "../catalog/testdata/queries/simulate.tsv"

// walkSimQueriesMinRecallPct is the queries subtest's recall floor.
const walkSimQueriesMinRecallPct = 2.0

// walkSimProgramsSuite is oracle/catalog's stored-program suite: procedures,
// functions, packages, triggers and types with bodies, and changes to the
// objects they read, on programs@objects.
const walkSimProgramsSuite = "../catalog/testdata/programs/simulate.tsv"

// walkSimProgramsMinRecallPct is the programs subtest's recall floor.
const walkSimProgramsMinRecallPct = 53.0

// walkSimSuite holds WALK_THROUGH to one of oracle/catalog's suites as
// TestWalkSimulation holds it to compat.tsv: no finding on a statement
// Oracle accepts, and a finding on a refused one maps to Oracle's code or
// claims only what the fixture holds.
func walkSimSuite(t *testing.T, w *walkSchemas, shop, allNames, name, path string, floor float64) {
	suite, err := simulate.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []string
	rowsByFixture := map[string][]simulate.Row{}
	for _, r := range suite.Rows {
		if rowsByFixture[r.Fixture] == nil {
			fixtures = append(fixtures, r.Fixture)
		}
		rowsByFixture[r.Fixture] = append(rowsByFixture[r.Fixture], r)
	}
	var refused, found int
	var problems []string
	for _, fxName := range fixtures {
		fxPath := filepath.Join("..", "catalog", "testdata", "fixtures", fxName+".sql")
		rawSQL, err := os.ReadFile(fxPath)
		if err != nil {
			t.Fatalf("read fixture %s: %v", fxPath, err)
		}
		fxStmts := splitFixtureStatements(t, string(rawSQL))
		resetSimFixture(t, w.adminConn, shop, fxStmts, false)
		e := &oracleEngine{
			w:      w,
			conn:   w.adminConn,
			target: review.Target{Schema: walkSync(t, w.adminConn, shop, allNames)},
		}
		inventory := simNormalizeInventory(e.Inventory(t, allNames))
		var fxRefused, fxFound int
		for _, r := range rowsByFixture[fxName] {
			verdict := execSimStatement(t, w.adminConn, r.Statement)
			if verdict == "ok" {
				resetSimFixture(t, w.adminConn, shop, fxStmts, false)
			}
			finding := walkSimReport(t, e, r.Statement)
			switch {
			case verdict == "ok" && finding != nil:
				problems = append(problems, fmt.Sprintf("[%s] FALSE POSITIVE %q on an Oracle-accepted statement:\n    %s", fxName, finding.Message, r.Statement))
			case verdict == "ok":
			case finding == nil:
				fxRefused++
			default:
				fxRefused++
				fxFound++
				if slices.Contains(e.Codes(finding.Message), verdict) {
					continue
				}
				v := e.Violation(t, r.Statement)
				if holds, why := simClaimsHold(walktest.Claims(v), inventory, e.Same); !holds {
					problems = append(problems, fmt.Sprintf("[%s] INVALID CLAIM %q on ORA-%s (%s):\n    %s", fxName, finding.Message, verdict, why, r.Statement))
				}
			}
		}
		refused += fxRefused
		found += fxFound
		t.Logf("fixture %-22s: %4d rows, caught %4d/%4d refusals", fxName, len(rowsByFixture[fxName]), fxFound, fxRefused)
	}
	if len(problems) > 0 {
		t.Errorf("%d problems (showing up to 50):\n%s", len(problems), strings.Join(problems[:min(50, len(problems))], "\n"))
	}
	pct := 0.0
	if refused > 0 {
		pct = 100 * float64(found) / float64(refused)
	}
	t.Logf("%s: caught %d/%d refusals (%.1f%%)", name, found, refused, pct)
	if pct < floor {
		t.Errorf("%s recall %.1f%% (%d/%d) dropped below floor %.1f%%", name, pct, found, refused, floor)
	}
}

// simNormalizeInventory marks materialized view entries in spaceObject as
// MATERIALIZED VIEW instead of TABLE. In Oracle's DBA_OBJECTS, a materialized
// view's container table sits in namespace 1 (object_type='TABLE') while the
// materialized view sits in namespace 19, so walkInventory's namespace=1 query
// emits TABLE for spaceObject even though spaceTable emits MATERIALIZED VIEW.
func simNormalizeInventory(objects []walktest.Object) []walktest.Object {
	mviews := map[[2]string]bool{}
	for _, o := range objects {
		if o.Key.Space == spaceTable && o.Kind == kindMView {
			mviews[[2]string{o.Key.Parent, o.Key.Name}] = true
		}
	}
	out := make([]walktest.Object, len(objects))
	copy(out, objects)
	for i := range out {
		if out[i].Key.Space == spaceObject && mviews[[2]string{out[i].Key.Parent, out[i].Key.Name}] {
			out[i].Kind = kindMView
		}
	}
	return out
}

func splitFixtureStatements(t testing.TB, sqlText string) []string {
	t.Helper()
	e := &oracleEngine{}
	stmts := e.Statements(sqlText)
	out := make([]string, len(stmts))
	for i, s := range stmts {
		if text, ok := strings.CutSuffix(s, "\n/"); ok {
			out[i] = text
		} else {
			out[i] = strings.TrimSuffix(s, ";")
		}
	}
	return out
}

// installSimReset creates a helper procedure in SYSTEM (cleaned up by
// testcontainer between tests) that drops every object in a schema using a
// single indexed DBA_OBJECTS query. This avoids resending a multi-KB anonymous
// block on every reset and avoids scanning ~8000 PUBLIC synonyms unless a
// statement touched PUBLIC.
func installSimReset(t testing.TB, conn *sql.Conn) {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	proc := `CREATE OR REPLACE PROCEDURE system.omni_sim_reset(p_schema IN VARCHAR2, p_pub IN NUMBER) AUTHID CURRENT_USER IS
  rc SYS_REFCURSOR;
  v_type VARCHAR2(128);
  v_name VARCHAR2(128);
BEGIN
  ROLLBACK;
  IF p_pub <> 0 THEN
    OPEN rc FOR 'SELECT synonym_name FROM dba_synonyms WHERE owner = ''PUBLIC'' AND (table_owner = :1 OR synonym_name IN (''A'', ''IF'')) ORDER BY synonym_name' USING p_schema;
    LOOP
      FETCH rc INTO v_name;
      EXIT WHEN rc%NOTFOUND;
      BEGIN
        EXECUTE IMMEDIATE 'DROP PUBLIC SYNONYM "' || v_name || '"';
      EXCEPTION WHEN OTHERS THEN NULL;
      END;
    END LOOP;
    CLOSE rc;
  END IF;
  OPEN rc FOR 'SELECT object_type, object_name FROM dba_objects WHERE owner = :1 AND object_type IN (''MATERIALIZED VIEW'', ''TABLE'', ''VIEW'', ''SEQUENCE'', ''SYNONYM'', ''TRIGGER'', ''INDEX'', ''PROCEDURE'', ''FUNCTION'', ''PACKAGE'', ''TYPE'', ''PACKAGE BODY'', ''TYPE BODY'') ORDER BY CASE object_type WHEN ''MATERIALIZED VIEW'' THEN 1 WHEN ''TRIGGER'' THEN 2 WHEN ''TABLE'' THEN 3 WHEN ''VIEW'' THEN 4 WHEN ''INDEX'' THEN 5 WHEN ''SYNONYM'' THEN 6 WHEN ''SEQUENCE'' THEN 7 WHEN ''PROCEDURE'' THEN 8 WHEN ''FUNCTION'' THEN 9 WHEN ''PACKAGE BODY'' THEN 10 WHEN ''PACKAGE'' THEN 11 WHEN ''TYPE BODY'' THEN 12 WHEN ''TYPE'' THEN 13 END, object_name' USING p_schema;
  LOOP
    FETCH rc INTO v_type, v_name;
    EXIT WHEN rc%NOTFOUND;
    BEGIN
      IF v_type = 'MATERIALIZED VIEW' THEN
        BEGIN
          EXECUTE IMMEDIATE 'ALTER TABLE "' || p_schema || '"."' || v_name || '" ENABLE TABLE LOCK';
        EXCEPTION WHEN OTHERS THEN NULL;
        END;
        BEGIN
          EXECUTE IMMEDIATE 'ALTER TABLE "' || p_schema || '"."' || v_name || '" READ WRITE';
        EXCEPTION WHEN OTHERS THEN NULL;
        END;
        EXECUTE IMMEDIATE 'DROP MATERIALIZED VIEW "' || p_schema || '"."' || v_name || '"';
      ELSIF v_type = 'TABLE' THEN
        BEGIN
          EXECUTE IMMEDIATE 'ALTER TABLE "' || p_schema || '"."' || v_name || '" ENABLE TABLE LOCK';
        EXCEPTION WHEN OTHERS THEN NULL;
        END;
        EXECUTE IMMEDIATE 'DROP TABLE "' || p_schema || '"."' || v_name || '" CASCADE CONSTRAINTS PURGE';
      ELSIF v_type = 'TYPE' THEN
        EXECUTE IMMEDIATE 'DROP TYPE "' || p_schema || '"."' || v_name || '" FORCE';
      ELSIF v_type = 'VIEW' THEN
        EXECUTE IMMEDIATE 'DROP VIEW "' || p_schema || '"."' || v_name || '" CASCADE CONSTRAINTS';
      ELSIF v_type = 'SYNONYM' THEN
        EXECUTE IMMEDIATE 'DROP SYNONYM "' || p_schema || '"."' || v_name || '" FORCE';
      ELSE
        EXECUTE IMMEDIATE 'DROP ' || v_type || ' "' || p_schema || '"."' || v_name || '"';
      END IF;
    EXCEPTION WHEN OTHERS THEN NULL;
    END;
  END LOOP;
  CLOSE rc;
  BEGIN
    EXECUTE IMMEDIATE 'PURGE DBA_RECYCLEBIN';
  EXCEPTION WHEN OTHERS THEN NULL;
  END;
END;`
	if _, err := conn.ExecContext(ctx, proc); err != nil {
		t.Fatalf("installSimReset: %v", err)
	}
}

// resetSimFixture drops every object in schema (and any public synonym named A,
// IF, or targeting schema when dropPub is set), then rebuilds the fixture's
// statements.
func resetSimFixture(t testing.TB, conn *sql.Conn, schema string, fxStmts []string, dropPub bool) {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	pub := 0
	if dropPub {
		pub = 1
	}
	if _, err := conn.ExecContext(ctx, fmt.Sprintf("BEGIN system.omni_sim_reset('%s', %d); END;", schema, pub)); err != nil {
		t.Fatalf("resetSimFixture cleanup: %v", err)
	}
	for _, stmt := range fxStmts {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("resetSimFixture %s: %v", stmt, err)
		}
	}
}

func execSimStatement(t testing.TB, conn *sql.Conn, stmt string) string {
	t.Helper()
	ctx, cancel := walkContext()
	defer cancel()
	_, err := conn.ExecContext(ctx, stmt)
	if err == nil {
		return "ok"
	}
	m := oraCode.FindStringSubmatch(err.Error())
	if m == nil || errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("exec %q: %v", stmt, err)
	}
	return m[1]
}

// walkSimReport runs Review with WalkThrough on stmt. A Syntax finding (from a
// parser gap already tracked in compat.tsv) yields nil, matching pg/mysql.
func walkSimReport(t testing.TB, e *oracleEngine, stmt string) *walktest.Finding {
	t.Helper()
	if strings.TrimSpace(stmt) == "" {
		return nil
	}
	result, err := Review(context.Background(), stmt, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{e.target})
	if err != nil {
		t.Fatalf("Review: %v", err)
	}
	var out *walktest.Finding
	for _, f := range result.Findings {
		if err := review.CheckCodes(f); err != nil {
			t.Errorf("%v\n%s", err, stmt)
		}
		switch f.Rule {
		case review.Syntax:
			return nil
		case review.WalkThrough:
			if out != nil {
				t.Errorf("more than one finding: %q and %q\n%s", out.Message, f.Message, stmt)
			}
			out = &walktest.Finding{Statement: f.Statement, Message: f.Message}
			s := result.Statements[f.Statement]
			e.stmt, e.anchor = stmt[s.Start:s.End], stmt[f.Range.Start:f.Range.End]
		}
	}
	return out
}

// simClaimsHold checks that every fact a violation claims holds in objects,
// matching walktest.listed.
func simClaimsHold(claims []walktest.Claim, objects []walktest.Object, same func(claimed, stored walk.Key) bool) (bool, string) {
	for _, c := range claims {
		if c.Fact == walk.Unknown {
			continue
		}
		var kinds []string
		present := false
		for _, o := range objects {
			if same(c.Key, o.Key) {
				present = true
				kinds = append(kinds, o.Kind)
			}
		}
		switch {
		case present != (c.Fact == walk.Present):
			return false, fmt.Sprintf("%s %q in %q is present: %v", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), present)
		case present && c.Kind != "" && !slices.Contains(kinds, c.Kind) && !slices.Contains(kinds, ""):
			return false, fmt.Sprintf("%s %q in %q is a %s, not a %s", c.Key.Space, c.Key.Name, walk.Names(c.Key.Parent), strings.Join(kinds, " and a "), c.Kind)
		}
	}
	return true, ""
}

// simPublicSynonymNames returns a comma-wrapped sorted list of DML target names
// (plus A, T, DUAL) for filtering Oracle's ~8000 built-in PUBLIC synonyms in
// walkSync and walkInventory without exceeding Oracle's 4000-byte SQL literal
// limit.
func simPublicSynonymNames(dmlTexts []string) string {
	seen := map[string]bool{"A": true, "T": true, "DUAL": true}
	for _, text := range dmlTexts {
		_, stmts, syntax := parseReview(text)
		if syntax != nil {
			continue
		}
		for _, s := range stmts {
			targets, _ := dmlTargets(s.node)
			for _, tgt := range targets {
				if tgt != nil && tgt.Name != "" {
					seen[tgt.Name] = true
				}
			}
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	slices.Sort(names)
	return "," + strings.Join(names, ",") + ","
}

package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/mssql/internal/testcontainer"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkSimFixtures lists the starting schemas in testdata/fixtures/. A compat.tsv
// statement uses identifier `a` for every object, column, index, constraint,
// schema, and type name, so each fixture places `a` in a distinct catalog role
// to exercise every branch of WALK_THROUGH against real SQL Server 2022 state.
var walkSimFixtures = []string{
	"one_table@missing",
	"one_table@type",
	"one_table@table",
	"one_table@column",
	"one_table@index",
	"objects@view",
	"objects@sequence",
	"objects@procedure",
	"objects@function",
	"objects@synonym",
	"objects@schema",
	"objects@schema_nonempty",
	"objects@constraint",
	"objects@trigger",
	"objects@xml_table",
	"objects@spatial_table",
}

// walkSimStatement selects every database-local DDL, trigger-state, SELECT INTO,
// and DML statement accepted by SQL Server 2022 in mssql/parser/testdata/compat.tsv.
var walkSimStatement = regexp.MustCompile(`(?i)^(?:` +
	`CREATE\s+(?:OR\s+ALTER\s+)?(?:TABLE|VIEW|PROC|PROCEDURE|FUNCTION|TRIGGER|SEQUENCE|SYNONYM|SCHEMA|TYPE|DEFAULT|RULE)\b` +
	`|CREATE\s+(?:UNIQUE\s+)?(?:(?:NON)?CLUSTERED\s+)?(?:COLUMNSTORE\s+)?INDEX\b` +
	`|CREATE\s+(?:PRIMARY\s+|SELECTIVE\s+)?XML\s+INDEX\b` +
	`|CREATE\s+SPATIAL\s+INDEX\b` +
	`|CREATE\s+STATISTICS\b` +
	`|ALTER\s+(?:TABLE|VIEW|PROC|PROCEDURE|FUNCTION|TRIGGER|SEQUENCE|SCHEMA|INDEX)\b` +
	`|DROP\s+(?:TABLE|VIEW|PROC|PROCEDURE|FUNCTION|TRIGGER|SEQUENCE|SYNONYM|SCHEMA|TYPE|INDEX|STATISTICS|DEFAULT|RULE)\b` +
	`|TRUNCATE\s+TABLE\b` +
	`|(?:ENABLE|DISABLE)\s+TRIGGER\b` +
	`|SELECT\b.*\bINTO\b` +
	`|(?:INSERT|UPDATE|DELETE|MERGE)\b` +
	`)`)

type walkSimRow struct {
	expected  string // "ok" or SQL Server error number
	fixture   string
	rule      string
	statement string
}

const walkSimVerdictsPath = "testdata/walk_simulation_verdicts.tsv"

// TestWalkSimulation holds MSSQL's WALK_THROUGH to each statement suite
// SQL Server 2022 judged, each reported on its own: compat, the
// grammar-derived statements of mssql/parser/testdata/compat.tsv on this
// package's fixtures (walkCompatSimulation), and mssql/catalog's
// simulate.tsv suites, checked offline on the in-memory catalog
// (walkCatalogSimulation).
func TestWalkSimulation(t *testing.T) {
	db := startMSSQL(t)
	t.Run("compat", func(t *testing.T) { walkCompatSimulation(t, db) })
	for _, suite := range walkSimSuites {
		t.Run(suite.name, func(t *testing.T) {
			res := walkCatalogSimulation(t, suite.path, suite.database)
			pct := 0.0
			if res.Refused > 0 {
				pct = 100 * float64(res.Found) / float64(res.Refused)
			}
			if res.Found < suite.minFound || pct < suite.minPct {
				t.Errorf("%s recall regressed: got %d/%d (%.1f%%), want at least %d (%.1f%%)",
					suite.name, res.Found, res.Refused, pct, suite.minFound, suite.minPct)
			}
		})
	}
}

// walkSimSuites are mssql/catalog's simulate.tsv suites, with the database
// their recording ran in (testcontainer.Database of the catalog's test) and
// the recall each holds.
var walkSimSuites = []struct {
	name, path, database string
	minPct               float64
	minFound             int
}{
	{"ddl", filepath.Join("..", "catalog", "testdata", "simulate.tsv"), "omni_TestSimulation", 91.2, 2697},
	// The queries and programs suites hold COLLATE rows SQL Server refuses
	// for an unknown collation (Msg 448), which the walk's ledger, a record
	// of which objects exist, does not judge.
	{"queries", filepath.Join("..", "catalog", "testdata", "queries", "simulate.tsv"), "omni_TestQuerySimulation", 61.6, 278},
	{"programs", filepath.Join("..", "catalog", "testdata", "programs", "simulate.tsv"), "omni_TestProgramSimulation", 68.4, 50},
}

// walkCompatSimulation holds MSSQL's WALK_THROUGH to every unique
// grammar-derived DDL/DML statement SQL Server 2022 accepts in
// mssql/parser/testdata/compat.tsv, across all fixtures in
// testdata/fixtures/:
//
//   - Hard invariant: for every (fixture, statement) pair SQL Server 2022
//     accepts, WALK_THROUGH MUST report 0 findings.
//   - Soundness on refusals: whenever WALK_THROUGH reports a finding on a
//     statement SQL Server 2022 rejects, the finding's codes must pass
//     review.CheckCodes and either match SQL Server's error number via
//     walkCodes or make only claims true of the fixture's real SQL Server
//     catalog.
//
// SQL Server's answers are recorded in testdata/tapes/TestWalkSimulation.tape
// and mirrored in testdata/walk_simulation_verdicts.tsv so offline CI runs with
// zero Docker. Run with OMNI_RECORD=1 to re-record against the pinned SQL
// Server 2022 container, or OMNI_RECORD=verify to verify read-only against it.
func walkCompatSimulation(t *testing.T, db *sql.DB) {
	mode, err := enginetest.RunMode()
	if err != nil {
		t.Fatal(err)
	}
	pin, err := enginetest.LoadPin(testcontainer.Pins, testcontainer.Default)
	if err != nil {
		t.Fatal(err)
	}

	stmts := loadWalkSimStatements(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	const simDB = "sim"
	var recorded []walkSimRow
	var total, okCount, refusedCount, foundCount, syntaxSkipped, falsePositives, unsound int
	we := &walkEngine{}

	for _, fxName := range walkSimFixtures {
		fxPath := filepath.Join("testdata", "fixtures", fxName+".sql")
		raw, err := os.ReadFile(fxPath)
		if err != nil {
			t.Fatalf("read fixture %s: %v", fxPath, err)
		}
		resetSimFixture(t, conn, simDB, string(raw))
		schema := walkSync(t, conn, simDB)
		cat := loadSimCatalog(t, conn, simDB)
		target := review.Target{Schema: schema, DefaultSchema: "dbo"}

		for _, s := range stmts {
			total++
			verdict := execSimStatement(t, conn, s.statement)
			recorded = append(recorded, walkSimRow{
				expected:  verdict,
				fixture:   fxName,
				rule:      s.rule,
				statement: s.statement,
			})

			finding, isSyntax := reviewSimStatement(t, s.statement, target)
			if isSyntax {
				syntaxSkipped++
			}
			if verdict == "ok" {
				okCount++
				if finding != nil {
					falsePositives++
					t.Errorf("[%s] (%s): WALK_THROUGH false positive %q on a statement SQL Server 2022 accepts:\n    %s",
						fxName, s.rule, finding.Message, s.statement)
				}
				continue
			}

			refusedCount++
			if finding == nil {
				continue
			}
			foundCount++
			if slices.Contains(we.Codes(finding.Message), verdict) {
				continue
			}
			v := simViolation(t, s.statement, target)
			if v == nil {
				unsound++
				t.Errorf("[%s] (%s): reported %q on SQL Server error %s with no violation to verify:\n    %s",
					fxName, s.rule, finding.Message, verdict, s.statement)
				continue
			}
			if ok, why := cat.holds(v); !ok {
				unsound++
				t.Errorf("[%s] (%s): reported %q on SQL Server error %s, and claim does not hold (%s):\n    %s",
					fxName, s.rule, finding.Message, verdict, why, s.statement)
			}
		}
	}

	// The optional-clause rows (the (every site) rows of tokens.tsv and
	// suite:optional-clauses) added 226 refusals the walk finds 167 of,
	// which moved the floor from 88.5% to 88.0%.
	const minWalkSimRecall = 88.0
	pct := 0.0
	if refusedCount > 0 {
		pct = 100 * float64(foundCount) / float64(refusedCount)
	}
	t.Logf("compat: %d rows (%d fixtures x %d statements): %d ok, %d refused; WALK_THROUGH found %d/%d refusals (%.1f%%), %d false positives, %d unsound, %d syntax-skipped",
		total, len(walkSimFixtures), len(stmts), okCount, refusedCount, foundCount, refusedCount, pct, falsePositives, unsound, syntaxSkipped)
	if pct < minWalkSimRecall {
		t.Errorf("WALK_THROUGH recall %.1f%% (%d/%d) dropped below %.1f%% floor", pct, foundCount, refusedCount, minWalkSimRecall)
	}

	if mode == enginetest.ModeRecord {
		if err := writeWalkSimVerdicts(walkSimVerdictsPath, pin.String(), recorded); err != nil {
			t.Fatalf("write verdicts: %v", err)
		}
		return
	}
	verifyWalkSimVerdicts(t, walkSimVerdictsPath, pin.String(), recorded)
}

type walkSimStmt struct {
	rule      string
	statement string
}

func loadWalkSimStatements(t testing.TB) []walkSimStmt {
	t.Helper()
	compatPath := filepath.Join("..", "parser", "testdata", "compat.tsv")
	f, err := compat.ReadFile(compatPath)
	if err != nil {
		t.Fatalf("read %s: %v", compatPath, err)
	}
	seen := make(map[string]bool)
	var out []walkSimStmt
	for _, r := range f.Rows {
		if r.Expected != compat.Accept || seen[r.Statement] || !walkSimStatement.MatchString(r.Statement) {
			continue
		}
		// INSERT BULK is SQL Server's internal TDS bulk-load wire command:
		// when accepted it switches the connection into BCP streaming mode
		// expecting binary rows on the next packet rather than T-SQL.
		if strings.HasPrefix(strings.ToUpper(r.Statement), "INSERT BULK ") {
			continue
		}
		seen[r.Statement] = true
		out = append(out, walkSimStmt{rule: r.Rule, statement: r.Statement})
	}
	if len(out) == 0 {
		t.Fatalf("no statements selected from %s", compatPath)
	}
	return out
}

func resetSimFixture(t testing.TB, conn *sql.Conn, dbName, fxSQL string) {
	t.Helper()
	ctx := context.Background()
	setup := []string{
		"IF @@TRANCOUNT > 0 ROLLBACK",
		"USE master",
		fmt.Sprintf("IF DB_ID('%[1]s') IS NOT NULL BEGIN ALTER DATABASE %[1]s SET SINGLE_USER WITH ROLLBACK IMMEDIATE; DROP DATABASE %[1]s; END", dbName),
		"CREATE DATABASE " + dbName + " COLLATE " + serverCollation,
		"USE " + dbName,
	}
	for _, q := range setup {
		if _, err := conn.ExecContext(ctx, q); err != nil {
			t.Fatalf("reset fixture (%s): %v", q, err)
		}
	}
	for _, batch := range bytebaseBatches(fxSQL) {
		text := strings.TrimSpace(batch.text)
		if text == "" {
			continue
		}
		if _, err := conn.ExecContext(ctx, text); err != nil {
			t.Fatalf("fixture batch %q: %v", text, err)
		}
	}
}

func execSimStatement(t testing.TB, conn *sql.Conn, stmt string) string {
	t.Helper()
	ctx := context.Background()
	if _, err := conn.ExecContext(ctx, "BEGIN TRANSACTION"); err != nil {
		t.Fatalf("BEGIN TRANSACTION before %q: %v", stmt, err)
	}
	_, execErr := conn.ExecContext(ctx, stmt)
	if _, err := conn.ExecContext(ctx, "IF @@TRANCOUNT > 0 ROLLBACK"); err != nil {
		t.Fatalf("ROLLBACK after %q: %v", stmt, err)
	}
	if execErr == nil {
		return "ok"
	}
	var me mssql.Error
	if !errors.As(execErr, &me) {
		t.Fatalf("non-mssql error on %q: %v", stmt, execErr)
	}
	first := me
	if len(me.All) > 0 {
		first = me.All[0]
	}
	return strconv.Itoa(int(first.Number))
}

func reviewSimStatement(t testing.TB, stmt string, target review.Target) (*walktest.Finding, bool) {
	t.Helper()
	result, err := Review(context.Background(), stmt, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{target})
	if err != nil {
		t.Fatalf("Review(%q): %v", stmt, err)
	}
	for _, f := range result.Failures {
		t.Errorf("Review(%q) could not evaluate target %d: %v", stmt, f.Target, f.Err)
	}
	var out *walktest.Finding
	for _, f := range result.Findings {
		if err := review.CheckCodes(f); err != nil {
			t.Error(err)
		}
		switch f.Rule {
		case review.Syntax:
			return nil, true
		case review.WalkThrough:
			if out != nil {
				t.Errorf("more than one finding on %q: %q and %q", stmt, out.Message, f.Message)
			}
			out = &walktest.Finding{Statement: f.Statement, Message: f.Message}
		}
	}
	return out, false
}

func simViolation(t testing.TB, stmt string, target review.Target) *walk.Violation {
	t.Helper()
	_, stmts, syntax := parseReview(stmt)
	if syntax != nil {
		return nil
	}
	v, err := walkThrough(context.Background(), stmt, stmts, target)
	if err != nil {
		t.Fatalf("walkThrough(%q): %v", stmt, err)
	}
	return v
}

// simCatalog caches the real SQL Server catalog of one fixture so claimHolds
// can check violations without issuing conditional queries that would depend
// on omni's output (which would break tape replay determinism).
type simCatalog struct {
	schemas     []catObject
	schemaByCls map[string]string // schema class -> schema_id ("*" if ambiguous)
	objects     map[string][]catObject
	relations   map[string][]catObject
	types       map[string][]catObject
	tableByPath map[string]string // Path(schemaClass, tableClass) -> oid ("*" if ambiguous)
	columns     map[string][]catObject
	indexes     map[string][]catObject
	constraints map[string][]catObject
	primaries   map[string][]catObject
}

func loadSimCatalog(t testing.TB, conn *sql.Conn, dbName string) *simCatalog {
	t.Helper()
	ctx := context.Background()
	c := &simCatalog{
		schemaByCls: make(map[string]string),
		objects:     make(map[string][]catObject),
		relations:   make(map[string][]catObject),
		types:       make(map[string][]catObject),
		tableByPath: make(map[string]string),
		columns:     make(map[string][]catObject),
		indexes:     make(map[string][]catObject),
		constraints: make(map[string][]catObject),
		primaries:   make(map[string][]catObject),
	}
	query := func(q string, scan func(*sql.Rows) error) {
		t.Helper()
		q = strings.ReplaceAll(q, "sim.", dbName+".")
		rows, err := conn.QueryContext(ctx, q)
		if err != nil {
			t.Fatalf("sim catalog %s: %v", q, err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := scan(rows); err != nil {
				t.Fatal(err)
			}
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}

	schemaNameToClass := make(map[string]string)
	query("SELECT name FROM sim.sys.schemas ORDER BY name", func(r *sql.Rows) error {
		var name string
		if err := r.Scan(&name); err != nil {
			return err
		}
		c.schemas = append(c.schemas, catObject{name: name})
		if n, ok := keyed(name); ok {
			if prev, exists := c.schemaByCls[n.class]; exists && prev != name {
				c.schemaByCls[n.class] = "*"
			} else {
				c.schemaByCls[n.class] = name
			}
			schemaNameToClass[name] = n.class
		}
		return nil
	})
	query("SELECT s.name, o.name, o.type FROM sim.sys.objects o JOIN sim.sys.schemas s ON s.schema_id = o.schema_id ORDER BY s.name, o.name, o.type", func(r *sql.Rows) error {
		var sname, name, typ string
		if err := r.Scan(&sname, &name, &typ); err != nil {
			return err
		}
		kind := objectKinds[strings.TrimSpace(typ)]
		if kind == "" {
			kind = "other"
		}
		c.objects[sname] = append(c.objects[sname], catObject{name: name, kind: kind})
		if isRelation(kind) {
			c.relations[sname] = append(c.relations[sname], catObject{name: name, kind: kind})
		}
		return nil
	})
	query("SELECT s.name, ty.name FROM sim.sys.types ty JOIN sim.sys.schemas s ON s.schema_id = ty.schema_id WHERE ty.is_user_defined = 1 ORDER BY s.name, ty.name", func(r *sql.Rows) error {
		var sname, name string
		if err := r.Scan(&sname, &name); err != nil {
			return err
		}
		c.types[sname] = append(c.types[sname], catObject{name: name})
		return nil
	})
	query("SELECT s.name, t.name FROM sim.sys.tables t JOIN sim.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name", func(r *sql.Rows) error {
		var sname, name string
		if err := r.Scan(&sname, &name); err != nil {
			return err
		}
		sc := schemaNameToClass[sname]
		if n, ok := keyed(name); ok && sc != "" {
			p := walk.Path(sc, n.class)
			tkey := sname + "\x00" + name
			if prev, exists := c.tableByPath[p]; exists && prev != tkey {
				c.tableByPath[p] = "*"
			} else {
				c.tableByPath[p] = tkey
			}
		}
		return nil
	})
	query("SELECT s.name, t.name, c.name FROM sim.sys.columns c JOIN sim.sys.tables t ON t.object_id = c.object_id JOIN sim.sys.schemas s ON s.schema_id = t.schema_id ORDER BY s.name, t.name, c.name", func(r *sql.Rows) error {
		var sname, tname, name string
		if err := r.Scan(&sname, &tname, &name); err != nil {
			return err
		}
		tkey := sname + "\x00" + tname
		c.columns[tkey] = append(c.columns[tkey], catObject{name: name})
		return nil
	})
	query("SELECT s.name, t.name, i.name, CASE WHEN i.is_primary_key = 1 THEN 'pk' WHEN i.is_unique_constraint = 1 THEN 'uq' ELSE 'index' END FROM sim.sys.indexes i JOIN sim.sys.tables t ON t.object_id = i.object_id JOIN sim.sys.schemas s ON s.schema_id = t.schema_id WHERE i.name IS NOT NULL ORDER BY s.name, t.name, i.name", func(r *sql.Rows) error {
		var sname, tname, name, kind string
		if err := r.Scan(&sname, &tname, &name, &kind); err != nil {
			return err
		}
		tkey := sname + "\x00" + tname
		c.indexes[tkey] = append(c.indexes[tkey], catObject{name: name, kind: kind})
		return nil
	})
	query("SELECT s.name, t.name, o.name, o.type FROM sim.sys.objects o JOIN sim.sys.tables t ON t.object_id = o.parent_object_id JOIN sim.sys.schemas s ON s.schema_id = t.schema_id WHERE o.type IN ('PK', 'UQ', 'F', 'C', 'D', 'EC') ORDER BY s.name, t.name, o.name", func(r *sql.Rows) error {
		var sname, tname, name, typ string
		if err := r.Scan(&sname, &tname, &name, &typ); err != nil {
			return err
		}
		tkey := sname + "\x00" + tname
		c.constraints[tkey] = append(c.constraints[tkey], catObject{name: name, kind: objectKinds[strings.TrimSpace(typ)]})
		return nil
	})
	query("SELECT s.name, t.name FROM sim.sys.key_constraints k JOIN sim.sys.tables t ON t.object_id = k.parent_object_id JOIN sim.sys.schemas s ON s.schema_id = t.schema_id WHERE k.type = 'PK' ORDER BY s.name, t.name", func(r *sql.Rows) error {
		var sname, tname string
		if err := r.Scan(&sname, &tname); err != nil {
			return err
		}
		tkey := sname + "\x00" + tname
		c.primaries[tkey] = append(c.primaries[tkey], catObject{name: primaryName})
		return nil
	})
	return c
}

func (c *simCatalog) lookup(k walk.Key) ([]catObject, bool) {
	path := walk.Names(k.Parent)
	switch k.Space {
	case spaceSchema:
		return c.schemas, true
	case spaceTemp:
		return nil, true
	}
	if len(path) == 0 {
		return nil, false
	}
	sid, ok := c.schemaByCls[path[0]]
	if !ok {
		return nil, true
	}
	if sid == "*" {
		return nil, false
	}
	switch k.Space {
	case spaceObject:
		return c.objects[sid], true
	case spaceRelation:
		return c.relations[sid], true
	case spaceType:
		return c.types[sid], true
	}
	if len(path) != 2 {
		return nil, false
	}
	oid, ok := c.tableByPath[k.Parent]
	if !ok {
		return nil, true
	}
	if oid == "*" {
		return nil, false
	}
	switch k.Space {
	case spaceColumn:
		return c.columns[oid], true
	case spaceIndex:
		return c.indexes[oid], true
	case spaceConstraint:
		return c.constraints[oid], true
	case spacePrimary:
		return c.primaries[oid], true
	}
	return nil, false
}

func (c *simCatalog) holds(v *walk.Violation) (bool, string) {
	for _, cl := range walktest.Claims(v) {
		if cl.Fact == walk.Unknown {
			continue
		}
		objs, ok := c.lookup(cl.Key)
		if !ok {
			return false, fmt.Sprintf("cannot read parent of %+v", cl.Key)
		}
		if ok, why := claimHolds(cl, objs, strings.EqualFold); !ok {
			return false, why
		}
	}
	return true, ""
}

func formatWalkSimVerdicts(pin string, rows []walkSimRow) string {
	var b strings.Builder
	b.WriteString("# pin: " + pin + "\n")
	b.WriteString("# columns: expected, fixture, rule, statement\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", r.expected, r.fixture, r.rule, r.statement)
	}
	return b.String()
}

func writeWalkSimVerdicts(path, pin string, rows []walkSimRow) error {
	mode, err := enginetest.RunMode()
	if err != nil {
		return err
	}
	if mode == enginetest.ModeVerify {
		return fmt.Errorf("walk simulation: cannot write %s with %s=verify", path, enginetest.RecordEnv)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(formatWalkSimVerdicts(pin, rows)), 0o644)
}

func verifyWalkSimVerdicts(t testing.TB, path, pin string, got []walkSimRow) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (record with: %s)", path, err, enginetest.RecordCommand(t))
	}
	want := formatWalkSimVerdicts(pin, got)
	if string(b) != want {
		t.Errorf("%s differs from current SQL Server / tape verdicts; re-record with:\n  %s", path, enginetest.RecordCommand(t))
	}
}

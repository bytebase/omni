package catalog

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"maps"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// The OceanBase oracle for this package: every TestOracle_* asks it.
//
// oceanBaseOracle opens the test's connection to OceanBase: a replay of its tape, or the
// pinned server in a recording (testcontainer.OpenWith). The pool holds one connection,
// so a database selected with USE stays selected, and nothing needs db.Conn to pin one:
// a helper that asks the pool while its caller holds that connection waits forever.
//
// The version is MySQL80, the only one, and the probes measure that it fits. 4.3.5 does
// report a 5.7 version string, and does print 5.7-style widths, but the stored form it
// prints is the one MySQL80's normalizer reproduces: 8.0 strips widths from both sides.
// A wrong normalizer shows up as a non-empty diff rather than as a phantom change in a
// user's release.
func oceanBaseOracle(t testing.TB) *oracleConn {
	t.Helper()
	o := &oracleConn{
		db:      testcontainer.OpenWith(t, testcontainer.Options{MultiStatements: true}),
		version: MySQL80,
		name:    "oceanbase",
	}
	top := t.Name()
	t.Cleanup(func() {
		// A run of chosen subtests leaves the others' entries unvisited.
		if t.Failed() || strings.Contains(testRunFlag(), "/") {
			return
		}
		for _, name := range slices.Sorted(maps.Keys(obRefused)) {
			if _, seen := obRefusedSeen.Load(name); !seen && strings.HasPrefix(name, top+"/") {
				t.Errorf("obRefused lists %s (%s), but the probe did not end in that refusal: remove the entry", name, obRefused[name].what)
			}
		}
		for _, name := range slices.Sorted(maps.Keys(obGaps)) {
			if _, seen := obGapsSeen.Load(name); !seen && strings.HasPrefix(name, top+"/") {
				t.Errorf("obGaps lists %s, but no such probe ran under knownGap: remove the entry", name)
			}
		}
	})
	return o
}

func testRunFlag() string {
	if f := flag.Lookup("test.run"); f != nil {
		return f.Value.String()
	}
	return ""
}

// A refusal is how OceanBase answers a form it will not take.
type refusal struct {
	errno uint16
	what  string // the construct, for the log
}

// obRefused lists the probes whose form OceanBase refuses, by subtest name. The probes
// were written for MySQL; where OceanBase will not take one, the engine stores nothing
// to compare and the probe ends there. That is a form the probe sets up, or a statement
// of a plan that the plan itself warns of (applyPlan). refused holds every entry to the
// engine in both directions: a refusal that is not listed fails its test, and so does a
// listed probe that no longer ends in its refusal.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6.
var obRefused = map[string]refusal{
	"TestOracle_DiffIdempotence/oceanbase/binary-modifier-after-charset":               {1064, "BINARY after CHARACTER SET in a column's type"},
	"TestOracle_DiffIdempotence/oceanbase/adjacent-literal-default-double-quoted":      {1064, `a double-quoted literal with a neighbour, "x" 'y', in DEFAULT`},
	"TestOracle_MigrationIdempotence/oceanbase/binary-modifier-after-charset":          {1064, "BINARY after CHARACTER SET in a column's type"},
	"TestOracle_MigrationIdempotence/oceanbase/adjacent-literal-default-double-quoted": {1064, `a double-quoted literal with a neighbour, "x" 'y', in DEFAULT`},
	"TestOracle_PartitionChangeDetected/oceanbase/toggle-linear":                       {1064, "LINEAR KEY partitioning"},

	"TestOracle_PartitionDiffIdempotence/oceanbase/key-algorithm-1":      {1064, "KEY ALGORITHM = n partitioning"},
	"TestOracle_PartitionDiffIdempotence/oceanbase/key-algorithm-2":      {1064, "KEY ALGORITHM = n partitioning"},
	"TestOracle_PartitionDiffIdempotence/oceanbase/linear-hash":          {1064, "LINEAR HASH partitioning"},
	"TestOracle_PartitionDiffIdempotence/oceanbase/linear-key":           {1064, "LINEAR KEY partitioning"},
	"TestOracle_PartitionMigrationIdempotence/oceanbase/key-algorithm-1": {1064, "KEY ALGORITHM = n partitioning"},
	"TestOracle_PartitionMigrationIdempotence/oceanbase/key-algorithm-2": {1064, "KEY ALGORITHM = n partitioning"},
	"TestOracle_PartitionMigrationIdempotence/oceanbase/linear-hash":     {1064, "LINEAR HASH partitioning"},
	"TestOracle_PartitionMigrationIdempotence/oceanbase/linear-key":      {1064, "LINEAR KEY partitioning"},

	"TestOracle_PhantomDiffElimination/oceanbase/char-binary-attribute-ascii-explicit": {1064, "BINARY after CHARACTER SET in a column's type"},
	"TestOracle_PhantomDiffElimination/oceanbase/char-binary-attribute-utf8-alias":     {1064, "BINARY after CHARACTER SET in a column's type"},
	"TestOracle_PhantomDiffElimination/oceanbase/functional-default-json":              {1101, "a default on a JSON column"},

	"TestOracle_PartitionApplyCorrectness/oceanbase/create-key-algorithm-2":    {1064, "KEY ALGORITHM = n partitioning"},
	"TestOracle_PartitionApplyCorrectness/oceanbase/create-linear-key":         {1064, "LINEAR KEY partitioning"},
	"TestOracle_PartitionApplyCorrectness/oceanbase/repartition-toggle-linear": {1064, "LINEAR KEY partitioning"},

	// Changes OceanBase has no statement for. The plan keeps the statement MySQL would run
	// and warns of it; applyPlan checks that OceanBase refuses the warned statement.
	"TestOracle_PartitionApplyCorrectness/oceanbase/remove-partitioning":                 {1064, "REMOVE PARTITIONING"},
	"TestOracle_PartitionApplyCorrectness/oceanbase/remove-partitioning-range":           {1064, "REMOVE PARTITIONING"},
	"TestOracle_PartitionApplyCorrectness/oceanbase/remove-partitioning-drop-key-column": {1064, "REMOVE PARTITIONING"},
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/add-srid-none-to-4326":            {1235, "a change of a spatial column's SRID"},
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/add-srid-none-to-zero":            {1235, "a change of a spatial column's SRID"},
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/change-srid-4326-to-zero":         {1235, "a change of a spatial column's SRID"},
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/change-srid-zero-to-4326":         {1235, "a change of a spatial column's SRID"},
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/remove-srid-4326-to-none":         {1235, "a change of a spatial column's SRID"},
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/remove-srid-zero-to-none":         {1235, "a change of a spatial column's SRID"},
	"TestOracle_MigrationApplyCorrectness/oceanbase/modify-generated-expr-under-check":   {1235, "a change of a generated column's expression"},
	"TestOracle_IndexMigrationApplyCorrectness/oceanbase/add-mv-index":                   {1235, "a multi-valued index added to a table that exists"},

	"TestOracle_EventIdempotence/oceanbase/every-disable-on-slave":            {1064, "DISABLE ON SLAVE on an event"},
	"TestOracle_EventApplyCorrectness/oceanbase/alter-status-slave":           {1064, "DISABLE ON SLAVE on an event"},
	"TestOracle_EventIdempotence/oceanbase/every-interval-quoted-hm":          {1064, "a quoted interval, EVERY '1:30' HOUR_MINUTE"},
	"TestOracle_EventIdempotence/oceanbase/every-interval-quoted-leadingzero": {1064, "a quoted interval, EVERY '02:03:04' HOUR_SECOND"},
	"TestOracle_EventIdempotence/oceanbase/every-interval-quoted-daysecond":   {1064, "a quoted interval, EVERY '1 12:30:00' DAY_SECOND"},
	"TestOracle_EventIdempotence/oceanbase/every-interval-single-quoted":      {1064, "a quoted interval, EVERY '05' MINUTE"},
}

var obRefusedSeen sync.Map

// refused ends a probe whose statement the engine would not take. err is the engine's
// answer, or nil when the helper that gave up kept it in o.refusal. The probe passes only
// if obRefused records this subtest with the error number the engine answered; anything
// else fails. The caller returns after it.
func (o *oracleConn) refused(t testing.TB, err error, format string, args ...any) {
	t.Helper()
	if err == nil {
		err = o.refusal
	}
	o.refusal = nil
	why := fmt.Sprintf(format, args...)
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		t.Fatalf("%s, and not because OceanBase refused a statement: %v", why, err)
	}
	want, ok := obRefused[t.Name()]
	if !ok {
		t.Fatalf("%s: OceanBase answers %v. obRefused does not record that for %s: name the construct there, or fix the probe", why, err, t.Name())
	}
	if want.errno != myErr.Number {
		t.Fatalf("%s: OceanBase answers %v, but obRefused records %d (%s)", why, err, want.errno, want.what)
	}
	obRefusedSeen.Store(t.Name(), true)
	t.Logf("OceanBase does not take %s (%d), so it stores no form to compare", want.what, myErr.Number)
}

// applyPlan runs a plan's statements on the engine, one at a time, and reports whether the
// plan ran to its end. A statement the plan does not warn about has to be taken. One it
// warns about has to be refused, with the error number obRefused records for the probe:
// the Warning says OceanBase has no statement for the change, and this asks the engine
// whether that still holds. The plan stops there, with nothing for the caller to compare.
func (o *oracleConn) applyPlan(t testing.TB, conn oracleSession, plan *MigrationPlan) bool {
	t.Helper()
	for _, op := range plan.Ops {
		_, err := conn.ExecContext(context.Background(), op.SQL)
		switch {
		case op.Warning != "" && err == nil:
			t.Fatalf("[%s] OceanBase took a statement the plan warns that it refuses:\n  stmt: %s\n  warning: %s", o.name, op.SQL, op.Warning)
		case op.Warning != "":
			o.refused(t, err, "the plan warns of %q: %s", op.SQL, op.Warning)
			return false
		case err != nil:
			t.Fatalf("[%s] APPLY FAILED:\n  stmt: %s\n  err: %v\n  full plan:\n%s", o.name, op.SQL, err, plan.SQL())
		}
	}
	return true
}

// A gap is a probe this catalog fails on OceanBase today.
type gap struct {
	why   string   // what the catalog gets wrong
	fails []string // a phrase from each failure the probe reports, in order
}

// obGaps lists the probes that fail on OceanBase today, by subtest name: bugs and missing
// pieces of this catalog and its parser that are not fixed yet, and two readbacks of
// OceanBase's that are not the text it was given. The engine's side of each is on the
// probe's tape. knownGap runs a listed probe and passes it only while it fails exactly
// as recorded, so an entry cannot hide anything else and cannot outlive its gap: a fix
// fails the probe until the entry is removed.
//
// This is the work list. An entry is not a way to make a new failure pass.
var obGaps = map[string]gap{
	"TestOracle_CheckDiffIdempotence/oceanbase/between": {
		"OceanBase stores BETWEEN over constant bounds as two comparisons, ((`a` >= 1) and (`a` <= 10)), in a CHECK, a generated column and a view; the catalog keeps BETWEEN, so the CHECK differs from its own readback",
		twoWays("CANONICALIZATION: user form vs stored form not empty", "CANONICALIZATION (reverse): stored vs user not empty"),
	},
	"TestOracle_PhantomDiffElimination/oceanbase/column-charset-echo-80": {
		"OceanBase's default collation for utf8mb4 is utf8mb4_general_ci, which a readback leaves out; the catalog takes the default to be MySQL 8.0's utf8mb4_0900_ai_ci, so a column written COLLATE utf8mb4_general_ci differs from its own readback",
		[]string{`PHANTOM DIFF on column "a"`},
	},

	// SHOW CREATE FUNCTION prints a COMMENT without its escapes, so the readback is not
	// the comment the routine has. information_schema.ROUTINES is where to read it.
	"TestOracle_RoutineDiffIdempotence/oceanbase/fn-comment-quote": {
		"OceanBase's SHOW CREATE FUNCTION prints a COMMENT's quote unescaped, COMMENT 'it's here', which does not parse",
		[]string{"reload of routine readback failed"},
	},
	"TestOracle_RoutineDiffIdempotence/oceanbase/fn-comment-backslash": {
		`OceanBase's SHOW CREATE FUNCTION prints a COMMENT's backslash unescaped, COMMENT 'a\b' for the comment a\b, which reads back as another comment`,
		twoWays("CANONICALIZATION: user vs stored not empty", "CANONICALIZATION (reverse): stored vs user not empty"),
	},

	// What OceanBase rewrites in a view's body. The canonical body does not follow yet,
	// so each of these views differs from its own readback.
	"TestOracle_ViewIdempotence/oceanbase/cte":                              viewBody("OceanBase prints WITH RECURSIVE for a CTE written without RECURSIVE"),
	"TestOracle_ViewIdempotence/oceanbase/explicit-columns":                 viewBody("OceanBase folds a view's column list into the aliases of its select list: CREATE VIEW v (p, q) AS SELECT a, b reads back as select a AS p, b AS q"),
	"TestOracle_ViewIdempotence/oceanbase/view-on-view-explicit-cols":       viewBody("OceanBase folds a view's column list into the aliases of its select list"),
	"TestOracle_ViewIdempotence/oceanbase/paren-subquery-compare-if":        viewBody("OceanBase stores IF(c, a, b) as (case when c then a else b end)"),
	"TestOracle_ViewIdempotence/oceanbase/paren-subquery-compare-if-depth2": viewBody("OceanBase stores IF(c, a, b) as (case when c then a else b end)"),
	"TestOracle_ViewIdempotence/oceanbase/paren-subquery-in-depth2":         viewBody("OceanBase stores x IN ((subquery)) as x = any (subquery)"),
	"TestOracle_ViewIdempotence/oceanbase/paren-subquery-setop-compare":     viewBody("OceanBase prints each branch of a UNION in parentheses"),
	"TestOracle_ViewIdempotence/oceanbase/union":                            viewBody("OceanBase prints each branch of a UNION in parentheses"),
	"TestOracle_ViewIdempotence/oceanbase/trim-keyword-forms":               viewBody("OceanBase stores TRIM('x' FROM s) as trim(both 'x' from s)"),
	"TestOracle_ViewIdempotence/oceanbase/weight-string-forms": {
		"OceanBase prints WEIGHT_STRING as its internal five-argument call, weight_string(`s`,0,4,64,0), which the parser does not read",
		[]string{"load stored SDL failed"},
	},

	// Events. OceanBase keeps an event's comment when ALTER EVENT gives an empty one, and
	// gives an altered event the server's sql_mode, not the session's.
	"TestOracle_EventApplyCorrectness/oceanbase/alter-remove-comment": {
		"OceanBase's ALTER EVENT ... COMMENT '' leaves the comment as it was, so the plan cannot take one away",
		twoWays("APPLY-CORRECTNESS FAILED for alter-remove-comment: result != to", "APPLY-CORRECTNESS (reverse) FAILED for alter-remove-comment"),
	},
	"TestOracle_EventAlterPreservesSQLMode/oceanbase": {
		"OceanBase gives an event altered under SET sql_mode the server's sql_mode, so the plan's SET does not carry the mode the event had",
		[]string{"event sql_mode not preserved across ALTER"},
	},
}

// twoWays is the pair of failures of a form that differs from its readback: once each way.
func twoWays(forward, reverse string) []string { return []string{forward, reverse} }

// viewBody is the gap of a view whose body OceanBase stores in other words.
func viewBody(why string) gap {
	return gap{why, twoWays("CANONICALIZATION: user vs stored not empty", "CANONICALIZATION (reverse): stored vs user not empty")}
}

var obGapsSeen sync.Map

// knownGap runs body as the probe t. For a probe obGaps lists, body runs on a recorder
// in place of t, and the probe passes only if body failed with the recorded failures and
// no others.
func knownGap(t *testing.T, body func(t testing.TB)) {
	t.Helper()
	g, known := obGaps[t.Name()]
	if !known {
		body(t)
		return
	}
	obGapsSeen.Store(t.Name(), true)
	if wrong := gapVerdict(g, recorded(t, body)); wrong != "" {
		t.Fatal(wrong)
	}
	t.Logf("known gap: %s", g.why)
}

// gapped makes body the body of a probe's subtest, run under knownGap:
// t.Run(name, gapped(func(t testing.TB) { ... })).
func gapped(body func(t testing.TB)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		knownGap(t, body)
	}
}

// recorded runs body on a recorder and returns the failures it reported. body gets a
// goroutine of its own, so a Fatalf on the recorder ends body and not the test.
func recorded(t testing.TB, body func(t testing.TB)) []string {
	rec := &recorder{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		body(rec)
	}()
	<-done
	return rec.failures
}

// gapVerdict says why failures are not what g records, or "" when they are: as many
// failures, each holding its phrase.
func gapVerdict(g gap, failures []string) string {
	if len(failures) == 0 {
		return fmt.Sprintf("obGaps lists this probe (%s), but it passes: the gap is closed, remove the entry", g.why)
	}
	matches := len(failures) == len(g.fails)
	for i := 0; matches && i < len(g.fails); i++ {
		matches = strings.Contains(failures[i], g.fails[i])
	}
	if matches {
		return ""
	}
	return fmt.Sprintf("obGaps records %d failures for this probe (%s): %q\nIt now fails %d times:\n%s",
		len(g.fails), g.why, g.fails, len(failures), strings.Join(failures, "\n---\n"))
}

// recorder is the t of a probe that is expected to fail: it keeps the failures.
type recorder struct {
	testing.TB
	failures []string
}

func (r *recorder) Errorf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}
func (r *recorder) Error(args ...any) { r.failures = append(r.failures, fmt.Sprint(args...)) }
func (r *recorder) Fail()             { r.failures = append(r.failures, "Fail") }
func (r *recorder) Failed() bool      { return len(r.failures) > 0 }
func (r *recorder) Fatalf(format string, args ...any) {
	r.Errorf(format, args...)
	runtime.Goexit()
}
func (r *recorder) Fatal(args ...any) {
	r.Error(args...)
	runtime.Goexit()
}
func (r *recorder) FailNow() {
	r.Fail()
	runtime.Goexit()
}

// A skip would leave the probe unrun, which is the one thing a recorded gap must not be.
func (r *recorder) Skipf(format string, args ...any) { r.Fatalf("skipped: "+format, args...) }
func (r *recorder) Skip(args ...any)                 { r.Fatal(append([]any{"skipped: "}, args...)...) }
func (r *recorder) SkipNow()                         { r.Fatal("skipped") }

// rehomed returns a readback taken in database phys as it reads in database logical.
// OceanBase prints a foreign key's referenced table with its database, REFERENCES
// `phys`.`p`, where MySQL prints the bare name for a table of the same database. A
// harness that loads a readback under another name has to take the reference along,
// or the load answers 1824: no such table.
func rehomed(ddl, phys, logical string) string {
	return strings.ReplaceAll(ddl, "`"+phys+"`.", "`"+logical+"`.")
}

// storeAndRead applies createSQL to OceanBase and returns the form it stores. Unlike
// oracleConn.showCreate it keeps the error, so the caller can tell a form the engine
// rejects from an oracle that has stopped answering.
func storeAndRead(ctx context.Context, o *oracleConn, dbName, createSQL, table string) (string, error) {
	for _, stmt := range []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName),
		fmt.Sprintf("CREATE DATABASE %s", dbName),
		fmt.Sprintf("USE %s", dbName),
		createSQL,
	} {
		if _, err := o.db.ExecContext(ctx, stmt); err != nil {
			return "", fmt.Errorf("exec %q: %w", stmt, err)
		}
	}
	var name, ddl string
	row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, table))
	if err := row.Scan(&name, &ddl); err != nil {
		return "", fmt.Errorf("show create %s: %w", table, err)
	}
	return ddl, nil
}

// TestOracle_OceanBaseIndexScopeEcho pins that the index scope OceanBase echoes does not
// read as a schema change. 4.3.5 prints LOCAL on every secondary index, partitioned table
// or not, whether or not the DDL asked for it, so an unstated scope has to canonicalize
// equal to that readback or every indexed table reports a spurious index MODIFY and the
// migration grows a DROP INDEX and ADD INDEX nobody asked for.
//
// This needs its own probes because diffIdempotenceProbes has no secondary index at all,
// which is how the diff carried the raw scope for a while without anything failing.
func TestOracle_OceanBaseIndexScopeEcho(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, probe := range []struct{ id, create string }{
		{"unpartitioned", "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))"},
		{"partitioned", "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a)) PARTITION BY KEY(id) PARTITIONS 2"},
		{"global written", "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a) GLOBAL) PARTITION BY KEY(id) PARTITIONS 2"},
		// A repeated scope reads back as its last one.
		{"repeated scope", "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a) LOCAL GLOBAL, KEY kb (a) GLOBAL LOCAL) PARTITION BY KEY(id) PARTITIONS 2"},
		{"explicit partitions", "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a)) PARTITION BY KEY(id) (PARTITION p0, PARTITION p1)"},
		// OceanBase prints a template back as a template over named partitions, while
		// this catalog renders it expanded onto each one. Both have to describe the same
		// layout or a templated table phantom-diffs against its own readback.
		{"subpartition template", "CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)"},
	} {
		t.Run(probe.id, gapped(func(t testing.TB) {
			readback, err := storeAndRead(ctx, o, "obscope_"+strings.ReplaceAll(probe.id, " ", "_"), probe.create, "t")
			if err != nil {
				t.Fatalf("OceanBase would not take the form: %v", err)
			}
			sc := serverCharsetFor(o.version)
			n := NormalizerFor(o.version)
			userCat, userTbl := loadOneTable(t, sc, probe.create, "t")
			storedCat, storedTbl := loadOneTable(t, sc, readback, "t")
			d := DiffWithNormalizer(userCat, storedCat, n)
			_, _ = userTbl, storedTbl
			// Only the index half is asserted here. A partitioned table also differs at
			// the partition level on this engine, which is recorded elsewhere; this probe
			// exists to catch the scope echo reading as a change.
			for _, te := range d.Tables {
				for _, ie := range te.Indexes {
					t.Errorf("index %s reported %s, so the echoed scope read as a change\n  user:   %s\n  stored: %s",
						ie.Name, ie.Action, strings.TrimSpace(probe.create), strings.TrimSpace(readback))
				}
				if te.PartitionChanged {
					if probe.id == "subpartition template" {
						t.Errorf("the template does not canonicalize against its readback\n  user:   %s\n  stored: %s",
							strings.TrimSpace(probe.create), strings.TrimSpace(readback))
					} else {
						t.Logf("known gap: the partition form also differs for %s", probe.id)
					}
				}
			}
		}))
	}
}

// assertDiffEmptyAgainstStored takes a form and what the engine stored for it, and
// requires the two to canonicalize equal: each self-diffs empty, and user against stored
// is empty in both directions, since Drop and Add are direction-sensitive.
func assertDiffEmptyAgainstStored(t testing.TB, o *oracleConn, readback, userDDL, table string) {
	t.Helper()
	sc := serverCharsetFor(o.version)
	n := NormalizerFor(o.version)
	userCat, _ := loadOneTable(t, sc, userDDL, table)
	storedCat, _ := loadOneTable(t, sc, readback, table)

	if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
		t.Errorf("[%s] IDEMPOTENCE: self-diff of stored form not empty for %s:\n  stored: %s\n  diff: %s",
			o.name, table, strings.TrimSpace(readback), describeDiff(d))
	}
	if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
		t.Errorf("[%s] self-diff of user form not empty for %s:\n  user: %s\n  diff: %s",
			o.name, table, strings.TrimSpace(userDDL), describeDiff(d))
	}
	// The headline property: user form vs engine-stored form must collapse to empty.
	if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
		t.Errorf("[%s] CANONICALIZATION: user form vs stored form not empty for %s:\n  user:   %s\n  stored: %s\n  diff: %s",
			o.name, table, strings.TrimSpace(userDDL), strings.TrimSpace(readback), describeDiff(d))
	}
	if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
		t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s:\n  diff: %s",
			o.name, table, describeDiff(d))
	}
}

// TestKnownGapHoldsAProbeToItsRecordedFailures pins what obGaps lets through: a probe
// that fails as recorded, and nothing else. It needs no engine.
func TestKnownGapHoldsAProbeToItsRecordedFailures(t *testing.T) {
	g := gap{why: "why", fails: []string{"first", "second"}}
	tests := []struct {
		name     string
		failures []string
		passes   bool
	}{
		{"fails as recorded", []string{"the first thing", "the second thing"}, true},
		{"passes: the gap is closed", nil, false},
		{"fails once more", []string{"first", "second", "third"}, false},
		{"fails once less", []string{"first"}, false},
		{"fails otherwise", []string{"first", "another"}, false},
		{"fails in another order", []string{"second", "first"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if wrong := gapVerdict(g, tc.failures); (wrong == "") != tc.passes {
				t.Errorf("verdict %q, want passes = %v", wrong, tc.passes)
			}
		})
	}

	// The recorder keeps every failure, stops at a fatal one, and reports nothing.
	reached := false
	got := recorded(t, func(t testing.TB) {
		t.Errorf("one %d", 1)
		t.Error("two")
		if !t.Failed() {
			panic("the recorder forgot its failures")
		}
		t.Fatalf("three")
		reached = true
	})
	if want := []string{"one 1", "two", "three"}; !slices.Equal(got, want) || reached {
		t.Errorf("recorded %q (ran past the fatal failure: %v), want %q", got, reached, want)
	}
	if got := recorded(t, func(t testing.TB) { t.Skipf("no %s", "readback") }); !slices.Equal(got, []string{"skipped: no readback"}) {
		t.Errorf("a skip was recorded as %q, want a failure", got)
	}
}

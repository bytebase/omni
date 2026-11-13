package catalog

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// The OceanBase oracle for this package: every TestOracle_* that has moved off the
// MySQL container asks it. The ones that have not moved still call connectOracle.
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
func oceanBaseOracle(t *testing.T) *oracleConn {
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
// to compare and the probe ends there. refused holds every entry to the engine in both
// directions: a refusal that is not listed fails its test, and so does a listed probe
// that no longer ends in its refusal.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6.
var obRefused = map[string]refusal{
	"TestOracle_DiffIdempotence/oceanbase/binary-modifier-after-charset":               {1064, "BINARY after CHARACTER SET in a column's type"},
	"TestOracle_DiffIdempotence/oceanbase/adjacent-literal-default-double-quoted":      {1064, `a double-quoted literal with a neighbour, "x" 'y', in DEFAULT`},
	"TestOracle_MigrationIdempotence/oceanbase/binary-modifier-after-charset":          {1064, "BINARY after CHARACTER SET in a column's type"},
	"TestOracle_MigrationIdempotence/oceanbase/adjacent-literal-default-double-quoted": {1064, `a double-quoted literal with a neighbour, "x" 'y', in DEFAULT`},
	"TestOracle_PartitionChangeDetected/oceanbase/toggle-linear":                       {1064, "LINEAR KEY partitioning"},
}

var obRefusedSeen sync.Map

// refused ends a probe whose statement the engine would not take. err is the engine's
// answer, or nil when the helper that gave up kept it in o.refusal. On OceanBase the
// probe passes only if obRefused records this subtest with the error number the engine
// answered; anything else fails. The caller returns after it.
//
// On the MySQL oracle, which the tests that have not moved still ask, it skips, as
// those tests always have.
func (o *oracleConn) refused(t *testing.T, err error, format string, args ...any) {
	t.Helper()
	if err == nil {
		err = o.refusal
	}
	o.refusal = nil
	why := fmt.Sprintf(format, args...)
	if o.name != "oceanbase" {
		t.Skipf("[%s] %s: %v", o.name, why, err)
	}
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
		t.Run(probe.id, func(t *testing.T) {
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
		})
	}
}

// assertDiffEmptyAgainstStored takes a form and what the engine stored for it, and
// requires the two to canonicalize equal: each self-diffs empty, and user against stored
// is empty in both directions, since Drop and Add are direction-sensitive.
func assertDiffEmptyAgainstStored(t *testing.T, o *oracleConn, readback, userDDL, table string) {
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

// rowFormats is what OceanBase stores for each ROW_FORMAT a table is written with.
// IgnoreRowFormat rests on the first three: DYNAMIC is what a table that declares
// nothing stores, so a readback's DYNAMIC is not a declaration.
var rowFormats = []struct{ written, stored string }{
	{"", "DYNAMIC"},
	{"DEFAULT", "DYNAMIC"},
	{"DYNAMIC", "DYNAMIC"},
	{"COMPACT", "COMPACT"},
	{"REDUNDANT", "REDUNDANT"},
	{"COMPRESSED", "COMPRESSED"},
	{"CONDENSED", "CONDENSED"},
}

// TestOracle_OceanBaseRowFormat holds rowFormats to the engine, and requires each
// written form to canonicalize equal to what the engine stores for it. Before
// IgnoreRowFormat knew DYNAMIC, every table of every schema differed from its own
// readback by one MODIFY.
//
// The second half is the way back: a table that drops its ROW_FORMAT is altered with
// ROW_FORMAT=DEFAULT, and that has to store what a table declaring none stores.
func TestOracle_OceanBaseRowFormat(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, rf := range rowFormats {
		name, create := "unstated", "CREATE TABLE t (id INT PRIMARY KEY)"
		if rf.written != "" {
			name, create = rf.written, create+" ROW_FORMAT="+rf.written
		}
		t.Run(name, func(t *testing.T) {
			readback, err := storeAndRead(ctx, o, "obrf_"+strings.ToLower(name), create, "t")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(readback, " ROW_FORMAT = "+rf.stored+" ") {
				t.Errorf("recorded ROW_FORMAT = %s, but the engine stored\n%s", rf.stored, readback)
			}
			assertDiffEmptyAgainstStored(t, o, readback, create, "t")
		})
	}
	n := NormalizerFor(o.version)
	for _, from := range []string{"COMPACT", "COMPRESSED"} {
		t.Run("from "+from+" back to unstated", func(t *testing.T) {
			assertApplyCorrect(t, o, n, migrationProbe{
				id:      "rowformat-reset-" + strings.ToLower(from),
				table:   "t",
				fromDDL: "CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=" + from,
				toDDL:   "CREATE TABLE t (id INT PRIMARY KEY)",
			})
		})
	}
}

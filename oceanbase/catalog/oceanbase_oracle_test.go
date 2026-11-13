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

	// BINARY before CHARACTER SET is OceanBase's order (`VARCHAR string_length_i
	// opt_binary opt_charset`), and the parser here reads only the order OceanBase
	// refuses, so it cannot load a column OceanBase takes.
	"TestOracle_DiffIdempotence/oceanbase/binary-modifier-before-charset":      binaryBeforeCharset,
	"TestOracle_MigrationIdempotence/oceanbase/binary-modifier-before-charset": binaryBeforeCharset,

	// A multi-valued index. The parser does not read json_query's RETURNING clause, so a
	// table that has such an index does not load from its own readback.
	"TestOracle_IndexDiffIdempotence/oceanbase/mv-unsigned":  multiValuedIndex,
	"TestOracle_IndexDiffIdempotence/oceanbase/mv-char":      multiValuedIndex,
	"TestOracle_IndexDiffIdempotence/oceanbase/mv-decimal":   multiValuedIndex,
	"TestOracle_IndexDiffIdempotence/oceanbase/mv-date":      multiValuedIndex,
	"TestOracle_IndexDiffIdempotence/oceanbase/mv-composite": multiValuedIndex,

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

	// ---- Migrations: what the generator writes and OceanBase will not run. ----

	// A foreign key has no index of its own on OceanBase. The catalog gives it one, as
	// MySQL does, and the plan then drops an index that was never there.
	"TestOracle_ForeignKeyMigrationApplyCorrectness/oceanbase/drop-fk-and-backing-index":           backingIndex,
	"TestOracle_ForeignKeyMigrationApplyCorrectness/oceanbase/drop-fk-unnamed":                     backingIndex,
	"TestOracle_ForeignKeyMigrationApplyCorrectness/oceanbase/modify-fk-column":                    backingIndex,
	"TestOracle_ForeignKeyMultiTableApplyCorrectness/oceanbase/drop-column-and-its-fk":             backingIndex,
	"TestOracle_ForeignKeyMultiTableApplyCorrectness/oceanbase/drop-fk-reuse-name-for-other-index": backingIndex,
	"TestOracle_ForeignKeyMultiTableApplyCorrectness/oceanbase/drop-referenced-table":              backingIndex,
	"TestOracle_ForeignKeyMultiTableApplyCorrectness/oceanbase/drop-two-fks-sharing-index":         backingIndex,
	"TestOracle_ForeignKeyMultiTableApplyCorrectness/oceanbase/drop-two-unnamed-fks-sharing-index": backingIndex,
	"TestOracle_AutoIncGroupingFKImplicitBacking/oceanbase/aig-fk-implicit-backing-drop":           backingIndex,
	"TestOracle_AutoIncGroupingFKImplicitBacking/oceanbase/aig-fk-rename-in-place-to-ai":           backingIndex,

	// One ALTER TABLE that MySQL needs whole, to keep an AUTO_INCREMENT column keyed at
	// every statement's end, and that OceanBase takes only in parts.
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-add-col-promote-pk":           primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-drop-composite-pk-and-column": primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-drop-pk-and-column":           primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-drop-pk-deai":                 primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-pk-keep-ai-demote-member":     primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-pk-relocate-ai-to-secondary":  primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-pk-swap-drop-secondary":       primaryKeyAndMore,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-multi-drop-plus-add":          dropAndAddIndex,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-replace-backing-key":          dropAndAddIndex,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-swap-backing-key":             dropAndAddIndex,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-widen-ai-key-replaced":        dropAndAddIndex,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-drop-key-and-column": {
		"the plan drops an index and a column in one ALTER TABLE, which OceanBase will not combine",
		[]string{"Error 1235 (0A000): table redefinition and drop index in single statment not supported"},
	},
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-myisam-shared-key-ai-migrate": {
		"the plan drops an index and a column in one ALTER TABLE, which OceanBase will not combine",
		[]string{"Error 1235 (0A000): table redefinition and drop index in single statment not supported"},
	},
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-modify-gain-ai-new-key": {
		"the plan makes a column AUTO_INCREMENT and adds its key in one ALTER TABLE, which OceanBase will not combine",
		[]string{"Error 1235 (0A000): modify auto_increment and add index in single statment not supported"},
	},
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-pk-drop-member-drop-newkey":           dropPrimaryKeyColumn,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-pk-shrink-member-drop":                dropPrimaryKeyColumn,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-ai-migrate-columns":                   moveAutoIncrement,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-myisam-shared-key-ai-migrate-survive": moveAutoIncrement,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-add-col-new-pk":                       addKeyedAutoIncrement,
	"TestOracle_AutoIncGroupingApplyCorrectness/oceanbase/aig-add-col-promote-pk-drop-old":          addKeyedAutoIncrement,

	// What OceanBase has no statement for.
	"TestOracle_PartitionApplyCorrectness/oceanbase/remove-partitioning":                 removePartitioning,
	"TestOracle_PartitionApplyCorrectness/oceanbase/remove-partitioning-range":           removePartitioning,
	"TestOracle_PartitionApplyCorrectness/oceanbase/remove-partitioning-drop-key-column": removePartitioning,
	"TestOracle_PartitionApplyCorrectness/oceanbase/repartition-drop-old-key-column":     removePartitioning,
	"TestOracle_PartitionApplyCorrectness/oceanbase/repartition-drop-special-col":        removePartitioning,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/add-srid-none-to-4326":            alterSRID,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/add-srid-none-to-zero":            alterSRID,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/change-srid-4326-to-zero":         alterSRID,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/change-srid-zero-to-4326":         alterSRID,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/remove-srid-4326-to-none":         alterSRID,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/remove-srid-zero-to-none":         alterSRID,
	"TestOracle_SRIDPresenceApplyCorrectness/oceanbase/add-column-with-srid-zero": {
		"OceanBase will not add a NOT NULL geometry column, even to an empty table",
		[]string{"Error 1138 (22004): Invalid use of NULL value"},
	},
	"TestOracle_MigrationApplyCorrectness/oceanbase/modify-generated-expr": {
		"OceanBase will not change a generated column's expression in place",
		[]string{"Error 1235 (0A000): Modify generated column definition not supported"},
	},

	// The parser again: a multi-valued index, and a trigger named with its database, which
	// is how the generator writes CREATE TRIGGER and how OceanBase then keeps it.
	"TestOracle_IndexMigrationApplyCorrectness/oceanbase/create-with-mv-index":         multiValuedIndex,
	"TestOracle_IndexMigrationApplyCorrectness/oceanbase/add-mv-index":                 multiValuedIndex,
	"TestOracle_IndexMigrationApplyCorrectness/oceanbase/drop-mv-index":                multiValuedIndex,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/create-before-insert":       qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/create-after-update":        qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/create-before-delete":       qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/create-begin-end":           qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/change-body":                qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/change-timing":              qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/change-event":               qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/multi-add-drop":             qualifiedTrigger,
	"TestOracle_TriggerMigrationApplyCorrectness/oceanbase/modify-table-move-drop-old": qualifiedTrigger,

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

var (
	backingIndex = gap{
		"OceanBase makes no index for a foreign key; the catalog gives it one, as MySQL does, and the plan drops what was never there",
		[]string{"Error 1091 (42000): Can't DROP"},
	}
	primaryKeyAndMore = gap{
		"the plan changes the primary key and something else in one ALTER TABLE, which OceanBase takes one at a time",
		[]string{"Error 1235 (0A000): Multiple complex DDLs about primary in single stmt not supported"},
	}
	dropAndAddIndex = gap{
		"the plan drops one index and adds another in one ALTER TABLE, which OceanBase refuses",
		[]string{"Error 4179 (HY000): Dropping and adding indexes at the same time"},
	}
	dropPrimaryKeyColumn = gap{
		"the plan drops a column of the primary key, which OceanBase will not do",
		[]string{"Error 1235 (0A000): drop rowkey column is not supported"},
	}
	moveAutoIncrement = gap{
		"the plan moves AUTO_INCREMENT from one column to another in one ALTER TABLE; OceanBase reads the second MODIFY before it has taken the first",
		[]string{"Error 1075 (42000): Incorrect table definition; there can be only one auto column"},
	}
	addKeyedAutoIncrement = gap{
		"the plan adds an AUTO_INCREMENT column and its primary key in one ALTER TABLE; OceanBase looks for the key before the column is there",
		[]string{"Error 1176 (42000): Key 'seq' doesn't exist in table 't'"},
	}
	removePartitioning = gap{
		"OceanBase has no ALTER TABLE ... REMOVE PARTITIONING, which the plan writes for a table that stops being partitioned",
		[]string{"near 'PARTITIONING' at line 1"},
	}
	alterSRID = gap{
		"OceanBase will not change a geometry column's SRID in place",
		[]string{"Error 1235 (0A000): alter geometry srid not supported"},
	}
	qualifiedTrigger = gap{
		"the parser does not read CREATE TRIGGER `db`.`name`, which the generator writes and OceanBase keeps as written, so the result does not load",
		[]string{"reload of result SDL failed"},
	}
)

var binaryBeforeCharset = gap{
	"OceanBase takes VARCHAR(10) BINARY CHARACTER SET ascii and refuses CHARACTER SET ascii BINARY; the parser reads the second and not the first",
	[]string{"LoadSQL failed"},
}

var multiValuedIndex = gap{
	"OceanBase prints a multi-valued index as json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch), which the parser does not read",
	[]string{"LoadSQL failed"},
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

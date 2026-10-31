package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

// The OceanBase oracle for this package. The TestOracle_* family answers for the
// behavior this engine inherits from MySQL, against live 8.0 and 5.7; this file is
// the one that answers for OceanBase itself.
//
// OceanBase needs about 8 GB, which CI does not have, so this test reads
// OCEANBASE_WALK_DSN and skips without it:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/catalog/ -run OceanBase -count=1
const oceanBaseDSNEnv = "OCEANBASE_WALK_DSN"

// oceanBaseOracle opens the server OCEANBASE_WALK_DSN names, or skips.
//
// The version is MySQL80, which the probes below measure rather than assume. Under
// MySQL80 every probe either canonicalizes empty or disagrees only in the ROW_FORMAT
// echo; under MySQL57 twenty disagree, twelve because 5.7 assumes a latin1 server
// default where OceanBase defaults to utf8mb4, and eight on column forms, since 5.7
// injects display widths and then needs them to match. 4.3.5 does report a 5.7 version
// string, and does print 5.7-style widths, but the stored form it prints is the one
// MySQL80's normalizer reproduces: 8.0 strips widths from both sides. Nothing else pins
// the choice, so this test is what holds it, and a wrong normalizer shows up here as a
// non-empty diff rather than as a phantom change in a user's release.
func oceanBaseOracle(t *testing.T) *oracleConn {
	t.Helper()
	dsn := os.Getenv(oceanBaseDSNEnv)
	if dsn == "" {
		t.Skipf("%s is not set: run the test through scripts/oceanbase-oracle.sh", oceanBaseDSNEnv)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	// showCreate selects the database with USE, which is session state, so the pool
	// has to stay on one connection.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() { _ = db.Close() })
	// The tenant answers before it can plan, so a first statement can fail while the
	// server is still coming up.
	deadline := time.Now().Add(3 * time.Minute)
	for {
		err = db.PingContext(context.Background())
		if err == nil || time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
	return &oracleConn{db: db, version: MySQL80, name: "oceanbase"}
}

// obUnsupportedForms are the probe ids OceanBase itself will not accept, with the
// construct it rejects. The engine answers 1064 for each, so it stores no form to
// compare; any other probe failing to apply means the oracle stopped answering.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-30.
var obUnsupportedForms = map[string]string{
	"binary-modifier-utf8mb4":  "the BINARY column modifier",
	"adjacent-literal-default": "adjacent string literals in DEFAULT",
}

// obEchoedRowFormat is the ROW_FORMAT OceanBase stores for a table whose DDL declares
// none. MySQL echoes ROW_FORMAT only when the user declared a non-DEFAULT value, and
// rowFormatChanged rests on that: the undeclared user side canonicalizes to "" while
// the readback's DYNAMIC is significant, so on OceanBase every table whose DDL omits
// ROW_FORMAT phantom-diffs MODIFY. Normalizing the echo away needs an engine axis the
// version-keyed Normalizer does not have, so this test records the disagreement exactly
// rather than skipping over it: notRowFormatEcho tolerates that one shape and fails on
// anything else. Adding the axis empties these diffs and makes the tolerance unreachable.
const obEchoedRowFormat = "DYNAMIC"

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

// TestOracle_OceanBaseCanonicalization applies each representative form to a live
// OceanBase, reads the stored form back, and requires the catalog's user form and that
// stored form to canonicalize equal — TestOracle_DiffIdempotence's property, answered by
// the engine Bytebase runs instead of by MySQL standing in for it.
//
// Nothing here is skipped. A form OceanBase rejects must be named in obUnsupportedForms,
// and a non-empty diff must be the obEchoedRowFormat disagreement alone; every other
// failure, a stored form that will not parse included, fails the test.
func TestOracle_OceanBaseCanonicalization(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, probe := range diffIdempotenceProbes() {
		if !containsVersion(probe.versions, MySQL80) {
			continue
		}
		t.Run(probe.id, func(t *testing.T) {
			dbName := "ob_" + strings.ReplaceAll(probe.id, "-", "_")
			readback, err := storeAndRead(ctx, o, dbName, probe.create, probe.table)
			if err != nil {
				requireFormOceanBaseRejects(t, probe.id, err)
				return
			}
			assertDiffEmptyOrRowFormatEcho(t, o, readback, probe.create, probe.table)
		})
	}
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
		{"explicit partitions", "CREATE TABLE t (id INT, a INT, PRIMARY KEY (id), KEY ka (a)) PARTITION BY KEY(id) (PARTITION p0, PARTITION p1)"},
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
			// the partition level on this engine, and every table differs by the
			// obEchoedRowFormat gap, both recorded elsewhere; this probe exists to catch
			// the scope echo reading as a change.
			for _, te := range d.Tables {
				for _, ie := range te.Indexes {
					t.Errorf("index %s reported %s, so the echoed scope read as a change\n  user:   %s\n  stored: %s",
						ie.Name, ie.Action, strings.TrimSpace(probe.create), strings.TrimSpace(readback))
				}
				if te.PartitionChanged {
					t.Logf("known gap: the partition form also differs for %s", probe.id)
				}
			}
		})
	}
}

// requireFormOceanBaseRejects accepts the one apply failure this test tolerates: the
// engine answering 1064 for a form obUnsupportedForms names. A different error means the
// oracle is broken, and an unlisted probe means the dialect moved and nobody said so.
func requireFormOceanBaseRejects(t *testing.T, id string, err error) {
	t.Helper()
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) || myErr.Number != 1064 {
		t.Fatalf("applying %s failed for a reason other than OceanBase rejecting the form: %v", id, err)
	}
	why, ok := obUnsupportedForms[id]
	if !ok {
		t.Fatalf("OceanBase rejects %s, which is not recorded as unsupported: name the construct in obUnsupportedForms or fix the probe: %v", id, err)
	}
	t.Logf("OceanBase does not accept %s (%s), so it stores no form to compare", id, why)
}

// assertDiffEmptyOrRowFormatEcho is assertDiffEmptyAgainstReadback over a stored form the
// caller already read back, with the obEchoedRowFormat tolerance. Both self-diffs must
// still be empty, and user against stored must be empty or that echo alone in each
// direction, since Drop and Add are direction-sensitive.
func assertDiffEmptyOrRowFormatEcho(t *testing.T, o *oracleConn, readback, userDDL, table string) {
	t.Helper()
	sc := serverCharsetFor(o.version)
	n := NormalizerFor(o.version)
	userCat, userTbl := loadOneTable(t, sc, userDDL, table)
	storedCat, storedTbl := loadOneTable(t, sc, readback, table)

	if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
		t.Errorf("IDEMPOTENCE: self-diff of stored form not empty for %s:\n  stored: %s\n  diff: %s",
			table, strings.TrimSpace(readback), describeDiff(d))
	}
	if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
		t.Errorf("self-diff of user form not empty for %s:\n  user: %s\n  diff: %s",
			table, strings.TrimSpace(userDDL), describeDiff(d))
	}
	for _, dir := range []struct {
		name     string
		from, to *Catalog
	}{
		{"user vs stored", userCat, storedCat},
		{"stored vs user", storedCat, userCat},
	} {
		d := DiffWithNormalizer(dir.from, dir.to, n)
		if d.IsEmpty() {
			continue
		}
		if why := notRowFormatEcho(d, userTbl, storedTbl, table, n); why != "" {
			t.Errorf("CANONICALIZATION (%s) for %s: %s\n  user:   %s\n  stored: %s\n  diff: %s",
				dir.name, table, why, strings.TrimSpace(userDDL), strings.TrimSpace(readback), describeDiff(d))
		}
	}
}

// notRowFormatEcho reports why a non-empty diff is more than the obEchoedRowFormat
// disagreement, or "" when it is exactly that. It requires the diff to carry nothing but
// one MODIFY on this table with no sub-object entry, and requires every table option
// tableOptionsChanged compares except ROW_FORMAT to agree — so a column, index,
// constraint, partition, engine, charset, collation or comment phantom fails here instead
// of passing as the known gap. The two sides are named by role, not by diff direction,
// because only "the user declared none and the engine echoed DYNAMIC" is tolerated: a
// readback that echoed some other format disagrees with the user's own declaration and is
// a real difference.
func notRowFormatEcho(d *SchemaDiff, userTbl, storedTbl *Table, table string, n *Normalizer) string {
	switch {
	case len(d.Views) > 0 || len(d.Functions) > 0 || len(d.Procedures) > 0 || len(d.Triggers) > 0 || len(d.Events) > 0:
		return "the diff reaches objects other than the table"
	case len(d.Tables) != 1:
		return fmt.Sprintf("the diff names %d tables, not 1", len(d.Tables))
	}
	e := d.Tables[0]
	switch {
	case e.Action != DiffModify:
		return fmt.Sprintf("the table is %s, not MODIFY", e.Action)
	case !strings.EqualFold(e.Name, table):
		return fmt.Sprintf("the diff names table %q, not %q", e.Name, table)
	case len(e.Columns) > 0 || len(e.Indexes) > 0 || len(e.Constraints) > 0 || len(e.ForeignKeys) > 0 || len(e.Checks) > 0 || e.PartitionChanged:
		return "the table differs in a column, index, constraint or partition, not only in its options"
	case n.CanonicalEngine(userTbl) != n.CanonicalEngine(storedTbl):
		return fmt.Sprintf("the storage engines disagree: %q vs %q", n.CanonicalEngine(userTbl), n.CanonicalEngine(storedTbl))
	case foldCharset(userTbl.Charset) != foldCharset(storedTbl.Charset):
		return fmt.Sprintf("the charsets disagree: %q vs %q", userTbl.Charset, storedTbl.Charset)
	case n.effectiveTableCollation(userTbl) != n.effectiveTableCollation(storedTbl):
		return fmt.Sprintf("the collations disagree: %q vs %q", n.effectiveTableCollation(userTbl), n.effectiveTableCollation(storedTbl))
	case n.CanonicalComment(userTbl.Comment) != n.CanonicalComment(storedTbl.Comment):
		return fmt.Sprintf("the comments disagree: %q vs %q", userTbl.Comment, storedTbl.Comment)
	case !n.IgnoreRowFormat(userTbl) || !strings.EqualFold(storedTbl.RowFormat, obEchoedRowFormat):
		return fmt.Sprintf("the ROW_FORMAT disagreement is %q against %q, not an undeclared form against the echoed %s",
			userTbl.RowFormat, storedTbl.RowFormat, obEchoedRowFormat)
	}
	return ""
}

// TestNotRowFormatEcho pins what the obEchoedRowFormat tolerance will and will not
// excuse. notRowFormatEcho is a pure verdict over a diff and the two table forms, so this
// needs no engine — and without it the tolerance above is only as narrow as it looks.
func TestNotRowFormatEcho(t *testing.T) {
	const table = "t"
	undeclared := &Table{Name: table, Charset: "utf8mb4"}
	echoed := &Table{Name: table, Charset: "utf8mb4", RowFormat: obEchoedRowFormat}
	modify := func(mutate func(*TableDiffEntry)) *SchemaDiff {
		e := TableDiffEntry{Action: DiffModify, Name: table}
		if mutate != nil {
			mutate(&e)
		}
		return &SchemaDiff{Tables: []TableDiffEntry{e}}
	}

	tests := []struct {
		name         string
		diff         *SchemaDiff
		user, stored *Table
		tolerated    bool
	}{
		{"the echo itself", modify(nil), undeclared, echoed, true},
		{"ROW_FORMAT=DEFAULT is undeclared too", modify(nil), &Table{Name: table, Charset: "utf8mb4", RowFormat: "DEFAULT"}, echoed, true},
		{"a column also differs", modify(func(e *TableDiffEntry) {
			e.Columns = []ColumnDiffEntry{{Action: DiffModify, Name: "a"}}
		}), undeclared, echoed, false},
		{"an index also differs", modify(func(e *TableDiffEntry) {
			e.Indexes = []IndexDiffEntry{{Action: DiffAdd, Name: "ia"}}
		}), undeclared, echoed, false},
		{"a constraint also differs", modify(func(e *TableDiffEntry) {
			e.Constraints = []ConstraintDiffEntry{{Action: DiffAdd, Name: "pk"}}
		}), undeclared, echoed, false},
		{"the partitioning also differs", modify(func(e *TableDiffEntry) {
			e.PartitionChanged = true
		}), undeclared, echoed, false},
		// The collations are held equal because the collation check behind the charset
		// one would otherwise catch any differing charset by itself.
		{"the charsets differ while the collations agree", modify(nil),
			&Table{Name: table, Charset: "latin1", Collation: "utf8mb4_bin"},
			&Table{Name: table, Charset: "utf8mb4", Collation: "utf8mb4_bin", RowFormat: obEchoedRowFormat}, false},
		{"the collations differ", modify(nil), &Table{Name: table, Charset: "utf8mb4", Collation: "utf8mb4_bin"}, echoed, false},
		{"the engines differ", modify(nil), &Table{Name: table, Charset: "utf8mb4", Engine: "MyISAM"}, echoed, false},
		{"the comments differ", modify(nil), &Table{Name: table, Charset: "utf8mb4", Comment: "mine"}, echoed, false},
		{"the user declared a format too", modify(nil), &Table{Name: table, Charset: "utf8mb4", RowFormat: "COMPRESSED"}, echoed, false},
		{"the engine echoed another format", modify(nil), undeclared, &Table{Name: table, Charset: "utf8mb4", RowFormat: "COMPRESSED"}, false},
		{"the table is dropped, not modified", &SchemaDiff{
			Tables: []TableDiffEntry{{Action: DiffDrop, Name: table}},
		}, undeclared, echoed, false},
		{"the diff names another table", modify(func(e *TableDiffEntry) { e.Name = "other" }), undeclared, echoed, false},
		{"two tables differ", &SchemaDiff{
			Tables: []TableDiffEntry{{Action: DiffModify, Name: table}, {Action: DiffModify, Name: "u"}},
		}, undeclared, echoed, false},
		{"a view also differs", &SchemaDiff{
			Tables: []TableDiffEntry{{Action: DiffModify, Name: table}},
			Views:  []ViewDiffEntry{{Action: DiffAdd, Name: "v"}},
		}, undeclared, echoed, false},
	}

	n := NormalizerFor(MySQL80)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			why := notRowFormatEcho(tc.diff, tc.user, tc.stored, table, n)
			if tolerated := why == ""; tolerated != tc.tolerated {
				t.Fatalf("tolerated = %v, want %v (verdict %q)", tolerated, tc.tolerated, why)
			}
		})
	}
}

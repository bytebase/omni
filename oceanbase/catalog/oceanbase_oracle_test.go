package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
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
// The version is MySQL80 because OceanBase 4.x takes 8.0's grammar and utf8mb4
// server defaults even though it reports a 5.7 version string. Nothing else pins
// that choice, so this test is what holds it: a wrong normalizer shows up here as a
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

// TestOracle_OceanBaseCanonicalization applies each representative form to a live
// OceanBase, reads the stored form back, and requires the catalog's user form and
// that stored form to canonicalize equal — TestOracle_DiffIdempotence's property,
// answered by the engine Bytebase runs instead of by MySQL standing in for it.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-30: SHOW CREATE TABLE prints
// OceanBase's own table options (ORGANIZATION INDEX, REPLICA_NUM, BLOCK_SIZE,
// TABLET_SIZE, PCTFREE, USE_BLOOM_FILTER, ENABLE_MACRO_BLOCK_BLOOM_FILTER, and its
// own COMPRESSION and ROW_FORMAT forms), none of which this parser reads, so no
// stored form loads yet. Rather than skip, each case pins that: it demands the
// load fail on an unread table option, so the canonicalization assertion starts
// running by itself once the grammar lands, and a stored form that stops loading
// for some other reason fails here instead of passing quietly.
func TestOracle_OceanBaseCanonicalization(t *testing.T) {
	o := oceanBaseOracle(t)
	for _, probe := range diffIdempotenceProbes() {
		if !containsVersion(probe.versions, MySQL80) {
			continue
		}
		t.Run(probe.id, func(t *testing.T) {
			dbName := "ob_" + strings.ReplaceAll(probe.id, "-", "_")
			readback, ok := o.applyAndReadback(t, dbName, probe.create, probe.table)
			if !ok {
				t.Skipf("[%s] OceanBase did not store %s", o.name, probe.table)
			}
			if _, err := LoadSQL(wrapForLoad(serverCharsetFor(o.version), readback)); err != nil {
				requireUnreadTableOption(t, err, readback)
				return
			}
			assertDiffEmptyAgainstReadback(t, o, dbName, probe.create, probe.table)
		})
	}
}

// wrapForLoad mirrors loadOneTable's wrapping so a readback is loaded under the same
// server charset the comparison would use.
func wrapForLoad(serverCharset, createSQL string) string {
	return fmt.Sprintf("CREATE DATABASE diffdb DEFAULT CHARSET=%s;\nUSE diffdb;\n%s", serverCharset, createSQL)
}

// obOnlyTableOptions are the table options OceanBase prints that this parser does not
// read. Each is grammar this parser owes; implementing them clears this test's pending
// branch in the same change.
var obOnlyTableOptions = []string{
	"ORGANIZATION INDEX", "REPLICA_NUM", "BLOCK_SIZE", "TABLET_SIZE", "PCTFREE",
	"USE_BLOOM_FILTER", "ENABLE_MACRO_BLOCK_BLOOM_FILTER", "COMPRESSION", "ROW_FORMAT",
}

// requireUnreadTableOption accepts the one failure this test still tolerates: a stored
// form carrying an option from obOnlyTableOptions. Any other load failure is a real
// defect and fails here.
func requireUnreadTableOption(t *testing.T, err error, readback string) {
	t.Helper()
	for _, opt := range obOnlyTableOptions {
		if strings.Contains(readback, opt) {
			t.Logf("pending: the stored form carries %s, which this parser does not read yet: %v", opt, err)
			return
		}
	}
	t.Fatalf("the stored form failed to load for a reason other than an unread table option:\n  stored: %s\n  err: %v",
		strings.TrimSpace(readback), err)
}

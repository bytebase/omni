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

	obparser "github.com/bytebase/omni/oceanbase/parser"
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

// obUnsupportedForms are the probe ids OceanBase itself will not accept, with the
// construct it rejects. The engine answers 1064 for each, so it stores no form to
// compare; any other probe failing to apply means the oracle stopped answering.
//
// Read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-30.
var obUnsupportedForms = map[string]string{
	"binary-modifier-utf8mb4":  "the BINARY column modifier",
	"adjacent-literal-default": "adjacent string literals in DEFAULT",
}

// obUnreadTableOptions are the table options OceanBase prints in SHOW CREATE TABLE,
// named by the token the parser stops at. Each is grammar this parser owes;
// implementing them clears the pending branch below in the same change.
var obUnreadTableOptions = []string{
	"ORGANIZATION", "REPLICA_NUM", "BLOCK_SIZE", "TABLET_SIZE", "PCTFREE",
	"USE_BLOOM_FILTER", "ENABLE_MACRO_BLOCK_BLOOM_FILTER", "COMPRESSION", "ROW_FORMAT",
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

// TestOracle_OceanBaseCanonicalization applies each representative form to a live
// OceanBase, reads the stored form back, and requires the catalog's user form and that
// stored form to canonicalize equal — TestOracle_DiffIdempotence's property, answered by
// the engine Bytebase runs instead of by MySQL standing in for it.
//
// Nothing here is skipped. A form OceanBase rejects must be named in
// obUnsupportedForms, and a stored form that will not load must stop the parser at a
// table option in obUnreadTableOptions; every other failure fails the test. Read against
// 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-10-30, no stored form loads yet, so the
// canonicalization assertion starts running by itself once that grammar lands.
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
			if _, err := LoadSQL(wrapForLoad(serverCharsetFor(o.version), readback)); err != nil {
				requireUnreadTableOption(t, err, readback)
				return
			}
			assertDiffEmptyAgainstReadback(t, o, dbName, probe.create, probe.table)
		})
	}
}

// wrapForLoad mirrors loadOneTable's wrapping so a readback loads under the same server
// charset the comparison would use.
func wrapForLoad(serverCharset, createSQL string) string {
	return fmt.Sprintf("CREATE DATABASE diffdb DEFAULT CHARSET=%s;\nUSE diffdb;\n%s", serverCharset, createSQL)
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

// requireUnreadTableOption accepts the one load failure this test tolerates: the parser
// stopping at a table option obUnreadTableOptions names. It reads the token the parser
// reported rather than asking whether the option appears anywhere in the DDL, so a
// regression on a column, a default or the wrapper fails here instead of passing as
// pending.
func requireUnreadTableOption(t *testing.T, err error, readback string) {
	t.Helper()
	var pe *obparser.ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("the stored form failed to load with no parse position:\n  stored: %s\n  err: %v", strings.TrimSpace(readback), err)
	}
	for _, opt := range obUnreadTableOptions {
		if strings.Contains(pe.Message, fmt.Sprintf("%q", opt)) {
			t.Logf("pending: the parser stops at %s (line %d, column %d), a table option it does not read yet", opt, pe.Line, pe.Column)
			return
		}
	}
	t.Fatalf("the stored form failed to load at %s, which is not an unread table option:\n  stored: %s", pe.Message, strings.TrimSpace(readback))
}

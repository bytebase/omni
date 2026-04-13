package catalog

import (
	"os"
	"strings"
	"testing"
)

// This file provides the shared infrastructure used by the "mysql-implicit-behavior"
// starmap scenario tests. Section workers in BATCHES 1-7 build on these helpers
// to run dual-assertion scenarios against both a real MySQL 8.0 container and
// the omni catalog.
//
// NOTE: asString is already defined in catalog_spotcheck_test.go and is reused
// here as-is; do not redeclare it.

// scenarioContainer wraps startContainer for naming consistency with the
// scenario helpers. The caller must defer the cleanup func.
func scenarioContainer(t *testing.T) (*mysqlContainer, func()) {
	t.Helper()
	return startContainer(t)
}

// scenarioReset drops and recreates the shared testdb database on the MySQL
// container and selects it. It uses t.Error rather than t.Fatal so that the
// calling test can continue and report additional diffs within one run.
func scenarioReset(t *testing.T, mc *mysqlContainer) {
	t.Helper()
	stmts := []string{
		"DROP DATABASE IF EXISTS testdb",
		"CREATE DATABASE testdb",
		"USE testdb",
	}
	for _, stmt := range stmts {
		if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
			t.Errorf("scenarioReset %q: %v", stmt, err)
		}
	}
}

// scenarioNewCatalog returns a fresh omni catalog with a testdb database
// created and selected. Uses t.Fatal on setup errors because without a
// working catalog nothing else can run.
func scenarioNewCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec("CREATE DATABASE testdb; USE testdb;", nil)
	if err != nil {
		t.Fatalf("scenarioNewCatalog parse error: %v", err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("scenarioNewCatalog exec error on stmt %d: %v", r.Index, r.Error)
		}
	}
	return c
}

// runOnBoth executes a (possibly multi-statement) DDL string on both the
// MySQL container and the omni catalog. Errors on either side are reported
// via t.Error so that the calling test can continue comparing remaining
// scenario state. Statements are split respecting quotes; individual
// statements are executed one at a time on the container side.
func runOnBoth(t *testing.T, mc *mysqlContainer, c *Catalog, ddl string) {
	t.Helper()

	for _, stmt := range splitStmts(ddl) {
		if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
			t.Errorf("mysql container DDL failed: %q: %v", stmt, err)
		}
	}

	results, err := c.Exec(ddl, nil)
	if err != nil {
		t.Errorf("omni catalog parse error for DDL %q: %v", ddl, err)
		return
	}
	for _, r := range results {
		if r.Error != nil {
			t.Errorf("omni catalog exec error on stmt %d: %v", r.Index, r.Error)
		}
	}
}

// oracleScan runs a single-row information_schema (or other) query against
// the MySQL container and scans into dests. Uses t.Error on failure so the
// test can continue.
func oracleScan(t *testing.T, mc *mysqlContainer, query string, dests ...any) {
	t.Helper()
	row := mc.db.QueryRowContext(mc.ctx, query)
	if err := row.Scan(dests...); err != nil {
		t.Errorf("oracleScan failed: %q: %v", query, err)
	}
}

// oracleRows runs a multi-row query against the MySQL container and returns
// the rows as a [][]any, converting []byte values to string for readability.
// Uses t.Error on failure and returns nil.
func oracleRows(t *testing.T, mc *mysqlContainer, query string) [][]any {
	t.Helper()
	rows, err := mc.db.QueryContext(mc.ctx, query)
	if err != nil {
		t.Errorf("oracleRows query failed: %q: %v", query, err)
		return nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Errorf("oracleRows columns failed: %v", err)
		return nil
	}

	var out [][]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			t.Errorf("oracleRows scan failed: %v", err)
			return out
		}
		for i, v := range vals {
			if b, ok := v.([]byte); ok {
				vals[i] = string(b)
			}
		}
		out = append(out, vals)
	}
	if err := rows.Err(); err != nil {
		t.Errorf("oracleRows iteration error: %v", err)
	}
	return out
}

// oracleShow runs a SHOW CREATE TABLE / VIEW / ... statement against the
// container and returns the second column (the CREATE statement text).
// The first column (name) and any trailing columns are discarded. Uses
// t.Error on failure and returns the empty string.
//
// Note: different SHOW CREATE variants return different numbers of columns
// (SHOW CREATE TABLE returns 2, SHOW CREATE VIEW returns 4, SHOW CREATE
// FUNCTION/PROCEDURE/TRIGGER/EVENT return 6 or 7). This helper scans
// dynamically so it works with any of them.
func oracleShow(t *testing.T, mc *mysqlContainer, stmt string) string {
	t.Helper()
	rows, err := mc.db.QueryContext(mc.ctx, stmt)
	if err != nil {
		t.Errorf("oracleShow query failed: %q: %v", stmt, err)
		return ""
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Errorf("oracleShow columns failed: %v", err)
		return ""
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			t.Errorf("oracleShow iteration error: %v", err)
		} else {
			t.Errorf("oracleShow returned no rows for %q", stmt)
		}
		return ""
	}
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		t.Errorf("oracleShow scan failed: %v", err)
		return ""
	}
	if len(vals) < 2 {
		t.Errorf("oracleShow %q: expected >=2 columns, got %d", stmt, len(cols))
		return ""
	}
	return asString(vals[1])
}

// assertStringEq reports a diff if got != want.
func assertStringEq(t *testing.T, label, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", label, got, want)
	}
}

// assertIntEq reports a diff if got != want.
func assertIntEq(t *testing.T, label string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %d, want %d", label, got, want)
	}
}

// assertBoolEq reports a diff if got != want.
func assertBoolEq(t *testing.T, label string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", label, got, want)
	}
}

// scenariosSkipIfShort skips the calling test when testing.Short() is true.
func scenariosSkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping scenario test in short mode")
	}
}

// scenariosSkipIfNoDocker skips the calling test when SKIP_SCENARIO_TESTS=1
// is set in the environment (used by CI lanes that cannot run docker).
func scenariosSkipIfNoDocker(t *testing.T) {
	t.Helper()
	if os.Getenv("SKIP_SCENARIO_TESTS") == "1" {
		t.Skip("SKIP_SCENARIO_TESTS=1 set; skipping scenario test")
	}
}

// splitStmts splits a possibly multi-statement DDL string into individual
// statements, respecting single quotes, double quotes, and backticks, and
// trimming empty results. It is a thin wrapper around splitStatements
// (defined in container_test.go) with extra trimming so scenario workers
// can write `splitStmts(ddl)` without importing two names.
func splitStmts(ddl string) []string {
	raw := splitStatements(ddl)
	out := raw[:0]
	for _, s := range raw {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

package catalog

import (
	"errors"
	"regexp"
	"slices"
	"strings"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// This file provides the shared infrastructure used by the "mysql-implicit-behavior"
// scenario tests. The scenarios are MySQL's
// (mysql/catalog/SCENARIOS-mysql-implicit-behavior.md); here each one runs on
// OceanBase and on the omni catalog, and asserts OceanBase's answer. Where
// the catalog still answers as MySQL does, the subtest asserts that beside
// it under a comment that starts "omni gap:".
//
// NOTE: asString is already defined in catalog_spotcheck_test.go and is reused
// here as-is; do not redeclare it.

// scenarioContainer is startContainer under the scenario tests' name for it.
// The scenarios rely on connection-scoped state (USE testdb, SET SESSION
// sql_mode), which holds because a tape's stream is one connection.
func scenarioContainer(t *testing.T) (*mysqlContainer, func()) {
	t.Helper()
	return startContainer(t)
}

// oracleRefuses runs stmt on OceanBase and reports unless it answers the
// error numbered code.
func oracleRefuses(t *testing.T, mc *mysqlContainer, stmt string, code uint16) {
	t.Helper()
	_, err := mc.db.ExecContext(mc.ctx, stmt)
	var me *mysqldriver.MySQLError
	switch {
	case err == nil:
		t.Errorf("oracle: %q ran, want error %d", stmt, code)
	case !errors.As(err, &me):
		t.Errorf("oracle: %q: %v", stmt, err)
	case me.Number != code:
		t.Errorf("oracle: %q: got %v, want error %d", stmt, err, code)
	}
}

// omniCode runs ddl on the catalog and returns the code of its first error:
// 0 when every statement runs, and -1 for an error with no code, such as a
// parse error.
func omniCode(c *Catalog, ddl string) int {
	results, err := c.Exec(ddl, nil)
	if err == nil {
		for _, r := range results {
			if r.Error != nil {
				err = r.Error
				break
			}
		}
	}
	if err == nil {
		return 0
	}
	var ce *Error
	if errors.As(err, &ce) {
		return ce.Code
	}
	return -1
}

// madeUpCount is the count a tape holds where OceanBase wrote the time into
// a name it made up: child_OBFK_0000000000000003 for the third such name a
// test reads (testcontainer.Codec.StandIn).
var madeUpCount = regexp.MustCompile(`(_OB(?:FK|CHECK)_)\d{16}`)

// withoutCount writes <n> for that count. Which made-up name a test happened
// to read third says nothing of the engine, so a test compares the name
// without it: child_OBFK_<n>.
func withoutCount(name string) string {
	return madeUpCount.ReplaceAllString(name, "${1}<n>")
}

// asOracleLists writes each name without its count and orders them as OceanBase's
// ORDER BY does, whatever their case, so the catalog's names read as the oracle's.
func asOracleLists(names []string) []string {
	out := make([]string, len(names))
	for i, name := range names {
		out[i] = withoutCount(name)
	}
	slices.SortFunc(out, func(a, b string) int { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) })
	return out
}

// scenarioReset drops and recreates the shared testdb database on OceanBase
// and selects it. It uses t.Error rather than t.Fatal so that the
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

// runOnBoth executes a (possibly multi-statement) DDL string on both
// OceanBase and the omni catalog. Errors on either side are reported
// via t.Error so that the calling test can continue comparing remaining
// scenario state. Statements are split respecting quotes; individual
// statements are executed one at a time on the container side.
func runOnBoth(t *testing.T, mc *mysqlContainer, c *Catalog, ddl string) {
	t.Helper()

	for _, stmt := range splitStmts(ddl) {
		if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
			t.Errorf("oracle DDL failed: %q: %v", stmt, err)
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
// OceanBase and scans into dests. Uses t.Error on failure so the
// test can continue.
func oracleScan(t *testing.T, mc *mysqlContainer, query string, dests ...any) {
	t.Helper()
	row := mc.db.QueryRowContext(mc.ctx, query)
	if err := row.Scan(dests...); err != nil {
		t.Errorf("oracleScan failed: %q: %v", query, err)
	}
}

// oracleRows runs a multi-row query against OceanBase and returns
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

// oracleShow runs a SHOW CREATE TABLE / VIEW / ... statement against
// OceanBase and returns the second column (the CREATE statement text).
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

// mysqlAtLeast reports whether the VERSION() string reports a server at
// or above the requested (major, minor, patch) triple. Non-numeric suffixes
// like "-log" or "-debug" are ignored. Malformed versions return false.
func mysqlAtLeast(ver string, wantMajor, wantMinor, wantPatch int) bool {
	// Strip anything after the first non-numeric/non-dot rune.
	cut := len(ver)
	for i, r := range ver {
		if (r < '0' || r > '9') && r != '.' {
			cut = i
			break
		}
	}
	parts := strings.Split(ver[:cut], ".")
	if len(parts) < 3 {
		return false
	}
	atoi := func(s string) int {
		n := 0
		for _, r := range s {
			if r < '0' || r > '9' {
				return -1
			}
			n = n*10 + int(r-'0')
		}
		return n
	}
	maj, min, pat := atoi(parts[0]), atoi(parts[1]), atoi(parts[2])
	if maj < 0 || min < 0 || pat < 0 {
		return false
	}
	if maj != wantMajor {
		return maj > wantMajor
	}
	if min != wantMinor {
		return min > wantMinor
	}
	return pat >= wantPatch
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

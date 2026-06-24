package catalog

import (
	"strings"
	"testing"
)

// TestApplicationTimePeriodCatalog: a named (application-time) PERIOD FOR <name>
// renders with a backtick-quoted period name and NOT NULL period columns, and
// coexists with SYSTEM_TIME (bitemporal). Grounded vs mariadb:11.8.8.
func TestApplicationTimePeriodCatalog(t *testing.T) {
	ddl := defineAndShow(t, "t", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e))")
	if !strings.Contains(ddl, "PERIOD FOR `app_time` (`s`, `e`)") {
		t.Errorf("missing app-time period clause:\n%s", ddl)
	}
	if !strings.Contains(ddl, "`s` date NOT NULL") || !strings.Contains(ddl, "`e` date NOT NULL") {
		t.Errorf("period columns should be NOT NULL:\n%s", ddl)
	}
}

// TestApplicationTimePeriodValidation pins the malformed-CREATE error codes,
// grounded vs mariadb:11.8.8.
func TestApplicationTimePeriodValidation(t *testing.T) {
	cases := []struct {
		name string
		ddl  string
		code int
	}{
		{"cols_dont_exist", "CREATE TABLE t (id INT, PERIOD FOR app_time(nope1, nope2))", 1054},
		{"two_app_periods", "CREATE TABLE t (id INT, s DATE, e DATE, s2 DATE, e2 DATE, PERIOD FOR p1(s,e), PERIOD FOR p2(s2,e2))", 4154},
		{"same_col", "CREATE TABLE t (id INT, s DATE, PERIOD FOR app_time(s, s))", 1110},
		{"wrong_type", "CREATE TABLE t (id INT, s VARCHAR(10), e VARCHAR(10), PERIOD FOR app_time(s, e))", 1063},
		{"name_is_column", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR id(s, e))", 1060},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := execErr(t, tc.ddl)
			catErr, ok := err.(*Error)
			if !ok {
				t.Fatalf("expected *Error, got %v", err)
			}
			if catErr.Code != tc.code {
				t.Errorf("Code = %d, want %d (message: %q)", catErr.Code, tc.code, catErr.Message)
			}
		})
	}
}

// TestWithoutOverlapsCatalog: WITHOUT OVERLAPS renders on the period key part.
func TestWithoutOverlapsCatalog(t *testing.T) {
	for _, tc := range []struct{ name, ddl, want string }{
		{"unique", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), UNIQUE (id, app_time WITHOUT OVERLAPS))", "`app_time` WITHOUT OVERLAPS"},
		{"primary", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), PRIMARY KEY (id, app_time WITHOUT OVERLAPS))", "`app_time` WITHOUT OVERLAPS"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if ddl := defineAndShow(t, "t", tc.ddl); !strings.Contains(ddl, tc.want) {
				t.Errorf("missing %q:\n%s", tc.want, ddl)
			}
		})
	}
}

// TestWithoutOverlapsValidation: a WITHOUT OVERLAPS key part must name the
// application-time period (else 4156). Grounded vs 11.8.8.
func TestWithoutOverlapsValidation(t *testing.T) {
	for _, tc := range []struct {
		name, ddl string
	}{
		{"non_period_col", "CREATE TABLE t (id INT, s DATE, e DATE, UNIQUE(id, e WITHOUT OVERLAPS))"},
		{"no_period_at_all", "CREATE TABLE t (id INT, x INT, UNIQUE(id, x WITHOUT OVERLAPS))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := execErr(t, tc.ddl)
			if catErr, ok := err.(*Error); !ok || catErr.Code != 4156 {
				t.Fatalf("want *Error 4156, got %v", err)
			}
		})
	}
}

// TestWithoutOverlapsValidationAllPaths: the 4156 period-name check applies to
// CREATE INDEX and ALTER ADD too, not just CREATE TABLE.
func TestWithoutOverlapsValidationAllPaths(t *testing.T) {
	for _, tc := range []struct{ name, stmt string }{
		{"create_unique_index", "CREATE UNIQUE INDEX ux ON t (id, e WITHOUT OVERLAPS)"},
		{"alter_add_unique_key", "ALTER TABLE t ADD UNIQUE KEY ux (id, e WITHOUT OVERLAPS)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec("CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e))", nil)
			r, _ := c.Exec(tc.stmt, &ExecOptions{ContinueOnError: true})
			if catErr, ok := r[0].Error.(*Error); !ok || catErr.Code != 4156 {
				t.Fatalf("want 4156, got %v", r[0].Error)
			}
		})
	}
}

// TestKeyPartZeroPrefixLength: an explicit zero-length key-part prefix `col(0)`
// is rejected with 1391 on every key path and every column type (grounded vs
// mariadb:11.8.8). `col(0)` is distinct from a bare `col` (no prefix) — the
// latter is valid. This also closes the WITHOUT OVERLAPS bare-part bypass:
// `app_time(0) WITHOUT OVERLAPS` parses past the bare-part check (Length 0) but
// must still be rejected 1391 (length-0 outranks the WO-shape 1064).
func TestKeyPartZeroPrefixLength(t *testing.T) {
	// Negative: every CREATE TABLE key path rejects col(0) with 1391.
	createCases := []struct{ name, ddl string }{
		{"primary_key", "CREATE TABLE t (c VARCHAR(20), PRIMARY KEY (c(0)))"},
		{"unique", "CREATE TABLE t (id INT, c VARCHAR(20), UNIQUE (c(0)))"},
		{"key", "CREATE TABLE t (id INT, c VARCHAR(20), KEY (c(0)))"},
		{"blob", "CREATE TABLE t (b BLOB, KEY (b(0)))"},                  // 1391, not 1170
		{"fulltext", "CREATE TABLE t (t1 TEXT, FULLTEXT KEY k (t1(0)))"}, // FULLTEXT allows a prefix; 0 is still 1391
		{"without_overlaps", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), UNIQUE (id, app_time(0) WITHOUT OVERLAPS))"},
	}
	for _, tc := range createCases {
		t.Run(tc.name, func(t *testing.T) {
			err := execErr(t, tc.ddl)
			if catErr, ok := err.(*Error); !ok || catErr.Code != 1391 {
				t.Fatalf("want *Error 1391, got %v", err)
			}
		})
	}

	// Negative: CREATE INDEX and ALTER ADD reject col(0) too.
	stmtCases := []struct{ name, stmt string }{
		{"create_index", "CREATE INDEX ix ON t (c(0))"},
		{"alter_add_key", "ALTER TABLE t ADD KEY ix (c(0))"},
	}
	for _, tc := range stmtCases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec("CREATE TABLE t (id INT, c VARCHAR(20))", nil)
			r, _ := c.Exec(tc.stmt, &ExecOptions{ContinueOnError: true})
			if catErr, ok := r[0].Error.(*Error); !ok || catErr.Code != 1391 {
				t.Fatalf("want 1391, got %v", r[0].Error)
			}
		})
	}

	// Positive: a real prefix, no prefix, and a bare WITHOUT OVERLAPS stay valid.
	okCases := []struct{ name, ddl string }{
		{"valid_prefix", "CREATE TABLE t (id INT, c VARCHAR(20), KEY (c(5)))"},
		{"no_prefix", "CREATE TABLE t (id INT, c VARCHAR(20), KEY (c))"},
		{"fulltext_valid_prefix", "CREATE TABLE t (t1 TEXT, FULLTEXT KEY k (t1(5)))"},
		{"bare_without_overlaps", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), UNIQUE (id, app_time WITHOUT OVERLAPS))"},
	}
	for _, tc := range okCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := execErr(t, tc.ddl); err != nil {
				t.Fatalf("want success, got %v", err)
			}
		})
	}
}

// TestCreateTableLikePreservesWithoutOverlaps: CREATE TABLE ... LIKE copies the
// application-time period and the WITHOUT OVERLAPS modifier on its key (MariaDB
// 11.8.8 preserves both). createTableLike copied the period but dropped the
// per-key-part WithoutOverlaps flag.
func TestCreateTableLikePreservesWithoutOverlaps(t *testing.T) {
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec("CREATE TABLE src (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), UNIQUE (id, app_time WITHOUT OVERLAPS))", nil)
	c.Exec("CREATE TABLE cpy LIKE src", nil)
	show := c.ShowCreateTable("test", "cpy")
	if !strings.Contains(show, "`app_time` WITHOUT OVERLAPS") {
		t.Fatalf("CREATE TABLE LIKE should preserve WITHOUT OVERLAPS:\n%s", show)
	}
}

// TestKeyPartZeroPrefixSpatialExcluded documents that the zero-length prefix
// check intentionally leaves SPATIAL keys alone. MariaDB rejects a prefix on a
// SPATIAL key part at parse time (1064 for any length), not 1391 — a separate
// pre-existing parser-level gap (omni currently accepts SPATIAL prefixes of any
// length). validateKeyPartPrefixes must NOT report SPATIAL `g(0)` as 1391, which
// would diverge from MariaDB's 1064.
func TestKeyPartZeroPrefixSpatialExcluded(t *testing.T) {
	err := execErr(t, "CREATE TABLE t (g GEOMETRY NOT NULL, SPATIAL KEY k (g(0)))")
	if catErr, ok := err.(*Error); ok && catErr.Code == 1391 {
		t.Fatalf("SPATIAL g(0) must not be reported as 1391 (MariaDB rejects it 1064 at parse); got %v", err)
	}
}

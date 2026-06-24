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

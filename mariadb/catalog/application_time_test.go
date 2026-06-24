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

// TestAlterAddApplicationTimePeriod: ALTER TABLE ... ADD PERIOD FOR <name>(s,e)
// adds an application-time period — same validation codes as CREATE, and the
// period columns become NOT NULL. Grounded vs mariadb:11.8.8.
func TestAlterAddApplicationTimePeriod(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t (id INT, s DATE, e DATE)", nil)
		if r, _ := c.Exec("ALTER TABLE t ADD PERIOD FOR app_time(s, e)", &ExecOptions{ContinueOnError: true}); r[0].Error != nil {
			t.Fatalf("add period: %v", r[0].Error)
		}
		show := c.ShowCreateTable("test", "t")
		if !strings.Contains(show, "PERIOD FOR `app_time` (`s`, `e`)") {
			t.Fatalf("missing period clause:\n%s", show)
		}
		if !strings.Contains(show, "`s` date NOT NULL") || !strings.Contains(show, "`e` date NOT NULL") {
			t.Fatalf("period columns should be NOT NULL:\n%s", show)
		}
	})
	for _, tc := range []struct {
		name, create, add string
		code              int
	}{
		{"already_has_period", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s,e))", "ALTER TABLE t ADD PERIOD FOR p2(s, e)", 4154},
		{"nonexistent_col", "CREATE TABLE t (id INT)", "ALTER TABLE t ADD PERIOD FOR app_time(s, e)", 1054},
		{"non_temporal", "CREATE TABLE t (id INT, s INT, e INT)", "ALTER TABLE t ADD PERIOD FOR app_time(s, e)", 1063},
		{"same_col", "CREATE TABLE t (id INT, s DATE)", "ALTER TABLE t ADD PERIOD FOR app_time(s, s)", 1110},
		{"name_is_column", "CREATE TABLE t (id INT, s DATE, e DATE)", "ALTER TABLE t ADD PERIOD FOR id(s, e)", 1060},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.create, nil)
			r, _ := c.Exec(tc.add, &ExecOptions{ContinueOnError: true})
			if catErr, ok := r[0].Error.(*Error); !ok || catErr.Code != tc.code {
				t.Fatalf("want %d, got %v", tc.code, r[0].Error)
			}
		})
	}
}

// TestAlterDropApplicationTimePeriod: ALTER TABLE ... DROP PERIOD FOR <name>
// removes the application-time period (1091 if absent; 4156 if a WITHOUT
// OVERLAPS key still references it). Grounded vs mariadb:11.8.8.
func TestAlterDropApplicationTimePeriod(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		c := New()
		c.Exec("CREATE DATABASE test", nil)
		c.SetCurrentDatabase("test")
		c.Exec("CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s,e))", nil)
		if r, _ := c.Exec("ALTER TABLE t DROP PERIOD FOR app_time", &ExecOptions{ContinueOnError: true}); r[0].Error != nil {
			t.Fatalf("drop period: %v", r[0].Error)
		}
		if show := c.ShowCreateTable("test", "t"); strings.Contains(show, "PERIOD FOR") {
			t.Fatalf("period should be gone:\n%s", show)
		}
	})
	for _, tc := range []struct {
		name, create, drop string
		code               int
	}{
		{"nonexistent", "CREATE TABLE t (id INT)", "ALTER TABLE t DROP PERIOD FOR app_time", 1091},
		{"wrong_name", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s,e))", "ALTER TABLE t DROP PERIOD FOR other", 1091},
		{"used_by_without_overlaps", "CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s,e), UNIQUE(id, app_time WITHOUT OVERLAPS))", "ALTER TABLE t DROP PERIOD FOR app_time", 4156},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.create, nil)
			r, _ := c.Exec(tc.drop, &ExecOptions{ContinueOnError: true})
			if catErr, ok := r[0].Error.(*Error); !ok || catErr.Code != tc.code {
				t.Fatalf("want %d, got %v", tc.code, r[0].Error)
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

// TestWithoutOverlapsOrdinaryParts: in a WITHOUT OVERLAPS key, the ordinary
// (non-WO) parts must be real columns (else 1072) that are NOT application-time
// period columns (else 4170). A system-time period column is allowed as an
// ordinary part (no 4170). Grounded vs mariadb:11.8.8.
func TestWithoutOverlapsOrdinaryParts(t *testing.T) {
	base := "CREATE TABLE t (id INT, x INT, s DATE, e DATE, PERIOD FOR app_time(s,e), %s)"
	bad := []struct {
		name, key string
		code      int
	}{
		{"period_start_col", "UNIQUE (s, app_time WITHOUT OVERLAPS)", 4170},
		{"period_end_col", "UNIQUE (e, app_time WITHOUT OVERLAPS)", 4170},
		{"period_col_after_valid", "UNIQUE (id, s, app_time WITHOUT OVERLAPS)", 4170},
		{"period_name_as_ordinary", "UNIQUE (app_time, app_time WITHOUT OVERLAPS)", 1072},
		{"nonexistent_ordinary", "UNIQUE (nope, app_time WITHOUT OVERLAPS)", 1072},
	}
	for _, tc := range bad {
		t.Run(tc.name, func(t *testing.T) {
			err := execErr(t, strings.Replace(base, "%s", tc.key, 1))
			if catErr, ok := err.(*Error); !ok || catErr.Code != tc.code {
				t.Fatalf("want *Error %d, got %v", tc.code, err)
			}
		})
	}
	for _, key := range []string{
		"UNIQUE (id, app_time WITHOUT OVERLAPS)",
		"UNIQUE (id, x, app_time WITHOUT OVERLAPS)",
	} {
		t.Run("ok_"+key, func(t *testing.T) {
			if err := execErr(t, strings.Replace(base, "%s", key, 1)); err != nil {
				t.Fatalf("want success, got %v", err)
			}
		})
	}
}

// TestWithoutOverlapsOrdinaryPartsAllPaths: the ordinary-part rules apply on
// CREATE INDEX and ALTER ADD, not only CREATE TABLE.
func TestWithoutOverlapsOrdinaryPartsAllPaths(t *testing.T) {
	for _, tc := range []struct {
		name, stmt string
		code       int
	}{
		{"create_index_period_col", "CREATE UNIQUE INDEX ux ON t (s, app_time WITHOUT OVERLAPS)", 4170},
		{"alter_add_period_col", "ALTER TABLE t ADD UNIQUE KEY ux (s, app_time WITHOUT OVERLAPS)", 4170},
		{"create_index_nonexistent", "CREATE UNIQUE INDEX ux ON t (nope, app_time WITHOUT OVERLAPS)", 1072},
		{"alter_add_nonexistent", "ALTER TABLE t ADD UNIQUE KEY ux (nope, app_time WITHOUT OVERLAPS)", 1072},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec("CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e))", nil)
			r, _ := c.Exec(tc.stmt, &ExecOptions{ContinueOnError: true})
			if catErr, ok := r[0].Error.(*Error); !ok || catErr.Code != tc.code {
				t.Fatalf("want %d, got %v", tc.code, r[0].Error)
			}
		})
	}
}

// TestDropColumnRemovesOrphanedWithoutOverlapsKey: dropping the sole ordinary
// column of a WITHOUT OVERLAPS key removes the whole key + constraint (MariaDB
// 11.8.8 leaves no period-only key). Regular keys are unaffected: a sole column
// drop removes the key; a composite key shrinks.
func TestDropColumnRemovesOrphanedWithoutOverlapsKey(t *testing.T) {
	for _, tc := range []struct {
		name, create, drop, wantAbsent, wantPresent string
	}{
		{
			"wo_single_ordinary",
			"CREATE TABLE t (id INT, s DATE, e DATE, PERIOD FOR app_time(s, e), UNIQUE (id, app_time WITHOUT OVERLAPS))",
			"ALTER TABLE t DROP COLUMN id", "WITHOUT OVERLAPS", "",
		},
		{
			"regular_single_col",
			"CREATE TABLE t (id INT, y INT, UNIQUE (id))",
			"ALTER TABLE t DROP COLUMN id", "UNIQUE KEY", "",
		},
		{
			"regular_composite_shrinks",
			"CREATE TABLE t (id INT, x INT, UNIQUE (id, x))",
			"ALTER TABLE t DROP COLUMN id", "", "`x`",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			c.Exec("CREATE DATABASE test", nil)
			c.SetCurrentDatabase("test")
			c.Exec(tc.create, nil)
			if r, _ := c.Exec(tc.drop, &ExecOptions{ContinueOnError: true}); r[0].Error != nil {
				t.Fatalf("drop column: %v", r[0].Error)
			}
			show := c.ShowCreateTable("test", "t")
			if tc.wantAbsent != "" && strings.Contains(show, tc.wantAbsent) {
				t.Fatalf("%q should be absent after drop:\n%s", tc.wantAbsent, show)
			}
			if tc.wantPresent != "" && !strings.Contains(show, tc.wantPresent) {
				t.Fatalf("%q should be present after drop:\n%s", tc.wantPresent, show)
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

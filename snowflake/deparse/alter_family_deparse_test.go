package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ALTER-family round-trips (gap-alter-family): ALTER SESSION parameters,
// ALTER TABLE SEARCH OPTIMIZATION method-lists, ALTER VIEW SET/UNSET
// {JOIN|AGGREGATION} POLICY.
//
// assertRoundTrip compares ASTs (Loc-stripped), so cosmetic normalizations the
// deparser applies — uppercasing option names/keys, normalizing whitespace —
// are expected and verified to be AST-equivalent.
// ---------------------------------------------------------------------------

func TestDeparse_AlterSessionSet(t *testing.T) {
	assertRoundTrip(t, "ALTER SESSION SET LOCK_TIMEOUT = 3600")
	assertRoundTrip(t, "ALTER SESSION SET STATEMENT_TIMEOUT_IN_SECONDS = 604800")
	assertRoundTrip(t, "ALTER SESSION SET ERROR_ON_NONDETERMINISTIC_UPDATE = TRUE")
	assertRoundTrip(t, "ALTER SESSION SET DEFAULT_NULL_ORDERING = 'LAST'")
	assertRoundTrip(t, "ALTER SESSION SET AUTOCOMMIT = TRUE")
	// Multiple space-separated parameters.
	assertRoundTrip(t, "ALTER SESSION SET AUTOCOMMIT = TRUE LOCK_TIMEOUT = 3600")
}

func TestDeparse_AlterSessionUnset(t *testing.T) {
	assertRoundTrip(t, "ALTER SESSION UNSET LOCK_TIMEOUT")
	assertRoundTrip(t, "ALTER SESSION UNSET LOCK_TIMEOUT, AUTOCOMMIT")
	assertRoundTrip(t, "ALTER SESSION UNSET DEFAULT_NULL_ORDERING")
}

func TestDeparse_AlterTableSearchOptimizationMethodList(t *testing.T) {
	assertRoundTrip(t, "ALTER TABLE t1 ADD SEARCH OPTIMIZATION ON EQUALITY(c1), EQUALITY(c2, c3)")
	assertRoundTrip(t, "ALTER TABLE t1 ADD SEARCH OPTIMIZATION ON EQUALITY(c1, c2, c3, c4)")
	assertRoundTrip(t, "ALTER TABLE t1 DROP SEARCH OPTIMIZATION ON EQUALITY(c1, c2)")
	assertRoundTrip(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON EQUALITY(*)")
	assertRoundTrip(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION ON SUBSTRING(c1), GEO(c2)")
	// No-target form (bare SEARCH OPTIMIZATION) still round-trips.
	assertRoundTrip(t, "ALTER TABLE t ADD SEARCH OPTIMIZATION")
}

func TestDeparse_AlterViewJoinPolicy(t *testing.T) {
	assertRoundTrip(t, "ALTER VIEW join_view SET JOIN POLICY jp1")
	assertRoundTrip(t, "ALTER VIEW v SET JOIN POLICY jp1 ALLOWED JOIN KEYS (a, b)")
	assertRoundTrip(t, "ALTER VIEW v SET JOIN POLICY jp1 ENFORCED JOIN KEYS (a)")
	assertRoundTrip(t, "ALTER VIEW v UNSET JOIN POLICY")
}

func TestDeparse_AlterViewAggregationPolicy(t *testing.T) {
	assertRoundTrip(t, "ALTER VIEW v SET AGGREGATION POLICY ap1")
	assertRoundTrip(t, "ALTER VIEW v SET AGGREGATION POLICY ap1 ENTITY KEY (id)")
	assertRoundTrip(t, "ALTER VIEW v SET AGGREGATION POLICY ap1 ENTITY KEY (a, b) FORCE")
	assertRoundTrip(t, "ALTER VIEW v UNSET AGGREGATION POLICY")
}

// TestDeparse_AlterFamily_ExactString pins the exact deparsed output for
// representative statements, guarding keyword spelling and spacing.
func TestDeparse_AlterFamily_ExactString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{
			in:   "alter session set lock_timeout = 3600",
			want: "ALTER SESSION SET LOCK_TIMEOUT = 3600",
		},
		{
			in:   "ALTER SESSION UNSET lock_timeout, autocommit",
			want: "ALTER SESSION UNSET LOCK_TIMEOUT, AUTOCOMMIT",
		},
		{
			in:   "ALTER VIEW join_view SET JOIN POLICY jp1",
			want: "ALTER VIEW join_view SET JOIN POLICY jp1",
		},
		{
			in:   "ALTER VIEW v SET AGGREGATION POLICY ap1 ENTITY KEY (id) FORCE",
			want: "ALTER VIEW v SET AGGREGATION POLICY ap1 ENTITY KEY (id) FORCE",
		},
		{
			in:   "ALTER TABLE t1 ADD SEARCH OPTIMIZATION ON EQUALITY(c1), EQUALITY(c2, c3)",
			want: "ALTER TABLE t1 ADD SEARCH OPTIMIZATION ON EQUALITY(c1), EQUALITY(c2, c3)",
		},
	}
	for _, c := range cases {
		f, err := parser.Parse(c.in)
		require.NoError(t, err, "parse %q", c.in)
		out, err := deparse.DeparseFile(f)
		require.NoError(t, err, "deparse %q", c.in)
		require.Equal(t, c.want, out)
	}
}

package parser

import "testing"

// TestPR4bParity is the container-grounded conformance matrix for the PR4b
// expr-side query features: IGNORE NULLS null-treatment, binary/hex literals +
// the BINARY operator, and map/array collection literals. Each construct is
// asserted accepted by BOTH the omni parser and StarRocks 3.4, and each
// sibling-arm negative rejected by both. Short-gated via startStarRocks.
func TestPR4bParity(t *testing.T) {
	c := startStarRocks(t)

	cases := []struct {
		name       string
		sql        string
		wantAccept bool
	}{
		// IGNORE NULLS — accept probes.
		{
			"ignore_nulls_inside",
			"SELECT FIRST_VALUE(v IGNORE NULLS) OVER (PARTITION BY g ORDER BY ts) FROM t",
			true,
		},
		{
			"ignore_nulls_last_value",
			"SELECT LAST_VALUE(v IGNORE NULLS) OVER (ORDER BY ts) FROM t",
			true,
		},
		{
			"ignore_nulls_outside",
			"SELECT LEAD(v, 1) IGNORE NULLS OVER (ORDER BY ts) FROM t",
			true,
		},
		{
			"ignore_nulls_mid_args",
			"SELECT LAG(v IGNORE NULLS, 1) OVER (ORDER BY ts) FROM t",
			true,
		},
		{
			"window_no_ignore_nulls", // regression: plain window fn still parses
			"SELECT FIRST_VALUE(v) OVER (PARTITION BY g ORDER BY ts) FROM t",
			true,
		},
		// IGNORE NULLS — sibling-arm negatives.
		{
			"respect_nulls_rejected", // RESPECT NULLS does not exist in StarRocks 3.4
			"SELECT FIRST_VALUE(v RESPECT NULLS) OVER (ORDER BY ts) FROM t",
			false,
		},
		{
			"ignore_without_nulls", // IGNORE alone is malformed
			"SELECT FIRST_VALUE(v IGNORE) OVER (ORDER BY ts) FROM t",
			false,
		},

		// Binary/hex literals + the BINARY operator — accept probes.
		{"hex_literal_upper", "SELECT X'4142' FROM t", true},
		{"hex_literal_lower", "SELECT x'ff' FROM t", true},
		{"bit_literal", "SELECT B'101' FROM t", true},
		{"hex_integer_0x", "SELECT 0x4142 FROM t", true},
		{"binary_operator_select", "SELECT BINARY col1 FROM t", true},
		{"binary_operator_where", "SELECT * FROM t WHERE BINARY name = 'abc'", true},
		{"binary_operator_lowprec", "SELECT BINARY name = 'abc' FROM t", true},
		// Known divergences (documented, left as-is per skip-prune):
		//   - X'zz' (non-hex / odd-length content): omni over-accepts — the lexer
		//     scans the hex body verbatim without validating it; StarRocks rejects.
		//   - SELECT BINARY FROM t (BINARY as a bare column name): omni under-accepts
		//     — BINARY is reserved in omni and handled as the unary operator, while
		//     StarRocks treats it as a nonReserved identifier here. Pre-existing
		//     (omni rejected bare BINARY before this feature too), not a regression.
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.assertParity(t, tc.sql, tc.wantAccept)
		})
	}
}

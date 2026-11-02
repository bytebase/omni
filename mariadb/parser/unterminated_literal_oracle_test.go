package parser

import "testing"

// TestUnterminatedLiteralOracle proves against a live MariaDB 11.8 container
// that the guards refuse only what the server refuses: every unterminated token
// is a 1064 syntax error there, and every closed counterpart parses. classify
// reads a non-1064 error as parsed, so a closed statement may still fail for a
// semantic reason such as an unknown column.
func TestUnterminatedLiteralOracle(t *testing.T) {
	o := startMariaDB(t)

	for _, tc := range unterminatedCases {
		t.Run(tc.name, func(t *testing.T) {
			v, code, msg := o.classify(tc.sql)
			if v == vErrored {
				t.Fatalf("ERRORED (infra, not a parse verdict) for %q: %s", tc.sql, msg)
			}
			if v != vRejected {
				t.Errorf("MariaDB parsed %q (code %d: %s), so omni must not reject it", tc.sql, code, msg)
			}
			if _, err := Parse(tc.sql); err == nil {
				t.Errorf("Parse(%q) accepted what MariaDB rejects", tc.sql)
			}
		})
	}

	for _, sql := range terminatedCounterparts {
		t.Run("closed/"+sql, func(t *testing.T) {
			v, code, msg := o.classify(sql)
			if v == vErrored {
				t.Fatalf("ERRORED (infra, not a parse verdict) for %q: %s", sql, msg)
			}
			if v == vRejected {
				t.Errorf("MariaDB rejected %q (code %d: %s), so it is not a valid control", sql, code, msg)
			}
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}

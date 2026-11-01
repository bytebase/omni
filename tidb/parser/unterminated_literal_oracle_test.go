package parser

import "testing"

// TestUnterminatedLiteralOracle proves against a live MySQL 8.0 container that
// the guards refuse only what the server refuses: every unterminated token is a
// 1064 syntax error there, and every closed counterpart parses. checkSyntax
// reads a non-1064 error as parsed, so a closed statement may still fail for a
// semantic reason such as an unknown column.
//
// The oracle is MySQL, not TiDB, so it reads /*T! ... */ as an ordinary block
// comment rather than TiDB SQL. That is enough for these cases: an unterminated
// comment is a syntax error either way, and a closed one is skipped.
func TestUnterminatedLiteralOracle(t *testing.T) {
	o := startParserOracle(t)

	for _, tc := range unterminatedCases {
		t.Run(tc.name, func(t *testing.T) {
			if o.checkSyntax(t, tc.sql) {
				t.Errorf("the oracle parsed %q, so omni must not reject it", tc.sql)
			}
			if _, err := Parse(tc.sql); err == nil {
				t.Errorf("Parse(%q) accepted what the oracle rejects", tc.sql)
			}
		})
	}

	for _, sql := range terminatedCounterparts {
		t.Run("closed/"+sql, func(t *testing.T) {
			if !o.checkSyntax(t, sql) {
				t.Errorf("the oracle rejected %q, so it is not a valid control", sql)
			}
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}

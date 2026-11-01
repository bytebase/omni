package parser

import "testing"

// TestUnterminatedLiteralOracle proves against a live TiDB v8.5.5 container that
// the guards refuse only what the server refuses: every unterminated token is a
// 1064 syntax error there, every closed counterpart parses, and an executable
// comment that never closes is accepted because TiDB drops the opener.
// tidbRejectedSyntax reads a non-1064 error as parsed, so a statement may still
// fail for a semantic reason such as an unknown column.
//
// The oracle is TiDB and not MySQL because /*T! ... */ and its feature gates are
// TiDB's own syntax: MySQL reads them as ordinary block comments and never takes
// the executable-comment branches these cases cover.
func TestUnterminatedLiteralOracle(t *testing.T) {
	tc := startTiDB(t)

	for _, c := range unterminatedCases {
		t.Run(c.name, func(t *testing.T) {
			if !tidbRejectedSyntax(tc.canExecute(c.sql)) {
				t.Errorf("TiDB parsed %q, so omni must not reject it", c.sql)
			}
			if _, err := Parse(c.sql); err == nil {
				t.Errorf("Parse(%q) accepted what TiDB rejects", c.sql)
			}
		})
	}

	for _, sql := range terminatedCounterparts {
		t.Run("closed/"+sql, func(t *testing.T) {
			if tidbRejectedSyntax(tc.canExecute(sql)) {
				t.Errorf("TiDB rejected %q, so it is not a valid control", sql)
			}
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}

	for _, sql := range unclosedExecCases {
		t.Run("unclosed/"+sql, func(t *testing.T) {
			if tidbRejectedSyntax(tc.canExecute(sql)) {
				t.Errorf("TiDB rejected %q, so omni must not accept it", sql)
			}
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}

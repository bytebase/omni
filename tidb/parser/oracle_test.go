package parser

import (
	"testing"
)

// parserOracle asks the package's TiDB server whether it accepts a statement's
// syntax. It is the same server tidb_container_test.go starts for every oracle
// test in this package, pinned to tidbImage.
type parserOracle struct {
	tc *tidbContainer
}

// startParserOracle returns the TiDB oracle, failing the test unless the server
// runs the release the package's verdicts were recorded on.
func startParserOracle(t *testing.T) *parserOracle {
	t.Helper()
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	return &parserOracle{tc: tc}
}

// checkSyntax reports whether TiDB parses sql. It prepares the statement rather
// than running it, so nothing it names has to exist, and reads the answer
// through prepareVerdict: 1064 and the codes TiDB's grammar actions raise are
// refusals, and an error raised after the parse means the statement was read.
//
// Anything that is not a server error — a dropped connection, a timeout — is no
// verdict at all and fails the test. Reporting it as acceptance would pass every
// case that expects the engine to accept, against a server that never answered.
func (o *parserOracle) checkSyntax(t *testing.T, sql string) bool {
	t.Helper()
	accepts, err := o.tc.prepareSyntax(sql)
	if err != nil {
		t.Fatalf("TiDB gave no verdict for %q: %v", sql, err)
	}
	return accepts
}

// assertOracleMatch parses with omni's parser and asks TiDB, and fails unless
// both accept or both reject.
func assertOracleMatch(t *testing.T, o *parserOracle, sql string) {
	t.Helper()

	_, omniErr := Parse(sql)
	omniAccepts := omniErr == nil

	tidbAccepts := o.checkSyntax(t, sql)

	if omniAccepts != tidbAccepts {
		t.Errorf("MISMATCH for %q:\n  omni accepts=%v (err=%v)\n  TiDB accepts=%v",
			sql, omniAccepts, omniErr, tidbAccepts)
	}
}

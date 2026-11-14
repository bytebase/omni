package compat

import (
	"errors"
	"testing"
)

// dorisError is a stand-in for the driver's error: what the server said.
type dorisError struct{ msg string }

func (e dorisError) Error() string { return "Error 1105 (HY000): " + e.msg }

func readDoris(err error) (string, bool) {
	var e dorisError
	if !errors.As(err, &e) {
		return "", false
	}
	return e.msg, true
}

// TestDorisVerdictReadsTheParsersExceptions pins how a verdict reads the
// answer. The messages are those of dyrnq/doris:3.1.4, to a statement sent
// inside a transaction.
func TestDorisVerdictReadsTheParsersExceptions(t *testing.T) {
	for _, c := range []struct {
		name string
		msg  string // "" for a statement that ran
		want string
	}{
		{"ran", "", Accept},
		{"parsed by the ANTLR parser", "errCode = 2, detailMessage = This is in a transaction, only insert, update, delete, commit, rollback is acceptable.", Accept},
		{"parsed by the CUP parser too", "TException, msg: org.apache.thrift.TException: This is in a transaction, only insert, update, delete, commit, rollback is acceptable.", Accept},
		{"mismatched input", "errCode = 2, detailMessage = \nmismatched input ')' expecting EQ(line 1, pos 38)\n", Reject},
		{"no viable alternative", "errCode = 2, detailMessage = \nno viable alternative at input 'SHOW zzz'(line 1, pos 5)\n", Reject},
		{"extraneous input", "errCode = 2, detailMessage = \nextraneous input 'zzz' expecting {<EOF>, ';'}(line 1, pos 35)\n", Reject},
		{"a missing token", "errCode = 2, detailMessage = \nmissing EQ at '\"v\"'(line 1, pos 92)\n", Reject},
		{"the AST builder refuses", "errCode = 2, detailMessage = \nEvery derived table must have its own alias(line 1, pos 13)\n\n== SQL ==\nSELECT 1 FROM(SELECT 1)\n-------------^^^\n", Reject},
		{"the CUP parser refuses", "errCode = 2, detailMessage = Syntax error in line 1:\nCREATE MATERIALIZED VIEW a SELECT CURRENT_DATE\n                           ^\nEncountered: SELECT\nExpected\n", Reject},
		{"the CUP parser fails", "errCode = 2, detailMessage = sql parsing error, please check your sql", Reject},
		{"a check of the AST builder", "errCode = 2, detailMessage = Unknown privilege type a", Accept},
		{"a DML statement ran", "errCode = 2, detailMessage = Table [a] does not exist in database [omni_compat].", Accept},
		{"a lookup", "errCode = 2, detailMessage = Could not find table function a", Accept},
	} {
		t.Run(c.name, func(t *testing.T) {
			var err error
			if c.msg != "" {
				err = dorisError{c.msg}
			}
			got, verr := DorisVerdict(err, readDoris)
			if verr != nil || got != c.want {
				t.Errorf("DorisVerdict(%q) = %q, %v; want %q", c.msg, got, verr, c.want)
			}
		})
	}
}

// TestDorisVerdictIsNoVerdictWithoutAnAnswer: an error the server did not
// send, such as a dropped connection, is not read as either verdict.
func TestDorisVerdictIsNoVerdictWithoutAnAnswer(t *testing.T) {
	dropped := errors.New("invalid connection")
	if got, err := DorisVerdict(dropped, readDoris); err == nil {
		t.Errorf("DorisVerdict(a dropped connection) = %q, want an error", got)
	}
}

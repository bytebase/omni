package compat

import (
	"errors"
	"testing"
)

// TestStarRocksVerdictReadsTheProbe pins how a verdict reads the answer:
// the probe's error is an acceptance, every other the parse's rejection.
// The messages are those of starrocks/allin1-ubuntu:3.4.10, to a statement
// sent behind StarRocksProbe.
func TestStarRocksVerdictReadsTheProbe(t *testing.T) {
	for _, c := range []struct {
		name string
		msg  string
		want string
	}{
		{"parsed, and the probe failed", "Getting analyzing error from line 1, column 7 to line 1, column 25. Detail message: No matching function with signature: omni_compat_probe().", Accept},
		{"a syntax error", "Getting syntax error at line 2, column 10. Detail message: Unexpected input '<EOF>', the most similar input is {a legal identifier}.", Reject},
		{"no viable statement", "Getting syntax error at line 2, column 5. Detail message: No viable statement for input 'SHOW zzz'.", Reject},
		{"a syntax error that quotes the probe's name", "Getting syntax error at line 2, column 5. Detail message: No viable statement for input 'SHOW omni_compat_probe'.", Reject},
		{"the AST builder refuses", "Getting analyzing error. Detail message: not support.", Reject},
		{"the AST builder fails", "Cannot cast com.starrocks.analysis.SlotRef to com.starrocks.analysis.Predicate", Reject},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := StarRocksVerdict(dorisError{c.msg}, readDoris)
			if err != nil || got != c.want {
				t.Errorf("StarRocksVerdict(%q) = %q, %v; want %q", c.msg, got, err, c.want)
			}
		})
	}
}

// TestStarRocksVerdictIsNoVerdictWithoutTheProbesAnswer: a request that ran,
// and an error the server did not send, are not read as either verdict.
func TestStarRocksVerdictIsNoVerdictWithoutTheProbesAnswer(t *testing.T) {
	if got, err := StarRocksVerdict(nil, readDoris); err == nil {
		t.Errorf("StarRocksVerdict(a request that ran) = %q, want an error", got)
	}
	if got, err := StarRocksVerdict(errors.New("invalid connection"), readDoris); err == nil {
		t.Errorf("StarRocksVerdict(a dropped connection) = %q, want an error", got)
	}
}

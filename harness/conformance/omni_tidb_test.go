package main

import "testing"

func TestOmniTiDBVerdict(t *testing.T) {
	v, errMsg := omniTiDBVerdict("SELECT 1")
	if v != VerdictAccept || errMsg != "" {
		t.Errorf("SELECT 1: got %v/%q, want accept", v, errMsg)
	}
	v, errMsg = omniTiDBVerdict("SELECT FROM WHERE")
	if v != VerdictReject || errMsg == "" {
		t.Errorf("garbage: got %v/%q, want reject with message", v, errMsg)
	}
}

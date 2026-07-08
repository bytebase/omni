package main

import "testing"

func TestExtractTiDBFixture(t *testing.T) {
	entries, err := extractTiDBFile("testdata/tidb_fixture.go.txt")
	if err != nil {
		t.Fatal(err)
	}
	var accepts, rejects, skips int
	for _, e := range entries {
		switch {
		case e.SkipReason != "":
			skips++
		case e.Expected == VerdictAccept:
			accepts++
		case e.Expected == VerdictReject:
			rejects++
		}
	}
	if accepts != 10 || rejects != 3 || skips != 4 {
		t.Fatalf("got accepts=%d rejects=%d skips=%d, want 10/3/4", accepts, rejects, skips)
	}
	// Every entry — skips included — must carry test-function provenance.
	fixtureTests := map[string]bool{"TestDMLStmt": true, "TestNonCompositeElement": true, "TestAppendForm": true, "TestErrMsgTable": true}
	for _, e := range entries {
		if !fixtureTests[e.TestName] {
			t.Errorf("test name = %q, want a fixture test function", e.TestName)
		}
	}
}

// TestExtractTiDBFixtureEntries pins exact SQL round-trips, provenance lines,
// and skip reasons for representative fixture rows.
func TestExtractTiDBFixtureEntries(t *testing.T) {
	entries, err := extractTiDBFile("testdata/tidb_fixture.go.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 17 {
		t.Fatalf("got %d entries, want 17", len(entries))
	}
	want := []struct {
		sql        string
		expected   Verdict
		line       int
		skipReason string
	}{
		{"SELECT 1", VerdictAccept, 16, ""},
		{"INSERT INTO t VALUES (1)", VerdictAccept, 17, ""}, // backtick raw string
		{"SELECT FROM WHERE", VerdictReject, 18, ""},
		{"DELETE FROM t", VerdictAccept, 20, ""}, // keyed form
		{"", VerdictNone, 22, "non_literal"},     // buildSQL() call
		{"SELECT 1 + 1", VerdictAccept, 24, ""},  // "a" + "b" concatenation
		// exact bytes: \r\n escapes and escaped quotes resolve
		{"SELECT 'a\r\nb' WHERE x = \"q\"", VerdictAccept, 26, ""},
		// exact bytes: multi-line raw string keeps interior newlines
		{"CREATE TABLE t (\n\ta INT\n)", VerdictAccept, 28, ""},
		{"SELECT 2", VerdictAccept, 37, ""},            // second table
		{"", VerdictNone, 38, "non_composite_element"}, // bare identifier element
		// bare testCase literal in append form: a literal src extracts for real
		{"SELECT 77", VerdictAccept, 46, ""},
		// bare literal with a non-literal src still yields a SKIP row
		{"", VerdictNone, 47, "non_literal"},
		// testErrMsgCase table: nil err means the parse must succeed — accept
		{"SELECT 100", VerdictAccept, 54, ""},
		// non-nil err (errors.New call): parse must fail — reject
		{"select 1/*", VerdictReject, 56, ""},
		// non-literal src in a testErrMsgCase still yields a SKIP row
		{"", VerdictNone, 58, "non_literal"},
		// keyed form with an ErrXxx selector: non-nil — reject
		{"SELECT FROM ERRMSG", VerdictReject, 60, ""},
		// keyed form with err omitted: the zero value is nil — accept
		{"SELECT 102", VerdictAccept, 62, ""},
	}
	for i, w := range want {
		e := entries[i]
		if e.SQL != w.sql || e.Expected != w.expected || e.Line != w.line || e.SkipReason != w.skipReason {
			t.Errorf("entry %d = {SQL:%q Expected:%q Line:%d SkipReason:%q}, want {SQL:%q Expected:%q Line:%d SkipReason:%q}",
				i, e.SQL, e.Expected, e.Line, e.SkipReason, w.sql, w.expected, w.line, w.skipReason)
		}
		if e.SourcePath != "testdata/tidb_fixture.go.txt" {
			t.Errorf("entry %d source path = %q", i, e.SourcePath)
		}
	}
}

// TestLiteralResolversNilSafe pins that the resolvers reject a nil expr
// (missing keyed field / short positional literal) instead of panicking.
func TestLiteralResolversNilSafe(t *testing.T) {
	if _, ok := stringValue(nil); ok {
		t.Error("stringValue(nil) ok = true, want false")
	}
	if _, ok := boolValue(nil); ok {
		t.Error("boolValue(nil) ok = true, want false")
	}
}

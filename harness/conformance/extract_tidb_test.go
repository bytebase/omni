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
	if accepts != 7 || rejects != 1 || skips != 2 {
		t.Fatalf("got accepts=%d rejects=%d skips=%d, want 7/1/2", accepts, rejects, skips)
	}
	// Every entry — skips included — must carry test-function provenance.
	for _, e := range entries {
		if e.TestName != "TestDMLStmt" && e.TestName != "TestNonCompositeElement" {
			t.Errorf("test name = %q, want TestDMLStmt or TestNonCompositeElement", e.TestName)
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
	if len(entries) != 10 {
		t.Fatalf("got %d entries, want 10", len(entries))
	}
	want := []struct {
		sql        string
		expected   Verdict
		line       int
		skipReason string
	}{
		{"SELECT 1", VerdictAccept, 11, ""},
		{"INSERT INTO t VALUES (1)", VerdictAccept, 12, ""}, // backtick raw string
		{"SELECT FROM WHERE", VerdictReject, 13, ""},
		{"DELETE FROM t", VerdictAccept, 15, ""}, // keyed form
		{"", VerdictNone, 17, "non_literal"},     // buildSQL() call
		{"SELECT 1 + 1", VerdictAccept, 19, ""},  // "a" + "b" concatenation
		// exact bytes: \r\n escapes and escaped quotes resolve
		{"SELECT 'a\r\nb' WHERE x = \"q\"", VerdictAccept, 21, ""},
		// exact bytes: multi-line raw string keeps interior newlines
		{"CREATE TABLE t (\n\ta INT\n)", VerdictAccept, 23, ""},
		{"SELECT 2", VerdictAccept, 32, ""},            // second table
		{"", VerdictNone, 33, "non_composite_element"}, // bare identifier element
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

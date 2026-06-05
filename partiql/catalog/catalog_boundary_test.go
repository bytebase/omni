package catalog

import (
	"reflect"
	"testing"
)

// TestAddTableEmptyName documents the current contract for an empty table
// name. The catalog is a thin wrapper over a map[string]struct{}; the empty
// string is a valid map key, so AddTable("") registers a real (if unusual)
// entry rather than being silently dropped. Callers that want to reject empty
// names must do so before reaching the catalog.
func TestAddTableEmptyName(t *testing.T) {
	c := New()

	c.AddTable("")
	if !c.HasTable("") {
		t.Error(`HasTable("") = false after AddTable(""), want true (empty string is a stored key)`)
	}

	tables := c.Tables()
	if len(tables) != 1 {
		t.Fatalf("Tables() length = %d after AddTable(\"\"), want 1; got %#v", len(tables), tables)
	}
	if tables[0] != "" {
		t.Errorf("Tables()[0] = %q, want \"\"", tables[0])
	}

	// The empty entry must not be confused with "no entry": a fresh catalog
	// reports HasTable("") == false.
	if New().HasTable("") {
		t.Error(`HasTable("") = true on a fresh catalog, want false`)
	}
}

// TestAddTableQuotedAndSpecialNames documents that table names are stored
// verbatim as opaque strings. The catalog performs no unquoting, trimming,
// case-folding, or normalization, so quoted identifiers, embedded quotes,
// whitespace, dots, and Unicode all round-trip exactly as supplied.
//
// (DynamoDB table names themselves are restricted to a narrow ASCII set, but
// the catalog API does not enforce that — it simply remembers whatever string
// it is handed. This test pins that "no normalization" behavior.)
func TestAddTableQuotedAndSpecialNames(t *testing.T) {
	names := []string{
		`"Music"`,        // double-quoted identifier kept verbatim (quotes included)
		`'Albums'`,       // single-quoted string kept verbatim
		"`Backtick`",     // backtick-quoted name kept verbatim
		"my table",       // embedded space
		"schema.Table",   // dotted / qualified name
		"weird-name_123", // hyphen + underscore + digits
		"emoji_🎵",        // multibyte Unicode
		"tab\tname",      // embedded tab
		"Mixed\"Quote",   // embedded double quote inside an unquoted name
	}

	c := New()
	for _, n := range names {
		c.AddTable(n)
	}

	for _, n := range names {
		if !c.HasTable(n) {
			t.Errorf("HasTable(%q) = false after AddTable(%q), want true (names stored verbatim)", n, n)
		}
	}

	if got := len(c.Tables()); got != len(names) {
		t.Errorf("Tables() length = %d, want %d (each distinct verbatim name is one entry); got %#v", got, len(names), c.Tables())
	}

	// A quoted form and its unquoted form are DISTINCT keys: the catalog does
	// not unquote. "Music" (with quotes) != Music (without).
	if c.HasTable("Music") {
		t.Error(`HasTable("Music") (unquoted) = true, want false: catalog stored only the quoted form "\"Music\"" and does not unquote`)
	}
}

// TestTablesOrderingAndDedup pins the Tables() contract: results are sorted by
// Go's byte-wise string order and duplicates collapse to a single entry,
// independent of insertion order.
func TestTablesOrderingAndDedup(t *testing.T) {
	c := New()
	// Insert out of order, with duplicates and case variants.
	for _, n := range []string{"Music", "albums", "Albums", "music", "Music", "Zebra", "Apple"} {
		c.AddTable(n)
	}

	got := c.Tables()

	// Expected: deduped set, sorted by byte value. Note byte-wise sort puts all
	// uppercase letters (A-Z, 0x41-0x5A) before all lowercase (a-z, 0x61-0x7A),
	// and case variants are distinct keys (see HasTable case-sensitivity test).
	want := []string{"Albums", "Apple", "Music", "Zebra", "albums", "music"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Tables() = %#v, want %#v", got, want)
	}

	// Tables() must return a fresh slice each call (callers may mutate without
	// corrupting catalog state). Mutating the returned slice must not affect a
	// subsequent call.
	got[0] = "MUTATED"
	if again := c.Tables(); reflect.DeepEqual(again, got) {
		t.Errorf("Tables() returned a slice aliased to internal state: after caller mutation, second call = %#v", again)
	}
}

// TestHasTableCaseSensitivity DOCUMENTS the current (U3-ruling) behavior:
// catalog HasTable is a plain map lookup and is therefore CASE-SENSITIVE.
// "Music" and "music" are independent entries.
//
// This deliberately CONTRASTS with PartiQL completion, whose matchesPrefix
// upper-cases both sides (partiql/completion/completion.go: matchesPrefix uses
// strings.HasPrefix(strings.ToUpper(text), strings.ToUpper(prefix))) and is
// therefore CASE-INSENSITIVE. Per the U3 ruling this divergence is left AS-IS;
// the catalog is the exact-match metadata store, while completion is the
// fuzzy/UX layer. This test exists to lock that intentional asymmetry in place
// so a future refactor that "unifies" the two is caught.
func TestHasTableCaseSensitivity(t *testing.T) {
	c := New()
	c.AddTable("Music")

	if !c.HasTable("Music") {
		t.Error(`HasTable("Music") = false, want true (exact match)`)
	}

	// Case-sensitive: a differently-cased query does NOT match.
	for _, miss := range []string{"music", "MUSIC", "MuSiC"} {
		if c.HasTable(miss) {
			t.Errorf("HasTable(%q) = true, want false: catalog HasTable is case-sensitive (exact map lookup); only %q is registered", miss, "Music")
		}
	}

	// Adding the lowercased variant creates a SECOND, independent entry.
	c.AddTable("music")
	if !c.HasTable("music") {
		t.Error(`HasTable("music") = false after AddTable("music"), want true`)
	}
	if got := len(c.Tables()); got != 2 {
		t.Errorf("Tables() length = %d after registering both \"Music\" and \"music\", want 2 (case variants are distinct keys); got %#v", got, c.Tables())
	}
}

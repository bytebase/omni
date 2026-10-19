package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// identFoldCases pairs an identifier as written in SQL with the name
// PostgreSQL 17 stores for it in a UTF-8 database. Unquoted identifiers fold
// only A-Z (downcase_identifier, scansup.c); every identifier, quoted or not,
// is cut to at most 63 bytes on a character boundary (truncate_identifier
// through pg_mbcliplen). TestIdentifierFoldingMatchesPG checks the expected
// names against a real server.
var identFoldCases = []struct {
	name string
	sql  string
	want string
}{
	{"ascii upper", "MixedCase", "mixedcase"},
	{"latin-1 upper keeps case", "Äb", "Äb"},
	{"accented word", "ÉCOLE", "École"},
	{"titlecase digraph", "ǅx", "ǅx"},
	{"turkish dotted capital I", "İSTANBUL", "İstanbul"},
	{"kelvin sign is not keyword KEY", "KEY", "Key"},
	{"long s is not keyword SELECT", "ſELECT", "ſelect"},
	{"quoted keeps case", `"MixedCase"`, "MixedCase"},
	{"ascii truncated to 63", strings.Repeat("A", 70), strings.Repeat("a", 63)},
	{"two-byte char straddles 63", strings.Repeat("a", 62) + "Äxyz", strings.Repeat("a", 62)},
	{"two-byte char ends at 63", strings.Repeat("a", 61) + "Äxyz", strings.Repeat("a", 61) + "Ä"},
	{"three-byte char ends at 63", strings.Repeat("a", 60) + "€x", strings.Repeat("a", 60) + "€"},
	{"three-byte char straddles 63", strings.Repeat("a", 61) + "€x", strings.Repeat("a", 61)},
	{"four-byte char straddles 63", strings.Repeat("a", 60) + "😀x", strings.Repeat("a", 60)},
	{"exactly 63 bytes kept", strings.Repeat("a", 63), strings.Repeat("a", 63)},
	{"quoted truncated", `"` + strings.Repeat("B", 62) + `Éxyz"`, strings.Repeat("B", 62)},
	{"quoted ascii truncated", `"` + strings.Repeat("B", 64) + `"`, strings.Repeat("B", 63)},
	{"unicode escape truncated after decoding", `U&"` + strings.Repeat("c", 62) + `\00C9xyz"`, strings.Repeat("c", 62)},
	{"unicode escape fits after decoding", `U&"` + strings.Repeat("c", 61) + `\00C9xyz"`, strings.Repeat("c", 61) + "É"},
	{"unicode escape with UESCAPE", `U&"d!00C9" UESCAPE '!'`, "dÉ"},
}

// identFoldName parses "SELECT 1 AS <ident>" and returns the alias.
func identFoldName(t *testing.T, ident string) string {
	t.Helper()
	sql := "SELECT 1 AS " + ident
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	sel := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	return sel.TargetList.Items[0].(*nodes.ResTarget).Name
}

func TestIdentifierFolding(t *testing.T) {
	for _, tc := range identFoldCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := identFoldName(t, tc.sql); got != tc.want {
				t.Errorf("alias %q: got %q, want %q", tc.sql, got, tc.want)
			}
		})
	}
}

// TestIdentifierFoldingRelation checks that a relation name folds the same
// way as a column alias, through qualified names and CREATE TABLE.
func TestIdentifierFoldingRelation(t *testing.T) {
	list, err := Parse(`CREATE TABLE ÄB.ÉCOLE (a int)`)
	if err != nil {
		t.Fatal(err)
	}
	rv := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateStmt).Relation
	if rv.Schemaname != "Äb" || rv.Relname != "École" {
		t.Errorf("got %q.%q, want %q.%q", rv.Schemaname, rv.Relname, "Äb", "École")
	}
}

// TestLookupKeywordASCIIOnly pins that keyword lookup folds only ASCII
// letters, as ScanKeywordLookup (src/common/kwlookup.c) does.
func TestLookupKeywordASCIIOnly(t *testing.T) {
	if kw := LookupKeyword("SeLeCt"); kw == nil || kw.Token != SELECT {
		t.Errorf("LookupKeyword(SeLeCt) = %v, want SELECT", kw)
	}
	for _, name := range []string{"Key", "ſelect", "İn"} {
		if kw := LookupKeyword(name); kw != nil {
			t.Errorf("LookupKeyword(%q) = %q, want nil", name, kw.Name)
		}
	}
}

// TestIdentifierFoldingMatchesPG asks PostgreSQL 17 for the name it gives
// each identifier, as the column name of SELECT 1 AS <ident>.
func TestIdentifierFoldingMatchesPG(t *testing.T) {
	o := startFirstSetOracle(t)
	for _, tc := range identFoldCases {
		t.Run(tc.name, func(t *testing.T) {
			rows, err := o.db.QueryContext(o.ctx, "SELECT 1 AS "+tc.sql)
			if err != nil {
				t.Fatalf("PG rejected alias %q: %v", tc.sql, err)
			}
			cols, err := rows.Columns()
			rows.Close()
			if err != nil {
				t.Fatal(err)
			}
			if len(cols) != 1 || cols[0] != tc.want {
				t.Errorf("PG named alias %q %q, want %q", tc.sql, cols, tc.want)
			}
		})
	}
}

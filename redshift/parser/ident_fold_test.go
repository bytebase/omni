package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestIdentifierFolding pins how an unquoted identifier folds: only the ASCII
// letters A-Z are lowercased, as in PostgreSQL's downcase_identifier and as
// Redshift documents ("ASCII letters in standard and delimited identifiers
// are case-insensitive and are folded to lowercase",
// https://docs.aws.amazon.com/redshift/latest/dg/r_names.html), and a long
// identifier is cut on a character boundary. Delimited identifiers are left
// as written; see namedatalen.
func TestIdentifierFolding(t *testing.T) {
	keep := namedatalen - 1
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{"ascii upper", "MixedCase", "mixedcase"},
		{"latin-1 upper keeps case", "Äb", "Äb"},
		{"accented word", "ÉCOLE", "École"},
		{"turkish dotted capital I", "İSTANBUL", "İstanbul"},
		{"kelvin sign is not keyword KEY", "KEY", "Key"},
		{"quoted keeps case", `"MixedCase"`, "MixedCase"},
		{"ascii truncated", strings.Repeat("A", keep+7), strings.Repeat("a", keep)},
		{"exactly the limit kept", strings.Repeat("a", keep), strings.Repeat("a", keep)},
		{"two-byte char straddles the limit", strings.Repeat("a", keep-1) + "Äxyz", strings.Repeat("a", keep-1)},
		{"two-byte char ends at the limit", strings.Repeat("a", keep-2) + "Äxyz", strings.Repeat("a", keep-2) + "Ä"},
		{"four-byte char straddles the limit", strings.Repeat("a", keep-3) + "😀x", strings.Repeat("a", keep-3)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sql := "SELECT 1 AS " + tc.sql
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
			sel := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
			if got := sel.TargetList.Items[0].(*nodes.ResTarget).Name; got != tc.want {
				t.Errorf("alias %q: got %q, want %q", tc.sql, got, tc.want)
			}
		})
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

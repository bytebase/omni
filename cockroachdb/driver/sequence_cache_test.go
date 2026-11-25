package driver

import "testing"

// TestSequenceCacheSkipsQuotedNames reads CACHE from the options, never from a
// quoted name holding the word.
func TestSequenceCacheSkipsQuotedNames(t *testing.T) {
	for statement, want := range map[string]string{
		`CREATE SEQUENCE public."orders CACHE 4" MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1 CACHE 8`: "8",
		`CREATE SEQUENCE public."orders CACHE 4" MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1`:         "",
		`CREATE SEQUENCE s.cached MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1 PER NODE CACHE 16`:       "",
	} {
		got := ""
		if m := sequenceCache.FindStringSubmatch(quotedName.ReplaceAllString(statement, `""`)); m != nil {
			got = m[1]
		}
		if got != want {
			t.Errorf("%s: cache %q, want %q", statement, got, want)
		}
	}
}

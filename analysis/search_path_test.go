package analysis

import (
	"slices"
	"testing"
)

func TestSearchPath(t *testing.T) {
	for _, tc := range []struct {
		setting string
		want    []string
	}{
		{"", []string{}},
		{`"$user", public`, []string{"public"}},
		{`"$user"`, nil},
		{`pg_catalog, "My""Schema" ,  Sales,,information_schema`, []string{`My"Schema`, "sales"}},
		{"pg_temp_1, pg_temp_2, pg_toast, pg_global", []string{}},
	} {
		if got := SearchPath(tc.setting); !slices.Equal(got, tc.want) || (got == nil) != (tc.want == nil) {
			t.Errorf("SearchPath(%q) = %#v, want %#v", tc.setting, got, tc.want)
		}
	}
}

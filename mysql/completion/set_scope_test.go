package completion

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mysql/catalog"
)

// TestSetScopeCompletion checks a later assignment offers its scope words, and that a word
// past the cursor does not make what comes before it a variable name.
func TestSetScopeCompletion(t *testing.T) {
	cat := catalog.New()
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"set @a = 1, <|>", []string{"GLOBAL", "SESSION"}},
		{"set a = 1, b = 2, <|>", []string{"GLOBAL", "SESSION"}},
		{"set transaction <|> = 1", []string{"ISOLATION", "READ"}},
		{"set @a = 1, names <|>", []string{"utf8mb4"}},
		{"set @a = 1, charset <|>", []string{"utf8mb4"}},
		{"set @a = 1, character set <|>", []string{"utf8mb4"}},
		{"set char set <|>", []string{"utf8mb4"}},
		{"set names utf8, <|>", []string{"GLOBAL", "SESSION"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			sql, cursor := splitCursor(t, tc.sql)
			var got []string
			for _, c := range Complete(sql, cursor, cat) {
				got = append(got, c.Text)
			}
			for _, want := range tc.want {
				if !slices.Contains(got, want) {
					t.Errorf("candidates = %v, want %s among them", got, want)
				}
			}
		})
	}
}
